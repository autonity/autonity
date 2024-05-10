package fixsizecache

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultBuckets = 97
	defaultEntries = 5
)

type entry[K comparable, V any] struct {
	key       K
	value     V
	timestamp time.Time
	active    bool
}

func (e *entry[K, V]) set(key K, value V) {
	e.active = true
	e.key = key
	e.value = value
	e.timestamp = time.Now()
}

type bucket[T comparable, V any] struct {
	sync.RWMutex
	entries []entry[T, V] // Array to store entries
	head    int
	tail    int
	full    bool
}

func newBucket[T comparable, V any](numEntries uint) *bucket[T, V] {
	return &bucket[T, V]{
		entries: make([]entry[T, V], numEntries),
	}
}

type Cache[T comparable, V any] struct {
	buckets      []*bucket[T, V]
	sum32        func(T) uint
	size         atomic.Int64
	evictionTime time.Duration
}

func New[T comparable, V any](numBuckets, numEntries uint, evictionDuration time.Duration, sum32 func(T) uint) *Cache[T, V] {
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

	if evictionDuration > 0 {
		go c.startClearing()
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
	if b.full {
		b.tail = (b.tail + 1) % len(b.entries)
	}
	if !b.entries[b.head].active {
		c.size.Add(1)
	}
	b.entries[b.head].set(key, value)
	b.head = (b.head + 1) % len(b.entries)
	b.full = b.head == b.tail
}

func (c *Cache[T, V]) Remove(key T) {
	sum := c.sum32(key)
	b := c.getBucket(sum)
	b.Lock()
	defer b.Unlock()
	for _, entry := range b.entries {
		if entry.key == key {
			entry.active = false
			c.size.Add(-1)
			break
		}
	}
}

func (c *Cache[T, V]) evictEntries(size *atomic.Int64) {
	for _, b := range c.buckets {
		b.Lock()
		var i int
		if !b.full {
			for i = b.tail; ; i = (i + 1) % len(b.entries) {
				if b.entries[i].active && time.Since(b.entries[i].timestamp) > c.evictionTime {
					b.entries[i].active = false
					size.Add(-1)
				}
				if i == b.head {
					break
				}
			}
		} else {
			for i = (b.tail - 1 + len(b.entries)) % len(b.entries); ; i = (i - 1 + len(b.entries)) % len(b.entries) {
				if b.entries[i].active && time.Since(b.entries[i].timestamp) > c.evictionTime {
					b.entries[i].active = false
					size.Add(-1)
				} else {
					break
				}
				if i == b.tail {
					break
				}
				if i == b.head {
					b.full = false
				}
			}
		}
		b.tail = i
		b.Unlock()
	}
}

// startClearing periodically clears entries from all shards.
func (c *Cache[T, V]) startClearing() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		c.evictEntries(&c.size)
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
