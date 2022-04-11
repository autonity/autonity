package acdefault

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/acdefault/generated"
)

func Governance() common.Address {
	return common.HexToAddress("0x1336000000000000000000000000000000000000")
}

func Bytecode() string {
	return generated.Bytecode
}

func ABI() string {
	return generated.Abi
}
