// Copyright 2018 The klaytn Authors
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

//go:generate abigen --sol contract/KlaytnReward.sol --pkg contract --out contract/KlaytnReward.go
//go:generate abigen --sol contract/AddressBook.sol --pkg contract --out contract/AddressBook.go

package reward

import (
	"errors"
	"fmt"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/contracts/reward/contract"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var logger = log.NewModuleLogger(log.Reward)

type Reward struct {
	*contract.KlaytnRewardSession
	contractBackend bind.ContractBackend
}

const (
	AddrNotFoundInCouncilNodes = -1
)

// StakingInfo contains staking information.
type StakingInfo struct {
	BlockNum uint64 // Block number where staking information of Council is fetched

	// Information retrieved from AddressBook smart contract
	CouncilNodeIds       []common.Address // NodeIds of Council
	CouncilStakingdAddrs []common.Address // Address of Staking account which holds staking balance
	CouncilRewardAddrs   []common.Address // Address of Council account which will get block reward
	KIRAddr              common.Address   // Address of KIR contract
	PoCAddr              common.Address   // Address of PoC contract

	UseGini bool
	Gini    float64 // gini coefficient

	// Derived from CouncilStakingAddrs
	CouncilStakingAmounts []uint64 // Staking amounts of Council
}

func (s *StakingInfo) String() string {
	str := make([]string, 0)

	header := fmt.Sprintf("StakingInfo:{BlockNum:%v", s.BlockNum)
	str = append(str, header)

	// nodeIds
	nodeIdsHeader := fmt.Sprintf(" CouncilNodeIds:[")
	str = append(str, nodeIdsHeader)
	nodeIds := make([]string, 0)
	for _, nodeId := range s.CouncilNodeIds {
		nodeIds = append(nodeIds, nodeId.String())
	}
	str = append(str, strings.Join(nodeIds, " "))
	str = append(str, "]")

	// stakingAddrs
	stakingAddrsHeader := fmt.Sprintf(", CouncilStakingAddrs:[")
	str = append(str, stakingAddrsHeader)
	stakingAddrs := make([]string, 0)
	for _, stakingAddr := range s.CouncilStakingdAddrs {
		stakingAddrs = append(stakingAddrs, stakingAddr.String())
	}
	str = append(str, strings.Join(stakingAddrs, " "))
	str = append(str, "]")

	// rewardAddrs
	rewardAddrsHeader := fmt.Sprintf(", CouncilRewardAddrs:[")
	str = append(str, rewardAddrsHeader)
	rewardAddrs := make([]string, 0)
	for _, rewardAddr := range s.CouncilRewardAddrs {
		rewardAddrs = append(rewardAddrs, rewardAddr.String())
	}
	str = append(str, strings.Join(rewardAddrs, " "))
	str = append(str, "]")

	// pocAddr and kirAddr
	pocAddr := fmt.Sprintf(", PoCAddr:%v", s.PoCAddr.String())
	str = append(str, pocAddr)
	kirAddr := fmt.Sprintf(", KIRAddr:%v", s.KIRAddr.String())
	str = append(str, kirAddr)

	// stakingAmounts
	stakingAmounts := fmt.Sprintf(", CouncilStakingAmounts:%v", s.CouncilStakingAmounts)
	str = append(str, stakingAmounts)

	gini := fmt.Sprintf(", Gini:%v", s.Gini)
	str = append(str, gini)

	useGini := fmt.Sprintf(", UseGini:%v }", s.UseGini)
	str = append(str, useGini)

	return strings.Join(str, " ")
}

func (s *StakingInfo) GetIndexByNodeId(nodeId common.Address) int {
	for i, addr := range s.CouncilNodeIds {
		if addr == nodeId {
			return i
		}
	}
	return AddrNotFoundInCouncilNodes
}

func (s *StakingInfo) GetStakingAmountByNodeId(nodeId common.Address) uint64 {
	i := s.GetIndexByNodeId(nodeId)
	if i != AddrNotFoundInCouncilNodes {
		return s.CouncilStakingAmounts[i]
	}
	return 0
}

func (s *StakingInfo) CalcGiniCoefficientOfValidators(validators []istanbul.Validator) {
	var stakingAmounts []uint64
	for _, val := range validators {
		i := s.GetIndexByNodeId(val.Address())
		if i != AddrNotFoundInCouncilNodes {
			stakingAmounts = append(stakingAmounts, s.CouncilStakingAmounts[i])
		}
	}
	s.Gini = calcGiniCoefficient(stakingAmounts)
}

func NewReward(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*Reward, error) {
	klaytnReward, err := contract.NewKlaytnReward(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &Reward{
		&contract.KlaytnRewardSession{
			Contract:     klaytnReward,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployReward(transactOpts *bind.TransactOpts, contractBackend bind.ContractBackend) (common.Address, *Reward, error) {

	rewardAddr, _, _, err := contract.DeployKlaytnReward(transactOpts, contractBackend)
	if err != nil {
		return rewardAddr, nil, err
	}

	reward, err := NewReward(transactOpts, rewardAddr, contractBackend)
	if err != nil {
		return rewardAddr, nil, err
	}

	return rewardAddr, reward, nil
}

type BalanceAdder interface {
	AddBalance(addr common.Address, v *big.Int)
}

// MintKLAY mints KLAY and gives the KLAY to the block proposer
func MintKLAY(b BalanceAdder, header *types.Header, gov *governance.Governance) error {

	unitPrice := big.NewInt(0)
	if r, err := gov.GetGovernanceItemAtNumber(header.Number.Uint64(), governance.GovernanceKeyMapReverse[params.UnitPrice]); err == nil {
		unitPrice.SetUint64(r.(uint64))
	} else {
		logger.Error("Couldn't get UnitPrice from governance", "err", err, "received", r)
		return err
	}

	mintingAmount := big.NewInt(0)
	if r, err := gov.GetGovernanceItemAtNumber(header.Number.Uint64(), governance.GovernanceKeyMapReverse[params.MintingAmount]); err == nil {
		mintingAmount.SetString(r.(string), 10)
	} else {
		logger.Error("Couldn't get MintingAmount from governance", "err", err, "received", r)
		return err
	}

	totalGasUsed := big.NewInt(0).SetUint64(header.GasUsed)
	totalTxFee := big.NewInt(0).Mul(totalGasUsed, unitPrice)
	blockReward := big.NewInt(0).Add(mintingAmount, totalTxFee)

	b.AddBalance(header.Rewardbase, blockReward)
	return nil
}

func isEmptyAddress(addr common.Address) bool {
	return addr == common.Address{}
}

// DistributeBlockReward distributes block reward to proposer, kirAddr and pocAddr.
func DistributeBlockReward(b BalanceAdder, header *types.Header, pocAddr common.Address, kirAddr common.Address, config *params.ChainConfig) {

	// Calculate total tx fee
	totalTxFee := common.Big0
	if config.Governance.DeferredTxFee() {
		totalGasUsed := big.NewInt(0).SetUint64(header.GasUsed)
		unitPrice := big.NewInt(0).SetUint64(config.UnitPrice)
		totalTxFee = big.NewInt(0).Mul(totalGasUsed, unitPrice)
	}

	distributeBlockReward(b, header, totalTxFee, pocAddr, kirAddr, config)
}

// distributeBlockReward mints KLAY and distribute newly minted KLAY and transaction fee to proposer, kirAddr and pocAddr.
func distributeBlockReward(b BalanceAdder, header *types.Header, totalTxFee *big.Int, pocAddr common.Address, kirAddr common.Address, config *params.ChainConfig) {
	proposer := header.Rewardbase
	rewardParams := getRewardGovernanceParameters(config, header)

	// Block reward
	blockReward := big.NewInt(0).Add(rewardParams.mintingAmount, totalTxFee)

	tmpInt := big.NewInt(0)

	tmpInt = tmpInt.Mul(blockReward, rewardParams.cnRewardRatio)
	cnReward := big.NewInt(0).Div(tmpInt, rewardParams.totalRatio)

	tmpInt = tmpInt.Mul(blockReward, rewardParams.pocRatio)
	pocIncentive := big.NewInt(0).Div(tmpInt, rewardParams.totalRatio)

	tmpInt = tmpInt.Mul(blockReward, rewardParams.kirRatio)
	kirIncentive := big.NewInt(0).Div(tmpInt, rewardParams.totalRatio)

	remaining := tmpInt.Sub(blockReward, cnReward)
	remaining = tmpInt.Sub(remaining, pocIncentive)
	remaining = tmpInt.Sub(remaining, kirIncentive)
	pocIncentive = pocIncentive.Add(pocIncentive, remaining)

	// CN reward
	b.AddBalance(proposer, cnReward)

	// Proposer gets PoC incentive and KIR incentive, if there is no PoC/KIR address.
	// PoC
	if isEmptyAddress(pocAddr) {
		pocAddr = proposer
	}
	b.AddBalance(pocAddr, pocIncentive)

	// KIR
	if isEmptyAddress(kirAddr) {
		kirAddr = proposer
	}
	b.AddBalance(kirAddr, kirIncentive)

	logger.Debug("Block reward",
		"Reward address of a proposer", proposer, "CN reward amount", cnReward,
		"PoC address", pocAddr, "Poc incentive", pocIncentive,
		"KIR address", kirAddr, "KIR incentive", kirIncentive)
}

func parseRewardRatio(ratio string) (int, int, int, error) {
	s := strings.Split(ratio, "/")
	if len(s) != 3 {
		return 0, 0, 0, errors.New("Invalid format")
	}
	cn, err1 := strconv.Atoi(s[0])
	poc, err2 := strconv.Atoi(s[1])
	kir, err3 := strconv.Atoi(s[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, errors.New("Parsing error")
	}
	return cn, poc, kir, nil
}

// Cache for parsed reward parameters from governance
type blockRewardParameters struct {
	blockNum uint64

	mintingAmount *big.Int
	cnRewardRatio *big.Int
	pocRatio      *big.Int
	kirRatio      *big.Int
	totalRatio    *big.Int
}

var blockRewardCache *blockRewardParameters
var blockRewardCacheLock sync.Mutex

// StakingCache
const (
	// TODO-Klaytn-Issue1166 Decide size of cache
	maxStakingCache   = 4 // TODO-Klaytn If you increase this value, please also improve add operation of stakingInfoCache
	chainHeadChanSize = 10
)

type stakingInfoCache struct {
	cells       map[uint64]*StakingInfo
	minBlockNum uint64
	lock        sync.RWMutex
}

var stakingCache *stakingInfoCache

func (sc *stakingInfoCache) get(blockNum uint64) *StakingInfo {
	sc.lock.RLock()
	defer sc.lock.RUnlock()

	if s, ok := sc.cells[blockNum]; ok {
		return s
	}
	return nil
}

func (sc *stakingInfoCache) add(stakingInfo *StakingInfo) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	// Assumption: stakingInfo should not be nil.

	if _, ok := sc.cells[stakingInfo.BlockNum]; ok {
		return
	}

	if len(sc.cells) < maxStakingCache {
		// empty room available
		sc.cells[stakingInfo.BlockNum] = stakingInfo
		if stakingInfo.BlockNum < sc.minBlockNum || len(sc.cells) == 1 {
			// new minBlockNum or newly inserted one is the first element
			sc.minBlockNum = stakingInfo.BlockNum
		}
		return
	}

	// evict one and insert new one
	delete(sc.cells, sc.minBlockNum)

	// update minBlockNum
	if stakingInfo.BlockNum < sc.minBlockNum {
		sc.minBlockNum = stakingInfo.BlockNum
	} else {
		min := stakingInfo.BlockNum
		for _, s := range sc.cells {
			if s.BlockNum < min {
				min = s.BlockNum
			}
		}
		sc.minBlockNum = min
	}
	sc.cells[stakingInfo.BlockNum] = stakingInfo

}

var chainHeadCh chan blockchain.ChainHeadEvent
var chainHeadSub event.Subscription
var blockchainForReward *blockchain.BlockChain

func init() {
	initStakingCache()
	allocBlockRewardCache()
}

func allocBlockRewardCache() {
	blockRewardCache = new(blockRewardParameters)

	blockRewardCache.mintingAmount = nil // We don't allocate mintingAmount, because we will use allocated mintingAmount in governance.
	blockRewardCache.cnRewardRatio = new(big.Int)
	blockRewardCache.pocRatio = new(big.Int)
	blockRewardCache.kirRatio = new(big.Int)
	blockRewardCache.totalRatio = new(big.Int)
}

// Subscribe setups a channel to listen chain head event and starts a goroutine to update staking cache.
func Subscribe(bc *blockchain.BlockChain) {
	blockchainForReward = bc
	chainHeadSub = bc.SubscribeChainHeadEvent(chainHeadCh)

	go waitHeadChain(bc)
}

func initStakingCache() {
	stakingCache = new(stakingInfoCache)
	stakingCache.cells = make(map[uint64]*StakingInfo)
	chainHeadCh = make(chan blockchain.ChainHeadEvent, chainHeadChanSize)
}

func waitHeadChain(bc *blockchain.BlockChain) {
	defer chainHeadSub.Unsubscribe()

	logger.Info("Start listening chain head event to update staking cache.")

	for {
		// A real event arrived, process interesting content
		select {
		// Handle ChainHeadEvent
		case ev := <-chainHeadCh:
			if bc.Config().Istanbul.ProposerPolicy == params.WeightedRandom && params.IsStakingUpdateInterval(ev.Block.NumberU64()) {
				blockNum := ev.Block.NumberU64()
				logger.Debug("ChainHeadEvent arrived and try to update staking cache.", "Block number", blockNum)
				if _, err := updateStakingCache(bc, blockNum); err != nil {
					logger.Error("Failed to update staking cache", "err", err)
				}
			}
		case <-chainHeadSub.Err():
			return
		}
	}
}

// GetStakingInfoFromStakingCache returns corresponding staking information for a block of blockNum.
func GetStakingInfoFromStakingCache(blockNum uint64) *StakingInfo {
	var stakingInfo *StakingInfo
	number := params.CalcStakingBlockNumber(blockNum)

	stakingInfo, err := updateStakingCache(blockchainForReward, number)
	if err != nil {
		logger.Error("Failed to get staking information", "Block number", blockNum, "number of staking block", number, "err", err)
		return nil
	}

	if stakingInfo.BlockNum != number {
		logger.Error("Invalid staking info from staking cache", "Block number", blockNum, "expected staking block number", number, "actual staking block number", stakingInfo.BlockNum)
		return nil
	}

	logger.Debug("Staking cache hit.", "Block number", blockNum, "number of staking block", number)
	return stakingInfo
}

// updateStakingCache updates staking cache with staking information of given block number.
func updateStakingCache(bc *blockchain.BlockChain, blockNum uint64) (*StakingInfo, error) {
	if cachedStakingInfo := stakingCache.get(blockNum); cachedStakingInfo != nil {
		// already updated
		return cachedStakingInfo, nil
	}

	stakingInfo, err := getStakingInfoFromAddressBook(bc, blockNum)
	if err != nil {
		return nil, err
	}

	stakingCache.add(stakingInfo)

	logger.Info("Add new staking information to staking cache", "stakingInfo", stakingInfo)

	return stakingInfo, nil
}

const (
	maxStakingLimit        = uint64(100000000000)
	DefaultGiniCoefficient = -1.0
)

var (
	maxStakingLimitBigInt = big.NewInt(0).SetUint64(maxStakingLimit)
)

func newEmptyStakingInfo(bc *blockchain.BlockChain, blockNum uint64) (*StakingInfo, error) {
	stakingInfo := &StakingInfo{
		BlockNum:              blockNum,
		CouncilNodeIds:        make([]common.Address, 0, 0),
		CouncilStakingdAddrs:  make([]common.Address, 0, 0),
		CouncilRewardAddrs:    make([]common.Address, 0, 0),
		KIRAddr:               common.Address{},
		PoCAddr:               common.Address{},
		CouncilStakingAmounts: make([]uint64, 0, 0),
		Gini:                  DefaultGiniCoefficient,
		UseGini:               false,
	}
	return stakingInfo, nil
}

func newStakingInfo(bc *blockchain.BlockChain, blockNum uint64, nodeIds []common.Address, stakingAddrs []common.Address, rewardAddrs []common.Address, KIRAddr common.Address, PoCAddr common.Address) (*StakingInfo, error) {

	// TODO-Klaytn-Issue1166 Disable all below Trace log later after all block reward implementation merged

	// Prepare
	intervalBlock := bc.GetBlockByNumber(blockNum)
	statedb, err := bc.StateAt(intervalBlock.Root())
	if err != nil {
		logger.Trace("Failed to make a state for interval block", "interval blockNum", blockNum, "err", err)
		return nil, err
	}

	// Get balance of rewardAddrs
	var stakingAmounts []uint64
	stakingAmounts = make([]uint64, len(stakingAddrs))
	for i, stakingAddr := range stakingAddrs {
		tempStakingAmount := big.NewInt(0).Div(statedb.GetBalance(stakingAddr), big.NewInt(0).SetUint64(params.KLAY))
		if tempStakingAmount.Cmp(maxStakingLimitBigInt) > 0 {
			tempStakingAmount.SetUint64(maxStakingLimit)
		}
		stakingAmounts[i] = tempStakingAmount.Uint64()
	}

	useGini := bc.Config().Governance.Reward.UseGiniCoeff
	gini := DefaultGiniCoefficient

	stakingInfo := &StakingInfo{
		BlockNum:              blockNum,
		CouncilNodeIds:        nodeIds,
		CouncilStakingdAddrs:  stakingAddrs,
		CouncilRewardAddrs:    rewardAddrs,
		KIRAddr:               KIRAddr,
		PoCAddr:               PoCAddr,
		CouncilStakingAmounts: stakingAmounts,
		Gini:                  gini,
		UseGini:               useGini,
	}
	return stakingInfo, nil
}

type uint64Slice []uint64

func (p uint64Slice) Len() int           { return len(p) }
func (p uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func calcGiniCoefficient(stakingAmount uint64Slice) float64 {
	sort.Sort(stakingAmount)

	// calculate gini coefficient
	sumOfAbsoluteDifferences := uint64(0)
	subSum := uint64(0)

	for i, x := range stakingAmount {
		temp := x*uint64(i) - subSum

		sumOfAbsoluteDifferences = sumOfAbsoluteDifferences + temp
		subSum = subSum + x
	}

	result := float64(sumOfAbsoluteDifferences) / float64(subSum) / float64(len(stakingAmount))
	result = math.Round(result*100) / 100

	return result
}

// getRewardGovernanceParameters retrieves reward parameters from governance. It also maintains a cache to reuse already parsed parameters.
func getRewardGovernanceParameters(config *params.ChainConfig, header *types.Header) *blockRewardParameters {
	blockRewardCacheLock.Lock()
	defer blockRewardCacheLock.Unlock()

	blockNum := header.Number.Uint64()

	// Cache hit condition
	// (1) blockNum is a key of cache.
	// (2) mintingAmount indicates whether cache entry is initialized or not
	// refresh at block (number -1) % epoch == 0 .
	// voting is calculated at epoch number in snapshot (which is 1 less than block header number)
	// cache refresh should be done after snapshot calculating vote.
	// so cache refresh for block header number should be 1 + epoch number
	// blockNumber cannot be 0 because this function is called by finalize() and finalize for genesis block isn't called
	if (blockNum-1)%config.Istanbul.Epoch == 0 || blockRewardCache.blockNum+config.Istanbul.Epoch < blockNum || blockRewardCache.mintingAmount == nil {
		// Cache missed or not initialized yet. Let's parse governance parameters and update cache
		cn, poc, kir, err := parseRewardRatio(config.Governance.Reward.Ratio)
		if err != nil {
			logger.Error("Error while parsing reward ratio of governance. Using default ratio", "err", err)

			cn = params.DefaultCNRewardRatio
			poc = params.DefaultPoCRewardRatio
			kir = params.DefaultKIRRewardRatio
		}

		// allocate new cache entry
		newBlockRewardCache := new(blockRewardParameters)
		newBlockRewardCache.cnRewardRatio = new(big.Int)
		newBlockRewardCache.pocRatio = new(big.Int)
		newBlockRewardCache.kirRatio = new(big.Int)
		newBlockRewardCache.totalRatio = new(big.Int)

		// update new cache entry
		if config.Governance.Reward.MintingAmount == nil {
			logger.Error("No minting amount defined in governance. Use default value.", "Default minting amount", params.DefaultMintedKLAY)
			newBlockRewardCache.mintingAmount = params.DefaultMintedKLAY
		} else {
			newBlockRewardCache.mintingAmount = config.Governance.Reward.MintingAmount
		}

		newBlockRewardCache.blockNum = blockNum
		newBlockRewardCache.cnRewardRatio.SetInt64(int64(cn))
		newBlockRewardCache.pocRatio.SetInt64(int64(poc))
		newBlockRewardCache.kirRatio.SetInt64(int64(kir))
		newBlockRewardCache.totalRatio.Add(newBlockRewardCache.cnRewardRatio, newBlockRewardCache.pocRatio)
		newBlockRewardCache.totalRatio.Add(newBlockRewardCache.totalRatio, newBlockRewardCache.kirRatio)

		// update cache
		blockRewardCache = newBlockRewardCache

		// TODO-Klaytn-RemoveLater Remove below trace later
		logger.Trace("Reward parameters updated from governance", "blockNum", newBlockRewardCache.blockNum, "minting amount", newBlockRewardCache.mintingAmount, "cn ratio", newBlockRewardCache.cnRewardRatio, "poc ratio", newBlockRewardCache.pocRatio, "kir ratio", newBlockRewardCache.kirRatio)
	}

	return blockRewardCache
}

func makeMsgToAddressBook() (*types.Transaction, error) {
	abiStr := contract.AddressBookABI
	abii, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}

	data, err := abii.Pack("getAllAddress")
	if err != nil {
		return nil, err
	}

	intrinsicGas, err := types.IntrinsicGas(data, false, true)
	if err != nil {
		return nil, err
	}

	addr := common.HexToAddress(contract.AddressBookContractAddress)

	// Create new call message
	// TODO-Klaytn-Issue1166 Decide who will be sender(i.e. from)
	msg := types.NewMessage(common.Address{}, &addr, 0, big.NewInt(0), 10000000, big.NewInt(0), data, false, intrinsicGas)

	return msg, nil
}

// addressType defined in AddressBook
const (
	addressTypeNodeID = iota
	addressTypeStakingAddr
	addressTypeRewardAddr
	addressTypePoCAddr
	addressTypeKIRAddr
)

var (
	errAddressBookEmptyResult = errors.New("got empty result, nothing to parse")
	errAddressBookIncomplete  = errors.New("incomplete node information from AddressBook")
)

func getAllAddressFromAddressBook(result []byte) ([]common.Address, []common.Address, []common.Address, common.Address, common.Address, error) {

	if result == nil {
		return nil, nil, nil, common.Address{}, common.Address{}, errAddressBookIncomplete
	}

	abiStr := contract.AddressBookABI
	abii, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, nil, nil, common.Address{}, common.Address{}, err
	}

	var (
		allTypeList    = new([]uint8)
		allAddressList = new([]common.Address)
	)
	out := &[]interface{}{
		allTypeList,
		allAddressList,
	}

	err = abii.Unpack(out, "getAllAddress", result)
	if err != nil {
		logger.Trace("abii.Unpack failed for getAllAddress in AddressBook")
		return nil, nil, nil, common.Address{}, common.Address{}, err
	}

	if len(*allTypeList) != len(*allAddressList) {
		return nil, nil, nil, common.Address{}, common.Address{}, errors.New(fmt.Sprintf("length of type list and address list differ. len(type)=%d, len(addrs)=%d", len(*allTypeList), len(*allAddressList)))
	}

	// Parse and construct node information
	nodeIds := []common.Address{}
	stakingAddrs := []common.Address{}
	rewardAddrs := []common.Address{}
	pocAddr := common.Address{}
	kirAddr := common.Address{}
	for i, addrType := range *allTypeList {
		switch addrType {
		case addressTypeNodeID:
			nodeIds = append(nodeIds, (*allAddressList)[i])
		case addressTypeStakingAddr:
			stakingAddrs = append(stakingAddrs, (*allAddressList)[i])
		case addressTypeRewardAddr:
			rewardAddrs = append(rewardAddrs, (*allAddressList)[i])
		case addressTypePoCAddr:
			pocAddr = (*allAddressList)[i]
		case addressTypeKIRAddr:
			kirAddr = (*allAddressList)[i]
		default:
			return nil, nil, nil, common.Address{}, common.Address{}, errors.New(fmt.Sprintf("invalid type from AddressBook: %d", addrType))
		}
	}

	// validate parsed node information
	if len(nodeIds) != len(stakingAddrs) ||
		len(nodeIds) != len(rewardAddrs) ||
		isEmptyAddress(pocAddr) ||
		isEmptyAddress(kirAddr) {
		// This is expected behavior when bootstrapping
		//logger.Trace("Incomplete node information from AddressBook.",
		//	"# of nodeIds", len(nodeIds),
		//	"# of stakingAddrs", len(stakingAddrs),
		//	"# of rewardAddrs", len(rewardAddrs),
		//	"PoC address", pocAddr.String(),
		//	"KIR address", kirAddr.String())

		return nil, nil, nil, common.Address{}, common.Address{}, errAddressBookIncomplete
	}

	return nodeIds, stakingAddrs, rewardAddrs, pocAddr, kirAddr, nil
}

// getStakingInfoFromAddressBook returns staking info when calling AddressBook
// succeeded. It returns an error otherwise.
func getStakingInfoFromAddressBook(bc *blockchain.BlockChain, blockNum uint64) (*StakingInfo, error) {

	var nodeIds []common.Address
	var stakingAddrs []common.Address
	var rewardAddrs []common.Address
	var KIRAddr = common.Address{}
	var PoCAddr = common.Address{}
	var err error

	if !params.IsStakingUpdateInterval(blockNum) {
		return nil, errors.New(fmt.Sprintf("not staking block number. blockNum: %d", blockNum))
	}

	// Prepare a message
	msg, err := makeMsgToAddressBook()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to make message for AddressBook. root err: %s", err))
	}

	// Prepare
	if bc == nil {
		return nil, errors.New(fmt.Sprintf("blockchain is not ready for staking info. blockNum: %d", blockNum))
	}
	chainConfig := bc.Config()
	intervalBlock := bc.GetBlockByNumber(blockNum)
	if intervalBlock == nil {
		return nil, errors.New("stateDB is not ready for staking info")
	}
	statedb, err := bc.StateAt(intervalBlock.Root())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to make a state for interval block. blockNum: %d, root err: %s", blockNum, err))
	}

	// Create a new context to be used in the EVM environment
	context := blockchain.NewEVMContext(msg, intervalBlock.Header(), bc, nil)
	evm := vm.NewEVM(context, statedb, chainConfig, &vm.Config{})

	res, gas, kerr := blockchain.ApplyMessage(evm, msg)
	logger.Trace("Call AddressBook contract", "used gas", gas, "kerr", kerr)
	err = kerr.ErrTxInvalid
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to call AddressBook contract. root err: %s", err))
	}

	nodeIds, stakingAddrs, rewardAddrs, PoCAddr, KIRAddr, err = getAllAddressFromAddressBook(res)
	if err != nil {
		if err == errAddressBookIncomplete {
			// This is expected behavior when smart contract is not setup yet.
			logger.Info("Use empty staking info instead of info from AddressBook", "reason", err)
		} else {
			logger.Error("Failed to parse result from AddressBook contract. Use empty staking info", "err", err)
		}
		return newEmptyStakingInfo(bc, blockNum)
	}

	return newStakingInfo(bc, blockNum, nodeIds, stakingAddrs, rewardAddrs, KIRAddr, PoCAddr)
}
