// Copyright 2021 The klaytn Authors
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

package state

import (
	"testing"

	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/stretchr/testify/assert"
)

func TestResetStateObjectEncoder(t *testing.T) {
	defer func() {
		resetStateObjectEncoder(stateObjEncoderDefaultWorkers, stateObjEncoderDefaultCap)
	}()

	firstChSize := 1
	secondChSize := 2
	testAcc, err := account.NewAccountWithType(account.ExternallyOwnedAccountType)
	if err != nil {
		t.Fatal("failed to create a test account", "err", err)
	}

	// reset stateObjectEncoder for test
	resetStateObjectEncoder(1, firstChSize)
	assert.Equal(t, firstChSize, cap(stateObjEncoder.tasksCh))
	assert.Equal(t, 0, len(stateObjEncoder.tasksCh))

	// getStateObjectEncoder(firstChSize) should not assign a new channel
	soe := getStateObjectEncoder(firstChSize)
	assert.Equal(t, firstChSize, cap(stateObjEncoder.tasksCh))
	soe.encode(&stateObject{account: testAcc})

	// getStateObjectEncoder(secondChSize) should assign a new channel
	soe = getStateObjectEncoder(secondChSize)
	assert.Equal(t, secondChSize, cap(stateObjEncoder.tasksCh))
	assert.Equal(t, 0, len(stateObjEncoder.tasksCh))

	soe.encode(&stateObject{account: testAcc})
	soe.encode(&stateObject{account: testAcc})
}
