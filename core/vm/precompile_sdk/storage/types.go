package storage

import (
	"math/big"

	"github.com/holiman/uint256"
)

type Uint256 struct {
	uint256.Int
}

func NewUint256FromInt(val uint64) Uint256 {
	return Uint256{*uint256.NewInt(val)}
}

func NewUint256FromBig(val *big.Int) Uint256 {
	u := Uint256{*uint256.NewInt(0)}
	u.SetFromBig(val)
	return u
}

func (u *Uint256) Get() uint64 {
	return u.Uint64()
}
