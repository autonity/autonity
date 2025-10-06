package params

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params/generated"
)

// wrapper for protocol contract addresses, to easily access their "name"
type ProtocolContract common.Address

func (contract ProtocolContract) String() string {
	switch common.Address(contract) {
	case AutonityContractAddress:
		return "AutonityContract"
	case AccountabilityContractAddress:
		return "AccountabilityContract"
	case OracleContractAddress:
		return "OracleContract"
	case ACUContractAddress:
		return "ACUContract"
	case SupplyControlContractAddress:
		return "SupplyControlContract"
	case StabilizationContractAddress:
		return "StabilizationContract"
	case UpgradeManagerContractAddress:
		return "UpgradeManagerContract"
	case InflationControllerContractAddress:
		return "InflationControllerContract"
	case OmissionAccountabilityContractAddress:
		return "OmissionAccountabilityContract"
	case AuctioneerContractAddress:
		return "AuctioneerContract"
	case ASMGroupAddress:
		return "ASMGroup"
	case ProtocolGroupAddress:
		return "ProtocolGroup"
	default:
		return "Unknown"
	}
}

func (contract ProtocolContract) Address() common.Address {
	return common.Address(contract)
}

type ProtocolContractVersion struct {
	Hash     common.Hash      // == codeHash for a single contract, hash(codeHash1,codeHash2,...) for a contract group (e.g. ASM, all contracts)
	Contract ProtocolContract // address for single contract, 0x000000 for a contract group (e.g. ASM, all contracts)
	Version  string           // semver version string
}

func (contractVersion ProtocolContractVersion) String() string {
	return fmt.Sprintf("%s-%s (address: %s hash: %s)", contractVersion.Contract.String(), contractVersion.Version, contractVersion.Contract.Address(), contractVersion.Hash)
}

var (
	DecimalPrecision = int64(18)
	SecondsInYear    = int64(365 * 24 * 60 * 60)
	SecondsInDay     = int64(24 * 60 * 60)
	DecimalFactor    = new(big.Int).Exp(big.NewInt(10), big.NewInt(DecimalPrecision), nil)

	//Oracle Contract defaults
	OracleVotePeriod           = uint64(30)
	OracleInitialSymbols       = []string{"AUD-USD", "CAD-USD", "EUR-USD", "GBP-USD", "JPY-USD", "SEK-USD", "ATN-USD", "NTN-USD", "NTN-ATN"}
	DefaultGenesisOracleConfig = &OracleContractGenesis{
		Symbols:                   OracleInitialSymbols,
		VotePeriod:                OracleVotePeriod,
		OutlierDetectionThreshold: 10,  // 10%
		OutlierSlashingThreshold:  225, // 15%
		BaseSlashingRate:          10,
		NonRevealThreshold:        3,
		RevealResetInterval:       10,
		SlashingRateCap:           1000, // 10%
	}

	DefaultNTNGenesisAllocation = new(big.Int).Mul(big.NewInt(60_000_000), DecimalFactor) // 60 mil NTN
	// TODO: update `DefautlGenesisBonding`
	DefaultGenesisBonding = new(big.Int).Mul(big.NewInt(0), DecimalFactor)

	// DefaultAcuContractGenesis contains the default values for the ASM ACU contract
	DefaultAcuContractGenesis = &AcuContractGenesis{
		Symbols:    []string{"AUD-USD", "CAD-USD", "EUR-USD", "GBP-USD", "JPY-USD", "SEK-USD", "USD-USD"},
		Quantities: []uint64{1_744_583, 1_598_986, 1_058_522, 886_091, 175_605_573, 12_318_802, 1_148_285},
		Scale:      uint64(7),
	}

	// DefaultStabilizationGenesis contains the default values for the ASM Stabilization contract
	DefaultStabilizationGenesis = &StabilizationContractGenesis{
		BorrowInterestRate:        (*math.HexOrDecimal256)(math.MustParseBig256("50_000_000_000_000_000")),
		AnnouncementWindow:        (*math.HexOrDecimal256)(math.MustParseBig256("3600")), // 1 hour
		LiquidationRatio:          (*math.HexOrDecimal256)(math.MustParseBig256("1_800_000_000_000_000_000")),
		MinCollateralizationRatio: (*math.HexOrDecimal256)(math.MustParseBig256("2_000_000_000_000_000_000")),
		MinDebtRequirement:        (*math.HexOrDecimal256)(math.MustParseBig256("1_000_000")),
		TargetPrice:               (*math.HexOrDecimal256)(math.MustParseBig256("1_618_034_000_000_000_000")),
		DefaultNTNATNPrice:        (*math.HexOrDecimal256)(math.MustParseBig256("1_000_000_000_000_000_000")),
		DefaultNTNUSDPrice:        (*math.HexOrDecimal256)(math.MustParseBig256("1_600_000_000_000_000_000")),
		DefaultACUUSDPrice:        (*math.HexOrDecimal256)(math.MustParseBig256("0")),
	}

	// ToDo: Add the real default values for the Auctioneer contract
	DefaultAuctioneerGenesis = &AuctioneerContractGenesis{
		LiquidationAuctionDuration: big.NewInt(60),                                        // 60 blocks
		InterestAuctionDuration:    big.NewInt(60),                                        // 60 blocks
		InterestAuctionDiscount:    new(big.Int).Exp(big.NewInt(10), big.NewInt(17), nil), // 0.1
		InterestAuctionThreshold:   new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), // 1 ATN
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
		InnocenceProofSubmissionWindow: 120, // 120 blocks
		Delta:                          10,  // 10 blocks
		Range:                          128, // 256 blocks
		BaseSlashingRateLow:            50,  // 0.5%
		BaseSlashingRateMid:            100, // 1%
		BaseSlashingRateHigh:           200, // 2%
		CollusionFactor:                25,  // 0.25%
		HistoryFactor:                  150, // 1.5%
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

	DeployerAddress                       = common.Address{}
	AutonityContractAddress               = crypto.CreateAddress(DeployerAddress, 0)
	AccountabilityContractAddress         = crypto.CreateAddress(DeployerAddress, 1)
	OracleContractAddress                 = crypto.CreateAddress(DeployerAddress, 2)
	ACUContractAddress                    = crypto.CreateAddress(DeployerAddress, 3)
	SupplyControlContractAddress          = crypto.CreateAddress(DeployerAddress, 4)
	StabilizationContractAddress          = crypto.CreateAddress(DeployerAddress, 5)
	UpgradeManagerContractAddress         = crypto.CreateAddress(DeployerAddress, 6)
	InflationControllerContractAddress    = crypto.CreateAddress(DeployerAddress, 7)
	OmissionAccountabilityContractAddress = crypto.CreateAddress(DeployerAddress, 8)
	AuctioneerContractAddress             = crypto.CreateAddress(DeployerAddress, 9)

	// mock addresses for contract groups. Used for contract versioning purposes
	ASMGroupAddress      = crypto.CreateAddress(DeployerAddress, math.MaxUint64-1)
	ProtocolGroupAddress = crypto.CreateAddress(DeployerAddress, math.MaxUint64)

	// NOTE: do not change order and add new contracts at the end if needed
	// otherwise contract groups version hashes will change
	ProtocolContracts = []ProtocolContract{
		ProtocolContract(AutonityContractAddress),
		ProtocolContract(AccountabilityContractAddress),
		ProtocolContract(OracleContractAddress),
		ProtocolContract(ACUContractAddress),
		ProtocolContract(SupplyControlContractAddress),
		ProtocolContract(StabilizationContractAddress),
		ProtocolContract(UpgradeManagerContractAddress),
		ProtocolContract(InflationControllerContractAddress),
		ProtocolContract(OmissionAccountabilityContractAddress),
		ProtocolContract(AuctioneerContractAddress),
	}

	ASMContracts = []ProtocolContract{
		ProtocolContract(ACUContractAddress),
		ProtocolContract(SupplyControlContractAddress),
		ProtocolContract(StabilizationContractAddress),
		ProtocolContract(InflationControllerContractAddress),
		ProtocolContract(AuctioneerContractAddress),
	}

	// maps codeHash --> (contract,version)
	// needs to be manually updated whenever there is an upgrade
	VersionHistory = map[common.Hash]ProtocolContractVersion{
		// version 1.0.0
		common.HexToHash("0xc74124cdea7c515bdde48bdf65e6f2a4d0f8eab278c8a4abc22f236c4391482d"): {
			Hash:     common.HexToHash("0xc74124cdea7c515bdde48bdf65e6f2a4d0f8eab278c8a4abc22f236c4391482d"),
			Contract: ProtocolContract(AutonityContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x4adac12d20f59528a61862a990546fbf5c95620662ad11291ce5697cfd2df094"): {
			Hash:     common.HexToHash("0x4adac12d20f59528a61862a990546fbf5c95620662ad11291ce5697cfd2df094"),
			Contract: ProtocolContract(AccountabilityContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x2a7bd44f6a7f6299461120eab732bd69cd5cd6fe1b8057b543a1ca093cf72202"): {
			Hash:     common.HexToHash("0x2a7bd44f6a7f6299461120eab732bd69cd5cd6fe1b8057b543a1ca093cf72202"),
			Contract: ProtocolContract(OracleContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x21a1125b1e7f814d380048581c22818ec638d06d1975ea773377abf53425669f"): {
			Hash:     common.HexToHash("0x21a1125b1e7f814d380048581c22818ec638d06d1975ea773377abf53425669f"),
			Contract: ProtocolContract(ACUContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x5541ab53bd6b62c3d52d69fe47005fc2ad0bcbf0134f463e982504031dfb6620"): {
			Hash:     common.HexToHash("0x5541ab53bd6b62c3d52d69fe47005fc2ad0bcbf0134f463e982504031dfb6620"),
			Contract: ProtocolContract(SupplyControlContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x22e0af6bcfb03c7522336066558a6caae329385d8fe68d4a3d1e0e796a2f6187"): {
			Hash:     common.HexToHash("0x22e0af6bcfb03c7522336066558a6caae329385d8fe68d4a3d1e0e796a2f6187"),
			Contract: ProtocolContract(StabilizationContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0xe372907e641e0f22d54a93466ce8487f9ecd1be2c6b8fc7fdad073451d00e0f2"): {
			Hash:     common.HexToHash("0xe372907e641e0f22d54a93466ce8487f9ecd1be2c6b8fc7fdad073451d00e0f2"),
			Contract: ProtocolContract(UpgradeManagerContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x00d2855a05fac35be477a8bd238bdd79535d372fedf5251300508be4c4eb76ff"): {
			Hash:     common.HexToHash("0x00d2855a05fac35be477a8bd238bdd79535d372fedf5251300508be4c4eb76ff"),
			Contract: ProtocolContract(InflationControllerContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x4118ecbe59133ce415ec85b037291963bd74a379887a24987f6d100f080ebf47"): {
			Hash:     common.HexToHash("0x4118ecbe59133ce415ec85b037291963bd74a379887a24987f6d100f080ebf47"),
			Contract: ProtocolContract(OmissionAccountabilityContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0xad6401716eded73fe64060a1b24ac8b2a48044989123c7778715ef17510d8713"): {
			Hash:     common.HexToHash("0xad6401716eded73fe64060a1b24ac8b2a48044989123c7778715ef17510d8713"),
			Contract: ProtocolContract(AuctioneerContractAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x47854920b5e065a4fc10cf14c4029c8b477f96f6fbfa2100d7904e4b83f4a032"): {
			Hash:     common.HexToHash("0x47854920b5e065a4fc10cf14c4029c8b477f96f6fbfa2100d7904e4b83f4a032"),
			Contract: ProtocolContract(ProtocolGroupAddress),
			Version:  "1.0.0",
		},
		common.HexToHash("0x194dba471e6af77d4012a15894673f962b4108892dd61043255ade30518f768b"): {
			Hash:     common.HexToHash("0x194dba471e6af77d4012a15894673f962b4108892dd61043255ade30518f768b"),
			Contract: ProtocolContract(ASMGroupAddress),
			Version:  "1.0.0",
		},
	}
)

type AutonityContractGenesis struct {
	Bytecode                 hexutil.Bytes         `json:"bytecode,omitempty" toml:",omitempty"`
	ABI                      *abi.ABI              `json:"abi,omitempty" toml:",omitempty"`
	MinBaseFee               uint64                `json:"minBaseFee"`
	EpochPeriod              uint64                `json:"epochPeriod"`
	UnbondingPeriod          uint64                `json:"unbondingPeriod"`
	BlockPeriod              uint64                `json:"blockPeriod"`
	MaxCommitteeSize         uint64                `json:"maxCommitteeSize"`
	MaxScheduleDuration      uint64                `json:"maxScheduleDuration"`
	GasLimit                 uint64                `json:"gasLimit"`
	GasLimitBoundDivisor     uint64                `json:"gasLimitBoundDivisor"`
	BaseFeeChangeDenominator uint64                `json:"baseFeeChangeDenominator"`
	ElasticityMultiplier     uint64                `json:"elasticityMultiplier"`
	ClusteringThreshold      uint64                `json:"clusteringThreshold"`
	Operator                 common.Address        `json:"operator"`
	Treasury                 common.Address        `json:"treasury"`
	WithheldRewardsPool      common.Address        `json:"withheldRewardsPool"`
	TreasuryFee              uint64                `json:"treasuryFee"`
	DelegationRate           uint64                `json:"delegationRate"`
	WithholdingThreshold     uint64                `json:"withholdingThreshold"`
	ProposerRewardRate       uint64                `json:"proposerRewardRate"`
	OracleRewardRate         uint64                `json:"oracleRewardRate"`
	InitialInflationReserve  *math.HexOrDecimal256 `json:"initialInflationReserve"`
	Validators               []*Validator          `json:"validators"` // todo: Can we change that to []Validator
	Schedules                []Schedule            `json:"schedules"`
	SkipGenesisVerification  bool                  `json:"skipGenesisVerification"`
	TokenMint                *math.HexOrDecimal256 `json:"tokenMint"`
	TokenBond                *math.HexOrDecimal256 `json:"tokenBond"`
}

type AccountabilityGenesis struct {
	InnocenceProofSubmissionWindow uint64 `json:"innocenceProofSubmissionWindow"`
	Delta                          uint64 `json:"delta"`
	Range                          uint64 `json:"range"`

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
	ConversionRatio          *big.Int
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
		ConversionRatio          *big.Int        `json:"conversionRatio"`
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
	v.ConversionRatio = dec.ConversionRatio

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
		ConversionRatio          *big.Int        `json:"conversionRatio"`
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
	enc.ConversionRatio = v.ConversionRatio
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

	if v.OracleAddress == common.ZeroAddress {
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
	// sanitize starting conversion ratio to CONVERSION_RATIO_SCALE_FACTOR (10 ^ 18)
	if v.ConversionRatio == nil || v.ConversionRatio.Cmp(big.NewInt(1e18)) != 0 {
		v.ConversionRatio = new(big.Int).SetUint64(1e18) // everyone starts at 1:1 at genesis
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
	NonRevealThreshold        uint64        `json:"nonRevealThreshold"`
	RevealResetInterval       uint64        `json:"revealResetInterval"`
	SlashingRateCap           uint64        `json:"slashingRateCap"`
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
	// at genesis, we allow some tolerance for missed reveal
	if g.NonRevealThreshold == 0 {
		g.NonRevealThreshold = DefaultGenesisOracleConfig.NonRevealThreshold
	}
	if g.RevealResetInterval == 0 {
		g.RevealResetInterval = DefaultGenesisOracleConfig.RevealResetInterval
	}
	if g.SlashingRateCap == 0 {
		g.SlashingRateCap = DefaultGenesisOracleConfig.SlashingRateCap
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
	AnnouncementWindow        *math.HexOrDecimal256
	LiquidationRatio          *math.HexOrDecimal256
	MinCollateralizationRatio *math.HexOrDecimal256
	MinDebtRequirement        *math.HexOrDecimal256
	TargetPrice               *math.HexOrDecimal256
	DefaultNTNATNPrice        *math.HexOrDecimal256
	DefaultNTNUSDPrice        *math.HexOrDecimal256
	DefaultACUUSDPrice        *math.HexOrDecimal256
}

func (s *StabilizationContractGenesis) SetDefaults() {
	if s.BorrowInterestRate == nil {
		s.BorrowInterestRate = DefaultStabilizationGenesis.BorrowInterestRate
	}
	if s.AnnouncementWindow == nil {
		s.AnnouncementWindow = DefaultStabilizationGenesis.AnnouncementWindow
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
	if s.DefaultNTNATNPrice == nil {
		s.DefaultNTNATNPrice = DefaultStabilizationGenesis.DefaultNTNATNPrice
	}
	if s.DefaultNTNUSDPrice == nil {
		s.DefaultNTNUSDPrice = DefaultStabilizationGenesis.DefaultNTNUSDPrice
	}
	if s.DefaultACUUSDPrice == nil {
		s.DefaultACUUSDPrice = DefaultStabilizationGenesis.DefaultACUUSDPrice
	}
}

type AuctioneerContractGenesis struct {
	LiquidationAuctionDuration *big.Int
	InterestAuctionDuration    *big.Int
	InterestAuctionDiscount    *big.Int // value between [0,1) with SCALE_FACTOR precision
	InterestAuctionThreshold   *big.Int // in ATN
}

func (a *AuctioneerContractGenesis) SetDefaults() {
	if a.LiquidationAuctionDuration == nil {
		a.LiquidationAuctionDuration = DefaultAuctioneerGenesis.LiquidationAuctionDuration
	}
	if a.InterestAuctionDuration == nil {
		a.InterestAuctionDuration = DefaultAuctioneerGenesis.InterestAuctionDuration
	}
	if a.InterestAuctionDiscount == nil {
		a.InterestAuctionDiscount = DefaultAuctioneerGenesis.InterestAuctionDiscount
	}
	if a.InterestAuctionThreshold == nil {
		a.InterestAuctionThreshold = DefaultAuctioneerGenesis.InterestAuctionThreshold
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
	if s.VaultAddress == common.ZeroAddress {
		return errors.New("vault address must be specified")
	}
	return nil
}
