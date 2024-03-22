package types

import (
	"fmt"
	"testing"
)

// TODO(lorenzo) do proper testing
// test also the various flatten
func TestCoefficients(t *testing.T) {
	t.Run("test increment", func(t *testing.T) {
		c := NewCoefficients(100)
		c.Increment(0)
		c.Increment(50)
		c.Increment(99)
		c.Increment(99)
		c.Increment(99)
		c.Increment(99)
		fmt.Println(c)
	})
}
