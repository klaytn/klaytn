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
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/otiai10/copy"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var errNoOriginalDataDir = errors.New("original data directory does not exist, aborting the test")

const (
	// All databases are compressed by Snappy, CompactionTableSize = 2MiB, CompactionTableSizeMultiplier = 1.0
	aspen500_orig = "aspen500_orig"
	// All databases are compressed by Snappy, CompactionTableSize = 4MiB, CompactionTableSizeMultiplier = 2.0
	baobab500_orig = "baobab500_orig"

	// Only receipt database is compressed by Snappy, CompactionTableSize = 2MiB, CompactionTableSizeMultiplier = 1.0
	candidate500LevelDB_orig = "candidate500LevelDB_orig"
	// Using BadgerDB with its default options.
	candidate500BadgerDB_orig = "candidate500BadgerDB_orig"

	// Same configuration as Baobab network, however only 10,000 accounts exist.
	baobab1_orig = "baobab1_orig"
)

// randomIndex is used to access data with random index.
func randomIndex(index, lenAddrs int) int {
	return rand.Intn(lenAddrs)
}

// sequentialIndex is used to access data with sequential index.
func sequentialIndex(index, lenAddrs int) int {
	return index % lenAddrs
}

// fixedIndex is used to access data with same index.
func fixedIndex(index int) func(int, int) int {
	return func(int, int) int {
		return index
	}
}

// makeTxsWithStateDB generates transactions with the nonce retrieved from stateDB.
// stateDB is used only once to initialize nonceMap, and then nonceMap is used instead of stateDB.
func makeTxsWithStateDB(isGenerate bool, stateDB *state.StateDB, fromAddrs []*common.Address, fromKeys []*ecdsa.PrivateKey, toAddrs []*common.Address, signer types.Signer, numTransactions int, indexPicker func(int, int) int) (types.Transactions, map[common.Address]uint64, error) {
	if len(fromAddrs) != len(fromKeys) {
		return nil, nil, fmt.Errorf("len(fromAddrs) %v != len(fromKeys) %v", len(fromAddrs), len(fromKeys))
	}

	// Use nonceMap, not to change the nonce of stateDB.
	nonceMap := make(map[common.Address]uint64)
	for _, addr := range fromAddrs {
		nonce := stateDB.GetNonce(*addr)
		nonceMap[*addr] = nonce
	}

	// Generate value transfer transactions from initial account to the given "toAddrs".
	return makeTxsWithNonceMap(isGenerate, nonceMap, fromAddrs, fromKeys, toAddrs, signer, numTransactions, indexPicker)
}

// makeTxsWithStateDB generates transactions with the nonce retrieved from nonceMap.
func makeTxsWithNonceMap(isGenerate bool, nonceMap map[common.Address]uint64, fromAddrs []*common.Address, fromKeys []*ecdsa.PrivateKey, toAddrs []*common.Address, signer types.Signer, numTransactions int, indexPicker func(int, int) int) (types.Transactions, map[common.Address]uint64, error) {
	txs := make(types.Transactions, 0, numTransactions)
	lenFromAddrs := len(fromAddrs)
	lenToAddrs := len(toAddrs)

	var transferValue *big.Int
	if isGenerate {
		transferValue = new(big.Int).Mul(big.NewInt(1e4), big.NewInt(params.KLAY))
	} else {
		transferValue = new(big.Int).Mul(big.NewInt(1e3), big.NewInt(params.Peb))
	}

	for i := 0; i < numTransactions; i++ {
		fromIdx := indexPicker(i, lenFromAddrs)
		toIdx := indexPicker(i, lenToAddrs)

		fromAddr := *fromAddrs[fromIdx]
		fromKey := fromKeys[fromIdx]
		fromNonce := nonceMap[fromAddr]

		toAddr := *toAddrs[toIdx]

		tx := types.NewTransaction(fromNonce, toAddr, transferValue, 1000000, new(big.Int).SetInt64(25000000000), nil)
		signedTx, err := types.SignTx(tx, signer, fromKey)
		if err != nil {
			return nil, nil, err
		}

		txs = append(txs, signedTx)
		nonceMap[fromAddr]++
	}

	return txs, nonceMap, nil
}

// setupTestDir does two things. If it is a data-generating test, it will just
// return the target path. If it is not a data-generating test, it will remove
// previously existing path and then copy the original data to the target path.
func setupTestDir(originalDataDirName string, isGenerateTest bool) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Original data directory should be located at github.com/klaytn
	// Therefore, it should be something like github.com/klaytn/testdata150_orig
	grandParentPath := filepath.Dir(filepath.Dir(wd))
	originalDataDirPath := path.Join(grandParentPath, originalDataDirName)

	// If it is generating test case, just returns the path.
	if isGenerateTest {
		return originalDataDirPath, nil
	}

	if _, err = os.Stat(originalDataDirPath); err != nil {
		return "", errNoOriginalDataDir
	}

	testDir := strings.Split(originalDataDirName, "_orig")[0]

	originalDataPath := path.Join(grandParentPath, originalDataDirName)
	testDataPath := path.Join(grandParentPath, testDir)

	os.RemoveAll(testDataPath)
	if err := copy.Copy(originalDataPath, testDataPath); err != nil {
		return "", err
	}
	return testDataPath, nil
}

type preGeneratedTC struct {
	isGenerateTest  bool
	testName        string
	originalDataDir string

	numTotalAccountsToGenerate int
	numTxsPerGen               int

	numTotalSenders    int // senders are loaded once at the test initialization time.
	numReceiversPerRun int // receivers are loaded repetitively for every tx generation run.

	filePicker func(int, int) int // determines the index of address file to use.
	addrPicker func(int, int) int // determines the index of address while making tx.

	dbc           *database.DBConfig
	levelDBOption *opt.Options
	cacheConfig   *blockchain.CacheConfig
}

// BenchmarkDataGeneration_Aspen generates the data with Aspen network's database configurations.
func BenchmarkDataGeneration_Aspen(b *testing.B) {
	tc := getGenerationTestDefaultTC()
	tc.testName = "BenchmarkDataGeneration_Aspen"
	tc.originalDataDir = aspen500_orig

	tc.cacheConfig = defaultCacheConfig()

	tc.dbc, tc.levelDBOption = genAspenOptions()

	dataGenerationTest(b, tc)
}

// BenchmarkDataGeneration_Baobab generates the data with Baobab network's database configurations.
func BenchmarkDataGeneration_Baobab(b *testing.B) {
	tc := getGenerationTestDefaultTC()
	tc.testName = "BenchmarkDataGeneration_Baobab"
	tc.originalDataDir = baobab500_orig

	tc.cacheConfig = defaultCacheConfig()

	tc.dbc, tc.levelDBOption = genBaobabOptions()

	dataGenerationTest(b, tc)
}

// BenchmarkDataGeneration_CandidateLevelDB generates the data for main-net's
// with candidate configurations, using LevelDB.
func BenchmarkDataGeneration_CandidateLevelDB(b *testing.B) {
	tc := getGenerationTestDefaultTC()
	tc.testName = "BenchmarkDataGeneration_CandidateLevelDB"
	tc.originalDataDir = candidate500LevelDB_orig

	tc.cacheConfig = defaultCacheConfig()

	tc.dbc, tc.levelDBOption = genCandidateLevelDBOptions()

	dataGenerationTest(b, tc)
}

// BenchmarkDataGeneration_CandidateBadgerDB generates the data for main-net's
// with candidate configurations, using BadgerDB.
func BenchmarkDataGeneration_CandidateBadgerDB(b *testing.B) {
	tc := getGenerationTestDefaultTC()
	tc.testName = "BenchmarkDataGeneration_CandidateBadgerDB"
	tc.originalDataDir = candidate500BadgerDB_orig

	tc.cacheConfig = defaultCacheConfig()

	tc.dbc, tc.levelDBOption = genCandidateBadgerDBOptions()

	dataGenerationTest(b, tc)
}

// BenchmarkDataGeneration_Baobab_ControlGroup generates the data with Baobab network's database configurations.
// To work as a control group, it only generates 10,000 accounts.
func BenchmarkDataGeneration_Baobab_ControlGroup(b *testing.B) {
	tc := getGenerationTestDefaultTC()
	tc.testName = "BenchmarkDataGeneration_Baobab_ControlGroup"
	tc.originalDataDir = baobab1_orig
	tc.numTotalAccountsToGenerate = 10000

	tc.cacheConfig = defaultCacheConfig()

	tc.dbc, tc.levelDBOption = genBaobabOptions()

	dataGenerationTest(b, tc)
}

// dataGenerationTest generates given number of accounts for pre-generated tests.
// Newly generated data directory will be located at "$GOPATH/src/github.com/klaytn/"
func dataGenerationTest(b *testing.B, tc *preGeneratedTC) {
	testDataDir, profileFile, err := setUpTest(tc)
	if err != nil {
		b.Fatal(err)
	}

	bcData, err := NewBCDataForPreGeneratedTest(testDataDir, tc)
	if err != nil {
		b.Fatal(err)
	}

	defer bcData.db.Close()
	defer bcData.bc.Stop()

	txPool := makeTxPool(bcData, tc.numTxsPerGen)
	signer := types.MakeSigner(bcData.bc.Config(), bcData.bc.CurrentHeader().Number)

	b.ResetTimer()
	b.StopTimer()

	numTxGenerationRuns := tc.numTotalAccountsToGenerate / tc.numTxsPerGen
	for run := 1; run < numTxGenerationRuns; run++ {
		toAddrs, _, err := makeOrGenerateAddrsAndKeys(testDataDir, run, tc)
		if err != nil {
			b.Fatal(err)
		}

		// Generate transactions
		stateDB, err := bcData.bc.State()
		if err != nil {
			b.Fatal(err)
		}

		txs, _, err := makeTxsWithStateDB(true, stateDB, bcData.addrs, bcData.privKeys, toAddrs, signer, tc.numTxsPerGen, tc.addrPicker)
		if err != nil {
			b.Fatal(err)
		}

		for _, tx := range txs {
			tx.AsMessageWithAccountKeyPicker(signer, stateDB, bcData.bc.CurrentBlock().NumberU64())
		}

		b.StartTimer()
		if run == numTxGenerationRuns {
			pprof.StartCPUProfile(profileFile)
		}

		txPool.AddRemotes(txs)

		for {
			if err := bcData.GenABlockWithTxPoolWithoutAccountMap(txPool); err != nil {
				if err == errEmptyPending {
					break
				}
				b.Fatal(err)
			}
		}

		if run == numTxGenerationRuns {
			pprof.StopCPUProfile()
		}
		b.StopTimer()
	}
}

// getGenerationTestDefaultTC returns default TC of data generation tests.
func getGenerationTestDefaultTC() *preGeneratedTC {
	return &preGeneratedTC{
		isGenerateTest:             true,
		numTotalAccountsToGenerate: 500 * 10000, numTxsPerGen: 10000,
		numTotalSenders: 10000, numReceiversPerRun: 10000,
		filePicker: sequentialIndex, addrPicker: sequentialIndex,
		cacheConfig: defaultCacheConfig(),
	}
}
