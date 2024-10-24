package types

import (
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/stretchr/testify/require"
	"log"
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
)

func TestEpoch_Equal(t *testing.T) {
	consensusKey1, err := blst.RandKey()
	require.NoError(t, err)
	consensusPubKey1 := consensusKey1.PublicKey()
	consensusPubKey1Bytes := consensusPubKey1.Marshal()

	consensusKey2, err := blst.RandKey()
	require.NoError(t, err)
	consensusPubKey2 := consensusKey2.PublicKey()
	consensusPubKey2Bytes := consensusPubKey2.Marshal()

	tests := []struct {
		name   string
		epoch1 *Epoch
		epoch2 *Epoch
		expect bool
	}{
		{
			name:   "both nil epochs",
			epoch1: nil,
			epoch2: nil,
			expect: true,
		},
		{
			name:   "one nil epoch",
			epoch1: nil,
			epoch2: &Epoch{},
			expect: false,
		},
		{
			name:   "different parent blocks",
			epoch1: &Epoch{PreviousEpochBlock: big.NewInt(1), NextEpochBlock: big.NewInt(2), Committee: &Committee{}},
			epoch2: &Epoch{PreviousEpochBlock: big.NewInt(2), NextEpochBlock: big.NewInt(2), Committee: &Committee{}},
			expect: false,
		},
		{
			name: "equal epochs",
			epoch1: &Epoch{
				PreviousEpochBlock: big.NewInt(1),
				NextEpochBlock:     big.NewInt(2),
				Committee:          &Committee{Members: []CommitteeMember{{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes}}},
			},
			epoch2: &Epoch{
				PreviousEpochBlock: big.NewInt(1),
				NextEpochBlock:     big.NewInt(2),
				Committee:          &Committee{Members: []CommitteeMember{{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes}}},
			},
			expect: true,
		},
		{
			name: "unequal epochs - different key",
			epoch1: &Epoch{
				PreviousEpochBlock: big.NewInt(1),
				NextEpochBlock:     big.NewInt(2),
				Committee:          &Committee{Members: []CommitteeMember{{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey2, ConsensusKeyBytes: consensusPubKey2Bytes}}},
			},
			epoch2: &Epoch{
				PreviousEpochBlock: big.NewInt(1),
				NextEpochBlock:     big.NewInt(2),
				Committee:          &Committee{Members: []CommitteeMember{{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes}}},
			},
			expect: false,
		},
		{
			name: "unequal epochs - different committee",
			epoch1: &Epoch{
				PreviousEpochBlock: big.NewInt(1),
				NextEpochBlock:     big.NewInt(2),
				Committee:          nil,
			},
			epoch2: &Epoch{
				PreviousEpochBlock: big.NewInt(1),
				NextEpochBlock:     big.NewInt(2),
				Committee:          &Committee{Members: []CommitteeMember{{Address: common.Address{2}, VotingPower: big.NewInt(10)}}},
			},
			expect: false,
		},
		{
			name: "unequal epochs - different committee",
			epoch1: &Epoch{
				PreviousEpochBlock: big.NewInt(1),
				NextEpochBlock:     big.NewInt(2),
				Committee:          &Committee{Members: []CommitteeMember{{Address: common.Address{1}, VotingPower: big.NewInt(10)}}},
			},
			epoch2: &Epoch{
				PreviousEpochBlock: big.NewInt(1),
				NextEpochBlock:     big.NewInt(2),
				Committee:          &Committee{Members: []CommitteeMember{{Address: common.Address{2}, VotingPower: big.NewInt(10)}}},
			},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.epoch1.Equal(tt.epoch2)
			if result != tt.expect {
				t.Errorf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestCommittee_Copy(t *testing.T) {
	consensusKey1, err := blst.RandKey()
	require.NoError(t, err)
	consensusPubKey1 := consensusKey1.PublicKey()
	consensusPubKey1Bytes := consensusPubKey1.Marshal()

	consensusKey2, err := blst.RandKey()
	require.NoError(t, err)
	consensusPubKey2 := consensusKey2.PublicKey()
	consensusPubKey2Bytes := consensusPubKey2.Marshal()

	original := &Committee{
		Members: []CommitteeMember{
			{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes},
			{Address: common.Address{2}, VotingPower: big.NewInt(20), ConsensusKeyBytes: consensusPubKey2Bytes, ConsensusKey: consensusPubKey2},
		},
	}

	clone := original.Copy()

	if !original.Equal(clone) {
		t.Errorf("expected original and clone to be equal")
	}

	// Modify the clone and check that the original is unaffected
	clone.Members[0].VotingPower.SetUint64(30)

	if original.Members[0].VotingPower.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("original should not be modified")
	}
}

func TestCommittee_Enrich(t *testing.T) {
	consensusKey1, err := blst.RandKey()
	require.NoError(t, err)
	consensusPubKey1 := consensusKey1.PublicKey()
	consensusPubKey1Bytes := consensusPubKey1.Marshal()

	member := CommitteeMember{
		Address:           common.Address{1},
		VotingPower:       big.NewInt(10),
		ConsensusKeyBytes: consensusPubKey1Bytes,
	}

	c := &Committee{
		Members: []CommitteeMember{member},
	}

	err = c.Enrich()
	require.NoError(t, err)

	if committee.Members[0].Index != 0 {
		t.Errorf("expected index to be 0, got %d", committee.Members[0].Index)
	}
}

func TestCommittee_TotalVotingPower(t *testing.T) {
	c := &Committee{
		Members: []CommitteeMember{
			{VotingPower: big.NewInt(10)},
			{VotingPower: big.NewInt(20)},
		},
	}

	total := c.TotalVotingPower()
	expected := big.NewInt(30)

	if total.Cmp(expected) != 0 {
		t.Errorf("expected total voting power to be %v, got %v", expected, total)
	}
}

func TestCommittee_Proposer(t *testing.T) {
	c := &Committee{
		Members: []CommitteeMember{
			{Address: common.Address{1}, VotingPower: big.NewInt(10)},
			{Address: common.Address{2}, VotingPower: big.NewInt(20)},
		},
	}

	proposer := c.Proposer(1, 0)
	if proposer != c.Members[1].Address {
		t.Errorf("expected proposer to be address %v, got %v", c.Members[1].Address, proposer)
	}
}

func TestSortCommitteeMembers(t *testing.T) {
	members := []CommitteeMember{
		{Address: common.Address{1}, VotingPower: big.NewInt(10)},
		{Address: common.Address{2}, VotingPower: big.NewInt(20)},
		{Address: common.Address{3}, VotingPower: big.NewInt(15)},
	}

	SortCommitteeMembers(members)

	if members[0].VotingPower.Cmp(members[1].VotingPower) != 1 {
		t.Errorf("expected first member to have higher voting power")
	}
}

func TestCommittee_Equal(t *testing.T) {
	consensusKey1, err := blst.RandKey()
	require.NoError(t, err)
	consensusPubKey1 := consensusKey1.PublicKey()
	consensusPubKey1Bytes := consensusPubKey1.Marshal()

	consensusKey2, err := blst.RandKey()
	require.NoError(t, err)
	consensusPubKey2 := consensusKey2.PublicKey()
	consensusPubKey2Bytes := consensusPubKey2.Marshal()

	tests := []struct {
		name       string
		committee1 *Committee
		committee2 *Committee
		expect     bool
	}{
		{
			name:       "both nil committees",
			committee1: nil,
			committee2: nil,
			expect:     true,
		},
		{
			name:       "one nil committee",
			committee1: nil,
			committee2: &Committee{},
			expect:     false,
		},
		{
			name:       "different number of members",
			committee1: &Committee{Members: []CommitteeMember{{Address: common.Address{1}}}},
			committee2: &Committee{Members: []CommitteeMember{{Address: common.Address{1}}, {Address: common.Address{2}}}},
			expect:     false,
		},
		{
			name: "equal committees",
			committee1: &Committee{
				Members: []CommitteeMember{
					{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes},
					{Address: common.Address{2}, VotingPower: big.NewInt(20), ConsensusKey: consensusPubKey2, ConsensusKeyBytes: consensusPubKey2Bytes},
				},
			},
			committee2: &Committee{
				Members: []CommitteeMember{
					{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes},
					{Address: common.Address{2}, VotingPower: big.NewInt(20), ConsensusKey: consensusPubKey2, ConsensusKeyBytes: consensusPubKey2Bytes},
				},
			},
			expect: true,
		},
		{
			name: "unequal committees - different consensus key",
			committee1: &Committee{
				Members: []CommitteeMember{
					{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes},
					{Address: common.Address{2}, VotingPower: big.NewInt(20), ConsensusKey: consensusPubKey2, ConsensusKeyBytes: consensusPubKey2Bytes},
				},
			},
			committee2: &Committee{
				Members: []CommitteeMember{
					{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey2, ConsensusKeyBytes: consensusPubKey2Bytes},
					{Address: common.Address{2}, VotingPower: big.NewInt(20), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes},
				},
			},
			expect: false,
		},
		{
			name: "unequal committees - different voting power",
			committee1: &Committee{
				Members: []CommitteeMember{
					{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes},
					{Address: common.Address{2}, VotingPower: big.NewInt(20), ConsensusKey: consensusPubKey2, ConsensusKeyBytes: consensusPubKey2Bytes},
				},
			},
			committee2: &Committee{
				Members: []CommitteeMember{
					{Address: common.Address{1}, VotingPower: big.NewInt(10), ConsensusKey: consensusPubKey1, ConsensusKeyBytes: consensusPubKey1Bytes},
					{Address: common.Address{2}, VotingPower: big.NewInt(30), ConsensusKey: consensusPubKey2, ConsensusKeyBytes: consensusPubKey2Bytes},
				},
			},
			expect: false,
		},
		{
			name: "unequal committees - different addresses",
			committee1: &Committee{
				Members: []CommitteeMember{
					{Address: common.Address{1}, VotingPower: big.NewInt(10)},
					{Address: common.Address{2}, VotingPower: big.NewInt(20)},
				},
			},
			committee2: &Committee{
				Members: []CommitteeMember{
					{Address: common.Address{3}, VotingPower: big.NewInt(10)},
					{Address: common.Address{2}, VotingPower: big.NewInt(20)},
				},
			},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.committee1.Equal(tt.committee2)
			if result != tt.expect {
				t.Errorf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestElectProposer(t *testing.T) {
	samePowers := []int{100, 100, 100, 100}
	linearPowers := []int{100, 200, 400, 800}

	t.Run("Proposer election should be deterministic", func(t *testing.T) {
		c := generateCommittee(samePowers)
		for h := uint64(0); h < uint64(100); h++ {
			for r := int64(0); r <= int64(3); r++ {
				proposer1 := c.Proposer(h, r)
				proposer2 := c.Proposer(h, r)
				require.Equal(t, proposer1, proposer2)
			}
		}
	})

	t.Run("Proposer selection, print and compare the scheduling rate with same stake", func(t *testing.T) {
		c := generateCommittee(samePowers)
		maxHeight := uint64(10000)
		maxRound := int64(4)
		//expectedRatioDelta := float64(0.01)
		counterMap := make(map[common.Address]int)
		counterMap[common.Address{}] = 1
		for h := uint64(0); h < maxHeight; h++ {
			for round := int64(0); round < maxRound; round++ {
				proposer := c.Proposer(h, round)
				_, ok := counterMap[proposer]
				if ok {
					counterMap[proposer]++
				} else {
					counterMap[proposer] = 1
				}
			}
		}

		totalStake := 0
		for _, s := range samePowers {
			totalStake += s
		}

		for i, c := range c.Members {
			stake := samePowers[i]
			scheduled := counterMap[c.Address]
			log.Print("electing ", "proposer: ", c.Address.String(), " stake: ", stake, " scheduled: ", scheduled)
		}
	})

	t.Run("Proposer selection, print and compare the scheduling rate with liner increasing stake", func(t *testing.T) {
		c := generateCommittee(linearPowers)
		maxHeight := uint64(1000000)
		maxRound := int64(4)
		//expectedRatioDelta := float64(0.01)
		counterMap := make(map[common.Address]int)
		counterMap[common.Address{}] = 1
		for h := uint64(0); h < maxHeight; h++ {
			for round := int64(0); round < maxRound; round++ {
				proposer := c.Proposer(h, round)
				_, ok := counterMap[proposer]
				if ok {
					counterMap[proposer]++
				} else {
					counterMap[proposer] = 1
				}
			}
		}

		totalStake := 0
		for _, s := range samePowers {
			totalStake += s
		}

		for _, m := range c.Members {
			stake := m.VotingPower.Uint64()
			scheduled := counterMap[m.Address]
			log.Print("electing ", "proposer: ", m.Address.String(), " stake: ", stake, " scheduled: ", scheduled)
		}
	})
}

func generateCommittee(powers []int) *Committee {
	vals := make([]CommitteeMember, len(powers))
	for i, p := range powers {
		privateKey, _ := crypto.GenerateKey()
		committeeMember := CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetInt64(int64(p)),
		}
		vals[i] = committeeMember
	}
	c := &Committee{Members: vals}
	SortCommitteeMembers(c.Members)
	return c
}
