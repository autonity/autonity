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
	mutex      sync.Mutex
	limits     map[common.Address]*rateRecord
	timeWindow time.Duration
	maxBurst   uint64
}

type rateRecord struct {
	count      uint64
	expiration time.Time
}

func NewTimeWindowLimiter(window time.Duration, maxBurst uint64) *TimeWindowLimiter {
	return &TimeWindowLimiter{
		limits:     make(map[common.Address]*rateRecord),
		timeWindow: window,
		maxBurst:   maxBurst,
	}
}

func (l *TimeWindowLimiter) Allow(sender common.Address) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	record, exists := l.limits[sender]

	if !exists || now.After(record.expiration) {
		l.limits[sender] = &rateRecord{
			count:      1,
			expiration: now.Add(l.timeWindow),
		}
		return nil
	}

	if record.count >= l.maxBurst {
		return ErrRateLimitExceeded
	}

	record.count++
	return nil
}

func (l *TimeWindowLimiter) Cleanup() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	for addr, record := range l.limits {
		if now.After(record.expiration) {
			delete(l.limits, addr)
		}
	}
}

type HeightBasedLimiter struct {
	mutex        sync.Mutex
	records      map[common.Address]map[uint64]uint64
	maxPerHeight uint64
	btl          uint64
}

func NewHeightBasedLimiter(maxPerHeight uint64, btl uint64) *HeightBasedLimiter {
	return &HeightBasedLimiter{
		records:      make(map[common.Address]map[uint64]uint64),
		maxPerHeight: maxPerHeight,
		btl:          btl,
	}
}

func (l *HeightBasedLimiter) Allow(sender common.Address, height uint64) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, exists := l.records[sender]; !exists {
		l.records[sender] = make(map[uint64]uint64)
	}

	if l.records[sender][height] >= l.maxPerHeight {
		return ErrHeightQuotaExhausted
	}

	l.records[sender][height]++
	return nil
}

func (l *HeightBasedLimiter) Cleanup(head uint64) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	staled := head - l.btl
	for addr, heights := range l.records {
		for h := range heights {
			if h <= staled {
				delete(heights, h)
			}
		}
		if len(heights) == 0 {
			delete(l.records, addr)
		}
	}
}

type DuplicateLimiter struct {
	mutex   sync.Mutex
	records map[common.Address]map[common.Hash]time.Time
	ttl     time.Duration
}

func NewDuplicateTracker(ttl time.Duration) *DuplicateLimiter {
	return &DuplicateLimiter{
		records: make(map[common.Address]map[common.Hash]time.Time),
		ttl:     ttl,
	}
}

func (t *DuplicateLimiter) Allow(sender common.Address, hash common.Hash) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if hashes, exists := t.records[sender]; exists {
		if _, duplicate := hashes[hash]; duplicate {
			return ErrDuplicateMessage
		}
	}

	if _, exists := t.records[sender]; !exists {
		t.records[sender] = make(map[common.Hash]time.Time)
	}

	t.records[sender][hash] = time.Now()
	return nil
}

func (t *DuplicateLimiter) Cleanup() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	cutoff := time.Now().Add(-t.ttl)
	for addr, hashes := range t.records {
		for h, timestamp := range hashes {
			if timestamp.Before(cutoff) {
				delete(hashes, h)
			}
		}
		if len(hashes) == 0 {
			delete(t.records, addr)
		}
	}
}

type AFDRateLimiter struct {
	timeLimiter      *TimeWindowLimiter
	heightLimiter    *HeightBasedLimiter
	duplicateLimiter *DuplicateLimiter
}

func NewAFDRateLimiter() *AFDRateLimiter {
	limiter := &AFDRateLimiter{
		// since communication channel is asynchronous, those pending write of off chain accusation msgs from a sender
		// could potentially be received once the peer connection get established from a disaster recovery, thus it
		// could exceed the number of accusation that could be produced by rule engine over a height, so we set higher
		// rate limit during 1 second to be tolerant for such case.
		// 8 accusations per 1s window for per client, rate limit reset per 1s.
		timeLimiter: NewTimeWindowLimiter(time.Second, maxAccusationPerHeight*2),
		// 4 accusations per height for per client.
		heightLimiter: NewHeightBasedLimiter(maxAccusationPerHeight, HeightRange),
		// duplicated accusation checker, reset per 5 minutes.
		duplicateLimiter: NewDuplicateTracker(time.Minute * 5),
	}

	return limiter
}

func (l *AFDRateLimiter) Cleanup(height uint64) {
	l.timeLimiter.Cleanup()
	l.heightLimiter.Cleanup(height)
	l.duplicateLimiter.Cleanup()
}
