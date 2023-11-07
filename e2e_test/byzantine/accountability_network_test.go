package byzantine

import (
	"context"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	proto "github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	bk "github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/require"
)

func newOffChainAccusationRulePVNBroadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &OffChainAccusationRulePVNBroadcaster{c.(*core.Core)}
}

type OffChainAccusationRulePVNBroadcaster struct {
	*core.Core
}

// PVN accusation is simulated by the removal of proposal and those corresponding quorum prevotes from msg store on a
// client, such client will rise accusation PVN over those client who prevote for the removed proposal.
func (s *OffChainAccusationRulePVNBroadcaster) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	_ = msg.DecodePayload()
	currentHeight := uint64(15)
	if msg.H() != currentHeight {
		return
	}

	// simulate accusation over height 7
	height := currentHeight - proto.DeltaBlocks + 2
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate off chain accusation PVN")
	}

	proposals := backEnd.MsgStore.Get(height, func(m *message.Message) bool {
		return m.Type() == proto.MsgProposal
	})

	for _, p := range proposals {
		preVotes := backEnd.MsgStore.Get(height, func(m *message.Message) bool {
			return m.Type() == proto.MsgPrevote && m.R() == p.R() && m.Value() == p.Value()
		})
		// remove proposal.
		backEnd.MsgStore.RemoveMsg(p.H(), p.R(), p.Code, p.Address)
		// remove over quorum corresponding prevotes.
		counter := 0
		for _, pv := range preVotes {
			if counter < len(preVotes)/2 {
				backEnd.MsgStore.RemoveMsg(pv.H(), pv.R(), pv.Code, pv.Address)
				counter++
			} else {
				break
			}
		}
	}
	s.Logger().Info("Off chain Accusation of PVN rule is simulated")
}

func newOffChainAccusationRuleC1Broadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &OffChainAccusationRuleC1Broadcaster{c.(*core.Core)}
}

type OffChainAccusationRuleC1Broadcaster struct {
	*core.Core
}

// C1 accusation is simulated by the removal of those corresponding quorum prevotes from msg store on a
// client, thus, the client will rise accusation C1 over those client who precommit for the corresponding proposal that
// there were no quorum prevotes of it.
func (s *OffChainAccusationRuleC1Broadcaster) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	_ = msg.DecodePayload()
	currentHeight := uint64(15)
	if msg.H() != currentHeight {
		return
	}

	// simulate accusation over height 7
	height := currentHeight - proto.DeltaBlocks + 2
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate off chain accusation C1")
	}

	proposals := backEnd.MsgStore.Get(height, func(m *message.Message) bool {
		return m.Type() == proto.MsgProposal
	})

	for _, p := range proposals {
		preVotes := backEnd.MsgStore.Get(height, func(m *message.Message) bool {
			return m.Type() == proto.MsgPrevote && m.R() == p.R() && m.Value() == p.Value()
		})

		// remove over quorum corresponding prevotes.
		counter := 0
		for _, pv := range preVotes {
			if counter < len(preVotes)/2 {
				backEnd.MsgStore.RemoveMsg(pv.H(), pv.R(), pv.Code, pv.Address)
				counter++
			} else {
				break
			}
		}
	}
	s.Logger().Info("Off chain Accusation of C1 rule is simulated")
}

func newOffChainAccusationGarbageBroadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &OffChainAccusationGarbageBroadcaster{c.(*core.Core)}
}

// send accusation with garbage accusation msg, the sender of the msg should get disconnected from receiver end.
type OffChainAccusationGarbageBroadcaster struct {
	*core.Core
}

func (s *OffChainAccusationGarbageBroadcaster) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	// construct garbage off chain accusation msg and sent it.
	_ = msg.DecodePayload()
	if msg.H() <= uint64(1) {
		return
	}
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate duplicated off chain accusation")
	}

	committee := s.Core.Committee().Committee()

	for _, c := range committee {
		if c.Address == s.Address() {
			continue
		}
		targets := make(map[common.Address]struct{})
		targets[c.Address] = struct{}{}
		peers := backEnd.Broadcaster.FindPeers(targets)
		payload, err := e2e.GenerateRandomBytes(2048)
		if err != nil {
			panic("Failed to generate random bytes ")
		}
		if len(peers) > 0 {
			// send garbage accusation msg.
			go peers[c.Address].Send(bk.AccountabilityMsg, payload) // nolint
			s.Logger().Info("Off chain Accusation garbage accusation is simulated")
		}
	}
}

func newOffChainDuplicatedAccusationBroadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &OffChainDuplicatedAccusationBroadcaster{c.(*core.Core)}
}

// send duplicated accusation msg from challenger, this would get the challenger removed from the peer connection.
type OffChainDuplicatedAccusationBroadcaster struct {
	*core.Core
}

func (s *OffChainDuplicatedAccusationBroadcaster) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	// construct duplicated accusation msg and send them.
	_ = msg.DecodePayload()
	if msg.H() <= uint64(1) {
		return
	}
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate duplicated off chain accusation")
	}

	preVotes := backEnd.MsgStore.Get(msg.H()-1, func(m *message.Message) bool {
		return m.Type() == proto.MsgPrevote && m.Address != msg.Address
	})

	for _, pv := range preVotes {
		targets := make(map[common.Address]struct{})
		targets[pv.Address] = struct{}{}
		peers := backEnd.Broadcaster.FindPeers(targets)
		if len(peers) > 0 {
			accusation := &accountability.Proof{
				Type:    autonity.Accusation,
				Rule:    autonity.PVN,
				Message: pv,
			}
			rProof, err := rlp.EncodeToBytes(accusation)
			if err != nil {
				panic("cannot encode accusation at e2e test for off chain accusation protocol")
			}
			// send duplicated msg.
			go peers[pv.Address].Send(bk.AccountabilityMsg, rProof) // nolint
			go peers[pv.Address].Send(bk.AccountabilityMsg, rProof) // nolint
			s.Logger().Info("Off chain Accusation duplicated accusation is simulated")
		}
	}
}

func newOffChainAccusationOverRatedBroadcaster(c interfaces.Tendermint) interfaces.Broadcaster {
	return &OffChainAccusationOverRatedBroadcaster{c.(*core.Core)}
}

type OffChainAccusationOverRatedBroadcaster struct {
	*core.Core
}

func (s *OffChainAccusationOverRatedBroadcaster) SignAndBroadcast(ctx context.Context, msg *message.Message) {
	e2e.DefaultSignAndBroadcast(ctx, s.Core, msg)
	// construct accusations and send them with high rate.
	_ = msg.DecodePayload()
	if msg.H() <= uint64(10) {
		return
	}
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate duplicated off chain accusation")
	}

	// collect some out of updated consensus msg
	for h := uint64(2); h <= msg.H(); h++ {
		preVotes := backEnd.MsgStore.Get(h, func(m *message.Message) bool {
			return m.Type() == proto.MsgPrevote && m.Address != msg.Address
		})
		for _, pv := range preVotes {
			targets := make(map[common.Address]struct{})
			targets[pv.Address] = struct{}{}
			peers := backEnd.Broadcaster.FindPeers(targets)
			if len(peers) > 0 {
				accusation := &accountability.Proof{
					Type:    autonity.Accusation,
					Rule:    autonity.PVN,
					Message: pv,
				}
				rProof, err := rlp.EncodeToBytes(accusation)
				if err != nil {
					panic("cannot encode accusation at e2e test for off chain accusation protocol")
				}
				// send msg.
				go peers[pv.Address].Send(bk.AccountabilityMsg, rProof) // nolint
				s.Logger().Info("Off chain Accusation over rated accusation is simulated")
			}
		}
	}
}

func OffChainAccusationTests(t *testing.T) {
	t.Run("OffChainAccusationRuleC1", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: newOffChainAccusationRuleC1Broadcaster}
		tp := autonity.Accusation
		rule := autonity.C1
		runOffChainAccountabilityEventTest(t, handler, tp, rule, 100)
	})

	t.Run("OffChainAccusationRulePVN", func(t *testing.T) {
		handler := &node.TendermintServices{Broadcaster: newOffChainAccusationRulePVNBroadcaster}
		tp := autonity.Accusation
		rule := autonity.PVN
		runOffChainAccountabilityEventTest(t, handler, tp, rule, 100)
	})

	// Following in belows, there are 3 testcases to observe if those malicious off-chain challenger's peer is dropped
	// due to the DoS protection of off chain accusation protocol. Due to the re-dail scheduler in P2P layer, those
	// dropped peer are reconnected after a short while which making the tests unstable from e2e point of view. So skip
	// them from CI job.
	t.Run("Test off chain accusation with garbage msg", func(t *testing.T) {
		t.Skip("dropped peer was reconnected after a while in the p2p layer causing this case unstable")
		handler := &node.TendermintServices{Broadcaster: newOffChainAccusationGarbageBroadcaster}
		runDropPeerConnectionTest(t, handler, 10)
	})

	t.Run("Test duplicated accusation msg from same peer", func(t *testing.T) {
		t.Skip("dropped peer was reconnected after a while in the p2p layer causing this case unstable")
		handler := &node.TendermintServices{Broadcaster: newOffChainDuplicatedAccusationBroadcaster}
		runDropPeerConnectionTest(t, handler, 15)
	})

	t.Run("Test over rated off chain accusation", func(t *testing.T) {
		t.Skip("dropped peer was reconnected after a while in the p2p layer causing this case unstable")
		handler := &node.TendermintServices{Broadcaster: newOffChainAccusationOverRatedBroadcaster}
		runDropPeerConnectionTest(t, handler, 20)
	})
}

func runDropPeerConnectionTest(t *testing.T, handler *node.TendermintServices, testPeriod uint64) { // nolint
	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	validators, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	// set malicious challenger
	challenger := 0
	validators[challenger].TendermintServices = handler
	// creates a network of 4 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(testPeriod, 20, false) // nolint

	// the challenger should get no peer connection left.
	n := network[1]
	client := n.WsClient
	count, err := client.PeerCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint(2), count)
}

func runOffChainAccountabilityEventTest(t *testing.T, handler *node.TendermintServices, tp autonity.AccountabilityEventType,
	rule autonity.Rule, testPeriod uint64) {

	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	users, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	// set malicious challenger
	challenger := 0
	users[challenger].TendermintServices = handler
	// creates a network of 4 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(testPeriod, 500, false) // nolint

	// accusation of PVN shouldn't be submitted on chain by challenger.
	challengerAddress := network[challenger].Address
	for _, n := range network {
		if n.Address == challengerAddress {
			continue
		}
		detected := e2e.AccountabilityEventDetected(t, n.Address, tp, rule, network)
		require.Equal(t, false, detected)
	}
}
