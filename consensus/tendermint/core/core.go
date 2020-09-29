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
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"

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

// New creates an Tendermint consensus core
func New(backend Backend, config *config.Config) *core {
	addr := backend.Address()
	logger := log.New("addr", addr.String())
	return &core{
		proposerPolicy:        config.ProposerPolicy,
		address:               addr,
		logger:                logger,
		backend:               backend,
		pendingUnminedBlocks:  make(map[uint64]*types.Block),
		pendingUnminedBlockCh: make(chan *types.Block),
		stopped:               make(chan struct{}, 4),
		committee:             nil,
	}
}

type core struct {
	proposerPolicy config.ProposerPolicy
	address        common.Address
	logger         log.Logger

	backend Backend
	cancel  context.CancelFunc

	eventsSub               *event.TypeMuxSubscription
	newUnminedBlockEventSub *event.TypeMuxSubscription
	syncEventSub            *event.TypeMuxSubscription
	stopped                 chan struct{}

	msgCache *messageCache
	// map[Height]UnminedBlock
	pendingUnminedBlocks     map[uint64]*types.Block
	pendingUnminedBlocksMu   sync.Mutex
	pendingUnminedBlockCh    chan *types.Block
	isWaitingForUnminedBlock bool

	committee  committee
	lastHeader *types.Header
	// height, round and committeeSet are the ONLY guarded fields.
	// everything else MUST be accessed only by the main thread.
	stateMu sync.RWMutex

	autonityContract *autonity.Contract

	height *big.Int
	algo   *algorithm.Algorithm
}

func (c *core) GetCurrentHeightMessages() []*Message {
	return c.msgCache.heightMessages(c.Height().Uint64())
}

func (c *core) finalizeMessage(msg *Message) ([]byte, error) {
	var err error

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) broadcast(ctx context.Context, m *algorithm.ConsensusMessage) {
	logger := c.logger.New("step", nil)

	var code uint64
	var internalMessage interface{}
	bigHeight := new(big.Int).SetUint64(m.Height)
	switch m.MsgType {
	case algorithm.Propose:
		code = msgProposal
		internalMessage = NewProposal(m.Round, bigHeight, m.ValidRound, c.msgCache.value(common.Hash(m.Value)))
	case algorithm.Prevote:
		code = msgPrevote
		internalMessage = &Vote{
			Round:             m.Round,
			Height:            bigHeight,
			ProposedBlockHash: common.Hash(m.Value),
		}
	case algorithm.Precommit:
		code = msgPrecommit
		internalMessage = &Vote{
			Round:             m.Round,
			Height:            bigHeight,
			ProposedBlockHash: common.Hash(m.Value),
		}
	}
	marshalledInternalMessage, err := Encode(internalMessage)
	if err != nil {
		panic(fmt.Sprintf("error while encoding consensus message: %v", err))
	}
	msg := &Message{
		Code:          code,
		Address:       c.address,
		CommittedSeal: []byte{}, // Not sure why this is empty but it seems to be set like this everywhere.
		Msg:           marshalledInternalMessage,
	}

	if m.MsgType == algorithm.Precommit {
		seal := PrepareCommittedSeal(common.Hash(m.Value), m.Round, bigHeight)
		msg.CommittedSeal, err = c.backend.Sign(seal)
		if err != nil {
			panic(fmt.Sprintf("error while signing committed seal: %v", err))
		}
	}

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		panic(fmt.Sprintf("Failed to finalize message: %+v err: %v", msg, err))
	}

	// Broadcast payload
	logger.Debug("broadcasting", "msg", msg.String())
	if err = c.backend.Broadcast(ctx, c.committeeSet().Committee(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

// check if msg sender is proposer for proposal handling.
func (c *core) isProposerMsg(round int64, msgAddress common.Address) bool {
	return c.committeeSet().GetProposer(round).Address == msgAddress
}

func (c *core) Commit(proposal *algorithm.ConsensusMessage) error {
	block := c.msgCache.value(common.Hash(proposal.Value))
	committedSeals := c.msgCache.signatures(block.Hash(), proposal.Round, block.NumberU64())
	// Sanity checks
	if block == nil {
		return fmt.Errorf("attempted to commit nil block")
	}
	if proposal.Round < 0 {
		return fmt.Errorf("Attempted to commit a block in a negative round: %d", proposal.Round)
	}
	if len(committedSeals) == 0 {
		return fmt.Errorf("attempted to commit block without any committed seals")
	}

	for _, seal := range committedSeals {
		if len(seal) != types.BFTExtraSeal {
			return fmt.Errorf("Attempted to commit block with a committed seal of invalid length: %s", hex.EncodeToString(seal))
		}
	}
	h := block.Header()
	h.CommittedSeals = committedSeals
	h.Round = uint64(proposal.Round)
	block = block.WithSeal(h)
	if err := c.backend.Commit(block, proposal.Round, committedSeals); err != nil {
		c.logger.Error("failed to commit a block", "err", err)
	}

	c.logger.Info("commit a block", "hash", block.Hash())
	return nil
}

func (c *core) commit(block *types.Block, round int64) {
}

// Metric collecton of round change and height change.
func (c *core) measureHeightRoundMetrics(round int64) {
	if round == 0 {
		tendermintHeightChangeMeter.Mark(1)
	}
	tendermintRoundChangeMeter.Mark(1)
}

func (c *core) updateLatestBlock() {
	lastBlockMined, _ := c.backend.LastCommittedProposal()
	c.setHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))

	lastHeader := lastBlockMined.Header()
	var committeeSet committee
	var err error
	var lastProposer common.Address
	switch c.proposerPolicy {
	case config.RoundRobin:
		if !lastHeader.IsGenesis() {
			lastProposer, err = types.Ecrecover(lastHeader)
			if err != nil {
				panic(fmt.Sprintf("unable to recover proposer address from header %q: %v", lastHeader, err))
			}
		}
		committeeSet, err = newRoundRobinSet(lastHeader.Committee, lastProposer)
		if err != nil {
			panic(fmt.Sprintf("failed to construct committee %v", err))
		}
	case config.WeightedRandomSampling:
		committeeSet = newWeightedRandomSamplingCommittee(lastBlockMined, c.autonityContract, c.backend.BlockChain())
	default:
		panic(fmt.Sprintf("unrecognised proposer policy %q", c.proposerPolicy))
	}

	c.lastHeader = lastHeader
	c.setCommitteeSet(committeeSet)

}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash, round int64, height *big.Int) []byte {
	var buf bytes.Buffer
	roundBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(roundBytes, uint64(round))
	buf.Write(roundBytes)
	buf.Write(height.Bytes())
	buf.Write(hash.Bytes())
	return buf.Bytes()
}

func (c *core) setHeight(height *big.Int) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.height = height
}
func (c *core) setCommitteeSet(set committee) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.committee = set
}

func (c *core) Height() *big.Int {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.height
}
func (c *core) committeeSet() committee {
	c.stateMu.RLock()
	defer c.stateMu.RUnlock()
	return c.committee
}

func (c *core) verifyCommittedSeal(addressMsg common.Address, committedSealMsg []byte, proposedBlockHash common.Hash, round int64, height *big.Int) error {
	committedSeal := PrepareCommittedSeal(proposedBlockHash, round, height)

	sealerAddress, err := types.GetSignatureAddress(committedSeal, committedSealMsg)
	if err != nil {
		c.logger.Error("Failed to get signer address", "err", err)
		return err
	}

	// ensure sender signed the committed seal
	if !bytes.Equal(sealerAddress.Bytes(), addressMsg.Bytes()) {
		c.logger.Error("verify precommit seal error", "got", addressMsg.String(), "expected", sealerAddress.String())

		return errInvalidSenderOfCommittedSeal
	}

	return nil
}
