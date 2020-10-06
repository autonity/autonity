package core

import (
	"context"
	"math/big"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
)

type oracle struct {
	lastHeader   *types.Header
	store        *messageCache
	committeeSet committee
	c            *core
	ctx          context.Context
}

func (o *oracle) FThresh(round int64) bool {
	return o.store.fail(round, o.lastHeader)
}

func (o *oracle) MatchingProposal(cm *algorithm.ConsensusMessage) *algorithm.ConsensusMessage {
	return o.store.matchingProposal(cm)
}

func (o *oracle) PrecommitQThresh(round int64, value *algorithm.ValueID) bool {
	return o.store.precommitQuorum((*common.Hash)(value), round, o.lastHeader)
}

func (o *oracle) PrevoteQThresh(round int64, value *algorithm.ValueID) bool {
	return o.store.prevoteQuorum((*common.Hash)(value), round, o.lastHeader)
}

func (o *oracle) Proposer(round int64, nodeID algorithm.NodeID) bool {
	return o.committeeSet.GetProposer(round).Address == common.Address(nodeID)
}

func (o *oracle) Valid(value algorithm.ValueID) bool {
	return o.store.isValid(common.Hash(value))
}

func (o *oracle) Value(height uint64) algorithm.ValueID {

	vchan := make(chan *types.Block)
	go func() {
		vchan <- o.c.AwaitValue(new(big.Int).SetUint64(height))
	}()
	select {
	case b := <-vchan:
		return algorithm.ValueID(b.Hash())
	case <-o.ctx.Done():
		return algorithm.ValueID{}
	}
	// return algorithm.ValueID(o.c.AwaitValue(new(big.Int).SetUint64(height)).Hash())
}
