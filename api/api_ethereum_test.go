package api

import (
	"context"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/klaytn/klaytn/crypto"

	"github.com/klaytn/klaytn/rlp"

	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/params"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/mocks"

	"github.com/golang/mock/gomock"
	mock_api "github.com/klaytn/klaytn/api/mocks"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/stretchr/testify/assert"
)

// TestEthereumAPI_GetHeaderByNumber tests GetHeaderByNumber.
func TestEthereumAPI_GetHeaderByNumber(t *testing.T) {
	testGetHeader(t, "GetHeaderByNumber")
}

// TestEthereumAPI_GetHeaderByHash tests GetHeaderByNumber.
func TestEthereumAPI_GetHeaderByHash(t *testing.T) {
	testGetHeader(t, "GetHeaderByHash")
}

// testGetHeader generates data to test GetHeader related functions in EthereumAPI
// and actually tests the API function passed as a parameter.
func testGetHeader(t *testing.T, testAPIName string) {
	mockCtrl := gomock.NewController(t)
	mockBackend := mock_api.NewMockBackend(mockCtrl)

	api := EthereumAPI{
		publicKlayAPI:       NewPublicKlayAPI(mockBackend),
		publicBlockChainAPI: NewPublicBlockChainAPI(mockBackend),
	}

	// Creates a MockEngine.
	mockEngine := mocks.NewMockEngine(mockCtrl)
	// GetHeader APIs calls internally below methods.
	mockBackend.EXPECT().Engine().Return(mockEngine)
	// Author is called when calculates miner field of Header.
	mockEngine.EXPECT().Author(gomock.Any()).Return(common.Address{}, nil)
	mockBackend.EXPECT().GetTd(gomock.Any()).Return(new(big.Int))

	header := types.CopyHeader(&types.Header{}) // Creates an empty header
	var blockParam interface{}
	switch testAPIName {
	case "GetHeaderByNumber":
		blockParam = rpc.BlockNumber(header.Number.Uint64())
		mockBackend.EXPECT().HeaderByNumber(gomock.Any(), gomock.Any()).Return(header, nil)
	case "GetHeaderByHash":
		blockParam = header.Hash()
		mockBackend.EXPECT().HeaderByHash(gomock.Any(), gomock.Any()).Return(header, nil)
	}

	results := reflect.ValueOf(&api).MethodByName(testAPIName).Call(
		[]reflect.Value{
			reflect.ValueOf(context.Background()),
			reflect.ValueOf(blockParam),
		},
	)
	ethHeader, ok := results[0].Interface().(map[string]interface{})
	assert.Equal(t, true, ok)
	assert.NotEqual(t, ethHeader, nil)

	expected := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal([]byte(`
    {
        "jsonrpc": "2.0",
        "id": 1,
        "result": {
            "baseFeePerGas": "0xf79eb1429",
            "difficulty": "0x29daf6647fd2b7",
            "extraData": "0x75732d656173742d37",
            "gasLimit": "0x1c9c32b",
            "gasUsed": "0x28b484",
            "hash": "0x5de0dc71dec2e724be002dcad135b602810769ce26e16b3b06862405e08ca71b",
            "logsBloom": "0x02200022800002050000084080014015001001004b0002440401060a0830000200014041044010180010430018800119120098000800200241c2090a4020011040004400002201081800440a340020a4000820100848081020003000892050105a05000002100000200012c0800408982000085100000c4040a03814000800200812210100200010004018410d80004214800123210400082002214620100021028800120309200802008291c8e000904210080008110900010100081000101000501002010a0080311886000008000000240900400000100200a402005830200010300804020200000002310000008004004080a58000550000508000000000",
            "miner": "0xea674fdde714fd979de3edf0f56aa9716b898ec8",
            "mixHash": "0x6d266344754999c95ad189f78257d31c276c63c689d864c31fdc62fcb481e5f0",
            "nonce": "0x8b232816a7045c51",
            "number": "0xd208de",
            "parentHash": "0x99fcd33dddd763835ba8bdc842853d973496a7e64ea2f6cf826bc2c338e23b0c",
            "receiptsRoot": "0xd3d70ed54a9274ba3191bf2d4fd8738c5d782fe17c8bfb45c03a25dc04120c35",
            "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
            "size": "0x23a",
            "stateRoot": "0x1076e6726164bd6f74720a717242584109f37c55017d004eefccf9ec3be76c18",
            "timestamp": "0x61b0a6c6",
            "totalDifficulty": "0x7a58841ac2bdc5d1e97",
            "transactionsRoot": "0x6ec8daca98c1005d9bbd7716b5e94180e2bf0e6b77770174563a166337369344"
        }
    }   
    `), &expected))
	checkEthereumBlockOrHeaderFormat(t, expected, ethHeader)
}

// checkEthereumHeaderFormat checks that ethBlockOrHeader returned from
// GetHeader or GetBlock APIs have all keys of the Ethereum data structure and
// also checks the immutable values are well-defined or not.
func checkEthereumBlockOrHeaderFormat(
	t *testing.T,
	expected map[string]interface{},
	ethBlockOrHeader map[string]interface{},
) {
	expectedResult, ok := expected["result"].(map[string]interface{})
	assert.True(t, ok)

	// Check immutable values.
	// these values only exists for supporting Ethereum compatible data structures.
	// 1. Check nonce
	nonce := ethBlockOrHeader["nonce"]
	assert.Equal(t, BlockNonce{}, nonce)

	mixHash := ethBlockOrHeader["mixHash"]
	assert.Equal(t, common.Hash{}, mixHash)

	emptyUncleHeaders, err := rlp.EncodeToBytes([]*types.Header(nil))
	assert.NoError(t, err)
	emptySha3Uncles := crypto.Keccak256Hash(emptyUncleHeaders)
	sha3Uncles := ethBlockOrHeader["sha3Uncles"]
	assert.Equal(t, common.HexToHash(EmptySha3Uncles), sha3Uncles)
	assert.Equal(t, emptySha3Uncles, sha3Uncles)

	extraData := ethBlockOrHeader["extraData"]
	assert.Equal(t, hexutil.Bytes(DummyExtraData), extraData)

	gasLimit := ethBlockOrHeader["gasLimit"]
	assert.Equal(t, hexutil.Uint64(DummyGasLimit), gasLimit)

	baseFeePerGas := ethBlockOrHeader["baseFeePerGas"]
	assert.Equal(t, (*hexutil.Big)(new(big.Int).SetUint64(params.BaseFee)), baseFeePerGas)

	// returnedKeys must have all keys of the Ethereum header.
	returnedKeys := make([]string, len(ethBlockOrHeader))
	for k := range ethBlockOrHeader {
		returnedKeys = append(returnedKeys, k)
	}
	existed := func(key string) bool {
		for _, ek := range returnedKeys {
			if key == ek {
				return true
			}
		}
		return false
	}
	existsOrNot := false
	for k := range expectedResult {
		existsOrNot = existed(k)
	}
	assert.True(t, existsOrNot)
}
