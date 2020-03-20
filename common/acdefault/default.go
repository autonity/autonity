package acdefault

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/acdefault/generated"
)

func Deployer() common.Address {
	return common.HexToAddress("0x1336000000000000000000000000000000000000")
}

func Governance() common.Address {
	return common.HexToAddress("0x1336000000000000000000000000000000000000")
}

func Bytecode() string {
	return generated.Bytecode
}

func ABI() string {
	return generated.Abi
}

// Generate the abi and bytecode

//go:generate mkdir -p generated
//go:generate bash -c "docker run --rm -v $(pwd)/../../contracts/autonity/contract/contracts:/contracts -v $(pwd)/generated:/output ethereum/solc:0.5.1 --overwrite --abi --bin -o /output /contracts/Autonity.sol"

//go:generate echo Generating generated/bytecode.go
//go:generate bash -c "echo 'package generated' > generated/bytecode.go"
//go:generate bash -c "echo '// Code generated DO NOT EDIT.' >> generated/bytecode.go"
//go:generate bash -c "echo -n 'const Bytecode = \"' >> generated/bytecode.go"
//go:generate bash -c "cat  generated/Autonity.bin >> generated/bytecode.go"
//go:generate bash -c "echo '\"' >> generated/bytecode.go"
//go:generate gofmt -s -w generated/bytecode.go

//go:generate echo Generating generated/abi.go
//go:generate bash -c "echo 'package generated' > generated/abi.go"
//go:generate bash -c "echo '// Code generated DO NOT EDIT.' >> generated/abi.go"
//go:generate bash -c "echo -n 'const Abi = `' >> generated/abi.go"
//go:generate bash -c "cat  generated/Autonity.abi | json_pp  >> generated/abi.go"
//go:generate bash -c "echo '`' >> generated/abi.go"
//go:generate gofmt -s -w generated/abi.go
