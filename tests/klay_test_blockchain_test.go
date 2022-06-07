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
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/work"

	istanbulBackend "github.com/klaytn/klaytn/consensus/istanbul/backend"
	istanbulCore "github.com/klaytn/klaytn/consensus/istanbul/core"
)

const transactionsJournalFilename = "transactions.rlp"

// If you don't want to remove 'chaindata', set removeChaindataOnExit = false
const removeChaindataOnExit = true

var errEmptyPending = errors.New("pending is empty")

type BCData struct {
	bc                 *blockchain.BlockChain
	addrs              []*common.Address
	privKeys           []*ecdsa.PrivateKey
	db                 database.DBManager
	rewardBase         *common.Address
	validatorAddresses []common.Address
	validatorPrivKeys  []*ecdsa.PrivateKey
	engine             consensus.Istanbul
	genesis            *blockchain.Genesis
	governance         governance.Engine
	rewardDistributor  *reward.RewardDistributor
}

var (
	dir      = "chaindata"
	nodeAddr = common.StringToAddress("nodeAddr")
)

func NewBCData(maxAccounts, numValidators int) (*BCData, error) {
	if numValidators > maxAccounts {
		return nil, errors.New("maxAccounts should be bigger numValidators!!")
	}

	// Remove leveldb dir if exists
	if _, err := os.Stat(dir); err == nil {
		os.RemoveAll(dir)
	}

	// Remove transactions.rlp if exists
	if _, err := os.Stat(transactionsJournalFilename); err == nil {
		os.RemoveAll(transactionsJournalFilename)
	}

	////////////////////////////////////////////////////////////////////////////////
	// Create a database
	chainDb := NewDatabase(dir, database.LevelDB)

	////////////////////////////////////////////////////////////////////////////////
	// Create a governance
	gov := generateGovernaceDataForTest()
	gov.SetNodeAddress(nodeAddr)
	////////////////////////////////////////////////////////////////////////////////
	// Create accounts as many as maxAccounts
	addrs, privKeys, err := createAccounts(maxAccounts)
	if err != nil {
		return nil, err
	}

	////////////////////////////////////////////////////////////////////////////////
	// Set the genesis address
	genesisAddr := *addrs[0]

	////////////////////////////////////////////////////////////////////////////////
	// Use the first `numValidators` accounts as validators
	validatorAddresses, validatorPrivKeys := getValidatorAddrsAndKeys(addrs, privKeys, numValidators)

	////////////////////////////////////////////////////////////////////////////////
	// Setup istanbul consensus backend
	engine := istanbulBackend.New(genesisAddr, istanbul.DefaultConfig, validatorPrivKeys[0], chainDb, gov, common.CONSENSUSNODE)

	////////////////////////////////////////////////////////////////////////////////
	// Make a blockchain
	bc, genesis, err := initBlockChain(chainDb, nil, addrs, validatorAddresses, nil, engine)
	if err != nil {
		return nil, err
	}

	engine.Start(bc, bc.CurrentBlock, bc.HasBadBlock)

	governance.AddGovernanceCacheForTest(gov, 0, genesis.Config)
	rewardDistributor := reward.NewRewardDistributor(gov)

	return &BCData{
		bc, addrs, privKeys, chainDb,
		&genesisAddr, validatorAddresses,
		validatorPrivKeys, engine, genesis, gov, rewardDistributor,
	}, nil
}

func (bcdata *BCData) Shutdown() {
	bcdata.bc.Stop()

	bcdata.db.Close()
	// Remove leveldb dir which was created for this test.
	if removeChaindataOnExit {
		os.RemoveAll(dir)
		os.RemoveAll(transactionsJournalFilename)
	}
}

func (bcdata *BCData) prepareHeader() (*types.Header, error) {
	tstart := time.Now()
	parent := bcdata.bc.CurrentBlock()

	tstamp := tstart.Unix()
	if parent.Time().Cmp(new(big.Int).SetInt64(tstamp)) >= 0 {
		tstamp = parent.Time().Int64() + 1
	}
	// this will ensure we're not going off too far in the future
	if now := time.Now().Unix(); tstamp > now {
		wait := time.Duration(tstamp-now) * time.Second
		time.Sleep(wait)
	}

	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		Time:       big.NewInt(tstamp),
		Governance: common.Hex2Bytes("b8dc7b22676f7665726e696e676e6f6465223a22307865373333636234643237396461363936663330643437306638633034646563623534666362306432222c22676f7665726e616e63656d6f6465223a2273696e676c65222c22726577617264223a7b226d696e74696e67616d6f756e74223a393630303030303030303030303030303030302c22726174696f223a2233342f33332f3333227d2c22626674223a7b2265706f6368223a33303030302c22706f6c696379223a302c22737562223a32317d2c22756e69745072696365223a32353030303030303030307d"),
		Vote:       common.Hex2Bytes("e194e733cb4d279da696f30d470f8c04decb54fcb0d28565706f6368853330303030"),
	}

	if err := bcdata.engine.Prepare(bcdata.bc, header); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to prepare header for mining %s.\n", err))
	}

	return header, nil
}

func (bcdata *BCData) MineABlock(transactions types.Transactions, signer types.Signer, prof *profile.Profiler) (*types.Block, types.Receipts, error) {
	// Set the block header
	start := time.Now()
	header, err := bcdata.prepareHeader()
	if err != nil {
		return nil, nil, err
	}
	prof.Profile("mine_prepareHeader", time.Now().Sub(start))

	statedb, err := bcdata.bc.State()
	if err != nil {
		return nil, nil, err
	}

	// Group transactions by the sender address
	start = time.Now()
	txs := make(map[common.Address]types.Transactions)
	for _, tx := range transactions {
		acc, err := types.Sender(signer, tx)
		if err != nil {
			return nil, nil, err
		}
		txs[acc] = append(txs[acc], tx)
	}
	prof.Profile("mine_groupTransactions", time.Now().Sub(start))

	// Create a transaction set where transactions are sorted by price and nonce
	start = time.Now()
	txset := types.NewTransactionsByPriceAndNonce(signer, txs)
	prof.Profile("mine_NewTransactionsByPriceAndNonce", time.Now().Sub(start))

	// Apply the set of transactions
	start = time.Now()
	task := work.NewTask(bcdata.bc.Config(), signer, statedb, header)
	task.ApplyTransactions(txset, bcdata.bc, *bcdata.rewardBase)
	newtxs := task.Transactions()
	receipts := task.Receipts()
	prof.Profile("mine_ApplyTransactions", time.Now().Sub(start))

	// Finalize the block
	start = time.Now()
	b, err := bcdata.engine.Finalize(bcdata.bc, header, statedb, newtxs, receipts)
	if err != nil {
		return nil, nil, err
	}
	prof.Profile("mine_finalize_block", time.Now().Sub(start))

	////////////////////////////////////////////////////////////////////////////////

	start = time.Now()
	b, err = sealBlock(b, bcdata.validatorPrivKeys)
	if err != nil {
		return nil, nil, err
	}
	prof.Profile("mine_seal_block", time.Now().Sub(start))

	return b, receipts, nil
}

func (bcdata *BCData) GenABlock(accountMap *AccountMap, opt *testOption,
	numTransactions int, prof *profile.Profiler,
) error {
	// Make a set of transactions
	start := time.Now()
	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)
	transactions, err := opt.makeTransactions(bcdata, accountMap, signer, numTransactions, nil, opt.txdata)
	if err != nil {
		return err
	}
	prof.Profile("main_makeTransactions", time.Now().Sub(start))

	return bcdata.GenABlockWithTransactions(accountMap, transactions, prof)
}

func (bcdata *BCData) GenABlockWithTxpool(accountMap *AccountMap, txpool *blockchain.TxPool,
	prof *profile.Profiler,
) error {
	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	pending, err := txpool.Pending()
	if err != nil {
		return err
	}
	if len(pending) == 0 {
		return errEmptyPending
	}
	pooltxs := types.NewTransactionsByPriceAndNonce(signer, pending)

	// Set the block header
	start := time.Now()
	header, err := bcdata.prepareHeader()
	if err != nil {
		return err
	}
	prof.Profile("mine_prepareHeader", time.Now().Sub(start))

	statedb, err := bcdata.bc.State()
	if err != nil {
		return err
	}

	start = time.Now()
	task := work.NewTask(bcdata.bc.Config(), signer, statedb, header)
	task.ApplyTransactions(pooltxs, bcdata.bc, *bcdata.rewardBase)
	newtxs := task.Transactions()
	receipts := task.Receipts()
	prof.Profile("mine_ApplyTransactions", time.Now().Sub(start))

	// Finalize the block
	start = time.Now()
	b, err := bcdata.engine.Finalize(bcdata.bc, header, statedb, newtxs, receipts)
	if err != nil {
		return err
	}
	prof.Profile("mine_finalize_block", time.Now().Sub(start))

	start = time.Now()
	b, err = sealBlock(b, bcdata.validatorPrivKeys)
	if err != nil {
		return err
	}
	prof.Profile("mine_seal_block", time.Now().Sub(start))

	// Update accountMap
	start = time.Now()
	if err := accountMap.Update(newtxs, signer, statedb, b.NumberU64()); err != nil {
		return err
	}
	prof.Profile("main_update_accountMap", time.Now().Sub(start))

	// Insert the block into the blockchain
	start = time.Now()
	if n, err := bcdata.bc.InsertChain(types.Blocks{b}); err != nil {
		return fmt.Errorf("err = %s, n = %d\n", err, n)
	}
	prof.Profile("main_insert_blockchain", time.Now().Sub(start))

	// Apply reward
	start = time.Now()
	if err := bcdata.rewardDistributor.MintKLAY(accountMap, header); err != nil {
		return err
	}
	prof.Profile("main_apply_reward", time.Now().Sub(start))

	// Verification with accountMap
	start = time.Now()
	statedbNew, err := bcdata.bc.State()
	if err != nil {
		return err
	}
	if err := accountMap.Verify(statedbNew); err != nil {
		return err
	}
	prof.Profile("main_verification", time.Now().Sub(start))

	return nil
}

func (bcdata *BCData) GenABlockWithTransactions(accountMap *AccountMap, transactions types.Transactions,
	prof *profile.Profiler,
) error {
	signer := types.MakeSigner(bcdata.bc.Config(), bcdata.bc.CurrentHeader().Number)

	statedb, err := bcdata.bc.State()
	if err != nil {
		return err
	}

	// Update accountMap
	start := time.Now()
	if err := accountMap.Update(transactions, signer, statedb, bcdata.bc.CurrentBlock().NumberU64()); err != nil {
		return err
	}
	prof.Profile("main_update_accountMap", time.Now().Sub(start))

	// Mine a block!
	start = time.Now()
	b, _, err := bcdata.MineABlock(transactions, signer, prof)
	if err != nil {
		return err
	}
	prof.Profile("main_mineABlock", time.Now().Sub(start))

	txs := make(types.Transactions, len(b.Transactions()))
	for i, tt := range b.Transactions() {
		encodedTx, err := rlp.EncodeToBytes(tt)
		if err != nil {
			return err
		}
		decodedTx := types.Transaction{}
		rlp.DecodeBytes(encodedTx, &decodedTx)
		txs[i] = &decodedTx
	}
	b = b.WithBody(txs)

	// Insert the block into the blockchain
	start = time.Now()
	if n, err := bcdata.bc.InsertChain(types.Blocks{b}); err != nil {
		return fmt.Errorf("err = %s, n = %d\n", err, n)
	}
	prof.Profile("main_insert_blockchain", time.Now().Sub(start))

	// Apply reward
	start = time.Now()
	if err := bcdata.rewardDistributor.MintKLAY(accountMap, b.Header()); err != nil {
		return err
	}
	prof.Profile("main_apply_reward", time.Now().Sub(start))

	// Verification with accountMap
	start = time.Now()
	statedb, err = bcdata.bc.State()
	if err != nil {
		return err
	}
	if err := accountMap.Verify(statedb); err != nil {
		return err
	}
	prof.Profile("main_verification", time.Now().Sub(start))

	return nil
}

////////////////////////////////////////////////////////////////////////////////
func NewDatabase(dir string, dbType database.DBType) database.DBManager {
	if dir == "" {
		return database.NewMemoryDBManager()
	} else {
		dbc := &database.DBConfig{
			Dir: dir, DBType: dbType, LevelDBCacheSize: 768,
			OpenFilesLimit: 1024, SingleDB: false, NumStateTrieShards: 4, ParallelDBWrite: true,
			LevelDBCompression: database.AllNoCompression, LevelDBBufferPool: true,
		}
		return database.NewDBManager(dbc)
	}
}

// Copied from consensus/istanbul/backend/engine.go
func prepareIstanbulExtra(validators []common.Address) ([]byte, error) {
	var buf bytes.Buffer

	buf.Write(bytes.Repeat([]byte{0x0}, types.IstanbulExtraVanity))

	ist := &types.IstanbulExtra{
		Validators:    validators,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}
	return append(buf.Bytes(), payload...), nil
}

func initBlockChain(db database.DBManager, cacheConfig *blockchain.CacheConfig, coinbaseAddrs []*common.Address, validators []common.Address,
	genesis *blockchain.Genesis, engine consensus.Engine,
) (*blockchain.BlockChain, *blockchain.Genesis, error) {
	extraData, err := prepareIstanbulExtra(validators)

	if genesis == nil {
		genesis = blockchain.DefaultGenesisBlock()
		genesis.Config = Forks["Byzantium"]
		genesis.ExtraData = extraData
		genesis.BlockScore = big.NewInt(1)
		genesis.Config.Governance = params.GetDefaultGovernanceConfig(params.UseIstanbul)
		genesis.Config.Istanbul = params.GetDefaultIstanbulConfig()
		genesis.Config.UnitPrice = 25 * params.Ston
	}

	alloc := make(blockchain.GenesisAlloc)
	for _, a := range coinbaseAddrs {
		alloc[*a] = blockchain.GenesisAccount{Balance: new(big.Int).Mul(big.NewInt(1e16), big.NewInt(params.KLAY))}
	}

	genesis.Alloc = alloc

	chainConfig, _, err := blockchain.SetupGenesisBlock(db, genesis, params.UnusedNetworkId, false, false)
	if _, ok := err.(*params.ConfigCompatError); err != nil && !ok {
		return nil, nil, err
	}

	// The chainConfig value has been modified while executing test. (ex, The test included executing applyTransaction())
	// Therefore, a deep copy is required to prevent the chainConfing value from being modified.
	var cfg params.ChainConfig
	b, err := json.Marshal(chainConfig)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return nil, nil, err
	}

	genesis.Config = &cfg

	chain, err := blockchain.NewBlockChain(db, cacheConfig, genesis.Config, engine, vm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return chain, genesis, nil
}

func createAccounts(numAccounts int) ([]*common.Address, []*ecdsa.PrivateKey, error) {
	accs := make([]*common.Address, numAccounts)
	privKeys := make([]*ecdsa.PrivateKey, numAccounts)

	for i := 0; i < numAccounts; i++ {
		k, err := crypto.GenerateKey()
		if err != nil {
			return nil, nil, err
		}
		keyAddr := crypto.PubkeyToAddress(k.PublicKey)

		accs[i] = &keyAddr
		privKeys[i] = k
	}

	return accs, privKeys, nil
}

// Copied from consensus/istanbul/backend/engine.go
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

// writeSeal writes the extra-data field of the given header with the given seals.
// Copied from consensus/istanbul/backend/engine.go
func writeSeal(h *types.Header, seal []byte) error {
	if len(seal)%types.IstanbulExtraSeal != 0 {
		return errors.New("invalid signature")
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}

// writeCommittedSeals writes the extra-data field of a block header with given committed seals.
// Copied from consensus/istanbul/backend/engine.go
func writeCommittedSeals(h *types.Header, committedSeals [][]byte) error {
	errInvalidCommittedSeals := errors.New("invalid committed seals")

	if len(committedSeals) == 0 {
		return errInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != types.IstanbulExtraSeal {
			return errInvalidCommittedSeals
		}
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(istanbulExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}

// sign implements istanbul.backend.Sign
// Copied from consensus/istanbul/backend/backend.go
func sign(data []byte, privkey *ecdsa.PrivateKey) ([]byte, error) {
	hashData := crypto.Keccak256([]byte(data))
	return crypto.Sign(hashData, privkey)
}

func makeCommittedSeal(h *types.Header, privKeys []*ecdsa.PrivateKey) ([][]byte, error) {
	committedSeals := make([][]byte, 0, 3)

	for i := 1; i < 4; i++ {
		seal := istanbulCore.PrepareCommittedSeal(h.Hash())
		committedSeal, err := sign(seal, privKeys[i])
		if err != nil {
			return nil, err
		}
		committedSeals = append(committedSeals, committedSeal)
	}

	return committedSeals, nil
}

func sealBlock(b *types.Block, privKeys []*ecdsa.PrivateKey) (*types.Block, error) {
	header := b.Header()

	seal, err := sign(sigHash(header).Bytes(), privKeys[0])
	if err != nil {
		return nil, err
	}

	err = writeSeal(header, seal)
	if err != nil {
		return nil, err
	}

	committedSeals, err := makeCommittedSeal(header, privKeys)
	if err != nil {
		return nil, err
	}

	err = writeCommittedSeals(header, committedSeals)
	if err != nil {
		return nil, err
	}

	return b.WithSeal(header), nil
}
