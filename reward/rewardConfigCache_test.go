package reward

import (
	"errors"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

type testGovernance struct {
	epoch         uint64
	mintingAmount string
	ratio         string
	unitPrice     uint64
	useGiniCoeff  bool
	deferredTxFee bool
}

func newDefaultTestGovernance() *testGovernance {
	return &testGovernance{
		epoch:         604800,
		mintingAmount: "9600000000000000000",
		ratio:         "34/54/12",
		unitPrice:     25000000000,
		useGiniCoeff:  true,
		deferredTxFee: true,
	}
}

func newTestGovernance(epoch uint64, mintingAmount string, ratio string, unitPrice uint64, useGiniCoeff bool, deferredTxFee bool) *testGovernance {
	return &testGovernance{
		epoch:         epoch,
		mintingAmount: mintingAmount,
		ratio:         ratio,
		unitPrice:     unitPrice,
		useGiniCoeff:  useGiniCoeff,
		deferredTxFee: deferredTxFee,
	}
}

func (governance *testGovernance) Epoch() uint64 {
	return governance.epoch
}

func (governance *testGovernance) GetItemAtNumberByIntKey(num uint64, key int) (interface{}, error) {
	switch key {
	case params.MintingAmount:
		return governance.mintingAmount, nil
	case params.Ratio:
		return governance.ratio, nil
	case params.UnitPrice:
		return governance.unitPrice, nil
	case params.Epoch:
		return governance.epoch, nil
	default:
		return nil, errors.New("Unhandled key on testGovernance")
	}
}

func (governance *testGovernance) DeferredTxFee() bool {
	return governance.deferredTxFee
}

func (governance *testGovernance) setTestGovernance(epoch uint64, mintingAmount string, ratio string, unitprice uint64, useGiniCoeff bool, deferredTxFee bool) {
	governance.epoch = epoch
	governance.mintingAmount = mintingAmount
	governance.ratio = ratio
	governance.unitPrice = unitprice
	governance.useGiniCoeff = useGiniCoeff
	governance.deferredTxFee = deferredTxFee
}

func TestRewardConfigCache_parseRewardRatio(t *testing.T) {
	testCases := []struct {
		s       string
		cn      int
		poc     int
		kir     int
		success bool
	}{
		{"34/54/12", 34, 54, 12, true},
		{"3/3/3", 3, 3, 3, true},
		{"10/20/30", 10, 20, 30, true},
		{"///", 0, 0, 0, false},
		{"1//", 0, 0, 0, false},
		{"/1/", 0, 0, 0, false},
		{"//1", 0, 0, 0, false},
		{"1/2/3/4/", 0, 0, 0, false},
		{"3.3/3.3/3.3", 0, 0, 0, false},
		{"a/b/c", 0, 0, 0, false},
	}
	rewardConfigCache := newRewardConfigCache(newDefaultTestGovernance())

	for i := 0; i < len(testCases); i++ {
		cn, poc, kir, error := rewardConfigCache.parseRewardRatio(testCases[i].s)

		if (error == nil) != testCases[i].success || cn != testCases[i].cn ||
			poc != testCases[i].poc || kir != testCases[i].kir {
			t.Errorf("test case %v fail. The result is different", testCases[i].s)
			t.Errorf("The parsed cn. Result : %v, Expected : %v", cn, testCases[i].cn)
			t.Errorf("The parsed poc. Result : %v, Expected : %v", poc, testCases[i].poc)
			t.Errorf("The parsed kir. Result : %v, Expected : %v", kir, testCases[i].kir)
		}
	}
}

func TestRewardConfigCache_newRewardConfig(t *testing.T) {
	testCases := []struct {
		testGovernance testGovernance
		result         rewardConfig
	}{
		{
			testGovernance{
				epoch:         604800,
				mintingAmount: "9600000000000000000",
				ratio:         "34/54/12",
				unitPrice:     25000000000,
				useGiniCoeff:  true,
				deferredTxFee: true,
			},
			rewardConfig{
				blockNum:      1,
				mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
				cnRatio:       big.NewInt(0).SetInt64(34),
				pocRatio:      big.NewInt(0).SetInt64(54),
				kirRatio:      big.NewInt(0).SetInt64(12),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(25000000000),
			},
		},
		{
			testGovernance{
				epoch:         30,
				mintingAmount: "10000",
				ratio:         "50/30/20",
				unitPrice:     50000000000,
				useGiniCoeff:  true,
				deferredTxFee: false,
			},
			rewardConfig{
				blockNum:      2,
				mintingAmount: big.NewInt(0).SetInt64(10000),
				cnRatio:       big.NewInt(0).SetInt64(50),
				pocRatio:      big.NewInt(0).SetInt64(30),
				kirRatio:      big.NewInt(0).SetInt64(20),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(50000000000),
			},
		},
		{
			testGovernance{
				epoch:         3000,
				mintingAmount: "100000000",
				ratio:         "10/35/55",
				unitPrice:     1500000000,
				useGiniCoeff:  false,
				deferredTxFee: true,
			},
			rewardConfig{
				blockNum:      3,
				mintingAmount: big.NewInt(0).SetInt64(100000000),
				cnRatio:       big.NewInt(0).SetInt64(10),
				pocRatio:      big.NewInt(0).SetInt64(35),
				kirRatio:      big.NewInt(0).SetInt64(55),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(1500000000),
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		testGovernance := &testCases[i].testGovernance
		rewardConfigCache := newRewardConfigCache(testGovernance)

		rewardConfig, error := rewardConfigCache.newRewardConfig(uint64(i) + 1)

		if error != nil {
			t.Errorf("error has occurred err : %v", error)
		}

		expectedResult := &testCases[i].result
		assert.Equal(t, expectedResult.blockNum, rewardConfig.blockNum)
		assert.Equal(t, expectedResult.mintingAmount, rewardConfig.mintingAmount)
		assert.Equal(t, expectedResult.cnRatio, rewardConfig.cnRatio)
		assert.Equal(t, expectedResult.pocRatio, rewardConfig.pocRatio)
		assert.Equal(t, expectedResult.kirRatio, rewardConfig.kirRatio)
		assert.Equal(t, expectedResult.totalRatio, rewardConfig.totalRatio)
		assert.Equal(t, expectedResult.unitPrice, rewardConfig.unitPrice)
	}
}

func TestRewardConfigCache_add(t *testing.T) {
	testCases := []rewardConfig{
		{
			blockNum:      1,
			mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
			cnRatio:       big.NewInt(0).SetInt64(34),
			pocRatio:      big.NewInt(0).SetInt64(54),
			kirRatio:      big.NewInt(0).SetInt64(12),
			totalRatio:    big.NewInt(0).SetInt64(100),
			unitPrice:     big.NewInt(0).SetInt64(25000000000),
		},
		{
			blockNum:      2,
			mintingAmount: big.NewInt(0).SetInt64(10000),
			cnRatio:       big.NewInt(0).SetInt64(50),
			pocRatio:      big.NewInt(0).SetInt64(30),
			kirRatio:      big.NewInt(0).SetInt64(20),
			totalRatio:    big.NewInt(0).SetInt64(100),
			unitPrice:     big.NewInt(0).SetInt64(50000000000),
		},
		{
			blockNum:      3,
			mintingAmount: big.NewInt(0).SetInt64(100000000),
			cnRatio:       big.NewInt(0).SetInt64(10),
			pocRatio:      big.NewInt(0).SetInt64(35),
			kirRatio:      big.NewInt(0).SetInt64(55),
			totalRatio:    big.NewInt(0).SetInt64(100),
			unitPrice:     big.NewInt(0).SetInt64(1500000000),
		},
	}

	testGovernance := newDefaultTestGovernance()
	rewardConfigCache := newRewardConfigCache(testGovernance)

	for i := 0; i < len(testCases); i++ {
		rewardConfigCache.add(uint64(i)+1, &testCases[i])
		assert.Equal(t, i+1, rewardConfigCache.cache.Len())
	}
}

func TestRewardConfigCache_add_sameNumber(t *testing.T) {
	rewardConfig := rewardConfig{
		blockNum:      1,
		mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
		cnRatio:       big.NewInt(0).SetInt64(34),
		pocRatio:      big.NewInt(0).SetInt64(54),
		kirRatio:      big.NewInt(0).SetInt64(12),
		totalRatio:    big.NewInt(0).SetInt64(100),
		unitPrice:     big.NewInt(0).SetInt64(25000000000),
	}

	testGovernance := newDefaultTestGovernance()
	rewardConfigCache := newRewardConfigCache(testGovernance)

	rewardConfigCache.add(1, &rewardConfig)
	rewardConfigCache.add(1, &rewardConfig)
	assert.Equal(t, 1, rewardConfigCache.cache.Len())
}

func TestRewardConfigCache_get_exist(t *testing.T) {
	testCases := []struct {
		blockNumber   uint64
		epoch         uint64
		mintingAmount string
		ratio         string
		unitprice     uint64
		useGiniCoeff  bool
		deferredTxFee bool
		result        rewardConfig
	}{
		{
			blockNumber:   1,
			epoch:         604800,
			mintingAmount: "9600000000000000000",
			ratio:         "34/54/12",
			unitprice:     25000000000,
			useGiniCoeff:  true,
			deferredTxFee: true,
			result: rewardConfig{
				blockNum:      0,
				mintingAmount: big.NewInt(0).SetUint64(9600000000000000000),
				cnRatio:       big.NewInt(0).SetInt64(34),
				pocRatio:      big.NewInt(0).SetInt64(54),
				kirRatio:      big.NewInt(0).SetInt64(12),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(25000000000),
			},
		},
		{
			blockNumber:   604805,
			epoch:         604800,
			mintingAmount: "9600000000000000",
			ratio:         "40/25/35",
			unitprice:     250000000,
			useGiniCoeff:  true,
			deferredTxFee: false,
			result: rewardConfig{
				blockNum:      604800,
				mintingAmount: big.NewInt(0).SetUint64(9600000000000000),
				cnRatio:       big.NewInt(0).SetInt64(40),
				pocRatio:      big.NewInt(0).SetInt64(25),
				kirRatio:      big.NewInt(0).SetInt64(35),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(250000000),
			},
		},
		{
			blockNumber:   1210000,
			epoch:         604800,
			mintingAmount: "100000000000000000",
			ratio:         "34/33/33",
			unitprice:     100000000000,
			useGiniCoeff:  false,
			deferredTxFee: true,
			result: rewardConfig{
				blockNum:      1209600,
				mintingAmount: big.NewInt(0).SetUint64(100000000000000000),
				cnRatio:       big.NewInt(0).SetInt64(34),
				pocRatio:      big.NewInt(0).SetInt64(33),
				kirRatio:      big.NewInt(0).SetInt64(33),
				totalRatio:    big.NewInt(0).SetInt64(100),
				unitPrice:     big.NewInt(0).SetInt64(100000000000),
			},
		},
	}

	testGovernance := newDefaultTestGovernance()
	rewardConfigCache := newRewardConfigCache(testGovernance)

	for i := 0; i < len(testCases); i++ {
		testGovernance.setTestGovernance(testCases[i].epoch, testCases[i].mintingAmount, testCases[i].ratio, testCases[i].unitprice, testCases[i].useGiniCoeff, testCases[i].deferredTxFee)
		blockNumber := testCases[i].blockNumber
		if blockNumber%testCases[i].epoch == 0 {
			blockNumber -= testCases[i].epoch
		} else {
			blockNumber -= (blockNumber % testCases[i].epoch)
		}
		rewardConfig, _ := rewardConfigCache.newRewardConfig(blockNumber)
		rewardConfigCache.add(blockNumber, rewardConfig)
	}
	for i := 0; i < len(testCases); i++ {
		rewardConfig, err := rewardConfigCache.get(testCases[i].blockNumber)
		if err != nil {
			t.Errorf("error has occurred. err : %v", err)
		}
		assert.Equal(t, testCases[i].result.blockNum, rewardConfig.blockNum)
		assert.Equal(t, testCases[i].result.mintingAmount, rewardConfig.mintingAmount)
		assert.Equal(t, testCases[i].result.cnRatio, rewardConfig.cnRatio)
		assert.Equal(t, testCases[i].result.pocRatio, rewardConfig.pocRatio)
		assert.Equal(t, testCases[i].result.kirRatio, rewardConfig.kirRatio)
		assert.Equal(t, testCases[i].result.totalRatio, rewardConfig.totalRatio)
		assert.Equal(t, testCases[i].result.unitPrice, rewardConfig.unitPrice)
	}
}
