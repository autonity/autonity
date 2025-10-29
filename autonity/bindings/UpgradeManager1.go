// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// UpgradeManager1version is an auto generated low-level Go binding around an user-defined struct.
type UpgradeManager1version struct {
	Number string
	Block  *big.Int
}

// UpgradeManager1MetaData contains all meta data concerning the UpgradeManager1 contract.
var UpgradeManager1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_hashes\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"number\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"}],\"internalType\":\"structUpgradeManager1.version[]\",\"name\":\"_versions\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newValue\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"oldValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newValue\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateBool\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"oldValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"newValue\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateInt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"appliesAtHeight\",\"type\":\"uint256\"}],\"name\":\"ConfigUpdateUint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"UpgradeResult\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getAutonity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"}],\"name\":\"getVersion\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"number\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"}],\"internalType\":\"structUpgradeManager1.version\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"setOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_hash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"number\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"}],\"internalType\":\"structUpgradeManager1.version\",\"name\":\"_version\",\"type\":\"tuple\"}],\"name\":\"setVersion\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_data\",\"type\":\"string\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_data\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_versionString\",\"type\":\"string\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_targets\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"_bytecodes\",\"type\":\"string[]\"}],\"name\":\"upgradeMultiple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_targets\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"_bytecodes\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"_versionStrings\",\"type\":\"string[]\"}],\"name\":\"upgradeMultiple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7c8ccebe": "getAutonity()",
		"e7f43c68": "getOperator()",
		"9aaf9f08": "getVersion(bytes32)",
		"b3ab15fb": "setOperator(address)",
		"b4062f8e": "setVersion(bytes32,(string,uint256))",
		"6e3d9ff0": "upgrade(address,string)",
		"82f75182": "upgrade(address,string,string)",
		"2011c815": "upgradeMultiple(address[],string[])",
		"9d897130": "upgradeMultiple(address[],string[],string[])",
	},
	Bin: "0x608060405234801561001057600080fd5b506040516117ab3803806117ab83398101604081905261002f916102ee565b80518251146100965760405162461bcd60e51b815260206004820152602960248201527f68617368657320616e642076657273696f6e73206861766520646966666572656044820152680dce840d8cadccee8d60bb1b606482015260840160405180910390fd5b60005b8251811015610111578181815181106100b4576100b46103b0565b6020026020010151600260008584815181106100d2576100d26103b0565b6020026020010151815260200190815260200160002060008201518160000190816100fd919061044f565b506020919091015160019182015501610099565b50505061050d565b634e487b7160e01b600052604160045260246000fd5b604080519081016001600160401b038111828210171561015157610151610119565b60405290565b604051601f8201601f191681016001600160401b038111828210171561017f5761017f610119565b604052919050565b60006001600160401b038211156101a0576101a0610119565b5060051b60200190565b600082601f8301126101bb57600080fd5b81516101ce6101c982610187565b610157565b8082825260208201915060208360051b8601019250858311156101f057600080fd5b602085015b838110156102e45780516001600160401b0381111561021357600080fd5b86016040818903601f1901121561022957600080fd5b61023161012f565b60208201516001600160401b0381111561024a57600080fd5b82016020810190603f018a1361025f57600080fd5b80516001600160401b0381111561027857610278610119565b61028b601f8201601f1916602001610157565b8181528b60208385010111156102a057600080fd5b60005b828110156102bf576020818501810151838301820152016102a3565b50600060209282018301528352604093909301518284015250845292830192016101f5565b5095945050505050565b6000806040838503121561030157600080fd5b82516001600160401b0381111561031757600080fd5b8301601f8101851361032857600080fd5b80516103366101c982610187565b8082825260208201915060208360051b85010192508783111561035857600080fd5b6020840193505b8284101561037a57835182526020938401939091019061035f565b6020870151909550925050506001600160401b0381111561039a57600080fd5b6103a6858286016101aa565b9150509250929050565b634e487b7160e01b600052603260045260246000fd5b600181811c908216806103da57607f821691505b6020821081036103fa57634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111561044a57806000526020600020601f840160051c810160208510156104275750805b601f840160051c820191505b818110156104475760008155600101610433565b50505b505050565b81516001600160401b0381111561046857610468610119565b61047c8161047684546103c6565b84610400565b6020601f8211600181146104b057600083156104985750848201515b600019600385901b1c1916600184901b178455610447565b600084815260208120601f198516915b828110156104e057878501518255602094850194600190920191016104c0565b50848210156104fe5786840151600019600387901b60f8161c191681555b50505050600190811b01905550565b61128f8061051c6000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c80639aaf9f0811610076578063b3ab15fb1161005b578063b3ab15fb1461015a578063b4062f8e1461016d578063e7f43c681461018057600080fd5b80639aaf9f08146101275780639d8971301461014757600080fd5b80632011c815146100a85780636e3d9ff0146100bd5780637c8ccebe146100d057806382f7518214610114575b600080fd5b6100bb6100b6366004610d0d565b61019e565b005b6100bb6100cb366004610d76565b61030a565b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100bb610122366004610dba565b610399565b61013a610135366004610e33565b610481565b60405161010b9190610e70565b6100bb610155366004610ed4565b610550565b6100bb610168366004610f5e565b6107ce565b6100bb61017b366004610f80565b610949565b60015473ffffffffffffffffffffffffffffffffffffffff166100ea565b60015473ffffffffffffffffffffffffffffffffffffffff163314610224576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064015b60405180910390fd5b80518251146102b5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603060248201527f61646472657373657320616e642062797465636f6465732073686f756c64206260448201527f6520696e2073616d65206e756d62657200000000000000000000000000000000606482015260840161021b565b60005b8251811015610305576102fd8382815181106102d6576102d6611025565b60200260200101518383815181106102f0576102f0611025565b60200260200101516109fb565b6001016102b8565b505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461038b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604482015260640161021b565b61039582826109fb565b5050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461041a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604482015260640161021b565b61042483836109fb565b6040805180820182528281524360208083019190915273ffffffffffffffffffffffffffffffffffffffff86163f60009081526002909152919091208151819061046e90826110f5565b5060208201518160010155905050505050565b6040805180820190915260608152600060208201526000828152600260205260409081902081518083019092528054829082906104bd90611054565b80601f01602080910402602001604051908101604052809291908181526020018280546104e990611054565b80156105365780601f1061050b57610100808354040283529160200191610536565b820191906000526020600020905b81548152906001019060200180831161051957829003601f168201915b505050505081526020016001820154815250509050919050565b60015473ffffffffffffffffffffffffffffffffffffffff1633146105d1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604482015260640161021b565b8151835114610662576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603060248201527f61646472657373657320616e642062797465636f6465732073686f756c64206260448201527f6520696e2073616d65206e756d62657200000000000000000000000000000000606482015260840161021b565b80518351146106f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603660248201527f61646472657373657320616e642076657273696f6e20737472696e677320736860448201527f6f756c6420626520696e2073616d65206e756d62657200000000000000000000606482015260840161021b565b60005b83518110156107c85761072e84828151811061071457610714611025565b60200260200101518483815181106102f0576102f0611025565b604051806040016040528083838151811061074b5761074b611025565b60200260200101518152602001438152506002600086848151811061077257610772611025565b602002602001015173ffffffffffffffffffffffffffffffffffffffff163f815260200190815260200160002060008201518160000190816107b491906110f5565b5060209190910151600191820155016106f6565b50505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610875576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f63616c6c6572206973206e6f7420746865204175746f6e69747920636f6e747260448201527f6163740000000000000000000000000000000000000000000000000000000000606482015260840161021b565b6001546040805160808082526008908201527f6f70657261746f7200000000000000000000000000000000000000000000000060a082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152436060830152517fe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a09181900360c00190a1600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60015473ffffffffffffffffffffffffffffffffffffffff1633146109ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f72000000000000604482015260640161021b565b6000828152600260205260409020815182919081906109e990826110f5565b50602082015181600101559050505050565b60405160f990600090610a14908590859060200161120e565b60405160208183030381529060405290506000806060600080855160208701885af4809350503d91506040519050602082018101604052818152816000602083013e8673ffffffffffffffffffffffffffffffffffffffff167f852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d393479884604051610aa0911515815260200190565b60405180910390a282610ab4578160208201fd5b50505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610b3357610b33610abd565b604052919050565b600067ffffffffffffffff821115610b5557610b55610abd565b5060051b60200190565b803573ffffffffffffffffffffffffffffffffffffffff81168114610b8357600080fd5b919050565b600082601f830112610b9957600080fd5b8135610bac610ba782610b3b565b610aec565b8082825260208201915060208360051b860101925085831115610bce57600080fd5b602085015b83811015610bf257610be481610b5f565b835260209283019201610bd3565b5095945050505050565b600082601f830112610c0d57600080fd5b813567ffffffffffffffff811115610c2757610c27610abd565b610c5860207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610aec565b818152846020838601011115610c6d57600080fd5b816020850160208301376000918101602001919091529392505050565b600082601f830112610c9b57600080fd5b8135610ca9610ba782610b3b565b8082825260208201915060208360051b860101925085831115610ccb57600080fd5b602085015b83811015610bf257803567ffffffffffffffff811115610cef57600080fd5b610cfe886020838a0101610bfc565b84525060209283019201610cd0565b60008060408385031215610d2057600080fd5b823567ffffffffffffffff811115610d3757600080fd5b610d4385828601610b88565b925050602083013567ffffffffffffffff811115610d6057600080fd5b610d6c85828601610c8a565b9150509250929050565b60008060408385031215610d8957600080fd5b610d9283610b5f565b9150602083013567ffffffffffffffff811115610dae57600080fd5b610d6c85828601610bfc565b600080600060608486031215610dcf57600080fd5b610dd884610b5f565b9250602084013567ffffffffffffffff811115610df457600080fd5b610e0086828701610bfc565b925050604084013567ffffffffffffffff811115610e1d57600080fd5b610e2986828701610bfc565b9150509250925092565b600060208284031215610e4557600080fd5b5035919050565b60005b83811015610e67578181015183820152602001610e4f565b50506000910152565b6020815260008251604060208401528051806060850152610e98816080860160208501610e4c565b60209490940151604084015250506080601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010190565b600080600060608486031215610ee957600080fd5b833567ffffffffffffffff811115610f0057600080fd5b610f0c86828701610b88565b935050602084013567ffffffffffffffff811115610f2957600080fd5b610f3586828701610c8a565b925050604084013567ffffffffffffffff811115610f5257600080fd5b610e2986828701610c8a565b600060208284031215610f7057600080fd5b610f7982610b5f565b9392505050565b60008060408385031215610f9357600080fd5b82359150602083013567ffffffffffffffff811115610fb157600080fd5b830160408186031215610fc357600080fd5b6040805190810167ffffffffffffffff81118282101715610fe657610fe6610abd565b604052813567ffffffffffffffff81111561100057600080fd5b61100c87828501610bfc565b8252506020820135602082015280925050509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061106857607f821691505b6020821081036110a1577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561030557806000526020600020601f840160051c810160208510156110ce5750805b601f840160051c820191505b818110156110ee57600081556001016110da565b5050505050565b815167ffffffffffffffff81111561110f5761110f610abd565b6111238161111d8454611054565b846110a7565b6020601f821160018114611175576000831561113f5750848201515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600385901b1c1916600184901b1784556110ee565b6000848152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08516915b828110156111c357878501518255602094850194600190920191016111a3565b50848210156111ff57868401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b60f8161c191681555b50505050600190811b01905550565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008360601b1681526000825161124b816014850160208701610e4c565b91909101601401939250505056fea26469706673582212201a854cdb664a83a91515e263c540f1c60ba555863bfc3ea26a4d327909cfddc464736f6c634300081e0033",
}

// UpgradeManager1ABI is the input ABI used to generate the binding from.
// Deprecated: Use UpgradeManager1MetaData.ABI instead.
var UpgradeManager1ABI = UpgradeManager1MetaData.ABI

// Deprecated: Use UpgradeManager1MetaData.Sigs instead.
// UpgradeManager1FuncSigs maps the 4-byte function signature to its string representation.
var UpgradeManager1FuncSigs = UpgradeManager1MetaData.Sigs

// UpgradeManager1Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UpgradeManager1MetaData.Bin instead.
var UpgradeManager1Bin = UpgradeManager1MetaData.Bin

// DeployUpgradeManager1 deploys a new Ethereum contract, binding an instance of UpgradeManager1 to it.
func DeployUpgradeManager1(auth *bind.TransactOpts, backend bind.ContractBackend, _hashes [][32]byte, _versions []UpgradeManager1version) (common.Address, *types.Transaction, *UpgradeManager1, error) {
	parsed, err := UpgradeManager1MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UpgradeManager1Bin), backend, _hashes, _versions)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UpgradeManager1{UpgradeManager1Caller: UpgradeManager1Caller{contract: contract}, UpgradeManager1Transactor: UpgradeManager1Transactor{contract: contract}, UpgradeManager1Filterer: UpgradeManager1Filterer{contract: contract}}, nil
}

// UpgradeManager1 is an auto generated Go binding around an Ethereum contract.
type UpgradeManager1 struct {
	UpgradeManager1Caller     // Read-only binding to the contract
	UpgradeManager1Transactor // Write-only binding to the contract
	UpgradeManager1Filterer   // Log filterer for contract events
}

// UpgradeManager1Caller is an auto generated read-only Go binding around an Ethereum contract.
type UpgradeManager1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManager1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type UpgradeManager1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManager1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UpgradeManager1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeManager1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UpgradeManager1Session struct {
	Contract     *UpgradeManager1  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpgradeManager1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UpgradeManager1CallerSession struct {
	Contract *UpgradeManager1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// UpgradeManager1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UpgradeManager1TransactorSession struct {
	Contract     *UpgradeManager1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// UpgradeManager1Raw is an auto generated low-level Go binding around an Ethereum contract.
type UpgradeManager1Raw struct {
	Contract *UpgradeManager1 // Generic contract binding to access the raw methods on
}

// UpgradeManager1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UpgradeManager1CallerRaw struct {
	Contract *UpgradeManager1Caller // Generic read-only contract binding to access the raw methods on
}

// UpgradeManager1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UpgradeManager1TransactorRaw struct {
	Contract *UpgradeManager1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeManager1 creates a new instance of UpgradeManager1, bound to a specific deployed contract.
func NewUpgradeManager1(address common.Address, backend bind.ContractBackend) (*UpgradeManager1, error) {
	contract, err := bindUpgradeManager1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1{UpgradeManager1Caller: UpgradeManager1Caller{contract: contract}, UpgradeManager1Transactor: UpgradeManager1Transactor{contract: contract}, UpgradeManager1Filterer: UpgradeManager1Filterer{contract: contract}}, nil
}

// NewUpgradeManager1Caller creates a new read-only instance of UpgradeManager1, bound to a specific deployed contract.
func NewUpgradeManager1Caller(address common.Address, caller bind.ContractCaller) (*UpgradeManager1Caller, error) {
	contract, err := bindUpgradeManager1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1Caller{contract: contract}, nil
}

// NewUpgradeManager1Transactor creates a new write-only instance of UpgradeManager1, bound to a specific deployed contract.
func NewUpgradeManager1Transactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeManager1Transactor, error) {
	contract, err := bindUpgradeManager1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1Transactor{contract: contract}, nil
}

// NewUpgradeManager1Filterer creates a new log filterer instance of UpgradeManager1, bound to a specific deployed contract.
func NewUpgradeManager1Filterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeManager1Filterer, error) {
	contract, err := bindUpgradeManager1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1Filterer{contract: contract}, nil
}

// bindUpgradeManager1 binds a generic wrapper to an already deployed contract.
func bindUpgradeManager1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeManager1ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeManager1 *UpgradeManager1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeManager1.Contract.UpgradeManager1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeManager1 *UpgradeManager1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.UpgradeManager1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeManager1 *UpgradeManager1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.UpgradeManager1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UpgradeManager1 *UpgradeManager1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UpgradeManager1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UpgradeManager1 *UpgradeManager1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UpgradeManager1 *UpgradeManager1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.contract.Transact(opts, method, params...)
}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager1 *UpgradeManager1Caller) GetAutonity(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpgradeManager1.contract.Call(opts, &out, "getAutonity")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager1 *UpgradeManager1Session) GetAutonity() (common.Address, error) {
	return _UpgradeManager1.Contract.GetAutonity(&_UpgradeManager1.CallOpts)
}

// GetAutonity is a free data retrieval call binding the contract method 0x7c8ccebe.
//
// Solidity: function getAutonity() view returns(address)
func (_UpgradeManager1 *UpgradeManager1CallerSession) GetAutonity() (common.Address, error) {
	return _UpgradeManager1.Contract.GetAutonity(&_UpgradeManager1.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager1 *UpgradeManager1Caller) GetOperator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UpgradeManager1.contract.Call(opts, &out, "getOperator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager1 *UpgradeManager1Session) GetOperator() (common.Address, error) {
	return _UpgradeManager1.Contract.GetOperator(&_UpgradeManager1.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0xe7f43c68.
//
// Solidity: function getOperator() view returns(address)
func (_UpgradeManager1 *UpgradeManager1CallerSession) GetOperator() (common.Address, error) {
	return _UpgradeManager1.Contract.GetOperator(&_UpgradeManager1.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x9aaf9f08.
//
// Solidity: function getVersion(bytes32 _hash) view returns((string,uint256))
func (_UpgradeManager1 *UpgradeManager1Caller) GetVersion(opts *bind.CallOpts, _hash [32]byte) (UpgradeManager1version, error) {
	var out []interface{}
	err := _UpgradeManager1.contract.Call(opts, &out, "getVersion", _hash)

	if err != nil {
		return *new(UpgradeManager1version), err
	}

	out0 := *abi.ConvertType(out[0], new(UpgradeManager1version)).(*UpgradeManager1version)

	return out0, err

}

// GetVersion is a free data retrieval call binding the contract method 0x9aaf9f08.
//
// Solidity: function getVersion(bytes32 _hash) view returns((string,uint256))
func (_UpgradeManager1 *UpgradeManager1Session) GetVersion(_hash [32]byte) (UpgradeManager1version, error) {
	return _UpgradeManager1.Contract.GetVersion(&_UpgradeManager1.CallOpts, _hash)
}

// GetVersion is a free data retrieval call binding the contract method 0x9aaf9f08.
//
// Solidity: function getVersion(bytes32 _hash) view returns((string,uint256))
func (_UpgradeManager1 *UpgradeManager1CallerSession) GetVersion(_hash [32]byte) (UpgradeManager1version, error) {
	return _UpgradeManager1.Contract.GetVersion(&_UpgradeManager1.CallOpts, _hash)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) SetOperator(opts *bind.TransactOpts, _account common.Address) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "setOperator", _account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager1 *UpgradeManager1Session) SetOperator(_account common.Address) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.SetOperator(&_UpgradeManager1.TransactOpts, _account)
}

// SetOperator is a paid mutator transaction binding the contract method 0xb3ab15fb.
//
// Solidity: function setOperator(address _account) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) SetOperator(_account common.Address) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.SetOperator(&_UpgradeManager1.TransactOpts, _account)
}

// SetVersion is a paid mutator transaction binding the contract method 0xb4062f8e.
//
// Solidity: function setVersion(bytes32 _hash, (string,uint256) _version) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) SetVersion(opts *bind.TransactOpts, _hash [32]byte, _version UpgradeManager1version) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "setVersion", _hash, _version)
}

// SetVersion is a paid mutator transaction binding the contract method 0xb4062f8e.
//
// Solidity: function setVersion(bytes32 _hash, (string,uint256) _version) returns()
func (_UpgradeManager1 *UpgradeManager1Session) SetVersion(_hash [32]byte, _version UpgradeManager1version) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.SetVersion(&_UpgradeManager1.TransactOpts, _hash, _version)
}

// SetVersion is a paid mutator transaction binding the contract method 0xb4062f8e.
//
// Solidity: function setVersion(bytes32 _hash, (string,uint256) _version) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) SetVersion(_hash [32]byte, _version UpgradeManager1version) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.SetVersion(&_UpgradeManager1.TransactOpts, _hash, _version)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) Upgrade(opts *bind.TransactOpts, _target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "upgrade", _target, _data)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager1 *UpgradeManager1Session) Upgrade(_target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.Upgrade(&_UpgradeManager1.TransactOpts, _target, _data)
}

// Upgrade is a paid mutator transaction binding the contract method 0x6e3d9ff0.
//
// Solidity: function upgrade(address _target, string _data) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) Upgrade(_target common.Address, _data string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.Upgrade(&_UpgradeManager1.TransactOpts, _target, _data)
}

// Upgrade0 is a paid mutator transaction binding the contract method 0x82f75182.
//
// Solidity: function upgrade(address _target, string _data, string _versionString) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) Upgrade0(opts *bind.TransactOpts, _target common.Address, _data string, _versionString string) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "upgrade0", _target, _data, _versionString)
}

// Upgrade0 is a paid mutator transaction binding the contract method 0x82f75182.
//
// Solidity: function upgrade(address _target, string _data, string _versionString) returns()
func (_UpgradeManager1 *UpgradeManager1Session) Upgrade0(_target common.Address, _data string, _versionString string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.Upgrade0(&_UpgradeManager1.TransactOpts, _target, _data, _versionString)
}

// Upgrade0 is a paid mutator transaction binding the contract method 0x82f75182.
//
// Solidity: function upgrade(address _target, string _data, string _versionString) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) Upgrade0(_target common.Address, _data string, _versionString string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.Upgrade0(&_UpgradeManager1.TransactOpts, _target, _data, _versionString)
}

// UpgradeMultiple is a paid mutator transaction binding the contract method 0x2011c815.
//
// Solidity: function upgradeMultiple(address[] _targets, string[] _bytecodes) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) UpgradeMultiple(opts *bind.TransactOpts, _targets []common.Address, _bytecodes []string) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "upgradeMultiple", _targets, _bytecodes)
}

// UpgradeMultiple is a paid mutator transaction binding the contract method 0x2011c815.
//
// Solidity: function upgradeMultiple(address[] _targets, string[] _bytecodes) returns()
func (_UpgradeManager1 *UpgradeManager1Session) UpgradeMultiple(_targets []common.Address, _bytecodes []string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.UpgradeMultiple(&_UpgradeManager1.TransactOpts, _targets, _bytecodes)
}

// UpgradeMultiple is a paid mutator transaction binding the contract method 0x2011c815.
//
// Solidity: function upgradeMultiple(address[] _targets, string[] _bytecodes) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) UpgradeMultiple(_targets []common.Address, _bytecodes []string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.UpgradeMultiple(&_UpgradeManager1.TransactOpts, _targets, _bytecodes)
}

// UpgradeMultiple0 is a paid mutator transaction binding the contract method 0x9d897130.
//
// Solidity: function upgradeMultiple(address[] _targets, string[] _bytecodes, string[] _versionStrings) returns()
func (_UpgradeManager1 *UpgradeManager1Transactor) UpgradeMultiple0(opts *bind.TransactOpts, _targets []common.Address, _bytecodes []string, _versionStrings []string) (*types.Transaction, error) {
	return _UpgradeManager1.contract.Transact(opts, "upgradeMultiple0", _targets, _bytecodes, _versionStrings)
}

// UpgradeMultiple0 is a paid mutator transaction binding the contract method 0x9d897130.
//
// Solidity: function upgradeMultiple(address[] _targets, string[] _bytecodes, string[] _versionStrings) returns()
func (_UpgradeManager1 *UpgradeManager1Session) UpgradeMultiple0(_targets []common.Address, _bytecodes []string, _versionStrings []string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.UpgradeMultiple0(&_UpgradeManager1.TransactOpts, _targets, _bytecodes, _versionStrings)
}

// UpgradeMultiple0 is a paid mutator transaction binding the contract method 0x9d897130.
//
// Solidity: function upgradeMultiple(address[] _targets, string[] _bytecodes, string[] _versionStrings) returns()
func (_UpgradeManager1 *UpgradeManager1TransactorSession) UpgradeMultiple0(_targets []common.Address, _bytecodes []string, _versionStrings []string) (*types.Transaction, error) {
	return _UpgradeManager1.Contract.UpgradeMultiple0(&_UpgradeManager1.TransactOpts, _targets, _bytecodes, _versionStrings)
}

// UpgradeManager1ConfigUpdateAddressIterator is returned from FilterConfigUpdateAddress and is used to iterate over the raw logs and unpacked data for ConfigUpdateAddress events raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateAddressIterator struct {
	Event *UpgradeManager1ConfigUpdateAddress // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1ConfigUpdateAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1ConfigUpdateAddress)
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
		it.Event = new(UpgradeManager1ConfigUpdateAddress)
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
func (it *UpgradeManager1ConfigUpdateAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1ConfigUpdateAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1ConfigUpdateAddress represents a ConfigUpdateAddress event raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateAddress struct {
	Name            string
	OldValue        common.Address
	NewValue        common.Address
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateAddress is a free log retrieval operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterConfigUpdateAddress(opts *bind.FilterOpts) (*UpgradeManager1ConfigUpdateAddressIterator, error) {

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1ConfigUpdateAddressIterator{contract: _UpgradeManager1.contract, event: "ConfigUpdateAddress", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateAddress is a free log subscription operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchConfigUpdateAddress(opts *bind.WatchOpts, sink chan<- *UpgradeManager1ConfigUpdateAddress) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "ConfigUpdateAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1ConfigUpdateAddress)
				if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
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

// ParseConfigUpdateAddress is a log parse operation binding the contract event 0xe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a0.
//
// Solidity: event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseConfigUpdateAddress(log types.Log) (*UpgradeManager1ConfigUpdateAddress, error) {
	event := new(UpgradeManager1ConfigUpdateAddress)
	if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManager1ConfigUpdateBoolIterator is returned from FilterConfigUpdateBool and is used to iterate over the raw logs and unpacked data for ConfigUpdateBool events raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateBoolIterator struct {
	Event *UpgradeManager1ConfigUpdateBool // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1ConfigUpdateBoolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1ConfigUpdateBool)
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
		it.Event = new(UpgradeManager1ConfigUpdateBool)
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
func (it *UpgradeManager1ConfigUpdateBoolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1ConfigUpdateBoolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1ConfigUpdateBool represents a ConfigUpdateBool event raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateBool struct {
	Name            string
	OldValue        bool
	NewValue        bool
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateBool is a free log retrieval operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterConfigUpdateBool(opts *bind.FilterOpts) (*UpgradeManager1ConfigUpdateBoolIterator, error) {

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1ConfigUpdateBoolIterator{contract: _UpgradeManager1.contract, event: "ConfigUpdateBool", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateBool is a free log subscription operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchConfigUpdateBool(opts *bind.WatchOpts, sink chan<- *UpgradeManager1ConfigUpdateBool) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "ConfigUpdateBool")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1ConfigUpdateBool)
				if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
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

// ParseConfigUpdateBool is a log parse operation binding the contract event 0x5edb308c5eddc69bcd31b4e689c5eed2fbd3155ae57915d1cad05425f6c1a39b.
//
// Solidity: event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseConfigUpdateBool(log types.Log) (*UpgradeManager1ConfigUpdateBool, error) {
	event := new(UpgradeManager1ConfigUpdateBool)
	if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateBool", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManager1ConfigUpdateIntIterator is returned from FilterConfigUpdateInt and is used to iterate over the raw logs and unpacked data for ConfigUpdateInt events raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateIntIterator struct {
	Event *UpgradeManager1ConfigUpdateInt // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1ConfigUpdateIntIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1ConfigUpdateInt)
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
		it.Event = new(UpgradeManager1ConfigUpdateInt)
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
func (it *UpgradeManager1ConfigUpdateIntIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1ConfigUpdateIntIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1ConfigUpdateInt represents a ConfigUpdateInt event raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateInt struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateInt is a free log retrieval operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterConfigUpdateInt(opts *bind.FilterOpts) (*UpgradeManager1ConfigUpdateIntIterator, error) {

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1ConfigUpdateIntIterator{contract: _UpgradeManager1.contract, event: "ConfigUpdateInt", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateInt is a free log subscription operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchConfigUpdateInt(opts *bind.WatchOpts, sink chan<- *UpgradeManager1ConfigUpdateInt) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "ConfigUpdateInt")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1ConfigUpdateInt)
				if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
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

// ParseConfigUpdateInt is a log parse operation binding the contract event 0xb5114472b89d1126433287defb6308bbabdb95b6ce5dd949ab2c151ed7b1ff4c.
//
// Solidity: event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseConfigUpdateInt(log types.Log) (*UpgradeManager1ConfigUpdateInt, error) {
	event := new(UpgradeManager1ConfigUpdateInt)
	if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateInt", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManager1ConfigUpdateUintIterator is returned from FilterConfigUpdateUint and is used to iterate over the raw logs and unpacked data for ConfigUpdateUint events raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateUintIterator struct {
	Event *UpgradeManager1ConfigUpdateUint // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1ConfigUpdateUintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1ConfigUpdateUint)
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
		it.Event = new(UpgradeManager1ConfigUpdateUint)
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
func (it *UpgradeManager1ConfigUpdateUintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1ConfigUpdateUintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1ConfigUpdateUint represents a ConfigUpdateUint event raised by the UpgradeManager1 contract.
type UpgradeManager1ConfigUpdateUint struct {
	Name            string
	OldValue        *big.Int
	NewValue        *big.Int
	AppliesAtHeight *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigUpdateUint is a free log retrieval operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterConfigUpdateUint(opts *bind.FilterOpts) (*UpgradeManager1ConfigUpdateUintIterator, error) {

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1ConfigUpdateUintIterator{contract: _UpgradeManager1.contract, event: "ConfigUpdateUint", logs: logs, sub: sub}, nil
}

// WatchConfigUpdateUint is a free log subscription operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchConfigUpdateUint(opts *bind.WatchOpts, sink chan<- *UpgradeManager1ConfigUpdateUint) (event.Subscription, error) {

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "ConfigUpdateUint")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1ConfigUpdateUint)
				if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
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

// ParseConfigUpdateUint is a log parse operation binding the contract event 0x207e45ce6f2191c3efffe5d0d91591cf81afd23cbf4bd746981d7a8a8dcfe1ba.
//
// Solidity: event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight)
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseConfigUpdateUint(log types.Log) (*UpgradeManager1ConfigUpdateUint, error) {
	event := new(UpgradeManager1ConfigUpdateUint)
	if err := _UpgradeManager1.contract.UnpackLog(event, "ConfigUpdateUint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeManager1UpgradeResultIterator is returned from FilterUpgradeResult and is used to iterate over the raw logs and unpacked data for UpgradeResult events raised by the UpgradeManager1 contract.
type UpgradeManager1UpgradeResultIterator struct {
	Event *UpgradeManager1UpgradeResult // Event containing the contract specifics and raw log

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
func (it *UpgradeManager1UpgradeResultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeManager1UpgradeResult)
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
		it.Event = new(UpgradeManager1UpgradeResult)
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
func (it *UpgradeManager1UpgradeResultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeManager1UpgradeResultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeManager1UpgradeResult represents a UpgradeResult event raised by the UpgradeManager1 contract.
type UpgradeManager1UpgradeResult struct {
	ContractAddress common.Address
	Success         bool
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpgradeResult is a free log retrieval operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager1 *UpgradeManager1Filterer) FilterUpgradeResult(opts *bind.FilterOpts, contractAddress []common.Address) (*UpgradeManager1UpgradeResultIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _UpgradeManager1.contract.FilterLogs(opts, "UpgradeResult", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeManager1UpgradeResultIterator{contract: _UpgradeManager1.contract, event: "UpgradeResult", logs: logs, sub: sub}, nil
}

// WatchUpgradeResult is a free log subscription operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager1 *UpgradeManager1Filterer) WatchUpgradeResult(opts *bind.WatchOpts, sink chan<- *UpgradeManager1UpgradeResult, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _UpgradeManager1.contract.WatchLogs(opts, "UpgradeResult", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeManager1UpgradeResult)
				if err := _UpgradeManager1.contract.UnpackLog(event, "UpgradeResult", log); err != nil {
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

// ParseUpgradeResult is a log parse operation binding the contract event 0x852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798.
//
// Solidity: event UpgradeResult(address indexed contractAddress, bool success)
func (_UpgradeManager1 *UpgradeManager1Filterer) ParseUpgradeResult(log types.Log) (*UpgradeManager1UpgradeResult, error) {
	event := new(UpgradeManager1UpgradeResult)
	if err := _UpgradeManager1.contract.UnpackLog(event, "UpgradeResult", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
