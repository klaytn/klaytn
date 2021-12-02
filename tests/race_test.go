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

// +build race

package tests

import (
	"os"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// TestRaceBetweenTxpoolAddAndCommitNewWork tests race conditions between `Txpool.add` and `commitNewWork`.
// Since both access to txpool pending concurrently, critical sections should be protected by mutex lock.
// This race test may need multiple trials and additional flags to avoid false alarms from sha3 package.
// For example, `go test -gcflags=all=-d=checkptr=0 -race -run TestRaceBetweenTxpoolAddAndCommitNewWork`.
func TestRaceBetweenTxpoolAddAndCommitNewWork(t *testing.T) {
	if testing.Verbose() {
		enableLog() // Change verbosity level in the function if needed
	}

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
				signer := types.NewEIP155Signer(chainId)
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
				signer := types.NewEIP155Signer(chainId)
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
