// Copyright 2018 The klaytn Authors
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
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/params"
)

var cpuprofile bool = false

func init() {
	flag.BoolVar(&cpuprofile, "addtx_profile", false, "Enable cpu profiling for AddTx")
}

////////////////////////////////////////////////////////////////////////////////
// BenchmarkAddTx
////////////////////////////////////////////////////////////////////////////////
func BenchmarkAddTx(b *testing.B) {
	cacheSender := []bool{false, true}
	maxAccounts := []int{1000 + 100, 1000 + 1000, 1000 + 10000}
	numValidators := []int{4}
	parallels := []string{"parallel", "sequential", "queueing"}
	numQueues := []int{1, 2, 4, 8}

	fmt.Printf("addtx_profile = %t\n", cpuprofile)

	for _, cs := range cacheSender {
		for _, ma := range maxAccounts {
			for _, nv := range numValidators {
				for _, p := range parallels {

					numQueues := numQueues
					if p != "queueing" {
						numQueues = []int{1}
					}
					for _, nq := range numQueues {
						testName := fmt.Sprintf("%t,%d,%d,%s,%d", cs, ma-1000, nv, p, nq)
						b.Run(testName, func(b *testing.B) {
							benchAddTx(b, ma, nv, p, nq, cs)
						})
					}
				}
			}
		}
	}
}

func txDispatcher(ch <-chan *types.Transaction, txpool *blockchain.TxPool, wait *sync.WaitGroup) {
	for {
		if t, ok := <-ch; ok {
			err := txpool.AddLocal(t)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			wait.Done()
		} else {
			break
		}
	}
}

func benchAddTx(b *testing.B, maxAccounts, numValidators int, parallel string, numQueue int,
	cacheSender bool,
) {
	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(maxAccounts, numValidators)
	if err != nil {
		b.Fatal(err)
	}
	profile.Prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		b.Fatal(err)
	}
	profile.Prof.Profile("main_init_accountMap", time.Now().Sub(start))

	poolConfig := blockchain.TxPoolConfig{
		Journal:         transactionsJournalFilename,
		JournalInterval: time.Hour,

		PriceLimit: 1,
		PriceBump:  10,

		ExecSlotsAccount:    16,
		ExecSlotsAll:        40000,
		NonExecSlotsAccount: 64,
		NonExecSlotsAll:     40000,

		Lifetime: 5 * time.Minute,
	}
	txpool := blockchain.NewTxPool(poolConfig, bcdata.bc.Config(), bcdata.bc)

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentBlock().Number())

	var wait sync.WaitGroup

	var txChs []chan *types.Transaction

	if parallel == "queueing" {
		txChs = make([]chan *types.Transaction, numQueue)
		for i := 0; i < numQueue; i++ {
			txChs[i] = make(chan *types.Transaction, 100)
		}

		for i := 0; i < numQueue; i++ {
			go txDispatcher(txChs[i], txpool, &wait)
		}
	}

	var f *os.File = nil
	if cpuprofile {
		f, err = os.Create(fmt.Sprintf("profile_cpu_%t_%d_%d_%s_%d.out",
			cacheSender, maxAccounts-1000, numValidators, parallel, numQueue))
		if err != nil {
			b.Fatal("could not create CPU profile :", err)
		}
	}

	txs := make([]types.Transactions, b.N)
	for i := 0; i < b.N; i++ {
		txs[i], err = makeTransactions(accountMap,
			bcdata.addrs[1000:maxAccounts], bcdata.privKeys[1000:maxAccounts],
			signer, bcdata.addrs[0:maxAccounts-1000], nil, i, cacheSender)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	if cpuprofile {
		if err := pprof.StartCPUProfile(f); err != nil {
			b.Fatal("could not start CPU profile : ", err)
		}
	}
	for i := 0; i < b.N; i++ {
		switch parallel {
		case "parallel":
			txParallel(txs[i], txpool)

		case "sequential":
			txSequential(txs[i], txpool)

		case "queueing":
			txQueue(txs[i], txpool, txChs, numQueue, &wait)
		}

		pending, _ := txpool.Stats()
		if pending%(maxAccounts-1000) != 0 {
			b.Fatalf("pending(%d) should be divided by %d!\n", pending, maxAccounts-1000)
		}
	}
	b.StopTimer()
	if cpuprofile {
		pprof.StopCPUProfile()
	}

	txpool.Stop()
	if cpuprofile {
		f.Close()
	}

	if testing.Verbose() {
		profile.Prof.PrintProfileInfo()
	}
}

func makeTransactions(accountMap *AccountMap, fromAddrs []*common.Address, privKeys []*ecdsa.PrivateKey,
	signer types.Signer, toAddrs []*common.Address, amount *big.Int, additionalNonce int,
	cacheSender bool,
) (types.Transactions, error) {
	txs := make(types.Transactions, 0, len(toAddrs))
	for i, from := range fromAddrs {
		nonce := accountMap.GetNonce(*from)
		nonce += uint64(additionalNonce)

		txamount := amount
		if txamount == nil {
			txamount = big.NewInt(rand.Int63n(10))
			txamount = txamount.Add(txamount, big.NewInt(1))
		}

		var gasLimit uint64 = 1000000
		gasPrice := new(big.Int).SetInt64(25 * params.Ston)
		data := []byte{}

		tx := types.NewTransaction(nonce, *toAddrs[i], txamount, gasLimit, gasPrice, data)
		signedTx, err := types.SignTx(tx, signer, privKeys[i])
		if err != nil {
			return nil, err
		}

		if cacheSender {
			signed_addr, err := types.Sender(signer, signedTx)
			if signed_addr != *from {
				fmt.Printf("signed address(%s) != from(%s)\n", signed_addr.Hex(), from.Hex())
			}
			if err != nil {
				return nil, err
			}
		}

		txs = append(txs, signedTx)
	}

	return txs, nil
}

func txParallel(txs types.Transactions, txpool *blockchain.TxPool) {
	var wait sync.WaitGroup

	wait.Add(len(txs))

	for _, tx := range txs {
		go func(t *types.Transaction) {
			err := txpool.AddLocal(t)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			wait.Done()
		}(tx)
	}

	wait.Wait()
}

func txSequential(txs types.Transactions, txpool *blockchain.TxPool) {
	for _, tx := range txs {
		err := txpool.AddLocal(tx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func txQueue(txs types.Transactions, txpool *blockchain.TxPool,
	txChs []chan *types.Transaction, numQueue int, wait *sync.WaitGroup,
) {
	wait.Add(len(txs))

	for i, tx := range txs {
		txChs[i%numQueue] <- tx
	}

	wait.Wait()
}
