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
	"errors"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
)

type PublicGovernanceAPI struct {
	governance Engine // Node interfaced by this API
}

type returnTally struct {
	Key                string
	Value              interface{}
	ApprovalPercentage float64
}

func NewGovernanceAPI(gov Engine) *PublicGovernanceAPI {
	return &PublicGovernanceAPI{governance: gov}
}

type GovernanceKlayAPI struct {
	governance Engine
	chain      blockChain
}

func NewGovernanceKlayAPI(gov Engine, chain blockChain) *GovernanceKlayAPI {
	return &GovernanceKlayAPI{governance: gov, chain: chain}
}

var (
	errUnknownBlock           = errors.New("Unknown block")
	errNotAvailableInThisMode = errors.New("In current governance mode, voting power is not available")
	errSetDefaultFailure      = errors.New("Failed to set a default value")
	errPermissionDenied       = errors.New("You don't have the right to vote")
	errRemoveSelf             = errors.New("You can't vote on removing yourself")
	errInvalidKeyValue        = errors.New("Your vote couldn't be placed. Please check your vote's key and value")
)

// TODO-Klaytn-Governance: Refine this API and consider the gas price of txpool
func (api *GovernanceKlayAPI) GasPriceAt(num *rpc.BlockNumber) (*hexutil.Big, error) {
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		ret := api.governance.UnitPrice()
		return (*hexutil.Big)(big.NewInt(0).SetUint64(ret)), nil
	} else {
		blockNum := num.Int64()

		if blockNum > api.chain.CurrentHeader().Number.Int64() {
			return nil, errUnknownBlock
		}

		if ret, err := api.GasPriceAtNumber(blockNum); err != nil {
			return nil, err
		} else {
			return (*hexutil.Big)(big.NewInt(0).SetUint64(ret)), nil
		}
	}
}

// Vote injects a new vote for governance targets such as unitprice and governingnode.
func (api *PublicGovernanceAPI) Vote(key string, val interface{}) (string, error) {
	gMode := api.governance.GovernanceMode()
	gNode := api.governance.GoverningNode()

	if GovernanceModeMap[gMode] == params.GovernanceMode_Single && gNode != api.governance.NodeAddress() {
		return "", errPermissionDenied
	}
	if strings.ToLower(key) == "governance.removevalidator" {
		if _, ok := api.governance.ValidateVote(&GovernanceVote{Key: key, Value: val}); !ok {
			return "", errInvalidKeyValue
		}
		if api.isRemovingSelf(val.(string)) {
			return "", errRemoveSelf
		}
	}
	if api.governance.AddVote(key, val) {
		return "Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate.", nil
	}
	return "", errInvalidKeyValue
}

func (api *PublicGovernanceAPI) isRemovingSelf(val string) bool {
	for _, str := range strings.Split(val, ",") {
		str = strings.Trim(str, " ")
		if common.HexToAddress(str) == api.governance.NodeAddress() {
			return true
		}
	}
	return false
}

func (api *PublicGovernanceAPI) ShowTally() []*returnTally {
	ret := []*returnTally{}

	for _, val := range api.governance.GetGovernanceTalliesCopy() {
		item := &returnTally{
			Key:                val.Key,
			Value:              val.Value,
			ApprovalPercentage: float64(val.Votes) / float64(api.governance.TotalVotingPower()) * 100,
		}
		ret = append(ret, item)
	}

	return ret
}

func (api *PublicGovernanceAPI) TotalVotingPower() (float64, error) {
	if !api.isGovernanceModeBallot() {
		return 0, errNotAvailableInThisMode
	}
	return float64(api.governance.TotalVotingPower()) / 1000.0, nil
}

func (api *PublicGovernanceAPI) ItemsAt(num *rpc.BlockNumber) (map[string]interface{}, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.governance.BlockChain().CurrentHeader().Number.Uint64()
	} else {
		blockNumber = uint64(num.Int64())
	}
	_, data, error := api.governance.ReadGovernance(blockNumber)
	if error == nil {
		return data, error
	} else {
		return nil, error
	}
}

func (api *PublicGovernanceAPI) GetStakingInfo(num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.governance.BlockChain().CurrentHeader().Number.Uint64()
	} else {
		blockNumber = uint64(num.Int64())
	}
	return reward.GetStakingInfo(blockNumber), nil
}

func (api *PublicGovernanceAPI) PendingChanges() map[string]interface{} {
	return api.governance.PendingChanges()
}

func (api *PublicGovernanceAPI) Votes() []GovernanceVote {
	return api.governance.Votes()
}

func (api *PublicGovernanceAPI) IdxCache() []uint64 {
	return api.governance.IdxCache()
}

func (api *PublicGovernanceAPI) IdxCacheFromDb() []uint64 {
	return api.governance.IdxCacheFromDb()
}

// TODO-Klaytn: Return error if invalid input is given such as pending or a too big number
func (api *PublicGovernanceAPI) ItemCacheFromDb(num *rpc.BlockNumber) map[string]interface{} {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.governance.BlockChain().CurrentHeader().Number.Uint64()
	} else {
		blockNumber = uint64(num.Int64())
	}
	ret, _ := api.governance.DB().ReadGovernance(blockNumber)
	return ret
}

type VoteList struct {
	Key      string
	Value    interface{}
	Casted   bool
	BlockNum uint64
}

func (api *PublicGovernanceAPI) MyVotes() []*VoteList {
	ret := []*VoteList{}

	for k, v := range api.governance.GetVoteMapCopy() {
		item := &VoteList{
			Key:      k,
			Value:    v.Value,
			Casted:   v.Casted,
			BlockNum: v.Num,
		}
		ret = append(ret, item)
	}

	return ret
}

func (api *PublicGovernanceAPI) MyVotingPower() (float64, error) {
	if !api.isGovernanceModeBallot() {
		return 0, errNotAvailableInThisMode
	}
	return float64(api.governance.MyVotingPower()) / 1000.0, nil
}

func (api *PublicGovernanceAPI) ChainConfig() *params.ChainConfig {
	return api.governance.InitialChainConfig()
}

func (api *PublicGovernanceAPI) NodeAddress() common.Address {
	return api.governance.NodeAddress()
}

func (api *PublicGovernanceAPI) isGovernanceModeBallot() bool {
	if GovernanceModeMap[api.governance.GovernanceMode()] == params.GovernanceMode_Ballot {
		return true
	}
	return false
}

func (api *GovernanceKlayAPI) GasPriceAtNumber(num int64) (uint64, error) {
	val, err := api.governance.GetGovernanceItemAtNumber(uint64(num), GovernanceKeyMapReverse[params.UnitPrice])
	if err != nil {
		logger.Error("Failed to retrieve unit price", "err", err)
		return 0, err
	}
	return val.(uint64), nil
}

// Disabled APIs
// func (api *GovernanceKlayAPI) GetTxGasHumanReadable(num *rpc.BlockNumber) (uint64, error) {
// 	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
// 		// If the value hasn't been set in governance, set it with default value
// 		if ret := api.governance.GetGovernanceValue(params.ConstTxGasHumanReadable); ret == nil {
// 			return api.setDefaultTxGasHumanReadable()
// 		} else {
// 			return ret.(uint64), nil
// 		}
// 	} else {
// 		blockNum := num.Int64()
//
// 		if blockNum > api.chain.CurrentHeader().Number.Int64() {
// 			return 0, errUnknownBlock
// 		}
//
// 		if ret, err := api.governance.GetGovernanceItemAtNumber(uint64(blockNum), GovernanceKeyMapReverse[params.ConstTxGasHumanReadable]); err == nil && ret != nil {
// 			return ret.(uint64), nil
// 		} else {
// 			logger.Error("Failed to retrieve TxGasHumanReadable, sending default value", "err", err)
// 			return api.setDefaultTxGasHumanReadable()
// 		}
// 	}
// }
//
// func (api *GovernanceKlayAPI) setDefaultTxGasHumanReadable() (uint64, error) {
// 	err := api.governance.currentSet.SetValue(params.ConstTxGasHumanReadable, params.TxGasHumanReadable)
// 	if err != nil {
// 		return 0, errSetDefaultFailure
// 	} else {
// 		return params.TxGasHumanReadable, nil
// 	}
// }
