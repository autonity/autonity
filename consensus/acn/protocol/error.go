package protocol

import (
	"errors"

	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/p2p"
)

const (
	errACNHandler = iota
)

const acnErrorSuspensionSpan = 60 // num of blocks, peer is not allowed to set up connection

var errorToString = map[int]string{
	errACNHandler: "acn message handling error",
}

func newACNError(backend Backend, err error) *p2p.ProtocolError {
	desc, ok := errorToString[errACNHandler]
	if !ok {
		panic("invalid error code")
	}
	pError := &p2p.ProtocolError{Suspension: func() uint64 {
		var suspension = uint64(acnErrorSuspensionSpan)
		if errors.Is(err, message.ErrBadSignature) {
			// TODO(lorenzo) implement more harsh exponential approach disconnection?
			suspension = backend.Chain().ProtocolContracts().Cache.EpochPeriod().Uint64()
		}
		//Note: Can add more errors here for different suspension span
		return suspension
	}, Code: errACNHandler, Message: desc}
	pError.Message += ": " + err.Error()
	return pError
}
