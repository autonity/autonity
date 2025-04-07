package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var SupplyControlBytecode = common.Hex2Bytes("60806040526040516109ad3803806109ad833981016040819052610022916100a5565b3460000361004357604051637c946ed760e01b815260040160405180910390fd5b600380546001600160a01b039485166001600160a01b031991821617909155600480549385169382169390931790925560018054919093169116179055346002556100e8565b80516001600160a01b03811681146100a057600080fd5b919050565b6000806000606084860312156100ba57600080fd5b6100c384610089565b92506100d160208501610089565b91506100df60408501610089565b90509250925092565b6108b6806100f76000396000f3fe6080604052600436106100705760003560e01c806380af17991161004e57806380af1799146100c7578063b3ab15fb14610101578063c4e41b2214610121578063db7f521a1461013657600080fd5b806340c10f191461007557806344df8e70146100975780637ecc2b561461009f575b600080fd5b34801561008157600080fd5b50610095610090366004610834565b610156565b005b610095610306565b3480156100ab57600080fd5b506100b46103d7565b6040519081526020015b60405180910390f35b3480156100d357600080fd5b506100dc61044e565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100be565b34801561010d57600080fd5b5061009561011c36600461085e565b6104d8565b34801561012d57600080fd5b506100b46105ff565b34801561014257600080fd5b5061009561015136600461085e565b610673565b60015473ffffffffffffffffffffffffffffffffffffffff1633146101a7576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101af61079a565b73ffffffffffffffffffffffffffffffffffffffff821615806101ec575060015473ffffffffffffffffffffffffffffffffffffffff8381169116145b15610223576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158061022f57504781115b15610266576040517f2c5211c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff83169082156108fc029083906000818181858888f193505050501580156102a9573d6000803e3d6000fd5b506040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a161030260008055565b5050565b34600003610340576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015473ffffffffffffffffffffffffffffffffffffffff163314610391576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61039961079a565b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a16103d560008055565b565b60008054600103610449576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e637920646574656374656400000060448201526064015b60405180910390fd5b504790565b600080546001036104bb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610440565b5060015473ffffffffffffffffffffffffffffffffffffffff1690565b60035473ffffffffffffffffffffffffffffffffffffffff163314610529576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61053161079a565b6004546040805160608082526008908201527f6f70657261746f72000000000000000000000000000000000000000000000000608082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152517f47f301f3da43f9b6e7d7fcbbc0d8a337b669877e8f0cb201934a64a196c94f429181900360a00190a1600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83161790556000805550565b6000805460010361066c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f72656164206f6e6c79207265656e7472616e63792064657465637465640000006044820152606401610440565b5060025490565b60035473ffffffffffffffffffffffffffffffffffffffff1633146106c4576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6106cc61079a565b600154604080516060808252600a908201527f73746162696c697a657200000000000000000000000000000000000000000000608082015273ffffffffffffffffffffffffffffffffffffffff928316602082015291831682820152517f47f301f3da43f9b6e7d7fcbbc0d8a337b669877e8f0cb201934a64a196c94f429181900360a00190a1600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83161790556000805550565b60005415610804576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f7265656e7472616e6379206465746563746564000000000000000000000000006044820152606401610440565b6001600055565b803573ffffffffffffffffffffffffffffffffffffffff8116811461082f57600080fd5b919050565b6000806040838503121561084757600080fd5b6108508361080b565b946020939093013593505050565b60006020828403121561087057600080fd5b6108798261080b565b939250505056fea2646970667358221220d837fa09358f3b4be019b8688f21f2e22ccea421d895181993163f540a3fc93664736f6c63430008160033")

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
