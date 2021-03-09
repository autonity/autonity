package faultdetector

import (
	"github.com/clearmatics/autonity/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuleEngine(t *testing.T) {
	addrNode := common.Address{0x0}
	/*
	height := uint64(100)
	round := int64(3)
	validRound := int64(1)
	heightPower := 100
	*/

	/*
	addrAlice := common.Address{0x1}
	addrBob := common.Address{0x2}
	noneNilValue := common.Hash{0x1}
	*/

	t.Run("getInnocentProof with unprovable rule id", func(t *testing.T) {
		fd := NewFaultDetector(nil, addrNode)
		var input = Proof{
			Rule: PVO,
		}

		_, err := fd.getInnocentProof(&input)
		assert.NotNil(t, err)
	})

	// assume that an accusation of PO is rise for node, node need to get the proof of innocence of PO.
	t.Run("GetInnocentProofOfPO", func(t *testing.T) {
		// PO: node propose an old value with an validRound, innocent proof of it should be:
		// there are quorum num of prevote for that value at the validRound.
		/*
		fd := NewFaultDetector(nil, addrNode)
		// set total power of height
		fd.totalPowers[height] = heightPower
		proposal := newProposalMsg(height, round, validRound, addrNode)
		var accusationProof = Proof{
			Type:     Accusation,
			Rule:     PO,
			Message:  *proposal,
		}
		value := proposal.Value()
		// prepare quorum prevotes for value at the valid round.
		*/

	})

}
