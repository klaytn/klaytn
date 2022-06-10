// Copyright 2021 The klaytn Authors
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

//go:build race
// +build race

package tests

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

// TestRaceBetweenTxpoolAddAndCommitNewWork tests race conditions between `Txpool.add` and `commitNewWork`.
// Since both access to txpool pending concurrently, critical sections should be protected by mutex lock.
// This race test may need multiple trials and additional flags to avoid false alarms from sha3 package.
// For example, `go test -gcflags=all=-d=checkptr=0 -race -run TestRaceBetweenTxpoolAddAndCommitNewWork`.
func TestRaceBetweenTxpoolAddAndCommitNewWork(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	numAccounts := 2
	fullNode, node, validator, chainId, workspace := newBlockchain(t)
	defer os.RemoveAll(workspace)

	// create account
	richAccount, accounts, _ := createAccount(t, numAccounts, validator)

	quitCh := make(chan struct{})
	iterNum := 1000

	go func() {
		var txList []*types.Transaction
		for i := 0; i < iterNum; i++ {
			{
				tx, _, err := generateDefaultTx(richAccount, accounts[1], types.TxTypeValueTransfer, common.Address{})
				if err != nil {
					t.Fatal(err)
				}
				signer := types.LatestSignerForChainID(chainId)
				if err := tx.Sign(signer, richAccount.Keys[0]); err != nil {
					t.Fatal(err)
				}
				txList = append(txList, tx)
			}
			{
				tx, _, err := generateDefaultTx(richAccount, accounts[1], types.TxTypeCancel, common.Address{})
				if err != nil {
					t.Fatal(err)
				}
				signer := types.LatestSignerForChainID(chainId)
				if err := tx.Sign(signer, richAccount.Keys[0]); err != nil {
					t.Fatal(err)
				}
				txList = append(txList, tx)
			}
			richAccount.AddNonce()
		}

		for _, tx := range txList {
			if err := node.TxPool().AddLocal(tx); err != nil {
				t.Fatal(err)
			}
		}
		quitCh <- struct{}{}
	}()

	<-quitCh
	time.Sleep(time.Second)

	// stop node before ending the test code
	if err := fullNode.Stop(); err != nil {
		t.Fatal(err)
	}
}

// TestRaceAsMessageWithAccountPickerForFeePayer tests calling AsMessageWithAccountPicker of a fee delegated transaction
// where a fee payer may be inserted wrongly due to concurrent issue.
func TestRaceAsMessageWithAccountPickerForFeePayer(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	// Configure and generate a sample block chain
	var (
		gendb = database.NewMemoryDBManager()

		// create a sender and a feepayer
		from, _     = createAnonymousAccount("a5c9a50938a089618167c9d67dbebc0deaffc3c76ddc6b40c2777ae594389999")
		feePayer, _ = createAnonymousAccount("ed580f5bd71a2ee4dae5cb43e331b7d0318596e561e6add7844271ed94156b20")

		funds = new(big.Int).Mul(big.NewInt(1e16), big.NewInt(params.KLAY))
		gspec = &blockchain.Genesis{
			Config: params.TestChainConfig,
			Alloc: blockchain.GenesisAlloc{
				from.GetAddr():     {Balance: funds},
				feePayer.GetAddr(): {Balance: funds},
			},
		}
		genesis = gspec.MustCommit(gendb)
		signer  = types.LatestSignerForChainID(gspec.Config.ChainID)
	)

	iterNum := 10000
	errCh := make(chan error, 2*iterNum)

	for i := 0; i < iterNum; i++ {
		tx, _ := genFeeDelegatedChainDataAnchoring(t, signer, from, nil, feePayer, big.NewInt(1234))
		for i := 0; i < 2; i++ {
			go func() {
				stateDB, err := state.New(genesis.Root(), state.NewDatabase(gendb))
				if err != nil {
					panic(err)
				}

				msg, err := tx.AsMessageWithAccountKeyPicker(signer, stateDB, 0)
				if err != nil {
					panic(err)
				}

				if msg.ValidatedFeePayer() != feePayer.GetAddr() {
					errCh <- fmt.Errorf("expected: %v, actual: %v", feePayer.GetAddr().String(), msg.ValidatedFeePayer().String())
				} else {
					errCh <- nil
				}
			}()
		}
	}

	for i := 0; i < 2*iterNum; i++ {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}
}
