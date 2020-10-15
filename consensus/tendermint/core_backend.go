package tendermint

import (
	"context"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/contracts/autonity"
	ethcore "github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
)

// Backend provides application specific functions for Istanbul core
type Backend interface {

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock *types.Block, proposer common.Address)

	HandleUnhandledMsgs(ctx context.Context)

	Post(ev interface{})

	Subscribe(types ...interface{}) *event.TypeMuxSubscription

	// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyProposal(types.Block) (time.Duration, error)

	//Used to set the blockchain on this
	SetBlockchain(bc *ethcore.BlockChain)
}

type Tendermint interface {
	Start(ctx context.Context, contract *autonity.Contract)
	SetValue(*types.Block)
	Stop()
}
