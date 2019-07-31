package core

import "github.com/clearmatics/autonity/consensus"

type QuorumAlwaysFalseEngine struct {
	consensus.Engine
}

func NewCoreQuorumAlwaysFalse(c consensus.Engine) *QuorumAlwaysFalseEngine {
	return &QuorumAlwaysFalseEngine{c}
}

func (c *QuorumAlwaysFalseEngine) Quorum(i int) bool {
	return false
}
