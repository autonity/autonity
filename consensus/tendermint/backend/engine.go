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
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/consensus/tendermint/events"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rpc"
	"github.com/autonity/autonity/trie"
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
	sealDelayBg                   = metrics.NewRegisteredBufferedGauge("work/seal/delay", nil, nil) // injected sleep delay before producing new candidate block
)

// Author retrieves the Ethereum address of the account that minted the given
// block.
func (sb *Backend) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *Backend) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, _ bool) error {
	// Short circuit if the header is known, or its parent not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// get latest epoch info for signature checks and epoch boundary checks latter on.
	epoch, err := chain.EpochOfHeight(header.Number.Uint64())
	if err != nil {
		return err
	}

	return sb.verifyHeader(chain, header, parent, epoch.Committee, epoch.EpochBlock.Uint64(), epoch.NextEpochBlock.Uint64())
}

// verifyHeader checks whether a header conforms to the consensus rules. It expects the parent header
// to be provided unless header is the genesis header.
func (sb *Backend) verifyHeader(chain consensus.ChainHeaderReader, header, parent *types.Header,
	committee *types.Committee, curEpochHead uint64, nextEpochHead uint64) error {
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
	if parent == nil {
		return errUnknownBlock
	}

	// Re-injecting historical blocks, as epoch info is a factor to compute the hash, thus we don't need to check epoch
	// , just double check header against its parent and the quorum certificates of this height.
	if chain.GetHeader(header.Hash(), header.Number.Uint64()) != nil {
		return sb.verifyHeaderAgainstLastView(header, parent, committee)
	}

	// for unknown headers, header number should pass the corresponding epoch boundary check.
	if header.Number.Uint64() <= curEpochHead || header.Number.Uint64() > nextEpochHead {
		sb.logger.Error("header is out of epoch range",
			"height", header.Number.Uint64(), "curEpochHead", curEpochHead, "nextEpochHead", nextEpochHead)
		return consensus.ErrOutOfEpochRange
	}

	// epoch bi-direction link check for epoch header and its parent epoch header.
	if header.IsEpochHeader() {
		if nextEpochHead != header.Number.Uint64() || header.Epoch.PreviousEpochBlock.Uint64() != curEpochHead {
			return consensus.ErrInvalidEpochBoundary
		}
	}

	// check quorum certificates of consensus participants
	return sb.verifyHeaderAgainstLastView(header, parent, committee)
}

// verifyHeaderAgainstLastView verifies that the given header is valid with respect to its parent block and the
// corresponding epoch's committee.
func (sb *Backend) verifyHeaderAgainstLastView(header, parent *types.Header, committee *types.Committee) error {
	if parent.Number.Uint64() != header.Number.Uint64()-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	if parent.Time+1 > header.Time { // Todo : fetch block period from contract
		return errInvalidTimestamp
	}

	if err := sb.verifySigner(header, committee); err != nil {
		return err
	}

	return sb.verifyQuorumCertificate(header, committee)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *Backend) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, _ []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{}, 1)
	results := make(chan error, len(headers))

	go func() {
		firstHead := headers[0].Number.Uint64()
		epoch, err := chain.EpochOfHeight(firstHead)
		// short circuit, if we cannot find the correct epoch for the 1st header, we quit this batch of verification.
		if err != nil {
			sb.logger.Error("VerifyHeaders", "cannot find epoch for the 1st header of the batch: ", err.Error(), "height", firstHead)
			results <- err
			return
		}

		committee := epoch.Committee
		curEpochBlock := epoch.EpochBlock.Uint64()
		nextEpochBlock := epoch.NextEpochBlock.Uint64()
		for i, header := range headers {
			var parent *types.Header
			switch {
			case i > 0:
				parent = headers[i-1]
			case i == 0:
				parent = chain.GetHeaderByHash(header.ParentHash)
			}

			if parent == nil {
				sb.logger.Error("VerifyHeaders", "cannot find parent header", header.ParentHash)
				err = consensus.ErrUnknownAncestor
			} else {
				err = sb.verifyHeader(chain, header, parent, committee, curEpochBlock, nextEpochBlock)
			}

			// cross epoch header check, update the committee and epoch boundary if current header is an epoch head.
			// the verification behind this header will be continued with the updated epoch info.
			if header.IsEpochHeader() {
				committee = header.Epoch.Committee
				curEpochBlock = header.Number.Uint64()
				nextEpochBlock = header.Epoch.NextEpochBlock.Uint64()
			}

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
func (sb *Backend) verifySigner(header *types.Header, committee *types.Committee) error {
	// resolve the authorization key and check against signers
	signer, err := types.ECRecover(header)
	if err != nil {
		return err
	}
	if header.Coinbase != signer {
		return errInvalidCoinbase
	}
	// Signer should be in the validator set of previous block's extraData.
	if committee.MemberByAddress(signer) != nil {
		return nil
	}
	return errUnauthorized
}

// verifyQuorumCertificate validates that the quorum certificate for header come from
// committee members and that the voting power constitute a quorum.
func (sb *Backend) verifyQuorumCertificate(header *types.Header, committee *types.Committee) error {
	// un-finalized proposals will have these fields set to nil
	if header.QuorumCertificate.Signature == nil || header.QuorumCertificate.Signers == nil {
		return types.ErrEmptyQuorumCertificate
	}
	quorumCertificate := header.QuorumCertificate.Copy() // copy so that we do not modify the header when doing Signers.Validate()
	if err := quorumCertificate.Signers.Validate(committee.Len()); err != nil {
		return fmt.Errorf("Invalid quorum certificate signers information: %w", err)
	}

	// The data that was signed over for this block
	headerSeal := message.PrepareCommittedSeal(header.Hash(), int64(header.Round), header.Number)

	// Total Voting power for this block
	power := new(big.Int)
	for _, index := range quorumCertificate.Signers.FlattenUniq() {
		power.Add(power, committee.Members[index].VotingPower)
	}

	// verify signature
	var keys []blst.PublicKey //nolint
	for _, index := range quorumCertificate.Signers.Flatten() {
		keys = append(keys, committee.Members[index].ConsensusKey)
	}
	aggregatedKey, err := blst.AggregatePublicKeys(keys)
	if err != nil {
		sb.logger.Error("Failed to aggregate keys from committee members", "err", err)
		return err
	}
	valid := quorumCertificate.Signature.Verify(aggregatedKey, headerSeal[:])
	if !valid {
		sb.logger.Error("block had invalid committed seal")
		return types.ErrInvalidQuorumCertificate
	}

	// We need at least a quorum for the block to be considered valid
	if power.Cmp(bft.Quorum(committee.TotalVotingPower())) < 0 {
		return types.ErrInvalidQuorumCertificate
	}

	return nil
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *Backend) Prepare(_ consensus.ChainHeaderReader, parentHeader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = sb.Address()
	header.Nonce = emptyNonce
	header.MixDigest = types.BFTDigest
	header.Difficulty = defaultDifficulty

	// set header's timestamp
	// todo: block period from contract
	header.Time = parentHeader.Time + 1
	if int64(header.Time) < time.Now().Unix() {
		header.Time = uint64(time.Now().Unix())
	}
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// Finaize doesn't modify the passed header.
func (sb *Backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	_ []*types.Header, receipts []*types.Receipt) (*types.Receipt, *types.Epoch, error) {

	receipt, epochInfo, err := sb.AutonityContractFinalize(header, chain, state, txs, receipts)
	if err != nil {
		return nil, nil, err
	}

	return receipt, epochInfo, nil
}

// FinalizeAndAssemble call Finaize to compute post transacation state modifications
// and assembles the final block.
func (sb *Backend) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, statedb *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts *[]*types.Receipt) (*types.Block, error) {

	statedb.Prepare(common.ACHash(header.Number), len(txs))
	receipt, epochInfo, err := sb.Finalize(chain, header, statedb, txs, uncles, *receipts)
	if err != nil {
		return nil, err
	}
	*receipts = append(*receipts, receipt)
	// No block rewards in BFT, so the state remains as is and uncles are dropped
	header.Root = statedb.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash
	header.Epoch = epochInfo

	return types.NewBlock(header, txs, nil, *receipts, new(trie.Trie)), nil
}

// AutonityContractFinalize is called to deploy the Autonity Contract at block #1. it returns as well the
// committee field containaining the list of committee members allowed to participate in consensus for the next block.
func (sb *Backend) AutonityContractFinalize(header *types.Header, chain consensus.ChainReader, state *state.StateDB,
	_ []*types.Transaction, _ []*types.Receipt) (*types.Receipt, *types.Epoch, error) {

	receipt, epochInfo, err := sb.blockchain.ProtocolContracts().FinalizeAndGetCommittee(header, state)
	if err != nil {
		sb.logger.Error("Autonity Contract finalize", "err", err)
		return nil, nil, err
	}

	return receipt, epochInfo, nil
}

func (sb *Backend) EpochByHeight(height *big.Int) (*types.EpochInfo, error) {
	header := sb.BlockChain().CurrentHeader()
	stateDB, err := sb.blockchain.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	return sb.BlockChain().ProtocolContracts().EpochByHeight(header, stateDB, height)
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *Backend) Seal(parent *types.Header, block *types.Block, _ chan<- *types.Block, stop <-chan struct{}) error {
	if !sb.coreRunning.Load() {
		return ErrStoppedEngine
	}

	if parent == nil {
		err := errors.New("unknown ancestor")
		return err
	}

	// we do the validator authorization later, just before sending proposal
	block, err := sb.AddSeal(block)
	if err != nil {
		sb.logger.Error("sealing error", "err", err.Error())
		return err
	}

	// wait for the timestamp of header, use this to adjust the block period
	delay := time.Unix(int64(block.Header().Time), 0).Sub(now())
	if metrics.Enabled {
		sealDelayBg.Add(delay.Nanoseconds())
	}
	select {
	case <-time.After(delay):
		// nothing to do
	case <-sb.stopped:
		return nil
	case <-stop:
		return nil
	}

	// post block into BFT engine
	sb.Post(events.NewCandidateBlockEvent{
		NewCandidateBlock: *block,
		CreatedAt:         time.Now(),
	})

	return nil
}

func (sb *Backend) SetProposalVerifiedEventChan(proposalVerifiedCh chan<- *types.Block) {
	sb.proposalVerifiedCh = proposalVerifiedCh
}

func (sb *Backend) ProposalVerified(block *types.Block) {
	sb.proposalVerifiedCh <- block
}

func (sb *Backend) IsProposalStateCached(hash common.Hash) bool {
	return sb.blockchain.IsProposalStateCached(hash)
}

func (sb *Backend) SetResultChan(results chan<- *types.Block) {
	sb.commitCh = results
}

func (sb *Backend) sendResultChan(block *types.Block) {
	sb.commitCh <- block
}

func (sb *Backend) isResultChanNil() bool {
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

func (sb *Backend) ProposedBlockHash() common.Hash {
	return sb.proposedBlockHash
}

// AddSeal update timestamp and signature of the block based on its number of transactions
func (sb *Backend) AddSeal(block *types.Block) (*types.Block, error) {
	header := block.Header()
	hashData := types.SigHash(header)
	signature, err := crypto.Sign(hashData[:], sb.nodeKey)
	if err != nil {
		return nil, err
	}
	if err := types.WriteSeal(header, signature); err != nil {
		return nil, err
	}
	return block.WithSeal(header), nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *Backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "tendermint",
		Version:   "1.0",
		Service:   &API{chain: chain, tendermint: sb},
		Public:    true,
	}}
}

// Start implements consensus.Start
// youssef: I'm not sure about the use case of this context in argument
func (sb *Backend) Start(ctx context.Context) error {
	if !sb.coreStarting.CompareAndSwap(false, true) {
		return ErrStartedEngine
	}

	sb.stopped = make(chan struct{})
	sb.UpdateStopChannel(sb.stopped)
	// clear previous data
	sb.proposedBlockHash = common.Hash{}

	sb.wg.Add(1)
	go sb.faultyValidatorsWatcher(ctx)

	// Start Tendermint
	sb.aggregator.start(ctx)
	sb.core.Start(ctx, sb.blockchain.ProtocolContracts())
	sb.coreRunning.CompareAndSwap(false, true)
	return nil
}

// Close signals core to stop all background threads.
func (sb *Backend) Close() error {
	if !sb.coreRunning.CompareAndSwap(true, false) {
		return ErrStoppedEngine
	}
	// We need to make sure we close sb.stopped before calling sb.core.Stop
	// otherwise we can end up with a deadlock where sb.core.Stop is waiting
	// for a routine to return from calling sb.AskSync but sb.AskSync will
	// never return because we did not close sb.stopped.
	close(sb.stopped)
	// Stop Tendermint
	sb.aggregator.stop()
	sb.core.Stop()
	sb.wg.Wait()
	sb.coreStarting.CompareAndSwap(true, false)
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
			// a fault proof against our own node has been finalized on-chain
			// we cannot do anything about it now, let's just write a summary for the validator operator
			if ev.Offender == sb.address {
				event, err := sb.blockchain.ProtocolContracts().Events(nil, ev.Id)
				if err != nil {
					// this should never happen
					sb.logger.Crit("Can't retrieve accountability event", "id", ev.Id)
				}
				eventType := autonity.AccountabilityEventType(event.EventType).String()
				rule := autonity.Rule(event.Rule).String()
				explanation := autonity.Rule(event.Rule).Explanation()
				sb.logger.Warn("Your validator has been found guilty of consensus misbehaviour", "address", event.Offender, "event id", ev.Id.Uint64(), "event type", eventType, "rule", rule, "block", event.Block.Uint64(), "epoch", event.Epoch.Uint64(), "faulty message hash", common.BigToHash(event.MessageHash))
				sb.logger.Warn(explanation)
			}
			sb.jailedLock.Lock()
			// a 0 value means that the validator is in a perpetual jailed state
			// which should only be temporary until it gets updated at the next
			// slashing event.
			sb.jailed[ev.Offender] = 0
			sb.jailedLock.Unlock()
		case ev := <-slashingEventCh:
			// local node got slashed, print out information about the slashing that can be correlated with the information about the fault proof above.
			if ev.Validator == sb.address {
				sb.logger.Warn("Your validator has been slashed", "amount", ev.Amount.Uint64(), "jail release block", ev.ReleaseBlock.Uint64(), "jailbound", ev.IsJailbound, "event id", ev.EventId.Uint64())
			}
			sb.jailedLock.Lock()
			if ev.IsJailbound {
				// the validator is jailed permanently, won't be able to enter committee and
				// his messages will be discarded at validation, no need to keep track
				delete(sb.jailed, ev.Validator)
			} else {
				sb.jailed[ev.Validator] = ev.ReleaseBlock.Uint64()
			}
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
