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
	"fmt"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/storage/statedb"
)

func BenchmarkDeriveSha(b *testing.B) {
	funcs := map[string]types.IDeriveSha{
		"Orig":   statedb.DeriveShaOrig{},
		"Simple": types.DeriveShaSimple{},
		"Concat": types.DeriveShaConcat{},
	}

	NTS := []int{1000}

	for k, f := range funcs {
		for _, nt := range NTS {
			testName := fmt.Sprintf("%s,%d", k, nt)
			b.Run(testName, func(b *testing.B) {
				benchDeriveSha(b, nt, 4, f)
			})
		}
	}
}

func benchDeriveSha(b *testing.B, numTransactions, numValidators int, sha types.IDeriveSha) {
	// Initialize blockchain
	start := time.Now()
	maxAccounts := numTransactions * 2
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

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentBlock().Number())

	txs, err := makeTransactions(accountMap,
		bcdata.addrs[0:numTransactions], bcdata.privKeys[0:numTransactions],
		signer, bcdata.addrs[numTransactions:numTransactions*2], nil, 0, false)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hash := sha.DeriveSha(txs)
		if testing.Verbose() {
			fmt.Printf("[%d] txhash = %s\n", i, hash.Hex())
		}
	}
}
