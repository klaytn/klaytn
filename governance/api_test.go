package governance

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

type testBlockChain struct {
	num    uint64
	config *params.ChainConfig
}

func newTestBlockchain(config *params.ChainConfig) *testBlockChain {
	return &testBlockChain{
		config: config,
	}
}

func newTestGovernanceApi() *PublicGovernanceAPI {
	config := params.CypressChainConfig
	config.Governance.KIP71 = params.GetDefaultKIP71Config()
	govApi := NewGovernanceAPI(NewMixedEngine(config, database.NewMemoryDBManager()))
	govApi.governance.SetNodeAddress(common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"))
	bc := newTestBlockchain(config)
	govApi.governance.SetBlockchain(bc)
	return govApi
}

func TestUpperBoundBaseFeeSet(t *testing.T) {
	govApi := newTestGovernanceApi()

	curLowerBoundBaseFee := govApi.governance.Params().LowerBoundBaseFee()
	// unexpected case : upperboundbasefee < lowerboundbasefee
	invalidUpperBoundBaseFee := curLowerBoundBaseFee - 100
	_, err := govApi.Vote("kip71.upperboundbasefee", invalidUpperBoundBaseFee)
	assert.Equal(t, err, errInvalidUpperBound)
}

func TestLowerBoundFeeSet(t *testing.T) {
	govApi := newTestGovernanceApi()

	curUpperBoundBaseFee := govApi.governance.Params().UpperBoundBaseFee()
	// unexpected case : upperboundbasefee < lowerboundbasefee
	invalidLowerBoundBaseFee := curUpperBoundBaseFee + 100
	_, err := govApi.Vote("kip71.lowerboundbasefee", invalidLowerBoundBaseFee)
	assert.Equal(t, err, errInvalidLowerBound)
}

func TestGetRewards(t *testing.T) {
	type expected = map[int]uint64
	type strMap = map[string]interface{}
	type override struct {
		num    int
		strMap strMap
	}
	type testcase struct {
		length   int // total number of blocks to simulate
		override []override
		expected expected
	}

	var (
		mintAmount = uint64(1)
		koreBlock  = uint64(9)
		epoch      = 3
		latestNum  = rpc.BlockNumber(-1)
		proposer   = common.HexToAddress("0x0000000000000000000000000000000000000000")
		config     = getTestConfig()
	)

	testcases := []testcase{
		{
			12,
			[]override{
				{
					3,
					strMap{
						"reward.mintingamount": "2",
					},
				},
				{
					6,
					strMap{
						"reward.mintingamount": "3",
					},
				},
			},
			map[int]uint64{
				1:  1,
				2:  1,
				3:  1,
				4:  1,
				5:  1,
				6:  1,
				7:  2, // 2 is minted from now
				8:  2,
				9:  3, // 3 is minted from now
				10: 3,
				11: 3,
				12: 3,
				13: 3,
			},
		},
	}

	for _, tc := range testcases {
		config.Governance.Reward.MintingAmount = new(big.Int).SetUint64(mintAmount)
		config.Istanbul.Epoch = uint64(epoch)
		config.KoreCompatibleBlock = new(big.Int).SetUint64(koreBlock)

		bc := newTestBlockchain(config)

		dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})

		e := NewMixedEngine(config, dbm)
		e.SetBlockchain(bc)
		e.UpdateParams(bc.CurrentBlock().NumberU64())

		// write initial gov items and overrides to database
		pset, _ := params.NewGovParamSetChainConfig(config)
		gset := NewGovernanceSet()
		gset.Import(pset.StrMap())
		e.headerGov.WriteGovernance(0, NewGovernanceSet(), gset)
		for _, o := range tc.override {
			override := NewGovernanceSet()
			override.Import(o.strMap)
			e.headerGov.WriteGovernance(uint64(o.num), gset, override)
		}

		govKlayApi := NewGovernanceKlayAPI(e, bc)

		for num := 1; num <= tc.length; num++ {
			bc.SetBlockNum(uint64(num))

			rewardSpec, err := govKlayApi.GetRewards(&latestNum)
			assert.Nil(t, err)

			minted := new(big.Int).SetUint64(tc.expected[num])
			expectedRewardSpec := &reward.RewardSpec{
				Minted:   minted,
				TotalFee: common.Big0,
				BurntFee: common.Big0,
				Proposer: minted,
				Stakers:  common.Big0,
				Kgf:      common.Big0,
				Kir:      common.Big0,
				Rewards: map[common.Address]*big.Int{
					proposer: minted,
				},
			}
			assert.Equal(t, expectedRewardSpec, rewardSpec, "wrong at block %d", num)
		}
	}
}

func (bc *testBlockChain) Engine() consensus.Engine                    { return nil }
func (bc *testBlockChain) GetHeader(common.Hash, uint64) *types.Header { return nil }
func (bc *testBlockChain) GetHeaderByNumber(val uint64) *types.Header {
	return &types.Header{
		Number: new(big.Int).SetUint64(val),
	}
}
func (bc *testBlockChain) GetBlockByNumber(num uint64) *types.Block         { return nil }
func (bc *testBlockChain) StateAt(root common.Hash) (*state.StateDB, error) { return nil, nil }
func (bc *testBlockChain) Config() *params.ChainConfig {
	return bc.config
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlockWithHeader(bc.CurrentHeader())
}

func (bc *testBlockChain) CurrentHeader() *types.Header {
	return &types.Header{
		Number: new(big.Int).SetUint64(bc.num),
	}
}

func (bc *testBlockChain) SetBlockNum(num uint64) {
	bc.num = num
}
