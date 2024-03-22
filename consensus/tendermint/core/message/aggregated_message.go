package message

import (
	"fmt"
	"math/big"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/rlp"
)

// sent on the wire
type extAggregateVote struct {
	Code      uint8
	Round     uint64
	Height    uint64
	Value     common.Hash
	Senders   *types.SendersInfo
	Signature *blst.BlsSignature
}

type AggregatePrevote struct {
	value common.Hash
	aggregateMsg
}

func (ap *AggregatePrevote) Code() uint8 {
	return AggregatePrevoteCode
}

func (ap *AggregatePrevote) Value() common.Hash {
	return ap.value
}
func (ap *AggregatePrevote) String() string {
	return fmt.Sprintf("{r:  %v, h: %v , power: %v, Code: %v, value: %v}",
		ap.round, ap.height, ap.power, ap.Code(), ap.value)
}

func (ap *AggregatePrevote) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}

	encoded := &extAggregateVote{}
	if err := rlp.DecodeBytes(payload, encoded); err != nil {
		return err
	}
	if encoded.Code != AggregatePrevoteCode {
		return constants.ErrInvalidMessage
	}
	if encoded.Signature == nil {
		return constants.ErrInvalidMessage
	}
	if encoded.Height == 0 {
		return constants.ErrInvalidMessage
	}
	if encoded.Round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	if encoded.Senders == nil || encoded.Senders.Bits == nil || encoded.Senders.Len() == 0 {
		return constants.ErrInvalidMessage
	}
	ap.height = encoded.Height
	ap.round = int64(encoded.Round)
	ap.value = encoded.Value
	ap.signature = encoded.Signature
	ap.senders = encoded.Senders
	// note: code for signature input is still prevote (this is correct)
	ap.payload = payload
	// precompute hash and signature hash
	signaturePayload, _ := rlp.EncodeToBytes([]any{PrevoteCode, encoded.Round, encoded.Height, encoded.Value})
	ap.signatureInput = crypto.Hash(signaturePayload)
	ap.hash = crypto.Hash(payload)
	return nil
}

type AggregatePrecommit struct {
	value common.Hash
	aggregateMsg
}

func (ap *AggregatePrecommit) Code() uint8 {
	return AggregatePrecommitCode
}

func (ap *AggregatePrecommit) Value() common.Hash {
	return ap.value
}
func (ap *AggregatePrecommit) String() string {
	return fmt.Sprintf("{r:  %v, h: %v , power: %v, Code: %v, value: %v}",
		ap.round, ap.height, ap.power, ap.Code(), ap.value)
}

func (ap *AggregatePrecommit) DecodeRLP(s *rlp.Stream) error {
	payload, err := s.Raw()
	if err != nil {
		return err
	}

	encoded := &extAggregateVote{}
	if err := rlp.DecodeBytes(payload, encoded); err != nil {
		return err
	}
	if encoded.Code != AggregatePrecommitCode {
		return constants.ErrInvalidMessage
	}
	if encoded.Signature == nil {
		return constants.ErrInvalidMessage
	}
	if encoded.Height == 0 {
		return constants.ErrInvalidMessage
	}
	if encoded.Round > constants.MaxRound {
		return constants.ErrInvalidMessage
	}
	if encoded.Senders == nil || encoded.Senders.Bits == nil || encoded.Senders.Len() == 0 {
		return constants.ErrInvalidMessage
	}
	ap.height = encoded.Height
	ap.round = int64(encoded.Round)
	ap.value = encoded.Value
	ap.signature = encoded.Signature
	ap.senders = encoded.Senders
	// note: code for signature input is still precommit (this is correct)
	ap.payload = payload
	// precompute hash and signature hash
	signaturePayload, _ := rlp.EncodeToBytes([]any{PrecommitCode, encoded.Round, encoded.Height, encoded.Value})
	ap.signatureInput = crypto.Hash(signaturePayload)
	ap.hash = crypto.Hash(payload)
	return nil
}

func NewAggregatePrevote(votes []Msg, header *types.Header) *AggregatePrevote {
	return NewAggregateVote[AggregatePrevote](votes, header)
}

func NewAggregatePrecommit(votes []Msg, header *types.Header) *AggregatePrecommit {
	return NewAggregateVote[AggregatePrecommit](votes, header)
}

// NOTE: this function assumes that:
// 1. all votes are for the same signature input (code,h,r,value)
// 2. all votes have previously been cryptographically verified
func NewAggregateVote[
	E AggregatePrevote | AggregatePrecommit,
	PE interface {
		*E
		Msg
	}](votes []Msg, header *types.Header) *E {
	code := PE(new(E)).Code()

	// TODO(lorenzo) aggregates for different heights might have different len(senders.bits)
	senders := types.NewSendersInfo(len(header.Committee))

	// compute new aggregated signature and related sender information
	var signatures []blst.Signature
	for _, vote := range votes {
		switch m := vote.(type) {
		case *Propose, *Prevote, *Precommit:
			senders.Increment(header, m.(IndividualMsg).SenderIndex())
		case *AggregatePrevote, *AggregatePrecommit:
			senders.Merge(m.(AggregateMsg).Senders())
		}
		signatures = append(signatures, vote.Signature())
	}
	aggregatedSignature := blst.Aggregate(signatures)

	// use votes[0] as a set representative
	vote := votes[0]
	h := vote.H()
	r := vote.R()
	value := vote.Value()
	signatureInput := vote.SignatureInput()

	payload, _ := rlp.EncodeToBytes(extAggregateVote{
		Code:      code,
		Round:     uint64(r),
		Height:    h,
		Value:     value,
		Senders:   senders,
		Signature: aggregatedSignature.(*blst.BlsSignature),
	})

	aggregateVote := E{
		value: value,
		aggregateMsg: aggregateMsg{
			senders: senders,
			base: base{
				height:         h,
				round:          r,
				signatureInput: signatureInput,
				signature:      aggregatedSignature,
				payload:        payload,
				hash:           crypto.Hash(payload),
				verified:       true,            // verified due to all votes being verified
				power:          senders.Power(), // aggregated power
				//TODO(lorenzo) missing sender key, maybe it is better to add it even if not strictly needed?
			},
		},
	}
	return &aggregateVote
}

type aggregateMsg struct {
	senders *types.SendersInfo
	base
}

func (am *aggregateMsg) Senders() *types.SendersInfo {
	return am.senders
}

func (am *aggregateMsg) PreValidate(header *types.Header) error {
	if am.senders.Len() != len(header.Committee) {
		return ErrInvalidSenders
	}
	// compute aggregated key
	indexes := am.senders.Flatten()
	keys := make([][]byte, len(indexes))
	for i, index := range indexes {
		keys[i] = header.Committee[index].ConsensusKeyBytes
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		panic("Error while aggregating keys from committee: " + err.Error())
	}

	indexesUniq := am.senders.FlattenUniq()
	addresses := make(map[int]common.Address)
	powers := make(map[int]*big.Int)

	for _, index := range indexesUniq {
		member := header.Committee[index]
		addresses[index] = member.Address
		powers[index] = member.VotingPower
	}

	am.senders.SetPowers(powers)
	am.senders.SetAddresses(addresses)

	am.power = am.senders.Power()
	am.senderKey = aggregatedKey
	return nil
}

func (am *aggregateMsg) Validate() error {
	if am.verified {
		return nil
	}

	if valid := am.signature.Verify(am.senderKey, am.signatureInput[:]); !valid {
		return ErrBadSignature
	}
	am.verified = true

	return nil
}
