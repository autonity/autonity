package backend

import (
	"math/big"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

// Instantiates a new EVM object which is required when creating or calling a deployed contract
func (sb *backend) getEVM(chain consensus.ChainReader, header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {

	coinbase, _ := sb.Author(header)
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     core.GetHashFn(header, chain),
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    new(big.Int).SetUint64(0x0),
	}
	evm := vm.NewEVM(evmContext, statedb, chain.Config(), *sb.vmConfig)
	return evm
}

// deployContract deploys the contract contained within the genesis field bytecode
func (sb *backend) deployContract(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(sb.config.Bytecode)
	evm := sb.getEVM(chain, header, sb.config.Deployer, statedb)
	sender := vm.AccountRef(sb.config.Deployer)

	var validators common.Addresses
	validators, _ = sb.retrieveSavedValidators(1)
	sort.Sort(validators)
	//We need to append to data the constructor's parameters
	//That should always be genesis validators

	SomaAbi, err := abi.JSON(strings.NewReader(sb.config.ABI))
	if err != nil {
		return common.Address{}, err
	}
	constructorParams, err := SomaAbi.Pack("", validators)
	if err != nil {
		return common.Address{}, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Soma validator governance contract
	_, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)

	if vmerr != nil {
		log.Error("Error Soma Governance Contract deployment")
		return contractAddress, vmerr
	}

	log.Info("Deployed Soma Governance Contract", "Address", contractAddress.String())

	return contractAddress, nil
}

func (sb *backend) contractGetValidators(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) ([]common.Address, error) {
	sender := vm.AccountRef(sb.config.Deployer)
	gas := uint64(0xFFFFFFFF)
	evm := sb.getEVM(chain, header, sb.config.Deployer, statedb)
	SomaAbi, err := abi.JSON(strings.NewReader(sb.config.ABI))
	input, err := SomaAbi.Pack("getValidators")

	if err != nil {
		return nil, err
	}
	value := new(big.Int).SetUint64(0x00)
	//A standard call is issued - we leave the possibility to modify the state
	ret, gas, vmerr := evm.Call(sender, sb.somaContract, input, gas, value)
	if vmerr != nil {
		log.Error("Error Soma Governance Contract GetValidators()")
		return nil, vmerr
	}

	var addresses []common.Address
	if err := SomaAbi.Unpack(&addresses, "getValidators", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getValidators returned value")
		return nil, err
	}

	sortableAddresses := common.Addresses(addresses) // Haven't found a way to avoid this
	sort.Sort(sortableAddresses)
	return sortableAddresses, nil
}

/*
// contractCallAddress queries the Soma contract, for any functions that take and address as argument.
// Returns true/false if the the address is an active validator and false if not.
func (sb *backend) contractCallAddress(functionSig string, chain consensus.ChainReader, userAddr common.Address, contractAddress common.Address, header *types.Header, db ethdb.Database) (bool, error) {
	// Instantiate new state database
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(header.Root, sdb)

	sender := vm.AccountRef(userAddr)
	gas := uint64(0xFFFFFFFF)
	evm := sb.getEVM(chain, header, userAddr, statedb)

	// Pad address for ABI encoding
	encodedAddress := [32]byte{}
	copy(encodedAddress[12:], userAddr[:])
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()[:4]
	inputData := append(input[:], encodedAddress[:]...)

	// Call ActiveValidators()
	ret, gas, vmerr := evm.StaticCall(sender, contractAddress, inputData, gas)
	if len(ret) == 0 {
		log.Info("contractCallAddress(): No return value", "Block", header.Number.Int64())
		return false, consensus.ErrPrunedAncestor
	}

	if vmerr != nil {
		return false, vmerr
	}

	const def = `[{ "name" : "method", "outputs": [{ "type": "bool" }] }]`
	funcAbi, err := abi.JSON(strings.NewReader(def))
	if err != nil {
		return false, vmerr
	}

	var output bool
	err = funcAbi.Unpack(&output, "method", ret)
	if err != nil {
		return false, err
	}

	return output, nil
}

// callContractDifficulty calls contract to find the difficulty for a specific validator returns an int either 1 or 2
func (sb *backend) callContractDifficulty(chain consensus.ChainReader, userAddr common.Address, contractAddress common.Address, header *types.Header, db ethdb.Database) (*big.Int, error) {
	// Signature of function being called defined by Soma interface
	functionSig := "calculateDifficulty(address)"

	// Instantiate new state database
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(header.Root, sdb)

	sender := vm.AccountRef(userAddr)
	gas := uint64(0xFFFFFFFF)
	evm := getEVM(chain, header, userAddr, userAddr, statedb)

	// Pad address for ABI encoding
	encodedAddress := [32]byte{}
	copy(encodedAddress[12:], userAddr[:])
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()[:4]
	inputData := append(input[:], encodedAddress[:]...)

	// Call ActiveValidators()
	ret, gas, vmerr := evm.StaticCall(sender, contractAddress, inputData, gas)
	if vmerr != nil {
		return big.NewInt(1), vmerr
	}

	const def = `[{"name" : "int", "constant" : false, "outputs": [ { "type": "uint256" } ]}]`
	funcAbi, err := abi.JSON(strings.NewReader(def))
	if err != nil {
		return big.NewInt(1), err
	}

	// marshal int
	var Int *big.Int
	err = funcAbi.Unpack(&Int, "int", ret)
	if err != nil {
		return big.NewInt(1), consensus.ErrPrunedAncestor
	}

	return Int, nil
}

// updateGovernance when a validator attempts to submit a block the
func (sb *backend) updateGovernance(chain consensus.ChainReader, userAddr common.Address, contractAddress common.Address, header *types.Header, statedb *state.StateDB) error {
	// Signature of function being called defined by Soma interface
	functionSig := "UpdateGovernance()"

	sender := vm.AccountRef(userAddr)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	evm := getEVM(chain, header, userAddr, userAddr, statedb)

	// Pad address for ABI encoding
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()

	// Call ActiveValidators()
	_, gas, vmerr := evm.Call(sender, contractAddress, input, gas, value)

	if vmerr != nil {
		return vmerr
	}

	return nil

}

// contractCall calls a contract function with the input functionSig this MUST NOT take any arguments
func (sb *backend) contractCall(functionSig string, chain consensus.ChainReader, userAddr common.Address, contractAddress common.Address, header *types.Header, db ethdb.Database) ([]byte, error) {
	// Instantiate new state database
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(header.Root, sdb)

	sender := vm.AccountRef(userAddr)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	evm := getEVM(chain, header, userAddr, userAddr, statedb)

	// Pad address for ABI encoding
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()

	// Call ActiveValidators()
	ret, gas, vmerr := evm.Call(sender, contractAddress, input, gas, value)
	if vmerr != nil {
		return nil, vmerr
	}

	return ret, nil
}
*/
