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
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	tendermintConfig "github.com/clearmatics/autonity/consensus/tendermint/config"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
	"github.com/hashicorp/golang-lru"
)

const (
	// fetcherID is the ID indicates the block is from BFT engine
	fetcherID = "tendermint"
)

var (
	// ErrUnauthorizedAddress is returned when given address cannot be found in
	// current validator set.
	ErrUnauthorizedAddress = errors.New("unauthorized address")
	// ErrStoppedEngine is returned if the engine is stopped
	ErrStoppedEngine = errors.New("stopped engine")
)

// New creates an Ethereum Backend for BFT core engine.
func New(config *tendermintConfig.Config, privateKey *ecdsa.PrivateKey, db ethdb.Database, chainConfig *params.ChainConfig, vmConfig *vm.Config) *Backend {
	log.Warn("new backend with public key",
		"backAddress", crypto.PubkeyToAddress(privateKey.PublicKey).String(),
	)

	if chainConfig.Tendermint.Epoch != 0 {
		config.Epoch = chainConfig.Tendermint.Epoch
	}

	if chainConfig.Tendermint.RequestTimeout != 0 {
		config.RequestTimeout = chainConfig.Tendermint.RequestTimeout
	}
	if chainConfig.Tendermint.BlockPeriod != 0 {
		config.BlockPeriod = chainConfig.Tendermint.BlockPeriod
	}

	config.SetProposerPolicy(tendermintConfig.ProposerPolicy(chainConfig.Tendermint.ProposerPolicy))

	recents, _ := lru.NewARC(inmemorySnapshots)
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)
	logger := log.New()
	backend := &Backend{
		config:         config,
		eventMux:       event.NewTypeMuxSilent(logger),
		privateKey:     privateKey,
		address:        crypto.PubkeyToAddress(privateKey.PublicKey),
		logger:         logger,
		db:             db,
		recents:        recents,
		coreStarted:    false,
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		vmConfig:       vmConfig,
	}

	return backend
}

// ----------------------------------------------------------------------------

type Backend struct {
	config       *tendermintConfig.Config
	eventMux     *event.TypeMuxSilent
	privateKey   *ecdsa.PrivateKey
	privateKeyMu sync.RWMutex
	address      common.Address
	logger       log.Logger
	db           ethdb.Database
	blockchain   *core.BlockChain
	currentBlock func() *types.Block
	hasBadBlock  func(hash common.Hash) bool

	// the channels for tendermint engine notifications
	commitCh          chan<- *types.Block
	proposedBlockHash common.Hash
	coreStarted       bool
	stopped           chan struct{}
	coreMu            sync.RWMutex

	// Snapshots for recent block to speed up reorgs
	recents *lru.ARCCache

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	//TODO: ARCChace is patented by IBM, so probably need to stop using it
	recentMessages *lru.ARCCache // the cache of peer's messages
	knownMessages  *lru.ARCCache // the cache of self messages

	autonityContractAddress common.Address // Ethereum address of the white list contract
	contractsMu             sync.RWMutex
	vmConfig                *vm.Config

	resend chan messageToPeers
}

// Address implements tendermint.Backend.Address
func (sb *Backend) Address() common.Address {
	sb.privateKeyMu.RLock()
	defer sb.privateKeyMu.RUnlock()
	return sb.address
}

func (sb *Backend) Validators(number uint64) validator.Set {
	validators, err := sb.retrieveSavedValidators(number, sb.blockchain)
	proposerPolicy := sb.config.GetProposerPolicy()
	if err != nil {
		return validator.NewSet(nil, proposerPolicy)
	}
	return validator.NewSet(validators, proposerPolicy)
}

// Broadcast implements tendermint.Backend.Broadcast
func (sb *Backend) Broadcast(ctx context.Context, valSet validator.Set, payload []byte) error {
	// send to others
	sb.Gossip(ctx, valSet, payload)
	// send to self
	msg := events.MessageEvent{
		Payload: payload,
	}
	sb.postEvent(msg)
	return nil
}

func (sb *Backend) postEvent(event interface{}) {
	go sb.Post(event)
}

const TTL = 10           //seconds
const retryInterval = 50 //milliseconds

// Broadcast implements tendermint.Backend.Gossip
func (sb *Backend) Gossip(ctx context.Context, valSet validator.Set, payload []byte) {
	hash := types.RLPHash(payload)
	sb.knownMessages.Add(hash, true)

	var targets []common.Address
	for _, val := range valSet.List() {
		if val.Address() != sb.Address() {
			targets = append(targets, val.Address())
		}
	}

	if sb.broadcaster != nil && len(targets) > 0 {
		sb.trySend(ctx, messageToPeers{
			message{
				hash,
				payload,
			},
			targets,
			time.Now(),
			time.Now(),
		})
	}
}

type messageToPeers struct {
	msg       message
	peers     []common.Address
	startTime time.Time
	lastTry   time.Time
}

func (m messageToPeers) String() string {
	msg := fmt.Sprintf("msg hash %s   length %d", m.msg.hash.String(), len(m.msg.payload))
	t := fmt.Sprintf("time %s", m.startTime.String())

	peers := peersToString(m.peers)

	return fmt.Sprintf("%s %s %s", msg, t, peers)
}

func peersToString(ps []common.Address) string {
	peersStr := fmt.Sprintf("peers %d: ", len(ps))
	for _, p := range ps {
		peersStr = fmt.Sprintf("%s%s ", peersStr, p.Hex())
	}

	return peersStr
}

// Commit implements tendermint.Backend.Commit
func (sb *Backend) Commit(proposal types.Block, seals [][]byte) error {
	// Check if the proposal is a valid block
	block := &proposal

	//if block == nil {
	//	sb.logger.Error("Invalid proposal, %v", proposal)
	//	return errInvalidProposal
	//}

	h := block.Header()
	// Append seals into extra-data
	err := types.WriteCommittedSeals(h, seals)
	if err != nil {
		return err
	}
	// update block's header
	block = block.WithSeal(h)

	sb.logger.Info("Committed", "address", sb.Address(), "hash", proposal.Hash(), "number", proposal.Number().Uint64())
	// - if the proposed and committed blocks are the same, send the proposed hash
	//   to commit channel, which is being watched inside the engine.Seal() function.
	// - otherwise, we try to insert the block.
	// -- if success, the ChainHeadEvent event will be broadcasted, try to build
	//    the next block and the previous Seal() will be stopped.
	// -- otherwise, a error will be returned and a round change event will be fired.
	if sb.proposedBlockHash == block.Hash() && !sb.isResultChanNil() {
		// feed block hash to Seal() and wait the Seal() result
		sb.sendResultChan(block)
		return nil
	}

	if sb.broadcaster != nil {
		sb.broadcaster.Enqueue(fetcherID, block)
	}
	return nil
}

func (sb *Backend) Post(ev interface{}) {
	sb.eventMux.Post(ev)
}

func (sb *Backend) Subscribe(types ...interface{}) *event.TypeMuxSubscription {
	return sb.eventMux.Subscribe(types...)
}

// VerifyProposal implements tendermint.Backend.VerifyProposal
func (sb *Backend) VerifyProposal(proposal types.Block) (time.Duration, error) {
	// Check if the proposal is a valid block
	// TODO: fix always false statement and check for non nil
	// TODO: use interface instead of type
	block := &proposal
	//if block == nil {
	//	sb.logger.Error("Invalid proposal, %v", proposal)
	//	return 0, errInvalidProposal
	//}

	// check bad block
	if sb.HasBadProposal(block.Hash()) {
		return 0, core.ErrBlacklistedHash
	}

	// verify the header of proposed block
	err := sb.VerifyHeader(sb.blockchain, block.Header(), false)
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || err == types.ErrEmptyCommittedSeals {
		var (
			receipts   types.Receipts
			validators []common.Address

			usedGas        = new(uint64)
			gp             = new(core.GasPool).AddGas(block.GasLimit())
			header         = block.Header()
			proposalNumber = header.Number.Uint64()
			parent         = sb.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
		)

		// We need to process all of the transaction to get the latest state to get the latest validators
		state, stateErr := sb.blockchain.StateAt(parent.Root())
		if stateErr != nil {
			return 0, stateErr
		}

		// Validate the body of the proposal
		if err = sb.blockchain.Validator().ValidateBody(block); err != nil {
			return 0, err
		}

		// sb.blockchain.Processor().Process() was not called because it calls back Finalize() and would have modified the proposal
		// Instead only the transactions are applied to the copied state
		for i, tx := range block.Transactions() {
			state.Prepare(tx.Hash(), block.Hash(), i)
			// Might be vulnerable to DoS Attack depending on gaslimit
			// Todo : Double check
			receipt, _, receiptErr := core.ApplyTransaction(sb.blockchain.Config(), sb.blockchain, nil, gp, state, header, tx, usedGas, *sb.vmConfig)
			if receiptErr != nil {
				return 0, receiptErr
			}
			receipts = append(receipts, receipt)
		}

		// Here the order of applying transaction matters
		// We need to ensure that the block transactions applied before the Autonity contract
		if proposalNumber == 1 {
			//Apply the same changes from consensus/tendermint/backend/engine.go:getValidator()349-369
			log.Info("Autonity Contract Deployer in test state", "Address", sb.blockchain.Config().AutonityContractConfig.Deployer)

			_, err = sb.blockchain.GetAutonityContract().DeployAutonityContract(sb.blockchain, header, state)
			if err != nil {
				return 0, err
			}
		} else if proposalNumber > 1 {
			err = sb.blockchain.GetAutonityContract().ApplyPerformRedistribution(block.Transactions(), receipts, block.Header(), state)
			if err != nil {
				return 0, err
			}
		}

		//Validate the state of the proposal
		if err = sb.blockchain.Validator().ValidateState(block, state, receipts, *usedGas); err != nil {
			return 0, err
		}

		if proposalNumber > 1 {
			validators, err = sb.blockchain.GetAutonityContract().ContractGetValidators(sb.blockchain, header, state)
			if err != nil {
				return 0, err
			}
		} else {
			validators, err = sb.retrieveSavedValidators(1, sb.blockchain) //genesis block and block #1 have the same validators
			if err != nil {
				return 0, err
			}
		}

		// Verify the validator set by comparing the validators in extra data and Soma-contract
		tendermintExtra, _ := types.ExtractBFTHeaderExtra(header)

		//Perform the actual comparison
		if len(tendermintExtra.Validators) != len(validators) {
			sb.logger.Error("wrong validator set",
				"extraLen", len(tendermintExtra.Validators),
				"currentLen", len(validators),
				"extra", tendermintExtra.Validators,
				"current", validators,
			)
			return 0, errInconsistentValidatorSet
		}

		for i := range validators {
			if tendermintExtra.Validators[i] != validators[i] {
				sb.logger.Error("wrong validator in the set",
					"index", i,
					"extraValidator", tendermintExtra.Validators[i],
					"currentValidator", validators[i],
					"extra", tendermintExtra.Validators,
					"current", validators,
				)
				return 0, errInconsistentValidatorSet
			}
		}
		// At this stage extradata field is consistent with the validator list returned by Soma-contract

		return 0, nil
	} else if err == consensus.ErrFutureBlock {
		return time.Unix(int64(block.Header().Time), 0).Sub(now()), consensus.ErrFutureBlock
	}
	return 0, err
}

// Sign implements tendermint.Backend.Sign
func (sb *Backend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, sb.GetPrivateKey())
}

// CheckSignature implements tendermint.Backend.CheckSignature
func (sb *Backend) CheckSignature(data []byte, address common.Address, sig []byte) error {
	signer, err := types.GetSignatureAddress(data, sig)
	if err != nil {
		log.Error("Failed to get signer address", "err", err)
		return err
	}
	// Compare derived addresses
	if signer != address {
		return types.ErrInvalidSignature
	}
	return nil
}

// GetProposer implements tendermint.Backend.GetProposer
func (sb *Backend) GetProposer(number uint64) common.Address {
	if h := sb.blockchain.GetHeaderByNumber(number); h != nil {
		a, _ := sb.Author(h)
		return a
	}
	return common.Address{}
}

func (sb *Backend) LastCommittedProposal() (*types.Block, common.Address) {
	block := sb.currentBlock()

	var proposer common.Address
	if block.Number().Cmp(common.Big0) > 0 {
		var err error
		proposer, err = sb.Author(block.Header())
		if err != nil {
			sb.logger.Error("Failed to get block proposer", "err", err)
			return new(types.Block), common.Address{}
		}
	}

	// Return header only block here since we don't need block body
	return block, proposer
}

func (sb *Backend) HasBadProposal(hash common.Hash) bool {
	if sb.hasBadBlock == nil {
		return false
	}
	return sb.hasBadBlock(hash)
}

// Whitelist for the current block
func (sb *Backend) WhiteList() []string {
	db, err := sb.blockchain.State()
	if err != nil {
		sb.logger.Error("Failed to get block white list", "err", err)
		return nil
	}

	enodes, err := sb.blockchain.GetAutonityContract().GetWhitelist(sb.blockchain.CurrentBlock(), db)
	if err != nil {
		sb.logger.Error("Failed to get block white list", "err", err)
		return nil
	}

	return enodes.StrList
}

func (sb *Backend) GetPrivateKey() *ecdsa.PrivateKey {
	sb.privateKeyMu.RLock()
	defer sb.privateKeyMu.RUnlock()

	pk := sb.privateKey.PublicKey
	d := big.NewInt(0).Set(sb.privateKey.D)
	return &ecdsa.PrivateKey{PublicKey: pk, D: d}
}

func (sb *Backend) SetPrivateKey(key *ecdsa.PrivateKey) {
	sb.privateKeyMu.Lock()
	defer sb.privateKeyMu.Unlock()

	sb.privateKey = key
	sb.address = crypto.PubkeyToAddress(key.PublicKey)
}

// Synchronize new connected peer with current height state
func (sb *Backend) SyncPeer(address common.Address, messages []*tendermintCore.Message) {
	if sb.broadcaster == nil {
		return
	}

	sb.logger.Info("Syncing", "peer", address)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	for _, msg := range messages {
		payload, err := msg.Payload()
		if err != nil {
			sb.logger.Error("Sending", "code", msg.GetCode(), "sig", msg.GetSignature(), "err", err)
			continue
		}
		hash := types.RLPHash(payload)
		now := time.Now()
		sb.trySend(ctx, messageToPeers{
			message{hash: hash, payload: payload},
			[]common.Address{address},
			now,
			now,
		})
	}
}

func (sb *Backend) ResetPeerCache(address common.Address) {
	ms, ok := sb.recentMessages.Get(address)
	var m *lru.ARCCache
	if ok {
		m, _ = ms.(*lru.ARCCache)
		m.Purge()
	}
}
