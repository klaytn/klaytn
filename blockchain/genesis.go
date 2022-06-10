// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/genesis.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
)

//go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var (
	errGenesisNoConfig = errors.New("genesis has no chain configuration")
	errNoGenesis       = errors.New("genesis block is not provided")
)

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config     *params.ChainConfig `json:"config"`
	Timestamp  uint64              `json:"timestamp"`
	ExtraData  []byte              `json:"extraData"`
	Governance []byte              `json:"governanceData"`
	BlockScore *big.Int            `json:"blockScore"`
	Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount

func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	BlockScore *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

// findBlockWithState returns the latest block with state.
func findBlockWithState(db database.DBManager) *types.Block {
	headBlock := db.ReadBlockByHash(db.ReadHeadBlockHash())
	if headBlock == nil {
		logger.Crit("failed to read head block by head block hash")
	}

	startBlock := headBlock
	for _, err := state.New(headBlock.Root(), state.NewDatabase(db), nil); err != nil; {
		if headBlock.NumberU64() == 0 {
			logger.Crit("failed to find state from the head block to the genesis block",
				"headBlockNum", headBlock.NumberU64(),
				"headBlockHash", headBlock.Hash().String(), "headBlockRoot", headBlock.Root().String())
		}
		headBlock = db.ReadBlockByNumber(headBlock.NumberU64() - 1)
		if headBlock == nil {
			logger.Crit("failed to read previous block by head block number")
		}
		logger.Warn("found previous block", "blockNum", headBlock.NumberU64())
	}
	logger.Info("found the latest block with state",
		"blockNum", headBlock.NumberU64(), "startedNum", startBlock.NumberU64())
	return headBlock
}

// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//                          genesis == nil                            genesis != nil
//                       +-------------------------------------------------------------------
//     db has no genesis |  main-net default, baobab if specified  |  genesis
//     db has genesis    |  from DB                                |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(db database.DBManager, genesis *Genesis, networkId uint64, isPrivate, overwriteGenesis bool) (*params.ChainConfig, common.Hash, error) {
	if genesis != nil && genesis.Config == nil {
		return params.AllGxhashProtocolChanges, common.Hash{}, errGenesisNoConfig
	}

	// Just commit the new block if there is no stored genesis block.
	stored := db.ReadCanonicalHash(0)
	if (stored == common.Hash{}) {
		if genesis == nil {
			switch {
			case isPrivate:
				logger.Error("No genesis is provided. --networkid should be omitted if you want to use preconfigured network")
				return params.AllGxhashProtocolChanges, common.Hash{}, errNoGenesis
			case networkId == params.BaobabNetworkId:
				logger.Info("Writing default baobab genesis block")
				genesis = DefaultBaobabGenesisBlock()
			case networkId == params.CypressNetworkId:
				fallthrough
			default:
				logger.Info("Writing default main-net genesis block")
				genesis = DefaultGenesisBlock()
			}
			if genesis.Config.Governance != nil {
				genesis.Governance = SetGenesisGovernance(genesis)
			}
		} else {
			logger.Info("Writing custom genesis block")
		}
		// Initialize DeriveSha implementation
		InitDeriveSha(genesis.Config.DeriveShaImpl)
		block, err := genesis.Commit(common.Hash{}, db)
		if err != nil {
			return genesis.Config, common.Hash{}, err
		}
		return genesis.Config, block.Hash(), err
	}

	// Check whether the genesis block is already written.
	if genesis != nil {
		// If overwriteGenesis is true, overwrite existing genesis block with the new one.
		// This is to run a test with pre-existing data.
		if overwriteGenesis {
			headBlock := findBlockWithState(db)
			logger.Warn("Trying to overwrite original genesis block with the new one",
				"headBlockHash", headBlock.Hash().String(), "headBlockNum", headBlock.NumberU64())
			newGenesisBlock, err := genesis.Commit(headBlock.Root(), db)
			return genesis.Config, newGenesisBlock.Hash(), err
		}
		// This is the usual path which does not overwrite genesis block with the new one.
		hash := genesis.ToBlock(common.Hash{}, nil).Hash()
		if hash != stored {
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
	}

	// Get the existing chain configuration.
	newcfg := genesis.configOrDefault(stored)
	if err := newcfg.CheckConfigForkOrder(); err != nil {
		return newcfg, common.Hash{}, err
	}
	storedcfg := db.ReadChainConfig(stored)
	if storedcfg == nil {
		logger.Info("Found genesis block without chain config")
		db.WriteChainConfig(stored, newcfg)
		return newcfg, stored, nil
	} else {
		if storedcfg.Governance == nil {
			logger.Crit("Failed to read governance. storedcfg.Governance == nil")
		}
		if storedcfg.Governance.Reward == nil {
			logger.Crit("Failed to read governance. storedcfg.Governance.Reward == nil")
		}
		if storedcfg.Governance.Reward.StakingUpdateInterval != 0 {
			params.SetStakingUpdateInterval(storedcfg.Governance.Reward.StakingUpdateInterval)
		}
		if storedcfg.Governance.Reward.ProposerUpdateInterval != 0 {
			params.SetProposerUpdateInterval(storedcfg.Governance.Reward.ProposerUpdateInterval)
		}
	}
	// Special case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	if genesis == nil && params.CypressGenesisHash != stored && params.BaobabGenesisHash != stored {
		return storedcfg, stored, nil
	}

	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := db.ReadHeaderNumber(db.ReadHeadHeaderHash())
	if height == nil {
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}
	compatErr := storedcfg.CheckCompatible(newcfg, *height)
	if compatErr != nil && *height != 0 && compatErr.RewindTo != 0 {
		return newcfg, stored, compatErr
	}
	db.WriteChainConfig(stored, newcfg)
	return newcfg, stored, nil
}

func (g *Genesis) configOrDefault(ghash common.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	case ghash == params.CypressGenesisHash:
		return params.CypressChainConfig
	case ghash == params.BaobabGenesisHash:
		return params.BaobabChainConfig
	default:
		return params.AllGxhashProtocolChanges
	}
}

// ToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func (g *Genesis) ToBlock(baseStateRoot common.Hash, db database.DBManager) *types.Block {
	if db == nil {
		db = database.NewMemoryDBManager()
	}
	stateDB, _ := state.New(baseStateRoot, state.NewDatabase(db), nil)
	for addr, account := range g.Alloc {
		if len(account.Code) != 0 {
			originalCode := stateDB.GetCode(addr)
			stateDB.SetCode(addr, account.Code)
			// If originalCode is not nil,
			// just update the code and don't change the other states
			if originalCode != nil {
				logger.Warn("this address already has a not nil code, now the code of this address has been changed", "addr", addr.String())
				continue
			}
		}
		for key, value := range account.Storage {
			stateDB.SetState(addr, key, value)
		}
		stateDB.AddBalance(addr, account.Balance)
		stateDB.SetNonce(addr, account.Nonce)
	}
	root := stateDB.IntermediateRoot(false)
	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Time:       new(big.Int).SetUint64(g.Timestamp),
		TimeFoS:    0,
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		Governance: g.Governance,
		GasUsed:    g.GasUsed,
		BlockScore: g.BlockScore,
		Root:       root,
	}
	if g.BlockScore == nil {
		head.BlockScore = params.GenesisBlockScore
	}
	stateDB.Commit(false)
	stateDB.Database().TrieDB().Commit(root, true, g.Number)

	return types.NewBlock(head, nil, nil)
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(baseStateRoot common.Hash, db database.DBManager) (*types.Block, error) {
	block := g.ToBlock(baseStateRoot, db)
	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}
	db.WriteTd(block.Hash(), block.NumberU64(), g.BlockScore)
	db.WriteBlock(block)
	db.WriteReceipts(block.Hash(), block.NumberU64(), nil)
	db.WriteCanonicalHash(block.Hash(), block.NumberU64())
	db.WriteHeadBlockHash(block.Hash())
	db.WriteHeadHeaderHash(block.Hash())

	config := g.Config
	if config == nil {
		config = params.AllGxhashProtocolChanges
	}
	if err := config.CheckConfigForkOrder(); err != nil {
		return nil, err
	}
	db.WriteChainConfig(block.Hash(), config)
	return block, nil
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db database.DBManager) *types.Block {
	config := g.Config
	if config == nil {
		config = params.AllGxhashProtocolChanges
	}
	InitDeriveSha(config.DeriveShaImpl)

	block, err := g.Commit(common.Hash{}, db)
	if err != nil {
		panic(err)
	}
	return block
}

// GenesisBlockForTesting creates and writes a block in which addr has the given peb balance.
func GenesisBlockForTesting(db database.DBManager, addr common.Address, balance *big.Int) *types.Block {
	g := Genesis{Alloc: GenesisAlloc{addr: {Balance: balance}}}
	return g.MustCommit(db)
}

// DefaultGenesisBlock returns the Cypress mainnet genesis block.
// It is also used for default genesis block.
func DefaultGenesisBlock() *Genesis {
	ret := &Genesis{}
	if err := json.Unmarshal(cypressGenesisJson, &ret); err != nil {
		logger.Error("Error in Unmarshalling Cypress Genesis Json", "err", err)
	}
	ret.Config = params.CypressChainConfig
	return ret
}

// DefaultBaobabGenesisBlock returns the Baobab testnet genesis block.
func DefaultBaobabGenesisBlock() *Genesis {
	ret := &Genesis{}
	if err := json.Unmarshal(baobabGenesisJson, &ret); err != nil {
		logger.Error("Error in Unmarshalling Baobab Genesis Json", "err", err)
		return nil
	}
	ret.Config = params.BaobabChainConfig
	return ret
}

func decodePrealloc(data string) GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
	}
	return ga
}

type GovernanceSet map[string]interface{}

func SetGenesisGovernance(genesis *Genesis) []byte {
	g := make(GovernanceSet)
	governance := genesis.Config.Governance
	g["governance.governancemode"] = governance.GovernanceMode
	g["governance.governingnode"] = governance.GoverningNode
	g["governance.unitprice"] = genesis.Config.UnitPrice
	g["reward.mintingamount"] = governance.Reward.MintingAmount.String()
	g["reward.minimumstake"] = governance.Reward.MinimumStake.String()
	g["reward.ratio"] = governance.Reward.Ratio
	g["reward.useginicoeff"] = governance.Reward.UseGiniCoeff
	g["reward.deferredtxfee"] = governance.Reward.DeferredTxFee
	g["reward.stakingupdateinterval"] = governance.Reward.StakingUpdateInterval
	g["reward.proposerupdateinterval"] = governance.Reward.ProposerUpdateInterval
	g["istanbul.epoch"] = genesis.Config.Istanbul.Epoch
	g["istanbul.policy"] = genesis.Config.Istanbul.ProposerPolicy
	g["istanbul.committeesize"] = genesis.Config.Istanbul.SubGroupSize

	data, err := json.Marshal(g)
	if err != nil {
		logger.Crit("Error in marshaling governance data", "err", err)
	}
	ret, err := rlp.EncodeToBytes(data)
	if err != nil {
		logger.Crit("Error in RLP Encoding governance data", "err", err)
	}
	return ret
}
