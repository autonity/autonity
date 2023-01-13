package bft

import (
	"math"
	"math/big"
	"testing"
)

func TestQuorum(t *testing.T) {

	tests := []uint64{1, 2, 3, 4, 5, 6, 7, 11, 15, 20, 100, 150, 2000, 2509, 3045, 7689, 10032, 12932, 15982, 301234, 301235, 301236}

	verify := func(vp uint64) uint64 {
		return uint64(math.Ceil((2 * float64(vp)) / 3.))
	}
	for _, tt := range tests {
		if got, want := Quorum(new(big.Int).SetUint64(tt)), verify(tt); got.Uint64() != want {
			t.Errorf("Quorum() = %v, want %v", got, want)
		}
	}
}

func TestF(t *testing.T) {

	tests := []uint64{1, 2, 3, 4, 5, 6, 7, 11, 15, 20, 100, 150, 2000, 2509, 3045, 7689, 10032, 12932, 15982, 301234, 301235, 301236}

	verify := func(vp uint64) uint64 {
		return uint64(math.Ceil(float64(vp)/3)) - 1
	}
	for _, tt := range tests {
		if got, want := F(new(big.Int).SetUint64(tt)), verify(tt); got.Uint64() != want {
			t.Errorf("Quorum() = %v, want %v", got, want)
		}
	}
}
