package bitutil

import (
	"errors"
	"fmt"
	"math/bits"
)

const (
	BitsPerItem  = 1 //NOTE: if this gets changed, major refactoring will be needed for this file. Proceed with caution.
	BitsInByte   = 8
	ItemsPerByte = BitsInByte / BitsPerItem
)

var ErrSizeMismatch = errors.New("comparing Bitmap with different sized Bitmap")

type Bitmap []byte

// Creates a bitmap to support `totalItem` items where each item is represented by a single bit
func NewBitmap(totalItem int) Bitmap { //nolint
	byteLength := (totalItem*BitsPerItem + BitsInByte - 1) / BitsInByte
	return make(Bitmap, byteLength)
}

// ensures that the Bitmap has the correct length compared to the committee size
func (bm Bitmap) Valid(totalItem int) bool {
	expectedByteLength := (totalItem*BitsPerItem + BitsInByte - 1) / BitsInByte
	return len(bm) == expectedByteLength
}

// Returns `true` if the Bitmap contains `item`
//
// this function needs to be modified if the constant `ItemsPerByte` is changed
func (bm Bitmap) HasItem(item int) bool {
	byteIndex, bitIndex := itemToBitMapPosition(item)

	// because of the constant `BitsPerItem = 1`, each validator takes a single bit
	// we just need to check if the bit is 1 or 0
	// note that the bits are numbered from LSB to MSB (in both `HasItem` and `AddItem`)
	return (bm[byteIndex] & (1 << bitIndex)) > 0
}

// Adds the `item` into the Bitmap.
//
// this function needs to be modified if the constant `ItemsPerByte` is changed
func (bm Bitmap) AddItem(item int) {
	byteIndex, bitIndex := itemToBitMapPosition(item)

	// because of the constant `BitsPerItem = 1`, each validator takes a single bit
	// we just need to set 1 in `bitIndex`
	// note that the bits are numbered from LSB to MSB (in both `HasItem` and `AddItem`)
	bm[byteIndex] = bm[byteIndex] | (1 << bitIndex)
}

// returns the total number of item in bitmap
func (bm Bitmap) ItemCount() int {
	countNonZero := 0
	for _, b := range bm {
		if b == 0 {
			continue
		}
		countNonZero += bits.OnesCount8(b)
	}
	return countNonZero
}

// returns `true` if `bm` contains `other`, or `other` is a subset of `bm`
func (bm Bitmap) Contains(other Bitmap) bool {
	if len(bm) != len(other) {
		panic(ErrSizeMismatch.Error())
	}
	// can we use `bitutil.ANDBytes` instead?
	for i, b := range other {
		if (b & bm[i]) != b {
			return false
		}
	}
	return true
}

// Returns a new common set of `bm` and `other`.
// `bm` and `other` will not be modified
func (bm Bitmap) CommonSet(other Bitmap) Bitmap {
	if len(bm) != len(other) {
		panic(ErrSizeMismatch.Error())
	}
	newBitmap := make([]byte, len(bm))
	for i, b := range other {
		newBitmap[i] = bm[i] & b
	}
	return newBitmap
}

// Returns a new union set of `bm` and `other`.
// `bm` and `other` will not be modified
func (bm Bitmap) UnionSet(other Bitmap) Bitmap {
	if len(bm) != len(other) {
		panic(ErrSizeMismatch.Error())
	}

	newBitmap := make([]byte, len(bm))
	for i, b := range other {
		newBitmap[i] = bm[i] | b
	}

	return newBitmap
}

// Merges `other` Bitmap with `bm`, modifying `bm` in the process.
//
// returns true if the `other` object contributes to the `bm` object
func (bm Bitmap) Merge(other Bitmap) bool {
	if len(bm) != len(other) {
		panic(ErrSizeMismatch.Error())
	}
	contributed := false
	for i, b := range other {
		if (bm[i] & b) != b {
			contributed = true
		}
		bm[i] |= b
	}
	return contributed
}

// Returns a list of items in the bitmap
func (bm Bitmap) ItemList() []int {
	items := make([]int, 0)

	for i, b := range bm {

		for b > 0 {
			// take the LSB of `b` and convert it into `item`

			bitIndex := ByteLSBPosition(b)
			item := bitMapPositionToItem(i, bitIndex)
			items = append(items, item)

			// remove the LSB of `b`
			b = b & (b - 1)
			// the above idea is inspired from Brian Kernighan's Algorithm
			// see more here: https://how.dev/answers/what-is-kernighans-algorithm
		}
	}

	return items
}

func (bm Bitmap) ItemCountBefore(lastItem int) int {
	byteIndex, bitIndex := itemToBitMapPosition(lastItem)
	count := 0
	for i := range byteIndex {
		count += bits.OnesCount8(bm[i])
	}
	return count + ByteOnesCountBeforeIndex(bm[byteIndex], bitIndex)
}

func (bm Bitmap) ItemListBefore(lastItem int) []int {
	fmt.Printf("start\n")
	items := make([]int, 0)

	for i, b := range bm {
		fmt.Printf("i %v, b %v\n", i, b)

		for b > 0 {
			// take the LSB of `b` and convert it into `item`

			bitIndex := ByteLSBPosition(b)
			item := bitMapPositionToItem(i, bitIndex)
			if item >= lastItem {
				return items
			}
			items = append(items, item)

			// remove the LSB of `b`
			b = b & (b - 1)
			// the above idea is inspired from Brian Kernighan's Algorithm
			// see more here: https://how.dev/answers/what-is-kernighans-algorithm
		}
	}
	fmt.Printf("end\n")

	return items
}

func (bm Bitmap) ForEachItem(callback func(itemIndex, item int), lastItem int) {
	for index, item := range bm.ItemListBefore(lastItem) {
		callback(index, item)
	}
}

func itemToBitMapPosition(index int) (int, int) {
	// the following line is the same as `index / ItemsPerByte`
	// but the following works because `ItemsPerByte = 8 = 2^3`
	byteIndex := index >> 3

	// the following line is the same as `index % ItemsPerByte`
	// but the following works because `ItemsPerByte = 8`, which is a power of 2
	bitIndex := index & (ItemsPerByte - 1)
	return byteIndex, bitIndex
}

func bitMapPositionToItem(byteIndex, bitIndex int) int {
	// the following line is the same as `byteIndex * ItemsPerByte + bitIndex`
	// but the following works because `ItemsPerByte = 8 = 2^3` and `bitIndex < ItemsPerByte`
	// `bitIndex < ItemsPerByte` should always be true
	return (byteIndex << 3) | bitIndex
}
