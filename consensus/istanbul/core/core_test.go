// Modifications Copyright 2022 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from quorum/consensus/istanbul/core/core_test.go (2019/09/30).
// Modified and improved for the klaytn development.

package core

import (
	"fmt"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
)

// What quorum size Q do we need in Byzantine setting?
//  * Liveness: Q <= N - f
//      As in non-Byzantine case, failed nodes might not reply
//  * Safety: Quorum intersection must contain one non-faulty node
//      Idea: out of f+1 nodes, at most one can be faulty
//      Hence:  2Q - N > f    (since f could be malicious)
//  So: N + f < 2Q <= 2(N - f)

func TestCore_QuorumSize(t *testing.T) {
	validatorAddrs, _ := genValidators(1)
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}
	defer istCore.Stop()

	valSet := istCore.valSet
	for i := 0; i <= 100; i++ {
		valSet.AddValidator(common.StringToAddress(fmt.Sprint(i)))
		valSet.SetSubGroupSize(valSet.Size())
		if 2*valSet.QuorumSize() <= int(valSet.Size())+valSet.F() || 2*valSet.QuorumSize() > 2*(int(valSet.Size())-valSet.F()) {
			t.Errorf("quorumSize constraint failed, expected value (2*QuorumSize > Size+F && 2*QuorumSize > 2*(Size-F)) to qs:%v, f: %v, size+f:%v, be:%v, got: %v, for size: %v", valSet.QuorumSize(), valSet.F(), int(valSet.Size())+valSet.F(), true, false, valSet.Size())
		}
	}
}
