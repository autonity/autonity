package economic

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/shopspring/decimal"
	"github.com/wcharczuk/go-chart/v2"
	"math/big"
	"os"
	"time"
)

var (
	defATNPriceInUSD, _ = decimal.NewFromString("211.28") // to be replaced with production price
	defNTNPriceInUSD, _ = decimal.NewFromString("981.25") // to be replaced with production price
	defSysParams        = systemParams{
		EpochPeriod: 30,

		// market price of ATN, NTN that we are targeting to.
		AtnPriceTarget: defATNPriceInUSD,
		NtnPriceTarget: defNTNPriceInUSD,

		// default atn gas fee market settings.
		GenesisGasLimit:          20_000_000,
		GasCeil:                  20_000_000, // 20M
		InitialBaseFee:           new(big.Int).SetUint64(1_000_000_000),
		MinBaseFee:               new(big.Int).SetUint64(500_000_000),
		BaseFeeChangeDenominator: 8,
		ElasticityMultiplier:     2,

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
	EpochPeriod uint64 `json:"epochPeriod"`

	// the market price of ATN and NTN that we are targeting to at genesis.
	AtnPriceTarget decimal.Decimal `json:"atnPriceTarget"`
	NtnPriceTarget decimal.Decimal `json:"ntnPriceTarget"`

	// atn gas fee market settings
	GenesisGasLimit          uint64   `json:"genesisGasLimit"`
	GasCeil                  uint64   `json:"gasCeil"` // default value is 20_000_000
	InitialBaseFee           *big.Int `json:"initialBaseFee"`
	MinBaseFee               *big.Int `json:"minBaseFee"`
	BaseFeeChangeDenominator uint64   `json:"baseFeeChangeDenominator"`
	ElasticityMultiplier     uint64   `json:"elasticityMultiplier"`

	// ntn inflation settings
	InflationRateInitial      *big.Int `json:"inflationRateInitial"`
	InflationRateTransition   *big.Int `json:"inflationRateTransition"`
	InflationCurveConvexity   *big.Int `json:"inflationCurveConvexity"`
	InflationTransitionPeriod *big.Int `json:"inflationTransitionPeriod"`
	InflationReserveDecayRate *big.Int `json:"inflationReserveDecayRate"`
	InflationReserves         *big.Int `json:"inflationReserves"`
	NTNCirculatingSupply      *big.Int `json:"ntnCirculatingSupply"`
}

// WriteToFile writes the systemParams object to a JSON file
func (sp *systemParams) WriteToFile(filename string) error {
	// Convert the struct to JSON
	jsonData, err := json.MarshalIndent(sp, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
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

func amountToATNUSD(atn *big.Int) (*big.Float, *big.Float) {
	oneATNInWei := new(big.Float).SetUint64(1e18)
	// ATN to USD exchange rate
	atnToUSD := new(big.Float).SetFloat64(defATNPriceInUSD.InexactFloat64())
	// Convert accumulatedATNRewards to big.Float
	atnRewardsFloat := new(big.Float).SetInt(atn)
	// Calculate accumulatedATNRewards in ATN and USD
	atnRewardsInATN := new(big.Float).Quo(atnRewardsFloat, oneATNInWei)
	atnRewardsInUSD := new(big.Float).Mul(atnRewardsInATN, atnToUSD)
	return atnRewardsInATN, atnRewardsInUSD
}

func amountToNTNUSD(ntn *big.Int) (*big.Float, *big.Float) {
	ntnDecimals := new(big.Float).SetUint64(1e18)
	// Convert accumulatedNTNRewards to big.Float
	ntnRewardsFloat := new(big.Float).SetInt(ntn)
	// NTN to USD exchange rate
	ntnToUSD := new(big.Float).SetFloat64(defNTNPriceInUSD.InexactFloat64())
	// Calculate accumulatedNTNRewards in NTN and USD
	ntnRewardsInNTN := new(big.Float).Quo(ntnRewardsFloat, ntnDecimals)
	ntnRewardsInUSD := new(big.Float).Mul(ntnRewardsInNTN, ntnToUSD)
	return ntnRewardsInNTN, ntnRewardsInUSD
}

func (s *state) string() string {
	/*
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
	*/

	atnRewardsInATN, atnRewardsInUSD := amountToATNUSD(s.accumulatedATNRewards)
	ntnRewardsInNTN, ntnRewardsInUSD := amountToNTNUSD(s.accumulatedNTNRewards)

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
		gasLimit:  s.params.GenesisGasLimit,
		baseFee:   s.params.InitialBaseFee,
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

		if i != 0 && i%s.params.EpochPeriod == 0 {
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
}

func (s *simulator) assembleData(index int) {
	// get the metrics, each returned column is a slice of float64, each variable name below represent the colum name.
	blocks, gasLimits, gasUseds, baseFeesATN, baseFeesUSD, baseFeeChgRates, spamCosts, accFeeRewards, accFeeRewardsUSD,
		accNtnRewards, accNtnRewardsUSD, mergedRewards := s.Metrics()
	// Dump the metrics into a CSV file
	err := writeMetricsToCSV(index, blocks, gasLimits, gasUseds, baseFeesATN, baseFeesUSD, baseFeeChgRates, spamCosts,
		accFeeRewards, accFeeRewardsUSD, accNtnRewards, accNtnRewardsUSD, mergedRewards)
	if err != nil {
		log.Crit("Failed to write metrics to CSV: %v", err)
	}

	// below, we try to render the data with go chart.
	// render gas limit and gas used in YAxis and block in XAxis.
	gasChart := chart.Chart{
		Title: "GasLimits and GasUsed over blocks",
		XAxis: chart.XAxis{
			Name: "Block/Time (s)",
		},
		YAxis: chart.YAxis{
			Name: "Gas",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "GasLimit",
				XValues: blocks,
				YValues: gasLimits,
			},
			chart.ContinuousSeries{
				Name:    "GasUsed",
				XValues: blocks,
				YValues: gasUseds,
			},
		},
	}
	renderChartToFile(gasChart, fmt.Sprintf("%d_gas_chart.png", index))
	// render baseFee(ATN) baseFee(USD) in YAxis and block in XAxis.
	baseFeeChart := chart.Chart{
		Title: "Base Fee (ATN & USD) over blocks",
		XAxis: chart.XAxis{
			Name: "Block/Time (s)",
		},
		YAxis: chart.YAxis{
			Name: "BaseFee",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "baseFee(ATN)",
				XValues: blocks,
				YValues: baseFeesATN,
			},
			chart.ContinuousSeries{
				Name:    "baseFee(USD)",
				XValues: blocks,
				YValues: baseFeesUSD,
			},
		},
	}
	renderChartToFile(baseFeeChart, fmt.Sprintf("%d_base_fee_chart.png", index))
	// render baseFeeChangeRate in YAxis and block in XAxis.
	baseFeeChgRateChart := chart.Chart{
		Title: "Base Fee Change Rate over blocks",
		XAxis: chart.XAxis{
			Name: "Block/Time (s)",
		},
		YAxis: chart.YAxis{
			Name: "BaseFeeChangeRate",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{Name: "baseFeeChangeRate", XValues: blocks, YValues: baseFeeChgRates},
		},
	}
	renderChartToFile(baseFeeChgRateChart, fmt.Sprintf("%d_base_fee_chg_rate.png", index))
	// render block spam cost(USD) in YAxis and block in XAxis.
	spamCostsChart := chart.Chart{
		Title:  "SPAM Costs (USD) Over blocks",
		XAxis:  chart.XAxis{Name: "Block/Time (s)"},
		YAxis:  chart.YAxis{Name: "spam cost (USD)"},
		Series: []chart.Series{chart.ContinuousSeries{Name: "spamCosts", XValues: blocks, YValues: spamCosts}},
	}
	renderChartToFile(spamCostsChart, fmt.Sprintf("%d_spam_costs.png", index))

	// Merge the 3 into 1 chart?
	// render accumulating fee rewards(ATN) & (USD) in YAxis and block in XAxis.
	// render accumulating ntn rewards(NTN) & (USD) in YAxis and block in XAxis.
	// render accumulating merged rewards(USD) in YAxis and block in XAxis.
	rewardChart := chart.Chart{
		Title: "ATN & NTN Rewards Over blocks",
		XAxis: chart.XAxis{Name: "Block/Time (s)"},
		YAxis: chart.YAxis{Name: "FeeRewards & NTNRewards & Merged (USD)"},
		Series: []chart.Series{
			//chart.ContinuousSeries{Name: "accumulating fee rewards (ATN)", XValues: blocks, YValues: accFeeRewards},
			chart.ContinuousSeries{Name: "accumulating fee rewards (USD)", XValues: blocks, YValues: accFeeRewardsUSD},
			//chart.ContinuousSeries{Name: "accumulating ntn rewards (NTN)", XValues: blocks, YValues: accNtnRewards},
			chart.ContinuousSeries{Name: "accumulating ntn rewards (USD)", XValues: blocks, YValues: accNtnRewardsUSD},
			chart.ContinuousSeries{Name: "merged rewards (USD)", XValues: blocks, YValues: mergedRewards},
		},
	}
	renderChartToFile(rewardChart, fmt.Sprintf("%d_rewards.png", index))
}

func (s *simulator) Metrics() (heights []float64, gasLimits []float64, gasUseds []float64, baseFeesATN []float64,
	baseFeesUSD []float64, baseFeeChangeRates []float64, spamCosts []float64, accumulatingFeeRewards []float64,
	accumulatingFeeRewardsUSD []float64, accumulatingNtnRewards []float64, accumulatingNtnRewardsUSD []float64,
	mergedRewardsUSD []float64) {

	for i, b := range s.blocks {
		heights = append(heights, float64(b.number))
		gasLimits = append(gasLimits, float64(b.gasLimit))
		gasUseds = append(gasUseds, float64(b.gasUsed))
		baseFeeInATN, _ := b.data.baseFeeInATN.Float64()
		baseFeesATN = append(baseFeesATN, baseFeeInATN)
		baseFeeInUSD, _ := b.data.baseFeeInUSD.Float64()
		baseFeesUSD = append(baseFeesUSD, baseFeeInUSD)
		baseFeeChangeRates = append(baseFeeChangeRates, b.data.baseFeeChangeRate.InexactFloat64())
		spamCost, _ := b.data.blockSpamCostInUSD.Float64()
		spamCosts = append(spamCosts, spamCost)

		// accumulating rewards for atn, and converted in usd.
		accFeeRewardATN, accFeeRewardUSD := amountToATNUSD(s.states[i].accumulatedATNRewards)
		accFeeRewardATNFloat, _ := accFeeRewardATN.Float64()
		accumulatingFeeRewards = append(accumulatingFeeRewards, accFeeRewardATNFloat)
		accFeeRewardUSDFloat, _ := accFeeRewardUSD.Float64()
		accumulatingFeeRewardsUSD = append(accumulatingFeeRewardsUSD, accFeeRewardUSDFloat)

		// accumulating rewards for ntn, and converted in usd.
		accNtnRewardNTN, accNtnRewardUSD := amountToNTNUSD(s.states[i].accumulatedNTNRewards)
		accFeeRewardNTNFloat, _ := accNtnRewardNTN.Float64()
		accumulatingNtnRewards = append(accumulatingNtnRewards, accFeeRewardNTNFloat)
		accNtnRewardUSDFloat, _ := accNtnRewardUSD.Float64()
		accumulatingNtnRewardsUSD = append(accumulatingNtnRewardsUSD, accNtnRewardUSDFloat)

		mergedRewardsUSD = append(mergedRewardsUSD, accFeeRewardUSDFloat+accNtnRewardUSDFloat)
	}
	return heights, gasLimits, gasUseds, baseFeesATN, baseFeesUSD, baseFeeChangeRates, spamCosts, accumulatingFeeRewards,
		accumulatingFeeRewardsUSD, accumulatingNtnRewards, accumulatingNtnRewardsUSD, mergedRewardsUSD
}

func fillBlock(parent *block, params *systemParams, f TXNPacker) (*block, *big.Int) {
	baseFee := CalcBaseFee(parent, params)
	baseFeeChgRate := baseFeeChangeRate(parent.baseFee, baseFee)
	gasLimit := core.CalcGasLimit(parent.gasLimit, params.GasCeil)
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

// Render the chart to a file
func renderChartToFile(c chart.Chart, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Crit("Failed to create file: %v", err)
	}
	defer f.Close()

	err = c.Render(chart.PNG, f)
	if err != nil {
		log.Crit("Failed to render chart: %v", err)
	}
}

// writeMetricsToCSV writes the metrics to a CSV file
func writeMetricsToCSV(index int, blocks, gasLimits, gasUseds, baseFeesATN, baseFeesUSD, baseFeeChgRates, spamCosts,
	accFeeRewards, accFeeRewardsUSD, accNtnRewards, accNtnRewardsUSD, mergedRewards []float64) error {
	// Create the CSV file
	filename := fmt.Sprintf("%d_metrics.csv", index)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{
		"Block", "GasLimit", "GasUsed", "BaseFee(ATN)", "BaseFee(USD)", "BaseFeeChangeRate", "SpamCosts(USD)",
		"AccFeeRewards(ATN)", "AccFeeRewards(USD)", "AccNtnRewards(NTN)", "AccNtnRewards(USD)", "MergedRewards(USD)",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write the data rows
	for i := 0; i < len(blocks); i++ {
		row := []string{
			fmt.Sprintf("%f", blocks[i]),
			fmt.Sprintf("%f", gasLimits[i]),
			fmt.Sprintf("%f", gasUseds[i]),
			fmt.Sprintf("%f", baseFeesATN[i]),
			fmt.Sprintf("%f", baseFeesUSD[i]),
			fmt.Sprintf("%f", baseFeeChgRates[i]),
			fmt.Sprintf("%f", spamCosts[i]),
			fmt.Sprintf("%f", accFeeRewards[i]),
			fmt.Sprintf("%f", accFeeRewardsUSD[i]),
			fmt.Sprintf("%f", accNtnRewards[i]),
			fmt.Sprintf("%f", accNtnRewardsUSD[i]),
			fmt.Sprintf("%f", mergedRewards[i]),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}
