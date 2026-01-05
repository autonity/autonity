package fba

import (
	"math/big"

	"github.com/autonity/autonity/common"
)

type OrderIntent struct {
	User      common.Address
	ProductID common.Hash
	Side     Side
	Price     *big.Int
	Quantity  *big.Int
	Nonce     uint64
	Signature []byte
}
