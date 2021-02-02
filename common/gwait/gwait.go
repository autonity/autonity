package gwait

import (
	"sync"
)

// Waiter provides a mechanism to wait for a group of goroutines to complete.
type Waiter struct {
	wg *sync.WaitGroup
}

// NewWaiter returns a new Waiter instance.
func NewWaiter() *Waiter {
	return &Waiter{
		wg: &sync.WaitGroup{},
	}
}

// Go runs the given func in a new goroutine.
func (w Waiter) Go(f func()) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		f()
	}()
}

// Wait waits for all goroutines started by Go to complete before returning.
// Waiter may be re-used after this call returns.
func (w *Waiter) Wait() {
	w.wg.Wait()
}
