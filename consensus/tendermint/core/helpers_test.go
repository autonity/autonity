package core

import (
	"github.com/autonity/autonity/common"
	proto "github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	mUtils "github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestOverQuorumVotes(t *testing.T) {
	t.Run("with duplicated votes, it returns none duplicated votes of just quorum ones", func(t *testing.T) {
		seats := 10
		committee, keys := helpers.GenerateCommittee(seats)
		quorum := bft.Quorum(new(big.Int).SetInt64(int64(seats)))
		height := uint64(1)
		round := int64(0)
		noneNilValue := common.Hash{0x1}
		var preVotes []*mUtils.Message
		for i := 0; i < len(committee); i++ {
			preVote := newVoteMsg(height, round, proto.MsgPrevote, keys[committee[i].Address], noneNilValue, committee)
			preVotes = append(preVotes, preVote)
		}

		// let duplicated msg happens, the counting should skip duplicated ones.
		preVotes = append(preVotes, preVotes...)

		overQuorumVotes := OverQuorumVotes(preVotes, quorum.Uint64())
		require.Equal(t, quorum.Uint64(), uint64(len(overQuorumVotes)))
	})

	t.Run("with less quorum votes, it returns no votes", func(t *testing.T) {
		seats := 10
		committee, keys := helpers.GenerateCommittee(seats)
		quorum := bft.Quorum(new(big.Int).SetInt64(int64(seats)))
		height := uint64(1)
		round := int64(0)
		noneNilValue := common.Hash{0x1}
		var preVotes []*mUtils.Message
		for i := 0; i < int(quorum.Uint64()-1); i++ {
			preVote := newVoteMsg(height, round, proto.MsgPrevote, keys[committee[i].Address], noneNilValue, committee)
			preVotes = append(preVotes, preVote)
		}

		// let duplicated msg happens, the counting should skip duplicated ones.
		preVotes = append(preVotes, preVotes...)

		overQuorumVotes := OverQuorumVotes(preVotes, quorum.Uint64())
		require.Nil(t, overQuorumVotes)
	})
}
