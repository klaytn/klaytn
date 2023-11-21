// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from params/config.go (2018/06/04).
// Modified and improved for the klaytn development.

package params

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
)

// Genesis hashes to enforce below configs on.
var (
	CypressGenesisHash      = common.HexToHash("0xc72e5293c3c3ba38ed8ae910f780e4caaa9fb95e79784f7ab74c3c262ea7137e") // cypress genesis hash to enforce below configs on
	BaobabGenesisHash       = common.HexToHash("0xe33ff05ceec2581ca9496f38a2bf9baad5d4eed629e896ccb33d1dc991bc4b4a") // baobab genesis hash to enforce below configs on
	AuthorAddressForTesting = common.HexToAddress("0xc0ea08a2d404d3172d2add29a45be56da40e2949")
	mintingAmount, _        = new(big.Int).SetString("9600000000000000000", 10)
	logger                  = log.NewModuleLogger(log.Governance)
)

var (
	// CypressChainConfig is the chain parameters to run a node on the cypress main network.
	CypressChainConfig = &ChainConfig{
		ChainID:                  big.NewInt(int64(CypressNetworkId)),
		IstanbulCompatibleBlock:  big.NewInt(86816005),
		LondonCompatibleBlock:    big.NewInt(86816005),
		EthTxTypeCompatibleBlock: big.NewInt(86816005),
		MagmaCompatibleBlock:     big.NewInt(99841497),
		KoreCompatibleBlock:      big.NewInt(119750400),
		ShanghaiCompatibleBlock:  big.NewInt(135456000),
		CancunCompatibleBlock:    nil, // TODO-Klaytn-Cancun: set Cypress CancunCompatibleBlock
		Kip103CompatibleBlock:    big.NewInt(119750400),
		Kip103ContractAddress:    common.HexToAddress("0xD5ad6D61Dd87EdabE2332607C328f5cc96aeCB95"),
		DeriveShaImpl:            2,
		Governance: &GovernanceConfig{
			GoverningNode:  common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"),
			GovernanceMode: "single",
			Reward: &RewardConfig{
				MintingAmount:          mintingAmount,
				Ratio:                  "34/54/12",
				UseGiniCoeff:           true,
				DeferredTxFee:          true,
				StakingUpdateInterval:  86400,
				ProposerUpdateInterval: 3600,
				MinimumStake:           big.NewInt(5000000),
			},
		},
		Istanbul: &IstanbulConfig{
			Epoch:          604800,
			ProposerPolicy: 2,
			SubGroupSize:   22,
		},
		UnitPrice: 25000000000,
	}

	// BaobabChainConfig contains the chain parameters to run a node on the Baobab test network.
	BaobabChainConfig = &ChainConfig{
		ChainID:                  big.NewInt(int64(BaobabNetworkId)),
		IstanbulCompatibleBlock:  big.NewInt(75373312),
		LondonCompatibleBlock:    big.NewInt(80295291),
		EthTxTypeCompatibleBlock: big.NewInt(86513895),
		MagmaCompatibleBlock:     big.NewInt(98347376),
		KoreCompatibleBlock:      big.NewInt(111736800),
		ShanghaiCompatibleBlock:  big.NewInt(131608000),
		CancunCompatibleBlock:    nil, // TODO-Klaytn-Cancun: set Baobab CancunCompatibleBlock
		Kip103CompatibleBlock:    big.NewInt(119145600),
		Kip103ContractAddress:    common.HexToAddress("0xD5ad6D61Dd87EdabE2332607C328f5cc96aeCB95"),
		DeriveShaImpl:            2,
		Governance: &GovernanceConfig{
			GoverningNode:  common.HexToAddress("0x99fb17d324fa0e07f23b49d09028ac0919414db6"),
			GovernanceMode: "single",
			Reward: &RewardConfig{
				MintingAmount:          mintingAmount,
				Ratio:                  "34/54/12",
				UseGiniCoeff:           true,
				DeferredTxFee:          true,
				StakingUpdateInterval:  86400,
				ProposerUpdateInterval: 3600,
				MinimumStake:           big.NewInt(5000000),
			},
		},
		Istanbul: &IstanbulConfig{
			Epoch:          604800,
			ProposerPolicy: 2,
			SubGroupSize:   22,
		},
		UnitPrice: 25000000000,
	}

	// AllGxhashProtocolChanges contains every protocol change (GxIPs) introduced
	// and accepted by the klaytn developers into the Klaytn consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllGxhashProtocolChanges = &ChainConfig{
		ChainID:  big.NewInt(0),
		Gxhash:   new(GxhashConfig),
		Clique:   nil,
		Istanbul: nil,
	}

	// AllCliqueProtocolChanges contains every protocol change (GxIPs) introduced
	// and accepted by the klaytn developers into the Clique consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllCliqueProtocolChanges = &ChainConfig{
		ChainID:  big.NewInt(0),
		Gxhash:   nil,
		Clique:   &CliqueConfig{Period: 0, Epoch: 30000},
		Istanbul: nil,
	}

	TestChainConfig = &ChainConfig{
		ChainID:       big.NewInt(1),
		Gxhash:        new(GxhashConfig),
		Clique:        nil,
		Istanbul:      nil,
		UnitPrice:     1, // NOTE-Klaytn Use UnitPrice 1 for tests
		DeriveShaImpl: 0,
	}
	TestRules = TestChainConfig.Rules(new(big.Int))

	// istanbul BFT
	BFTTestChainConfig = &ChainConfig{
		ChainID:  big.NewInt(1),
		Gxhash:   new(GxhashConfig),
		Clique:   nil,
		Istanbul: nil,
	}
)

// VMLogTarget sets the output target of vmlog.
// The values below can be OR'ed.
//   - 0x0: no output (default)
//   - 0x1: file (DATADIR/logs/vm.log)
//   - 0x2: stdout (like logger.DEBUG)
var VMLogTarget = 0x0

const (
	VMLogToFile   = 0x1
	VMLogToStdout = 0x2
	VMLogToAll    = VMLogToFile | VMLogToStdout

	UpperGasLimit = uint64(999999999999)
)

const (
	PasswordLength = 16
)

// ChainConfig is the blockchain config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	// "Compatible" means that it is EVM compatible(the opcode and precompiled contracts are the same as Ethereum EVM).
	// In other words, not all the hard fork items are included.
	IstanbulCompatibleBlock  *big.Int `json:"istanbulCompatibleBlock,omitempty"`  // IstanbulCompatibleBlock switch block (nil = no fork, 0 = already on istanbul)
	LondonCompatibleBlock    *big.Int `json:"londonCompatibleBlock,omitempty"`    // LondonCompatibleBlock switch block (nil = no fork, 0 = already on london)
	EthTxTypeCompatibleBlock *big.Int `json:"ethTxTypeCompatibleBlock,omitempty"` // EthTxTypeCompatibleBlock switch block (nil = no fork, 0 = already on ethTxType)
	MagmaCompatibleBlock     *big.Int `json:"magmaCompatibleBlock,omitempty"`     // MagmaCompatible switch block (nil = no fork, 0 already on Magma)
	KoreCompatibleBlock      *big.Int `json:"koreCompatibleBlock,omitempty"`      // KoreCompatible switch block (nil = no fork, 0 already on Kore)
	ShanghaiCompatibleBlock  *big.Int `json:"shanghaiCompatibleBlock,omitempty"`  // ShanghaiCompatible switch block (nil = no fork, 0 already on shanghai)
	CancunCompatibleBlock    *big.Int `json:"cancunCompatibleBlock,omitempty"`    // CancunCompatible switch block (nil = no fork, 0 already on Cancun)

	// KIP103 is a special purpose hardfork feature that can be executed only once
	// Both Kip103CompatibleBlock and Kip103ContractAddress should be specified to enable KIP103
	Kip103CompatibleBlock *big.Int       `json:"kip103CompatibleBlock,omitempty"` // Kip103Compatible activate block (nil = no fork)
	Kip103ContractAddress common.Address `json:"kip103ContractAddress,omitempty"` // Kip103 contract address already deployed on the network

	// Randao is an optional hardfork
	// RandaoCompatibleBlock, RandaoRegistryRecords and RandaoRegistryOwner all must be specified to enable Randao
	RandaoCompatibleBlock *big.Int        `json:"randaoCompatibleBlock,omitempty"` // RandaoCompatible activate block (nil = no fork)
	RandaoRegistry        *RegistryConfig `json:"randaoRegistry,omitempty"`        // Registry initial states

	// Various consensus engines
	Gxhash   *GxhashConfig   `json:"gxhash,omitempty"` // (deprecated) not supported engine
	Clique   *CliqueConfig   `json:"clique,omitempty"`
	Istanbul *IstanbulConfig `json:"istanbul,omitempty"`

	UnitPrice     uint64            `json:"unitPrice"`
	DeriveShaImpl int               `json:"deriveShaImpl"`
	Governance    *GovernanceConfig `json:"governance"`
}

// GovernanceConfig stores governance information for a network
type GovernanceConfig struct {
	GoverningNode    common.Address `json:"governingNode"`
	GovernanceMode   string         `json:"governanceMode"`
	GovParamContract common.Address `json:"govParamContract"`
	Reward           *RewardConfig  `json:"reward,omitempty"`
	KIP71            *KIP71Config   `json:"kip71,omitempty"`
}

func (g *GovernanceConfig) DeferredTxFee() bool {
	return g.Reward.DeferredTxFee
}

// RewardConfig stores information about the network's token economy
type RewardConfig struct {
	MintingAmount          *big.Int `json:"mintingAmount"`
	Ratio                  string   `json:"ratio"`                  // Define how much portion of reward be distributed to CN/KFF/KCF
	Kip82Ratio             string   `json:"kip82ratio,omitempty"`   // Define how much portion of reward be distributed to proposer/stakers
	UseGiniCoeff           bool     `json:"useGiniCoeff"`           // Decide if Gini Coefficient will be used or not
	DeferredTxFee          bool     `json:"deferredTxFee"`          // Decide if TX fee will be handled instantly or handled later at block finalization
	StakingUpdateInterval  uint64   `json:"stakingUpdateInterval"`  // Interval when staking information is updated
	ProposerUpdateInterval uint64   `json:"proposerUpdateInterval"` // Interval when proposer information is updated
	MinimumStake           *big.Int `json:"minimumStake"`           // Minimum amount of peb to join CCO
}

// Magma governance parameters
type KIP71Config struct {
	LowerBoundBaseFee         uint64 `json:"lowerboundbasefee"`         // Minimum base fee for dynamic gas price
	UpperBoundBaseFee         uint64 `json:"upperboundbasefee"`         // Maximum base fee for dynamic gas price
	GasTarget                 uint64 `json:"gastarget"`                 // Gauge parameter increasing or decreasing gas price
	MaxBlockGasUsedForBaseFee uint64 `json:"maxblockgasusedforbasefee"` // Maximum network and process capacity to allow in a block
	BaseFeeDenominator        uint64 `json:"basefeedenominator"`        // For normalizing effect of the rapid change like impulse gas used
}

// IstanbulConfig is the consensus engine configs for Istanbul based sealing.
type IstanbulConfig struct {
	Epoch          uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
	ProposerPolicy uint64 `json:"policy"` // The policy for proposer selection; 0: Round Robin, 1: Sticky, 2: Weighted Random
	SubGroupSize   uint64 `json:"sub"`
}

// RegistryConfig is the initial KIP-149 system contract registry states.
// It is installed at block (RandaoCompatibleBlock - 1). Initial states are not applied if RandaoCompatibleBlock is nil or 0.
// To install the initial states from the block 0, use the AllocRegistry to generate GenesisAlloc.
//
// This struct only represents a special case of Registry state where:
// - there is only one record for each name
// - the activation of all records is zero
// - the names array is lexicographically sorted
type RegistryConfig struct {
	Records map[string]common.Address `json:"records"`
	Owner   common.Address            `json:"owner"`
}

// GxhashConfig is the consensus engine configs for proof-of-work based sealing.
// Deprecated: Use IstanbulConfig or CliqueConfig.
type GxhashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *GxhashConfig) String() string {
	return "gxhash"
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// String implements the stringer interface, returning the consensus engine details.
func (c *IstanbulConfig) String() string {
	return "istanbul"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Gxhash != nil:
		engine = c.Gxhash
	case c.Clique != nil:
		engine = c.Clique
	case c.Istanbul != nil:
		engine = c.Istanbul
	default:
		engine = "unknown"
	}

	kip103 := fmt.Sprintf("KIP103CompatibleBlock: %v KIP103ContractAddress %s", c.Kip103CompatibleBlock, c.Kip103ContractAddress.String())

	if c.Istanbul != nil {
		return fmt.Sprintf("{ChainID: %v IstanbulCompatibleBlock: %v LondonCompatibleBlock: %v EthTxTypeCompatibleBlock: %v MagmaCompatibleBlock: %v KoreCompatibleBlock: %v ShanghaiCompatibleBlock: %v CancunCompatibleBlock: %v RandaoCompatibleBlock: %v %s SubGroupSize: %d UnitPrice: %d DeriveShaImpl: %d Engine: %v}",
			c.ChainID,
			c.IstanbulCompatibleBlock,
			c.LondonCompatibleBlock,
			c.EthTxTypeCompatibleBlock,
			c.MagmaCompatibleBlock,
			c.KoreCompatibleBlock,
			c.ShanghaiCompatibleBlock,
			c.CancunCompatibleBlock,
			c.RandaoCompatibleBlock,
			kip103,
			c.Istanbul.SubGroupSize,
			c.UnitPrice,
			c.DeriveShaImpl,
			engine,
		)
	} else {
		return fmt.Sprintf("{ChainID: %v IstanbulCompatibleBlock: %v LondonCompatibleBlock: %v EthTxTypeCompatibleBlock: %v MagmaCompatibleBlock: %v KoreCompatibleBlock: %v ShanghaiCompatibleBlock: %v CancunCompatibleBlock: %v RandaoCompatibleBlock: %v %s UnitPrice: %d DeriveShaImpl: %d Engine: %v }",
			c.ChainID,
			c.IstanbulCompatibleBlock,
			c.LondonCompatibleBlock,
			c.EthTxTypeCompatibleBlock,
			c.MagmaCompatibleBlock,
			c.KoreCompatibleBlock,
			c.ShanghaiCompatibleBlock,
			c.CancunCompatibleBlock,
			c.RandaoCompatibleBlock,
			kip103,
			c.UnitPrice,
			c.DeriveShaImpl,
			engine,
		)
	}
}

func (c *ChainConfig) Copy() *ChainConfig {
	r := &ChainConfig{}
	j, _ := json.Marshal(c)
	json.Unmarshal(j, r)
	return r
}

// IsIstanbulForkEnabled returns whether num is either equal to the istanbul block or greater.
func (c *ChainConfig) IsIstanbulForkEnabled(num *big.Int) bool {
	return isForked(c.IstanbulCompatibleBlock, num)
}

// IsLondonForkEnabled returns whether num is either equal to the london block or greater.
func (c *ChainConfig) IsLondonForkEnabled(num *big.Int) bool {
	return isForked(c.LondonCompatibleBlock, num)
}

// IsEthTxTypeForkEnabled returns whether num is either equal to the ethTxType block or greater.
func (c *ChainConfig) IsEthTxTypeForkEnabled(num *big.Int) bool {
	return isForked(c.EthTxTypeCompatibleBlock, num)
}

// IsMagmaForkEnabled returns whether num is either equal to the magma block or greater.
func (c *ChainConfig) IsMagmaForkEnabled(num *big.Int) bool {
	return isForked(c.MagmaCompatibleBlock, num)
}

// IsKoreForkEnabled returns whether num is either equal to the kore block or greater.
func (c *ChainConfig) IsKoreForkEnabled(num *big.Int) bool {
	return isForked(c.KoreCompatibleBlock, num)
}

// IsShanghaiForkEnabled returns whether num is either equal to the shanghai block or greater.
func (c *ChainConfig) IsShanghaiForkEnabled(num *big.Int) bool {
	return isForked(c.ShanghaiCompatibleBlock, num)
}

// IsCancunForkEnabled returns whether num is either equal to the cancun block or greater.
func (c *ChainConfig) IsCancunForkEnabled(num *big.Int) bool {
	return isForked(c.CancunCompatibleBlock, num)
}

// IsRandaoForkEnabled returns whether num is either equal to the randao block or greater.
func (c *ChainConfig) IsRandaoForkEnabled(num *big.Int) bool {
	return isForked(c.RandaoCompatibleBlock, num)
}

// IsKIP103ForkBlock returns whether num is equal to the kip103 block.
func (c *ChainConfig) IsKIP103ForkBlock(num *big.Int) bool {
	if c.Kip103CompatibleBlock == nil || num == nil {
		return false
	}
	return c.Kip103CompatibleBlock.Cmp(num) == 0
}

// IsRandaoForkBlockParent returns whethere num is one block before the randao block.
func (c *ChainConfig) IsRandaoForkBlockParent(num *big.Int) bool {
	if c.RandaoCompatibleBlock == nil || num == nil {
		return false
	}
	nextNum := new(big.Int).Add(num, common.Big1)
	return c.RandaoCompatibleBlock.Cmp(nextNum) == 0 // randao == num + 1
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name     string
		block    *big.Int
		optional bool // if true, the fork may be nil and next fork is still allowed
	}
	var lastFork fork
	for _, cur := range []fork{
		{name: "istanbulBlock", block: c.IstanbulCompatibleBlock},
		{name: "londonBlock", block: c.LondonCompatibleBlock},
		{name: "ethTxTypeBlock", block: c.EthTxTypeCompatibleBlock},
		{name: "magmaBlock", block: c.MagmaCompatibleBlock},
		{name: "koreBlock", block: c.KoreCompatibleBlock},
		{name: "shanghaiBlock", block: c.ShanghaiCompatibleBlock},
		{name: "cancunBlock", block: c.CancunCompatibleBlock},
		{name: "randaoBlock", block: c.RandaoCompatibleBlock, optional: true},
	} {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || cur.block != nil {
			lastFork = cur
		}
	}
	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.IstanbulCompatibleBlock, newcfg.IstanbulCompatibleBlock, head) {
		return newCompatError("Istanbul Block", c.IstanbulCompatibleBlock, newcfg.IstanbulCompatibleBlock)
	}
	if isForkIncompatible(c.LondonCompatibleBlock, newcfg.LondonCompatibleBlock, head) {
		return newCompatError("London Block", c.LondonCompatibleBlock, newcfg.LondonCompatibleBlock)
	}
	if isForkIncompatible(c.EthTxTypeCompatibleBlock, newcfg.EthTxTypeCompatibleBlock, head) {
		return newCompatError("EthTxType Block", c.EthTxTypeCompatibleBlock, newcfg.EthTxTypeCompatibleBlock)
	}
	if isForkIncompatible(c.MagmaCompatibleBlock, newcfg.MagmaCompatibleBlock, head) {
		return newCompatError("Magma Block", c.MagmaCompatibleBlock, newcfg.MagmaCompatibleBlock)
	}
	if isForkIncompatible(c.KoreCompatibleBlock, newcfg.KoreCompatibleBlock, head) {
		return newCompatError("Kore Block", c.KoreCompatibleBlock, newcfg.KoreCompatibleBlock)
	}
	// We have intentionally skipped kip103Block in the fork ordering check since kip103 is designed
	// as an optional hardfork and there are no dependency with other forks.
	if isForkIncompatible(c.ShanghaiCompatibleBlock, newcfg.ShanghaiCompatibleBlock, head) {
		return newCompatError("Shanghai Block", c.ShanghaiCompatibleBlock, newcfg.ShanghaiCompatibleBlock)
	}
	if isForkIncompatible(c.CancunCompatibleBlock, newcfg.CancunCompatibleBlock, head) {
		return newCompatError("Cancun Block", c.CancunCompatibleBlock, newcfg.CancunCompatibleBlock)
	}
	if isForkIncompatible(c.RandaoCompatibleBlock, newcfg.RandaoCompatibleBlock, head) {
		return newCompatError("Randao Block", c.RandaoCompatibleBlock, newcfg.RandaoCompatibleBlock)
	}
	return nil
}

// SetDefaultsForGenesis fills undefined chain config with default values.
// Only used for generating genesis.
// Empty values from genesis.json will be left out from genesis.
func (c *ChainConfig) SetDefaultsForGenesis() {
	if c.Clique == nil && c.Istanbul == nil {
		c.Istanbul = GetDefaultIstanbulConfig()
		logger.Warn("Override the default Istanbul config to the chain config")
	}

	if c.Governance == nil {
		c.Governance = GetDefaultGovernanceConfigForGenesis()
		logger.Warn("Override the default governance config to the chain config")
	}

	if c.Governance.Reward == nil {
		c.Governance.Reward = GetDefaultRewardConfigForGenesis()
		logger.Warn("Override the default governance reward config to the chain config", "reward",
			c.Governance.Reward)
	}

	// StakingUpdateInterval must be nonzero because it is used as denominator
	if c.Governance.Reward.StakingUpdateInterval == 0 {
		c.Governance.Reward.StakingUpdateInterval = StakingUpdateInterval()
		logger.Warn("Override the default staking update interval to the chain config", "interval",
			c.Governance.Reward.StakingUpdateInterval)
	}

	// ProposerUpdateInterval must be nonzero because it is used as denominator
	if c.Governance.Reward.ProposerUpdateInterval == 0 {
		c.Governance.Reward.ProposerUpdateInterval = ProposerUpdateInterval()
		logger.Warn("Override the default proposer update interval to the chain config", "interval",
			c.Governance.Reward.ProposerUpdateInterval)
	}
}

// SetDefaults fills undefined chain config with default values
// so that nil pointer does not exist in the chain config
func (c *ChainConfig) SetDefaults() {
	c.SetDefaultsForGenesis()

	if c.Governance.KIP71 == nil {
		c.Governance.KIP71 = GetDefaultKIP71Config()
	}
	if c.Governance.Reward.Kip82Ratio == "" {
		c.Governance.Reward.Kip82Ratio = DefaultKip82Ratio
	}
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID     *big.Int
	IsIstanbul  bool
	IsLondon    bool
	IsEthTxType bool
	IsMagma     bool
	IsKore      bool
	IsShanghai  bool
	IsCancun    bool
	IsRandao    bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:     new(big.Int).Set(chainID),
		IsIstanbul:  c.IsIstanbulForkEnabled(num),
		IsLondon:    c.IsLondonForkEnabled(num),
		IsEthTxType: c.IsEthTxTypeForkEnabled(num),
		IsMagma:     c.IsMagmaForkEnabled(num),
		IsKore:      c.IsKoreForkEnabled(num),
		IsShanghai:  c.IsShanghaiForkEnabled(num),
		IsCancun:    c.IsCancunForkEnabled(num),
		IsRandao:    c.IsRandaoForkEnabled(num),
	}
}

// cypress genesis config
func GetDefaultGovernanceConfigForGenesis() *GovernanceConfig {
	gov := &GovernanceConfig{
		GovernanceMode: DefaultGovernanceMode,
		GoverningNode:  common.HexToAddress(DefaultGoverningNode),
		Reward:         GetDefaultRewardConfigForGenesis(),
	}
	return gov
}

func GetDefaultGovernanceConfig() *GovernanceConfig {
	gov := &GovernanceConfig{
		GovernanceMode:   DefaultGovernanceMode,
		GoverningNode:    common.HexToAddress(DefaultGoverningNode),
		GovParamContract: common.HexToAddress(DefaultGovParamContract),
		Reward:           GetDefaultRewardConfig(),
		KIP71:            GetDefaultKIP71Config(),
	}
	return gov
}

func GetDefaultIstanbulConfig() *IstanbulConfig {
	return &IstanbulConfig{
		Epoch:          DefaultEpoch,
		ProposerPolicy: DefaultProposerPolicy,
		SubGroupSize:   DefaultSubGroupSize,
	}
}

func GetDefaultRewardConfigForGenesis() *RewardConfig {
	return &RewardConfig{
		MintingAmount:          DefaultMintingAmount,
		Ratio:                  DefaultRatio,
		UseGiniCoeff:           DefaultUseGiniCoeff,
		DeferredTxFee:          DefaultDeferredTxFee,
		StakingUpdateInterval:  DefaultStakeUpdateInterval,
		ProposerUpdateInterval: DefaultProposerRefreshInterval,
		MinimumStake:           DefaultMinimumStake,
	}
}

func GetDefaultRewardConfig() *RewardConfig {
	return &RewardConfig{
		MintingAmount:          DefaultMintingAmount,
		Ratio:                  DefaultRatio,
		Kip82Ratio:             DefaultKip82Ratio,
		UseGiniCoeff:           DefaultUseGiniCoeff,
		DeferredTxFee:          DefaultDeferredTxFee,
		StakingUpdateInterval:  DefaultStakeUpdateInterval,
		ProposerUpdateInterval: DefaultProposerRefreshInterval,
		MinimumStake:           DefaultMinimumStake,
	}
}

func GetDefaultKIP71Config() *KIP71Config {
	return &KIP71Config{
		LowerBoundBaseFee:         DefaultLowerBoundBaseFee,
		UpperBoundBaseFee:         DefaultUpperBoundBaseFee,
		GasTarget:                 DefaultGasTarget,
		MaxBlockGasUsedForBaseFee: DefaultMaxBlockGasUsedForBaseFee,
		BaseFeeDenominator:        DefaultBaseFeeDenominator,
	}
}

func GetDefaultCliqueConfig() *CliqueConfig {
	return &CliqueConfig{
		Epoch:  DefaultEpoch,
		Period: DefaultPeriod,
	}
}
