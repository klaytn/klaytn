package governance

import (
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

// TODO: remove initialConfig, initialParams, currentParams if not necessary
type MixedEngine struct {
	initialConfig *params.ChainConfig
	initialParams *params.GovParamSet
	currentParams *params.GovParamSet

	db    database.DBManager
	chain blockChain

	// Subordinate engines
	contractGov *ContractEngine
	defaultGov  *Governance
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

	// Setup subordinate engines
	if doInit {
		e.defaultGov = NewGovernanceInitialize(config, db)
	} else {
		e.defaultGov = NewGovernance(config, db)
	}

	e.contractGov = NewContractEngine(config, e.defaultGov)

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
	if e.isContractEnabledAtNext() {
		return e.contractGov.Params()
	} else {
		return e.defaultGov.Params()
	}
}

func (e *MixedEngine) ParamsAt(num uint64) (*params.GovParamSet, error) {
	if e.isContractEnabledAt(num) {
		return e.contractGov.ParamsAt(num)
	} else {
		return e.defaultGov.ParamsAt(num)
	}
}

func (e *MixedEngine) UpdateParams() error {
	if e.isContractEnabledAtNext() {
		return e.contractGov.UpdateParams()
	} else {
		return e.defaultGov.UpdateParams()
	}
}

func (e *MixedEngine) SetBlockchain(chain blockChain) {
	e.chain = chain
	e.contractGov.SetBlockchain(chain)
	e.defaultGov.SetBlockchain(chain)
}

func (e *MixedEngine) isContractEnabledAt(num uint64) bool {
	return e.initialConfig.IsContractGovForkEnabled(new(big.Int).SetUint64(num))
}

func (e *MixedEngine) isContractEnabledAtNext() bool {
	if e.chain == nil {
		return false
	}
	head := e.chain.CurrentHeader().Number.Uint64()
	return e.isContractEnabledAt(head + 1)
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
	header *types.Header, proposer common.Address, self common.Address,
) (
	istanbul.ValidatorSet, []GovernanceVote, []GovernanceTallyItem,
) {
	return e.defaultGov.HandleGovernanceVote(valset, votes, tally, header, proposer, self)
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

func (e *MixedEngine) SetTxPool(txpool txPool) {
	e.defaultGov.SetTxPool(txpool)
}
