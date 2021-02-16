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

package state

import (
	"math"
	"runtime"

	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/statedb"
)

var stateObjEncoderDefaultWorkers = calcNumStateObjectEncoderWorkers()

func calcNumStateObjectEncoderWorkers() int {
	numWorkers := math.Ceil(float64(runtime.NumCPU()) / 4.0)
	if numWorkers > stateObjEncoderMaxWorkers {
		return stateObjEncoderMaxWorkers
	}
	return int(numWorkers)
}

const stateObjEncoderMaxWorkers = 16
const stateObjEncoderDefaultCap = 20000

var stateObjEncoder = newStateObjectEncoder(stateObjEncoderDefaultWorkers, stateObjEncoderDefaultCap)

// newStateObjectEncoder generates a stateObjectEncoder and spawns goroutines
// which encode stateObject in parallel manner.
func newStateObjectEncoder(numGoRoutines, tasksChSize int) *stateObjectEncoder {
	soe := &stateObjectEncoder{
		tasksCh: make(chan *stateObject, tasksChSize),
	}

	for i := 0; i < numGoRoutines; i++ {
		go encodeStateObject(soe.tasksCh)
	}

	return soe
}

func getStateObjectEncoder(requiredChSize int) *stateObjectEncoder {
	if requiredChSize <= cap(stateObjEncoder.tasksCh) {
		return stateObjEncoder
	}
	return resetStateObjectEncoder(stateObjEncoderDefaultWorkers, requiredChSize)
}

// resetStateObjectEncoder closes existing tasksCh and assigns a new stateObjectEncoder.
func resetStateObjectEncoder(numGoRoutines, tasksChSize int) *stateObjectEncoder {
	close(stateObjEncoder.tasksCh)
	stateObjEncoder = newStateObjectEncoder(numGoRoutines, tasksChSize)
	return stateObjEncoder
}

// stateObjectEncoder handles tasksCh and resultsCh
// to distribute the tasks and gather the results.
type stateObjectEncoder struct {
	tasksCh chan *stateObject
}

func (soe *stateObjectEncoder) encode(so *stateObject) {
	soe.tasksCh <- so
}

// encodeStateObject encodes the given stateObject and generates its hashKey and hexKey.
func encodeStateObject(tasksCh <-chan *stateObject) {
	for stateObj := range tasksCh {
		data, err := rlp.EncodeToBytes(stateObj)
		if err != nil {
			stateObj.encoded.Store(&encodedData{err: err})
			continue
		}
		addr := stateObj.Address()
		hashKey, hexKey := statedb.GetHashAndHexKey(addr[:])
		stateObj.encoded.Store(&encodedData{data: data, trieHashKey: hashKey, trieHexKey: hexKey})
	}
}
