package cache

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

type Recipients interface {
	Get(key string) (Entry, bool)
	Set(key string, recipients []common.Address)
	UpdateLastUsed(key string)
	Invalidate()
	Cleanup()
}

type Entry struct {
	Recipients []common.Address
	Version    int64
	LastUsed   time.Time
}

type peerCache struct {
	recipients map[string]Entry
	sync.RWMutex
	cacheVersion int64
}

func New() Recipients {
	return &peerCache{
		recipients:   make(map[string]Entry),
		cacheVersion: 1,
	}
}

func (c *peerCache) Get(key string) (Entry, bool) {
	c.RLock()
	defer c.RUnlock()
	entry, exists := c.recipients[key]
	if entry.Version == c.cacheVersion {
		return entry, exists
	}
	return Entry{}, false
}

func (c *peerCache) Set(key string, recipients []common.Address) {
	c.Lock()
	defer c.Unlock()
	c.recipients[key] = Entry{
		Recipients: recipients,
		Version:    c.cacheVersion,
		LastUsed:   time.Now(),
	}
}

func (c *peerCache) UpdateLastUsed(key string) {
	c.Lock()
	defer c.Unlock()
	if entry, exists := c.recipients[key]; exists {
		entry.LastUsed = time.Now()
		c.recipients[key] = entry
	}
}

func (c *peerCache) Invalidate() {
	c.Lock()
	defer c.Unlock()
	c.cacheVersion++
}

func (c *peerCache) Cleanup() {
	c.Lock()
	defer c.Unlock()
	now := time.Now()
	removed := 0
	for key, entry := range c.recipients {
		if now.Sub(entry.LastUsed) > cacheEntryTTL || entry.Version < c.cacheVersion {
			delete(c.recipients, key)
			removed++
		}
	}
	if removed > 0 {
		log.Debug("peerCache: cleaned up cache entries", "count", removed, "remaining", len(c.recipients))
	}
}

func GenerateKey(senderType int, isProposal bool) string {
	return fmt.Sprintf("%d-%v", senderType, isProposal)
}
