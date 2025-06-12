package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/router/network"
	"github.com/autonity/autonity/consensus/tendermint/router/ping"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/p2p/enode"
)

type NetworkProvider interface {
	Clusters() network.Clusters
	UpdateClusters(network.Clusters)
}

type LatencyProvider interface {
	SetBroadcaster(broadcaster PeerFinder)
	Fetch(validators []common.Address, self common.Address) (map[common.Address]uint, []common.Address, error)
	Pinger() ping.Pinger
	SetPinger(ping.Pinger)
}

type PeerFinder interface {
	FindPeers([]common.Address) map[common.Address]consensus.Peer
	FindPeer(addr common.Address) (consensus.Peer, bool)
	CommitteeEnodes() []*enode.Node
}

type PeerSelector interface {
	SelectPeers(committee *types.Committee, msg message.Msg, from common.Address) ([]common.Address, error)
	SetBroadcaster(broadcaster PeerFinder)
}

type BlockChainProvider interface {
	LatestEpoch() (*types.EpochInfo, error)
	SubscribeEpochHeadEvent(chan<- core.EpochHeadEvent) event.Subscription
}
