package messages

import "math/big"

type MessageType uint8

const (
	ProposalType MessageType = iota
	PrevoteType
	PrecommitType
)

type Message struct {
	Type     MessageType
	Proposal ProposalMessage
	Vote     VoteMessage
}

type ProposalMessage struct {
	Round      int64
	Height     *big.Int
	ValidRound int64
	NodeId     []byte
	Value      []byte
}

type VoteMessage struct {
	Round   int64
	Height  *big.Int
	NodeId  []byte
	ValueId []byte
}
