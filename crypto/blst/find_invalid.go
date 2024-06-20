package blst

import (
	"fmt"
	"sync"
)

func FindInvalid(signatures []Signature, publicKeys []PublicKey, msg [32]byte) []uint {
	if len(signatures) != len(publicKeys) {
		panic(fmt.Sprintf("invalid arguments, length mismatch. signatures: %d, public keys: %d\n", len(signatures), len(publicKeys)))
	}
	if len(signatures) == 0 {
		return nil
	}
	return findInvalid(signatures, publicKeys, msg, 0, uint(len(signatures)))
}

func findInvalid(
	signatures []Signature,
	pks []PublicKey,
	msg [32]byte,
	start, end uint,
) []uint {
	// if we have two elements, no point in further splitting since we need to do
	// two verifications anyway
	if end-start <= 2 {

		var ret []uint

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

		return ret
	}

	pivot := start + ((end - start) / 2)

	var leftInvalid []uint
	var rightInvalid []uint

	wg := sync.WaitGroup{}

	//left
	wg.Add(1)

	go func() {
		verified := FastAggregateVerifyBatch(signatures[start:pivot], pks[start:pivot], msg)

		if !verified {
			leftInvalid = findInvalid(signatures, pks, msg, start, pivot)
		}
		wg.Done()

	}()

	//right
	wg.Add(1)

	go func() {
		verified := FastAggregateVerifyBatch(signatures[pivot:end], pks[pivot:end], msg)

		if !verified {

			rightInvalid = findInvalid(signatures, pks, msg, pivot, end)

		}
		wg.Done()
	}()

	wg.Wait()

	return append(leftInvalid, rightInvalid...)
}

// NOTE: functions below currently not used, they use AggregateVerify instead of FastAggregateVerify. They are used in benchmarks.
// 			IMPORTANT: If we start to use them we need to make sure that we always pass distinct messages to AggregateVerify

func FindInvalidSignatures(signatures []Signature, pks []PublicKey, msgs [][32]byte) ([]uint, error) {
	if len(signatures) != len(pks) && len(pks) != len(msgs) {
		return nil, fmt.Errorf("invalid arguments, length mismatch")
	}
	if len(signatures) == 0 {
		return nil, nil
	}
	return findInvalidSignaturesRecursive(signatures, pks, msgs, 0, uint(len(signatures)))
}

func findInvalidSignaturesRecursive(
	signatures []Signature,
	pks []PublicKey,
	msgs [][32]byte,
	start, end uint,
) ([]uint, error) {

	// if we have two elements, no point in further splitting since we need to do
	// two verifications anyway
	if end-start <= 2 {

		var ret []uint

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

	}
	pivot := start + ((end - start) / 2)

	var leftInvalid []uint
	var rightInvalid []uint

	wg := sync.WaitGroup{}

	//left

	var leftErr error
	var rightErr error

	wg.Add(1)

	go func() {
		aggSig := AggregateSignatures(signatures[start:pivot])
		// TODO: make sure that msgs are all distinct
		verified := aggSig.AggregateVerify(pks[start:pivot], msgs[start:pivot])

		if !verified {
			leftInvalid, leftErr = findInvalidSignaturesRecursive(signatures, pks, msgs, start, pivot)
		}
		wg.Done()

	}()

	//right
	wg.Add(1)

	go func() {
		aggSig := AggregateSignatures(signatures[pivot:end])
		// TODO: make sure that messages are all distinct
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
