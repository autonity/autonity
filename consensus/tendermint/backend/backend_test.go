// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	lru "github.com/hashicorp/golang-lru"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
)

func TestAskSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// We are testing for a Quorum Q of peers to be asked for sync.
	header := newTestHeader(7) // N=7, F=2, Q=5
	validators := header.Committee
	addresses := make([]common.Address, 0, len(validators))
	peers := make(map[common.Address]consensus.Peer)
	counter := uint64(0)
	for _, val := range validators {
		addresses = append(addresses, val.Address)
		mockedPeer := consensus.NewMockPeer(ctrl)
		mockedPeer.EXPECT().Send(uint64(TendermintSyncMsg), gomock.Eq([]byte{})).Do(func(_, _ interface{}) {
			atomic.AddUint64(&counter, 1)
		}).MaxTimes(1)
		peers[val.Address] = mockedPeer
	}

	m := make(map[common.Address]struct{})
	for _, p := range addresses {
		m[p] = struct{}{}
	}
	knownMessages, err := lru.NewARC(inmemoryMessages)
	if err != nil {
		t.Fatalf("Expected <nil>, got %v", err)
	}

	broadcaster := consensus.NewMockBroadcaster(ctrl)
	broadcaster.EXPECT().FindPeers(m).Return(peers)
	b := &Backend{
		knownMessages: knownMessages,
		logger:        log.New("backend", "test", "id", 0),
	}
	b.SetBroadcaster(broadcaster)
	b.AskSync(header)
	<-time.NewTimer(2 * time.Second).C
	if atomic.LoadUint64(&counter) != 5 {
		t.Fatalf("ask sync message transmission failure")
	}
}

func TestGossip(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	header := newTestHeader(5)
	validators := header.Committee
	payload, err := rlp.EncodeToBytes([]byte("data"))
	hash := types.RLPHash(payload)
	if err != nil {
		t.Fatalf("Expected <nil>, got %v", err)
	}
	addresses := make([]common.Address, 0, len(validators))
	peers := make(map[common.Address]consensus.Peer)
	counter := uint64(0)
	for i, val := range validators {
		addresses = append(addresses, val.Address)
		mockedPeer := consensus.NewMockPeer(ctrl)
		// Address n3 is supposed to already have this message
		if i == 3 {
			mockedPeer.EXPECT().Send(gomock.Any(), gomock.Any()).Times(0)
		} else {
			mockedPeer.EXPECT().Send(gomock.Any(), gomock.Any()).Do(func(msgCode, data interface{}) {
				// We want to make sure the payload is correct AND that no other messages is sent.
				if msgCode == uint64(TendermintMsg) && reflect.DeepEqual(data, payload) {
					atomic.AddUint64(&counter, 1)
				}
			}).Times(1)
		}
		peers[val.Address] = mockedPeer
	}

	m := make(map[common.Address]struct{})
	for _, p := range addresses {
		m[p] = struct{}{}
	}

	broadcaster := consensus.NewMockBroadcaster(ctrl)
	broadcaster.EXPECT().FindPeers(m).Return(peers)
	broadcaster.EXPECT().BroadcastToLocal(payload)

	knownMessages, err := lru.NewARC(inmemoryMessages)
	if err != nil {
		t.Fatalf("Expected <nil>, got %v", err)
	}
	recentMessages, err := lru.NewARC(inmemoryMessages)
	if err != nil {
		t.Fatalf("Expected <nil>, got %v", err)
	}
	address3Cache, err := lru.NewARC(inmemoryMessages)
	if err != nil {
		t.Fatalf("Expected <nil>, got %v", err)
	}
	address3Cache.Add(hash, true)
	recentMessages.Add(addresses[3], address3Cache)
	b := &Backend{
		knownMessages:  knownMessages,
		recentMessages: recentMessages,
	}
	b.SetBroadcaster(broadcaster)

	b.Gossip(context.Background(), validators, payload)
	<-time.NewTimer(2 * time.Second).C
	if atomic.LoadUint64(&counter) != 4 {
		t.Fatalf("gossip message transmission failure")
	}
}

func TestVerifyProposal(t *testing.T) {
	blockchain, backend := newBlockChain(1)
	blocks := make([]*types.Block, 5)

	for i := range blocks {
		var parent *types.Block
		if i == 0 {
			parent = blockchain.Genesis()
		} else {
			parent = blocks[i-1]
		}

		block, errBlock := makeBlockWithoutSeal(blockchain, backend, parent)
		if errBlock != nil {
			t.Fatalf("could not create block %d, err=%s", i, errBlock)
		}
		header := block.Header()

		seal, errS := backend.Sign(types.SigHash(header).Bytes())
		if errS != nil {
			t.Fatalf("could not sign %d, err=%s", i, errS)
		}
		if err := types.WriteSeal(header, seal); err != nil {
			t.Fatalf("could not write seal %d, err=%s", i, err)
		}
		block = block.WithSeal(header)

		// We need to sleep to avoid verifying a block in the future
		time.Sleep(time.Duration(backend.config.BlockPeriod) * time.Second)
		if _, err := backend.VerifyProposal(*block); err != nil {
			t.Fatalf("could not verify block %d, err=%s", i, err)
		}
		// VerifyProposal dont need committed seals
		committedSeal, errSC := backend.Sign(PrepareCommittedSeal(block.Hash()))
		if errSC != nil {
			t.Fatalf("could not sign commit %d, err=%s", i, errS)
		}
		// Append seals into extra-data
		if err := types.WriteCommittedSeals(header, [][]byte{committedSeal}); err != nil {
			t.Fatalf("could not write committed seal %d, err=%s", i, err)
		}
		block = block.WithSeal(header)

		state, stateErr := blockchain.State()
		if stateErr != nil {
			t.Fatalf("could not retrieve state %d, err=%s", i, stateErr)
		}
		if status, errW := blockchain.WriteBlockWithState(block, nil, nil, state, false); status != core.CanonStatTy && errW != nil {
			t.Fatalf("write block failure %d, err=%s", i, errW)
		}
		blocks[i] = block
	}

}
func TestResetPeerCache(t *testing.T) {
	addr := common.HexToAddress("0x01234567890")
	msgCache, err := lru.NewARC(inmemoryMessages)
	if err != nil {
		t.Fatalf("Expected <nil>, got %v", err)
	}
	msgCache.Add(addr, addr)

	recentMessages, err := lru.NewARC(inmemoryMessages)
	if err != nil {
		t.Fatalf("Expected <nil>, got %v", err)
	}
	recentMessages.Add(addr, msgCache)

	b := &Backend{
		recentMessages: recentMessages,
	}

	b.ResetPeerCache(addr)
	if msgCache.Contains(addr) {
		t.Fatalf("expected empty cache")
	}
}

func TestHasBadProposal(t *testing.T) {
	t.Run("callback is not set, false returned", func(t *testing.T) {
		b := &Backend{}
		if b.HasBadProposal(common.HexToHash("0x01234567890")) {
			t.Fatalf("expected <false>, got <true>")
		}
	})

	t.Run("callback is set, true returned", func(t *testing.T) {
		b := &Backend{
			hasBadBlock: func(hash common.Hash) bool {
				return true
			},
		}
		if !b.HasBadProposal(common.HexToHash("0x01234567890")) {
			t.Fatalf("expected <true>, got <false>")
		}
	})
}

func TestSign(t *testing.T) {
	_, b := newBlockChain(4)
	data := []byte("Here is a string....")
	sig, err := b.Sign(data)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	//Check signature recover
	hashData := crypto.Keccak256(data)
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	if signer != b.address {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
	}
}

func TestCheckSignature(t *testing.T) {
	key, _ := generatePrivateKey()
	data := []byte("Here is a string....")
	hashData := crypto.Keccak256(data)
	sig, _ := crypto.Sign(hashData, key)
	_, b := newBlockChain(4)
	a := getAddress()
	err := b.CheckSignature(data, a, sig)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	a = getInvalidAddress()
	err = b.CheckSignature(data, a, sig)
	if err != types.ErrInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, types.ErrInvalidSignature)
	}
}

func TestCommit(t *testing.T) {
	t.Run("broadcaster is not set", func(t *testing.T) {
		_, backend := newBlockChain(4)

		commitCh := make(chan *types.Block, 1)
		backend.setResultChan(commitCh)

		// Case: it's a proposer, so the Backend.commit will receive channel result from Backend.Commit function
		testCases := []struct {
			expectedErr       error
			expectedSignature [][]byte
			expectedBlock     func() *types.Block
		}{
			{
				// normal case
				nil,
				[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.BFTExtraSeal-1)...)},
				func() *types.Block {
					chain, engine := newBlockChain(1)
					block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
					if err != nil {
						t.Fatal(err)
					}
					expectedBlock, _ := engine.AddSeal(block)
					return expectedBlock
				},
			},
			{
				// invalid signature
				types.ErrInvalidCommittedSeals,
				nil,
				func() *types.Block {
					chain, engine := newBlockChain(1)
					block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
					if err != nil {
						t.Fatal(err)
					}
					expectedBlock, _ := engine.AddSeal(block)
					return expectedBlock
				},
			},
		}

		for _, test := range testCases {
			expBlock := test.expectedBlock()

			backend.proposedBlockHash = expBlock.Hash()
			if err := backend.Commit(expBlock, 0, test.expectedSignature); err != nil {
				if err != test.expectedErr {
					t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
				}
			}

			if test.expectedErr == nil {
				// to avoid race condition is occurred by goroutine
				select {
				case result := <-commitCh:
					if result.Hash() != expBlock.Hash() {
						t.Errorf("hash mismatch: have %v, want %v", result.Hash(), expBlock.Hash())
					}
				case <-time.After(10 * time.Second):
					t.Fatal("timeout")
				}
			}
		}
	})

	t.Run("broadcaster is set", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		blockFactory := func() *types.Block {
			chain, engine := newBlockChain(1)
			block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
			if err != nil {
				t.Fatal(err)
			}
			expectedBlock, _ := engine.AddSeal(block)
			return expectedBlock
		}

		newBlock := blockFactory()
		seals := [][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.BFTExtraSeal-1)...)}

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().Enqueue(fetcherID, gomock.Any())

		b := &Backend{
			broadcaster: broadcaster,
			logger:      log.New("backend", "test", "id", 0),
		}
		b.SetBroadcaster(broadcaster)

		err := b.Commit(newBlock, 0, seals)
		if err != nil {
			t.Fatalf("expected <nil>, got %v", err)
		}
	})
}

func TestSyncPeer(t *testing.T) {
	t.Run("no broadcaster set, nothing done", func(t *testing.T) {
		b := &Backend{}
		b.SyncPeer(common.HexToAddress("0x0123456789"))
	})

	t.Run("valid params given, messages sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		peerAddr1 := common.HexToAddress("0x0123456789")
		messages := []*tendermintCore.Message{
			{
				Address: peerAddr1,
			},
		}

		peersAddrMap := make(map[common.Address]struct{})
		peersAddrMap[peerAddr1] = struct{}{}

		payload := messages[0].Payload()

		peer1Mock := consensus.NewMockPeer(ctrl)
		peer1Mock.EXPECT().Send(uint64(TendermintMsg), payload)

		peers := make(map[common.Address]consensus.Peer)
		peers[peerAddr1] = peer1Mock

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeers(peersAddrMap).Return(peers)

		recentMessages, err := lru.NewARC(inmemoryPeers)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		tendermintC := tendermintCore.NewMockTendermint(ctrl)
		tendermintC.EXPECT().GetCurrentHeightMessages().Return(messages)

		b := &Backend{
			logger:         log.New("backend", "test", "id", 0),
			recentMessages: recentMessages,
			core:           tendermintC,
		}
		b.SetBroadcaster(broadcaster)

		b.SyncPeer(peerAddr1)

		wait := time.NewTimer(time.Second)
		<-wait.C
	})
}

func TestBackendLastCommittedProposal(t *testing.T) {
	t.Run("block number 0, block returned", func(t *testing.T) {
		block := types.NewBlockWithHeader(&types.Header{})

		b := &Backend{
			currentBlock: func() *types.Block {
				return block
			},
			logger: log.New("backend", "test", "id", 0),
		}

		bl, _ := b.LastCommittedProposal()
		if !reflect.DeepEqual(bl, block) {
			t.Fatalf("expected %v, got %v", block, bl)
		}
	})

	t.Run("block number is greater than 0, empty block returned", func(t *testing.T) {
		block := types.NewBlockWithHeader(&types.Header{
			Number: big.NewInt(1),
		})

		b := &Backend{
			currentBlock: func() *types.Block {
				return block
			},
			logger: log.New("backend", "test", "id", 0),
		}

		bl, _ := b.LastCommittedProposal()
		if !reflect.DeepEqual(bl, &types.Block{}) {
			t.Fatalf("expected empty block, got %v", bl)
		}
	})
}

// Test get contract ABI, it should have the default abi before contract upgrade.
func TestBackendGetContractABI(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	_, err = chain.InsertChain(types.Blocks{block})
	if err != nil {
		t.Fatal(err)
	}
	contractABI := engine.GetContractABI()
	expectedABI := chain.Config().AutonityContractConfig.ABI
	if contractABI != expectedABI {
		t.Fatalf("unexpected returned ABI")
	}
}

func TestBackendWhiteList(t *testing.T) {
	//Very shallow test for the time being, running only with 1 validator
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}
	_, err = chain.InsertChain(types.Blocks{block})
	if err != nil {
		t.Fatal(err)
	}
	whitelist := engine.WhiteList()
	if len(whitelist) != 1 {
		t.Fatalf("unexpected returned whitelist")
	}
	expectedWhitelist := chain.Config().AutonityContractConfig.Users[0].Enode
	if strings.Compare(whitelist[0], expectedWhitelist) != 0 {
		t.Fatalf("unexpected returned whitelist")
	}
}

/**
 * SimpleBackend
 * Private key: bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1
 * Public key: 04a2bfb0f7da9e1b9c0c64e14f87e8fb82eb0144e97c25fe3a977a921041a50976984d18257d2495e7bfd3d4b280220217f429287d25ecdf2b0d7c0f7aae9aa624
 * Address: 0x70524d664ffe731100208a0154e556f9bb679ae6
 */
func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func getInvalidAddress() common.Address {
	return common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
}

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func newTestHeader(n int) *types.Header {
	// generate committee
	addrs := make(types.Committee, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		addrs[i] = types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetUint64(1),
		}
	}
	h := &types.Header{
		Committee: addrs,
	}
	return h
}

type Keys []*ecdsa.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func newBlockChain(n int) (*core.BlockChain, *Backend) {
	genesis, nodeKeys := getGenesisAndKeys(n)
	memDB := rawdb.NewMemoryDatabase()
	// Use the first key as private key
	b := New(genesis.Config.Tendermint, nodeKeys[0], memDB, genesis.Config, &vm.Config{})

	genesis.MustCommit(memDB)
	blockchain, err := core.NewBlockChain(memDB, nil, genesis.Config, b, vm.Config{}, nil, core.NewTxSenderCacher(), nil)
	if err != nil {
		panic(err)
	}
	b.SetBlockchain(blockchain)

	err = b.Start(context.Background())
	if err != nil {
		panic(err)
	}

	return blockchain, b
}

func getGenesisAndKeys(n int) (*core.Genesis, []*ecdsa.PrivateKey) {
	genesis := core.DefaultGenesisBlock()
	// Setup committee
	var nodeKeys = make([]*ecdsa.PrivateKey, n)
	var addrs = make([]common.Address, n)
	for i := 0; i < n; i++ {
		nodeKeys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
		genesis.Alloc[addrs[i]] = core.GenesisAccount{Balance: new(big.Int).SetUint64(uint64(math.Pow10(18)))}
	}

	// generate genesis block

	genesis.Config = params.TestChainConfig
	genesis.GasLimit = 10000000
	genesis.Config.AutonityContractConfig = &params.AutonityContractGenesis{}
	// force enable Tendermint engine
	genesis.Config.Tendermint = config.DefaultConfig()
	genesis.Difficulty = defaultDifficulty
	genesis.Nonce = emptyNonce.Uint64()
	genesis.Mixhash = types.BFTDigest

	AppendValidators(genesis, addrs)
	err := genesis.Config.AutonityContractConfig.Prepare()
	if err != nil {
		panic(err)
	}

	return genesis, nodeKeys
}

const EnodeStub = "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303"

func AppendValidators(genesis *core.Genesis, addrs []common.Address) {
	if genesis.Config == nil {
		genesis.Config = &params.ChainConfig{}
	}
	if genesis.Config.AutonityContractConfig == nil {
		genesis.Config.AutonityContractConfig = &params.AutonityContractGenesis{}
	}

	for i := range addrs {
		genesis.Config.AutonityContractConfig.Users = append(
			genesis.Config.AutonityContractConfig.Users,
			params.User{
				Address: &addrs[i],
				Type:    params.UserValidator,
				Enode:   EnodeStub,
				Stake:   100,
			})
	}
}

func makeHeader(parent *types.Block, config *config.Config) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, 8000000, 8000000),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(big.NewInt(int64(parent.Time())), new(big.Int).SetUint64(config.BlockPeriod)).Uint64(),
		Difficulty: defaultDifficulty,
		MixDigest:  types.BFTDigest,
		Round:      0,
	}
	return header
}

func makeBlock(chain *core.BlockChain, engine *Backend, parent *types.Block) (*types.Block, error) {
	block, err := makeBlockWithoutSeal(chain, engine, parent)
	if err != nil {
		return nil, err
	}

	resultCh := make(chan *types.Block)
	err = engine.Seal(chain, block, resultCh, nil)
	if err != nil {
		return nil, err
	}

	return <-resultCh, nil
}

func makeBlockWithoutSeal(chain *core.BlockChain, engine *Backend, parent *types.Block) (*types.Block, error) {
	header := makeHeader(parent, engine.config)
	_ = engine.Prepare(chain, header)

	state, errS := chain.StateAt(parent.Root())
	if errS != nil {
		return nil, errS
	}

	//add a few txs
	txs := make(types.Transactions, 5)
	nonce := state.GetNonce(engine.address)
	gasPrice := new(big.Int).SetUint64(1000000)
	gasPool := new(core.GasPool).AddGas(header.GasLimit)
	var receipts []*types.Receipt
	for i := range txs {
		amount := new(big.Int).SetUint64((nonce + 1) * 1000000000)
		tx := types.NewTransaction(nonce, common.Address{}, amount, params.TxGas, gasPrice, []byte{})
		tx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), engine.privateKey)
		if err != nil {
			return nil, err
		}
		txs[i] = tx
		receipt, err := core.ApplyTransaction(chain.Config(), chain, nil, gasPool, state, header, txs[i], &header.GasUsed, *engine.vmConfig)
		if err != nil {
			return nil, err
		}
		nonce++
		receipts = append(receipts, receipt)
	}
	block, err := engine.FinalizeAndAssemble(chain, header, state, txs, nil, &receipts)
	if err != nil {
		return nil, err
	}

	// Write state changes to db
	root, err := state.Commit(chain.Config().IsEIP158(block.Header().Number))
	if err != nil {
		return nil, fmt.Errorf("state write error: %v", err)
	}
	if err := state.Database().TrieDB().Commit(root, false, nil); err != nil {
		return nil, fmt.Errorf("trie write error: %v", err)
	}

	return block, nil
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(2)})
	return buf.Bytes()
}
