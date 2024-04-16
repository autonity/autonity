package collusion

import (
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
	"math/big"
	"math/rand"
	"sync"
)

/**
 * Collusion test, it simulates 1/3 faulty validators working together to manipulate voting of consensus.
 * The collusion framework is base on a faulty party which has 1/3 faulty nodes in it, the framework select a leader out
 * from the faulty nodes, as a leader of faulty nodes, it needs to coordinate the misbehaviour of all faulty followers
 * by according to the accountability rules in the test context, by putting the messages in the queues of each steps of
 * each faulty followers, the followers pick up them from the queue of per step then broadcast it to the committee. Thus,
 * we have such collusion simulated to verify if the accountability module can detect such kind of misbehaviour and to
 * verify the correctness and live-ness of consensus.
 */

type collusionPlaner interface {
	plan(h uint64, r int64, v *types.Block, members []*gengen.Validator) (map[common.Address]map[uint64]map[int64]map[core.Step]message.Msg, *gengen.Validator, []*gengen.Validator)
}

// configurations for different testcases
var colludedActions = map[autonity.Rule]*colludedBehaviours{}
var colludedActionsLock = sync.RWMutex{}

type colludedBehaviours struct {
	leader       *gengen.Validator
	followers    []*gengen.Validator
	h            uint64
	r            int64
	rule         autonity.Rule
	invalidValue *types.Block
	//lock         sync.RWMutex
	// message queues for faulty members, each member will consume message and broadcast it.
	colludedBehaviours map[common.Address]map[uint64]map[int64]map[core.Step]message.Msg
}

func collusionBehaviour(rule autonity.Rule) *colludedBehaviours {
	colludedActionsLock.RLock()
	defer colludedActionsLock.RUnlock()
	return colludedActions[rule]
}

/*
func (cp *colludedBehaviours) removeMessage(actor common.Address, height uint64, round int64, step core.Step) message.Msg {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	// Check if the actor exists in the map
	if actorMap, ok := cp.colludedBehaviours[actor]; ok {
		// Check if the height exists in the actor map
		if heightMap, ok := actorMap[height]; ok {
			// Check if the round exists in the height map
			if roundMap, ok := heightMap[round]; ok {
				// Check if the step exists in the round map
				if msg, ok := roundMap[step]; ok {
					// Remove the message from the map
					delete(roundMap, step)
					return msg
				}
			}
		}
	}
	return nil
}*/

func newBlockHeader(height uint64) *types.Header {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[i] = byte(rand.Intn(256)) //nolint
	}
	return &types.Header{
		Number: new(big.Int).SetUint64(height),
		Nonce:  nonce,
	}
}

func createCollusionParty(vals []*gengen.Validator, h uint64, r int64, rule autonity.Rule, planer collusionPlaner) {
	f := bft.F(new(big.Int).SetUint64(uint64(len(vals))))
	if f.Uint64() < 2 {
		panic("collusion test requires at least two faulty nodes")
	}
	var faultyMembers []*gengen.Validator
	for i := uint64(0); i < f.Uint64(); i++ {
		faultyMembers = append(faultyMembers, vals[i])
	}
	header := newBlockHeader(h)
	block := types.NewBlockWithHeader(header)
	e2e.FuzBlock(block, new(big.Int).SetUint64(h))

	b := &colludedBehaviours{
		h:            h,
		r:            r,
		rule:         rule,
		invalidValue: block,
	}

	b.colludedBehaviours, b.leader, b.followers = planer.plan(h, r, block, faultyMembers)
	colludedActionsLock.Lock()
	defer colludedActionsLock.Unlock()
	colludedActions[rule] = b
}
