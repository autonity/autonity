package message

import (
	"fmt"
	"io"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/blst"
)

// TODO(lorenzo) refinements, removed the lock, check that after the aggregator no-one modifies fields

type base struct {
	// populated at decoding phase
	height         uint64
	round          int64
	signatureInput common.Hash
	signature      blst.Signature
	payload        []byte
	hash           common.Hash
	verified       bool
	preverified    bool

	// populated at PreValidate() phase
	senderKey blst.PublicKey
}

func (b *base) H() uint64 {
	return b.height
}

func (b *base) R() int64 {
	return b.round
}

func (b *base) SignatureInput() common.Hash {
	return b.signatureInput
}

func (b *base) Signature() blst.Signature {
	return b.signature
}

func (b *base) Payload() []byte {
	return b.payload
}

func (b *base) EncodeRLP(w io.Writer) error {
	_, err := w.Write(b.payload)
	return err
}

func (b *base) Hash() common.Hash {
	return b.hash
}

// Bls key that needs to be used to verify the signature. Can be an aggregated key.
func (b *base) SenderKey() blst.PublicKey {
	if !b.preverified {
		panic("Trying to access sender key on not preverified message")
	}
	return b.senderKey
}

func (b *base) Validate() error {
	if !b.preverified {
		panic("Trying to verify a message that was not previously pre-verified")
	}
	if b.verified {
		return nil
	}

	if valid := b.signature.Verify(b.senderKey, b.signatureInput[:]); !valid {
		return ErrBadSignature
	}

	b.verified = true

	return nil
}

func (b *base) String() string {
	return fmt.Sprintf("h: %v, r: %v, verified: %v",
		b.height, b.round, b.verified)
}
