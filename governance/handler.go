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
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

type check struct {
	t         reflect.Type
	validator func(k string, v interface{}) bool
	trigger   func(g *Governance, k string, v interface{})
}

var (
	stringT  = reflect.TypeOf("")
	uint64T  = reflect.TypeOf(uint64(0))
	addressT = reflect.TypeOf(common.StringToAddress("0x0"))
	boolT    = reflect.TypeOf(true)
	float64T = reflect.TypeOf(float64(0.0))
)

var GovernanceItems = map[int]check{
	params.GovernanceMode:          {stringT, checkGovernanceMode, nil},
	params.GoverningNode:           {addressT, checkAddress, nil},
	params.UnitPrice:               {uint64T, checkUint64andBool, updateUnitPrice},
	params.AddValidator:            {addressT, checkAddressOrListOfUniqueAddresses, nil},
	params.RemoveValidator:         {addressT, checkAddressOrListOfUniqueAddresses, nil},
	params.MintingAmount:           {stringT, checkBigInt, nil},
	params.Ratio:                   {stringT, checkRatio, nil},
	params.UseGiniCoeff:            {boolT, checkUint64andBool, updateUseGiniCoeff},
	params.DeferredTxFee:           {boolT, checkUint64andBool, nil},
	params.MinimumStake:            {stringT, checkRewardMinimumStake, nil},
	params.StakeUpdateInterval:     {uint64T, checkUint64andBool, updateStakingUpdateInterval},
	params.ProposerRefreshInterval: {uint64T, checkUint64andBool, updateProposerUpdateInterval},
	params.Epoch:                   {uint64T, checkUint64andBool, nil},
	params.Policy:                  {uint64T, checkUint64andBool, updateProposerPolicy},
	params.CommitteeSize:           {uint64T, checkCommitteeSize, nil},
	params.ConstTxGasHumanReadable: {uint64T, checkUint64andBool, updateTxGasHumanReadable},
	params.Timeout:                 {uint64T, checkUint64andBool, nil},
}

func updateTxGasHumanReadable(g *Governance, k string, v interface{}) {
	params.TxGasHumanReadable = v.(uint64)
	logger.Info("TxGasHumanReadable changed", "New value", params.TxGasHumanReadable)
}

func updateUnitPrice(g *Governance, k string, v interface{}) {
	newPrice := v.(uint64)
	if g.TxPool != nil {
		g.TxPool.SetGasPrice(big.NewInt(0).SetUint64(newPrice))
	}
}

func updateUseGiniCoeff(g *Governance, k string, v interface{}) {
	if g.blockChain != nil {
		g.blockChain.SetUseGiniCoeff(g.UseGiniCoeff())
	}
}

func updateStakingUpdateInterval(g *Governance, k string, v interface{}) {
	params.SetStakingUpdateInterval(g.StakingUpdateInterval())
}

func updateProposerUpdateInterval(g *Governance, k string, v interface{}) {
	params.SetProposerUpdateInterval(g.ProposerUpdateInterval())
}

func updateProposerPolicy(g *Governance, k string, v interface{}) {
	if g.blockChain != nil {
		g.blockChain.SetProposerPolicy(g.ProposerPolicy())
	}
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

	// When an int value comes from JS console, it comes as a float64
	if GovernanceItems[k].t == uint64T {
		v, ok := val.(float64)
		if !ok {
			return val
		}
		if float64(uint64(v)) == v {
			return uint64(v)
		}
		return val
	}

	// Otherwise, it comes as a string
	v, ok := val.(string)
	if !ok {
		return val
	}
	if GovernanceItems[k].t == addressT {
		addresses := strings.Split(v, ",")
		switch len(addresses) {
		case 0:
			return val
		case 1:
			str := strings.Trim(v, " ")
			if common.IsHexAddress(str) {
				return common.HexToAddress(str)
			} else {
				return val
			}
		default:
			var nodeAddresses []common.Address
			for _, str := range addresses {
				str = strings.Trim(str, " ")
				if common.IsHexAddress(str) {
					nodeAddresses = append(nodeAddresses, common.HexToAddress(str))
				} else {
					return val
				}
			}
			return nodeAddresses
		}
	} else {
		// If a string text come as uppercase, make it into lowercase
		return strings.ToLower(v)
	}
}

func checkValueType(v interface{}, expectType reflect.Type) bool {
	var ok bool
	switch expectType {
	case uint64T:
		_, ok = v.(uint64)
	case stringT:
		_, ok = v.(string)
	case addressT:
		_, ok = v.(common.Address)
		if !ok {
			_, ok = v.([]common.Address)
		}
	case boolT:
		_, ok = v.(bool)
	default:
		ok = false
	}
	return ok
}

func checkKey(k string) bool {
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

	if checkKey(vote.Key) && checkValueType(vote.Value, GovernanceItems[key].t) {
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

func checkCommitteeSize(k string, v interface{}) bool {
	if !checkUint64andBool(k, v) {
		return false
	}
	if v == uint64(0) {
		return false
	}
	return true
}

func checkRewardMinimumStake(k string, v interface{}) bool {
	if !checkBigInt(k, v) {
		return false
	}
	if v, ok := new(big.Int).SetString(v.(string), 10); ok {
		if v.Cmp(common.Big0) < 0 {
			return false
		}
	}
	return true
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
	if _, ok := v.(common.Address); ok {
		return true
	}
	return false
}

func checkAddressOrListOfUniqueAddresses(k string, v interface{}) bool {
	if checkAddress(k, v) {
		return true
	}
	if _, ok := v.([]common.Address); !ok {
		return false
	}

	// there should not be duplicated addresses, if value contains multiple addresses
	addressExists := make(map[common.Address]bool)
	for _, address := range v.([]common.Address) {
		if _, ok := addressExists[address]; ok {
			return false
		} else {
			addressExists[address] = true
		}
	}
	return true
}

func isEqualValue(k int, v1 interface{}, v2 interface{}) bool {
	if reflect.TypeOf(v1) != reflect.TypeOf(v2) {
		return false
	}

	if GovernanceItems[k].t == addressT && reflect.TypeOf(v1) != addressT {
		value1, ok1 := v1.([]common.Address)
		value2, ok2 := v2.([]common.Address)
		if ok1 == false || ok2 == false {
			return false
		}
		return reflect.DeepEqual(value1, value2)
	}

	return v1 == v2
}

func (gov *Governance) HandleGovernanceVote(valset istanbul.ValidatorSet, votes []GovernanceVote, tally []GovernanceTallyItem, header *types.Header, proposer common.Address, self common.Address) (istanbul.ValidatorSet, []GovernanceVote, []GovernanceTallyItem) {
	gVote := new(GovernanceVote)

	if len(header.Vote) > 0 {
		var err error

		if err := rlp.DecodeBytes(header.Vote, gVote); err != nil {
			logger.Error("Failed to decode a vote. This vote will be ignored", "number", header.Number)
			return valset, votes, tally
		}
		if gVote, err = gov.ParseVoteValue(gVote); err != nil {
			logger.Error("Failed to parse a vote value. This vote will be ignored", "number", header.Number)
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
			v, ok := gVote.Value.(common.Address)
			if !ok {
				logger.Warn("Invalid value Type", "number", header.Number, "Validator", gVote.Validator, "key", gVote.Key, "value", gVote.Value)
				return valset, votes, tally
			}
			_, addr := valset.GetByAddress(v)
			if addr == nil {
				logger.Warn("Invalid governing node address", "number", header.Number, "Validator", gVote.Validator, "key", gVote.Key, "value", gVote.Value)
				return valset, votes, tally
			}
		case params.AddValidator, params.RemoveValidator:
			authorize := key == params.AddValidator
			if addr, ok := gVote.Value.(common.Address); ok {
				if !gov.checkVote(addr, authorize, valset) {
					if proposer == self {
						logger.Warn("A meaningless vote has been proposed. It is being removed without further handling", "key", gVote.Key, "value", gVote.Value)
						gov.removeDuplicatedVote(gVote, header.Number.Uint64())
					}
					return valset, votes, tally
				}
			} else if addresses, ok := gVote.Value.([]common.Address); ok {
				for _, address := range addresses {
					if !gov.checkVote(address, authorize, valset) {
						if proposer == self {
							logger.Warn("A meaningless vote has been proposed. It is being removed without further handling", "key", gVote.Key, "value", gVote.Value)
							gov.removeDuplicatedVote(gVote, header.Number.Uint64())
						}
						return valset, votes, tally
					}
				}
			} else {
				logger.Warn("Invalid value Type", "number", header.Number, "Validator", gVote.Validator, "key", gVote.Key, "value", gVote.Value)
			}
		}

		number := header.Number.Uint64()
		// Check vote's validity
		if gVote, ok := gov.ValidateVote(gVote); ok {
			governanceMode := GovernanceModeMap[gov.GovernanceMode()]
			governingNode := gov.GoverningNode()

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
	if validator == nil {
		_, validator = valset.GetDemotedByAddress(address)
	}
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
		if v.Key == key && isEqualValue(GovernanceKeyMap[key], v.Value, value) {
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
				//reward.GetStakingInfo()
				if addr, ok := gVote.Value.(common.Address); ok {
					valset.AddValidator(addr)
				} else {
					for _, address := range gVote.Value.([]common.Address) {
						valset.AddValidator(address)
					}
				}
			case params.RemoveValidator:
				if addr, ok := gVote.Value.(common.Address); ok {
					valset.RemoveValidator(addr)
				} else {
					for _, target := range gVote.Value.([]common.Address) {
						valset.RemoveValidator(target)
						votes = gov.removeVotesFromRemovedNode(votes, target)
					}
				}
			case params.Timeout:
				timeout := gVote.Value.(uint64)
				atomic.StoreUint64(&istanbul.DefaultConfig.Timeout, timeout)
				if blockNum > atomic.LoadUint64(&gov.lastGovernanceStateBlock) {
					gov.ReflectVotes(*gVote)
				}
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

func (gov *Governance) GetItemAtNumberByIntKey(num uint64, key int) (interface{}, error) {
	return gov.GetGovernanceItemAtNumber(num, GovernanceKeyMapReverse[key])
}

// GetGoverningInfoAtNumber returns whether the governing mode is single or not and the governing node.
func (gov *Governance) GetGoverningInfoAtNumber(num uint64) (bool, common.Address, error) {
	govMode, err := gov.GetItemAtNumberByIntKey(num, params.GovernanceMode)
	if err != nil {
		return false, common.Address{}, err
	}

	if GovernanceModeMap[govMode.(string)] != params.GovernanceMode_Single {
		return false, common.Address{}, nil
	}

	govNode, err := gov.GetItemAtNumberByIntKey(num, params.GoverningNode)
	if err != nil {
		return true, common.Address{}, err
	}

	return true, govNode.(common.Address), nil
}

func (gov *Governance) GetMinimumStakingAtNumber(num uint64) (uint64, error) {
	minStaking, err := gov.GetItemAtNumberByIntKey(num, params.MinimumStake)
	if err != nil {
		return 0, err
	}
	bigMinStaking, ok := new(big.Int).SetString(minStaking.(string), 10)
	if !ok {
		return 0, fmt.Errorf("invalid number string: %v", minStaking)
	}
	return bigMinStaking.Uint64(), nil
}
