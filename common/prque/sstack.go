// CookieJar - A contestant's algorithm toolbox
// Copyright (c) 2013 Peter Szilagyi. All rights reserved.
//
// CookieJar is dual licensed: use of this source code is governed by a BSD
// license that can be found in the LICENSE file. Alternatively, the CookieJar
// toolbox may be used in accordance with the terms and conditions contained
// in a signed written agreement between you and the author(s).

// This is a duplicated and slightly modified version of "gopkg.in/karalabe/cookiejar.v2/collections/prque".

package prque

import (
	"bytes"
)

// The size of a block of data
const blockSize = 4096

// A prioritized item in the sorted stack.
//
// Note: priorities can "wrap around" the int64 range, a comes before b if (a.priority - b.priority) > 0.
// The difference between the lowest and highest priorities in the queue at any point should be less than 2^63.
type item struct {
	value    interface{}
	priority interface{}
}

// Internal sortable stack data structure. Implements the Push and Pop ops for
// the stack (heap) functionality and the Len, Less and Swap methods for the
// sortability requirements of the heaps.
// To support a new type, add a code in Less().
type sstack struct {
	size     int
	capacity int
	offset   int
	reverse  bool // reverse the result of Less()

	blocks [][]*item
	active []*item
}

// Creates a new, empty stack.
func newSstack(reverse bool) *sstack {
	active := make([]*item, blockSize)
	return &sstack{
		size:     0,
		capacity: blockSize,
		offset:   0,
		reverse:  reverse,
		blocks:   [][]*item{active},
		active:   active,
	}
}

// Pushes a value onto the stack, expanding it if necessary. Required by
// heap.Interface.
func (s *sstack) Push(data interface{}) {
	if s.size == s.capacity {
		s.active = make([]*item, blockSize)
		s.blocks = append(s.blocks, s.active)
		s.capacity += blockSize
		s.offset = 0
	} else if s.offset == blockSize {
		s.active = s.blocks[s.size/blockSize]
		s.offset = 0
	}
	s.active[s.offset] = data.(*item)
	s.offset++
	s.size++
}

// Pops a value off the stack and returns it. Currently no shrinking is done.
// Required by heap.Interface.
func (s *sstack) Pop() (res interface{}) {
	s.size--
	s.offset--
	if s.offset < 0 {
		s.offset = blockSize - 1
		s.active = s.blocks[s.size/blockSize]
	}
	res, s.active[s.offset] = s.active[s.offset], nil
	return
}

// Returns the length of the stack. Required by sort.Interface.
func (s *sstack) Len() int {
	return s.size
}

// Compares the priority of two elements of the stack (higher is first).
// Required by sort.Interface.
// To support a new type, add a switch case.
func (s *sstack) Less(i, j int) bool {
	iIntfPriority := s.blocks[i/blockSize][i%blockSize].priority
	jIntfPriority := s.blocks[j/blockSize][j%blockSize].priority

	var result bool
	switch iPriority := iIntfPriority.(type) {
	case int:
		// If an type assertion error occurred, check types in Push().
		// Same type should be pushed.
		jPriority := jIntfPriority.(int)
		result = iPriority > jPriority
	case int64:
		jPriority := jIntfPriority.(int64)
		result = iPriority > jPriority
	case uint64:
		jPriority := jIntfPriority.(uint64)
		result = iPriority > jPriority
	case []byte:
		jPriority := jIntfPriority.([]byte)
		result = bytes.Compare(iPriority, jPriority) > 0
	}

	if s.reverse {
		return !result
	}
	return result
}

// Swaps two elements in the stack. Required by sort.Interface.
func (s *sstack) Swap(i, j int) {
	ib, io, jb, jo := i/blockSize, i%blockSize, j/blockSize, j%blockSize
	s.blocks[ib][io], s.blocks[jb][jo] = s.blocks[jb][jo], s.blocks[ib][io]
}

// Resets the stack, effectively clearing its contents.
func (s *sstack) Reset() {
	*s = *newSstack(s.reverse)
}
