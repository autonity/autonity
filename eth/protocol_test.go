// Copyright 2014 The go-ethereum Authors
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

package eth

import (
	"fmt"
	"github.com/clearmatics/autonity/consensus/ethash"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
	"sync"
	"testing"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/eth/downloader"
	"github.com/clearmatics/autonity/p2p"
	"github.com/clearmatics/autonity/rlp"
)

func init() {
	// log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false))))
}

var testAccount, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

// Tests that handshake failures are detected and reported correctly.
func TestStatusMsgErrors62(t *testing.T) { testStatusMsgErrors(t, 62) }
func TestStatusMsgErrors63(t *testing.T) { testStatusMsgErrors(t, 63) }

func testStatusMsgErrors(t *testing.T, protocol int) {
	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, 0, nil, nil, nil)
	var (
		genesis = pm.blockchain.Genesis()
		head    = pm.blockchain.CurrentHeader()
		td      = pm.blockchain.GetTd(head.Hash(), head.Number.Uint64())
	)
	defer pm.Stop()

	tests := []struct {
		code      uint64
		data      interface{}
		wantError error
	}{
		{
			code: TxMsg, data: []interface{}{},
			wantError: errResp(ErrNoStatusMsg, "first msg has code 2 (!= 0)"),
		},
		{
			code: StatusMsg, data: statusData{10, DefaultConfig.NetworkId, td, head.Hash(), genesis.Hash()},
			wantError: errResp(ErrProtocolVersionMismatch, "10 (!= %d)", protocol),
		},
		{
			code: StatusMsg, data: statusData{uint32(protocol), 999, td, head.Hash(), genesis.Hash()},
			wantError: errResp(ErrNetworkIdMismatch, "999 (!= %d)", DefaultConfig.NetworkId),
		},
		{
			code: StatusMsg, data: statusData{uint32(protocol), DefaultConfig.NetworkId, td, head.Hash(), common.Hash{3}},
			wantError: errResp(ErrGenesisBlockMismatch, "0300000000000000 (!= %x)", genesis.Hash().Bytes()[:8]),
		},
	}

	for i, test := range tests {
		p2pPeer := newTestP2PPeer("peer")
		p, errc := newTestPeer(p2pPeer, protocol, pm, false)
		// The send call might hang until reset because
		// the protocol might not read the payload.
		go p2p.Send(p.app, test.code, test.data)

		select {
		case err := <-errc:
			if err == nil {
				t.Errorf("test %d: protocol returned nil error, want %q", i, test.wantError)
			} else if err.Error() != test.wantError.Error() {
				t.Errorf("test %d: wrong error: got %q, want %q", i, err, test.wantError)
			}
		case <-time.After(2 * time.Second):
			t.Errorf("protocol did not shut down within 2 seconds")
		}
		p.close()
	}
}

// This test checks that received transactions are added to the local pool.
func TestRecvTransactions62(t *testing.T) { testRecvTransactions(t, 62) }
func TestRecvTransactions63(t *testing.T) { testRecvTransactions(t, 63) }

func testRecvTransactions(t *testing.T, protocol int) {
	txAdded := make(chan []*types.Transaction)

	p2pPeer := newTestP2PPeer("peer")

	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, 0, nil, txAdded, []string{p2pPeer.Info().Enode})
	pm.acceptTxs = 1 // mark synced to accept transactions
	p, _ := newTestPeer(p2pPeer, protocol, pm, true)
	defer pm.Stop()
	defer p.close()

	tx := newTestTransaction(testAccount, 0, 0)
	if err := p2p.Send(p.app, TxMsg, []interface{}{tx}); err != nil {
		t.Fatalf("send error: %v", err)
	}
	select {
	case added := <-txAdded:
		if len(added) != 1 {
			t.Errorf("wrong number of added transactions: got %d, want 1", len(added))
		} else if added[0].Hash() != tx.Hash() {
			t.Errorf("added wrong tx hash: got %v, want %v", added[0].Hash(), tx.Hash())
		}
	case <-time.After(2 * time.Second):
		t.Errorf("no NewTxsEvent received within 2 seconds")
	}
}

// This test checks that pending transactions are sent.
func TestSendTransactions62(t *testing.T) { testSendTransactions(t, 62) }
func TestSendTransactions63(t *testing.T) { testSendTransactions(t, 63) }

func testSendTransactions(t *testing.T, protocol int) {
	var (
		evmux   = new(event.TypeMux)
		pow     = ethash.NewFaker()
		db      = rawdb.NewMemoryDatabase()
		config  = &params.ChainConfig{}
		gspec   = &core.Genesis{Config: config}
	)

	var p2pPeers []*p2p.Peer
	var enodes []string
	for i := 0; i < 3; i++ {
		p2pPeers = append(p2pPeers, newTestP2PPeer(fmt.Sprintf("peer #%d", i)))
		enodes = append(enodes, p2pPeers[i].Info().Enode)
	}

	gspec.Config.EnodeWhitelist = enodes
	gspec.MustCommit(db)

	blockchain, err := core.NewBlockChain(db, nil, config, pow, vm.Config{}, nil)
	if err != nil {
		t.Fatalf("failed to create new blockchain: %v", err)
	}

	const txCount = 100
	txAdded := make(chan []*types.Transaction, txCount)
	pm, err := NewProtocolManager(config, nil, downloader.FullSync, DefaultConfig.NetworkId, evmux, &testTxPool{added: txAdded}, pow, blockchain, db, 1, nil, DefaultConfig.OpenNetwork)
	if err != nil {
		t.Fatalf("failed to start test protocol manager: %v", err)
	}
	pm.Start(1000)
	defer pm.Stop()

	var peers []*testPeer
	for _, p2pPeer := range p2pPeers {
		p, _ := newTestPeer(p2pPeer, protocol, pm, true)
		peers = append(peers, p)
	}

	// Fill the pool with big transactions.
	const txsize = txsyncPackSize / 10
	alltxs := make([]*types.Transaction, txCount)
	for nonce := range alltxs {
		alltxs[nonce] = newTestTransaction(testAccount, uint64(nonce), txsize)
	}
	pm.BroadcastTxs(alltxs)

	// Connect several peers. They should all receive the pending transactions.
	var wg sync.WaitGroup

	type message struct {
		msg p2p.Msg
		err error
	}

	checktxs := func(p *testPeer) {
		defer wg.Done()
		defer p.close()

		seen := make(map[common.Hash]bool)
		for _, tx := range alltxs {
			seen[tx.Hash()] = false
		}
		msgCh := make(chan message, 1)

		for n := 0; n < len(alltxs) && !t.Failed(); {
			var txs []*types.Transaction

			go func() {
				msg, err := p.app.ReadMsg()
				msgCh <- message{msg, err}
			}()

			var readMsg message
			select {
			case readMsg = <-msgCh:
				break
			case <-time.After(60 * time.Second):
				t.Fatalf("timeout for peer %v: read error: %v. Tx in pool %d", p.Peer, err, len(txAdded))
			}

			msg := readMsg.msg

			if readMsg.err != nil {
				t.Fatalf("%v: read error: %v, Msg: %v", p.Peer, err, msg.String())
			} else if msg.Code == 7 {
				log.Debug("genesis block message")
				continue
			} else if msg.Code != TxMsg {
				t.Errorf("%v: got code %d, want TxMsg: %v", p.Peer, msg.Code, msg)
			}
			if err := msg.Decode(&txs); err != nil {
				t.Errorf("%v: %v", p.Peer, err)
			}
			for _, tx := range txs {
				hash := tx.Hash()
				seentx, want := seen[hash]
				if seentx {
					t.Errorf("%v: got tx more than once: %x", p.Peer, hash)
				}
				if !want {
					t.Errorf("%v: got unexpected tx: %x", p.Peer, hash)
				}
				seen[hash] = true
				n++
			}
		}
	}
	for _, p := range peers {
		wg.Add(1)
		go checktxs(p)
	}
	wg.Wait()
}

// Tests that the custom union field encoder and decoder works correctly.
func TestGetBlockHeadersDataEncodeDecode(t *testing.T) {
	// Create a "random" hash for testing
	var hash common.Hash
	for i := range hash {
		hash[i] = byte(i)
	}
	// Assemble some table driven tests
	tests := []struct {
		packet *getBlockHeadersData
		fail   bool
	}{
		// Providing the origin as either a hash or a number should both work
		{fail: false, packet: &getBlockHeadersData{Origin: hashOrNumber{Number: 314}}},
		{fail: false, packet: &getBlockHeadersData{Origin: hashOrNumber{Hash: hash}}},

		// Providing arbitrary query field should also work
		{fail: false, packet: &getBlockHeadersData{Origin: hashOrNumber{Number: 314}, Amount: 314, Skip: 1, Reverse: true}},
		{fail: false, packet: &getBlockHeadersData{Origin: hashOrNumber{Hash: hash}, Amount: 314, Skip: 1, Reverse: true}},

		// Providing both the origin hash and origin number must fail
		{fail: true, packet: &getBlockHeadersData{Origin: hashOrNumber{Hash: hash, Number: 314}}},
	}
	// Iterate over each of the tests and try to encode and then decode
	for i, tt := range tests {
		bytes, err := rlp.EncodeToBytes(tt.packet)
		if err != nil && !tt.fail {
			t.Fatalf("test %d: failed to encode packet: %v", i, err)
		} else if err == nil && tt.fail {
			t.Fatalf("test %d: encode should have failed", i)
		}
		if !tt.fail {
			packet := new(getBlockHeadersData)
			if err := rlp.DecodeBytes(bytes, packet); err != nil {
				t.Fatalf("test %d: failed to decode packet: %v", i, err)
			}
			if packet.Origin.Hash != tt.packet.Origin.Hash || packet.Origin.Number != tt.packet.Origin.Number || packet.Amount != tt.packet.Amount ||
				packet.Skip != tt.packet.Skip || packet.Reverse != tt.packet.Reverse {
				t.Fatalf("test %d: encode decode mismatch: have %+v, want %+v", i, packet, tt.packet)
			}
		}
	}
}
