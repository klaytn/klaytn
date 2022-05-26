package governance

import (
	"testing"

	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestMixedEngine(t *testing.T, config *params.ChainConfig) (*MixedEngine, database.DBManager, *Governance) {
	config.Istanbul.Epoch = 30
	// TODO: if ContractEngine is added, make sure to disable it in this test.

	db := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})

	e := NewMixedEngine(config, db)
	defaultGov := e.defaultGov.(*Governance) // to manipulate internal fields

	require.NotNil(t, e)
	require.NotNil(t, defaultGov)

	return e, db, defaultGov
}

// Without ContractGov, Check that
// - From a fresh MixedEngine instance, Params() and ParamsAt(0) returns the
//   initial config value.
func TestMixedEngine_Header_New(t *testing.T) {
	valueA := uint64(0x11)

	config := getTestConfig()
	config.Istanbul.SubGroupSize = valueA
	e, _, _ := newTestMixedEngine(t, config)

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
	config.Istanbul.SubGroupSize = valueA
	e, _, defaultGov := newTestMixedEngine(t, config)

	defaultGov.currentSet.SetValue(params.CommitteeSize, valueB)
	err := e.UpdateParams()
	assert.Nil(t, err)

	pset := e.Params()
	assert.Equal(t, valueB, pset.CommitteeSize())
}

// Without ContractGov, Check that
// - after DB is written at [n - epoch], ParamsAt(n) returns the new value
// - ParamsAt(n) == ReadGovernance(n)
func TestMixedEngine_Header_ParamsAt(t *testing.T) {
	valueA := uint64(0x11)
	valueB := uint64(0x22)

	config := getTestConfig()
	config.Istanbul.SubGroupSize = valueA
	e, _, defaultGov := newTestMixedEngine(t, config)

	// Write to database. Note that we must use gov.WriteGovernance(), not db.WriteGovernance()
	// The reason is that gov.ReadGovernance() depends on the caches, and that
	// gov.WriteGovernance() sets idxCache accordingly, whereas db.WriteGovernance don't
	items := e.Params().StrMap()
	items["istanbul.committeesize"] = valueB
	gset := NewGovernanceSet()
	gset.Import(items)
	defaultGov.WriteGovernance(30, NewGovernanceSet(), gset)

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
		pset, err := e.ParamsAt(tc.num)
		assert.Nil(t, err)
		assert.Equal(t, tc.value, pset.CommitteeSize())

		// Check that defaultGov.ReadGovernance() == tc
		// and by extension defaultGov.ReadGovernance() == e.ParamsAt().
		_, strMap, err := defaultGov.ReadGovernance(tc.num)
		assert.Nil(t, err)
		pset, err = params.NewGovParamSetStrMap(strMap)
		assert.Nil(t, err)
		assert.Equal(t, tc.value, pset.CommitteeSize())
	}
}
