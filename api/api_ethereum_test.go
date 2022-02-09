package api

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/governance"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/accounts"
	mock_accounts "github.com/klaytn/klaytn/accounts/mocks"
	mock_api "github.com/klaytn/klaytn/api/mocks"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus/mocks"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

// TestEthereumAPI_Etherbase tests Etherbase.
func TestEthereumAPI_Etherbase(t *testing.T) {
	testNodeAddress(t, "Etherbase")
}

// TestEthereumAPI_Coinbase tests Coinbase.
func TestEthereumAPI_Coinbase(t *testing.T) {
	testNodeAddress(t, "Coinbase")
}

// testNodeAddress generates nodeAddress and tests Etherbase and Coinbase.
func testNodeAddress(t *testing.T, testAPIName string) {
	gov := governance.NewGovernance(
		&params.ChainConfig{},
		database.NewMemoryDBManager(),
	)
	key, _ := crypto.GenerateKey()
	nodeAddress := crypto.PubkeyToAddress(key.PublicKey)
	gov.SetNodeAddress(nodeAddress)

	api := EthereumAPI{publicGovernanceAPI: governance.NewGovernanceAPI(gov)}
	results := reflect.ValueOf(&api).MethodByName(testAPIName).Call([]reflect.Value{})
	result, ok := results[0].Interface().(common.Address)
	assert.True(t, ok)
	assert.Equal(t, nodeAddress, result)
}

// TestEthereumAPI_Hashrate tests Hasharate.
func TestEthereumAPI_Hashrate(t *testing.T) {
	api := &EthereumAPI{}
	assert.Equal(t, hexutil.Uint64(ZeroHashrate), api.Hashrate())
}

// TestEthereumAPI_Mining tests Mining.
func TestEthereumAPI_Mining(t *testing.T) {
	api := &EthereumAPI{}
	assert.Equal(t, false, api.Mining())
}

// TestEthereumAPI_GetWork tests GetWork.
func TestEthereumAPI_GetWork(t *testing.T) {
	api := &EthereumAPI{}
	_, err := api.GetWork()
	assert.Equal(t, errNoMiningWork, err)
}

// TestEthereumAPI_SubmitWork tests SubmitWork.
func TestEthereumAPI_SubmitWork(t *testing.T) {
	api := &EthereumAPI{}
	assert.Equal(t, false, api.SubmitWork(BlockNonce{}, common.Hash{}, common.Hash{}))
}

// TestEthereumAPI_SubmitHashrate tests SubmitHashrate.
func TestEthereumAPI_SubmitHashrate(t *testing.T) {
	api := &EthereumAPI{}
	assert.Equal(t, false, api.SubmitHashrate(hexutil.Uint64(0), common.Hash{}))
}

// TestEthereumAPI_GetHashrate tests GetHashrate.
func TestEthereumAPI_GetHashrate(t *testing.T) {
	api := &EthereumAPI{}
	assert.Equal(t, ZeroHashrate, api.GetHashrate())
}

// TestEthereumAPI_GetUncleByBlockNumberAndIndex tests GetUncleByBlockNumberAndIndex.
func TestEthereumAPI_GetUncleByBlockNumberAndIndex(t *testing.T) {
	api := &EthereumAPI{}
	uncleBlock, err := api.GetUncleByBlockNumberAndIndex(context.Background(), rpc.BlockNumber(0), hexutil.Uint(0))
	assert.NoError(t, err)
	assert.Nil(t, uncleBlock)
}

// TestEthereumAPI_GetUncleByBlockHashAndIndex tests GetUncleByBlockHashAndIndex.
func TestEthereumAPI_GetUncleByBlockHashAndIndex(t *testing.T) {
	api := &EthereumAPI{}
	uncleBlock, err := api.GetUncleByBlockHashAndIndex(context.Background(), common.Hash{}, hexutil.Uint(0))
	assert.NoError(t, err)
	assert.Nil(t, uncleBlock)
}

// TestTestEthereumAPI_GetUncleCountByBlockNumber tests GetUncleCountByBlockNumber.
func TestTestEthereumAPI_GetUncleCountByBlockNumber(t *testing.T) {
	api := &EthereumAPI{}
	uncleCount := hexutil.Uint(ZeroUncleCount)
	assert.Equal(t, uncleCount, *api.GetUncleCountByBlockNumber(context.Background(), rpc.BlockNumber(0)))
}

// TestTestEthereumAPI_GetUncleCountByBlockHash tests GetUncleCountByBlockHash.
func TestTestEthereumAPI_GetUncleCountByBlockHash(t *testing.T) {
	api := &EthereumAPI{}
	uncleCount := hexutil.Uint(ZeroUncleCount)
	assert.Equal(t, uncleCount, *api.GetUncleCountByBlockHash(context.Background(), common.Hash{}))
}

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
	dummyMiner := common.HexToAddress("0x9712f943b296758aaae79944ec975884188d3a96")
	mockEngine.EXPECT().Author(gomock.Any()).Return(dummyMiner, nil)
	var dummyTotalDifficulty uint64 = 5
	mockBackend.EXPECT().GetTd(gomock.Any()).Return(new(big.Int).SetUint64(dummyTotalDifficulty))

	// Create dummy header
	header := types.CopyHeader(&types.Header{
		ParentHash:  common.HexToHash("0xc8036293065bacdfce87debec0094a71dbbe40345b078d21dcc47adb4513f348"),
		Rewardbase:  common.Address{},
		TxHash:      common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"),
		Root:        common.HexToHash("0xad31c32942fa033166e4ef588ab973dbe26657c594de4ba98192108becf0fec9"),
		ReceiptHash: common.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"),
		Bloom:       types.Bloom{},
		BlockScore:  new(big.Int).SetUint64(1),
		Number:      new(big.Int).SetUint64(4),
		GasUsed:     uint64(10000),
		Time:        new(big.Int).SetUint64(1641363540),
		TimeFoS:     uint8(85),
		Extra:       common.Hex2Bytes("0xd983010701846b6c617988676f312e31362e338664617277696e000000000000f89ed5949712f943b296758aaae79944ec975884188d3a96b8415a0614be7fd5ea40f11ce558e02993bd55f11ae72a3cfbc861875a57483ec5ec3adda3e5845fd7ab271d670c755480f9ef5b8dd731f4e1f032fff5d165b763ac01f843b8418867d3733167a0c737fa5b62dcc59ec3b0af5748bcc894e7990a0b5a642da4546713c9127b3358cdfe7894df1ca1db5a97560599986d7f1399003cd63660b98200"),
		Governance:  []byte{},
		Vote:        []byte{},
	})

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
		"baseFeePerGas": "0x0",
		"difficulty": "0x1",
		"extraData": "0x",
		"gasLimit": "0x3b9ac9ff",
		"gasUsed": "0x2710",
		"hash": "0xd6d5d8295f612bc824762f1945f4271c73aee9306bcf91e151d269369526ba60",
		"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"miner": "0x9712f943b296758aaae79944ec975884188d3a96",
		"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce": "0x0000000000000000",
		"number": "0x4",
		"parentHash": "0xc8036293065bacdfce87debec0094a71dbbe40345b078d21dcc47adb4513f348",
		"receiptsRoot": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size": "0x20c",
		"stateRoot": "0xad31c32942fa033166e4ef588ab973dbe26657c594de4ba98192108becf0fec9",
		"timestamp": "0x61d53854",
		"totalDifficulty": "0x5",
		"transactionsRoot": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
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
	marshaledExpectedResult, err := json.Marshal(expectedResult)
	assert.NoError(t, err)
	t.Logf("expectedResult: %s\n", marshaledExpectedResult)
	marshaledEthBlockOrHeader, err := json.Marshal(ethBlockOrHeader)
	assert.NoError(t, err)
	t.Logf("actualResult: %s\n", marshaledEthBlockOrHeader)
	assert.Equal(t, marshaledExpectedResult, marshaledEthBlockOrHeader)
}

// TestEthereumAPI_GetTransactionByBlockNumberAndIndex tests GetTransactionByBlockNumberAndIndex.
func TestEthereumAPI_GetTransactionByBlockNumberAndIndex(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, _, _, _ := createTestData(t)

	// Mock Backend functions.
	mockBackend.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil).Times(txs.Len())

	// Get transaction by block number and index for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		ethTx := api.GetTransactionByBlockNumberAndIndex(context.Background(), rpc.BlockNumber(block.NumberU64()), hexutil.Uint(i))
		checkEthRPCTransactionFormat(t, block, ethTx, txs[i], hexutil.Uint64(i))
	}

	mockCtrl.Finish()
}

// TestEthereumAPI_GetTransactionByBlockHashAndIndex tests GetTransactionByBlockHashAndIndex.
func TestEthereumAPI_GetTransactionByBlockHashAndIndex(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, _, _, _ := createTestData(t)

	// Mock Backend functions.
	mockBackend.EXPECT().BlockByHash(gomock.Any(), gomock.Any()).Return(block, nil).Times(txs.Len())

	// Get transaction by block hash and index for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		ethTx := api.GetTransactionByBlockHashAndIndex(context.Background(), block.Hash(), hexutil.Uint(i))
		checkEthRPCTransactionFormat(t, block, ethTx, txs[i], hexutil.Uint64(i))
	}

	mockCtrl.Finish()
}

// TestEthereumAPI_GetTransactionByHash tests GetTransactionByHash.
func TestEthereumAPI_GetTransactionByHash(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, txHashMap, _, _ := createTestData(t)

	// Define queryFromPool for ReadTxAndLookupInfo function return tx from hash map.
	// MockDatabaseManager will initiate data with txHashMap, block and queryFromPool.
	// If queryFromPool is true, MockDatabaseManager will return nil to query transactions from transaction pool,
	// otherwise return a transaction from txHashMap.
	mockDBManager := &MockDatabaseManager{txHashMap: txHashMap, blockData: block, queryFromPool: false}

	// Mock Backend functions.
	mockBackend.EXPECT().ChainDB().Return(mockDBManager).Times(txs.Len())

	// Get transaction by hash for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		ethTx, err := api.GetTransactionByHash(context.Background(), txs[i].Hash())
		if err != nil {
			t.Fatal(err)
		}
		checkEthRPCTransactionFormat(t, block, ethTx, txs[i], hexutil.Uint64(i))
	}

	mockCtrl.Finish()
}

// TestEthereumAPI_GetTransactionByHash tests GetTransactionByHash from transaction pool.
func TestEthereumAPI_GetTransactionByHashFromPool(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, txHashMap, _, _ := createTestData(t)

	// Define queryFromPool for ReadTxAndLookupInfo function return nil.
	// MockDatabaseManager will initiate data with txHashMap, block and queryFromPool.
	// If queryFromPool is true, MockDatabaseManager will return nil to query transactions from transaction pool,
	// otherwise return a transaction from txHashMap.
	mockDBManager := &MockDatabaseManager{txHashMap: txHashMap, blockData: block, queryFromPool: true}

	// Mock Backend functions.
	mockBackend.EXPECT().ChainDB().Return(mockDBManager).Times(txs.Len())
	mockBackend.EXPECT().GetPoolTransaction(gomock.Any()).DoAndReturn(
		func(hash common.Hash) *types.Transaction {
			return txHashMap[hash]
		},
	).Times(txs.Len())

	//  Get transaction by hash from the transaction pool for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		ethTx, err := api.GetTransactionByHash(context.Background(), txs[i].Hash())
		if err != nil {
			t.Fatal(err)
		}
		checkEthRPCTransactionFormat(t, nil, ethTx, txs[i], 0)
	}

	mockCtrl.Finish()
}

// TestEthereumAPI_PendingTransactionstests PendingTransactions.
func TestEthereumAPI_PendingTransactions(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	_, txs, txHashMap, _, _ := createTestData(t)

	mockAccountManager := mock_accounts.NewMockAccountManager(mockCtrl)
	mockBackend.EXPECT().AccountManager().Return(mockAccountManager)

	mockBackend.EXPECT().GetPoolTransactions().Return(txs, nil)

	wallets := make([]accounts.Wallet, 1)
	wallets[0] = NewMockWallet(txs)
	mockAccountManager.EXPECT().Wallets().Return(wallets)

	pendingTxs, err := api.PendingTransactions()
	if err != nil {
		t.Fatal(err)
	}

	for _, pt := range pendingTxs {
		checkEthRPCTransactionFormat(t, nil, pt, txHashMap[pt.Hash], 0)
	}

	mockCtrl.Finish()
}

// TestEthereumAPI_GetTransactionReceipt tests GetTransactionReceipt.
func TestEthereumAPI_GetTransactionReceipt(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, txHashMap, receiptMap, receipts := createTestData(t)

	// Mock Backend functions.
	mockBackend.EXPECT().GetTxLookupInfoAndReceipt(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
			txLookupInfo := txHashMap[hash]
			idx := txLookupInfo.Nonce() // Assume idx of the transaction is nonce
			return txLookupInfo, block.Hash(), block.NumberU64(), idx, receiptMap[hash]
		},
	).Times(txs.Len())
	mockBackend.EXPECT().GetBlockReceipts(gomock.Any(), gomock.Any()).Return(receipts).Times(txs.Len())

	// Get receipt for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		receipt, err := api.GetTransactionReceipt(context.Background(), txs[i].Hash())
		if err != nil {
			t.Fatal(err)
		}
		txIdx := uint64(i)
		checkEthTransactionReceiptFormat(t, block, receipts, receipt, RpcOutputReceipt(txs[i], block.Hash(), block.NumberU64(), txIdx, receiptMap[txs[i].Hash()]), txIdx)
	}

	mockCtrl.Finish()
}

func testInitForEthApi(t *testing.T) (*gomock.Controller, *mock_api.MockBackend, EthereumAPI) {
	mockCtrl := gomock.NewController(t)
	mockBackend := mock_api.NewMockBackend(mockCtrl)

	blockchain.InitDeriveSha(types.ImplDeriveShaOriginal)

	api := EthereumAPI{
		publicTransactionPoolAPI: NewPublicTransactionPoolAPI(mockBackend, new(AddrLocker)),
	}
	return mockCtrl, mockBackend, api
}

func checkEthRPCTransactionFormat(t *testing.T, block *types.Block, ethTx *EthRPCTransaction, tx *types.Transaction, expectedIndex hexutil.Uint64) {
	// All Klaytn transaction types must be returned as TxTypeLegacyTransaction types.
	assert.Equal(t, types.TxType(ethTx.Type), types.TxTypeLegacyTransaction)

	// Check the data of common fields of the transaction.
	from := getFrom(tx)
	assert.Equal(t, from, ethTx.From)
	assert.Equal(t, hexutil.Uint64(tx.Gas()), ethTx.Gas)
	assert.Equal(t, tx.GasPrice(), ethTx.GasPrice.ToInt())
	assert.Equal(t, tx.Hash(), ethTx.Hash)
	assert.Equal(t, tx.GetTxInternalData().RawSignatureValues()[0].V, ethTx.V.ToInt())
	assert.Equal(t, tx.GetTxInternalData().RawSignatureValues()[0].R, ethTx.R.ToInt())
	assert.Equal(t, tx.GetTxInternalData().RawSignatureValues()[0].S, ethTx.S.ToInt())
	assert.Equal(t, hexutil.Uint64(tx.Nonce()), ethTx.Nonce)

	// Check the optional field of Klaytn transactions.
	assert.Equal(t, 0, bytes.Compare(ethTx.Input, tx.Data()))

	to := tx.To()
	switch tx.Type() {
	case types.TxTypeAccountUpdate, types.TxTypeFeeDelegatedAccountUpdate, types.TxTypeFeeDelegatedAccountUpdateWithRatio,
		types.TxTypeCancel, types.TxTypeFeeDelegatedCancel, types.TxTypeFeeDelegatedCancelWithRatio,
		types.TxTypeChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		assert.Equal(t, &from, ethTx.To)
	default:
		assert.Equal(t, to, ethTx.To)
	}
	value := tx.Value()
	assert.Equal(t, value, ethTx.Value.ToInt())

	// If it is not a pending transaction and has already been processed and added into a block,
	// the following fields should be returned.
	if block != nil {
		assert.Equal(t, block.Hash().String(), ethTx.BlockHash.String())
		assert.Equal(t, block.NumberU64(), ethTx.BlockNumber.ToInt().Uint64())
		assert.Equal(t, expectedIndex, *ethTx.TransactionIndex)
	}

	// Fields additionally used for Ethereum transaction types are not used
	// when returning Klaytn transactions.
	assert.Equal(t, true, reflect.ValueOf(ethTx.Accesses).IsNil())
	assert.Equal(t, true, reflect.ValueOf(ethTx.ChainID).IsNil())
	assert.Equal(t, true, reflect.ValueOf(ethTx.GasFeeCap).IsNil())
	assert.Equal(t, true, reflect.ValueOf(ethTx.GasTipCap).IsNil())
}

func checkEthTransactionReceiptFormat(t *testing.T, block *types.Block, receipts []*types.Receipt, ethReceipt map[string]interface{}, kReceipt map[string]interface{}, idx uint64) {
	tx := block.Transactions()[idx]

	// Check the common receipt fields.
	blockHash, ok := ethReceipt["blockHash"]
	if !ok {
		t.Fatal("blockHash is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, blockHash, kReceipt["blockHash"])

	blockNumber, ok := ethReceipt["blockNumber"]
	if !ok {
		t.Fatal("blockNumber is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, blockNumber.(hexutil.Uint64), hexutil.Uint64(kReceipt["blockNumber"].(*hexutil.Big).ToInt().Uint64()))

	transactionHash, ok := ethReceipt["transactionHash"]
	if !ok {
		t.Fatal("transactionHash is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, transactionHash, kReceipt["transactionHash"])

	transactionIndex, ok := ethReceipt["transactionIndex"]
	if !ok {
		t.Fatal("transactionIndex is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, transactionIndex, hexutil.Uint64(kReceipt["transactionIndex"].(hexutil.Uint)))

	from, ok := ethReceipt["from"]
	if !ok {
		t.Fatal("from is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, from, kReceipt["from"])

	// Klaytn transactions that do not use the 'To' field
	// fill in 'To' with from during converting format.
	toInTx := tx.To()
	fromAddress := getFrom(tx)
	to, ok := ethReceipt["to"]
	if !ok {
		t.Fatal("to is not defined in Ethereum transaction receipt format.")
	}
	switch tx.Type() {
	case types.TxTypeAccountUpdate, types.TxTypeFeeDelegatedAccountUpdate, types.TxTypeFeeDelegatedAccountUpdateWithRatio,
		types.TxTypeCancel, types.TxTypeFeeDelegatedCancel, types.TxTypeFeeDelegatedCancelWithRatio,
		types.TxTypeChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		assert.Equal(t, &fromAddress, to)
	default:
		assert.Equal(t, toInTx, to)
	}

	gasUsed, ok := ethReceipt["gasUsed"]
	if !ok {
		t.Fatal("gasUsed is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, gasUsed, kReceipt["gasUsed"])

	// Compare with the calculated cumulative gas used value
	// to check whether the cumulativeGasUsed value is calculated properly.
	cumulativeGasUsed, ok := ethReceipt["cumulativeGasUsed"]
	if !ok {
		t.Fatal("cumulativeGasUsed is not defined in Ethereum transaction receipt format.")
	}
	calculatedCumulativeGas := uint64(0)
	for i := 0; i <= int(idx); i++ {
		calculatedCumulativeGas += receipts[i].GasUsed
	}
	assert.Equal(t, cumulativeGasUsed, hexutil.Uint64(calculatedCumulativeGas))

	contractAddress, ok := ethReceipt["contractAddress"]
	if !ok {
		t.Fatal("contractAddress is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, contractAddress, kReceipt["contractAddress"])

	logs, ok := ethReceipt["logs"]
	if !ok {
		t.Fatal("logs is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, logs, kReceipt["logs"])

	logsBloom, ok := ethReceipt["logsBloom"]
	if !ok {
		t.Fatal("logsBloom is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, logsBloom, kReceipt["logsBloom"])

	typeInt, ok := ethReceipt["type"]
	if !ok {
		t.Fatal("type is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, types.TxType(typeInt.(hexutil.Uint)), types.TxTypeLegacyTransaction)

	effectiveGasPrice, ok := ethReceipt["effectiveGasPrice"]
	if !ok {
		t.Fatal("effectiveGasPrice is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, effectiveGasPrice, kReceipt["gasPrice"].(*hexutil.Big).ToInt())

	status, ok := ethReceipt["status"]
	if !ok {
		t.Fatal("status is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, status, kReceipt["status"])

	// Check the receipt fields that should be removed.
	var shouldNotExisted []string
	shouldNotExisted = append(shouldNotExisted, "gas", "gasPrice", "senderTxHash", "signatures", "txError", "typeInt", "feePayer", "feePayerSignatures", "feeRatio", "input", "value", "codeFormat", "humanReadable", "key", "inputJSON")
	for i := 0; i < len(shouldNotExisted); i++ {
		k := shouldNotExisted[i]
		_, ok = ethReceipt[k]
		if ok {
			t.Fatal(k, " should not be defined in the Ethereum transaction receipt format.")
		}
	}
}

func createTestData(t *testing.T) (*types.Block, types.Transactions, map[common.Hash]*types.Transaction, map[common.Hash]*types.Receipt, []*types.Receipt) {
	var txs types.Transactions

	var gasPrice = big.NewInt(25 * params.Ston)
	var deployData = "0x60806040526000805534801561001457600080fd5b506101ea806100246000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306661abd1461007257806342cbb15c1461009d578063767800de146100c8578063b22636271461011f578063d14e62b814610150575b600080fd5b34801561007e57600080fd5b5061008761017d565b6040518082815260200191505060405180910390f35b3480156100a957600080fd5b506100b2610183565b6040518082815260200191505060405180910390f35b3480156100d457600080fd5b506100dd61018b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561012b57600080fd5b5061014e60048036038101908080356000191690602001909291905050506101b1565b005b34801561015c57600080fd5b5061017b600480360381019080803590602001909291905050506101b4565b005b60005481565b600043905090565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b50565b80600081905550505600a165627a7a7230582053c65686a3571c517e2cf4f741d842e5ee6aa665c96ce70f46f9a594794f11eb0029"
	var executeData = "0xa9059cbb0000000000000000000000008a4c9c443bb0645df646a2d5bb55def0ed1e885a0000000000000000000000000000000000000000000000000000000000003039"
	var anchorData []byte

	var txHashMap = make(map[common.Hash]*types.Transaction)
	var receiptMap = make(map[common.Hash]*types.Receipt)
	var receipts []*types.Receipt

	// Create test data for chainDataAnchoring tx
	{
		dummyBlock := types.NewBlock(&types.Header{}, nil, nil)
		scData, err := types.NewAnchoringDataType0(dummyBlock, 0, uint64(dummyBlock.Transactions().Len()))
		if err != nil {
			t.Fatal(err)
		}
		anchorData, _ = rlp.EncodeToBytes(scData)
	}

	// Make test transactions data
	{
		// TxTypeLegacyTransaction
		values := map[types.TxValueKeyType]interface{}{
			// Simply set the nonce to txs.Len() to have a different nonce for each transaction type.
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyTo:       common.StringToAddress("0xe0680cfce04f80a386f1764a55c833b108770490"),
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(30000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     []byte("0xe3197e8f000000000000000000000000e0bef99b4a22286e27630b485d06c5a147cee931000000000000000000000000158beff8c8cdebd64654add5f6a1d9937e73536c0000000000000000000000000000000000000000000029bb5e7fb6beae32cf8000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000180000000000000000000000000000000000000000000000000001b60fb4614a22e000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000158beff8c8cdebd64654add5f6a1d9937e73536c00000000000000000000000074ba03198fed2b15a51af242b9c63faf3c8f4d3400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeLegacyTransaction, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])

	}
	{
		// TxTypeValueTransfer
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x520af902892196a3449b06ead301daeaf67e77e8"),
			types.TxValueKeyTo:       common.StringToAddress("0xa06fa690d92788cac4953da5f2dfbc4a2b3871db"),
			types.TxValueKeyAmount:   big.NewInt(5),
			types.TxValueKeyGasLimit: uint64(10000000),
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeValueTransferMemo
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0xc05e11f9075d453b4fc87023f815fa6a63f7aa0c"),
			types.TxValueKeyTo:       common.StringToAddress("0xb5a2d79e9228f3d278cb36b5b15930f24fe8bae8"),
			types.TxValueKeyAmount:   big.NewInt(3),
			types.TxValueKeyGasLimit: uint64(20000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     []byte(string("hello")),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeAccountUpdate
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      uint64(txs.Len()),
			types.TxValueKeyFrom:       common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:   uint64(20000000),
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: accountkey.NewAccountKeyLegacy(),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeSmartContractDeploy
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         uint64(txs.Len()),
			types.TxValueKeyFrom:          common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        big.NewInt(0),
			types.TxValueKeyGasLimit:      uint64(100000000),
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyData:          common.Hex2Bytes(deployData),
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		var r = createReceipt(t, tx, tx.Gas())
		fromAddress, err := tx.From()
		if err != nil {
			t.Fatal(err)
		}
		tx.FillContractAddress(fromAddress, r)
		receiptMap[tx.Hash()] = r
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeSmartContractExecution
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:       common.StringToAddress("0x00ca1eee49a4d2b04e6562222eab95e9ed29c4bf"),
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(50000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     common.Hex2Bytes(executeData),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeCancel
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit: uint64(50000000),
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeCancel, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeChainDataAnchoring
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:        uint64(txs.Len()),
			types.TxValueKeyFrom:         common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:     uint64(50000000),
			types.TxValueKeyGasPrice:     gasPrice,
			types.TxValueKeyAnchoredData: anchorData,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeChainDataAnchoring, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedValueTransfer
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x520af902892196a3449b06ead301daeaf67e77e8"),
			types.TxValueKeyTo:       common.StringToAddress("0xa06fa690d92788cac4953da5f2dfbc4a2b3871db"),
			types.TxValueKeyAmount:   big.NewInt(5),
			types.TxValueKeyGasLimit: uint64(10000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedValueTransferMemo
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0xc05e11f9075d453b4fc87023f815fa6a63f7aa0c"),
			types.TxValueKeyTo:       common.StringToAddress("0xb5a2d79e9228f3d278cb36b5b15930f24fe8bae8"),
			types.TxValueKeyAmount:   big.NewInt(3),
			types.TxValueKeyGasLimit: uint64(20000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     []byte(string("hello")),
			types.TxValueKeyFeePayer: common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemo, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx

		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedAccountUpdate
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      uint64(txs.Len()),
			types.TxValueKeyFrom:       common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:   uint64(20000000),
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: accountkey.NewAccountKeyLegacy(),
			types.TxValueKeyFeePayer:   common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedSmartContractDeploy
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         uint64(txs.Len()),
			types.TxValueKeyFrom:          common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        big.NewInt(0),
			types.TxValueKeyGasLimit:      uint64(100000000),
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyData:          common.Hex2Bytes(deployData),
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
			types.TxValueKeyFeePayer:      common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		var r = createReceipt(t, tx, tx.Gas())
		fromAddress, err := tx.From()
		if err != nil {
			t.Fatal(err)
		}
		tx.FillContractAddress(fromAddress, r)
		receiptMap[tx.Hash()] = r
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedSmartContractExecution
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:       common.StringToAddress("0x00ca1eee49a4d2b04e6562222eab95e9ed29c4bf"),
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(50000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     common.Hex2Bytes(executeData),
			types.TxValueKeyFeePayer: common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecution, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedCancel
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit: uint64(50000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancel, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedChainDataAnchoring
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:        uint64(txs.Len()),
			types.TxValueKeyFrom:         common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:     uint64(50000000),
			types.TxValueKeyGasPrice:     gasPrice,
			types.TxValueKeyAnchoredData: anchorData,
			types.TxValueKeyFeePayer:     common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedChainDataAnchoring, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx

		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedValueTransferWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x520af902892196a3449b06ead301daeaf67e77e8"),
			types.TxValueKeyTo:                 common.StringToAddress("0xa06fa690d92788cac4953da5f2dfbc4a2b3871db"),
			types.TxValueKeyAmount:             big.NewInt(5),
			types.TxValueKeyGasLimit:           uint64(10000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedValueTransferMemoWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0xc05e11f9075d453b4fc87023f815fa6a63f7aa0c"),
			types.TxValueKeyTo:                 common.StringToAddress("0xb5a2d79e9228f3d278cb36b5b15930f24fe8bae8"),
			types.TxValueKeyAmount:             big.NewInt(3),
			types.TxValueKeyGasLimit:           uint64(20000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               []byte(string("hello")),
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemoWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedAccountUpdateWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:           uint64(20000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyAccountKey:         accountkey.NewAccountKeyLegacy(),
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedSmartContractDeployWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:                 (*common.Address)(nil),
			types.TxValueKeyAmount:             big.NewInt(0),
			types.TxValueKeyGasLimit:           uint64(100000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               common.Hex2Bytes(deployData),
			types.TxValueKeyHumanReadable:      false,
			types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		var r = createReceipt(t, tx, tx.Gas())
		fromAddress, err := tx.From()
		if err != nil {
			t.Fatal(err)
		}
		tx.FillContractAddress(fromAddress, r)
		receiptMap[tx.Hash()] = r
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedSmartContractExecutionWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:                 common.StringToAddress("0x00ca1eee49a4d2b04e6562222eab95e9ed29c4bf"),
			types.TxValueKeyAmount:             big.NewInt(0),
			types.TxValueKeyGasLimit:           uint64(50000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               common.Hex2Bytes(executeData),
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecutionWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedCancelWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:           uint64(50000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedChainDataAnchoringWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:           uint64(50000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyAnchoredData:       anchorData,
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedChainDataAnchoringWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}

	// Create a block which includes all transaction data.
	block := types.NewBlock(&types.Header{Number: big.NewInt(1)}, txs, nil)

	return block, txs, txHashMap, receiptMap, receipts
}

func createReceipt(t *testing.T, tx *types.Transaction, gasUsed uint64) *types.Receipt {
	rct := types.NewReceipt(uint(0), tx.Hash(), gasUsed)
	rct.Logs = []*types.Log{}
	rct.Bloom = types.Bloom{}
	return rct
}

// MockDatabaseManager is a mock of database.DBManager interface for overriding the ReadTxAndLookupInfo function.
type MockDatabaseManager struct {
	database.DBManager

	txHashMap     map[common.Hash]*types.Transaction
	blockData     *types.Block
	queryFromPool bool
}

// GetTxLookupInfoAndReceipt retrieves a tx and lookup info and receipt for a given transaction hash.
func (dbm *MockDatabaseManager) ReadTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	// If queryFromPool, return nil to query from pool after this function
	if dbm.queryFromPool {
		return nil, common.Hash{}, 0, 0
	}

	txFromHashMap := dbm.txHashMap[hash]
	if txFromHashMap == nil {
		return nil, common.Hash{}, 0, 0
	}
	return txFromHashMap, dbm.blockData.Hash(), dbm.blockData.NumberU64(), txFromHashMap.Nonce()
}

// MockWallet is a mock of accounts.Wallet interface for overriding the Accounts function.
type MockWallet struct {
	accounts.Wallet

	accounts []accounts.Account
}

// NewMockWallet prepares accounts based on tx from.
func NewMockWallet(txs types.Transactions) *MockWallet {
	mw := &MockWallet{}

	for _, t := range txs {
		mw.accounts = append(mw.accounts, accounts.Account{Address: getFrom(t)})
	}
	return mw
}

// Accounts implements accounts.Wallet, returning an account list.
func (mw *MockWallet) Accounts() []accounts.Account {
	return mw.accounts
}
