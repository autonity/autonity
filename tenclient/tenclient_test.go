package tenclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContractAddress(t *testing.T) {

	test := &testCase{
		name:          "no malicious",
		numValidators: 3,
		numBlocks:     5,
		txPerPeer:     1,
	}
	nodes := setupNodes(t, test)
	assert.NotNil(t, nodes)
}
