package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
	"github.com/influxdata/influxdb/pkg/deep"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestMsgEncodingDecodingProposal(t *testing.T) {
	committeeSet, keys := NewTestCommitteeSetWithKeys(4)
	addr := committeeSet.Committee()[0].Address // round 3 - height 1 proposer
	height := uint64(1)
	round := int64(3)
	signer := makeSigner(keys[addr].consensus)
	signerMember := &committeeSet.Committee()[0]

	proposal := generateBlockProposal(round, new(big.Int).SetUint64(height), -1, false, signer, signerMember)
	m := ExtMsg{Msg: proposal, Verified: true}
	msgBytes, err := rlp.EncodeToBytes(&m)
	require.NoError(t, err)

	var decodedMsg ExtMsg
	err = rlp.DecodeBytes(msgBytes, &decodedMsg)
	require.NoError(t, err)
	deep.Equal(m, decodedMsg)
}

func TestMsgEncodingDecodingPreVote(t *testing.T) {
	committeeSet, keys := NewTestCommitteeSetWithKeys(4)
	addr := committeeSet.Committee()[0].Address // round 3 - height 1 proposer
	height := uint64(1)
	round := int64(3)
	signer := makeSigner(keys[addr].consensus)
	//signerMember := &committeeSet.Committee()[0]
	csize := len(committeeSet.Committee())

	vote := message.NewPrevote(round, height, common.Hash{}, signer, &committeeSet.Committee()[0], csize)

	m := ExtMsg{Msg: vote, Verified: true}
	msgBytes, err := rlp.EncodeToBytes(&m)
	require.NoError(t, err)

	var decodedMsg ExtMsg
	err = rlp.DecodeBytes(msgBytes, &decodedMsg)
	require.NoError(t, err)
	deep.Equal(m, decodedMsg)
}

func TestMsgEncodingDecodingPreCommit(t *testing.T) {
	committeeSet, keys := NewTestCommitteeSetWithKeys(4)
	addr := committeeSet.Committee()[0].Address // round 3 - height 1 proposer
	height := uint64(1)
	round := int64(3)
	signer := makeSigner(keys[addr].consensus)
	//signerMember := &committeeSet.Committee()[0]
	csize := len(committeeSet.Committee())

	vote := message.NewPrecommit(round, height, common.Hash{}, signer, &committeeSet.Committee()[0], csize)

	m := ExtMsg{Msg: vote, Verified: true}
	msgBytes, err := rlp.EncodeToBytes(&m)
	require.NoError(t, err)

	var decodedMsg ExtMsg
	err = rlp.DecodeBytes(msgBytes, &decodedMsg)
	require.NoError(t, err)
	deep.Equal(m, decodedMsg)
}
