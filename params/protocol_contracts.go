package params

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/autonity/autonity/crypto/blst"
	"math/big"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params/generated"
)

var (
	//Global default protocol parameters
	GovernanceOperator = common.HexToAddress("0x1336000000000000000000000000000000000000")

	//Oracle Contract defaults
	OracleVotePeriod           = uint64(30)
	OracleInitialSymbols       = []string{"AUD-USD", "CAD-USD", "EUR-USD", "GBP-USD", "JPY-USD", "SEK-USD", "ATN-USD", "NTN-USD", "NTN-ATN"}
	DefaultGenesisOracleConfig = &OracleContractGenesis{
		Bytecode:   generated.OracleBytecode,
		ABI:        &generated.OracleAbi,
		Symbols:    OracleInitialSymbols,
		VotePeriod: OracleVotePeriod,
	}

	// DefaultAcuContractGenesis contains the default values for the ASM ACU contract
	DefaultAcuContractGenesis = &AcuContractGenesis{
		Symbols:    []string{"AUD-USD", "CAD-USD", "EUR-USD", "GBP-USD", "JPY-USD", "USD-USD", "SEK-USD"},
		Quantities: []uint64{21_300, 18_700, 14_300, 10_400, 1_760_000, 18_000, 141_000},
		Scale:      uint64(5),
	}

	// DefaultStabilizationGenesis contains the default values for the ASM Stabilization contract
	DefaultStabilizationGenesis = &StabilizationContractGenesis{
		BorrowInterestRate:        (*math.HexOrDecimal256)(math.MustParseBig256("50_000_000_000_000_000")),
		LiquidationRatio:          (*math.HexOrDecimal256)(math.MustParseBig256("1_800_000_000_000_000_000")),
		MinCollateralizationRatio: (*math.HexOrDecimal256)(math.MustParseBig256("2_000_000_000_000_000_000")),
		MinDebtRequirement:        (*math.HexOrDecimal256)(math.MustParseBig256("1_000_000")),
		TargetPrice:               (*math.HexOrDecimal256)(math.MustParseBig256("1_000_000_000_000_000_000")),
	}

	DefaultSupplyControlGenesis = &SupplyControlGenesis{
		InitialAllocation: (*math.HexOrDecimal256)(new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), common.Big1)),
	}

	DefaultAccountabilityConfig = &AccountabilityGenesis{
		InnocenceProofSubmissionWindow: 100,
		BaseSlashingRateLow:            1000, // 10%
		BaseSlashingRateMid:            2000, // 20%
		CollusionFactor:                500,  // 5%
		HistoryFactor:                  750,  // 7.5%
		JailFactor:                     48,   // 1 day with 30 mins epoch
		SlashingRatePrecision:          10_000,
	}
)

type AutonityContractGenesis struct {
	Bytecode         hexutil.Bytes  `json:"bytecode,omitempty" toml:",omitempty"`
	ABI              *abi.ABI       `json:"abi,omitempty" toml:",omitempty"`
	MinBaseFee       uint64         `json:"minBaseFee"`
	EpochPeriod      uint64         `json:"epochPeriod"`
	UnbondingPeriod  uint64         `json:"unbondingPeriod"`
	BlockPeriod      uint64         `json:"blockPeriod"`
	MaxCommitteeSize uint64         `json:"maxCommitteeSize"`
	Operator         common.Address `json:"operator"`
	Treasury         common.Address `json:"treasury"`
	TreasuryFee      uint64         `json:"treasuryFee"`
	DelegationRate   uint64         `json:"delegationRate"`
	Validators       []*Validator   `json:"validators"`
}

type AccountabilityGenesis struct {
	InnocenceProofSubmissionWindow uint64 `json:"innocenceProofSubmissionWindow"`
	// Slashing parameters
	BaseSlashingRateLow   uint64 `json:"baseSlashingRateLow"`
	BaseSlashingRateMid   uint64 `json:"baseSlashingRateMid"`
	CollusionFactor       uint64 `json:"collusionFactor"`
	HistoryFactor         uint64 `json:"historyFactor"`
	JailFactor            uint64 `json:"jailFactor"`
	SlashingRatePrecision uint64 `json:"slashingRatePrecision"`
}

// Prepare prepares the AutonityContractGenesis by filling in missing fields.
// It returns an error if the configuration is invalid.
func (g *AutonityContractGenesis) Prepare() error {
	if g.Bytecode == nil && g.ABI != nil || g.Bytecode != nil && g.ABI == nil {
		return errors.New("autonity contract abi or bytecode missing")
	}
	if g.Bytecode == nil && g.ABI == nil {
		g.ABI = &generated.AutonityAbi
		g.Bytecode = generated.AutonityBytecode
	}
	if g.Operator == (common.Address{}) {
		g.Operator = GovernanceOperator
	}
	if g.MaxCommitteeSize == 0 {
		return errors.New("invalid max committee size")
	}
	if len(g.Validators) == 0 {
		return errors.New("no initial validators")
	}
	for i, v := range g.Validators {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("error parsing validator %d, err: %v", i+1, err)
		}
	}
	return nil
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
	LiquidContract           *common.Address
	LiquidSupply             *big.Int
	RegistrationBlock        *big.Int
	TotalSlashed             *big.Int
	JailReleaseBlock         *big.Int
	ProvableFaultCount       *big.Int
	ValidatorKey             []byte //ABI packing does not support hexutil.Bytes, thus we need to introduce customized JSON Marshal/UnMarshal methods.
	State                    *uint8
}

// UnmarshalJSON and MarshalJSON are customized marshal and unmarshal methods to parse validators with activity key in
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
		LiquidContract           *common.Address `json:"liquidContract"`
		LiquidSupply             *big.Int        `json:"liquidSupply"`
		RegistrationBlock        *big.Int        `json:"registrationBlock"`
		TotalSlashed             *big.Int        `json:"totalSlashed"`
		JailReleaseBlock         *big.Int        `json:"jailReleaseBlock"`
		ProvableFaultCount       *big.Int        `json:"provableFaultCount"`
		ValidatorKey             hexutil.Bytes   `json:"validatorKey"`
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
	v.LiquidContract = dec.LiquidContract
	v.LiquidSupply = dec.LiquidSupply
	v.RegistrationBlock = dec.RegistrationBlock
	v.TotalSlashed = dec.TotalSlashed
	v.JailReleaseBlock = dec.JailReleaseBlock
	v.ProvableFaultCount = dec.ProvableFaultCount
	v.ValidatorKey = dec.ValidatorKey
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
		LiquidContract           *common.Address `json:"liquidContract"`
		LiquidSupply             *big.Int        `json:"liquidSupply"`
		RegistrationBlock        *big.Int        `json:"registrationBlock"`
		TotalSlashed             *big.Int        `json:"totalSlashed"`
		JailReleaseBlock         *big.Int        `json:"jailReleaseBlock"`
		ProvableFaultCount       *big.Int        `json:"provableFaultCount"`
		ValidatorKey             hexutil.Bytes   `json:"validatorKey"`
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
	enc.LiquidContract = v.LiquidContract
	enc.LiquidSupply = v.LiquidSupply
	enc.RegistrationBlock = v.RegistrationBlock
	enc.TotalSlashed = v.TotalSlashed
	enc.JailReleaseBlock = v.JailReleaseBlock
	enc.ProvableFaultCount = v.ProvableFaultCount
	enc.ValidatorKey = v.ValidatorKey
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

	_, err := blst.PublicKeyFromBytes(v.ValidatorKey)
	if err != nil {
		return errors.New("cannot decode bls public key")
	}

	if v.BondedStake == nil || v.BondedStake.Cmp(new(big.Int)) == 0 {
		return errors.New("bonded stake must be specified")
	}
	a, err := v.AddressFromEnode()
	if err != nil {
		return err
	}
	// If address is set check it matches the address from the enode
	if v.NodeAddress != nil && *v.NodeAddress != a {
		return fmt.Errorf("mismatching address %q and address from enode %q", v.NodeAddress.String(), a.String())
	}
	v.NodeAddress = &a
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
	if v.LiquidContract == nil {
		v.LiquidContract = new(common.Address)
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
	if v.ProvableFaultCount == nil {
		v.ProvableFaultCount = new(big.Int)
	}
	if v.CommissionRate != nil && v.CommissionRate.Cmp(big.NewInt(0)) != 0 {
		return fmt.Errorf("commission rate for enode %q not allowed", a.String())
	}
	return nil
}

// OracleContractGenesis Autonity contract config. It'is used for deployment.
type OracleContractGenesis struct {
	Bytecode   hexutil.Bytes `json:"bytecode,omitempty" toml:",omitempty"`
	ABI        *abi.ABI      `json:"abi,omitempty" toml:",omitempty"`
	Symbols    []string      `json:"symbols"`
	VotePeriod uint64        `json:"votePeriod"`
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
