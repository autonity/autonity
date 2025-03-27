package core

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestSyncMsg(t *testing.T) {
	round1 := &message.RoundMsgView{
		Round:             1,
		Proposal:          common.HexToHash("0x1234"),
		Prevotes:          []common.Hash{common.HexToHash("0xabcd"), common.HexToHash("0xefgh")},
		PrevotesSigners:   []*big.Int{big.NewInt(1), big.NewInt(2)},
		Precommits:        []common.Hash{common.HexToHash("0xijkl")},
		PrecommitsSigners: []*big.Int{big.NewInt(3)},
	}

	round2 := &message.RoundMsgView{
		Round:             2,
		Proposal:          common.Hash{},
		Prevotes:          []common.Hash{},
		PrevotesSigners:   []*big.Int{},
		Precommits:        []common.Hash{},
		PrecommitsSigners: []*big.Int{},
	}

	syncMsg := &message.SyncMsg{
		Height:      100,
		RoundsViews: []*message.RoundMsgView{round1, round2},
	}

	encoded, err := rlp.EncodeToBytes(syncMsg)
	require.NoError(t, err)

	var decodedSyncMsg message.SyncMsg
	err = rlp.DecodeBytes(encoded, &decodedSyncMsg)
	require.NoError(t, err)

	require.Equal(t, syncMsg.Height, decodedSyncMsg.Height)
	require.Equal(t, len(syncMsg.RoundsViews), len(decodedSyncMsg.RoundsViews))

	for i := range syncMsg.RoundsViews {
		round1 := syncMsg.RoundsViews[i]
		round2 := decodedSyncMsg.RoundsViews[i]

		require.Equal(t, round1.Round, round2.Round)
		require.Equal(t, round1.Proposal, round2.Proposal)

		require.Equal(t, len(round1.Prevotes), len(round2.Prevotes))
		for j := range round1.Prevotes {
			require.Equal(t, round1.Prevotes[j], round2.Prevotes[j])
		}

		require.Equal(t, len(round1.PrevotesSigners), len(round2.PrevotesSigners))
		for j := range round1.PrevotesSigners {
			require.Equal(t, round1.PrevotesSigners[j].String(), round2.PrevotesSigners[j].String())
		}

		require.Equal(t, len(round1.Precommits), len(round2.Precommits))
		for j := range round1.Precommits {
			require.Equal(t, round1.Precommits[j], round2.Precommits[j])
		}

		require.Equal(t, len(round1.PrecommitsSigners), len(round2.PrecommitsSigners))
		for j := range round1.PrecommitsSigners {
			require.Equal(t, round1.PrecommitsSigners[j].String(), round2.PrecommitsSigners[j].String())
		}
	}
}
