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
