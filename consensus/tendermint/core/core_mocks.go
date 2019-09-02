package core

import (
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
)

type VerifyHeaderAlwaysTrueEngine struct {
	*core
}

func NewVerifyHeaderAlwaysTrueEngine(c consensus.Engine) *VerifyHeaderAlwaysTrueEngine {
	basicCore, ok := c.(*core)
	if !ok {
		panic("*core type is expected")
	}
	return &VerifyHeaderAlwaysTrueEngine{basicCore}
}

func (c *VerifyHeaderAlwaysTrueEngine) VerifyHeader(_ consensus.ChainReader, _ *types.Header, _ bool) error {
	return nil
}
