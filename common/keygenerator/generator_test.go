package keygenerator

import (
	"testing"

	"github.com/clearmatics/autonity/crypto"
)

func TestGenerator_Next(t *testing.T) {
	for i := 0; i < 200; i++ {
		_, err := Next()
		if err != nil {
			t.Fatal(i, err)
		}
	}
}

func BenchmarkKeygenerator(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := Next()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEcdsaGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := crypto.GenerateKey()
		if err != nil {
			b.Fatal(err)
		}
	}
}
