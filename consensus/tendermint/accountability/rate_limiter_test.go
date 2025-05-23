package accountability

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/autonity/autonity/common"
)

// TimeWindowLimiter Tests
func TestTimeWindowLimiter_BasicAllowance(t *testing.T) {
	limiter := NewTimeWindowLimiter(time.Second, 3, time.Minute)
	addr := common.HexToAddress("0x1")

	// under rate.
	for i := 0; i < 3; i++ {
		require.NoError(t, limiter.Allow(addr))
	}

	// over rate
	require.Equal(t, ErrRateLimitExceeded, limiter.Allow(addr))
}

func TestTimeWindowLimiter_TTLExpiration(t *testing.T) {
	limiter := NewTimeWindowLimiter(time.Second, 1, 500*time.Millisecond)
	addr := common.HexToAddress("0x2")

	// allow to
	require.NoError(t, limiter.Allow(addr))

	// rejected
	require.Equal(t, ErrRateLimitExceeded, limiter.Allow(addr))

	// allowed again as the last access expired.
	time.Sleep(600 * time.Millisecond)
	require.NoError(t, limiter.Allow(addr))
}

// HeightBasedLimiter Tests
func TestHeightBasedLimiter_HeightQuota(t *testing.T) {
	limiter := NewHeightBasedLimiter(2, 10)
	addr := common.HexToAddress("0x3")
	height := uint64(100)

	// allow for the 1st two tries.
	require.NoError(t, limiter.Allow(addr, height))
	require.NoError(t, limiter.Allow(addr, height))

	// over rated.
	require.Equal(t, ErrHeightQuotaExhausted, limiter.Allow(addr, height))

	// allowed for new height.
	require.NoError(t, limiter.Allow(addr, height+1))
}

func TestHeightBasedLimiter_CleanupLogic(t *testing.T) {
	limiter := NewHeightBasedLimiter(1, 5)
	addr := common.HexToAddress("0x4")

	require.NoError(t, limiter.Allow(addr, 99))
	require.NoError(t, limiter.Allow(addr, 100))

	// clean up
	limiter.Cleanup(105)

	limiter.rwMutex.Lock()
	defer limiter.rwMutex.Unlock()
	require.Empty(t, limiter.accusations[addr])
}

// DuplicateTracker Tests
func TestDuplicateTracker_DuplicateDetection(t *testing.T) {
	tracker := NewDuplicateTracker(time.Minute)
	addr := common.HexToAddress("0x5")
	hash := common.HexToHash("0xabc")

	// allow to
	require.NoError(t, tracker.Allow(addr, hash))

	// rejected as msg is duplicated.
	require.Equal(t, ErrDuplicateMessage, tracker.Allow(addr, hash))
}

func TestDuplicateTracker_TTLCleanup(t *testing.T) {
	tracker := NewDuplicateTracker(500 * time.Millisecond)
	addr := common.HexToAddress("0x6")
	hash := common.HexToHash("0xdef")

	require.NoError(t, tracker.Allow(addr, hash))

	tracker.Cleanup()
	require.Equal(t, ErrDuplicateMessage, tracker.Allow(addr, hash))

	time.Sleep(600 * time.Millisecond)
	tracker.Cleanup()
	require.NoError(t, tracker.Allow(addr, hash))
}
