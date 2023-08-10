package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var SupplyControlBytecode = common.Hex2Bytes("608060405260405161089638038061089683398181016040528101906100259190610151565b6000340361005f576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550346002819055505050610191565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061011e826100f3565b9050919050565b61012e81610113565b811461013957600080fd5b50565b60008151905061014b81610125565b92915050565b60008060408385031215610168576101676100ee565b5b60006101768582860161013c565b92505060206101878582860161013c565b9150509250929050565b6106f6806101a06000396000f3fe6080604052600436106100555760003560e01c806318160ddd1461005a57806340c10f191461008557806344df8e70146100ae578063570ca735146100b85780637ecc2b56146100e3578063b3ab15fb1461010e575b600080fd5b34801561006657600080fd5b5061006f610137565b60405161007c9190610556565b60405180910390f35b34801561009157600080fd5b506100ac60048036038101906100a79190610600565b61013d565b005b6100b661034c565b005b3480156100c457600080fd5b506100cd610446565b6040516100da919061064f565b60405180910390f35b3480156100ef57600080fd5b506100f861046c565b6040516101059190610556565b60405180910390f35b34801561011a57600080fd5b506101356004803603810190610130919061066a565b610474565b005b60025481565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101c4576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16148061024c5750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16145b15610283576040517f9c8d2cd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600081148061029157504781115b156102c8576040517f2c5211c600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f1935050505015801561030e573d6000803e3d6000fd5b507f0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d41213968858282604051610340929190610697565b60405180910390a15050565b60003403610386576040517f7c946ed700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461040d576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7fb90306ad06b2a6ff86ddc9327db583062895ef6540e62dc50add009db5b356eb3460405161043c9190610556565b60405180910390a1565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600047905090565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104f9576040517f82b4290000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000819050919050565b6105508161053d565b82525050565b600060208201905061056b6000830184610547565b92915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006105a182610576565b9050919050565b6105b181610596565b81146105bc57600080fd5b50565b6000813590506105ce816105a8565b92915050565b6105dd8161053d565b81146105e857600080fd5b50565b6000813590506105fa816105d4565b92915050565b6000806040838503121561061757610616610571565b5b6000610625858286016105bf565b9250506020610636858286016105eb565b9150509250929050565b61064981610596565b82525050565b60006020820190506106646000830184610640565b92915050565b6000602082840312156106805761067f610571565b5b600061068e848285016105bf565b91505092915050565b60006040820190506106ac6000830185610640565b6106b96020830184610547565b939250505056fea2646970667358221220b0de66be4b9acd7e705b53b5d8fe51a4eb93a5a5b34d094feff4d4d2b6827c2564736f6c63430008130033")

var SupplyControlAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "admin",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "operator_",
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
      "inputs" : [],
      "name" : "operator",
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
            "name" : "operator_",
            "type" : "address"
         }
      ],
      "name" : "setOperator",
      "outputs" : [],
      "stateMutability" : "nonpayable",
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
