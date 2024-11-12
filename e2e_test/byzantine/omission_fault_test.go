package byzantine

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/backend"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	tdmcore "github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	e2e "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/ethclient"
	"github.com/autonity/autonity/params"
)

const active = uint8(0)
const jailed = uint8(2)
const jailbound = uint8(3)
const jailedForInactivity = uint8(4)
const jailboundForInactivity = uint8(5)

func createNetwork(t *testing.T, nodes int, start bool, options ...gengen.GenesisOption) e2e.Network {
	validators, err := e2e.Validators(t, nodes, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	network, err := e2e.NewNetworkFromValidators(t, validators, start, options...)
	require.NoError(t, err)
	return network
}

// returns contract bindings
func contracts(t *testing.T, node *e2e.Node) (*autonity.Autonity, *autonity.OmissionAccountability) {
	endPoint := node.WsClient
	autonityContract, err := autonity.NewAutonity(params.AutonityContractAddress, endPoint)
	require.NoError(t, err)
	omissionContract, err := autonity.NewOmissionAccountability(params.OmissionAccountabilityContractAddress, endPoint)
	require.NoError(t, err)
	return autonityContract, omissionContract
}

func epochID(t *testing.T, autonity *autonity.Autonity) uint64 {
	epochID, err := autonity.EpochID(nil)
	require.NoError(t, err)
	return epochID.Uint64()
}

func inactivityScore(t *testing.T, omission *autonity.OmissionAccountability, validator common.Address) uint64 {
	score, err := omission.GetInactivityScore(nil, validator)
	require.NoError(t, err)
	return score.Uint64()
}

func effort(t *testing.T, omission *autonity.OmissionAccountability, validator common.Address) *big.Int {
	effort, err := omission.ProposerEffort(nil, validator)
	require.NoError(t, err)
	return effort
}

func totalEffort(t *testing.T, omission *autonity.OmissionAccountability) *big.Int {
	effort, err := omission.TotalEffort(nil)
	require.NoError(t, err)
	return effort
}

func inactivityCounter(t *testing.T, omission *autonity.OmissionAccountability, validator common.Address) uint64 {
	score, err := omission.InactivityCounter(nil, validator)
	require.NoError(t, err)
	return score.Uint64()
}

func validator(t *testing.T, autonity *autonity.Autonity, validator common.Address) autonity.AutonityValidator {
	val, err := autonity.GetValidator(nil, validator)
	require.NoError(t, err)
	return val
}

func committee(t *testing.T, autonity *autonity.Autonity) []autonity.AutonityCommitteeMember {
	committee, err := autonity.GetCommittee(nil)
	require.NoError(t, err)
	return committee
}

func offences(t *testing.T, omission *autonity.OmissionAccountability, validator common.Address) uint64 {
	offences, err := omission.RepeatedOffences(nil, validator)
	require.NoError(t, err)
	return offences.Uint64()
}

func probation(t *testing.T, omission *autonity.OmissionAccountability, validator common.Address) uint64 {
	probation, err := omission.ProbationPeriods(nil, validator)
	require.NoError(t, err)
	return probation.Uint64()
}

func collusionDegree(t *testing.T, omission *autonity.OmissionAccountability, epochID uint64) uint64 {
	collusionDegree, err := omission.EpochCollusionDegree(nil, new(big.Int).SetUint64(epochID))
	require.NoError(t, err)
	return collusionDegree.Uint64()
}

func isAbsent(t *testing.T, omission *autonity.OmissionAccountability, height uint64, validator common.Address) bool {
	isAbsent, err := omission.InactiveValidators(nil, new(big.Int).SetUint64(height), validator)
	require.NoError(t, err)
	return isAbsent
}

// checks if the proposer was marked as faulty for the passed height
func IsProposerFaulty(t *testing.T, omission *autonity.OmissionAccountability, height uint64) bool {
	isFaulty, err := omission.FaultyProposers(nil, new(big.Int).SetUint64(height))
	require.NoError(t, err)
	return isFaulty
}

func scaleFactor(t *testing.T, omission *autonity.OmissionAccountability) uint64 {
	scaleFactor, err := omission.GetScaleFactor(nil)
	require.NoError(t, err)
	return scaleFactor.Uint64()
}

// assumes usage of default params and mimic solidity precision loss
func computeInactivityScore(inactiveBlocks uint64, scaleFactor uint64) uint64 {
	newFloat := func(value uint64) *big.Float {
		return new(big.Float).SetPrec(256).SetUint64(value)
	}

	pastPerformanceWeight := new(big.Float).Quo(newFloat(defaultPastPerformanceWeight), newFloat(scaleFactor))
	currentPerformanceWeight := new(big.Float).Sub(newFloat(1), pastPerformanceWeight)

	score := new(big.Float).Quo(newFloat(inactiveBlocks), newFloat(defaultEpochPeriod-defaultDelta-defaultDelta+1))
	expectedInactivityScoreFloat := new(big.Float).Mul(score, currentPerformanceWeight)
	expectedInactivityScoreFloatScaled := new(big.Float).Mul(expectedInactivityScoreFloat, newFloat(scaleFactor))
	expectedInactivityScore, _ := expectedInactivityScoreFloatScaled.Uint64()
	return expectedInactivityScore
}

func atnBalanceOf(t *testing.T, wsClient *ethclient.Client, target common.Address) *big.Int {
	balance, err := wsClient.BalanceAt(context.Background(), target, nil)
	require.NoError(t, err)
	return balance
}

func ntnBalanceOf(t *testing.T, autonity *autonity.Autonity, target common.Address) *big.Int {
	balance, err := autonity.BalanceOf(nil, target)
	require.NoError(t, err)
	return balance
}

func ntnSelfBonded(t *testing.T, autonity *autonity.Autonity, validator common.Address) *big.Int {
	val, err := autonity.GetValidator(nil, validator)
	require.NoError(t, err)
	return val.SelfBondedStake
}

const defaultEpochPeriod = 40
const defaultLookbackWindow = 10
const defaultInitialJailingPeriod = 30
const defaultInitialProbationPeriod = 2
const defaultDelta = 10
const defaultInactivityThreshold = 1000
const defaultInitialSlashingRate = 100
const defaultPastPerformanceWeight = 1000

func defaultGenesisOptions(genesis *core.Genesis) {
	genesis.Config.AutonityContractConfig.EpochPeriod = defaultEpochPeriod
	genesis.Config.OmissionAccountabilityConfig.LookbackWindow = defaultLookbackWindow
	genesis.Config.OmissionAccountabilityConfig.InitialJailingPeriod = defaultInitialJailingPeriod
	genesis.Config.OmissionAccountabilityConfig.InitialProbationPeriod = defaultInitialProbationPeriod
	genesis.Config.OmissionAccountabilityConfig.Delta = defaultDelta
	genesis.Config.OmissionAccountabilityConfig.InactivityThreshold = defaultInactivityThreshold
	genesis.Config.OmissionAccountabilityConfig.InitialSlashingRate = defaultInitialSlashingRate
	genesis.Config.OmissionAccountabilityConfig.PastPerformanceWeight = defaultPastPerformanceWeight
}

// assumes the usage of defaultEpochPeriod as epoch period
func waitForEpochEnd(t *testing.T, network e2e.Network, autonity *autonity.Autonity) {
	currentEpochID := epochID(t, autonity)

	err := network.WaitForHeight(defaultEpochPeriod*(currentEpochID+1), defaultEpochPeriod*2)
	require.NoError(t, err)

	newEpochID := epochID(t, autonity)
	require.Equal(t, currentEpochID+1, newEpochID)
}

func activateValidator(t *testing.T, network e2e.Network, autonity *autonity.Autonity, validatorID common.Address, treasuryKey *ecdsa.PrivateKey) {
	currentHeight := network[1].GetChainHeight()

	val := validator(t, autonity, validatorID)
	require.Equal(t, validatorID, val.NodeAddress)
	require.Equal(t, crypto.PubkeyToAddress(treasuryKey.PublicKey), val.Treasury)
	require.NotEqual(t, active, val.State)
	if val.State == jailed || val.State == jailedForInactivity {
		require.LessOrEqual(t, val.JailReleaseBlock.Uint64(), currentHeight)
	}
	require.NotEqual(t, jailbound, val.State)
	require.NotEqual(t, jailboundForInactivity, val.State)

	transactor, err := bind.NewKeyedTransactorWithChainID(treasuryKey, params.TestChainConfig.ChainID)
	require.NoError(t, err)
	transactor.NoSend = true
	tx, err := autonity.ActivateValidator(transactor, validatorID)
	require.NoError(t, err)
	err = network[1].WsClient.SendTransaction(context.Background(), tx)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = network[1].AwaitTransactions(ctx, tx)
	require.NoError(t, err)
	require.Equal(t, active, validator(t, autonity, validatorID).State)
}

// No omission fault, no validator is addressed as faulty.
func TestOmissionHappyCase(t *testing.T) {
	network := createNetwork(t, 6, true, defaultGenesisOptions)
	defer network.Shutdown(t)

	autonityContract, omissionContract := contracts(t, network[0])

	err := network.WaitForHeight((defaultEpochPeriod*2)+20, (defaultEpochPeriod*2)+20)
	require.NoError(t, err)

	// should have reached epoch 2
	require.Equal(t, uint64(2), epochID(t, autonityContract))

	// inactivity counters for current epoch should all be 0
	// scores compute in epoch 1 should be 0
	for _, n := range network {
		require.Equal(t, uint64(0), inactivityCounter(t, omissionContract, n.Address))
		require.Equal(t, uint64(0), inactivityScore(t, omissionContract, n.Address))
	}
}

// Let a single node be omission faulty, and then recover it.
func TestOmissionFaultyNode(t *testing.T) {
	network := createNetwork(t, 6, true, defaultGenesisOptions)
	defer network.Shutdown(t)

	autonityContract, omissionContract := contracts(t, network[1])

	// wait to reach the last delta blocks of the epoch, which are not accountable
	err := network.WaitForHeight(defaultEpochPeriod-defaultDelta+1, defaultEpochPeriod*2)
	require.NoError(t, err)

	// stop node 0
	faultyNode := network[0].Address
	err = network[0].Close(false)
	require.NoError(t, err)
	network[0].Wait()

	// close epoch 0
	waitForEpochEnd(t, network, autonityContract)

	// network should be up and continue to mine blocks, although with slower pace
	err = network.WaitForHeight((defaultEpochPeriod*2)-10, defaultEpochPeriod*2)
	require.NoError(t, err)

	// we should still be in the epoch 1
	require.Equal(t, uint64(1), epochID(t, autonityContract))

	// inactivity score of past epoch should be = 0
	// inactivity counter of current epoch should be > 0
	score := inactivityScore(t, omissionContract, faultyNode)
	counter := inactivityCounter(t, omissionContract, faultyNode)
	t.Logf("inactivity counter: %d", counter)
	t.Logf("inactivity score: %d", score)
	require.Equal(t, uint64(0), score)
	require.Greater(t, counter, uint64(0))

	// close epoch 1
	waitForEpochEnd(t, network, autonityContract)

	// inactivity score of past epoch should be > 0
	score = inactivityScore(t, omissionContract, faultyNode)
	t.Logf("inactivity score: %d", score)
	require.Greater(t, score, uint64(0))

	// faulty node should be jailed and outside the committee - no slashing for now
	faultyValidator := validator(t, autonityContract, faultyNode)
	require.Equal(t, jailedForInactivity, faultyValidator.State)
	require.Equal(t, uint64(defaultEpochPeriod*2+defaultInitialJailingPeriod), faultyValidator.JailReleaseBlock.Uint64())
	require.Equal(t, uint64(0), faultyValidator.TotalSlashed.Uint64())
	require.Equal(t, uint64(1), offences(t, omissionContract, faultyNode))
	require.Equal(t, uint64(1), collusionDegree(t, omissionContract, 1))
	require.Equal(t, uint64(defaultInitialProbationPeriod), probation(t, omissionContract, faultyNode))
	for _, member := range committee(t, autonityContract) {
		if member.Addr == faultyNode {
			t.Fatal("Omission faulty node is in the committee")
		}
		// active validator scores should be 0
		require.Equal(t, uint64(0), inactivityCounter(t, omissionContract, faultyNode))
	}
	// the faulty validator should be in the absent list for the entirety of epoch 1 (except for last delta blocks that are not accountable)
	for h := uint64(defaultEpochPeriod) + 1; h <= (defaultEpochPeriod*2)-defaultDelta; h++ {
		t.Logf("checking absent list for height: %d", h)
		require.True(t, isAbsent(t, omissionContract, h, faultyNode))
		// all other members should never be absent
		for _, member := range committee(t, autonityContract) {
			require.False(t, isAbsent(t, omissionContract, h, member.Addr))
		}
	}

	// we should be able to re-activate at next epoch
	require.Less(t, faultyValidator.JailReleaseBlock.Uint64(), uint64(defaultEpochPeriod*3))

	// restart faulty node, but we did not re-activate yet - the validator is outside the committee
	err = network[0].Start()
	require.NoError(t, err)

	// close epoch 2
	waitForEpochEnd(t, network, autonityContract)
	// close epoch 3
	waitForEpochEnd(t, network, autonityContract)

	// we should still be out of the committee and our probation period should have remained the same
	require.Equal(t, uint64(defaultInitialProbationPeriod), probation(t, omissionContract, faultyNode))
	for _, member := range committee(t, autonityContract) {
		if member.Addr == faultyNode {
			t.Fatal("Omission faulty node is in the committee")
		}
		// active validator scores should be 0
		require.Equal(t, uint64(0), inactivityCounter(t, omissionContract, faultyNode))
	}

	// re-activate validator
	t.Log("Re-activating validator 0")
	require.Equal(t, jailedForInactivity, validator(t, autonityContract, faultyNode).State) // should still be jailed
	activateValidator(t, network, autonityContract, faultyNode, network[0].TreasuryKey)

	// close epoch 4 - we should be back into the committee
	waitForEpochEnd(t, network, autonityContract)

	found := false
	for _, member := range committee(t, autonityContract) {
		if member.Addr == faultyNode {
			found = true
		}
		// active validator scores should be 0
		require.Equal(t, uint64(0), inactivityCounter(t, omissionContract, faultyNode))
	}
	require.True(t, found)

	// close epoch 5
	// - our probation period should reduce
	// - our inactivity score should reduce and should be lower than the threshold
	waitForEpochEnd(t, network, autonityContract)
	require.Equal(t, uint64(defaultInitialProbationPeriod)-1, probation(t, omissionContract, faultyNode))
	require.Less(t, inactivityScore(t, omissionContract, faultyNode), score)
	require.Less(t, inactivityScore(t, omissionContract, faultyNode), uint64(defaultInactivityThreshold))

	// stop node 0 again, it should get slashed at the end of this epoch now
	err = network[0].Close(false)
	require.NoError(t, err)
	network[0].Wait()

	// stop node 2 as well (node 1 is used for state retrieval calls)
	err = network[2].Close(false)
	require.NoError(t, err)
	network[2].Wait()

	// close epoch 6
	// mine some blocks to avoid hitting deadline exceed in `waitForEpochEnd`. This is because the mining rate is reduced due to the 2 offline vals
	err = network.WaitToMineNBlocks(20, 40, false)
	require.NoError(t, err)
	waitForEpochEnd(t, network, autonityContract)

	// inactivity score of past epoch should be > threshold
	score = inactivityScore(t, omissionContract, faultyNode)
	t.Logf("inactivity score of node 0: %d", score)
	require.Greater(t, score, uint64(defaultInactivityThreshold))
	score = inactivityScore(t, omissionContract, network[2].Address)
	t.Logf("inactivity score of node 2: %d", score)
	require.Greater(t, score, uint64(defaultInactivityThreshold))

	// faulty node should also have gotten slashed now
	faultyValidator = validator(t, autonityContract, faultyNode)
	require.Equal(t, jailedForInactivity, faultyValidator.State)
	require.Equal(t, uint64(2), offences(t, omissionContract, faultyNode))
	require.Equal(t, uint64(2), collusionDegree(t, omissionContract, 6))
	require.Equal(t, uint64(defaultInitialProbationPeriod*4)+1, probation(t, omissionContract, faultyNode))
	require.Equal(t, uint64(defaultEpochPeriod*7+defaultInitialJailingPeriod*4), faultyValidator.JailReleaseBlock.Uint64())
	require.NotEqual(t, uint64(0), faultyValidator.TotalSlashed.Uint64())

	// node 2 should have 1 offence
	require.Equal(t, uint64(1), offences(t, omissionContract, network[2].Address))

	// restart val2 and re-activate it
	err = network[2].Start()
	require.NoError(t, err)

	// we should be able to re-activate after epoch 7 is closed
	require.Less(t, validator(t, autonityContract, network[2].Address).JailReleaseBlock.Uint64(), uint64(defaultEpochPeriod*8))

	// close epoch 7
	waitForEpochEnd(t, network, autonityContract)

	// re-activate validator 2
	t.Log("Re-activating validator 2")
	require.Equal(t, jailedForInactivity, validator(t, autonityContract, network[2].Address).State) // should still be jailed
	activateValidator(t, network, autonityContract, network[2].Address, network[2].TreasuryKey)

	// close epoch 8 - val2 should be back into the committee
	waitForEpochEnd(t, network, autonityContract)

	// probation and #offences should still be un-modified
	require.Equal(t, uint64(defaultInitialProbationPeriod), probation(t, omissionContract, network[2].Address))
	require.Equal(t, uint64(1), offences(t, omissionContract, network[2].Address))

	// close epoch 9 - probation should reduce
	waitForEpochEnd(t, network, autonityContract)
	require.Equal(t, uint64(defaultInitialProbationPeriod)-1, probation(t, omissionContract, network[2].Address))

	// close epoch 9 - probation should have gotten to 0, and offence counter should have been reset as well
	waitForEpochEnd(t, network, autonityContract)
	require.Equal(t, uint64(0), probation(t, omissionContract, network[2].Address))
	require.Equal(t, uint64(0), offences(t, omissionContract, network[2].Address))
}

// Edge case: make 1/3 members be omission faulty out of the committee, make the net proposer effort and total net
// effort to zero during the period. The network should keep up and faulty nodes should be detected.
func TestOmissionFFaultyNodes(t *testing.T) {
	numOfNodes := 6
	numOfFaultyNodes := 2
	network := createNetwork(t, numOfNodes, false, defaultGenesisOptions)
	defer network.Shutdown(t)

	// start quorum nodes only, let F nodes be faulty from the very beginning.
	faultyNodes := make([]common.Address, numOfFaultyNodes)
	for i, n := range network {
		if i >= numOfFaultyNodes {
			err := n.Start()
			require.NoError(t, err)
		} else {
			faultyNodes[i] = n.Address
		}
	}

	// network should be up and continue to mine blocks
	err := network.WaitToMineNBlocks(defaultEpochPeriod*2, defaultEpochPeriod*4, false)
	require.NoError(t, err)

	_, omissionContract := contracts(t, network[numOfFaultyNodes])

	for _, n := range network[numOfFaultyNodes:] {
		require.Equal(t, uint64(0), inactivityScore(t, omissionContract, n.Address))
	}

	for _, n := range faultyNodes {
		score := inactivityScore(t, omissionContract, n)
		t.Logf("inactivity score of node %s: %d", n, score)
		require.Greater(t, score, uint64(0))
	}
}

// Edge case, a node keeps resetting for every 30 seconds.
func TestOmissionNodeKeepResetting(t *testing.T) {
	t.Run("Node that does not hold quorum keeps resetting", func(t *testing.T) {
		runResettingNodeTest(t, 5)
	})
	t.Run("Node that does hold quorum keeps resetting", func(t *testing.T) {
		// same as before but the offline node holds 1/3 of voting power (making the net proposer effort 0)
		runResettingNodeTest(t, 3)
	})
}

func runResettingNodeTest(t *testing.T, numNodes int) {
	// increase epoch period for this test, we want to test the behaviour of the system when a proposer is not able to provide an activity proof due to crashing
	network := createNetwork(t, numNodes, true, defaultGenesisOptions, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = 100
	})
	defer network.Shutdown(t)

	_, omissionContract := contracts(t, network[1])

	for network[1].Eth.BlockChain().CurrentHeader().Number.Uint64() <= 200 {
		err := network[0].Restart()
		require.NoError(t, err)
		time.Sleep(30 * time.Second)
	}

	score := inactivityScore(t, omissionContract, network[0].Address)
	t.Logf("inactivity score of node %s: %d", network[0].Address, score)
	require.Greater(t, score, uint64(0))
}

// Supposing that a node crashes by accident in an epoch and comes back up right away, ideally we would not want to punish him (or at least not too hard)
// this can be controlled by tweaking the lookback window. A reasonable value for look-back window depends on many factors:
// 1. The amount of chain data to be downloaded to get client synced again.
// 2. The HW ability of the client, computing power, memory, etc...
// 3. The latencies in between ACN and execution peers.
// this test can help us make an estimate for an appropriate value (however a real scenario will differ a lot)
func TestOmissionNodeHasOneResetWithinEpoch(t *testing.T) {
	t.Skip("Not important for now") //TODO(lorenzo) fix or remove
	network := createNetwork(t, 5, true, defaultGenesisOptions, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = uint64(300)
		genesis.Config.OmissionAccountabilityConfig.LookbackWindow = uint64(100) // lowering the lookback window will make the resetting node pass the threshold
	})
	defer network.Shutdown(t)

	resetHeight := uint64(60)
	for network[1].Eth.BlockChain().CurrentHeader().Number.Uint64() <= 300 {
		height := network[1].Eth.BlockChain().CurrentHeader().Number.Uint64()
		if height == resetHeight {
			err := network[0].Restart()
			require.NoError(t, err)
		}
		time.Sleep(1 * time.Second)
	}

	_, omissionContract := contracts(t, network[1])

	score := inactivityScore(t, omissionContract, network[0].Address)
	t.Logf("inactivity score of node %s: %d", network[0].Address, score)
	// TODO: another option is to allow score > 0 but verify that it is < InactivityThreshold
	require.Equal(t, uint64(0), score)
}

type lockedSlice struct {
	slice []uint64
	sync.RWMutex
}

func newNoActivityProposalSender(ls *lockedSlice) func(c interfaces.Core) interfaces.Proposer {
	return func(c interfaces.Core) interfaces.Proposer {
		return &noActivityProposalSender{ls, c.(*tdmcore.Core), c.Proposer()}
	}
}

type noActivityProposalSender struct {
	*lockedSlice
	*tdmcore.Core
	interfaces.Proposer
}

func (c *noActivityProposalSender) SendProposal(ctx context.Context, p *types.Block) {
	header := p.Header()

	var proposal *types.Block
	if header.ActivityProof != nil {
		c.Core.Logger().Warn("Scrubbing activity proof from proposal...")
		c.lockedSlice.Lock()
		defer c.lockedSlice.Unlock()

		// mark height as faulty
		c.lockedSlice.slice = append(c.lockedSlice.slice, header.Number.Uint64())

		// scrub activity proof, recompute state root and resign block
		header.ActivityProof = nil
		header.ActivityProofRound = 0
		backend, ok := c.Core.Backend().(*backend.Backend)
		if !ok {
			panic("cannot cast backend")
		}
		parent := backend.BlockChain().GetHeaderByHash(header.ParentHash)
		statedb, err := backend.BlockChain().StateAt(parent.Root)
		if err != nil {
			panic("Cannot fetch parent state: " + err.Error())
		}
		finalizedBlock, err := backend.FinalizeAndAssemble(c.Core.Backend().BlockChain(), header, statedb, []*types.Transaction{}, []*types.Header{}, &[]*types.Receipt{})
		if err != nil {
			panic("cannot re-finalize block: " + err.Error())
		}
		proposal, err = c.Core.Backend().AddSeal(finalizedBlock)
		if err != nil {
			panic("cannot re-sign block: " + err.Error())
		}
	} else {
		proposal = p
	}
	c.Proposer.SendProposal(ctx, proposal)
}

// a node is always online but never provides valid activity proofs
// he should still get a relatively high (depending on how often he is proposer) inactivity score
func TestOmissionProposerFaulty(t *testing.T) {
	validators, err := e2e.Validators(t, 4, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	faultyHeights := &lockedSlice{slice: make([]uint64, 0)}

	validators[0].TendermintServices = &interfaces.Services{Proposer: newNoActivityProposalSender(faultyHeights)}

	network, err := e2e.NewNetworkFromValidators(t, validators, true, defaultGenesisOptions)
	require.NoError(t, err)
	defer network.Shutdown(t)

	autonityContract, omissionContract := contracts(t, network[1])

	waitForEpochEnd(t, network, autonityContract)

	// this locking will stop the faulty proposer to propose new proposals with empty activity proof, therefore we can do our checks correctly
	faultyHeights.RLock()
	defer faultyHeights.RUnlock()

	t.Logf("faulty heights %v", faultyHeights.slice)

	// all faulty heights should be marked as faulty in the contract as well
	for _, height := range faultyHeights.slice {
		targetHeight := height - defaultDelta
		t.Logf("Checking height %d, targetHeight: %d", height, targetHeight)
		require.True(t, IsProposerFaulty(t, omissionContract, targetHeight))
	}

	expectedInactivityScore := computeInactivityScore(uint64(len(faultyHeights.slice)), scaleFactor(t, omissionContract))
	actualScore := inactivityScore(t, omissionContract, network[0].Address)
	t.Logf("expectedInactivityScore %v, inactivityScore %v", expectedInactivityScore, actualScore)
	require.Equal(t, expectedInactivityScore, actualScore)
}

type lockedBigInt struct {
	value *big.Int
	sync.RWMutex
}

func newEffortTrackerProposer(lb *lockedBigInt, epochPeriod uint64) func(c interfaces.Core) interfaces.Proposer {
	return func(c interfaces.Core) interfaces.Proposer {
		return &effortTrackerProposer{epochPeriod: epochPeriod, lockedBigInt: lb, Core: c.(*tdmcore.Core), Proposer: c.Proposer()}
	}
}

type effortTrackerProposer struct {
	epochPeriod uint64
	*lockedBigInt
	*tdmcore.Core
	interfaces.Proposer
}

func (c *effortTrackerProposer) SendProposal(ctx context.Context, proposal *types.Block) {
	// do not send proposal if we closed the first epoch
	if proposal.Header().Number.Uint64() > c.epochPeriod {
		return
	}
	backend, ok := c.Core.Backend().(*backend.Backend)
	if !ok {
		panic("cannot cast backend")
	}

	header := proposal.Header()
	epoch, err := backend.BlockChain().LatestEpoch()
	if err != nil {
		panic("cannot fetch latest epoch: " + err.Error())
	}
	committee := epoch.Committee
	quorum := bft.Quorum(committee.TotalVotingPower())

	if header.ActivityProof != nil {
		c.Core.Logger().Warn("Local node is the proposer, saving activity proof effort!")

		c.lockedBigInt.Lock()
		defer c.lockedBigInt.Unlock()

		effort := new(big.Int).Sub(header.ActivityProof.Signers.Power(), quorum)
		c.lockedBigInt.value = new(big.Int).Add(c.lockedBigInt.value, effort)
	}
	c.Proposer.SendProposal(ctx, proposal)
}

func TestOmissionRewards(t *testing.T) {
	// test that proposer effort is correctly compute and related rewards are distributed correctly
	t.Run("verify effort and rewards computations", func(t *testing.T) {
		runRewardTest(t, 4, 0)
	})
	t.Run("check rewards withholding", func(t *testing.T) {
		runRewardTest(t, 4, 1)
	})
}

func runRewardTest(t *testing.T, numNodes int, numOffline int) {
	isOnline := func(i int) bool {
		return i < (numNodes - numOffline)
	}

	validators, err := e2e.Validators(t, numNodes, "10e18,v,10000,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	customEpochPeriod := uint64(70)

	effortTrackers := make([]*lockedBigInt, numNodes)
	initialAtn := make([]*big.Int, numNodes)
	initialNtnSelfBonded := make([]*big.Int, numNodes)
	for i := 0; i < numNodes; i++ {
		effortTrackers[i] = &lockedBigInt{value: new(big.Int)}
		validators[i].TendermintServices = &interfaces.Services{Proposer: newEffortTrackerProposer(effortTrackers[i], customEpochPeriod)}
		initialAtn[i] = new(big.Int).Set(validators[i].InitialEth)
		initialNtnSelfBonded[i] = new(big.Int).SetUint64(validators[i].Stake)
	}

	autonityAtns := new(big.Int).SetInt64(54644123324456465) // random amount
	network, err := e2e.NewNetworkFromValidators(t, validators, false, defaultGenesisOptions, func(genesis *core.Genesis) {
		genesis.Config.AutonityContractConfig.EpochPeriod = customEpochPeriod
		genesis.Config.AutonityContractConfig.ProposerRewardRate = 10000          // all rewards are distributed based on proposer effort for ease of testing
		genesis.Config.AutonityContractConfig.MaxCommitteeSize = uint64(numNodes) // make the committee full to not have reward reduction due to committee factor
		genesis.Config.AutonityContractConfig.TreasuryFee = 0
		genesis.Alloc[params.AutonityContractAddress] = core.GenesisAccount{Balance: autonityAtns} // give some ATNs to AC for rewards
	})
	require.NoError(t, err)
	defer network.Shutdown(t)

	for i, node := range network {
		if isOnline(i) {
			err = node.Start()
			require.NoError(t, err)
		}
	}

	autonityContract, omissionContract := contracts(t, network[0])

	autonityTreasury, err := autonityContract.GetTreasuryAccount(nil)
	require.NoError(t, err)
	autonityTreasuryInitialAtn := atnBalanceOf(t, network[0].WsClient, autonityTreasury)
	autonityTreasuryInitialNtn := ntnBalanceOf(t, autonityContract, autonityTreasury)

	// stop before the epoch end, as proposer effort is reset at epoch end
	for network[0].Eth.BlockChain().CurrentHeader().Number.Uint64() < customEpochPeriod-20 {
		time.Sleep(100 * time.Millisecond)
	}

	// this locking will make the network lose liveness, therefore we can do our checks correctly
	for i := 0; i < numNodes; i++ {
		effortTrackers[i].RLock()
	}

	totalEffortTracked := new(big.Int)
	for i, node := range network {
		t.Logf("node %d, effort: %d", i, effortTrackers[i].value.Uint64())
		require.Equal(t, effortTrackers[i].value.String(), effort(t, omissionContract, node.Address).String())
		totalEffortTracked.Add(totalEffortTracked, effortTrackers[i].value)

		if isOnline(i) {
			require.Equal(t, uint64(0), inactivityCounter(t, omissionContract, node.Address))
		} else {
			require.NotEqual(t, uint64(0), inactivityCounter(t, omissionContract, node.Address))
		}
	}

	require.Equal(t, totalEffortTracked.String(), totalEffort(t, omissionContract).String())

	mintedStakeCh := make(chan *autonity.AutonityMintedStake)
	sub, err := autonityContract.WatchMintedStake(nil, mintedStakeCh, []common.Address{params.AutonityContractAddress})
	require.NoError(t, err)
	defer sub.Unsubscribe()

	// resume network liveness, close the epoch and verify reward computations
	for i := 0; i < numNodes; i++ {
		effortTrackers[i].RUnlock()
	}

	// close the epoch
	err = network.WaitForHeight(customEpochPeriod, int(customEpochPeriod))
	require.NoError(t, err)
	require.Equal(t, uint64(1), epochID(t, autonityContract))

	for i := 0; i < numNodes; i++ {
		effortTrackers[i].RLock()
	}

	ev := <-mintedStakeCh
	ntnRewards := ev.Amount

	totalEpochEffort := new(big.Int)
	for _, effortTracker := range effortTrackers {
		totalEpochEffort.Add(totalEpochEffort, effortTracker.value)
	}

	expectedWithheldRewardsAtn := new(big.Int)
	expectedWithheldRewardsNtn := new(big.Int)
	for i, effortTracker := range effortTrackers {
		atnExpectedIncrement := new(big.Int).Mul(effortTracker.value, autonityAtns)
		atnExpectedIncrement.Div(atnExpectedIncrement, totalEpochEffort)
		ntnExpectedIncrement := new(big.Int).Mul(effortTracker.value, ntnRewards)
		ntnExpectedIncrement.Div(ntnExpectedIncrement, totalEpochEffort)

		if isOnline(i) {
			atnExpectedBalance := new(big.Int).Add(initialAtn[i], atnExpectedIncrement)
			ntnExpectedBalance := new(big.Int).Add(initialNtnSelfBonded[i], ntnExpectedIncrement)
			treasury := crypto.PubkeyToAddress(network[i].TreasuryKey.PublicKey)
			require.Equal(t, atnExpectedBalance.String(), atnBalanceOf(t, network[0].WsClient, treasury).String())
			require.Equal(t, ntnExpectedBalance.String(), ntnSelfBonded(t, autonityContract, network[i].Address).String())
		} else {
			expectedWithheldRewardsAtn.Add(expectedWithheldRewardsAtn, atnExpectedIncrement)
			expectedWithheldRewardsNtn.Add(expectedWithheldRewardsNtn, ntnExpectedIncrement)
		}
	}

	expectedAutonityTreasuryBalanceAtn := new(big.Int).Add(autonityTreasuryInitialAtn, expectedWithheldRewardsAtn)
	expectedAutonityTreasuryBalanceNtn := new(big.Int).Add(autonityTreasuryInitialNtn, expectedWithheldRewardsNtn)
	require.Equal(t, expectedAutonityTreasuryBalanceAtn.String(), atnBalanceOf(t, network[0].WsClient, autonityTreasury).String())
	require.Equal(t, expectedAutonityTreasuryBalanceNtn.String(), ntnBalanceOf(t, autonityContract, autonityTreasury).String())

	for i := 0; i < numNodes; i++ {
		effortTrackers[i].RUnlock()
	}
}
