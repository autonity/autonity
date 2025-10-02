package byzantine

import (
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	bk "github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/rlp"
)

func newOffChainAccusationFuzzer(c interfaces.Core) interfaces.Broadcaster {
	return &OffChainAccusationFuzzer{c.(*core.Core)}
}

// send accusation with garbage accusation msg, the sender of the msg should get disconnected from receiver end.
type OffChainAccusationFuzzer struct {
	*core.Core
}

func fuzzedMessages() []message.Msg {
	num := rand.Intn(10) + 1
	var msgs []message.Msg
	f := fuzz.New()

	for i := 0; i < num; i++ {
		m := &message.Fake{}
		f.Fuzz(&m.FakePayload)
		f.Fuzz(&m.FakeCode)
		// todo: (Jason) add fuzz for the proposal contains in the accusation accountability message.
		// since our RLP encoding of accusation msg does not allow a full proposal, thus we skip the proposal code.
		if m.FakeCode == message.ProposalCode {
			m.FakeCode++
		}
		msgs = append(msgs, m)
	}
	return msgs
}

func (s *OffChainAccusationFuzzer) Broadcast(msg message.Msg) {
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
		peer, ok := backEnd.Broadcaster.FindPeer(c.Address)

		evidences := fuzzedMessages()
		accusation := &accountability.Proof{
			Type:          autonity.AccountabilityEventType(rand.Int()),
			Rule:          autonity.Rule(rand.Int()),
			Evidences:     evidences,
			Message:       evidences[0],
			OffenderIndex: rand.Int(),
		}
		proof, err := rlp.EncodeToBytes(accusation)
		if err != nil {
			panic("Failed to generate random bytes " + err.Error())
		}
		if ok {
			// send fuzzed accusation msg.
			go peer.Send(message.AccountabilityNetworkMsg, proof) // nolint
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
	committee, err := backEnd.BlockChain().CommitteeByHeight(msg.H())
	if err != nil {
		panic(err)
	}

	self := committee.MemberByAddress(s.Address()).Index
	preVotes := backEnd.MsgStore.GetPrevotes(msg.H()-1, func(m *message.Prevote) bool {
		return true
	})

	for _, pv := range preVotes {
		for _, signerIndex := range pv.Signers().FlattenUniq() {
			if uint64(signerIndex) == self {
				continue
			}
			peer, ok := backEnd.Broadcaster.FindPeer(committee.Members[signerIndex].Address)
			if ok {
				accusation := &accountability.Proof{
					Type:          autonity.Accusation,
					Rule:          autonity.PVN,
					Message:       pv,
					OffenderIndex: signerIndex,
				}
				rProof, err := rlp.EncodeToBytes(accusation)
				if err != nil {
					panic("cannot encode accusation at e2e test for off chain accusation protocol")
				}
				// send duplicated msg.
				go peer.Send(message.AccountabilityNetworkMsg, rProof) // nolint
				go peer.Send(message.AccountabilityNetworkMsg, rProof) // nolint
				s.Logger().Info("Off chain Accusation duplicated accusation is simulated")
			}
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
		committee, err := backEnd.BlockChain().CommitteeByHeight(h)
		if err != nil {
			panic(err)
		}

		self := committee.MemberByAddress(s.Address()).Index
		preVotes := backEnd.MsgStore.GetPrevotes(h, func(m *message.Prevote) bool {
			return true
		})
		for _, pv := range preVotes {
			for _, signerIndex := range pv.Signers().FlattenUniq() {
				if uint64(signerIndex) == self {
					continue
				}
				peer, ok := backEnd.Broadcaster.FindPeer(committee.Members[signerIndex].Address)
				if ok {
					accusation := &accountability.Proof{
						Type:          autonity.Accusation,
						Rule:          autonity.PVN,
						Message:       pv,
						OffenderIndex: signerIndex,
					}
					rProof, err := rlp.EncodeToBytes(accusation)
					if err != nil {
						panic("cannot encode accusation at e2e test for off chain accusation protocol")
					}
					// send msg.
					go peer.Send(message.AccountabilityNetworkMsg, rProof) // nolint
					s.Logger().Info("Off chain Accusation over rated accusation is simulated")
				}
			}
		}
	}
}

func TestOffChainAccusation(t *testing.T) {
	// some helper functions
	msgStore := func(node *e2e.Node) *core.MsgStore {
		return node.Eth.Engine().(*bk.Backend).MsgStore
	}
	validatorToSigner := func(v *gengen.Validator) message.Signer {
		return func(h common.Hash) blst.Signature {
			return v.ConsensusKey.Sign(h[:])
		}
	}
	fakeHeader := func() *types.Header {
		return &types.Header{
			ParentHash:         common.Hash{0xde, 0xad},
			UncleHash:          common.Hash{0xbe, 0xef},
			Coinbase:           common.Address{0xca, 0xfe},
			Root:               common.Hash{},
			TxHash:             common.Hash{},
			ReceiptHash:        common.Hash{},
			Bloom:              types.Bloom{},
			Difficulty:         new(big.Int).SetUint64(15),
			Number:             new(big.Int).SetUint64(15),
			GasLimit:           50_000_000,
			GasUsed:            100,
			Time:               0,
			Extra:              nil,
			MixDigest:          types.BFTDigest,
			Nonce:              types.BlockNonce{},
			BaseFee:            nil,
			ProposerSeal:       nil,
			Round:              0,
			ActivityProofRound: 0,
			QuorumCertificate:  nil,
			Epoch:              nil,
			ActivityProof:      nil,
		}
	}
	committeeFromValidators := func(validators []*gengen.Validator) *types.Committee {
		c := &types.Committee{Members: make([]types.CommitteeMember, 0, len(validators))}
		for i, v := range validators {
			c.Members = append(c.Members, types.CommitteeMember{
				Address:           crypto.PubkeyToAddress(v.NodeKey.PublicKey),
				VotingPower:       new(big.Int).SetUint64(v.Stake),
				ConsensusKeyBytes: v.ConsensusKey.PublicKey().Marshal(),
				ConsensusKey:      v.ConsensusKey.PublicKey(),
				Index:             uint64(i),
			})

		}
		return c
	}

	t.Run("off-chain accusation - C1 rule", func(t *testing.T) {
		validators, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)

		c := committeeFromValidators(validators)

		// prepare artificial precommit so that C1 accusation can be triggered at height 15
		accuserIndex := 0
		accusedIndex := 1
		r := int64(99) // high round, so that equivocation does not come into the picture
		h := uint64(15)
		value := common.Hash{0xca, 0xfe}
		accusableVote := message.NewPrecommit(r, h, value, validatorToSigner(validators[accusedIndex]), &c.Members[accusedIndex], c)
		innocenceProof := message.AggregatePrevotesSingle([]message.Vote{
			message.NewPrevote(r, h, value, validatorToSigner(validators[0]), &c.Members[0], c),
			message.NewPrevote(r, h, value, validatorToSigner(validators[1]), &c.Members[1], c),
			message.NewPrevote(r, h, value, validatorToSigner(validators[2]), &c.Members[2], c),
			message.NewPrevote(r, h, value, validatorToSigner(validators[3]), &c.Members[3], c),
		}) // has quorum --> no PVN accusation gets raised

		network, err := e2e.NewNetworkFromValidators(t, validators, true)
		require.NoError(t, err)
		defer network.Shutdown(t)

		// wait for some heights to be mined to correctly set the first buffered height in the msgStore
		err = network.WaitForHeight(2, 30)
		require.NoError(t, err)

		// insert messages to trigger accusations and relative innocence proof
		msgStore(network[accuserIndex]).Save(accusableVote)
		msgStore(network[accusedIndex]).Save(innocenceProof)

		accuserAddress := network[accuserIndex].Address
		accusedAddress := network[accusedIndex].Address

		// if an off-chain accountability accusation is sent
		// the AccountabilityEvent will be posted by the backend handler of the accused
		var receivedOffChainAccusations atomic.Uint64
		var wg sync.WaitGroup
		var sub *event.TypeMuxSubscription
		defer func() {
			sub.Unsubscribe()
			wg.Wait()
		}()
		sub = network[accusedIndex].Eth.Engine().(*bk.Backend).Subscribe(events.AccountabilityEvent{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case ev, ok := <-sub.Chan():
					if !ok {
						return // channel closed
					}
					require.Equal(t, accuserAddress, ev.Data.(events.AccountabilityEvent).Sender)
					receivedOffChainAccusations.Add(1)
				}
			}
		}()

		// network should be up and continue to mine blocks
		err = network.WaitToMineNBlocks(100, 200, false)
		require.NoError(t, err)

		// at least one off-chain accusations should have been received by the validators
		receivedOffChainAccusationsUint64 := receivedOffChainAccusations.Load()
		t.Logf("received off-chain accusations: %d", receivedOffChainAccusationsUint64)
		require.Greater(t, receivedOffChainAccusationsUint64, uint64(0))

		// accusation of C1 should not end up on-chain, it should be resolved off-chain
		err = e2e.AccountabilityEventDetected(t, accusedAddress, autonity.Accusation, autonity.C1, network)
		require.ErrorIs(t, err, e2e.ErrAccountabilityEventMissing)
	})
	t.Run("off-chain accusation - PVN rule", func(t *testing.T) {
		validators, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s,%s,%s")
		require.NoError(t, err)

		c := committeeFromValidators(validators)

		// prepare artificial prevote so that PVN accusation can be triggered at height 15
		accuserIndex := 0
		accusedIndex := 1
		otherProposer := 2 // someone has to propose the new proposal, and will be misbehaving of PN
		r := int64(99)     // high round, so that equivocation does not come into the picture
		h := uint64(15)
		header := fakeHeader()
		block := types.NewBlockWithHeader(header)
		accusableVote := message.NewPrevote(r, h, block.Hash(), validatorToSigner(validators[accusedIndex]), &c.Members[accusedIndex], c)
		innocenceProof := message.NewPropose(r, h, -1, block, validatorToSigner(validators[otherProposer]), &c.Members[otherProposer])

		network, err := e2e.NewNetworkFromValidators(t, validators, true)
		require.NoError(t, err)
		defer network.Shutdown(t)

		// wait for some heights to be mined to correctly set the first buffered height in the msgStore
		err = network.WaitForHeight(2, 30)
		require.NoError(t, err)

		// insert messages to trigger accusations and relative innocence proof
		msgStore(network[accuserIndex]).Save(accusableVote)
		msgStore(network[accusedIndex]).Save(innocenceProof)

		accuserAddress := network[accuserIndex].Address
		accusedAddress := network[accusedIndex].Address

		// if an off-chain accountability accusation is sent
		// the AccountabilityEvent will be posted by the backend handler of the accused
		var receivedOffChainAccusations atomic.Uint64
		var wg sync.WaitGroup
		var sub *event.TypeMuxSubscription
		defer func() {
			sub.Unsubscribe()
			wg.Wait()
		}()
		sub = network[accusedIndex].Eth.Engine().(*bk.Backend).Subscribe(events.AccountabilityEvent{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case ev, ok := <-sub.Chan():
					if !ok {
						return // channel closed
					}
					require.Equal(t, accuserAddress, ev.Data.(events.AccountabilityEvent).Sender)
					receivedOffChainAccusations.Add(1)
				}
			}
		}()

		// network should be up and continue to mine blocks
		err = network.WaitToMineNBlocks(100, 200, false)
		require.NoError(t, err)

		// at least one off-chain accusations should have been received by the accused
		receivedOffChainAccusationsUint64 := receivedOffChainAccusations.Load()
		t.Logf("received off-chain accusations: %d", receivedOffChainAccusationsUint64)
		require.Greater(t, receivedOffChainAccusationsUint64, uint64(0))

		// accusation of PVN should not end up on-chain, it should be resolved off-chain
		err = e2e.AccountabilityEventDetected(t, accusedAddress, autonity.Accusation, autonity.PVN, network)
		require.ErrorIs(t, err, e2e.ErrAccountabilityEventMissing)
	})

	t.Run("Test off chain accusation with fuzzed msg", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: newOffChainAccusationFuzzer}
		runDropPeerConnectionTest(t, handler, 30, 60)
	})

	t.Run("Test duplicated accusation msg from same peer", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: newOffChainDuplicatedAccusationBroadcaster}
		runDropPeerConnectionTest(t, handler, 30, 60)
	})

	t.Run("Test over rated off chain accusation", func(t *testing.T) {
		handler := &interfaces.Services{Broadcaster: newOverRatedOffChainAccusation}
		runDropPeerConnectionTest(t, handler, 30, 60)
	})
}

func runDropPeerConnectionTest(t *testing.T, handler *interfaces.Services, testPeriod uint64, numSec int) { // nolint
	validators, err := e2e.Validators(t, 4, "10e36,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	// set malicious
	spammer := 0
	validators[spammer].TendermintServices = handler

	// creates a network of 4 validators and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, validators, true)
	require.NoError(t, err)
	defer network.Shutdown(t)
	n := network[spammer]
	count := n.ConsensusServer().PeerCount()
	require.Equal(t, 3, count)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(testPeriod, numSec, false)
	require.NoError(t, err)

	// the challenger should get no peer connection left.
	count = n.ConsensusServer().PeerCount()
	require.Equal(t, 0, count)
}
