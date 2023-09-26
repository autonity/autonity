package backend

import (
	"github.com/autonity/autonity/params/generated"
	"go.uber.org/mock/gomock"
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/rpc"

	"github.com/stretchr/testify/assert"
)

func TestGetCommittee(t *testing.T) {
	want := types.Committee{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := consensus.NewMockChainReader(ctrl)
	h := &types.Header{Number: big.NewInt(1)}
	c.EXPECT().GetHeaderByNumber(uint64(1)).Return(h)
	api := &API{
		chain: c,
		getCommittee: func(header *types.Header, chain consensus.ChainReader) (types.Committee, error) {
			if header == h && chain == c {
				return want, nil
			}
			return nil, nil
		},
	}

	bn := rpc.BlockNumber(1)

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

		c := consensus.NewMockChainReader(ctrl)
		h := &types.Header{Number: big.NewInt(1)}
		c.EXPECT().GetHeaderByHash(hash).Return(h)

		want := types.Committee{}

		api := &API{
			chain: c,
			getCommittee: func(header *types.Header, chain consensus.ChainReader) (types.Committee, error) {
				if header == h && chain == c {
					return want, nil
				}
				return nil, nil
			},
		}

		got, err := api.GetCommitteeAtHash(hash)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestAPIGetContractABI(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	assert.Nil(t, err)
	_, err = chain.InsertChain(types.Blocks{block})
	assert.Nil(t, err)

	want := generated.AutonityAbi

	api := &API{
		tendermint: engine,
	}

	got := *api.GetContractABI()
	assert.Equal(t, want, got)
}

func TestAPIGetContractAddress(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	assert.Nil(t, err)
	_, err = chain.InsertChain(types.Blocks{block})
	assert.Nil(t, err)

	want := autonity.AutonityContractAddress

	API := &API{
		tendermint: engine,
	}

	got := API.GetContractAddress()
	assert.Equal(t, want, got)
}
