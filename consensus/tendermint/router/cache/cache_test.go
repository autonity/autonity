package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/autonity/autonity/common"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := New()
	addr := common.HexToAddress("0x123")
	key := GenerateKey(addr, 1, 2)
	recipients := []common.Address{common.HexToAddress("0x456"), common.HexToAddress("0x789")}

	cache.Set(key, recipients)
	entry, exists := cache.Get(key)

	assert.True(t, exists, "Expected entry to exist")
	assert.Equal(t, recipients, entry.Recipients, "Recipients should match")
	assert.Equal(t, int64(1), entry.Version, "Version should be 1")
	assert.False(t, entry.LastUsed.IsZero(), "LastUsed should be set")
}

func TestCache_UpdateLastUsed(t *testing.T) {
	cache := New()
	addr := common.HexToAddress("0x123")
	key := GenerateKey(addr, 1, 42)
	recipients := []common.Address{common.HexToAddress("0x456")}

	// Set initial entry
	cache.Set(key, recipients)
	entry, _ := cache.Get(key)
	originalTime := entry.LastUsed

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	cache.UpdateLastUsed(key)
	updatedEntry, exists := cache.Get(key)

	assert.True(t, exists, "Expected entry to exist after update")
	assert.True(t, updatedEntry.LastUsed.After(originalTime), "LastUsed should be updated")
	assert.Equal(t, recipients, updatedEntry.Recipients, "Recipients should remain unchanged")
}

func TestCache_Invalidate(t *testing.T) {
	cache := New()
	addr := common.HexToAddress("0x123")
	key := GenerateKey(addr, 1, 42)
	recipients := []common.Address{common.HexToAddress("0x456")}

	// Set and verify entry
	cache.Set(key, recipients)
	_, exists := cache.Get(key)
	assert.True(t, exists, "Expected entry to exist before invalidation")

	// Invalidate and check
	cache.Invalidate()
	_, exists = cache.Get(key)
	assert.False(t, exists, "Expected entry to be invalidated")
}

func TestCache_Cleanup(t *testing.T) {
	cache := New()
	addr1 := common.HexToAddress("0x123")
	addr2 := common.HexToAddress("0x456")
	key1 := GenerateKey(addr1, 1, 42)
	key2 := GenerateKey(addr2, 1, 42)
	recipients := []common.Address{common.HexToAddress("0x789")}

	// Set entries with different LastUsed times
	cache.Set(key1, recipients)
	time.Sleep(10 * time.Millisecond)
	cache.Set(key2, recipients)

	// Manipulate LastUsed for key1 to simulate expiration
	c := cache.(*peerCache)
	c.Lock()
	entry1 := c.recipients[key1]
	entry1.LastUsed = time.Now().Add(-31 * time.Minute) // Beyond TTL
	c.recipients[key1] = entry1
	c.Unlock()

	// Run cleanup
	cache.Cleanup()

	// Check results
	_, exists1 := cache.Get(key1)
	_, exists2 := cache.Get(key2)
	assert.False(t, exists1, "Expected expired entry to be removed")
	assert.True(t, exists2, "Expected non-expired entry to remain")
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := New()
	addr := common.HexToAddress("0x123")
	key := GenerateKey(addr, 1, 42)
	recipients := []common.Address{common.HexToAddress("0x456")}

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent Set and Get
	for i := 0; i < numGoroutines; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			cache.Set(key, recipients)
		}()
		go func() {
			defer wg.Done()
			cache.Get(key)
		}()
	}

	// Concurrent UpdateLastUsed
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.UpdateLastUsed(key)
		}()
	}

	// Concurrent Invalidate and Cleanup
	wg.Add(2)
	go func() {
		defer wg.Done()
		cache.Invalidate()
	}()
	go func() {
		defer wg.Done()
		cache.Cleanup()
	}()

	wg.Wait()

	// Verify cache integrity
	_, exists := cache.Get(key)
	if exists {
		// If the entry exists, it should be valid
		entry, _ := cache.Get(key)
		assert.Equal(t, recipients, entry.Recipients, "Recipients should match after concurrent operations")
	}
}

func TestGenerateKey(t *testing.T) {
	addr := common.HexToAddress("0x123")
	senderType := 1
	msgCode := uint8(42)

	key := GenerateKey(addr, senderType, msgCode)
	expected := addr.Hex() + "-1-42"
	assert.Equal(t, expected, key, "Generated key should match expected format")
}
