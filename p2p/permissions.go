package p2p

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// attachDb creates an instance of the statedb returning the latest copy with the latest block header
func attachDb(datadir string) (*state.StateDB, *types.Header, error) {
	db, err := ethdb.NewLDBDatabase(datadir+"/geth/chaindata", 768, 512)
	if err != nil {
		return nil, nil, err
	}

	hash := rawdb.ReadHeadBlockHash(db)
	number := rawdb.ReadHeaderNumber(db, hash)
	header := rawdb.ReadHeader(db, hash, *number)

	sdb := state.NewDatabase(db)
	statedb, err := state.New(header.Root, sdb)
	if err != nil {
		return nil, nil, err
	}
	return statedb, header, nil
}

// callGlienick queries the Soma contract, for any functions that take and address as argument.
// Returns true/false if the the address is an active validator and false if not.
func callGlienicke(functionSig string, node string, contractAddress common.Address, header *types.Header, statedb *state.StateDB) (bool, error) {
	sender := vm.AccountRef(contractAddress)
	gas := uint64(0xFFFFFFFF)
	evm := getEVM(header, contractAddress, contractAddress, statedb)

	// Pad address for ABI encoding
	input, err := packInputData(node)
	if err != nil {
		return false, err
	}

	// Call ActiveValidators()
	ret, gas, vmerr := evm.StaticCall(sender, contractAddress, input, gas)
	if len(ret) == 0 {
		log.Info("contractCallAddress(): No return value", "Block", header.Number.Int64())
		return false, nil
	}
	if vmerr != nil {
		return false, vmerr
	}

	const def = `[{ "name" : "method", "outputs": [{ "type": "bool" }] }]`
	funcAbi, err := abi.JSON(strings.NewReader(def))
	if err != nil {
		return false, err
	}

	var output bool
	err = funcAbi.Unpack(&output, "method", ret)
	if err != nil {
		return false, err
	}

	return output, nil
}

// packInputData
func packInputData(input string) ([]byte, error) {
	const def = `[
	{ "type" : "function", "name" : "IsAllowed", "constant" : true, "inputs" : [ { "name" : "str", "type" : "string" } ] }
	]`

	funcAbi, err := abi.JSON(strings.NewReader(def))
	if err != nil {
		log.Error("Error Creating ABI JSON", "Error", err)
		return nil, err
	}

	data, err := funcAbi.Pack("IsAllowed", input)
	if err != nil {
		log.Error("Error Packing Input Data", "Error", err)
		return nil, err
	}

	return data, nil
}

// getEVM Instantiates a new EVM object which is required when creating or calling a deployed contract
func getEVM(header *types.Header, coinbase, origin common.Address, statedb *state.StateDB) *vm.EVM {
	gasPrice := new(big.Int).SetUint64(0x0)
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(n uint64) common.Hash { return header.Root }, // since this a one time thing, no point in adding complex functions to get the hash
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    gasPrice,
	}

	chainConfig := params.AllCliqueProtocolChanges
	vmconfig := vm.Config{}
	evm := vm.NewEVM(evmContext, statedb, chainConfig, vmconfig)
	return evm
}
