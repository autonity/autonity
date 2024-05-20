package fixsizecache

import (
	"math"

	farmhash "github.com/leemcloughlin/gofarmhash"

	"github.com/autonity/autonity/common"
)

type keyConstraint interface {
	common.Hash | [32]byte
}

// HashKey is an helper which calculates the bucket index for a given key
func HashKey[T keyConstraint](key T) uint {
	return uint(farmhash.Hash64(key[:]))
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
	for i := 5; float64(i) <= math.Sqrt(float64(n)); i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// NextPrime returns the smallest prime number greater than or equal to n
func NextPrime(n int) uint {
	if n <= 2 {
		return 2
	}
	for {
		if isPrime(n) {
			return uint(n)
		}
		n++
	}
}

// PreviousPrime returns the largest prime number less than or equal to n
func PreviousPrime(n int) uint {
	for {
		if n <= 2 {
			return 2
		}
		if isPrime(n) {
			return uint(n)
		}
		n--
	}
}
