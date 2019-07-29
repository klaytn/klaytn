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

package governance

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
)

type check struct {
	t         reflect.Type
	validator func(k string, v interface{}) bool
	trigger   func(g *Governance, k string, v interface{}) bool
}

var (
	stringT  = reflect.TypeOf("")
	uint64T  = reflect.TypeOf(uint64(0))
	addressT = reflect.TypeOf(common.StringToAddress("0x0"))
	boolT    = reflect.TypeOf(true)
	float64T = reflect.TypeOf(float64(0.0))
)

var GovernanceItems = map[int]check{
	params.GovernanceMode:          {stringT, checkGovernanceMode, updateGovernanceConfig},
	params.GoverningNode:           {addressT, checkAddress, updateGovernanceConfig},
	params.UnitPrice:               {uint64T, checkUint64andBool, updateGovernanceConfig},
	params.AddValidator:            {addressT, checkAddress, updateGovernanceConfig},
	params.RemoveValidator:         {addressT, checkAddress, updateGovernanceConfig},
	params.MintingAmount:           {stringT, checkBigInt, updateGovernanceConfig},
	params.Ratio:                   {stringT, checkRatio, updateGovernanceConfig},
	params.UseGiniCoeff:            {boolT, checkUint64andBool, updateGovernanceConfig},
	params.DeferredTxFee:           {boolT, checkUint64andBool, updateGovernanceConfig},
	params.MinimumStake:            {stringT, checkBigInt, updateGovernanceConfig},
	params.StakeUpdateInterval:     {uint64T, checkUint64andBool, updateGovernanceConfig},
	params.ProposerRefreshInterval: {uint64T, checkUint64andBool, updateGovernanceConfig},
	params.Epoch:                   {uint64T, checkUint64andBool, updateGovernanceConfig},
	params.Policy:                  {uint64T, checkUint64andBool, updateGovernanceConfig},
	params.CommitteeSize:           {uint64T, checkUint64andBool, updateGovernanceConfig},
	params.ConstTxGasHumanReadable: {uint64T, checkUint64andBool, updateParams},
}

func updateParams(g *Governance, k string, v interface{}) bool {
	switch GovernanceKeyMap[k] {
	case params.ConstTxGasHumanReadable:
		params.TxGasHumanReadable = v.(uint64)
		logger.Info("TxGasHumanReadable changed", "New value", params.TxGasHumanReadable)
	}
	return true
}

func updateGovernanceConfig(g *Governance, k string, v interface{}) bool {
	switch GovernanceKeyMap[k] {
	case params.GovernanceMode:
		g.ChainConfig.Governance.GovernanceMode = v.(string)
	case params.GoverningNode:
		g.ChainConfig.Governance.GoverningNode = v.(common.Address)
	case params.UnitPrice:
		newPrice := v.(uint64)
		g.TxPool.SetGasPrice(big.NewInt(0).SetUint64(newPrice))
		g.ChainConfig.UnitPrice = newPrice
	case params.MintingAmount:
		g.ChainConfig.Governance.Reward.MintingAmount, _ = new(big.Int).SetString(v.(string), 10)
	case params.Ratio:
		g.ChainConfig.Governance.Reward.Ratio = v.(string)
	case params.UseGiniCoeff:
		g.ChainConfig.Governance.Reward.UseGiniCoeff = v.(bool)
		g.blockChain.Config().Governance.Reward.UseGiniCoeff = g.ChainConfig.Governance.Reward.UseGiniCoeff
	case params.DeferredTxFee:
		g.ChainConfig.Governance.Reward.DeferredTxFee = v.(bool)
	case params.MinimumStake:
		g.ChainConfig.Governance.Reward.MinimumStake, _ = new(big.Int).SetString(v.(string), 10)
	case params.StakeUpdateInterval:
		g.ChainConfig.Governance.Reward.StakingUpdateInterval = v.(uint64)
		params.SetStakingUpdateInterval(g.ChainConfig.Governance.Reward.StakingUpdateInterval)
	case params.ProposerRefreshInterval:
		g.ChainConfig.Governance.Reward.ProposerUpdateInterval = v.(uint64)
		params.SetProposerUpdateInterval(g.ChainConfig.Governance.Reward.ProposerUpdateInterval)
	case params.Epoch:
		g.ChainConfig.Istanbul.Epoch = v.(uint64)
	case params.Policy:
		g.ChainConfig.Istanbul.ProposerPolicy = uint64(v.(uint64))
		g.blockChain.Config().Istanbul.ProposerPolicy = g.ChainConfig.Istanbul.ProposerPolicy
	case params.CommitteeSize:
		g.ChainConfig.Istanbul.SubGroupSize = v.(uint64)
	}
	return true
}

// AddVote adds a vote to the voteMap
func (g *Governance) AddVote(key string, val interface{}) bool {
	key = g.getKey(key)

	// If the key is forbidden, stop processing it
	if _, ok := GovernanceForbiddenKeyMap[key]; ok {
		return false
	}

	vote := &GovernanceVote{Key: key, Value: val}
	var ok bool
	if vote, ok = g.ValidateVote(vote); ok {
		g.voteMap.SetValue(key, VoteStatus{
			Value:  vote.Value,
			Casted: false,
			Num:    0,
		})
		return true
	}
	return false
}

func (g *Governance) adjustValueType(key string, val interface{}) interface{} {
	k := GovernanceKeyMap[key]
	reqType := GovernanceItems[k].t
	var x interface{}

	// When an int value comes from JS console, it comes as a float64
	if reqType == uint64T && reflect.TypeOf(val) == float64T {
		x = uint64(val.(float64))
		if float64(x.(uint64)) == val.(float64) {
			return x
		}
		return val
	}

	// address comes as a form of string from JS console
	if reqType == addressT && reflect.TypeOf(val) == stringT {
		if common.IsHexAddress(val.(string)) {
			x = common.HexToAddress(val.(string))
			return x
		}
	}

	// If a string text come as uppercase, make it into lowercase
	if reflect.TypeOf(val) == stringT {
		x = strings.ToLower(val.(string))
		return x
	}
	return val
}

func (gov *Governance) checkType(vote *GovernanceVote) bool {
	key := GovernanceKeyMap[vote.Key]
	return GovernanceItems[key].t == reflect.TypeOf(vote.Value)
}

func (gov *Governance) checkKey(k string) bool {
	key := GovernanceKeyMap[k]
	if _, ok := GovernanceItems[key]; ok {
		return true
	}
	return false
}

func (gov *Governance) ValidateVote(vote *GovernanceVote) (*GovernanceVote, bool) {
	vote.Key = gov.getKey(vote.Key)
	key := GovernanceKeyMap[vote.Key]
	vote.Value = gov.adjustValueType(vote.Key, vote.Value)

	if gov.checkKey(vote.Key) && gov.checkType(vote) {
		return vote, GovernanceItems[key].validator(vote.Key, vote.Value)
	}
	return vote, false
}

func checkRatio(k string, v interface{}) bool {
	x := strings.Split(v.(string), "/")
	if len(x) != params.RewardSliceCount {
		return false
	}
	var sum uint64
	for _, item := range x {
		v, err := strconv.ParseUint(item, 10, 64)
		if err != nil {
			return false
		}
		sum += v
	}
	if sum == 100 {
		return true
	} else {
		return false
	}
}

func checkGovernanceMode(k string, v interface{}) bool {
	if _, ok := GovernanceModeMap[v.(string)]; ok {
		return true
	}
	return false
}

func checkUint64andBool(k string, v interface{}) bool {
	// for Uint64 and Bool, no more check is needed
	if reflect.TypeOf(v) == uint64T || reflect.TypeOf(v) == boolT {
		return true
	}
	return false
}

func checkProposerPolicy(k string, v interface{}) bool {
	if _, ok := ProposerPolicyMap[v.(string)]; ok {
		return true
	}
	return false
}

func checkBigInt(k string, v interface{}) bool {
	x := new(big.Int)
	if _, ok := x.SetString(v.(string), 10); ok {
		return true
	}
	return false
}

func checkAddress(k string, v interface{}) bool {
	return true
}

func (gov *Governance) HandleGovernanceVote(valset istanbul.ValidatorSet, votes []GovernanceVote, tally []GovernanceTallyItem, header *types.Header, proposer common.Address, self common.Address) (istanbul.ValidatorSet, []GovernanceVote, []GovernanceTallyItem) {
	gVote := new(GovernanceVote)

	if len(header.Vote) > 0 {
		var err error

		if err := rlp.DecodeBytes(header.Vote, gVote); err != nil {
			logger.Error("Failed to decode a vote. This vote will be ignored", "number", header.Number, "key", gVote.Key, "value", gVote.Value, "validator", gVote.Validator)
			return valset, votes, tally
		}
		if gVote, err = gov.ParseVoteValue(gVote); err != nil {
			logger.Error("Failed to parse a vote value. This vote will be ignored", "number", header.Number, "key", gVote.Key, "value", gVote.Value, "validator", gVote.Validator)
			return valset, votes, tally
		}

		// If the given key is forbidden, stop processing
		if _, ok := GovernanceForbiddenKeyMap[gVote.Key]; ok {
			logger.Warn("Forbidden vote key was received", "key", gVote.Key, "value", gVote.Value, "from", gVote.Validator)
			return valset, votes, tally
		}

		key := GovernanceKeyMap[gVote.Key]
		switch key {
		case params.GoverningNode:
			_, addr := valset.GetByAddress(gVote.Value.(common.Address))
			if addr == nil {
				logger.Warn("Invalid governing node address", "number", header.Number, "Validator", gVote.Validator, "key", gVote.Key, "value", gVote.Value)
				return valset, votes, tally
			}
		case params.AddValidator:
			if !gov.checkVote(gVote.Value.(common.Address), true, valset) {
				return valset, votes, tally
			}
		case params.RemoveValidator:
			if !gov.checkVote(gVote.Value.(common.Address), false, valset) {
				return valset, votes, tally
			}
		}

		number := header.Number.Uint64()
		// Check vote's validity
		if gVote, ok := gov.ValidateVote(gVote); ok {
			governanceMode := GovernanceModeMap[gov.ChainConfig.Governance.GovernanceMode]
			governingNode := gov.ChainConfig.Governance.GoverningNode

			// Remove old vote with same validator and key
			votes, tally = gov.removePreviousVote(valset, votes, tally, proposer, gVote, governanceMode, governingNode)

			// Add new Vote to snapshot.GovernanceVotes
			votes = append(votes, *gVote)

			// Tally up the new vote. This will be cleared when Epoch ends.
			// Add to GovernanceTallies if it doesn't exist
			valset, votes, tally = gov.addNewVote(valset, votes, tally, gVote, governanceMode, governingNode, number)

			// If this vote was casted by this node, remove it
			if self == proposer {
				gov.removeDuplicatedVote(gVote, header.Number.Uint64())
			}
		} else {
			logger.Warn("Received Vote was invalid", "number", header.Number, "Validator", gVote.Validator, "key", gVote.Key, "value", gVote.Value)
		}
		if number > atomic.LoadUint64(&gov.lastGovernanceStateBlock) {
			gov.GovernanceVotes.Import(votes)
			gov.GovernanceTallies.Import(tally)
		}
	}
	return valset, votes, tally
}

func (gov *Governance) checkVote(address common.Address, authorize bool, valset istanbul.ValidatorSet) bool {
	_, validator := valset.GetByAddress(address)
	return (validator != nil && !authorize) || (validator == nil && authorize)
}

func (gov *Governance) isGovernanceModeSingleOrNone(governanceMode int, governingNode common.Address, voter common.Address) bool {
	return governanceMode == params.GovernanceMode_None || (governanceMode == params.GovernanceMode_Single && voter == governingNode)
}

func (gov *Governance) removePreviousVote(valset istanbul.ValidatorSet, votes []GovernanceVote, tally []GovernanceTallyItem, validator common.Address, gVote *GovernanceVote, governanceMode int, governingNode common.Address) ([]GovernanceVote, []GovernanceTallyItem) {
	ret := make([]GovernanceVote, len(votes))
	copy(ret, votes)

	// Removing duplicated previous GovernanceVotes
	for idx, vote := range votes {
		// Check if previous vote from same validator exists
		if vote.Validator == validator && vote.Key == gVote.Key {
			// Reduce Tally
			_, v := valset.GetByAddress(vote.Validator)
			vp := v.VotingPower()
			var currentVotes uint64
			currentVotes, tally = gov.changeGovernanceTally(tally, vote.Key, vote.Value, vp, false)

			// Remove the old vote from GovernanceVotes
			ret = append(votes[:idx], votes[idx+1:]...)
			if gov.isGovernanceModeSingleOrNone(governanceMode, governingNode, gVote.Validator) ||
				(governanceMode == params.GovernanceMode_Ballot && currentVotes <= valset.TotalVotingPower()/2) {
				if v, ok := gov.changeSet.GetValue(GovernanceKeyMap[vote.Key]); ok && v == vote.Value {
					gov.changeSet.RemoveItem(vote.Key)
				}
			}
			break
		}
	}
	return ret, tally
}

// changeGovernanceTally updates snapshot's tally for governance votes.
func (gov *Governance) changeGovernanceTally(tally []GovernanceTallyItem, key string, value interface{}, vp uint64, isAdd bool) (uint64, []GovernanceTallyItem) {
	found := false
	var currentVote uint64
	ret := make([]GovernanceTallyItem, len(tally))
	copy(ret, tally)

	for idx, v := range tally {
		if v.Key == key && v.Value == value {
			if isAdd {
				ret[idx].Votes += vp
			} else {
				if ret[idx].Votes > vp {
					ret[idx].Votes -= vp
				} else {
					ret[idx].Votes = uint64(0)
				}
			}

			currentVote = ret[idx].Votes

			if currentVote == 0 {
				ret = append(tally[:idx], tally[idx+1:]...)
			}
			found = true
			break
		}
	}

	if !found && isAdd {
		ret = append(ret, GovernanceTallyItem{Key: key, Value: value, Votes: vp})
		return vp, ret
	} else {
		return currentVote, ret
	}
}

func (gov *Governance) addNewVote(valset istanbul.ValidatorSet, votes []GovernanceVote, tally []GovernanceTallyItem, gVote *GovernanceVote, governanceMode int, governingNode common.Address, blockNum uint64) (istanbul.ValidatorSet, []GovernanceVote, []GovernanceTallyItem) {
	_, v := valset.GetByAddress(gVote.Validator)
	if v != nil {
		vp := v.VotingPower()
		var currentVotes uint64
		currentVotes, tally = gov.changeGovernanceTally(tally, gVote.Key, gVote.Value, vp, true)
		if gov.isGovernanceModeSingleOrNone(governanceMode, governingNode, gVote.Validator) ||
			(governanceMode == params.GovernanceMode_Ballot && currentVotes > valset.TotalVotingPower()/2) {
			switch GovernanceKeyMap[gVote.Key] {
			case params.AddValidator:
				valset.AddValidator(gVote.Value.(common.Address))
			case params.RemoveValidator:
				target := gVote.Value.(common.Address)
				valset.RemoveValidator(target)
				votes = gov.removeVotesFromRemovedNode(votes, target)
			default:
				if blockNum > atomic.LoadUint64(&gov.lastGovernanceStateBlock) {
					gov.ReflectVotes(*gVote)
				}
			}
		}
	}
	return valset, votes, tally
}

func (gov *Governance) removeVotesFromRemovedNode(votes []GovernanceVote, addr common.Address) []GovernanceVote {
	ret := make([]GovernanceVote, len(votes))
	copy(ret, votes)

	for i := 0; i < len(votes); i++ {
		if votes[i].Validator == addr {
			// Uncast the vote from the chronological list
			ret = append(votes[:i], votes[i+1:]...)
			i--
		}
	}
	return ret
}

func (gov *Governance) GetGovernanceItemAtNumber(num uint64, key string) (interface{}, error) {
	_, data, err := gov.ReadGovernance(num)
	if err != nil {
		return nil, err
	}

	if item, ok := data[key]; ok {
		if item == nil {
			return nil, ErrItemNil
		}
		return item, nil
	} else {
		return nil, ErrItemNotFound
	}
}

func (gov *Governance) GetLatestGovernanceItem(key int) interface{} {
	ret, _ := gov.currentSet.GetValue(key)
	return ret
}
