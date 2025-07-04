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
func ByteOnesCountLowerLSB(n, c uint8) int {
	return ByteOnesCountLowerIndex(n, ByteLSBPosition(n))
}

// returns the number of 1 bits in `n` at lower positions than the `index`
func ByteOnesCountLowerIndex(n uint8, index int) int {
	if index <= 0 {
		return 0
	}
	if index > 7 {
		return bits.OnesCount8(n)
	}
	// take only the bits at lower positions than the `index`
	return bits.OnesCount8(n & (1 << index))
}

func ByteLSBPosition(n uint8) int {
	if n == 0 {
		return -1
	}
	// remove everything but LSB of the `n`

	// also add 1 bit to all the positions in `n` which are
	// lower than the LSB, i.e., calculate `pow(2,x+1) - 1`,
	// where `x` is the position of the LSB of `n`
	lsbBecomesMSB := n ^ (n - 1)

	// we want to know `x`

	// because `lsbBecomesMSB = pow(2,x+1) - 1`, it has `x+1` 1 bits
	return bits.OnesCount8(lsbBecomesMSB) - 1
}

func Uint16LSBPosition(n uint16) int {
	if n == 0 {
		return -1
	}
	// remove everything but LSB of the `n`

	// also add 1 bit to all the positions in `n` which are
	// lower than the LSB, i.e., calculate `pow(2,x+1) - 1`,
	// where `x` is the position of the LSB of `n`
	lsbBecomesMSB := n ^ (n - 1)

	// we want to know `x`

	// because `lsbBecomesMSB = pow(2,x+1) - 1`, it has `x+1` 1 bits
	return bits.OnesCount16(lsbBecomesMSB) - 1
}

func Uint32LSBPosition(n uint32) int {
	if n == 0 {
		return -1
	}
	// remove everything but LSB of the `n`

	// also add 1 bit to all the positions in `n` which are
	// lower than the LSB, i.e., calculate `pow(2,x+1) - 1`,
	// where `x` is the position of the LSB of `n`
	lsbBecomesMSB := n ^ (n - 1)

	// we want to know `x`

	// because `lsbBecomesMSB = pow(2,x+1) - 1`, it has `x+1` 1 bits
	return bits.OnesCount32(lsbBecomesMSB) - 1
}

func Uint16MSBPosition(n uint16) int {
	for n > 0 {
		// check if `n` is a power of 2, then we have the MSB
		if (n & (n - 1)) == 0 {
			return Uint16LSBPosition(n)
		}
		n = n & (n - 1)
	}
	// 0 input
	return -1
}

func Uint32MSBPosition(n uint32) int {
	for n > 0 {
		// check if `n` is a power of 2, then we have the MSB
		if (n & (n - 1)) == 0 {
			return Uint32LSBPosition(n)
		}
		n = n & (n - 1)
	}
	// 0 input
	return -1
}
