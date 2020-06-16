package backend

import (
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/acdefault"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetCommittee(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	assert.Nil(t, err)
	_, err = chain.InsertChain(types.Blocks{block})
	assert.Nil(t, err)

	want, err := engine.Committee(1)
	assert.Nil(t, err)

	backendAPI := &API{tendermint: engine}

	bn := rpc.BlockNumber(0)
	got, err := backendAPI.GetCommittee(&bn)
	assert.Nil(t, err)

	assert.Equal(t, want.Committee(), got)
}

func TestGetCommitteeAtHash(t *testing.T) {
	t.Run("unknown block given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hash := common.HexToHash("0x0123456789")

		chain := consensus.NewMockChainReader(ctrl)
		chain.EXPECT().GetHeaderByHash(hash).Return(nil)

		API := &API{
			chain: chain,
		}

		_, err := API.GetCommitteeAtHash(hash)
		if err != errUnknownBlock {
			t.Fatalf("expected %v, got %v", errUnknownBlock, err)
		}
	})

	t.Run("valid block given, committee returned", func(t *testing.T) {
		chain, engine := newBlockChain(1)
		block, err := makeBlock(chain, engine, chain.Genesis())
		assert.Nil(t, err)
		_, err = chain.InsertChain(types.Blocks{block})
		assert.Nil(t, err)

		want, err := engine.Committee(1)
		assert.Nil(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		h := block.Hash()
		mockChain := consensus.NewMockChainReader(ctrl)
		mockChain.EXPECT().GetHeaderByHash(h).Return(block.Header())

		API := &API{
			chain:      mockChain,
			tendermint: engine,
		}

		got, err := API.GetCommitteeAtHash(h)
		assert.Nil(t, err)

		assert.Equal(t, want.Committee(), got)
	})
}

func TestAPIGetContractABI(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	assert.Nil(t, err)
	_, err = chain.InsertChain(types.Blocks{block})
	assert.Nil(t, err)

	want := acdefault.ABI()

	API := &API{
		tendermint: engine,
	}

	got := API.GetContractABI()
	assert.Equal(t, want, got)
}

func TestAPIGetContractAddress(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	assert.Nil(t, err)
	_, err = chain.InsertChain(types.Blocks{block})
	assert.Nil(t, err)

	want := autonity.ContractAddress

	API := &API{
		tendermint: engine,
	}

	got := API.GetContractAddress()
	assert.Equal(t, want, got)
}

func TestAPIGetWhitelist(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	assert.Nil(t, err)
	_, err = chain.InsertChain(types.Blocks{block})
	assert.Nil(t, err)

	want := []string{"enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303"}

	API := &API{
		tendermint: engine,
	}

	got := API.GetWhitelist()

	assert.Equal(t, want, got)
}
