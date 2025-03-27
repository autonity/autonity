package core

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/rlp"
)

func TestSyncMsg(t *testing.T) {
	round1 := &message.RoundMsgView{
		Round:             1,
		HaveProposal:      true,
		Prevotes:          []common.Hash{common.HexToHash("0xabcd"), common.HexToHash("0xefgh")},
		PrevotesSigners:   []*big.Int{big.NewInt(1), big.NewInt(2)},
		Precommits:        []common.Hash{common.HexToHash("0xijkl")},
		PrecommitsSigners: []*big.Int{big.NewInt(3)},
	}

	round2 := &message.RoundMsgView{
		Round:             2,
		HaveProposal:      false,
		Prevotes:          []common.Hash{},
		PrevotesSigners:   []*big.Int{},
		Precommits:        []common.Hash{},
		PrecommitsSigners: []*big.Int{},
	}

	syncMsg := &message.AskSyncMsg{
		Height:        100,
		KnownMessages: []*message.RoundMsgView{round1, round2},
	}

	encoded, err := rlp.EncodeToBytes(syncMsg)
	require.NoError(t, err)

	var decodedSyncMsg message.AskSyncMsg
	err = rlp.DecodeBytes(encoded, &decodedSyncMsg)
	require.NoError(t, err)

	require.Equal(t, syncMsg.Height, decodedSyncMsg.Height)
	require.Equal(t, len(syncMsg.KnownMessages), len(decodedSyncMsg.KnownMessages))

	for i := range syncMsg.KnownMessages {
		round1 := syncMsg.KnownMessages[i]
		round2 := decodedSyncMsg.KnownMessages[i]

		require.Equal(t, round1.Round, round2.Round)
		require.Equal(t, round1.HaveProposal, round2.HaveProposal)

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
