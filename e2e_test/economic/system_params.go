package economic

import (
	"fmt"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

var (
	defATNPriceInUSD, _ = decimal.NewFromString("211.28") // to be replaced with production price
	defNTNPriceInUSD, _ = decimal.NewFromString("985.25") // to be replaced with production price
	defSysParams        = systemParams{
		epochPeriod: 30,

		// market price of ATN, NTN that we are targeting to.
		atnPriceTarget: defATNPriceInUSD,
		ntnPriceTarget: defNTNPriceInUSD,

		// default atn gas fee market settings.
		genesisGasLimit:          20_000_000,
		gasCeil:                  20_000_000, // 20M
		initialBaseFee:           new(big.Int).SetUint64(1_000_000_000),
		minBaseFee:               new(big.Int).SetUint64(500_000_000),
		baseFeeChangeDenominator: 8,
		elasticityMultiplier:     2,

		// default ntn inflation settings.
		InflationRateInitial:      (*big.Int)(params.DefaultInflationControllerGenesis.InflationRateInitial),
		InflationRateTransition:   (*big.Int)(params.DefaultInflationControllerGenesis.InflationRateTransition),
		InflationCurveConvexity:   (*big.Int)(params.DefaultInflationControllerGenesis.InflationCurveConvexity),
		InflationTransitionPeriod: (*big.Int)(params.DefaultInflationControllerGenesis.InflationTransitionPeriod),
		InflationReserveDecayRate: (*big.Int)(params.DefaultInflationControllerGenesis.InflationReserveDecayRate),
		InflationReserves:         (*big.Int)(params.TestAutonityContractConfig.InitialInflationReserve),
		NTNCirculatingSupply:      new(big.Int).Mul(big.NewInt(60_000_000), params.NTNDecimalFactor), // NTN precision is 18
	}
)

type systemParams struct {
	epochPeriod uint64

	// the market price of ATN and NTN that we are targeting to at genesis.
	atnPriceTarget decimal.Decimal
	ntnPriceTarget decimal.Decimal

	// atn gas fee market settings
	genesisGasLimit          uint64
	gasCeil                  uint64 // default value is 20_000_000
	initialBaseFee           *big.Int
	minBaseFee               *big.Int
	baseFeeChangeDenominator uint64
	elasticityMultiplier     uint64

	// ntn inflation settings
	InflationRateInitial      *big.Int
	InflationRateTransition   *big.Int
	InflationCurveConvexity   *big.Int
	InflationTransitionPeriod *big.Int
	InflationReserveDecayRate *big.Int
	InflationReserves         *big.Int
	NTNCirculatingSupply      *big.Int
}

type block struct {
	timestamp         int64
	number            uint64
	numOfTXNs         uint64
	gasLimit          uint64
	gasUsed           uint64
	baseFee           *big.Int
	baseFeeChangeRate decimal.Decimal
	data              *blockData
}

type blockData struct {
	H                  uint64
	gasLimit           uint64
	gasUsed            uint64
	baseFeeChangeRate  decimal.Decimal
	baseFee            *big.Int
	baseFeeInATN       *big.Float
	baseFeeInUSD       *big.Float
	blockSpamCostInUSD *big.Float
}

func (b *block) collectData() {
	// Convert 1e18 to a big.Float for floating-point division
	oneATNInWei := new(big.Float).SetUint64(1e18)

	// Convert baseFee to big.Float
	baseFeeFloat := new(big.Float).SetInt(b.baseFee)

	// Calculate baseFee in ATN with decimals
	baseFeeInATN := new(big.Float).Quo(baseFeeFloat, oneATNInWei)

	// ATN to USD exchange rate (1 ATN = 1.28 USD)
	price := defATNPriceInUSD.InexactFloat64()
	atnToUSD := new(big.Float).SetFloat64(price)

	// Calculate baseFee in USD
	baseFeeInUSD := new(big.Float).Mul(baseFeeInATN, atnToUSD)

	// Block spam cost.
	gasLimit := new(big.Float).SetUint64(b.gasLimit)
	blockSpamCostInUSD := new(big.Float).Mul(baseFeeInUSD, gasLimit)

	b.data = &blockData{
		H:                  b.number,
		gasLimit:           b.gasLimit,
		gasUsed:            b.gasUsed,
		baseFeeChangeRate:  b.baseFeeChangeRate,
		baseFee:            b.baseFee,
		baseFeeInATN:       baseFeeInATN,
		baseFeeInUSD:       baseFeeInUSD,
		blockSpamCostInUSD: blockSpamCostInUSD,
	}
}

func (b *block) string() string {
	b.collectData()
	return fmt.Sprintf(
		"H: %d, GL: %d, GU: %d, baseFeeChgRate: %s%%, baseFee(Wei): %s, baseFee(ATN): %.9f, baseFee(USD): %.9f, blockSpamCost(USD): %.9f",
		b.data.H, b.data.gasLimit, b.data.gasUsed, b.data.baseFeeChangeRate.String(), b.data.baseFee.String(), b.data.baseFeeInATN, b.data.baseFeeInUSD, b.data.blockSpamCostInUSD,
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
		"feeRWD: %s (Wei), %.9f (ATN), %.9f (USD) "+
			"inflationRWD: %s, %.9f (NTN), %.9f (USD)",
		s.accumulatedATNRewards.String(), atnRewardsInATN, atnRewardsInUSD,
		s.accumulatedNTNRewards.String(), ntnRewardsInNTN, ntnRewardsInUSD,
	)
}

type simulator struct {
	name        string
	genesisTime int64
	numOfBlocks uint64
	params      *systemParams
	txnPacker   TXNPacker

	inflationEngine inflationEngine

	blocks []*block
	states []*state
}

func newSimulator(name string, blocks uint64, sp *systemParams, filler TXNPacker) *simulator {
	genesisTime := time.Now().Unix()

	inflationParams := &autonity.InflationControllerParams{
		InflationRateInitial:      sp.InflationRateInitial,
		InflationRateTransition:   sp.InflationRateTransition,
		InflationCurveConvexity:   sp.InflationCurveConvexity,
		InflationTransitionPeriod: sp.InflationTransitionPeriod,
		InflationReserveDecayRate: sp.InflationReserveDecayRate,
	}

	inflationCore := newInflationEngine(inflationParams, new(big.Int).SetInt64(genesisTime))

	return &simulator{
		name,
		genesisTime,
		blocks,
		sp,
		filler,
		inflationCore,
		nil,
		nil}
}

func (s *simulator) genesisBlock() *block {
	return &block{
		timestamp: s.genesisTime,
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

	lastState := &state{
		accumulatedATNRewards: new(big.Int).SetUint64(0),
		accumulatedNTNRewards: new(big.Int).SetUint64(0),
	}
	s.states = append(s.states, lastState)

	preBlock := genesisBlock
	log.Info("", "genesis", genesisBlock.string())

	// todo: add the simulation of dynamic adjustment of circulatingSupply.
	circulatingSupply := new(big.Int).Set(s.params.NTNCirculatingSupply)
	lastEpochTime := new(big.Int).SetInt64(genesisBlock.timestamp)
	inflationReserve := new(big.Int).Set(s.params.InflationReserves)
	for i := uint64(0); i < s.numOfBlocks; i++ {
		b, feeReward := fillBlock(preBlock, s.params, s.txnPacker)

		newState := &state{
			accumulatedATNRewards: new(big.Int).Add(lastState.accumulatedATNRewards, feeReward),
			accumulatedNTNRewards: new(big.Int).Set(lastState.accumulatedNTNRewards),
		}

		if i != 0 && i%s.params.epochPeriod == 0 {
			currentTime := new(big.Int).SetInt64(b.timestamp)
			inflationReward := s.inflationEngine.calculateSupplyDelta(circulatingSupply, inflationReserve, lastEpochTime, currentTime)
			lastEpochTime = currentTime
			newState.accumulatedNTNRewards = new(big.Int).Add(lastState.accumulatedNTNRewards, inflationReward)
			inflationReserve.Sub(inflationReserve, inflationReward)
			circulatingSupply.Add(circulatingSupply, inflationReward)
		}

		preBlock = b
		s.blocks = append(s.blocks, b)
		s.states = append(s.states, newState)
		lastState = newState

		log.Info("", "block data", b.string(), "accumulating rewards", newState.string())
	}

	// render data charts for analysis.
	s.renderDataCharts()
}

func (s *simulator) renderDataCharts() {
	// render gas limit and gas used in YAxis and block in XAxis.

	// render baseFee(ATN) baseFee(USD) in YAxis and block in XAxis.

	// render baseFeeChangeRate in YAxis and block in XAxis.

	// render block spam cost(USD) in YAxis and block in XAxis.

	// Merge the 3 into 1 chart?
	// render accumulating fee rewards(ATN) & (USD) in YAxis and block in XAxis.
	// render accumulating ntn rewards(NTN) & (USD) in YAxis and block in XAxis.
	// render accumulating merged rewards(USD) in YAxis and block in XAxis.
}

func fillBlock(parent *block, params *systemParams, f TXNPacker) (*block, *big.Int) {
	baseFee := CalcBaseFee(parent, params)
	baseFeeChgRate := baseFeeChangeRate(parent.baseFee, baseFee)
	gasLimit := core.CalcGasLimit(parent.gasLimit, params.gasCeil)
	numOfTXN, gasUsed := f.packTXNs(parent)
	atnRewards := new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), baseFee)
	b := &block{
		timestamp:         parent.timestamp + 1,
		number:            parent.number + 1,
		gasLimit:          gasLimit,
		baseFee:           baseFee,
		baseFeeChangeRate: baseFeeChgRate,
		numOfTXNs:         numOfTXN,
		gasUsed:           gasUsed,
	}
	return b, atnRewards
}
