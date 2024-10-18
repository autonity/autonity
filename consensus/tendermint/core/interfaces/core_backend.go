package interfaces

import (
	"context"
	"math/big"
	"time"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	ethcore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
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
	Commit(proposalBlock *types.Block, round int64, quorumCertificate types.AggregateSignature) error

	// GetContractABI returns the Autonity Contract ABI
	GetContractABI() *abi.ABI

	// Gossip sends a message to all validators (exclude self)
	Gossip(committee *types.Committee, message message.Msg)

	KnownMsgHash() []common.Hash

	HandleUnhandledMsgs(ctx context.Context)

	// HeadBlock retrieves latest committed proposal and the address of proposer
	HeadBlock() *types.Block

	Post(ev any)

	ProposedBlockHash() common.Hash
	// SetProposedBlockHash is a setter for the proposed block hash
	SetProposedBlockHash(hash common.Hash)

	// Sign signs input data with the backend's private key
	Sign(hash common.Hash) blst.Signature

	Subscribe(types ...any) *event.TypeMuxSubscription

	SyncPeer(address common.Address)

	// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyProposal(*types.Block) (time.Duration, error)

	// Returns the main blockchain object.
	BlockChain() *ethcore.BlockChain

	EpochOfHeight(height uint64) (*types.EpochInfo, error)

	// SetBlockchain is used to set the blockchain on this object
	SetBlockchain(bc *ethcore.BlockChain)

	// Logger returns the object used for logging purposes.
	Logger() log.Logger

	// IsJailed returns true if the address belongs to the jailed validator list.
	IsJailed(address common.Address) bool

	// Gossiper returns gossiper object
	Gossiper() Gossiper

	// re-injects buffered future height messages
	ProcessFutureMsgs(height uint64)

	// returns future height buffered messages. Called by core for tendermint state dump
	FutureMsgs() []message.Msg

	// returns the channel used to pass messages between peer sessions and the aggregator
	MessageCh() <-chan events.UnverifiedMessageEvent

	// ProposalVerified notifies miner a proposal is verified
	ProposalVerified(block *types.Block)

	// IsProposalStateCached checks if the proposal is cached in the blockchain
	IsProposalStateCached(hash common.Hash) bool
}

type Core interface {
	Start(ctx context.Context, contract *autonity.ProtocolContracts)
	Stop()
	CoreState() CoreState
	Broadcaster() Broadcaster
	Proposer() Proposer
	Prevoter() Prevoter
	Precommiter() Precommiter
	Height() *big.Int
	Round() int64
	CurrentHeightMessages() []message.Msg

	// Used by the aggregator
	Power(h uint64, r int64) *message.AggregatedPower
	VotesPower(h uint64, r int64, code uint8) *message.AggregatedPower
	VotesPowerFor(h uint64, r int64, code uint8, v common.Hash) *message.AggregatedPower
}

type EventDispatcher interface {
	Post(ev any)
}
