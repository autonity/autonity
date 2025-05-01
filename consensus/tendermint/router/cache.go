package router

import (
	"fmt"
	"sync"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
)

const (
	cacheEntryTTL = 30 * time.Minute
)

type CacheEntry struct {
	Recipients []common.Address
	Version    int64
	LastUsed   time.Time
}

type PeerSelectionCache struct {
	cache        map[string]CacheEntry
	cacheMu      sync.RWMutex
	cacheVersion int64
}

func NewPeerSelectionCache() *PeerSelectionCache {
	return &PeerSelectionCache{
		cache:        make(map[string]CacheEntry),
		cacheVersion: 1,
	}
}

func (c *PeerSelectionCache) Get(key string) (CacheEntry, bool) {
	c.cacheMu.RLock()
	defer c.cacheMu.RUnlock()
	entry, exists := c.cache[key]
	if entry.Version == c.cacheVersion {
		return entry, exists
	} else {
		return CacheEntry{}, false
	}
}

func (c *PeerSelectionCache) Set(key string, recipients []common.Address) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.cache[key] = CacheEntry{
		Recipients: recipients,
		Version:    c.cacheVersion,
		LastUsed:   time.Now(),
	}
}

func (c *PeerSelectionCache) UpdateLastUsed(key string) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	if entry, exists := c.cache[key]; exists {
		entry.LastUsed = time.Now()
		c.cache[key] = entry
	}
}

func (c *PeerSelectionCache) Invalidate() {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.cacheVersion++
}

func (c *PeerSelectionCache) Cleanup() {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	now := time.Now()
	removed := 0
	for key, entry := range c.cache {
		if now.Sub(entry.LastUsed) > cacheEntryTTL {
			delete(c.cache, key)
			removed++
		}
	}
	if removed > 0 {
		log.Debug("PeerSelectionCache: cleaned up cache entries", "count", removed, "remaining", len(c.cache))
	}
}

func GenerateCacheKey(from common.Address, senderType SenderType, msgCode uint8) string {
	return fmt.Sprintf("%s-%d-%d", from.Hex(), senderType, msgCode)
}
