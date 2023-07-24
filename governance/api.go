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
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
)

type GovernanceAPI struct {
	governance Engine // Node interfaced by this API
}

type returnTally struct {
	Key                string
	Value              interface{}
	ApprovalPercentage float64
}

func NewGovernanceAPI(gov Engine) *GovernanceAPI {
	return &GovernanceAPI{governance: gov}
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

func (api *GovernanceKlayAPI) GetChainConfig(num *rpc.BlockNumber) *params.ChainConfig {
	return getChainConfig(api.governance, num)
}

func (api *GovernanceKlayAPI) GetStakingInfo(num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	return getStakingInfo(api.governance, num)
}

func (api *GovernanceKlayAPI) GetParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return getParams(api.governance, num)
}

func (api *GovernanceKlayAPI) NodeAddress() common.Address {
	return api.governance.NodeAddress()
}

// GetRewards returns detailed information of the block reward at a given block number.
func (api *GovernanceKlayAPI) GetRewards(num *rpc.BlockNumber) (*reward.RewardSpec, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber {
		blockNumber = api.chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	header := api.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNumber)
	}

	rules := api.chain.Config().Rules(new(big.Int).SetUint64(blockNumber))
	pset, err := api.governance.EffectiveParams(blockNumber)
	if err != nil {
		return nil, err
	}
	rewardParamNum := reward.CalcRewardParamBlock(header.Number.Uint64(), pset.Epoch(), rules)
	rewardParamSet, err := api.governance.EffectiveParams(rewardParamNum)
	if err != nil {
		return nil, err
	}

	return reward.GetBlockReward(header, rules, rewardParamSet)
}

type AccumulatedRewards struct {
	FirstBlockTime string   `json:"firstBlockTime"`
	LastBlockTime  string   `json:"lastBlockTime"`
	FirstBlock     *big.Int `json:"firstBlock"`
	LastBlock      *big.Int `json:"lastBlock"`

	// TotalMinted + TotalTxFee - TotalBurntTxFee = TotalProposerRewards + TotalStakingRewards + TotalKFFRewards + TotalKCFRewards
	TotalMinted          *big.Int                    `json:"totalMinted"`
	TotalTxFee           *big.Int                    `json:"totalTxFee"`
	TotalBurntTxFee      *big.Int                    `json:"totalBurntTxFee"`
	TotalProposerRewards *big.Int                    `json:"totalProposerRewards"`
	TotalStakingRewards  *big.Int                    `json:"totalStakingRewards"`
	TotalKFFRewards      *big.Int                    `json:"totalKFFRewards"`
	TotalKCFRewards      *big.Int                    `json:"totalKCFRewards"`
	Rewards              map[common.Address]*big.Int `json:"rewards"`
}

// GetRewardsAccumulated returns accumulated rewards data in the block range of [first, last].
func (api *GovernanceAPI) GetRewardsAccumulated(first rpc.BlockNumber, last rpc.BlockNumber) (*AccumulatedRewards, error) {
	blockchain := api.governance.BlockChain()
	govKlayAPI := NewGovernanceKlayAPI(api.governance, blockchain)

	currentBlock := blockchain.CurrentBlock().NumberU64()

	firstBlock := currentBlock
	if first >= rpc.EarliestBlockNumber {
		firstBlock = uint64(first.Int64())
	}

	lastBlock := currentBlock
	if last >= rpc.EarliestBlockNumber {
		lastBlock = uint64(last.Int64())
	}

	if firstBlock > lastBlock {
		return nil, errors.New("the last block number should be equal or larger the first block number")
	}

	if lastBlock > currentBlock {
		return nil, errors.New("the last block number should be equal or less than the current block number")
	}

	blockCount := lastBlock - firstBlock + 1
	if blockCount > 604800 { // 7 days. naive resource protection
		return nil, errors.New("block range should be equal or less than 604800")
	}

	// initialize structures before request a job
	accumRewards := &AccumulatedRewards{}
	blockRewards := reward.NewRewardSpec()
	mu := sync.Mutex{} // protect blockRewards

	numWorkers := runtime.NumCPU()
	reqCh := make(chan uint64, numWorkers)
	errCh := make(chan error, numWorkers+1)
	wg := sync.WaitGroup{}

	// introduce the worker pattern to prevent resource exhaustion
	for i := 0; i < numWorkers; i++ {
		go func() {
			// the minimum digit of request is period to avoid current access to an accArray item
			for num := range reqCh {
				bn := rpc.BlockNumber(num)
				blockReward, err := govKlayAPI.GetRewards(&bn)
				if err != nil {
					errCh <- err
					wg.Done()
					return
				}

				mu.Lock()
				blockRewards.Add(blockReward)
				mu.Unlock()

				wg.Done()
			}
		}()
	}

	// write the information of the first block
	header := blockchain.GetHeaderByNumber(firstBlock)
	if header == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", firstBlock)
	}
	accumRewards.FirstBlock = header.Number
	accumRewards.FirstBlockTime = time.Unix(header.Time.Int64(), 0).String()

	// write the information of the last block
	header = blockchain.GetHeaderByNumber(lastBlock)
	if header == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", lastBlock)
	}
	accumRewards.LastBlock = header.Number
	accumRewards.LastBlockTime = time.Unix(header.Time.Int64(), 0).String()

	for num := firstBlock; num <= lastBlock; num++ {
		wg.Add(1)
		reqCh <- num
	}

	// generate a goroutine to return error early
	go func() {
		wg.Wait()
		errCh <- nil
	}()

	if err := <-errCh; err != nil {
		return nil, err
	}

	// collect the accumulated rewards information
	accumRewards.Rewards = blockRewards.Rewards
	accumRewards.TotalMinted = blockRewards.Minted
	accumRewards.TotalTxFee = blockRewards.TotalFee
	accumRewards.TotalBurntTxFee = blockRewards.BurntFee
	accumRewards.TotalProposerRewards = blockRewards.Proposer
	accumRewards.TotalStakingRewards = blockRewards.Stakers
	accumRewards.TotalKFFRewards = blockRewards.KFF
	accumRewards.TotalKCFRewards = blockRewards.KCF

	return accumRewards, nil
}

// Vote injects a new vote for governance targets such as unitprice and governingnode.
func (api *GovernanceAPI) Vote(key string, val interface{}) (string, error) {
	blockNumber := api.governance.BlockChain().CurrentBlock().NumberU64()
	pset, err := api.governance.EffectiveParams(blockNumber + 1)
	if err != nil {
		return "", err
	}
	gMode := pset.GovernanceModeInt()
	gNode := pset.GoverningNode()

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
		if vote.Value.(uint64) > pset.UpperBoundBaseFee() {
			return "", errInvalidLowerBound
		}
	}
	if vote.Key == "kip71.upperboundbasefee" {
		if vote.Value.(uint64) < pset.LowerBoundBaseFee() {
			return "", errInvalidUpperBound
		}
	}
	if api.governance.AddVote(key, val) {
		return "Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate.", nil
	}
	return "", errInvalidKeyValue
}

func (api *GovernanceAPI) isRemovingSelf(val string) bool {
	for _, str := range strings.Split(val, ",") {
		str = strings.Trim(str, " ")
		if common.HexToAddress(str) == api.governance.NodeAddress() {
			return true
		}
	}
	return false
}

func (api *GovernanceAPI) ShowTally() []*returnTally {
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

func (api *GovernanceAPI) TotalVotingPower() (float64, error) {
	if !api.isGovernanceModeBallot() {
		return 0, errNotAvailableInThisMode
	}
	return float64(api.governance.TotalVotingPower()) / 1000.0, nil
}

func (api *GovernanceAPI) GetParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return getParams(api.governance, num)
}

func getParams(governance Engine, num *rpc.BlockNumber) (map[string]interface{}, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = governance.BlockChain().CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	pset, err := governance.EffectiveParams(blockNumber)
	if err != nil {
		return nil, err
	}
	return pset.StrMap(), nil
}

func (api *GovernanceAPI) GetStakingInfo(num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	return getStakingInfo(api.governance, num)
}

func getStakingInfo(governance Engine, num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = governance.BlockChain().CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}
	return reward.GetStakingInfo(blockNumber), nil
}

func (api *GovernanceAPI) PendingChanges() map[string]interface{} {
	return api.governance.PendingChanges()
}

func (api *GovernanceAPI) Votes() []GovernanceVote {
	return api.governance.Votes()
}

func (api *GovernanceAPI) IdxCache() []uint64 {
	return api.governance.IdxCache()
}

func (api *GovernanceAPI) IdxCacheFromDb() []uint64 {
	return api.governance.IdxCacheFromDb()
}

// TODO-Klaytn: Return error if invalid input is given such as pending or a too big number
func (api *GovernanceAPI) ItemCacheFromDb(num *rpc.BlockNumber) map[string]interface{} {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.governance.BlockChain().CurrentBlock().NumberU64()
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

func (api *GovernanceAPI) MyVotes() []*VoteList {
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

func (api *GovernanceAPI) MyVotingPower() (float64, error) {
	if !api.isGovernanceModeBallot() {
		return 0, errNotAvailableInThisMode
	}
	return float64(api.governance.MyVotingPower()) / 1000.0, nil
}

func (api *GovernanceAPI) GetChainConfig(num *rpc.BlockNumber) *params.ChainConfig {
	return getChainConfig(api.governance, num)
}

func getChainConfig(governance Engine, num *rpc.BlockNumber) *params.ChainConfig {
	var blocknum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blocknum = governance.BlockChain().CurrentBlock().NumberU64()
	} else {
		blocknum = num.Uint64()
	}

	pset, err := governance.EffectiveParams(blocknum)
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
	config.Kip103CompatibleBlock = latestConfig.Kip103CompatibleBlock
	config.Kip103ContractAddress = latestConfig.Kip103ContractAddress
	config.MantleCompatibleBlock = latestConfig.MantleCompatibleBlock

	return config
}

func (api *GovernanceAPI) NodeAddress() common.Address {
	return api.governance.NodeAddress()
}

func (api *GovernanceAPI) isGovernanceModeBallot() bool {
	blockNumber := api.governance.BlockChain().CurrentBlock().NumberU64()
	pset, err := api.governance.EffectiveParams(blockNumber + 1)
	if err != nil {
		return false
	}
	gMode := pset.GovernanceModeInt()
	return gMode == params.GovernanceMode_Ballot
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
// 		if blockNum > api.chain.CurrentBlock().NumberU64() {
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
