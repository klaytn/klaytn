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

package reward

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stakingInfoTestCase struct {
	stakingInfo          *StakingInfo
	expectedConsolidated *ConsolidatedStakingInfo
	expectedAmounts      map[common.Address]uint64
}

var stakingInfoTestCases = generateStakingInfoTestCases()

func generateStakingInfoTestCases() []stakingInfoTestCase {
	var (
		n1 = common.HexToAddress("0x8aD8F547fa00f58A8c4fb3B671Ee5f1A75bA028a")
		n2 = common.HexToAddress("0xB2AAda7943919e82143324296987f6091F3FDC9e")
		n3 = common.HexToAddress("0xD95c70710f07A3DaF7ae11cFBa10c789da3564D0")
		n4 = common.HexToAddress("0xC704765db1d21C2Ea6F7359dcB8FD5233DeD16b5")

		s1 = common.HexToAddress("0x4dd324F9821485caE941640B32c3Bcf1fA6E93E6")
		s2 = common.HexToAddress("0x0d5Df5086B5f86f748dFaed5779c3f862C075B1f")
		s3 = common.HexToAddress("0xD3Ff05f00491571E86A3cc8b0c320aA76D7413A5")
		s4 = common.HexToAddress("0x11EF8e61d10365c2ECAe0E95b5fFa9ed4D68d64f")

		r1 = common.HexToAddress("0x241c793A9AD555f52f6C3a83afe6178408796ab2")
		r2 = common.HexToAddress("0x79b427Fb77077A9716E08D049B0e8f36Abfc8E2E")
		r3 = common.HexToAddress("0x62E47d858bf8513fc401886B94E33e7DCec2Bfb7")
		r4 = common.HexToAddress("0xf275f9f4c0d375F9E3E50370f93b504A1e45dB09")

		kcf = common.HexToAddress("0x136807B12327a8AfF9831F09617dA1B9D398cda2")
		kff = common.HexToAddress("0x46bA8F7538CD0749e572b2631F9FB4Ce3653AFB8")

		a0 uint64 = 0
		aL uint64 = 1000000  // less than minstaking
		aM uint64 = 2000000  // exactly minstaking (params.DefaultMinimumStake)
		a1 uint64 = 10000000 // above minstaking. Using 1,2,4,8 to uniquely spot errors
		a2 uint64 = 20000000
		a3 uint64 = 40000000
		a4 uint64 = 80000000
	)
	if aM != params.DefaultMinimumStake.Uint64() {
		panic("broken test assumption")
	}

	return []stakingInfoTestCase{
		// Empty
		{
			stakingInfo: newEmptyStakingInfo(0),
			expectedConsolidated: &ConsolidatedStakingInfo{
				nodes:     make([]consolidatedNode, 0),
				nodeIndex: make(map[common.Address]int),
			},
			expectedAmounts: make(map[common.Address]uint64),
		},

		// 1 entry
		{
			stakingInfo: &StakingInfo{
				BlockNum:              86400,
				CouncilNodeAddrs:      []common.Address{n1},
				CouncilStakingAddrs:   []common.Address{s1},
				CouncilRewardAddrs:    []common.Address{r1},
				KCFAddr:               kcf,
				KFFAddr:               kff,
				UseGini:               true,
				Gini:                  0.00,
				CouncilStakingAmounts: []uint64{a1},
			},
			expectedConsolidated: &ConsolidatedStakingInfo{
				nodes: []consolidatedNode{
					{[]common.Address{n1}, []common.Address{s1}, r1, a1},
				},
				nodeIndex: map[common.Address]int{n1: 0},
			},
			expectedAmounts: map[common.Address]uint64{n1: a1},
		},

		// Ordinary 4-entry info
		{
			stakingInfo: &StakingInfo{
				BlockNum:              2 * 86400,
				CouncilNodeAddrs:      []common.Address{n1, n2, n3, n4},
				CouncilStakingAddrs:   []common.Address{s1, s2, s3, s4},
				CouncilRewardAddrs:    []common.Address{r1, r2, r3, r4},
				KCFAddr:               kcf,
				KFFAddr:               kff,
				UseGini:               true,
				Gini:                  0.38, // Gini(10, 20, 40, 80)
				CouncilStakingAmounts: []uint64{a1, a2, a3, a4},
			},
			expectedConsolidated: &ConsolidatedStakingInfo{
				nodes: []consolidatedNode{
					{[]common.Address{n1}, []common.Address{s1}, r1, a1},
					{[]common.Address{n2}, []common.Address{s2}, r2, a2},
					{[]common.Address{n3}, []common.Address{s3}, r3, a3},
					{[]common.Address{n4}, []common.Address{s4}, r4, a4},
				},
				nodeIndex: map[common.Address]int{n1: 0, n2: 1, n3: 2, n4: 3},
			},
			expectedAmounts: map[common.Address]uint64{n1: a1, n2: a2, n3: a3, n4: a4},
		},

		// 4-entry with common reward addrs
		{
			stakingInfo: &StakingInfo{
				BlockNum:              3 * 86400,
				CouncilNodeAddrs:      []common.Address{n1, n2, n3, n4},
				CouncilStakingAddrs:   []common.Address{s1, s2, s3, s4},
				CouncilRewardAddrs:    []common.Address{r1, r2, r1, r2}, // r1 and r2 used twice each
				KCFAddr:               kcf,
				KFFAddr:               kff,
				UseGini:               true,
				Gini:                  0.17, // Gini(50, 100)
				CouncilStakingAmounts: []uint64{a1, a2, a3, a4},
			},
			expectedConsolidated: &ConsolidatedStakingInfo{
				nodes: []consolidatedNode{
					{[]common.Address{n1, n3}, []common.Address{s1, s3}, r1, a1 + a3}, // n1 & n3
					{[]common.Address{n2, n4}, []common.Address{s2, s4}, r2, a2 + a4}, // n2 & n4
				},
				nodeIndex: map[common.Address]int{n1: 0, n2: 1, n3: 0, n4: 1},
			},
			expectedAmounts: map[common.Address]uint64{n1: a1 + a3, n2: a2 + a4, n3: a1 + a3, n4: a2 + a4},
		},

		// 4-entry with less-than-minstaking amounts
		{
			stakingInfo: &StakingInfo{
				BlockNum:              4 * 86400,
				CouncilNodeAddrs:      []common.Address{n1, n2, n3, n4},
				CouncilStakingAddrs:   []common.Address{s1, s2, s3, s4},
				CouncilRewardAddrs:    []common.Address{r1, r2, r3, r4},
				KCFAddr:               kcf,
				KFFAddr:               kff,
				UseGini:               true,
				Gini:                  0.41,                     // Gini(20, 2)
				CouncilStakingAmounts: []uint64{a2, aM, aL, a0}, // aL and a0 should be ignored in Gini calculation
			},
			expectedConsolidated: &ConsolidatedStakingInfo{
				nodes: []consolidatedNode{
					{[]common.Address{n1}, []common.Address{s1}, r1, a2},
					{[]common.Address{n2}, []common.Address{s2}, r2, aM},
					{[]common.Address{n3}, []common.Address{s3}, r3, aL},
					{[]common.Address{n4}, []common.Address{s4}, r4, a0},
				},
				nodeIndex: map[common.Address]int{n1: 0, n2: 1, n3: 2, n4: 3},
			},
			expectedAmounts: map[common.Address]uint64{n1: a2, n2: aM, n3: aL, n4: a0},
		},
	}
}

func TestStakingInfo_GetIndexByNodeAddress(t *testing.T) {
	testdata := []common.Address{
		common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"),
		common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"),
		common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"),
		common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"),
	}
	testCases := []struct {
		address common.Address
		index   int
		err     error
	}{
		{common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"), 0, nil},
		{common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"), 1, nil},
		{common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"), 2, nil},
		{common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"), 3, nil},
		{common.StringToAddress("0x027AbB8c9f952cfFf01B1707fF14E2CB5D439502"), AddrNotFoundInCouncilNodes, ErrAddrNotInStakingInfo},
	}

	stakingInfo := newEmptyStakingInfo(0)
	stakingInfo.CouncilNodeAddrs = testdata

	for i := 0; i < len(testCases); i++ {
		result, err := stakingInfo.GetIndexByNodeAddress(testCases[i].address)
		assert.Equal(t, testCases[i].index, result)
		assert.Equal(t, testCases[i].err, err)
	}
}

func TestStakingInfo_GetStakingAmountByNodeId(t *testing.T) {
	testdata := struct {
		address       []common.Address
		stakingAmount []uint64
	}{
		[]common.Address{
			common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"),
			common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"),
			common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"),
			common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"),
		},
		[]uint64{
			100, 200, 300, 400,
		},
	}
	testCases := []struct {
		address       common.Address
		stakingAmount uint64
		err           error
	}{
		{common.StringToAddress("0xB55e5986b972Be438b4A91d6e8726aA50AD55EDc"), 100, nil},
		{common.StringToAddress("0xaDfc427080B4a66b5a629cd633d48C5d734572cA"), 200, nil},
		{common.StringToAddress("0x994daB8EB6f3FaE044cC0c9a0AB1A038e136b0B6"), 300, nil},
		{common.StringToAddress("0xD527822212Fded72c5fE89f46281d5355BD58235"), 400, nil},
		{common.StringToAddress("0x027AbB8c9f952cfFf01B1707fF14E2CB5D439502"), 0, ErrAddrNotInStakingInfo},
	}

	stakingInfo := newEmptyStakingInfo(0)
	stakingInfo.CouncilNodeAddrs = testdata.address
	stakingInfo.CouncilStakingAmounts = testdata.stakingAmount

	for i := 0; i < len(testCases); i++ {
		result, err := stakingInfo.GetStakingAmountByNodeId(testCases[i].address)
		assert.Equal(t, testCases[i].stakingAmount, result)
		assert.Equal(t, testCases[i].err, err)
	}
}

func TestStakingInfo_String(t *testing.T) {
	// No information loss in String() -> Unmarshal() round trip
	for _, testcase := range stakingInfoTestCases {
		resultStr := testcase.stakingInfo.String()
		t.Logf("%s", resultStr)

		resultByteArr := []byte(resultStr)
		resultStakingInfo := &StakingInfo{}
		err := json.Unmarshal(resultByteArr, resultStakingInfo)
		assert.NoError(t, err)

		assert.Equal(t, testcase.stakingInfo, resultStakingInfo)
	}
}

func TestCalcGiniCoefficient(t *testing.T) {
	testCase := []struct {
		testdata []float64
		result   float64
	}{
		{[]float64{1, 1, 1}, 0.0},
		{[]float64{0, 8, 0, 0, 0}, 0.8},
		{[]float64{5, 4, 3, 2, 1}, 0.27},
	}

	for i := 0; i < len(testCase); i++ {
		result := CalcGiniCoefficient(testCase[i].testdata)
		assert.Equal(t, testCase[i].result, result)
	}
}

func TestGiniReflectToExpectedCCO(t *testing.T) {
	testCase := []struct {
		ccoToken        []float64
		beforeReflected []float64
		adjustment      []float64
		afterReflected  []float64
	}{
		{
			[]float64{
				66666667, 233333333, 5000000, 5000000, 5000000,
				77777778, 5000000, 33333333, 20000000, 16666667,
				10000000, 5000000, 5000000, 5000000, 5000000,
				5000000, 5000000, 5000000, 5000000, 5000000,
				5000000,
			},
			[]float64{13, 44, 1, 1, 1, 15, 1, 6, 4, 3, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]float64{42612, 89426, 9202, 9202, 9202, 46682, 9202, 28275, 20900, 18762, 13868, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202},
			[]float64{11, 23, 2, 2, 2, 12, 2, 7, 5, 5, 4, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		},
		{
			[]float64{
				400000000, 233333333, 233333333, 150000000, 108333333,
				83333333, 66666667, 33333333, 20000000, 16666667,
				10000000, 5000000, 5000000, 5000000, 5000000,
				5000000, 5000000, 5000000, 5000000, 5000000,
				5000000,
			},
			[]float64{28, 17, 17, 11, 8, 6, 5, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]float64{123020, 89426, 89426, 68853, 56793, 48627, 42612, 28275, 20900, 18762, 13868, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202, 9202},
			[]float64{18, 13, 13, 10, 8, 7, 6, 4, 3, 3, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
	}
	for i := 0; i < len(testCase); i++ {
		stakingInfo := newEmptyStakingInfo(uint64(1))

		weights := make([]float64, len(testCase[i].ccoToken))
		tokenListToCalcGini := make([]float64, len(testCase[i].ccoToken))
		totalAmount := 0.0
		for j := 0; j < len(testCase[i].ccoToken); j++ {
			totalAmount += float64(testCase[i].ccoToken[j])
			tokenListToCalcGini[j] = testCase[i].ccoToken[j]
		}

		for j := 0; j < len(testCase[i].ccoToken); j++ {
			weights[j] = math.Round(float64(testCase[i].ccoToken[j]) * 100 / totalAmount)
			if weights[j] < 1 {
				weights[j] = 1
			}
			if weights[j] != testCase[i].beforeReflected[j] {
				t.Errorf("normal weight is incorrect. result : %v expected : %v", weights[j], testCase[i].beforeReflected[j])
			}
		}

		stakingAmountsGiniReflected := make([]float64, len(testCase[i].ccoToken))
		totalAmountGiniReflected := 0.0
		stakingInfo.Gini = CalcGiniCoefficient(tokenListToCalcGini)

		for j := 0; j < len(stakingAmountsGiniReflected); j++ {
			stakingAmountsGiniReflected[j] = math.Round(math.Pow(float64(testCase[i].ccoToken[j]), 1.0/(1+stakingInfo.Gini)))
			totalAmountGiniReflected += stakingAmountsGiniReflected[j]
		}

		for j := 0; j < len(testCase[i].ccoToken); j++ {
			if stakingAmountsGiniReflected[j] != testCase[i].adjustment[j] {
				t.Errorf("staking amount reflected gini is different. result : %v expected : %v", stakingAmountsGiniReflected[j], testCase[i].adjustment[j])
			}
		}

		for j := 0; j < len(testCase[i].ccoToken); j++ {
			stakingAmountsGiniReflected[j] = math.Round(stakingAmountsGiniReflected[j] * 100 / totalAmountGiniReflected)
			if stakingAmountsGiniReflected[j] != testCase[i].afterReflected[j] {
				t.Errorf("weight reflected gini is different. result : %v expected : %v", stakingAmountsGiniReflected[j], testCase[i].afterReflected[j])
			}
		}
	}
}

// TestStakingInfoJSON tests marshalling and unmarshalling StakingInfo
// StakingInfo is marshaled before storing to DB.
func TestStakingInfoJSON(t *testing.T) {
	// No information loss in json.Marshal() -> json.Unmarshal() round trip
	for _, testcase := range stakingInfoTestCases {
		src := testcase.stakingInfo

		s, err := json.Marshal(src)
		require.Nil(t, err)

		dst := new(StakingInfo)
		err = json.Unmarshal(s, dst)
		require.Nil(t, err)

		assert.Equal(t, src, dst)
	}
}

func TestConsolidatedStakingInfo(t *testing.T) {
	for _, testcase := range stakingInfoTestCases {
		expected := testcase.expectedConsolidated
		c := testcase.stakingInfo.GetConsolidatedStakingInfo()

		// Test ConsolidatedStakingInfo
		assert.Equal(t, expected.nodes, c.nodes)
		assert.Equal(t, expected.nodeIndex, c.nodeIndex)

		// Test CalcGiniCoefficient()
		expectedGini := testcase.stakingInfo.Gini
		gini := c.CalcGiniCoefficientMinStake(2000000)
		assert.Equal(t, expectedGini, gini)

		// Test GetConsolidatedNode()
		for addr, expectedAmount := range testcase.expectedAmounts {
			node := c.GetConsolidatedNode(addr)
			require.NotNil(t, node)
			assert.Equal(t, expectedAmount, node.StakingAmount)
		}
	}
}

// oldStakingInfo is a legacy of StakingInfo providing backward-compatibility.
// Since json tags of StakingInfo were changed, a node may fail to unmarshal stored data without this.
// oldStakingInfo's field names are the same with StakingInfo's names, but json tag is different.
type oldStakingInfo struct {
	BlockNum              uint64           `json:"BlockNum"`
	CouncilNodeAddrs      []common.Address `json:"CouncilNodeAddrs"`
	CouncilStakingAddrs   []common.Address `json:"CouncilStakingAddrs"`
	CouncilRewardAddrs    []common.Address `json:"CouncilRewardAddrs"`
	KCFAddr               common.Address   `json:"KIRAddr"` // KIRAddr -> KCFAddr from v1.10.2
	KFFAddr               common.Address   `json:"PoCAddr"` // PoCAddr -> KFFAddr from v1.10.2
	UseGini               bool             `json:"UseGini"`
	Gini                  float64          `json:"Gini"`
	CouncilStakingAmounts []uint64         `json:"CouncilStakingAmounts"`
}

var oldInfo = oldStakingInfo{
	2880,
	[]common.Address{
		common.HexToAddress("0x159ae5ccda31b77475c64d88d4499c86f77b7ecc"),
		common.HexToAddress("0x181deb121304b0430d99328ff1a9122df9f09d7f"),
		common.HexToAddress("0x324ec8f2681cd73642cc55057970540a1f4393e0"),
		common.HexToAddress("0x11191029025d3fcd21001746f949b25c6e8435cc"),
	},
	[]common.Address{
		common.HexToAddress("0x70e051c46ea76b9af9977407bb32192319907f9e"),
		common.HexToAddress("0xe4a0c3821a2711758306ed57c2f4900aa9ddbb3d"),
		common.HexToAddress("0xf3ba3a33b3bf7cf2085890315b41cc788770feb3"),
		common.HexToAddress("0x9285a85777d0ae7e12bee3ffd7842908b2295f45"),
	},
	[]common.Address{
		common.HexToAddress("0xd155d4277c99fa837c54a37a40a383f71a3d082a"),
		common.HexToAddress("0x2b8cc0ca62537fa5e49dce197acc8a15d3c5d4a8"),
		common.HexToAddress("0x7d892f470ecde693f52588dd0cfe46c3d26b6219"),
		common.HexToAddress("0xa0f7354a0cef878246820b6caa19d2bdef74a0cc"),
	},
	common.HexToAddress("0x673003e5f9a852d3dc85b83d16ef62d45497fb96"),
	common.HexToAddress("0x576dc0c2afeb1661da3cf53a60e76dd4e32c7ab1"),
	false,
	-1,
	[]uint64{5000000, 5000000, 5000000, 5000000},
}

var newInfo = StakingInfo{
	oldInfo.BlockNum,
	[]common.Address{
		common.HexToAddress("0x70e051c46ea76b9af9977407bb32192319907f9e"),
		common.HexToAddress("0xe4a0c3821a2711758306ed57c2f4900aa9ddbb3d"),
		common.HexToAddress("0xf3ba3a33b3bf7cf2085890315b41cc788770feb3"),
		common.HexToAddress("0x11191029025d3fcd21001746f949b25c6e8435cc"),
	},
	[]common.Address{
		common.HexToAddress("0x7d892f470ecde693f52588dd0cfe46c3d26b6219"),
		common.HexToAddress("0x2b8cc0ca62537fa5e49dce197acc8a15d3c5d4a8"),
		common.HexToAddress("0xa0f7354a0cef878246820b6caa19d2bdef74a0cc"),
		common.HexToAddress("0x576dc0c2afeb1661da3cf53a60e76dd4e32c7ab1"),
	},
	[]common.Address{
		common.HexToAddress("0xd155d4277c99fa837c54a37a40a383f71a3d082a"),
		common.HexToAddress("0x159ae5ccda31b77475c64d88d4499c86f77b7ecc"),
		common.HexToAddress("0x181deb121304b0430d99328ff1a9122df9f09d7f"),
		common.HexToAddress("0x673003e5f9a852d3dc85b83d16ef62d45497fb96"),
	},
	common.HexToAddress("0x324ec8f2681cd73642cc55057970540a1f4393e0"),
	common.HexToAddress("0x9285a85777d0ae7e12bee3ffd7842908b2295f45"),
	false,
	0.3,
	[]uint64{15000000, 4000000, 25000000, 35000000},
}

// TestGetStakingInfoFromDB tests whether the node can read oldStakingInfo and StakingInfo data or not.
func TestGetStakingInfoFromDB(t *testing.T) {
	oldStakingManager := GetStakingManager()
	defer SetTestStakingManager(oldStakingManager)

	for _, info := range []interface{}{oldInfo, newInfo} {
		// reset database
		SetTestStakingManagerWithDB(database.NewMemoryDBManager())

		infoBytes, err := json.Marshal(info)
		if err != nil {
			t.Fatal(err)
		}

		stakingManager.stakingInfoDB.WriteStakingInfo(oldInfo.BlockNum, infoBytes)
		retrievedInfo, err := getStakingInfoFromDB(oldInfo.BlockNum)
		if err != nil {
			t.Fatal(err)
		}

		vInfo := reflect.ValueOf(info)
		vRetriedInfo := reflect.ValueOf(*retrievedInfo)

		assert.Equal(t, vInfo.NumField(), vRetriedInfo.NumField())
		for i := 0; i < vInfo.NumField(); i++ {
			assert.Equal(t, vInfo.Field(i).Interface(), vRetriedInfo.Field(i).Interface())
		}
	}
}

// TestStakingInfo_MarshalJSON tests marshal/unmarshal staking info data.
func TestStakingInfo_MarshalJSON(t *testing.T) {
	// old marshalled data, new unmarshal method
	{
		oldInfoByte, err := json.Marshal(oldInfo)
		if err != nil {
			t.Fatal(err)
		}

		var unmarshalled StakingInfo
		if err := json.Unmarshal(oldInfoByte, &unmarshalled); err != nil {
			t.Fatal(err)
		}
		checkStakingInfoValues(t, oldInfo, unmarshalled)
	}

	// new marshalled data, old unmarshal method
	{
		newInfoByte, err := json.Marshal(newInfo)
		if err != nil {
			t.Fatal(err)
		}

		var unmarshalled oldStakingInfo
		if err := json.Unmarshal(newInfoByte, &unmarshalled); err != nil {
			t.Fatal(err)
		}
		checkStakingInfoValues(t, newInfo, unmarshalled)
	}

	// new marshalled data, new unmarshal method
	{
		newInfoByte, err := json.Marshal(newInfo)
		if err != nil {
			t.Fatal(err)
		}

		var unmarshalled StakingInfo
		if err := json.Unmarshal(newInfoByte, &unmarshalled); err != nil {
			t.Fatal(err)
		}
		checkStakingInfoValues(t, newInfo, unmarshalled)
	}
}

func checkStakingInfoValues(t *testing.T, info interface{}, stakingInfo interface{}) {
	vOld := reflect.ValueOf(info)
	vNew := reflect.ValueOf(stakingInfo)
	assert.Equal(t, vOld.NumField(), vNew.NumField())

	for i := 0; i < vOld.NumField(); i++ {
		field := reflect.TypeOf(info).Field(i).Name
		assert.Equal(t, vOld.FieldByName(field).Interface(), vNew.FieldByName(field).Interface())
	}
}
