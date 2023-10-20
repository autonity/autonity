package core

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/stretchr/testify/require"
	"math/big"
	"math/rand"
	"testing"
)

func TestCore_MeasureHeightRoundMetrics(t *testing.T) {
	t.Run("measure metrics of new height", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.MeasureHeightRoundMetrics(0)
		if m := metrics.Get("tendermint/height/change"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics of new round", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.MeasureHeightRoundMetrics(1)
		if m := metrics.Get("tendermint/round/change"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}

func TestCore_measureMetricsOnTimeOut(t *testing.T) {
	t.Run("measure metrics on Timeout of propose", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(message.MsgProposal, 2)
		if m := metrics.Get("tendermint/timer/propose"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on Timeout of prevote", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(message.MsgPrevote, 2)
		if m := metrics.Get("tendermint/timer/prevote"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on Timeout of precommit", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   types.NewTimeout(types.Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   types.NewTimeout(types.Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: types.NewTimeout(types.Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(message.MsgPrecommit, 2)
		if m := metrics.Get("tendermint/timer/precommit"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}

func TestCore_Setters(t *testing.T) {
	t.Run("SetStep", func(t *testing.T) {
		c := &Core{
			logger: log.New("addr", "tendermint core"),
		}

		c.SetStep(types.Propose)
		require.Equal(t, types.Propose, c.step)
	})

	t.Run("setRound", func(t *testing.T) {
		c := &Core{}
		c.setRound(0)
		require.Equal(t, int64(0), c.Round())
	})

	t.Run("setHeight", func(t *testing.T) {
		c := &Core{}
		c.setHeight(new(big.Int).SetUint64(10))
		require.Equal(t, uint64(10), c.height.Uint64())
	})

	t.Run("setCommitteeSet", func(t *testing.T) {
		c := &Core{}
		committeeSizeAndMaxRound := maxSize
		committeeSet, _ := prepareCommittee(t, committeeSizeAndMaxRound)
		c.setCommitteeSet(committeeSet)
		require.Equal(t, committeeSet, c.CommitteeSet())
	})

	t.Run("setLastHeader", func(t *testing.T) {
		c := &Core{}
		prevHeight := big.NewInt(int64(rand.Intn(100) + 1)) //nolint
		prevBlock := generateBlock(prevHeight)
		c.setLastHeader(prevBlock.Header())
		require.Equal(t, prevBlock.Header(), c.LastHeader())
	})
}

func TestCore_AcceptVote(t *testing.T) {

	t.Run("AcceptPreVote", func(t *testing.T) {
		messagesMap := message.NewMessagesMap()
		roundMessage := messagesMap.GetOrCreate(0)
		c := &Core{
			messages:         messagesMap,
			curRoundMessages: roundMessage,
		}
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1)) //nolint
		currentRound := int64(0)
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		clientAddr := crypto.PubkeyToAddress(key.PublicKey)

		prevoteMsg, _, _ := prepareVote(t, message.MsgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, key)
		c.AcceptVote(c.CurRoundMessages(), types.Prevote, common.Hash{}, *prevoteMsg)
		require.Equal(t, 1, len(c.CurRoundMessages().GetMessages()))
	})

	t.Run("AcceptPreCommit", func(t *testing.T) {
		messagesMap := message.NewMessagesMap()
		roundMessage := messagesMap.GetOrCreate(0)
		c := &Core{
			messages:         messagesMap,
			curRoundMessages: roundMessage,
		}
		currentHeight := big.NewInt(int64(rand.Intn(maxSize) + 1)) //nolint
		currentRound := int64(0)
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		clientAddr := crypto.PubkeyToAddress(key.PublicKey)

		prevoteMsg, _, _ := prepareVote(t, message.MsgPrecommit, currentRound, currentHeight, common.Hash{}, clientAddr, key)
		c.AcceptVote(c.CurRoundMessages(), types.Precommit, common.Hash{}, *prevoteMsg)
		require.Equal(t, 1, len(c.CurRoundMessages().GetMessages()))
	})
}

func signer(prv *ecdsa.PrivateKey) func(data []byte) ([]byte, error) {
	return func(data []byte) ([]byte, error) {
		hashData := crypto.Keccak256(data)
		return crypto.Sign(hashData, prv)
	}
}
