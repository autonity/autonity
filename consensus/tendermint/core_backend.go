package tendermint

import (
	"context"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
)

// Backend provides application specific functions for Istanbul core
type Backend interface {

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock *types.Block, proposer common.Address)

	Post(ev interface{})

	Subscribe(types ...interface{}) *event.TypeMuxSubscription
}

type Tendermint interface {
	Start(ctx context.Context, contract *autonity.Contract)
	SetValue(*types.Block)
	Stop()
}
