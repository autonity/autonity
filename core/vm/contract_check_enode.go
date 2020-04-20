package vm

import (
	"fmt"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/clearmatics/autonity/params"
)

// checkEnode implemented as a native contract.
type checkEnode struct{}

func (c checkEnode) RequiredGas(_ []byte) uint64 {
	return params.EnodeCheckGas
}
func (c checkEnode) Run(input []byte) ([]byte, error) {
	if len(input) == 0 {
		panic(fmt.Errorf("invalid enode - empty"))
	}
	input = common.TrimPrefixAndSuffix(input, []byte("enode:"), []byte{'\x00'})
	nodeStr := string(input)

	if _, err := enode.ParseV4SkipResolve(nodeStr); err != nil {
		return false32Byte, fmt.Errorf("invalid enode %q: %v", nodeStr, err)
	}
	return true32Byte, nil
}
