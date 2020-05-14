package acdefault

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/acdefault/generated"
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
