package message

import (
	"fmt"
	"io"
	"math/big"

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

	// populated at PreValidate() phase
	senderKey blst.PublicKey
	power     *big.Int
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
	return b.senderKey
}

func (b *base) Power() *big.Int {
	return b.power
}

func (b *base) Validate() error {
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
	return fmt.Sprintf("h: %v, r: %v, power: %v, verified: %v",
		b.height, b.round, b.power, b.verified)
}

/* //TODO(lorenzo) refinements, for later
// used by tests to simulate unverified messages
func (b *base) unvalidate() {
	b.verified = false
	b.power = new(big.Int)
	b.sender = common.Address{}
}*/
