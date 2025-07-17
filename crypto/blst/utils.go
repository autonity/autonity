package blst

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/autonity/autonity/common"
)

func ToScalars(coefficients []*big.Int) []*blstScalar {
	scalars := make([]*blstScalar, 0, len(coefficients))
	bytes := make([]byte, BlstScalarBytes)

	for _, c := range coefficients {
		// coefficients that exceed 4 bytes (uint32) should get filtered before arriving here
		if c.BitLen() > common.QuorumCap {
			panic(fmt.Sprintf("coefficient too big: trying to fit %d bits in %d", c.BitLen(), common.QuorumCap))
		}
		// use little endian, as these coefficients need to multiplied as they are
		binary.LittleEndian.PutUint32(bytes, uint32(c.Uint64())) //nolint:gosec
		scalar := new(blstScalar)
		scalar.FromLEndian(bytes)
		scalars = append(scalars, scalar)
	}
	return scalars
}
