package economic

import (
	"testing"
)

func TestFullFilledBlocks(t *testing.T) {
	blocks := uint64(600)
	copyParams := defSysParams
	copyParams.baseFeeChangeDenominator = 8
	sim := newSimulator(blocks, &copyParams, &fullFiller{params: &copyParams})
	sim.start()
}

func TestHalfFilledBlocks(t *testing.T) {
	blocks := uint64(600)
	copyParam := defSysParams
	systemParam := &copyParam
	sim := newSimulator(blocks, systemParam, &halfFiller{params: systemParam})
	sim.start()
}

func TestDustFilledBlocks(t *testing.T) {
	blocks := uint64(600)
	systemParam := &defSysParams
	sim := newSimulator(blocks, systemParam, &dustFiller{params: systemParam})
	sim.start()
}

// todo, address the cost of spam.

// todo, add helpers to render data in a diagram.
