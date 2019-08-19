package backend

import (
	"math/big"
	"sort"
	"strings"

	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
)

// Instantiates a new EVM object which is required when creating or calling a deployed contract
func (sb *Backend) getEVM(chain consensus.ChainReader, header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {

	coinbase, _ := sb.Author(header)
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     core.GetHashFn(header, chain),
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        big.NewInt(int64(header.Time)),
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    new(big.Int).SetUint64(0x0),
	}
	evm := vm.NewEVM(evmContext, statedb, chain.Config(), *sb.vmConfig)
	return evm
}

// deployContract deploys the contract contained within the genesis field bytecode
func (sb *Backend) deployContract(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(sb.config.Bytecode)
	evm := sb.getEVM(chain, header, sb.config.Deployer, statedb)
	sender := vm.AccountRef(sb.config.Deployer)

	var validators common.Addresses
	validators, _ = sb.retrieveSavedValidators(1, chain)
	sort.Sort(validators)
	//We need to append to data the constructor's parameters
	//That should always be genesis validators

	somaAbi, err := abi.JSON(strings.NewReader(sb.config.ABI))
	if err != nil {
		return common.Address{}, err
	}

	constructorParams, err := somaAbi.Pack("", validators)
	if err != nil {
		return common.Address{}, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Soma validator governance contract
	_, contractAddress, _, vmerr := evm.Create(sender, data, gas, value)
	if vmerr != nil {
		log.Error("Error Soma Governance Contract deployment")
		return contractAddress, vmerr
	}

	log.Info("Deployed Soma Governance Contract", "Address", contractAddress.String())

	return contractAddress, nil
}

func (sb *Backend) contractGetValidators(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB) ([]common.Address, error) {
	sender := vm.AccountRef(sb.config.Deployer)
	gas := uint64(0xFFFFFFFF)
	evm := sb.getEVM(chain, header, sb.config.Deployer, statedb)

	somaAbi, err := abi.JSON(strings.NewReader(sb.config.ABI))
	if err != nil {
		return nil, err
	}

	input, err := somaAbi.Pack("getValidators")
	if err != nil {
		return nil, err
	}

	value := new(big.Int).SetUint64(0x00)
	//A standard call is issued - we leave the possibility to modify the state
	ret, _, vmerr := evm.Call(sender, sb.somaContract, input, gas, value)
	if vmerr != nil {
		log.Error("Error Soma Governance Contract GetValidators()")
		return nil, vmerr
	}

	var addresses []common.Address
	if err := somaAbi.Unpack(&addresses, "getValidators", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getValidators returned value")
		return nil, err
	}

	sortableAddresses := common.Addresses(addresses)
	sort.Sort(sortableAddresses)
	return sortableAddresses, nil
}
