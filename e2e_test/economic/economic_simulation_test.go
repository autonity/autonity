package economic

import (
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/log"
	"github.com/shopspring/decimal"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

var (
	defATNPriceInUSD, _ = decimal.NewFromString("1.28")
	defNTNPriceInUSD, _ = decimal.NewFromString("1.25")
	defSysParams        = systemParams{
		genesisGasLimit:          30_000_000,
		gasCeil:                  20_000_000,
		blockGasTarget:           20_000_000,
		initialBaseFee:           new(big.Int).SetUint64(1_000_000_000),
		minBaseFee:               new(big.Int).SetUint64(500_000_000),
		baseFeeChangeDenominator: 8,
		elasticityMultiplier:     2,
		atnPriceTarget:           defATNPriceInUSD,
		ntnPriceTarget:           defNTNPriceInUSD,
	}
)

type systemParams struct {
	genesisGasLimit          uint64
	gasCeil                  uint64 // default value is 20_000_000
	blockGasTarget           uint64
	initialBaseFee           *big.Int
	minBaseFee               *big.Int
	baseFeeChangeDenominator uint64
	elasticityMultiplier     uint64
	atnPriceTarget           decimal.Decimal
	ntnPriceTarget           decimal.Decimal
}

type block struct {
	number    uint64
	numOfTXNs uint64
	gasLimit  uint64
	gasUsed   uint64
	baseFee   *big.Int
}

func (b *block) string() string {
	// Convert 1e18 to a big.Float for floating-point division
	oneATNInWei := new(big.Float).SetUint64(1e18)

	// Convert baseFee to big.Float
	baseFeeFloat := new(big.Float).SetInt(b.baseFee)

	// Calculate baseFee in ATN with decimals
	baseFeeInETH := new(big.Float).Quo(baseFeeFloat, oneATNInWei)

	// ATN to USD exchange rate (1 ATN = 1.28 USD)
	price := defATNPriceInUSD.InexactFloat64()
	atnToUSD := new(big.Float).SetFloat64(price)

	// Calculate baseFee in USD
	baseFeeInUSD := new(big.Float).Mul(baseFeeInETH, atnToUSD)

	return fmt.Sprintf(
		"number: %d, gasLimit: %d, gasUsed: %d, baseFee(Wei): %s, baseFee(ATN): %.9f, baseFee(USD): %.9f",
		b.number, b.gasLimit, b.gasUsed, b.baseFee.String(), baseFeeInETH, baseFeeInUSD,
	)
}

type state struct {
	accumulatedATNRewards *big.Int // fee rewards.
	accumulatedNTNRewards *big.Int // inflation rewards.
}

func (s *state) string() string {
	// Convert 1e18 to a big.Float for floating-point division
	oneATNInWei := new(big.Float).SetUint64(1e18)
	ntnDecimals := new(big.Float).SetUint64(1e18)

	// ATN to USD exchange rate (1 ATN = 1.28 USD)
	atnToUSD := new(big.Float).SetFloat64(defATNPriceInUSD.InexactFloat64())

	// Convert accumulatedATNRewards to big.Float
	atnRewardsFloat := new(big.Float).SetInt(s.accumulatedATNRewards)

	// Calculate accumulatedATNRewards in ATN and USD
	atnRewardsInATN := new(big.Float).Quo(atnRewardsFloat, oneATNInWei)
	atnRewardsInUSD := new(big.Float).Mul(atnRewardsInATN, atnToUSD)

	// Convert accumulatedNTNRewards to big.Float
	ntnRewardsFloat := new(big.Float).SetInt(s.accumulatedNTNRewards)

	// NTN to USD exchange rate (1 NTN = 1.25 USD)
	ntnToUSD := new(big.Float).SetFloat64(defNTNPriceInUSD.InexactFloat64())

	// Calculate accumulatedNTNRewards in NTN and USD
	ntnRewardsInNTN := new(big.Float).Quo(ntnRewardsFloat, ntnDecimals)
	ntnRewardsInUSD := new(big.Float).Mul(ntnRewardsInNTN, ntnToUSD)

	return fmt.Sprintf(
		"accumulatedATNRewards(Wei): %s, accumulatedATNRewards(ATN): %.9f, accumulatedATNRewards(USD): %.9f"+
			"accumulatedNTNRewards: %s, accumulatedNTNRewards(NTN): %.9f, accumulatedNTNRewards(USD): %.9f",
		s.accumulatedATNRewards.String(), atnRewardsInATN, atnRewardsInUSD,
		s.accumulatedNTNRewards.String(), ntnRewardsInNTN, ntnRewardsInUSD,
	)
}

// CalcBaseFee calculates the basefee of the header.
func CalcBaseFee(parent *block, params *systemParams) *big.Int {
	var (
		parentGasTarget          = parent.gasLimit / params.elasticityMultiplier
		parentGasTargetBig       = new(big.Int).SetUint64(parentGasTarget)
		baseFeeChangeDenominator = new(big.Int).SetUint64(params.baseFeeChangeDenominator)
	)
	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.gasUsed == parentGasTarget {
		return new(big.Int).Set(parent.baseFee)
	}
	if parent.gasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		gasUsedDelta := new(big.Int).SetUint64(parent.gasUsed - parentGasTarget)
		x := new(big.Int).Mul(parent.baseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := math.BigMax(
			x.Div(y, baseFeeChangeDenominator),
			common.Big1,
		)
		return x.Add(parent.baseFee, baseFeeDelta)
	} else {
		// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
		gasUsedDelta := new(big.Int).SetUint64(parentGasTarget - parent.gasUsed)
		x := new(big.Int).Mul(parent.baseFee, gasUsedDelta)
		y := x.Div(x, parentGasTargetBig)
		baseFeeDelta := x.Div(y, baseFeeChangeDenominator)

		minBaseFee := big.NewInt(0).SetUint64(params.minBaseFee.Uint64())
		return math.BigMax(
			x.Sub(parent.baseFee, baseFeeDelta),
			minBaseFee,
		)
	}
}

type blockFiller interface {
	fillBlock(parent *block) (*block, *big.Int, *big.Int)
}

type fullFiller struct {
	params *systemParams
}

func (f *fullFiller) genTXNs(parent *block) (uint64, uint64) {
	// Define the gas cost per transaction
	const gasPerTxn = 21000

	// Generate a random number of additional transactions (1 to 5)
	rand.Seed(time.Now().UnixNano())
	additionalTxns := rand.Intn(5) + 1

	// Calculate the minimum gas required to exceed the parent's gas limit
	minGasRequired := parent.gasLimit + uint64(additionalTxns*gasPerTxn)

	// Calculate the number of transactions needed to exceed the gas limit
	numTxns := minGasRequired / gasPerTxn

	// Calculate the total gas used
	usedGas := numTxns * gasPerTxn

	return numTxns, usedGas
}

func (f *fullFiller) fillBlock(parent *block) (*block, *big.Int, *big.Int) {
	baseFee := CalcBaseFee(parent, f.params)
	numOfTXN, gasUsed := f.genTXNs(parent)

	atnRewards := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), baseFee)
	inflationRewards := new(big.Int).SetUint64(0)

	return &block{
		number:    parent.number + 1,
		gasUsed:   gasUsed,
		numOfTXNs: numOfTXN,
		gasLimit:  core.CalcGasLimit(parent.gasLimit, f.params.gasCeil),
		baseFee:   baseFee,
	}, atnRewards, inflationRewards
}

type simulator struct {
	numOfBlocks uint64
	coreState   *state
	params      *systemParams
	blockFiler  blockFiller
	blocks      []*block
}

func newSimulator(blocks uint64, params *systemParams, filler blockFiller) *simulator {
	coreState := &state{
		new(big.Int).SetUint64(0),
		new(big.Int).SetUint64(0),
	}
	return &simulator{blocks, coreState, params, filler, nil}
}

func (s *simulator) genesisBlock() *block {
	return &block{
		number:    0,
		numOfTXNs: 0,
		gasUsed:   0,
		gasLimit:  s.params.genesisGasLimit,
		baseFee:   s.params.initialBaseFee,
	}
}

// start the simulation, and collect runtime data.
func (s *simulator) start() {
	genesisBlock := s.genesisBlock()
	s.blocks = append(s.blocks, genesisBlock)
	preBlock := genesisBlock
	log.Info("genesis", "data", genesisBlock.string())

	for i := uint64(0); i < s.numOfBlocks; i++ {
		b, feeReward, inflationReward := s.blockFiler.fillBlock(preBlock)
		preBlock = b
		s.coreState.accumulatedATNRewards = s.coreState.accumulatedATNRewards.Add(s.coreState.accumulatedATNRewards, feeReward)
		s.coreState.accumulatedNTNRewards = s.coreState.accumulatedNTNRewards.Add(s.coreState.accumulatedNTNRewards, inflationReward)
		s.blocks = append(s.blocks, b)

		log.Info("block", "data", b.string(), "rewards", s.coreState.string())
	}

	return
}

func TestFullBlocks(t *testing.T) {
	blocks := uint64(900)
	systemParam := &defSysParams
	sim := newSimulator(blocks, systemParam, &fullFiller{params: systemParam})
	sim.start()
}
