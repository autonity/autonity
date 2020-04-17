package committee

import (
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/core/types"
	"math/big"
	"testing"
)

//Idea: using add/remove validators we can make proposing of one of the validators more frequent
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

	numOfMaliciousValidators := 0
	nextLastProposer := addr2
	numOfBlocks := 10000
	for i := 0; i < numOfBlocks; i++ {
		committeeSet, _ = NewSet(nextCommitee, config.WeightedRoundRobin, nextLastProposer)
		c := committeeSet.GetProposer(r, big.NewInt(h))
		nextLastProposer = c.Address
		if c.Address == addr1 {
			numOfMaliciousValidators++
		}

		//try to manipulate adding or removing validators
		for j := 5; j > 3; j-- {
			committeeSet, _ = NewSet(vals[:j], config.WeightedRoundRobin, nextLastProposer)
			c = committeeSet.GetProposer(r, big.NewInt(h+1))
			if c.Address == addr1 {
				nextCommitee = vals[:j]
				//fmt.Println("c mal")
				break
			}
		}
		h++
	}

	t.Log(numOfMaliciousValidators)
	if numOfMaliciousValidators > numOfBlocks/4 {
		t.Fatal(fmt.Sprintf("freq: %f, num of blocks: %v, num of expected: %v", float64(numOfMaliciousValidators)/float64(numOfBlocks)*100, numOfMaliciousValidators, numOfBlocks/4))
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
	val4 := types.CommitteeMember{Address: addr4, VotingPower: new(big.Int).SetUint64(100)}
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
