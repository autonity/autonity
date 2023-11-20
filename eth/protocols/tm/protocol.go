package tm

import (
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"math/big"
)

// Constants to match up protocol versions and messages
const (
	CNS1 = 1
)

// ProtocolName is the official short name of the `snap` protocol used during
// devp2p capability negotiation.
const ProtocolName = "aut-consensus"

// ProtocolVersions are the supported versions of the `snap` protocol (first
// is primary).
var ProtocolVersions = []uint{CNS1}

// protocolLengths are the number of implemented message corresponding to
// different protocol versions.
var protocolLengths = map[uint]uint64{CNS1: 8}

// MaxMessageSize is the maximum cap on the size of a consensus protocol message.
const MaxMessageSize = 10 * 1024 * 1024

// TODO: fix below constants numbering
const (
	NewBlockHashesMsg = 0x01
	NewBlockMsg       = 0x07
)

var (
	errMsgTooLarge    = errors.New("message too long")
	errDecode         = errors.New("invalid message")
	errInvalidMsgCode = errors.New("invalid message code")
	errBadRequest     = errors.New("bad request")
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
