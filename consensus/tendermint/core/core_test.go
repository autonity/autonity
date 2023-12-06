package core

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCore_MeasureHeightRoundMetrics(t *testing.T) {
	t.Run("measure metrics of new height", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureHeightRoundMetrics(0)
		if m := metrics.Get("tendermint/height/change"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics of new round", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureHeightRoundMetrics(1)
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
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(Propose, 2)
		if m := metrics.Get("tendermint/timer/propose"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on Timeout of prevote", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(Prevote, 2)
		if m := metrics.Get("tendermint/timer/prevote"); m == nil {
			t.Fatalf("test case failed.")
		}
	})

	t.Run("measure metrics on Timeout of precommit", func(t *testing.T) {
		c := &Core{
			address:          common.Address{},
			logger:           log.New("Core", "test", "id", 0),
			proposeTimeout:   NewTimeout(Propose, log.New("Core", "test", "id", 0)),
			prevoteTimeout:   NewTimeout(Prevote, log.New("Core", "test", "id", 0)),
			precommitTimeout: NewTimeout(Precommit, log.New("Core", "test", "id", 0)),
			round:            0,
			height:           big.NewInt(1),
		}
		c.measureMetricsOnTimeOut(Precommit, 2)
		if m := metrics.Get("tendermint/timer/precommit"); m == nil {
			t.Fatalf("test case failed.")
		}
	})
}

func TestCore_Setters(t *testing.T) {
	t.Run("SetStep", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		backendMock := interfaces.NewMockBackend(ctrl)
		// it might get called or not based on how the go routine is scheduled
		backendMock.EXPECT().ProcessFutureMsgs(uint64(2)).MaxTimes(1)

		c := &Core{
			logger:  log.New("addr", "tendermint core"),
			height:  big.NewInt(2),
			backend: backendMock,
		}

		c.SetStep(Propose)
		require.Equal(t, Propose, c.step)
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

/*TODO(lorenzo) remove?
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

		prevoteMsg, _ := prepareVote(t, consensus.MsgPrevote, currentRound, currentHeight, common.Hash{}, clientAddr, key)
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

		prevoteMsg, _ := prepareVote(t, consensus.MsgPrecommit, currentRound, currentHeight, common.Hash{}, clientAddr, key)
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
*/
