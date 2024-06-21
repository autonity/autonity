package types

import (
	"errors"
	"golang.org/x/crypto/blake2b"

	"github.com/autonity/autonity/crypto"

	lru "github.com/hashicorp/golang-lru"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/rlp"
)

var (
	// BFTDigest represents a hash of "Tendermint byzantine fault tolerance"
	// to identify whether the block is from BFT consensus engine
	//TODO differentiate the digest between IBFT and Tendermint
	BFTDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	BFTExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity

	inmemoryAddresses  = 500 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)

	// ErrInvalidBFTHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidBFTHeaderExtra = errors.New("invalid pos header extra-data")
	// ErrInvalidSignature is returned when given signature is not signed by given address.
	ErrInvalidSignature = errors.New("invalid signature")
	// ErrInvalidQuorumCertificate is returned if the committed seal is not signed by any of parent validators.
	ErrInvalidQuorumCertificate = errors.New("invalid quorum certificate")
	// ErrEmptyQuorumCertificate is returned if the field of quorum certificate is empty.
	ErrEmptyQuorumCertificate = errors.New("empty quorum certificate")
	// ErrNegativeRound is returned if the round field is negative
	ErrNegativeRound = errors.New("negative round")
)

// BFTFilteredHeader returns a filtered header which some information (like proposerSeal, quorumCertificate)
// are clean to fulfill the BFT hash rules.
func BFTFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	if !keepSeal {
		newHeader.ProposerSeal = []byte{}
	}
	newHeader.QuorumCertificate = AggregateSignature{}
	newHeader.Round = 0
	newHeader.Extra = []byte{}
	return newHeader
}

// SigHash returns the hash which is used as input for the BFT
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func SigHash(header *Header) common.Hash {
	inputData, err := rlp.EncodeToBytes(BFTFilteredHeader(header, false))
	if err != nil {
		log.Error("can't encode header", "err", err, "header", header)
		return common.Hash{}
	}
	return blake2b.Sum256(inputData)
}

// ECRecover extracts the Ethereum account address from a signed header.
func ECRecover(header *Header) (common.Address, error) {
	hash := SigHash(header)
	// note that we can use the hash alone for the key
	// as the hash is over an object containing the proposer address
	// e.g. every proposer will have a different header hash for the same block content.
	if addr, ok := recentAddresses.Get(hash); ok {
		return addr.(common.Address), nil
	}
	addr, err := crypto.SigToAddr(hash[:], header.ProposerSeal)
	if err != nil {
		return addr, err
	}
	recentAddresses.Add(hash, addr)
	return addr, nil
}

// TODO: All these Write* functions do useless checks as we always create the input ourselves. Remove them?

// WriteSeal writes the extra-data field of the given header with the given seals.
func WriteSeal(h *Header, seal []byte) error {
	if len(seal) != common.SealLength {
		return ErrInvalidSignature
	}
	h.ProposerSeal = make([]byte, len(seal))
	copy(h.ProposerSeal, seal)
	return nil
}

// WriteRound writes the round field of the block header.
func WriteRound(h *Header, round int64) error {
	if round < 0 {
		return ErrNegativeRound
	}
	h.Round = uint64(round)
	return nil
}

// WriteQuorumCertificate writes the extra-data field of a block header with given quorumCertificate
func WriteQuorumCertificate(h *Header, quorumCertificate AggregateSignature) error {
	// TODO(Lorenzo) I think these cases can never happen and could be removed. same goes for the other Write* functions
	if quorumCertificate.Signature == nil || quorumCertificate.Signers == nil || quorumCertificate.Signers.Len() == 0 {
		return ErrInvalidQuorumCertificate
	}
	h.QuorumCertificate = quorumCertificate.Copy()
	return nil
}

// WriteActivityProof writes the extra-data field of a block header with given activity proof
func WriteActivityProof(h *Header, proof AggregateSignature, round uint64) {
	// if the proof is empty, do not copy anything
	if proof.Signature == nil || proof.Signers == nil || proof.Signers.Len() == 0 {
		return
	}
	h.ActivityProof = proof.Copy()
	h.ActivityProofRound = round
}
