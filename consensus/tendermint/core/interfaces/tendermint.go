package interfaces

import (
	"context"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
)

type Tendermint interface {
	Start(ctx context.Context, contract *autonity.ProtocolContracts)
	Stop()
	CurrentHeightMessages() []*message.Message
	CoreState() core.TendermintState
}
