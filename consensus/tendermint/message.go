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
	// Signature is a signature over the the Keccak256 hash of
	// ConsensusMessage.
	Signature []byte
	// Value is a serialised types.Block, this is only set if the consensus
	// message is a proposal.
	Value []byte
	// ProposerSeal is a signature over the types.Block.Hash() of the block
	// serialised in Value, this is only set if Value is set, (i.e if this
	// message is a proposal message).
	ProposerSeal []byte
}

// message is a convenience type that holds the information contained in a
// signed message unpacked into a useful form. It is never serialised directly
// it is only used internally by the bridge to represent a message.
type message struct {
	hash             common.Hash
	address          common.Address
	consensusMessage *algorithm.ConsensusMessage
	signature        []byte
	proposerSeal     []byte
	value            *types.Block
}

func (m *message) String() string {
	return fmt.Sprintf("%s %v", addr(m.address), m.consensusMessage)
}

func badMessageErr(description string, m *algorithm.ConsensusMessage, address common.Address) error {
	return fmt.Errorf("%s - message: %q, from: %q", description, m, addr(address))
}

// encodeSignedMessage constructs a marshalled signed message from the given
// parameters, the proposalBlock is expected when the consensus message is a
// proposal.
func encodeSignedMessage(cm *algorithm.ConsensusMessage, key *ecdsa.PrivateKey, proposalBlock *types.Block) ([]byte, error) {

	// If we are building a proposal then serialise the value and generate the
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

	// Encode the consensus message.
	cmBytes, err := encodeConsensusMessage(cm.Height, cm.Round, cm.ValidRound, cm.MsgType, cm.Value)
	if err != nil {
		return nil, err
	}

	// The message signature is over the hash of serialised consensus message.
	// We don't need to sign over the value because we can verify that the
	// value hash matches cm.Value and cm was signed over.
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

// decodeSignedMessage decodes the given bytes into a *message. Note, we cannot
// verify the signature at this point, verification happens when we see if the
// signer address is part of the committee.
func decodeSignedMessage(msgBytes []byte) (*message, error) {
	sm := &signedMessage{}
	err := rlp.Decode(bytes.NewBuffer(msgBytes), sm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signed message: %v", err)
	}

	// Decode the consensus message
	rlpMsg := &rlpConsensusMesage{}
	err = rlp.Decode(bytes.NewBuffer(sm.ConsensusMessage), rlpMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode consensus message message: %v", err)
	}
	cm := &algorithm.ConsensusMessage{
		MsgType:    rlpMsg.MsgType,
		Height:     rlpMsg.Height,
		Round:      int64(rlpMsg.Round),
		Value:      rlpMsg.Value,
		ValidRound: int64(rlpMsg.ValidRound),
	}

	// Get the sender address, note we are not verifying the signature here,
	// we are simply recovering the sender address. Verification will occur at
	// the point when we check to see if the sender address is part of the
	// committee for the block in question.
	msgHash := crypto.Keccak256(sm.ConsensusMessage)
	signerAddress, err := types.GetSignatureAddressHash(msgHash, sm.Signature)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve address from message signature: %v - message: %q",
			err,
			cm.String(),
		)
	}

	var value *types.Block
	if cm.MsgType == algorithm.Propose {
		// Get the proposer address, note we are not verifying the signature
		// here, we are simply recovering the proposer address. Verification
		// will occur at the point when we check to see if the propser address
		// matches the selected proposer for this height and round.
		proposerAddress, err := types.GetSignatureAddressHash(cm.Value[:], sm.ProposerSeal)
		if err != nil {
			return nil, badMessageErr(fmt.Sprintf("received proposal message with invalid proposer seal: %v", err),
				cm,
				signerAddress,
			)
		}
		// Check that the proposer seal address matches the signer address.
		if proposerAddress != signerAddress {
			return nil, badMessageErr(fmt.Sprintf("proposer seal address %q does not match sender address", proposerAddress.String()),
				cm,
				signerAddress,
			)
		}

		// Check that a block was provided.
		if len(sm.Value) == 0 {
			return nil, badMessageErr(
				"received proposal message without associated value",
				cm,
				signerAddress,
			)
		}
		value = &types.Block{}
		err = rlp.Decode(bytes.NewBuffer(sm.Value), value)
		if err != nil {
			return nil, err
		}

		// Check proposal message hash and block hash match.
		if value.Hash() != common.Hash(cm.Value) {
			return nil, badMessageErr(
				fmt.Sprintf("received proposal message with mismatching value having hash %q", value.Hash().String()),
				cm,
				signerAddress,
			)
		}
	}
	return &message{
		hash:             common.Hash(sha256.Sum256(msgBytes)),
		signature:        sm.Signature,
		consensusMessage: cm,
		value:            value,
		proposerSeal:     sm.ProposerSeal,
		address:          signerAddress,
	}, nil
}
