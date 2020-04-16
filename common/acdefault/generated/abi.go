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
            "type" : "uint256[]",
            "name" : "_participantStake",
            "internalType" : "uint256[]"
         },
         {
            "name" : "_commissionRate",
            "type" : "uint256[]",
            "internalType" : "uint256[]"
         },
         {
            "name" : "_operatorAccount",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "name" : "_minGasPrice",
            "type" : "uint256",
            "internalType" : "uint256"
         },
         {
            "type" : "uint256",
            "name" : "_bondingPeriod",
            "internalType" : "uint256"
         },
         {
            "type" : "uint256",
            "name" : "_committeeSize",
            "internalType" : "uint256"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_contractVersion"
         }
      ]
   },
   {
      "type" : "event",
      "name" : "AddParticipant",
      "inputs" : [
         {
            "internalType" : "address",
            "indexed" : false,
            "type" : "address",
            "name" : "_address"
         },
         {
            "type" : "uint256",
            "name" : "_stake",
            "indexed" : false,
            "internalType" : "uint256"
         }
      ],
      "anonymous" : false
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "indexed" : false,
            "name" : "_address"
         },
         {
            "type" : "uint256",
            "indexed" : false,
            "name" : "_stake",
            "internalType" : "uint256"
         }
      ],
      "name" : "AddStakeholder",
      "type" : "event"
   },
   {
      "type" : "event",
      "name" : "AddValidator",
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "indexed" : false,
            "internalType" : "address"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "indexed" : false,
            "name" : "_stake"
         }
      ],
      "anonymous" : false
   },
   {
      "name" : "BlockReward",
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "internalType" : "address",
            "indexed" : false,
            "type" : "address",
            "name" : "_address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256",
            "indexed" : false
         }
      ]
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_address",
            "indexed" : false
         },
         {
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8",
            "indexed" : false,
            "name" : "_oldType"
         },
         {
            "internalType" : "enum Autonity.UserType",
            "name" : "_newType",
            "type" : "uint8",
            "indexed" : false
         }
      ],
      "type" : "event",
      "name" : "ChangeUserType"
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
            "type" : "uint256",
            "name" : "_amount",
            "internalType" : "uint256"
         }
      ],
      "anonymous" : false,
      "name" : "MintStake",
      "type" : "event"
   },
   {
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
            "type" : "uint256",
            "indexed" : false,
            "name" : "_amount"
         }
      ],
      "name" : "RedeemStake",
      "type" : "event"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "indexed" : false,
            "name" : "_address"
         },
         {
            "indexed" : false,
            "type" : "uint8",
            "name" : "_type",
            "internalType" : "enum Autonity.UserType"
         }
      ],
      "anonymous" : false,
      "name" : "RemoveUser",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "indexed" : false,
            "name" : "_address"
         },
         {
            "type" : "uint256",
            "indexed" : false,
            "name" : "_value",
            "internalType" : "uint256"
         }
      ],
      "type" : "event",
      "name" : "SetCommissionRate"
   },
   {
      "type" : "event",
      "name" : "SetMinimumGasPrice",
      "inputs" : [
         {
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256",
            "name" : "_gasPrice"
         }
      ],
      "anonymous" : false
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "from",
            "type" : "address",
            "indexed" : true
         },
         {
            "name" : "to",
            "type" : "address",
            "indexed" : true,
            "internalType" : "address"
         },
         {
            "indexed" : false,
            "type" : "uint256",
            "name" : "value",
            "internalType" : "uint256"
         }
      ],
      "anonymous" : false,
      "type" : "event",
      "name" : "Transfer"
   },
   {
      "inputs" : [
         {
            "indexed" : false,
            "type" : "string",
            "name" : "version",
            "internalType" : "string"
         }
      ],
      "anonymous" : false,
      "name" : "Version",
      "type" : "event"
   },
   {
      "stateMutability" : "payable",
      "type" : "fallback"
   },
   {
      "outputs" : [],
      "type" : "function",
      "name" : "addParticipant",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "address payable",
            "type" : "address",
            "name" : "_address"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_enode"
         }
      ]
   },
   {
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_address",
            "type" : "address"
         },
         {
            "internalType" : "string",
            "name" : "_enode",
            "type" : "string"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_stake"
         }
      ],
      "type" : "function",
      "name" : "addStakeholder",
      "outputs" : []
   },
   {
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_address",
            "type" : "address"
         },
         {
            "type" : "uint256",
            "name" : "_stake",
            "internalType" : "uint256"
         },
         {
            "internalType" : "string",
            "name" : "_enode",
            "type" : "string"
         }
      ],
      "name" : "addValidator",
      "type" : "function",
      "outputs" : []
   },
   {
      "outputs" : [
         {
            "type" : "uint256",
            "name" : "",
            "internalType" : "uint256"
         }
      ],
      "inputs" : [],
      "stateMutability" : "view",
      "name" : "bondingPeriod",
      "type" : "function"
   },
   {
      "name" : "changeUserType",
      "type" : "function",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_address"
         },
         {
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8",
            "name" : "newUserType"
         }
      ],
      "stateMutability" : "nonpayable",
      "outputs" : []
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         }
      ],
      "stateMutability" : "view",
      "type" : "function",
      "name" : "checkMember",
      "outputs" : [
         {
            "name" : "",
            "type" : "bool",
            "internalType" : "bool"
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
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "name" : "committeeSize"
   },
   {
      "stateMutability" : "nonpayable",
      "inputs" : [],
      "name" : "computeCommittee",
      "type" : "function",
      "outputs" : []
   },
   {
      "name" : "contractVersion",
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : [],
      "outputs" : [
         {
            "type" : "string",
            "name" : "",
            "internalType" : "string"
         }
      ]
   },
   {
      "name" : "deployer",
      "type" : "function",
      "inputs" : [],
      "stateMutability" : "view",
      "outputs" : [
         {
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ]
   },
   {
      "outputs" : [
         {
            "internalType" : "struct Autonity.EconomicsMetricData",
            "type" : "tuple",
            "name" : "economics",
            "components" : [
               {
                  "internalType" : "address[]",
                  "name" : "accounts",
                  "type" : "address[]"
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
                  "internalType" : "uint256[]",
                  "type" : "uint256[]",
                  "name" : "commissionrates"
               },
               {
                  "internalType" : "uint256",
                  "name" : "mingasprice",
                  "type" : "uint256"
               },
               {
                  "name" : "stakesupply",
                  "type" : "uint256",
                  "internalType" : "uint256"
               }
            ]
         }
      ],
      "inputs" : [],
      "stateMutability" : "view",
      "name" : "dumpEconomicsMetricData",
      "type" : "function"
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "name" : "enodesWhitelist",
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ]
   },
   {
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "type" : "uint256",
            "name" : "_amount",
            "internalType" : "uint256"
         }
      ],
      "name" : "finalize",
      "type" : "function",
      "outputs" : [
         {
            "internalType" : "bool",
            "type" : "bool",
            "name" : ""
         },
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
      ]
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "inputs" : [
         {
            "type" : "address",
            "name" : "_account",
            "internalType" : "address"
         }
      ],
      "stateMutability" : "view",
      "name" : "getAccountStake",
      "type" : "function"
   },
   {
      "outputs" : [
         {
            "components" : [
               {
                  "type" : "address",
                  "name" : "addr",
                  "internalType" : "address payable"
               },
               {
                  "internalType" : "uint256",
                  "name" : "votingPower",
                  "type" : "uint256"
               }
            ],
            "name" : "",
            "type" : "tuple[]",
            "internalType" : "struct Autonity.CommitteeMember[]"
         }
      ],
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "name" : "getCommittee"
   },
   {
      "type" : "function",
      "name" : "getCurrentCommiteeSize",
      "stateMutability" : "view",
      "inputs" : [],
      "outputs" : [
         {
            "type" : "uint256",
            "name" : "",
            "internalType" : "uint256"
         }
      ]
   },
   {
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "name" : "getMaxCommitteeSize",
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
      "inputs" : [],
      "name" : "getMinimumGasPrice",
      "type" : "function",
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
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "name" : "getRate",
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         }
      ]
   },
   {
      "name" : "getStake",
      "type" : "function",
      "inputs" : [],
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
      "stateMutability" : "view",
      "inputs" : [],
      "name" : "getStakeholders",
      "type" : "function",
      "outputs" : [
         {
            "type" : "address[]",
            "name" : "",
            "internalType" : "address[]"
         }
      ]
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "address[]",
            "internalType" : "address[]"
         }
      ],
      "stateMutability" : "view",
      "inputs" : [],
      "name" : "getValidators",
      "type" : "function"
   },
   {
      "type" : "function",
      "name" : "getVersion",
      "stateMutability" : "view",
      "inputs" : [],
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ]
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "string[]",
            "internalType" : "string[]"
         }
      ],
      "stateMutability" : "view",
      "inputs" : [],
      "type" : "function",
      "name" : "getWhitelist"
   },
   {
      "name" : "mintStake",
      "type" : "function",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "name" : "_account",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "type" : "uint256",
            "name" : "_amount",
            "internalType" : "uint256"
         }
      ],
      "outputs" : []
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType"
         }
      ],
      "inputs" : [],
      "stateMutability" : "view",
      "name" : "myUserType",
      "type" : "function"
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "address",
            "internalType" : "address"
         }
      ],
      "name" : "operatorAccount",
      "type" : "function",
      "inputs" : [],
      "stateMutability" : "view"
   },
   {
      "type" : "function",
      "name" : "redeemStake",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "nonpayable",
      "outputs" : []
   },
   {
      "name" : "removeUser",
      "type" : "function",
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_address",
            "type" : "address"
         }
      ],
      "stateMutability" : "nonpayable",
      "outputs" : []
   },
   {
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "name" : "retrieveContract",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         },
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
            "internalType" : "address[]",
            "type" : "address[]",
            "name" : ""
         },
         {
            "name" : "",
            "type" : "string[]",
            "internalType" : "string[]"
         },
         {
            "internalType" : "uint256[]",
            "name" : "",
            "type" : "uint256[]"
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
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         },
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "name" : "retrieveState",
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : []
   },
   {
      "outputs" : [
         {
            "internalType" : "bool",
            "type" : "bool",
            "name" : ""
         }
      ],
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "name" : "_recipient",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "name" : "_amount",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "type" : "function",
      "name" : "send"
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "bool",
            "internalType" : "bool"
         }
      ],
      "type" : "function",
      "name" : "setCommissionRate",
      "inputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_rate"
         }
      ],
      "stateMutability" : "nonpayable"
   },
   {
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_size"
         }
      ],
      "name" : "setCommitteeSize",
      "type" : "function",
      "outputs" : []
   },
   {
      "outputs" : [],
      "inputs" : [
         {
            "type" : "uint256",
            "name" : "_value",
            "internalType" : "uint256"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function",
      "name" : "setMinimumGasPrice"
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "stateMutability" : "view",
      "inputs" : [],
      "name" : "totalSupply",
      "type" : "function"
   },
   {
      "type" : "function",
      "name" : "upgradeContract",
      "inputs" : [
         {
            "internalType" : "string",
            "name" : "_bytecode",
            "type" : "string"
         },
         {
            "type" : "string",
            "name" : "_abi",
            "internalType" : "string"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_version"
         }
      ],
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
      "stateMutability" : "payable",
      "type" : "receive"
   }
]
`
