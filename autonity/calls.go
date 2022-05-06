package autonity

import (
	"math"
	"math/big"
	"reflect"
	"sort"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
)

type raw []byte

type Config struct {
	OperatorAccount common.Address `abi:"operatorAccount"`
	TreasuryAccount common.Address `abi:"treasuryAccount"`
	TreasuryFee     *big.Int       `abi:"treasuryFee"`
	MinBaseFee      *big.Int       `abi:"minBaseFee"`
	DelegationRate  *big.Int       `abi:"delegationRate"`
	EpochPeriod     *big.Int       `abi:"epochPeriod"`
	UnbondingPeriod *big.Int       `abi:"unbondingPeriod"`
	CommitteeSize   *big.Int       `abi:"committeeSize"`
	ContractVersion string         `abi:"contractVersion"`
	BlockPeriod     *big.Int       `abi:"blockPeriod"`
}

func DeployContract(abi *abi.ABI, genesisConfig *params.AutonityContractGenesis, evm *vm.EVM) error {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(genesisConfig.Bytecode)
	defaultVersion := "v0.0.0"

	contractConfig := Config{
		OperatorAccount: genesisConfig.Operator,
		TreasuryAccount: genesisConfig.Treasury,
		TreasuryFee:     new(big.Int).SetUint64(genesisConfig.TreasuryFee),
		MinBaseFee:      new(big.Int).SetUint64(genesisConfig.MinBaseFee),
		DelegationRate:  new(big.Int).SetUint64(genesisConfig.DelegationRate),
		EpochPeriod:     new(big.Int).SetUint64(genesisConfig.EpochPeriod),
		UnbondingPeriod: new(big.Int).SetUint64(genesisConfig.UnbondingPeriod),
		CommitteeSize:   new(big.Int).SetUint64(genesisConfig.MaxCommitteeSize),
		ContractVersion: defaultVersion,
		BlockPeriod:     new(big.Int).SetUint64(genesisConfig.BlockPeriod),
	}

	vals := genesisConfig.GetValidatorsCopy()

	constructorParams, err := abi.Pack("", vals, contractConfig)
	if err != nil {
		log.Error("contractABI.Pack returns err", "err", err)
		return err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Autonity contract
	_, _, _, vmerr := evm.Create(vm.AccountRef(Deployer), data, gas, value)
	if vmerr != nil {
		log.Error("DeployAutonityContract evm.Create", "err", vmerr)
		return vmerr
	}
	log.Info("Deployed Autonity Contract", "Address", ContractAddress.String())

	return nil
}

func (ac *Contract) updateAutonityContract(header *types.Header, statedb *state.StateDB, bytecode []byte) error {
	evm := ac.evmProvider.EVM(header, Deployer, statedb)
	_, _, _, vmerr := evm.Replace(vm.AccountRef(Deployer), bytecode, ContractAddress)
	if vmerr != nil {
		log.Error("updateAutonityContract evm.Create", "err", vmerr)
		return vmerr
	}
	return nil
}

// AutonityContractCall calls the specified function of the autonity contract
// with the given args, and returns the output unpacked into the result
// interface.
func (ac *Contract) AutonityContractCall(statedb *state.StateDB, header *types.Header, function string, result interface{}, args ...interface{}) error {

	packedArgs, err := ac.contractABI.Pack(function, args...)
	if err != nil {
		return err
	}

	ret, err := ac.CallContractFunc(statedb, header, function, packedArgs)
	if err != nil {
		return err
	}
	// if result's type is "raw" then bypass unpacking
	if reflect.TypeOf(result) == reflect.TypeOf(&raw{}) {
		rawPtr := result.(*raw)
		*rawPtr = raw(ret)
		return nil
	}

	if err := ac.contractABI.UnpackIntoInterface(result, function, ret); err != nil {
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
func (ac *Contract) CallContractFunc(statedb *state.StateDB, header *types.Header, function string, packedArgs []byte) ([]byte, error) {
	gas := uint64(math.MaxUint64)
	evm := ac.evmProvider.EVM(header, Deployer, statedb)
	packedResult, _, err := evm.Call(vm.AccountRef(Deployer), ContractAddress, packedArgs, gas, new(big.Int))
	return packedResult, err
}

func (ac *Contract) callGetCommitteeEnodes(state *state.StateDB, header *types.Header) (*types.Nodes, error) {
	var returnedEnodes []string
	err := ac.AutonityContractCall(state, header, "getCommitteeEnodes", &returnedEnodes)
	if err != nil {
		return nil, err
	}
	return types.NewNodes(returnedEnodes), nil
}

func (ac *Contract) callGetMinimumBaseFee(state *state.StateDB, header *types.Header) (uint64, error) {
	minBaseFee := new(big.Int)
	err := ac.AutonityContractCall(state, header, "getMinimumBaseFee", &minBaseFee)
	if err != nil {
		return 0, err
	}
	return minBaseFee.Uint64(), nil
}

func (ac *Contract) callGetProposer(state *state.StateDB, header *types.Header, height uint64, round int64) common.Address {
	var proposer common.Address
	h := new(big.Int).SetUint64(height)
	r := new(big.Int).SetInt64(round)
	err := ac.AutonityContractCall(state, header, "getProposer", &proposer, h, r)
	if err != nil {
		log.Error("get proposer failed from contract.", "error", err)
		return common.Address{}
	}
	return proposer
}

func (ac *Contract) callFinalize(state *state.StateDB, header *types.Header, blockGas *big.Int) (bool, types.Committee, error) {

	var updateReady bool
	var committee types.Committee

	err := ac.AutonityContractCall(state, header, "finalize", &[]interface{}{&updateReady, &committee}, blockGas)
	if err != nil {
		return false, nil, err
	}
	sort.Sort(committee)
	// submit the final reward distribution metrics.
	//ac.metrics.SubmitRewardDistributionMetrics(&v, header.Number.Uint64())
	return updateReady, committee, nil
}

func (ac *Contract) callRetrieveContract(state *state.StateDB, header *types.Header) ([]byte, string, error) {
	var bytecode []byte
	var abi string
	err := ac.AutonityContractCall(state, header, "getNewContract", &[]interface{}{&bytecode, &abi})
	if err != nil {
		return nil, "", err
	}
	return bytecode, abi, nil
}

func (ac *Contract) callSetMinimumBaseFee(state *state.StateDB, header *types.Header, price *big.Int) error {
	// Needs to be refactored somehow
	gas := uint64(0xFFFFFFFF)
	evm := ac.evmProvider.EVM(header, Deployer, state)

	input, err := ac.contractABI.Pack("setMinimumGasPrice")
	if err != nil {
		return err
	}

	_, _, vmerr := evm.Call(vm.AccountRef(Deployer), ContractAddress, input, gas, price)
	if vmerr != nil {
		log.Error("Error Autonity Contract getMinimumGasPrice()")
		return vmerr
	}
	return nil
}
