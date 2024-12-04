package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var SupplyControlBytecode = common.Hex2Bytes("6080604052604051610775380380610775833981016040819052610022916100a8565b3460000361004357604051637c946ed760e01b815260040160405180910390fd5b600380546001600160a01b03199081166001600160a01b0396871617909155600480548216948616949094179093556000805490931691909316179055346001556002556100f3565b80516001600160a01b03811681146100a357600080fd5b919050565b600080600080608085870312156100be57600080fd5b6100c78561008c565b93506100d56020860161008c565b92506100e36040860161008c565b6060959095015193969295505050565b610673806101026000396000f3fe60806040526004361061007b5760003560e01c80637e47961c1161004e5780637e47961c146100e95780637ecc2b561461013b578063b3ab15fb1461014e578063db7f521a1461016e57600080fd5b806318160ddd1461008057806340c10f19146100a957806344df8e70146100cb57806359599203146100d3575b600080fd5b34801561008c57600080fd5b5061009660015481565b6040519081526020015b60405180910390f35b3480156100b557600080fd5b506100c96100c43660046105b1565b61018e565b005b6100c961032d565b3480156100df57600080fd5b5061009660025481565b3480156100f557600080fd5b506000546101169073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100a0565b34801561014757600080fd5b5047610096565b34801561015a57600080fd5b506100c96101693660046105db565b610458565b34801561017a57600080fd5b506100c96101893660046105db565b6104f0565b60005473ffffffffffffffffffffffffffffffffffffffff1633146101df576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff8216158061021c575060005473ffffffffffffffffffffffffffffffffffffffff8381169116145b15610253576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80158061025f57504781115b15610296576040517f2c5211c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405173ffffffffffffffffffffffffffffffffffffffff83169082156108fc029083906000818181858888f193505050501580156102d9573d6000803e3d6000fd5b506040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885910160405180910390a15050565b34600003610367576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005473ffffffffffffffffffffffffffffffffffffffff1633146103b8576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600254156104235760025434908111156103db57506002805460009091556103f3565b34600260008282546103ed91906105fd565b90915550505b60405160009082156108fc0290839083818181858288f19350505050158015610420573d6000803e3d6000fd5b50505b6040513481527fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb9060200160405180910390a1565b60035473ffffffffffffffffffffffffffffffffffffffff1633146104a9576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600480547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60045473ffffffffffffffffffffffffffffffffffffffff163314610541576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b803573ffffffffffffffffffffffffffffffffffffffff811681146105ac57600080fd5b919050565b600080604083850312156105c457600080fd5b6105cd83610588565b946020939093013593505050565b6000602082840312156105ed57600080fd5b6105f682610588565b9392505050565b81810381811115610637577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea26469706673582212207a41861153c1d05d803ab060551d1037b25925e3d192c927cf4d649faf7e0e0264736f6c63430008150033")

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
         },
         {
            "internalType" : "uint256",
            "name" : "burnableGenesisSupply_",
            "type" : "uint256"
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
      "inputs" : [],
      "name" : "burnableGenesisSupply",
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
