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
	Bin: "0x608060405260006011553480156200001657600080fd5b5060405162003cbe38038062003cbe8339810160408190526200003991620000f7565b600180546001600160a01b0319166001600160a01b03841690811790915560408051636fd8d26960e11b8152905163dfb1a4d2916004808201926020929091908290030181865afa15801562000093573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000b99190620001bf565b6000558051600355602081015160045560408101516005556060810151600655608081015160075560a081015160085560c0015160095550620001d9565b6000808284036101008112156200010d57600080fd5b83516001600160a01b03811681146200012557600080fd5b925060e0601f19820112156200013a57600080fd5b5060405160e081016001600160401b03811182821017156200016c57634e487b7160e01b600052604160045260246000fd5b80604052506020840151815260408401516020820152606084015160408201526080840151606082015260a0840151608082015260c084015160a082015260e084015160c0820152809150509250929050565b600060208284031215620001d257600080fd5b5051919050565b613ad580620001e96000396000f3fe6080604052600436106100c25760003560e01c80636c9789b01161007f5780639cb22b06116100595780639cb22b0614610291578063b5b7a184146102be578063bebaa8fc146102e2578063e7bb0b521461030f57600080fd5b80636c9789b0146101d657806379502c55146101f65780637ccecadd1461025a57600080fd5b806301567739146100c7578063069f68631461011a5780630b7914301461013c5780631de9d9b6146101735780634108a95a146101865780636b5f444c146101b6575b600080fd5b3480156100d357600080fd5b506100fd6100e2366004612e99565b600a602052600090815260409020546001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b34801561012657600080fd5b5061013a610135366004612fff565b610347565b005b34801561014857600080fd5b5061015c610157366004613100565b6106af565b6040516101119b9a999897969594939291906131a3565b61013a610181366004612e99565b6107bc565b34801561019257600080fd5b506101a66101a1366004613229565b6109c2565b6040519015158152602001610111565b3480156101c257600080fd5b5061013a6101d1366004613100565b610a72565b3480156101e257600080fd5b5061013a6101f1366004613267565b610aa1565b34801561020257600080fd5b506003546004546005546006546007546008546009546102259695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e001610111565b34801561026657600080fd5b5061027a610275366004613229565b610ae4565b604080519215158352602083019190915201610111565b34801561029d57600080fd5b506102b16102ac366004612e99565b610c39565b6040516101119190613354565b3480156102ca57600080fd5b506102d460005481565b604051908152602001610111565b3480156102ee57600080fd5b506103026102fd366004612e99565b610e56565b6040516101119190613367565b34801561031b57600080fd5b506102d461032a3660046133c9565b600e60209081526000928352604080842090915290825290205481565b600154604051630c825d9760e11b81523360048201526000916001600160a01b031690631904bb2e90602401600060405180830381865afa158015610390573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103b8919081019061345c565b60208101519091506001600160a01b031633146104325760405162461bcd60e51b815260206004820152602d60248201527f66756e6374696f6e207265737472696374656420746f2061207265676973746560448201526c3932b2103b30b634b230ba37b960991b60648201526084015b60405180910390fd5b60808201516001600160a01b0316331461048e5760405162461bcd60e51b815260206004820152601d60248201527f6576656e74207265706f72746572206d7573742062652063616c6c65720000006044820152606401610429565b6001826000015160ff1611156106395760006104a9836110eb565b9050806104b557505050565b336000908152600d6020908152604091829020825161016081018452815460ff8082168352610100820481169483019490945290939192918401916201000090910416600281111561050957610509613119565b600281111561051a5761051a613119565b815281546020909101906301000000900460ff16600a81111561053f5761053f613119565b600a81111561055057610550613119565b815281546001600160a01b03600160201b90910481166020830152600183015416604082015260028201805460609092019161058b906135da565b80601f01602080910402602001604051908101604052809291908181526020018280546105b7906135da565b80156106045780601f106105d957610100808354040283529160200191610604565b820191906000526020600020905b8154815290600101906020018083116105e757829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815250509250505b60008260400151600281111561065157610651613119565b036106635761065f82611319565b5050565b60018260400151600281111561067b5761067b613119565b036106895761065f826114c5565b6002826040015160028111156106a1576106a1613119565b0361065f5761065f82611652565b600281815481106106bf57600080fd5b600091825260209091206007909102018054600182015460028301805460ff8085169650610100850481169562010000860482169563010000008104909216946001600160a01b03600160201b909304831694921692909190610721906135da565b80601f016020809104026020016040519081016040528092919081815260200182805461074d906135da565b801561079a5780601f1061076f5761010080835404028352916020019161079a565b820191906000526020600020905b81548152906001019060200180831161077d57829003601f168201915b505050505090806003015490806004015490806005015490806006015490508b565b6001546001600160a01b031633146107e65760405162461bcd60e51b815260040161042990613614565b6001546001600160a01b038281166000908152600a6020526040808220549051630c825d9760e11b8152908316600482015290929190911690631904bb2e90602401600060405180830381865afa158015610845573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261086d919081019061345c565b516040519091506000906001600160a01b038316906108fc90349084818181858888f193505050503d80600081146108c1576040519150601f19603f3d011682016040523d82523d6000602084013e6108c6565b606091505b505090508061099957600160009054906101000a90046001600160a01b03166001600160a01b031663f7866ee36040518163ffffffff1660e01b8152600401602060405180830381865afa158015610922573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109469190613658565b6001600160a01b03163460405160006040518083038185875af1925050503d8060008114610990576040519150601f19603f3d011682016040523d82523d6000602084013e610995565b606091505b5050505b50506001600160a01b03166000908152600a6020526040902080546001600160a01b0319169055565b6000806109ce84611768565b6001546040516396b477cb60e01b8152600481018690529192506000916001600160a01b03909116906396b477cb90602401602060405180830381865afa158015610a1d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a419190613675565b6001600160a01b0387166000908152600e602090815260408083209383529290522054919091109150509392505050565b6001546001600160a01b03163314610a9c5760405162461bcd60e51b815260040161042990613614565b600055565b6001546001600160a01b03163314610acb5760405162461bcd60e51b815260040161042990613614565b610ad36117d3565b8015610ae157610ae1611b21565b50565b6000806000610af285611768565b6001546040516396b477cb60e01b8152600481018790529192506000916001600160a01b03909116906396b477cb90602401602060405180830381865afa158015610b41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b659190613675565b6001600160a01b0388166000908152600e602090815260408083208484529091529020549091508211610b9f576000935060009250610c2f565b6001600160a01b0387166000908152600c602052604090205415610c26576001600160a01b0387166000908152600c6020526040812054600290610be5906001906136a4565b81548110610bf557610bf56136b7565b90600052602060002090600702019050600094506003600001548160030154610c1e91906136cd565b935050610c2f565b60019350600092505b5050935093915050565b610c41612db3565b6001600160a01b0382166000908152600c6020526040902054610c965760405162461bcd60e51b815260206004820152600d60248201526c37379030b1b1bab9b0ba34b7b760991b6044820152606401610429565b6001600160a01b0382166000908152600c6020526040902054600290610cbe906001906136a4565b81548110610cce57610cce6136b7565b600091825260209182902060408051610160810182526007909302909101805460ff808216855261010082048116958501959095529293909291840191620100009004166002811115610d2357610d23613119565b6002811115610d3457610d34613119565b815281546020909101906301000000900460ff16600a811115610d5957610d59613119565b600a811115610d6a57610d6a613119565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191610da5906135da565b80601f0160208091040260200160405190810160405280929190818152602001828054610dd1906135da565b8015610e1e5780601f10610df357610100808354040283529160200191610e1e565b820191906000526020600020905b815481529060010190602001808311610e0157829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815250509050919050565b6001600160a01b0381166000908152600b60205260408120546060919067ffffffffffffffff811115610e8b57610e8b612ebd565b604051908082528060200260200182016040528015610ec457816020015b610eb1612db3565b815260200190600190039081610ea95790505b50905060005b6001600160a01b0384166000908152600b60205260409020548110156110e4576001600160a01b0384166000908152600b6020526040902080546002919083908110610f1857610f186136b7565b906000526020600020015481548110610f3357610f336136b7565b600091825260209182902060408051610160810182526007909302909101805460ff808216855261010082048116958501959095529293909291840191620100009004166002811115610f8857610f88613119565b6002811115610f9957610f99613119565b815281546020909101906301000000900460ff16600a811115610fbe57610fbe613119565b600a811115610fcf57610fcf613119565b815281546001600160a01b03600160201b90910481166020830152600183015416604082015260028201805460609092019161100a906135da565b80601f0160208091040260200160405190810160405280929190818152602001828054611036906135da565b80156110835780601f1061105857610100808354040283529160200191611083565b820191906000526020600020905b81548152906001019060200180831161106657829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815250508282815181106110c6576110c66136b7565b602002602001018190525080806110dc906136e0565b915050610eca565b5092915050565b6000816020015160ff1660000361121557336000908152600d6020908152604091829020845181549286015160ff9081166101000261ffff1990941691161791909117808255918401518492829062ff000019166201000083600281111561115557611155613119565b021790555060608201518154829063ff0000001916630100000083600a81111561118157611181613119565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c082015160028201906111df9082613744565b5060e082015160038201556101008201516004820155610120820151600582015561014090910151600690910155506000919050565b602080830151336000908152600d90925260409091205460ff91821691611243916101009004166001613804565b60ff16146112935760405162461bcd60e51b815260206004820152601960248201527f6368756e6b73206d75737420626520636f6e746967756f7573000000000000006044820152606401610429565b336000908152600d6020526040902060c08301516112b49160020190611df9565b336000908152600d6020526040902080546001919082906112de908290610100900460ff16613804565b92506101000a81548160ff021916908360ff160217905550816000015160ff168260200151600161130f9190613804565b60ff161492915050565b600080600080600061133060fe8760c00151611f43565b94509450945094509450846113875760405162461bcd60e51b815260206004820152601960248201527f6661696c65642070726f6f6620766572696669636174696f6e000000000000006044820152606401610429565b8560a001516001600160a01b0316846001600160a01b0316146113bc5760405162461bcd60e51b81526004016104299061381d565b8560600151600a8111156113d2576113d2613119565b83146113f05760405162461bcd60e51b815260040161042990613848565b43821061140f5760405162461bcd60e51b815260040161042990613872565b81600003611425576114226001436136a4565b91505b6001546040516396b477cb60e01b8152600481018490526000916001600160a01b0316906396b477cb90602401602060405180830381865afa15801561146f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114939190613675565b60e08801849052610100880181905243610120890152610140880183905290506114bc87611faa565b50505050505050565b60008060008060006114dc60fc8760c00151611f43565b94509450945094509450846115335760405162461bcd60e51b815260206004820152601e60248201527f6661696c65642061636375736174696f6e20766572696669636174696f6e00006044820152606401610429565b8560a001516001600160a01b0316846001600160a01b0316146115685760405162461bcd60e51b81526004016104299061381d565b8560600151600a81111561157e5761157e613119565b831461159c5760405162461bcd60e51b815260040161042990613848565b4382106115bb5760405162461bcd60e51b815260040161042990613872565b6001546040516396b477cb60e01b8152600481018490526000916001600160a01b0316906396b477cb90602401602060405180830381865afa158015611605573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116299190613675565b60e08801849052610100880181905243610120890152610140880183905290506114bc8761221e565b600080600080600061166960fd8760c00151611f43565b94509450945094509450846116c05760405162461bcd60e51b815260206004820152601d60248201527f6661696c656420696e6e6f63656e636520766572696669636174696f6e0000006044820152606401610429565b8560a001516001600160a01b0316846001600160a01b0316146116f55760405162461bcd60e51b81526004016104299061381d565b8560600151600a81111561170b5761170b613119565b83146117295760405162461bcd60e51b815260040161042990613848565b4382106117485760405162461bcd60e51b815260040161042990613872565b60e086018290526101408601819052611760866124c7565b505050505050565b6000600a82600a81111561177e5761177e613119565b0361178c5760025b92915050565b600082600a8111156117a0576117a0613119565b036117ac576002611786565b600182600a8111156117c0576117c0613119565b036117cc576002611786565b6002611786565b6011545b601054811015611b1c576000601082815481106117f6576117f66136b7565b90600052602060002001549050806000036118115750611b0a565b61181c6001826136a4565b9050600060028281548110611833576118336136b7565b600091825260209182902060408051610160810182526007909302909101805460ff80821685526101008204811695850195909552929390929184019162010000900416600281111561188857611888613119565b600281111561189957611899613119565b815281546020909101906301000000900460ff16600a8111156118be576118be613119565b600a8111156118cf576118cf613119565b815281546001600160a01b03600160201b90910481166020830152600183015416604082015260028201805460609092019161190a906135da565b80601f0160208091040260200160405190810160405280929190818152602001828054611936906135da565b80156119835780601f1061195857610100808354040283529160200191611983565b820191906000526020600020905b81548152906001019060200180831161196657829003601f168201915b505050505081526020016003820154815260200160048201548152602001600582015481526020016006820154815250509050436003600001548261012001516119cd91906136cd565b11156119da575050601155565b60a08101516001600160a01b03166000908152600c602052604081208190556060820151611a0790611768565b60a08301516001600160a01b03166000908152600e6020908152604080832061010087015184529091529020549091508111611a4557505050611b0a565b60a0820180516001600160a01b039081166000908152600e6020908152604080832061010088015184528252808320869055845184168352600b825280832080546001808201835591855283852001899055600f805491820181559093527f8d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac80290920187905592518151858152938401879052909116917f6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f910160405180910390a25050505b80611b14816136e0565b9150506117d7565b601155565b600080600160009054906101000a90046001600160a01b03166001600160a01b031663c9d97af46040518163ffffffff1660e01b8152600401602060405180830381865afa158015611b77573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b9b9190613675565b905060005b600f54811015611c1257816002600f8381548110611bc057611bc06136b7565b906000526020600020015481548110611bdb57611bdb6136b7565b90600052602060002090600702016004015403611c0057611bfd6001846136cd565b92505b80611c0a816136e0565b915050611ba0565b5060005b600f54811015611dec57611dda6002600f8381548110611c3857611c386136b7565b906000526020600020015481548110611c5357611c536136b7565b600091825260209182902060408051610160810182526007909302909101805460ff808216855261010082048116958501959095529293909291840191620100009004166002811115611ca857611ca8613119565b6002811115611cb957611cb9613119565b815281546020909101906301000000900460ff16600a811115611cde57611cde613119565b600a811115611cef57611cef613119565b815281546001600160a01b03600160201b909104811660208301526001830154166040820152600282018054606090920191611d2a906135da565b80601f0160208091040260200160405190810160405280929190818152602001828054611d56906135da565b8015611da35780601f10611d7857610100808354040283529160200191611da3565b820191906000526020600020905b815481529060010190602001808311611d8657829003601f168201915b50505050508152602001600382015481526020016004820154815260200160058201548152602001600682015481525050846127d8565b80611de4816136e0565b915050611c16565b5061065f600f6000612e1f565b815460026001808316156101000203821604825180820160208110602084100160028114611ea35760018114611ec8578660005260208404602060002001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b80821015611e8d5781518455600184019350602082019150611e70565b815191036101000a9081900402909155506114bc565b60028302826020036101000a846020036101000a6020890151040201850187556114bc565b8660005260208404602060002001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b80821015611f295781518455600184019350602082019150611f0c565b815191036101000a90819004029091555050505050505050565b60008060008060008086516020611f5a91906136cd565b9050611f64612e3d565b60a081838a8c5afa611f7557600080fd5b8051600103611f8357600196505b602081015160408201516060830151608090930151989b919a509850909695509350505050565b6000611fb98260600151611768565b60a08301516001600160a01b03166000908152600e60209081526040808320610100870151845290915290205490915081116120075760405162461bcd60e51b8152600401610429906138a2565b600280546001810182556000829052835160079091027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054602086015160ff9081166101000261ffff1990921693169290921791909117808255604085015185939091839162ff000019909116906201000090849081111561208d5761208d613119565b021790555060608201518154829063ff0000001916630100000083600a8111156120b9576120b9613119565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c082015160028201906121179082613744565b5060e082015160038201556101008201516004820155610120820151600582015561014090910151600690910155600254600090612157906001906136a4565b60a0840180516001600160a01b039081166000908152600b6020908152604080832080546001808201835591855283852001879055600f8054918201815584527f8d1108e10bcb7c27dddfc02ed9d693a074039d026cf4ea4240b40f7d581ac80201869055845184168352600e82528083206101008a01518452825291829020879055925181518781529384018590529394509216917f6b7783718ab8e152c193eb08bf76eed1191fcd1677a23a7fe9d338265aad132f91015b60405180910390a2505050565b60a08101516001600160a01b03166000908152600c6020526040902054156122885760405162461bcd60e51b815260206004820181905260248201527f616c72656164792070726f63657373696e6720616e2061636375736174696f6e6044820152606401610429565b60006122978260600151611768565b60a08301516001600160a01b03166000908152600e60209081526040808320610100870151845290915290205490915081116122e55760405162461bcd60e51b8152600401610429906138a2565b600280546001810182556000829052835160079091027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054602086015160ff9081166101000261ffff1990921693169290921791909117808255604085015185939091839162ff000019909116906201000090849081111561236b5761236b613119565b021790555060608201518154829063ff0000001916630100000083600a81111561239757612397613119565b021790555060808201518154640100000000600160c01b031916600160201b6001600160a01b039283160217825560a08301516001830180546001600160a01b0319169190921617905560c082015160028201906123f59082613744565b5060e082015160038201556101008201516004820155610120820151600582015561014090910151600690910155600254600090612435906001906136a4565b90506124428160016136cd565b60a08401516001600160a01b03166000908152600c6020526040902055601061246c8260016136cd565b81546001810183556000928352602092839020015560a0840151604080518581529283018490526001600160a01b03909116917f2e8e354b41470731dafa7c3df150e9498a8d5b9c51ff0259fbf77f721ba403519101612211565b60a08101516001600160a01b03166000908152600c6020526040812054908190036125345760405162461bcd60e51b815260206004820152601860248201527f6e6f206173736f6369617465642061636375736174696f6e00000000000000006044820152606401610429565b8160600151600a81111561254a5761254a613119565b60026125576001846136a4565b81548110612567576125676136b7565b60009182526020909120600790910201546301000000900460ff16600a81111561259357612593613119565b146125f05760405162461bcd60e51b815260206004820152602760248201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e206044820152661c9d5b19481a5960ca1b6064820152608401610429565b60e082015160026126026001846136a4565b81548110612612576126126136b7565b9060005260206000209060070201600301541461267f5760405162461bcd60e51b815260206004820152602560248201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e20604482015264626c6f636b60d81b6064820152608401610429565b61014082015160026126926001846136a4565b815481106126a2576126a26136b7565b9060005260206000209060070201600601541461270d5760405162461bcd60e51b8152602060048201526024808201527f756e6d61746368696e672070726f6f6620616e642061636375736174696f6e206044820152630d0c2e6d60e31b6064820152608401610429565b6011545b60105481101561277757816010828154811061272f5761272f6136b7565b90600052602060002001540361276557600060108281548110612754576127546136b7565b600091825260209091200155612777565b8061276f816136e0565b915050612711565b5060a0820180516001600160a01b039081166000908152600c602090815260408083208390559351935191825292909116917f1fa96beb8dddcb7d4484dd00c4059e872439f7a474a2ecf49c430fc6e86c9e1f910160405180910390a25050565b60015460a0830151604051630c825d9760e11b81526001600160a01b0391821660048201526000929190911690631904bb2e90602401600060405180830381865afa15801561282b573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052612853919081019061345c565b608084015160a08501516001600160a01b039081166000908152600a6020526040902080546001600160a01b031916919092161790559050600381610260015160038111156128a4576128a4613119565b036128ae57505050565b60006128c56128c08560600151611768565b612d5a565b610220830151600754919250906000906128df90836138e6565b6006546128ec90876138e6565b6128f690856136cd565b61290091906136cd565b60095490915081111561291257506009545b60008461012001518560c001518660a0015161292e91906136cd565b61293891906136cd565b60095490915060009061294b83856138e6565b61295591906138fd565b905060008111801561296657508181145b15612a8c57600060a087018190526101008701819052610120870181905260c08701526101e08601805182919061299e9083906136cd565b90525061022086018051600191906129b79083906136cd565b905250600361026087015260006102008701526001546040516301adf0b760e51b81526001600160a01b03909116906335be16e0906129fa90899060040161392f565b600060405180830381600087803b158015612a1457600080fd5b505af1158015612a28573d6000803e3d6000fd5b50505050602086810151604080516001600160a01b0390921682529181018390526000818301526001606082015290517f876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee89181900360800190a15050505050505050565b61012086015181908111612ab957808761012001818151612aad91906136a4565b90525060009050612ad4565b610120870151612ac990826136a4565b600061012089015290505b8015612b51578087610100015110612b1c57808761010001818151612af991906136a4565b90525060a087018051829190612b109083906136a4565b90525060009050612b51565b610100870151612b2c90826136a4565b90508661010001518760a001818151612b4591906136a4565b90525060006101008801525b600081118015612b74575060008760a001518860c00151612b7291906136cd565b115b15612c205760008760a001518860c00151612b8f91906136cd565b60c0890151612b9e90846138e6565b612ba891906138fd565b905060008860a001518960c00151612bc091906136cd565b60a08a0151612bcf90856138e6565b612bd991906138fd565b9050818960c001818151612bed91906136a4565b90525060a089018051829190612c049083906136a4565b905250612c1181836136cd565b612c1b90846136a4565b925050505b612c2a81836136a4565b915081876101e001818151612c3f91906136cd565b9052506102208701805160019190612c589083906136cd565b905250600054610220880151600854612c7191906138e6565b612c7b91906138e6565b612c8590436136cd565b61020088015260026102608801526001546040516301adf0b760e51b81526001600160a01b03909116906335be16e090612cc3908a9060040161392f565b600060405180830381600087803b158015612cdd57600080fd5b505af1158015612cf1573d6000803e3d6000fd5b50505050602087810151610200890151604080516001600160a01b039093168352928201859052818301526000606082015290517f876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee89181900360800190a1505050505050505050565b600081612d6957505060055490565b60018203612d7957505060055490565b60028203612d8957505060055490565b60038203612d9957505060055490565b60048203612daa5750612710919050565b50612710919050565b6040805161016081018252600080825260208201819052909182019081526020016000815260200160006001600160a01b0316815260200160006001600160a01b0316815260200160608152602001600081526020016000815260200160008152602001600081525090565b5080546000825590600052602060002090810190610ae19190612e5b565b6040518060a001604052806005906020820280368337509192915050565b5b80821115612e705760008155600101612e5c565b5090565b6001600160a01b0381168114610ae157600080fd5b8035612e9481612e74565b919050565b600060208284031215612eab57600080fd5b8135612eb681612e74565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604051610160810167ffffffffffffffff81118282101715612ef757612ef7612ebd565b60405290565b604051610280810167ffffffffffffffff81118282101715612ef757612ef7612ebd565b604051601f8201601f1916810167ffffffffffffffff81118282101715612f4a57612f4a612ebd565b604052919050565b803560ff81168114612e9457600080fd5b803560038110612e9457600080fd5b8035600b8110612e9457600080fd5b600067ffffffffffffffff821115612f9b57612f9b612ebd565b50601f01601f191660200190565b600082601f830112612fba57600080fd5b8135612fcd612fc882612f81565b612f21565b818152846020838601011115612fe257600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561301157600080fd5b813567ffffffffffffffff8082111561302957600080fd5b90830190610160828603121561303e57600080fd5b613046612ed3565b61304f83612f52565b815261305d60208401612f52565b602082015261306e60408401612f63565b604082015261307f60608401612f72565b606082015261309060808401612e89565b60808201526130a160a08401612e89565b60a082015260c0830135828111156130b857600080fd5b6130c487828601612fa9565b60c08301525060e08381013590820152610100808401359082015261012080840135908201526101409283013592810192909252509392505050565b60006020828403121561311257600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b6003811061313f5761313f613119565b9052565b600b811061313f5761313f613119565b60005b8381101561316e578181015183820152602001613156565b50506000910152565b6000815180845261318f816020860160208601613153565b601f01601f19169290920160200192915050565b600061016060ff8e16835260ff8d1660208401526131c4604084018d61312f565b6131d1606084018c613143565b6001600160a01b038a81166080850152891660a084015260c083018190526131fb81840189613177565b9150508560e08301528461010083015283610120830152826101408301529c9b505050505050505050505050565b60008060006060848603121561323e57600080fd5b833561324981612e74565b925061325760208501612f72565b9150604084013590509250925092565b60006020828403121561327957600080fd5b81358015158114612eb657600080fd5b805160ff168252600061016060208301516132a9602086018260ff169052565b5060408301516132bc604086018261312f565b5060608301516132cf6060860182613143565b5060808301516132ea60808601826001600160a01b03169052565b5060a083015161330560a08601826001600160a01b03169052565b5060c08301518160c086015261331d82860182613177565b60e0858101519087015261010080860151908701526101208086015190870152610140948501519490950193909352509192915050565b602081526000612eb66020830184613289565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156133bc57603f198886030184526133aa858351613289565b9450928501929085019060010161338e565b5092979650505050505050565b600080604083850312156133dc57600080fd5b82356133e781612e74565b946020939093013593505050565b8051612e9481612e74565b600082601f83011261341157600080fd5b815161341f612fc882612f81565b81815284602083860101111561343457600080fd5b613445826020830160208701613153565b949350505050565b805160048110612e9457600080fd5b60006020828403121561346e57600080fd5b815167ffffffffffffffff8082111561348657600080fd5b90830190610280828603121561349b57600080fd5b6134a3612efd565b6134ac836133f5565b81526134ba602084016133f5565b60208201526134cb604084016133f5565b60408201526060830151828111156134e257600080fd5b6134ee87828601613400565b6060830152506080830151608082015260a083015160a082015260c083015160c082015260e083015160e082015261010080840151818301525061012080840151818301525061014080840151818301525061016080840151818301525061018061355a8185016133f5565b908201526101a083810151908201526101c080840151908201526101e080840151908201526102008084015190820152610220808401519082015261024080840151838111156135a957600080fd5b6135b588828701613400565b82840152505061026091506135cb82840161344d565b91810191909152949350505050565b600181811c908216806135ee57607f821691505b60208210810361360e57634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526024908201527f66756e6374696f6e207265737472696374656420746f207468652076616c696460408201526330ba37b960e11b606082015260800190565b60006020828403121561366a57600080fd5b8151612eb681612e74565b60006020828403121561368757600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b818103818111156117865761178661368e565b634e487b7160e01b600052603260045260246000fd5b808201808211156117865761178661368e565b6000600182016136f2576136f261368e565b5060010190565b601f82111561373f57600081815260208120601f850160051c810160208610156137205750805b601f850160051c820191505b818110156117605782815560010161372c565b505050565b815167ffffffffffffffff81111561375e5761375e612ebd565b6137728161376c84546135da565b846136f9565b602080601f8311600181146137a7576000841561378f5750858301515b600019600386901b1c1916600185901b178555611760565b600085815260208120601f198616915b828110156137d6578886015182559484019460019091019084016137b7565b50858210156137f45787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b60ff81811683821601908111156117865761178661368e565b6020808252601190820152700decccccadcc8cae440dad2e6dac2e8c6d607b1b604082015260600190565b60208082526010908201526f0e4ead8ca40d2c840dad2e6dac2e8c6d60831b604082015260600190565b60208082526016908201527563616e277420626520696e207468652066757475726560501b604082015260600190565b60208082526024908201527f616c726561647920736c6173686564206174207468652070726f6f66277320656040820152630e0dec6d60e31b606082015260800190565b80820281158282048414176117865761178661368e565b60008261391a57634e487b7160e01b600052601260045260246000fd5b500490565b6004811061313f5761313f613119565b602081526139496020820183516001600160a01b03169052565b6000602083015161396560408401826001600160a01b03169052565b5060408301516001600160a01b03811660608401525060608301516102808060808501526139976102a0850183613177565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a0613a1e818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f190161026080880191909152909150613a808483613177565b935080870151915050613a958286018261391f565b509094935050505056fea264697066735822122053c55cbbf35b9b9960584b8709c48b4e91dd095507fd7f796a2db0aa018987ba64736f6c63430008150033",
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
// Solidity: function events(uint256 ) view returns(uint8 chunks, uint8 chunkId, uint8 eventType, uint8 rule, address reporter, address offender, bytes rawProof, uint256 block, uint256 epoch, uint256 reportingBlock, uint256 messageHash)
func (_Accountability *Accountability) Events(opts *runOptions, arg0 *big.Int) (struct {
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
	outstruct.Block = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Epoch = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.ReportingBlock = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)
	outstruct.MessageHash = *abi.ConvertType(out[10], new(*big.Int)).(**big.Int)
	return *outstruct, consumed, err

}

// GetValidatorAccusation is a free data retrieval call binding the contract method 0x9cb22b06.
//
// Solidity: function getValidatorAccusation(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256))
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
// Solidity: function getValidatorFaults(address _val) view returns((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256)[])
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

// HandleEvent is a paid mutator transaction binding the contract method 0x069f6863.
//
// Solidity: function handleEvent((uint8,uint8,uint8,uint8,address,address,bytes,uint256,uint256,uint256,uint256) _event) returns()
func (_Accountability *Accountability) HandleEvent(opts *runOptions, _event AccountabilityEvent) (uint64, error) {
	_, consumed, err := _Accountability.call(opts, "handleEvent", _event)
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
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterSlashingEvent is a free log retrieval operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
 		func (_Accountability *Accountability) FilterSlashingEvent(opts *bind.FilterOpts) (*AccountabilitySlashingEventIterator, error) {






			logs, sub, err := _Accountability.contract.FilterLogs(opts, "SlashingEvent")
			if err != nil {
				return nil, err
			}
			return &AccountabilitySlashingEventIterator{contract: _Accountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
 		}

		// WatchSlashingEvent is a free log subscription operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
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

		// ParseSlashingEvent is a log parse operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator[]\",\"name\":\"_validators\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Config\",\"name\":\"_config\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"ActivatedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"BondingRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewUnbondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeInitialization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getEpochFromBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_signatures\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAccountability\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAccountabilityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIACU\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAcuContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setOracleContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIStabilization\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setStabilizationContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISupplyControl\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSupplyControlContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractUpgradeManager\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setUpgradeManagerContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"_val\",\"type\":\"tuple\"}],\"name\":\"updateValidatorAndTransferSlashedFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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
	Bin: "0x60806040526000600b556000600c553480156200001b57600080fd5b506040516200a1343803806200a1348339810160408190526200003e9162000e64565b601c546000036200006757602a80546001600160a01b031916331790556200006782826200006f565b505062001325565b80518051600d55602080820151600e55604080830151600f55606080840151601055608093840151601180546001600160a01b03199081166001600160a01b03938416179091558487015180516012805484169185169190911790558086015160138054841691851691909117905580850151601480548416918516919091179055808401516015805484169185169190911790559586015160168054831691841691909117905560a0909501516017805487169183169190911790558286015180516018805490971692169190911790945591830151601955820151601a5590810151601b55810151601c5560005b825181101562000420576000838281518110620001805762000180620010a1565b602002602001015160a0015190506000848381518110620001a557620001a5620010a1565b60200260200101516101a00181815250506000848381518110620001cd57620001cd620010a1565b602002602001015161018001906001600160a01b031690816001600160a01b0316815250506000848381518110620002095762000209620010a1565b602002602001015160a00181815250506000848381518110620002305762000230620010a1565b60209081029190910101516101c00152600f5484518590849081106200025a576200025a620010a1565b602002602001015160800181815250506000848381518110620002815762000281620010a1565b602002602001015161026001906003811115620002a257620002a2620010b7565b90816003811115620002b857620002b8620010b7565b815250506000848381518110620002d357620002d3620010a1565b602002602001015161016001818152505062000311848381518110620002fd57620002fd620010a1565b60200260200101516200042560201b60201c565b6200033e8483815181106200032a576200032a620010a1565b60200260200101516200056060201b60201c565b8060276000868581518110620003585762000358620010a1565b6020026020010151600001516001600160a01b03166001600160a01b031681526020019081526020016000206000828254620003959190620010e3565b925050819055508060296000828254620003b09190620010e3565b925050819055506200040a848381518110620003d057620003d0620010a1565b60200260200101516020015182868581518110620003f257620003f2620010a1565b6020026020010151600001516200079160201b60201c565b50806200041781620010ff565b9150506200015f565b505050565b60006200043c82606001516200097c60201b60201c565b6001600160a01b03909116602084015290508015620004905760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b60448201526064015b60405180910390fd5b6020808301516001600160a01b03908116600090815260289092526040909120600101541615620005045760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000487565b612710826080015111156200055c5760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000487565b5050565b6101808101516001600160a01b0316620005e357601d546000906200058590620009ca565b905081602001518260000151836080015183604051620005a59062000aea565b620005b494939291906200111b565b604051809103906000f080158015620005d1573d6000803e3d6000fd5b506001600160a01b0316610180830152505b60208181018051601d80546001808201835560009283527f6d4407e7be21f808e6509aa9fa9143369579dd7d760fe20a2c09680fc146134f90910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602890955260409182902086518154871690851617815593519084018054861691841691909117905584015160028301805490941691161790915560608201518291906003820190620006959082620011ff565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f8201556102008201516010820155610220820151601182015561024082015160128201906200075f9082620011ff565b5061026082015160138201805460ff19166001836003811115620007875762000787620010b7565b0217905550505050565b60008211620007ef5760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000487565b6001600160a01b038116600090815260276020526040902054821115620008595760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000487565b6001600160a01b0381166000908152602760205260408120805484929062000883908490620012cb565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580546000908152600394859052978820875181549088166001600160a01b03199182161782559551600182018054919098169616959095179095559051600284015551910155805491926200090f83620010ff565b90915550506001600160a01b03848116600081815260286020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d910160405180910390a35050505050565b6000806200098962000af8565b60008060ff9050604083875160208901845afa620009a657600080fd5b505080516020909101516c0100000000000000000000000090910494909350915050565b606081600003620009f25750506040805180820190915260018152600360fc1b602082015290565b8160005b811562000a22578062000a0981620010ff565b915062000a1a9050600a83620012f7565b9150620009f6565b6000816001600160401b0381111562000a3f5762000a3f62000b16565b6040519080825280601f01601f19166020018201604052801562000a6a576020820181803683370190505b5090505b841562000ae25762000a82600183620012cb565b915062000a91600a866200130e565b62000a9e906030620010e3565b60f81b81838151811062000ab65762000ab6620010a1565b60200101906001600160f81b031916908160001a90535062000ada600a86620012f7565b945062000a6e565b949350505050565b6115318062008c0383390190565b60405180604001604052806002906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b604051608081016001600160401b038111828210171562000b515762000b5162000b16565b60405290565b60405160a081016001600160401b038111828210171562000b515762000b5162000b16565b60405161028081016001600160401b038111828210171562000b515762000b5162000b16565b604051601f8201601f191681016001600160401b038111828210171562000bcd5762000bcd62000b16565b604052919050565b6001600160a01b038116811462000beb57600080fd5b50565b805162000bfb8162000bd5565b919050565b60005b8381101562000c1d57818101518382015260200162000c03565b50506000910152565b600082601f83011262000c3857600080fd5b81516001600160401b0381111562000c545762000c5462000b16565b62000c69601f8201601f191660200162000ba2565b81815284602083860101111562000c7f57600080fd5b62000ae282602083016020870162000c00565b80516004811062000bfb57600080fd5b600060c0828403121562000cb557600080fd5b60405160c081016001600160401b038111828210171562000cda5762000cda62000b16565b8060405250809150825162000cef8162000bd5565b8152602083015162000d018162000bd5565b6020820152604083015162000d168162000bd5565b6040820152606083015162000d2b8162000bd5565b6060820152608083015162000d408162000bd5565b608082015260a083015162000d558162000bd5565b60a0919091015292915050565b60006080828403121562000d7557600080fd5b62000d7f62000b2c565b9050815162000d8e8162000bd5565b8082525060208201516020820152604082015160408201526060820151606082015292915050565b600081830361020081121562000dcb57600080fd5b62000dd562000b2c565b915060a081121562000de657600080fd5b5062000df162000b57565b82518152602083015160208201526040830151604082015260608301516060820152608083015162000e238162000bd5565b6080820152815262000e398360a0840162000ca2565b602082015262000e4e83610160840162000d62565b60408201526101e0820151606082015292915050565b60008061022080848603121562000e7a57600080fd5b83516001600160401b038082111562000e9257600080fd5b818601915086601f83011262000ea757600080fd5b815160208282111562000ebe5762000ebe62000b16565b8160051b62000ecf82820162000ba2565b928352848101820192828101908b85111562000eea57600080fd5b83870192505b848310156200107f5782518681111562000f0957600080fd5b8701610280818e03601f1901121562000f2157600080fd5b62000f2b62000b7c565b62000f3886830162000bee565b815262000f486040830162000bee565b8682015262000f5a6060830162000bee565b604082015260808201518881111562000f7257600080fd5b62000f828f888386010162000c26565b60608301525060a0820151608082015260c082015160a082015260e082015160c082015261010082015160e082015261012082015161010082015261014082015161012082015261016082015161014082015261018082015161016082015262000ff06101a0830162000bee565b6101808201526101c08201516101a08201526101e08201516101c08201526102008201516101e0820152898201516102008201526102408201518a820152610260820151888111156200104257600080fd5b620010528f888386010162000c26565b6102408301525062001068610280830162000c92565b610260820152835250918301919083019062000ef0565b8099505050506200109389828a0162000db6565b955050505050509250929050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052601160045260246000fd5b80820180821115620010f957620010f9620010cd565b92915050565b600060018201620011145762001114620010cd565b5060010190565b600060018060a01b0380871683528086166020840152508360408301526080606083015282518060808401526200115a8160a085016020870162000c00565b601f01601f19169190910160a00195945050505050565b600181811c908216806200118657607f821691505b602082108103620011a757634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200042057600081815260208120601f850160051c81016020861015620011d65750805b601f850160051c820191505b81811015620011f757828155600101620011e2565b505050505050565b81516001600160401b038111156200121b576200121b62000b16565b62001233816200122c845462001171565b84620011ad565b602080601f8311600181146200126b5760008415620012525750858301515b600019600386901b1c1916600185901b178555620011f7565b600085815260208120601f198616915b828110156200129c578886015182559484019460019091019084016200127b565b5085821015620012bb5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b81810381811115620010f957620010f9620010cd565b634e487b7160e01b600052601260045260246000fd5b600082620013095762001309620012e1565b500490565b600082620013205762001320620012e1565b500690565b6178ce80620013356000396000f3fe608060405260043610620003db5760003560e01c8063872cf0591162000203578063b66b3e791162000117578063d372c07e11620000a7578063dd62ed3e1162000075578063dd62ed3e1462000cbd578063dfb1a4d21462000d07578063e7f43c681462000d1e578063f7866ee31462000d3e57005b8063d372c07e1462000c39578063d5f394881462000c5e578063d861b0e81462000c80578063d886f8a21462000c9857005b8063cb696f5411620000e5578063cb696f541462000bb2578063ceaad4551462000bd7578063cf9c57191462000bfc578063cfd19fb91462000c1457005b8063b66b3e791462000b42578063b7ab4db51462000b6a578063c2362dd51462000b82578063c9d97af41462000b9a57005b8063a5d059ca1162000193578063ae1f5fa01162000161578063ae1f5fa01462000aac578063b2ea9adb1462000ad3578063b3ecbadd1462000af8578063b46e55201462000b1d57005b8063a5d059ca1462000a14578063a8b2216e1462000a39578063a9059cbb1462000a60578063ab8f6ffe1462000a8557005b80639bb851c011620001d15780639bb851c0146200099a5780639c98e47114620009b25780639dc29fac14620009ca578063a515366a14620009ef57005b8063872cf05914620008fe5780638bac7dad146200091657806395d89b41146200093b57806396b477cb146200096957005b806340c10f1911620002fb57806370a08231116200028b578063819b64631162000259578063819b6463146200087d578063833b1fce146200089457806384467fdb14620008b4578063852c484914620008d957005b806370a082311462000729578063731b3a03146200076357806377e741c7146200077a57806379502c55146200079f57005b8063520fdbbc11620002c9578063520fdbbc146200068a5780635f7d394914620006af5780636b5f444c14620006ed5780636fd2c80b146200071257005b806340c10f191462000600578063436459691462000625578063496ccd9b146200063c5780634bb278f3146200066157005b80631604e416116200037757806329070c6d116200034557806329070c6d146200058e5780632f2c3f2e14620005a5578063313ce56714620005bd57806335be16e014620005db57005b80631604e416146200050657806318160ddd146200051e5780631904bb2e146200053557806323b872dd146200056957005b80630d8e6e2c11620003b55780630d8e6e2c14620004845780631122063314620004a5578063114eaf5514620004bc5780631250a28d14620004e157005b806306fdde0314620003e5578063095ea7b314620004295780630ae65e7a146200045f57005b36620003e357005b005b348015620003f257600080fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b6040516200042091906200531e565b60405180910390f35b3480156200043657600080fd5b506200044e6200044836600462005350565b62000d5e565b604051901515815260200162000420565b3480156200046c57600080fd5b50620003e36200047e3660046200537f565b62000d77565b3480156200049157600080fd5b50601c545b60405190815260200162000420565b348015620004b257600080fd5b50600e5462000496565b348015620004c957600080fd5b50620003e3620004db3660046200539f565b62000e0a565b348015620004ee57600080fd5b50620003e3620005003660046200537f565b62000e3c565b3480156200051357600080fd5b506200049660245481565b3480156200052b57600080fd5b5060295462000496565b3480156200054257600080fd5b506200055a620005543660046200537f565b62000e8b565b604051620004209190620053f2565b3480156200057657600080fd5b506200044e620005883660046200556b565b62001139565b3480156200059b57600080fd5b50600d5462000496565b348015620005b257600080fd5b506200049661271081565b348015620005ca57600080fd5b506040516012815260200162000420565b348015620005e857600080fd5b50620003e3620005fa366004620055b1565b62001193565b3480156200060d57600080fd5b50620003e36200061f36600462005350565b6200136f565b3480156200063257600080fd5b50601a5462000496565b3480156200064957600080fd5b50620003e36200065b3660046200537f565b62001429565b3480156200066e57600080fd5b506200067962001533565b604051620004209392919062005672565b3480156200069757600080fd5b50620003e3620006a93660046200537f565b62001923565b348015620006bc57600080fd5b50620006d4620006ce3660046200569f565b62001b1e565b6040516001600160a01b03909116815260200162000420565b348015620006fa57600080fd5b50620003e36200070c3660046200539f565b62001d35565b3480156200071f57600080fd5b5060105462000496565b3480156200073657600080fd5b5062000496620007483660046200537f565b6001600160a01b031660009081526027602052604090205490565b3480156200077057600080fd5b5060205462000496565b3480156200078757600080fd5b50620003e3620007993660046200539f565b62001ebb565b348015620007ac57600080fd5b506040805160a08082018352600d548252600e54602080840191909152600f54838501526010546060808501919091526011546001600160a01b03908116608080870191909152865160c0810188526012548316815260135483168186015260145483168189015260155483168185015260165483168183015260175483169581019590955286519081018752601854909116815260195492810192909252601a5494820194909452601b5493810193909352601c546200086b939084565b604051620004209493929190620056c2565b3480156200088a57600080fd5b50601b5462000496565b348015620008a157600080fd5b506013546001600160a01b0316620006d4565b348015620008c157600080fd5b50620003e3620008d33660046200583c565b62001eed565b348015620008e657600080fd5b50620003e3620008f836600462005350565b6200201c565b3480156200090b57600080fd5b50620003e3620021ad565b3480156200092357600080fd5b50620003e3620009353660046200539f565b620021e9565b3480156200094857600080fd5b50604080518082019091526003815262272a2760e91b602082015262000411565b3480156200097657600080fd5b5062000496620009883660046200539f565b6000908152601f602052604090205490565b348015620009a757600080fd5b506200049660235481565b348015620009bf57600080fd5b506200049660215481565b348015620009d757600080fd5b50620003e3620009e936600462005350565b6200226d565b348015620009fc57600080fd5b50620003e362000a0e36600462005350565b62002383565b34801562000a2157600080fd5b50620003e362000a3336600462005350565b62002456565b34801562000a4657600080fd5b5062000a51620024ed565b604051620004209190620058e5565b34801562000a6d57600080fd5b506200044e62000a7f36600462005350565b620025d0565b34801562000a9257600080fd5b5062000a9d620025df565b6040516200042091906200594b565b34801562000ab957600080fd5b5062000ac4620026f1565b60405162000420919062005960565b34801562000ae057600080fd5b50620003e362000af2366004620059af565b62002937565b34801562000b0557600080fd5b50620003e362000b173660046200537f565b6200297e565b34801562000b2a57600080fd5b50620003e362000b3c3660046200537f565b620029cd565b34801562000b4f57600080fd5b5062000b5a62002c50565b6040516200042092919062005a1a565b34801562000b7757600080fd5b5062000ac462002d87565b34801562000b8f57600080fd5b506200049660205481565b34801562000ba757600080fd5b5062000496601e5481565b34801562000bbf57600080fd5b50620003e362000bd13660046200539f565b62002deb565b34801562000be457600080fd5b50620003e362000bf63660046200537f565b62002e4e565b34801562000c0957600080fd5b50620003e362002e9d565b34801562000c2157600080fd5b50620003e362000c333660046200537f565b62002ef1565b34801562000c4657600080fd5b50620003e362000c583660046200537f565b62002f40565b34801562000c6b57600080fd5b50602a54620006d4906001600160a01b031681565b34801562000c8d57600080fd5b50620003e362002f8f565b34801562000ca557600080fd5b50620003e362000cb73660046200537f565b62002fd0565b34801562000cca57600080fd5b506200049662000cdc36600462005a4c565b6001600160a01b03918216600090815260266020908152604080832093909416825291909152205490565b34801562000d1457600080fd5b5060195462000496565b34801562000d2b57600080fd5b506018546001600160a01b0316620006d4565b34801562000d4b57600080fd5b506011546001600160a01b0316620006d4565b600062000d6d3384846200301f565b5060015b92915050565b6001600160a01b038082166000818152602860205260409020600101549091161462000dc05760405162461bcd60e51b815260040162000db79062005a8a565b60405180910390fd5b6001600160a01b0381811660009081526028602052604090205416331462000dfc5760405162461bcd60e51b815260040162000db79062005ac1565b62000e078162003148565b50565b6018546001600160a01b0316331462000e375760405162461bcd60e51b815260040162000db79062005b0d565b601055565b6018546001600160a01b0316331462000e695760405162461bcd60e51b815260040162000db79062005b0d565b601280546001600160a01b0319166001600160a01b0392909216919091179055565b62000e95620050f7565b6001600160a01b038083166000818152602860205260409020600101549091161462000ed55760405162461bcd60e51b815260040162000db79062005b44565b6001600160a01b03808316600090815260286020908152604091829020825161028081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000f389062005b7b565b80601f016020809104026020016040519081016040528092919081815260200182805462000f669062005b7b565b801562000fb75780601f1062000f8b5761010080835404028352916020019162000fb7565b820191906000526020600020905b81548152906001019060200180831162000f9957829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546080820152600882015460a0820152600982015460c0820152600a82015460e0820152600b820154610100820152600c8201546001600160a01b0316610120820152600d820154610140820152600e820154610160820152600f82015461018082015260108201546101a082015260118201546101c08201526012820180546101e090920191620010739062005b7b565b80601f0160208091040260200160405190810160405280929190818152602001828054620010a19062005b7b565b8015620010f25780601f10620010c657610100808354040283529160200191620010f2565b820191906000526020600020905b815481529060010190602001808311620010d457829003601f168201915b5050509183525050601382015460209091019060ff1660038111156200111c576200111c620053b9565b6003811115620011305762001130620053b9565b90525092915050565b6000620011488484846200321f565b6001600160a01b03841660009081526026602090815260408083203384529091528120546200117990849062005bcd565b9050620011888533836200301f565b506001949350505050565b6012546001600160a01b03163314620011fb5760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f742074686520736c617368696e6720636f6e74726044820152621858dd60ea1b606482015260840162000db7565b60006101208201356028826200121860408601602087016200537f565b6001600160a01b03166001600160a01b031681526020019081526020016000206009015462001248919062005bcd565b60c0830135602860006200126360408701602088016200537f565b6001600160a01b03166001600160a01b031681526020019081526020016000206006015462001293919062005bcd565b60a084013560286000620012ae60408801602089016200537f565b6001600160a01b03166001600160a01b0316815260200190815260200160002060050154620012de919062005bcd565b620012ea919062005be3565b620012f6919062005be3565b6011546001600160a01b03166000908152602760205260408120805492935083929091906200132790849062005be3565b90915550829050602860006200134460408401602085016200537f565b6001600160a01b03168152602081019190915260400160002062001369828262005ded565b50505050565b6018546001600160a01b031633146200139c5760405162461bcd60e51b815260040162000db79062005b0d565b6001600160a01b03821660009081526027602052604081208054839290620013c690849062005be3565b925050819055508060296000828254620013e1919062005be3565b90915550506040518181526001600160a01b038316907f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf906020015b60405180910390a25050565b6018546001600160a01b03163314620014565760405162461bcd60e51b815260040162000db79062005b0d565b601380546001600160a01b0319166001600160a01b03838116918217909255601454604051637adbf97360e01b8152600481019290925290911690637adbf97390602401600060405180830381600087803b158015620014b557600080fd5b505af1158015620014ca573d6000803e3d6000fd5b5050601654604051637adbf97360e01b81526001600160a01b0385811660048301529091169250637adbf97391506024015b600060405180830381600087803b1580156200151757600080fd5b505af11580156200152c573d6000803e3d6000fd5b5050505050565b602a5460009060609082906001600160a01b03163314620015685760405162461bcd60e51b815260040162000db79062005f43565b601e54436000818152601f6020908152604082209390935560195492549092620015929162005be3565b6012546040516306c9789b60e41b8152929091146004830181905292506001600160a01b031690636c9789b090602401600060405180830381600087803b158015620015dd57600080fd5b505af1158015620015f2573d6000803e3d6000fd5b505050508015620016e7576200160762003328565b6200161162003712565b6200161b62003801565b600062001627620026f1565b60135460405163422811f960e11b81529192506001600160a01b03169063845023f2906200165a90849060040162005960565b600060405180830381600087803b1580156200167557600080fd5b505af11580156200168a573d6000803e3d6000fd5b50505050436020819055506001601e6000828254620016aa919062005be3565b9091555050601e546040519081527febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e3359060200160405180910390a1505b60135460408051634bb278f360e01b815290516000926001600160a01b031691634bb278f3916004808301926020929190829003018187875af115801562001733573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001759919062005f86565b90508015620017d5576014546040805163a2e6204560e01b815290516001600160a01b039092169163a2e620459160048082019260209290919082900301816000875af1925050508015620017cd575060408051601f3d908101601f19168201909252620017ca9181019062005f86565b60015b15620017d557505b81156200190a576002546020805460228054604080518286028101860190915281815260ff90951694919384919060009084015b82821015620018f7576000848152602090819020604080516060810182526003860290920180546001600160a01b03168352600181015493830193909352600283018054929392918401916200185f9062005b7b565b80601f01602080910402602001604051908101604052809291908181526020018280546200188d9062005b7b565b8015620018de5780601f10620018b257610100808354040283529160200191620018de565b820191906000526020600020905b815481529060010190602001808311620018c057829003601f168201915b5050505050815250508152602001906001019062001809565b5050505091509450945094505050909192565b505060025460205460ff90911694606094509092509050565b6018546001600160a01b03163314620019505760405162461bcd60e51b815260040162000db79062005b0d565b601880546001600160a01b0319166001600160a01b0383811691821790925560135460405163b3ab15fb60e01b815260048101929092529091169063b3ab15fb90602401600060405180830381600087803b158015620019af57600080fd5b505af1158015620019c4573d6000803e3d6000fd5b505060145460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001a1057600080fd5b505af115801562001a25573d6000803e3d6000fd5b505060155460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001a7157600080fd5b505af115801562001a86573d6000803e3d6000fd5b505060165460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001ad257600080fd5b505af115801562001ae7573d6000803e3d6000fd5b505060175460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401620014fc565b600080805b60225481101562001b7a576022818154811062001b445762001b4462005faa565b9060005260206000209060030201600101548262001b63919062005be3565b91508062001b718162005fc0565b91505062001b23565b508060000362001bcd5760405162461bcd60e51b815260206004820152601c60248201527f54686520636f6d6d6974746565206973206e6f74207374616b696e6700000000604482015260640162000db7565b60008362001bdd60638762005fdc565b62001be9919062005be3565b905060008160405160200162001c0191815260200190565b60408051601f1981840301815291905280516020909101209050600062001c2984836200600c565b90506000805b60225481101562001cd9576022818154811062001c505762001c5062005faa565b9060005260206000209060030201600101548262001c6f919062005be3565b915062001c7e60018362005bcd565b831162001cc4576022818154811062001c9b5762001c9b62005faa565b60009182526020909120600390910201546001600160a01b0316965062000d7195505050505050565b8062001cd08162005fc0565b91505062001c2f565b5060405162461bcd60e51b815260206004820152602960248201527f5468657265206973206e6f2076616c696461746f72206c65667420696e20746860448201526865206e6574776f726b60b81b606482015260840162000db7565b6018546001600160a01b0316331462001d625760405162461bcd60e51b815260040162000db79062005b0d565b60195481101562001e19578060205462001d7d919062005be3565b431062001e195760405162461bcd60e51b815260206004820152605760248201527f63757272656e7420636861696e2068656164206578636565642074686520776960448201527f6e646f773a206c617374426c6f636b45706f6368202b205f6e6577506572696f60648201527f642c2074727920616761696e206c6174746572206f6e2e000000000000000000608482015260a40162000db7565b6019819055601254604051631ad7d11360e21b8152600481018390526001600160a01b0390911690636b5f444c90602401600060405180830381600087803b15801562001e6557600080fd5b505af115801562001e7a573d6000803e3d6000fd5b505050507fd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f818160405162001eb091815260200190565b60405180910390a150565b6018546001600160a01b0316331462001ee85760405162461bcd60e51b815260040162000db79062005b0d565b600d55565b6000604051806102800160405280336001600160a01b0316815260200160006001600160a01b03168152602001856001600160a01b03168152602001868152602001600d6000016002015481526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b03168152602001600081526020014381526020016000815260200160008152602001600081526020018481526020016000600381111562001fbc5762001fbc620053b9565b9052905062001fcc818362003920565b60208101516101808201516040517f8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c926200200d92339289918b9162006023565b60405180910390a15050505050565b6001600160a01b03808316600081815260286020526040902060010154909116146200205c5760405162461bcd60e51b815260040162000db79062005a8a565b6001600160a01b03828116600090815260286020526040902054163314620020985760405162461bcd60e51b815260040162000db79062005ac1565b612710811115620020ec5760405162461bcd60e51b815260206004820152601f60248201527f7265717569726520636f727265637420636f6d6d697373696f6e207261746500604482015260640162000db7565b604080516060810182526001600160a01b038481168252436020808401918252838501868152600c80546000908152600a909352958220855181546001600160a01b0319169516949094178455915160018085019190915591516002909301929092558354929390929091906200216590849062005be3565b90915550506040518281526001600160a01b038416907f4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf9060200160405180910390a2505050565b6018546001600160a01b03163314620021da5760405162461bcd60e51b815260040162000db79062005b0d565b6002805460ff19166001179055565b6018546001600160a01b03163314620022165760405162461bcd60e51b815260040162000db79062005b0d565b60008111620022685760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000db7565b601b55565b6018546001600160a01b031633146200229a5760405162461bcd60e51b815260040162000db79062005b0d565b6001600160a01b038216600090815260276020526040902054811115620022fd5760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000db7565b6001600160a01b038216600090815260276020526040812080548392906200232790849062005bcd565b92505081905550806029600082825462002342919062005bcd565b90915550506040518181526001600160a01b038316907f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3906020016200141d565b6001600160a01b0380831660008181526028602052604090206001015490911614620023c35760405162461bcd60e51b815260040162000db79062005b44565b6001600160a01b03821660009081526028602052604081206013015460ff166003811115620023f657620023f6620053b9565b14620024455760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000db7565b6200245282823362003d7a565b5050565b6001600160a01b0380831660008181526028602052604090206001015490911614620024965760405162461bcd60e51b815260040162000db79062005b44565b60008111620024e05760405162461bcd60e51b81526020600482015260156024820152740756e626f6e64696e6720616d6f756e74206973203605c1b604482015260640162000db7565b6200245282823362003f66565b60606025805480602002602001604051908101604052809291908181526020016000905b82821015620025c7578382906000526020600020018054620025339062005b7b565b80601f0160208091040260200160405190810160405280929190818152602001828054620025619062005b7b565b8015620025b25780601f106200258657610100808354040283529160200191620025b2565b820191906000526020600020905b8154815290600101906020018083116200259457829003601f168201915b50505050508152602001906001019062002511565b50505050905090565b600062000d6d3384846200321f565b60606022805480602002602001604051908101604052809291908181526020016000905b82821015620025c7576000848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620026599062005b7b565b80601f0160208091040260200160405190810160405280929190818152602001828054620026879062005b7b565b8015620026d85780601f10620026ac57610100808354040283529160200191620026d8565b820191906000526020600020905b815481529060010190602001808311620026ba57829003601f168201915b5050505050815250508152602001906001019062002603565b602a546060906001600160a01b03163314620027215760405162461bcd60e51b815260040162000db79062005f43565b601d54620027725760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000db7565b6200277c620051c7565b601b546080820152601d8152602860208201526022604082015260216060820152620027a88162004299565b620027b660256000620051e5565b60225480620027fd5760405162461bcd60e51b8152602060048201526012602482015271636f6d6d697474656520697320656d70747960701b604482015260640162000db7565b60008167ffffffffffffffff8111156200281b576200281b62005791565b60405190808252806020026020018201604052801562002845578160200160208202803683370190505b50905060005b828110156200292f576000602860006022848154811062002870576200287062005faa565b60009182526020808320600392830201546001600160a01b031684528301939093526040909101812060258054600181018255925292507f401968ff42a154441da5f6c4c935ac46b8671f0e062baaa62a7545ba53bb6e4c0190620028d8908301826200606c565b50600281015483516001600160a01b039091169084908490811062002901576200290162005faa565b6001600160a01b03909216602092830291909101909101525080620029268162005fc0565b9150506200284b565b509250505090565b6018546001600160a01b03163314620029645760405162461bcd60e51b815260040162000db79062005b0d565b62002971600083620042b4565b62002452600182620042b4565b6018546001600160a01b03163314620029ab5760405162461bcd60e51b815260040162000db79062005b0d565b601580546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038082166000818152602860205260409020600101549091161462002a0d5760405162461bcd60e51b815260040162000db79062005a8a565b6001600160a01b0380821660009081526028602052604090208054909116331462002a935760405162461bcd60e51b815260206004820152602f60248201527f726571756972652063616c6c657220746f2062652076616c696461746f72207460448201526e1c99585cdd5c9e481858d8dbdd5b9d608a1b606482015260840162000db7565b6000601382015460ff16600381111562002ab15762002ab1620053b9565b0362002b005760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f7220616c7265616479206163746976650000000000000000604482015260640162000db7565b6002601382015460ff16600381111562002b1e5762002b1e620053b9565b14801562002b2f5750438160100154115b1562002b7e5760405162461bcd60e51b815260206004820152601760248201527f76616c696461746f72207374696c6c20696e206a61696c000000000000000000604482015260640162000db7565b6003601382015460ff16600381111562002b9c5762002b9c620053b9565b0362002beb5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206a61696c6564207065726d616e656e746c7900000000604482015260640162000db7565b60138101805460ff1916905580546019546020546001600160a01b038581169316917f60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b59162002c3b919062005be3565b60405190815260200160405180910390a35050565b6060806000600181805462002c659062005b7b565b80601f016020809104026020016040519081016040528092919081815260200182805462002c939062005b7b565b801562002ce45780601f1062002cb85761010080835404028352916020019162002ce4565b820191906000526020600020905b81548152906001019060200180831162002cc657829003601f168201915b5050505050915080805462002cf99062005b7b565b80601f016020809104026020016040519081016040528092919081815260200182805462002d279062005b7b565b801562002d785780601f1062002d4c5761010080835404028352916020019162002d78565b820191906000526020600020905b81548152906001019060200180831162002d5a57829003601f168201915b50505050509050915091509091565b6060601d80548060200260200160405190810160405280929190818152602001828054801562002de157602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831162002dc2575b5050505050905090565b6018546001600160a01b0316331462002e185760405162461bcd60e51b815260040162000db79062005b0d565b600e8190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200162001eb0565b6018546001600160a01b0316331462002e7b5760405162461bcd60e51b815260040162000db79062005b0d565b601780546001600160a01b0319166001600160a01b0392909216919091179055565b6018546001600160a01b0316331462002eca5760405162461bcd60e51b815260040162000db79062005b0d565b62002ed760008062005205565b62002ee56001600062005205565b6002805460ff19169055565b6018546001600160a01b0316331462002f1e5760405162461bcd60e51b815260040162000db79062005b0d565b601680546001600160a01b0319166001600160a01b0392909216919091179055565b6018546001600160a01b0316331462002f6d5760405162461bcd60e51b815260040162000db79062005b0d565b601480546001600160a01b0319166001600160a01b0392909216919091179055565b602a546001600160a01b0316331462002fbc5760405162461bcd60e51b815260040162000db79062005f43565b62002fc662003712565b62000e07620026f1565b6018546001600160a01b0316331462002ffd5760405162461bcd60e51b815260040162000db79062005b0d565b601180546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038316620030835760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000db7565b6001600160a01b038216620030e65760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000db7565b6001600160a01b0383811660008181526026602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b038116600090815260286020526040812090601382015460ff1660038111156200317d576200317d620053b9565b14620031cc5760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265206163746976650000000000000000604482015260640162000db7565b60138101805460ff1916600117905580546019546020546001600160a01b038581169316917f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c9162002c3b919062005be3565b6001600160a01b038316600090815260276020526040902054811115620032825760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000db7565b6001600160a01b03831660009081526027602052604081208054839290620032ac90849062005bcd565b90915550506001600160a01b03821660009081526027602052604081208054839290620032db90849062005be3565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200313b91815260200190565b476000036200333357565b600d544790600090670de0b6b3a7640000906200335290849062005fdc565b6200335e919062006155565b90508015620033dc576011546040516000916001600160a01b03169083908381818185875af1925050503d8060008114620033b6576040519150601f19603f3d011682016040523d82523d6000602084013e620033bb565b606091505b5090915050801515600103620033da57620033d7828462005bcd565b92505b505b8160236000828254620033f0919062005be3565b90915550600090505b6022548110156200370d576000602860006022848154811062003420576200342062005faa565b600091825260208083206003909202909101546001600160a01b031683528201929092526040018120602154602280549294509091879190869081106200346b576200346b62005faa565b90600052602060002090600302016001015462003489919062005fdc565b62003495919062006155565b90508015620036f5576002601383015460ff166003811115620034bc57620034bc620053b9565b1480620034e457506003601383015460ff166003811115620034e257620034e2620053b9565b145b156200358c57601254602280546001600160a01b0390921691631de9d9b69184918790811062003518576200351862005faa565b600091825260209091206003909102015460405160e084901b6001600160e01b03191681526001600160a01b0390911660048201526024016000604051808303818588803b1580156200356a57600080fd5b505af11580156200357f573d6000803e3d6000fd5b50505050505050620036f8565b60008260050154828460080154620035a5919062005fdc565b620035b1919062006155565b90506000620035c1828462005bcd565b90508115620036265783546040516001600160a01b03909116906108fc9084906000818181858888f193505050503d80600081146200361d576040519150601f19603f3d011682016040523d82523d6000602084013e62003622565b606091505b5050505b8015620036ae5783600c0160009054906101000a90046001600160a01b03166001600160a01b031663fb489a7b826040518263ffffffff1660e01b815260040160206040518083038185885af115801562003685573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190620036ac91906200616c565b505b60018401546040518481526001600160a01b03909116907fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe5639060200160405180910390a250505b50505b80620037048162005fc0565b915050620033f9565b505050565b6004545b6005548110156200373f576200373981620037318162005fc0565b925062004407565b62003716565b50600554600455600854600754036200375457565b6009545b60085481101562003781576200377b81620037738162005fc0565b92506200462e565b62003758565b50600854600955600754805b600854811015620037fb576010546000828152600660205260409020600401544391620037ba9162005be3565b11620037e057620037cb8162004929565b620037d860018362005be3565b9150620037e6565b620037fb565b80620037f28162005fc0565b9150506200378d565b50600755565b600c54600b5410156200391e57600b546000908152600a6020526040902060105460018201544391620038349162005be3565b11156200383e5750565b600281015481546001600160a01b03908116600090815260286020526040808220600490810185905585548416835291819020600c015490516319fac8fd60e01b81529216926319fac8fd9262003899920190815260200190565b600060405180830381600087803b158015620038b457600080fd5b505af1158015620038c9573d6000803e3d6000fd5b5050600b80546000908152600a6020526040812080546001600160a01b0319168155600180820183905560029091018290558254909450919250906200391190849062005be3565b9091555062003801915050565b565b60e28151146200396a5760405162461bcd60e51b8152602060048201526014602482015273092dcecc2d8d2c840e0e4dedecc40d8cadccee8d60631b604482015260640162000db7565b60308261024001515114620039c25760405162461bcd60e51b815260206004820152601c60248201527f496e76616c696420636f6e73656e737573206b6579206c656e67746800000000604482015260640162000db7565b620039cd8262004a74565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a00000000000060208083019190915284519251919260009262003a30920160609190911b6bffffffffffffffffffffffff1916815260140190565b604051602081830303815290604052905060008262003a50835162004ba1565b8360405160200162003a659392919062006186565b60408051601f198184030181528282528051602091820120600280855260608501845290945060009392909183019080368337019050509050600080808062003abd898262003ab76041600262005fdc565b62004cc2565b9050600062003adc8a62003ad46041600262005fdc565b606062004cc2565b905060205b825181101562003bad5762003af7838262004ddb565b6040805160008152602081018083528d905260ff8316918101919091526060810184905260808101839052929850909650945060019060a0016020604051602081039080840390855afa15801562003b53573d6000803e3d6000fd5b5050604051601f19015190508762003b6d60418462006155565b8151811062003b805762003b8062005faa565b6001600160a01b039092166020928302919091019091015262003ba560418262005be3565b905062003ae1565b508a602001516001600160a01b03168660008151811062003bd25762003bd262005faa565b60200260200101516001600160a01b03161462003c445760405162461bcd60e51b815260206004820152602960248201527f496e76616c6964206e6f6465206b6579206f776e6572736869702070726f6f66604482015268081c1c9bdd9a59195960ba1b606482015260840162000db7565b8a604001516001600160a01b03168660018151811062003c685762003c6862005faa565b60200260200101516001600160a01b03161462003cdc5760405162461bcd60e51b815260206004820152602b60248201527f496e76616c6964206f7261636c65206b6579206f776e6572736869702070726f60448201526a1bd9881c1c9bdd9a59195960aa1b606482015260840162000db7565b600162003cf48c6102400151838e6000015162004e12565b1462003d625760405162461bcd60e51b815260206004820152603660248201527f496e76616c696420636f6e73656e737573206b6579206f776e65727368697020604482015275383937b7b3103337b9103932b3b4b9ba3930ba34b7b760511b606482015260840162000db7565b62003d6d8b62004e81565b5050505050505050505050565b6000821162003dd85760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000db7565b6001600160a01b03811660009081526027602052604090205482111562003e425760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000db7565b6001600160a01b0381166000908152602760205260408120805484929062003e6c90849062005bcd565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580546000908152600394859052978820875181549088166001600160a01b031991821617825595516001820180549190981696169590951790955590516002840155519101558054919262003ef88362005fc0565b90915550506001600160a01b03848116600081815260286020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d91015b60405180910390a35050505050565b6001600160a01b0380841660009081526028602052604090208054909183811691161480620040dc57600c820154604051631092ab9160e31b81526001600160a01b03858116600483015260009216906384955c8890602401602060405180830381865afa15801562003fdd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200400391906200616c565b9050848110156200406b5760405162461bcd60e51b815260206004820152602b60248201527f696e73756666696369656e7420756e6c6f636b6564204c6971756964204e657760448201526a746f6e2062616c616e636560a81b606482015260840162000db7565b600c83015460405163282d3fdf60e01b81526001600160a01b038681166004830152602482018890529091169063282d3fdf90604401600060405180830381600087803b158015620040bc57600080fd5b505af1158015620040d1573d6000803e3d6000fd5b50505050506200416f565b8382600b01548360080154620040f3919062005bcd565b1015620041535760405162461bcd60e51b815260206004820152602760248201527f696e73756666696369656e742073656c6620626f6e646564206e6577746f6e2060448201526662616c616e636560c81b606482015260840162000db7565b8382600b01600082825462004169919062005be3565b90915550505b6040805160e0810182526001600160a01b0380861682528781166020808401918252838501898152600060608601818152436080880190815260a088018381528a151560c08a019081526008805486526006909752998420985189549089166001600160a01b0319918216178a55965160018a01805491909916971696909617909655915160028701559051600386015592516004850155905160059093018054945115156101000261ff00199415159490941661ffff199095169490941792909217909255805491620042438362005fc0565b9190505550826001600160a01b0316856001600160a01b03167f63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc838760405162003f579291909115158252602082015260400190565b60fa60a06000808285855af46200370d573d6000803e3d6000fd5b81546002600180831615610100020382160482518082016020811060208410016002811462004363576001811462004389578660005260208404602060002001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b808210156200434c57815184556001840193506020820191506200432d565b815191036101000a908190040290915550620043fe565b60028302826020036101000a846020036101000a602089015104020185018755620043fe565b8660005260208404602060002001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b80821015620043ec5781518455600184019350602082019150620043cd565b815191036101000a9081900402909155505b50505050505050565b600081815260036020908152604080832060018101546001600160a01b03168452602890925282209091601382015460ff1660038111156200444d576200444d620053b9565b14620044ed57600282015482546001600160a01b0316600090815260276020526040812080549091906200448390849062005be3565b909155505081546001830154600284015460138401546040517f1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f87834294620044e0946001600160a01b0391821694911692909160ff90911690620061cf565b60405180910390a1505050565b805482546001600160a01b03908116911614620045ea57600080826008015483600501546200451d919062005bcd565b9050806000036200453557836002015491506200455b565b80846002015484600d01546200454c919062005fdc565b62004558919062006155565b91505b600c83015484546040516340c10f1960e01b81526001600160a01b039182166004820152602481018590529116906340c10f1990604401600060405180830381600087803b158015620045ad57600080fd5b505af1158015620045c2573d6000803e3d6000fd5b505050508183600d016000828254620045dc919062005be3565b909155506200460a92505050565b816002015481600801600082825462004604919062005be3565b90915550505b816002015481600501600082825462004624919062005be3565b9091555050505050565b600081815260066020908152604080832060018101546001600160a01b0316845260289092528220600582015491929091610100900460ff1662004822576002830154600c8301548454604051637eee288d60e01b81526001600160a01b03918216600482015260248101849052911690637eee288d90604401600060405180830381600087803b158015620046c357600080fd5b505af1158015620046d8573d6000803e3d6000fd5b50505050600c8301548454604051632770a7eb60e21b81526001600160a01b03918216600482015260248101849052911690639dc29fac90604401600060405180830381600087803b1580156200472e57600080fd5b505af115801562004743573d6000803e3d6000fd5b505050506000836008015484600501546200475f919062005bcd565b600d85015490915062004773828462005fdc565b6200477f919062006155565b92508184600d01600082825462004797919062005bcd565b90915550506006840154600003620047b65760038501839055620047de565b60068401546007850154620047cc908562005fdc565b620047d8919062006155565b60038601555b82846006016000828254620047f4919062005be3565b909155505060038501546007850180546000906200481490849062005be3565b90915550620048f792505050565b50600282015460088201548111156200483c575060088101545b81600901546000036200485657600383018190556200487e565b6009820154600a8301546200486c908362005fdc565b62004878919062006155565b60038401555b8082600901600082825462004894919062005be3565b90915550506003830154600a83018054600090620048b490849062005be3565b9250508190555080826008016000828254620048d1919062005bcd565b90915550506002830154600b83018054600090620048f190849062005bcd565b90915550505b6005808401805460ff19166001179055820180548291906000906200491e90849062005bcd565b909155505050505050565b6000818152600660205260408120600381015490910362004948575050565b60018101546001600160a01b031660009081526028602052604081206005830154909190610100900460ff16620049e45781600701548260060154846003015462004994919062005fdc565b620049a0919062006155565b905080826006016000828254620049b8919062005bcd565b90915550506003830154600783018054600090620049d890849062005bcd565b9091555062004a499050565b81600a015482600901548460030154620049ff919062005fdc565b62004a0b919062006155565b90508082600901600082825462004a23919062005bcd565b90915550506003830154600a8301805460009062004a4390849062005bcd565b90915550505b82546001600160a01b0316600090815260276020526040812080548392906200491e90849062005be3565b600062004a858260600151620050b2565b6001600160a01b0390911660208401529050801562004ad55760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000db7565b6020808301516001600160a01b0390811660009081526028909252604090912060010154161562004b495760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000db7565b61271082608001511115620024525760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000db7565b60608160000362004bc95750506040805180820190915260018152600360fc1b602082015290565b8160005b811562004bf9578062004be08162005fc0565b915062004bf19050600a8362006155565b915062004bcd565b60008167ffffffffffffffff81111562004c175762004c1762005791565b6040519080825280601f01601f19166020018201604052801562004c42576020820181803683370190505b5090505b841562004cba5762004c5a60018362005bcd565b915062004c69600a866200600c565b62004c7690603062005be3565b60f81b81838151811062004c8e5762004c8e62005faa565b60200101906001600160f81b031916908160001a90535062004cb2600a8662006155565b945062004c46565b949350505050565b60608162004cd281601f62005be3565b101562004d135760405162461bcd60e51b815260206004820152600e60248201526d736c6963655f6f766572666c6f7760901b604482015260640162000db7565b62004d1f828462005be3565b8451101562004d655760405162461bcd60e51b8152602060048201526011602482015270736c6963655f6f75744f66426f756e647360781b604482015260640162000db7565b60608215801562004d86576040519150600082526020820160405262004dd2565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101562004dc157805183526020928301920162004da7565b5050858452601f01601f1916604052505b50949350505050565b8181018051602082015160409092015190919060001a601b81101562004e0b5762004e08601b82620061fd565b90505b9250925092565b600062004e1e62005244565b600085858560405160200162004e379392919062006219565b6040516020818303038152906040529050600060fb905060008251602062004e60919062005be3565b90506020848285855afa62004e7457600080fd5b5050905195945050505050565b6101808101516001600160a01b031662004f0457601d5460009062004ea69062004ba1565b90508160200151826000015183608001518360405162004ec69062005262565b62004ed5949392919062006268565b604051809103906000f08015801562004ef2573d6000803e3d6000fd5b506001600160a01b0316610180830152505b60208181018051601d80546001808201835560009283527f6d4407e7be21f808e6509aa9fa9143369579dd7d760fe20a2c09680fc146134f90910180546001600160a01b03199081166001600160a01b039586161790915584518416835260289095526040918290208651815487169085161781559351908401805486169184169190911790558401516002830180549094169116179091556060820151829190600382019062004fb69082620062a7565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620050809082620062a7565b5061026082015160138201805460ff19166001836003811115620050a857620050a8620053b9565b0217905550505050565b600080620050bf62005270565b60008060ff9050604083875160208901845afa620050dc57600080fd5b50508051602090910151600160601b90910494909350915050565b60405180610280016040528060006001600160a01b0316815260200160006001600160a01b0316815260200160006001600160a01b0316815260200160608152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b0316815260200160008152602001600081526020016000815260200160008152602001600081526020016060815260200160006003811115620051c257620051c2620053b9565b905290565b6040518060a001604052806005906020820280368337509192915050565b508054600082559060005260206000209081019062000e0791906200528e565b508054620052139062005b7b565b6000825580601f1062005224575050565b601f01602090049060005260206000209081019062000e079190620052b3565b60405180602001604052806001906020820280368337509192915050565b611531806200636883390190565b60405180604001604052806002906020820280368337509192915050565b80821115620052af576000620052a5828262005205565b506001016200528e565b5090565b5b80821115620052af5760008155600101620052b4565b60005b83811015620052e7578181015183820152602001620052cd565b50506000910152565b600081518084526200530a816020860160208601620052ca565b601f01601f19169290920160200192915050565b602081526000620053336020830184620052f0565b9392505050565b6001600160a01b038116811462000e0757600080fd5b600080604083850312156200536457600080fd5b823562005371816200533a565b946020939093013593505050565b6000602082840312156200539257600080fd5b813562005333816200533a565b600060208284031215620053b257600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b60048110620053ee57634e487b7160e01b600052602160045260246000fd5b9052565b602081526200540d6020820183516001600160a01b03169052565b600060208301516200542a60408401826001600160a01b03169052565b5060408301516001600160a01b03811660608401525060608301516102808060808501526200545e6102a0850183620052f0565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a0620054e6818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f1901610260808801919091529091506200554a8483620052f0565b9350808701519150506200556182860182620053cf565b5090949350505050565b6000806000606084860312156200558157600080fd5b83356200558e816200533a565b92506020840135620055a0816200533a565b929592945050506040919091013590565b600060208284031215620055c457600080fd5b813567ffffffffffffffff811115620055dc57600080fd5b820161028081850312156200533357600080fd5b600081518084526020808501808196508360051b8101915082860160005b8581101562005665578284038952815180516001600160a01b0316855285810151868601526040908101516060918601829052906200565081870183620052f0565b9a87019a95505050908401906001016200560e565b5091979650505050505050565b83151581526060602082015260006200568f6060830185620055f0565b9050826040830152949350505050565b60008060408385031215620056b357600080fd5b50508035926020909101359150565b845181526020808601518183015260408087015181840152606080880151818501526080808901516001600160a01b03908116828701528851811660a08088019190915294890151811660c087015292880151831660e086015290870151821661010085015286015116610120830152840151610200820190620057526101408401826001600160a01b03169052565b5083516001600160a01b0316610160830152602084015161018083015260408401516101a08301526060909301516101c08201526101e0015292915050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620057b957600080fd5b813567ffffffffffffffff80821115620057d757620057d762005791565b604051601f8301601f19908116603f0116810190828211818310171562005802576200580262005791565b816040528381528660208588010111156200581c57600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080608085870312156200585357600080fd5b843567ffffffffffffffff808211156200586c57600080fd5b6200587a88838901620057a7565b9550602087013591506200588e826200533a565b90935060408601359080821115620058a557600080fd5b620058b388838901620057a7565b93506060870135915080821115620058ca57600080fd5b50620058d987828801620057a7565b91505092959194509250565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156200593e57603f198886030184526200592b858351620052f0565b945092850192908501906001016200590c565b5092979650505050505050565b602081526000620053336020830184620055f0565b6020808252825182820181905260009190848201906040850190845b81811015620059a35783516001600160a01b0316835292840192918401916001016200597c565b50909695505050505050565b60008060408385031215620059c357600080fd5b823567ffffffffffffffff80821115620059dc57600080fd5b620059ea86838701620057a7565b9350602085013591508082111562005a0157600080fd5b5062005a1085828601620057a7565b9150509250929050565b60408152600062005a2f6040830185620052f0565b828103602084015262005a438185620052f0565b95945050505050565b6000806040838503121562005a6057600080fd5b823562005a6d816200533a565b9150602083013562005a7f816200533a565b809150509250929050565b6020808252601c908201527f76616c696461746f72206d757374206265207265676973746572656400000000604082015260600190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b60208082526018908201527f76616c696461746f72206e6f7420726567697374657265640000000000000000604082015260600190565b600181811c9082168062005b9057607f821691505b60208210810362005bb157634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000d715762000d7162005bb7565b8082018082111562000d715762000d7162005bb7565b6000813562000d71816200533a565b80546001600160a01b0319166001600160a01b0392909216919091179055565b6000808335601e1984360301811262005c4057600080fd5b83018035915067ffffffffffffffff82111562005c5c57600080fd5b60200191503681900382131562005c7257600080fd5b9250929050565b601f8211156200370d57600081815260208120601f850160051c8101602086101562005ca25750805b601f850160051c820191505b8181101562005cc35782815560010162005cae565b505050505050565b600019600383901b1c191660019190911b1790565b67ffffffffffffffff83111562005cfb5762005cfb62005791565b62005d138362005d0c835462005b7b565b8362005c79565b6000601f84116001811462005d46576000851562005d315750838201355b62005d3d868262005ccb565b8455506200152c565b600083815260209020601f19861690835b8281101562005d79578685013582556020948501946001909201910162005d57565b508682101562005d975760001960f88860031b161c19848701351681555b505060018560011b0183555050505050565b600081356004811062000d7157600080fd5b6004821062005dda57634e487b7160e01b600052602160045260246000fd5b60ff1981541660ff831681178255505050565b62005e0362005dfc8362005bf9565b8262005c08565b62005e1f62005e156020840162005bf9565b6001830162005c08565b62005e3b62005e316040840162005bf9565b6002830162005c08565b62005e4a606083018362005c28565b62005e5a81836003860162005ce0565b50506080820135600482015560a0820135600582015560c0820135600682015560e0820135600782015561010082013560088201556101208201356009820155610140820135600a820155610160820135600b82015562005ecd62005ec3610180840162005bf9565b600c830162005c08565b6101a0820135600d8201556101c0820135600e8201556101e0820135600f8201556102008201356010820155610220820135601182015562005f1461024083018362005c28565b62005f2481836012860162005ce0565b50506200245262005f39610260840162005da9565b6013830162005dbb565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b60006020828403121562005f9957600080fd5b815180151581146200533357600080fd5b634e487b7160e01b600052603260045260246000fd5b60006001820162005fd55762005fd562005bb7565b5060010190565b808202811582820484141762000d715762000d7162005bb7565b634e487b7160e01b600052601260045260246000fd5b6000826200601e576200601e62005ff6565b500690565b600060018060a01b0380881683528087166020840152808616604084015260a060608401526200605760a0840186620052f0565b91508084166080840152509695505050505050565b81810362006078575050565b62006084825462005b7b565b67ffffffffffffffff8111156200609f576200609f62005791565b620060b781620060b0845462005b7b565b8462005c79565b6000601f821160018114620060ea5760008315620060d55750848201545b620060e1848262005ccb565b8555506200152c565b600085815260209020601f19841690600086815260209020845b8381101562006126578286015482556001958601959091019060200162006104565b5085831015620061455781850154600019600388901b60f8161c191681555b5050505050600190811b01905550565b60008262006167576200616762005ff6565b500490565b6000602082840312156200617f57600080fd5b5051919050565b600084516200619a818460208901620052ca565b845190830190620061b0818360208901620052ca565b8451910190620061c5818360208801620052ca565b0195945050505050565b6001600160a01b03858116825284166020820152604081018390526080810162005a436060830184620053cf565b60ff818116838216019081111562000d715762000d7162005bb7565b600084516200622d818460208901620052ca565b84519083019062006243818360208901620052ca565b60609490941b6bffffffffffffffffffffffff19169301928352505060140192915050565b6001600160a01b03858116825284166020820152604081018390526080606082018190526000906200629d90830184620052f0565b9695505050505050565b815167ffffffffffffffff811115620062c457620062c462005791565b620062d581620060b0845462005b7b565b602080601f831160018114620063095760008415620062f45750858301515b62006300858262005ccb565b86555062005cc3565b600085815260208120601f198616915b828110156200633a5788860151825594840194600190910190840162006319565b50858210156200614557939096015160001960f8600387901b161c19169092555050600190811b0190555056fe60806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220018b8c3d2155d9eb67d018108e35d70682e98011384859e6cefa8fd80d2593e364736f6c63430008150033a26469706673582212204f87ee3fa792f13a44901188d11f188ef1496b6e8fca82d35195dc4d5a64b67964736f6c6343000815003360806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220018b8c3d2155d9eb67d018108e35d70682e98011384859e6cefa8fd80d2593e364736f6c63430008150033",
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
// Solidity: function computeCommittee() returns(address[])
func (_Autonity *Autonity) ComputeCommittee(opts *runOptions) (uint64, error) {
	_, consumed, err := _Autonity.call(opts, "computeCommittee")
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256,bytes)[], uint256)
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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"ActivatedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegatee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"BondingRejected\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurnedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"CommissionRateChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"period\",\"type\":\"uint256\"}],\"name\":\"EpochPeriodUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"}],\"name\":\"MinimumBaseFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MintedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"selfBonded\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewUnbondingRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"effectiveBlock\",\"type\":\"uint256\"}],\"name\":\"PausedValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"treasury\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"liquidContract\",\"type\":\"address\"}],\"name\":\"RegisteredValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Rewarded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"COMMISSION_RATE_PRECISION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"activateValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rate\",\"type\":\"uint256\"}],\"name\":\"changeCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"completeContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"computeCommittee\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"config\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"treasuryFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingPeriod\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"treasuryAccount\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Policy\",\"name\":\"policy\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"contractIAccountability\",\"name\":\"accountabilityContract\",\"type\":\"address\"},{\"internalType\":\"contractIOracle\",\"name\":\"oracleContract\",\"type\":\"address\"},{\"internalType\":\"contractIACU\",\"name\":\"acuContract\",\"type\":\"address\"},{\"internalType\":\"contractISupplyControl\",\"name\":\"supplyControlContract\",\"type\":\"address\"},{\"internalType\":\"contractIStabilization\",\"name\":\"stabilizationContract\",\"type\":\"address\"},{\"internalType\":\"contractUpgradeManager\",\"name\":\"upgradeManagerContract\",\"type\":\"address\"}],\"internalType\":\"structAutonity.Contracts\",\"name\":\"contracts\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"operatorAccount\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"epochPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockPeriod\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"committeeSize\",\"type\":\"uint256\"}],\"internalType\":\"structAutonity.Protocol\",\"name\":\"protocol\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"contractVersion\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochTotalBondedStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeInitialization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommittee\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"votingPower\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"}],\"internalType\":\"structAutonity.CommitteeMember[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCommitteeEnodes\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"}],\"name\":\"getEpochFromBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getEpochPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCommitteeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumBaseFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewContract\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAccount\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUnbondingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getValidator\",\"outputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastEpochBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"pauseValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_enode\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_oracleAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_signatures\",\"type\":\"bytes\"}],\"name\":\"registerValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resetContractUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIAccountability\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAccountabilityContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIACU\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setAcuContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"name\":\"setCommitteeSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setEpochPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"setMinimumBaseFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperatorAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setOracleContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIStabilization\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setStabilizationContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISupplyControl\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSupplyControlContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setTreasuryAccount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_treasuryFee\",\"type\":\"uint256\"}],\"name\":\"setTreasuryFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setUnbondingPeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractUpgradeManager\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setUpgradeManagerContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalRedistributed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"unbond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"treasury\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nodeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"enode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfBondedStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStake\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingShares\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"selfUnbondingStakeLocked\",\"type\":\"uint256\"},{\"internalType\":\"contractLiquid\",\"name\":\"liquidContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSlashed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"jailReleaseBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"provableFaultCount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"consensusKey\",\"type\":\"bytes\"},{\"internalType\":\"enumValidatorState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAutonity.Validator\",\"name\":\"_val\",\"type\":\"tuple\"}],\"name\":\"updateValidatorAndTransferSlashedFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_bytecode\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_abi\",\"type\":\"string\"}],\"name\":\"upgradeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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
	Bin: "0x60806040526000600b556000600c553480156200001b57600080fd5b50604080516000808252602082019092529062000051565b6200003d62000d4e565b815260200190600190039081620000335790505b506040805161012081018252600d546080808301918252600e5460a080850191909152600f5460c08086019190915260105460e08601526011546001600160a01b039081166101008701529385528551908101865260125484168152601354841660208281019190915260145485168288015260155485166060808401919091526016548616838601526017548616938301939093528086019190915285519283018652601854909316825260195492820192909252601a5481850152601b548183015292820192909252601c54918101829052906000036200014c57602a80546001600160a01b031916331790556200014c828262000169565b5050601c546001036200016357620001636200051f565b6200115e565b80518051600d55602080820151600e55604080830151600f55606080840151601055608093840151601180546001600160a01b03199081166001600160a01b03938416179091558487015180516012805484169185169190911790558086015160138054841691851691909117905580850151601480548416918516919091179055808401516015805484169185169190911790559586015160168054831691841691909117905560a0909501516017805487169183169190911790558286015180516018805490971692169190911790945591830151601955820151601a5590810151601b55810151601c5560005b82518110156200051a5760008382815181106200027a576200027a62000ebd565b602002602001015160a00151905060008483815181106200029f576200029f62000ebd565b60200260200101516101a00181815250506000848381518110620002c757620002c762000ebd565b602002602001015161018001906001600160a01b031690816001600160a01b031681525050600084838151811062000303576200030362000ebd565b602002602001015160a001818152505060008483815181106200032a576200032a62000ebd565b60209081029190910101516101c00152600f54845185908490811062000354576200035462000ebd565b6020026020010151608001818152505060008483815181106200037b576200037b62000ebd565b6020026020010151610260019060038111156200039c576200039c62000ed3565b90816003811115620003b257620003b262000ed3565b815250506000848381518110620003cd57620003cd62000ebd565b60200260200101516101600181815250506200040b848381518110620003f757620003f762000ebd565b60200260200101516200068d60201b60201c565b6200043884838151811062000424576200042462000ebd565b6020026020010151620007c460201b60201c565b806027600086858151811062000452576200045262000ebd565b6020026020010151600001516001600160a01b03166001600160a01b0316815260200190815260200160002060008282546200048f919062000eff565b925050819055508060296000828254620004aa919062000eff565b9250508190555062000504848381518110620004ca57620004ca62000ebd565b60200260200101516020015182868581518110620004ec57620004ec62000ebd565b602002602001015160000151620009f560201b60201c565b5080620005118162000f1b565b91505062000259565b505050565b602a546001600160a01b031633146200058b5760405162461bcd60e51b815260206004820152602360248201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60448201526218dbdb60ea1b60648201526084015b60405180910390fd5b600260286000601d600181548110620005a857620005a862000ebd565b60009182526020808320909101546001600160a01b0316835282019290925260400181206005018054909190620005e190849062000f4d565b92505081905550600260286000601d60018154811062000605576200060562000ebd565b60009182526020808320909101546001600160a01b03168352820192909252604001812060080180549091906200063e90849062000f4d565b90915550506002601c556018546001600160a01b031660009081526027602052604081206103e8905562000673908062000e1e565b620006816001600062000e1e565b6002805460ff19169055565b6000620006a4826060015162000be060201b60201c565b6001600160a01b03909116602084015290508015620006f45760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000582565b6020808301516001600160a01b03908116600090815260289092526040909120600101541615620007685760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000582565b61271082608001511115620007c05760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000582565b5050565b6101808101516001600160a01b03166200084757601d54600090620007e99062000c2e565b905081602001518260000151836080015183604051620008099062000e60565b62000818949392919062000f64565b604051809103906000f08015801562000835573d6000803e3d6000fd5b506001600160a01b0316610180830152505b60208181018051601d80546001808201835560009283527f6d4407e7be21f808e6509aa9fa9143369579dd7d760fe20a2c09680fc146134f90910180546001600160a01b03199081166001600160a01b0395861617909155845184168352602890955260409182902086518154871690851617815593519084018054861691841691909117905584015160028301805490941691161790915560608201518291906003820190620008f9908262001065565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620009c3908262001065565b5061026082015160138201805460ff19166001836003811115620009eb57620009eb62000ed3565b0217905550505050565b6000821162000a535760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000582565b6001600160a01b03811660009081526027602052604090205482111562000abd5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000582565b6001600160a01b0381166000908152602760205260408120805484929062000ae790849062001131565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580546000908152600394859052978820875181549088166001600160a01b031991821617825595516001820180549190981696169590951790955590516002840155519101558054919262000b738362000f1b565b90915550506001600160a01b03848116600081815260286020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d910160405180910390a35050505050565b60008062000bed62000e6e565b60008060ff9050604083875160208901845afa62000c0a57600080fd5b505080516020909101516c0100000000000000000000000090910494909350915050565b60608160000362000c565750506040805180820190915260018152600360fc1b602082015290565b8160005b811562000c86578062000c6d8162000f1b565b915062000c7e9050600a8362000f4d565b915062000c5a565b6000816001600160401b0381111562000ca35762000ca362000ea7565b6040519080825280601f01601f19166020018201604052801562000cce576020820181803683370190505b5090505b841562000d465762000ce660018362001131565b915062000cf5600a8662001147565b62000d0290603062000eff565b60f81b81838151811062000d1a5762000d1a62000ebd565b60200101906001600160f81b031916908160001a90535062000d3e600a8662000f4d565b945062000cd2565b949350505050565b60405180610280016040528060006001600160a01b0316815260200160006001600160a01b0316815260200160006001600160a01b0316815260200160608152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b031681526020016000815260200160008152602001600081526020016000815260200160008152602001606081526020016000600381111562000e195762000e1962000ed3565b905290565b50805462000e2c9062000fd7565b6000825580601f1062000e3d575050565b601f01602090049060005260206000209081019062000e5d919062000e8c565b50565b6115318062008a6583390190565b60405180604001604052806002906020820280368337509192915050565b5b8082111562000ea3576000815560010162000e8d565b5090565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052601160045260246000fd5b8082018082111562000f155762000f1562000ee9565b92915050565b60006001820162000f305762000f3062000ee9565b5060010190565b634e487b7160e01b600052601260045260246000fd5b60008262000f5f5762000f5f62000f37565b500490565b600060018060a01b038087168352602081871681850152856040850152608060608501528451915081608085015260005b8281101562000fb35785810182015185820160a00152810162000f95565b5050600060a0828501015260a0601f19601f83011684010191505095945050505050565b600181811c9082168062000fec57607f821691505b6020821081036200100d57634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156200051a57600081815260208120601f850160051c810160208610156200103c5750805b601f850160051c820191505b818110156200105d5782815560010162001048565b505050505050565b81516001600160401b0381111562001081576200108162000ea7565b620010998162001092845462000fd7565b8462001013565b602080601f831160018114620010d15760008415620010b85750858301515b600019600386901b1c1916600185901b1785556200105d565b600085815260208120601f198616915b828110156200110257888601518255948401946001909101908401620010e1565b5085821015620011215787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b8181038181111562000f155762000f1562000ee9565b60008262001159576200115962000f37565b500690565b6178f7806200116e6000396000f3fe608060405260043610620003db5760003560e01c8063872cf0591162000203578063b66b3e791162000117578063d372c07e11620000a7578063dd62ed3e1162000075578063dd62ed3e1462000cbd578063dfb1a4d21462000d07578063e7f43c681462000d1e578063f7866ee31462000d3e57005b8063d372c07e1462000c39578063d5f394881462000c5e578063d861b0e81462000c80578063d886f8a21462000c9857005b8063cb696f5411620000e5578063cb696f541462000bb2578063ceaad4551462000bd7578063cf9c57191462000bfc578063cfd19fb91462000c1457005b8063b66b3e791462000b42578063b7ab4db51462000b6a578063c2362dd51462000b82578063c9d97af41462000b9a57005b8063a5d059ca1162000193578063ae1f5fa01162000161578063ae1f5fa01462000aac578063b2ea9adb1462000ad3578063b3ecbadd1462000af8578063b46e55201462000b1d57005b8063a5d059ca1462000a14578063a8b2216e1462000a39578063a9059cbb1462000a60578063ab8f6ffe1462000a8557005b80639bb851c011620001d15780639bb851c0146200099a5780639c98e47114620009b25780639dc29fac14620009ca578063a515366a14620009ef57005b8063872cf05914620008fe5780638bac7dad146200091657806395d89b41146200093b57806396b477cb146200096957005b806340c10f1911620002fb57806370a08231116200028b578063819b64631162000259578063819b6463146200087d578063833b1fce146200089457806384467fdb14620008b4578063852c484914620008d957005b806370a082311462000729578063731b3a03146200076357806377e741c7146200077a57806379502c55146200079f57005b8063520fdbbc11620002c9578063520fdbbc146200068a5780635f7d394914620006af5780636b5f444c14620006ed5780636fd2c80b146200071257005b806340c10f191462000600578063436459691462000625578063496ccd9b146200063c5780634bb278f3146200066157005b80631604e416116200037757806329070c6d116200034557806329070c6d146200058e5780632f2c3f2e14620005a5578063313ce56714620005bd57806335be16e014620005db57005b80631604e416146200050657806318160ddd146200051e5780631904bb2e146200053557806323b872dd146200056957005b80630d8e6e2c11620003b55780630d8e6e2c14620004845780631122063314620004a5578063114eaf5514620004bc5780631250a28d14620004e157005b806306fdde0314620003e5578063095ea7b314620004295780630ae65e7a146200045f57005b36620003e357005b005b348015620003f257600080fd5b506040805180820190915260068152652732bbba37b760d11b60208201525b60405162000420919062005347565b60405180910390f35b3480156200043657600080fd5b506200044e6200044836600462005379565b62000d5e565b604051901515815260200162000420565b3480156200046c57600080fd5b50620003e36200047e366004620053a8565b62000d77565b3480156200049157600080fd5b50601c545b60405190815260200162000420565b348015620004b257600080fd5b50600e5462000496565b348015620004c957600080fd5b50620003e3620004db366004620053c8565b62000e0a565b348015620004ee57600080fd5b50620003e362000500366004620053a8565b62000e3c565b3480156200051357600080fd5b506200049660245481565b3480156200052b57600080fd5b5060295462000496565b3480156200054257600080fd5b506200055a62000554366004620053a8565b62000e8b565b6040516200042091906200541b565b3480156200057657600080fd5b506200044e6200058836600462005594565b62001139565b3480156200059b57600080fd5b50600d5462000496565b348015620005b257600080fd5b506200049661271081565b348015620005ca57600080fd5b506040516012815260200162000420565b348015620005e857600080fd5b50620003e3620005fa366004620055da565b62001193565b3480156200060d57600080fd5b50620003e36200061f36600462005379565b6200136f565b3480156200063257600080fd5b50601a5462000496565b3480156200064957600080fd5b50620003e36200065b366004620053a8565b62001429565b3480156200066e57600080fd5b506200067962001533565b60405162000420939291906200569b565b3480156200069757600080fd5b50620003e3620006a9366004620053a8565b62001923565b348015620006bc57600080fd5b50620006d4620006ce366004620056c8565b62001b1e565b6040516001600160a01b03909116815260200162000420565b348015620006fa57600080fd5b50620003e36200070c366004620053c8565b62001d35565b3480156200071f57600080fd5b5060105462000496565b3480156200073657600080fd5b506200049662000748366004620053a8565b6001600160a01b031660009081526027602052604090205490565b3480156200077057600080fd5b5060205462000496565b3480156200078757600080fd5b50620003e362000799366004620053c8565b62001ebb565b348015620007ac57600080fd5b506040805160a08082018352600d548252600e54602080840191909152600f54838501526010546060808501919091526011546001600160a01b03908116608080870191909152865160c0810188526012548316815260135483168186015260145483168189015260155483168185015260165483168183015260175483169581019590955286519081018752601854909116815260195492810192909252601a5494820194909452601b5493810193909352601c546200086b939084565b604051620004209493929190620056eb565b3480156200088a57600080fd5b50601b5462000496565b348015620008a157600080fd5b506013546001600160a01b0316620006d4565b348015620008c157600080fd5b50620003e3620008d336600462005865565b62001eed565b348015620008e657600080fd5b50620003e3620008f836600462005379565b6200201c565b3480156200090b57600080fd5b50620003e3620021ad565b3480156200092357600080fd5b50620003e362000935366004620053c8565b620021e9565b3480156200094857600080fd5b50604080518082019091526003815262272a2760e91b602082015262000411565b3480156200097657600080fd5b506200049662000988366004620053c8565b6000908152601f602052604090205490565b348015620009a757600080fd5b506200049660235481565b348015620009bf57600080fd5b506200049660215481565b348015620009d757600080fd5b50620003e3620009e936600462005379565b6200226d565b348015620009fc57600080fd5b50620003e362000a0e36600462005379565b62002383565b34801562000a2157600080fd5b50620003e362000a3336600462005379565b62002456565b34801562000a4657600080fd5b5062000a51620024ed565b6040516200042091906200590e565b34801562000a6d57600080fd5b506200044e62000a7f36600462005379565b620025d0565b34801562000a9257600080fd5b5062000a9d620025df565b60405162000420919062005974565b34801562000ab957600080fd5b5062000ac4620026f1565b60405162000420919062005989565b34801562000ae057600080fd5b50620003e362000af2366004620059d8565b62002937565b34801562000b0557600080fd5b50620003e362000b17366004620053a8565b6200297e565b34801562000b2a57600080fd5b50620003e362000b3c366004620053a8565b620029cd565b34801562000b4f57600080fd5b5062000b5a62002c50565b6040516200042092919062005a43565b34801562000b7757600080fd5b5062000ac462002d87565b34801562000b8f57600080fd5b506200049660205481565b34801562000ba757600080fd5b5062000496601e5481565b34801562000bbf57600080fd5b50620003e362000bd1366004620053c8565b62002deb565b34801562000be457600080fd5b50620003e362000bf6366004620053a8565b62002e4e565b34801562000c0957600080fd5b50620003e362002e9d565b34801562000c2157600080fd5b50620003e362000c33366004620053a8565b62002ef1565b34801562000c4657600080fd5b50620003e362000c58366004620053a8565b62002f40565b34801562000c6b57600080fd5b50602a54620006d4906001600160a01b031681565b34801562000c8d57600080fd5b50620003e362002f8f565b34801562000ca557600080fd5b50620003e362000cb7366004620053a8565b62002fd0565b34801562000cca57600080fd5b506200049662000cdc36600462005a75565b6001600160a01b03918216600090815260266020908152604080832093909416825291909152205490565b34801562000d1457600080fd5b5060195462000496565b34801562000d2b57600080fd5b506018546001600160a01b0316620006d4565b34801562000d4b57600080fd5b506011546001600160a01b0316620006d4565b600062000d6d3384846200301f565b5060015b92915050565b6001600160a01b038082166000818152602860205260409020600101549091161462000dc05760405162461bcd60e51b815260040162000db79062005ab3565b60405180910390fd5b6001600160a01b0381811660009081526028602052604090205416331462000dfc5760405162461bcd60e51b815260040162000db79062005aea565b62000e078162003148565b50565b6018546001600160a01b0316331462000e375760405162461bcd60e51b815260040162000db79062005b36565b601055565b6018546001600160a01b0316331462000e695760405162461bcd60e51b815260040162000db79062005b36565b601280546001600160a01b0319166001600160a01b0392909216919091179055565b62000e9562005120565b6001600160a01b038083166000818152602860205260409020600101549091161462000ed55760405162461bcd60e51b815260040162000db79062005b6d565b6001600160a01b03808316600090815260286020908152604091829020825161028081018452815485168152600182015485169281019290925260028101549093169181019190915260038201805491929160608401919062000f389062005ba4565b80601f016020809104026020016040519081016040528092919081815260200182805462000f669062005ba4565b801562000fb75780601f1062000f8b5761010080835404028352916020019162000fb7565b820191906000526020600020905b81548152906001019060200180831162000f9957829003601f168201915b505050918352505060048201546020820152600582015460408201526006820154606082015260078201546080820152600882015460a0820152600982015460c0820152600a82015460e0820152600b820154610100820152600c8201546001600160a01b0316610120820152600d820154610140820152600e820154610160820152600f82015461018082015260108201546101a082015260118201546101c08201526012820180546101e090920191620010739062005ba4565b80601f0160208091040260200160405190810160405280929190818152602001828054620010a19062005ba4565b8015620010f25780601f10620010c657610100808354040283529160200191620010f2565b820191906000526020600020905b815481529060010190602001808311620010d457829003601f168201915b5050509183525050601382015460209091019060ff1660038111156200111c576200111c620053e2565b6003811115620011305762001130620053e2565b90525092915050565b6000620011488484846200321f565b6001600160a01b03841660009081526026602090815260408083203384529091528120546200117990849062005bf6565b9050620011888533836200301f565b506001949350505050565b6012546001600160a01b03163314620011fb5760405162461bcd60e51b815260206004820152602360248201527f63616c6c6572206973206e6f742074686520736c617368696e6720636f6e74726044820152621858dd60ea1b606482015260840162000db7565b6000610120820135602882620012186040860160208701620053a8565b6001600160a01b03166001600160a01b031681526020019081526020016000206009015462001248919062005bf6565b60c083013560286000620012636040870160208801620053a8565b6001600160a01b03166001600160a01b031681526020019081526020016000206006015462001293919062005bf6565b60a084013560286000620012ae6040880160208901620053a8565b6001600160a01b03166001600160a01b0316815260200190815260200160002060050154620012de919062005bf6565b620012ea919062005c0c565b620012f6919062005c0c565b6011546001600160a01b03166000908152602760205260408120805492935083929091906200132790849062005c0c565b9091555082905060286000620013446040840160208501620053a8565b6001600160a01b03168152602081019190915260400160002062001369828262005e16565b50505050565b6018546001600160a01b031633146200139c5760405162461bcd60e51b815260040162000db79062005b36565b6001600160a01b03821660009081526027602052604081208054839290620013c690849062005c0c565b925050819055508060296000828254620013e1919062005c0c565b90915550506040518181526001600160a01b038316907f48490b4407bb949b708ec5f514b4167f08f4969baaf78d53b05028adf369bfcf906020015b60405180910390a25050565b6018546001600160a01b03163314620014565760405162461bcd60e51b815260040162000db79062005b36565b601380546001600160a01b0319166001600160a01b03838116918217909255601454604051637adbf97360e01b8152600481019290925290911690637adbf97390602401600060405180830381600087803b158015620014b557600080fd5b505af1158015620014ca573d6000803e3d6000fd5b5050601654604051637adbf97360e01b81526001600160a01b0385811660048301529091169250637adbf97391506024015b600060405180830381600087803b1580156200151757600080fd5b505af11580156200152c573d6000803e3d6000fd5b5050505050565b602a5460009060609082906001600160a01b03163314620015685760405162461bcd60e51b815260040162000db79062005f6c565b601e54436000818152601f6020908152604082209390935560195492549092620015929162005c0c565b6012546040516306c9789b60e41b8152929091146004830181905292506001600160a01b031690636c9789b090602401600060405180830381600087803b158015620015dd57600080fd5b505af1158015620015f2573d6000803e3d6000fd5b505050508015620016e7576200160762003351565b620016116200373b565b6200161b6200382a565b600062001627620026f1565b60135460405163422811f960e11b81529192506001600160a01b03169063845023f2906200165a90849060040162005989565b600060405180830381600087803b1580156200167557600080fd5b505af11580156200168a573d6000803e3d6000fd5b50505050436020819055506001601e6000828254620016aa919062005c0c565b9091555050601e546040519081527febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e3359060200160405180910390a1505b60135460408051634bb278f360e01b815290516000926001600160a01b031691634bb278f3916004808301926020929190829003018187875af115801562001733573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001759919062005faf565b90508015620017d5576014546040805163a2e6204560e01b815290516001600160a01b039092169163a2e620459160048082019260209290919082900301816000875af1925050508015620017cd575060408051601f3d908101601f19168201909252620017ca9181019062005faf565b60015b15620017d557505b81156200190a576002546020805460228054604080518286028101860190915281815260ff90951694919384919060009084015b82821015620018f7576000848152602090819020604080516060810182526003860290920180546001600160a01b03168352600181015493830193909352600283018054929392918401916200185f9062005ba4565b80601f01602080910402602001604051908101604052809291908181526020018280546200188d9062005ba4565b8015620018de5780601f10620018b257610100808354040283529160200191620018de565b820191906000526020600020905b815481529060010190602001808311620018c057829003601f168201915b5050505050815250508152602001906001019062001809565b5050505091509450945094505050909192565b505060025460205460ff90911694606094509092509050565b6018546001600160a01b03163314620019505760405162461bcd60e51b815260040162000db79062005b36565b601880546001600160a01b0319166001600160a01b0383811691821790925560135460405163b3ab15fb60e01b815260048101929092529091169063b3ab15fb90602401600060405180830381600087803b158015620019af57600080fd5b505af1158015620019c4573d6000803e3d6000fd5b505060145460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001a1057600080fd5b505af115801562001a25573d6000803e3d6000fd5b505060155460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001a7157600080fd5b505af115801562001a86573d6000803e3d6000fd5b505060165460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401600060405180830381600087803b15801562001ad257600080fd5b505af115801562001ae7573d6000803e3d6000fd5b505060175460405163b3ab15fb60e01b81526001600160a01b038581166004830152909116925063b3ab15fb9150602401620014fc565b600080805b60225481101562001b7a576022818154811062001b445762001b4462005fd3565b9060005260206000209060030201600101548262001b63919062005c0c565b91508062001b718162005fe9565b91505062001b23565b508060000362001bcd5760405162461bcd60e51b815260206004820152601c60248201527f54686520636f6d6d6974746565206973206e6f74207374616b696e6700000000604482015260640162000db7565b60008362001bdd60638762006005565b62001be9919062005c0c565b905060008160405160200162001c0191815260200190565b60408051601f1981840301815291905280516020909101209050600062001c29848362006035565b90506000805b60225481101562001cd9576022818154811062001c505762001c5062005fd3565b9060005260206000209060030201600101548262001c6f919062005c0c565b915062001c7e60018362005bf6565b831162001cc4576022818154811062001c9b5762001c9b62005fd3565b60009182526020909120600390910201546001600160a01b0316965062000d7195505050505050565b8062001cd08162005fe9565b91505062001c2f565b5060405162461bcd60e51b815260206004820152602960248201527f5468657265206973206e6f2076616c696461746f72206c65667420696e20746860448201526865206e6574776f726b60b81b606482015260840162000db7565b6018546001600160a01b0316331462001d625760405162461bcd60e51b815260040162000db79062005b36565b60195481101562001e19578060205462001d7d919062005c0c565b431062001e195760405162461bcd60e51b815260206004820152605760248201527f63757272656e7420636861696e2068656164206578636565642074686520776960448201527f6e646f773a206c617374426c6f636b45706f6368202b205f6e6577506572696f60648201527f642c2074727920616761696e206c6174746572206f6e2e000000000000000000608482015260a40162000db7565b6019819055601254604051631ad7d11360e21b8152600481018390526001600160a01b0390911690636b5f444c90602401600060405180830381600087803b15801562001e6557600080fd5b505af115801562001e7a573d6000803e3d6000fd5b505050507fd7f1279ded354dbf22a69fcc2fd661763a6e2956a5d2891af9410af880fa5f818160405162001eb091815260200190565b60405180910390a150565b6018546001600160a01b0316331462001ee85760405162461bcd60e51b815260040162000db79062005b36565b600d55565b6000604051806102800160405280336001600160a01b0316815260200160006001600160a01b03168152602001856001600160a01b03168152602001868152602001600d6000016002015481526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b03168152602001600081526020014381526020016000815260200160008152602001600081526020018481526020016000600381111562001fbc5762001fbc620053e2565b9052905062001fcc818362003949565b60208101516101808201516040517f8ad8bd2eb6950e5f332fd3a6dca48cb358ecfe3057848902b98cbdfe455c915c926200200d92339289918b916200604c565b60405180910390a15050505050565b6001600160a01b03808316600081815260286020526040902060010154909116146200205c5760405162461bcd60e51b815260040162000db79062005ab3565b6001600160a01b03828116600090815260286020526040902054163314620020985760405162461bcd60e51b815260040162000db79062005aea565b612710811115620020ec5760405162461bcd60e51b815260206004820152601f60248201527f7265717569726520636f727265637420636f6d6d697373696f6e207261746500604482015260640162000db7565b604080516060810182526001600160a01b038481168252436020808401918252838501868152600c80546000908152600a909352958220855181546001600160a01b0319169516949094178455915160018085019190915591516002909301929092558354929390929091906200216590849062005c0c565b90915550506040518281526001600160a01b038416907f4fba51c92fa3d6ad8374d394f6cd5766857552e153d7384a8f23aa4ce9a8a7cf9060200160405180910390a2505050565b6018546001600160a01b03163314620021da5760405162461bcd60e51b815260040162000db79062005b36565b6002805460ff19166001179055565b6018546001600160a01b03163314620022165760405162461bcd60e51b815260040162000db79062005b36565b60008111620022685760405162461bcd60e51b815260206004820152601960248201527f636f6d6d69747465652073697a652063616e2774206265203000000000000000604482015260640162000db7565b601b55565b6018546001600160a01b031633146200229a5760405162461bcd60e51b815260040162000db79062005b36565b6001600160a01b038216600090815260276020526040902054811115620022fd5760405162461bcd60e51b8152602060048201526016602482015275416d6f756e7420657863656564732062616c616e636560501b604482015260640162000db7565b6001600160a01b038216600090815260276020526040812080548392906200232790849062005bf6565b92505081905550806029600082825462002342919062005bf6565b90915550506040518181526001600160a01b038316907f5024dbeedf0c06664c9bd7be836915730c955e936972c020683dadf11d5488a3906020016200141d565b6001600160a01b0380831660008181526028602052604090206001015490911614620023c35760405162461bcd60e51b815260040162000db79062005b6d565b6001600160a01b03821660009081526028602052604081206013015460ff166003811115620023f657620023f6620053e2565b14620024455760405162461bcd60e51b815260206004820152601b60248201527f76616c696461746f72206e65656420746f206265206163746976650000000000604482015260640162000db7565b6200245282823362003da3565b5050565b6001600160a01b0380831660008181526028602052604090206001015490911614620024965760405162461bcd60e51b815260040162000db79062005b6d565b60008111620024e05760405162461bcd60e51b81526020600482015260156024820152740756e626f6e64696e6720616d6f756e74206973203605c1b604482015260640162000db7565b6200245282823362003f8f565b60606025805480602002602001604051908101604052809291908181526020016000905b82821015620025c7578382906000526020600020018054620025339062005ba4565b80601f0160208091040260200160405190810160405280929190818152602001828054620025619062005ba4565b8015620025b25780601f106200258657610100808354040283529160200191620025b2565b820191906000526020600020905b8154815290600101906020018083116200259457829003601f168201915b50505050508152602001906001019062002511565b50505050905090565b600062000d6d3384846200321f565b60606022805480602002602001604051908101604052809291908181526020016000905b82821015620025c7576000848152602090819020604080516060810182526003860290920180546001600160a01b0316835260018101549383019390935260028301805492939291840191620026599062005ba4565b80601f0160208091040260200160405190810160405280929190818152602001828054620026879062005ba4565b8015620026d85780601f10620026ac57610100808354040283529160200191620026d8565b820191906000526020600020905b815481529060010190602001808311620026ba57829003601f168201915b5050505050815250508152602001906001019062002603565b602a546060906001600160a01b03163314620027215760405162461bcd60e51b815260040162000db79062005f6c565b601d54620027725760405162461bcd60e51b815260206004820152601860248201527f5468657265206d7573742062652076616c696461746f72730000000000000000604482015260640162000db7565b6200277c620051f0565b601b546080820152601d8152602860208201526022604082015260216060820152620027a881620042c2565b620027b6602560006200520e565b60225480620027fd5760405162461bcd60e51b8152602060048201526012602482015271636f6d6d697474656520697320656d70747960701b604482015260640162000db7565b60008167ffffffffffffffff8111156200281b576200281b620057ba565b60405190808252806020026020018201604052801562002845578160200160208202803683370190505b50905060005b828110156200292f576000602860006022848154811062002870576200287062005fd3565b60009182526020808320600392830201546001600160a01b031684528301939093526040909101812060258054600181018255925292507f401968ff42a154441da5f6c4c935ac46b8671f0e062baaa62a7545ba53bb6e4c0190620028d89083018262006095565b50600281015483516001600160a01b039091169084908490811062002901576200290162005fd3565b6001600160a01b03909216602092830291909101909101525080620029268162005fe9565b9150506200284b565b509250505090565b6018546001600160a01b03163314620029645760405162461bcd60e51b815260040162000db79062005b36565b62002971600083620042dd565b62002452600182620042dd565b6018546001600160a01b03163314620029ab5760405162461bcd60e51b815260040162000db79062005b36565b601580546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038082166000818152602860205260409020600101549091161462002a0d5760405162461bcd60e51b815260040162000db79062005ab3565b6001600160a01b0380821660009081526028602052604090208054909116331462002a935760405162461bcd60e51b815260206004820152602f60248201527f726571756972652063616c6c657220746f2062652076616c696461746f72207460448201526e1c99585cdd5c9e481858d8dbdd5b9d608a1b606482015260840162000db7565b6000601382015460ff16600381111562002ab15762002ab1620053e2565b0362002b005760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f7220616c7265616479206163746976650000000000000000604482015260640162000db7565b6002601382015460ff16600381111562002b1e5762002b1e620053e2565b14801562002b2f5750438160100154115b1562002b7e5760405162461bcd60e51b815260206004820152601760248201527f76616c696461746f72207374696c6c20696e206a61696c000000000000000000604482015260640162000db7565b6003601382015460ff16600381111562002b9c5762002b9c620053e2565b0362002beb5760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f72206a61696c6564207065726d616e656e746c7900000000604482015260640162000db7565b60138101805460ff1916905580546019546020546001600160a01b038581169316917f60fcbf2d07dc712a93e59fb28f1edb626d7c2497c57ba71a8c0b3999ecb9a3b59162002c3b919062005c0c565b60405190815260200160405180910390a35050565b6060806000600181805462002c659062005ba4565b80601f016020809104026020016040519081016040528092919081815260200182805462002c939062005ba4565b801562002ce45780601f1062002cb85761010080835404028352916020019162002ce4565b820191906000526020600020905b81548152906001019060200180831162002cc657829003601f168201915b5050505050915080805462002cf99062005ba4565b80601f016020809104026020016040519081016040528092919081815260200182805462002d279062005ba4565b801562002d785780601f1062002d4c5761010080835404028352916020019162002d78565b820191906000526020600020905b81548152906001019060200180831162002d5a57829003601f168201915b50505050509050915091509091565b6060601d80548060200260200160405190810160405280929190818152602001828054801562002de157602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831162002dc2575b5050505050905090565b6018546001600160a01b0316331462002e185760405162461bcd60e51b815260040162000db79062005b36565b600e8190556040518181527f1f4d2fc7529047a5bd96d3229bfea127fd18b7748f13586e097c69fccd3891289060200162001eb0565b6018546001600160a01b0316331462002e7b5760405162461bcd60e51b815260040162000db79062005b36565b601780546001600160a01b0319166001600160a01b0392909216919091179055565b6018546001600160a01b0316331462002eca5760405162461bcd60e51b815260040162000db79062005b36565b62002ed76000806200522e565b62002ee5600160006200522e565b6002805460ff19169055565b6018546001600160a01b0316331462002f1e5760405162461bcd60e51b815260040162000db79062005b36565b601680546001600160a01b0319166001600160a01b0392909216919091179055565b6018546001600160a01b0316331462002f6d5760405162461bcd60e51b815260040162000db79062005b36565b601480546001600160a01b0319166001600160a01b0392909216919091179055565b602a546001600160a01b0316331462002fbc5760405162461bcd60e51b815260040162000db79062005f6c565b62002fc66200373b565b62000e07620026f1565b6018546001600160a01b0316331462002ffd5760405162461bcd60e51b815260040162000db79062005b36565b601180546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b038316620030835760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840162000db7565b6001600160a01b038216620030e65760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840162000db7565b6001600160a01b0383811660008181526026602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b6001600160a01b038116600090815260286020526040812090601382015460ff1660038111156200317d576200317d620053e2565b14620031cc5760405162461bcd60e51b815260206004820152601860248201527f76616c696461746f72206d757374206265206163746976650000000000000000604482015260640162000db7565b60138101805460ff1916600117905580546019546020546001600160a01b038581169316917f75bdcdbe540758778e669d108fbcb7ede734f27f46e4e5525eeb8ecf91849a9c9162002c3b919062005c0c565b6001600160a01b038316600090815260276020526040902054811115620032825760405162461bcd60e51b8152602060048201526016602482015275616d6f756e7420657863656564732062616c616e636560501b604482015260640162000db7565b6001600160a01b03831660009081526027602052604081208054839290620032ac90849062005bf6565b925050819055508060296000828254620032c7919062005c0c565b90915550620032da905081600262006005565b6001600160a01b038316600090815260276020526040812080549091906200330490849062005c0c565b92505081905550816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516200313b91815260200190565b476000036200335c57565b600d544790600090670de0b6b3a7640000906200337b90849062006005565b6200338791906200617e565b9050801562003405576011546040516000916001600160a01b03169083908381818185875af1925050503d8060008114620033df576040519150601f19603f3d011682016040523d82523d6000602084013e620033e4565b606091505b5090915050801515600103620034035762003400828462005bf6565b92505b505b816023600082825462003419919062005c0c565b90915550600090505b60225481101562003736576000602860006022848154811062003449576200344962005fd3565b600091825260208083206003909202909101546001600160a01b0316835282019290925260400181206021546022805492945090918791908690811062003494576200349462005fd3565b906000526020600020906003020160010154620034b2919062006005565b620034be91906200617e565b905080156200371e576002601383015460ff166003811115620034e557620034e5620053e2565b14806200350d57506003601383015460ff1660038111156200350b576200350b620053e2565b145b15620035b557601254602280546001600160a01b0390921691631de9d9b69184918790811062003541576200354162005fd3565b600091825260209091206003909102015460405160e084901b6001600160e01b03191681526001600160a01b0390911660048201526024016000604051808303818588803b1580156200359357600080fd5b505af1158015620035a8573d6000803e3d6000fd5b5050505050505062003721565b60008260050154828460080154620035ce919062006005565b620035da91906200617e565b90506000620035ea828462005bf6565b905081156200364f5783546040516001600160a01b03909116906108fc9084906000818181858888f193505050503d806000811462003646576040519150601f19603f3d011682016040523d82523d6000602084013e6200364b565b606091505b5050505b8015620036d75783600c0160009054906101000a90046001600160a01b03166001600160a01b031663fb489a7b826040518263ffffffff1660e01b815260040160206040518083038185885af1158015620036ae573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190620036d5919062006195565b505b60018401546040518481526001600160a01b03909116907fb3b7a071186534c03b40695710096f289fd4ed6c1a374aff0bb648955e4fe5639060200160405180910390a250505b50505b806200372d8162005fe9565b91505062003422565b505050565b6004545b600554811015620037685762003762816200375a8162005fe9565b925062004430565b6200373f565b50600554600455600854600754036200377d57565b6009545b600854811015620037aa57620037a4816200379c8162005fe9565b925062004657565b62003781565b50600854600955600754805b60085481101562003824576010546000828152600660205260409020600401544391620037e39162005c0c565b116200380957620037f48162004952565b6200380160018362005c0c565b91506200380f565b62003824565b806200381b8162005fe9565b915050620037b6565b50600755565b600c54600b5410156200394757600b546000908152600a60205260409020601054600182015443916200385d9162005c0c565b1115620038675750565b600281015481546001600160a01b03908116600090815260286020526040808220600490810185905585548416835291819020600c015490516319fac8fd60e01b81529216926319fac8fd92620038c2920190815260200190565b600060405180830381600087803b158015620038dd57600080fd5b505af1158015620038f2573d6000803e3d6000fd5b5050600b80546000908152600a6020526040812080546001600160a01b0319168155600180820183905560029091018290558254909450919250906200393a90849062005c0c565b909155506200382a915050565b565b60e2815114620039935760405162461bcd60e51b8152602060048201526014602482015273092dcecc2d8d2c840e0e4dedecc40d8cadccee8d60631b604482015260640162000db7565b60308261024001515114620039eb5760405162461bcd60e51b815260206004820152601c60248201527f496e76616c696420636f6e73656e737573206b6579206c656e67746800000000604482015260640162000db7565b620039f68262004a9d565b604080518082018252601a81527f19457468657265756d205369676e6564204d6573736167653a0a00000000000060208083019190915284519251919260009262003a59920160609190911b6bffffffffffffffffffffffff1916815260140190565b604051602081830303815290604052905060008262003a79835162004bca565b8360405160200162003a8e93929190620061af565b60408051601f198184030181528282528051602091820120600280855260608501845290945060009392909183019080368337019050509050600080808062003ae6898262003ae06041600262006005565b62004ceb565b9050600062003b058a62003afd6041600262006005565b606062004ceb565b905060205b825181101562003bd65762003b20838262004e04565b6040805160008152602081018083528d905260ff8316918101919091526060810184905260808101839052929850909650945060019060a0016020604051602081039080840390855afa15801562003b7c573d6000803e3d6000fd5b5050604051601f19015190508762003b966041846200617e565b8151811062003ba95762003ba962005fd3565b6001600160a01b039092166020928302919091019091015262003bce60418262005c0c565b905062003b0a565b508a602001516001600160a01b03168660008151811062003bfb5762003bfb62005fd3565b60200260200101516001600160a01b03161462003c6d5760405162461bcd60e51b815260206004820152602960248201527f496e76616c6964206e6f6465206b6579206f776e6572736869702070726f6f66604482015268081c1c9bdd9a59195960ba1b606482015260840162000db7565b8a604001516001600160a01b03168660018151811062003c915762003c9162005fd3565b60200260200101516001600160a01b03161462003d055760405162461bcd60e51b815260206004820152602b60248201527f496e76616c6964206f7261636c65206b6579206f776e6572736869702070726f60448201526a1bd9881c1c9bdd9a59195960aa1b606482015260840162000db7565b600162003d1d8c6102400151838e6000015162004e3b565b1462003d8b5760405162461bcd60e51b815260206004820152603660248201527f496e76616c696420636f6e73656e737573206b6579206f776e65727368697020604482015275383937b7b3103337b9103932b3b4b9ba3930ba34b7b760511b606482015260840162000db7565b62003d968b62004eaa565b5050505050505050505050565b6000821162003e015760405162461bcd60e51b815260206004820152602360248201527f616d6f756e74206e65656420746f206265207374726963746c7920706f73697460448201526269766560e81b606482015260840162000db7565b6001600160a01b03811660009081526027602052604090205482111562003e6b5760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e74204e6577746f6e2062616c616e63650000000000604482015260640162000db7565b6001600160a01b0381166000908152602760205260408120805484929062003e9590849062005bf6565b9091555050604080516080810182526001600160a01b03808416825285811660208084019182528385018781524360608601908152600580546000908152600394859052978820875181549088166001600160a01b031991821617825595516001820180549190981696169590951790955590516002840155519101558054919262003f218362005fe9565b90915550506001600160a01b03848116600081815260286020908152604091829020548251908516948716948514808252918101889052909392917fc46aaee12f38035617ad448c04a7956119f7c7ed395ecc347b898817451ddb8d91015b60405180910390a35050505050565b6001600160a01b03808416600090815260286020526040902080549091838116911614806200410557600c820154604051631092ab9160e31b81526001600160a01b03858116600483015260009216906384955c8890602401602060405180830381865afa15801562004006573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200402c919062006195565b905084811015620040945760405162461bcd60e51b815260206004820152602b60248201527f696e73756666696369656e7420756e6c6f636b6564204c6971756964204e657760448201526a746f6e2062616c616e636560a81b606482015260840162000db7565b600c83015460405163282d3fdf60e01b81526001600160a01b038681166004830152602482018890529091169063282d3fdf90604401600060405180830381600087803b158015620040e557600080fd5b505af1158015620040fa573d6000803e3d6000fd5b505050505062004198565b8382600b015483600801546200411c919062005bf6565b10156200417c5760405162461bcd60e51b815260206004820152602760248201527f696e73756666696369656e742073656c6620626f6e646564206e6577746f6e2060448201526662616c616e636560c81b606482015260840162000db7565b8382600b01600082825462004192919062005c0c565b90915550505b6040805160e0810182526001600160a01b0380861682528781166020808401918252838501898152600060608601818152436080880190815260a088018381528a151560c08a019081526008805486526006909752998420985189549089166001600160a01b0319918216178a55965160018a01805491909916971696909617909655915160028701559051600386015592516004850155905160059093018054945115156101000261ff00199415159490941661ffff1990951694909417929092179092558054916200426c8362005fe9565b9190505550826001600160a01b0316856001600160a01b03167f63f8870909f7c59c9c4932bf98dbd491647c8d2e89ca0a032aacdd943a13e2fc838760405162003f809291909115158252602082015260400190565b60fa60a06000808285855af462003736573d6000803e3d6000fd5b8154600260018083161561010002038216048251808201602081106020841001600281146200438c5760018114620043b2578660005260208404602060002001600160028402018855602085068060200390508088018589016001836101000a0392508282511684540184556001840193506020820191505b8082101562004375578151845560018401935060208201915062004356565b815191036101000a90819004029091555062004427565b60028302826020036101000a846020036101000a60208901510402018501875562004427565b8660005260208404602060002001600160028402018855846020038088018589016001836101000a0392508282511660ff198a160184556020820191506001840193505b80821015620044155781518455600184019350602082019150620043f6565b815191036101000a9081900402909155505b50505050505050565b600081815260036020908152604080832060018101546001600160a01b03168452602890925282209091601382015460ff166003811115620044765762004476620053e2565b146200451657600282015482546001600160a01b031660009081526027602052604081208054909190620044ac90849062005c0c565b909155505081546001830154600284015460138401546040517f1ff2b052afa4bb37ce30d9aaccde416a700b97e632d089111749af937f8783429462004509946001600160a01b0391821694911692909160ff90911690620061f8565b60405180910390a1505050565b805482546001600160a01b0390811691161462004613576000808260080154836005015462004546919062005bf6565b9050806000036200455e578360020154915062004584565b80846002015484600d015462004575919062006005565b6200458191906200617e565b91505b600c83015484546040516340c10f1960e01b81526001600160a01b039182166004820152602481018590529116906340c10f1990604401600060405180830381600087803b158015620045d657600080fd5b505af1158015620045eb573d6000803e3d6000fd5b505050508183600d01600082825462004605919062005c0c565b909155506200463392505050565b81600201548160080160008282546200462d919062005c0c565b90915550505b81600201548160050160008282546200464d919062005c0c565b9091555050505050565b600081815260066020908152604080832060018101546001600160a01b0316845260289092528220600582015491929091610100900460ff166200484b576002830154600c8301548454604051637eee288d60e01b81526001600160a01b03918216600482015260248101849052911690637eee288d90604401600060405180830381600087803b158015620046ec57600080fd5b505af115801562004701573d6000803e3d6000fd5b50505050600c8301548454604051632770a7eb60e21b81526001600160a01b03918216600482015260248101849052911690639dc29fac90604401600060405180830381600087803b1580156200475757600080fd5b505af11580156200476c573d6000803e3d6000fd5b5050505060008360080154846005015462004788919062005bf6565b600d8501549091506200479c828462006005565b620047a891906200617e565b92508184600d016000828254620047c0919062005bf6565b90915550506006840154600003620047df576003850183905562004807565b60068401546007850154620047f5908562006005565b6200480191906200617e565b60038601555b828460060160008282546200481d919062005c0c565b909155505060038501546007850180546000906200483d90849062005c0c565b909155506200492092505050565b506002820154600882015481111562004865575060088101545b81600901546000036200487f5760038301819055620048a7565b6009820154600a83015462004895908362006005565b620048a191906200617e565b60038401555b80826009016000828254620048bd919062005c0c565b90915550506003830154600a83018054600090620048dd90849062005c0c565b9250508190555080826008016000828254620048fa919062005bf6565b90915550506002830154600b830180546000906200491a90849062005bf6565b90915550505b6005808401805460ff19166001179055820180548291906000906200494790849062005bf6565b909155505050505050565b6000818152600660205260408120600381015490910362004971575050565b60018101546001600160a01b031660009081526028602052604081206005830154909190610100900460ff1662004a0d57816007015482600601548460030154620049bd919062006005565b620049c991906200617e565b905080826006016000828254620049e1919062005bf6565b9091555050600383015460078301805460009062004a0190849062005bf6565b9091555062004a729050565b81600a01548260090154846003015462004a28919062006005565b62004a3491906200617e565b90508082600901600082825462004a4c919062005bf6565b90915550506003830154600a8301805460009062004a6c90849062005bf6565b90915550505b82546001600160a01b0316600090815260276020526040812080548392906200494790849062005c0c565b600062004aae8260600151620050db565b6001600160a01b0390911660208401529050801562004afe5760405162461bcd60e51b815260206004820152600b60248201526a32b737b2329032b93937b960a91b604482015260640162000db7565b6020808301516001600160a01b0390811660009081526028909252604090912060010154161562004b725760405162461bcd60e51b815260206004820152601c60248201527f76616c696461746f7220616c7265616479207265676973746572656400000000604482015260640162000db7565b61271082608001511115620024525760405162461bcd60e51b815260206004820152601760248201527f696e76616c696420636f6d6d697373696f6e2072617465000000000000000000604482015260640162000db7565b60608160000362004bf25750506040805180820190915260018152600360fc1b602082015290565b8160005b811562004c22578062004c098162005fe9565b915062004c1a9050600a836200617e565b915062004bf6565b60008167ffffffffffffffff81111562004c405762004c40620057ba565b6040519080825280601f01601f19166020018201604052801562004c6b576020820181803683370190505b5090505b841562004ce35762004c8360018362005bf6565b915062004c92600a8662006035565b62004c9f90603062005c0c565b60f81b81838151811062004cb75762004cb762005fd3565b60200101906001600160f81b031916908160001a90535062004cdb600a866200617e565b945062004c6f565b949350505050565b60608162004cfb81601f62005c0c565b101562004d3c5760405162461bcd60e51b815260206004820152600e60248201526d736c6963655f6f766572666c6f7760901b604482015260640162000db7565b62004d48828462005c0c565b8451101562004d8e5760405162461bcd60e51b8152602060048201526011602482015270736c6963655f6f75744f66426f756e647360781b604482015260640162000db7565b60608215801562004daf576040519150600082526020820160405262004dfb565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101562004dea57805183526020928301920162004dd0565b5050858452601f01601f1916604052505b50949350505050565b8181018051602082015160409092015190919060001a601b81101562004e345762004e31601b8262006226565b90505b9250925092565b600062004e476200526d565b600085858560405160200162004e609392919062006242565b6040516020818303038152906040529050600060fb905060008251602062004e89919062005c0c565b90506020848285855afa62004e9d57600080fd5b5050905195945050505050565b6101808101516001600160a01b031662004f2d57601d5460009062004ecf9062004bca565b90508160200151826000015183608001518360405162004eef906200528b565b62004efe949392919062006291565b604051809103906000f08015801562004f1b573d6000803e3d6000fd5b506001600160a01b0316610180830152505b60208181018051601d80546001808201835560009283527f6d4407e7be21f808e6509aa9fa9143369579dd7d760fe20a2c09680fc146134f90910180546001600160a01b03199081166001600160a01b039586161790915584518416835260289095526040918290208651815487169085161781559351908401805486169184169190911790558401516002830180549094169116179091556060820151829190600382019062004fdf9082620062d0565b506080820151600482015560a0820151600582015560c0820151600682015560e0820151600782015561010082015160088201556101208201516009820155610140820151600a820155610160820151600b820155610180820151600c820180546001600160a01b0319166001600160a01b039092169190911790556101a0820151600d8201556101c0820151600e8201556101e0820151600f820155610200820151601082015561022082015160118201556102408201516012820190620050a99082620062d0565b5061026082015160138201805460ff19166001836003811115620050d157620050d1620053e2565b0217905550505050565b600080620050e862005299565b60008060ff9050604083875160208901845afa6200510557600080fd5b50508051602090910151600160601b90910494909350915050565b60405180610280016040528060006001600160a01b0316815260200160006001600160a01b0316815260200160006001600160a01b0316815260200160608152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081526020016000815260200160006001600160a01b0316815260200160008152602001600081526020016000815260200160008152602001600081526020016060815260200160006003811115620051eb57620051eb620053e2565b905290565b6040518060a001604052806005906020820280368337509192915050565b508054600082559060005260206000209081019062000e079190620052b7565b5080546200523c9062005ba4565b6000825580601f106200524d575050565b601f01602090049060005260206000209081019062000e079190620052dc565b60405180602001604052806001906020820280368337509192915050565b611531806200639183390190565b60405180604001604052806002906020820280368337509192915050565b80821115620052d8576000620052ce82826200522e565b50600101620052b7565b5090565b5b80821115620052d85760008155600101620052dd565b60005b8381101562005310578181015183820152602001620052f6565b50506000910152565b6000815180845262005333816020860160208601620052f3565b601f01601f19169290920160200192915050565b6020815260006200535c602083018462005319565b9392505050565b6001600160a01b038116811462000e0757600080fd5b600080604083850312156200538d57600080fd5b82356200539a8162005363565b946020939093013593505050565b600060208284031215620053bb57600080fd5b81356200535c8162005363565b600060208284031215620053db57600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b600481106200541757634e487b7160e01b600052602160045260246000fd5b9052565b60208152620054366020820183516001600160a01b03169052565b600060208301516200545360408401826001600160a01b03169052565b5060408301516001600160a01b0381166060840152506060830151610280806080850152620054876102a085018362005319565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a06200550f818701836001600160a01b03169052565b8601516101c0868101919091528601516101e080870191909152860151610200808701919091528601516102208087019190915286015161024080870191909152860151858403601f19016102608088019190915290915062005573848362005319565b9350808701519150506200558a82860182620053f8565b5090949350505050565b600080600060608486031215620055aa57600080fd5b8335620055b78162005363565b92506020840135620055c98162005363565b929592945050506040919091013590565b600060208284031215620055ed57600080fd5b813567ffffffffffffffff8111156200560557600080fd5b820161028081850312156200535c57600080fd5b600081518084526020808501808196508360051b8101915082860160005b858110156200568e578284038952815180516001600160a01b031685528581015186860152604090810151606091860182905290620056798187018362005319565b9a87019a955050509084019060010162005637565b5091979650505050505050565b8315158152606060208201526000620056b8606083018562005619565b9050826040830152949350505050565b60008060408385031215620056dc57600080fd5b50508035926020909101359150565b845181526020808601518183015260408087015181840152606080880151818501526080808901516001600160a01b03908116828701528851811660a08088019190915294890151811660c087015292880151831660e0860152908701518216610100850152860151166101208301528401516102008201906200577b6101408401826001600160a01b03169052565b5083516001600160a01b0316610160830152602084015161018083015260408401516101a08301526060909301516101c08201526101e0015292915050565b634e487b7160e01b600052604160045260246000fd5b600082601f830112620057e257600080fd5b813567ffffffffffffffff80821115620058005762005800620057ba565b604051601f8301601f19908116603f011681019082821181831017156200582b576200582b620057ba565b816040528381528660208588010111156200584557600080fd5b836020870160208301376000602085830101528094505050505092915050565b600080600080608085870312156200587c57600080fd5b843567ffffffffffffffff808211156200589557600080fd5b620058a388838901620057d0565b955060208701359150620058b78262005363565b90935060408601359080821115620058ce57600080fd5b620058dc88838901620057d0565b93506060870135915080821115620058f357600080fd5b506200590287828801620057d0565b91505092959194509250565b6000602080830181845280855180835260408601915060408160051b870101925083870160005b828110156200596757603f198886030184526200595485835162005319565b9450928501929085019060010162005935565b5092979650505050505050565b6020815260006200535c602083018462005619565b6020808252825182820181905260009190848201906040850190845b81811015620059cc5783516001600160a01b031683529284019291840191600101620059a5565b50909695505050505050565b60008060408385031215620059ec57600080fd5b823567ffffffffffffffff8082111562005a0557600080fd5b62005a1386838701620057d0565b9350602085013591508082111562005a2a57600080fd5b5062005a3985828601620057d0565b9150509250929050565b60408152600062005a58604083018562005319565b828103602084015262005a6c818562005319565b95945050505050565b6000806040838503121562005a8957600080fd5b823562005a968162005363565b9150602083013562005aa88162005363565b809150509250929050565b6020808252601c908201527f76616c696461746f72206d757374206265207265676973746572656400000000604082015260600190565b6020808252602c908201527f726571756972652063616c6c657220746f2062652076616c696461746f72206160408201526b191b5a5b881858d8dbdd5b9d60a21b606082015260800190565b6020808252601a908201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604082015260600190565b60208082526018908201527f76616c696461746f72206e6f7420726567697374657265640000000000000000604082015260600190565b600181811c9082168062005bb957607f821691505b60208210810362005bda57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000d715762000d7162005be0565b8082018082111562000d715762000d7162005be0565b6000813562000d718162005363565b80546001600160a01b0319166001600160a01b0392909216919091179055565b6000808335601e1984360301811262005c6957600080fd5b83018035915067ffffffffffffffff82111562005c8557600080fd5b60200191503681900382131562005c9b57600080fd5b9250929050565b601f8211156200373657600081815260208120601f850160051c8101602086101562005ccb5750805b601f850160051c820191505b8181101562005cec5782815560010162005cd7565b505050505050565b600019600383901b1c191660019190911b1790565b67ffffffffffffffff83111562005d245762005d24620057ba565b62005d3c8362005d35835462005ba4565b8362005ca2565b6000601f84116001811462005d6f576000851562005d5a5750838201355b62005d66868262005cf4565b8455506200152c565b600083815260209020601f19861690835b8281101562005da2578685013582556020948501946001909201910162005d80565b508682101562005dc05760001960f88860031b161c19848701351681555b505060018560011b0183555050505050565b600081356004811062000d7157600080fd5b6004821062005e0357634e487b7160e01b600052602160045260246000fd5b60ff1981541660ff831681178255505050565b62005e2c62005e258362005c22565b8262005c31565b62005e4862005e3e6020840162005c22565b6001830162005c31565b62005e6462005e5a6040840162005c22565b6002830162005c31565b62005e73606083018362005c51565b62005e8381836003860162005d09565b50506080820135600482015560a0820135600582015560c0820135600682015560e0820135600782015561010082013560088201556101208201356009820155610140820135600a820155610160820135600b82015562005ef662005eec610180840162005c22565b600c830162005c31565b6101a0820135600d8201556101c0820135600e8201556101e0820135600f8201556102008201356010820155610220820135601182015562005f3d61024083018362005c51565b62005f4d81836012860162005d09565b50506200245262005f62610260840162005dd2565b6013830162005de4565b60208082526023908201527f66756e6374696f6e207265737472696374656420746f207468652070726f746f60408201526218dbdb60ea1b606082015260800190565b60006020828403121562005fc257600080fd5b815180151581146200535c57600080fd5b634e487b7160e01b600052603260045260246000fd5b60006001820162005ffe5762005ffe62005be0565b5060010190565b808202811582820484141762000d715762000d7162005be0565b634e487b7160e01b600052601260045260246000fd5b6000826200604757620060476200601f565b500690565b600060018060a01b0380881683528087166020840152808616604084015260a060608401526200608060a084018662005319565b91508084166080840152509695505050505050565b818103620060a1575050565b620060ad825462005ba4565b67ffffffffffffffff811115620060c857620060c8620057ba565b620060e081620060d9845462005ba4565b8462005ca2565b6000601f821160018114620061135760008315620060fe5750848201545b6200610a848262005cf4565b8555506200152c565b600085815260209020601f19841690600086815260209020845b838110156200614f57828601548255600195860195909101906020016200612d565b50858310156200616e5781850154600019600388901b60f8161c191681555b5050505050600190811b01905550565b6000826200619057620061906200601f565b500490565b600060208284031215620061a857600080fd5b5051919050565b60008451620061c3818460208901620052f3565b845190830190620061d9818360208901620052f3565b8451910190620061ee818360208801620052f3565b0195945050505050565b6001600160a01b03858116825284166020820152604081018390526080810162005a6c6060830184620053f8565b60ff818116838216019081111562000d715762000d7162005be0565b6000845162006256818460208901620052f3565b8451908301906200626c818360208901620052f3565b60609490941b6bffffffffffffffffffffffff19169301928352505060140192915050565b6001600160a01b0385811682528416602082015260408101839052608060608201819052600090620062c69083018462005319565b9695505050505050565b815167ffffffffffffffff811115620062ed57620062ed620057ba565b620062fe81620060d9845462005ba4565b602080601f8311600181146200633257600084156200631d5750858301515b62006329858262005cf4565b86555062005cec565b600085815260208120601f198616915b82811015620063635788860151825594840194600190910190840162006342565b50858210156200616e57939096015160001960f8600387901b161c19169092555050600190811b0190555056fe60806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220018b8c3d2155d9eb67d018108e35d70682e98011384859e6cefa8fd80d2593e364736f6c63430008150033a2646970667358221220ae7f19a50c64bde52a9480de7de6252890232586c5aa23e239cd92f69bdf42ce64736f6c6343000815003360806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220018b8c3d2155d9eb67d018108e35d70682e98011384859e6cefa8fd80d2593e364736f6c63430008150033",
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
// Solidity: function computeCommittee() returns(address[])
func (_AutonityUpgradeTest *AutonityUpgradeTest) ComputeCommittee(opts *runOptions) (uint64, error) {
	_, consumed, err := _AutonityUpgradeTest.call(opts, "computeCommittee")
	return consumed, err
}

// Finalize is a paid mutator transaction binding the contract method 0x4bb278f3.
//
// Solidity: function finalize() returns(bool, (address,uint256,bytes)[], uint256)
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
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterSlashingEvent is a free log retrieval operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
 		func (_IAccountability *IAccountability) FilterSlashingEvent(opts *bind.FilterOpts) (*IAccountabilitySlashingEventIterator, error) {






			logs, sub, err := _IAccountability.contract.FilterLogs(opts, "SlashingEvent")
			if err != nil {
				return nil, err
			}
			return &IAccountabilitySlashingEventIterator{contract: _IAccountability.contract, event: "SlashingEvent", logs: logs, sub: sub}, nil
 		}

		// WatchSlashingEvent is a free log subscription operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
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

		// ParseSlashingEvent is a log parse operation binding the contract event 0x876592f0699764d3d26c3faee6876cfc5b7b3ef0bd909da978591ce77829aee8.
		//
		// Solidity: event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound)
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
	Bin: "0x60806040523480156200001157600080fd5b506040516200153138038062001531833981016040819052620000349162000151565b6127108211156200004457600080fd5b600a80546001600160a01b038087166001600160a01b031992831617909255600b805492861692909116919091179055600c8290556040516200008c9082906020016200023e565b60405160208183030381529060405260089081620000ab9190620002fc565b5080604051602001620000bf91906200023e565b60405160208183030381529060405260099081620000de9190620002fc565b5050600080546001600160a01b0319163317905550620003c8915050565b6001600160a01b03811681146200011257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60005b83811015620001485781810151838201526020016200012e565b50506000910152565b600080600080608085870312156200016857600080fd5b84516200017581620000fc565b60208601519094506200018881620000fc565b6040860151606087015191945092506001600160401b0380821115620001ad57600080fd5b818701915087601f830112620001c257600080fd5b815181811115620001d757620001d762000115565b604051601f8201601f19908116603f0116810190838211818310171562000202576200020262000115565b816040528281528a60208487010111156200021c57600080fd5b6200022f8360208301602088016200012b565b979a9699509497505050505050565b644c4e544e2d60d81b815260008251620002608160058501602087016200012b565b9190910160050192915050565b600181811c908216806200028257607f821691505b602082108103620002a357634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620002f757600081815260208120601f850160051c81016020861015620002d25750805b601f850160051c820191505b81811015620002f357828155600101620002de565b5050505b505050565b81516001600160401b0381111562000318576200031862000115565b62000330816200032984546200026d565b84620002a9565b602080601f8311600181146200036857600084156200034f5750858301515b600019600386901b1c1916600185901b178555620002f3565b600085815260208120601f198616915b82811015620003995788860151825594840194600190910190840162000378565b5085821015620003b85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61115980620003d86000396000f3fe60806040526004361061014b5760003560e01c806359355736116100b6578063949813b81161006f578063949813b8146103c557806395d89b41146103e55780639dc29fac146103fa578063a9059cbb1461041a578063dd62ed3e1461043a578063fb489a7b1461048057600080fd5b806359355736146102e35780635ea1d6f81461031957806361d027b31461032f57806370a082311461034f5780637eee288d1461038557806384955c88146103a557600080fd5b8063282d3fdf11610108578063282d3fdf146102245780632f2c3f2e14610244578063313ce5671461025a578063372500ab146102765780633a5381b51461028b57806340c10f19146102c357600080fd5b806306fdde0314610150578063095ea7b31461017b57806318160ddd146101ab578063187cf4d7146101ca57806319fac8fd146101e257806323b872dd14610204575b600080fd5b34801561015c57600080fd5b50610165610488565b6040516101729190610ece565b60405180910390f35b34801561018757600080fd5b5061019b610196366004610f38565b61051a565b6040519015158152602001610172565b3480156101b757600080fd5b506004545b604051908152602001610172565b3480156101d657600080fd5b506101bc633b9aca0081565b3480156101ee57600080fd5b506102026101fd366004610f62565b610531565b005b34801561021057600080fd5b5061019b61021f366004610f7b565b610569565b34801561023057600080fd5b5061020261023f366004610f38565b61065c565b34801561025057600080fd5b506101bc61271081565b34801561026657600080fd5b5060405160128152602001610172565b34801561028257600080fd5b50610202610741565b34801561029757600080fd5b50600a546102ab906001600160a01b031681565b6040516001600160a01b039091168152602001610172565b3480156102cf57600080fd5b506102026102de366004610f38565b6107ef565b3480156102ef57600080fd5b506101bc6102fe366004610fb7565b6001600160a01b031660009081526002602052604090205490565b34801561032557600080fd5b506101bc600c5481565b34801561033b57600080fd5b50600b546102ab906001600160a01b031681565b34801561035b57600080fd5b506101bc61036a366004610fb7565b6001600160a01b031660009081526001602052604090205490565b34801561039157600080fd5b506102026103a0366004610f38565b610857565b3480156103b157600080fd5b506101bc6103c0366004610fb7565b61091d565b3480156103d157600080fd5b506101bc6103e0366004610fb7565b61094b565b3480156103f157600080fd5b50610165610979565b34801561040657600080fd5b50610202610415366004610f38565b610988565b34801561042657600080fd5b5061019b610435366004610f38565b6109e8565b34801561044657600080fd5b506101bc610455366004610fd9565b6001600160a01b03918216600090815260036020908152604080832093909416825291909152205490565b6101bc610a35565b6060600880546104979061100c565b80601f01602080910402602001604051908101604052809291908181526020018280546104c39061100c565b80156105105780601f106104e557610100808354040283529160200191610510565b820191906000526020600020905b8154815290600101906020018083116104f357829003601f168201915b5050505050905090565b6000610527338484610b9d565b5060015b92915050565b6000546001600160a01b031633146105645760405162461bcd60e51b815260040161055b90611046565b60405180910390fd5b600c55565b6001600160a01b0383166000908152600360209081526040808320338452909152812054828110156105ee5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b606482015260840161055b565b61060285336105fd86856110a4565b610b9d565b61060c8584610cc1565b6106168484610da9565b836001600160a01b0316856001600160a01b03166000805160206111048339815191528560405161064991815260200190565b60405180910390a3506001949350505050565b6000546001600160a01b031633146106865760405162461bcd60e51b815260040161055b90611046565b6001600160a01b03821660009081526002602090815260408083205460019092529091205482916106b6916110a4565b10156107105760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201526361626c6560e01b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110b7565b90915550505050565b600061074c33610df4565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d806000811461079e576040519150601f19603f3d011682016040523d82523d6000602084013e6107a3565b606091505b50509050806107eb5760405162461bcd60e51b81526020600482015260146024820152732330b4b632b2103a379039b2b7321022ba3432b960611b604482015260640161055b565b5050565b6000546001600160a01b031633146108195760405162461bcd60e51b815260040161055b90611046565b6108238282610da9565b6040518181526001600160a01b03831690600090600080516020611104833981519152906020015b60405180910390a35050565b6000546001600160a01b031633146108815760405162461bcd60e51b815260040161055b90611046565b6001600160a01b0382166000908152600260205260409020548111156108f55760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f636044820152621ad95960ea1b606482015260840161055b565b6001600160a01b038216600090815260026020526040812080548392906107389084906110a4565b6001600160a01b038116600090815260026020908152604080832054600190925282205461052b91906110a4565b600061095682610e59565b6001600160a01b03831660009081526005602052604090205461052b91906110b7565b6060600980546104979061100c565b6000546001600160a01b031633146109b25760405162461bcd60e51b815260040161055b90611046565b6109bc8282610cc1565b6040518181526000906001600160a01b038416906000805160206111048339815191529060200161084b565b60006109f43383610cc1565b6109fe8383610da9565b6040518281526001600160a01b0384169033906000805160206111048339815191529060200160405180910390a350600192915050565b600080546001600160a01b03163314610a605760405162461bcd60e51b815260040161055b90611046565b600c54349060009061271090610a7690846110ca565b610a8091906110e1565b905081811115610ad25760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f72207265776172640000000000000000604482015260640161055b565b610adc81836110a4565b600b546040519193506001600160a01b0316906108fc9083906000818181858888f193505050503d8060008114610b2f576040519150601f19603f3d011682016040523d82523d6000602084013e610b34565b606091505b505060045460009150610b4b633b9aca00856110ca565b610b5591906110e1565b905080600754610b6591906110b7565b600755600454600090633b9aca0090610b7e90846110ca565b610b8891906110e1565b9050610b9481846110b7565b94505050505090565b6001600160a01b038316610bff5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161055b565b6001600160a01b038216610c605760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161055b565b6001600160a01b0383811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610cca82610df4565b506001600160a01b038216600090815260016020908152604080832054600290925290912054610cfa90826110a4565b821115610d495760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e64730000000000604482015260640161055b565b610d5382826110a4565b6001600160a01b038416600090815260016020526040902055808203610d8d576001600160a01b0383166000908152600660205260408120555b8160046000828254610d9f91906110a4565b9091555050505050565b610db282610df4565b506001600160a01b03821660009081526001602052604081208054839290610ddb9084906110b7565b92505081905550806004600082825461073891906110b7565b600080610e0083610e59565b6001600160a01b038416600090815260056020526040902054909150610e279082906110b7565b6001600160a01b0390931660009081526005602090815260408083208690556007546006909252909120555090919050565b6001600160a01b038116600090815260016020526040812054808203610e825750600092915050565b6001600160a01b038316600090815260066020526040812054600754610ea891906110a4565b90506000633b9aca00610ebb84846110ca565b610ec591906110e1565b95945050505050565b600060208083528351808285015260005b81811015610efb57858101830151858201604001528201610edf565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b0381168114610f3357600080fd5b919050565b60008060408385031215610f4b57600080fd5b610f5483610f1c565b946020939093013593505050565b600060208284031215610f7457600080fd5b5035919050565b600080600060608486031215610f9057600080fd5b610f9984610f1c565b9250610fa760208501610f1c565b9150604084013590509250925092565b600060208284031215610fc957600080fd5b610fd282610f1c565b9392505050565b60008060408385031215610fec57600080fd5b610ff583610f1c565b915061100360208401610f1c565b90509250929050565b600181811c9082168061102057607f821691505b60208210810361104057634e487b7160e01b600052602260045260246000fd5b50919050565b60208082526028908201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060408201526710dbdb9d1c9858dd60c21b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b8181038181111561052b5761052b61108e565b8082018082111561052b5761052b61108e565b808202811582820484141761052b5761052b61108e565b6000826110fe57634e487b7160e01b600052601260045260246000fd5b50049056feddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa2646970667358221220018b8c3d2155d9eb67d018108e35d70682e98011384859e6cefa8fd80d2593e364736f6c63430008150033",
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
	Bin: "0x608060405234801561001057600080fd5b5060405161041238038061041283398101604081905261002f91610058565b600061003b82826101aa565b5050610269565b634e487b7160e01b600052604160045260246000fd5b6000602080838503121561006b57600080fd5b82516001600160401b038082111561008257600080fd5b818501915085601f83011261009657600080fd5b8151818111156100a8576100a8610042565b604051601f8201601f19908116603f011681019083821181831017156100d0576100d0610042565b8160405282815288868487010111156100e857600080fd5b600093505b8284101561010a57848401860151818501870152928501926100ed565b600086848301015280965050505050505092915050565b600181811c9082168061013557607f821691505b60208210810361015557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156101a557600081815260208120601f850160051c810160208610156101825750805b601f850160051c820191505b818110156101a15782815560010161018e565b5050505b505050565b81516001600160401b038111156101c3576101c3610042565b6101d7816101d18454610121565b8461015b565b602080601f83116001811461020c57600084156101f45750858301515b600019600386901b1c1916600185901b1785556101a1565b600085815260208120601f198616915b8281101561023b5788860151825594840194600190910190840161021c565b50858210156102595787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61019a806102786000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063bfb4ebcf14610030575b600080fd5b61003861004e565b60405161004591906100dc565b60405180910390f35b6000805461005b9061012a565b80601f01602080910402602001604051908101604052809291908181526020018280546100879061012a565b80156100d45780601f106100a9576101008083540402835291602001916100d4565b820191906000526020600020905b8154815290600101906020018083116100b757829003601f168201915b505050505081565b600060208083528351808285015260005b81811015610109578581018301518582016040015282016100ed565b506000604082860101526040601f19601f8301168501019250505092915050565b600181811c9082168061013e57607f821691505b60208210810361015e57634e487b7160e01b600052602260045260246000fd5b5091905056fea26469706673582212202cc3b202f69ac4e15a8d9d29e8d6a487680dc277758463ce9725edc64656c36764736f6c63430008150033",
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
	Bin: "0x608060405234801561001057600080fd5b5060405161068838038061068883398101604081905261002f9161010e565b80600061003c82826101fa565b506001905061004b83826101fa565b5050506102b9565b634e487b7160e01b600052604160045260246000fd5b600082601f83011261007a57600080fd5b81516001600160401b038082111561009457610094610053565b604051601f8301601f19908116603f011681019082821181831017156100bc576100bc610053565b816040528381526020925086838588010111156100d857600080fd5b600091505b838210156100fa57858201830151818301840152908201906100dd565b600093810190920192909252949350505050565b6000806040838503121561012157600080fd5b82516001600160401b038082111561013857600080fd5b61014486838701610069565b9350602085015191508082111561015a57600080fd5b5061016785828601610069565b9150509250929050565b600181811c9082168061018557607f821691505b6020821081036101a557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156101f557600081815260208120601f850160051c810160208610156101d25750805b601f850160051c820191505b818110156101f1578281556001016101de565b5050505b505050565b81516001600160401b0381111561021357610213610053565b610227816102218454610171565b846101ab565b602080601f83116001811461025c57600084156102445750858301515b600019600386901b1c1916600185901b1785556101f1565b600085815260208120601f198616915b8281101561028b5788860151825594840194600190910190840161026c565b50858210156102a95787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6103c0806102c86000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80631e4f339514610046578063b0a378b01461005b578063bfb4ebcf14610079575b600080fd5b610059610054366004610142565b610081565b005b610063610091565b60405161007091906101f3565b60405180910390f35b61006361011f565b600061008d82826102ca565b5050565b6001805461009e90610241565b80601f01602080910402602001604051908101604052809291908181526020018280546100ca90610241565b80156101175780601f106100ec57610100808354040283529160200191610117565b820191906000526020600020905b8154815290600101906020018083116100fa57829003601f168201915b505050505081565b6000805461009e90610241565b634e487b7160e01b600052604160045260246000fd5b60006020828403121561015457600080fd5b813567ffffffffffffffff8082111561016c57600080fd5b818401915084601f83011261018057600080fd5b8135818111156101925761019261012c565b604051601f8201601f19908116603f011681019083821181831017156101ba576101ba61012c565b816040528281528760208487010111156101d357600080fd5b826020860160208301376000928101602001929092525095945050505050565b600060208083528351808285015260005b8181101561022057858101830151858201604001528201610204565b506000604082860101526040601f19601f8301168501019250505092915050565b600181811c9082168061025557607f821691505b60208210810361027557634e487b7160e01b600052602260045260246000fd5b50919050565b601f8211156102c557600081815260208120601f850160051c810160208610156102a25750805b601f850160051c820191505b818110156102c1578281556001016102ae565b5050505b505050565b815167ffffffffffffffff8111156102e4576102e461012c565b6102f8816102f28454610241565b8461027b565b602080601f83116001811461032d57600084156103155750858301515b600019600386901b1c1916600185901b1785556102c1565b600085815260208120601f198616915b8281101561035c5788860151825594840194600190910190840161033d565b508582101561037a5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea26469706673582212208b04b4bf1ad6b8710ee7d7320308028284fa45655b15c64a1494b25e8580837364736f6c63430008150033",
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
