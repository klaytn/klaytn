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
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash      = common.HexToHash("// todo generate new hash for mainnet") // Mainnet genesis hash to enforce below configs on
	TestnetGenesisHash      = common.HexToHash("// todo generate new hash for testnet") // Testnet genesis hash to enforce below configs on
	AuthorAddressForTesting = common.HexToAddress("0xc0ea08a2d404d3172d2add29a45be56da40e2949")
)

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID: big.NewInt(1),
		Gxhash:  new(GxhashConfig),
	}

	// TestnetChainConfig contains the chain parameters to run a node on the Ropsten test network.
	TestnetChainConfig = &ChainConfig{
		ChainID: big.NewInt(2),
		Gxhash:  new(GxhashConfig),
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

var (
	// VMLogTarget sets the output target of vmlog.
	// The values below can be OR'ed.
	//  - 0x0: no output (default)
	//  - 0x1: file (DATADIR/logs/vm.log)
	//  - 0x2: stdout (like logger.DEBUG)
	VMLogTarget = 0x0
)

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
	GoverningNode  common.Address `json:"governingNode"`
	GovernanceMode string         `json:"governanceMode"`
	Reward         *RewardConfig  `json:"reward,omitempty"`
}

func (g *GovernanceConfig) DeferredTxFee() bool {
	return g.Reward.DeferredTxFee
}

// RewardConfig stores information about the network's token economy
type RewardConfig struct {
	MintingAmount          *big.Int `json:"mintingAmount"`
	Ratio                  string   `json:"ratio"`                  // Define how much portion of reward be distributed to CN/PoC/KIR
	UseGiniCoeff           bool     `json:"useGiniCoeff"`           // Decide if Gini Coefficient will be used or not
	DeferredTxFee          bool     `json:"deferredTxFee"`          // Decide if TX fee will be handled instantly or handled later at block finalization
	StakingUpdateInterval  uint64   `json:"stakingUpdateInterval"`  // Interval when staking information is updated
	ProposerUpdateInterval uint64   `json:"proposerUpdateInterval"` // Interval when proposer information is updated
	MinimumStake           *big.Int `json:"minimumStake"`           // Minimum amount of peb to join CCO
}

// IstanbulConfig is the consensus engine configs for Istanbul based sealing.
type IstanbulConfig struct {
	Epoch          uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
	ProposerPolicy uint64 `json:"policy"` // The policy for proposer selection
	SubGroupSize   uint64 `json:"sub"`
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
	if c.Istanbul != nil {
		return fmt.Sprintf("{ChainID: %v Engine: %v SubGroupSize: %d UnitPrice: %d DeriveShaImpl: %d}",
			c.ChainID,
			engine,
			c.Istanbul.SubGroupSize,
			c.UnitPrice,
			c.DeriveShaImpl,
		)
	} else {
		return fmt.Sprintf("{ChainID: %v Engine: %v UnitPrice: %d DeriveShaImpl: %d}",
			c.ChainID,
			engine,
			c.UnitPrice,
			c.DeriveShaImpl,
		)
	}
}

// GasTable returns the gas table corresponding to the current phase.
//
// The returned GasTable's fields shouldn't, under any circumstances, be changed.
func (c *ChainConfig) GasTable(num *big.Int) GasTable {
	return GasTableCypress
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

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	return nil
}

// GetConsensusEngine returns the consensus engine type specified in ChainConfig.
// It returns Unknown type if none of engine type is configured or more than one type is configured.
func (c *ChainConfig) GetConsensusEngine() EngineType {
	switch {
	case c.Clique != nil && c.Istanbul == nil:
		return UseClique
	case c.Clique == nil && c.Istanbul != nil:
		return UseIstanbul
	default:
		return Unknown
	}
}

// SetDefaults fills undefined chain config with default values.
func (c *ChainConfig) SetDefaults() {
	logger := log.NewModuleLogger(log.Governance)

	if c.GetConsensusEngine() == Unknown && c.Istanbul == nil {
		c.Istanbul = GetDefaultIstanbulConfig()
		logger.Warn("Override the default Istanbul config to the chain config")
	}

	if c.Governance == nil {
		engineType := c.GetConsensusEngine()
		c.Governance = GetDefaultGovernanceConfig(engineType)
		logger.Warn("Override the default governance config to the chain config", "engineType", engineType)
	}

	if c.Governance.Reward == nil {
		c.Governance.Reward = GetDefaultRewardConfig()
		logger.Warn("Override the default governance reward config to the chain config", "reward",
			c.Governance.Reward)
	}

	if c.Governance.Reward.StakingUpdateInterval == 0 {
		c.Governance.Reward.StakingUpdateInterval = StakingUpdateInterval()
		logger.Warn("Override the default staking update interval to the chain config", "interval",
			c.Governance.Reward.StakingUpdateInterval)
	}

	if c.Governance.Reward.ProposerUpdateInterval == 0 {
		c.Governance.Reward.ProposerUpdateInterval = ProposerUpdateInterval()
		logger.Warn("Override the default proposer update interval to the chain config", "interval",
			c.Governance.Reward.ProposerUpdateInterval)
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

// Rules wraps ChainConfig and is merely syntatic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID *big.Int
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{ChainID: new(big.Int).Set(chainID)}
}

// Copy copies self to a new governance config and return it
func (g *GovernanceConfig) Copy() *GovernanceConfig {
	newConfig := &GovernanceConfig{
		Reward: &RewardConfig{},
	}
	newConfig.GovernanceMode = g.GovernanceMode
	newConfig.Reward.MintingAmount = big.NewInt(0).Set(g.Reward.MintingAmount)
	newConfig.Reward.Ratio = g.Reward.Ratio
	newConfig.Reward.UseGiniCoeff = g.Reward.UseGiniCoeff
	newConfig.Reward.DeferredTxFee = g.Reward.DeferredTxFee
	newConfig.GoverningNode = g.GoverningNode

	return newConfig
}

func (c *IstanbulConfig) Copy() *IstanbulConfig {
	newIC := &IstanbulConfig{}

	newIC.Epoch = c.Epoch
	newIC.SubGroupSize = c.SubGroupSize
	newIC.ProposerPolicy = c.ProposerPolicy

	return newIC
}

// TODO-Klaytn-Governance: Remove input parameter if not needed anymore
func GetDefaultGovernanceConfig(engine EngineType) *GovernanceConfig {
	gov := &GovernanceConfig{
		GovernanceMode: DefaultGovernanceMode,
		GoverningNode:  common.HexToAddress(DefaultGoverningNode),
		Reward:         GetDefaultRewardConfig(),
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

func GetDefaultRewardConfig() *RewardConfig {
	return &RewardConfig{
		MintingAmount:          big.NewInt(DefaultMintingAmount),
		Ratio:                  DefaultRatio,
		UseGiniCoeff:           DefaultUseGiniCoeff,
		DeferredTxFee:          DefaultDefferedTxFee,
		StakingUpdateInterval:  uint64(86400),
		ProposerUpdateInterval: uint64(3600),
		MinimumStake:           big.NewInt(2000000),
	}
}

func GetDefaultCliqueConfig() *CliqueConfig {
	return &CliqueConfig{
		Epoch:  DefaultEpoch,
		Period: DefaultPeriod,
	}
}
