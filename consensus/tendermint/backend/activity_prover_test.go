package backend

import (
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

}
