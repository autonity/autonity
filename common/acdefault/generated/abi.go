package generated

const Abi = `[
   {
      "inputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : "_participantAddress"
         },
         {
            "name" : "_participantEnode",
            "internalType" : "string[]",
            "type" : "string[]"
         },
         {
            "type" : "uint256[]",
            "name" : "_participantType",
            "internalType" : "uint256[]"
         },
         {
            "type" : "uint256[]",
            "name" : "_participantStake",
            "internalType" : "uint256[]"
         },
         {
            "internalType" : "uint256[]",
            "name" : "_commissionRate",
            "type" : "uint256[]"
         },
         {
            "name" : "_operatorAccount",
            "internalType" : "address",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_minGasPrice",
            "type" : "uint256"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_bondingPeriod"
         },
         {
            "type" : "uint256",
            "name" : "_committeeSize",
            "internalType" : "uint256"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_contractVersion"
         }
      ],
      "type" : "constructor",
      "stateMutability" : "nonpayable"
   },
   {
      "name" : "AddParticipant",
      "inputs" : [
         {
            "type" : "address",
            "indexed" : false,
            "internalType" : "address",
            "name" : "_address"
         },
         {
            "indexed" : false,
            "type" : "uint256",
            "name" : "_stake",
            "internalType" : "uint256"
         }
      ],
      "anonymous" : false,
      "type" : "event"
   },
   {
      "name" : "AddStakeholder",
      "inputs" : [
         {
            "name" : "_address",
            "internalType" : "address",
            "type" : "address",
            "indexed" : false
         },
         {
            "internalType" : "uint256",
            "name" : "_stake",
            "indexed" : false,
            "type" : "uint256"
         }
      ],
      "anonymous" : false,
      "type" : "event"
   },
   {
      "inputs" : [
         {
            "indexed" : false,
            "type" : "address",
            "name" : "_address",
            "internalType" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_stake",
            "indexed" : false,
            "type" : "uint256"
         }
      ],
      "name" : "AddValidator",
      "anonymous" : false,
      "type" : "event"
   },
   {
      "type" : "event",
      "name" : "BlockReward",
      "inputs" : [
         {
            "name" : "_address",
            "internalType" : "address",
            "indexed" : false,
            "type" : "address"
         },
         {
            "name" : "_amount",
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256"
         }
      ],
      "anonymous" : false
   },
   {
      "name" : "ChangeUserType",
      "inputs" : [
         {
            "indexed" : false,
            "type" : "address",
            "internalType" : "address",
            "name" : "_address"
         },
         {
            "internalType" : "enum Autonity.UserType",
            "name" : "_oldType",
            "indexed" : false,
            "type" : "uint8"
         },
         {
            "type" : "uint8",
            "indexed" : false,
            "internalType" : "enum Autonity.UserType",
            "name" : "_newType"
         }
      ],
      "anonymous" : false,
      "type" : "event"
   },
   {
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "type" : "address",
            "name" : "_address",
            "internalType" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256",
            "indexed" : false
         }
      ],
      "name" : "MintStake"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "type" : "address",
            "indexed" : false,
            "internalType" : "address",
            "name" : "_address"
         },
         {
            "indexed" : false,
            "type" : "uint256",
            "name" : "_amount",
            "internalType" : "uint256"
         }
      ],
      "name" : "RedeemStake",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "name" : "RemoveUser",
      "inputs" : [
         {
            "name" : "_address",
            "internalType" : "address",
            "type" : "address",
            "indexed" : false
         },
         {
            "type" : "uint8",
            "indexed" : false,
            "internalType" : "enum Autonity.UserType",
            "name" : "_type"
         }
      ],
      "type" : "event"
   },
   {
      "type" : "event",
      "inputs" : [
         {
            "name" : "_address",
            "internalType" : "address",
            "type" : "address",
            "indexed" : false
         },
         {
            "name" : "_value",
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256"
         }
      ],
      "name" : "SetCommissionRate",
      "anonymous" : false
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "type" : "uint256",
            "name" : "_gasPrice",
            "internalType" : "uint256"
         }
      ],
      "name" : "SetMinimumGasPrice",
      "type" : "event"
   },
   {
      "inputs" : [
         {
            "name" : "from",
            "internalType" : "address",
            "type" : "address",
            "indexed" : true
         },
         {
            "indexed" : true,
            "type" : "address",
            "internalType" : "address",
            "name" : "to"
         },
         {
            "type" : "uint256",
            "indexed" : false,
            "name" : "value",
            "internalType" : "uint256"
         }
      ],
      "name" : "Transfer",
      "anonymous" : false,
      "type" : "event"
   },
   {
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "type" : "string",
            "internalType" : "string",
            "name" : "version"
         }
      ],
      "name" : "Version"
   },
   {
      "stateMutability" : "payable",
      "type" : "fallback"
   },
   {
      "outputs" : [],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_address",
            "type" : "address"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_enode"
         }
      ],
      "name" : "addParticipant"
   },
   {
      "name" : "addStakeholder",
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_address",
            "type" : "address"
         },
         {
            "name" : "_enode",
            "internalType" : "string",
            "type" : "string"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_stake"
         }
      ],
      "type" : "function",
      "outputs" : [],
      "stateMutability" : "nonpayable"
   },
   {
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address payable",
            "name" : "_address"
         },
         {
            "internalType" : "uint256",
            "name" : "_stake",
            "type" : "uint256"
         },
         {
            "type" : "string",
            "name" : "_enode",
            "internalType" : "string"
         }
      ],
      "name" : "addValidator",
      "stateMutability" : "nonpayable",
      "outputs" : [],
      "type" : "function"
   },
   {
      "name" : "bondingPeriod",
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ]
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : [],
      "name" : "changeUserType",
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_address",
            "type" : "address"
         },
         {
            "type" : "uint8",
            "name" : "newUserType",
            "internalType" : "enum Autonity.UserType"
         }
      ]
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_account",
            "type" : "address"
         }
      ],
      "name" : "checkMember",
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "bool",
            "internalType" : "bool",
            "name" : ""
         }
      ],
      "type" : "function"
   },
   {
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : [],
      "name" : "committeeSize"
   },
   {
      "type" : "function",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "name" : "computeCommittee",
      "inputs" : []
   },
   {
      "inputs" : [],
      "name" : "contractVersion",
      "outputs" : [
         {
            "type" : "string",
            "name" : "",
            "internalType" : "string"
         }
      ],
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "name" : "deployer",
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "type" : "address",
            "name" : "",
            "internalType" : "address"
         }
      ],
      "stateMutability" : "view"
   },
   {
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "type" : "tuple",
            "components" : [
               {
                  "type" : "address[]",
                  "internalType" : "address[]",
                  "name" : "accounts"
               },
               {
                  "name" : "usertypes",
                  "internalType" : "enum Autonity.UserType[]",
                  "type" : "uint8[]"
               },
               {
                  "internalType" : "uint256[]",
                  "name" : "stakes",
                  "type" : "uint256[]"
               },
               {
                  "type" : "uint256[]",
                  "name" : "commissionrates",
                  "internalType" : "uint256[]"
               },
               {
                  "type" : "uint256",
                  "internalType" : "uint256",
                  "name" : "mingasprice"
               },
               {
                  "type" : "uint256",
                  "internalType" : "uint256",
                  "name" : "stakesupply"
               }
            ],
            "name" : "economics",
            "internalType" : "struct Autonity.EconomicsMetricData"
         }
      ],
      "name" : "dumpEconomicsMetricData",
      "inputs" : []
   },
   {
      "inputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "name" : "enodesWhitelist",
      "stateMutability" : "view",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "type" : "uint256",
            "name" : "_amount",
            "internalType" : "uint256"
         }
      ],
      "name" : "finalize",
      "outputs" : [
         {
            "type" : "bool",
            "internalType" : "bool",
            "name" : ""
         },
         {
            "internalType" : "struct Autonity.CommitteeMember[]",
            "name" : "",
            "components" : [
               {
                  "type" : "address",
                  "name" : "addr",
                  "internalType" : "address payable"
               },
               {
                  "type" : "uint256",
                  "name" : "votingPower",
                  "internalType" : "uint256"
               }
            ],
            "type" : "tuple[]"
         }
      ],
      "type" : "function",
      "stateMutability" : "nonpayable"
   },
   {
      "stateMutability" : "view",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "type" : "function",
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address",
            "name" : "_account"
         }
      ],
      "name" : "getAccountStake"
   },
   {
      "outputs" : [
         {
            "components" : [
               {
                  "type" : "address",
                  "internalType" : "address payable",
                  "name" : "addr"
               },
               {
                  "name" : "votingPower",
                  "internalType" : "uint256",
                  "type" : "uint256"
               }
            ],
            "type" : "tuple[]",
            "internalType" : "struct Autonity.CommitteeMember[]",
            "name" : ""
         }
      ],
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : [],
      "name" : "getCommittee"
   },
   {
      "name" : "getCurrentCommiteeSize",
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "type" : "uint256",
            "name" : "",
            "internalType" : "uint256"
         }
      ],
      "stateMutability" : "view"
   },
   {
      "name" : "getMaxCommitteeSize",
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "type" : "uint256",
            "name" : "",
            "internalType" : "uint256"
         }
      ]
   },
   {
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "view",
      "name" : "getMinimumGasPrice",
      "inputs" : []
   },
   {
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "view",
      "name" : "getRate",
      "inputs" : [
         {
            "type" : "address",
            "name" : "_account",
            "internalType" : "address"
         }
      ]
   },
   {
      "inputs" : [],
      "name" : "getStake",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "type" : "function",
      "outputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "name" : "getStakeholders",
      "inputs" : []
   },
   {
      "inputs" : [],
      "name" : "getValidators",
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : ""
         }
      ],
      "type" : "function"
   },
   {
      "name" : "getVersion",
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "stateMutability" : "view"
   },
   {
      "outputs" : [
         {
            "type" : "string[]",
            "name" : "",
            "internalType" : "string[]"
         }
      ],
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : [],
      "name" : "getWhitelist"
   },
   {
      "name" : "mintStake",
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_account",
            "type" : "address"
         },
         {
            "type" : "uint256",
            "name" : "_amount",
            "internalType" : "uint256"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : []
   },
   {
      "inputs" : [],
      "name" : "myUserType",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8"
         }
      ],
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "operatorAccount",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "address",
            "type" : "address"
         }
      ],
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
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_amount"
         }
      ],
      "name" : "redeemStake",
      "outputs" : [],
      "type" : "function",
      "stateMutability" : "nonpayable"
   },
   {
      "stateMutability" : "nonpayable",
      "outputs" : [],
      "type" : "function",
      "inputs" : [
         {
            "type" : "address",
            "name" : "_address",
            "internalType" : "address"
         }
      ],
      "name" : "removeUser"
   },
   {
      "inputs" : [],
      "name" : "retrieveContract",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         },
         {
            "name" : "",
            "internalType" : "string",
            "type" : "string"
         }
      ],
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : ""
         },
         {
            "name" : "",
            "internalType" : "string[]",
            "type" : "string[]"
         },
         {
            "internalType" : "uint256[]",
            "name" : "",
            "type" : "uint256[]"
         },
         {
            "name" : "",
            "internalType" : "uint256[]",
            "type" : "uint256[]"
         },
         {
            "name" : "",
            "internalType" : "uint256[]",
            "type" : "uint256[]"
         },
         {
            "name" : "",
            "internalType" : "address",
            "type" : "address"
         },
         {
            "type" : "uint256",
            "name" : "",
            "internalType" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         },
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : ""
         }
      ],
      "type" : "function",
      "inputs" : [],
      "name" : "retrieveState"
   },
   {
      "name" : "send",
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
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : [
         {
            "type" : "bool",
            "internalType" : "bool",
            "name" : ""
         }
      ]
   },
   {
      "type" : "function",
      "outputs" : [
         {
            "type" : "bool",
            "name" : "",
            "internalType" : "bool"
         }
      ],
      "stateMutability" : "nonpayable",
      "name" : "setCommissionRate",
      "inputs" : [
         {
            "name" : "_rate",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ]
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_size",
            "type" : "uint256"
         }
      ],
      "name" : "setCommitteeSize",
      "stateMutability" : "nonpayable",
      "outputs" : [],
      "type" : "function"
   },
   {
      "name" : "setMinimumGasPrice",
      "inputs" : [
         {
            "type" : "uint256",
            "name" : "_value",
            "internalType" : "uint256"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : []
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
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "outputs" : [
         {
            "name" : "",
            "internalType" : "bool",
            "type" : "bool"
         }
      ],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "string",
            "name" : "_bytecode",
            "type" : "string"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_abi"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_version"
         }
      ],
      "name" : "upgradeContract"
   },
   {
      "stateMutability" : "payable",
      "type" : "receive"
   }
]
`
