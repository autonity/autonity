package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
)

type Gossiper interface {
	Gossip(committee *types.Committee, message message.Msg)
	GossipPayload(committee *types.Committee, code uint8, hash common.Hash, payload []byte)
	AskSync(committee *types.Committee)
	SetBroadcaster(broadcaster consensus.Broadcaster)
	Broadcaster() consensus.Broadcaster
	KnownMessages() *fixsizecache.Cache[common.Hash, bool]
	Address() common.Address
	UpdateStopChannel(chan struct{})
}
