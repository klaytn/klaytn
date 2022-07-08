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
	"math/big"
	"testing"
	"time"

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

// TestValidatingUnavailableContractExecution tests validation logic of invalid contract execution transaction.
// TxPool will invalidate contract execution transactions sending to un-executable account even though the recipient is a contract account.
func TestValidatingUnavailableContractExecution(t *testing.T) {
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

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// various accounts preparation
	contract, err := createAnonymousAccount("98275a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78991dab")
	contract.Addr = common.Address{}

	invalidContract, err := createAnonymousAccount("11115a145bc1726eb0445433088f5f882f8a4a9499135239cfb4040e78911111")
	invalidContract.Addr = common.Address{}

	EOA, err := createDefaultAccount(accountkey.AccountKeyTypeLegacy)
	assert.Equal(t, nil, err)

	code = "0x608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"

	// Prepare. creates various types of accounts tried to be executed
	{
		var txs types.Transactions

		// tx to create a contract account
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.GetNonce(),
			types.TxValueKeyAmount:        new(big.Int).SetUint64(0),
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyFrom:          reservoir.GetAddr(),
			types.TxValueKeyData:          common.FromHex(code),
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		contract.Addr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1

		// tx2 to create a invalid contract account
		values = map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         reservoir.GetNonce(),
			types.TxValueKeyAmount:        new(big.Int).SetUint64(0),
			types.TxValueKeyGasLimit:      gasLimit,
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyFrom:          reservoir.GetAddr(),
			types.TxValueKeyData:          []byte{}, // the invalid contract doesn't have contract code
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx2, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		err = tx2.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx2)

		invalidContract.Addr = crypto.CreateAddress(reservoir.Addr, reservoir.Nonce)

		reservoir.Nonce += 1

		// tx3 to create an EOA account
		tx3, _, err := generateDefaultTx(reservoir, EOA, types.TxTypeValueTransfer, contract.Addr)
		assert.Equal(t, nil, err)

		err = tx3.SignWithKeys(signer, reservoir.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx3)
		reservoir.Nonce += 1

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
	}

	// make TxPool to test validation in 'TxPool add' process
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	// 1. contract execution transaction to the contract account.
	{
		tx, _ := genSmartContractExecution(t, signer, reservoir, contract, nil, gasPrice)

		err = txpool.AddRemote(tx)
		assert.Equal(t, nil, err)
	}

	// 2. contract execution transaction to the invalid contract account.
	{
		tx, _ := genSmartContractExecution(t, signer, reservoir, invalidContract, nil, gasPrice)

		err = txpool.AddRemote(tx)
		assert.Equal(t, kerrors.ErrNotProgramAccount, err)
	}

	// 3. contract execution transaction to the EOA account.
	{
		tx, _ := genSmartContractExecution(t, signer, reservoir, EOA, nil, gasPrice)

		err = txpool.AddRemote(tx)
		assert.Equal(t, kerrors.ErrNotProgramAccount, err)
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}
