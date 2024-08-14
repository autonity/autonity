package types

import (
	"github.com/autonity/autonity/crypto/blst"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"

	"github.com/autonity/autonity/common"
)

func TestEpoch_IsEpochHeader(t *testing.T) {
	tests := []struct {
		name   string
		epoch  Epoch
		expect bool
	}{
		{
			name:   "valid epoch header",
			epoch:  Epoch{ParentEpochBlock: big.NewInt(1), NextEpochBlock: big.NewInt(2), Committee: &Committee{Members: []CommitteeMember{{}}}},
			expect: true,
		},
		{
			name:   "missing parent block",
			epoch:  Epoch{NextEpochBlock: big.NewInt(2), Committee: &Committee{Members: []CommitteeMember{{}}}},
			expect: false,
		},
		{
			name:   "missing next block",
			epoch:  Epoch{ParentEpochBlock: big.NewInt(1), Committee: &Committee{Members: []CommitteeMember{{}}}},
			expect: false,
		},
		{
			name:   "empty committee",
			epoch:  Epoch{ParentEpochBlock: big.NewInt(1), NextEpochBlock: big.NewInt(2), Committee: &Committee{Members: []CommitteeMember{}}},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.epoch.IsEpochHeader()
			if result != tt.expect {
				t.Errorf("expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestCommittee_Copy(t *testing.T) {
	original := &Committee{
		Members: []CommitteeMember{
			{Address: common.Address{1}, VotingPower: big.NewInt(10)},
			{Address: common.Address{2}, VotingPower: big.NewInt(20)},
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
	member := CommitteeMember{
		Address:           common.Address{1},
		VotingPower:       big.NewInt(10),
		ConsensusKeyBytes: []byte{0x01, 0x02, 0x03},
	}

	committee := &Committee{
		Members: []CommitteeMember{member},
	}

	err := committee.Enrich()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if committee.Members[0].Index != 0 {
		t.Errorf("expected index to be 0, got %d", committee.Members[0].Index)
	}
}

func TestCommittee_TotalVotingPower(t *testing.T) {
	committee := &Committee{
		Members: []CommitteeMember{
			{VotingPower: big.NewInt(10)},
			{VotingPower: big.NewInt(20)},
		},
	}

	total := committee.TotalVotingPower()
	expected := big.NewInt(30)

	if total.Cmp(expected) != 0 {
		t.Errorf("expected total voting power to be %v, got %v", expected, total)
	}
}

func TestCommittee_Proposer(t *testing.T) {
	committee := &Committee{
		Members: []CommitteeMember{
			{Address: common.Address{1}, VotingPower: big.NewInt(10)},
			{Address: common.Address{2}, VotingPower: big.NewInt(20)},
		},
	}

	proposer := committee.Proposer(1, 0)
	if proposer != committee.Members[1].Address {
		t.Errorf("expected proposer to be address %v, got %v", committee.Members[1].Address, proposer)
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
