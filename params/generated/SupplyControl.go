package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var SupplyControlBytecode = common.Hex2Bytes("608060405260405161069e38038061069e833981016040819052610022916100a5565b3460000361004357604051637c946ed760e01b815260040160405180910390fd5b600280546001600160a01b039485166001600160a01b031991821617909155600380549385169382169390931790925560008054919093169116179055346001556100e8565b80516001600160a01b03811681146100a057600080fd5b919050565b6000806000606084860312156100ba57600080fd5b6100c384610089565b92506100d160208501610089565b91506100df60408501610089565b90509250925092565b6105a7806100f76000396000f3fe6080604052600436106100705760003560e01c80637e47961c1161004e5780637e47961c146100c85780637ecc2b561461011a578063b3ab15fb1461012d578063db7f521a1461014d57600080fd5b806318160ddd1461007557806340c10f191461009e57806344df8e70146100c0575b600080fd5b34801561008157600080fd5b5061008b60015481565b6040519081526020015b60405180910390f35b3480156100aa57600080fd5b506100be6100b9366004610525565b61016d565b005b6100be61030c565b3480156100d457600080fd5b506000546100f59073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610095565b34801561012657600080fd5b504761008b565b34801561013957600080fd5b506100be61014836600461054f565b6103cc565b34801561015957600080fd5b506100be61016836600461054f565b610464565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101be576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff821615806101fb575060005473ffffffffffffffffffffffffffffffffffffffff8381169116145b15610232576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158061023e57504781115b15610275576040517f2c5211c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff83169082156108fc029083906000818181858888f193505050501580156102b8573d6000803e3d6000fd5b506040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a15050565b34600003610346576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005473ffffffffffffffffffffffffffffffffffffffff163314610397576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a1565b60025473ffffffffffffffffffffffffffffffffffffffff16331461041d576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60035473ffffffffffffffffffffffffffffffffffffffff1633146104b5576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b803573ffffffffffffffffffffffffffffffffffffffff8116811461052057600080fd5b919050565b6000806040838503121561053857600080fd5b610541836104fc565b946020939093013593505050565b60006020828403121561056157600080fd5b61056a826104fc565b939250505056fea2646970667358221220da79857220355f97f29cc6b6d0da24a4c5e2c6c44592025346a088c81bcd132d64736f6c63430008150033")

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
