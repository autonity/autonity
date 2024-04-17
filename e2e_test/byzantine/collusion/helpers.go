package collusion

import (
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
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
	setupRoles(members []*gengen.Validator) (*gengen.Validator, []*gengen.Validator)
}

// configurations for different testcases
var collusionContexts = map[autonity.Rule]*collusionContext{}
var collusionContextsLock = sync.RWMutex{}

type collusionContext struct {
	leader       *gengen.Validator
	followers    []*gengen.Validator
	h            uint64
	r            int64
	rule         autonity.Rule
	invalidValue *types.Block
	lock         sync.RWMutex
}

func (c *collusionContext) context() (uint64, int64, *types.Block) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.h, c.r, c.invalidValue
}

func (c *collusionContext) setupContext(h uint64, r int64, invalidValue *types.Block) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.h = h
	c.r = r
	c.invalidValue = invalidValue
}

func (c *collusionContext) isReady() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.invalidValue != nil
}

func getCollusionContext(rule autonity.Rule) *collusionContext {
	collusionContextsLock.RLock()
	defer collusionContextsLock.RUnlock()
	return collusionContexts[rule]
}

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

func initCollusionContext(vals []*gengen.Validator, rule autonity.Rule, planer collusionPlaner) {
	f := bft.F(new(big.Int).SetUint64(uint64(len(vals))))
	if f.Uint64() < 2 {
		panic("collusion test requires at least two faulty nodes")
	}
	var faultyMembers []*gengen.Validator
	for i := uint64(0); i < f.Uint64(); i++ {
		faultyMembers = append(faultyMembers, vals[i])
	}

	b := &collusionContext{
		rule: rule,
	}
	b.leader, b.followers = planer.setupRoles(faultyMembers)

	collusionContextsLock.Lock()
	defer collusionContextsLock.Unlock()
	collusionContexts[rule] = b
}

func validProposer(address common.Address, h uint64, r int64, core *core.Core) bool {
	contract := core.Backend().BlockChain().ProtocolContracts()
	db := core.Backend().BlockChain().StateCache()
	statedb, err := state.New(core.LastHeader().Root, db, nil)
	if err != nil {
		panic("cannot load state from block chain.")
	}
	return address == contract.Proposer(core.LastHeader(), statedb, h, r)
}
