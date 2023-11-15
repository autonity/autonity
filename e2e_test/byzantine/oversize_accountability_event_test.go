package byzantine

import (
	"context"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
)

type OversizedEventBroadcaster struct {
	*core.Core
	faultSimulated bool
}

func NewOversizedEventBroadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &OversizedEventBroadcaster{c.(*core.Core), false}
}

func (s *OversizedEventBroadcaster) SignAndBroadcast(msg *message.Message) {
	e2e.DefaultSignAndBroadcast(s.Core, msg)
	if s.faultSimulated {
		return
	}
	logger := s.Logger().New("step", s.Step())
	logger.Info("Broadcasting accountable oversize garbage msg bytes")

	// simulate an oversize accountability event.
	randomBytes, err := e2e.GenerateRandomBytes(1024 * 9)
	if err != nil {
		logger.Crit("Failed to generate random bytes ", "err", err)
		return
	}
	maliciousMsg := &message.Message{
		Code:          1,
		Payload:       randomBytes,
		Address:       s.Address(),
		CommittedSeal: []byte{},
	}

	logger.Info("Misbehaviour of AccountableGarbageMessage rule is simulated by oversize garbage msg.")
	e2e.DefaultSignAndBroadcast(s.Core, maliciousMsg)
	s.faultSimulated = true
}

func TestAccountableOversizeGarbageMsgTests(t *testing.T) {
	t.Run("TestTBFTMisbehaviourRuleAccountableOversizeGarbageProposal", func(t *testing.T) {
		tp := autonity.Misbehaviour
		rule := autonity.GarbageMessage
		handler := &node.TendermintServices{Broadcaster: NewOversizedEventBroadcaster}
		runTest(t, handler, tp, rule, 120)
	})
}
