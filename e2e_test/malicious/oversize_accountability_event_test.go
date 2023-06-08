package malicious

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	et "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"testing"
)

type accountableOversizeGarbageProposalBroadcaster struct {
	*core.Core
	faultSimulated bool
}

func simulateAccountableOversizeGarbagePropoosalMsg(ctx context.Context, c *core.Core, msg *messageutils.Message, code uint8) bool {

	logger := c.Logger().New("step", c.Step())
	logger.Info("Broadcasting accountable oversize garbage msg bytes")

	// simulate an oversize accountability event.
	randomBytes, err := et.GenerateRandomBytes(1024 * 512)
	if err != nil {
		logger.Error("Failed to generate random bytes ", "err", err)
		return false
	}

	maliciousMsg := &messageutils.Message{
		Code:          code,
		TbftMsgBytes:  randomBytes,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	}
	maliciousBytes, err := c.FinalizeMessage(maliciousMsg)
	if err != nil {
		logger.Error("Failed to finalize accountable garbage maliciousMsg msg", "err", err)
		return false
	}
	logger.Info("Misbehaviour of AccountableGarbageMessage rule is simulated by oversize garbage msg.")
	et.DefaultBehaviour(ctx, c, msg)
	_ = c.Backend().Broadcast(ctx, c.CommitteeSet().Committee(), maliciousBytes)
	return true
}

func (s *accountableOversizeGarbageProposalBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	if s.faultSimulated {
		et.DefaultBehaviour(ctx, s.Core, msg)
		return
	}
	s.faultSimulated = simulateAccountableOversizeGarbagePropoosalMsg(ctx, s.Core, msg, consensus.MsgProposal)
}

func TestAccountableOversizeGarbageMsgTests(t *testing.T) {
	t.Run("TestTBFTMisbehaviourRuleAccountableOversizeGarbageProposal", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &accountableOversizeGarbageProposalBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.AccountableGarbageMessage
		runAccountabilityEventTest(t, handler, tp, rule, 120)
	})
}
