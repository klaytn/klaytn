package governance

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

// Mixed engine consists of multiple governance engines
//
// Each parameter is added to a parameter set from one of the following sources:
// The highest priority is 1, and falls back to lower ones if non-existent
//  1. contractParams: ContractEngine items (when enabled)
//  2. headerParams:   Header Governance items
//  3. initialParams:  initial ChainConfig from genesis.json
//  4. defaultParams:  Default params such as params.Default*
//                     Note that some items are not backed by defaultParams.
//
type MixedEngine struct {
	initialConfig *params.ChainConfig

	initialParams *params.GovParamSet // initial ChainConfig
	defaultParams *params.GovParamSet // constants used as last fallback

	currentParams *params.GovParamSet // latest params to be returned by Params()

	db database.DBManager

	// Subordinate engines
	// TODO: Add ContractEngine
	headerGov *Governance
}

// newMixedEngine instantiate a new MixedEngine struct.
// Only if doInit is true, subordinate engines will be initialized.
func newMixedEngine(config *params.ChainConfig, db database.DBManager, doInit bool) *MixedEngine {
	e := &MixedEngine{
		initialConfig: config,
		db:            db,
	}

	if p, err := params.NewGovParamSetChainConfig(config); err == nil {
		e.initialParams = p
		e.currentParams = p
	} else {
		logger.Crit("Error parsing initial ChainConfig", "err", err)
	}

	constMap := map[int]interface{}{
		params.LowerBoundBaseFee:         params.DefaultLowerBoundBaseFee,
		params.UpperBoundBaseFee:         params.DefaultUpperBoundBaseFee,
		params.GasTarget:                 params.DefaultGasTarget,
		params.MaxBlockGasUsedForBaseFee: params.DefaultMaxBlockGasUsedForBaseFee,
		params.BaseFeeDenominator:        params.DefaultBaseFeeDenominator,
	}
	if p, err := params.NewGovParamSetIntMap(constMap); err == nil {
		e.defaultParams = p
	} else {
		logger.Crit("Error parsing initial ParamSet", "err", err)
	}

	// Setup subordinate engines
	if doInit {
		e.headerGov = NewGovernanceInitialize(config, db)
	} else {
		e.headerGov = NewGovernance(config, db)
	}

	// Load last state
	e.UpdateParams()

	return e
}

// Developers are encouraged to call this constructor in most cases.
func NewMixedEngine(config *params.ChainConfig, db database.DBManager) *MixedEngine {
	return newMixedEngine(config, db, true)
}

// Does not load initial data for test purposes
func NewMixedEngineNoInit(config *params.ChainConfig, db database.DBManager) *MixedEngine {
	return newMixedEngine(config, db, false)
}

func (e *MixedEngine) Params() *params.GovParamSet {
	headerParams := e.headerGov.Params()
	return e.assembleParams(headerParams)
}

func (e *MixedEngine) ParamsAt(num uint64) (*params.GovParamSet, error) {
	headerParams, err := e.headerGov.ParamsAt(num)
	if err != nil {
		return nil, err
	}

	// TODO-Klaytn-Kore: merge contractParams
	return e.assembleParams(headerParams), nil
}

func (e *MixedEngine) UpdateParams() error {
	if err := e.headerGov.UpdateParams(); err != nil {
		return err
	}

	headerParams := e.headerGov.Params()

	// TODO-Klaytn-Kore: merge contractParams
	e.currentParams = e.assembleParams(headerParams)
	return nil
}

func (e *MixedEngine) assembleParams(headerParams *params.GovParamSet) *params.GovParamSet {
	// Refer to the comments above `type MixedEngine` for assembly order
	p := params.NewGovParamSet()
	p = params.NewGovParamSetMerged(p, e.defaultParams)
	p = params.NewGovParamSetMerged(p, e.initialParams)
	p = params.NewGovParamSetMerged(p, headerParams)
	return p
}

// Pass-through to HeaderEngine
func (e *MixedEngine) AddVote(key string, val interface{}) bool {
	return e.headerGov.AddVote(key, val)
}

func (e *MixedEngine) ValidateVote(vote *GovernanceVote) (*GovernanceVote, bool) {
	return e.headerGov.ValidateVote(vote)
}

func (e *MixedEngine) CanWriteGovernanceState(num uint64) bool {
	return e.headerGov.CanWriteGovernanceState(num)
}

func (e *MixedEngine) WriteGovernanceState(num uint64, isCheckpoint bool) error {
	return e.headerGov.WriteGovernanceState(num, isCheckpoint)
}

func (e *MixedEngine) ReadGovernance(num uint64) (uint64, map[string]interface{}, error) {
	return e.headerGov.ReadGovernance(num)
}

func (e *MixedEngine) WriteGovernance(num uint64, data GovernanceSet, delta GovernanceSet) error {
	return e.headerGov.WriteGovernance(num, data, delta)
}

func (e *MixedEngine) GetEncodedVote(addr common.Address, number uint64) []byte {
	return e.headerGov.GetEncodedVote(addr, number)
}

func (e *MixedEngine) GetGovernanceChange() map[string]interface{} {
	return e.headerGov.GetGovernanceChange()
}

func (e *MixedEngine) VerifyGovernance(received []byte) error {
	return e.headerGov.VerifyGovernance(received)
}

func (e *MixedEngine) ClearVotes(num uint64) {
	e.headerGov.ClearVotes(num)
}

func (e *MixedEngine) WriteGovernanceForNextEpoch(number uint64, governance []byte) {
	e.headerGov.WriteGovernanceForNextEpoch(number, governance)
}

func (e *MixedEngine) UpdateCurrentSet(num uint64) {
	e.headerGov.UpdateCurrentSet(num)
}

func (e *MixedEngine) HandleGovernanceVote(
	valset istanbul.ValidatorSet, votes []GovernanceVote, tally []GovernanceTallyItem,
	header *types.Header, proposer common.Address, self common.Address, writable bool,
) (
	istanbul.ValidatorSet, []GovernanceVote, []GovernanceTallyItem,
) {
	return e.headerGov.HandleGovernanceVote(valset, votes, tally, header, proposer, self, writable)
}

func (e *MixedEngine) ChainId() uint64 {
	return e.headerGov.ChainId()
}

func (e *MixedEngine) InitialChainConfig() *params.ChainConfig {
	return e.headerGov.InitialChainConfig()
}

func (e *MixedEngine) GetVoteMapCopy() map[string]VoteStatus {
	return e.headerGov.GetVoteMapCopy()
}

func (e *MixedEngine) GetGovernanceTalliesCopy() []GovernanceTallyItem {
	return e.headerGov.GetGovernanceTalliesCopy()
}

func (e *MixedEngine) CurrentSetCopy() map[string]interface{} {
	return e.headerGov.CurrentSetCopy()
}

func (e *MixedEngine) PendingChanges() map[string]interface{} {
	return e.headerGov.PendingChanges()
}

func (e *MixedEngine) Votes() []GovernanceVote {
	return e.headerGov.Votes()
}

func (e *MixedEngine) IdxCache() []uint64 {
	return e.headerGov.IdxCache()
}

func (e *MixedEngine) IdxCacheFromDb() []uint64 {
	return e.headerGov.IdxCacheFromDb()
}

func (e *MixedEngine) NodeAddress() common.Address {
	return e.headerGov.NodeAddress()
}

func (e *MixedEngine) TotalVotingPower() uint64 {
	return e.headerGov.TotalVotingPower()
}

func (e *MixedEngine) MyVotingPower() uint64 {
	return e.headerGov.MyVotingPower()
}

func (e *MixedEngine) BlockChain() blockChain {
	return e.headerGov.BlockChain()
}

func (e *MixedEngine) DB() database.DBManager {
	return e.headerGov.DB()
}

func (e *MixedEngine) SetNodeAddress(addr common.Address) {
	e.headerGov.SetNodeAddress(addr)
}

func (e *MixedEngine) SetTotalVotingPower(t uint64) {
	e.headerGov.SetTotalVotingPower(t)
}

func (e *MixedEngine) SetMyVotingPower(t uint64) {
	e.headerGov.SetMyVotingPower(t)
}

func (e *MixedEngine) SetBlockchain(chain blockChain) {
	e.headerGov.SetBlockchain(chain)
}

func (e *MixedEngine) SetTxPool(txpool txPool) {
	e.headerGov.SetTxPool(txpool)
}

func (e *MixedEngine) GetTxPool() txPool {
	return e.headerGov.GetTxPool()
}
