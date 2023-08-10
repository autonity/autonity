package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var SupplyControlBytecode = common.Hex2Bytes("60806040526040516106b23803806106b28339810160408190526100229161010e565b345f0361004257604051637c946ed760e01b815260040160405180910390fd5b600180546001600160a01b038085166001600160a01b03199283168117909355600280549185169190921617905560408051631cfe878d60e31b8152905163e7f43c68916004808201926020929091908290030181865afa1580156100a9573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906100cd9190610146565b5f80546001600160a01b0319166001600160a01b0392909216919091179055505034600355610168565b6001600160a01b038116811461010b575f80fd5b50565b5f806040838503121561011f575f80fd5b825161012a816100f7565b602084015190925061013b816100f7565b809150509250929050565b5f60208284031215610156575f80fd5b8151610161816100f7565b9392505050565b61053d806101755f395ff3fe60806040526004361061006e575f3560e01c80637e47961c1161004c5780637e47961c146100c35780637ecc2b5614610114578063b3ab15fb14610126578063db7f521a14610186575f80fd5b806318160ddd1461007257806340c10f191461009a57806344df8e70146100bb575b5f80fd5b34801561007d575f80fd5b5061008760035481565b6040519081526020015b60405180910390f35b3480156100a5575f80fd5b506100b96100b43660046104bf565b6101a5565b005b6100b9610341565b3480156100ce575f80fd5b506002546100ef9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610091565b34801561011f575f80fd5b5047610087565b348015610131575f80fd5b506100b96101403660046104e7565b5f80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b348015610191575f80fd5b506100b96101a03660046104e7565b610400565b60025473ffffffffffffffffffffffffffffffffffffffff1633146101f6576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff82161580610233575060025473ffffffffffffffffffffffffffffffffffffffff8381169116145b1561026a576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158061027657504781115b156102ad576040517f2c5211c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff83169082156108fc029083905f818181858888f193505050501580156102ed573d5f803e3d5ffd5b506040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a15050565b345f0361037a576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60025473ffffffffffffffffffffffffffffffffffffffff1633146103cb576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a1565b5f5473ffffffffffffffffffffffffffffffffffffffff163314610450576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b803573ffffffffffffffffffffffffffffffffffffffff811681146104ba575f80fd5b919050565b5f80604083850312156104d0575f80fd5b6104d983610497565b946020939093013593505050565b5f602082840312156104f7575f80fd5b61050082610497565b939250505056fea26469706673582212206e181310260cfb20e2108410f486b4a7f0b3b6771164dc1d4b973f98f8cdbef964736f6c63430008150033")

var SupplyControlAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_autonity",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_stabilizer",
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
            "name" : "_recipient",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
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
            "name" : "_operator",
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
            "name" : "_stabilizer",
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
