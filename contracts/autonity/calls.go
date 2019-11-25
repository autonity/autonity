package autonity

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
	"math"
	"math/big"
	"reflect"
)

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
	for _, v := range chain.Config().AutonityContractConfig.Users {
		validators = append(validators, v.Address)
		enodes = append(enodes, v.Enode)
		accTypes = append(accTypes, big.NewInt(int64(v.Type.GetID())))
		participantStake = append(participantStake, big.NewInt(int64(v.Stake)))
	}

	constructorParams, err := contractABI.Pack("",
		validators,
		enodes,
		accTypes,
		participantStake,
		chain.Config().AutonityContractConfig.Operator,
		new(big.Int).SetUint64(chain.Config().AutonityContractConfig.MinGasPrice))
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

	return contractAddress, nil
}

func (ac *Contract) UpdateAutonityContract(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) error {
	caller := ac.bc.Config().AutonityContractConfig.Deployer
	evm := ac.getEVM(header, caller, statedb)
	_, contractAddress, _, vmerr := evm.CreateWithAddress(sender, data, gas, value, ac.Address())
}

func (ac *Contract) AutonityContractCall(statedb *state.StateDB, header *types.Header, function string, result interface{}) error {
	caller := ac.bc.Config().AutonityContractConfig.Deployer
	contractABI, err := ac.abi()
	if err != nil {
		return err
	}

	gas := uint64(math.MaxUint64)
	evm := ac.getEVM(header, caller, statedb)

	input, err := contractABI.Pack(function)
	if err != nil {
		return err
	}

	ret, gas, vmerr := evm.Call(vm.AccountRef(caller), ac.Address(), input, gas, new(big.Int))
	if vmerr != nil {
		log.Error("Error Autonity Contract", "function", function)
		return vmerr
	}
	// if result's type is "raw" then bypass unpacking
	if reflect.TypeOf(result) == reflect.TypeOf(&raw{}) {
		log.Info("meme type")
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
	var upgradeReady bool
	err := ac.AutonityContractCall(state, header, "finalize", &upgradeReady)
	if err != nil {
		return false, err
	}
	return upgradeReady, nil
}

func (ac *Contract) callRetrieveState(statedb *state.StateDB, header *types.Header) ([]byte, error) {
	var state raw
	err := ac.AutonityContractCall(statedb, header, "retrieveState", &state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (ac *Contract) callRetrieveContract(state *state.StateDB, header *types.Header) ([]byte, []byte, error) {
	var bytecode []byte
	var abi []byte
	err := ac.AutonityContractCall(state, header, "retrieveContract", &[]interface{}{&bytecode, &abi})
	if err != nil {
		return nil, nil, err
	}
	return bytecode, abi, nil
}
