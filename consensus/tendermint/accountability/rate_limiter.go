package accountability

import (
	"errors"
	"github.com/autonity/autonity/common"
	"sync"
	"time"
)

var (
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrHeightQuotaExhausted = errors.New("block height quota exhausted")
	ErrDuplicateMessage     = errors.New("duplicate message detected")
)

type TimeWindowLimiter struct {
	sync.RWMutex
	limits     map[common.Address]*rateRecord
	timeWindow time.Duration
	burst      int
	ttl        time.Duration
}

type rateRecord struct {
	count      int
	last       time.Time
	expiration time.Time
}

func NewTimeWindowLimiter(window time.Duration, burst int, ttl time.Duration) *TimeWindowLimiter {
	return &TimeWindowLimiter{
		limits:     make(map[common.Address]*rateRecord),
		timeWindow: window,
		burst:      burst,
		ttl:        ttl,
	}
}

func (l *TimeWindowLimiter) Allow(sender common.Address) error {
	l.Lock()
	defer l.Unlock()

	now := time.Now()
	record, exists := l.limits[sender]

	if !exists || now.After(record.expiration) {
		l.limits[sender] = &rateRecord{
			count:      1,
			last:       now,
			expiration: now.Add(l.ttl),
		}
		return nil
	}

	if record.count >= l.burst {
		return ErrRateLimitExceeded
	}

	record.count++
	record.last = now
	return nil
}

func (l *TimeWindowLimiter) Cleanup() {
	l.Lock()
	defer l.Unlock()

	now := time.Now()
	for addr, record := range l.limits {
		if now.After(record.expiration) {
			delete(l.limits, addr)
		}
	}
}

type HeightBasedLimiter struct {
	sync.RWMutex
	accusations  map[common.Address]map[uint64]int
	maxPerHeight int
	btl          uint64
}

func NewHeightBasedLimiter(maxPerHeight int, btl uint64) *HeightBasedLimiter {
	return &HeightBasedLimiter{
		accusations:  make(map[common.Address]map[uint64]int),
		maxPerHeight: maxPerHeight,
		btl:          btl,
	}
}

func (l *HeightBasedLimiter) Allow(sender common.Address, height uint64) error {
	l.Lock()
	defer l.Unlock()

	if _, exists := l.accusations[sender]; !exists {
		l.accusations[sender] = make(map[uint64]int)
	}

	if l.accusations[sender][height] >= l.maxPerHeight {
		return ErrHeightQuotaExhausted
	}

	l.accusations[sender][height]++
	return nil
}

func (l *HeightBasedLimiter) Cleanup(head uint64) {
	l.Lock()
	defer l.Unlock()

	staled := head - l.btl
	for addr, heights := range l.accusations {
		for h := range heights {
			if h <= staled {
				delete(heights, h)
			}
		}
		if len(heights) == 0 {
			delete(l.accusations, addr)
		}
	}
}

type DuplicateTracker struct {
	sync.RWMutex
	messages map[common.Address]map[common.Hash]time.Time
	ttl      time.Duration
}

func NewDuplicateTracker(ttl time.Duration) *DuplicateTracker {
	return &DuplicateTracker{
		messages: make(map[common.Address]map[common.Hash]time.Time),
		ttl:      ttl,
	}
}

func (t *DuplicateTracker) Allow(sender common.Address, hash common.Hash) error {
	t.RLock()
	if hashes, exists := t.messages[sender]; exists {
		if _, duplicate := hashes[hash]; duplicate {
			t.RUnlock()
			return ErrDuplicateMessage
		}
	}
	t.RUnlock()

	t.Lock()
	defer t.Unlock()

	if _, exists := t.messages[sender]; !exists {
		t.messages[sender] = make(map[common.Hash]time.Time)
	}

	t.messages[sender][hash] = time.Now()
	return nil
}

func (t *DuplicateTracker) Cleanup() {
	t.Lock()
	defer t.Unlock()

	cutoff := time.Now().Add(-t.ttl)
	for addr, hashes := range t.messages {
		for h, timestamp := range hashes {
			if timestamp.Before(cutoff) {
				delete(hashes, h)
			}
		}
		if len(hashes) == 0 {
			delete(t.messages, addr)
		}
	}
}

type AFDRateLimiter struct {
	timeLimiter    *TimeWindowLimiter
	heightLimiter  *HeightBasedLimiter
	duplicateCheck *DuplicateTracker
}

func NewAFDRateLimiter() *AFDRateLimiter {
	limiter := &AFDRateLimiter{
		// since communication channel is asynchronous, those pending write of off chain accusation msgs from a sender
		// could potentially be received once the peer connection get established from a disaster recovery, thus it
		// could exceed the number of accusation that could be produced by rule engine over a height, so we set higher
		// rate limit during 1 second to be tolerant for such case.
		// 8 accusations per 1s window for per client, rate limit reset per 1s.
		timeLimiter: NewTimeWindowLimiter(time.Second, maxAccusationPerHeight*2, time.Second),
		// 4 accusations per height for per client.
		heightLimiter: NewHeightBasedLimiter(maxAccusationPerHeight, msgGCInterval),
		// duplicated accusation checker, reset per 10 minutes.
		duplicateCheck: NewDuplicateTracker(time.Minute * 10),
	}

	return limiter
}

func (l *AFDRateLimiter) Cleanup(height uint64) {
	l.timeLimiter.Cleanup()
	l.heightLimiter.Cleanup(height)
	l.duplicateCheck.Cleanup()
}
