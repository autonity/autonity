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
	"math/big"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/consensus/tendermint"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/core/vm"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/params"
)

const (
	// fetcherID is the ID indicates the block is from PoS engine
	fetcherID = "tendermint"
)

// New creates an Ethereum Backend for PoS core engine.
func New(config *tendermint.Config, privateKey *ecdsa.PrivateKey, db ethdb.Database, chainConfig *params.ChainConfig, vmConfig *vm.Config) *Backend {
	if chainConfig.Tendermint.Epoch != 0 {
		config.Epoch = chainConfig.Tendermint.Epoch
	}

	if chainConfig.Tendermint.Bytecode != "" && chainConfig.Tendermint.ABI != "" {
		config.Bytecode = chainConfig.Tendermint.Bytecode
		config.ABI = chainConfig.Tendermint.ABI
		log.Info("Default Validator smart contract set")
	} else {
		log.Info("User specified Validator smart contract set")
	}

	if chainConfig.Tendermint.RequestTimeout != 0 {
		config.RequestTimeout = chainConfig.Tendermint.RequestTimeout
	}
	if chainConfig.Tendermint.BlockPeriod != 0 {
		config.BlockPeriod = chainConfig.Tendermint.BlockPeriod
	}

	config.SetProposerPolicy(tendermint.ProposerPolicy(chainConfig.Tendermint.ProposerPolicy))

	recents, _ := lru.NewARC(inmemorySnapshots)
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)
	logger := log.New()
	backend := &Backend{
		config:         config,
		eventMux:       event.NewTypeMuxSilent(logger),
		privateKey:     privateKey,
		address:        crypto.PubkeyToAddress(privateKey.PublicKey),
		logger:         log.New(),
		db:             db,
		recents:        recents,
		coreStarted:    false,
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		vmConfig:       vmConfig,
		resend:         make(chan messageToPeers, 1024),
	}
	backend.core = tendermintCore.New(backend, backend.config)

	go backend.ReSend(10)

	return backend
}

// ----------------------------------------------------------------------------

type Backend struct {
	config       *tendermint.Config
	eventMux     *event.TypeMuxSilent
	privateKey   *ecdsa.PrivateKey
	address      common.Address
	core         tendermintCore.Engine
	logger       log.Logger
	db           ethdb.Database
	blockchain   *core.BlockChain
	currentBlock func() *types.Block
	hasBadBlock  func(hash common.Hash) bool

	// the channels for tendermint engine notifications
	commitCh          chan<- *types.Block
	proposedBlockHash common.Hash
	coreStarted       bool
	coreMu            sync.RWMutex

	// Snapshots for recent block to speed up reorgs
	recents *lru.ARCCache

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	//TODO: ARCChace is patented by IBM, so probably need to stop using it
	recentMessages *lru.ARCCache // the cache of peer's messages
	knownMessages  *lru.ARCCache // the cache of self messages

	somaContract      common.Address // Ethereum address of the governance contract
	glienickeContract common.Address // Ethereum address of the white list contract
	vmConfig          *vm.Config

	resend chan messageToPeers
}

// Address implements tendermint.Backend.Address
func (sb *Backend) Address() common.Address {
	return sb.address
}

func (sb *Backend) Validators(number uint64) tendermint.ValidatorSet {
	validators, err := sb.retrieveSavedValidators(number, sb.blockchain)
	proposerPolicy := sb.config.GetProposerPolicy()
	if err != nil {
		return validator.NewSet(nil, proposerPolicy)
	}
	return validator.NewSet(validators, proposerPolicy)
}

// Broadcast implements tendermint.Backend.Broadcast
func (sb *Backend) Broadcast(valSet tendermint.ValidatorSet, payload []byte) error {
	// send to others
	sb.Gossip(valSet, payload)
	// send to self
	msg := tendermint.MessageEvent{
		Payload: payload,
	}
	go sb.eventMux.Post(msg)
	return nil
}

const TTL = 10 //seconds

// Broadcast implements tendermint.Backend.Gossip
func (sb *Backend) Gossip(valSet tendermint.ValidatorSet, payload []byte) {
	hash := types.RLPHash(payload)
	sb.knownMessages.Add(hash, true)

	targets := make(map[common.Address]bool)
	for _, val := range valSet.List() {
		if val.Address() != sb.Address() {
			targets[val.Address()] = true
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	if sb.broadcaster != nil && len(targets) > 0 {
		ps, notConnected := sb.broadcaster.FindPeers(targets)
		for addr, p := range ps {
			sb.sendToPeer(ctx, cancel, addr, hash, payload, p)
		}

		if len(notConnected) > 0 {
			sb.logger.Info("gossip. got not connected peers", "peers", notConnected)

			sb.resend <- messageToPeers{
				message{
					hash,
					payload,
				},
				notConnected,
				time.Now(),
			}
		}
	}
}

type messageToPeers struct {
	msg       message
	peers     []common.Address
	startTime time.Time
}

type message struct {
	hash    common.Hash
	payload []byte
}

type peerError struct {
	error
	addr common.Address
}

func (sb *Backend) sendToPeer(ctx context.Context, cancelFunc context.CancelFunc, addr common.Address, hash common.Hash, payload []byte, p consensus.Peer) chan error {
	ms, ok := sb.recentMessages.Get(addr)
	errCh := make(chan error, 1)
	var m *lru.ARCCache
	if ok {
		m, _ = ms.(*lru.ARCCache)
		if _, k := m.Get(hash); k {
			// This peer had this event, skip it
			errCh <- nil
			return errCh
		}
	} else {
		m, _ = lru.NewARC(inmemoryMessages)
	}

	go func(p consensus.Peer, m *lru.ARCCache) {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		var err error

	SenderLoop:
		for {
			select {
			case <-ticker.C:
				sb.logger.Info("inner sender loop", "msg", payload)

				if err := p.Send(tendermintMsg, payload); err != nil {
					err = peerError{errors.New("error while sending tendermintMsg message to the peer: " + err.Error()), addr}
				} else {
					err = nil
					cancelFunc()
					break SenderLoop
				}
			case <-ctx.Done():
				err = peerError{errors.New("error while sending tendermintMsg message to the peer: " + ctx.Err().Error()), addr}
			}
		}

		if err == nil {
			m.Add(hash, true)
			sb.recentMessages.Add(addr, m)
		}

		errCh <- err
	}(p, m)

	return errCh
}

func (sb *Backend) ReSend(numberOfWorkers int) {
	workers := sync.WaitGroup{}
	workers.Add(numberOfWorkers)

	for i := 0; i < numberOfWorkers; i++ {
		go func() {
			sb.workerSendLoop()
			workers.Done()
		}()
	}

	workers.Wait()
}

func (sb *Backend) workerSendLoop() {
	for msgToPeers := range sb.resend {
		if int(time.Now().Sub(msgToPeers.startTime).Seconds()) > TTL {
			continue
		}

		sb.logger.Info("worker loop. resending", "msg", msgToPeers)

		m := make(map[common.Address]bool)
		for _, p := range msgToPeers.peers {
			m[p] = true
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		ps, notConnected := sb.broadcaster.FindPeers(m)
		sb.logger.Info("worker loop. got not connected peers", "peers", notConnected)

		var errChs []chan error
		if sb.broadcaster != nil && len(ps) > 0 {
			for addr, p := range ps {
				errCh := sb.sendToPeer(ctx, cancel, addr, msgToPeers.msg.hash, msgToPeers.msg.payload, p)
				errChs = append(errChs, errCh)
			}
		}

		wg := sync.WaitGroup{}
		wg.Add(len(errChs))
		notConnectedCh := make(chan common.Address, len(errChs))

		for _, errCh := range errChs {
			go func(errCh chan error) {
				err := <-errCh
				if err != nil {
					pe, ok := err.(peerError)
					if ok {
						notConnectedCh <- pe.addr

						sb.logger.Error(pe.Error(), "peer", pe.addr)
					}
				}

				wg.Done()
			}(errCh)
		}
		wg.Wait()

		for addr := range notConnectedCh {
			notConnected = append(notConnected, addr)
		}

		if int(time.Now().Sub(msgToPeers.startTime).Seconds()) > TTL {
			continue
		}
		if len(notConnected) > 0 {
			sb.logger.Info("worker loop. got not connected and error peers", "peers", notConnected)

			sb.resend <- messageToPeers{
				message{
					msgToPeers.msg.hash,
					msgToPeers.msg.payload,
				},
				notConnected,
				msgToPeers.startTime,
			}
		}
	}
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
	if sb.proposedBlockHash == block.Hash() && sb.commitCh != nil {
		// feed block hash to Seal() and wait the Seal() result
		sb.commitCh <- block
		return nil
	}

	if sb.broadcaster != nil {
		sb.broadcaster.Enqueue(fetcherID, block)
	}
	return nil
}

// EventMux implements tendermint.Backend.EventMux
func (sb *Backend) EventMux() *event.TypeMuxSilent {
	return sb.eventMux
}

// Verify implements tendermint.Backend.Verify
func (sb *Backend) Verify(proposal types.Block) (time.Duration, error) {
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

	// check block body
	txnHash := types.DeriveSha(block.Transactions())
	uncleHash := types.CalcUncleHash(block.Uncles())
	if txnHash != block.Header().TxHash {
		return 0, errMismatchTxhashes
	}
	if uncleHash != nilUncleHash {
		return 0, errInvalidUncleHash
	}

	// verify the header of proposed block
	err := sb.VerifyHeader(sb.blockchain, block.Header(), false)
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || err == types.ErrEmptyCommittedSeals {
		// the current blockchain state is synchronized with PoS's state
		// and we know that the proposed block was mined by a valid validator
		header := block.Header()
		//We need at this point to process all the transactions in the block
		//in order to extract the list of the next validators and validate the extradata field
		var validators []common.Address
		if header.Number.Uint64() > 1 {

			state, _ := sb.blockchain.State()
			state = state.Copy() // copy the state, we don't want to save modifications
			gp := new(core.GasPool).AddGas(block.GasLimit())
			usedGas := new(uint64)
			// blockchain.Processor().Process() would have been a better choice but it calls back Finalize()
			for i, tx := range block.Transactions() {
				state.Prepare(tx.Hash(), block.Hash(), i)
				// Might be vulnerable to DoS Attack depending on gaslimit
				// Todo : Double check
				_, _, err = core.ApplyTransaction(sb.blockchain.Config(), sb.blockchain, nil,
					gp, state, header, tx, usedGas, *sb.vmConfig)

				if err != nil {
					return 0, err
				}
			}

			validators, err = sb.contractGetValidators(sb.blockchain, header, state)
			if err != nil {
				return 0, err
			}
		} else {
			validators, err = sb.retrieveSavedValidators(1, sb.blockchain) //genesis block and block #1 have the same validators
			if err != nil {
				return 0, err
			}
		}
		tendermintExtra, _ := types.ExtractPoSExtra(header)

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
		return time.Unix(block.Header().Time.Int64(), 0).Sub(now()), consensus.ErrFutureBlock
	}
	return 0, err
}

// Sign implements tendermint.Backend.Sign
func (sb *Backend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, sb.privateKey)
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

// HasPropsal implements tendermint.Backend.HashBlock
func (sb *Backend) HasPropsal(hash common.Hash, number *big.Int) bool {
	return sb.blockchain.GetHeader(hash, number.Uint64()) != nil
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

	enodes, err := sb.blockchain.GetWhitelist(sb.blockchain.CurrentBlock(), db)
	if err != nil {
		sb.logger.Error("Failed to get block white list", "err", err)
		return nil
	}

	return enodes.StrList
}
