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
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/compiler"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
)

type deployedContract struct {
	abi     string
	name    string
	address common.Address
}

func deployContract(filename string, bcdata *BCData, accountMap *AccountMap,
	prof *profile.Profiler,
) (map[string]*deployedContract, error) {
	contracts, err := compiler.CompileSolidityOrLoad("", filename)
	if err != nil {
		return nil, err
	}

	cont := make(map[string]*deployedContract)
	transactions := make(types.Transactions, 0, 10)

	userAddr := bcdata.addrs[0]
	nonce := accountMap.GetNonce(*userAddr)

	// create a contract tx
	for name, contract := range contracts {

		abiStr, err := json.Marshal(contract.Info.AbiDefinition)
		if err != nil {
			return nil, err
		}

		header := bcdata.bc.CurrentHeader()

		contractAddr := crypto.CreateAddress(*userAddr, nonce)

		signer := types.MakeSigner(bcdata.bc.Config(), header.Number)
		tx := types.NewContractCreation(nonce,
			big.NewInt(0), 50000000, big.NewInt(0), common.FromHex(contract.Code))
		signedTx, err := types.SignTx(tx, signer, bcdata.privKeys[0])
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, signedTx)

		cont[name] = &deployedContract{
			abi:     string(abiStr),
			name:    name,
			address: contractAddr,
		}

		nonce += 1
	}

	bcdata.GenABlockWithTransactions(accountMap, transactions, prof)

	return cont, nil
}

func callContract(bcdata *BCData, tx *types.Transaction) ([]byte, error) {
	header := bcdata.bc.CurrentHeader()
	statedb, err := bcdata.bc.State()
	if err != nil {
		return nil, err
	}

	signer := types.MakeSigner(bcdata.bc.Config(), header.Number)
	msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, bcdata.bc.CurrentBlock().NumberU64())
	if err != nil {
		return nil, err
	}

	evmContext := blockchain.NewEVMContext(msg, header, bcdata.bc, nil)
	vmenv := vm.NewEVM(evmContext, statedb, bcdata.bc.Config(), &vm.Config{})

	ret, _, kerr := blockchain.NewStateTransition(vmenv, msg).TransitionDb()
	err = kerr.ErrTxInvalid
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func makeRewardTransactions(c *deployedContract, accountMap *AccountMap, bcdata *BCData,
	numTransactions int,
) (types.Transactions, error) {
	abii, err := abi.JSON(strings.NewReader(c.abi))
	if err != nil {
		return nil, err
	}

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	transactions := make(types.Transactions, numTransactions)

	numAddrs := len(bcdata.addrs)
	fromNonces := make([]uint64, numAddrs)
	for i, addr := range bcdata.addrs {
		fromNonces[i] = accountMap.GetNonce(*addr)
	}
	for i := 0; i < numTransactions; i++ {
		idx := i % numAddrs

		addr := bcdata.addrs[idx]
		data, err := abii.Pack("reward", addr)
		if err != nil {
			return nil, err
		}

		tx := types.NewTransaction(fromNonces[idx], c.address, big.NewInt(10), 5000000, big.NewInt(0), data)
		signedTx, err := types.SignTx(tx, signer, bcdata.privKeys[idx])
		if err != nil {
			return nil, err
		}

		transactions[i] = signedTx
		fromNonces[idx]++
	}

	return transactions, nil
}

func executeRewardTransactions(c *deployedContract, transactions types.Transactions, prof *profile.Profiler, bcdata *BCData,
	accountMap *AccountMap,
) error {
	return bcdata.GenABlockWithTransactions(accountMap, transactions, prof)
}

func makeBalanceOf(c *deployedContract, accountMap *AccountMap, bcdata *BCData,
	numTransactions int,
) (types.Transactions, error) {
	abii, err := abi.JSON(strings.NewReader(c.abi))
	if err != nil {
		return nil, err
	}

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	transactions := make(types.Transactions, numTransactions)

	numAddrs := len(bcdata.addrs)
	fromNonces := make([]uint64, numAddrs)
	for i, addr := range bcdata.addrs {
		fromNonces[i] = accountMap.GetNonce(*addr)
	}
	for i := 0; i < numTransactions; i++ {
		idx := i % numAddrs

		addr := bcdata.addrs[idx]
		data, err := abii.Pack("balanceOf", addr)
		if err != nil {
			return nil, err
		}

		tx := types.NewTransaction(fromNonces[idx], c.address, big.NewInt(0), 5000000, big.NewInt(0), data)
		signedTx, err := types.SignTx(tx, signer, bcdata.privKeys[idx])
		if err != nil {
			return nil, err
		}

		transactions[i] = signedTx

		// This is not required because the transactions will not be inserted into the blockchain.
		// fromNonces[idx]++
	}

	return transactions, nil
}

func executeBalanceOf(c *deployedContract, transactions types.Transactions, prof *profile.Profiler, bcdata *BCData,
	accountMap *AccountMap,
) error {
	abii, err := abi.JSON(strings.NewReader(c.abi))
	if err != nil {
		return err
	}

	for _, tx := range transactions {
		ret, err := callContract(bcdata, tx)
		if err != nil {
			return err
		}

		balance := new(big.Int)
		abii.Unpack(&balance, "balanceOf", ret)
	}

	return nil
}

func makeQuickSortTransactions(c *deployedContract, accountMap *AccountMap, bcdata *BCData,
	numTransactions int,
) (types.Transactions, error) {
	abii, err := abi.JSON(strings.NewReader(c.abi))
	if err != nil {
		return nil, err
	}

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	transactions := make(types.Transactions, numTransactions)

	numAddrs := len(bcdata.addrs)
	fromNonces := make([]uint64, numAddrs)
	for i, addr := range bcdata.addrs {
		fromNonces[i] = accountMap.GetNonce(*addr)
	}
	for i := 0; i < numTransactions; i++ {
		idx := i % numAddrs

		data, err := abii.Pack("sort", big.NewInt(100), big.NewInt(123))
		if err != nil {
			return nil, err
		}

		tx := types.NewTransaction(fromNonces[idx], c.address, nil, 10000000, big.NewInt(0), data)
		signedTx, err := types.SignTx(tx, signer, bcdata.privKeys[idx])
		if err != nil {
			return nil, err
		}

		transactions[i] = signedTx
		fromNonces[idx]++
	}

	return transactions, nil
}

func executeQuickSortTransactions(c *deployedContract, transactions types.Transactions, prof *profile.Profiler, bcdata *BCData,
	accountMap *AccountMap,
) error {
	return bcdata.GenABlockWithTransactions(accountMap, transactions, prof)
}

func executeSmartContract(b *testing.B, opt *ContractExecutionOption, prof *profile.Profiler) {
	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(2000, 4)
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

	start = time.Now()
	contracts, err := deployContract(opt.filepath, bcdata, accountMap, prof)
	if err != nil {
		b.Fatal(err)
	}
	prof.Profile("main_deployContract", time.Now().Sub(start))

	b.StopTimer()
	b.ResetTimer()
	for _, c := range contracts {
		start = time.Now()
		transactions, err := opt.makeTx(c, accountMap, bcdata, b.N)
		if err != nil {
			b.Fatal(err)
		}
		prof.Profile("main_makeTx", time.Now().Sub(start))

		start = time.Now()
		b.StartTimer()
		opt.executeTx(c, transactions, prof, bcdata, accountMap)
		b.StopTimer()
		prof.Profile("main_executeTx", time.Now().Sub(start))
	}
}

type ContractExecutionOption struct {
	name      string
	filepath  string
	makeTx    func(c *deployedContract, accountMap *AccountMap, bcdata *BCData, numTransactions int) (types.Transactions, error)
	executeTx func(c *deployedContract, transactions types.Transactions, prof *profile.Profiler, bcdata *BCData, accountMap *AccountMap) error
}

func BenchmarkSmartContractExecute(b *testing.B) {
	prof := profile.NewProfiler()

	benches := []ContractExecutionOption{
		{"KlaytnReward:reward", "../contracts/reward/contract/KlaytnReward.sol", makeRewardTransactions, executeRewardTransactions},
		{"KlaytnReward:balanceOf", "../contracts/reward/contract/KlaytnReward.sol", makeBalanceOf, executeBalanceOf},
		{"QuickSort:sort", "./testdata/contracts/sort/QuickSort.sol", makeQuickSortTransactions, executeQuickSortTransactions},
	}

	for _, bench := range benches {
		b.Run(bench.name, func(b *testing.B) {
			executeSmartContract(b, &bench, prof)
		})
	}

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

func BenchmarkStorageTrieStore(b *testing.B) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	benchOption := ContractExecutionOption{
		"StorageTrieStore",
		"../contracts/storagetrie/StorageTrieStoreTest.sol",
		makeStorageTrieTransactions,
		nil,
	}

	executeSmartContractForStorageTrie(b, &benchOption, prof)

	if testing.Verbose() {
		prof.PrintProfileInfo()
	}
}

func executeSmartContractForStorageTrie(b *testing.B, opt *ContractExecutionOption, prof *profile.Profiler) {
	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCData(2000, 4)
	if err != nil {
		b.Fatal(err)
	}
	prof.Profile("main_init_blockchain", time.Now().Sub(start))

	defer bcdata.db.Close()
	defer bcdata.bc.Stop()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		b.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	start = time.Now()
	contracts, err := deployContract(opt.filepath, bcdata, accountMap, prof)
	if err != nil {
		b.Fatal(err)
	}

	if len(contracts) != 1 {
		b.Fatalf("contracts length should be 1 but %v!!", len(contracts))
	}

	prof.Profile("main_deployContract", time.Now().Sub(start))

	b.StopTimer()
	b.ResetTimer()

	timeNow := time.Now()
	f, err := os.Create(opt.name + "_" + timeNow.Format("2006-01-02-1504") + ".cpu.out")
	if err != nil {
		b.Fatalf("failed to create file for cpu profiling, err: %v", err)
	}

	totalRuns := 1

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)
	txPool := makeTxPool(bcdata, txPoolSize)
	for _, c := range contracts {
		start = time.Now()
		for run := 1; run <= totalRuns; run++ {
			fmt.Printf("run %v started \n", run)
			transactions, err := opt.makeTx(c, accountMap, bcdata, 20000)
			if err != nil {
				b.Fatal(err)
			}
			fmt.Printf("run %v tx generated \n", run)

			state, _ := bcdata.bc.State()
			for _, tx := range transactions {
				if _, err := tx.AsMessageWithAccountKeyPicker(signer, state, bcdata.bc.CurrentBlock().NumberU64()); err != nil {
					b.Fatal(err)
				}
			}
			fmt.Printf("run %v tx validated \n", run)

			start = time.Now()
			b.StartTimer()

			if run == totalRuns {
				pprof.StartCPUProfile(f)
			}

			if err := executeTxs(bcdata, txPool, transactions); err != nil {
				b.Fatal(err)
			}

			b.StopTimer()
			if run == totalRuns {
				pprof.StopCPUProfile()
			}
		}
	}
}

func makeStorageTrieTransactions(c *deployedContract, accountMap *AccountMap, bcdata *BCData,
	numTransactions int,
) (types.Transactions, error) {
	abii, err := abi.JSON(strings.NewReader(c.abi))
	if err != nil {
		return nil, err
	}

	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	transactions := make(types.Transactions, numTransactions)

	stateDB, _ := bcdata.bc.State()

	numAddrs := len(bcdata.addrs)
	fromNonces := make([]uint64, numAddrs)
	for i, addr := range bcdata.addrs {
		fromNonces[i] = stateDB.GetNonce(*addr)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numTransactions; i++ {
		idx := i % numAddrs

		// function insertIdentity(string _serialNumber, string _publicKey, string _hash)
		data, err := abii.Pack("insertIdentity", randomString(39, r), randomString(814, r), randomString(40, r))
		if err != nil {
			return nil, err
		}

		tx := types.NewTransaction(fromNonces[idx], c.address, nil, 10000000, big.NewInt(25000000000), data)
		signedTx, err := types.SignTx(tx, signer, bcdata.privKeys[idx])
		if err != nil {
			return nil, err
		}

		transactions[i] = signedTx
		fromNonces[idx]++
	}

	return transactions, nil
}

func randomBytes(n int, rand *rand.Rand) []byte {
	r := make([]byte, n)
	if _, err := rand.Read(r); err != nil {
		panic("rand.Read failed: " + err.Error())
	}
	return r
}

func randomString(n int, rand *rand.Rand) string {
	b := randomBytes(n, rand)
	return string(b)
}
