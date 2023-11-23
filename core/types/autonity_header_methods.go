package types

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto/bls"
	"github.com/autonity/autonity/log"
	"math/big"
)

func (h *Header) IsGenesis() bool {
	return h.Number.Uint64() == 0
}

// CommitteeMember returns the committee member having the given address or nil if there is none.
func (h *Header) CommitteeMember(address common.Address) *CommitteeMember {
	// if we are not on an epoch-header, crash
	if h.LastEpochBlock.Cmp(h.Number) != 0 {
		log.Crit("calling committee related function on a non-epoch header")
	}

	return h.Committee.CommitteeMember(address)
}

// AggregatedValidatorKey returns the aggregated validator public key of the committee members
func (h *Header) AggregatedValidatorKey() bls.PublicKey {
	// if we are not on an epoch-header, crash
	if h.LastEpochBlock.Cmp(h.Number) != 0 {
		log.Crit("calling committee related function on a non-epoch header")
	}
	return h.Committee.AggregatedValidatorKey()
}

// TotalVotingPower returns the total voting power contained in the committee
// for the block associated with this header.
func (h *Header) TotalVotingPower() *big.Int {
	// if we are not on an epoch-head, crash
	if h.LastEpochBlock.Cmp(h.Number) != 0 {
		log.Crit("calling committee related function on a non-epoch header")
	}
	return h.Committee.TotalVotingPower()
}
