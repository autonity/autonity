package backend

import (
	"errors"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/crypto/blst"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/consensus/tendermint"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
)

func TestAssembleProof(t *testing.T) {
	t.Run("for the first delta blocks of the epoch, assembling should return empty proof", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		proof, round, err := backend.assembleActivityProof(0)
		require.Equal(t, types.AggregateSignature{}, proof)
		require.Nil(t, proof.Signature)
		require.Nil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)

		for i := 0; i < tendermint.DeltaBlocks-1; i++ {
			mineOneBlock(t, chain, backend)
		}

		proof, round, err = backend.assembleActivityProof(tendermint.DeltaBlocks)
		require.Equal(t, types.AggregateSignature{}, proof)
		require.Nil(t, proof.Signature)
		require.Nil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)
	})
	t.Run("from block Delta+1 of the epoch, assembling should return a valid proof", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		self := &chain.Genesis().Header().Committee[0]

		for i := 0; i < tendermint.DeltaBlocks; i++ {
			mineOneBlock(t, chain, backend)
			header := chain.CurrentHeader()

			// insert precommit in msgStore to be able to assemble proof later
			precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), backend.Sign, self, 1)
			backend.MsgStore.Save(precommit)
		}

		proof, round, err := backend.assembleActivityProof(tendermint.DeltaBlocks + 1)
		require.NotEqual(t, types.AggregateSignature{}, proof)
		require.NotNil(t, proof.Signature)
		require.NotNil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)

		for i := 0; i < tendermint.DeltaBlocks; i++ {
			mineOneBlock(t, chain, backend)
			header := chain.CurrentHeader()

			// insert precommit in msgStore to be able to assemble proof later
			precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), backend.Sign, self, 1)
			backend.MsgStore.Save(precommit)
		}

		proof, round, err = backend.assembleActivityProof(tendermint.DeltaBlocks*2 + 1)
		require.NotEqual(t, types.AggregateSignature{}, proof)
		require.NotNil(t, proof.Signature)
		require.NotNil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)

		// NOTE: we are not updating the msg store with the new precommits here.
		// this is to simulate a node crash. Assemble proof should still return no error, but it should return an empty proof
		for i := 0; i < tendermint.DeltaBlocks; i++ {
			mineOneBlock(t, chain, backend)
		}

		proof, round, err = backend.assembleActivityProof(tendermint.DeltaBlocks*3 + 1)
		require.Equal(t, types.AggregateSignature{}, proof)
		require.Nil(t, proof.Signature)
		require.Nil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)

	})
	t.Run("proof should be empty if we do not have quorum precommits to provide", func(t *testing.T) {
		chain, backend := newBlockChain(1)

		self := &chain.Genesis().Header().Committee[0]

		for i := 0; i < tendermint.DeltaBlocks; i++ {
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

		proof, round, err := backend.assembleActivityProof(tendermint.DeltaBlocks + 1)
		require.Equal(t, types.AggregateSignature{}, proof)
		require.Nil(t, proof.Signature)
		require.Nil(t, proof.Signers)
		require.Equal(t, uint64(0), round)
		require.Equal(t, err, nil)
	})
}

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

		for i := 0; i < tendermint.DeltaBlocks; i++ {
			mineOneBlock(t, chain, backend)
			header := chain.CurrentHeader()

			// insert precommit in msgStore to be able to assemble proof later
			precommit := message.NewPrecommit(int64(header.Round), header.Number.Uint64(), header.Hash(), backend.Sign, self, 1)
			backend.MsgStore.Save(precommit)
		}

		isProposerFaulty, effort, absents, err := backend.validateActivityProof(types.AggregateSignature{}, tendermint.DeltaBlocks+1, 0)
		require.True(t, isProposerFaulty)
		require.Equal(t, common.Big0, effort)
		require.Equal(t, 0, len(absents))
		require.Nil(t, err)

		// reject the proposal if the proof is only partially filled in
		randomSig := backend.Sign(common.Hash{})
		proof := types.AggregateSignature{Signature: randomSig.(*blst.BlsSignature), Signers: nil}
		isProposerFaulty, effort, absents, err = backend.validateActivityProof(proof, tendermint.DeltaBlocks+1, 0)
		require.False(t, isProposerFaulty)
		require.Equal(t, common.Big0, effort)
		require.Equal(t, 0, len(absents))
		require.Equal(t, ErrIncompleteActivityProof, err)

		proof.Signers = types.NewSigners(10)
		isProposerFaulty, effort, absents, err = backend.validateActivityProof(proof, tendermint.DeltaBlocks+1, 0)
		require.False(t, isProposerFaulty)
		require.Equal(t, common.Big0, effort)
		require.Equal(t, 0, len(absents))
		require.NotNil(t, err)

		proof, round, err := backend.assembleActivityProof(tendermint.DeltaBlocks + 1)
		require.NotEqual(t, types.AggregateSignature{}, proof)
		require.Equal(t, uint64(0), round)
		require.Nil(t, err)

		// valid proof!!
		expectedEffort := new(big.Int).Sub(self.VotingPower, bft.Quorum(committee.TotalVotingPower()))
		isProposerFaulty, effort, absents, err = backend.validateActivityProof(proof, tendermint.DeltaBlocks+1, 0)
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

	// pass the first DeltaBlocks
	for i := 0; i < tendermint.DeltaBlocks; i++ {
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
	require.Equal(t, []common.Address{}, signers)
	require.Equal(t, common.Big0, effort)
	require.Equal(t, ErrInvalidActivityProofSignature, err)

	// valid proof, verification should go smooth
	proof, round, err := backend.assembleActivityProof(tendermint.DeltaBlocks + 1)
	require.NotEqual(t, types.AggregateSignature{}, proof)
	require.Equal(t, uint64(0), round)
	require.Nil(t, err)
	expectedEffort := new(big.Int).Sub(self.VotingPower, bft.Quorum(committee.TotalVotingPower()))
	signers, effort, err = backend.verifyActivityProof(proof, committee, 1, 0)
	require.Equal(t, []common.Address{self.Address}, signers)
	require.Equal(t, expectedEffort, effort)
	require.NoError(t, err)

	// simulate insufficient voting power by artificially inflating the committee
	var keys []blst.SecretKey
	for i := 0; i < 5; i++ {
		key, err := blst.RandKey()
		require.NoError(t, err)
		committee = append(committee, types.CommitteeMember{
			VotingPower:  new(big.Int).Set(self.VotingPower),
			Address:      common.Address{byte(i)},
			ConsensusKey: key.PublicKey(),
		})
		keys = append(keys, key)
	}
	signers, effort, err = backend.verifyActivityProof(proof, committee, 1, 0)
	require.Equal(t, []common.Address{}, signers)
	require.Equal(t, common.Big0, effort)
	require.Equal(t, ErrInsufficientActivityProof, err)

	// check that the signers returned are the correct one
	// inflate the voting power of the signer so that we achieve quorum
	committee[0].VotingPower = new(big.Int).Mul(committee.TotalVotingPower(), common.Big5)
	expectedEffort = new(big.Int).Sub(committee[0].VotingPower, bft.Quorum(committee.TotalVotingPower()))
	signers, effort, err = backend.verifyActivityProof(proof, committee, 1, 0)
	require.Equal(t, []common.Address{committee[0].Address}, signers)
	require.Equal(t, expectedEffort, effort)
	require.NoError(t, err)

}

func TestComputeAbsents(t *testing.T) {
	var keys []blst.SecretKey
	n := 10
	var committee []types.CommitteeMember

	for i := 0; i < n; i++ {
		key, err := blst.RandKey()
		require.NoError(t, err)
		committee = append(committee, types.CommitteeMember{Address: common.Address{byte(i)}})
		keys = append(keys, key)
	}

	var signers []common.Address
	signers = append(signers, committee[0].Address)
	signers = append(signers, committee[1].Address)
	signers = append(signers, committee[7].Address)
	absents := computeAbsents(signers, committee)
	require.Equal(t, n-len(signers), len(absents))

	found := make(map[common.Address]bool)
	for _, absent := range absents {
		found[absent] = true
		for _, signer := range signers {
			found[signer] = true // repetitive but whatever
			if absent == signer {
				t.Fatal("signer has been included into absents")
			}
		}
	}

	// all members should be either signers or absents
	require.Equal(t, len(committee), len(found))

}
