package backend

import (
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	"github.com/autonity/autonity/rpc"
	"github.com/stretchr/testify/assert"
)

func TestGetCommittee(t *testing.T) {
	chain, engine := newBlockChain(1)
	want := chain.Genesis().Header().Epoch.Committee
	bn := rpc.BlockNumber(0)
	api := &API{
		chain:      chain,
		tendermint: engine,
	}
	got, err := api.GetCommittee(&bn)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGetCommitteeAtHash(t *testing.T) {
	t.Run("unknown block given, error returned", func(t *testing.T) {
		hash := common.HexToHash("0x0123456789")
		chain, engine := newBlockChain(1)
		api := &API{
			chain:      chain,
			tendermint: engine,
		}
		_, err := api.GetCommitteeAtHash(hash)
		if err != errUnknownBlock {
			t.Fatalf("expected %v, got %v", errUnknownBlock, err)
		}
	})

	t.Run("valid block given, committee returned", func(t *testing.T) {
		chain, engine := newBlockChain(1)
		api := &API{
			chain:      chain,
			tendermint: engine,
		}

		hash := chain.Genesis().Hash()
		want := chain.Genesis().Header().Epoch.Committee
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
