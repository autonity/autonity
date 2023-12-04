package tm

import (
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/forkid"
	"github.com/autonity/autonity/core/types"
	"math/big"
)

// Constants to match up protocol versions and messages
const (
	CNS1 = 1
)

// TODO: better name for protocol
// ProtocolName is the official short name of the `snap` protocol used during
// devp2p capability negotiation.
const ProtocolName = "aut-consensus"

// ProtocolVersions are the supported versions of the `snap` protocol (first
// is primary).
var ProtocolVersions = []uint{CNS1}

// protocolLengths are the number of implemented message corresponding to
// different protocol versions.
var protocolLengths = map[uint]uint64{CNS1: 22}

// MaxMessageSize is the maximum cap on the size of a consensus protocol message.
const MaxMessageSize = 10 * 1024 * 1024

// TODO: fix below constants numbering
const (
	StatusMsg         = 0x00
	NewBlockHashesMsg = 0x01
	NewBlockMsg       = 0x07
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

// NewBlockHashesPacket is the network packet for the block announcements.
type NewBlockHashesPacket []struct {
	Hash   common.Hash // Hash of one particular block being announced
	Number uint64      // Number of one particular block being announced
}

// NewBlockPacket is the network packet for the block propagation message.
type NewBlockPacket struct {
	Block *types.Block
	TD    *big.Int
}

// Packet represents a p2p message in the `consensus` protocol.
type Packet interface {
	Name() string // Name returns a string corresponding to the message type.
	Kind() byte   // Kind returns the message type.
}

// StatusPacket is the network packet for the status message for eth/64 and later.
type StatusPacket struct {
	ProtocolVersion uint32
	NetworkID       uint64
	TD              *big.Int
	Head            common.Hash
	Genesis         common.Hash
	ForkID          forkid.ID
}
