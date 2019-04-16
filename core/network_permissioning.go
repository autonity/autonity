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
	"github.com/pkg/errors"
	"math/big"
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

var GlienickeContractError = errors.New("could not call Glienicke contract")

func (bc *BlockChain) updateEnodesWhitelist(state *state.StateDB, block *types.Block) error {
	newWhitelist, err := bc.GetWhitelist(block, state)
	if err != nil {
		log.Error("could not call Glienicke contract", "err", err)
		return GlienickeContractError
	}

	rawdb.WriteEnodeWhitelist(bc.db, newWhitelist)
	go bc.glienickeFeed.Send(GlienickeEvent{Whitelist: newWhitelist.List})

	return nil
}

func (bc *BlockChain) GetWhitelist(block *types.Block, db *state.StateDB) (*types.Nodes, error) {
	var (
		newWhitelist *types.Nodes
		err error
	)

	if block.Number().Uint64() == 1 {
		// use genesis block whitelist
		newWhitelist = rawdb.ReadEnodeWhitelist(bc.db, false)
	} else {
		// call retrieveWhitelist contract function
		newWhitelist, err = bc.callGlienickeContract(db, block.Header())
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
func (bc *BlockChain) DeployGlienickeContract(state *state.StateDB, header *types.Header) (*types.Nodes, common.Address, error) {
	//if bytecode or abi is missing use default one
	glienickeByteCode := bc.chainConfig.GlienickeBytecode
	glienickeABI := bc.chainConfig.GlienickeABI
	if bc.chainConfig.GlienickeBytecode == "" || bc.chainConfig.GlienickeABI == "" {
		glienickeByteCode = params.GlienickeDefaultBytecode
		glienickeABI = params.GlienickeDefaultABI
	}
	bc.chainConfig.GlienickeABI = glienickeABI

	//Same for deployer
	glienickeDeployer := bc.chainConfig.GlienickeDeployer
	if glienickeDeployer == (common.Address{}) {
		glienickeDeployer = params.GlienickeDefaultDeployer
	}

	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(glienickeByteCode)
	evm := bc.getEVM(header, glienickeDeployer, state)
	sender := vm.AccountRef(glienickeDeployer)

	enodesWhitelist := rawdb.ReadEnodeWhitelist(bc.db, false)

	glienickeAbi, err := abi.JSON(strings.NewReader(glienickeABI))
	if err != nil {
		return nil, common.Address{}, err
	}

	constructorParams, err := glienickeAbi.Pack("", enodesWhitelist.StrList)
	if err != nil {
		return nil, common.Address{}, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Glienicke validator governance contract
	_, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)
	if vmerr != nil {
		log.Error("Error Glienicke Contract deployment")
		return nil, common.Address{}, vmerr
	}

	log.Info("Deployed Glienicke Contract", "Address", contractAddress.String())

	return enodesWhitelist, contractAddress, nil
}

func (bc *BlockChain) callGlienickeContract(state *state.StateDB, header *types.Header) (*types.Nodes, error) {
	// Needs to be refactored somehow
	glienickeDeployer := bc.chainConfig.GlienickeDeployer
	if glienickeDeployer == (common.Address{}) {
		glienickeDeployer = params.GlienickeDefaultDeployer
	}

	glienickeABI := bc.chainConfig.GlienickeABI
	if bc.chainConfig.GlienickeABI == "" {
		glienickeABI = params.GlienickeDefaultABI
	}

	sender := vm.AccountRef(glienickeDeployer)
	gas := uint64(0xFFFFFFFF)
	evm := bc.getEVM(header, glienickeDeployer, state)

	ABI, err := abi.JSON(strings.NewReader(glienickeABI))
	if err != nil {
		return nil, err
	}

	input, err := ABI.Pack("getWhitelist")
	if err != nil {
		return nil, err
	}

	glienickeAddress := crypto.CreateAddress(glienickeDeployer, 0)

	ret, gas, vmerr := evm.StaticCall(sender, glienickeAddress, input, gas)
	if vmerr != nil {
		log.Error("Error Glienicke Contract getWhitelist()")
		return nil, vmerr
	}

	var returnedEnodes []string
	if err := ABI.Unpack(returnedEnodes, "getWhitelist", ret); err != nil { // can't work with aliased types
		log.Error("Could not unpack getWhitelist returned value")
		return nil, err
	}

	return types.NewNodes(returnedEnodes, false), nil
}
