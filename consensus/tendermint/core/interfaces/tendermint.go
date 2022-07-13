package interfaces

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
)

type Tendermint interface {
	Start(ctx context.Context, contract *autonity.Contract)
	Stop()
	GetCurrentHeightMessages() []*messageutils.Message
	CoreState() types.TendermintState
}
