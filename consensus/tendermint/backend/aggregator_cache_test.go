package backend

import (
	"errors"
	"math/big"
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

type mockCommittee []mockCommitteeMember

func (mc mockCommittee) ToCommittee() *types.Committee {
	members := make([]types.CommitteeMember, len(mc))
	for i, member := range mc {
		members[i] = member.CommitteeMember
	}
	return &types.Committee{
		Members: members,
	}
}

func fromCommittee(committee *types.Committee) []mockCommitteeMember {
	members := make([]mockCommitteeMember, len(committee.Members))
	for i, member := range committee.Members {
		key, err := blst.RandKey()
		if err != nil {
			panic(errors.New("failed to generate random key for mock committee member"))
		}
		members[i] = mockCommitteeMember{
			CommitteeMember: types.CommitteeMember{
				Address:      member.Address,
				VotingPower:  member.VotingPower,
				ConsensusKey: key.PublicKey(),
				Index:        member.Index,
			},
			privateKey: key,
		}
	}
	return members
}
func TestAggregatorCache(t *testing.T) {
	require.NoError(t, committee.Enrich())
	mockCommittee := fromCommittee(committee)

	t.Run("should add a prevote to the correct round and height", func(t *testing.T) {
		// Initialize the aggregator cache
		cache := newAggregatorCache()

		// create a mock message
		r := int64(0)
		h := uint64(100)
		value := testrand.Hash()
		signer := mockCommittee[0].Sign
		prevote := message.NewPrevote(r, h, value, signer, &committee.Members[0], committee.Len())

		cache.markCommittee(h, committee)
		cache.addVote(prevote, stepReceived)

		require.True(t, cache.voteCaches[message.PrevoteCode][stepReceived].contains(
			h, r, value, prevote.Signers().Bits.ToSingleBitmap(committee.Len()),
		))
	})

	t.Run("it should not add a prevote when a vote with a superset of signers was added", func(t *testing.T) {
		// Initialize the aggregator cache
		cache := newAggregatorCache()

		// create a mock message
		r := int64(0)
		h := uint64(100)
		value := testrand.Hash()

		prevote1 := message.NewPrevote(r, h, value, mockCommittee[0].Sign, &committee.Members[0], committee.Len())
		prevote2 := message.NewPrevote(r, h, value, mockCommittee[1].Sign, &committee.Members[1], committee.Len())

		aggregatedVote := message.AggregatePrevotes([]message.Vote{prevote1, prevote2})

		cache.markCommittee(h, committee)
		cache.addVote(aggregatedVote, stepReceived)

		require.True(t, cache.voteCaches[message.PrevoteCode][stepReceived].contains(
			h,
			r,
			prevote1.Value(),
			prevote1.Signers().Bits.ToSingleBitmap(committee.Len()),
		))
		require.True(t, cache.voteCaches[message.PrevoteCode][stepReceived].contains(
			h,
			r,
			prevote2.Value(),
			prevote2.Signers().Bits.ToSingleBitmap(committee.Len()),
		))

		// This should filter the prevote, as the signer is already included in the aggregated vote
		errCh := make(chan<- error)
		defer close(errCh)

		prevoteFiltered := cache.filter(
			committee.Len(),
			events.UnverifiedMessageEvent{
				Message: prevote1,
				Sender:  committee.Members[0].Address,
				Posted:  time.Now(),
				ErrCh:   errCh,
			},
		)
		require.True(t, prevoteFiltered)

		require.Len(t, cache.emptyFiltered(h, r, message.PrevoteCode, value), 1)
		// should be empty now
		require.Len(t, cache.filtered[h][filteredCacheKey{round: r, code: message.PrevoteCode, value: value}], 0)
	})

	t.Run("should discard messages if dispatched and not save them in queue", func(t *testing.T) {
		// Initialize the aggregator cache
		cache := newAggregatorCache()

		// create a mock message
		r := int64(0)
		h := uint64(100)
		value := testrand.Hash()
		signer := mockCommittee[0].Sign
		prevote := message.NewPrevote(r, h, value, signer, &committee.Members[0], committee.Len())

		cache.markCommittee(h, committee)
		cache.addVote(prevote, stepDispatched)

		require.False(t, cache.voteCaches[message.PrevoteCode][stepReceived].contains(
			h,
			r,
			prevote.Value(),
			prevote.Signers().Bits.ToSingleBitmap(committee.Len()),
		))
		require.True(t, cache.voteCaches[message.PrevoteCode][stepDispatched].contains(
			h,
			r,
			prevote.Value(),
			prevote.Signers().Bits.ToSingleBitmap(committee.Len()),
		))

		errCh := make(chan<- error)
		defer close(errCh)

		require.True(t, cache.filter(
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

func TestAggregatorCachePowerCalculations(t *testing.T) {
	require.NoError(t, committee.Enrich())
	mc := fromCommittee(committee)

	t.Run("should calculate the correct power for a specific value", func(t *testing.T) {
		// Initialize the aggregator cache
		cache := newAggregatorCache()

		r := int64(0)
		h := uint64(100)
		value := testrand.Hash()

		votersA := []int{1, 3, 5} // Mock voters from the committee
		votersB := []int{0, 1, 4} // Another set of mock voters
		var votesA []message.Vote
		for _, voter := range votersA {
			votesA = append(votesA, message.NewPrevote(r, h, value, mc[voter].Sign, &committee.Members[voter], committee.Len()))
		}
		voteA := message.AggregatePrevotes(votesA)

		var votesB []message.Vote
		for _, voter := range votersB {
			votesB = append(votesB, message.NewPrevote(r, h, value, mc[voter].Sign, &committee.Members[voter], committee.Len()))
		}
		voteB := message.AggregatePrevotes(votesB)

		cache.markCommittee(h, committee)

		cache.addVote(voteA, stepReceived)
		cache.addVote(voteB, stepReceived)

		power := cache.presentPowerForValue(h, r, value, message.PrevoteCode, stepReceived)
		require.NotNil(t, power, "power should not be nil")

		expectedPower := big.NewInt(0)
		for _, voter := range []int{0, 1, 3, 4, 5} {
			expectedPower.Add(expectedPower, committee.Members[voter].VotingPower)
		}
		require.Equal(t, expectedPower, power, "should match the expected voting power for the value")
	})

	t.Run("should calculate the proper round power for all messages", func(t *testing.T) {
		// Initialize the aggregator cache
		cache := newAggregatorCache()

		r := int64(0)
		h := uint64(100)
		cache.markCommittee(h, committee)

		valueA := testrand.Hash()
		valueB := testrand.Hash()

		proposers := []int{5, 6}  // Mock proposers from the committee
		votersA := []int{1, 3, 5} // Mock voters from the committee
		votersB := []int{0, 1, 4} // Another set of mock voters
		for _, proposer := range proposers {
			proposal := message.NewPropose(
				r, h, -1,
				types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(h))}),
				mc[proposer].Sign,
				&committee.Members[proposer],
			)
			cache.addEvent(events.UnverifiedMessageEvent{Message: proposal}, stepReceived)
		}
		var votesA []message.Vote
		for _, voter := range votersA {
			votesA = append(votesA, message.NewPrevote(r, h, valueA, mc[voter].Sign, &committee.Members[voter], committee.Len()))
		}
		voteA := message.AggregatePrevotes(votesA)

		var votesB []message.Vote
		for _, voter := range votersB {
			votesB = append(votesB, message.NewPrecommit(r, h, valueB, mc[voter].Sign, &committee.Members[voter], committee.Len()))
		}
		voteB := message.AggregatePrecommits(votesB)

		cache.addEvent(events.UnverifiedMessageEvent{Message: voteA}, stepReceived)
		cache.addEvent(events.UnverifiedMessageEvent{Message: voteB}, stepReceived)

		power := cache.totalPowerForRound(h, r, stepReceived)
		require.NotNil(t, power, "power should not be nil")

		expectedPower := big.NewInt(0)
		for _, voter := range []int{0, 1, 3, 4, 5, 6} {
			expectedPower.Add(expectedPower, committee.Members[voter].VotingPower)
		}

		require.Equal(t, expectedPower, power, "should match the expected voting power for the round")
	})

	t.Run("should calculate the correct power for a specific code", func(t *testing.T) {
		cache := newAggregatorCache()

		// add a bunch of stuff to the cache
		r := int64(0)
		h := uint64(100)
		value := testrand.Hash()

		cache.markCommittee(h, mockCommittee(mc).ToCommittee())
		cache.markCommittee(h-1, mockCommittee(mc).ToCommittee())

		voteRound0Prevote := newSignedTestMsg(t, h, r, value, message.PrevoteCode, mc, []int{0, 1, 2})
		voteRound0PrevoteB := newSignedTestMsg(t, h, r, common.Hash{}, message.PrevoteCode, mc, []int{3, 4})
		voteRound0Precommit := newSignedTestMsg(t, h, r, value, message.PrecommitCode, mc, []int{0, 1, 4, 5})
		voteRound1 := newSignedTestMsg(t, h, r+1, value, message.PrevoteCode, mc, []int{0, 1, 2})
		voteRound2 := newSignedTestMsg(t, h, r+2, value, message.PrevoteCode, mc, []int{0, 1})
		voteRound3 := newSignedTestMsg(t, h-1, r+3, value, message.PrevoteCode, mc, []int{0})

		for _, event := range []events.UnverifiedMessageEvent{voteRound0Prevote, voteRound0PrevoteB, voteRound0Precommit, voteRound1, voteRound2, voteRound3} {
			cache.addEvent(event, stepReceived)
		}
		powerPrevote := cache.totalPowerForCode(h, r, message.PrevoteCode, stepReceived)
		require.NotNil(t, powerPrevote, "power should not be nil for prevote")

		// should be 5, as we have 5 members with voting power in the prevote round (for different values)
		require.Equal(t, big.NewInt(5), powerPrevote, "should match the expected voting power for the prevote code")
	})
}

func TestBitmapLogic(t *testing.T) {
	t.Run("signers should align with bitmap indexes", func(t *testing.T) {
		signers := types.NewSigners(committee.Len())
		signers.Increment(&committee.Members[0])
		signers.Increment(&committee.Members[3])
		signers.Increment(&committee.Members[5])

		bm := oneBitmap(signers.Bits.ToSingleBitmap(committee.Len()))
		indexes := bm.presentIndexes()
		require.Equal(t, 3, len(indexes), "should have 3 indexes present in the bitmap")
		require.Equal(t, []int{0, 3, 5}, indexes, "indexes should match the signers present in the bitmap")
	})

	t.Run("should properly merge bitmaps with different sets", func(t *testing.T) {
		signersA := types.NewSigners(committee.Len())
		signersA.Increment(&committee.Members[0])
		signersA.Increment(&committee.Members[3])
		signersA.Increment(&committee.Members[5])

		bmA := oneBitmap(signersA.Bits.ToSingleBitmap(committee.Len()))

		signersB := types.NewSigners(committee.Len())
		signersB.Increment(&committee.Members[1])
		signersB.Increment(&committee.Members[4])
		signersB.Increment(&committee.Members[5])
		bmB := oneBitmap(signersB.Bits.ToSingleBitmap(committee.Len()))

		require.False(t, bmA.contains(bmB))

		merged := bmA.merge(bmB)
		require.True(t, merged.contains(bmA), "merged bitmap should contain the first bitmap")
		require.True(t, merged.contains(bmB), "merged bitmap should contain the second bitmap")

		indexes := merged.presentIndexes()
		require.Equal(t, []int{0, 1, 3, 4, 5}, indexes, "merged bitmap should have indexes from both bitmaps")
	})

}
