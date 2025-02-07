package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var UpgradeManagerBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b5060405161064738038061064783398101604081905261002f9161007c565b600080546001600160a01b039384166001600160a01b031991821617909155600180549290931691161790556100af565b80516001600160a01b038116811461007757600080fd5b919050565b6000806040838503121561008f57600080fd5b61009883610060565b91506100a660208401610060565b90509250929050565b610589806100be6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80636e3d9ff0146100515780637c8ccebe14610066578063b3ab15fb146100a9578063e7f43c68146100bc575b600080fd5b61006461005f3660046103f3565b6100da565b005b60005473ffffffffffffffffffffffffffffffffffffffff165b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b6100646100b73660046104d3565b610220565b60015473ffffffffffffffffffffffffffffffffffffffff16610080565b60015473ffffffffffffffffffffffffffffffffffffffff163314610160576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064015b60405180910390fd5b60405160f99060009061017990859085906020016104f5565b60405160208183030381529060405290506000806060600080855160208701885af4809350503d91506040519050602082018101604052818152816000602083013e8673ffffffffffffffffffffffffffffffffffffffff167f852faebd7b6261599c399849128319e84cd07abf61cfdea131405ee1d393479884604051610205911515815260200190565b60405180910390a282610219578160208201fd5b8160208201f35b60005473ffffffffffffffffffffffffffffffffffffffff1633146102c7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f63616c6c6572206973206e6f7420746865204175746f6e69747920636f6e747260448201527f61637400000000000000000000000000000000000000000000000000000000006064820152608401610157565b6001546040805160808082526008908201527f6f70657261746f7200000000000000000000000000000000000000000000000060a082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152436060830152517fe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a09181900360c00190a1600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b803573ffffffffffffffffffffffffffffffffffffffff811681146103bf57600080fd5b919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561040657600080fd5b61040f8361039b565b9150602083013567ffffffffffffffff8082111561042c57600080fd5b818501915085601f83011261044057600080fd5b813581811115610452576104526103c4565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715610498576104986103c4565b816040528281528860208487010111156104b157600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b6000602082840312156104e557600080fd5b6104ee8261039b565b9392505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008360601b1681526000825160005b818110156105415760208186018101516014868401015201610524565b5060009201601401918252509291505056fea2646970667358221220ed4e025342330119b0545c0b6a8a26a4c911a8943c8b57584d3f060fef39309464736f6c63430008160033")

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
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "appliesAtHeight",
            "type" : "uint256"
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
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "appliesAtHeight",
            "type" : "uint256"
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
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "appliesAtHeight",
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
