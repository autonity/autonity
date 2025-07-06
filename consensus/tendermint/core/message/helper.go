package message

import (
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
)

func (e *EvidenceVote) ToPrevote() *Prevote {
	signers, err := e.signers.ToVoteSigners()
	if err != nil {
		panic(err.Error())
	}
	payload, _ := rlp.EncodeToBytes(extVote[uint16]{
		Code:      PrevoteCode,
		Round:     uint64(e.round), // #nosec
		Height:    e.height,
		Value:     e.value,
		Signers:   signers.SignersBase,
		Signature: e.signature.(*blst.BlsSignature),
	})

	return &Prevote{
		value: e.value,
		vote: vote[uint16]{
			signers: signers.SignersBase,
			base: base{
				height:         e.height,
				round:          e.round,
				signatureInput: e.signatureInput,
				signature:      e.signature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				verified:       e.verified,
				preverified:    e.preverified,
				signerKey:      e.signerKey,
			},
		},
	}
}
