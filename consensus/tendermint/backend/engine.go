package backend

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/misc"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core/constants"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/crypto"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rpc"
	"github.com/autonity/autonity/trie"
)

const (
	inmemorySnapshots = 128 // Number of recent vote snapshots to keep in memory
	inmemoryPeers     = 40
	inmemoryMessages  = 1024
)

// ErrStartedEngine is returned if the engine is already started
var ErrStartedEngine = errors.New("started engine")

var (
	// errInvalidProposal is returned when a prposal is malformed.
	//errInvalidProposal = errors.New("invalid proposal")
	// errUnknownBlock is returned when the list of committee is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")
	// errInvalidCoindbase is returned if the signer is not the coinbase address,
	errInvalidCoinbase = errors.New("invalid coinbase")
	// errInvalidDifficulty is returned if the difficulty of a block is not 1
	errInvalidDifficulty = errors.New("invalid difficulty")
	// errInvalidMixDigest is returned if a block's mix digest is not BFT digest.
	errInvalidMixDigest = errors.New("invalid BFT mix digest")
	// errInvalidNonce is returned if a block's nonce is invalid
	errInvalidNonce = errors.New("invalid nonce")
	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")
	// errInvalidTimestamp is returned if the timestamp of a block is lower than the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")
	// errInvalidRound is returned if the round exceed maximum round number.
	errInvalidRound = errors.New("invalid round")
)
var (
	defaultDifficulty             = big.NewInt(1)
	allowedFutureBlockTimeSeconds = int64(1)
	nilUncleHash                  = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce                    = types.BlockNonce{}
	now                           = time.Now
)

// BFT returns true if the engine is an implementation of BFT consensus algorithm.
func (sb *Backend) BFT() bool {
	return true
}

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (sb *Backend) Author(header *types.Header) (common.Address, error) {
	return types.Ecrecover(header)
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *Backend) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, _ bool) error {
	// Short circuit if the header is known, or its parent or epoch head are not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}

	epochHead, parent, err := chain.EpochHeadAndParentHead(number)
	if err != nil {
		return consensus.ErrUnknownAncestor
	}
	return sb.verifyHeader(chain, header, parent, epochHead)
}

// verifyHeader checks whether a header conforms to the consensus rules. It
// expects the parent header to be provided unless header is the genesis
// header.
func (sb *Backend) verifyHeader(chain consensus.ChainHeaderReader, header, parent, epochHead *types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	if header.Round > constants.MaxRound {
		return errInvalidRound
	}
	// Don't waste time checking blocks from the future

	if header.Time > uint64(now().Unix()+allowedFutureBlockTimeSeconds) {
		return consensus.ErrFutureTimestampBlock
	}

	// Ensure that the coinbase is valid
	if header.Nonce != emptyNonce {
		return errInvalidNonce
	}
	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.BFTDigest {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in BFT
	if header.UncleHash != nilUncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}
	// Verify that the gas limit is <= 2^63-1
	if header.GasLimit > params.MaxGasLimit {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, params.MaxGasLimit)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}
	// Verify London hard fork attributes
	// minbasefee is only checked when processing a proposal
	if err := misc.VerifyEip1559Header(chain.Config(), nil, parent, header); err != nil {
		// Verify the header's EIP-1559 attributes.
		return err
	}

	// If this is the genesis block there is no further verification to be
	// done.
	if header.IsGenesis() {
		return nil
	}
	// We expect the parent to be non nil when header is not the genesis header.
	if parent == nil || epochHead == nil {
		return errUnknownBlock
	}

	if err := sb.verifyHeaderAgainstParent(header, parent); err != nil {
		return err
	}

	if err := sb.verifySigner(header, epochHead); err != nil {
		return err
	}

	return sb.verifyCommittedSeals(header, epochHead)
}

// verifyHeaderAgainstParent verifies that the given header is valid with respect to its parent.
func (sb *Backend) verifyHeaderAgainstParent(header, parent *types.Header) error {
	if parent.Number.Uint64() != header.Number.Uint64()-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	if parent.Time+1 > header.Time { // Todo : fetch block period from contract
		return errInvalidTimestamp
	}

	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *Backend) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{}, 1)
	results := make(chan error, len(headers))
	go func() {
		epochHeaders := make(map[uint64]*types.Header)
		for i, header := range headers {
			if header.Number.Cmp(header.LastEpochBlock) == 0 {
				epochHeaders[header.Number.Uint64()] = header
			}
			var parent *types.Header
			switch {
			case i > 0:
				parent = headers[i-1]
			case i == 0:
				parent = chain.GetHeaderByHash(header.ParentHash)
			}
			// resolve correct epoch head, and proceed the header verification, verifyHeader will return err if epoch header is nil.
			epochHead := chain.GetHeaderByNumber(parent.LastEpochBlock.Uint64())
			if epochHead == nil {
				epochHead = epochHeaders[parent.LastEpochBlock.Uint64()]
			}

			err := sb.verifyHeader(chain, header, parent, epochHead)
			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (sb *Backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

// verifySigner checks that the signer is part of the committee.
func (sb *Backend) verifySigner(header, epochHead *types.Header) error {
	// resolve the authorization key and check against signers
	signer, err := types.Ecrecover(header)
	if err != nil {
		return err
	}

	if header.Coinbase != signer {
		return errInvalidCoinbase
	}

	// Signer should be in the validator set of previous block's extraData.
	if epochHead.CommitteeMember(signer) != nil {
		return nil
	}

	return errUnauthorized
}

// verifyCommittedSeals validates that the committed seals for header come from
// committee members and that the voting power of the committed seals constitutes
// a quorum.
func (sb *Backend) verifyCommittedSeals(header, epochHead *types.Header) error {
	// The length of Committed seals should be larger than 0
	if len(header.CommittedSeals) == 0 {
		return types.ErrEmptyCommittedSeals
	}

	// Setup map to track votes made by committee members
	votes := make(map[common.Address]int, len(epochHead.Committee.Members))

	// Calculate total voting power
	totalVotingPower := epochHead.TotalVotingPower()

	// Total Voting power for this block
	power := new(big.Int)
	// The data that was sined over for this block
	headerSeal := helpers.PrepareCommittedSeal(header.Hash(), int64(header.Round), header.Number)

	// 1. Get committed seals from current header
	for _, signedSeal := range header.CommittedSeals {
		// 2. Get the address from signature
		addr, err := types.GetSignatureAddress(headerSeal, signedSeal)
		if err != nil {
			sb.logger.Error("not a valid address", "err", err)
			return types.ErrInvalidSignature
		}

		member := epochHead.CommitteeMember(addr)
		if member == nil {
			sb.logger.Error(fmt.Sprintf("block had seal from non committee member %q", addr))
			return types.ErrInvalidCommittedSeals
		}

		votes[member.Address]++
		if votes[member.Address] > 1 {
			sb.logger.Error(fmt.Sprintf("committee member %q had multiple seals on block", addr))
			return types.ErrInvalidCommittedSeals
		}
		power.Add(power, member.VotingPower)
	}

	// We need at least a quorum for the block to be considered valid
	if power.Cmp(bft.Quorum(totalVotingPower)) < 0 {
		return types.ErrInvalidCommittedSeals
	}

	return nil
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *Backend) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = sb.Address()
	header.Nonce = emptyNonce
	header.MixDigest = types.BFTDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same difficulty for all blocks
	header.Difficulty = defaultDifficulty

	// set header's timestamp
	// todo: block period from contract
	header.Time = new(big.Int).Add(big.NewInt(int64(parent.Time)), new(big.Int).SetUint64(1)).Uint64()
	if int64(header.Time) < time.Now().Unix() {
		header.Time = uint64(time.Now().Unix())
	}
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// Finaize doesn't modify the passed header.
func (sb *Backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	_ []*types.Header, receipts []*types.Receipt) (*types.Committee, *types.Receipt, *big.Int, error) {

	committeeSet, receipt, lastEpochBlock, err := sb.AutonityContractFinalize(header, chain, state, txs, receipts)
	if err != nil {
		return nil, nil, nil, err
	}

	// if we are at epoch change, convert the slice of CommitteeMember to types.Committee (slice of *CommitteeMember)
	// if not, leave committee to its default value
	committee := &types.Committee{}
	if len(committeeSet) != 0 {
		committee.Members = make([]*types.CommitteeMember, len(committeeSet))
		for i, m := range committeeSet {
			committee.Members[i] = &types.CommitteeMember{
				Address:      m.Address,
				VotingPower:  new(big.Int).Set(m.VotingPower),
				ValidatorKey: m.ValidatorKey,
			}
		}
		committee.Sort()
		sb.logger.Debug("Finalized epoch change block", "committee", committee)
	}

	return committee, receipt, lastEpochBlock, nil
}

// FinalizeAndAssemble call Finaize to compute post transacation state modifications
// and assembles the final block.
func (sb *Backend) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts *[]*types.Receipt) (*types.Block, error) {

	statedb.Prepare(common.ACHash(header.Number), len(txs))
	committee, receipt, lastEpochBlock, err := sb.Finalize(chain, header, statedb, txs, uncles, *receipts)
	if err != nil {
		return nil, err
	}
	*receipts = append(*receipts, receipt)
	// No block rewards in BFT, so the state remains as is and uncles are dropped
	header.Root = statedb.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// add committee to extraData's committee section
	header.Committee = committee
	header.LastEpochBlock = lastEpochBlock
	return types.NewBlock(header, txs, nil, *receipts, new(trie.Trie)), nil
}

// AutonityContractFinalize is called to deploy the Autonity Contract at block #1. it returns as well the
// committee field containaining the list of committee members allowed to participate in consensus for the next block.
func (sb *Backend) AutonityContractFinalize(header *types.Header, chain consensus.ChainReader, state *state.StateDB,
	_ []*types.Transaction, _ []*types.Receipt) ([]types.CommitteeMember, *types.Receipt, *big.Int, error) {
	sb.contractsMu.Lock()
	defer sb.contractsMu.Unlock()

	committeeSet, receipt, lastEpochBlock, err := sb.blockchain.ProtocolContracts().FinalizeAndGetCommittee(header, state)
	if err != nil {
		sb.logger.Error("Autonity Contract finalize", "err", err)
		return nil, nil, nil, err
	}
	return committeeSet, receipt, lastEpochBlock, nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *Backend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	sb.coreMu.RLock()
	isStarted := sb.coreStarted
	sb.coreMu.RUnlock()
	if !isStarted {
		return ErrStoppedEngine
	}

	// update the block header and signature and propose the block to core engine
	header := block.Header()

	epochHead, _, err := chain.EpochHeadAndParentHead(header.Number.Uint64())
	if err != nil {
		sb.logger.Error("Error ancestor")
		return consensus.ErrUnknownAncestor
	}

	nodeAddress := sb.Address()
	if epochHead.CommitteeMember(nodeAddress) == nil {
		sb.logger.Error("error validator errUnauthorized", "addr", sb.address)
		return errUnauthorized
	}

	block, err = sb.AddSeal(block)
	if err != nil {
		sb.logger.Error("seal error updateBlock", "err", err.Error())
		return err
	}

	// wait for the timestamp of header, use this to adjust the block period
	delay := time.Unix(int64(block.Header().Time), 0).Sub(now())
	select {
	case <-time.After(delay):
		// nothing to do
	case <-sb.stopped:
		return nil
	case <-stop:
		return nil
	}
	sb.setResultChan(results)
	// post block into BFT engine
	sb.postEvent(events.NewCandidateBlockEvent{
		NewCandidateBlock: *block,
	})
	return nil
}

func (sb *Backend) setResultChan(results chan<- *types.Block) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	sb.commitCh = results
}

func (sb *Backend) sendResultChan(block *types.Block) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()

	sb.commitCh <- block
}

func (sb *Backend) isResultChanNil() bool {
	sb.coreMu.RLock()
	defer sb.coreMu.RUnlock()

	return sb.commitCh == nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the blockchain and the
// current signer.
func (sb *Backend) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return defaultDifficulty
}

func (sb *Backend) SetProposedBlockHash(hash common.Hash) {
	sb.proposedBlockHash = hash
}

// update timestamp and signature of the block based on its number of transactions
func (sb *Backend) AddSeal(block *types.Block) (*types.Block, error) {
	header := block.Header()

	err := crypto.SignHeader(header, sb.privateKey)
	if err != nil {
		return nil, err
	}

	return block.WithSeal(header), nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *Backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "tendermint",
		Version:   "1.0",
		Service:   &API{chain: chain, tendermint: sb, getCommittee: getCommittee},
		Public:    true,
	}}
}

// getCommittee retrieves the committee for the given header.
func getCommittee(header *types.Header, chain consensus.ChainReader) (*types.Committee, error) {
	epochHead, _, err := chain.EpochHeadAndParentHead(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	return epochHead.Committee, nil
}

// Start implements consensus.Start
// youssef: I'm not sure about the use case of this context in argument
func (sb *Backend) Start(ctx context.Context) error {
	// the mutex along with coreStarted should prevent double start
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return ErrStartedEngine
	}
	sb.stopped = make(chan struct{})
	// clear previous data
	sb.proposedBlockHash = common.Hash{}
	// Start Tendermint
	go sb.faultyValidatorsWatcher(ctx)
	sb.wg.Add(1)
	sb.core.Start(ctx, sb.blockchain.ProtocolContracts())
	sb.coreStarted = true
	return nil
}

// Close signals core to stop all background threads.
func (sb *Backend) Close() error {
	// the mutex along with coreStarted should prevent double stop
	sb.coreMu.Lock()
	if !sb.coreStarted {
		sb.coreMu.Unlock()
		return ErrStoppedEngine
	}
	sb.coreStarted = false
	sb.coreMu.Unlock()
	// We need to make sure we close sb.stopped before calling sb.core.Stop
	// otherwise we can end up with a deadlock where sb.core.Stop is waiting
	// for a routine to return from calling sb.AskSync but sb.AskSync will
	// never return because we did not close sb.stopped.
	close(sb.stopped)
	// Stop Tendermint
	sb.core.Stop()
	sb.wg.Wait()
	return nil
}

func (sb *Backend) SealHash(header *types.Header) common.Hash {
	return types.SigHash(header)
}

func (sb *Backend) SetBlockchain(bc *core.BlockChain) {
	sb.blockchain = bc
	sb.currentBlock = bc.CurrentBlock
	sb.hasBadBlock = bc.HasBadBlock
}

func (sb *Backend) faultyValidatorsWatcher(ctx context.Context) {
	if !sb.BFT() {
		sb.logger.Info("skip watching faulty validators for none BFT consensus engine")
		return
	}
	var subscriptions event.SubscriptionScope
	newFaultProofCh := make(chan *autonity.AccountabilityNewFaultProof)
	slashingEventCh := make(chan *autonity.AccountabilitySlashingEvent)
	chainHeadCh := make(chan core.ChainHeadEvent)

	subNewFaultProofs, _ := sb.blockchain.ProtocolContracts().WatchNewFaultProof(nil, newFaultProofCh, nil)
	subSlashigEvent, _ := sb.blockchain.ProtocolContracts().WatchSlashingEvent(nil, slashingEventCh)
	subChainhead := sb.blockchain.SubscribeChainHeadEvent(chainHeadCh)
	subscriptions.Track(subNewFaultProofs)
	subscriptions.Track(subSlashigEvent)
	subscriptions.Track(subChainhead)

	defer func() {
		subscriptions.Close()
		sb.wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-subNewFaultProofs.Err():
			return
		case <-sb.stopped:
			return
		case ev := <-newFaultProofCh:
			sb.jailedLock.Lock()
			// a 0 value means that the validator is in a perpetual jailed state
			// which should only be temporary until it gets updated at the next
			// slashing event.
			sb.jailed[ev.Offender] = 0
			sb.jailedLock.Unlock()
		case ev := <-slashingEventCh:
			sb.jailedLock.Lock()
			sb.jailed[ev.Validator] = ev.ReleaseBlock.Uint64()
			sb.jailedLock.Unlock()
		case ev := <-chainHeadCh:
			sb.jailedLock.Lock()
			for k, v := range sb.jailed {
				if v < ev.Block.NumberU64() && v != 0 {
					delete(sb.jailed, k)
				}
			}
			sb.jailedLock.Unlock()
		}
	}
}

func (sb *Backend) IsJailed(address common.Address) bool {
	sb.jailedLock.RLock()
	defer sb.jailedLock.RUnlock()
	_, ok := sb.jailed[address]
	return ok
}
