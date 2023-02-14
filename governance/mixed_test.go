package governance

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/common"
	govcontract "github.com/klaytn/klaytn/contracts/gov"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMixedEngine(t *testing.T, config *params.ChainConfig) (*MixedEngine, *bind.TransactOpts, *backends.SimulatedBackend, *govcontract.GovParam) {
	config.IstanbulCompatibleBlock = common.Big0
	config.LondonCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0

	accounts, sim, addr, contract := prepareSimulatedContract(t)

	config.Governance.GovParamContract = addr

	db := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	e := NewMixedEngine(config, db)
	require.NotNil(t, e)
	require.NotNil(t, e.headerGov)

	e.SetBlockchain(sim.BlockChain())

	return e, accounts[0], sim, contract
}

func newTestMixedEngineNoContractEngine(t *testing.T, config *params.ChainConfig) *MixedEngine {
	// disable ContractEngine
	config.KoreCompatibleBlock = new(big.Int).SetUint64(0xffffffff)

	db := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})

	e := NewMixedEngine(config, db)
	headerGov := e.headerGov // to manipulate internal fields

	require.NotNil(t, e)
	require.NotNil(t, headerGov)

	return e
}

// Without ContractGov, Check that
//   - From a fresh MixedEngine instance, Params() and ParamsAt(0) returns the
//     initial config value.
func TestMixedEngine_Header_New(t *testing.T) {
	valueA := uint64(0x11)

	config := getTestConfig()
	config.Istanbul.SubGroupSize = valueA
	e := newTestMixedEngineNoContractEngine(t, config)

	// Params() should work even before explicitly calling UpdateParams().
	// For instance in cn.New().
	pset := e.Params()
	assert.Equal(t, valueA, pset.CommitteeSize())

	pset, err := e.ParamsAt(0)
	assert.Nil(t, err)
	assert.Equal(t, valueA, pset.CommitteeSize())
}

// Without ContractGov, Check that
// - after UpdateParams(), Params() returns the new value
func TestMixedEngine_Header_Params(t *testing.T) {
	valueA := uint64(0x11)
	valueB := uint64(0x22)

	config := getTestConfig()
	config.Governance.KIP71.GasTarget = valueA
	e := newTestMixedEngineNoContractEngine(t, config)
	assert.Equal(t, valueA, e.Params().GasTarget())

	items := e.Params().StrMap()
	items["kip71.gastarget"] = valueB
	gset := NewGovernanceSet()
	gset.Import(items)
	err := e.headerGov.WriteGovernance(e.Params().Epoch(), NewGovernanceSet(), gset)
	assert.Nil(t, err)
	err = e.UpdateParams(e.Params().Epoch() * 2)
	assert.Nil(t, err)

	assert.Equal(t, valueB, e.Params().GasTarget())
	// check if config is updated as well
	assert.Equal(t, valueB, config.Governance.KIP71.GasTarget)
}

// Before Kore hardfork (i.e., without ContractGov), check that
// - after DB is written at [n - epoch], ParamsAt(n+1) returns the new value
// - ParamsAt(n+1) == ReadGovernance(n)
func TestMixedEngine_Header_ParamsAt(t *testing.T) {
	valueA := uint64(0x11)
	valueB := uint64(0x22)

	config := getTestConfig()
	config.Istanbul.Epoch = 30
	config.Istanbul.SubGroupSize = valueA
	e := newTestMixedEngineNoContractEngine(t, config)

	// Write to database. Note that we must use gov.WriteGovernance(), not db.WriteGovernance()
	// The reason is that gov.ReadGovernance() depends on the caches, and that
	// gov.WriteGovernance() sets idxCache accordingly, whereas db.WriteGovernance don't
	items := e.Params().StrMap()
	items["istanbul.committeesize"] = valueB
	gset := NewGovernanceSet()
	gset.Import(items)
	e.headerGov.WriteGovernance(30, NewGovernanceSet(), gset)

	testcases := []struct {
		num   uint64
		value uint64
	}{
		{0, valueA},
		{30, valueA},
		{31, valueA},
		{59, valueA},
		{60, valueB},
		{61, valueB},
	}
	for _, tc := range testcases {
		// Check that e.ParamsAt() == tc
		pset, err := e.ParamsAt(tc.num + 1)
		assert.Nil(t, err)
		assert.Equal(t, tc.value, pset.CommitteeSize())

		// Check that headerGov.ReadGovernance() == tc
		_, strMap, err := e.headerGov.ReadGovernance(tc.num)
		assert.Nil(t, err)
		pset, err = params.NewGovParamSetStrMap(strMap)
		assert.Nil(t, err)
		assert.Equal(t, tc.value, pset.CommitteeSize())
	}
}

// TestMixedEngine_Params tests if Params() conforms to the fallback mechanism
func TestMixedEngine_Params(t *testing.T) {
	var (
		valueA      = uint64(0xa)
		valueB      = uint64(0xbb)
		valueC      = uint64(0xcccccc)
		valueCBytes = []byte{0xcc, 0xcc, 0xcc}
	)
	config := getTestConfig()
	config.Governance.KIP71.GasTarget = valueA
	config.Istanbul.Epoch = 4

	// 1. fallback to ChainConfig because headerGov and contractGov don't have the param
	e, owner, sim, contract := newTestMixedEngine(t, config)
	require.Equal(t, valueA, e.Params().GasTarget(), "fallback to ChainConfig failed")

	// 2. fallback to headerGov because contractGov doesn't have the param
	delta := NewGovernanceSet()
	delta.SetValue(params.GasTarget, valueB)
	e.headerGov.WriteGovernance(4, e.headerGov.currentSet, delta)
	for e.headerGov.blockChain.CurrentBlock().NumberU64() < 7 {
		sim.Commit()
	}
	err := e.UpdateParams(e.headerGov.blockChain.CurrentBlock().NumberU64())
	assert.Nil(t, err)
	require.Equal(t, valueB, e.Params().GasTarget(), "fallback to headerGov failed")

	// 3. use contractGov
	_, err = contract.SetParamIn(owner, "kip71.gastarget", true, valueCBytes, big.NewInt(1))
	assert.Nil(t, err)

	sim.Commit() // mine SetParamIn

	err = e.UpdateParams(e.headerGov.blockChain.CurrentBlock().NumberU64())
	assert.Nil(t, err)
	require.Equal(t, valueC, e.Params().GasTarget(), "fallback to contractGov failed")
}

// TestMixedEngine_ParamsAt tests if ParamsAt() returns correct values
// given headerBlock and contractBlock;
//
//	at headerBlock, params are inserted to DB via WriteGovernance()
//	at contractBlock, params are inserted to GovParam via SetParamIn() contract call.
//
// valueA is set in ChainConfig
// valueB is set in DB
// valueC is set in GovParam contract
//
//	chainConfig     headerBlock    contractBlock       now
//
// Block |---------------|--------------|---------------|
//
// ............valueA          valueB         valueC
func TestMixedEngine_ParamsAt(t *testing.T) {
	var (
		name        = "kip71.gastarget"
		valueA      = uint64(0xa)
		valueB      = uint64(0xbb)
		valueC      = uint64(0xcccccc)
		valueCBytes = []byte{0xcc, 0xcc, 0xcc}
	)

	config := getTestConfig()
	config.Istanbul.Epoch = 1

	// set initial value
	config.Governance.KIP71.GasTarget = valueA

	e, owner, sim, contract := newTestMixedEngine(t, config)

	// write minimal params for test to headerGov
	// note that mainnet will have all parameters in the headerGov db
	headerBlock := sim.BlockChain().CurrentBlock().NumberU64()
	e.headerGov.db.WriteGovernance(map[string]interface{}{
		name:                          valueB,
		"governance.govparamcontract": config.Governance.GovParamContract,
	}, headerBlock)
	err := e.UpdateParams(headerBlock)
	assert.Nil(t, err)

	// forward a few blocks
	for i := 0; i < 3; i++ {
		sim.Commit()
	}

	// write to contractGov
	_, err = contract.SetParamIn(owner, name, true, valueCBytes, big.NewInt(1))
	assert.Nil(t, err)

	sim.Commit() // mine SetParamIn
	contractBlock := sim.BlockChain().CurrentHeader().Number.Uint64()

	// forward a few blocks
	for i := 0; i < 3; i++ {
		sim.Commit()
	}

	now := sim.BlockChain().CurrentHeader().Number.Uint64()
	for i := uint64(0); i < now; i++ {
		pset, err := e.ParamsAt(i)
		assert.Nil(t, err)

		val, ok := pset.Get(params.GasTarget)
		assert.True(t, ok)

		var expected uint64

		switch {
		case i <= headerBlock:
			expected = valueA
		case i <= contractBlock:
			expected = valueB
		default:
			expected = valueC
		}

		assert.Equal(t, expected, val,
			"ParamsAt(%d) failed (headerBlock=%d contractBlock=%d)",
			i, headerBlock, contractBlock)
	}
}
