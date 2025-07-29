package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var SupplyControlBytecode = common.Hex2Bytes("608060405260405161091a38038061091a833981016040819052610022916100a5565b3460000361004357604051637c946ed760e01b815260040160405180910390fd5b600380546001600160a01b039485166001600160a01b031991821617909155600480549385169382169390931790925560018054919093169116179055346002556100e8565b80516001600160a01b03811681146100a057600080fd5b919050565b6000806000606084860312156100ba57600080fd5b6100c384610089565b92506100d160208501610089565b91506100df60408501610089565b90509250925092565b610823806100f76000396000f3fe6080604052600436106100705760003560e01c806380af17991161004e57806380af1799146100c7578063b3ab15fb146100fc578063c4e41b221461011c578063db7f521a1461013157600080fd5b806340c10f191461007557806344df8e70146100975780637ecc2b561461009f575b600080fd5b34801561008157600080fd5b506100956100903660046107a1565b610151565b005b610095610301565b3480156100ab57600080fd5b506100b46103d2565b6040519081526020015b60405180910390f35b3480156100d357600080fd5b5060015460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100be565b34801561010857600080fd5b506100956101173660046107cb565b610449565b34801561012857600080fd5b506100b461056e565b34801561013d57600080fd5b5061009561014c3660046107cb565b6105e2565b60015473ffffffffffffffffffffffffffffffffffffffff1633146101a2576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101aa610707565b73ffffffffffffffffffffffffffffffffffffffff821615806101e7575060015473ffffffffffffffffffffffffffffffffffffffff8381169116145b1561021e576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158061022a57504781115b15610261576040517f2c5211c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff83169082156108fc029083906000818181858888f193505050501580156102a4573d6000803e3d6000fd5b506040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a16102fd60008055565b5050565b3460000361033b576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015473ffffffffffffffffffffffffffffffffffffffff16331461038c576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610394610707565b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a16103d060008055565b565b60008054600103610444576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e637920646574656374656400000060448201526064015b60405180910390fd5b504790565b60035473ffffffffffffffffffffffffffffffffffffffff16331461049a576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004546040805160808082526008908201527f6f70657261746f7200000000000000000000000000000000000000000000000060a082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152436060830152517fe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a09181900360c00190a1600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b600080546001036105db576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e6379206465746563746564000000604482015260640161043b565b5060025490565b60035473ffffffffffffffffffffffffffffffffffffffff163314610633576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600154604080516080808252600a908201527f73746162696c697a65720000000000000000000000000000000000000000000060a082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152436060830152517fe821ac8084a7329d09d00cf1380cba50edea5f54a7bd453d49a9c354f565d4a09181900360c00190a1600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005415610771576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f7265656e7472616e637920646574656374656400000000000000000000000000604482015260640161043b565b6001600055565b803573ffffffffffffffffffffffffffffffffffffffff8116811461079c57600080fd5b919050565b600080604083850312156107b457600080fd5b6107bd83610778565b946020939093013593505050565b6000602082840312156107dd57600080fd5b6107e682610778565b939250505056fea2646970667358221220c61557eb263a1f5ad4980cfe26a82ee3a5ea57ba11b2111e4875dec518a285c364736f6c63430008160033")

var SupplyControlAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "autonity",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "operator",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "stabilizer_",
            "type" : "address"
         }
      ],
      "stateMutability" : "payable",
      "type" : "constructor"
   },
   {
      "inputs" : [],
      "name" : "InvalidAmount",
      "type" : "error"
   },
   {
      "inputs" : [],
      "name" : "InvalidRecipient",
      "type" : "error"
   },
   {
      "inputs" : [],
      "name" : "Unauthorized",
      "type" : "error"
   },
   {
      "inputs" : [],
      "name" : "ZeroValue",
      "type" : "error"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         }
      ],
      "name" : "Burn",
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
            "internalType" : "bool",
            "name" : "oldValue",
            "type" : "bool"
         },
         {
            "indexed" : false,
            "internalType" : "bool",
            "name" : "newValue",
            "type" : "bool"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "appliesAtHeight",
            "type" : "uint256"
         }
      ],
      "name" : "ConfigUpdateBool",
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
            "indexed" : false,
            "internalType" : "address",
            "name" : "recipient",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         }
      ],
      "name" : "Mint",
      "type" : "event"
   },
   {
      "inputs" : [],
      "name" : "availableSupply",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "burn",
      "outputs" : [],
      "stateMutability" : "payable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getStabilizer",
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
      "name" : "getTotalSupply",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "recipient",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         }
      ],
      "name" : "mint",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "operator",
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
            "name" : "stabilizer_",
            "type" : "address"
         }
      ],
      "name" : "setStabilizer",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   }
]
`))
