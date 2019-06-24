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
	"crypto/ecdsa"
	"fmt"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/params"
	"github.com/clearmatics/autonity/rlp"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
)

func TestSign(t *testing.T) {
	b := newBackend()
	data := []byte("Here is a string....")
	sig, err := b.Sign(data)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	//Check signature recover
	hashData := crypto.Keccak256([]byte(data))
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	if signer != getAddress() {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
	}
}

func TestCheckSignature(t *testing.T) {
	key, _ := generatePrivateKey()
	data := []byte("Here is a string....")
	hashData := crypto.Keccak256([]byte(data))
	sig, _ := crypto.Sign(hashData, key)
	b := newBackend()
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

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256([]byte(data))
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		// CheckValidatorSignature should succeed
		addr, err := tendermint.CheckValidatorSignature(vset, data, sig)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		validator := vset.GetByIndex(uint64(i))
		if addr != validator.Address() {
			t.Errorf("validator address mismatch: have %v, want %v", addr, validator.Address())
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig, err := crypto.Sign(hashData, key)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := tendermint.CheckValidatorSignature(vset, data, sig)
	if err != tendermint.ErrUnauthorizedAddress {
		t.Errorf("error mismatch: have %v, want %v", err, tendermint.ErrUnauthorizedAddress)
	}
	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCommit(t *testing.T) {
	backend := newBackend()

	commitCh := make(chan *types.Block, 1)
	backend.commitCh = commitCh

	// Case: it's a proposer, so the Backend.commit will receive channel result from Backend.Commit function
	testCases := []struct {
		expectedErr       error
		expectedSignature [][]byte
		expectedBlock     func() types.Block
	}{
		{
			// normal case
			nil,
			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.PoSExtraSeal-1)...)},
			func() types.Block {
				chain, engine := newBlockChain(1)
				block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				if err != nil {
					t.Fatal(err)
				}
				expectedBlock, _ := engine.updateBlock(block)
				return *expectedBlock
			},
		},
		{
			// invalid signature
			types.ErrInvalidCommittedSeals,
			nil,
			func() types.Block {
				chain, engine := newBlockChain(1)
				block, err := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				if err != nil {
					t.Fatal(err)
				}
				expectedBlock, _ := engine.updateBlock(block)
				return *expectedBlock
			},
		},
	}

	for _, test := range testCases {
		expBlock := test.expectedBlock()

		backend.proposedBlockHash = expBlock.Hash()
		if err := backend.Commit(expBlock, test.expectedSignature); err != nil {
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
}

func TestGetProposer(t *testing.T) {
	chain, engine := newBlockChain(1)
	block, err := makeBlock(chain, engine, chain.Genesis())
	if err != nil {
		t.Fatal(err)
	}

	chain.InsertChain(types.Blocks{block})
	expected := engine.GetProposer(1)
	actual := engine.Address()
	if actual != expected {
		t.Errorf("proposer mismatch: have %v, want %v", actual.Hex(), expected.Hex())
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

func newTestValidatorSet(n int) (tendermint.ValidatorSet, []*ecdsa.PrivateKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		addrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
	vset := validator.NewSet(addrs, tendermint.RoundRobin)
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
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

func newBackend() (b *Backend) {
	_, b = newBlockChain(4)
	key, _ := generatePrivateKey()
	b.privateKey = key
	return
}


// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func newBlockChain(n int) (*core.BlockChain, *Backend) {
	genesis, nodeKeys := getGenesisAndKeys(n)
	memDB := ethdb.NewMemDatabase()
	config := tendermint.DefaultConfig
	// Use the first key as private key
	b := New(config, nodeKeys[0], memDB, genesis.Config, &vm.Config{})
	genesis.MustCommit(memDB)
	blockchain, err := core.NewBlockChain(memDB, nil, genesis.Config, b, vm.Config{}, nil)
	if err != nil {
		panic(err)
	}
	b.Start(blockchain, blockchain.CurrentBlock, blockchain.HasBadBlock)
	validators := b.Validators(0)
	if validators.Size() == 0 {
		panic("failed to get validators")
	}
	proposerAddr := validators.GetProposer().Address()

	// find proposer key
	for _, key := range nodeKeys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		if addr.String() == proposerAddr.String() {
			b.privateKey = key
			b.address = addr
		}
	}

	return blockchain, b
}

func getGenesisAndKeys(n int) (*core.Genesis, []*ecdsa.PrivateKey) {
	// Setup validators
	var nodeKeys = make([]*ecdsa.PrivateKey, n)
	var addrs = make([]common.Address, n)
	for i := 0; i < n; i++ {
		nodeKeys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
	}

	// generate genesis block
	genesis := core.DefaultGenesisBlock()
	genesis.Config = params.TestChainConfig
	// force enable Istanbul engine
	genesis.Config.Tendermint = &params.TendermintConfig{}
	genesis.Config.Ethash = nil
	genesis.Difficulty = defaultDifficulty
	genesis.Nonce = emptyNonce.Uint64()
	genesis.Mixhash = types.PoSDigest

	AppendValidators(genesis, addrs)
	return genesis, nodeKeys
}

func AppendValidators(genesis *core.Genesis, addrs []common.Address) {

	if len(genesis.ExtraData) < types.PoSExtraVanity {
		genesis.ExtraData = append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.PoSExtraVanity)...)
	}
	genesis.ExtraData = genesis.ExtraData[:types.PoSExtraVanity]

	ist := &types.PoSExtra{
		Validators:    addrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode tendermint extra")
	}
	genesis.ExtraData = append(genesis.ExtraData, istPayload...)
}

func makeHeader(parent *types.Block, config *tendermint.Config) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasLimit:   core.CalcGasLimit(parent, 8000000, 8000000),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(parent.Time(), new(big.Int).SetUint64(config.BlockPeriod)),
		Difficulty: defaultDifficulty,
	}
	return header
}

func makeBlock(chain *core.BlockChain, engine *Backend, parent *types.Block) (*types.Block, error) {
	block, err := makeBlockWithoutSeal(chain, engine, parent)
	if err != nil {
		return nil, err
	}

	resultCh := make(chan *types.Block)
	engine.Seal(chain, block, resultCh, nil)

	return <-resultCh, nil
}

func makeBlockWithoutSeal(chain *core.BlockChain, engine *Backend, parent *types.Block) (*types.Block, error) {
	header := makeHeader(parent, engine.config)
	engine.Prepare(chain, header)
	state, err := chain.StateAt(parent.Root())
	block, _ := engine.Finalize(chain, header, state, nil, nil, nil)

	// Write state changes to db
	root, err := state.Commit(chain.Config().IsEIP158(block.Header().Number))
	if err != nil {
		return nil, fmt.Errorf("state write error: %v", err)
	}
	if err := state.Database().TrieDB().Commit(root, false); err != nil {
		return nil, fmt.Errorf("trie write error: %v", err)
	}

	return block, nil
}