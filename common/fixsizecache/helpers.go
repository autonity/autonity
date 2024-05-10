package fixsizecache

import (
	"math"

	"github.com/autonity/autonity/common"
)

type keyConstraint interface {
	common.Hash | [32]byte
}

// HashKey is an helper which calculates the bucket index for a given key
func HashKey[T keyConstraint](key T) uint {
	var result int
	for i, k := range key {
		result |= int(k) << (i * 8) // Shift each byte by its position and OR with the result
	}
	return uint(result)
}

// isPrime checks if a number is prime
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := int(5); float64(i) <= math.Sqrt(float64(n)); i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// NextPrime returns the smallest prime number greater than n
func NextPrime(n int) uint {
	if n <= 0 {
		return 0
	}
	for {
		n++
		if isPrime(n) {
			return uint(n)
		}
	}
}

// PreviousPrime returns the largest prime number less than n
func PreviousPrime(n int) uint {
	for {
		if n <= 0 {
			return 0
		}
		n--
		if isPrime(n) {
			return uint(n)
		}
	}
}
