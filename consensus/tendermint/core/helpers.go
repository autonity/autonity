package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"math/big"
)

// OverQuorumVotes compute voting power out from a set of prevotes or precommits of a certain round and height, the caller
// should make sure that the votes belong to a certain round and height, it returns a set of votes that the corresponding
// voting power is over quorum, otherwise it returns nil.
func OverQuorumVotes(msgs []message.Msg, quorum *big.Int) (overQuorumVotes []message.Msg) {
	votingPower := new(big.Int)
	counted := make(map[common.Address]struct{})
	for _, v := range msgs {
		if _, ok := counted[v.Sender()]; ok {
			continue
		}
		counted[v.Sender()] = struct{}{}
		votingPower = votingPower.Add(votingPower, v.Power())
		overQuorumVotes = append(overQuorumVotes, v)
		if votingPower.Cmp(quorum) >= 0 {
			return overQuorumVotes
		}
	}
	return nil
}
