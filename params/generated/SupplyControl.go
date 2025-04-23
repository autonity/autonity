package generated

import (
	"strings"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
)

var SupplyControlBytecode = common.Hex2Bytes("608060405260405161067f38038061067f833981016040819052610022916100a2565b345f0361004257604051637c946ed760e01b815260040160405180910390fd5b600280546001600160a01b039485166001600160a01b03199182161790915560038054938516938216939093179092555f8054919093169116179055346001556100e2565b80516001600160a01b038116811461009d575f80fd5b919050565b5f805f606084860312156100b4575f80fd5b6100bd84610087565b92506100cb60208501610087565b91506100d960408501610087565b90509250925092565b610590806100ef5f395ff3fe60806040526004361061006e575f3560e01c80637e47961c1161004c5780637e47961c146100c35780637ecc2b5614610113578063b3ab15fb14610125578063db7f521a14610144575f80fd5b806318160ddd1461007257806340c10f191461009a57806344df8e70146100bb575b5f80fd5b34801561007d575f80fd5b5061008760015481565b6040519081526020015b60405180910390f35b3480156100a5575f80fd5b506100b96100b4366004610512565b610163565b005b6100b96102fd565b3480156100ce575f80fd5b505f546100ee9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610091565b34801561011e575f80fd5b5047610087565b348015610130575f80fd5b506100b961013f36600461053a565b6103bb565b34801561014f575f80fd5b506100b961015e36600461053a565b610453565b5f5473ffffffffffffffffffffffffffffffffffffffff1633146101b3576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821615806101ef57505f5473ffffffffffffffffffffffffffffffffffffffff8381169116145b15610226576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158061023257504781115b15610269576040517f2c5211c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff83169082156108fc029083905f818181858888f193505050501580156102a9573d5f803e3d5ffd5b506040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a15050565b345f03610336576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5473ffffffffffffffffffffffffffffffffffffffff163314610386576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a1565b60025473ffffffffffffffffffffffffffffffffffffffff16331461040c576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60035473ffffffffffffffffffffffffffffffffffffffff1633146104a4576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b803573ffffffffffffffffffffffffffffffffffffffff8116811461050d575f80fd5b919050565b5f8060408385031215610523575f80fd5b61052c836104ea565b946020939093013593505050565b5f6020828403121561054a575f80fd5b610553826104ea565b939250505056fea26469706673582212201e2594abf9aa17057d879a563d4f536600d9790bb4e92c0b1ec04d6ad59d8bdf64736f6c63430008150033")

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
   },
   {
      "inputs" : [],
      "name" : "stabilizer",
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
      "name" : "totalSupply",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   }
]
`))
