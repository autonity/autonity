package types

//
//import (
//	"bytes"
//	"encoding/binary"
//	"fmt"
//	"io"
//	"math/big"
//	"sort"
//
//	"github.com/autonity/autonity/rlp"
//)
//
//type Signers struct {
//	// SignatureCounts maps a validator's index to their number of signatures.
//	SignatureCounts map[int]uint16
//
//	committeeSize int              `rlp:"-"`
//	powers        map[int]*big.Int `rlp:"-"`
//	power         *big.Int         `rlp:"-"`
//
//	// auxiliary data structures flags
//	validated     bool `rlp:"-"` // if true --> committeeSize is assigned and the data is sensible.
//	powerAssigned bool `rlp:"-"` // if true --> powers and power have been assigned.
//}
//
//// NewSigners creates a new, empty Signers struct for a given committee size.
//func NewSigners(committeeSize int) *Signers {
//	if committeeSize > maxUint16 {
//		panic("Unsupported committee size")
//	}
//	return &Signers{
//		SignatureCounts: make(map[int]uint16),
//		committeeSize:   committeeSize,
//		powers:          make(map[int]*big.Int),
//		power:           new(big.Int),
//		validated:       true,
//		powerAssigned:   true, // Locally created instances are assumed to have power assigned (initially zero).
//	}
//}
//
//// EncodeRLP encodes the Signers struct into RLP format.
//// It uses a variable-integer encoding for a compact representation.
//func (s *Signers) EncodeRLP(w io.Writer) error {
//	payload := s.encodeVarintPayload()
//	return rlp.Encode(w, payload)
//}
//
//// ForEachDistinctSigner iterates over each signer and calls the callback function.
//func (s *Signers) ForEachDistinctSigner(callback func(signerIndex int)) {
//	for index, count := range s.SignatureCounts {
//		if count > 0 {
//			callback(index)
//		}
//	}
//}
//
//// encodeVarintPayload creates a compact byte representation of the signers.
//// Format: [committeeSize, index1, count1, index2, count2, ...]
//func (s *Signers) encodeVarintPayload() []byte {
//	var buf bytes.Buffer
//	buf.Grow(len(s.SignatureCounts) * 4) // Estimate size.
//
//	varintBuffer := make([]byte, binary.MaxVarintLen64)
//
//	// Write committee size.
//	n := binary.PutUvarint(varintBuffer, uint64(s.committeeSize))
//	buf.Write(varintBuffer[:n])
//
//	// sort indices to always create a deterministic payload for same content
//	indices := make([]int, 0, len(s.SignatureCounts))
//	for index := range s.SignatureCounts {
//		indices = append(indices, index)
//	}
//	sort.Ints(indices)
//
//	// Write each signer's index and count.
//	for _, i := range indices {
//		count := s.SignatureCounts[i]
//		if count == 0 {
//			continue
//		}
//		n = binary.PutUvarint(varintBuffer, uint64(i))
//		buf.Write(varintBuffer[:n])
//
//		n = binary.PutUvarint(varintBuffer, uint64(count))
//		buf.Write(varintBuffer[:n])
//	}
//	return buf.Bytes()
//}
//
//// DecodeRLP decodes RLP-encoded bytes into the Signers struct.
//func (s *Signers) DecodeRLP(stream *rlp.Stream) error {
//	payload, err := stream.Bytes()
//	if err != nil {
//		return err
//	}
//	return s.decodeVarintPayload(payload)
//}
//
//// decodeVarintPayload parses the compact byte representation into the Signers struct.
//func (s *Signers) decodeVarintPayload(payload []byte) error {
//	reader := bytes.NewReader(payload)
//	committeeSize, err := binary.ReadUvarint(reader)
//	if err != nil {
//		return fmt.Errorf("failed to read committee size: %w", err)
//	}
//
//	s.committeeSize = int(committeeSize)
//	s.SignatureCounts = make(map[int]uint16, s.committeeSize)
//
//	for reader.Len() > 0 {
//		index, err := binary.ReadUvarint(reader)
//		if err != nil {
//			return fmt.Errorf("error reading index: %w", err)
//		}
//		if index >= uint64(s.committeeSize) {
//			return fmt.Errorf("index out of bounds: %d >= %d", index, s.committeeSize)
//		}
//
//		count, err := binary.ReadUvarint(reader)
//		if err != nil {
//			return fmt.Errorf("error reading count: %w", err)
//		}
//		if count > uint64(maxUint16) {
//			return fmt.Errorf("count out of bounds: %d", count)
//		}
//
//		if count > 0 {
//			s.SignatureCounts[int(index)] = uint16(count)
//		}
//	}
//	return nil
//}
//
//// Increment records a signature from a committee member.
//func (s *Signers) Increment(member *CommitteeMember) {
//	if !s.validated {
//		panic("Attempting to use un-validated signer information")
//	}
//	if !s.powerAssigned {
//		panic("Power has not been assigned")
//	}
//	index := int(member.Index)
//	if index >= s.committeeSize {
//		panic("Invalid validator index")
//	}
//
//	// If this is the first signature from this member, update power.
//	if _, exists := s.SignatureCounts[index]; !exists {
//		s.powers[index] = member.VotingPower
//		s.power.Add(s.power, member.VotingPower)
//	}
//	s.SignatureCounts[index]++
//}
//
//// Merge combines another Signers struct into this one.
//func (s *Signers) Merge(other *Signers) {
//	if err := safetyCheck(s, other); err != nil {
//		panic(err.Error())
//	}
//
//	for otherIndex, otherCount := range other.SignatureCounts {
//		currentCount := s.SignatureCounts[otherIndex]
//		if int(currentCount+otherCount) > s.committeeSize {
//			panic("signature count exceeds committee size")
//		}
//
//		// If the signer from 'other' is new to 's', add their power.
//		if currentCount == 0 && otherCount > 0 {
//			if power, ok := other.powers[otherIndex]; ok {
//				s.powers[otherIndex] = power
//				s.power.Add(s.power, power)
//			}
//		}
//		s.SignatureCounts[otherIndex] += otherCount
//	}
//}
//
//// Set explicitly sets the signature count for a validator.
//func (s *Signers) Set(validatorIndex int, value uint16) {
//	s.SignatureCounts[validatorIndex] = value
//}
//
//// SanityCheck performs basic validation on the Signers struct.
//func (s *Signers) SanityCheck() error {
//	if s.SignatureCounts == nil {
//		return ErrNilSigners
//	}
//	if s.committeeSize > maxAllowedSigners {
//		return fmt.Errorf("invalid committee size: %d", s.committeeSize)
//	}
//	return nil
//}
//
//// Contains checks if a validator at a given index has signed.
//func (s *Signers) Contains(index int) bool {
//	if !s.validated {
//		panic("Attempting to use un-validated signer information")
//	}
//	if index >= s.committeeSize {
//		panic("Invalid validator index")
//	}
//	count, _ := s.SignatureCounts[index]
//	if count > 0 {
//		return true
//	}
//	return false
//}
//
//// Validate checks if the Signers data is consistent with a given committee size.
//func (s *Signers) Validate(committeeSize int) error {
//	if err := s.validate(committeeSize); err != nil {
//		return err
//	}
//	s.committeeSize = committeeSize
//	s.validated = true
//	return nil
//}
//
//// validate performs internal consistency checks.
//func (s *Signers) validate(committeeSize int) error {
//	if s.SignatureCounts == nil {
//		return ErrNilSigners
//	}
//
//	sum := 0
//	countNonZero := 0
//	for index, count := range s.SignatureCounts {
//		if index >= committeeSize {
//			return fmt.Errorf("index out of bounds: %d >= %d", index, committeeSize)
//		}
//		if count > 0 {
//			countNonZero++
//			sum += int(count)
//		}
//		if int(count) > committeeSize {
//			return ErrInvalidCoefficient
//		}
//	}
//
//	if sum == 0 {
//		return ErrEmptySigners
//	}
//
//	if countNonZero == 1 && sum != 1 {
//		return ErrInvalidSingleSig
//	}
//
//	return nil
//}
//
//// Power returns the total aggregated voting power of all signers.
//func (s *Signers) Power() *big.Int {
//	if !s.powerAssigned {
//		panic("Power has not been assigned")
//	}
//	return s.power
//}
//
//// Copy creates a deep copy of the Signers struct.
//// This operation is now much cheaper due to the sparse map representation.
//func (s *Signers) Copy() *Signers {
//	signatureCounts := make(map[int]uint16, len(s.SignatureCounts))
//	for idx, count := range s.SignatureCounts {
//		if count > 0 {
//			signatureCounts[idx] = count
//		}
//	}
//
//	// Copy the powers map.
//	var powers map[int]*big.Int
//	if s.powers != nil {
//		powers = make(map[int]*big.Int, len(s.powers))
//		for idx, power := range s.powers {
//			powers[idx] = new(big.Int).Set(power)
//		}
//	}
//
//	copySigner := &Signers{
//		SignatureCounts: signatureCounts,
//		committeeSize:   s.committeeSize,
//		powers:          powers,
//		validated:       s.validated,
//		powerAssigned:   s.powerAssigned,
//	}
//
//	if s.power != nil {
//		copySigner.power = new(big.Int).Set(s.power)
//	}
//
//	return copySigner
//}
//
//// Flatten returns a slice containing the index of each signature.
//// If a member signed N times, their index appears N times.
//func (s *Signers) Flatten() []int {
//	if !s.validated {
//		panic("Using un-validated signers information")
//	}
//	return s.flatten()
//}
//
//func (s *Signers) flatten() []int{
//	sum := 0
//	for _, count := range s.SignatureCounts {
//		sum += int(count)
//	}
//	indexes := make([]int, 0, sum)
//	for i, count := range s.SignatureCounts {
//		for j := 0; j < int(count); j++ {
//			indexes = append(indexes, i)
//		}
//	}
//	return indexes
//
//}
//
//// FlattenUniq returns a slice of unique signer indexes.
//func (s *Signers) FlattenUniq() []int {
//	if !s.validated {
//		panic("Using un-validated signers information")
//	}
//	indexes := make([]int, 0, len(s.SignatureCounts))
//	for i, count := range s.SignatureCounts {
//		if count > 0 {
//			indexes = append(indexes, i)
//		}
//	}
//	return indexes
//}
//
//// Len returns the number of distinct signers.
//func (s *Signers) Len() int {
//	if !s.validated {
//		panic("Using un-validated signers information")
//	}
//	return len(s.SignatureCounts)
//}
//
//// IsComplex checks if any signer has more than one signature.
//func (s *Signers) IsComplex() bool {
//	if !s.validated {
//		panic("Using un-validated signers information")
//	}
//	for _, count := range s.SignatureCounts {
//		if count > 1 {
//			return true
//		}
//	}
//	return false
//}
//
//// Powers returns the map of signer indexes to their voting power.
//func (s *Signers) Powers() map[int]*big.Int {
//	if !s.validated {
//		panic("Using un-validated signers information")
//	}
//	return s.powers
//}
//
//// CommitteeSize returns the configured committee size.
//func (s *Signers) CommitteeSize() int {
//	if !s.validated {
//		panic("Using un-validated signers information")
//	}
//	return s.committeeSize
//}
//
//func (s *Signers) String() string {
//	return fmt.Sprintf("SignatureCounts: %v, power: %v, validated: %v, powerAssigned: %v", s.SignatureCounts, s.power, s.validated, s.powerAssigned)
//}
//
//// safetyCheck ensures two Signers structs can be compared or merged.
//func safetyCheck(first *Signers, second *Signers) error {
//	if !first.validated || !second.validated {
//		return ErrNotValidated
//	}
//	if first.committeeSize != second.committeeSize {
//		return ErrDifferentSize
//	}
//	return nil
//}
//
//// RespectsBoundaries checks if merging with another Signers would not exceed count limits.
//func (s *Signers) RespectsBoundaries(other *Signers) bool {
//	if err := safetyCheck(s, other); err != nil {
//		panic(err.Error())
//	}
//	for i, otherCount := range other.SignatureCounts {
//		if int(s.SignatureCounts[i]+otherCount) > s.committeeSize {
//			return false
//		}
//	}
//	return true
//}
//
//// AddsInformation checks if the other Signers struct contains signatures not present in this one.
//func (s *Signers) AddsInformation(other *Signers) bool {
//	if err := safetyCheck(s, other); err != nil {
//		panic(err.Error())
//	}
//	for i, otherCount := range other.SignatureCounts {
//		if otherCount > 0 && s.SignatureCounts[i] == noSignature {
//			return true
//		}
//	}
//	return false
//}
//
//// CanMergeSimple checks if two Signers can be merged without creating complex signatures (counts > 1).
//func (s *Signers) CanMergeSimple(other *Signers) bool {
//	if err := safetyCheck(s, other); err != nil {
//		panic(err.Error())
//	}
//	for i, otherCount := range other.SignatureCounts {
//		if s.SignatureCounts[i]+otherCount > oneSignature {
//			return false
//		}
//	}
//
//	// Also check s for existing complex signatures
//	for _, sCount := range s.SignatureCounts {
//		if sCount > oneSignature {
//			return false
//		}
//	}
//	return true
//}
//
//func (s *Signers) AssignPower(powers map[int]*big.Int, power *big.Int) {
//	s.powers = powers
//	s.power = power
//	s.powerAssigned = true
//}