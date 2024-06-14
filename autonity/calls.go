package autonity

import (
	"math/big"
	"reflect"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

type raw []byte

func (c *EVMContract) replaceAutonityBytecode(header *types.Header, statedb vm.StateDB, bytecode []byte) error {
	evm := c.evmProvider(header, params.DeployerAddress, statedb)
	_, _, _, vmerr := evm.Replace(vm.AccountRef(params.DeployerAddress), bytecode, params.AutonityContractAddress)
	if vmerr != nil {
		log.Error("replaceAutonityBytecode evm.Create", "err", vmerr)
		return vmerr
	}
	return nil
}

// AutonityContractCall calls the specified function of the autonity contract
// with the given args, and returns the output unpacked into the result
// interface.
func AutonityContractCall(evm *vm.EVM, function string, result any, args ...any) error {
	packedArgs, err := generated.AutonityAbi.Pack(function, args...)
	if err != nil {
		return err
	}
	ret, _, err := evm.Call(vm.AccountRef(params.DeployerAddress), params.AutonityContractAddress, packedArgs, uint64(math.MaxUint64), common.Big0)
	if err != nil {
		return err
	}
	// if result's type is "raw" then bypass unpacking
	if reflect.TypeOf(result) == reflect.TypeOf(&raw{}) {
		rawPtr := result.(*raw)
		*rawPtr = ret
		return nil
	}
	if err := generated.AutonityAbi.UnpackIntoInterface(result, function, ret); err != nil {
		log.Error("Could not unpack returned value", "function", function)
		return err
	}
	return nil
}

// CallContractFunc creates an evm object, uses it to call the
// specified function of the autonity contract with packedArgs and returns the
// packed result. If there is an error making the evm call it will be returned.
// Callers should use the autonity contract ABI to pack and unpack the args and
// result.
func (c *EVMContract) CallContractFunc(state vm.StateDB, header *types.Header, contractAddress common.Address, packedArgs []byte) ([]byte, uint64, error) {
	gas := uint64(math.MaxUint64)
	evm := c.evmProvider(header, params.DeployerAddress, state)
	return evm.Call(vm.AccountRef(params.DeployerAddress), contractAddress, packedArgs, gas, new(big.Int))
}

func (c *EVMContract) CallContractFuncAs(state vm.StateDB, header *types.Header, contractAddress common.Address, origin common.Address, packedArgs []byte) ([]byte, error) {
	gas := uint64(math.MaxUint64)
	evm := c.evmProvider(header, origin, state)
	packedResult, _, err := evm.Call(vm.AccountRef(origin), contractAddress, packedArgs, gas, new(big.Int))
	return packedResult, err
}

func (c *AutonityContract) CallGetCommitteeEnodes(state vm.StateDB, header *types.Header, asACN bool) (*types.Nodes, error) {
	var returnedEnodes []string
	err := AutonityContractCall(c.evmProvider(header, params.DeployerAddress, state), "getCommitteeEnodes", &returnedEnodes)
	if err != nil {
		return nil, err
	}
	return types.NewNodes(returnedEnodes, asACN), nil
}

func (c *AutonityContract) CallGetMinimumBaseFee(state vm.StateDB, header *types.Header) (*big.Int, error) {
	minBaseFee := new(big.Int)
	err := AutonityContractCall(c.evmProvider(header, params.DeployerAddress, state), "getMinimumBaseFee", &minBaseFee)
	if err != nil {
		return nil, err
	}
	return minBaseFee, nil
}

func (c *AutonityContract) CallGetEpochPeriod(state vm.StateDB, header *types.Header) (*big.Int, error) {
	epochPeriod := new(big.Int)
	err := AutonityContractCall(c.evmProvider(header, params.DeployerAddress, state), "getEpochPeriod", &epochPeriod)
	if err != nil {
		return nil, err
	}
	return epochPeriod, nil
}

func (c *AutonityContract) CallFinalize(state vm.StateDB, header *types.Header) (bool, types.Committee, error) {
	var updateReady bool
	var committee types.Committee
	if err := AutonityContractCall(c.evmProvider(header, params.DeployerAddress, state), "finalize", &[]any{&updateReady, &committee}); err != nil {
		return false, nil, err
	}
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error())
	}
	return updateReady, committee, nil
}

func CallGetCommittee(evm *vm.EVM) ([]types.CommitteeMember, error) {
	var committee types.Committee
	if err := AutonityContractCall(evm, "getCommittee", &committee); err != nil {
		return nil, err
	}
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error()) //nolint
	}
	return committee, nil
}

func (c *AutonityContract) CallRetrieveContract(state vm.StateDB, header *types.Header) ([]byte, string, error) {
	var bytecode []byte
	var abi string
	if err := AutonityContractCall(c.evmProvider(header, params.DeployerAddress, state), "getNewContract", &[]any{&bytecode, &abi}); err != nil {
		return nil, "", err
	}
	return bytecode, abi, nil
}
