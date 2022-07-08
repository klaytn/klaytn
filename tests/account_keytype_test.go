// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package tests

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

// createDefaultAccount creates a default account with a specific account key type.
func createDefaultAccount(accountKeyType accountkey.AccountKeyType) (*TestAccountType, error) {
	var err error

	// prepare  keys
	keys := genTestKeys(3)
	weights := []uint{1, 1, 1}
	weightedKeys := make(accountkey.WeightedPublicKeys, 3)
	threshold := uint(2)

	for i := range keys {
		weightedKeys[i] = accountkey.NewWeightedPublicKey(weights[i], (*accountkey.PublicKeySerializable)(&keys[i].PublicKey))
	}

	// a role-based key
	roleAccKey := accountkey.AccountKeyRoleBased{
		accountkey.NewAccountKeyPublicWithValue(&keys[accountkey.RoleTransaction].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&keys[accountkey.RoleAccountUpdate].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&keys[accountkey.RoleFeePayer].PublicKey),
	}

	// default account setting
	account := &TestAccountType{
		Addr:   crypto.PubkeyToAddress(keys[0].PublicKey), // default
		Keys:   []*ecdsa.PrivateKey{keys[0]},              // default
		Nonce:  uint64(0),                                 // default
		AccKey: nil,
	}

	// set an account key and a private key
	switch accountKeyType {
	case accountkey.AccountKeyTypeNil:
		account.AccKey, err = accountkey.NewAccountKey(accountKeyType)
	case accountkey.AccountKeyTypeLegacy:
		account.AccKey, err = accountkey.NewAccountKey(accountKeyType)
	case accountkey.AccountKeyTypePublic:
		account.AccKey = accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey)
	case accountkey.AccountKeyTypeFail:
		account.AccKey, err = accountkey.NewAccountKey(accountKeyType)
	case accountkey.AccountKeyTypeWeightedMultiSig:
		account.Keys = keys
		account.AccKey = accountkey.NewAccountKeyWeightedMultiSigWithValues(threshold, weightedKeys)
	case accountkey.AccountKeyTypeRoleBased:
		account.Keys = keys
		account.AccKey = accountkey.NewAccountKeyRoleBasedWithValues(roleAccKey)
	default:
		return nil, kerrors.ErrDifferentAccountKeyType
	}
	if err != nil {
		return nil, err
	}

	return account, err
}

// generateDefaultTx returns a Tx with default values of txTypes.
// If txType is a kind of account update, it will return an account to update.
// Otherwise, it will return (tx, nil, nil).
// For contract execution Txs, TxValueKeyTo value is set to "contract" as a default.
// The address "contact" should exist before calling this function.
func generateDefaultTx(sender *TestAccountType, recipient *TestAccountType, txType types.TxType, contractAddr common.Address) (*types.Transaction, *TestAccountType, error) {
	gasPrice := new(big.Int).SetUint64(25 * params.Ston)

	// For Dynamic fee tx.
	gasTipCap := new(big.Int).SetUint64(25 * params.Ston)
	gasFeeCap := new(big.Int).SetUint64(25 * params.Ston)

	gasLimit := uint64(10000000)
	amount := new(big.Int).SetUint64(1)

	// generate a new account for account creation/update Txs or contract deploy Txs
	senderAccType := accountkey.AccountKeyTypeLegacy
	if sender.AccKey != nil {
		senderAccType = sender.AccKey.Type()
	}
	newAcc, err := createDefaultAccount(senderAccType)
	if err != nil {
		return nil, nil, err
	}

	// Smart contract data for TxTypeSmartContractDeploy, TxTypeSmartContractExecution Txs
	var code string
	var abiStr string

	if isCompilerAvailable() {
		filename := string("../contracts/reward/contract/KlaytnReward.sol")
		codes, abistrings := compileSolidity(filename)
		code = codes[0]
		abiStr = abistrings[0]
	} else {
		// Falling back to use compiled code.
		code = "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"
		abiStr = `[{"constant":true,"inputs":[],"name":"totalAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"receiver","type":"address"}],"name":"reward","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"safeWithdrawal","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`
	}

	abii, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		return nil, nil, err
	}

	dataABI, err := abii.Pack("reward", recipient.Addr)
	if err != nil {
		return nil, nil, err
	}

	// generate a legacy tx
	if txType == types.TxTypeLegacyTransaction {
		tx := types.NewTransaction(sender.Nonce, recipient.Addr, amount, gasLimit, gasPrice, []byte{})
		return tx, nil, nil
	}

	// Default valuesMap setting
	amountZero := new(big.Int).SetUint64(0)
	ratio := types.FeeRatio(30)
	dataMemo := []byte("hello")
	dataAnchor := []byte{0x11, 0x22}
	dataCode := common.FromHex(code)
	values := map[types.TxValueKeyType]interface{}{}

	switch txType {
	case types.TxTypeValueTransfer:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = recipient.Addr
		values[types.TxValueKeyAmount] = amount
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
	case types.TxTypeFeeDelegatedValueTransfer:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = recipient.Addr
		values[types.TxValueKeyAmount] = amount
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyFeePayer] = recipient.Addr
	case types.TxTypeFeeDelegatedValueTransferWithRatio:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = recipient.Addr
		values[types.TxValueKeyAmount] = amount
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyFeePayer] = recipient.Addr
		values[types.TxValueKeyFeeRatioOfFeePayer] = ratio
	case types.TxTypeValueTransferMemo:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = recipient.Addr
		values[types.TxValueKeyAmount] = amount
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyData] = dataMemo
	case types.TxTypeFeeDelegatedValueTransferMemo:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = recipient.Addr
		values[types.TxValueKeyAmount] = amount
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyData] = dataMemo
		values[types.TxValueKeyFeePayer] = recipient.Addr
	case types.TxTypeFeeDelegatedValueTransferMemoWithRatio:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = recipient.Addr
		values[types.TxValueKeyAmount] = amount
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyData] = dataMemo
		values[types.TxValueKeyFeePayer] = recipient.Addr
		values[types.TxValueKeyFeeRatioOfFeePayer] = ratio
	case types.TxTypeAccountUpdate:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyAccountKey] = newAcc.AccKey
	case types.TxTypeFeeDelegatedAccountUpdate:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyAccountKey] = newAcc.AccKey
		values[types.TxValueKeyFeePayer] = recipient.Addr
	case types.TxTypeFeeDelegatedAccountUpdateWithRatio:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyAccountKey] = newAcc.AccKey
		values[types.TxValueKeyFeePayer] = recipient.Addr
		values[types.TxValueKeyFeeRatioOfFeePayer] = ratio
	case types.TxTypeSmartContractDeploy:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = (*common.Address)(nil)
		values[types.TxValueKeyAmount] = amountZero
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = amountZero
		values[types.TxValueKeyData] = dataCode
		values[types.TxValueKeyHumanReadable] = false
		values[types.TxValueKeyCodeFormat] = params.CodeFormatEVM
	case types.TxTypeFeeDelegatedSmartContractDeploy:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = (*common.Address)(nil)
		values[types.TxValueKeyAmount] = amountZero
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = amountZero
		values[types.TxValueKeyData] = dataCode
		values[types.TxValueKeyHumanReadable] = false
		values[types.TxValueKeyFeePayer] = recipient.Addr
		values[types.TxValueKeyCodeFormat] = params.CodeFormatEVM
	case types.TxTypeFeeDelegatedSmartContractDeployWithRatio:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = (*common.Address)(nil)
		values[types.TxValueKeyAmount] = amountZero
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = amountZero
		values[types.TxValueKeyData] = dataCode
		values[types.TxValueKeyHumanReadable] = false
		values[types.TxValueKeyFeePayer] = recipient.Addr
		values[types.TxValueKeyFeeRatioOfFeePayer] = ratio
		values[types.TxValueKeyCodeFormat] = params.CodeFormatEVM
	case types.TxTypeSmartContractExecution:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = contractAddr
		values[types.TxValueKeyAmount] = amountZero
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = amountZero
		values[types.TxValueKeyData] = dataABI
	case types.TxTypeFeeDelegatedSmartContractExecution:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = contractAddr
		values[types.TxValueKeyAmount] = amountZero
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = amountZero
		values[types.TxValueKeyData] = dataABI
		values[types.TxValueKeyFeePayer] = recipient.Addr
	case types.TxTypeFeeDelegatedSmartContractExecutionWithRatio:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyTo] = contractAddr
		values[types.TxValueKeyAmount] = amountZero
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = amountZero
		values[types.TxValueKeyData] = dataABI
		values[types.TxValueKeyFeePayer] = recipient.Addr
		values[types.TxValueKeyFeeRatioOfFeePayer] = ratio
	case types.TxTypeCancel:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
	case types.TxTypeFeeDelegatedCancel:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyFeePayer] = recipient.Addr
	case types.TxTypeFeeDelegatedCancelWithRatio:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyFeePayer] = recipient.Addr
		values[types.TxValueKeyFeeRatioOfFeePayer] = ratio
	case types.TxTypeChainDataAnchoring:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyAnchoredData] = dataAnchor
	case types.TxTypeFeeDelegatedChainDataAnchoring:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyAnchoredData] = dataAnchor
		values[types.TxValueKeyFeePayer] = recipient.Addr
	case types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyFrom] = sender.Addr
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyAnchoredData] = dataAnchor
		values[types.TxValueKeyFeePayer] = recipient.Addr
		values[types.TxValueKeyFeeRatioOfFeePayer] = ratio
	case types.TxTypeEthereumAccessList:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyTo] = &recipient.Addr
		values[types.TxValueKeyAmount] = amount
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasPrice] = gasPrice
		values[types.TxValueKeyChainID] = big.NewInt(1)
		values[types.TxValueKeyData] = dataCode
		values[types.TxValueKeyAccessList] = types.AccessList{}
	case types.TxTypeEthereumDynamicFee:
		values[types.TxValueKeyNonce] = sender.Nonce
		values[types.TxValueKeyTo] = &recipient.Addr
		values[types.TxValueKeyAmount] = amount
		values[types.TxValueKeyGasLimit] = gasLimit
		values[types.TxValueKeyGasFeeCap] = gasFeeCap
		values[types.TxValueKeyGasTipCap] = gasTipCap
		values[types.TxValueKeyChainID] = big.NewInt(1)
		values[types.TxValueKeyData] = dataCode
		values[types.TxValueKeyAccessList] = types.AccessList{}
	}

	tx, err := types.NewTransactionWithMap(txType, values)
	if err != nil {
		return nil, nil, err
	}

	// the function returns an updated sender account for account update Txs
	if txType.IsAccountUpdate() {
		// For the account having a legacy key, its private key will not be updated since it is coupled with its address.
		if newAcc.AccKey.Type().IsLegacyAccountKey() {
			newAcc.Keys = sender.Keys
		}
		newAcc.Addr = sender.Addr
		newAcc.Nonce = sender.Nonce
		return tx, newAcc, err
	}

	return tx, nil, err
}

// expectedTestResultForDefaultTx returns expected validity of tx which generated from (accountKeyType, txType) pair.
func expectedTestResultForDefaultTx(accountKeyType accountkey.AccountKeyType, txType types.TxType) error {
	switch accountKeyType {
	// case accountkey.AccountKeyTypeNil:                     // not supported type
	case accountkey.AccountKeyTypeFail:
		if txType.IsAccountUpdate() {
			return kerrors.ErrAccountKeyFailNotUpdatable
		}
		return types.ErrInvalidSigSender
	}
	return nil
}

func signTxWithVariousKeyTypes(signer types.Signer, tx *types.Transaction, sender *TestAccountType) (*types.Transaction, error) {
	var err error
	txType := tx.Type()
	accKeyType := sender.AccKey.Type()

	if accKeyType == accountkey.AccountKeyTypeWeightedMultiSig {
		if txType.IsLegacyTransaction() {
			err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{sender.Keys[0]})
		} else {
			err = tx.SignWithKeys(signer, sender.Keys)
		}
	} else if accKeyType == accountkey.AccountKeyTypeRoleBased {
		if txType.IsAccountUpdate() {
			err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{sender.Keys[accountkey.RoleAccountUpdate]})
		} else {
			err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{sender.Keys[accountkey.RoleTransaction]})
		}
	} else {
		err = tx.SignWithKeys(signer, sender.Keys)
	}
	return tx, err
}

// TestDefaultTxsWithDefaultAccountKey tests most of transactions types with most of account key types.
// The test creates a default account for each account key type, and generates default Tx for each Tx type.
// AccountKeyTypeNil is excluded because it cannot be used for account creation.
func TestDefaultTxsWithDefaultAccountKey(t *testing.T) {
	gasPrice := new(big.Int).SetUint64(25 * params.Ston)
	gasLimit := uint64(100000000)

	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	bcdata.bc.Config().IstanbulCompatibleBlock = big.NewInt(0)
	bcdata.bc.Config().LondonCompatibleBlock = big.NewInt(0)
	bcdata.bc.Config().EthTxTypeCompatibleBlock = big.NewInt(0)
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// smart contact account
	contractAddr := common.Address{}

	// smart contract code
	var code string

	if isCompilerAvailable() {
		filename := string("../contracts/reward/contract/KlaytnReward.sol")
		codes, _ := compileSolidity(filename)
		code = codes[0]
	} else {
		// Falling back to use compiled code.
		code = "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	// create a smart contract account for contract execution test
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(0)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.Nonce,
			types.TxValueKeyFrom:          reservoir.Addr,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      uint64(50 * uint64(params.Ston)),
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyData:          common.FromHex(code),
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		err = bcdata.GenABlockWithTransactions(accountMap, txs, prof)
		assert.Equal(t, nil, err)

		contractAddr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1
	}
	// select account key types to be tested
	accountKeyTypes := []accountkey.AccountKeyType{
		// accountkey.AccountKeyTypeNil, // not supported type
		accountkey.AccountKeyTypeLegacy,
		accountkey.AccountKeyTypePublic,
		accountkey.AccountKeyTypeFail,
		accountkey.AccountKeyTypeWeightedMultiSig,
		accountkey.AccountKeyTypeRoleBased,
	}

	txTypes := []types.TxType{}
	for i := types.TxTypeLegacyTransaction; i < types.TxTypeEthereumLast; i++ {
		if i == types.TxTypeKlaytnLast {
			i = types.TxTypeEthereumAccessList
		}

		_, err := types.NewTxInternalData(i)
		if err == nil {
			txTypes = append(txTypes, i)
		}
	}

	// tests for all accountKeyTypes
	for _, accountKeyType := range accountKeyTypes {
		// a sender account
		sender, err := createDefaultAccount(accountKeyType)
		assert.Equal(t, nil, err)

		// senderLegacy provides a coupled (address, key pair) will be used by sender
		senderLegacy, err := createAnonymousAccount(getRandomPrivateKeyString(t))
		assert.Equal(t, nil, err)

		// assign senderLegacy address to sender
		sender.Addr = senderLegacy.Addr

		if testing.Verbose() {
			fmt.Println("reservoirAddr = ", reservoir.Addr.String())
			fmt.Println("senderAddr = ", sender.Addr.String())
		}

		// send KLAY to sender
		{
			var txs types.Transactions

			amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
			tx := types.NewTransaction(reservoir.GetNonce(),
				sender.Addr, amount, gasLimit, gasPrice, []byte{})

			err := tx.SignWithKeys(signer, reservoir.Keys)
			assert.Equal(t, nil, err)
			txs = append(txs, tx)

			if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
				t.Fatal(err)
			}
			reservoir.AddNonce()
		}

		if senderLegacy.AccKey.Type() != accountKeyType {
			// update sender's account key
			{
				var txs types.Transactions

				values := map[types.TxValueKeyType]interface{}{
					types.TxValueKeyNonce:      sender.Nonce,
					types.TxValueKeyFrom:       sender.Addr,
					types.TxValueKeyGasLimit:   gasLimit,
					types.TxValueKeyGasPrice:   gasPrice,
					types.TxValueKeyAccountKey: sender.AccKey,
				}
				tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
				assert.Equal(t, nil, err)

				err = tx.SignWithKeys(signer, senderLegacy.Keys)
				assert.Equal(t, nil, err)

				txs = append(txs, tx)

				if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
					t.Fatal(err)
				}
				sender.AddNonce()
			}
		} else {
			sender.Keys = senderLegacy.Keys
		}

		// tests for all txTypes
		for _, txType := range txTypes {
			// skip if tx type is legacy transaction and sender is not legacy.
			if (txType.IsLegacyTransaction() || txType.IsEthTypedTransaction()) &&
				!sender.AccKey.Type().IsLegacyAccountKey() {
				continue
			}

			if testing.Verbose() {
				fmt.Println("Testing... accountKeyType: ", accountKeyType, ", txType: ", txType)
			}

			// generate a default transaction
			tx, _, err := generateDefaultTx(sender, reservoir, txType, contractAddr)
			assert.Equal(t, nil, err)

			// sign a tx
			tx, err = signTxWithVariousKeyTypes(signer, tx, sender)
			assert.Equal(t, nil, err)

			if txType.IsFeeDelegatedTransaction() {
				err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
				assert.Equal(t, nil, err)
			}

			expectedError := expectedTestResultForDefaultTx(accountKeyType, txType)

			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, expectedError, err)

			if err == nil {
				assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
			}
		}
	}
	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdateMultiSigKeyMaxKey tests multiSig key update with maximum private keys.
// A multiSig account supports maximum 10 different private keys.
// Update an account key to a multiSig key with 11 different private keys (more than 10 -> failed)
func TestAccountUpdateMultiSigKeyMaxKey(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// anonymous account
	anon, err := createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594389999")
	assert.Equal(t, nil, err)

	// multisig setting
	threshold := uint(10)
	weights := []uint{1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 0, 1}
	multisigAddr := common.HexToAddress("0xbbfa38050bf3167c887c086758f448ce067ea8ea")
	prvKeys := []string{
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380000",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380001",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380002",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380003",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380004",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380005",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380006",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594300007",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594300008",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594300009",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594300010",
	}

	// multi-sig account
	multisig, err := createMultisigAccount(threshold, weights, prvKeys, multisigAddr)
	assert.Equal(t, nil, err)

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("multisigAddr = ", multisig.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// update key to a multiSig account with 11 different private keys (more than 10 -> failed)
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: multisig.AccKey,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrMaxKeysExceed, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, kerrors.ErrMaxKeysExceed, err)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
		}

		anon.Nonce += 1
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdateMultiSigKeyBigThreshold tests multiSig key update with abnormal threshold.
// When a multisig key is updated, a threshold value should be less or equal to the total weight of private keys.
// If not, the account cannot creates any valid signatures.
// The test update an account key to a multisig key with a threshold (10) and the total weight (6). (failed case)
func TestAccountUpdateMultiSigKeyBigThreshold(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// anonymous account
	anon, err := createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594389999")
	assert.Equal(t, nil, err)

	// multisig setting
	threshold := uint(10)
	weights := []uint{1, 2, 3}
	multisigAddr := common.HexToAddress("0xbbfa38050bf3167c887c086758f448ce067ea8ea")
	prvKeys := []string{
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380000",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380001",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380002",
	}

	// multi-sig account
	multisig, err := createMultisigAccount(threshold, weights, prvKeys, multisigAddr)

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("multisigAddr = ", multisig.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// update key to a multisig key with a threshold (10) and the total weight (6). (failed case)
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: multisig.AccKey,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrUnsatisfiableThreshold, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrUnsatisfiableThreshold, err)
		}
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdateMultiSigKeyDupPrvKeys tests multiSig key update with duplicated private keys.
// A multisig key consists of  all different private keys, therefore account update with duplicated private keys should be failed.
// The test supposed the case when two same private keys are used in creation processes.
func TestAccountUpdateMultiSigKeyDupPrvKeys(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// anonymous account
	anon, err := createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594389999")
	assert.Equal(t, nil, err)

	// the case when two same private keys are used in creation processes.
	threshold := uint(2)
	weights := []uint{1, 1, 2}
	multisigAddr := common.HexToAddress("0xbbfa38050bf3167c887c086758f448ce067ea8ea")
	prvKeys := []string{
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380000",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380001",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380001",
	}

	// multi-sig account
	multisig, err := createMultisigAccount(threshold, weights, prvKeys, multisigAddr)
	assert.Equal(t, nil, err)

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// 1. Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// 2. Update to a multisig key which has two same private keys.
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: multisig.AccKey,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrDuplicatedKey, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrDuplicatedKey, err)
		}
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdateMultiSigKeyWeightOverflow tests multiSig key update with weight overflow.
// If the sum of weights is overflowed, the test should fail.
func TestAccountUpdateMultiSigKeyWeightOverflow(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// Simply check & set the maximum value of uint
	MAX := uint(math.MaxUint32)
	if strconv.IntSize == 64 {
		MAX = math.MaxUint64
	}

	// anonymous account
	anon, err := createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594389999")
	assert.Equal(t, nil, err)

	// multisig setting
	threshold := uint(MAX)
	weights := []uint{MAX / 2, MAX / 2, MAX / 2}
	multisigAddr := common.HexToAddress("0xbbfa38050bf3167c887c086758f448ce067ea8ea")
	prvKeys := []string{
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380000",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380001",
		"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594380002",
	}

	// multi-sig account
	multisig, err := createMultisigAccount(threshold, weights, prvKeys, multisigAddr)
	assert.Equal(t, nil, err)

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// 1. Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// 2. update toc a multisig key with a threshold, uint(MAX), and the total weight, uint(MAX/2)*3. (failed case)
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: multisig.AccKey,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrWeightedSumOverflow, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrWeightedSumOverflow, err)
		}
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdateRoleBasedKeyInvalidNumKey tests account update with a RoleBased key which contains invalid number of sub-keys.
// A RoleBased key can contain 1 ~ 3 sub-keys, otherwise it will fail to the account creation.
// 1. try to create an account with a RoleBased key which contains 4 sub-keys.
// 2. try to create an account with a RoleBased key which contains 0 sub-key.
func TestAccountUpdateRoleBasedKeyInvalidNumKey(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// anonymous account
	anon, err := createAnonymousAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	assert.Equal(t, nil, err)

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("anonAddr = ", anon.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// 1. Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// 2. update to a RoleBased key which contains 4 sub-keys.
	{
		keys := genTestKeys(4)
		roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
			accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
			accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
			accountkey.NewAccountKeyPublicWithValue(&keys[2].PublicKey),
			accountkey.NewAccountKeyPublicWithValue(&keys[3].PublicKey),
		})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrLengthTooLong, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrLengthTooLong, err)
		}
	}

	// 2. update to a RoleBased key which contains 0 sub-key.
	{
		roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrZeroLength, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrZeroLength, err)
		}
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdateRoleBasedKeyInvalidTypeKey tests account key update with a RoleBased key contains types of sub-keys.
// As a sub-key type, a RoleBased key can have AccountKeyFail keys but not AccountKeyNil keys.
// 1. a RoleBased key contains an AccountKeyNil type sub-key as a first sub-key. (fail)
// 2. a RoleBased key contains an AccountKeyNil type sub-key as a second sub-key. (fail)
// 3. a RoleBased key contains an AccountKeyNil type sub-key as a third sub-key. (fail)
// 4. a RoleBased key contains an AccountKeyFail type sub-key as a first sub-key. (success)
// 5. a RoleBased key contains an AccountKeyFail type sub-key as a second sub-key. (success)
// 6. a RoleBased key contains an AccountKeyFail type sub-key as a third sub-key. (success)
func TestAccountUpdateRoleBasedKeyInvalidTypeKey(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// anonymous account
	anon, err := createAnonymousAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	assert.Equal(t, nil, err)

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("anonAddr = ", anon.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)
	keys := genTestKeys(2)

	// 0. Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// 1. a RoleBased key contains an AccountKeyNil type sub-key as a first sub-key. (fail)
	{
		roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
			accountkey.NewAccountKeyNil(),
			accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
			accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
		})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrAccountKeyNilUninitializable, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrAccountKeyNilUninitializable, err)
		}
	}

	// 2. a RoleBased key contains an AccountKeyNil type sub-key as a second sub-key. (fail)
	{
		roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
			accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
			accountkey.NewAccountKeyNil(),
			accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
		})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrAccountKeyNilUninitializable, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrAccountKeyNilUninitializable, err)
		}
	}

	// 3. a RoleBased key contains an AccountKeyNil type sub-key as a third sub-key. (fail)
	{
		roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
			accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
			accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
			accountkey.NewAccountKeyNil(),
		})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrAccountKeyNilUninitializable, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrAccountKeyNilUninitializable, err)
		}
	}

	// 4. a RoleBased key contains an AccountKeyFail type sub-key as a first sub-key. (success)
	{
		roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
			accountkey.NewAccountKeyFail(),
			accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
			accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
		})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, nil, err)
		assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	// 5. a RoleBased key contains an AccountKeyFail type sub-key as a second sub-key. (success)
	{
		roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
			accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
			accountkey.NewAccountKeyFail(),
			accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
		})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, nil, err)
		assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	// 6. a RoleBased key contains an AccountKeyFail type sub-key as a third sub-key. (success)
	{
		roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
			accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
			accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
			accountkey.NewAccountKeyFail(),
		})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, nil, err)
		assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdateWithRoleBasedKey tests account update with a roleBased key.
// A roleBased key contains three types of sub-keys, and only RoleAccountUpdate key is used for update.
// Other sub-keys are not used for the account update.
// 0. create an account and update its key to a roleBased key.
// 1. try to update the account with a RoleTransaction key. (fail)
// 2. try to update the account with a RoleFeePayer key. (fail)
// 3. try to update the account with a RoleAccountUpdate key. (success)
func TestAccountUpdateRoleBasedKey(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// anonymous account
	anon, err := createAnonymousAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	assert.Equal(t, nil, err)

	// generate a roleBased key
	keys := genTestKeys(3)
	roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
		accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&keys[2].PublicKey),
	})

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("anonAddr = ", anon.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// update the account with a roleBased key.
	{
		var txs types.Transactions
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		anon.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// 1. try to update the account with a RoleTransaction key. (fail)
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{keys[accountkey.RoleTransaction]})
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, types.ErrInvalidSigSender, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, types.ErrInvalidSigSender, err)
		}
	}

	// 2. try to update the account with a RoleFeePayer key. (fail)
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{keys[accountkey.RoleFeePayer]})
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, types.ErrInvalidSigSender, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, types.ErrInvalidSigSender, err)
		}
	}

	// 3. try to update the account with a RoleAccountUpdate key. (success)
	{
		var txs types.Transactions
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{keys[accountkey.RoleAccountUpdate]})
		assert.Equal(t, nil, err)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		anon.Nonce += 1
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdateRoleBasedKeyNested tests account update with a nested RoleBasedKey.
// Nested RoleBasedKey is not allowed in Klaytn.
// 1. Create an account with a RoleBasedKey.
// 2. Update an accountKey with a nested RoleBasedKey
func TestAccountUpdateRoleBasedKeyNested(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// anonymous account
	anon, err := createAnonymousAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	assert.Equal(t, nil, err)

	// roleBasedKeys and a nested roleBasedKey
	roleKey, err := createDefaultAccount(accountkey.AccountKeyTypeRoleBased)
	assert.Equal(t, nil, err)

	nestedAccKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
		roleKey.AccKey,
	})

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("roleAddr = ", roleKey.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			anon.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// update the account with a roleBased key.
	{
		var txs types.Transactions
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey.AccKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		anon.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// 2. Update an accountKey with a nested RoleBasedKey.
	{
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: nestedAccKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{roleKey.Keys[accountkey.RoleAccountUpdate]})
		assert.Equal(t, nil, err)

		// For tx pool validation test
		{
			err = txpool.AddRemote(tx)
			assert.Equal(t, kerrors.ErrNestedCompositeType, err)
		}

		// For block tx validation test
		{
			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.Equal(t, (*types.Receipt)(nil), receipt)
			assert.Equal(t, kerrors.ErrNestedCompositeType, err)
		}
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestRoleBasedKeySendTx tests signing transactions with a role-based key.
// A role-based key contains three types of sub-keys: RoleTransaction, RoleAccountUpdate, RoleFeePayer.
// Only RoleTransaction can generate valid signature as a sender except account update txs.
// RoleAccountUpdate can generate valid signature for account update txs.
func TestRoleBasedKeySendTx(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	gasPrice := new(big.Int).SetUint64(25 * params.Ston)

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// main account with a role-based key
	roleBased, err := createAnonymousAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	assert.Equal(t, nil, err)

	// smart contract account
	contractAddr := common.Address{}

	// generate a role-based key
	prvKeys := genTestKeys(3)
	roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[0].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[1].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[2].PublicKey),
	})

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("roleBasedAddr = ", roleBased.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	txTypes := []types.TxType{}
	for i := types.TxTypeLegacyTransaction; i < types.TxTypeEthereumLast; i++ {
		if i == types.TxTypeKlaytnLast {
			i = types.TxTypeEthereumAccessList
		}

		if i.IsLegacyTransaction() || i.IsEthTypedTransaction() {
			continue // accounts with role-based key cannot send the legacy tx and ethereum typed tx.
		}
		_, err := types.NewTxInternalData(i)
		if err == nil {
			txTypes = append(txTypes, i)
		}
	}

	// deploy a contract to test smart contract execution.
	{
		var txs types.Transactions
		valueMap, _ := genMapForTxTypes(reservoir, reservoir, types.TxTypeSmartContractDeploy)
		valueMap[types.TxValueKeyTo] = (*common.Address)(nil)

		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, valueMap)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		contractAddr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)
		reservoir.Nonce += 1
	}

	// transfer (reservoir -> roleBased) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			roleBased.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// update to an roleBased account with a role-based key.
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      roleBased.Nonce,
			types.TxValueKeyFrom:       roleBased.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, roleBased.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		roleBased.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// test fee delegation txs for each role of role-based key.
	// only RoleFeePayer type can generate valid signature as a fee payer.
	for keyType, key := range prvKeys {
		for _, txType := range txTypes {
			valueMap, _ := genMapForTxTypes(roleBased, reservoir, txType)
			valueMap[types.TxValueKeyGasLimit] = uint64(1000000)

			if txType.IsFeeDelegatedTransaction() {
				valueMap[types.TxValueKeyFeePayer] = reservoir.Addr
			}

			// Currently, test VM is not working properly when the GasPrice is not 0.
			basicType := toBasicType(txType)
			if keyType == int(accountkey.RoleTransaction) {
				if basicType == types.TxTypeSmartContractDeploy || basicType == types.TxTypeSmartContractExecution {
					valueMap[types.TxValueKeyGasPrice] = new(big.Int).SetUint64(0)
				}
			}

			if basicType == types.TxTypeSmartContractExecution {
				valueMap[types.TxValueKeyTo] = contractAddr
			}

			tx, err := types.NewTransactionWithMap(txType, valueMap)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{key})
			assert.Equal(t, nil, err)

			if txType.IsFeeDelegatedTransaction() {
				err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
				assert.Equal(t, nil, err)
			}

			// Only RoleTransaction can generate valid signature as a sender except account update txs.
			// RoleAccountUpdate can generate valid signature for account update txs.
			if keyType == int(accountkey.RoleAccountUpdate) && txType.IsAccountUpdate() ||
				keyType == int(accountkey.RoleTransaction) && !txType.IsAccountUpdate() {
				// Do not make a block since account update tx can change sender's keys.
				receipt, _, err := applyTransaction(t, bcdata, tx)
				assert.Equal(t, nil, err)
				assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
			} else {
				// For tx pool validation test
				{
					err = txpool.AddRemote(tx)
					assert.Equal(t, types.ErrInvalidSigSender, err)
				}

				// For block tx validation test
				{
					receipt, _, err := applyTransaction(t, bcdata, tx)
					assert.Equal(t, types.ErrInvalidSigSender, err)
					assert.Equal(t, (*types.Receipt)(nil), receipt)
				}
			}
		}
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestRoleBasedKeyFeeDelegation tests fee delegation with a role-based key.
// A role-based key contains three types of sub-keys: RoleTransaction, RoleAccountUpdate, RoleFeePayer.
// Only RoleFeePayer can sign txs as a fee payer.
func TestRoleBasedKeyFeeDelegation(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	gasPrice := new(big.Int).SetUint64(25 * params.Ston)

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// main account with a role-based key
	roleBased, err := createAnonymousAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	assert.Equal(t, nil, err)

	// smart contract account
	contractAddr := common.Address{}

	// generate a role-based key
	prvKeys := genTestKeys(3)
	roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[0].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[1].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[2].PublicKey),
	})

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("roleBasedAddr = ", roleBased.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	feeTxTypes := []types.TxType{
		types.TxTypeFeeDelegatedValueTransfer,
		types.TxTypeFeeDelegatedValueTransferMemo,
		types.TxTypeFeeDelegatedSmartContractDeploy,
		types.TxTypeFeeDelegatedSmartContractExecution,
		types.TxTypeFeeDelegatedAccountUpdate,
		types.TxTypeFeeDelegatedCancel,

		types.TxTypeFeeDelegatedValueTransferWithRatio,
		types.TxTypeFeeDelegatedValueTransferMemoWithRatio,
		types.TxTypeFeeDelegatedSmartContractDeployWithRatio,
		types.TxTypeFeeDelegatedSmartContractExecutionWithRatio,
		types.TxTypeFeeDelegatedAccountUpdateWithRatio,
		types.TxTypeFeeDelegatedCancelWithRatio,
	}

	// deploy a contract to test smart contract execution.
	{
		var txs types.Transactions
		valueMap, _ := genMapForTxTypes(reservoir, reservoir, types.TxTypeSmartContractDeploy)
		valueMap[types.TxValueKeyTo] = (*common.Address)(nil)

		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, valueMap)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		contractAddr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1
	}

	// transfer (reservoir -> roleBased) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KLAY))
		tx := types.NewTransaction(reservoir.Nonce,
			roleBased.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)
		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// update to an roleBased account with a role-based key.
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      roleBased.Nonce,
			types.TxValueKeyFrom:       roleBased.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, roleBased.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		roleBased.Nonce += 1
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// test fee delegation txs for each role of role-based key.
	// only RoleFeePayer type can generate valid signature as a fee payer.
	for keyType, key := range prvKeys {
		for _, txType := range feeTxTypes {
			valueMap, _ := genMapForTxTypes(reservoir, reservoir, txType)
			valueMap[types.TxValueKeyFeePayer] = roleBased.GetAddr()
			valueMap[types.TxValueKeyGasLimit] = uint64(1000000)

			// Currently, test VM is not working properly when the GasPrice is not 0.
			basicType := toBasicType(txType)
			if keyType == int(accountkey.RoleFeePayer) {
				if basicType == types.TxTypeSmartContractDeploy || basicType == types.TxTypeSmartContractExecution {
					valueMap[types.TxValueKeyGasPrice] = new(big.Int).SetUint64(0)
				}
			}

			if basicType == types.TxTypeSmartContractExecution {
				valueMap[types.TxValueKeyTo] = contractAddr
			}

			tx, err := types.NewTransactionWithMap(txType, valueMap)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.Keys)
			assert.Equal(t, nil, err)

			err = tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{key})
			assert.Equal(t, nil, err)

			if keyType == int(accountkey.RoleFeePayer) {
				// Do not make a block since account update tx can change sender's keys.
				receipt, _, err := applyTransaction(t, bcdata, tx)
				assert.Equal(t, nil, err)
				assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
			} else {
				// For tx pool validation test
				{
					err = txpool.AddRemote(tx)
					assert.Equal(t, blockchain.ErrInvalidFeePayer, err)
				}

				// For block tx validation test
				{
					receipt, _, err := applyTransaction(t, bcdata, tx)
					assert.Equal(t, types.ErrInvalidSigFeePayer, err)
					assert.Equal(t, (*types.Receipt)(nil), receipt)
				}
			}
		}
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

func TestAccountKeyUpdateLegacyToPublic(t *testing.T) {
	gasPrice := new(big.Int).SetUint64(25 * params.Ston)
	gasLimit := uint64(1000000)

	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(6, 4)
	if err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	var txs types.Transactions

	// make TxPool to test validation in 'TxPool add' process
	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	{
		addr := *bcdata.addrs[1]
		feepayer := *bcdata.addrs[0]
		acckey := accountkey.NewAccountKeyWeightedMultiSigWithValues(1,
			accountkey.WeightedPublicKeys{
				accountkey.NewWeightedPublicKey(1, (*accountkey.PublicKeySerializable)(&bcdata.privKeys[1].PublicKey)),
				accountkey.NewWeightedPublicKey(1, (*accountkey.PublicKeySerializable)(&bcdata.privKeys[2].PublicKey)),
			})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      uint64(0),
			types.TxValueKeyFrom:       addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: acckey,
			types.TxValueKeyFeePayer:   feepayer,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[1]})
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)
	}
	{
		addr := *bcdata.addrs[1]
		feepayer := *bcdata.addrs[0]
		acckey := accountkey.NewAccountKeyWeightedMultiSigWithValues(1,
			accountkey.WeightedPublicKeys{
				accountkey.NewWeightedPublicKey(1, (*accountkey.PublicKeySerializable)(&bcdata.privKeys[1].PublicKey)),
				accountkey.NewWeightedPublicKey(1, (*accountkey.PublicKeySerializable)(&bcdata.privKeys[3].PublicKey)),
			})

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      uint64(1),
			types.TxValueKeyFrom:       addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: acckey,
			types.TxValueKeyFeePayer:   feepayer,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[1]})
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, []*ecdsa.PrivateKey{bcdata.privKeys[0]})
		assert.Equal(t, nil, err)

		txs = append(txs, tx)
	}
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}

	// select account key types to be tested
	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}
