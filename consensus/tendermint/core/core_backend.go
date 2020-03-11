package core

import (
	"context"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/committee"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"time"
)

// Backend provides application specific functions for Istanbul core
type Backend interface {
	Address() common.Address

	AddSeal(block *types.Block) (*types.Block, error)

	AskSync(set committee.Set)

	// Broadcast sends a message to all validators (include self)
	Broadcast(ctx context.Context, valSet committee.Set, payload []byte) error

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock *types.Block, round int64, seals [][]byte) error

	// Validators returns the committee set
	Committee(number uint64) (committee.Set, error)

	// TODO: change the location later once mock_engine_modify_validator_list.go is better understood
	FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Gossip sends a message to all validators (exclude self)
	Gossip(ctx context.Context, valSet committee.Set, payload []byte)

	HandleUnhandledMsgs(ctx context.Context)

	// LastCommittedProposal retrieves latest committed proposal and the address of proposer
	LastCommittedProposal() (*types.Block, common.Address)

	Post(ev interface{})

	// Setter for proposed block hash
	SetProposedBlockHash(hash common.Hash)

	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)

	Subscribe(types ...interface{}) *event.TypeMuxSubscription

	SyncPeer(address common.Address)

	// TODO: change the location later once mock_backend_verfiy_self_proposal_always_true.go is better understood
	VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error

	// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyProposal(types.Block) (time.Duration, error)
}

type CoreEngine interface {
	Start(ctx context.Context) error
	Stop() error
	GetCurrentHeightMessages() []*Message
}
