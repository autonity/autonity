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

// The signature is over the message and the ProposerSeal, it doesn't include
// value because Value is only sent for proposals and the hash of the value is
// contained in Message.
type signedMessage struct {
	Signature        []byte
	ConsensusMessage []byte
	Value            []byte
	ProposerSeal     []byte
	// We don't need the committed seal separate from the signature since we
	// can just recreate the signature by constructing a precommit message from
	// the info in a block.
}

type message struct {
	hash             common.Hash
	address          common.Address
	consensusMessage *algorithm.ConsensusMessage
	signature        []byte
	proposerSeal     []byte
	value            *types.Block
}

func (m *message) String() string {
	return "bla"
}

// func (cm *algorithm.ConsensusMessage, )
func badMessageErr(description string, m *algorithm.ConsensusMessage, address common.Address) error {
	return fmt.Errorf("%s - message: %q, from: %q", description, m, addr(address))
}

// TODO think about the values in the block, calling hash on a header
// returns the hash but filters out the bft fields. Except for the proposer
// seal.  This means that if we are building a proposal message at this
// point we may have a fresh proposal from our miner, in which case it will
// not have a proposer seal or it could be that we are reproposing a valid
// value. In which case we need to set the proposer seal to be empty before
// we hash it and add our seal. This means that the message we receive here
// from the algorithm will be wrong, because when we propose a valid value
// we are using the hash of a block that had a different proposer's seal. But we can catch this here.
//
// Alternatively we can ensure that the hashes that are passed around by
// the algorithm do not include proposer seals. We could store the proposer
// seal alongside the block and then participants could insert it into the
// block before they add their committed seal. This seems nice actually
// because it makes verifying the proposer seal easier. So that means that
// all values in the message store are without any kind of seal.

func encodeSignedMessage(cm *algorithm.ConsensusMessage, key *ecdsa.PrivateKey, store *messageStore) ([]byte, error) {

	// First generate the proposer seal, which is either our seal if we are
	// building a proposal, or its the seal of the proposal we are voting for,
	// if we are not voting for nil.
	var proposerSeal []byte
	var err error
	switch {
	case cm.MsgType == algorithm.Propose:
		proposerSeal, err = crypto.Sign(cm.Value[:], key)
		if err != nil {
			return nil, err
		}
	case cm.Value != algorithm.NilValue:
		proposerSeal = store.matchingProposal(cm).proposerSeal
	}

	// Now encode the consensus message and generate the message signature.
	rlpM := &rlpConsensusMesage{
		MsgType:    cm.MsgType,
		Height:     cm.Height,
		Round:      uint64(cm.Round),
		Value:      cm.Value,
		ValidRound: uint64(cm.ValidRound),
	}
	cmBytes, err := rlp.EncodeToBytes(rlpM)
	if err != nil {
		return nil, err
	}

	// The message signature is over the consensus message and the proposer seal.
	msgHash := crypto.Keccak256(append(cmBytes, proposerSeal...))
	signature, err := crypto.Sign(msgHash, key)
	if err != nil {
		return nil, err
	}

	sm := &signedMessage{
		ConsensusMessage: cmBytes,
		ProposerSeal:     proposerSeal,
		Signature:        signature,
	}

	// Add the encoded value if we are building a proposal. We don't need to
	// sign over the value because we can verify that the value hash matches
	// cm.Value and cm was signed over.
	var value []byte
	if cm.MsgType == algorithm.Propose {
		block := store.value(common.Hash(cm.Value))
		if block == nil {
			return nil, fmt.Errorf("unable to find block for proposal: %s", cm.String())
		}
		value, err = rlp.EncodeToBytes(block)
		if err != nil {
			return nil, err
		}
	}
	sm.Value = value
	return rlp.EncodeToBytes(sm)
}

func BuildCommitment(proposerSeal []byte, height uint64, round int64, value algorithm.ValueID) ([]byte, error) {
	rlpM := &rlpConsensusMesage{
		MsgType: algorithm.Precommit,
		Height:  height,
		Round:   uint64(round),
		Value:   value,
	}
	cmBytes, err := rlp.EncodeToBytes(rlpM)
	if err != nil {
		return nil, err
	}
	// return crypto.Keccak256(append(cmBytes, proposerSeal...)), nil
	return crypto.Keccak256(append(cmBytes, proposerSeal...)), nil
}

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

	//println("da msg", cm.String())

	// Verify the signatures
	// Get the sender address
	msgHash := crypto.Keccak256(append(sm.ConsensusMessage, sm.ProposerSeal...))
	signerAddress, err := types.GetSignatureAddressHash(msgHash, sm.Signature)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve address from message signature: %v - message: %q",
			err,
			cm.String(),
		)
	}

	var value *types.Block
	// Votes for nil value do not have a proposer seal and cannot be of type
	// algorithm.Propose.
	if cm.Value != algorithm.NilValue {

		// Get the proposer address
		proposerAddress, err := types.GetSignatureAddressHash(cm.Value[:], sm.ProposerSeal)
		if err != nil {
			return nil, badMessageErr(fmt.Sprintf("received proposal message with invalid proposer seal: %v", err),
				cm,
				signerAddress,
			)
		}
		if cm.MsgType == algorithm.Propose {
			// Check that the proposer seal address matches the signer address.
			if proposerAddress != signerAddress {
				return nil, badMessageErr(fmt.Sprintf("proposer seal address %q does not match sender address", proposerAddress.String()),
					cm,
					signerAddress,
				)
			}

			// Check that a block was provided if the message is a proposal.
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
