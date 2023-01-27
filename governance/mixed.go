// Copyright 2022 The klaytn Authors
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

package governance

import (
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

// MixedEngine consists of multiple governance engines
//
// Each parameter is added to a parameter set from one of the following sources:
// The highest priority is 1, and falls back to lower ones if non-existent
//  1. contractParams: ContractEngine items (when enabled)
//  2. headerParams:   Header Governance items
//  3. initialParams:  initial ChainConfig from genesis.json
//  4. defaultParams:  Default params such as params.Default*
//     Note that some items are not backed by defaultParams.
type MixedEngine struct {
	// The same ChainConfig instance as Blockchain.chainConfig, {cn, worker}.config
	config *params.ChainConfig

	initialParams *params.GovParamSet // initial ChainConfig
	defaultParams *params.GovParamSet // default constants used as last fallback

	currentParams *params.GovParamSet // latest params to be returned by Params()

	db database.DBManager

	// Subordinate engines
	// contractGov is enabled when all the following conditions are met:
	//   - Kore hardfork block has passed
	//   - GovParamContract has been set
	// contractGov can be ignored even if it is enabled for various reasons. To name a few:
	//   - GovParam returns invalid parameters
	//   - Calling GovParam reverted
	// contractGov can throw critical error:
	//   - headerGov.BlockChain() is nil -> Fix by calling gov.SetBlockChain(bc)
	// headerGov cannot be disabled. However, its parameters can be ignored
	// by the prior contract parameters
	contractGov *ContractEngine
	headerGov   *Governance

	// for param update
	txpool     txPool
	blockchain blockChain
}

// newMixedEngine instantiate a new MixedEngine struct.
// Only if doInit is true, subordinate engines will be initialized.
func newMixedEngine(config *params.ChainConfig, db database.DBManager, doInit bool) *MixedEngine {
	e := &MixedEngine{
		config: config,
		db:     db,
	}

	if p, err := params.NewGovParamSetChainConfig(config); err == nil {
		e.initialParams = p
		e.currentParams = p
	} else {
		logger.Crit("Error parsing initial ChainConfig", "err", err)
	}

	defaultMap := map[int]interface{}{
		params.LowerBoundBaseFee:         params.DefaultLowerBoundBaseFee,
		params.UpperBoundBaseFee:         params.DefaultUpperBoundBaseFee,
		params.GasTarget:                 params.DefaultGasTarget,
		params.MaxBlockGasUsedForBaseFee: params.DefaultMaxBlockGasUsedForBaseFee,
		params.BaseFeeDenominator:        params.DefaultBaseFeeDenominator,
		params.GovParamContract:          params.DefaultGovParamContract,
		params.Kip82Ratio:                params.DefaultKip82Ratio,
	}
	if p, err := params.NewGovParamSetIntMap(defaultMap); err == nil {
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

	e.contractGov = NewContractEngine(e.headerGov)

	return e
}

// NewMixedEngine creates a governance engine using both contract-based and haeder-based gov.
// Developers are encouraged to call this constructor in most cases.
func NewMixedEngine(config *params.ChainConfig, db database.DBManager) *MixedEngine {
	return newMixedEngine(config, db, true)
}

// NewMixedEngineNoInit creates a MixedEngine without initializing governance.
func NewMixedEngineNoInit(config *params.ChainConfig, db database.DBManager) *MixedEngine {
	return newMixedEngine(config, db, false)
}

func (e *MixedEngine) Params() *params.GovParamSet {
	return e.currentParams
}

func (e *MixedEngine) ParamsAt(num uint64) (*params.GovParamSet, error) {
	var contractParams *params.GovParamSet
	var err error

	if e.config.IsKoreForkEnabled(new(big.Int).SetUint64(num)) {
		contractParams, err = e.contractGov.ParamsAt(num)
		if err != nil {
			logger.Error("contractGov.ParamsAt() failed", "err", err)
			return nil, err
		}
	} else {
		contractParams = params.NewGovParamSet()
	}

	headerParams, err := e.headerGov.ParamsAt(num)
	if err != nil {
		logger.Error("headerGov.ParamsAt() failed", "err", err)
		return nil, err
	}

	return e.assembleParams(headerParams, contractParams), nil
}

func (e *MixedEngine) UpdateParams(num uint64) error {
	var contractParams *params.GovParamSet
	numBigInt := big.NewInt(int64(num))

	if e.config.IsKoreForkEnabled(numBigInt) {
		if err := e.contractGov.UpdateParams(num); err != nil {
			logger.Error("contractGov.UpdateParams(num) failed", "num", num, "err", err)
			return err
		}
		contractParams = e.contractGov.Params()
	} else {
		contractParams = params.NewGovParamSet()
	}

	if err := e.headerGov.UpdateParams(num); err != nil {
		logger.Error("headerGov.UpdateParams(num) failed", "num", num, "err", err)
		return err
	}

	headerParams := e.headerGov.Params()

	newParams := e.assembleParams(headerParams, contractParams)
	e.handleParamUpdate(e.currentParams, newParams)

	e.currentParams = newParams

	return nil
}

func (e *MixedEngine) assembleParams(headerParams, contractParams *params.GovParamSet) *params.GovParamSet {
	// Refer to the comments above `type MixedEngine` for assembly order
	p := params.NewGovParamSet()
	p = params.NewGovParamSetMerged(p, e.defaultParams)
	p = params.NewGovParamSetMerged(p, e.initialParams)
	p = params.NewGovParamSetMerged(p, headerParams)
	p = params.NewGovParamSetMerged(p, contractParams)
	return p
}

func (e *MixedEngine) handleParamUpdate(old, new *params.GovParamSet) {
	// NOTE: key set must be the same, which is guaranteed at NewMixedEngine
	for k, oldval := range old.IntMap() {
		if newval := new.MustGet(k); oldval != newval {
			switch k {
			// config.Istanbul
			case params.Epoch:
				e.config.Istanbul.Epoch = new.Epoch()
			case params.Policy:
				e.config.Istanbul.ProposerPolicy = new.Policy()
			case params.CommitteeSize:
				e.config.Istanbul.SubGroupSize = new.CommitteeSize()
			// config.Governance
			case params.GoverningNode:
				e.config.Governance.GoverningNode = new.GoverningNode()
			case params.GovernanceMode:
				e.config.Governance.GovernanceMode = new.GovernanceModeStr()
			case params.GovParamContract:
				e.config.Governance.GovParamContract = new.GovParamContract()
			// config.Governance.Reward
			case params.MintingAmount:
				e.config.Governance.Reward.MintingAmount = new.MintingAmountBig()
			case params.Ratio:
				e.config.Governance.Reward.Ratio = new.Ratio()
			case params.Kip82Ratio:
				e.config.Governance.Reward.Kip82Ratio = new.Kip82Ratio()
			case params.UseGiniCoeff:
				e.config.Governance.Reward.UseGiniCoeff = new.UseGiniCoeff()
			case params.DeferredTxFee:
				e.config.Governance.Reward.DeferredTxFee = new.DeferredTxFee()
			case params.MinimumStake:
				e.config.Governance.Reward.MinimumStake = new.MinimumStakeBig()
			case params.StakeUpdateInterval:
				e.config.Governance.Reward.StakingUpdateInterval = new.StakeUpdateInterval()
				params.SetStakingUpdateInterval(new.StakeUpdateInterval())
			case params.ProposerRefreshInterval:
				e.config.Governance.Reward.ProposerUpdateInterval = new.ProposerRefreshInterval()
				params.SetProposerUpdateInterval(new.ProposerRefreshInterval())
			// config.Governance.KIP17
			case params.LowerBoundBaseFee:
				e.config.Governance.KIP71.LowerBoundBaseFee = new.LowerBoundBaseFee()
			case params.UpperBoundBaseFee:
				e.config.Governance.KIP71.UpperBoundBaseFee = new.UpperBoundBaseFee()
			case params.GasTarget:
				e.config.Governance.KIP71.GasTarget = new.GasTarget()
			case params.MaxBlockGasUsedForBaseFee:
				e.config.Governance.KIP71.MaxBlockGasUsedForBaseFee = new.MaxBlockGasUsedForBaseFee()
			case params.BaseFeeDenominator:
				e.config.Governance.KIP71.BaseFeeDenominator = new.BaseFeeDenominator()
			// others
			case params.UnitPrice:
				e.config.UnitPrice = new.UnitPrice()
				if e.txpool != nil {
					e.txpool.SetGasPrice(big.NewInt(0).SetUint64(new.UnitPrice()))
				}
			case params.DeriveShaImpl:
				e.config.DeriveShaImpl = new.DeriveShaImpl()
			}
		}
	}
}

func (e *MixedEngine) HeaderGov() HeaderEngine {
	return e.headerGov
}

func (e *MixedEngine) ContractGov() ReaderEngine {
	return e.contractGov
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
	e.blockchain = chain
	e.headerGov.SetBlockchain(chain)
}

func (e *MixedEngine) SetTxPool(txpool txPool) {
	e.txpool = txpool
	e.headerGov.SetTxPool(txpool)
}

func (e *MixedEngine) GetTxPool() txPool {
	return e.headerGov.GetTxPool()
}
