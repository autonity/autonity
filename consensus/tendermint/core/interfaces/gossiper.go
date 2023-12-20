package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	lru "github.com/hashicorp/golang-lru"
)

type Gossiper interface {
	Gossip(committee *types.Committee, message message.Msg)
	AskSync(committee *types.Committee)
	SetBroadcaster(broadcaster consensus.Broadcaster)
	Broadcaster() consensus.Broadcaster
	RecentMessages() *lru.ARCCache
	KnownMessages() *lru.ARCCache
	Address() common.Address
}
