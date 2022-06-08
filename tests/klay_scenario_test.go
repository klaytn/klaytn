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
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/compiler"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

var (
	to       = common.HexToAddress("7b65B75d204aBed71587c9E519a89277766EE1d0")
	feePayer = common.HexToAddress("5A0043070275d9f6054307Ee7348bD660849D90f")
)

const (
	nonce    = uint64(1234)
	gasLimit = uint64(100000000000)
)

type TestAccountType struct {
	Addr   common.Address
	Keys   []*ecdsa.PrivateKey
	Nonce  uint64
	AccKey accountkey.AccountKey
}

type TestCreateMultisigAccountParam struct {
	Threshold uint
	Weights   []uint
	PrvKeys   []string
}

// createAnonymousAccount creates an account whose address is derived from the private key.
func createAnonymousAccount(prvKeyHex string) (*TestAccountType, error) {
	key, err := crypto.HexToECDSA(prvKeyHex)
	if err != nil {
		return nil, err
	}

	addr := crypto.PubkeyToAddress(key.PublicKey)

	return &TestAccountType{
		Addr:   addr,
		Keys:   []*ecdsa.PrivateKey{key},
		Nonce:  uint64(0),
		AccKey: accountkey.NewAccountKeyLegacy(),
	}, nil
}

// createDecoupledAccount creates an account whose address is decoupled with its private key.
func createDecoupledAccount(prvKeyHex string, addr common.Address) (*TestAccountType, error) {
	key, err := crypto.HexToECDSA(prvKeyHex)
	if err != nil {
		return nil, err
	}

	return &TestAccountType{
		Addr:   addr,
		Keys:   []*ecdsa.PrivateKey{key},
		Nonce:  uint64(0),
		AccKey: accountkey.NewAccountKeyPublicWithValue(&key.PublicKey),
	}, nil
}

//

// createMultisigAccount creates an account having multiple of keys.
func createMultisigAccount(threshold uint, weights []uint, prvKeys []string, addr common.Address) (*TestAccountType, error) {
	var err error

	keys := make([]*ecdsa.PrivateKey, len(prvKeys))
	weightedKeys := make(accountkey.WeightedPublicKeys, len(prvKeys))

	for i, p := range prvKeys {
		keys[i], err = crypto.HexToECDSA(p)
		if err != nil {
			return nil, err
		}
		weightedKeys[i] = accountkey.NewWeightedPublicKey(weights[i], (*accountkey.PublicKeySerializable)(&keys[i].PublicKey))
	}

	return &TestAccountType{
		Addr:   addr,
		Keys:   keys,
		Nonce:  uint64(0),
		AccKey: accountkey.NewAccountKeyWeightedMultiSigWithValues(threshold, weightedKeys),
	}, nil
}

// createRoleBasedAccountWithAccountKeyPublic creates an account having keys that have role with AccountKeyPublic.
func createRoleBasedAccountWithAccountKeyPublic(prvKeys []string, addr common.Address) (*TestRoleBasedAccountType, error) {
	var err error

	if len(prvKeys) != 3 {
		return nil, errors.New("Need three key value for create role-based account")
	}

	keys := make([]*ecdsa.PrivateKey, len(prvKeys))

	for i, p := range prvKeys {
		keys[i], err = crypto.HexToECDSA(p)
		if err != nil {
			return nil, err
		}
	}

	accKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
		accountkey.NewAccountKeyPublicWithValue(&keys[0].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&keys[1].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&keys[2].PublicKey),
	})

	return &TestRoleBasedAccountType{
		Addr:       addr,
		TxKeys:     []*ecdsa.PrivateKey{keys[0]},
		UpdateKeys: []*ecdsa.PrivateKey{keys[1]},
		FeeKeys:    []*ecdsa.PrivateKey{keys[2]},
		Nonce:      uint64(0),
		AccKey:     accKey,
	}, nil
}

// createRoleBasedAccountWithAccountKeyWeightedMultisig creates an account having keys that have role with AccountKeyWeightedMultisig.
func createRoleBasedAccountWithAccountKeyWeightedMultiSig(multisigs []TestCreateMultisigAccountParam, addr common.Address) (*TestRoleBasedAccountType, error) {
	var err error

	if len(multisigs) != 3 {
		return nil, errors.New("Need three key value for create role-based account")
	}

	prvKeys := make([][]*ecdsa.PrivateKey, len(multisigs))
	multisigKeys := make([]*accountkey.AccountKeyWeightedMultiSig, len(multisigs))

	for idx, multisig := range multisigs {
		keys := make([]*ecdsa.PrivateKey, len(multisig.PrvKeys))
		weightedKeys := make(accountkey.WeightedPublicKeys, len(multisig.PrvKeys))

		for i, p := range multisig.PrvKeys {
			keys[i], err = crypto.HexToECDSA(p)
			if err != nil {
				return nil, err
			}
			weightedKeys[i] = accountkey.NewWeightedPublicKey(multisig.Weights[i], (*accountkey.PublicKeySerializable)(&keys[i].PublicKey))
		}
		prvKeys[idx] = keys
		multisigKeys[idx] = accountkey.NewAccountKeyWeightedMultiSigWithValues(multisig.Threshold, weightedKeys)
	}

	accKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{multisigKeys[0], multisigKeys[1], multisigKeys[2]})

	return &TestRoleBasedAccountType{
		Addr:       addr,
		TxKeys:     prvKeys[0],
		UpdateKeys: prvKeys[1],
		FeeKeys:    prvKeys[2],
		Nonce:      uint64(0),
		AccKey:     accKey,
	}, nil
}

// TestFeeDelegatedWithSmallBalance tests the case that an account having a small amount of tokens transfers
// all the tokens to another account with a fee payer.
// This kinds of transactions were discarded in TxPool.promoteExecutable() because the total cost of
// the transaction is larger than the amount of tokens the sender has.
// Since we provide fee-delegated transactions, it is not true in the above case.
// This test code should succeed.
func TestFeeDelegatedWithSmallBalance(t *testing.T) {
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
	gasPrice := big.NewInt(25 * params.Ston)

	// 1. Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(10000)
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

	// 2. Transfer (anon -> reservoir) using a TxTypeFeeDelegatedValueTransfer
	{
		amount := new(big.Int).SetUint64(10000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    anon.Nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyFeePayer: reservoir.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		p := makeTxPool(bcdata, 10)

		p.AddRemote(tx)

		if err := bcdata.GenABlockWithTxpool(accountMap, p, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1
	}

	state, err := bcdata.bc.State()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(0), state.GetBalance(anon.Addr).Uint64())
}

// TestSmartContractDeployAddress checks that the smart contract is deployed to the given address or not by
// checking receipt.ContractAddress.
func TestSmartContractDeployAddress(t *testing.T) {
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

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
	}

	contractAddr := common.Address{}

	gasPrice := new(big.Int).SetUint64(0)
	gasLimit := uint64(100000000000)

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	var code string

	if isCompilerAvailable() {
		filename := string("../contracts/reward/contract/KlaytnReward.sol")
		codes, _ := compileSolidity(filename)
		code = codes[0]
	} else {
		// Falling back to use compiled code.
		code = "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"
	}

	// 1. Deploy smart contract (reservoir -> contract)
	{
		amount := new(big.Int).SetUint64(0)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.Nonce,
			types.TxValueKeyFrom:          reservoir.Addr,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyData:          common.FromHex(code),
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		contractAddr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		// check receipt
		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, nil, err)
		assert.Equal(t, contractAddr, receipt.ContractAddress)
	}
}

// TestSmartContractScenario tests the following scenario:
// 1. Deploy smart contract (reservoir -> contract)
// 2. Check the smart contract is deployed well.
// 3. Execute "reward" function with amountToSend
// 4. Validate "reward" function is executed correctly by executing "balanceOf".
func TestSmartContractScenario(t *testing.T) {
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

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
	}

	contractAddr := common.Address{}

	gasPrice := new(big.Int).SetUint64(0)
	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	var abiStr string
	var code string

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

	// 1. Deploy smart contract (reservoir -> contract)
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(0)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.Nonce,
			types.TxValueKeyFrom:          reservoir.Addr,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
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

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		contractAddr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1
	}

	// 2. Check the smart contract is deployed well.
	{
		statedb, err := bcdata.bc.State()
		if err != nil {
			t.Fatal(err)
		}
		code := statedb.GetCode(contractAddr)
		assert.Equal(t, 478, len(code))
	}

	// 3. Execute "reward" function with amountToSend
	amountToSend := new(big.Int).SetUint64(10)
	{
		var txs types.Transactions

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("reward", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       contractAddr,
			types.TxValueKeyAmount:   amountToSend,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     data,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// 4. Validate "reward" function is executed correctly by executing "balanceOf".
	{
		amount := new(big.Int).SetUint64(0)

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("balanceOf", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       contractAddr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     data,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		assert.Equal(t, nil, err)

		balance := new(big.Int)
		abii.Unpack(&balance, "balanceOf", ret)

		assert.Equal(t, amountToSend, balance)
		reservoir.Nonce += 1
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestSmartContractSign tests value transfer and fee delegation of smart contract accounts.
// It performs the following scenario:
// 1. Deploy smart contract (reservoir -> contract)
// 2. Check the smart contract is deployed well.
// 3. Try value transfer. It should be failed.
// 4. Try fee delegation. It should be failed.
func TestSmartContractSign(t *testing.T) {
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

	reservoir2 := &TestAccountType{
		Addr:  *bcdata.addrs[1],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[1]},
		Nonce: uint64(0),
	}

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
	}

	contract, err := createAnonymousAccount("ed34b0cf47a0021e9897760f0a904a69260c2f638e0bcc805facb745ec3ff9ab")
	assert.Equal(t, nil, err)

	gasPrice := new(big.Int).SetUint64(0)
	gasLimit := uint64(100000000000)

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	var code string

	if isCompilerAvailable() {
		filename := string("../contracts/reward/contract/KlaytnReward.sol")
		codes, _ := compileSolidity(filename)
		code = codes[0]
	} else {
		// Falling back to use compiled code.
		code = "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"
	}

	// 1. Deploy smart contract (reservoir -> contract)
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(0)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.Nonce,
			types.TxValueKeyFrom:          reservoir.Addr,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
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

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		contract.Addr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1
	}

	// 2. Check the smart contract is deployed well.
	{
		statedb, err := bcdata.bc.State()
		if err != nil {
			t.Fatal(err)
		}
		code := statedb.GetCode(contract.Addr)
		assert.Equal(t, 478, len(code))
	}

	// 3. Try value transfer. It should be failed.
	{
		amount := new(big.Int).SetUint64(100000000000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    contract.Nonce,
			types.TxValueKeyFrom:     contract.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)
		err = tx.SignWithKeys(signer, contract.Keys)
		assert.Equal(t, nil, err)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, (*types.Receipt)(nil), receipt)
		assert.Equal(t, types.ErrInvalidSigSender, err)
	}

	// 4. Try fee delegation. It should be failed.
	{
		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       reservoir2.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: contract.Addr,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, contract.Keys)
		assert.Equal(t, nil, err)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, (*types.Receipt)(nil), receipt)
		assert.Equal(t, types.ErrInvalidSigFeePayer, err)
	}
}

// TestFeeDelegatedSmartContractScenario tests the following scenario:
// 1. Deploy smart contract (reservoir -> contract) with fee-delegation.
// 2. Check the smart contract is deployed well.
// 3. Execute "reward" function with amountToSend with fee-delegation.
// 4. Validate "reward" function is executed correctly by executing "balanceOf".
func TestFeeDelegatedSmartContractScenario(t *testing.T) {
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

	reservoir2 := &TestAccountType{
		Addr:  *bcdata.addrs[1],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[1]},
		Nonce: uint64(0),
	}

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
	}

	contractAddr := common.Address{}

	gasPrice := new(big.Int).SetUint64(0)
	gasLimit := uint64(100000000000)

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	var abiStr string
	var code string

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

	// 1. Deploy smart contract (reservoir -> contract) with fee-delegation.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(0)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.Nonce,
			types.TxValueKeyFrom:          reservoir.Addr,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyData:          common.FromHex(code),
			types.TxValueKeyFeePayer:      reservoir2.Addr,
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		contractAddr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1
	}

	// 2. Check the smart contract is deployed well.
	{
		statedb, err := bcdata.bc.State()
		if err != nil {
			t.Fatal(err)
		}
		code := statedb.GetCode(contractAddr)
		assert.Equal(t, 478, len(code))
	}

	// 3. Execute "reward" function with amountToSend with fee-delegation.
	amountToSend := new(big.Int).SetUint64(10)
	{
		var txs types.Transactions

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("reward", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       contractAddr,
			types.TxValueKeyAmount:   amountToSend,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     data,
			types.TxValueKeyFeePayer: reservoir2.Addr,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecution, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// 4. Validate "reward" function is executed correctly by executing "balanceOf".
	{
		amount := new(big.Int).SetUint64(0)

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("balanceOf", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       contractAddr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     data,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		assert.Equal(t, nil, err)

		balance := new(big.Int)
		abii.Unpack(&balance, "balanceOf", ret)

		assert.Equal(t, amountToSend, balance)
		reservoir.Nonce += 1
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestFeeDelegatedSmartContractScenarioWithRatio tests the following scenario:
// 1. Deploy smart contract (reservoir -> contract) with fee-delegation.
// 2. Check the smart contract is deployed well.
// 3. Execute "reward" function with amountToSend with fee-delegation.
// 4. Validate "reward" function is executed correctly by executing "balanceOf".
func TestFeeDelegatedSmartContractScenarioWithRatio(t *testing.T) {
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

	reservoir2 := &TestAccountType{
		Addr:  *bcdata.addrs[1],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[1]},
		Nonce: uint64(0),
	}

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
	}

	contractAddr := common.Address{}

	gasPrice := new(big.Int).SetUint64(0)
	gasLimit := uint64(100000000000)

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)

	var abiStr string
	var code string

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

	// 1. Deploy smart contract (reservoir -> contract) with fee-delegation.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(0)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 (*common.Address)(nil),
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyHumanReadable:      false,
			types.TxValueKeyData:               common.FromHex(code),
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
			types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		contractAddr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1
	}

	// 2. Check the smart contract is deployed well.
	{
		statedb, err := bcdata.bc.State()
		if err != nil {
			t.Fatal(err)
		}
		code := statedb.GetCode(contractAddr)
		assert.Equal(t, 478, len(code))
	}

	// 3. Execute "reward" function with amountToSend with fee-delegation.
	amountToSend := new(big.Int).SetUint64(10)
	{
		var txs types.Transactions

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("reward", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              reservoir.Nonce,
			types.TxValueKeyFrom:               reservoir.Addr,
			types.TxValueKeyTo:                 contractAddr,
			types.TxValueKeyAmount:             amountToSend,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               data,
			types.TxValueKeyFeePayer:           reservoir2.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecutionWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir2.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		reservoir.Nonce += 1
	}

	// 4. Validate "reward" function is executed correctly by executing "balanceOf".
	{
		amount := new(big.Int).SetUint64(0)

		abii, err := abi.JSON(strings.NewReader(string(abiStr)))
		assert.Equal(t, nil, err)

		data, err := abii.Pack("balanceOf", reservoir.Addr)
		assert.Equal(t, nil, err)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    reservoir.Nonce,
			types.TxValueKeyFrom:     reservoir.Addr,
			types.TxValueKeyTo:       contractAddr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     data,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		ret, err := callContract(bcdata, tx)
		assert.Equal(t, nil, err)

		balance := new(big.Int)
		abii.Unpack(&balance, "balanceOf", ret)

		assert.Equal(t, amountToSend, balance)
		reservoir.Nonce += 1
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestAccountUpdate tests a following scenario:
// 1. Transfer (reservoir -> anon) using a legacy transaction.
// 2. Key update of decoupled using AccountUpdate
// 3. Transfer (anon -> reservoir) using TxTypeValueTransfer.
// 4. Key update of anon using AccountUpdate with multisig keys.
// 5. Transfer (anon-> reservoir) using TxTypeValueTransfer.
func TestAccountUpdate(t *testing.T) {
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
	gasLimit := uint64(100000)

	// 1. Transfer (reservoir -> anon) using a legacy transaction.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(params.KLAY)
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

	// 2. Key update of anon using AccountUpdate
	{
		var txs types.Transactions

		newKey, err := crypto.HexToECDSA("41bd2b972564206658eab115f26ff4db617e6eb39c81a557adc18d8305d2f867")
		if err != nil {
			t.Fatal(err)
		}

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&newKey.PublicKey),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1

		anon.Keys = []*ecdsa.PrivateKey{newKey}
	}

	// 3. Transfer (anon -> reservoir) using TxTypeValueTransfer.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    anon.Nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1
	}

	// 4. Key update of anon using AccountUpdate with multisig keys.
	{
		var txs types.Transactions

		k1, err := crypto.HexToECDSA("41bd2b972564206658eab115f26ff4db617e6eb39c81a557adc18d8305d2f867")
		if err != nil {
			t.Fatal(err)
		}
		k2, err := crypto.HexToECDSA("c64f2cd1196e2a1791365b00c4bc07ab8f047b73152e4617c6ed06ac221a4b0c")
		if err != nil {
			panic(err)
		}
		threshold := uint(2)
		keys := accountkey.WeightedPublicKeys{
			accountkey.NewWeightedPublicKey(1, (*accountkey.PublicKeySerializable)(&k1.PublicKey)),
			accountkey.NewWeightedPublicKey(1, (*accountkey.PublicKeySerializable)(&k2.PublicKey)),
		}
		newKey := accountkey.NewAccountKeyWeightedMultiSigWithValues(threshold, keys)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: newKey,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1

		anon.Keys = []*ecdsa.PrivateKey{k1, k2}
		anon.AccKey = newKey
	}

	// 5. Transfer (anon-> reservoir) using TxTypeValueTransfer.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(10000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    anon.Nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestFeeDelegatedAccountUpdate tests a following scenario:
// 1. Transfer (reservoir -> anon) using a legacy transaction.
// 2. Key update of anon using TxTypeFeeDelegatedAccountUpdate
// 3. Transfer (anon -> reservoir) using TxTypeValueTransfer.
func TestFeeDelegatedAccountUpdate(t *testing.T) {
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

	// 2. Key update of anon using TxTypeFeeDelegatedAccountUpdate
	{
		var txs types.Transactions

		newKey, err := crypto.HexToECDSA("41bd2b972564206658eab115f26ff4db617e6eb39c81a557adc18d8305d2f867")
		if err != nil {
			t.Fatal(err)
		}

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      anon.Nonce,
			types.TxValueKeyFrom:       anon.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&newKey.PublicKey),
			types.TxValueKeyFeePayer:   reservoir.Addr,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1

		anon.Keys = []*ecdsa.PrivateKey{newKey}
	}

	// 3. Transfer (anon -> reservoir) using TxTypeValueTransfer.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    anon.Nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestFeeDelegatedAccountUpdateWithRatio tests a following scenario:
// 1. Transfer (reservoir -> anon) using a legacy transaction.
// 2. Key update of anon using TxTypeFeeDelegatedAccountUpdateWithRatio.
// 3. Transfer (anon -> reservoir) using TxTypeValueTransfer.
func TestFeeDelegatedAccountUpdateWithRatio(t *testing.T) {
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

	// 2. Key update of decoupled using TxTypeFeeDelegatedAccountUpdateWithRatio.
	{
		var txs types.Transactions

		newKey, err := crypto.HexToECDSA("41bd2b972564206658eab115f26ff4db617e6eb39c81a557adc18d8305d2f867")
		if err != nil {
			t.Fatal(err)
		}

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              anon.Nonce,
			types.TxValueKeyFrom:               anon.Addr,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyAccountKey:         accountkey.NewAccountKeyPublicWithValue(&newKey.PublicKey),
			types.TxValueKeyFeePayer:           reservoir.Addr,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1

		anon.Keys = []*ecdsa.PrivateKey{newKey}
	}

	// 3. Transfer (anon -> reservoir) using TxTypeValueTransfer.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    anon.Nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestMultisigScenario tests a test case for a multi-sig accounts.
// 1. Create an account anon using LegacyTransaction.
// 2. Update the account with multisig key.
// 2. Transfer (anon -> reservoir) using TxTypeValueTransfer.
// 3. Transfer (anon -> reservoir) using TxTypeValueTransfer with only two keys.
// 4. FAILED-CASE: Transfer (anon -> reservoir) using TxTypeValueTransfer with only one key.
func TestMultisigScenario(t *testing.T) {
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

	// multi-sig account
	multisig, err := createMultisigAccount(uint(2),
		[]uint{1, 1, 1},
		[]string{
			"bb113e82881499a7a361e8354a5b68f6c6885c7bcba09ea2b0891480396c322e",
			"a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae59438e989",
			"c32c471b732e2f56103e2f8e8cfd52792ef548f05f326e546a7d1fbf9d0419ec",
		},
		common.HexToAddress("0xbbfa38050bf3167c887c086758f448ce067ea8ea"))

	if testing.Verbose() {
		fmt.Println("reservoirAddr = ", reservoir.Addr.String())
		fmt.Println("multisigAddr = ", multisig.Addr.String())
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// 1. Create an account anon using LegacyTransaction.
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

	// 2. Update the account with multisig key.
	{
		var txs types.Transactions

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

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1

		anon.AccKey = multisig.AccKey
		anon.Keys = multisig.Keys
	}

	// 2. Transfer (anon -> reservoir) using TxTypeValueTransfer.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    anon.Nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1
	}

	// 3. Transfer (anon -> reservoir) using TxTypeValueTransfer with only two keys.
	{
		var txs types.Transactions

		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    anon.Nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys[0:2])
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		anon.Nonce += 1
	}

	// 4. FAILED-CASE: Transfer (anon -> reservoir) using TxTypeValueTransfer with only one key.
	{
		amount := new(big.Int).SetUint64(1000)
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    anon.Nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyTo:       reservoir.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys[:1])
		assert.Equal(t, nil, err)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, types.ErrInvalidSigSender, err)
		assert.Equal(t, (*types.Receipt)(nil), receipt)
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

// TestValidateSender tests ValidateSender with all transaction types.
func TestValidateSender(t *testing.T) {
	// anonymous account
	anon, err := createAnonymousAccount("1da6dfcb52128060cdd2108edb786ca0aff4ef1fa537574286eeabe5c2ebd5ca")
	assert.Equal(t, nil, err)

	// decoupled account
	decoupled, err := createDecoupledAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab",
		common.HexToAddress("0x5104711f7faa9e2dadf593e43db1577a2887636f"))
	assert.Equal(t, nil, err)

	initialBalance := big.NewInt(1000000)

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil)
	statedb.CreateEOA(anon.Addr, false, anon.AccKey)
	statedb.SetNonce(anon.Addr, nonce)
	statedb.SetBalance(anon.Addr, initialBalance)

	statedb.CreateEOA(decoupled.Addr, false, decoupled.AccKey)
	statedb.SetNonce(decoupled.Addr, rand.Uint64())
	statedb.SetBalance(decoupled.Addr, initialBalance)

	signer := types.MakeSigner(params.BFTTestChainConfig, big.NewInt(32))
	gasPrice := new(big.Int).SetUint64(0)
	gasLimit := uint64(100000000000)
	amount := new(big.Int).SetUint64(10000)

	// LegacyTransaction
	{
		amount := new(big.Int).SetUint64(100000000000)
		tx := types.NewTransaction(anon.Nonce,
			decoupled.Addr, amount, gasLimit, gasPrice, []byte{})

		err := tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		_, err = tx.ValidateSender(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, anon.Addr, tx.ValidatedSender())
	}

	// TxTypeValueTransfer
	{
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    nonce,
			types.TxValueKeyFrom:     anon.Addr,
			types.TxValueKeyTo:       decoupled.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		})
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		_, err = tx.ValidateSender(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, anon.Addr, tx.ValidatedSender())
	}

	// TxTypeValueTransfer
	{
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    nonce,
			types.TxValueKeyFrom:     decoupled.Addr,
			types.TxValueKeyTo:       decoupled.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		})
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, decoupled.Keys)
		assert.Equal(t, nil, err)

		_, err = tx.ValidateSender(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, decoupled.Addr, tx.ValidatedSender())
	}

	// TxTypeSmartContractDeploy
	{
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         nonce,
			types.TxValueKeyFrom:          decoupled.Addr,
			types.TxValueKeyTo:            &anon.Addr,
			types.TxValueKeyAmount:        amount,
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyHumanReadable: false,
			// The binary below is a compiled binary of contracts/reward/contract/KlaytnReward.sol.
			types.TxValueKeyData:       common.Hex2Bytes("608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"),
			types.TxValueKeyCodeFormat: params.CodeFormatEVM,
		})
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, decoupled.Keys)
		assert.Equal(t, nil, err)

		_, err = tx.ValidateSender(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, decoupled.Addr, tx.ValidatedSender())
	}

	// TxTypeSmartContractExecution
	{
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    nonce,
			types.TxValueKeyFrom:     decoupled.Addr,
			types.TxValueKeyTo:       anon.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
			// A abi-packed bytes calling "reward" of contracts/reward/contract/KlaytnReward.sol with an address "bc5951f055a85f41a3b62fd6f68ab7de76d299b2".
			types.TxValueKeyData: common.Hex2Bytes("6353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2"),
		})
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, decoupled.Keys)
		assert.Equal(t, nil, err)

		_, err = tx.ValidateSender(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, decoupled.Addr, tx.ValidatedSender())
	}

	// TxTypeChainDataAnchoring
	{
		dummyBlock := types.NewBlock(&types.Header{}, nil, nil)

		scData, err := types.NewAnchoringDataType0(dummyBlock, 0, uint64(dummyBlock.Transactions().Len()))
		if err != nil {
			t.Fatal(err)
		}
		dataAnchoredRLP, _ := rlp.EncodeToBytes(scData)

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:        anon.Nonce,
			types.TxValueKeyFrom:         anon.Addr,
			types.TxValueKeyGasLimit:     gasLimit,
			types.TxValueKeyGasPrice:     gasPrice,
			types.TxValueKeyAnchoredData: dataAnchoredRLP,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeChainDataAnchoring, values)
		if err != nil {
			t.Fatal(err)
		}
		err = tx.SignWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		_, err = tx.ValidateSender(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, anon.Addr, tx.ValidatedSender())
	}

	// TxTypeFeeDelegatedValueTransfer
	{
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    nonce,
			types.TxValueKeyFrom:     decoupled.Addr,
			types.TxValueKeyFeePayer: anon.Addr,
			types.TxValueKeyTo:       decoupled.Addr,
			types.TxValueKeyAmount:   amount,
			types.TxValueKeyGasLimit: gasLimit,
			types.TxValueKeyGasPrice: gasPrice,
		})
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, decoupled.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		_, err = tx.ValidateSender(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, decoupled.Addr, tx.ValidatedSender())

		_, err = tx.ValidateFeePayer(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, anon.Addr, tx.ValidatedFeePayer())
	}

	// TxTypeFeeDelegatedValueTransferWithRatio
	{
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              nonce,
			types.TxValueKeyFrom:               decoupled.Addr,
			types.TxValueKeyFeePayer:           anon.Addr,
			types.TxValueKeyTo:                 decoupled.Addr,
			types.TxValueKeyAmount:             amount,
			types.TxValueKeyGasLimit:           gasLimit,
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(30),
		})
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, decoupled.Keys)
		assert.Equal(t, nil, err)

		err = tx.SignFeePayerWithKeys(signer, anon.Keys)
		assert.Equal(t, nil, err)

		_, err = tx.ValidateSender(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, decoupled.Addr, tx.ValidatedSender())

		_, err = tx.ValidateFeePayer(signer, statedb, 0)
		assert.Equal(t, nil, err)
		assert.Equal(t, anon.Addr, tx.ValidatedFeePayer())
	}
}

func isCompilerAvailable() bool {
	solc, err := compiler.SolidityVersion("")
	if err != nil {
		fmt.Println("Solidity version check failed. Skipping this test", err)
		return false
	}
	if solc.Version != "0.4.24" && solc.Version != "0.4.25" {
		if testing.Verbose() {
			fmt.Println("solc version mismatch. Supported versions are 0.4.24 and 0.4.25.", "version", solc.Version)
		}
		return false
	}

	return true
}

func compileSolidity(filename string) (code []string, abiStr []string) {
	contracts, err := compiler.CompileSolidity("", filename)
	if err != nil {
		panic(err)
	}

	code = make([]string, 0, len(contracts))
	abiStr = make([]string, 0, len(contracts))

	for _, c := range contracts {
		abiBytes, err := json.Marshal(c.Info.AbiDefinition)
		if err != nil {
			panic(err)
		}

		code = append(code, c.Code)
		abiStr = append(abiStr, string(abiBytes))
	}

	return
}

// applyTransaction setups variables to call block.ApplyTransaction() for tests.
// It directly returns values from block.ApplyTransaction().
func applyTransaction(t *testing.T, bcdata *BCData, tx *types.Transaction) (*types.Receipt, uint64, error) {
	state, err := bcdata.bc.State()
	assert.Equal(t, nil, err)

	vmConfig := &vm.Config{
		JumpTable: vm.ConstantinopleInstructionSet,
	}
	parent := bcdata.bc.CurrentBlock()
	num := parent.Number()
	author := bcdata.addrs[0]
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(parent.Time(), common.Big1),
		BlockScore: big.NewInt(0),
	}
	usedGas := uint64(0)
	receipt, gas, _, err := bcdata.bc.ApplyTransaction(bcdata.bc.Config(), author, state, header, tx, &usedGas, vmConfig)
	return receipt, gas, err
}
