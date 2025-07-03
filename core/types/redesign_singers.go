package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	"github.com/autonity/autonity/rlp"
)

type Signers struct {
	SignatureCounts []uint16 // support up to 65535 committee members

	// these fields are not serialized, but instead computed at preValidate steps
	committeeSize int              `rlp:"-"`
	length        int              `rlp:"-"` // number of distinct signers
	powers        map[int]*big.Int `rlp:"-"`
	power         *big.Int         `rlp:"-"` // aggregated power of all senders

	// auxiliary data structures flags
	validated     bool `rlp:"-"` // if true --> Bits and Coefficients have correct length + committeeSize and length is assigned
	powerAssigned bool `rlp:"-"` // if true --> powers and power assigned
}

func NewSigners(committeeSize int) *Signers {
	if committeeSize > maxUint16 {
		panic("Unsupported committee size")
	}
	return &Signers{
		SignatureCounts: make([]uint16, committeeSize),
		committeeSize:   committeeSize,
		powers:          make(map[int]*big.Int),
		power:           new(big.Int),
		length:          0,
		validated:       true,
		powerAssigned:   true, // when we are locally creating a sender info, we are ok with power being 0 initially
	}
}

func (s *Signers) EncodeRLP(w io.Writer) error {
	payload := s.encodeVarintPayload()
	return rlp.Encode(w, payload)
}

func (s *Signers) ForEachDistinctSigner(callback func(signerIndex int)) {
	for i, count := range s.SignatureCounts {
		if count > 0 {
			callback(i)
		}
	}
}

func (s *Signers) encodeVarintPayload() []byte {
	var buf bytes.Buffer
	buf.Grow(s.length * 2) // estimate size based on number of signers

	varintBuffer := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(varintBuffer, uint64(s.committeeSize))
	buf.Write(varintBuffer[:n])
	for index, count := range s.SignatureCounts {
		if count == 0 {
			continue
		}
		// todo: review should we use u16
		n := binary.PutUvarint(varintBuffer, uint64(index))
		buf.Write(varintBuffer[:n])

		n = binary.PutUvarint(varintBuffer, uint64(count))
		buf.Write(varintBuffer[:n])
	}
	return buf.Bytes()
}

func (s *Signers) DecodeRLP(stream *rlp.Stream) error {
	payload, err := stream.Bytes()
	if err != nil {
		panic(err)
	}
	return s.decodeVarintPayload(payload)
}

func (s *Signers) decodeVarintPayload(payload []byte) error {
	reader := bytes.NewReader(payload)
	committeeSize, err := binary.ReadUvarint(reader)
	if err != nil {
		return fmt.Errorf("failed to read committee size: %w", err)
	}
	s.committeeSize = int(committeeSize)
	s.SignatureCounts = make([]uint16, s.committeeSize)
	s.length = 0
	for reader.Len() > 0 {
		index, err := binary.ReadUvarint(reader)
		if err != nil {
			return fmt.Errorf("error reading index: %w", err)
		}
		if index >= uint64(s.committeeSize) {
			return fmt.Errorf("index out of bounds: %d >= %d", index, s.committeeSize)
		}

		count, err := binary.ReadUvarint(reader)
		if err != nil {
			return fmt.Errorf("error reading count: %w", err)
		}
		if count > uint64(s.committeeSize) {
			return fmt.Errorf("count out of bounds: %d > %d", count, s.committeeSize)
		}

		s.SignatureCounts[index] = uint16(count)
		s.length++
	}
	return nil
}

func (s *Signers) Increment(member *CommitteeMember) {
	if !s.validated {
		panic("Attempting to use un validated signer information")
	}
	if !s.powerAssigned {
		panic("Power has not been assigned")
	}
	index := int(member.Index)
	if index >= s.committeeSize {
		panic("Invalid validator index")
	}
	if int(s.SignatureCounts[index]) >= s.committeeSize {
		panic("Signature count exceeds committee size")
	}
	if s.SignatureCounts[index] == 0 {
		s.length++
		s.powers[index] = member.VotingPower
		s.power.Add(s.power, member.VotingPower)
	}
	s.SignatureCounts[index]++
}

func (s *Signers) Merge(other *Signers) {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}
	for i := 0; i < s.committeeSize; i++ {
		if other.SignatureCounts[i] == 0 {
			continue
		}
		if int(s.SignatureCounts[i]+other.SignatureCounts[i]) > s.committeeSize {
			panic("signature count exceeds committee size")
		}
		// new signer, add length and update power contribution
		if s.SignatureCounts[i] == 0 && other.SignatureCounts[i] > 0 {
			s.length++
			s.powers[i] = other.powers[i]
			s.power.Add(s.power, other.power)
		}
		s.SignatureCounts[i] += other.SignatureCounts[i]
	}
}

func (s *Signers) Get(index int) uint16 {
	if !s.validated {
		panic("Attempting to use un validated signer information")
	}
	if index >= s.committeeSize {
		panic("Invalid validator index")
	}
	return s.SignatureCounts[index]
}

func (s *Signers) SanityCheck() error {
	if s.SignatureCounts == nil {
		return ErrNilSigners
	}
	if len(s.SignatureCounts) > maxAllowedSigners {
		return fmt.Errorf("invalid SignatureCounts length: %d", len(s.SignatureCounts))
	}
	return nil
}

func (s *Signers) Contains(index int) bool {
	if !s.validated {
		panic("Attempting to use un validated signer information")
	}
	if index >= s.committeeSize {
		panic("Invalid validator index")
	}
	return s.SignatureCounts[index] > 0
}

func (s *Signers) Validate(committeeSize int) error {
	countNonZero, err := s.validate(committeeSize)
	if err != nil {
		return err
	}
	s.committeeSize = committeeSize
	s.length = countNonZero
	s.validated = true
	return nil
}

func (s *Signers) validate(committeeSize int) (int, error) {
	if s.SignatureCounts == nil {
		return 0, ErrNilSigners
	}
	if len(s.SignatureCounts) != committeeSize {
		return 0, ErrWrongSizeSigners
	}
	countNonZero := 0
	sum := 0
	for _, count := range s.SignatureCounts {
		if count > 0 {
			countNonZero++
			sum += int(count)
		}
		if int(count) > committeeSize {
			return 0, ErrInvalidCoefficient
		}
	}
	if sum == 0 {
		return 0, ErrEmptySigners
	}
	if countNonZero == 1 && sum != 1 {
		return 0, ErrInvalidSingleSig
	}
	return countNonZero, nil
}

func (s *Signers) Power() *big.Int {
	if !s.powerAssigned {
		panic("Power has not been assigned")
	}
	return s.power
}

func (s *Signers) RespectsBoundaries(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}
	for i := 0; i < s.committeeSize; i++ {
		if int(s.SignatureCounts[i]+other.SignatureCounts[i]) > s.committeeSize {
			return false
		}
	}
	return true
}

func (s *Signers) AddsInformation(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}
	for i := 0; i < s.committeeSize; i++ {
		if byte(s.SignatureCounts[i]) == noSignature && byte(other.SignatureCounts[i]) != noSignature {
			return true
		}
	}
	return false
}

func (s *Signers) CanMergeSimple(other *Signers) bool {
	if err := safetyCheck(s, other); err != nil {
		panic(err.Error())
	}
	//todo(review) : if one signature has multiple counts but others not it should be still mergeable,
	// though in current implementation that scenario is not possible
	for i := 0; i < s.committeeSize; i++ {
		if byte(s.SignatureCounts[i]+other.SignatureCounts[i]) > oneSignature {
			return false
		}
	}
	return true
}

func (s *Signers) AssignPower(powers map[int]*big.Int, power *big.Int) {
	s.powers = powers
	s.power = power
	s.powerAssigned = true
}

func (s *Signers) Copy() *Signers {
	var powers map[int]*big.Int
	if s.powers != nil {
		powers = make(map[int]*big.Int, len(s.powers))
		for idx, power := range s.powers {
			powers[idx] = power
		}
	}
	signatureCounts := make([]uint16, len(s.SignatureCounts))
	copy(signatureCounts, s.SignatureCounts)

	copySigner := &Signers{
		SignatureCounts: signatureCounts,
		committeeSize:   s.committeeSize,
		length:          s.length,
		powers:          powers,
		validated:       s.validated,
		powerAssigned:   s.powerAssigned,
	}
	if s.power != nil {
		copySigner.power = new(big.Int).Set(s.power)
	}
	return copySigner

}

func (s *Signers) flatten() []int {
	sum := 0
	for _, count := range s.SignatureCounts {
		sum += int(count)
	}
	indexes := make([]int, 0, sum)
	for i, count := range s.SignatureCounts {
		for j := 0; j < int(count); j++ {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (s *Signers) Flatten() []int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.flatten()
}

func (s *Signers) flattenUniq() []int {
	indexes := make([]int, 0, s.length)
	for i, count := range s.SignatureCounts {
		if count > 0 {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

func (s *Signers) FlattenUniq() []int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.flattenUniq()
}

func (s *Signers) Len() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.length
}

func (s *Signers) IsComplex() bool {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	for _, count := range s.SignatureCounts {
		if count > 1 {
			return true
		}
	}
	return false
}

func (s *Signers) Powers() map[int]*big.Int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.powers
}

func (s *Signers) CommitteeSize() int {
	if !s.validated {
		panic("Using un-validated signers information")
	}
	return s.committeeSize
}

func (s *Signers) String() string {
	return fmt.Sprintf("SignatureCounts: %v, power: %v, validated: %v, powerAssigned: %v", s.SignatureCounts, s.power, s.validated, s.powerAssigned)
}

func safetyCheck(first *Signers, second *Signers) error {
	if !first.validated || !second.validated {
		return ErrNotValidated
	}
	if first.committeeSize != second.committeeSize {
		return ErrDifferentSize
	}
	return nil
}
