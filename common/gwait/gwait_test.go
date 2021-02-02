package gwait_test

import (
	"sync/atomic"
	"testing"

	"github.com/clearmatics/autonity/common/gwait"
	"github.com/stretchr/testify/assert"
)

func TestWaiting(t *testing.T) {
	w := gwait.NewWaiter()

	var expected int32 = 1
	var x int32 = 0
	// Run function in goroutine
	w.Go(func() {
		atomic.StoreInt32(&x, expected)
	})
	// Ensure we wait for it to complete
	w.Wait()

	// Chekc that the value of x was updated
	assert.Equal(t, atomic.LoadInt32(&x), expected)
}
