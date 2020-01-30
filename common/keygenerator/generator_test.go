package keygenerator

import "testing"

func TestGenerator_Next(t *testing.T) {
	for i := 0; i < 200; i++ {
		_, err := Next()
		if err != nil {
			t.Fatal(i, err)
		}
	}
}
