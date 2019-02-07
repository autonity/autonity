package core

import (
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/p2p/enode"
)

func (bc *BlockChain) updateEnodesWhitelist(state *state.StateDB, block *types.Block) {
	var newWhitelist []*enode.Node
	if block.Number().Uint64() == 1 {
		// deploy the contract
	}else{
		// call retrieveWhitelist contract function
	}
	rawdb.WriteEnodeWhitelist(bc.db, newWhitelist)
}
