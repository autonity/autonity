package gwait_test

import "github.com/autonity/autonity/common/gwait"

func ExampleWaiter() {
	w := gwait.NewWaiter()

	// Run function in a goroutine
	w.Go(func() {
		println("hello")
	})

	// Run another function in a goroutine
	w.Go(func() {
		println("byebye")
	})

	// Ensure we wait for the goroutines to complete
	w.Wait()
}
