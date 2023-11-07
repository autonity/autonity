package interfaces

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	lru "github.com/hashicorp/golang-lru"
)

type Gossiper interface {
	// Gossip consensus message to the other committee members
	Gossip(committee types.Committee, payload []byte)
	AskSync(header *types.Header)
	SetBroadcaster(broadcaster consensus.Broadcaster)
	Broadcaster() consensus.Broadcaster
	RecentMessages() *lru.ARCCache
	KnownMessages() *lru.ARCCache
	Address() common.Address
}
