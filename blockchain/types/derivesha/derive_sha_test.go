// Copyright 2022 The klaytn Authors
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

package derivesha

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"gotest.tools/assert"
)

var dummyList = types.Transactions([]*types.Transaction{
	types.NewTransaction(1, common.Address{}, big.NewInt(123), 21000, big.NewInt(25e9), nil),
})

type testGov struct{}

// Mimic the governance vote situation
var testGovSchedule = map[uint64]int{
	0: types.ImplDeriveShaOriginal,
	1: types.ImplDeriveShaOriginal,
	2: types.ImplDeriveShaOriginal,
	3: types.ImplDeriveShaSimple,
	4: types.ImplDeriveShaSimple,
	5: types.ImplDeriveShaSimple,
	6: types.ImplDeriveShaConcat,
	7: types.ImplDeriveShaConcat,
	8: types.ImplDeriveShaConcat,
}

func (e *testGov) EffectiveParams(num uint64) (*params.GovParamSet, error) {
	return params.NewGovParamSetIntMap(map[int]interface{}{
		params.DeriveShaImpl: testGovSchedule[num],
	})
}

func TestEmptyRoot(t *testing.T) {
	assert.Equal(t,
		DeriveShaOrig{}.DeriveSha(types.Transactions{}).Hex(),
		types.EmptyRootHashOriginal.Hex())
	assert.Equal(t,
		DeriveShaOrig{}.DeriveSha(types.Transactions{}).Hex(),
		"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	assert.Equal(t,
		DeriveShaSimple{}.DeriveSha(types.Transactions{}).Hex(),
		"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
	assert.Equal(t,
		DeriveShaConcat{}.DeriveSha(types.Transactions{}).Hex(),
		"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
}

func TestMuxChainConfig(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

	for implType, impl := range impls {
		InitDeriveSha(&params.ChainConfig{DeriveShaImpl: implType}, nil)
		assert.Equal(t,
			DeriveShaMux(dummyList, big.NewInt(0)),
			impl.DeriveSha(dummyList),
		)
		assert.Equal(t,
			EmptyRootHashMux(big.NewInt(0)),
			impl.DeriveSha(types.Transactions{}),
		)
	}
}

func TestMuxGovernance(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

	InitDeriveSha(
		&params.ChainConfig{DeriveShaImpl: testGovSchedule[0]},
		&testGov{})

	for num := uint64(0); num < 9; num++ {
		implType := testGovSchedule[num]
		impl := impls[implType]

		assert.Equal(t,
			DeriveShaMux(dummyList, new(big.Int).SetUint64(num)),
			impl.DeriveSha(dummyList),
		)
		assert.Equal(t,
			EmptyRootHashMux(new(big.Int).SetUint64(num)),
			impl.DeriveSha(types.Transactions{}),
		)
	}
}
