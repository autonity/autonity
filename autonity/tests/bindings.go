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
	AccountabilityContract common.Address
	OracleContract         common.Address
	AcuContract            common.Address
	SupplyControlContract  common.Address
	StabilizationContract  common.Address
	UpgradeManagerContract common.Address
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
	Bin: "0x60806040523480156200001157600080fd5b50604051620018e6380380620018e683398101604081905262000034916200036d565b858580518251146200005957604051634ff799c560e01b815260040160405180910390fd5b60005b8151811015620000c0576001600160ff1b038282815181106200008357620000836200052c565b60200260200101511115620000ab57604051634ff799c560e01b815260040160405180910390fd5b80620000b78162000558565b9150506200005c565b508751620000d69060039060208b01906200014b565b508651620000ec9060049060208a0190620001a8565b506001869055620000ff86600a62000673565b6002555050600680546001600160a01b039485166001600160a01b03199182161790915560078054938516938216939093179092556008805491909316911617905550620007e3915050565b82805482825590600052602060002090810192821562000196579160200282015b8281111562000196578251829062000185908262000717565b50916020019190600101906200016c565b50620001a4929150620001f4565b5090565b828054828255906000526020600020908101928215620001e6579160200282015b82811115620001e6578251825591602001919060010190620001c9565b50620001a492915062000215565b80821115620001a45760006200020b82826200022c565b50600101620001f4565b5b80821115620001a4576000815560010162000216565b5080546200023a9062000688565b6000825580601f106200024b575050565b601f0160209004906000526020600020908101906200026b919062000215565b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620002af57620002af6200026e565b604052919050565b60006001600160401b03821115620002d357620002d36200026e565b5060051b60200190565b600082601f830112620002ef57600080fd5b81516020620003086200030283620002b7565b62000284565b82815260059290921b840181019181810190868411156200032857600080fd5b8286015b848110156200034557805183529183019183016200032c565b509695505050505050565b80516001600160a01b03811681146200036857600080fd5b919050565b60008060008060008060c087890312156200038757600080fd5b86516001600160401b038111156200039e57600080fd5b8701601f81018913620003b057600080fd5b8051620003c16200030282620002b7565b808282526020820191508b60208460051b8601011115620003e157600080fd5b602084015b60208460051b860101811015620004b85780516001600160401b038111156200040e57600080fd5b8d603f82880101126200042057600080fd5b858101602001516001600160401b038111156200044157620004416200026e565b62000456601f8201601f191660200162000284565b8181528f604083858b01010111156200046e57600080fd5b60005b828110156200049657604081858b010101516020828401015260208101905062000471565b50600060208383010152808652505050602083019250602081019050620003e6565b5060208b0151909950925050506001600160401b03811115620004da57600080fd5b620004e889828a01620002dd565b95505060408701519350620005006060880162000350565b9250620005106080880162000350565b91506200052060a0880162000350565b90509295509295509295565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016200056d576200056d62000542565b5060010190565b600181815b80851115620005b557816000190482111562000599576200059962000542565b80851615620005a757918102915b93841c939080029062000579565b509250929050565b600082620005ce575060016200066d565b81620005dd575060006200066d565b8160018114620005f65760028114620006015762000621565b60019150506200066d565b60ff84111562000615576200061562000542565b50506001821b6200066d565b5060208310610133831016604e8410600b841016171562000646575081810a6200066d565b62000652838362000574565b806000190482111562000669576200066962000542565b0290505b92915050565b6000620006818383620005bd565b9392505050565b600181811c908216806200069d57607f821691505b602082108103620006be57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200071257600081815260208120601f850160051c81016020861015620006ed5750805b601f850160051c820191505b818110156200070e57828155600101620006f9565b5050505b505050565b81516001600160401b038111156200073357620007336200026e565b6200074b8162000744845462000688565b84620006c4565b602080601f8311600181146200078357600084156200076a5750858301515b600019600386901b1c1916600185901b1785556200070e565b600085815260208120601f198616915b82811015620007b45788860151825594840194600190910190840162000793565b5085821015620007d35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6110f380620007f36000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c80637adbf973116100665780637adbf973146100fe578063a2e6204514610111578063b3ab15fb14610129578063d54d27991461013c578063f51e181a1461015157600080fd5b806307039ff9146100a3578063146ca531146100c15780633fa4f245146100d857806344b4708a146100e0578063683dd191146100f5575b600080fd5b6100ab61015a565b6040516100b8919061093f565b60405180910390f35b6100ca60005481565b6040519081526020016100b8565b6100ca610233565b6100f36100ee366004610a2f565b61025e565b005b6100ca60025481565b6100f361010c366004610b61565b61038a565b6101196103d6565b60405190151581526020016100b8565b6100f3610137366004610b61565b610707565b610144610753565b6040516100b89190610bc5565b6100ca60015481565b60606003805480602002602001604051908101604052809291908181526020016000905b8282101561022a57838290600052602060002001805461019d90610bd8565b80601f01602080910402602001604051908101604052809291908181526020018280546101c990610bd8565b80156102165780601f106101eb57610100808354040283529160200191610216565b820191906000526020600020905b8154815290600101906020018083116101f957829003601f168201915b50505050508152602001906001019061017e565b50505050905090565b6000805460000361025757604051631d3e00bb60e11b815260040160405180910390fd5b5060055490565b8282805182511461028257604051634ff799c560e01b815260040160405180910390fd5b60005b81518110156102e1576001600160ff1b038282815181106102a8576102a8610c12565b602002602001015111156102cf57604051634ff799c560e01b815260040160405180910390fd5b806102d981610c3e565b915050610285565b506007546001600160a01b0316331461030c576040516282b42960e81b815260040160405180910390fd5b845161031f9060039060208801906107ab565b508351610333906004906020870190610801565b50600183905561034483600a610d3d565b6002556040517fdbdcd10543a20811a4a332247f28d03b22686d3281043de35824a06075c06c099061037b90879087908790610d49565b60405180910390a15050505050565b6006546001600160a01b031633146103b4576040516282b42960e81b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0392909216919091179055565b6006546000906001600160a01b03163314610403576040516282b42960e81b815260040160405180910390fd5b60006001600860009054906101000a90046001600160a01b03166001600160a01b0316639f8743f76040518163ffffffff1660e01b8152600401602060405180830381865afa15801561045a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061047e9190610d7f565b6104889190610d98565b9050806000541061049b57600091505090565b600080600860009054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156104f1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105159190610d7f565b905060005b6003548110156106a257604051661554d10b5554d160ca1b6020820152600090602701604051602081830303815290604052805190602001206003838154811061056657610566610c12565b906000526020600020016040516020016105809190610dab565b60405160208183030381529060405280519060200120036105a2575081610659565b600854600380546000926001600160a01b031691633c8510fd91899190879081106105cf576105cf610c12565b906000526020600020016040518363ffffffff1660e01b81526004016105f6929190610e21565b608060405180830381865afa158015610613573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106379190610eb4565b90508060600151600014610652576000965050505050505090565b6020015190505b6004828154811061066c5761066c610c12565b9060005260206000200154816106829190610f1a565b61068c9085610f4a565b935050808061069a90610c3e565b91505061051a565b506106ad8183610f72565b60058190556000849055604080514381524260208201528082018690526060810192909252517f23f161ca67071b3e902d4fa7afade82672c6160677e89d373a830145bdda6d269181900360800190a16001935050505090565b6006546001600160a01b03163314610731576040516282b42960e81b815260040160405180910390fd5b600780546001600160a01b0319166001600160a01b0392909216919091179055565b606060048054806020026020016040519081016040528092919081815260200182805480156107a157602002820191906000526020600020905b81548152602001906001019080831161078d575b5050505050905090565b8280548282559060005260206000209081019282156107f1579160200282015b828111156107f157825182906107e19082610ffd565b50916020019190600101906107cb565b506107fd929150610848565b5090565b82805482825590600052602060002090810192821561083c579160200282015b8281111561083c578251825591602001919060010190610821565b506107fd929150610865565b808211156107fd57600061085c828261087a565b50600101610848565b5b808211156107fd5760008155600101610866565b50805461088690610bd8565b6000825580601f10610896575050565b601f0160209004906000526020600020908101906108b49190610865565b50565b600082825180855260208086019550808260051b8401018186016000805b8581101561093157601f1980888603018b5283518051808752845b8181101561090b578281018901518882018a015288016108f0565b5086810188018590529b87019b601f0190911690940185019350918401916001016108d5565b509198975050505050505050565b60208152600061095260208301846108b7565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561099857610998610959565b604052919050565b600067ffffffffffffffff8211156109ba576109ba610959565b5060051b60200190565b600082601f8301126109d557600080fd5b813560206109ea6109e5836109a0565b61096f565b82815260059290921b84018101918181019086841115610a0957600080fd5b8286015b84811015610a245780358352918301918301610a0d565b509695505050505050565b600080600060608486031215610a4457600080fd5b833567ffffffffffffffff80821115610a5c57600080fd5b818601915086601f830112610a7057600080fd5b81356020610a806109e5836109a0565b82815260059290921b8401810191818101908a841115610a9f57600080fd5b8286015b84811015610b2c57803586811115610abb5760008081fd5b8701603f81018d13610acd5760008081fd5b84810135604088821115610ae357610ae3610959565b610af5601f8301601f1916880161096f565b8281528f82848601011115610b0a5760008081fd5b8282850189830137600092810188019290925250845250918301918301610aa3565b5097505087013592505080821115610b4357600080fd5b50610b50868287016109c4565b925050604084013590509250925092565b600060208284031215610b7357600080fd5b81356001600160a01b038116811461095257600080fd5b600081518084526020808501945080840160005b83811015610bba57815187529582019590820190600101610b9e565b509495945050505050565b6020815260006109526020830184610b8a565b600181811c90821680610bec57607f821691505b602082108103610c0c57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201610c5057610c50610c28565b5060010190565b600181815b80851115610c92578160001904821115610c7857610c78610c28565b80851615610c8557918102915b93841c9390800290610c5c565b509250929050565b600082610ca957506001610d37565b81610cb657506000610d37565b8160018114610ccc5760028114610cd657610cf2565b6001915050610d37565b60ff841115610ce757610ce7610c28565b50506001821b610d37565b5060208310610133831016604e8410600b8410161715610d15575081810a610d37565b610d1f8383610c57565b8060001904821115610d3357610d33610c28565b0290505b92915050565b60006109528383610c9a565b606081526000610d5c60608301866108b7565b8281036020840152610d6e8186610b8a565b915050826040830152949350505050565b600060208284031215610d9157600080fd5b5051919050565b81810381811115610d3757610d37610c28565b6000808354610db981610bd8565b60018281168015610dd15760018114610de657610e15565b60ff1984168752821515830287019450610e15565b8760005260208060002060005b85811015610e0c5781548a820152908401908201610df3565b50505082870194505b50929695505050505050565b8281526000602060408184015260008454610e3b81610bd8565b8060408701526060600180841660008114610e5d5760018114610e7757610ea5565b60ff1985168984015283151560051b890183019550610ea5565b896000528660002060005b85811015610e9d5781548b8201860152908301908801610e82565b8a0184019650505b50939998505050505050505050565b600060808284031215610ec657600080fd5b6040516080810181811067ffffffffffffffff82111715610ee957610ee9610959565b8060405250825181526020830151602082015260408301516040820152606083015160608201528091505092915050565b80820260008212600160ff1b84141615610f3657610f36610c28565b8181058314821517610d3757610d37610c28565b8082018281126000831280158216821582161715610f6a57610f6a610c28565b505092915050565b600082610f8f57634e487b7160e01b600052601260045260246000fd5b600160ff1b821460001984141615610fa957610fa9610c28565b500590565b601f821115610ff857600081815260208120601f850160051c81016020861015610fd55750805b601f850160051c820191505b81811015610ff457828155600101610fe1565b5050505b505050565b815167ffffffffffffffff81111561101757611017610959565b61102b816110258454610bd8565b84610fae565b602080601f83116001811461106057600084156110485750858301515b600019600386901b1c1916600185901b178555610ff4565b600085815260208120601f198616915b8281101561108f57888601518255948401946001909101908401611070565b50858210156110ad5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea264697066735822122090a9b0d0854200a11135cda6f702d02a3d6c311a3a5d80dc2547b2edb0dca66864736f6c63430008150033",
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
	return address, gasConsumed, &ACU{contract: c}, nil
}

// ACU is an auto generated Go binding around an Ethereum contract.
type ACU struct {
	*contract
}

// Quantities is a free data retrieval call binding the contract method 0xd54d2799.
//
// Solidity: function quantities() view returns(uint256[])
func (_ACU *ACU) Quantities(opts *runOptions) ([]*big.Int, uint64, error) {
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
func (_ACU *ACU) Round(opts *runOptions) (*big.Int, uint64, error) {
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
func (_ACU *ACU) Scale(opts *runOptions) (*big.Int, uint64, error) {
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
func (_ACU *ACU) ScaleFactor(opts *runOptions) (*big.Int, uint64, error) {
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
func (_ACU *ACU) Symbols(opts *runOptions) ([]string, uint64, error) {
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
func (_ACU *ACU) Value(opts *runOptions) (*big.Int, uint64, error) {
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
func (_ACU *ACU) ModifyBasket(opts *runOptions, symbols_ []string, quantities_ []*big.Int, scale_ *big.Int) (uint64, error) {
	_, consumed, err := _ACU.call(opts, "modifyBasket", symbols_, quantities_, scale_)
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ACU *ACU) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _ACU.call(opts, "setOperator", operator)
	return consumed, err
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_ACU *ACU) SetOracle(opts *runOptions, oracle common.Address) (uint64, error) {
	_, consumed, err := _ACU.call(opts, "setOracle", oracle)
	return consumed, err
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_ACU *ACU) Update(opts *runOptions) (uint64, error) {
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
 		func (_ACU *ACU) FilterBasketModified(opts *bind.FilterOpts) (*ACUBasketModifiedIterator, error) {





			logs, sub, err := _ACU.contract.FilterLogs(opts, "BasketModified")
			if err != nil {
				return nil, err
			}
			return &ACUBasketModifiedIterator{contract: _ACU.contract, event: "BasketModified", logs: logs, sub: sub}, nil
 		}

		// WatchBasketModified is a free log subscription operation binding the contract event 0xdbdcd10543a20811a4a332247f28d03b22686d3281043de35824a06075c06c09.
		//
		// Solidity: event BasketModified(string[] symbols, uint256[] quantities, uint256 scale)
		func (_ACU *ACU) WatchBasketModified(opts *bind.WatchOpts, sink chan<- *ACUBasketModified) (event.Subscription, error) {





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
		func (_ACU *ACU) ParseBasketModified(log types.Log) (*ACUBasketModified, error) {
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
 		func (_ACU *ACU) FilterUpdated(opts *bind.FilterOpts) (*ACUUpdatedIterator, error) {






			logs, sub, err := _ACU.contract.FilterLogs(opts, "Updated")
			if err != nil {
				return nil, err
			}
			return &ACUUpdatedIterator{contract: _ACU.contract, event: "Updated", logs: logs, sub: sub}, nil
 		}

		// WatchUpdated is a free log subscription operation binding the contract event 0x23f161ca67071b3e902d4fa7afade82672c6160677e89d373a830145bdda6d26.
		//
		// Solidity: event Updated(uint256 height, uint256 timestamp, uint256 round, int256 value)
		func (_ACU *ACU) WatchUpdated(opts *bind.WatchOpts, sink chan<- *ACUUpdated) (event.Subscription, error) {






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
		func (_ACU *ACU) ParseUpdated(log types.Log) (*ACUUpdated, error) {
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
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_autonity\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"innocenceProofSubmissionWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateLow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateMid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collusionFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"historyFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashingRatePrecision\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"InnocenceProven\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewAccusation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewFaultProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"releaseBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isJailbound\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"SlashingEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"beneficiaries\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"_rule\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"canAccuse\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_result\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_deadline\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"_rule\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"canSlash\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"innocenceProofSubmissionWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateLow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseSlashingRateMid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"collusionFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"historyFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailFactor\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashingRatePrecision\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"events\",\"outputs\":[{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_epochEnd\",\"type\":\"bool\"}],\"name\":\"finalize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_val\",\"type\":\"address\"}],\"name\":\"getValidatorAccusation\",\"outputs\":[{\"components\":[{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_val\",\"type\":\"address\"}],\"name\":\"getValidatorFaults\",\"outputs\":[{\"components\":[{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event\",\"name\":\"_event\",\"type\":\"tuple\"}],\"name\":\"handleAccusation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event\",\"name\":\"_event\",\"type\":\"tuple\"}],\"name\":\"handleInnocenceProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumAccountability.EventType\",\"name\":\"eventType\",\"type\":\"uint8\"},{\"internalType\":\"enumAccountability.Rule\",\"name\":\"rule\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"reporter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"offender\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"rawProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reportingBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"messageHash\",\"type\":\"uint256\"}],\"internalType\":\"structAccountability.Event\",\"name\":\"_event\",\"type\":\"tuple\"}],\"name\":\"handleMisbehaviour\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_newCommittee\",\"type\":\"address[]\"}],\"name\":\"setCommittee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newPeriod\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"slashingHistory\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
		"1e85f2e5": "handleAccusation((uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256))",
		"e05f8786": "handleInnocenceProof((uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256))",
		"a48a9b50": "handleMisbehaviour((uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256))",
		"e08b14ed": "setCommittee(address[])",
		"6b5f444c": "setEpochPeriod(uint256)",
		"e7bb0b52": "slashingHistory(address,uint256)",
	},
	Bin: "0x608060405260006014553480156200001657600080fd5b5060405162004241380380620042418339810160408190526200003991620002e7565b600180546001600160a01b0319166001600160a01b03841690811790915560408051636fd8d26960e11b8152905163dfb1a4d2916004808201926020929091908290030181865afa15801562000093573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000b9919062000378565b6000908155600154604080516355c7b7ff60e11b815290516001600160a01b039092169163ab8f6ffe9160048082019286929091908290030181865afa15801562000108573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405262000132919081019062000392565b905060005b8151811015620001f857600e82828151811062000158576200015862000530565b6020908102919091018101515182546001808201855560009485529284200180546001600160a01b0319166001600160a01b0390921691909117905583519091600f91859085908110620001b057620001b062000530565b602090810291909101810151516001600160a01b03168252810191909152604001600020805460ff191691151591909117905580620001ef8162000546565b91505062000137565b50508051600355602081015160045560408101516005556060810151600655608081015160075560a081015160085560c00151600955506200056e565b6001600160a01b03811681146200024b57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60405160e081016001600160401b03811182821017156200028957620002896200024e565b60405290565b604051606081016001600160401b03811182821017156200028957620002896200024e565b604051601f8201601f191681016001600160401b0381118282101715620002df57620002df6200024e565b604052919050565b600080828403610100811215620002fd57600080fd5b83516200030a8162000235565b925060e0601f19820112156200031f57600080fd5b506200032a62000264565b6020840151815260408401516020820152606084015160408201526080840151606082015260a0840151608082015260c084015160a082015260e084015160c0820152809150509250929050565b6000602082840312156200038b57600080fd5b5051919050565b60006020808385031215620003a657600080fd5b82516001600160401b0380821115620003be57600080fd5b818501915085601f830112620003d357600080fd5b815181811115620003e857620003e86200024e565b8060051b620003f9858201620002b4565b91825283810185019185810190898411156200041457600080fd5b86860192505b838310156200052357825185811115620004345760008081fd5b86016060601f19828d0381018213156200044e5760008081fd5b620004586200028f565b8a840151620004678162000235565b81526040848101518c830152928401519289841115620004875760008081fd5b83850194508e603f8601126200049f57600093508384fd5b8b850151935089841115620004b857620004b86200024e565b620004ca8c84601f87011601620002b4565b92508383528e81858701011115620004e25760008081fd5b60005b8481101562000502578581018201518482018e01528c01620004e5565b5060009383018c01939093529182015283525091860191908601906200041a565b9998505050505050505050565b634e487b7160e01b600052603260045260246000fd5b6000600182016200056757634e487b7160e01b600052601160045260246000fd5b5060010190565b613cc3806200057e6000396000f3fe6080604052600436106100f35760003560e01c80637ccecadd1161008a578063bebaa8fc11610059578063bebaa8fc14610332578063e05f87861461035f578063e08b14ed1461037f578063e7bb0b521461039f57600080fd5b80637ccecadd1461028a5780639cb22b06146102c1578063a48a9b50146102ee578063b5b7a1841461030e57600080fd5b80634108a95a116100c65780634108a95a146101b65780636b5f444c146101e65780636c9789b01461020657806379502c551461022657600080fd5b806301567739146100f85780630b7914301461014b5780631de9d9b6146101815780631e85f2e514610196575b600080fd5b34801561010457600080fd5b5061012e610113366004612ff0565b600a602052600090815260409020546001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b34801561015757600080fd5b5061016b610166366004613014565b6103d7565b6040516101429a999897969594939291906130b7565b61019461018f366004612ff0565b6104d4565b005b3480156101a257600080fd5b506101946101b136600461325e565b6106e3565b3480156101c257600080fd5b506101d66101d1366004613346565b6107b1565b6040519015158152602001610142565b3480156101f257600080fd5b50610194610201366004613014565b6108d9565b34801561021257600080fd5b50610194610221366004613384565b610908565b34801561023257600080fd5b506003546004546005546006546007546008546009546102559695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e001610142565b34801561029657600080fd5b506102aa6102a5366004613346565b610948565b604080519215158352602083019190915201610142565b3480156102cd57600080fd5b506102e16102dc366004612ff0565b610b14565b6040516101429190613459565b3480156102fa57600080fd5b5061019461030936600461325e565b610d17565b34801561031a57600080fd5b5061032460005481565b604051908152602001610142565b34801561033e57600080fd5b5061035261034d366004612ff0565b610dec565b604051610142919061346c565b34801561036b57600080fd5b5061019461037a36600461325e565b611067565b34801561038b57600080fd5b5061019461039a3660046134ce565b61113e565b3480156103ab57600080fd5b506103246103ba366004613580565b601160209081526000928352604080842090915290825290205481565b600281815481106103e757600080fd5b600091825260209091206008909102018054600182015460028301805460ff808516965061010085041694620100009094046001600160a01b03908116949316929190610433906135ac565b80601f016020809104026020016040519081016040528092919081815260200182805461045f906135ac565b80156104ac5780601f10610481576101008083540402835291602001916104ac565b820191906000526020600020905b81548152906001019060200180831161048f57829003601f168201915b505050505090806003015490806004015490806005015490806006015490806007015490508a565b6001546001600160a01b031633146105075760405162461bcd60e51b81526004016104fe906135e6565b60405180910390fd5b6001546001600160a01b038281166000908152600a6020526040808220549051630c825d9760e11b8152908316600482015290929190911690631904bb2e90602401600060405180830381865afa158015610566573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261058e9190810190613691565b516040519091506000906001600160a01b038316906108fc90349084818181858888f193505050503d80600081146105e2576040519150601f19603f3d011682016040523d82523d6000602084013e6105e7565b606091505b50509050806106ba57600160009054906101000a90046001600160a01b03166001600160a01b031663f7866ee36040518163ffffffff1660e01b8152600401602060405180830381865afa158015610643573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610667919061380f565b6001600160a01b03163460405160006040518083038185875af1925050503d80600081146106b1576040519150601f19603f3d011682016040523d82523d6000602084013e6106b6565b606091505b5050505b50506001600160a01b03166000908152600a6020526040902080546001600160a01b0319169055565b336000908152600f602052604090205460ff1615156001146107175760405162461bcd60e51b81526004016104fe9061382c565b60408101516001600160a01b031633146107435760405162461bcd60e51b81526004016104fe90613875565b6001815160028111156107585761075861302d565b146107a55760405162461bcd60e51b815260206004820152601f60248201527f77726f6e67206576656e74207479706520666f722061636375736174696f6e0060448201526064016104fe565b6107ae81611277565b50565b6000808360098111156107c6576107c661302d565b101580156107e6575060098360098111156107e3576107e361302d565b11155b61082a5760405162461bcd60e51b81526020600482015260156024820152741c9d5b19481a59081b5d5cdd081899481d985b1a59605a1b60448201526064016104fe565b60006108358461157a565b6001546040516396b477cb60e01b8152600481018690529192506000916001600160a01b03909116906396b477cb90602401602060405180830381865afa158015610884573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108a891906138ac565b6001600160a01b03871660009081526011602090815260408083209383529290522054919091109150509392505050565b6001546001600160a01b031633146109035760405162461bcd60e51b81526004016104fe906135e6565b600055565b6001546001600160a01b031633146109325760405162461bcd60e51b81526004016104fe906135e6565b61093a6115e5565b80156107ae576107ae611917565b6000808084600981111561095e5761095e61302d565b1015801561097e5750600984600981111561097b5761097b61302d565b11155b6109c25760405162461bcd60e51b81526020600482015260156024820152741c9d5b19481a59081b5d5cdd081899481d985b1a59605a1b60448201526064016104fe565b60006109cd8561157a565b6001546040516396b477cb60e01b8152600481018790529192506000916001600160a01b03909116906396b477cb90602401602060405180830381865afa158015610a1c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a4091906138ac565b6001600160a01b03881660009081526011602090815260408083208484529091529020549091508211610a7a576000935060009250610b0a565b6001600160a01b0387166000908152600c602052604090205415610b01576001600160a01b0387166000908152600c6020526040812054600290610ac0906001906138db565b81548110610ad057610ad06138ee565b90600052602060002090600802019050600094506003600001548160040154610af99190613904565b935050610b0a565b60019350600092505b5050935093915050565b610b1c612e6f565b6001600160a01b0382166000908152600c6020526040902054610b715760405162461bcd60e51b815260206004820152600d60248201526c37379030b1b1bab9b0ba34b7b760991b60448201526064016104fe565b6001600160a01b0382166000908152600c6020526040902054600290610b99906001906138db565b81548110610ba957610ba96138ee565b600091825260209091206040805161014081019091526008909202018054829060ff166002811115610bdd57610bdd61302d565b6002811115610bee57610bee61302d565b81528154602090910190610100900460ff166009811115610c1157610c1161302d565b6009811115610c2257610c2261302d565b815281546001600160a01b0362010000909104811660208301526001830154166040820152600282018054606090920191610c5c906135ac565b80601f0160208091040260200160405190810160405280929190818152602001828054610c88906135ac565b8015610cd55780601f10610caa57610100808354040283529160200191610cd5565b820191906000526020600020905b815481529060010190602001808311610cb857829003601f168201915b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820154815250509050919050565b336000908152600f602052604090205460ff161515600114610d4b5760405162461bcd60e51b81526004016104fe9061382c565b60408101516001600160a01b03163314610d775760405162461bcd60e51b81526004016104fe90613875565b600081516002811115610d8c57610d8c61302d565b14610de35760405162461bcd60e51b815260206004820152602160248201527f77726f6e67206576656e74207479706520666f72206d69736265686176696f756044820152603960f91b60648201526084016104fe565b6107ae81611bd5565b6001600160a01b0381166000908152600b60205260408120546060919067ffffffffffffffff811115610e2157610e2161312d565b604051908082528060200260200182016040528015610e5a57816020015b610e47612e6f565b815260200190600190039081610e3f5790505b50905060005b6001600160a01b0384166000908152600b6020526040902054811015611060576001600160a01b0384166000908152600b6020526040902080546002919083908110610eae57610eae6138ee565b906000526020600020015481548110610ec957610ec96138ee565b600091825260209091206040805161014081019091526008909202018054829060ff166002811115610efd57610efd61302d565b6002811115610f0e57610f0e61302d565b81528154602090910190610100900460ff166009811115610f3157610f3161302d565b6009811115610f4257610f4261302d565b815281546001600160a01b0362010000909104811660208301526001830154166040820152600282018054606090920191610f7c906135ac565b80601f0160208091040260200160405190810160405280929190818152602001828054610fa8906135ac565b8015610ff55780601f10610fca57610100808354040283529160200191610ff5565b820191906000526020600020905b815481529060010190602001808311610fd857829003601f168201915b5050505050815260200160038201548152602001600482015481526020016005820154815260200160068201548152602001600782015481525050828281518110611042576110426138ee565b6020026020010181905250808061105890613917565b915050610e60565b5092915050565b336000908152600f602052604090205460ff16151560011461109b5760405162461bcd60e51b81526004016104fe9061382c565b60408101516001600160a01b031633146110c75760405162461bcd60e51b81526004016104fe90613875565b6002815160028111156110dc576110dc61302d565b146111355760405162461bcd60e51b8152602060048201526024808201527f77726f6e67206576656e74207479706520666f7220696e6e6f63656e636520706044820152633937b7b360e11b60648201526084016104fe565b6107ae81611dd9565b6001546001600160a01b031633146111685760405162461bcd60e51b81526004016104fe906135e6565b6111cb600d8054806020026020016040519081016040528092919081815260200182805480156111c157602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116111a3575b5050505050611f1f565b61122e600e80548060200260200160405190810160405280929190818152602001828054801561122457602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611206575b5050505050611fa8565b61123a600d6000612ed5565b600e805461124a91600d91612ef3565b5061125481611fa8565b611260600e6000612ed5565b805161127390600e906020840190612f43565b5050565b600080600080600061128e60fc876080015161203a565b94509450945094509450846112e55760405162461bcd60e51b815260206004820152601e60248201527f6661696c65642061636375736174696f6e20766572696669636174696f6e000060448201526064016104fe565b85606001516001600160a01b0316846001600160a01b03161461131a5760405162461bcd60e51b81526004016104fe90613930565b856020015160098111156113305761133061302d565b831461134e5760405162461bcd60e51b81526004016104fe9061395b565b6001546040516396b477cb60e01b8152600481018490526000916001600160a01b0316906396b477cb90602401602060405180830381865afa158015611398573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113bc91906138ac565b90506000600160009054906101000a90046001600160a01b03166001600160a01b031663c9d97af46040518163ffffffff1660e01b8152600401602060405180830381865afa158015611413573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061143791906138ac565b60c0890185905260e08901839052436101008a015261012089018490529050600061146460026001613904565b600e54611473906001906138db565b61147d9190613985565b9050818960e0015110156114b35761149760026001613904565b600d546114a6906001906138db565b6114b09190613985565b90505b3360009081526010602090815260408083208684529091529020548110156115365760405162461bcd60e51b815260206004820152603060248201527f76616c696461746f722073686f756c646e27742061627573652061636375736160448201526f0e8d2dedc40d2dc40c2dc40cae0dec6d60831b60648201526084016104fe565b61153f896120a1565b336000908152601060209081526040808320868452909152812080546001929061156a908490613904565b9091555050505050505050505050565b600060098260098111156115905761159061302d565b0361159e5760025b92915050565b60008260098111156115b2576115b261302d565b036115be576002611598565b60018260098111156115d2576115d261302d565b036115de576002611598565b6002611598565b6014545b60135481101561191257600060138281548110611608576116086138ee565b90600052602060002001549050806000036116235750611900565b61162e6001826138db565b9050600060028281548110611645576116456138ee565b600091825260209091206040805161014081019091526008909202018054829060ff1660028111156116795761167961302d565b600281111561168a5761168a61302d565b81528154602090910190610100900460ff1660098111156116ad576116ad61302d565b60098111156116be576116be61302d565b815281546001600160a01b03620100009091048116602083015260018301541660408201526002820180546060909201916116f8906135ac565b80601f0160208091040260200160405190810160405280929190818152602001828054611724906135ac565b80156117715780601f1061174657610100808354040283529160200191611771565b820191906000526020600020905b81548152906001019060200180831161175457829003601f168201915b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820154815250509050436003600001548261010001516117c59190613904565b11156117d2575050601455565b60608101516001600160a01b03166000908152600c6020908152604082208290558201516117ff9061157a565b60608301516001600160a01b0316600090815260116020908152604080832060e08701518452909152902054909150811161183c57505050611900565b6060820180516001600160a01b03908116600090815260116020908152604080832060e088015184528252808320869055845184168352600b8252808320805460018082018355918552838520018990556012805491820181559093527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec344490920187905592518151858152938401879052909116917f6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f910160405180910390a25050505b8061190a81613917565b9150506115e9565b601455565b600080600160009054906101000a90046001600160a01b03166001600160a01b031663c9d97af46040518163ffffffff1660e01b8152600401602060405180830381865afa15801561196d573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061199191906138ac565b905060005b601254811015611a0857816002601283815481106119b6576119b66138ee565b9060005260206000200154815481106119d1576119d16138ee565b906000526020600020906008020160050154036119f6576119f3600184613904565b92505b80611a0081613917565b915050611996565b5060005b601254811015611bc857611bb6600260128381548110611a2e57611a2e6138ee565b906000526020600020015481548110611a4957611a496138ee565b600091825260209091206040805161014081019091526008909202018054829060ff166002811115611a7d57611a7d61302d565b6002811115611a8e57611a8e61302d565b81528154602090910190610100900460ff166009811115611ab157611ab161302d565b6009811115611ac257611ac261302d565b815281546001600160a01b0362010000909104811660208301526001830154166040820152600282018054606090920191611afc906135ac565b80601f0160208091040260200160405190810160405280929190818152602001828054611b28906135ac565b8015611b755780601f10611b4a57610100808354040283529160200191611b75565b820191906000526020600020905b815481529060010190602001808311611b5857829003601f168201915b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820154815250508461232a565b80611bc081613917565b915050611a0c565b5061127360126000612ed5565b6000806000806000611bec60fe876080015161203a565b9450945094509450945084611c435760405162461bcd60e51b815260206004820152601960248201527f6661696c65642070726f6f6620766572696669636174696f6e0000000000000060448201526064016104fe565b85606001516001600160a01b0316846001600160a01b031614611c785760405162461bcd60e51b81526004016104fe90613930565b85602001516009811115611c8e57611c8e61302d565b8314611cac5760405162461bcd60e51b81526004016104fe9061395b565b438210611cf45760405162461bcd60e51b815260206004820152601660248201527563616e277420626520696e207468652066757475726560501b60448201526064016104fe565b60008211611d3a5760405162461bcd60e51b815260206004820152601360248201527263616e27742062652061742067656e6573697360681b60448201526064016104fe565b6001546040516396b477cb60e01b8152600481018490526000916001600160a01b0316906396b477cb90602401602060405180830381865afa158015611d84573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611da891906138ac565b60c0880184905260e088018190524361010089015261012088018390529050611dd0876128cc565b50505050505050565b6000806000806000611df060fd876080015161203a565b9450945094509450945084611e475760405162461bcd60e51b815260206004820152601d60248201527f6661696c656420696e6e6f63656e636520766572696669636174696f6e00000060448201526064016104fe565b85606001516001600160a01b0316846001600160a01b031614611e7c5760405162461bcd60e51b81526004016104fe90613930565b85602001516009811115611e9257611e9261302d565b8314611eb05760405162461bcd60e51b81526004016104fe9061395b565b438210611ef85760405162461bcd60e51b815260206004820152601660248201527563616e277420626520696e207468652066757475726560501b60448201526064016104fe565b60c08601829052610120860181905243610100870152611f1786612b0e565b505050505050565b6001546001600160a01b03163314611f495760405162461bcd60e51b81526004016104fe906135e6565b60005b815181101561127357600f6000838381518110611f6b57611f6b6138ee565b6020908102919091018101516001600160a01b03168252810191909152604001600020805460ff1916905580611fa081613917565b915050611f4c565b6001546001600160a01b03163314611fd25760405162461bcd60e51b81526004016104fe906135e6565b60005b8151811015611273576001600f6000848481518110611ff657611ff66138ee565b6020908102919091018101516001600160a01b03168252810191909152604001600020805460ff19169115159190911790558061203281613917565b915050611fd5565b600080600080600080865160206120519190613904565b905061205b612f98565b60a081838a8c5afa61206c57600080fd5b805160010361207a57600196505b602081015160408201516060830151608090930151989b919a509850909695509350505050565b60608101516001600160a01b03166000908152600c60205260409020541561210b5760405162461bcd60e51b815260206004820181905260248201527f616c72656164792070726f63657373696e6720616e2061636375736174696f6e60448201526064016104fe565b600061211a826020015161157a565b60608301516001600160a01b0316600090815260116020908152604080832060e0870151845290915290205490915081116121675760405162461bcd60e51b81526004016104fe9061399c565b6002805460a08401819052600180820183556000839052845160089092027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace0180548694919392849260ff199092169184908111156121c8576121c861302d565b021790555060208201518154829061ff0019166101008360098111156121f0576121f061302d565b02179055506040820151815462010000600160b01b031916620100006001600160a01b039283160217825560608301516001830180546001600160a01b031916919092161790556080820151600282019061224b9082613a2b565b5060a082810151600383015560c0830151600483015560e08301516005830155610100830151600683015561012090920151600790910155820151612291906001613904565b60608301516001600160a01b03166000908152600c602052604090205560a08201516013906122c1906001613904565b815460018101835560009283526020928390200155606083015160a084015160408051858152938401919091526001600160a01b03909116917f2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba4035191015b60405180910390a25050565b6001546060830151604051630c825d9760e11b81526001600160a01b0391821660048201526000929190911690631904bb2e90602401600060405180830381865afa15801561237d573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526123a59190810190613691565b60408481015160608601516001600160a01b039081166000908152600a6020529290922080546001600160a01b031916929091169190911790559050600381610260015160038111156123fa576123fa61302d565b0361240457505050565b600061241b612416856020015161157a565b612e16565b610220830151600754919250906000906124359083613985565b6006546124429087613985565b61244c9085613904565b6124569190613904565b60095490915081111561246857506009545b60008461012001518560c001518660a001516124849190613904565b61248e9190613904565b6009549091506000906124a18385613985565b6124ab9190613aeb565b90506000811180156124bc57508181145b156125ee57600060a087018190526101008701819052610120870181905260c08701526101e0860180518291906124f4908390613904565b905250610220860180516001919061250d908390613904565b905250600361026087015260006102008701526001546040516301adf0b760e51b81526001600160a01b03909116906335be16e090612550908990600401613b1d565b600060405180830381600087803b15801561256a57600080fd5b505af115801561257e573d6000803e3d6000fd5b5050505060208681015160a08a810151604080516001600160a01b03909416845293830185905260008385015260016060840152608083015291517f6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9929181900390910190a15050505050505050565b6101208601518190811161261b5780876101200181815161260f91906138db565b90525060009050612636565b61012087015161262b90826138db565b600061012089015290505b80156126b357808761010001511061267e5780876101000181815161265b91906138db565b90525060a0870180518291906126729083906138db565b905250600090506126b3565b61010087015161268e90826138db565b90508661010001518760a0018181516126a791906138db565b90525060006101008801525b6000811180156126d6575060008760a001518860c001516126d49190613904565b115b156127825760008760a001518860c001516126f19190613904565b60c08901516127009084613985565b61270a9190613aeb565b905060008860a001518960c001516127229190613904565b60a08a01516127319085613985565b61273b9190613aeb565b9050818960c00181815161274f91906138db565b90525060a0890180518291906127669083906138db565b9052506127738183613904565b61277d90846138db565b925050505b61278c81836138db565b915081876101e0018181516127a19190613904565b90525061022087018051600191906127ba908390613904565b9052506000546102208801516008546127d39190613985565b6127dd9190613985565b6127e79043613904565b61020088015260026102608801526001546040516301adf0b760e51b81526001600160a01b03909116906335be16e090612825908a90600401613b1d565b600060405180830381600087803b15801561283f57600080fd5b505af1158015612853573d6000803e3d6000fd5b5050506020808901516102008a015160a0808e0151604080516001600160a01b039095168552948401889052938301919091526000606083015260808201929092527f6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f992500160405180910390a1505050505050505050565b60006128db826020015161157a565b60608301516001600160a01b0316600090815260116020908152604080832060e0870151845290915290205490915081116129285760405162461bcd60e51b81526004016104fe9061399c565b6002805460a08401819052600180820183556000839052845160089092027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace0180548694919392849260ff199092169184908111156129895761298961302d565b021790555060208201518154829061ff0019166101008360098111156129b1576129b161302d565b02179055506040820151815462010000600160b01b031916620100006001600160a01b039283160217825560608301516001830180546001600160a01b0319169190921617905560808201516002820190612a0c9082613a2b565b5060a082810151600383015560c0830151600483015560e08084015160058401556101008401516006840155610120909301516007909201919091556060840180516001600160a01b039081166000908152600b602090815260408083209589018051875460018181018a559886528486200155805160128054988901815585527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec344490970196909655845184168352601182528083209689015183529590528490208590559051915192519116917f6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f9161231e91858252602082015260400190565b60608101516001600160a01b03166000908152600c602052604081205490819003612b7b5760405162461bcd60e51b815260206004820152601860248201527f6e6f206173736f6369617465642061636375736174696f6e000000000000000060448201526064016104fe565b81602001516009811115612b9157612b9161302d565b6002612b9e6001846138db565b81548110612bae57612bae6138ee565b6000918252602090912060089091020154610100900460ff166009811115612bd857612bd861302d565b14612c355760405162461bcd60e51b815260206004820152602760248201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e206044820152661c9d5b19481a5960ca1b60648201526084016104fe565b60c08201516002612c476001846138db565b81548110612c5757612c576138ee565b90600052602060002090600802016004015414612cc45760405162461bcd60e51b815260206004820152602560248201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e20604482015264626c6f636b60d81b60648201526084016104fe565b6101208201516002612cd76001846138db565b81548110612ce757612ce76138ee565b90600052602060002090600802016007015414612d525760405162461bcd60e51b8152602060048201526024808201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e206044820152630d0c2e6d60e31b60648201526084016104fe565b6014545b601354811015612dbc578160138281548110612d7457612d746138ee565b906000526020600020015403612daa57600060138281548110612d9957612d996138ee565b600091825260209091200155612dbc565b80612db481613917565b915050612d56565b506060820180516001600160a01b039081166000908152600c602090815260408083208390559351935191825292909116917f1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f910161231e565b600081612e2557505060055490565b60018203612e3557505060055490565b60028203612e4557505060055490565b60038203612e5557505060055490565b60048203612e665750612710919050565b50612710919050565b60408051610140810190915280600081526020016000815260200160006001600160a01b0316815260200160006001600160a01b031681526020016060815260200160008152602001600081526020016000815260200160008152602001600081525090565b50805460008255906000526020600020908101906107ae9190612fb6565b828054828255906000526020600020908101928215612f335760005260206000209182015b82811115612f33578254825591600101919060010190612f18565b50612f3f929150612fb6565b5090565b828054828255906000526020600020908101928215612f33579160200282015b82811115612f3357825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190612f63565b6040518060a001604052806005906020820280368337509192915050565b5b80821115612f3f5760008155600101612fb7565b6001600160a01b03811681146107ae57600080fd5b8035612feb81612fcb565b919050565b60006020828403121561300257600080fd5b813561300d81612fcb565b9392505050565b60006020828403121561302657600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b600381106130535761305361302d565b9052565b600a81106130535761305361302d565b60005b8381101561308257818101518382015260200161306a565b50506000910152565b600081518084526130a3816020860160208601613067565b601f01601f19169290920160200192915050565b60006101406130c6838e613043565b6130d3602084018d613057565b6001600160a01b038b811660408501528a166060840152608083018190526130fd8184018a61308b565b60a0840198909852505060c081019490945260e08401929092526101008301526101209091015295945050505050565b634e487b7160e01b600052604160045260246000fd5b604051610140810167ffffffffffffffff811182821017156131675761316761312d565b60405290565b604051610280810167ffffffffffffffff811182821017156131675761316761312d565b604051601f8201601f1916810167ffffffffffffffff811182821017156131ba576131ba61312d565b604052919050565b803560038110612feb57600080fd5b8035600a8110612feb57600080fd5b600067ffffffffffffffff8211156131fa576131fa61312d565b50601f01601f191660200190565b600082601f83011261321957600080fd5b813561322c613227826131e0565b613191565b81815284602083860101111561324157600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561327057600080fd5b813567ffffffffffffffff8082111561328857600080fd5b90830190610140828603121561329d57600080fd5b6132a5613143565b6132ae836131c2565b81526132bc602084016131d1565b60208201526132cd60408401612fe0565b60408201526132de60608401612fe0565b60608201526080830135828111156132f557600080fd5b61330187828601613208565b60808301525060a0838101359082015260c0808401359082015260e0808401359082015261010080840135908201526101209283013592810192909252509392505050565b60008060006060848603121561335b57600080fd5b833561336681612fcb565b9250613374602085016131d1565b9150604084013590509250925092565b60006020828403121561339657600080fd5b8135801515811461300d57600080fd5b60006101406133b6848451613043565b60208301516133c86020860182613057565b5060408301516133e360408601826001600160a01b03169052565b5060608301516133fe60608601826001600160a01b03169052565b5060808301518160808601526134168286018261308b565b91505060a083015160a085015260c083015160c085015260e083015160e08501526101008084015181860152506101208084015181860152508091505092915050565b60208152600061300d60208301846133a6565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156134c157603f198886030184526134af8583516133a6565b94509285019290850190600101613493565b5092979650505050505050565b600060208083850312156134e157600080fd5b823567ffffffffffffffff808211156134f957600080fd5b818501915085601f83011261350d57600080fd5b81358181111561351f5761351f61312d565b8060051b9150613530848301613191565b818152918301840191848101908884111561354a57600080fd5b938501935b83851015613574578435925061356483612fcb565b828252938501939085019061354f565b98975050505050505050565b6000806040838503121561359357600080fd5b823561359e81612fcb565b946020939093013593505050565b600181811c908216806135c057607f821691505b6020821081036135e057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526024908201527f66756e6374696f6e207265737472696374656420746f207468652076616c696460408201526330ba37b960e11b606082015260800190565b8051612feb81612fcb565b600082601f83011261364657600080fd5b8151613654613227826131e0565b81815284602083860101111561366957600080fd5b61367a826020830160208701613067565b949350505050565b805160048110612feb57600080fd5b6000602082840312156136a357600080fd5b815167ffffffffffffffff808211156136bb57600080fd5b9083019061028082860312156136d057600080fd5b6136d861316d565b6136e18361362a565b81526136ef6020840161362a565b60208201526137006040840161362a565b604082015260608301518281111561371757600080fd5b61372387828601613635565b6060830152506080830151608082015260a083015160a082015260c083015160c082015260e083015160e082015261010080840151818301525061012080840151818301525061014080840151818301525061016080840151818301525061018061378f81850161362a565b908201526101a083810151908201526101c080840151908201526101e080840151908201526102008084015190820152610220808401519082015261024080840151838111156137de57600080fd5b6137ea88828701613635565b8284015250506102609150613800828401613682565b91810191909152949350505050565b60006020828403121561382157600080fd5b815161300d81612fcb565b60208082526029908201527f66756e6374696f6e207265737472696374656420746f206120636f6d6d69747460408201526832b29036b2b6b132b960b91b606082015260800190565b6020808252601d908201527f6576656e74207265706f72746572206d7573742062652063616c6c6572000000604082015260600190565b6000602082840312156138be57600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b81810381811115611598576115986138c5565b634e487b7160e01b600052603260045260246000fd5b80820180821115611598576115986138c5565b600060018201613929576139296138c5565b5060010190565b6020808252601190820152700decccccadcc8cae440dad2e6dac2e8c6d607b1b604082015260600190565b60208082526010908201526f0e4ead8ca40d2c840dad2e6dac2e8c6d60831b604082015260600190565b8082028115828204841417611598576115986138c5565b60208082526024908201527f616c726561647920736c6173686564206174207468652070726f6f66277320656040820152630e0dec6d60e31b606082015260800190565b601f821115613a2657600081815260208120601f850160051c81016020861015613a075750805b601f850160051c820191505b81811015611f1757828155600101613a13565b505050565b815167ffffffffffffffff811115613a4557613a4561312d565b613a5981613a5384546135ac565b846139e0565b602080601f831160018114613a8e5760008415613a765750858301515b600019600386901b1c1916600185901b178555611f17565b600085815260208120601f198616915b82811015613abd57888601518255948401946001909101908401613a9e565b5085821015613adb5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b600082613b0857634e487b7160e01b600052601260045260246000fd5b500490565b600481106130535761305361302d565b60208152613b376020820183516001600160a01b03169052565b60006020830151613b5360408401826001600160a01b03169052565b5060408301516001600160a01b0381166060840152506060830151610280806080850152613b856102a085018361308b565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a0613c0c818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f190161026080880191909152909150613c6e848361308b565b935080870151915050613c8382860182613b0d565b509094935050505056fea26469706673582212201e36a873559d31abc0186c0e85d14a73ad8c0bc3ebc1016844095063c338bda264736f6c63430008150033",
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
	return address, gasConsumed, &Accountability{contract: c}, nil
}

// Accountability is an auto generated Go binding around an Ethereum contract.
type Accountability struct {
	*contract
}

// Beneficiaries is a free data retrieval call binding the contract method 0x01567739.
//
// Solidity: function beneficiaries(address ) view returns(address)
func (_Accountability *Accountability) Beneficiaries(opts *runOptions, arg0 common.Address) (common.Address, uint64, error) {
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
func (_Accountability *Accountability) CanAccuse(opts *runOptions, _offender common.Address, _rule uint8, _block *big.Int) (struct {
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
func (_Accountability *Accountability) CanSlash(opts *runOptions, _offender common.Address, _rule uint8, _block *big.Int) (bool, uint64, error) {
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
func (_Accountability *Accountability) Config(opts *runOptions) (struct {
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
func (_Accountability *Accountability) EpochPeriod(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "epochPeriod")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Events is a free data retrieval call binding the contract method 0x0b791430.
//
// Solidity: function events(uint256 ) view returns(uint8 eventType, uint8 rule, address reporter, address offender, bytes rawProof, uint256 id, uint256 block, uint256 epoch, uint256 reportingBlock, uint256 messageHash)
func (_Accountability *Accountability) Events(opts *runOptions, arg0 *big.Int) (struct {
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

	outstruct.EventType = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Rule = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Reporter = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Offender = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.RawProof = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.Id = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Block = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Epoch = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.ReportingBlock = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.MessageHash = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// GetValidatorAccusation is a free data retrieval call binding the contract method 0x9cb22b06.
//
// Solidity: function getValidatorAccusation(address _val) view returns((uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256))
func (_Accountability *Accountability) GetValidatorAccusation(opts *runOptions, _val common.Address) (AccountabilityEvent, uint64, error) {
	out, consumed, err := _Accountability.call(opts, "getValidatorAccusation", _val)

	if err != nil {
		return *new(AccountabilityEvent), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(AccountabilityEvent)).(*AccountabilityEvent)
	return out0, consumed, err

}

// GetValidatorFaults is a free data retrieval call binding the contract method 0xbebaa8fc.
//
// Solidity: function getValidatorFaults(address _val) view returns((uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256)[])
func (_Accountability *Accountability) GetValidatorFaults(opts *runOptions, _val common.Address) ([]AccountabilityEvent, uint64, error) {
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
func (_Accountability *Accountability) SlashingHistory(opts *runOptions, arg0 common.Address, arg1 *big.Int) (*big.Int, uint64, error) {
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
func (_Accountability *Accountability) DistributeRewards(opts *runOptions, _validator common.Address) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "distributeRewards", _validator)
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_Accountability *Accountability) Finalize(opts *runOptions, _epochEnd bool) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "finalize", _epochEnd)
	return consumed, err
}

// HandleAccusation is a paid mutator transaction binding the contract method 0x1e85f2e5.
//
// Solidity: function handleAccusation((uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256) _event) returns()
func (_Accountability *Accountability) HandleAccusation(opts *runOptions, _event AccountabilityEvent) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "handleAccusation", _event)
	return consumed, err
}

// HandleInnocenceProof is a paid mutator transaction binding the contract method 0xe05f8786.
//
// Solidity: function handleInnocenceProof((uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256) _event) returns()
func (_Accountability *Accountability) HandleInnocenceProof(opts *runOptions, _event AccountabilityEvent) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "handleInnocenceProof", _event)
	return consumed, err
}

// HandleMisbehaviour is a paid mutator transaction binding the contract method 0xa48a9b50.
//
// Solidity: function handleMisbehaviour((uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256,uint256) _event) returns()
func (_Accountability *Accountability) HandleMisbehaviour(opts *runOptions, _event AccountabilityEvent) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "handleMisbehaviour", _event)
	return consumed, err
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe08b14ed.
//
// Solidity: function setCommittee(address[] _newCommittee) returns()
func (_Accountability *Accountability) SetCommittee(opts *runOptions, _newCommittee []common.Address) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "setCommittee", _newCommittee)
	return consumed, err
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_Accountability *Accountability) SetEpochPeriod(opts *runOptions, _newPeriod *big.Int) (uint64, error) {
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
 		func (_Accountability *Accountability) FilterInnocenceProven(opts *bind.FilterOpts, _offender []common.Address) (*AccountabilityInnocenceProvenIterator, error) {

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
		func (_Accountability *Accountability) WatchInnocenceProven(opts *bind.WatchOpts, sink chan<- *AccountabilityInnocenceProven, _offender []common.Address) (event.Subscription, error) {

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
		func (_Accountability *Accountability) ParseInnocenceProven(log types.Log) (*AccountabilityInnocenceProven, error) {
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
 		func (_Accountability *Accountability) FilterNewAccusation(opts *bind.FilterOpts, _offender []common.Address) (*AccountabilityNewAccusationIterator, error) {

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
		func (_Accountability *Accountability) WatchNewAccusation(opts *bind.WatchOpts, sink chan<- *AccountabilityNewAccusation, _offender []common.Address) (event.Subscription, error) {

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
		func (_Accountability *Accountability) ParseNewAccusation(log types.Log) (*AccountabilityNewAccusation, error) {
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
 		func (_Accountability *Accountability) FilterNewFaultProof(opts *bind.FilterOpts, _offender []common.Address) (*AccountabilityNewFaultProofIterator, error) {

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
		func (_Accountability *Accountability) WatchNewFaultProof(opts *bind.WatchOpts, sink chan<- *AccountabilityNewFaultProof, _offender []common.Address) (event.Subscription, error) {

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
		func (_Accountability *Accountability) ParseNewFaultProof(log types.Log) (*AccountabilityNewFaultProof, error) {
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
 		func (_Accountability *Accountability) FilterSlashingEvent(opts *bind.FilterOpts) (*AccountabilitySlashingEventIterator, error) {







			logs, sub, err := _Accountability.contract.FilterLogs(opts, "SlashingEvent")
			if err != nil {
				return nil, err
			}
			return &AccountabilitySlashingEventIterator{contract: _Accountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
 		}

		// WatchSlashingEvent is a free log subscription operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
		func (_Accountability *Accountability) WatchSlashingEvent(opts *bind.WatchOpts, sink chan<- *AccountabilitySlashingEvent) (event.Subscription, error) {







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
		func (_Accountability *Accountability) ParseSlashingEvent(log types.Log) (*AccountabilitySlashingEvent, error) {
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator[]\",\"name\":\"_validators\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"ActivatedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"BondingRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewUnbondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeInitialization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getEpochFromBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_signatures\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAccountability\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAccountabilityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIACU\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAcuContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setOracleContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIStabilization\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setStabilizationContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISupplyControl\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSupplyControlContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractUpgradeManager\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setUpgradeManagerContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"_val\",\"type\":\"tuple\"}],\"name\":\"updateValidatorAndTransferSlashedFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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
		"9bb851c0": "totalRedistributed()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"a5d059ca": "unbond(address,uint256)",
		"35be16e0": "updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8))",
		"b2ea9adb": "upgradeContract(bytes,string)",
	},
	Bin: "0x60806040526000600b556000600c553480156200001b57600080fd5b506040516200a2303803806200a2308339810160408190526200003e9162000e64565b601c546000036200006757602a80546001600160a01b031916331790556200006782826200006f565b505062001325565b80518051600d55602080820151600e55604080830151600f55606080840151601055608093840151601180546001600160a01b03199081166001600160a01b03938416179091558487015180516012805484169185169190911790558086015160138054841691851691909117905580850151601480548416918516919091179055808401516015805484169185169190911790559586015160168054831691841691909117905560a0909501516017805487169183169190911790558286015180516018805490971692169190911790945591830151601955820151601a5590810151601b55810151601c5560005b825181101562000420576000838281518110620001805762000180620010a1565b602002602001015160a0015190506000848381518110620001a557620001a5620010a1565b60200260200101516101a00181815250506000848381518110620001cd57620001cd620010a1565b602002602001015161018001906001600160a01b031690816001600160a01b0316815250506000848381518110620002095762000209620010a1565b602002602001015160a00181815250506000848381518110620002305762000230620010a1565b60209081029190910101516101c00152600f5484518590849081106200025a576200025a620010a1565b602002602001015160800181815250506000848381518110620002815762000281620010a1565b602002602001015161026001906003811115620002a257620002a2620010b7565b90816003811115620002b857620002b8620010b7565b815250506000848381518110620002d357620002d3620010a1565b602002602001015161016001818152505062000311848381518110620002fd57620002fd620010a1565b60200260200101516200042560201b60201c565b6200033e8483815181106200032a576200032a620010a1565b60200260200101516200056060201b60201c565b8060276000868581518110620003585762000358620010a1565b6020026020010151600001516001600160a01b03166001600160a01b031681526020019081526020016000206000828254620003959190620010e3565b925050819055508060296000828254620003b09190620010e3565b925050819055506200040a848381518110620003d057620003d0620010a1565b60200260200101516020015182868581518110620003f257620003f2620010a1565b6020026020010151600001516200079160201b60201c565b50806200041781620010ff565b9150506200015f565b505050565b60006200043c82606001516200097c60201b60201c565b6001600160a01b03909116602084015290508015620004905760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b60448201526064015b60405180910390fd5b6020808301516001600160a01b03908116600090815260289092526040909120600101541615620005045760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000487565b612710826080015111156200055c5760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000487565b5050565b6101808101516001600160a01b0316620005e357601d546000906200058590620009ca565b905081602001518260000151836080015183604051620005a59062000aea565b620005b494939291906200111b565b604051809103906000f080158015620005d1573d6000803e3d6000fd5b506001600160a01b0316610180830152505b60208181018051601d80546001808201835560009283527f6d4407e7be21f808e6509aa9fa9143369579dd7d760fe20a2c09680fc146134f90910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602890955260409182902086518154871690851617815593519084018054861691841691909117905584015160028301805490941691161790915560608201518291906003820190620006959082620011ff565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f8201556102008201516010820155610220820151601182015561024082015160128201906200075f9082620011ff565b5061026082015160138201805460ff19166001836003811115620007875762000787620010b7565b0217905550505050565b60008211620007ef5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000487565b6001600160a01b038116600090815260276020526040902054821115620008595760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000487565b6001600160a01b0381166000908152602760205260408120805484929062000883908490620012cb565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580546000908152600394859052978820875181549088166001600160a01b03199182161782559551600182018054919098169616959095179095559051600284015551910155805491926200090f83620010ff565b90915550506001600160a01b03848116600081815260286020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d910160405180910390a35050505050565b6000806200098962000af8565b60008060ff9050604083875160208901845afa620009a657600080fd5b505080516020909101516c0100000000000000000000000090910494909350915050565b606081600003620009f25750506040805180820190915260018152600360fc1b602082015290565b8160005b811562000a22578062000a0981620010ff565b915062000a1a9050600a83620012f7565b9150620009f6565b6000816001600160401b0381111562000a3f5762000a3f62000b16565b6040519080825280601f01601f19166020018201604052801562000a6a576020820181803683370190505b5090505b841562000ae25762000a82600183620012cb565b915062000a91600a866200130e565b62000a9e906030620010e3565b60f81b81838151811062000ab65762000ab6620010a1565b60200101906001600160f81b031916908160001a90535062000ada600a86620012f7565b945062000a6e565b949350505050565b6115318062008cff83390190565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b038111828210171562000b515762000b5162000b16565b60405290565b60405160a081016001600160401b038111828210171562000b515762000b5162000b16565b60405161028081016001600160401b038111828210171562000b515762000b5162000b16565b604051601f8201601f191681016001600160401b038111828210171562000bcd5762000bcd62000b16565b604052919050565b6001600160a01b038116811462000beb57600080fd5b50565b805162000bfb8162000bd5565b919050565b60005b8381101562000c1d57818101518382015260200162000c03565b50506000910152565b600082601f83011262000c3857600080fd5b81516001600160401b0381111562000c545762000c5462000b16565b62000c69601f8201601f191660200162000ba2565b81815284602083860101111562000c7f57600080fd5b62000ae282602083016020870162000c00565b80516004811062000bfb57600080fd5b600060c0828403121562000cb557600080fd5b60405160c081016001600160401b038111828210171562000cda5762000cda62000b16565b8060405250809150825162000cef8162000bd5565b8152602083015162000d018162000bd5565b6020820152604083015162000d168162000bd5565b6040820152606083015162000d2b8162000bd5565b6060820152608083015162000d408162000bd5565b608082015260a083015162000d558162000bd5565b60a0919091015292915050565b60006080828403121562000d7557600080fd5b62000d7f62000b2c565b9050815162000d8e8162000bd5565b8082525060208201516020820152604082015160408201526060820151606082015292915050565b600081830361020081121562000dcb57600080fd5b62000dd562000b2c565b915060a081121562000de657600080fd5b5062000df162000b57565b82518152602083015160208201526040830151604082015260608301516060820152608083015162000e238162000bd5565b6080820152815262000e398360a0840162000ca2565b602082015262000e4e83610160840162000d62565b60408201526101e0820151606082015292915050565b60008061022080848603121562000e7a57600080fd5b83516001600160401b038082111562000e9257600080fd5b818601915086601f83011262000ea757600080fd5b815160208282111562000ebe5762000ebe62000b16565b8160051b62000ecf82820162000ba2565b928352848101820192828101908b85111562000eea57600080fd5b83870192505b848310156200107f5782518681111562000f0957600080fd5b8701610280818e03601f1901121562000f2157600080fd5b62000f2b62000b7c565b62000f3886830162000bee565b815262000f486040830162000bee565b8682015262000f5a6060830162000bee565b604082015260808201518881111562000f7257600080fd5b62000f828f888386010162000c26565b60608301525060a0820151608082015260c082015160a082015260e082015160c082015261010082015160e082015261012082015161010082015261014082015161012082015261016082015161014082015261018082015161016082015262000ff06101a0830162000bee565b6101808201526101c08201516101a08201526101e08201516101c08201526102008201516101e0820152898201516102008201526102408201518a820152610260820151888111156200104257600080fd5b620010528f888386010162000c26565b6102408301525062001068610280830162000c92565b610260820152835250918301919083019062000ef0565b8099505050506200109389828a0162000db6565b955050505050509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052601160045260246000fd5b80820180821115620010f957620010f9620010cd565b92915050565b600060018201620011145762001114620010cd565b5060010190565b600060018060a01b0380871683528086166020840152508360408301526080606083015282518060808401526200115a8160a085016020870162000c00565b601f01601f19169190910160a00195945050505050565b600181811c908216806200118657607f821691505b602082108103620011a757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200042057600081815260208120601f850160051c81016020861015620011d65750805b601f850160051c820191505b81811015620011f757828155600101620011e2565b505050505050565b81516001600160401b038111156200121b576200121b62000b16565b62001233816200122c845462001171565b84620011ad565b602080601f8311600181146200126b5760008415620012525750858301515b600019600386901b1c1916600185901b178555620011f7565b600085815260208120601f198616915b828110156200129c578886015182559484019460019091019084016200127b565b5085821015620012bb5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b81810381811115620010f957620010f9620010cd565b634e487b7160e01b600052601260045260246000fd5b600082620013095762001309620012e1565b500490565b600082620013205762001320620012e1565b500690565b6179ca80620013356000396000f3fe608060405260043610620003db5760003560e01c8063872cf0591162000203578063b66b3e791162000117578063d372c07e11620000a7578063dd62ed3e1162000075578063dd62ed3e1462000ccc578063dfb1a4d21462000d16578063e7f43c681462000d2d578063f7866ee31462000d4d57005b8063d372c07e1462000c48578063d5f394881462000c6d578063d861b0e81462000c8f578063d886f8a21462000ca757005b8063cb696f5411620000e5578063cb696f541462000bc1578063ceaad4551462000be6578063cf9c57191462000c0b578063cfd19fb91462000c2357005b8063b66b3e791462000b42578063b7ab4db51462000b6a578063c2362dd51462000b91578063c9d97af41462000ba957005b8063a5d059ca1162000193578063ae1f5fa01162000161578063ae1f5fa01462000aab578063b2ea9adb1462000ad3578063b3ecbadd1462000af8578063b46e55201462000b1d57005b8063a5d059ca1462000a13578063a8b2216e1462000a38578063a9059cbb1462000a5f578063ab8f6ffe1462000a8457005b80639bb851c011620001d15780639bb851c014620009995780639c98e47114620009b15780639dc29fac14620009c9578063a515366a14620009ee57005b8063872cf05914620008fd5780638bac7dad146200091557806395d89b41146200093a57806396b477cb146200096857005b806340c10f1911620002fb57806370a08231116200028b578063819b64631162000259578063819b6463146200087c578063833b1fce146200089357806384467fdb14620008b3578063852c484914620008d857005b806370a082311462000728578063731b3a03146200076257806377e741c7146200077957806379502c55146200079e57005b8063520fdbbc11620002c9578063520fdbbc14620006895780635f7d394914620006ae5780636b5f444c14620006ec5780636fd2c80b146200071157005b806340c10f191462000600578063436459691462000625578063496ccd9b146200063c5780634bb278f3146200066157005b80631604e416116200037757806329070c6d116200034557806329070c6d146200058e5780632f2c3f2e14620005a5578063313ce56714620005bd57806335be16e014620005db57005b80631604e416146200050657806318160ddd146200051e5780631904bb2e146200053557806323b872dd146200056957005b80630d8e6e2c11620003b55780630d8e6e2c14620004845780631122063314620004a5578063114eaf5514620004bc5780631250a28d14620004e157005b806306fdde0314620003e5578063095ea7b314620004295780630ae65e7a146200045f57005b36620003e357005b005b348015620003f257600080fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b604051620004209190620053f5565b60405180910390f35b3480156200043657600080fd5b506200044e6200044836600462005427565b62000d6d565b604051901515815260200162000420565b3480156200046c57600080fd5b50620003e36200047e36600462005456565b62000d86565b3480156200049157600080fd5b50601c545b60405190815260200162000420565b348015620004b257600080fd5b50600e5462000496565b348015620004c957600080fd5b50620003e3620004db36600462005476565b62000e19565b348015620004ee57600080fd5b50620003e36200050036600462005456565b62000e4b565b3480156200051357600080fd5b506200049660245481565b3480156200052b57600080fd5b5060295462000496565b3480156200054257600080fd5b506200055a6200055436600462005456565b62000e9a565b604051620004209190620054c9565b3480156200057657600080fd5b506200044e6200058836600462005642565b62001148565b3480156200059b57600080fd5b50600d5462000496565b348015620005b257600080fd5b506200049661271081565b348015620005ca57600080fd5b506040516012815260200162000420565b348015620005e857600080fd5b50620003e3620005fa36600462005688565b620011a2565b3480156200060d57600080fd5b50620003e36200061f36600462005427565b6200137e565b3480156200063257600080fd5b50601a5462000496565b3480156200064957600080fd5b50620003e36200065b36600462005456565b62001438565b3480156200066e57600080fd5b506200067962001542565b6040516200042092919062005749565b3480156200069657600080fd5b50620003e3620006a836600462005456565b62001974565b348015620006bb57600080fd5b50620006d3620006cd36600462005766565b62001b6f565b6040516001600160a01b03909116815260200162000420565b348015620006f957600080fd5b50620003e36200070b36600462005476565b62001d86565b3480156200071e57600080fd5b5060105462000496565b3480156200073557600080fd5b50620004966200074736600462005456565b6001600160a01b031660009081526027602052604090205490565b3480156200076f57600080fd5b5060205462000496565b3480156200078657600080fd5b50620003e36200079836600462005476565b62001f0c565b348015620007ab57600080fd5b506040805160a08082018352600d548252600e54602080840191909152600f54838501526010546060808501919091526011546001600160a01b03908116608080870191909152865160c0810188526012548316815260135483168186015260145483168189015260155483168185015260165483168183015260175483169581019590955286519081018752601854909116815260195492810192909252601a5494820194909452601b5493810193909352601c546200086a939084565b60405162000420949392919062005789565b3480156200088957600080fd5b50601b5462000496565b348015620008a057600080fd5b506013546001600160a01b0316620006d3565b348015620008c057600080fd5b50620003e3620008d236600462005903565b62001f3e565b348015620008e557600080fd5b50620003e3620008f736600462005427565b6200206d565b3480156200090a57600080fd5b50620003e3620021fe565b3480156200092257600080fd5b50620003e36200093436600462005476565b6200223a565b3480156200094757600080fd5b50604080518082019091526003815262272a2760e91b602082015262000411565b3480156200097557600080fd5b50620004966200098736600462005476565b6000908152601f602052604090205490565b348015620009a657600080fd5b506200049660235481565b348015620009be57600080fd5b506200049660215481565b348015620009d657600080fd5b50620003e3620009e836600462005427565b620022be565b348015620009fb57600080fd5b50620003e362000a0d36600462005427565b620023d4565b34801562000a2057600080fd5b50620003e362000a3236600462005427565b620024a7565b34801562000a4557600080fd5b5062000a506200253e565b604051620004209190620059ac565b34801562000a6c57600080fd5b506200044e62000a7e36600462005427565b62002621565b34801562000a9157600080fd5b5062000a9c62002630565b60405162000420919062005a12565b34801562000ab857600080fd5b5062000ac362002742565b6040516200042092919062005a6d565b34801562000ae057600080fd5b50620003e362000af236600462005a9f565b62002a0e565b34801562000b0557600080fd5b50620003e362000b1736600462005456565b62002a55565b34801562000b2a57600080fd5b50620003e362000b3c36600462005456565b62002aa4565b34801562000b4f57600080fd5b5062000b5a62002d27565b6040516200042092919062005b0a565b34801562000b7757600080fd5b5062000b8262002e5e565b60405162000420919062005b33565b34801562000b9e57600080fd5b506200049660205481565b34801562000bb657600080fd5b5062000496601e5481565b34801562000bce57600080fd5b50620003e362000be036600462005476565b62002ec2565b34801562000bf357600080fd5b50620003e362000c0536600462005456565b62002f25565b34801562000c1857600080fd5b50620003e362002f74565b34801562000c3057600080fd5b50620003e362000c4236600462005456565b62002fc8565b34801562000c5557600080fd5b50620003e362000c6736600462005456565b62003017565b34801562000c7a57600080fd5b50602a54620006d3906001600160a01b031681565b34801562000c9c57600080fd5b50620003e362003066565b34801562000cb457600080fd5b50620003e362000cc636600462005456565b620030a7565b34801562000cd957600080fd5b506200049662000ceb36600462005b48565b6001600160a01b03918216600090815260266020908152604080832093909416825291909152205490565b34801562000d2357600080fd5b5060195462000496565b34801562000d3a57600080fd5b506018546001600160a01b0316620006d3565b34801562000d5a57600080fd5b506011546001600160a01b0316620006d3565b600062000d7c338484620030f6565b5060015b92915050565b6001600160a01b038082166000818152602860205260409020600101549091161462000dcf5760405162461bcd60e51b815260040162000dc69062005b86565b60405180910390fd5b6001600160a01b0381811660009081526028602052604090205416331462000e0b5760405162461bcd60e51b815260040162000dc69062005bbd565b62000e16816200321f565b50565b6018546001600160a01b0316331462000e465760405162461bcd60e51b815260040162000dc69062005c09565b601055565b6018546001600160a01b0316331462000e785760405162461bcd60e51b815260040162000dc69062005c09565b601280546001600160a01b0319166001600160a01b0392909216919091179055565b62000ea4620051ce565b6001600160a01b038083166000818152602860205260409020600101549091161462000ee45760405162461bcd60e51b815260040162000dc69062005c40565b6001600160a01b03808316600090815260286020908152604091829020825161028081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000f479062005c77565b80601f016020809104026020016040519081016040528092919081815260200182805462000f759062005c77565b801562000fc65780601f1062000f9a5761010080835404028352916020019162000fc6565b820191906000526020600020905b81548152906001019060200180831162000fa857829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546080820152600882015460a0820152600982015460c0820152600a82015460e0820152600b820154610100820152600c8201546001600160a01b0316610120820152600d820154610140820152600e820154610160820152600f82015461018082015260108201546101a082015260118201546101c08201526012820180546101e090920191620010829062005c77565b80601f0160208091040260200160405190810160405280929190818152602001828054620010b09062005c77565b8015620011015780601f10620010d55761010080835404028352916020019162001101565b820191906000526020600020905b815481529060010190602001808311620010e357829003601f168201915b5050509183525050601382015460209091019060ff1660038111156200112b576200112b62005490565b60038111156200113f576200113f62005490565b90525092915050565b600062001157848484620032f6565b6001600160a01b03841660009081526026602090815260408083203384529091528120546200118890849062005cc9565b905062001197853383620030f6565b506001949350505050565b6012546001600160a01b031633146200120a5760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f742074686520736c617368696e6720636f6e74726044820152621858dd60ea1b606482015260840162000dc6565b600061012082013560288262001227604086016020870162005456565b6001600160a01b03166001600160a01b031681526020019081526020016000206009015462001257919062005cc9565b60c08301356028600062001272604087016020880162005456565b6001600160a01b03166001600160a01b0316815260200190815260200160002060060154620012a2919062005cc9565b60a084013560286000620012bd604088016020890162005456565b6001600160a01b03166001600160a01b0316815260200190815260200160002060050154620012ed919062005cc9565b620012f9919062005cdf565b62001305919062005cdf565b6011546001600160a01b03166000908152602760205260408120805492935083929091906200133690849062005cdf565b909155508290506028600062001353604084016020850162005456565b6001600160a01b03168152602081019190915260400160002062001378828262005ee9565b50505050565b6018546001600160a01b03163314620013ab5760405162461bcd60e51b815260040162000dc69062005c09565b6001600160a01b03821660009081526027602052604081208054839290620013d590849062005cdf565b925050819055508060296000828254620013f0919062005cdf565b90915550506040518181526001600160a01b038316907f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf906020015b60405180910390a25050565b6018546001600160a01b03163314620014655760405162461bcd60e51b815260040162000dc69062005c09565b601380546001600160a01b0319166001600160a01b03838116918217909255601454604051637adbf97360e01b8152600481019290925290911690637adbf97390602401600060405180830381600087803b158015620014c457600080fd5b505af1158015620014d9573d6000803e3d6000fd5b5050601654604051637adbf97360e01b81526001600160a01b0385811660048301529091169250637adbf97391506024015b600060405180830381600087803b1580156200152657600080fd5b505af11580156200153b573d6000803e3d6000fd5b5050505050565b602a546000906060906001600160a01b03163314620015755760405162461bcd60e51b815260040162000dc6906200603f565b601e54436000818152601f60209081526040822093909355601954925490926200159f9162005cdf565b6012546040516306c9789b60e41b8152929091146004830181905292506001600160a01b031690636c9789b090602401600060405180830381600087803b158015620015ea57600080fd5b505af1158015620015ff573d6000803e3d6000fd5b5050505080156200175f5762001614620033ff565b6200161e620037e9565b62001628620038d8565b6000806200163562002742565b60135460405163422811f960e11b81529294509092506001600160a01b03169063845023f2906200166b90859060040162005b33565b600060405180830381600087803b1580156200168657600080fd5b505af11580156200169b573d6000803e3d6000fd5b505060125460405163e08b14ed60e01b81526001600160a01b03909116925063e08b14ed9150620016d190849060040162005b33565b600060405180830381600087803b158015620016ec57600080fd5b505af115801562001701573d6000803e3d6000fd5b50505050436020819055506001601e600082825462001721919062005cdf565b9091555050601e546040519081527febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e3359060200160405180910390a150505b60135460408051634bb278f360e01b815290516000926001600160a01b031691634bb278f3916004808301926020929190829003018187875af1158015620017ab573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620017d1919062006082565b905080156200184d576014546040805163a2e6204560e01b815290516001600160a01b039092169163a2e620459160048082019260209290919082900301816000875af192505050801562001845575060408051601f3d908101601f19168201909252620018429181019062006082565b60015b156200184d57505b600254602280546040805160208084028201810190925282815260ff9094169391839160009084015b8282101562001964576000848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620018cc9062005c77565b80601f0160208091040260200160405190810160405280929190818152602001828054620018fa9062005c77565b80156200194b5780601f106200191f576101008083540402835291602001916200194b565b820191906000526020600020905b8154815290600101906020018083116200192d57829003601f168201915b5050505050815250508152602001906001019062001876565b5050505090509350935050509091565b6018546001600160a01b03163314620019a15760405162461bcd60e51b815260040162000dc69062005c09565b601880546001600160a01b0319166001600160a01b0383811691821790925560135460405163b3ab15fb60e01b815260048101929092529091169063b3ab15fb90602401600060405180830381600087803b15801562001a0057600080fd5b505af115801562001a15573d6000803e3d6000fd5b505060145460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001a6157600080fd5b505af115801562001a76573d6000803e3d6000fd5b505060155460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001ac257600080fd5b505af115801562001ad7573d6000803e3d6000fd5b505060165460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001b2357600080fd5b505af115801562001b38573d6000803e3d6000fd5b505060175460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024016200150b565b600080805b60225481101562001bcb576022818154811062001b955762001b95620060a6565b9060005260206000209060030201600101548262001bb4919062005cdf565b91508062001bc281620060bc565b91505062001b74565b508060000362001c1e5760405162461bcd60e51b815260206004820152601c60248201527f54686520636f6d6d6974746565206973206e6f74207374616b696e6700000000604482015260640162000dc6565b60008362001c2e606387620060d8565b62001c3a919062005cdf565b905060008160405160200162001c5291815260200190565b60408051601f1981840301815291905280516020909101209050600062001c7a848362006108565b90506000805b60225481101562001d2a576022818154811062001ca15762001ca1620060a6565b9060005260206000209060030201600101548262001cc0919062005cdf565b915062001ccf60018362005cc9565b831162001d15576022818154811062001cec5762001cec620060a6565b60009182526020909120600390910201546001600160a01b0316965062000d8095505050505050565b8062001d2181620060bc565b91505062001c80565b5060405162461bcd60e51b815260206004820152602960248201527f5468657265206973206e6f2076616c696461746f72206c65667420696e20746860448201526865206e6574776f726b60b81b606482015260840162000dc6565b6018546001600160a01b0316331462001db35760405162461bcd60e51b815260040162000dc69062005c09565b60195481101562001e6a578060205462001dce919062005cdf565b431062001e6a5760405162461bcd60e51b815260206004820152605760248201527f63757272656e7420636861696e2068656164206578636565642074686520776960448201527f6e646f773a206c617374426c6f636b45706f6368202b205f6e6577506572696f60648201527f642c2074727920616761696e206c6174746572206f6e2e000000000000000000608482015260a40162000dc6565b6019819055601254604051631ad7d11360e21b8152600481018390526001600160a01b0390911690636b5f444c90602401600060405180830381600087803b15801562001eb657600080fd5b505af115801562001ecb573d6000803e3d6000fd5b505050507fd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f818160405162001f0191815260200190565b60405180910390a150565b6018546001600160a01b0316331462001f395760405162461bcd60e51b815260040162000dc69062005c09565b600d55565b6000604051806102800160405280336001600160a01b0316815260200160006001600160a01b03168152602001856001600160a01b03168152602001868152602001600d6000016002015481526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b0316815260200160008152602001438152602001600081526020016000815260200160008152602001848152602001600060038111156200200d576200200d62005490565b905290506200201d8183620039f7565b60208101516101808201516040517f8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c926200205e92339289918b916200611f565b60405180910390a15050505050565b6001600160a01b0380831660008181526028602052604090206001015490911614620020ad5760405162461bcd60e51b815260040162000dc69062005b86565b6001600160a01b03828116600090815260286020526040902054163314620020e95760405162461bcd60e51b815260040162000dc69062005bbd565b6127108111156200213d5760405162461bcd60e51b815260206004820152601f60248201527f7265717569726520636f727265637420636f6d6d697373696f6e207261746500604482015260640162000dc6565b604080516060810182526001600160a01b038481168252436020808401918252838501868152600c80546000908152600a909352958220855181546001600160a01b031916951694909417845591516001808501919091559151600290930192909255835492939092909190620021b690849062005cdf565b90915550506040518281526001600160a01b038416907f4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf9060200160405180910390a2505050565b6018546001600160a01b031633146200222b5760405162461bcd60e51b815260040162000dc69062005c09565b6002805460ff19166001179055565b6018546001600160a01b03163314620022675760405162461bcd60e51b815260040162000dc69062005c09565b60008111620022b95760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000dc6565b601b55565b6018546001600160a01b03163314620022eb5760405162461bcd60e51b815260040162000dc69062005c09565b6001600160a01b0382166000908152602760205260409020548111156200234e5760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000dc6565b6001600160a01b038216600090815260276020526040812080548392906200237890849062005cc9565b92505081905550806029600082825462002393919062005cc9565b90915550506040518181526001600160a01b038316907f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3906020016200142c565b6001600160a01b0380831660008181526028602052604090206001015490911614620024145760405162461bcd60e51b815260040162000dc69062005c40565b6001600160a01b03821660009081526028602052604081206013015460ff16600381111562002447576200244762005490565b14620024965760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000dc6565b620024a382823362003e51565b5050565b6001600160a01b0380831660008181526028602052604090206001015490911614620024e75760405162461bcd60e51b815260040162000dc69062005c40565b60008111620025315760405162461bcd60e51b81526020600482015260156024820152740756e626f6e64696e6720616d6f756e74206973203605c1b604482015260640162000dc6565b620024a38282336200403d565b60606025805480602002602001604051908101604052809291908181526020016000905b8282101562002618578382906000526020600020018054620025849062005c77565b80601f0160208091040260200160405190810160405280929190818152602001828054620025b29062005c77565b8015620026035780601f10620025d75761010080835404028352916020019162002603565b820191906000526020600020905b815481529060010190602001808311620025e557829003601f168201915b50505050508152602001906001019062002562565b50505050905090565b600062000d7c338484620032f6565b60606022805480602002602001604051908101604052809291908181526020016000905b8282101562002618576000848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620026aa9062005c77565b80601f0160208091040260200160405190810160405280929190818152602001828054620026d89062005c77565b8015620027295780601f10620026fd5761010080835404028352916020019162002729565b820191906000526020600020905b8154815290600101906020018083116200270b57829003601f168201915b5050505050815250508152602001906001019062002654565b602a5460609081906001600160a01b03163314620027745760405162461bcd60e51b815260040162000dc6906200603f565b601d54620027c55760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000dc6565b620027cf6200529e565b601b546080820152601d8152602860208201526022604082015260216060820152620027fb8162004370565b6200280960256000620052bc565b60225480620028505760405162461bcd60e51b8152602060048201526012602482015271636f6d6d697474656520697320656d70747960701b604482015260640162000dc6565b60008167ffffffffffffffff8111156200286e576200286e62005858565b60405190808252806020026020018201604052801562002898578160200160208202803683370190505b50905060008267ffffffffffffffff811115620028b957620028b962005858565b604051908082528060200260200182016040528015620028e3578160200160208202803683370190505b50905060005b8381101562002a0257600060286000602284815481106200290e576200290e620060a6565b60009182526020808320600392830201546001600160a01b031684528301939093526040909101812060258054600181018255925292507f401968ff42a154441da5f6c4c935ac46b8671f0e062baaa62a7545ba53bb6e4c0190620029769083018262006168565b50600281015484516001600160a01b03909116908590849081106200299f576200299f620060a6565b6001600160a01b03928316602091820292909201015260018201548451911690849084908110620029d457620029d4620060a6565b6001600160a01b03909216602092830291909101909101525080620029f981620060bc565b915050620028e9565b50909450925050509091565b6018546001600160a01b0316331462002a3b5760405162461bcd60e51b815260040162000dc69062005c09565b62002a486000836200438b565b620024a36001826200438b565b6018546001600160a01b0316331462002a825760405162461bcd60e51b815260040162000dc69062005c09565b601580546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038082166000818152602860205260409020600101549091161462002ae45760405162461bcd60e51b815260040162000dc69062005b86565b6001600160a01b0380821660009081526028602052604090208054909116331462002b6a5760405162461bcd60e51b815260206004820152602f60248201527f726571756972652063616c6c657220746f2062652076616c696461746f72207460448201526e1c99585cdd5c9e481858d8dbdd5b9d608a1b606482015260840162000dc6565b6000601382015460ff16600381111562002b885762002b8862005490565b0362002bd75760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f7220616c7265616479206163746976650000000000000000604482015260640162000dc6565b6002601382015460ff16600381111562002bf55762002bf562005490565b14801562002c065750438160100154115b1562002c555760405162461bcd60e51b815260206004820152601760248201527f76616c696461746f72207374696c6c20696e206a61696c000000000000000000604482015260640162000dc6565b6003601382015460ff16600381111562002c735762002c7362005490565b0362002cc25760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206a61696c6564207065726d616e656e746c7900000000604482015260640162000dc6565b60138101805460ff1916905580546019546020546001600160a01b038581169316917f60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b59162002d12919062005cdf565b60405190815260200160405180910390a35050565b6060806000600181805462002d3c9062005c77565b80601f016020809104026020016040519081016040528092919081815260200182805462002d6a9062005c77565b801562002dbb5780601f1062002d8f5761010080835404028352916020019162002dbb565b820191906000526020600020905b81548152906001019060200180831162002d9d57829003601f168201915b5050505050915080805462002dd09062005c77565b80601f016020809104026020016040519081016040528092919081815260200182805462002dfe9062005c77565b801562002e4f5780601f1062002e235761010080835404028352916020019162002e4f565b820191906000526020600020905b81548152906001019060200180831162002e3157829003601f168201915b50505050509050915091509091565b6060601d80548060200260200160405190810160405280929190818152602001828054801562002eb857602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831162002e99575b5050505050905090565b6018546001600160a01b0316331462002eef5760405162461bcd60e51b815260040162000dc69062005c09565b600e8190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200162001f01565b6018546001600160a01b0316331462002f525760405162461bcd60e51b815260040162000dc69062005c09565b601780546001600160a01b0319166001600160a01b0392909216919091179055565b6018546001600160a01b0316331462002fa15760405162461bcd60e51b815260040162000dc69062005c09565b62002fae600080620052dc565b62002fbc60016000620052dc565b6002805460ff19169055565b6018546001600160a01b0316331462002ff55760405162461bcd60e51b815260040162000dc69062005c09565b601680546001600160a01b0319166001600160a01b0392909216919091179055565b6018546001600160a01b03163314620030445760405162461bcd60e51b815260040162000dc69062005c09565b601480546001600160a01b0319166001600160a01b0392909216919091179055565b602a546001600160a01b03163314620030935760405162461bcd60e51b815260040162000dc6906200603f565b6200309d620037e9565b620024a362002742565b6018546001600160a01b03163314620030d45760405162461bcd60e51b815260040162000dc69062005c09565b601180546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b0383166200315a5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000dc6565b6001600160a01b038216620031bd5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000dc6565b6001600160a01b0383811660008181526026602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b038116600090815260286020526040812090601382015460ff16600381111562003254576200325462005490565b14620032a35760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265206163746976650000000000000000604482015260640162000dc6565b60138101805460ff1916600117905580546019546020546001600160a01b038581169316917f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c9162002d12919062005cdf565b6001600160a01b038316600090815260276020526040902054811115620033595760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000dc6565b6001600160a01b038316600090815260276020526040812080548392906200338390849062005cc9565b90915550506001600160a01b03821660009081526027602052604081208054839290620033b290849062005cdf565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200321291815260200190565b476000036200340a57565b600d544790600090670de0b6b3a76400009062003429908490620060d8565b62003435919062006251565b90508015620034b3576011546040516000916001600160a01b03169083908381818185875af1925050503d80600081146200348d576040519150601f19603f3d011682016040523d82523d6000602084013e62003492565b606091505b5090915050801515600103620034b157620034ae828462005cc9565b92505b505b8160236000828254620034c7919062005cdf565b90915550600090505b602254811015620037e45760006028600060228481548110620034f757620034f7620060a6565b600091825260208083206003909202909101546001600160a01b03168352820192909252604001812060215460228054929450909187919086908110620035425762003542620060a6565b906000526020600020906003020160010154620035609190620060d8565b6200356c919062006251565b90508015620037cc576002601383015460ff16600381111562003593576200359362005490565b1480620035bb57506003601383015460ff166003811115620035b957620035b962005490565b145b156200366357601254602280546001600160a01b0390921691631de9d9b691849187908110620035ef57620035ef620060a6565b600091825260209091206003909102015460405160e084901b6001600160e01b03191681526001600160a01b0390911660048201526024016000604051808303818588803b1580156200364157600080fd5b505af115801562003656573d6000803e3d6000fd5b50505050505050620037cf565b600082600501548284600801546200367c9190620060d8565b62003688919062006251565b9050600062003698828462005cc9565b90508115620036fd5783546040516001600160a01b03909116906108fc9084906000818181858888f193505050503d8060008114620036f4576040519150601f19603f3d011682016040523d82523d6000602084013e620036f9565b606091505b5050505b8015620037855783600c0160009054906101000a90046001600160a01b03166001600160a01b031663fb489a7b826040518263ffffffff1660e01b815260040160206040518083038185885af11580156200375c573d6000803e3d6000fd5b50505050506040513d601f19601f8201168201806040525081019062003783919062006268565b505b60018401546040518481526001600160a01b03909116907fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe5639060200160405180910390a250505b50505b80620037db81620060bc565b915050620034d0565b505050565b6004545b600554811015620038165762003810816200380881620060bc565b9250620044de565b620037ed565b50600554600455600854600754036200382b57565b6009545b600854811015620038585762003852816200384a81620060bc565b925062004705565b6200382f565b50600854600955600754805b600854811015620038d2576010546000828152600660205260409020600401544391620038919162005cdf565b11620038b757620038a28162004a00565b620038af60018362005cdf565b9150620038bd565b620038d2565b80620038c981620060bc565b91505062003864565b50600755565b600c54600b541015620039f557600b546000908152600a60205260409020601054600182015443916200390b9162005cdf565b1115620039155750565b600281015481546001600160a01b03908116600090815260286020526040808220600490810185905585548416835291819020600c015490516319fac8fd60e01b81529216926319fac8fd9262003970920190815260200190565b600060405180830381600087803b1580156200398b57600080fd5b505af1158015620039a0573d6000803e3d6000fd5b5050600b80546000908152600a6020526040812080546001600160a01b031916815560018082018390556002909101829055825490945091925090620039e890849062005cdf565b90915550620038d8915050565b565b60e281511462003a415760405162461bcd60e51b8152602060048201526014602482015273092dcecc2d8d2c840e0e4dedecc40d8cadccee8d60631b604482015260640162000dc6565b6030826102400151511462003a995760405162461bcd60e51b815260206004820152601c60248201527f496e76616c696420636f6e73656e737573206b6579206c656e67746800000000604482015260640162000dc6565b62003aa48262004b4b565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a00000000000060208083019190915284519251919260009262003b07920160609190911b6bffffffffffffffffffffffff1916815260140190565b604051602081830303815290604052905060008262003b27835162004c78565b8360405160200162003b3c9392919062006282565b60408051601f198184030181528282528051602091820120600280855260608501845290945060009392909183019080368337019050509050600080808062003b94898262003b8e60416002620060d8565b62004d99565b9050600062003bb38a62003bab60416002620060d8565b606062004d99565b905060205b825181101562003c845762003bce838262004eb2565b6040805160008152602081018083528d905260ff8316918101919091526060810184905260808101839052929850909650945060019060a0016020604051602081039080840390855afa15801562003c2a573d6000803e3d6000fd5b5050604051601f19015190508762003c4460418462006251565b8151811062003c575762003c57620060a6565b6001600160a01b039092166020928302919091019091015262003c7c60418262005cdf565b905062003bb8565b508a602001516001600160a01b03168660008151811062003ca95762003ca9620060a6565b60200260200101516001600160a01b03161462003d1b5760405162461bcd60e51b815260206004820152602960248201527f496e76616c6964206e6f6465206b6579206f776e6572736869702070726f6f66604482015268081c1c9bdd9a59195960ba1b606482015260840162000dc6565b8a604001516001600160a01b03168660018151811062003d3f5762003d3f620060a6565b60200260200101516001600160a01b03161462003db35760405162461bcd60e51b815260206004820152602b60248201527f496e76616c6964206f7261636c65206b6579206f776e6572736869702070726f60448201526a1bd9881c1c9bdd9a59195960aa1b606482015260840162000dc6565b600162003dcb8c6102400151838e6000015162004ee9565b1462003e395760405162461bcd60e51b815260206004820152603660248201527f496e76616c696420636f6e73656e737573206b6579206f776e65727368697020604482015275383937b7b3103337b9103932b3b4b9ba3930ba34b7b760511b606482015260840162000dc6565b62003e448b62004f58565b5050505050505050505050565b6000821162003eaf5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000dc6565b6001600160a01b03811660009081526027602052604090205482111562003f195760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000dc6565b6001600160a01b0381166000908152602760205260408120805484929062003f4390849062005cc9565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580546000908152600394859052978820875181549088166001600160a01b031991821617825595516001820180549190981696169590951790955590516002840155519101558054919262003fcf83620060bc565b90915550506001600160a01b03848116600081815260286020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d91015b60405180910390a35050505050565b6001600160a01b0380841660009081526028602052604090208054909183811691161480620041b357600c820154604051631092ab9160e31b81526001600160a01b03858116600483015260009216906384955c8890602401602060405180830381865afa158015620040b4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620040da919062006268565b905084811015620041425760405162461bcd60e51b815260206004820152602b60248201527f696e73756666696369656e7420756e6c6f636b6564204c6971756964204e657760448201526a746f6e2062616c616e636560a81b606482015260840162000dc6565b600c83015460405163282d3fdf60e01b81526001600160a01b038681166004830152602482018890529091169063282d3fdf90604401600060405180830381600087803b1580156200419357600080fd5b505af1158015620041a8573d6000803e3d6000fd5b505050505062004246565b8382600b01548360080154620041ca919062005cc9565b10156200422a5760405162461bcd60e51b815260206004820152602760248201527f696e73756666696369656e742073656c6620626f6e646564206e6577746f6e2060448201526662616c616e636560c81b606482015260840162000dc6565b8382600b01600082825462004240919062005cdf565b90915550505b6040805160e0810182526001600160a01b0380861682528781166020808401918252838501898152600060608601818152436080880190815260a088018381528a151560c08a019081526008805486526006909752998420985189549089166001600160a01b0319918216178a55965160018a01805491909916971696909617909655915160028701559051600386015592516004850155905160059093018054945115156101000261ff00199415159490941661ffff1990951694909417929092179092558054916200431a83620060bc565b9190505550826001600160a01b0316856001600160a01b03167f63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc83876040516200402e9291909115158252602082015260400190565b60fa60a06000808285855af4620037e4573d6000803e3d6000fd5b8154600260018083161561010002038216048251808201602081106020841001600281146200443a576001811462004460578660005260208404602060002001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b8082101562004423578151845560018401935060208201915062004404565b815191036101000a908190040290915550620044d5565b60028302826020036101000a846020036101000a602089015104020185018755620044d5565b8660005260208404602060002001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b80821015620044c35781518455600184019350602082019150620044a4565b815191036101000a9081900402909155505b50505050505050565b600081815260036020908152604080832060018101546001600160a01b03168452602890925282209091601382015460ff16600381111562004524576200452462005490565b14620045c457600282015482546001600160a01b0316600090815260276020526040812080549091906200455a90849062005cdf565b909155505081546001830154600284015460138401546040517f1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f87834294620045b7946001600160a01b0391821694911692909160ff90911690620062cb565b60405180910390a1505050565b805482546001600160a01b03908116911614620046c15760008082600801548360050154620045f4919062005cc9565b9050806000036200460c578360020154915062004632565b80846002015484600d0154620046239190620060d8565b6200462f919062006251565b91505b600c83015484546040516340c10f1960e01b81526001600160a01b039182166004820152602481018590529116906340c10f1990604401600060405180830381600087803b1580156200468457600080fd5b505af115801562004699573d6000803e3d6000fd5b505050508183600d016000828254620046b3919062005cdf565b90915550620046e192505050565b8160020154816008016000828254620046db919062005cdf565b90915550505b8160020154816005016000828254620046fb919062005cdf565b9091555050505050565b600081815260066020908152604080832060018101546001600160a01b0316845260289092528220600582015491929091610100900460ff16620048f9576002830154600c8301548454604051637eee288d60e01b81526001600160a01b03918216600482015260248101849052911690637eee288d90604401600060405180830381600087803b1580156200479a57600080fd5b505af1158015620047af573d6000803e3d6000fd5b50505050600c8301548454604051632770a7eb60e21b81526001600160a01b03918216600482015260248101849052911690639dc29fac90604401600060405180830381600087803b1580156200480557600080fd5b505af11580156200481a573d6000803e3d6000fd5b5050505060008360080154846005015462004836919062005cc9565b600d8501549091506200484a8284620060d8565b62004856919062006251565b92508184600d0160008282546200486e919062005cc9565b909155505060068401546000036200488d5760038501839055620048b5565b60068401546007850154620048a39085620060d8565b620048af919062006251565b60038601555b82846006016000828254620048cb919062005cdf565b90915550506003850154600785018054600090620048eb90849062005cdf565b90915550620049ce92505050565b506002820154600882015481111562004913575060088101545b81600901546000036200492d576003830181905562004955565b6009820154600a830154620049439083620060d8565b6200494f919062006251565b60038401555b808260090160008282546200496b919062005cdf565b90915550506003830154600a830180546000906200498b90849062005cdf565b9250508190555080826008016000828254620049a8919062005cc9565b90915550506002830154600b83018054600090620049c890849062005cc9565b90915550505b6005808401805460ff1916600117905582018054829190600090620049f590849062005cc9565b909155505050505050565b6000818152600660205260408120600381015490910362004a1f575050565b60018101546001600160a01b031660009081526028602052604081206005830154909190610100900460ff1662004abb5781600701548260060154846003015462004a6b9190620060d8565b62004a77919062006251565b90508082600601600082825462004a8f919062005cc9565b9091555050600383015460078301805460009062004aaf90849062005cc9565b9091555062004b209050565b81600a01548260090154846003015462004ad69190620060d8565b62004ae2919062006251565b90508082600901600082825462004afa919062005cc9565b90915550506003830154600a8301805460009062004b1a90849062005cc9565b90915550505b82546001600160a01b031660009081526027602052604081208054839290620049f590849062005cdf565b600062004b5c826060015162005189565b6001600160a01b0390911660208401529050801562004bac5760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000dc6565b6020808301516001600160a01b0390811660009081526028909252604090912060010154161562004c205760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000dc6565b61271082608001511115620024a35760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000dc6565b60608160000362004ca05750506040805180820190915260018152600360fc1b602082015290565b8160005b811562004cd0578062004cb781620060bc565b915062004cc89050600a8362006251565b915062004ca4565b60008167ffffffffffffffff81111562004cee5762004cee62005858565b6040519080825280601f01601f19166020018201604052801562004d19576020820181803683370190505b5090505b841562004d915762004d3160018362005cc9565b915062004d40600a8662006108565b62004d4d90603062005cdf565b60f81b81838151811062004d655762004d65620060a6565b60200101906001600160f81b031916908160001a90535062004d89600a8662006251565b945062004d1d565b949350505050565b60608162004da981601f62005cdf565b101562004dea5760405162461bcd60e51b815260206004820152600e60248201526d736c6963655f6f766572666c6f7760901b604482015260640162000dc6565b62004df6828462005cdf565b8451101562004e3c5760405162461bcd60e51b8152602060048201526011602482015270736c6963655f6f75744f66426f756e647360781b604482015260640162000dc6565b60608215801562004e5d576040519150600082526020820160405262004ea9565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101562004e9857805183526020928301920162004e7e565b5050858452601f01601f1916604052505b50949350505050565b8181018051602082015160409092015190919060001a601b81101562004ee25762004edf601b82620062f9565b90505b9250925092565b600062004ef56200531b565b600085858560405160200162004f0e9392919062006315565b6040516020818303038152906040529050600060fb905060008251602062004f37919062005cdf565b90506020848285855afa62004f4b57600080fd5b5050905195945050505050565b6101808101516001600160a01b031662004fdb57601d5460009062004f7d9062004c78565b90508160200151826000015183608001518360405162004f9d9062005339565b62004fac949392919062006364565b604051809103906000f08015801562004fc9573d6000803e3d6000fd5b506001600160a01b0316610180830152505b60208181018051601d80546001808201835560009283527f6d4407e7be21f808e6509aa9fa9143369579dd7d760fe20a2c09680fc146134f90910180546001600160a01b03199081166001600160a01b03958616179091558451841683526028909552604091829020865181548716908516178155935190840180548616918416919091179055840151600283018054909416911617909155606082015182919060038201906200508d9082620063a3565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620051579082620063a3565b5061026082015160138201805460ff191660018360038111156200517f576200517f62005490565b0217905550505050565b6000806200519662005347565b60008060ff9050604083875160208901845afa620051b357600080fd5b50508051602090910151600160601b90910494909350915050565b60405180610280016040528060006001600160a01b0316815260200160006001600160a01b0316815260200160006001600160a01b0316815260200160608152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b031681526020016000815260200160008152602001600081526020016000815260200160008152602001606081526020016000600381111562005299576200529962005490565b905290565b6040518060a001604052806005906020820280368337509192915050565b508054600082559060005260206000209081019062000e16919062005365565b508054620052ea9062005c77565b6000825580601f10620052fb575050565b601f01602090049060005260206000209081019062000e1691906200538a565b60405180602001604052806001906020820280368337509192915050565b611531806200646483390190565b60405180604001604052806002906020820280368337509192915050565b80821115620053865760006200537c8282620052dc565b5060010162005365565b5090565b5b808211156200538657600081556001016200538b565b60005b83811015620053be578181015183820152602001620053a4565b50506000910152565b60008151808452620053e1816020860160208601620053a1565b601f01601f19169290920160200192915050565b6020815260006200540a6020830184620053c7565b9392505050565b6001600160a01b038116811462000e1657600080fd5b600080604083850312156200543b57600080fd5b8235620054488162005411565b946020939093013593505050565b6000602082840312156200546957600080fd5b81356200540a8162005411565b6000602082840312156200548957600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b60048110620054c557634e487b7160e01b600052602160045260246000fd5b9052565b60208152620054e46020820183516001600160a01b03169052565b600060208301516200550160408401826001600160a01b03169052565b5060408301516001600160a01b0381166060840152506060830151610280806080850152620055356102a0850183620053c7565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a0620055bd818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f190161026080880191909152909150620056218483620053c7565b9350808701519150506200563882860182620054a6565b5090949350505050565b6000806000606084860312156200565857600080fd5b8335620056658162005411565b92506020840135620056778162005411565b929592945050506040919091013590565b6000602082840312156200569b57600080fd5b813567ffffffffffffffff811115620056b357600080fd5b820161028081850312156200540a57600080fd5b600081518084526020808501808196508360051b8101915082860160005b858110156200573c578284038952815180516001600160a01b0316855285810151868601526040908101516060918601829052906200572781870183620053c7565b9a87019a9550505090840190600101620056e5565b5091979650505050505050565b821515815260406020820152600062004d916040830184620056c7565b600080604083850312156200577a57600080fd5b50508035926020909101359150565b845181526020808601518183015260408087015181840152606080880151818501526080808901516001600160a01b03908116828701528851811660a08088019190915294890151811660c087015292880151831660e086015290870151821661010085015286015116610120830152840151610200820190620058196101408401826001600160a01b03169052565b5083516001600160a01b0316610160830152602084015161018083015260408401516101a08301526060909301516101c08201526101e0015292915050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200588057600080fd5b813567ffffffffffffffff808211156200589e576200589e62005858565b604051601f8301601f19908116603f01168101908282118183101715620058c957620058c962005858565b81604052838152866020858801011115620058e357600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080608085870312156200591a57600080fd5b843567ffffffffffffffff808211156200593357600080fd5b62005941888389016200586e565b955060208701359150620059558262005411565b909350604086013590808211156200596c57600080fd5b6200597a888389016200586e565b935060608701359150808211156200599157600080fd5b50620059a0878288016200586e565b91505092959194509250565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101562005a0557603f19888603018452620059f2858351620053c7565b94509285019290850190600101620059d3565b5092979650505050505050565b6020815260006200540a6020830184620056c7565b600081518084526020808501945080840160005b8381101562005a625781516001600160a01b03168752958201959082019060010162005a3b565b509495945050505050565b60408152600062005a82604083018562005a27565b828103602084015262005a96818562005a27565b95945050505050565b6000806040838503121562005ab357600080fd5b823567ffffffffffffffff8082111562005acc57600080fd5b62005ada868387016200586e565b9350602085013591508082111562005af157600080fd5b5062005b00858286016200586e565b9150509250929050565b60408152600062005b1f6040830185620053c7565b828103602084015262005a968185620053c7565b6020815260006200540a602083018462005a27565b6000806040838503121562005b5c57600080fd5b823562005b698162005411565b9150602083013562005b7b8162005411565b809150509250929050565b6020808252601c908201527f76616c696461746f72206d757374206265207265676973746572656400000000604082015260600190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b60208082526018908201527f76616c696461746f72206e6f7420726567697374657265640000000000000000604082015260600190565b600181811c9082168062005c8c57607f821691505b60208210810362005cad57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000d805762000d8062005cb3565b8082018082111562000d805762000d8062005cb3565b6000813562000d808162005411565b80546001600160a01b0319166001600160a01b0392909216919091179055565b6000808335601e1984360301811262005d3c57600080fd5b83018035915067ffffffffffffffff82111562005d5857600080fd5b60200191503681900382131562005d6e57600080fd5b9250929050565b601f821115620037e457600081815260208120601f850160051c8101602086101562005d9e5750805b601f850160051c820191505b8181101562005dbf5782815560010162005daa565b505050505050565b600019600383901b1c191660019190911b1790565b67ffffffffffffffff83111562005df75762005df762005858565b62005e0f8362005e08835462005c77565b8362005d75565b6000601f84116001811462005e42576000851562005e2d5750838201355b62005e39868262005dc7565b8455506200153b565b600083815260209020601f19861690835b8281101562005e75578685013582556020948501946001909201910162005e53565b508682101562005e935760001960f88860031b161c19848701351681555b505060018560011b0183555050505050565b600081356004811062000d8057600080fd5b6004821062005ed657634e487b7160e01b600052602160045260246000fd5b60ff1981541660ff831681178255505050565b62005eff62005ef88362005cf5565b8262005d04565b62005f1b62005f116020840162005cf5565b6001830162005d04565b62005f3762005f2d6040840162005cf5565b6002830162005d04565b62005f46606083018362005d24565b62005f5681836003860162005ddc565b50506080820135600482015560a0820135600582015560c0820135600682015560e0820135600782015561010082013560088201556101208201356009820155610140820135600a820155610160820135600b82015562005fc962005fbf610180840162005cf5565b600c830162005d04565b6101a0820135600d8201556101c0820135600e8201556101e0820135600f820155610200820135601082015561022082013560118201556200601061024083018362005d24565b6200602081836012860162005ddc565b5050620024a362006035610260840162005ea5565b6013830162005eb7565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b6000602082840312156200609557600080fd5b815180151581146200540a57600080fd5b634e487b7160e01b600052603260045260246000fd5b600060018201620060d157620060d162005cb3565b5060010190565b808202811582820484141762000d805762000d8062005cb3565b634e487b7160e01b600052601260045260246000fd5b6000826200611a576200611a620060f2565b500690565b600060018060a01b0380881683528087166020840152808616604084015260a060608401526200615360a0840186620053c7565b91508084166080840152509695505050505050565b81810362006174575050565b62006180825462005c77565b67ffffffffffffffff8111156200619b576200619b62005858565b620061b381620061ac845462005c77565b8462005d75565b6000601f821160018114620061e65760008315620061d15750848201545b620061dd848262005dc7565b8555506200153b565b600085815260209020601f19841690600086815260209020845b8381101562006222578286015482556001958601959091019060200162006200565b5085831015620062415781850154600019600388901b60f8161c191681555b5050505050600190811b01905550565b600082620062635762006263620060f2565b500490565b6000602082840312156200627b57600080fd5b5051919050565b6000845162006296818460208901620053a1565b845190830190620062ac818360208901620053a1565b8451910190620062c1818360208801620053a1565b0195945050505050565b6001600160a01b03858116825284166020820152604081018390526080810162005a966060830184620054a6565b60ff818116838216019081111562000d805762000d8062005cb3565b6000845162006329818460208901620053a1565b8451908301906200633f818360208901620053a1565b60609490941b6bffffffffffffffffffffffff19169301928352505060140192915050565b6001600160a01b03858116825284166020820152604081018390526080606082018190526000906200639990830184620053c7565b9695505050505050565b815167ffffffffffffffff811115620063c057620063c062005858565b620063d181620061ac845462005c77565b602080601f831160018114620064055760008415620063f05750858301515b620063fc858262005dc7565b86555062005dbf565b600085815260208120601f198616915b82811015620064365788860151825594840194600190910190840162006415565b50858210156200624157939096015160001960f8600387901b161c19169092555050600190811b0190555056fe60806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220e1e1658d868fb41912eb4d5ecdf5bed6ec95bd7db81bea3e11358a02e3aa535c64736f6c63430008150033a2646970667358221220459de632b4e4dbdff5e8aa8b8794e2ee5442764ce8fc73d97b9483908f48dce264736f6c6343000815003360806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220e1e1658d868fb41912eb4d5ecdf5bed6ec95bd7db81bea3e11358a02e3aa535c64736f6c63430008150033",
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
	return address, gasConsumed, &Autonity{contract: c}, nil
}

// Autonity is an auto generated Go binding around an Ethereum contract.
type Autonity struct {
	*contract
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Autonity *Autonity) COMMISSIONRATEPRECISION(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) Allowance(opts *runOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_Autonity *Autonity) BalanceOf(opts *runOptions, _addr common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "balanceOf", _addr)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,address) policy, (address,address,address,address,address,address) contracts, (address,uint256,uint256,uint256) protocol, uint256 contractVersion)
func (_Autonity *Autonity) Config(opts *runOptions) (struct {
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
func (_Autonity *Autonity) Decimals(opts *runOptions) (uint8, uint64, error) {
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
func (_Autonity *Autonity) Deployer(opts *runOptions) (common.Address, uint64, error) {
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
func (_Autonity *Autonity) EpochID(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) EpochReward(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) EpochTotalBondedStake(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetBlockPeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetCommittee(opts *runOptions) ([]AutonityCommitteeMember, uint64, error) {
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
func (_Autonity *Autonity) GetCommitteeEnodes(opts *runOptions) ([]string, uint64, error) {
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
func (_Autonity *Autonity) GetEpochFromBlock(opts *runOptions, _block *big.Int) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetEpochPeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetLastEpochBlock(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetMaxCommitteeSize(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetMinimumBaseFee(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetNewContract(opts *runOptions) ([]byte, string, uint64, error) {
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
func (_Autonity *Autonity) GetOperator(opts *runOptions) (common.Address, uint64, error) {
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
func (_Autonity *Autonity) GetOracle(opts *runOptions) (common.Address, uint64, error) {
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
func (_Autonity *Autonity) GetProposer(opts *runOptions, height *big.Int, round *big.Int) (common.Address, uint64, error) {
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
func (_Autonity *Autonity) GetTreasuryAccount(opts *runOptions) (common.Address, uint64, error) {
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
func (_Autonity *Autonity) GetTreasuryFee(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetUnbondingPeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) GetValidator(opts *runOptions, _addr common.Address) (AutonityValidator, uint64, error) {
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
func (_Autonity *Autonity) GetValidators(opts *runOptions) ([]common.Address, uint64, error) {
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
func (_Autonity *Autonity) GetVersion(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "getVersion")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_Autonity *Autonity) LastEpochBlock(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) Name(opts *runOptions) (string, uint64, error) {
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
func (_Autonity *Autonity) Symbol(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "symbol")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// TotalRedistributed is a free data retrieval call binding the contract method 0x9bb851c0.
//
// Solidity: function totalRedistributed() view returns(uint256)
func (_Autonity *Autonity) TotalRedistributed(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _Autonity.call(opts, "totalRedistributed")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Autonity *Autonity) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Autonity *Autonity) ActivateValidator(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "activateValidator", _address)
	return consumed, err
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Autonity *Autonity) Approve(opts *runOptions, spender common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "approve", spender, amount)
	return consumed, err
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_Autonity *Autonity) Bond(opts *runOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "bond", _validator, _amount)
	return consumed, err
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_Autonity *Autonity) Burn(opts *runOptions, _addr common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "burn", _addr, _amount)
	return consumed, err
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_Autonity *Autonity) ChangeCommissionRate(opts *runOptions, _validator common.Address, _rate *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "changeCommissionRate", _validator, _rate)
	return consumed, err
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_Autonity *Autonity) CompleteContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "completeContractUpgrade")
	return consumed, err
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[], address[])
func (_Autonity *Autonity) ComputeCommittee(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "computeCommittee")
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256,bytes)[])
func (_Autonity *Autonity) Finalize(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "finalize")
	return consumed, err
}

// FinalizeInitialization is a paid mutator transaction binding the contract method 0xd861b0e8.
//
// Solidity: function finalizeInitialization() returns()
func (_Autonity *Autonity) FinalizeInitialization(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "finalizeInitialization")
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_Autonity *Autonity) Mint(opts *runOptions, _addr common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "mint", _addr, _amount)
	return consumed, err
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_Autonity *Autonity) PauseValidator(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "pauseValidator", _address)
	return consumed, err
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x84467fdb.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_Autonity *Autonity) RegisterValidator(opts *runOptions, _enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "registerValidator", _enode, _oracleAddress, _consensusKey, _signatures)
	return consumed, err
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Autonity *Autonity) ResetContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "resetContractUpgrade")
	return consumed, err
}

// SetAccountabilityContract is a paid mutator transaction binding the contract method 0x1250a28d.
//
// Solidity: function setAccountabilityContract(address _address) returns()
func (_Autonity *Autonity) SetAccountabilityContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setAccountabilityContract", _address)
	return consumed, err
}

// SetAcuContract is a paid mutator transaction binding the contract method 0xd372c07e.
//
// Solidity: function setAcuContract(address _address) returns()
func (_Autonity *Autonity) SetAcuContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setAcuContract", _address)
	return consumed, err
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_Autonity *Autonity) SetCommitteeSize(opts *runOptions, _size *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setCommitteeSize", _size)
	return consumed, err
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_Autonity *Autonity) SetEpochPeriod(opts *runOptions, _period *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setEpochPeriod", _period)
	return consumed, err
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_Autonity *Autonity) SetMinimumBaseFee(opts *runOptions, _price *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setMinimumBaseFee", _price)
	return consumed, err
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_Autonity *Autonity) SetOperatorAccount(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setOperatorAccount", _account)
	return consumed, err
}

// SetOracleContract is a paid mutator transaction binding the contract method 0x496ccd9b.
//
// Solidity: function setOracleContract(address _address) returns()
func (_Autonity *Autonity) SetOracleContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setOracleContract", _address)
	return consumed, err
}

// SetStabilizationContract is a paid mutator transaction binding the contract method 0xcfd19fb9.
//
// Solidity: function setStabilizationContract(address _address) returns()
func (_Autonity *Autonity) SetStabilizationContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setStabilizationContract", _address)
	return consumed, err
}

// SetSupplyControlContract is a paid mutator transaction binding the contract method 0xb3ecbadd.
//
// Solidity: function setSupplyControlContract(address _address) returns()
func (_Autonity *Autonity) SetSupplyControlContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setSupplyControlContract", _address)
	return consumed, err
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_Autonity *Autonity) SetTreasuryAccount(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setTreasuryAccount", _account)
	return consumed, err
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_Autonity *Autonity) SetTreasuryFee(opts *runOptions, _treasuryFee *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setTreasuryFee", _treasuryFee)
	return consumed, err
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_Autonity *Autonity) SetUnbondingPeriod(opts *runOptions, _period *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setUnbondingPeriod", _period)
	return consumed, err
}

// SetUpgradeManagerContract is a paid mutator transaction binding the contract method 0xceaad455.
//
// Solidity: function setUpgradeManagerContract(address _address) returns()
func (_Autonity *Autonity) SetUpgradeManagerContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "setUpgradeManagerContract", _address)
	return consumed, err
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_Autonity *Autonity) Transfer(opts *runOptions, _recipient common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "transfer", _recipient, _amount)
	return consumed, err
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Autonity *Autonity) TransferFrom(opts *runOptions, sender common.Address, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "transferFrom", sender, recipient, amount)
	return consumed, err
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_Autonity *Autonity) Unbond(opts *runOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "unbond", _validator, _amount)
	return consumed, err
}

// UpdateValidatorAndTransferSlashedFunds is a paid mutator transaction binding the contract method 0x35be16e0.
//
// Solidity: function updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8) _val) returns()
func (_Autonity *Autonity) UpdateValidatorAndTransferSlashedFunds(opts *runOptions, _val AutonityValidator) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "updateValidatorAndTransferSlashedFunds", _val)
	return consumed, err
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Autonity *Autonity) UpgradeContract(opts *runOptions, _bytecode []byte, _abi string) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "upgradeContract", _bytecode, _abi)
	return consumed, err
}

// Fallback is a paid mutator transaction binding the contract fallback function.
// WARNING! UNTESTED
// Solidity: fallback() payable returns()
func (_Autonity *Autonity) Fallback(opts *runOptions, calldata []byte) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "", calldata)
	return consumed, err
}

// Receive is a paid mutator transaction binding the contract receive function.
// WARNING! UNTESTED
// Solidity: receive() payable returns()
func (_Autonity *Autonity) Receive(opts *runOptions) (uint64, error) {
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
 		func (_Autonity *Autonity) FilterActivatedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*AutonityActivatedValidatorIterator, error) {

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
		func (_Autonity *Autonity) WatchActivatedValidator(opts *bind.WatchOpts, sink chan<- *AutonityActivatedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseActivatedValidator(log types.Log) (*AutonityActivatedValidator, error) {
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
 		func (_Autonity *Autonity) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AutonityApprovalIterator, error) {

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
		func (_Autonity *Autonity) WatchApproval(opts *bind.WatchOpts, sink chan<- *AutonityApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseApproval(log types.Log) (*AutonityApproval, error) {
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
 		func (_Autonity *Autonity) FilterBondingRejected(opts *bind.FilterOpts) (*AutonityBondingRejectedIterator, error) {






			logs, sub, err := _Autonity.contract.FilterLogs(opts, "BondingRejected")
			if err != nil {
				return nil, err
			}
			return &AutonityBondingRejectedIterator{contract: _Autonity.contract, event: "BondingRejected", logs: logs, sub: sub}, nil
 		}

		// WatchBondingRejected is a free log subscription operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
		//
		// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount, uint8 state)
		func (_Autonity *Autonity) WatchBondingRejected(opts *bind.WatchOpts, sink chan<- *AutonityBondingRejected) (event.Subscription, error) {






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
		func (_Autonity *Autonity) ParseBondingRejected(log types.Log) (*AutonityBondingRejected, error) {
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
 		func (_Autonity *Autonity) FilterBurnedStake(opts *bind.FilterOpts, addr []common.Address) (*AutonityBurnedStakeIterator, error) {

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
		func (_Autonity *Autonity) WatchBurnedStake(opts *bind.WatchOpts, sink chan<- *AutonityBurnedStake, addr []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseBurnedStake(log types.Log) (*AutonityBurnedStake, error) {
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
 		func (_Autonity *Autonity) FilterCommissionRateChange(opts *bind.FilterOpts, validator []common.Address) (*AutonityCommissionRateChangeIterator, error) {

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
		func (_Autonity *Autonity) WatchCommissionRateChange(opts *bind.WatchOpts, sink chan<- *AutonityCommissionRateChange, validator []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseCommissionRateChange(log types.Log) (*AutonityCommissionRateChange, error) {
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
 		func (_Autonity *Autonity) FilterEpochPeriodUpdated(opts *bind.FilterOpts) (*AutonityEpochPeriodUpdatedIterator, error) {



			logs, sub, err := _Autonity.contract.FilterLogs(opts, "EpochPeriodUpdated")
			if err != nil {
				return nil, err
			}
			return &AutonityEpochPeriodUpdatedIterator{contract: _Autonity.contract, event: "EpochPeriodUpdated", logs: logs, sub: sub}, nil
 		}

		// WatchEpochPeriodUpdated is a free log subscription operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
		//
		// Solidity: event EpochPeriodUpdated(uint256 period)
		func (_Autonity *Autonity) WatchEpochPeriodUpdated(opts *bind.WatchOpts, sink chan<- *AutonityEpochPeriodUpdated) (event.Subscription, error) {



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
		func (_Autonity *Autonity) ParseEpochPeriodUpdated(log types.Log) (*AutonityEpochPeriodUpdated, error) {
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
 		func (_Autonity *Autonity) FilterMinimumBaseFeeUpdated(opts *bind.FilterOpts) (*AutonityMinimumBaseFeeUpdatedIterator, error) {



			logs, sub, err := _Autonity.contract.FilterLogs(opts, "MinimumBaseFeeUpdated")
			if err != nil {
				return nil, err
			}
			return &AutonityMinimumBaseFeeUpdatedIterator{contract: _Autonity.contract, event: "MinimumBaseFeeUpdated", logs: logs, sub: sub}, nil
 		}

		// WatchMinimumBaseFeeUpdated is a free log subscription operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
		//
		// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
		func (_Autonity *Autonity) WatchMinimumBaseFeeUpdated(opts *bind.WatchOpts, sink chan<- *AutonityMinimumBaseFeeUpdated) (event.Subscription, error) {



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
		func (_Autonity *Autonity) ParseMinimumBaseFeeUpdated(log types.Log) (*AutonityMinimumBaseFeeUpdated, error) {
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
 		func (_Autonity *Autonity) FilterMintedStake(opts *bind.FilterOpts, addr []common.Address) (*AutonityMintedStakeIterator, error) {

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
		func (_Autonity *Autonity) WatchMintedStake(opts *bind.WatchOpts, sink chan<- *AutonityMintedStake, addr []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseMintedStake(log types.Log) (*AutonityMintedStake, error) {
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
 		func (_Autonity *Autonity) FilterNewBondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*AutonityNewBondingRequestIterator, error) {

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
		func (_Autonity *Autonity) WatchNewBondingRequest(opts *bind.WatchOpts, sink chan<- *AutonityNewBondingRequest, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseNewBondingRequest(log types.Log) (*AutonityNewBondingRequest, error) {
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
 		func (_Autonity *Autonity) FilterNewEpoch(opts *bind.FilterOpts) (*AutonityNewEpochIterator, error) {



			logs, sub, err := _Autonity.contract.FilterLogs(opts, "NewEpoch")
			if err != nil {
				return nil, err
			}
			return &AutonityNewEpochIterator{contract: _Autonity.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
 		}

		// WatchNewEpoch is a free log subscription operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
		//
		// Solidity: event NewEpoch(uint256 epoch)
		func (_Autonity *Autonity) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *AutonityNewEpoch) (event.Subscription, error) {



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
		func (_Autonity *Autonity) ParseNewEpoch(log types.Log) (*AutonityNewEpoch, error) {
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
 		func (_Autonity *Autonity) FilterNewUnbondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*AutonityNewUnbondingRequestIterator, error) {

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
		func (_Autonity *Autonity) WatchNewUnbondingRequest(opts *bind.WatchOpts, sink chan<- *AutonityNewUnbondingRequest, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseNewUnbondingRequest(log types.Log) (*AutonityNewUnbondingRequest, error) {
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
 		func (_Autonity *Autonity) FilterPausedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*AutonityPausedValidatorIterator, error) {

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
		func (_Autonity *Autonity) WatchPausedValidator(opts *bind.WatchOpts, sink chan<- *AutonityPausedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParsePausedValidator(log types.Log) (*AutonityPausedValidator, error) {
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
 		func (_Autonity *Autonity) FilterRegisteredValidator(opts *bind.FilterOpts) (*AutonityRegisteredValidatorIterator, error) {







			logs, sub, err := _Autonity.contract.FilterLogs(opts, "RegisteredValidator")
			if err != nil {
				return nil, err
			}
			return &AutonityRegisteredValidatorIterator{contract: _Autonity.contract, event: "RegisteredValidator", logs: logs, sub: sub}, nil
 		}

		// WatchRegisteredValidator is a free log subscription operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
		//
		// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
		func (_Autonity *Autonity) WatchRegisteredValidator(opts *bind.WatchOpts, sink chan<- *AutonityRegisteredValidator) (event.Subscription, error) {







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
		func (_Autonity *Autonity) ParseRegisteredValidator(log types.Log) (*AutonityRegisteredValidator, error) {
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
 		func (_Autonity *Autonity) FilterRewarded(opts *bind.FilterOpts, addr []common.Address) (*AutonityRewardedIterator, error) {

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
		func (_Autonity *Autonity) WatchRewarded(opts *bind.WatchOpts, sink chan<- *AutonityRewarded, addr []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseRewarded(log types.Log) (*AutonityRewarded, error) {
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
 		func (_Autonity *Autonity) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutonityTransferIterator, error) {

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
		func (_Autonity *Autonity) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AutonityTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

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
		func (_Autonity *Autonity) ParseTransfer(log types.Log) (*AutonityTransfer, error) {
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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"ActivatedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"BondingRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewUnbondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeInitialization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getEpochFromBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_signatures\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAccountability\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAccountabilityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIACU\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAcuContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setOracleContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIStabilization\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setStabilizationContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISupplyControl\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSupplyControlContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractUpgradeManager\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setUpgradeManagerContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"_val\",\"type\":\"tuple\"}],\"name\":\"updateValidatorAndTransferSlashedFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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
		"9bb851c0": "totalRedistributed()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"a5d059ca": "unbond(address,uint256)",
		"35be16e0": "updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8))",
		"b2ea9adb": "upgradeContract(bytes,string)",
	},
	Bin: "0x60806040526000600b556000600c553480156200001b57600080fd5b50604080516000808252602082019092529062000051565b6200003d62000d4e565b815260200190600190039081620000335790505b506040805161012081018252600d546080808301918252600e5460a080850191909152600f5460c08086019190915260105460e08601526011546001600160a01b039081166101008701529385528551908101865260125484168152601354841660208281019190915260145485168288015260155485166060808401919091526016548616838601526017548616938301939093528086019190915285519283018652601854909316825260195492820192909252601a5481850152601b548183015292820192909252601c54918101829052906000036200014c57602a80546001600160a01b031916331790556200014c828262000169565b5050601c546001036200016357620001636200051f565b6200115e565b80518051600d55602080820151600e55604080830151600f55606080840151601055608093840151601180546001600160a01b03199081166001600160a01b03938416179091558487015180516012805484169185169190911790558086015160138054841691851691909117905580850151601480548416918516919091179055808401516015805484169185169190911790559586015160168054831691841691909117905560a0909501516017805487169183169190911790558286015180516018805490971692169190911790945591830151601955820151601a5590810151601b55810151601c5560005b82518110156200051a5760008382815181106200027a576200027a62000ebd565b602002602001015160a00151905060008483815181106200029f576200029f62000ebd565b60200260200101516101a00181815250506000848381518110620002c757620002c762000ebd565b602002602001015161018001906001600160a01b031690816001600160a01b031681525050600084838151811062000303576200030362000ebd565b602002602001015160a001818152505060008483815181106200032a576200032a62000ebd565b60209081029190910101516101c00152600f54845185908490811062000354576200035462000ebd565b6020026020010151608001818152505060008483815181106200037b576200037b62000ebd565b6020026020010151610260019060038111156200039c576200039c62000ed3565b90816003811115620003b257620003b262000ed3565b815250506000848381518110620003cd57620003cd62000ebd565b60200260200101516101600181815250506200040b848381518110620003f757620003f762000ebd565b60200260200101516200068d60201b60201c565b6200043884838151811062000424576200042462000ebd565b6020026020010151620007c460201b60201c565b806027600086858151811062000452576200045262000ebd565b6020026020010151600001516001600160a01b03166001600160a01b0316815260200190815260200160002060008282546200048f919062000eff565b925050819055508060296000828254620004aa919062000eff565b9250508190555062000504848381518110620004ca57620004ca62000ebd565b60200260200101516020015182868581518110620004ec57620004ec62000ebd565b602002602001015160000151620009f560201b60201c565b5080620005118162000f1b565b91505062000259565b505050565b602a546001600160a01b031633146200058b5760405162461bcd60e51b815260206004820152602360248201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60448201526218dbdb60ea1b60648201526084015b60405180910390fd5b600260286000601d600181548110620005a857620005a862000ebd565b60009182526020808320909101546001600160a01b0316835282019290925260400181206005018054909190620005e190849062000f4d565b92505081905550600260286000601d60018154811062000605576200060562000ebd565b60009182526020808320909101546001600160a01b03168352820192909252604001812060080180549091906200063e90849062000f4d565b90915550506002601c556018546001600160a01b031660009081526027602052604081206103e8905562000673908062000e1e565b620006816001600062000e1e565b6002805460ff19169055565b6000620006a4826060015162000be060201b60201c565b6001600160a01b03909116602084015290508015620006f45760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000582565b6020808301516001600160a01b03908116600090815260289092526040909120600101541615620007685760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000582565b61271082608001511115620007c05760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000582565b5050565b6101808101516001600160a01b03166200084757601d54600090620007e99062000c2e565b905081602001518260000151836080015183604051620008099062000e60565b62000818949392919062000f64565b604051809103906000f08015801562000835573d6000803e3d6000fd5b506001600160a01b0316610180830152505b60208181018051601d80546001808201835560009283527f6d4407e7be21f808e6509aa9fa9143369579dd7d760fe20a2c09680fc146134f90910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602890955260409182902086518154871690851617815593519084018054861691841691909117905584015160028301805490941691161790915560608201518291906003820190620008f9908262001065565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620009c3908262001065565b5061026082015160138201805460ff19166001836003811115620009eb57620009eb62000ed3565b0217905550505050565b6000821162000a535760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000582565b6001600160a01b03811660009081526027602052604090205482111562000abd5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000582565b6001600160a01b0381166000908152602760205260408120805484929062000ae790849062001131565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580546000908152600394859052978820875181549088166001600160a01b031991821617825595516001820180549190981696169590951790955590516002840155519101558054919262000b738362000f1b565b90915550506001600160a01b03848116600081815260286020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d910160405180910390a35050505050565b60008062000bed62000e6e565b60008060ff9050604083875160208901845afa62000c0a57600080fd5b505080516020909101516c0100000000000000000000000090910494909350915050565b60608160000362000c565750506040805180820190915260018152600360fc1b602082015290565b8160005b811562000c86578062000c6d8162000f1b565b915062000c7e9050600a8362000f4d565b915062000c5a565b6000816001600160401b0381111562000ca35762000ca362000ea7565b6040519080825280601f01601f19166020018201604052801562000cce576020820181803683370190505b5090505b841562000d465762000ce660018362001131565b915062000cf5600a8662001147565b62000d0290603062000eff565b60f81b81838151811062000d1a5762000d1a62000ebd565b60200101906001600160f81b031916908160001a90535062000d3e600a8662000f4d565b945062000cd2565b949350505050565b60405180610280016040528060006001600160a01b0316815260200160006001600160a01b0316815260200160006001600160a01b0316815260200160608152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b031681526020016000815260200160008152602001600081526020016000815260200160008152602001606081526020016000600381111562000e195762000e1962000ed3565b905290565b50805462000e2c9062000fd7565b6000825580601f1062000e3d575050565b601f01602090049060005260206000209081019062000e5d919062000e8c565b50565b6115318062008b6183390190565b60405180604001604052806002906020820280368337509192915050565b5b8082111562000ea3576000815560010162000e8d565b5090565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052601160045260246000fd5b8082018082111562000f155762000f1562000ee9565b92915050565b60006001820162000f305762000f3062000ee9565b5060010190565b634e487b7160e01b600052601260045260246000fd5b60008262000f5f5762000f5f62000f37565b500490565b600060018060a01b038087168352602081871681850152856040850152608060608501528451915081608085015260005b8281101562000fb35785810182015185820160a00152810162000f95565b5050600060a0828501015260a0601f19601f83011684010191505095945050505050565b600181811c9082168062000fec57607f821691505b6020821081036200100d57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200051a57600081815260208120601f850160051c810160208610156200103c5750805b601f850160051c820191505b818110156200105d5782815560010162001048565b505050505050565b81516001600160401b0381111562001081576200108162000ea7565b620010998162001092845462000fd7565b8462001013565b602080601f831160018114620010d15760008415620010b85750858301515b600019600386901b1c1916600185901b1785556200105d565b600085815260208120601f198616915b828110156200110257888601518255948401946001909101908401620010e1565b5085821015620011215787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b8181038181111562000f155762000f1562000ee9565b60008262001159576200115962000f37565b500690565b6179f3806200116e6000396000f3fe608060405260043610620003db5760003560e01c8063872cf0591162000203578063b66b3e791162000117578063d372c07e11620000a7578063dd62ed3e1162000075578063dd62ed3e1462000ccc578063dfb1a4d21462000d16578063e7f43c681462000d2d578063f7866ee31462000d4d57005b8063d372c07e1462000c48578063d5f394881462000c6d578063d861b0e81462000c8f578063d886f8a21462000ca757005b8063cb696f5411620000e5578063cb696f541462000bc1578063ceaad4551462000be6578063cf9c57191462000c0b578063cfd19fb91462000c2357005b8063b66b3e791462000b42578063b7ab4db51462000b6a578063c2362dd51462000b91578063c9d97af41462000ba957005b8063a5d059ca1162000193578063ae1f5fa01162000161578063ae1f5fa01462000aab578063b2ea9adb1462000ad3578063b3ecbadd1462000af8578063b46e55201462000b1d57005b8063a5d059ca1462000a13578063a8b2216e1462000a38578063a9059cbb1462000a5f578063ab8f6ffe1462000a8457005b80639bb851c011620001d15780639bb851c014620009995780639c98e47114620009b15780639dc29fac14620009c9578063a515366a14620009ee57005b8063872cf05914620008fd5780638bac7dad146200091557806395d89b41146200093a57806396b477cb146200096857005b806340c10f1911620002fb57806370a08231116200028b578063819b64631162000259578063819b6463146200087c578063833b1fce146200089357806384467fdb14620008b3578063852c484914620008d857005b806370a082311462000728578063731b3a03146200076257806377e741c7146200077957806379502c55146200079e57005b8063520fdbbc11620002c9578063520fdbbc14620006895780635f7d394914620006ae5780636b5f444c14620006ec5780636fd2c80b146200071157005b806340c10f191462000600578063436459691462000625578063496ccd9b146200063c5780634bb278f3146200066157005b80631604e416116200037757806329070c6d116200034557806329070c6d146200058e5780632f2c3f2e14620005a5578063313ce56714620005bd57806335be16e014620005db57005b80631604e416146200050657806318160ddd146200051e5780631904bb2e146200053557806323b872dd146200056957005b80630d8e6e2c11620003b55780630d8e6e2c14620004845780631122063314620004a5578063114eaf5514620004bc5780631250a28d14620004e157005b806306fdde0314620003e5578063095ea7b314620004295780630ae65e7a146200045f57005b36620003e357005b005b348015620003f257600080fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b6040516200042091906200541e565b60405180910390f35b3480156200043657600080fd5b506200044e6200044836600462005450565b62000d6d565b604051901515815260200162000420565b3480156200046c57600080fd5b50620003e36200047e3660046200547f565b62000d86565b3480156200049157600080fd5b50601c545b60405190815260200162000420565b348015620004b257600080fd5b50600e5462000496565b348015620004c957600080fd5b50620003e3620004db3660046200549f565b62000e19565b348015620004ee57600080fd5b50620003e3620005003660046200547f565b62000e4b565b3480156200051357600080fd5b506200049660245481565b3480156200052b57600080fd5b5060295462000496565b3480156200054257600080fd5b506200055a620005543660046200547f565b62000e9a565b604051620004209190620054f2565b3480156200057657600080fd5b506200044e620005883660046200566b565b62001148565b3480156200059b57600080fd5b50600d5462000496565b348015620005b257600080fd5b506200049661271081565b348015620005ca57600080fd5b506040516012815260200162000420565b348015620005e857600080fd5b50620003e3620005fa366004620056b1565b620011a2565b3480156200060d57600080fd5b50620003e36200061f36600462005450565b6200137e565b3480156200063257600080fd5b50601a5462000496565b3480156200064957600080fd5b50620003e36200065b3660046200547f565b62001438565b3480156200066e57600080fd5b506200067962001542565b6040516200042092919062005772565b3480156200069657600080fd5b50620003e3620006a83660046200547f565b62001974565b348015620006bb57600080fd5b50620006d3620006cd3660046200578f565b62001b6f565b6040516001600160a01b03909116815260200162000420565b348015620006f957600080fd5b50620003e36200070b3660046200549f565b62001d86565b3480156200071e57600080fd5b5060105462000496565b3480156200073557600080fd5b5062000496620007473660046200547f565b6001600160a01b031660009081526027602052604090205490565b3480156200076f57600080fd5b5060205462000496565b3480156200078657600080fd5b50620003e3620007983660046200549f565b62001f0c565b348015620007ab57600080fd5b506040805160a08082018352600d548252600e54602080840191909152600f54838501526010546060808501919091526011546001600160a01b03908116608080870191909152865160c0810188526012548316815260135483168186015260145483168189015260155483168185015260165483168183015260175483169581019590955286519081018752601854909116815260195492810192909252601a5494820194909452601b5493810193909352601c546200086a939084565b604051620004209493929190620057b2565b3480156200088957600080fd5b50601b5462000496565b348015620008a057600080fd5b506013546001600160a01b0316620006d3565b348015620008c057600080fd5b50620003e3620008d23660046200592c565b62001f3e565b348015620008e557600080fd5b50620003e3620008f736600462005450565b6200206d565b3480156200090a57600080fd5b50620003e3620021fe565b3480156200092257600080fd5b50620003e3620009343660046200549f565b6200223a565b3480156200094757600080fd5b50604080518082019091526003815262272a2760e91b602082015262000411565b3480156200097557600080fd5b5062000496620009873660046200549f565b6000908152601f602052604090205490565b348015620009a657600080fd5b506200049660235481565b348015620009be57600080fd5b506200049660215481565b348015620009d657600080fd5b50620003e3620009e836600462005450565b620022be565b348015620009fb57600080fd5b50620003e362000a0d36600462005450565b620023d4565b34801562000a2057600080fd5b50620003e362000a3236600462005450565b620024a7565b34801562000a4557600080fd5b5062000a506200253e565b604051620004209190620059d5565b34801562000a6c57600080fd5b506200044e62000a7e36600462005450565b62002621565b34801562000a9157600080fd5b5062000a9c62002630565b60405162000420919062005a3b565b34801562000ab857600080fd5b5062000ac362002742565b6040516200042092919062005a96565b34801562000ae057600080fd5b50620003e362000af236600462005ac8565b62002a0e565b34801562000b0557600080fd5b50620003e362000b173660046200547f565b62002a55565b34801562000b2a57600080fd5b50620003e362000b3c3660046200547f565b62002aa4565b34801562000b4f57600080fd5b5062000b5a62002d27565b6040516200042092919062005b33565b34801562000b7757600080fd5b5062000b8262002e5e565b60405162000420919062005b5c565b34801562000b9e57600080fd5b506200049660205481565b34801562000bb657600080fd5b5062000496601e5481565b34801562000bce57600080fd5b50620003e362000be03660046200549f565b62002ec2565b34801562000bf357600080fd5b50620003e362000c053660046200547f565b62002f25565b34801562000c1857600080fd5b50620003e362002f74565b34801562000c3057600080fd5b50620003e362000c423660046200547f565b62002fc8565b34801562000c5557600080fd5b50620003e362000c673660046200547f565b62003017565b34801562000c7a57600080fd5b50602a54620006d3906001600160a01b031681565b34801562000c9c57600080fd5b50620003e362003066565b34801562000cb457600080fd5b50620003e362000cc63660046200547f565b620030a7565b34801562000cd957600080fd5b506200049662000ceb36600462005b71565b6001600160a01b03918216600090815260266020908152604080832093909416825291909152205490565b34801562000d2357600080fd5b5060195462000496565b34801562000d3a57600080fd5b506018546001600160a01b0316620006d3565b34801562000d5a57600080fd5b506011546001600160a01b0316620006d3565b600062000d7c338484620030f6565b5060015b92915050565b6001600160a01b038082166000818152602860205260409020600101549091161462000dcf5760405162461bcd60e51b815260040162000dc69062005baf565b60405180910390fd5b6001600160a01b0381811660009081526028602052604090205416331462000e0b5760405162461bcd60e51b815260040162000dc69062005be6565b62000e16816200321f565b50565b6018546001600160a01b0316331462000e465760405162461bcd60e51b815260040162000dc69062005c32565b601055565b6018546001600160a01b0316331462000e785760405162461bcd60e51b815260040162000dc69062005c32565b601280546001600160a01b0319166001600160a01b0392909216919091179055565b62000ea4620051f7565b6001600160a01b038083166000818152602860205260409020600101549091161462000ee45760405162461bcd60e51b815260040162000dc69062005c69565b6001600160a01b03808316600090815260286020908152604091829020825161028081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000f479062005ca0565b80601f016020809104026020016040519081016040528092919081815260200182805462000f759062005ca0565b801562000fc65780601f1062000f9a5761010080835404028352916020019162000fc6565b820191906000526020600020905b81548152906001019060200180831162000fa857829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546080820152600882015460a0820152600982015460c0820152600a82015460e0820152600b820154610100820152600c8201546001600160a01b0316610120820152600d820154610140820152600e820154610160820152600f82015461018082015260108201546101a082015260118201546101c08201526012820180546101e090920191620010829062005ca0565b80601f0160208091040260200160405190810160405280929190818152602001828054620010b09062005ca0565b8015620011015780601f10620010d55761010080835404028352916020019162001101565b820191906000526020600020905b815481529060010190602001808311620010e357829003601f168201915b5050509183525050601382015460209091019060ff1660038111156200112b576200112b620054b9565b60038111156200113f576200113f620054b9565b90525092915050565b600062001157848484620032f6565b6001600160a01b03841660009081526026602090815260408083203384529091528120546200118890849062005cf2565b905062001197853383620030f6565b506001949350505050565b6012546001600160a01b031633146200120a5760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f742074686520736c617368696e6720636f6e74726044820152621858dd60ea1b606482015260840162000dc6565b60006101208201356028826200122760408601602087016200547f565b6001600160a01b03166001600160a01b031681526020019081526020016000206009015462001257919062005cf2565b60c0830135602860006200127260408701602088016200547f565b6001600160a01b03166001600160a01b0316815260200190815260200160002060060154620012a2919062005cf2565b60a084013560286000620012bd60408801602089016200547f565b6001600160a01b03166001600160a01b0316815260200190815260200160002060050154620012ed919062005cf2565b620012f9919062005d08565b62001305919062005d08565b6011546001600160a01b03166000908152602760205260408120805492935083929091906200133690849062005d08565b90915550829050602860006200135360408401602085016200547f565b6001600160a01b03168152602081019190915260400160002062001378828262005f12565b50505050565b6018546001600160a01b03163314620013ab5760405162461bcd60e51b815260040162000dc69062005c32565b6001600160a01b03821660009081526027602052604081208054839290620013d590849062005d08565b925050819055508060296000828254620013f0919062005d08565b90915550506040518181526001600160a01b038316907f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf906020015b60405180910390a25050565b6018546001600160a01b03163314620014655760405162461bcd60e51b815260040162000dc69062005c32565b601380546001600160a01b0319166001600160a01b03838116918217909255601454604051637adbf97360e01b8152600481019290925290911690637adbf97390602401600060405180830381600087803b158015620014c457600080fd5b505af1158015620014d9573d6000803e3d6000fd5b5050601654604051637adbf97360e01b81526001600160a01b0385811660048301529091169250637adbf97391506024015b600060405180830381600087803b1580156200152657600080fd5b505af11580156200153b573d6000803e3d6000fd5b5050505050565b602a546000906060906001600160a01b03163314620015755760405162461bcd60e51b815260040162000dc69062006068565b601e54436000818152601f60209081526040822093909355601954925490926200159f9162005d08565b6012546040516306c9789b60e41b8152929091146004830181905292506001600160a01b031690636c9789b090602401600060405180830381600087803b158015620015ea57600080fd5b505af1158015620015ff573d6000803e3d6000fd5b5050505080156200175f576200161462003428565b6200161e62003812565b6200162862003901565b6000806200163562002742565b60135460405163422811f960e11b81529294509092506001600160a01b03169063845023f2906200166b90859060040162005b5c565b600060405180830381600087803b1580156200168657600080fd5b505af11580156200169b573d6000803e3d6000fd5b505060125460405163e08b14ed60e01b81526001600160a01b03909116925063e08b14ed9150620016d190849060040162005b5c565b600060405180830381600087803b158015620016ec57600080fd5b505af115801562001701573d6000803e3d6000fd5b50505050436020819055506001601e600082825462001721919062005d08565b9091555050601e546040519081527febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e3359060200160405180910390a150505b60135460408051634bb278f360e01b815290516000926001600160a01b031691634bb278f3916004808301926020929190829003018187875af1158015620017ab573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620017d19190620060ab565b905080156200184d576014546040805163a2e6204560e01b815290516001600160a01b039092169163a2e620459160048082019260209290919082900301816000875af192505050801562001845575060408051601f3d908101601f191682019092526200184291810190620060ab565b60015b156200184d57505b600254602280546040805160208084028201810190925282815260ff9094169391839160009084015b8282101562001964576000848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620018cc9062005ca0565b80601f0160208091040260200160405190810160405280929190818152602001828054620018fa9062005ca0565b80156200194b5780601f106200191f576101008083540402835291602001916200194b565b820191906000526020600020905b8154815290600101906020018083116200192d57829003601f168201915b5050505050815250508152602001906001019062001876565b5050505090509350935050509091565b6018546001600160a01b03163314620019a15760405162461bcd60e51b815260040162000dc69062005c32565b601880546001600160a01b0319166001600160a01b0383811691821790925560135460405163b3ab15fb60e01b815260048101929092529091169063b3ab15fb90602401600060405180830381600087803b15801562001a0057600080fd5b505af115801562001a15573d6000803e3d6000fd5b505060145460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001a6157600080fd5b505af115801562001a76573d6000803e3d6000fd5b505060155460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001ac257600080fd5b505af115801562001ad7573d6000803e3d6000fd5b505060165460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001b2357600080fd5b505af115801562001b38573d6000803e3d6000fd5b505060175460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb91506024016200150b565b600080805b60225481101562001bcb576022818154811062001b955762001b95620060cf565b9060005260206000209060030201600101548262001bb4919062005d08565b91508062001bc281620060e5565b91505062001b74565b508060000362001c1e5760405162461bcd60e51b815260206004820152601c60248201527f54686520636f6d6d6974746565206973206e6f74207374616b696e6700000000604482015260640162000dc6565b60008362001c2e60638762006101565b62001c3a919062005d08565b905060008160405160200162001c5291815260200190565b60408051601f1981840301815291905280516020909101209050600062001c7a848362006131565b90506000805b60225481101562001d2a576022818154811062001ca15762001ca1620060cf565b9060005260206000209060030201600101548262001cc0919062005d08565b915062001ccf60018362005cf2565b831162001d15576022818154811062001cec5762001cec620060cf565b60009182526020909120600390910201546001600160a01b0316965062000d8095505050505050565b8062001d2181620060e5565b91505062001c80565b5060405162461bcd60e51b815260206004820152602960248201527f5468657265206973206e6f2076616c696461746f72206c65667420696e20746860448201526865206e6574776f726b60b81b606482015260840162000dc6565b6018546001600160a01b0316331462001db35760405162461bcd60e51b815260040162000dc69062005c32565b60195481101562001e6a578060205462001dce919062005d08565b431062001e6a5760405162461bcd60e51b815260206004820152605760248201527f63757272656e7420636861696e2068656164206578636565642074686520776960448201527f6e646f773a206c617374426c6f636b45706f6368202b205f6e6577506572696f60648201527f642c2074727920616761696e206c6174746572206f6e2e000000000000000000608482015260a40162000dc6565b6019819055601254604051631ad7d11360e21b8152600481018390526001600160a01b0390911690636b5f444c90602401600060405180830381600087803b15801562001eb657600080fd5b505af115801562001ecb573d6000803e3d6000fd5b505050507fd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f818160405162001f0191815260200190565b60405180910390a150565b6018546001600160a01b0316331462001f395760405162461bcd60e51b815260040162000dc69062005c32565b600d55565b6000604051806102800160405280336001600160a01b0316815260200160006001600160a01b03168152602001856001600160a01b03168152602001868152602001600d6000016002015481526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b0316815260200160008152602001438152602001600081526020016000815260200160008152602001848152602001600060038111156200200d576200200d620054b9565b905290506200201d818362003a20565b60208101516101808201516040517f8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c926200205e92339289918b9162006148565b60405180910390a15050505050565b6001600160a01b0380831660008181526028602052604090206001015490911614620020ad5760405162461bcd60e51b815260040162000dc69062005baf565b6001600160a01b03828116600090815260286020526040902054163314620020e95760405162461bcd60e51b815260040162000dc69062005be6565b6127108111156200213d5760405162461bcd60e51b815260206004820152601f60248201527f7265717569726520636f727265637420636f6d6d697373696f6e207261746500604482015260640162000dc6565b604080516060810182526001600160a01b038481168252436020808401918252838501868152600c80546000908152600a909352958220855181546001600160a01b031916951694909417845591516001808501919091559151600290930192909255835492939092909190620021b690849062005d08565b90915550506040518281526001600160a01b038416907f4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf9060200160405180910390a2505050565b6018546001600160a01b031633146200222b5760405162461bcd60e51b815260040162000dc69062005c32565b6002805460ff19166001179055565b6018546001600160a01b03163314620022675760405162461bcd60e51b815260040162000dc69062005c32565b60008111620022b95760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000dc6565b601b55565b6018546001600160a01b03163314620022eb5760405162461bcd60e51b815260040162000dc69062005c32565b6001600160a01b0382166000908152602760205260409020548111156200234e5760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000dc6565b6001600160a01b038216600090815260276020526040812080548392906200237890849062005cf2565b92505081905550806029600082825462002393919062005cf2565b90915550506040518181526001600160a01b038316907f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3906020016200142c565b6001600160a01b0380831660008181526028602052604090206001015490911614620024145760405162461bcd60e51b815260040162000dc69062005c69565b6001600160a01b03821660009081526028602052604081206013015460ff166003811115620024475762002447620054b9565b14620024965760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000dc6565b620024a382823362003e7a565b5050565b6001600160a01b0380831660008181526028602052604090206001015490911614620024e75760405162461bcd60e51b815260040162000dc69062005c69565b60008111620025315760405162461bcd60e51b81526020600482015260156024820152740756e626f6e64696e6720616d6f756e74206973203605c1b604482015260640162000dc6565b620024a382823362004066565b60606025805480602002602001604051908101604052809291908181526020016000905b8282101562002618578382906000526020600020018054620025849062005ca0565b80601f0160208091040260200160405190810160405280929190818152602001828054620025b29062005ca0565b8015620026035780601f10620025d75761010080835404028352916020019162002603565b820191906000526020600020905b815481529060010190602001808311620025e557829003601f168201915b50505050508152602001906001019062002562565b50505050905090565b600062000d7c338484620032f6565b60606022805480602002602001604051908101604052809291908181526020016000905b8282101562002618576000848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620026aa9062005ca0565b80601f0160208091040260200160405190810160405280929190818152602001828054620026d89062005ca0565b8015620027295780601f10620026fd5761010080835404028352916020019162002729565b820191906000526020600020905b8154815290600101906020018083116200270b57829003601f168201915b5050505050815250508152602001906001019062002654565b602a5460609081906001600160a01b03163314620027745760405162461bcd60e51b815260040162000dc69062006068565b601d54620027c55760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000dc6565b620027cf620052c7565b601b546080820152601d8152602860208201526022604082015260216060820152620027fb8162004399565b6200280960256000620052e5565b60225480620028505760405162461bcd60e51b8152602060048201526012602482015271636f6d6d697474656520697320656d70747960701b604482015260640162000dc6565b60008167ffffffffffffffff8111156200286e576200286e62005881565b60405190808252806020026020018201604052801562002898578160200160208202803683370190505b50905060008267ffffffffffffffff811115620028b957620028b962005881565b604051908082528060200260200182016040528015620028e3578160200160208202803683370190505b50905060005b8381101562002a0257600060286000602284815481106200290e576200290e620060cf565b60009182526020808320600392830201546001600160a01b031684528301939093526040909101812060258054600181018255925292507f401968ff42a154441da5f6c4c935ac46b8671f0e062baaa62a7545ba53bb6e4c0190620029769083018262006191565b50600281015484516001600160a01b03909116908590849081106200299f576200299f620060cf565b6001600160a01b03928316602091820292909201015260018201548451911690849084908110620029d457620029d4620060cf565b6001600160a01b03909216602092830291909101909101525080620029f981620060e5565b915050620028e9565b50909450925050509091565b6018546001600160a01b0316331462002a3b5760405162461bcd60e51b815260040162000dc69062005c32565b62002a48600083620043b4565b620024a3600182620043b4565b6018546001600160a01b0316331462002a825760405162461bcd60e51b815260040162000dc69062005c32565b601580546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038082166000818152602860205260409020600101549091161462002ae45760405162461bcd60e51b815260040162000dc69062005baf565b6001600160a01b0380821660009081526028602052604090208054909116331462002b6a5760405162461bcd60e51b815260206004820152602f60248201527f726571756972652063616c6c657220746f2062652076616c696461746f72207460448201526e1c99585cdd5c9e481858d8dbdd5b9d608a1b606482015260840162000dc6565b6000601382015460ff16600381111562002b885762002b88620054b9565b0362002bd75760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f7220616c7265616479206163746976650000000000000000604482015260640162000dc6565b6002601382015460ff16600381111562002bf55762002bf5620054b9565b14801562002c065750438160100154115b1562002c555760405162461bcd60e51b815260206004820152601760248201527f76616c696461746f72207374696c6c20696e206a61696c000000000000000000604482015260640162000dc6565b6003601382015460ff16600381111562002c735762002c73620054b9565b0362002cc25760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206a61696c6564207065726d616e656e746c7900000000604482015260640162000dc6565b60138101805460ff1916905580546019546020546001600160a01b038581169316917f60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b59162002d12919062005d08565b60405190815260200160405180910390a35050565b6060806000600181805462002d3c9062005ca0565b80601f016020809104026020016040519081016040528092919081815260200182805462002d6a9062005ca0565b801562002dbb5780601f1062002d8f5761010080835404028352916020019162002dbb565b820191906000526020600020905b81548152906001019060200180831162002d9d57829003601f168201915b5050505050915080805462002dd09062005ca0565b80601f016020809104026020016040519081016040528092919081815260200182805462002dfe9062005ca0565b801562002e4f5780601f1062002e235761010080835404028352916020019162002e4f565b820191906000526020600020905b81548152906001019060200180831162002e3157829003601f168201915b50505050509050915091509091565b6060601d80548060200260200160405190810160405280929190818152602001828054801562002eb857602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831162002e99575b5050505050905090565b6018546001600160a01b0316331462002eef5760405162461bcd60e51b815260040162000dc69062005c32565b600e8190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200162001f01565b6018546001600160a01b0316331462002f525760405162461bcd60e51b815260040162000dc69062005c32565b601780546001600160a01b0319166001600160a01b0392909216919091179055565b6018546001600160a01b0316331462002fa15760405162461bcd60e51b815260040162000dc69062005c32565b62002fae60008062005305565b62002fbc6001600062005305565b6002805460ff19169055565b6018546001600160a01b0316331462002ff55760405162461bcd60e51b815260040162000dc69062005c32565b601680546001600160a01b0319166001600160a01b0392909216919091179055565b6018546001600160a01b03163314620030445760405162461bcd60e51b815260040162000dc69062005c32565b601480546001600160a01b0319166001600160a01b0392909216919091179055565b602a546001600160a01b03163314620030935760405162461bcd60e51b815260040162000dc69062006068565b6200309d62003812565b620024a362002742565b6018546001600160a01b03163314620030d45760405162461bcd60e51b815260040162000dc69062005c32565b601180546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b0383166200315a5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000dc6565b6001600160a01b038216620031bd5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000dc6565b6001600160a01b0383811660008181526026602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b038116600090815260286020526040812090601382015460ff166003811115620032545762003254620054b9565b14620032a35760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265206163746976650000000000000000604482015260640162000dc6565b60138101805460ff1916600117905580546019546020546001600160a01b038581169316917f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c9162002d12919062005d08565b6001600160a01b038316600090815260276020526040902054811115620033595760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000dc6565b6001600160a01b038316600090815260276020526040812080548392906200338390849062005cf2565b9250508190555080602960008282546200339e919062005d08565b90915550620033b1905081600262006101565b6001600160a01b03831660009081526027602052604081208054909190620033db90849062005d08565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200321291815260200190565b476000036200343357565b600d544790600090670de0b6b3a7640000906200345290849062006101565b6200345e91906200627a565b90508015620034dc576011546040516000916001600160a01b03169083908381818185875af1925050503d8060008114620034b6576040519150601f19603f3d011682016040523d82523d6000602084013e620034bb565b606091505b5090915050801515600103620034da57620034d7828462005cf2565b92505b505b8160236000828254620034f0919062005d08565b90915550600090505b6022548110156200380d5760006028600060228481548110620035205762003520620060cf565b600091825260208083206003909202909101546001600160a01b031683528201929092526040018120602154602280549294509091879190869081106200356b576200356b620060cf565b90600052602060002090600302016001015462003589919062006101565b6200359591906200627a565b90508015620037f5576002601383015460ff166003811115620035bc57620035bc620054b9565b1480620035e457506003601383015460ff166003811115620035e257620035e2620054b9565b145b156200368c57601254602280546001600160a01b0390921691631de9d9b691849187908110620036185762003618620060cf565b600091825260209091206003909102015460405160e084901b6001600160e01b03191681526001600160a01b0390911660048201526024016000604051808303818588803b1580156200366a57600080fd5b505af11580156200367f573d6000803e3d6000fd5b50505050505050620037f8565b60008260050154828460080154620036a5919062006101565b620036b191906200627a565b90506000620036c1828462005cf2565b90508115620037265783546040516001600160a01b03909116906108fc9084906000818181858888f193505050503d80600081146200371d576040519150601f19603f3d011682016040523d82523d6000602084013e62003722565b606091505b5050505b8015620037ae5783600c0160009054906101000a90046001600160a01b03166001600160a01b031663fb489a7b826040518263ffffffff1660e01b815260040160206040518083038185885af115801562003785573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190620037ac919062006291565b505b60018401546040518481526001600160a01b03909116907fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe5639060200160405180910390a250505b50505b806200380481620060e5565b915050620034f9565b505050565b6004545b6005548110156200383f5762003839816200383181620060e5565b925062004507565b62003816565b50600554600455600854600754036200385457565b6009545b60085481101562003881576200387b816200387381620060e5565b92506200472e565b62003858565b50600854600955600754805b600854811015620038fb576010546000828152600660205260409020600401544391620038ba9162005d08565b11620038e057620038cb8162004a29565b620038d860018362005d08565b9150620038e6565b620038fb565b80620038f281620060e5565b9150506200388d565b50600755565b600c54600b54101562003a1e57600b546000908152600a6020526040902060105460018201544391620039349162005d08565b11156200393e5750565b600281015481546001600160a01b03908116600090815260286020526040808220600490810185905585548416835291819020600c015490516319fac8fd60e01b81529216926319fac8fd9262003999920190815260200190565b600060405180830381600087803b158015620039b457600080fd5b505af1158015620039c9573d6000803e3d6000fd5b5050600b80546000908152600a6020526040812080546001600160a01b03191681556001808201839055600290910182905582549094509192509062003a1190849062005d08565b9091555062003901915050565b565b60e281511462003a6a5760405162461bcd60e51b8152602060048201526014602482015273092dcecc2d8d2c840e0e4dedecc40d8cadccee8d60631b604482015260640162000dc6565b6030826102400151511462003ac25760405162461bcd60e51b815260206004820152601c60248201527f496e76616c696420636f6e73656e737573206b6579206c656e67746800000000604482015260640162000dc6565b62003acd8262004b74565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a00000000000060208083019190915284519251919260009262003b30920160609190911b6bffffffffffffffffffffffff1916815260140190565b604051602081830303815290604052905060008262003b50835162004ca1565b8360405160200162003b6593929190620062ab565b60408051601f198184030181528282528051602091820120600280855260608501845290945060009392909183019080368337019050509050600080808062003bbd898262003bb76041600262006101565b62004dc2565b9050600062003bdc8a62003bd46041600262006101565b606062004dc2565b905060205b825181101562003cad5762003bf7838262004edb565b6040805160008152602081018083528d905260ff8316918101919091526060810184905260808101839052929850909650945060019060a0016020604051602081039080840390855afa15801562003c53573d6000803e3d6000fd5b5050604051601f19015190508762003c6d6041846200627a565b8151811062003c805762003c80620060cf565b6001600160a01b039092166020928302919091019091015262003ca560418262005d08565b905062003be1565b508a602001516001600160a01b03168660008151811062003cd25762003cd2620060cf565b60200260200101516001600160a01b03161462003d445760405162461bcd60e51b815260206004820152602960248201527f496e76616c6964206e6f6465206b6579206f776e6572736869702070726f6f66604482015268081c1c9bdd9a59195960ba1b606482015260840162000dc6565b8a604001516001600160a01b03168660018151811062003d685762003d68620060cf565b60200260200101516001600160a01b03161462003ddc5760405162461bcd60e51b815260206004820152602b60248201527f496e76616c6964206f7261636c65206b6579206f776e6572736869702070726f60448201526a1bd9881c1c9bdd9a59195960aa1b606482015260840162000dc6565b600162003df48c6102400151838e6000015162004f12565b1462003e625760405162461bcd60e51b815260206004820152603660248201527f496e76616c696420636f6e73656e737573206b6579206f776e65727368697020604482015275383937b7b3103337b9103932b3b4b9ba3930ba34b7b760511b606482015260840162000dc6565b62003e6d8b62004f81565b5050505050505050505050565b6000821162003ed85760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000dc6565b6001600160a01b03811660009081526027602052604090205482111562003f425760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000dc6565b6001600160a01b0381166000908152602760205260408120805484929062003f6c90849062005cf2565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580546000908152600394859052978820875181549088166001600160a01b031991821617825595516001820180549190981696169590951790955590516002840155519101558054919262003ff883620060e5565b90915550506001600160a01b03848116600081815260286020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d91015b60405180910390a35050505050565b6001600160a01b0380841660009081526028602052604090208054909183811691161480620041dc57600c820154604051631092ab9160e31b81526001600160a01b03858116600483015260009216906384955c8890602401602060405180830381865afa158015620040dd573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062004103919062006291565b9050848110156200416b5760405162461bcd60e51b815260206004820152602b60248201527f696e73756666696369656e7420756e6c6f636b6564204c6971756964204e657760448201526a746f6e2062616c616e636560a81b606482015260840162000dc6565b600c83015460405163282d3fdf60e01b81526001600160a01b038681166004830152602482018890529091169063282d3fdf90604401600060405180830381600087803b158015620041bc57600080fd5b505af1158015620041d1573d6000803e3d6000fd5b50505050506200426f565b8382600b01548360080154620041f3919062005cf2565b1015620042535760405162461bcd60e51b815260206004820152602760248201527f696e73756666696369656e742073656c6620626f6e646564206e6577746f6e2060448201526662616c616e636560c81b606482015260840162000dc6565b8382600b01600082825462004269919062005d08565b90915550505b6040805160e0810182526001600160a01b0380861682528781166020808401918252838501898152600060608601818152436080880190815260a088018381528a151560c08a019081526008805486526006909752998420985189549089166001600160a01b0319918216178a55965160018a01805491909916971696909617909655915160028701559051600386015592516004850155905160059093018054945115156101000261ff00199415159490941661ffff1990951694909417929092179092558054916200434383620060e5565b9190505550826001600160a01b0316856001600160a01b03167f63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc8387604051620040579291909115158252602082015260400190565b60fa60a06000808285855af46200380d573d6000803e3d6000fd5b81546002600180831615610100020382160482518082016020811060208410016002811462004463576001811462004489578660005260208404602060002001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b808210156200444c57815184556001840193506020820191506200442d565b815191036101000a908190040290915550620044fe565b60028302826020036101000a846020036101000a602089015104020185018755620044fe565b8660005260208404602060002001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b80821015620044ec5781518455600184019350602082019150620044cd565b815191036101000a9081900402909155505b50505050505050565b600081815260036020908152604080832060018101546001600160a01b03168452602890925282209091601382015460ff1660038111156200454d576200454d620054b9565b14620045ed57600282015482546001600160a01b0316600090815260276020526040812080549091906200458390849062005d08565b909155505081546001830154600284015460138401546040517f1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f87834294620045e0946001600160a01b0391821694911692909160ff90911690620062f4565b60405180910390a1505050565b805482546001600160a01b03908116911614620046ea57600080826008015483600501546200461d919062005cf2565b9050806000036200463557836002015491506200465b565b80846002015484600d01546200464c919062006101565b6200465891906200627a565b91505b600c83015484546040516340c10f1960e01b81526001600160a01b039182166004820152602481018590529116906340c10f1990604401600060405180830381600087803b158015620046ad57600080fd5b505af1158015620046c2573d6000803e3d6000fd5b505050508183600d016000828254620046dc919062005d08565b909155506200470a92505050565b816002015481600801600082825462004704919062005d08565b90915550505b816002015481600501600082825462004724919062005d08565b9091555050505050565b600081815260066020908152604080832060018101546001600160a01b0316845260289092528220600582015491929091610100900460ff1662004922576002830154600c8301548454604051637eee288d60e01b81526001600160a01b03918216600482015260248101849052911690637eee288d90604401600060405180830381600087803b158015620047c357600080fd5b505af1158015620047d8573d6000803e3d6000fd5b50505050600c8301548454604051632770a7eb60e21b81526001600160a01b03918216600482015260248101849052911690639dc29fac90604401600060405180830381600087803b1580156200482e57600080fd5b505af115801562004843573d6000803e3d6000fd5b505050506000836008015484600501546200485f919062005cf2565b600d85015490915062004873828462006101565b6200487f91906200627a565b92508184600d01600082825462004897919062005cf2565b90915550506006840154600003620048b65760038501839055620048de565b60068401546007850154620048cc908562006101565b620048d891906200627a565b60038601555b82846006016000828254620048f4919062005d08565b909155505060038501546007850180546000906200491490849062005d08565b90915550620049f792505050565b50600282015460088201548111156200493c575060088101545b81600901546000036200495657600383018190556200497e565b6009820154600a8301546200496c908362006101565b6200497891906200627a565b60038401555b8082600901600082825462004994919062005d08565b90915550506003830154600a83018054600090620049b490849062005d08565b9250508190555080826008016000828254620049d1919062005cf2565b90915550506002830154600b83018054600090620049f190849062005cf2565b90915550505b6005808401805460ff191660011790558201805482919060009062004a1e90849062005cf2565b909155505050505050565b6000818152600660205260408120600381015490910362004a48575050565b60018101546001600160a01b031660009081526028602052604081206005830154909190610100900460ff1662004ae45781600701548260060154846003015462004a94919062006101565b62004aa091906200627a565b90508082600601600082825462004ab8919062005cf2565b9091555050600383015460078301805460009062004ad890849062005cf2565b9091555062004b499050565b81600a01548260090154846003015462004aff919062006101565b62004b0b91906200627a565b90508082600901600082825462004b23919062005cf2565b90915550506003830154600a8301805460009062004b4390849062005cf2565b90915550505b82546001600160a01b03166000908152602760205260408120805483929062004a1e90849062005d08565b600062004b858260600151620051b2565b6001600160a01b0390911660208401529050801562004bd55760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000dc6565b6020808301516001600160a01b0390811660009081526028909252604090912060010154161562004c495760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000dc6565b61271082608001511115620024a35760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000dc6565b60608160000362004cc95750506040805180820190915260018152600360fc1b602082015290565b8160005b811562004cf9578062004ce081620060e5565b915062004cf19050600a836200627a565b915062004ccd565b60008167ffffffffffffffff81111562004d175762004d1762005881565b6040519080825280601f01601f19166020018201604052801562004d42576020820181803683370190505b5090505b841562004dba5762004d5a60018362005cf2565b915062004d69600a8662006131565b62004d7690603062005d08565b60f81b81838151811062004d8e5762004d8e620060cf565b60200101906001600160f81b031916908160001a90535062004db2600a866200627a565b945062004d46565b949350505050565b60608162004dd281601f62005d08565b101562004e135760405162461bcd60e51b815260206004820152600e60248201526d736c6963655f6f766572666c6f7760901b604482015260640162000dc6565b62004e1f828462005d08565b8451101562004e655760405162461bcd60e51b8152602060048201526011602482015270736c6963655f6f75744f66426f756e647360781b604482015260640162000dc6565b60608215801562004e86576040519150600082526020820160405262004ed2565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101562004ec157805183526020928301920162004ea7565b5050858452601f01601f1916604052505b50949350505050565b8181018051602082015160409092015190919060001a601b81101562004f0b5762004f08601b8262006322565b90505b9250925092565b600062004f1e62005344565b600085858560405160200162004f37939291906200633e565b6040516020818303038152906040529050600060fb905060008251602062004f60919062005d08565b90506020848285855afa62004f7457600080fd5b5050905195945050505050565b6101808101516001600160a01b03166200500457601d5460009062004fa69062004ca1565b90508160200151826000015183608001518360405162004fc69062005362565b62004fd594939291906200638d565b604051809103906000f08015801562004ff2573d6000803e3d6000fd5b506001600160a01b0316610180830152505b60208181018051601d80546001808201835560009283527f6d4407e7be21f808e6509aa9fa9143369579dd7d760fe20a2c09680fc146134f90910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602890955260409182902086518154871690851617815593519084018054861691841691909117905584015160028301805490941691161790915560608201518291906003820190620050b69082620063cc565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620051809082620063cc565b5061026082015160138201805460ff19166001836003811115620051a857620051a8620054b9565b0217905550505050565b600080620051bf62005370565b60008060ff9050604083875160208901845afa620051dc57600080fd5b50508051602090910151600160601b90910494909350915050565b60405180610280016040528060006001600160a01b0316815260200160006001600160a01b0316815260200160006001600160a01b0316815260200160608152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b0316815260200160008152602001600081526020016000815260200160008152602001600081526020016060815260200160006003811115620052c257620052c2620054b9565b905290565b6040518060a001604052806005906020820280368337509192915050565b508054600082559060005260206000209081019062000e1691906200538e565b508054620053139062005ca0565b6000825580601f1062005324575050565b601f01602090049060005260206000209081019062000e169190620053b3565b60405180602001604052806001906020820280368337509192915050565b611531806200648d83390190565b60405180604001604052806002906020820280368337509192915050565b80821115620053af576000620053a5828262005305565b506001016200538e565b5090565b5b80821115620053af5760008155600101620053b4565b60005b83811015620053e7578181015183820152602001620053cd565b50506000910152565b600081518084526200540a816020860160208601620053ca565b601f01601f19169290920160200192915050565b602081526000620054336020830184620053f0565b9392505050565b6001600160a01b038116811462000e1657600080fd5b600080604083850312156200546457600080fd5b823562005471816200543a565b946020939093013593505050565b6000602082840312156200549257600080fd5b813562005433816200543a565b600060208284031215620054b257600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b60048110620054ee57634e487b7160e01b600052602160045260246000fd5b9052565b602081526200550d6020820183516001600160a01b03169052565b600060208301516200552a60408401826001600160a01b03169052565b5060408301516001600160a01b03811660608401525060608301516102808060808501526200555e6102a0850183620053f0565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a0620055e6818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f1901610260808801919091529091506200564a8483620053f0565b9350808701519150506200566182860182620054cf565b5090949350505050565b6000806000606084860312156200568157600080fd5b83356200568e816200543a565b92506020840135620056a0816200543a565b929592945050506040919091013590565b600060208284031215620056c457600080fd5b813567ffffffffffffffff811115620056dc57600080fd5b820161028081850312156200543357600080fd5b600081518084526020808501808196508360051b8101915082860160005b8581101562005765578284038952815180516001600160a01b0316855285810151868601526040908101516060918601829052906200575081870183620053f0565b9a87019a95505050908401906001016200570e565b5091979650505050505050565b821515815260406020820152600062004dba6040830184620056f0565b60008060408385031215620057a357600080fd5b50508035926020909101359150565b845181526020808601518183015260408087015181840152606080880151818501526080808901516001600160a01b03908116828701528851811660a08088019190915294890151811660c087015292880151831660e086015290870151821661010085015286015116610120830152840151610200820190620058426101408401826001600160a01b03169052565b5083516001600160a01b0316610160830152602084015161018083015260408401516101a08301526060909301516101c08201526101e0015292915050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620058a957600080fd5b813567ffffffffffffffff80821115620058c757620058c762005881565b604051601f8301601f19908116603f01168101908282118183101715620058f257620058f262005881565b816040528381528660208588010111156200590c57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080608085870312156200594357600080fd5b843567ffffffffffffffff808211156200595c57600080fd5b6200596a8883890162005897565b9550602087013591506200597e826200543a565b909350604086013590808211156200599557600080fd5b620059a38883890162005897565b93506060870135915080821115620059ba57600080fd5b50620059c98782880162005897565b91505092959194509250565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b8281101562005a2e57603f1988860301845262005a1b858351620053f0565b94509285019290850190600101620059fc565b5092979650505050505050565b602081526000620054336020830184620056f0565b600081518084526020808501945080840160005b8381101562005a8b5781516001600160a01b03168752958201959082019060010162005a64565b509495945050505050565b60408152600062005aab604083018562005a50565b828103602084015262005abf818562005a50565b95945050505050565b6000806040838503121562005adc57600080fd5b823567ffffffffffffffff8082111562005af557600080fd5b62005b038683870162005897565b9350602085013591508082111562005b1a57600080fd5b5062005b298582860162005897565b9150509250929050565b60408152600062005b486040830185620053f0565b828103602084015262005abf8185620053f0565b60208152600062005433602083018462005a50565b6000806040838503121562005b8557600080fd5b823562005b92816200543a565b9150602083013562005ba4816200543a565b809150509250929050565b6020808252601c908201527f76616c696461746f72206d757374206265207265676973746572656400000000604082015260600190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b60208082526018908201527f76616c696461746f72206e6f7420726567697374657265640000000000000000604082015260600190565b600181811c9082168062005cb557607f821691505b60208210810362005cd657634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000d805762000d8062005cdc565b8082018082111562000d805762000d8062005cdc565b6000813562000d80816200543a565b80546001600160a01b0319166001600160a01b0392909216919091179055565b6000808335601e1984360301811262005d6557600080fd5b83018035915067ffffffffffffffff82111562005d8157600080fd5b60200191503681900382131562005d9757600080fd5b9250929050565b601f8211156200380d57600081815260208120601f850160051c8101602086101562005dc75750805b601f850160051c820191505b8181101562005de85782815560010162005dd3565b505050505050565b600019600383901b1c191660019190911b1790565b67ffffffffffffffff83111562005e205762005e2062005881565b62005e388362005e31835462005ca0565b8362005d9e565b6000601f84116001811462005e6b576000851562005e565750838201355b62005e62868262005df0565b8455506200153b565b600083815260209020601f19861690835b8281101562005e9e578685013582556020948501946001909201910162005e7c565b508682101562005ebc5760001960f88860031b161c19848701351681555b505060018560011b0183555050505050565b600081356004811062000d8057600080fd5b6004821062005eff57634e487b7160e01b600052602160045260246000fd5b60ff1981541660ff831681178255505050565b62005f2862005f218362005d1e565b8262005d2d565b62005f4462005f3a6020840162005d1e565b6001830162005d2d565b62005f6062005f566040840162005d1e565b6002830162005d2d565b62005f6f606083018362005d4d565b62005f7f81836003860162005e05565b50506080820135600482015560a0820135600582015560c0820135600682015560e0820135600782015561010082013560088201556101208201356009820155610140820135600a820155610160820135600b82015562005ff262005fe8610180840162005d1e565b600c830162005d2d565b6101a0820135600d8201556101c0820135600e8201556101e0820135600f820155610200820135601082015561022082013560118201556200603961024083018362005d4d565b6200604981836012860162005e05565b5050620024a36200605e610260840162005ece565b6013830162005ee0565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b600060208284031215620060be57600080fd5b815180151581146200543357600080fd5b634e487b7160e01b600052603260045260246000fd5b600060018201620060fa57620060fa62005cdc565b5060010190565b808202811582820484141762000d805762000d8062005cdc565b634e487b7160e01b600052601260045260246000fd5b6000826200614357620061436200611b565b500690565b600060018060a01b0380881683528087166020840152808616604084015260a060608401526200617c60a0840186620053f0565b91508084166080840152509695505050505050565b8181036200619d575050565b620061a9825462005ca0565b67ffffffffffffffff811115620061c457620061c462005881565b620061dc81620061d5845462005ca0565b8462005d9e565b6000601f8211600181146200620f5760008315620061fa5750848201545b62006206848262005df0565b8555506200153b565b600085815260209020601f19841690600086815260209020845b838110156200624b578286015482556001958601959091019060200162006229565b50858310156200626a5781850154600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000826200628c576200628c6200611b565b500490565b600060208284031215620062a457600080fd5b5051919050565b60008451620062bf818460208901620053ca565b845190830190620062d5818360208901620053ca565b8451910190620062ea818360208801620053ca565b0195945050505050565b6001600160a01b03858116825284166020820152604081018390526080810162005abf6060830184620054cf565b60ff818116838216019081111562000d805762000d8062005cdc565b6000845162006352818460208901620053ca565b84519083019062006368818360208901620053ca565b60609490941b6bffffffffffffffffffffffff19169301928352505060140192915050565b6001600160a01b0385811682528416602082015260408101839052608060608201819052600090620063c290830184620053f0565b9695505050505050565b815167ffffffffffffffff811115620063e957620063e962005881565b620063fa81620061d5845462005ca0565b602080601f8311600181146200642e5760008415620064195750858301515b62006425858262005df0565b86555062005de8565b600085815260208120601f198616915b828110156200645f578886015182559484019460019091019084016200643e565b50858210156200626a57939096015160001960f8600387901b161c19169092555050600190811b0190555056fe60806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220e1e1658d868fb41912eb4d5ecdf5bed6ec95bd7db81bea3e11358a02e3aa535c64736f6c63430008150033a2646970667358221220e56302c787f4ddad24b291019fb7508bc51df5e921c6d39ac4ddfdff6344cbe864736f6c6343000815003360806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220e1e1658d868fb41912eb4d5ecdf5bed6ec95bd7db81bea3e11358a02e3aa535c64736f6c63430008150033",
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
	return address, gasConsumed, &AutonityUpgradeTest{contract: c}, nil
}

// AutonityUpgradeTest is an auto generated Go binding around an Ethereum contract.
type AutonityUpgradeTest struct {
	*contract
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTest) COMMISSIONRATEPRECISION(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) Allowance(opts *runOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _addr) view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTest) BalanceOf(opts *runOptions, _addr common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "balanceOf", _addr)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// Config is a free data retrieval call binding the contract method 0x79502c55.
//
// Solidity: function config() view returns((uint256,uint256,uint256,uint256,address) policy, (address,address,address,address,address,address) contracts, (address,uint256,uint256,uint256) protocol, uint256 contractVersion)
func (_AutonityUpgradeTest *AutonityUpgradeTest) Config(opts *runOptions) (struct {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) Decimals(opts *runOptions) (uint8, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) Deployer(opts *runOptions) (common.Address, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) EpochID(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) EpochReward(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) EpochTotalBondedStake(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetBlockPeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetCommittee(opts *runOptions) ([]AutonityCommitteeMember, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetCommitteeEnodes(opts *runOptions) ([]string, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetEpochFromBlock(opts *runOptions, _block *big.Int) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetEpochPeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetLastEpochBlock(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetMaxCommitteeSize(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetMinimumBaseFee(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetNewContract(opts *runOptions) ([]byte, string, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetOperator(opts *runOptions) (common.Address, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetOracle(opts *runOptions) (common.Address, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetProposer(opts *runOptions, height *big.Int, round *big.Int) (common.Address, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetTreasuryAccount(opts *runOptions) (common.Address, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetTreasuryFee(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetUnbondingPeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetValidator(opts *runOptions, _addr common.Address) (AutonityValidator, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetValidators(opts *runOptions) ([]common.Address, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) GetVersion(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "getVersion")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// LastEpochBlock is a free data retrieval call binding the contract method 0xc2362dd5.
//
// Solidity: function lastEpochBlock() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTest) LastEpochBlock(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) Name(opts *runOptions) (string, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) Symbol(opts *runOptions) (string, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "symbol")

	if err != nil {
		return *new(string), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, consumed, err

}

// TotalRedistributed is a free data retrieval call binding the contract method 0x9bb851c0.
//
// Solidity: function totalRedistributed() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTest) TotalRedistributed(opts *runOptions) (*big.Int, uint64, error) {
	out, consumed, err := _AutonityUpgradeTest.call(opts, "totalRedistributed")

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AutonityUpgradeTest *AutonityUpgradeTest) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
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
func (_AutonityUpgradeTest *AutonityUpgradeTest) ActivateValidator(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "activateValidator", _address)
	return consumed, err
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_AutonityUpgradeTest *AutonityUpgradeTest) Approve(opts *runOptions, spender common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "approve", spender, amount)
	return consumed, err
}

// Bond is a paid mutator transaction binding the contract method 0xa515366a.
//
// Solidity: function bond(address _validator, uint256 _amount) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) Bond(opts *runOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "bond", _validator, _amount)
	return consumed, err
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _addr, uint256 _amount) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) Burn(opts *runOptions, _addr common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "burn", _addr, _amount)
	return consumed, err
}

// ChangeCommissionRate is a paid mutator transaction binding the contract method 0x852c4849.
//
// Solidity: function changeCommissionRate(address _validator, uint256 _rate) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) ChangeCommissionRate(opts *runOptions, _validator common.Address, _rate *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "changeCommissionRate", _validator, _rate)
	return consumed, err
}

// CompleteContractUpgrade is a paid mutator transaction binding the contract method 0x872cf059.
//
// Solidity: function completeContractUpgrade() returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) CompleteContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "completeContractUpgrade")
	return consumed, err
}

// ComputeCommittee is a paid mutator transaction binding the contract method 0xae1f5fa0.
//
// Solidity: function computeCommittee() returns(address[], address[])
func (_AutonityUpgradeTest *AutonityUpgradeTest) ComputeCommittee(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "computeCommittee")
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256,bytes)[])
func (_AutonityUpgradeTest *AutonityUpgradeTest) Finalize(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "finalize")
	return consumed, err
}

// FinalizeInitialization is a paid mutator transaction binding the contract method 0xd861b0e8.
//
// Solidity: function finalizeInitialization() returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) FinalizeInitialization(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "finalizeInitialization")
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _addr, uint256 _amount) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) Mint(opts *runOptions, _addr common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "mint", _addr, _amount)
	return consumed, err
}

// PauseValidator is a paid mutator transaction binding the contract method 0x0ae65e7a.
//
// Solidity: function pauseValidator(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) PauseValidator(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "pauseValidator", _address)
	return consumed, err
}

// RegisterValidator is a paid mutator transaction binding the contract method 0x84467fdb.
//
// Solidity: function registerValidator(string _enode, address _oracleAddress, bytes _consensusKey, bytes _signatures) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) RegisterValidator(opts *runOptions, _enode string, _oracleAddress common.Address, _consensusKey []byte, _signatures []byte) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "registerValidator", _enode, _oracleAddress, _consensusKey, _signatures)
	return consumed, err
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) ResetContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "resetContractUpgrade")
	return consumed, err
}

// SetAccountabilityContract is a paid mutator transaction binding the contract method 0x1250a28d.
//
// Solidity: function setAccountabilityContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetAccountabilityContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setAccountabilityContract", _address)
	return consumed, err
}

// SetAcuContract is a paid mutator transaction binding the contract method 0xd372c07e.
//
// Solidity: function setAcuContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetAcuContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setAcuContract", _address)
	return consumed, err
}

// SetCommitteeSize is a paid mutator transaction binding the contract method 0x8bac7dad.
//
// Solidity: function setCommitteeSize(uint256 _size) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetCommitteeSize(opts *runOptions, _size *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setCommitteeSize", _size)
	return consumed, err
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _period) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetEpochPeriod(opts *runOptions, _period *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setEpochPeriod", _period)
	return consumed, err
}

// SetMinimumBaseFee is a paid mutator transaction binding the contract method 0xcb696f54.
//
// Solidity: function setMinimumBaseFee(uint256 _price) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetMinimumBaseFee(opts *runOptions, _price *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setMinimumBaseFee", _price)
	return consumed, err
}

// SetOperatorAccount is a paid mutator transaction binding the contract method 0x520fdbbc.
//
// Solidity: function setOperatorAccount(address _account) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetOperatorAccount(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setOperatorAccount", _account)
	return consumed, err
}

// SetOracleContract is a paid mutator transaction binding the contract method 0x496ccd9b.
//
// Solidity: function setOracleContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetOracleContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setOracleContract", _address)
	return consumed, err
}

// SetStabilizationContract is a paid mutator transaction binding the contract method 0xcfd19fb9.
//
// Solidity: function setStabilizationContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetStabilizationContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setStabilizationContract", _address)
	return consumed, err
}

// SetSupplyControlContract is a paid mutator transaction binding the contract method 0xb3ecbadd.
//
// Solidity: function setSupplyControlContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetSupplyControlContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setSupplyControlContract", _address)
	return consumed, err
}

// SetTreasuryAccount is a paid mutator transaction binding the contract method 0xd886f8a2.
//
// Solidity: function setTreasuryAccount(address _account) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetTreasuryAccount(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setTreasuryAccount", _account)
	return consumed, err
}

// SetTreasuryFee is a paid mutator transaction binding the contract method 0x77e741c7.
//
// Solidity: function setTreasuryFee(uint256 _treasuryFee) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetTreasuryFee(opts *runOptions, _treasuryFee *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setTreasuryFee", _treasuryFee)
	return consumed, err
}

// SetUnbondingPeriod is a paid mutator transaction binding the contract method 0x114eaf55.
//
// Solidity: function setUnbondingPeriod(uint256 _period) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetUnbondingPeriod(opts *runOptions, _period *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setUnbondingPeriod", _period)
	return consumed, err
}

// SetUpgradeManagerContract is a paid mutator transaction binding the contract method 0xceaad455.
//
// Solidity: function setUpgradeManagerContract(address _address) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) SetUpgradeManagerContract(opts *runOptions, _address common.Address) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "setUpgradeManagerContract", _address)
	return consumed, err
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _recipient, uint256 _amount) returns(bool)
func (_AutonityUpgradeTest *AutonityUpgradeTest) Transfer(opts *runOptions, _recipient common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "transfer", _recipient, _amount)
	return consumed, err
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_AutonityUpgradeTest *AutonityUpgradeTest) TransferFrom(opts *runOptions, sender common.Address, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "transferFrom", sender, recipient, amount)
	return consumed, err
}

// Unbond is a paid mutator transaction binding the contract method 0xa5d059ca.
//
// Solidity: function unbond(address _validator, uint256 _amount) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) Unbond(opts *runOptions, _validator common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "unbond", _validator, _amount)
	return consumed, err
}

// UpdateValidatorAndTransferSlashedFunds is a paid mutator transaction binding the contract method 0x35be16e0.
//
// Solidity: function updateValidatorAndTransferSlashedFunds((address,address,address,string,uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256,address,uint256,uint256,uint256,uint256,uint256,bytes,uint8) _val) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) UpdateValidatorAndTransferSlashedFunds(opts *runOptions, _val AutonityValidator) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "updateValidatorAndTransferSlashedFunds", _val)
	return consumed, err
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) UpgradeContract(opts *runOptions, _bytecode []byte, _abi string) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "upgradeContract", _bytecode, _abi)
	return consumed, err
}

// Fallback is a paid mutator transaction binding the contract fallback function.
// WARNING! UNTESTED
// Solidity: fallback() payable returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) Fallback(opts *runOptions, calldata []byte) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "", calldata)
	return consumed, err
}

// Receive is a paid mutator transaction binding the contract receive function.
// WARNING! UNTESTED
// Solidity: receive() payable returns()
func (_AutonityUpgradeTest *AutonityUpgradeTest) Receive(opts *runOptions) (uint64, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterActivatedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*AutonityUpgradeTestActivatedValidatorIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchActivatedValidator(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestActivatedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseActivatedValidator(log types.Log) (*AutonityUpgradeTestActivatedValidator, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AutonityUpgradeTestApprovalIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchApproval(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseApproval(log types.Log) (*AutonityUpgradeTestApproval, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterBondingRejected(opts *bind.FilterOpts) (*AutonityUpgradeTestBondingRejectedIterator, error) {






			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "BondingRejected")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestBondingRejectedIterator{contract: _AutonityUpgradeTest.contract, event: "BondingRejected", logs: logs, sub: sub}, nil
 		}

		// WatchBondingRejected is a free log subscription operation binding the contract event 0x1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f878342.
		//
		// Solidity: event BondingRejected(address delegator, address delegatee, uint256 amount, uint8 state)
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchBondingRejected(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestBondingRejected) (event.Subscription, error) {






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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseBondingRejected(log types.Log) (*AutonityUpgradeTestBondingRejected, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterBurnedStake(opts *bind.FilterOpts, addr []common.Address) (*AutonityUpgradeTestBurnedStakeIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchBurnedStake(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestBurnedStake, addr []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseBurnedStake(log types.Log) (*AutonityUpgradeTestBurnedStake, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterCommissionRateChange(opts *bind.FilterOpts, validator []common.Address) (*AutonityUpgradeTestCommissionRateChangeIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchCommissionRateChange(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestCommissionRateChange, validator []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseCommissionRateChange(log types.Log) (*AutonityUpgradeTestCommissionRateChange, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterEpochPeriodUpdated(opts *bind.FilterOpts) (*AutonityUpgradeTestEpochPeriodUpdatedIterator, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "EpochPeriodUpdated")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestEpochPeriodUpdatedIterator{contract: _AutonityUpgradeTest.contract, event: "EpochPeriodUpdated", logs: logs, sub: sub}, nil
 		}

		// WatchEpochPeriodUpdated is a free log subscription operation binding the contract event 0xd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f81.
		//
		// Solidity: event EpochPeriodUpdated(uint256 period)
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchEpochPeriodUpdated(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestEpochPeriodUpdated) (event.Subscription, error) {



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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseEpochPeriodUpdated(log types.Log) (*AutonityUpgradeTestEpochPeriodUpdated, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterMinimumBaseFeeUpdated(opts *bind.FilterOpts) (*AutonityUpgradeTestMinimumBaseFeeUpdatedIterator, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "MinimumBaseFeeUpdated")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestMinimumBaseFeeUpdatedIterator{contract: _AutonityUpgradeTest.contract, event: "MinimumBaseFeeUpdated", logs: logs, sub: sub}, nil
 		}

		// WatchMinimumBaseFeeUpdated is a free log subscription operation binding the contract event 0x1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd389128.
		//
		// Solidity: event MinimumBaseFeeUpdated(uint256 gasPrice)
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchMinimumBaseFeeUpdated(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestMinimumBaseFeeUpdated) (event.Subscription, error) {



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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseMinimumBaseFeeUpdated(log types.Log) (*AutonityUpgradeTestMinimumBaseFeeUpdated, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterMintedStake(opts *bind.FilterOpts, addr []common.Address) (*AutonityUpgradeTestMintedStakeIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchMintedStake(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestMintedStake, addr []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseMintedStake(log types.Log) (*AutonityUpgradeTestMintedStake, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterNewBondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*AutonityUpgradeTestNewBondingRequestIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchNewBondingRequest(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestNewBondingRequest, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseNewBondingRequest(log types.Log) (*AutonityUpgradeTestNewBondingRequest, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterNewEpoch(opts *bind.FilterOpts) (*AutonityUpgradeTestNewEpochIterator, error) {



			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "NewEpoch")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestNewEpochIterator{contract: _AutonityUpgradeTest.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
 		}

		// WatchNewEpoch is a free log subscription operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
		//
		// Solidity: event NewEpoch(uint256 epoch)
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestNewEpoch) (event.Subscription, error) {



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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseNewEpoch(log types.Log) (*AutonityUpgradeTestNewEpoch, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterNewUnbondingRequest(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*AutonityUpgradeTestNewUnbondingRequestIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchNewUnbondingRequest(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestNewUnbondingRequest, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseNewUnbondingRequest(log types.Log) (*AutonityUpgradeTestNewUnbondingRequest, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterPausedValidator(opts *bind.FilterOpts, treasury []common.Address, addr []common.Address) (*AutonityUpgradeTestPausedValidatorIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchPausedValidator(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestPausedValidator, treasury []common.Address, addr []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParsePausedValidator(log types.Log) (*AutonityUpgradeTestPausedValidator, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterRegisteredValidator(opts *bind.FilterOpts) (*AutonityUpgradeTestRegisteredValidatorIterator, error) {







			logs, sub, err := _AutonityUpgradeTest.contract.FilterLogs(opts, "RegisteredValidator")
			if err != nil {
				return nil, err
			}
			return &AutonityUpgradeTestRegisteredValidatorIterator{contract: _AutonityUpgradeTest.contract, event: "RegisteredValidator", logs: logs, sub: sub}, nil
 		}

		// WatchRegisteredValidator is a free log subscription operation binding the contract event 0x8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c.
		//
		// Solidity: event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract)
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchRegisteredValidator(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestRegisteredValidator) (event.Subscription, error) {







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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseRegisteredValidator(log types.Log) (*AutonityUpgradeTestRegisteredValidator, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterRewarded(opts *bind.FilterOpts, addr []common.Address) (*AutonityUpgradeTestRewardedIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchRewarded(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestRewarded, addr []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseRewarded(log types.Log) (*AutonityUpgradeTestRewarded, error) {
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
 		func (_AutonityUpgradeTest *AutonityUpgradeTest) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AutonityUpgradeTestTransferIterator, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AutonityUpgradeTestTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

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
		func (_AutonityUpgradeTest *AutonityUpgradeTest) ParseTransfer(log types.Log) (*AutonityUpgradeTestTransfer, error) {
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
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220831cbf1d093207c822ca0bb8903ee29676712d9a6abe4b340d9c5562c8a851e464736f6c63430008150033",
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
	return address, gasConsumed, &BytesLib{contract: c}, nil
}

// BytesLib is an auto generated Go binding around an Ethereum contract.
type BytesLib struct {
	*contract
}

/* EVENTS ARE NOT YET SUPPORTED

 */

// HelpersMetaData contains all meta data concerning the Helpers contract.
var HelpersMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220f08ec662b3647da9f8d9adbbfb5ae571333af0be6a27cf20a23d7dccab133b1b64736f6c63430008150033",
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
	return address, gasConsumed, &Helpers{contract: c}, nil
}

// Helpers is an auto generated Go binding around an Ethereum contract.
type Helpers struct {
	*contract
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
	*contract
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IACU *IACU) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _IACU.call(opts, "setOperator", operator)
	return consumed, err
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IACU *IACU) SetOracle(opts *runOptions, oracle common.Address) (uint64, error) {
	_, consumed, err := _IACU.call(opts, "setOracle", oracle)
	return consumed, err
}

// Update is a paid mutator transaction binding the contract method 0xa2e62045.
//
// Solidity: function update() returns(bool status)
func (_IACU *IACU) Update(opts *runOptions) (uint64, error) {
	_, consumed, err := _IACU.call(opts, "update")
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

 */

// IAccountabilityMetaData contains all meta data concerning the IAccountability contract.
var IAccountabilityMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"InnocenceProven\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewAccusation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_offender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_severity\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"NewFaultProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"releaseBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isJailbound\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"SlashingEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"distributeRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_epochEnd\",\"type\":\"bool\"}],\"name\":\"finalize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_committee\",\"type\":\"address[]\"}],\"name\":\"setCommittee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_newPeriod\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1de9d9b6": "distributeRewards(address)",
		"6c9789b0": "finalize(bool)",
		"e08b14ed": "setCommittee(address[])",
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
	*contract
}

// DistributeRewards is a paid mutator transaction binding the contract method 0x1de9d9b6.
//
// Solidity: function distributeRewards(address _validator) payable returns()
func (_IAccountability *IAccountability) DistributeRewards(opts *runOptions, _validator common.Address) (uint64, error) {
	_, consumed, err := _IAccountability.call(opts, "distributeRewards", _validator)
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x6c9789b0.
//
// Solidity: function finalize(bool _epochEnd) returns()
func (_IAccountability *IAccountability) Finalize(opts *runOptions, _epochEnd bool) (uint64, error) {
	_, consumed, err := _IAccountability.call(opts, "finalize", _epochEnd)
	return consumed, err
}

// SetCommittee is a paid mutator transaction binding the contract method 0xe08b14ed.
//
// Solidity: function setCommittee(address[] _committee) returns()
func (_IAccountability *IAccountability) SetCommittee(opts *runOptions, _committee []common.Address) (uint64, error) {
	_, consumed, err := _IAccountability.call(opts, "setCommittee", _committee)
	return consumed, err
}

// SetEpochPeriod is a paid mutator transaction binding the contract method 0x6b5f444c.
//
// Solidity: function setEpochPeriod(uint256 _newPeriod) returns()
func (_IAccountability *IAccountability) SetEpochPeriod(opts *runOptions, _newPeriod *big.Int) (uint64, error) {
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
 		func (_IAccountability *IAccountability) FilterInnocenceProven(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityInnocenceProvenIterator, error) {

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
		func (_IAccountability *IAccountability) WatchInnocenceProven(opts *bind.WatchOpts, sink chan<- *IAccountabilityInnocenceProven, _offender []common.Address) (event.Subscription, error) {

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
		func (_IAccountability *IAccountability) ParseInnocenceProven(log types.Log) (*IAccountabilityInnocenceProven, error) {
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
 		func (_IAccountability *IAccountability) FilterNewAccusation(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityNewAccusationIterator, error) {

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
		func (_IAccountability *IAccountability) WatchNewAccusation(opts *bind.WatchOpts, sink chan<- *IAccountabilityNewAccusation, _offender []common.Address) (event.Subscription, error) {

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
		func (_IAccountability *IAccountability) ParseNewAccusation(log types.Log) (*IAccountabilityNewAccusation, error) {
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
 		func (_IAccountability *IAccountability) FilterNewFaultProof(opts *bind.FilterOpts, _offender []common.Address) (*IAccountabilityNewFaultProofIterator, error) {

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
		func (_IAccountability *IAccountability) WatchNewFaultProof(opts *bind.WatchOpts, sink chan<- *IAccountabilityNewFaultProof, _offender []common.Address) (event.Subscription, error) {

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
		func (_IAccountability *IAccountability) ParseNewFaultProof(log types.Log) (*IAccountabilityNewFaultProof, error) {
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
 		func (_IAccountability *IAccountability) FilterSlashingEvent(opts *bind.FilterOpts) (*IAccountabilitySlashingEventIterator, error) {







			logs, sub, err := _IAccountability.contract.FilterLogs(opts, "SlashingEvent")
			if err != nil {
				return nil, err
			}
			return &IAccountabilitySlashingEventIterator{contract: _IAccountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
 		}

		// WatchSlashingEvent is a free log subscription operation binding the contract event 0x6617e612ea2d01b5a235997fa4963b56b1097df6f968a82972433e9ff852e0f9.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId)
		func (_IAccountability *IAccountability) WatchSlashingEvent(opts *bind.WatchOpts, sink chan<- *IAccountabilitySlashingEvent) (event.Subscription, error) {







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
		func (_IAccountability *IAccountability) ParseSlashingEvent(log types.Log) (*IAccountabilitySlashingEvent, error) {
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
	*contract
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_IAutonity *IAutonity) GetOperator(opts *runOptions) (common.Address, uint64, error) {
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
func (_IAutonity *IAutonity) GetOracle(opts *runOptions) (common.Address, uint64, error) {
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
	*contract
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20) Allowance(opts *runOptions, owner common.Address, spender common.Address) (*big.Int, uint64, error) {
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
func (_IERC20 *IERC20) BalanceOf(opts *runOptions, account common.Address) (*big.Int, uint64, error) {
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
func (_IERC20 *IERC20) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
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
func (_IERC20 *IERC20) Approve(opts *runOptions, spender common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _IERC20.call(opts, "approve", spender, amount)
	return consumed, err
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20) Transfer(opts *runOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _IERC20.call(opts, "transfer", recipient, amount)
	return consumed, err
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20) TransferFrom(opts *runOptions, sender common.Address, recipient common.Address, amount *big.Int) (uint64, error) {
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
 		func (_IERC20 *IERC20) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IERC20ApprovalIterator, error) {

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
		func (_IERC20 *IERC20) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

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
		func (_IERC20 *IERC20) ParseApproval(log types.Log) (*IERC20Approval, error) {
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
 		func (_IERC20 *IERC20) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IERC20TransferIterator, error) {

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
		func (_IERC20 *IERC20) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

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
		func (_IERC20 *IERC20) ParseTransfer(log types.Log) (*IERC20Transfer, error) {
			event := new(IERC20Transfer)
			if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


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
	*contract
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() view returns(uint256)
func (_IOracle *IOracle) GetPrecision(opts *runOptions) (*big.Int, uint64, error) {
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
func (_IOracle *IOracle) GetRound(opts *runOptions) (*big.Int, uint64, error) {
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
func (_IOracle *IOracle) GetRoundData(opts *runOptions, _round *big.Int, _symbol string) (IOracleRoundData, uint64, error) {
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
func (_IOracle *IOracle) GetSymbols(opts *runOptions) ([]string, uint64, error) {
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
func (_IOracle *IOracle) GetVotePeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_IOracle *IOracle) GetVoters(opts *runOptions) ([]common.Address, uint64, error) {
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
func (_IOracle *IOracle) LatestRoundData(opts *runOptions, _symbol string) (IOracleRoundData, uint64, error) {
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
func (_IOracle *IOracle) Finalize(opts *runOptions) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "finalize")
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_IOracle *IOracle) SetOperator(opts *runOptions, _operator common.Address) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "setOperator", _operator)
	return consumed, err
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_IOracle *IOracle) SetSymbols(opts *runOptions, _symbols []string) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "setSymbols", _symbols)
	return consumed, err
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_IOracle *IOracle) SetVoters(opts *runOptions, _newVoters []common.Address) (uint64, error) {
	_, consumed, err := _IOracle.call(opts, "setVoters", _newVoters)
	return consumed, err
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_IOracle *IOracle) Vote(opts *runOptions, _commit *big.Int, _reports []*big.Int, _salt *big.Int) (uint64, error) {
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
 		func (_IOracle *IOracle) FilterNewRound(opts *bind.FilterOpts) (*IOracleNewRoundIterator, error) {






			logs, sub, err := _IOracle.contract.FilterLogs(opts, "NewRound")
			if err != nil {
				return nil, err
			}
			return &IOracleNewRoundIterator{contract: _IOracle.contract, event: "NewRound", logs: logs, sub: sub}, nil
 		}

		// WatchNewRound is a free log subscription operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
		//
		// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
		func (_IOracle *IOracle) WatchNewRound(opts *bind.WatchOpts, sink chan<- *IOracleNewRound) (event.Subscription, error) {






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
		func (_IOracle *IOracle) ParseNewRound(log types.Log) (*IOracleNewRound, error) {
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
 		func (_IOracle *IOracle) FilterNewSymbols(opts *bind.FilterOpts) (*IOracleNewSymbolsIterator, error) {




			logs, sub, err := _IOracle.contract.FilterLogs(opts, "NewSymbols")
			if err != nil {
				return nil, err
			}
			return &IOracleNewSymbolsIterator{contract: _IOracle.contract, event: "NewSymbols", logs: logs, sub: sub}, nil
 		}

		// WatchNewSymbols is a free log subscription operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
		//
		// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
		func (_IOracle *IOracle) WatchNewSymbols(opts *bind.WatchOpts, sink chan<- *IOracleNewSymbols) (event.Subscription, error) {




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
		func (_IOracle *IOracle) ParseNewSymbols(log types.Log) (*IOracleNewSymbols, error) {
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
 		func (_IOracle *IOracle) FilterVoted(opts *bind.FilterOpts, _voter []common.Address) (*IOracleVotedIterator, error) {

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
		func (_IOracle *IOracle) WatchVoted(opts *bind.WatchOpts, sink chan<- *IOracleVoted, _voter []common.Address) (event.Subscription, error) {

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
		func (_IOracle *IOracle) ParseVoted(log types.Log) (*IOracleVoted, error) {
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
	*contract
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_IStabilization *IStabilization) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _IStabilization.call(opts, "setOperator", operator)
	return consumed, err
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_IStabilization *IStabilization) SetOracle(opts *runOptions, oracle common.Address) (uint64, error) {
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
	*contract
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_ISupplyControl *ISupplyControl) AvailableSupply(opts *runOptions) (*big.Int, uint64, error) {
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
func (_ISupplyControl *ISupplyControl) Stabilizer(opts *runOptions) (common.Address, uint64, error) {
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
func (_ISupplyControl *ISupplyControl) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
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
func (_ISupplyControl *ISupplyControl) Burn(opts *runOptions) (uint64, error) {
	_, consumed, err := _ISupplyControl.call(opts, "burn")
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_ISupplyControl *ISupplyControl) Mint(opts *runOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _ISupplyControl.call(opts, "mint", recipient, amount)
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_ISupplyControl *ISupplyControl) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _ISupplyControl.call(opts, "setOperator", operator)
	return consumed, err
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_ISupplyControl *ISupplyControl) SetStabilizer(opts *runOptions, stabilizer_ common.Address) (uint64, error) {
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
 		func (_ISupplyControl *ISupplyControl) FilterBurn(opts *bind.FilterOpts) (*ISupplyControlBurnIterator, error) {



			logs, sub, err := _ISupplyControl.contract.FilterLogs(opts, "Burn")
			if err != nil {
				return nil, err
			}
			return &ISupplyControlBurnIterator{contract: _ISupplyControl.contract, event: "Burn", logs: logs, sub: sub}, nil
 		}

		// WatchBurn is a free log subscription operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
		//
		// Solidity: event Burn(uint256 amount)
		func (_ISupplyControl *ISupplyControl) WatchBurn(opts *bind.WatchOpts, sink chan<- *ISupplyControlBurn) (event.Subscription, error) {



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
		func (_ISupplyControl *ISupplyControl) ParseBurn(log types.Log) (*ISupplyControlBurn, error) {
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
 		func (_ISupplyControl *ISupplyControl) FilterMint(opts *bind.FilterOpts) (*ISupplyControlMintIterator, error) {




			logs, sub, err := _ISupplyControl.contract.FilterLogs(opts, "Mint")
			if err != nil {
				return nil, err
			}
			return &ISupplyControlMintIterator{contract: _ISupplyControl.contract, event: "Mint", logs: logs, sub: sub}, nil
 		}

		// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
		//
		// Solidity: event Mint(address recipient, uint256 amount)
		func (_ISupplyControl *ISupplyControl) WatchMint(opts *bind.WatchOpts, sink chan<- *ISupplyControlMint) (event.Subscription, error) {




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
		func (_ISupplyControl *ISupplyControl) ParseMint(log types.Log) (*ISupplyControlMint, error) {
			event := new(ISupplyControlMint)
			if err := _ISupplyControl.contract.UnpackLog(event, "Mint", log); err != nil {
				return nil, err
			}
			event.Raw = log
			return event, nil
		}


*/

// LiquidMetaData contains all meta data concerning the Liquid contract.
var LiquidMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"_treasury\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_index\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_FACTOR_UNIT_RECIP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"lockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"redistribute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"setCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"unclaimedRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_delegator\",\"type\":\"address\"}],\"name\":\"unlockedBalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
	Bin: "0x60806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220e1e1658d868fb41912eb4d5ecdf5bed6ec95bd7db81bea3e11358a02e3aa535c64736f6c63430008150033",
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
	return address, gasConsumed, &Liquid{contract: c}, nil
}

// Liquid is an auto generated Go binding around an Ethereum contract.
type Liquid struct {
	*contract
}

// COMMISSIONRATEPRECISION is a free data retrieval call binding the contract method 0x2f2c3f2e.
//
// Solidity: function COMMISSION_RATE_PRECISION() view returns(uint256)
func (_Liquid *Liquid) COMMISSIONRATEPRECISION(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Liquid *Liquid) FEEFACTORUNITRECIP(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Liquid *Liquid) Allowance(opts *runOptions, _owner common.Address, _spender common.Address) (*big.Int, uint64, error) {
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
func (_Liquid *Liquid) BalanceOf(opts *runOptions, _delegator common.Address) (*big.Int, uint64, error) {
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
func (_Liquid *Liquid) CommissionRate(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Liquid *Liquid) Decimals(opts *runOptions) (uint8, uint64, error) {
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
func (_Liquid *Liquid) LockedBalanceOf(opts *runOptions, _delegator common.Address) (*big.Int, uint64, error) {
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
func (_Liquid *Liquid) Name(opts *runOptions) (string, uint64, error) {
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
func (_Liquid *Liquid) Symbol(opts *runOptions) (string, uint64, error) {
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
func (_Liquid *Liquid) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Liquid *Liquid) Treasury(opts *runOptions) (common.Address, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "treasury")

	if err != nil {
		return *new(common.Address), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, consumed, err

}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address _account) view returns(uint256)
func (_Liquid *Liquid) UnclaimedRewards(opts *runOptions, _account common.Address) (*big.Int, uint64, error) {
	out, consumed, err := _Liquid.call(opts, "unclaimedRewards", _account)

	if err != nil {
		return *new(*big.Int), consumed, err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return out0, consumed, err

}

// UnlockedBalanceOf is a free data retrieval call binding the contract method 0x84955c88.
//
// Solidity: function unlockedBalanceOf(address _delegator) view returns(uint256)
func (_Liquid *Liquid) UnlockedBalanceOf(opts *runOptions, _delegator common.Address) (*big.Int, uint64, error) {
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
func (_Liquid *Liquid) Validator(opts *runOptions) (common.Address, uint64, error) {
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
func (_Liquid *Liquid) Approve(opts *runOptions, _spender common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "approve", _spender, _amount)
	return consumed, err
}

// Burn is a paid mutator transaction binding the contract method 0x9dc29fac.
//
// Solidity: function burn(address _account, uint256 _amount) returns()
func (_Liquid *Liquid) Burn(opts *runOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "burn", _account, _amount)
	return consumed, err
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_Liquid *Liquid) ClaimRewards(opts *runOptions) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "claimRewards")
	return consumed, err
}

// Lock is a paid mutator transaction binding the contract method 0x282d3fdf.
//
// Solidity: function lock(address _account, uint256 _amount) returns()
func (_Liquid *Liquid) Lock(opts *runOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "lock", _account, _amount)
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address _account, uint256 _amount) returns()
func (_Liquid *Liquid) Mint(opts *runOptions, _account common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "mint", _account, _amount)
	return consumed, err
}

// Redistribute is a paid mutator transaction binding the contract method 0xfb489a7b.
//
// Solidity: function redistribute() payable returns(uint256)
func (_Liquid *Liquid) Redistribute(opts *runOptions) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "redistribute")
	return consumed, err
}

// SetCommissionRate is a paid mutator transaction binding the contract method 0x19fac8fd.
//
// Solidity: function setCommissionRate(uint256 _rate) returns()
func (_Liquid *Liquid) SetCommissionRate(opts *runOptions, _rate *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "setCommissionRate", _rate)
	return consumed, err
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _amount) returns(bool _success)
func (_Liquid *Liquid) Transfer(opts *runOptions, _to common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "transfer", _to, _amount)
	return consumed, err
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _sender, address _recipient, uint256 _amount) returns(bool _success)
func (_Liquid *Liquid) TransferFrom(opts *runOptions, _sender common.Address, _recipient common.Address, _amount *big.Int) (uint64, error) {
	_, consumed, err := _Liquid.call(opts, "transferFrom", _sender, _recipient, _amount)
	return consumed, err
}

// Unlock is a paid mutator transaction binding the contract method 0x7eee288d.
//
// Solidity: function unlock(address _account, uint256 _amount) returns()
func (_Liquid *Liquid) Unlock(opts *runOptions, _account common.Address, _amount *big.Int) (uint64, error) {
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
 		func (_Liquid *Liquid) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*LiquidApprovalIterator, error) {

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
		func (_Liquid *Liquid) WatchApproval(opts *bind.WatchOpts, sink chan<- *LiquidApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

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
		func (_Liquid *Liquid) ParseApproval(log types.Log) (*LiquidApproval, error) {
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
 		func (_Liquid *Liquid) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*LiquidTransferIterator, error) {

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
		func (_Liquid *Liquid) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LiquidTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

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
		func (_Liquid *Liquid) ParseTransfer(log types.Log) (*LiquidTransfer, error) {
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
	Bin: "0x6080604052600160ff1b600755600160ff1b6008553480156200002157600080fd5b5060405162002df838038062002df8833981016040819052620000449162000639565b600280546001600160a01b038087166001600160a01b03199283161790925560038054928616929091169190911790558151620000899060009060208501906200035f565b5081516200009f9060019060208501906200035f565b5080600981905550620000c485600060018851620000be91906200074e565b62000181565b8451620000d9906004906020880190620003bc565b508451620000ef906005906020880190620003bc565b5060016006819055600d8054909101815560009081525b855181101562000175576001600b60008884815181106200012b576200012b6200076a565b6020908102919091018101516001600160a01b03168252810191909152604001600020600201805460ff1916911515919091179055806200016c8162000780565b91505062000106565b505050505050620009c3565b8082126200018e57505050565b81816000856002620001a185856200079c565b620001ad9190620007c6565b620001b9908762000806565b81518110620001cc57620001cc6200076a565b602002602001015190505b8183136200032b575b806001600160a01b0316868481518110620001ff57620001ff6200076a565b60200260200101516001600160a01b031610156200022c5782620002238162000831565b935050620001e0565b806001600160a01b03168683815181106200024b576200024b6200076a565b60200260200101516001600160a01b031611156200027857816200026f816200084c565b9250506200022c565b81831362000325578582815181106200029557620002956200076a565b6020026020010151868481518110620002b257620002b26200076a565b6020026020010151878581518110620002cf57620002cf6200076a565b60200260200101888581518110620002eb57620002eb6200076a565b6001600160a01b0393841660209182029290920101529116905282620003118162000831565b935050818062000321906200084c565b9250505b620001d7565b8185121562000341576200034186868462000181565b8383121562000357576200035786848662000181565b505050505050565b828054828255906000526020600020908101928215620003aa579160200282015b82811115620003aa5782518290620003999082620008f7565b509160200191906001019062000380565b50620003b892915062000422565b5090565b82805482825590600052602060002090810192821562000414579160200282015b828111156200041457825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620003dd565b50620003b892915062000443565b80821115620003b85760006200043982826200045a565b5060010162000422565b5b80821115620003b8576000815560010162000444565b50805462000468906200086c565b6000825580601f1062000479575050565b601f01602090049060005260206000209081019062000499919062000443565b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620004dd57620004dd6200049c565b604052919050565b60006001600160401b038211156200050157620005016200049c565b5060051b60200190565b80516001600160a01b03811681146200052357600080fd5b919050565b6000601f83818401126200053b57600080fd5b82516020620005546200054e83620004e5565b620004b2565b82815260059290921b850181019181810190878411156200057457600080fd5b8287015b848110156200062d5780516001600160401b03808211156200059a5760008081fd5b818a0191508a603f830112620005b05760008081fd5b8582015181811115620005c757620005c76200049c565b620005da818a01601f19168801620004b2565b915080825260408c81838601011115620005f45760008081fd5b60005b8281101562000614578481018201518482018a01528801620005f7565b5050600090820187015284525091830191830162000578565b50979650505050505050565b600080600080600060a086880312156200065257600080fd5b85516001600160401b03808211156200066a57600080fd5b818801915088601f8301126200067f57600080fd5b81516020620006926200054e83620004e5565b82815260059290921b8401810191818101908c841115620006b257600080fd5b948201945b83861015620006db57620006cb866200050b565b82529482019490820190620006b7565b9950620006ec90508a82016200050b565b97505050620006fe604089016200050b565b945060608801519150808211156200071557600080fd5b50620007248882890162000528565b925050608086015190509295509295909350565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000764576200076462000738565b92915050565b634e487b7160e01b600052603260045260246000fd5b60006001820162000795576200079562000738565b5060010190565b8181036000831280158383131683831282161715620007bf57620007bf62000738565b5092915050565b600082620007e457634e487b7160e01b600052601260045260246000fd5b600160ff1b82146000198414161562000801576200080162000738565b500590565b808201828112600083128015821682158216171562000829576200082962000738565b505092915050565b60006001600160ff1b01820162000795576200079562000738565b6000600160ff1b820162000864576200086462000738565b506000190190565b600181811c908216806200088157607f821691505b602082108103620008a257634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620008f257600081815260208120601f850160051c81016020861015620008d15750805b601f850160051c820191505b818110156200035757828155600101620008dd565b505050565b81516001600160401b038111156200091357620009136200049c565b6200092b816200092484546200086c565b84620008a8565b602080601f8311600181146200096357600084156200094a5750858301515b600019600386901b1c1916600185901b17855562000357565b600085815260208120601f198616915b82811015620009945788860151825594840194600190910190840162000973565b5085821015620009b35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61242580620009d36000396000f3fe6080604052600436106101225760003560e01c80638d4f75d2116100a5578063b3ab15fb1161006c578063b3ab15fb1461038a578063b78dec52146103aa578063ccce413b146103bf578063cdd72253146103df578063df7f710e14610401578063e6a02a281461042357005b80638d4f75d2146103135780639670c0bc146103335780639f8743f714610349578063a78135871461035e578063aa2f89b51461037457005b80634bb278f3116100e95780634bb278f3146101fd5780634c56ea56146102225780635281b5c61461026a5780635412b3ae14610297578063845023f2146102f357005b806308f21ff51461012b578063146ca53114610154578063307de9b61461016a57806333f98c771461018a5780633c8510fd146101dd57005b3661012957005b005b34801561013757600080fd5b5061014160085481565b6040519081526020015b60405180910390f35b34801561016057600080fd5b5061014160065481565b34801561017657600080fd5b50610129610185366004611a3a565b610439565b34801561019657600080fd5b506101aa6101a5366004611b76565b610680565b60405161014b91908151815260208083015190820152604080830151908201526060918201519181019190915260800190565b3480156101e957600080fd5b506101aa6101f8366004611bab565b6107a3565b34801561020957600080fd5b506102126108ad565b604051901515815260200161014b565b34801561022e57600080fd5b5061014161023d366004611c0e565b8151602081840181018051600c825292820194820194909420919093529091526000908152604090205481565b34801561027657600080fd5b5061028a610285366004611c5c565b610a50565b60405161014b9190611cc5565b3480156102a357600080fd5b506102d66102b2366004611cdf565b600b6020526000908152604090208054600182015460029092015490919060ff1683565b60408051938452602084019290925215159082015260600161014b565b3480156102ff57600080fd5b5061012961030e366004611d1e565b610afc565b34801561031f57600080fd5b5061012961032e366004611dbb565b610ba5565b34801561033f57600080fd5b5062989680610141565b34801561035557600080fd5b50600654610141565b34801561036a57600080fd5b5061014160095481565b34801561038057600080fd5b5061014160075481565b34801561039657600080fd5b506101296103a5366004611cdf565b610d15565b3480156103b657600080fd5b50600954610141565b3480156103cb57600080fd5b5061028a6103da366004611c5c565b610d61565b3480156103eb57600080fd5b506103f4610d71565b60405161014b9190611e6c565b34801561040d57600080fd5b50610416610dd3565b60405161014b9190611f0e565b34801561042f57600080fd5b50610141600a5481565b336000908152600b602052604090206002015460ff166104a05760405162461bcd60e51b815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f746572730000000000000060448201526064015b60405180910390fd5b600654336000908152600b6020526040902054036104f05760405162461bcd60e51b815260206004820152600d60248201526c185b1c9958591e481d9bdd1959609a1b6044820152606401610497565b336000908152600b602052604081206001810180549087905581546006549092559181900361052057505061067a565b600054841461053057505061067a565b600160065461053f9190611f37565b8114158061057c57508484843360405160200161055f9493929190611f4a565b6040516020818303038152906040528051906020012060001c8214155b156105f85760005b6000548110156105f0576001600160ff1b03600c600083815481106105ab576105ab611f95565b906000526020600020016040516105c29190611fe5565b90815260408051602092819003830190203360009081529252902055806105e88161205b565b915050610584565b50505061067a565b60005b848110156106765785858281811061061557610615611f95565b90506020020135600c6000838154811061063157610631611f95565b906000526020600020016040516106489190611fe5565b908152604080516020928190038301902033600090815292529020558061066e8161205b565b9150506105fb565b5050505b50505050565b6106ab6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000600d60016006546106be9190611f37565b815481106106ce576106ce611f95565b90600052602060002001836040516106e69190612074565b908152602001604051809103902060405180606001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff16600181111561073757610737612090565b600181111561074857610748612090565b81525050905060006040518060800160405280600160065461076a9190611f37565b815260200183600001518152602001836020015181526020018360400151600181111561079957610799612090565b9052949350505050565b6107ce6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000600d84815481106107e3576107e3611f95565b90600052602060002001836040516107fb9190612074565b908152602001604051809103902060405180606001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff16600181111561084c5761084c612090565b600181111561085d5761085d612090565b815250509050600060405180608001604052808681526020018360000151815260200183602001518152602001836040015160018111156108a0576108a0612090565b9052925050505b92915050565b6002546000906001600160a01b031633146108da5760405162461bcd60e51b8152600401610497906120a6565b600954600a546108ea91906120e9565b4310610a4a5760005b6000548110156109185761090681610f91565b6109116001826120e9565b90506108f3565b50600654600754036109955760005b600554811015610993576001600b60006005848154811061094a5761094a611f95565b6000918252602080832091909101546001600160a01b031683528201929092526040019020600201805460ff19169115159190911790558061098b8161205b565b915050610927565b505b6006546007546109a69060016120fc565b036109b3576109b36112c6565b43600a819055506001600660008282546109cd91906120e9565b90915550506008546109e09060026120fc565b600654036109fa57600180546109f89160009161188c565b505b60065460095460408051928352436020840152429083015260608201527fb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e59060800160405180910390a150600190565b50600090565b60018181548110610a6057600080fd5b906000526020600020016000915090508054610a7b90611fab565b80601f0160208091040260200160405190810160405280929190818152602001828054610aa790611fab565b8015610af45780601f10610ac957610100808354040283529160200191610af4565b820191906000526020600020905b815481529060010190602001808311610ad757829003601f168201915b505050505081565b6002546001600160a01b03163314610b265760405162461bcd60e51b8152600401610497906120a6565b8051600003610b6f5760405162461bcd60e51b8152602060048201526015602482015274566f746572732063616e277420626520656d70747960581b6044820152606401610497565b610b8881600060018451610b839190611f37565b6114a3565b8051610b9b9060059060208401906118e4565b5050600654600755565b6003546001600160a01b03163314610bf85760405162461bcd60e51b81526020600482015260166024820152753932b9ba3934b1ba32b2103a379037b832b930ba37b960511b6044820152606401610497565b8051600003610c425760405162461bcd60e51b815260206004820152601660248201527573796d626f6c732063616e277420626520656d70747960501b6044820152606401610497565b600654600854610c539060016120fc565b14158015610c65575060065460085414155b610cb15760405162461bcd60e51b815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e6400006044820152606401610497565b8051610cc4906001906020840190611945565b5060065460088190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d908290610cfc9060016120e9565b604051610d0a929190612124565b60405180910390a150565b6002546001600160a01b03163314610d3f5760405162461bcd60e51b8152600401610497906120a6565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b60008181548110610a6057600080fd5b60606005805480602002602001604051908101604052809291908181526020018280548015610dc957602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610dab575b5050505050905090565b60606006546008546001610de791906120fc565b03610ec3576001805480602002602001604051908101604052809291908181526020016000905b82821015610eba578382906000526020600020018054610e2d90611fab565b80601f0160208091040260200160405190810160405280929190818152602001828054610e5990611fab565b8015610ea65780601f10610e7b57610100808354040283529160200191610ea6565b820191906000526020600020905b815481529060010190602001808311610e8957829003601f168201915b505050505081526020019060010190610e0e565b50505050905090565b6000805480602002602001604051908101604052809291908181526020016000905b82821015610eba578382906000526020600020018054610f0490611fab565b80601f0160208091040260200160405190810160405280929190818152602001828054610f3090611fab565b8015610f7d5780601f10610f5257610100808354040283529160200191610f7d565b820191906000526020600020905b815481529060010190602001808311610f6057829003601f168201915b505050505081526020019060010190610ee5565b6000808281548110610fa557610fa5611f95565b906000526020600020018054610fba90611fab565b80601f0160208091040260200160405190810160405280929190818152602001828054610fe690611fab565b80156110335780601f1061100857610100808354040283529160200191611033565b820191906000526020600020905b81548152906001019060200180831161101657829003601f168201915b50505050509050600060048054905067ffffffffffffffff81111561105a5761105a611abf565b604051908082528060200260200182016040528015611083578160200160208202803683370190505b5090506000805b60045481101561119a576000600482815481106110a9576110a9611f95565b60009182526020808320909101546006546001600160a01b03909116808452600b90925260409092205490925014158061111d57506001600160ff1b03600c866040516110f69190612074565b90815260408051602092819003830190206001600160a01b03851660009081529252902054145b156111285750611188565b600c856040516111389190612074565b90815260408051602092819003830190206001600160a01b0384166000908152925290205484846111688161205b565b95508151811061117a5761117a611f95565b602002602001018181525050505b806111928161205b565b91505061108a565b506000600d60016006546111ae9190611f37565b815481106111be576111be611f95565b90600052602060002001846040516111d69190612074565b90815260405190819003602001902054905060018215611201576111fa8484611650565b9150600090505b600d805460019081018255600091909152604080516060810182528481524260208201529190820190839081111561123b5761123b612090565b815250600d6006548154811061125357611253611f95565b906000526020600020018660405161126b9190612074565b9081526020016040518091039020600082015181600001556020820151816001015560408201518160020160006101000a81548160ff021916908360018111156112b7576112b7612090565b02179055505050505050505050565b6000805b600454821080156112dc575060055481105b1561142157600581815481106112f4576112f4611f95565b600091825260209091200154600480546001600160a01b03909216918490811061132057611320611f95565b6000918252602090912001546001600160a01b03160361135a57816113448161205b565b92505080806113529061205b565b9150506112ca565b6005818154811061136d5761136d611f95565b600091825260209091200154600480546001600160a01b03909216918490811061139957611399611f95565b6000918252602090912001546001600160a01b0316101561141757600b6000600484815481106113cb576113cb611f95565b60009182526020808320909101546001600160a01b0316835282019290925260400181208181556001810191909155600201805460ff191690558161140f8161205b565b9250506112ca565b806113528161205b565b60045482101561148e57600b60006004848154811061144257611442611f95565b60009182526020808320909101546001600160a01b0316835282019290925260400181208181556001810191909155600201805460ff19169055816114868161205b565b925050611421565b6005805461149e9160049161198b565b505050565b8082126114af57505050565b818160008560026114c08585612146565b6114ca9190612183565b6114d490876120fc565b815181106114e4576114e4611f95565b602002602001015190505b818313611622575b806001600160a01b031686848151811061151357611513611f95565b60200260200101516001600160a01b0316101561153c5782611534816121b1565b9350506114f7565b806001600160a01b031686838151811061155857611558611f95565b60200260200101516001600160a01b031611156115815781611579816121c9565b92505061153c565b81831361161d5785828151811061159a5761159a611f95565b60200260200101518684815181106115b4576115b4611f95565b60200260200101518785815181106115ce576115ce611f95565b602002602001018885815181106115e7576115e7611f95565b6001600160a01b039384166020918202929092010152911690528261160b816121b1565b9350508180611619906121c9565b9250505b6114ef565b81851215611635576116358686846114a3565b83831215611648576116488684866114a3565b505050505050565b600081600003611662575060006108a7565b611678836000611673600186611f37565b611714565b60006116856002846121e6565b90506116926002846121fa565b156116b6578381815181106116a9576116a9611f95565b602002602001015161170c565b60028482815181106116ca576116ca611f95565b6020026020010151856001846116e09190611f37565b815181106116f0576116f0611f95565b602002602001015161170291906120fc565b61170c9190612183565b949350505050565b8181808203611724575050505050565b60008560026117338787612146565b61173d9190612183565b61174790876120fc565b8151811061175757611757611f95565b602002602001015190505b818313611866575b8086848151811061177d5761177d611f95565b6020026020010151121561179d5782611795816121b1565b93505061176a565b8582815181106117af576117af611f95565b60200260200101518112156117d057816117c8816121c9565b92505061179d565b818313611861578582815181106117e9576117e9611f95565b602002602001015186848151811061180357611803611f95565b602002602001015187858151811061181d5761181d611f95565b6020026020010188858151811061183657611836611f95565b6020908102919091010191909152528161184f816121c9565b925050828061185d906121b1565b9350505b611762565b8185121561187957611879868684611714565b8383121561164857611648868486611714565b8280548282559060005260206000209081019282156118d45760005260206000209182015b828111156118d457816118c48482612254565b50916001019190600101906118b1565b506118e09291506119cb565b5090565b828054828255906000526020600020908101928215611939579160200282015b8281111561193957825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611904565b506118e09291506119e8565b8280548282559060005260206000209081019282156118d4579160200282015b828111156118d4578251829061197b9082612335565b5091602001919060010190611965565b8280548282559060005260206000209081019282156119395760005260206000209182015b828111156119395782548255916001019190600101906119b0565b808211156118e05760006119df82826119fd565b506001016119cb565b5b808211156118e057600081556001016119e9565b508054611a0990611fab565b6000825580601f10611a19575050565b601f016020900490600052602060002090810190611a3791906119e8565b50565b60008060008060608587031215611a5057600080fd5b84359350602085013567ffffffffffffffff80821115611a6f57600080fd5b818701915087601f830112611a8357600080fd5b813581811115611a9257600080fd5b8860208260051b8501011115611aa757600080fd5b95986020929092019750949560400135945092505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715611afe57611afe611abf565b604052919050565b600082601f830112611b1757600080fd5b813567ffffffffffffffff811115611b3157611b31611abf565b611b44601f8201601f1916602001611ad5565b818152846020838601011115611b5957600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215611b8857600080fd5b813567ffffffffffffffff811115611b9f57600080fd5b61170c84828501611b06565b60008060408385031215611bbe57600080fd5b82359150602083013567ffffffffffffffff811115611bdc57600080fd5b611be885828601611b06565b9150509250929050565b80356001600160a01b0381168114611c0957600080fd5b919050565b60008060408385031215611c2157600080fd5b823567ffffffffffffffff811115611c3857600080fd5b611c4485828601611b06565b925050611c5360208401611bf2565b90509250929050565b600060208284031215611c6e57600080fd5b5035919050565b60005b83811015611c90578181015183820152602001611c78565b50506000910152565b60008151808452611cb1816020860160208601611c75565b601f01601f19169290920160200192915050565b602081526000611cd86020830184611c99565b9392505050565b600060208284031215611cf157600080fd5b611cd882611bf2565b600067ffffffffffffffff821115611d1457611d14611abf565b5060051b60200190565b60006020808385031215611d3157600080fd5b823567ffffffffffffffff811115611d4857600080fd5b8301601f81018513611d5957600080fd5b8035611d6c611d6782611cfa565b611ad5565b81815260059190911b82018301908381019087831115611d8b57600080fd5b928401925b82841015611db057611da184611bf2565b82529284019290840190611d90565b979650505050505050565b60006020808385031215611dce57600080fd5b823567ffffffffffffffff80821115611de657600080fd5b818501915085601f830112611dfa57600080fd5b8135611e08611d6782611cfa565b81815260059190911b83018401908481019088831115611e2757600080fd5b8585015b83811015611e5f57803585811115611e435760008081fd5b611e518b89838a0101611b06565b845250918601918601611e2b565b5098975050505050505050565b6020808252825182820181905260009190848201906040850190845b81811015611ead5783516001600160a01b031683529284019291840191600101611e88565b50909695505050505050565b600081518084526020808501808196508360051b8101915082860160005b85811015611f01578284038952611eef848351611c99565b98850198935090840190600101611ed7565b5091979650505050505050565b602081526000611cd86020830184611eb9565b634e487b7160e01b600052601160045260246000fd5b818103818111156108a7576108a7611f21565b60008186825b87811015611f6e578135835260209283019290910190600101611f50565b5050938452505060601b6bffffffffffffffffffffffff1916602082015260340192915050565b634e487b7160e01b600052603260045260246000fd5b600181811c90821680611fbf57607f821691505b602082108103611fdf57634e487b7160e01b600052602260045260246000fd5b50919050565b6000808354611ff381611fab565b6001828116801561200b57600181146120205761204f565b60ff198416875282151583028701945061204f565b8760005260208060002060005b858110156120465781548a82015290840190820161202d565b50505082870194505b50929695505050505050565b60006001820161206d5761206d611f21565b5060010190565b60008251612086818460208701611c75565b9190910192915050565b634e487b7160e01b600052602160045260246000fd5b60208082526023908201527f7265737472696374656420746f20746865206175746f6e69747920636f6e74726040820152621858dd60ea1b606082015260800190565b808201808211156108a7576108a7611f21565b808201828112600083128015821682158216171561211c5761211c611f21565b505092915050565b6040815260006121376040830185611eb9565b90508260208301529392505050565b818103600083128015838313168383128216171561216657612166611f21565b5092915050565b634e487b7160e01b600052601260045260246000fd5b6000826121925761219261216d565b600160ff1b8214600019841416156121ac576121ac611f21565b500590565b60006001600160ff1b01820161206d5761206d611f21565b6000600160ff1b82016121de576121de611f21565b506000190190565b6000826121f5576121f561216d565b500490565b6000826122095761220961216d565b500690565b601f82111561149e57600081815260208120601f850160051c810160208610156122355750805b601f850160051c820191505b8181101561164857828155600101612241565b81810361225f575050565b6122698254611fab565b67ffffffffffffffff81111561228157612281611abf565b6122958161228f8454611fab565b8461220e565b6000601f8211600181146122c957600083156122b15750848201545b600019600385901b1c1916600184901b17845561232e565b600085815260209020601f19841690600086815260209020845b8381101561230357828601548255600195860195909101906020016122e3565b50858310156123215781850154600019600388901b60f8161c191681555b50505060018360011b0184555b5050505050565b815167ffffffffffffffff81111561234f5761234f611abf565b61235d8161228f8454611fab565b602080601f831160018114612392576000841561237a5750858301515b600019600386901b1c1916600185901b178555611648565b600085815260208120601f198616915b828110156123c1578886015182559484019460019091019084016123a2565b50858210156123df5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea264697066735822122068031b136be48d163f8475c337216a54400c89da096cc3efb3384bf7d9c52f0d64736f6c63430008150033",
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
	return address, gasConsumed, &Oracle{contract: c}, nil
}

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	*contract
}

// GetPrecision is a free data retrieval call binding the contract method 0x9670c0bc.
//
// Solidity: function getPrecision() pure returns(uint256)
func (_Oracle *Oracle) GetPrecision(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) GetRound(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) GetRoundData(opts *runOptions, _round *big.Int, _symbol string) (IOracleRoundData, uint64, error) {
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
func (_Oracle *Oracle) GetSymbols(opts *runOptions) ([]string, uint64, error) {
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
func (_Oracle *Oracle) GetVotePeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) GetVoters(opts *runOptions) ([]common.Address, uint64, error) {
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
func (_Oracle *Oracle) LastRoundBlock(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) LastVoterUpdateRound(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) LatestRoundData(opts *runOptions, _symbol string) (IOracleRoundData, uint64, error) {
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
func (_Oracle *Oracle) NewSymbols(opts *runOptions, arg0 *big.Int) (string, uint64, error) {
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
func (_Oracle *Oracle) Reports(opts *runOptions, arg0 string, arg1 common.Address) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) Round(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) SymbolUpdatedRound(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) Symbols(opts *runOptions, arg0 *big.Int) (string, uint64, error) {
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
func (_Oracle *Oracle) VotePeriod(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Oracle *Oracle) VotingInfo(opts *runOptions, arg0 common.Address) (struct {
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
func (_Oracle *Oracle) Finalize(opts *runOptions) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "finalize")
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _operator) returns()
func (_Oracle *Oracle) SetOperator(opts *runOptions, _operator common.Address) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "setOperator", _operator)
	return consumed, err
}

// SetSymbols is a paid mutator transaction binding the contract method 0x8d4f75d2.
//
// Solidity: function setSymbols(string[] _symbols) returns()
func (_Oracle *Oracle) SetSymbols(opts *runOptions, _symbols []string) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "setSymbols", _symbols)
	return consumed, err
}

// SetVoters is a paid mutator transaction binding the contract method 0x845023f2.
//
// Solidity: function setVoters(address[] _newVoters) returns()
func (_Oracle *Oracle) SetVoters(opts *runOptions, _newVoters []common.Address) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "setVoters", _newVoters)
	return consumed, err
}

// Vote is a paid mutator transaction binding the contract method 0x307de9b6.
//
// Solidity: function vote(uint256 _commit, int256[] _reports, uint256 _salt) returns()
func (_Oracle *Oracle) Vote(opts *runOptions, _commit *big.Int, _reports []*big.Int, _salt *big.Int) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "vote", _commit, _reports, _salt)
	return consumed, err
}

// Fallback is a paid mutator transaction binding the contract fallback function.
// WARNING! UNTESTED
// Solidity: fallback() payable returns()
func (_Oracle *Oracle) Fallback(opts *runOptions, calldata []byte) (uint64, error) {
	_, consumed, err := _Oracle.call(opts, "", calldata)
	return consumed, err
}

// Receive is a paid mutator transaction binding the contract receive function.
// WARNING! UNTESTED
// Solidity: receive() payable returns()
func (_Oracle *Oracle) Receive(opts *runOptions) (uint64, error) {
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
 		func (_Oracle *Oracle) FilterNewRound(opts *bind.FilterOpts) (*OracleNewRoundIterator, error) {






			logs, sub, err := _Oracle.contract.FilterLogs(opts, "NewRound")
			if err != nil {
				return nil, err
			}
			return &OracleNewRoundIterator{contract: _Oracle.contract, event: "NewRound", logs: logs, sub: sub}, nil
 		}

		// WatchNewRound is a free log subscription operation binding the contract event 0xb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e5.
		//
		// Solidity: event NewRound(uint256 _round, uint256 _height, uint256 _timestamp, uint256 _votePeriod)
		func (_Oracle *Oracle) WatchNewRound(opts *bind.WatchOpts, sink chan<- *OracleNewRound) (event.Subscription, error) {






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
		func (_Oracle *Oracle) ParseNewRound(log types.Log) (*OracleNewRound, error) {
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
 		func (_Oracle *Oracle) FilterNewSymbols(opts *bind.FilterOpts) (*OracleNewSymbolsIterator, error) {




			logs, sub, err := _Oracle.contract.FilterLogs(opts, "NewSymbols")
			if err != nil {
				return nil, err
			}
			return &OracleNewSymbolsIterator{contract: _Oracle.contract, event: "NewSymbols", logs: logs, sub: sub}, nil
 		}

		// WatchNewSymbols is a free log subscription operation binding the contract event 0xaa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d.
		//
		// Solidity: event NewSymbols(string[] _symbols, uint256 _round)
		func (_Oracle *Oracle) WatchNewSymbols(opts *bind.WatchOpts, sink chan<- *OracleNewSymbols) (event.Subscription, error) {




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
		func (_Oracle *Oracle) ParseNewSymbols(log types.Log) (*OracleNewSymbols, error) {
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
 		func (_Oracle *Oracle) FilterVoted(opts *bind.FilterOpts, _voter []common.Address) (*OracleVotedIterator, error) {

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
		func (_Oracle *Oracle) WatchVoted(opts *bind.WatchOpts, sink chan<- *OracleVoted, _voter []common.Address) (event.Subscription, error) {

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
		func (_Oracle *Oracle) ParseVoted(log types.Log) (*OracleVoted, error) {
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
	Bin: "0x61012561003a600b82828239805160001a60731461002d57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe730000000000000000000000000000000000000000301460806040526004361060885760003560e01c8063925c549211605f578063925c54921460c6578063a4ad5d911460cd578063c13974e11460d4578063d0a6d1a61460db57600080fd5b80632090a44214608d5780634dc925d31460b157806350d937201460b85780638e153dc31460bf575b600080fd5b609460fa81565b6040516001600160a01b0390911681526020015b60405180910390f35b609460fc81565b609460fb81565b609460fd81565b609460fe81565b609460f981565b609460ff81565b60e2600181565b60405190815260200160a856fea2646970667358221220a470277bc4c1e7b278582ab14c86b1a0c084269fa795894c29ef5f41b896f86e64736f6c63430008150033",
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
	return address, gasConsumed, &Precompiled{contract: c}, nil
}

// Precompiled is an auto generated Go binding around an Ethereum contract.
type Precompiled struct {
	*contract
}

// ACCUSATIONCONTRACT is a free data retrieval call binding the contract method 0x4dc925d3.
//
// Solidity: function ACCUSATION_CONTRACT() view returns(address)
func (_Precompiled *Precompiled) ACCUSATIONCONTRACT(opts *runOptions) (common.Address, uint64, error) {
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
func (_Precompiled *Precompiled) COMPUTECOMMITTEECONTRACT(opts *runOptions) (common.Address, uint64, error) {
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
func (_Precompiled *Precompiled) ENODEVERIFIERCONTRACT(opts *runOptions) (common.Address, uint64, error) {
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
func (_Precompiled *Precompiled) INNOCENCECONTRACT(opts *runOptions) (common.Address, uint64, error) {
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
func (_Precompiled *Precompiled) MISBEHAVIOURCONTRACT(opts *runOptions) (common.Address, uint64, error) {
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
func (_Precompiled *Precompiled) POPVERIFIERCONTRACT(opts *runOptions) (common.Address, uint64, error) {
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
func (_Precompiled *Precompiled) SUCCESS(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Precompiled *Precompiled) UPGRADERCONTRACT(opts *runOptions) (common.Address, uint64, error) {
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
	Bin: "0x60806040523480156200001157600080fd5b50604051620027013803806200270183398101604081905262000034916200012e565b8560400151806000036200005b57604051630309cb8760e51b815260040160405180910390fd5b866020015187604001518082106200008657604051630309cb8760e51b815260040160405180910390fd5b5050865160005550602086015160015560408601516002556060860151600355608090950151600455600780546001600160a01b03199081166001600160a01b039687161790915560088054821694861694909417909355600a8054841692851692909217909155600b8054831691841691909117905560098054909116919092161790556200021c565b80516001600160a01b03811681146200012957600080fd5b919050565b6000806000806000808688036101408112156200014a57600080fd5b60a08112156200015957600080fd5b5060405160a081016001600160401b03811182821017156200018b57634e487b7160e01b600052604160045260246000fd5b8060405250875181526020880151602082015260408801516040820152606088015160608201526080880151608082015280965050620001ce60a0880162000111565b9450620001de60c0880162000111565b9350620001ee60e0880162000111565b9250620001ff610100880162000111565b915062000210610120880162000111565b90509295509295509295565b6124d5806200022c6000396000f3fe6080604052600436106101665760003560e01c806368cd03f6116100d1578063946ce8cd1161008a578063c5ebeaec11610064578063c5ebeaec1461046b578063ce4b5bbe1461048b578063eced5526146104a0578063fbbe6991146104b557600080fd5b8063946ce8cd1461040b578063b3ab15fb1461042b578063b6b55f251461044b57600080fd5b806368cd03f6146102d857806379502c55146102fa5780637adbf973146103495780637b44646a1461036957806383baa17414610389578063840c7e24146103a957600080fd5b806350bf06bf1161012357806350bf06bf1461022b57806352e5a0501461024b57806353afe81d1461026b57806354a9f42c1461028b5780635891de72146102ab5780635dcc9391146102c057600080fd5b8063042e02cf1461016b57806308796696146101a057806315184245146101ce5780632e1a7d4d146101ee5780632f86556814610210578063402d888314610223575b600080fd5b34801561017757600080fd5b5061018b61018636600461211e565b6104d5565b60405190151581526020015b60405180910390f35b3480156101ac57600080fd5b506101c06101bb366004612139565b61051d565b604051908152602001610197565b3480156101da57600080fd5b506101c06101e9366004612165565b610587565b3480156101fa57600080fd5b5061020e610209366004612197565b610617565b005b61020e61021e36600461211e565b6107e5565b61020e610a27565b34801561023757600080fd5b506101c06102463660046121b0565b610bff565b34801561025757600080fd5b5061020e61026636600461211e565b610c6b565b34801561027757600080fd5b5061020e610286366004612197565b610cb7565b34801561029757600080fd5b506101c06102a636600461211e565b610ce6565b3480156102b757600080fd5b506101c0610d5d565b3480156102cc57600080fd5b506101c06301e1338081565b3480156102e457600080fd5b506102ed610fff565b60405161019791906121da565b34801561030657600080fd5b50600054600154600254600354600454610321949392919085565b604080519586526020860194909452928401919091526060830152608082015260a001610197565b34801561035557600080fd5b5061020e61036436600461211e565b611061565b34801561037557600080fd5b5061020e610384366004612197565b6110ad565b34801561039557600080fd5b506101c06103a4366004612165565b611125565b3480156103b557600080fd5b506103eb6103c436600461211e565b60056020526000908152604090208054600182015460028301546003909301549192909184565b604080519485526020850193909352918301526060820152608001610197565b34801561041757600080fd5b5061020e610426366004612197565b61117b565b34801561043757600080fd5b5061020e61044636600461211e565b6111d1565b34801561045757600080fd5b5061020e610466366004612197565b61121d565b34801561047757600080fd5b5061020e610486366004612197565b61141c565b34801561049757600080fd5b506101c06115e8565b3480156104ac57600080fd5b506101c0601281565b3480156104c157600080fd5b5061018b6104d0366004612165565b6115f7565b6001600160a01b0381166000908152600560205260408120816104f8824261164e565b509050610515826001015461050b610d5d565b60015484906115f7565b949350505050565b600082806000036105405760405162bfc92160e01b815260040160405180910390fd5b83158061054b575082155b1561056957604051630309cb8760e51b815260040160405180910390fd5b83610574848761223d565b61057e919061226a565b95945050505050565b6000818311156105aa57604051630309cb8760e51b815260040160405180910390fd5b848460006105cc6301e133806105c66105c3898961228c565b90565b906116c4565b905060006105e26105dd84846116e3565b6116f2565b905060006106096106026105fb6105c36012600a612383565b8490611745565b86906116e3565b9a9950505050505050505050565b80806000036106395760405163162908e360e11b815260040160405180910390fd5b336000908152600560205260409020600181015483111561066d5760405163162908e360e11b815260040160405180910390fd5b6000610679824261164e565b5090506000610686610d5d565b905061069d836001015482846000600101546115f7565b156106bb57604051636229415360e01b815260040160405180910390fd5b6106cf83600201548260006002015461051d565b8584600101546106df919061228c565b10156106fe57604051633a23d82560e01b815260040160405180910390fd5b84836001016000828254610712919061228c565b909155505060095460405163a9059cbb60e01b8152336004820152602481018790526001600160a01b039091169063a9059cbb906044016020604051808303816000875af1158015610768573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061078c919061238f565b6107a9576040516312171d8360e31b815260040160405180910390fd5b60405185815233907f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649060200160405180910390a25050505050565b3460000361080657604051637c946ed760e01b815260040160405180910390fd5b6001600160a01b0381166000908152600560205260408120600281015490910361084357604051638aa5baf360e01b815260040160405180910390fd5b600080610850834261164e565b9150915061086e8360010154610864610d5d565b60015485906115f7565b61088b57604051636ef5bcdd60e11b815260040160405180910390fd5b6000610897833461228c565b905060018401805442865560009182905560028601829055600386019190915560095460405163a9059cbb60e01b8152336004820152602481018390526001600160a01b039091169063a9059cbb906044016020604051808303816000875af1158015610908573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061092c919061238f565b610949576040516312171d8360e31b815260040160405180910390fd5b600b546001600160a01b03166344df8e70610964858761228c565b6040518263ffffffff1660e01b81526004016000604051808303818588803b15801561098f57600080fd5b505af11580156109a3573d6000803e3d6000fd5b505050505060008211156109e057604051339083156108fc029084906000818181858888f193505050501580156109de573d6000803e3d6000fd5b505b6040513381526001600160a01b038716907fc3d81b2125598b9a2b024afe09e33981f0aa5b7bcbe3e30c4303a4dec209ddb4906020015b60405180910390a2505050505050565b34600003610a4857604051637c946ed760e01b815260040160405180910390fd5b3360009081526005602052604081206002810154909103610a7c57604051638aa5baf360e01b815260040160405180910390fd5b600080610a89834261164e565b915091508134108015610aa65750600354610aa4348461228c565b105b15610ac45760405163e6bd447960e01b815260040160405180910390fd5b80836003016000828254610ad891906123b1565b909155505042835560008080610aee8634611754565b92509250925081866002016000828254610b08919061228c565b9250508190555082866003016000828254610b23919061228c565b90915550508115610b9857600b60009054906101000a90046001600160a01b03166001600160a01b03166344df8e70836040518263ffffffff1660e01b81526004016000604051808303818588803b158015610b7e57600080fd5b505af1158015610b92573d6000803e3d6000fd5b50505050505b8015610bcd57604051339082156108fc029083906000818181858888f19350505050158015610bcb573d6000803e3d6000fd5b505b60405134815233907f5c16de4f8b59bd9caf0f49a545f25819a895ed223294290b408242e72a59423190602001610a17565b6001600160a01b0382166000908152600560205260408120805484918491821015610c3d57604051630309cb8760e51b815260040160405180910390fd5b6001600160a01b0386166000908152600560205260409020610c5f818761164e565b50979650505050505050565b6008546001600160a01b03163314610c95576040516282b42960e81b815260040160405180910390fd5b600b80546001600160a01b0319166001600160a01b0392909216919091179055565b6008546001600160a01b03163314610ce1576040516282b42960e81b815260040160405180910390fd5b600355565b6040516350bf06bf60e01b81526001600160a01b038216600482015242602482015260009030906350bf06bf90604401602060405180830381865afa158015610d33573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d5791906123c4565b92915050565b600a546040805180820182526007815266272a2716a0aa2760c91b602082015290516333f98c7760e01b815260009283926001600160a01b03909116916333f98c7791610dac916004016123dd565b608060405180830381865afa158015610dc9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ded919061242b565b90508060600151600014610e145760405163cb08be8160e01b815260040160405180910390fd5b6000816020015113610e385760405162bfc92160e01b815260040160405180910390fd5b600a60009054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e8b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610eaf91906123c4565b610ebb6012600a612383565b1115610f6357600a60009054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f14573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f3891906123c4565b610f446012600a612383565b610f4e919061226a565b8160200151610f5d919061223d565b91505090565b610f6f6012600a612383565b600a60009054906101000a90046001600160a01b03166001600160a01b0316639670c0bc6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610fc2573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fe691906123c4565b610ff0919061226a565b8160200151610f5d919061226a565b6060600680548060200260200160405190810160405280929190818152602001828054801561105757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611039575b5050505050905090565b6007546001600160a01b0316331461108b576040516282b42960e81b815260040160405180910390fd5b600a80546001600160a01b0319166001600160a01b0392909216919091179055565b80806000036110cf57604051630309cb8760e51b815260040160405180910390fd5b600154828082106110f357604051630309cb8760e51b815260040160405180910390fd5b6008546001600160a01b0316331461111d576040516282b42960e81b815260040160405180910390fd5b505050600255565b6000831580611132575081155b1561115057604051630309cb8760e51b815260040160405180910390fd5b61115c6012600a612383565b611166908361223d565b83611171868861223d565b610574919061223d565b60025481908082106111a057604051630309cb8760e51b815260040160405180910390fd5b6008546001600160a01b031633146111ca576040516282b42960e81b815260040160405180910390fd5b5050600155565b6007546001600160a01b031633146111fb576040516282b42960e81b815260040160405180910390fd5b600880546001600160a01b0319166001600160a01b0392909216919091179055565b808060000361123f5760405163162908e360e11b815260040160405180910390fd5b600954604051636eb1769f60e11b815233600482015230602482015283916001600160a01b03169063dd62ed3e90604401602060405180830381865afa15801561128d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112b191906123c4565b10156112d0576040516313be252b60e01b815260040160405180910390fd5b336000908152600560205260408120805490910361132b57600680546001810182556000919091527ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f0180546001600160a01b031916331790555b4281556001810180548491906000906113459084906123b1565b90915550506009546040516323b872dd60e01b8152336004820152306024820152604481018590526001600160a01b03909116906323b872dd906064016020604051808303816000875af11580156113a1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113c5919061238f565b6113e2576040516312171d8360e31b815260040160405180910390fd5b60405183815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2505050565b808060000361143e5760405163162908e360e11b815260040160405180910390fd5b3360009081526005602052604081209080611459834261164e565b909250905061146885836123b1565b60035490925082101561148e5760405163e6bd447960e01b815260040160405180910390fd5b6000611498610d5d565b90506114af846001015482856000600101546115f7565b156114cd57604051636229415360e01b815260040160405180910390fd5b60006114e9856001015483600060040154600060020154611125565b90508084111561150c57604051633a23d82560e01b815260040160405180910390fd5b4285556002850180548891906000906115269084906123b1565b925050819055508285600301600082825461154191906123b1565b9091555050600b546040516340c10f1960e01b8152336004820152602481018990526001600160a01b03909116906340c10f1990604401600060405180830381600087803b15801561159257600080fd5b505af11580156115a6573d6000803e3d6000fd5b50506040518981523392507fcbc04eca7e9da35cb1393a6135a199ca52e450d5e9251cbd99f7847d33a36750915060200160405180910390a250505050505050565b6115f46012600a612383565b81565b6000838060000361161a5760405162bfc92160e01b815260040160405180910390fd5b8360000361162b5760009150611645565b8284611637878961223d565b611641919061226a565b1091505b50949350505050565b6000808260000361167257604051630309cb8760e51b815260040160405180910390fd5b60008460030154856002015461168891906123b1565b8554909150840361169c57600091506116b0565b60005485546116ad91839187610587565b91505b6116ba82826123b1565b9250509250929050565b60006116dc6105c384670de0b6b3a7640000856117c8565b9392505050565b60006116dc6105c3848461189b565b600081680736ea4425c11ac63081111561172757604051630d7b1d6560e11b8152600481018490526024015b60405180910390fd5b6714057b7ef767814f8102610515670de0b6b3a76400008204611951565b60006116dc6105c3838561228c565b6000806000808560030154866002015461176e91906123b1565b905085600301548510611785578560030154611787565b845b935080851061179a5785600201546117a4565b6117a4848661228c565b92508085116117b45760006117be565b6117be818661228c565b9150509250925092565b6000808060001985870985870292508281108382030391505080600003611802578382816117f8576117f8612254565b04925050506116dc565b83811061183357604051630c740aef60e31b815260048101879052602481018690526044810185905260640161171e565b600084868809600260036001881981018916988990049182028318808302840302808302840302808302840302808302840302808302840302918202909203026000889003889004909101858311909403939093029303949094049190911702949350505050565b60008080600019848609848602925082811083820303915050806000036118cf5750670de0b6b3a764000090049050610d57565b670de0b6b3a7640000811061190157604051635173648d60e01b8152600481018690526024810185905260440161171e565b6000670de0b6b3a764000085870962040000818503049310909103600160ee1b02919091177faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106690291505092915050565b600081680a688906bd8affffff8111156119815760405163b3b6ba1f60e01b81526004810184905260240161171e565b6000611999670de0b6b3a7640000604084901b61226a565b90506105156105c382600160bf1b67ff00000000000000821615611aaf576780000000000000008216156119d65768016a09e667f3bcc9090260401c5b6740000000000000008216156119f5576801306fe0a31b7152df0260401c5b672000000000000000821615611a14576801172b83c7d517adce0260401c5b671000000000000000821615611a335768010b5586cf9890f62a0260401c5b670800000000000000821615611a52576801059b0d31585743ae0260401c5b670400000000000000821615611a7157680102c9a3e778060ee70260401c5b670200000000000000821615611a905768010163da9fb33356d80260401c5b670100000000000000821615611aaf57680100b1afa5abcbed610260401c5b66ff000000000000821615611bae576680000000000000821615611adc5768010058c86da1c09ea20260401c5b6640000000000000821615611afa576801002c605e2e8cec500260401c5b6620000000000000821615611b1857680100162f3904051fa10260401c5b6610000000000000821615611b36576801000b175effdc76ba0260401c5b6608000000000000821615611b5457680100058ba01fb9f96d0260401c5b6604000000000000821615611b725768010002c5cc37da94920260401c5b6602000000000000821615611b90576801000162e525ee05470260401c5b6601000000000000821615611bae5768010000b17255775c040260401c5b65ff0000000000821615611ca45765800000000000821615611bd9576801000058b91b5bc9ae0260401c5b65400000000000821615611bf657680100002c5c89d5ec6d0260401c5b65200000000000821615611c135768010000162e43f4f8310260401c5b65100000000000821615611c3057680100000b1721bcfc9a0260401c5b65080000000000821615611c4d5768010000058b90cf1e6e0260401c5b65040000000000821615611c6a576801000002c5c863b73f0260401c5b65020000000000821615611c8757680100000162e430e5a20260401c5b65010000000000821615611ca4576801000000b1721835510260401c5b64ff00000000821615611d9157648000000000821615611ccd57680100000058b90c0b490260401c5b644000000000821615611ce95768010000002c5c8601cc0260401c5b642000000000821615611d05576801000000162e42fff00260401c5b641000000000821615611d215768010000000b17217fbb0260401c5b640800000000821615611d3d576801000000058b90bfce0260401c5b640400000000821615611d5957680100000002c5c85fe30260401c5b640200000000821615611d755768010000000162e42ff10260401c5b640100000000821615611d9157680100000000b17217f80260401c5b63ff000000821615611e75576380000000821615611db85768010000000058b90bfc0260401c5b6340000000821615611dd3576801000000002c5c85fe0260401c5b6320000000821615611dee57680100000000162e42ff0260401c5b6310000000821615611e09576801000000000b17217f0260401c5b6308000000821615611e2457680100000000058b90c00260401c5b6304000000821615611e3f5768010000000002c5c8600260401c5b6302000000821615611e5a576801000000000162e4300260401c5b6301000000821615611e755768010000000000b172180260401c5b62ff0000821615611f505762800000821615611e9a576801000000000058b90c0260401c5b62400000821615611eb457680100000000002c5c860260401c5b62200000821615611ece5768010000000000162e430260401c5b62100000821615611ee857680100000000000b17210260401c5b62080000821615611f025768010000000000058b910260401c5b62040000821615611f1c576801000000000002c5c80260401c5b62020000821615611f3657680100000000000162e40260401c5b62010000821615611f50576801000000000000b1720260401c5b61ff0082161561202257618000821615611f7357680100000000000058b90260401c5b614000821615611f8c5768010000000000002c5d0260401c5b612000821615611fa5576801000000000000162e0260401c5b611000821615611fbe5768010000000000000b170260401c5b610800821615611fd7576801000000000000058c0260401c5b610400821615611ff057680100000000000002c60260401c5b61020082161561200957680100000000000001630260401c5b61010082161561202257680100000000000000b10260401c5b60ff8216156120eb57608082161561204357680100000000000000590260401c5b604082161561205b576801000000000000002c0260401c5b602082161561207357680100000000000000160260401c5b601082161561208b576801000000000000000b0260401c5b60088216156120a357680100000000000000060260401c5b60048216156120bb57680100000000000000030260401c5b60028216156120d357680100000000000000010260401c5b60018216156120eb57680100000000000000010260401c5b670de0b6b3a76400000260409190911c60bf031c90565b80356001600160a01b038116811461211957600080fd5b919050565b60006020828403121561213057600080fd5b6116dc82612102565b60008060006060848603121561214e57600080fd5b505081359360208301359350604090920135919050565b6000806000806080858703121561217b57600080fd5b5050823594602084013594506040840135936060013592509050565b6000602082840312156121a957600080fd5b5035919050565b600080604083850312156121c357600080fd5b6121cc83612102565b946020939093013593505050565b6020808252825182820181905260009190848201906040850190845b8181101561221b5783516001600160a01b0316835292840192918401916001016121f6565b50909695505050505050565b634e487b7160e01b600052601160045260246000fd5b8082028115828204841417610d5757610d57612227565b634e487b7160e01b600052601260045260246000fd5b60008261228757634e487b7160e01b600052601260045260246000fd5b500490565b81810381811115610d5757610d57612227565b600181815b808511156122da5781600019048211156122c0576122c0612227565b808516156122cd57918102915b93841c93908002906122a4565b509250929050565b6000826122f157506001610d57565b816122fe57506000610d57565b8160018114612314576002811461231e5761233a565b6001915050610d57565b60ff84111561232f5761232f612227565b50506001821b610d57565b5060208310610133831016604e8410600b841016171561235d575081810a610d57565b612367838361229f565b806000190482111561237b5761237b612227565b029392505050565b60006116dc83836122e2565b6000602082840312156123a157600080fd5b815180151581146116dc57600080fd5b80820180821115610d5757610d57612227565b6000602082840312156123d657600080fd5b5051919050565b600060208083528351808285015260005b8181101561240a578581018301518582016040015282016123ee565b506000604082860101526040601f19601f8301168501019250505092915050565b60006080828403121561243d57600080fd5b6040516080810181811067ffffffffffffffff8211171561246e57634e487b7160e01b600052604160045260246000fd5b806040525082518152602083015160208201526040830151604082015260608301516060820152809150509291505056fea26469706673582212205e76940183cf0f9ee63393d5c96d3129f42207d506dafc89e2104ea4228d5bda64736f6c63430008150033",
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
	return address, gasConsumed, &Stabilization{contract: c}, nil
}

// Stabilization is an auto generated Go binding around an Ethereum contract.
type Stabilization struct {
	*contract
}

// SCALE is a free data retrieval call binding the contract method 0xeced5526.
//
// Solidity: function SCALE() view returns(uint256)
func (_Stabilization *Stabilization) SCALE(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) SCALEFACTOR(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) SECONDSINYEAR(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) Accounts(opts *runOptions) ([]common.Address, uint64, error) {
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
func (_Stabilization *Stabilization) BorrowLimit(opts *runOptions, collateral *big.Int, price *big.Int, targetPrice *big.Int, mcr *big.Int) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) Cdps(opts *runOptions, arg0 common.Address) (struct {
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
func (_Stabilization *Stabilization) CollateralPrice(opts *runOptions) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) Config(opts *runOptions) (struct {
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
func (_Stabilization *Stabilization) DebtAmount(opts *runOptions, account common.Address, timestamp *big.Int) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) DebtAmount0(opts *runOptions, account common.Address) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) InterestDue(opts *runOptions, debt *big.Int, rate *big.Int, timeBorrow *big.Int, timeDue *big.Int) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) IsLiquidatable(opts *runOptions, account common.Address) (bool, uint64, error) {
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
func (_Stabilization *Stabilization) MinimumCollateral(opts *runOptions, principal *big.Int, price *big.Int, mcr *big.Int) (*big.Int, uint64, error) {
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
func (_Stabilization *Stabilization) UnderCollateralized(opts *runOptions, collateral *big.Int, price *big.Int, debt *big.Int, liquidationRatio *big.Int) (bool, uint64, error) {
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
func (_Stabilization *Stabilization) Borrow(opts *runOptions, amount *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "borrow", amount)
	return consumed, err
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Stabilization *Stabilization) Deposit(opts *runOptions, amount *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "deposit", amount)
	return consumed, err
}

// Liquidate is a paid mutator transaction binding the contract method 0x2f865568.
//
// Solidity: function liquidate(address account) payable returns()
func (_Stabilization *Stabilization) Liquidate(opts *runOptions, account common.Address) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "liquidate", account)
	return consumed, err
}

// Repay is a paid mutator transaction binding the contract method 0x402d8883.
//
// Solidity: function repay() payable returns()
func (_Stabilization *Stabilization) Repay(opts *runOptions) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "repay")
	return consumed, err
}

// SetLiquidationRatio is a paid mutator transaction binding the contract method 0x946ce8cd.
//
// Solidity: function setLiquidationRatio(uint256 ratio) returns()
func (_Stabilization *Stabilization) SetLiquidationRatio(opts *runOptions, ratio *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setLiquidationRatio", ratio)
	return consumed, err
}

// SetMinCollateralizationRatio is a paid mutator transaction binding the contract method 0x7b44646a.
//
// Solidity: function setMinCollateralizationRatio(uint256 ratio) returns()
func (_Stabilization *Stabilization) SetMinCollateralizationRatio(opts *runOptions, ratio *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setMinCollateralizationRatio", ratio)
	return consumed, err
}

// SetMinDebtRequirement is a paid mutator transaction binding the contract method 0x53afe81d.
//
// Solidity: function setMinDebtRequirement(uint256 amount) returns()
func (_Stabilization *Stabilization) SetMinDebtRequirement(opts *runOptions, amount *big.Int) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setMinDebtRequirement", amount)
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_Stabilization *Stabilization) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setOperator", operator)
	return consumed, err
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address oracle) returns()
func (_Stabilization *Stabilization) SetOracle(opts *runOptions, oracle common.Address) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setOracle", oracle)
	return consumed, err
}

// SetSupplyControl is a paid mutator transaction binding the contract method 0x52e5a050.
//
// Solidity: function setSupplyControl(address supplyControl) returns()
func (_Stabilization *Stabilization) SetSupplyControl(opts *runOptions, supplyControl common.Address) (uint64, error) {
	_, consumed, err := _Stabilization.call(opts, "setSupplyControl", supplyControl)
	return consumed, err
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_Stabilization *Stabilization) Withdraw(opts *runOptions, amount *big.Int) (uint64, error) {
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
 		func (_Stabilization *Stabilization) FilterBorrow(opts *bind.FilterOpts, account []common.Address) (*StabilizationBorrowIterator, error) {

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
		func (_Stabilization *Stabilization) WatchBorrow(opts *bind.WatchOpts, sink chan<- *StabilizationBorrow, account []common.Address) (event.Subscription, error) {

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
		func (_Stabilization *Stabilization) ParseBorrow(log types.Log) (*StabilizationBorrow, error) {
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
 		func (_Stabilization *Stabilization) FilterDeposit(opts *bind.FilterOpts, account []common.Address) (*StabilizationDepositIterator, error) {

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
		func (_Stabilization *Stabilization) WatchDeposit(opts *bind.WatchOpts, sink chan<- *StabilizationDeposit, account []common.Address) (event.Subscription, error) {

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
		func (_Stabilization *Stabilization) ParseDeposit(log types.Log) (*StabilizationDeposit, error) {
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
 		func (_Stabilization *Stabilization) FilterLiquidate(opts *bind.FilterOpts, account []common.Address) (*StabilizationLiquidateIterator, error) {

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
		func (_Stabilization *Stabilization) WatchLiquidate(opts *bind.WatchOpts, sink chan<- *StabilizationLiquidate, account []common.Address) (event.Subscription, error) {

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
		func (_Stabilization *Stabilization) ParseLiquidate(log types.Log) (*StabilizationLiquidate, error) {
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
 		func (_Stabilization *Stabilization) FilterRepay(opts *bind.FilterOpts, account []common.Address) (*StabilizationRepayIterator, error) {

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
		func (_Stabilization *Stabilization) WatchRepay(opts *bind.WatchOpts, sink chan<- *StabilizationRepay, account []common.Address) (event.Subscription, error) {

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
		func (_Stabilization *Stabilization) ParseRepay(log types.Log) (*StabilizationRepay, error) {
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
 		func (_Stabilization *Stabilization) FilterWithdraw(opts *bind.FilterOpts, account []common.Address) (*StabilizationWithdrawIterator, error) {

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
		func (_Stabilization *Stabilization) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *StabilizationWithdraw, account []common.Address) (event.Subscription, error) {

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
		func (_Stabilization *Stabilization) ParseWithdraw(log types.Log) (*StabilizationWithdraw, error) {
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
	Bin: "0x6080604052604051610512380380610512833981016040819052610022916100a5565b3460000361004357604051637c946ed760e01b815260040160405180910390fd5b600280546001600160a01b039485166001600160a01b031991821617909155600380549385169382169390931790925560008054919093169116179055346001556100e8565b80516001600160a01b03811681146100a057600080fd5b919050565b6000806000606084860312156100ba57600080fd5b6100c384610089565b92506100d160208501610089565b91506100df60408501610089565b90509250925092565b61041b806100f76000396000f3fe6080604052600436106100705760003560e01c80637e47961c1161004e5780637e47961c146100c85780637ecc2b5614610100578063b3ab15fb14610113578063db7f521a1461013357600080fd5b806318160ddd1461007557806340c10f191461009e57806344df8e70146100c0575b600080fd5b34801561008157600080fd5b5061008b60015481565b6040519081526020015b60405180910390f35b3480156100aa57600080fd5b506100be6100b9366004610399565b610153565b005b6100be610265565b3480156100d457600080fd5b506000546100e8906001600160a01b031681565b6040516001600160a01b039091168152602001610095565b34801561010c57600080fd5b504761008b565b34801561011f57600080fd5b506100be61012e3660046103c3565b6102e5565b34801561013f57600080fd5b506100be61014e3660046103c3565b610331565b6000546001600160a01b0316331461017d576040516282b42960e81b815260040160405180910390fd5b6001600160a01b03821615806101a057506000546001600160a01b038381169116145b156101be57604051634e46966960e11b815260040160405180910390fd5b8015806101ca57504781115b156101e85760405163162908e360e11b815260040160405180910390fd5b6040516001600160a01b0383169082156108fc029083906000818181858888f1935050505015801561021e573d6000803e3d6000fd5b50604080516001600160a01b0384168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a15050565b3460000361028657604051637c946ed760e01b815260040160405180910390fd5b6000546001600160a01b031633146102b0576040516282b42960e81b815260040160405180910390fd5b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a1565b6002546001600160a01b0316331461030f576040516282b42960e81b815260040160405180910390fd5b600380546001600160a01b0319166001600160a01b0392909216919091179055565b6003546001600160a01b0316331461035b576040516282b42960e81b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b0392909216919091179055565b80356001600160a01b038116811461039457600080fd5b919050565b600080604083850312156103ac57600080fd5b6103b58361037d565b946020939093013593505050565b6000602082840312156103d557600080fd5b6103de8261037d565b939250505056fea26469706673582212207646b2e6881a2d6269951386532b4fc6a3409fe50e89f6a55b709c4060cd997f64736f6c63430008150033",
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
	return address, gasConsumed, &SupplyControl{contract: c}, nil
}

// SupplyControl is an auto generated Go binding around an Ethereum contract.
type SupplyControl struct {
	*contract
}

// AvailableSupply is a free data retrieval call binding the contract method 0x7ecc2b56.
//
// Solidity: function availableSupply() view returns(uint256)
func (_SupplyControl *SupplyControl) AvailableSupply(opts *runOptions) (*big.Int, uint64, error) {
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
func (_SupplyControl *SupplyControl) Stabilizer(opts *runOptions) (common.Address, uint64, error) {
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
func (_SupplyControl *SupplyControl) TotalSupply(opts *runOptions) (*big.Int, uint64, error) {
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
func (_SupplyControl *SupplyControl) Burn(opts *runOptions) (uint64, error) {
	_, consumed, err := _SupplyControl.call(opts, "burn")
	return consumed, err
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address recipient, uint256 amount) returns()
func (_SupplyControl *SupplyControl) Mint(opts *runOptions, recipient common.Address, amount *big.Int) (uint64, error) {
	_, consumed, err := _SupplyControl.call(opts, "mint", recipient, amount)
	return consumed, err
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address operator) returns()
func (_SupplyControl *SupplyControl) SetOperator(opts *runOptions, operator common.Address) (uint64, error) {
	_, consumed, err := _SupplyControl.call(opts, "setOperator", operator)
	return consumed, err
}

// SetStabilizer is a paid mutator transaction binding the contract method 0xdb7f521a.
//
// Solidity: function setStabilizer(address stabilizer_) returns()
func (_SupplyControl *SupplyControl) SetStabilizer(opts *runOptions, stabilizer_ common.Address) (uint64, error) {
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
 		func (_SupplyControl *SupplyControl) FilterBurn(opts *bind.FilterOpts) (*SupplyControlBurnIterator, error) {



			logs, sub, err := _SupplyControl.contract.FilterLogs(opts, "Burn")
			if err != nil {
				return nil, err
			}
			return &SupplyControlBurnIterator{contract: _SupplyControl.contract, event: "Burn", logs: logs, sub: sub}, nil
 		}

		// WatchBurn is a free log subscription operation binding the contract event 0xb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb.
		//
		// Solidity: event Burn(uint256 amount)
		func (_SupplyControl *SupplyControl) WatchBurn(opts *bind.WatchOpts, sink chan<- *SupplyControlBurn) (event.Subscription, error) {



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
		func (_SupplyControl *SupplyControl) ParseBurn(log types.Log) (*SupplyControlBurn, error) {
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
 		func (_SupplyControl *SupplyControl) FilterMint(opts *bind.FilterOpts) (*SupplyControlMintIterator, error) {




			logs, sub, err := _SupplyControl.contract.FilterLogs(opts, "Mint")
			if err != nil {
				return nil, err
			}
			return &SupplyControlMintIterator{contract: _SupplyControl.contract, event: "Mint", logs: logs, sub: sub}, nil
 		}

		// WatchMint is a free log subscription operation binding the contract event 0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885.
		//
		// Solidity: event Mint(address recipient, uint256 amount)
		func (_SupplyControl *SupplyControl) WatchMint(opts *bind.WatchOpts, sink chan<- *SupplyControlMint) (event.Subscription, error) {




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
		func (_SupplyControl *SupplyControl) ParseMint(log types.Log) (*SupplyControlMint, error) {
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
	Bin: "0x608060405234801561001057600080fd5b5060405161041238038061041283398101604081905261002f91610058565b600061003b82826101aa565b5050610269565b634e487b7160e01b600052604160045260246000fd5b6000602080838503121561006b57600080fd5b82516001600160401b038082111561008257600080fd5b818501915085601f83011261009657600080fd5b8151818111156100a8576100a8610042565b604051601f8201601f19908116603f011681019083821181831017156100d0576100d0610042565b8160405282815288868487010111156100e857600080fd5b600093505b8284101561010a57848401860151818501870152928501926100ed565b600086848301015280965050505050505092915050565b600181811c9082168061013557607f821691505b60208210810361015557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156101a557600081815260208120601f850160051c810160208610156101825750805b601f850160051c820191505b818110156101a15782815560010161018e565b5050505b505050565b81516001600160401b038111156101c3576101c3610042565b6101d7816101d18454610121565b8461015b565b602080601f83116001811461020c57600084156101f45750858301515b600019600386901b1c1916600185901b1785556101a1565b600085815260208120601f198616915b8281101561023b5788860151825594840194600190910190840161021c565b50858210156102595787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61019a806102786000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063bfb4ebcf14610030575b600080fd5b61003861004e565b60405161004591906100dc565b60405180910390f35b6000805461005b9061012a565b80601f01602080910402602001604051908101604052809291908181526020018280546100879061012a565b80156100d45780601f106100a9576101008083540402835291602001916100d4565b820191906000526020600020905b8154815290600101906020018083116100b757829003601f168201915b505050505081565b600060208083528351808285015260005b81811015610109578581018301518582016040015282016100ed565b506000604082860101526040601f19601f8301168501019250505092915050565b600181811c9082168061013e57607f821691505b60208210810361015e57634e487b7160e01b600052602260045260246000fd5b5091905056fea2646970667358221220721b96e4eacfa27ee5c77c0364c12e322da6414667595c874a886675378f899a64736f6c63430008150033",
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
	return address, gasConsumed, &TestBase{contract: c}, nil
}

// TestBase is an auto generated Go binding around an Ethereum contract.
type TestBase struct {
	*contract
}

// Foo is a free data retrieval call binding the contract method 0xbfb4ebcf.
//
// Solidity: function Foo() view returns(string)
func (_TestBase *TestBase) Foo(opts *runOptions) (string, uint64, error) {
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
	Bin: "0x608060405234801561001057600080fd5b5060405161068838038061068883398101604081905261002f9161010e565b80600061003c82826101fa565b506001905061004b83826101fa565b5050506102b9565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261007a57600080fd5b81516001600160401b038082111561009457610094610053565b604051601f8301601f19908116603f011681019082821181831017156100bc576100bc610053565b816040528381526020925086838588010111156100d857600080fd5b600091505b838210156100fa57858201830151818301840152908201906100dd565b600093810190920192909252949350505050565b6000806040838503121561012157600080fd5b82516001600160401b038082111561013857600080fd5b61014486838701610069565b9350602085015191508082111561015a57600080fd5b5061016785828601610069565b9150509250929050565b600181811c9082168061018557607f821691505b6020821081036101a557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156101f557600081815260208120601f850160051c810160208610156101d25750805b601f850160051c820191505b818110156101f1578281556001016101de565b5050505b505050565b81516001600160401b0381111561021357610213610053565b610227816102218454610171565b846101ab565b602080601f83116001811461025c57600084156102445750858301515b600019600386901b1c1916600185901b1785556101f1565b600085815260208120601f198616915b8281101561028b5788860151825594840194600190910190840161026c565b50858210156102a95787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6103c0806102c86000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80631e4f339514610046578063b0a378b01461005b578063bfb4ebcf14610079575b600080fd5b610059610054366004610142565b610081565b005b610063610091565b60405161007091906101f3565b60405180910390f35b61006361011f565b600061008d82826102ca565b5050565b6001805461009e90610241565b80601f01602080910402602001604051908101604052809291908181526020018280546100ca90610241565b80156101175780601f106100ec57610100808354040283529160200191610117565b820191906000526020600020905b8154815290600101906020018083116100fa57829003601f168201915b505050505081565b6000805461009e90610241565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561015457600080fd5b813567ffffffffffffffff8082111561016c57600080fd5b818401915084601f83011261018057600080fd5b8135818111156101925761019261012c565b604051601f8201601f19908116603f011681019083821181831017156101ba576101ba61012c565b816040528281528760208487010111156101d357600080fd5b826020860160208301376000928101602001929092525095945050505050565b600060208083528351808285015260005b8181101561022057858101830151858201604001528201610204565b506000604082860101526040601f19601f8301168501019250505092915050565b600181811c9082168061025557607f821691505b60208210810361027557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156102c557600081815260208120601f850160051c810160208610156102a25750805b601f850160051c820191505b818110156102c1578281556001016102ae565b5050505b505050565b815167ffffffffffffffff8111156102e4576102e461012c565b6102f8816102f28454610241565b8461027b565b602080601f83116001811461032d57600084156103155750858301515b600019600386901b1c1916600185901b1785556102c1565b600085815260208120601f198616915b8281101561035c5788860151825594840194600190910190840161033d565b508582101561037a5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea2646970667358221220df9dc5ee36a6177eac85a385b417844be4615597f2759cb94fb3ea2fef95019264736f6c63430008150033",
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
	return address, gasConsumed, &TestUpgraded{contract: c}, nil
}

// TestUpgraded is an auto generated Go binding around an Ethereum contract.
type TestUpgraded struct {
	*contract
}

// Bar is a free data retrieval call binding the contract method 0xb0a378b0.
//
// Solidity: function Bar() view returns(string)
func (_TestUpgraded *TestUpgraded) Bar(opts *runOptions) (string, uint64, error) {
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
func (_TestUpgraded *TestUpgraded) Foo(opts *runOptions) (string, uint64, error) {
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
func (_TestUpgraded *TestUpgraded) FooBar(opts *runOptions, _foo string) (uint64, error) {
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
	Bin: "0x608060405234801561001057600080fd5b5060405161044338038061044383398101604081905261002f9161007c565b600080546001600160a01b039384166001600160a01b031991821617909155600180549290931691161790556100af565b80516001600160a01b038116811461007757600080fd5b919050565b6000806040838503121561008f57600080fd5b61009883610060565b91506100a660208401610060565b90509250929050565b610385806100be6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806355463ceb14610051578063570ca735146100805780636e3d9ff014610093578063b3ab15fb146100a8575b600080fd5b600054610064906001600160a01b031681565b6040516001600160a01b03909116815260200160405180910390f35b600154610064906001600160a01b031681565b6100a66100a1366004610220565b6100bb565b005b6100a66100b63660046102e2565b610166565b6001546001600160a01b0316331461011a5760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064015b60405180910390fd5b60405160f9906000906101339085908590602001610304565b6040516020818303038152906040529050600080825160208401855af43d6000803e808015610161573d6000f35b3d6000fd5b6000546001600160a01b031633146101cc5760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f7420746865204175746f6e69747920636f6e74726044820152621858dd60ea1b6064820152608401610111565b600180546001600160a01b0319166001600160a01b0392909216919091179055565b80356001600160a01b038116811461020557600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561023357600080fd5b61023c836101ee565b9150602083013567ffffffffffffffff8082111561025957600080fd5b818501915085601f83011261026d57600080fd5b81358181111561027f5761027f61020a565b604051601f8201601f19908116603f011681019083821181831017156102a7576102a761020a565b816040528281528860208487010111156102c057600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6000602082840312156102f457600080fd5b6102fd826101ee565b9392505050565b6bffffffffffffffffffffffff198360601b1681526000825160005b8181101561033d5760208186018101516014868401015201610320565b5060009201601401918252509291505056fea26469706673582212207cb99b8c663b7b3f7201f7c37a4ee5d8f83506bfd644debc45c05e5883ef0c7064736f6c63430008150033",
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
	return address, gasConsumed, &UpgradeManager{contract: c}, nil
}

// UpgradeManager is an auto generated Go binding around an Ethereum contract.
type UpgradeManager struct {
	*contract
}

// Autonity is a free data retrieval call binding the contract method 0x55463ceb.
//
// Solidity: function autonity() view returns(address)
func (_UpgradeManager *UpgradeManager) Autonity(opts *runOptions) (common.Address, uint64, error) {
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
func (_UpgradeManager *UpgradeManager) Operator(opts *runOptions) (common.Address, uint64, error) {
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
func (_UpgradeManager *UpgradeManager) SetOperator(opts *runOptions, _account common.Address) (uint64, error) {
	_, consumed, err := _UpgradeManager.call(opts, "setOperator", _account)
	return consumed, err
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager *UpgradeManager) Upgrade(opts *runOptions, _target common.Address, _data string) (uint64, error) {
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
	*contract
}

// GetNewContract is a free data retrieval call binding the contract method 0xb66b3e79.
//
// Solidity: function getNewContract() view returns(bytes, string)
func (_Upgradeable *Upgradeable) GetNewContract(opts *runOptions) ([]byte, string, uint64, error) {
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
func (_Upgradeable *Upgradeable) CompleteContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _Upgradeable.call(opts, "completeContractUpgrade")
	return consumed, err
}

// ResetContractUpgrade is a paid mutator transaction binding the contract method 0xcf9c5719.
//
// Solidity: function resetContractUpgrade() returns()
func (_Upgradeable *Upgradeable) ResetContractUpgrade(opts *runOptions) (uint64, error) {
	_, consumed, err := _Upgradeable.call(opts, "resetContractUpgrade")
	return consumed, err
}

// UpgradeContract is a paid mutator transaction binding the contract method 0xb2ea9adb.
//
// Solidity: function upgradeContract(bytes _bytecode, string _abi) returns()
func (_Upgradeable *Upgradeable) UpgradeContract(opts *runOptions, _bytecode []byte, _abi string) (uint64, error) {
	_, consumed, err := _Upgradeable.call(opts, "upgradeContract", _bytecode, _abi)
	return consumed, err
}

/* EVENTS ARE NOT YET SUPPORTED

 */
