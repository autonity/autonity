package tendermint

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	common "github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/algorithm"
	types "github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/rlp"
)

// This is used as a means to get around rlp's restriction on encoding negative
// integers. We rlp encode and decode consensus messages via this struct
// casting the round and valid round to and from int64 to uint64.
type rlpConsensusMesage struct {
	MsgType    algorithm.Step
	Height     uint64
	Round      uint64
	Value      algorithm.ValueID
	ValidRound uint64
}

// signedMessage is what is serialised an sent between nodes.
type signedMessage struct {
	// ConsensusMessage is a serialised rlpConsensusMesage.
	ConsensusMessage []byte
	// Signature is a Signature over the the Keccak256 Hash of
	// ConsensusMessage.
	Signature []byte
	// Value is a serialised types.Block, this is only set if the consensus
	// Message is a proposal.
	Value []byte
	// ProposerSeal is a Signature over the types.Block.Hash() of the block
	// serialised in Value, this is only set if Value is set, (i.e if this
	// Message is a proposal Message).
	ProposerSeal []byte
}

// Message is a convenience type that holds the information contained in a
// signed Message unpacked into a useful form. It is never serialised directly
// it is only used internally by the bridge to represent a Message.
type Message struct {
	Hash             common.Hash
	Address          common.Address
	ConsensusMessage *algorithm.ConsensusMessage
	Signature        []byte
	ProposerSeal     []byte
	Value            *types.Block
}

func (m *Message) String() string {
	return fmt.Sprintf("%s %v", addr(m.Address), m.ConsensusMessage)
}

func (m *Message) H() uint64 {
	return m.ConsensusMessage.Height
}

func (m *Message) R() int64 {
	return m.ConsensusMessage.Round
}

func (m *Message) VR() int64 {
	return m.ConsensusMessage.ValidRound
}

func (m *Message) Type() algorithm.Step {
	return m.ConsensusMessage.MsgType
}

func (m *Message) Sender() common.Address {
	return m.Address
}

func (m *Message) V() common.Hash {
	return common.Hash(m.ConsensusMessage.Value)
}

func (m *Message) MsgHash() common.Hash {
	return m.Hash
}

func (m *Message) Payload() []byte {
	var value []byte
	var err error
	if m.Type() == algorithm.Propose {
		value, err = rlp.EncodeToBytes(m.Value)
		if err != nil {
			return nil
		}
	}

	cmBytes, err := encodeConsensusMessage(m.H(), m.R(), m.VR(), m.Type(), algorithm.ValueID(m.V()))
	if err != nil {
		return nil
	}

	sm := &signedMessage{
		ConsensusMessage: cmBytes,
		Signature:        m.Signature,
		Value:            value,
		ProposerSeal:     m.ProposerSeal,
	}

	b, err := rlp.EncodeToBytes(sm)
	if err != nil {
		return nil
	}
	return b
}

// func (cm *algorithm.ConsensusMessage, )
func badMessageErr(description string, m *algorithm.ConsensusMessage, address common.Address) error {
	return fmt.Errorf("%s - Message: %q, from: %q", description, m, addr(address))
}

// TODO think about the values in the block, calling Hash on a header
// returns the Hash but filters out the bft fields. Except for the proposer
// seal.  This means that if we are building a proposal Message at this
// point we may have a fresh proposal from our miner, in which case it will
// not have a proposer seal or it could be that we are reproposing a valid
// Value. In which case we need to set the proposer seal to be empty before
// we Hash it and add our seal. This means that the Message we receive here
// from the algorithm will be wrong, because when we propose a valid Value
// we are using the Hash of a block that had a different proposer's seal. But we can catch this here.
//
// Alternatively we can ensure that the hashes that are passed around by
// the algorithm do not include proposer seals. We could store the proposer
// seal alongside the block and then participants could insert it into the
// block before they add their committed seal. This seems nice actually
// because it makes verifying the proposer seal easier. So that means that
// all values in the Message store are without any kind of seal.

// Ok I don't think we actually need the proposer seal on messages other than
// proposals. So we can just pass the block here rather than the store if we
// are serialising a proposal.
func EncodeSignedMessage(cm *algorithm.ConsensusMessage, key *ecdsa.PrivateKey, proposalBlock *types.Block) ([]byte, error) {

	// If we are building a proposal then serialise the Value and generate the
	// proposer seal.
	var value []byte
	var proposerSeal []byte
	var err error
	if cm.MsgType == algorithm.Propose {
		value, err = rlp.EncodeToBytes(proposalBlock)
		if err != nil {
			return nil, err
		}
		proposerSeal, err = crypto.Sign(cm.Value[:], key)
		if err != nil {
			return nil, err
		}
	}

	// Encode the consensus Message.
	cmBytes, err := encodeConsensusMessage(cm.Height, cm.Round, cm.ValidRound, cm.MsgType, cm.Value)
	if err != nil {
		return nil, err
	}

	// The Message Signature is over the Hash of serialised consensus Message.
	// We don't need to sign over the Value because we can verify that the
	// Value Hash matches cm.Value and cm was signed over.
	msgHash := crypto.Keccak256(cmBytes)
	signature, err := crypto.Sign(msgHash, key)
	if err != nil {
		return nil, err
	}

	sm := &signedMessage{
		ConsensusMessage: cmBytes,
		Signature:        signature,
		ProposerSeal:     proposerSeal,
		Value:            value,
	}
	return rlp.EncodeToBytes(sm)
}

func encodeConsensusMessage(height uint64, round, validRound int64, step algorithm.Step, value algorithm.ValueID) ([]byte, error) {
	rlpM := &rlpConsensusMesage{
		MsgType:    step,
		Height:     height,
		Round:      uint64(round),
		Value:      value,
		ValidRound: uint64(validRound),
	}
	return rlp.EncodeToBytes(rlpM)
}

func DecodeSignedMessage(msgBytes []byte) (*Message, error) {
	sm := &signedMessage{}
	err := rlp.Decode(bytes.NewBuffer(msgBytes), sm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signed Message: %v", err)
	}

	// Decode the consensus Message
	rlpMsg := &rlpConsensusMesage{}
	err = rlp.Decode(bytes.NewBuffer(sm.ConsensusMessage), rlpMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode consensus Message Message: %v", err)
	}
	cm := &algorithm.ConsensusMessage{
		MsgType:    rlpMsg.MsgType,
		Height:     rlpMsg.Height,
		Round:      int64(rlpMsg.Round),
		Value:      rlpMsg.Value,
		ValidRound: int64(rlpMsg.ValidRound),
	}

	// Get the sender Address, note we are not verifying the Signature here,
	// we are simply recovering the sender Address. Verification will occur at
	// the point when we check to see if the sender Address is part of the
	// committee for the block in question.
	msgHash := crypto.Keccak256(sm.ConsensusMessage)
	signerAddress, err := types.GetSignatureAddressHash(msgHash, sm.Signature)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve Address from Message Signature: %v - Message: %q",
			err,
			cm.String(),
		)
	}

	var value *types.Block
	if cm.MsgType == algorithm.Propose {
		// Get the proposer Address, note we are not verifying the Signature
		// here, we are simply recovering the proposer Address. Verification
		// will occur at the point when we check to see if the propser Address
		// matches the selected proposer for this height and round.
		proposerAddress, err := types.GetSignatureAddressHash(cm.Value[:], sm.ProposerSeal)
		if err != nil {
			return nil, badMessageErr(fmt.Sprintf("received proposal Message with invalid proposer seal: %v", err),
				cm,
				signerAddress,
			)
		}
		// Check that the proposer seal Address matches the signer Address.
		if proposerAddress != signerAddress {
			return nil, badMessageErr(fmt.Sprintf("proposer seal Address %q does not match sender Address", proposerAddress.String()),
				cm,
				signerAddress,
			)
		}

		// Check that a block was provided.
		if len(sm.Value) == 0 {
			return nil, badMessageErr(
				"received proposal Message without associated Value",
				cm,
				signerAddress,
			)
		}
		value = &types.Block{}
		err = rlp.Decode(bytes.NewBuffer(sm.Value), value)
		if err != nil {
			return nil, err
		}

		// Check proposal Message Hash and block Hash match.
		if value.Hash() != common.Hash(cm.Value) {
			return nil, badMessageErr(
				fmt.Sprintf("received proposal Message with mismatching Value having Hash %q", value.Hash().String()),
				cm,
				signerAddress,
			)
		}
	}
	return &Message{
		Hash:             common.Hash(sha256.Sum256(msgBytes)),
		Signature:        sm.Signature,
		ConsensusMessage: cm,
		Value:            value,
		ProposerSeal:     sm.ProposerSeal,
		Address:          signerAddress,
	}, nil
}
