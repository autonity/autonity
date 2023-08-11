package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var SupplyControlBytecode = common.Hex2Bytes("608060405260405161065f38038061065f8339810160408190526100229161008d565b345f0361004257604051637c946ed760e01b815260040160405180910390fd5b5f80546001600160a01b039384166001600160a01b03199182161790915560018054929093169116179055346002556100c5565b6001600160a01b038116811461008a575f80fd5b50565b5f806040838503121561009e575f80fd5b82516100a981610076565b60208401519092506100ba81610076565b809150509250929050565b61058d806100d25f395ff3fe608060405260043610610063575f3560e01c80637e47961c116100415780637e47961c146100b85780637ecc2b5614610109578063db7f521a1461011b575f80fd5b806318160ddd1461006757806340c10f191461008f57806344df8e70146100b0575b5f80fd5b348015610072575f80fd5b5061007c60025481565b6040519081526020015b60405180910390f35b34801561009a575f80fd5b506100ae6100a93660046104f0565b61013a565b005b6100ae6102d6565b3480156100c3575f80fd5b506001546100e49073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610086565b348015610114575f80fd5b504761007c565b348015610126575f80fd5b506100ae61013536600461051a565b610395565b60015473ffffffffffffffffffffffffffffffffffffffff16331461018b576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821615806101c8575060015473ffffffffffffffffffffffffffffffffffffffff8381169116145b156101ff576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158061020b57504781115b15610242576040517f2c5211c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff83169082156108fc029083905f818181858888f19350505050158015610282573d5f803e3d5ffd5b506040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a15050565b345f0361030f576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60015473ffffffffffffffffffffffffffffffffffffffff163314610360576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a1565b5f8054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e7f43c686040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103fd573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610421919061053c565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610485576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b73ffffffffffffffffffffffffffffffffffffffff811681146104ed575f80fd5b50565b5f8060408385031215610501575f80fd5b823561050c816104cc565b946020939093013593505050565b5f6020828403121561052a575f80fd5b8135610535816104cc565b9392505050565b5f6020828403121561054c575f80fd5b8151610535816104cc56fea264697066735822122051b79a12b1d9baee2b321e7c7db3a2bc28944e3a861323d5ab2b615ebff3a60d64736f6c63430008150033")

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
