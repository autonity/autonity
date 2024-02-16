package accountability

import (
	"errors"
	"fmt"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/eth/protocols/eth"
	"github.com/autonity/autonity/rlp"
)

var (
	errNoParentHeader              = errors.New("no parent header")
	errInvalidAccusation           = errors.New("invalid accusation")
	errPeerDuplicatedAccusation    = errors.New("remote peer is sending duplicated accusation")
	errInvalidInnocenceProof       = errors.New("invalid proof of innocence")
	errAccusationRateMalicious     = errors.New("malicious accusation msg rate, peer to be dropped")
	errAccusationFromNoneValidator = errors.New("accusation from none validator node")
)

type AccusationRateLimiter struct {
	// to track the number of accusation sent by a challenger over a specific height.
	accusationsPerHeight map[common.Address]map[uint64]int
	// malicious one might use out of updated accusations to DoS node, so we track those recently accusation rates.
	accusationRates map[common.Address]int
	// track if one send duplicated accusation.
	peerProcessedAccusations map[common.Address]map[common.Hash]struct{}
}

func NewAccusationRateLimiter() *AccusationRateLimiter {
	l := &AccusationRateLimiter{
		accusationRates:          make(map[common.Address]int),
		accusationsPerHeight:     make(map[common.Address]map[uint64]int),
		peerProcessedAccusations: make(map[common.Address]map[common.Hash]struct{}),
	}
	return l
}

// although we have rate limit over per height, but since malicious node can use out of updated consensus msg to send
// accusation to DoS node, thus we have to track the rate limit of accusation during the recently period, 1 seconds.
func (r *AccusationRateLimiter) validAccusationRate(sender common.Address) error {
	// get accusation counters of the last 1 seconds.
	times, ok := r.accusationRates[sender]
	if !ok {
		r.accusationRates[sender] = 1
		return nil
	}

	// since communication channel is asynchronous, those pending write of off chain accusation msgs from a sender could
	// potentially be received once the network session get established from a disaster, thus it could exceed the number
	// of accusation that could be produced by rule engine over a height, so we set higher rate limit during 1 second
	// to be tolerant for such case.
	if times > maxAccusationRatePerHeight*2 {
		return errAccusationRateMalicious
	}

	r.accusationRates[sender]++
	return nil
}

// rate limit counters are reset on each 1 seconds.
func (r *AccusationRateLimiter) resetRateLimiter() {
	for k := range r.accusationRates {
		delete(r.accusationRates, k)
	}
}

func (r *AccusationRateLimiter) checkPeerDuplicatedAccusation(sender common.Address, msgHash common.Hash) error {
	msgMap, ok := r.peerProcessedAccusations[sender]
	if !ok {
		msgMap = make(map[common.Hash]struct{})
		r.peerProcessedAccusations[sender] = msgMap
		r.peerProcessedAccusations[sender][msgHash] = struct{}{}
		return nil
	}
	_, ok = msgMap[msgHash]
	if !ok {
		msgMap[msgHash] = struct{}{}
		return nil
	}
	return errPeerDuplicatedAccusation
}

// justified accusations are reset on every 60 blocks.
func (r *AccusationRateLimiter) resetPeerJustifiedAccusations() {
	for k := range r.peerProcessedAccusations {
		delete(r.peerProcessedAccusations, k)
	}
}

// reset rate limiter of per height on each 60 blocks.
func (r *AccusationRateLimiter) resetHeightRateLimiter() {
	for k := range r.accusationsPerHeight {
		delete(r.accusationsPerHeight, k)
	}
}

func (r *AccusationRateLimiter) checkHeightAccusationRate(sender common.Address, height uint64) error {
	hMap, ok := r.accusationsPerHeight[sender]
	if !ok {
		hMap = make(map[uint64]int)
		r.accusationsPerHeight[sender] = hMap
		r.accusationsPerHeight[sender][height] = 1
		return nil
	}

	times, ok := hMap[height]
	if !ok {
		hMap[height] = 1
		return nil
	}
	hMap[height] = times + 1

	if hMap[height] > maxAccusationRatePerHeight {
		return errAccusationRateMalicious
	}

	return nil
}

type InnocenceProofBuffer struct {
	accusationList []common.Hash
	proofs         map[common.Hash][]byte
}

func NewInnocenceProofBuffer() *InnocenceProofBuffer {
	buff := &InnocenceProofBuffer{
		proofs: make(map[common.Hash][]byte),
	}
	return buff
}

func (i *InnocenceProofBuffer) cacheInnocenceProof(challengeHash common.Hash, payload []byte) {
	if len(i.accusationList) >= maxNumOfInnocenceProofCached {
		// remove the LRU one.
		delete(i.proofs, i.accusationList[0])
		i.accusationList = i.accusationList[1:]
	}
	i.accusationList = append(i.accusationList, challengeHash)
	i.proofs[challengeHash] = payload
}

func (i *InnocenceProofBuffer) getInnocenceProofFromCache(challengeHash common.Hash) []byte {
	proof, ok := i.proofs[challengeHash]
	if !ok {
		return nil
	}
	return proof
}

// this function take accountability events: an accusation or an innocence proof event to handle off chain accusation
// protocol. It returns error to freeze remote peer for 30s by according to dev p2p protocol to prevent from DoS attack.
func (fd *FaultDetector) handleOffChainAccountabilityEvent(payload []byte, sender common.Address) error {
	// drop peer if the accusation exceed the rate limit during the last 1 seconds.
	err := fd.rateLimiter.validAccusationRate(sender)
	if err != nil {
		fd.logger.Error("accountability abuse detected!", "err", err)
		return err
	}

	// drop peer if it sent duplicated accusation event.
	msgHash := crypto.Hash(payload)
	err = fd.rateLimiter.checkPeerDuplicatedAccusation(sender, msgHash)
	if err != nil {
		fd.logger.Error("duplicated accusation from peer", "err", err)
		return err
	}

	// send response if we have processed it recently.
	cachedProof := fd.innocenceProofBuff.getInnocenceProofFromCache(msgHash)
	if cachedProof != nil {
		fd.sendOffChainInnocenceProof(sender, cachedProof)
		return nil
	}

	// handle a brand-new accusation event from rlp decoding of proof.
	proof, err := decodeRawProof(payload)
	if err != nil {
		return err
	}

	// drop peer if the event is not from validator node.
	msgHeight := proof.Message.H()
	lastHeader := fd.blockchain.GetHeaderByNumber(msgHeight - 1)
	if lastHeader == nil {
		return errNoParentHeader
	}
	// TODO(lorenzo) verify if this can be remvoed. It is checked again in `verifyProofSignatures` but maybe
	// we need it also here for the rate limiting. Actually sender here is the p2p sender?
	memberShip := lastHeader.CommitteeMember(sender)
	if memberShip == nil {
		return errAccusationFromNoneValidator
	}

	// drop peer if one send more than the number of accusations could be produced by rule engine over a height.
	err = fd.rateLimiter.checkHeightAccusationRate(sender, msgHeight)
	if err != nil {
		fd.logger.Info("over rated accusation over a height", "error", err)
		return err
	}
	if err = verifyProofSignatures(fd.blockchain, proof); err != nil {
		return err
	}
	// handle accusation and provide innocence proof.
	if proof.Type == autonity.Accusation {
		return fd.handleOffChainAccusation(proof, sender, msgHash)
	}

	// handle innocence proof and to withdraw those pending accusation.
	if proof.Type == autonity.Innocence {
		return fd.handleOffChainProofOfInnocence(proof, sender)
	}
	return fmt.Errorf("wrong proof type for off chain accusation events")
}

func (fd *FaultDetector) handleOffChainAccusation(accusation *Proof, sender common.Address, accusationHash common.Hash) error {
	// if the suspected msg's sender is not current peer, then it would be a DoS attack, drop the peer with an error returned.
	if accusation.Message.Sender() != fd.address {
		return errInvalidAccusation
	}

	// TODO(lorenzo) decide whether to disconnect and whether to merge 1st and 3rd param
	if err := preVerifyAccusation(fd.blockchain, accusation.Message, fd.blockchain.CurrentBlock().NumberU64()); err != nil {
		return nil
	}

	// check if the accusation sent by remote peer is valid or not, an invalid accusation will drop sender's peer.
	if !verifyAccusation(accusation) {
		return errInvalidAccusation
	}

	// query innocence proof for accusation from msg store.
	ev, err := fd.innocenceProof(accusation)
	if err != nil {
		fd.logger.Warn("cannot collect ev of innocence for the accusation", "err", err)
		return nil
	}

	if len(ev.RawProof) > eth.MaxMessageSize {
		fd.logger.Error("the innocence ev is over size than 10MB")
		return nil
	}

	// buffer the innocence proof for such accusation in case of same accusation from other peers.
	fd.innocenceProofBuff.cacheInnocenceProof(accusationHash, ev.RawProof)
	// send the innocence proof to challenger.
	fd.sendOffChainInnocenceProof(sender, ev.RawProof)
	return nil
}

func (fd *FaultDetector) handleOffChainProofOfInnocence(proof *Proof, sender common.Address) error {
	// if the sender is not the one being challenged against, then drop the peer by returning error.
	if proof.Message.Sender() != sender {
		return errInvalidInnocenceProof
	}
	// check if the proof is valid, an invalid proof of innocence will freeze the peer connection.
	if !verifyInnocenceProof(proof, fd.blockchain) {
		return errInvalidInnocenceProof
	}
	// the proof is valid, withdraw the off chain challenge.
	fd.removeOffChainAccusation(proof)
	return nil
}

func (fd *FaultDetector) addOffChainAccusation(accusation *Proof) {
	fd.offChainAccusationsMu.Lock()
	defer fd.offChainAccusationsMu.Unlock()
	fd.offChainAccusations = append(fd.offChainAccusations, accusation)
}

// remove off chain accusation is called when there is valid innocence proof been received or on when the timer is expired,
// the accusation to be escalated as on chain accusation.
func (fd *FaultDetector) removeOffChainAccusation(innocenceProof *Proof) {
	fd.offChainAccusationsMu.Lock()
	defer fd.offChainAccusationsMu.Unlock()
	i := 0
	find := false
	for ; i < len(fd.offChainAccusations); i++ {
		if fd.offChainAccusations[i].Rule == innocenceProof.Rule && fd.offChainAccusations[i].Type == autonity.Accusation &&
			fd.offChainAccusations[i].Message.Hash() == innocenceProof.Message.Hash() {
			// release the event's memory
			fd.offChainAccusations[i] = nil
			find = true
			break
		}
	}

	// release the pointer from slice.
	if find {
		fd.offChainAccusations = append(fd.offChainAccusations[:i], fd.offChainAccusations[i+1:]...)
	}
}

func (fd *FaultDetector) getExpiredOffChainAccusation(currentChainHeight uint64) []*Proof {
	fd.offChainAccusationsMu.RLock()
	defer fd.offChainAccusationsMu.RUnlock()
	var expiredOnes []*Proof
	for _, proof := range fd.offChainAccusations {
		// NOTE: accusations for message at height h is generated at height h + delta by the fault detector
		// then we have up to h + delta + offchainWindow to resolve it offchain
		if currentChainHeight-proof.Message.H() > (DeltaBlocks + offChainAccusationProofWindow) {
			expiredOnes = append(expiredOnes, proof)
		}
	}
	return expiredOnes
}

// if those off chain challenge have no innocence proof within the proof window, then escalate them on-chain.
func (fd *FaultDetector) escalateExpiredAccusations(currentChainHeight uint64) {
	escalatedOnes := fd.getExpiredOffChainAccusation(currentChainHeight)
	for _, accusation := range escalatedOnes {
		fd.removeOffChainAccusation(accusation)
		p := fd.eventFromProof(accusation)
		// push it to the on chain accountability event list
		fd.pendingEvents = append(fd.pendingEvents, p)
	}
}

// send the off chain accusation msg to the peer suspected
func (fd *FaultDetector) sendOffChainAccusationMsg(accusation *Proof) {
	// send the off chain accusation msg to the suspected one,
	if fd.broadcaster == nil {
		fd.logger.Warn("p2p protocol handler is not ready yet")
		return
	}

	targets := make(map[common.Address]struct{})
	targets[accusation.Message.Sender()] = struct{}{}
	peers := fd.broadcaster.FindPeers(targets)
	if len(peers) == 0 {
		//todo: if we need to gossip this message in case of there are no direct peer connection.
		fd.logger.Debug("No direct p2p connection with suspect")
		return
	}

	rProof, err := rlp.EncodeToBytes(accusation)
	if err != nil {
		fd.logger.Warn("cannot rlp encode accusation", "err", err)
		return
	}

	fd.logger.Info("Attempting direct p2p resolution..", "suspect", accusation.Message.Sender())
	go peers[accusation.Message.Sender()].Send(backend.AccountabilityNetworkMsg, rProof) //nolint
}

// sendOffChainInnocenceProof, send an innocence proof to receiver peer.
func (fd *FaultDetector) sendOffChainInnocenceProof(receiver common.Address, payload []byte) {
	if fd.broadcaster == nil {
		fd.logger.Warn("p2p protocol handler is not ready yet")
		return
	}
	targets := make(map[common.Address]struct{})
	targets[receiver] = struct{}{}
	peers := fd.broadcaster.FindPeers(targets)
	if len(peers) == 0 {
		//todo: if we need to gossip this message in case of there are no direct peer connection.
		fd.logger.Debug("no peer connection for off chain innocence proof event")
		return
	}

	fd.logger.Info("Sending requested innocence proof", "addr", receiver)
	go peers[receiver].Send(backend.AccountabilityNetworkMsg, payload) //nolint
}
