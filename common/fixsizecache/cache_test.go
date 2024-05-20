package fixsizecache

import (
	"crypto/rand"
	"crypto/sha256"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/stretchr/testify/require"
)

const (
	evictionDuration = 0
	bucketCount      = 1997
	maxEntries       = 20
)

func generateKey(i int) [32]byte {
	key := "key" + strconv.Itoa(i%(bucketCount*maxEntries))
	k := sha256.Sum256([]byte(key))
	return k
}

func generateKeyRand() [32]byte {
	var key [32]byte
	rand.Read(key[:])
	return key
}

func setupFullCache(c *Cache[[32]byte, bool]) {
	for i := 0; i < bucketCount*maxEntries; i++ {
		c.Add(generateKey(i), true)
	}
}

func addEntries(c *Cache[[32]byte, bool], num int) {
	for i := 0; i < num; i++ {
		c.Add(generateKey(i), true)
	}
}

func cleanupFullCache(c *Cache[[32]byte, bool]) { //nolint
	for i := 0; i < bucketCount*maxEntries; i++ {
		c.Remove(generateKey(i))
	}
}

func setupFullLRU(c *lru.LRU[[32]byte, bool]) {
	for i := 0; i < bucketCount*maxEntries; i++ {
		c.Add(generateKey(i), true)
	}
}

func cleanupFullLRU(c *lru.LRU[[32]byte, bool]) { //nolint
	for i := 0; i < bucketCount*maxEntries; i++ {
		c.Remove(generateKey(i))
	}
}

// simpleHash is a simple hash function for testing purposes.
func simpleHash(key string) uint {
	h := sha256.Sum256([]byte(key))
	return uint(h[0])
}

func TestCacheConcurrentGetSet(t *testing.T) {
	c := New[[32]byte, bool](bucketCount, maxEntries, evictionDuration, HashKey[[32]byte])
	setupFullCache(c)

	var wg sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := generateKeyRand()
			c.Add(key, true)
			if _, ok := c.Get(key); !ok {
				t.Errorf("Expected to get key %v, but didn't", key)
			}
		}()
	}
	wg.Wait()
}

func TestCache_Eviction(t *testing.T) {
	c := New[[32]byte, bool](bucketCount, maxEntries, 2*time.Second, HashKey[[32]byte])
	setupFullCache(c)
	//t.Log("full cache size", c.Size(), "expected", bucketCount*maxEntries)
	time.Sleep(4 * time.Second)
	require.Equal(t, int64(0), c.Size())
	addEntries(c, bucketCount*2)
	time.Sleep(4 * time.Second)
	require.Equal(t, int64(0), c.Size())
}

func TestSmallCache(t *testing.T) {
	c := New[string, bool](1, 1, evictionDuration, simpleHash)
	c.Add("key1", true)
	require.True(t, c.Contains("key1"))
	c.Add("key2", true)
	require.Equal(t, int64(1), c.Size())
	c.Contains("key2")
	require.True(t, c.Contains("key2"))
	require.Equal(t, int64(1), c.Size())
}

// TestKeysEmptyCache tests the Keys method on an empty cache.
func TestKeysEmptyCache(t *testing.T) {
	cache := New[string, int](10, 10, 5*time.Minute, simpleHash)
	keys := cache.Keys()
	if len(keys) != 0 {
		t.Errorf("Expected 0 keys, got %d", len(keys))
	}
}

// TestKeysSingleEntry tests the Keys method on a cache with a single entry.
func TestKeysSingleEntry(t *testing.T) {
	cache := New[string, int](10, 10, 5*time.Minute, simpleHash)
	cache.Add("key1", 123)
	keys := cache.Keys()
	if len(keys) != 1 || keys[0] != "key1" {
		t.Errorf("Expected [key1], got %v", keys)
	}
}

// TestKeysMultipleEntries tests the Keys method on a cache with multiple entries.
func TestKeysMultipleEntries(t *testing.T) {
	cache := New[string, int](10, 10, 5*time.Minute, simpleHash)
	cache.Add("key1", 123)
	cache.Add("key2", 456)
	cache.Add("key3", 789)
	keys := cache.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
	expectedKeys := map[string]bool{"key1": true, "key2": true, "key3": true}
	for _, key := range keys {
		if !expectedKeys[key] {
			t.Errorf("Unexpected key: %s", key)
		}
	}
}

// TestKeysWithEviction tests the Keys method on a cache where some entries have been evicted.
func TestKeysWithEviction(t *testing.T) {
	cache := New[string, int](10, 10, 1*time.Second, simpleHash)
	cache.Add("key1", 123)
	time.Sleep(2 * time.Second) // Wait for eviction time to pass
	cache.Add("key2", 456)
	keys := cache.Keys()
	if len(keys) != 1 || keys[0] != "key2" {
		t.Errorf("Expected [key2] after eviction, got %v", keys)
	}
}

func BenchmarkCache_SetGet(b *testing.B) {
	cache := New[[32]byte, bool](bucketCount, maxEntries, evictionDuration, HashKey[[32]byte])

	// Fill
	for i := 0; i < bucketCount; i++ {
		key := generateKey(i)
		value := i%2 == 0 // Alternate bool values
		cache.Add(key, value)
	}

	// Benchmark setGet operation
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		key := generateKey(i)
		value := i%2 != 0 // Alternate bool values (different keys)
		b.StartTimer()
		cache.Add(key, value)
		_, ok := cache.Get(key)
		require.True(b, ok)
	}
}

// Keys are not present worse case
func BenchmarkConcurrentGet(b *testing.B) {
	c := New[[32]byte, bool](bucketCount, maxEntries, evictionDuration, HashKey[[32]byte])
	setupFullCache(c)

	//b.Log("full cache size", c.Size(), "expected", bucketCount*maxEntries)
	i := atomic.Int64{}
	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(32)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := generateKey(int(i.Add(1)))
			_ = c.Contains(key)
		}
	})
}

func BenchmarkGet_LRU(b *testing.B) {
	c := lru.NewLRU[[32]byte, bool](bucketCount, nil, time.Second*10)
	setupFullLRU(c)

	b.ResetTimer()
	i := atomic.Int64{}
	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(32)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := generateKey(int(i.Add(1)))
			_, _ = c.Get(key)
		}
	})
}

func BenchmarkSet(b *testing.B) {
	c := New[[32]byte, bool](bucketCount, maxEntries, evictionDuration, HashKey[[32]byte])
	setupFullCache(c)
	//b.Log("full cache size", c.Size(), "expected", bucketCount*maxEntries)

	i := atomic.Int64{}
	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(32)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := generateKey(int(i.Add(1)))
			c.Add(key, true)
		}
	})
}

func BenchmarkSet_LRU(b *testing.B) {
	c := lru.NewLRU[[32]byte, bool](bucketCount, nil, time.Second*2)
	setupFullLRU(c)
	//cleanupFullLRU(c)

	i := atomic.Int64{}
	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(32)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := generateKey(int(i.Add(1)))
			c.Add(key, true)
		}
	})
}

func BenchmarkConcurrentGetSet(b *testing.B) {
	c := New[[32]byte, bool](bucketCount, maxEntries, evictionDuration, HashKey[[32]byte])
	setupFullCache(c)
	//b.Log("full cache size", c.Size(), "expected", bucketCount*maxEntries)
	i := atomic.Int64{}
	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(32)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := generateKey(int(i.Add(1)))
			c.Add(key, true)
			_ = c.Contains(key)
		}
	})
}

func BenchmarkConcurrentGetSet_LRU(b *testing.B) {
	c := lru.NewLRU[[32]byte, bool](bucketCount*maxEntries, nil, time.Second*2)
	setupFullLRU(c)

	i := atomic.Int64{}
	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(32)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := generateKey(int(i.Add(1)))
			c.Add(key, true)
			_ = c.Contains(key)
		}
	})
}
