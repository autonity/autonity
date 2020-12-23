package tendermint

import (
	"context"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
)

// Backend provides application specific functions for Istanbul core
type Backend interface {

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock *types.Block, proposer common.Address)
}

type Tendermint interface {
	consensus.Handler
	Start(ctx context.Context, contract *autonity.Contract, blockchain *core.BlockChain)
	SetValue(*types.Block)
	Stop()
}
