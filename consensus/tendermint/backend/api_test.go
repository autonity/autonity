package backend

import (
	"context"
	"github.com/clearmatics/autonity/common/acdefault"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/params"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/rpc"
)

func TestGetCommittee(t *testing.T) {
	genesis, nodeKeys := getGenesisAndKeys(rand.Intn(30) + 1)
	// Add non validators to the genesis users
	part1, part2 := common.HexToAddress("0xabcd1235145"), common.HexToAddress("0x124214643")
	genesis.Config.AutonityContractConfig.Users = append(genesis.Config.AutonityContractConfig.Users, []params.User{
		{Address: &part1, Type: params.UserParticipant},
		{Address: &part2, Type: params.UserStakeHolder, Stake: 2},
	}...)

	memDB := rawdb.NewMemoryDatabase()
	backend := New(genesis.Config.Tendermint, nodeKeys[0], memDB, genesis.Config, &vm.Config{})
	genesis.MustCommit(memDB)
	blockchain, err := core.NewBlockChain(memDB, nil, genesis.Config, backend, vm.Config{}, nil, core.NewTxSenderCacher())
	assert.Nil(t, err)

	err = backend.Start(context.Background(), blockchain, blockchain.CurrentBlock, blockchain.HasBadBlock)
	assert.Nil(t, err)

	want, err := backend.Committee(0)
	assert.Nil(t, err)

	backendAPI := &API{tendermint: backend}

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
		genesis, nodeKeys := getGenesisAndKeys(rand.Intn(30) + 1)
		// Add non validators to the genesis users
		part1, part2 := common.HexToAddress("0xabcd1235145"), common.HexToAddress("0x124214643")
		genesis.Config.AutonityContractConfig.Users = append(genesis.Config.AutonityContractConfig.Users, []params.User{
			{Address: &part1, Type: params.UserParticipant},
			{Address: &part2, Type: params.UserStakeHolder, Stake: 2},
		}...)

		memDB := rawdb.NewMemoryDatabase()
		backend := New(genesis.Config.Tendermint, nodeKeys[0], memDB, genesis.Config, &vm.Config{})
		genesis.MustCommit(memDB)
		blockchain, err := core.NewBlockChain(memDB, nil, genesis.Config, backend, vm.Config{}, nil, core.NewTxSenderCacher())
		assert.Nil(t, err)

		err = backend.Start(context.Background(), blockchain, blockchain.CurrentBlock, blockchain.HasBadBlock)
		assert.Nil(t, err)

		want, err := backend.Committee(0)
		assert.Nil(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		genHash := blockchain.Genesis().Hash()
		chain := consensus.NewMockChainReader(ctrl)
		chain.EXPECT().GetHeaderByHash(genHash).Return(blockchain.Genesis().Header())

		API := &API{
			chain:      chain,
			tendermint: backend,
		}

		got, err := API.GetCommitteeAtHash(genHash)
		assert.Nil(t, err)

		assert.Equal(t, want.Committee(), got)
	})
}

func TestAPIGetContractABI(t *testing.T) {
	genesis, nodeKeys := getGenesisAndKeys(rand.Intn(30) + 1)
	// Add non validators to the genesis users
	part1, part2 := common.HexToAddress("0xabcd1235145"), common.HexToAddress("0x124214643")
	genesis.Config.AutonityContractConfig.Users = append(genesis.Config.AutonityContractConfig.Users, []params.User{
		{Address: &part1, Type: params.UserParticipant},
		{Address: &part2, Type: params.UserStakeHolder, Stake: 2},
	}...)

	memDB := rawdb.NewMemoryDatabase()
	backend := New(genesis.Config.Tendermint, nodeKeys[0], memDB, genesis.Config, &vm.Config{})
	genesis.MustCommit(memDB)
	blockchain, err := core.NewBlockChain(memDB, nil, genesis.Config, backend, vm.Config{}, nil, core.NewTxSenderCacher())
	assert.Nil(t, err)

	err = backend.Start(context.Background(), blockchain, blockchain.CurrentBlock, blockchain.HasBadBlock)
	assert.Nil(t, err)

	want := acdefault.ABI()

	API := &API{
		tendermint: backend,
	}

	got := API.GetContractABI()
	assert.Equal(t, want, got)
}

//
//func TestAPIGetContractAddress(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	want := common.HexToAddress("0x0123456789")
//
//	backend := tendermintCore.NewMockBackend(ctrl)
//	backend.EXPECT().GetContractAddress().Return(want)
//
//	API := &API{
//		tendermint: backend,
//	}
//
//	got := API.GetContractAddress()
//	if !reflect.DeepEqual(got, want) {
//		t.Fatalf("want %v, got %v", want, got)
//	}
//}
//
//func TestAPIGetWhitelist(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	want := []string{"d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303"}
//
//	backend := tendermintCore.NewMockBackend(ctrl)
//	backend.EXPECT().WhiteList().Return(want)
//
//	API := &API{
//		tendermint: backend,
//	}
//
//	got := API.GetWhitelist()
//	if !reflect.DeepEqual(got, want) {
//		t.Fatalf("want %v, got %v", want, got)
//	}
//}
