package committee

import (
	"fmt"
	"github.com/clearmatics/autonity/core/types"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
)

var (
	validator1 = "70524d664ffe731100208a0154e556f9bb679ae6"
	validator2 = "b37866a925bccd69cfa98d43b510f1d23d78a851"
	validator3 = "70524d664ffe731100208a0154e556f9bb679ae0"
	validator4 = "b37866a925bccd69cfa98d43b510f1d23d78a850"
	validator5 = "70524d664ffe731100208a0154e556f9bb679ae1"
	validator6 = "b37866a925bccd69cfa98d43b510f1d23d78a853"
)

func TestWeightedRoundRobinProposer(t *testing.T) {
	testZeroVotingPower(t)
	testWRRDetermination(t)
	testNotScheduleZeroStakeHolder(t)
	testWRRSchedule(t)
}

func TestWeightedRoundRobinProposerOperatorAttackVector(t *testing.T) {
	b1 := common.Hex2Bytes(validator1)
	b2 := common.Hex2Bytes(validator2)
	b3 := common.Hex2Bytes(validator3)
	b4 := common.Hex2Bytes(validator4)
	b5 := common.Hex2Bytes(validator5)
	b6 := common.Hex2Bytes(validator6)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	addr3 := common.BytesToAddress(b3)
	addr4 := common.BytesToAddress(b4)
	addr5 := common.BytesToAddress(b5)
	addr6 := common.BytesToAddress(b6)
	val1 := types.CommitteeMember{Address: addr1, VotingPower: new(big.Int).SetUint64(100)}
	val2 := types.CommitteeMember{Address: addr2, VotingPower: new(big.Int).SetUint64(100)}
	val3 := types.CommitteeMember{Address: addr3, VotingPower: new(big.Int).SetUint64(100)}
	val4 := types.CommitteeMember{Address: addr4, VotingPower: new(big.Int).SetUint64(100)}
	val5 := types.CommitteeMember{Address: addr5, VotingPower: new(big.Int).SetUint64(100)}
	val6 := types.CommitteeMember{Address: addr6, VotingPower: new(big.Int).SetUint64(100)}

	//val1 - malicious

	vals := types.Committee{val1, val2, val3, val4, val5, val6}
	committeeSet, err := NewSet(types.Committee{val1, val2, val3, val4}, config.WeightedRoundRobin, val2.Address)
	if committeeSet == nil || err != nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}
	h := int64(1)
	r := int64(1)
	nextCommitee := vals[:4]

	nextLastProposer := addr2
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		committeeSet, _ = NewSet(nextCommitee, config.WeightedRoundRobin, nextLastProposer)
		c := committeeSet.GetProposer(r, big.NewInt(h))
		nextLastProposer = c.Address
		if c.Address == addr1 {
			fmt.Println("mal")
		} else {
			//fmt.Println("not")
		}
		for j := 5; j > 3; j-- {
			committeeSet, _ = NewSet(vals[:j], config.WeightedRoundRobin, c.Address)
			c = committeeSet.GetProposer(r, big.NewInt(h))
			if c.Address == addr1 {
				nextCommitee = vals[:j]
				//fmt.Println("c mal")
				break
			} else {
				//fmt.Println("c not")
			}
		}
		h++
	}
}

func TestWeightedRoundRobinProposerValidatorsAttackVector(t *testing.T) {
	b1 := common.Hex2Bytes(validator1)
	b2 := common.Hex2Bytes(validator2)
	b3 := common.Hex2Bytes(validator3)
	b4 := common.Hex2Bytes(validator4)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	addr3 := common.BytesToAddress(b3)
	addr4 := common.BytesToAddress(b4)
	val1 := types.CommitteeMember{Address: addr1, VotingPower: new(big.Int).SetUint64(100)}
	val2 := types.CommitteeMember{Address: addr2, VotingPower: new(big.Int).SetUint64(100)}
	val3 := types.CommitteeMember{Address: addr3, VotingPower: new(big.Int).SetUint64(100)}
	val4 := types.CommitteeMember{Address: addr4, VotingPower: new(big.Int).SetUint64(200)}
	committeeSet, err := NewSet(types.Committee{val1, val2, val3, val4}, config.WeightedRoundRobin, val1.Address)
	if committeeSet == nil || err != nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}
	// v1, v2  is maligous
	c := committeeSet.GetProposer(1, big.NewInt(1))
	t.Log("round ", 1, "height", 1, c, "maligous validator block")
	//ok
	c2 := committeeSet.GetProposer(1, big.NewInt(2))
	t.Log("round ", 1, "height", 2, c2, "maligous validator block")
	//ok
	c3 := committeeSet.GetProposer(1, big.NewInt(3))
	t.Log("round ", 1, "height", 3, c3, "vote agains it")
	//attack v1&v2 against block

	c4 := committeeSet.GetProposer(2, big.NewInt(3))
	t.Log("round ", 2, "height", 3, c4, "maligous validator block")

	//attack v1&v2 against block
	c5 := committeeSet.GetProposer(1, big.NewInt(4))
	t.Log("round ", 1, "height", 4, c5, "maligous validator block")
}

func TestWeightedRoundRobinProposerValidatorsAttackVector2(t *testing.T) {
	b1 := common.Hex2Bytes(validator1)
	b2 := common.Hex2Bytes(validator2)
	b3 := common.Hex2Bytes(validator3)
	b4 := common.Hex2Bytes(validator4)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	addr3 := common.BytesToAddress(b3)
	addr4 := common.BytesToAddress(b4)
	val1 := types.CommitteeMember{Address: addr1, VotingPower: new(big.Int).SetUint64(50)}
	val2 := types.CommitteeMember{Address: addr2, VotingPower: new(big.Int).SetUint64(40)}
	val3 := types.CommitteeMember{Address: addr3, VotingPower: new(big.Int).SetUint64(100)}
	val4 := types.CommitteeMember{Address: addr4, VotingPower: new(big.Int).SetUint64(200)}
	committeeSet, err := NewSet(types.Committee{val1, val2, val3, val4}, config.WeightedRoundRobin, val1.Address)
	if committeeSet == nil || err != nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}
	h := int64(1)
	r := int64(1)
	for i := 0; i < 100; i++ {
		c := committeeSet.GetProposer(r, big.NewInt(h))
		if c.Address == addr1 || c.Address == addr2 {
			t.Log("round ", r, "height", h, c, "maligous validator block")
			h++
			r = 1
		} else {
			t.Log("round ", r, "height", h, c, "vote agains it")
			r++
		}
	}
}

func testWRRDetermination(t *testing.T) {
	b1 := common.Hex2Bytes(validator1)
	b2 := common.Hex2Bytes(validator2)
	b3 := common.Hex2Bytes(validator3)
	b4 := common.Hex2Bytes(validator4)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	addr3 := common.BytesToAddress(b3)
	addr4 := common.BytesToAddress(b4)
	val1 := types.CommitteeMember{Address: addr1, VotingPower: new(big.Int).SetUint64(100)}
	val2 := types.CommitteeMember{Address: addr2, VotingPower: new(big.Int).SetUint64(100)}
	val3 := types.CommitteeMember{Address: addr3, VotingPower: new(big.Int).SetUint64(100)}
	val4 := types.CommitteeMember{Address: addr4, VotingPower: new(big.Int).SetUint64(200)}

	committeeSet, err := NewSet(types.Committee{val1, val2, val3, val4}, config.WeightedRoundRobin, val1.Address)
	if committeeSet == nil || err != nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	for height := int64(0); height < 1000; height++ {
		for round := int64(0); round < 100; round++ {
			v1 := committeeSet.GetProposer(round, big.NewInt(height))
			v2 := committeeSet.GetProposer(round, big.NewInt(height))
			if !reflect.DeepEqual(v1, v2) {
				t.Errorf("validator mismatch: have %v, want %v", v1, v2)
			}
		}

	}
}

func testZeroVotingPower(t *testing.T) {
	b1 := common.Hex2Bytes(validator1)
	b2 := common.Hex2Bytes(validator2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	val1 := types.CommitteeMember{Address: addr1, VotingPower: new(big.Int).SetUint64(0)}
	val2 := types.CommitteeMember{Address: addr2, VotingPower: new(big.Int).SetUint64(0)}

	committeeSet, err := NewSet(types.Committee{val1, val2}, config.WeightedRoundRobin, val1.Address)
	if committeeSet == nil || err != nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	for height := int64(0); height < 1000; height++ {
		for round := int64(0); round < 100; round++ {
			if round%2 == 0 {
				if validator := committeeSet.GetProposer(round, big.NewInt(height)); !reflect.DeepEqual(validator, val2) {
					t.Errorf("validator mismatch: have %v, want %v", validator, val2)
				}

			} else {
				if validator := committeeSet.GetProposer(round, big.NewInt(height)); !reflect.DeepEqual(validator, val1) {
					t.Errorf("validator mismatch: have %v, want %v", validator, val1)
				}
			}
		}
	}

}

func testNotScheduleZeroStakeHolder(t *testing.T) {
	b1 := common.Hex2Bytes(validator1)
	b2 := common.Hex2Bytes(validator2)
	b3 := common.Hex2Bytes(validator3)
	b4 := common.Hex2Bytes(validator4)
	b5 := common.Hex2Bytes(validator5)
	b6 := common.Hex2Bytes(validator6)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	addr3 := common.BytesToAddress(b3)
	addr4 := common.BytesToAddress(b4)
	addr5 := common.BytesToAddress(b5)
	addr6 := common.BytesToAddress(b6)
	val1 := types.CommitteeMember{Address: addr1, VotingPower: new(big.Int).SetUint64(100)}
	val2 := types.CommitteeMember{Address: addr2, VotingPower: new(big.Int).SetUint64(0)}
	val3 := types.CommitteeMember{Address: addr3, VotingPower: new(big.Int).SetUint64(0)}
	val4 := types.CommitteeMember{Address: addr4, VotingPower: new(big.Int).SetUint64(200)}
	val5 := types.CommitteeMember{Address: addr5, VotingPower: new(big.Int).SetUint64(100)}
	val6 := types.CommitteeMember{Address: addr6, VotingPower: new(big.Int).SetUint64(0)}

	committeeSet, err := NewSet(types.Committee{val1, val2, val3, val4, val5, val6}, config.WeightedRoundRobin, val1.Address)
	if committeeSet == nil || err != nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	for height := int64(0); height < 10000; height++ {
		for round := int64(0); round < 10; round++ {
			proposer := committeeSet.GetProposer(round, big.NewInt(height))
			if reflect.DeepEqual(proposer, val2) {
				t.Errorf("scheduled zero stake validator: %v, checking %v", proposer, val2)
			}

			if reflect.DeepEqual(proposer, val3) {
				t.Errorf("scheduled zero stake validator: %v, checking %v", proposer, val3)
			}

			if reflect.DeepEqual(proposer, val6) {
				t.Errorf("scheduled zero stake validator: %v, checking %v", proposer, val6)
			}
		}
	}
}

func testWRRSchedule(t *testing.T) {
	b1 := common.Hex2Bytes(validator1)
	b2 := common.Hex2Bytes(validator2)
	b3 := common.Hex2Bytes(validator3)
	b4 := common.Hex2Bytes(validator4)
	b5 := common.Hex2Bytes(validator5)
	b6 := common.Hex2Bytes(validator6)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)
	addr3 := common.BytesToAddress(b3)
	addr4 := common.BytesToAddress(b4)
	addr5 := common.BytesToAddress(b5)
	addr6 := common.BytesToAddress(b6)
	val1 := types.CommitteeMember{Address: addr1, VotingPower: new(big.Int).SetUint64(100)}
	val2 := types.CommitteeMember{Address: addr2, VotingPower: new(big.Int).SetUint64(100)}

	val3 := types.CommitteeMember{Address: addr3, VotingPower: new(big.Int).SetUint64(200)}
	val4 := types.CommitteeMember{Address: addr4, VotingPower: new(big.Int).SetUint64(100)}
	val5 := types.CommitteeMember{Address: addr5, VotingPower: new(big.Int).SetUint64(100)}
	val6 := types.CommitteeMember{Address: addr6, VotingPower: new(big.Int).SetUint64(200)}

	committeeSet, err := NewSet(types.Committee{val1, val2, val3, val4, val5, val6}, config.WeightedRoundRobin, val1.Address)
	if committeeSet == nil || err != nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	totalPower := uint64(0)
	valSet := committeeSet.Committee()
	for i := range valSet {
		totalPower += valSet[i].VotingPower.Uint64()
	}

	expectedRates := make(map[common.Address]float64)
	for i := range valSet {
		expectedRates[valSet[i].Address] = float64(valSet[i].VotingPower.Uint64()) / float64(totalPower)
	}

	mapHits := make(map[common.Address]int64)

	maxHeight := 100000
	maxRound := 10
	totalElection := int64(maxHeight * maxRound)
	for height := int64(0); height < int64(maxHeight); height++ {
		for round := int64(0); round < int64(maxRound); round++ {
			proposer := committeeSet.GetProposer(round, big.NewInt(height))
			_, ok := mapHits[proposer.Address]
			if !ok {
				mapHits[proposer.Address] = 1
			} else {
				mapHits[proposer.Address]++
			}
		}
	}

	for k, scheduled := range mapHits {
		expected := expectedRates[k] * 100
		actualRate := (float64(scheduled) / float64(totalElection)) * 100
		t.Logf("address: %s, scheduled: %d times, expected rate: %f%%, actual rate: %f%%.", k.String(),
			scheduled, expected, actualRate)
		// if the schedule rate is more than 1% unexpected, fail the case.
		delta := expected - actualRate
		if math.Abs(delta) > 1.0 {
			t.Errorf("the schedule rate delta is unexpected: %f", math.Abs(delta))
		}
	}
}
