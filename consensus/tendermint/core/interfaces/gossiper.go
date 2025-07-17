package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
)

type Gossiper interface {
	Gossip(committee *types.Committee, message message.Msg, sender common.Address)
	AskSync(committee *types.Committee, syncMsg *message.AskSyncMsg, clusteringThreshold int64) error
	SlowGossip(committee *types.Committee, message message.Msg, sender common.Address, clusteringThreshold int64)
	SetBroadcaster(broadcaster consensus.Broadcaster)
	Broadcaster() consensus.Broadcaster
	KnownMessages() *fixsizecache.Cache[common.Hash, bool]
	Address() common.Address
	UpdateStopChannel(chan struct{})
}
