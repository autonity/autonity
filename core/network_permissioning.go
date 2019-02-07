package core

import (
	"github.com/clearmatics/autonity/accounts/abi"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
	"math/big"
	"sort"
	"strings"
)

func (bc *BlockChain) updateEnodesWhitelist(state *state.StateDB, block *types.Block) {
	var newWhitelist []*enode.Node
	var err error
	if block.Number().Uint64() == 1 {
		// deploy the contract
		newWhitelist, err = bc.deployContract(state, block.Header())
		if err != nil {
			log.Crit("Could not deploy Glieniecke contract.", "error", err.Error())
		}
	}else{
		// call retrieveWhitelist contract function
		newWhitelist, err = bc.callGlienieckeContract(state, block.Header())
		if err != nil {
			log.Crit("Could not call Glienecke contract.", "error", err.Error())
		}
	}

	rawdb.WriteEnodeWhitelist(bc.db, newWhitelist)

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

// todo : think about hash of genesis and verify on state after deployment (inc. Soma)
// deployContract deploys the contract contained within the genesis field bytecode
func (bc *BlockChain) deployContract(state *state.StateDB, header *types.Header) ([]*enode.Node, error) {
	// Convert the contract bytecode from hex into bytes
	contractBytecode := common.Hex2Bytes(bc.chainConfig.GlienickeBytecode)
	evm := bc.getEVM(header, bc.config.Deployer, state)
	sender := vm.AccountRef(bc.config.Deployer)

	currentWhitelist :=  rawdb.ReadEnodeWhitelist(bc.db)

	GlienickeAbi, err := abi.JSON(strings.NewReader(bc.chainConfig.GlienickeABI))
	if err != nil {
		return nil, err
	}
	constructorParams, err := GlienickeAbi.Pack("", currentWhitelist)
	if err != nil {
		return nil, err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Glieniecke validator governance contract
	_, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)

	if vmerr != nil {
		log.Error("Error Glienicke Contract deployment")
		return nil, vmerr
	}

	log.Info("Deployed Glienicke Contract", "Address", contractAddress.String())

	return currentWhitelist, nil
}

