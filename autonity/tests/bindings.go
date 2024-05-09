// Code generated for internal testing purposes only - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tests

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AccountabilityConfig is an auto generated low-level Go binding around an user-defined struct.
type AccountabilityConfig struct {
	InnocenceProofSubmissionWindow *big.Int
	BaseSlashingRateLow            *big.Int
	BaseSlashingRateMid            *big.Int
	CollusionFactor                *big.Int
	HistoryFactor                  *big.Int
	JailFactor                     *big.Int
	SlashingRatePrecision          *big.Int
}

// AccountabilityEvent is an auto generated low-level Go binding around an user-defined struct.
type AccountabilityEvent struct {
	Chunks         uint8
	ChunkId        uint8
	EventType      uint8
	Rule           uint8
	Reporter       common.Address
	Offender       common.Address
	RawProof       []byte
	Id             *big.Int
	Block          *big.Int
	Epoch          *big.Int
	ReportingBlock *big.Int
	MessageHash    *big.Int
}

// AutonityCommitteeMember is an auto generated low-level Go binding around an user-defined struct.
type AutonityCommitteeMember struct {
	Addr         common.Address
	VotingPower  *big.Int
	ConsensusKey []byte
}

// AutonityConfig is an auto generated low-level Go binding around an user-defined struct.
type AutonityConfig struct {
	Policy          AutonityPolicy
	Contracts       AutonityContracts
	Protocol        AutonityProtocol
	ContractVersion *big.Int
}

// AutonityContracts is an auto generated low-level Go binding around an user-defined struct.
type AutonityContracts struct {
	AccountabilityContract      common.Address
	OracleContract              common.Address
	AcuContract                 common.Address
	SupplyControlContract       common.Address
	StabilizationContract       common.Address
	UpgradeManagerContract      common.Address
	InflationControllerContract common.Address
}

// AutonityPolicy is an auto generated low-level Go binding around an user-defined struct.
type AutonityPolicy struct {
	TreasuryFee             *big.Int
	MinBaseFee              *big.Int
	DelegationRate          *big.Int
	UnbondingPeriod         *big.Int
	InitialInflationReserve *big.Int
	TreasuryAccount         common.Address
}

// AutonityProtocol is an auto generated low-level Go binding around an user-defined struct.
type AutonityProtocol struct {
	OperatorAccount common.Address
	EpochPeriod     *big.Int
	BlockPeriod     *big.Int
	CommitteeSize   *big.Int
}

// AutonityValidator is an auto generated low-level Go binding around an user-defined struct.
type AutonityValidator struct {
	Treasury                 common.Address
	NodeAddress              common.Address
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
	LiquidContract           common.Address
	LiquidSupply             *big.Int
	RegistrationBlock        *big.Int
	TotalSlashed             *big.Int
	JailReleaseBlock         *big.Int
	ProvableFaultCount       *big.Int
	ConsensusKey             []byte
	State                    uint8
}

// IOracleRoundData is an auto generated low-level Go binding around an user-defined struct.
type IOracleRoundData struct {
	Round     *big.Int
	Price     *big.Int
	Timestamp *big.Int
	Status    *big.Int
}

// InflationControllerParams is an auto generated low-level Go binding around an user-defined struct.
type InflationControllerParams struct {
	IInit  *big.Int
	ITrans *big.Int
	AE     *big.Int
	T      *big.Int
	IPerm  *big.Int
}

// StabilizationConfig is an auto generated low-level Go binding around an user-defined struct.
type StabilizationConfig struct {
	BorrowInterestRate        *big.Int
	LiquidationRatio          *big.Int
	MinCollateralizationRatio *big.Int
	MinDebtRequirement        *big.Int
	TargetPrice               *big.Int
}

// ACUMetaData contains all meta data concerning the ACU contract.
var ACUMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"symbols_\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"quantities_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"scale_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidBasket\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoACUValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"symbols\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"quantities\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"scale\",\"type\":\"uint256\"}],\"name\":\"BasketModified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"Updated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"symbols_\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"quantities_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"scale_\",\"type\":\"uint256\"}],\"name\":\"modifyBasket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"quantities\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"round\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scale\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"scaleFactor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbols\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"update\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"44b4708a": "modifyBasket(string[],uint256[],uint256)",
		"d54d2799": "quantities()",
		"146ca531": "round()",
		"f51e181a": "scale()",
		"683dd191": "scaleFactor()",
		"b3ab15fb": "setOperator(address)",
		"7adbf973": "setOracle(address)",
		"07039ff9": "symbols()",
		"a2e62045": "update()",
		"3fa4f245": "value()",
	},
	Bin: "0x608060405234801562000010575f80fd5b50604051620018433803806200184383398101604081905262000033916200035b565b858580518251146200005857604051634ff799c560e01b815260040160405180910390fd5b5f5b8151811015620000be576001600160ff1b038282815181106200008157620000816200050d565b60200260200101511115620000a957604051634ff799c560e01b815260040160405180910390fd5b80620000b58162000535565b9150506200005a565b508751620000d49060039060208b019062000149565b508651620000ea9060049060208a0190620001a4565b506001869055620000fd86600a6200064b565b6002555050600680546001600160a01b039485166001600160a01b03199182161790915560078054938516938216939093179092556008805491909316911617905550620007b3915050565b828054828255905f5260205f2090810192821562000192579160200282015b82811115620001925782518290620001819082620006eb565b509160200191906001019062000168565b50620001a0929150620001ee565b5090565b828054828255905f5260205f20908101928215620001e0579160200282015b82811115620001e0578251825591602001919060010190620001c3565b50620001a09291506200020e565b80821115620001a0575f62000204828262000224565b50600101620001ee565b5b80821115620001a0575f81556001016200020f565b50805462000232906200065f565b5f825580601f1062000242575050565b601f0160209004905f5260205f20908101906200026091906200020e565b50565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f191681016001600160401b0381118282101715620002a257620002a262000263565b604052919050565b5f6001600160401b03821115620002c557620002c562000263565b5060051b60200190565b5f82601f830112620002df575f80fd5b81516020620002f8620002f283620002aa565b62000277565b82815260059290921b8401810191818101908684111562000317575f80fd5b8286015b848110156200033457805183529183019183016200031b565b509695505050505050565b80516001600160a01b038116811462000356575f80fd5b919050565b5f805f805f8060c0878903121562000371575f80fd5b86516001600160401b0381111562000387575f80fd5b8701601f8101891362000398575f80fd5b8051620003a9620002f282620002aa565b808282526020820191508b60208460051b8601011115620003c8575f80fd5b602084015b60208460051b8601018110156200049a5780516001600160401b03811115620003f4575f80fd5b8d603f828801011262000405575f80fd5b858101602001516001600160401b0381111562000426576200042662000263565b6200043b601f8201601f191660200162000277565b8181528f604083858b010101111562000452575f80fd5b5f5b828110156200047957604081858b010101516020828401015260208101905062000454565b505f60208383010152808652505050602083019250602081019050620003cd565b5060208b0151909950925050506001600160401b03811115620004bb575f80fd5b620004c989828a01620002cf565b95505060408701519350620004e1606088016200033f565b9250620004f1608088016200033f565b91506200050160a088016200033f565b90509295509295509295565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b5f6001820162000549576200054962000521565b5060010190565b600181815b808511156200059057815f190482111562000574576200057462000521565b808516156200058257918102915b93841c939080029062000555565b509250929050565b5f82620005a85750600162000645565b81620005b657505f62000645565b8160018114620005cf5760028114620005da57620005fa565b600191505062000645565b60ff841115620005ee57620005ee62000521565b50506001821b62000645565b5060208310610133831016604e8410600b84101617156200061f575081810a62000645565b6200062b838362000550565b805f190482111562000641576200064162000521565b0290505b92915050565b5f62000658838362000598565b9392505050565b600181811c908216806200067457607f821691505b6020821081036200069357634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620006e6575f81815260208120601f850160051c81016020861015620006c15750805b601f850160051c820191505b81811015620006e257828155600101620006cd565b5050505b505050565b81516001600160401b0381111562000707576200070762000263565b6200071f816200071884546200065f565b8462000699565b602080601f83116001811462000755575f84156200073d5750858301515b5f19600386901b1c1916600185901b178555620006e2565b5f85815260208120601f198616915b82811015620007855788860151825594840194600190910190840162000764565b5085821015620007a357878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61108280620007c15f395ff3fe608060405234801561000f575f80fd5b506004361061009b575f3560e01c80637adbf973116100635780637adbf973146100f9578063a2e620451461010c578063b3ab15fb14610124578063d54d279914610137578063f51e181a1461014c575f80fd5b806307039ff91461009f578063146ca531146100bd5780633fa4f245146100d357806344b4708a146100db578063683dd191146100f0575b5f80fd5b6100a7610155565b6040516100b4919061090c565b60405180910390f35b6100c55f5481565b6040519081526020016100b4565b6100c5610229565b6100ee6100e93660046109f5565b610252565b005b6100c560025481565b6100ee610107366004610b1c565b61037d565b6101146103c9565b60405190151581526020016100b4565b6100ee610132366004610b1c565b6106e1565b61013f61072d565b6040516100b49190610b7b565b6100c560015481565b60606003805480602002602001604051908101604052809291908181526020015f905b82821015610220578382905f5260205f2001805461019590610b8d565b80601f01602080910402602001604051908101604052809291908181526020018280546101c190610b8d565b801561020c5780601f106101e35761010080835404028352916020019161020c565b820191905f5260205f20905b8154815290600101906020018083116101ef57829003601f168201915b505050505081526020019060010190610178565b50505050905090565b5f80545f0361024b57604051631d3e00bb60e11b815260040160405180910390fd5b5060055490565b8282805182511461027657604051634ff799c560e01b815260040160405180910390fd5b5f5b81518110156102d4576001600160ff1b0382828151811061029b5761029b610bc5565b602002602001015111156102c257604051634ff799c560e01b815260040160405180910390fd5b806102cc81610bed565b915050610278565b506007546001600160a01b031633146102ff576040516282b42960e81b815260040160405180910390fd5b8451610312906003906020880190610783565b5083516103269060049060208701906107d7565b50600183905561033783600a610ce7565b6002556040517fdbdcd10543a20811a4a332247f28d03b22686d3281043de35824a06075c06c099061036e90879087908790610cf2565b60405180910390a15050505050565b6006546001600160a01b031633146103a7576040516282b42960e81b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0392909216919091179055565b6006545f906001600160a01b031633146103f5576040516282b42960e81b815260040160405180910390fd5b5f600160085f9054906101000a90046001600160a01b03166001600160a01b0316639f8743f76040518163ffffffff1660e01b8152600401602060405180830381865afa158015610448573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061046c9190610d27565b6104769190610d3e565b9050805f5410610487575f91505090565b5f8060085f9054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156104d9573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104fd9190610d27565b90505f5b60035481101561067d57604051661554d10b5554d160ca1b60208201525f90602701604051602081830303815290604052805190602001206003838154811061054c5761054c610bc5565b905f5260205f20016040516020016105649190610d51565b6040516020818303038152906040528051906020012003610586575081610636565b600854600380545f926001600160a01b031691633c8510fd91899190879081106105b2576105b2610bc5565b905f5260205f20016040518363ffffffff1660e01b81526004016105d7929190610dc3565b608060405180830381865afa1580156105f2573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106169190610e50565b905080606001515f1461062f575f965050505050505090565b6020015190505b6004828154811061064957610649610bc5565b905f5260205f2001548161065d9190610eb4565b6106679085610ee3565b935050808061067590610bed565b915050610501565b506106888183610f0a565b60058190555f849055604080514381524260208201528082018690526060810192909252517f23f161ca67071b3e902d4fa7afade82672c6160677e89d373a830145bdda6d269181900360800190a16001935050505090565b6006546001600160a01b0316331461070b576040516282b42960e81b815260040160405180910390fd5b600780546001600160a01b0319166001600160a01b0392909216919091179055565b6060600480548060200260200160405190810160405280929190818152602001828054801561077957602002820191905f5260205f20905b815481526020019060010190808311610765575b5050505050905090565b828054828255905f5260205f209081019282156107c7579160200282015b828111156107c757825182906107b79082610f90565b50916020019190600101906107a1565b506107d392915061081c565b5090565b828054828255905f5260205f20908101928215610810579160200282015b828111156108105782518255916020019190600101906107f5565b506107d3929150610838565b808211156107d3575f61082f828261084c565b5060010161081c565b5b808211156107d3575f8155600101610839565b50805461085890610b8d565b5f825580601f10610867575050565b601f0160209004905f5260205f20908101906108839190610838565b50565b5f82825180855260208086019550808260051b8401018186015f805b858110156108fe57601f1980888603018b5283518051808752845b818110156108d8578281018901518882018a015288016108bd565b5086810188018590529b87019b601f0190911690940185019350918401916001016108a2565b509198975050505050505050565b602081525f61091e6020830184610886565b9392505050565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561096257610962610925565b604052919050565b5f67ffffffffffffffff82111561098357610983610925565b5060051b60200190565b5f82601f83011261099c575f80fd5b813560206109b16109ac8361096a565b610939565b82815260059290921b840181019181810190868411156109cf575f80fd5b8286015b848110156109ea57803583529183019183016109d3565b509695505050505050565b5f805f60608486031215610a07575f80fd5b833567ffffffffffffffff80821115610a1e575f80fd5b818601915086601f830112610a31575f80fd5b81356020610a416109ac8361096a565b82815260059290921b8401810191818101908a841115610a5f575f80fd5b8286015b84811015610ae857803586811115610a7a575f8081fd5b8701603f81018d13610a8b575f8081fd5b84810135604088821115610aa157610aa1610925565b610ab3601f8301601f19168801610939565b8281528f82848601011115610ac7575f8081fd5b82828501898301375f92810188019290925250845250918301918301610a63565b5097505087013592505080821115610afe575f80fd5b50610b0b8682870161098d565b925050604084013590509250925092565b5f60208284031215610b2c575f80fd5b81356001600160a01b038116811461091e575f80fd5b5f8151808452602080850194508084015f5b83811015610b7057815187529582019590820190600101610b54565b509495945050505050565b602081525f61091e6020830184610b42565b600181811c90821680610ba157607f821691505b602082108103610bbf57634e487b7160e01b5f52602260045260245ffd5b50919050565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b5f60018201610bfe57610bfe610bd9565b5060010190565b600181815b80851115610c3f57815f1904821115610c2557610c25610bd9565b80851615610c3257918102915b93841c9390800290610c0a565b509250929050565b5f82610c5557506001610ce1565b81610c6157505f610ce1565b8160018114610c775760028114610c8157610c9d565b6001915050610ce1565b60ff841115610c9257610c92610bd9565b50506001821b610ce1565b5060208310610133831016604e8410600b8410161715610cc0575081810a610ce1565b610cca8383610c05565b805f1904821115610cdd57610cdd610bd9565b0290505b92915050565b5f61091e8383610c47565b606081525f610d046060830186610886565b8281036020840152610d168186610b42565b915050826040830152949350505050565b5f60208284031215610d37575f80fd5b5051919050565b81810381811115610ce157610ce1610bd9565b5f808354610d5e81610b8d565b60018281168015610d765760018114610d8b57610db7565b60ff1984168752821515830287019450610db7565b875f526020805f205f5b85811015610dae5781548a820152908401908201610d95565b50505082870194505b50929695505050505050565b8281525f60206040818401525f8454610ddb81610b8d565b806040870152606060018084165f8114610dfc5760018114610e1657610e41565b60ff1985168984015283151560051b890183019550610e41565b895f52865f205f5b85811015610e395781548b8201860152908301908801610e1e565b8a0184019650505b50939998505050505050505050565b5f60808284031215610e60575f80fd5b6040516080810181811067ffffffffffffffff82111715610e8357610e83610925565b8060405250825181526020830151602082015260408301516040820152606083015160608201528091505092915050565b8082025f8212600160ff1b84141615610ecf57610ecf610bd9565b8181058314821517610ce157610ce1610bd9565b8082018281125f831280158216821582161715610f0257610f02610bd9565b505092915050565b5f82610f2457634e487b7160e01b5f52601260045260245ffd5b600160ff1b82145f1984141615610f3d57610f3d610bd9565b500590565b601f821115610f8b575f81815260208120601f850160051c81016020861015610f685750805b601f850160051c820191505b81811015610f8757828155600101610f74565b5050505b505050565b815167ffffffffffffffff811115610faa57610faa610925565b610fbe81610fb88454610b8d565b84610f42565b602080601f831160018114610ff1575f8415610fda5750858301515b5f19600386901b1c1916600185901b178555610f87565b5f85815260208120601f198616915b8281101561101f57888601518255948401946001909101908401611000565b508582101561103c57878501515f19600388901b60f8161c191681555b5050505050600190811b0190555056fea26469706673582212209f40da3fbd9c0f2eb11fe36ef9fe2b55db72c911314ad8683252eef825700df264736f6c63430008150033",
}

// ACUABI is the input ABI used to generate the binding from.
// Deprecated: Use ACUMetaData.ABI instead.
var ACUABI = ACUMetaData.ABI

// Deprecated: Use ACUMetaData.Sigs instead.
// ACUFuncSigs maps the 4-byte function signature to its string representation.
var ACUFuncSigs = ACUMetaData.Sigs

// ACUBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ACUMetaData.Bin instead.
var ACUBin = ACUMetaData.Bin

// DeployACU deploys a new Ethereum contract, binding an instance of ACU to it.
func (r *runner) deployACU(opts *runOptions, symbols_ []string, quantities_ []*big.Int, scale_ *big.Int, autonity common.Address, operator common.Address, oracle common.Address) (common.Address, uint64, *ACU, error) {
	parsed, err := ACUMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(ACUBin), symbols_, quantities_, scale_, autonity, operator, oracle)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &ACU{ACUCaller: ACUCaller{contract: c}, ACUTransactor: ACUTransactor{contract: c}, ACUFilterer: ACUFilterer{contract: c}}, nil
}

// ACU is an auto generated Go binding around an Ethereum contract.
type ACU struct {
	ACUCaller     // Read-only binding to the contract
	ACUTransactor // Write-only binding to the contract
	ACUFilterer   // Log filterer for contract events
}

// ACUCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type ACUCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// ACUTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ACUTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// ACUFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ACUFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// Quantities is a free data retrieval call binding the contract method 0xd54d2799.
//
// Solidity: function quantities() view returns(uint256[])
func (_ACU *ACUCaller) Quantities(opts *runOptions) ([]*big.Int, uint64, error) {
	out, consumed, err := _ACU.call(opts, "quantities")

	if err != nil {
		return *new([]*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	return out0, consumed, err

}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_ACU *ACUCaller) Round(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _ACU.call(opts, "round")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Scale is a free data retrieval call binding the contract method 0xf51e181a.
//
// Solidity: function scale() view returns(uint256)
func (_ACU *ACUCaller) Scale(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _ACU.call(opts, "scale")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// ScaleFactor is a free data retrieval call binding the contract method 0x683dd191.
//
// Solidity: function scaleFactor() view returns(uint256)
func (_ACU *ACUCaller) ScaleFactor(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _ACU.call(opts, "scaleFactor")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Symbols is a free data retrieval call binding the contract method 0x07039ff9.
//
// Solidity: function symbols() view returns(string[])
func (_ACU *ACUCaller) Symbols(opts *runOptions) ([]string, uint64, error) {
	out, consumed, err := _ACU.call(opts, "symbols")

	if err != nil {
		return *new([]string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, consumed, err

}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(int256)
func (_ACU *ACUCaller) Value(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _ACU.call(opts, "value")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// ModifyBasket is a paid mutator transaction binding the contract method 0x44b4708a.
//
// Solidity: function modifyBasket(string[] symbols_, uint256[] quantities_, uint256 scale_) returns()
func (_ACU *ACUTransactor) ModifyBasket(opts *runOptions, symbols_ []string, quantities_ []*big.Int, scale_ *big.Int) (uint64, error) {
	_, consumed, err := _ACU.call(opts, "modifyBasket", symbols_, quantities_, scale_)
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ACU *ACUTransactor) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _ACU.call(opts, "setOperator", operator)
	return consumed, err
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_ACU *ACUTransactor) SetOracle(opts *runOptions, oracle common.Address) (uint64, error) {
	_, consumed, err := _ACU.call(opts, "setOracle", oracle)
	return consumed, err
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_ACU *ACUTransactor) Update(opts *runOptions) (uint64, error) {
	_, consumed, err := _ACU.call(opts, "update")
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// ACUBasketModifiedIterator is returned from FilterBasketModified and is used to iterate over the raw logs and unpacked data for BasketModified events raised by the ACU contract.
		type ACUBasketModifiedIterator struct {
			Event *ACUBasketModified // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *ACUBasketModifiedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(ACUBasketModified)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(ACUBasketModified)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *ACUBasketModifiedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *ACUBasketModifiedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// ACUBasketModified represents a BasketModified event raised by the ACU contract.
		type ACUBasketModified struct {
			Symbols []string;
			Quantities []*big.Int;
			Scale *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterBasketModified is a free log retrieval operation binding the contract event 0xdbdcd10543a20811a4a332247f28d03b22686d3281043de35824a06075c06c09.
		//
		// Solidity: event BasketModified(string[] symbols, uint256[] quantities, uint256 scale)
 		func (_ACU *ACUFilterer) FilterBasketModified(opts *bind.FilterOpts) (*ACUBasketModifiedIterator, error) {





			logs, sub, err := _ACU.contract.FilterLogs(opts, "BasketModified")
			if err != nil {
				return nil, err
			}
			return &ACUBasketModifiedIterator{contract: _ACU.contract, event: "BasketModified", logs: logs, sub: sub}, nil
 		}

		// WatchBasketModified is a free log subscription operation binding the contract event 0xdbdcd10543a20811a4a332247f28d03b22686d3281043de35824a06075c06c09.
		//
		// Solidity: event BasketModified(string[] symbols, uint256[] quantities, uint256 scale)
		func (_ACU *ACUFilterer) WatchBasketModified(opts *bind.WatchOpts, sink chan<- *ACUBasketModified) (event.Subscription, error) {





			logs, sub, err := _ACU.contract.WatchLogs(opts, "BasketModified")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(ACUBasketModified)
						if err := _ACU.contract.UnpackLog(event, "BasketModified", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseBasketModified is a log parse operation binding the contract event 0xdbdcd10543a20811a4a332247f28d03b22686d3281043de35824a06075c06c09.
		//
		// Solidity: event BasketModified(string[] symbols, uint256[] quantities, uint256 scale)
		func (_ACU *ACUFilterer) ParseBasketModified(log types.Log) (*ACUBasketModified, error) {
			event := new(ACUBasketModified)
			if err := _ACU.contract.UnpackLog(event, "BasketModified", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// ACUUpdatedIterator is returned from FilterUpdated and is used to iterate over the raw logs and unpacked data for Updated events raised by the ACU contract.
		type ACUUpdatedIterator struct {
			Event *ACUUpdated // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *ACUUpdatedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(ACUUpdated)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(ACUUpdated)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *ACUUpdatedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *ACUUpdatedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// ACUUpdated represents a Updated event raised by the ACU contract.
		type ACUUpdated struct {
			Height *big.Int;
			Timestamp *big.Int;
			Round *big.Int;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterUpdated is a free log retrieval operation binding the contract event 0x23f161ca67071b3e902d4fa7afade82672c6160677e89d373a830145bdda6d26.
		//
		// Solidity: event Updated(uint256 height, uint256 timestamp, uint256 round, int256 value)
 		func (_ACU *ACUFilterer) FilterUpdated(opts *bind.FilterOpts) (*ACUUpdatedIterator, error) {






			logs, sub, err := _ACU.contract.FilterLogs(opts, "Updated")
			if err != nil {
				return nil, err
			}
			return &ACUUpdatedIterator{contract: _ACU.contract, event: "Updated", logs: logs, sub: sub}, nil
 		}

		// WatchUpdated is a free log subscription operation binding the contract event 0x23f161ca67071b3e902d4fa7afade82672c6160677e89d373a830145bdda6d26.
		//
		// Solidity: event Updated(uint256 height, uint256 timestamp, uint256 round, int256 value)
		func (_ACU *ACUFilterer) WatchUpdated(opts *bind.WatchOpts, sink chan<- *ACUUpdated) (event.Subscription, error) {






			logs, sub, err := _ACU.contract.WatchLogs(opts, "Updated")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(ACUUpdated)
						if err := _ACU.contract.UnpackLog(event, "Updated", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseUpdated is a log parse operation binding the contract event 0x23f161ca67071b3e902d4fa7afade82672c6160677e89d373a830145bdda6d26.
		//
		// Solidity: event Updated(uint256 height, uint256 timestamp, uint256 round, int256 value)
		func (_ACU *ACUFilterer) ParseUpdated(log types.Log) (*ACUUpdated, error) {
			event := new(ACUUpdated)
			if err := _ACU.contract.UnpackLog(event, "Updated", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// AccountabilityMetaData contains all meta data concerning the Accountability contract.
var AccountabilityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_autonity\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"innocenceProofSubmissionWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateLow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateMid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collusionFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"historyFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashingRatePrecision\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"InnocenceProven\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewAccusation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewFaultProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"releaseBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isJailbound\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"SlashingEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"beneficiaries\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"_rule\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"canAccuse\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_deadline\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"_rule\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"canSlash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"innocenceProofSubmissionWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateLow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateMid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collusionFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"historyFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashingRatePrecision\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"events\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"chunkId\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_epochEnd\",\"type\":\"bool\"}],\"name\":\"finalize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_val\",\"type\":\"address\"}],\"name\":\"getValidatorAccusation\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"chunkId\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_val\",\"type\":\"address\"}],\"name\":\"getValidatorFaults\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"chunkId\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"chunkId\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event\",\"name\":\"_event\",\"type\":\"tuple\"}],\"name\":\"handleEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newPeriod\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"slashingHistory\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"01567739": "beneficiaries(address)",
		"7ccecadd": "canAccuse(address,uint8,uint256)",
		"4108a95a": "canSlash(address,uint8,uint256)",
		"79502c55": "config()",
		"1de9d9b6": "distributeRewards(address)",
		"b5b7a184": "epochPeriod()",
		"0b791430": "events(uint256)",
		"6c9789b0": "finalize(bool)",
		"9cb22b06": "getValidatorAccusation(address)",
		"bebaa8fc": "getValidatorFaults(address)",
		"c50d21f0": "handleEvent((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256))",
		"6b5f444c": "setEpochPeriod(uint256)",
		"e7bb0b52": "slashingHistory(address,uint256)",
	},
	Bin: "0x60806040525f60115534801562000014575f80fd5b5060405162003c6538038062003c658339810160408190526200003791620000f2565b600180546001600160a01b0319166001600160a01b03841690811790915560408051636fd8d26960e11b8152905163dfb1a4d2916004808201926020929091908290030181865afa1580156200008f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620000b59190620001b4565b5f558051600355602081015160045560408101516005556060810151600655608081015160075560a081015160085560c0015160095550620001cc565b5f8082840361010081121562000106575f80fd5b83516001600160a01b03811681146200011d575f80fd5b925060e0601f198201121562000131575f80fd5b5060405160e081016001600160401b03811182821017156200016157634e487b7160e01b5f52604160045260245ffd5b80604052506020840151815260408401516020820152606084015160408201526080840151606082015260a0840151608082015260c084015160a082015260e084015160c0820152809150509250929050565b5f60208284031215620001c5575f80fd5b5051919050565b613a8b80620001da5f395ff3fe6080604052600436106100bf575f3560e01c806379502c551161007c578063b5b7a18411610057578063b5b7a18414610292578063bebaa8fc146102b4578063c50d21f0146102e0578063e7bb0b52146102ff575f80fd5b806379502c55146101cd5780637ccecadd146102305780639cb22b0614610266575f80fd5b806301567739146100c35780630b791430146101145780631de9d9b61461014b5780634108a95a146101605780636b5f444c1461018f5780636c9789b0146101ae575b5f80fd5b3480156100ce575f80fd5b506100f76100dd366004612ea6565b600a6020525f90815260409020546001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b34801561011f575f80fd5b5061013361012e366004612ec8565b610335565b60405161010b9c9b9a99989796959493929190612f64565b61015e610159366004612ea6565b610444565b005b34801561016b575f80fd5b5061017f61017a366004612ffd565b610644565b604051901515815260200161010b565b34801561019a575f80fd5b5061015e6101a9366004612ec8565b6106ef565b3480156101b9575f80fd5b5061015e6101c8366004613038565b61071d565b3480156101d8575f80fd5b506003546004546005546006546007546008546009546101fb9695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e00161010b565b34801561023b575f80fd5b5061024f61024a366004612ffd565b610760565b60408051921515835260208301919091520161010b565b348015610271575f80fd5b50610285610280366004612ea6565b6108a6565b60405161010b919061312c565b34801561029d575f80fd5b506102a65f5481565b60405190815260200161010b565b3480156102bf575f80fd5b506102d36102ce366004612ea6565b610ac8565b60405161010b919061313e565b3480156102eb575f80fd5b5061015e6102fa3660046132c8565b610d5e565b34801561030a575f80fd5b506102a66103193660046133cf565b600e60209081525f928352604080842090915290825290205481565b60028181548110610344575f80fd5b5f91825260209091206008909102018054600182015460028301805460ff8085169650610100850481169562010000860482169563010000008104909216946001600160a01b03600160201b9093048316949216929091906103a5906133f9565b80601f01602080910402602001604051908101604052809291908181526020018280546103d1906133f9565b801561041c5780601f106103f35761010080835404028352916020019161041c565b820191905f5260205f20905b8154815290600101906020018083116103ff57829003601f168201915b505050505090806003015490806004015490806005015490806006015490806007015490508c565b6001546001600160a01b031633146104775760405162461bcd60e51b815260040161046e90613431565b60405180910390fd5b6001546001600160a01b038281165f908152600a6020526040808220549051630c825d9760e11b8152908316600482015290929190911690631904bb2e906024015f60405180830381865afa1580156104d2573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526104f991908101906134d8565b516040519091505f906001600160a01b038316906108fc90349084818181858888f193505050503d805f811461054a576040519150601f19603f3d011682016040523d82523d5f602084013e61054f565b606091505b505090508061061c5760015f9054906101000a90046001600160a01b03166001600160a01b031663f7866ee36040518163ffffffff1660e01b8152600401602060405180830381865afa1580156105a8573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105cc9190613650565b6001600160a01b0316346040515f6040518083038185875af1925050503d805f8114610613576040519150601f19603f3d011682016040523d82523d5f602084013e610618565b606091505b5050505b50506001600160a01b03165f908152600a6020526040902080546001600160a01b0319169055565b5f8061064f846110c0565b6001546040516396b477cb60e01b8152600481018690529192505f916001600160a01b03909116906396b477cb90602401602060405180830381865afa15801561069b573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106bf919061366b565b6001600160a01b0387165f908152600e602090815260408083209383529290522054919091109150509392505050565b6001546001600160a01b031633146107195760405162461bcd60e51b815260040161046e90613431565b5f55565b6001546001600160a01b031633146107475760405162461bcd60e51b815260040161046e90613431565b61074f611129565b801561075d5761075d611475565b50565b5f805f61076c856110c0565b6001546040516396b477cb60e01b8152600481018790529192505f916001600160a01b03909116906396b477cb90602401602060405180830381865afa1580156107b8573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107dc919061366b565b6001600160a01b0388165f908152600e602090815260408083208484529091529020549091508211610813575f93505f925061089c565b6001600160a01b0387165f908152600c602052604090205415610894576001600160a01b0387165f908152600c602052604081205460029061085790600190613696565b81548110610867576108676136a9565b905f5260205f20906008020190505f945060035f0154816004015461088c91906136bd565b93505061089c565b600193505f92505b5050935093915050565b6108ae612dc7565b6001600160a01b0382165f908152600c60205260409020546109025760405162461bcd60e51b815260206004820152600d60248201526c37379030b1b1bab9b0ba34b7b760991b604482015260640161046e565b6001600160a01b0382165f908152600c602052604090205460029061092990600190613696565b81548110610939576109396136a9565b5f91825260209182902060408051610180810182526008909302909101805460ff80821685526101008204811695850195909552929390929184019162010000900416600281111561098d5761098d612edf565b600281111561099e5761099e612edf565b815281546020909101906301000000900460ff1660098111156109c3576109c3612edf565b60098111156109d4576109d4612edf565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191610a0f906133f9565b80601f0160208091040260200160405190810160405280929190818152602001828054610a3b906133f9565b8015610a865780601f10610a5d57610100808354040283529160200191610a86565b820191905f5260205f20905b815481529060010190602001808311610a6957829003601f168201915b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820154815250509050919050565b6001600160a01b0381165f908152600b60205260408120546060919067ffffffffffffffff811115610afc57610afc61319e565b604051908082528060200260200182016040528015610b3557816020015b610b22612dc7565b815260200190600190039081610b1a5790505b5090505f5b6001600160a01b0384165f908152600b6020526040902054811015610d57576001600160a01b0384165f908152600b6020526040902080546002919083908110610b8657610b866136a9565b905f5260205f20015481548110610b9f57610b9f6136a9565b5f91825260209182902060408051610180810182526008909302909101805460ff808216855261010082048116958501959095529293909291840191620100009004166002811115610bf357610bf3612edf565b6002811115610c0457610c04612edf565b815281546020909101906301000000900460ff166009811115610c2957610c29612edf565b6009811115610c3a57610c3a612edf565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191610c75906133f9565b80601f0160208091040260200160405190810160405280929190818152602001828054610ca1906133f9565b8015610cec5780601f10610cc357610100808354040283529160200191610cec565b820191905f5260205f20905b815481529060010190602001808311610ccf57829003601f168201915b5050505050815260200160038201548152602001600482015481526020016005820154815260200160068201548152602001600782015481525050828281518110610d3957610d396136a9565b60200260200101819052508080610d4f906136d0565b915050610b3a565b5092915050565b600154604051630c825d9760e11b81523360048201525f916001600160a01b031690631904bb2e906024015f60405180830381865afa158015610da3573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f19168201604052610dca91908101906134d8565b60208101519091506001600160a01b03163314610e3f5760405162461bcd60e51b815260206004820152602d60248201527f66756e6374696f6e207265737472696374656420746f2061207265676973746560448201526c3932b2103b30b634b230ba37b960991b606482015260840161046e565b60808201516001600160a01b03163314610e9b5760405162461bcd60e51b815260206004820152601d60248201527f6576656e74207265706f72746572206d7573742062652063616c6c6572000000604482015260640161046e565b6001825f015160ff16111561104b575f610eb483611747565b905080610ec057505050565b335f908152600d6020908152604091829020825161018081018452815460ff80821683526101008204811694830194909452909391929184019162010000909104166002811115610f1357610f13612edf565b6002811115610f2457610f24612edf565b815281546020909101906301000000900460ff166009811115610f4957610f49612edf565b6009811115610f5a57610f5a612edf565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191610f95906133f9565b80601f0160208091040260200160405190810160405280929190818152602001828054610fc1906133f9565b801561100c5780601f10610fe35761010080835404028352916020019161100c565b820191905f5260205f20905b815481529060010190602001808311610fef57829003601f168201915b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820154815250509250505b5f8260400151600281111561106257611062612edf565b036110745761107082611978565b5050565b60018260400151600281111561108c5761108c612edf565b0361109a5761107082611b77565b6002826040015160028111156110b2576110b2612edf565b036110705761107082611ce0565b5f60098260098111156110d5576110d5612edf565b036110e35760025b92915050565b5f8260098111156110f6576110f6612edf565b036111025760026110dd565b600182600981111561111657611116612edf565b036111225760026110dd565b60026110dd565b6011545b601054811015611470575f6010828154811061114b5761114b6136a9565b905f5260205f2001549050805f03611163575061145e565b61116e600182613696565b90505f60028281548110611184576111846136a9565b5f91825260209182902060408051610180810182526008909302909101805460ff8082168552610100820481169585019590955292939092918401916201000090041660028111156111d8576111d8612edf565b60028111156111e9576111e9612edf565b815281546020909101906301000000900460ff16600981111561120e5761120e612edf565b600981111561121f5761121f612edf565b815281546001600160a01b03600160201b90910481166020830152600183015416604082015260028201805460609092019161125a906133f9565b80601f0160208091040260200160405190810160405280929190818152602001828054611286906133f9565b80156112d15780601f106112a8576101008083540402835291602001916112d1565b820191905f5260205f20905b8154815290600101906020018083116112b457829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815260200160078201548152505090504360035f015482610140015161132491906136bd565b1115611331575050601155565b60a08101516001600160a01b03165f908152600c60205260408120819055606082015161135d906110c0565b60a08301516001600160a01b03165f908152600e602090815260408083206101208701518452909152902054909150811161139a5750505061145e565b60a0820180516001600160a01b039081165f908152600e6020908152604080832061012088015184528252808320869055845184168352600b825280832080546001808201835591855283852001899055600f805491820181559093527f8d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac80290920187905592518151858152938401879052909116917f6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f910160405180910390a25050505b80611468816136d0565b91505061112d565b601155565b5f8060015f9054906101000a90046001600160a01b03166001600160a01b031663c9d97af46040518163ffffffff1660e01b8152600401602060405180830381865afa1580156114c7573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906114eb919061366b565b90505f5b600f5481101561155d57816002600f838154811061150f5761150f6136a9565b905f5260205f20015481548110611528576115286136a9565b905f5260205f209060080201600501540361154b576115486001846136bd565b92505b80611555816136d0565b9150506114ef565b505f5b600f5481101561173b576117296002600f8381548110611582576115826136a9565b905f5260205f2001548154811061159b5761159b6136a9565b5f91825260209182902060408051610180810182526008909302909101805460ff8082168552610100820481169585019590955292939092918401916201000090041660028111156115ef576115ef612edf565b600281111561160057611600612edf565b815281546020909101906301000000900460ff16600981111561162557611625612edf565b600981111561163657611636612edf565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191611671906133f9565b80601f016020809104026020016040519081016040528092919081815260200182805461169d906133f9565b80156116e85780601f106116bf576101008083540402835291602001916116e8565b820191905f5260205f20905b8154815290600101906020018083116116cb57829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815260200160078201548152505084611e1d565b80611733816136d0565b915050611560565b50611070600f5f612e31565b5f816020015160ff165f0361187857335f908152600d6020908152604091829020845181549286015160ff9081166101000261ffff1990941691161791909117808255918401518492829062ff00001916620100008360028111156117ae576117ae612edf565b021790555060608201518154829063ff000000191663010000008360098111156117da576117da612edf565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c082015160028201906118389082613732565b5060e0820151600382015561010082015160048201556101208201516005820155610140820151600682015561016090910151600790910155505f919050565b602080830151335f908152600d90925260409091205460ff918216916118a59161010090041660016137ee565b60ff16146118f55760405162461bcd60e51b815260206004820152601960248201527f6368756e6b73206d75737420626520636f6e746967756f757300000000000000604482015260640161046e565b335f908152600d6020526040902060c08301516119159160020190612395565b335f908152600d60205260409020805460019190829061193e908290610100900460ff166137ee565b92506101000a81548160ff021916908360ff160217905550815f015160ff168260200151600161196e91906137ee565b60ff161492915050565b5f805f805f61198c60fe8760c001516124db565b94509450945094509450846119e35760405162461bcd60e51b815260206004820152601960248201527f6661696c65642070726f6f6620766572696669636174696f6e00000000000000604482015260640161046e565b8560a001516001600160a01b0316846001600160a01b031614611a185760405162461bcd60e51b815260040161046e90613807565b85606001516009811115611a2e57611a2e612edf565b8314611a4c5760405162461bcd60e51b815260040161046e90613832565b438210611a945760405162461bcd60e51b815260206004820152601660248201527563616e277420626520696e207468652066757475726560501b604482015260640161046e565b5f8211611ad95760405162461bcd60e51b815260206004820152601360248201527263616e27742062652061742067656e6573697360681b604482015260640161046e565b6001546040516396b477cb60e01b8152600481018490525f916001600160a01b0316906396b477cb90602401602060405180830381865afa158015611b20573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611b44919061366b565b610100880184905261012088018190524361014089015261016088018390529050611b6e8761253e565b50505050505050565b5f805f805f611b8b60fc8760c001516124db565b9450945094509450945084611be25760405162461bcd60e51b815260206004820152601e60248201527f6661696c65642061636375736174696f6e20766572696669636174696f6e0000604482015260640161046e565b8560a001516001600160a01b0316846001600160a01b031614611c175760405162461bcd60e51b815260040161046e90613807565b85606001516009811115611c2d57611c2d612edf565b8314611c4b5760405162461bcd60e51b815260040161046e90613832565b6001546040516396b477cb60e01b8152600481018490525f916001600160a01b0316906396b477cb90602401602060405180830381865afa158015611c92573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611cb6919061366b565b610100880184905261012088018190524361014089015261016088018390529050611b6e876127be565b5f805f805f611cf460fd8760c001516124db565b9450945094509450945084611d4b5760405162461bcd60e51b815260206004820152601d60248201527f6661696c656420696e6e6f63656e636520766572696669636174696f6e000000604482015260640161046e565b8560a001516001600160a01b0316846001600160a01b031614611d805760405162461bcd60e51b815260040161046e90613807565b85606001516009811115611d9657611d96612edf565b8314611db45760405162461bcd60e51b815260040161046e90613832565b438210611dfc5760405162461bcd60e51b815260206004820152601660248201527563616e277420626520696e207468652066757475726560501b604482015260640161046e565b61010086018290526101608601819052611e1586612a6f565b505050505050565b60015460a0830151604051630c825d9760e11b81526001600160a01b0391821660048201525f929190911690631904bb2e906024015f60405180830381865afa158015611e6c573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f19168201604052611e9391908101906134d8565b608084015160a08501516001600160a01b039081165f908152600a6020526040902080546001600160a01b03191691909216179055905060038161026001516003811115611ee357611ee3612edf565b03611eed57505050565b5f611f03611efe85606001516110c0565b612d6f565b610220830151600754919250905f90611f1c908361385c565b600654611f29908761385c565b611f3390856136bd565b611f3d91906136bd565b600954909150811115611f4f57506009545b5f8461012001518560c001518660a00151611f6a91906136bd565b611f7491906136bd565b6009549091505f90611f86838561385c565b611f909190613873565b90505f81118015611fa057508181145b156120c8575f60a087018190526101008701819052610120870181905260c08701526101e086018051829190611fd79083906136bd565b9052506102208601805160019190611ff09083906136bd565b90525060036102608701525f6102008701526001546040516301adf0b760e51b81526001600160a01b03909116906335be16e0906120329089906004016138a2565b5f604051808303815f87803b158015612049575f80fd5b505af115801561205b573d5f803e3d5ffd5b5050505060208681015160e08a0151604080516001600160a01b0390931683529282018490525f8284015260016060830152608082015290517f6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f99181900360a00190a15050505050505050565b610120860151819081116120f4578087610120018181516120e99190613696565b9052505f905061210e565b6101208701516121049082613696565b5f61012089015290505b8015612189578087610100015110612155578087610100018181516121339190613696565b90525060a08701805182919061214a908390613696565b9052505f9050612189565b6101008701516121659082613696565b90508661010001518760a00181815161217e9190613696565b9052505f6101008801525b5f811180156121aa57505f8760a001518860c001516121a891906136bd565b115b15612254575f8760a001518860c001516121c491906136bd565b60c08901516121d3908461385c565b6121dd9190613873565b90505f8860a001518960c001516121f491906136bd565b60a08a0151612203908561385c565b61220d9190613873565b9050818960c0018181516122219190613696565b90525060a089018051829190612238908390613696565b90525061224581836136bd565b61224f9084613696565b925050505b61225e8183613696565b915081876101e00181815161227391906136bd565b905250610220870180516001919061228c9083906136bd565b9052505f546102208801516008546122a4919061385c565b6122ae919061385c565b6122b890436136bd565b61020088015260026102608801526001546040516301adf0b760e51b81526001600160a01b03909116906335be16e0906122f6908a906004016138a2565b5f604051808303815f87803b15801561230d575f80fd5b505af115801561231f573d5f803e3d5ffd5b5050506020808901516102008a015160e08d0151604080516001600160a01b039094168452938301879052928201525f606082015260808101919091527f6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9915060a00160405180910390a1505050505050505050565b81546002600180831615610100020382160482518082016020811060208410016002811461243d576001811461246257865f526020840460205f2001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b80821015612427578151845560018401935060208201915061240a565b815191036101000a908190040290915550611b6e565b60028302826020036101000a846020036101000a602089015104020185018755611b6e565b865f526020840460205f2001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b808210156124c157815184556001840193506020820191506124a4565b815191036101000a90819004029091555050505050505050565b5f805f805f80865160206124ef91906136bd565b90506124f9612e4c565b60a081838a8c5afa612509575f80fd5b805160010361251757600196505b602081015160408201516060830151608090930151989b919a509850909695509350505050565b5f61254c82606001516110c0565b60a08301516001600160a01b03165f908152600e60209081526040808320610120870151845290915290205490915081116125995760405162461bcd60e51b815260040161046e90613a11565b6002805460e084018190526001810182555f829052835160089091027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054602086015160ff9081166101000261ffff1990921693169290921791909117808255604085015185939091839162ff000019909116906201000090849081111561262557612625612edf565b021790555060608201518154829063ff0000001916630100000083600981111561265157612651612edf565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c082015160028201906126af9082613732565b5060e0828101516003830155610100830151600483015561012080840151600584015561014084015160068401556101609093015160079092019190915560a0840180516001600160a01b039081165f908152600b602090815260408083209589018051875460018181018a5598865284862001558051600f8054988901815585527f8d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac80290970196909655845184168352600e82528083209689015183529590528490208590559051915192519116917f6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f916127b291858252602082015260400190565b60405180910390a25050565b60a08101516001600160a01b03165f908152600c6020526040902054156128275760405162461bcd60e51b815260206004820181905260248201527f616c72656164792070726f63657373696e6720616e2061636375736174696f6e604482015260640161046e565b5f61283582606001516110c0565b60a08301516001600160a01b03165f908152600e60209081526040808320610120870151845290915290205490915081116128825760405162461bcd60e51b815260040161046e90613a11565b6002805460e084018190526001810182555f829052835160089091027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054602086015160ff9081166101000261ffff1990921693169290921791909117808255604085015185939091839162ff000019909116906201000090849081111561290e5761290e612edf565b021790555060608201518154829063ff0000001916630100000083600981111561293a5761293a612edf565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c082015160028201906129989082613732565b5060e0828101516003830155610100830151600483015561012083015160058301556101408301516006830155610160909201516007909101558201516129e09060016136bd565b60a08301516001600160a01b03165f908152600c602052604090205560e0820151601090612a0f9060016136bd565b81546001810183555f928352602092839020015560a083015160e084015160408051858152938401919091526001600160a01b03909116917f2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba4035191016127b2565b60a08101516001600160a01b03165f908152600c602052604081205490819003612adb5760405162461bcd60e51b815260206004820152601860248201527f6e6f206173736f6369617465642061636375736174696f6e0000000000000000604482015260640161046e565b81606001516009811115612af157612af1612edf565b6002612afe600184613696565b81548110612b0e57612b0e6136a9565b5f9182526020909120600890910201546301000000900460ff166009811115612b3957612b39612edf565b14612b965760405162461bcd60e51b815260206004820152602760248201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e206044820152661c9d5b19481a5960ca1b606482015260840161046e565b6101008201516002612ba9600184613696565b81548110612bb957612bb96136a9565b905f5260205f2090600802016004015414612c245760405162461bcd60e51b815260206004820152602560248201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e20604482015264626c6f636b60d81b606482015260840161046e565b6101608201516002612c37600184613696565b81548110612c4757612c476136a9565b905f5260205f2090600802016007015414612cb05760405162461bcd60e51b8152602060048201526024808201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e206044820152630d0c2e6d60e31b606482015260840161046e565b6011545b601054811015612d16578160108281548110612cd257612cd26136a9565b905f5260205f20015403612d04575f60108281548110612cf457612cf46136a9565b5f91825260209091200155612d16565b80612d0e816136d0565b915050612cb4565b5060a0820180516001600160a01b039081165f908152600c602090815260408083208390559351935191825292909116917f1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f91016127b2565b5f81612d7d57505060055490565b60018203612d8d57505060055490565b60028203612d9d57505060055490565b60038203612dad57505060055490565b60048203612dbe5750612710919050565b50612710919050565b60408051610180810182525f80825260208201819052909182019081526020015f81526020015f6001600160a01b031681526020015f6001600160a01b03168152602001606081526020015f81526020015f81526020015f81526020015f81526020015f81525090565b5080545f8255905f5260205f209081019061075d9190612e6a565b6040518060a001604052806005906020820280368337509192915050565b5b80821115612e7e575f8155600101612e6b565b5090565b6001600160a01b038116811461075d575f80fd5b8035612ea181612e82565b919050565b5f60208284031215612eb6575f80fd5b8135612ec181612e82565b9392505050565b5f60208284031215612ed8575f80fd5b5035919050565b634e487b7160e01b5f52602160045260245ffd5b60038110612f0357612f03612edf565b9052565b600a8110612f0357612f03612edf565b5f5b83811015612f31578181015183820152602001612f19565b50505f910152565b5f8151808452612f50816020860160208601612f17565b601f01601f19169290920160200192915050565b5f61018060ff8f16835260ff8e166020840152612f84604084018e612ef3565b612f91606084018d612f07565b6001600160a01b038b811660808501528a1660a084015260c08301819052612fbb8184018a612f39565b60e0840198909852505061010081019490945261012084019290925261014083015261016090910152979650505050505050565b8035600a8110612ea1575f80fd5b5f805f6060848603121561300f575f80fd5b833561301a81612e82565b925061302860208501612fef565b9150604084013590509250925092565b5f60208284031215613048575f80fd5b81358015158114612ec1575f80fd5b805160ff1682525f6101806020830151613076602086018260ff169052565b5060408301516130896040860182612ef3565b50606083015161309c6060860182612f07565b5060808301516130b760808601826001600160a01b03169052565b5060a08301516130d260a08601826001600160a01b03169052565b5060c08301518160c08601526130ea82860182612f39565b60e08581015190870152610100808601519087015261012080860151908701526101408086015190870152610160948501519490950193909352509192915050565b602081525f612ec16020830184613057565b5f602080830181845280855180835260408601915060408160051b87010192508387015f5b8281101561319157603f1988860301845261317f858351613057565b94509285019290850190600101613163565b5092979650505050505050565b634e487b7160e01b5f52604160045260245ffd5b604051610180810167ffffffffffffffff811182821017156131d6576131d661319e565b60405290565b604051610280810167ffffffffffffffff811182821017156131d6576131d661319e565b604051601f8201601f1916810167ffffffffffffffff811182821017156132295761322961319e565b604052919050565b803560ff81168114612ea1575f80fd5b803560038110612ea1575f80fd5b5f67ffffffffffffffff8211156132685761326861319e565b50601f01601f191660200190565b5f82601f830112613285575f80fd5b81356132986132938261324f565b613200565b8181528460208386010111156132ac575f80fd5b816020850160208301375f918101602001919091529392505050565b5f602082840312156132d8575f80fd5b813567ffffffffffffffff808211156132ef575f80fd5b908301906101808286031215613303575f80fd5b61330b6131b2565b61331483613231565b815261332260208401613231565b602082015261333360408401613241565b604082015261334460608401612fef565b606082015261335560808401612e96565b608082015261336660a08401612e96565b60a082015260c08301358281111561337c575f80fd5b61338887828601613276565b60c08301525060e083810135908201526101008084013590820152610120808401359082015261014080840135908201526101609283013592810192909252509392505050565b5f80604083850312156133e0575f80fd5b82356133eb81612e82565b946020939093013593505050565b600181811c9082168061340d57607f821691505b60208210810361342b57634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526024908201527f66756e6374696f6e207265737472696374656420746f207468652076616c696460408201526330ba37b960e11b606082015260800190565b8051612ea181612e82565b5f82601f83011261348f575f80fd5b815161349d6132938261324f565b8181528460208386010111156134b1575f80fd5b6134c2826020830160208701612f17565b949350505050565b805160048110612ea1575f80fd5b5f602082840312156134e8575f80fd5b815167ffffffffffffffff808211156134ff575f80fd5b908301906102808286031215613513575f80fd5b61351b6131dc565b61352483613475565b815261353260208401613475565b602082015261354360408401613475565b6040820152606083015182811115613559575f80fd5b61356587828601613480565b6060830152506080830151608082015260a083015160a082015260c083015160c082015260e083015160e08201526101008084015181830152506101208084015181830152506101408084015181830152506101608084015181830152506101806135d1818501613475565b908201526101a083810151908201526101c080840151908201526101e0808401519082015261020080840151908201526102208084015190820152610240808401518381111561361f575f80fd5b61362b88828701613480565b82840152505061026091506136418284016134ca565b91810191909152949350505050565b5f60208284031215613660575f80fd5b8151612ec181612e82565b5f6020828403121561367b575f80fd5b5051919050565b634e487b7160e01b5f52601160045260245ffd5b818103818111156110dd576110dd613682565b634e487b7160e01b5f52603260045260245ffd5b808201808211156110dd576110dd613682565b5f600182016136e1576136e1613682565b5060010190565b601f82111561372d575f81815260208120601f850160051c8101602086101561370e5750805b601f850160051c820191505b81811015611e155782815560010161371a565b505050565b815167ffffffffffffffff81111561374c5761374c61319e565b6137608161375a84546133f9565b846136e8565b602080601f831160018114613793575f841561377c5750858301515b5f19600386901b1c1916600185901b178555611e15565b5f85815260208120601f198616915b828110156137c1578886015182559484019460019091019084016137a2565b50858210156137de57878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b60ff81811683821601908111156110dd576110dd613682565b6020808252601190820152700decccccadcc8cae440dad2e6dac2e8c6d607b1b604082015260600190565b60208082526010908201526f0e4ead8ca40d2c840dad2e6dac2e8c6d60831b604082015260600190565b80820281158282048414176110dd576110dd613682565b5f8261388d57634e487b7160e01b5f52601260045260245ffd5b500490565b60048110612f0357612f03612edf565b602081526138bc6020820183516001600160a01b03169052565b5f60208301516138d760408401826001600160a01b03169052565b5060408301516001600160a01b03811660608401525060608301516102808060808501526139096102a0850183612f39565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a0613990818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f1901610260808801919091529091506139f28483612f39565b935080870151915050613a0782860182613892565b5090949350505050565b60208082526024908201527f616c726561647920736c6173686564206174207468652070726f6f66277320656040820152630e0dec6d60e31b60608201526080019056fea264697066735822122003bcffc82be5f9e5152e75d7799a6d631d5941e2e899e2b8ee3c60e642010af664736f6c63430008150033",
}

// AccountabilityABI is the input ABI used to generate the binding from.
// Deprecated: Use AccountabilityMetaData.ABI instead.
var AccountabilityABI = AccountabilityMetaData.ABI

// Deprecated: Use AccountabilityMetaData.Sigs instead.
// AccountabilityFuncSigs maps the 4-byte function signature to its string representation.
var AccountabilityFuncSigs = AccountabilityMetaData.Sigs

// AccountabilityBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AccountabilityMetaData.Bin instead.
var AccountabilityBin = AccountabilityMetaData.Bin

// DeployAccountability deploys a new Ethereum contract, binding an instance of Accountability to it.
func (r *runner) deployAccountability(opts *runOptions, _autonity common.Address, _config AccountabilityConfig) (common.Address, uint64, *Accountability, error) {
	parsed, err := AccountabilityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(AccountabilityBin), _autonity, _config)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &Accountability{AccountabilityCaller: AccountabilityCaller{contract: c}, AccountabilityTransactor: AccountabilityTransactor{contract: c}, AccountabilityFilterer: AccountabilityFilterer{contract: c}}, nil
}

// Accountability is an auto generated Go binding around an Ethereum contract.
type Accountability struct {
	AccountabilityCaller     // Read-only binding to the contract
	AccountabilityTransactor // Write-only binding to the contract
	AccountabilityFilterer   // Log filterer for contract events
}

// AccountabilityCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type AccountabilityCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// AccountabilityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccountabilityTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// AccountabilityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccountabilityFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// Beneficiaries is a free data retrieval call binding the contract method 0x01567739.
//
// Solidity: function beneficiaries(address ) view returns(address)
func (_Accountability *AccountabilityCaller) Beneficiaries(opts *runOptions, arg0 common.Address) (common.Address, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "beneficiaries", arg0)

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// CanAccuse is a free data retrieval call binding the contract method 0x7ccecadd.
//
// Solidity: function canAccuse(address _offender, uint8 _rule, uint256 _block) view returns(bool _result, uint256 _deadline)
func (_Accountability *AccountabilityCaller) CanAccuse(opts *runOptions, _offender common.Address, _rule uint8, _block *big.Int) (struct {
	Result   bool
	Deadline *big.Int
}, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "canAccuse", _offender, _rule, _block)

	outstruct := new(struct {
		Result   bool
		Deadline *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.Result = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Deadline = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// CanSlash is a free data retrieval call binding the contract method 0x4108a95a.
//
// Solidity: function canSlash(address _offender, uint8 _rule, uint256 _block) view returns(bool)
func (_Accountability *AccountabilityCaller) CanSlash(opts *runOptions, _offender common.Address, _rule uint8, _block *big.Int) (bool, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "canSlash", _offender, _rule, _block)

	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(uint256 innocenceProofSubmissionWindow, uint256 baseSlashingRateLow, uint256 baseSlashingRateMid, uint256 collusionFactor, uint256 historyFactor, uint256 jailFactor, uint256 slashingRatePrecision)
func (_Accountability *AccountabilityCaller) Config(opts *runOptions) (struct {
	InnocenceProofSubmissionWindow *big.Int
	BaseSlashingRateLow            *big.Int
	BaseSlashingRateMid            *big.Int
	CollusionFactor                *big.Int
	HistoryFactor                  *big.Int
	JailFactor                     *big.Int
	SlashingRatePrecision          *big.Int
}, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "config")

	outstruct := new(struct {
		InnocenceProofSubmissionWindow *big.Int
		BaseSlashingRateLow            *big.Int
		BaseSlashingRateMid            *big.Int
		CollusionFactor                *big.Int
		HistoryFactor                  *big.Int
		JailFactor                     *big.Int
		SlashingRatePrecision          *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.InnocenceProofSubmissionWindow = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BaseSlashingRateLow = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BaseSlashingRateMid = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.CollusionFactor = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.HistoryFactor = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.JailFactor = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.SlashingRatePrecision = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// EpochPeriod is a free data retrieval call binding the contract method 0xb5b7a184.
//
// Solidity: function epochPeriod() view returns(uint256)
func (_Accountability *AccountabilityCaller) EpochPeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "epochPeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Events is a free data retrieval call binding the contract method 0x0b791430.
//
// Solidity: function events(uint256 ) view returns(uint8 chunks, uint8 chunkId, uint8 eventType, uint8 rule, address reporter, address offender, bytes rawProof, uint256 id, uint256 block, uint256 epoch, uint256 reportingBlock, uint256 messageHash)
func (_Accountability *AccountabilityCaller) Events(opts *runOptions, arg0 *big.Int) (struct {
	Chunks         uint8
	ChunkId        uint8
	EventType      uint8
	Rule           uint8
	Reporter       common.Address
	Offender       common.Address
	RawProof       []byte
	Id             *big.Int
	Block          *big.Int
	Epoch          *big.Int
	ReportingBlock *big.Int
	MessageHash    *big.Int
}, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "events", arg0)

	outstruct := new(struct {
		Chunks         uint8
		ChunkId        uint8
		EventType      uint8
		Rule           uint8
		Reporter       common.Address
		Offender       common.Address
		RawProof       []byte
		Id             *big.Int
		Block          *big.Int
		Epoch          *big.Int
		ReportingBlock *big.Int
		MessageHash    *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.Chunks = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.ChunkId = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.EventType = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.Rule = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.Reporter = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Offender = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.RawProof = *abi.ConvertType(out[6], new([]byte)).(*[]byte)
	outstruct.Id = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Block = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.Epoch = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)
	outstruct.ReportingBlock = *abi.ConvertType(out[10], new(*big.Int)).(**big.Int)
	outstruct.MessageHash = *abi.ConvertType(out[11], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// GetValidatorAccusation is a free data retrieval call binding the contract method 0x9cb22b06.
//
// Solidity: function getValidatorAccusation(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256))
func (_Accountability *AccountabilityCaller) GetValidatorAccusation(opts *runOptions, _val common.Address) (AccountabilityEvent, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "getValidatorAccusation", _val)

	if err != nil {
		return *new(AccountabilityEvent), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(AccountabilityEvent)).(*AccountabilityEvent)
	return out0, consumed, err

}

// GetValidatorFaults is a free data retrieval call binding the contract method 0xbebaa8fc.
//
// Solidity: function getValidatorFaults(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256)[])
func (_Accountability *AccountabilityCaller) GetValidatorFaults(opts *runOptions, _val common.Address) ([]AccountabilityEvent, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "getValidatorFaults", _val)

	if err != nil {
		return *new([]AccountabilityEvent), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]AccountabilityEvent)).(*[]AccountabilityEvent)
	return out0, consumed, err

}

// SlashingHistory is a free data retrieval call binding the contract method 0xe7bb0b52.
//
// Solidity: function slashingHistory(address , uint256 ) view returns(uint256)
func (_Accountability *AccountabilityCaller) SlashingHistory(opts *runOptions, arg0 common.Address, arg1 *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "slashingHistory", arg0, arg1)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_Accountability *AccountabilityTransactor) DistributeRewards(opts *runOptions, _validator common.Address) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "distributeRewards", _validator)
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_Accountability *AccountabilityTransactor) Finalize(opts *runOptions, _epochEnd bool) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "finalize", _epochEnd)
	return consumed, err
}

// HandleEvent is a paid mutator transaction binding the contract method 0xc50d21f0.
//
// Solidity: function handleEvent((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256) _event) returns()
func (_Accountability *AccountabilityTransactor) HandleEvent(opts *runOptions, _event AccountabilityEvent) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "handleEvent", _event)
	return consumed, err
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_Accountability *AccountabilityTransactor) SetEpochPeriod(opts *runOptions, _newPeriod *big.Int) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "setEpochPeriod", _newPeriod)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// AccountabilityInnocenceProvenIterator is returned from FilterInnocenceProven and is used to iterate over the raw logs and unpacked data for InnocenceProven events raised by the Accountability contract.
		type AccountabilityInnocenceProvenIterator struct {
			Event *AccountabilityInnocenceProven // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AccountabilityInnocenceProvenIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AccountabilityInnocenceProven)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AccountabilityInnocenceProven)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AccountabilityInnocenceProvenIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AccountabilityInnocenceProvenIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AccountabilityInnocenceProven represents a InnocenceProven event raised by the Accountability contract.
		type AccountabilityInnocenceProven struct {
			Offender common.Address;
			Id *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterInnocenceProven is a free log retrieval operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
		//
		// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
 		func (_Accountability *AccountabilityFilterer) FilterInnocenceProven(opts *bind.FilterOpts, _offender []common.Address) (*AccountabilityInnocenceProvenIterator, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}


			logs, sub, err := _Accountability.contract.FilterLogs(opts, "InnocenceProven", _offenderRule)
			if err != nil {
				return nil, err
			}
			return &AccountabilityInnocenceProvenIterator{contract: _Accountability.contract, event: "InnocenceProven", logs: logs, sub: sub}, nil
 		}

		// WatchInnocenceProven is a free log subscription operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
		//
		// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
		func (_Accountability *AccountabilityFilterer) WatchInnocenceProven(opts *bind.WatchOpts, sink chan<- *AccountabilityInnocenceProven, _offender []common.Address) (event.Subscription, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}


			logs, sub, err := _Accountability.contract.WatchLogs(opts, "InnocenceProven", _offenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AccountabilityInnocenceProven)
						if err := _Accountability.contract.UnpackLog(event, "InnocenceProven", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseInnocenceProven is a log parse operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
		//
		// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
		func (_Accountability *AccountabilityFilterer) ParseInnocenceProven(log types.Log) (*AccountabilityInnocenceProven, error) {
			event := new(AccountabilityInnocenceProven)
			if err := _Accountability.contract.UnpackLog(event, "InnocenceProven", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AccountabilityNewAccusationIterator is returned from FilterNewAccusation and is used to iterate over the raw logs and unpacked data for NewAccusation events raised by the Accountability contract.
		type AccountabilityNewAccusationIterator struct {
			Event *AccountabilityNewAccusation // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AccountabilityNewAccusationIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AccountabilityNewAccusation)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AccountabilityNewAccusation)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AccountabilityNewAccusationIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AccountabilityNewAccusationIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AccountabilityNewAccusation represents a NewAccusation event raised by the Accountability contract.
		type AccountabilityNewAccusation struct {
			Offender common.Address;
			Severity *big.Int;
			Id *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewAccusation is a free log retrieval operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
		//
		// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
 		func (_Accountability *AccountabilityFilterer) FilterNewAccusation(opts *bind.FilterOpts, _offender []common.Address) (*AccountabilityNewAccusationIterator, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}



			logs, sub, err := _Accountability.contract.FilterLogs(opts, "NewAccusation", _offenderRule)
			if err != nil {
				return nil, err
			}
			return &AccountabilityNewAccusationIterator{contract: _Accountability.contract, event: "NewAccusation", logs: logs, sub: sub}, nil
 		}

		// WatchNewAccusation is a free log subscription operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
		//
		// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
		func (_Accountability *AccountabilityFilterer) WatchNewAccusation(opts *bind.WatchOpts, sink chan<- *AccountabilityNewAccusation, _offender []common.Address) (event.Subscription, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}



			logs, sub, err := _Accountability.contract.WatchLogs(opts, "NewAccusation", _offenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AccountabilityNewAccusation)
						if err := _Accountability.contract.UnpackLog(event, "NewAccusation", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewAccusation is a log parse operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
		//
		// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
		func (_Accountability *AccountabilityFilterer) ParseNewAccusation(log types.Log) (*AccountabilityNewAccusation, error) {
			event := new(AccountabilityNewAccusation)
			if err := _Accountability.contract.UnpackLog(event, "NewAccusation", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AccountabilityNewFaultProofIterator is returned from FilterNewFaultProof and is used to iterate over the raw logs and unpacked data for NewFaultProof events raised by the Accountability contract.
		type AccountabilityNewFaultProofIterator struct {
			Event *AccountabilityNewFaultProof // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AccountabilityNewFaultProofIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AccountabilityNewFaultProof)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AccountabilityNewFaultProof)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AccountabilityNewFaultProofIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AccountabilityNewFaultProofIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AccountabilityNewFaultProof represents a NewFaultProof event raised by the Accountability contract.
		type AccountabilityNewFaultProof struct {
			Offender common.Address;
			Severity *big.Int;
			Id *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewFaultProof is a free log retrieval operation binding the contract event 0x6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f.
		//
		// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id)
 		func (_Accountability *AccountabilityFilterer) FilterNewFaultProof(opts *bind.FilterOpts, _offender []common.Address) (*AccountabilityNewFaultProofIterator, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}



			logs, sub, err := _Accountability.contract.FilterLogs(opts, "NewFaultProof", _offenderRule)
			if err != nil {
				return nil, err
			}
			return &AccountabilityNewFaultProofIterator{contract: _Accountability.contract, event: "NewFaultProof", logs: logs, sub: sub}, nil
 		}

		// WatchNewFaultProof is a free log subscription operation binding the contract event 0x6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f.
		//
		// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id)
		func (_Accountability *AccountabilityFilterer) WatchNewFaultProof(opts *bind.WatchOpts, sink chan<- *AccountabilityNewFaultProof, _offender []common.Address) (event.Subscription, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}



			logs, sub, err := _Accountability.contract.WatchLogs(opts, "NewFaultProof", _offenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AccountabilityNewFaultProof)
						if err := _Accountability.contract.UnpackLog(event, "NewFaultProof", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewFaultProof is a log parse operation binding the contract event 0x6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f.
		//
		// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id)
		func (_Accountability *AccountabilityFilterer) ParseNewFaultProof(log types.Log) (*AccountabilityNewFaultProof, error) {
			event := new(AccountabilityNewFaultProof)
			if err := _Accountability.contract.UnpackLog(event, "NewFaultProof", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AccountabilitySlashingEventIterator is returned from FilterSlashingEvent and is used to iterate over the raw logs and unpacked data for SlashingEvent events raised by the Accountability contract.
		type AccountabilitySlashingEventIterator struct {
			Event *AccountabilitySlashingEvent // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AccountabilitySlashingEventIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AccountabilitySlashingEvent)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AccountabilitySlashingEvent)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AccountabilitySlashingEventIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AccountabilitySlashingEventIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AccountabilitySlashingEvent represents a SlashingEvent event raised by the Accountability contract.
		type AccountabilitySlashingEvent struct {
			Validator common.Address;
			Amount *big.Int;
			ReleaseBlock *big.Int;
			IsJailbound bool;
			EventId *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterSlashingEvent is a free log retrieval operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
 		func (_Accountability *AccountabilityFilterer) FilterSlashingEvent(opts *bind.FilterOpts) (*AccountabilitySlashingEventIterator, error) {







			logs, sub, err := _Accountability.contract.FilterLogs(opts, "SlashingEvent")
			if err != nil {
				return nil, err
			}
			return &AccountabilitySlashingEventIterator{contract: _Accountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
 		}

		// WatchSlashingEvent is a free log subscription operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
		func (_Accountability *AccountabilityFilterer) WatchSlashingEvent(opts *bind.WatchOpts, sink chan<- *AccountabilitySlashingEvent) (event.Subscription, error) {







			logs, sub, err := _Accountability.contract.WatchLogs(opts, "SlashingEvent")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AccountabilitySlashingEvent)
						if err := _Accountability.contract.UnpackLog(event, "SlashingEvent", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseSlashingEvent is a log parse operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
		func (_Accountability *AccountabilityFilterer) ParseSlashingEvent(log types.Log) (*AccountabilitySlashingEvent, error) {
			event := new(AccountabilitySlashingEvent)
			if err := _Accountability.contract.UnpackLog(event, "SlashingEvent", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// AutonityMetaData contains all meta data concerning the Autonity contract.
var AutonityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator[]\",\"name\":\"_validators\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialInflationReserve\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"},{\"internalType\":\"contractIInflationController\",\"name\":\"inflationControllerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"ActivatedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"BondingRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewUnbondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"atnTotalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialInflationReserve\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"},{\"internalType\":\"contractIInflationController\",\"name\":\"inflationControllerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeInitialization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getEpochFromBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inflationReserve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_signatures\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAccountability\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAccountabilityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIACU\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAcuContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setOracleContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIStabilization\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setStabilizationContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISupplyControl\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSupplyControlContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractUpgradeManager\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setUpgradeManagerContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"_val\",\"type\":\"tuple\"}],\"name\":\"updateValidatorAndTransferSlashedFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"b46e5520": "activateValidator(address)",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"e98712b2": "atnTotalRedistributed()",
		"70a08231": "balanceOf(address)",
		"a515366a": "bond(address,uint256)",
		"9dc29fac": "burn(address,uint256)",
		"852c4849": "changeCommissionRate(address,uint256)",
		"872cf059": "completeContractUpgrade()",
		"ae1f5fa0": "computeCommittee()",
		"79502c55": "config()",
		"313ce567": "decimals()",
		"d5f39488": "deployer()",
		"c9d97af4": "epochID()",
		"1604e416": "epochReward()",
		"9c98e471": "epochTotalBondedStake()",
		"4bb278f3": "finalize()",
		"d861b0e8": "finalizeInitialization()",
		"43645969": "getBlockPeriod()",
		"ab8f6ffe": "getCommittee()",
		"a8b2216e": "getCommitteeEnodes()",
		"96b477cb": "getEpochFromBlock(uint256)",
		"dfb1a4d2": "getEpochPeriod()",
		"731b3a03": "getLastEpochBlock()",
		"819b6463": "getMaxCommitteeSize()",
		"11220633": "getMinimumBaseFee()",
		"b66b3e79": "getNewContract()",
		"e7f43c68": "getOperator()",
		"833b1fce": "getOracle()",
		"5f7d3949": "getProposer(uint256,uint256)",
		"f7866ee3": "getTreasuryAccount()",
		"29070c6d": "getTreasuryFee()",
		"6fd2c80b": "getUnbondingPeriod()",
		"1904bb2e": "getValidator(address)",
		"b7ab4db5": "getValidators()",
		"0d8e6e2c": "getVersion()",
		"9ac31074": "inflationReserve()",
		"c2362dd5": "lastEpochBlock()",
		"40c10f19": "mint(address,uint256)",
		"06fdde03": "name()",
		"0ae65e7a": "pauseValidator(address)",
		"84467fdb": "registerValidator(string,address,bytes,bytes)",
		"cf9c5719": "resetContractUpgrade()",
		"1250a28d": "setAccountabilityContract(address)",
		"d372c07e": "setAcuContract(address)",
		"8bac7dad": "setCommitteeSize(uint256)",
		"6b5f444c": "setEpochPeriod(uint256)",
		"cb696f54": "setMinimumBaseFee(uint256)",
		"520fdbbc": "setOperatorAccount(address)",
		"496ccd9b": "setOracleContract(address)",
		"cfd19fb9": "setStabilizationContract(address)",
		"b3ecbadd": "setSupplyControlContract(address)",
		"d886f8a2": "setTreasuryAccount(address)",
		"77e741c7": "setTreasuryFee(uint256)",
		"114eaf55": "setUnbondingPeriod(uint256)",
		"ceaad455": "setUpgradeManagerContract(address)",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"a5d059ca": "unbond(address,uint256)",
		"35be16e0": "updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8))",
		"b2ea9adb": "upgradeContract(bytes,string)",
	},
	Bin: "0x60806040525f600b555f600c5534801562000018575f80fd5b506040516200a68e3803806200a68e8339810160408190526200003b9162000e68565b601e545f036200006357602d80546001600160a01b031916331790556200006382826200006b565b50506200130d565b80518051600d55602080820151600e55604080830151600f55606080840151601055608080850151601181905560a095860151601280546001600160a01b03199081166001600160a01b0393841617909155868901518051601380548416918516919091179055808801516014805484169185169190911790558087015160158054841691851691909117905580860151601680548416918516919091179055938401516017805483169184169190911790559683015160188054891691831691909117905560c090920151601980548816918416919091179055838701518051601a805490981693169290921790955592830151601b5590820151601c5590810151601d55820151601e55602c555f5b825181101562000430575f8382815181106200019c576200019c6200109c565b602002602001015160a0015190505f848381518110620001c057620001c06200109c565b60200260200101516101a00181815250505f848381518110620001e757620001e76200109c565b602002602001015161018001906001600160a01b031690816001600160a01b0316815250505f8483815181106200022257620002226200109c565b602002602001015160a00181815250505f8483815181106200024857620002486200109c565b60209081029190910101516101c00152600f5484518590849081106200027257620002726200109c565b602002602001015160800181815250505f8483815181106200029857620002986200109c565b602002602001015161026001906003811115620002b957620002b9620010b0565b90816003811115620002cf57620002cf620010b0565b815250505f848381518110620002e957620002e96200109c565b6020026020010151610160018181525050620003278483815181106200031357620003136200109c565b60200260200101516200043560201b60201c565b620003548483815181106200034057620003406200109c565b60200260200101516200056e60201b60201c565b8060295f8685815181106200036d576200036d6200109c565b60200260200101515f01516001600160a01b03166001600160a01b031681526020019081526020015f205f828254620003a79190620010d8565b9250508190555080602b5f828254620003c19190620010d8565b925050819055506200041a848381518110620003e157620003e16200109c565b602002602001015160200151828685815181106200040357620004036200109c565b60200260200101515f01516200079960201b60201c565b50806200042781620010f4565b9150506200017c565b505050565b5f6200044b82606001516200097f60201b60201c565b6001600160a01b039091166020840152905080156200049f5760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b60448201526064015b60405180910390fd5b6020808301516001600160a01b039081165f908152602a9092526040909120600101541615620005125760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000496565b612710826080015111156200056a5760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000496565b5050565b6101808101516001600160a01b0316620005ec57601f545f906200059290620009ca565b90508160200151825f0151836080015183604051620005b19062000ae6565b620005c094939291906200110f565b604051809103905ff080158015620005da573d5f803e3d5ffd5b506001600160a01b0316610180830152505b60208181018051601f8054600180820183555f9283527fa03837a25210ee280c2113ff4b77ca23440b19d4866cca721c801278fd08d80790910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602a909552604091829020865181548716908516178155935190840180548616918416919091179055840151600283018054909416911617909155606082015182919060038201906200069d9082620011ef565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620007679082620011ef565b5061026082015160138201805460ff191660018360038111156200078f576200078f620010b0565b0217905550505050565b5f8211620007f65760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000496565b6001600160a01b0381165f908152602960205260409020548211156200085f5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000496565b6001600160a01b0381165f908152602960205260408120805484929062000888908490620012b7565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580545f908152600394859052978820875181549088166001600160a01b03199182161782559551600182018054919098169616959095179095559051600284015551910155805491926200091383620010f4565b90915550506001600160a01b038481165f818152602a6020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d910160405180910390a35050505050565b5f806200098b62000af4565b5f8060ff9050604083875160208901845afa620009a6575f80fd5b505080516020909101516c0100000000000000000000000090910494909350915050565b6060815f03620009f15750506040805180820190915260018152600360fc1b602082015290565b815f5b811562000a20578062000a0781620010f4565b915062000a189050600a83620012e1565b9150620009f4565b5f816001600160401b0381111562000a3c5762000a3c62000b12565b6040519080825280601f01601f19166020018201604052801562000a67576020820181803683370190505b5090505b841562000ade5762000a7f600183620012b7565b915062000a8e600a86620012f7565b62000a9b906030620010d8565b60f81b81838151811062000ab35762000ab36200109c565b60200101906001600160f81b03191690815f1a90535062000ad6600a86620012e1565b945062000a6b565b949350505050565b6117d68062008eb883390190565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b5f52604160045260245ffd5b604051608081016001600160401b038111828210171562000b4b5762000b4b62000b12565b60405290565b60405160c081016001600160401b038111828210171562000b4b5762000b4b62000b12565b60405161028081016001600160401b038111828210171562000b4b5762000b4b62000b12565b604051601f8201601f191681016001600160401b038111828210171562000bc75762000bc762000b12565b604052919050565b6001600160a01b038116811462000be4575f80fd5b50565b805162000bf48162000bcf565b919050565b5f5b8381101562000c1557818101518382015260200162000bfb565b50505f910152565b5f82601f83011262000c2d575f80fd5b81516001600160401b0381111562000c495762000c4962000b12565b62000c5e601f8201601f191660200162000b9c565b81815284602083860101111562000c73575f80fd5b62000ade82602083016020870162000bf9565b80516004811062000bf4575f80fd5b5f60e0828403121562000ca6575f80fd5b60405160e081016001600160401b038111828210171562000ccb5762000ccb62000b12565b8060405250809150825162000ce08162000bcf565b8152602083015162000cf28162000bcf565b6020820152604083015162000d078162000bcf565b6040820152606083015162000d1c8162000bcf565b606082015262000d2f6080840162000be7565b608082015262000d4260a0840162000be7565b60a082015262000d5560c0840162000be7565b60c08201525092915050565b5f6080828403121562000d72575f80fd5b62000d7c62000b26565b9050815162000d8b8162000bcf565b8082525060208201516020820152604082015160408201526060820151606082015292915050565b5f81830361024081121562000dc6575f80fd5b62000dd062000b26565b915060c081121562000de0575f80fd5b5062000deb62000b51565b825181526020830151602082015260408301516040820152606083015160608201526080830151608082015260a083015162000e278162000bcf565b60a0820152815262000e3d8360c0840162000c95565b602082015262000e52836101a0840162000d61565b6040820152610220820151606082015292915050565b5f8061026080848603121562000e7c575f80fd5b83516001600160401b038082111562000e93575f80fd5b818601915086601f83011262000ea7575f80fd5b815160208282111562000ebe5762000ebe62000b12565b8160051b62000ecf82820162000b9c565b928352848101820192828101908b85111562000ee9575f80fd5b83870192505b848310156200107a5782518681111562000f07575f80fd5b8701610280818e03601f1901121562000f1e575f80fd5b62000f2862000b76565b62000f3586830162000be7565b815262000f456040830162000be7565b8682015262000f576060830162000be7565b604082015260808201518881111562000f6e575f80fd5b62000f7e8f888386010162000c1d565b60608301525060a0820151608082015260c082015160a082015260e082015160c082015261010082015160e082015261012082015161010082015261014082015161012082015261016082015161014082015261018082015161016082015262000fec6101a0830162000be7565b6101808201526101c08201516101a08201526101e08201516101c08201526102008201516101e082015261022082015161020082015261024082015161022082015289820151888111156200103f575f80fd5b6200104f8f888386010162000c1d565b6102408301525062001065610280830162000c86565b818b0152835250918301919083019062000eef565b8099505050506200108e89828a0162000db3565b955050505050509250929050565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b80820180821115620010ee57620010ee620010c4565b92915050565b5f60018201620011085762001108620010c4565b5060010190565b5f60018060a01b0380871683528086166020840152508360408301526080606083015282518060808401526200114d8160a085016020870162000bf9565b601f01601f19169190910160a00195945050505050565b600181811c908216806200117957607f821691505b6020821081036200119857634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111562000430575f81815260208120601f850160051c81016020861015620011c65750805b601f850160051c820191505b81811015620011e757828155600101620011d2565b505050505050565b81516001600160401b038111156200120b576200120b62000b12565b62001223816200121c845462001164565b846200119e565b602080601f83116001811462001259575f8415620012415750858301515b5f19600386901b1c1916600185901b178555620011e7565b5f85815260208120601f198616915b82811015620012895788860151825594840194600190910190840162001268565b5085821015620012a757878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b81810381811115620010ee57620010ee620010c4565b634e487b7160e01b5f52601260045260245ffd5b5f82620012f257620012f2620012cd565b500490565b5f82620013085762001308620012cd565b500690565b617b9d806200131b5f395ff3fe608060405260043610620003e6575f3560e01c80638bac7dad1162000202578063b7ab4db51162000116578063d5f3948811620000a6578063dfb1a4d21162000074578063dfb1a4d21462000da4578063e7f43c681462000dba578063e98712b21462000dd9578063f7866ee31462000df057005b8063d5f394881462000d00578063d861b0e81462000d21578063d886f8a21462000d38578063dd62ed3e1462000d5c57005b8063ceaad45511620000e4578063ceaad4551462000c7d578063cf9c57191462000ca1578063cfd19fb91462000cb8578063d372c07e1462000cdc57005b8063b7ab4db51462000c14578063c2362dd51462000c2b578063c9d97af41462000c42578063cb696f541462000c5957005b8063a8b2216e1162000192578063b2ea9adb1162000160578063b2ea9adb1462000b81578063b3ecbadd1462000ba5578063b46e55201462000bc9578063b66b3e791462000bed57005b8063a8b2216e1462000aeb578063a9059cbb1462000b11578063ab8f6ffe1462000b35578063ae1f5fa01462000b5b57005b80639c98e47111620001d05780639c98e4711462000a685780639dc29fac1462000a7f578063a515366a1462000aa3578063a5d059ca1462000ac757005b80638bac7dad14620009d157806395d89b4114620009f557806396b477cb1462000a225780639ac310741462000a5157005b80634364596911620002fa578063731b3a03116200028a578063833b1fce1162000258578063833b1fce146200095357806384467fdb1462000972578063852c48491462000996578063872cf05914620009ba57005b8063731b3a03146200075457806377e741c7146200076a57806379502c55146200078e578063819b6463146200093d57005b80635f7d394911620002c85780635f7d394914620006a55780636b5f444c14620006e25780636fd2c80b146200070657806370a08231146200071c57005b8063436459691462000620578063496ccd9b14620006365780634bb278f3146200065a578063520fdbbc146200068157005b806318160ddd11620003765780632f2c3f2e11620003445780632f2c3f2e14620005a4578063313ce56714620005bb57806335be16e014620005d857806340c10f1914620005fc57005b806318160ddd14620005215780631904bb2e146200053757806323b872dd146200056a57806329070c6d146200058e57005b80631122063311620003b45780631122063314620004ac578063114eaf5514620004c25780631250a28d14620004e65780631604e416146200050a57005b806306fdde0314620003f0578063095ea7b314620004335780630ae65e7a14620004685780630d8e6e2c146200048c57005b36620003ee57005b005b348015620003fc575f80fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b6040516200042a91906200545f565b60405180910390f35b3480156200043f575f80fd5b5062000457620004513660046200548f565b62000e0f565b60405190151581526020016200042a565b34801562000474575f80fd5b50620003ee62000486366004620054bc565b62000e27565b34801562000498575f80fd5b50601e545b6040519081526020016200042a565b348015620004b8575f80fd5b50600e546200049d565b348015620004ce575f80fd5b50620003ee620004e0366004620054da565b62000eb8565b348015620004f2575f80fd5b50620003ee62000504366004620054bc565b62000eea565b34801562000516575f80fd5b506200049d60265481565b3480156200052d575f80fd5b50602b546200049d565b34801562000543575f80fd5b506200055b62000555366004620054bc565b62000f39565b6040516200042a919062005527565b34801562000576575f80fd5b5062000457620005883660046200569f565b620011e1565b3480156200059a575f80fd5b50600d546200049d565b348015620005b0575f80fd5b506200049d61271081565b348015620005c7575f80fd5b50604051601281526020016200042a565b348015620005e4575f80fd5b50620003ee620005f6366004620056e2565b62001287565b34801562000608575f80fd5b50620003ee6200061a3660046200548f565b6200145a565b3480156200062c575f80fd5b50601c546200049d565b34801562000642575f80fd5b50620003ee62000654366004620054bc565b62001497565b34801562000666575f80fd5b506200067162001597565b6040516200042a9291906200579d565b3480156200068d575f80fd5b50620003ee6200069f366004620054bc565b62001a0b565b348015620006b1575f80fd5b50620006c9620006c3366004620057b9565b62001bf2565b6040516001600160a01b0390911681526020016200042a565b348015620006ee575f80fd5b50620003ee62000700366004620054da565b62001dfe565b34801562000712575f80fd5b506010546200049d565b34801562000728575f80fd5b506200049d6200073a366004620054bc565b6001600160a01b03165f9081526029602052604090205490565b34801562000760575f80fd5b506022546200049d565b34801562000776575f80fd5b50620003ee62000788366004620054da565b62001f7f565b3480156200079a575f80fd5b506040805160c08082018352600d548252600e54602080840191909152600f54838501526010546060808501919091526011546080808601919091526012546001600160a01b0390811660a080880191909152875160e081018952601354831681526014548316818701526015548316818a015260165483168186015260175483168185015260185483169181019190915260195482169581019590955286519182018752601a54168152601b5492810192909252601c5494820194909452601d5493810193909352601e546200086f939084565b6040805185518152602080870151818301528683015182840152606080880151818401526080808901518185015260a0988901516001600160a01b039081168a8601528851811660c080870191909152848a0151821660e0870152898701518216610100870152838a0151821661012087015291890151811661014086015298880151891661016085015290960151871661018083015284519096166101a0820152948301516101c08601528201516101e0850152910151610200830152610220820152610240016200042a565b34801562000949575f80fd5b50601d546200049d565b3480156200095f575f80fd5b506014546001600160a01b0316620006c9565b3480156200097e575f80fd5b50620003ee620009903660046200587f565b62001fb1565b348015620009a2575f80fd5b50620003ee620009b43660046200548f565b620020d0565b348015620009c6575f80fd5b50620003ee6200225e565b348015620009dd575f80fd5b50620003ee620009ef366004620054da565b6200229a565b34801562000a01575f80fd5b50604080518082019091526003815262272a2760e91b60208201526200041b565b34801562000a2e575f80fd5b506200049d62000a40366004620054da565b5f9081526021602052604090205490565b34801562000a5d575f80fd5b506200049d602c5481565b34801562000a74575f80fd5b506200049d60235481565b34801562000a8b575f80fd5b50620003ee62000a9d3660046200548f565b6200231d565b34801562000aaf575f80fd5b50620003ee62000ac13660046200548f565b62002437565b34801562000ad3575f80fd5b50620003ee62000ae53660046200548f565b62002504565b34801562000af7575f80fd5b5062000b0262002599565b6040516200042a919062005922565b34801562000b1d575f80fd5b506200045762000b2f3660046200548f565b62002677565b34801562000b41575f80fd5b5062000b4c620026ce565b6040516200042a919062005986565b34801562000b67575f80fd5b5062000b72620027dc565b6040516200042a91906200599a565b34801562000b8d575f80fd5b50620003ee62000b9f366004620059e8565b62002a1c565b34801562000bb1575f80fd5b50620003ee62000bc3366004620054bc565b62002a62565b34801562000bd5575f80fd5b50620003ee62000be7366004620054bc565b62002ab1565b34801562000bf9575f80fd5b5062000c0462002d31565b6040516200042a92919062005a4f565b34801562000c20575f80fd5b5062000b7262002e63565b34801562000c37575f80fd5b506200049d60225481565b34801562000c4e575f80fd5b506200049d60205481565b34801562000c65575f80fd5b50620003ee62000c77366004620054da565b62002ec5565b34801562000c89575f80fd5b50620003ee62000c9b366004620054bc565b62002f28565b34801562000cad575f80fd5b50620003ee62002f77565b34801562000cc4575f80fd5b50620003ee62000cd6366004620054bc565b62002fc9565b34801562000ce8575f80fd5b50620003ee62000cfa366004620054bc565b62003018565b34801562000d0c575f80fd5b50602d54620006c9906001600160a01b031681565b34801562000d2d575f80fd5b50620003ee62003067565b34801562000d44575f80fd5b50620003ee62000d56366004620054bc565b620030a8565b34801562000d68575f80fd5b506200049d62000d7a36600462005a80565b6001600160a01b039182165f90815260286020908152604080832093909416825291909152205490565b34801562000db0575f80fd5b50601b546200049d565b34801562000dc6575f80fd5b50601a546001600160a01b0316620006c9565b34801562000de5575f80fd5b506200049d60255481565b34801562000dfc575f80fd5b506012546001600160a01b0316620006c9565b5f62000e1d338484620030f7565b5060015b92915050565b6001600160a01b038082165f818152602a60205260409020600101549091161462000e6f5760405162461bcd60e51b815260040162000e669062005abc565b60405180910390fd5b6001600160a01b038181165f908152602a602052604090205416331462000eaa5760405162461bcd60e51b815260040162000e669062005af3565b62000eb5816200321e565b50565b601a546001600160a01b0316331462000ee55760405162461bcd60e51b815260040162000e669062005b3f565b601055565b601a546001600160a01b0316331462000f175760405162461bcd60e51b815260040162000e669062005b3f565b601380546001600160a01b0319166001600160a01b0392909216919091179055565b62000f4362005255565b6001600160a01b038083165f818152602a60205260409020600101549091161462000f825760405162461bcd60e51b815260040162000e669062005b76565b6001600160a01b038083165f908152602a6020908152604091829020825161028081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000fe49062005bad565b80601f0160208091040260200160405190810160405280929190818152602001828054620010129062005bad565b8015620010615780601f10620010375761010080835404028352916020019162001061565b820191905f5260205f20905b8154815290600101906020018083116200104357829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546080820152600882015460a0820152600982015460c0820152600a82015460e0820152600b820154610100820152600c8201546001600160a01b0316610120820152600d820154610140820152600e820154610160820152600f82015461018082015260108201546101a082015260118201546101c08201526012820180546101e0909201916200111d9062005bad565b80601f01602080910402602001604051908101604052809291908181526020018280546200114b9062005bad565b80156200119a5780601f1062001170576101008083540402835291602001916200119a565b820191905f5260205f20905b8154815290600101906020018083116200117c57829003601f168201915b5050509183525050601382015460209091019060ff166003811115620011c457620011c4620054f2565b6003811115620011d857620011d8620054f2565b90525092915050565b5f620011ef848484620032f4565b6001600160a01b0384165f9081526028602090815260408083203384529091528120546200121f90849062005bfb565b90506200122e853383620030f7565b836001600160a01b0316856001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040516200127491815260200190565b60405180910390a3506001949350505050565b6013546001600160a01b03163314620012ef5760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f742074686520736c617368696e6720636f6e74726044820152621858dd60ea1b606482015260840162000e66565b5f610120820135602a826200130b6040860160208701620054bc565b6001600160a01b03166001600160a01b031681526020019081526020015f20600901546200133a919062005bfb565b60c0830135602a5f620013546040870160208801620054bc565b6001600160a01b03166001600160a01b031681526020019081526020015f206006015462001383919062005bfb565b60a0840135602a5f6200139d6040880160208901620054bc565b6001600160a01b03166001600160a01b031681526020019081526020015f2060050154620013cc919062005bfb565b620013d8919062005c11565b620013e4919062005c11565b6012546001600160a01b03165f908152602960205260408120805492935083929091906200141490849062005c11565b90915550829050602a5f620014306040840160208501620054bc565b6001600160a01b0316815260208101919091526040015f2062001454828262005e0c565b50505050565b601a546001600160a01b03163314620014875760405162461bcd60e51b815260040162000e669062005b3f565b620014938282620033b7565b5050565b601a546001600160a01b03163314620014c45760405162461bcd60e51b815260040162000e669062005b3f565b601480546001600160a01b0319166001600160a01b03838116918217909255601554604051637adbf97360e01b8152600481019290925290911690637adbf973906024015f604051808303815f87803b15801562001520575f80fd5b505af115801562001533573d5f803e3d5ffd5b5050601754604051637adbf97360e01b81526001600160a01b0385811660048301529091169250637adbf97391506024015b5f604051808303815f87803b1580156200157d575f80fd5b505af115801562001590573d5f803e3d5ffd5b5050505050565b602d545f906060906001600160a01b03163314620015c95760405162461bcd60e51b815260040162000e669062005f62565b60208054435f81815260219093526040832091909155601b54602254620015f1919062005c11565b6013546040516306c9789b60e41b8152929091146004830181905292506001600160a01b031690636c9789b0906024015f604051808303815f87803b15801562001639575f80fd5b505af11580156200164c573d5f803e3d5ffd5b505050508015620017fe57601954602b54602c546022546040516392eff3cd60e01b81526004810193909352602483019190915260448201524360648201525f916001600160a01b0316906392eff3cd90608401602060405180830381865afa158015620016bc573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620016e2919062005fa5565b905080602c541015620016f45750602c545b620017003082620033b7565b80602c5f82825462001713919062005bfb565b9091555062001725905047826200343b565b6200172f620038de565b62001739620039cc565b5f62001744620027dc565b60145460405163422811f960e11b81529192506001600160a01b03169063845023f290620017779084906004016200599a565b5f604051808303815f87803b1580156200178f575f80fd5b505af1158015620017a2573d5f803e3d5ffd5b5050505043602281905550600160205f828254620017c1919062005c11565b9091555050602080546040519081527febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335910160405180910390a150505b60145460408051634bb278f360e01b815290515f926001600160a01b031691634bb278f3916004808301926020929190829003018187875af115801562001847573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906200186d919062005fbd565b90508015620018e8576015546040805163a2e6204560e01b815290516001600160a01b039092169163a2e620459160048082019260209290919082900301815f875af1925050508015620018e0575060408051601f3d908101601f19168201909252620018dd9181019062005fbd565b60015b15620018e857505b600254602480546040805160208084028201810190925282815260ff909416939183915f9084015b82821015620019fb575f848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620019659062005bad565b80601f0160208091040260200160405190810160405280929190818152602001828054620019939062005bad565b8015620019e25780601f10620019b857610100808354040283529160200191620019e2565b820191905f5260205f20905b815481529060010190602001808311620019c457829003601f168201915b5050505050815250508152602001906001019062001910565b5050505090509350935050509091565b601a546001600160a01b0316331462001a385760405162461bcd60e51b815260040162000e669062005b3f565b601a80546001600160a01b0319166001600160a01b0383811691821790925560145460405163b3ab15fb60e01b815260048101929092529091169063b3ab15fb906024015f604051808303815f87803b15801562001a94575f80fd5b505af115801562001aa7573d5f803e3d5ffd5b505060155460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024015f604051808303815f87803b15801562001af0575f80fd5b505af115801562001b03573d5f803e3d5ffd5b505060165460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024015f604051808303815f87803b15801562001b4c575f80fd5b505af115801562001b5f573d5f803e3d5ffd5b505060175460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024015f604051808303815f87803b15801562001ba8575f80fd5b505af115801562001bbb573d5f803e3d5ffd5b505060185460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb915060240162001565565b5f80805b60245481101562001c4b576024818154811062001c175762001c1762005fde565b905f5260205f209060030201600101548262001c34919062005c11565b91508062001c428162005ff2565b91505062001bf6565b50805f0362001c9d5760405162461bcd60e51b815260206004820152601c60248201527f54686520636f6d6d6974746565206973206e6f74207374616b696e6700000000604482015260640162000e66565b5f8362001cac6063876200600d565b62001cb8919062005c11565b90505f8160405160200162001ccf91815260200190565b60408051601f19818403018152919052805160209091012090505f62001cf684836200603b565b90505f805b60245481101562001da2576024818154811062001d1c5762001d1c62005fde565b905f5260205f209060030201600101548262001d39919062005c11565b915062001d4860018362005bfb565b831162001d8d576024818154811062001d655762001d6562005fde565b5f9182526020909120600390910201546001600160a01b0316965062000e2195505050505050565b8062001d998162005ff2565b91505062001cfb565b5060405162461bcd60e51b815260206004820152602960248201527f5468657265206973206e6f2076616c696461746f72206c65667420696e20746860448201526865206e6574776f726b60b81b606482015260840162000e66565b601a546001600160a01b0316331462001e2b5760405162461bcd60e51b815260040162000e669062005b3f565b601b5481101562001ee2578060225462001e46919062005c11565b431062001ee25760405162461bcd60e51b815260206004820152605760248201527f63757272656e7420636861696e2068656164206578636565642074686520776960448201527f6e646f773a206c617374426c6f636b45706f6368202b205f6e6577506572696f60648201527f642c2074727920616761696e206c6174746572206f6e2e000000000000000000608482015260a40162000e66565b601b819055601354604051631ad7d11360e21b8152600481018390526001600160a01b0390911690636b5f444c906024015f604051808303815f87803b15801562001f2b575f80fd5b505af115801562001f3e573d5f803e3d5ffd5b505050507fd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f818160405162001f7491815260200190565b60405180910390a150565b601a546001600160a01b0316331462001fac5760405162461bcd60e51b815260040162000e669062005b3f565b600d55565b5f604051806102800160405280336001600160a01b031681526020015f6001600160a01b03168152602001856001600160a01b03168152602001868152602001600d5f016002015481526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f81526020014381526020015f81526020015f81526020015f81526020018481526020015f6003811115620020705762002070620054f2565b9052905062002080818362003ae3565b60208101516101808201516040517f8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c92620020c192339289918b9162006051565b60405180910390a15050505050565b6001600160a01b038083165f818152602a6020526040902060010154909116146200210f5760405162461bcd60e51b815260040162000e669062005abc565b6001600160a01b038281165f908152602a60205260409020541633146200214a5760405162461bcd60e51b815260040162000e669062005af3565b6127108111156200219e5760405162461bcd60e51b815260206004820152601f60248201527f7265717569726520636f727265637420636f6d6d697373696f6e207261746500604482015260640162000e66565b604080516060810182526001600160a01b038481168252436020808401918252838501868152600c80545f908152600a909352958220855181546001600160a01b0319169516949094178455915160018085019190915591516002909301929092558354929390929091906200221690849062005c11565b90915550506040518281526001600160a01b038416907f4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf9060200160405180910390a2505050565b601a546001600160a01b031633146200228b5760405162461bcd60e51b815260040162000e669062005b3f565b6002805460ff19166001179055565b601a546001600160a01b03163314620022c75760405162461bcd60e51b815260040162000e669062005b3f565b5f8111620023185760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000e66565b601d55565b601a546001600160a01b031633146200234a5760405162461bcd60e51b815260040162000e669062005b3f565b6001600160a01b0382165f90815260296020526040902054811115620023ac5760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000e66565b6001600160a01b0382165f9081526029602052604081208054839290620023d590849062005bfb565b9250508190555080602b5f828254620023ef919062005bfb565b90915550506040518181526001600160a01b038316907f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3906020015b60405180910390a25050565b6001600160a01b038083165f818152602a602052604090206001015490911614620024765760405162461bcd60e51b815260040162000e669062005b76565b6001600160a01b0382165f908152602a602052604081206013015460ff166003811115620024a857620024a8620054f2565b14620024f75760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000e66565b6200149382823362003f33565b6001600160a01b038083165f818152602a602052604090206001015490911614620025435760405162461bcd60e51b815260040162000e669062005b76565b5f81116200258c5760405162461bcd60e51b81526020600482015260156024820152740756e626f6e64696e6720616d6f756e74206973203605c1b604482015260640162000e66565b620014938282336200411a565b60606027805480602002602001604051908101604052809291908181526020015f905b828210156200266e578382905f5260205f20018054620025dc9062005bad565b80601f01602080910402602001604051908101604052809291908181526020018280546200260a9062005bad565b8015620026595780601f106200262f5761010080835404028352916020019162002659565b820191905f5260205f20905b8154815290600101906020018083116200263b57829003601f168201915b505050505081526020019060010190620025bc565b50505050905090565b5f62002685338484620032f4565b6040518281526001600160a01b0384169033907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a350600192915050565b60606024805480602002602001604051908101604052809291908181526020015f905b828210156200266e575f848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620027469062005bad565b80601f0160208091040260200160405190810160405280929190818152602001828054620027749062005bad565b8015620027c35780601f106200279957610100808354040283529160200191620027c3565b820191905f5260205f20905b815481529060010190602001808311620027a557829003601f168201915b50505050508152505081526020019060010190620026f1565b602d546060906001600160a01b031633146200280c5760405162461bcd60e51b815260040162000e669062005f62565b601f546200285d5760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000e66565b6200286762005313565b601d546080820152601f8152602a60208201526024604082015260236060820152620028938162004442565b620028a060275f62005331565b60245480620028e75760405162461bcd60e51b8152602060048201526012602482015271636f6d6d697474656520697320656d70747960701b604482015260640162000e66565b5f8167ffffffffffffffff811115620029045762002904620057da565b6040519080825280602002602001820160405280156200292e578160200160208202803683370190505b5090505f5b8281101562002a14575f602a5f6024848154811062002956576200295662005fde565b5f9182526020808320600392830201546001600160a01b031684528301939093526040909101812060278054600181018255925292507f98a476f1687bc3d60a2da2adbcba2c46958e61fa2fb4042cd7bc5816a710195b0190620029bd9083018262006099565b50600281015483516001600160a01b0390911690849084908110620029e657620029e662005fde565b6001600160a01b0390921660209283029190910190910152508062002a0b8162005ff2565b91505062002933565b509250505090565b601a546001600160a01b0316331462002a495760405162461bcd60e51b815260040162000e669062005b3f565b62002a555f836200445f565b620014936001826200445f565b601a546001600160a01b0316331462002a8f5760405162461bcd60e51b815260040162000e669062005b3f565b601680546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038082165f818152602a60205260409020600101549091161462002af05760405162461bcd60e51b815260040162000e669062005abc565b6001600160a01b038082165f908152602a602052604090208054909116331462002b755760405162461bcd60e51b815260206004820152602f60248201527f726571756972652063616c6c657220746f2062652076616c696461746f72207460448201526e1c99585cdd5c9e481858d8dbdd5b9d608a1b606482015260840162000e66565b5f601382015460ff16600381111562002b925762002b92620054f2565b0362002be15760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f7220616c7265616479206163746976650000000000000000604482015260640162000e66565b6002601382015460ff16600381111562002bff5762002bff620054f2565b14801562002c105750438160100154115b1562002c5f5760405162461bcd60e51b815260206004820152601760248201527f76616c696461746f72207374696c6c20696e206a61696c000000000000000000604482015260640162000e66565b6003601382015460ff16600381111562002c7d5762002c7d620054f2565b0362002ccc5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206a61696c6564207065726d616e656e746c7900000000604482015260640162000e66565b60138101805460ff191690558054601b546022546001600160a01b038581169316917f60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b59162002d1c919062005c11565b60405190815260200160405180910390a35050565b6060805f600181805462002d459062005bad565b80601f016020809104026020016040519081016040528092919081815260200182805462002d739062005bad565b801562002dc25780601f1062002d985761010080835404028352916020019162002dc2565b820191905f5260205f20905b81548152906001019060200180831162002da457829003601f168201915b5050505050915080805462002dd79062005bad565b80601f016020809104026020016040519081016040528092919081815260200182805462002e059062005bad565b801562002e545780601f1062002e2a5761010080835404028352916020019162002e54565b820191905f5260205f20905b81548152906001019060200180831162002e3657829003601f168201915b50505050509050915091509091565b6060601f80548060200260200160405190810160405280929190818152602001828054801562002ebb57602002820191905f5260205f20905b81546001600160a01b0316815260019091019060200180831162002e9c575b5050505050905090565b601a546001600160a01b0316331462002ef25760405162461bcd60e51b815260040162000e669062005b3f565b600e8190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200162001f74565b601a546001600160a01b0316331462002f555760405162461bcd60e51b815260040162000e669062005b3f565b601880546001600160a01b0319166001600160a01b0392909216919091179055565b601a546001600160a01b0316331462002fa45760405162461bcd60e51b815260040162000e669062005b3f565b62002fb05f806200534e565b62002fbd60015f6200534e565b6002805460ff19169055565b601a546001600160a01b0316331462002ff65760405162461bcd60e51b815260040162000e669062005b3f565b601780546001600160a01b0319166001600160a01b0392909216919091179055565b601a546001600160a01b03163314620030455760405162461bcd60e51b815260040162000e669062005b3f565b601580546001600160a01b0319166001600160a01b0392909216919091179055565b602d546001600160a01b03163314620030945760405162461bcd60e51b815260040162000e669062005f62565b6200309e620038de565b62000eb5620027dc565b601a546001600160a01b03163314620030d55760405162461bcd60e51b815260040162000e669062005b3f565b601280546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b0383166200315b5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000e66565b6001600160a01b038216620031be5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000e66565b6001600160a01b038381165f8181526028602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6001600160a01b0381165f908152602a6020526040812090601382015460ff166003811115620032525762003252620054f2565b14620032a15760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265206163746976650000000000000000604482015260640162000e66565b60138101805460ff191660011790558054601b546022546001600160a01b038581169316917f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c9162002d1c919062005c11565b6001600160a01b0383165f90815260296020526040902054811115620033565760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000e66565b6001600160a01b0383165f90815260296020526040812080548392906200337f90849062005bfb565b90915550506001600160a01b0382165f9081526029602052604081208054839290620033ad90849062005c11565b9091555050505050565b6001600160a01b0382165f9081526029602052604081208054839290620033e090849062005c11565b9250508190555080602b5f828254620033fa919062005c11565b90915550506040518181526001600160a01b038316907f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf906020016200242b565b8115801562003448575080155b1562003452575050565b600d545f90670de0b6b3a7640000906200346e9085906200600d565b6200347a91906200617d565b90508015620034f5576012546040515f916001600160a01b03169083908381818185875af1925050503d805f8114620034cf576040519150601f19603f3d011682016040523d82523d5f602084013e620034d4565b606091505b5090915050801515600103620034f357620034f0828562005bfb565b93505b505b8260255f82825462003508919062005c11565b909155505f90505b60245481101562001454575f602a5f6024848154811062003535576200353562005fde565b5f91825260208083206003909202909101546001600160a01b031683528201929092526040018120602354602480549294509091889190869081106200357f576200357f62005fde565b905f5260205f209060030201600101546200359b91906200600d565b620035a791906200617d565b90505f6023548660248681548110620035c457620035c462005fde565b905f5260205f20906003020160010154620035e091906200600d565b620035ec91906200617d565b90505f821180620035fc57505f81115b15620038c5576002601384015460ff166003811115620036205762003620620054f2565b14806200364857506003601384015460ff166003811115620036465762003646620054f2565b145b15620036ec57601354602480546001600160a01b0390921691631de9d9b6918591889081106200367c576200367c62005fde565b5f91825260209091206003909102015460405160e084901b6001600160e01b03191681526001600160a01b0390911660048201526024015f604051808303818588803b158015620036cb575f80fd5b505af1158015620036de573d5f803e3d5ffd5b5050505050505050620038c9565b5f83600501548385600801546200370491906200600d565b6200371091906200617d565b90508015620037725783546040516001600160a01b03909116906108fc9083905f818181858888f193505050503d805f811462003769576040519150601f19603f3d011682016040523d82523d5f602084013e6200376e565b606091505b5050505b5f84600501548386600801546200378a91906200600d565b6200379691906200617d565b90508015620037b8578454620037b89030906001600160a01b031683620032f4565b5f620037c5828562005bfb565b90505f620037d4848762005bfb565b90505f811180620037e457505f82115b156200387c57600c870154620038069030906001600160a01b031684620032f4565b600c87015460405163a0ce552d60e01b8152600481018490526001600160a01b039091169063a0ce552d908390602401604080518083038185885af115801562003852573d5f803e3d5ffd5b50505050506040513d601f19601f8201168201806040525081019062003879919062006193565b50505b60018701546040518781526001600160a01b03909116907fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe5639060200160405180910390a2505050505b5050505b80620038d58162005ff2565b91505062003510565b6004545b6005548110156200390b576200390581620038fd8162005ff2565b9250620045ae565b620038e2565b50600554600455600854600754036200392057565b6009545b6008548110156200394d5762003947816200393f8162005ff2565b9250620047bf565b62003924565b50600854600955600754805b600854811015620039c6576010545f828152600660205260409020600401544391620039859162005c11565b11620039ab57620039968162004aa4565b620039a360018362005c11565b9150620039b1565b620039c6565b80620039bd8162005ff2565b91505062003959565b50600755565b600c54600b54101562003ae157600b545f908152600a6020526040902060105460018201544391620039fe9162005c11565b111562003a085750565b600281015481546001600160a01b039081165f908152602a6020526040808220600490810185905585548416835291819020600c015490516319fac8fd60e01b81529216926319fac8fd9262003a62920190815260200190565b5f604051808303815f87803b15801562003a7a575f80fd5b505af115801562003a8d573d5f803e3d5ffd5b5050600b80545f908152600a6020526040812080546001600160a01b03191681556001808201839055600290910182905582549094509192509062003ad490849062005c11565b90915550620039cc915050565b565b60e281511462003b2d5760405162461bcd60e51b8152602060048201526014602482015273092dcecc2d8d2c840e0e4dedecc40d8cadccee8d60631b604482015260640162000e66565b6030826102400151511462003b855760405162461bcd60e51b815260206004820152601c60248201527f496e76616c696420636f6e73656e737573206b6579206c656e67746800000000604482015260640162000e66565b62003b908262004be8565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a0000000000006020808301919091528451925191925f9262003bf2920160609190911b6bffffffffffffffffffffffff1916815260140190565b60405160208183030381529060405290505f8262003c11835162004d13565b8360405160200162003c2693929190620061b6565b60408051601f19818403018152828252805160209182012060028085526060850184529094505f93929091830190803683370190505090505f80808062003c7c898262003c76604160026200600d565b62004e30565b90505f62003c9a8a62003c92604160026200600d565b606062004e30565b905060205b825181101562003d685762003cb5838262004f48565b604080515f8152602081018083528d905260ff8316918101919091526060810184905260808101839052929850909650945060019060a0016020604051602081039080840390855afa15801562003d0e573d5f803e3d5ffd5b5050604051601f19015190508762003d286041846200617d565b8151811062003d3b5762003d3b62005fde565b6001600160a01b039092166020928302919091019091015262003d6060418262005c11565b905062003c9f565b508a602001516001600160a01b0316865f8151811062003d8c5762003d8c62005fde565b60200260200101516001600160a01b03161462003dfe5760405162461bcd60e51b815260206004820152602960248201527f496e76616c6964206e6f6465206b6579206f776e6572736869702070726f6f66604482015268081c1c9bdd9a59195960ba1b606482015260840162000e66565b8a604001516001600160a01b03168660018151811062003e225762003e2262005fde565b60200260200101516001600160a01b03161462003e965760405162461bcd60e51b815260206004820152602b60248201527f496e76616c6964206f7261636c65206b6579206f776e6572736869702070726f60448201526a1bd9881c1c9bdd9a59195960aa1b606482015260840162000e66565b600162003ead8c6102400151838e5f015162004f7e565b1462003f1b5760405162461bcd60e51b815260206004820152603660248201527f496e76616c696420636f6e73656e737573206b6579206f776e65727368697020604482015275383937b7b3103337b9103932b3b4b9ba3930ba34b7b760511b606482015260840162000e66565b62003f268b62004fe8565b5050505050505050505050565b5f821162003f905760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000e66565b6001600160a01b0381165f9081526029602052604090205482111562003ff95760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000e66565b6001600160a01b0381165f90815260296020526040812080548492906200402290849062005bfb565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580545f908152600394859052978820875181549088166001600160a01b0319918216178255955160018201805491909816961695909517909555905160028401555191015580549192620040ad8362005ff2565b90915550506001600160a01b038481165f818152602a6020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d91015b60405180910390a35050505050565b6001600160a01b038084165f908152602a6020526040902080549091838116911614806200428757600c820154604051631092ab9160e31b81526001600160a01b0385811660048301525f9216906384955c8890602401602060405180830381865afa1580156200418d573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620041b3919062005fa5565b9050848110156200421b5760405162461bcd60e51b815260206004820152602b60248201527f696e73756666696369656e7420756e6c6f636b6564204c6971756964204e657760448201526a746f6e2062616c616e636560a81b606482015260840162000e66565b600c83015460405163282d3fdf60e01b81526001600160a01b038681166004830152602482018890529091169063282d3fdf906044015f604051808303815f87803b15801562004269575f80fd5b505af11580156200427c573d5f803e3d5ffd5b505050505062004319565b8382600b015483600801546200429e919062005bfb565b1015620042fe5760405162461bcd60e51b815260206004820152602760248201527f696e73756666696369656e742073656c6620626f6e646564206e6577746f6e2060448201526662616c616e636560c81b606482015260840162000e66565b8382600b015f82825462004313919062005c11565b90915550505b6040805160e0810182526001600160a01b03808616825287811660208084019182528385018981525f60608601818152436080880190815260a088018381528a151560c08a019081526008805486526006909752998420985189549089166001600160a01b0319918216178a55965160018a01805491909916971696909617909655915160028701559051600386015592516004850155905160059093018054945115156101000261ff00199415159490941661ffff199095169490941792909217909255805491620043ec8362005ff2565b9190505550826001600160a01b0316856001600160a01b03167f63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc83876040516200410b9291909115158252602082015260400190565b60fa60a05f808285855af46200445a573d5f803e3d5ffd5b505050565b8154600260018083161561010002038216048251808201602081106020841001600281146200450c57600181146200453257865f526020840460205f2001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b80821015620044f55781518455600184019350602082019150620044d6565b815191036101000a908190040290915550620045a5565b60028302826020036101000a846020036101000a602089015104020185018755620045a5565b865f526020840460205f2001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b8082101562004593578151845560018401935060208201915062004574565b815191036101000a9081900402909155505b50505050505050565b5f81815260036020908152604080832060018101546001600160a01b03168452602a90925282209091601382015460ff166003811115620045f357620045f3620054f2565b146200469257600282015482546001600160a01b03165f90815260296020526040812080549091906200462890849062005c11565b909155505081546001830154600284015460138401546040517f1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f8783429462004685946001600160a01b0391821694911692909160ff90911690620061fe565b60405180910390a1505050565b805482546001600160a01b0390811691161462004787575f8082600801548360050154620046c1919062005bfb565b9050805f03620046d85783600201549150620046fe565b80846002015484600d0154620046ef91906200600d565b620046fb91906200617d565b91505b600c83015484546040516340c10f1960e01b81526001600160a01b039182166004820152602481018590529116906340c10f19906044015f604051808303815f87803b1580156200474d575f80fd5b505af115801562004760573d5f803e3d5ffd5b505050508183600d015f82825462004779919062005c11565b90915550620047a692505050565b8160020154816008015f828254620047a0919062005c11565b90915550505b8160020154816005015f828254620033ad919062005c11565b5f81815260066020908152604080832060018101546001600160a01b03168452602a9092528220600582015491929091610100900460ff16620049a3576002830154600c8301548454604051637eee288d60e01b81526001600160a01b03918216600482015260248101849052911690637eee288d906044015f604051808303815f87803b15801562004850575f80fd5b505af115801562004863573d5f803e3d5ffd5b50505050600c8301548454604051632770a7eb60e21b81526001600160a01b03918216600482015260248101849052911690639dc29fac906044015f604051808303815f87803b158015620048b6575f80fd5b505af1158015620048c9573d5f803e3d5ffd5b505050505f83600801548460050154620048e4919062005bfb565b600d850154909150620048f882846200600d565b6200490491906200617d565b92508184600d015f8282546200491b919062005bfb565b909155505060068401545f0362004939576003850183905562004961565b600684015460078501546200494f90856200600d565b6200495b91906200617d565b60038601555b82846006015f82825462004976919062005c11565b909155505060038501546007850180545f906200499590849062005c11565b9091555062004a7392505050565b5060028201546008820154811115620049bd575060088101545b81600901545f03620049d65760038301819055620049fe565b6009820154600a830154620049ec90836200600d565b620049f891906200617d565b60038401555b80826009015f82825462004a13919062005c11565b90915550506003830154600a830180545f9062004a3290849062005c11565b9250508190555080826008015f82825462004a4e919062005bfb565b90915550506002830154600b830180545f9062004a6d90849062005bfb565b90915550505b6005808401805460ff19166001179055820180548291905f9062004a9990849062005bfb565b909155505050505050565b5f818152600660205260408120600381015490910362004ac2575050565b60018101546001600160a01b03165f908152602a602052604081206005830154909190610100900460ff1662004b5b5781600701548260060154846003015462004b0d91906200600d565b62004b1991906200617d565b905080826006015f82825462004b30919062005bfb565b909155505060038301546007830180545f9062004b4f90849062005bfb565b9091555062004bbe9050565b81600a01548260090154846003015462004b7691906200600d565b62004b8291906200617d565b905080826009015f82825462004b99919062005bfb565b90915550506003830154600a830180545f9062004bb890849062005bfb565b90915550505b82546001600160a01b03165f908152602960205260408120805483929062004a9990849062005c11565b5f62004bf8826060015162005213565b6001600160a01b0390911660208401529050801562004c485760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000e66565b6020808301516001600160a01b039081165f908152602a909252604090912060010154161562004cbb5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000e66565b61271082608001511115620014935760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000e66565b6060815f0362004d3a5750506040805180820190915260018152600360fc1b602082015290565b815f5b811562004d69578062004d508162005ff2565b915062004d619050600a836200617d565b915062004d3d565b5f8167ffffffffffffffff81111562004d865762004d86620057da565b6040519080825280601f01601f19166020018201604052801562004db1576020820181803683370190505b5090505b841562004e285762004dc960018362005bfb565b915062004dd8600a866200603b565b62004de590603062005c11565b60f81b81838151811062004dfd5762004dfd62005fde565b60200101906001600160f81b03191690815f1a90535062004e20600a866200617d565b945062004db5565b949350505050565b60608162004e4081601f62005c11565b101562004e815760405162461bcd60e51b815260206004820152600e60248201526d736c6963655f6f766572666c6f7760901b604482015260640162000e66565b62004e8d828462005c11565b8451101562004ed35760405162461bcd60e51b8152602060048201526011602482015270736c6963655f6f75744f66426f756e647360781b604482015260640162000e66565b60608215801562004ef35760405191505f82526020820160405262004f3f565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101562004f2e57805183526020928301920162004f14565b5050858452601f01601f1916604052505b50949350505050565b818101805160208201516040909201519091905f1a601b81101562004f775762004f74601b826200622c565b90505b9250925092565b5f62004f896200538a565b5f85858560405160200162004fa19392919062006248565b60405160208183030381529060405290505f60fb90505f8251602062004fc8919062005c11565b90506020848285855afa62004fdb575f80fd5b5050905195945050505050565b6101808101516001600160a01b03166200506657601f545f906200500c9062004d13565b90508160200151825f01518360800151836040516200502b90620053a8565b6200503a949392919062006296565b604051809103905ff08015801562005054573d5f803e3d5ffd5b506001600160a01b0316610180830152505b60208181018051601f8054600180820183555f9283527fa03837a25210ee280c2113ff4b77ca23440b19d4866cca721c801278fd08d80790910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602a90955260409182902086518154871690851617815593519084018054861691841691909117905584015160028301805490941691161790915560608201518291906003820190620051179082620062d4565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620051e19082620062d4565b5061026082015160138201805460ff19166001836003811115620052095762005209620054f2565b0217905550505050565b5f806200521f620053b6565b5f8060ff9050604083875160208901845afa6200523a575f80fd5b50508051602090910151600160601b90910494909350915050565b6040518061028001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b03168152602001606081526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081526020015f60038111156200530e576200530e620054f2565b905290565b6040518060a001604052806005906020820280368337509192915050565b5080545f8255905f5260205f209081019062000eb59190620053d4565b5080546200535c9062005bad565b5f825580601f106200536c575050565b601f0160209004905f5260205f209081019062000eb59190620053f8565b60405180602001604052806001906020820280368337509192915050565b6117d6806200639283390190565b60405180604001604052806002906020820280368337509192915050565b80821115620053f4575f620053ea82826200534e565b50600101620053d4565b5090565b5b80821115620053f4575f8155600101620053f9565b5f5b838110156200542a57818101518382015260200162005410565b50505f910152565b5f81518084526200544b8160208601602086016200540e565b601f01601f19169290920160200192915050565b602081525f62005473602083018462005432565b9392505050565b6001600160a01b038116811462000eb5575f80fd5b5f8060408385031215620054a1575f80fd5b8235620054ae816200547a565b946020939093013593505050565b5f60208284031215620054cd575f80fd5b813562005473816200547a565b5f60208284031215620054eb575f80fd5b5035919050565b634e487b7160e01b5f52602160045260245ffd5b600481106200552357634e487b7160e01b5f52602160045260245ffd5b9052565b60208152620055426020820183516001600160a01b03169052565b5f60208301516200555e60408401826001600160a01b03169052565b5060408301516001600160a01b0381166060840152506060830151610280806080850152620055926102a085018362005432565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a06200561a818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f1901610260808801919091529091506200567e848362005432565b935080870151915050620056958286018262005506565b5090949350505050565b5f805f60608486031215620056b2575f80fd5b8335620056bf816200547a565b92506020840135620056d1816200547a565b929592945050506040919091013590565b5f60208284031215620056f3575f80fd5b813567ffffffffffffffff8111156200570a575f80fd5b8201610280818503121562005473575f80fd5b5f81518084526020808501808196508360051b810191508286015f5b8581101562005790578284038952815180516001600160a01b0316855285810151868601526040908101516060918601829052906200577b8187018362005432565b9a87019a955050509084019060010162005739565b5091979650505050505050565b8215158152604060208201525f62004e2860408301846200571d565b5f8060408385031215620057cb575f80fd5b50508035926020909101359150565b634e487b7160e01b5f52604160045260245ffd5b5f82601f830112620057fe575f80fd5b813567ffffffffffffffff808211156200581c576200581c620057da565b604051601f8301601f19908116603f01168101908282118183101715620058475762005847620057da565b8160405283815286602085880101111562005860575f80fd5b836020870160208301375f602085830101528094505050505092915050565b5f805f806080858703121562005893575f80fd5b843567ffffffffffffffff80821115620058ab575f80fd5b620058b988838901620057ee565b955060208701359150620058cd826200547a565b90935060408601359080821115620058e3575f80fd5b620058f188838901620057ee565b9350606087013591508082111562005907575f80fd5b506200591687828801620057ee565b91505092959194509250565b5f602080830181845280855180835260408601915060408160051b87010192508387015f5b828110156200597957603f198886030184526200596685835162005432565b9450928501929085019060010162005947565b5092979650505050505050565b602081525f6200547360208301846200571d565b602080825282518282018190525f9190848201906040850190845b81811015620059dc5783516001600160a01b031683529284019291840191600101620059b5565b50909695505050505050565b5f8060408385031215620059fa575f80fd5b823567ffffffffffffffff8082111562005a12575f80fd5b62005a2086838701620057ee565b9350602085013591508082111562005a36575f80fd5b5062005a4585828601620057ee565b9150509250929050565b604081525f62005a63604083018562005432565b828103602084015262005a77818562005432565b95945050505050565b5f806040838503121562005a92575f80fd5b823562005a9f816200547a565b9150602083013562005ab1816200547a565b809150509250929050565b6020808252601c908201527f76616c696461746f72206d757374206265207265676973746572656400000000604082015260600190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b60208082526018908201527f76616c696461746f72206e6f7420726567697374657265640000000000000000604082015260600190565b600181811c9082168062005bc257607f821691505b60208210810362005be157634e487b7160e01b5f52602260045260245ffd5b50919050565b634e487b7160e01b5f52601160045260245ffd5b8181038181111562000e215762000e2162005be7565b8082018082111562000e215762000e2162005be7565b5f813562000e21816200547a565b80546001600160a01b0319166001600160a01b0392909216919091179055565b5f808335601e1984360301811262005c6b575f80fd5b83018035915067ffffffffffffffff82111562005c86575f80fd5b60200191503681900382131562005c9b575f80fd5b9250929050565b601f8211156200445a575f81815260208120601f850160051c8101602086101562005cca5750805b601f850160051c820191505b8181101562005ceb5782815560010162005cd6565b505050505050565b5f19600383901b1c191660019190911b1790565b67ffffffffffffffff83111562005d225762005d22620057da565b62005d3a8362005d33835462005bad565b8362005ca2565b5f601f84116001811462005d6b575f851562005d565750838201355b62005d62868262005cf3565b84555062001590565b5f83815260209020601f19861690835b8281101562005d9d578685013582556020948501946001909201910162005d7b565b508682101562005dba575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b5f81356004811062000e21575f80fd5b6004821062005df957634e487b7160e01b5f52602160045260245ffd5b60ff1981541660ff831681178255505050565b62005e2262005e1b8362005c27565b8262005c35565b62005e3e62005e346020840162005c27565b6001830162005c35565b62005e5a62005e506040840162005c27565b6002830162005c35565b62005e69606083018362005c55565b62005e7981836003860162005d07565b50506080820135600482015560a0820135600582015560c0820135600682015560e0820135600782015561010082013560088201556101208201356009820155610140820135600a820155610160820135600b82015562005eec62005ee2610180840162005c27565b600c830162005c35565b6101a0820135600d8201556101c0820135600e8201556101e0820135600f8201556102008201356010820155610220820135601182015562005f3361024083018362005c55565b62005f4381836012860162005d07565b50506200149362005f58610260840162005dcc565b6013830162005ddc565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b5f6020828403121562005fb6575f80fd5b5051919050565b5f6020828403121562005fce575f80fd5b8151801515811462005473575f80fd5b634e487b7160e01b5f52603260045260245ffd5b5f6001820162006006576200600662005be7565b5060010190565b808202811582820484141762000e215762000e2162005be7565b634e487b7160e01b5f52601260045260245ffd5b5f826200604c576200604c62006027565b500690565b5f60018060a01b0380881683528087166020840152808616604084015260a060608401526200608460a084018662005432565b91508084166080840152509695505050505050565b818103620060a5575050565b620060b1825462005bad565b67ffffffffffffffff811115620060cc57620060cc620057da565b620060e481620060dd845462005bad565b8462005ca2565b5f601f82116001811462006115575f8315620061005750848201545b6200610c848262005cf3565b85555062001590565b5f85815260209020601f198416905f86815260209020845b838110156200614f57828601548255600195860195909101906020016200612d565b50858310156200616d57818501545f19600388901b60f8161c191681555b5050505050600190811b01905550565b5f826200618e576200618e62006027565b500490565b5f8060408385031215620061a5575f80fd5b505080516020909101519092909150565b5f8451620061c98184602089016200540e565b845190830190620061df8183602089016200540e565b8451910190620061f48183602088016200540e565b0195945050505050565b6001600160a01b03858116825284166020820152604081018390526080810162005a77606083018462005506565b60ff818116838216019081111562000e215762000e2162005be7565b5f84516200625b8184602089016200540e565b845190830190620062718183602089016200540e565b60609490941b6bffffffffffffffffffffffff19169301928352505060140192915050565b6001600160a01b03858116825284166020820152604081018390526080606082018190525f90620062ca9083018462005432565b9695505050505050565b815167ffffffffffffffff811115620062f157620062f1620057da565b6200630281620060dd845462005bad565b602080601f83116001811462006335575f8415620063205750858301515b6200632c858262005cf3565b86555062005ceb565b5f85815260208120601f198616915b82811015620063655788860151825594840194600190910190840162006344565b50858210156200616d5793909601515f1960f8600387901b161c19169092555050600190811b0190555056fe608060405234801562000010575f80fd5b50604051620017d6380380620017d6833981016040819052620000339162000149565b61271082111562000042575f80fd5b600d80546001600160a01b038087166001600160a01b031992831617909255600e805492861692909116919091179055600f8290556040516200008a90829060200162000230565b604051602081830303815290604052600b9081620000a99190620002ea565b5080604051602001620000bd919062000230565b604051602081830303815290604052600c9081620000dc9190620002ea565b50505f80546001600160a01b0319163317905550620003b2915050565b6001600160a01b03811681146200010e575f80fd5b50565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156200014157818101518382015260200162000127565b50505f910152565b5f805f80608085870312156200015d575f80fd5b84516200016a81620000f9565b60208601519094506200017d81620000f9565b6040860151606087015191945092506001600160401b0380821115620001a1575f80fd5b818701915087601f830112620001b5575f80fd5b815181811115620001ca57620001ca62000111565b604051601f8201601f19908116603f01168101908382118183101715620001f557620001f562000111565b816040528281528a60208487010111156200020e575f80fd5b6200022183602083016020880162000125565b979a9699509497505050505050565b644c4e544e2d60d81b81525f82516200025181600585016020870162000125565b9190910160050192915050565b600181811c908216806200027357607f821691505b6020821081036200029257634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620002e5575f81815260208120601f850160051c81016020861015620002c05750805b601f850160051c820191505b81811015620002e157828155600101620002cc565b5050505b505050565b81516001600160401b0381111562000306576200030662000111565b6200031e816200031784546200025e565b8462000298565b602080601f83116001811462000354575f84156200033c5750858301515b5f19600386901b1c1916600185901b178555620002e1565b5f85815260208120601f198616915b82811015620003845788860151825594840194600190910190840162000363565b5085821015620003a257878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61141680620003c05f395ff3fe608060405260043610610147575f3560e01c806359355736116100b3578063949813b81161006d578063949813b8146103ac57806395d89b41146103e05780639dc29fac146103f4578063a0ce552d14610413578063a9059cbb14610426578063dd62ed3e14610445575f80fd5b806359355736146102d25780635ea1d6f81461030657806361d027b31461031b57806370a082311461033a5780637eee288d1461036e57806384955c881461038d575f80fd5b8063282d3fdf11610104578063282d3fdf146102195780632f2c3f2e14610238578063313ce5671461024d578063372500ab146102685780633a5381b51461027c57806340c10f19146102b3575f80fd5b806306fdde031461014b578063095ea7b31461017557806318160ddd146101a4578063187cf4d7146101c257806319fac8fd146101d957806323b872dd146101fa575b5f80fd5b348015610156575f80fd5b5061015f610489565b60405161016c9190611182565b60405180910390f35b348015610180575f80fd5b5061019461018f3660046111e8565b610519565b604051901515815260200161016c565b3480156101af575f80fd5b506004545b60405190815260200161016c565b3480156101cd575f80fd5b506101b4633b9aca0081565b3480156101e4575f80fd5b506101f86101f3366004611210565b61052f565b005b348015610205575f80fd5b50610194610214366004611227565b610566565b348015610224575f80fd5b506101f86102333660046111e8565b610657565b348015610243575f80fd5b506101b461271081565b348015610258575f80fd5b506040516012815260200161016c565b348015610273575f80fd5b506101f8610739565b348015610287575f80fd5b50600d5461029b906001600160a01b031681565b6040516001600160a01b03909116815260200161016c565b3480156102be575f80fd5b506101f86102cd3660046111e8565b6108a9565b3480156102dd575f80fd5b506101b46102ec366004611260565b6001600160a01b03165f9081526002602052604090205490565b348015610311575f80fd5b506101b4600f5481565b348015610326575f80fd5b50600e5461029b906001600160a01b031681565b348015610345575f80fd5b506101b4610354366004611260565b6001600160a01b03165f9081526001602052604090205490565b348015610379575f80fd5b506101f86103883660046111e8565b61090e565b348015610398575f80fd5b506101b46103a7366004611260565b6109d1565b3480156103b7575f80fd5b506103cb6103c6366004611260565b6109fe565b6040805192835260208301919091520161016c565b3480156103eb575f80fd5b5061015f610a62565b3480156103ff575f80fd5b506101f861040e3660046111e8565b610a71565b6103cb610421366004611210565b610ace565b348015610431575f80fd5b506101946104403660046111e8565b610d78565b348015610450575f80fd5b506101b461045f366004611280565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b6060600b8054610498906112b1565b80601f01602080910402602001604051908101604052809291908181526020018280546104c4906112b1565b801561050f5780601f106104e65761010080835404028352916020019161050f565b820191905f5260205f20905b8154815290600101906020018083116104f257829003601f168201915b5050505050905090565b5f610525338484610dc3565b5060015b92915050565b5f546001600160a01b031633146105615760405162461bcd60e51b8152600401610558906112e9565b60405180910390fd5b600f55565b6001600160a01b0383165f908152600360209081526040808320338452909152812054828110156105ea5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b6064820152608401610558565b6105fe85336105f98685611345565b610dc3565b6106088584610ee6565b6106128484610fd8565b836001600160a01b0316856001600160a01b03165f805160206113c18339815191528560405161064491815260200190565b60405180910390a3506001949350505050565b5f546001600160a01b031633146106805760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f9081526002602090815260408083205460019092529091205482916106af91611345565b10156107095760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611358565b90915550505050565b5f8061074433611022565b335f818152600560209081526040808320839055600890915280822082905551939550919350909184908381818185875af1925050503d805f81146107a4576040519150601f19603f3d011682016040523d82523d5f602084013e6107a9565b606091505b50509050806107ef5760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b7321020aa2760711b6044820152606401610558565b5f5460405163a9059cbb60e01b8152336004820152602481018490526001600160a01b039091169063a9059cbb906044016020604051808303815f875af115801561083c573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610860919061136b565b9050806108a45760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b73210272a2760711b6044820152606401610558565b505050565b5f546001600160a01b031633146108d25760405162461bcd60e51b8152600401610558906112e9565b6108dc8282610fd8565b6040518181526001600160a01b038316905f905f805160206113c1833981519152906020015b60405180910390a35050565b5f546001600160a01b031633146109375760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f908152600260205260409020548111156109aa5760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611345565b6001600160a01b0381165f9081526002602090815260408083205460019092528220546105299190611345565b5f805f80610a0b856110c7565b6001600160a01b0387165f908152600560205260409020549193509150610a33908390611358565b6001600160a01b0386165f90815260086020526040902054909450610a59908290611358565b92505050915091565b6060600c8054610498906112b1565b5f546001600160a01b03163314610a9a5760405162461bcd60e51b8152600401610558906112e9565b610aa48282610ee6565b6040518181525f906001600160a01b038416905f805160206113c183398151915290602001610902565b5f805481906001600160a01b03163314610afa5760405162461bcd60e51b8152600401610558906112e9565b600f5434905f9061271090610b0f908461138a565b610b1991906113a1565b905081811115610b665760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610b708183611345565b600e546040519193506001600160a01b0316906108fc9083905f818181858888f193505050503d805f8114610bc0576040519150601f19603f3d011682016040523d82523d5f602084013e610bc5565b606091505b5050505f612710600f5487610bda919061138a565b610be491906113a1565b905085811115610c315760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610c3b8187611345565b5f54600e5460405163a9059cbb60e01b81526001600160a01b03918216600482015260248101859052929850169063a9059cbb906044016020604051808303815f875af1158015610c8e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610cb2919061136b565b506004545f90610cc6633b9aca008661138a565b610cd091906113a1565b905080600754610ce09190611358565b6007556004545f90610cf6633b9aca008a61138a565b610d0091906113a1565b905080600a54610d109190611358565b600a556004545f90633b9aca0090610d28908561138a565b610d3291906113a1565b90505f633b9aca0060045484610d48919061138a565b610d5291906113a1565b9050610d5e8287611358565b610d688287611358565b9850985050505050505050915091565b5f610d833383610ee6565b610d8d8383610fd8565b6040518281526001600160a01b0384169033905f805160206113c18339815191529060200160405180910390a350600192915050565b6001600160a01b038316610e255760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610558565b6001600160a01b038216610e865760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610558565b6001600160a01b038381165f8181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610eef82611022565b50506001600160a01b0382165f90815260016020908152604080832054600290925290912054610f1f9082611345565b821115610f6e5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e647300000000006044820152606401610558565b610f788282611345565b6001600160a01b0384165f90815260016020526040902055808203610fbd576001600160a01b0383165f90815260066020908152604080832083905560099091528120555b8160045f828254610fce9190611345565b9091555050505050565b610fe182611022565b50506001600160a01b0382165f908152600160205260408120805483929061100a908490611358565b925050819055508060045f8282546107309190611358565b5f805f8061102f856110c7565b6001600160a01b0387165f908152600560205260409020549193509150611057908390611358565b6001600160a01b0386165f90815260056020908152604080832084905560075460068352818420556008909152902054909450611095908290611358565b6001600160a01b039095165f908152600860209081526040808320889055600a54600990925290912055509193915050565b6001600160a01b0381165f9081526001602052604081205481908082036110f357505f93849350915050565b6001600160a01b0384165f908152600660205260408120546007546111189190611345565b6001600160a01b0386165f90815260096020526040812054600a5492935090916111429190611345565b9050633b9aca00611153848461138a565b61115d91906113a1565b9450633b9aca0061116e848361138a565b61117891906113a1565b9350505050915091565b5f6020808352835180828501525f5b818110156111ad57858101830151858201604001528201611191565b505f604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146111e3575f80fd5b919050565b5f80604083850312156111f9575f80fd5b611202836111cd565b946020939093013593505050565b5f60208284031215611220575f80fd5b5035919050565b5f805f60608486031215611239575f80fd5b611242846111cd565b9250611250602085016111cd565b9150604084013590509250925092565b5f60208284031215611270575f80fd5b611279826111cd565b9392505050565b5f8060408385031215611291575f80fd5b61129a836111cd565b91506112a8602084016111cd565b90509250929050565b600181811c908216806112c557607f821691505b6020821081036112e357634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561052957610529611331565b8082018082111561052957610529611331565b5f6020828403121561137b575f80fd5b81518015158114611279575f80fd5b808202811582820484141761052957610529611331565b5f826113bb57634e487b7160e01b5f52601260045260245ffd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220fcf7fe14f5b055004178aa5c27d7983104c707dfee4e3b218ce03dfc1bc4bea864736f6c63430008150033a26469706673582212205f4113bd21fa31fa68f74638c7188ee318bf2bddeaec3d29e65fb108f6035e8264736f6c63430008150033608060405234801562000010575f80fd5b50604051620017d6380380620017d6833981016040819052620000339162000149565b61271082111562000042575f80fd5b600d80546001600160a01b038087166001600160a01b031992831617909255600e805492861692909116919091179055600f8290556040516200008a90829060200162000230565b604051602081830303815290604052600b9081620000a99190620002ea565b5080604051602001620000bd919062000230565b604051602081830303815290604052600c9081620000dc9190620002ea565b50505f80546001600160a01b0319163317905550620003b2915050565b6001600160a01b03811681146200010e575f80fd5b50565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156200014157818101518382015260200162000127565b50505f910152565b5f805f80608085870312156200015d575f80fd5b84516200016a81620000f9565b60208601519094506200017d81620000f9565b6040860151606087015191945092506001600160401b0380821115620001a1575f80fd5b818701915087601f830112620001b5575f80fd5b815181811115620001ca57620001ca62000111565b604051601f8201601f19908116603f01168101908382118183101715620001f557620001f562000111565b816040528281528a60208487010111156200020e575f80fd5b6200022183602083016020880162000125565b979a9699509497505050505050565b644c4e544e2d60d81b81525f82516200025181600585016020870162000125565b9190910160050192915050565b600181811c908216806200027357607f821691505b6020821081036200029257634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620002e5575f81815260208120601f850160051c81016020861015620002c05750805b601f850160051c820191505b81811015620002e157828155600101620002cc565b5050505b505050565b81516001600160401b0381111562000306576200030662000111565b6200031e816200031784546200025e565b8462000298565b602080601f83116001811462000354575f84156200033c5750858301515b5f19600386901b1c1916600185901b178555620002e1565b5f85815260208120601f198616915b82811015620003845788860151825594840194600190910190840162000363565b5085821015620003a257878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61141680620003c05f395ff3fe608060405260043610610147575f3560e01c806359355736116100b3578063949813b81161006d578063949813b8146103ac57806395d89b41146103e05780639dc29fac146103f4578063a0ce552d14610413578063a9059cbb14610426578063dd62ed3e14610445575f80fd5b806359355736146102d25780635ea1d6f81461030657806361d027b31461031b57806370a082311461033a5780637eee288d1461036e57806384955c881461038d575f80fd5b8063282d3fdf11610104578063282d3fdf146102195780632f2c3f2e14610238578063313ce5671461024d578063372500ab146102685780633a5381b51461027c57806340c10f19146102b3575f80fd5b806306fdde031461014b578063095ea7b31461017557806318160ddd146101a4578063187cf4d7146101c257806319fac8fd146101d957806323b872dd146101fa575b5f80fd5b348015610156575f80fd5b5061015f610489565b60405161016c9190611182565b60405180910390f35b348015610180575f80fd5b5061019461018f3660046111e8565b610519565b604051901515815260200161016c565b3480156101af575f80fd5b506004545b60405190815260200161016c565b3480156101cd575f80fd5b506101b4633b9aca0081565b3480156101e4575f80fd5b506101f86101f3366004611210565b61052f565b005b348015610205575f80fd5b50610194610214366004611227565b610566565b348015610224575f80fd5b506101f86102333660046111e8565b610657565b348015610243575f80fd5b506101b461271081565b348015610258575f80fd5b506040516012815260200161016c565b348015610273575f80fd5b506101f8610739565b348015610287575f80fd5b50600d5461029b906001600160a01b031681565b6040516001600160a01b03909116815260200161016c565b3480156102be575f80fd5b506101f86102cd3660046111e8565b6108a9565b3480156102dd575f80fd5b506101b46102ec366004611260565b6001600160a01b03165f9081526002602052604090205490565b348015610311575f80fd5b506101b4600f5481565b348015610326575f80fd5b50600e5461029b906001600160a01b031681565b348015610345575f80fd5b506101b4610354366004611260565b6001600160a01b03165f9081526001602052604090205490565b348015610379575f80fd5b506101f86103883660046111e8565b61090e565b348015610398575f80fd5b506101b46103a7366004611260565b6109d1565b3480156103b7575f80fd5b506103cb6103c6366004611260565b6109fe565b6040805192835260208301919091520161016c565b3480156103eb575f80fd5b5061015f610a62565b3480156103ff575f80fd5b506101f861040e3660046111e8565b610a71565b6103cb610421366004611210565b610ace565b348015610431575f80fd5b506101946104403660046111e8565b610d78565b348015610450575f80fd5b506101b461045f366004611280565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b6060600b8054610498906112b1565b80601f01602080910402602001604051908101604052809291908181526020018280546104c4906112b1565b801561050f5780601f106104e65761010080835404028352916020019161050f565b820191905f5260205f20905b8154815290600101906020018083116104f257829003601f168201915b5050505050905090565b5f610525338484610dc3565b5060015b92915050565b5f546001600160a01b031633146105615760405162461bcd60e51b8152600401610558906112e9565b60405180910390fd5b600f55565b6001600160a01b0383165f908152600360209081526040808320338452909152812054828110156105ea5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b6064820152608401610558565b6105fe85336105f98685611345565b610dc3565b6106088584610ee6565b6106128484610fd8565b836001600160a01b0316856001600160a01b03165f805160206113c18339815191528560405161064491815260200190565b60405180910390a3506001949350505050565b5f546001600160a01b031633146106805760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f9081526002602090815260408083205460019092529091205482916106af91611345565b10156107095760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611358565b90915550505050565b5f8061074433611022565b335f818152600560209081526040808320839055600890915280822082905551939550919350909184908381818185875af1925050503d805f81146107a4576040519150601f19603f3d011682016040523d82523d5f602084013e6107a9565b606091505b50509050806107ef5760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b7321020aa2760711b6044820152606401610558565b5f5460405163a9059cbb60e01b8152336004820152602481018490526001600160a01b039091169063a9059cbb906044016020604051808303815f875af115801561083c573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610860919061136b565b9050806108a45760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b73210272a2760711b6044820152606401610558565b505050565b5f546001600160a01b031633146108d25760405162461bcd60e51b8152600401610558906112e9565b6108dc8282610fd8565b6040518181526001600160a01b038316905f905f805160206113c1833981519152906020015b60405180910390a35050565b5f546001600160a01b031633146109375760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f908152600260205260409020548111156109aa5760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611345565b6001600160a01b0381165f9081526002602090815260408083205460019092528220546105299190611345565b5f805f80610a0b856110c7565b6001600160a01b0387165f908152600560205260409020549193509150610a33908390611358565b6001600160a01b0386165f90815260086020526040902054909450610a59908290611358565b92505050915091565b6060600c8054610498906112b1565b5f546001600160a01b03163314610a9a5760405162461bcd60e51b8152600401610558906112e9565b610aa48282610ee6565b6040518181525f906001600160a01b038416905f805160206113c183398151915290602001610902565b5f805481906001600160a01b03163314610afa5760405162461bcd60e51b8152600401610558906112e9565b600f5434905f9061271090610b0f908461138a565b610b1991906113a1565b905081811115610b665760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610b708183611345565b600e546040519193506001600160a01b0316906108fc9083905f818181858888f193505050503d805f8114610bc0576040519150601f19603f3d011682016040523d82523d5f602084013e610bc5565b606091505b5050505f612710600f5487610bda919061138a565b610be491906113a1565b905085811115610c315760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610c3b8187611345565b5f54600e5460405163a9059cbb60e01b81526001600160a01b03918216600482015260248101859052929850169063a9059cbb906044016020604051808303815f875af1158015610c8e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610cb2919061136b565b506004545f90610cc6633b9aca008661138a565b610cd091906113a1565b905080600754610ce09190611358565b6007556004545f90610cf6633b9aca008a61138a565b610d0091906113a1565b905080600a54610d109190611358565b600a556004545f90633b9aca0090610d28908561138a565b610d3291906113a1565b90505f633b9aca0060045484610d48919061138a565b610d5291906113a1565b9050610d5e8287611358565b610d688287611358565b9850985050505050505050915091565b5f610d833383610ee6565b610d8d8383610fd8565b6040518281526001600160a01b0384169033905f805160206113c18339815191529060200160405180910390a350600192915050565b6001600160a01b038316610e255760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610558565b6001600160a01b038216610e865760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610558565b6001600160a01b038381165f8181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610eef82611022565b50506001600160a01b0382165f90815260016020908152604080832054600290925290912054610f1f9082611345565b821115610f6e5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e647300000000006044820152606401610558565b610f788282611345565b6001600160a01b0384165f90815260016020526040902055808203610fbd576001600160a01b0383165f90815260066020908152604080832083905560099091528120555b8160045f828254610fce9190611345565b9091555050505050565b610fe182611022565b50506001600160a01b0382165f908152600160205260408120805483929061100a908490611358565b925050819055508060045f8282546107309190611358565b5f805f8061102f856110c7565b6001600160a01b0387165f908152600560205260409020549193509150611057908390611358565b6001600160a01b0386165f90815260056020908152604080832084905560075460068352818420556008909152902054909450611095908290611358565b6001600160a01b039095165f908152600860209081526040808320889055600a54600990925290912055509193915050565b6001600160a01b0381165f9081526001602052604081205481908082036110f357505f93849350915050565b6001600160a01b0384165f908152600660205260408120546007546111189190611345565b6001600160a01b0386165f90815260096020526040812054600a5492935090916111429190611345565b9050633b9aca00611153848461138a565b61115d91906113a1565b9450633b9aca0061116e848361138a565b61117891906113a1565b9350505050915091565b5f6020808352835180828501525f5b818110156111ad57858101830151858201604001528201611191565b505f604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146111e3575f80fd5b919050565b5f80604083850312156111f9575f80fd5b611202836111cd565b946020939093013593505050565b5f60208284031215611220575f80fd5b5035919050565b5f805f60608486031215611239575f80fd5b611242846111cd565b9250611250602085016111cd565b9150604084013590509250925092565b5f60208284031215611270575f80fd5b611279826111cd565b9392505050565b5f8060408385031215611291575f80fd5b61129a836111cd565b91506112a8602084016111cd565b90509250929050565b600181811c908216806112c557607f821691505b6020821081036112e357634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561052957610529611331565b8082018082111561052957610529611331565b5f6020828403121561137b575f80fd5b81518015158114611279575f80fd5b808202811582820484141761052957610529611331565b5f826113bb57634e487b7160e01b5f52601260045260245ffd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220fcf7fe14f5b055004178aa5c27d7983104c707dfee4e3b218ce03dfc1bc4bea864736f6c63430008150033",
}

// AutonityABI is the input ABI used to generate the binding from.
// Deprecated: Use AutonityMetaData.ABI instead.
var AutonityABI = AutonityMetaData.ABI

// Deprecated: Use AutonityMetaData.Sigs instead.
// AutonityFuncSigs maps the 4-byte function signature to its string representation.
var AutonityFuncSigs = AutonityMetaData.Sigs

// AutonityBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AutonityMetaData.Bin instead.
var AutonityBin = AutonityMetaData.Bin

// DeployAutonity deploys a new Ethereum contract, binding an instance of Autonity to it.
func (r *runner) deployAutonity(opts *runOptions, _validators []AutonityValidator, _config AutonityConfig) (common.Address, uint64, *Autonity, error) {
	parsed, err := AutonityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(AutonityBin), _validators, _config)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &Autonity{AutonityCaller: AutonityCaller{contract: c}, AutonityTransactor: AutonityTransactor{contract: c}, AutonityFilterer: AutonityFilterer{contract: c}}, nil
}

// Autonity is an auto generated Go binding around an Ethereum contract.
type Autonity struct {
	AutonityCaller     // Read-only binding to the contract
	AutonityTransactor // Write-only binding to the contract
	AutonityFilterer   // Log filterer for contract events
}

// AutonityCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type AutonityCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// AutonityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AutonityTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// AutonityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AutonityFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Autonity *AutonityCaller) COMMISSIONRATEPRECISION(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "COMMISSION_RATE_PRECISION")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Autonity *AutonityCaller) Allowance(opts *runOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// AtnTotalRedistributed is a free data retrieval call binding the contract method 0xe98712b2.
//
// Solidity: function atnTotalRedistributed() view returns(uint256)
func (_Autonity *AutonityCaller) AtnTotalRedistributed(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "atnTotalRedistributed")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_Autonity *AutonityCaller) BalanceOf(opts *runOptions, _addr common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "balanceOf", _addr)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,uint256,address) policy, (address,address,address,address,address,address,address) contracts, (address,uint256,uint256,uint256) protocol, uint256 contractVersion)
func (_Autonity *AutonityCaller) Config(opts *runOptions) (struct {
	Policy          AutonityPolicy
	Contracts       AutonityContracts
	Protocol        AutonityProtocol
	ContractVersion *big.Int
}, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "config")

	outstruct := new(struct {
		Policy          AutonityPolicy
		Contracts       AutonityContracts
		Protocol        AutonityProtocol
		ContractVersion *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.Policy = *abi.ConvertType(out[0], new(AutonityPolicy)).(*AutonityPolicy)
	outstruct.Contracts = *abi.ConvertType(out[1], new(AutonityContracts)).(*AutonityContracts)
	outstruct.Protocol = *abi.ConvertType(out[2], new(AutonityProtocol)).(*AutonityProtocol)
	outstruct.ContractVersion = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Autonity *AutonityCaller) Decimals(opts *runOptions) (uint8, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "decimals")

	if err != nil {
		return *new(uint8), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, consumed, err

}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() view returns(address)
func (_Autonity *AutonityCaller) Deployer(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "deployer")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// EpochID is a free data retrieval call binding the contract method 0xc9d97af4.
//
// Solidity: function epochID() view returns(uint256)
func (_Autonity *AutonityCaller) EpochID(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "epochID")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// EpochReward is a free data retrieval call binding the contract method 0x1604e416.
//
// Solidity: function epochReward() view returns(uint256)
func (_Autonity *AutonityCaller) EpochReward(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "epochReward")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// EpochTotalBondedStake is a free data retrieval call binding the contract method 0x9c98e471.
//
// Solidity: function epochTotalBondedStake() view returns(uint256)
func (_Autonity *AutonityCaller) EpochTotalBondedStake(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "epochTotalBondedStake")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetBlockPeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getBlockPeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256,bytes)[])
func (_Autonity *AutonityCaller) GetCommittee(opts *runOptions) ([]AutonityCommitteeMember, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getCommittee")

	if err != nil {
		return *new([]AutonityCommitteeMember), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]AutonityCommitteeMember)).(*[]AutonityCommitteeMember)
	return out0, consumed, err

}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_Autonity *AutonityCaller) GetCommitteeEnodes(opts *runOptions) ([]string, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getCommitteeEnodes")

	if err != nil {
		return *new([]string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, consumed, err

}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_Autonity *AutonityCaller) GetEpochFromBlock(opts *runOptions, _block *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getEpochFromBlock", _block)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetEpochPeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getEpochPeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCaller) GetLastEpochBlock(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getLastEpochBlock")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_Autonity *AutonityCaller) GetMaxCommitteeSize(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getMaxCommitteeSize")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_Autonity *AutonityCaller) GetMinimumBaseFee(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getMinimumBaseFee")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Autonity *AutonityCaller) GetNewContract(opts *runOptions) ([]byte, string, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getNewContract")

	if err != nil {
		return *new([]byte), *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	return out0, out1, consumed, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_Autonity *AutonityCaller) GetOperator(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getOperator")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_Autonity *AutonityCaller) GetOracle(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getOracle")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetProposer is a free data retrieval call binding the contract method 0x5f7d3949.
//
// Solidity: function getProposer(uint256 height, uint256 round) view returns(address)
func (_Autonity *AutonityCaller) GetProposer(opts *runOptions, height *big.Int, round *big.Int) (common.Address, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getProposer", height, round)

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_Autonity *AutonityCaller) GetTreasuryAccount(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getTreasuryAccount")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_Autonity *AutonityCaller) GetTreasuryFee(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getTreasuryFee")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetUnbondingPeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getUnbondingPeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8))
func (_Autonity *AutonityCaller) GetValidator(opts *runOptions, _addr common.Address) (AutonityValidator, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getValidator", _addr)

	if err != nil {
		return *new(AutonityValidator), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(AutonityValidator)).(*AutonityValidator)
	return out0, consumed, err

}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Autonity *AutonityCaller) GetValidators(opts *runOptions) ([]common.Address, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getValidators")

	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_Autonity *AutonityCaller) GetVersion(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getVersion")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// InflationReserve is a free data retrieval call binding the contract method 0x9ac31074.
//
// Solidity: function inflationReserve() view returns(uint256)
func (_Autonity *AutonityCaller) InflationReserve(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "inflationReserve")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCaller) LastEpochBlock(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "lastEpochBlock")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Autonity *AutonityCaller) Name(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "name")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Autonity *AutonityCaller) Symbol(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "symbol")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Autonity *AutonityCaller) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_Autonity *AutonityTransactor) ActivateValidator(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "activateValidator", _address)
	return consumed, err
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactor) Approve(opts *runOptions, spender common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "approve", spender, amount)
	return consumed, err
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Bond(opts *runOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "bond", _validator, _amount)
	return consumed, err
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Burn(opts *runOptions, _addr common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "burn", _addr, _amount)
	return consumed, err
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_Autonity *AutonityTransactor) ChangeCommissionRate(opts *runOptions, _validator common.Address, _rate *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "changeCommissionRate", _validator, _rate)
	return consumed, err
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Autonity *AutonityTransactor) CompleteContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "completeContractUpgrade")
	return consumed, err
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[])
func (_Autonity *AutonityTransactor) ComputeCommittee(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "computeCommittee")
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256,bytes)[])
func (_Autonity *AutonityTransactor) Finalize(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "finalize")
	return consumed, err
}

// FinalizeInitialization is a paid mutator transaction binding the contract method 0xd861b0e8.
//
// Solidity: function finalizeInitialization() returns()
func (_Autonity *AutonityTransactor) FinalizeInitialization(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "finalizeInitialization")
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Mint(opts *runOptions, _addr common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "mint", _addr, _amount)
	return consumed, err
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_Autonity *AutonityTransactor) PauseValidator(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "pauseValidator", _address)
	return consumed, err
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x84467fdb.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_Autonity *AutonityTransactor) RegisterValidator(opts *runOptions, _enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "registerValidator", _enode, _oracleAddress, _consensusKey, _signatures)
	return consumed, err
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Autonity *AutonityTransactor) ResetContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "resetContractUpgrade")
	return consumed, err
}

// SetAccountabilityContract is a paid mutator transaction binding the contract method 0x1250a28d.
//
// Solidity: function setAccountabilityContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetAccountabilityContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setAccountabilityContract", _address)
	return consumed, err
}

// SetAcuContract is a paid mutator transaction binding the contract method 0xd372c07e.
//
// Solidity: function setAcuContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetAcuContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setAcuContract", _address)
	return consumed, err
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_Autonity *AutonityTransactor) SetCommitteeSize(opts *runOptions, _size *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setCommitteeSize", _size)
	return consumed, err
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactor) SetEpochPeriod(opts *runOptions, _period *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setEpochPeriod", _period)
	return consumed, err
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_Autonity *AutonityTransactor) SetMinimumBaseFee(opts *runOptions, _price *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setMinimumBaseFee", _price)
	return consumed, err
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_Autonity *AutonityTransactor) SetOperatorAccount(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setOperatorAccount", _account)
	return consumed, err
}

// SetOracleContract is a paid mutator transaction binding the contract method 0x496ccd9b.
//
// Solidity: function setOracleContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetOracleContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setOracleContract", _address)
	return consumed, err
}

// SetStabilizationContract is a paid mutator transaction binding the contract method 0xcfd19fb9.
//
// Solidity: function setStabilizationContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetStabilizationContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setStabilizationContract", _address)
	return consumed, err
}

// SetSupplyControlContract is a paid mutator transaction binding the contract method 0xb3ecbadd.
//
// Solidity: function setSupplyControlContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetSupplyControlContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setSupplyControlContract", _address)
	return consumed, err
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_Autonity *AutonityTransactor) SetTreasuryAccount(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setTreasuryAccount", _account)
	return consumed, err
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_Autonity *AutonityTransactor) SetTreasuryFee(opts *runOptions, _treasuryFee *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setTreasuryFee", _treasuryFee)
	return consumed, err
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactor) SetUnbondingPeriod(opts *runOptions, _period *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setUnbondingPeriod", _period)
	return consumed, err
}

// SetUpgradeManagerContract is a paid mutator transaction binding the contract method 0xceaad455.
//
// Solidity: function setUpgradeManagerContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetUpgradeManagerContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setUpgradeManagerContract", _address)
	return consumed, err
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonityTransactor) Transfer(opts *runOptions, _recipient common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "transfer", _recipient, _amount)
	return consumed, err
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonityTransactor) TransferFrom(opts *runOptions, _sender common.Address, _recipient common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "transferFrom", _sender, _recipient, _amount)
	return consumed, err
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Unbond(opts *runOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "unbond", _validator, _amount)
	return consumed, err
}

// UpdateValidatorAndTransferSlashedFunds is a paid mutator transaction binding the contract method 0x35be16e0.
//
// Solidity: function updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8) _val) returns()
func (_Autonity *AutonityTransactor) UpdateValidatorAndTransferSlashedFunds(opts *runOptions, _val AutonityValidator) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "updateValidatorAndTransferSlashedFunds", _val)
	return consumed, err
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Autonity *AutonityTransactor) UpgradeContract(opts *runOptions, _bytecode []byte, _abi string) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "upgradeContract", _bytecode, _abi)
	return consumed, err
}

// Fallback is a paid mutator transaction binding the contract fallback function.
// WARNING! UNTESTED
// Solidity: fallback() payable returns()
func (_Autonity *AutonityTransactor) Fallback(opts *runOptions, calldata []byte) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "", calldata)
	return consumed, err
}

// Receive is a paid mutator transaction binding the contract receive function.
// WARNING! UNTESTED
// Solidity: receive() payable returns()
func (_Autonity *AutonityTransactor) Receive(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "")
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// AutonityActivatedValidatorIterator is returned from FilterActivatedValidator and is used to iterate over the raw logs and unpacked data for ActivatedValidator events raised by the Autonity contract.
		type AutonityActivatedValidatorIterator struct {
			Event *AutonityActivatedValidator // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityActivatedValidatorIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityActivatedValidator)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityActivatedValidator)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityActivatedValidatorIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityActivatedValidatorIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityActivatedValidator represents a ActivatedValidator event raised by the Autonity contract.
		type AutonityActivatedValidator struct {
			Treasury common.Address;
			Addr common.Address;
			EffectiveBlock *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterActivatedValidator is a free log retrieval operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
		//
		// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
 		func (_Autonity *AutonityFilterer) FilterActivatedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*AutonityActivatedValidatorIterator, error) {

			var treasuryRule []interface{}
			for _, treasuryItem := range treasury {
				treasuryRule = append(treasuryRule, treasuryItem)
			}
			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.FilterLogs(opts, "ActivatedValidator", treasuryRule, addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityActivatedValidatorIterator{contract: _Autonity.contract, event: "ActivatedValidator", logs: logs, sub: sub}, nil
 		}

		// WatchActivatedValidator is a free log subscription operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
		//
		// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
		func (_Autonity *AutonityFilterer) WatchActivatedValidator(opts *bind.WatchOpts, sink chan<- *AutonityActivatedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

			var treasuryRule []interface{}
			for _, treasuryItem := range treasury {
				treasuryRule = append(treasuryRule, treasuryItem)
			}
			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.WatchLogs(opts, "ActivatedValidator", treasuryRule, addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityActivatedValidator)
						if err := _Autonity.contract.UnpackLog(event, "ActivatedValidator", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseActivatedValidator is a log parse operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
		//
		// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
		func (_Autonity *AutonityFilterer) ParseActivatedValidator(log types.Log) (*AutonityActivatedValidator, error) {
			event := new(AutonityActivatedValidator)
			if err := _Autonity.contract.UnpackLog(event, "ActivatedValidator", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Autonity contract.
		type AutonityApprovalIterator struct {
			Event *AutonityApproval // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityApprovalIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityApproval)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityApproval)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityApprovalIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityApprovalIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityApproval represents a Approval event raised by the Autonity contract.
		type AutonityApproval struct {
			Owner common.Address;
			Spender common.Address;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
 		func (_Autonity *AutonityFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AutonityApprovalIterator, error) {

			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var spenderRule []interface{}
			for _, spenderItem := range spender {
				spenderRule = append(spenderRule, spenderItem)
			}


			logs, sub, err := _Autonity.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
			if err != nil {
				return nil, err
			}
			return &AutonityApprovalIterator{contract: _Autonity.contract, event: "Approval", logs: logs, sub: sub}, nil
 		}

		// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
		func (_Autonity *AutonityFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AutonityApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var spenderRule []interface{}
			for _, spenderItem := range spender {
				spenderRule = append(spenderRule, spenderItem)
			}


			logs, sub, err := _Autonity.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityApproval)
						if err := _Autonity.contract.UnpackLog(event, "Approval", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
		func (_Autonity *AutonityFilterer) ParseApproval(log types.Log) (*AutonityApproval, error) {
			event := new(AutonityApproval)
			if err := _Autonity.contract.UnpackLog(event, "Approval", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityBondingRejectedIterator is returned from FilterBondingRejected and is used to iterate over the raw logs and unpacked data for BondingRejected events raised by the Autonity contract.
		type AutonityBondingRejectedIterator struct {
			Event *AutonityBondingRejected // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityBondingRejectedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityBondingRejected)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityBondingRejected)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityBondingRejectedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityBondingRejectedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityBondingRejected represents a BondingRejected event raised by the Autonity contract.
		type AutonityBondingRejected struct {
			Delegator common.Address;
			Delegatee common.Address;
			Amount *big.Int;
			State uint8;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterBondingRejected is a free log retrieval operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
		//
		// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount, uint8 state)
 		func (_Autonity *AutonityFilterer) FilterBondingRejected(opts *bind.FilterOpts) (*AutonityBondingRejectedIterator, error) {






			logs, sub, err := _Autonity.contract.FilterLogs(opts, "BondingRejected")
			if err != nil {
				return nil, err
			}
			return &AutonityBondingRejectedIterator{contract: _Autonity.contract, event: "BondingRejected", logs: logs, sub: sub}, nil
 		}

		// WatchBondingRejected is a free log subscription operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
		//
		// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount, uint8 state)
		func (_Autonity *AutonityFilterer) WatchBondingRejected(opts *bind.WatchOpts, sink chan<- *AutonityBondingRejected) (event.Subscription, error) {






			logs, sub, err := _Autonity.contract.WatchLogs(opts, "BondingRejected")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityBondingRejected)
						if err := _Autonity.contract.UnpackLog(event, "BondingRejected", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseBondingRejected is a log parse operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
		//
		// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount, uint8 state)
		func (_Autonity *AutonityFilterer) ParseBondingRejected(log types.Log) (*AutonityBondingRejected, error) {
			event := new(AutonityBondingRejected)
			if err := _Autonity.contract.UnpackLog(event, "BondingRejected", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityBurnedStakeIterator is returned from FilterBurnedStake and is used to iterate over the raw logs and unpacked data for BurnedStake events raised by the Autonity contract.
		type AutonityBurnedStakeIterator struct {
			Event *AutonityBurnedStake // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityBurnedStakeIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityBurnedStake)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityBurnedStake)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityBurnedStakeIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityBurnedStakeIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityBurnedStake represents a BurnedStake event raised by the Autonity contract.
		type AutonityBurnedStake struct {
			Addr common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterBurnedStake is a free log retrieval operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
		//
		// Solidity: event BurnedStake(address indexed addr, uint256 amount)
 		func (_Autonity *AutonityFilterer) FilterBurnedStake(opts *bind.FilterOpts, addr []common.Address) (*AutonityBurnedStakeIterator, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.FilterLogs(opts, "BurnedStake", addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityBurnedStakeIterator{contract: _Autonity.contract, event: "BurnedStake", logs: logs, sub: sub}, nil
 		}

		// WatchBurnedStake is a free log subscription operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
		//
		// Solidity: event BurnedStake(address indexed addr, uint256 amount)
		func (_Autonity *AutonityFilterer) WatchBurnedStake(opts *bind.WatchOpts, sink chan<- *AutonityBurnedStake, addr []common.Address) (event.Subscription, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.WatchLogs(opts, "BurnedStake", addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityBurnedStake)
						if err := _Autonity.contract.UnpackLog(event, "BurnedStake", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseBurnedStake is a log parse operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
		//
		// Solidity: event BurnedStake(address indexed addr, uint256 amount)
		func (_Autonity *AutonityFilterer) ParseBurnedStake(log types.Log) (*AutonityBurnedStake, error) {
			event := new(AutonityBurnedStake)
			if err := _Autonity.contract.UnpackLog(event, "BurnedStake", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityCommissionRateChangeIterator is returned from FilterCommissionRateChange and is used to iterate over the raw logs and unpacked data for CommissionRateChange events raised by the Autonity contract.
		type AutonityCommissionRateChangeIterator struct {
			Event *AutonityCommissionRateChange // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityCommissionRateChangeIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityCommissionRateChange)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityCommissionRateChange)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityCommissionRateChangeIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityCommissionRateChangeIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityCommissionRateChange represents a CommissionRateChange event raised by the Autonity contract.
		type AutonityCommissionRateChange struct {
			Validator common.Address;
			Rate *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterCommissionRateChange is a free log retrieval operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
		//
		// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
 		func (_Autonity *AutonityFilterer) FilterCommissionRateChange(opts *bind.FilterOpts, validator []common.Address) (*AutonityCommissionRateChangeIterator, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}


			logs, sub, err := _Autonity.contract.FilterLogs(opts, "CommissionRateChange", validatorRule)
			if err != nil {
				return nil, err
			}
			return &AutonityCommissionRateChangeIterator{contract: _Autonity.contract, event: "CommissionRateChange", logs: logs, sub: sub}, nil
 		}

		// WatchCommissionRateChange is a free log subscription operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
		//
		// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
		func (_Autonity *AutonityFilterer) WatchCommissionRateChange(opts *bind.WatchOpts, sink chan<- *AutonityCommissionRateChange, validator []common.Address) (event.Subscription, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}


			logs, sub, err := _Autonity.contract.WatchLogs(opts, "CommissionRateChange", validatorRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityCommissionRateChange)
						if err := _Autonity.contract.UnpackLog(event, "CommissionRateChange", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseCommissionRateChange is a log parse operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
		//
		// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
		func (_Autonity *AutonityFilterer) ParseCommissionRateChange(log types.Log) (*AutonityCommissionRateChange, error) {
			event := new(AutonityCommissionRateChange)
			if err := _Autonity.contract.UnpackLog(event, "CommissionRateChange", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityEpochPeriodUpdatedIterator is returned from FilterEpochPeriodUpdated and is used to iterate over the raw logs and unpacked data for EpochPeriodUpdated events raised by the Autonity contract.
		type AutonityEpochPeriodUpdatedIterator struct {
			Event *AutonityEpochPeriodUpdated // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityEpochPeriodUpdatedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityEpochPeriodUpdated)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityEpochPeriodUpdated)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityEpochPeriodUpdatedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityEpochPeriodUpdatedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityEpochPeriodUpdated represents a EpochPeriodUpdated event raised by the Autonity contract.
		type AutonityEpochPeriodUpdated struct {
			Period *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterEpochPeriodUpdated is a free log retrieval operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
		//
		// Solidity: event EpochPeriodUpdated(uint256 period)
 		func (_Autonity *AutonityFilterer) FilterEpochPeriodUpdated(opts *bind.FilterOpts) (*AutonityEpochPeriodUpdatedIterator, error) {



			logs, sub, err := _Autonity.contract.FilterLogs(opts, "EpochPeriodUpdated")
			if err != nil {
				return nil, err
			}
			return &AutonityEpochPeriodUpdatedIterator{contract: _Autonity.contract, event: "EpochPeriodUpdated", logs: logs, sub: sub}, nil
 		}

		// WatchEpochPeriodUpdated is a free log subscription operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
		//
		// Solidity: event EpochPeriodUpdated(uint256 period)
		func (_Autonity *AutonityFilterer) WatchEpochPeriodUpdated(opts *bind.WatchOpts, sink chan<- *AutonityEpochPeriodUpdated) (event.Subscription, error) {



			logs, sub, err := _Autonity.contract.WatchLogs(opts, "EpochPeriodUpdated")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityEpochPeriodUpdated)
						if err := _Autonity.contract.UnpackLog(event, "EpochPeriodUpdated", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseEpochPeriodUpdated is a log parse operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
		//
		// Solidity: event EpochPeriodUpdated(uint256 period)
		func (_Autonity *AutonityFilterer) ParseEpochPeriodUpdated(log types.Log) (*AutonityEpochPeriodUpdated, error) {
			event := new(AutonityEpochPeriodUpdated)
			if err := _Autonity.contract.UnpackLog(event, "EpochPeriodUpdated", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityMinimumBaseFeeUpdatedIterator is returned from FilterMinimumBaseFeeUpdated and is used to iterate over the raw logs and unpacked data for MinimumBaseFeeUpdated events raised by the Autonity contract.
		type AutonityMinimumBaseFeeUpdatedIterator struct {
			Event *AutonityMinimumBaseFeeUpdated // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityMinimumBaseFeeUpdatedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityMinimumBaseFeeUpdated)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityMinimumBaseFeeUpdated)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityMinimumBaseFeeUpdatedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityMinimumBaseFeeUpdatedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityMinimumBaseFeeUpdated represents a MinimumBaseFeeUpdated event raised by the Autonity contract.
		type AutonityMinimumBaseFeeUpdated struct {
			GasPrice *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterMinimumBaseFeeUpdated is a free log retrieval operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
		//
		// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
 		func (_Autonity *AutonityFilterer) FilterMinimumBaseFeeUpdated(opts *bind.FilterOpts) (*AutonityMinimumBaseFeeUpdatedIterator, error) {



			logs, sub, err := _Autonity.contract.FilterLogs(opts, "MinimumBaseFeeUpdated")
			if err != nil {
				return nil, err
			}
			return &AutonityMinimumBaseFeeUpdatedIterator{contract: _Autonity.contract, event: "MinimumBaseFeeUpdated", logs: logs, sub: sub}, nil
 		}

		// WatchMinimumBaseFeeUpdated is a free log subscription operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
		//
		// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
		func (_Autonity *AutonityFilterer) WatchMinimumBaseFeeUpdated(opts *bind.WatchOpts, sink chan<- *AutonityMinimumBaseFeeUpdated) (event.Subscription, error) {



			logs, sub, err := _Autonity.contract.WatchLogs(opts, "MinimumBaseFeeUpdated")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityMinimumBaseFeeUpdated)
						if err := _Autonity.contract.UnpackLog(event, "MinimumBaseFeeUpdated", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseMinimumBaseFeeUpdated is a log parse operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
		//
		// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
		func (_Autonity *AutonityFilterer) ParseMinimumBaseFeeUpdated(log types.Log) (*AutonityMinimumBaseFeeUpdated, error) {
			event := new(AutonityMinimumBaseFeeUpdated)
			if err := _Autonity.contract.UnpackLog(event, "MinimumBaseFeeUpdated", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityMintedStakeIterator is returned from FilterMintedStake and is used to iterate over the raw logs and unpacked data for MintedStake events raised by the Autonity contract.
		type AutonityMintedStakeIterator struct {
			Event *AutonityMintedStake // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityMintedStakeIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityMintedStake)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityMintedStake)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityMintedStakeIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityMintedStakeIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityMintedStake represents a MintedStake event raised by the Autonity contract.
		type AutonityMintedStake struct {
			Addr common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterMintedStake is a free log retrieval operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
		//
		// Solidity: event MintedStake(address indexed addr, uint256 amount)
 		func (_Autonity *AutonityFilterer) FilterMintedStake(opts *bind.FilterOpts, addr []common.Address) (*AutonityMintedStakeIterator, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.FilterLogs(opts, "MintedStake", addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityMintedStakeIterator{contract: _Autonity.contract, event: "MintedStake", logs: logs, sub: sub}, nil
 		}

		// WatchMintedStake is a free log subscription operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
		//
		// Solidity: event MintedStake(address indexed addr, uint256 amount)
		func (_Autonity *AutonityFilterer) WatchMintedStake(opts *bind.WatchOpts, sink chan<- *AutonityMintedStake, addr []common.Address) (event.Subscription, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.WatchLogs(opts, "MintedStake", addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityMintedStake)
						if err := _Autonity.contract.UnpackLog(event, "MintedStake", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseMintedStake is a log parse operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
		//
		// Solidity: event MintedStake(address indexed addr, uint256 amount)
		func (_Autonity *AutonityFilterer) ParseMintedStake(log types.Log) (*AutonityMintedStake, error) {
			event := new(AutonityMintedStake)
			if err := _Autonity.contract.UnpackLog(event, "MintedStake", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityNewBondingRequestIterator is returned from FilterNewBondingRequest and is used to iterate over the raw logs and unpacked data for NewBondingRequest events raised by the Autonity contract.
		type AutonityNewBondingRequestIterator struct {
			Event *AutonityNewBondingRequest // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityNewBondingRequestIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityNewBondingRequest)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityNewBondingRequest)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityNewBondingRequestIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityNewBondingRequestIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityNewBondingRequest represents a NewBondingRequest event raised by the Autonity contract.
		type AutonityNewBondingRequest struct {
			Validator common.Address;
			Delegator common.Address;
			SelfBonded bool;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewBondingRequest is a free log retrieval operation binding the contract event 0xc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d.
		//
		// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
 		func (_Autonity *AutonityFilterer) FilterNewBondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*AutonityNewBondingRequestIterator, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}
			var delegatorRule []interface{}
			for _, delegatorItem := range delegator {
				delegatorRule = append(delegatorRule, delegatorItem)
			}



			logs, sub, err := _Autonity.contract.FilterLogs(opts, "NewBondingRequest", validatorRule, delegatorRule)
			if err != nil {
				return nil, err
			}
			return &AutonityNewBondingRequestIterator{contract: _Autonity.contract, event: "NewBondingRequest", logs: logs, sub: sub}, nil
 		}

		// WatchNewBondingRequest is a free log subscription operation binding the contract event 0xc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d.
		//
		// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
		func (_Autonity *AutonityFilterer) WatchNewBondingRequest(opts *bind.WatchOpts, sink chan<- *AutonityNewBondingRequest, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}
			var delegatorRule []interface{}
			for _, delegatorItem := range delegator {
				delegatorRule = append(delegatorRule, delegatorItem)
			}



			logs, sub, err := _Autonity.contract.WatchLogs(opts, "NewBondingRequest", validatorRule, delegatorRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityNewBondingRequest)
						if err := _Autonity.contract.UnpackLog(event, "NewBondingRequest", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewBondingRequest is a log parse operation binding the contract event 0xc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d.
		//
		// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
		func (_Autonity *AutonityFilterer) ParseNewBondingRequest(log types.Log) (*AutonityNewBondingRequest, error) {
			event := new(AutonityNewBondingRequest)
			if err := _Autonity.contract.UnpackLog(event, "NewBondingRequest", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the Autonity contract.
		type AutonityNewEpochIterator struct {
			Event *AutonityNewEpoch // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityNewEpochIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityNewEpoch)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityNewEpoch)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityNewEpochIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityNewEpochIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityNewEpoch represents a NewEpoch event raised by the Autonity contract.
		type AutonityNewEpoch struct {
			Epoch *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewEpoch is a free log retrieval operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
		//
		// Solidity: event NewEpoch(uint256 epoch)
 		func (_Autonity *AutonityFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*AutonityNewEpochIterator, error) {



			logs, sub, err := _Autonity.contract.FilterLogs(opts, "NewEpoch")
			if err != nil {
				return nil, err
			}
			return &AutonityNewEpochIterator{contract: _Autonity.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
 		}

		// WatchNewEpoch is a free log subscription operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
		//
		// Solidity: event NewEpoch(uint256 epoch)
		func (_Autonity *AutonityFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *AutonityNewEpoch) (event.Subscription, error) {



			logs, sub, err := _Autonity.contract.WatchLogs(opts, "NewEpoch")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityNewEpoch)
						if err := _Autonity.contract.UnpackLog(event, "NewEpoch", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewEpoch is a log parse operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
		//
		// Solidity: event NewEpoch(uint256 epoch)
		func (_Autonity *AutonityFilterer) ParseNewEpoch(log types.Log) (*AutonityNewEpoch, error) {
			event := new(AutonityNewEpoch)
			if err := _Autonity.contract.UnpackLog(event, "NewEpoch", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityNewUnbondingRequestIterator is returned from FilterNewUnbondingRequest and is used to iterate over the raw logs and unpacked data for NewUnbondingRequest events raised by the Autonity contract.
		type AutonityNewUnbondingRequestIterator struct {
			Event *AutonityNewUnbondingRequest // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityNewUnbondingRequestIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityNewUnbondingRequest)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityNewUnbondingRequest)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityNewUnbondingRequestIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityNewUnbondingRequestIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityNewUnbondingRequest represents a NewUnbondingRequest event raised by the Autonity contract.
		type AutonityNewUnbondingRequest struct {
			Validator common.Address;
			Delegator common.Address;
			SelfBonded bool;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewUnbondingRequest is a free log retrieval operation binding the contract event 0x63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc.
		//
		// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
 		func (_Autonity *AutonityFilterer) FilterNewUnbondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*AutonityNewUnbondingRequestIterator, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}
			var delegatorRule []interface{}
			for _, delegatorItem := range delegator {
				delegatorRule = append(delegatorRule, delegatorItem)
			}



			logs, sub, err := _Autonity.contract.FilterLogs(opts, "NewUnbondingRequest", validatorRule, delegatorRule)
			if err != nil {
				return nil, err
			}
			return &AutonityNewUnbondingRequestIterator{contract: _Autonity.contract, event: "NewUnbondingRequest", logs: logs, sub: sub}, nil
 		}

		// WatchNewUnbondingRequest is a free log subscription operation binding the contract event 0x63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc.
		//
		// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
		func (_Autonity *AutonityFilterer) WatchNewUnbondingRequest(opts *bind.WatchOpts, sink chan<- *AutonityNewUnbondingRequest, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}
			var delegatorRule []interface{}
			for _, delegatorItem := range delegator {
				delegatorRule = append(delegatorRule, delegatorItem)
			}



			logs, sub, err := _Autonity.contract.WatchLogs(opts, "NewUnbondingRequest", validatorRule, delegatorRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityNewUnbondingRequest)
						if err := _Autonity.contract.UnpackLog(event, "NewUnbondingRequest", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewUnbondingRequest is a log parse operation binding the contract event 0x63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc.
		//
		// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
		func (_Autonity *AutonityFilterer) ParseNewUnbondingRequest(log types.Log) (*AutonityNewUnbondingRequest, error) {
			event := new(AutonityNewUnbondingRequest)
			if err := _Autonity.contract.UnpackLog(event, "NewUnbondingRequest", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityPausedValidatorIterator is returned from FilterPausedValidator and is used to iterate over the raw logs and unpacked data for PausedValidator events raised by the Autonity contract.
		type AutonityPausedValidatorIterator struct {
			Event *AutonityPausedValidator // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityPausedValidatorIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityPausedValidator)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityPausedValidator)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityPausedValidatorIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityPausedValidatorIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityPausedValidator represents a PausedValidator event raised by the Autonity contract.
		type AutonityPausedValidator struct {
			Treasury common.Address;
			Addr common.Address;
			EffectiveBlock *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterPausedValidator is a free log retrieval operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
		//
		// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
 		func (_Autonity *AutonityFilterer) FilterPausedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*AutonityPausedValidatorIterator, error) {

			var treasuryRule []interface{}
			for _, treasuryItem := range treasury {
				treasuryRule = append(treasuryRule, treasuryItem)
			}
			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.FilterLogs(opts, "PausedValidator", treasuryRule, addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityPausedValidatorIterator{contract: _Autonity.contract, event: "PausedValidator", logs: logs, sub: sub}, nil
 		}

		// WatchPausedValidator is a free log subscription operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
		//
		// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
		func (_Autonity *AutonityFilterer) WatchPausedValidator(opts *bind.WatchOpts, sink chan<- *AutonityPausedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

			var treasuryRule []interface{}
			for _, treasuryItem := range treasury {
				treasuryRule = append(treasuryRule, treasuryItem)
			}
			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.WatchLogs(opts, "PausedValidator", treasuryRule, addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityPausedValidator)
						if err := _Autonity.contract.UnpackLog(event, "PausedValidator", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParsePausedValidator is a log parse operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
		//
		// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
		func (_Autonity *AutonityFilterer) ParsePausedValidator(log types.Log) (*AutonityPausedValidator, error) {
			event := new(AutonityPausedValidator)
			if err := _Autonity.contract.UnpackLog(event, "PausedValidator", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityRegisteredValidatorIterator is returned from FilterRegisteredValidator and is used to iterate over the raw logs and unpacked data for RegisteredValidator events raised by the Autonity contract.
		type AutonityRegisteredValidatorIterator struct {
			Event *AutonityRegisteredValidator // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityRegisteredValidatorIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityRegisteredValidator)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityRegisteredValidator)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityRegisteredValidatorIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityRegisteredValidatorIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityRegisteredValidator represents a RegisteredValidator event raised by the Autonity contract.
		type AutonityRegisteredValidator struct {
			Treasury common.Address;
			Addr common.Address;
			OracleAddress common.Address;
			Enode string;
			LiquidContract common.Address;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterRegisteredValidator is a free log retrieval operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
		//
		// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
 		func (_Autonity *AutonityFilterer) FilterRegisteredValidator(opts *bind.FilterOpts) (*AutonityRegisteredValidatorIterator, error) {







			logs, sub, err := _Autonity.contract.FilterLogs(opts, "RegisteredValidator")
			if err != nil {
				return nil, err
			}
			return &AutonityRegisteredValidatorIterator{contract: _Autonity.contract, event: "RegisteredValidator", logs: logs, sub: sub}, nil
 		}

		// WatchRegisteredValidator is a free log subscription operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
		//
		// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
		func (_Autonity *AutonityFilterer) WatchRegisteredValidator(opts *bind.WatchOpts, sink chan<- *AutonityRegisteredValidator) (event.Subscription, error) {







			logs, sub, err := _Autonity.contract.WatchLogs(opts, "RegisteredValidator")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityRegisteredValidator)
						if err := _Autonity.contract.UnpackLog(event, "RegisteredValidator", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseRegisteredValidator is a log parse operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
		//
		// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
		func (_Autonity *AutonityFilterer) ParseRegisteredValidator(log types.Log) (*AutonityRegisteredValidator, error) {
			event := new(AutonityRegisteredValidator)
			if err := _Autonity.contract.UnpackLog(event, "RegisteredValidator", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityRewardedIterator is returned from FilterRewarded and is used to iterate over the raw logs and unpacked data for Rewarded events raised by the Autonity contract.
		type AutonityRewardedIterator struct {
			Event *AutonityRewarded // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityRewardedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityRewarded)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityRewarded)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityRewardedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityRewardedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityRewarded represents a Rewarded event raised by the Autonity contract.
		type AutonityRewarded struct {
			Addr common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterRewarded is a free log retrieval operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
		//
		// Solidity: event Rewarded(address indexed addr, uint256 amount)
 		func (_Autonity *AutonityFilterer) FilterRewarded(opts *bind.FilterOpts, addr []common.Address) (*AutonityRewardedIterator, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.FilterLogs(opts, "Rewarded", addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityRewardedIterator{contract: _Autonity.contract, event: "Rewarded", logs: logs, sub: sub}, nil
 		}

		// WatchRewarded is a free log subscription operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
		//
		// Solidity: event Rewarded(address indexed addr, uint256 amount)
		func (_Autonity *AutonityFilterer) WatchRewarded(opts *bind.WatchOpts, sink chan<- *AutonityRewarded, addr []common.Address) (event.Subscription, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _Autonity.contract.WatchLogs(opts, "Rewarded", addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityRewarded)
						if err := _Autonity.contract.UnpackLog(event, "Rewarded", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseRewarded is a log parse operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
		//
		// Solidity: event Rewarded(address indexed addr, uint256 amount)
		func (_Autonity *AutonityFilterer) ParseRewarded(log types.Log) (*AutonityRewarded, error) {
			event := new(AutonityRewarded)
			if err := _Autonity.contract.UnpackLog(event, "Rewarded", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Autonity contract.
		type AutonityTransferIterator struct {
			Event *AutonityTransfer // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityTransferIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityTransfer)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityTransfer)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityTransferIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityTransferIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityTransfer represents a Transfer event raised by the Autonity contract.
		type AutonityTransfer struct {
			From common.Address;
			To common.Address;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
 		func (_Autonity *AutonityFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutonityTransferIterator, error) {

			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}


			logs, sub, err := _Autonity.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
			if err != nil {
				return nil, err
			}
			return &AutonityTransferIterator{contract: _Autonity.contract, event: "Transfer", logs: logs, sub: sub}, nil
 		}

		// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
		func (_Autonity *AutonityFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AutonityTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}


			logs, sub, err := _Autonity.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityTransfer)
						if err := _Autonity.contract.UnpackLog(event, "Transfer", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
		func (_Autonity *AutonityFilterer) ParseTransfer(log types.Log) (*AutonityTransfer, error) {
			event := new(AutonityTransfer)
			if err := _Autonity.contract.UnpackLog(event, "Transfer", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// AutonityUpgradeTestMetaData contains all meta data concerning the AutonityUpgradeTest contract.
var AutonityUpgradeTestMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"ActivatedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"BondingRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewUnbondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"atnTotalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initialInflationReserve\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"},{\"internalType\":\"contractIInflationController\",\"name\":\"inflationControllerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeInitialization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getEpochFromBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inflationReserve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_signatures\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAccountability\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAccountabilityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIACU\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAcuContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setOracleContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIStabilization\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setStabilizationContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISupplyControl\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSupplyControlContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractUpgradeManager\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setUpgradeManagerContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"_val\",\"type\":\"tuple\"}],\"name\":\"updateValidatorAndTransferSlashedFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"b46e5520": "activateValidator(address)",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"e98712b2": "atnTotalRedistributed()",
		"70a08231": "balanceOf(address)",
		"a515366a": "bond(address,uint256)",
		"9dc29fac": "burn(address,uint256)",
		"852c4849": "changeCommissionRate(address,uint256)",
		"872cf059": "completeContractUpgrade()",
		"ae1f5fa0": "computeCommittee()",
		"79502c55": "config()",
		"313ce567": "decimals()",
		"d5f39488": "deployer()",
		"c9d97af4": "epochID()",
		"1604e416": "epochReward()",
		"9c98e471": "epochTotalBondedStake()",
		"4bb278f3": "finalize()",
		"d861b0e8": "finalizeInitialization()",
		"43645969": "getBlockPeriod()",
		"ab8f6ffe": "getCommittee()",
		"a8b2216e": "getCommitteeEnodes()",
		"96b477cb": "getEpochFromBlock(uint256)",
		"dfb1a4d2": "getEpochPeriod()",
		"731b3a03": "getLastEpochBlock()",
		"819b6463": "getMaxCommitteeSize()",
		"11220633": "getMinimumBaseFee()",
		"b66b3e79": "getNewContract()",
		"e7f43c68": "getOperator()",
		"833b1fce": "getOracle()",
		"5f7d3949": "getProposer(uint256,uint256)",
		"f7866ee3": "getTreasuryAccount()",
		"29070c6d": "getTreasuryFee()",
		"6fd2c80b": "getUnbondingPeriod()",
		"1904bb2e": "getValidator(address)",
		"b7ab4db5": "getValidators()",
		"0d8e6e2c": "getVersion()",
		"9ac31074": "inflationReserve()",
		"c2362dd5": "lastEpochBlock()",
		"40c10f19": "mint(address,uint256)",
		"06fdde03": "name()",
		"0ae65e7a": "pauseValidator(address)",
		"84467fdb": "registerValidator(string,address,bytes,bytes)",
		"cf9c5719": "resetContractUpgrade()",
		"1250a28d": "setAccountabilityContract(address)",
		"d372c07e": "setAcuContract(address)",
		"8bac7dad": "setCommitteeSize(uint256)",
		"6b5f444c": "setEpochPeriod(uint256)",
		"cb696f54": "setMinimumBaseFee(uint256)",
		"520fdbbc": "setOperatorAccount(address)",
		"496ccd9b": "setOracleContract(address)",
		"cfd19fb9": "setStabilizationContract(address)",
		"b3ecbadd": "setSupplyControlContract(address)",
		"d886f8a2": "setTreasuryAccount(address)",
		"77e741c7": "setTreasuryFee(uint256)",
		"114eaf55": "setUnbondingPeriod(uint256)",
		"ceaad455": "setUpgradeManagerContract(address)",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"a5d059ca": "unbond(address,uint256)",
		"35be16e0": "updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8))",
		"b2ea9adb": "upgradeContract(bytes,string)",
	},
	Bin: "0x60806040525f600b555f600c5534801562000018575f80fd5b50604080515f80825260208201909252906200004d565b6200003962000d59565b8152602001906001900390816200002f5790505b506040805161014081018252600d546080808301918252600e5460a080850191909152600f5460c08086019190915260105460e0808701919091526011546101008701526012546001600160a01b039081166101208801529486528651908101875260135485168152601454851660208281019190915260155486168289015260165486166060808401919091526017548716838701526018548716948301949094526019548616928201929092528186015285519283018652601a549093168252601b5492820192909252601c5481850152601d548183015292820192909252601e54918101829052905f036200015d57602d80546001600160a01b031916331790556200015d82826200017a565b5050601e5460010362000174576200017462000544565b6200113c565b80518051600d55602080820151600e55604080830151600f55606080840151601055608080850151601181905560a095860151601280546001600160a01b03199081166001600160a01b0393841617909155868901518051601380548416918516919091179055808801516014805484169185169190911790558087015160158054841691851691909117905580860151601680548416918516919091179055938401516017805483169184169190911790559683015160188054891691831691909117905560c090920151601980548816918416919091179055838701518051601a805490981693169290921790955592830151601b5590820151601c5590810151601d55820151601e55602c555f5b82518110156200053f575f838281518110620002ab57620002ab62000eb0565b602002602001015160a0015190505f848381518110620002cf57620002cf62000eb0565b60200260200101516101a00181815250505f848381518110620002f657620002f662000eb0565b602002602001015161018001906001600160a01b031690816001600160a01b0316815250505f84838151811062000331576200033162000eb0565b602002602001015160a00181815250505f84838151811062000357576200035762000eb0565b60209081029190910101516101c00152600f54845185908490811062000381576200038162000eb0565b602002602001015160800181815250505f848381518110620003a757620003a762000eb0565b602002602001015161026001906003811115620003c857620003c862000ec4565b90816003811115620003de57620003de62000ec4565b815250505f848381518110620003f857620003f862000eb0565b60200260200101516101600181815250506200043684838151811062000422576200042262000eb0565b6020026020010151620006ac60201b60201c565b620004638483815181106200044f576200044f62000eb0565b6020026020010151620007e160201b60201c565b8060295f8685815181106200047c576200047c62000eb0565b60200260200101515f01516001600160a01b03166001600160a01b031681526020019081526020015f205f828254620004b6919062000eec565b9250508190555080602b5f828254620004d0919062000eec565b9250508190555062000529848381518110620004f057620004f062000eb0565b6020026020010151602001518286858151811062000512576200051262000eb0565b60200260200101515f015162000a0c60201b60201c565b5080620005368162000f08565b9150506200028b565b505050565b602d546001600160a01b03163314620005b05760405162461bcd60e51b815260206004820152602360248201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60448201526218dbdb60ea1b60648201526084015b60405180910390fd5b6002602a5f601f600181548110620005cc57620005cc62000eb0565b5f9182526020808320909101546001600160a01b03168352820192909252604001812060050180549091906200060490849062000f37565b925050819055506002602a5f601f60018154811062000627576200062762000eb0565b5f9182526020808320909101546001600160a01b03168352820192909252604001812060080180549091906200065f90849062000f37565b90915550506002601e55601a546001600160a01b03165f9081526029602052604081206103e8905562000693908062000e17565b620006a060015f62000e17565b6002805460ff19169055565b5f620006c2826060015162000bf260201b60201c565b6001600160a01b03909116602084015290508015620007125760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b6044820152606401620005a7565b6020808301516001600160a01b039081165f908152602a9092526040909120600101541615620007855760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c72656164792072656769737465726564000000006044820152606401620005a7565b61271082608001511115620007dd5760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e20726174650000000000000000006044820152606401620005a7565b5050565b6101808101516001600160a01b03166200085f57601f545f90620008059062000c3d565b90508160200151825f0151836080015183604051620008249062000e56565b62000833949392919062000f4d565b604051809103905ff0801580156200084d573d5f803e3d5ffd5b506001600160a01b0316610180830152505b60208181018051601f8054600180820183555f9283527fa03837a25210ee280c2113ff4b77ca23440b19d4866cca721c801278fd08d80790910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602a9095526040918290208651815487169085161781559351908401805486169184169190911790558401516002830180549094169116179091556060820151829190600382019062000910908262001048565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620009da908262001048565b5061026082015160138201805460ff1916600183600381111562000a025762000a0262000ec4565b0217905550505050565b5f821162000a695760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b6064820152608401620005a7565b6001600160a01b0381165f9081526029602052604090205482111562000ad25760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e636500000000006044820152606401620005a7565b6001600160a01b0381165f908152602960205260408120805484929062000afb90849062001110565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580545f908152600394859052978820875181549088166001600160a01b031991821617825595516001820180549190981696169590951790955590516002840155519101558054919262000b868362000f08565b90915550506001600160a01b038481165f818152602a6020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d910160405180910390a35050505050565b5f8062000bfe62000e64565b5f8060ff9050604083875160208901845afa62000c19575f80fd5b505080516020909101516c0100000000000000000000000090910494909350915050565b6060815f0362000c645750506040805180820190915260018152600360fc1b602082015290565b815f5b811562000c93578062000c7a8162000f08565b915062000c8b9050600a8362000f37565b915062000c67565b5f816001600160401b0381111562000caf5762000caf62000e9c565b6040519080825280601f01601f19166020018201604052801562000cda576020820181803683370190505b5090505b841562000d515762000cf260018362001110565b915062000d01600a8662001126565b62000d0e90603062000eec565b60f81b81838151811062000d265762000d2662000eb0565b60200101906001600160f81b03191690815f1a90535062000d49600a8662000f37565b945062000cde565b949350505050565b6040518061028001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b03168152602001606081526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081526020015f600381111562000e125762000e1262000ec4565b905290565b50805462000e259062000fbd565b5f825580601f1062000e35575050565b601f0160209004905f5260205f209081019062000e53919062000e82565b50565b6117d68062008d5d83390190565b60405180604001604052806002906020820280368337509192915050565b5b8082111562000e98575f815560010162000e83565b5090565b634e487b7160e01b5f52604160045260245ffd5b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b8082018082111562000f025762000f0262000ed8565b92915050565b5f6001820162000f1c5762000f1c62000ed8565b5060010190565b634e487b7160e01b5f52601260045260245ffd5b5f8262000f485762000f4862000f23565b500490565b5f60018060a01b03808716835260208187168185015285604085015260806060850152845191508160808501525f5b8281101562000f9a5785810182015185820160a00152810162000f7c565b50505f60a0828501015260a0601f19601f83011684010191505095945050505050565b600181811c9082168062000fd257607f821691505b60208210810362000ff157634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200053f575f81815260208120601f850160051c810160208610156200101f5750805b601f850160051c820191505b8181101562001040578281556001016200102b565b505050505050565b81516001600160401b0381111562001064576200106462000e9c565b6200107c8162001075845462000fbd565b8462000ff7565b602080601f831160018114620010b2575f84156200109a5750858301515b5f19600386901b1c1916600185901b17855562001040565b5f85815260208120601f198616915b82811015620010e257888601518255948401946001909101908401620010c1565b50858210156200110057878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b8181038181111562000f025762000f0262000ed8565b5f8262001137576200113762000f23565b500690565b617c13806200114a5f395ff3fe608060405260043610620003e6575f3560e01c80638bac7dad1162000202578063b7ab4db51162000116578063d5f3948811620000a6578063dfb1a4d21162000074578063dfb1a4d21462000da4578063e7f43c681462000dba578063e98712b21462000dd9578063f7866ee31462000df057005b8063d5f394881462000d00578063d861b0e81462000d21578063d886f8a21462000d38578063dd62ed3e1462000d5c57005b8063ceaad45511620000e4578063ceaad4551462000c7d578063cf9c57191462000ca1578063cfd19fb91462000cb8578063d372c07e1462000cdc57005b8063b7ab4db51462000c14578063c2362dd51462000c2b578063c9d97af41462000c42578063cb696f541462000c5957005b8063a8b2216e1162000192578063b2ea9adb1162000160578063b2ea9adb1462000b81578063b3ecbadd1462000ba5578063b46e55201462000bc9578063b66b3e791462000bed57005b8063a8b2216e1462000aeb578063a9059cbb1462000b11578063ab8f6ffe1462000b35578063ae1f5fa01462000b5b57005b80639c98e47111620001d05780639c98e4711462000a685780639dc29fac1462000a7f578063a515366a1462000aa3578063a5d059ca1462000ac757005b80638bac7dad14620009d157806395d89b4114620009f557806396b477cb1462000a225780639ac310741462000a5157005b80634364596911620002fa578063731b3a03116200028a578063833b1fce1162000258578063833b1fce146200095357806384467fdb1462000972578063852c48491462000996578063872cf05914620009ba57005b8063731b3a03146200075457806377e741c7146200076a57806379502c55146200078e578063819b6463146200093d57005b80635f7d394911620002c85780635f7d394914620006a55780636b5f444c14620006e25780636fd2c80b146200070657806370a08231146200071c57005b8063436459691462000620578063496ccd9b14620006365780634bb278f3146200065a578063520fdbbc146200068157005b806318160ddd11620003765780632f2c3f2e11620003445780632f2c3f2e14620005a4578063313ce56714620005bb57806335be16e014620005d857806340c10f1914620005fc57005b806318160ddd14620005215780631904bb2e146200053757806323b872dd146200056a57806329070c6d146200058e57005b80631122063311620003b45780631122063314620004ac578063114eaf5514620004c25780631250a28d14620004e65780631604e416146200050a57005b806306fdde0314620003f0578063095ea7b314620004335780630ae65e7a14620004685780630d8e6e2c146200048c57005b36620003ee57005b005b348015620003fc575f80fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b6040516200042a9190620054d5565b60405180910390f35b3480156200043f575f80fd5b50620004576200045136600462005505565b62000e0f565b60405190151581526020016200042a565b34801562000474575f80fd5b50620003ee6200048636600462005532565b62000e27565b34801562000498575f80fd5b50601e545b6040519081526020016200042a565b348015620004b8575f80fd5b50600e546200049d565b348015620004ce575f80fd5b50620003ee620004e036600462005550565b62000eb8565b348015620004f2575f80fd5b50620003ee6200050436600462005532565b62000eea565b34801562000516575f80fd5b506200049d60265481565b3480156200052d575f80fd5b50602b546200049d565b34801562000543575f80fd5b506200055b6200055536600462005532565b62000f39565b6040516200042a91906200559d565b34801562000576575f80fd5b50620004576200058836600462005715565b620011e1565b3480156200059a575f80fd5b50600d546200049d565b348015620005b0575f80fd5b506200049d61271081565b348015620005c7575f80fd5b50604051601281526020016200042a565b348015620005e4575f80fd5b50620003ee620005f636600462005758565b62001287565b34801562000608575f80fd5b50620003ee6200061a36600462005505565b6200145a565b3480156200062c575f80fd5b50601c546200049d565b34801562000642575f80fd5b50620003ee6200065436600462005532565b62001497565b34801562000666575f80fd5b506200067162001597565b6040516200042a92919062005813565b3480156200068d575f80fd5b50620003ee6200069f36600462005532565b62001a0b565b348015620006b1575f80fd5b50620006c9620006c33660046200582f565b62001bf2565b6040516001600160a01b0390911681526020016200042a565b348015620006ee575f80fd5b50620003ee6200070036600462005550565b62001dfe565b34801562000712575f80fd5b506010546200049d565b34801562000728575f80fd5b506200049d6200073a36600462005532565b6001600160a01b03165f9081526029602052604090205490565b34801562000760575f80fd5b506022546200049d565b34801562000776575f80fd5b50620003ee6200078836600462005550565b62001f7f565b3480156200079a575f80fd5b506040805160c08082018352600d548252600e54602080840191909152600f54838501526010546060808501919091526011546080808601919091526012546001600160a01b0390811660a080880191909152875160e081018952601354831681526014548316818701526015548316818a015260165483168186015260175483168185015260185483169181019190915260195482169581019590955286519182018752601a54168152601b5492810192909252601c5494820194909452601d5493810193909352601e546200086f939084565b6040805185518152602080870151818301528683015182840152606080880151818401526080808901518185015260a0988901516001600160a01b039081168a8601528851811660c080870191909152848a0151821660e0870152898701518216610100870152838a0151821661012087015291890151811661014086015298880151891661016085015290960151871661018083015284519096166101a0820152948301516101c08601528201516101e0850152910151610200830152610220820152610240016200042a565b34801562000949575f80fd5b50601d546200049d565b3480156200095f575f80fd5b506014546001600160a01b0316620006c9565b3480156200097e575f80fd5b50620003ee62000990366004620058f5565b62001fb1565b348015620009a2575f80fd5b50620003ee620009b436600462005505565b620020d0565b348015620009c6575f80fd5b50620003ee6200225e565b348015620009dd575f80fd5b50620003ee620009ef36600462005550565b6200229a565b34801562000a01575f80fd5b50604080518082019091526003815262272a2760e91b60208201526200041b565b34801562000a2e575f80fd5b506200049d62000a4036600462005550565b5f9081526021602052604090205490565b34801562000a5d575f80fd5b506200049d602c5481565b34801562000a74575f80fd5b506200049d60235481565b34801562000a8b575f80fd5b50620003ee62000a9d36600462005505565b6200231d565b34801562000aaf575f80fd5b50620003ee62000ac136600462005505565b62002437565b34801562000ad3575f80fd5b50620003ee62000ae536600462005505565b62002504565b34801562000af7575f80fd5b5062000b0262002599565b6040516200042a919062005998565b34801562000b1d575f80fd5b506200045762000b2f36600462005505565b62002677565b34801562000b41575f80fd5b5062000b4c620026ce565b6040516200042a9190620059fc565b34801562000b67575f80fd5b5062000b72620027dc565b6040516200042a919062005a10565b34801562000b8d575f80fd5b50620003ee62000b9f36600462005a5e565b62002a1c565b34801562000bb1575f80fd5b50620003ee62000bc336600462005532565b62002a62565b34801562000bd5575f80fd5b50620003ee62000be736600462005532565b62002ab1565b34801562000bf9575f80fd5b5062000c0462002d31565b6040516200042a92919062005ac5565b34801562000c20575f80fd5b5062000b7262002e63565b34801562000c37575f80fd5b506200049d60225481565b34801562000c4e575f80fd5b506200049d60205481565b34801562000c65575f80fd5b50620003ee62000c7736600462005550565b62002ec5565b34801562000c89575f80fd5b50620003ee62000c9b36600462005532565b62002f28565b34801562000cad575f80fd5b50620003ee62002f77565b34801562000cc4575f80fd5b50620003ee62000cd636600462005532565b62002fc9565b34801562000ce8575f80fd5b50620003ee62000cfa36600462005532565b62003018565b34801562000d0c575f80fd5b50602d54620006c9906001600160a01b031681565b34801562000d2d575f80fd5b50620003ee62003067565b34801562000d44575f80fd5b50620003ee62000d5636600462005532565b620030a8565b34801562000d68575f80fd5b506200049d62000d7a36600462005af6565b6001600160a01b039182165f90815260286020908152604080832093909416825291909152205490565b34801562000db0575f80fd5b50601b546200049d565b34801562000dc6575f80fd5b50601a546001600160a01b0316620006c9565b34801562000de5575f80fd5b506200049d60255481565b34801562000dfc575f80fd5b506012546001600160a01b0316620006c9565b5f62000e1d338484620030f7565b5060015b92915050565b6001600160a01b038082165f818152602a60205260409020600101549091161462000e6f5760405162461bcd60e51b815260040162000e669062005b32565b60405180910390fd5b6001600160a01b038181165f908152602a602052604090205416331462000eaa5760405162461bcd60e51b815260040162000e669062005b69565b62000eb5816200321f565b50565b601a546001600160a01b0316331462000ee55760405162461bcd60e51b815260040162000e669062005bb5565b601055565b601a546001600160a01b0316331462000f175760405162461bcd60e51b815260040162000e669062005bb5565b601380546001600160a01b0319166001600160a01b0392909216919091179055565b62000f43620052cb565b6001600160a01b038083165f818152602a60205260409020600101549091161462000f825760405162461bcd60e51b815260040162000e669062005bec565b6001600160a01b038083165f908152602a6020908152604091829020825161028081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000fe49062005c23565b80601f0160208091040260200160405190810160405280929190818152602001828054620010129062005c23565b8015620010615780601f10620010375761010080835404028352916020019162001061565b820191905f5260205f20905b8154815290600101906020018083116200104357829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546080820152600882015460a0820152600982015460c0820152600a82015460e0820152600b820154610100820152600c8201546001600160a01b0316610120820152600d820154610140820152600e820154610160820152600f82015461018082015260108201546101a082015260118201546101c08201526012820180546101e0909201916200111d9062005c23565b80601f01602080910402602001604051908101604052809291908181526020018280546200114b9062005c23565b80156200119a5780601f1062001170576101008083540402835291602001916200119a565b820191905f5260205f20905b8154815290600101906020018083116200117c57829003601f168201915b5050509183525050601382015460209091019060ff166003811115620011c457620011c462005568565b6003811115620011d857620011d862005568565b90525092915050565b5f620011ef848484620032f5565b6001600160a01b0384165f9081526028602090815260408083203384529091528120546200121f90849062005c71565b90506200122e853383620030f7565b836001600160a01b0316856001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040516200127491815260200190565b60405180910390a3506001949350505050565b6013546001600160a01b03163314620012ef5760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f742074686520736c617368696e6720636f6e74726044820152621858dd60ea1b606482015260840162000e66565b5f610120820135602a826200130b604086016020870162005532565b6001600160a01b03166001600160a01b031681526020019081526020015f20600901546200133a919062005c71565b60c0830135602a5f62001354604087016020880162005532565b6001600160a01b03166001600160a01b031681526020019081526020015f206006015462001383919062005c71565b60a0840135602a5f6200139d604088016020890162005532565b6001600160a01b03166001600160a01b031681526020019081526020015f2060050154620013cc919062005c71565b620013d8919062005c87565b620013e4919062005c87565b6012546001600160a01b03165f908152602960205260408120805492935083929091906200141490849062005c87565b90915550829050602a5f62001430604084016020850162005532565b6001600160a01b0316815260208101919091526040015f2062001454828262005e82565b50505050565b601a546001600160a01b03163314620014875760405162461bcd60e51b815260040162000e669062005bb5565b62001493828262003423565b5050565b601a546001600160a01b03163314620014c45760405162461bcd60e51b815260040162000e669062005bb5565b601480546001600160a01b0319166001600160a01b03838116918217909255601554604051637adbf97360e01b8152600481019290925290911690637adbf973906024015f604051808303815f87803b15801562001520575f80fd5b505af115801562001533573d5f803e3d5ffd5b5050601754604051637adbf97360e01b81526001600160a01b0385811660048301529091169250637adbf97391506024015b5f604051808303815f87803b1580156200157d575f80fd5b505af115801562001590573d5f803e3d5ffd5b5050505050565b602d545f906060906001600160a01b03163314620015c95760405162461bcd60e51b815260040162000e669062005fd8565b60208054435f81815260219093526040832091909155601b54602254620015f1919062005c87565b6013546040516306c9789b60e41b8152929091146004830181905292506001600160a01b031690636c9789b0906024015f604051808303815f87803b15801562001639575f80fd5b505af11580156200164c573d5f803e3d5ffd5b505050508015620017fe57601954602b54602c546022546040516392eff3cd60e01b81526004810193909352602483019190915260448201524360648201525f916001600160a01b0316906392eff3cd90608401602060405180830381865afa158015620016bc573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620016e291906200601b565b905080602c541015620016f45750602c545b62001700308262003423565b80602c5f82825462001713919062005c71565b909155506200172590504782620034a7565b6200172f6200394a565b6200173962003a38565b5f62001744620027dc565b60145460405163422811f960e11b81529192506001600160a01b03169063845023f2906200177790849060040162005a10565b5f604051808303815f87803b1580156200178f575f80fd5b505af1158015620017a2573d5f803e3d5ffd5b5050505043602281905550600160205f828254620017c1919062005c87565b9091555050602080546040519081527febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335910160405180910390a150505b60145460408051634bb278f360e01b815290515f926001600160a01b031691634bb278f3916004808301926020929190829003018187875af115801562001847573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906200186d919062006033565b90508015620018e8576015546040805163a2e6204560e01b815290516001600160a01b039092169163a2e620459160048082019260209290919082900301815f875af1925050508015620018e0575060408051601f3d908101601f19168201909252620018dd9181019062006033565b60015b15620018e857505b600254602480546040805160208084028201810190925282815260ff909416939183915f9084015b82821015620019fb575f848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620019659062005c23565b80601f0160208091040260200160405190810160405280929190818152602001828054620019939062005c23565b8015620019e25780601f10620019b857610100808354040283529160200191620019e2565b820191905f5260205f20905b815481529060010190602001808311620019c457829003601f168201915b5050505050815250508152602001906001019062001910565b5050505090509350935050509091565b601a546001600160a01b0316331462001a385760405162461bcd60e51b815260040162000e669062005bb5565b601a80546001600160a01b0319166001600160a01b0383811691821790925560145460405163b3ab15fb60e01b815260048101929092529091169063b3ab15fb906024015f604051808303815f87803b15801562001a94575f80fd5b505af115801562001aa7573d5f803e3d5ffd5b505060155460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024015f604051808303815f87803b15801562001af0575f80fd5b505af115801562001b03573d5f803e3d5ffd5b505060165460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024015f604051808303815f87803b15801562001b4c575f80fd5b505af115801562001b5f573d5f803e3d5ffd5b505060175460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024015f604051808303815f87803b15801562001ba8575f80fd5b505af115801562001bbb573d5f803e3d5ffd5b505060185460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb915060240162001565565b5f80805b60245481101562001c4b576024818154811062001c175762001c1762006054565b905f5260205f209060030201600101548262001c34919062005c87565b91508062001c428162006068565b91505062001bf6565b50805f0362001c9d5760405162461bcd60e51b815260206004820152601c60248201527f54686520636f6d6d6974746565206973206e6f74207374616b696e6700000000604482015260640162000e66565b5f8362001cac60638762006083565b62001cb8919062005c87565b90505f8160405160200162001ccf91815260200190565b60408051601f19818403018152919052805160209091012090505f62001cf68483620060b1565b90505f805b60245481101562001da2576024818154811062001d1c5762001d1c62006054565b905f5260205f209060030201600101548262001d39919062005c87565b915062001d4860018362005c71565b831162001d8d576024818154811062001d655762001d6562006054565b5f9182526020909120600390910201546001600160a01b0316965062000e2195505050505050565b8062001d998162006068565b91505062001cfb565b5060405162461bcd60e51b815260206004820152602960248201527f5468657265206973206e6f2076616c696461746f72206c65667420696e20746860448201526865206e6574776f726b60b81b606482015260840162000e66565b601a546001600160a01b0316331462001e2b5760405162461bcd60e51b815260040162000e669062005bb5565b601b5481101562001ee2578060225462001e46919062005c87565b431062001ee25760405162461bcd60e51b815260206004820152605760248201527f63757272656e7420636861696e2068656164206578636565642074686520776960448201527f6e646f773a206c617374426c6f636b45706f6368202b205f6e6577506572696f60648201527f642c2074727920616761696e206c6174746572206f6e2e000000000000000000608482015260a40162000e66565b601b819055601354604051631ad7d11360e21b8152600481018390526001600160a01b0390911690636b5f444c906024015f604051808303815f87803b15801562001f2b575f80fd5b505af115801562001f3e573d5f803e3d5ffd5b505050507fd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f818160405162001f7491815260200190565b60405180910390a150565b601a546001600160a01b0316331462001fac5760405162461bcd60e51b815260040162000e669062005bb5565b600d55565b5f604051806102800160405280336001600160a01b031681526020015f6001600160a01b03168152602001856001600160a01b03168152602001868152602001600d5f016002015481526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f81526020014381526020015f81526020015f81526020015f81526020018481526020015f600381111562002070576200207062005568565b9052905062002080818362003b4f565b60208101516101808201516040517f8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c92620020c192339289918b91620060c7565b60405180910390a15050505050565b6001600160a01b038083165f818152602a6020526040902060010154909116146200210f5760405162461bcd60e51b815260040162000e669062005b32565b6001600160a01b038281165f908152602a60205260409020541633146200214a5760405162461bcd60e51b815260040162000e669062005b69565b6127108111156200219e5760405162461bcd60e51b815260206004820152601f60248201527f7265717569726520636f727265637420636f6d6d697373696f6e207261746500604482015260640162000e66565b604080516060810182526001600160a01b038481168252436020808401918252838501868152600c80545f908152600a909352958220855181546001600160a01b0319169516949094178455915160018085019190915591516002909301929092558354929390929091906200221690849062005c87565b90915550506040518281526001600160a01b038416907f4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf9060200160405180910390a2505050565b601a546001600160a01b031633146200228b5760405162461bcd60e51b815260040162000e669062005bb5565b6002805460ff19166001179055565b601a546001600160a01b03163314620022c75760405162461bcd60e51b815260040162000e669062005bb5565b5f8111620023185760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000e66565b601d55565b601a546001600160a01b031633146200234a5760405162461bcd60e51b815260040162000e669062005bb5565b6001600160a01b0382165f90815260296020526040902054811115620023ac5760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000e66565b6001600160a01b0382165f9081526029602052604081208054839290620023d590849062005c71565b9250508190555080602b5f828254620023ef919062005c71565b90915550506040518181526001600160a01b038316907f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3906020015b60405180910390a25050565b6001600160a01b038083165f818152602a602052604090206001015490911614620024765760405162461bcd60e51b815260040162000e669062005bec565b6001600160a01b0382165f908152602a602052604081206013015460ff166003811115620024a857620024a862005568565b14620024f75760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000e66565b6200149382823362003f9f565b6001600160a01b038083165f818152602a602052604090206001015490911614620025435760405162461bcd60e51b815260040162000e669062005bec565b5f81116200258c5760405162461bcd60e51b81526020600482015260156024820152740756e626f6e64696e6720616d6f756e74206973203605c1b604482015260640162000e66565b6200149382823362004186565b60606027805480602002602001604051908101604052809291908181526020015f905b828210156200266e578382905f5260205f20018054620025dc9062005c23565b80601f01602080910402602001604051908101604052809291908181526020018280546200260a9062005c23565b8015620026595780601f106200262f5761010080835404028352916020019162002659565b820191905f5260205f20905b8154815290600101906020018083116200263b57829003601f168201915b505050505081526020019060010190620025bc565b50505050905090565b5f62002685338484620032f5565b6040518281526001600160a01b0384169033907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a350600192915050565b60606024805480602002602001604051908101604052809291908181526020015f905b828210156200266e575f848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620027469062005c23565b80601f0160208091040260200160405190810160405280929190818152602001828054620027749062005c23565b8015620027c35780601f106200279957610100808354040283529160200191620027c3565b820191905f5260205f20905b815481529060010190602001808311620027a557829003601f168201915b50505050508152505081526020019060010190620026f1565b602d546060906001600160a01b031633146200280c5760405162461bcd60e51b815260040162000e669062005fd8565b601f546200285d5760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000e66565b6200286762005389565b601d546080820152601f8152602a602082015260246040820152602360608201526200289381620044ae565b620028a060275f620053a7565b60245480620028e75760405162461bcd60e51b8152602060048201526012602482015271636f6d6d697474656520697320656d70747960701b604482015260640162000e66565b5f8167ffffffffffffffff81111562002904576200290462005850565b6040519080825280602002602001820160405280156200292e578160200160208202803683370190505b5090505f5b8281101562002a14575f602a5f6024848154811062002956576200295662006054565b5f9182526020808320600392830201546001600160a01b031684528301939093526040909101812060278054600181018255925292507f98a476f1687bc3d60a2da2adbcba2c46958e61fa2fb4042cd7bc5816a710195b0190620029bd908301826200610f565b50600281015483516001600160a01b0390911690849084908110620029e657620029e662006054565b6001600160a01b0390921660209283029190910190910152508062002a0b8162006068565b91505062002933565b509250505090565b601a546001600160a01b0316331462002a495760405162461bcd60e51b815260040162000e669062005bb5565b62002a555f83620044cb565b62001493600182620044cb565b601a546001600160a01b0316331462002a8f5760405162461bcd60e51b815260040162000e669062005bb5565b601680546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038082165f818152602a60205260409020600101549091161462002af05760405162461bcd60e51b815260040162000e669062005b32565b6001600160a01b038082165f908152602a602052604090208054909116331462002b755760405162461bcd60e51b815260206004820152602f60248201527f726571756972652063616c6c657220746f2062652076616c696461746f72207460448201526e1c99585cdd5c9e481858d8dbdd5b9d608a1b606482015260840162000e66565b5f601382015460ff16600381111562002b925762002b9262005568565b0362002be15760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f7220616c7265616479206163746976650000000000000000604482015260640162000e66565b6002601382015460ff16600381111562002bff5762002bff62005568565b14801562002c105750438160100154115b1562002c5f5760405162461bcd60e51b815260206004820152601760248201527f76616c696461746f72207374696c6c20696e206a61696c000000000000000000604482015260640162000e66565b6003601382015460ff16600381111562002c7d5762002c7d62005568565b0362002ccc5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206a61696c6564207065726d616e656e746c7900000000604482015260640162000e66565b60138101805460ff191690558054601b546022546001600160a01b038581169316917f60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b59162002d1c919062005c87565b60405190815260200160405180910390a35050565b6060805f600181805462002d459062005c23565b80601f016020809104026020016040519081016040528092919081815260200182805462002d739062005c23565b801562002dc25780601f1062002d985761010080835404028352916020019162002dc2565b820191905f5260205f20905b81548152906001019060200180831162002da457829003601f168201915b5050505050915080805462002dd79062005c23565b80601f016020809104026020016040519081016040528092919081815260200182805462002e059062005c23565b801562002e545780601f1062002e2a5761010080835404028352916020019162002e54565b820191905f5260205f20905b81548152906001019060200180831162002e3657829003601f168201915b50505050509050915091509091565b6060601f80548060200260200160405190810160405280929190818152602001828054801562002ebb57602002820191905f5260205f20905b81546001600160a01b0316815260019091019060200180831162002e9c575b5050505050905090565b601a546001600160a01b0316331462002ef25760405162461bcd60e51b815260040162000e669062005bb5565b600e8190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200162001f74565b601a546001600160a01b0316331462002f555760405162461bcd60e51b815260040162000e669062005bb5565b601880546001600160a01b0319166001600160a01b0392909216919091179055565b601a546001600160a01b0316331462002fa45760405162461bcd60e51b815260040162000e669062005bb5565b62002fb05f80620053c4565b62002fbd60015f620053c4565b6002805460ff19169055565b601a546001600160a01b0316331462002ff65760405162461bcd60e51b815260040162000e669062005bb5565b601780546001600160a01b0319166001600160a01b0392909216919091179055565b601a546001600160a01b03163314620030455760405162461bcd60e51b815260040162000e669062005bb5565b601580546001600160a01b0319166001600160a01b0392909216919091179055565b602d546001600160a01b03163314620030945760405162461bcd60e51b815260040162000e669062005fd8565b6200309e6200394a565b62000eb5620027dc565b601a546001600160a01b03163314620030d55760405162461bcd60e51b815260040162000e669062005bb5565b601280546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b0383166200315b5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000e66565b6001600160a01b038216620031be5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000e66565b6001600160a01b038381165f8181526028602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b0381165f908152602a6020526040812090601382015460ff16600381111562003253576200325362005568565b14620032a25760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265206163746976650000000000000000604482015260640162000e66565b60138101805460ff191660011790558054601b546022546001600160a01b038581169316917f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c9162002d1c919062005c87565b6001600160a01b0383165f90815260296020526040902054811115620033575760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000e66565b6001600160a01b0383165f90815260296020526040812080548392906200338090849062005c71565b9250508190555080602b5f8282546200339a919062005c87565b90915550620033ad905081600262006083565b6001600160a01b0383165f9081526029602052604081208054909190620033d690849062005c87565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200321291815260200190565b6001600160a01b0382165f90815260296020526040812080548392906200344c90849062005c87565b9250508190555080602b5f82825462003466919062005c87565b90915550506040518181526001600160a01b038316907f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf906020016200242b565b81158015620034b4575080155b15620034be575050565b600d545f90670de0b6b3a764000090620034da90859062006083565b620034e69190620061f3565b9050801562003561576012546040515f916001600160a01b03169083908381818185875af1925050503d805f81146200353b576040519150601f19603f3d011682016040523d82523d5f602084013e62003540565b606091505b50909150508015156001036200355f576200355c828562005c71565b93505b505b8260255f82825462003574919062005c87565b909155505f90505b60245481101562001454575f602a5f60248481548110620035a157620035a162006054565b5f91825260208083206003909202909101546001600160a01b03168352820192909252604001812060235460248054929450909188919086908110620035eb57620035eb62006054565b905f5260205f2090600302016001015462003607919062006083565b620036139190620061f3565b90505f602354866024868154811062003630576200363062006054565b905f5260205f209060030201600101546200364c919062006083565b620036589190620061f3565b90505f8211806200366857505f81115b1562003931576002601384015460ff1660038111156200368c576200368c62005568565b1480620036b457506003601384015460ff166003811115620036b257620036b262005568565b145b156200375857601354602480546001600160a01b0390921691631de9d9b691859188908110620036e857620036e862006054565b5f91825260209091206003909102015460405160e084901b6001600160e01b03191681526001600160a01b0390911660048201526024015f604051808303818588803b15801562003737575f80fd5b505af11580156200374a573d5f803e3d5ffd5b505050505050505062003935565b5f836005015483856008015462003770919062006083565b6200377c9190620061f3565b90508015620037de5783546040516001600160a01b03909116906108fc9083905f818181858888f193505050503d805f8114620037d5576040519150601f19603f3d011682016040523d82523d5f602084013e620037da565b606091505b5050505b5f8460050154838660080154620037f6919062006083565b620038029190620061f3565b9050801562003824578454620038249030906001600160a01b031683620032f5565b5f62003831828562005c71565b90505f62003840848762005c71565b90505f8111806200385057505f82115b15620038e857600c870154620038729030906001600160a01b031684620032f5565b600c87015460405163a0ce552d60e01b8152600481018490526001600160a01b039091169063a0ce552d908390602401604080518083038185885af1158015620038be573d5f803e3d5ffd5b50505050506040513d601f19601f82011682018060405250810190620038e5919062006209565b50505b60018701546040518781526001600160a01b03909116907fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe5639060200160405180910390a2505050505b5050505b80620039418162006068565b9150506200357c565b6004545b60055481101562003977576200397181620039698162006068565b92506200461a565b6200394e565b50600554600455600854600754036200398c57565b6009545b600854811015620039b957620039b381620039ab8162006068565b925062004835565b62003990565b50600854600955600754805b60085481101562003a32576010545f828152600660205260409020600401544391620039f19162005c87565b1162003a175762003a028162004b1a565b62003a0f60018362005c87565b915062003a1d565b62003a32565b8062003a298162006068565b915050620039c5565b50600755565b600c54600b54101562003b4d57600b545f908152600a602052604090206010546001820154439162003a6a9162005c87565b111562003a745750565b600281015481546001600160a01b039081165f908152602a6020526040808220600490810185905585548416835291819020600c015490516319fac8fd60e01b81529216926319fac8fd9262003ace920190815260200190565b5f604051808303815f87803b15801562003ae6575f80fd5b505af115801562003af9573d5f803e3d5ffd5b5050600b80545f908152600a6020526040812080546001600160a01b03191681556001808201839055600290910182905582549094509192509062003b4090849062005c87565b9091555062003a38915050565b565b60e281511462003b995760405162461bcd60e51b8152602060048201526014602482015273092dcecc2d8d2c840e0e4dedecc40d8cadccee8d60631b604482015260640162000e66565b6030826102400151511462003bf15760405162461bcd60e51b815260206004820152601c60248201527f496e76616c696420636f6e73656e737573206b6579206c656e67746800000000604482015260640162000e66565b62003bfc8262004c5e565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a0000000000006020808301919091528451925191925f9262003c5e920160609190911b6bffffffffffffffffffffffff1916815260140190565b60405160208183030381529060405290505f8262003c7d835162004d89565b8360405160200162003c92939291906200622c565b60408051601f19818403018152828252805160209182012060028085526060850184529094505f93929091830190803683370190505090505f80808062003ce8898262003ce26041600262006083565b62004ea6565b90505f62003d068a62003cfe6041600262006083565b606062004ea6565b905060205b825181101562003dd45762003d21838262004fbe565b604080515f8152602081018083528d905260ff8316918101919091526060810184905260808101839052929850909650945060019060a0016020604051602081039080840390855afa15801562003d7a573d5f803e3d5ffd5b5050604051601f19015190508762003d94604184620061f3565b8151811062003da75762003da762006054565b6001600160a01b039092166020928302919091019091015262003dcc60418262005c87565b905062003d0b565b508a602001516001600160a01b0316865f8151811062003df85762003df862006054565b60200260200101516001600160a01b03161462003e6a5760405162461bcd60e51b815260206004820152602960248201527f496e76616c6964206e6f6465206b6579206f776e6572736869702070726f6f66604482015268081c1c9bdd9a59195960ba1b606482015260840162000e66565b8a604001516001600160a01b03168660018151811062003e8e5762003e8e62006054565b60200260200101516001600160a01b03161462003f025760405162461bcd60e51b815260206004820152602b60248201527f496e76616c6964206f7261636c65206b6579206f776e6572736869702070726f60448201526a1bd9881c1c9bdd9a59195960aa1b606482015260840162000e66565b600162003f198c6102400151838e5f015162004ff4565b1462003f875760405162461bcd60e51b815260206004820152603660248201527f496e76616c696420636f6e73656e737573206b6579206f776e65727368697020604482015275383937b7b3103337b9103932b3b4b9ba3930ba34b7b760511b606482015260840162000e66565b62003f928b6200505e565b5050505050505050505050565b5f821162003ffc5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000e66565b6001600160a01b0381165f90815260296020526040902054821115620040655760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000e66565b6001600160a01b0381165f90815260296020526040812080548492906200408e90849062005c71565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580545f908152600394859052978820875181549088166001600160a01b0319918216178255955160018201805491909816961695909517909555905160028401555191015580549192620041198362006068565b90915550506001600160a01b038481165f818152602a6020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d91015b60405180910390a35050505050565b6001600160a01b038084165f908152602a602052604090208054909183811691161480620042f357600c820154604051631092ab9160e31b81526001600160a01b0385811660048301525f9216906384955c8890602401602060405180830381865afa158015620041f9573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906200421f91906200601b565b905084811015620042875760405162461bcd60e51b815260206004820152602b60248201527f696e73756666696369656e7420756e6c6f636b6564204c6971756964204e657760448201526a746f6e2062616c616e636560a81b606482015260840162000e66565b600c83015460405163282d3fdf60e01b81526001600160a01b038681166004830152602482018890529091169063282d3fdf906044015f604051808303815f87803b158015620042d5575f80fd5b505af1158015620042e8573d5f803e3d5ffd5b505050505062004385565b8382600b015483600801546200430a919062005c71565b10156200436a5760405162461bcd60e51b815260206004820152602760248201527f696e73756666696369656e742073656c6620626f6e646564206e6577746f6e2060448201526662616c616e636560c81b606482015260840162000e66565b8382600b015f8282546200437f919062005c87565b90915550505b6040805160e0810182526001600160a01b03808616825287811660208084019182528385018981525f60608601818152436080880190815260a088018381528a151560c08a019081526008805486526006909752998420985189549089166001600160a01b0319918216178a55965160018a01805491909916971696909617909655915160028701559051600386015592516004850155905160059093018054945115156101000261ff00199415159490941661ffff199095169490941792909217909255805491620044588362006068565b9190505550826001600160a01b0316856001600160a01b03167f63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc8387604051620041779291909115158252602082015260400190565b60fa60a05f808285855af4620044c6573d5f803e3d5ffd5b505050565b8154600260018083161561010002038216048251808201602081106020841001600281146200457857600181146200459e57865f526020840460205f2001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b8082101562004561578151845560018401935060208201915062004542565b815191036101000a90819004029091555062004611565b60028302826020036101000a846020036101000a60208901510402018501875562004611565b865f526020840460205f2001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b80821015620045ff5781518455600184019350602082019150620045e0565b815191036101000a9081900402909155505b50505050505050565b5f81815260036020908152604080832060018101546001600160a01b03168452602a90925282209091601382015460ff1660038111156200465f576200465f62005568565b14620046fe57600282015482546001600160a01b03165f90815260296020526040812080549091906200469490849062005c87565b909155505081546001830154600284015460138401546040517f1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f87834294620046f1946001600160a01b0391821694911692909160ff9091169062006274565b60405180910390a1505050565b805482546001600160a01b03908116911614620047f3575f80826008015483600501546200472d919062005c71565b9050805f036200474457836002015491506200476a565b80846002015484600d01546200475b919062006083565b620047679190620061f3565b91505b600c83015484546040516340c10f1960e01b81526001600160a01b039182166004820152602481018590529116906340c10f19906044015f604051808303815f87803b158015620047b9575f80fd5b505af1158015620047cc573d5f803e3d5ffd5b505050508183600d015f828254620047e5919062005c87565b909155506200481292505050565b8160020154816008015f8282546200480c919062005c87565b90915550505b8160020154816005015f8282546200482b919062005c87565b9091555050505050565b5f81815260066020908152604080832060018101546001600160a01b03168452602a9092528220600582015491929091610100900460ff1662004a19576002830154600c8301548454604051637eee288d60e01b81526001600160a01b03918216600482015260248101849052911690637eee288d906044015f604051808303815f87803b158015620048c6575f80fd5b505af1158015620048d9573d5f803e3d5ffd5b50505050600c8301548454604051632770a7eb60e21b81526001600160a01b03918216600482015260248101849052911690639dc29fac906044015f604051808303815f87803b1580156200492c575f80fd5b505af11580156200493f573d5f803e3d5ffd5b505050505f836008015484600501546200495a919062005c71565b600d8501549091506200496e828462006083565b6200497a9190620061f3565b92508184600d015f82825462004991919062005c71565b909155505060068401545f03620049af5760038501839055620049d7565b60068401546007850154620049c5908562006083565b620049d19190620061f3565b60038601555b82846006015f828254620049ec919062005c87565b909155505060038501546007850180545f9062004a0b90849062005c87565b9091555062004ae992505050565b506002820154600882015481111562004a33575060088101545b81600901545f0362004a4c576003830181905562004a74565b6009820154600a83015462004a62908362006083565b62004a6e9190620061f3565b60038401555b80826009015f82825462004a89919062005c87565b90915550506003830154600a830180545f9062004aa890849062005c87565b9250508190555080826008015f82825462004ac4919062005c71565b90915550506002830154600b830180545f9062004ae390849062005c71565b90915550505b6005808401805460ff19166001179055820180548291905f9062004b0f90849062005c71565b909155505050505050565b5f818152600660205260408120600381015490910362004b38575050565b60018101546001600160a01b03165f908152602a602052604081206005830154909190610100900460ff1662004bd15781600701548260060154846003015462004b83919062006083565b62004b8f9190620061f3565b905080826006015f82825462004ba6919062005c71565b909155505060038301546007830180545f9062004bc590849062005c71565b9091555062004c349050565b81600a01548260090154846003015462004bec919062006083565b62004bf89190620061f3565b905080826009015f82825462004c0f919062005c71565b90915550506003830154600a830180545f9062004c2e90849062005c71565b90915550505b82546001600160a01b03165f908152602960205260408120805483929062004b0f90849062005c87565b5f62004c6e826060015162005289565b6001600160a01b0390911660208401529050801562004cbe5760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000e66565b6020808301516001600160a01b039081165f908152602a909252604090912060010154161562004d315760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000e66565b61271082608001511115620014935760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000e66565b6060815f0362004db05750506040805180820190915260018152600360fc1b602082015290565b815f5b811562004ddf578062004dc68162006068565b915062004dd79050600a83620061f3565b915062004db3565b5f8167ffffffffffffffff81111562004dfc5762004dfc62005850565b6040519080825280601f01601f19166020018201604052801562004e27576020820181803683370190505b5090505b841562004e9e5762004e3f60018362005c71565b915062004e4e600a86620060b1565b62004e5b90603062005c87565b60f81b81838151811062004e735762004e7362006054565b60200101906001600160f81b03191690815f1a90535062004e96600a86620061f3565b945062004e2b565b949350505050565b60608162004eb681601f62005c87565b101562004ef75760405162461bcd60e51b815260206004820152600e60248201526d736c6963655f6f766572666c6f7760901b604482015260640162000e66565b62004f03828462005c87565b8451101562004f495760405162461bcd60e51b8152602060048201526011602482015270736c6963655f6f75744f66426f756e647360781b604482015260640162000e66565b60608215801562004f695760405191505f82526020820160405262004fb5565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101562004fa457805183526020928301920162004f8a565b5050858452601f01601f1916604052505b50949350505050565b818101805160208201516040909201519091905f1a601b81101562004fed5762004fea601b82620062a2565b90505b9250925092565b5f62004fff62005400565b5f8585856040516020016200501793929190620062be565b60405160208183030381529060405290505f60fb90505f825160206200503e919062005c87565b90506020848285855afa62005051575f80fd5b5050905195945050505050565b6101808101516001600160a01b0316620050dc57601f545f90620050829062004d89565b90508160200151825f0151836080015183604051620050a1906200541e565b620050b094939291906200630c565b604051809103905ff080158015620050ca573d5f803e3d5ffd5b506001600160a01b0316610180830152505b60208181018051601f8054600180820183555f9283527fa03837a25210ee280c2113ff4b77ca23440b19d4866cca721c801278fd08d80790910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602a909552604091829020865181548716908516178155935190840180548616918416919091179055840151600283018054909416911617909155606082015182919060038201906200518d90826200634a565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f8201556102008201516010820155610220820151601182015561024082015160128201906200525790826200634a565b5061026082015160138201805460ff191660018360038111156200527f576200527f62005568565b0217905550505050565b5f80620052956200542c565b5f8060ff9050604083875160208901845afa620052b0575f80fd5b50508051602090910151600160601b90910494909350915050565b6040518061028001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b03168152602001606081526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f81526020015f81526020015f81526020015f81526020015f8152602001606081526020015f600381111562005384576200538462005568565b905290565b6040518060a001604052806005906020820280368337509192915050565b5080545f8255905f5260205f209081019062000eb591906200544a565b508054620053d29062005c23565b5f825580601f10620053e2575050565b601f0160209004905f5260205f209081019062000eb591906200546e565b60405180602001604052806001906020820280368337509192915050565b6117d6806200640883390190565b60405180604001604052806002906020820280368337509192915050565b808211156200546a575f620054608282620053c4565b506001016200544a565b5090565b5b808211156200546a575f81556001016200546f565b5f5b83811015620054a057818101518382015260200162005486565b50505f910152565b5f8151808452620054c181602086016020860162005484565b601f01601f19169290920160200192915050565b602081525f620054e96020830184620054a8565b9392505050565b6001600160a01b038116811462000eb5575f80fd5b5f806040838503121562005517575f80fd5b82356200552481620054f0565b946020939093013593505050565b5f6020828403121562005543575f80fd5b8135620054e981620054f0565b5f6020828403121562005561575f80fd5b5035919050565b634e487b7160e01b5f52602160045260245ffd5b600481106200559957634e487b7160e01b5f52602160045260245ffd5b9052565b60208152620055b86020820183516001600160a01b03169052565b5f6020830151620055d460408401826001600160a01b03169052565b5060408301516001600160a01b0381166060840152506060830151610280806080850152620056086102a0850183620054a8565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a062005690818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f190161026080880191909152909150620056f48483620054a8565b9350808701519150506200570b828601826200557c565b5090949350505050565b5f805f6060848603121562005728575f80fd5b83356200573581620054f0565b925060208401356200574781620054f0565b929592945050506040919091013590565b5f6020828403121562005769575f80fd5b813567ffffffffffffffff81111562005780575f80fd5b82016102808185031215620054e9575f80fd5b5f81518084526020808501808196508360051b810191508286015f5b8581101562005806578284038952815180516001600160a01b031685528581015186860152604090810151606091860182905290620057f181870183620054a8565b9a87019a9550505090840190600101620057af565b5091979650505050505050565b8215158152604060208201525f62004e9e604083018462005793565b5f806040838503121562005841575f80fd5b50508035926020909101359150565b634e487b7160e01b5f52604160045260245ffd5b5f82601f83011262005874575f80fd5b813567ffffffffffffffff8082111562005892576200589262005850565b604051601f8301601f19908116603f01168101908282118183101715620058bd57620058bd62005850565b81604052838152866020858801011115620058d6575f80fd5b836020870160208301375f602085830101528094505050505092915050565b5f805f806080858703121562005909575f80fd5b843567ffffffffffffffff8082111562005921575f80fd5b6200592f8883890162005864565b9550602087013591506200594382620054f0565b9093506040860135908082111562005959575f80fd5b620059678883890162005864565b935060608701359150808211156200597d575f80fd5b506200598c8782880162005864565b91505092959194509250565b5f602080830181845280855180835260408601915060408160051b87010192508387015f5b82811015620059ef57603f19888603018452620059dc858351620054a8565b94509285019290850190600101620059bd565b5092979650505050505050565b602081525f620054e9602083018462005793565b602080825282518282018190525f9190848201906040850190845b8181101562005a525783516001600160a01b03168352928401929184019160010162005a2b565b50909695505050505050565b5f806040838503121562005a70575f80fd5b823567ffffffffffffffff8082111562005a88575f80fd5b62005a968683870162005864565b9350602085013591508082111562005aac575f80fd5b5062005abb8582860162005864565b9150509250929050565b604081525f62005ad96040830185620054a8565b828103602084015262005aed8185620054a8565b95945050505050565b5f806040838503121562005b08575f80fd5b823562005b1581620054f0565b9150602083013562005b2781620054f0565b809150509250929050565b6020808252601c908201527f76616c696461746f72206d757374206265207265676973746572656400000000604082015260600190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b60208082526018908201527f76616c696461746f72206e6f7420726567697374657265640000000000000000604082015260600190565b600181811c9082168062005c3857607f821691505b60208210810362005c5757634e487b7160e01b5f52602260045260245ffd5b50919050565b634e487b7160e01b5f52601160045260245ffd5b8181038181111562000e215762000e2162005c5d565b8082018082111562000e215762000e2162005c5d565b5f813562000e2181620054f0565b80546001600160a01b0319166001600160a01b0392909216919091179055565b5f808335601e1984360301811262005ce1575f80fd5b83018035915067ffffffffffffffff82111562005cfc575f80fd5b60200191503681900382131562005d11575f80fd5b9250929050565b601f821115620044c6575f81815260208120601f850160051c8101602086101562005d405750805b601f850160051c820191505b8181101562005d615782815560010162005d4c565b505050505050565b5f19600383901b1c191660019190911b1790565b67ffffffffffffffff83111562005d985762005d9862005850565b62005db08362005da9835462005c23565b8362005d18565b5f601f84116001811462005de1575f851562005dcc5750838201355b62005dd8868262005d69565b84555062001590565b5f83815260209020601f19861690835b8281101562005e13578685013582556020948501946001909201910162005df1565b508682101562005e30575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b5f81356004811062000e21575f80fd5b6004821062005e6f57634e487b7160e01b5f52602160045260245ffd5b60ff1981541660ff831681178255505050565b62005e9862005e918362005c9d565b8262005cab565b62005eb462005eaa6020840162005c9d565b6001830162005cab565b62005ed062005ec66040840162005c9d565b6002830162005cab565b62005edf606083018362005ccb565b62005eef81836003860162005d7d565b50506080820135600482015560a0820135600582015560c0820135600682015560e0820135600782015561010082013560088201556101208201356009820155610140820135600a820155610160820135600b82015562005f6262005f58610180840162005c9d565b600c830162005cab565b6101a0820135600d8201556101c0820135600e8201556101e0820135600f8201556102008201356010820155610220820135601182015562005fa961024083018362005ccb565b62005fb981836012860162005d7d565b50506200149362005fce610260840162005e42565b6013830162005e52565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b5f602082840312156200602c575f80fd5b5051919050565b5f6020828403121562006044575f80fd5b81518015158114620054e9575f80fd5b634e487b7160e01b5f52603260045260245ffd5b5f600182016200607c576200607c62005c5d565b5060010190565b808202811582820484141762000e215762000e2162005c5d565b634e487b7160e01b5f52601260045260245ffd5b5f82620060c257620060c26200609d565b500690565b5f60018060a01b0380881683528087166020840152808616604084015260a06060840152620060fa60a0840186620054a8565b91508084166080840152509695505050505050565b8181036200611b575050565b62006127825462005c23565b67ffffffffffffffff81111562006142576200614262005850565b6200615a8162006153845462005c23565b8462005d18565b5f601f8211600181146200618b575f8315620061765750848201545b62006182848262005d69565b85555062001590565b5f85815260209020601f198416905f86815260209020845b83811015620061c55782860154825560019586019590910190602001620061a3565b5085831015620061e357818501545f19600388901b60f8161c191681555b5050505050600190811b01905550565b5f826200620457620062046200609d565b500490565b5f80604083850312156200621b575f80fd5b505080516020909101519092909150565b5f84516200623f81846020890162005484565b8451908301906200625581836020890162005484565b84519101906200626a81836020880162005484565b0195945050505050565b6001600160a01b03858116825284166020820152604081018390526080810162005aed60608301846200557c565b60ff818116838216019081111562000e215762000e2162005c5d565b5f8451620062d181846020890162005484565b845190830190620062e781836020890162005484565b60609490941b6bffffffffffffffffffffffff19169301928352505060140192915050565b6001600160a01b03858116825284166020820152604081018390526080606082018190525f906200634090830184620054a8565b9695505050505050565b815167ffffffffffffffff81111562006367576200636762005850565b620063788162006153845462005c23565b602080601f831160018114620063ab575f8415620063965750858301515b620063a2858262005d69565b86555062005d61565b5f85815260208120601f198616915b82811015620063db57888601518255948401946001909101908401620063ba565b5085821015620061e35793909601515f1960f8600387901b161c19169092555050600190811b0190555056fe608060405234801562000010575f80fd5b50604051620017d6380380620017d6833981016040819052620000339162000149565b61271082111562000042575f80fd5b600d80546001600160a01b038087166001600160a01b031992831617909255600e805492861692909116919091179055600f8290556040516200008a90829060200162000230565b604051602081830303815290604052600b9081620000a99190620002ea565b5080604051602001620000bd919062000230565b604051602081830303815290604052600c9081620000dc9190620002ea565b50505f80546001600160a01b0319163317905550620003b2915050565b6001600160a01b03811681146200010e575f80fd5b50565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156200014157818101518382015260200162000127565b50505f910152565b5f805f80608085870312156200015d575f80fd5b84516200016a81620000f9565b60208601519094506200017d81620000f9565b6040860151606087015191945092506001600160401b0380821115620001a1575f80fd5b818701915087601f830112620001b5575f80fd5b815181811115620001ca57620001ca62000111565b604051601f8201601f19908116603f01168101908382118183101715620001f557620001f562000111565b816040528281528a60208487010111156200020e575f80fd5b6200022183602083016020880162000125565b979a9699509497505050505050565b644c4e544e2d60d81b81525f82516200025181600585016020870162000125565b9190910160050192915050565b600181811c908216806200027357607f821691505b6020821081036200029257634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620002e5575f81815260208120601f850160051c81016020861015620002c05750805b601f850160051c820191505b81811015620002e157828155600101620002cc565b5050505b505050565b81516001600160401b0381111562000306576200030662000111565b6200031e816200031784546200025e565b8462000298565b602080601f83116001811462000354575f84156200033c5750858301515b5f19600386901b1c1916600185901b178555620002e1565b5f85815260208120601f198616915b82811015620003845788860151825594840194600190910190840162000363565b5085821015620003a257878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61141680620003c05f395ff3fe608060405260043610610147575f3560e01c806359355736116100b3578063949813b81161006d578063949813b8146103ac57806395d89b41146103e05780639dc29fac146103f4578063a0ce552d14610413578063a9059cbb14610426578063dd62ed3e14610445575f80fd5b806359355736146102d25780635ea1d6f81461030657806361d027b31461031b57806370a082311461033a5780637eee288d1461036e57806384955c881461038d575f80fd5b8063282d3fdf11610104578063282d3fdf146102195780632f2c3f2e14610238578063313ce5671461024d578063372500ab146102685780633a5381b51461027c57806340c10f19146102b3575f80fd5b806306fdde031461014b578063095ea7b31461017557806318160ddd146101a4578063187cf4d7146101c257806319fac8fd146101d957806323b872dd146101fa575b5f80fd5b348015610156575f80fd5b5061015f610489565b60405161016c9190611182565b60405180910390f35b348015610180575f80fd5b5061019461018f3660046111e8565b610519565b604051901515815260200161016c565b3480156101af575f80fd5b506004545b60405190815260200161016c565b3480156101cd575f80fd5b506101b4633b9aca0081565b3480156101e4575f80fd5b506101f86101f3366004611210565b61052f565b005b348015610205575f80fd5b50610194610214366004611227565b610566565b348015610224575f80fd5b506101f86102333660046111e8565b610657565b348015610243575f80fd5b506101b461271081565b348015610258575f80fd5b506040516012815260200161016c565b348015610273575f80fd5b506101f8610739565b348015610287575f80fd5b50600d5461029b906001600160a01b031681565b6040516001600160a01b03909116815260200161016c565b3480156102be575f80fd5b506101f86102cd3660046111e8565b6108a9565b3480156102dd575f80fd5b506101b46102ec366004611260565b6001600160a01b03165f9081526002602052604090205490565b348015610311575f80fd5b506101b4600f5481565b348015610326575f80fd5b50600e5461029b906001600160a01b031681565b348015610345575f80fd5b506101b4610354366004611260565b6001600160a01b03165f9081526001602052604090205490565b348015610379575f80fd5b506101f86103883660046111e8565b61090e565b348015610398575f80fd5b506101b46103a7366004611260565b6109d1565b3480156103b7575f80fd5b506103cb6103c6366004611260565b6109fe565b6040805192835260208301919091520161016c565b3480156103eb575f80fd5b5061015f610a62565b3480156103ff575f80fd5b506101f861040e3660046111e8565b610a71565b6103cb610421366004611210565b610ace565b348015610431575f80fd5b506101946104403660046111e8565b610d78565b348015610450575f80fd5b506101b461045f366004611280565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b6060600b8054610498906112b1565b80601f01602080910402602001604051908101604052809291908181526020018280546104c4906112b1565b801561050f5780601f106104e65761010080835404028352916020019161050f565b820191905f5260205f20905b8154815290600101906020018083116104f257829003601f168201915b5050505050905090565b5f610525338484610dc3565b5060015b92915050565b5f546001600160a01b031633146105615760405162461bcd60e51b8152600401610558906112e9565b60405180910390fd5b600f55565b6001600160a01b0383165f908152600360209081526040808320338452909152812054828110156105ea5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b6064820152608401610558565b6105fe85336105f98685611345565b610dc3565b6106088584610ee6565b6106128484610fd8565b836001600160a01b0316856001600160a01b03165f805160206113c18339815191528560405161064491815260200190565b60405180910390a3506001949350505050565b5f546001600160a01b031633146106805760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f9081526002602090815260408083205460019092529091205482916106af91611345565b10156107095760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611358565b90915550505050565b5f8061074433611022565b335f818152600560209081526040808320839055600890915280822082905551939550919350909184908381818185875af1925050503d805f81146107a4576040519150601f19603f3d011682016040523d82523d5f602084013e6107a9565b606091505b50509050806107ef5760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b7321020aa2760711b6044820152606401610558565b5f5460405163a9059cbb60e01b8152336004820152602481018490526001600160a01b039091169063a9059cbb906044016020604051808303815f875af115801561083c573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610860919061136b565b9050806108a45760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b73210272a2760711b6044820152606401610558565b505050565b5f546001600160a01b031633146108d25760405162461bcd60e51b8152600401610558906112e9565b6108dc8282610fd8565b6040518181526001600160a01b038316905f905f805160206113c1833981519152906020015b60405180910390a35050565b5f546001600160a01b031633146109375760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f908152600260205260409020548111156109aa5760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611345565b6001600160a01b0381165f9081526002602090815260408083205460019092528220546105299190611345565b5f805f80610a0b856110c7565b6001600160a01b0387165f908152600560205260409020549193509150610a33908390611358565b6001600160a01b0386165f90815260086020526040902054909450610a59908290611358565b92505050915091565b6060600c8054610498906112b1565b5f546001600160a01b03163314610a9a5760405162461bcd60e51b8152600401610558906112e9565b610aa48282610ee6565b6040518181525f906001600160a01b038416905f805160206113c183398151915290602001610902565b5f805481906001600160a01b03163314610afa5760405162461bcd60e51b8152600401610558906112e9565b600f5434905f9061271090610b0f908461138a565b610b1991906113a1565b905081811115610b665760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610b708183611345565b600e546040519193506001600160a01b0316906108fc9083905f818181858888f193505050503d805f8114610bc0576040519150601f19603f3d011682016040523d82523d5f602084013e610bc5565b606091505b5050505f612710600f5487610bda919061138a565b610be491906113a1565b905085811115610c315760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610c3b8187611345565b5f54600e5460405163a9059cbb60e01b81526001600160a01b03918216600482015260248101859052929850169063a9059cbb906044016020604051808303815f875af1158015610c8e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610cb2919061136b565b506004545f90610cc6633b9aca008661138a565b610cd091906113a1565b905080600754610ce09190611358565b6007556004545f90610cf6633b9aca008a61138a565b610d0091906113a1565b905080600a54610d109190611358565b600a556004545f90633b9aca0090610d28908561138a565b610d3291906113a1565b90505f633b9aca0060045484610d48919061138a565b610d5291906113a1565b9050610d5e8287611358565b610d688287611358565b9850985050505050505050915091565b5f610d833383610ee6565b610d8d8383610fd8565b6040518281526001600160a01b0384169033905f805160206113c18339815191529060200160405180910390a350600192915050565b6001600160a01b038316610e255760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610558565b6001600160a01b038216610e865760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610558565b6001600160a01b038381165f8181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610eef82611022565b50506001600160a01b0382165f90815260016020908152604080832054600290925290912054610f1f9082611345565b821115610f6e5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e647300000000006044820152606401610558565b610f788282611345565b6001600160a01b0384165f90815260016020526040902055808203610fbd576001600160a01b0383165f90815260066020908152604080832083905560099091528120555b8160045f828254610fce9190611345565b9091555050505050565b610fe182611022565b50506001600160a01b0382165f908152600160205260408120805483929061100a908490611358565b925050819055508060045f8282546107309190611358565b5f805f8061102f856110c7565b6001600160a01b0387165f908152600560205260409020549193509150611057908390611358565b6001600160a01b0386165f90815260056020908152604080832084905560075460068352818420556008909152902054909450611095908290611358565b6001600160a01b039095165f908152600860209081526040808320889055600a54600990925290912055509193915050565b6001600160a01b0381165f9081526001602052604081205481908082036110f357505f93849350915050565b6001600160a01b0384165f908152600660205260408120546007546111189190611345565b6001600160a01b0386165f90815260096020526040812054600a5492935090916111429190611345565b9050633b9aca00611153848461138a565b61115d91906113a1565b9450633b9aca0061116e848361138a565b61117891906113a1565b9350505050915091565b5f6020808352835180828501525f5b818110156111ad57858101830151858201604001528201611191565b505f604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146111e3575f80fd5b919050565b5f80604083850312156111f9575f80fd5b611202836111cd565b946020939093013593505050565b5f60208284031215611220575f80fd5b5035919050565b5f805f60608486031215611239575f80fd5b611242846111cd565b9250611250602085016111cd565b9150604084013590509250925092565b5f60208284031215611270575f80fd5b611279826111cd565b9392505050565b5f8060408385031215611291575f80fd5b61129a836111cd565b91506112a8602084016111cd565b90509250929050565b600181811c908216806112c557607f821691505b6020821081036112e357634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561052957610529611331565b8082018082111561052957610529611331565b5f6020828403121561137b575f80fd5b81518015158114611279575f80fd5b808202811582820484141761052957610529611331565b5f826113bb57634e487b7160e01b5f52601260045260245ffd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220fcf7fe14f5b055004178aa5c27d7983104c707dfee4e3b218ce03dfc1bc4bea864736f6c63430008150033a2646970667358221220e3f6236419f0934f760ab164479d0bdcc6cb5206f5543ce9d6eab55455ef7ff164736f6c63430008150033608060405234801562000010575f80fd5b50604051620017d6380380620017d6833981016040819052620000339162000149565b61271082111562000042575f80fd5b600d80546001600160a01b038087166001600160a01b031992831617909255600e805492861692909116919091179055600f8290556040516200008a90829060200162000230565b604051602081830303815290604052600b9081620000a99190620002ea565b5080604051602001620000bd919062000230565b604051602081830303815290604052600c9081620000dc9190620002ea565b50505f80546001600160a01b0319163317905550620003b2915050565b6001600160a01b03811681146200010e575f80fd5b50565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156200014157818101518382015260200162000127565b50505f910152565b5f805f80608085870312156200015d575f80fd5b84516200016a81620000f9565b60208601519094506200017d81620000f9565b6040860151606087015191945092506001600160401b0380821115620001a1575f80fd5b818701915087601f830112620001b5575f80fd5b815181811115620001ca57620001ca62000111565b604051601f8201601f19908116603f01168101908382118183101715620001f557620001f562000111565b816040528281528a60208487010111156200020e575f80fd5b6200022183602083016020880162000125565b979a9699509497505050505050565b644c4e544e2d60d81b81525f82516200025181600585016020870162000125565b9190910160050192915050565b600181811c908216806200027357607f821691505b6020821081036200029257634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620002e5575f81815260208120601f850160051c81016020861015620002c05750805b601f850160051c820191505b81811015620002e157828155600101620002cc565b5050505b505050565b81516001600160401b0381111562000306576200030662000111565b6200031e816200031784546200025e565b8462000298565b602080601f83116001811462000354575f84156200033c5750858301515b5f19600386901b1c1916600185901b178555620002e1565b5f85815260208120601f198616915b82811015620003845788860151825594840194600190910190840162000363565b5085821015620003a257878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61141680620003c05f395ff3fe608060405260043610610147575f3560e01c806359355736116100b3578063949813b81161006d578063949813b8146103ac57806395d89b41146103e05780639dc29fac146103f4578063a0ce552d14610413578063a9059cbb14610426578063dd62ed3e14610445575f80fd5b806359355736146102d25780635ea1d6f81461030657806361d027b31461031b57806370a082311461033a5780637eee288d1461036e57806384955c881461038d575f80fd5b8063282d3fdf11610104578063282d3fdf146102195780632f2c3f2e14610238578063313ce5671461024d578063372500ab146102685780633a5381b51461027c57806340c10f19146102b3575f80fd5b806306fdde031461014b578063095ea7b31461017557806318160ddd146101a4578063187cf4d7146101c257806319fac8fd146101d957806323b872dd146101fa575b5f80fd5b348015610156575f80fd5b5061015f610489565b60405161016c9190611182565b60405180910390f35b348015610180575f80fd5b5061019461018f3660046111e8565b610519565b604051901515815260200161016c565b3480156101af575f80fd5b506004545b60405190815260200161016c565b3480156101cd575f80fd5b506101b4633b9aca0081565b3480156101e4575f80fd5b506101f86101f3366004611210565b61052f565b005b348015610205575f80fd5b50610194610214366004611227565b610566565b348015610224575f80fd5b506101f86102333660046111e8565b610657565b348015610243575f80fd5b506101b461271081565b348015610258575f80fd5b506040516012815260200161016c565b348015610273575f80fd5b506101f8610739565b348015610287575f80fd5b50600d5461029b906001600160a01b031681565b6040516001600160a01b03909116815260200161016c565b3480156102be575f80fd5b506101f86102cd3660046111e8565b6108a9565b3480156102dd575f80fd5b506101b46102ec366004611260565b6001600160a01b03165f9081526002602052604090205490565b348015610311575f80fd5b506101b4600f5481565b348015610326575f80fd5b50600e5461029b906001600160a01b031681565b348015610345575f80fd5b506101b4610354366004611260565b6001600160a01b03165f9081526001602052604090205490565b348015610379575f80fd5b506101f86103883660046111e8565b61090e565b348015610398575f80fd5b506101b46103a7366004611260565b6109d1565b3480156103b7575f80fd5b506103cb6103c6366004611260565b6109fe565b6040805192835260208301919091520161016c565b3480156103eb575f80fd5b5061015f610a62565b3480156103ff575f80fd5b506101f861040e3660046111e8565b610a71565b6103cb610421366004611210565b610ace565b348015610431575f80fd5b506101946104403660046111e8565b610d78565b348015610450575f80fd5b506101b461045f366004611280565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b6060600b8054610498906112b1565b80601f01602080910402602001604051908101604052809291908181526020018280546104c4906112b1565b801561050f5780601f106104e65761010080835404028352916020019161050f565b820191905f5260205f20905b8154815290600101906020018083116104f257829003601f168201915b5050505050905090565b5f610525338484610dc3565b5060015b92915050565b5f546001600160a01b031633146105615760405162461bcd60e51b8152600401610558906112e9565b60405180910390fd5b600f55565b6001600160a01b0383165f908152600360209081526040808320338452909152812054828110156105ea5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b6064820152608401610558565b6105fe85336105f98685611345565b610dc3565b6106088584610ee6565b6106128484610fd8565b836001600160a01b0316856001600160a01b03165f805160206113c18339815191528560405161064491815260200190565b60405180910390a3506001949350505050565b5f546001600160a01b031633146106805760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f9081526002602090815260408083205460019092529091205482916106af91611345565b10156107095760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611358565b90915550505050565b5f8061074433611022565b335f818152600560209081526040808320839055600890915280822082905551939550919350909184908381818185875af1925050503d805f81146107a4576040519150601f19603f3d011682016040523d82523d5f602084013e6107a9565b606091505b50509050806107ef5760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b7321020aa2760711b6044820152606401610558565b5f5460405163a9059cbb60e01b8152336004820152602481018490526001600160a01b039091169063a9059cbb906044016020604051808303815f875af115801561083c573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610860919061136b565b9050806108a45760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b73210272a2760711b6044820152606401610558565b505050565b5f546001600160a01b031633146108d25760405162461bcd60e51b8152600401610558906112e9565b6108dc8282610fd8565b6040518181526001600160a01b038316905f905f805160206113c1833981519152906020015b60405180910390a35050565b5f546001600160a01b031633146109375760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f908152600260205260409020548111156109aa5760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611345565b6001600160a01b0381165f9081526002602090815260408083205460019092528220546105299190611345565b5f805f80610a0b856110c7565b6001600160a01b0387165f908152600560205260409020549193509150610a33908390611358565b6001600160a01b0386165f90815260086020526040902054909450610a59908290611358565b92505050915091565b6060600c8054610498906112b1565b5f546001600160a01b03163314610a9a5760405162461bcd60e51b8152600401610558906112e9565b610aa48282610ee6565b6040518181525f906001600160a01b038416905f805160206113c183398151915290602001610902565b5f805481906001600160a01b03163314610afa5760405162461bcd60e51b8152600401610558906112e9565b600f5434905f9061271090610b0f908461138a565b610b1991906113a1565b905081811115610b665760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610b708183611345565b600e546040519193506001600160a01b0316906108fc9083905f818181858888f193505050503d805f8114610bc0576040519150601f19603f3d011682016040523d82523d5f602084013e610bc5565b606091505b5050505f612710600f5487610bda919061138a565b610be491906113a1565b905085811115610c315760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610c3b8187611345565b5f54600e5460405163a9059cbb60e01b81526001600160a01b03918216600482015260248101859052929850169063a9059cbb906044016020604051808303815f875af1158015610c8e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610cb2919061136b565b506004545f90610cc6633b9aca008661138a565b610cd091906113a1565b905080600754610ce09190611358565b6007556004545f90610cf6633b9aca008a61138a565b610d0091906113a1565b905080600a54610d109190611358565b600a556004545f90633b9aca0090610d28908561138a565b610d3291906113a1565b90505f633b9aca0060045484610d48919061138a565b610d5291906113a1565b9050610d5e8287611358565b610d688287611358565b9850985050505050505050915091565b5f610d833383610ee6565b610d8d8383610fd8565b6040518281526001600160a01b0384169033905f805160206113c18339815191529060200160405180910390a350600192915050565b6001600160a01b038316610e255760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610558565b6001600160a01b038216610e865760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610558565b6001600160a01b038381165f8181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610eef82611022565b50506001600160a01b0382165f90815260016020908152604080832054600290925290912054610f1f9082611345565b821115610f6e5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e647300000000006044820152606401610558565b610f788282611345565b6001600160a01b0384165f90815260016020526040902055808203610fbd576001600160a01b0383165f90815260066020908152604080832083905560099091528120555b8160045f828254610fce9190611345565b9091555050505050565b610fe182611022565b50506001600160a01b0382165f908152600160205260408120805483929061100a908490611358565b925050819055508060045f8282546107309190611358565b5f805f8061102f856110c7565b6001600160a01b0387165f908152600560205260409020549193509150611057908390611358565b6001600160a01b0386165f90815260056020908152604080832084905560075460068352818420556008909152902054909450611095908290611358565b6001600160a01b039095165f908152600860209081526040808320889055600a54600990925290912055509193915050565b6001600160a01b0381165f9081526001602052604081205481908082036110f357505f93849350915050565b6001600160a01b0384165f908152600660205260408120546007546111189190611345565b6001600160a01b0386165f90815260096020526040812054600a5492935090916111429190611345565b9050633b9aca00611153848461138a565b61115d91906113a1565b9450633b9aca0061116e848361138a565b61117891906113a1565b9350505050915091565b5f6020808352835180828501525f5b818110156111ad57858101830151858201604001528201611191565b505f604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146111e3575f80fd5b919050565b5f80604083850312156111f9575f80fd5b611202836111cd565b946020939093013593505050565b5f60208284031215611220575f80fd5b5035919050565b5f805f60608486031215611239575f80fd5b611242846111cd565b9250611250602085016111cd565b9150604084013590509250925092565b5f60208284031215611270575f80fd5b611279826111cd565b9392505050565b5f8060408385031215611291575f80fd5b61129a836111cd565b91506112a8602084016111cd565b90509250929050565b600181811c908216806112c557607f821691505b6020821081036112e357634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561052957610529611331565b8082018082111561052957610529611331565b5f6020828403121561137b575f80fd5b81518015158114611279575f80fd5b808202811582820484141761052957610529611331565b5f826113bb57634e487b7160e01b5f52601260045260245ffd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220fcf7fe14f5b055004178aa5c27d7983104c707dfee4e3b218ce03dfc1bc4bea864736f6c63430008150033",
}

// AutonityUpgradeTestABI is the input ABI used to generate the binding from.
// Deprecated: Use AutonityUpgradeTestMetaData.ABI instead.
var AutonityUpgradeTestABI = AutonityUpgradeTestMetaData.ABI

// Deprecated: Use AutonityUpgradeTestMetaData.Sigs instead.
// AutonityUpgradeTestFuncSigs maps the 4-byte function signature to its string representation.
var AutonityUpgradeTestFuncSigs = AutonityUpgradeTestMetaData.Sigs

// AutonityUpgradeTestBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AutonityUpgradeTestMetaData.Bin instead.
var AutonityUpgradeTestBin = AutonityUpgradeTestMetaData.Bin

// DeployAutonityUpgradeTest deploys a new Ethereum contract, binding an instance of AutonityUpgradeTest to it.
func (r *runner) deployAutonityUpgradeTest(opts *runOptions) (common.Address, uint64, *AutonityUpgradeTest, error) {
	parsed, err := AutonityUpgradeTestMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(AutonityUpgradeTestBin))
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &AutonityUpgradeTest{AutonityUpgradeTestCaller: AutonityUpgradeTestCaller{contract: c}, AutonityUpgradeTestTransactor: AutonityUpgradeTestTransactor{contract: c}, AutonityUpgradeTestFilterer: AutonityUpgradeTestFilterer{contract: c}}, nil
}

// AutonityUpgradeTest is an auto generated Go binding around an Ethereum contract.
type AutonityUpgradeTest struct {
	AutonityUpgradeTestCaller     // Read-only binding to the contract
	AutonityUpgradeTestTransactor // Write-only binding to the contract
	AutonityUpgradeTestFilterer   // Log filterer for contract events
}

// AutonityUpgradeTestCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type AutonityUpgradeTestCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// AutonityUpgradeTestTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AutonityUpgradeTestTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// AutonityUpgradeTestFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AutonityUpgradeTestFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) COMMISSIONRATEPRECISION(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "COMMISSION_RATE_PRECISION")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) Allowance(opts *runOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// AtnTotalRedistributed is a free data retrieval call binding the contract method 0xe98712b2.
//
// Solidity: function atnTotalRedistributed() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) AtnTotalRedistributed(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "atnTotalRedistributed")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) BalanceOf(opts *runOptions, _addr common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "balanceOf", _addr)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,uint256,address) policy, (address,address,address,address,address,address,address) contracts, (address,uint256,uint256,uint256) protocol, uint256 contractVersion)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) Config(opts *runOptions) (struct {
	Policy          AutonityPolicy
	Contracts       AutonityContracts
	Protocol        AutonityProtocol
	ContractVersion *big.Int
}, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "config")

	outstruct := new(struct {
		Policy          AutonityPolicy
		Contracts       AutonityContracts
		Protocol        AutonityProtocol
		ContractVersion *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.Policy = *abi.ConvertType(out[0], new(AutonityPolicy)).(*AutonityPolicy)
	outstruct.Contracts = *abi.ConvertType(out[1], new(AutonityContracts)).(*AutonityContracts)
	outstruct.Protocol = *abi.ConvertType(out[2], new(AutonityProtocol)).(*AutonityProtocol)
	outstruct.ContractVersion = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) Decimals(opts *runOptions) (uint8, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "decimals")

	if err != nil {
		return *new(uint8), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, consumed, err

}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() view returns(address)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) Deployer(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "deployer")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// EpochID is a free data retrieval call binding the contract method 0xc9d97af4.
//
// Solidity: function epochID() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) EpochID(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "epochID")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// EpochReward is a free data retrieval call binding the contract method 0x1604e416.
//
// Solidity: function epochReward() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) EpochReward(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "epochReward")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// EpochTotalBondedStake is a free data retrieval call binding the contract method 0x9c98e471.
//
// Solidity: function epochTotalBondedStake() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) EpochTotalBondedStake(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "epochTotalBondedStake")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetBlockPeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getBlockPeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256,bytes)[])
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetCommittee(opts *runOptions) ([]AutonityCommitteeMember, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getCommittee")

	if err != nil {
		return *new([]AutonityCommitteeMember), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]AutonityCommitteeMember)).(*[]AutonityCommitteeMember)
	return out0, consumed, err

}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetCommitteeEnodes(opts *runOptions) ([]string, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getCommitteeEnodes")

	if err != nil {
		return *new([]string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, consumed, err

}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetEpochFromBlock(opts *runOptions, _block *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getEpochFromBlock", _block)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetEpochPeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getEpochPeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetLastEpochBlock(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getLastEpochBlock")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetMaxCommitteeSize(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getMaxCommitteeSize")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetMinimumBaseFee(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getMinimumBaseFee")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetNewContract(opts *runOptions) ([]byte, string, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getNewContract")

	if err != nil {
		return *new([]byte), *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	return out0, out1, consumed, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetOperator(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getOperator")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetOracle(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getOracle")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetProposer is a free data retrieval call binding the contract method 0x5f7d3949.
//
// Solidity: function getProposer(uint256 height, uint256 round) view returns(address)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetProposer(opts *runOptions, height *big.Int, round *big.Int) (common.Address, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getProposer", height, round)

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetTreasuryAccount(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getTreasuryAccount")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetTreasuryFee(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getTreasuryFee")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetUnbondingPeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getUnbondingPeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8))
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetValidator(opts *runOptions, _addr common.Address) (AutonityValidator, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getValidator", _addr)

	if err != nil {
		return *new(AutonityValidator), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(AutonityValidator)).(*AutonityValidator)
	return out0, consumed, err

}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetValidators(opts *runOptions) ([]common.Address, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getValidators")

	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) GetVersion(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getVersion")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// InflationReserve is a free data retrieval call binding the contract method 0x9ac31074.
//
// Solidity: function inflationReserve() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) InflationReserve(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "inflationReserve")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) LastEpochBlock(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "lastEpochBlock")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) Name(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "name")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) Symbol(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "symbol")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTestCaller) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) ActivateValidator(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "activateValidator", _address)
	return consumed, err
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Approve(opts *runOptions, spender common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "approve", spender, amount)
	return consumed, err
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Bond(opts *runOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "bond", _validator, _amount)
	return consumed, err
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Burn(opts *runOptions, _addr common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "burn", _addr, _amount)
	return consumed, err
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) ChangeCommissionRate(opts *runOptions, _validator common.Address, _rate *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "changeCommissionRate", _validator, _rate)
	return consumed, err
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) CompleteContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "completeContractUpgrade")
	return consumed, err
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[])
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) ComputeCommittee(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "computeCommittee")
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256,bytes)[])
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Finalize(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "finalize")
	return consumed, err
}

// FinalizeInitialization is a paid mutator transaction binding the contract method 0xd861b0e8.
//
// Solidity: function finalizeInitialization() returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) FinalizeInitialization(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "finalizeInitialization")
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Mint(opts *runOptions, _addr common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "mint", _addr, _amount)
	return consumed, err
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) PauseValidator(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "pauseValidator", _address)
	return consumed, err
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x84467fdb.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) RegisterValidator(opts *runOptions, _enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "registerValidator", _enode, _oracleAddress, _consensusKey, _signatures)
	return consumed, err
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) ResetContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "resetContractUpgrade")
	return consumed, err
}

// SetAccountabilityContract is a paid mutator transaction binding the contract method 0x1250a28d.
//
// Solidity: function setAccountabilityContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetAccountabilityContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setAccountabilityContract", _address)
	return consumed, err
}

// SetAcuContract is a paid mutator transaction binding the contract method 0xd372c07e.
//
// Solidity: function setAcuContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetAcuContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setAcuContract", _address)
	return consumed, err
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetCommitteeSize(opts *runOptions, _size *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setCommitteeSize", _size)
	return consumed, err
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetEpochPeriod(opts *runOptions, _period *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setEpochPeriod", _period)
	return consumed, err
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetMinimumBaseFee(opts *runOptions, _price *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setMinimumBaseFee", _price)
	return consumed, err
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetOperatorAccount(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setOperatorAccount", _account)
	return consumed, err
}

// SetOracleContract is a paid mutator transaction binding the contract method 0x496ccd9b.
//
// Solidity: function setOracleContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetOracleContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setOracleContract", _address)
	return consumed, err
}

// SetStabilizationContract is a paid mutator transaction binding the contract method 0xcfd19fb9.
//
// Solidity: function setStabilizationContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetStabilizationContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setStabilizationContract", _address)
	return consumed, err
}

// SetSupplyControlContract is a paid mutator transaction binding the contract method 0xb3ecbadd.
//
// Solidity: function setSupplyControlContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetSupplyControlContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setSupplyControlContract", _address)
	return consumed, err
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetTreasuryAccount(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setTreasuryAccount", _account)
	return consumed, err
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetTreasuryFee(opts *runOptions, _treasuryFee *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setTreasuryFee", _treasuryFee)
	return consumed, err
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetUnbondingPeriod(opts *runOptions, _period *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setUnbondingPeriod", _period)
	return consumed, err
}

// SetUpgradeManagerContract is a paid mutator transaction binding the contract method 0xceaad455.
//
// Solidity: function setUpgradeManagerContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) SetUpgradeManagerContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setUpgradeManagerContract", _address)
	return consumed, err
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Transfer(opts *runOptions, _recipient common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "transfer", _recipient, _amount)
	return consumed, err
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool)
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) TransferFrom(opts *runOptions, _sender common.Address, _recipient common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "transferFrom", _sender, _recipient, _amount)
	return consumed, err
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Unbond(opts *runOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "unbond", _validator, _amount)
	return consumed, err
}

// UpdateValidatorAndTransferSlashedFunds is a paid mutator transaction binding the contract method 0x35be16e0.
//
// Solidity: function updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8) _val) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) UpdateValidatorAndTransferSlashedFunds(opts *runOptions, _val AutonityValidator) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "updateValidatorAndTransferSlashedFunds", _val)
	return consumed, err
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) UpgradeContract(opts *runOptions, _bytecode []byte, _abi string) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "upgradeContract", _bytecode, _abi)
	return consumed, err
}

// Fallback is a paid mutator transaction binding the contract fallback function.
// WARNING! UNTESTED
// Solidity: fallback() payable returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Fallback(opts *runOptions, calldata []byte) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "", calldata)
	return consumed, err
}

// Receive is a paid mutator transaction binding the contract receive function.
// WARNING! UNTESTED
// Solidity: receive() payable returns()
func (_AutonityUpgradeTest *AutonityUpgradeTestTransactor) Receive(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "")
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// AutonityUpgradeTestActivatedValidatorIterator is returned from FilterActivatedValidator and is used to iterate over the raw logs and unpacked data for ActivatedValidator events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestActivatedValidatorIterator struct {
			Event *AutonityUpgradeTestActivatedValidator // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestActivatedValidatorIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestActivatedValidator)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestActivatedValidator)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestActivatedValidatorIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestActivatedValidatorIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestActivatedValidator represents a ActivatedValidator event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestActivatedValidator struct {
			Treasury common.Address;
			Addr common.Address;
			EffectiveBlock *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterActivatedValidator is a free log retrieval operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
		//
		// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterActivatedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*AutonityUpgradeTestActivatedValidatorIterator, error) {

			var treasuryRule []interface{}
			for _, treasuryItem := range treasury {
				treasuryRule = append(treasuryRule, treasuryItem)
			}
			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "ActivatedValidator", treasuryRule, addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestActivatedValidatorIterator{contract: _AutonityUpgradeTest.contract, event: "ActivatedValidator", logs: logs, sub: sub}, nil
 		}

		// WatchActivatedValidator is a free log subscription operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
		//
		// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchActivatedValidator(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestActivatedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

			var treasuryRule []interface{}
			for _, treasuryItem := range treasury {
				treasuryRule = append(treasuryRule, treasuryItem)
			}
			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "ActivatedValidator", treasuryRule, addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestActivatedValidator)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "ActivatedValidator", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseActivatedValidator is a log parse operation binding the contract event 0x60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b5.
		//
		// Solidity: event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseActivatedValidator(log types.Log) (*AutonityUpgradeTestActivatedValidator, error) {
			event := new(AutonityUpgradeTestActivatedValidator)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "ActivatedValidator", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestApprovalIterator struct {
			Event *AutonityUpgradeTestApproval // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestApprovalIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestApproval)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestApproval)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestApprovalIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestApprovalIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestApproval represents a Approval event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestApproval struct {
			Owner common.Address;
			Spender common.Address;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AutonityUpgradeTestApprovalIterator, error) {

			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var spenderRule []interface{}
			for _, spenderItem := range spender {
				spenderRule = append(spenderRule, spenderItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestApprovalIterator{contract: _AutonityUpgradeTest.contract, event: "Approval", logs: logs, sub: sub}, nil
 		}

		// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var spenderRule []interface{}
			for _, spenderItem := range spender {
				spenderRule = append(spenderRule, spenderItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestApproval)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "Approval", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseApproval(log types.Log) (*AutonityUpgradeTestApproval, error) {
			event := new(AutonityUpgradeTestApproval)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "Approval", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestBondingRejectedIterator is returned from FilterBondingRejected and is used to iterate over the raw logs and unpacked data for BondingRejected events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestBondingRejectedIterator struct {
			Event *AutonityUpgradeTestBondingRejected // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestBondingRejectedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestBondingRejected)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestBondingRejected)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestBondingRejectedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestBondingRejectedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestBondingRejected represents a BondingRejected event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestBondingRejected struct {
			Delegator common.Address;
			Delegatee common.Address;
			Amount *big.Int;
			State uint8;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterBondingRejected is a free log retrieval operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
		//
		// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount, uint8 state)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterBondingRejected(opts *bind.FilterOpts) (*AutonityUpgradeTestBondingRejectedIterator, error) {






			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "BondingRejected")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestBondingRejectedIterator{contract: _AutonityUpgradeTest.contract, event: "BondingRejected", logs: logs, sub: sub}, nil
 		}

		// WatchBondingRejected is a free log subscription operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
		//
		// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount, uint8 state)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchBondingRejected(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestBondingRejected) (event.Subscription, error) {






			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "BondingRejected")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestBondingRejected)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "BondingRejected", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseBondingRejected is a log parse operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
		//
		// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount, uint8 state)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseBondingRejected(log types.Log) (*AutonityUpgradeTestBondingRejected, error) {
			event := new(AutonityUpgradeTestBondingRejected)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "BondingRejected", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestBurnedStakeIterator is returned from FilterBurnedStake and is used to iterate over the raw logs and unpacked data for BurnedStake events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestBurnedStakeIterator struct {
			Event *AutonityUpgradeTestBurnedStake // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestBurnedStakeIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestBurnedStake)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestBurnedStake)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestBurnedStakeIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestBurnedStakeIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestBurnedStake represents a BurnedStake event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestBurnedStake struct {
			Addr common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterBurnedStake is a free log retrieval operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
		//
		// Solidity: event BurnedStake(address indexed addr, uint256 amount)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterBurnedStake(opts *bind.FilterOpts, addr []common.Address) (*AutonityUpgradeTestBurnedStakeIterator, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "BurnedStake", addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestBurnedStakeIterator{contract: _AutonityUpgradeTest.contract, event: "BurnedStake", logs: logs, sub: sub}, nil
 		}

		// WatchBurnedStake is a free log subscription operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
		//
		// Solidity: event BurnedStake(address indexed addr, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchBurnedStake(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestBurnedStake, addr []common.Address) (event.Subscription, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "BurnedStake", addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestBurnedStake)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "BurnedStake", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseBurnedStake is a log parse operation binding the contract event 0x5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3.
		//
		// Solidity: event BurnedStake(address indexed addr, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseBurnedStake(log types.Log) (*AutonityUpgradeTestBurnedStake, error) {
			event := new(AutonityUpgradeTestBurnedStake)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "BurnedStake", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestCommissionRateChangeIterator is returned from FilterCommissionRateChange and is used to iterate over the raw logs and unpacked data for CommissionRateChange events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestCommissionRateChangeIterator struct {
			Event *AutonityUpgradeTestCommissionRateChange // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestCommissionRateChangeIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestCommissionRateChange)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestCommissionRateChange)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestCommissionRateChangeIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestCommissionRateChangeIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestCommissionRateChange represents a CommissionRateChange event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestCommissionRateChange struct {
			Validator common.Address;
			Rate *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterCommissionRateChange is a free log retrieval operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
		//
		// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterCommissionRateChange(opts *bind.FilterOpts, validator []common.Address) (*AutonityUpgradeTestCommissionRateChangeIterator, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "CommissionRateChange", validatorRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestCommissionRateChangeIterator{contract: _AutonityUpgradeTest.contract, event: "CommissionRateChange", logs: logs, sub: sub}, nil
 		}

		// WatchCommissionRateChange is a free log subscription operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
		//
		// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchCommissionRateChange(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestCommissionRateChange, validator []common.Address) (event.Subscription, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "CommissionRateChange", validatorRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestCommissionRateChange)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "CommissionRateChange", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseCommissionRateChange is a log parse operation binding the contract event 0x4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf.
		//
		// Solidity: event CommissionRateChange(address indexed validator, uint256 rate)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseCommissionRateChange(log types.Log) (*AutonityUpgradeTestCommissionRateChange, error) {
			event := new(AutonityUpgradeTestCommissionRateChange)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "CommissionRateChange", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestEpochPeriodUpdatedIterator is returned from FilterEpochPeriodUpdated and is used to iterate over the raw logs and unpacked data for EpochPeriodUpdated events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestEpochPeriodUpdatedIterator struct {
			Event *AutonityUpgradeTestEpochPeriodUpdated // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestEpochPeriodUpdatedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestEpochPeriodUpdated)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestEpochPeriodUpdated)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestEpochPeriodUpdatedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestEpochPeriodUpdatedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestEpochPeriodUpdated represents a EpochPeriodUpdated event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestEpochPeriodUpdated struct {
			Period *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterEpochPeriodUpdated is a free log retrieval operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
		//
		// Solidity: event EpochPeriodUpdated(uint256 period)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterEpochPeriodUpdated(opts *bind.FilterOpts) (*AutonityUpgradeTestEpochPeriodUpdatedIterator, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "EpochPeriodUpdated")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestEpochPeriodUpdatedIterator{contract: _AutonityUpgradeTest.contract, event: "EpochPeriodUpdated", logs: logs, sub: sub}, nil
 		}

		// WatchEpochPeriodUpdated is a free log subscription operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
		//
		// Solidity: event EpochPeriodUpdated(uint256 period)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchEpochPeriodUpdated(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestEpochPeriodUpdated) (event.Subscription, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "EpochPeriodUpdated")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestEpochPeriodUpdated)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "EpochPeriodUpdated", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseEpochPeriodUpdated is a log parse operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
		//
		// Solidity: event EpochPeriodUpdated(uint256 period)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseEpochPeriodUpdated(log types.Log) (*AutonityUpgradeTestEpochPeriodUpdated, error) {
			event := new(AutonityUpgradeTestEpochPeriodUpdated)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "EpochPeriodUpdated", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestMinimumBaseFeeUpdatedIterator is returned from FilterMinimumBaseFeeUpdated and is used to iterate over the raw logs and unpacked data for MinimumBaseFeeUpdated events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestMinimumBaseFeeUpdatedIterator struct {
			Event *AutonityUpgradeTestMinimumBaseFeeUpdated // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestMinimumBaseFeeUpdatedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestMinimumBaseFeeUpdated)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestMinimumBaseFeeUpdated)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestMinimumBaseFeeUpdatedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestMinimumBaseFeeUpdatedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestMinimumBaseFeeUpdated represents a MinimumBaseFeeUpdated event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestMinimumBaseFeeUpdated struct {
			GasPrice *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterMinimumBaseFeeUpdated is a free log retrieval operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
		//
		// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterMinimumBaseFeeUpdated(opts *bind.FilterOpts) (*AutonityUpgradeTestMinimumBaseFeeUpdatedIterator, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "MinimumBaseFeeUpdated")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestMinimumBaseFeeUpdatedIterator{contract: _AutonityUpgradeTest.contract, event: "MinimumBaseFeeUpdated", logs: logs, sub: sub}, nil
 		}

		// WatchMinimumBaseFeeUpdated is a free log subscription operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
		//
		// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchMinimumBaseFeeUpdated(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestMinimumBaseFeeUpdated) (event.Subscription, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "MinimumBaseFeeUpdated")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestMinimumBaseFeeUpdated)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "MinimumBaseFeeUpdated", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseMinimumBaseFeeUpdated is a log parse operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
		//
		// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseMinimumBaseFeeUpdated(log types.Log) (*AutonityUpgradeTestMinimumBaseFeeUpdated, error) {
			event := new(AutonityUpgradeTestMinimumBaseFeeUpdated)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "MinimumBaseFeeUpdated", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestMintedStakeIterator is returned from FilterMintedStake and is used to iterate over the raw logs and unpacked data for MintedStake events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestMintedStakeIterator struct {
			Event *AutonityUpgradeTestMintedStake // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestMintedStakeIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestMintedStake)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestMintedStake)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestMintedStakeIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestMintedStakeIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestMintedStake represents a MintedStake event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestMintedStake struct {
			Addr common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterMintedStake is a free log retrieval operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
		//
		// Solidity: event MintedStake(address indexed addr, uint256 amount)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterMintedStake(opts *bind.FilterOpts, addr []common.Address) (*AutonityUpgradeTestMintedStakeIterator, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "MintedStake", addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestMintedStakeIterator{contract: _AutonityUpgradeTest.contract, event: "MintedStake", logs: logs, sub: sub}, nil
 		}

		// WatchMintedStake is a free log subscription operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
		//
		// Solidity: event MintedStake(address indexed addr, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchMintedStake(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestMintedStake, addr []common.Address) (event.Subscription, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "MintedStake", addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestMintedStake)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "MintedStake", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseMintedStake is a log parse operation binding the contract event 0x48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf.
		//
		// Solidity: event MintedStake(address indexed addr, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseMintedStake(log types.Log) (*AutonityUpgradeTestMintedStake, error) {
			event := new(AutonityUpgradeTestMintedStake)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "MintedStake", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestNewBondingRequestIterator is returned from FilterNewBondingRequest and is used to iterate over the raw logs and unpacked data for NewBondingRequest events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestNewBondingRequestIterator struct {
			Event *AutonityUpgradeTestNewBondingRequest // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestNewBondingRequestIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestNewBondingRequest)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestNewBondingRequest)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestNewBondingRequestIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestNewBondingRequestIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestNewBondingRequest represents a NewBondingRequest event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestNewBondingRequest struct {
			Validator common.Address;
			Delegator common.Address;
			SelfBonded bool;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewBondingRequest is a free log retrieval operation binding the contract event 0xc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d.
		//
		// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterNewBondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*AutonityUpgradeTestNewBondingRequestIterator, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}
			var delegatorRule []interface{}
			for _, delegatorItem := range delegator {
				delegatorRule = append(delegatorRule, delegatorItem)
			}



			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "NewBondingRequest", validatorRule, delegatorRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestNewBondingRequestIterator{contract: _AutonityUpgradeTest.contract, event: "NewBondingRequest", logs: logs, sub: sub}, nil
 		}

		// WatchNewBondingRequest is a free log subscription operation binding the contract event 0xc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d.
		//
		// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchNewBondingRequest(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestNewBondingRequest, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}
			var delegatorRule []interface{}
			for _, delegatorItem := range delegator {
				delegatorRule = append(delegatorRule, delegatorItem)
			}



			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "NewBondingRequest", validatorRule, delegatorRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestNewBondingRequest)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "NewBondingRequest", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewBondingRequest is a log parse operation binding the contract event 0xc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d.
		//
		// Solidity: event NewBondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseNewBondingRequest(log types.Log) (*AutonityUpgradeTestNewBondingRequest, error) {
			event := new(AutonityUpgradeTestNewBondingRequest)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "NewBondingRequest", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestNewEpochIterator struct {
			Event *AutonityUpgradeTestNewEpoch // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestNewEpochIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestNewEpoch)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestNewEpoch)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestNewEpochIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestNewEpochIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestNewEpoch represents a NewEpoch event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestNewEpoch struct {
			Epoch *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewEpoch is a free log retrieval operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
		//
		// Solidity: event NewEpoch(uint256 epoch)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*AutonityUpgradeTestNewEpochIterator, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "NewEpoch")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestNewEpochIterator{contract: _AutonityUpgradeTest.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
 		}

		// WatchNewEpoch is a free log subscription operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
		//
		// Solidity: event NewEpoch(uint256 epoch)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestNewEpoch) (event.Subscription, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "NewEpoch")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestNewEpoch)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "NewEpoch", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewEpoch is a log parse operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
		//
		// Solidity: event NewEpoch(uint256 epoch)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseNewEpoch(log types.Log) (*AutonityUpgradeTestNewEpoch, error) {
			event := new(AutonityUpgradeTestNewEpoch)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "NewEpoch", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestNewUnbondingRequestIterator is returned from FilterNewUnbondingRequest and is used to iterate over the raw logs and unpacked data for NewUnbondingRequest events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestNewUnbondingRequestIterator struct {
			Event *AutonityUpgradeTestNewUnbondingRequest // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestNewUnbondingRequestIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestNewUnbondingRequest)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestNewUnbondingRequest)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestNewUnbondingRequestIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestNewUnbondingRequestIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestNewUnbondingRequest represents a NewUnbondingRequest event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestNewUnbondingRequest struct {
			Validator common.Address;
			Delegator common.Address;
			SelfBonded bool;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewUnbondingRequest is a free log retrieval operation binding the contract event 0x63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc.
		//
		// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterNewUnbondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*AutonityUpgradeTestNewUnbondingRequestIterator, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}
			var delegatorRule []interface{}
			for _, delegatorItem := range delegator {
				delegatorRule = append(delegatorRule, delegatorItem)
			}



			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "NewUnbondingRequest", validatorRule, delegatorRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestNewUnbondingRequestIterator{contract: _AutonityUpgradeTest.contract, event: "NewUnbondingRequest", logs: logs, sub: sub}, nil
 		}

		// WatchNewUnbondingRequest is a free log subscription operation binding the contract event 0x63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc.
		//
		// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchNewUnbondingRequest(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestNewUnbondingRequest, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

			var validatorRule []interface{}
			for _, validatorItem := range validator {
				validatorRule = append(validatorRule, validatorItem)
			}
			var delegatorRule []interface{}
			for _, delegatorItem := range delegator {
				delegatorRule = append(delegatorRule, delegatorItem)
			}



			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "NewUnbondingRequest", validatorRule, delegatorRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestNewUnbondingRequest)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "NewUnbondingRequest", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewUnbondingRequest is a log parse operation binding the contract event 0x63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc.
		//
		// Solidity: event NewUnbondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseNewUnbondingRequest(log types.Log) (*AutonityUpgradeTestNewUnbondingRequest, error) {
			event := new(AutonityUpgradeTestNewUnbondingRequest)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "NewUnbondingRequest", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestPausedValidatorIterator is returned from FilterPausedValidator and is used to iterate over the raw logs and unpacked data for PausedValidator events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestPausedValidatorIterator struct {
			Event *AutonityUpgradeTestPausedValidator // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestPausedValidatorIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestPausedValidator)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestPausedValidator)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestPausedValidatorIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestPausedValidatorIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestPausedValidator represents a PausedValidator event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestPausedValidator struct {
			Treasury common.Address;
			Addr common.Address;
			EffectiveBlock *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterPausedValidator is a free log retrieval operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
		//
		// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterPausedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*AutonityUpgradeTestPausedValidatorIterator, error) {

			var treasuryRule []interface{}
			for _, treasuryItem := range treasury {
				treasuryRule = append(treasuryRule, treasuryItem)
			}
			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "PausedValidator", treasuryRule, addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestPausedValidatorIterator{contract: _AutonityUpgradeTest.contract, event: "PausedValidator", logs: logs, sub: sub}, nil
 		}

		// WatchPausedValidator is a free log subscription operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
		//
		// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchPausedValidator(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestPausedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

			var treasuryRule []interface{}
			for _, treasuryItem := range treasury {
				treasuryRule = append(treasuryRule, treasuryItem)
			}
			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "PausedValidator", treasuryRule, addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestPausedValidator)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "PausedValidator", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParsePausedValidator is a log parse operation binding the contract event 0x75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c.
		//
		// Solidity: event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParsePausedValidator(log types.Log) (*AutonityUpgradeTestPausedValidator, error) {
			event := new(AutonityUpgradeTestPausedValidator)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "PausedValidator", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestRegisteredValidatorIterator is returned from FilterRegisteredValidator and is used to iterate over the raw logs and unpacked data for RegisteredValidator events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestRegisteredValidatorIterator struct {
			Event *AutonityUpgradeTestRegisteredValidator // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestRegisteredValidatorIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestRegisteredValidator)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestRegisteredValidator)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestRegisteredValidatorIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestRegisteredValidatorIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestRegisteredValidator represents a RegisteredValidator event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestRegisteredValidator struct {
			Treasury common.Address;
			Addr common.Address;
			OracleAddress common.Address;
			Enode string;
			LiquidContract common.Address;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterRegisteredValidator is a free log retrieval operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
		//
		// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterRegisteredValidator(opts *bind.FilterOpts) (*AutonityUpgradeTestRegisteredValidatorIterator, error) {







			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "RegisteredValidator")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestRegisteredValidatorIterator{contract: _AutonityUpgradeTest.contract, event: "RegisteredValidator", logs: logs, sub: sub}, nil
 		}

		// WatchRegisteredValidator is a free log subscription operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
		//
		// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchRegisteredValidator(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestRegisteredValidator) (event.Subscription, error) {







			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "RegisteredValidator")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestRegisteredValidator)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "RegisteredValidator", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseRegisteredValidator is a log parse operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
		//
		// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseRegisteredValidator(log types.Log) (*AutonityUpgradeTestRegisteredValidator, error) {
			event := new(AutonityUpgradeTestRegisteredValidator)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "RegisteredValidator", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestRewardedIterator is returned from FilterRewarded and is used to iterate over the raw logs and unpacked data for Rewarded events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestRewardedIterator struct {
			Event *AutonityUpgradeTestRewarded // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestRewardedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestRewarded)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestRewarded)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestRewardedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestRewardedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestRewarded represents a Rewarded event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestRewarded struct {
			Addr common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterRewarded is a free log retrieval operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
		//
		// Solidity: event Rewarded(address indexed addr, uint256 amount)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterRewarded(opts *bind.FilterOpts, addr []common.Address) (*AutonityUpgradeTestRewardedIterator, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "Rewarded", addrRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestRewardedIterator{contract: _AutonityUpgradeTest.contract, event: "Rewarded", logs: logs, sub: sub}, nil
 		}

		// WatchRewarded is a free log subscription operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
		//
		// Solidity: event Rewarded(address indexed addr, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchRewarded(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestRewarded, addr []common.Address) (event.Subscription, error) {

			var addrRule []interface{}
			for _, addrItem := range addr {
				addrRule = append(addrRule, addrItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "Rewarded", addrRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestRewarded)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "Rewarded", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseRewarded is a log parse operation binding the contract event 0xb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe563.
		//
		// Solidity: event Rewarded(address indexed addr, uint256 amount)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseRewarded(log types.Log) (*AutonityUpgradeTestRewarded, error) {
			event := new(AutonityUpgradeTestRewarded)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "Rewarded", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// AutonityUpgradeTestTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestTransferIterator struct {
			Event *AutonityUpgradeTestTransfer // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *AutonityUpgradeTestTransferIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(AutonityUpgradeTestTransfer)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(AutonityUpgradeTestTransfer)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *AutonityUpgradeTestTransferIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *AutonityUpgradeTestTransferIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// AutonityUpgradeTestTransfer represents a Transfer event raised by the AutonityUpgradeTest contract.
		type AutonityUpgradeTestTransfer struct {
			From common.Address;
			To common.Address;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
 		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutonityUpgradeTestTransferIterator, error) {

			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestTransferIterator{contract: _AutonityUpgradeTest.contract, event: "Transfer", logs: logs, sub: sub}, nil
 		}

		// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}


			logs, sub, err := _AutonityUpgradeTest.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(AutonityUpgradeTestTransfer)
						if err := _AutonityUpgradeTest.contract.UnpackLog(event, "Transfer", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
		func (_AutonityUpgradeTest *AutonityUpgradeTestFilterer) ParseTransfer(log types.Log) (*AutonityUpgradeTestTransfer, error) {
			event := new(AutonityUpgradeTestTransfer)
			if err := _AutonityUpgradeTest.contract.UnpackLog(event, "Transfer", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// BytesLibMetaData contains all meta data concerning the BytesLib contract.
var BytesLibMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220a382116c91b2d99ddbd669d16236b0462d4240ef2dacde934e5c8103df52ee0a64736f6c63430008150033",
}

// BytesLibABI is the input ABI used to generate the binding from.
// Deprecated: Use BytesLibMetaData.ABI instead.
var BytesLibABI = BytesLibMetaData.ABI

// BytesLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BytesLibMetaData.Bin instead.
var BytesLibBin = BytesLibMetaData.Bin

// DeployBytesLib deploys a new Ethereum contract, binding an instance of BytesLib to it.
func (r *runner) deployBytesLib(opts *runOptions) (common.Address, uint64, *BytesLib, error) {
	parsed, err := BytesLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(BytesLibBin))
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &BytesLib{BytesLibCaller: BytesLibCaller{contract: c}, BytesLibTransactor: BytesLibTransactor{contract: c}, BytesLibFilterer: BytesLibFilterer{contract: c}}, nil
}

// BytesLib is an auto generated Go binding around an Ethereum contract.
type BytesLib struct {
	BytesLibCaller     // Read-only binding to the contract
	BytesLibTransactor // Write-only binding to the contract
	BytesLibFilterer   // Log filterer for contract events
}

// BytesLibCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type BytesLibCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// BytesLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BytesLibTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// BytesLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BytesLibFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

/* EVENTS ARE NOT YET SUPPORTED

 */

// HelpersMetaData contains all meta data concerning the Helpers contract.
var HelpersMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220ccdb043a0f3f32a8e520ed4d013ed348218883d602fba52fbd74595fbecf146664736f6c63430008150033",
}

// HelpersABI is the input ABI used to generate the binding from.
// Deprecated: Use HelpersMetaData.ABI instead.
var HelpersABI = HelpersMetaData.ABI

// HelpersBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use HelpersMetaData.Bin instead.
var HelpersBin = HelpersMetaData.Bin

// DeployHelpers deploys a new Ethereum contract, binding an instance of Helpers to it.
func (r *runner) deployHelpers(opts *runOptions) (common.Address, uint64, *Helpers, error) {
	parsed, err := HelpersMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(HelpersBin))
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &Helpers{HelpersCaller: HelpersCaller{contract: c}, HelpersTransactor: HelpersTransactor{contract: c}, HelpersFilterer: HelpersFilterer{contract: c}}, nil
}

// Helpers is an auto generated Go binding around an Ethereum contract.
type Helpers struct {
	HelpersCaller     // Read-only binding to the contract
	HelpersTransactor // Write-only binding to the contract
	HelpersFilterer   // Log filterer for contract events
}

// HelpersCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type HelpersCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// HelpersTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HelpersTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// HelpersFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HelpersFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

/* EVENTS ARE NOT YET SUPPORTED

 */

// IACUMetaData contains all meta data concerning the IACU contract.
var IACUMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"update\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"b3ab15fb": "setOperator(address)",
		"7adbf973": "setOracle(address)",
		"a2e62045": "update()",
	},
}

// IACUABI is the input ABI used to generate the binding from.
// Deprecated: Use IACUMetaData.ABI instead.
var IACUABI = IACUMetaData.ABI

// Deprecated: Use IACUMetaData.Sigs instead.
// IACUFuncSigs maps the 4-byte function signature to its string representation.
var IACUFuncSigs = IACUMetaData.Sigs

// IACU is an auto generated Go binding around an Ethereum contract.
type IACU struct {
	IACUCaller     // Read-only binding to the contract
	IACUTransactor // Write-only binding to the contract
	IACUFilterer   // Log filterer for contract events
}

// IACUCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type IACUCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// IACUTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IACUTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// IACUFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IACUFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACUTransactor) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _IACU.call(opts, "setOperator", operator)
	return consumed, err
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACUTransactor) SetOracle(opts *runOptions, oracle common.Address) (uint64, error) {
	_, consumed, err := _IACU.call(opts, "setOracle", oracle)
	return consumed, err
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACUTransactor) Update(opts *runOptions) (uint64, error) {
	_, consumed, err := _IACU.call(opts, "update")
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

 */

// IAccountabilityMetaData contains all meta data concerning the IAccountability contract.
var IAccountabilityMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"InnocenceProven\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewAccusation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewFaultProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"releaseBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isJailbound\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"SlashingEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_epochEnd\",\"type\":\"bool\"}],\"name\":\"finalize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newPeriod\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1de9d9b6": "distributeRewards(address)",
		"6c9789b0": "finalize(bool)",
		"6b5f444c": "setEpochPeriod(uint256)",
	},
}

// IAccountabilityABI is the input ABI used to generate the binding from.
// Deprecated: Use IAccountabilityMetaData.ABI instead.
var IAccountabilityABI = IAccountabilityMetaData.ABI

// Deprecated: Use IAccountabilityMetaData.Sigs instead.
// IAccountabilityFuncSigs maps the 4-byte function signature to its string representation.
var IAccountabilityFuncSigs = IAccountabilityMetaData.Sigs

// IAccountability is an auto generated Go binding around an Ethereum contract.
type IAccountability struct {
	IAccountabilityCaller     // Read-only binding to the contract
	IAccountabilityTransactor // Write-only binding to the contract
	IAccountabilityFilterer   // Log filterer for contract events
}

// IAccountabilityCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type IAccountabilityCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// IAccountabilityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAccountabilityTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// IAccountabilityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAccountabilityFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_IAccountability *IAccountabilityTransactor) DistributeRewards(opts *runOptions, _validator common.Address) (uint64, error) {
	_, consumed, err := _IAccountability.call(opts, "distributeRewards", _validator)
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_IAccountability *IAccountabilityTransactor) Finalize(opts *runOptions, _epochEnd bool) (uint64, error) {
	_, consumed, err := _IAccountability.call(opts, "finalize", _epochEnd)
	return consumed, err
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_IAccountability *IAccountabilityTransactor) SetEpochPeriod(opts *runOptions, _newPeriod *big.Int) (uint64, error) {
	_, consumed, err := _IAccountability.call(opts, "setEpochPeriod", _newPeriod)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// IAccountabilityInnocenceProvenIterator is returned from FilterInnocenceProven and is used to iterate over the raw logs and unpacked data for InnocenceProven events raised by the IAccountability contract.
		type IAccountabilityInnocenceProvenIterator struct {
			Event *IAccountabilityInnocenceProven // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IAccountabilityInnocenceProvenIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IAccountabilityInnocenceProven)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IAccountabilityInnocenceProven)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IAccountabilityInnocenceProvenIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IAccountabilityInnocenceProvenIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IAccountabilityInnocenceProven represents a InnocenceProven event raised by the IAccountability contract.
		type IAccountabilityInnocenceProven struct {
			Offender common.Address;
			Id *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterInnocenceProven is a free log retrieval operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
		//
		// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
 		func (_IAccountability *IAccountabilityFilterer) FilterInnocenceProven(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityInnocenceProvenIterator, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}


			logs, sub, err := _IAccountability.contract.FilterLogs(opts, "InnocenceProven", _offenderRule)
			if err != nil {
				return nil, err
			}
			return &IAccountabilityInnocenceProvenIterator{contract: _IAccountability.contract, event: "InnocenceProven", logs: logs, sub: sub}, nil
 		}

		// WatchInnocenceProven is a free log subscription operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
		//
		// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
		func (_IAccountability *IAccountabilityFilterer) WatchInnocenceProven(opts *bind.WatchOpts, sink chan<- *IAccountabilityInnocenceProven, _offender []common.Address) (event.Subscription, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}


			logs, sub, err := _IAccountability.contract.WatchLogs(opts, "InnocenceProven", _offenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IAccountabilityInnocenceProven)
						if err := _IAccountability.contract.UnpackLog(event, "InnocenceProven", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseInnocenceProven is a log parse operation binding the contract event 0x1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f.
		//
		// Solidity: event InnocenceProven(address indexed _offender, uint256 _id)
		func (_IAccountability *IAccountabilityFilterer) ParseInnocenceProven(log types.Log) (*IAccountabilityInnocenceProven, error) {
			event := new(IAccountabilityInnocenceProven)
			if err := _IAccountability.contract.UnpackLog(event, "InnocenceProven", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// IAccountabilityNewAccusationIterator is returned from FilterNewAccusation and is used to iterate over the raw logs and unpacked data for NewAccusation events raised by the IAccountability contract.
		type IAccountabilityNewAccusationIterator struct {
			Event *IAccountabilityNewAccusation // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IAccountabilityNewAccusationIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IAccountabilityNewAccusation)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IAccountabilityNewAccusation)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IAccountabilityNewAccusationIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IAccountabilityNewAccusationIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IAccountabilityNewAccusation represents a NewAccusation event raised by the IAccountability contract.
		type IAccountabilityNewAccusation struct {
			Offender common.Address;
			Severity *big.Int;
			Id *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewAccusation is a free log retrieval operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
		//
		// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
 		func (_IAccountability *IAccountabilityFilterer) FilterNewAccusation(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityNewAccusationIterator, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}



			logs, sub, err := _IAccountability.contract.FilterLogs(opts, "NewAccusation", _offenderRule)
			if err != nil {
				return nil, err
			}
			return &IAccountabilityNewAccusationIterator{contract: _IAccountability.contract, event: "NewAccusation", logs: logs, sub: sub}, nil
 		}

		// WatchNewAccusation is a free log subscription operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
		//
		// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
		func (_IAccountability *IAccountabilityFilterer) WatchNewAccusation(opts *bind.WatchOpts, sink chan<- *IAccountabilityNewAccusation, _offender []common.Address) (event.Subscription, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}



			logs, sub, err := _IAccountability.contract.WatchLogs(opts, "NewAccusation", _offenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IAccountabilityNewAccusation)
						if err := _IAccountability.contract.UnpackLog(event, "NewAccusation", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewAccusation is a log parse operation binding the contract event 0x2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba40351.
		//
		// Solidity: event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id)
		func (_IAccountability *IAccountabilityFilterer) ParseNewAccusation(log types.Log) (*IAccountabilityNewAccusation, error) {
			event := new(IAccountabilityNewAccusation)
			if err := _IAccountability.contract.UnpackLog(event, "NewAccusation", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// IAccountabilityNewFaultProofIterator is returned from FilterNewFaultProof and is used to iterate over the raw logs and unpacked data for NewFaultProof events raised by the IAccountability contract.
		type IAccountabilityNewFaultProofIterator struct {
			Event *IAccountabilityNewFaultProof // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IAccountabilityNewFaultProofIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IAccountabilityNewFaultProof)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IAccountabilityNewFaultProof)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IAccountabilityNewFaultProofIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IAccountabilityNewFaultProofIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IAccountabilityNewFaultProof represents a NewFaultProof event raised by the IAccountability contract.
		type IAccountabilityNewFaultProof struct {
			Offender common.Address;
			Severity *big.Int;
			Id *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewFaultProof is a free log retrieval operation binding the contract event 0x6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f.
		//
		// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id)
 		func (_IAccountability *IAccountabilityFilterer) FilterNewFaultProof(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityNewFaultProofIterator, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}



			logs, sub, err := _IAccountability.contract.FilterLogs(opts, "NewFaultProof", _offenderRule)
			if err != nil {
				return nil, err
			}
			return &IAccountabilityNewFaultProofIterator{contract: _IAccountability.contract, event: "NewFaultProof", logs: logs, sub: sub}, nil
 		}

		// WatchNewFaultProof is a free log subscription operation binding the contract event 0x6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f.
		//
		// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id)
		func (_IAccountability *IAccountabilityFilterer) WatchNewFaultProof(opts *bind.WatchOpts, sink chan<- *IAccountabilityNewFaultProof, _offender []common.Address) (event.Subscription, error) {

			var _offenderRule []interface{}
			for _, _offenderItem := range _offender {
				_offenderRule = append(_offenderRule, _offenderItem)
			}



			logs, sub, err := _IAccountability.contract.WatchLogs(opts, "NewFaultProof", _offenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IAccountabilityNewFaultProof)
						if err := _IAccountability.contract.UnpackLog(event, "NewFaultProof", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewFaultProof is a log parse operation binding the contract event 0x6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f.
		//
		// Solidity: event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id)
		func (_IAccountability *IAccountabilityFilterer) ParseNewFaultProof(log types.Log) (*IAccountabilityNewFaultProof, error) {
			event := new(IAccountabilityNewFaultProof)
			if err := _IAccountability.contract.UnpackLog(event, "NewFaultProof", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// IAccountabilitySlashingEventIterator is returned from FilterSlashingEvent and is used to iterate over the raw logs and unpacked data for SlashingEvent events raised by the IAccountability contract.
		type IAccountabilitySlashingEventIterator struct {
			Event *IAccountabilitySlashingEvent // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IAccountabilitySlashingEventIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IAccountabilitySlashingEvent)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IAccountabilitySlashingEvent)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IAccountabilitySlashingEventIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IAccountabilitySlashingEventIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IAccountabilitySlashingEvent represents a SlashingEvent event raised by the IAccountability contract.
		type IAccountabilitySlashingEvent struct {
			Validator common.Address;
			Amount *big.Int;
			ReleaseBlock *big.Int;
			IsJailbound bool;
			EventId *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterSlashingEvent is a free log retrieval operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
 		func (_IAccountability *IAccountabilityFilterer) FilterSlashingEvent(opts *bind.FilterOpts) (*IAccountabilitySlashingEventIterator, error) {







			logs, sub, err := _IAccountability.contract.FilterLogs(opts, "SlashingEvent")
			if err != nil {
				return nil, err
			}
			return &IAccountabilitySlashingEventIterator{contract: _IAccountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
 		}

		// WatchSlashingEvent is a free log subscription operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
		func (_IAccountability *IAccountabilityFilterer) WatchSlashingEvent(opts *bind.WatchOpts, sink chan<- *IAccountabilitySlashingEvent) (event.Subscription, error) {







			logs, sub, err := _IAccountability.contract.WatchLogs(opts, "SlashingEvent")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IAccountabilitySlashingEvent)
						if err := _IAccountability.contract.UnpackLog(event, "SlashingEvent", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseSlashingEvent is a log parse operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
		func (_IAccountability *IAccountabilityFilterer) ParseSlashingEvent(log types.Log) (*IAccountabilitySlashingEvent, error) {
			event := new(IAccountabilitySlashingEvent)
			if err := _IAccountability.contract.UnpackLog(event, "SlashingEvent", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// IAutonityMetaData contains all meta data concerning the IAutonity contract.
var IAutonityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e7f43c68": "getOperator()",
		"833b1fce": "getOracle()",
	},
}

// IAutonityABI is the input ABI used to generate the binding from.
// Deprecated: Use IAutonityMetaData.ABI instead.
var IAutonityABI = IAutonityMetaData.ABI

// Deprecated: Use IAutonityMetaData.Sigs instead.
// IAutonityFuncSigs maps the 4-byte function signature to its string representation.
var IAutonityFuncSigs = IAutonityMetaData.Sigs

// IAutonity is an auto generated Go binding around an Ethereum contract.
type IAutonity struct {
	IAutonityCaller     // Read-only binding to the contract
	IAutonityTransactor // Write-only binding to the contract
	IAutonityFilterer   // Log filterer for contract events
}

// IAutonityCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type IAutonityCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// IAutonityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAutonityTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// IAutonityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAutonityFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonityCaller) GetOperator(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _IAutonity.call(opts, "getOperator")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_IAutonity *IAutonityCaller) GetOracle(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _IAutonity.call(opts, "getOracle")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

/* EVENTS ARE NOT YET SUPPORTED

 */

// IERC20MetaData contains all meta data concerning the IERC20 contract.
var IERC20MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
	},
}

// IERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC20MetaData.ABI instead.
var IERC20ABI = IERC20MetaData.ABI

// Deprecated: Use IERC20MetaData.Sigs instead.
// IERC20FuncSigs maps the 4-byte function signature to its string representation.
var IERC20FuncSigs = IERC20MetaData.Sigs

// IERC20 is an auto generated Go binding around an Ethereum contract.
type IERC20 struct {
	IERC20Caller     // Read-only binding to the contract
	IERC20Transactor // Write-only binding to the contract
	IERC20Filterer   // Log filterer for contract events
}

// IERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type IERC20Caller struct {
	*contract // Generic contract wrapper for the low level calls
}

// IERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type IERC20Transactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// IERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IERC20Filterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Caller) Allowance(opts *runOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _IERC20.call(opts, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Caller) BalanceOf(opts *runOptions, account common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _IERC20.call(opts, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Caller) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _IERC20.call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Approve(opts *runOptions, spender common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _IERC20.call(opts, "approve", spender, amount)
	return consumed, err
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Transfer(opts *runOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _IERC20.call(opts, "transfer", recipient, amount)
	return consumed, err
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) TransferFrom(opts *runOptions, sender common.Address, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _IERC20.call(opts, "transferFrom", sender, recipient, amount)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// IERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IERC20 contract.
		type IERC20ApprovalIterator struct {
			Event *IERC20Approval // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IERC20ApprovalIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IERC20Approval)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IERC20Approval)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IERC20ApprovalIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IERC20ApprovalIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IERC20Approval represents a Approval event raised by the IERC20 contract.
		type IERC20Approval struct {
			Owner common.Address;
			Spender common.Address;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
 		func (_IERC20 *IERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IERC20ApprovalIterator, error) {

			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var spenderRule []interface{}
			for _, spenderItem := range spender {
				spenderRule = append(spenderRule, spenderItem)
			}


			logs, sub, err := _IERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
			if err != nil {
				return nil, err
			}
			return &IERC20ApprovalIterator{contract: _IERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
 		}

		// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
		func (_IERC20 *IERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var spenderRule []interface{}
			for _, spenderItem := range spender {
				spenderRule = append(spenderRule, spenderItem)
			}


			logs, sub, err := _IERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IERC20Approval)
						if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
		func (_IERC20 *IERC20Filterer) ParseApproval(log types.Log) (*IERC20Approval, error) {
			event := new(IERC20Approval)
			if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// IERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IERC20 contract.
		type IERC20TransferIterator struct {
			Event *IERC20Transfer // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IERC20TransferIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IERC20Transfer)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IERC20Transfer)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IERC20TransferIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IERC20TransferIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IERC20Transfer represents a Transfer event raised by the IERC20 contract.
		type IERC20Transfer struct {
			From common.Address;
			To common.Address;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
 		func (_IERC20 *IERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IERC20TransferIterator, error) {

			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}


			logs, sub, err := _IERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
			if err != nil {
				return nil, err
			}
			return &IERC20TransferIterator{contract: _IERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
 		}

		// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
		func (_IERC20 *IERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}


			logs, sub, err := _IERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IERC20Transfer)
						if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
		func (_IERC20 *IERC20Filterer) ParseTransfer(log types.Log) (*IERC20Transfer, error) {
			event := new(IERC20Transfer)
			if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// IInflationControllerMetaData contains all meta data concerning the IInflationController contract.
var IInflationControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_currentSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_inflationReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lastEpochTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_currentEpochTime\",\"type\":\"uint256\"}],\"name\":\"calculateSupplyDelta\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"92eff3cd": "calculateSupplyDelta(uint256,uint256,uint256,uint256)",
	},
}

// IInflationControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use IInflationControllerMetaData.ABI instead.
var IInflationControllerABI = IInflationControllerMetaData.ABI

// Deprecated: Use IInflationControllerMetaData.Sigs instead.
// IInflationControllerFuncSigs maps the 4-byte function signature to its string representation.
var IInflationControllerFuncSigs = IInflationControllerMetaData.Sigs

// IInflationController is an auto generated Go binding around an Ethereum contract.
type IInflationController struct {
	IInflationControllerCaller     // Read-only binding to the contract
	IInflationControllerTransactor // Write-only binding to the contract
	IInflationControllerFilterer   // Log filterer for contract events
}

// IInflationControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type IInflationControllerCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// IInflationControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IInflationControllerTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// IInflationControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IInflationControllerFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// CalculateSupplyDelta is a free data retrieval call binding the contract method 0x92eff3cd.
//
// Solidity: function calculateSupplyDelta(uint256 _currentSupply, uint256 _inflationReserve, uint256 _lastEpochTime, uint256 _currentEpochTime) view returns(uint256)
func (_IInflationController *IInflationControllerCaller) CalculateSupplyDelta(opts *runOptions, _currentSupply *big.Int, _inflationReserve *big.Int, _lastEpochTime *big.Int, _currentEpochTime *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _IInflationController.call(opts, "calculateSupplyDelta", _currentSupply, _inflationReserve, _lastEpochTime, _currentEpochTime)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

/* EVENTS ARE NOT YET SUPPORTED

 */

// IOracleMetaData contains all meta data concerning the IOracle contract.
var IOracleMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"}],\"name\":\"NewSymbols\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"_votes\",\"type\":\"int256[]\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPrecision\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"getRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSymbols\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVotePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"latestRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"name\":\"setSymbols\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_newVoters\",\"type\":\"address[]\"}],\"name\":\"setVoters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_commit\",\"type\":\"uint256\"},{\"internalType\":\"int256[]\",\"name\":\"_reports\",\"type\":\"int256[]\"},{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"4bb278f3": "finalize()",
		"9670c0bc": "getPrecision()",
		"9f8743f7": "getRound()",
		"3c8510fd": "getRoundData(uint256,string)",
		"df7f710e": "getSymbols()",
		"b78dec52": "getVotePeriod()",
		"cdd72253": "getVoters()",
		"33f98c77": "latestRoundData(string)",
		"b3ab15fb": "setOperator(address)",
		"8d4f75d2": "setSymbols(string[])",
		"845023f2": "setVoters(address[])",
		"307de9b6": "vote(uint256,int256[],uint256)",
	},
}

// IOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use IOracleMetaData.ABI instead.
var IOracleABI = IOracleMetaData.ABI

// Deprecated: Use IOracleMetaData.Sigs instead.
// IOracleFuncSigs maps the 4-byte function signature to its string representation.
var IOracleFuncSigs = IOracleMetaData.Sigs

// IOracle is an auto generated Go binding around an Ethereum contract.
type IOracle struct {
	IOracleCaller     // Read-only binding to the contract
	IOracleTransactor // Write-only binding to the contract
	IOracleFilterer   // Log filterer for contract events
}

// IOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type IOracleCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// IOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IOracleTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// IOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IOracleFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() view returns(uint256)
func (_IOracle *IOracleCaller) GetPrecision(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _IOracle.call(opts, "getPrecision")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_IOracle *IOracleCaller) GetRound(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _IOracle.call(opts, "getRound")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleCaller) GetRoundData(opts *runOptions, _round *big.Int, _symbol string) (IOracleRoundData, uint64, error) {
	out, consumed, err := _IOracle.call(opts, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)
	return out0, consumed, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[] _symbols)
func (_IOracle *IOracleCaller) GetSymbols(opts *runOptions) ([]string, uint64, error) {
	out, consumed, err := _IOracle.call(opts, "getSymbols")

	if err != nil {
		return *new([]string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, consumed, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_IOracle *IOracleCaller) GetVotePeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _IOracle.call(opts, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_IOracle *IOracleCaller) GetVoters(opts *runOptions) ([]common.Address, uint64, error) {
	out, consumed, err := _IOracle.call(opts, "getVoters")

	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleCaller) LatestRoundData(opts *runOptions, _symbol string) (IOracleRoundData, uint64, error) {
	out, consumed, err := _IOracle.call(opts, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)
	return out0, consumed, err

}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracleTransactor) Finalize(opts *runOptions) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "finalize")
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracleTransactor) SetOperator(opts *runOptions, _operator common.Address) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "setOperator", _operator)
	return consumed, err
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracleTransactor) SetSymbols(opts *runOptions, _symbols []string) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "setSymbols", _symbols)
	return consumed, err
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_IOracle *IOracleTransactor) SetVoters(opts *runOptions, _newVoters []common.Address) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "setVoters", _newVoters)
	return consumed, err
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_IOracle *IOracleTransactor) Vote(opts *runOptions, _commit *big.Int, _reports []*big.Int, _salt *big.Int) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "vote", _commit, _reports, _salt)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// IOracleNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the IOracle contract.
		type IOracleNewRoundIterator struct {
			Event *IOracleNewRound // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IOracleNewRoundIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IOracleNewRound)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IOracleNewRound)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IOracleNewRoundIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IOracleNewRoundIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IOracleNewRound represents a NewRound event raised by the IOracle contract.
		type IOracleNewRound struct {
			Round *big.Int;
			Height *big.Int;
			Timestamp *big.Int;
			VotePeriod *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewRound is a free log retrieval operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
		//
		// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
 		func (_IOracle *IOracleFilterer) FilterNewRound(opts *bind.FilterOpts) (*IOracleNewRoundIterator, error) {






			logs, sub, err := _IOracle.contract.FilterLogs(opts, "NewRound")
			if err != nil {
				return nil, err
			}
			return &IOracleNewRoundIterator{contract: _IOracle.contract, event: "NewRound", logs: logs, sub: sub}, nil
 		}

		// WatchNewRound is a free log subscription operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
		//
		// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
		func (_IOracle *IOracleFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *IOracleNewRound) (event.Subscription, error) {






			logs, sub, err := _IOracle.contract.WatchLogs(opts, "NewRound")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IOracleNewRound)
						if err := _IOracle.contract.UnpackLog(event, "NewRound", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewRound is a log parse operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
		//
		// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
		func (_IOracle *IOracleFilterer) ParseNewRound(log types.Log) (*IOracleNewRound, error) {
			event := new(IOracleNewRound)
			if err := _IOracle.contract.UnpackLog(event, "NewRound", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// IOracleNewSymbolsIterator is returned from FilterNewSymbols and is used to iterate over the raw logs and unpacked data for NewSymbols events raised by the IOracle contract.
		type IOracleNewSymbolsIterator struct {
			Event *IOracleNewSymbols // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IOracleNewSymbolsIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IOracleNewSymbols)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IOracleNewSymbols)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IOracleNewSymbolsIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IOracleNewSymbolsIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IOracleNewSymbols represents a NewSymbols event raised by the IOracle contract.
		type IOracleNewSymbols struct {
			Symbols []string;
			Round *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewSymbols is a free log retrieval operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
		//
		// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
 		func (_IOracle *IOracleFilterer) FilterNewSymbols(opts *bind.FilterOpts) (*IOracleNewSymbolsIterator, error) {




			logs, sub, err := _IOracle.contract.FilterLogs(opts, "NewSymbols")
			if err != nil {
				return nil, err
			}
			return &IOracleNewSymbolsIterator{contract: _IOracle.contract, event: "NewSymbols", logs: logs, sub: sub}, nil
 		}

		// WatchNewSymbols is a free log subscription operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
		//
		// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
		func (_IOracle *IOracleFilterer) WatchNewSymbols(opts *bind.WatchOpts, sink chan<- *IOracleNewSymbols) (event.Subscription, error) {




			logs, sub, err := _IOracle.contract.WatchLogs(opts, "NewSymbols")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IOracleNewSymbols)
						if err := _IOracle.contract.UnpackLog(event, "NewSymbols", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewSymbols is a log parse operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
		//
		// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
		func (_IOracle *IOracleFilterer) ParseNewSymbols(log types.Log) (*IOracleNewSymbols, error) {
			event := new(IOracleNewSymbols)
			if err := _IOracle.contract.UnpackLog(event, "NewSymbols", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// IOracleVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the IOracle contract.
		type IOracleVotedIterator struct {
			Event *IOracleVoted // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *IOracleVotedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(IOracleVoted)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(IOracleVoted)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *IOracleVotedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *IOracleVotedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// IOracleVoted represents a Voted event raised by the IOracle contract.
		type IOracleVoted struct {
			Voter common.Address;
			Votes []*big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterVoted is a free log retrieval operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
		//
		// Solidity: event Voted(address indexed _voter, int256[] _votes)
 		func (_IOracle *IOracleFilterer) FilterVoted(opts *bind.FilterOpts, _voter []common.Address) (*IOracleVotedIterator, error) {

			var _voterRule []interface{}
			for _, _voterItem := range _voter {
				_voterRule = append(_voterRule, _voterItem)
			}


			logs, sub, err := _IOracle.contract.FilterLogs(opts, "Voted", _voterRule)
			if err != nil {
				return nil, err
			}
			return &IOracleVotedIterator{contract: _IOracle.contract, event: "Voted", logs: logs, sub: sub}, nil
 		}

		// WatchVoted is a free log subscription operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
		//
		// Solidity: event Voted(address indexed _voter, int256[] _votes)
		func (_IOracle *IOracleFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *IOracleVoted, _voter []common.Address) (event.Subscription, error) {

			var _voterRule []interface{}
			for _, _voterItem := range _voter {
				_voterRule = append(_voterRule, _voterItem)
			}


			logs, sub, err := _IOracle.contract.WatchLogs(opts, "Voted", _voterRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(IOracleVoted)
						if err := _IOracle.contract.UnpackLog(event, "Voted", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseVoted is a log parse operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
		//
		// Solidity: event Voted(address indexed _voter, int256[] _votes)
		func (_IOracle *IOracleFilterer) ParseVoted(log types.Log) (*IOracleVoted, error) {
			event := new(IOracleVoted)
			if err := _IOracle.contract.UnpackLog(event, "Voted", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// IStabilizationMetaData contains all meta data concerning the IStabilization contract.
var IStabilizationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"b3ab15fb": "setOperator(address)",
		"7adbf973": "setOracle(address)",
	},
}

// IStabilizationABI is the input ABI used to generate the binding from.
// Deprecated: Use IStabilizationMetaData.ABI instead.
var IStabilizationABI = IStabilizationMetaData.ABI

// Deprecated: Use IStabilizationMetaData.Sigs instead.
// IStabilizationFuncSigs maps the 4-byte function signature to its string representation.
var IStabilizationFuncSigs = IStabilizationMetaData.Sigs

// IStabilization is an auto generated Go binding around an Ethereum contract.
type IStabilization struct {
	IStabilizationCaller     // Read-only binding to the contract
	IStabilizationTransactor // Write-only binding to the contract
	IStabilizationFilterer   // Log filterer for contract events
}

// IStabilizationCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type IStabilizationCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// IStabilizationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IStabilizationTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// IStabilizationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IStabilizationFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilizationTransactor) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _IStabilization.call(opts, "setOperator", operator)
	return consumed, err
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilizationTransactor) SetOracle(opts *runOptions, oracle common.Address) (uint64, error) {
	_, consumed, err := _IStabilization.call(opts, "setOracle", oracle)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

 */

// ISupplyControlMetaData contains all meta data concerning the ISupplyControl contract.
var ISupplyControlMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"availableSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stabilizer_\",\"type\":\"address\"}],\"name\":\"setStabilizer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stabilizer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7ecc2b56": "availableSupply()",
		"44df8e70": "burn()",
		"40c10f19": "mint(address,uint256)",
		"b3ab15fb": "setOperator(address)",
		"db7f521a": "setStabilizer(address)",
		"7e47961c": "stabilizer()",
		"18160ddd": "totalSupply()",
	},
}

// ISupplyControlABI is the input ABI used to generate the binding from.
// Deprecated: Use ISupplyControlMetaData.ABI instead.
var ISupplyControlABI = ISupplyControlMetaData.ABI

// Deprecated: Use ISupplyControlMetaData.Sigs instead.
// ISupplyControlFuncSigs maps the 4-byte function signature to its string representation.
var ISupplyControlFuncSigs = ISupplyControlMetaData.Sigs

// ISupplyControl is an auto generated Go binding around an Ethereum contract.
type ISupplyControl struct {
	ISupplyControlCaller     // Read-only binding to the contract
	ISupplyControlTransactor // Write-only binding to the contract
	ISupplyControlFilterer   // Log filterer for contract events
}

// ISupplyControlCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type ISupplyControlCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// ISupplyControlTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISupplyControlTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// ISupplyControlFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISupplyControlFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCaller) AvailableSupply(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _ISupplyControl.call(opts, "availableSupply")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Stabilizer is a free data retrieval call binding the contract method 0x7e47961c.
//
// Solidity: function stabilizer() view returns(address)
func (_ISupplyControl *ISupplyControlCaller) Stabilizer(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _ISupplyControl.call(opts, "stabilizer")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCaller) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _ISupplyControl.call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControlTransactor) Burn(opts *runOptions) (uint64, error) {
	_, consumed, err := _ISupplyControl.call(opts, "burn")
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControlTransactor) Mint(opts *runOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _ISupplyControl.call(opts, "mint", recipient, amount)
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControlTransactor) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _ISupplyControl.call(opts, "setOperator", operator)
	return consumed, err
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControlTransactor) SetStabilizer(opts *runOptions, stabilizer_ common.Address) (uint64, error) {
	_, consumed, err := _ISupplyControl.call(opts, "setStabilizer", stabilizer_)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// ISupplyControlBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the ISupplyControl contract.
		type ISupplyControlBurnIterator struct {
			Event *ISupplyControlBurn // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *ISupplyControlBurnIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(ISupplyControlBurn)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(ISupplyControlBurn)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *ISupplyControlBurnIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *ISupplyControlBurnIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// ISupplyControlBurn represents a Burn event raised by the ISupplyControl contract.
		type ISupplyControlBurn struct {
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterBurn is a free log retrieval operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
		//
		// Solidity: event Burn(uint256 amount)
 		func (_ISupplyControl *ISupplyControlFilterer) FilterBurn(opts *bind.FilterOpts) (*ISupplyControlBurnIterator, error) {



			logs, sub, err := _ISupplyControl.contract.FilterLogs(opts, "Burn")
			if err != nil {
				return nil, err
			}
			return &ISupplyControlBurnIterator{contract: _ISupplyControl.contract, event: "Burn", logs: logs, sub: sub}, nil
 		}

		// WatchBurn is a free log subscription operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
		//
		// Solidity: event Burn(uint256 amount)
		func (_ISupplyControl *ISupplyControlFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *ISupplyControlBurn) (event.Subscription, error) {



			logs, sub, err := _ISupplyControl.contract.WatchLogs(opts, "Burn")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(ISupplyControlBurn)
						if err := _ISupplyControl.contract.UnpackLog(event, "Burn", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseBurn is a log parse operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
		//
		// Solidity: event Burn(uint256 amount)
		func (_ISupplyControl *ISupplyControlFilterer) ParseBurn(log types.Log) (*ISupplyControlBurn, error) {
			event := new(ISupplyControlBurn)
			if err := _ISupplyControl.contract.UnpackLog(event, "Burn", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// ISupplyControlMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the ISupplyControl contract.
		type ISupplyControlMintIterator struct {
			Event *ISupplyControlMint // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *ISupplyControlMintIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(ISupplyControlMint)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(ISupplyControlMint)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *ISupplyControlMintIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *ISupplyControlMintIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// ISupplyControlMint represents a Mint event raised by the ISupplyControl contract.
		type ISupplyControlMint struct {
			Recipient common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterMint is a free log retrieval operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
		//
		// Solidity: event Mint(address recipient, uint256 amount)
 		func (_ISupplyControl *ISupplyControlFilterer) FilterMint(opts *bind.FilterOpts) (*ISupplyControlMintIterator, error) {




			logs, sub, err := _ISupplyControl.contract.FilterLogs(opts, "Mint")
			if err != nil {
				return nil, err
			}
			return &ISupplyControlMintIterator{contract: _ISupplyControl.contract, event: "Mint", logs: logs, sub: sub}, nil
 		}

		// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
		//
		// Solidity: event Mint(address recipient, uint256 amount)
		func (_ISupplyControl *ISupplyControlFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *ISupplyControlMint) (event.Subscription, error) {




			logs, sub, err := _ISupplyControl.contract.WatchLogs(opts, "Mint")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(ISupplyControlMint)
						if err := _ISupplyControl.contract.UnpackLog(event, "Mint", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseMint is a log parse operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
		//
		// Solidity: event Mint(address recipient, uint256 amount)
		func (_ISupplyControl *ISupplyControlFilterer) ParseMint(log types.Log) (*ISupplyControlMint, error) {
			event := new(ISupplyControlMint)
			if err := _ISupplyControl.contract.UnpackLog(event, "Mint", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// InflationControllerMetaData contains all meta data concerning the InflationController contract.
var InflationControllerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"SD59x18\",\"name\":\"iInit\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"iTrans\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"aE\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"T\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"iPerm\",\"type\":\"int256\"}],\"internalType\":\"structInflationController.Params\",\"name\":\"_params\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Convert_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Convert_Underflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Div_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Div_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Exp_InputTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PRBMath_SD59x18_Mul_InputTooSmall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"SD59x18\",\"name\":\"x\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"y\",\"type\":\"int256\"}],\"name\":\"PRBMath_SD59x18_Mul_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_currentSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_inflationReserve\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lastEpochTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_currentEpochTime\",\"type\":\"uint256\"}],\"name\":\"calculateSupplyDelta\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_currentSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_lastEpochTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_currentEpochTime\",\"type\":\"uint256\"}],\"name\":\"calculateSupplyDeltaOLD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"params\",\"outputs\":[{\"internalType\":\"SD59x18\",\"name\":\"iInit\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"iTrans\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"aE\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"T\",\"type\":\"int256\"},{\"internalType\":\"SD59x18\",\"name\":\"iPerm\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"92eff3cd": "calculateSupplyDelta(uint256,uint256,uint256,uint256)",
		"5ed8f8db": "calculateSupplyDeltaOLD(uint256,uint256,uint256)",
		"cff0ab96": "params()",
	},
	Bin: "0x608060405234801561000f575f80fd5b506040516111a03803806111a083398101604081905261002e91610056565b80515f55602081015160015560408101516002556060810151600355608001516004556100cf565b5f60a08284031215610066575f80fd5b60405160a081016001600160401b038111828210171561009457634e487b7160e01b5f52604160045260245ffd5b806040525082518152602083015160208201526040830151604082015260608301516060820152608083015160808201528091505092915050565b6110c4806100dc5f395ff3fe608060405234801561000f575f80fd5b506004361061003f575f3560e01c80635ed8f8db1461004357806392eff3cd14610069578063cff0ab961461007c575b5f80fd5b610056610051366004610f76565b6100bd565b6040519081526020015b60405180910390f35b610056610077366004610f9f565b6101f9565b5f54600154600254600354600454610095949392919085565b604080519586526020860194909452928401919091526060830152608082015260a001610060565b5f806100c884610298565b90505f6100d484610298565b90505f6100f06100e85f6002015485610321565b6003546103e6565b90505f6101046100e85f6002015485610321565b90505f6101bf61016d6101235f80015461011e888a61048c565b610321565b61016861014461013a5f8001545f6001015461048c565b61011e8a8c61048c565b6101636101545f6002015461049a565b61015e6001610298565b61048c565b6103e6565b6104ee565b6101686101a461018f6101875f600101545f80015461048c565b600354610321565b61011e61019b8861049a565b61015e8a61049a565b6101636101b76101545f6002015461049a565b600254610321565b90506101ea6101e56101dc6101d38461049a565b61011e8d610298565b61015e8c610298565b610502565b955050505050505b9392505050565b5f8061020484610298565b90505f61021084610298565b600354909150811361023057610227878383610515565b92505050610290565b60035482128015610242575060035481135b15610280575f61025788845f60030154610515565b90505f610269885f60030154856105d3565b90506102758183610fe2565b945050505050610290565b61028b8683836105d3565b925050505b949350505050565b5f6102af670de0b6b3a7640000600160ff1b611009565b8212156102d7576040516399474eeb60e01b8152600481018390526024015b60405180910390fd5b6102f0670de0b6b3a76400006001600160ff1b03611009565b82131561031357604051639d58109160e01b8152600481018390526024016102ce565b50670de0b6b3a76400000290565b5f8282600160ff1b8214806103395750600160ff1b81145b156103575760405163a6070c2560e01b815260040160405180910390fd5b5f805f8412610366578361036a565b835f035b91505f8312610379578261037d565b825f035b90505f61038a83836105ee565b90506001600160ff1b038111156103be5760405163120b5b4360e01b815260048101899052602481018890526044016102ce565b5f19858518136103d9816103d457825f036103d6565b825b90565b9998505050505050505050565b5f8282600160ff1b8214806103fe5750600160ff1b81145b1561041c576040516309fe2b4560e41b815260040160405180910390fd5b5f805f841261042b578361042f565b835f035b91505f831261043e5782610442565b825f035b90505f61045883670de0b6b3a7640000846106a0565b90506001600160ff1b038111156103be5760405163d49c26b360e01b815260048101899052602481018890526044016102ce565b5f6101f26103d68385611041565b5f81680736ea4425c11ac6308113156104c95760405163ca7ec0c560e01b8152600481018490526024016102ce565b6714057b7ef767814f81026102906104e9670de0b6b3a7640000835b0590565b61076f565b5f6101f26103d68385611067565b92915050565b5f6104fc670de0b6b3a764000083611009565b5f8061052c5f600201546105285f610298565b1490565b1561055a575f546001546105539190610168906100e89061054d908461048c565b88610321565b90506105a3565b5f61056c6100e85f6002015487610321565b90505f61057e6101446101548461049a565b5f5460015491925061059e9161016890610598908361048c565b84610321565b925050505b5f6105be6105b48361011e89610298565b61011e868861048c565b90506105c981610502565b9695505050505050565b5f6102906101e56105e66105b487610298565b600454610321565b5f80805f19848609848602925082811083820303915050805f0361061f5750670de0b6b3a7640000900490506104fc565b670de0b6b3a7640000811061065157604051635173648d60e01b815260048101869052602481018590526044016102ce565b5f670de0b6b3a764000085870962040000818503049310909103600160ee1b02919091177faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106690291505092915050565b5f80805f19858709858702925082811083820303915050805f036106d7578382816106cd576106cd610ff5565b04925050506101f2565b83811061070857604051630c740aef60e31b81526004810187905260248101869052604481018590526064016102ce565b5f8486880960026001871981018816978890046003810283188082028403028082028403028082028403028082028403028082028403029081029092039091025f889003889004909101858311909403939093029303949094049190911702949350505050565b5f81818112156107c95768033dd1780914b971141981121561079357505f92915050565b6107c26107a56103d66104e9845f0390565b6ec097ce7bc90715b34b9f1000000000816104e5576104e5610ff5565b9150610810565b680a688906bd8affffff8113156107f557604051626c1a0560e31b8152600481018490526024016102ce565b670de0b6b3a7640000604082901b056102906103d682610816565b50919050565b600160bf1b67ff000000000000008216156109235767800000000000000082161561084a5768016a09e667f3bcc9090260401c5b674000000000000000821615610869576801306fe0a31b7152df0260401c5b672000000000000000821615610888576801172b83c7d517adce0260401c5b6710000000000000008216156108a75768010b5586cf9890f62a0260401c5b6708000000000000008216156108c6576801059b0d31585743ae0260401c5b6704000000000000008216156108e557680102c9a3e778060ee70260401c5b6702000000000000008216156109045768010163da9fb33356d80260401c5b67010000000000000082161561092357680100b1afa5abcbed610260401c5b66ff000000000000821615610a225766800000000000008216156109505768010058c86da1c09ea20260401c5b664000000000000082161561096e576801002c605e2e8cec500260401c5b662000000000000082161561098c57680100162f3904051fa10260401c5b66100000000000008216156109aa576801000b175effdc76ba0260401c5b66080000000000008216156109c857680100058ba01fb9f96d0260401c5b66040000000000008216156109e65768010002c5cc37da94920260401c5b6602000000000000821615610a04576801000162e525ee05470260401c5b6601000000000000821615610a225768010000b17255775c040260401c5b65ff0000000000821615610b185765800000000000821615610a4d576801000058b91b5bc9ae0260401c5b65400000000000821615610a6a57680100002c5c89d5ec6d0260401c5b65200000000000821615610a875768010000162e43f4f8310260401c5b65100000000000821615610aa457680100000b1721bcfc9a0260401c5b65080000000000821615610ac15768010000058b90cf1e6e0260401c5b65040000000000821615610ade576801000002c5c863b73f0260401c5b65020000000000821615610afb57680100000162e430e5a20260401c5b65010000000000821615610b18576801000000b1721835510260401c5b64ff00000000821615610c0557648000000000821615610b4157680100000058b90c0b490260401c5b644000000000821615610b5d5768010000002c5c8601cc0260401c5b642000000000821615610b79576801000000162e42fff00260401c5b641000000000821615610b955768010000000b17217fbb0260401c5b640800000000821615610bb1576801000000058b90bfce0260401c5b640400000000821615610bcd57680100000002c5c85fe30260401c5b640200000000821615610be95768010000000162e42ff10260401c5b640100000000821615610c0557680100000000b17217f80260401c5b63ff000000821615610ce9576380000000821615610c2c5768010000000058b90bfc0260401c5b6340000000821615610c47576801000000002c5c85fe0260401c5b6320000000821615610c6257680100000000162e42ff0260401c5b6310000000821615610c7d576801000000000b17217f0260401c5b6308000000821615610c9857680100000000058b90c00260401c5b6304000000821615610cb35768010000000002c5c8600260401c5b6302000000821615610cce576801000000000162e4300260401c5b6301000000821615610ce95768010000000000b172180260401c5b62ff0000821615610dc45762800000821615610d0e576801000000000058b90c0260401c5b62400000821615610d2857680100000000002c5c860260401c5b62200000821615610d425768010000000000162e430260401c5b62100000821615610d5c57680100000000000b17210260401c5b62080000821615610d765768010000000000058b910260401c5b62040000821615610d90576801000000000002c5c80260401c5b62020000821615610daa57680100000000000162e40260401c5b62010000821615610dc4576801000000000000b1720260401c5b61ff00821615610e9657618000821615610de757680100000000000058b90260401c5b614000821615610e005768010000000000002c5d0260401c5b612000821615610e19576801000000000000162e0260401c5b611000821615610e325768010000000000000b170260401c5b610800821615610e4b576801000000000000058c0260401c5b610400821615610e6457680100000000000002c60260401c5b610200821615610e7d57680100000000000001630260401c5b610100821615610e9657680100000000000000b10260401c5b60ff821615610f5f576080821615610eb757680100000000000000590260401c5b6040821615610ecf576801000000000000002c0260401c5b6020821615610ee757680100000000000000160260401c5b6010821615610eff576801000000000000000b0260401c5b6008821615610f1757680100000000000000060260401c5b6004821615610f2f57680100000000000000030260401c5b6002821615610f4757680100000000000000010260401c5b6001821615610f5f57680100000000000000010260401c5b670de0b6b3a76400000260409190911c60bf031c90565b5f805f60608486031215610f88575f80fd5b505081359360208301359350604090920135919050565b5f805f8060808587031215610fb2575f80fd5b5050823594602084013594506040840135936060013592509050565b634e487b7160e01b5f52601160045260245ffd5b808201808211156104fc576104fc610fce565b634e487b7160e01b5f52601260045260245ffd5b5f8261102357634e487b7160e01b5f52601260045260245ffd5b600160ff1b82145f198414161561103c5761103c610fce565b500590565b8181035f83128015838313168383128216171561106057611060610fce565b5092915050565b8082018281125f83128015821682158216171561108657611086610fce565b50509291505056fea2646970667358221220b31e4c3dc45d7f25913eb8f235580549b239cb91edbef74700dab5f4fe45f79964736f6c63430008150033",
}

// InflationControllerABI is the input ABI used to generate the binding from.
// Deprecated: Use InflationControllerMetaData.ABI instead.
var InflationControllerABI = InflationControllerMetaData.ABI

// Deprecated: Use InflationControllerMetaData.Sigs instead.
// InflationControllerFuncSigs maps the 4-byte function signature to its string representation.
var InflationControllerFuncSigs = InflationControllerMetaData.Sigs

// InflationControllerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use InflationControllerMetaData.Bin instead.
var InflationControllerBin = InflationControllerMetaData.Bin

// DeployInflationController deploys a new Ethereum contract, binding an instance of InflationController to it.
func (r *runner) deployInflationController(opts *runOptions, _params InflationControllerParams) (common.Address, uint64, *InflationController, error) {
	parsed, err := InflationControllerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(InflationControllerBin), _params)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &InflationController{InflationControllerCaller: InflationControllerCaller{contract: c}, InflationControllerTransactor: InflationControllerTransactor{contract: c}, InflationControllerFilterer: InflationControllerFilterer{contract: c}}, nil
}

// InflationController is an auto generated Go binding around an Ethereum contract.
type InflationController struct {
	InflationControllerCaller     // Read-only binding to the contract
	InflationControllerTransactor // Write-only binding to the contract
	InflationControllerFilterer   // Log filterer for contract events
}

// InflationControllerCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type InflationControllerCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// InflationControllerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InflationControllerTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// InflationControllerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InflationControllerFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// CalculateSupplyDelta is a free data retrieval call binding the contract method 0x92eff3cd.
//
// Solidity: function calculateSupplyDelta(uint256 _currentSupply, uint256 _inflationReserve, uint256 _lastEpochTime, uint256 _currentEpochTime) view returns(uint256)
func (_InflationController *InflationControllerCaller) CalculateSupplyDelta(opts *runOptions, _currentSupply *big.Int, _inflationReserve *big.Int, _lastEpochTime *big.Int, _currentEpochTime *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _InflationController.call(opts, "calculateSupplyDelta", _currentSupply, _inflationReserve, _lastEpochTime, _currentEpochTime)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// CalculateSupplyDeltaOLD is a free data retrieval call binding the contract method 0x5ed8f8db.
//
// Solidity: function calculateSupplyDeltaOLD(uint256 _currentSupply, uint256 _lastEpochTime, uint256 _currentEpochTime) view returns(uint256)
func (_InflationController *InflationControllerCaller) CalculateSupplyDeltaOLD(opts *runOptions, _currentSupply *big.Int, _lastEpochTime *big.Int, _currentEpochTime *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _InflationController.call(opts, "calculateSupplyDeltaOLD", _currentSupply, _lastEpochTime, _currentEpochTime)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns(int256 iInit, int256 iTrans, int256 aE, int256 T, int256 iPerm)
func (_InflationController *InflationControllerCaller) Params(opts *runOptions) (struct {
	IInit  *big.Int
	ITrans *big.Int
	AE     *big.Int
	T      *big.Int
	IPerm  *big.Int
}, uint64, error) {
	out, consumed, err := _InflationController.call(opts, "params")

	outstruct := new(struct {
		IInit  *big.Int
		ITrans *big.Int
		AE     *big.Int
		T      *big.Int
		IPerm  *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.IInit = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ITrans = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.AE = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.T = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.IPerm = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

/* EVENTS ARE NOT YET SUPPORTED

 */

// LiquidMetaData contains all meta data concerning the Liquid contract.
var LiquidMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"_treasury\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_index\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_FACTOR_UNIT_RECIP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"lockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_ntnReward\",\"type\":\"uint256\"}],\"name\":\"redistribute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unclaimedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_unclaimedATN\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_unclaimedNTN\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"unlockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"187cf4d7": "FEE_FACTOR_UNIT_RECIP()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"9dc29fac": "burn(address,uint256)",
		"372500ab": "claimRewards()",
		"5ea1d6f8": "commissionRate()",
		"313ce567": "decimals()",
		"282d3fdf": "lock(address,uint256)",
		"59355736": "lockedBalanceOf(address)",
		"40c10f19": "mint(address,uint256)",
		"06fdde03": "name()",
		"a0ce552d": "redistribute(uint256)",
		"19fac8fd": "setCommissionRate(uint256)",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"61d027b3": "treasury()",
		"949813b8": "unclaimedRewards(address)",
		"7eee288d": "unlock(address,uint256)",
		"84955c88": "unlockedBalanceOf(address)",
		"3a5381b5": "validator()",
	},
	Bin: "0x608060405234801562000010575f80fd5b50604051620017d6380380620017d6833981016040819052620000339162000149565b61271082111562000042575f80fd5b600d80546001600160a01b038087166001600160a01b031992831617909255600e805492861692909116919091179055600f8290556040516200008a90829060200162000230565b604051602081830303815290604052600b9081620000a99190620002ea565b5080604051602001620000bd919062000230565b604051602081830303815290604052600c9081620000dc9190620002ea565b50505f80546001600160a01b0319163317905550620003b2915050565b6001600160a01b03811681146200010e575f80fd5b50565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156200014157818101518382015260200162000127565b50505f910152565b5f805f80608085870312156200015d575f80fd5b84516200016a81620000f9565b60208601519094506200017d81620000f9565b6040860151606087015191945092506001600160401b0380821115620001a1575f80fd5b818701915087601f830112620001b5575f80fd5b815181811115620001ca57620001ca62000111565b604051601f8201601f19908116603f01168101908382118183101715620001f557620001f562000111565b816040528281528a60208487010111156200020e575f80fd5b6200022183602083016020880162000125565b979a9699509497505050505050565b644c4e544e2d60d81b81525f82516200025181600585016020870162000125565b9190910160050192915050565b600181811c908216806200027357607f821691505b6020821081036200029257634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620002e5575f81815260208120601f850160051c81016020861015620002c05750805b601f850160051c820191505b81811015620002e157828155600101620002cc565b5050505b505050565b81516001600160401b0381111562000306576200030662000111565b6200031e816200031784546200025e565b8462000298565b602080601f83116001811462000354575f84156200033c5750858301515b5f19600386901b1c1916600185901b178555620002e1565b5f85815260208120601f198616915b82811015620003845788860151825594840194600190910190840162000363565b5085821015620003a257878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61141680620003c05f395ff3fe608060405260043610610147575f3560e01c806359355736116100b3578063949813b81161006d578063949813b8146103ac57806395d89b41146103e05780639dc29fac146103f4578063a0ce552d14610413578063a9059cbb14610426578063dd62ed3e14610445575f80fd5b806359355736146102d25780635ea1d6f81461030657806361d027b31461031b57806370a082311461033a5780637eee288d1461036e57806384955c881461038d575f80fd5b8063282d3fdf11610104578063282d3fdf146102195780632f2c3f2e14610238578063313ce5671461024d578063372500ab146102685780633a5381b51461027c57806340c10f19146102b3575f80fd5b806306fdde031461014b578063095ea7b31461017557806318160ddd146101a4578063187cf4d7146101c257806319fac8fd146101d957806323b872dd146101fa575b5f80fd5b348015610156575f80fd5b5061015f610489565b60405161016c9190611182565b60405180910390f35b348015610180575f80fd5b5061019461018f3660046111e8565b610519565b604051901515815260200161016c565b3480156101af575f80fd5b506004545b60405190815260200161016c565b3480156101cd575f80fd5b506101b4633b9aca0081565b3480156101e4575f80fd5b506101f86101f3366004611210565b61052f565b005b348015610205575f80fd5b50610194610214366004611227565b610566565b348015610224575f80fd5b506101f86102333660046111e8565b610657565b348015610243575f80fd5b506101b461271081565b348015610258575f80fd5b506040516012815260200161016c565b348015610273575f80fd5b506101f8610739565b348015610287575f80fd5b50600d5461029b906001600160a01b031681565b6040516001600160a01b03909116815260200161016c565b3480156102be575f80fd5b506101f86102cd3660046111e8565b6108a9565b3480156102dd575f80fd5b506101b46102ec366004611260565b6001600160a01b03165f9081526002602052604090205490565b348015610311575f80fd5b506101b4600f5481565b348015610326575f80fd5b50600e5461029b906001600160a01b031681565b348015610345575f80fd5b506101b4610354366004611260565b6001600160a01b03165f9081526001602052604090205490565b348015610379575f80fd5b506101f86103883660046111e8565b61090e565b348015610398575f80fd5b506101b46103a7366004611260565b6109d1565b3480156103b7575f80fd5b506103cb6103c6366004611260565b6109fe565b6040805192835260208301919091520161016c565b3480156103eb575f80fd5b5061015f610a62565b3480156103ff575f80fd5b506101f861040e3660046111e8565b610a71565b6103cb610421366004611210565b610ace565b348015610431575f80fd5b506101946104403660046111e8565b610d78565b348015610450575f80fd5b506101b461045f366004611280565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b6060600b8054610498906112b1565b80601f01602080910402602001604051908101604052809291908181526020018280546104c4906112b1565b801561050f5780601f106104e65761010080835404028352916020019161050f565b820191905f5260205f20905b8154815290600101906020018083116104f257829003601f168201915b5050505050905090565b5f610525338484610dc3565b5060015b92915050565b5f546001600160a01b031633146105615760405162461bcd60e51b8152600401610558906112e9565b60405180910390fd5b600f55565b6001600160a01b0383165f908152600360209081526040808320338452909152812054828110156105ea5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b6064820152608401610558565b6105fe85336105f98685611345565b610dc3565b6106088584610ee6565b6106128484610fd8565b836001600160a01b0316856001600160a01b03165f805160206113c18339815191528560405161064491815260200190565b60405180910390a3506001949350505050565b5f546001600160a01b031633146106805760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f9081526002602090815260408083205460019092529091205482916106af91611345565b10156107095760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611358565b90915550505050565b5f8061074433611022565b335f818152600560209081526040808320839055600890915280822082905551939550919350909184908381818185875af1925050503d805f81146107a4576040519150601f19603f3d011682016040523d82523d5f602084013e6107a9565b606091505b50509050806107ef5760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b7321020aa2760711b6044820152606401610558565b5f5460405163a9059cbb60e01b8152336004820152602481018490526001600160a01b039091169063a9059cbb906044016020604051808303815f875af115801561083c573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610860919061136b565b9050806108a45760405162461bcd60e51b81526020600482015260126024820152712330b4b632b2103a379039b2b73210272a2760711b6044820152606401610558565b505050565b5f546001600160a01b031633146108d25760405162461bcd60e51b8152600401610558906112e9565b6108dc8282610fd8565b6040518181526001600160a01b038316905f905f805160206113c1833981519152906020015b60405180910390a35050565b5f546001600160a01b031633146109375760405162461bcd60e51b8152600401610558906112e9565b6001600160a01b0382165f908152600260205260409020548111156109aa5760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b6064820152608401610558565b6001600160a01b0382165f9081526002602052604081208054839290610730908490611345565b6001600160a01b0381165f9081526002602090815260408083205460019092528220546105299190611345565b5f805f80610a0b856110c7565b6001600160a01b0387165f908152600560205260409020549193509150610a33908390611358565b6001600160a01b0386165f90815260086020526040902054909450610a59908290611358565b92505050915091565b6060600c8054610498906112b1565b5f546001600160a01b03163314610a9a5760405162461bcd60e51b8152600401610558906112e9565b610aa48282610ee6565b6040518181525f906001600160a01b038416905f805160206113c183398151915290602001610902565b5f805481906001600160a01b03163314610afa5760405162461bcd60e51b8152600401610558906112e9565b600f5434905f9061271090610b0f908461138a565b610b1991906113a1565b905081811115610b665760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610b708183611345565b600e546040519193506001600160a01b0316906108fc9083905f818181858888f193505050503d805f8114610bc0576040519150601f19603f3d011682016040523d82523d5f602084013e610bc5565b606091505b5050505f612710600f5487610bda919061138a565b610be491906113a1565b905085811115610c315760405162461bcd60e51b81526020600482015260186024820152771a5b9d985b1a59081d985b1a59185d1bdc881c995dd85c9960421b6044820152606401610558565b610c3b8187611345565b5f54600e5460405163a9059cbb60e01b81526001600160a01b03918216600482015260248101859052929850169063a9059cbb906044016020604051808303815f875af1158015610c8e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610cb2919061136b565b506004545f90610cc6633b9aca008661138a565b610cd091906113a1565b905080600754610ce09190611358565b6007556004545f90610cf6633b9aca008a61138a565b610d0091906113a1565b905080600a54610d109190611358565b600a556004545f90633b9aca0090610d28908561138a565b610d3291906113a1565b90505f633b9aca0060045484610d48919061138a565b610d5291906113a1565b9050610d5e8287611358565b610d688287611358565b9850985050505050505050915091565b5f610d833383610ee6565b610d8d8383610fd8565b6040518281526001600160a01b0384169033905f805160206113c18339815191529060200160405180910390a350600192915050565b6001600160a01b038316610e255760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b6064820152608401610558565b6001600160a01b038216610e865760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b6064820152608401610558565b6001600160a01b038381165f8181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610eef82611022565b50506001600160a01b0382165f90815260016020908152604080832054600290925290912054610f1f9082611345565b821115610f6e5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e647300000000006044820152606401610558565b610f788282611345565b6001600160a01b0384165f90815260016020526040902055808203610fbd576001600160a01b0383165f90815260066020908152604080832083905560099091528120555b8160045f828254610fce9190611345565b9091555050505050565b610fe182611022565b50506001600160a01b0382165f908152600160205260408120805483929061100a908490611358565b925050819055508060045f8282546107309190611358565b5f805f8061102f856110c7565b6001600160a01b0387165f908152600560205260409020549193509150611057908390611358565b6001600160a01b0386165f90815260056020908152604080832084905560075460068352818420556008909152902054909450611095908290611358565b6001600160a01b039095165f908152600860209081526040808320889055600a54600990925290912055509193915050565b6001600160a01b0381165f9081526001602052604081205481908082036110f357505f93849350915050565b6001600160a01b0384165f908152600660205260408120546007546111189190611345565b6001600160a01b0386165f90815260096020526040812054600a5492935090916111429190611345565b9050633b9aca00611153848461138a565b61115d91906113a1565b9450633b9aca0061116e848361138a565b61117891906113a1565b9350505050915091565b5f6020808352835180828501525f5b818110156111ad57858101830151858201604001528201611191565b505f604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b03811681146111e3575f80fd5b919050565b5f80604083850312156111f9575f80fd5b611202836111cd565b946020939093013593505050565b5f60208284031215611220575f80fd5b5035919050565b5f805f60608486031215611239575f80fd5b611242846111cd565b9250611250602085016111cd565b9150604084013590509250925092565b5f60208284031215611270575f80fd5b611279826111cd565b9392505050565b5f8060408385031215611291575f80fd5b61129a836111cd565b91506112a8602084016111cd565b90509250929050565b600181811c908216806112c557607f821691505b6020821081036112e357634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561052957610529611331565b8082018082111561052957610529611331565b5f6020828403121561137b575f80fd5b81518015158114611279575f80fd5b808202811582820484141761052957610529611331565b5f826113bb57634e487b7160e01b5f52601260045260245ffd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220fcf7fe14f5b055004178aa5c27d7983104c707dfee4e3b218ce03dfc1bc4bea864736f6c63430008150033",
}

// LiquidABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidMetaData.ABI instead.
var LiquidABI = LiquidMetaData.ABI

// Deprecated: Use LiquidMetaData.Sigs instead.
// LiquidFuncSigs maps the 4-byte function signature to its string representation.
var LiquidFuncSigs = LiquidMetaData.Sigs

// LiquidBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LiquidMetaData.Bin instead.
var LiquidBin = LiquidMetaData.Bin

// DeployLiquid deploys a new Ethereum contract, binding an instance of Liquid to it.
func (r *runner) deployLiquid(opts *runOptions, _validator common.Address, _treasury common.Address, _commissionRate *big.Int, _index string) (common.Address, uint64, *Liquid, error) {
	parsed, err := LiquidMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(LiquidBin), _validator, _treasury, _commissionRate, _index)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &Liquid{LiquidCaller: LiquidCaller{contract: c}, LiquidTransactor: LiquidTransactor{contract: c}, LiquidFilterer: LiquidFilterer{contract: c}}, nil
}

// Liquid is an auto generated Go binding around an Ethereum contract.
type Liquid struct {
	LiquidCaller     // Read-only binding to the contract
	LiquidTransactor // Write-only binding to the contract
	LiquidFilterer   // Log filterer for contract events
}

// LiquidCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type LiquidCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// LiquidTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// LiquidFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Liquid *LiquidCaller) COMMISSIONRATEPRECISION(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "COMMISSION_RATE_PRECISION")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// FEEFACTORUNITRECIP is a free data retrieval call binding the contract method 0x187cf4d7.
//
// Solidity: function FEE_FACTOR_UNIT_RECIP() view returns(uint256)
func (_Liquid *LiquidCaller) FEEFACTORUNITRECIP(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "FEE_FACTOR_UNIT_RECIP")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_Liquid *LiquidCaller) Allowance(opts *runOptions, _owner common.Address, _spender common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_Liquid *LiquidCaller) BalanceOf(opts *runOptions, _delegator common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "balanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_Liquid *LiquidCaller) CommissionRate(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "commissionRate")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Liquid *LiquidCaller) Decimals(opts *runOptions) (uint8, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "decimals")

	if err != nil {
		return *new(uint8), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, consumed, err

}

// LockedBalanceOf is a free data retrieval call binding the contract method 0x59355736.
//
// Solidity: function lockedBalanceOf(address _delegator) view returns(uint256)
func (_Liquid *LiquidCaller) LockedBalanceOf(opts *runOptions, _delegator common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "lockedBalanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Liquid *LiquidCaller) Name(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "name")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Liquid *LiquidCaller) Symbol(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "symbol")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Liquid *LiquidCaller) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Liquid *LiquidCaller) Treasury(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "treasury")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256 _unclaimedATN, uint256 _unclaimedNTN)
func (_Liquid *LiquidCaller) UnclaimedRewards(opts *runOptions, _account common.Address) (struct {
	UnclaimedATN *big.Int
	UnclaimedNTN *big.Int
}, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "unclaimedRewards", _account)

	outstruct := new(struct {
		UnclaimedATN *big.Int
		UnclaimedNTN *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.UnclaimedATN = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.UnclaimedNTN = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_Liquid *LiquidCaller) UnlockedBalanceOf(opts *runOptions, _delegator common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "unlockedBalanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_Liquid *LiquidCaller) Validator(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "validator")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_Liquid *LiquidTransactor) Approve(opts *runOptions, _spender common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "approve", _spender, _amount)
	return consumed, err
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_Liquid *LiquidTransactor) Burn(opts *runOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "burn", _account, _amount)
	return consumed, err
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_Liquid *LiquidTransactor) ClaimRewards(opts *runOptions) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "claimRewards")
	return consumed, err
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_Liquid *LiquidTransactor) Lock(opts *runOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "lock", _account, _amount)
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_Liquid *LiquidTransactor) Mint(opts *runOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "mint", _account, _amount)
	return consumed, err
}

// Redistribute is a paid mutator transaction binding the contract method 0xa0ce552d.
//
// Solidity: function redistribute(uint256 _ntnReward) payable returns(uint256, uint256)
func (_Liquid *LiquidTransactor) Redistribute(opts *runOptions, _ntnReward *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "redistribute", _ntnReward)
	return consumed, err
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_Liquid *LiquidTransactor) SetCommissionRate(opts *runOptions, _rate *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "setCommissionRate", _rate)
	return consumed, err
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_Liquid *LiquidTransactor) Transfer(opts *runOptions, _to common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "transfer", _to, _amount)
	return consumed, err
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_Liquid *LiquidTransactor) TransferFrom(opts *runOptions, _sender common.Address, _recipient common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "transferFrom", _sender, _recipient, _amount)
	return consumed, err
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_Liquid *LiquidTransactor) Unlock(opts *runOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "unlock", _account, _amount)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// LiquidApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Liquid contract.
		type LiquidApprovalIterator struct {
			Event *LiquidApproval // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *LiquidApprovalIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(LiquidApproval)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(LiquidApproval)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *LiquidApprovalIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *LiquidApprovalIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// LiquidApproval represents a Approval event raised by the Liquid contract.
		type LiquidApproval struct {
			Owner common.Address;
			Spender common.Address;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
 		func (_Liquid *LiquidFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LiquidApprovalIterator, error) {

			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var spenderRule []interface{}
			for _, spenderItem := range spender {
				spenderRule = append(spenderRule, spenderItem)
			}


			logs, sub, err := _Liquid.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
			if err != nil {
				return nil, err
			}
			return &LiquidApprovalIterator{contract: _Liquid.contract, event: "Approval", logs: logs, sub: sub}, nil
 		}

		// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
		func (_Liquid *LiquidFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LiquidApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var spenderRule []interface{}
			for _, spenderItem := range spender {
				spenderRule = append(spenderRule, spenderItem)
			}


			logs, sub, err := _Liquid.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(LiquidApproval)
						if err := _Liquid.contract.UnpackLog(event, "Approval", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
		func (_Liquid *LiquidFilterer) ParseApproval(log types.Log) (*LiquidApproval, error) {
			event := new(LiquidApproval)
			if err := _Liquid.contract.UnpackLog(event, "Approval", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// LiquidTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Liquid contract.
		type LiquidTransferIterator struct {
			Event *LiquidTransfer // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *LiquidTransferIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(LiquidTransfer)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(LiquidTransfer)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *LiquidTransferIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *LiquidTransferIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// LiquidTransfer represents a Transfer event raised by the Liquid contract.
		type LiquidTransfer struct {
			From common.Address;
			To common.Address;
			Value *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
 		func (_Liquid *LiquidFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidTransferIterator, error) {

			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}


			logs, sub, err := _Liquid.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
			if err != nil {
				return nil, err
			}
			return &LiquidTransferIterator{contract: _Liquid.contract, event: "Transfer", logs: logs, sub: sub}, nil
 		}

		// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
		func (_Liquid *LiquidFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LiquidTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}


			logs, sub, err := _Liquid.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(LiquidTransfer)
						if err := _Liquid.contract.UnpackLog(event, "Transfer", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
		func (_Liquid *LiquidFilterer) ParseTransfer(log types.Log) (*LiquidTransfer, error) {
			event := new(LiquidTransfer)
			if err := _Liquid.contract.UnpackLog(event, "Transfer", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// OracleMetaData contains all meta data concerning the Oracle contract.
var OracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_voters\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"_autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_votePeriod\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"}],\"name\":\"NewSymbols\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256[]\",\"name\":\"_votes\",\"type\":\"int256[]\"}],\"name\":\"Voted\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPrecision\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_round\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"getRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSymbols\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVotePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVoters\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRoundBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastVoterUpdateRound\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"}],\"name\":\"latestRoundData\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"price\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"}],\"internalType\":\"structIOracle.RoundData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newSymbols\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"reports\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"round\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_symbols\",\"type\":\"string[]\"}],\"name\":\"setSymbols\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_newVoters\",\"type\":\"address[]\"}],\"name\":\"setVoters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbolUpdatedRound\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"symbols\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_commit\",\"type\":\"uint256\"},{\"internalType\":\"int256[]\",\"name\":\"_reports\",\"type\":\"int256[]\"},{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"votePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"votingInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"commit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isVoter\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"4bb278f3": "finalize()",
		"9670c0bc": "getPrecision()",
		"9f8743f7": "getRound()",
		"3c8510fd": "getRoundData(uint256,string)",
		"df7f710e": "getSymbols()",
		"b78dec52": "getVotePeriod()",
		"cdd72253": "getVoters()",
		"e6a02a28": "lastRoundBlock()",
		"aa2f89b5": "lastVoterUpdateRound()",
		"33f98c77": "latestRoundData(string)",
		"5281b5c6": "newSymbols(uint256)",
		"4c56ea56": "reports(string,address)",
		"146ca531": "round()",
		"b3ab15fb": "setOperator(address)",
		"8d4f75d2": "setSymbols(string[])",
		"845023f2": "setVoters(address[])",
		"08f21ff5": "symbolUpdatedRound()",
		"ccce413b": "symbols(uint256)",
		"307de9b6": "vote(uint256,int256[],uint256)",
		"a7813587": "votePeriod()",
		"5412b3ae": "votingInfo(address)",
	},
	Bin: "0x6080604052600160ff1b600755600160ff1b60085534801562000020575f80fd5b5060405162002cd338038062002cd383398101604081905262000043916200061d565b600280546001600160a01b038087166001600160a01b0319928316179092556003805492861692909116919091179055815162000087905f90602085019062000358565b5081516200009d90600190602085019062000358565b5080600981905550620000c1855f60018851620000bb919062000728565b6200017b565b8451620000d6906004906020880190620003b3565b508451620000ec906005906020880190620003b3565b5060016006819055600d805490910181555f9081525b85518110156200016f576001600b5f88848151811062000126576200012662000744565b6020908102919091018101516001600160a01b031682528101919091526040015f20600201805460ff191691151591909117905580620001668162000758565b91505062000102565b5050505050506200098a565b8082126200018857505050565b81815f8560026200019a858562000773565b620001a691906200079c565b620001b29087620007d8565b81518110620001c557620001c562000744565b602002602001015190505b81831362000324575b806001600160a01b0316868481518110620001f857620001f862000744565b60200260200101516001600160a01b031610156200022557826200021c8162000802565b935050620001d9565b806001600160a01b031686838151811062000244576200024462000744565b60200260200101516001600160a01b0316111562000271578162000268816200081c565b92505062000225565b8183136200031e578582815181106200028e576200028e62000744565b6020026020010151868481518110620002ab57620002ab62000744565b6020026020010151878581518110620002c857620002c862000744565b60200260200101888581518110620002e457620002e462000744565b6001600160a01b03938416602091820292909201015291169052826200030a8162000802565b93505081806200031a906200081c565b9250505b620001d0565b818512156200033a576200033a8686846200017b565b838312156200035057620003508684866200017b565b505050505050565b828054828255905f5260205f20908101928215620003a1579160200282015b82811115620003a15782518290620003909082620008c2565b509160200191906001019062000377565b50620003af92915062000417565b5090565b828054828255905f5260205f2090810192821562000409579160200282015b828111156200040957825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620003d2565b50620003af92915062000437565b80821115620003af575f6200042d82826200044d565b5060010162000417565b5b80821115620003af575f815560010162000438565b5080546200045b906200083a565b5f825580601f106200046b575050565b601f0160209004905f5260205f209081019062000489919062000437565b50565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f191681016001600160401b0381118282101715620004cb57620004cb6200048c565b604052919050565b5f6001600160401b03821115620004ee57620004ee6200048c565b5060051b60200190565b80516001600160a01b03811681146200050f575f80fd5b919050565b5f601f838184011262000525575f80fd5b825160206200053e6200053883620004d3565b620004a0565b82815260059290921b850181019181810190878411156200055d575f80fd5b8287015b84811015620006115780516001600160401b038082111562000582575f8081fd5b818a0191508a603f83011262000597575f8081fd5b8582015181811115620005ae57620005ae6200048c565b620005c1818a01601f19168801620004a0565b915080825260408c81838601011115620005da575f8081fd5b5f5b82811015620005f9578481018201518482018a01528801620005dc565b50505f90820187015284525091830191830162000561565b50979650505050505050565b5f805f805f60a0868803121562000632575f80fd5b85516001600160401b038082111562000649575f80fd5b818801915088601f8301126200065d575f80fd5b81516020620006706200053883620004d3565b82815260059290921b8401810191818101908c8411156200068f575f80fd5b948201945b83861015620006b857620006a886620004f8565b8252948201949082019062000694565b9950620006c990508a8201620004f8565b97505050620006db60408901620004f8565b94506060880151915080821115620006f1575f80fd5b50620007008882890162000514565b925050608086015190509295509295909350565b634e487b7160e01b5f52601160045260245ffd5b818103818111156200073e576200073e62000714565b92915050565b634e487b7160e01b5f52603260045260245ffd5b5f600182016200076c576200076c62000714565b5060010190565b8181035f83128015838313168383128216171562000795576200079562000714565b5092915050565b5f82620007b757634e487b7160e01b5f52601260045260245ffd5b600160ff1b82145f1984141615620007d357620007d362000714565b500590565b8082018281125f831280158216821582161715620007fa57620007fa62000714565b505092915050565b5f6001600160ff1b0182016200076c576200076c62000714565b5f600160ff1b820162000833576200083362000714565b505f190190565b600181811c908216806200084f57607f821691505b6020821081036200086e57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620008bd575f81815260208120601f850160051c810160208610156200089c5750805b601f850160051c820191505b818110156200035057828155600101620008a8565b505050565b81516001600160401b03811115620008de57620008de6200048c565b620008f681620008ef84546200083a565b8462000874565b602080601f8311600181146200092c575f8415620009145750858301515b5f19600386901b1c1916600185901b17855562000350565b5f85815260208120601f198616915b828110156200095c578886015182559484019460019091019084016200093b565b50858210156200097a57878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61233b80620009985f395ff3fe608060405260043610610121575f3560e01c80638d4f75d2116100a4578063b3ab15fb1161006b578063b3ab15fb14610378578063b78dec5214610397578063ccce413b146103ab578063cdd72253146103ca578063df7f710e146103eb578063e6a02a281461040c57005b80638d4f75d2146103065780639670c0bc146103255780639f8743f71461033a578063a78135871461034e578063aa2f89b51461036357005b80634bb278f3116100e85780634bb278f3146101f75780634c56ea561461021b5780635281b5c6146102615780635412b3ae1461028d578063845023f2146102e757005b806308f21ff51461012a578063146ca53114610152578063307de9b61461016757806333f98c77146101865780633c8510fd146101d857005b3661012857005b005b348015610135575f80fd5b5061013f60085481565b6040519081526020015b60405180910390f35b34801561015d575f80fd5b5061013f60065481565b348015610172575f80fd5b506101286101813660046119a5565b610421565b348015610191575f80fd5b506101a56101a0366004611ad4565b610658565b60405161014991908151815260208083015190820152604080830151908201526060918201519181019190915260800190565b3480156101e3575f80fd5b506101a56101f2366004611b06565b610770565b348015610202575f80fd5b5061020b61086f565b6040519015158152602001610149565b348015610226575f80fd5b5061013f610235366004611b65565b8151602081840181018051600c825292820194820194909420919093529091525f908152604090205481565b34801561026c575f80fd5b5061028061027b366004611bb0565b610a09565b6040516101499190611c14565b348015610298575f80fd5b506102ca6102a7366004611c2d565b600b6020525f908152604090208054600182015460029092015490919060ff1683565b604080519384526020840192909252151590820152606001610149565b3480156102f2575f80fd5b50610128610301366004611c69565b610aaf565b348015610311575f80fd5b50610128610320366004611d01565b610b56565b348015610330575f80fd5b506298968061013f565b348015610345575f80fd5b5060065461013f565b348015610359575f80fd5b5061013f60095481565b34801561036e575f80fd5b5061013f60075481565b348015610383575f80fd5b50610128610392366004611c2d565b610cc5565b3480156103a2575f80fd5b5060095461013f565b3480156103b6575f80fd5b506102806103c5366004611bb0565b610d11565b3480156103d5575f80fd5b506103de610d1f565b6040516101499190611dac565b3480156103f6575f80fd5b506103ff610d7f565b6040516101499190611e4b565b348015610417575f80fd5b5061013f600a5481565b335f908152600b602052604090206002015460ff166104875760405162461bcd60e51b815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f746572730000000000000060448201526064015b60405180910390fd5b600654335f908152600b6020526040902054036104d65760405162461bcd60e51b815260206004820152600d60248201526c185b1c9958591e481d9bdd1959609a1b604482015260640161047e565b335f908152600b6020526040812060018101805490879055815460065490925591819003610505575050610652565b5f548414610514575050610652565b60016006546105239190611e71565b8114158061055f5750848484336040516020016105439493929190611e84565b604051602081830303815290604052805190602001205f1c8214155b156105d5575f5b5f548110156105cd576001600160ff1b03600c5f838154811061058b5761058b611ece565b905f5260205f20016040516105a09190611f1a565b9081526040805160209281900383019020335f9081529252902055806105c581611f8c565b915050610566565b505050610652565b5f5b8481101561064e578585828181106105f1576105f1611ece565b90506020020135600c5f838154811061060c5761060c611ece565b905f5260205f20016040516106219190611f1a565b9081526040805160209281900383019020335f90815292529020558061064681611f8c565b9150506105d7565b5050505b50505050565b61067f60405180608001604052805f81526020015f81526020015f81526020015f81525090565b5f600d60016006546106919190611e71565b815481106106a1576106a1611ece565b905f5260205f2001836040516106b79190611fa4565b90815260200160405180910390206040518060600160405290815f820154815260200160018201548152602001600282015f9054906101000a900460ff16600181111561070657610706611fbf565b600181111561071757610717611fbf565b8152505090505f604051806080016040528060016006546107389190611e71565b8152602001835f01518152602001836020015181526020018360400151600181111561076657610766611fbf565b9052949350505050565b61079760405180608001604052805f81526020015f81526020015f81526020015f81525090565b5f600d84815481106107ab576107ab611ece565b905f5260205f2001836040516107c19190611fa4565b90815260200160405180910390206040518060600160405290815f820154815260200160018201548152602001600282015f9054906101000a900460ff16600181111561081057610810611fbf565b600181111561082157610821611fbf565b8152505090505f6040518060800160405280868152602001835f01518152602001836020015181526020018360400151600181111561086257610862611fbf565b9052925050505b92915050565b6002545f906001600160a01b0316331461089b5760405162461bcd60e51b815260040161047e90611fd3565b600954600a546108ab9190612016565b4310610a04575f5b5f548110156108d7576108c581610f32565b6108d0600182612016565b90506108b3565b5060065460075403610951575f5b60055481101561094f576001600b5f6005848154811061090757610907611ece565b5f918252602080832091909101546001600160a01b031683528201929092526040019020600201805460ff19169115159190911790558061094781611f8c565b9150506108e5565b505b600654600754610962906001612029565b0361096f5761096f611252565b43600a81905550600160065f8282546109889190612016565b909155505060085461099b906002612029565b600654036109b457600180546109b2915f91611808565b505b60065460095460408051928352436020840152429083015260608201527fb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e59060800160405180910390a150600190565b505f90565b60018181548110610a18575f80fd5b905f5260205f20015f915090508054610a3090611ee2565b80601f0160208091040260200160405190810160405280929190818152602001828054610a5c90611ee2565b8015610aa75780601f10610a7e57610100808354040283529160200191610aa7565b820191905f5260205f20905b815481529060010190602001808311610a8a57829003601f168201915b505050505081565b6002546001600160a01b03163314610ad95760405162461bcd60e51b815260040161047e90611fd3565b80515f03610b215760405162461bcd60e51b8152602060048201526015602482015274566f746572732063616e277420626520656d70747960581b604482015260640161047e565b610b39815f60018451610b349190611e71565b611426565b8051610b4c90600590602084019061185c565b5050600654600755565b6003546001600160a01b03163314610ba95760405162461bcd60e51b81526020600482015260166024820152753932b9ba3934b1ba32b2103a379037b832b930ba37b960511b604482015260640161047e565b80515f03610bf25760405162461bcd60e51b815260206004820152601660248201527573796d626f6c732063616e277420626520656d70747960501b604482015260640161047e565b600654600854610c03906001612029565b14158015610c15575060065460085414155b610c615760405162461bcd60e51b815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e640000604482015260640161047e565b8051610c749060019060208401906118bb565b5060065460088190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d908290610cac906001612016565b604051610cba929190612050565b60405180910390a150565b6002546001600160a01b03163314610cef5760405162461bcd60e51b815260040161047e90611fd3565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b5f8181548110610a18575f80fd5b60606005805480602002602001604051908101604052809291908181526020018280548015610d7557602002820191905f5260205f20905b81546001600160a01b03168152600190910190602001808311610d57575b5050505050905090565b60606006546008546001610d939190612029565b03610e6a576001805480602002602001604051908101604052809291908181526020015f905b82821015610e61578382905f5260205f20018054610dd690611ee2565b80601f0160208091040260200160405190810160405280929190818152602001828054610e0290611ee2565b8015610e4d5780601f10610e2457610100808354040283529160200191610e4d565b820191905f5260205f20905b815481529060010190602001808311610e3057829003601f168201915b505050505081526020019060010190610db9565b50505050905090565b5f805480602002602001604051908101604052809291908181526020015f905b82821015610e61578382905f5260205f20018054610ea790611ee2565b80601f0160208091040260200160405190810160405280929190818152602001828054610ed390611ee2565b8015610f1e5780601f10610ef557610100808354040283529160200191610f1e565b820191905f5260205f20905b815481529060010190602001808311610f0157829003601f168201915b505050505081526020019060010190610e8a565b5f808281548110610f4557610f45611ece565b905f5260205f20018054610f5890611ee2565b80601f0160208091040260200160405190810160405280929190818152602001828054610f8490611ee2565b8015610fcf5780601f10610fa657610100808354040283529160200191610fcf565b820191905f5260205f20905b815481529060010190602001808311610fb257829003601f168201915b505050505090505f60048054905067ffffffffffffffff811115610ff557610ff5611a23565b60405190808252806020026020018201604052801561101e578160200160208202803683370190505b5090505f805b600454811015611130575f6004828154811061104257611042611ece565b5f9182526020808320909101546006546001600160a01b03909116808452600b9092526040909220549092501415806110b457506001600160ff1b03600c8660405161108e9190611fa4565b90815260408051602092819003830190206001600160a01b0385165f9081529252902054145b156110bf575061111e565b600c856040516110cf9190611fa4565b90815260408051602092819003830190206001600160a01b0384165f908152925290205484846110fe81611f8c565b95508151811061111057611110611ece565b602002602001018181525050505b8061112881611f8c565b915050611024565b505f600d60016006546111439190611e71565b8154811061115357611153611ece565b905f5260205f2001846040516111699190611fa4565b908152604051908190036020019020549050600182156111935761118d84846115d2565b91505f90505b600d8054600190810182555f9190915260408051606081018252848152426020820152919082019083908111156111cc576111cc611fbf565b815250600d600654815481106111e4576111e4611ece565b905f5260205f2001866040516111fa9190611fa4565b90815260200160405180910390205f820151815f0155602082015181600101556040820151816002015f6101000a81548160ff0219169083600181111561124357611243611fbf565b02179055505050505050505050565b5f805b60045482108015611267575060055481105b156113a6576005818154811061127f5761127f611ece565b5f91825260209091200154600480546001600160a01b0390921691849081106112aa576112aa611ece565b5f918252602090912001546001600160a01b0316036112e357816112cd81611f8c565b92505080806112db90611f8c565b915050611255565b600581815481106112f6576112f6611ece565b5f91825260209091200154600480546001600160a01b03909216918490811061132157611321611ece565b5f918252602090912001546001600160a01b0316101561139c57600b5f6004848154811061135157611351611ece565b5f9182526020808320909101546001600160a01b0316835282019290925260400181208181556001810191909155600201805460ff191690558161139481611f8c565b925050611255565b806112db81611f8c565b60045482101561141157600b5f600484815481106113c6576113c6611ece565b5f9182526020808320909101546001600160a01b0316835282019290925260400181208181556001810191909155600201805460ff191690558161140981611f8c565b9250506113a6565b60058054611421916004916118ff565b505050565b80821261143257505050565b81815f8560026114428585612071565b61144c91906120ab565b6114569087612029565b8151811061146657611466611ece565b602002602001015190505b8183136115a4575b806001600160a01b031686848151811061149557611495611ece565b60200260200101516001600160a01b031610156114be57826114b6816120d7565b935050611479565b806001600160a01b03168683815181106114da576114da611ece565b60200260200101516001600160a01b0316111561150357816114fb816120ee565b9250506114be565b81831361159f5785828151811061151c5761151c611ece565b602002602001015186848151811061153657611536611ece565b602002602001015187858151811061155057611550611ece565b6020026020010188858151811061156957611569611ece565b6001600160a01b039384166020918202929092010152911690528261158d816120d7565b935050818061159b906120ee565b9250505b611471565b818512156115b7576115b7868684611426565b838312156115ca576115ca868486611426565b505050505050565b5f815f036115e157505f610869565b6115f6835f6115f1600186611e71565b611691565b5f611602600284612109565b905061160f60028461211c565b156116335783818151811061162657611626611ece565b6020026020010151611689565b600284828151811061164757611647611ece565b60200260200101518560018461165d9190611e71565b8151811061166d5761166d611ece565b602002602001015161167f9190612029565b61168991906120ab565b949350505050565b81818082036116a1575050505050565b5f8560026116af8787612071565b6116b991906120ab565b6116c39087612029565b815181106116d3576116d3611ece565b602002602001015190505b8183136117e2575b808684815181106116f9576116f9611ece565b602002602001015112156117195782611711816120d7565b9350506116e6565b85828151811061172b5761172b611ece565b602002602001015181121561174c5781611744816120ee565b925050611719565b8183136117dd5785828151811061176557611765611ece565b602002602001015186848151811061177f5761177f611ece565b602002602001015187858151811061179957611799611ece565b602002602001018885815181106117b2576117b2611ece565b602090810291909101019190915252816117cb816120ee565b92505082806117d9906120d7565b9350505b6116de565b818512156117f5576117f5868684611691565b838312156115ca576115ca868486611691565b828054828255905f5260205f2090810192821561184c575f5260205f209182015b8281111561184c578161183c8482612174565b5091600101919060010190611829565b5061185892915061193b565b5090565b828054828255905f5260205f209081019282156118af579160200282015b828111156118af57825182546001600160a01b0319166001600160a01b0390911617825560209092019160019091019061187a565b50611858929150611957565b828054828255905f5260205f2090810192821561184c579160200282015b8281111561184c57825182906118ef908261224f565b50916020019190600101906118d9565b828054828255905f5260205f209081019282156118af575f5260205f209182015b828111156118af578254825591600101919060010190611920565b80821115611858575f61194e828261196b565b5060010161193b565b5b80821115611858575f8155600101611958565b50805461197790611ee2565b5f825580601f10611986575050565b601f0160209004905f5260205f20908101906119a29190611957565b50565b5f805f80606085870312156119b8575f80fd5b84359350602085013567ffffffffffffffff808211156119d6575f80fd5b818701915087601f8301126119e9575f80fd5b8135818111156119f7575f80fd5b8860208260051b8501011115611a0b575f80fd5b95986020929092019750949560400135945092505050565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f1916810167ffffffffffffffff81118282101715611a6057611a60611a23565b604052919050565b5f82601f830112611a77575f80fd5b813567ffffffffffffffff811115611a9157611a91611a23565b611aa4601f8201601f1916602001611a37565b818152846020838601011115611ab8575f80fd5b816020850160208301375f918101602001919091529392505050565b5f60208284031215611ae4575f80fd5b813567ffffffffffffffff811115611afa575f80fd5b61168984828501611a68565b5f8060408385031215611b17575f80fd5b82359150602083013567ffffffffffffffff811115611b34575f80fd5b611b4085828601611a68565b9150509250929050565b80356001600160a01b0381168114611b60575f80fd5b919050565b5f8060408385031215611b76575f80fd5b823567ffffffffffffffff811115611b8c575f80fd5b611b9885828601611a68565b925050611ba760208401611b4a565b90509250929050565b5f60208284031215611bc0575f80fd5b5035919050565b5f5b83811015611be1578181015183820152602001611bc9565b50505f910152565b5f8151808452611c00816020860160208601611bc7565b601f01601f19169290920160200192915050565b602081525f611c266020830184611be9565b9392505050565b5f60208284031215611c3d575f80fd5b611c2682611b4a565b5f67ffffffffffffffff821115611c5f57611c5f611a23565b5060051b60200190565b5f6020808385031215611c7a575f80fd5b823567ffffffffffffffff811115611c90575f80fd5b8301601f81018513611ca0575f80fd5b8035611cb3611cae82611c46565b611a37565b81815260059190911b82018301908381019087831115611cd1575f80fd5b928401925b82841015611cf657611ce784611b4a565b82529284019290840190611cd6565b979650505050505050565b5f6020808385031215611d12575f80fd5b823567ffffffffffffffff80821115611d29575f80fd5b818501915085601f830112611d3c575f80fd5b8135611d4a611cae82611c46565b81815260059190911b83018401908481019088831115611d68575f80fd5b8585015b83811015611d9f57803585811115611d83575f8081fd5b611d918b89838a0101611a68565b845250918601918601611d6c565b5098975050505050505050565b602080825282518282018190525f9190848201906040850190845b81811015611dec5783516001600160a01b031683529284019291840191600101611dc7565b50909695505050505050565b5f81518084526020808501808196508360051b810191508286015f5b85811015611e3e578284038952611e2c848351611be9565b98850198935090840190600101611e14565b5091979650505050505050565b602081525f611c266020830184611df8565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561086957610869611e5d565b5f8186825b87811015611ea7578135835260209283019290910190600101611e89565b5050938452505060601b6bffffffffffffffffffffffff1916602082015260340192915050565b634e487b7160e01b5f52603260045260245ffd5b600181811c90821680611ef657607f821691505b602082108103611f1457634e487b7160e01b5f52602260045260245ffd5b50919050565b5f808354611f2781611ee2565b60018281168015611f3f5760018114611f5457611f80565b60ff1984168752821515830287019450611f80565b875f526020805f205f5b85811015611f775781548a820152908401908201611f5e565b50505082870194505b50929695505050505050565b5f60018201611f9d57611f9d611e5d565b5060010190565b5f8251611fb5818460208701611bc7565b9190910192915050565b634e487b7160e01b5f52602160045260245ffd5b60208082526023908201527f7265737472696374656420746f20746865206175746f6e69747920636f6e74726040820152621858dd60ea1b606082015260800190565b8082018082111561086957610869611e5d565b8082018281125f83128015821682158216171561204857612048611e5d565b505092915050565b604081525f6120626040830185611df8565b90508260208301529392505050565b8181035f83128015838313168383128216171561209057612090611e5d565b5092915050565b634e487b7160e01b5f52601260045260245ffd5b5f826120b9576120b9612097565b600160ff1b82145f19841416156120d2576120d2611e5d565b500590565b5f6001600160ff1b018201611f9d57611f9d611e5d565b5f600160ff1b820161210257612102611e5d565b505f190190565b5f8261211757612117612097565b500490565b5f8261212a5761212a612097565b500690565b601f821115611421575f81815260208120601f850160051c810160208610156121555750805b601f850160051c820191505b818110156115ca57828155600101612161565b81810361217f575050565b6121898254611ee2565b67ffffffffffffffff8111156121a1576121a1611a23565b6121b5816121af8454611ee2565b8461212f565b5f601f8211600181146121e6575f83156121cf5750848201545b5f19600385901b1c1916600184901b178455612248565b5f85815260209020601f198416905f86815260209020845b8381101561221e57828601548255600195860195909101906020016121fe565b508583101561223b57818501545f19600388901b60f8161c191681555b50505060018360011b0184555b5050505050565b815167ffffffffffffffff81111561226957612269611a23565b612277816121af8454611ee2565b602080601f8311600181146122aa575f84156122935750858301515b5f19600386901b1c1916600185901b1785556115ca565b5f85815260208120601f198616915b828110156122d8578886015182559484019460019091019084016122b9565b50858210156122f557878501515f19600388901b60f8161c191681555b5050505050600190811b0190555056fea2646970667358221220ed6f0ed927459bafb9b0c82d428ec27e6660970a4ee16d8633e1d2278975a03964736f6c63430008150033",
}

// OracleABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleMetaData.ABI instead.
var OracleABI = OracleMetaData.ABI

// Deprecated: Use OracleMetaData.Sigs instead.
// OracleFuncSigs maps the 4-byte function signature to its string representation.
var OracleFuncSigs = OracleMetaData.Sigs

// OracleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OracleMetaData.Bin instead.
var OracleBin = OracleMetaData.Bin

// DeployOracle deploys a new Ethereum contract, binding an instance of Oracle to it.
func (r *runner) deployOracle(opts *runOptions, _voters []common.Address, _autonity common.Address, _operator common.Address, _symbols []string, _votePeriod *big.Int) (common.Address, uint64, *Oracle, error) {
	parsed, err := OracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(OracleBin), _voters, _autonity, _operator, _symbols, _votePeriod)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &Oracle{OracleCaller: OracleCaller{contract: c}, OracleTransactor: OracleTransactor{contract: c}, OracleFilterer: OracleFilterer{contract: c}}, nil
}

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	OracleCaller     // Read-only binding to the contract
	OracleTransactor // Write-only binding to the contract
	OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type OracleCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() pure returns(uint256)
func (_Oracle *OracleCaller) GetPrecision(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "getPrecision")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle *OracleCaller) GetRound(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "getRound")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCaller) GetRoundData(opts *runOptions, _round *big.Int, _symbol string) (IOracleRoundData, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)
	return out0, consumed, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle *OracleCaller) GetSymbols(opts *runOptions) ([]string, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "getSymbols")

	if err != nil {
		return *new([]string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)
	return out0, consumed, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle *OracleCaller) GetVotePeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle *OracleCaller) GetVoters(opts *runOptions) ([]common.Address, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "getVoters")

	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// LastRoundBlock is a free data retrieval call binding the contract method 0xe6a02a28.
//
// Solidity: function lastRoundBlock() view returns(uint256)
func (_Oracle *OracleCaller) LastRoundBlock(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "lastRoundBlock")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// LastVoterUpdateRound is a free data retrieval call binding the contract method 0xaa2f89b5.
//
// Solidity: function lastVoterUpdateRound() view returns(int256)
func (_Oracle *OracleCaller) LastVoterUpdateRound(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "lastVoterUpdateRound")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCaller) LatestRoundData(opts *runOptions, _symbol string) (IOracleRoundData, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)
	return out0, consumed, err

}

// NewSymbols is a free data retrieval call binding the contract method 0x5281b5c6.
//
// Solidity: function newSymbols(uint256 ) view returns(string)
func (_Oracle *OracleCaller) NewSymbols(opts *runOptions, arg0 *big.Int) (string, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "newSymbols", arg0)

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// Reports is a free data retrieval call binding the contract method 0x4c56ea56.
//
// Solidity: function reports(string , address ) view returns(int256)
func (_Oracle *OracleCaller) Reports(opts *runOptions, arg0 string, arg1 common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "reports", arg0, arg1)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_Oracle *OracleCaller) Round(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "round")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// SymbolUpdatedRound is a free data retrieval call binding the contract method 0x08f21ff5.
//
// Solidity: function symbolUpdatedRound() view returns(int256)
func (_Oracle *OracleCaller) SymbolUpdatedRound(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "symbolUpdatedRound")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Symbols is a free data retrieval call binding the contract method 0xccce413b.
//
// Solidity: function symbols(uint256 ) view returns(string)
func (_Oracle *OracleCaller) Symbols(opts *runOptions, arg0 *big.Int) (string, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "symbols", arg0)

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// VotePeriod is a free data retrieval call binding the contract method 0xa7813587.
//
// Solidity: function votePeriod() view returns(uint256)
func (_Oracle *OracleCaller) VotePeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "votePeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// VotingInfo is a free data retrieval call binding the contract method 0x5412b3ae.
//
// Solidity: function votingInfo(address ) view returns(uint256 round, uint256 commit, bool isVoter)
func (_Oracle *OracleCaller) VotingInfo(opts *runOptions, arg0 common.Address) (struct {
	Round   *big.Int
	Commit  *big.Int
	IsVoter bool
}, uint64, error) {
	out, consumed, err := _Oracle.call(opts, "votingInfo", arg0)

	outstruct := new(struct {
		Round   *big.Int
		Commit  *big.Int
		IsVoter bool
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.Round = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Commit = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IsVoter = *abi.ConvertType(out[2], new(bool)).(*bool)
	return *outstruct, consumed, err

}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_Oracle *OracleTransactor) Finalize(opts *runOptions) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "finalize")
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle *OracleTransactor) SetOperator(opts *runOptions, _operator common.Address) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "setOperator", _operator)
	return consumed, err
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle *OracleTransactor) SetSymbols(opts *runOptions, _symbols []string) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "setSymbols", _symbols)
	return consumed, err
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_Oracle *OracleTransactor) SetVoters(opts *runOptions, _newVoters []common.Address) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "setVoters", _newVoters)
	return consumed, err
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_Oracle *OracleTransactor) Vote(opts *runOptions, _commit *big.Int, _reports []*big.Int, _salt *big.Int) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "vote", _commit, _reports, _salt)
	return consumed, err
}

// Fallback is a paid mutator transaction binding the contract fallback function.
// WARNING! UNTESTED
// Solidity: fallback() payable returns()
func (_Oracle *OracleTransactor) Fallback(opts *runOptions, calldata []byte) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "", calldata)
	return consumed, err
}

// Receive is a paid mutator transaction binding the contract receive function.
// WARNING! UNTESTED
// Solidity: receive() payable returns()
func (_Oracle *OracleTransactor) Receive(opts *runOptions) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "")
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// OracleNewRoundIterator is returned from FilterNewRound and is used to iterate over the raw logs and unpacked data for NewRound events raised by the Oracle contract.
		type OracleNewRoundIterator struct {
			Event *OracleNewRound // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *OracleNewRoundIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(OracleNewRound)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(OracleNewRound)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *OracleNewRoundIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *OracleNewRoundIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// OracleNewRound represents a NewRound event raised by the Oracle contract.
		type OracleNewRound struct {
			Round *big.Int;
			Height *big.Int;
			Timestamp *big.Int;
			VotePeriod *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewRound is a free log retrieval operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
		//
		// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
 		func (_Oracle *OracleFilterer) FilterNewRound(opts *bind.FilterOpts) (*OracleNewRoundIterator, error) {






			logs, sub, err := _Oracle.contract.FilterLogs(opts, "NewRound")
			if err != nil {
				return nil, err
			}
			return &OracleNewRoundIterator{contract: _Oracle.contract, event: "NewRound", logs: logs, sub: sub}, nil
 		}

		// WatchNewRound is a free log subscription operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
		//
		// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
		func (_Oracle *OracleFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *OracleNewRound) (event.Subscription, error) {






			logs, sub, err := _Oracle.contract.WatchLogs(opts, "NewRound")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(OracleNewRound)
						if err := _Oracle.contract.UnpackLog(event, "NewRound", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewRound is a log parse operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
		//
		// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
		func (_Oracle *OracleFilterer) ParseNewRound(log types.Log) (*OracleNewRound, error) {
			event := new(OracleNewRound)
			if err := _Oracle.contract.UnpackLog(event, "NewRound", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// OracleNewSymbolsIterator is returned from FilterNewSymbols and is used to iterate over the raw logs and unpacked data for NewSymbols events raised by the Oracle contract.
		type OracleNewSymbolsIterator struct {
			Event *OracleNewSymbols // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *OracleNewSymbolsIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(OracleNewSymbols)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(OracleNewSymbols)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *OracleNewSymbolsIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *OracleNewSymbolsIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// OracleNewSymbols represents a NewSymbols event raised by the Oracle contract.
		type OracleNewSymbols struct {
			Symbols []string;
			Round *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterNewSymbols is a free log retrieval operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
		//
		// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
 		func (_Oracle *OracleFilterer) FilterNewSymbols(opts *bind.FilterOpts) (*OracleNewSymbolsIterator, error) {




			logs, sub, err := _Oracle.contract.FilterLogs(opts, "NewSymbols")
			if err != nil {
				return nil, err
			}
			return &OracleNewSymbolsIterator{contract: _Oracle.contract, event: "NewSymbols", logs: logs, sub: sub}, nil
 		}

		// WatchNewSymbols is a free log subscription operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
		//
		// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
		func (_Oracle *OracleFilterer) WatchNewSymbols(opts *bind.WatchOpts, sink chan<- *OracleNewSymbols) (event.Subscription, error) {




			logs, sub, err := _Oracle.contract.WatchLogs(opts, "NewSymbols")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(OracleNewSymbols)
						if err := _Oracle.contract.UnpackLog(event, "NewSymbols", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseNewSymbols is a log parse operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
		//
		// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
		func (_Oracle *OracleFilterer) ParseNewSymbols(log types.Log) (*OracleNewSymbols, error) {
			event := new(OracleNewSymbols)
			if err := _Oracle.contract.UnpackLog(event, "NewSymbols", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// OracleVotedIterator is returned from FilterVoted and is used to iterate over the raw logs and unpacked data for Voted events raised by the Oracle contract.
		type OracleVotedIterator struct {
			Event *OracleVoted // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *OracleVotedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(OracleVoted)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(OracleVoted)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *OracleVotedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *OracleVotedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// OracleVoted represents a Voted event raised by the Oracle contract.
		type OracleVoted struct {
			Voter common.Address;
			Votes []*big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterVoted is a free log retrieval operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
		//
		// Solidity: event Voted(address indexed _voter, int256[] _votes)
 		func (_Oracle *OracleFilterer) FilterVoted(opts *bind.FilterOpts, _voter []common.Address) (*OracleVotedIterator, error) {

			var _voterRule []interface{}
			for _, _voterItem := range _voter {
				_voterRule = append(_voterRule, _voterItem)
			}


			logs, sub, err := _Oracle.contract.FilterLogs(opts, "Voted", _voterRule)
			if err != nil {
				return nil, err
			}
			return &OracleVotedIterator{contract: _Oracle.contract, event: "Voted", logs: logs, sub: sub}, nil
 		}

		// WatchVoted is a free log subscription operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
		//
		// Solidity: event Voted(address indexed _voter, int256[] _votes)
		func (_Oracle *OracleFilterer) WatchVoted(opts *bind.WatchOpts, sink chan<- *OracleVoted, _voter []common.Address) (event.Subscription, error) {

			var _voterRule []interface{}
			for _, _voterItem := range _voter {
				_voterRule = append(_voterRule, _voterItem)
			}


			logs, sub, err := _Oracle.contract.WatchLogs(opts, "Voted", _voterRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(OracleVoted)
						if err := _Oracle.contract.UnpackLog(event, "Voted", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseVoted is a log parse operation binding the contract event 0xd0d8560f1076ac6b216b1091a2571d6f9bc3e0889f4dbdbe1c7d1be7136714d3.
		//
		// Solidity: event Voted(address indexed _voter, int256[] _votes)
		func (_Oracle *OracleFilterer) ParseVoted(log types.Log) (*OracleVoted, error) {
			event := new(OracleVoted)
			if err := _Oracle.contract.UnpackLog(event, "Voted", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// PrecompiledMetaData contains all meta data concerning the Precompiled contract.
var PrecompiledMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ACCUSATION_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"COMPUTE_COMMITTEE_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ENODE_VERIFIER_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INNOCENCE_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MISBEHAVIOUR_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POP_VERIFIER_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SUCCESS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADER_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"4dc925d3": "ACCUSATION_CONTRACT()",
		"2090a442": "COMPUTE_COMMITTEE_CONTRACT()",
		"c13974e1": "ENODE_VERIFIER_CONTRACT()",
		"8e153dc3": "INNOCENCE_CONTRACT()",
		"925c5492": "MISBEHAVIOUR_CONTRACT()",
		"50d93720": "POP_VERIFIER_CONTRACT()",
		"d0a6d1a6": "SUCCESS()",
		"a4ad5d91": "UPGRADER_CONTRACT()",
	},
	Bin: "0x610122610035600b8282823980515f1a60731461002957634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106086575f3560e01c8063925c549211605d578063925c54921460c3578063a4ad5d911460ca578063c13974e11460d1578063d0a6d1a61460d8575f80fd5b80632090a44214608a5780634dc925d31460ae57806350d937201460b55780638e153dc31460bc575b5f80fd5b609160fa81565b6040516001600160a01b0390911681526020015b60405180910390f35b609160fc81565b609160fb81565b609160fd81565b609160fe81565b609160f981565b609160ff81565b60df600181565b60405190815260200160a556fea2646970667358221220322ce5d034a4ee14d724706036655b098fc38338efa4bf43a833b4d6743bb67464736f6c63430008150033",
}

// PrecompiledABI is the input ABI used to generate the binding from.
// Deprecated: Use PrecompiledMetaData.ABI instead.
var PrecompiledABI = PrecompiledMetaData.ABI

// Deprecated: Use PrecompiledMetaData.Sigs instead.
// PrecompiledFuncSigs maps the 4-byte function signature to its string representation.
var PrecompiledFuncSigs = PrecompiledMetaData.Sigs

// PrecompiledBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PrecompiledMetaData.Bin instead.
var PrecompiledBin = PrecompiledMetaData.Bin

// DeployPrecompiled deploys a new Ethereum contract, binding an instance of Precompiled to it.
func (r *runner) deployPrecompiled(opts *runOptions) (common.Address, uint64, *Precompiled, error) {
	parsed, err := PrecompiledMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(PrecompiledBin))
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &Precompiled{PrecompiledCaller: PrecompiledCaller{contract: c}, PrecompiledTransactor: PrecompiledTransactor{contract: c}, PrecompiledFilterer: PrecompiledFilterer{contract: c}}, nil
}

// Precompiled is an auto generated Go binding around an Ethereum contract.
type Precompiled struct {
	PrecompiledCaller     // Read-only binding to the contract
	PrecompiledTransactor // Write-only binding to the contract
	PrecompiledFilterer   // Log filterer for contract events
}

// PrecompiledCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type PrecompiledCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// PrecompiledTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PrecompiledTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// PrecompiledFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PrecompiledFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) ACCUSATIONCONTRACT(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Precompiled.call(opts, "ACCUSATION_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// COMPUTECOMMITTEECONTRACT is a free data retrieval call binding the contract method 0x2090a442.
//
// Solidity: function COMPUTE_COMMITTEE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) COMPUTECOMMITTEECONTRACT(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Precompiled.call(opts, "COMPUTE_COMMITTEE_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// ENODEVERIFIERCONTRACT is a free data retrieval call binding the contract method 0xc13974e1.
//
// Solidity: function ENODE_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) ENODEVERIFIERCONTRACT(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Precompiled.call(opts, "ENODE_VERIFIER_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// INNOCENCECONTRACT is a free data retrieval call binding the contract method 0x8e153dc3.
//
// Solidity: function INNOCENCE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) INNOCENCECONTRACT(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Precompiled.call(opts, "INNOCENCE_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// MISBEHAVIOURCONTRACT is a free data retrieval call binding the contract method 0x925c5492.
//
// Solidity: function MISBEHAVIOUR_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) MISBEHAVIOURCONTRACT(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Precompiled.call(opts, "MISBEHAVIOUR_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// POPVERIFIERCONTRACT is a free data retrieval call binding the contract method 0x50d93720.
//
// Solidity: function POP_VERIFIER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) POPVERIFIERCONTRACT(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Precompiled.call(opts, "POP_VERIFIER_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// SUCCESS is a free data retrieval call binding the contract method 0xd0a6d1a6.
//
// Solidity: function SUCCESS() view returns(uint256)
func (_Precompiled *PrecompiledCaller) SUCCESS(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Precompiled.call(opts, "SUCCESS")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UPGRADERCONTRACT is a free data retrieval call binding the contract method 0xa4ad5d91.
//
// Solidity: function UPGRADER_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) UPGRADERCONTRACT(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Precompiled.call(opts, "UPGRADER_CONTRACT")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

/* EVENTS ARE NOT YET SUPPORTED

 */

// StabilizationMetaData contains all meta data concerning the Stabilization contract.
var StabilizationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"borrowInterestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minCollateralizationRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minDebtRequirement\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetPrice\",\"type\":\"uint256\"}],\"internalType\":\"structStabilization.Config\",\"name\":\"config_\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"supplyControl\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"collateralToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientCollateral\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientPayment\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDebtPosition\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidParameter\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Liquidatable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoDebtPosition\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotLiquidatable\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv18_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"denominator\",\"type\":\"uint256\"}],\"name\":\"PRBMath_MulDiv_Overflow\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp2_InputTooBig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"UD60x18\",\"name\":\"x\",\"type\":\"uint256\"}],\"name\":\"PRBMath_UD60x18_Exp_InputTooBig\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PriceUnavailable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroValue\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Borrow\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidator\",\"type\":\"address\"}],\"name\":\"Liquidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Repay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"SCALE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SCALE_FACTOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SECONDS_IN_YEAR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accounts\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"borrow\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mcr\",\"type\":\"uint256\"}],\"name\":\"borrowLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"cdps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interest\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"collateralPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"borrowInterestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minCollateralizationRatio\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minDebtRequirement\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetPrice\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"debtAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"debtAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeBorrow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeDue\",\"type\":\"uint256\"}],\"name\":\"interestDue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isLiquidatable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"liquidate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mcr\",\"type\":\"uint256\"}],\"name\":\"minimumCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"repay\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ratio\",\"type\":\"uint256\"}],\"name\":\"setLiquidationRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ratio\",\"type\":\"uint256\"}],\"name\":\"setMinCollateralizationRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"setMinDebtRequirement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"supplyControl\",\"type\":\"address\"}],\"name\":\"setSupplyControl\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"collateral\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"debt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationRatio\",\"type\":\"uint256\"}],\"name\":\"underCollateralized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"eced5526": "SCALE()",
		"ce4b5bbe": "SCALE_FACTOR()",
		"5dcc9391": "SECONDS_IN_YEAR()",
		"68cd03f6": "accounts()",
		"c5ebeaec": "borrow(uint256)",
		"83baa174": "borrowLimit(uint256,uint256,uint256,uint256)",
		"840c7e24": "cdps(address)",
		"5891de72": "collateralPrice()",
		"79502c55": "config()",
		"54a9f42c": "debtAmount(address)",
		"50bf06bf": "debtAmount(address,uint256)",
		"b6b55f25": "deposit(uint256)",
		"15184245": "interestDue(uint256,uint256,uint256,uint256)",
		"042e02cf": "isLiquidatable(address)",
		"2f865568": "liquidate(address)",
		"08796696": "minimumCollateral(uint256,uint256,uint256)",
		"402d8883": "repay()",
		"946ce8cd": "setLiquidationRatio(uint256)",
		"7b44646a": "setMinCollateralizationRatio(uint256)",
		"53afe81d": "setMinDebtRequirement(uint256)",
		"b3ab15fb": "setOperator(address)",
		"7adbf973": "setOracle(address)",
		"52e5a050": "setSupplyControl(address)",
		"fbbe6991": "underCollateralized(uint256,uint256,uint256,uint256)",
		"2e1a7d4d": "withdraw(uint256)",
	},
	Bin: "0x608060405234801562000010575f80fd5b50604051620026363803806200263683398101604081905262000033916200012a565b8560400151805f036200005957604051630309cb8760e51b815260040160405180910390fd5b866020015187604001518082106200008457604051630309cb8760e51b815260040160405180910390fd5b505086515f5550602086015160015560408601516002556060860151600355608090950151600455600780546001600160a01b03199081166001600160a01b039687161790915560088054821694861694909417909355600a8054841692851692909217909155600b80548316918416919091179055600980549091169190921617905562000211565b80516001600160a01b038116811462000125575f80fd5b919050565b5f805f805f8086880361014081121562000142575f80fd5b60a081121562000150575f80fd5b5060405160a081016001600160401b03811182821017156200018057634e487b7160e01b5f52604160045260245ffd5b8060405250875181526020880151602082015260408801516040820152606088015160608201526080880151608082015280965050620001c360a088016200010e565b9450620001d360c088016200010e565b9350620001e360e088016200010e565b9250620001f461010088016200010e565b91506200020561012088016200010e565b90509295509295509295565b612417806200021f5f395ff3fe608060405260043610610161575f3560e01c806368cd03f6116100cd578063946ce8cd11610087578063c5ebeaec11610062578063c5ebeaec14610450578063ce4b5bbe1461046f578063eced552614610483578063fbbe699114610497575f80fd5b8063946ce8cd146103f3578063b3ab15fb14610412578063b6b55f2514610431575f80fd5b806368cd03f6146102c857806379502c55146102e95780637adbf973146103365780637b44646a1461035557806383baa17414610374578063840c7e2414610393575f80fd5b806350bf06bf1161011e57806350bf06bf1461022157806352e5a0501461024057806353afe81d1461025f57806354a9f42c1461027e5780635891de721461029d5780635dcc9391146102b1575f80fd5b8063042e02cf14610165578063087966961461019957806315184245146101c65780632e1a7d4d146101e55780632f86556814610206578063402d888314610219575b5f80fd5b348015610170575f80fd5b5061018461017f366004612085565b6104b6565b60405190151581526020015b60405180910390f35b3480156101a4575f80fd5b506101b86101b336600461209e565b6104fd565b604051908152602001610190565b3480156101d1575f80fd5b506101b86101e03660046120c7565b610565565b3480156101f0575f80fd5b506102046101ff3660046120f6565b6105f1565b005b610204610214366004612085565b6107b5565b6102046109e7565b34801561022c575f80fd5b506101b861023b36600461210d565b610bb0565b34801561024b575f80fd5b5061020461025a366004612085565b610c1a565b34801561026a575f80fd5b506102046102793660046120f6565b610c66565b348015610289575f80fd5b506101b8610298366004612085565b610c95565b3480156102a8575f80fd5b506101b8610d09565b3480156102bc575f80fd5b506101b86301e1338081565b3480156102d3575f80fd5b506102dc610f9d565b6040516101909190612135565b3480156102f4575f80fd5b505f5460015460025460035460045461030e949392919085565b604080519586526020860194909452928401919091526060830152608082015260a001610190565b348015610341575f80fd5b50610204610350366004612085565b610ffd565b348015610360575f80fd5b5061020461036f3660046120f6565b611049565b34801561037f575f80fd5b506101b861038e3660046120c7565b6110c0565b34801561039e575f80fd5b506103d36103ad366004612085565b60056020525f908152604090208054600182015460028301546003909301549192909184565b604080519485526020850193909352918301526060820152608001610190565b3480156103fe575f80fd5b5061020461040d3660046120f6565b611115565b34801561041d575f80fd5b5061020461042c366004612085565b61116b565b34801561043c575f80fd5b5061020461044b3660046120f6565b6111b7565b34801561045b575f80fd5b5061020461046a3660046120f6565b6113ad565b34801561047a575f80fd5b506101b861156b565b34801561048e575f80fd5b506101b8601281565b3480156104a2575f80fd5b506101846104b13660046120c7565b61157a565b6001600160a01b0381165f908152600560205260408120816104d882426115cd565b5090506104f582600101546104eb610d09565b600154849061157a565b949350505050565b5f82805f0361051e5760405162bfc92160e01b815260040160405180910390fd5b831580610529575082155b1561054757604051630309cb8760e51b815260040160405180910390fd5b836105528487612195565b61055c91906121c0565b95945050505050565b5f8183111561058757604051630309cb8760e51b815260040160405180910390fd5b84845f6105a86301e133806105a261059f89896121df565b90565b9061163e565b90505f6105bd6105b8848461165c565b61166a565b90505f6105e36105dc6105d561059f6012600a6122d2565b84906116bc565b869061165c565b9a9950505050505050505050565b80805f036106125760405163162908e360e11b815260040160405180910390fd5b335f90815260056020526040902060018101548311156106455760405163162908e360e11b815260040160405180910390fd5b5f61065082426115cd565b5090505f61065c610d09565b9050610672836001015482845f6001015461157a565b1561069057604051636229415360e01b815260040160405180910390fd5b6106a38360020154825f600201546104fd565b8584600101546106b391906121df565b10156106d257604051633a23d82560e01b815260040160405180910390fd5b84836001015f8282546106e591906121df565b909155505060095460405163a9059cbb60e01b8152336004820152602481018790526001600160a01b039091169063a9059cbb906044016020604051808303815f875af1158015610738573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061075c91906122dd565b610779576040516312171d8360e31b815260040160405180910390fd5b60405185815233907f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649060200160405180910390a25050505050565b345f036107d557604051637c946ed760e01b815260040160405180910390fd5b6001600160a01b0381165f908152600560205260408120600281015490910361081157604051638aa5baf360e01b815260040160405180910390fd5b5f8061081d83426115cd565b9150915061083b8360010154610831610d09565b600154859061157a565b61085857604051636ef5bcdd60e11b815260040160405180910390fd5b5f61086383346121df565b90506001840180544286555f9182905560028601829055600386019190915560095460405163a9059cbb60e01b8152336004820152602481018390526001600160a01b039091169063a9059cbb906044016020604051808303815f875af11580156108d0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906108f491906122dd565b610911576040516312171d8360e31b815260040160405180910390fd5b600b546001600160a01b03166344df8e7061092c85876121df565b6040518263ffffffff1660e01b81526004015f604051808303818588803b158015610955575f80fd5b505af1158015610967573d5f803e3d5ffd5b50505050505f8211156109a057604051339083156108fc029084905f818181858888f1935050505015801561099e573d5f803e3d5ffd5b505b6040513381526001600160a01b038716907fc3d81b2125598b9a2b024afe09e33981f0aa5b7bcbe3e30c4303a4dec209ddb4906020015b60405180910390a2505050505050565b345f03610a0757604051637c946ed760e01b815260040160405180910390fd5b335f9081526005602052604081206002810154909103610a3a57604051638aa5baf360e01b815260040160405180910390fd5b5f80610a4683426115cd565b915091508134108015610a635750600354610a6134846121df565b105b15610a815760405163e6bd447960e01b815260040160405180910390fd5b80836003015f828254610a9491906122fc565b90915550504283555f8080610aa986346116ca565b92509250925081866002015f828254610ac291906121df565b9250508190555082866003015f828254610adc91906121df565b90915550508115610b4c57600b5f9054906101000a90046001600160a01b03166001600160a01b03166344df8e70836040518263ffffffff1660e01b81526004015f604051808303818588803b158015610b34575f80fd5b505af1158015610b46573d5f803e3d5ffd5b50505050505b8015610b7e57604051339082156108fc029083905f818181858888f19350505050158015610b7c573d5f803e3d5ffd5b505b60405134815233907f5c16de4f8b59bd9caf0f49a545f25819a895ed223294290b408242e72a594231906020016109d7565b6001600160a01b0382165f908152600560205260408120805484918491821015610bed57604051630309cb8760e51b815260040160405180910390fd5b6001600160a01b0386165f908152600560205260409020610c0e81876115cd565b50979650505050505050565b6008546001600160a01b03163314610c44576040516282b42960e81b815260040160405180910390fd5b600b80546001600160a01b0319166001600160a01b0392909216919091179055565b6008546001600160a01b03163314610c90576040516282b42960e81b815260040160405180910390fd5b600355565b6040516350bf06bf60e01b81526001600160a01b03821660048201524260248201525f9030906350bf06bf90604401602060405180830381865afa158015610cdf573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d03919061230f565b92915050565b600a546040805180820182526007815266272a2716a0aa2760c91b602082015290516333f98c7760e01b81525f9283926001600160a01b03909116916333f98c7791610d5791600401612326565b608060405180830381865afa158015610d72573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d969190612371565b905080606001515f14610dbc5760405163cb08be8160e01b815260040160405180910390fd5b5f816020015113610ddf5760405162bfc92160e01b815260040160405180910390fd5b600a5f9054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e2f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e53919061230f565b610e5f6012600a6122d2565b1115610f0457600a5f9054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610eb5573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ed9919061230f565b610ee56012600a6122d2565b610eef91906121c0565b8160200151610efe9190612195565b91505090565b610f106012600a6122d2565b600a5f9054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f60573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610f84919061230f565b610f8e91906121c0565b8160200151610efe91906121c0565b60606006805480602002602001604051908101604052809291908181526020018280548015610ff357602002820191905f5260205f20905b81546001600160a01b03168152600190910190602001808311610fd5575b5050505050905090565b6007546001600160a01b03163314611027576040516282b42960e81b815260040160405180910390fd5b600a80546001600160a01b0319166001600160a01b0392909216919091179055565b80805f0361106a57604051630309cb8760e51b815260040160405180910390fd5b6001548280821061108e57604051630309cb8760e51b815260040160405180910390fd5b6008546001600160a01b031633146110b8576040516282b42960e81b815260040160405180910390fd5b505050600255565b5f8315806110cc575081155b156110ea57604051630309cb8760e51b815260040160405180910390fd5b6110f66012600a6122d2565b6111009083612195565b8361110b8688612195565b6105529190612195565b600254819080821061113a57604051630309cb8760e51b815260040160405180910390fd5b6008546001600160a01b03163314611164576040516282b42960e81b815260040160405180910390fd5b5050600155565b6007546001600160a01b03163314611195576040516282b42960e81b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0392909216919091179055565b80805f036111d85760405163162908e360e11b815260040160405180910390fd5b600954604051636eb1769f60e11b815233600482015230602482015283916001600160a01b03169063dd62ed3e90604401602060405180830381865afa158015611224573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611248919061230f565b1015611267576040516313be252b60e01b815260040160405180910390fd5b335f90815260056020526040812080549091036112c057600680546001810182555f919091527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0180546001600160a01b031916331790555b4281556001810180548491905f906112d99084906122fc565b90915550506009546040516323b872dd60e01b8152336004820152306024820152604481018590526001600160a01b03909116906323b872dd906064016020604051808303815f875af1158015611332573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061135691906122dd565b611373576040516312171d8360e31b815260040160405180910390fd5b60405183815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2505050565b80805f036113ce5760405163162908e360e11b815260040160405180910390fd5b335f90815260056020526040812090806113e883426115cd565b90925090506113f785836122fc565b60035490925082101561141d5760405163e6bd447960e01b815260040160405180910390fd5b5f611426610d09565b905061143c846001015482855f6001015461157a565b1561145a57604051636229415360e01b815260040160405180910390fd5b5f6114738560010154835f600401545f600201546110c0565b90508084111561149657604051633a23d82560e01b815260040160405180910390fd5b4285556002850180548891905f906114af9084906122fc565b9250508190555082856003015f8282546114c991906122fc565b9091555050600b546040516340c10f1960e01b8152336004820152602481018990526001600160a01b03909116906340c10f19906044015f604051808303815f87803b158015611517575f80fd5b505af1158015611529573d5f803e3d5ffd5b50506040518981523392507fcbc04eca7e9da35cb1393a6135a199ca52e450d5e9251cbd99f7847d33a36750915060200160405180910390a250505050505050565b6115776012600a6122d2565b81565b5f83805f0361159b5760405162bfc92160e01b815260040160405180910390fd5b835f036115aa575f91506115c4565b82846115b68789612195565b6115c091906121c0565b1091505b50949350505050565b5f80825f036115ef57604051630309cb8760e51b815260040160405180910390fd5b5f8460030154856002015461160491906122fc565b85549091508403611617575f915061162a565b5f54855461162791839187610565565b91505b61163482826122fc565b9250509250929050565b5f61165561059f84670de0b6b3a76400008561173b565b9392505050565b5f61165561059f8484611809565b5f81680736ea4425c11ac63081111561169e57604051630d7b1d6560e11b8152600481018490526024015b60405180910390fd5b6714057b7ef767814f81026104f5670de0b6b3a764000082046118bb565b5f61165561059f83856121df565b5f805f80856003015486600201546116e291906122fc565b9050856003015485106116f95785600301546116fb565b845b935080851061170e578560020154611718565b61171884866121df565b9250808511611727575f611731565b61173181866121df565b9150509250925092565b5f80805f19858709858702925082811083820303915050805f0361177257838281611768576117686121ac565b0492505050611655565b8381106117a357604051630c740aef60e31b8152600481018790526024810186905260448101859052606401611695565b5f84868809600260036001881981018916988990049182028318808302840302808302840302808302840302808302840302808302840302918202909203025f889003889004909101858311909403939093029303949094049190911702949350505050565b5f80805f19848609848602925082811083820303915050805f0361183a5750670de0b6b3a764000090049050610d03565b670de0b6b3a7640000811061186c57604051635173648d60e01b81526004810186905260248101859052604401611695565b5f670de0b6b3a764000085870962040000818503049310909103600160ee1b02919091177faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106690291505092915050565b5f81680a688906bd8affffff8111156118ea5760405163b3b6ba1f60e01b815260048101849052602401611695565b5f611901670de0b6b3a7640000604084901b6121c0565b90506104f561059f82600160bf1b67ff00000000000000821615611a175767800000000000000082161561193e5768016a09e667f3bcc9090260401c5b67400000000000000082161561195d576801306fe0a31b7152df0260401c5b67200000000000000082161561197c576801172b83c7d517adce0260401c5b67100000000000000082161561199b5768010b5586cf9890f62a0260401c5b6708000000000000008216156119ba576801059b0d31585743ae0260401c5b6704000000000000008216156119d957680102c9a3e778060ee70260401c5b6702000000000000008216156119f85768010163da9fb33356d80260401c5b670100000000000000821615611a1757680100b1afa5abcbed610260401c5b66ff000000000000821615611b16576680000000000000821615611a445768010058c86da1c09ea20260401c5b6640000000000000821615611a62576801002c605e2e8cec500260401c5b6620000000000000821615611a8057680100162f3904051fa10260401c5b6610000000000000821615611a9e576801000b175effdc76ba0260401c5b6608000000000000821615611abc57680100058ba01fb9f96d0260401c5b6604000000000000821615611ada5768010002c5cc37da94920260401c5b6602000000000000821615611af8576801000162e525ee05470260401c5b6601000000000000821615611b165768010000b17255775c040260401c5b65ff0000000000821615611c0c5765800000000000821615611b41576801000058b91b5bc9ae0260401c5b65400000000000821615611b5e57680100002c5c89d5ec6d0260401c5b65200000000000821615611b7b5768010000162e43f4f8310260401c5b65100000000000821615611b9857680100000b1721bcfc9a0260401c5b65080000000000821615611bb55768010000058b90cf1e6e0260401c5b65040000000000821615611bd2576801000002c5c863b73f0260401c5b65020000000000821615611bef57680100000162e430e5a20260401c5b65010000000000821615611c0c576801000000b1721835510260401c5b64ff00000000821615611cf957648000000000821615611c3557680100000058b90c0b490260401c5b644000000000821615611c515768010000002c5c8601cc0260401c5b642000000000821615611c6d576801000000162e42fff00260401c5b641000000000821615611c895768010000000b17217fbb0260401c5b640800000000821615611ca5576801000000058b90bfce0260401c5b640400000000821615611cc157680100000002c5c85fe30260401c5b640200000000821615611cdd5768010000000162e42ff10260401c5b640100000000821615611cf957680100000000b17217f80260401c5b63ff000000821615611ddd576380000000821615611d205768010000000058b90bfc0260401c5b6340000000821615611d3b576801000000002c5c85fe0260401c5b6320000000821615611d5657680100000000162e42ff0260401c5b6310000000821615611d71576801000000000b17217f0260401c5b6308000000821615611d8c57680100000000058b90c00260401c5b6304000000821615611da75768010000000002c5c8600260401c5b6302000000821615611dc2576801000000000162e4300260401c5b6301000000821615611ddd5768010000000000b172180260401c5b62ff0000821615611eb85762800000821615611e02576801000000000058b90c0260401c5b62400000821615611e1c57680100000000002c5c860260401c5b62200000821615611e365768010000000000162e430260401c5b62100000821615611e5057680100000000000b17210260401c5b62080000821615611e6a5768010000000000058b910260401c5b62040000821615611e84576801000000000002c5c80260401c5b62020000821615611e9e57680100000000000162e40260401c5b62010000821615611eb8576801000000000000b1720260401c5b61ff00821615611f8a57618000821615611edb57680100000000000058b90260401c5b614000821615611ef45768010000000000002c5d0260401c5b612000821615611f0d576801000000000000162e0260401c5b611000821615611f265768010000000000000b170260401c5b610800821615611f3f576801000000000000058c0260401c5b610400821615611f5857680100000000000002c60260401c5b610200821615611f7157680100000000000001630260401c5b610100821615611f8a57680100000000000000b10260401c5b60ff821615612053576080821615611fab57680100000000000000590260401c5b6040821615611fc3576801000000000000002c0260401c5b6020821615611fdb57680100000000000000160260401c5b6010821615611ff3576801000000000000000b0260401c5b600882161561200b57680100000000000000060260401c5b600482161561202357680100000000000000030260401c5b600282161561203b57680100000000000000010260401c5b600182161561205357680100000000000000010260401c5b670de0b6b3a76400000260409190911c60bf031c90565b80356001600160a01b0381168114612080575f80fd5b919050565b5f60208284031215612095575f80fd5b6116558261206a565b5f805f606084860312156120b0575f80fd5b505081359360208301359350604090920135919050565b5f805f80608085870312156120da575f80fd5b5050823594602084013594506040840135936060013592509050565b5f60208284031215612106575f80fd5b5035919050565b5f806040838503121561211e575f80fd5b6121278361206a565b946020939093013593505050565b602080825282518282018190525f9190848201906040850190845b818110156121755783516001600160a01b031683529284019291840191600101612150565b50909695505050505050565b634e487b7160e01b5f52601160045260245ffd5b8082028115828204841417610d0357610d03612181565b634e487b7160e01b5f52601260045260245ffd5b5f826121da57634e487b7160e01b5f52601260045260245ffd5b500490565b81810381811115610d0357610d03612181565b600181815b8085111561222c57815f190482111561221257612212612181565b8085161561221f57918102915b93841c93908002906121f7565b509250929050565b5f8261224257506001610d03565b8161224e57505f610d03565b8160018114612264576002811461226e5761228a565b6001915050610d03565b60ff84111561227f5761227f612181565b50506001821b610d03565b5060208310610133831016604e8410600b84101617156122ad575081810a610d03565b6122b783836121f2565b805f19048211156122ca576122ca612181565b029392505050565b5f6116558383612234565b5f602082840312156122ed575f80fd5b81518015158114611655575f80fd5b80820180821115610d0357610d03612181565b5f6020828403121561231f575f80fd5b5051919050565b5f6020808352835180828501525f5b8181101561235157858101830151858201604001528201612335565b505f604082860101526040601f19601f8301168501019250505092915050565b5f60808284031215612381575f80fd5b6040516080810181811067ffffffffffffffff821117156123b057634e487b7160e01b5f52604160045260245ffd5b806040525082518152602083015160208201526040830151604082015260608301516060820152809150509291505056fea264697066735822122046771e55cda039f88fda83aeced3ba0e5f684493f60bb2161075b27661e1c5c064736f6c63430008150033",
}

// StabilizationABI is the input ABI used to generate the binding from.
// Deprecated: Use StabilizationMetaData.ABI instead.
var StabilizationABI = StabilizationMetaData.ABI

// Deprecated: Use StabilizationMetaData.Sigs instead.
// StabilizationFuncSigs maps the 4-byte function signature to its string representation.
var StabilizationFuncSigs = StabilizationMetaData.Sigs

// StabilizationBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StabilizationMetaData.Bin instead.
var StabilizationBin = StabilizationMetaData.Bin

// DeployStabilization deploys a new Ethereum contract, binding an instance of Stabilization to it.
func (r *runner) deployStabilization(opts *runOptions, config_ StabilizationConfig, autonity common.Address, operator common.Address, oracle common.Address, supplyControl common.Address, collateralToken common.Address) (common.Address, uint64, *Stabilization, error) {
	parsed, err := StabilizationMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(StabilizationBin), config_, autonity, operator, oracle, supplyControl, collateralToken)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &Stabilization{StabilizationCaller: StabilizationCaller{contract: c}, StabilizationTransactor: StabilizationTransactor{contract: c}, StabilizationFilterer: StabilizationFilterer{contract: c}}, nil
}

// Stabilization is an auto generated Go binding around an Ethereum contract.
type Stabilization struct {
	StabilizationCaller     // Read-only binding to the contract
	StabilizationTransactor // Write-only binding to the contract
	StabilizationFilterer   // Log filterer for contract events
}

// StabilizationCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type StabilizationCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// StabilizationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StabilizationTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// StabilizationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StabilizationFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// SCALE is a free data retrieval call binding the contract method 0xeced5526.
//
// Solidity: function SCALE() view returns(uint256)
func (_Stabilization *StabilizationCaller) SCALE(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "SCALE")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// SCALEFACTOR is a free data retrieval call binding the contract method 0xce4b5bbe.
//
// Solidity: function SCALE_FACTOR() view returns(uint256)
func (_Stabilization *StabilizationCaller) SCALEFACTOR(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "SCALE_FACTOR")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// SECONDSINYEAR is a free data retrieval call binding the contract method 0x5dcc9391.
//
// Solidity: function SECONDS_IN_YEAR() view returns(uint256)
func (_Stabilization *StabilizationCaller) SECONDSINYEAR(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "SECONDS_IN_YEAR")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Accounts is a free data retrieval call binding the contract method 0x68cd03f6.
//
// Solidity: function accounts() view returns(address[])
func (_Stabilization *StabilizationCaller) Accounts(opts *runOptions) ([]common.Address, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "accounts")

	if err != nil {
		return *new([]common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	return out0, consumed, err

}

// BorrowLimit is a free data retrieval call binding the contract method 0x83baa174.
//
// Solidity: function borrowLimit(uint256 collateral, uint256 price, uint256 targetPrice, uint256 mcr) pure returns(uint256)
func (_Stabilization *StabilizationCaller) BorrowLimit(opts *runOptions, collateral *big.Int, price *big.Int, targetPrice *big.Int, mcr *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "borrowLimit", collateral, price, targetPrice, mcr)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Cdps is a free data retrieval call binding the contract method 0x840c7e24.
//
// Solidity: function cdps(address ) view returns(uint256 timestamp, uint256 collateral, uint256 principal, uint256 interest)
func (_Stabilization *StabilizationCaller) Cdps(opts *runOptions, arg0 common.Address) (struct {
	Timestamp  *big.Int
	Collateral *big.Int
	Principal  *big.Int
	Interest   *big.Int
}, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "cdps", arg0)

	outstruct := new(struct {
		Timestamp  *big.Int
		Collateral *big.Int
		Principal  *big.Int
		Interest   *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Collateral = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Principal = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Interest = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// CollateralPrice is a free data retrieval call binding the contract method 0x5891de72.
//
// Solidity: function collateralPrice() view returns(uint256 price)
func (_Stabilization *StabilizationCaller) CollateralPrice(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "collateralPrice")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(uint256 borrowInterestRate, uint256 liquidationRatio, uint256 minCollateralizationRatio, uint256 minDebtRequirement, uint256 targetPrice)
func (_Stabilization *StabilizationCaller) Config(opts *runOptions) (struct {
	BorrowInterestRate        *big.Int
	LiquidationRatio          *big.Int
	MinCollateralizationRatio *big.Int
	MinDebtRequirement        *big.Int
	TargetPrice               *big.Int
}, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "config")

	outstruct := new(struct {
		BorrowInterestRate        *big.Int
		LiquidationRatio          *big.Int
		MinCollateralizationRatio *big.Int
		MinDebtRequirement        *big.Int
		TargetPrice               *big.Int
	})
	if err != nil {
		return *outstruct, consumed, err
	}

	outstruct.BorrowInterestRate = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LiquidationRatio = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MinCollateralizationRatio = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.MinDebtRequirement = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TargetPrice = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// DebtAmount is a free data retrieval call binding the contract method 0x50bf06bf.
//
// Solidity: function debtAmount(address account, uint256 timestamp) view returns(uint256 debt)
func (_Stabilization *StabilizationCaller) DebtAmount(opts *runOptions, account common.Address, timestamp *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "debtAmount", account, timestamp)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// DebtAmount0 is a free data retrieval call binding the contract method 0x54a9f42c.
//
// Solidity: function debtAmount(address account) view returns(uint256 debt)
func (_Stabilization *StabilizationCaller) DebtAmount0(opts *runOptions, account common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "debtAmount0", account)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// InterestDue is a free data retrieval call binding the contract method 0x15184245.
//
// Solidity: function interestDue(uint256 debt, uint256 rate, uint256 timeBorrow, uint256 timeDue) pure returns(uint256)
func (_Stabilization *StabilizationCaller) InterestDue(opts *runOptions, debt *big.Int, rate *big.Int, timeBorrow *big.Int, timeDue *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "interestDue", debt, rate, timeBorrow, timeDue)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// IsLiquidatable is a free data retrieval call binding the contract method 0x042e02cf.
//
// Solidity: function isLiquidatable(address account) view returns(bool)
func (_Stabilization *StabilizationCaller) IsLiquidatable(opts *runOptions, account common.Address) (bool, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "isLiquidatable", account)

	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// MinimumCollateral is a free data retrieval call binding the contract method 0x08796696.
//
// Solidity: function minimumCollateral(uint256 principal, uint256 price, uint256 mcr) pure returns(uint256)
func (_Stabilization *StabilizationCaller) MinimumCollateral(opts *runOptions, principal *big.Int, price *big.Int, mcr *big.Int) (*big.Int, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "minimumCollateral", principal, price, mcr)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UnderCollateralized is a free data retrieval call binding the contract method 0xfbbe6991.
//
// Solidity: function underCollateralized(uint256 collateral, uint256 price, uint256 debt, uint256 liquidationRatio) pure returns(bool)
func (_Stabilization *StabilizationCaller) UnderCollateralized(opts *runOptions, collateral *big.Int, price *big.Int, debt *big.Int, liquidationRatio *big.Int) (bool, uint64, error) {
	out, consumed, err := _Stabilization.call(opts, "underCollateralized", collateral, price, debt, liquidationRatio)

	if err != nil {
		return *new(bool), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, consumed, err

}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_Stabilization *StabilizationTransactor) Borrow(opts *runOptions, amount *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "borrow", amount)
	return consumed, err
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Stabilization *StabilizationTransactor) Deposit(opts *runOptions, amount *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "deposit", amount)
	return consumed, err
}

// Liquidate is a paid mutator transaction binding the contract method 0x2f865568.
//
// Solidity: function liquidate(address account) payable returns()
func (_Stabilization *StabilizationTransactor) Liquidate(opts *runOptions, account common.Address) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "liquidate", account)
	return consumed, err
}

// Repay is a paid mutator transaction binding the contract method 0x402d8883.
//
// Solidity: function repay() payable returns()
func (_Stabilization *StabilizationTransactor) Repay(opts *runOptions) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "repay")
	return consumed, err
}

// SetLiquidationRatio is a paid mutator transaction binding the contract method 0x946ce8cd.
//
// Solidity: function setLiquidationRatio(uint256 ratio) returns()
func (_Stabilization *StabilizationTransactor) SetLiquidationRatio(opts *runOptions, ratio *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setLiquidationRatio", ratio)
	return consumed, err
}

// SetMinCollateralizationRatio is a paid mutator transaction binding the contract method 0x7b44646a.
//
// Solidity: function setMinCollateralizationRatio(uint256 ratio) returns()
func (_Stabilization *StabilizationTransactor) SetMinCollateralizationRatio(opts *runOptions, ratio *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setMinCollateralizationRatio", ratio)
	return consumed, err
}

// SetMinDebtRequirement is a paid mutator transaction binding the contract method 0x53afe81d.
//
// Solidity: function setMinDebtRequirement(uint256 amount) returns()
func (_Stabilization *StabilizationTransactor) SetMinDebtRequirement(opts *runOptions, amount *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setMinDebtRequirement", amount)
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_Stabilization *StabilizationTransactor) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setOperator", operator)
	return consumed, err
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_Stabilization *StabilizationTransactor) SetOracle(opts *runOptions, oracle common.Address) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setOracle", oracle)
	return consumed, err
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_Stabilization *StabilizationTransactor) SetSupplyControl(opts *runOptions, supplyControl common.Address) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setSupplyControl", supplyControl)
	return consumed, err
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_Stabilization *StabilizationTransactor) Withdraw(opts *runOptions, amount *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "withdraw", amount)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// StabilizationBorrowIterator is returned from FilterBorrow and is used to iterate over the raw logs and unpacked data for Borrow events raised by the Stabilization contract.
		type StabilizationBorrowIterator struct {
			Event *StabilizationBorrow // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *StabilizationBorrowIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StabilizationBorrow)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(StabilizationBorrow)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *StabilizationBorrowIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StabilizationBorrowIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StabilizationBorrow represents a Borrow event raised by the Stabilization contract.
		type StabilizationBorrow struct {
			Account common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterBorrow is a free log retrieval operation binding the contract event 0xcbc04eca7e9da35cb1393a6135a199ca52e450d5e9251cbd99f7847d33a36750.
		//
		// Solidity: event Borrow(address indexed account, uint256 amount)
 		func (_Stabilization *StabilizationFilterer) FilterBorrow(opts *bind.FilterOpts, account []common.Address) (*StabilizationBorrowIterator, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.FilterLogs(opts, "Borrow", accountRule)
			if err != nil {
				return nil, err
			}
			return &StabilizationBorrowIterator{contract: _Stabilization.contract, event: "Borrow", logs: logs, sub: sub}, nil
 		}

		// WatchBorrow is a free log subscription operation binding the contract event 0xcbc04eca7e9da35cb1393a6135a199ca52e450d5e9251cbd99f7847d33a36750.
		//
		// Solidity: event Borrow(address indexed account, uint256 amount)
		func (_Stabilization *StabilizationFilterer) WatchBorrow(opts *bind.WatchOpts, sink chan<- *StabilizationBorrow, account []common.Address) (event.Subscription, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.WatchLogs(opts, "Borrow", accountRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StabilizationBorrow)
						if err := _Stabilization.contract.UnpackLog(event, "Borrow", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseBorrow is a log parse operation binding the contract event 0xcbc04eca7e9da35cb1393a6135a199ca52e450d5e9251cbd99f7847d33a36750.
		//
		// Solidity: event Borrow(address indexed account, uint256 amount)
		func (_Stabilization *StabilizationFilterer) ParseBorrow(log types.Log) (*StabilizationBorrow, error) {
			event := new(StabilizationBorrow)
			if err := _Stabilization.contract.UnpackLog(event, "Borrow", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// StabilizationDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Stabilization contract.
		type StabilizationDepositIterator struct {
			Event *StabilizationDeposit // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *StabilizationDepositIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StabilizationDeposit)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(StabilizationDeposit)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *StabilizationDepositIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StabilizationDepositIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StabilizationDeposit represents a Deposit event raised by the Stabilization contract.
		type StabilizationDeposit struct {
			Account common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
		//
		// Solidity: event Deposit(address indexed account, uint256 amount)
 		func (_Stabilization *StabilizationFilterer) FilterDeposit(opts *bind.FilterOpts, account []common.Address) (*StabilizationDepositIterator, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.FilterLogs(opts, "Deposit", accountRule)
			if err != nil {
				return nil, err
			}
			return &StabilizationDepositIterator{contract: _Stabilization.contract, event: "Deposit", logs: logs, sub: sub}, nil
 		}

		// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
		//
		// Solidity: event Deposit(address indexed account, uint256 amount)
		func (_Stabilization *StabilizationFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *StabilizationDeposit, account []common.Address) (event.Subscription, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.WatchLogs(opts, "Deposit", accountRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StabilizationDeposit)
						if err := _Stabilization.contract.UnpackLog(event, "Deposit", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
		//
		// Solidity: event Deposit(address indexed account, uint256 amount)
		func (_Stabilization *StabilizationFilterer) ParseDeposit(log types.Log) (*StabilizationDeposit, error) {
			event := new(StabilizationDeposit)
			if err := _Stabilization.contract.UnpackLog(event, "Deposit", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// StabilizationLiquidateIterator is returned from FilterLiquidate and is used to iterate over the raw logs and unpacked data for Liquidate events raised by the Stabilization contract.
		type StabilizationLiquidateIterator struct {
			Event *StabilizationLiquidate // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *StabilizationLiquidateIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StabilizationLiquidate)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(StabilizationLiquidate)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *StabilizationLiquidateIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StabilizationLiquidateIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StabilizationLiquidate represents a Liquidate event raised by the Stabilization contract.
		type StabilizationLiquidate struct {
			Account common.Address;
			Liquidator common.Address;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterLiquidate is a free log retrieval operation binding the contract event 0xc3d81b2125598b9a2b024afe09e33981f0aa5b7bcbe3e30c4303a4dec209ddb4.
		//
		// Solidity: event Liquidate(address indexed account, address liquidator)
 		func (_Stabilization *StabilizationFilterer) FilterLiquidate(opts *bind.FilterOpts, account []common.Address) (*StabilizationLiquidateIterator, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.FilterLogs(opts, "Liquidate", accountRule)
			if err != nil {
				return nil, err
			}
			return &StabilizationLiquidateIterator{contract: _Stabilization.contract, event: "Liquidate", logs: logs, sub: sub}, nil
 		}

		// WatchLiquidate is a free log subscription operation binding the contract event 0xc3d81b2125598b9a2b024afe09e33981f0aa5b7bcbe3e30c4303a4dec209ddb4.
		//
		// Solidity: event Liquidate(address indexed account, address liquidator)
		func (_Stabilization *StabilizationFilterer) WatchLiquidate(opts *bind.WatchOpts, sink chan<- *StabilizationLiquidate, account []common.Address) (event.Subscription, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.WatchLogs(opts, "Liquidate", accountRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StabilizationLiquidate)
						if err := _Stabilization.contract.UnpackLog(event, "Liquidate", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseLiquidate is a log parse operation binding the contract event 0xc3d81b2125598b9a2b024afe09e33981f0aa5b7bcbe3e30c4303a4dec209ddb4.
		//
		// Solidity: event Liquidate(address indexed account, address liquidator)
		func (_Stabilization *StabilizationFilterer) ParseLiquidate(log types.Log) (*StabilizationLiquidate, error) {
			event := new(StabilizationLiquidate)
			if err := _Stabilization.contract.UnpackLog(event, "Liquidate", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// StabilizationRepayIterator is returned from FilterRepay and is used to iterate over the raw logs and unpacked data for Repay events raised by the Stabilization contract.
		type StabilizationRepayIterator struct {
			Event *StabilizationRepay // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *StabilizationRepayIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StabilizationRepay)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(StabilizationRepay)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *StabilizationRepayIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StabilizationRepayIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StabilizationRepay represents a Repay event raised by the Stabilization contract.
		type StabilizationRepay struct {
			Account common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterRepay is a free log retrieval operation binding the contract event 0x5c16de4f8b59bd9caf0f49a545f25819a895ed223294290b408242e72a594231.
		//
		// Solidity: event Repay(address indexed account, uint256 amount)
 		func (_Stabilization *StabilizationFilterer) FilterRepay(opts *bind.FilterOpts, account []common.Address) (*StabilizationRepayIterator, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.FilterLogs(opts, "Repay", accountRule)
			if err != nil {
				return nil, err
			}
			return &StabilizationRepayIterator{contract: _Stabilization.contract, event: "Repay", logs: logs, sub: sub}, nil
 		}

		// WatchRepay is a free log subscription operation binding the contract event 0x5c16de4f8b59bd9caf0f49a545f25819a895ed223294290b408242e72a594231.
		//
		// Solidity: event Repay(address indexed account, uint256 amount)
		func (_Stabilization *StabilizationFilterer) WatchRepay(opts *bind.WatchOpts, sink chan<- *StabilizationRepay, account []common.Address) (event.Subscription, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.WatchLogs(opts, "Repay", accountRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StabilizationRepay)
						if err := _Stabilization.contract.UnpackLog(event, "Repay", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseRepay is a log parse operation binding the contract event 0x5c16de4f8b59bd9caf0f49a545f25819a895ed223294290b408242e72a594231.
		//
		// Solidity: event Repay(address indexed account, uint256 amount)
		func (_Stabilization *StabilizationFilterer) ParseRepay(log types.Log) (*StabilizationRepay, error) {
			event := new(StabilizationRepay)
			if err := _Stabilization.contract.UnpackLog(event, "Repay", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// StabilizationWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Stabilization contract.
		type StabilizationWithdrawIterator struct {
			Event *StabilizationWithdraw // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *StabilizationWithdrawIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StabilizationWithdraw)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(StabilizationWithdraw)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *StabilizationWithdrawIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StabilizationWithdrawIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StabilizationWithdraw represents a Withdraw event raised by the Stabilization contract.
		type StabilizationWithdraw struct {
			Account common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterWithdraw is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
		//
		// Solidity: event Withdraw(address indexed account, uint256 amount)
 		func (_Stabilization *StabilizationFilterer) FilterWithdraw(opts *bind.FilterOpts, account []common.Address) (*StabilizationWithdrawIterator, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.FilterLogs(opts, "Withdraw", accountRule)
			if err != nil {
				return nil, err
			}
			return &StabilizationWithdrawIterator{contract: _Stabilization.contract, event: "Withdraw", logs: logs, sub: sub}, nil
 		}

		// WatchWithdraw is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
		//
		// Solidity: event Withdraw(address indexed account, uint256 amount)
		func (_Stabilization *StabilizationFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *StabilizationWithdraw, account []common.Address) (event.Subscription, error) {

			var accountRule []interface{}
			for _, accountItem := range account {
				accountRule = append(accountRule, accountItem)
			}


			logs, sub, err := _Stabilization.contract.WatchLogs(opts, "Withdraw", accountRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StabilizationWithdraw)
						if err := _Stabilization.contract.UnpackLog(event, "Withdraw", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseWithdraw is a log parse operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
		//
		// Solidity: event Withdraw(address indexed account, uint256 amount)
		func (_Stabilization *StabilizationFilterer) ParseWithdraw(log types.Log) (*StabilizationWithdraw, error) {
			event := new(StabilizationWithdraw)
			if err := _Stabilization.contract.UnpackLog(event, "Withdraw", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// SupplyControlMetaData contains all meta data concerning the SupplyControl contract.
var SupplyControlMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stabilizer_\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRecipient\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroValue\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"availableSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stabilizer_\",\"type\":\"address\"}],\"name\":\"setStabilizer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stabilizer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7ecc2b56": "availableSupply()",
		"44df8e70": "burn()",
		"40c10f19": "mint(address,uint256)",
		"b3ab15fb": "setOperator(address)",
		"db7f521a": "setStabilizer(address)",
		"7e47961c": "stabilizer()",
		"18160ddd": "totalSupply()",
	},
	Bin: "0x60806040526040516104f33803806104f3833981016040819052610022916100a2565b345f0361004257604051637c946ed760e01b815260040160405180910390fd5b600280546001600160a01b039485166001600160a01b03199182161790915560038054938516938216939093179092555f8054919093169116179055346001556100e2565b80516001600160a01b038116811461009d575f80fd5b919050565b5f805f606084860312156100b4575f80fd5b6100bd84610087565b92506100cb60208501610087565b91506100d960408501610087565b90509250925092565b610404806100ef5f395ff3fe60806040526004361061006e575f3560e01c80637e47961c1161004c5780637e47961c146100c35780637ecc2b56146100f9578063b3ab15fb1461010b578063db7f521a1461012a575f80fd5b806318160ddd1461007257806340c10f191461009a57806344df8e70146100bb575b5f80fd5b34801561007d575f80fd5b5061008760015481565b6040519081526020015b60405180910390f35b3480156100a5575f80fd5b506100b96100b4366004610386565b610149565b005b6100b9610256565b3480156100ce575f80fd5b505f546100e1906001600160a01b031681565b6040516001600160a01b039091168152602001610091565b348015610104575f80fd5b5047610087565b348015610116575f80fd5b506100b96101253660046103ae565b6102d4565b348015610135575f80fd5b506100b96101443660046103ae565b610320565b5f546001600160a01b03163314610172576040516282b42960e81b815260040160405180910390fd5b6001600160a01b038216158061019457505f546001600160a01b038381169116145b156101b257604051634e46966960e11b815260040160405180910390fd5b8015806101be57504781115b156101dc5760405163162908e360e11b815260040160405180910390fd5b6040516001600160a01b0383169082156108fc029083905f818181858888f1935050505015801561020f573d5f803e3d5ffd5b50604080516001600160a01b0384168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a15050565b345f0361027657604051637c946ed760e01b815260040160405180910390fd5b5f546001600160a01b0316331461029f576040516282b42960e81b815260040160405180910390fd5b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a1565b6002546001600160a01b031633146102fe576040516282b42960e81b815260040160405180910390fd5b600380546001600160a01b0319166001600160a01b0392909216919091179055565b6003546001600160a01b0316331461034a576040516282b42960e81b815260040160405180910390fd5b5f80546001600160a01b0319166001600160a01b0392909216919091179055565b80356001600160a01b0381168114610381575f80fd5b919050565b5f8060408385031215610397575f80fd5b6103a08361036b565b946020939093013593505050565b5f602082840312156103be575f80fd5b6103c78261036b565b939250505056fea26469706673582212200dff1b3da3366e61beb6eef65970a8b8f6042ad26792c809eaeb3628313ab6d764736f6c63430008150033",
}

// SupplyControlABI is the input ABI used to generate the binding from.
// Deprecated: Use SupplyControlMetaData.ABI instead.
var SupplyControlABI = SupplyControlMetaData.ABI

// Deprecated: Use SupplyControlMetaData.Sigs instead.
// SupplyControlFuncSigs maps the 4-byte function signature to its string representation.
var SupplyControlFuncSigs = SupplyControlMetaData.Sigs

// SupplyControlBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SupplyControlMetaData.Bin instead.
var SupplyControlBin = SupplyControlMetaData.Bin

// DeploySupplyControl deploys a new Ethereum contract, binding an instance of SupplyControl to it.
func (r *runner) deploySupplyControl(opts *runOptions, autonity common.Address, operator common.Address, stabilizer_ common.Address) (common.Address, uint64, *SupplyControl, error) {
	parsed, err := SupplyControlMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(SupplyControlBin), autonity, operator, stabilizer_)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &SupplyControl{SupplyControlCaller: SupplyControlCaller{contract: c}, SupplyControlTransactor: SupplyControlTransactor{contract: c}, SupplyControlFilterer: SupplyControlFilterer{contract: c}}, nil
}

// SupplyControl is an auto generated Go binding around an Ethereum contract.
type SupplyControl struct {
	SupplyControlCaller     // Read-only binding to the contract
	SupplyControlTransactor // Write-only binding to the contract
	SupplyControlFilterer   // Log filterer for contract events
}

// SupplyControlCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type SupplyControlCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// SupplyControlTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SupplyControlTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// SupplyControlFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SupplyControlFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_SupplyControl *SupplyControlCaller) AvailableSupply(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _SupplyControl.call(opts, "availableSupply")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Stabilizer is a free data retrieval call binding the contract method 0x7e47961c.
//
// Solidity: function stabilizer() view returns(address)
func (_SupplyControl *SupplyControlCaller) Stabilizer(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _SupplyControl.call(opts, "stabilizer")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SupplyControl *SupplyControlCaller) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _SupplyControl.call(opts, "totalSupply")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_SupplyControl *SupplyControlTransactor) Burn(opts *runOptions) (uint64, error) {
	_, consumed, err := _SupplyControl.call(opts, "burn")
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_SupplyControl *SupplyControlTransactor) Mint(opts *runOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _SupplyControl.call(opts, "mint", recipient, amount)
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_SupplyControl *SupplyControlTransactor) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _SupplyControl.call(opts, "setOperator", operator)
	return consumed, err
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_SupplyControl *SupplyControlTransactor) SetStabilizer(opts *runOptions, stabilizer_ common.Address) (uint64, error) {
	_, consumed, err := _SupplyControl.call(opts, "setStabilizer", stabilizer_)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

		// SupplyControlBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the SupplyControl contract.
		type SupplyControlBurnIterator struct {
			Event *SupplyControlBurn // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *SupplyControlBurnIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(SupplyControlBurn)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(SupplyControlBurn)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *SupplyControlBurnIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *SupplyControlBurnIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// SupplyControlBurn represents a Burn event raised by the SupplyControl contract.
		type SupplyControlBurn struct {
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterBurn is a free log retrieval operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
		//
		// Solidity: event Burn(uint256 amount)
 		func (_SupplyControl *SupplyControlFilterer) FilterBurn(opts *bind.FilterOpts) (*SupplyControlBurnIterator, error) {



			logs, sub, err := _SupplyControl.contract.FilterLogs(opts, "Burn")
			if err != nil {
				return nil, err
			}
			return &SupplyControlBurnIterator{contract: _SupplyControl.contract, event: "Burn", logs: logs, sub: sub}, nil
 		}

		// WatchBurn is a free log subscription operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
		//
		// Solidity: event Burn(uint256 amount)
		func (_SupplyControl *SupplyControlFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *SupplyControlBurn) (event.Subscription, error) {



			logs, sub, err := _SupplyControl.contract.WatchLogs(opts, "Burn")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(SupplyControlBurn)
						if err := _SupplyControl.contract.UnpackLog(event, "Burn", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseBurn is a log parse operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
		//
		// Solidity: event Burn(uint256 amount)
		func (_SupplyControl *SupplyControlFilterer) ParseBurn(log types.Log) (*SupplyControlBurn, error) {
			event := new(SupplyControlBurn)
			if err := _SupplyControl.contract.UnpackLog(event, "Burn", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


		// SupplyControlMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the SupplyControl contract.
		type SupplyControlMintIterator struct {
			Event *SupplyControlMint // Event containing the contract specifics and raw log

			contract *bind.BoundContract // Generic contract to use for unpacking event data
			event    string              // Event name to use for unpacking event data

			logs chan types.Log        // Log channel receiving the found contract events
			sub  ethereum.Subscription // Subscription for errors, completion and termination
			done bool                  // Whether the subscription completed delivering logs
			fail error                 // Occurred error to stop iteration
		}
		// Next advances the iterator to the subsequent event, returning whether there
		// are any more events found. In case of a retrieval or parsing error, false is
		// returned and Error() can be queried for the exact failure.
		func (it *SupplyControlMintIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(SupplyControlMint)
					if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
						it.fail = err
						return false
					}
					it.Event.Raw = log
					return true

				default:
					return false
				}
			}
			// Iterator still in progress, wait for either a data or an error event
			select {
			case log := <-it.logs:
				it.Event = new(SupplyControlMint)
				if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
					it.fail = err
					return false
				}
				it.Event.Raw = log
				return true

			case err := <-it.sub.Err():
				it.done = true
				it.fail = err
				return it.Next()
			}
		}
		// Error returns any retrieval or parsing error occurred during filtering.
		func (it *SupplyControlMintIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *SupplyControlMintIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// SupplyControlMint represents a Mint event raised by the SupplyControl contract.
		type SupplyControlMint struct {
			Recipient common.Address;
			Amount *big.Int;
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterMint is a free log retrieval operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
		//
		// Solidity: event Mint(address recipient, uint256 amount)
 		func (_SupplyControl *SupplyControlFilterer) FilterMint(opts *bind.FilterOpts) (*SupplyControlMintIterator, error) {




			logs, sub, err := _SupplyControl.contract.FilterLogs(opts, "Mint")
			if err != nil {
				return nil, err
			}
			return &SupplyControlMintIterator{contract: _SupplyControl.contract, event: "Mint", logs: logs, sub: sub}, nil
 		}

		// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
		//
		// Solidity: event Mint(address recipient, uint256 amount)
		func (_SupplyControl *SupplyControlFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *SupplyControlMint) (event.Subscription, error) {




			logs, sub, err := _SupplyControl.contract.WatchLogs(opts, "Mint")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(SupplyControlMint)
						if err := _SupplyControl.contract.UnpackLog(event, "Mint", log); err != nil {
							return err
						}
						event.Raw = log

						select {
						case sink <- event:
						case err := <-sub.Err():
							return err
						case <-quit:
							return nil
						}
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			}), nil
		}

		// ParseMint is a log parse operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
		//
		// Solidity: event Mint(address recipient, uint256 amount)
		func (_SupplyControl *SupplyControlFilterer) ParseMint(log types.Log) (*SupplyControlMint, error) {
			event := new(SupplyControlMint)
			if err := _SupplyControl.contract.UnpackLog(event, "Mint", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// TestBaseMetaData contains all meta data concerning the TestBase contract.
var TestBaseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_foo\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Foo\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"bfb4ebcf": "Foo()",
	},
	Bin: "0x608060405234801561000f575f80fd5b506040516103f33803806103f383398101604081905261002e91610054565b5f610039828261019c565b5050610257565b634e487b7160e01b5f52604160045260245ffd5b5f6020808385031215610065575f80fd5b82516001600160401b038082111561007b575f80fd5b818501915085601f83011261008e575f80fd5b8151818111156100a0576100a0610040565b604051601f8201601f19908116603f011681019083821181831017156100c8576100c8610040565b8160405282815288868487010111156100df575f80fd5b5f93505b8284101561010057848401860151818501870152928501926100e3565b5f86848301015280965050505050505092915050565b600181811c9082168061012a57607f821691505b60208210810361014857634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115610197575f81815260208120601f850160051c810160208610156101745750805b601f850160051c820191505b8181101561019357828155600101610180565b5050505b505050565b81516001600160401b038111156101b5576101b5610040565b6101c9816101c38454610116565b8461014e565b602080601f8311600181146101fc575f84156101e55750858301515b5f19600386901b1c1916600185901b178555610193565b5f85815260208120601f198616915b8281101561022a5788860151825594840194600190910190840161020b565b508582101561024757878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61018f806102645f395ff3fe608060405234801561000f575f80fd5b5060043610610029575f3560e01c8063bfb4ebcf1461002d575b5f80fd5b61003561004b565b60405161004291906100d6565b60405180910390f35b5f805461005790610121565b80601f016020809104026020016040519081016040528092919081815260200182805461008390610121565b80156100ce5780601f106100a5576101008083540402835291602001916100ce565b820191905f5260205f20905b8154815290600101906020018083116100b157829003601f168201915b505050505081565b5f6020808352835180828501525f5b81811015610101578581018301518582016040015282016100e5565b505f604082860101526040601f19601f8301168501019250505092915050565b600181811c9082168061013557607f821691505b60208210810361015357634e487b7160e01b5f52602260045260245ffd5b5091905056fea26469706673582212209a99d9f5d9d4b6b63e2672e6456a22f8df59dade50c3793f97a6bc55d17a443e64736f6c63430008150033",
}

// TestBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use TestBaseMetaData.ABI instead.
var TestBaseABI = TestBaseMetaData.ABI

// Deprecated: Use TestBaseMetaData.Sigs instead.
// TestBaseFuncSigs maps the 4-byte function signature to its string representation.
var TestBaseFuncSigs = TestBaseMetaData.Sigs

// TestBaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestBaseMetaData.Bin instead.
var TestBaseBin = TestBaseMetaData.Bin

// DeployTestBase deploys a new Ethereum contract, binding an instance of TestBase to it.
func (r *runner) deployTestBase(opts *runOptions, _foo string) (common.Address, uint64, *TestBase, error) {
	parsed, err := TestBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(TestBaseBin), _foo)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &TestBase{TestBaseCaller: TestBaseCaller{contract: c}, TestBaseTransactor: TestBaseTransactor{contract: c}, TestBaseFilterer: TestBaseFilterer{contract: c}}, nil
}

// TestBase is an auto generated Go binding around an Ethereum contract.
type TestBase struct {
	TestBaseCaller     // Read-only binding to the contract
	TestBaseTransactor // Write-only binding to the contract
	TestBaseFilterer   // Log filterer for contract events
}

// TestBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type TestBaseCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// TestBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestBaseTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// TestBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestBaseFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// Foo is a free data retrieval call binding the contract method 0xbfb4ebcf.
//
// Solidity: function Foo() view returns(string)
func (_TestBase *TestBaseCaller) Foo(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _TestBase.call(opts, "Foo")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

/* EVENTS ARE NOT YET SUPPORTED

 */

// TestUpgradedMetaData contains all meta data concerning the TestUpgraded contract.
var TestUpgradedMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_bar\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_foo\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Bar\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"Foo\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_foo\",\"type\":\"string\"}],\"name\":\"FooBar\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"b0a378b0": "Bar()",
		"bfb4ebcf": "Foo()",
		"1e4f3395": "FooBar(string)",
	},
	Bin: "0x608060405234801561000f575f80fd5b5060405161065938038061065983398101604081905261002e91610105565b805f61003a82826101ea565b506001905061004983826101ea565b5050506102a5565b634e487b7160e01b5f52604160045260245ffd5b5f82601f830112610074575f80fd5b81516001600160401b038082111561008e5761008e610051565b604051601f8301601f19908116603f011681019082821181831017156100b6576100b6610051565b816040528381526020925086838588010111156100d1575f80fd5b5f91505b838210156100f257858201830151818301840152908201906100d5565b5f93810190920192909252949350505050565b5f8060408385031215610116575f80fd5b82516001600160401b038082111561012c575f80fd5b61013886838701610065565b9350602085015191508082111561014d575f80fd5b5061015a85828601610065565b9150509250929050565b600181811c9082168061017857607f821691505b60208210810361019657634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156101e5575f81815260208120601f850160051c810160208610156101c25750805b601f850160051c820191505b818110156101e1578281556001016101ce565b5050505b505050565b81516001600160401b0381111561020357610203610051565b610217816102118454610164565b8461019c565b602080601f83116001811461024a575f84156102335750858301515b5f19600386901b1c1916600185901b1785556101e1565b5f85815260208120601f198616915b8281101561027857888601518255948401946001909101908401610259565b508582101561029557878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b6103a7806102b25f395ff3fe608060405234801561000f575f80fd5b506004361061003f575f3560e01c80631e4f339514610043578063b0a378b014610058578063bfb4ebcf14610076575b5f80fd5b610056610051366004610139565b61007e565b005b61006061008d565b60405161006d91906101e4565b60405180910390f35b610060610119565b5f61008982826102b5565b5050565b6001805461009a9061022f565b80601f01602080910402602001604051908101604052809291908181526020018280546100c69061022f565b80156101115780601f106100e857610100808354040283529160200191610111565b820191905f5260205f20905b8154815290600101906020018083116100f457829003601f168201915b505050505081565b5f805461009a9061022f565b634e487b7160e01b5f52604160045260245ffd5b5f60208284031215610149575f80fd5b813567ffffffffffffffff80821115610160575f80fd5b818401915084601f830112610173575f80fd5b81358181111561018557610185610125565b604051601f8201601f19908116603f011681019083821181831017156101ad576101ad610125565b816040528281528760208487010111156101c5575f80fd5b826020860160208301375f928101602001929092525095945050505050565b5f6020808352835180828501525f5b8181101561020f578581018301518582016040015282016101f3565b505f604082860101526040601f19601f8301168501019250505092915050565b600181811c9082168061024357607f821691505b60208210810361026157634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156102b0575f81815260208120601f850160051c8101602086101561028d5750805b601f850160051c820191505b818110156102ac57828155600101610299565b5050505b505050565b815167ffffffffffffffff8111156102cf576102cf610125565b6102e3816102dd845461022f565b84610267565b602080601f831160018114610316575f84156102ff5750858301515b5f19600386901b1c1916600185901b1785556102ac565b5f85815260208120601f198616915b8281101561034457888601518255948401946001909101908401610325565b508582101561036157878501515f19600388901b60f8161c191681555b5050505050600190811b0190555056fea26469706673582212200430be316b59537ec3f0ff42cc3e527d519fc2950a1a7674b8afcb3b5f3abfa064736f6c63430008150033",
}

// TestUpgradedABI is the input ABI used to generate the binding from.
// Deprecated: Use TestUpgradedMetaData.ABI instead.
var TestUpgradedABI = TestUpgradedMetaData.ABI

// Deprecated: Use TestUpgradedMetaData.Sigs instead.
// TestUpgradedFuncSigs maps the 4-byte function signature to its string representation.
var TestUpgradedFuncSigs = TestUpgradedMetaData.Sigs

// TestUpgradedBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestUpgradedMetaData.Bin instead.
var TestUpgradedBin = TestUpgradedMetaData.Bin

// DeployTestUpgraded deploys a new Ethereum contract, binding an instance of TestUpgraded to it.
func (r *runner) deployTestUpgraded(opts *runOptions, _bar string, _foo string) (common.Address, uint64, *TestUpgraded, error) {
	parsed, err := TestUpgradedMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(TestUpgradedBin), _bar, _foo)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &TestUpgraded{TestUpgradedCaller: TestUpgradedCaller{contract: c}, TestUpgradedTransactor: TestUpgradedTransactor{contract: c}, TestUpgradedFilterer: TestUpgradedFilterer{contract: c}}, nil
}

// TestUpgraded is an auto generated Go binding around an Ethereum contract.
type TestUpgraded struct {
	TestUpgradedCaller     // Read-only binding to the contract
	TestUpgradedTransactor // Write-only binding to the contract
	TestUpgradedFilterer   // Log filterer for contract events
}

// TestUpgradedCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type TestUpgradedCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// TestUpgradedTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestUpgradedTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// TestUpgradedFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestUpgradedFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// Bar is a free data retrieval call binding the contract method 0xb0a378b0.
//
// Solidity: function Bar() view returns(string)
func (_TestUpgraded *TestUpgradedCaller) Bar(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _TestUpgraded.call(opts, "Bar")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// Foo is a free data retrieval call binding the contract method 0xbfb4ebcf.
//
// Solidity: function Foo() view returns(string)
func (_TestUpgraded *TestUpgradedCaller) Foo(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _TestUpgraded.call(opts, "Foo")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// FooBar is a paid mutator transaction binding the contract method 0x1e4f3395.
//
// Solidity: function FooBar(string _foo) returns()
func (_TestUpgraded *TestUpgradedTransactor) FooBar(opts *runOptions, _foo string) (uint64, error) {
	_, consumed, err := _TestUpgraded.call(opts, "FooBar", _foo)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

 */

// UpgradeManagerMetaData contains all meta data concerning the UpgradeManager contract.
var UpgradeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_autonity\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"autonity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_data\",\"type\":\"string\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"55463ceb": "autonity()",
		"570ca735": "operator()",
		"b3ab15fb": "setOperator(address)",
		"6e3d9ff0": "upgrade(address,string)",
	},
	Bin: "0x608060405234801561000f575f80fd5b5060405161042438038061042483398101604081905261002e91610079565b5f80546001600160a01b039384166001600160a01b031991821617909155600180549290931691161790556100aa565b80516001600160a01b0381168114610074575f80fd5b919050565b5f806040838503121561008a575f80fd5b6100938361005e565b91506100a16020840161005e565b90509250929050565b61036d806100b75f395ff3fe608060405234801561000f575f80fd5b506004361061004a575f3560e01c806355463ceb1461004e578063570ca7351461007c5780636e3d9ff01461008f578063b3ab15fb146100a4575b5f80fd5b5f54610060906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b600154610060906001600160a01b031681565b6100a261009d366004610213565b6100b7565b005b6100a26100b23660046102cf565b61015d565b6001546001600160a01b031633146101165760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064015b60405180910390fd5b60405160f9905f9061012e90859085906020016102ef565b60405160208183030381529060405290505f80825160208401855af43d5f803e808015610159573d5ff35b3d5ffd5b5f546001600160a01b031633146101c25760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f7420746865204175746f6e69747920636f6e74726044820152621858dd60ea1b606482015260840161010d565b600180546001600160a01b0319166001600160a01b0392909216919091179055565b80356001600160a01b03811681146101fa575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f8060408385031215610224575f80fd5b61022d836101e4565b9150602083013567ffffffffffffffff80821115610249575f80fd5b818501915085601f83011261025c575f80fd5b81358181111561026e5761026e6101ff565b604051601f8201601f19908116603f01168101908382118183101715610296576102966101ff565b816040528281528860208487010111156102ae575f80fd5b826020860160208301375f6020848301015280955050505050509250929050565b5f602082840312156102df575f80fd5b6102e8826101e4565b9392505050565b6bffffffffffffffffffffffff198360601b1681525f82515f5b818110156103265760208186018101516014868401015201610309565b505f9201601401918252509291505056fea2646970667358221220343f4d0feccb061dc53785b45d9ca7b5c2f59b8907eebd003689afdae133576564736f6c63430008150033",
}

// UpgradeManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeManagerMetaData.ABI instead.
var UpgradeManagerABI = UpgradeManagerMetaData.ABI

// Deprecated: Use UpgradeManagerMetaData.Sigs instead.
// UpgradeManagerFuncSigs maps the 4-byte function signature to its string representation.
var UpgradeManagerFuncSigs = UpgradeManagerMetaData.Sigs

// UpgradeManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpgradeManagerMetaData.Bin instead.
var UpgradeManagerBin = UpgradeManagerMetaData.Bin

// DeployUpgradeManager deploys a new Ethereum contract, binding an instance of UpgradeManager to it.
func (r *runner) deployUpgradeManager(opts *runOptions, _autonity common.Address, _operator common.Address) (common.Address, uint64, *UpgradeManager, error) {
	parsed, err := UpgradeManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	if parsed == nil {
		return common.Address{}, 0, nil, errors.New("GetABI returned nil")
	}

	address, gasConsumed, c, err := r.deployContract(opts, parsed, common.FromHex(UpgradeManagerBin), _autonity, _operator)
	if err != nil {
		return common.Address{}, 0, nil, err
	}
	return address, gasConsumed, &UpgradeManager{UpgradeManagerCaller: UpgradeManagerCaller{contract: c}, UpgradeManagerTransactor: UpgradeManagerTransactor{contract: c}, UpgradeManagerFilterer: UpgradeManagerFilterer{contract: c}}, nil
}

// UpgradeManager is an auto generated Go binding around an Ethereum contract.
type UpgradeManager struct {
	UpgradeManagerCaller     // Read-only binding to the contract
	UpgradeManagerTransactor // Write-only binding to the contract
	UpgradeManagerFilterer   // Log filterer for contract events
}

// UpgradeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type UpgradeManagerCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// UpgradeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeManagerTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// UpgradeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeManagerFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// Autonity is a free data retrieval call binding the contract method 0x55463ceb.
//
// Solidity: function autonity() view returns(address)
func (_UpgradeManager *UpgradeManagerCaller) Autonity(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _UpgradeManager.call(opts, "autonity")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_UpgradeManager *UpgradeManagerCaller) Operator(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _UpgradeManager.call(opts, "operator")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager *UpgradeManagerTransactor) SetOperator(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _UpgradeManager.call(opts, "setOperator", _account)
	return consumed, err
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager *UpgradeManagerTransactor) Upgrade(opts *runOptions, _target common.Address, _data string) (uint64, error) {
	_, consumed, err := _UpgradeManager.call(opts, "upgrade", _target, _data)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

 */

// UpgradeableMetaData contains all meta data concerning the Upgradeable contract.
var UpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"872cf059": "completeContractUpgrade()",
		"b66b3e79": "getNewContract()",
		"cf9c5719": "resetContractUpgrade()",
		"b2ea9adb": "upgradeContract(bytes,string)",
	},
}

// UpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeableMetaData.ABI instead.
var UpgradeableABI = UpgradeableMetaData.ABI

// Deprecated: Use UpgradeableMetaData.Sigs instead.
// UpgradeableFuncSigs maps the 4-byte function signature to its string representation.
var UpgradeableFuncSigs = UpgradeableMetaData.Sigs

// Upgradeable is an auto generated Go binding around an Ethereum contract.
type Upgradeable struct {
	UpgradeableCaller     // Read-only binding to the contract
	UpgradeableTransactor // Write-only binding to the contract
	UpgradeableFilterer   // Log filterer for contract events
}

// UpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
// Todo: make it staticall
type UpgradeableCaller struct {
	*contract // Generic contract wrapper for the low level calls
}

// UpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeableTransactor struct {
	*contract // Generic contract wrapper for the low level calls
}

// UpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeableFilterer struct {
	*contract // Generic contract wrapper for the low level calls
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Upgradeable *UpgradeableCaller) GetNewContract(opts *runOptions) ([]byte, string, uint64, error) {
	out, consumed, err := _Upgradeable.call(opts, "getNewContract")

	if err != nil {
		return *new([]byte), *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	return out0, out1, consumed, err

}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactor) CompleteContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _Upgradeable.call(opts, "completeContractUpgrade")
	return consumed, err
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactor) ResetContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _Upgradeable.call(opts, "resetContractUpgrade")
	return consumed, err
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Upgradeable *UpgradeableTransactor) UpgradeContract(opts *runOptions, _bytecode []byte, _abi string) (uint64, error) {
	_, consumed, err := _Upgradeable.call(opts, "upgradeContract", _bytecode, _abi)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

 */
