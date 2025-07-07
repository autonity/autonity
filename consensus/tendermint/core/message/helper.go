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

func (p *Prevote) ToEvidence() *EvidenceVote {
	signers := p.signers.ToQuorumSigners()
	payload, _ := rlp.EncodeToBytes(extVote[uint32]{
		Code:      EvidenceVoteCode,
		Round:     uint64(p.round), // #nosec
		Height:    p.height,
		Value:     p.value,
		Signers:   signers.SignersBase,
		Signature: p.signature.(*blst.BlsSignature),
	})

	return &EvidenceVote{
		value: p.value,
		vote: vote[uint32]{
			signers: signers.SignersBase,
			base: base{
				height:         p.height,
				round:          p.round,
				signatureInput: p.signatureInput,
				signature:      p.signature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				verified:       p.verified,
				preverified:    p.preverified,
				signerKey:      p.signerKey,
			},
		},
	}
}
