package tendermint

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	"github.com/clearmatics/autonity/core/types"
)

type oracle struct {
	lastHeader   *types.Header
	store        *messageStore
	committeeSet committee
}

func newOracle(lh *types.Header, s *messageStore, cs committee) *oracle {
	return &oracle{
		lastHeader:   lh,
		store:        s,
		committeeSet: cs,
	}
}

func (o *oracle) FThresh(round int64) bool {
	return o.store.fail(round, o.lastHeader)
}

func (o *oracle) MatchingProposal(cm *algorithm.ConsensusMessage) *algorithm.ConsensusMessage {
	return o.store.matchingProposal(cm).consensusMessage
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

func (o *oracle) Height() uint64 {
	return o.lastHeader.Number.Uint64() + 1
}

func (o *oracle) Value() algorithm.ValueID {
	//return o.c.AwaitValue()
	return [32]byte{}
}
