package message

import (
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto/blst"
)

// Msg can represent both:
// 1. an individual message as sent by a validator
// 2. an aggregate message, resulting from bls signature aggregation of multiple individual messages
type Msg interface {
	// Code returns the message code, it must always match the concrete type.
	Code() uint8

	// R returns the message round.
	R() int64

	// H returns the message height.
	H() uint64

	// Value returns the block hash being voted for.
	Value() common.Hash

	// Power returns the message voting power.
	Power() *big.Int

	// String returns a string description of the message.
	String() string

	// Hash returns the hash of the message. This is not available if the underlying payload
	// hasn't be assigned.
	Hash() common.Hash

	// Payload returns the rlp-encoded payload ready to be broadcasted.
	Payload() []byte

	// Signature returns the signature of this message
	Signature() blst.Signature

	//TODO(lorenzo) add comments
	Validate() error
	PreValidate(header *types.Header) error
	SignatureInput() common.Hash
	SenderKey() blst.PublicKey
}

type IndividualMsg interface {
	Sender() common.Address
	SenderIndex() int
	Msg
}

type AggregateMsg interface {
	Senders() *types.SendersInfo
	Msg
}
