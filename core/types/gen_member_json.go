// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto/blst"
)

var _ = (*committeeMemberMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (c CommitteeMember) MarshalJSON() ([]byte, error) {
	type CommitteeMember struct {
		Address           common.Address `json:"address"            gencodec:"required"       abi:"addr"`
		VotingPower       *hexutil.Big   `json:"votingPower"        gencodec:"required"`
		ConsensusKeyBytes hexutil.Bytes  `json:"consensusKey"       gencodec:"required"       abi:"consensusKey"`
		ConsensusKey      blst.PublicKey `json:"-" rlp:"-"`
		Index             uint64         `json:"-" rlp:"-"`
	}
	var enc CommitteeMember
	enc.Address = c.Address
	enc.VotingPower = (*hexutil.Big)(c.VotingPower)
	enc.ConsensusKeyBytes = c.ConsensusKeyBytes
	enc.ConsensusKey = c.ConsensusKey
	enc.Index = c.Index
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (c *CommitteeMember) UnmarshalJSON(input []byte) error {
	type CommitteeMember struct {
		Address           *common.Address `json:"address"            gencodec:"required"       abi:"addr"`
		VotingPower       *hexutil.Big    `json:"votingPower"        gencodec:"required"`
		ConsensusKeyBytes *hexutil.Bytes  `json:"consensusKey"       gencodec:"required"       abi:"consensusKey"`
		ConsensusKey      blst.PublicKey  `json:"-" rlp:"-"`
		Index             *uint64         `json:"-" rlp:"-"`
	}
	var dec CommitteeMember
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Address == nil {
		return errors.New("missing required field 'address' for CommitteeMember")
	}
	c.Address = *dec.Address
	if dec.VotingPower == nil {
		return errors.New("missing required field 'votingPower' for CommitteeMember")
	}
	c.VotingPower = (*big.Int)(dec.VotingPower)
	if dec.ConsensusKeyBytes == nil {
		return errors.New("missing required field 'consensusKey' for CommitteeMember")
	}
	c.ConsensusKeyBytes = *dec.ConsensusKeyBytes
	if dec.ConsensusKey != nil {
		c.ConsensusKey = dec.ConsensusKey
	}
	if dec.Index != nil {
		c.Index = *dec.Index
	}
	return nil
}