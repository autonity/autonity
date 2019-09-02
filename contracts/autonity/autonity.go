package autonity

import (
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
	"math/big"
	"sort"
	"strings"
)

func NewAutonityContract(
	bc Blockchainer,
	canTransfer func(db vm.StateDB, addr common.Address, amount *big.Int) bool,
	transfer func(db vm.StateDB, sender, recipient common.Address, amount *big.Int),
	GetHashFn func(ref *types.Header, chain ChainContext) func(n uint64) common.Hash,
) *AutonityContract {
	return &AutonityContract{
		bc:          bc,
		canTransfer: canTransfer,
		transfer:    transfer,
		GetHashFn:   GetHashFn,
		//SavedValidatorsRetriever: SavedValidatorsRetriever,

	}
}

type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}
type Blockchainer interface {
	ChainContext
	GetVMConfig() *vm.Config
	Config() *params.ChainConfig

	UpdateEnodeWhitelist(newWhitelist *types.Nodes)
	ReadEnodeWhitelist(openNetwork bool) *types.Nodes
}

type AutonityContract struct {
	Address                  common.Address
	bc                       Blockchainer
	SavedValidatorsRetriever func(i uint64) ([]common.Address, error)

	canTransfer func(db vm.StateDB, addr common.Address, amount *big.Int) bool
	transfer    func(db vm.StateDB, sender, recipient common.Address, amount *big.Int)
	GetHashFn   func(ref *types.Header, chain ChainContext) func(n uint64) common.Hash
}

var Sl = vm.NewStructLogger(&vm.LogConfig{
	Debug: true,
})

//// Instantiates a new EVM object which is required when creating or calling a deployed contract
func (ac *AutonityContract) getEVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {

	coinbase, _ := types.Ecrecover(header)
	evmContext := vm.Context{
		CanTransfer: ac.canTransfer,
		Transfer:    ac.transfer,
		GetHash:     ac.GetHashFn(header, ac.bc),
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    new(big.Int).SetUint64(0x0),
	}
	vmConfig := *ac.bc.GetVMConfig()
	vmConfig.Debug = true
	vmConfig.Tracer = Sl
	evm := vm.NewEVM(evmContext, statedb, ac.bc.Config(), vmConfig)
	return evm
}

// deployContract deploys the contract contained within the genesis field bytecode
func (ac *AutonityContract) DeployAutonityContract(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(chain.Config().AutonityContractConfig.Bytecode)
	evm := ac.getEVM(header, chain.Config().AutonityContractConfig.Deployer, statedb)
	sender := vm.AccountRef(chain.Config().AutonityContractConfig.Deployer)

	//todo do we need it?
	//validators, err = ac.SavedValidatorsRetriever(1)
	//sort.Sort(validators)

	//We need to append to data the constructor's parameters
	//That should always be genesis validators
	contractABI, err := abi.JSON(strings.NewReader(chain.Config().AutonityContractConfig.ABI))
	if err != nil {
		return common.Address{}, err
	}

	ln := len(chain.Config().AutonityContractConfig.GetValidatorUsers())
	validators := make(common.Addresses, 0, ln)
	enodes := make([]string, 0, ln)
	accTypes := make([]*big.Int, 0, ln)
	participantStake := make([]*big.Int, 0, ln)
	for _, v := range chain.Config().AutonityContractConfig.GetValidatorUsers() {
		validators = append(validators, v.Address)
		enodes = append(enodes, v.Enode)
		accTypes = append(accTypes, big.NewInt(int64(v.Type)))
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
		return common.Address{}, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Soma validator governance contract
	_, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)
	if vmerr != nil {
		return contractAddress, vmerr
	}
	ac.Address = contractAddress
	log.Info("Deployed Autonity Contract", "Address", contractAddress.String())

	return contractAddress, nil
}

func (ac *AutonityContract) ContractGetValidators(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) ([]common.Address, error) {
	sender := vm.AccountRef(chain.Config().AutonityContractConfig.Deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, chain.Config().AutonityContractConfig.Deployer, statedb)
	contractABI, err := abi.JSON(strings.NewReader(chain.Config().AutonityContractConfig.ABI))
	if err != nil {
		fmt.Print("11")
		return nil, err
	}

	input, err := contractABI.Pack("GetValidators")
	if err != nil {
		fmt.Print("22")
		return nil, err
	}

	value := new(big.Int).SetUint64(0x00)
	//A standard call is issued - we leave the possibility to modify the state
	ret, gas, vmerr := evm.Call(sender, ac.Address, input, gas, value)
	if vmerr != nil {
		log.Error("Error Contract GetValidators()")
		fmt.Print("33")
		return nil, vmerr
	}

	var addresses []common.Address
	if err := contractABI.Unpack(&addresses, "GetValidators", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getValidators returned value", err)
		fmt.Print("44")
		return nil, err
	}

	sortableAddresses := common.Addresses(addresses)
	sort.Sort(sortableAddresses)
	return sortableAddresses, nil
}

var AutonityContractError = errors.New("could not call Autonity contract")

func (ac *AutonityContract) UpdateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	newWhitelist, err := ac.GetWhitelist(block, state)
	if err != nil {
		log.Error("could not call contract", "err", err)
		return AutonityContractError
	}

	ac.bc.UpdateEnodeWhitelist(newWhitelist)
	return nil
}

func (ac *AutonityContract) GetWhitelist(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	var (
		newWhitelist *types.Nodes
		err          error
	)

	if block.Number().Uint64() == 1 {
		// use genesis block whitelist
		newWhitelist = ac.bc.ReadEnodeWhitelist(false)
	} else {
		// call retrieveWhitelist contract function
		newWhitelist, err = ac.callGetWhitelist(db, block.Header())
	}

	return newWhitelist, err
}

//blockchain

func (ac *AutonityContract) callGetWhitelist(state *state.StateDB, header *types.Header) (*types.Nodes, error) {
	// Needs to be refactored somehow
	deployer := ac.bc.Config().AutonityContractConfig.Deployer
	var contractABI = ac.bc.Config().AutonityContractConfig.ABI

	sender := vm.AccountRef(deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, deployer, state)

	ABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	input, err := ABI.Pack("getWhitelist")
	if err != nil {
		return nil, err
	}

	ret, gas, vmerr := evm.StaticCall(sender, ac.Address, input, gas)
	if vmerr != nil {
		log.Error("Error Autonity Contract getWhitelist()")
		return nil, vmerr
	}

	var returnedEnodes []string
	if err := ABI.Unpack(&returnedEnodes, "getWhitelist", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getWhitelist returned value")
		return nil, err
	}

	return types.NewNodes(returnedEnodes, false), nil
}
