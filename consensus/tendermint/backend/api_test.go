package backend

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

func TestGetCommittee(t *testing.T) {
	want := types.Committee{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := consensus.NewMockChainReader(ctrl)
	h := &types.Header{Number: big.NewInt(1)}
	c.EXPECT().GetHeaderByNumber(uint64(1)).Return(h)
	API := &API{
		chain: c,
		getCommittee: func(header *types.Header, chain consensus.ChainReader) (types.Committee, error) {
			if header == h && chain == c {
				return want, nil
			}
			return nil, nil
		},
	}

	bn := rpc.BlockNumber(1)

	got, err := API.GetCommittee(&bn)
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

		API := &API{
			chain: chain,
		}

		_, err := API.GetCommitteeAtHash(hash)
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

		API := &API{
			chain: c,
			getCommittee: func(header *types.Header, chain consensus.ChainReader) (types.Committee, error) {
				if header == h && chain == c {
					return want, nil
				}
				return nil, nil
			},
		}

		got, err := API.GetCommitteeAtHash(hash)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestAPIGetContractABI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	want := "CONTRACT ABI DATA"

	backend := core.NewMockBackend(ctrl)
	backend.EXPECT().GetContractABI().Return(want)

	API := &API{
		tendermint: backend,
	}

	got := API.GetContractABI()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestAPIGetContractAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	backend := core.NewMockBackend(ctrl)

	API := &API{
		tendermint: backend,
	}

	got := API.GetContractAddress()
	if !reflect.DeepEqual(got, autonity.ContractAddress) {
		t.Fatalf("want %v, got %v", autonity.ContractAddress, got)
	}
}

func TestAPIGetWhitelist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	want := []string{"d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303"}

	backend := core.NewMockBackend(ctrl)
	backend.EXPECT().WhiteList().Return(want)

	API := &API{
		tendermint: backend,
	}

	got := API.GetWhitelist()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}
