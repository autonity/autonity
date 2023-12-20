package backend

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"os"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/consensus/misc"
	"github.com/autonity/autonity/consensus/tendermint"
	tdmcore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	lru "github.com/hashicorp/golang-lru"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
)

var (
	testAddress = common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
	testKey, _  = crypto.HexToECDSA("bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1")
	testSigner  = func(data common.Hash) ([]byte, common.Address) {
		out, _ := crypto.Sign(data[:], testKey)
		return out, testAddress
	}
)

func newCommittee(committeeSize int) *types.Committee {
	c := new(types.Committee)
	for i := 0; i < committeeSize; i++ {
		privateKey, _ := crypto.GenerateKey()
		consensusKey, _ := blst.RandKey()
		committeeMember := types.CommitteeMember{
			Address:      crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower:  new(big.Int).SetUint64(1),
			ConsensusKey: consensusKey.PublicKey().Marshal(),
		}
		c.Members = append(c.Members, &committeeMember)
	}
	c.Sort()
	return c
}

func TestAskSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// We are testing for a Quorum Q of peers to be asked for sync.
	committee := newCommittee(7) // N=7, F=2, Q=5
	addresses := make([]common.Address, 0, committee.Len())
	peers := make(map[common.Address]ethereum.Peer)
	counter := uint64(0)
	for _, val := range committee.Members {
		addresses = append(addresses, val.Address)
		mockedPeer := tendermint.NewMockPeer(ctrl)
		mockedPeer.EXPECT().Send(SyncNetworkMsg, gomock.Eq([]byte{})).Do(func(_, _ interface{}) {
			atomic.AddUint64(&counter, 1)
		}).MaxTimes(1)
		peers[val.Address] = mockedPeer
	}

	m := make(map[common.Address]struct{})
	for _, p := range addresses {
		m[p] = struct{}{}
	}
	knownMessages, err := lru.NewARC(inmemoryMessages)
	require.NoError(t, err)
	recentMessages, err := lru.NewARC(inmemoryMessages)
	require.NoError(t, err)

	broadcaster := consensus.NewMockBroadcaster(ctrl)
	broadcaster.EXPECT().FindPeers(m).Return(peers)
	b := &Backend{
		knownMessages: knownMessages,
		gossiper:      NewGossiper(recentMessages, knownMessages, common.Address{}, log.New(), make(chan struct{})),
		logger:        log.New("backend", "test", "id", 0),
	}
	b.SetBroadcaster(broadcaster)
	b.AskSync(committee)
	<-time.NewTimer(2 * time.Second).C
	if atomic.LoadUint64(&counter) != 5 {
		t.Fatalf("ask sync message transmission failure")
	}
}

func TestGossip(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	committee := newCommittee(5)
	msg := message.NewPrevote(1, 1, common.Hash{}, testSigner)

	addresses := make([]common.Address, 0, committee.Len())
	peers := make(map[common.Address]ethereum.Peer)
	counter := uint64(0)
	for i, val := range committee.Members {
		addresses = append(addresses, val.Address)
		mockedPeer := tendermint.NewMockPeer(ctrl)
		// Address n3 is supposed to already have this message
		if i == 3 {
			mockedPeer.EXPECT().SendRaw(gomock.Any(), gomock.Any()).Times(0)
		} else {
			mockedPeer.EXPECT().SendRaw(gomock.Any(), gomock.Any()).Do(func(msgCode, data interface{}) {
				// We want to make sure the payload is correct AND that no other messages is sent.
				if msgCode == PrevoteNetworkMsg && reflect.DeepEqual(data, msg.Payload()) {
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

	knownMessages, err := lru.NewARC(inmemoryMessages)
	require.NoError(t, err)
	recentMessages, err := lru.NewARC(inmemoryMessages)
	require.NoError(t, err)
	address3Cache, err := lru.NewARC(inmemoryMessages)
	require.NoError(t, err)

	address3Cache.Add(msg.Hash(), true)
	recentMessages.Add(addresses[3], address3Cache)
	b := &Backend{
		knownMessages:  knownMessages,
		recentMessages: recentMessages,
		gossiper:       NewGossiper(recentMessages, knownMessages, common.Address{}, log.New(), make(chan struct{})),
	}
	b.SetBroadcaster(broadcaster)

	b.Gossip(committee, msg)
	<-time.NewTimer(2 * time.Second).C
	if c := atomic.LoadUint64(&counter); c != 4 {
		t.Fatal("Gossip message transmission failure", "have", c, "want", 4)
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
		seal, _ := backend.Sign(types.SigHash(header))
		if err := types.WriteSeal(header, seal); err != nil {
			t.Fatalf("could not write seal %d, err=%s", i, err)
		}
		block = block.WithSeal(header)

		// We need to sleep to avoid verifying a block in the future
		time.Sleep(time.Duration(1) * time.Second)
		if _, err := backend.VerifyProposal(block); err != nil {
			t.Fatalf("could not verify block %d, err=%s", i, err)
		}
		// VerifyProposal don't need committed seals
		committedSeal, address := backend.Sign(message.PrepareCommittedSeal(block.Hash(), 0, block.Number()))
		if address != backend.address {
			t.Fatal("did not return signing address")
		}
		// Append seals into extra-data
		if err := types.WriteCommittedSeals(header, [][]byte{committedSeal}); err != nil {
			t.Fatalf("could not write committed seal %d, err=%s", i, err)
		}
		block = block.WithSeal(header)

		if _, errW := blockchain.InsertChain(types.Blocks{block}); errW != nil {
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
	data := common.HexToHash("0x12345")
	sig, addr := b.Sign(data)
	if addr != b.address {
		t.Error("error mismatch of addresses")
	}
	//Check signature recovery
	signer, _ := crypto.SigToAddr(data[:], sig)
	if signer != b.address {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), testAddress)
	}
}

func TestCommit(t *testing.T) {
	t.Run("Broadcaster is not set", func(t *testing.T) {
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

	t.Run("Broadcaster is set", func(t *testing.T) {
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

		gossiper := interfaces.NewMockGossiper(ctrl)
		gossiper.EXPECT().SetBroadcaster(broadcaster).Times(1)
		b := &Backend{
			Broadcaster: broadcaster,
			gossiper:    gossiper,
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
	t.Run("no Broadcaster set, nothing done", func(t *testing.T) {
		b := &Backend{}
		b.SyncPeer(common.HexToAddress("0x0123456789"))
	})

	t.Run("valid params given, messages sent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		peerAddr1 := common.HexToAddress("0x0123456789")
		messages := []message.Msg{
			message.NewPrevote(7, 8, common.HexToHash("0x1227"), dummySigner),
		}

		peersAddrMap := make(map[common.Address]struct{})
		peersAddrMap[peerAddr1] = struct{}{}

		payload := messages[0].Payload()

		peer1Mock := tendermint.NewMockPeer(ctrl)
		peer1Mock.EXPECT().SendRaw(PrevoteNetworkMsg, payload)

		peers := make(map[common.Address]ethereum.Peer)
		peers[peerAddr1] = peer1Mock

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		broadcaster.EXPECT().FindPeers(peersAddrMap).Return(peers)

		recentMessages, err := lru.NewARC(inmemoryPeers)
		if err != nil {
			t.Fatalf("Expected <nil>, got %v", err)
		}

		tendermintC := interfaces.NewMockCore(ctrl)
		tendermintC.EXPECT().CurrentHeightMessages().Return(messages)

		gossiper := interfaces.NewMockGossiper(ctrl)
		gossiper.EXPECT().SetBroadcaster(broadcaster).Times(1)
		b := &Backend{
			logger:         log.New("backend", "test", "id", 0),
			gossiper:       gossiper,
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
	t.Run("return current block", func(t *testing.T) {
		block := types.NewBlockWithHeader(&types.Header{})

		b := &Backend{
			currentBlock: func() *types.Block {
				return block
			},
			logger: log.New("backend", "test", "id", 0),
		}

		bl := b.HeadBlock()
		if !reflect.DeepEqual(bl, block) {
			t.Fatalf("expected %v, got %v", block, bl)
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

// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func newBlockChain(n int) (*core.BlockChain, *Backend) {
	genesis, nodeKeys := getGenesisAndKeys(n)

	memDB := rawdb.NewMemoryDatabase()
	msgStore := new(tdmcore.MsgStore)
	// Use the first key as private key
	b := New(nodeKeys[0], &vm.Config{}, nil, new(event.TypeMux), msgStore, log.Root())
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	genesis.MustCommit(memDB)
	blockchain, err := core.NewBlockChain(memDB, nil, genesis.Config, b, vm.Config{}, nil, core.NewTxSenderCacher(), nil, backends.NewInternalBackend(nil), log.Root())
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
	genesis.Config.AutonityContractConfig.Validators = nil
	genesis.Config.Ethash = nil
	genesis.GasLimit = 10000000
	genesis.Nonce = emptyNonce.Uint64()
	genesis.Mixhash = types.BFTDigest
	genesis.Timestamp = 1

	AppendValidators(genesis, nodeKeys)
	err := genesis.Config.AutonityContractConfig.Prepare()
	if err != nil {
		panic(err)
	}

	return genesis, nodeKeys
}

func AppendValidators(genesis *core.Genesis, keys []*ecdsa.PrivateKey) {
	if genesis.Config == nil {
		genesis.Config = &params.ChainConfig{}
	}
	if genesis.Config.AutonityContractConfig == nil {
		genesis.Config.AutonityContractConfig = &params.AutonityContractGenesis{}
	}

	if genesis.Config.OracleContractConfig == nil {
		genesis.Config.OracleContractConfig = &params.OracleContractGenesis{}
	}

	for i := range keys {
		nodeAddr := crypto.PubkeyToAddress(keys[i].PublicKey)
		node := enode.NewV4(&keys[i].PublicKey, nil, 0, 0)
		blsKey, err := blst.RandKey()
		if err != nil {
			panic(err)
		}
		oracleKey, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		treasuryAddr := nodeAddr
		pop, err := crypto.AutonityPOPProof(keys[i], oracleKey, treasuryAddr.Hex(), blsKey)
		if err != nil {
			panic(err)
		}

		genesis.Config.AutonityContractConfig.Validators = append(
			genesis.Config.AutonityContractConfig.Validators,
			&params.Validator{
				NodeAddress:   &nodeAddr,
				OracleAddress: crypto.PubkeyToAddress(oracleKey.PublicKey),
				Pop:           pop,
				Treasury:      nodeAddr,
				Enode:         node.URLv4(),
				BondedStake:   new(big.Int).SetUint64(100),
				ConsensusKey:  blsKey.PublicKey().Marshal(),
			})
	}
}

func makeHeader(parent *types.Block) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent.GasLimit(), 8000000),
		GasUsed:    0,
		BaseFee:    misc.CalcBaseFee(params.TestChainConfig, parent.Header(), nil),
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(big.NewInt(int64(parent.Time())), new(big.Int).SetUint64(1)).Uint64(),
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
	header := makeHeader(parent)
	_ = engine.Prepare(chain, header)

	state, errS := chain.StateAt(parent.Root())
	if errS != nil {
		return nil, errS
	}

	//add a few txs
	txs := make(types.Transactions, 5)
	nonce := state.GetNonce(engine.address)
	gasPrice := new(big.Int).Set(header.BaseFee)
	gasPool := new(core.GasPool).AddGas(header.GasLimit)
	var receipts []*types.Receipt
	for i := range txs {
		amount := new(big.Int).SetUint64((nonce + 1) * 1000000000)
		tx := types.NewTransaction(nonce, common.Address{}, amount, params.TxGas, gasPrice, []byte{})
		tx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1337)), engine.privateKey)
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

func dummySigner(_ common.Hash) ([]byte, common.Address) {
	return nil, common.Address{}
}
