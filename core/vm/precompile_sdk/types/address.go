package types

import "github.com/autonity/autonity/common"

type Address common.Address

func NewAddressFromBytes(bytes []byte) Address {
	return Address(common.BytesToAddress(bytes))
}

func (a *Address) ToCommonAddress() common.Address {
	return common.Address(*a)
}
