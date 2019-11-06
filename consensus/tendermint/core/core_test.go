package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/log"
	"math/big"
	"testing"
)

func TestCore_MeasureHeightRoundMetrics(t *testing.T) {
	t.Run("measure metrics of new height", func(t *testing.T) {
		c := &core{
			address:           common.Address{},
			logger:            log.New("core", "test", "id", 0),
			proposeTimeout:    newTimeout(propose),
			prevoteTimeout:    newTimeout(prevote),
			precommitTimeout:  newTimeout(precommit),
			currentRoundState: NewRoundState(big.NewInt(0), big.NewInt(1)),
		}
		c.measureHeightRoundMetrics(common.Big0)
	})

	t.Run("measure metrics of new round", func(t *testing.T) {
		c := &core{
			address:           common.Address{},
			logger:            log.New("core", "test", "id", 0),
			proposeTimeout:    newTimeout(propose),
			prevoteTimeout:    newTimeout(prevote),
			precommitTimeout:  newTimeout(precommit),
			currentRoundState: NewRoundState(big.NewInt(0), big.NewInt(1)),
		}
		c.measureHeightRoundMetrics(common.Big1)
	})
}
