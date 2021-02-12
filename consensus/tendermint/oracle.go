package tendermint

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
)

type oracle struct {
	height     uint64
	lastHeader *types.Header
	store      *messageStore
	ba         *blockAwaiter
}

func newOracle(lh *types.Header, s *messageStore, ba *blockAwaiter) *oracle {
	return &oracle{
		height:     lh.Number.Uint64() + 1,
		lastHeader: lh,
		store:      s,
		ba:         ba,
	}
}

func (o *oracle) FThresh(round int64) bool {
	return o.store.fail(round, o.lastHeader)
}

func (o *oracle) MatchingProposal(cm *algorithm.ConsensusMessage) *algorithm.ConsensusMessage {
	mp := o.store.matchingProposal(cm)
	if mp != nil {
		return mp.consensusMessage
	}
	return nil
}

func (o *oracle) PrecommitQThresh(round int64, value *algorithm.ValueID) bool {
	return o.store.precommitQuorum((*common.Hash)(value), round, o.lastHeader)
}

func (o *oracle) PrevoteQThresh(round int64, value *algorithm.ValueID) bool {
	return o.store.prevoteQuorum((*common.Hash)(value), round, o.lastHeader)
}

func (o *oracle) Valid(value algorithm.ValueID) bool {
	return o.store.isValid(common.Hash(value))
}

func (o *oracle) Height() uint64 {
	return o.height
}
