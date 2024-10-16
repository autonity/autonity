package types

import (
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

	inmemoryAddresses  = 500 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)
)

// BFTFilteredHeader returns a filtered header which some information (like proposerSeal, quorumCertificate)
// are clean to fulfill the BFT hash rules.
func BFTFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	if !keepSeal {
		newHeader.ProposerSeal = []byte{}
	}
	newHeader.QuorumCertificate = nil
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
