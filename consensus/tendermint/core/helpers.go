package core

import (
	"math/big"

	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

// OverQuorumVotes compute voting power out from a set of prevotes or precommits of a certain round and height, the caller
// should make sure that the votes belong to a certain round and height, it returns a set of votes that the corresponding
// voting power is over quorum, otherwise it returns nil.
func OverQuorumVotes(msgs []message.Msg, quorum *big.Int) (overQuorumVotes []message.Msg) {
	if message.Power(msgs).Cmp(quorum) >= 0 {
		return msgs
	}
	return nil
}
