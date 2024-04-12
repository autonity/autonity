package interfaces

import (
	lru "github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
)

type Gossiper interface {
	Gossip(committee types.Committee, message message.Msg)
	AskSync(header *types.Header)
	SetBroadcaster(broadcaster consensus.Broadcaster)
	Broadcaster() consensus.Broadcaster
	RecentMessages() *lru.LRU[common.Address, *lru.LRU[common.Hash, bool]]
	KnownMessages() *lru.LRU[common.Hash, bool]
	Address() common.Address
	UpdateStopChannel(chan struct{})
}
