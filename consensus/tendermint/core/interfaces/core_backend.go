package interfaces

import (
	"context"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/core/message"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	ethcore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

// Backend is the interface used by Core
type Backend interface {
	Address() common.Address

	AddSeal(block *types.Block) (*types.Block, error)

	AskSync(committee *types.Committee)

	// Broadcast sends a message to all validators (include self)
	Broadcast(committee *types.Committee, message message.Msg)

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock *types.Block, round int64, seals [][]byte) error

	// GetContractABI returns the Autonity Contract ABI
	GetContractABI() *abi.ABI

	// Gossip sends a message to all validators (exclude self)
	Gossip(committee *types.Committee, message message.Msg)

	KnownMsgHash() []common.Hash

	HandleUnhandledMsgs(ctx context.Context)

	// HeadBlock retrieves latest committed proposal and the address of proposer
	HeadBlock() *types.Block

	Post(ev any)

	// SetProposedBlockHash is a setter for the proposed block hash
	SetProposedBlockHash(hash common.Hash)

	// Sign signs input data with the backend's private key
	Sign(hash common.Hash) ([]byte, common.Address)

	Subscribe(types ...any) *event.TypeMuxSubscription

	SyncPeer(address common.Address)

	// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyProposal(*types.Block) (time.Duration, error)

	// Returns the main blockchain object.
	BlockChain() *ethcore.BlockChain

	// SetBlockchain is used to set the blockchain on this object
	SetBlockchain(bc *ethcore.BlockChain)

	// RemoveMessageFromLocalCache removes a local message from the known messages cache.
	// It is called by Core when some unprocessed messages are removed from the untrusted backlog buffer.
	RemoveMessageFromLocalCache(message message.Msg)

	// Logger returns the object used for logging purposes.
	Logger() log.Logger

	// IsJailed returns true if the address belongs to the jailed validator list.
	IsJailed(address common.Address) bool

	// Gossiper returns gossiper object
	Gossiper() Gossiper
}

type Core interface {
	Start(ctx context.Context, contract *autonity.ProtocolContracts)
	Stop()
	CurrentHeightMessages() []message.Msg
	CoreState() CoreState
	Broadcaster() Broadcaster
	Proposer() Proposer
	Prevoter() Prevoter
	Precommiter() Precommiter
}
