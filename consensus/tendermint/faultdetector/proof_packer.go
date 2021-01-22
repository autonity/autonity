package faultdetector

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core/types"
)

// todo: this packing / unpacking logic can be replaced by native RLP encode/decode, going to remove it..
/* proof_packer.go provide pack/unpack functions for precompiled contract to decode proof from and to byte array.
	Pack
	type Proof struct {
		Rule     Rule
		Message  message
		Evidence []message
	}

	into a byte array which is represented by above elements for static unpacking.

	Rule uint8;				 // take the 1st byte as rule id.
	NumberOfMsg uint32;		 // take the next 4 bytes as number of msgs.
	LengthOfEach []uint32;   // take the next "NumberOfMsg" * 4 bytes as the set of lengths for each message.
	MessagePayloads []bytes; // take the left bytes as payload of each message tightly packed, slotted by own lengths.
*/

func decodeMessage(payload []byte) (*message, error) {
	// decode message with rlp.
	msg := new(core.Message)
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

	m := new(message)
	return m, nil
}

func lengthOfBuffer(proof *Proof) int {
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
func PackProof(proof *Proof) []byte {
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

func UnpackProof(packedProof []byte) (*Proof, error) {
	proof := new(Proof)

	// unpack rule id.
	proof.Rule = Rule(packedProof[1])
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


