package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var LiquidLogicBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b5060405161134238038061134283398101604081905261002f91610054565b600080546001600160a01b0319166001600160a01b0392909216919091179055610084565b60006020828403121561006657600080fd5b81516001600160a01b038116811461007d57600080fd5b9392505050565b6112af806100936000396000f3fe6080604052600436106100f35760003560e01c80635ea1d6f81161008a578063949813b811610059578063949813b8146102a85780639dc29fac146102c8578063a9059cbb146102e8578063fb489a7b1461030857600080fd5b80635ea1d6f81461020057806361d027b3146102165780637eee288d1461026857806384955c881461028857600080fd5b8063282d3fdf116100c6578063282d3fdf146101955780632f2c3f2e146101b5578063372500ab146101cb57806340c10f19146101e057600080fd5b8063095ea7b3146100f8578063187cf4d71461012d57806319fac8fd1461015357806323b872dd14610175575b600080fd5b34801561010457600080fd5b50610118610113366004611131565b610310565b60405190151581526020015b60405180910390f35b34801561013957600080fd5b50610145633b9aca0081565b604051908152602001610124565b34801561015f57600080fd5b5061017361016e36600461115b565b610327565b005b34801561018157600080fd5b50610118610190366004611174565b6103be565b3480156101a157600080fd5b506101736101b0366004611131565b6104ff565b3480156101c157600080fd5b5061014561271081565b3480156101d757600080fd5b5061017361067a565b3480156101ec57600080fd5b506101736101fb366004611131565b610731565b34801561020c57600080fd5b5061014560095481565b34801561022257600080fd5b506008546102439073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610124565b34801561027457600080fd5b50610173610283366004611131565b61081b565b34801561029457600080fd5b506101456102a33660046111b0565b610978565b3480156102b457600080fd5b506101456102c33660046111b0565b6109b3565b3480156102d457600080fd5b506101736102e3366004611131565b6109ee565b3480156102f457600080fd5b50610118610303366004611131565b610ad0565b610145610b3c565b600061031d338484610d14565b5060015b92915050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146103b95760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b600955565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600360209081526040808320338452909152812054828110156104655760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e742065786365656473206160448201527f6c6c6f77616e636500000000000000000000000000000000000000000000000060648201526084016103b0565b61047985336104748685611201565b610d14565b6104838584610e93565b61048d8484610fa2565b8373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef856040516104ec91815260200190565b60405180910390a3506001949350505050565b60005473ffffffffffffffffffffffffffffffffffffffff16331461058c5760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016103b0565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602090815260408083205460019092529091205482916105c991611201565b101561063c5760405162461bcd60e51b8152602060048201526024808201527f63616e2774206c6f636b206d6f72652066756e6473207468616e20617661696c60448201527f61626c650000000000000000000000000000000000000000000000000000000060648201526084016103b0565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081208054839290610671908490611214565b90915550505050565b600061068533610ffa565b33600081815260056020526040808220829055519293509183908381818185875af1925050503d80600081146106d7576040519150601f19603f3d011682016040523d82523d6000602084013e6106dc565b606091505b505090508061072d5760405162461bcd60e51b815260206004820152601460248201527f4661696c656420746f2073656e6420457468657200000000000000000000000060448201526064016103b0565b5050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146107be5760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016103b0565b6107c88282610fa2565b60405181815273ffffffffffffffffffffffffffffffffffffffff8316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef906020015b60405180910390a35050565b60005473ffffffffffffffffffffffffffffffffffffffff1633146108a85760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016103b0565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600260205260409020548111156109435760405162461bcd60e51b815260206004820152602360248201527f63616e277420756e6c6f636b206d6f72652066756e6473207468616e206c6f6360448201527f6b6564000000000000000000000000000000000000000000000000000000000060648201526084016103b0565b73ffffffffffffffffffffffffffffffffffffffff821660009081526002602052604081208054839290610671908490611201565b73ffffffffffffffffffffffffffffffffffffffff811660009081526002602090815260408083205460019092528220546103219190611201565b60006109be82611079565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600560205260409020546103219190611214565b60005473ffffffffffffffffffffffffffffffffffffffff163314610a7b5760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016103b0565b610a858282610e93565b60405181815260009073ffffffffffffffffffffffffffffffffffffffff8416907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200161080f565b6000610adc3383610e93565b610ae68383610fa2565b60405182815273ffffffffffffffffffffffffffffffffffffffff84169033907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a350600192915050565b6000805473ffffffffffffffffffffffffffffffffffffffff163314610bca5760405162461bcd60e51b815260206004820152602860248201527f43616c6c207265737472696374656420746f20746865204175746f6e6974792060448201527f436f6e747261637400000000000000000000000000000000000000000000000060648201526084016103b0565b600954349060009061271090610be09084611227565b610bea919061123e565b905081811115610c3c5760405162461bcd60e51b815260206004820152601860248201527f696e76616c69642076616c696461746f7220726577617264000000000000000060448201526064016103b0565b610c468183611201565b60085460405191935073ffffffffffffffffffffffffffffffffffffffff16906108fc9083906000818181858888f193505050503d8060008114610ca6576040519150601f19603f3d011682016040523d82523d6000602084013e610cab565b606091505b505060045460009150610cc2633b9aca0085611227565b610ccc919061123e565b905080600754610cdc9190611214565b600755600454600090633b9aca0090610cf59084611227565b610cff919061123e565b9050610d0b8184611214565b94505050505090565b73ffffffffffffffffffffffffffffffffffffffff8316610d9c5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460448201527f726573730000000000000000000000000000000000000000000000000000000060648201526084016103b0565b73ffffffffffffffffffffffffffffffffffffffff8216610e255760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f20616464726560448201527f737300000000000000000000000000000000000000000000000000000000000060648201526084016103b0565b73ffffffffffffffffffffffffffffffffffffffff83811660008181526003602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b610e9c82610ffa565b5073ffffffffffffffffffffffffffffffffffffffff8216600090815260016020908152604080832054600290925290912054610ed99082611201565b821115610f285760405162461bcd60e51b815260206004820152601b60248201527f696e73756666696369656e7420756e6c6f636b65642066756e6473000000000060448201526064016103b0565b610f328282611201565b73ffffffffffffffffffffffffffffffffffffffff8416600090815260016020526040902055808203610f865773ffffffffffffffffffffffffffffffffffffffff83166000908152600660205260408120555b8160046000828254610f989190611201565b9091555050505050565b610fab82610ffa565b5073ffffffffffffffffffffffffffffffffffffffff821660009081526001602052604081208054839290610fe1908490611214565b9250508190555080600460008282546106719190611214565b60008061100683611079565b73ffffffffffffffffffffffffffffffffffffffff841660009081526005602052604090205490915061103a908290611214565b73ffffffffffffffffffffffffffffffffffffffff90931660009081526005602090815260408083208690556007546006909252909120555090919050565b73ffffffffffffffffffffffffffffffffffffffff81166000908152600160205260408120548082036110af5750600092915050565b73ffffffffffffffffffffffffffffffffffffffff83166000908152600660205260408120546007546110e29190611201565b90506000633b9aca006110f58484611227565b6110ff919061123e565b95945050505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461112c57600080fd5b919050565b6000806040838503121561114457600080fd5b61114d83611108565b946020939093013593505050565b60006020828403121561116d57600080fd5b5035919050565b60008060006060848603121561118957600080fd5b61119284611108565b92506111a060208501611108565b9150604084013590509250925092565b6000602082840312156111c257600080fd5b6111cb82611108565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b81810381811115610321576103216111d2565b80820180821115610321576103216111d2565b8082028115828204841417610321576103216111d2565b600082611274577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b50049056fea26469706673582212208d59d9062f607a7150e7e519c283130154a58028cd535732dca868ce004e171e64736f6c63430008170033")

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
