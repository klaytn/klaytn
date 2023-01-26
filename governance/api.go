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
	"fmt"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/common/hexutil"

	"github.com/klaytn/klaytn/common"
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
	errInvalidLowerBound      = errors.New("lowerboundbasefee cannot be set exceeding upperboundbasefee")
	errInvalidUpperBound      = errors.New("upperboundbasefee cannot be set lower than lowerboundbasefee")
)

func (api *GovernanceKlayAPI) GetStakingInfo(num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	return getStakingInfo(api.governance, num)
}

func (api *GovernanceKlayAPI) GovParamsAt(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return itemsAt(api.governance, num)
}

func (api *GovernanceKlayAPI) NodeAddress() common.Address {
	return api.governance.NodeAddress()
}

// GasPriceAt returns the base fee of the given block in peb,
// or returns unit price by using governance if there is no base fee set in header,
// or returns gas price of txpool if the block is pending block.
func (api *GovernanceKlayAPI) GasPriceAt(num *rpc.BlockNumber) (*hexutil.Big, error) {
	if num == nil || *num == rpc.LatestBlockNumber {
		header := api.chain.CurrentHeader()
		if header.BaseFee == nil {
			return (*hexutil.Big)(new(big.Int).SetUint64(api.governance.Params().UnitPrice())), nil
		}
		return (*hexutil.Big)(header.BaseFee), nil
	} else if *num == rpc.PendingBlockNumber {
		txpool := api.governance.GetTxPool()
		return (*hexutil.Big)(txpool.GasPrice()), nil
	} else {
		blockNum := num.Uint64()

		// Return the BaseFee in header at the block number
		header := api.chain.GetHeaderByNumber(blockNum)
		if blockNum > api.chain.CurrentHeader().Number.Uint64() || header == nil {
			return nil, errUnknownBlock
		} else if header.BaseFee != nil {
			return (*hexutil.Big)(header.BaseFee), nil
		}

		// Return the UnitPrice in governance data at the block number
		if ret, err := api.GasPriceAtNumber(blockNum); err != nil {
			return nil, err
		} else {
			return (*hexutil.Big)(new(big.Int).SetUint64(ret)), nil
		}
	}
}

// GetRewards returns detailed information of the block reward at a given block number.
func (api *GovernanceKlayAPI) GetRewards(num *rpc.BlockNumber) (*reward.RewardSpec, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber {
		blockNumber = api.governance.BlockChain().CurrentHeader().Number.Uint64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	header := api.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNumber)
	}

	rules := api.chain.Config().Rules(new(big.Int).SetUint64(blockNumber))
	pset, err := api.governance.ParamsAt(blockNumber)
	if err != nil {
		return nil, err
	}
	rewardParamNum := reward.CalcRewardParamBlock(header.Number.Uint64(), pset.Epoch(), rules)
	rewardParamSet, err := api.governance.ParamsAt(rewardParamNum)
	if err != nil {
		return nil, err
	}

	return reward.GetBlockReward(header, rules, rewardParamSet)
}

func (api *GovernanceKlayAPI) ChainConfig() *params.ChainConfig {
	num := rpc.LatestBlockNumber
	return chainConfigAt(api.governance, &num)
}

func (api *GovernanceKlayAPI) ChainConfigAt(num *rpc.BlockNumber) *params.ChainConfig {
	return chainConfigAt(api.governance, num)
}

// Vote injects a new vote for governance targets such as unitprice and governingnode.
func (api *PublicGovernanceAPI) Vote(key string, val interface{}) (string, error) {
	gMode := api.governance.Params().GovernanceModeInt()
	gNode := api.governance.Params().GoverningNode()

	if gMode == params.GovernanceMode_Single && gNode != api.governance.NodeAddress() {
		return "", errPermissionDenied
	}
	vote, ok := api.governance.ValidateVote(&GovernanceVote{Key: strings.ToLower(key), Value: val})
	if !ok {
		return "", errInvalidKeyValue
	}
	if vote.Key == "governance.removevalidator" {
		if api.isRemovingSelf(val.(string)) {
			return "", errRemoveSelf
		}
	}
	if vote.Key == "kip71.lowerboundbasefee" {
		if vote.Value.(uint64) > api.governance.Params().UpperBoundBaseFee() {
			return "", errInvalidLowerBound
		}
	}
	if vote.Key == "kip71.upperboundbasefee" {
		if vote.Value.(uint64) < api.governance.Params().LowerBoundBaseFee() {
			return "", errInvalidUpperBound
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
	return itemsAt(api.governance, num)
}

func itemsAt(governance Engine, num *rpc.BlockNumber) (map[string]interface{}, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = governance.BlockChain().CurrentHeader().Number.Uint64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	pset, err := governance.ParamsAt(blockNumber)
	if err != nil {
		return nil, err
	}
	return pset.StrMap(), nil
}

func (api *PublicGovernanceAPI) GetStakingInfo(num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	return getStakingInfo(api.governance, num)
}

func getStakingInfo(governance Engine, num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = governance.BlockChain().CurrentHeader().Number.Uint64()
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
	num := rpc.LatestBlockNumber
	return chainConfigAt(api.governance, &num)
}

func (api *PublicGovernanceAPI) ChainConfigAt(num *rpc.BlockNumber) *params.ChainConfig {
	return chainConfigAt(api.governance, num)
}

func chainConfigAt(governance Engine, num *rpc.BlockNumber) *params.ChainConfig {
	var blocknum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blocknum = governance.BlockChain().CurrentHeader().Number.Uint64()
	} else {
		blocknum = num.Uint64()
	}

	pset, err := governance.ParamsAt(blocknum)
	if err != nil {
		return nil
	}

	latestConfig := governance.BlockChain().Config()
	config := pset.ToChainConfig()
	config.ChainID = latestConfig.ChainID
	config.IstanbulCompatibleBlock = latestConfig.IstanbulCompatibleBlock
	config.LondonCompatibleBlock = latestConfig.LondonCompatibleBlock
	config.EthTxTypeCompatibleBlock = latestConfig.EthTxTypeCompatibleBlock
	config.MagmaCompatibleBlock = latestConfig.MagmaCompatibleBlock
	config.KoreCompatibleBlock = latestConfig.KoreCompatibleBlock

	return config
}

func (api *PublicGovernanceAPI) NodeAddress() common.Address {
	return api.governance.NodeAddress()
}

func (api *PublicGovernanceAPI) isGovernanceModeBallot() bool {
	gMode := api.governance.Params().GovernanceModeInt()
	return gMode == params.GovernanceMode_Ballot
}

func (api *GovernanceKlayAPI) GasPriceAtNumber(num uint64) (uint64, error) {
	pset, err := api.governance.ParamsAt(num)
	if err != nil {
		logger.Error("Failed to retrieve unit price", "err", err)
		return 0, err
	}
	return pset.UnitPrice(), nil
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
