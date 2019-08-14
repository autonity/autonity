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
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
	"math/big"
	"sort"
	"strings"
)

//go:generate abigen --sol contract/autonity/contract/contracts/Autonity.sol --exc contract/autonity/contract/contracts/SafeMath.sol:SafeMath --pkg autonity --out contract/autonity2.go

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
	ByteCode                 string
	ABI                      string
	Deployer                 common.Address
	Address                  common.Address
	bc                       Blockchainer
	SavedValidatorsRetriever func(i uint64) ([]common.Address, error)

	canTransfer func(db vm.StateDB, addr common.Address, amount *big.Int) bool
	transfer    func(db vm.StateDB, sender, recipient common.Address, amount *big.Int)
	GetHashFn   func(ref *types.Header, chain ChainContext) func(n uint64) common.Hash
}

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
	evm := vm.NewEVM(evmContext, statedb, ac.bc.Config(), *ac.bc.GetVMConfig())
	return evm
}

// deployContract deploys the contract contained within the genesis field bytecode
func (ac *AutonityContract) DeployAutonityContract(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(ac.ByteCode)
	evm := ac.getEVM(header, ac.Deployer, statedb)
	sender := vm.AccountRef(ac.Deployer)

	var validators common.Addresses
	var err error
	validators, err = ac.SavedValidatorsRetriever(1)
	sort.Sort(validators)
	fmt.Println("contracts/autonity/autonity.go:101 DeployAutonityContract", err, validators)
	//We need to append to data the constructor's parameters
	//That should always be genesis validators

	fmt.Println("ac.ABI", ac.ABI)
	contractABI, err := abi.JSON(strings.NewReader(ac.ABI))
	if err != nil {
		fmt.Println("ABI", ac.ABI)
		fmt.Println("a1", err)
		return common.Address{}, err
	}

	constructorParams, err := contractABI.Pack("", validators)
	if err != nil {
		fmt.Println("a2", err)
		return common.Address{}, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Soma validator governance contract
	_, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)
	if vmerr != nil {
		fmt.Println("contracts/autonity/autonity.go:127 deployment err", err)
		return contractAddress, vmerr
	}

	fmt.Println("Successful deployment")
	log.Info("Deployed Autonity Contract", "Address", contractAddress.String())

	return contractAddress, nil
}

func (ac *AutonityContract) ContractGetValidators(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) ([]common.Address, error) {
	sender := vm.AccountRef(ac.Deployer)
	gas := uint64(0xFFFFFFFF)
	evm := ac.getEVM(header, ac.Deployer, statedb)

	somaAbi, err := abi.JSON(strings.NewReader(ac.ABI))
	if err != nil {
		return nil, err
	}

	input, err := somaAbi.Pack("getValidators")
	if err != nil {
		return nil, err
	}

	value := new(big.Int).SetUint64(0x00)
	//A standard call is issued - we leave the possibility to modify the state
	ret, gas, vmerr := evm.Call(sender, ac.Address, input, gas, value)
	if vmerr != nil {
		log.Error("Error Soma Governance Contract GetValidators()")
		return nil, vmerr
	}

	var addresses []common.Address
	fmt.Println("ret", ret)
	if err := somaAbi.Unpack(&addresses, "getValidators", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getValidators returned value")
		return nil, err
	}

	sortableAddresses := common.Addresses(addresses)
	sort.Sort(sortableAddresses)
	return sortableAddresses, nil
}

var GlienickeContractError = errors.New("could not call Glienicke contract")

func (ac *AutonityContract) UpdateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	newWhitelist, err := ac.GetWhitelist(block, state)
	if err != nil {
		log.Error("could not call Glienicke contract", "err", err)
		return GlienickeContractError
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
		newWhitelist, err = ac.callContract(db, block.Header())
	}

	return newWhitelist, err
}

//blockchain

func (ac *AutonityContract) callContract(state *state.StateDB, header *types.Header) (*types.Nodes, error) {
	// Needs to be refactored somehow
	deployer := ac.Deployer
	if deployer == (common.Address{}) {
		deployer = params.GlienickeDefaultDeployer
	}

	var contractABI = ac.ABI
	if contractABI == "" {
		contractABI = params.GlienickeDefaultABI
	}

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

	glienickeAddress := crypto.CreateAddress(deployer, 0)

	ret, gas, vmerr := evm.StaticCall(sender, glienickeAddress, input, gas)
	if vmerr != nil {
		log.Error("Error Glienicke Contract getWhitelist()")
		return nil, vmerr
	}

	var returnedEnodes []string
	if err := ABI.Unpack(&returnedEnodes, "getWhitelist", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getWhitelist returned value")
		return nil, err
	}

	return types.NewNodes(returnedEnodes, false), nil
}

// Instantiates a new EVM object which is required when creating or calling a deployed contract
//func (sb *AutonityContract) getEVM(chain consensus.ChainReader, header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
//
//	coinbase, _ := sb.Author(header)
//	evmContext := vm.Context{
//		CanTransfer: core.CanTransfer,
//		Transfer:    core.Transfer,
//		GetHash:     core.GetHashFn(header, chain),
//		Origin:      origin,
//		Coinbase:    coinbase,
//		BlockNumber: header.Number,
//		Time:        header.Time,
//		GasLimit:    header.GasLimit,
//		Difficulty:  header.Difficulty,
//		GasPrice:    new(big.Int).SetUint64(0x0),
//	}
//	evm := vm.NewEVM(evmContext, statedb, chain.Config(), *sb.vmConfig)
//	return evm
//}

//deployContract deploys the contract contained within the genesis field bytecode
//func (bc *AutonityContract) DeployGlienickeContract(state *state.StateDB, header *types.Header) (*types.Nodes, common.Address, error) {
//	//if bytecode or abi is missing use default one
//	glienickeByteCode := bc.chainConfig.GlienickeBytecode
//	glienickeABI := bc.chainConfig.GlienickeABI
//	if bc.chainConfig.GlienickeBytecode == "" || bc.chainConfig.GlienickeABI == "" {
//		glienickeByteCode = params.GlienickeDefaultBytecode
//		glienickeABI = params.GlienickeDefaultABI
//	}
//	bc.chainConfig.GlienickeABI = glienickeABI
//
//	//Same for deployer
//	glienickeDeployer := bc.chainConfig.GlienickeDeployer
//	if glienickeDeployer == (common.Address{}) {
//		glienickeDeployer = params.GlienickeDefaultDeployer
//	}
//
//	// Convert the contract bytecode from hex into bytes
//	contractBytecode := common.Hex2Bytes(glienickeByteCode)
//	evm := bc.getEVM(header, glienickeDeployer, state)
//	sender := vm.AccountRef(glienickeDeployer)
//
//	enodesWhitelist := rawdb.ReadEnodeWhitelist(bc.db, false)
//
//	glienickeAbi, err := abi.JSON(strings.NewReader(glienickeABI))
//	if err != nil {
//		return nil, common.Address{}, err
//	}
//
//	sort.Strings(enodesWhitelist.StrList)
//	constructorParams, err := glienickeAbi.Pack("", enodesWhitelist.StrList)
//	if err != nil {
//		return nil, common.Address{}, err
//	}
//
//	data := append(contractBytecode, constructorParams...)
//	gas := uint64(0xFFFFFFFF)
//	value := new(big.Int).SetUint64(0x00)
//
//	// Deploy the Glienicke validator governance contract
//	_, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)
//	if vmerr != nil {
//		log.Error("Error Glienicke Contract deployment")
//		return nil, common.Address{}, vmerr
//	}
//
//	log.Info("Deployed Glienicke Contract", "Address", contractAddress.String())
//
//	return enodesWhitelist, contractAddress, nil
//}
