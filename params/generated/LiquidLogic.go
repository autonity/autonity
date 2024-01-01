package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var LiquidLogicBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b5060405161146d38038061146d83398101604081905261002f91610054565b600080546001600160a01b0319166001600160a01b0392909216919091179055610084565b60006020828403121561006657600080fd5b81516001600160a01b038116811461007d57600080fd5b9392505050565b6113da806100936000396000f3fe6080604052600436106101445760003560e01c80635ea1d6f8116100c0578063949813b811610074578063a9059cbb11610059578063a9059cbb1461038d578063dd62ed3e146103ad578063fb489a7b1461040057600080fd5b8063949813b81461034d5780639dc29fac1461036d57600080fd5b806370a08231116100a557806370a08231146102ca5780637eee288d1461030d57806384955c881461032d57600080fd5b80635ea1d6f81461026257806361d027b31461027857600080fd5b806323b872dd116101175780632f2c3f2e116100fc5780632f2c3f2e14610217578063372500ab1461022d57806340c10f191461024257600080fd5b806323b872dd146101d7578063282d3fdf146101f757600080fd5b8063095ea7b31461014957806318160ddd1461017e578063187cf4d71461019d57806319fac8fd146101b5575b600080fd5b34801561015557600080fd5b50610169610164366004611229565b610408565b60405190151581526020015b60405180910390f35b34801561018a57600080fd5b506004545b604051908152602001610175565b3480156101a957600080fd5b5061018f633b9aca0081565b3480156101c157600080fd5b506101d56101d0366004611253565b61041f565b005b3480156101e357600080fd5b506101696101f236600461126c565b6104b6565b34801561020357600080fd5b506101d5610212366004611229565b6105f7565b34801561022357600080fd5b5061018f61271081565b34801561023957600080fd5b506101d5610772565b34801561024e57600080fd5b506101d561025d366004611229565b610829565b34801561026e57600080fd5b5061018f60095481565b34801561028457600080fd5b506008546102a59073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610175565b3480156102d657600080fd5b5061018f6102e53660046112a8565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604090205490565b34801561031957600080fd5b506101d5610328366004611229565b610913565b34801561033957600080fd5b5061018f6103483660046112a8565b610a70565b34801561035957600080fd5b5061018f6103683660046112a8565b610aab565b34801561037957600080fd5b506101d5610388366004611229565b610ae6565b34801561039957600080fd5b506101696103a8366004611229565b610bc8565b3480156103b957600080fd5b5061018f6103c83660046112ca565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260036020908152604080832093909416825291909152205490565b61018f610c34565b6000610415338484610e0c565b5060015b92915050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146104b15760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b600955565b73ffffffffffffffffffffffffffffffffffffffff831660009081526003602090815260408083203384529091528120548281101561055d5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206160448201527f6c6c6f77616e636500000000000000000000000000000000000000000000000060648201526084016104a8565b610571853361056c868561132c565b610e0c565b61057b8584610f8b565b610585848461109a565b8373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040516105e491815260200190565b60405180910390a3506001949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146106845760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016104a8565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602090815260408083205460019092529091205482916106c19161132c565b10156107345760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201527f61626c650000000000000000000000000000000000000000000000000000000060648201526084016104a8565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600260205260408120805483929061076990849061133f565b90915550505050565b600061077d336110f2565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d80600081146107cf576040519150601f19603f3d011682016040523d82523d6000602084013e6107d4565b606091505b50509050806108255760405162461bcd60e51b815260206004820152601460248201527f4661696c656420746f2073656e6420457468657200000000000000000000000060448201526064016104a8565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108b65760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016104a8565b6108c0828261109a565b60405181815273ffffffffffffffffffffffffffffffffffffffff8316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef906020015b60405180910390a35050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146109a05760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016104a8565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260026020526040902054811115610a3b5760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f6360448201527f6b6564000000000000000000000000000000000000000000000000000000000060648201526084016104a8565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600260205260408120805483929061076990849061132c565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600260209081526040808320546001909252822054610419919061132c565b6000610ab682611171565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260056020526040902054610419919061133f565b60005473ffffffffffffffffffffffffffffffffffffffff163314610b735760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016104a8565b610b7d8282610f8b565b60405181815260009073ffffffffffffffffffffffffffffffffffffffff8416907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610907565b6000610bd43383610f8b565b610bde838361109a565b60405182815273ffffffffffffffffffffffffffffffffffffffff84169033907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a350600192915050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314610cc25760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016104a8565b600954349060009061271090610cd89084611352565b610ce29190611369565b905081811115610d345760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f7220726577617264000000000000000060448201526064016104a8565b610d3e818361132c565b60085460405191935073ffffffffffffffffffffffffffffffffffffffff16906108fc9083906000818181858888f193505050503d8060008114610d9e576040519150601f19603f3d011682016040523d82523d6000602084013e610da3565b606091505b505060045460009150610dba633b9aca0085611352565b610dc49190611369565b905080600754610dd4919061133f565b600755600454600090633b9aca0090610ded9084611352565b610df79190611369565b9050610e03818461133f565b94505050505090565b73ffffffffffffffffffffffffffffffffffffffff8316610e945760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f726573730000000000000000000000000000000000000000000000000000000060648201526084016104a8565b73ffffffffffffffffffffffffffffffffffffffff8216610f1d5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f737300000000000000000000000000000000000000000000000000000000000060648201526084016104a8565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610f94826110f2565b5073ffffffffffffffffffffffffffffffffffffffff8216600090815260016020908152604080832054600290925290912054610fd1908261132c565b8211156110205760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e6473000000000060448201526064016104a8565b61102a828261132c565b73ffffffffffffffffffffffffffffffffffffffff841660009081526001602052604090205580820361107e5773ffffffffffffffffffffffffffffffffffffffff83166000908152600660205260408120555b8160046000828254611090919061132c565b9091555050505050565b6110a3826110f2565b5073ffffffffffffffffffffffffffffffffffffffff8216600090815260016020526040812080548392906110d990849061133f565b925050819055508060046000828254610769919061133f565b6000806110fe83611171565b73ffffffffffffffffffffffffffffffffffffffff841660009081526005602052604090205490915061113290829061133f565b73ffffffffffffffffffffffffffffffffffffffff90931660009081526005602090815260408083208690556007546006909252909120555090919050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600160205260408120548082036111a75750600092915050565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600660205260408120546007546111da919061132c565b90506000633b9aca006111ed8484611352565b6111f79190611369565b95945050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461122457600080fd5b919050565b6000806040838503121561123c57600080fd5b61124583611200565b946020939093013593505050565b60006020828403121561126557600080fd5b5035919050565b60008060006060848603121561128157600080fd5b61128a84611200565b925061129860208501611200565b9150604084013590509250925092565b6000602082840312156112ba57600080fd5b6112c382611200565b9392505050565b600080604083850312156112dd57600080fd5b6112e683611200565b91506112f460208401611200565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b81810381811115610419576104196112fd565b80820180821115610419576104196112fd565b8082028115828204841417610419576104196112fd565b60008261139f577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea26469706673582212205b66ef255e9691b068bca62c5c1ccdb180030a64cda8552e187b8d8cce36462b64736f6c63430008170033")

var LiquidLogicAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_autonity",
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
            "indexed" : true,
            "internalType" : "address",
            "name" : "owner",
            "type" : "address"
         },
         {
            "indexed" : true,
            "internalType" : "address",
            "name" : "spender",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "value",
            "type" : "uint256"
         }
      ],
      "name" : "Approval",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : true,
            "internalType" : "address",
            "name" : "from",
            "type" : "address"
         },
         {
            "indexed" : true,
            "internalType" : "address",
            "name" : "to",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "value",
            "type" : "uint256"
         }
      ],
      "name" : "Transfer",
      "type" : "event"
   },
   {
      "inputs" : [],
      "name" : "COMMISSION_RATE_PRECISION",
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
      "name" : "FEE_FACTOR_UNIT_RECIP",
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
            "name" : "_owner",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_spender",
            "type" : "address"
         }
      ],
      "name" : "allowance",
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
            "name" : "_spender",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "approve",
      "outputs" : [
         {
            "internalType" : "bool",
            "name" : "",
            "type" : "bool"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_delegator",
            "type" : "address"
         }
      ],
      "name" : "balanceOf",
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
            "name" : "_account",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "burn",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "claimRewards",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "commissionRate",
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
            "name" : "_account",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "lock",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_account",
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
      "inputs" : [],
      "name" : "redistribute",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "payable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_rate",
            "type" : "uint256"
         }
      ],
      "name" : "setCommissionRate",
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
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_to",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "transfer",
      "outputs" : [
         {
            "internalType" : "bool",
            "name" : "_success",
            "type" : "bool"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_sender",
            "type" : "address"
         },
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
      "name" : "transferFrom",
      "outputs" : [
         {
            "internalType" : "bool",
            "name" : "_success",
            "type" : "bool"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "treasury",
      "outputs" : [
         {
            "internalType" : "address payable",
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
      "name" : "unclaimedRewards",
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
            "name" : "_account",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "unlock",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_delegator",
            "type" : "address"
         }
      ],
      "name" : "unlockedBalanceOf",
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
