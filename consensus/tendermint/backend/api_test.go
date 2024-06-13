package backend

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	"github.com/autonity/autonity/rpc"
)

func TestGetCommittee(t *testing.T) {
	want := &types.Committee{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bn := rpc.BlockNumber(1)
	c := consensus.NewMockChainReader(ctrl)
	c.EXPECT().CommitteeOfHeight(uint64(bn.Int64())).Return(want, nil)
	api := &API{
		chain: c,
	}
	got, err := api.GetCommittee(&bn)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetCommitteeAtHash(t *testing.T) {
	t.Run("unknown block given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hash := common.HexToHash("0x0123456789")

		chain := consensus.NewMockChainReader(ctrl)
		chain.EXPECT().GetHeaderByHash(hash).Return(nil)

		api := &API{
			chain: chain,
		}

		_, err := api.GetCommitteeAtHash(hash)
		if err != errUnknownBlock {
			t.Fatalf("expected %v, got %v", errUnknownBlock, err)
		}
	})

	t.Run("valid block given, committee returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hash := common.HexToHash("0x0123456789")
		want := &types.Committee{}

		c := consensus.NewMockChainReader(ctrl)
		h := &types.Header{Number: big.NewInt(1)}
		c.EXPECT().GetHeaderByHash(hash).Return(h)
		c.EXPECT().CommitteeOfHeight(h.Number.Uint64()).Return(want, nil)
		api := &API{
			chain: c,
		}

		got, err := api.GetCommitteeAtHash(hash)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestAPIGetContractABI(t *testing.T) {
	_, engine := newBlockChain(1)
	api := &API{
		tendermint: engine,
	}
	want := generated.AutonityAbi
	got := *api.GetContractABI()
	assert.Equal(t, want, got)
}

func TestAPIGetContractAddress(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	assert.Nil(t, err)
	_, err = chain.InsertChain(types.Blocks{block})
	assert.Nil(t, err)

	want := params.AutonityContractAddress

	API := &API{
		tendermint: engine,
	}

	got := API.GetContractAddress()
	assert.Equal(t, want, got)
}
