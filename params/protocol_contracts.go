package params

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/influxdata/influxdb/pkg/deep"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params/generated"
)

var (
	DecimalPrecision = int64(18)
	SecondsInYear    = int64(365 * 24 * 60 * 60)
	SecondsInDay     = int64(24 * 60 * 60)
	DecimalFactor    = new(big.Int).Exp(big.NewInt(10), big.NewInt(DecimalPrecision), nil)
	NTNDecimalFactor = new(big.Int).SetUint64(Ether)

	//Oracle Contract defaults
	OracleVotePeriod           = uint64(30)
	OracleInitialSymbols       = []string{"AUD-USD", "CAD-USD", "EUR-USD", "GBP-USD", "JPY-USD", "SEK-USD", "ATN-USD", "NTN-USD", "NTN-ATN"}
	DefaultGenesisOracleConfig = &OracleContractGenesis{
		Symbols:                   OracleInitialSymbols,
		VotePeriod:                OracleVotePeriod,
		OutlierDetectionThreshold: 10,  // 10%
		OutlierSlashingThreshold:  225, // 15%
		BaseSlashingRate:          10,
	}

	// DefaultAcuContractGenesis contains the default values for the ASM ACU contract
	DefaultAcuContractGenesis = &AcuContractGenesis{
		Symbols:    []string{"AUD-USD", "CAD-USD", "EUR-USD", "GBP-USD", "JPY-USD", "SEK-USD", "USD-USD"},
		Quantities: []uint64{1_744_583, 1_598_986, 1_058_522, 886_091, 175_605_573, 12_318_802, 1_148_285},
		Scale:      uint64(7),
	}

	// DefaultStabilizationGenesis contains the default values for the ASM Stabilization contract
	DefaultStabilizationGenesis = &StabilizationContractGenesis{
		BorrowInterestRate:        (*math.HexOrDecimal256)(math.MustParseBig256("50_000_000_000_000_000")),
		LiquidationRatio:          (*math.HexOrDecimal256)(math.MustParseBig256("1_800_000_000_000_000_000")),
		MinCollateralizationRatio: (*math.HexOrDecimal256)(math.MustParseBig256("2_000_000_000_000_000_000")),
		MinDebtRequirement:        (*math.HexOrDecimal256)(math.MustParseBig256("1_000_000")),
		TargetPrice:               (*math.HexOrDecimal256)(math.MustParseBig256("1_618_034_000_000_000_000")),
	}

	DefaultSupplyControlGenesis = &SupplyControlGenesis{
		InitialAllocation: (*math.HexOrDecimal256)(new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), common.Big1)),
	}

	DefaultInflationControllerGenesis = &InflationControllerGenesis{
		InflationRateInitial:      (*math.HexOrDecimal256)(new(big.Int).Div(new(big.Int).Mul(big.NewInt(75), DecimalFactor), big.NewInt(1000*SecondsInYear))),        // 7.5% AR
		InflationRateTransition:   (*math.HexOrDecimal256)(new(big.Int).Div(new(big.Int).Mul(big.NewInt(55), DecimalFactor), big.NewInt(1000*SecondsInYear))),        // 5.5% AR
		InflationReserveDecayRate: (*math.HexOrDecimal256)(new(big.Int).Div(new(big.Int).Mul(big.NewInt(17_317), DecimalFactor), big.NewInt(100_000*SecondsInYear))), // 17.3168186793% AR
		InflationTransitionPeriod: (*math.HexOrDecimal256)(new(big.Int).Mul(big.NewInt(4*SecondsInYear+SecondsInDay), DecimalFactor)),                                // (4+1/365) years
		InflationCurveConvexity:   (*math.HexOrDecimal256)(new(big.Int).Div(new(big.Int).Mul(big.NewInt(-1_779), DecimalFactor), big.NewInt(1_000))),                 // -1.7794797758
	}

	// all percentage parameters needs to be scaled according to SLASHING_RATE_PRECISION
	DefaultAccountabilityConfig = &AccountabilityGenesis{
		InnocenceProofSubmissionWindow: 100, // 100 blocks
		BaseSlashingRateLow:            400, // 4%
		BaseSlashingRateMid:            600, // 6%
		BaseSlashingRateHigh:           800, // 8%
		CollusionFactor:                200, // 2%
		HistoryFactor:                  500, // 5%
		JailFactor:                     48,  // 48 epochs, i.e. 1 day with 30 mins epoch
	}

	/*
	* 1. InactivityThreshold and PastPerformanceWeight need to be scaled based on the OmissionAccountability.sol SCALE_FACTOR
	* 2. InitialSlashingRate needs to be scaled based on the Slasher.sol SLASHING_RATE_PRECISION
	* 3. the following equation needs to be respected: pastPerformanceWeight <= InactivityThreshold
	*    this ensures that a validator with 100% inactivity in epoch x and 0% inactivity in epoch x+n,
	*    will not be considered inactive again at epoch x+n
	* 4. lookbackWindow needs to be >= 1
	* 5. delta needs to be >= 2
	 */
	DefaultOmissionAccountabilityConfig = &OmissionAccountabilityGenesis{
		InactivityThreshold:    1000,   // 10%
		LookbackWindow:         40,     // 40 blocks
		PastPerformanceWeight:  1000,   // 10%
		InitialJailingPeriod:   10_000, // 10000 blocks
		InitialProbationPeriod: 24,     // 24 epochs
		InitialSlashingRate:    25,     // 0.25%
		Delta:                  5,      // 5 blocks
	}

	// same as the previous one, but with InactivityThreshold raised to 50%.
	// This is a very conservative threshold, and should be used for the first testnet we will run with omission enabled.
	// the idea is to avoid mass jailings at network startup. We can then lower down the threshold to 10% gradually at network runtime.
	ConservativeOmissionAccountabilityConfig = &OmissionAccountabilityGenesis{
		InactivityThreshold:    5000,   // 50%
		LookbackWindow:         40,     // 40 blocks
		PastPerformanceWeight:  1000,   // 10%
		InitialJailingPeriod:   10_000, // 10000 blocks
		InitialProbationPeriod: 24,     // 24 epochs
		InitialSlashingRate:    25,     // 0.25%
		Delta:                  5,      // 5 blocks
	}

	DefaultNonStakeableVestingGenesis = &NonStakeableVestingGenesis{
		NonStakeableContracts: make([]NonStakeableVestingData, 0),
	}

	DefaultStakeableVestingGenesis = &StakeableVestingGenesis{
		StakeableContracts: make([]StakeableVestingData, 0),
		TotalNominal:       big.NewInt(0),
	}

	DeployerAddress                        = common.Address{}
	AutonityContractAddress                = crypto.CreateAddress(DeployerAddress, 0)
	AccountabilityContractAddress          = crypto.CreateAddress(DeployerAddress, 1)
	OracleContractAddress                  = crypto.CreateAddress(DeployerAddress, 2)
	ACUContractAddress                     = crypto.CreateAddress(DeployerAddress, 3)
	SupplyControlContractAddress           = crypto.CreateAddress(DeployerAddress, 4)
	StabilizationContractAddress           = crypto.CreateAddress(DeployerAddress, 5)
	UpgradeManagerContractAddress          = crypto.CreateAddress(DeployerAddress, 6)
	InflationControllerContractAddress     = crypto.CreateAddress(DeployerAddress, 7)
	StakeableVestingManagerContractAddress = crypto.CreateAddress(DeployerAddress, 8)
	NonStakeableVestingContractAddress     = crypto.CreateAddress(DeployerAddress, 9)
	OmissionAccountabilityContractAddress  = crypto.CreateAddress(DeployerAddress, 10)
)

type AutonityContractGenesis struct {
	Bytecode                hexutil.Bytes         `json:"bytecode,omitempty" toml:",omitempty"`
	ABI                     *abi.ABI              `json:"abi,omitempty" toml:",omitempty"`
	MinBaseFee              uint64                `json:"minBaseFee"`
	EpochPeriod             uint64                `json:"epochPeriod"`
	UnbondingPeriod         uint64                `json:"unbondingPeriod"`
	BlockPeriod             uint64                `json:"blockPeriod"`
	MaxCommitteeSize        uint64                `json:"maxCommitteeSize"`
	MaxScheduleDuration     uint64                `json:"maxScheduleDuration"`
	Operator                common.Address        `json:"operator"`
	Treasury                common.Address        `json:"treasury"`
	WithheldRewardsPool     common.Address        `json:"withheldRewardsPool"`
	TreasuryFee             uint64                `json:"treasuryFee"`
	DelegationRate          uint64                `json:"delegationRate"`
	WithholdingThreshold    uint64                `json:"withholdingThreshold"`
	ProposerRewardRate      uint64                `json:"proposerRewardRate"`
	OracleRewardRate        uint64                `json:"oracleRewardRate"`
	InitialInflationReserve *math.HexOrDecimal256 `json:"initialInflationReserve"`
	Validators              []*Validator          `json:"validators"` // todo: Can we change that to []Validator
	Schedules               []Schedule            `json:"schedules"`
}

type AccountabilityGenesis struct {
	InnocenceProofSubmissionWindow uint64 `json:"innocenceProofSubmissionWindow"`

	// Slashing parameters
	BaseSlashingRateLow  uint64 `json:"baseSlashingRateLow"`
	BaseSlashingRateMid  uint64 `json:"baseSlashingRateMid"`
	BaseSlashingRateHigh uint64 `json:"baseSlashingRateHigh"`

	// Factors
	CollusionFactor uint64 `json:"collusionFactor"`
	HistoryFactor   uint64 `json:"historyFactor"`
	JailFactor      uint64 `json:"jailFactor"`
}

// OmissionAccountabilityGenesis defines the omission fault detection parameters
type OmissionAccountabilityGenesis struct {
	InactivityThreshold    uint64 `json:"inactivityThreshold"`
	LookbackWindow         uint64 `json:"LookbackWindow"`
	PastPerformanceWeight  uint64 `json:"pastPerformanceWeight"` // k belong to [0, 1), after scaling in the contract
	InitialJailingPeriod   uint64 `json:"initialJailingPeriod"`
	InitialProbationPeriod uint64 `json:"initialProbationPeriod"`
	InitialSlashingRate    uint64 `json:"initialSlashingRate"`
	Delta                  uint64 `json:"delta"`
}

type Validator struct {
	Treasury                 common.Address
	NodeAddress              *common.Address
	OracleAddress            common.Address
	Enode                    string
	CommissionRate           *big.Int
	BondedStake              *big.Int
	UnbondingStake           *big.Int
	UnbondingShares          *big.Int
	SelfBondedStake          *big.Int
	SelfUnbondingStake       *big.Int
	SelfUnbondingShares      *big.Int
	SelfUnbondingStakeLocked *big.Int
	LiquidStateContract      *common.Address
	LiquidSupply             *big.Int
	RegistrationBlock        *big.Int
	TotalSlashed             *big.Int
	JailReleaseBlock         *big.Int
	ConsensusKey             []byte //ABI packing does not support hexutil.Bytes, thus we need to introduce customized JSON Marshal/UnMarshal methods.
	State                    *uint8
}

// UnmarshalJSON and MarshalJSON are customized marshal and unmarshal methods to parse validators with validator key in
// hex string from genesis file, the raw type []byte is replaced by hexutil.Bytes.
func (v *Validator) UnmarshalJSON(input []byte) error {
	type validator struct {
		Treasury                 common.Address  `json:"treasury"`
		NodeAddress              *common.Address `json:"nodeAddress"`
		OracleAddress            common.Address  `json:"oracleAddress"`
		Enode                    string          `json:"enode"`
		CommissionRate           *big.Int        `json:"commissionRate"`
		BondedStake              *big.Int        `json:"bondedStake"`
		UnbondingStake           *big.Int        `json:"unbondingStake"`
		UnbondingShares          *big.Int        `json:"unbondingShares"`
		SelfBondedStake          *big.Int        `json:"selfBondedStake"`
		SelfUnbondingStake       *big.Int        `json:"selfUnbondingStake"`
		SelfUnbondingShares      *big.Int        `json:"selfUnbondingShares"`
		SelfUnbondingStakeLocked *big.Int        `json:"selfUnbondingStakeLocked"`
		LiquidStateContract      *common.Address `json:"liquidStateContract"`
		LiquidSupply             *big.Int        `json:"liquidSupply"`
		RegistrationBlock        *big.Int        `json:"registrationBlock"`
		TotalSlashed             *big.Int        `json:"totalSlashed"`
		JailReleaseBlock         *big.Int        `json:"jailReleaseBlock"`
		ConsensusKey             hexutil.Bytes   `json:"consensusKey"`
		State                    *uint8          `json:"state"`
	}

	var dec validator
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	v.Treasury = dec.Treasury
	v.NodeAddress = dec.NodeAddress
	v.OracleAddress = dec.OracleAddress
	v.Enode = dec.Enode
	v.CommissionRate = dec.CommissionRate
	v.BondedStake = dec.BondedStake
	v.UnbondingStake = dec.UnbondingStake
	v.UnbondingShares = dec.UnbondingShares
	v.SelfBondedStake = dec.SelfBondedStake
	v.SelfUnbondingStake = dec.SelfUnbondingStake
	v.SelfUnbondingShares = dec.SelfUnbondingShares
	v.SelfUnbondingStakeLocked = dec.SelfUnbondingStakeLocked
	v.LiquidStateContract = dec.LiquidStateContract
	v.LiquidSupply = dec.LiquidSupply
	v.RegistrationBlock = dec.RegistrationBlock
	v.TotalSlashed = dec.TotalSlashed
	v.JailReleaseBlock = dec.JailReleaseBlock
	v.ConsensusKey = dec.ConsensusKey
	v.State = dec.State

	return nil
}

func (v *Validator) MarshalJSON() ([]byte, error) {
	type validator struct {
		Treasury                 common.Address  `json:"treasury"`
		NodeAddress              *common.Address `json:"nodeAddress"`
		OracleAddress            common.Address  `json:"oracleAddress"`
		Enode                    string          `json:"enode"`
		CommissionRate           *big.Int        `json:"commissionRate"`
		BondedStake              *big.Int        `json:"bondedStake"`
		UnbondingStake           *big.Int        `json:"unbondingStake"`
		UnbondingShares          *big.Int        `json:"unbondingShares"`
		SelfBondedStake          *big.Int        `json:"selfBondedStake"`
		SelfUnbondingStake       *big.Int        `json:"selfUnbondingStake"`
		SelfUnbondingShares      *big.Int        `json:"selfUnbondingShares"`
		SelfUnbondingStakeLocked *big.Int        `json:"selfUnbondingStakeLocked"`
		LiquidStateContract      *common.Address `json:"liquidStateContract"`
		LiquidSupply             *big.Int        `json:"liquidSupply"`
		RegistrationBlock        *big.Int        `json:"registrationBlock"`
		TotalSlashed             *big.Int        `json:"totalSlashed"`
		JailReleaseBlock         *big.Int        `json:"jailReleaseBlock"`
		ConsensusKey             hexutil.Bytes   `json:"consensusKey"`
		State                    *uint8          `json:"state"`
	}

	var enc validator
	enc.Treasury = v.Treasury
	enc.NodeAddress = v.NodeAddress
	enc.OracleAddress = v.OracleAddress
	enc.Enode = v.Enode
	enc.CommissionRate = v.CommissionRate
	enc.BondedStake = v.BondedStake
	enc.UnbondingStake = v.UnbondingStake
	enc.UnbondingShares = v.UnbondingShares
	enc.SelfBondedStake = v.SelfBondedStake
	enc.SelfUnbondingStake = v.SelfUnbondingStake
	enc.SelfUnbondingShares = v.SelfUnbondingShares
	enc.SelfUnbondingStakeLocked = v.SelfUnbondingStakeLocked
	enc.LiquidStateContract = v.LiquidStateContract
	enc.LiquidSupply = v.LiquidSupply
	enc.RegistrationBlock = v.RegistrationBlock
	enc.TotalSlashed = v.TotalSlashed
	enc.JailReleaseBlock = v.JailReleaseBlock
	enc.ConsensusKey = v.ConsensusKey
	enc.State = v.State
	return json.Marshal(&enc)
}

// AddressFromEnode gets the account address from the user enode.
func (v *Validator) AddressFromEnode() (common.Address, error) {
	n, err := enode.ParseV4NoResolve(v.Enode)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to parse enode %q, error:%v", v.Enode, err)
	}
	return crypto.PubkeyToAddress(*n.Pubkey()), nil
}

func (v *Validator) Validate() error {
	if len(v.Enode) == 0 {
		return errors.New("enode must be specified")
	}

	if v.BondedStake == nil || v.BondedStake.Cmp(new(big.Int)) == 0 {
		return errors.New("bonded stake must be specified")
	}
	nodeAddr, err := v.AddressFromEnode()
	if err != nil {
		return err
	}
	// If address is set check it matches the address from the enode
	if v.NodeAddress != nil && *v.NodeAddress != nodeAddr {
		return fmt.Errorf("mismatching address %q and address from enode %q", v.NodeAddress.String(), nodeAddr.String())
	}
	v.NodeAddress = &nodeAddr

	if deep.Equal(v.OracleAddress, common.Address{}) {
		return fmt.Errorf("missing oracle address from genesis for node %q", nodeAddr.String())
	}

	if v.TotalSlashed == nil {
		v.TotalSlashed = new(big.Int)
	}
	if v.LiquidSupply == nil {
		v.LiquidSupply = new(big.Int)
	}
	if v.CommissionRate == nil {
		v.CommissionRate = new(big.Int)
	}
	if v.RegistrationBlock == nil {
		v.RegistrationBlock = new(big.Int)
	}
	if v.LiquidStateContract == nil {
		v.LiquidStateContract = new(common.Address)
	}
	if v.State == nil {
		v.State = new(uint8)
	}
	if v.SelfBondedStake == nil {
		v.SelfBondedStake = new(big.Int)
	}
	if v.SelfUnbondingStakeLocked == nil {
		v.SelfUnbondingStakeLocked = new(big.Int)
	}
	if v.SelfUnbondingStake == nil {
		v.SelfUnbondingStake = new(big.Int)
	}
	if v.SelfUnbondingShares == nil {
		v.SelfUnbondingShares = new(big.Int)
	}
	if v.UnbondingShares == nil {
		v.UnbondingShares = new(big.Int)
	}
	if v.UnbondingStake == nil {
		v.UnbondingStake = new(big.Int)
	}
	if v.JailReleaseBlock == nil {
		v.JailReleaseBlock = new(big.Int)
	}
	if v.CommissionRate != nil && v.CommissionRate.Cmp(big.NewInt(0)) != 0 {
		return fmt.Errorf("commission rate for enode %q not allowed", nodeAddr.String())
	}
	if _, err = blst.PublicKeyFromBytes(v.ConsensusKey); err != nil {
		return fmt.Errorf("cant decode bls public key: %w", err)
	}
	return nil
}

// OracleContractGenesis Autonity contract config. It is used for deployment.
type OracleContractGenesis struct {
	Bytecode                  hexutil.Bytes `json:"bytecode,omitempty" toml:",omitempty"`
	ABI                       *abi.ABI      `json:"abi,omitempty" toml:",omitempty"`
	Symbols                   []string      `json:"symbols"`
	VotePeriod                uint64        `json:"votePeriod"`
	OutlierDetectionThreshold uint64        `json:"outlierDetectionThreshold"`
	OutlierSlashingThreshold  uint64        `json:"outlierSlashingThreshold"`
	BaseSlashingRate          uint64        `json:"baseSlashingRate"`
}

// SetDefaults prepares the AutonityContractGenesis by filling in missing fields.
// It returns an error if the configuration is invalid.
func (g *OracleContractGenesis) SetDefaults() error {
	if g.Bytecode == nil && g.ABI != nil || g.Bytecode != nil && g.ABI == nil {
		return errors.New("it is an error to set only of oracle contract abi or bytecode")
	}
	if g.Bytecode == nil && g.ABI == nil {
		g.ABI = &generated.OracleAbi
		g.Bytecode = generated.OracleBytecode
	}
	if len(g.Symbols) == 0 {
		g.Symbols = OracleInitialSymbols
	}
	if g.VotePeriod == 0 {
		g.VotePeriod = OracleVotePeriod
	}
	if g.OutlierSlashingThreshold == 0 {
		g.OutlierSlashingThreshold = DefaultGenesisOracleConfig.OutlierSlashingThreshold
	}
	if g.BaseSlashingRate == 0 {
		g.BaseSlashingRate = DefaultGenesisOracleConfig.BaseSlashingRate
	}
	if g.OutlierDetectionThreshold == 0 {
		g.OutlierDetectionThreshold = DefaultGenesisOracleConfig.OutlierDetectionThreshold
	}
	return nil
}

type AcuContractGenesis struct {
	Symbols    []string
	Quantities []uint64
	Scale      uint64
}

func (acu *AcuContractGenesis) SetDefaults() {
	if acu.Symbols == nil {
		acu.Symbols = DefaultAcuContractGenesis.Symbols
	}
	if acu.Quantities == nil {
		acu.Quantities = DefaultAcuContractGenesis.Quantities
	}
}

type StabilizationContractGenesis struct {
	BorrowInterestRate        *math.HexOrDecimal256
	LiquidationRatio          *math.HexOrDecimal256
	MinCollateralizationRatio *math.HexOrDecimal256
	MinDebtRequirement        *math.HexOrDecimal256
	TargetPrice               *math.HexOrDecimal256
}

func (s *StabilizationContractGenesis) SetDefaults() {
	if s.BorrowInterestRate == nil {
		s.BorrowInterestRate = DefaultStabilizationGenesis.BorrowInterestRate
	}
	if s.LiquidationRatio == nil {
		s.LiquidationRatio = DefaultStabilizationGenesis.LiquidationRatio
	}
	if s.MinCollateralizationRatio == nil {
		s.MinCollateralizationRatio = DefaultStabilizationGenesis.MinCollateralizationRatio
	}
	if s.MinDebtRequirement == nil {
		s.MinDebtRequirement = DefaultStabilizationGenesis.MinDebtRequirement
	}
	if s.TargetPrice == nil {
		s.TargetPrice = DefaultStabilizationGenesis.TargetPrice
	}
}

type SupplyControlGenesis struct {
	InitialAllocation *math.HexOrDecimal256
}

func (s *SupplyControlGenesis) SetDefaults() {
	if s.InitialAllocation == nil {
		s.InitialAllocation = DefaultSupplyControlGenesis.InitialAllocation
	}
}

type InflationControllerGenesis struct {
	// Those parameters need to be compatible with the solidity SD59x18 format
	InflationRateInitial      *math.HexOrDecimal256 `json:"inflationRateInitial"`
	InflationRateTransition   *math.HexOrDecimal256 `json:"inflationRateTransition"`
	InflationReserveDecayRate *math.HexOrDecimal256 `json:"inflationReserveDecayRate"`
	InflationTransitionPeriod *math.HexOrDecimal256 `json:"inflationTransitionPeriod"`
	InflationCurveConvexity   *math.HexOrDecimal256 `json:"inflationCurveConvexity"`
}

func (s *InflationControllerGenesis) SetDefaults() {
	if s.InflationRateInitial == nil {
		s.InflationRateInitial = DefaultInflationControllerGenesis.InflationRateInitial
	}
	if s.InflationRateTransition == nil {
		s.InflationRateTransition = DefaultInflationControllerGenesis.InflationRateTransition
	}
	if s.InflationReserveDecayRate == nil {
		s.InflationReserveDecayRate = DefaultInflationControllerGenesis.InflationReserveDecayRate
	}
	if s.InflationTransitionPeriod == nil {
		s.InflationTransitionPeriod = DefaultInflationControllerGenesis.InflationTransitionPeriod
	}
	if s.InflationCurveConvexity == nil {
		s.InflationCurveConvexity = DefaultInflationControllerGenesis.InflationCurveConvexity
	}
}

type NonStakeableVestingGenesis struct {
	NonStakeableContracts []NonStakeableVestingData `json:"nonStakeableVestingContracts"`
}

type Schedule struct {
	Start         *big.Int       `json:"startTime"`
	TotalDuration *big.Int       `json:"totalDuration"`
	Amount        *big.Int       `json:"amount"`
	VaultAddress  common.Address `json:"vaultAddress"`
}

func (s *Schedule) Validate() error {
	if s.Start == nil {
		return errors.New("start time must be specified")
	}
	if s.TotalDuration == nil {
		return errors.New("total duration must be specified")
	}
	if s.Amount == nil {
		return errors.New("amount must be specified")
	}
	if deep.Equal(s.VaultAddress, common.Address{}) {
		s.VaultAddress = NonStakeableVestingContractAddress
	}
	return nil
}

type NonStakeableVestingData struct {
	Beneficiary   common.Address `json:"beneficiary"`
	Amount        *big.Int       `json:"amount"`
	ScheduleID    *big.Int       `json:"scheduleID"`
	CliffDuration *big.Int       `json:"cliffDuration"`
}

type StakeableVestingGenesis struct {
	TotalNominal       *big.Int               `json:"totalNominal"`
	StakeableContracts []StakeableVestingData `json:"stakeableVestingContracts"`
}

type StakeableVestingData struct {
	Beneficiary   common.Address `json:"beneficiary"`
	Amount        *big.Int       `json:"amount"`
	Start         *big.Int       `json:"startTime"`
	CliffDuration *big.Int       `json:"cliffDuration"`
	TotalDuration *big.Int       `json:"totalDuration"`
}

func (s *StakeableVestingGenesis) SetDefaults() {
	if s.TotalNominal == nil {
		s.TotalNominal = DefaultStakeableVestingGenesis.TotalNominal
	}
}
