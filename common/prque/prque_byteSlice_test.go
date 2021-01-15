// CookieJar - A contestant's algorithm toolbox
// Copyright (c) 2013 Peter Szilagyi. All rights reserved.
//
// CookieJar is dual licensed: use of this source code is governed by a BSD
// license that can be found in the LICENSE file. Alternatively, the CookieJar
// toolbox may be used in accordance with the terms and conditions contained
// in a signed written agreement between you and the author(s).

package prque

import (
	"math/rand"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
)

func TestPrqueByteSlice(t *testing.T) {
	// Generate a batch of random data and a specific priority order
	size := 16 * blockSize
	data := make([]int, size)
	prio := make([][]byte, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Int()
		prio[i] = common.MakeRandomBytes(256)
	}
	queue := NewByteSlice()
	for rep := 0; rep < 2; rep++ {
		// Fill a priority queue with the above data
		for i := 0; i < size; i++ {
			queue.Push(data[i], prio[i])
			if queue.Size() != i+1 {
				t.Errorf("queue size mismatch: have %v, want %v.", queue.Size(), i+1)
			}
		}
		// Create a map the values to the priorities for easier verification
		dict := make(map[uint64]int)
		for i := 0; i < size; i++ {
			dict[common.ByteBigEndianToUint64(prio[i])] = data[i]
		}
		// Pop out the elements in priority order and verify them
		prevPrio := uint64(math.MaxUint64)
		for !queue.Empty() {
			val, prio := queue.Pop()
			if common.ByteBigEndianToUint64(prio) > prevPrio {
				t.Errorf("invalid priority order: %v after %v.", prio, prevPrio)
			}
			prevPrio = common.ByteBigEndianToUint64(prio)
			if val != dict[prevPrio] {
				t.Errorf("push/pop mismatch: have %v, want %v.", val, dict[prevPrio])
			}
			delete(dict, prevPrio)
		}
		if len(dict) != 0 {
			t.Errorf("items left in dict: dict %v.", dict)
		}
	}
}

func TestPrqueByteSliceInverted(t *testing.T) {
	// Generate a batch of random data and a specific priority order
	size := 16 * blockSize
	data := make([]int, size)
	prio := make([][]byte, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Int()
		prio[i] = common.MakeRandomBytes(256)
	}
	queue := NewByteSliceInverted()
	for rep := 0; rep < 2; rep++ {
		// Fill a priority queue with the above data
		for i := 0; i < size; i++ {
			queue.Push(data[i], prio[i])
			if queue.Size() != i+1 {
				t.Errorf("queue size mismatch: have %v, want %v.", queue.Size(), i+1)
			}
		}
		// Create a map the values to the priorities for easier verification
		dict := make(map[uint64]int)
		for i := 0; i < size; i++ {
			dict[common.ByteBigEndianToUint64(prio[i])] = data[i]
		}
		// Pop out the elements in priority order and verify them
		prevPrio := uint64(0)
		for !queue.Empty() {
			val, prio := queue.Pop()
			if common.ByteBigEndianToUint64(prio) < prevPrio {
				t.Errorf("invalid priority order: %v after %v.", prio, prevPrio)
			}
			prevPrio = common.ByteBigEndianToUint64(prio)
			if val != dict[prevPrio] {
				t.Errorf("push/pop mismatch: have %v, want %v.", val, dict[prevPrio])
			}
			delete(dict, prevPrio)
		}
		if len(dict) != 0 {
			t.Errorf("items left in dict: dict %v.", dict)
		}
	}
}

func BenchmarkPrqueByteSlicePush(b *testing.B) {
	// Create some initial data
	data := make([]int, b.N)
	prio := make([][]byte, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Int()
		prio[i] = common.MakeRandomBytes(256)
	}
	// Execute the benchmark
	b.ResetTimer()
	queue := NewByteSlice()
	for i := 0; i < len(data); i++ {
		queue.Push(data[i], prio[i])
	}
}

func BenchmarkPrqueByteSlicePop(b *testing.B) {
	// Create some initial data
	data := make([]int, b.N)
	prio := make([][]byte, b.N)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Int()
		prio[i] = common.MakeRandomBytes(256)
	}
	queue := NewByteSlice()
	for i := 0; i < len(data); i++ {
		queue.Push(data[i], prio[i])
	}
	// Execute the benchmark
	b.ResetTimer()
	for !queue.Empty() {
		queue.Pop()
	}
}
