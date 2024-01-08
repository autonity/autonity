// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package autonity

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
	Block          *big.Int
	Epoch          *big.Int
	ReportingBlock *big.Int
	MessageHash    *big.Int
}

// AutonityCommitteeMember is an auto generated low-level Go binding around an user-defined struct.
type AutonityCommitteeMember struct {
	Addr        common.Address
	VotingPower *big.Int
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
	AccountabilityContract common.Address
	OracleContract         common.Address
	AcuContract            common.Address
	SupplyControlContract  common.Address
	StabilizationContract  common.Address
}

// AutonityPolicy is an auto generated low-level Go binding around an user-defined struct.
type AutonityPolicy struct {
	TreasuryFee     *big.Int
	MinBaseFee      *big.Int
	DelegationRate  *big.Int
	UnbondingPeriod *big.Int
	TreasuryAccount common.Address
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
	State                    uint8
}

// IOracleRoundData is an auto generated low-level Go binding around an user-defined struct.
type IOracleRoundData struct {
	Round     *big.Int
	Price     *big.Int
	Timestamp *big.Int
	Status    *big.Int
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
	Bin: "0x608060405234801562000010575f80fd5b50604051620017fb380380620017fb833981016040819052620000339162000353565b858580518251146200005857604051634ff799c560e01b815260040160405180910390fd5b5f5b8151811015620000b2576001600160ff1b0382828151811062000081576200008162000505565b60200260200101511115620000a957604051634ff799c560e01b815260040160405180910390fd5b6001016200005a565b508751620000c89060039060208b01906200013d565b508651620000de9060049060208a019062000198565b506001869055620000f186600a62000628565b6002555050600680546001600160a01b039485166001600160a01b0319918216179091556007805493851693821693909317909255600880549190931691161790555062000792915050565b828054828255905f5260205f2090810192821562000186579160200282015b82811115620001865782518290620001759082620006c6565b50916020019190600101906200015c565b5062000194929150620001e2565b5090565b828054828255905f5260205f20908101928215620001d4579160200282015b82811115620001d4578251825591602001919060010190620001b7565b506200019492915062000202565b8082111562000194575f620001f8828262000218565b50600101620001e2565b5b8082111562000194575f815560010162000203565b50805462000226906200063c565b5f825580601f1062000236575050565b601f0160209004905f5260205f209081019062000254919062000202565b50565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f191681016001600160401b038111828210171562000296576200029662000257565b604052919050565b5f6001600160401b03821115620002b957620002b962000257565b5060051b60200190565b5f82601f830112620002d3575f80fd5b81516020620002ec620002e6836200029e565b6200026b565b8083825260208201915060208460051b8701019350868411156200030e575f80fd5b602086015b848110156200032c578051835291830191830162000313565b509695505050505050565b80516001600160a01b03811681146200034e575f80fd5b919050565b5f805f805f8060c0878903121562000369575f80fd5b86516001600160401b038111156200037f575f80fd5b8701601f8101891362000390575f80fd5b8051620003a1620002e6826200029e565b808282526020820191508b60208460051b8601011115620003c0575f80fd5b602084015b60208460051b860101811015620004925780516001600160401b03811115620003ec575f80fd5b8d603f8288010112620003fd575f80fd5b858101602001516001600160401b038111156200041e576200041e62000257565b62000433601f8201601f19166020016200026b565b8181528f604083858b01010111156200044a575f80fd5b5f5b828110156200047157604081858b01010151602082840101526020810190506200044c565b505f60208383010152808652505050602083019250602081019050620003c5565b5060208b0151909950925050506001600160401b03811115620004b3575f80fd5b620004c189828a01620002c3565b95505060408701519350620004d96060880162000337565b9250620004e96080880162000337565b9150620004f960a0880162000337565b90509295509295509295565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b600181815b808511156200056d57815f190482111562000551576200055162000519565b808516156200055f57918102915b93841c939080029062000532565b509250929050565b5f82620005855750600162000622565b816200059357505f62000622565b8160018114620005ac5760028114620005b757620005d7565b600191505062000622565b60ff841115620005cb57620005cb62000519565b50506001821b62000622565b5060208310610133831016604e8410600b8410161715620005fc575081810a62000622565b6200060883836200052d565b805f19048211156200061e576200061e62000519565b0290505b92915050565b5f62000635838362000575565b9392505050565b600181811c908216806200065157607f821691505b6020821081036200067057634e487b7160e01b5f52602260045260245ffd5b50919050565b601f821115620006c157805f5260205f20601f840160051c810160208510156200069d5750805b601f840160051c820191505b81811015620006be575f8155600101620006a9565b50505b505050565b81516001600160401b03811115620006e257620006e262000257565b620006fa81620006f384546200063c565b8462000676565b602080601f83116001811462000730575f8415620007185750858301515b5f19600386901b1c1916600185901b1785556200078a565b5f85815260208120601f198616915b8281101562000760578886015182559484019460019091019084016200073f565b50858210156200077e57878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b61105b80620007a05f395ff3fe608060405234801561000f575f80fd5b506004361061009b575f3560e01c80637adbf973116100635780637adbf973146100f9578063a2e620451461010c578063b3ab15fb14610124578063d54d279914610137578063f51e181a1461014c575f80fd5b806307039ff91461009f578063146ca531146100bd5780633fa4f245146100d357806344b4708a146100db578063683dd191146100f0575b5f80fd5b6100a7610155565b6040516100b491906108f4565b60405180910390f35b6100c55f5481565b6040519081526020016100b4565b6100c5610229565b6100ee6100e93660046109e1565b610252565b005b6100c560025481565b6100ee610107366004610b05565b610373565b6101146103bf565b60405190151581526020016100b4565b6100ee610132366004610b05565b6106cc565b61013f610718565b6040516100b49190610b65565b6100c560015481565b60606003805480602002602001604051908101604052809291908181526020015f905b82821015610220578382905f5260205f2001805461019590610b77565b80601f01602080910402602001604051908101604052809291908181526020018280546101c190610b77565b801561020c5780601f106101e35761010080835404028352916020019161020c565b820191905f5260205f20905b8154815290600101906020018083116101ef57829003601f168201915b505050505081526020019060010190610178565b50505050905090565b5f80545f0361024b57604051631d3e00bb60e11b815260040160405180910390fd5b5060055490565b8282805182511461027657604051634ff799c560e01b815260040160405180910390fd5b5f5b81518110156102ca576001600160ff1b0382828151811061029b5761029b610baf565b602002602001015111156102c257604051634ff799c560e01b815260040160405180910390fd5b600101610278565b506007546001600160a01b031633146102f5576040516282b42960e81b815260040160405180910390fd5b845161030890600390602088019061076e565b50835161031c9060049060208701906107c2565b50600183905561032d83600a610cb9565b6002556040517fdbdcd10543a20811a4a332247f28d03b22686d3281043de35824a06075c06c099061036490879087908790610cc4565b60405180910390a15050505050565b6006546001600160a01b0316331461039d576040516282b42960e81b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0392909216919091179055565b6006545f906001600160a01b031633146103eb576040516282b42960e81b815260040160405180910390fd5b5f600160085f9054906101000a90046001600160a01b03166001600160a01b0316639f8743f76040518163ffffffff1660e01b8152600401602060405180830381865afa15801561043e573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104629190610cf9565b61046c9190610d10565b9050805f541061047d575f91505090565b5f8060085f9054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156104cf573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104f39190610cf9565b90505f5b60035481101561066857604051661554d10b5554d160ca1b60208201525f90602701604051602081830303815290604052805190602001206003838154811061054257610542610baf565b905f5260205f200160405160200161055a9190610d23565b604051602081830303815290604052805190602001200361057c57508161062c565b600854600380545f926001600160a01b031691633c8510fd91899190879081106105a8576105a8610baf565b905f5260205f20016040518363ffffffff1660e01b81526004016105cd929190610d95565b608060405180830381865afa1580156105e8573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061060c9190610e27565b905080606001515f14610625575f965050505050505090565b6020015190505b6004828154811061063f5761063f610baf565b905f5260205f200154816106539190610e8b565b61065d9085610eba565b9350506001016104f7565b506106738183610ee1565b60058190555f849055604080514381524260208201528082018690526060810192909252517f23f161ca67071b3e902d4fa7afade82672c6160677e89d373a830145bdda6d269181900360800190a16001935050505090565b6006546001600160a01b031633146106f6576040516282b42960e81b815260040160405180910390fd5b600780546001600160a01b0319166001600160a01b0392909216919091179055565b6060600480548060200260200160405190810160405280929190818152602001828054801561076457602002820191905f5260205f20905b815481526020019060010190808311610750575b5050505050905090565b828054828255905f5260205f209081019282156107b2579160200282015b828111156107b257825182906107a29082610f65565b509160200191906001019061078c565b506107be929150610807565b5090565b828054828255905f5260205f209081019282156107fb579160200282015b828111156107fb5782518255916020019190600101906107e0565b506107be929150610823565b808211156107be575f61081a8282610837565b50600101610807565b5b808211156107be575f8155600101610824565b50805461084390610b77565b5f825580601f10610852575050565b601f0160209004905f5260205f209081019061086e9190610823565b50565b5f82825180855260208086019550808260051b8401018186015f5b848110156108e757601f1980878503018a52825180518086525f5b818110156108c25782810188015187820189015287016108a7565b505f8682018801529a86019a601f01909116909301840192509083019060010161088c565b5090979650505050505050565b602081525f6109066020830184610871565b9392505050565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561094a5761094a61090d565b604052919050565b5f67ffffffffffffffff82111561096b5761096b61090d565b5060051b60200190565b5f82601f830112610984575f80fd5b8135602061099961099483610952565b610921565b8083825260208201915060208460051b8701019350868411156109ba575f80fd5b602086015b848110156109d657803583529183019183016109bf565b509695505050505050565b5f805f606084860312156109f3575f80fd5b833567ffffffffffffffff80821115610a0a575f80fd5b818601915086601f830112610a1d575f80fd5b81356020610a2d61099483610952565b82815260059290921b8401810191818101908a841115610a4b575f80fd5b8286015b84811015610ad157803586811115610a65575f80fd5b8701603f81018d13610a75575f80fd5b84810135604088821115610a8b57610a8b61090d565b610a9d601f8301601f19168801610921565b8281528f82848601011115610ab0575f80fd5b82828501898301375f92810188019290925250845250918301918301610a4f565b5097505087013592505080821115610ae7575f80fd5b50610af486828701610975565b925050604084013590509250925092565b5f60208284031215610b15575f80fd5b81356001600160a01b0381168114610906575f80fd5b5f815180845260208085019450602084015f5b83811015610b5a57815187529582019590820190600101610b3e565b509495945050505050565b602081525f6109066020830184610b2b565b600181811c90821680610b8b57607f821691505b602082108103610ba957634e487b7160e01b5f52602260045260245ffd5b50919050565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b600181815b80851115610c1157815f1904821115610bf757610bf7610bc3565b80851615610c0457918102915b93841c9390800290610bdc565b509250929050565b5f82610c2757506001610cb3565b81610c3357505f610cb3565b8160018114610c495760028114610c5357610c6f565b6001915050610cb3565b60ff841115610c6457610c64610bc3565b50506001821b610cb3565b5060208310610133831016604e8410600b8410161715610c92575081810a610cb3565b610c9c8383610bd7565b805f1904821115610caf57610caf610bc3565b0290505b92915050565b5f6109068383610c19565b606081525f610cd66060830186610871565b8281036020840152610ce88186610b2b565b915050826040830152949350505050565b5f60208284031215610d09575f80fd5b5051919050565b81810381811115610cb357610cb3610bc3565b5f808354610d3081610b77565b60018281168015610d485760018114610d5d57610d89565b60ff1984168752821515830287019450610d89565b875f526020805f205f5b85811015610d805781548a820152908401908201610d67565b50505082870194505b50929695505050505050565b8281525f6020604060208401525f8454610dae81610b77565b806040870152606060018084165f8114610dcf5760018114610deb57610e18565b60ff19851660608a0152606084151560051b8a01019550610e18565b895f5260205f205f5b85811015610e0f5781548b8201860152908301908801610df4565b8a016060019650505b50939998505050505050505050565b5f60808284031215610e37575f80fd5b6040516080810181811067ffffffffffffffff82111715610e5a57610e5a61090d565b8060405250825181526020830151602082015260408301516040820152606083015160608201528091505092915050565b8082025f8212600160ff1b84141615610ea657610ea6610bc3565b8181058314821517610cb357610cb3610bc3565b8082018281125f831280158216821582161715610ed957610ed9610bc3565b505092915050565b5f82610efb57634e487b7160e01b5f52601260045260245ffd5b600160ff1b82145f1984141615610f1457610f14610bc3565b500590565b601f821115610f6057805f5260205f20601f840160051c81016020851015610f3e5750805b601f840160051c820191505b81811015610f5d575f8155600101610f4a565b50505b505050565b815167ffffffffffffffff811115610f7f57610f7f61090d565b610f9381610f8d8454610b77565b84610f19565b602080601f831160018114610fc6575f8415610faf5750858301515b5f19600386901b1c1916600185901b17855561101d565b5f85815260208120601f198616915b82811015610ff457888601518255948401946001909101908401610fd5565b508582101561101157878501515f19600388901b60f8161c191681555b505060018460011b0185555b50505050505056fea264697066735822122099c00422e6fb46c734bdd3cee7b3dd68e35402f188a00485fb31b5a4565929bd64736f6c63430008170033",
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
func DeployACU(auth *bind.TransactOpts, backend bind.ContractBackend, symbols_ []string, quantities_ []*big.Int, scale_ *big.Int, autonity common.Address, operator common.Address, oracle common.Address) (common.Address, *types.Transaction, *ACU, error) {
	parsed, err := ACUMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ACUBin), backend, symbols_, quantities_, scale_, autonity, operator, oracle)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ACU{ACUCaller: ACUCaller{contract: contract}, ACUTransactor: ACUTransactor{contract: contract}, ACUFilterer: ACUFilterer{contract: contract}}, nil
}

// ACU is an auto generated Go binding around an Ethereum contract.
type ACU struct {
	ACUCaller     // Read-only binding to the contract
	ACUTransactor // Write-only binding to the contract
	ACUFilterer   // Log filterer for contract events
}

// ACUCaller is an auto generated read-only Go binding around an Ethereum contract.
type ACUCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ACUTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ACUTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ACUFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ACUFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ACUSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ACUSession struct {
	Contract     *ACU              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ACUCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ACUCallerSession struct {
	Contract *ACUCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ACUTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ACUTransactorSession struct {
	Contract     *ACUTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ACURaw is an auto generated low-level Go binding around an Ethereum contract.
type ACURaw struct {
	Contract *ACU // Generic contract binding to access the raw methods on
}

// ACUCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ACUCallerRaw struct {
	Contract *ACUCaller // Generic read-only contract binding to access the raw methods on
}

// ACUTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ACUTransactorRaw struct {
	Contract *ACUTransactor // Generic write-only contract binding to access the raw methods on
}

// NewACU creates a new instance of ACU, bound to a specific deployed contract.
func NewACU(address common.Address, backend bind.ContractBackend) (*ACU, error) {
	contract, err := bindACU(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ACU{ACUCaller: ACUCaller{contract: contract}, ACUTransactor: ACUTransactor{contract: contract}, ACUFilterer: ACUFilterer{contract: contract}}, nil
}

// NewACUCaller creates a new read-only instance of ACU, bound to a specific deployed contract.
func NewACUCaller(address common.Address, caller bind.ContractCaller) (*ACUCaller, error) {
	contract, err := bindACU(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ACUCaller{contract: contract}, nil
}

// NewACUTransactor creates a new write-only instance of ACU, bound to a specific deployed contract.
func NewACUTransactor(address common.Address, transactor bind.ContractTransactor) (*ACUTransactor, error) {
	contract, err := bindACU(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ACUTransactor{contract: contract}, nil
}

// NewACUFilterer creates a new log filterer instance of ACU, bound to a specific deployed contract.
func NewACUFilterer(address common.Address, filterer bind.ContractFilterer) (*ACUFilterer, error) {
	contract, err := bindACU(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ACUFilterer{contract: contract}, nil
}

// bindACU binds a generic wrapper to an already deployed contract.
func bindACU(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ACUABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ACU *ACURaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ACU.Contract.ACUCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ACU *ACURaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ACU.Contract.ACUTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ACU *ACURaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ACU.Contract.ACUTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ACU *ACUCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ACU.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ACU *ACUTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ACU.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ACU *ACUTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ACU.Contract.contract.Transact(opts, method, params...)
}

// Quantities is a free data retrieval call binding the contract method 0xd54d2799.
//
// Solidity: function quantities() view returns(uint256[])
func (_ACU *ACUCaller) Quantities(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _ACU.contract.Call(opts, &out, "quantities")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// Quantities is a free data retrieval call binding the contract method 0xd54d2799.
//
// Solidity: function quantities() view returns(uint256[])
func (_ACU *ACUSession) Quantities() ([]*big.Int, error) {
	return _ACU.Contract.Quantities(&_ACU.CallOpts)
}

// Quantities is a free data retrieval call binding the contract method 0xd54d2799.
//
// Solidity: function quantities() view returns(uint256[])
func (_ACU *ACUCallerSession) Quantities() ([]*big.Int, error) {
	return _ACU.Contract.Quantities(&_ACU.CallOpts)
}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_ACU *ACUCaller) Round(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ACU.contract.Call(opts, &out, "round")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_ACU *ACUSession) Round() (*big.Int, error) {
	return _ACU.Contract.Round(&_ACU.CallOpts)
}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_ACU *ACUCallerSession) Round() (*big.Int, error) {
	return _ACU.Contract.Round(&_ACU.CallOpts)
}

// Scale is a free data retrieval call binding the contract method 0xf51e181a.
//
// Solidity: function scale() view returns(uint256)
func (_ACU *ACUCaller) Scale(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ACU.contract.Call(opts, &out, "scale")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Scale is a free data retrieval call binding the contract method 0xf51e181a.
//
// Solidity: function scale() view returns(uint256)
func (_ACU *ACUSession) Scale() (*big.Int, error) {
	return _ACU.Contract.Scale(&_ACU.CallOpts)
}

// Scale is a free data retrieval call binding the contract method 0xf51e181a.
//
// Solidity: function scale() view returns(uint256)
func (_ACU *ACUCallerSession) Scale() (*big.Int, error) {
	return _ACU.Contract.Scale(&_ACU.CallOpts)
}

// ScaleFactor is a free data retrieval call binding the contract method 0x683dd191.
//
// Solidity: function scaleFactor() view returns(uint256)
func (_ACU *ACUCaller) ScaleFactor(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ACU.contract.Call(opts, &out, "scaleFactor")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScaleFactor is a free data retrieval call binding the contract method 0x683dd191.
//
// Solidity: function scaleFactor() view returns(uint256)
func (_ACU *ACUSession) ScaleFactor() (*big.Int, error) {
	return _ACU.Contract.ScaleFactor(&_ACU.CallOpts)
}

// ScaleFactor is a free data retrieval call binding the contract method 0x683dd191.
//
// Solidity: function scaleFactor() view returns(uint256)
func (_ACU *ACUCallerSession) ScaleFactor() (*big.Int, error) {
	return _ACU.Contract.ScaleFactor(&_ACU.CallOpts)
}

// Symbols is a free data retrieval call binding the contract method 0x07039ff9.
//
// Solidity: function symbols() view returns(string[])
func (_ACU *ACUCaller) Symbols(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _ACU.contract.Call(opts, &out, "symbols")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// Symbols is a free data retrieval call binding the contract method 0x07039ff9.
//
// Solidity: function symbols() view returns(string[])
func (_ACU *ACUSession) Symbols() ([]string, error) {
	return _ACU.Contract.Symbols(&_ACU.CallOpts)
}

// Symbols is a free data retrieval call binding the contract method 0x07039ff9.
//
// Solidity: function symbols() view returns(string[])
func (_ACU *ACUCallerSession) Symbols() ([]string, error) {
	return _ACU.Contract.Symbols(&_ACU.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(int256)
func (_ACU *ACUCaller) Value(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ACU.contract.Call(opts, &out, "value")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(int256)
func (_ACU *ACUSession) Value() (*big.Int, error) {
	return _ACU.Contract.Value(&_ACU.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(int256)
func (_ACU *ACUCallerSession) Value() (*big.Int, error) {
	return _ACU.Contract.Value(&_ACU.CallOpts)
}

// ModifyBasket is a paid mutator transaction binding the contract method 0x44b4708a.
//
// Solidity: function modifyBasket(string[] symbols_, uint256[] quantities_, uint256 scale_) returns()
func (_ACU *ACUTransactor) ModifyBasket(opts *bind.TransactOpts, symbols_ []string, quantities_ []*big.Int, scale_ *big.Int) (*types.Transaction, error) {
	return _ACU.contract.Transact(opts, "modifyBasket", symbols_, quantities_, scale_)
}

// ModifyBasket is a paid mutator transaction binding the contract method 0x44b4708a.
//
// Solidity: function modifyBasket(string[] symbols_, uint256[] quantities_, uint256 scale_) returns()
func (_ACU *ACUSession) ModifyBasket(symbols_ []string, quantities_ []*big.Int, scale_ *big.Int) (*types.Transaction, error) {
	return _ACU.Contract.ModifyBasket(&_ACU.TransactOpts, symbols_, quantities_, scale_)
}

// ModifyBasket is a paid mutator transaction binding the contract method 0x44b4708a.
//
// Solidity: function modifyBasket(string[] symbols_, uint256[] quantities_, uint256 scale_) returns()
func (_ACU *ACUTransactorSession) ModifyBasket(symbols_ []string, quantities_ []*big.Int, scale_ *big.Int) (*types.Transaction, error) {
	return _ACU.Contract.ModifyBasket(&_ACU.TransactOpts, symbols_, quantities_, scale_)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ACU *ACUTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ACU.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ACU *ACUSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _ACU.Contract.SetOperator(&_ACU.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ACU *ACUTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _ACU.Contract.SetOperator(&_ACU.TransactOpts, operator)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_ACU *ACUTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _ACU.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_ACU *ACUSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _ACU.Contract.SetOracle(&_ACU.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_ACU *ACUTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _ACU.Contract.SetOracle(&_ACU.TransactOpts, oracle)
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_ACU *ACUTransactor) Update(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ACU.contract.Transact(opts, "update")
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_ACU *ACUSession) Update() (*types.Transaction, error) {
	return _ACU.Contract.Update(&_ACU.TransactOpts)
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_ACU *ACUTransactorSession) Update() (*types.Transaction, error) {
	return _ACU.Contract.Update(&_ACU.TransactOpts)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Symbols    []string
	Quantities []*big.Int
	Scale      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Height    *big.Int
	Timestamp *big.Int
	Round     *big.Int
	Value     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
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

// AccountabilityMetaData contains all meta data concerning the Accountability contract.
var AccountabilityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_autonity\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"innocenceProofSubmissionWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateLow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateMid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collusionFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"historyFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashingRatePrecision\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"InnocenceProven\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewAccusation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewFaultProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"releaseBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isJailbound\",\"type\":\"bool\"}],\"name\":\"SlashingEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"beneficiaries\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"_rule\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"canAccuse\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_deadline\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"_rule\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"canSlash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"innocenceProofSubmissionWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateLow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateMid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collusionFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"historyFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashingRatePrecision\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"events\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"chunkId\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_epochEnd\",\"type\":\"bool\"}],\"name\":\"finalize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_val\",\"type\":\"address\"}],\"name\":\"getValidatorAccusation\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"chunkId\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_val\",\"type\":\"address\"}],\"name\":\"getValidatorFaults\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"chunkId\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"chunks\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"chunkId\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event\",\"name\":\"_event\",\"type\":\"tuple\"}],\"name\":\"handleEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newPeriod\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"slashingHistory\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
		"069f6863": "handleEvent((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256))",
		"6b5f444c": "setEpochPeriod(uint256)",
		"e7bb0b52": "slashingHistory(address,uint256)",
	},
	Bin: "0x60806040525f60115534801562000014575f80fd5b5060405162003afb38038062003afb8339810160408190526200003791620000f2565b600180546001600160a01b0319166001600160a01b03841690811790915560408051636fd8d26960e11b8152905163dfb1a4d2916004808201926020929091908290030181865afa1580156200008f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620000b59190620001b4565b5f558051600355602081015160045560408101516005556060810151600655608081015160075560a081015160085560c0015160095550620001cc565b5f8082840361010081121562000106575f80fd5b83516001600160a01b03811681146200011d575f80fd5b925060e0601f198201121562000131575f80fd5b5060405160e081016001600160401b03811182821017156200016157634e487b7160e01b5f52604160045260245ffd5b80604052506020840151815260408401516020820152606084015160408201526080840151606082015260a0840151608082015260c084015160a082015260e084015160c0820152809150509250929050565b5f60208284031215620001c5575f80fd5b5051919050565b61392180620001da5f395ff3fe6080604052600436106100bf575f3560e01c80636c9789b01161007c5780639cb22b06116100575780639cb22b0614610284578063b5b7a184146102b0578063bebaa8fc146102d2578063e7bb0b52146102fe575f80fd5b80636c9789b0146101cc57806379502c55146101eb5780637ccecadd1461024e575f80fd5b806301567739146100c3578063069f6863146101145780630b791430146101355780631de9d9b61461016b5780634108a95a1461017e5780636b5f444c146101ad575b5f80fd5b3480156100ce575f80fd5b506100f76100dd366004612d86565b600a6020525f90815260409020546001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b34801561011f575f80fd5b5061013361012e366004612ee0565b610334565b005b348015610140575f80fd5b5061015461014f366004612fdc565b610691565b60405161010b9b9a99989796959493929190613078565b610133610179366004612d86565b61079a565b348015610189575f80fd5b5061019d6101983660046130fd565b610991565b604051901515815260200161010b565b3480156101b8575f80fd5b506101336101c7366004612fdc565b610a3c565b3480156101d7575f80fd5b506101336101e6366004613138565b610a6a565b3480156101f6575f80fd5b506003546004546005546006546007546008546009546102199695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e00161010b565b348015610259575f80fd5b5061026d6102683660046130fd565b610aad565b60408051921515835260208301919091520161010b565b34801561028f575f80fd5b506102a361029e366004612d86565b610bf3565b60405161010b9190613221565b3480156102bb575f80fd5b506102c45f5481565b60405190815260200161010b565b3480156102dd575f80fd5b506102f16102ec366004612d86565b610e0b565b60405161010b9190613233565b348015610309575f80fd5b506102c4610318366004613295565b600e60209081525f928352604080842090915290825290205481565b600154604051630c825d9760e11b81523360048201525f916001600160a01b031690631904bb2e906024015f60405180830381865afa158015610379573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526103a09190810190613322565b60208101519091506001600160a01b0316331461041a5760405162461bcd60e51b815260206004820152602d60248201527f66756e6374696f6e207265737472696374656420746f2061207265676973746560448201526c3932b2103b30b634b230ba37b960991b60648201526084015b60405180910390fd5b60808201516001600160a01b031633146104765760405162461bcd60e51b815260206004820152601d60248201527f6576656e74207265706f72746572206d7573742062652063616c6c65720000006044820152606401610411565b6001825f015160ff16111561061c575f61048f8361108c565b90508061049b57505050565b335f908152600d6020908152604091829020825161016081018452815460ff808216835261010082048116948301949094529093919291840191620100009091041660028111156104ee576104ee612ff3565b60028111156104ff576104ff612ff3565b815281546020909101906301000000900460ff16600d81111561052457610524612ff3565b600d81111561053557610535612ff3565b815281546001600160a01b03600160201b90910481166020830152600183015416604082015260028201805460609092019161057090613475565b80601f016020809104026020016040519081016040528092919081815260200182805461059c90613475565b80156105e75780601f106105be576101008083540402835291602001916105e7565b820191905f5260205f20905b8154815290600101906020018083116105ca57829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815250509250505b5f8260400151600281111561063357610633612ff3565b0361064557610641826112b2565b5050565b60018260400151600281111561065d5761065d612ff3565b0361066b5761064182611457565b60028260400151600281111561068357610683612ff3565b0361064157610641826115de565b600281815481106106a0575f80fd5b5f91825260209091206007909102018054600182015460028301805460ff8085169650610100850481169562010000860482169563010000008104909216946001600160a01b03600160201b90930483169492169290919061070190613475565b80601f016020809104026020016040519081016040528092919081815260200182805461072d90613475565b80156107785780601f1061074f57610100808354040283529160200191610778565b820191905f5260205f20905b81548152906001019060200180831161075b57829003601f168201915b505050505090806003015490806004015490806005015490806006015490508b565b6001546001600160a01b031633146107c45760405162461bcd60e51b8152600401610411906134ad565b6001546001600160a01b038281165f908152600a6020526040808220549051630c825d9760e11b8152908316600482015290929190911690631904bb2e906024015f60405180830381865afa15801561081f573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526108469190810190613322565b516040519091505f906001600160a01b038316906108fc90349084818181858888f193505050503d805f8114610897576040519150601f19603f3d011682016040523d82523d5f602084013e61089c565b606091505b50509050806109695760015f9054906101000a90046001600160a01b03166001600160a01b031663f7866ee36040518163ffffffff1660e01b8152600401602060405180830381865afa1580156108f5573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061091991906134f1565b6001600160a01b0316346040515f6040518083038185875af1925050503d805f8114610960576040519150601f19603f3d011682016040523d82523d5f602084013e610965565b606091505b5050505b50506001600160a01b03165f908152600a6020526040902080546001600160a01b0319169055565b5f8061099c846116f1565b6001546040516396b477cb60e01b8152600481018690529192505f916001600160a01b03909116906396b477cb90602401602060405180830381865afa1580156109e8573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610a0c919061350c565b6001600160a01b0387165f908152600e602090815260408083209383529290522054919091109150509392505050565b6001546001600160a01b03163314610a665760405162461bcd60e51b8152600401610411906134ad565b5f55565b6001546001600160a01b03163314610a945760405162461bcd60e51b8152600401610411906134ad565b610a9c61175a565b8015610aaa57610aaa611a92565b50565b5f805f610ab9856116f1565b6001546040516396b477cb60e01b8152600481018790529192505f916001600160a01b03909116906396b477cb90602401602060405180830381865afa158015610b05573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610b29919061350c565b6001600160a01b0388165f908152600e602090815260408083208484529091529020549091508211610b60575f93505f9250610be9565b6001600160a01b0387165f908152600c602052604090205415610be1576001600160a01b0387165f908152600c6020526040812054600290610ba490600190613537565b81548110610bb457610bb461354a565b905f5260205f20906007020190505f945060035f01548160030154610bd9919061355e565b935050610be9565b600193505f92505b5050935093915050565b610bfb612cad565b6001600160a01b0382165f908152600c6020526040902054610c4f5760405162461bcd60e51b815260206004820152600d60248201526c37379030b1b1bab9b0ba34b7b760991b6044820152606401610411565b6001600160a01b0382165f908152600c6020526040902054600290610c7690600190613537565b81548110610c8657610c8661354a565b5f91825260209182902060408051610160810182526007909302909101805460ff808216855261010082048116958501959095529293909291840191620100009004166002811115610cda57610cda612ff3565b6002811115610ceb57610ceb612ff3565b815281546020909101906301000000900460ff16600d811115610d1057610d10612ff3565b600d811115610d2157610d21612ff3565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191610d5c90613475565b80601f0160208091040260200160405190810160405280929190818152602001828054610d8890613475565b8015610dd35780601f10610daa57610100808354040283529160200191610dd3565b820191905f5260205f20905b815481529060010190602001808311610db657829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815250509050919050565b6001600160a01b0381165f908152600b60205260408120546060919067ffffffffffffffff811115610e3f57610e3f612da8565b604051908082528060200260200182016040528015610e7857816020015b610e65612cad565b815260200190600190039081610e5d5790505b5090505f5b6001600160a01b0384165f908152600b6020526040902054811015611085576001600160a01b0384165f908152600b6020526040902080546002919083908110610ec957610ec961354a565b905f5260205f20015481548110610ee257610ee261354a565b5f91825260209182902060408051610160810182526007909302909101805460ff808216855261010082048116958501959095529293909291840191620100009004166002811115610f3657610f36612ff3565b6002811115610f4757610f47612ff3565b815281546020909101906301000000900460ff16600d811115610f6c57610f6c612ff3565b600d811115610f7d57610f7d612ff3565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191610fb890613475565b80601f0160208091040260200160405190810160405280929190818152602001828054610fe490613475565b801561102f5780601f106110065761010080835404028352916020019161102f565b820191905f5260205f20905b81548152906001019060200180831161101257829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815250508282815181106110725761107261354a565b6020908102919091010152600101610e7d565b5092915050565b5f816020015160ff165f036111b257335f908152600d6020908152604091829020845181549286015160ff9081166101000261ffff1990941691161791909117808255918401518492829062ff00001916620100008360028111156110f3576110f3612ff3565b021790555060608201518154829063ff0000001916630100000083600d81111561111f5761111f612ff3565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c0820151600282019061117d90826135bd565b5060e082015160038201556101008201516004820155610120820151600582015561014090910151600690910155505f919050565b602080830151335f908152600d90925260409091205460ff918216916111df916101009004166001613679565b60ff161461122f5760405162461bcd60e51b815260206004820152601960248201527f6368756e6b73206d75737420626520636f6e746967756f7573000000000000006044820152606401610411565b335f908152600d6020526040902060c083015161124f9160020190611d46565b335f908152600d602052604090208054600191908290611278908290610100900460ff16613679565b92506101000a81548160ff021916908360ff160217905550815f015160ff16826020015160016112a89190613679565b60ff161492915050565b5f805f805f6112c660fe8760c00151611e8c565b945094509450945094508461131d5760405162461bcd60e51b815260206004820152601960248201527f6661696c65642070726f6f6620766572696669636174696f6e000000000000006044820152606401610411565b8560a001516001600160a01b0316846001600160a01b0316146113525760405162461bcd60e51b815260040161041190613692565b8560600151600d81111561136857611368612ff3565b83146113865760405162461bcd60e51b8152600401610411906136bd565b4382106113a55760405162461bcd60e51b8152600401610411906136e7565b815f036113ba576113b7600143613537565b91505b6001546040516396b477cb60e01b8152600481018490525f916001600160a01b0316906396b477cb90602401602060405180830381865afa158015611401573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611425919061350c565b60e088018490526101008801819052436101208901526101408801839052905061144e87611eef565b50505050505050565b5f805f805f61146b60fc8760c00151611e8c565b94509450945094509450846114c25760405162461bcd60e51b815260206004820152601e60248201527f6661696c65642061636375736174696f6e20766572696669636174696f6e00006044820152606401610411565b8560a001516001600160a01b0316846001600160a01b0316146114f75760405162461bcd60e51b815260040161041190613692565b8560600151600d81111561150d5761150d612ff3565b831461152b5760405162461bcd60e51b8152600401610411906136bd565b43821061154a5760405162461bcd60e51b8152600401610411906136e7565b6001546040516396b477cb60e01b8152600481018490525f916001600160a01b0316906396b477cb90602401602060405180830381865afa158015611591573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906115b5919061350c565b60e088018490526101008801819052436101208901526101408801839052905061144e8761215e565b5f805f805f6115f260fd8760c00151611e8c565b94509450945094509450846116495760405162461bcd60e51b815260206004820152601d60248201527f6661696c656420696e6e6f63656e636520766572696669636174696f6e0000006044820152606401610411565b8560a001516001600160a01b0316846001600160a01b03161461167e5760405162461bcd60e51b815260040161041190613692565b8560600151600d81111561169457611694612ff3565b83146116b25760405162461bcd60e51b8152600401610411906136bd565b4382106116d15760405162461bcd60e51b8152600401610411906136e7565b60e0860182905261014086018190526116e986612400565b505050505050565b5f600a82600d81111561170657611706612ff3565b036117145760025b92915050565b5f82600d81111561172757611727612ff3565b0361173357600261170e565b600182600d81111561174757611747612ff3565b0361175357600261170e565b600261170e565b6011545b601054811015611a8d575f6010828154811061177c5761177c61354a565b905f5260205f2001549050805f036117945750611a85565b61179f600182613537565b90505f600282815481106117b5576117b561354a565b5f91825260209182902060408051610160810182526007909302909101805460ff80821685526101008204811695850195909552929390929184019162010000900416600281111561180957611809612ff3565b600281111561181a5761181a612ff3565b815281546020909101906301000000900460ff16600d81111561183f5761183f612ff3565b600d81111561185057611850612ff3565b815281546001600160a01b03600160201b90910481166020830152600183015416604082015260028201805460609092019161188b90613475565b80601f01602080910402602001604051908101604052809291908181526020018280546118b790613475565b80156119025780601f106118d957610100808354040283529160200191611902565b820191905f5260205f20905b8154815290600101906020018083116118e557829003601f168201915b5050505050815260200160038201548152602001600482015481526020016005820154815260200160068201548152505090504360035f015482610120015161194b919061355e565b1115611958575050601155565b60a08101516001600160a01b03165f908152600c602052604081208190556060820151611984906116f1565b60a08301516001600160a01b03165f908152600e60209081526040808320610100870151845290915290205490915081116119c157505050611a85565b60a0820180516001600160a01b039081165f908152600e6020908152604080832061010088015184528252808320869055845184168352600b825280832080546001808201835591855283852001899055600f805491820181559093527f8d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac80290920187905592518151858152938401879052909116917f6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f910160405180910390a25050505b60010161175e565b601155565b5f8060015f9054906101000a90046001600160a01b03166001600160a01b031663c9d97af46040518163ffffffff1660e01b8152600401602060405180830381865afa158015611ae4573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611b08919061350c565b90505f5b600f54811015611b7057816002600f8381548110611b2c57611b2c61354a565b905f5260205f20015481548110611b4557611b4561354a565b905f5260205f2090600702016004015403611b6857611b6560018461355e565b92505b600101611b0c565b505f5b600f54811015611d3a57611d326002600f8381548110611b9557611b9561354a565b905f5260205f20015481548110611bae57611bae61354a565b5f91825260209182902060408051610160810182526007909302909101805460ff808216855261010082048116958501959095529293909291840191620100009004166002811115611c0257611c02612ff3565b6002811115611c1357611c13612ff3565b815281546020909101906301000000900460ff16600d811115611c3857611c38612ff3565b600d811115611c4957611c49612ff3565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191611c8490613475565b80601f0160208091040260200160405190810160405280929190818152602001828054611cb090613475565b8015611cfb5780601f10611cd257610100808354040283529160200191611cfb565b820191905f5260205f20905b815481529060010190602001808311611cde57829003601f168201915b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481525050846126fc565b600101611b73565b50610641600f5f612d11565b815460026001808316156101000203821604825180820160208110602084100160028114611dee5760018114611e1357865f526020840460205f2001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b80821015611dd85781518455600184019350602082019150611dbb565b815191036101000a90819004029091555061144e565b60028302826020036101000a846020036101000a60208901510402018501875561144e565b865f526020840460205f2001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b80821015611e725781518455600184019350602082019150611e55565b815191036101000a90819004029091555050505050505050565b5f805f805f8086516020611ea0919061355e565b9050611eaa612d2c565b60a081838a8c5afa611eba575f80fd5b8051600103611ec857600196505b602081015160408201516060830151608090930151989b919a509850909695509350505050565b5f611efd82606001516116f1565b60a08301516001600160a01b03165f908152600e6020908152604080832061010087015184529091529020549091508111611f4a5760405162461bcd60e51b815260040161041190613717565b600280546001810182555f829052835160079091027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054602086015160ff9081166101000261ffff1990921693169290921791909117808255604085015185939091839162ff0000199091169062010000908490811115611fcf57611fcf612ff3565b021790555060608201518154829063ff0000001916630100000083600d811115611ffb57611ffb612ff3565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c0820151600282019061205990826135bd565b5060e0820151600382015561010082015160048201556101208201516005820155610140909101516006909101556002545f9061209890600190613537565b60a0840180516001600160a01b039081165f908152600b6020908152604080832080546001808201835591855283852001879055600f8054918201815584527f8d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac80201869055845184168352600e82528083206101008a01518452825291829020879055925181518781529384018590529394509216917f6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f91015b60405180910390a2505050565b60a08101516001600160a01b03165f908152600c6020526040902054156121c75760405162461bcd60e51b815260206004820181905260248201527f616c72656164792070726f63657373696e6720616e2061636375736174696f6e6044820152606401610411565b5f6121d582606001516116f1565b60a08301516001600160a01b03165f908152600e60209081526040808320610100870151845290915290205490915081116122225760405162461bcd60e51b815260040161041190613717565b600280546001810182555f829052835160079091027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054602086015160ff9081166101000261ffff1990921693169290921791909117808255604085015185939091839162ff00001990911690620100009084908111156122a7576122a7612ff3565b021790555060608201518154829063ff0000001916630100000083600d8111156122d3576122d3612ff3565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c0820151600282019061233190826135bd565b5060e0820151600382015561010082015160048201556101208201516005820155610140909101516006909101556002545f9061237090600190613537565b905061237d81600161355e565b60a08401516001600160a01b03165f908152600c602052604090205560106123a682600161355e565b81546001810183555f928352602092839020015560a0840151604080518581529283018490526001600160a01b03909116917f2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba403519101612151565b60a08101516001600160a01b03165f908152600c60205260408120549081900361246c5760405162461bcd60e51b815260206004820152601860248201527f6e6f206173736f6369617465642061636375736174696f6e00000000000000006044820152606401610411565b8160600151600d81111561248257612482612ff3565b600261248f600184613537565b8154811061249f5761249f61354a565b5f9182526020909120600790910201546301000000900460ff16600d8111156124ca576124ca612ff3565b146125275760405162461bcd60e51b815260206004820152602760248201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e206044820152661c9d5b19481a5960ca1b6064820152608401610411565b60e08201516002612539600184613537565b815481106125495761254961354a565b905f5260205f20906007020160030154146125b45760405162461bcd60e51b815260206004820152602560248201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e20604482015264626c6f636b60d81b6064820152608401610411565b61014082015160026125c7600184613537565b815481106125d7576125d761354a565b905f5260205f20906007020160060154146126405760405162461bcd60e51b8152602060048201526024808201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e206044820152630d0c2e6d60e31b6064820152608401610411565b6011545b60105481101561269c5781601082815481106126625761266261354a565b905f5260205f20015403612694575f601082815481106126845761268461354a565b5f9182526020909120015561269c565b600101612644565b5060a0820180516001600160a01b039081165f908152600c602090815260408083208390559351935191825292909116917f1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f910160405180910390a25050565b60015460a0830151604051630c825d9760e11b81526001600160a01b0391821660048201525f929190911690631904bb2e906024015f60405180830381865afa15801561274b573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526127729190810190613322565b608084015160a08501516001600160a01b039081165f908152600a6020526040902080546001600160a01b031916919092161790559050600381610240015160038111156127c2576127c2612ff3565b036127cc57505050565b5f6127e26127dd85606001516116f1565b612c55565b610220830151600754919250905f906127fb908361375b565b600654612808908761375b565b612812908561355e565b61281c919061355e565b60095490915081111561282e57506009545b5f8461012001518560c001518660a00151612849919061355e565b612853919061355e565b6009549091505f90612865838561375b565b61286f9190613772565b90505f8111801561287f57508181145b15612996575f60a087018190526101008701819052610120870181905260c08701526101e0860180518291906128b690839061355e565b90525061022086018051600191906128cf90839061355e565b905250600361024087015260015460405163f5371f4b60e01b81526001600160a01b039091169063f5371f4b9061290a9089906004016137a1565b5f604051808303815f87803b158015612921575f80fd5b505af1158015612933573d5f803e3d5ffd5b50505050602086810151604080516001600160a01b0390921682529181018390525f818301526001606082015290517f876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee89181900360800190a15050505050505050565b610120860151819081116129c2578087610120018181516129b79190613537565b9052505f90506129dc565b6101208701516129d29082613537565b5f61012089015290505b8015612a57578087610100015110612a2357808761010001818151612a019190613537565b90525060a087018051829190612a18908390613537565b9052505f9050612a57565b610100870151612a339082613537565b90508661010001518760a001818151612a4c9190613537565b9052505f6101008801525b5f81118015612a7857505f8760a001518860c00151612a76919061355e565b115b15612b22575f8760a001518860c00151612a92919061355e565b60c0890151612aa1908461375b565b612aab9190613772565b90505f8860a001518960c00151612ac2919061355e565b60a08a0151612ad1908561375b565b612adb9190613772565b9050818960c001818151612aef9190613537565b90525060a089018051829190612b06908390613537565b905250612b13818361355e565b612b1d9084613537565b925050505b612b2c8183613537565b915081876101e001818151612b41919061355e565b9052506102208701805160019190612b5a90839061355e565b9052505f54610220880151600854612b72919061375b565b612b7c919061375b565b612b86904361355e565b610200880152600261024088015260015460405163f5371f4b60e01b81526001600160a01b039091169063f5371f4b90612bc4908a906004016137a1565b5f604051808303815f87803b158015612bdb575f80fd5b505af1158015612bed573d5f803e3d5ffd5b50505050602087810151610200890151604080516001600160a01b039093168352928201859052818301525f606082015290517f876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee89181900360800190a1505050505050505050565b5f81612c6357505060055490565b60018203612c7357505060055490565b60028203612c8357505060055490565b60038203612c9357505060055490565b60048203612ca45750612710919050565b50612710919050565b60408051610160810182525f80825260208201819052909182019081526020015f81526020015f6001600160a01b031681526020015f6001600160a01b03168152602001606081526020015f81526020015f81526020015f81526020015f81525090565b5080545f8255905f5260205f2090810190610aaa9190612d4a565b6040518060a001604052806005906020820280368337509192915050565b5b80821115612d5e575f8155600101612d4b565b5090565b6001600160a01b0381168114610aaa575f80fd5b8035612d8181612d62565b919050565b5f60208284031215612d96575f80fd5b8135612da181612d62565b9392505050565b634e487b7160e01b5f52604160045260245ffd5b604051610160810167ffffffffffffffff81118282101715612de057612de0612da8565b60405290565b604051610260810167ffffffffffffffff81118282101715612de057612de0612da8565b604051601f8201601f1916810167ffffffffffffffff81118282101715612e3357612e33612da8565b604052919050565b803560ff81168114612d81575f80fd5b803560038110612d81575f80fd5b8035600e8110612d81575f80fd5b5f67ffffffffffffffff821115612e8057612e80612da8565b50601f01601f191660200190565b5f82601f830112612e9d575f80fd5b8135612eb0612eab82612e67565b612e0a565b818152846020838601011115612ec4575f80fd5b816020850160208301375f918101602001919091529392505050565b5f60208284031215612ef0575f80fd5b813567ffffffffffffffff80821115612f07575f80fd5b908301906101608286031215612f1b575f80fd5b612f23612dbc565b612f2c83612e3b565b8152612f3a60208401612e3b565b6020820152612f4b60408401612e4b565b6040820152612f5c60608401612e59565b6060820152612f6d60808401612d76565b6080820152612f7e60a08401612d76565b60a082015260c083013582811115612f94575f80fd5b612fa087828601612e8e565b60c08301525060e08381013590820152610100808401359082015261012080840135908201526101409283013592810192909252509392505050565b5f60208284031215612fec575f80fd5b5035919050565b634e487b7160e01b5f52602160045260245ffd5b6003811061301757613017612ff3565b9052565b600e811061301757613017612ff3565b5f5b8381101561304557818101518382015260200161302d565b50505f910152565b5f815180845261306481602086016020860161302b565b601f01601f19169290920160200192915050565b5f61016060ff8e16835260ff8d166020840152613098604084018d613007565b6130a5606084018c61301b565b6001600160a01b038a81166080850152891660a084015260c083018190526130cf8184018961304d565b9150508560e08301528461010083015283610120830152826101408301529c9b505050505050505050505050565b5f805f6060848603121561310f575f80fd5b833561311a81612d62565b925061312860208501612e59565b9150604084013590509250925092565b5f60208284031215613148575f80fd5b81358015158114612da1575f80fd5b805160ff1682525f6101606020830151613176602086018260ff169052565b5060408301516131896040860182613007565b50606083015161319c606086018261301b565b5060808301516131b760808601826001600160a01b03169052565b5060a08301516131d260a08601826001600160a01b03169052565b5060c08301518160c08601526131ea8286018261304d565b60e0858101519087015261010080860151908701526101208086015190870152610140948501519490950193909352509192915050565b602081525f612da16020830184613157565b5f60208083016020845280855180835260408601915060408160051b8701019250602087015f5b8281101561328857603f19888603018452613276858351613157565b9450928501929085019060010161325a565b5092979650505050505050565b5f80604083850312156132a6575f80fd5b82356132b181612d62565b946020939093013593505050565b8051612d8181612d62565b5f82601f8301126132d9575f80fd5b81516132e7612eab82612e67565b8181528460208386010111156132fb575f80fd5b61330c82602083016020870161302b565b949350505050565b805160048110612d81575f80fd5b5f60208284031215613332575f80fd5b815167ffffffffffffffff80821115613349575f80fd5b90830190610260828603121561335d575f80fd5b613365612de6565b61336e836132bf565b815261337c602084016132bf565b602082015261338d604084016132bf565b60408201526060830151828111156133a3575f80fd5b6133af878286016132ca565b6060830152506080838101519082015260a0808401519082015260c0808401519082015260e08084015190820152610100808401519082015261012080840151908201526101408084015190820152610160808401519082015261018091506134198284016132bf565b918101919091526101a082810151908201526101c080830151908201526101e080830151908201526102008083015190820152610220808301519082015261024090613466828401613314565b91810191909152949350505050565b600181811c9082168061348957607f821691505b6020821081036134a757634e487b7160e01b5f52602260045260245ffd5b50919050565b60208082526024908201527f66756e6374696f6e207265737472696374656420746f207468652076616c696460408201526330ba37b960e11b606082015260800190565b5f60208284031215613501575f80fd5b8151612da181612d62565b5f6020828403121561351c575f80fd5b5051919050565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561170e5761170e613523565b634e487b7160e01b5f52603260045260245ffd5b8082018082111561170e5761170e613523565b601f8211156135b857805f5260205f20601f840160051c810160208510156135965750805b601f840160051c820191505b818110156135b5575f81556001016135a2565b50505b505050565b815167ffffffffffffffff8111156135d7576135d7612da8565b6135eb816135e58454613475565b84613571565b602080601f83116001811461361e575f84156136075750858301515b5f19600386901b1c1916600185901b1785556116e9565b5f85815260208120601f198616915b8281101561364c5788860151825594840194600190910190840161362d565b508582101561366957878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b60ff818116838216019081111561170e5761170e613523565b6020808252601190820152700decccccadcc8cae440dad2e6dac2e8c6d607b1b604082015260600190565b60208082526010908201526f0e4ead8ca40d2c840dad2e6dac2e8c6d60831b604082015260600190565b60208082526016908201527563616e277420626520696e207468652066757475726560501b604082015260600190565b60208082526024908201527f616c726561647920736c6173686564206174207468652070726f6f66277320656040820152630e0dec6d60e31b606082015260800190565b808202811582820484141761170e5761170e613523565b5f8261378c57634e487b7160e01b5f52601260045260245ffd5b500490565b6004811061301757613017612ff3565b602081526137bb6020820183516001600160a01b03169052565b5f60208301516137d660408401826001600160a01b03169052565b5060408301516001600160a01b038116606084015250606083015161026080608085015261380861028085018361304d565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a061388f818701836001600160a01b03169052565b8601516101c0868101919091528601516101e08087019190915286015161020080870191909152860151610220808701919091528601516102408087019190915286015190506138e182860182613791565b509094935050505056fea26469706673582212202a3dfe47d85c3cbd2a124f8563e19f8983820b40c697faabe8851678fde17de464736f6c63430008170033",
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
func DeployAccountability(auth *bind.TransactOpts, backend bind.ContractBackend, _autonity common.Address, _config AccountabilityConfig) (common.Address, *types.Transaction, *Accountability, error) {
	parsed, err := AccountabilityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AccountabilityBin), backend, _autonity, _config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Accountability{AccountabilityCaller: AccountabilityCaller{contract: contract}, AccountabilityTransactor: AccountabilityTransactor{contract: contract}, AccountabilityFilterer: AccountabilityFilterer{contract: contract}}, nil
}

// Accountability is an auto generated Go binding around an Ethereum contract.
type Accountability struct {
	AccountabilityCaller     // Read-only binding to the contract
	AccountabilityTransactor // Write-only binding to the contract
	AccountabilityFilterer   // Log filterer for contract events
}

// AccountabilityCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccountabilityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountabilityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccountabilityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountabilityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccountabilityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountabilitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccountabilitySession struct {
	Contract     *Accountability   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccountabilityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccountabilityCallerSession struct {
	Contract *AccountabilityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// AccountabilityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccountabilityTransactorSession struct {
	Contract     *AccountabilityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AccountabilityRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccountabilityRaw struct {
	Contract *Accountability // Generic contract binding to access the raw methods on
}

// AccountabilityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccountabilityCallerRaw struct {
	Contract *AccountabilityCaller // Generic read-only contract binding to access the raw methods on
}

// AccountabilityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccountabilityTransactorRaw struct {
	Contract *AccountabilityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccountability creates a new instance of Accountability, bound to a specific deployed contract.
func NewAccountability(address common.Address, backend bind.ContractBackend) (*Accountability, error) {
	contract, err := bindAccountability(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Accountability{AccountabilityCaller: AccountabilityCaller{contract: contract}, AccountabilityTransactor: AccountabilityTransactor{contract: contract}, AccountabilityFilterer: AccountabilityFilterer{contract: contract}}, nil
}

// NewAccountabilityCaller creates a new read-only instance of Accountability, bound to a specific deployed contract.
func NewAccountabilityCaller(address common.Address, caller bind.ContractCaller) (*AccountabilityCaller, error) {
	contract, err := bindAccountability(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccountabilityCaller{contract: contract}, nil
}

// NewAccountabilityTransactor creates a new write-only instance of Accountability, bound to a specific deployed contract.
func NewAccountabilityTransactor(address common.Address, transactor bind.ContractTransactor) (*AccountabilityTransactor, error) {
	contract, err := bindAccountability(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccountabilityTransactor{contract: contract}, nil
}

// NewAccountabilityFilterer creates a new log filterer instance of Accountability, bound to a specific deployed contract.
func NewAccountabilityFilterer(address common.Address, filterer bind.ContractFilterer) (*AccountabilityFilterer, error) {
	contract, err := bindAccountability(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccountabilityFilterer{contract: contract}, nil
}

// bindAccountability binds a generic wrapper to an already deployed contract.
func bindAccountability(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountabilityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Accountability *AccountabilityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Accountability.Contract.AccountabilityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Accountability *AccountabilityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Accountability.Contract.AccountabilityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Accountability *AccountabilityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Accountability.Contract.AccountabilityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Accountability *AccountabilityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Accountability.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Accountability *AccountabilityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Accountability.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Accountability *AccountabilityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Accountability.Contract.contract.Transact(opts, method, params...)
}

// Beneficiaries is a free data retrieval call binding the contract method 0x01567739.
//
// Solidity: function beneficiaries(address ) view returns(address)
func (_Accountability *AccountabilityCaller) Beneficiaries(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "beneficiaries", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Beneficiaries is a free data retrieval call binding the contract method 0x01567739.
//
// Solidity: function beneficiaries(address ) view returns(address)
func (_Accountability *AccountabilitySession) Beneficiaries(arg0 common.Address) (common.Address, error) {
	return _Accountability.Contract.Beneficiaries(&_Accountability.CallOpts, arg0)
}

// Beneficiaries is a free data retrieval call binding the contract method 0x01567739.
//
// Solidity: function beneficiaries(address ) view returns(address)
func (_Accountability *AccountabilityCallerSession) Beneficiaries(arg0 common.Address) (common.Address, error) {
	return _Accountability.Contract.Beneficiaries(&_Accountability.CallOpts, arg0)
}

// CanAccuse is a free data retrieval call binding the contract method 0x7ccecadd.
//
// Solidity: function canAccuse(address _offender, uint8 _rule, uint256 _block) view returns(bool _result, uint256 _deadline)
func (_Accountability *AccountabilityCaller) CanAccuse(opts *bind.CallOpts, _offender common.Address, _rule uint8, _block *big.Int) (struct {
	Result   bool
	Deadline *big.Int
}, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "canAccuse", _offender, _rule, _block)

	outstruct := new(struct {
		Result   bool
		Deadline *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Result = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Deadline = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// CanAccuse is a free data retrieval call binding the contract method 0x7ccecadd.
//
// Solidity: function canAccuse(address _offender, uint8 _rule, uint256 _block) view returns(bool _result, uint256 _deadline)
func (_Accountability *AccountabilitySession) CanAccuse(_offender common.Address, _rule uint8, _block *big.Int) (struct {
	Result   bool
	Deadline *big.Int
}, error) {
	return _Accountability.Contract.CanAccuse(&_Accountability.CallOpts, _offender, _rule, _block)
}

// CanAccuse is a free data retrieval call binding the contract method 0x7ccecadd.
//
// Solidity: function canAccuse(address _offender, uint8 _rule, uint256 _block) view returns(bool _result, uint256 _deadline)
func (_Accountability *AccountabilityCallerSession) CanAccuse(_offender common.Address, _rule uint8, _block *big.Int) (struct {
	Result   bool
	Deadline *big.Int
}, error) {
	return _Accountability.Contract.CanAccuse(&_Accountability.CallOpts, _offender, _rule, _block)
}

// CanSlash is a free data retrieval call binding the contract method 0x4108a95a.
//
// Solidity: function canSlash(address _offender, uint8 _rule, uint256 _block) view returns(bool)
func (_Accountability *AccountabilityCaller) CanSlash(opts *bind.CallOpts, _offender common.Address, _rule uint8, _block *big.Int) (bool, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "canSlash", _offender, _rule, _block)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanSlash is a free data retrieval call binding the contract method 0x4108a95a.
//
// Solidity: function canSlash(address _offender, uint8 _rule, uint256 _block) view returns(bool)
func (_Accountability *AccountabilitySession) CanSlash(_offender common.Address, _rule uint8, _block *big.Int) (bool, error) {
	return _Accountability.Contract.CanSlash(&_Accountability.CallOpts, _offender, _rule, _block)
}

// CanSlash is a free data retrieval call binding the contract method 0x4108a95a.
//
// Solidity: function canSlash(address _offender, uint8 _rule, uint256 _block) view returns(bool)
func (_Accountability *AccountabilityCallerSession) CanSlash(_offender common.Address, _rule uint8, _block *big.Int) (bool, error) {
	return _Accountability.Contract.CanSlash(&_Accountability.CallOpts, _offender, _rule, _block)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(uint256 innocenceProofSubmissionWindow, uint256 baseSlashingRateLow, uint256 baseSlashingRateMid, uint256 collusionFactor, uint256 historyFactor, uint256 jailFactor, uint256 slashingRatePrecision)
func (_Accountability *AccountabilityCaller) Config(opts *bind.CallOpts) (struct {
	InnocenceProofSubmissionWindow *big.Int
	BaseSlashingRateLow            *big.Int
	BaseSlashingRateMid            *big.Int
	CollusionFactor                *big.Int
	HistoryFactor                  *big.Int
	JailFactor                     *big.Int
	SlashingRatePrecision          *big.Int
}, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "config")

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
		return *outstruct, err
	}

	outstruct.InnocenceProofSubmissionWindow = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BaseSlashingRateLow = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BaseSlashingRateMid = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.CollusionFactor = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.HistoryFactor = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.JailFactor = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.SlashingRatePrecision = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(uint256 innocenceProofSubmissionWindow, uint256 baseSlashingRateLow, uint256 baseSlashingRateMid, uint256 collusionFactor, uint256 historyFactor, uint256 jailFactor, uint256 slashingRatePrecision)
func (_Accountability *AccountabilitySession) Config() (struct {
	InnocenceProofSubmissionWindow *big.Int
	BaseSlashingRateLow            *big.Int
	BaseSlashingRateMid            *big.Int
	CollusionFactor                *big.Int
	HistoryFactor                  *big.Int
	JailFactor                     *big.Int
	SlashingRatePrecision          *big.Int
}, error) {
	return _Accountability.Contract.Config(&_Accountability.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(uint256 innocenceProofSubmissionWindow, uint256 baseSlashingRateLow, uint256 baseSlashingRateMid, uint256 collusionFactor, uint256 historyFactor, uint256 jailFactor, uint256 slashingRatePrecision)
func (_Accountability *AccountabilityCallerSession) Config() (struct {
	InnocenceProofSubmissionWindow *big.Int
	BaseSlashingRateLow            *big.Int
	BaseSlashingRateMid            *big.Int
	CollusionFactor                *big.Int
	HistoryFactor                  *big.Int
	JailFactor                     *big.Int
	SlashingRatePrecision          *big.Int
}, error) {
	return _Accountability.Contract.Config(&_Accountability.CallOpts)
}

// EpochPeriod is a free data retrieval call binding the contract method 0xb5b7a184.
//
// Solidity: function epochPeriod() view returns(uint256)
func (_Accountability *AccountabilityCaller) EpochPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "epochPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochPeriod is a free data retrieval call binding the contract method 0xb5b7a184.
//
// Solidity: function epochPeriod() view returns(uint256)
func (_Accountability *AccountabilitySession) EpochPeriod() (*big.Int, error) {
	return _Accountability.Contract.EpochPeriod(&_Accountability.CallOpts)
}

// EpochPeriod is a free data retrieval call binding the contract method 0xb5b7a184.
//
// Solidity: function epochPeriod() view returns(uint256)
func (_Accountability *AccountabilityCallerSession) EpochPeriod() (*big.Int, error) {
	return _Accountability.Contract.EpochPeriod(&_Accountability.CallOpts)
}

// Events is a free data retrieval call binding the contract method 0x0b791430.
//
// Solidity: function events(uint256 ) view returns(uint8 chunks, uint8 chunkId, uint8 eventType, uint8 rule, address reporter, address offender, bytes rawProof, uint256 block, uint256 epoch, uint256 reportingBlock, uint256 messageHash)
func (_Accountability *AccountabilityCaller) Events(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Chunks         uint8
	ChunkId        uint8
	EventType      uint8
	Rule           uint8
	Reporter       common.Address
	Offender       common.Address
	RawProof       []byte
	Block          *big.Int
	Epoch          *big.Int
	ReportingBlock *big.Int
	MessageHash    *big.Int
}, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "events", arg0)

	outstruct := new(struct {
		Chunks         uint8
		ChunkId        uint8
		EventType      uint8
		Rule           uint8
		Reporter       common.Address
		Offender       common.Address
		RawProof       []byte
		Block          *big.Int
		Epoch          *big.Int
		ReportingBlock *big.Int
		MessageHash    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Chunks = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.ChunkId = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.EventType = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.Rule = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.Reporter = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Offender = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.RawProof = *abi.ConvertType(out[6], new([]byte)).(*[]byte)
	outstruct.Block = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Epoch = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.ReportingBlock = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)
	outstruct.MessageHash = *abi.ConvertType(out[10], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Events is a free data retrieval call binding the contract method 0x0b791430.
//
// Solidity: function events(uint256 ) view returns(uint8 chunks, uint8 chunkId, uint8 eventType, uint8 rule, address reporter, address offender, bytes rawProof, uint256 block, uint256 epoch, uint256 reportingBlock, uint256 messageHash)
func (_Accountability *AccountabilitySession) Events(arg0 *big.Int) (struct {
	Chunks         uint8
	ChunkId        uint8
	EventType      uint8
	Rule           uint8
	Reporter       common.Address
	Offender       common.Address
	RawProof       []byte
	Block          *big.Int
	Epoch          *big.Int
	ReportingBlock *big.Int
	MessageHash    *big.Int
}, error) {
	return _Accountability.Contract.Events(&_Accountability.CallOpts, arg0)
}

// Events is a free data retrieval call binding the contract method 0x0b791430.
//
// Solidity: function events(uint256 ) view returns(uint8 chunks, uint8 chunkId, uint8 eventType, uint8 rule, address reporter, address offender, bytes rawProof, uint256 block, uint256 epoch, uint256 reportingBlock, uint256 messageHash)
func (_Accountability *AccountabilityCallerSession) Events(arg0 *big.Int) (struct {
	Chunks         uint8
	ChunkId        uint8
	EventType      uint8
	Rule           uint8
	Reporter       common.Address
	Offender       common.Address
	RawProof       []byte
	Block          *big.Int
	Epoch          *big.Int
	ReportingBlock *big.Int
	MessageHash    *big.Int
}, error) {
	return _Accountability.Contract.Events(&_Accountability.CallOpts, arg0)
}

// GetValidatorAccusation is a free data retrieval call binding the contract method 0x9cb22b06.
//
// Solidity: function getValidatorAccusation(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256))
func (_Accountability *AccountabilityCaller) GetValidatorAccusation(opts *bind.CallOpts, _val common.Address) (AccountabilityEvent, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "getValidatorAccusation", _val)

	if err != nil {
		return *new(AccountabilityEvent), err
	}

	out0 := *abi.ConvertType(out[0], new(AccountabilityEvent)).(*AccountabilityEvent)

	return out0, err

}

// GetValidatorAccusation is a free data retrieval call binding the contract method 0x9cb22b06.
//
// Solidity: function getValidatorAccusation(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256))
func (_Accountability *AccountabilitySession) GetValidatorAccusation(_val common.Address) (AccountabilityEvent, error) {
	return _Accountability.Contract.GetValidatorAccusation(&_Accountability.CallOpts, _val)
}

// GetValidatorAccusation is a free data retrieval call binding the contract method 0x9cb22b06.
//
// Solidity: function getValidatorAccusation(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256))
func (_Accountability *AccountabilityCallerSession) GetValidatorAccusation(_val common.Address) (AccountabilityEvent, error) {
	return _Accountability.Contract.GetValidatorAccusation(&_Accountability.CallOpts, _val)
}

// GetValidatorFaults is a free data retrieval call binding the contract method 0xbebaa8fc.
//
// Solidity: function getValidatorFaults(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256)[])
func (_Accountability *AccountabilityCaller) GetValidatorFaults(opts *bind.CallOpts, _val common.Address) ([]AccountabilityEvent, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "getValidatorFaults", _val)

	if err != nil {
		return *new([]AccountabilityEvent), err
	}

	out0 := *abi.ConvertType(out[0], new([]AccountabilityEvent)).(*[]AccountabilityEvent)

	return out0, err

}

// GetValidatorFaults is a free data retrieval call binding the contract method 0xbebaa8fc.
//
// Solidity: function getValidatorFaults(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256)[])
func (_Accountability *AccountabilitySession) GetValidatorFaults(_val common.Address) ([]AccountabilityEvent, error) {
	return _Accountability.Contract.GetValidatorFaults(&_Accountability.CallOpts, _val)
}

// GetValidatorFaults is a free data retrieval call binding the contract method 0xbebaa8fc.
//
// Solidity: function getValidatorFaults(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256)[])
func (_Accountability *AccountabilityCallerSession) GetValidatorFaults(_val common.Address) ([]AccountabilityEvent, error) {
	return _Accountability.Contract.GetValidatorFaults(&_Accountability.CallOpts, _val)
}

// SlashingHistory is a free data retrieval call binding the contract method 0xe7bb0b52.
//
// Solidity: function slashingHistory(address , uint256 ) view returns(uint256)
func (_Accountability *AccountabilityCaller) SlashingHistory(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Accountability.contract.Call(opts, &out, "slashingHistory", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SlashingHistory is a free data retrieval call binding the contract method 0xe7bb0b52.
//
// Solidity: function slashingHistory(address , uint256 ) view returns(uint256)
func (_Accountability *AccountabilitySession) SlashingHistory(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Accountability.Contract.SlashingHistory(&_Accountability.CallOpts, arg0, arg1)
}

// SlashingHistory is a free data retrieval call binding the contract method 0xe7bb0b52.
//
// Solidity: function slashingHistory(address , uint256 ) view returns(uint256)
func (_Accountability *AccountabilityCallerSession) SlashingHistory(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Accountability.Contract.SlashingHistory(&_Accountability.CallOpts, arg0, arg1)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_Accountability *AccountabilityTransactor) DistributeRewards(opts *bind.TransactOpts, _validator common.Address) (*types.Transaction, error) {
	return _Accountability.contract.Transact(opts, "distributeRewards", _validator)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_Accountability *AccountabilitySession) DistributeRewards(_validator common.Address) (*types.Transaction, error) {
	return _Accountability.Contract.DistributeRewards(&_Accountability.TransactOpts, _validator)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_Accountability *AccountabilityTransactorSession) DistributeRewards(_validator common.Address) (*types.Transaction, error) {
	return _Accountability.Contract.DistributeRewards(&_Accountability.TransactOpts, _validator)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_Accountability *AccountabilityTransactor) Finalize(opts *bind.TransactOpts, _epochEnd bool) (*types.Transaction, error) {
	return _Accountability.contract.Transact(opts, "finalize", _epochEnd)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_Accountability *AccountabilitySession) Finalize(_epochEnd bool) (*types.Transaction, error) {
	return _Accountability.Contract.Finalize(&_Accountability.TransactOpts, _epochEnd)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_Accountability *AccountabilityTransactorSession) Finalize(_epochEnd bool) (*types.Transaction, error) {
	return _Accountability.Contract.Finalize(&_Accountability.TransactOpts, _epochEnd)
}

// HandleEvent is a paid mutator transaction binding the contract method 0x069f6863.
//
// Solidity: function handleEvent((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256) _event) returns()
func (_Accountability *AccountabilityTransactor) HandleEvent(opts *bind.TransactOpts, _event AccountabilityEvent) (*types.Transaction, error) {
	return _Accountability.contract.Transact(opts, "handleEvent", _event)
}

// HandleEvent is a paid mutator transaction binding the contract method 0x069f6863.
//
// Solidity: function handleEvent((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256) _event) returns()
func (_Accountability *AccountabilitySession) HandleEvent(_event AccountabilityEvent) (*types.Transaction, error) {
	return _Accountability.Contract.HandleEvent(&_Accountability.TransactOpts, _event)
}

// HandleEvent is a paid mutator transaction binding the contract method 0x069f6863.
//
// Solidity: function handleEvent((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256) _event) returns()
func (_Accountability *AccountabilityTransactorSession) HandleEvent(_event AccountabilityEvent) (*types.Transaction, error) {
	return _Accountability.Contract.HandleEvent(&_Accountability.TransactOpts, _event)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_Accountability *AccountabilityTransactor) SetEpochPeriod(opts *bind.TransactOpts, _newPeriod *big.Int) (*types.Transaction, error) {
	return _Accountability.contract.Transact(opts, "setEpochPeriod", _newPeriod)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_Accountability *AccountabilitySession) SetEpochPeriod(_newPeriod *big.Int) (*types.Transaction, error) {
	return _Accountability.Contract.SetEpochPeriod(&_Accountability.TransactOpts, _newPeriod)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_Accountability *AccountabilityTransactorSession) SetEpochPeriod(_newPeriod *big.Int) (*types.Transaction, error) {
	return _Accountability.Contract.SetEpochPeriod(&_Accountability.TransactOpts, _newPeriod)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Offender common.Address
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Offender common.Address
	Severity *big.Int
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Offender common.Address
	Severity *big.Int
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Validator    common.Address
	Amount       *big.Int
	ReleaseBlock *big.Int
	IsJailbound  bool
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSlashingEvent is a free log retrieval operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
func (_Accountability *AccountabilityFilterer) FilterSlashingEvent(opts *bind.FilterOpts) (*AccountabilitySlashingEventIterator, error) {

	logs, sub, err := _Accountability.contract.FilterLogs(opts, "SlashingEvent")
	if err != nil {
		return nil, err
	}
	return &AccountabilitySlashingEventIterator{contract: _Accountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
}

// WatchSlashingEvent is a free log subscription operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
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

// ParseSlashingEvent is a log parse operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
func (_Accountability *AccountabilityFilterer) ParseSlashingEvent(log types.Log) (*AccountabilitySlashingEvent, error) {
	event := new(AccountabilitySlashingEvent)
	if err := _Accountability.contract.UnpackLog(event, "SlashingEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AutonityMetaData contains all meta data concerning the Autonity contract.
var AutonityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquidState\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator[]\",\"name\":\"_validators\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"ActivatedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BondingRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewUnbondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeInitialization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getEpochFromBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquidState\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"liquidLogicAddress\",\"outputs\":[{\"internalType\":\"contractLiquidLogic\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_multisig\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAccountability\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAccountabilityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIACU\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAcuContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setLiquidNewtonLogicContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setOracleContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIStabilization\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setStabilizationContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISupplyControl\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSupplyControlContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquidState\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"_val\",\"type\":\"tuple\"}],\"name\":\"updateValidatorAndTransferSlashedFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"b46e5520": "activateValidator(address)",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
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
		"c2362dd5": "lastEpochBlock()",
		"66d3bd27": "liquidLogicAddress()",
		"40c10f19": "mint(address,uint256)",
		"06fdde03": "name()",
		"0ae65e7a": "pauseValidator(address)",
		"ad722d4d": "registerValidator(string,address,bytes)",
		"cf9c5719": "resetContractUpgrade()",
		"1250a28d": "setAccountabilityContract(address)",
		"d372c07e": "setAcuContract(address)",
		"8bac7dad": "setCommitteeSize(uint256)",
		"6b5f444c": "setEpochPeriod(uint256)",
		"07b5c33b": "setLiquidNewtonLogicContract(address)",
		"cb696f54": "setMinimumBaseFee(uint256)",
		"520fdbbc": "setOperatorAccount(address)",
		"496ccd9b": "setOracleContract(address)",
		"cfd19fb9": "setStabilizationContract(address)",
		"b3ecbadd": "setSupplyControlContract(address)",
		"d886f8a2": "setTreasuryAccount(address)",
		"77e741c7": "setTreasuryFee(uint256)",
		"114eaf55": "setUnbondingPeriod(uint256)",
		"95d89b41": "symbol()",
		"9bb851c0": "totalRedistributed()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"a5d059ca": "unbond(address,uint256)",
		"f5371f4b": "updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,uint8))",
		"b2ea9adb": "upgradeContract(bytes,string)",
	},
	Bin: "0x60806040525f600b555f600c5534801562000018575f80fd5b506040516200b7353803806200b7358339810160408190526200003b9162000e1f565b601b545f036200006357602980546001600160a01b031916331790556200006382826200006b565b5050620012a7565b80518051600d55602080820151600e55604080830151600f55606080840151601055608093840151601180546001600160a01b03199081166001600160a01b03938416179091558487015180516012805484169185169190911790558086015160138054841691851691909117905580850151601480548416918516919091179055808401516015805484169185169190911790559095015160168054871691831691909117905582860151805160178054909716921691909117909455918301516018558281015160195591810151601a55820151601b55513090620001529062000ae2565b6001600160a01b039091168152602001604051809103905ff0801580156200017c573d5f803e3d5ffd5b50602a80546001600160a01b0319166001600160a01b03929092169190911790555f5b82518110156200041a575f838281518110620001bf57620001bf6200102b565b602002602001015160a0015190505f848381518110620001e357620001e36200102b565b60200260200101516101a00181815250505f8483815181106200020a576200020a6200102b565b602002602001015161018001906001600160a01b031690816001600160a01b0316815250505f8483815181106200024557620002456200102b565b602002602001015160a00181815250505f8483815181106200026b576200026b6200102b565b60209081029190910101516101c00152600f5484518590849081106200029557620002956200102b565b602002602001015160800181815250505f848381518110620002bb57620002bb6200102b565b602002602001015161024001906003811115620002dc57620002dc6200103f565b90816003811115620002f257620002f26200103f565b815250505f8483815181106200030c576200030c6200102b565b60200260200101516101600181815250506200034a8483815181106200033657620003366200102b565b60200260200101516200041f60201b60201c565b8060265f8685815181106200036357620003636200102b565b60200260200101515f01516001600160a01b03166001600160a01b031681526020019081526020015f205f8282546200039d919062001067565b925050819055508060285f828254620003b7919062001067565b9250508190555062000410848381518110620003d757620003d76200102b565b60200260200101516020015182868581518110620003f957620003f96200102b565b60200260200101515f01516200043860201b60201c565b506001016200019f565b505050565b6200042a8162000622565b620004358162000757565b50565b5f8211620004995760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b60648201526084015b60405180910390fd5b6001600160a01b0381165f90815260266020526040902054821115620005025760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000490565b6001600160a01b0381165f90815260266020526040812080548492906200052b90849062001083565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580545f908152600394859052978820875181549088166001600160a01b0319918216178255955160018201805491909816961695909517909555905160028401555191015580549192620005b68362001099565b90915550506001600160a01b038481165f81815260276020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d910160405180910390a35050505050565b5f6200063882606001516200098060201b60201c565b6001600160a01b03909116602084015290508015620006885760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000490565b6020808301516001600160a01b039081165f90815260279092526040909120600101541615620006fb5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000490565b61271082608001511115620007535760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000490565b5050565b6101808101516001600160a01b0316620007eb57601c545f906200077b90620009c6565b9050602a5f9054906101000a90046001600160a01b03168260200151835f0151846080015184604051620007af9062000af0565b620007bf959493929190620010b4565b604051809103905ff080158015620007d9573d5f803e3d5ffd5b506001600160a01b0316610180830152505b60208181018051601c8054600180820183555f9283527f0e4562a10381dec21b205ed72637e6b1b523bdd0e4d4d50af5cd23dd4500a21190910180546001600160a01b03199081166001600160a01b03958616179091558451841683526027909552604091829020865181548716908516178155935190840180548616918416919091179055840151600283018054909416911617909155606082015182919060038201906200089c90826200119b565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f8201556102008201516010820155610220820151601182015561024082015160128201805460ff191660018360038111156200097657620009766200103f565b0217905550505050565b5f806200098c62000afe565b5f60408286516020880160ff5afa620009a3575f80fd5b5080516020909101516c0100000000000000000000000090910494909350915050565b6060815f03620009ed5750506040805180820190915260018152600360fc1b602082015290565b815f5b811562000a1c578062000a038162001099565b915062000a149050600a836200127b565b9150620009f0565b5f816001600160401b0381111562000a385762000a3862000b1c565b6040519080825280601f01601f19166020018201604052801562000a63576020820181803683370190505b5090505b841562000ada5762000a7b60018362001083565b915062000a8a600a8662001291565b62000a9790603062001067565b60f81b81838151811062000aaf5762000aaf6200102b565b60200101906001600160f81b03191690815f1a90535062000ad2600a866200127b565b945062000a67565b949350505050565b610f808062008e6183390190565b6119548062009de183390190565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b5f52604160045260245ffd5b60405160a081016001600160401b038111828210171562000b555762000b5562000b1c565b60405290565b604051608081016001600160401b038111828210171562000b555762000b5562000b1c565b60405161026081016001600160401b038111828210171562000b555762000b5562000b1c565b604051601f8201601f191681016001600160401b038111828210171562000bd15762000bd162000b1c565b604052919050565b6001600160a01b038116811462000435575f80fd5b805162000bfb8162000bd9565b919050565b5f5b8381101562000c1c57818101518382015260200162000c02565b50505f910152565b5f82601f83011262000c34575f80fd5b81516001600160401b0381111562000c505762000c5062000b1c565b62000c65601f8201601f191660200162000ba6565b81815284602083860101111562000c7a575f80fd5b62000ada82602083016020870162000c00565b80516004811062000bfb575f80fd5b5f60a0828403121562000cad575f80fd5b62000cb762000b30565b9050815162000cc68162000bd9565b8152602082015162000cd88162000bd9565b6020820152604082015162000ced8162000bd9565b6040820152606082015162000d028162000bd9565b6060820152608082015162000d178162000bd9565b608082015292915050565b5f6080828403121562000d33575f80fd5b62000d3d62000b5b565b9050815162000d4c8162000bd9565b8082525060208201516020820152604082015160408201526060820151606082015292915050565b5f8183036101e081121562000d87575f80fd5b62000d9162000b5b565b915060a081121562000da1575f80fd5b5062000dac62000b30565b82518152602083015160208201526040830151604082015260608301516060820152608083015162000dde8162000bd9565b6080820152815262000df48360a0840162000c9c565b602082015262000e0983610140840162000d22565b60408201526101c0820151606082015292915050565b5f8061020080848603121562000e33575f80fd5b83516001600160401b038082111562000e4a575f80fd5b818601915086601f83011262000e5e575f80fd5b815160208282111562000e755762000e7562000b1c565b8160051b62000e8682820162000ba6565b928352848101820192828101908b85111562000ea0575f80fd5b83870192505b84831015620010095782518681111562000ebe575f80fd5b8701610260818e03601f1901121562000ed5575f80fd5b62000edf62000b80565b62000eec86830162000bee565b815262000efc6040830162000bee565b8682015262000f0e6060830162000bee565b604082015260808201518881111562000f25575f80fd5b62000f358f888386010162000c24565b60608301525060a0820151608082015260c082015160a082015260e082015160c08201526101008083015160e083015261012080840151828401526101409150818401518184015250610160808401518284015261018091508184015181840152506101a062000fa781850162000bee565b828401526101c091508184015181840152506101e080840151828401528b840151818401525050610220808301518b830152610240808401518284015262000ff3610260850162000c8d565b9083015250835250918301919083019062000ea6565b8099505050506200101d89828a0162000d74565b955050505050509250929050565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b634e487b7160e01b5f52601160045260245ffd5b808201808211156200107d576200107d62001053565b92915050565b818103818111156200107d576200107d62001053565b5f60018201620010ad57620010ad62001053565b5060010190565b5f60018060a01b038088168352808716602084015280861660408401525083606083015260a0608083015282518060a0840152620010fa8160c085016020870162000c00565b601f01601f19169190910160c0019695505050505050565b600181811c908216806200112757607f821691505b6020821081036200114657634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200041a57805f5260205f20601f840160051c81016020851015620011735750805b601f840160051c820191505b8181101562001194575f81556001016200117f565b5050505050565b81516001600160401b03811115620011b757620011b762000b1c565b620011cf81620011c8845462001112565b846200114c565b602080601f83116001811462001205575f8415620011ed5750858301515b5f19600386901b1c1916600185901b1785556200125f565b5f85815260208120601f198616915b82811015620012355788860151825594840194600190910190840162001214565b50858210156200125357878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b634e487b7160e01b5f52601260045260245ffd5b5f826200128c576200128c62001267565b500490565b5f82620012a257620012a262001267565b500690565b617bac80620012b55f395ff3fe608060405260043610620003e6575f3560e01c80638bac7dad1162000202578063b66b3e791162000116578063d5f3948811620000a6578063dfb1a4d21162000074578063dfb1a4d21462000d58578063e7f43c681462000d6e578063f5371f4b1462000d8d578063f7866ee31462000db157005b8063d5f394881462000cb4578063d861b0e81462000cd5578063d886f8a21462000cec578063dd62ed3e1462000d1057005b8063cb696f5411620000e4578063cb696f541462000c31578063cf9c57191462000c55578063cfd19fb91462000c6c578063d372c07e1462000c9057005b8063b66b3e791462000bc5578063b7ab4db51462000bec578063c2362dd51462000c03578063c9d97af41462000c1a57005b8063a8b2216e1162000192578063ae1f5fa01162000160578063ae1f5fa01462000b33578063b2ea9adb1462000b59578063b3ecbadd1462000b7d578063b46e55201462000ba157005b8063a8b2216e1462000a9f578063a9059cbb1462000ac5578063ab8f6ffe1462000ae9578063ad722d4d1462000b0f57005b80639c98e47111620001d05780639c98e4711462000a1c5780639dc29fac1462000a33578063a515366a1462000a57578063a5d059ca1462000a7b57005b80638bac7dad146200098557806395d89b4114620009a957806396b477cb14620009d65780639bb851c01462000a0557005b80634364596911620002fa57806370a08231116200028a578063819b64631162000258578063819b64631462000915578063833b1fce146200092b578063852c4849146200094a578063872cf059146200096e57005b806370a08231146200073d578063731b3a03146200077557806377e741c7146200078b57806379502c5514620007af57005b80635f7d394911620002c85780635f7d394914620006a557806366d3bd2714620006e25780636b5f444c14620007035780636fd2c80b146200072757005b8063436459691462000620578063496ccd9b14620006365780634bb278f3146200065a578063520fdbbc146200068157005b80631604e416116200037657806329070c6d116200034457806329070c6d14620005b25780632f2c3f2e14620005c8578063313ce56714620005df57806340c10f1914620005fc57005b80631604e416146200052e57806318160ddd14620005455780631904bb2e146200055b57806323b872dd146200058e57005b80630d8e6e2c11620003b45780630d8e6e2c14620004b05780631122063314620004d0578063114eaf5514620004e65780631250a28d146200050a57005b806306fdde0314620003f057806307b5c33b1462000433578063095ea7b314620004575780630ae65e7a146200048c57005b36620003ee57005b005b348015620003fc575f80fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b6040516200042a91906200545a565b60405180910390f35b3480156200043f575f80fd5b50620003ee620004513660046200548a565b62000dd0565b34801562000463575f80fd5b506200047b62000475366004620054a8565b62000e28565b60405190151581526020016200042a565b34801562000498575f80fd5b50620003ee620004aa3660046200548a565b62000e40565b348015620004bc575f80fd5b50601b545b6040519081526020016200042a565b348015620004dc575f80fd5b50600e54620004c1565b348015620004f2575f80fd5b50620003ee62000504366004620054d5565b62000ec8565b34801562000516575f80fd5b50620003ee620005283660046200548a565b62000efa565b3480156200053a575f80fd5b50620004c160235481565b34801562000551575f80fd5b50602854620004c1565b34801562000567575f80fd5b506200057f620005793660046200548a565b62000f49565b6040516200042a919062005522565b3480156200059a575f80fd5b506200047b620005ac36600462005673565b62001155565b348015620005be575f80fd5b50600d54620004c1565b348015620005d4575f80fd5b50620004c161271081565b348015620005eb575f80fd5b50604051601281526020016200042a565b34801562000608575f80fd5b50620003ee6200061a366004620054a8565b620011ad565b3480156200062c575f80fd5b50601954620004c1565b34801562000642575f80fd5b50620003ee620006543660046200548a565b62001265565b34801562000666575f80fd5b506200067162001365565b6040516200042a92919062005706565b3480156200068d575f80fd5b50620003ee6200069f3660046200548a565b62001672565b348015620006b1575f80fd5b50620006c9620006c336600462005722565b620017fd565b6040516001600160a01b0390911681526020016200042a565b348015620006ee575f80fd5b50602a54620006c9906001600160a01b031681565b3480156200070f575f80fd5b50620003ee62000721366004620054d5565b620019f1565b34801562000733575f80fd5b50601054620004c1565b34801562000749575f80fd5b50620004c16200075b3660046200548a565b6001600160a01b03165f9081526026602052604090205490565b34801562000781575f80fd5b50601f54620004c1565b34801562000797575f80fd5b50620003ee620007a9366004620054d5565b62001b72565b348015620007bb575f80fd5b506040805160a08082018352600d548252600e54602080840191909152600f54838501526010546060808501919091526011546001600160a01b0390811660808087019190915286519485018752601254821685526013548216858501526014548216858801526015548216858401526016548216858201528651908101875260175490911681526018549281019290925260195494820194909452601a5493810193909352601b546200086d939084565b6040805185518152602080870151818301528683015182840152606080880151818401526080978801516001600160a01b03908116898501528751811660a085015282880151811660c085015287850151811660e0850152818801518116610100850152979096015187166101208301528451909616610140820152948301516101608601528201516101808501529101516101a08301526101c08201526101e0016200042a565b34801562000921575f80fd5b50601a54620004c1565b34801562000937575f80fd5b506013546001600160a01b0316620006c9565b34801562000956575f80fd5b50620003ee62000968366004620054a8565b62001ba4565b3480156200097a575f80fd5b50620003ee62001d32565b34801562000991575f80fd5b50620003ee620009a3366004620054d5565b62001d6e565b348015620009b5575f80fd5b50604080518082019091526003815262272a2760e91b60208201526200041b565b348015620009e2575f80fd5b50620004c1620009f4366004620054d5565b5f908152601e602052604090205490565b34801562000a11575f80fd5b50620004c160225481565b34801562000a28575f80fd5b50620004c160205481565b34801562000a3f575f80fd5b50620003ee62000a51366004620054a8565b62001df1565b34801562000a63575f80fd5b50620003ee62000a75366004620054a8565b62001f04565b34801562000a87575f80fd5b50620003ee62000a99366004620054a8565b62001fd5565b34801562000aab575f80fd5b5062000ab66200206a565b6040516200042a919062005743565b34801562000ad1575f80fd5b506200047b62000ae3366004620054a8565b62002148565b34801562000af5575f80fd5b5062000b0062002156565b6040516200042a9190620057a9565b34801562000b1b575f80fd5b50620003ee62000b2d36600462005862565b620021c2565b34801562000b3f575f80fd5b5062000b4a620022da565b6040516200042a9190620058de565b34801562000b65575f80fd5b50620003ee62000b773660046200592c565b62002a8d565b34801562000b89575f80fd5b50620003ee62000b9b3660046200548a565b62002ad3565b34801562000bad575f80fd5b50620003ee62000bbf3660046200548a565b62002b22565b34801562000bd1575f80fd5b5062000bdc62002da2565b6040516200042a92919062005993565b34801562000bf8575f80fd5b5062000b4a62002ed4565b34801562000c0f575f80fd5b50620004c1601f5481565b34801562000c26575f80fd5b50620004c1601d5481565b34801562000c3d575f80fd5b50620003ee62000c4f366004620054d5565b62002f36565b34801562000c61575f80fd5b50620003ee62002f99565b34801562000c78575f80fd5b50620003ee62000c8a3660046200548a565b62002feb565b34801562000c9c575f80fd5b50620003ee62000cae3660046200548a565b6200303a565b34801562000cc0575f80fd5b50602954620006c9906001600160a01b031681565b34801562000ce1575f80fd5b50620003ee62003089565b34801562000cf8575f80fd5b50620003ee62000d0a3660046200548a565b620030ca565b34801562000d1c575f80fd5b50620004c162000d2e366004620059c4565b6001600160a01b039182165f90815260256020908152604080832093909416825291909152205490565b34801562000d64575f80fd5b50601854620004c1565b34801562000d7a575f80fd5b506017546001600160a01b0316620006c9565b34801562000d99575f80fd5b50620003ee62000dab36600462005a00565b62003119565b34801562000dbd575f80fd5b506011546001600160a01b0316620006c9565b6017546001600160a01b0316331462000e065760405162461bcd60e51b815260040162000dfd9062005a3b565b60405180910390fd5b602a80546001600160a01b0319166001600160a01b0392909216919091179055565b5f62000e36338484620032ec565b5060015b92915050565b6001600160a01b038082165f818152602760205260409020600101549091161462000e7f5760405162461bcd60e51b815260040162000dfd9062005a72565b6001600160a01b038181165f9081526027602052604090205416331462000eba5760405162461bcd60e51b815260040162000dfd9062005aa9565b62000ec58162003414565b50565b6017546001600160a01b0316331462000ef55760405162461bcd60e51b815260040162000dfd9062005a3b565b601055565b6017546001600160a01b0316331462000f275760405162461bcd60e51b815260040162000dfd9062005a3b565b601280546001600160a01b0319166001600160a01b0392909216919091179055565b62000f536200524c565b6001600160a01b038083165f818152602760205260409020600101549091161462000f925760405162461bcd60e51b815260040162000dfd9062005af5565b6001600160a01b038083165f90815260276020908152604091829020825161026081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000ff49062005b2c565b80601f0160208091040260200160405190810160405280929190818152602001828054620010229062005b2c565b8015620010715780601f10620010475761010080835404028352916020019162001071565b820191905f5260205f20905b8154815290600101906020018083116200105357829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546080820152600882015460a0820152600982015460c0820152600a82015460e0820152600b820154610100820152600c8201546001600160a01b0316610120820152600d820154610140820152600e820154610160820152600f82015461018082015260108201546101a082015260118201546101c082015260128201546101e09091019060ff166003811115620011385762001138620054ed565b60038111156200114c576200114c620054ed565b90525092915050565b5f62001163848484620034ea565b6001600160a01b0384165f9081526025602090815260408083203384529091528120546200119390849062005b7a565b9050620011a2853383620032ec565b506001949350505050565b6017546001600160a01b03163314620011da5760405162461bcd60e51b815260040162000dfd9062005a3b565b6001600160a01b0382165f90815260266020526040812080548392906200120390849062005b90565b925050819055508060285f8282546200121d919062005b90565b90915550506040518181526001600160a01b038316907f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf906020015b60405180910390a25050565b6017546001600160a01b03163314620012925760405162461bcd60e51b815260040162000dfd9062005a3b565b601380546001600160a01b0319166001600160a01b03838116918217909255601454604051637adbf97360e01b8152600481019290925290911690637adbf973906024015f604051808303815f87803b158015620012ee575f80fd5b505af115801562001301573d5f803e3d5ffd5b5050601654604051637adbf97360e01b81526001600160a01b0385811660048301529091169250637adbf97391506024015b5f604051808303815f87803b1580156200134b575f80fd5b505af11580156200135e573d5f803e3d5ffd5b5050505050565b6029545f906060906001600160a01b03163314620013975760405162461bcd60e51b815260040162000dfd9062005ba6565b601d54435f818152601e6020526040812092909255601854601f54620013be919062005b90565b6012546040516306c9789b60e41b8152929091146004830181905292506001600160a01b031690636c9789b0906024015f604051808303815f87803b15801562001406575f80fd5b505af115801562001419573d5f803e3d5ffd5b50505050801562001507576200142e620035f0565b62001438620039c8565b6200144262003aaa565b5f6200144d620022da565b60135460405163422811f960e11b81529192506001600160a01b03169063845023f29062001480908490600401620058de565b5f604051808303815f87803b15801562001498575f80fd5b505af1158015620014ab573d5f803e3d5ffd5b5050505043601f819055506001601d5f828254620014ca919062005b90565b9091555050601d546040519081527febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e3359060200160405180910390a1505b60135460408051634bb278f360e01b815290515f926001600160a01b031691634bb278f3916004808301926020929190829003018187875af115801562001550573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019062001576919062005be9565b90508015620015f1576014546040805163a2e6204560e01b815290516001600160a01b039092169163a2e620459160048082019260209290919082900301815f875af1925050508015620015e9575060408051601f3d908101601f19168201909252620015e69181019062005be9565b60015b15620015f157505b600254602180546040805160208084028201810190925282815260ff909416939183915f9084015b8282101562001662575f848152602090819020604080518082019091526002850290910180546001600160a01b0316825260019081015482840152908352909201910162001619565b5050505090509350935050509091565b6017546001600160a01b031633146200169f5760405162461bcd60e51b815260040162000dfd9062005a3b565b601780546001600160a01b0319166001600160a01b0383811691821790925560135460405163b3ab15fb60e01b815260048101929092529091169063b3ab15fb906024015f604051808303815f87803b158015620016fb575f80fd5b505af11580156200170e573d5f803e3d5ffd5b505060145460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024015f604051808303815f87803b15801562001757575f80fd5b505af11580156200176a573d5f803e3d5ffd5b505060155460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024015f604051808303815f87803b158015620017b3575f80fd5b505af1158015620017c6573d5f803e3d5ffd5b505060165460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb915060240162001333565b5f80805b6021548110156200184a576021818154811062001822576200182262005c0a565b905f5260205f20906002020160010154826200183f919062005b90565b915060010162001801565b50805f036200189c5760405162461bcd60e51b815260206004820152601c60248201527f54686520636f6d6d6974746565206973206e6f74207374616b696e6700000000604482015260640162000dfd565b5f83620018ab60638762005c1e565b620018b7919062005b90565b90505f81604051602001620018ce91815260200190565b60408051601f19818403018152919052805160209091012090505f620018f5848362005c4c565b90505f805b6021548110156200199557602181815481106200191b576200191b62005c0a565b905f5260205f209060020201600101548262001938919062005b90565b91506200194760018362005b7a565b83116200198c576021818154811062001964576200196462005c0a565b5f9182526020909120600290910201546001600160a01b0316965062000e3a95505050505050565b600101620018fa565b5060405162461bcd60e51b815260206004820152602960248201527f5468657265206973206e6f2076616c696461746f72206c65667420696e20746860448201526865206e6574776f726b60b81b606482015260840162000dfd565b6017546001600160a01b0316331462001a1e5760405162461bcd60e51b815260040162000dfd9062005a3b565b60185481101562001ad55780601f5462001a39919062005b90565b431062001ad55760405162461bcd60e51b815260206004820152605760248201527f63757272656e7420636861696e2068656164206578636565642074686520776960448201527f6e646f773a206c617374426c6f636b45706f6368202b205f6e6577506572696f60648201527f642c2074727920616761696e206c6174746572206f6e2e000000000000000000608482015260a40162000dfd565b6018819055601254604051631ad7d11360e21b8152600481018390526001600160a01b0390911690636b5f444c906024015f604051808303815f87803b15801562001b1e575f80fd5b505af115801562001b31573d5f803e3d5ffd5b505050507fd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f818160405162001b6791815260200190565b60405180910390a150565b6017546001600160a01b0316331462001b9f5760405162461bcd60e51b815260040162000dfd9062005a3b565b600d55565b6001600160a01b038083165f818152602760205260409020600101549091161462001be35760405162461bcd60e51b815260040162000dfd9062005a72565b6001600160a01b038281165f9081526027602052604090205416331462001c1e5760405162461bcd60e51b815260040162000dfd9062005aa9565b61271081111562001c725760405162461bcd60e51b815260206004820152601f60248201527f7265717569726520636f727265637420636f6d6d697373696f6e207261746500604482015260640162000dfd565b604080516060810182526001600160a01b038481168252436020808401918252838501868152600c80545f908152600a909352958220855181546001600160a01b03191695169490941784559151600180850191909155915160029093019290925583549293909290919062001cea90849062005b90565b90915550506040518281526001600160a01b038416907f4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf9060200160405180910390a2505050565b6017546001600160a01b0316331462001d5f5760405162461bcd60e51b815260040162000dfd9062005a3b565b6002805460ff19166001179055565b6017546001600160a01b0316331462001d9b5760405162461bcd60e51b815260040162000dfd9062005a3b565b5f811162001dec5760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000dfd565b601a55565b6017546001600160a01b0316331462001e1e5760405162461bcd60e51b815260040162000dfd9062005a3b565b6001600160a01b0382165f9081526026602052604090205481111562001e805760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000dfd565b6001600160a01b0382165f908152602660205260408120805483929062001ea990849062005b7a565b925050819055508060285f82825462001ec3919062005b7a565b90915550506040518181526001600160a01b038316907f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a39060200162001259565b6001600160a01b038083165f818152602760205260409020600101549091161462001f435760405162461bcd60e51b815260040162000dfd9062005af5565b6001600160a01b0382165f9081526027602052604081206012015460ff16600381111562001f755762001f75620054ed565b1462001fc45760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000dfd565b62001fd182823362003bc1565b5050565b6001600160a01b038083165f8181526027602052604090206001015490911614620020145760405162461bcd60e51b815260040162000dfd9062005af5565b5f81116200205d5760405162461bcd60e51b81526020600482015260156024820152740756e626f6e64696e6720616d6f756e74206973203605c1b604482015260640162000dfd565b62001fd182823362003da8565b60606024805480602002602001604051908101604052809291908181526020015f905b828210156200213f578382905f5260205f20018054620020ad9062005b2c565b80601f0160208091040260200160405190810160405280929190818152602001828054620020db9062005b2c565b80156200212a5780601f1062002100576101008083540402835291602001916200212a565b820191905f5260205f20905b8154815290600101906020018083116200210c57829003601f168201915b5050505050815260200190600101906200208d565b50505050905090565b5f62000e36338484620034ea565b60606021805480602002602001604051908101604052809291908181526020015f905b828210156200213f575f848152602090819020604080518082019091526002850290910180546001600160a01b0316825260019081015482840152908352909201910162002179565b5f604051806102600160405280336001600160a01b031681526020015f6001600160a01b03168152602001846001600160a01b03168152602001858152602001600d5f016002015481526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f81526020014381526020015f81526020015f81526020015f81526020015f60038111156200227b576200227b620054ed565b905290506200228b8183620040d1565b60208101516101808201516040517f8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c92620022cc92339288918a9162005c62565b60405180910390a150505050565b6029546060906001600160a01b031633146200230a5760405162461bcd60e51b815260040162000dfd9062005ba6565b601c546200235b5760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000dfd565b5f805b601c548110156200242c575f60275f601c848154811062002383576200238362005c0a565b5f9182526020808320909101546001600160a01b0316835282019290925260400190206012015460ff166003811115620023c157620023c1620054ed565b1480156200240d57505f60275f601c8481548110620023e457620023e462005c0a565b5f9182526020808320909101546001600160a01b03168352820192909252604001902060050154115b156200242357816200241f8162005caa565b9250505b6001016200235e565b50601a548181106200243b5750805b5f8267ffffffffffffffff811115620024585762002458620057bd565b6040519080825280602002602001820160405280156200249557816020015b620024816200524c565b815260200190600190039081620024775790505b5090505f8267ffffffffffffffff811115620024b557620024b5620057bd565b604051908082528060200260200182016040528015620024f257816020015b620024de6200524c565b815260200190600190039081620024d45790505b5090505f8367ffffffffffffffff811115620025125762002512620057bd565b6040519080825280602002602001820160405280156200253c578160200160208202803683370190505b5090505f805b601c548110156200281a575f60275f601c848154811062002567576200256762005c0a565b5f9182526020808320909101546001600160a01b0316835282019290925260400190206012015460ff166003811115620025a557620025a5620054ed565b148015620025f157505f60275f601c8481548110620025c857620025c862005c0a565b5f9182526020808320909101546001600160a01b03168352820192909252604001902060050154115b1562002811575f60275f601c848154811062002611576200261162005c0a565b5f9182526020808320909101546001600160a01b0390811684528382019490945260409283019091208251610260810184528154851681526001820154851692810192909252600281015490931691810191909152600382018054919291606084019190620026809062005b2c565b80601f0160208091040260200160405190810160405280929190818152602001828054620026ae9062005b2c565b8015620026fd5780601f10620026d357610100808354040283529160200191620026fd565b820191905f5260205f20905b815481529060010190602001808311620026df57829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546080820152600882015460a0820152600982015460c0820152600a82015460e0820152600b820154610100820152600c8201546001600160a01b0316610120820152600d820154610140820152600e820154610160820152600f82015461018082015260108201546101a082015260118201546101c082015260128201546101e09091019060ff166003811115620027c457620027c4620054ed565b6003811115620027d857620027d8620054ed565b81525050905080868481518110620027f457620027f462005c0a565b602002602001018190525082806200280c9062005caa565b935050505b60010162002542565b50601a54845111156200288c57620028328462004407565b5f5b601a54811015620028855784818151811062002854576200285462005c0a565b602002602001015184828151811062002871576200287162005c0a565b602090810291909101015260010162002834565b5062002890565b8392505b6200289d60215f62005303565b620028aa60245f62005323565b5f60208190555b8581101562002a81575f6040518060400160405280868481518110620028db57620028db62005c0a565b6020026020010151602001516001600160a01b031681526020018684815181106200290a576200290a62005c0a565b60209081029190910181015160a00151909152602180546001810182555f9190915282517f3a6357012c1a3ae0a17d304c9920310382d968ebcc4b1771f41c6b304205b570600290920291820180546001600160a01b0319166001600160a01b03909216919091179055908201517f3a6357012c1a3ae0a17d304c9920310382d968ebcc4b1771f41c6b304205b571909101558551909150602490869084908110620029ba57620029ba62005c0a565b6020908102919091018101516060015182546001810184555f938452919092200190620029e8908262005d0d565b50848281518110620029fe57620029fe62005c0a565b60200260200101516040015184838151811062002a1f5762002a1f62005c0a565b60200260200101906001600160a01b031690816001600160a01b03168152505084828151811062002a545762002a5462005c0a565b602002602001015160a0015160205f82825462002a72919062005b90565b909155505050600101620028b1565b50909550505050505090565b6017546001600160a01b0316331462002aba5760405162461bcd60e51b815260040162000dfd9062005a3b565b62002ac65f8362004423565b62001fd160018262004423565b6017546001600160a01b0316331462002b005760405162461bcd60e51b815260040162000dfd9062005a3b565b601580546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038082165f818152602760205260409020600101549091161462002b615760405162461bcd60e51b815260040162000dfd9062005a72565b6001600160a01b038082165f9081526027602052604090208054909116331462002be65760405162461bcd60e51b815260206004820152602f60248201527f726571756972652063616c6c657220746f2062652076616c696461746f72207460448201526e1c99585cdd5c9e481858d8dbdd5b9d608a1b606482015260840162000dfd565b5f601282015460ff16600381111562002c035762002c03620054ed565b0362002c525760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f7220616c7265616479206163746976650000000000000000604482015260640162000dfd565b6002601282015460ff16600381111562002c705762002c70620054ed565b14801562002c815750438160100154115b1562002cd05760405162461bcd60e51b815260206004820152601760248201527f76616c696461746f72207374696c6c20696e206a61696c000000000000000000604482015260640162000dfd565b6003601282015460ff16600381111562002cee5762002cee620054ed565b0362002d3d5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206a61696c6564207065726d616e656e746c7900000000604482015260640162000dfd565b60128101805460ff191690558054601854601f546001600160a01b038581169316917f60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b59162002d8d919062005b90565b60405190815260200160405180910390a35050565b6060805f600181805462002db69062005b2c565b80601f016020809104026020016040519081016040528092919081815260200182805462002de49062005b2c565b801562002e335780601f1062002e095761010080835404028352916020019162002e33565b820191905f5260205f20905b81548152906001019060200180831162002e1557829003601f168201915b5050505050915080805462002e489062005b2c565b80601f016020809104026020016040519081016040528092919081815260200182805462002e769062005b2c565b801562002ec55780601f1062002e9b5761010080835404028352916020019162002ec5565b820191905f5260205f20905b81548152906001019060200180831162002ea757829003601f168201915b50505050509050915091509091565b6060601c80548060200260200160405190810160405280929190818152602001828054801562002f2c57602002820191905f5260205f20905b81546001600160a01b0316815260019091019060200180831162002f0d575b5050505050905090565b6017546001600160a01b0316331462002f635760405162461bcd60e51b815260040162000dfd9062005a3b565b600e8190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200162001b67565b6017546001600160a01b0316331462002fc65760405162461bcd60e51b815260040162000dfd9062005a3b565b62002fd25f8062005340565b62002fdf60015f62005340565b6002805460ff19169055565b6017546001600160a01b03163314620030185760405162461bcd60e51b815260040162000dfd9062005a3b565b601680546001600160a01b0319166001600160a01b0392909216919091179055565b6017546001600160a01b03163314620030675760405162461bcd60e51b815260040162000dfd9062005a3b565b601480546001600160a01b0319166001600160a01b0392909216919091179055565b6029546001600160a01b03163314620030b65760405162461bcd60e51b815260040162000dfd9062005ba6565b620030c0620039c8565b62000ec5620022da565b6017546001600160a01b03163314620030f75760405162461bcd60e51b815260040162000dfd9062005a3b565b601180546001600160a01b0319166001600160a01b0392909216919091179055565b6012546001600160a01b03163314620031815760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f742074686520736c617368696e6720636f6e74726044820152621858dd60ea1b606482015260840162000dfd565b5f6101208201356027826200319d60408601602087016200548a565b6001600160a01b03166001600160a01b031681526020019081526020015f2060090154620031cc919062005b7a565b60c083013560275f620031e660408701602088016200548a565b6001600160a01b03166001600160a01b031681526020019081526020015f206006015462003215919062005b7a565b60a084013560275f6200322f60408801602089016200548a565b6001600160a01b03166001600160a01b031681526020019081526020015f20600501546200325e919062005b7a565b6200326a919062005b90565b62003276919062005b90565b6011546001600160a01b03165f90815260266020526040812080549293508392909190620032a690849062005b90565b9091555082905060275f620032c260408401602085016200548a565b6001600160a01b0316815260208101919091526040015f20620032e6828262005f58565b50505050565b6001600160a01b038316620033505760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000dfd565b6001600160a01b038216620033b35760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000dfd565b6001600160a01b038381165f8181526025602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b0381165f90815260276020526040812090601282015460ff166003811115620034485762003448620054ed565b14620034975760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265206163746976650000000000000000604482015260640162000dfd565b60128101805460ff191660011790558054601854601f546001600160a01b038581169316917f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c9162002d8d919062005b90565b6001600160a01b0383165f908152602660205260409020548111156200354c5760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000dfd565b6001600160a01b0383165f90815260266020526040812080548392906200357590849062005b7a565b90915550506001600160a01b0382165f9081526026602052604081208054839290620035a390849062005b90565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200340791815260200190565b475f03620035fa57565b600d5447905f90670de0b6b3a7640000906200361890849062005c1e565b6200362491906200608c565b905080156200369f576011546040515f916001600160a01b03169083908381818185875af1925050503d805f811462003679576040519150601f19603f3d011682016040523d82523d5f602084013e6200367e565b606091505b50909150508015156001036200369d576200369a828462005b7a565b92505b505b8160225f828254620036b2919062005b90565b909155505f90505b602154811015620039c3575f60275f60218481548110620036df57620036df62005c0a565b905f5260205f2090600202015f015f9054906101000a90046001600160a01b03166001600160a01b03166001600160a01b031681526020019081526020015f2090505f60205485602185815481106200373c576200373c62005c0a565b905f5260205f2090600202016001015462003758919062005c1e565b6200376491906200608c565b90508015620039b7576002601283015460ff1660038111156200378b576200378b620054ed565b1480620037b357506003601283015460ff166003811115620037b157620037b1620054ed565b145b156200385657601254602180546001600160a01b0390921691631de9d9b691849187908110620037e757620037e762005c0a565b5f91825260209091206002909102015460405160e084901b6001600160e01b03191681526001600160a01b0390911660048201526024015f604051808303818588803b15801562003836575f80fd5b505af115801562003849573d5f803e3d5ffd5b50505050505050620039ba565b5f82600501548284600801546200386e919062005c1e565b6200387a91906200608c565b90505f62003889828462005b7a565b90508115620038eb5783546040516001600160a01b03909116906108fc9084905f818181858888f193505050503d805f8114620038e2576040519150601f19603f3d011682016040523d82523d5f602084013e620038e7565b606091505b5050505b8015620039705783600c015f9054906101000a90046001600160a01b03166001600160a01b031663fb489a7b826040518263ffffffff1660e01b815260040160206040518083038185885af115801562003947573d5f803e3d5ffd5b50505050506040513d601f19601f820116820180604052508101906200396e9190620060a2565b505b60018401546040518481526001600160a01b03909116907fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe5639060200160405180910390a250505b50505b600101620036ba565b505050565b6004545b600554811015620039f557620039ef81620039e78162005caa565b925062004572565b620039cc565b506005546004556008546007540362003a0a57565b6009545b60085481101562003a375762003a318162003a298162005caa565b925062004782565b62003a0e565b50600854600955600754805b60085481101562003aa4576010545f82815260066020526040902060040154439162003a6f9162005b90565b1162003a955762003a808162004a67565b62003a8d60018362005b90565b915062003a9b565b62003aa4565b60010162003a43565b50600755565b600c54600b54101562003bbf57600b545f908152600a602052604090206010546001820154439162003adc9162005b90565b111562003ae65750565b600281015481546001600160a01b039081165f90815260276020526040808220600490810185905585548416835291819020600c015490516319fac8fd60e01b81529216926319fac8fd9262003b40920190815260200190565b5f604051808303815f87803b15801562003b58575f80fd5b505af115801562003b6b573d5f803e3d5ffd5b5050600b80545f908152600a6020526040812080546001600160a01b03191681556001808201839055600290910182905582549094509192509062003bb290849062005b90565b9091555062003aaa915050565b565b5f821162003c1e5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000dfd565b6001600160a01b0381165f9081526026602052604090205482111562003c875760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000dfd565b6001600160a01b0381165f908152602660205260408120805484929062003cb090849062005b7a565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580545f908152600394859052978820875181549088166001600160a01b031991821617825595516001820180549190981696169590951790955590516002840155519101558054919262003d3b8362005caa565b90915550506001600160a01b038481165f81815260276020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d91015b60405180910390a35050505050565b6001600160a01b038084165f908152602760205260409020805490918381169116148062003f1657600c820154604051631092ab9160e31b81526001600160a01b0385811660048301525f9216906384955c88906024016020604051808303815f875af115801562003e1c573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019062003e429190620060a2565b90508481101562003eaa5760405162461bcd60e51b815260206004820152602b60248201527f696e73756666696369656e7420756e6c6f636b6564204c6971756964204e657760448201526a746f6e2062616c616e636560a81b606482015260840162000dfd565b600c83015460405163282d3fdf60e01b81526001600160a01b038681166004830152602482018890529091169063282d3fdf906044015f604051808303815f87803b15801562003ef8575f80fd5b505af115801562003f0b573d5f803e3d5ffd5b505050505062003fa8565b8382600b0154836008015462003f2d919062005b7a565b101562003f8d5760405162461bcd60e51b815260206004820152602760248201527f696e73756666696369656e742073656c6620626f6e646564206e6577746f6e2060448201526662616c616e636560c81b606482015260840162000dfd565b8382600b015f82825462003fa2919062005b90565b90915550505b6040805160e0810182526001600160a01b03808616825287811660208084019182528385018981525f60608601818152436080880190815260a088018381528a151560c08a019081526008805486526006909752998420985189549089166001600160a01b0319918216178a55965160018a01805491909916971696909617909655915160028701559051600386015592516004850155905160059093018054945115156101000261ff00199415159490941661ffff1990951694909417929092179092558054916200407b8362005caa565b9190505550826001600160a01b0316856001600160a01b03167f63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc838760405162003d999291909115158252602082015260400190565b80516082146200411b5760405162461bcd60e51b8152602060048201526014602482015273092dcecc2d8d2c840e0e4dedecc40d8cadccee8d60631b604482015260640162000dfd565b620041268262004bab565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a0000000000006020808301919091528451925191925f9262004188920160609190911b6bffffffffffffffffffffffff1916815260140190565b60405160208183030381529060405290505f82620041a7835162004cd6565b83604051602001620041bc93929190620060ba565b60408051601f19818403018152828252805160209182012060028085526060850184529094505f93929091830190803683370190505090505f808060205b8851811015620042c35762004210898262004df3565b604080515f8152602081018083528b905260ff8316918101919091526060810184905260808101839052929650909450925060019060a0016020604051602081039080840390855afa15801562004269573d5f803e3d5ffd5b5050604051601f190151905085620042836041846200608c565b8151811062004296576200429662005c0a565b6001600160a01b0390921660209283029190910190910152620042bb60418262005b90565b9050620041fa565b5088602001516001600160a01b0316845f81518110620042e757620042e762005c0a565b60200260200101516001600160a01b031614620043595760405162461bcd60e51b815260206004820152602960248201527f496e76616c6964206e6f6465206b6579206f776e6572736869702070726f6f66604482015268081c1c9bdd9a59195960ba1b606482015260840162000dfd565b88604001516001600160a01b0316846001815181106200437d576200437d62005c0a565b60200260200101516001600160a01b031614620043f15760405162461bcd60e51b815260206004820152602b60248201527f496e76616c6964206f7261636c65206b6579206f776e6572736869702070726f60448201526a1bd9881c1c9bdd9a59195960aa1b606482015260840162000dfd565b620043fc8962004e29565b505050505050505050565b62000ec5815f600184516200441d919062005b7a565b62005052565b815460026001808316156101000203821604825180820160208110602084100160028114620044d05760018114620044f657865f526020840460205f2001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b80821015620044b957815184556001840193506020820191506200449a565b815191036101000a90819004029091555062004569565b60028302826020036101000a846020036101000a60208901510402018501875562004569565b865f526020840460205f2001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b8082101562004557578151845560018401935060208201915062004538565b815191036101000a9081900402909155505b50505050505050565b5f81815260036020818152604080842060018101546001600160a01b03168552602790925290922090601282015460ff166003811115620045b757620045b7620054ed565b036200464b57600282015482546001600160a01b03165f9081526026602052604081208054909190620045ec90849062005b90565b9091555050815460018301546002840154604080516001600160a01b03948516815293909216602084015282820152517f14d528d8a5c5378dc43eafa0ea18cefb7d35f17cad6e4ca6062a4f8a71ea3dd99181900360600190a1505050565b805482546001600160a01b0390811691161462004740575f80826008015483600501546200467a919062005b7a565b9050805f03620046915783600201549150620046b7565b80846002015484600d0154620046a8919062005c1e565b620046b491906200608c565b91505b600c83015484546040516340c10f1960e01b81526001600160a01b039182166004820152602481018590529116906340c10f19906044015f604051808303815f87803b15801562004706575f80fd5b505af115801562004719573d5f803e3d5ffd5b505050508183600d015f82825462004732919062005b90565b909155506200475f92505050565b8160020154816008015f82825462004759919062005b90565b90915550505b8160020154816005015f82825462004778919062005b90565b9091555050505050565b5f81815260066020908152604080832060018101546001600160a01b0316845260279092528220600582015491929091610100900460ff1662004966576002830154600c8301548454604051637eee288d60e01b81526001600160a01b03918216600482015260248101849052911690637eee288d906044015f604051808303815f87803b15801562004813575f80fd5b505af115801562004826573d5f803e3d5ffd5b50505050600c8301548454604051632770a7eb60e21b81526001600160a01b03918216600482015260248101849052911690639dc29fac906044015f604051808303815f87803b15801562004879575f80fd5b505af11580156200488c573d5f803e3d5ffd5b505050505f83600801548460050154620048a7919062005b7a565b600d850154909150620048bb828462005c1e565b620048c791906200608c565b92508184600d015f828254620048de919062005b7a565b909155505060068401545f03620048fc576003850183905562004924565b6006840154600785015462004912908562005c1e565b6200491e91906200608c565b60038601555b82846006015f82825462004939919062005b90565b909155505060038501546007850180545f906200495890849062005b90565b9091555062004a3692505050565b506002820154600882015481111562004980575060088101545b81600901545f03620049995760038301819055620049c1565b6009820154600a830154620049af908362005c1e565b620049bb91906200608c565b60038401555b80826009015f828254620049d6919062005b90565b90915550506003830154600a830180545f90620049f590849062005b90565b9250508190555080826008015f82825462004a11919062005b7a565b90915550506002830154600b830180545f9062004a3090849062005b7a565b90915550505b6005808401805460ff19166001179055820180548291905f9062004a5c90849062005b7a565b909155505050505050565b5f818152600660205260408120600381015490910362004a85575050565b60018101546001600160a01b03165f9081526027602052604081206005830154909190610100900460ff1662004b1e5781600701548260060154846003015462004ad0919062005c1e565b62004adc91906200608c565b905080826006015f82825462004af3919062005b7a565b909155505060038301546007830180545f9062004b1290849062005b7a565b9091555062004b819050565b81600a01548260090154846003015462004b39919062005c1e565b62004b4591906200608c565b905080826009015f82825462004b5c919062005b7a565b90915550506003830154600a830180545f9062004b7b90849062005b7a565b90915550505b82546001600160a01b03165f908152602660205260408120805483929062004a5c90849062005b90565b5f62004bbb82606001516200520f565b6001600160a01b0390911660208401529050801562004c0b5760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000dfd565b6020808301516001600160a01b039081165f9081526027909252604090912060010154161562004c7e5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000dfd565b6127108260800151111562001fd15760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000dfd565b6060815f0362004cfd5750506040805180820190915260018152600360fc1b602082015290565b815f5b811562004d2c578062004d138162005caa565b915062004d249050600a836200608c565b915062004d00565b5f8167ffffffffffffffff81111562004d495762004d49620057bd565b6040519080825280601f01601f19166020018201604052801562004d74576020820181803683370190505b5090505b841562004deb5762004d8c60018362005b7a565b915062004d9b600a8662005c4c565b62004da890603062005b90565b60f81b81838151811062004dc05762004dc062005c0a565b60200101906001600160f81b03191690815f1a90535062004de3600a866200608c565b945062004d78565b949350505050565b818101805160208201516040909201519091905f1a601b81101562004e225762004e1f601b8262006102565b90505b9250925092565b6101808101516001600160a01b031662004ebd57601c545f9062004e4d9062004cd6565b9050602a5f9054906101000a90046001600160a01b03168260200151835f015184608001518460405162004e81906200537c565b62004e919594939291906200611e565b604051809103905ff08015801562004eab573d5f803e3d5ffd5b506001600160a01b0316610180830152505b60208181018051601c8054600180820183555f9283527f0e4562a10381dec21b205ed72637e6b1b523bdd0e4d4d50af5cd23dd4500a21190910180546001600160a01b03199081166001600160a01b039586161790915584518416835260279095526040918290208651815487169085161781559351908401805486169184169190911790558401516002830180549094169116179091556060820151829190600382019062004f6e908262005d0d565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f8201556102008201516010820155610220820151601182015561024082015160128201805460ff19166001836003811115620050485762005048620054ed565b0217905550505050565b818180820362005063575050505050565b5f85600262005073878762006165565b6200507f91906200618e565b6200508b9087620061c0565b815181106200509e576200509e62005c0a565b602002602001015160a0015190505b818313620051db575b80868481518110620050cc57620050cc62005c0a565b602002602001015160a001511115620050f45782620050eb81620061ea565b935050620050b6565b85828151811062005109576200510962005c0a565b602002602001015160a00151811115620051325781620051298162006204565b925050620050f4565b818313620051d5578582815181106200514f576200514f62005c0a565b60200260200101518684815181106200516c576200516c62005c0a565b602002602001015187858151811062005189576200518962005c0a565b60200260200101888581518110620051a557620051a562005c0a565b6020026020010182905282905250508280620051c190620061ea565b9350508180620051d19062006204565b9250505b620050ad565b81851215620051f157620051f186868462005052565b8383121562005207576200520786848662005052565b505050505050565b5f806200521b6200538a565b5f60408286516020880160ff5afa62005232575f80fd5b508051602090910151600160601b90910494909350915050565b6040518061026001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b03168152602001606081526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6003811115620052fe57620052fe620054ed565b905290565b5080545f8255600202905f5260205f209081019062000ec59190620053a8565b5080545f8255905f5260205f209081019062000ec59190620053d3565b5080546200534e9062005b2c565b5f825580601f106200535e575050565b601f0160209004905f5260205f209081019062000ec59190620053f3565b611954806200622383390190565b60405180604001604052806002906020820280368337509192915050565b5b80821115620053cf5780546001600160a01b03191681555f6001820155600201620053a9565b5090565b80821115620053cf575f620053e9828262005340565b50600101620053d3565b5b80821115620053cf575f8155600101620053f4565b5f5b83811015620054255781810151838201526020016200540b565b50505f910152565b5f81518084526200544681602086016020860162005409565b601f01601f19169290920160200192915050565b602081525f6200546e60208301846200542d565b9392505050565b6001600160a01b038116811462000ec5575f80fd5b5f602082840312156200549b575f80fd5b81356200546e8162005475565b5f8060408385031215620054ba575f80fd5b8235620054c78162005475565b946020939093013593505050565b5f60208284031215620054e6575f80fd5b5035919050565b634e487b7160e01b5f52602160045260245ffd5b600481106200551e57634e487b7160e01b5f52602160045260245ffd5b9052565b602081526200553d6020820183516001600160a01b03169052565b5f60208301516200555960408401826001600160a01b03169052565b5060408301516001600160a01b03811660608401525060608301516102608060808501526200558d6102808501836200542d565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a062005615818701836001600160a01b03169052565b8601516101c0868101919091528601516101e0808701919091528601516102008087019190915286015161022080870191909152860151610240808701919091528601519050620056698286018262005501565b5090949350505050565b5f805f6060848603121562005686575f80fd5b8335620056938162005475565b92506020840135620056a58162005475565b929592945050506040919091013590565b5f815180845260208085019450602084015f5b83811015620056fb57815180516001600160a01b031688528301518388015260409096019590820190600101620056c9565b509495945050505050565b8215158152604060208201525f62004deb6040830184620056b6565b5f806040838503121562005734575f80fd5b50508035926020909101359150565b5f60208083016020845280855180835260408601915060408160051b8701019250602087015f5b828110156200579c57603f19888603018452620057898583516200542d565b945092850192908501906001016200576a565b5092979650505050505050565b602081525f6200546e6020830184620056b6565b634e487b7160e01b5f52604160045260245ffd5b5f82601f830112620057e1575f80fd5b813567ffffffffffffffff80821115620057ff57620057ff620057bd565b604051601f8301601f19908116603f011681019082821181831017156200582a576200582a620057bd565b8160405283815286602085880101111562005843575f80fd5b836020870160208301375f602085830101528094505050505092915050565b5f805f6060848603121562005875575f80fd5b833567ffffffffffffffff808211156200588d575f80fd5b6200589b87838801620057d1565b945060208601359150620058af8262005475565b90925060408501359080821115620058c5575f80fd5b50620058d486828701620057d1565b9150509250925092565b602080825282518282018190525f9190848201906040850190845b81811015620059205783516001600160a01b031683529284019291840191600101620058f9565b50909695505050505050565b5f80604083850312156200593e575f80fd5b823567ffffffffffffffff8082111562005956575f80fd5b6200596486838701620057d1565b935060208501359150808211156200597a575f80fd5b506200598985828601620057d1565b9150509250929050565b604081525f620059a760408301856200542d565b8281036020840152620059bb81856200542d565b95945050505050565b5f8060408385031215620059d6575f80fd5b8235620059e38162005475565b91506020830135620059f58162005475565b809150509250929050565b5f6020828403121562005a11575f80fd5b813567ffffffffffffffff81111562005a28575f80fd5b820161026081850312156200546e575f80fd5b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b6020808252601c908201527f76616c696461746f72206d757374206265207265676973746572656400000000604082015260600190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b60208082526018908201527f76616c696461746f72206e6f7420726567697374657265640000000000000000604082015260600190565b600181811c9082168062005b4157607f821691505b60208210810362005b6057634e487b7160e01b5f52602260045260245ffd5b50919050565b634e487b7160e01b5f52601160045260245ffd5b8181038181111562000e3a5762000e3a62005b66565b8082018082111562000e3a5762000e3a62005b66565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b5f6020828403121562005bfa575f80fd5b815180151581146200546e575f80fd5b634e487b7160e01b5f52603260045260245ffd5b808202811582820484141762000e3a5762000e3a62005b66565b634e487b7160e01b5f52601260045260245ffd5b5f8262005c5d5762005c5d62005c38565b500690565b5f60018060a01b0380881683528087166020840152808616604084015260a0606084015262005c9560a08401866200542d565b91508084166080840152509695505050505050565b5f6001820162005cbe5762005cbe62005b66565b5060010190565b601f821115620039c357805f5260205f20601f840160051c8101602085101562005cec5750805b601f840160051c820191505b818110156200135e575f815560010162005cf8565b815167ffffffffffffffff81111562005d2a5762005d2a620057bd565b62005d428162005d3b845462005b2c565b8462005cc5565b602080601f83116001811462005d78575f841562005d605750858301515b5f19600386901b1c1916600185901b17855562005207565b5f85815260208120601f198616915b8281101562005da85788860151825594840194600190910190840162005d87565b508582101562005dc657878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b5f813562000e3a8162005475565b80546001600160a01b0319166001600160a01b0392909216919091179055565b5f808335601e1984360301811262005e1a575f80fd5b83018035915067ffffffffffffffff82111562005e35575f80fd5b60200191503681900382131562005e4a575f80fd5b9250929050565b67ffffffffffffffff83111562005e6c5762005e6c620057bd565b62005e848362005e7d835462005b2c565b8362005cc5565b5f601f84116001811462005eb8575f851562005ea05750838201355b5f19600387901b1c1916600186901b1783556200135e565b5f83815260208120601f198716915b8281101562005ee9578685013582556020948501946001909201910162005ec7565b508682101562005f06575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b5f81356004811062000e3a575f80fd5b6004821062005f4557634e487b7160e01b5f52602160045260245ffd5b60ff1981541660ff831681178255505050565b62005f6e62005f678362005dd6565b8262005de4565b62005f8a62005f806020840162005dd6565b6001830162005de4565b62005fa662005f9c6040840162005dd6565b6002830162005de4565b62005fb5606083018362005e04565b62005fc581836003860162005e51565b50506080820135600482015560a0820135600582015560c0820135600682015560e0820135600782015561010082013560088201556101208201356009820155610140820135600a820155610160820135600b820155620060386200602e610180840162005dd6565b600c830162005de4565b6101a0820135600d8201556101c0820135600e8201556101e0820135600f8201556102008201356010820155610220820135601182015562001fd162006082610240840162005f18565b6012830162005f28565b5f826200609d576200609d62005c38565b500490565b5f60208284031215620060b3575f80fd5b5051919050565b5f8451620060cd81846020890162005409565b845190830190620060e381836020890162005409565b8451910190620060f881836020880162005409565b0195945050505050565b60ff818116838216019081111562000e3a5762000e3a62005b66565b6001600160a01b0386811682528581166020830152841660408201526060810183905260a0608082018190525f906200615a908301846200542d565b979650505050505050565b8181035f83128015838313168383128216171562006187576200618762005b66565b5092915050565b5f826200619f576200619f62005c38565b600160ff1b82145f1984141615620061bb57620061bb62005b66565b500590565b8082018281125f831280158216821582161715620061e257620061e262005b66565b505092915050565b5f6001600160ff1b01820162005cbe5762005cbe62005b66565b5f600160ff1b82016200621b576200621b62005b66565b505f19019056fe608060405234801562000010575f80fd5b5060405162001954380380620019548339810160408190526200003391620001bc565b5f856001600160a01b0316632f2c3f2e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000071573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620000979190620002ba565b905080831115620000a6575f80fd5b600d80546001600160a01b038089166001600160a01b031992831617909255600c805488841690831617905560088054928716929091169190911790556009839055604051620000fb908390602001620002d2565b604051602081830303815290604052600a90816200011a91906200038a565b50816040516020016200012e9190620002d2565b604051602081830303815290604052600b90816200014d91906200038a565b50505f80546001600160a01b0319163317905550620004569350505050565b6001600160a01b038116811462000181575f80fd5b50565b634e487b7160e01b5f52604160045260245ffd5b5f5b83811015620001b45781810151838201526020016200019a565b50505f910152565b5f805f805f60a08688031215620001d1575f80fd5b8551620001de816200016c565b6020870151909550620001f1816200016c565b604087015190945062000204816200016c565b6060870151608088015191945092506001600160401b038082111562000228575f80fd5b818801915088601f8301126200023c575f80fd5b81518181111562000251576200025162000184565b604051601f8201601f19908116603f011681019083821181831017156200027c576200027c62000184565b816040528281528b602084870101111562000295575f80fd5b620002a883602083016020880162000198565b80955050505050509295509295909350565b5f60208284031215620002cb575f80fd5b5051919050565b644c4e544e2d60d81b81525f8251620002f381600585016020870162000198565b9190910160050192915050565b600181811c908216806200031557607f821691505b6020821081036200033457634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200038557805f5260205f20601f840160051c81016020851015620003615750805b601f840160051c820191505b8181101562000382575f81556001016200036d565b50505b505050565b81516001600160401b03811115620003a657620003a662000184565b620003be81620003b7845462000300565b846200033a565b602080601f831160018114620003f4575f8415620003dc5750858301515b5f19600386901b1c1916600185901b1785556200044e565b5f85815260208120601f198616915b82811015620004245788860151825594840194600190910190840162000403565b50858210156200044257878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b6114f080620004645f395ff3fe60806040526004361061013c575f3560e01c806359355736116100b3578063949813b81161006d578063949813b81461039457806395d89b41146103b35780639dc29fac146103c7578063a9059cbb146103e6578063dd62ed3e14610405578063fb489a7b14610449575f80fd5b806359355736146102ba5780635ea1d6f8146102ee57806361d027b31461030357806370a08231146103225780637eee288d1461035657806384955c8814610375575f80fd5b806323b872dd1161010457806323b872dd1461020f578063282d3fdf1461022e578063313ce5671461024d578063372500ab146102685780633a5381b51461027c57806340c10f191461029b575f80fd5b806306fdde0314610140578063095ea7b31461016a5780630c3c6fee1461019957806318160ddd146101d057806319fac8fd146101ee575b5f80fd5b34801561014b575f80fd5b50610154610451565b60405161016191906112d3565b60405180910390f35b348015610175575f80fd5b50610189610184366004611320565b6104e1565b6040519015158152602001610161565b3480156101a4575f80fd5b50600d546101b8906001600160a01b031681565b6040516001600160a01b039091168152602001610161565b3480156101db575f80fd5b506004545b604051908152602001610161565b3480156101f9575f80fd5b5061020d610208366004611348565b610606565b005b34801561021a575f80fd5b5061018961022936600461135f565b610740565b348015610239575f80fd5b5061020d610248366004611320565b610877565b348015610258575f80fd5b5060405160128152602001610161565b348015610273575f80fd5b5061020d6109a6565b348015610287575f80fd5b50600c546101b8906001600160a01b031681565b3480156102a6575f80fd5b5061020d6102b5366004611320565b610aa3565b3480156102c5575f80fd5b506101e06102d4366004611398565b6001600160a01b03165f9081526002602052604090205490565b3480156102f9575f80fd5b506101e060095481565b34801561030e575f80fd5b506008546101b8906001600160a01b031681565b34801561032d575f80fd5b506101e061033c366004611398565b6001600160a01b03165f9081526001602052604090205490565b348015610361575f80fd5b5061020d610370366004611320565b610bcc565b348015610380575f80fd5b506101e061038f366004611398565b610cf5565b34801561039f575f80fd5b506101e06103ae366004611398565b610e20565b3480156103be575f80fd5b50610154610f2e565b3480156103d2575f80fd5b5061020d6103e1366004611320565b610f3d565b3480156103f1575f80fd5b50610189610400366004611320565b611066565b348015610410575f80fd5b506101e061041f3660046113b8565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b6101e0611173565b6060600a8054610460906113e9565b80601f016020809104026020016040519081016040528092919081815260200182805461048c906113e9565b80156104d75780601f106104ae576101008083540402835291602001916104d7565b820191905f5260205f20905b8154815290600101906020018083116104ba57829003601f168201915b5050505050905090565b600d546040516001600160a01b038481166024830152604482018490525f92839283929091169060640160408051601f198184030181529181526020820180516001600160e01b031663095ea7b360e01b179052516105409190611421565b5f60405180830381855af49150503d805f8114610578576040519150601f19603f3d011682016040523d82523d5f602084013e61057d565b606091505b5090925090508115155f036105e95780511561059c5780518082602001fd5b60405162461bcd60e51b815260206004820181905260248201527f46756e6374696f6e2063616c6c2027617070726f76652720726576657274656460448201526064015b60405180910390fd5b808060200190518101906105fd919061143c565b95945050505050565b5f546001600160a01b0316331461062f5760405162461bcd60e51b81526004016105e09061145b565b600d54604051602481018390525f9182916001600160a01b039091169060440160408051601f198184030181529181526020820180516001600160e01b03166319fac8fd60e01b179052516106849190611421565b5f60405180830381855af49150503d805f81146106bc576040519150601f19603f3d011682016040523d82523d5f602084013e6106c1565b606091505b5090925090508115155f0361073b578051156106e05780518082602001fd5b60405162461bcd60e51b815260206004820152602a60248201527f46756e6374696f6e2063616c6c2027736574436f6d6d697373696f6e5261746560448201526909c81c995d995c9d195960b21b60648201526084016105e0565b505050565b600d546040516001600160a01b0385811660248301528481166044830152606482018490525f92839283929091169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516107a79190611421565b5f60405180830381855af49150503d805f81146107df576040519150601f19603f3d011682016040523d82523d5f602084013e6107e4565b606091505b5090925090508115155f03610859578051156108035780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c20277472616e7366657246726f6d2720726576604482015264195c9d195960da1b60648201526084016105e0565b8080602001905181019061086d919061143c565b9695505050505050565b5f546001600160a01b031633146108a05760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b031663282d3fdf60e01b179052516108fc9190611421565b5f60405180830381855af49150503d805f8114610934576040519150601f19603f3d011682016040523d82523d5f602084013e610939565b606091505b5090925090508115155f036109a0578051156109585780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276c6f636b2720726576657274656400000060448201526064016105e0565b50505050565b600d5460408051600481526024810182526020810180516001600160e01b031663372500ab60e01b17905290515f9283926001600160a01b03909116916109ed9190611421565b5f60405180830381855af49150503d805f8114610a25576040519150601f19603f3d011682016040523d82523d5f602084013e610a2a565b606091505b5090925090508115155f03610a9f57805115610a495780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c2027636c61696d526577617264732720726576604482015264195c9d195960da1b60648201526084016105e0565b5050565b5f546001600160a01b03163314610acc5760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b03166340c10f1960e01b17905251610b289190611421565b5f60405180830381855af49150503d805f8114610b60576040519150601f19603f3d011682016040523d82523d5f602084013e610b65565b606091505b5090925090508115155f036109a057805115610b845780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276d696e742720726576657274656400000060448201526064016105e0565b5f546001600160a01b03163314610bf55760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b0316637eee288d60e01b17905251610c519190611421565b5f60405180830381855af49150503d805f8114610c89576040519150601f19603f3d011682016040523d82523d5f602084013e610c8e565b606091505b5090925090508115155f036109a057805115610cad5780518082602001fd5b60405162461bcd60e51b815260206004820152601f60248201527f46756e6374696f6e2063616c6c2027756e6c6f636b272072657665727465640060448201526064016105e0565b600d546040516001600160a01b0383811660248301525f92839283929091169060440160408051601f198184030181529181526020820180516001600160e01b0316631092ab9160e31b17905251610d4d9190611421565b5f60405180830381855af49150503d805f8114610d85576040519150601f19603f3d011682016040523d82523d5f602084013e610d8a565b606091505b5090925090508115155f03610e0457805115610da95780518082602001fd5b60405162461bcd60e51b815260206004820152602a60248201527f46756e6374696f6e2063616c6c2027756e6c6f636b656442616c616e63654f6660448201526909c81c995d995c9d195960b21b60648201526084016105e0565b80806020019051810190610e1891906114a3565b949350505050565b600d546040516001600160a01b0383811660248301525f92839283929091169060440160408051601f198184030181529181526020820180516001600160e01b0316631293027760e31b17905251610e789190611421565b5f60405180830381855af49150503d805f8114610eb0576040519150601f19603f3d011682016040523d82523d5f602084013e610eb5565b606091505b5090925090508115155f03610e0457805115610ed45780518082602001fd5b60405162461bcd60e51b815260206004820152602960248201527f46756e6374696f6e2063616c6c2027756e636c61696d65645265776172647327604482015268081c995d995c9d195960ba1b60648201526084016105e0565b6060600b8054610460906113e9565b5f546001600160a01b03163314610f665760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b0316632770a7eb60e21b17905251610fc29190611421565b5f60405180830381855af49150503d805f8114610ffa576040519150601f19603f3d011682016040523d82523d5f602084013e610fff565b606091505b5090925090508115155f036109a05780511561101e5780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276275726e2720726576657274656400000060448201526064016105e0565b600d546040516001600160a01b038481166024830152604482018490525f92839283929091169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b179052516110c59190611421565b5f60405180830381855af49150503d805f81146110fd576040519150601f19603f3d011682016040523d82523d5f602084013e611102565b606091505b5090925090508115155f036105e9578051156111215780518082602001fd5b60405162461bcd60e51b815260206004820152602160248201527f46756e6374696f6e2063616c6c20277472616e736665722720726576657274656044820152601960fa1b60648201526084016105e0565b5f80546001600160a01b0316331461119d5760405162461bcd60e51b81526004016105e09061145b565b600d5460408051600481526024810182526020810180516001600160e01b031663fb489a7b60e01b17905290515f9283926001600160a01b03909116916111e49190611421565b5f60405180830381855af49150503d805f811461121c576040519150601f19603f3d011682016040523d82523d5f602084013e611221565b606091505b5090925090508115155f03611296578051156112405780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c20277265646973747269627574652720726576604482015264195c9d195960da1b60648201526084016105e0565b808060200190518101906112aa91906114a3565b9250505090565b5f5b838110156112cb5781810151838201526020016112b3565b50505f910152565b602081525f82518060208401526112f18160408501602087016112b1565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461131b575f80fd5b919050565b5f8060408385031215611331575f80fd5b61133a83611305565b946020939093013593505050565b5f60208284031215611358575f80fd5b5035919050565b5f805f60608486031215611371575f80fd5b61137a84611305565b925061138860208501611305565b9150604084013590509250925092565b5f602082840312156113a8575f80fd5b6113b182611305565b9392505050565b5f80604083850312156113c9575f80fd5b6113d283611305565b91506113e060208401611305565b90509250929050565b600181811c908216806113fd57607f821691505b60208210810361141b57634e487b7160e01b5f52602260045260245ffd5b50919050565b5f82516114328184602087016112b1565b9190910192915050565b5f6020828403121561144c575f80fd5b815180151581146113b1575f80fd5b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b5f602082840312156114b3575f80fd5b505191905056fea26469706673582212200f950bd4f4d9457214149c5d71464d2a2df907f4dd143e5feaf92aaf3d4852a864736f6c63430008170033a2646970667358221220c0004f911e6c33f2ab64113bc3e58a2facccf10d389997754737b49555abfa2064736f6c63430008170033608060405234801561000f575f80fd5b50604051610f80380380610f8083398101604081905261002e91610052565b5f80546001600160a01b0319166001600160a01b039290921691909117905561007f565b5f60208284031215610062575f80fd5b81516001600160a01b0381168114610078575f80fd5b9392505050565b610ef48061008c5f395ff3fe608060405260043610610110575f3560e01c80635ea1d6f81161009d578063949813b811610062578063949813b8146102e25780639dc29fac14610301578063a9059cbb14610320578063dd62ed3e1461033f578063fb489a7b14610383575f80fd5b80635ea1d6f81461022457806361d027b31461023957806370a08231146102705780637eee288d146102a457806384955c88146102c3575f80fd5b806323b872dd116100e357806323b872dd1461019e578063282d3fdf146101bd5780632f2c3f2e146101dc578063372500ab146101f157806340c10f1914610205575f80fd5b8063095ea7b31461011457806318160ddd14610148578063187cf4d71461016657806319fac8fd1461017d575b5f80fd5b34801561011f575f80fd5b5061013361012e366004610d1d565b61038b565b60405190151581526020015b60405180910390f35b348015610153575f80fd5b506004545b60405190815260200161013f565b348015610171575f80fd5b50610158633b9aca0081565b348015610188575f80fd5b5061019c610197366004610d45565b6103a1565b005b3480156101a9575f80fd5b506101336101b8366004610d5c565b6103d8565b3480156101c8575f80fd5b5061019c6101d7366004610d1d565b6104c9565b3480156101e7575f80fd5b5061015861271081565b3480156101fc575f80fd5b5061019c6105ab565b348015610210575f80fd5b5061019c61021f366004610d1d565b610655565b34801561022f575f80fd5b5061015860095481565b348015610244575f80fd5b50600854610258906001600160a01b031681565b6040516001600160a01b03909116815260200161013f565b34801561027b575f80fd5b5061015861028a366004610d95565b6001600160a01b03165f9081526001602052604090205490565b3480156102af575f80fd5b5061019c6102be366004610d1d565b6106ba565b3480156102ce575f80fd5b506101586102dd366004610d95565b61077d565b3480156102ed575f80fd5b506101586102fc366004610d95565b6107aa565b34801561030c575f80fd5b5061019c61031b366004610d1d565b6107d6565b34801561032b575f80fd5b5061013361033a366004610d1d565b610833565b34801561034a575f80fd5b50610158610359366004610db5565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b61015861087e565b5f6103973384846109df565b5060015b92915050565b5f546001600160a01b031633146103d35760405162461bcd60e51b81526004016103ca90610de6565b60405180910390fd5b600955565b6001600160a01b0383165f9081526003602090815260408083203384529091528120548281101561045c5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b60648201526084016103ca565b610470853361046b8685610e42565b6109df565b61047a8584610b02565b6104848484610be6565b836001600160a01b0316856001600160a01b03165f80516020610e9f833981519152856040516104b691815260200190565b60405180910390a3506001949350505050565b5f546001600160a01b031633146104f25760405162461bcd60e51b81526004016103ca90610de6565b6001600160a01b0382165f90815260026020908152604080832054600190925290912054829161052191610e42565b101561057b5760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b60648201526084016103ca565b6001600160a01b0382165f90815260026020526040812080548392906105a2908490610e55565b90915550505050565b5f6105b533610c2f565b335f81815260056020526040808220829055519293509183908381818185875af1925050503d805f8114610604576040519150601f19603f3d011682016040523d82523d5f602084013e610609565b606091505b50509050806106515760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b60448201526064016103ca565b5050565b5f546001600160a01b0316331461067e5760405162461bcd60e51b81526004016103ca90610de6565b6106888282610be6565b6040518181526001600160a01b038316905f905f80516020610e9f833981519152906020015b60405180910390a35050565b5f546001600160a01b031633146106e35760405162461bcd60e51b81526004016103ca90610de6565b6001600160a01b0382165f908152600260205260409020548111156107565760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b60648201526084016103ca565b6001600160a01b0382165f90815260026020526040812080548392906105a2908490610e42565b6001600160a01b0381165f90815260026020908152604080832054600190925282205461039b9190610e42565b5f6107b482610c91565b6001600160a01b0383165f9081526005602052604090205461039b9190610e55565b5f546001600160a01b031633146107ff5760405162461bcd60e51b81526004016103ca90610de6565b6108098282610b02565b6040518181525f906001600160a01b038416905f80516020610e9f833981519152906020016106ae565b5f61083e3383610b02565b6108488383610be6565b6040518281526001600160a01b0384169033905f80516020610e9f8339815191529060200160405180910390a350600192915050565b5f80546001600160a01b031633146108a85760405162461bcd60e51b81526004016103ca90610de6565b60095434905f90612710906108bd9084610e68565b6108c79190610e7f565b9050818111156109195760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f7220726577617264000000000000000060448201526064016103ca565b6109238183610e42565b6008546040519193506001600160a01b0316906108fc9083905f818181858888f193505050503d805f8114610973576040519150601f19603f3d011682016040523d82523d5f602084013e610978565b606091505b50506004545f915061098e633b9aca0085610e68565b6109989190610e7f565b9050806007546109a89190610e55565b6007556004545f90633b9aca00906109c09084610e68565b6109ca9190610e7f565b90506109d68184610e55565b94505050505090565b6001600160a01b038316610a415760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b60648201526084016103ca565b6001600160a01b038216610aa25760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b60648201526084016103ca565b6001600160a01b038381165f8181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610b0b82610c2f565b506001600160a01b0382165f90815260016020908152604080832054600290925290912054610b3a9082610e42565b821115610b895760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e6473000000000060448201526064016103ca565b610b938282610e42565b6001600160a01b0384165f90815260016020526040902055808203610bcb576001600160a01b0383165f908152600660205260408120555b8160045f828254610bdc9190610e42565b9091555050505050565b610bef82610c2f565b506001600160a01b0382165f9081526001602052604081208054839290610c17908490610e55565b925050819055508060045f8282546105a29190610e55565b5f80610c3a83610c91565b6001600160a01b0384165f90815260056020526040902054909150610c60908290610e55565b6001600160a01b039093165f9081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b0381165f90815260016020526040812054808203610cb857505f92915050565b6001600160a01b0383165f90815260066020526040812054600754610cdd9190610e42565b90505f633b9aca00610cef8484610e68565b610cf99190610e7f565b95945050505050565b80356001600160a01b0381168114610d18575f80fd5b919050565b5f8060408385031215610d2e575f80fd5b610d3783610d02565b946020939093013593505050565b5f60208284031215610d55575f80fd5b5035919050565b5f805f60608486031215610d6e575f80fd5b610d7784610d02565b9250610d8560208501610d02565b9150604084013590509250925092565b5f60208284031215610da5575f80fd5b610dae82610d02565b9392505050565b5f8060408385031215610dc6575f80fd5b610dcf83610d02565b9150610ddd60208401610d02565b90509250929050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561039b5761039b610e2e565b8082018082111561039b5761039b610e2e565b808202811582820484141761039b5761039b610e2e565b5f82610e9957634e487b7160e01b5f52601260045260245ffd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220a9171ec1ffebc6ff9e96eb0c1adf7b9e4e81605dc47e7a79c63c4befebd2e08064736f6c63430008170033608060405234801562000010575f80fd5b5060405162001954380380620019548339810160408190526200003391620001bc565b5f856001600160a01b0316632f2c3f2e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000071573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620000979190620002ba565b905080831115620000a6575f80fd5b600d80546001600160a01b038089166001600160a01b031992831617909255600c805488841690831617905560088054928716929091169190911790556009839055604051620000fb908390602001620002d2565b604051602081830303815290604052600a90816200011a91906200038a565b50816040516020016200012e9190620002d2565b604051602081830303815290604052600b90816200014d91906200038a565b50505f80546001600160a01b0319163317905550620004569350505050565b6001600160a01b038116811462000181575f80fd5b50565b634e487b7160e01b5f52604160045260245ffd5b5f5b83811015620001b45781810151838201526020016200019a565b50505f910152565b5f805f805f60a08688031215620001d1575f80fd5b8551620001de816200016c565b6020870151909550620001f1816200016c565b604087015190945062000204816200016c565b6060870151608088015191945092506001600160401b038082111562000228575f80fd5b818801915088601f8301126200023c575f80fd5b81518181111562000251576200025162000184565b604051601f8201601f19908116603f011681019083821181831017156200027c576200027c62000184565b816040528281528b602084870101111562000295575f80fd5b620002a883602083016020880162000198565b80955050505050509295509295909350565b5f60208284031215620002cb575f80fd5b5051919050565b644c4e544e2d60d81b81525f8251620002f381600585016020870162000198565b9190910160050192915050565b600181811c908216806200031557607f821691505b6020821081036200033457634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200038557805f5260205f20601f840160051c81016020851015620003615750805b601f840160051c820191505b8181101562000382575f81556001016200036d565b50505b505050565b81516001600160401b03811115620003a657620003a662000184565b620003be81620003b7845462000300565b846200033a565b602080601f831160018114620003f4575f8415620003dc5750858301515b5f19600386901b1c1916600185901b1785556200044e565b5f85815260208120601f198616915b82811015620004245788860151825594840194600190910190840162000403565b50858210156200044257878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b6114f080620004645f395ff3fe60806040526004361061013c575f3560e01c806359355736116100b3578063949813b81161006d578063949813b81461039457806395d89b41146103b35780639dc29fac146103c7578063a9059cbb146103e6578063dd62ed3e14610405578063fb489a7b14610449575f80fd5b806359355736146102ba5780635ea1d6f8146102ee57806361d027b31461030357806370a08231146103225780637eee288d1461035657806384955c8814610375575f80fd5b806323b872dd1161010457806323b872dd1461020f578063282d3fdf1461022e578063313ce5671461024d578063372500ab146102685780633a5381b51461027c57806340c10f191461029b575f80fd5b806306fdde0314610140578063095ea7b31461016a5780630c3c6fee1461019957806318160ddd146101d057806319fac8fd146101ee575b5f80fd5b34801561014b575f80fd5b50610154610451565b60405161016191906112d3565b60405180910390f35b348015610175575f80fd5b50610189610184366004611320565b6104e1565b6040519015158152602001610161565b3480156101a4575f80fd5b50600d546101b8906001600160a01b031681565b6040516001600160a01b039091168152602001610161565b3480156101db575f80fd5b506004545b604051908152602001610161565b3480156101f9575f80fd5b5061020d610208366004611348565b610606565b005b34801561021a575f80fd5b5061018961022936600461135f565b610740565b348015610239575f80fd5b5061020d610248366004611320565b610877565b348015610258575f80fd5b5060405160128152602001610161565b348015610273575f80fd5b5061020d6109a6565b348015610287575f80fd5b50600c546101b8906001600160a01b031681565b3480156102a6575f80fd5b5061020d6102b5366004611320565b610aa3565b3480156102c5575f80fd5b506101e06102d4366004611398565b6001600160a01b03165f9081526002602052604090205490565b3480156102f9575f80fd5b506101e060095481565b34801561030e575f80fd5b506008546101b8906001600160a01b031681565b34801561032d575f80fd5b506101e061033c366004611398565b6001600160a01b03165f9081526001602052604090205490565b348015610361575f80fd5b5061020d610370366004611320565b610bcc565b348015610380575f80fd5b506101e061038f366004611398565b610cf5565b34801561039f575f80fd5b506101e06103ae366004611398565b610e20565b3480156103be575f80fd5b50610154610f2e565b3480156103d2575f80fd5b5061020d6103e1366004611320565b610f3d565b3480156103f1575f80fd5b50610189610400366004611320565b611066565b348015610410575f80fd5b506101e061041f3660046113b8565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b6101e0611173565b6060600a8054610460906113e9565b80601f016020809104026020016040519081016040528092919081815260200182805461048c906113e9565b80156104d75780601f106104ae576101008083540402835291602001916104d7565b820191905f5260205f20905b8154815290600101906020018083116104ba57829003601f168201915b5050505050905090565b600d546040516001600160a01b038481166024830152604482018490525f92839283929091169060640160408051601f198184030181529181526020820180516001600160e01b031663095ea7b360e01b179052516105409190611421565b5f60405180830381855af49150503d805f8114610578576040519150601f19603f3d011682016040523d82523d5f602084013e61057d565b606091505b5090925090508115155f036105e95780511561059c5780518082602001fd5b60405162461bcd60e51b815260206004820181905260248201527f46756e6374696f6e2063616c6c2027617070726f76652720726576657274656460448201526064015b60405180910390fd5b808060200190518101906105fd919061143c565b95945050505050565b5f546001600160a01b0316331461062f5760405162461bcd60e51b81526004016105e09061145b565b600d54604051602481018390525f9182916001600160a01b039091169060440160408051601f198184030181529181526020820180516001600160e01b03166319fac8fd60e01b179052516106849190611421565b5f60405180830381855af49150503d805f81146106bc576040519150601f19603f3d011682016040523d82523d5f602084013e6106c1565b606091505b5090925090508115155f0361073b578051156106e05780518082602001fd5b60405162461bcd60e51b815260206004820152602a60248201527f46756e6374696f6e2063616c6c2027736574436f6d6d697373696f6e5261746560448201526909c81c995d995c9d195960b21b60648201526084016105e0565b505050565b600d546040516001600160a01b0385811660248301528481166044830152606482018490525f92839283929091169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516107a79190611421565b5f60405180830381855af49150503d805f81146107df576040519150601f19603f3d011682016040523d82523d5f602084013e6107e4565b606091505b5090925090508115155f03610859578051156108035780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c20277472616e7366657246726f6d2720726576604482015264195c9d195960da1b60648201526084016105e0565b8080602001905181019061086d919061143c565b9695505050505050565b5f546001600160a01b031633146108a05760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b031663282d3fdf60e01b179052516108fc9190611421565b5f60405180830381855af49150503d805f8114610934576040519150601f19603f3d011682016040523d82523d5f602084013e610939565b606091505b5090925090508115155f036109a0578051156109585780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276c6f636b2720726576657274656400000060448201526064016105e0565b50505050565b600d5460408051600481526024810182526020810180516001600160e01b031663372500ab60e01b17905290515f9283926001600160a01b03909116916109ed9190611421565b5f60405180830381855af49150503d805f8114610a25576040519150601f19603f3d011682016040523d82523d5f602084013e610a2a565b606091505b5090925090508115155f03610a9f57805115610a495780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c2027636c61696d526577617264732720726576604482015264195c9d195960da1b60648201526084016105e0565b5050565b5f546001600160a01b03163314610acc5760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b03166340c10f1960e01b17905251610b289190611421565b5f60405180830381855af49150503d805f8114610b60576040519150601f19603f3d011682016040523d82523d5f602084013e610b65565b606091505b5090925090508115155f036109a057805115610b845780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276d696e742720726576657274656400000060448201526064016105e0565b5f546001600160a01b03163314610bf55760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b0316637eee288d60e01b17905251610c519190611421565b5f60405180830381855af49150503d805f8114610c89576040519150601f19603f3d011682016040523d82523d5f602084013e610c8e565b606091505b5090925090508115155f036109a057805115610cad5780518082602001fd5b60405162461bcd60e51b815260206004820152601f60248201527f46756e6374696f6e2063616c6c2027756e6c6f636b272072657665727465640060448201526064016105e0565b600d546040516001600160a01b0383811660248301525f92839283929091169060440160408051601f198184030181529181526020820180516001600160e01b0316631092ab9160e31b17905251610d4d9190611421565b5f60405180830381855af49150503d805f8114610d85576040519150601f19603f3d011682016040523d82523d5f602084013e610d8a565b606091505b5090925090508115155f03610e0457805115610da95780518082602001fd5b60405162461bcd60e51b815260206004820152602a60248201527f46756e6374696f6e2063616c6c2027756e6c6f636b656442616c616e63654f6660448201526909c81c995d995c9d195960b21b60648201526084016105e0565b80806020019051810190610e1891906114a3565b949350505050565b600d546040516001600160a01b0383811660248301525f92839283929091169060440160408051601f198184030181529181526020820180516001600160e01b0316631293027760e31b17905251610e789190611421565b5f60405180830381855af49150503d805f8114610eb0576040519150601f19603f3d011682016040523d82523d5f602084013e610eb5565b606091505b5090925090508115155f03610e0457805115610ed45780518082602001fd5b60405162461bcd60e51b815260206004820152602960248201527f46756e6374696f6e2063616c6c2027756e636c61696d65645265776172647327604482015268081c995d995c9d195960ba1b60648201526084016105e0565b6060600b8054610460906113e9565b5f546001600160a01b03163314610f665760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b0316632770a7eb60e21b17905251610fc29190611421565b5f60405180830381855af49150503d805f8114610ffa576040519150601f19603f3d011682016040523d82523d5f602084013e610fff565b606091505b5090925090508115155f036109a05780511561101e5780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276275726e2720726576657274656400000060448201526064016105e0565b600d546040516001600160a01b038481166024830152604482018490525f92839283929091169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b179052516110c59190611421565b5f60405180830381855af49150503d805f81146110fd576040519150601f19603f3d011682016040523d82523d5f602084013e611102565b606091505b5090925090508115155f036105e9578051156111215780518082602001fd5b60405162461bcd60e51b815260206004820152602160248201527f46756e6374696f6e2063616c6c20277472616e736665722720726576657274656044820152601960fa1b60648201526084016105e0565b5f80546001600160a01b0316331461119d5760405162461bcd60e51b81526004016105e09061145b565b600d5460408051600481526024810182526020810180516001600160e01b031663fb489a7b60e01b17905290515f9283926001600160a01b03909116916111e49190611421565b5f60405180830381855af49150503d805f811461121c576040519150601f19603f3d011682016040523d82523d5f602084013e611221565b606091505b5090925090508115155f03611296578051156112405780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c20277265646973747269627574652720726576604482015264195c9d195960da1b60648201526084016105e0565b808060200190518101906112aa91906114a3565b9250505090565b5f5b838110156112cb5781810151838201526020016112b3565b50505f910152565b602081525f82518060208401526112f18160408501602087016112b1565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461131b575f80fd5b919050565b5f8060408385031215611331575f80fd5b61133a83611305565b946020939093013593505050565b5f60208284031215611358575f80fd5b5035919050565b5f805f60608486031215611371575f80fd5b61137a84611305565b925061138860208501611305565b9150604084013590509250925092565b5f602082840312156113a8575f80fd5b6113b182611305565b9392505050565b5f80604083850312156113c9575f80fd5b6113d283611305565b91506113e060208401611305565b90509250929050565b600181811c908216806113fd57607f821691505b60208210810361141b57634e487b7160e01b5f52602260045260245ffd5b50919050565b5f82516114328184602087016112b1565b9190910192915050565b5f6020828403121561144c575f80fd5b815180151581146113b1575f80fd5b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b5f602082840312156114b3575f80fd5b505191905056fea26469706673582212200f950bd4f4d9457214149c5d71464d2a2df907f4dd143e5feaf92aaf3d4852a864736f6c63430008170033",
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
func DeployAutonity(auth *bind.TransactOpts, backend bind.ContractBackend, _validators []AutonityValidator, _config AutonityConfig) (common.Address, *types.Transaction, *Autonity, error) {
	parsed, err := AutonityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AutonityBin), backend, _validators, _config)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Autonity{AutonityCaller: AutonityCaller{contract: contract}, AutonityTransactor: AutonityTransactor{contract: contract}, AutonityFilterer: AutonityFilterer{contract: contract}}, nil
}

// Autonity is an auto generated Go binding around an Ethereum contract.
type Autonity struct {
	AutonityCaller     // Read-only binding to the contract
	AutonityTransactor // Write-only binding to the contract
	AutonityFilterer   // Log filterer for contract events
}

// AutonityCaller is an auto generated read-only Go binding around an Ethereum contract.
type AutonityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AutonityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AutonityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AutonityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AutonityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AutonitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AutonitySession struct {
	Contract     *Autonity         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AutonityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AutonityCallerSession struct {
	Contract *AutonityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// AutonityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AutonityTransactorSession struct {
	Contract     *AutonityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AutonityRaw is an auto generated low-level Go binding around an Ethereum contract.
type AutonityRaw struct {
	Contract *Autonity // Generic contract binding to access the raw methods on
}

// AutonityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AutonityCallerRaw struct {
	Contract *AutonityCaller // Generic read-only contract binding to access the raw methods on
}

// AutonityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AutonityTransactorRaw struct {
	Contract *AutonityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAutonity creates a new instance of Autonity, bound to a specific deployed contract.
func NewAutonity(address common.Address, backend bind.ContractBackend) (*Autonity, error) {
	contract, err := bindAutonity(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Autonity{AutonityCaller: AutonityCaller{contract: contract}, AutonityTransactor: AutonityTransactor{contract: contract}, AutonityFilterer: AutonityFilterer{contract: contract}}, nil
}

// NewAutonityCaller creates a new read-only instance of Autonity, bound to a specific deployed contract.
func NewAutonityCaller(address common.Address, caller bind.ContractCaller) (*AutonityCaller, error) {
	contract, err := bindAutonity(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AutonityCaller{contract: contract}, nil
}

// NewAutonityTransactor creates a new write-only instance of Autonity, bound to a specific deployed contract.
func NewAutonityTransactor(address common.Address, transactor bind.ContractTransactor) (*AutonityTransactor, error) {
	contract, err := bindAutonity(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AutonityTransactor{contract: contract}, nil
}

// NewAutonityFilterer creates a new log filterer instance of Autonity, bound to a specific deployed contract.
func NewAutonityFilterer(address common.Address, filterer bind.ContractFilterer) (*AutonityFilterer, error) {
	contract, err := bindAutonity(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AutonityFilterer{contract: contract}, nil
}

// bindAutonity binds a generic wrapper to an already deployed contract.
func bindAutonity(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AutonityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Autonity *AutonityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Autonity.Contract.AutonityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Autonity *AutonityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.Contract.AutonityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Autonity *AutonityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Autonity.Contract.AutonityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Autonity *AutonityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Autonity.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Autonity *AutonityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Autonity *AutonityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Autonity.Contract.contract.Transact(opts, method, params...)
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Autonity *AutonityCaller) COMMISSIONRATEPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "COMMISSION_RATE_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Autonity *AutonitySession) COMMISSIONRATEPRECISION() (*big.Int, error) {
	return _Autonity.Contract.COMMISSIONRATEPRECISION(&_Autonity.CallOpts)
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Autonity *AutonityCallerSession) COMMISSIONRATEPRECISION() (*big.Int, error) {
	return _Autonity.Contract.COMMISSIONRATEPRECISION(&_Autonity.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Autonity *AutonityCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Autonity *AutonitySession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Autonity.Contract.Allowance(&_Autonity.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Autonity *AutonityCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Autonity.Contract.Allowance(&_Autonity.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_Autonity *AutonityCaller) BalanceOf(opts *bind.CallOpts, _addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "balanceOf", _addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_Autonity *AutonitySession) BalanceOf(_addr common.Address) (*big.Int, error) {
	return _Autonity.Contract.BalanceOf(&_Autonity.CallOpts, _addr)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_Autonity *AutonityCallerSession) BalanceOf(_addr common.Address) (*big.Int, error) {
	return _Autonity.Contract.BalanceOf(&_Autonity.CallOpts, _addr)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,address) policy, (address,address,address,address,address) contracts, (address,uint256,uint256,uint256) protocol, uint256 contractVersion)
func (_Autonity *AutonityCaller) Config(opts *bind.CallOpts) (struct {
	Policy          AutonityPolicy
	Contracts       AutonityContracts
	Protocol        AutonityProtocol
	ContractVersion *big.Int
}, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "config")

	outstruct := new(struct {
		Policy          AutonityPolicy
		Contracts       AutonityContracts
		Protocol        AutonityProtocol
		ContractVersion *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Policy = *abi.ConvertType(out[0], new(AutonityPolicy)).(*AutonityPolicy)
	outstruct.Contracts = *abi.ConvertType(out[1], new(AutonityContracts)).(*AutonityContracts)
	outstruct.Protocol = *abi.ConvertType(out[2], new(AutonityProtocol)).(*AutonityProtocol)
	outstruct.ContractVersion = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,address) policy, (address,address,address,address,address) contracts, (address,uint256,uint256,uint256) protocol, uint256 contractVersion)
func (_Autonity *AutonitySession) Config() (struct {
	Policy          AutonityPolicy
	Contracts       AutonityContracts
	Protocol        AutonityProtocol
	ContractVersion *big.Int
}, error) {
	return _Autonity.Contract.Config(&_Autonity.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,address) policy, (address,address,address,address,address) contracts, (address,uint256,uint256,uint256) protocol, uint256 contractVersion)
func (_Autonity *AutonityCallerSession) Config() (struct {
	Policy          AutonityPolicy
	Contracts       AutonityContracts
	Protocol        AutonityProtocol
	ContractVersion *big.Int
}, error) {
	return _Autonity.Contract.Config(&_Autonity.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Autonity *AutonityCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Autonity *AutonitySession) Decimals() (uint8, error) {
	return _Autonity.Contract.Decimals(&_Autonity.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Autonity *AutonityCallerSession) Decimals() (uint8, error) {
	return _Autonity.Contract.Decimals(&_Autonity.CallOpts)
}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() view returns(address)
func (_Autonity *AutonityCaller) Deployer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "deployer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() view returns(address)
func (_Autonity *AutonitySession) Deployer() (common.Address, error) {
	return _Autonity.Contract.Deployer(&_Autonity.CallOpts)
}

// Deployer is a free data retrieval call binding the contract method 0xd5f39488.
//
// Solidity: function deployer() view returns(address)
func (_Autonity *AutonityCallerSession) Deployer() (common.Address, error) {
	return _Autonity.Contract.Deployer(&_Autonity.CallOpts)
}

// EpochID is a free data retrieval call binding the contract method 0xc9d97af4.
//
// Solidity: function epochID() view returns(uint256)
func (_Autonity *AutonityCaller) EpochID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "epochID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochID is a free data retrieval call binding the contract method 0xc9d97af4.
//
// Solidity: function epochID() view returns(uint256)
func (_Autonity *AutonitySession) EpochID() (*big.Int, error) {
	return _Autonity.Contract.EpochID(&_Autonity.CallOpts)
}

// EpochID is a free data retrieval call binding the contract method 0xc9d97af4.
//
// Solidity: function epochID() view returns(uint256)
func (_Autonity *AutonityCallerSession) EpochID() (*big.Int, error) {
	return _Autonity.Contract.EpochID(&_Autonity.CallOpts)
}

// EpochReward is a free data retrieval call binding the contract method 0x1604e416.
//
// Solidity: function epochReward() view returns(uint256)
func (_Autonity *AutonityCaller) EpochReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "epochReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochReward is a free data retrieval call binding the contract method 0x1604e416.
//
// Solidity: function epochReward() view returns(uint256)
func (_Autonity *AutonitySession) EpochReward() (*big.Int, error) {
	return _Autonity.Contract.EpochReward(&_Autonity.CallOpts)
}

// EpochReward is a free data retrieval call binding the contract method 0x1604e416.
//
// Solidity: function epochReward() view returns(uint256)
func (_Autonity *AutonityCallerSession) EpochReward() (*big.Int, error) {
	return _Autonity.Contract.EpochReward(&_Autonity.CallOpts)
}

// EpochTotalBondedStake is a free data retrieval call binding the contract method 0x9c98e471.
//
// Solidity: function epochTotalBondedStake() view returns(uint256)
func (_Autonity *AutonityCaller) EpochTotalBondedStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "epochTotalBondedStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochTotalBondedStake is a free data retrieval call binding the contract method 0x9c98e471.
//
// Solidity: function epochTotalBondedStake() view returns(uint256)
func (_Autonity *AutonitySession) EpochTotalBondedStake() (*big.Int, error) {
	return _Autonity.Contract.EpochTotalBondedStake(&_Autonity.CallOpts)
}

// EpochTotalBondedStake is a free data retrieval call binding the contract method 0x9c98e471.
//
// Solidity: function epochTotalBondedStake() view returns(uint256)
func (_Autonity *AutonityCallerSession) EpochTotalBondedStake() (*big.Int, error) {
	return _Autonity.Contract.EpochTotalBondedStake(&_Autonity.CallOpts)
}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetBlockPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getBlockPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_Autonity *AutonitySession) GetBlockPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetBlockPeriod(&_Autonity.CallOpts)
}

// GetBlockPeriod is a free data retrieval call binding the contract method 0x43645969.
//
// Solidity: function getBlockPeriod() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetBlockPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetBlockPeriod(&_Autonity.CallOpts)
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256)[])
func (_Autonity *AutonityCaller) GetCommittee(opts *bind.CallOpts) ([]AutonityCommitteeMember, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getCommittee")

	if err != nil {
		return *new([]AutonityCommitteeMember), err
	}

	out0 := *abi.ConvertType(out[0], new([]AutonityCommitteeMember)).(*[]AutonityCommitteeMember)

	return out0, err

}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256)[])
func (_Autonity *AutonitySession) GetCommittee() ([]AutonityCommitteeMember, error) {
	return _Autonity.Contract.GetCommittee(&_Autonity.CallOpts)
}

// GetCommittee is a free data retrieval call binding the contract method 0xab8f6ffe.
//
// Solidity: function getCommittee() view returns((address,uint256)[])
func (_Autonity *AutonityCallerSession) GetCommittee() ([]AutonityCommitteeMember, error) {
	return _Autonity.Contract.GetCommittee(&_Autonity.CallOpts)
}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_Autonity *AutonityCaller) GetCommitteeEnodes(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getCommitteeEnodes")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_Autonity *AutonitySession) GetCommitteeEnodes() ([]string, error) {
	return _Autonity.Contract.GetCommitteeEnodes(&_Autonity.CallOpts)
}

// GetCommitteeEnodes is a free data retrieval call binding the contract method 0xa8b2216e.
//
// Solidity: function getCommitteeEnodes() view returns(string[])
func (_Autonity *AutonityCallerSession) GetCommitteeEnodes() ([]string, error) {
	return _Autonity.Contract.GetCommitteeEnodes(&_Autonity.CallOpts)
}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_Autonity *AutonityCaller) GetEpochFromBlock(opts *bind.CallOpts, _block *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getEpochFromBlock", _block)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_Autonity *AutonitySession) GetEpochFromBlock(_block *big.Int) (*big.Int, error) {
	return _Autonity.Contract.GetEpochFromBlock(&_Autonity.CallOpts, _block)
}

// GetEpochFromBlock is a free data retrieval call binding the contract method 0x96b477cb.
//
// Solidity: function getEpochFromBlock(uint256 _block) view returns(uint256)
func (_Autonity *AutonityCallerSession) GetEpochFromBlock(_block *big.Int) (*big.Int, error) {
	return _Autonity.Contract.GetEpochFromBlock(&_Autonity.CallOpts, _block)
}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetEpochPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getEpochPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_Autonity *AutonitySession) GetEpochPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetEpochPeriod(&_Autonity.CallOpts)
}

// GetEpochPeriod is a free data retrieval call binding the contract method 0xdfb1a4d2.
//
// Solidity: function getEpochPeriod() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetEpochPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetEpochPeriod(&_Autonity.CallOpts)
}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCaller) GetLastEpochBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getLastEpochBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_Autonity *AutonitySession) GetLastEpochBlock() (*big.Int, error) {
	return _Autonity.Contract.GetLastEpochBlock(&_Autonity.CallOpts)
}

// GetLastEpochBlock is a free data retrieval call binding the contract method 0x731b3a03.
//
// Solidity: function getLastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetLastEpochBlock() (*big.Int, error) {
	return _Autonity.Contract.GetLastEpochBlock(&_Autonity.CallOpts)
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_Autonity *AutonityCaller) GetMaxCommitteeSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getMaxCommitteeSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_Autonity *AutonitySession) GetMaxCommitteeSize() (*big.Int, error) {
	return _Autonity.Contract.GetMaxCommitteeSize(&_Autonity.CallOpts)
}

// GetMaxCommitteeSize is a free data retrieval call binding the contract method 0x819b6463.
//
// Solidity: function getMaxCommitteeSize() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetMaxCommitteeSize() (*big.Int, error) {
	return _Autonity.Contract.GetMaxCommitteeSize(&_Autonity.CallOpts)
}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_Autonity *AutonityCaller) GetMinimumBaseFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getMinimumBaseFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_Autonity *AutonitySession) GetMinimumBaseFee() (*big.Int, error) {
	return _Autonity.Contract.GetMinimumBaseFee(&_Autonity.CallOpts)
}

// GetMinimumBaseFee is a free data retrieval call binding the contract method 0x11220633.
//
// Solidity: function getMinimumBaseFee() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetMinimumBaseFee() (*big.Int, error) {
	return _Autonity.Contract.GetMinimumBaseFee(&_Autonity.CallOpts)
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Autonity *AutonityCaller) GetNewContract(opts *bind.CallOpts) ([]byte, string, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getNewContract")

	if err != nil {
		return *new([]byte), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Autonity *AutonitySession) GetNewContract() ([]byte, string, error) {
	return _Autonity.Contract.GetNewContract(&_Autonity.CallOpts)
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Autonity *AutonityCallerSession) GetNewContract() ([]byte, string, error) {
	return _Autonity.Contract.GetNewContract(&_Autonity.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_Autonity *AutonityCaller) GetOperator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getOperator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_Autonity *AutonitySession) GetOperator() (common.Address, error) {
	return _Autonity.Contract.GetOperator(&_Autonity.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_Autonity *AutonityCallerSession) GetOperator() (common.Address, error) {
	return _Autonity.Contract.GetOperator(&_Autonity.CallOpts)
}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_Autonity *AutonityCaller) GetOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_Autonity *AutonitySession) GetOracle() (common.Address, error) {
	return _Autonity.Contract.GetOracle(&_Autonity.CallOpts)
}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_Autonity *AutonityCallerSession) GetOracle() (common.Address, error) {
	return _Autonity.Contract.GetOracle(&_Autonity.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0x5f7d3949.
//
// Solidity: function getProposer(uint256 height, uint256 round) view returns(address)
func (_Autonity *AutonityCaller) GetProposer(opts *bind.CallOpts, height *big.Int, round *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getProposer", height, round)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProposer is a free data retrieval call binding the contract method 0x5f7d3949.
//
// Solidity: function getProposer(uint256 height, uint256 round) view returns(address)
func (_Autonity *AutonitySession) GetProposer(height *big.Int, round *big.Int) (common.Address, error) {
	return _Autonity.Contract.GetProposer(&_Autonity.CallOpts, height, round)
}

// GetProposer is a free data retrieval call binding the contract method 0x5f7d3949.
//
// Solidity: function getProposer(uint256 height, uint256 round) view returns(address)
func (_Autonity *AutonityCallerSession) GetProposer(height *big.Int, round *big.Int) (common.Address, error) {
	return _Autonity.Contract.GetProposer(&_Autonity.CallOpts, height, round)
}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_Autonity *AutonityCaller) GetTreasuryAccount(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getTreasuryAccount")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_Autonity *AutonitySession) GetTreasuryAccount() (common.Address, error) {
	return _Autonity.Contract.GetTreasuryAccount(&_Autonity.CallOpts)
}

// GetTreasuryAccount is a free data retrieval call binding the contract method 0xf7866ee3.
//
// Solidity: function getTreasuryAccount() view returns(address)
func (_Autonity *AutonityCallerSession) GetTreasuryAccount() (common.Address, error) {
	return _Autonity.Contract.GetTreasuryAccount(&_Autonity.CallOpts)
}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_Autonity *AutonityCaller) GetTreasuryFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getTreasuryFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_Autonity *AutonitySession) GetTreasuryFee() (*big.Int, error) {
	return _Autonity.Contract.GetTreasuryFee(&_Autonity.CallOpts)
}

// GetTreasuryFee is a free data retrieval call binding the contract method 0x29070c6d.
//
// Solidity: function getTreasuryFee() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetTreasuryFee() (*big.Int, error) {
	return _Autonity.Contract.GetTreasuryFee(&_Autonity.CallOpts)
}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_Autonity *AutonityCaller) GetUnbondingPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getUnbondingPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_Autonity *AutonitySession) GetUnbondingPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetUnbondingPeriod(&_Autonity.CallOpts)
}

// GetUnbondingPeriod is a free data retrieval call binding the contract method 0x6fd2c80b.
//
// Solidity: function getUnbondingPeriod() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetUnbondingPeriod() (*big.Int, error) {
	return _Autonity.Contract.GetUnbondingPeriod(&_Autonity.CallOpts)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Autonity *AutonityCaller) GetValidator(opts *bind.CallOpts, _addr common.Address) (AutonityValidator, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getValidator", _addr)

	if err != nil {
		return *new(AutonityValidator), err
	}

	out0 := *abi.ConvertType(out[0], new(AutonityValidator)).(*AutonityValidator)

	return out0, err

}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Autonity *AutonitySession) GetValidator(_addr common.Address) (AutonityValidator, error) {
	return _Autonity.Contract.GetValidator(&_Autonity.CallOpts, _addr)
}

// GetValidator is a free data retrieval call binding the contract method 0x1904bb2e.
//
// Solidity: function getValidator(address _addr) view returns((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,uint8))
func (_Autonity *AutonityCallerSession) GetValidator(_addr common.Address) (AutonityValidator, error) {
	return _Autonity.Contract.GetValidator(&_Autonity.CallOpts, _addr)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Autonity *AutonityCaller) GetValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getValidators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Autonity *AutonitySession) GetValidators() ([]common.Address, error) {
	return _Autonity.Contract.GetValidators(&_Autonity.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Autonity *AutonityCallerSession) GetValidators() ([]common.Address, error) {
	return _Autonity.Contract.GetValidators(&_Autonity.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_Autonity *AutonityCaller) GetVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "getVersion")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_Autonity *AutonitySession) GetVersion() (*big.Int, error) {
	return _Autonity.Contract.GetVersion(&_Autonity.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(uint256)
func (_Autonity *AutonityCallerSession) GetVersion() (*big.Int, error) {
	return _Autonity.Contract.GetVersion(&_Autonity.CallOpts)
}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCaller) LastEpochBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "lastEpochBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_Autonity *AutonitySession) LastEpochBlock() (*big.Int, error) {
	return _Autonity.Contract.LastEpochBlock(&_Autonity.CallOpts)
}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_Autonity *AutonityCallerSession) LastEpochBlock() (*big.Int, error) {
	return _Autonity.Contract.LastEpochBlock(&_Autonity.CallOpts)
}

// LiquidLogicAddress is a free data retrieval call binding the contract method 0x66d3bd27.
//
// Solidity: function liquidLogicAddress() view returns(address)
func (_Autonity *AutonityCaller) LiquidLogicAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "liquidLogicAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LiquidLogicAddress is a free data retrieval call binding the contract method 0x66d3bd27.
//
// Solidity: function liquidLogicAddress() view returns(address)
func (_Autonity *AutonitySession) LiquidLogicAddress() (common.Address, error) {
	return _Autonity.Contract.LiquidLogicAddress(&_Autonity.CallOpts)
}

// LiquidLogicAddress is a free data retrieval call binding the contract method 0x66d3bd27.
//
// Solidity: function liquidLogicAddress() view returns(address)
func (_Autonity *AutonityCallerSession) LiquidLogicAddress() (common.Address, error) {
	return _Autonity.Contract.LiquidLogicAddress(&_Autonity.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Autonity *AutonityCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Autonity *AutonitySession) Name() (string, error) {
	return _Autonity.Contract.Name(&_Autonity.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_Autonity *AutonityCallerSession) Name() (string, error) {
	return _Autonity.Contract.Name(&_Autonity.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Autonity *AutonityCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Autonity *AutonitySession) Symbol() (string, error) {
	return _Autonity.Contract.Symbol(&_Autonity.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_Autonity *AutonityCallerSession) Symbol() (string, error) {
	return _Autonity.Contract.Symbol(&_Autonity.CallOpts)
}

// TotalRedistributed is a free data retrieval call binding the contract method 0x9bb851c0.
//
// Solidity: function totalRedistributed() view returns(uint256)
func (_Autonity *AutonityCaller) TotalRedistributed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "totalRedistributed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalRedistributed is a free data retrieval call binding the contract method 0x9bb851c0.
//
// Solidity: function totalRedistributed() view returns(uint256)
func (_Autonity *AutonitySession) TotalRedistributed() (*big.Int, error) {
	return _Autonity.Contract.TotalRedistributed(&_Autonity.CallOpts)
}

// TotalRedistributed is a free data retrieval call binding the contract method 0x9bb851c0.
//
// Solidity: function totalRedistributed() view returns(uint256)
func (_Autonity *AutonityCallerSession) TotalRedistributed() (*big.Int, error) {
	return _Autonity.Contract.TotalRedistributed(&_Autonity.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Autonity *AutonityCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Autonity.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Autonity *AutonitySession) TotalSupply() (*big.Int, error) {
	return _Autonity.Contract.TotalSupply(&_Autonity.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Autonity *AutonityCallerSession) TotalSupply() (*big.Int, error) {
	return _Autonity.Contract.TotalSupply(&_Autonity.CallOpts)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_Autonity *AutonityTransactor) ActivateValidator(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "activateValidator", _address)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_Autonity *AutonitySession) ActivateValidator(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.ActivateValidator(&_Autonity.TransactOpts, _address)
}

// ActivateValidator is a paid mutator transaction binding the contract method 0xb46e5520.
//
// Solidity: function activateValidator(address _address) returns()
func (_Autonity *AutonityTransactorSession) ActivateValidator(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.ActivateValidator(&_Autonity.TransactOpts, _address)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Autonity *AutonitySession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Approve(&_Autonity.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Approve(&_Autonity.TransactOpts, spender, amount)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Bond(opts *bind.TransactOpts, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "bond", _validator, _amount)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonitySession) Bond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Bond(&_Autonity.TransactOpts, _validator, _amount)
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) Bond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Bond(&_Autonity.TransactOpts, _validator, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Burn(opts *bind.TransactOpts, _addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "burn", _addr, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_Autonity *AutonitySession) Burn(_addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Burn(&_Autonity.TransactOpts, _addr, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) Burn(_addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Burn(&_Autonity.TransactOpts, _addr, _amount)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_Autonity *AutonityTransactor) ChangeCommissionRate(opts *bind.TransactOpts, _validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "changeCommissionRate", _validator, _rate)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_Autonity *AutonitySession) ChangeCommissionRate(_validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.ChangeCommissionRate(&_Autonity.TransactOpts, _validator, _rate)
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_Autonity *AutonityTransactorSession) ChangeCommissionRate(_validator common.Address, _rate *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.ChangeCommissionRate(&_Autonity.TransactOpts, _validator, _rate)
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Autonity *AutonityTransactor) CompleteContractUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "completeContractUpgrade")
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Autonity *AutonitySession) CompleteContractUpgrade() (*types.Transaction, error) {
	return _Autonity.Contract.CompleteContractUpgrade(&_Autonity.TransactOpts)
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Autonity *AutonityTransactorSession) CompleteContractUpgrade() (*types.Transaction, error) {
	return _Autonity.Contract.CompleteContractUpgrade(&_Autonity.TransactOpts)
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[])
func (_Autonity *AutonityTransactor) ComputeCommittee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "computeCommittee")
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[])
func (_Autonity *AutonitySession) ComputeCommittee() (*types.Transaction, error) {
	return _Autonity.Contract.ComputeCommittee(&_Autonity.TransactOpts)
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[])
func (_Autonity *AutonityTransactorSession) ComputeCommittee() (*types.Transaction, error) {
	return _Autonity.Contract.ComputeCommittee(&_Autonity.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256)[])
func (_Autonity *AutonityTransactor) Finalize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "finalize")
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256)[])
func (_Autonity *AutonitySession) Finalize() (*types.Transaction, error) {
	return _Autonity.Contract.Finalize(&_Autonity.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256)[])
func (_Autonity *AutonityTransactorSession) Finalize() (*types.Transaction, error) {
	return _Autonity.Contract.Finalize(&_Autonity.TransactOpts)
}

// FinalizeInitialization is a paid mutator transaction binding the contract method 0xd861b0e8.
//
// Solidity: function finalizeInitialization() returns()
func (_Autonity *AutonityTransactor) FinalizeInitialization(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "finalizeInitialization")
}

// FinalizeInitialization is a paid mutator transaction binding the contract method 0xd861b0e8.
//
// Solidity: function finalizeInitialization() returns()
func (_Autonity *AutonitySession) FinalizeInitialization() (*types.Transaction, error) {
	return _Autonity.Contract.FinalizeInitialization(&_Autonity.TransactOpts)
}

// FinalizeInitialization is a paid mutator transaction binding the contract method 0xd861b0e8.
//
// Solidity: function finalizeInitialization() returns()
func (_Autonity *AutonityTransactorSession) FinalizeInitialization() (*types.Transaction, error) {
	return _Autonity.Contract.FinalizeInitialization(&_Autonity.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Mint(opts *bind.TransactOpts, _addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "mint", _addr, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_Autonity *AutonitySession) Mint(_addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Mint(&_Autonity.TransactOpts, _addr, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) Mint(_addr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Mint(&_Autonity.TransactOpts, _addr, _amount)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_Autonity *AutonityTransactor) PauseValidator(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "pauseValidator", _address)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_Autonity *AutonitySession) PauseValidator(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.PauseValidator(&_Autonity.TransactOpts, _address)
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_Autonity *AutonityTransactorSession) PauseValidator(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.PauseValidator(&_Autonity.TransactOpts, _address)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xad722d4d.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _multisig) returns()
func (_Autonity *AutonityTransactor) RegisterValidator(opts *bind.TransactOpts, _enode string, _oracleAddress common.Address, _multisig []byte) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "registerValidator", _enode, _oracleAddress, _multisig)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xad722d4d.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _multisig) returns()
func (_Autonity *AutonitySession) RegisterValidator(_enode string, _oracleAddress common.Address, _multisig []byte) (*types.Transaction, error) {
	return _Autonity.Contract.RegisterValidator(&_Autonity.TransactOpts, _enode, _oracleAddress, _multisig)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xad722d4d.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _multisig) returns()
func (_Autonity *AutonityTransactorSession) RegisterValidator(_enode string, _oracleAddress common.Address, _multisig []byte) (*types.Transaction, error) {
	return _Autonity.Contract.RegisterValidator(&_Autonity.TransactOpts, _enode, _oracleAddress, _multisig)
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Autonity *AutonityTransactor) ResetContractUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "resetContractUpgrade")
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Autonity *AutonitySession) ResetContractUpgrade() (*types.Transaction, error) {
	return _Autonity.Contract.ResetContractUpgrade(&_Autonity.TransactOpts)
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Autonity *AutonityTransactorSession) ResetContractUpgrade() (*types.Transaction, error) {
	return _Autonity.Contract.ResetContractUpgrade(&_Autonity.TransactOpts)
}

// SetAccountabilityContract is a paid mutator transaction binding the contract method 0x1250a28d.
//
// Solidity: function setAccountabilityContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetAccountabilityContract(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setAccountabilityContract", _address)
}

// SetAccountabilityContract is a paid mutator transaction binding the contract method 0x1250a28d.
//
// Solidity: function setAccountabilityContract(address _address) returns()
func (_Autonity *AutonitySession) SetAccountabilityContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetAccountabilityContract(&_Autonity.TransactOpts, _address)
}

// SetAccountabilityContract is a paid mutator transaction binding the contract method 0x1250a28d.
//
// Solidity: function setAccountabilityContract(address _address) returns()
func (_Autonity *AutonityTransactorSession) SetAccountabilityContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetAccountabilityContract(&_Autonity.TransactOpts, _address)
}

// SetAcuContract is a paid mutator transaction binding the contract method 0xd372c07e.
//
// Solidity: function setAcuContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetAcuContract(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setAcuContract", _address)
}

// SetAcuContract is a paid mutator transaction binding the contract method 0xd372c07e.
//
// Solidity: function setAcuContract(address _address) returns()
func (_Autonity *AutonitySession) SetAcuContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetAcuContract(&_Autonity.TransactOpts, _address)
}

// SetAcuContract is a paid mutator transaction binding the contract method 0xd372c07e.
//
// Solidity: function setAcuContract(address _address) returns()
func (_Autonity *AutonityTransactorSession) SetAcuContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetAcuContract(&_Autonity.TransactOpts, _address)
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_Autonity *AutonityTransactor) SetCommitteeSize(opts *bind.TransactOpts, _size *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setCommitteeSize", _size)
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_Autonity *AutonitySession) SetCommitteeSize(_size *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetCommitteeSize(&_Autonity.TransactOpts, _size)
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_Autonity *AutonityTransactorSession) SetCommitteeSize(_size *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetCommitteeSize(&_Autonity.TransactOpts, _size)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactor) SetEpochPeriod(opts *bind.TransactOpts, _period *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setEpochPeriod", _period)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_Autonity *AutonitySession) SetEpochPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetEpochPeriod(&_Autonity.TransactOpts, _period)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactorSession) SetEpochPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetEpochPeriod(&_Autonity.TransactOpts, _period)
}

// SetLiquidNewtonLogicContract is a paid mutator transaction binding the contract method 0x07b5c33b.
//
// Solidity: function setLiquidNewtonLogicContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetLiquidNewtonLogicContract(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setLiquidNewtonLogicContract", _address)
}

// SetLiquidNewtonLogicContract is a paid mutator transaction binding the contract method 0x07b5c33b.
//
// Solidity: function setLiquidNewtonLogicContract(address _address) returns()
func (_Autonity *AutonitySession) SetLiquidNewtonLogicContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetLiquidNewtonLogicContract(&_Autonity.TransactOpts, _address)
}

// SetLiquidNewtonLogicContract is a paid mutator transaction binding the contract method 0x07b5c33b.
//
// Solidity: function setLiquidNewtonLogicContract(address _address) returns()
func (_Autonity *AutonityTransactorSession) SetLiquidNewtonLogicContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetLiquidNewtonLogicContract(&_Autonity.TransactOpts, _address)
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_Autonity *AutonityTransactor) SetMinimumBaseFee(opts *bind.TransactOpts, _price *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setMinimumBaseFee", _price)
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_Autonity *AutonitySession) SetMinimumBaseFee(_price *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetMinimumBaseFee(&_Autonity.TransactOpts, _price)
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_Autonity *AutonityTransactorSession) SetMinimumBaseFee(_price *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetMinimumBaseFee(&_Autonity.TransactOpts, _price)
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_Autonity *AutonityTransactor) SetOperatorAccount(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setOperatorAccount", _account)
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_Autonity *AutonitySession) SetOperatorAccount(_account common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetOperatorAccount(&_Autonity.TransactOpts, _account)
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_Autonity *AutonityTransactorSession) SetOperatorAccount(_account common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetOperatorAccount(&_Autonity.TransactOpts, _account)
}

// SetOracleContract is a paid mutator transaction binding the contract method 0x496ccd9b.
//
// Solidity: function setOracleContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetOracleContract(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setOracleContract", _address)
}

// SetOracleContract is a paid mutator transaction binding the contract method 0x496ccd9b.
//
// Solidity: function setOracleContract(address _address) returns()
func (_Autonity *AutonitySession) SetOracleContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetOracleContract(&_Autonity.TransactOpts, _address)
}

// SetOracleContract is a paid mutator transaction binding the contract method 0x496ccd9b.
//
// Solidity: function setOracleContract(address _address) returns()
func (_Autonity *AutonityTransactorSession) SetOracleContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetOracleContract(&_Autonity.TransactOpts, _address)
}

// SetStabilizationContract is a paid mutator transaction binding the contract method 0xcfd19fb9.
//
// Solidity: function setStabilizationContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetStabilizationContract(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setStabilizationContract", _address)
}

// SetStabilizationContract is a paid mutator transaction binding the contract method 0xcfd19fb9.
//
// Solidity: function setStabilizationContract(address _address) returns()
func (_Autonity *AutonitySession) SetStabilizationContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetStabilizationContract(&_Autonity.TransactOpts, _address)
}

// SetStabilizationContract is a paid mutator transaction binding the contract method 0xcfd19fb9.
//
// Solidity: function setStabilizationContract(address _address) returns()
func (_Autonity *AutonityTransactorSession) SetStabilizationContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetStabilizationContract(&_Autonity.TransactOpts, _address)
}

// SetSupplyControlContract is a paid mutator transaction binding the contract method 0xb3ecbadd.
//
// Solidity: function setSupplyControlContract(address _address) returns()
func (_Autonity *AutonityTransactor) SetSupplyControlContract(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setSupplyControlContract", _address)
}

// SetSupplyControlContract is a paid mutator transaction binding the contract method 0xb3ecbadd.
//
// Solidity: function setSupplyControlContract(address _address) returns()
func (_Autonity *AutonitySession) SetSupplyControlContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetSupplyControlContract(&_Autonity.TransactOpts, _address)
}

// SetSupplyControlContract is a paid mutator transaction binding the contract method 0xb3ecbadd.
//
// Solidity: function setSupplyControlContract(address _address) returns()
func (_Autonity *AutonityTransactorSession) SetSupplyControlContract(_address common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetSupplyControlContract(&_Autonity.TransactOpts, _address)
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_Autonity *AutonityTransactor) SetTreasuryAccount(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setTreasuryAccount", _account)
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_Autonity *AutonitySession) SetTreasuryAccount(_account common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetTreasuryAccount(&_Autonity.TransactOpts, _account)
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_Autonity *AutonityTransactorSession) SetTreasuryAccount(_account common.Address) (*types.Transaction, error) {
	return _Autonity.Contract.SetTreasuryAccount(&_Autonity.TransactOpts, _account)
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_Autonity *AutonityTransactor) SetTreasuryFee(opts *bind.TransactOpts, _treasuryFee *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setTreasuryFee", _treasuryFee)
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_Autonity *AutonitySession) SetTreasuryFee(_treasuryFee *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetTreasuryFee(&_Autonity.TransactOpts, _treasuryFee)
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_Autonity *AutonityTransactorSession) SetTreasuryFee(_treasuryFee *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetTreasuryFee(&_Autonity.TransactOpts, _treasuryFee)
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactor) SetUnbondingPeriod(opts *bind.TransactOpts, _period *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "setUnbondingPeriod", _period)
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_Autonity *AutonitySession) SetUnbondingPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetUnbondingPeriod(&_Autonity.TransactOpts, _period)
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_Autonity *AutonityTransactorSession) SetUnbondingPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.SetUnbondingPeriod(&_Autonity.TransactOpts, _period)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonityTransactor) Transfer(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "transfer", _recipient, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonitySession) Transfer(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Transfer(&_Autonity.TransactOpts, _recipient, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *AutonityTransactorSession) Transfer(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Transfer(&_Autonity.TransactOpts, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Autonity *AutonitySession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.TransferFrom(&_Autonity.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Autonity *AutonityTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.TransferFrom(&_Autonity.TransactOpts, sender, recipient, amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactor) Unbond(opts *bind.TransactOpts, _validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "unbond", _validator, _amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonitySession) Unbond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Unbond(&_Autonity.TransactOpts, _validator, _amount)
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_Autonity *AutonityTransactorSession) Unbond(_validator common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Autonity.Contract.Unbond(&_Autonity.TransactOpts, _validator, _amount)
}

// UpdateValidatorAndTransferSlashedFunds is a paid mutator transaction binding the contract method 0xf5371f4b.
//
// Solidity: function updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,uint8) _val) returns()
func (_Autonity *AutonityTransactor) UpdateValidatorAndTransferSlashedFunds(opts *bind.TransactOpts, _val AutonityValidator) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "updateValidatorAndTransferSlashedFunds", _val)
}

// UpdateValidatorAndTransferSlashedFunds is a paid mutator transaction binding the contract method 0xf5371f4b.
//
// Solidity: function updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,uint8) _val) returns()
func (_Autonity *AutonitySession) UpdateValidatorAndTransferSlashedFunds(_val AutonityValidator) (*types.Transaction, error) {
	return _Autonity.Contract.UpdateValidatorAndTransferSlashedFunds(&_Autonity.TransactOpts, _val)
}

// UpdateValidatorAndTransferSlashedFunds is a paid mutator transaction binding the contract method 0xf5371f4b.
//
// Solidity: function updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,uint8) _val) returns()
func (_Autonity *AutonityTransactorSession) UpdateValidatorAndTransferSlashedFunds(_val AutonityValidator) (*types.Transaction, error) {
	return _Autonity.Contract.UpdateValidatorAndTransferSlashedFunds(&_Autonity.TransactOpts, _val)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Autonity *AutonityTransactor) UpgradeContract(opts *bind.TransactOpts, _bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Autonity.contract.Transact(opts, "upgradeContract", _bytecode, _abi)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Autonity *AutonitySession) UpgradeContract(_bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Autonity.Contract.UpgradeContract(&_Autonity.TransactOpts, _bytecode, _abi)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Autonity *AutonityTransactorSession) UpgradeContract(_bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Autonity.Contract.UpgradeContract(&_Autonity.TransactOpts, _bytecode, _abi)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Autonity *AutonityTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Autonity.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Autonity *AutonitySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Autonity.Contract.Fallback(&_Autonity.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Autonity *AutonityTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Autonity.Contract.Fallback(&_Autonity.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Autonity *AutonityTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Autonity.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Autonity *AutonitySession) Receive() (*types.Transaction, error) {
	return _Autonity.Contract.Receive(&_Autonity.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Autonity *AutonityTransactorSession) Receive() (*types.Transaction, error) {
	return _Autonity.Contract.Receive(&_Autonity.TransactOpts)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Treasury       common.Address
	Addr           common.Address
	EffectiveBlock *big.Int
	Raw            types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Delegator common.Address
	Delegatee common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBondingRejected is a free log retrieval operation binding the contract event 0x14d528d8a5c5378dc43eafa0ea18cefb7d35f17cad6e4ca6062a4f8a71ea3dd9.
//
// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount)
func (_Autonity *AutonityFilterer) FilterBondingRejected(opts *bind.FilterOpts) (*AutonityBondingRejectedIterator, error) {

	logs, sub, err := _Autonity.contract.FilterLogs(opts, "BondingRejected")
	if err != nil {
		return nil, err
	}
	return &AutonityBondingRejectedIterator{contract: _Autonity.contract, event: "BondingRejected", logs: logs, sub: sub}, nil
}

// WatchBondingRejected is a free log subscription operation binding the contract event 0x14d528d8a5c5378dc43eafa0ea18cefb7d35f17cad6e4ca6062a4f8a71ea3dd9.
//
// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount)
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

// ParseBondingRejected is a log parse operation binding the contract event 0x14d528d8a5c5378dc43eafa0ea18cefb7d35f17cad6e4ca6062a4f8a71ea3dd9.
//
// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount)
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Validator common.Address
	Rate      *big.Int
	Raw       types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Period *big.Int
	Raw    types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	GasPrice *big.Int
	Raw      types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Validator  common.Address
	Delegator  common.Address
	SelfBonded bool
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Epoch *big.Int
	Raw   types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Validator  common.Address
	Delegator  common.Address
	SelfBonded bool
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Treasury       common.Address
	Addr           common.Address
	EffectiveBlock *big.Int
	Raw            types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Treasury       common.Address
	Addr           common.Address
	OracleAddress  common.Address
	Enode          string
	LiquidContract common.Address
	Raw            types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Addr   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
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

// BytesLibMetaData contains all meta data concerning the BytesLib contract.
var BytesLibMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220ed6d39ae665eec6e6f925cfad4aa9b561fbc3eaf5123811c801a25734ad9047164736f6c63430008170033",
}

// BytesLibABI is the input ABI used to generate the binding from.
// Deprecated: Use BytesLibMetaData.ABI instead.
var BytesLibABI = BytesLibMetaData.ABI

// BytesLibBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BytesLibMetaData.Bin instead.
var BytesLibBin = BytesLibMetaData.Bin

// DeployBytesLib deploys a new Ethereum contract, binding an instance of BytesLib to it.
func DeployBytesLib(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BytesLib, error) {
	parsed, err := BytesLibMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BytesLibBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BytesLib{BytesLibCaller: BytesLibCaller{contract: contract}, BytesLibTransactor: BytesLibTransactor{contract: contract}, BytesLibFilterer: BytesLibFilterer{contract: contract}}, nil
}

// BytesLib is an auto generated Go binding around an Ethereum contract.
type BytesLib struct {
	BytesLibCaller     // Read-only binding to the contract
	BytesLibTransactor // Write-only binding to the contract
	BytesLibFilterer   // Log filterer for contract events
}

// BytesLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type BytesLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BytesLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BytesLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BytesLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BytesLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BytesLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BytesLibSession struct {
	Contract     *BytesLib         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BytesLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BytesLibCallerSession struct {
	Contract *BytesLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// BytesLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BytesLibTransactorSession struct {
	Contract     *BytesLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BytesLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type BytesLibRaw struct {
	Contract *BytesLib // Generic contract binding to access the raw methods on
}

// BytesLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BytesLibCallerRaw struct {
	Contract *BytesLibCaller // Generic read-only contract binding to access the raw methods on
}

// BytesLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BytesLibTransactorRaw struct {
	Contract *BytesLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBytesLib creates a new instance of BytesLib, bound to a specific deployed contract.
func NewBytesLib(address common.Address, backend bind.ContractBackend) (*BytesLib, error) {
	contract, err := bindBytesLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BytesLib{BytesLibCaller: BytesLibCaller{contract: contract}, BytesLibTransactor: BytesLibTransactor{contract: contract}, BytesLibFilterer: BytesLibFilterer{contract: contract}}, nil
}

// NewBytesLibCaller creates a new read-only instance of BytesLib, bound to a specific deployed contract.
func NewBytesLibCaller(address common.Address, caller bind.ContractCaller) (*BytesLibCaller, error) {
	contract, err := bindBytesLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BytesLibCaller{contract: contract}, nil
}

// NewBytesLibTransactor creates a new write-only instance of BytesLib, bound to a specific deployed contract.
func NewBytesLibTransactor(address common.Address, transactor bind.ContractTransactor) (*BytesLibTransactor, error) {
	contract, err := bindBytesLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BytesLibTransactor{contract: contract}, nil
}

// NewBytesLibFilterer creates a new log filterer instance of BytesLib, bound to a specific deployed contract.
func NewBytesLibFilterer(address common.Address, filterer bind.ContractFilterer) (*BytesLibFilterer, error) {
	contract, err := bindBytesLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BytesLibFilterer{contract: contract}, nil
}

// bindBytesLib binds a generic wrapper to an already deployed contract.
func bindBytesLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BytesLibABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BytesLib *BytesLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BytesLib.Contract.BytesLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BytesLib *BytesLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BytesLib.Contract.BytesLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BytesLib *BytesLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BytesLib.Contract.BytesLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BytesLib *BytesLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BytesLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BytesLib *BytesLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BytesLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BytesLib *BytesLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BytesLib.Contract.contract.Transact(opts, method, params...)
}

// HelpersMetaData contains all meta data concerning the Helpers contract.
var HelpersMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea26469706673582212206fd1007332792aef66022ee699272b0c6e6d1fbb8d9206f36306cdab53d33f2d64736f6c63430008170033",
}

// HelpersABI is the input ABI used to generate the binding from.
// Deprecated: Use HelpersMetaData.ABI instead.
var HelpersABI = HelpersMetaData.ABI

// HelpersBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use HelpersMetaData.Bin instead.
var HelpersBin = HelpersMetaData.Bin

// DeployHelpers deploys a new Ethereum contract, binding an instance of Helpers to it.
func DeployHelpers(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Helpers, error) {
	parsed, err := HelpersMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HelpersBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Helpers{HelpersCaller: HelpersCaller{contract: contract}, HelpersTransactor: HelpersTransactor{contract: contract}, HelpersFilterer: HelpersFilterer{contract: contract}}, nil
}

// Helpers is an auto generated Go binding around an Ethereum contract.
type Helpers struct {
	HelpersCaller     // Read-only binding to the contract
	HelpersTransactor // Write-only binding to the contract
	HelpersFilterer   // Log filterer for contract events
}

// HelpersCaller is an auto generated read-only Go binding around an Ethereum contract.
type HelpersCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelpersTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HelpersTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelpersFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HelpersFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelpersSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HelpersSession struct {
	Contract     *Helpers          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HelpersCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HelpersCallerSession struct {
	Contract *HelpersCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// HelpersTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HelpersTransactorSession struct {
	Contract     *HelpersTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// HelpersRaw is an auto generated low-level Go binding around an Ethereum contract.
type HelpersRaw struct {
	Contract *Helpers // Generic contract binding to access the raw methods on
}

// HelpersCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HelpersCallerRaw struct {
	Contract *HelpersCaller // Generic read-only contract binding to access the raw methods on
}

// HelpersTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HelpersTransactorRaw struct {
	Contract *HelpersTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHelpers creates a new instance of Helpers, bound to a specific deployed contract.
func NewHelpers(address common.Address, backend bind.ContractBackend) (*Helpers, error) {
	contract, err := bindHelpers(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Helpers{HelpersCaller: HelpersCaller{contract: contract}, HelpersTransactor: HelpersTransactor{contract: contract}, HelpersFilterer: HelpersFilterer{contract: contract}}, nil
}

// NewHelpersCaller creates a new read-only instance of Helpers, bound to a specific deployed contract.
func NewHelpersCaller(address common.Address, caller bind.ContractCaller) (*HelpersCaller, error) {
	contract, err := bindHelpers(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HelpersCaller{contract: contract}, nil
}

// NewHelpersTransactor creates a new write-only instance of Helpers, bound to a specific deployed contract.
func NewHelpersTransactor(address common.Address, transactor bind.ContractTransactor) (*HelpersTransactor, error) {
	contract, err := bindHelpers(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HelpersTransactor{contract: contract}, nil
}

// NewHelpersFilterer creates a new log filterer instance of Helpers, bound to a specific deployed contract.
func NewHelpersFilterer(address common.Address, filterer bind.ContractFilterer) (*HelpersFilterer, error) {
	contract, err := bindHelpers(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HelpersFilterer{contract: contract}, nil
}

// bindHelpers binds a generic wrapper to an already deployed contract.
func bindHelpers(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HelpersABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Helpers *HelpersRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Helpers.Contract.HelpersCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Helpers *HelpersRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Helpers.Contract.HelpersTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Helpers *HelpersRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Helpers.Contract.HelpersTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Helpers *HelpersCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Helpers.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Helpers *HelpersTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Helpers.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Helpers *HelpersTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Helpers.Contract.contract.Transact(opts, method, params...)
}

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
type IACUCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IACUTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IACUTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IACUFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IACUFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IACUSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IACUSession struct {
	Contract     *IACU             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IACUCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IACUCallerSession struct {
	Contract *IACUCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IACUTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IACUTransactorSession struct {
	Contract     *IACUTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IACURaw is an auto generated low-level Go binding around an Ethereum contract.
type IACURaw struct {
	Contract *IACU // Generic contract binding to access the raw methods on
}

// IACUCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IACUCallerRaw struct {
	Contract *IACUCaller // Generic read-only contract binding to access the raw methods on
}

// IACUTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IACUTransactorRaw struct {
	Contract *IACUTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIACU creates a new instance of IACU, bound to a specific deployed contract.
func NewIACU(address common.Address, backend bind.ContractBackend) (*IACU, error) {
	contract, err := bindIACU(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IACU{IACUCaller: IACUCaller{contract: contract}, IACUTransactor: IACUTransactor{contract: contract}, IACUFilterer: IACUFilterer{contract: contract}}, nil
}

// NewIACUCaller creates a new read-only instance of IACU, bound to a specific deployed contract.
func NewIACUCaller(address common.Address, caller bind.ContractCaller) (*IACUCaller, error) {
	contract, err := bindIACU(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IACUCaller{contract: contract}, nil
}

// NewIACUTransactor creates a new write-only instance of IACU, bound to a specific deployed contract.
func NewIACUTransactor(address common.Address, transactor bind.ContractTransactor) (*IACUTransactor, error) {
	contract, err := bindIACU(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IACUTransactor{contract: contract}, nil
}

// NewIACUFilterer creates a new log filterer instance of IACU, bound to a specific deployed contract.
func NewIACUFilterer(address common.Address, filterer bind.ContractFilterer) (*IACUFilterer, error) {
	contract, err := bindIACU(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IACUFilterer{contract: contract}, nil
}

// bindIACU binds a generic wrapper to an already deployed contract.
func bindIACU(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IACUABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IACU *IACURaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IACU.Contract.IACUCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IACU *IACURaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IACU.Contract.IACUTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IACU *IACURaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IACU.Contract.IACUTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IACU *IACUCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IACU.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IACU *IACUTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IACU.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IACU *IACUTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IACU.Contract.contract.Transact(opts, method, params...)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACUTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _IACU.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACUSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IACU.Contract.SetOperator(&_IACU.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACUTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IACU.Contract.SetOperator(&_IACU.TransactOpts, operator)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACUTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _IACU.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACUSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IACU.Contract.SetOracle(&_IACU.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACUTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IACU.Contract.SetOracle(&_IACU.TransactOpts, oracle)
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACUTransactor) Update(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IACU.contract.Transact(opts, "update")
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACUSession) Update() (*types.Transaction, error) {
	return _IACU.Contract.Update(&_IACU.TransactOpts)
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACUTransactorSession) Update() (*types.Transaction, error) {
	return _IACU.Contract.Update(&_IACU.TransactOpts)
}

// IAccountabilityMetaData contains all meta data concerning the IAccountability contract.
var IAccountabilityMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"InnocenceProven\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewAccusation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewFaultProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"releaseBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isJailbound\",\"type\":\"bool\"}],\"name\":\"SlashingEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_epochEnd\",\"type\":\"bool\"}],\"name\":\"finalize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newPeriod\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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
type IAccountabilityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountabilityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAccountabilityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountabilityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAccountabilityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAccountabilitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAccountabilitySession struct {
	Contract     *IAccountability  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAccountabilityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAccountabilityCallerSession struct {
	Contract *IAccountabilityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IAccountabilityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAccountabilityTransactorSession struct {
	Contract     *IAccountabilityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IAccountabilityRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAccountabilityRaw struct {
	Contract *IAccountability // Generic contract binding to access the raw methods on
}

// IAccountabilityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAccountabilityCallerRaw struct {
	Contract *IAccountabilityCaller // Generic read-only contract binding to access the raw methods on
}

// IAccountabilityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAccountabilityTransactorRaw struct {
	Contract *IAccountabilityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAccountability creates a new instance of IAccountability, bound to a specific deployed contract.
func NewIAccountability(address common.Address, backend bind.ContractBackend) (*IAccountability, error) {
	contract, err := bindIAccountability(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAccountability{IAccountabilityCaller: IAccountabilityCaller{contract: contract}, IAccountabilityTransactor: IAccountabilityTransactor{contract: contract}, IAccountabilityFilterer: IAccountabilityFilterer{contract: contract}}, nil
}

// NewIAccountabilityCaller creates a new read-only instance of IAccountability, bound to a specific deployed contract.
func NewIAccountabilityCaller(address common.Address, caller bind.ContractCaller) (*IAccountabilityCaller, error) {
	contract, err := bindIAccountability(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityCaller{contract: contract}, nil
}

// NewIAccountabilityTransactor creates a new write-only instance of IAccountability, bound to a specific deployed contract.
func NewIAccountabilityTransactor(address common.Address, transactor bind.ContractTransactor) (*IAccountabilityTransactor, error) {
	contract, err := bindIAccountability(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityTransactor{contract: contract}, nil
}

// NewIAccountabilityFilterer creates a new log filterer instance of IAccountability, bound to a specific deployed contract.
func NewIAccountabilityFilterer(address common.Address, filterer bind.ContractFilterer) (*IAccountabilityFilterer, error) {
	contract, err := bindIAccountability(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAccountabilityFilterer{contract: contract}, nil
}

// bindIAccountability binds a generic wrapper to an already deployed contract.
func bindIAccountability(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAccountabilityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAccountability *IAccountabilityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAccountability.Contract.IAccountabilityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAccountability *IAccountabilityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccountability.Contract.IAccountabilityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAccountability *IAccountabilityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAccountability.Contract.IAccountabilityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAccountability *IAccountabilityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAccountability.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAccountability *IAccountabilityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAccountability.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAccountability *IAccountabilityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAccountability.Contract.contract.Transact(opts, method, params...)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_IAccountability *IAccountabilityTransactor) DistributeRewards(opts *bind.TransactOpts, _validator common.Address) (*types.Transaction, error) {
	return _IAccountability.contract.Transact(opts, "distributeRewards", _validator)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_IAccountability *IAccountabilitySession) DistributeRewards(_validator common.Address) (*types.Transaction, error) {
	return _IAccountability.Contract.DistributeRewards(&_IAccountability.TransactOpts, _validator)
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_IAccountability *IAccountabilityTransactorSession) DistributeRewards(_validator common.Address) (*types.Transaction, error) {
	return _IAccountability.Contract.DistributeRewards(&_IAccountability.TransactOpts, _validator)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_IAccountability *IAccountabilityTransactor) Finalize(opts *bind.TransactOpts, _epochEnd bool) (*types.Transaction, error) {
	return _IAccountability.contract.Transact(opts, "finalize", _epochEnd)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_IAccountability *IAccountabilitySession) Finalize(_epochEnd bool) (*types.Transaction, error) {
	return _IAccountability.Contract.Finalize(&_IAccountability.TransactOpts, _epochEnd)
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_IAccountability *IAccountabilityTransactorSession) Finalize(_epochEnd bool) (*types.Transaction, error) {
	return _IAccountability.Contract.Finalize(&_IAccountability.TransactOpts, _epochEnd)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_IAccountability *IAccountabilityTransactor) SetEpochPeriod(opts *bind.TransactOpts, _newPeriod *big.Int) (*types.Transaction, error) {
	return _IAccountability.contract.Transact(opts, "setEpochPeriod", _newPeriod)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_IAccountability *IAccountabilitySession) SetEpochPeriod(_newPeriod *big.Int) (*types.Transaction, error) {
	return _IAccountability.Contract.SetEpochPeriod(&_IAccountability.TransactOpts, _newPeriod)
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_IAccountability *IAccountabilityTransactorSession) SetEpochPeriod(_newPeriod *big.Int) (*types.Transaction, error) {
	return _IAccountability.Contract.SetEpochPeriod(&_IAccountability.TransactOpts, _newPeriod)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Offender common.Address
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Offender common.Address
	Severity *big.Int
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Offender common.Address
	Severity *big.Int
	Id       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Validator    common.Address
	Amount       *big.Int
	ReleaseBlock *big.Int
	IsJailbound  bool
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSlashingEvent is a free log retrieval operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
func (_IAccountability *IAccountabilityFilterer) FilterSlashingEvent(opts *bind.FilterOpts) (*IAccountabilitySlashingEventIterator, error) {

	logs, sub, err := _IAccountability.contract.FilterLogs(opts, "SlashingEvent")
	if err != nil {
		return nil, err
	}
	return &IAccountabilitySlashingEventIterator{contract: _IAccountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
}

// WatchSlashingEvent is a free log subscription operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
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

// ParseSlashingEvent is a log parse operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
//
// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
func (_IAccountability *IAccountabilityFilterer) ParseSlashingEvent(log types.Log) (*IAccountabilitySlashingEvent, error) {
	event := new(IAccountabilitySlashingEvent)
	if err := _IAccountability.contract.UnpackLog(event, "SlashingEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

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
type IAutonityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAutonityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IAutonityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAutonityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IAutonityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAutonitySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IAutonitySession struct {
	Contract     *IAutonity        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAutonityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IAutonityCallerSession struct {
	Contract *IAutonityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IAutonityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IAutonityTransactorSession struct {
	Contract     *IAutonityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IAutonityRaw is an auto generated low-level Go binding around an Ethereum contract.
type IAutonityRaw struct {
	Contract *IAutonity // Generic contract binding to access the raw methods on
}

// IAutonityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IAutonityCallerRaw struct {
	Contract *IAutonityCaller // Generic read-only contract binding to access the raw methods on
}

// IAutonityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IAutonityTransactorRaw struct {
	Contract *IAutonityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAutonity creates a new instance of IAutonity, bound to a specific deployed contract.
func NewIAutonity(address common.Address, backend bind.ContractBackend) (*IAutonity, error) {
	contract, err := bindIAutonity(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAutonity{IAutonityCaller: IAutonityCaller{contract: contract}, IAutonityTransactor: IAutonityTransactor{contract: contract}, IAutonityFilterer: IAutonityFilterer{contract: contract}}, nil
}

// NewIAutonityCaller creates a new read-only instance of IAutonity, bound to a specific deployed contract.
func NewIAutonityCaller(address common.Address, caller bind.ContractCaller) (*IAutonityCaller, error) {
	contract, err := bindIAutonity(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAutonityCaller{contract: contract}, nil
}

// NewIAutonityTransactor creates a new write-only instance of IAutonity, bound to a specific deployed contract.
func NewIAutonityTransactor(address common.Address, transactor bind.ContractTransactor) (*IAutonityTransactor, error) {
	contract, err := bindIAutonity(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAutonityTransactor{contract: contract}, nil
}

// NewIAutonityFilterer creates a new log filterer instance of IAutonity, bound to a specific deployed contract.
func NewIAutonityFilterer(address common.Address, filterer bind.ContractFilterer) (*IAutonityFilterer, error) {
	contract, err := bindIAutonity(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAutonityFilterer{contract: contract}, nil
}

// bindIAutonity binds a generic wrapper to an already deployed contract.
func bindIAutonity(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IAutonityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAutonity *IAutonityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutonity.Contract.IAutonityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAutonity *IAutonityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutonity.Contract.IAutonityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAutonity *IAutonityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutonity.Contract.IAutonityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAutonity *IAutonityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAutonity.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAutonity *IAutonityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAutonity.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAutonity *IAutonityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAutonity.Contract.contract.Transact(opts, method, params...)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonityCaller) GetOperator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getOperator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonitySession) GetOperator() (common.Address, error) {
	return _IAutonity.Contract.GetOperator(&_IAutonity.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonityCallerSession) GetOperator() (common.Address, error) {
	return _IAutonity.Contract.GetOperator(&_IAutonity.CallOpts)
}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_IAutonity *IAutonityCaller) GetOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAutonity.contract.Call(opts, &out, "getOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_IAutonity *IAutonitySession) GetOracle() (common.Address, error) {
	return _IAutonity.Contract.GetOracle(&_IAutonity.CallOpts)
}

// GetOracle is a free data retrieval call binding the contract method 0x833b1fce.
//
// Solidity: function getOracle() view returns(address)
func (_IAutonity *IAutonityCallerSession) GetOracle() (common.Address, error) {
	return _IAutonity.Contract.GetOracle(&_IAutonity.CallOpts)
}

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
type IERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type IERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IERC20Session struct {
	Contract     *IERC20           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IERC20CallerSession struct {
	Contract *IERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IERC20TransactorSession struct {
	Contract     *IERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type IERC20Raw struct {
	Contract *IERC20 // Generic contract binding to access the raw methods on
}

// IERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IERC20CallerRaw struct {
	Contract *IERC20Caller // Generic read-only contract binding to access the raw methods on
}

// IERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IERC20TransactorRaw struct {
	Contract *IERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC20 creates a new instance of IERC20, bound to a specific deployed contract.
func NewIERC20(address common.Address, backend bind.ContractBackend) (*IERC20, error) {
	contract, err := bindIERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20{IERC20Caller: IERC20Caller{contract: contract}, IERC20Transactor: IERC20Transactor{contract: contract}, IERC20Filterer: IERC20Filterer{contract: contract}}, nil
}

// NewIERC20Caller creates a new read-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Caller(address common.Address, caller bind.ContractCaller) (*IERC20Caller, error) {
	contract, err := bindIERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Caller{contract: contract}, nil
}

// NewIERC20Transactor creates a new write-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC20Transactor, error) {
	contract, err := bindIERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Transactor{contract: contract}, nil
}

// NewIERC20Filterer creates a new log filterer instance of IERC20, bound to a specific deployed contract.
func NewIERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC20Filterer, error) {
	contract, err := bindIERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC20Filterer{contract: contract}, nil
}

// bindIERC20 binds a generic wrapper to an already deployed contract.
func bindIERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20.Contract.IERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Session) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, sender, recipient, amount)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
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
type IOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IOracleSession struct {
	Contract     *IOracle          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IOracleCallerSession struct {
	Contract *IOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IOracleTransactorSession struct {
	Contract     *IOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type IOracleRaw struct {
	Contract *IOracle // Generic contract binding to access the raw methods on
}

// IOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IOracleCallerRaw struct {
	Contract *IOracleCaller // Generic read-only contract binding to access the raw methods on
}

// IOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IOracleTransactorRaw struct {
	Contract *IOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIOracle creates a new instance of IOracle, bound to a specific deployed contract.
func NewIOracle(address common.Address, backend bind.ContractBackend) (*IOracle, error) {
	contract, err := bindIOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IOracle{IOracleCaller: IOracleCaller{contract: contract}, IOracleTransactor: IOracleTransactor{contract: contract}, IOracleFilterer: IOracleFilterer{contract: contract}}, nil
}

// NewIOracleCaller creates a new read-only instance of IOracle, bound to a specific deployed contract.
func NewIOracleCaller(address common.Address, caller bind.ContractCaller) (*IOracleCaller, error) {
	contract, err := bindIOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IOracleCaller{contract: contract}, nil
}

// NewIOracleTransactor creates a new write-only instance of IOracle, bound to a specific deployed contract.
func NewIOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*IOracleTransactor, error) {
	contract, err := bindIOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IOracleTransactor{contract: contract}, nil
}

// NewIOracleFilterer creates a new log filterer instance of IOracle, bound to a specific deployed contract.
func NewIOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*IOracleFilterer, error) {
	contract, err := bindIOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IOracleFilterer{contract: contract}, nil
}

// bindIOracle binds a generic wrapper to an already deployed contract.
func bindIOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOracle *IOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOracle.Contract.IOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOracle *IOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOracle.Contract.IOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOracle *IOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOracle.Contract.IOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IOracle *IOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IOracle *IOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IOracle *IOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IOracle.Contract.contract.Transact(opts, method, params...)
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() view returns(uint256)
func (_IOracle *IOracleCaller) GetPrecision(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getPrecision")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() view returns(uint256)
func (_IOracle *IOracleSession) GetPrecision() (*big.Int, error) {
	return _IOracle.Contract.GetPrecision(&_IOracle.CallOpts)
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() view returns(uint256)
func (_IOracle *IOracleCallerSession) GetPrecision() (*big.Int, error) {
	return _IOracle.Contract.GetPrecision(&_IOracle.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_IOracle *IOracleCaller) GetRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_IOracle *IOracleSession) GetRound() (*big.Int, error) {
	return _IOracle.Contract.GetRound(&_IOracle.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_IOracle *IOracleCallerSession) GetRound() (*big.Int, error) {
	return _IOracle.Contract.GetRound(&_IOracle.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleCaller) GetRoundData(opts *bind.CallOpts, _round *big.Int, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.GetRoundData(&_IOracle.CallOpts, _round, _symbol)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleCallerSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.GetRoundData(&_IOracle.CallOpts, _round, _symbol)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[] _symbols)
func (_IOracle *IOracleCaller) GetSymbols(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getSymbols")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[] _symbols)
func (_IOracle *IOracleSession) GetSymbols() ([]string, error) {
	return _IOracle.Contract.GetSymbols(&_IOracle.CallOpts)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[] _symbols)
func (_IOracle *IOracleCallerSession) GetSymbols() ([]string, error) {
	return _IOracle.Contract.GetSymbols(&_IOracle.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_IOracle *IOracleCaller) GetVotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_IOracle *IOracleSession) GetVotePeriod() (*big.Int, error) {
	return _IOracle.Contract.GetVotePeriod(&_IOracle.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_IOracle *IOracleCallerSession) GetVotePeriod() (*big.Int, error) {
	return _IOracle.Contract.GetVotePeriod(&_IOracle.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_IOracle *IOracleCaller) GetVoters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "getVoters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_IOracle *IOracleSession) GetVoters() ([]common.Address, error) {
	return _IOracle.Contract.GetVoters(&_IOracle.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_IOracle *IOracleCallerSession) GetVoters() ([]common.Address, error) {
	return _IOracle.Contract.GetVoters(&_IOracle.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleCaller) LatestRoundData(opts *bind.CallOpts, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _IOracle.contract.Call(opts, &out, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.LatestRoundData(&_IOracle.CallOpts, _symbol)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_IOracle *IOracleCallerSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _IOracle.Contract.LatestRoundData(&_IOracle.CallOpts, _symbol)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracleTransactor) Finalize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "finalize")
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracleSession) Finalize() (*types.Transaction, error) {
	return _IOracle.Contract.Finalize(&_IOracle.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_IOracle *IOracleTransactorSession) Finalize() (*types.Transaction, error) {
	return _IOracle.Contract.Finalize(&_IOracle.TransactOpts)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracleTransactor) SetOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "setOperator", _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracleSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _IOracle.Contract.SetOperator(&_IOracle.TransactOpts, _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracleTransactorSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _IOracle.Contract.SetOperator(&_IOracle.TransactOpts, _operator)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracleTransactor) SetSymbols(opts *bind.TransactOpts, _symbols []string) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "setSymbols", _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracleSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _IOracle.Contract.SetSymbols(&_IOracle.TransactOpts, _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracleTransactorSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _IOracle.Contract.SetSymbols(&_IOracle.TransactOpts, _symbols)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_IOracle *IOracleTransactor) SetVoters(opts *bind.TransactOpts, _newVoters []common.Address) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "setVoters", _newVoters)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_IOracle *IOracleSession) SetVoters(_newVoters []common.Address) (*types.Transaction, error) {
	return _IOracle.Contract.SetVoters(&_IOracle.TransactOpts, _newVoters)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_IOracle *IOracleTransactorSession) SetVoters(_newVoters []common.Address) (*types.Transaction, error) {
	return _IOracle.Contract.SetVoters(&_IOracle.TransactOpts, _newVoters)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_IOracle *IOracleTransactor) Vote(opts *bind.TransactOpts, _commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _IOracle.contract.Transact(opts, "vote", _commit, _reports, _salt)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_IOracle *IOracleSession) Vote(_commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.Vote(&_IOracle.TransactOpts, _commit, _reports, _salt)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_IOracle *IOracleTransactorSession) Vote(_commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _IOracle.Contract.Vote(&_IOracle.TransactOpts, _commit, _reports, _salt)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Round      *big.Int
	Height     *big.Int
	Timestamp  *big.Int
	VotePeriod *big.Int
	Raw        types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Symbols []string
	Round   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Voter common.Address
	Votes []*big.Int
	Raw   types.Log // Blockchain specific contextual infos
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
type IStabilizationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStabilizationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IStabilizationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStabilizationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IStabilizationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IStabilizationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IStabilizationSession struct {
	Contract     *IStabilization   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IStabilizationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IStabilizationCallerSession struct {
	Contract *IStabilizationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IStabilizationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IStabilizationTransactorSession struct {
	Contract     *IStabilizationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IStabilizationRaw is an auto generated low-level Go binding around an Ethereum contract.
type IStabilizationRaw struct {
	Contract *IStabilization // Generic contract binding to access the raw methods on
}

// IStabilizationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IStabilizationCallerRaw struct {
	Contract *IStabilizationCaller // Generic read-only contract binding to access the raw methods on
}

// IStabilizationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IStabilizationTransactorRaw struct {
	Contract *IStabilizationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIStabilization creates a new instance of IStabilization, bound to a specific deployed contract.
func NewIStabilization(address common.Address, backend bind.ContractBackend) (*IStabilization, error) {
	contract, err := bindIStabilization(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IStabilization{IStabilizationCaller: IStabilizationCaller{contract: contract}, IStabilizationTransactor: IStabilizationTransactor{contract: contract}, IStabilizationFilterer: IStabilizationFilterer{contract: contract}}, nil
}

// NewIStabilizationCaller creates a new read-only instance of IStabilization, bound to a specific deployed contract.
func NewIStabilizationCaller(address common.Address, caller bind.ContractCaller) (*IStabilizationCaller, error) {
	contract, err := bindIStabilization(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IStabilizationCaller{contract: contract}, nil
}

// NewIStabilizationTransactor creates a new write-only instance of IStabilization, bound to a specific deployed contract.
func NewIStabilizationTransactor(address common.Address, transactor bind.ContractTransactor) (*IStabilizationTransactor, error) {
	contract, err := bindIStabilization(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IStabilizationTransactor{contract: contract}, nil
}

// NewIStabilizationFilterer creates a new log filterer instance of IStabilization, bound to a specific deployed contract.
func NewIStabilizationFilterer(address common.Address, filterer bind.ContractFilterer) (*IStabilizationFilterer, error) {
	contract, err := bindIStabilization(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IStabilizationFilterer{contract: contract}, nil
}

// bindIStabilization binds a generic wrapper to an already deployed contract.
func bindIStabilization(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IStabilizationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IStabilization *IStabilizationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IStabilization.Contract.IStabilizationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IStabilization *IStabilizationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IStabilization.Contract.IStabilizationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IStabilization *IStabilizationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IStabilization.Contract.IStabilizationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IStabilization *IStabilizationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IStabilization.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IStabilization *IStabilizationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IStabilization.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IStabilization *IStabilizationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IStabilization.Contract.contract.Transact(opts, method, params...)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilizationTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _IStabilization.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilizationSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetOperator(&_IStabilization.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilizationTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetOperator(&_IStabilization.TransactOpts, operator)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilizationTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _IStabilization.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilizationSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetOracle(&_IStabilization.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilizationTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _IStabilization.Contract.SetOracle(&_IStabilization.TransactOpts, oracle)
}

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
type ISupplyControlCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISupplyControlTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISupplyControlTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISupplyControlFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISupplyControlFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISupplyControlSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISupplyControlSession struct {
	Contract     *ISupplyControl   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ISupplyControlCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISupplyControlCallerSession struct {
	Contract *ISupplyControlCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ISupplyControlTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISupplyControlTransactorSession struct {
	Contract     *ISupplyControlTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ISupplyControlRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISupplyControlRaw struct {
	Contract *ISupplyControl // Generic contract binding to access the raw methods on
}

// ISupplyControlCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISupplyControlCallerRaw struct {
	Contract *ISupplyControlCaller // Generic read-only contract binding to access the raw methods on
}

// ISupplyControlTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISupplyControlTransactorRaw struct {
	Contract *ISupplyControlTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISupplyControl creates a new instance of ISupplyControl, bound to a specific deployed contract.
func NewISupplyControl(address common.Address, backend bind.ContractBackend) (*ISupplyControl, error) {
	contract, err := bindISupplyControl(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISupplyControl{ISupplyControlCaller: ISupplyControlCaller{contract: contract}, ISupplyControlTransactor: ISupplyControlTransactor{contract: contract}, ISupplyControlFilterer: ISupplyControlFilterer{contract: contract}}, nil
}

// NewISupplyControlCaller creates a new read-only instance of ISupplyControl, bound to a specific deployed contract.
func NewISupplyControlCaller(address common.Address, caller bind.ContractCaller) (*ISupplyControlCaller, error) {
	contract, err := bindISupplyControl(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISupplyControlCaller{contract: contract}, nil
}

// NewISupplyControlTransactor creates a new write-only instance of ISupplyControl, bound to a specific deployed contract.
func NewISupplyControlTransactor(address common.Address, transactor bind.ContractTransactor) (*ISupplyControlTransactor, error) {
	contract, err := bindISupplyControl(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISupplyControlTransactor{contract: contract}, nil
}

// NewISupplyControlFilterer creates a new log filterer instance of ISupplyControl, bound to a specific deployed contract.
func NewISupplyControlFilterer(address common.Address, filterer bind.ContractFilterer) (*ISupplyControlFilterer, error) {
	contract, err := bindISupplyControl(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISupplyControlFilterer{contract: contract}, nil
}

// bindISupplyControl binds a generic wrapper to an already deployed contract.
func bindISupplyControl(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISupplyControlABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISupplyControl *ISupplyControlRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISupplyControl.Contract.ISupplyControlCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISupplyControl *ISupplyControlRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISupplyControl.Contract.ISupplyControlTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISupplyControl *ISupplyControlRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISupplyControl.Contract.ISupplyControlTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISupplyControl *ISupplyControlCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISupplyControl.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISupplyControl *ISupplyControlTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISupplyControl.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISupplyControl *ISupplyControlTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISupplyControl.Contract.contract.Transact(opts, method, params...)
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCaller) AvailableSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ISupplyControl.contract.Call(opts, &out, "availableSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlSession) AvailableSupply() (*big.Int, error) {
	return _ISupplyControl.Contract.AvailableSupply(&_ISupplyControl.CallOpts)
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCallerSession) AvailableSupply() (*big.Int, error) {
	return _ISupplyControl.Contract.AvailableSupply(&_ISupplyControl.CallOpts)
}

// Stabilizer is a free data retrieval call binding the contract method 0x7e47961c.
//
// Solidity: function stabilizer() view returns(address)
func (_ISupplyControl *ISupplyControlCaller) Stabilizer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ISupplyControl.contract.Call(opts, &out, "stabilizer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Stabilizer is a free data retrieval call binding the contract method 0x7e47961c.
//
// Solidity: function stabilizer() view returns(address)
func (_ISupplyControl *ISupplyControlSession) Stabilizer() (common.Address, error) {
	return _ISupplyControl.Contract.Stabilizer(&_ISupplyControl.CallOpts)
}

// Stabilizer is a free data retrieval call binding the contract method 0x7e47961c.
//
// Solidity: function stabilizer() view returns(address)
func (_ISupplyControl *ISupplyControlCallerSession) Stabilizer() (common.Address, error) {
	return _ISupplyControl.Contract.Stabilizer(&_ISupplyControl.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ISupplyControl.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlSession) TotalSupply() (*big.Int, error) {
	return _ISupplyControl.Contract.TotalSupply(&_ISupplyControl.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControlCallerSession) TotalSupply() (*big.Int, error) {
	return _ISupplyControl.Contract.TotalSupply(&_ISupplyControl.CallOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControlTransactor) Burn(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISupplyControl.contract.Transact(opts, "burn")
}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControlSession) Burn() (*types.Transaction, error) {
	return _ISupplyControl.Contract.Burn(&_ISupplyControl.TransactOpts)
}

// Burn is a paid mutator transaction binding the contract method 0x44df8e70.
//
// Solidity: function burn() payable returns()
func (_ISupplyControl *ISupplyControlTransactorSession) Burn() (*types.Transaction, error) {
	return _ISupplyControl.Contract.Burn(&_ISupplyControl.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControlTransactor) Mint(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISupplyControl.contract.Transact(opts, "mint", recipient, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControlSession) Mint(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISupplyControl.Contract.Mint(&_ISupplyControl.TransactOpts, recipient, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControlTransactorSession) Mint(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ISupplyControl.Contract.Mint(&_ISupplyControl.TransactOpts, recipient, amount)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControlTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _ISupplyControl.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControlSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _ISupplyControl.Contract.SetOperator(&_ISupplyControl.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControlTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _ISupplyControl.Contract.SetOperator(&_ISupplyControl.TransactOpts, operator)
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControlTransactor) SetStabilizer(opts *bind.TransactOpts, stabilizer_ common.Address) (*types.Transaction, error) {
	return _ISupplyControl.contract.Transact(opts, "setStabilizer", stabilizer_)
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControlSession) SetStabilizer(stabilizer_ common.Address) (*types.Transaction, error) {
	return _ISupplyControl.Contract.SetStabilizer(&_ISupplyControl.TransactOpts, stabilizer_)
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControlTransactorSession) SetStabilizer(stabilizer_ common.Address) (*types.Transaction, error) {
	return _ISupplyControl.Contract.SetStabilizer(&_ISupplyControl.TransactOpts, stabilizer_)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
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

// LiquidLogicMetaData contains all meta data concerning the LiquidLogic contract.
var LiquidLogicMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_autonity\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_FACTOR_UNIT_RECIP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"redistribute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unclaimedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"unlockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"2f2c3f2e": "COMMISSION_RATE_PRECISION()",
		"187cf4d7": "FEE_FACTOR_UNIT_RECIP()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"9dc29fac": "burn(address,uint256)",
		"372500ab": "claimRewards()",
		"5ea1d6f8": "commissionRate()",
		"282d3fdf": "lock(address,uint256)",
		"40c10f19": "mint(address,uint256)",
		"fb489a7b": "redistribute()",
		"19fac8fd": "setCommissionRate(uint256)",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"61d027b3": "treasury()",
		"949813b8": "unclaimedRewards(address)",
		"7eee288d": "unlock(address,uint256)",
		"84955c88": "unlockedBalanceOf(address)",
	},
	Bin: "0x608060405234801561000f575f80fd5b50604051610f80380380610f8083398101604081905261002e91610052565b5f80546001600160a01b0319166001600160a01b039290921691909117905561007f565b5f60208284031215610062575f80fd5b81516001600160a01b0381168114610078575f80fd5b9392505050565b610ef48061008c5f395ff3fe608060405260043610610110575f3560e01c80635ea1d6f81161009d578063949813b811610062578063949813b8146102e25780639dc29fac14610301578063a9059cbb14610320578063dd62ed3e1461033f578063fb489a7b14610383575f80fd5b80635ea1d6f81461022457806361d027b31461023957806370a08231146102705780637eee288d146102a457806384955c88146102c3575f80fd5b806323b872dd116100e357806323b872dd1461019e578063282d3fdf146101bd5780632f2c3f2e146101dc578063372500ab146101f157806340c10f1914610205575f80fd5b8063095ea7b31461011457806318160ddd14610148578063187cf4d71461016657806319fac8fd1461017d575b5f80fd5b34801561011f575f80fd5b5061013361012e366004610d1d565b61038b565b60405190151581526020015b60405180910390f35b348015610153575f80fd5b506004545b60405190815260200161013f565b348015610171575f80fd5b50610158633b9aca0081565b348015610188575f80fd5b5061019c610197366004610d45565b6103a1565b005b3480156101a9575f80fd5b506101336101b8366004610d5c565b6103d8565b3480156101c8575f80fd5b5061019c6101d7366004610d1d565b6104c9565b3480156101e7575f80fd5b5061015861271081565b3480156101fc575f80fd5b5061019c6105ab565b348015610210575f80fd5b5061019c61021f366004610d1d565b610655565b34801561022f575f80fd5b5061015860095481565b348015610244575f80fd5b50600854610258906001600160a01b031681565b6040516001600160a01b03909116815260200161013f565b34801561027b575f80fd5b5061015861028a366004610d95565b6001600160a01b03165f9081526001602052604090205490565b3480156102af575f80fd5b5061019c6102be366004610d1d565b6106ba565b3480156102ce575f80fd5b506101586102dd366004610d95565b61077d565b3480156102ed575f80fd5b506101586102fc366004610d95565b6107aa565b34801561030c575f80fd5b5061019c61031b366004610d1d565b6107d6565b34801561032b575f80fd5b5061013361033a366004610d1d565b610833565b34801561034a575f80fd5b50610158610359366004610db5565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b61015861087e565b5f6103973384846109df565b5060015b92915050565b5f546001600160a01b031633146103d35760405162461bcd60e51b81526004016103ca90610de6565b60405180910390fd5b600955565b6001600160a01b0383165f9081526003602090815260408083203384529091528120548281101561045c5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b60648201526084016103ca565b610470853361046b8685610e42565b6109df565b61047a8584610b02565b6104848484610be6565b836001600160a01b0316856001600160a01b03165f80516020610e9f833981519152856040516104b691815260200190565b60405180910390a3506001949350505050565b5f546001600160a01b031633146104f25760405162461bcd60e51b81526004016103ca90610de6565b6001600160a01b0382165f90815260026020908152604080832054600190925290912054829161052191610e42565b101561057b5760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b60648201526084016103ca565b6001600160a01b0382165f90815260026020526040812080548392906105a2908490610e55565b90915550505050565b5f6105b533610c2f565b335f81815260056020526040808220829055519293509183908381818185875af1925050503d805f8114610604576040519150601f19603f3d011682016040523d82523d5f602084013e610609565b606091505b50509050806106515760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b60448201526064016103ca565b5050565b5f546001600160a01b0316331461067e5760405162461bcd60e51b81526004016103ca90610de6565b6106888282610be6565b6040518181526001600160a01b038316905f905f80516020610e9f833981519152906020015b60405180910390a35050565b5f546001600160a01b031633146106e35760405162461bcd60e51b81526004016103ca90610de6565b6001600160a01b0382165f908152600260205260409020548111156107565760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b60648201526084016103ca565b6001600160a01b0382165f90815260026020526040812080548392906105a2908490610e42565b6001600160a01b0381165f90815260026020908152604080832054600190925282205461039b9190610e42565b5f6107b482610c91565b6001600160a01b0383165f9081526005602052604090205461039b9190610e55565b5f546001600160a01b031633146107ff5760405162461bcd60e51b81526004016103ca90610de6565b6108098282610b02565b6040518181525f906001600160a01b038416905f80516020610e9f833981519152906020016106ae565b5f61083e3383610b02565b6108488383610be6565b6040518281526001600160a01b0384169033905f80516020610e9f8339815191529060200160405180910390a350600192915050565b5f80546001600160a01b031633146108a85760405162461bcd60e51b81526004016103ca90610de6565b60095434905f90612710906108bd9084610e68565b6108c79190610e7f565b9050818111156109195760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f7220726577617264000000000000000060448201526064016103ca565b6109238183610e42565b6008546040519193506001600160a01b0316906108fc9083905f818181858888f193505050503d805f8114610973576040519150601f19603f3d011682016040523d82523d5f602084013e610978565b606091505b50506004545f915061098e633b9aca0085610e68565b6109989190610e7f565b9050806007546109a89190610e55565b6007556004545f90633b9aca00906109c09084610e68565b6109ca9190610e7f565b90506109d68184610e55565b94505050505090565b6001600160a01b038316610a415760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b60648201526084016103ca565b6001600160a01b038216610aa25760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b60648201526084016103ca565b6001600160a01b038381165f8181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610b0b82610c2f565b506001600160a01b0382165f90815260016020908152604080832054600290925290912054610b3a9082610e42565b821115610b895760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e6473000000000060448201526064016103ca565b610b938282610e42565b6001600160a01b0384165f90815260016020526040902055808203610bcb576001600160a01b0383165f908152600660205260408120555b8160045f828254610bdc9190610e42565b9091555050505050565b610bef82610c2f565b506001600160a01b0382165f9081526001602052604081208054839290610c17908490610e55565b925050819055508060045f8282546105a29190610e55565b5f80610c3a83610c91565b6001600160a01b0384165f90815260056020526040902054909150610c60908290610e55565b6001600160a01b039093165f9081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b0381165f90815260016020526040812054808203610cb857505f92915050565b6001600160a01b0383165f90815260066020526040812054600754610cdd9190610e42565b90505f633b9aca00610cef8484610e68565b610cf99190610e7f565b95945050505050565b80356001600160a01b0381168114610d18575f80fd5b919050565b5f8060408385031215610d2e575f80fd5b610d3783610d02565b946020939093013593505050565b5f60208284031215610d55575f80fd5b5035919050565b5f805f60608486031215610d6e575f80fd5b610d7784610d02565b9250610d8560208501610d02565b9150604084013590509250925092565b5f60208284031215610da5575f80fd5b610dae82610d02565b9392505050565b5f8060408385031215610dc6575f80fd5b610dcf83610d02565b9150610ddd60208401610d02565b90509250929050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561039b5761039b610e2e565b8082018082111561039b5761039b610e2e565b808202811582820484141761039b5761039b610e2e565b5f82610e9957634e487b7160e01b5f52601260045260245ffd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220a9171ec1ffebc6ff9e96eb0c1adf7b9e4e81605dc47e7a79c63c4befebd2e08064736f6c63430008170033",
}

// LiquidLogicABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidLogicMetaData.ABI instead.
var LiquidLogicABI = LiquidLogicMetaData.ABI

// Deprecated: Use LiquidLogicMetaData.Sigs instead.
// LiquidLogicFuncSigs maps the 4-byte function signature to its string representation.
var LiquidLogicFuncSigs = LiquidLogicMetaData.Sigs

// LiquidLogicBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LiquidLogicMetaData.Bin instead.
var LiquidLogicBin = LiquidLogicMetaData.Bin

// DeployLiquidLogic deploys a new Ethereum contract, binding an instance of LiquidLogic to it.
func DeployLiquidLogic(auth *bind.TransactOpts, backend bind.ContractBackend, _autonity common.Address) (common.Address, *types.Transaction, *LiquidLogic, error) {
	parsed, err := LiquidLogicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidLogicBin), backend, _autonity)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LiquidLogic{LiquidLogicCaller: LiquidLogicCaller{contract: contract}, LiquidLogicTransactor: LiquidLogicTransactor{contract: contract}, LiquidLogicFilterer: LiquidLogicFilterer{contract: contract}}, nil
}

// LiquidLogic is an auto generated Go binding around an Ethereum contract.
type LiquidLogic struct {
	LiquidLogicCaller     // Read-only binding to the contract
	LiquidLogicTransactor // Write-only binding to the contract
	LiquidLogicFilterer   // Log filterer for contract events
}

// LiquidLogicCaller is an auto generated read-only Go binding around an Ethereum contract.
type LiquidLogicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidLogicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidLogicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidLogicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidLogicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidLogicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LiquidLogicSession struct {
	Contract     *LiquidLogic      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LiquidLogicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LiquidLogicCallerSession struct {
	Contract *LiquidLogicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// LiquidLogicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LiquidLogicTransactorSession struct {
	Contract     *LiquidLogicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// LiquidLogicRaw is an auto generated low-level Go binding around an Ethereum contract.
type LiquidLogicRaw struct {
	Contract *LiquidLogic // Generic contract binding to access the raw methods on
}

// LiquidLogicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LiquidLogicCallerRaw struct {
	Contract *LiquidLogicCaller // Generic read-only contract binding to access the raw methods on
}

// LiquidLogicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LiquidLogicTransactorRaw struct {
	Contract *LiquidLogicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLiquidLogic creates a new instance of LiquidLogic, bound to a specific deployed contract.
func NewLiquidLogic(address common.Address, backend bind.ContractBackend) (*LiquidLogic, error) {
	contract, err := bindLiquidLogic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LiquidLogic{LiquidLogicCaller: LiquidLogicCaller{contract: contract}, LiquidLogicTransactor: LiquidLogicTransactor{contract: contract}, LiquidLogicFilterer: LiquidLogicFilterer{contract: contract}}, nil
}

// NewLiquidLogicCaller creates a new read-only instance of LiquidLogic, bound to a specific deployed contract.
func NewLiquidLogicCaller(address common.Address, caller bind.ContractCaller) (*LiquidLogicCaller, error) {
	contract, err := bindLiquidLogic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidLogicCaller{contract: contract}, nil
}

// NewLiquidLogicTransactor creates a new write-only instance of LiquidLogic, bound to a specific deployed contract.
func NewLiquidLogicTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidLogicTransactor, error) {
	contract, err := bindLiquidLogic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidLogicTransactor{contract: contract}, nil
}

// NewLiquidLogicFilterer creates a new log filterer instance of LiquidLogic, bound to a specific deployed contract.
func NewLiquidLogicFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidLogicFilterer, error) {
	contract, err := bindLiquidLogic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidLogicFilterer{contract: contract}, nil
}

// bindLiquidLogic binds a generic wrapper to an already deployed contract.
func bindLiquidLogic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LiquidLogicABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidLogic *LiquidLogicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidLogic.Contract.LiquidLogicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidLogic *LiquidLogicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidLogic.Contract.LiquidLogicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidLogic *LiquidLogicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidLogic.Contract.LiquidLogicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidLogic *LiquidLogicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidLogic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidLogic *LiquidLogicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidLogic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidLogic *LiquidLogicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidLogic.Contract.contract.Transact(opts, method, params...)
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_LiquidLogic *LiquidLogicCaller) COMMISSIONRATEPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "COMMISSION_RATE_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_LiquidLogic *LiquidLogicSession) COMMISSIONRATEPRECISION() (*big.Int, error) {
	return _LiquidLogic.Contract.COMMISSIONRATEPRECISION(&_LiquidLogic.CallOpts)
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_LiquidLogic *LiquidLogicCallerSession) COMMISSIONRATEPRECISION() (*big.Int, error) {
	return _LiquidLogic.Contract.COMMISSIONRATEPRECISION(&_LiquidLogic.CallOpts)
}

// FEEFACTORUNITRECIP is a free data retrieval call binding the contract method 0x187cf4d7.
//
// Solidity: function FEE_FACTOR_UNIT_RECIP() view returns(uint256)
func (_LiquidLogic *LiquidLogicCaller) FEEFACTORUNITRECIP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "FEE_FACTOR_UNIT_RECIP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEFACTORUNITRECIP is a free data retrieval call binding the contract method 0x187cf4d7.
//
// Solidity: function FEE_FACTOR_UNIT_RECIP() view returns(uint256)
func (_LiquidLogic *LiquidLogicSession) FEEFACTORUNITRECIP() (*big.Int, error) {
	return _LiquidLogic.Contract.FEEFACTORUNITRECIP(&_LiquidLogic.CallOpts)
}

// FEEFACTORUNITRECIP is a free data retrieval call binding the contract method 0x187cf4d7.
//
// Solidity: function FEE_FACTOR_UNIT_RECIP() view returns(uint256)
func (_LiquidLogic *LiquidLogicCallerSession) FEEFACTORUNITRECIP() (*big.Int, error) {
	return _LiquidLogic.Contract.FEEFACTORUNITRECIP(&_LiquidLogic.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_LiquidLogic *LiquidLogicCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_LiquidLogic *LiquidLogicSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _LiquidLogic.Contract.Allowance(&_LiquidLogic.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_LiquidLogic *LiquidLogicCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _LiquidLogic.Contract.Allowance(&_LiquidLogic.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_LiquidLogic *LiquidLogicCaller) BalanceOf(opts *bind.CallOpts, _delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "balanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_LiquidLogic *LiquidLogicSession) BalanceOf(_delegator common.Address) (*big.Int, error) {
	return _LiquidLogic.Contract.BalanceOf(&_LiquidLogic.CallOpts, _delegator)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_LiquidLogic *LiquidLogicCallerSession) BalanceOf(_delegator common.Address) (*big.Int, error) {
	return _LiquidLogic.Contract.BalanceOf(&_LiquidLogic.CallOpts, _delegator)
}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_LiquidLogic *LiquidLogicCaller) CommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "commissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_LiquidLogic *LiquidLogicSession) CommissionRate() (*big.Int, error) {
	return _LiquidLogic.Contract.CommissionRate(&_LiquidLogic.CallOpts)
}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_LiquidLogic *LiquidLogicCallerSession) CommissionRate() (*big.Int, error) {
	return _LiquidLogic.Contract.CommissionRate(&_LiquidLogic.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LiquidLogic *LiquidLogicCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LiquidLogic *LiquidLogicSession) TotalSupply() (*big.Int, error) {
	return _LiquidLogic.Contract.TotalSupply(&_LiquidLogic.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LiquidLogic *LiquidLogicCallerSession) TotalSupply() (*big.Int, error) {
	return _LiquidLogic.Contract.TotalSupply(&_LiquidLogic.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LiquidLogic *LiquidLogicCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LiquidLogic *LiquidLogicSession) Treasury() (common.Address, error) {
	return _LiquidLogic.Contract.Treasury(&_LiquidLogic.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LiquidLogic *LiquidLogicCallerSession) Treasury() (common.Address, error) {
	return _LiquidLogic.Contract.Treasury(&_LiquidLogic.CallOpts)
}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_LiquidLogic *LiquidLogicCaller) UnclaimedRewards(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "unclaimedRewards", _account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_LiquidLogic *LiquidLogicSession) UnclaimedRewards(_account common.Address) (*big.Int, error) {
	return _LiquidLogic.Contract.UnclaimedRewards(&_LiquidLogic.CallOpts, _account)
}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_LiquidLogic *LiquidLogicCallerSession) UnclaimedRewards(_account common.Address) (*big.Int, error) {
	return _LiquidLogic.Contract.UnclaimedRewards(&_LiquidLogic.CallOpts, _account)
}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_LiquidLogic *LiquidLogicCaller) UnlockedBalanceOf(opts *bind.CallOpts, _delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidLogic.contract.Call(opts, &out, "unlockedBalanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_LiquidLogic *LiquidLogicSession) UnlockedBalanceOf(_delegator common.Address) (*big.Int, error) {
	return _LiquidLogic.Contract.UnlockedBalanceOf(&_LiquidLogic.CallOpts, _delegator)
}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_LiquidLogic *LiquidLogicCallerSession) UnlockedBalanceOf(_delegator common.Address) (*big.Int, error) {
	return _LiquidLogic.Contract.UnlockedBalanceOf(&_LiquidLogic.CallOpts, _delegator)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_LiquidLogic *LiquidLogicTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "approve", _spender, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_LiquidLogic *LiquidLogicSession) Approve(_spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Approve(&_LiquidLogic.TransactOpts, _spender, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_LiquidLogic *LiquidLogicTransactorSession) Approve(_spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Approve(&_LiquidLogic.TransactOpts, _spender, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicTransactor) Burn(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "burn", _account, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicSession) Burn(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Burn(&_LiquidLogic.TransactOpts, _account, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicTransactorSession) Burn(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Burn(&_LiquidLogic.TransactOpts, _account, _amount)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_LiquidLogic *LiquidLogicTransactor) ClaimRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "claimRewards")
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_LiquidLogic *LiquidLogicSession) ClaimRewards() (*types.Transaction, error) {
	return _LiquidLogic.Contract.ClaimRewards(&_LiquidLogic.TransactOpts)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_LiquidLogic *LiquidLogicTransactorSession) ClaimRewards() (*types.Transaction, error) {
	return _LiquidLogic.Contract.ClaimRewards(&_LiquidLogic.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicTransactor) Lock(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "lock", _account, _amount)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicSession) Lock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Lock(&_LiquidLogic.TransactOpts, _account, _amount)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicTransactorSession) Lock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Lock(&_LiquidLogic.TransactOpts, _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicTransactor) Mint(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "mint", _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicSession) Mint(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Mint(&_LiquidLogic.TransactOpts, _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicTransactorSession) Mint(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Mint(&_LiquidLogic.TransactOpts, _account, _amount)
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_LiquidLogic *LiquidLogicTransactor) Redistribute(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "redistribute")
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_LiquidLogic *LiquidLogicSession) Redistribute() (*types.Transaction, error) {
	return _LiquidLogic.Contract.Redistribute(&_LiquidLogic.TransactOpts)
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_LiquidLogic *LiquidLogicTransactorSession) Redistribute() (*types.Transaction, error) {
	return _LiquidLogic.Contract.Redistribute(&_LiquidLogic.TransactOpts)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_LiquidLogic *LiquidLogicTransactor) SetCommissionRate(opts *bind.TransactOpts, _rate *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "setCommissionRate", _rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_LiquidLogic *LiquidLogicSession) SetCommissionRate(_rate *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.SetCommissionRate(&_LiquidLogic.TransactOpts, _rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_LiquidLogic *LiquidLogicTransactorSession) SetCommissionRate(_rate *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.SetCommissionRate(&_LiquidLogic.TransactOpts, _rate)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_LiquidLogic *LiquidLogicTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "transfer", _to, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_LiquidLogic *LiquidLogicSession) Transfer(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Transfer(&_LiquidLogic.TransactOpts, _to, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_LiquidLogic *LiquidLogicTransactorSession) Transfer(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Transfer(&_LiquidLogic.TransactOpts, _to, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_LiquidLogic *LiquidLogicTransactor) TransferFrom(opts *bind.TransactOpts, _sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "transferFrom", _sender, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_LiquidLogic *LiquidLogicSession) TransferFrom(_sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.TransferFrom(&_LiquidLogic.TransactOpts, _sender, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_LiquidLogic *LiquidLogicTransactorSession) TransferFrom(_sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.TransferFrom(&_LiquidLogic.TransactOpts, _sender, _recipient, _amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicTransactor) Unlock(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.contract.Transact(opts, "unlock", _account, _amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicSession) Unlock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Unlock(&_LiquidLogic.TransactOpts, _account, _amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_LiquidLogic *LiquidLogicTransactorSession) Unlock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidLogic.Contract.Unlock(&_LiquidLogic.TransactOpts, _account, _amount)
}

// LiquidLogicApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the LiquidLogic contract.
type LiquidLogicApprovalIterator struct {
	Event *LiquidLogicApproval // Event containing the contract specifics and raw log

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
func (it *LiquidLogicApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidLogicApproval)
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
		it.Event = new(LiquidLogicApproval)
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
func (it *LiquidLogicApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidLogicApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidLogicApproval represents a Approval event raised by the LiquidLogic contract.
type LiquidLogicApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LiquidLogic *LiquidLogicFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LiquidLogicApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LiquidLogic.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LiquidLogicApprovalIterator{contract: _LiquidLogic.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LiquidLogic *LiquidLogicFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LiquidLogicApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LiquidLogic.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidLogicApproval)
				if err := _LiquidLogic.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_LiquidLogic *LiquidLogicFilterer) ParseApproval(log types.Log) (*LiquidLogicApproval, error) {
	event := new(LiquidLogicApproval)
	if err := _LiquidLogic.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidLogicTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the LiquidLogic contract.
type LiquidLogicTransferIterator struct {
	Event *LiquidLogicTransfer // Event containing the contract specifics and raw log

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
func (it *LiquidLogicTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidLogicTransfer)
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
		it.Event = new(LiquidLogicTransfer)
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
func (it *LiquidLogicTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidLogicTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidLogicTransfer represents a Transfer event raised by the LiquidLogic contract.
type LiquidLogicTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LiquidLogic *LiquidLogicFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidLogicTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LiquidLogic.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LiquidLogicTransferIterator{contract: _LiquidLogic.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LiquidLogic *LiquidLogicFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LiquidLogicTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LiquidLogic.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidLogicTransfer)
				if err := _LiquidLogic.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_LiquidLogic *LiquidLogicFilterer) ParseTransfer(log types.Log) (*LiquidLogicTransfer, error) {
	event := new(LiquidLogicTransfer)
	if err := _LiquidLogic.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStateMetaData contains all meta data concerning the LiquidState contract.
var LiquidStateMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"_treasury\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_index\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"liquidLogic\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"lockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"redistribute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unclaimedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"unlockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"9dc29fac": "burn(address,uint256)",
		"372500ab": "claimRewards()",
		"5ea1d6f8": "commissionRate()",
		"313ce567": "decimals()",
		"0c3c6fee": "liquidLogic()",
		"282d3fdf": "lock(address,uint256)",
		"59355736": "lockedBalanceOf(address)",
		"40c10f19": "mint(address,uint256)",
		"06fdde03": "name()",
		"fb489a7b": "redistribute()",
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
	Bin: "0x608060405234801562000010575f80fd5b5060405162001954380380620019548339810160408190526200003391620001bc565b5f856001600160a01b0316632f2c3f2e6040518163ffffffff1660e01b8152600401602060405180830381865afa15801562000071573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620000979190620002ba565b905080831115620000a6575f80fd5b600d80546001600160a01b038089166001600160a01b031992831617909255600c805488841690831617905560088054928716929091169190911790556009839055604051620000fb908390602001620002d2565b604051602081830303815290604052600a90816200011a91906200038a565b50816040516020016200012e9190620002d2565b604051602081830303815290604052600b90816200014d91906200038a565b50505f80546001600160a01b0319163317905550620004569350505050565b6001600160a01b038116811462000181575f80fd5b50565b634e487b7160e01b5f52604160045260245ffd5b5f5b83811015620001b45781810151838201526020016200019a565b50505f910152565b5f805f805f60a08688031215620001d1575f80fd5b8551620001de816200016c565b6020870151909550620001f1816200016c565b604087015190945062000204816200016c565b6060870151608088015191945092506001600160401b038082111562000228575f80fd5b818801915088601f8301126200023c575f80fd5b81518181111562000251576200025162000184565b604051601f8201601f19908116603f011681019083821181831017156200027c576200027c62000184565b816040528281528b602084870101111562000295575f80fd5b620002a883602083016020880162000198565b80955050505050509295509295909350565b5f60208284031215620002cb575f80fd5b5051919050565b644c4e544e2d60d81b81525f8251620002f381600585016020870162000198565b9190910160050192915050565b600181811c908216806200031557607f821691505b6020821081036200033457634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200038557805f5260205f20601f840160051c81016020851015620003615750805b601f840160051c820191505b8181101562000382575f81556001016200036d565b50505b505050565b81516001600160401b03811115620003a657620003a662000184565b620003be81620003b7845462000300565b846200033a565b602080601f831160018114620003f4575f8415620003dc5750858301515b5f19600386901b1c1916600185901b1785556200044e565b5f85815260208120601f198616915b82811015620004245788860151825594840194600190910190840162000403565b50858210156200044257878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b6114f080620004645f395ff3fe60806040526004361061013c575f3560e01c806359355736116100b3578063949813b81161006d578063949813b81461039457806395d89b41146103b35780639dc29fac146103c7578063a9059cbb146103e6578063dd62ed3e14610405578063fb489a7b14610449575f80fd5b806359355736146102ba5780635ea1d6f8146102ee57806361d027b31461030357806370a08231146103225780637eee288d1461035657806384955c8814610375575f80fd5b806323b872dd1161010457806323b872dd1461020f578063282d3fdf1461022e578063313ce5671461024d578063372500ab146102685780633a5381b51461027c57806340c10f191461029b575f80fd5b806306fdde0314610140578063095ea7b31461016a5780630c3c6fee1461019957806318160ddd146101d057806319fac8fd146101ee575b5f80fd5b34801561014b575f80fd5b50610154610451565b60405161016191906112d3565b60405180910390f35b348015610175575f80fd5b50610189610184366004611320565b6104e1565b6040519015158152602001610161565b3480156101a4575f80fd5b50600d546101b8906001600160a01b031681565b6040516001600160a01b039091168152602001610161565b3480156101db575f80fd5b506004545b604051908152602001610161565b3480156101f9575f80fd5b5061020d610208366004611348565b610606565b005b34801561021a575f80fd5b5061018961022936600461135f565b610740565b348015610239575f80fd5b5061020d610248366004611320565b610877565b348015610258575f80fd5b5060405160128152602001610161565b348015610273575f80fd5b5061020d6109a6565b348015610287575f80fd5b50600c546101b8906001600160a01b031681565b3480156102a6575f80fd5b5061020d6102b5366004611320565b610aa3565b3480156102c5575f80fd5b506101e06102d4366004611398565b6001600160a01b03165f9081526002602052604090205490565b3480156102f9575f80fd5b506101e060095481565b34801561030e575f80fd5b506008546101b8906001600160a01b031681565b34801561032d575f80fd5b506101e061033c366004611398565b6001600160a01b03165f9081526001602052604090205490565b348015610361575f80fd5b5061020d610370366004611320565b610bcc565b348015610380575f80fd5b506101e061038f366004611398565b610cf5565b34801561039f575f80fd5b506101e06103ae366004611398565b610e20565b3480156103be575f80fd5b50610154610f2e565b3480156103d2575f80fd5b5061020d6103e1366004611320565b610f3d565b3480156103f1575f80fd5b50610189610400366004611320565b611066565b348015610410575f80fd5b506101e061041f3660046113b8565b6001600160a01b039182165f90815260036020908152604080832093909416825291909152205490565b6101e0611173565b6060600a8054610460906113e9565b80601f016020809104026020016040519081016040528092919081815260200182805461048c906113e9565b80156104d75780601f106104ae576101008083540402835291602001916104d7565b820191905f5260205f20905b8154815290600101906020018083116104ba57829003601f168201915b5050505050905090565b600d546040516001600160a01b038481166024830152604482018490525f92839283929091169060640160408051601f198184030181529181526020820180516001600160e01b031663095ea7b360e01b179052516105409190611421565b5f60405180830381855af49150503d805f8114610578576040519150601f19603f3d011682016040523d82523d5f602084013e61057d565b606091505b5090925090508115155f036105e95780511561059c5780518082602001fd5b60405162461bcd60e51b815260206004820181905260248201527f46756e6374696f6e2063616c6c2027617070726f76652720726576657274656460448201526064015b60405180910390fd5b808060200190518101906105fd919061143c565b95945050505050565b5f546001600160a01b0316331461062f5760405162461bcd60e51b81526004016105e09061145b565b600d54604051602481018390525f9182916001600160a01b039091169060440160408051601f198184030181529181526020820180516001600160e01b03166319fac8fd60e01b179052516106849190611421565b5f60405180830381855af49150503d805f81146106bc576040519150601f19603f3d011682016040523d82523d5f602084013e6106c1565b606091505b5090925090508115155f0361073b578051156106e05780518082602001fd5b60405162461bcd60e51b815260206004820152602a60248201527f46756e6374696f6e2063616c6c2027736574436f6d6d697373696f6e5261746560448201526909c81c995d995c9d195960b21b60648201526084016105e0565b505050565b600d546040516001600160a01b0385811660248301528481166044830152606482018490525f92839283929091169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516107a79190611421565b5f60405180830381855af49150503d805f81146107df576040519150601f19603f3d011682016040523d82523d5f602084013e6107e4565b606091505b5090925090508115155f03610859578051156108035780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c20277472616e7366657246726f6d2720726576604482015264195c9d195960da1b60648201526084016105e0565b8080602001905181019061086d919061143c565b9695505050505050565b5f546001600160a01b031633146108a05760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b031663282d3fdf60e01b179052516108fc9190611421565b5f60405180830381855af49150503d805f8114610934576040519150601f19603f3d011682016040523d82523d5f602084013e610939565b606091505b5090925090508115155f036109a0578051156109585780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276c6f636b2720726576657274656400000060448201526064016105e0565b50505050565b600d5460408051600481526024810182526020810180516001600160e01b031663372500ab60e01b17905290515f9283926001600160a01b03909116916109ed9190611421565b5f60405180830381855af49150503d805f8114610a25576040519150601f19603f3d011682016040523d82523d5f602084013e610a2a565b606091505b5090925090508115155f03610a9f57805115610a495780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c2027636c61696d526577617264732720726576604482015264195c9d195960da1b60648201526084016105e0565b5050565b5f546001600160a01b03163314610acc5760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b03166340c10f1960e01b17905251610b289190611421565b5f60405180830381855af49150503d805f8114610b60576040519150601f19603f3d011682016040523d82523d5f602084013e610b65565b606091505b5090925090508115155f036109a057805115610b845780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276d696e742720726576657274656400000060448201526064016105e0565b5f546001600160a01b03163314610bf55760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b0316637eee288d60e01b17905251610c519190611421565b5f60405180830381855af49150503d805f8114610c89576040519150601f19603f3d011682016040523d82523d5f602084013e610c8e565b606091505b5090925090508115155f036109a057805115610cad5780518082602001fd5b60405162461bcd60e51b815260206004820152601f60248201527f46756e6374696f6e2063616c6c2027756e6c6f636b272072657665727465640060448201526064016105e0565b600d546040516001600160a01b0383811660248301525f92839283929091169060440160408051601f198184030181529181526020820180516001600160e01b0316631092ab9160e31b17905251610d4d9190611421565b5f60405180830381855af49150503d805f8114610d85576040519150601f19603f3d011682016040523d82523d5f602084013e610d8a565b606091505b5090925090508115155f03610e0457805115610da95780518082602001fd5b60405162461bcd60e51b815260206004820152602a60248201527f46756e6374696f6e2063616c6c2027756e6c6f636b656442616c616e63654f6660448201526909c81c995d995c9d195960b21b60648201526084016105e0565b80806020019051810190610e1891906114a3565b949350505050565b600d546040516001600160a01b0383811660248301525f92839283929091169060440160408051601f198184030181529181526020820180516001600160e01b0316631293027760e31b17905251610e789190611421565b5f60405180830381855af49150503d805f8114610eb0576040519150601f19603f3d011682016040523d82523d5f602084013e610eb5565b606091505b5090925090508115155f03610e0457805115610ed45780518082602001fd5b60405162461bcd60e51b815260206004820152602960248201527f46756e6374696f6e2063616c6c2027756e636c61696d65645265776172647327604482015268081c995d995c9d195960ba1b60648201526084016105e0565b6060600b8054610460906113e9565b5f546001600160a01b03163314610f665760405162461bcd60e51b81526004016105e09061145b565b600d546040516001600160a01b038481166024830152604482018490525f92839291169060640160408051601f198184030181529181526020820180516001600160e01b0316632770a7eb60e21b17905251610fc29190611421565b5f60405180830381855af49150503d805f8114610ffa576040519150601f19603f3d011682016040523d82523d5f602084013e610fff565b606091505b5090925090508115155f036109a05780511561101e5780518082602001fd5b60405162461bcd60e51b815260206004820152601d60248201527f46756e6374696f6e2063616c6c20276275726e2720726576657274656400000060448201526064016105e0565b600d546040516001600160a01b038481166024830152604482018490525f92839283929091169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b179052516110c59190611421565b5f60405180830381855af49150503d805f81146110fd576040519150601f19603f3d011682016040523d82523d5f602084013e611102565b606091505b5090925090508115155f036105e9578051156111215780518082602001fd5b60405162461bcd60e51b815260206004820152602160248201527f46756e6374696f6e2063616c6c20277472616e736665722720726576657274656044820152601960fa1b60648201526084016105e0565b5f80546001600160a01b0316331461119d5760405162461bcd60e51b81526004016105e09061145b565b600d5460408051600481526024810182526020810180516001600160e01b031663fb489a7b60e01b17905290515f9283926001600160a01b03909116916111e49190611421565b5f60405180830381855af49150503d805f811461121c576040519150601f19603f3d011682016040523d82523d5f602084013e611221565b606091505b5090925090508115155f03611296578051156112405780518082602001fd5b60405162461bcd60e51b815260206004820152602560248201527f46756e6374696f6e2063616c6c20277265646973747269627574652720726576604482015264195c9d195960da1b60648201526084016105e0565b808060200190518101906112aa91906114a3565b9250505090565b5f5b838110156112cb5781810151838201526020016112b3565b50505f910152565b602081525f82518060208401526112f18160408501602087016112b1565b601f01601f19169190910160400192915050565b80356001600160a01b038116811461131b575f80fd5b919050565b5f8060408385031215611331575f80fd5b61133a83611305565b946020939093013593505050565b5f60208284031215611358575f80fd5b5035919050565b5f805f60608486031215611371575f80fd5b61137a84611305565b925061138860208501611305565b9150604084013590509250925092565b5f602082840312156113a8575f80fd5b6113b182611305565b9392505050565b5f80604083850312156113c9575f80fd5b6113d283611305565b91506113e060208401611305565b90509250929050565b600181811c908216806113fd57607f821691505b60208210810361141b57634e487b7160e01b5f52602260045260245ffd5b50919050565b5f82516114328184602087016112b1565b9190910192915050565b5f6020828403121561144c575f80fd5b815180151581146113b1575f80fd5b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b5f602082840312156114b3575f80fd5b505191905056fea26469706673582212200f950bd4f4d9457214149c5d71464d2a2df907f4dd143e5feaf92aaf3d4852a864736f6c63430008170033",
}

// LiquidStateABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidStateMetaData.ABI instead.
var LiquidStateABI = LiquidStateMetaData.ABI

// Deprecated: Use LiquidStateMetaData.Sigs instead.
// LiquidStateFuncSigs maps the 4-byte function signature to its string representation.
var LiquidStateFuncSigs = LiquidStateMetaData.Sigs

// LiquidStateBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LiquidStateMetaData.Bin instead.
var LiquidStateBin = LiquidStateMetaData.Bin

// DeployLiquidState deploys a new Ethereum contract, binding an instance of LiquidState to it.
func DeployLiquidState(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, _validator common.Address, _treasury common.Address, _commissionRate *big.Int, _index string) (common.Address, *types.Transaction, *LiquidState, error) {
	parsed, err := LiquidStateMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidStateBin), backend, _logic, _validator, _treasury, _commissionRate, _index)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LiquidState{LiquidStateCaller: LiquidStateCaller{contract: contract}, LiquidStateTransactor: LiquidStateTransactor{contract: contract}, LiquidStateFilterer: LiquidStateFilterer{contract: contract}}, nil
}

// LiquidState is an auto generated Go binding around an Ethereum contract.
type LiquidState struct {
	LiquidStateCaller     // Read-only binding to the contract
	LiquidStateTransactor // Write-only binding to the contract
	LiquidStateFilterer   // Log filterer for contract events
}

// LiquidStateCaller is an auto generated read-only Go binding around an Ethereum contract.
type LiquidStateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidStateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidStateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidStateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidStateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidStateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LiquidStateSession struct {
	Contract     *LiquidState      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LiquidStateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LiquidStateCallerSession struct {
	Contract *LiquidStateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// LiquidStateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LiquidStateTransactorSession struct {
	Contract     *LiquidStateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// LiquidStateRaw is an auto generated low-level Go binding around an Ethereum contract.
type LiquidStateRaw struct {
	Contract *LiquidState // Generic contract binding to access the raw methods on
}

// LiquidStateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LiquidStateCallerRaw struct {
	Contract *LiquidStateCaller // Generic read-only contract binding to access the raw methods on
}

// LiquidStateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LiquidStateTransactorRaw struct {
	Contract *LiquidStateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLiquidState creates a new instance of LiquidState, bound to a specific deployed contract.
func NewLiquidState(address common.Address, backend bind.ContractBackend) (*LiquidState, error) {
	contract, err := bindLiquidState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LiquidState{LiquidStateCaller: LiquidStateCaller{contract: contract}, LiquidStateTransactor: LiquidStateTransactor{contract: contract}, LiquidStateFilterer: LiquidStateFilterer{contract: contract}}, nil
}

// NewLiquidStateCaller creates a new read-only instance of LiquidState, bound to a specific deployed contract.
func NewLiquidStateCaller(address common.Address, caller bind.ContractCaller) (*LiquidStateCaller, error) {
	contract, err := bindLiquidState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidStateCaller{contract: contract}, nil
}

// NewLiquidStateTransactor creates a new write-only instance of LiquidState, bound to a specific deployed contract.
func NewLiquidStateTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidStateTransactor, error) {
	contract, err := bindLiquidState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidStateTransactor{contract: contract}, nil
}

// NewLiquidStateFilterer creates a new log filterer instance of LiquidState, bound to a specific deployed contract.
func NewLiquidStateFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidStateFilterer, error) {
	contract, err := bindLiquidState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidStateFilterer{contract: contract}, nil
}

// bindLiquidState binds a generic wrapper to an already deployed contract.
func bindLiquidState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LiquidStateABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidState *LiquidStateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidState.Contract.LiquidStateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidState *LiquidStateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidState.Contract.LiquidStateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidState *LiquidStateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidState.Contract.LiquidStateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LiquidState *LiquidStateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LiquidState.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LiquidState *LiquidStateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidState.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LiquidState *LiquidStateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LiquidState.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_LiquidState *LiquidStateCaller) Allowance(opts *bind.CallOpts, _owner common.Address, _spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "allowance", _owner, _spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_LiquidState *LiquidStateSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _LiquidState.Contract.Allowance(&_LiquidState.CallOpts, _owner, _spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address _owner, address _spender) view returns(uint256)
func (_LiquidState *LiquidStateCallerSession) Allowance(_owner common.Address, _spender common.Address) (*big.Int, error) {
	return _LiquidState.Contract.Allowance(&_LiquidState.CallOpts, _owner, _spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_LiquidState *LiquidStateCaller) BalanceOf(opts *bind.CallOpts, _delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "balanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_LiquidState *LiquidStateSession) BalanceOf(_delegator common.Address) (*big.Int, error) {
	return _LiquidState.Contract.BalanceOf(&_LiquidState.CallOpts, _delegator)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _delegator) view returns(uint256)
func (_LiquidState *LiquidStateCallerSession) BalanceOf(_delegator common.Address) (*big.Int, error) {
	return _LiquidState.Contract.BalanceOf(&_LiquidState.CallOpts, _delegator)
}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_LiquidState *LiquidStateCaller) CommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "commissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_LiquidState *LiquidStateSession) CommissionRate() (*big.Int, error) {
	return _LiquidState.Contract.CommissionRate(&_LiquidState.CallOpts)
}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_LiquidState *LiquidStateCallerSession) CommissionRate() (*big.Int, error) {
	return _LiquidState.Contract.CommissionRate(&_LiquidState.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_LiquidState *LiquidStateCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_LiquidState *LiquidStateSession) Decimals() (uint8, error) {
	return _LiquidState.Contract.Decimals(&_LiquidState.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_LiquidState *LiquidStateCallerSession) Decimals() (uint8, error) {
	return _LiquidState.Contract.Decimals(&_LiquidState.CallOpts)
}

// LiquidLogic is a free data retrieval call binding the contract method 0x0c3c6fee.
//
// Solidity: function liquidLogic() view returns(address)
func (_LiquidState *LiquidStateCaller) LiquidLogic(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "liquidLogic")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LiquidLogic is a free data retrieval call binding the contract method 0x0c3c6fee.
//
// Solidity: function liquidLogic() view returns(address)
func (_LiquidState *LiquidStateSession) LiquidLogic() (common.Address, error) {
	return _LiquidState.Contract.LiquidLogic(&_LiquidState.CallOpts)
}

// LiquidLogic is a free data retrieval call binding the contract method 0x0c3c6fee.
//
// Solidity: function liquidLogic() view returns(address)
func (_LiquidState *LiquidStateCallerSession) LiquidLogic() (common.Address, error) {
	return _LiquidState.Contract.LiquidLogic(&_LiquidState.CallOpts)
}

// LockedBalanceOf is a free data retrieval call binding the contract method 0x59355736.
//
// Solidity: function lockedBalanceOf(address _delegator) view returns(uint256)
func (_LiquidState *LiquidStateCaller) LockedBalanceOf(opts *bind.CallOpts, _delegator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "lockedBalanceOf", _delegator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedBalanceOf is a free data retrieval call binding the contract method 0x59355736.
//
// Solidity: function lockedBalanceOf(address _delegator) view returns(uint256)
func (_LiquidState *LiquidStateSession) LockedBalanceOf(_delegator common.Address) (*big.Int, error) {
	return _LiquidState.Contract.LockedBalanceOf(&_LiquidState.CallOpts, _delegator)
}

// LockedBalanceOf is a free data retrieval call binding the contract method 0x59355736.
//
// Solidity: function lockedBalanceOf(address _delegator) view returns(uint256)
func (_LiquidState *LiquidStateCallerSession) LockedBalanceOf(_delegator common.Address) (*big.Int, error) {
	return _LiquidState.Contract.LockedBalanceOf(&_LiquidState.CallOpts, _delegator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LiquidState *LiquidStateCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LiquidState *LiquidStateSession) Name() (string, error) {
	return _LiquidState.Contract.Name(&_LiquidState.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LiquidState *LiquidStateCallerSession) Name() (string, error) {
	return _LiquidState.Contract.Name(&_LiquidState.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LiquidState *LiquidStateCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LiquidState *LiquidStateSession) Symbol() (string, error) {
	return _LiquidState.Contract.Symbol(&_LiquidState.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LiquidState *LiquidStateCallerSession) Symbol() (string, error) {
	return _LiquidState.Contract.Symbol(&_LiquidState.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LiquidState *LiquidStateCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LiquidState *LiquidStateSession) TotalSupply() (*big.Int, error) {
	return _LiquidState.Contract.TotalSupply(&_LiquidState.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LiquidState *LiquidStateCallerSession) TotalSupply() (*big.Int, error) {
	return _LiquidState.Contract.TotalSupply(&_LiquidState.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LiquidState *LiquidStateCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LiquidState *LiquidStateSession) Treasury() (common.Address, error) {
	return _LiquidState.Contract.Treasury(&_LiquidState.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LiquidState *LiquidStateCallerSession) Treasury() (common.Address, error) {
	return _LiquidState.Contract.Treasury(&_LiquidState.CallOpts)
}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_LiquidState *LiquidStateCaller) Validator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LiquidState.contract.Call(opts, &out, "validator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_LiquidState *LiquidStateSession) Validator() (common.Address, error) {
	return _LiquidState.Contract.Validator(&_LiquidState.CallOpts)
}

// Validator is a free data retrieval call binding the contract method 0x3a5381b5.
//
// Solidity: function validator() view returns(address)
func (_LiquidState *LiquidStateCallerSession) Validator() (common.Address, error) {
	return _LiquidState.Contract.Validator(&_LiquidState.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_LiquidState *LiquidStateTransactor) Approve(opts *bind.TransactOpts, _spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "approve", _spender, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_LiquidState *LiquidStateSession) Approve(_spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Approve(&_LiquidState.TransactOpts, _spender, _amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _spender, uint256 _amount) returns(bool)
func (_LiquidState *LiquidStateTransactorSession) Approve(_spender common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Approve(&_LiquidState.TransactOpts, _spender, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateTransactor) Burn(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "burn", _account, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateSession) Burn(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Burn(&_LiquidState.TransactOpts, _account, _amount)
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateTransactorSession) Burn(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Burn(&_LiquidState.TransactOpts, _account, _amount)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_LiquidState *LiquidStateTransactor) ClaimRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "claimRewards")
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_LiquidState *LiquidStateSession) ClaimRewards() (*types.Transaction, error) {
	return _LiquidState.Contract.ClaimRewards(&_LiquidState.TransactOpts)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_LiquidState *LiquidStateTransactorSession) ClaimRewards() (*types.Transaction, error) {
	return _LiquidState.Contract.ClaimRewards(&_LiquidState.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateTransactor) Lock(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "lock", _account, _amount)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateSession) Lock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Lock(&_LiquidState.TransactOpts, _account, _amount)
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateTransactorSession) Lock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Lock(&_LiquidState.TransactOpts, _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateTransactor) Mint(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "mint", _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateSession) Mint(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Mint(&_LiquidState.TransactOpts, _account, _amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateTransactorSession) Mint(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Mint(&_LiquidState.TransactOpts, _account, _amount)
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_LiquidState *LiquidStateTransactor) Redistribute(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "redistribute")
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_LiquidState *LiquidStateSession) Redistribute() (*types.Transaction, error) {
	return _LiquidState.Contract.Redistribute(&_LiquidState.TransactOpts)
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_LiquidState *LiquidStateTransactorSession) Redistribute() (*types.Transaction, error) {
	return _LiquidState.Contract.Redistribute(&_LiquidState.TransactOpts)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_LiquidState *LiquidStateTransactor) SetCommissionRate(opts *bind.TransactOpts, _rate *big.Int) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "setCommissionRate", _rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_LiquidState *LiquidStateSession) SetCommissionRate(_rate *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.SetCommissionRate(&_LiquidState.TransactOpts, _rate)
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_LiquidState *LiquidStateTransactorSession) SetCommissionRate(_rate *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.SetCommissionRate(&_LiquidState.TransactOpts, _rate)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_LiquidState *LiquidStateTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "transfer", _to, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_LiquidState *LiquidStateSession) Transfer(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Transfer(&_LiquidState.TransactOpts, _to, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_LiquidState *LiquidStateTransactorSession) Transfer(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Transfer(&_LiquidState.TransactOpts, _to, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_LiquidState *LiquidStateTransactor) TransferFrom(opts *bind.TransactOpts, _sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "transferFrom", _sender, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_LiquidState *LiquidStateSession) TransferFrom(_sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.TransferFrom(&_LiquidState.TransactOpts, _sender, _recipient, _amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_LiquidState *LiquidStateTransactorSession) TransferFrom(_sender common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.TransferFrom(&_LiquidState.TransactOpts, _sender, _recipient, _amount)
}

// UnclaimedRewards is a paid mutator transaction binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) returns(uint256)
func (_LiquidState *LiquidStateTransactor) UnclaimedRewards(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "unclaimedRewards", _account)
}

// UnclaimedRewards is a paid mutator transaction binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) returns(uint256)
func (_LiquidState *LiquidStateSession) UnclaimedRewards(_account common.Address) (*types.Transaction, error) {
	return _LiquidState.Contract.UnclaimedRewards(&_LiquidState.TransactOpts, _account)
}

// UnclaimedRewards is a paid mutator transaction binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) returns(uint256)
func (_LiquidState *LiquidStateTransactorSession) UnclaimedRewards(_account common.Address) (*types.Transaction, error) {
	return _LiquidState.Contract.UnclaimedRewards(&_LiquidState.TransactOpts, _account)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateTransactor) Unlock(opts *bind.TransactOpts, _account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "unlock", _account, _amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateSession) Unlock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Unlock(&_LiquidState.TransactOpts, _account, _amount)
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_LiquidState *LiquidStateTransactorSession) Unlock(_account common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _LiquidState.Contract.Unlock(&_LiquidState.TransactOpts, _account, _amount)
}

// UnlockedBalanceOf is a paid mutator transaction binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) returns(uint256)
func (_LiquidState *LiquidStateTransactor) UnlockedBalanceOf(opts *bind.TransactOpts, _delegator common.Address) (*types.Transaction, error) {
	return _LiquidState.contract.Transact(opts, "unlockedBalanceOf", _delegator)
}

// UnlockedBalanceOf is a paid mutator transaction binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) returns(uint256)
func (_LiquidState *LiquidStateSession) UnlockedBalanceOf(_delegator common.Address) (*types.Transaction, error) {
	return _LiquidState.Contract.UnlockedBalanceOf(&_LiquidState.TransactOpts, _delegator)
}

// UnlockedBalanceOf is a paid mutator transaction binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) returns(uint256)
func (_LiquidState *LiquidStateTransactorSession) UnlockedBalanceOf(_delegator common.Address) (*types.Transaction, error) {
	return _LiquidState.Contract.UnlockedBalanceOf(&_LiquidState.TransactOpts, _delegator)
}

// LiquidStateApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the LiquidState contract.
type LiquidStateApprovalIterator struct {
	Event *LiquidStateApproval // Event containing the contract specifics and raw log

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
func (it *LiquidStateApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStateApproval)
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
		it.Event = new(LiquidStateApproval)
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
func (it *LiquidStateApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStateApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStateApproval represents a Approval event raised by the LiquidState contract.
type LiquidStateApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LiquidState *LiquidStateFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LiquidStateApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LiquidState.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStateApprovalIterator{contract: _LiquidState.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_LiquidState *LiquidStateFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LiquidStateApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _LiquidState.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStateApproval)
				if err := _LiquidState.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_LiquidState *LiquidStateFilterer) ParseApproval(log types.Log) (*LiquidStateApproval, error) {
	event := new(LiquidStateApproval)
	if err := _LiquidState.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidStateTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the LiquidState contract.
type LiquidStateTransferIterator struct {
	Event *LiquidStateTransfer // Event containing the contract specifics and raw log

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
func (it *LiquidStateTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidStateTransfer)
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
		it.Event = new(LiquidStateTransfer)
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
func (it *LiquidStateTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidStateTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidStateTransfer represents a Transfer event raised by the LiquidState contract.
type LiquidStateTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LiquidState *LiquidStateFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidStateTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LiquidState.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LiquidStateTransferIterator{contract: _LiquidState.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_LiquidState *LiquidStateFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LiquidStateTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LiquidState.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidStateTransfer)
				if err := _LiquidState.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_LiquidState *LiquidStateFilterer) ParseTransfer(log types.Log) (*LiquidStateTransfer, error) {
	event := new(LiquidStateTransfer)
	if err := _LiquidState.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

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
	Bin: "0x6080604052600160ff1b600755600160ff1b60085534801562000020575f80fd5b5060405162002c8d38038062002c8d83398101604081905262000043916200060f565b600280546001600160a01b038087166001600160a01b0319928316179092556003805492861692909116919091179055815162000087905f9060208501906200034c565b5081516200009d9060019060208501906200034c565b5080600981905550620000c1855f60018851620000bb91906200071a565b6200016f565b8451620000d6906004906020880190620003a7565b508451620000ec906005906020880190620003a7565b5060016006819055600d805490910181555f9081525b855181101562000163576001600b5f88848151811062000126576200012662000736565b6020908102919091018101516001600160a01b031682528101919091526040015f20600201805460ff191691151591909117905560010162000102565b5050505050506200096a565b8082126200017c57505050565b81815f8560026200018e85856200074a565b6200019a919062000773565b620001a69087620007af565b81518110620001b957620001b962000736565b602002602001015190505b81831362000318575b806001600160a01b0316868481518110620001ec57620001ec62000736565b60200260200101516001600160a01b031610156200021957826200021081620007d9565b935050620001cd565b806001600160a01b031686838151811062000238576200023862000736565b60200260200101516001600160a01b031611156200026557816200025c81620007fa565b92505062000219565b818313620003125785828151811062000282576200028262000736565b60200260200101518684815181106200029f576200029f62000736565b6020026020010151878581518110620002bc57620002bc62000736565b60200260200101888581518110620002d857620002d862000736565b6001600160a01b0393841660209182029290920101529116905282620002fe81620007d9565b93505081806200030e90620007fa565b9250505b620001c4565b818512156200032e576200032e8686846200016f565b838312156200034457620003448684866200016f565b505050505050565b828054828255905f5260205f2090810192821562000395579160200282015b82811115620003955782518290620003849082620008a2565b50916020019190600101906200036b565b50620003a39291506200040b565b5090565b828054828255905f5260205f20908101928215620003fd579160200282015b82811115620003fd57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620003c6565b50620003a39291506200042b565b80821115620003a3575f62000421828262000441565b506001016200040b565b5b80821115620003a3575f81556001016200042c565b5080546200044f9062000818565b5f825580601f106200045f575050565b601f0160209004905f5260205f20908101906200047d91906200042b565b50565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f191681016001600160401b0381118282101715620004bf57620004bf62000480565b604052919050565b5f6001600160401b03821115620004e257620004e262000480565b5060051b60200190565b80516001600160a01b038116811462000503575f80fd5b919050565b5f601f83601f8401126200051a575f80fd5b82516020620005336200052d83620004c7565b62000494565b82815260059290921b8501810191818101908784111562000552575f80fd5b8287015b84811015620006035780516001600160401b038082111562000576575f80fd5b818a0191508a603f8301126200058a575f80fd5b8582015181811115620005a157620005a162000480565b620005b4818a01601f1916880162000494565b915080825260408c81838601011115620005cc575f80fd5b5f5b82811015620005eb578481018201518482018a01528801620005ce565b50505f90820187015284525091830191830162000556565b50979650505050505050565b5f805f805f60a0868803121562000624575f80fd5b85516001600160401b03808211156200063b575f80fd5b818801915088601f8301126200064f575f80fd5b81516020620006626200052d83620004c7565b82815260059290921b8401810191818101908c84111562000681575f80fd5b948201945b83861015620006aa576200069a86620004ec565b8252948201949082019062000686565b9950620006bb90508a8201620004ec565b97505050620006cd60408901620004ec565b94506060880151915080821115620006e3575f80fd5b50620006f28882890162000508565b925050608086015190509295509295909350565b634e487b7160e01b5f52601160045260245ffd5b8181038181111562000730576200073062000706565b92915050565b634e487b7160e01b5f52603260045260245ffd5b8181035f8312801583831316838312821617156200076c576200076c62000706565b5092915050565b5f826200078e57634e487b7160e01b5f52601260045260245ffd5b600160ff1b82145f1984141615620007aa57620007aa62000706565b500590565b8082018281125f831280158216821582161715620007d157620007d162000706565b505092915050565b5f6001600160ff1b018201620007f357620007f362000706565b5060010190565b5f600160ff1b820162000811576200081162000706565b505f190190565b600181811c908216806200082d57607f821691505b6020821081036200084c57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156200089d57805f5260205f20601f840160051c81016020851015620008795750805b601f840160051c820191505b818110156200089a575f815560010162000885565b50505b505050565b81516001600160401b03811115620008be57620008be62000480565b620008d681620008cf845462000818565b8462000852565b602080601f8311600181146200090c575f8415620008f45750858301515b5f19600386901b1c1916600185901b17855562000344565b5f85815260208120601f198616915b828110156200093c578886015182559484019460019091019084016200091b565b50858210156200095a57878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b61231580620009785f395ff3fe608060405260043610610121575f3560e01c80638d4f75d2116100a4578063b3ab15fb1161006b578063b3ab15fb14610378578063b78dec5214610397578063ccce413b146103ab578063cdd72253146103ca578063df7f710e146103eb578063e6a02a281461040c57005b80638d4f75d2146103065780639670c0bc146103255780639f8743f71461033a578063a78135871461034e578063aa2f89b51461036357005b80634bb278f3116100e85780634bb278f3146101f75780634c56ea561461021b5780635281b5c6146102615780635412b3ae1461028d578063845023f2146102e757005b806308f21ff51461012a578063146ca53114610152578063307de9b61461016757806333f98c77146101865780633c8510fd146101d857005b3661012857005b005b348015610135575f80fd5b5061013f60085481565b6040519081526020015b60405180910390f35b34801561015d575f80fd5b5061013f60065481565b348015610172575f80fd5b5061012861018136600461197d565b610421565b348015610191575f80fd5b506101a56101a0366004611aac565b610644565b60405161014991908151815260208083015190820152604080830151908201526060918201519181019190915260800190565b3480156101e3575f80fd5b506101a56101f2366004611ade565b61075c565b348015610202575f80fd5b5061020b61085b565b6040519015158152602001610149565b348015610226575f80fd5b5061013f610235366004611b3d565b8151602081840181018051600c825292820194820194909420919093529091525f908152604090205481565b34801561026c575f80fd5b5061028061027b366004611b88565b6109eb565b6040516101499190611bec565b348015610298575f80fd5b506102ca6102a7366004611c05565b600b6020525f908152604090208054600182015460029092015490919060ff1683565b604080519384526020840192909252151590820152606001610149565b3480156102f2575f80fd5b50610128610301366004611c41565b610a91565b348015610311575f80fd5b50610128610320366004611cd9565b610b38565b348015610330575f80fd5b506298968061013f565b348015610345575f80fd5b5060065461013f565b348015610359575f80fd5b5061013f60095481565b34801561036e575f80fd5b5061013f60075481565b348015610383575f80fd5b50610128610392366004611c05565b610ca7565b3480156103a2575f80fd5b5060095461013f565b3480156103b6575f80fd5b506102806103c5366004611b88565b610cf3565b3480156103d5575f80fd5b506103de610d01565b6040516101499190611d83565b3480156103f6575f80fd5b506103ff610d61565b6040516101499190611e27565b348015610417575f80fd5b5061013f600a5481565b335f908152600b602052604090206002015460ff166104875760405162461bcd60e51b815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f746572730000000000000060448201526064015b60405180910390fd5b600654335f908152600b6020526040902054036104d65760405162461bcd60e51b815260206004820152600d60248201526c185b1c9958591e481d9bdd1959609a1b604482015260640161047e565b335f908152600b602052604081206001810180549087905581546006549092559181900361050557505061063e565b5f54841461051457505061063e565b60016006546105239190611e4d565b8114158061055f5750848484336040516020016105439493929190611e60565b604051602081830303815290604052805190602001205f1c8214155b156105cb575f5b5f548110156105c3576001600160ff1b03600c5f838154811061058b5761058b611eaa565b905f5260205f20016040516105a09190611ef6565b9081526040805160209281900383019020335f9081529252902055600101610566565b50505061063e565b5f5b8481101561063a578585828181106105e7576105e7611eaa565b90506020020135600c5f838154811061060257610602611eaa565b905f5260205f20016040516106179190611ef6565b9081526040805160209281900383019020335f90815292529020556001016105cd565b5050505b50505050565b61066b60405180608001604052805f81526020015f81526020015f81526020015f81525090565b5f600d600160065461067d9190611e4d565b8154811061068d5761068d611eaa565b905f5260205f2001836040516106a39190611f68565b90815260200160405180910390206040518060600160405290815f820154815260200160018201548152602001600282015f9054906101000a900460ff1660018111156106f2576106f2611f83565b600181111561070357610703611f83565b8152505090505f604051806080016040528060016006546107249190611e4d565b8152602001835f01518152602001836020015181526020018360400151600181111561075257610752611f83565b9052949350505050565b61078360405180608001604052805f81526020015f81526020015f81526020015f81525090565b5f600d848154811061079757610797611eaa565b905f5260205f2001836040516107ad9190611f68565b90815260200160405180910390206040518060600160405290815f820154815260200160018201548152602001600282015f9054906101000a900460ff1660018111156107fc576107fc611f83565b600181111561080d5761080d611f83565b8152505090505f6040518060800160405280868152602001835f01518152602001836020015181526020018360400151600181111561084e5761084e611f83565b9052925050505b92915050565b6002545f906001600160a01b031633146108875760405162461bcd60e51b815260040161047e90611f97565b600954600a546108979190611fda565b43106109e6575f5b5f548110156108c3576108b181610f14565b6108bc600182611fda565b905061089f565b5060065460075403610933575f5b600554811015610931576001600b5f600584815481106108f3576108f3611eaa565b5f918252602080832091909101546001600160a01b031683528201929092526040019020600201805460ff19169115159190911790556001016108d1565b505b600654600754610944906001611fed565b036109515761095161122a565b43600a81905550600160065f82825461096a9190611fda565b909155505060085461097d906002611fed565b600654036109965760018054610994915f916117e0565b505b60065460095460408051928352436020840152429083015260608201527fb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e59060800160405180910390a150600190565b505f90565b600181815481106109fa575f80fd5b905f5260205f20015f915090508054610a1290611ebe565b80601f0160208091040260200160405190810160405280929190818152602001828054610a3e90611ebe565b8015610a895780601f10610a6057610100808354040283529160200191610a89565b820191905f5260205f20905b815481529060010190602001808311610a6c57829003601f168201915b505050505081565b6002546001600160a01b03163314610abb5760405162461bcd60e51b815260040161047e90611f97565b80515f03610b035760405162461bcd60e51b8152602060048201526015602482015274566f746572732063616e277420626520656d70747960581b604482015260640161047e565b610b1b815f60018451610b169190611e4d565b6113fe565b8051610b2e906005906020840190611834565b5050600654600755565b6003546001600160a01b03163314610b8b5760405162461bcd60e51b81526020600482015260166024820152753932b9ba3934b1ba32b2103a379037b832b930ba37b960511b604482015260640161047e565b80515f03610bd45760405162461bcd60e51b815260206004820152601660248201527573796d626f6c732063616e277420626520656d70747960501b604482015260640161047e565b600654600854610be5906001611fed565b14158015610bf7575060065460085414155b610c435760405162461bcd60e51b815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e640000604482015260640161047e565b8051610c56906001906020840190611893565b5060065460088190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d908290610c8e906001611fda565b604051610c9c929190612014565b60405180910390a150565b6002546001600160a01b03163314610cd15760405162461bcd60e51b815260040161047e90611f97565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b5f81815481106109fa575f80fd5b60606005805480602002602001604051908101604052809291908181526020018280548015610d5757602002820191905f5260205f20905b81546001600160a01b03168152600190910190602001808311610d39575b5050505050905090565b60606006546008546001610d759190611fed565b03610e4c576001805480602002602001604051908101604052809291908181526020015f905b82821015610e43578382905f5260205f20018054610db890611ebe565b80601f0160208091040260200160405190810160405280929190818152602001828054610de490611ebe565b8015610e2f5780601f10610e0657610100808354040283529160200191610e2f565b820191905f5260205f20905b815481529060010190602001808311610e1257829003601f168201915b505050505081526020019060010190610d9b565b50505050905090565b5f805480602002602001604051908101604052809291908181526020015f905b82821015610e43578382905f5260205f20018054610e8990611ebe565b80601f0160208091040260200160405190810160405280929190818152602001828054610eb590611ebe565b8015610f005780601f10610ed757610100808354040283529160200191610f00565b820191905f5260205f20905b815481529060010190602001808311610ee357829003601f168201915b505050505081526020019060010190610e6c565b5f808281548110610f2757610f27611eaa565b905f5260205f20018054610f3a90611ebe565b80601f0160208091040260200160405190810160405280929190818152602001828054610f6690611ebe565b8015610fb15780601f10610f8857610100808354040283529160200191610fb1565b820191905f5260205f20905b815481529060010190602001808311610f9457829003601f168201915b505050505090505f60048054905067ffffffffffffffff811115610fd757610fd76119fb565b604051908082528060200260200182016040528015611000578160200160208202803683370190505b5090505f805b600454811015611108575f6004828154811061102457611024611eaa565b5f9182526020808320909101546006546001600160a01b03909116808452600b90925260409092205490925014158061109657506001600160ff1b03600c866040516110709190611f68565b90815260408051602092819003830190206001600160a01b0385165f9081529252902054145b156110a15750611100565b600c856040516110b19190611f68565b90815260408051602092819003830190206001600160a01b0384165f908152925290205484846110e081612035565b9550815181106110f2576110f2611eaa565b602002602001018181525050505b600101611006565b505f600d600160065461111b9190611e4d565b8154811061112b5761112b611eaa565b905f5260205f2001846040516111419190611f68565b9081526040519081900360200190205490506001821561116b5761116584846115aa565b91505f90505b600d8054600190810182555f9190915260408051606081018252848152426020820152919082019083908111156111a4576111a4611f83565b815250600d600654815481106111bc576111bc611eaa565b905f5260205f2001866040516111d29190611f68565b90815260200160405180910390205f820151815f0155602082015181600101556040820151816002015f6101000a81548160ff0219169083600181111561121b5761121b611f83565b02179055505050505050505050565b5f805b6004548210801561123f575060055481105b1561137e576005818154811061125757611257611eaa565b5f91825260209091200154600480546001600160a01b03909216918490811061128257611282611eaa565b5f918252602090912001546001600160a01b0316036112bb57816112a581612035565b92505080806112b390612035565b91505061122d565b600581815481106112ce576112ce611eaa565b5f91825260209091200154600480546001600160a01b0390921691849081106112f9576112f9611eaa565b5f918252602090912001546001600160a01b0316101561137457600b5f6004848154811061132957611329611eaa565b5f9182526020808320909101546001600160a01b0316835282019290925260400181208181556001810191909155600201805460ff191690558161136c81612035565b92505061122d565b806112b381612035565b6004548210156113e957600b5f6004848154811061139e5761139e611eaa565b5f9182526020808320909101546001600160a01b0316835282019290925260400181208181556001810191909155600201805460ff19169055816113e181612035565b92505061137e565b600580546113f9916004916118d7565b505050565b80821261140a57505050565b81815f85600261141a858561204d565b6114249190612087565b61142e9087611fed565b8151811061143e5761143e611eaa565b602002602001015190505b81831361157c575b806001600160a01b031686848151811061146d5761146d611eaa565b60200260200101516001600160a01b03161015611496578261148e816120b3565b935050611451565b806001600160a01b03168683815181106114b2576114b2611eaa565b60200260200101516001600160a01b031611156114db57816114d3816120ca565b925050611496565b818313611577578582815181106114f4576114f4611eaa565b602002602001015186848151811061150e5761150e611eaa565b602002602001015187858151811061152857611528611eaa565b6020026020010188858151811061154157611541611eaa565b6001600160a01b0393841660209182029290920101529116905282611565816120b3565b9350508180611573906120ca565b9250505b611449565b8185121561158f5761158f8686846113fe565b838312156115a2576115a28684866113fe565b505050505050565b5f815f036115b957505f610855565b6115ce835f6115c9600186611e4d565b611669565b5f6115da6002846120e5565b90506115e76002846120f8565b1561160b578381815181106115fe576115fe611eaa565b6020026020010151611661565b600284828151811061161f5761161f611eaa565b6020026020010151856001846116359190611e4d565b8151811061164557611645611eaa565b60200260200101516116579190611fed565b6116619190612087565b949350505050565b8181808203611679575050505050565b5f856002611687878761204d565b6116919190612087565b61169b9087611fed565b815181106116ab576116ab611eaa565b602002602001015190505b8183136117ba575b808684815181106116d1576116d1611eaa565b602002602001015112156116f157826116e9816120b3565b9350506116be565b85828151811061170357611703611eaa565b6020026020010151811215611724578161171c816120ca565b9250506116f1565b8183136117b55785828151811061173d5761173d611eaa565b602002602001015186848151811061175757611757611eaa565b602002602001015187858151811061177157611771611eaa565b6020026020010188858151811061178a5761178a611eaa565b602090810291909101019190915252816117a3816120ca565b92505082806117b1906120b3565b9350505b6116b6565b818512156117cd576117cd868684611669565b838312156115a2576115a2868486611669565b828054828255905f5260205f20908101928215611824575f5260205f209182015b8281111561182457816118148482612156565b5091600101919060010190611801565b50611830929150611913565b5090565b828054828255905f5260205f20908101928215611887579160200282015b8281111561188757825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611852565b5061183092915061192f565b828054828255905f5260205f20908101928215611824579160200282015b8281111561182457825182906118c7908261222b565b50916020019190600101906118b1565b828054828255905f5260205f20908101928215611887575f5260205f209182015b828111156118875782548255916001019190600101906118f8565b80821115611830575f6119268282611943565b50600101611913565b5b80821115611830575f8155600101611930565b50805461194f90611ebe565b5f825580601f1061195e575050565b601f0160209004905f5260205f209081019061197a919061192f565b50565b5f805f8060608587031215611990575f80fd5b84359350602085013567ffffffffffffffff808211156119ae575f80fd5b818701915087601f8301126119c1575f80fd5b8135818111156119cf575f80fd5b8860208260051b85010111156119e3575f80fd5b95986020929092019750949560400135945092505050565b634e487b7160e01b5f52604160045260245ffd5b604051601f8201601f1916810167ffffffffffffffff81118282101715611a3857611a386119fb565b604052919050565b5f82601f830112611a4f575f80fd5b813567ffffffffffffffff811115611a6957611a696119fb565b611a7c601f8201601f1916602001611a0f565b818152846020838601011115611a90575f80fd5b816020850160208301375f918101602001919091529392505050565b5f60208284031215611abc575f80fd5b813567ffffffffffffffff811115611ad2575f80fd5b61166184828501611a40565b5f8060408385031215611aef575f80fd5b82359150602083013567ffffffffffffffff811115611b0c575f80fd5b611b1885828601611a40565b9150509250929050565b80356001600160a01b0381168114611b38575f80fd5b919050565b5f8060408385031215611b4e575f80fd5b823567ffffffffffffffff811115611b64575f80fd5b611b7085828601611a40565b925050611b7f60208401611b22565b90509250929050565b5f60208284031215611b98575f80fd5b5035919050565b5f5b83811015611bb9578181015183820152602001611ba1565b50505f910152565b5f8151808452611bd8816020860160208601611b9f565b601f01601f19169290920160200192915050565b602081525f611bfe6020830184611bc1565b9392505050565b5f60208284031215611c15575f80fd5b611bfe82611b22565b5f67ffffffffffffffff821115611c3757611c376119fb565b5060051b60200190565b5f6020808385031215611c52575f80fd5b823567ffffffffffffffff811115611c68575f80fd5b8301601f81018513611c78575f80fd5b8035611c8b611c8682611c1e565b611a0f565b81815260059190911b82018301908381019087831115611ca9575f80fd5b928401925b82841015611cce57611cbf84611b22565b82529284019290840190611cae565b979650505050505050565b5f6020808385031215611cea575f80fd5b823567ffffffffffffffff80821115611d01575f80fd5b818501915085601f830112611d14575f80fd5b8135611d22611c8682611c1e565b81815260059190911b83018401908481019088831115611d40575f80fd5b8585015b83811015611d7657803585811115611d5a575f80fd5b611d688b89838a0101611a40565b845250918601918601611d44565b5098975050505050505050565b602080825282518282018190525f9190848201906040850190845b81811015611dc35783516001600160a01b031683529284019291840191600101611d9e565b50909695505050505050565b5f8282518085526020808601955060208260051b840101602086015f5b84811015611e1a57601f19868403018952611e08838351611bc1565b98840198925090830190600101611dec565b5090979650505050505050565b602081525f611bfe6020830184611dcf565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561085557610855611e39565b5f8186825b87811015611e83578135835260209283019290910190600101611e65565b5050938452505060601b6bffffffffffffffffffffffff1916602082015260340192915050565b634e487b7160e01b5f52603260045260245ffd5b600181811c90821680611ed257607f821691505b602082108103611ef057634e487b7160e01b5f52602260045260245ffd5b50919050565b5f808354611f0381611ebe565b60018281168015611f1b5760018114611f3057611f5c565b60ff1984168752821515830287019450611f5c565b875f526020805f205f5b85811015611f535781548a820152908401908201611f3a565b50505082870194505b50929695505050505050565b5f8251611f79818460208701611b9f565b9190910192915050565b634e487b7160e01b5f52602160045260245ffd5b60208082526023908201527f7265737472696374656420746f20746865206175746f6e69747920636f6e74726040820152621858dd60ea1b606082015260800190565b8082018082111561085557610855611e39565b8082018281125f83128015821682158216171561200c5761200c611e39565b505092915050565b604081525f6120266040830185611dcf565b90508260208301529392505050565b5f6001820161204657612046611e39565b5060010190565b8181035f83128015838313168383128216171561206c5761206c611e39565b5092915050565b634e487b7160e01b5f52601260045260245ffd5b5f8261209557612095612073565b600160ff1b82145f19841416156120ae576120ae611e39565b500590565b5f6001600160ff1b01820161204657612046611e39565b5f600160ff1b82016120de576120de611e39565b505f190190565b5f826120f3576120f3612073565b500490565b5f8261210657612106612073565b500690565b601f8211156113f957805f5260205f20601f840160051c810160208510156121305750805b601f840160051c820191505b8181101561214f575f815560010161213c565b5050505050565b818103612161575050565b61216b8254611ebe565b67ffffffffffffffff811115612183576121836119fb565b612197816121918454611ebe565b8461210b565b5f601f8211600181146121c8575f83156121b15750848201545b5f19600385901b1c1916600184901b17845561214f565b5f8581526020808220868352908220601f198616925b838110156121fe57828601548255600195860195909101906020016121de565b508583101561221b57818501545f19600388901b60f8161c191681555b5050505050600190811b01905550565b815167ffffffffffffffff811115612245576122456119fb565b612253816121918454611ebe565b602080601f831160018114612286575f841561226f5750858301515b5f19600386901b1c1916600185901b1785556115a2565b5f85815260208120601f198616915b828110156122b457888601518255948401946001909101908401612295565b508582101561221b5793909601515f1960f8600387901b161c19169092555050600190811b0190555056fea264697066735822122041186a07c007283e8ae4b68295e14bfd62b1bdc745a9ae28d6925996c7af0b3764736f6c63430008170033",
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
func DeployOracle(auth *bind.TransactOpts, backend bind.ContractBackend, _voters []common.Address, _autonity common.Address, _operator common.Address, _symbols []string, _votePeriod *big.Int) (common.Address, *types.Transaction, *Oracle, error) {
	parsed, err := OracleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OracleBin), backend, _voters, _autonity, _operator, _symbols, _votePeriod)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	OracleCaller     // Read-only binding to the contract
	OracleTransactor // Write-only binding to the contract
	OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleSession struct {
	Contract     *Oracle           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleCallerSession struct {
	Contract *OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleTransactorSession struct {
	Contract     *OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRaw struct {
	Contract *Oracle // Generic contract binding to access the raw methods on
}

// OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleCallerRaw struct {
	Contract *OracleCaller // Generic read-only contract binding to access the raw methods on
}

// OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleTransactorRaw struct {
	Contract *OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle creates a new instance of Oracle, bound to a specific deployed contract.
func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
	contract, err := bindOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// NewOracleCaller creates a new read-only instance of Oracle, bound to a specific deployed contract.
func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
	contract, err := bindOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleCaller{contract: contract}, nil
}

// NewOracleTransactor creates a new write-only instance of Oracle, bound to a specific deployed contract.
func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
	contract, err := bindOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleTransactor{contract: contract}, nil
}

// NewOracleFilterer creates a new log filterer instance of Oracle, bound to a specific deployed contract.
func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
	contract, err := bindOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleFilterer{contract: contract}, nil
}

// bindOracle binds a generic wrapper to an already deployed contract.
func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transact(opts, method, params...)
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() pure returns(uint256)
func (_Oracle *OracleCaller) GetPrecision(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getPrecision")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() pure returns(uint256)
func (_Oracle *OracleSession) GetPrecision() (*big.Int, error) {
	return _Oracle.Contract.GetPrecision(&_Oracle.CallOpts)
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() pure returns(uint256)
func (_Oracle *OracleCallerSession) GetPrecision() (*big.Int, error) {
	return _Oracle.Contract.GetPrecision(&_Oracle.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle *OracleCaller) GetRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle *OracleSession) GetRound() (*big.Int, error) {
	return _Oracle.Contract.GetRound(&_Oracle.CallOpts)
}

// GetRound is a free data retrieval call binding the contract method 0x9f8743f7.
//
// Solidity: function getRound() view returns(uint256)
func (_Oracle *OracleCallerSession) GetRound() (*big.Int, error) {
	return _Oracle.Contract.GetRound(&_Oracle.CallOpts)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCaller) GetRoundData(opts *bind.CallOpts, _round *big.Int, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getRoundData", _round, _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _Oracle.Contract.GetRoundData(&_Oracle.CallOpts, _round, _symbol)
}

// GetRoundData is a free data retrieval call binding the contract method 0x3c8510fd.
//
// Solidity: function getRoundData(uint256 _round, string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCallerSession) GetRoundData(_round *big.Int, _symbol string) (IOracleRoundData, error) {
	return _Oracle.Contract.GetRoundData(&_Oracle.CallOpts, _round, _symbol)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle *OracleCaller) GetSymbols(opts *bind.CallOpts) ([]string, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getSymbols")

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle *OracleSession) GetSymbols() ([]string, error) {
	return _Oracle.Contract.GetSymbols(&_Oracle.CallOpts)
}

// GetSymbols is a free data retrieval call binding the contract method 0xdf7f710e.
//
// Solidity: function getSymbols() view returns(string[])
func (_Oracle *OracleCallerSession) GetSymbols() ([]string, error) {
	return _Oracle.Contract.GetSymbols(&_Oracle.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle *OracleCaller) GetVotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getVotePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle *OracleSession) GetVotePeriod() (*big.Int, error) {
	return _Oracle.Contract.GetVotePeriod(&_Oracle.CallOpts)
}

// GetVotePeriod is a free data retrieval call binding the contract method 0xb78dec52.
//
// Solidity: function getVotePeriod() view returns(uint256)
func (_Oracle *OracleCallerSession) GetVotePeriod() (*big.Int, error) {
	return _Oracle.Contract.GetVotePeriod(&_Oracle.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle *OracleCaller) GetVoters(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getVoters")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle *OracleSession) GetVoters() ([]common.Address, error) {
	return _Oracle.Contract.GetVoters(&_Oracle.CallOpts)
}

// GetVoters is a free data retrieval call binding the contract method 0xcdd72253.
//
// Solidity: function getVoters() view returns(address[])
func (_Oracle *OracleCallerSession) GetVoters() ([]common.Address, error) {
	return _Oracle.Contract.GetVoters(&_Oracle.CallOpts)
}

// LastRoundBlock is a free data retrieval call binding the contract method 0xe6a02a28.
//
// Solidity: function lastRoundBlock() view returns(uint256)
func (_Oracle *OracleCaller) LastRoundBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "lastRoundBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastRoundBlock is a free data retrieval call binding the contract method 0xe6a02a28.
//
// Solidity: function lastRoundBlock() view returns(uint256)
func (_Oracle *OracleSession) LastRoundBlock() (*big.Int, error) {
	return _Oracle.Contract.LastRoundBlock(&_Oracle.CallOpts)
}

// LastRoundBlock is a free data retrieval call binding the contract method 0xe6a02a28.
//
// Solidity: function lastRoundBlock() view returns(uint256)
func (_Oracle *OracleCallerSession) LastRoundBlock() (*big.Int, error) {
	return _Oracle.Contract.LastRoundBlock(&_Oracle.CallOpts)
}

// LastVoterUpdateRound is a free data retrieval call binding the contract method 0xaa2f89b5.
//
// Solidity: function lastVoterUpdateRound() view returns(int256)
func (_Oracle *OracleCaller) LastVoterUpdateRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "lastVoterUpdateRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastVoterUpdateRound is a free data retrieval call binding the contract method 0xaa2f89b5.
//
// Solidity: function lastVoterUpdateRound() view returns(int256)
func (_Oracle *OracleSession) LastVoterUpdateRound() (*big.Int, error) {
	return _Oracle.Contract.LastVoterUpdateRound(&_Oracle.CallOpts)
}

// LastVoterUpdateRound is a free data retrieval call binding the contract method 0xaa2f89b5.
//
// Solidity: function lastVoterUpdateRound() view returns(int256)
func (_Oracle *OracleCallerSession) LastVoterUpdateRound() (*big.Int, error) {
	return _Oracle.Contract.LastVoterUpdateRound(&_Oracle.CallOpts)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCaller) LatestRoundData(opts *bind.CallOpts, _symbol string) (IOracleRoundData, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "latestRoundData", _symbol)

	if err != nil {
		return *new(IOracleRoundData), err
	}

	out0 := *abi.ConvertType(out[0], new(IOracleRoundData)).(*IOracleRoundData)

	return out0, err

}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _Oracle.Contract.LatestRoundData(&_Oracle.CallOpts, _symbol)
}

// LatestRoundData is a free data retrieval call binding the contract method 0x33f98c77.
//
// Solidity: function latestRoundData(string _symbol) view returns((uint256,int256,uint256,uint256) data)
func (_Oracle *OracleCallerSession) LatestRoundData(_symbol string) (IOracleRoundData, error) {
	return _Oracle.Contract.LatestRoundData(&_Oracle.CallOpts, _symbol)
}

// NewSymbols is a free data retrieval call binding the contract method 0x5281b5c6.
//
// Solidity: function newSymbols(uint256 ) view returns(string)
func (_Oracle *OracleCaller) NewSymbols(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "newSymbols", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NewSymbols is a free data retrieval call binding the contract method 0x5281b5c6.
//
// Solidity: function newSymbols(uint256 ) view returns(string)
func (_Oracle *OracleSession) NewSymbols(arg0 *big.Int) (string, error) {
	return _Oracle.Contract.NewSymbols(&_Oracle.CallOpts, arg0)
}

// NewSymbols is a free data retrieval call binding the contract method 0x5281b5c6.
//
// Solidity: function newSymbols(uint256 ) view returns(string)
func (_Oracle *OracleCallerSession) NewSymbols(arg0 *big.Int) (string, error) {
	return _Oracle.Contract.NewSymbols(&_Oracle.CallOpts, arg0)
}

// Reports is a free data retrieval call binding the contract method 0x4c56ea56.
//
// Solidity: function reports(string , address ) view returns(int256)
func (_Oracle *OracleCaller) Reports(opts *bind.CallOpts, arg0 string, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "reports", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reports is a free data retrieval call binding the contract method 0x4c56ea56.
//
// Solidity: function reports(string , address ) view returns(int256)
func (_Oracle *OracleSession) Reports(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Oracle.Contract.Reports(&_Oracle.CallOpts, arg0, arg1)
}

// Reports is a free data retrieval call binding the contract method 0x4c56ea56.
//
// Solidity: function reports(string , address ) view returns(int256)
func (_Oracle *OracleCallerSession) Reports(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Oracle.Contract.Reports(&_Oracle.CallOpts, arg0, arg1)
}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_Oracle *OracleCaller) Round(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "round")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_Oracle *OracleSession) Round() (*big.Int, error) {
	return _Oracle.Contract.Round(&_Oracle.CallOpts)
}

// Round is a free data retrieval call binding the contract method 0x146ca531.
//
// Solidity: function round() view returns(uint256)
func (_Oracle *OracleCallerSession) Round() (*big.Int, error) {
	return _Oracle.Contract.Round(&_Oracle.CallOpts)
}

// SymbolUpdatedRound is a free data retrieval call binding the contract method 0x08f21ff5.
//
// Solidity: function symbolUpdatedRound() view returns(int256)
func (_Oracle *OracleCaller) SymbolUpdatedRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "symbolUpdatedRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SymbolUpdatedRound is a free data retrieval call binding the contract method 0x08f21ff5.
//
// Solidity: function symbolUpdatedRound() view returns(int256)
func (_Oracle *OracleSession) SymbolUpdatedRound() (*big.Int, error) {
	return _Oracle.Contract.SymbolUpdatedRound(&_Oracle.CallOpts)
}

// SymbolUpdatedRound is a free data retrieval call binding the contract method 0x08f21ff5.
//
// Solidity: function symbolUpdatedRound() view returns(int256)
func (_Oracle *OracleCallerSession) SymbolUpdatedRound() (*big.Int, error) {
	return _Oracle.Contract.SymbolUpdatedRound(&_Oracle.CallOpts)
}

// Symbols is a free data retrieval call binding the contract method 0xccce413b.
//
// Solidity: function symbols(uint256 ) view returns(string)
func (_Oracle *OracleCaller) Symbols(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "symbols", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbols is a free data retrieval call binding the contract method 0xccce413b.
//
// Solidity: function symbols(uint256 ) view returns(string)
func (_Oracle *OracleSession) Symbols(arg0 *big.Int) (string, error) {
	return _Oracle.Contract.Symbols(&_Oracle.CallOpts, arg0)
}

// Symbols is a free data retrieval call binding the contract method 0xccce413b.
//
// Solidity: function symbols(uint256 ) view returns(string)
func (_Oracle *OracleCallerSession) Symbols(arg0 *big.Int) (string, error) {
	return _Oracle.Contract.Symbols(&_Oracle.CallOpts, arg0)
}

// VotePeriod is a free data retrieval call binding the contract method 0xa7813587.
//
// Solidity: function votePeriod() view returns(uint256)
func (_Oracle *OracleCaller) VotePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "votePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VotePeriod is a free data retrieval call binding the contract method 0xa7813587.
//
// Solidity: function votePeriod() view returns(uint256)
func (_Oracle *OracleSession) VotePeriod() (*big.Int, error) {
	return _Oracle.Contract.VotePeriod(&_Oracle.CallOpts)
}

// VotePeriod is a free data retrieval call binding the contract method 0xa7813587.
//
// Solidity: function votePeriod() view returns(uint256)
func (_Oracle *OracleCallerSession) VotePeriod() (*big.Int, error) {
	return _Oracle.Contract.VotePeriod(&_Oracle.CallOpts)
}

// VotingInfo is a free data retrieval call binding the contract method 0x5412b3ae.
//
// Solidity: function votingInfo(address ) view returns(uint256 round, uint256 commit, bool isVoter)
func (_Oracle *OracleCaller) VotingInfo(opts *bind.CallOpts, arg0 common.Address) (struct {
	Round   *big.Int
	Commit  *big.Int
	IsVoter bool
}, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "votingInfo", arg0)

	outstruct := new(struct {
		Round   *big.Int
		Commit  *big.Int
		IsVoter bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Round = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Commit = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IsVoter = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// VotingInfo is a free data retrieval call binding the contract method 0x5412b3ae.
//
// Solidity: function votingInfo(address ) view returns(uint256 round, uint256 commit, bool isVoter)
func (_Oracle *OracleSession) VotingInfo(arg0 common.Address) (struct {
	Round   *big.Int
	Commit  *big.Int
	IsVoter bool
}, error) {
	return _Oracle.Contract.VotingInfo(&_Oracle.CallOpts, arg0)
}

// VotingInfo is a free data retrieval call binding the contract method 0x5412b3ae.
//
// Solidity: function votingInfo(address ) view returns(uint256 round, uint256 commit, bool isVoter)
func (_Oracle *OracleCallerSession) VotingInfo(arg0 common.Address) (struct {
	Round   *big.Int
	Commit  *big.Int
	IsVoter bool
}, error) {
	return _Oracle.Contract.VotingInfo(&_Oracle.CallOpts, arg0)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_Oracle *OracleTransactor) Finalize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "finalize")
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_Oracle *OracleSession) Finalize() (*types.Transaction, error) {
	return _Oracle.Contract.Finalize(&_Oracle.TransactOpts)
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool)
func (_Oracle *OracleTransactorSession) Finalize() (*types.Transaction, error) {
	return _Oracle.Contract.Finalize(&_Oracle.TransactOpts)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle *OracleTransactor) SetOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setOperator", _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle *OracleSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetOperator(&_Oracle.TransactOpts, _operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle *OracleTransactorSession) SetOperator(_operator common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetOperator(&_Oracle.TransactOpts, _operator)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle *OracleTransactor) SetSymbols(opts *bind.TransactOpts, _symbols []string) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setSymbols", _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle *OracleSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _Oracle.Contract.SetSymbols(&_Oracle.TransactOpts, _symbols)
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle *OracleTransactorSession) SetSymbols(_symbols []string) (*types.Transaction, error) {
	return _Oracle.Contract.SetSymbols(&_Oracle.TransactOpts, _symbols)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_Oracle *OracleTransactor) SetVoters(opts *bind.TransactOpts, _newVoters []common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setVoters", _newVoters)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_Oracle *OracleSession) SetVoters(_newVoters []common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetVoters(&_Oracle.TransactOpts, _newVoters)
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_Oracle *OracleTransactorSession) SetVoters(_newVoters []common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetVoters(&_Oracle.TransactOpts, _newVoters)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_Oracle *OracleTransactor) Vote(opts *bind.TransactOpts, _commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "vote", _commit, _reports, _salt)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_Oracle *OracleSession) Vote(_commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Vote(&_Oracle.TransactOpts, _commit, _reports, _salt)
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_Oracle *OracleTransactorSession) Vote(_commit *big.Int, _reports []*big.Int, _salt *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.Vote(&_Oracle.TransactOpts, _commit, _reports, _salt)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Oracle.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Oracle.Contract.Fallback(&_Oracle.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Oracle.Contract.Fallback(&_Oracle.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleSession) Receive() (*types.Transaction, error) {
	return _Oracle.Contract.Receive(&_Oracle.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleTransactorSession) Receive() (*types.Transaction, error) {
	return _Oracle.Contract.Receive(&_Oracle.TransactOpts)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Round      *big.Int
	Height     *big.Int
	Timestamp  *big.Int
	VotePeriod *big.Int
	Raw        types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Symbols []string
	Round   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Voter common.Address
	Votes []*big.Int
	Raw   types.Log // Blockchain specific contextual infos
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

// PrecompiledMetaData contains all meta data concerning the Precompiled contract.
var PrecompiledMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ACCUSATION_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INNOCENCE_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MISBEHAVIOUR_CONTRACT\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"4dc925d3": "ACCUSATION_CONTRACT()",
		"8e153dc3": "INNOCENCE_CONTRACT()",
		"925c5492": "MISBEHAVIOUR_CONTRACT()",
	},
	Bin: "0x60b1610034600b8282823980515f1a60731461002857634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe73000000000000000000000000000000000000000030146080604052600436106046575f3560e01c80634dc925d314604a5780638e153dc314606d578063925c5492146074575b5f80fd5b605160fc81565b6040516001600160a01b03909116815260200160405180910390f35b605160fd81565b605160fe8156fea2646970667358221220805fa2bb2e53207cd25df6a1713841caca03c3a4d5f0e414ba2cb69a6b15e68c64736f6c63430008170033",
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
func DeployPrecompiled(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Precompiled, error) {
	parsed, err := PrecompiledMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PrecompiledBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Precompiled{PrecompiledCaller: PrecompiledCaller{contract: contract}, PrecompiledTransactor: PrecompiledTransactor{contract: contract}, PrecompiledFilterer: PrecompiledFilterer{contract: contract}}, nil
}

// Precompiled is an auto generated Go binding around an Ethereum contract.
type Precompiled struct {
	PrecompiledCaller     // Read-only binding to the contract
	PrecompiledTransactor // Write-only binding to the contract
	PrecompiledFilterer   // Log filterer for contract events
}

// PrecompiledCaller is an auto generated read-only Go binding around an Ethereum contract.
type PrecompiledCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompiledTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PrecompiledTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompiledFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PrecompiledFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PrecompiledSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PrecompiledSession struct {
	Contract     *Precompiled      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PrecompiledCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PrecompiledCallerSession struct {
	Contract *PrecompiledCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PrecompiledTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PrecompiledTransactorSession struct {
	Contract     *PrecompiledTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PrecompiledRaw is an auto generated low-level Go binding around an Ethereum contract.
type PrecompiledRaw struct {
	Contract *Precompiled // Generic contract binding to access the raw methods on
}

// PrecompiledCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PrecompiledCallerRaw struct {
	Contract *PrecompiledCaller // Generic read-only contract binding to access the raw methods on
}

// PrecompiledTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PrecompiledTransactorRaw struct {
	Contract *PrecompiledTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPrecompiled creates a new instance of Precompiled, bound to a specific deployed contract.
func NewPrecompiled(address common.Address, backend bind.ContractBackend) (*Precompiled, error) {
	contract, err := bindPrecompiled(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Precompiled{PrecompiledCaller: PrecompiledCaller{contract: contract}, PrecompiledTransactor: PrecompiledTransactor{contract: contract}, PrecompiledFilterer: PrecompiledFilterer{contract: contract}}, nil
}

// NewPrecompiledCaller creates a new read-only instance of Precompiled, bound to a specific deployed contract.
func NewPrecompiledCaller(address common.Address, caller bind.ContractCaller) (*PrecompiledCaller, error) {
	contract, err := bindPrecompiled(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PrecompiledCaller{contract: contract}, nil
}

// NewPrecompiledTransactor creates a new write-only instance of Precompiled, bound to a specific deployed contract.
func NewPrecompiledTransactor(address common.Address, transactor bind.ContractTransactor) (*PrecompiledTransactor, error) {
	contract, err := bindPrecompiled(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PrecompiledTransactor{contract: contract}, nil
}

// NewPrecompiledFilterer creates a new log filterer instance of Precompiled, bound to a specific deployed contract.
func NewPrecompiledFilterer(address common.Address, filterer bind.ContractFilterer) (*PrecompiledFilterer, error) {
	contract, err := bindPrecompiled(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PrecompiledFilterer{contract: contract}, nil
}

// bindPrecompiled binds a generic wrapper to an already deployed contract.
func bindPrecompiled(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PrecompiledABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Precompiled *PrecompiledRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Precompiled.Contract.PrecompiledCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Precompiled *PrecompiledRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Precompiled.Contract.PrecompiledTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Precompiled *PrecompiledRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Precompiled.Contract.PrecompiledTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Precompiled *PrecompiledCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Precompiled.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Precompiled *PrecompiledTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Precompiled.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Precompiled *PrecompiledTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Precompiled.Contract.contract.Transact(opts, method, params...)
}

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) ACCUSATIONCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "ACCUSATION_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) ACCUSATIONCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.ACCUSATIONCONTRACT(&_Precompiled.CallOpts)
}

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) ACCUSATIONCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.ACCUSATIONCONTRACT(&_Precompiled.CallOpts)
}

// INNOCENCECONTRACT is a free data retrieval call binding the contract method 0x8e153dc3.
//
// Solidity: function INNOCENCE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) INNOCENCECONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "INNOCENCE_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// INNOCENCECONTRACT is a free data retrieval call binding the contract method 0x8e153dc3.
//
// Solidity: function INNOCENCE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) INNOCENCECONTRACT() (common.Address, error) {
	return _Precompiled.Contract.INNOCENCECONTRACT(&_Precompiled.CallOpts)
}

// INNOCENCECONTRACT is a free data retrieval call binding the contract method 0x8e153dc3.
//
// Solidity: function INNOCENCE_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) INNOCENCECONTRACT() (common.Address, error) {
	return _Precompiled.Contract.INNOCENCECONTRACT(&_Precompiled.CallOpts)
}

// MISBEHAVIOURCONTRACT is a free data retrieval call binding the contract method 0x925c5492.
//
// Solidity: function MISBEHAVIOUR_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCaller) MISBEHAVIOURCONTRACT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Precompiled.contract.Call(opts, &out, "MISBEHAVIOUR_CONTRACT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MISBEHAVIOURCONTRACT is a free data retrieval call binding the contract method 0x925c5492.
//
// Solidity: function MISBEHAVIOUR_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledSession) MISBEHAVIOURCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.MISBEHAVIOURCONTRACT(&_Precompiled.CallOpts)
}

// MISBEHAVIOURCONTRACT is a free data retrieval call binding the contract method 0x925c5492.
//
// Solidity: function MISBEHAVIOUR_CONTRACT() view returns(address)
func (_Precompiled *PrecompiledCallerSession) MISBEHAVIOURCONTRACT() (common.Address, error) {
	return _Precompiled.Contract.MISBEHAVIOURCONTRACT(&_Precompiled.CallOpts)
}

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
	Bin: "0x608060405234801562000010575f80fd5b50604051620026373803806200263783398101604081905262000033916200012a565b8560400151805f036200005957604051630309cb8760e51b815260040160405180910390fd5b866020015187604001518082106200008457604051630309cb8760e51b815260040160405180910390fd5b505086515f5550602086015160015560408601516002556060860151600355608090950151600455600780546001600160a01b03199081166001600160a01b039687161790915560088054821694861694909417909355600a8054841692851692909217909155600b80548316918416919091179055600980549091169190921617905562000211565b80516001600160a01b038116811462000125575f80fd5b919050565b5f805f805f8086880361014081121562000142575f80fd5b60a081121562000150575f80fd5b5060405160a081016001600160401b03811182821017156200018057634e487b7160e01b5f52604160045260245ffd5b8060405250875181526020880151602082015260408801516040820152606088015160608201526080880151608082015280965050620001c360a088016200010e565b9450620001d360c088016200010e565b9350620001e360e088016200010e565b9250620001f461010088016200010e565b91506200020561012088016200010e565b90509295509295509295565b612418806200021f5f395ff3fe608060405260043610610161575f3560e01c806368cd03f6116100cd578063946ce8cd11610087578063c5ebeaec11610062578063c5ebeaec14610450578063ce4b5bbe1461046f578063eced552614610483578063fbbe699114610497575f80fd5b8063946ce8cd146103f3578063b3ab15fb14610412578063b6b55f2514610431575f80fd5b806368cd03f6146102c857806379502c55146102e95780637adbf973146103365780637b44646a1461035557806383baa17414610374578063840c7e2414610393575f80fd5b806350bf06bf1161011e57806350bf06bf1461022157806352e5a0501461024057806353afe81d1461025f57806354a9f42c1461027e5780635891de721461029d5780635dcc9391146102b1575f80fd5b8063042e02cf14610165578063087966961461019957806315184245146101c65780632e1a7d4d146101e55780632f86556814610206578063402d888314610219575b5f80fd5b348015610170575f80fd5b5061018461017f366004612085565b6104b6565b60405190151581526020015b60405180910390f35b3480156101a4575f80fd5b506101b86101b336600461209e565b6104fd565b604051908152602001610190565b3480156101d1575f80fd5b506101b86101e03660046120c7565b610565565b3480156101f0575f80fd5b506102046101ff3660046120f6565b6105f1565b005b610204610214366004612085565b6107b5565b6102046109e7565b34801561022c575f80fd5b506101b861023b36600461210d565b610bb0565b34801561024b575f80fd5b5061020461025a366004612085565b610c1a565b34801561026a575f80fd5b506102046102793660046120f6565b610c66565b348015610289575f80fd5b506101b8610298366004612085565b610c95565b3480156102a8575f80fd5b506101b8610d09565b3480156102bc575f80fd5b506101b86301e1338081565b3480156102d3575f80fd5b506102dc610f9d565b6040516101909190612135565b3480156102f4575f80fd5b505f5460015460025460035460045461030e949392919085565b604080519586526020860194909452928401919091526060830152608082015260a001610190565b348015610341575f80fd5b50610204610350366004612085565b610ffd565b348015610360575f80fd5b5061020461036f3660046120f6565b611049565b34801561037f575f80fd5b506101b861038e3660046120c7565b6110c0565b34801561039e575f80fd5b506103d36103ad366004612085565b60056020525f908152604090208054600182015460028301546003909301549192909184565b604080519485526020850193909352918301526060820152608001610190565b3480156103fe575f80fd5b5061020461040d3660046120f6565b611115565b34801561041d575f80fd5b5061020461042c366004612085565b61116b565b34801561043c575f80fd5b5061020461044b3660046120f6565b6111b7565b34801561045b575f80fd5b5061020461046a3660046120f6565b6113ad565b34801561047a575f80fd5b506101b861156b565b34801561048e575f80fd5b506101b8601281565b3480156104a2575f80fd5b506101846104b13660046120c7565b61157a565b6001600160a01b0381165f908152600560205260408120816104d882426115cd565b5090506104f582600101546104eb610d09565b600154849061157a565b949350505050565b5f82805f0361051e5760405162bfc92160e01b815260040160405180910390fd5b831580610529575082155b1561054757604051630309cb8760e51b815260040160405180910390fd5b836105528487612195565b61055c91906121c0565b95945050505050565b5f8183111561058757604051630309cb8760e51b815260040160405180910390fd5b84845f6105a86301e133806105a261059f89896121df565b90565b9061163e565b90505f6105bd6105b8848461165c565b61166a565b90505f6105e36105dc6105d561059f6012600a6122d2565b84906116bc565b869061165c565b9a9950505050505050505050565b80805f036106125760405163162908e360e11b815260040160405180910390fd5b335f90815260056020526040902060018101548311156106455760405163162908e360e11b815260040160405180910390fd5b5f61065082426115cd565b5090505f61065c610d09565b9050610672836001015482845f6001015461157a565b1561069057604051636229415360e01b815260040160405180910390fd5b6106a38360020154825f600201546104fd565b8584600101546106b391906121df565b10156106d257604051633a23d82560e01b815260040160405180910390fd5b84836001015f8282546106e591906121df565b909155505060095460405163a9059cbb60e01b8152336004820152602481018790526001600160a01b039091169063a9059cbb906044016020604051808303815f875af1158015610738573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061075c91906122dd565b610779576040516312171d8360e31b815260040160405180910390fd5b60405185815233907f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649060200160405180910390a25050505050565b345f036107d557604051637c946ed760e01b815260040160405180910390fd5b6001600160a01b0381165f908152600560205260408120600281015490910361081157604051638aa5baf360e01b815260040160405180910390fd5b5f8061081d83426115cd565b9150915061083b8360010154610831610d09565b600154859061157a565b61085857604051636ef5bcdd60e11b815260040160405180910390fd5b5f61086383346121df565b90506001840180544286555f9182905560028601829055600386019190915560095460405163a9059cbb60e01b8152336004820152602481018390526001600160a01b039091169063a9059cbb906044016020604051808303815f875af11580156108d0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906108f491906122dd565b610911576040516312171d8360e31b815260040160405180910390fd5b600b546001600160a01b03166344df8e7061092c85876121df565b6040518263ffffffff1660e01b81526004015f604051808303818588803b158015610955575f80fd5b505af1158015610967573d5f803e3d5ffd5b50505050505f8211156109a057604051339083156108fc029084905f818181858888f1935050505015801561099e573d5f803e3d5ffd5b505b6040513381526001600160a01b038716907fc3d81b2125598b9a2b024afe09e33981f0aa5b7bcbe3e30c4303a4dec209ddb4906020015b60405180910390a2505050505050565b345f03610a0757604051637c946ed760e01b815260040160405180910390fd5b335f9081526005602052604081206002810154909103610a3a57604051638aa5baf360e01b815260040160405180910390fd5b5f80610a4683426115cd565b915091508134108015610a635750600354610a6134846121df565b105b15610a815760405163e6bd447960e01b815260040160405180910390fd5b80836003015f828254610a9491906122fc565b90915550504283555f8080610aa986346116ca565b92509250925081866002015f828254610ac291906121df565b9250508190555082866003015f828254610adc91906121df565b90915550508115610b4c57600b5f9054906101000a90046001600160a01b03166001600160a01b03166344df8e70836040518263ffffffff1660e01b81526004015f604051808303818588803b158015610b34575f80fd5b505af1158015610b46573d5f803e3d5ffd5b50505050505b8015610b7e57604051339082156108fc029083905f818181858888f19350505050158015610b7c573d5f803e3d5ffd5b505b60405134815233907f5c16de4f8b59bd9caf0f49a545f25819a895ed223294290b408242e72a594231906020016109d7565b6001600160a01b0382165f908152600560205260408120805484918491821015610bed57604051630309cb8760e51b815260040160405180910390fd5b6001600160a01b0386165f908152600560205260409020610c0e81876115cd565b50979650505050505050565b6008546001600160a01b03163314610c44576040516282b42960e81b815260040160405180910390fd5b600b80546001600160a01b0319166001600160a01b0392909216919091179055565b6008546001600160a01b03163314610c90576040516282b42960e81b815260040160405180910390fd5b600355565b6040516350bf06bf60e01b81526001600160a01b03821660048201524260248201525f9030906350bf06bf90604401602060405180830381865afa158015610cdf573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d03919061230f565b92915050565b600a546040805180820182526007815266272a2716a0aa2760c91b602082015290516333f98c7760e01b81525f9283926001600160a01b03909116916333f98c7791610d5791600401612326565b608060405180830381865afa158015610d72573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d969190612372565b905080606001515f14610dbc5760405163cb08be8160e01b815260040160405180910390fd5b5f816020015113610ddf5760405162bfc92160e01b815260040160405180910390fd5b600a5f9054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e2f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e53919061230f565b610e5f6012600a6122d2565b1115610f0457600a5f9054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610eb5573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ed9919061230f565b610ee56012600a6122d2565b610eef91906121c0565b8160200151610efe9190612195565b91505090565b610f106012600a6122d2565b600a5f9054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f60573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610f84919061230f565b610f8e91906121c0565b8160200151610efe91906121c0565b60606006805480602002602001604051908101604052809291908181526020018280548015610ff357602002820191905f5260205f20905b81546001600160a01b03168152600190910190602001808311610fd5575b5050505050905090565b6007546001600160a01b03163314611027576040516282b42960e81b815260040160405180910390fd5b600a80546001600160a01b0319166001600160a01b0392909216919091179055565b80805f0361106a57604051630309cb8760e51b815260040160405180910390fd5b6001548280821061108e57604051630309cb8760e51b815260040160405180910390fd5b6008546001600160a01b031633146110b8576040516282b42960e81b815260040160405180910390fd5b505050600255565b5f8315806110cc575081155b156110ea57604051630309cb8760e51b815260040160405180910390fd5b6110f66012600a6122d2565b6111009083612195565b8361110b8688612195565b6105529190612195565b600254819080821061113a57604051630309cb8760e51b815260040160405180910390fd5b6008546001600160a01b03163314611164576040516282b42960e81b815260040160405180910390fd5b5050600155565b6007546001600160a01b03163314611195576040516282b42960e81b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0392909216919091179055565b80805f036111d85760405163162908e360e11b815260040160405180910390fd5b600954604051636eb1769f60e11b815233600482015230602482015283916001600160a01b03169063dd62ed3e90604401602060405180830381865afa158015611224573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611248919061230f565b1015611267576040516313be252b60e01b815260040160405180910390fd5b335f90815260056020526040812080549091036112c057600680546001810182555f919091527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0180546001600160a01b031916331790555b4281556001810180548491905f906112d99084906122fc565b90915550506009546040516323b872dd60e01b8152336004820152306024820152604481018590526001600160a01b03909116906323b872dd906064016020604051808303815f875af1158015611332573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061135691906122dd565b611373576040516312171d8360e31b815260040160405180910390fd5b60405183815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2505050565b80805f036113ce5760405163162908e360e11b815260040160405180910390fd5b335f90815260056020526040812090806113e883426115cd565b90925090506113f785836122fc565b60035490925082101561141d5760405163e6bd447960e01b815260040160405180910390fd5b5f611426610d09565b905061143c846001015482855f6001015461157a565b1561145a57604051636229415360e01b815260040160405180910390fd5b5f6114738560010154835f600401545f600201546110c0565b90508084111561149657604051633a23d82560e01b815260040160405180910390fd5b4285556002850180548891905f906114af9084906122fc565b9250508190555082856003015f8282546114c991906122fc565b9091555050600b546040516340c10f1960e01b8152336004820152602481018990526001600160a01b03909116906340c10f19906044015f604051808303815f87803b158015611517575f80fd5b505af1158015611529573d5f803e3d5ffd5b50506040518981523392507fcbc04eca7e9da35cb1393a6135a199ca52e450d5e9251cbd99f7847d33a36750915060200160405180910390a250505050505050565b6115776012600a6122d2565b81565b5f83805f0361159b5760405162bfc92160e01b815260040160405180910390fd5b835f036115aa575f91506115c4565b82846115b68789612195565b6115c091906121c0565b1091505b50949350505050565b5f80825f036115ef57604051630309cb8760e51b815260040160405180910390fd5b5f8460030154856002015461160491906122fc565b85549091508403611617575f915061162a565b5f54855461162791839187610565565b91505b61163482826122fc565b9250509250929050565b5f61165561059f84670de0b6b3a76400008561173b565b9392505050565b5f61165561059f8484611809565b5f81680736ea4425c11ac63081111561169e57604051630d7b1d6560e11b8152600481018490526024015b60405180910390fd5b6714057b7ef767814f81026104f5670de0b6b3a764000082046118bb565b5f61165561059f83856121df565b5f805f80856003015486600201546116e291906122fc565b9050856003015485106116f95785600301546116fb565b845b935080851061170e578560020154611718565b61171884866121df565b9250808511611727575f611731565b61173181866121df565b9150509250925092565b5f80805f19858709858702925082811083820303915050805f0361177257838281611768576117686121ac565b0492505050611655565b8381106117a357604051630c740aef60e31b8152600481018790526024810186905260448101859052606401611695565b5f84868809600260036001881981018916988990049182028318808302840302808302840302808302840302808302840302808302840302918202909203025f889003889004909101858311909403939093029303949094049190911702949350505050565b5f80805f19848609848602925082811083820303915050805f0361183a5750670de0b6b3a764000090049050610d03565b670de0b6b3a7640000811061186c57604051635173648d60e01b81526004810186905260248101859052604401611695565b5f670de0b6b3a764000085870962040000818503049310909103600160ee1b02919091177faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106690291505092915050565b5f81680a688906bd8affffff8111156118ea5760405163b3b6ba1f60e01b815260048101849052602401611695565b5f611901670de0b6b3a7640000604084901b6121c0565b90506104f561059f82600160bf1b67ff00000000000000821615611a175767800000000000000082161561193e5768016a09e667f3bcc9090260401c5b67400000000000000082161561195d576801306fe0a31b7152df0260401c5b67200000000000000082161561197c576801172b83c7d517adce0260401c5b67100000000000000082161561199b5768010b5586cf9890f62a0260401c5b6708000000000000008216156119ba576801059b0d31585743ae0260401c5b6704000000000000008216156119d957680102c9a3e778060ee70260401c5b6702000000000000008216156119f85768010163da9fb33356d80260401c5b670100000000000000821615611a1757680100b1afa5abcbed610260401c5b66ff000000000000821615611b16576680000000000000821615611a445768010058c86da1c09ea20260401c5b6640000000000000821615611a62576801002c605e2e8cec500260401c5b6620000000000000821615611a8057680100162f3904051fa10260401c5b6610000000000000821615611a9e576801000b175effdc76ba0260401c5b6608000000000000821615611abc57680100058ba01fb9f96d0260401c5b6604000000000000821615611ada5768010002c5cc37da94920260401c5b6602000000000000821615611af8576801000162e525ee05470260401c5b6601000000000000821615611b165768010000b17255775c040260401c5b65ff0000000000821615611c0c5765800000000000821615611b41576801000058b91b5bc9ae0260401c5b65400000000000821615611b5e57680100002c5c89d5ec6d0260401c5b65200000000000821615611b7b5768010000162e43f4f8310260401c5b65100000000000821615611b9857680100000b1721bcfc9a0260401c5b65080000000000821615611bb55768010000058b90cf1e6e0260401c5b65040000000000821615611bd2576801000002c5c863b73f0260401c5b65020000000000821615611bef57680100000162e430e5a20260401c5b65010000000000821615611c0c576801000000b1721835510260401c5b64ff00000000821615611cf957648000000000821615611c3557680100000058b90c0b490260401c5b644000000000821615611c515768010000002c5c8601cc0260401c5b642000000000821615611c6d576801000000162e42fff00260401c5b641000000000821615611c895768010000000b17217fbb0260401c5b640800000000821615611ca5576801000000058b90bfce0260401c5b640400000000821615611cc157680100000002c5c85fe30260401c5b640200000000821615611cdd5768010000000162e42ff10260401c5b640100000000821615611cf957680100000000b17217f80260401c5b63ff000000821615611ddd576380000000821615611d205768010000000058b90bfc0260401c5b6340000000821615611d3b576801000000002c5c85fe0260401c5b6320000000821615611d5657680100000000162e42ff0260401c5b6310000000821615611d71576801000000000b17217f0260401c5b6308000000821615611d8c57680100000000058b90c00260401c5b6304000000821615611da75768010000000002c5c8600260401c5b6302000000821615611dc2576801000000000162e4300260401c5b6301000000821615611ddd5768010000000000b172180260401c5b62ff0000821615611eb85762800000821615611e02576801000000000058b90c0260401c5b62400000821615611e1c57680100000000002c5c860260401c5b62200000821615611e365768010000000000162e430260401c5b62100000821615611e5057680100000000000b17210260401c5b62080000821615611e6a5768010000000000058b910260401c5b62040000821615611e84576801000000000002c5c80260401c5b62020000821615611e9e57680100000000000162e40260401c5b62010000821615611eb8576801000000000000b1720260401c5b61ff00821615611f8a57618000821615611edb57680100000000000058b90260401c5b614000821615611ef45768010000000000002c5d0260401c5b612000821615611f0d576801000000000000162e0260401c5b611000821615611f265768010000000000000b170260401c5b610800821615611f3f576801000000000000058c0260401c5b610400821615611f5857680100000000000002c60260401c5b610200821615611f7157680100000000000001630260401c5b610100821615611f8a57680100000000000000b10260401c5b60ff821615612053576080821615611fab57680100000000000000590260401c5b6040821615611fc3576801000000000000002c0260401c5b6020821615611fdb57680100000000000000160260401c5b6010821615611ff3576801000000000000000b0260401c5b600882161561200b57680100000000000000060260401c5b600482161561202357680100000000000000030260401c5b600282161561203b57680100000000000000010260401c5b600182161561205357680100000000000000010260401c5b670de0b6b3a76400000260409190911c60bf031c90565b80356001600160a01b0381168114612080575f80fd5b919050565b5f60208284031215612095575f80fd5b6116558261206a565b5f805f606084860312156120b0575f80fd5b505081359360208301359350604090920135919050565b5f805f80608085870312156120da575f80fd5b5050823594602084013594506040840135936060013592509050565b5f60208284031215612106575f80fd5b5035919050565b5f806040838503121561211e575f80fd5b6121278361206a565b946020939093013593505050565b602080825282518282018190525f9190848201906040850190845b818110156121755783516001600160a01b031683529284019291840191600101612150565b50909695505050505050565b634e487b7160e01b5f52601160045260245ffd5b8082028115828204841417610d0357610d03612181565b634e487b7160e01b5f52601260045260245ffd5b5f826121da57634e487b7160e01b5f52601260045260245ffd5b500490565b81810381811115610d0357610d03612181565b600181815b8085111561222c57815f190482111561221257612212612181565b8085161561221f57918102915b93841c93908002906121f7565b509250929050565b5f8261224257506001610d03565b8161224e57505f610d03565b8160018114612264576002811461226e5761228a565b6001915050610d03565b60ff84111561227f5761227f612181565b50506001821b610d03565b5060208310610133831016604e8410600b84101617156122ad575081810a610d03565b6122b783836121f2565b805f19048211156122ca576122ca612181565b029392505050565b5f6116558383612234565b5f602082840312156122ed575f80fd5b81518015158114611655575f80fd5b80820180821115610d0357610d03612181565b5f6020828403121561231f575f80fd5b5051919050565b5f602080835283518060208501525f5b8181101561235257858101830151858201604001528201612336565b505f604082860101526040601f19601f8301168501019250505092915050565b5f60808284031215612382575f80fd5b6040516080810181811067ffffffffffffffff821117156123b157634e487b7160e01b5f52604160045260245ffd5b806040525082518152602083015160208201526040830151604082015260608301516060820152809150509291505056fea2646970667358221220d58a26ee101a1a9e0e41fcdf8ad5bfed87c35a226ded3e39063cba7795213df864736f6c63430008170033",
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
func DeployStabilization(auth *bind.TransactOpts, backend bind.ContractBackend, config_ StabilizationConfig, autonity common.Address, operator common.Address, oracle common.Address, supplyControl common.Address, collateralToken common.Address) (common.Address, *types.Transaction, *Stabilization, error) {
	parsed, err := StabilizationMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StabilizationBin), backend, config_, autonity, operator, oracle, supplyControl, collateralToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Stabilization{StabilizationCaller: StabilizationCaller{contract: contract}, StabilizationTransactor: StabilizationTransactor{contract: contract}, StabilizationFilterer: StabilizationFilterer{contract: contract}}, nil
}

// Stabilization is an auto generated Go binding around an Ethereum contract.
type Stabilization struct {
	StabilizationCaller     // Read-only binding to the contract
	StabilizationTransactor // Write-only binding to the contract
	StabilizationFilterer   // Log filterer for contract events
}

// StabilizationCaller is an auto generated read-only Go binding around an Ethereum contract.
type StabilizationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StabilizationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StabilizationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StabilizationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StabilizationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StabilizationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StabilizationSession struct {
	Contract     *Stabilization    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StabilizationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StabilizationCallerSession struct {
	Contract *StabilizationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StabilizationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StabilizationTransactorSession struct {
	Contract     *StabilizationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StabilizationRaw is an auto generated low-level Go binding around an Ethereum contract.
type StabilizationRaw struct {
	Contract *Stabilization // Generic contract binding to access the raw methods on
}

// StabilizationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StabilizationCallerRaw struct {
	Contract *StabilizationCaller // Generic read-only contract binding to access the raw methods on
}

// StabilizationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StabilizationTransactorRaw struct {
	Contract *StabilizationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStabilization creates a new instance of Stabilization, bound to a specific deployed contract.
func NewStabilization(address common.Address, backend bind.ContractBackend) (*Stabilization, error) {
	contract, err := bindStabilization(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Stabilization{StabilizationCaller: StabilizationCaller{contract: contract}, StabilizationTransactor: StabilizationTransactor{contract: contract}, StabilizationFilterer: StabilizationFilterer{contract: contract}}, nil
}

// NewStabilizationCaller creates a new read-only instance of Stabilization, bound to a specific deployed contract.
func NewStabilizationCaller(address common.Address, caller bind.ContractCaller) (*StabilizationCaller, error) {
	contract, err := bindStabilization(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StabilizationCaller{contract: contract}, nil
}

// NewStabilizationTransactor creates a new write-only instance of Stabilization, bound to a specific deployed contract.
func NewStabilizationTransactor(address common.Address, transactor bind.ContractTransactor) (*StabilizationTransactor, error) {
	contract, err := bindStabilization(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StabilizationTransactor{contract: contract}, nil
}

// NewStabilizationFilterer creates a new log filterer instance of Stabilization, bound to a specific deployed contract.
func NewStabilizationFilterer(address common.Address, filterer bind.ContractFilterer) (*StabilizationFilterer, error) {
	contract, err := bindStabilization(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StabilizationFilterer{contract: contract}, nil
}

// bindStabilization binds a generic wrapper to an already deployed contract.
func bindStabilization(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StabilizationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stabilization *StabilizationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stabilization.Contract.StabilizationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stabilization *StabilizationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stabilization.Contract.StabilizationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stabilization *StabilizationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stabilization.Contract.StabilizationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stabilization *StabilizationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Stabilization.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stabilization *StabilizationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stabilization.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stabilization *StabilizationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stabilization.Contract.contract.Transact(opts, method, params...)
}

// SCALE is a free data retrieval call binding the contract method 0xeced5526.
//
// Solidity: function SCALE() view returns(uint256)
func (_Stabilization *StabilizationCaller) SCALE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "SCALE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SCALE is a free data retrieval call binding the contract method 0xeced5526.
//
// Solidity: function SCALE() view returns(uint256)
func (_Stabilization *StabilizationSession) SCALE() (*big.Int, error) {
	return _Stabilization.Contract.SCALE(&_Stabilization.CallOpts)
}

// SCALE is a free data retrieval call binding the contract method 0xeced5526.
//
// Solidity: function SCALE() view returns(uint256)
func (_Stabilization *StabilizationCallerSession) SCALE() (*big.Int, error) {
	return _Stabilization.Contract.SCALE(&_Stabilization.CallOpts)
}

// SCALEFACTOR is a free data retrieval call binding the contract method 0xce4b5bbe.
//
// Solidity: function SCALE_FACTOR() view returns(uint256)
func (_Stabilization *StabilizationCaller) SCALEFACTOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "SCALE_FACTOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SCALEFACTOR is a free data retrieval call binding the contract method 0xce4b5bbe.
//
// Solidity: function SCALE_FACTOR() view returns(uint256)
func (_Stabilization *StabilizationSession) SCALEFACTOR() (*big.Int, error) {
	return _Stabilization.Contract.SCALEFACTOR(&_Stabilization.CallOpts)
}

// SCALEFACTOR is a free data retrieval call binding the contract method 0xce4b5bbe.
//
// Solidity: function SCALE_FACTOR() view returns(uint256)
func (_Stabilization *StabilizationCallerSession) SCALEFACTOR() (*big.Int, error) {
	return _Stabilization.Contract.SCALEFACTOR(&_Stabilization.CallOpts)
}

// SECONDSINYEAR is a free data retrieval call binding the contract method 0x5dcc9391.
//
// Solidity: function SECONDS_IN_YEAR() view returns(uint256)
func (_Stabilization *StabilizationCaller) SECONDSINYEAR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "SECONDS_IN_YEAR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SECONDSINYEAR is a free data retrieval call binding the contract method 0x5dcc9391.
//
// Solidity: function SECONDS_IN_YEAR() view returns(uint256)
func (_Stabilization *StabilizationSession) SECONDSINYEAR() (*big.Int, error) {
	return _Stabilization.Contract.SECONDSINYEAR(&_Stabilization.CallOpts)
}

// SECONDSINYEAR is a free data retrieval call binding the contract method 0x5dcc9391.
//
// Solidity: function SECONDS_IN_YEAR() view returns(uint256)
func (_Stabilization *StabilizationCallerSession) SECONDSINYEAR() (*big.Int, error) {
	return _Stabilization.Contract.SECONDSINYEAR(&_Stabilization.CallOpts)
}

// Accounts is a free data retrieval call binding the contract method 0x68cd03f6.
//
// Solidity: function accounts() view returns(address[])
func (_Stabilization *StabilizationCaller) Accounts(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "accounts")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// Accounts is a free data retrieval call binding the contract method 0x68cd03f6.
//
// Solidity: function accounts() view returns(address[])
func (_Stabilization *StabilizationSession) Accounts() ([]common.Address, error) {
	return _Stabilization.Contract.Accounts(&_Stabilization.CallOpts)
}

// Accounts is a free data retrieval call binding the contract method 0x68cd03f6.
//
// Solidity: function accounts() view returns(address[])
func (_Stabilization *StabilizationCallerSession) Accounts() ([]common.Address, error) {
	return _Stabilization.Contract.Accounts(&_Stabilization.CallOpts)
}

// BorrowLimit is a free data retrieval call binding the contract method 0x83baa174.
//
// Solidity: function borrowLimit(uint256 collateral, uint256 price, uint256 targetPrice, uint256 mcr) pure returns(uint256)
func (_Stabilization *StabilizationCaller) BorrowLimit(opts *bind.CallOpts, collateral *big.Int, price *big.Int, targetPrice *big.Int, mcr *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "borrowLimit", collateral, price, targetPrice, mcr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BorrowLimit is a free data retrieval call binding the contract method 0x83baa174.
//
// Solidity: function borrowLimit(uint256 collateral, uint256 price, uint256 targetPrice, uint256 mcr) pure returns(uint256)
func (_Stabilization *StabilizationSession) BorrowLimit(collateral *big.Int, price *big.Int, targetPrice *big.Int, mcr *big.Int) (*big.Int, error) {
	return _Stabilization.Contract.BorrowLimit(&_Stabilization.CallOpts, collateral, price, targetPrice, mcr)
}

// BorrowLimit is a free data retrieval call binding the contract method 0x83baa174.
//
// Solidity: function borrowLimit(uint256 collateral, uint256 price, uint256 targetPrice, uint256 mcr) pure returns(uint256)
func (_Stabilization *StabilizationCallerSession) BorrowLimit(collateral *big.Int, price *big.Int, targetPrice *big.Int, mcr *big.Int) (*big.Int, error) {
	return _Stabilization.Contract.BorrowLimit(&_Stabilization.CallOpts, collateral, price, targetPrice, mcr)
}

// Cdps is a free data retrieval call binding the contract method 0x840c7e24.
//
// Solidity: function cdps(address ) view returns(uint256 timestamp, uint256 collateral, uint256 principal, uint256 interest)
func (_Stabilization *StabilizationCaller) Cdps(opts *bind.CallOpts, arg0 common.Address) (struct {
	Timestamp  *big.Int
	Collateral *big.Int
	Principal  *big.Int
	Interest   *big.Int
}, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "cdps", arg0)

	outstruct := new(struct {
		Timestamp  *big.Int
		Collateral *big.Int
		Principal  *big.Int
		Interest   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Collateral = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Principal = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Interest = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Cdps is a free data retrieval call binding the contract method 0x840c7e24.
//
// Solidity: function cdps(address ) view returns(uint256 timestamp, uint256 collateral, uint256 principal, uint256 interest)
func (_Stabilization *StabilizationSession) Cdps(arg0 common.Address) (struct {
	Timestamp  *big.Int
	Collateral *big.Int
	Principal  *big.Int
	Interest   *big.Int
}, error) {
	return _Stabilization.Contract.Cdps(&_Stabilization.CallOpts, arg0)
}

// Cdps is a free data retrieval call binding the contract method 0x840c7e24.
//
// Solidity: function cdps(address ) view returns(uint256 timestamp, uint256 collateral, uint256 principal, uint256 interest)
func (_Stabilization *StabilizationCallerSession) Cdps(arg0 common.Address) (struct {
	Timestamp  *big.Int
	Collateral *big.Int
	Principal  *big.Int
	Interest   *big.Int
}, error) {
	return _Stabilization.Contract.Cdps(&_Stabilization.CallOpts, arg0)
}

// CollateralPrice is a free data retrieval call binding the contract method 0x5891de72.
//
// Solidity: function collateralPrice() view returns(uint256 price)
func (_Stabilization *StabilizationCaller) CollateralPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "collateralPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CollateralPrice is a free data retrieval call binding the contract method 0x5891de72.
//
// Solidity: function collateralPrice() view returns(uint256 price)
func (_Stabilization *StabilizationSession) CollateralPrice() (*big.Int, error) {
	return _Stabilization.Contract.CollateralPrice(&_Stabilization.CallOpts)
}

// CollateralPrice is a free data retrieval call binding the contract method 0x5891de72.
//
// Solidity: function collateralPrice() view returns(uint256 price)
func (_Stabilization *StabilizationCallerSession) CollateralPrice() (*big.Int, error) {
	return _Stabilization.Contract.CollateralPrice(&_Stabilization.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(uint256 borrowInterestRate, uint256 liquidationRatio, uint256 minCollateralizationRatio, uint256 minDebtRequirement, uint256 targetPrice)
func (_Stabilization *StabilizationCaller) Config(opts *bind.CallOpts) (struct {
	BorrowInterestRate        *big.Int
	LiquidationRatio          *big.Int
	MinCollateralizationRatio *big.Int
	MinDebtRequirement        *big.Int
	TargetPrice               *big.Int
}, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "config")

	outstruct := new(struct {
		BorrowInterestRate        *big.Int
		LiquidationRatio          *big.Int
		MinCollateralizationRatio *big.Int
		MinDebtRequirement        *big.Int
		TargetPrice               *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BorrowInterestRate = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LiquidationRatio = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MinCollateralizationRatio = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.MinDebtRequirement = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TargetPrice = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(uint256 borrowInterestRate, uint256 liquidationRatio, uint256 minCollateralizationRatio, uint256 minDebtRequirement, uint256 targetPrice)
func (_Stabilization *StabilizationSession) Config() (struct {
	BorrowInterestRate        *big.Int
	LiquidationRatio          *big.Int
	MinCollateralizationRatio *big.Int
	MinDebtRequirement        *big.Int
	TargetPrice               *big.Int
}, error) {
	return _Stabilization.Contract.Config(&_Stabilization.CallOpts)
}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns(uint256 borrowInterestRate, uint256 liquidationRatio, uint256 minCollateralizationRatio, uint256 minDebtRequirement, uint256 targetPrice)
func (_Stabilization *StabilizationCallerSession) Config() (struct {
	BorrowInterestRate        *big.Int
	LiquidationRatio          *big.Int
	MinCollateralizationRatio *big.Int
	MinDebtRequirement        *big.Int
	TargetPrice               *big.Int
}, error) {
	return _Stabilization.Contract.Config(&_Stabilization.CallOpts)
}

// DebtAmount is a free data retrieval call binding the contract method 0x50bf06bf.
//
// Solidity: function debtAmount(address account, uint256 timestamp) view returns(uint256 debt)
func (_Stabilization *StabilizationCaller) DebtAmount(opts *bind.CallOpts, account common.Address, timestamp *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "debtAmount", account, timestamp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DebtAmount is a free data retrieval call binding the contract method 0x50bf06bf.
//
// Solidity: function debtAmount(address account, uint256 timestamp) view returns(uint256 debt)
func (_Stabilization *StabilizationSession) DebtAmount(account common.Address, timestamp *big.Int) (*big.Int, error) {
	return _Stabilization.Contract.DebtAmount(&_Stabilization.CallOpts, account, timestamp)
}

// DebtAmount is a free data retrieval call binding the contract method 0x50bf06bf.
//
// Solidity: function debtAmount(address account, uint256 timestamp) view returns(uint256 debt)
func (_Stabilization *StabilizationCallerSession) DebtAmount(account common.Address, timestamp *big.Int) (*big.Int, error) {
	return _Stabilization.Contract.DebtAmount(&_Stabilization.CallOpts, account, timestamp)
}

// DebtAmount0 is a free data retrieval call binding the contract method 0x54a9f42c.
//
// Solidity: function debtAmount(address account) view returns(uint256 debt)
func (_Stabilization *StabilizationCaller) DebtAmount0(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "debtAmount0", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DebtAmount0 is a free data retrieval call binding the contract method 0x54a9f42c.
//
// Solidity: function debtAmount(address account) view returns(uint256 debt)
func (_Stabilization *StabilizationSession) DebtAmount0(account common.Address) (*big.Int, error) {
	return _Stabilization.Contract.DebtAmount0(&_Stabilization.CallOpts, account)
}

// DebtAmount0 is a free data retrieval call binding the contract method 0x54a9f42c.
//
// Solidity: function debtAmount(address account) view returns(uint256 debt)
func (_Stabilization *StabilizationCallerSession) DebtAmount0(account common.Address) (*big.Int, error) {
	return _Stabilization.Contract.DebtAmount0(&_Stabilization.CallOpts, account)
}

// InterestDue is a free data retrieval call binding the contract method 0x15184245.
//
// Solidity: function interestDue(uint256 debt, uint256 rate, uint256 timeBorrow, uint256 timeDue) pure returns(uint256)
func (_Stabilization *StabilizationCaller) InterestDue(opts *bind.CallOpts, debt *big.Int, rate *big.Int, timeBorrow *big.Int, timeDue *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "interestDue", debt, rate, timeBorrow, timeDue)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InterestDue is a free data retrieval call binding the contract method 0x15184245.
//
// Solidity: function interestDue(uint256 debt, uint256 rate, uint256 timeBorrow, uint256 timeDue) pure returns(uint256)
func (_Stabilization *StabilizationSession) InterestDue(debt *big.Int, rate *big.Int, timeBorrow *big.Int, timeDue *big.Int) (*big.Int, error) {
	return _Stabilization.Contract.InterestDue(&_Stabilization.CallOpts, debt, rate, timeBorrow, timeDue)
}

// InterestDue is a free data retrieval call binding the contract method 0x15184245.
//
// Solidity: function interestDue(uint256 debt, uint256 rate, uint256 timeBorrow, uint256 timeDue) pure returns(uint256)
func (_Stabilization *StabilizationCallerSession) InterestDue(debt *big.Int, rate *big.Int, timeBorrow *big.Int, timeDue *big.Int) (*big.Int, error) {
	return _Stabilization.Contract.InterestDue(&_Stabilization.CallOpts, debt, rate, timeBorrow, timeDue)
}

// IsLiquidatable is a free data retrieval call binding the contract method 0x042e02cf.
//
// Solidity: function isLiquidatable(address account) view returns(bool)
func (_Stabilization *StabilizationCaller) IsLiquidatable(opts *bind.CallOpts, account common.Address) (bool, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "isLiquidatable", account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsLiquidatable is a free data retrieval call binding the contract method 0x042e02cf.
//
// Solidity: function isLiquidatable(address account) view returns(bool)
func (_Stabilization *StabilizationSession) IsLiquidatable(account common.Address) (bool, error) {
	return _Stabilization.Contract.IsLiquidatable(&_Stabilization.CallOpts, account)
}

// IsLiquidatable is a free data retrieval call binding the contract method 0x042e02cf.
//
// Solidity: function isLiquidatable(address account) view returns(bool)
func (_Stabilization *StabilizationCallerSession) IsLiquidatable(account common.Address) (bool, error) {
	return _Stabilization.Contract.IsLiquidatable(&_Stabilization.CallOpts, account)
}

// MinimumCollateral is a free data retrieval call binding the contract method 0x08796696.
//
// Solidity: function minimumCollateral(uint256 principal, uint256 price, uint256 mcr) pure returns(uint256)
func (_Stabilization *StabilizationCaller) MinimumCollateral(opts *bind.CallOpts, principal *big.Int, price *big.Int, mcr *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "minimumCollateral", principal, price, mcr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimumCollateral is a free data retrieval call binding the contract method 0x08796696.
//
// Solidity: function minimumCollateral(uint256 principal, uint256 price, uint256 mcr) pure returns(uint256)
func (_Stabilization *StabilizationSession) MinimumCollateral(principal *big.Int, price *big.Int, mcr *big.Int) (*big.Int, error) {
	return _Stabilization.Contract.MinimumCollateral(&_Stabilization.CallOpts, principal, price, mcr)
}

// MinimumCollateral is a free data retrieval call binding the contract method 0x08796696.
//
// Solidity: function minimumCollateral(uint256 principal, uint256 price, uint256 mcr) pure returns(uint256)
func (_Stabilization *StabilizationCallerSession) MinimumCollateral(principal *big.Int, price *big.Int, mcr *big.Int) (*big.Int, error) {
	return _Stabilization.Contract.MinimumCollateral(&_Stabilization.CallOpts, principal, price, mcr)
}

// UnderCollateralized is a free data retrieval call binding the contract method 0xfbbe6991.
//
// Solidity: function underCollateralized(uint256 collateral, uint256 price, uint256 debt, uint256 liquidationRatio) pure returns(bool)
func (_Stabilization *StabilizationCaller) UnderCollateralized(opts *bind.CallOpts, collateral *big.Int, price *big.Int, debt *big.Int, liquidationRatio *big.Int) (bool, error) {
	var out []interface{}
	err := _Stabilization.contract.Call(opts, &out, "underCollateralized", collateral, price, debt, liquidationRatio)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UnderCollateralized is a free data retrieval call binding the contract method 0xfbbe6991.
//
// Solidity: function underCollateralized(uint256 collateral, uint256 price, uint256 debt, uint256 liquidationRatio) pure returns(bool)
func (_Stabilization *StabilizationSession) UnderCollateralized(collateral *big.Int, price *big.Int, debt *big.Int, liquidationRatio *big.Int) (bool, error) {
	return _Stabilization.Contract.UnderCollateralized(&_Stabilization.CallOpts, collateral, price, debt, liquidationRatio)
}

// UnderCollateralized is a free data retrieval call binding the contract method 0xfbbe6991.
//
// Solidity: function underCollateralized(uint256 collateral, uint256 price, uint256 debt, uint256 liquidationRatio) pure returns(bool)
func (_Stabilization *StabilizationCallerSession) UnderCollateralized(collateral *big.Int, price *big.Int, debt *big.Int, liquidationRatio *big.Int) (bool, error) {
	return _Stabilization.Contract.UnderCollateralized(&_Stabilization.CallOpts, collateral, price, debt, liquidationRatio)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_Stabilization *StabilizationTransactor) Borrow(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "borrow", amount)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_Stabilization *StabilizationSession) Borrow(amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.Borrow(&_Stabilization.TransactOpts, amount)
}

// Borrow is a paid mutator transaction binding the contract method 0xc5ebeaec.
//
// Solidity: function borrow(uint256 amount) returns()
func (_Stabilization *StabilizationTransactorSession) Borrow(amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.Borrow(&_Stabilization.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Stabilization *StabilizationTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Stabilization *StabilizationSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.Deposit(&_Stabilization.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Stabilization *StabilizationTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.Deposit(&_Stabilization.TransactOpts, amount)
}

// Liquidate is a paid mutator transaction binding the contract method 0x2f865568.
//
// Solidity: function liquidate(address account) payable returns()
func (_Stabilization *StabilizationTransactor) Liquidate(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "liquidate", account)
}

// Liquidate is a paid mutator transaction binding the contract method 0x2f865568.
//
// Solidity: function liquidate(address account) payable returns()
func (_Stabilization *StabilizationSession) Liquidate(account common.Address) (*types.Transaction, error) {
	return _Stabilization.Contract.Liquidate(&_Stabilization.TransactOpts, account)
}

// Liquidate is a paid mutator transaction binding the contract method 0x2f865568.
//
// Solidity: function liquidate(address account) payable returns()
func (_Stabilization *StabilizationTransactorSession) Liquidate(account common.Address) (*types.Transaction, error) {
	return _Stabilization.Contract.Liquidate(&_Stabilization.TransactOpts, account)
}

// Repay is a paid mutator transaction binding the contract method 0x402d8883.
//
// Solidity: function repay() payable returns()
func (_Stabilization *StabilizationTransactor) Repay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "repay")
}

// Repay is a paid mutator transaction binding the contract method 0x402d8883.
//
// Solidity: function repay() payable returns()
func (_Stabilization *StabilizationSession) Repay() (*types.Transaction, error) {
	return _Stabilization.Contract.Repay(&_Stabilization.TransactOpts)
}

// Repay is a paid mutator transaction binding the contract method 0x402d8883.
//
// Solidity: function repay() payable returns()
func (_Stabilization *StabilizationTransactorSession) Repay() (*types.Transaction, error) {
	return _Stabilization.Contract.Repay(&_Stabilization.TransactOpts)
}

// SetLiquidationRatio is a paid mutator transaction binding the contract method 0x946ce8cd.
//
// Solidity: function setLiquidationRatio(uint256 ratio) returns()
func (_Stabilization *StabilizationTransactor) SetLiquidationRatio(opts *bind.TransactOpts, ratio *big.Int) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "setLiquidationRatio", ratio)
}

// SetLiquidationRatio is a paid mutator transaction binding the contract method 0x946ce8cd.
//
// Solidity: function setLiquidationRatio(uint256 ratio) returns()
func (_Stabilization *StabilizationSession) SetLiquidationRatio(ratio *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.SetLiquidationRatio(&_Stabilization.TransactOpts, ratio)
}

// SetLiquidationRatio is a paid mutator transaction binding the contract method 0x946ce8cd.
//
// Solidity: function setLiquidationRatio(uint256 ratio) returns()
func (_Stabilization *StabilizationTransactorSession) SetLiquidationRatio(ratio *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.SetLiquidationRatio(&_Stabilization.TransactOpts, ratio)
}

// SetMinCollateralizationRatio is a paid mutator transaction binding the contract method 0x7b44646a.
//
// Solidity: function setMinCollateralizationRatio(uint256 ratio) returns()
func (_Stabilization *StabilizationTransactor) SetMinCollateralizationRatio(opts *bind.TransactOpts, ratio *big.Int) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "setMinCollateralizationRatio", ratio)
}

// SetMinCollateralizationRatio is a paid mutator transaction binding the contract method 0x7b44646a.
//
// Solidity: function setMinCollateralizationRatio(uint256 ratio) returns()
func (_Stabilization *StabilizationSession) SetMinCollateralizationRatio(ratio *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.SetMinCollateralizationRatio(&_Stabilization.TransactOpts, ratio)
}

// SetMinCollateralizationRatio is a paid mutator transaction binding the contract method 0x7b44646a.
//
// Solidity: function setMinCollateralizationRatio(uint256 ratio) returns()
func (_Stabilization *StabilizationTransactorSession) SetMinCollateralizationRatio(ratio *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.SetMinCollateralizationRatio(&_Stabilization.TransactOpts, ratio)
}

// SetMinDebtRequirement is a paid mutator transaction binding the contract method 0x53afe81d.
//
// Solidity: function setMinDebtRequirement(uint256 amount) returns()
func (_Stabilization *StabilizationTransactor) SetMinDebtRequirement(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "setMinDebtRequirement", amount)
}

// SetMinDebtRequirement is a paid mutator transaction binding the contract method 0x53afe81d.
//
// Solidity: function setMinDebtRequirement(uint256 amount) returns()
func (_Stabilization *StabilizationSession) SetMinDebtRequirement(amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.SetMinDebtRequirement(&_Stabilization.TransactOpts, amount)
}

// SetMinDebtRequirement is a paid mutator transaction binding the contract method 0x53afe81d.
//
// Solidity: function setMinDebtRequirement(uint256 amount) returns()
func (_Stabilization *StabilizationTransactorSession) SetMinDebtRequirement(amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.SetMinDebtRequirement(&_Stabilization.TransactOpts, amount)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_Stabilization *StabilizationTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "setOperator", operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_Stabilization *StabilizationSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _Stabilization.Contract.SetOperator(&_Stabilization.TransactOpts, operator)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_Stabilization *StabilizationTransactorSession) SetOperator(operator common.Address) (*types.Transaction, error) {
	return _Stabilization.Contract.SetOperator(&_Stabilization.TransactOpts, operator)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_Stabilization *StabilizationTransactor) SetOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "setOracle", oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_Stabilization *StabilizationSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _Stabilization.Contract.SetOracle(&_Stabilization.TransactOpts, oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_Stabilization *StabilizationTransactorSession) SetOracle(oracle common.Address) (*types.Transaction, error) {
	return _Stabilization.Contract.SetOracle(&_Stabilization.TransactOpts, oracle)
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_Stabilization *StabilizationTransactor) SetSupplyControl(opts *bind.TransactOpts, supplyControl common.Address) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "setSupplyControl", supplyControl)
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_Stabilization *StabilizationSession) SetSupplyControl(supplyControl common.Address) (*types.Transaction, error) {
	return _Stabilization.Contract.SetSupplyControl(&_Stabilization.TransactOpts, supplyControl)
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_Stabilization *StabilizationTransactorSession) SetSupplyControl(supplyControl common.Address) (*types.Transaction, error) {
	return _Stabilization.Contract.SetSupplyControl(&_Stabilization.TransactOpts, supplyControl)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_Stabilization *StabilizationTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_Stabilization *StabilizationSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.Withdraw(&_Stabilization.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_Stabilization *StabilizationTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _Stabilization.Contract.Withdraw(&_Stabilization.TransactOpts, amount)
}

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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Account    common.Address
	Liquidator common.Address
	Raw        types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
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
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
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
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
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
type UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradeableSession struct {
	Contract     *Upgradeable      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradeableCallerSession struct {
	Contract *UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// UpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradeableTransactorSession struct {
	Contract     *UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// UpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradeableRaw struct {
	Contract *Upgradeable // Generic contract binding to access the raw methods on
}

// UpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradeableCallerRaw struct {
	Contract *UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradeableTransactorRaw struct {
	Contract *UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeable creates a new instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeable(address common.Address, backend bind.ContractBackend) (*Upgradeable, error) {
	contract, err := bindUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Upgradeable{UpgradeableCaller: UpgradeableCaller{contract: contract}, UpgradeableTransactor: UpgradeableTransactor{contract: contract}, UpgradeableFilterer: UpgradeableFilterer{contract: contract}}, nil
}

// NewUpgradeableCaller creates a new read-only instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*UpgradeableCaller, error) {
	contract, err := bindUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeableCaller{contract: contract}, nil
}

// NewUpgradeableTransactor creates a new write-only instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeableTransactor, error) {
	contract, err := bindUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeableTransactor{contract: contract}, nil
}

// NewUpgradeableFilterer creates a new log filterer instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeableFilterer, error) {
	contract, err := bindUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeableFilterer{contract: contract}, nil
}

// bindUpgradeable binds a generic wrapper to an already deployed contract.
func bindUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Upgradeable *UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Upgradeable.Contract.UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Upgradeable *UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Upgradeable *UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Upgradeable *UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Upgradeable *UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Upgradeable *UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Upgradeable.Contract.contract.Transact(opts, method, params...)
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Upgradeable *UpgradeableCaller) GetNewContract(opts *bind.CallOpts) ([]byte, string, error) {
	var out []interface{}
	err := _Upgradeable.contract.Call(opts, &out, "getNewContract")

	if err != nil {
		return *new([]byte), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Upgradeable *UpgradeableSession) GetNewContract() ([]byte, string, error) {
	return _Upgradeable.Contract.GetNewContract(&_Upgradeable.CallOpts)
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Upgradeable *UpgradeableCallerSession) GetNewContract() ([]byte, string, error) {
	return _Upgradeable.Contract.GetNewContract(&_Upgradeable.CallOpts)
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactor) CompleteContractUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "completeContractUpgrade")
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Upgradeable *UpgradeableSession) CompleteContractUpgrade() (*types.Transaction, error) {
	return _Upgradeable.Contract.CompleteContractUpgrade(&_Upgradeable.TransactOpts)
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactorSession) CompleteContractUpgrade() (*types.Transaction, error) {
	return _Upgradeable.Contract.CompleteContractUpgrade(&_Upgradeable.TransactOpts)
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactor) ResetContractUpgrade(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "resetContractUpgrade")
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Upgradeable *UpgradeableSession) ResetContractUpgrade() (*types.Transaction, error) {
	return _Upgradeable.Contract.ResetContractUpgrade(&_Upgradeable.TransactOpts)
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Upgradeable *UpgradeableTransactorSession) ResetContractUpgrade() (*types.Transaction, error) {
	return _Upgradeable.Contract.ResetContractUpgrade(&_Upgradeable.TransactOpts)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Upgradeable *UpgradeableTransactor) UpgradeContract(opts *bind.TransactOpts, _bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Upgradeable.contract.Transact(opts, "upgradeContract", _bytecode, _abi)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Upgradeable *UpgradeableSession) UpgradeContract(_bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeContract(&_Upgradeable.TransactOpts, _bytecode, _abi)
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Upgradeable *UpgradeableTransactorSession) UpgradeContract(_bytecode []byte, _abi string) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeContract(&_Upgradeable.TransactOpts, _bytecode, _abi)
}
