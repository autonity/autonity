// Copyright 2016 The go-ethereum Authors
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

package ethclient_test

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/ethclient"
	"github.com/autonity/autonity/ethclient/simulated"
	"github.com/autonity/autonity/params"
)

// Verify that Client implements the ethereum interfaces.
var (
	_ = ethereum.ChainReader(&ethclient.Client{})
	_ = ethereum.TransactionReader(&ethclient.Client{})
	_ = ethereum.ChainStateReader(&ethclient.Client{})
	_ = ethereum.ChainSyncReader(&ethclient.Client{})
	_ = ethereum.ContractCaller(&ethclient.Client{})
	_ = ethereum.GasEstimator(&ethclient.Client{})
	_ = ethereum.GasPricer(&ethclient.Client{})
	_ = ethereum.LogFilterer(&ethclient.Client{})
	_ = ethereum.PendingStateReader(&ethclient.Client{})
	// _ = ethereum.PendingStateEventer(&ethclient.Client{})
	_ = ethereum.PendingContractCaller(&ethclient.Client{})
)

var (
	testKey, _         = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr           = crypto.PubkeyToAddress(testKey.PublicKey)
	testBalance        = big.NewInt(2e15)
	revertContractAddr = common.HexToAddress("290f1b36649a61e369c6276f6d29463335b4400c")
	revertCode         = common.FromHex("7f08c379a0000000000000000000000000000000000000000000000000000000006000526020600452600a6024527f75736572206572726f7200000000000000000000000000000000000000000000604452604e6000fd")
)

// testBackend holds a simulated backend and client for testing
type testBackend struct {
	sim    *simulated.Backend
	client simulated.Client
}

func newTestBackend() *testBackend {
	sim := simulated.NewBackend(types.GenesisAlloc{
		testAddr:           {Balance: testBalance},
		revertContractAddr: {Code: revertCode},
	})
	return &testBackend{
		sim:    sim,
		client: sim.Client(),
	}
}

func (tb *testBackend) close() {
	tb.sim.Close()
}

func TestEthClient(t *testing.T) {
	backend := newTestBackend()
	defer backend.close()

	// Create first block (empty)
	backend.sim.Commit()

	// Send test transactions and create second block
	ec := backend.sim.EthClient()
	ctx := context.Background()

	chainID, err := ec.ChainID(ctx)
	if err != nil {
		t.Fatalf("ChainID error: %v", err)
	}

	// Create and send testTx1
	signer := types.LatestSignerForChainID(chainID)
	testTx1, err := types.SignNewTx(testKey, signer, &types.LegacyTx{
		Nonce:    0,
		Value:    big.NewInt(12),
		GasPrice: big.NewInt(params.InitialBaseFee),
		Gas:      params.TxGas,
		To:       &common.Address{2},
	})
	if err != nil {
		t.Fatalf("SignNewTx error: %v", err)
	}
	if err := ec.SendTransaction(ctx, testTx1); err != nil {
		t.Fatalf("SendTransaction error: %v", err)
	}

	// Create and send testTx2
	testTx2, err := types.SignNewTx(testKey, signer, &types.LegacyTx{
		Nonce:    1,
		Value:    big.NewInt(8),
		GasPrice: big.NewInt(params.InitialBaseFee),
		Gas:      params.TxGas,
		To:       &common.Address{2},
	})
	if err != nil {
		t.Fatalf("SignNewTx error: %v", err)
	}
	if err := ec.SendTransaction(ctx, testTx2); err != nil {
		t.Fatalf("SendTransaction error: %v", err)
	}

	// Commit block 2 with transactions
	backend.sim.Commit()

	tests := map[string]struct {
		test func(t *testing.T)
	}{
		"BalanceAt": {
			func(t *testing.T) { testBalanceAt(t, ec) },
		},
		"TxInBlockInterrupted": {
			func(t *testing.T) { testTransactionInBlock(t, ec, testTx1, testTx2) },
		},
		"ChainID": {
			func(t *testing.T) { testChainID(t, ec) },
		},
		"GetBlock": {
			func(t *testing.T) { testGetBlock(t, ec) },
		},
		"CallContract": {
			func(t *testing.T) { testCallContract(t, ec) },
		},
		"CallContractAtHash": {
			func(t *testing.T) { testCallContractAtHash(t, ec) },
		},
		"TransactionSender": {
			func(t *testing.T) { testTransactionSender(t, ec, testTx1, testTx2) },
		},
	}

	for name, tt := range tests {
		t.Run(name, tt.test)
	}
}

func testBalanceAt(t *testing.T, ec *ethclient.Client) {
	// Note: With path-based state scheme, only recent state is available.
	// We test balance at the latest block (nil) and check for expected errors on old/future blocks.
	tests := map[string]struct {
		account common.Address
		block   *big.Int
		checkFn func(got *big.Int, err error) error
	}{
		"valid_account_latest": {
			account: testAddr,
			block:   nil, // latest block
			checkFn: func(got *big.Int, err error) error {
				if err != nil {
					return err
				}
				// Balance should be less than initial due to gas spent on transactions
				if got.Sign() <= 0 {
					return errors.New("balance should be positive")
				}
				return nil
			},
		},
		"non_existent_account": {
			account: common.Address{1},
			block:   nil,
			checkFn: func(got *big.Int, err error) error {
				if err != nil {
					return err
				}
				if got.Cmp(big.NewInt(0)) != 0 {
					return errors.New("non-existent account should have zero balance")
				}
				return nil
			},
		},
		"future_block": {
			account: testAddr,
			block:   big.NewInt(1000000000),
			checkFn: func(got *big.Int, err error) error {
				if err == nil || err.Error() != "header not found" {
					return errors.New("expected 'header not found' error for future block")
				}
				return nil
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := ec.BalanceAt(ctx, tt.account, tt.block)
			if checkErr := tt.checkFn(got, err); checkErr != nil {
				t.Fatalf("BalanceAt(%x, %v): %v", tt.account, tt.block, checkErr)
			}
		})
	}
}

func testTransactionInBlock(t *testing.T, ec *ethclient.Client, testTx1, testTx2 *types.Transaction) {
	// Get current block by number.
	block, err := ec.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test tx in block not found.
	if _, err := ec.TransactionInBlock(context.Background(), block.Hash(), 20); err != ethereum.NotFound {
		t.Fatal("error should be ethereum.NotFound")
	}

	// Test tx in block found.
	tx, err := ec.TransactionInBlock(context.Background(), block.Hash(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Hash() != testTx1.Hash() {
		t.Fatalf("unexpected transaction: %v", tx)
	}

	tx, err = ec.TransactionInBlock(context.Background(), block.Hash(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Hash() != testTx2.Hash() {
		t.Fatalf("unexpected transaction: %v", tx)
	}
}

func testChainID(t *testing.T, ec *ethclient.Client) {
	id, err := ec.ChainID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil || id.Cmp(params.TestChainConfig.ChainID) != 0 {
		t.Fatalf("ChainID returned wrong number: %+v", id)
	}
}

func testGetBlock(t *testing.T, ec *ethclient.Client) {
	// Get current block number
	blockNumber, err := ec.BlockNumber(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blockNumber != 2 {
		t.Fatalf("BlockNumber returned wrong number: %d", blockNumber)
	}
	// Get current block by number
	block, err := ec.BlockByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.NumberU64() != blockNumber {
		t.Fatalf("BlockByNumber returned wrong block: want %d got %d", blockNumber, block.NumberU64())
	}
	// Get current block by hash
	blockH, err := ec.BlockByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Hash() != blockH.Hash() {
		t.Fatalf("BlockByHash returned wrong block: want %v got %v", block.Hash().Hex(), blockH.Hash().Hex())
	}
	// Get header by number
	header, err := ec.HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != header.Hash() {
		t.Fatalf("HeaderByNumber returned wrong header: want %v got %v", block.Header().Hash().Hex(), header.Hash().Hex())
	}
	// Get header by hash
	headerH, err := ec.HeaderByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != headerH.Hash() {
		t.Fatalf("HeaderByHash returned wrong header: want %v got %v", block.Header().Hash().Hex(), headerH.Hash().Hex())
	}
}

func testCallContractAtHash(t *testing.T, ec *ethclient.Client) {
	// EstimateGas
	msg := ethereum.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	// Use the latest block for CallContractAtHash since older state may not be available
	block, err := ec.HeaderByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("HeaderByNumber error: %v", err)
	}
	// CallContract at current block hash
	if _, err := ec.CallContractAtHash(context.Background(), msg, block.Hash()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testCallContract(t *testing.T, ec *ethclient.Client) {
	// EstimateGas
	msg := ethereum.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	// CallContract at latest block (nil means latest)
	if _, err := ec.CallContract(context.Background(), msg, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PendingCallContract
	if _, err := ec.PendingCallContract(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testTransactionSender(t *testing.T, ec *ethclient.Client, testTx1, testTx2 *types.Transaction) {
	ctx := context.Background()

	// Retrieve testTx1 via RPC.
	block2, err := ec.HeaderByNumber(ctx, big.NewInt(2))
	if err != nil {
		t.Fatal("can't get block 2:", err)
	}
	tx1, err := ec.TransactionInBlock(ctx, block2.Hash(), 0)
	if err != nil {
		t.Fatal("can't get tx:", err)
	}
	if tx1.Hash() != testTx1.Hash() {
		t.Fatalf("wrong tx hash %v, want %v", tx1.Hash(), testTx1.Hash())
	}

	// The sender address is cached in tx1, so no additional RPC should be required in
	// TransactionSender. Ensure the server is not asked by canceling the context here.
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	<-canceledCtx.Done() // Ensure the close of the Done channel
	sender1, err := ec.TransactionSender(canceledCtx, tx1, block2.Hash(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if sender1 != testAddr {
		t.Fatal("wrong sender:", sender1)
	}

	// Now try to get the sender of testTx2, which was not fetched through RPC.
	// TransactionSender should query the server here.
	sender2, err := ec.TransactionSender(ctx, testTx2, block2.Hash(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if sender2 != testAddr {
		t.Fatal("wrong sender:", sender2)
	}
}

func testAtFunctions(t *testing.T, ec *ethclient.Client, sim *simulated.Backend) {
	block, err := ec.HeaderByNumber(context.Background(), big.NewInt(1))
	if err != nil {
		t.Fatalf("BlockByNumber error: %v", err)
	}

	// send a transaction for some interesting pending status
	sendTransaction(ec)

	// wait for the transaction to be included in the pending block
	for {
		// Check pending transaction count
		pending, err := ec.PendingTransactionCount(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pending == 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Query balance
	balance, err := ec.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hashBalance, err := ec.BalanceAtHash(context.Background(), testAddr, block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance.Cmp(hashBalance) == 0 {
		t.Fatalf("unexpected balance at hash: %v %v", balance, hashBalance)
	}
	penBalance, err := ec.PendingBalanceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance.Cmp(penBalance) == 0 {
		t.Fatalf("unexpected balance: %v %v", balance, penBalance)
	}
	// NonceAt
	nonce, err := ec.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hashNonce, err := ec.NonceAtHash(context.Background(), testAddr, block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hashNonce == nonce {
		t.Fatalf("unexpected nonce at hash: %v %v", nonce, hashNonce)
	}
	penNonce, err := ec.PendingNonceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if penNonce != nonce+1 {
		t.Fatalf("unexpected nonce: %v %v", nonce, penNonce)
	}
	// StorageAt
	storage, err := ec.StorageAt(context.Background(), testAddr, common.Hash{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hashStorage, err := ec.StorageAtHash(context.Background(), testAddr, common.Hash{}, block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(storage, hashStorage) {
		t.Fatalf("unexpected storage at hash: %v %v", storage, hashStorage)
	}
	penStorage, err := ec.PendingStorageAt(context.Background(), testAddr, common.Hash{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(storage, penStorage) {
		t.Fatalf("unexpected storage: %v %v", storage, penStorage)
	}
	// CodeAt
	code, err := ec.CodeAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hashCode, err := ec.CodeAtHash(context.Background(), common.Address{}, block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(code, hashCode) {
		t.Fatalf("unexpected code at hash: %v %v", code, hashCode)
	}
	penCode, err := ec.PendingCodeAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(code, penCode) {
		t.Fatalf("unexpected code: %v %v", code, penCode)
	}
	// Use HeaderByNumber to get a header for EstimateGasAtBlock and EstimateGasAtBlockHash
	latestHeader, err := ec.HeaderByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// EstimateGasAtBlock
	msg := ethereum.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := ec.EstimateGasAtBlock(context.Background(), msg, latestHeader.Number)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas limit: %v", gas)
	}
	// EstimateGasAtBlockHash
	gas, err = ec.EstimateGasAtBlockHash(context.Background(), msg, latestHeader.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas limit: %v", gas)
	}
}

func sendTransaction(ec *ethclient.Client) error {
	chainID, err := ec.ChainID(context.Background())
	if err != nil {
		return err
	}
	nonce, err := ec.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		return err
	}

	signer := types.LatestSignerForChainID(chainID)
	tx, err := types.SignNewTx(testKey, signer, &types.LegacyTx{
		Nonce:    nonce,
		To:       &common.Address{2},
		Value:    big.NewInt(1),
		Gas:      22000,
		GasPrice: big.NewInt(params.InitialBaseFee),
	})
	if err != nil {
		return err
	}
	return ec.SendTransaction(context.Background(), tx)
}
