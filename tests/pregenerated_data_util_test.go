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
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	istanbulBackend "github.com/klaytn/klaytn/consensus/istanbul/backend"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/klaytn/klaytn/work"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	numValidatorsForTest = 4

	addressDirectory    = "addrs"
	privateKeyDirectory = "privatekeys"

	addressFilePrefix    = "addrs_"
	privateKeyFilePrefix = "privateKeys_"

	chainDataDir = "chaindata"
)

var totalTxs = 0

// writeToFile writes addresses and private keys to designated directories with given fileNum.
// Addresses are stored in a file like `addrs_0` and private keys are stored in a file like `privateKeys_0`.
func writeToFile(addrs []*common.Address, privKeys []*ecdsa.PrivateKey, fileNum int, dir string) {
	_ = os.Mkdir(path.Join(dir, addressDirectory), os.ModePerm)
	_ = os.Mkdir(path.Join(dir, privateKeyDirectory), os.ModePerm)

	addrsFile, err := os.Create(path.Join(dir, addressDirectory, addressFilePrefix+strconv.Itoa(fileNum)))
	if err != nil {
		panic(err)
	}

	privateKeysFile, err := os.Create(path.Join(dir, privateKeyDirectory, privateKeyFilePrefix+strconv.Itoa(fileNum)))
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(2)

	syncSize := len(addrs) / 2

	go func() {
		for i, b := range addrs {
			addrsFile.WriteString(b.String() + "\n")
			if (i+1)%syncSize == 0 {
				addrsFile.Sync()
			}
		}

		addrsFile.Close()
		wg.Done()
	}()

	go func() {
		for i, key := range privKeys {
			privateKeysFile.WriteString(hex.EncodeToString(crypto.FromECDSA(key)) + "\n")
			if (i+1)%syncSize == 0 {
				privateKeysFile.Sync()
			}
		}

		privateKeysFile.Close()
		wg.Done()
	}()

	wg.Wait()
}

func readAddrsFromFile(dir string, num int) ([]*common.Address, error) {
	var addrs []*common.Address

	addrsFile, err := os.Open(path.Join(dir, addressDirectory, addressFilePrefix+strconv.Itoa(num)))
	if err != nil {
		return nil, err
	}

	defer addrsFile.Close()

	scanner := bufio.NewScanner(addrsFile)
	for scanner.Scan() {
		keyStr := scanner.Text()
		addr := common.HexToAddress(keyStr)
		addrs = append(addrs, &addr)
	}

	return addrs, nil
}

func readPrivateKeysFromFile(dir string, num int) ([]*ecdsa.PrivateKey, error) {
	var privKeys []*ecdsa.PrivateKey
	privateKeysFile, err := os.Open(path.Join(dir, privateKeyDirectory, privateKeyFilePrefix+strconv.Itoa(num)))
	if err != nil {
		return nil, err
	}

	defer privateKeysFile.Close()

	scanner := bufio.NewScanner(privateKeysFile)
	for scanner.Scan() {
		keyStr := scanner.Text()

		key, err := hex.DecodeString(keyStr)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}

		if pk, err := crypto.ToECDSA(key); err != nil {
			return nil, fmt.Errorf("%v", err)
		} else {
			privKeys = append(privKeys, pk)
		}
	}

	return privKeys, nil
}

func readAddrsAndPrivateKeysFromFile(dir string, num int) ([]*common.Address, []*ecdsa.PrivateKey, error) {
	addrs, err := readAddrsFromFile(dir, num)
	if err != nil {
		return nil, nil, err
	}

	privateKeys, err := readPrivateKeysFromFile(dir, num)
	if err != nil {
		return nil, nil, err
	}

	return addrs, privateKeys, nil
}

// getAddrsAndKeysFromFile extracts the address stored in file by numAccounts.
func getAddrsAndKeysFromFile(numAccounts int, testDataDir string, run int, filePicker func(int, int) int) ([]*common.Address, []*ecdsa.PrivateKey, error) {
	addrs := make([]*common.Address, 0, numAccounts)
	privKeys := make([]*ecdsa.PrivateKey, 0, numAccounts)

	files, err := ioutil.ReadDir(path.Join(testDataDir, addressDirectory))
	if err != nil {
		return nil, nil, err
	}

	numFiles := len(files)
	remain := numAccounts
	for i := run; remain > 0; i++ {
		fileIndex := filePicker(i, numFiles)

		// Read recipient addresses from file.
		addrsPerFile, privKeysPerFile, err := readAddrsAndPrivateKeysFromFile(testDataDir, fileIndex)
		if err != nil {
			return nil, nil, err
		}

		fmt.Println("Read addresses from " + addressFilePrefix + strconv.Itoa(fileIndex) + ", len(addrs)=" + strconv.Itoa(len(addrsPerFile)))

		partSize := int(math.Min(float64(len(addrsPerFile)), float64(remain)))
		addrs = append(addrs, addrsPerFile[:partSize]...)
		privKeys = append(privKeys, privKeysPerFile[:partSize]...)
		remain -= partSize
	}

	return addrs, privKeys, nil
}

func makeOrGenerateAddrsAndKeys(testDataDir string, run int, tc *preGeneratedTC) ([]*common.Address, []*ecdsa.PrivateKey, error) {
	// If it is execution test, makeOrGenerateAddrsAndKeys is called to build validator list.
	if !tc.isGenerateTest {
		return getAddrsAndKeysFromFile(numValidatorsForTest, testDataDir, run, tc.filePicker)
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	// If addrs directory exists in the current working directory for reuse,
	// use it instead of generating addresses.
	addrPathInWD := path.Join(wd, addressDirectory)
	if _, err := os.Stat(addrPathInWD); err == nil {
		return getAddrsAndKeysFromFile(tc.numReceiversPerRun, wd, run, tc.filePicker)
	}

	// No addrs directory exists. Generating and saving addresses and keys.
	var newPrivateKeys []*ecdsa.PrivateKey
	toAddrs, newPrivateKeys, err := createAccounts(tc.numTxsPerGen)
	if err != nil {
		return nil, nil, err
	}

	writeToFile(toAddrs, newPrivateKeys, run, testDataDir)
	return toAddrs, newPrivateKeys, nil
}

// getValidatorAddrsAndKeys returns the first `numValidators` addresses and private keys
// for validators.
func getValidatorAddrsAndKeys(addrs []*common.Address, privateKeys []*ecdsa.PrivateKey, numValidators int) ([]common.Address, []*ecdsa.PrivateKey) {
	validatorAddresses := make([]common.Address, numValidators)
	validatorPrivateKeys := make([]*ecdsa.PrivateKey, numValidators)

	for i := 0; i < numValidators; i++ {
		validatorPrivateKeys[i] = privateKeys[i]
		validatorAddresses[i] = *addrs[i]
	}

	return validatorAddresses, validatorPrivateKeys
}

// GenABlockWithTxPoolWithoutAccountMap basically does the same thing with GenABlockWithTxPool,
// however, it does not accept AccountMap which validates the outcome with stateDB.
// This is to remove the overhead of AccountMap management.
func (bcdata *BCData) GenABlockWithTxPoolWithoutAccountMap(txPool *blockchain.TxPool) error {
	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	pending, err := txPool.Pending()
	if err != nil {
		return err
	}
	if len(pending) == 0 {
		return errEmptyPending
	}

	pooltxs := types.NewTransactionsByPriceAndNonce(signer, pending)

	// Set the block header
	header, err := bcdata.prepareHeader()
	if err != nil {
		return err
	}

	stateDB, err := bcdata.bc.StateAt(bcdata.bc.CurrentBlock().Root())
	if err != nil {
		return err
	}

	task := work.NewTask(bcdata.bc.Config(), signer, stateDB, header)
	task.ApplyTransactions(pooltxs, bcdata.bc, *bcdata.rewardBase)
	newtxs := task.Transactions()
	receipts := task.Receipts()

	if len(newtxs) == 0 {
		return errEmptyPending
	}

	// Finalize the block.
	b, err := bcdata.engine.Finalize(bcdata.bc, header, stateDB, newtxs, receipts)
	if err != nil {
		return err
	}

	// Seal the block.
	b, err = sealBlock(b, bcdata.validatorPrivKeys)
	if err != nil {
		return err
	}

	// Write the block with state.
	result, err := bcdata.bc.WriteBlockWithState(b, receipts, stateDB)
	if err != nil {
		return fmt.Errorf("err = %s", err)
	}

	if result.Status == blockchain.SideStatTy {
		return fmt.Errorf("forked block is generated! number: %v, hash: %v, txs: %v", b.Number(), b.Hash(), len(b.Transactions()))
	}

	// Trigger post chain events after successful writing.
	logs := stateDB.Logs()
	events := []interface{}{blockchain.ChainEvent{Block: b, Hash: b.Hash(), Logs: logs}}
	events = append(events, blockchain.ChainHeadEvent{Block: b})
	bcdata.bc.PostChainEvents(events, logs)

	totalTxs += len(newtxs)
	fmt.Println("blockNum", b.NumberU64(), "numTxs", len(newtxs), "totalTxs", totalTxs)

	return nil
}

// NewBCDataForPreGeneratedTest returns a new BCData pointer constructed either 1) from the scratch or 2) from the existing data.
func NewBCDataForPreGeneratedTest(testDataDir string, tc *preGeneratedTC) (*BCData, error) {
	totalTxs = 0

	if numValidatorsForTest > tc.numTotalSenders {
		return nil, errors.New("numTotalSenders should be bigger numValidatorsForTest")
	}

	// Remove test data directory if 1) exists and and 2) generating test.
	if _, err := os.Stat(testDataDir); err == nil && tc.isGenerateTest {
		os.RemoveAll(testDataDir)
	}

	// Remove transactions.rlp if exists
	if _, err := os.Stat(transactionsJournalFilename); err == nil {
		os.RemoveAll(transactionsJournalFilename)
	}

	////////////////////////////////////////////////////////////////////////////////
	// Create a database
	tc.dbc.Dir = path.Join(testDataDir, chainDataDir)
	fmt.Println("DBDir", tc.dbc.Dir)

	var chainDB database.DBManager
	var err error
	if tc.dbc.DBType == database.LevelDB {
		chainDB, err = database.NewLevelDBManagerForTest(tc.dbc, tc.levelDBOption)
	} else if tc.dbc.DBType == database.BadgerDB {
		chainDB = database.NewDBManager(tc.dbc)
	}

	if err != nil {
		return nil, err
	}

	////////////////////////////////////////////////////////////////////////////////
	// Create a governance
	gov := generateGovernaceDataForTest()

	////////////////////////////////////////////////////////////////////////////////
	// Prepare sender addresses and private keys
	// 1) If generating test, create accounts and private keys as many as numTotalSenders
	// 2) If executing test, load accounts and private keys from file as many as numTotalSenders
	addrs, privKeys, err := makeOrGenerateAddrsAndKeys(testDataDir, 0, tc)
	if err != nil {
		return nil, err
	}

	////////////////////////////////////////////////////////////////////////////////
	// Set the genesis address
	genesisAddr := *addrs[0]

	////////////////////////////////////////////////////////////////////////////////
	// Use the first `numValidatorsForTest` accounts as validators
	validatorAddresses, validatorPrivKeys := getValidatorAddrsAndKeys(addrs, privKeys, numValidatorsForTest)

	////////////////////////////////////////////////////////////////////////////////
	// Setup istanbul consensus backend
	engine := istanbulBackend.New(genesisAddr, istanbul.DefaultConfig, validatorPrivKeys[0], chainDB, gov, common.CONSENSUSNODE)

	////////////////////////////////////////////////////////////////////////////////
	// Make a blockChain
	// 1) If generating test, call initBlockChain
	// 2) If executing test, call blockchain.NewBlockChain
	var bc *blockchain.BlockChain
	var genesis *blockchain.Genesis
	if tc.isGenerateTest {
		bc, genesis, err = initBlockChain(chainDB, tc.cacheConfig, addrs, validatorAddresses, nil, engine)
	} else {
		chainConfig, err := getChainConfig(chainDB)
		if err != nil {
			return nil, err
		}
		genesis = blockchain.DefaultGenesisBlock()
		genesis.Config = chainConfig
		bc, err = blockchain.NewBlockChain(chainDB, tc.cacheConfig, chainConfig, engine, vm.Config{})
	}

	if err != nil {
		return nil, err
	}

	rewardDistributor := reward.NewRewardDistributor(gov)

	return &BCData{
		bc, addrs, privKeys, chainDB,
		&genesisAddr, validatorAddresses,
		validatorPrivKeys, engine, genesis, gov, rewardDistributor,
	}, nil
}

// genAspenOptions returns database configurations of Aspen network.
func genAspenOptions() (*database.DBConfig, *opt.Options) {
	aspenDBConfig := defaultDBConfig()
	aspenDBConfig.DBType = database.LevelDB
	aspenDBConfig.LevelDBCompression = database.AllSnappyCompression

	aspenLevelDBOptions := &opt.Options{WriteBuffer: 256 * opt.MiB}

	return aspenDBConfig, aspenLevelDBOptions
}

// genBaobabOptions returns database configurations of Baobab network.
func genBaobabOptions() (*database.DBConfig, *opt.Options) {
	dbc, opts := genAspenOptions()

	opts.CompactionTableSize = 4 * opt.MiB
	opts.CompactionTableSizeMultiplier = 2

	dbc.LevelDBCompression = database.AllSnappyCompression

	return dbc, opts
}

// genCandidateLevelDBOptions returns candidate database configurations of main-net, using LevelDB.
func genCandidateLevelDBOptions() (*database.DBConfig, *opt.Options) {
	dbc, opts := genAspenOptions()

	dbc.LevelDBBufferPool = true
	dbc.LevelDBCompression = database.AllNoCompression

	return dbc, opts
}

// genCandidateLevelDBOptions returns candidate database configurations of main-net, using BadgerDB.
func genCandidateBadgerDBOptions() (*database.DBConfig, *opt.Options) {
	dbc := defaultDBConfig()
	dbc.DBType = database.BadgerDB

	return dbc, &opt.Options{}
}

// defaultDBConfig returns default database.DBConfig for pre-generated tests.
func defaultDBConfig() *database.DBConfig {
	return &database.DBConfig{SingleDB: false, ParallelDBWrite: true, NumStateTrieShards: 4}
}

// getChainConfig returns chain config from chainDB.
func getChainConfig(chainDB database.DBManager) (*params.ChainConfig, error) {
	stored := chainDB.ReadBlockByNumber(0)
	if stored == nil {
		return nil, errors.New("chainDB.ReadBlockByNumber(0) == nil")
	}

	chainConfig := chainDB.ReadChainConfig(stored.Hash())
	if chainConfig == nil {
		return nil, errors.New("chainConfig == nil")
	}

	return chainConfig, nil
}

// defaultCacheConfig returns cache config for data generation tests.
func defaultCacheConfig() *blockchain.CacheConfig {
	return &blockchain.CacheConfig{
		ArchiveMode:   false,
		CacheSize:     512,
		BlockInterval: blockchain.DefaultBlockInterval,
		TriesInMemory: blockchain.DefaultTriesInMemory,
		TrieNodeCacheConfig: &statedb.TrieNodeCacheConfig{
			CacheType:          statedb.CacheTypeLocal,
			LocalCacheSizeMiB:  4096,
			FastCacheFileDir:   "",
			RedisEndpoints:     nil,
			RedisClusterEnable: false,
		},
		SnapshotCacheSize: 512,
	}
}

// generateGovernaceDataForTest returns governance.Engine for test.
func generateGovernaceDataForTest() governance.Engine {
	dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})

	return governance.NewMixedEngine(&params.ChainConfig{
		ChainID:       big.NewInt(2018),
		UnitPrice:     25000000000,
		DeriveShaImpl: 0,
		Istanbul: &params.IstanbulConfig{
			Epoch:          istanbul.DefaultConfig.Epoch,
			ProposerPolicy: uint64(istanbul.DefaultConfig.ProposerPolicy),
			SubGroupSize:   istanbul.DefaultConfig.SubGroupSize,
		},
		Governance: params.GetDefaultGovernanceConfig(params.UseIstanbul),
	}, dbm)
}

// setUpTest sets up test data directory, verbosity and profile file.
func setUpTest(tc *preGeneratedTC) (string, *os.File, error) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	testDataDir, err := setupTestDir(tc.originalDataDir, tc.isGenerateTest)
	if err != nil {
		return "", nil, fmt.Errorf("err: %v, dir: %v", err, testDataDir)
	}

	timeNow := time.Now()
	f, err := os.Create(tc.testName + "_" + timeNow.Format("2006-01-02-1504") + ".cpu.out")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create file for cpu profiling, err: %v", err)
	}

	return testDataDir, f, nil
}
