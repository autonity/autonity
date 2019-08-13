// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"github.com/clearmatics/autonity/common/math"
	"math/big"
	"sort"
	"strings"

	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
)

var AutonityContractError = errors.New("could not call Autonity contract")

func (bc *BlockChain) updateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	newWhitelist, err := bc.GetWhitelist(block, state)
	if err != nil {
		log.Error("could not call Autonity contract", "err", err)
		return AutonityContractError
	}

	rawdb.WriteEnodeWhitelist(bc.db, newWhitelist)
	go bc.autonityContractFeed.Send(whitelistEvent{Whitelist: newWhitelist.List})

	return nil
}

func (bc *BlockChain) GetWhitelist(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	var (
		newWhitelist *types.Nodes
		err          error
	)

	if block.Number().Uint64() == 1 {
		// use genesis block whitelist
		newWhitelist = rawdb.ReadEnodeWhitelist(bc.db, false)
	} else {
		// call retrieveWhitelist contract function
		newWhitelist, err = bc.ACgetWhitelist(db, block.Header())
	}

	return newWhitelist, err
}

// Instantiates a new EVM object which is required when creating or calling a deployed contract
func (bc *BlockChain) getEVM(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {

	coinbase, _ := bc.engine.Author(header)
	evmContext := vm.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, bc),
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    new(big.Int).SetUint64(0x0),
	}
	evm := vm.NewEVM(evmContext, statedb, bc.Config(), bc.vmConfig)
	return evm
}

// deployContract deploys the contract contained within the genesis field bytecode
func (bc *BlockChain) DeployAutonityContract(state *state.StateDB, header *types.Header) (*types.Nodes, common.Address, error) {
	//if bytecode or abi is missing use default one
	autonityConfig := bc.chainConfig.AutonityContractConfig
	autonityByteCode := autonityConfig.Bytecode
	autonityABI := autonityConfig.ABI
	if autonityConfig.Bytecode == "" || autonityConfig.ABI == "" {
		autonityByteCode = params.AutonityDefaultBytecode
		autonityABI = params.AutonityDefaultABI
	}
	//Same for deployer
	autonityDeployer := autonityConfig.Deployer
	if autonityDeployer == (common.Address{}) {
		autonityDeployer = params.AutonityDefaultDeployer
	}

	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(autonityByteCode)
	evm := bc.getEVM(header, autonityDeployer, state)
	sender := vm.AccountRef(autonityDeployer)

	participantAddress := make([]common.Address, len(autonityConfig.Users))
	participantEnode := make([]string, len(autonityConfig.Users))
	participantType := make([]uint64, len(autonityConfig.Users))
	participantStake := make([]uint64, len(autonityConfig.Users))

	for i, user := range autonityConfig.Users {
		participantAddress[i] = user.Address
		participantEnode[i] = user.Enode
		participantType[i] = user.Type.Uint()
		participantStake[i] = user.Stake
	}

	autonityAbi, err := abi.JSON(strings.NewReader(autonityABI))
	if err != nil {
		return nil, common.Address{}, err
	}

	constructorParams, err := autonityAbi.Pack("",
		participantAddress,
		participantEnode,
		participantType,
		participantStake,
		autonityConfig.Operator,
		autonityConfig.MinGasPrice)

	if err != nil {
		return nil, common.Address{}, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(math.MaxUint64)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Autonity contract
	_, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)
	if vmerr != nil {
		log.Error("Error Autonity Contract deployment")
		return nil, common.Address{}, vmerr
	}

	log.Info("Deployed Autonity Contract", "Address", contractAddress.String())

	enodesWhitelist := rawdb.ReadEnodeWhitelist(bc.db, false)

	return enodesWhitelist, contractAddress, nil
}

func (bc *BlockChain) AutonityContractCall(state *state.StateDB, header *types.Header, function string, result interface{}) error {
	autonityConfig := bc.chainConfig.AutonityContractConfig
	autonityDeployer := autonityConfig.Deployer
	if autonityDeployer == (common.Address{}) {
		autonityDeployer = params.AutonityDefaultDeployer
	}

	autonityABI := autonityConfig.ABI
	if autonityABI == "" {
		autonityABI = params.AutonityDefaultABI
	}

	sender := vm.AccountRef(autonityDeployer)
	gas := uint64(math.MaxUint64)
	evm := bc.getEVM(header, autonityDeployer, state)

	ABI, err := abi.JSON(strings.NewReader(autonityABI))
	if err != nil {
		return err
	}
	autonityAddress := crypto.CreateAddress(autonityDeployer, 0)

	input, err := ABI.Pack(function)
	if err != nil {
		return err
	}

	ret, gas, vmerr := evm.StaticCall(sender, autonityAddress, input, gas)
	if vmerr != nil {
		log.Error("Error Autonity Contract", "function", function)
		return vmerr
	}

	if err := ABI.Unpack(result, function, ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack returned value", "function", function)
		return  err
	}

	return nil
}

func (bc *BlockChain) ACgetWhitelist(state *state.StateDB, header *types.Header) (*types.Nodes, error) {

	var returnedEnodes []string
	err := bc.AutonityContractCall(state, header, "GetWhitelist", &returnedEnodes)
	if err != nil {
		return nil, err
	}
	return types.NewNodes(returnedEnodes, false), nil
}

func (bc *BlockChain) ACgetValidators(statedb *state.StateDB, header *types.Header) ([]common.Address, error) {

	var addresses []common.Address
	err := bc.AutonityContractCall(statedb, header, "GetValidators", &addresses)
	if err != nil {
		return nil, err
	}
	sortableAddresses := common.Addresses(addresses)
	sort.Sort(sortableAddresses)
	return sortableAddresses, nil
}
