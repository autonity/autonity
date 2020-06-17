package core

import (
	"context"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
)

// Backend provides application specific functions for Istanbul core
type Backend interface {

	// Backend functions
	Address() common.Address
	// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyProposal(types.Block) (time.Duration, error)
	// Setter for proposed block hash
	SetProposedBlockHash(hash common.Hash)
	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock *types.Block, round int64, seals [][]byte) error
	// LastCommittedProposal retrieves latest committed proposal and the address of proposer
	LastCommittedProposal() (*types.Block, common.Address)
	Post(ev interface{})
	Subscribe(types ...interface{}) *event.TypeMuxSubscription

	// Broadcaster functions

	// Broadcast sends a message to all validators (include self)
	Broadcast(ctx context.Context, valSet *committee.Set, payload []byte) error
	// Gossip sends a message to all validators (exclude self)
	Gossip(ctx context.Context, valSet *committee.Set, payload []byte)

	// Tendermint crypto functions
	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)

	// Validators returns the committee set
	Committee(number uint64) (*committee.Set, error)

	HandleUnhandledMsgs(ctx context.Context)

	// Syncer interface, syncer will make use of broadcaster
	AskSync(set *committee.Set)
	SyncPeer(address common.Address)
}

type Tendermint interface {
	Start(ctx context.Context)
	Stop()
	GetCurrentHeightMessages() []*Message
}
