package fixsizecache

import (
	"sync"
	"sync/atomic"
)

const (
	defaultBuckets = 97
	defaultEntries = 5
)

type entry[K comparable, V any] struct {
	key    K
	value  V
	active bool
}

func (e *entry[K, V]) set(key K, value V) {
	e.active = true
	e.key = key
	e.value = value
}

type bucket[T comparable, V any] struct {
	sync.RWMutex
	entries []entry[T, V] // Array to store entries
	index   int
}

func newBucket[T comparable, V any](numEntries uint) *bucket[T, V] {
	return &bucket[T, V]{
		entries: make([]entry[T, V], numEntries),
	}
}

type Cache[T comparable, V any] struct {
	buckets []*bucket[T, V]
	sum32   func(T) uint
	size    atomic.Int64
}

func New[T comparable, V any](numBuckets, numEntries uint, sum32 func(T) uint) *Cache[T, V] {
	if numBuckets < 1 {
		numBuckets = defaultBuckets
	}
	if numEntries < 1 {
		numEntries = defaultEntries
	}

	c := &Cache[T, V]{sum32: sum32}
	c.buckets = make([]*bucket[T, V], numBuckets)
	for i := range c.buckets {
		c.buckets[i] = newBucket[T, V](numEntries)
	}

	return c
}

func (c *Cache[T, V]) getBucket(keysum uint) *bucket[T, V] {
	return c.buckets[keysum%uint(len(c.buckets))]
}

func (c *Cache[T, V]) Add(key T, value V) {
	sum := c.sum32(key)
	b := c.getBucket(sum)
	b.Lock()
	defer b.Unlock()

	if !b.entries[b.index].active {
		c.size.Add(1)
	}
	b.entries[b.index].set(key, value)
	b.index = (b.index + 1) % len(b.entries)
}

func (c *Cache[T, V]) Remove(key T) {
	sum := c.sum32(key)
	b := c.getBucket(sum)
	b.Lock()
	defer b.Unlock()
	for i, entry := range b.entries {
		if entry.active && entry.key == key {
			// Deactivate the entry
			b.entries[i].active = false
			c.size.Add(-1)
			if i == b.index {
				break
			}

			// Shift the subsequent entries
			for j := i; j != b.index && (j+1)%len(b.entries) != b.index; j = (j + 1) % len(b.entries) {
				nextIndex := (j + 1) % len(b.entries)
				b.entries[j] = b.entries[nextIndex]
				b.entries[nextIndex].active = false
			}

			// Update the index to point to the new position for the next entry
			b.index = (b.index - 1 + len(b.entries)) % len(b.entries)
			break
		}
	}
}

func (c *Cache[T, V]) Size() int64 {
	return c.size.Load()
}

func (c *Cache[T, V]) Get(key T) (any, bool) {
	sum := c.sum32(key)
	b := c.getBucket(sum)
	b.RLock()
	defer b.RUnlock()
	for _, entry := range b.entries {
		if entry.active && entry.key == key {
			return entry.value, true
		}
	}
	return "", false
}

func (c *Cache[T, V]) Contains(key T) bool {
	sum := c.sum32(key)
	b := c.getBucket(sum)
	b.RLock()
	defer b.RUnlock()
	for _, entry := range b.entries {
		if entry.active && entry.key == key {
			return true
		}
	}
	return false
}

func (c *Cache[T, V]) Keys() []T {
	keys := make([]T, 0, c.Size())
	for _, b := range c.buckets {
		b.RLock()
		for _, e := range b.entries {
			if e.active {
				keys = append(keys, e.key)
			}
		}
		b.RUnlock()
	}
	return keys
}

func (c *Cache[T, V]) Purge() {
	for _, b := range c.buckets {
		b.RLock()
		for i := range b.entries {
			if b.entries[i].active {
				c.size.Add(-1)
			}
			b.entries[i].active = false
		}
		b.index = 0
		b.RUnlock()
	}
}
