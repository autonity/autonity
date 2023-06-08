package malicious

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	proto "github.com/autonity/autonity/consensus"
	bk "github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/core"
	mUtils "github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/consensus/tendermint/misbehaviourdetector"
	et "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/rlp"
	"github.com/autonity/autonity/test"
	"github.com/stretchr/testify/require"
	"testing"
)

type OffChainAccusationRulePVNBroadcaster struct {
	*core.Core
}

// PVN accusation is simulated by the removal of proposal and those corresponding quorum prevotes from msg store on a
// client, such client will rise accusation PVN over those client who prevote for the removed proposal.
func (s *OffChainAccusationRulePVNBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	et.DefaultBehaviour(ctx, s.Core, msg)
	decodedMsg := et.DecodeMsg(msg, s.Core)
	currentHeight := uint64(15)
	if decodedMsg.H() != currentHeight {
		return
	}

	// simulate accusation over height 7
	height := currentHeight - proto.DeltaBlocks + 2
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate off chain accusation PVN")
	}

	proposals := backEnd.MsgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgProposal
	})

	for _, p := range proposals {
		preVotes := backEnd.MsgStore.Get(height, func(m *mUtils.Message) bool {
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

type OffChainAccusationRuleC1Broadcaster struct {
	*core.Core
}

// C1 accusation is simulated by the removal of those corresponding quorum prevotes from msg store on a
// client, thus, the client will rise accusation C1 over those client who precommit for the corresponding proposal that
// there were no quorum prevotes of it.
func (s *OffChainAccusationRuleC1Broadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	et.DefaultBehaviour(ctx, s.Core, msg)
	decodedMsg := et.DecodeMsg(msg, s.Core)
	currentHeight := uint64(15)
	if decodedMsg.H() != currentHeight {
		return
	}

	// simulate accusation over height 7
	height := currentHeight - proto.DeltaBlocks + 2
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate off chain accusation C1")
	}

	proposals := backEnd.MsgStore.Get(height, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgProposal
	})

	for _, p := range proposals {
		preVotes := backEnd.MsgStore.Get(height, func(m *mUtils.Message) bool {
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

// send accusation with garbage accusation msg, the sender of the msg should get disconnected from receiver end.
type OffChainAccusationGarbageBroadcaster struct {
	*core.Core
}

func (s *OffChainAccusationGarbageBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	et.DefaultBehaviour(ctx, s.Core, msg)
	// construct garbage off chain accusation msg and sent it.
	decodedMsg := et.DecodeMsg(msg, s.Core)
	if decodedMsg.H() <= uint64(1) {
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
		payload, err := et.GenerateRandomBytes(2048)
		if err != nil {
			panic("Failed to generate random bytes ")
		}
		if len(peers) > 0 {
			// send garbage accusation msg.
			go peers[c.Address].Send(bk.TendermintOffChainAccountabilityMsg, payload) // nolint
			s.Logger().Info("Off chain Accusation garbage accusation is simulated")
		}
	}
}

// send duplicated accusation msg from challenger, this would get the challenger removed from the peer connection.
type OffChainDuplicatedAccusationBroadcaster struct {
	*core.Core
}

func (s *OffChainDuplicatedAccusationBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	et.DefaultBehaviour(ctx, s.Core, msg)
	// construct duplicated accusation msg and send them.
	decodedMsg := et.DecodeMsg(msg, s.Core)
	if decodedMsg.H() <= uint64(1) {
		return
	}
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate duplicated off chain accusation")
	}

	preVotes := backEnd.MsgStore.Get(decodedMsg.H()-1, func(m *mUtils.Message) bool {
		return m.Type() == proto.MsgPrevote && m.Address != msg.Address
	})

	for _, pv := range preVotes {
		targets := make(map[common.Address]struct{})
		targets[pv.Address] = struct{}{}
		peers := backEnd.Broadcaster.FindPeers(targets)
		if len(peers) > 0 {
			accusation := &misbehaviourdetector.AccountabilityProof{
				Type:    autonity.Accusation,
				Rule:    autonity.PVN,
				Message: pv,
			}
			rProof, err := rlp.EncodeToBytes(accusation)
			if err != nil {
				panic("cannot encode accusation at e2e test for off chain accusation protocol")
			}
			// send duplicated msg.
			go peers[pv.Address].Send(bk.TendermintOffChainAccountabilityMsg, rProof) // nolint
			go peers[pv.Address].Send(bk.TendermintOffChainAccountabilityMsg, rProof) // nolint
			s.Logger().Info("Off chain Accusation duplicated accusation is simulated")
		}
	}
}

type OffChainAccusationOverRatedBroadcaster struct {
	*core.Core
}

func (s *OffChainAccusationOverRatedBroadcaster) Broadcast(ctx context.Context, msg *mUtils.Message) {
	et.DefaultBehaviour(ctx, s.Core, msg)
	// construct accusations and send them with high rate.
	decodedMsg := et.DecodeMsg(msg, s.Core)
	if decodedMsg.H() <= uint64(10) {
		return
	}
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate duplicated off chain accusation")
	}

	// collect some out of updated consensus msg
	for h := uint64(2); h <= decodedMsg.H(); h++ {
		preVotes := backEnd.MsgStore.Get(h, func(m *mUtils.Message) bool {
			return m.Type() == proto.MsgPrevote && m.Address != msg.Address
		})
		for _, pv := range preVotes {
			targets := make(map[common.Address]struct{})
			targets[pv.Address] = struct{}{}
			peers := backEnd.Broadcaster.FindPeers(targets)
			if len(peers) > 0 {
				accusation := &misbehaviourdetector.AccountabilityProof{
					Type:    autonity.Accusation,
					Rule:    autonity.PVN,
					Message: pv,
				}
				rProof, err := rlp.EncodeToBytes(accusation)
				if err != nil {
					panic("cannot encode accusation at e2e test for off chain accusation protocol")
				}
				// send msg.
				go peers[pv.Address].Send(bk.TendermintOffChainAccountabilityMsg, rProof) // nolint
				s.Logger().Info("Off chain Accusation over rated accusation is simulated")
			}
		}
	}
}

func TestTBFTOffChainAccusationTests(t *testing.T) {
	t.Run("TestTBFTOffChainAccusationRuleC1", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &OffChainAccusationRuleC1Broadcaster{}}
		tp := autonity.Accusation
		rule := autonity.C1
		runOffChainAccountabilityEventTest(t, handler, tp, rule, 100)
	})

	t.Run("TestTBFTOffChainAccusationRulePVN", func(t *testing.T) {
		handler := &node.CustomHandler{Broadcaster: &OffChainAccusationRulePVNBroadcaster{}}
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
		handler := &node.CustomHandler{Broadcaster: &OffChainAccusationGarbageBroadcaster{}}
		runDropPeerConnectionTest(t, handler, 10)
	})

	t.Run("Test duplicated accusation msg from same peer", func(t *testing.T) {
		t.Skip("dropped peer was reconnected after a while in the p2p layer causing this case unstable")
		handler := &node.CustomHandler{Broadcaster: &OffChainDuplicatedAccusationBroadcaster{}}
		runDropPeerConnectionTest(t, handler, 15)
	})

	t.Run("Test over rated off chain accusation", func(t *testing.T) {
		t.Skip("dropped peer was reconnected after a while in the p2p layer causing this case unstable")
		handler := &node.CustomHandler{Broadcaster: &OffChainAccusationOverRatedBroadcaster{}}
		runDropPeerConnectionTest(t, handler, 20)
	})
}

func runDropPeerConnectionTest(t *testing.T, handler *node.CustomHandler, testPeriod uint64) { // nolint
	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	users, err := test.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	// set malicious challenger
	challenger := 0
	users[challenger].CustHandler = handler
	// creates a network of 4 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(testPeriod, 20) // nolint

	// the challenger should get no peer connection left.
	n := network[1]
	client := n.WsClient
	count, err := client.PeerCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint(2), count)
}

func runOffChainAccountabilityEventTest(t *testing.T, handler *node.CustomHandler, tp autonity.AccountabilityEventType,
	rule autonity.Rule, testPeriod uint64) {

	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	users, err := test.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	// set malicious challenger
	challenger := 0
	users[challenger].CustHandler = handler
	// creates a network of 4 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	network.WaitToMineNBlocks(testPeriod, 500) // nolint

	// accusation of PVN shouldn't be submitted on chain by challenger.
	challengerAddress := network[challenger].Address
	for _, n := range network {
		if n.Address == challengerAddress {
			continue
		}
		detected := et.AccountabilityEventDetected(t, n.Address, tp, rule, network)
		require.Equal(t, false, detected)
	}
}
