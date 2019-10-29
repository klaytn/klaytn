package api

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/accounts/mocks"
	"github.com/klaytn/klaytn/api/mocks"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"testing"
)

// test tx types and internal data to be supported by APIs in PublicTransactionPoolAPI.
var internalDataTypes = map[types.TxType]interface{}{
	types.TxTypeLegacyTransaction:                           types.TxInternalDataLegacy{},
	types.TxTypeValueTransfer:                               types.TxInternalDataValueTransfer{},
	types.TxTypeFeeDelegatedValueTransfer:                   types.TxInternalDataFeeDelegatedValueTransfer{},
	types.TxTypeFeeDelegatedValueTransferWithRatio:          types.TxInternalDataFeeDelegatedValueTransferWithRatio{},
	types.TxTypeValueTransferMemo:                           types.TxInternalDataValueTransferMemo{},
	types.TxTypeFeeDelegatedValueTransferMemo:               types.TxInternalDataFeeDelegatedValueTransferMemo{},
	types.TxTypeFeeDelegatedValueTransferMemoWithRatio:      types.TxInternalDataFeeDelegatedValueTransferMemoWithRatio{},
	types.TxTypeAccountUpdate:                               types.TxInternalDataAccountUpdate{},
	types.TxTypeFeeDelegatedAccountUpdate:                   types.TxInternalDataFeeDelegatedAccountUpdate{},
	types.TxTypeFeeDelegatedAccountUpdateWithRatio:          types.TxInternalDataFeeDelegatedAccountUpdateWithRatio{},
	types.TxTypeSmartContractDeploy:                         types.TxInternalDataSmartContractDeploy{},
	types.TxTypeFeeDelegatedSmartContractDeploy:             types.TxInternalDataFeeDelegatedSmartContractDeploy{},
	types.TxTypeFeeDelegatedSmartContractDeployWithRatio:    types.TxInternalDataFeeDelegatedSmartContractDeployWithRatio{},
	types.TxTypeSmartContractExecution:                      types.TxInternalDataSmartContractExecution{},
	types.TxTypeFeeDelegatedSmartContractExecution:          types.TxInternalDataFeeDelegatedSmartContractExecution{},
	types.TxTypeFeeDelegatedSmartContractExecutionWithRatio: types.TxInternalDataFeeDelegatedSmartContractExecutionWithRatio{},
	types.TxTypeCancel:                                      types.TxInternalDataCancel{},
	types.TxTypeFeeDelegatedCancel:                          types.TxInternalDataFeeDelegatedCancel{},
	types.TxTypeFeeDelegatedCancelWithRatio:                 types.TxInternalDataFeeDelegatedCancelWithRatio{},
}

// test values of tx field.
var (
	testNonce         = hexutil.Uint64(0)
	testGasPrice      = (*hexutil.Big)(big.NewInt(25 * params.Ston))
	testValue         = (*hexutil.Big)(big.NewInt(1))
	testTo            = common.StringToAddress("1234")
	testFeePayer      = common.StringToAddress("5678")
	testFeeRatio      = types.FeeRatio(30)
	testData          = hexutil.Bytes{0x11, 0x99}
	testCodeFormat    = params.CodeFormatEVM
	testHumanReadable = false
	testAccountKey    = hexutil.Bytes{0x01, 0xc0}
	testFrom          = common.HexToAddress("0xa7Eb6992c5FD55F43305B24Ee67150Bf4910d329")

	senderPrvKey, _ = crypto.HexToECDSA("95a21e86efa290d6665a9dbce06ae56319335540d13540fb1b01e28a5b2c8460")
)

// TestTxTypeSupport tests tx type support of APIs in PublicTransactionPoolAPI.
func TestTxTypeSupport(t *testing.T) {
	var ctx context.Context
	chainConf := params.ChainConfig{ChainID: big.NewInt(1)}

	// generate a keystore and an active account
	dir, err := ioutil.TempDir("", "klay-keystore-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	ks := keystore.NewKeyStore(dir, 2, 1)
	password := ""
	acc, err := ks.ImportECDSA(senderPrvKey, password)
	if err != nil {
		t.Fatal(err)
	}
	if err := ks.Unlock(acc, password); err != nil {
		t.Fatal(err)
	}

	// mock Backend and AccountManager for easy test
	mockCtrl := gomock.NewController(t)
	mockBackend := mock_api.NewMockBackend(mockCtrl)
	mockAccountManager := mock_accounts.NewMockAccountManager(mockCtrl)

	mockBackend.EXPECT().AccountManager().Return(mockAccountManager).AnyTimes()
	mockBackend.EXPECT().SuggestPrice(ctx).Return((*big.Int)(testGasPrice), nil).AnyTimes()
	mockBackend.EXPECT().GetPoolNonce(ctx, gomock.Any()).Return(uint64(testNonce)).AnyTimes()
	mockBackend.EXPECT().SendTx(ctx, gomock.Any()).Return(nil).AnyTimes()
	mockBackend.EXPECT().ChainConfig().Return(&chainConf).AnyTimes()
	mockAccountManager.EXPECT().Find(gomock.Any()).Return(ks.Wallets()[0], nil).AnyTimes()

	// APIs in PublicTransactionPoolAPI will be tested
	api := PublicTransactionPoolAPI{
		b:         mockBackend,
		nonceLock: new(AddrLocker),
	}

	// test for all possible tx types
	for txType, internalData := range internalDataTypes {
		// args contains values of tx fields
		args := SendTxArgs{
			TypeInt: &txType,
			From:    testFrom,
		}

		// set required fields of each typed tx
		internalType := reflect.TypeOf(internalData)
		for i := 0; i < internalType.NumField(); i++ {
			switch internalType.Field(i).Name {
			case "Amount":
				args.Value = testValue
			case "Recipient":
				args.To = &testTo
			case "FeePayer":
				args.FeePayer = &testFeePayer
			case "FeeRatio":
				args.FeeRatio = &testFeeRatio
			case "Payload":
				args.Data = &testData
			case "CodeFormat":
				args.CodeFormat = &testCodeFormat
			case "HumanReadable":
				args.HumanReadable = &testHumanReadable
			case "Key":
				args.AccountKey = &testAccountKey
			}
		}
		// tests for non-fee-delegation types
		if !txType.IsFeeDelegatedTransaction() {
			testTxTypeSupport_normalCase(t, api, ctx, args)
		}
		// TODO - more test cases will be added
	}
}

// testTxTypeSupport_normalCase test APIs with proper SendTxArgs values.
func testTxTypeSupport_normalCase(t *testing.T, api PublicTransactionPoolAPI, ctx context.Context, args SendTxArgs) {
	// test tx type support of SignTransaction
	{
		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, nil, err)
	}

	// test tx type support of SendTransaction
	{
		_, err := api.SendTransaction(ctx, args)
		assert.Equal(t, nil, err)
	}
}
