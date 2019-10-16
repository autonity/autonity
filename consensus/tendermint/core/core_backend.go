package core

import (
	"context"
	"math/big"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/rpc"
)

func (c *core) Author(header *types.Header) (common.Address, error) {
	return c.backend.Author(header)
}

func (c *core) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return c.backend.VerifyHeader(chain, header, seal)
}

func (c *core) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return c.backend.VerifyHeaders(chain, headers, seals)
}

func (c *core) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return c.backend.VerifyUncles(chain, block)
}

func (c *core) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return c.backend.VerifySeal(chain, header)
}

func (c *core) Prepare(chain consensus.ChainReader, header *types.Header) error {
	return c.backend.Prepare(chain, header)
}

func (c *core) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	c.backend.Finalize(chain, header, state, txs, uncles)
}

func (c *core) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	return c.backend.FinalizeAndAssemble(chain, header, state, txs, uncles, receipts)
}

func (c *core) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return c.backend.Seal(chain, block, results, stop)
}

func (c *core) SealHash(header *types.Header) common.Hash {
	return c.backend.SealHash(header)
}

func (c *core) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return c.backend.CalcDifficulty(chain, time, parent)
}

func (c *core) APIs(chain consensus.ChainReader) []rpc.API {
	return c.backend.APIs(chain)
}

func (c *core) Close() error {
	return c.Stop()
}

func (c *core) NewChainHead() error {
	return c.backend.NewChainHead()
}

func (c *core) HandleMsg(address common.Address, data p2p.Msg) (bool, error) {
	return c.backend.HandleMsg(address, data)
}

func (c *core) SetBroadcaster(b consensus.Broadcaster) {
	c.backend.SetBroadcaster(b)
}

func (c *core) Protocol() (protocolName string, extraMsgCodes uint64) {
	return c.backend.Protocol()
}

// Synchronize new connected peer with current height state
func (c *core) SyncPeer(address common.Address) {
	if c.IsValidator(address) {
		c.backend.SyncPeer(address, c.GetCurrentHeightMessages())
	}
}

func (c *core) ResetPeerCache(address common.Address) {
	c.backend.ResetPeerCache(address)
}

// Backend provides application specific functions for Istanbul core
type Backend interface {
	consensus.Engine
	consensus.Handler
	Start(ctx context.Context, chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error

	// Address returns the owner's address
	Address() common.Address

	// Validators returns the validator set
	Validators(number uint64) validator.Set

	Subscribe(types ...interface{}) *event.TypeMuxSubscription

	Post(ev interface{})

	// Broadcast sends a message to all validators (include self)
	Broadcast(ctx context.Context, valSet validator.Set, payload []byte) error

	// Gossip sends a message to all validators (exclude self)
	Gossip(ctx context.Context, valSet validator.Set, payload []byte)

	// Commit delivers an approved proposal to backend.
	// The delivered proposal will be put into blockchain.
	Commit(proposalBlock types.Block, seals [][]byte) error

	// VerifyProposal verifies the proposal. If a consensus.ErrFutureBlock error is returned,
	// the time difference of the proposal and current time is also returned.
	VerifyProposal(types.Block) (time.Duration, error)

	// Sign signs input data with the backend's private key
	Sign([]byte) ([]byte, error)

	// CheckSignature verifies the signature by checking if it's signed by
	// the given validator
	CheckSignature(data []byte, addr common.Address, sig []byte) error

	// LastCommittedProposal retrieves latest committed proposal and the address of proposer
	LastCommittedProposal() (*types.Block, common.Address)

	// GetProposer returns the proposer of the given block height
	GetProposer(number uint64) common.Address

	// HasBadBlock returns whether the block with the hash is a bad block
	HasBadProposal(hash common.Hash) bool

	// Setter for proposed block hash
	SetProposedBlockHash(hash common.Hash)

	SyncPeer(address common.Address, messages []*Message)

	ResetPeerCache(address common.Address)

	AskSync(set validator.Set)

	GetContractAddress() common.Address

	GetContractABI() string

	WhiteList() []string
}
