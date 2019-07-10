package reward

import (
	"fmt"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"math"
	"math/big"
	"sort"
	"strings"
)

const (
	AddrNotFoundInCouncilNodes = -1
	maxStakingLimit            = uint64(100000000000)
	DefaultGiniCoefficient     = -1.0
)

var (
	maxStakingLimitBigInt = big.NewInt(0).SetUint64(maxStakingLimit)
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

func newEmptyStakingInfo(blockNum uint64) (*StakingInfo, error) {
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

func newStakingInfo(bc *blockchain.BlockChain, helper governanceHelper, blockNum uint64, nodeIds []common.Address, stakingAddrs []common.Address, rewardAddrs []common.Address, KIRAddr common.Address, PoCAddr common.Address) (*StakingInfo, error) {

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

	var useGini bool

	if res, err := helper.GetItemAtNumberByIntKey(blockNum, params.UseGiniCoeff); err != nil {
		logger.Trace("Failed to get useGiniCoeff from governance", "blockNum", blockNum, "err", err)
		return nil, err
	} else {
		useGini = res.(bool)
	}
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

type uint64Slice []uint64

func (p uint64Slice) Len() int           { return len(p) }
func (p uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func CalcGiniCoefficient(stakingAmount uint64Slice) float64 {
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
