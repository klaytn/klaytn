// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from core/genesis_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

// TestDefaultGenesisBlock tests the genesis block generation functions: DefaultGenesisBlock, DefaultBaobabGenesisBlock
func TestDefaultGenesisBlock(t *testing.T) {
	block := genCypressGenesisBlock().ToBlock(common.Hash{}, nil)
	if block.Hash() != params.CypressGenesisHash {
		t.Errorf("wrong cypress genesis hash, got %v, want %v", block.Hash(), params.CypressGenesisHash)
	}
	block = genBaobabGenesisBlock().ToBlock(common.Hash{}, nil)
	if block.Hash() != params.BaobabGenesisHash {
		t.Errorf("wrong baobab genesis hash, got %v, want %v", block.Hash(), params.BaobabGenesisHash)
	}
}

// TestHardCodedChainConfigUpdate tests the public network's chainConfig update.
func TestHardCodedChainConfigUpdate(t *testing.T) {
	cypressGenesisBlock, baobabGenesisBlock := genCypressGenesisBlock(), genBaobabGenesisBlock()
	tests := []struct {
		name             string
		newHFBlock       *big.Int
		originHFBlock    *big.Int
		fn               func(database.DBManager, *big.Int) (*params.ChainConfig, common.Hash, error)
		wantConfig       *params.ChainConfig // expect value of the SetupGenesisBlock's first return value
		wantHash         common.Hash
		wantErr          error
		wantStoredConfig *params.ChainConfig // expect value of the stored config in DB
		resetFn          func(*big.Int)
	}{
		{
			name:       "cypress chainConfig update",
			newHFBlock: big.NewInt(3),
			fn: func(db database.DBManager, newHFBlock *big.Int) (*params.ChainConfig, common.Hash, error) {
				cypressGenesisBlock.MustCommit(db)
				cypressGenesisBlock.Config.IstanbulCompatibleBlock = newHFBlock
				return SetupGenesisBlock(db, cypressGenesisBlock, params.CypressNetworkId, false, false)
			},
			wantHash:         params.CypressGenesisHash,
			wantConfig:       cypressGenesisBlock.Config,
			wantStoredConfig: cypressGenesisBlock.Config,
		},
		// TODO-klaytn: add more cypress test cases after cypress hard fork block numbers are added
		{
			// Because of the fork-ordering check logic, the istanbulCompatibleBlock should be less than the londonCompatibleBlock
			name:       "baobab chainConfig update - correct hard-fork block number order",
			newHFBlock: big.NewInt(79999999),
			fn: func(db database.DBManager, newHFBlock *big.Int) (*params.ChainConfig, common.Hash, error) {
				baobabGenesisBlock.MustCommit(db)
				baobabGenesisBlock.Config.IstanbulCompatibleBlock = newHFBlock
				return SetupGenesisBlock(db, baobabGenesisBlock, params.BaobabNetworkId, false, false)
			},
			wantHash:         params.BaobabGenesisHash,
			wantConfig:       baobabGenesisBlock.Config,
			wantStoredConfig: baobabGenesisBlock.Config,
		},
		{
			// This test fails because the new istanbulCompatibleBlock(90909999) is larger than londonCompatibleBlock(80295291)
			name:       "baobab chainConfig update - wrong hard-fork block number order",
			newHFBlock: big.NewInt(90909999),
			fn: func(db database.DBManager, newHFBlock *big.Int) (*params.ChainConfig, common.Hash, error) {
				baobabGenesisBlock.MustCommit(db)
				baobabGenesisBlock.Config.IstanbulCompatibleBlock = newHFBlock
				return SetupGenesisBlock(db, baobabGenesisBlock, params.BaobabNetworkId, false, false)
			},
			wantHash:         common.Hash{},
			wantConfig:       baobabGenesisBlock.Config,
			wantStoredConfig: nil,
			wantErr: fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
				"istanbulBlock", big.NewInt(90909999), "londonBlock", big.NewInt(80295291)),
		},
		{
			name:       "incompatible config in DB",
			newHFBlock: big.NewInt(3),
			fn: func(db database.DBManager, newHFBlock *big.Int) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with Istanbul transition at #2.
				// Advance to block #4, past the Istanbul transition block of customGenesis.
				genesis := genCypressGenesisBlock()
				genesisBlock := genesis.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, genesis.Config, gxhash.NewFullFaker(), vm.Config{})
				defer bc.Stop()

				blocks, _ := GenerateChain(genesis.Config, genesisBlock, gxhash.NewFaker(), db, 4, nil)
				bc.InsertChain(blocks)
				// This should return a compatibility error.
				newConfig := *genesis
				newConfig.Config.IstanbulCompatibleBlock = newHFBlock
				return SetupGenesisBlock(db, &newConfig, params.CypressNetworkId, true, false)
			},
			wantHash:         params.CypressGenesisHash,
			wantConfig:       cypressGenesisBlock.Config,
			wantStoredConfig: params.CypressChainConfig,
			wantErr: &params.ConfigCompatError{
				What:         "Istanbul Block",
				StoredConfig: params.CypressChainConfig.IstanbulCompatibleBlock,
				NewConfig:    big.NewInt(3),
				RewindTo:     2,
			},
		},
	}

	for _, test := range tests {
		db := database.NewMemoryDBManager()
		config, hash, err := test.fn(db, test.newHFBlock)

		// Check the return values
		assert.Equal(t, test.wantErr, err, test.name+": err is mismatching")
		assert.Equal(t, test.wantConfig, config, test.name+": config is mismatching")
		assert.Equal(t, test.wantHash, hash, test.name+": hash is mismatching")

		// Check stored genesis block
		if test.wantHash != (common.Hash{}) {
			stored := db.ReadBlock(test.wantHash, 0)
			assert.Equal(t, test.wantHash, stored.Hash(), test.name+": stored genesis block is not compatible")
		}

		// Check stored chainConfig
		storedChainConfig := db.ReadChainConfig(test.wantHash)
		assert.Equal(t, test.wantStoredConfig, storedChainConfig, test.name+": stored chainConfig is not compatible")
	}
}

func TestSetupGenesis(t *testing.T) {
	var (
		customGenesisHash = common.HexToHash("0x4eb4035b7a09619a9950c9a4751cc331843f2373ef38263d676b4a132ba4059c")
		customChainId     = uint64(4343)
		customGenesis     = genCustomGenesisBlock(customChainId)
	)
	tests := []struct {
		name       string
		fn         func(database.DBManager) (*params.ChainConfig, common.Hash, error)
		wantConfig *params.ChainConfig
		wantHash   common.Hash
		wantErr    error
	}{
		{
			name: "genesis without ChainConfig",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, new(Genesis), params.UnusedNetworkId, false, false)
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: params.AllGxhashProtocolChanges,
		},
		{
			name: "no block in DB, genesis == nil, cypress networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, nil, params.CypressNetworkId, false, false)
			},
			wantHash:   params.CypressGenesisHash,
			wantConfig: params.CypressChainConfig,
		},
		{
			name: "no block in DB, genesis == nil, baobab networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, nil, params.BaobabNetworkId, false, false)
			},
			wantHash:   params.BaobabGenesisHash,
			wantConfig: params.BaobabChainConfig,
		},
		{
			name: "no block in DB, genesis == customGenesis, private network",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, customGenesis, customChainId, true, false)
			},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "cypress block in DB, genesis == nil, cypress networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genCypressGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil, params.CypressNetworkId, false, false)
			},
			wantHash:   params.CypressGenesisHash,
			wantConfig: params.CypressChainConfig,
		},
		{
			name: "baobab block in DB, genesis == nil, baobab networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genBaobabGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil, params.BaobabNetworkId, false, false)
			},
			wantHash:   params.BaobabGenesisHash,
			wantConfig: params.BaobabChainConfig,
		},
		{
			name: "custom block in DB, genesis == nil, custom networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, nil, customChainId, true, false)
			},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "cypress block in DB, genesis == baobab",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genCypressGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, genBaobabGenesisBlock(), params.BaobabNetworkId, false, false)
			},
			wantErr:    &GenesisMismatchError{Stored: params.CypressGenesisHash, New: params.BaobabGenesisHash},
			wantHash:   params.BaobabGenesisHash,
			wantConfig: params.BaobabChainConfig,
		},
		{
			name: "baobab block in DB, genesis == cypress",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genBaobabGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, genCypressGenesisBlock(), params.CypressNetworkId, false, false)
			},
			wantErr:    &GenesisMismatchError{Stored: params.BaobabGenesisHash, New: params.CypressGenesisHash},
			wantHash:   params.CypressGenesisHash,
			wantConfig: params.CypressChainConfig,
		},
		{
			name: "cypress block in DB, genesis == custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genCypressGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, genCustomGenesisBlock(customChainId), customChainId, true, false)
			},
			wantErr:    &GenesisMismatchError{Stored: params.CypressGenesisHash, New: customGenesisHash},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "baobab block in DB, genesis == custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genBaobabGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, customGenesis, customChainId, true, false)
			},
			wantErr:    &GenesisMismatchError{Stored: params.BaobabGenesisHash, New: customGenesisHash},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "custom block in DB, genesis == cypress",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, genCypressGenesisBlock(), params.CypressNetworkId, false, false)
			},
			wantErr:    &GenesisMismatchError{Stored: customGenesisHash, New: params.CypressGenesisHash},
			wantHash:   params.CypressGenesisHash,
			wantConfig: params.CypressChainConfig,
		},
		{
			name: "custom block in DB, genesis == baobab",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, genBaobabGenesisBlock(), params.BaobabNetworkId, false, false)
			},
			wantErr:    &GenesisMismatchError{Stored: customGenesisHash, New: params.BaobabGenesisHash},
			wantHash:   params.BaobabGenesisHash,
			wantConfig: params.BaobabChainConfig,
		},
		{
			name: "compatible config in DB",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, customGenesis, customChainId, true, false)
			},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "incompatible config in DB",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with Istanbul transition at #2.
				// Advance to block #4, past the Istanbul transition block of customGenesis.
				genesis := customGenesis.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, customGenesis.Config, gxhash.NewFullFaker(), vm.Config{})
				defer bc.Stop()

				blocks, _ := GenerateChain(customGenesis.Config, genesis, gxhash.NewFaker(), db, 4, nil)
				bc.InsertChain(blocks)
				// This should return a compatibility error.
				newConfig := *customGenesis
				newConfig.Config.IstanbulCompatibleBlock = big.NewInt(3)
				return SetupGenesisBlock(db, &newConfig, customChainId, true, false)
			},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
			wantErr: &params.ConfigCompatError{
				What:         "Istanbul Block",
				StoredConfig: big.NewInt(2),
				NewConfig:    big.NewInt(3),
				RewindTo:     1,
			},
		},
	}

	for _, test := range tests {
		db := database.NewMemoryDBManager()
		config, hash, err := test.fn(db)
		// Check the return values.
		if !reflect.DeepEqual(err, test.wantErr) {
			spew := spew.ConfigState{DisablePointerAddresses: true, DisableCapacities: true}
			t.Errorf("%s: returned error %#v, want %#v", test.name, spew.NewFormatter(err), spew.NewFormatter(test.wantErr))
		}
		if !reflect.DeepEqual(config, test.wantConfig) {
			t.Errorf("%s:\nreturned %v\nwant     %v", test.name, config, test.wantConfig)
		}
		if hash != test.wantHash {
			t.Errorf("%s: returned hash %s, want %s", test.name, hash.Hex(), test.wantHash.Hex())
		} else if err == nil {
			// Check database content.
			stored := db.ReadBlock(test.wantHash, 0)
			if stored.Hash() != test.wantHash {
				t.Errorf("%s: block in DB has hash %s, want %s", test.name, stored.Hash(), test.wantHash)
			}
		}
	}
}

func genCypressGenesisBlock() *Genesis {
	copyOfCypressChainConfig := *params.CypressChainConfig
	genesis := DefaultGenesisBlock()
	genesis.Config = &copyOfCypressChainConfig
	genesis.Governance = SetGenesisGovernance(genesis)
	InitDeriveSha(genesis.Config.DeriveShaImpl)
	return genesis
}

func genBaobabGenesisBlock() *Genesis {
	copyOfBaobabChainConfig := *params.BaobabChainConfig
	genesis := DefaultBaobabGenesisBlock()
	genesis.Config = &copyOfBaobabChainConfig
	genesis.Governance = SetGenesisGovernance(genesis)
	InitDeriveSha(genesis.Config.DeriveShaImpl)
	return genesis
}

func genCustomGenesisBlock(customChainId uint64) *Genesis {
	genesis := &Genesis{
		Config: &params.ChainConfig{
			ChainID:                 new(big.Int).SetUint64(customChainId),
			IstanbulCompatibleBlock: big.NewInt(2),
			DeriveShaImpl:           types.ImplDeriveShaConcat,
		},
		Alloc: GenesisAlloc{
			common.HexToAddress("0x0100000000000000000000000000000000000000"): {
				Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}},
			},
		},
	}
	genesis.Config.SetDefaults()
	genesis.Governance = SetGenesisGovernance(genesis)
	InitDeriveSha(genesis.Config.DeriveShaImpl)
	return genesis
}
