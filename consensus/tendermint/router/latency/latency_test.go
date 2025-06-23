package latency

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/router/mocks"
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
)

func TestFetcher_Fetch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pinger := mocks.NewMockPinger(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	fetcher := NewFetcher(pinger)
	fetcher.SetBroadcaster(peerFinder)

	pubkey1, _ := crypto.HexToECDSA("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	pubkey2, _ := crypto.HexToECDSA("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	enode1 := enode.NewV4(&pubkey1.PublicKey, net.ParseIP("192.168.1.1"), 30303, 0)
	enode2 := enode.NewV4(&pubkey2.PublicKey, net.ParseIP("192.168.1.2"), 30304, 0)

	self := common.HexToAddress("0x111")
	validator1 := crypto.PubkeyToAddress(pubkey1.PublicKey)
	validator2 := crypto.PubkeyToAddress(pubkey2.PublicKey)
	validators := []common.Address{self, validator1, validator2}

	peerFinder.EXPECT().CommitteeEnodes().Return([]*enode.Node{enode1, enode2}).Times(1)
	peerFinder.EXPECT().FindPeer(validator1).Return(consensus.NewMockPeer(ctrl), true).Times(1)
	peerFinder.EXPECT().FindPeer(validator2).Return(nil, false).Times(1)

	pinger.EXPECT().Ping(gomock.Any(), ping.Target{IP: "192.168.1.1", Port: 30303}).
		Return(ping.Result{Latency: 50 * time.Millisecond}).Times(1)
	pinger.EXPECT().Ping(gomock.Any(), ping.Target{IP: "192.168.1.2", Port: 30304}).
		Return(ping.Result{Err: errors.New("ping failed")}).Times(0) // Not called due to FindPeer failure

	latency, failedNodes, err := fetcher.Fetch(validators, self)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, map[common.Address]uint{
		validator1: 50,
	}, latency, "Latency map should match expected values")
	assert.Equal(t, []common.Address{validator2}, failedNodes, "Failed nodes should include validator2")
}

func TestFetcher_FetchNilPeerFinder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pinger := mocks.NewMockPinger(ctrl)
	fetcher := NewFetcher(pinger)

	_, _, err := fetcher.Fetch([]common.Address{common.HexToAddress("0x111")}, common.HexToAddress("0x222"))
	assert.Error(t, err, "Expected error when peerFinder is nil")
	assert.Contains(t, err.Error(), "broadcaster not set", "Error message should mention broadcaster")
}

func TestFetcher_SetBroadcaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pinger := mocks.NewMockPinger(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	fetcher := NewFetcher(pinger)

	fetcher.SetBroadcaster(peerFinder)
	assert.Equal(t, peerFinder, fetcher.peerFinder, "PeerFinder should be set")
}

func TestFetcher_PingPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pinger := mocks.NewMockPinger(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	fetcher := NewFetcher(pinger)
	fetcher.SetBroadcaster(peerFinder)

	targets := []ping.Target{
		{IP: "192.168.1.1", Port: 30303},
		{IP: "", Port: 0},
		{IP: "192.168.1.2", Port: 30304},
	}

	pinger.EXPECT().Ping(gomock.Any(), targets[0]).
		Return(ping.Result{Latency: 100 * time.Millisecond}).Times(1)
	pinger.EXPECT().Ping(gomock.Any(), targets[2]).
		Return(ping.Result{Latency: 200 * time.Millisecond}).Times(1)

	results := fetcher.pingPeers(context.Background(), targets)

	assert.Equal(t, 3, len(results), "Expected results for all targets")
	assert.Equal(t, 100*time.Millisecond, results[0].Latency, "Expected latency for target 0")
	assert.Error(t, results[1].Err, "Expected error for empty target")
	assert.Equal(t, 200*time.Millisecond, results[2].Latency, "Expected latency for target 2")
}

func TestEnodeByAddress(t *testing.T) {
	key, _ := crypto.HexToECDSA("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	addr1 := crypto.PubkeyToAddress(key.PublicKey)
	addr2 := common.HexToAddress("0x222")

	// Create an enode for addr1
	enode1 := enode.NewV4(&key.PublicKey, net.ParseIP("192.168.1.1"), 30303, 0)
	committeeEnodes := []*enode.Node{enode1}

	node, found, _ := findEnode(0, committeeEnodes, addr1)
	assert.True(t, found, "Expected to find enode for addr1")
	assert.Equal(t, enode1, node, "Expected correct enode for addr1")

	_, found, _ = findEnode(0, committeeEnodes, addr2)
	assert.False(t, found, "Expected not to find enode for addr2")
}

func TestFindSortedEnode(t *testing.T) {
	addresses := make(map[common.Address]struct{})
	enodes := make([]*enode.Node, 0)
	total := 6
	for total > 0 {
		total--
		for {
			pubKey := randomPublicKey(t)
			if _, ok := addresses[crypto.PubkeyToAddress(pubKey)]; ok {
				continue
			}
			addresses[crypto.PubkeyToAddress(pubKey)] = struct{}{}
			enodes = append(enodes, enode.NewV4(&pubKey, net.ParseIP("192.168.1.1"), 30303, 0))
			break
		}
	}

	validators := make([]common.Address, 0)
	for i := 0; i < len(enodes)/2; i++ {
		validators = append(validators, crypto.PubkeyToAddress(*enodes[i].Pubkey()))
	}

	for i := 0; i < 3; i++ {
		validators = append(validators, common.BigToAddress(big.NewInt(int64(123*i))))
	}
	validators = append(validators, common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff"))

	validators, enodes = sortByAddress(validators, enodes)

	enodeAt := 0

	for _, member := range validators {
		_, found, biggerEnodeIndex := findEnode(enodeAt, enodes, member)
		enodeAt = biggerEnodeIndex
		_, ok := addresses[member]
		require.Equal(t, ok, found)
	}

}

func TestFetcher_ConcurrentPingPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pinger := mocks.NewMockPinger(ctrl)
	peerFinder := mocks.NewMockPeerFinder(ctrl)
	fetcher := NewFetcher(pinger)
	fetcher.SetBroadcaster(peerFinder)

	targets := make([]ping.Target, 10)
	for i := 0; i < 10; i++ {
		targets[i] = ping.Target{IP: "192.168.1." + string(rune(1+i)), Port: 30303 + i}
		pinger.EXPECT().Ping(gomock.Any(), targets[i]).
			Return(ping.Result{Latency: time.Duration(50+i) * time.Millisecond}).Times(1)
	}

	results := fetcher.pingPeers(context.Background(), targets)

	assert.Equal(t, 10, len(results), "Expected results for all targets")
	for i, result := range results {
		assert.Equal(t, time.Duration(50+i)*time.Millisecond, result.Latency, "Expected correct latency for target %d", i)
	}
}

func randomPublicKey(t *testing.T) ecdsa.PublicKey {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	return privateKey.PublicKey
}
