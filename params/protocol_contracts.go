package params

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params/generated"
	"math/big"
)

var (
	//Global default protocol parameters
	GovernanceOperator = common.HexToAddress("0x1336000000000000000000000000000000000000")

	//Oracle Contract defaults
	OracleVotePeriod     = uint64(30)
	OracleInitialSymbols = []string{"NTN/USD", "NTN/AUD", "NTN/CAD", "NTN/EUR", "NTN/GBP", "NTN/JPY", "NTN/SEK"}
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
	Treasury           common.Address  `abi:"treasury"`
	NodeAddress        *common.Address `abi:"nodeAddress"`
	OracleAddress      common.Address  `abi:"oracleAddress"`
	Enode              string          `abi:"enode"`
	CommissionRate     *big.Int        `abi:"commissionRate"`
	BondedStake        *big.Int        `abi:"bondedStake"`
	SelfBondedStake    *big.Int        `abi:"selfBondedStake"`
	TotalSlashed       *big.Int        `abi:"totalSlashed"`
	LiquidContract     *common.Address `abi:"liquidContract"`
	LiquidSupply       *big.Int        `abi:"liquidSupply"`
	RegistrationBlock  *big.Int        `abi:"registrationBlock"`
	JailReleaseBlock   *big.Int        `abi:"jailReleaseBlock"`
	ProvableFaultCount *big.Int        `abi:"provableFaultCount"`
	State              *uint8          `abi:"state"`
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

// Prepare prepares the AutonityContractGenesis by filling in missing fields.
// It returns an error if the configuration is invalid.
func (g *OracleContractGenesis) Prepare() error {
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

func DefaultGenesisOracleConfig() *OracleContractGenesis {
	return &OracleContractGenesis{
		Bytecode:   generated.OracleBytecode,
		ABI:        &generated.OracleAbi,
		Symbols:    OracleInitialSymbols,
		VotePeriod: OracleVotePeriod,
	}
}
