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
//  2. defaultParams:  Header Governance items
//  3. initialParams:  initial ChainConfig from genesis.json
//  4. constParams:    Constants such as params.Default*
//                     Note that some items are not backed by constParams.
//
type MixedEngine struct {
	initialConfig *params.ChainConfig

	initialParams *params.GovParamSet // initial ChainConfig
	constParams   *params.GovParamSet // constants used as last fallback

	currentParams *params.GovParamSet // latest params to be returned by Params()

	db database.DBManager

	// Subordinate engines
	// TODO: Add ContractEngine
	defaultGov *Governance
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
		e.constParams = p
	} else {
		logger.Crit("Error parsing initial ParamSet", "err", err)
	}

	// Setup subordinate engines
	if doInit {
		e.defaultGov = NewGovernanceInitialize(config, db)
	} else {
		e.defaultGov = NewGovernance(config, db)
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
	return e.defaultGov.Params()
}

func (e *MixedEngine) ParamsAt(num uint64) (*params.GovParamSet, error) {
	defaultParams, err := e.defaultGov.ParamsAt(num)
	if err != nil {
		return nil, err
	}

	// TODO-Klaytn-Kore: merge contractParams
	return e.assembleParams(defaultParams), nil
}

func (e *MixedEngine) UpdateParams() error {
	if err := e.defaultGov.UpdateParams(); err != nil {
		return err
	}

	defaultParams := e.defaultGov.Params()

	// TODO-Klaytn-Kore: merge contractParams
	e.currentParams = e.assembleParams(defaultParams)
	return nil
}

func (e *MixedEngine) assembleParams(defaultParams *params.GovParamSet) *params.GovParamSet {
	// Refer to the comments above `type MixedEngine` for assembly order
	p := params.NewGovParamSet()
	p = params.NewGovParamSetMerged(p, e.constParams)
	p = params.NewGovParamSetMerged(p, e.initialParams)
	p = params.NewGovParamSetMerged(p, defaultParams)
	return p
}

// Pass-through to HeaderEngine
func (e *MixedEngine) AddVote(key string, val interface{}) bool {
	return e.defaultGov.AddVote(key, val)
}

func (e *MixedEngine) ValidateVote(vote *GovernanceVote) (*GovernanceVote, bool) {
	return e.defaultGov.ValidateVote(vote)
}

func (e *MixedEngine) CanWriteGovernanceState(num uint64) bool {
	return e.defaultGov.CanWriteGovernanceState(num)
}

func (e *MixedEngine) WriteGovernanceState(num uint64, isCheckpoint bool) error {
	return e.defaultGov.WriteGovernanceState(num, isCheckpoint)
}

func (e *MixedEngine) ReadGovernance(num uint64) (uint64, map[string]interface{}, error) {
	return e.defaultGov.ReadGovernance(num)
}

func (e *MixedEngine) WriteGovernance(num uint64, data GovernanceSet, delta GovernanceSet) error {
	return e.defaultGov.WriteGovernance(num, data, delta)
}

func (e *MixedEngine) GetEncodedVote(addr common.Address, number uint64) []byte {
	return e.defaultGov.GetEncodedVote(addr, number)
}

func (e *MixedEngine) GetGovernanceChange() map[string]interface{} {
	return e.defaultGov.GetGovernanceChange()
}

func (e *MixedEngine) VerifyGovernance(received []byte) error {
	return e.defaultGov.VerifyGovernance(received)
}

func (e *MixedEngine) ClearVotes(num uint64) {
	e.defaultGov.ClearVotes(num)
}

func (e *MixedEngine) WriteGovernanceForNextEpoch(number uint64, governance []byte) {
	e.defaultGov.WriteGovernanceForNextEpoch(number, governance)
}

func (e *MixedEngine) UpdateCurrentSet(num uint64) {
	e.defaultGov.UpdateCurrentSet(num)
}

func (e *MixedEngine) HandleGovernanceVote(
	valset istanbul.ValidatorSet, votes []GovernanceVote, tally []GovernanceTallyItem,
	header *types.Header, proposer common.Address, self common.Address, writable bool,
) (
	istanbul.ValidatorSet, []GovernanceVote, []GovernanceTallyItem,
) {
	return e.defaultGov.HandleGovernanceVote(valset, votes, tally, header, proposer, self, writable)
}

func (e *MixedEngine) ChainId() uint64 {
	return e.defaultGov.ChainId()
}

func (e *MixedEngine) InitialChainConfig() *params.ChainConfig {
	return e.defaultGov.InitialChainConfig()
}

func (e *MixedEngine) GetVoteMapCopy() map[string]VoteStatus {
	return e.defaultGov.GetVoteMapCopy()
}

func (e *MixedEngine) GetGovernanceTalliesCopy() []GovernanceTallyItem {
	return e.defaultGov.GetGovernanceTalliesCopy()
}

func (e *MixedEngine) CurrentSetCopy() map[string]interface{} {
	return e.defaultGov.CurrentSetCopy()
}

func (e *MixedEngine) PendingChanges() map[string]interface{} {
	return e.defaultGov.PendingChanges()
}

func (e *MixedEngine) Votes() []GovernanceVote {
	return e.defaultGov.Votes()
}

func (e *MixedEngine) IdxCache() []uint64 {
	return e.defaultGov.IdxCache()
}

func (e *MixedEngine) IdxCacheFromDb() []uint64 {
	return e.defaultGov.IdxCacheFromDb()
}

func (e *MixedEngine) NodeAddress() common.Address {
	return e.defaultGov.NodeAddress()
}

func (e *MixedEngine) TotalVotingPower() uint64 {
	return e.defaultGov.TotalVotingPower()
}

func (e *MixedEngine) MyVotingPower() uint64 {
	return e.defaultGov.MyVotingPower()
}

func (e *MixedEngine) BlockChain() blockChain {
	return e.defaultGov.BlockChain()
}

func (e *MixedEngine) DB() database.DBManager {
	return e.defaultGov.DB()
}

func (e *MixedEngine) SetNodeAddress(addr common.Address) {
	e.defaultGov.SetNodeAddress(addr)
}

func (e *MixedEngine) SetTotalVotingPower(t uint64) {
	e.defaultGov.SetTotalVotingPower(t)
}

func (e *MixedEngine) SetMyVotingPower(t uint64) {
	e.defaultGov.SetMyVotingPower(t)
}

func (e *MixedEngine) SetBlockchain(chain blockChain) {
	e.defaultGov.SetBlockchain(chain)
}

func (e *MixedEngine) SetTxPool(txpool txPool) {
	e.defaultGov.SetTxPool(txpool)
}

func (e *MixedEngine) GetTxPool() txPool {
	return e.defaultGov.GetTxPool()
}

func (e *MixedEngine) GovernanceMode() string {
	return e.defaultGov.GovernanceMode()
}

func (e *MixedEngine) GoverningNode() common.Address {
	return e.defaultGov.GoverningNode()
}

func (e *MixedEngine) UnitPrice() uint64 {
	return e.defaultGov.UnitPrice()
}

func (e *MixedEngine) CommitteeSize() uint64 {
	return e.defaultGov.CommitteeSize()
}

func (e *MixedEngine) Epoch() uint64 {
	return e.defaultGov.Epoch()
}

func (e *MixedEngine) ProposerPolicy() uint64 {
	return e.defaultGov.ProposerPolicy()
}

func (e *MixedEngine) DeferredTxFee() bool {
	return e.defaultGov.DeferredTxFee()
}

func (e *MixedEngine) MinimumStake() string {
	return e.defaultGov.MinimumStake()
}

func (e *MixedEngine) MintingAmount() string {
	return e.defaultGov.MintingAmount()
}

func (e *MixedEngine) ProposerUpdateInterval() uint64 {
	return e.defaultGov.ProposerUpdateInterval()
}

func (e *MixedEngine) Ratio() string {
	return e.defaultGov.Ratio()
}

func (e *MixedEngine) StakingUpdateInterval() uint64 {
	return e.defaultGov.StakingUpdateInterval()
}

func (e *MixedEngine) UseGiniCoeff() bool {
	return e.defaultGov.UseGiniCoeff()
}

func (e *MixedEngine) LowerBoundBaseFee() uint64 {
	return e.defaultGov.LowerBoundBaseFee()
}

func (e *MixedEngine) UpperBoundBaseFee() uint64 {
	return e.defaultGov.UpperBoundBaseFee()
}

func (e *MixedEngine) GasTarget() uint64 {
	return e.defaultGov.GasTarget()
}

func (e *MixedEngine) MaxBlockGasUsedForBaseFee() uint64 {
	return e.defaultGov.MaxBlockGasUsedForBaseFee()
}

func (e *MixedEngine) BaseFeeDenominator() uint64 {
	return e.defaultGov.BaseFeeDenominator()
}

func (e *MixedEngine) GetGovernanceValue(key int) interface{} {
	return e.defaultGov.GetGovernanceValue(key)
}

func (e *MixedEngine) GetGovernanceItemAtNumber(num uint64, key string) (interface{}, error) {
	return e.defaultGov.GetGovernanceItemAtNumber(num, key)
}

func (e *MixedEngine) GetItemAtNumberByIntKey(num uint64, key int) (interface{}, error) {
	return e.defaultGov.GetItemAtNumberByIntKey(num, key)
}

func (e *MixedEngine) GetGoverningInfoAtNumber(num uint64) (bool, common.Address, error) {
	return e.defaultGov.GetGoverningInfoAtNumber(num)
}

func (e *MixedEngine) GetMinimumStakingAtNumber(num uint64) (uint64, error) {
	return e.defaultGov.GetMinimumStakingAtNumber(num)
}
