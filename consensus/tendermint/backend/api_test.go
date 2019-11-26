package backend

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rpc"
)

func TestGetValidators(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addr := common.HexToAddress("0x0123456789")
	want := []common.Address{addr}

	val := validator.NewMockValidator(ctrl)
	val.EXPECT().Address().Return(addr)

	valSet := validator.NewMockSet(ctrl)
	valSet.EXPECT().List().Return([]validator.Validator{val})

	backend := core.NewMockBackend(ctrl)
	backend.EXPECT().Validators(uint64(1)).Return(valSet)

	API := &API{
		tendermint: backend,
	}

	bn := rpc.BlockNumber(1)

	got, err := API.GetValidators(&bn)
	if err != nil {
		t.Fatalf("expected <nil>, got %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestGetValidatorsAtHash(t *testing.T) {
	t.Run("unknown block given, error returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		hash := common.HexToHash("0x0123456789")

		chain := consensus.NewMockChainReader(ctrl)
		chain.EXPECT().GetHeaderByHash(hash).Return(nil)

		API := &API{
			chain: chain,
		}

		_, err := API.GetValidatorsAtHash(hash)
		if err != errUnknownBlock {
			t.Fatalf("expected %v, got %v", errUnknownBlock, err)
		}
	})

	t.Run("valid block given, validators returned", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addr := common.HexToAddress("0x0123456789")
		want := []common.Address{addr}

		hash := common.HexToHash("0x0123456789")

		chain := consensus.NewMockChainReader(ctrl)
		chain.EXPECT().GetHeaderByHash(hash).Return(&types.Header{Number: big.NewInt(1)})

		val := validator.NewMockValidator(ctrl)
		val.EXPECT().Address().Return(addr)

		valSet := validator.NewMockSet(ctrl)
		valSet.EXPECT().List().Return([]validator.Validator{val})

		backend := core.NewMockBackend(ctrl)
		backend.EXPECT().Validators(uint64(1)).Return(valSet)

		API := &API{
			chain:      chain,
			tendermint: backend,
		}

		got, err := API.GetValidatorsAtHash(hash)
		if err != nil {
			t.Fatalf("expected <nil>, got %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("want %v, got %v", want, got)
		}
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

	want := common.HexToAddress("0x0123456789")

	backend := core.NewMockBackend(ctrl)
	backend.EXPECT().GetContractAddress().Return(want)

	API := &API{
		tendermint: backend,
	}

	got := API.GetContractAddress()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
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
