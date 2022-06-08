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
	"runtime/pprof"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

const txPoolSize = 32768

// BenchmarkDataExecution_Aspen generates the data with Aspen network's database configurations.
func BenchmarkDataExecution_Aspen(b *testing.B) {
	tc := getExecutionTestDefaultTC()
	tc.testName = "BenchmarkDataExecution_Aspen"
	tc.originalDataDir = aspen500_orig
	tc.dbc, tc.levelDBOption = genAspenOptions()

	dataExecutionTest(b, tc)
}

// BenchmarkDataExecution_Baobab generates the data with Baobab network's database configurations.
func BenchmarkDataExecution_Baobab(b *testing.B) {
	tc := getExecutionTestDefaultTC()
	tc.testName = "BenchmarkDataExecution_Baobab"
	tc.originalDataDir = baobab500_orig

	tc.dbc, tc.levelDBOption = genBaobabOptions()

	dataExecutionTest(b, tc)
}

// BenchmarkDataExecution_CandidateLevelDB generates the data for main-net's
// with candidate configurations, using LevelDB.
func BenchmarkDataExecution_CandidateLevelDB(b *testing.B) {
	tc := getExecutionTestDefaultTC()
	tc.testName = "BenchmarkDataExecution_CandidateLevelDB"
	tc.originalDataDir = candidate500LevelDB_orig

	tc.dbc, tc.levelDBOption = genCandidateLevelDBOptions()

	dataExecutionTest(b, tc)
}

// BenchmarkDataExecution_CandidateBadgerDB generates the data for main-net's
// with candidate configurations, using BadgerDB.
func BenchmarkDataExecution_CandidateBadgerDB(b *testing.B) {
	tc := getExecutionTestDefaultTC()
	tc.testName = "BenchmarkDataExecution_CandidateBadgerDB"
	tc.originalDataDir = candidate500BadgerDB_orig

	tc.dbc, tc.levelDBOption = genCandidateBadgerDBOptions()

	dataExecutionTest(b, tc)
}

// BenchmarkDataExecution_Baobab_ControlGroup generates the data with Baobab network's database configurations.
// To work as a control group, it only generates 10,000 accounts.
func BenchmarkDataExecution_Baobab_ControlGroup(b *testing.B) {
	tc := getExecutionTestDefaultTC()
	tc.testName = "BenchmarkDataExecution_Baobab_ControlGroup"
	tc.originalDataDir = baobab1_orig

	tc.dbc, tc.levelDBOption = genBaobabOptions()

	// ControlGroup specific setting
	tc.numReceiversPerRun = 10000

	dataExecutionTest(b, tc)
}

// Static variables, not to read same addresses and keys from files repeatedly.
// If there are saved addresses and keys and the given numReceiversPerRun and testDataDir
// match with saved ones, it will reuse saved addresses and keys.
var savedAddresses []*common.Address = nil

var (
	savedKeys               []*ecdsa.PrivateKey = nil
	savedNumReceiversPerRun int
	savedTestDataDir        string
)

// dataExecutionTest is to check the performance of Klaytn with pre-generated data.
// It generates warmUpTxs and executionTxs first, and then initialize blockchain and database to
// remove any effects caused by generating transactions. And then it executes warmUpTxs and executionTxs.
// To run the test, original data directory should be located at "$GOPATH/src/github.com/klaytn/"
func dataExecutionTest(b *testing.B, tc *preGeneratedTC) {
	testDataDir, profileFile, err := setUpTest(tc)
	if err != nil {
		b.Fatal(err)
	}

	///////////////////////////////////////////////////////////////////////////////////
	///  Tx Generation Process. Generate warmUpTxs and executionTxs in this phase.  ///
	///////////////////////////////////////////////////////////////////////////////////

	// activeAddrs is used for both sender and receiver.
	// len(activeAddrs) = tc.numReceiversPerRun
	if savedAddresses == nil || savedKeys == nil || tc.numReceiversPerRun != savedNumReceiversPerRun || testDataDir != savedTestDataDir {
		fmt.Println("Start reading addresses from files", "testDataDir", testDataDir, "numReceiversPerRun", tc.numReceiversPerRun)
		savedAddresses, savedKeys, err = getAddrsAndKeysFromFile(tc.numReceiversPerRun, testDataDir, 0, tc.filePicker)
		savedNumReceiversPerRun = tc.numReceiversPerRun
		savedTestDataDir = testDataDir
		if err != nil {
			b.Fatal(err)
		}
		fmt.Println("End reading addresses from files")
	} else {
		fmt.Println("Reuse previously saved addresses and keys", "len(addrs)", len(savedAddresses), "len(keys)", len(savedKeys))
	}

	activeAddrs, activeKeys := savedAddresses, savedKeys

	bcData, err := NewBCDataForPreGeneratedTest(testDataDir, tc)
	if err != nil {
		b.Fatal(err)
	}

	// Generate two different list of transactions, warmUpTxs and executionTxs
	// warmUpTxs is to activate caching.
	// executionTxs is to measure the performance of test after caching.
	signer := types.MakeSigner(bcData.bc.Config(), bcData.bc.CurrentHeader().Number)
	stateDB, err := bcData.bc.State()
	if err != nil {
		b.Fatal(err)
	}

	// len(warmUpTxs) = tc.numReceiversPerRun
	warmUpTxs, nonceMap, err := makeTxsWithStateDB(false, stateDB, activeAddrs, activeKeys, activeAddrs, signer, tc.numReceiversPerRun, sequentialIndex)
	if err != nil {
		b.Fatal(err)
	}

	// len(executionTxs) = tc.numTxsPerGen
	executionTxs, _, err := makeTxsWithNonceMap(false, nonceMap, activeAddrs, activeKeys, activeAddrs, signer, tc.numTxsPerGen, randomIndex)
	if err != nil {
		b.Fatal(err)
	}

	fmt.Println("len(warmUpTxs)", len(warmUpTxs), "len(executionTxs)", len(executionTxs))

	bcData.bc.Stop()
	bcData.db.Close()

	/////////////////////////////////////////////////////////////////////////////////////////////
	///  Tx Execution Process. Shutdown and initialize DB to execute txs in fresh condition.  ///
	/////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Println("Re-setting up test data dir to make database to initial condition.")
	testDataDir, err = setupTestDir(tc.originalDataDir, tc.isGenerateTest)
	if err != nil {
		b.Fatalf("err: %v, dir: %v", err, testDataDir)
	}

	bcData, err = NewBCDataForPreGeneratedTest(testDataDir, tc)
	if err != nil {
		b.Fatal(err)
	}
	defer bcData.db.Close()
	defer bcData.bc.Stop()

	signer = types.MakeSigner(bcData.bc.Config(), bcData.bc.CurrentHeader().Number)
	stateDB, err = bcData.bc.State()
	if err != nil {
		b.Fatal(err)
	}

	fmt.Println("Call AsMessageWithAccountKeyPicker for warmUpTxs")
	for _, tx := range warmUpTxs {
		if _, err = tx.AsMessageWithAccountKeyPicker(signer, stateDB, bcData.bc.CurrentBlock().NumberU64()); err != nil {
			b.Fatal(err)
		}
	}

	fmt.Println("Call AsMessageWithAccountKeyPicker for executionTxs")
	for _, tx := range executionTxs {
		if _, err = tx.AsMessageWithAccountKeyPicker(signer, stateDB, bcData.bc.CurrentBlock().NumberU64()); err != nil {
			b.Fatal(err)
		}
	}

	txPool := makeTxPool(bcData, txPoolSize)

	// Run warmUpTxs.
	fmt.Println("Start warming-up phase")
	if err := executeTxs(bcData, txPool, warmUpTxs); err != nil {
		b.Fatal(err)
	}
	fmt.Println("End warming-up phase. Start execution phase.")
	// Run executionTxs and measure the execution time and profiling.
	// Start timer and profiler from here.
	b.ResetTimer()
	b.StartTimer()
	defer b.StopTimer()
	pprof.StartCPUProfile(profileFile)
	defer pprof.StopCPUProfile()

	if err := executeTxs(bcData, txPool, executionTxs); err != nil {
		b.Fatal(err)
	}
}

// executeTxs repeats pushing and executing certain number of transactions until
// there is no transaction left.
func executeTxs(bcData *BCData, txPool *blockchain.TxPool, txs types.Transactions) error {
	for i := 0; i < len(txs); i += txPoolSize {
		end := i + txPoolSize
		if end > len(txs) {
			end = len(txs)
		}
		txPool.AddRemotes(txs[i:end])
		for {
			if err := bcData.GenABlockWithTxPoolWithoutAccountMap(txPool); err != nil {
				if err == errEmptyPending {
					break
				}
				return err
			}
		}
	}
	return nil
}

// getExecutionTestDefaultTC returns default TC of data execution tests.
func getExecutionTestDefaultTC() *preGeneratedTC {
	numActiveAccounts := 100 * 10000
	numExecPhaseTxs := txPoolSize * 10

	return &preGeneratedTC{
		isGenerateTest:     false,
		numReceiversPerRun: numActiveAccounts, // number of accounts used for warming-up phase, which means "active accounts"
		numTxsPerGen:       numExecPhaseTxs,   // number of transactions executed during execution phase
		numTotalSenders:    10000,
		filePicker:         sequentialIndex,
		addrPicker:         sequentialIndex,
		cacheConfig:        defaultCacheConfig(),
	}
}
