package interfaces

import (
	"context"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"time"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	ethcore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
)

// Backend provides application specific functions for Istanbul Core
type Backend interface {
	Address() common.Address

	AddSeal(block *types.Block) (*types.Block, error)

	AskSync(header *types.Header)

	// Broadcast sends a message to all validators (include self)
	Broadcast(ctx context.Context, committee types.Committee, message message.Message) error

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock *types.Block, round int64, seals [][]byte) error

	// GetContractABI returns the Autonity Contract ABI
	GetContractABI() *abi.ABI

	// Gossip sends a message to all validators (exclude self)
	Gossip(ctx context.Context, committee types.Committee, payload []byte)

	KnownMsgHash() []common.Hash

	HandleUnhandledMsgs(ctx context.Context)

	// HeadBlock retrieves latest committed proposal and the address of proposer
	HeadBlock() (*types.Block, common.Address)

	Post(ev any)

	// Setter for proposed block hash
	SetProposedBlockHash(hash common.Hash)

	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)

	Subscribe(types ...any) *event.TypeMuxSubscription

	SyncPeer(address common.Address)

	// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyProposal(*types.Block) (time.Duration, error)

	// Returns the main blockchain object.
	BlockChain() *ethcore.BlockChain

	//Used to set the blockchain on this
	SetBlockchain(bc *ethcore.BlockChain)

	// RemoveMessageFromLocalCache removes a local message from the known messages cache.
	// It is called by Core when some unprocessed messages are removed from the untrusted backlog buffer.
	RemoveMessageFromLocalCache(payload []byte)

	// Logger returns the object used for logging purposes.
	Logger() log.Logger

	// IsJailed returns true if the address belongs to the jailed validator list.
	IsJailed(address common.Address) bool
}
