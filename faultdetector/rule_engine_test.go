package faultdetector

import (
	"github.com/clearmatics/autonity/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuleEngine(t *testing.T) {
	addrNode := common.Address{0x0}
	t.Run("getInnocentProof with unprovable rule id", func(t *testing.T) {
		fd := NewFaultDetector(nil, addrNode)
		var input = Proof{
			Rule: PVO,
		}

		_, err := fd.getInnocentProof(&input)
		assert.NotNil(t, err)
	})
}
