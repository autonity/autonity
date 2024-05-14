package backend

import (
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

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/consensus/misc"
	"github.com/autonity/autonity/consensus/tendermint"
	tdmcore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/crypto/blst"

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
	testKey, _  = blst.RandKey()
	testSigner  = func(data common.Hash) (blst.Signature, common.Address) {
		signature := testKey.Sign(data[:])
		return signature, testAddress
	}
	testSignatureBytes = common.Hex2Bytes("8ff38c5915e56029ace231f12e6911587fac4b5618077f3dfe8068138ff1dc7a7ea45a5e0d6a51747cc5f4d990c9d4de1242f4efa93d8165936bfe111f86aaafeea5eda0c38fa3dc2f854576dde63214d7438ea398e48072bc6a0c8e6c2830ef")
	testSignature, _   = blst.SignatureFromBytes(testSignatureBytes)
)

func newTestHeader(committeeSize int) *types.Header {
	validators := make(types.Committee, committeeSize)
	for i := 0; i < committeeSize; i++ {
		privateKey, _ := crypto.GenerateKey()
		committeeMember := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetUint64(1),
		}
		validators[i] = committeeMember
	}
	return &types.Header{
		Number:    new(big.Int).SetUint64(7),
		Committee: validators,
	}
}

func TestAskSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// We are testing for a Quorum Q of peers to be asked for sync.
	header := newTestHeader(7) // N=7, F=2, Q=5
	validators := header.Committee
	addresses := make([]common.Address, 0, len(validators))
	peers := make(map[common.Address]ethereum.Peer)
	counter := uint64(0)
	for _, val := range validators {
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
	msg := message.NewPrevote(1, 1, common.Hash{}, testSigner)

	addresses := make([]common.Address, 0, len(validators))
	peers := make(map[common.Address]ethereum.Peer)
	counter := uint64(0)
	for i, val := range validators {
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

	b.Gossip(validators, msg)
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
		block, err := backend.AddSeal(block)
		require.NoError(t, err)

		// We need to sleep to avoid verifying a block in the future
		time.Sleep(time.Duration(1) * time.Second)
		if _, err := backend.VerifyProposal(block); err != nil {
			t.Fatalf("could not verify block %d, err=%s", i, err)
		}

		// VerifyProposal does not need committed seals, but InsertChain does
		committedSeal, address := backend.Sign(message.PrepareCommittedSeal(block.Hash(), 0, block.Number()))
		if address != backend.address {
			t.Fatal("did not return signing address")
		}
		// Append committed seals into extra-data
		committedSeals := make(types.Signatures)
		committedSeals[address] = committedSeal.(*blst.BlsSignature)
		header := block.Header()
		if err := types.WriteCommittedSeals(header, committedSeals); err != nil {
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
	//Check signature verification
	publicKey := b.consensusKey.PublicKey()
	valid := sig.Verify(publicKey, data.Bytes())
	require.True(t, valid)
}

func TestCommit(t *testing.T) {
	t.Run("Broadcaster is not set", func(t *testing.T) {
		_, backend := newBlockChain(4)

		commitCh := make(chan *types.Block, 1)
		backend.SetResultChan(commitCh)

		// signature is not verified when committing, therefore we can just insert a bogus sig
		seals := make(types.Signatures)
		seals[backend.address] = testSignature.(*blst.BlsSignature)

		// Case: it's a proposer, so the Backend.commit will receive channel result from Backend.Commit function
		testCases := []struct {
			expectedErr   error
			expectedSeals types.Signatures
			expectedBlock func() *types.Block
		}{
			{
				// normal case
				nil,
				seals,
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
			if err := backend.Commit(expBlock, 0, test.expectedSeals); err != nil {
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

		broadcaster := consensus.NewMockBroadcaster(ctrl)
		enqueuer := consensus.NewMockEnqueuer(ctrl)
		enqueuer.EXPECT().Enqueue(fetcherID, gomock.Any())

		gossiper := interfaces.NewMockGossiper(ctrl)
		gossiper.EXPECT().SetBroadcaster(broadcaster).Times(1)
		b := &Backend{
			Broadcaster: broadcaster,
			gossiper:    gossiper,
			logger:      log.New("backend", "test", "id", 0),
		}
		b.SetBroadcaster(broadcaster)
		b.SetEnqueuer(enqueuer)

		// signature is not verified when committing, therefore we can just insert a bogus sig
		seals := make(types.Signatures)
		seals[b.address] = testSignature.(*blst.BlsSignature)

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
			message.NewPrevote(7, 8, common.HexToHash("0x1227"), testSigner),
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
	genesis, nodeKeys, consensusKeys := getGenesisAndKeys(n)

	memDB := rawdb.NewMemoryDatabase()
	msgStore := new(tdmcore.MsgStore)
	// Use the first key as private key
	b := New(nodeKeys[0], consensusKeys[0], &vm.Config{}, nil, new(event.TypeMux), msgStore, log.Root(), false)
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

func getGenesisAndKeys(n int) (*core.Genesis, []*ecdsa.PrivateKey, []blst.SecretKey) {
	genesis := core.DefaultGenesisBlock()
	// Setup committee
	var nodeKeys = make([]*ecdsa.PrivateKey, n)
	var addrs = make([]common.Address, n)
	var consensusKeys = make([]blst.SecretKey, n)
	for i := 0; i < n; i++ {
		nodeKeys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
		genesis.Alloc[addrs[i]] = core.GenesisAccount{Balance: new(big.Int).SetUint64(uint64(math.Pow10(18)))}
		consensusKey, err := blst.RandKey()
		if err != nil {
			panic(err)
		}
		consensusKeys[i] = consensusKey
	}

	// generate genesis block

	genesis.Config = params.TestChainConfig
	genesis.Config.AutonityContractConfig.Validators = nil
	genesis.Config.Ethash = nil
	genesis.GasLimit = 10000000
	genesis.Nonce = emptyNonce.Uint64()
	genesis.Mixhash = types.BFTDigest
	genesis.Timestamp = 1

	AppendValidators(genesis, nodeKeys, consensusKeys)
	err := genesis.Config.AutonityContractConfig.Prepare()
	if err != nil {
		panic(err)
	}

	return genesis, nodeKeys, consensusKeys
}

func AppendValidators(genesis *core.Genesis, keys []*ecdsa.PrivateKey, consensusKeys []blst.SecretKey) {
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
		oracleKey, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}

		genesis.Config.AutonityContractConfig.Validators = append(
			genesis.Config.AutonityContractConfig.Validators,
			&params.Validator{
				NodeAddress:   &nodeAddr,
				OracleAddress: crypto.PubkeyToAddress(oracleKey.PublicKey),
				Treasury:      nodeAddr,
				Enode:         node.URLv4(),
				BondedStake:   new(big.Int).SetUint64(100),
				ConsensusKey:  consensusKeys[i].PublicKey().Marshal(),
			})
	}
}

func makeHeader(parent *types.Block, feeGetter misc.BaseFeeGetter) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent.GasLimit(), 8000000),
		GasUsed:    0,
		BaseFee:    misc.CalcBaseFee(params.TestChainConfig, parent.Header(), feeGetter),
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
	engine.SetResultChan(resultCh)
	err = engine.Seal(chain, block, resultCh, nil)
	if err != nil {
		return nil, err
	}

	return <-resultCh, nil
}

func makeBlockWithoutSeal(chain *core.BlockChain, engine *Backend, parent *types.Block) (*types.Block, error) {
	header := makeHeader(parent, chain)
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
		tx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1337)), engine.nodeKey)
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
