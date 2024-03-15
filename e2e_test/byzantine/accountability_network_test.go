package byzantine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	bk "github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/rlp"
)

func newPVNOffChainAccusation(c interfaces.Core) interfaces.Broadcaster {
	return &PVNOffChainAccusation{c.(*core.Core)}
}

type PVNOffChainAccusation struct {
	*core.Core
}

// PVN accusation is simulated by the removal of proposal and those corresponding quorum prevotes from msg store on a
// client, such client will rise accusation PVN over those client who prevote for the removed proposal.
func (s *PVNOffChainAccusation) Broadcast(msg message.Msg) {
	//TODO(lorenzo) fix this test. PVN accusation will not be raised anymore because the block has been mined
	s.BroadcastAll(msg)
	currentHeight := uint64(15)
	if msg.H() != currentHeight {
		return
	}

	// simulate accusation over height 13 (will be scanned at height 23)
	height := currentHeight - accountability.DeltaBlocks + 8
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate off chain accusation PVN")
	}

	proposals := backEnd.MsgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.ProposalCode
	})

	for _, p := range proposals {
		preVotes := backEnd.MsgStore.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.R() == p.R() && m.Value() == p.Value()
		})
		// remove proposal.
		backEnd.MsgStore.RemoveMsg(p.H(), p.R(), p.Code(), p.Sender())
		// remove over quorum corresponding prevotes.
		counter := 0
		for _, pv := range preVotes {
			if counter < len(preVotes)/2 {
				backEnd.MsgStore.RemoveMsg(pv.H(), pv.R(), pv.Code(), pv.Sender())
				counter++
			} else {
				break
			}
		}
	}
	s.Logger().Info("MsgStore manipulated to cause accusation of PVN rule to be raised later on", "accusationHeight", height)
}

func newC1OffChainAccusation(c interfaces.Core) interfaces.Broadcaster {
	return &C1OffChainAccusation{c.(*core.Core)}
}

type C1OffChainAccusation struct {
	*core.Core
}

// C1 accusation is simulated by the removal of those corresponding quorum prevotes from msg store on a
// client, thus, the client will rise accusation C1 over those client who precommit for the corresponding proposal that
// there were no quorum prevotes of it.
func (s *C1OffChainAccusation) Broadcast(msg message.Msg) {
	//TODO(lorenzo) fix this test. C1 accusation will not be raised anymore because the block has been mined
	s.BroadcastAll(msg)
	currentHeight := uint64(15)
	if msg.H() != currentHeight {
		return
	}

	// simulate accusation over height 13 (will be scanned at height 23)
	height := currentHeight - accountability.DeltaBlocks + 8

	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate off chain accusation C1")
	}

	proposals := backEnd.MsgStore.Get(height, func(m message.Msg) bool {
		return m.Code() == message.ProposalCode
	})

	for _, p := range proposals {
		preVotes := backEnd.MsgStore.Get(height, func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.R() == p.R() && m.Value() == p.Value()
		})

		// remove over quorum corresponding prevotes.
		counter := 0
		for _, pv := range preVotes {
			if counter < len(preVotes)/2 {
				backEnd.MsgStore.RemoveMsg(pv.H(), pv.R(), pv.Code(), pv.Sender())
				counter++
			} else {
				break
			}
		}
	}
	s.Logger().Info("MsgStore manipulated to cause accusation of C1 rule to be raised later on", "accusationHeight", height)
}

func newGarbageOffChainAccusation(c interfaces.Core) interfaces.Broadcaster {
	return &GarbageOffChainAccusation{c.(*core.Core)}
}

// send accusation with garbage accusation msg, the sender of the msg should get disconnected from receiver end.
type GarbageOffChainAccusation struct {
	*core.Core
}

func (s *GarbageOffChainAccusation) Broadcast(msg message.Msg) {
	s.BroadcastAll(msg)
	// construct garbage off chain accusation msg and sent it.
	if msg.H() <= uint64(1) {
		return
	}
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate duplicated off chain accusation")
	}

	committee := s.Core.CommitteeSet().Committee()

	for _, c := range committee.Members {
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
			go peers[c.Address].Send(bk.AccountabilityNetworkMsg, payload) // nolint
			s.Logger().Info("Off chain Accusation garbage accusation is simulated")
		}
	}
}

func newOffChainDuplicatedAccusationBroadcaster(c interfaces.Core) interfaces.Broadcaster {
	return &OffChainDuplicatedAccusationBroadcaster{c.(*core.Core)}
}

// send duplicated accusation msg from challenger, this would get the challenger removed from the peer connection.
type OffChainDuplicatedAccusationBroadcaster struct {
	*core.Core
}

func (s *OffChainDuplicatedAccusationBroadcaster) Broadcast(msg message.Msg) {
	s.BroadcastAll(msg)
	// construct duplicated accusation msg and send them.
	if msg.H() <= uint64(1) {
		return
	}
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate duplicated off chain accusation")
	}

	preVotes := backEnd.MsgStore.Get(msg.H()-1, func(m message.Msg) bool {
		return m.Code() == message.PrevoteCode && m.Sender() != msg.Sender()
	})

	for _, pv := range preVotes {
		targets := make(map[common.Address]struct{})
		targets[pv.Sender()] = struct{}{}
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
			go peers[pv.Sender()].Send(bk.AccountabilityNetworkMsg, rProof) // nolint
			go peers[pv.Sender()].Send(bk.AccountabilityNetworkMsg, rProof) // nolint
			s.Logger().Info("Off chain Accusation duplicated accusation is simulated")
		}
	}
}

func newOverRatedOffChainAccusation(c interfaces.Core) interfaces.Broadcaster {
	return &OverRatedOffChainAccusation{c.(*core.Core)}
}

type OverRatedOffChainAccusation struct {
	*core.Core
}

func (s *OverRatedOffChainAccusation) Broadcast(msg message.Msg) {
	s.BroadcastAll(msg)
	// construct accusations and send them with high rate.
	if msg.H() <= uint64(10) {
		return
	}
	backEnd, ok := s.Core.Backend().(*bk.Backend)
	if !ok {
		panic("cannot simulate duplicated off chain accusation")
	}

	// collect some out of updated consensus msg
	for h := uint64(2); h <= msg.H(); h++ {
		preVotes := backEnd.MsgStore.Get(h, func(m message.Msg) bool {
			return m.Code() == message.PrevoteCode && m.Sender() != msg.Sender()
		})
		for _, pv := range preVotes {
			targets := make(map[common.Address]struct{})
			targets[pv.Sender()] = struct{}{}
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
				go peers[pv.Sender()].Send(bk.AccountabilityNetworkMsg, rProof) // nolint
				s.Logger().Info("Off chain Accusation over rated accusation is simulated")
			}
		}
	}
}

// TODO(lorenzo): add test to check the maximum accusations per height
func TestOffChainAccusation(t *testing.T) {
	// TODO(lorenzo) we should add a check that an offchain accountability message is actually sent
	// if something prevents this from happening, these tests will keep passing because no proof is ever raised
	t.Run("OffChainAccusationRuleC1", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: newC1OffChainAccusation}
		tp := autonity.Accusation
		rule := autonity.C1
		runOffChainAccountabilityEventTest(t, handler, tp, rule, 100)
	})

	t.Run("OffChainAccusationRulePVN", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: newPVNOffChainAccusation}
		tp := autonity.Accusation
		rule := autonity.PVN
		runOffChainAccountabilityEventTest(t, handler, tp, rule, 100)
	})

	// TODO(lorenzo) attempt to restore
	// Following in belows, there are 3 testcases to observe if those malicious off-chain challenger's peer is dropped
	// due to the DoS protection of off chain accusation protocol. Due to the re-dail scheduler in P2P layer, those
	// dropped peer are reconnected after a short while which making the tests unstable from e2e point of view. So skip
	// them from CI job.
	t.Run("Test off chain accusation with garbage msg", func(t *testing.T) {
		t.Skip("dropped peer was reconnected after a while in the p2p layer causing this case unstable")
		handler := &interfaces.Services{Broadcaster: newGarbageOffChainAccusation}
		runDropPeerConnectionTest(t, handler, 10)
	})

	t.Run("Test duplicated accusation msg from same peer", func(t *testing.T) {
		t.Skip("dropped peer was reconnected after a while in the p2p layer causing this case unstable")
		handler := &interfaces.Services{Broadcaster: newOffChainDuplicatedAccusationBroadcaster}
		runDropPeerConnectionTest(t, handler, 15)
	})

	t.Run("Test over rated off chain accusation", func(t *testing.T) {
		t.Skip("dropped peer was reconnected after a while in the p2p layer causing this case unstable")
		handler := &interfaces.Services{Broadcaster: newOverRatedOffChainAccusation}
		runDropPeerConnectionTest(t, handler, 20)
	})
}

func runDropPeerConnectionTest(t *testing.T, handler *interfaces.Services, testPeriod uint64) { // nolint
	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	validators, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// set malicious challenger
	challenger := 0
	validators[challenger].TendermintServices = handler
	// creates a network of 4 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(testPeriod, 20, false)
	require.NoError(t, err)

	// the challenger should get no peer connection left.
	n := network[1]
	client := n.WsClient
	count, err := client.PeerCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint(2), count)
}

func runOffChainAccountabilityEventTest(t *testing.T, handler *interfaces.Services, tp autonity.AccountabilityEventType,
	rule autonity.Rule, testPeriod uint64) {

	//log.Root().SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	users, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// set malicious challenger
	challenger := 0
	users[challenger].TendermintServices = handler
	// creates a network of 4 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(testPeriod, 500, false)
	require.NoError(t, err)

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
