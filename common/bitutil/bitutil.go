// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Adapted from: https://golang.org/src/crypto/cipher/xor.go

// Package bitutil implements fast bitwise operations.
package bitutil

import (
	"math/bits"
	"runtime"
	"unsafe"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

// XORBytes xors the bytes in a and b. The destination is assumed to have enough
// space. Returns the number of bytes xor'd.
func XORBytes(dst, a, b []byte) int {
	if supportsUnaligned {
		return fastXORBytes(dst, a, b)
	}
	return safeXORBytes(dst, a, b)
}

// fastXORBytes xors in bulk. It only works on architectures that support
// unaligned read/writes.
func fastXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] ^ bw[i]
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

// safeXORBytes xors one by one. It works on all architectures, independent if
// it supports unaligned read/writes or not.
func safeXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

// ANDBytes ands the bytes in a and b. The destination is assumed to have enough
// space. Returns the number of bytes and'd.
func ANDBytes(dst, a, b []byte) int {
	if supportsUnaligned {
		return fastANDBytes(dst, a, b)
	}
	return safeANDBytes(dst, a, b)
}

// fastANDBytes ands in bulk. It only works on architectures that support
// unaligned read/writes.
func fastANDBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] & bw[i]
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] & b[i]
	}
	return n
}

// safeANDBytes ands one by one. It works on all architectures, independent if
// it supports unaligned read/writes or not.
func safeANDBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] & b[i]
	}
	return n
}

// ORBytes ors the bytes in a and b. The destination is assumed to have enough
// space. Returns the number of bytes or'd.
func ORBytes(dst, a, b []byte) int {
	if supportsUnaligned {
		return fastORBytes(dst, a, b)
	}
	return safeORBytes(dst, a, b)
}

// fastORBytes ors in bulk. It only works on architectures that support
// unaligned read/writes.
func fastORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] | bw[i]
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] | b[i]
	}
	return n
}

// safeORBytes ors one by one. It works on all architectures, independent if
// it supports unaligned read/writes or not.
func safeORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] | b[i]
	}
	return n
}

// TestBytes tests whether any bit is set in the input byte slice.
func TestBytes(p []byte) bool {
	if supportsUnaligned {
		return fastTestBytes(p)
	}
	return safeTestBytes(p)
}

// fastTestBytes tests for set bits in bulk. It only works on architectures that
// support unaligned read/writes.
func fastTestBytes(p []byte) bool {
	n := len(p)
	w := n / wordSize
	if w > 0 {
		pw := *(*[]uintptr)(unsafe.Pointer(&p))
		for i := 0; i < w; i++ {
			if pw[i] != 0 {
				return true
			}
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		if p[i] != 0 {
			return true
		}
	}
	return false
}

// safeTestBytes tests for set bits one byte at a time. It works on all
// architectures, independent if it supports unaligned read/writes or not.
func safeTestBytes(p []byte) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return true
		}
	}
	return false
}

// returns the number of 1 bits in `n` at lower positions than the LSB of `c`
// where the set of 1 bits in `c` excluding the LSB must be a subset of set of 1 bits in `n`.
// the caller is responsible to make sure the mentioned condition is true.
func OnesCountLowerLSB8(n, c uint8) int {
	if c == 0 {
		return 0
	}
	// remove all 1 bits from `n` which are common with `c`
	commonBitsRemoved := n ^ c
	// `commonBitsRemoved` has all the 1 bits of `n` which are not present in `c`

	// all the positions in `commonBitsRemoved` which are lower than the LSB of `c`
	// have 1 only if `n` has 1 in those positions

	// so we need to know the count of 1 bits in `commonBitsRemoved` at lower positions
	// than the LSB of `c`

	// `commonBitsRemoved` and `c` have no 1 bits in common or have exaclty one in common
	// at the position of the LSB of `c`

	// for example, if `n = 10101100` and `c = 10001000` (`c` is a subset of `n`),
	// then `commonBitsRemoved = 00100100` (no 1 bit common between `c` and `commonBitsRemoved`)
	// again, if `n = 10101100` and `c = 00110000` (`c`, excluding LSB, is a subset of `n`),
	// then `commonBitsRemoved = 10011100` (only the LSB of `c` is common between `c` and `commonBitsRemoved`)

	onlyLowerBits := commonBitsRemoved & (c - 1)
	// all the positions in `c-1` which are lower than the LSB of `c` have 1
	// and all the positions in `c-1` which greater than the LSB of `c`
	// have the same bit as `c` in those positions. `c-1` has 0 in the position
	// of the LSB of `c`.

	// for example if `c = 10100110000`, then `c-1 = 10100101111`

	// As `commonBitsRemoved` and `c` have no common 1 bits except LSB,
	// `commonBitsRemoved` and `c-1` have common bits only in the
	// lower positions than the LSB of `c`, because all the positions in `c-1`
	// which are lower than the LSB of `c` has 1.

	// Again `onlyLowerBits` is a subset of both `commonBitsRemoved` and `c-1`
	// i.e. `onlyLowerBits` only has the common bits of `commonBitsRemoved` and `c-1`

	// Now, because all the positions in `c-1` which are greater than or equal to the LSB of `c`
	// have no common 1 bits with `commonBitsRemoved`, `onlyLowerBits` has 0 in all those
	// positions. And because all the positions in `c-1` which are lower than the LSB of `c` have 1,
	// `onlyLowerBits` has 1 in those positions only if `commonBitsRemoved` has 1 in those positions

	// So we just need to count the 1 bits in the whole number `onlyLowerBits`
	return bits.OnesCount8(onlyLowerBits)
}

// returns the number of 1 bits in `n` at lower positions than the `index`
func OnesCountBeforeIndex8(n, index uint8) int {
	if index == 0 {
		return 0
	}
	if index > 7 {
		return bits.OnesCount8(n)
	}
	// take only the bits at lower positions than the `index`
	return bits.OnesCount8(n & (1 << index))
}
