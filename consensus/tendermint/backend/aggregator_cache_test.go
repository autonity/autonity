package backend

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/internal/testrand"
)

type mockCommitteeMember struct {
	types.CommitteeMember
	privateKey blst.SecretKey
}

func (m *mockCommitteeMember) Sign(hash common.Hash) blst.Signature {
	return m.privateKey.Sign(hash[:])
}

func fromCommittee(committee *types.Committee) []mockCommitteeMember {
	members := make([]mockCommitteeMember, len(committee.Members))
	for i, member := range committee.Members {
		key, err := blst.RandKey()
		if err != nil {
			panic(errors.New("failed to generate random key for mock committee member"))
		}
		members[i] = mockCommitteeMember{
			CommitteeMember: member,
			privateKey:      key,
		}
	}
	return members
}
func TestAggregatorCache(t *testing.T) {
	// Initialize the aggregator cache
	cache := newAggregatorCache()
	require.NoError(t, committee.Enrich())
	mockCommittee := fromCommittee(committee)

	t.Run("should add a prevote to the correct round and height", func(t *testing.T) {
		// create a mock message
		r := int64(0)
		h := uint64(100)
		value := testrand.Hash()
		signer := mockCommittee[0].Sign
		prevote := message.NewPrevote(r, h, value, signer, &committee.Members[0], committee.Len())

		cache.MarkCommittee(h, committee)
		cache.AddVote(prevote, StepReceived)

		require.True(t, cache.voteCaches[message.PrevoteCode][StepReceived].Contains(h, r, committee.Len(), prevote))
	})

	t.Run("it should not add a prevote when a vote with a superset of signers was added", func(t *testing.T) {
		// create a mock message
		r := int64(0)
		h := uint64(100)
		value := testrand.Hash()

		prevote1 := message.NewPrevote(r, h, value, mockCommittee[0].Sign, &committee.Members[0], committee.Len())
		prevote2 := message.NewPrevote(r, h, value, mockCommittee[1].Sign, &committee.Members[1], committee.Len())

		aggregatedVote := message.AggregatePrevotes([]message.Vote{prevote1, prevote2})

		cache.MarkCommittee(h, committee)
		cache.AddVote(aggregatedVote, StepReceived)

		require.True(t, cache.voteCaches[message.PrevoteCode][StepReceived].Contains(h, r, committee.Len(), prevote1))
		require.True(t, cache.voteCaches[message.PrevoteCode][StepReceived].Contains(h, r, committee.Len(), prevote2))

		// This should filter the prevote, as the signer is already included in the aggregated vote
		errCh := make(chan<- error)
		defer close(errCh)

		prevoteFiltered := cache.Filter(
			committee.Len(),
			events.UnverifiedMessageEvent{
				Message: prevote1,
				Sender:  committee.Members[0].Address,
				Posted:  time.Now(),
				ErrCh:   errCh,
			},
		)
		require.True(t, prevoteFiltered)

		require.Len(t, cache.EmptyFiltered(h, r, message.PrevoteCode, value), 1)
		// should be empty now
		require.Len(t, cache.filtered[h][filteredCacheKey{round: r, code: message.PrevoteCode, value: value}], 0)
	})

	t.Run("should discard messages if dispatched and not save them in queue", func(t *testing.T) {
		// create a mock message
		r := int64(0)
		h := uint64(100)
		value := testrand.Hash()
		signer := mockCommittee[0].Sign
		prevote := message.NewPrevote(r, h, value, signer, &committee.Members[0], committee.Len())

		cache.MarkCommittee(h, committee)
		cache.AddVote(prevote, StepDispatched)

		require.False(t, cache.voteCaches[message.PrevoteCode][StepReceived].Contains(h, r, committee.Len(), prevote))
		require.True(t, cache.voteCaches[message.PrevoteCode][StepDispatched].Contains(h, r, committee.Len(), prevote))

		errCh := make(chan<- error)
		defer close(errCh)

		require.True(t, cache.Filter(
			committee.Len(),
			events.UnverifiedMessageEvent{
				Message: prevote,
				Sender:  committee.Members[0].Address,
				Posted:  time.Now(),
				ErrCh:   errCh,
			},
		))
		require.Len(t, cache.filtered[h][filteredCacheKey{round: r, code: message.PrevoteCode, value: value}], 0)
	})

}
