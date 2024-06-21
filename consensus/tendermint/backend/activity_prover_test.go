package backend

import (
	"github.com/autonity/autonity/common"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
)

func TestAssembleProof(t *testing.T) {
	t.Run("for the first delta blocks of the epoch, assembling should return empty proof", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		delta := chain.ProtocolContracts().OmissionDelta().Uint64()

		proof, round, err := backend.assembleActivityProof(0)
		require.Equal(t, types.AggregateSignature{}, proof)
		require.Nil(t, proof.Signature)
		require.Nil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)

		for i := 0; i < int(delta)-1; i++ {
			mineOneBlock(t, chain, backend)
		}

		proof, round, err = backend.assembleActivityProof(delta)
		require.Equal(t, types.AggregateSignature{}, proof)
		require.Nil(t, proof.Signature)
		require.Nil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)
	})
	t.Run("from block Delta+1 of the epoch, assembling should return a valid proof", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		self := &chain.Genesis().Header().Epoch.Committee.Members[0]
		delta := chain.ProtocolContracts().OmissionDelta().Uint64()

		for i := 0; i < int(delta); i++ {
			mineOneBlock(t, chain, backend)
			header := chain.CurrentHeader()

			// insert precommit in msgStore to be able to assemble proof later
			precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), backend.Sign, self, 1)
			backend.MsgStore.Save(precommit)
		}

		proof, round, err := backend.assembleActivityProof(delta + 1)
		require.NotEqual(t, types.AggregateSignature{}, proof)
		require.NotNil(t, proof.Signature)
		require.NotNil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)

		for i := 0; i < int(delta); i++ {
			mineOneBlock(t, chain, backend)
			header := chain.CurrentHeader()

			// insert precommit in msgStore to be able to assemble proof later
			precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), backend.Sign, self, 1)
			backend.MsgStore.Save(precommit)
		}

		proof, round, err = backend.assembleActivityProof(delta*2 + 1)
		require.NotEqual(t, types.AggregateSignature{}, proof)
		require.NotNil(t, proof.Signature)
		require.NotNil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)

		// NOTE: we are not updating the msg store with the new precommits here.
		// this is to simulate a node crash. Assemble proof should still return no error, but it should return an empty proof
		for i := 0; i < int(delta); i++ {
			mineOneBlock(t, chain, backend)
		}

		proof, round, err = backend.assembleActivityProof(delta*3 + 1)
		require.Equal(t, types.AggregateSignature{}, proof)
		require.Nil(t, proof.Signature)
		require.Nil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)

	})
	t.Run("proof should be empty if we do not have quorum precommits to provide", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		self := &chain.Genesis().Header().Epoch.Committee.Members[0]
		delta := chain.ProtocolContracts().OmissionDelta().Uint64()

		for i := 0; i < int(delta); i++ {
			mineOneBlock(t, chain, backend)
			header := chain.CurrentHeader()

			// add fake precommit with low voting power, to simulate not enough voting power to build activity proof
			precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), backend.Sign, self, 1)
			signers := precommit.Signers()
			powers := make(map[int]*big.Int)
			powers[0] = common.Big0
			signers.AssignPower(powers, common.Big1)
			precommitWithLowPower := message.NewFakePrecommit(message.Fake{
				FakeCode:           message.PrecommitCode,
				FakeRound:          uint64(precommit.R()),
				FakeHeight:         precommit.H(),
				FakeValue:          precommit.Value(),
				FakePayload:        precommit.Payload(),
				FakeHash:           precommit.Hash(),
				FakeSigners:        precommit.Signers(),
				FakeSignature:      precommit.Signature(),
				FakeSignatureInput: precommit.SignatureInput(),
				FakeSignerKey:      precommit.SignerKey(),
			})
			backend.MsgStore.Save(precommitWithLowPower)
		}

		proof, round, err := backend.assembleActivityProof(delta + 1)
		require.Equal(t, types.AggregateSignature{}, proof)
		require.Nil(t, proof.Signature)
		require.Nil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)
	})
}

/* //TODO(lorenzo) fix these tests
func TestValidateProof(t *testing.T) {
	t.Run("for the first delta blocks of the epoch, the proof should be rejected if not empty", func(t *testing.T) {
		_, backend := newBlockChain(1)

		// proof empty, proposal is accepted
		isProposerFaulty, effort, absents, err := backend.validateActivityProof(types.AggregateSignature{}, 1, 0)
		require.False(t, isProposerFaulty)
		require.Equal(t, common.Big0, effort)
		require.Equal(t, []common.Address{}, absents)
		require.Nil(t, err)

		// proof not empty, proposal is rejected
		proof := types.AggregateSignature{Signature: new(blst.BlsSignature), Signers: types.NewSigners(1)}
		isProposerFaulty, effort, absents, err = backend.validateActivityProof(proof, 1, 0)
		require.False(t, isProposerFaulty)
		require.Equal(t, common.Big0, effort)
		require.Equal(t, []common.Address{}, absents)
		require.True(t, errors.Is(err, ErrNotEmptyActivityProof))
	})
	t.Run("from block Delta+1 of the epoch, proof should be a valid aggregate signature", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		committee := chain.Genesis().Header().Committee
		self := &chain.Genesis().Header().Committee[0]
		delta := chain.ProtocolContracts().OmissionDelta().Uint64()

		for i := 0; i < int(delta); i++ {
			mineOneBlock(t, chain, backend)
			header := chain.CurrentHeader()

			// insert precommit in msgStore to be able to assemble proof later
			precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), backend.Sign, self, 1)
			backend.MsgStore.Save(precommit)
		}

		isProposerFaulty, effort, absents, err := backend.validateActivityProof(types.AggregateSignature{}, delta+1, 0)
		require.True(t, isProposerFaulty)
		require.Equal(t, common.Big0, effort)
		require.Equal(t, 0, len(absents))
		require.Nil(t, err)

		// proof is invalid, the signers object is too big wrt the committee
		randomSig := backend.Sign(common.Hash{})
		proof := types.AggregateSignature{Signature: randomSig.(*blst.BlsSignature), Signers: types.NewSigners(10)}
		isProposerFaulty, effort, absents, err = backend.validateActivityProof(proof, delta+1, 0)
		require.False(t, isProposerFaulty)
		require.Equal(t, common.Big0, effort)
		require.Equal(t, 0, len(absents))
		require.NotNil(t, err)

		proof, round, err := backend.assembleActivityProof(delta + 1)
		require.NotEqual(t, types.AggregateSignature{}, proof)
		require.Equal(t, uint64(0), round)
		require.Nil(t, err)

		// valid proof!!
		expectedEffort := new(big.Int).Sub(self.VotingPower, bft.Quorum(committee.TotalVotingPower()))
		isProposerFaulty, effort, absents, err = backend.validateActivityProof(proof, delta+1, 0)
		require.False(t, isProposerFaulty)
		require.Equal(t, expectedEffort, effort)
		require.Equal(t, 0, len(absents))
		require.Nil(t, err)

	})
}

func TestVerifyActivityProof(t *testing.T) {
	chain, backend := newBlockChain(1)
	committee := chain.Genesis().Header().Committee
	self := &chain.Genesis().Header().Committee[0]
	delta := chain.ProtocolContracts().OmissionDelta().Uint64()

	// pass the first DeltaBlocks
	for i := 0; i < int(delta); i++ {
		mineOneBlock(t, chain, backend)
		header := chain.CurrentHeader()

		// insert precommit in msgStore to be able to assemble proof later
		precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), backend.Sign, self, 1)
		backend.MsgStore.Save(precommit)
	}

	// invalid signature, verification should fail
	randomSig := backend.Sign(common.Hash{})
	signersBitmap := types.NewSigners(1)
	signersBitmap.Increment(self)
	proof := types.AggregateSignature{Signature: randomSig.(*blst.BlsSignature), Signers: signersBitmap}
	signers, effort, err := backend.verifyActivityProof(proof, committee, 1, 0)
	require.Equal(t, make(map[common.Address]struct{}), signers)
	require.Equal(t, common.Big0, effort)
	require.Equal(t, ErrInvalidActivityProofSignature, err)

	// valid proof, verification should go smooth
	proof, round, err := backend.assembleActivityProof(delta + 1)
	require.NotEqual(t, types.AggregateSignature{}, proof)
	require.Equal(t, uint64(0), round)
	require.Nil(t, err)
	expectedEffort := new(big.Int).Sub(self.VotingPower, bft.Quorum(committee.TotalVotingPower()))
	signers, effort, err = backend.verifyActivityProof(proof, committee, 1, 0)
	require.Equal(t, 1, len(signers))
	_, ok := signers[self.Address]
	require.True(t, ok)
	require.Equal(t, expectedEffort, effort)
	require.NoError(t, err)

	// simulate insufficient voting power by artificially inflating the committee
	for i := 0; i < 5; i++ {
		key, err := blst.RandKey()
		require.NoError(t, err)
		committee = append(committee, types.CommitteeMember{
			VotingPower:  new(big.Int).Set(self.VotingPower),
			Address:      common.Address{byte(i)},
			ConsensusKey: key.PublicKey(),
		})
	}
	signers, effort, err = backend.verifyActivityProof(proof, committee, 1, 0)
	require.Equal(t, make(map[common.Address]struct{}), signers)
	require.Equal(t, common.Big0, effort)
	require.Equal(t, ErrInsufficientActivityProof, err)

	// check that the signers returned are the correct one
	// inflate the voting power of the signer so that we achieve quorum
	committee[0].VotingPower = new(big.Int).Mul(committee.TotalVotingPower(), common.Big5)
	expectedEffort = new(big.Int).Sub(committee[0].VotingPower, bft.Quorum(committee.TotalVotingPower()))
	signers, effort, err = backend.verifyActivityProof(proof, committee, 1, 0)
	require.Equal(t, 1, len(signers))
	_, ok = signers[committee[0].Address]
	require.True(t, ok)
	require.Equal(t, expectedEffort, effort)
	require.NoError(t, err)

}

func TestComputeAbsents(t *testing.T) {
	n := 10
	var committee []types.CommitteeMember

	for i := 0; i < n; i++ {
		committee = append(committee, types.CommitteeMember{Address: common.Address{byte(i)}})
	}

	signers := make(map[common.Address]struct{})
	signers[committee[0].Address] = struct{}{}
	signers[committee[1].Address] = struct{}{}
	signers[committee[7].Address] = struct{}{}
	absents := computeAbsents(signers, committee)
	require.Equal(t, n-len(signers), len(absents))

	found := make(map[common.Address]bool)
	for _, absent := range absents {
		found[absent] = true
		_, foundInSigners := signers[absent]
		if foundInSigners {
			t.Fatal("signer has been included into absents")
		}
	}

	for signer := range signers {
		found[signer] = true
	}

	// all members should be either signers or absents
	require.Equal(t, len(committee), len(found))

}
*/
