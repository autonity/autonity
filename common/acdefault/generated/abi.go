package generated

const Abi = `[
   {
      "type" : "constructor",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "address[]",
            "name" : "_participantAddress",
            "type" : "address[]"
         },
         {
            "internalType" : "string[]",
            "name" : "_participantEnode",
            "type" : "string[]"
         },
         {
            "internalType" : "uint256[]",
            "name" : "_participantType",
            "type" : "uint256[]"
         },
         {
            "internalType" : "uint256[]",
            "type" : "uint256[]",
            "name" : "_participantStake"
         },
         {
            "internalType" : "uint256[]",
            "type" : "uint256[]",
            "name" : "_commissionRate"
         },
         {
            "type" : "address",
            "name" : "_operatorAccount",
            "internalType" : "address"
         },
         {
            "name" : "_minGasPrice",
            "type" : "uint256",
            "internalType" : "uint256"
         },
         {
            "name" : "_bondingPeriod",
            "type" : "uint256",
            "internalType" : "uint256"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_committeeSize"
         },
         {
            "name" : "_contractVersion",
            "type" : "string",
            "internalType" : "string"
         }
      ]
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "indexed" : false,
            "name" : "_address",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_stake"
         }
      ],
      "anonymous" : false,
      "type" : "event",
      "name" : "AddParticipant"
   },
   {
      "name" : "AddStakeholder",
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "type" : "address",
            "name" : "_address",
            "indexed" : false,
            "internalType" : "address"
         },
         {
            "name" : "_stake",
            "type" : "uint256",
            "indexed" : false,
            "internalType" : "uint256"
         }
      ]
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "indexed" : false,
            "type" : "address",
            "name" : "_address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_stake"
         }
      ],
      "anonymous" : false,
      "name" : "AddValidator",
      "type" : "event"
   },
   {
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "_address",
            "type" : "address"
         },
         {
            "name" : "_amount",
            "type" : "uint256",
            "indexed" : false,
            "internalType" : "uint256"
         }
      ],
      "name" : "BlockReward",
      "type" : "event",
      "anonymous" : false
   },
   {
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "internalType" : "address",
            "indexed" : false
         },
         {
            "internalType" : "enum Autonity.UserType",
            "indexed" : false,
            "name" : "_oldType",
            "type" : "uint8"
         },
         {
            "type" : "uint8",
            "name" : "_newType",
            "internalType" : "enum Autonity.UserType",
            "indexed" : false
         }
      ],
      "name" : "ChangeUserType",
      "type" : "event",
      "anonymous" : false
   },
   {
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "address",
            "type" : "address",
            "name" : "_address"
         },
         {
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256",
            "name" : "_amount"
         }
      ],
      "name" : "MintStake",
      "type" : "event",
      "anonymous" : false
   },
   {
      "anonymous" : false,
      "type" : "event",
      "name" : "RedeemStake",
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "internalType" : "address",
            "indexed" : false
         },
         {
            "type" : "uint256",
            "name" : "_amount",
            "internalType" : "uint256",
            "indexed" : false
         }
      ]
   },
   {
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "indexed" : false,
            "internalType" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8",
            "name" : "_type"
         }
      ],
      "type" : "event",
      "name" : "RemoveUser",
      "anonymous" : false
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "indexed" : false,
            "type" : "address",
            "name" : "_address"
         },
         {
            "internalType" : "uint256",
            "indexed" : false,
            "name" : "_value",
            "type" : "uint256"
         }
      ],
      "anonymous" : false,
      "type" : "event",
      "name" : "SetCommissionRate"
   },
   {
      "inputs" : [
         {
            "type" : "uint256",
            "name" : "_gasPrice",
            "indexed" : false,
            "internalType" : "uint256"
         }
      ],
      "type" : "event",
      "name" : "SetMinimumGasPrice",
      "anonymous" : false
   },
   {
      "anonymous" : false,
      "type" : "event",
      "name" : "Transfer",
      "inputs" : [
         {
            "name" : "from",
            "type" : "address",
            "internalType" : "address",
            "indexed" : true
         },
         {
            "indexed" : true,
            "internalType" : "address",
            "type" : "address",
            "name" : "to"
         },
         {
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256",
            "name" : "value"
         }
      ]
   },
   {
      "anonymous" : false,
      "type" : "event",
      "name" : "Version",
      "inputs" : [
         {
            "internalType" : "string",
            "indexed" : false,
            "type" : "string",
            "name" : "version"
         }
      ]
   },
   {
      "type" : "fallback",
      "stateMutability" : "payable"
   },
   {
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "internalType" : "address payable"
         },
         {
            "name" : "_enode",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "name" : "addParticipant",
      "type" : "function"
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "name" : "addStakeholder",
      "outputs" : [],
      "inputs" : [
         {
            "type" : "address",
            "name" : "_address",
            "internalType" : "address payable"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_enode"
         },
         {
            "type" : "uint256",
            "name" : "_stake",
            "internalType" : "uint256"
         }
      ]
   },
   {
      "inputs" : [
         {
            "internalType" : "address payable",
            "type" : "address",
            "name" : "_address"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_stake"
         },
         {
            "name" : "_enode",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function",
      "name" : "addValidator"
   },
   {
      "name" : "bondingPeriod",
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "inputs" : []
   },
   {
      "name" : "changeUserType",
      "type" : "function",
      "stateMutability" : "nonpayable",
      "outputs" : [],
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "name" : "newUserType",
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType"
         }
      ]
   },
   {
      "outputs" : [
         {
            "internalType" : "bool",
            "type" : "bool",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "name" : "checkMember",
      "type" : "function",
      "inputs" : [
         {
            "name" : "_account",
            "type" : "address",
            "internalType" : "address"
         }
      ]
   },
   {
      "inputs" : [],
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "name" : "committeeSize",
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "inputs" : [],
      "outputs" : [],
      "type" : "function",
      "name" : "computeCommittee",
      "stateMutability" : "nonpayable"
   },
   {
      "inputs" : [],
      "type" : "function",
      "name" : "contractVersion",
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "string",
            "name" : "",
            "internalType" : "string"
         }
      ]
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "address",
            "internalType" : "address"
         }
      ],
      "name" : "deployer",
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : []
   },
   {
      "inputs" : [],
      "stateMutability" : "view",
      "name" : "dumpEconomicsMetricData",
      "type" : "function",
      "outputs" : [
         {
            "internalType" : "struct Autonity.EconomicsMetricData",
            "name" : "economics",
            "type" : "tuple",
            "components" : [
               {
                  "type" : "address[]",
                  "name" : "accounts",
                  "internalType" : "address[]"
               },
               {
                  "internalType" : "enum Autonity.UserType[]",
                  "name" : "usertypes",
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
                  "name" : "mingasprice",
                  "internalType" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "type" : "uint256",
                  "name" : "stakesupply"
               }
            ]
         }
      ]
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
         }
      ],
      "outputs" : [
         {
            "internalType" : "string",
            "type" : "string",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "name" : "enodesWhitelist",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "name" : "_amount",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "type" : "function",
      "name" : "finalize",
      "stateMutability" : "nonpayable",
      "outputs" : [
         {
            "internalType" : "bool",
            "name" : "",
            "type" : "bool"
         },
         {
            "internalType" : "struct Autonity.CommitteeMember[]",
            "name" : "",
            "type" : "tuple[]",
            "components" : [
               {
                  "name" : "addr",
                  "type" : "address",
                  "internalType" : "address payable"
               },
               {
                  "name" : "votingPower",
                  "type" : "uint256",
                  "internalType" : "uint256"
               }
            ]
         }
      ]
   },
   {
      "inputs" : [
         {
            "name" : "_account",
            "type" : "address",
            "internalType" : "address"
         }
      ],
      "type" : "function",
      "name" : "getAccountStake",
      "stateMutability" : "view",
      "outputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
         }
      ]
   },
   {
      "outputs" : [
         {
            "type" : "tuple[]",
            "name" : "",
            "components" : [
               {
                  "name" : "addr",
                  "type" : "address",
                  "internalType" : "address payable"
               },
               {
                  "internalType" : "uint256",
                  "name" : "votingPower",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct Autonity.CommitteeMember[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function",
      "name" : "getCommittee",
      "inputs" : []
   },
   {
      "inputs" : [],
      "type" : "function",
      "name" : "getCurrentCommiteeSize",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ]
   },
   {
      "stateMutability" : "view",
      "type" : "function",
      "name" : "getMaxCommitteeSize",
      "outputs" : [
         {
            "type" : "uint256",
            "name" : "",
            "internalType" : "uint256"
         }
      ],
      "inputs" : []
   },
   {
      "stateMutability" : "view",
      "type" : "function",
      "name" : "getMinimumGasPrice",
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "inputs" : []
   },
   {
      "stateMutability" : "view",
      "name" : "getRate",
      "type" : "function",
      "outputs" : [
         {
            "type" : "uint256",
            "name" : "",
            "internalType" : "uint256"
         }
      ],
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         }
      ]
   },
   {
      "inputs" : [],
      "outputs" : [
         {
            "type" : "uint256",
            "name" : "",
            "internalType" : "uint256"
         }
      ],
      "type" : "function",
      "name" : "getStake",
      "stateMutability" : "view"
   },
   {
      "outputs" : [
         {
            "internalType" : "address[]",
            "type" : "address[]",
            "name" : ""
         }
      ],
      "type" : "function",
      "name" : "getStakeholders",
      "stateMutability" : "view",
      "inputs" : []
   },
   {
      "inputs" : [],
      "stateMutability" : "view",
      "name" : "getValidators",
      "type" : "function",
      "outputs" : [
         {
            "internalType" : "address[]",
            "type" : "address[]",
            "name" : ""
         }
      ]
   },
   {
      "name" : "getVersion",
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "inputs" : []
   },
   {
      "stateMutability" : "view",
      "type" : "function",
      "name" : "getWhitelist",
      "outputs" : [
         {
            "type" : "string[]",
            "name" : "",
            "internalType" : "string[]"
         }
      ],
      "inputs" : []
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         },
         {
            "type" : "uint256",
            "name" : "_amount",
            "internalType" : "uint256"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function",
      "name" : "mintStake",
      "outputs" : []
   },
   {
      "inputs" : [],
      "name" : "myUserType",
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8",
            "name" : ""
         }
      ]
   },
   {
      "inputs" : [],
      "name" : "operatorAccount",
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "address",
            "name" : "",
            "internalType" : "address"
         }
      ]
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_account",
            "type" : "address"
         },
         {
            "name" : "_amount",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "name" : "redeemStake",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_address",
            "type" : "address"
         }
      ],
      "stateMutability" : "nonpayable",
      "name" : "removeUser",
      "type" : "function",
      "outputs" : []
   },
   {
      "inputs" : [],
      "outputs" : [
         {
            "type" : "string",
            "name" : "",
            "internalType" : "string"
         },
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "stateMutability" : "view",
      "name" : "retrieveContract",
      "type" : "function"
   },
   {
      "outputs" : [
         {
            "internalType" : "address[]",
            "name" : "",
            "type" : "address[]"
         },
         {
            "internalType" : "string[]",
            "type" : "string[]",
            "name" : ""
         },
         {
            "internalType" : "uint256[]",
            "type" : "uint256[]",
            "name" : ""
         },
         {
            "type" : "uint256[]",
            "name" : "",
            "internalType" : "uint256[]"
         },
         {
            "internalType" : "uint256[]",
            "name" : "",
            "type" : "uint256[]"
         },
         {
            "name" : "",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
         },
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : ""
         }
      ],
      "type" : "function",
      "name" : "retrieveState",
      "stateMutability" : "view",
      "inputs" : []
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "name" : "send",
      "outputs" : [
         {
            "internalType" : "bool",
            "name" : "",
            "type" : "bool"
         }
      ],
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_recipient"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ]
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "bool",
            "internalType" : "bool"
         }
      ],
      "name" : "setCommissionRate",
      "type" : "function",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "type" : "uint256",
            "name" : "_rate",
            "internalType" : "uint256"
         }
      ]
   },
   {
      "inputs" : [
         {
            "name" : "_size",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "name" : "setCommitteeSize",
      "type" : "function"
   },
   {
      "name" : "setMinimumGasPrice",
      "type" : "function",
      "stateMutability" : "nonpayable",
      "outputs" : [],
      "inputs" : [
         {
            "name" : "_value",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ]
   },
   {
      "outputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
         }
      ],
      "name" : "totalSupply",
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : []
   },
   {
      "inputs" : [
         {
            "type" : "string",
            "name" : "_bytecode",
            "internalType" : "string"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_abi"
         },
         {
            "name" : "_version",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "type" : "function",
      "name" : "upgradeContract",
      "stateMutability" : "nonpayable",
      "outputs" : [
         {
            "name" : "",
            "type" : "bool",
            "internalType" : "bool"
         }
      ]
   },
   {
      "type" : "receive",
      "stateMutability" : "payable"
   }
]
`
