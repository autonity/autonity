package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var UpgradeManagerBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b506040516107ae3803806107ae83398101604081905261002f9161007c565b600180546001600160a01b039384166001600160a01b031991821617909155600280549290931691161790556100af565b80516001600160a01b038116811461007757600080fd5b919050565b6000806040838503121561008f57600080fd5b61009883610060565b91506100a660208401610060565b90509250929050565b6106f0806100be6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80636e3d9ff0146100515780637c8ccebe14610066578063b3ab15fb14610097578063e7f43c68146100aa575b600080fd5b61006461005f36600461055a565b6100b2565b005b61006e610200565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b6100646100a536600461063a565b61028a565b61006e610407565b6100ba610491565b60025473ffffffffffffffffffffffffffffffffffffffff163314610140576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064015b60405180910390fd5b60405160f990600090610159908590859060200161065c565b60405160208183030381529060405290506000806060600080855160208701885af4809350503d91506040519050602082018101604052818152816000602083013e8673ffffffffffffffffffffffffffffffffffffffff167f852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d3934798846040516101e5911515815260200190565b60405180910390a2826101f9578160208201fd5b8160208201f35b6000805460010361026d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610137565b5060015473ffffffffffffffffffffffffffffffffffffffff1690565b610292610491565b60015473ffffffffffffffffffffffffffffffffffffffff163314610339576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f63616c6c6572206973206e6f7420746865204175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610137565b6002546040805160608082526008908201527f6f70657261746f72000000000000000000000000000000000000000000000000608082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152517f47f301f3da43f9b6e7d7fcbbc0d8a337b669877e8f0cb201934a64a196c94f429181900360a00190a1600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83161790556000805550565b60008054600103610474576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610137565b5060025473ffffffffffffffffffffffffffffffffffffffff1690565b600054156104fb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f7265656e7472616e6379206465746563746564000000000000000000000000006044820152606401610137565b6001600055565b803573ffffffffffffffffffffffffffffffffffffffff8116811461052657600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561056d57600080fd5b61057683610502565b9150602083013567ffffffffffffffff8082111561059357600080fd5b818501915085601f8301126105a757600080fd5b8135818111156105b9576105b961052b565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f011681019083821181831017156105ff576105ff61052b565b8160405282815288602084870101111561061857600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60006020828403121561064c57600080fd5b61065582610502565b9392505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008360601b1681526000825160005b818110156106a8576020818601810151601486840101520161068b565b5060009201601401918252509291505056fea2646970667358221220a0723fea6c5534c1a7f3f5b9139f7ac4f721bf771e329a7791ff1d523045973e64736f6c63430008160033")

var UpgradeManagerAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_autonity",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_operator",
            "type" : "address"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "constructor"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "string",
            "name" : "name",
            "type" : "string"
         },
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "oldValue",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "newValue",
            "type" : "address"
         }
      ],
      "name" : "ConfigUpdateAddress",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "string",
            "name" : "name",
            "type" : "string"
         },
         {
            "indexed" : false,
            "internalType" : "int256",
            "name" : "oldValue",
            "type" : "int256"
         },
         {
            "indexed" : false,
            "internalType" : "int256",
            "name" : "newValue",
            "type" : "int256"
         }
      ],
      "name" : "ConfigUpdateInt",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "string",
            "name" : "name",
            "type" : "string"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "oldValue",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "newValue",
            "type" : "uint256"
         }
      ],
      "name" : "ConfigUpdateUint",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : true,
            "internalType" : "address",
            "name" : "contractAddress",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "bool",
            "name" : "success",
            "type" : "bool"
         }
      ],
      "name" : "UpgradeResult",
      "type" : "event"
   },
   {
      "inputs" : [],
      "name" : "getAutonity",
      "outputs" : [
         {
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getOperator",
      "outputs" : [
         {
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_account",
            "type" : "address"
         }
      ],
      "name" : "setOperator",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_target",
            "type" : "address"
         },
         {
            "internalType" : "string",
            "name" : "_data",
            "type" : "string"
         }
      ],
      "name" : "upgrade",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   }
]
`))
