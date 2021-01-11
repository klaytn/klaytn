// CookieJar - A contestant's algorithm toolbox
// Copyright (c) 2013 Peter Szilagyi. All rights reserved.
//
// CookieJar is dual licensed: use of this source code is governed by a BSD
// license that can be found in the LICENSE file. Alternatively, the CookieJar
// toolbox may be used in accordance with the terms and conditions contained
// in a signed written agreement between you and the author(s).

// This is a duplicated and slightly modified version of "gopkg.in/karalabe/cookiejar.v2/collections/prque".

// Package prque implements a priority queue data structure supporting arbitrary
// value and priorities types.
//
// If you would like to use a min-priority queue, simply put false as a invert parameter.
//
// Internally the queue is based on the standard heap package working on a
// sortable version of the block based stack.

package prque

import (
	"container/heap"
)

// Priority queue data structure. The item's type can be int, int64, uint64,
// []byte and etc.
// An type assertion error could occur if different type of items are pushed.
type Prque struct {
	cont *sstack
}

// New creates a new priority queue.
// The priority queue can have any type specified in sstack.Less.
func New(invert bool) *Prque {
	return &Prque{newSstack(invert)}
}

// Pushes a value with a given priority into the queue, expanding if necessary.
// The priority queue can have any type specified in sstack.Less.
// In a queue, same type should be pushed. Otherwise, a type assertion error
// will occur. Number types (int, int64, uint64) are different types.
func (p *Prque) Push(data interface{}, priority interface{}) {
	heap.Push(p.cont, &item{data, priority})
}

// Peek returns the value with the greatest priority but does not pop it off.
func (p *Prque) Peek() (interface{}, interface{}) {
	item := p.cont.blocks[0][0]
	return item.value, item.priority
}

// Pops the value with the greates priority off the stack and returns it.
// Currently no shrinking is done.
func (p *Prque) Pop() (interface{}, interface{}) {
	item := heap.Pop(p.cont).(*item)
	return item.value, item.priority
}

// Pops only the item from the queue, dropping the associated priority value.
func (p *Prque) PopItem() interface{} {
	return heap.Pop(p.cont).(*item).value
}

// Checks whether the priority queue is empty.
func (p *Prque) Empty() bool {
	return p.cont.Len() == 0
}

// Returns the number of element in the priority queue.
func (p *Prque) Size() int {
	return p.cont.Len()
}

// Clears the contents of the priority queue.
func (p *Prque) Reset() {
	*p = *New(p.cont.invert)
}
