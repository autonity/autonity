package core

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
	Address() common.Address

	AddSeal(block *types.Block) (*types.Block, error)

	AskSync(header *types.Header)

	// Broadcast sends a message to all validators (include self)
	Broadcast(ctx context.Context, committee types.Committee, payload []byte) error

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock *types.Block, proposer common.Address)

	GetContractABI() string

	// Gossip sends a message to all validators (exclude self)
	Gossip(ctx context.Context, committee types.Committee, payload []byte)

	HandleUnhandledMsgs(ctx context.Context)

	// LastCommittedProposal retrieves latest committed proposal and the address of proposer
	LastCommittedProposal() (*types.Block, common.Address)

	Post(ev interface{})

	// Setter for proposed block hash
	SetProposedBlockHash(hash common.Hash)

	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)

	Subscribe(types ...interface{}) *event.TypeMuxSubscription

	SyncPeer(address common.Address, messages [][]byte)

	// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyProposal(types.Block) (time.Duration, error)

	WhiteList() []string

	BlockChain() *ethcore.BlockChain

	//Used to set the blockchain on this
	SetBlockchain(bc *ethcore.BlockChain)
}

type Tendermint interface {
	Start(ctx context.Context, contract *autonity.Contract)
	Stop()
}
