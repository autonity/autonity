package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/params"
)

// checkChallenge implemented as a native contract to take an on-chain challenge.
type checkChallenge struct{
	chainContext ChainContext
}
func (c *checkChallenge) InitChainContext(chain ChainContext) {
	c.chainContext = chain
	return
}
func (c *checkChallenge) RequiredGas(_ []byte) uint64 {
	return params.TakeChallengeGas
}

// checkChallenge, take challenge from AC by copy the packed byte array, decode and
// validate it, the on challenge client should send the proof of innocent via a transaction.
func (c *checkChallenge) Run(input []byte) ([]byte, error) {
	if len(input) == 0 {
		panic(fmt.Errorf("invalid proof of innocent - empty"))
	}

	err := CheckChallenge(input, c.chainContext)
	if err != nil {
		return false32Byte, fmt.Errorf("invalid proof of challenge %v", err)
	}

	return true32Byte, nil
}

// checkProof implemented as a native contract to validate an on-chain innocent proof.
type checkProof struct{
	chainContext ChainContext
}
func (c *checkProof) InitChainContext(chain ChainContext) {
	c.chainContext = chain
	return
}
func (c *checkProof) RequiredGas(_ []byte) uint64 {
	return params.CheckInnocentGas
}

// checkProof, take proof from AC by copy the packed byte array, decode and validate it.
func (c *checkProof) Run(input []byte) ([]byte, error) {
	// take an on-chain innocent proof, tell the results of the checking
	if len(input) == 0 {
		panic(fmt.Errorf("invalid proof of innocent - empty"))
	}

	err := CheckProof(input, c.chainContext)
	if err != nil {
		return false32Byte, fmt.Errorf("invalid proof of innocent %v", err)
	}

	return true32Byte, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// validate the proof is a valid challenge.
func validateChallenge(c *types.Proof, chain ChainContext) error {
	// get committee from block header
	h, err := c.Message.Height()
	if err != nil {
		return nil
	}

	// todo: double check h - 1?
	header := chain.GetHeader(c.ParentHash, h.Uint64())

	// check if evidences senders are presented in committee.
	for i:=0; i < len(c.Evidence); i++ {
		member := header.CommitteeMember(c.Evidence[i].Address)
		if member == nil {
			return fmt.Errorf("invalid evidence for proof of susipicous message")
		}
	}

	// todo: check if the suspicious message is proved by given evidence as a valid suspicion.
	return nil
}

// validate the innocent proof is valid.
func validateInnocentProof(in *types.Proof, chain ChainContext) error {
	// get committee from block header
	h, err := in.Message.Height()
	if err != nil {
		return nil
	}
	header := chain.GetHeader(in.ParentHash, h.Uint64())

	// check if evidences senders are presented in committee.
	for i:=0; i < len(in.Evidence); i++ {
		member := header.CommitteeMember(in.Evidence[i].Address)
		if member == nil {
			return fmt.Errorf("invalid evidence for proof of susipicous message")
		}
	}
	// todo: check if the suspicious message is proved by given evidence as an innocent behavior.
	return nil
}

// used by those who is on-challenge, to get innocent proof from msg store.
func resolveChallenge(c *types.Proof) (innocentProof *types.Proof, err error) {
	// todo: get innocent proof from msg store. Would need to distinguish the different rules.
	return nil, nil
}

// Check the proof of innocent, it is called from precompiled contracts of EVM package.
func CheckProof(packedProof []byte, chain ChainContext) error {
	innocentProof, err := UnpackProof(packedProof)
	if err != nil {
		return err
	}
	err = validateInnocentProof(innocentProof, chain)
	if err != nil {
		return err
	}

	return nil
}

// validate challenge, call from EVM package.
func CheckChallenge(packedProof []byte, chain ChainContext) error {
	challenge, err := UnpackProof(packedProof)
	if err != nil {
		return err
	}

	err = validateChallenge(challenge, chain)
	if err != nil {
		return err
	}
	return nil
}


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// pack / unpack, rlp encode / decode
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// todo: this packing / unpacking logic can be replaced by native RLP encode/decode, going to remove it..
/* provide pack/unpack functions for precompiled contract to decode proof from and to byte array.
Pack
type Proof struct {
	Rule     Rule
	ConsensusMessage  message
	Evidence []message
}

into a byte array which is represented by above elements for static unpacking.

Rule uint8;				 // take the 1st byte as rule id.
NumberOfMsg uint32;		 // take the next 4 bytes as number of msgs.
LengthOfEach []uint32;   // take the next "NumberOfMsg" * 4 bytes as the set of lengths for each message.
MessagePayloads []bytes; // take the left bytes as payload of each message tightly packed, slotted by own lengths.
*/

func decodeMessage(payload []byte) (*types.ConsensusMessage, error) {
	// decode message with rlp.
	msg := new(types.ConsensusMessage)
	err := msg.FromPayload(payload)
	if err != nil {
		return nil, err
	}

	// check if message is signed correctly.
	var noSigPayload []byte
	noSigPayload, err = msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}

	signer, err := types.GetSignatureAddress(noSigPayload, msg.Signature)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(msg.Address.Bytes(), signer.Bytes()) {
		return nil, fmt.Errorf("message is not signed by valid sender")
	}

	m := new(types.ConsensusMessage)
	return m, nil
}

func lengthOfBuffer(proof *types.Proof) int {
	length := 0
	// Rule uint8 take 1 byte
	length += 1
	// NumberOfMsg uint32 take 4 bytes
	length += 4
	// lengthOfEach []uint32 takes (len(proof.Evidence) + 1) * 4
	length += (len(proof.Evidence) + 1) * 4

	payloads := 0
	payloads += len(proof.Message.Payload())

	for i:=0; i < len(proof.Evidence); i++  {
		payloads += len(proof.Evidence[i].Payload())
	}
	return length + payloads
}

// should be used by node who want to provide proof of a challenge or an innocent.
func PackProof(proof *types.Proof) []byte {
	length := lengthOfBuffer(proof)
	output := make([]byte, length)

	// pack rule id uint8 into first byte of byte array.
	output[0] = byte(proof.Rule)

	// pack NumOfMsg into the next 4 bytes of byte array.
	binary.BigEndian.PutUint32(output[1:5], uint32(len(proof.Evidence) + 1))

	// pack lengthOfEachMsg into the next (len(proof.Evidence) + 1) * 4 bytes
	lenOffset := 5
	msgOffSet := 5 + (len(proof.Evidence) + 1) * 4

	for i := 0; i < len(proof.Evidence) ; i++  {
		// pack len of each msg's payload into lengthOfEach array.
		lenPayload := len(proof.Evidence[i].Payload())
		binary.BigEndian.PutUint32(output[lenOffset : lenOffset + 4], uint32(lenPayload))
		lenOffset += 4 // step forward offset of length for each msg payload.

		// pack payloads to payload arrays.
		copy(output[msgOffSet : msgOffSet+lenPayload], proof.Evidence[i].Payload())
		msgOffSet += lenPayload // step forward offset of payloads
	}

	// Pack the message which is considered to be misbehavior at the last slot of messages set.
	binary.BigEndian.PutUint32(output[lenOffset : lenOffset + 4], uint32(len(proof.Message.Payload())))
	copy(output[msgOffSet : msgOffSet+len(proof.Message.Payload())], proof.Message.Payload())

	return output
}

func UnpackProof(packedProof []byte) (*types.Proof, error) {
	proof := new(types.Proof)

	// unpack rule id.
	proof.Rule = packedProof[1]
	// unpack num of msg from buffer.
	numOfMsg := binary.BigEndian.Uint32(packedProof[1:5])

	// unpack evidence msgs.
	lenOffset := 5
	msgOffSet := 5 + numOfMsg * 4

	for i := 0; i < int(numOfMsg) - 1; i++ {
		payloadLen := binary.BigEndian.Uint32(packedProof[lenOffset: lenOffset+4])
		lenOffset += 4

		m, err := decodeMessage(packedProof[msgOffSet: msgOffSet+payloadLen])
		if err != nil {
			return nil, err
		}

		proof.Evidence = append(proof.Evidence, *m)

		// proof.Evidence = append(proof.Evidence, packedProof[msgOffSet: msgOffSet+payloadLen])
		msgOffSet += payloadLen
	}

	// unpack the msg which is considered to be misbehavior.
	message, err := decodeMessage(packedProof[msgOffSet:])
	if err != nil {
		return nil, err
	}
	proof.Message = *message

	return proof, nil
}