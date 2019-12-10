package autonity

import (
	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
	"math"
	"math/big"
	"reflect"
	"strings"
)

/*
 * ContractState is a unified structure to represent the autonity contract state.
 * By using a unified structure, the new state meta introduced in the Autonity.sol
 * should be synced with this structure.
 */
type ContractState struct {
	Users []common.Address `abi:"users"`
	Enodes []string `abi:"enodes"`
	Types []*big.Int `abi:"types"`
	Stakes []*big.Int `abi:"stakes"`
	CommissionRates []*big.Int `abi:"commisionrates"`
	Operator common.Address `abi:"operator"`
	Deployer common.Address `abi:"deployer"`
	MinGasPrice *big.Int `abi:"mingasprice"`
	BondingPeriod *big.Int `abi:"bondingperiod"`
}

type raw []byte

//// Instantiates a new EVM object which is required when creating or calling a deployed contract
func (ac *Contract) getEVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
	coinbase, _ := types.Ecrecover(header)
	evmContext := vm.Context{
		CanTransfer: ac.canTransfer,
		Transfer:    ac.transfer,
		GetHash:     ac.GetHashFn(header, ac.bc),
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        new(big.Int).SetUint64(header.Time),
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    new(big.Int).SetUint64(0x0),
	}
	vmConfig := *ac.bc.GetVMConfig()
	evm := vm.NewEVM(evmContext, statedb, ac.bc.Config(), vmConfig)
	return evm
}

// deployContract deploys the contract contained within the genesis field bytecode
func (ac *Contract) DeployAutonityContract(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(chain.Config().AutonityContractConfig.Bytecode)
	evm := ac.getEVM(header, chain.Config().AutonityContractConfig.Deployer, statedb)
	sender := vm.AccountRef(chain.Config().AutonityContractConfig.Deployer)

	//todo do we need it?
	//validators, err = ac.SavedValidatorsRetriever(1)
	//sort.Sort(validators)

	//We need to append to data the constructor's parameters
	//That should always be genesis validators

	contractABI, err := ac.abi()

	if err != nil {
		log.Error("abi.JSON returns err", "err", err)
		return common.Address{}, err
	}

	ln := len(chain.Config().AutonityContractConfig.GetValidatorUsers())
	validators := make(common.Addresses, 0, ln)
	enodes := make([]string, 0, ln)
	accTypes := make([]*big.Int, 0, ln)
	participantStake := make([]*big.Int, 0, ln)
	commissionRate := make([]*big.Int, 0, ln)

	// Default bond period is 100.
	defaultBondPeriod := big.NewInt(100)

	for _, v := range chain.Config().AutonityContractConfig.Users {
		validators = append(validators, v.Address)
		enodes = append(enodes, v.Enode)
		accTypes = append(accTypes, big.NewInt(int64(v.Type.GetID())))
		participantStake = append(participantStake, big.NewInt(int64(v.Stake)))

		//TODO: default commission rate is 0, should use a config file...
		commissionRate = append(commissionRate, common.Big0)
	}

	constructorParams, err := contractABI.Pack("",
		validators,
		enodes,
		accTypes,
		participantStake,
		commissionRate,
		chain.Config().AutonityContractConfig.Operator,
		chain.Config().AutonityContractConfig.Deployer,
		new(big.Int).SetUint64(chain.Config().AutonityContractConfig.MinGasPrice),
		defaultBondPeriod)
	if err != nil {
		log.Error("contractABI.Pack returns err", "err", err)
		return common.Address{}, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Autonity contract
	_, contractAddress, _, vmerr := evm.Create(sender, data, gas, value)
	if vmerr != nil {
		log.Error("evm.Create returns err", "err", vmerr)
		return contractAddress, vmerr
	}
	ac.Lock()
	ac.address = contractAddress
	ac.Unlock()
	log.Info("Deployed Autonity Contract", "Address", contractAddress.String())

	///////////////////////////////////////////////////////////////
	// only for testing, need to be removed once feature done!
	contractBin := chain.Config().AutonityContractConfig.Bytecode
	contractAbi := chain.Config().AutonityContractConfig.ABI
	contractErr := ac.callSetContractForTesting(statedb, header, contractBin, contractAbi)
	if contractErr != nil {
		log.Error("set contract binary failed", "err", contractErr)
		return contractAddress, contractErr
	}
	////////////////////////////////////////////////////////////////

	return contractAddress, nil
}

func (ac *Contract) UpdateAutonityContractV2(header *types.Header, statedb *state.StateDB, bytecode string, newAbi string, cs ContractState) error {

	if header == nil || statedb == nil || len(bytecode) == 0 || len(newAbi) == 0 {
		return ErrWrongParameter
	}

	caller := ac.bc.Config().AutonityContractConfig.Deployer
	evm := ac.getEVM(header, caller, statedb)
	contractBytecode := common.Hex2Bytes(bytecode)

	// get new abi, and prepare the constructor of new contract.
	newContractABI, err := abi.JSON(strings.NewReader(newAbi))
	if err != nil {
		return err
	}

	//TODO: Better to fix and use Youssef's direct byte packing upgrade, otherwise the construction is less flexible.
	constructorParams, err := newContractABI.Pack("",
		cs.Users,
		cs.Enodes,
		cs.Types,
		cs.Stakes,
		cs.CommissionRates,
		cs.Operator,
		cs.Deployer,
		cs.MinGasPrice,
		cs.BondingPeriod)
	if err != nil {
		log.Error("contractABI.Pack returns err", "err", err)
		return err
	}
	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)
	_, _, _, vmerr := evm.CreateWithAddress(vm.AccountRef(caller), data, gas, value, ac.Address())
	if vmerr != nil {
		log.Error("evm.Create returns err", "err", vmerr)
		return vmerr
	}
	return nil
}

func (ac *Contract) UpdateAutonityContract(header *types.Header, statedb *state.StateDB, bytecode string, abi string, state []byte) error {
	caller := ac.bc.Config().AutonityContractConfig.Deployer
	evm := ac.getEVM(header, caller, statedb)
	contractBytecode := common.Hex2Bytes(bytecode)
	data := append(contractBytecode, state...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)
	_, _, _, vmerr := evm.CreateWithAddress(vm.AccountRef(caller), data, gas, value, ac.Address())
	if vmerr != nil {
		log.Error("evm.Create returns err", "err", vmerr)
		return vmerr
	}
	return nil
}

func (ac *Contract) AutonityContractCall(statedb *state.StateDB, header *types.Header, function string, result interface{}, args ...interface{}) error {
	caller := ac.bc.Config().AutonityContractConfig.Deployer
	contractABI, err := ac.abi()
	if err != nil {
		return err
	}

	gas := uint64(math.MaxUint64)
	evm := ac.getEVM(header, caller, statedb)

	input, err := contractABI.Pack(function, args...)
	if err != nil {
		return err
	}

	ret, _, vmerr := evm.Call(vm.AccountRef(caller), ac.Address(), input, gas, new(big.Int))
	if vmerr != nil {
		log.Error("Error Autonity Contract", "function", function)
		return vmerr
	}
	// if result's type is "raw" then bypass unpacking
	if reflect.TypeOf(result) == reflect.TypeOf(&[]byte{}) {
		log.Info("meme type")
		// TODO: copy the slice of byte into slice interface. GO's typing system is complicated, using Unpack for now.
		//reflect.ValueOf(result).SetBytes(ret)
		//result = append([]byte(nil), ret...)
		//copy(result, ret)
		return nil
	}

	log.Info(reflect.TypeOf(result).String())
	if err := contractABI.Unpack(result, function, ret); err != nil {
		log.Error("Could not unpack returned value", "function", function)
		return err
	}

	return nil
}

func (ac *Contract) callGetWhitelist(state *state.StateDB, header *types.Header) (*types.Nodes, error) {
	var returnedEnodes []string
	err := ac.AutonityContractCall(state, header, "getWhitelist", &returnedEnodes)
	if err != nil {
		return nil, err
	}
	return types.NewNodes(returnedEnodes, false), nil
}

func (ac *Contract) callGetMinimumGasPrice(state *state.StateDB, header *types.Header) (uint64, error) {
	minGasPrice := new(big.Int)
	err := ac.AutonityContractCall(state, header, "getMinimumGasPrice", &minGasPrice)
	if err != nil {
		return 0, err
	}
	return minGasPrice.Uint64(), nil
}

func (ac *Contract) callFinalize(state *state.StateDB, header *types.Header, blockGas *big.Int) (bool, error) {
	v := RewardDistributionMetaData{}
	err := ac.AutonityContractCall(state, header, "finalize", &v, blockGas)
	if err != nil {
		return false, err
	}

	// submit the final reward distribution metrics.
	ac.metrics.SubmitRewardDistributionMetrics(&v, header.Number.Uint64())
	return v.Result, nil
}

func (ac *Contract) callRetrieveStateV2(statedb *state.StateDB, header *types.Header) (ContractState, error) {
	cs := ContractState{}
	err := ac.AutonityContractCall(statedb, header, "retrieveStateV2", &cs)
	if err != nil {
		return cs, err
	}
	return cs, nil
}

func (ac *Contract) callRetrieveState(statedb *state.StateDB, header *types.Header) ([]byte, error) {
	var state raw

	err := ac.AutonityContractCall(statedb, header, "retrieveState", &state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (ac *Contract) callRetrieveContract(state *state.StateDB, header *types.Header) (string, string, error) {
	var bytecode string
	var abi string
	err := ac.AutonityContractCall(state, header, "retrieveContract", &[]interface{}{&bytecode, &abi})
	if err != nil {
		return "", "", err
	}
	return bytecode, abi, nil
}

func (ac *Contract) callSetContractForTesting(state *state.StateDB, header *types.Header, byteCode string, abi string) error {
	if state == nil || header == nil {
		return ErrWrongParameter
	}

	var result bool
	err := ac.AutonityContractCall(state, header, "upgradeContract", &result, byteCode, abi)
	if err != nil {
		return err
	}
	if result {
		return nil
	}
	return ErrAutonityContract
}

func (ac *Contract) callSetMinimumGasPrice(state *state.StateDB, header *types.Header, price *big.Int) error {
	// Needs to be refactored somehow
	deployer := ac.bc.Config().AutonityContractConfig.Deployer
	sender := vm.AccountRef(deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, deployer, state)

	ABI, err := ac.abi()
	if err != nil {
		return err
	}

	input, err := ABI.Pack("setMinimumGasPrice")
	if err != nil {
		return err
	}

	_, _, vmerr := evm.Call(sender, ac.Address(), input, gas, price)
	if vmerr != nil {
		log.Error("Error Autonity Contract getMinimumGasPrice()")
		return vmerr
	}
	return nil
}
