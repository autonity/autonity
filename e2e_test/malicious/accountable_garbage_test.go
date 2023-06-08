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

type accountableGarbageProposalBroadcaster struct {
	*core.Core
}

func simulateAccountableGarbageMsg(ctx context.Context, c *core.Core, msg *messageutils.Message, code uint8) {
	logger := c.Logger().New("step", c.Step())
	logger.Info("Broadcasting accountable garbage msg bytes")

	randomBytes, err := et.GenerateRandomBytes(2048)
	if err != nil {
		logger.Error("Failed to generate random bytes ", "err", err)
		return
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
		return
	}
	logger.Info("Misbehaviour of AccountableGarbageMessage rule is simulated by garbage maliciousMsg msg.")
	et.DefaultBehaviour(ctx, c, msg)
	_ = c.Backend().Broadcast(ctx, c.CommitteeSet().Committee(), maliciousBytes)

}

func (s *accountableGarbageProposalBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	simulateAccountableGarbageMsg(ctx, s.Core, msg, consensus.MsgProposal)
}

type accountableGarbagePrevoteBroadcaster struct {
	*core.Core
}

func (s *accountableGarbagePrevoteBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	simulateAccountableGarbageMsg(ctx, s.Core, msg, consensus.MsgPrevote)
}

type accountableGarbagePrecommitBroadcaster struct {
	*core.Core
}

func (s *accountableGarbagePrecommitBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	simulateAccountableGarbageMsg(ctx, s.Core, msg, consensus.MsgPrecommit)
}

type invalidMsgCodeBroadcaster struct {
	*core.Core
}

func (s *invalidMsgCodeBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	simulateAccountableGarbageMsg(ctx, s.Core, msg, consensus.MsgPrecommit+100)
}

func TestTBFTAccountableGarbageMsgTests(t *testing.T) {
	t.Run("TestTBFTMisbehaviourRuleAccountableGarbageProposal", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &accountableGarbageProposalBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.AccountableGarbageMessage
		runAccountabilityEventTest(t, handler, tp, rule, 45)
	})

	t.Run("TestTBFTMisbehaviourRuleAccountableGarbagePrevote", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &accountableGarbagePrevoteBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.AccountableGarbageMessage
		runAccountabilityEventTest(t, handler, tp, rule, 45)
	})

	t.Run("TestTBFTMisbehaviourRuleAccountableGarbagePrecommit", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &accountableGarbagePrecommitBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.AccountableGarbageMessage
		runAccountabilityEventTest(t, handler, tp, rule, 45)
	})

	t.Run("TestTBFTMisbehaviourRuleAccountableGarbagePrecommit with invalid consensus msg code", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &invalidMsgCodeBroadcaster{}}
		tp := autonity.Misbehaviour
		rule := autonity.AccountableGarbageMessage
		runAccountabilityEventTest(t, handler, tp, rule, 45)
	})
}
