package api

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/keystore"
	mock_accounts "github.com/klaytn/klaytn/accounts/mocks"
	mock_api "github.com/klaytn/klaytn/api/mocks"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
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
	types.TxTypeChainDataAnchoring:                          types.TxInternalDataChainDataAnchoring{},
	types.TxTypeFeeDelegatedChainDataAnchoring:              types.TxInternalDataFeeDelegatedChainDataAnchoring{},
	types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:     types.TxInternalDataFeeDelegatedChainDataAnchoringWithRatio{},
}

// test values of tx field.
var (
	testNonce         = hexutil.Uint64(0)
	testGas           = hexutil.Uint64(900000)
	testGasPrice      = (*hexutil.Big)(big.NewInt(25 * params.Ston))
	testValue         = (*hexutil.Big)(big.NewInt(1))
	testTo            = common.StringToAddress("1234")
	testFeePayer      = common.HexToAddress("0x819104a190255e0cedbdd9d5f59a557633d79db1")
	testFeeRatio      = types.FeeRatio(30)
	testData          = hexutil.Bytes{0x11, 0x99}
	testCodeFormat    = params.CodeFormatEVM
	testHumanReadable = false
	testAccountKey    = hexutil.Bytes{0x01, 0xc0}
	testFrom          = common.HexToAddress("0xa7Eb6992c5FD55F43305B24Ee67150Bf4910d329")
	testSig           = types.TxSignatures{&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)}}.ToJSON()
	senderPrvKey, _   = crypto.HexToECDSA("95a21e86efa290d6665a9dbce06ae56319335540d13540fb1b01e28a5b2c8460")
	feePayerPrvKey, _ = crypto.HexToECDSA("aebb680a5e596c1d1a01bac78a3985b62c685c5e995d780c176138cb2679ba3e")
)

// TestTxTypeSupport tests tx type support of APIs in PublicTransactionPoolAPI.
func TestTxTypeSupport(t *testing.T) {
	var ctx context.Context
	chainConf := params.ChainConfig{ChainID: big.NewInt(1)}

	// generate a keystore and active accounts
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
	accFeePayer, err := ks.ImportECDSA(feePayerPrvKey, password)
	if err != nil {
		t.Fatal(err)
	}
	if err := ks.Unlock(accFeePayer, password); err != nil {
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
	mockAccountManager.EXPECT().Find(accounts.Account{Address: acc.Address}).Return(ks.Wallets()[0], nil).AnyTimes()
	mockAccountManager.EXPECT().Find(accounts.Account{Address: accFeePayer.Address}).Return(ks.Wallets()[1], nil).AnyTimes()

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
			case "AccountNonce":
				args.AccountNonce = &testNonce
			case "Amount":
				args.Amount = testValue
			case "Recipient":
				args.Recipient = &testTo
			case "FeePayer":
				args.FeePayer = &testFeePayer
			case "FeeRatio":
				args.FeeRatio = &testFeeRatio
			case "GasLimit":
				args.GasLimit = &testGas
			case "Price":
				args.Price = testGasPrice
			case "Payload":
				args.Payload = &testData
			case "CodeFormat":
				args.CodeFormat = &testCodeFormat
			case "HumanReadable":
				args.HumanReadable = &testHumanReadable
			case "Key":
				args.Key = &testAccountKey
			}
		}
		if txType.IsFeeDelegatedTransaction() {
			args.TxSignatures = testSig
		}

		testTxTypeSupport_normalCase(t, api, ctx, args)
		testTxTypeSupport_setDefault(t, api, ctx, args)
		testTxTypeSupport_noFieldValues(t, api, ctx, args)
		testTxTypeSupport_unnecessaryFieldValues(t, api, ctx, args)
	}
}

// testTxTypeSupport_normalCase tests APIs with proper SendTxArgs values.
func testTxTypeSupport_normalCase(t *testing.T, api PublicTransactionPoolAPI, ctx context.Context, args SendTxArgs) {
	var err error

	// test APIs for non-fee-delegation txs
	if !args.TypeInt.IsFeeDelegatedTransaction() {
		_, err = api.SendTransaction(ctx, args)
		assert.Equal(t, nil, err)

		// test APIs for fee delegation txs
	} else {
		_, err := api.SignTransactionAsFeePayer(ctx, args)
		assert.Equal(t, nil, err)

		_, err = api.SendTransactionAsFeePayer(ctx, args)
		assert.Equal(t, nil, err)
	}

	// test for all txs
	_, err = api.SignTransaction(ctx, args)
	assert.Equal(t, nil, err)
}

// testTxTypeSupport_setDefault tests the setDefault function which auto-assign some values of tx.
func testTxTypeSupport_setDefault(t *testing.T, api PublicTransactionPoolAPI, ctx context.Context, args SendTxArgs) {
	args.AccountNonce = nil
	args.GasLimit = nil
	args.Price = nil

	_, err := api.SignTransaction(ctx, args)
	assert.Equal(t, nil, err)
}

// testTxTypeSupport_noFieldValues tests error handling for not assigned field values.
func testTxTypeSupport_noFieldValues(t *testing.T, api PublicTransactionPoolAPI, ctx context.Context, oriArgs SendTxArgs) {
	// fields of legacy tx will not be checked in the checkArgs function
	if *oriArgs.TypeInt == types.TxTypeLegacyTransaction {
		return
	}

	args := oriArgs
	if args.Recipient != nil && !(*args.TypeInt).IsContractDeploy() {
		args.Recipient = nil

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"to\" is required for "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.Payload != nil {
		args.Payload = nil

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"input\" is required for "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.Amount != nil {
		args.Amount = nil

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"value\" is required for "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.CodeFormat != nil {
		args.CodeFormat = nil

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"codeFormat\" is required for "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.HumanReadable != nil {
		args.HumanReadable = nil

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"humanReadable\" is required for "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.Key != nil {
		args.Key = nil

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"key\" is required for "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.FeePayer != nil {
		args.FeePayer = nil

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"feePayer\" is required for "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.FeeRatio != nil {
		args.FeeRatio = nil

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"feeRatio\" is required for "+(*args.TypeInt).String(), err.Error())
	}
}

// testTxTypeSupport_unnecessaryFieldValues tests error handling for not assigned field values.
func testTxTypeSupport_unnecessaryFieldValues(t *testing.T, api PublicTransactionPoolAPI, ctx context.Context, oriArgs SendTxArgs) {
	// fields of legacy tx will not be checked in the checkArgs function
	if *oriArgs.TypeInt == types.TxTypeLegacyTransaction {
		return
	}

	args := oriArgs
	if args.Recipient == nil {
		args.Recipient = &testTo

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"to\" is not a field of "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.Payload == nil {
		args.Payload = &testData

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"input\" is not a field of "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.Amount == nil {
		args.Amount = testValue

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"value\" is not a field of "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.CodeFormat == nil {
		args.CodeFormat = &testCodeFormat

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"codeFormat\" is not a field of "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.HumanReadable == nil {
		args.HumanReadable = &testHumanReadable

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"humanReadable\" is not a field of "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.Key == nil {
		args.Key = &testAccountKey

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"key\" is not a field of "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.FeePayer == nil {
		args.FeePayer = &testFeePayer

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"feePayer\" is not a field of "+(*args.TypeInt).String(), err.Error())
	}

	args = oriArgs
	if args.FeeRatio == nil {
		args.FeeRatio = &testFeeRatio

		_, err := api.SignTransaction(ctx, args)
		assert.Equal(t, "json:\"feeRatio\" is not a field of "+(*args.TypeInt).String(), err.Error())
	}
}
