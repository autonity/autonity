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
            "name" : "_participantType",
            "type" : "uint256[]",
            "internalType" : "uint256[]"
         },
         {
            "name" : "_participantStake",
            "internalType" : "uint256[]",
            "type" : "uint256[]"
         },
         {
            "name" : "_commissionRate",
            "internalType" : "uint256[]",
            "type" : "uint256[]"
         },
         {
            "type" : "address",
            "internalType" : "address",
            "name" : "_operatorAccount"
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
            "name" : "_committeeSize",
            "type" : "uint256",
            "internalType" : "uint256"
         },
         {
            "name" : "_contractVersion",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "constructor"
   },
   {
      "name" : "AddParticipant",
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "name" : "_address",
            "indexed" : false,
            "type" : "address",
            "internalType" : "address"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_stake",
            "indexed" : false
         }
      ]
   },
   {
      "name" : "AddStakeholder",
      "inputs" : [
         {
            "indexed" : false,
            "name" : "_address",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_stake",
            "indexed" : false
         }
      ],
      "anonymous" : false,
      "type" : "event"
   },
   {
      "type" : "event",
      "inputs" : [
         {
            "name" : "_address",
            "indexed" : false,
            "internalType" : "address",
            "type" : "address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "indexed" : false,
            "name" : "_stake"
         }
      ],
      "anonymous" : false,
      "name" : "AddValidator"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_address",
            "indexed" : false
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_amount",
            "indexed" : false
         }
      ],
      "anonymous" : false,
      "type" : "event",
      "name" : "BlockReward"
   },
   {
      "type" : "event",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "indexed" : false,
            "name" : "_address"
         },
         {
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType",
            "indexed" : false,
            "name" : "_oldType"
         },
         {
            "name" : "_newType",
            "indexed" : false,
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType"
         }
      ],
      "anonymous" : false,
      "name" : "ChangeUserType"
   },
   {
      "type" : "event",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "indexed" : false,
            "name" : "_address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_amount",
            "indexed" : false
         }
      ],
      "anonymous" : false,
      "name" : "MintStake"
   },
   {
      "name" : "RedeemStake",
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_address",
            "indexed" : false
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_amount",
            "indexed" : false
         }
      ]
   },
   {
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "indexed" : false,
            "name" : "_address"
         },
         {
            "name" : "_type",
            "indexed" : false,
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType"
         }
      ],
      "name" : "RemoveUser"
   },
   {
      "name" : "SetCommissionRate",
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address",
            "indexed" : false,
            "name" : "_address"
         },
         {
            "name" : "_value",
            "indexed" : false,
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ]
   },
   {
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "name" : "_gasPrice",
            "indexed" : false,
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "name" : "SetMinimumGasPrice"
   },
   {
      "name" : "Transfer",
      "inputs" : [
         {
            "name" : "from",
            "indexed" : true,
            "type" : "address",
            "internalType" : "address"
         },
         {
            "indexed" : true,
            "name" : "to",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "value",
            "indexed" : false
         }
      ],
      "anonymous" : false,
      "type" : "event"
   },
   {
      "name" : "Version",
      "type" : "event",
      "anonymous" : false,
      "inputs" : [
         {
            "name" : "version",
            "indexed" : false,
            "internalType" : "string",
            "type" : "string"
         }
      ]
   },
   {
      "type" : "fallback",
      "stateMutability" : "payable"
   },
   {
      "name" : "addParticipant",
      "inputs" : [
         {
            "internalType" : "address payable",
            "type" : "address",
            "name" : "_address"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_enode"
         }
      ],
      "outputs" : [],
      "type" : "function",
      "stateMutability" : "nonpayable"
   },
   {
      "name" : "addStakeholder",
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
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_stake"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : []
   },
   {
      "name" : "addValidator",
      "outputs" : [],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "internalType" : "address payable"
         },
         {
            "name" : "_stake",
            "internalType" : "uint256",
            "type" : "uint256"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_enode"
         }
      ]
   },
   {
      "name" : "bondingPeriod",
      "inputs" : [],
      "type" : "function",
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
      "inputs" : [
         {
            "name" : "_address",
            "internalType" : "address",
            "type" : "address"
         },
         {
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType",
            "name" : "newUserType"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : [],
      "name" : "changeUserType"
   },
   {
      "name" : "checkMember",
      "inputs" : [
         {
            "name" : "_account",
            "internalType" : "address",
            "type" : "address"
         }
      ],
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "type" : "bool",
            "internalType" : "bool"
         }
      ]
   },
   {
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "name" : "committeeSize"
   },
   {
      "inputs" : [],
      "type" : "function",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "name" : "computeCommittee"
   },
   {
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "name" : "contractVersion"
   },
   {
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "type" : "address",
            "internalType" : "address"
         }
      ],
      "stateMutability" : "view",
      "name" : "deployer"
   },
   {
      "name" : "dumpEconomicsMetricData",
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "internalType" : "struct Autonity.EconomicsMetricData",
            "type" : "tuple",
            "components" : [
               {
                  "name" : "accounts",
                  "internalType" : "address[]",
                  "type" : "address[]"
               },
               {
                  "internalType" : "enum Autonity.UserType[]",
                  "type" : "uint8[]",
                  "name" : "usertypes"
               },
               {
                  "type" : "uint256[]",
                  "internalType" : "uint256[]",
                  "name" : "stakes"
               },
               {
                  "type" : "uint256[]",
                  "internalType" : "uint256[]",
                  "name" : "commissionrates"
               },
               {
                  "type" : "uint256",
                  "internalType" : "uint256",
                  "name" : "mingasprice"
               },
               {
                  "internalType" : "uint256",
                  "type" : "uint256",
                  "name" : "stakesupply"
               }
            ],
            "name" : "economics"
         }
      ]
   },
   {
      "name" : "enodesWhitelist",
      "inputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "string",
            "type" : "string"
         }
      ]
   },
   {
      "name" : "finalize",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "bool",
            "type" : "bool"
         },
         {
            "name" : "",
            "components" : [
               {
                  "name" : "addr",
                  "internalType" : "address payable",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "type" : "uint256",
                  "name" : "votingPower"
               }
            ],
            "internalType" : "struct Autonity.CommitteeMember[]",
            "type" : "tuple[]"
         }
      ],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_amount"
         }
      ]
   },
   {
      "name" : "getAccountStake",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         }
      ],
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ],
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "name" : "getCommittee",
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "components" : [
               {
                  "internalType" : "address payable",
                  "type" : "address",
                  "name" : "addr"
               },
               {
                  "internalType" : "uint256",
                  "type" : "uint256",
                  "name" : "votingPower"
               }
            ],
            "internalType" : "struct Autonity.CommitteeMember[]",
            "type" : "tuple[]"
         }
      ],
      "stateMutability" : "view"
   },
   {
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "inputs" : [],
      "name" : "getCurrentCommiteeSize"
   },
   {
      "name" : "getMaxCommitteeSize",
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
      "name" : "getMinimumGasPrice",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "type" : "function",
      "stateMutability" : "view",
      "inputs" : []
   },
   {
      "name" : "getRate",
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address",
            "name" : "_account"
         }
      ]
   },
   {
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "inputs" : [],
      "name" : "getStake"
   },
   {
      "name" : "getStakeholders",
      "type" : "function",
      "outputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "inputs" : []
   },
   {
      "name" : "getValidators",
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : ""
         }
      ]
   },
   {
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "stateMutability" : "view",
      "inputs" : [],
      "name" : "getVersion"
   },
   {
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "internalType" : "string[]",
            "type" : "string[]",
            "name" : ""
         }
      ],
      "name" : "getWhitelist"
   },
   {
      "name" : "mintStake",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_amount"
         }
      ],
      "type" : "function",
      "outputs" : [],
      "stateMutability" : "nonpayable"
   },
   {
      "name" : "myUserType",
      "inputs" : [],
      "outputs" : [
         {
            "name" : "",
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType"
         }
      ],
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "inputs" : [],
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "type" : "address",
            "internalType" : "address"
         }
      ],
      "name" : "operatorAccount"
   },
   {
      "name" : "redeemStake",
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : [],
      "inputs" : [
         {
            "name" : "_account",
            "internalType" : "address",
            "type" : "address"
         },
         {
            "name" : "_amount",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ]
   },
   {
      "name" : "removeUser",
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address",
            "name" : "_address"
         }
      ],
      "type" : "function",
      "outputs" : [],
      "stateMutability" : "nonpayable"
   },
   {
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         },
         {
            "name" : "",
            "internalType" : "string",
            "type" : "string"
         }
      ],
      "inputs" : [],
      "name" : "retrieveContract"
   },
   {
      "name" : "retrieveState",
      "stateMutability" : "view",
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "address[]",
            "type" : "address[]"
         },
         {
            "internalType" : "string[]",
            "type" : "string[]",
            "name" : ""
         },
         {
            "name" : "",
            "internalType" : "uint256[]",
            "type" : "uint256[]"
         },
         {
            "internalType" : "uint256[]",
            "type" : "uint256[]",
            "name" : ""
         },
         {
            "name" : "",
            "type" : "uint256[]",
            "internalType" : "uint256[]"
         },
         {
            "internalType" : "address",
            "type" : "address",
            "name" : ""
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         },
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         },
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         },
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "inputs" : []
   },
   {
      "name" : "send",
      "type" : "function",
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
            "type" : "address",
            "internalType" : "address",
            "name" : "_recipient"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_amount"
         }
      ]
   },
   {
      "name" : "setCommissionRate",
      "outputs" : [
         {
            "name" : "",
            "type" : "bool",
            "internalType" : "bool"
         }
      ],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_rate"
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
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : [],
      "name" : "setCommitteeSize"
   },
   {
      "name" : "setMinimumGasPrice",
      "outputs" : [],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_value"
         }
      ]
   },
   {
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "name" : "totalSupply"
   },
   {
      "name" : "upgradeContract",
      "stateMutability" : "nonpayable",
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "bool",
            "type" : "bool"
         }
      ],
      "inputs" : [
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_bytecode"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_abi"
         },
         {
            "name" : "_version",
            "internalType" : "string",
            "type" : "string"
         }
      ]
   },
   {
      "type" : "receive",
      "stateMutability" : "payable"
   }
]
`
