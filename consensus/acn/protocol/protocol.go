package protocol

import (
	"errors"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/forkid"
)

// Constants to match up protocol versions and messages
const (
	ACNv1 = 1
)

// ProtocolName is the official short name of the autonity consensus network protocol used during
// devp2p capability negotiation.
const ProtocolName = "acn"

// ProtocolVersions are the supported versions of the `snap` protocol (first
// is primary).
var ProtocolVersions = []uint{ACNv1}

// todo(piyush): length for ACN should be 6 because of 1 status message(0x00) and
// and 5 protocol message which have legacy codes(staring from 0x11) i.e. length 22 for now.
// protocolLengths are the number of implemented message corresponding to
// different protocol versions.
var protocolLengths = map[uint]uint64{ACNv1: 22}

// MaxMessageSize is the maximum cap on the size of a consensus protocol message.
const MaxMessageSize = 10 * 1024 * 1024

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
)

// StatusPacket is the network packet for the status message for eth/64 and later.
type StatusPacket struct {
	ProtocolVersion uint32
	NetworkID       uint64
	TD              *big.Int
	Head            common.Hash
	Genesis         common.Hash
	ForkID          forkid.ID
}
