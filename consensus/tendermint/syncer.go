package tendermint

import (
	time "time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	"github.com/clearmatics/autonity/core/types"
)

type Syncer interface {
	Stop()
	AskSync(lastestHeader *types.Header)
	SyncPeer(peerAddr common.Address, messages [][]byte)
}

type DefaultSyncer struct {
	address common.Address
	peers   consensus.Peers
	stopped chan struct{}
}

func NewSyncer(peers consensus.Peers, address common.Address) *DefaultSyncer {
	return &DefaultSyncer{
		peers:   peers,
		address: address,
		stopped: make(chan struct{}),
	}
}

func (s *DefaultSyncer) AskSync(latestHeader *types.Header) {
	var count uint64

	// Determine if there should be any other peers
	potentialPeerCount := len(latestHeader.Committee)
	if latestHeader.CommitteeMember(s.address) != nil {
		// Remove ourselves from the other potential peers
		potentialPeerCount--
	}
	// Exit if there are no other peers
	if potentialPeerCount == 0 {
		return
	}
	peers := s.peers.Peers()
	// Wait for there to be peers
	for len(peers) == 0 {
		t := time.NewTimer(10 * time.Millisecond)
		select {
		case <-t.C:
			peers = s.peers.Peers()
			continue
		case <-s.stopped:
			return
		}
	}

	// Ask for sync to peers
	for _, p := range peers {
		//ask to a quorum nodes to sync, 1 must then be honest and updated
		if count >= bft.Quorum(latestHeader.TotalVotingPower()) {
			break
		}
		go p.Send(tendermintSyncMsg, []byte{}) //nolint
		member := latestHeader.CommitteeMember(p.Address())
		if member == nil {
			continue
		}
		count += member.VotingPower.Uint64()
	}
}

// Synchronize new connected peer with current height state
func (s *DefaultSyncer) SyncPeer(address common.Address, messages [][]byte) {
	for _, p := range s.peers.Peers() {
		if address == p.Address() {
			for _, msg := range messages {
				go p.Send(tendermintMsg, msg) //nolint
			}
			break
		}
	}
}

func (s *DefaultSyncer) Stop() {
	close(s.stopped)
}
