// Copyright 2019 The go-ethereum Authors
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

package p2p

import (
	"container/heap"
	"sync"

	"github.com/autonity/autonity/common/mclock"
)

type Expirable interface {
	mclock.AbsTime | uint64
}

// expHeap tracks strings and their expiry time.
type expHeap[T Expirable] []expItem[T]

// expItem is an entry in addrHistory.
type expItem[T Expirable] struct {
	item string
	exp  T
}

// nextExpiry returns the next expiry time.
func (h *expHeap[T]) nextExpiry() T {
	return (*h)[0].exp
}

// add adds an item and sets its expiry time.
func (h *expHeap[T]) add(item string, exp T) {
	heap.Push(h, expItem[T]{item, exp})
}

// contains checks whether an item is present.
func (h expHeap[T]) contains(item string) bool {
	for _, v := range h {
		if v.item == item {
			return true
		}
	}
	return false
}

// expire removes items with expiry time before 'now'.
func (h *expHeap[T]) expire(now T, onExp func(string)) {
	for h.Len() > 0 && h.nextExpiry() < now {
		item := heap.Pop(h).(expItem[T])
		if onExp != nil {
			onExp(item.item)
		}
	}
}

// heap.Interface boilerplate
func (h expHeap[T]) Len() int            { return len(h) }
func (h expHeap[T]) Less(i, j int) bool  { return h[i].exp < (h[j].exp) }
func (h expHeap[T]) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *expHeap[T]) Push(x interface{}) { *h = append(*h, x.(expItem[T])) }
func (h *expHeap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Thread-safe version of ExpHeap
type safeExpHeap[T Expirable] struct {
	expHeap[T]
	sync.RWMutex
}

// add adds an item and sets its expiry time.
func (h *safeExpHeap[T]) add(item string, exp T) {
	h.Lock()
	defer h.Unlock()
	h.expHeap.add(item, exp)
}

// contains checks whether an item is present.
func (h *safeExpHeap[T]) contains(item string) bool {
	h.RLock()
	defer h.RUnlock()
	return h.expHeap.contains(item)
}

// expire removes items with expiry time before 'now'.
func (h *safeExpHeap[T]) expire(now T, onExp func(string)) {
	h.Lock()
	defer h.Unlock()
	h.expHeap.expire(now, onExp)
}
