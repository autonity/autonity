package blst

import (
	"encoding/binary"

	blstbind "github.com/supranational/blst/bindings/go"
)

func ToBlstScalars[T uint16 | uint32](coefficients []T) []*blstbind.Scalar {
	scalars := make([]*blstbind.Scalar, 0, len(coefficients))
	// always use 32 bytes, otherwise it breaks
	bytes := make([]byte, 32)

	for _, c := range coefficients {
		// use little endian, as these coefficients need to multiplied as they are
		binary.LittleEndian.PutUint32(bytes, uint32(c))
		scalar := new(blstbind.Scalar)
		scalar.FromLEndian(bytes)
		scalars = append(scalars, scalar)
	}
	return scalars
}
