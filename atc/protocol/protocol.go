package protocol

import (
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/forkid"
)

// Constants to match up protocol versions and messages
const (
	ATC_V1 = 1
)

// TODO: better name for protocol
// ProtocolName is the official short name of the `snap` protocol used during
// devp2p capability negotiation.
const ProtocolName = "atc"

// ProtocolVersions are the supported versions of the `snap` protocol (first
// is primary).
var ProtocolVersions = []uint{ATC_V1}

// protocolLengths are the number of implemented message corresponding to
// different protocol versions.
var protocolLengths = map[uint]uint64{ATC_V1: 22}

// MaxMessageSize is the maximum cap on the size of a consensus protocol message.
const MaxMessageSize = 10 * 1024 * 1024

// TODO: fix below constants numbering
const (
	StatusMsg = 0x00
)

var (
	errNoStatusMsg             = errors.New("no status message")
	errMsgTooLarge             = errors.New("message too long")
	errDecode                  = errors.New("invalid message")
	errInvalidMsgCode          = errors.New("invalid message code")
	errProtocolVersionMismatch = errors.New("protocol version mismatch")
	errNetworkIDMismatch       = errors.New("network ID mismatch")
	errGenesisMismatch         = errors.New("genesis mismatch")
	errForkIDRejected          = errors.New("fork ID rejected")
	errBadRequest              = errors.New("bad request")
)

// StatusPacket is the network packet for the status message for eth/64 and later.
type StatusPacket struct {
	ProtocolVersion uint32
	NetworkID       uint64
	Head            common.Hash
	Genesis         common.Hash
	ForkID          forkid.ID
}
