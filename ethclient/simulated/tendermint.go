// Copyright 2024 The go-ethereum Authors
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

package simulated

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/misc"
	"github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/txpool"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/eth"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
	"github.com/holiman/uint256"
)

var defaultDifficulty = big.NewInt(1)

// simulatorValidator holds signing keys for block production
type simulatorValidator struct {
	nodeKey      *ecdsa.PrivateKey
	consensusKey blst.SecretKey
	address      common.Address
}

// newSimulatorValidator generates a new validator with ECDSA and BLS keys
func newSimulatorValidator() (*simulatorValidator, error) {
	nodeKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate node key: %w", err)
	}

	consensusKey, err := blst.RandKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate consensus key: %w", err)
	}

	return &simulatorValidator{
		nodeKey:      nodeKey,
		consensusKey: consensusKey,
		address:      crypto.PubkeyToAddress(nodeKey.PublicKey),
	}, nil
}

// tendermintSimulator produces Tendermint-compatible blocks
type tendermintSimulator struct {
	validator  *simulatorValidator
	eth        *eth.Ethereum
	forkParent *types.Header // for Fork() support
	timeAdjust time.Duration // for AdjustTime() support
}

// newTendermintSimulator creates a new simulator with the given eth backend
func newTendermintSimulator(validator *simulatorValidator, ethBackend *eth.Ethereum) *tendermintSimulator {
	return &tendermintSimulator{
		validator: validator,
		eth:       ethBackend,
	}
}

// createBlock creates a new block including pending transactions
func (ts *tendermintSimulator) createBlock() (*types.Block, error) {
	chain := ts.eth.BlockChain()

	// 1. Get parent (forkParent or current head)
	parent := ts.forkParent
	if parent == nil {
		parent = chain.CurrentBlock()
	}

	// 2. Get state at parent
	var statedb *state.StateDB
	var err error

	// If forking to an earlier block, the state might not be available
	if ts.forkParent != nil && !chain.HasState(parent.Root) {
		// For fork operations, try to get state from the block
		block := chain.GetBlockByHash(parent.Hash())
		if block != nil {
			statedb, err = chain.StateAt(block.Root())
		}
		if statedb == nil || err != nil {
			return nil, fmt.Errorf("failed to get state at parent root %s: state not available for fork", parent.Root.Hex())
		}
	} else {
		statedb, err = chain.StateAt(parent.Root)
		if err != nil {
			return nil, fmt.Errorf("failed to get state at parent root: %w", err)
		}
	}

	// 3. Get EIP-1559 parameters
	eip1559Params, err := chain.Eip1559ParamsByHeight(parent.Number.Uint64() + 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get eip1559 params: %w", err)
	}

	// 4. Create header with Tendermint fields
	header := ts.makeHeader(parent, eip1559Params)

	// Save the intended timestamp before Prepare (which may override it)
	intendedTime := header.Time

	// 5. Prepare header (sets coinbase, activity proof, etc.)
	engine := ts.eth.Engine().(*backend.Backend)
	err = engine.Prepare(chain, parent, header, statedb)
	if err != nil {
		return nil, fmt.Errorf("engine prepare failed: %w", err)
	}

	// Override coinbase to match our simulator validator (engine.Prepare sets it to engine's address)
	header.Coinbase = ts.validator.address
	// Restore intended timestamp (Prepare may have set it to current time if it was in the past)
	header.Time = intendedTime

	// 6. Include pending transactions from txpool
	txs, receipts, err := ts.fillTransactions(statedb, header)
	if err != nil {
		return nil, fmt.Errorf("failed to fill transactions: %w", err)
	}

	// 7. Finalize and assemble block
	body := &types.Body{Transactions: txs}
	block, _, err := engine.FinalizeAndAssemble(chain, header, statedb, body, &receipts)
	if err != nil {
		return nil, fmt.Errorf("finalize and assemble failed: %w", err)
	}

	// 8. Add proposer seal
	block, err = ts.addProposerSeal(block)
	if err != nil {
		return nil, fmt.Errorf("failed to add proposer seal: %w", err)
	}

	// 9. Add quorum certificate
	block, err = ts.addQuorumCertificate(block)
	if err != nil {
		return nil, fmt.Errorf("failed to add quorum certificate: %w", err)
	}

	// 10. Commit state to database
	root, err := statedb.Commit(header.Number.Uint64(), chain.Config().IsEIP158(header.Number), chain.Config().IsCancun(header.Number))
	if err != nil {
		return nil, fmt.Errorf("state commit error: %w", err)
	}
	if err := statedb.Database().TrieDB().Commit(root, false); err != nil {
		return nil, fmt.Errorf("trie commit error: %w", err)
	}

	// 11. Reset fork parent and time adjust
	ts.forkParent = nil
	ts.timeAdjust = 0

	return block, nil
}

// makeHeader creates a header with proper Tendermint fields
func (ts *tendermintSimulator) makeHeader(parent *types.Header, eip1559Params *types.Eip1559Params) *types.Header {
	var timestamp uint64
	if ts.timeAdjust > 0 {
		// When time adjustment is specified, use it directly
		timestamp = parent.Time + uint64(ts.timeAdjust.Seconds())
	} else {
		// Default: 1 second after parent
		timestamp = parent.Time + 1
	}

	return &types.Header{
		ParentHash: parent.Hash(),
		Number:     new(big.Int).Add(parent.Number, common.Big1),
		GasLimit:   core.CalcGasLimit(parent.GasLimit, 8000000, 1024),
		GasUsed:    0,
		BaseFee:    misc.CalcBaseFee(params.TestChainConfig, parent, eip1559Params),
		Extra:      parent.Extra,
		Time:       timestamp,
		Difficulty: defaultDifficulty,
		MixDigest:  types.BFTDigest,
		Round:      0,
	}
}

// fillTransactions retrieves pending transactions and applies them
func (ts *tendermintSimulator) fillTransactions(statedb *state.StateDB, header *types.Header) (types.Transactions, []*types.Receipt, error) {
	chain := ts.eth.BlockChain()
	pool := ts.eth.TxPool()

	// Get pending transactions
	filter := txpool.PendingFilter{
		MinTip:       uint256.NewInt(0),
		OnlyPlainTxs: true,
	}
	if header.BaseFee != nil {
		filter.BaseFee = uint256.MustFromBig(header.BaseFee)
	}

	pendingTxs := pool.Pending(filter)

	var txs types.Transactions
	var receipts []*types.Receipt
	gasPool := new(core.GasPool).AddGas(header.GasLimit)
	blockCtx := core.NewEVMBlockContext(header, chain, nil)
	evm := vm.NewEVM(blockCtx, statedb, chain.Config(), *ts.eth.BlockChain().GetVMConfig())

	// Apply transactions
	for _, lazyTxs := range pendingTxs {
		for _, lazyTx := range lazyTxs {
			tx := lazyTx.Resolve()
			if tx == nil {
				continue
			}

			// Check if we have enough gas
			if gasPool.Gas() < tx.Gas() {
				continue
			}

			snap := statedb.Snapshot()
			receipt, err := core.ApplyTransaction(evm, gasPool, statedb, header, tx, &header.GasUsed)
			if err != nil {
				statedb.RevertToSnapshot(snap)
				continue
			}
			txs = append(txs, tx)
			receipts = append(receipts, receipt)
		}
	}

	return txs, receipts, nil
}

// addProposerSeal signs the block with the validator's ECDSA key
func (ts *tendermintSimulator) addProposerSeal(block *types.Block) (*types.Block, error) {
	header := block.Header()
	hashData := types.SigHash(header)
	signature, err := crypto.Sign(hashData[:], ts.validator.nodeKey)
	if err != nil {
		return nil, err
	}
	header.ProposerSeal = signature
	return block.WithSeal(header), nil
}

// addQuorumCertificate adds a BLS quorum certificate and saves the precommit
func (ts *tendermintSimulator) addQuorumCertificate(block *types.Block) (*types.Block, error) {
	chain := ts.eth.BlockChain()
	epoch, err := chain.EpochByHeight(block.NumberU64())
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch: %w", err)
	}

	self := &epoch.Committee.Members[0]
	header := block.Header()

	// Create signer function using the validator's consensus key
	signer := func(hash common.Hash) blst.Signature {
		return ts.validator.consensusKey.Sign(hash[:])
	}

	// Create precommit with BLS signature
	precommit := message.NewPrecommit(
		int64(header.Round),
		header.Number.Uint64(),
		header.Hash(),
		signer,
		self,
		epoch.Committee,
	)

	sig, err := precommit.Signature()
	if err != nil {
		return nil, fmt.Errorf("failed to get precommit signature: %w", err)
	}

	header.QuorumCertificate = types.NewAggregateSignature(
		sig.(*blst.BlsSignature),
		precommit.Signers(),
	)

	// Save precommit for activity proofs on subsequent blocks
	engine := ts.eth.Engine().(*backend.Backend)
	engine.MsgStore.Save(precommit)

	return block.WithSeal(header), nil
}

// setForkParent sets the parent for the next block to create a fork
func (ts *tendermintSimulator) setForkParent(parentHash common.Hash) error {
	parent := ts.eth.BlockChain().GetHeaderByHash(parentHash)
	if parent == nil {
		return fmt.Errorf("parent block not found: %s", parentHash.Hex())
	}
	ts.forkParent = parent
	return nil
}

// adjustTime sets the time adjustment for the next block
func (ts *tendermintSimulator) adjustTime(adjustment time.Duration) {
	ts.timeAdjust = adjustment
}

// rollback clears any pending state (fork parent and time adjustment)
func (ts *tendermintSimulator) rollback() {
	ts.forkParent = nil
	ts.timeAdjust = 0
}

// appendSimulatorValidator adds the simulator validator to the genesis config
func appendSimulatorValidator(genesis *core.Genesis, v *simulatorValidator) {
	if genesis.Config == nil {
		genesis.Config = &params.ChainConfig{}
	}
	if genesis.Config.AutonityContractConfig == nil {
		genesis.Config.AutonityContractConfig = &params.AutonityContractGenesis{}
	}
	if genesis.Config.OracleContractConfig == nil {
		genesis.Config.OracleContractConfig = &params.OracleContractGenesis{}
	}

	nodeAddr := v.address
	node := enode.NewV4(&v.nodeKey.PublicKey, nil, 0, 0)

	genesis.Config.AutonityContractConfig.Validators = []*params.Validator{{
		NodeAddress:   &nodeAddr,
		OracleAddress: nodeAddr,
		Treasury:      nodeAddr,
		Enode:         node.URLv4(),
		BondedStake:   new(big.Int).SetUint64(100),
		ConsensusKey:  v.consensusKey.PublicKey().Marshal(),
	}}

	// Give validator initial balance for gas
	if genesis.Alloc == nil {
		genesis.Alloc = make(types.GenesisAlloc)
	}
	if existing, ok := genesis.Alloc[nodeAddr]; ok {
		// Preserve existing balance if set
		if existing.Balance == nil || existing.Balance.Sign() == 0 {
			genesis.Alloc[nodeAddr] = types.Account{
				Balance: new(big.Int).SetUint64(1e18),
			}
		}
	} else {
		genesis.Alloc[nodeAddr] = types.Account{
			Balance: new(big.Int).SetUint64(1e18),
		}
	}
}
