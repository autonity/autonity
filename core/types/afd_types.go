package types

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/rlp"
	"io"
	"math/big"
)

type ProofType uint64
const (
	InnocentProof ProofType = iota
	ChallengeProof
)

type Rule uint8
const (
	PN Rule = iota
	PO
	PVN
	PVO
	C
	GarbageMessage		// message was signed by valid member, but it cannot be decoded.
	Equivocation        // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
	InvalidProposal     // The value proposed by proposer cannot pass the blockchain's validation.
	InvalidProposer     // A proposal sent from none proposer nodes of the committee.
	InvalidMsgSender    // Msg is not from committee member of the height.
)

// The proof used by accountability precompiled contract to validate the proof of innocent or misbehavior.
// Since precompiled contract take raw bytes as input, so it should be rlp encoded into bytes before client send the
// proof into autonity contract.
type RawProof struct {
	Rule       uint8       // rule id.
	Message    []byte      // the raw rlp encoded msg to be considered as suspicious one
	Evidence   [][]byte    // the raw rlp encoded msgs as proof of innocent or misbehavior.
}

func (p *RawProof) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{p.Rule, p.Message, p.Evidence})
}

func (p *RawProof) DecodeRLP(s *rlp.Stream) error {
	var proof struct {
		Rule uint8
		Message []byte
		Evidence [][]byte
	}
	if err := s.Decode(&proof); err != nil {
		return err
	}

	p.Rule, p.Message, p.Evidence = proof.Rule, proof.Message, proof.Evidence
	return nil
}

// RawProof to be decode into Proof for processing.
type Proof struct {
	Rule       uint8
	Message    ConsensusMessage   // the msg to be considered as suspicious one
	Evidence   []ConsensusMessage // the msgs as proof of innocent or misbehavior.
}

// OnChainProof to be stored by autonity contract for on-chain proof management.
type OnChainProof struct {
	// identities to address an unique proof on contract.
	Height *big.Int
	Round uint64
	MsgType uint64
	Sender common.Address
	Rule uint8

	// rlp enoded bytes for struct RawProof object.
	RawProofBytes []byte
}

// event to submit proofs via standard transaction.
type SubmitProofEvent struct {
	Proofs []OnChainProof
	Type ProofType
}