// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	time "time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
)

var (
	// errNotFromProposer is returned when received message is supposed to be from
	// proposer.
	errNotFromProposer = errors.New("message does not come from proposer")
	// errOldHeightMessage is returned when the received message's view is earlier
	// than curRoundMessages view.
	errOldHeightMessage = errors.New("old height message")
	// errInvalidMessage is returned when the message is malformed.
	errInvalidMessage = errors.New("invalid message")
	// errInvalidSenderOfCommittedSeal is returned when the committed seal is not from the sender of the message.
	errInvalidSenderOfCommittedSeal = errors.New("invalid sender of committed seal")
	// errFailedDecodeProposal is returned when the PROPOSAL message is malformed.
	errFailedDecodeProposal = errors.New("failed to decode PROPOSAL")
	// errFailedDecodePrevote is returned when the PREVOTE message is malformed.
	errFailedDecodePrevote = errors.New("failed to decode PREVOTE")
	// errFailedDecodePrecommit is returned when the PRECOMMIT message is malformed.
	errFailedDecodePrecommit = errors.New("failed to decode PRECOMMIT")
)

const (
	MaxRound = 99 // consequence of backlog priority
)

func addr(a common.Address) string {
	return hex.EncodeToString(a[:3])
}

// New creates an Tendermint consensus core
func New(backend Backend, config *config.Config, key *ecdsa.PrivateKey) *core {
	addr := backend.Address()
	logger := log.New("addr", addr.String())
	c := &core{
		key:                   key,
		proposerPolicy:        config.ProposerPolicy,
		address:               addr,
		logger:                logger,
		backend:               backend,
		pendingUnminedBlocks:  make(map[uint64]*types.Block),
		pendingUnminedBlockCh: make(chan *types.Block),
		valueSet:              sync.NewCond(&sync.Mutex{}),
		msgCache:              newMessageStore(),
	}
	o := &oracle{
		c:     c,
		store: c.msgCache,
	}
	c.ora = o
	return c
}

type core struct {
	key            *ecdsa.PrivateKey
	proposerPolicy config.ProposerPolicy
	address        common.Address
	logger         log.Logger

	backend Backend
	cancel  context.CancelFunc

	eventsSub               *event.TypeMuxSubscription
	newUnminedBlockEventSub *event.TypeMuxSubscription
	syncEventSub            *event.TypeMuxSubscription
	wg                      *sync.WaitGroup

	msgCache  *messageCache
	syncTimer *time.Timer

	// map[Height]UnminedBlock
	pendingUnminedBlocks     map[uint64]*types.Block
	pendingUnminedBlocksMu   sync.Mutex
	pendingUnminedBlockCh    chan *types.Block
	isWaitingForUnminedBlock bool

	committee  committee
	lastHeader *types.Header

	autonityContract *autonity.Contract

	height *big.Int
	algo   *algorithm.Algorithm
	ora    *oracle

	valueSet     *sync.Cond
	value        *types.Block
	currentBlock *types.Block
}

func (c *core) SetValue(b *types.Block) {
	c.valueSet.L.Lock()
	defer c.valueSet.L.Unlock()
	if c.value == nil {
		c.valueSet.Signal()
	}
	c.value = b
	println(addr(c.address), c.height, "setting value", c.value.Hash().String()[2:8], "value height", c.value.Number().String())
}

func (c *core) AwaitValue(ctx context.Context, height *big.Int) (*types.Block, error) {
	c.valueSet.L.Lock()
	defer c.valueSet.L.Unlock()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if c.value == nil || c.value.Number().Cmp(height) != 0 {
				c.value = nil
				if c.value == nil {
					println(addr(c.address), c.height.String(), "awaiting vlaue", "valueisnil")
				} else {
					println(addr(c.address), c.height.String(), "awaiting vlaue", "value height", c.value.Number().String(), "awaited height", height.String())
				}
				c.valueSet.Wait()
			} else {
				v := c.value
				println(addr(c.address), c.height, "received awaited vlaue", c.value.Hash().String()[2:8], "value height", c.value.Number().String(), "awaited height", height.String())

				// We put the value in the store here since this is called from the main
				// thread of the algorithm, and so we don't end up needing to syncronise
				// the store.  TODO this is a potential memory leak. We are adding a value
				// without it being referenced by a message that is tied to a height, so it
				// may never be cleared.
				c.msgCache.addValue(v.Hash(), v)
				// We assume our own suggestions are valid
				c.msgCache.setValid(v.Hash())
				c.value = nil
				return v, nil
			}
		}
	}
}

func (c *core) Commit(proposal *algorithm.ConsensusMessage) (*types.Block, error) {
	block := c.msgCache.value(common.Hash(proposal.Value))
	committedSeals := c.msgCache.signatures(algorithm.ValueID(block.Hash()), proposal.Round, block.NumberU64())
	// Sanity checks
	if block == nil {
		return nil, fmt.Errorf("attempted to commit nil block")
	}
	if proposal.Round < 0 {
		return nil, fmt.Errorf("attempted to commit a block in a negative round: %d", proposal.Round)
	}
	if len(committedSeals) == 0 {
		return nil, fmt.Errorf("attempted to commit block without any committed seals")
	}

	for _, seal := range committedSeals {
		if len(seal) != types.BFTExtraSeal {
			return nil, fmt.Errorf("attempted to commit block with a committed seal of invalid length: %s", hex.EncodeToString(seal))
		}
	}
	h := block.Header()
	h.CommittedSeals = committedSeals
	h.Round = uint64(proposal.Round)
	block = block.WithSeal(h)
	c.backend.Commit(block, c.committee.GetProposer(proposal.Round).Address)

	c.logger.Info("commit a block", "hash", block.Hash())
	return block, nil
}

// Metric collecton of round change and height change.
func (c *core) measureHeightRoundMetrics(round int64) {
	if round == 0 {
		tendermintHeightChangeMeter.Mark(1)
	}
	tendermintRoundChangeMeter.Mark(1)
}

func (c *core) createCommittee(block *types.Block) committee {
	var committeeSet committee
	var err error
	var lastProposer common.Address
	header := block.Header()
	switch c.proposerPolicy {
	case config.RoundRobin:
		if !header.IsGenesis() {
			lastProposer, err = types.Ecrecover(header)
			if err != nil {
				panic(fmt.Sprintf("unable to recover proposer address from header %q: %v", header, err))
			}
		}
		committeeSet, err = newRoundRobinSet(header.Committee, lastProposer)
		if err != nil {
			panic(fmt.Sprintf("failed to construct committee %v", err))
		}
	case config.WeightedRandomSampling:
		committeeSet = newWeightedRandomSamplingCommittee(block, c.autonityContract, c.backend.BlockChain())
	default:
		panic(fmt.Sprintf("unrecognised proposer policy %q", c.proposerPolicy))
	}
	return committeeSet
}
