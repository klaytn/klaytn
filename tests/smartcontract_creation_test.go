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
	"math/big"
	"path/filepath"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/compiler"
	"github.com/klaytn/klaytn/common/profile"
)

type testData struct {
	name string
	opt  testOption
}

func makeContractCreationTransactions(bcdata *BCData, accountMap *AccountMap, signer types.Signer,
	numTransactions int, amount *big.Int, data []byte,
) (types.Transactions, error) {
	numAddrs := len(bcdata.addrs)
	fromAddrs := bcdata.addrs
	fromNonces := make([]uint64, numAddrs)

	for i, addr := range fromAddrs {
		fromNonces[i] = accountMap.GetNonce(*addr)
	}

	txs := make(types.Transactions, 0, numTransactions)

	for i := 0; i < numTransactions; i++ {
		idx := i % numAddrs

		txamount := new(big.Int).SetInt64(0)

		var gasLimit uint64 = 1000000
		gasPrice := new(big.Int).SetInt64(0)

		tx := types.NewContractCreation(fromNonces[idx], txamount, gasLimit, gasPrice, data)
		signedTx, err := types.SignTx(tx, signer, bcdata.privKeys[idx])
		if err != nil {
			return nil, err
		}

		txs = append(txs, signedTx)

		fromNonces[idx]++
	}

	return txs, nil
}

func genOptions(b *testing.B) ([]testData, error) {
	solFiles := []string{"../contracts/reward/contract/KlaytnReward.sol"}

	opts := make([]testData, len(solFiles))
	for i, filename := range solFiles {
		contracts, err := compiler.CompileSolidity("", filename)
		if err != nil {
			return nil, err
		}

		for name, contract := range contracts {
			testName := filepath.Base(name)
			opts[i] = testData{testName, testOption{
				b.N, 2000, 4, 1, common.FromHex(contract.Code), makeContractCreationTransactions,
			}}
		}
	}

	return opts, nil
}

func deploySmartContract(b *testing.B, opt *testOption, prof *profile.Profiler) {
	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(opt.numMaxAccounts, opt.numValidators)
	if err != nil {
		b.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))
	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		b.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	b.ResetTimer()
	for i := 0; i < b.N/txPerBlock; i++ {
		// fmt.Printf("iteration %d tx %d\n", i, opt.numTransactions)
		err := bcdata.GenABlock(accountMap, opt, txPerBlock, prof)
		if err != nil {
			b.Fatal(err)
		}
	}

	genBlocks := b.N / txPerBlock
	remainTxs := b.N % txPerBlock
	if remainTxs != 0 {
		err := bcdata.GenABlock(accountMap, opt, remainTxs, prof)
		if err != nil {
			b.Fatal(err)
		}
		genBlocks++
	}

	bcHeight := int(bcdata.bc.CurrentHeader().Number.Uint64())
	if bcHeight != genBlocks {
		b.Fatalf("generated blocks should be %d, but %d.\n", genBlocks, bcHeight)
	}
}

func BenchmarkSmartContractDeploy(b *testing.B) {
	prof := profile.NewProfiler()

	benches, err := genOptions(b)
	if err != nil {
		b.Fatal(err)
	}

	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			bench.opt.numTransactions = b.N
			deploySmartContract(b, &bench.opt, prof)
		})
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}
