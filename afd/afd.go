package afd

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"math/rand"
	"sync"
	"time"
)

var (
	errFutureMsg = errors.New("future height msg")
)

// Fault detector, it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	// use below 3 members to send proof via transaction issuing.
	wg sync.WaitGroup
	afdFeed event.Feed
	scope event.SubscriptionScope

	// use below 2 members to forward consensus msg from protocol manager to afd.
	tendermintMsgSub *event.TypeMuxSubscription
	tendermintMsgMux *event.TypeMuxSilent

	// below 2 members subscribe block event to trigger execution
	// of rule engine and make proof of innocent.
	blockChan chan core.ChainEvent
	blockSub event.Subscription

	// chain context to validate consensus msgs.
	blockchain *core.BlockChain

	// node address
	address common.Address

	// msg store
	msgStore *MsgStore

	// rule engine
	ruleEngine *RuleEngine

	logger log.Logger
}

var(
	// todo: to be configured at genesis.
	randomDelayWindow = 10000
)

// call by ethereum object to create fd instance.
func NewFaultDetector(chain *core.BlockChain, nodeAddress common.Address) *FaultDetector {
	logger := log.New("afd_addr", nodeAddress)
	fd := &FaultDetector{
		address: nodeAddress,
		blockChan:  make(chan core.ChainEvent, 300),
		blockchain: chain,
		msgStore: new(MsgStore),
		ruleEngine: new(RuleEngine),
		logger:logger,
		tendermintMsgMux:  event.NewTypeMuxSilent(logger),
	}

	// init accountability precompiled contracts.
	initAccountabilityContracts(chain)

	// subscribe tendermint msg
	s := fd.tendermintMsgMux.Subscribe(events.MessageEvent{})
	fd.tendermintMsgSub = s
	return fd
}

// listen for new block events from block-chain, do the tasks like take challenge and provide proof for innocent, the
// AFD rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)

	for {
		select {
		case ev := <-fd.blockChan:
			// take my challenge from latest state DB, and provide innocent proof if there are any.
			err := fd.handleMyChallenges(ev.Block, ev.Hash)
			if err != nil {
				fd.logger.Warn("handle challenge","error", err)
			}

			// todo: tell rule engine to run patterns over msg store on each new height.
			fd.ruleEngine.run()
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			// take consensus msg from p2p protocol manager.
			if !ok {
				return
			}
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				msg := new(types.ConsensusMessage)
				if err := msg.FromPayload(e.Payload); err != nil {
					fd.logger.Error("invalid payload", "err", err)
					continue
				}

				if err := fd.processMsg(msg); err != nil {
					fd.logger.Error("process consensus msg", "err", err)
					continue
				}
			}

		case <-fd.blockSub.Err():
			return
		}
	}
}

// HandleMsg is called by p2p protocol manager to deliver the consensus msg to afd.
func (fd *FaultDetector) HandleMsg(addr common.Address, msg p2p.Msg) {
	if msg.Code != types.TendermintMsg {
		return
	}

	var data []byte
	if err := msg.Decode(&data); err != nil {
		log.Error("cannot decode consensus msg", "from", addr)
		return
	}

	// post consensus event to event loop.
	fd.tendermintMsgMux.Post(events.MessageEvent{Payload:data})
	return
}

func (fd *FaultDetector) Stop() {
	fd.scope.Close()
	fd.blockSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.wg.Wait()
	cleanContracts()
}

// call by ethereum object to subscribe proofs Events.
func (fd *FaultDetector) SubscribeAFDEvents(ch chan<- types.SubmitProofEvent) event.Subscription {
	return fd.scope.Track(fd.afdFeed.Subscribe(ch))
}

// processMsg it decode consensus msg, apply auto-incriminating, equivocation rules to it,
// and store it to msg store.
func (fd *FaultDetector) processMsg(m *types.ConsensusMessage) error {
	err := fd.preProcessMsg(m)
	if err != nil {
		if err == errFutureMsg {
			// todo: buffer the msg until we get synced with latest block,
			//  then process it.
			return nil
		}
		return err
	}

	// todo: pre-process proposal, prevote, precommit for auto-incriminating, equivocation.

	// todo: save valid msg into msg-store.

	return nil
}

//pre-process msg, it check if msg is from valid member of the committee, it return
func (fd *FaultDetector) preProcessMsg(m *types.ConsensusMessage) error {
	msgHeight, err := m.Height()
	if err != nil {
		return err
	}

	header := fd.blockchain.CurrentHeader()
	if msgHeight.Cmp(header.Number) > 0 {
		return errFutureMsg
	}

	lastHeader := fd.blockchain.GetHeaderByNumber(msgHeight.Uint64())

	if _, err = m.Validate(crypto.CheckValidatorSignature, lastHeader); err != nil {
		fd.logger.Error("Msg is not from committee member", "err", err)
		return err
	}
	return nil
}

// get challenges from blockchain via autonityContract calls.
func (fd *FaultDetector) handleMyChallenges(block *types.Block, hash common.Hash) error {
	var innocentProofs []types.OnChainProof
	state, err := fd.blockchain.StateAt(hash)
	if err != nil {
		return err
	}

	challenges := fd.blockchain.GetAutonityContract().GetChallenges(block.Header(), state)
	for i:=0; i < len(challenges); i++ {
		if challenges[i].Sender == fd.address {
			p, err := fd.proveInnocent(challenges[i])
			if err != nil {
				continue
			}
			innocentProofs = append(innocentProofs, p)
		}
	}

	// send proofs via standard transaction.
	fd.sendProofs(types.InnocentProof, innocentProofs)
	return nil
}

// get proof of innocent over msg store.
func (fd *FaultDetector) proveInnocent(challenge types.OnChainProof) (types.OnChainProof, error) {
	// todo: get proof from msg store over the rule.
	var proof types.OnChainProof
	return proof, nil
}

// send proofs via event which will handled by ethereum object to signed the TX to send proof.
func (fd *FaultDetector) sendProofs(t types.ProofType,  proofs[]types.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		if t == types.InnocentProof {
			fd.afdFeed.Send(types.SubmitProofEvent{Proofs:proofs, Type:t})
		}

		if t == types.ChallengeProof {
			// wait for random milliseconds (under the range of 10 seconds) to check if need to rise challenge.
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(randomDelayWindow)
			time.Sleep(time.Duration(n)*time.Millisecond)
			unPresented := fd.filterUnPresentedChallenges(&proofs)
			if len(unPresented) != 0 {
				fd.afdFeed.Send(types.SubmitProofEvent{Proofs:unPresented, Type:t})
			}
		}
	}()
}

func (fd *FaultDetector) filterUnPresentedChallenges(proofs *[]types.OnChainProof) []types.OnChainProof {
	// get latest chain state.
	var result []types.OnChainProof
	state, err := fd.blockchain.State()
	if err != nil {
		return nil
	}

	header := fd.blockchain.CurrentBlock().Header()
	challenges := fd.blockchain.GetAutonityContract().GetChallenges(header, state)

	for i:=0; i < len(*proofs); i++ {
		present := false
		for i:=0; i < len(challenges); i++ {
			if (*proofs)[i].Sender == challenges[i].Sender && (*proofs)[i].Height == challenges[i].Height &&
				(*proofs)[i].Round == challenges[i].Round && (*proofs)[i].Rule == challenges[i].Rule &&
				(*proofs)[i].MsgType == challenges[i].MsgType {
				present = true
			}
		}
		if !present {
			result = append(result, (*proofs)[i])
		}
	}

	return result
}
