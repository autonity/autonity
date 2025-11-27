package storage

import (
	"math/big"

	"github.com/holiman/uint256"

	"github.com/autonity/autonity/common"
)

type Address common.Address

func NewAddressFromBytes(bytes []byte) Address {
	return Address(common.BytesToAddress(bytes))
}

func (a *Address) ToCommonAddress() common.Address {
	return common.Address(*a)
}

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
