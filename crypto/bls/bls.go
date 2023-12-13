package bls

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"

	"github.com/autonity/autonity/crypto/bls/blst"
	"github.com/autonity/autonity/crypto/bls/common"
	"github.com/pkg/errors"
)

// SecretKeyFromECDSAKey generate a BLS private key by sourcing from an ECDSA private key.
func SecretKeyFromECDSAKey(sk *ecdsa.PrivateKey) (SecretKey, error) {
	return blst.SecretKeyFromECDSAKey(sk.D.Bytes())
}

// SecretKeyFromBytes creates a BLS private key from a BigEndian byte slice.
func SecretKeyFromBytes(privKey []byte) (SecretKey, error) {
	return blst.SecretKeyFromBytes(privKey)
}

// SecretKeyFromBigNum takes in a big number string and creates a BLS private key.
func SecretKeyFromBigNum(s string) (SecretKey, error) {
	num := new(big.Int)
	num, ok := num.SetString(s, 10)
	if !ok {
		return nil, errors.New("could not set big int from string")
	}
	bts := num.Bytes()
	if len(bts) != 32 {
		return nil, errors.Errorf("provided big number string sets to a key unequal to 32 bytes: %d != 32", len(bts))
	}
	return SecretKeyFromBytes(bts)
}

// PublicKeyBytesFromString decode a BLS public key from a string.
func PublicKeyBytesFromString(pubKey string) ([]byte, error) {
	return blst.PublicKeyBytesFromString(pubKey)
}

// PublicKeyFromBytes creates a BLS public key from a  BigEndian byte slice.
func PublicKeyFromBytes(pubKey []byte) (PublicKey, error) {
	return blst.PublicKeyFromBytes(pubKey)
}

// SignatureFromBytes creates a BLS signature from a LittleEndian byte slice.
func SignatureFromBytes(sig []byte) (Signature, error) {
	return blst.SignatureFromBytes(sig)
}

// AggregatePublicKeys aggregates the provided raw public keys into a single key.
func AggregatePublicKeys(pubs [][]byte) (PublicKey, error) {
	return blst.AggregatePublicKeys(pubs)
}

// AggregateSignatures converts a list of signatures into a single, aggregated sig.
func AggregateSignatures(sigs []common.BLSSignature) common.BLSSignature {
	return blst.AggregateSignatures(sigs)
}

// VerifyMultipleSignatures verifies multiple signatures for distinct messages securely.
func VerifyMultipleSignatures(sigs [][]byte, msgs [][32]byte, pubKeys []common.BLSPublicKey) (bool, error) {
	return blst.VerifyMultipleSignatures(sigs, msgs, pubKeys)
}

// NewAggregateSignature creates a blank aggregate signature.
func NewAggregateSignature() common.BLSSignature {
	return blst.NewAggregateSignature()
}

// RandKey creates a new private key using a random input.
func RandKey() (common.BLSSecretKey, error) {
	return blst.RandKey()
}

func FindInvalidSignatures(signatures []common.BLSSignature, pks []common.BLSPublicKey, msgs [][32]byte) ([]int, error) {
	if len(signatures) != len(pks) && len(pks) != len(msgs) {
		return nil, fmt.Errorf("invalid arguments, length mismatch")
	}
	if len(signatures) == 0 {
		return nil, nil
	}
	return findInvalidSignaturesRecursive(signatures, pks, msgs, 0, len(signatures))
}

func FindFastInvalidSignatures(signatures []common.BLSSignature, pks []common.BLSPublicKey, msg [32]byte) ([]int, error) {
	if len(signatures) != len(pks) {
		return nil, fmt.Errorf("invalid arguments, length mismatch")
	}
	if len(signatures) == 0 {
		return nil, nil
	}
	return findFastInvalidSignaturesRecursive(signatures, pks, msg, 0, len(signatures))
}

func findInvalidSignaturesRecursive(
	signatures []common.BLSSignature,
	pks []common.BLSPublicKey,
	msgs [][32]byte,
	start, end int,
) ([]int, error) {

	// if we have two elements, no point in further splitting since we need to do
	// two verifications anyway
	if end-start <= 2 {

		var ret []int

		valid := signatures[start].Verify(pks[start], msgs[start][:])
		if !valid {
			ret = append(ret, start)
		}

		if end-start == 2 {
			valid := signatures[end-1].Verify(pks[end-1], msgs[end-1][:])
			if !valid {
				ret = append(ret, end-1)
			}
		}

		return ret, nil

	} else {
		pivot := start + ((end - start) / 2)

		var leftInvalid []int
		var rightInvalid []int

		wg := sync.WaitGroup{}

		//left

		var leftErr error
		var rightErr error

		wg.Add(1)

		go func() {
			//func() {
			aggSig := AggregateSignatures(signatures[start:pivot])
			verified := aggSig.AggregateVerify(pks[start:pivot], msgs[start:pivot])

			if !verified {
				leftInvalid, leftErr = findInvalidSignaturesRecursive(signatures, pks, msgs, start, pivot)
			}
			wg.Done()

		}()

		//right
		wg.Add(1)

		go func() {
			//func() {

			aggSig := AggregateSignatures(signatures[pivot:end])
			verified := aggSig.AggregateVerify(pks[pivot:end], msgs[pivot:end])

			if !verified {

				rightInvalid, rightErr = findInvalidSignaturesRecursive(signatures, pks, msgs, pivot, end)

			}
			wg.Done()
		}()

		wg.Wait()

		if leftErr != nil {
			return nil, fmt.Errorf("recursive check failed for %d %d", start, pivot)
		}
		if rightErr != nil {
			return nil, fmt.Errorf("recursive check failed for %d %d", pivot, end)
		}

		return append(leftInvalid, rightInvalid...), nil
	}
}

func findFastInvalidSignaturesRecursive(
	signatures []common.BLSSignature,
	pks []common.BLSPublicKey,
	msg [32]byte,
	start, end int,
) ([]int, error) {

	// if we have two elements, no point in further splitting since we need to do
	// two verifications anyway
	if end-start <= 2 {

		var ret []int

		valid := signatures[start].Verify(pks[start], msg[:])
		if !valid {
			ret = append(ret, start)
		}

		if end-start == 2 {
			valid := signatures[end-1].Verify(pks[end-1], msg[:])
			if !valid {
				ret = append(ret, end-1)
			}
		}

		return ret, nil

	} else {
		pivot := start + ((end - start) / 2)

		var leftInvalid []int
		var rightInvalid []int

		wg := sync.WaitGroup{}

		//left

		var leftErr error
		var rightErr error

		wg.Add(1)

		go func() {
			//func() {
			aggSig := AggregateSignatures(signatures[start:pivot])
			verified := aggSig.FastAggregateVerify(pks[start:pivot], msg)

			if !verified {
				leftInvalid, leftErr = findFastInvalidSignaturesRecursive(signatures, pks, msg, start, pivot)
			}
			wg.Done()

		}()

		//right
		wg.Add(1)

		go func() {
			//func() {

			aggSig := AggregateSignatures(signatures[pivot:end])
			verified := aggSig.FastAggregateVerify(pks[pivot:end], msg)

			if !verified {

				rightInvalid, rightErr = findFastInvalidSignaturesRecursive(signatures, pks, msg, pivot, end)

			}
			wg.Done()
		}()

		wg.Wait()

		if leftErr != nil {
			return nil, fmt.Errorf("recursive check failed for %d %d", start, pivot)
		}
		if rightErr != nil {
			return nil, fmt.Errorf("recursive check failed for %d %d", pivot, end)
		}

		return append(leftInvalid, rightInvalid...), nil
	}
}
