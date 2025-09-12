package message

import (
	"fmt"
	"io"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/blst"
)

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

	// populated on demand
	signerKey blst.PublicKey
}

func (b *base) Verified() bool {
	return b.verified
}

func (b *base) PreVerified() bool {
	return b.preverified
}

func (b *base) H() uint64 {
	return b.height
}

func (b *base) R() int64 {
	return b.round
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

func (b *base) Validate() error {
	if !b.preverified {
		panic("Trying to verify a message that was not previously pre-verified")
	}
	if b.verified {
		return nil
	}

	if valid := b.signature.Verify(b.signerKey, b.signatureInput[:], blst.DefaultAssumeZeroValid); !valid {
		return ErrBadSignature
	}

	b.verified = true

	return nil
}

func (b *base) String() string {
	return fmt.Sprintf("h: %v, r: %v, hash: %v, preverified: %v, verified: %v",
		b.height, b.round, b.hash, b.preverified, b.verified)
}
