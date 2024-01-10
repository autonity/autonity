package blst

import (
	mrand "math/rand"
	"testing"
)

// benchmarks

// benchmark bls public key aggregation
func benchmarkAggregatePK(n int, b *testing.B) {
	// print out parameters
	b.Logf("n: %d\n", n)

	// setup public keys
	var pks [][]byte
	for i := 0; i < n; i++ {
		sk, err := RandKey()
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey().Marshal()
		pks = append(pks, pk)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := AggregatePublicKeys(pks); err != nil {
			b.Fatal("failed pks aggregation: ", err)
		}
	}

}

func BenchmarkAggregatePK0_100(b *testing.B) { benchmarkAggregatePK(100, b) }
func BenchmarkAggregatePK1_100(b *testing.B) { benchmarkAggregatePK(100, b) }
func BenchmarkAggregatePK2_100(b *testing.B) { benchmarkAggregatePK(100, b) }
func BenchmarkAggregatePK0_200(b *testing.B) { benchmarkAggregatePK(200, b) }
func BenchmarkAggregatePK1_200(b *testing.B) { benchmarkAggregatePK(200, b) }
func BenchmarkAggregatePK2_200(b *testing.B) { benchmarkAggregatePK(200, b) }
func BenchmarkAggregatePK0_300(b *testing.B) { benchmarkAggregatePK(300, b) }
func BenchmarkAggregatePK1_300(b *testing.B) { benchmarkAggregatePK(300, b) }
func BenchmarkAggregatePK2_300(b *testing.B) { benchmarkAggregatePK(300, b) }

// benchmark bls signature aggregation
func benchmarkAggregateSigs(seed int64, n int, b *testing.B) {
	// print out parameters
	b.Logf("seed: %d, n: %d\n", seed, n)

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(seed)) //nolint

	// generate msg
	var msg [32]byte
	rand.Read(msg[:])

	// generate signatures over constant msg
	var sigs []Signature
	for i := 0; i < n; i++ {
		sk, err := RandKey()
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		sig := sk.Sign(msg[:])
		sigs = append(sigs, sig)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AggregateSignatures(sigs)
	}

}

func BenchmarkAggregateSigs0_3(b *testing.B)   { benchmarkAggregateSigs(0, 3, b) }
func BenchmarkAggregateSigs0_100(b *testing.B) { benchmarkAggregateSigs(0, 100, b) }
func BenchmarkAggregateSigs1_100(b *testing.B) { benchmarkAggregateSigs(1, 100, b) }
func BenchmarkAggregateSigs2_100(b *testing.B) { benchmarkAggregateSigs(2, 100, b) }
func BenchmarkAggregateSigs0_200(b *testing.B) { benchmarkAggregateSigs(0, 200, b) }
func BenchmarkAggregateSigs1_200(b *testing.B) { benchmarkAggregateSigs(1, 200, b) }
func BenchmarkAggregateSigs2_200(b *testing.B) { benchmarkAggregateSigs(2, 200, b) }
func BenchmarkAggregateSigs0_300(b *testing.B) { benchmarkAggregateSigs(0, 300, b) }
func BenchmarkAggregateSigs1_300(b *testing.B) { benchmarkAggregateSigs(1, 300, b) }
func BenchmarkAggregateSigs2_300(b *testing.B) { benchmarkAggregateSigs(2, 300, b) }

// benchmark simple signature verification
func BenchmarkSigVerify(b *testing.B) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	// generate msg
	var msg [32]byte
	rand.Read(msg[:])

	// generate signature over constant msg
	sk, err := RandKey()
	if err != nil {
		b.Fatal("Failed to generate random bls key: ", err)
	}
	sig := sk.Sign(msg[:])

	// start the actual verification benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := sig.Verify(sk.PublicKey(), msg[:])
		if !res {
			b.Fatal("failed signature verification")
		}
	}

}

// benchmark aggregated signature verification on 1 message
func BenchmarkSigVerifyAgg(b *testing.B) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	// generate msg
	var msg [32]byte
	rand.Read(msg[:])

	// generate signatures over constant msg
	var sigs []Signature
	var pks [][]byte
	for i := 0; i < 100; i++ {
		sk, err := RandKey()
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey().Marshal()
		pks = append(pks, pk)
		sig := sk.Sign(msg[:])
		sigs = append(sigs, sig)
	}

	sig := AggregateSignatures(sigs)
	aggPk, err := AggregatePublicKeys(pks)
	if err != nil {
		b.Fatal("failed pks aggregation: ", err)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := sig.Verify(aggPk, msg[:])
		if !res {
			b.Fatal("failed signature verification")
		}
	}
}

// benchmark aggregated signature verification on 3 messages (propose, prevote, precommit)
func BenchmarkAggregateVerify3(b *testing.B) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	// create step messages
	var propose [32]byte
	rand.Read(propose[:])
	var prevote [32]byte
	rand.Read(prevote[:])
	var precommit [32]byte
	rand.Read(precommit[:])

	// generate signature over propose
	proposeSk, err := RandKey()
	if err != nil {
		b.Fatal("Failed to generate random bls key: ", err)
	}
	sigPropose := proposeSk.Sign(propose[:])

	// generate signatures over votes
	var sigsPrevote []Signature
	var sigsPrecommit []Signature
	var pks [][]byte
	for i := 0; i < 100; i++ {
		sk, err := RandKey()
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey().Marshal()
		pks = append(pks, pk)
		sigPrevote := sk.Sign(prevote[:])
		sigsPrevote = append(sigsPrevote, sigPrevote)
		sigPrecommit := sk.Sign(precommit[:])
		sigsPrecommit = append(sigsPrecommit, sigPrecommit)
	}

	aggSigPrevote := AggregateSignatures(sigsPrevote)
	aggSigPrecommit := AggregateSignatures(sigsPrecommit)
	aggSig := AggregateSignatures([]Signature{sigPropose, aggSigPrevote, aggSigPrecommit})
	aggPk, err := AggregatePublicKeys(pks)
	if err != nil {
		b.Fatal("failed pks aggregation: ", err)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := aggSig.AggregateVerify([]PublicKey{proposeSk.PublicKey(), aggPk, aggPk}, [][32]byte{propose, prevote, precommit})
		if !res {
			b.Fatal("failed signature verification")
		}
	}

}

// benchmark aggregated signature verification on 15 messages
func BenchmarkAggregateVerify15(b *testing.B) {
	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	n := 15

	// create messages
	var msgs [][32]byte
	var msg [32]byte
	for i := 0; i < n; i++ {
		rand.Read(msg[:])
		msgs = append(msgs, msg)
	}

	// generate signatures over msgs
	sigs := make([][]Signature, n)
	var pks [][]byte
	for i := 0; i < 100; i++ {
		sk, err := RandKey()
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey().Marshal()
		pks = append(pks, pk)
		for j := 0; j < len(msgs); j++ {
			sig := sk.Sign(msgs[j][:])
			sigs[j] = append(sigs[j], sig)
		}
	}

	aggPk, err := AggregatePublicKeys(pks)
	if err != nil {
		b.Fatal("failed pks aggregation: ", err)
	}
	var aggSigs []Signature
	for i := 0; i < len(sigs); i++ {
		aggSig := AggregateSignatures(sigs[i])
		aggSigs = append(aggSigs, aggSig)
	}
	aggSigTotal := AggregateSignatures(aggSigs)

	var pksTotal []PublicKey
	for i := 0; i < n; i++ {
		pksTotal = append(pksTotal, aggPk)
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := aggSigTotal.AggregateVerify(pksTotal, msgs)
		if !res {
			b.Fatal("failed signature verification")
		}
	}

}

func BenchmarkAggregateVerifyInvalid(b *testing.B) {

	b.Run("1 in 500", func(b *testing.B) {
		b.Run("simple", func(b *testing.B) {
			runBenchmarkAggregateVerifyInvalid(b, 500, 21)
		})
		b.Run("fast", func(b *testing.B) {
			runBenchmarkFastAggregateVerifyInvalid(b, 500, 21)
		})
	})

	b.Run("1 in 1000", func(b *testing.B) {
		b.Run("simple", func(b *testing.B) {
			runBenchmarkAggregateVerifyInvalid(b, 1000, 21)
		})
		b.Run("fast", func(b *testing.B) {
			runBenchmarkFastAggregateVerifyInvalid(b, 1000, 21)

		})
	})

	b.Run("5 in 1000", func(b *testing.B) {
		b.Run("simple", func(b *testing.B) {
			runBenchmarkAggregateVerifyInvalid(b, 1000, 21, 37, 434, 795)
		})
		b.Run("fast", func(b *testing.B) {
			runBenchmarkFastAggregateVerifyInvalid(b, 1000, 21, 37, 434, 795)
		})
	})

	b.Run("10 in 1000", func(b *testing.B) {
		b.Run("simple", func(b *testing.B) {
			runBenchmarkAggregateVerifyInvalid(b, 1000, 21, 37, 155, 245, 343, 664, 712, 713, 800, 911)
		})
		b.Run("fast", func(b *testing.B) {
			runBenchmarkFastAggregateVerifyInvalid(b, 1000, 21, 37, 155, 245, 343, 664, 712, 713, 800, 911)
		})
	})

	b.Run("1/3 uniformly distributed of 1000", func(b *testing.B) {

		failed := make([]int, 333)

		for i := 0; i < 333; i++ {
			failed[i] = i * 3
		}

		b.Run("simple", func(b *testing.B) {
			runBenchmarkAggregateVerifyInvalid(b, 1000, failed...)
		})
		b.Run("fast", func(b *testing.B) {
			runBenchmarkFastAggregateVerifyInvalid(b, 1000, failed...)
		})
	})

	b.Run("1/3 concentrated of 1000", func(b *testing.B) {

		failed := make([]int, 333)

		for i := 0; i < 333; i++ {
			failed[i] = i
		}

		b.Run("simple", func(b *testing.B) {
			runBenchmarkAggregateVerifyInvalid(b, 1000, failed...)
		})
		b.Run("fast", func(b *testing.B) {
			runBenchmarkFastAggregateVerifyInvalid(b, 1000, failed...)
		})
	})

	b.Run("1 in 5000", func(b *testing.B) {
		b.Run("simple", func(b *testing.B) {
			runBenchmarkAggregateVerifyInvalid(b, 5000, 21)
		})
		b.Run("fast", func(b *testing.B) {
			runBenchmarkFastAggregateVerifyInvalid(b, 5000, 21)
		})
	})

}

func runBenchmarkAggregateVerifyInvalid(b *testing.B, size int, failedNs ...int) {

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	n := size
	//failedNs := []int{21, 37, 434, 795}

	invalidSk, err := RandKey()
	if err != nil {
		b.Fatal("Failed to generate random  key. error: ", err)
	}

	// create messages
	var msgs [][32]byte
	var msg [32]byte
	for i := 0; i < n; i++ {
		rand.Read(msg[:])
		msgs = append(msgs, msg)
	}

	// generate signatures over msgs
	sigs := make([]Signature, 0, n)

	var pks []PublicKey
	for i := 0; i < n; i++ {
		sk, err := RandKey()
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey()
		pks = append(pks, pk)

		sig := sk.Sign(msgs[i][:])
		sigs = append(sigs, sig)
	}

	for _, i := range failedNs {
		invalidSig := invalidSk.Sign(msgs[failedNs[0]][:])

		sigs[i] = invalidSig
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err := FindInvalidSignatures(sigs, pks, msgs)
		if err != nil {
			b.Fatalf("not expected err %s", err)
		}
	}

}

func runBenchmarkFastAggregateVerifyInvalid(b *testing.B, size int, failedNs ...int) {

	// initialize deterministic randomness
	rand := mrand.New(mrand.NewSource(0)) //nolint

	n := size

	invalidSk, err := RandKey()
	if err != nil {
		b.Fatal("Failed to generate random key. error: ", err)
	}

	var msg [32]byte
	rand.Read(msg[:])

	// generate signatures over msgs
	sigs := make([]Signature, 0, n)

	var pks []PublicKey
	for i := 0; i < n; i++ {
		sk, err := RandKey()
		if err != nil {
			b.Fatal("Failed to generate random bls key: ", err)
		}
		pk := sk.PublicKey()
		pks = append(pks, pk)

		sig := sk.Sign(msg[:])
		sigs = append(sigs, sig)
	}

	for _, i := range failedNs {
		invalidSig := invalidSk.Sign(msg[:])

		sigs[i] = invalidSig
	}

	// start the actual aggregation benchmarking
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err := FindFastInvalidSignatures(sigs, pks, msg)
		if err != nil {
			b.Fatalf("not expected err %s", err)
		}
	}

}
