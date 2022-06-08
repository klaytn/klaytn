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
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
)

// BenchmarkResendNilDereference checks nil pointer dereferences while
// simultaneous execution of TxPool.AddRemotes() and TxPool.CachedPendingTxsByCount().
// Since TxPool.CachedPendingTxsByCount() updates m.cache in the TxPool.TxList without
// obtaining TxPool.mu, m.cache can be updated in both execution of TxPool.AddRemotes()
// and TxPool.CachedPendingTxsByCount().
// By obtaining TxPool.mu in TxPool.CachedPendingTxsByCount(),
// the nil pointer dereference errors disappeared.
//
// This does not need to be executed on CI, make this test a benchmark.
// To execute this,
// $ go test -run XXX -bench BenchmarkResendNilDereference
func BenchmarkResendNilDereference(t *testing.B) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	numTransactions := 20000
	numAccounts := 2000

	opt := testOption{numTransactions, numAccounts, 4, 1, []byte{}, makeNewTransactionsToRing}

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(opt.numMaxAccounts, opt.numValidators)
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

	// make txpool
	txpool := makeTxPool(bcdata, 50000)
	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	wg := sync.WaitGroup{}
	wg.Add(3)

	gasPrice := new(big.Int).SetUint64(25 * params.Ston)
	txAmount := big.NewInt(3)
	// reservoir account
	reservoir := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	go func(num int) {
		// fmt.Println("creating transaction start!")
		time.Sleep(1 * time.Nanosecond)
		// make ring transactions
		for i := 0; i < num; i++ {
			// fmt.Println("tx gen num", i)

			tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:    reservoir.Nonce,
				types.TxValueKeyTo:       reservoir.Addr,
				types.TxValueKeyAmount:   txAmount,
				types.TxValueKeyGasLimit: gasLimit,
				types.TxValueKeyGasPrice: gasPrice,
				types.TxValueKeyFrom:     reservoir.Addr,
			})
			if err != nil {
				fmt.Println("tx gen err", err)
				t.Error(err)
				return
			}

			signedTx, err := types.SignTx(tx, signer, bcdata.privKeys[0])
			if err != nil {
				fmt.Println("sign err", err)
			}
			time.Sleep(1 * time.Nanosecond)
			txpool.AddRemotes(types.Transactions{signedTx})

			reservoir.Nonce++
		}
		// fmt.Println("creating transaction done!")
		wg.Done()
	}(100000)

	go func(iter int) {
		// fmt.Println("Thread1 start!")
		for i := 0; i < iter; i++ {
			// fmt.Println("CachedPendingTxsByCount num", i)
			time.Sleep(1 * time.Nanosecond)
			txpool.CachedPendingTxsByCount(20000)
		}
		wg.Done()
		// fmt.Println("Thread1 done!")
	}(100000)

	go func(iter int) {
		// fmt.Println("Thread2 start!")
		for i := 0; i < iter; i++ {
			// fmt.Println("Content num", i)
			time.Sleep(1 * time.Nanosecond)
			txpool.Content()
		}
		wg.Done()
		// fmt.Println("Thread2 done!")
	}(100000)

	wg.Wait()
	fmt.Println(txpool.Stats())

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}
