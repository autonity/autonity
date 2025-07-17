package types

import (
	"fmt"
	"io"
	"math/big"
	"math/bits"

	"github.com/autonity/autonity/rlp"
)

type Bitmap big.Int

func NewBitmap() *Bitmap {
	return (*Bitmap)(new(big.Int))
}

// check if b contains c
func (b *Bitmap) Contains(c *Bitmap) bool {
	bBig := (*big.Int)(b)
	cBig := (*big.Int)(c)
	return new(big.Int).And(bBig, cBig).Cmp(cBig) == 0
}

// Diff does the set minus operation
// B \ C = B & ~C gives the elements that are in B but not in C
func (b *Bitmap) Diff(c *Bitmap) *Bitmap {
	cBig := (*big.Int)(c)
	bBig := (*big.Int)(b)
	diff := new(big.Int).And(bBig, new(big.Int).Not(cBig))
	return (*Bitmap)(diff)
}

func (b *Bitmap) Set(index int) {
	bBig := (*big.Int)(b)
	bBig.SetBit(bBig, index, 1)
}

func (b *Bitmap) Unset(index int) {
	bBig := (*big.Int)(b)
	bBig.SetBit(bBig, index, 0)
}

func (b *Bitmap) IsSet(index int) bool {
	return (*big.Int)(b).Bit(index) == 1
}

func (b *Bitmap) Copy() *Bitmap {
	return (*Bitmap)(new(big.Int).Set((*big.Int)(b)))
}

func (b *Bitmap) Count() int {
	bBig := (*big.Int)(b)
	countNonZero := 0
	for _, b := range bBig.Bits() {
		if b == 0 {
			continue
		}
		// number of signers in `b`
		countNonZero += bits.OnesCount(uint(b))
	}
	return countNonZero
}

func (b *Bitmap) RightmostIndex() int {
	return int((*big.Int)(b).TrailingZeroBits()) //nolint:gosec
}

// Merge is a set union operation
func (b *Bitmap) Merge(c *Bitmap) *Bitmap {
	bBig := (*big.Int)(b)
	cBig := (*big.Int)(c)
	merged := new(big.Int).Or(bBig, cBig)
	return (*Bitmap)(merged)
}

func (b *Bitmap) Len() int {
	return (*big.Int)(b).BitLen()
}

func (b *Bitmap) Sign() int {
	return (*big.Int)(b).Sign()
}

func (b *Bitmap) String() string {
	return fmt.Sprintf("%08b", (*big.Int)(b).Bytes())
}

func (b *Bitmap) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, (*big.Int)(b))
}

func (b *Bitmap) DecodeRLP(stream *rlp.Stream) error {
	return stream.Decode((*big.Int)(b))
}

func (b *Bitmap) MarshalJSON() ([]byte, error) {
	return (*big.Int)(b).MarshalJSON()
}

func (b *Bitmap) UnmarshalJSON(data []byte) error {
	return (*big.Int)(b).UnmarshalJSON(data)
}
