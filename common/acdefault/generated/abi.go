package generated

const Abi = `[
   {
      "stateMutability" : "nonpayable",
      "type" : "constructor",
      "inputs" : [
         {
            "internalType" : "address[]",
            "type" : "address[]",
            "name" : "_participantAddress"
         },
         {
            "type" : "string[]",
            "internalType" : "string[]",
            "name" : "_participantEnode"
         },
         {
            "type" : "uint256[]",
            "internalType" : "uint256[]",
            "name" : "_participantType"
         },
         {
            "internalType" : "uint256[]",
            "type" : "uint256[]",
            "name" : "_participantStake"
         },
         {
            "type" : "uint256[]",
            "internalType" : "uint256[]",
            "name" : "_commissionRate"
         },
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_operatorAccount"
         },
         {
            "name" : "_minGasPrice",
            "internalType" : "uint256",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_bondingPeriod"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_committeeSize"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_contractVersion"
         }
      ]
   },
   {
      "name" : "AddParticipant",
      "type" : "event",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_address",
            "indexed" : false
         },
         {
            "indexed" : false,
            "name" : "_stake",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "anonymous" : false
   },
   {
      "name" : "AddStakeholder",
      "type" : "event",
      "inputs" : [
         {
            "indexed" : false,
            "name" : "_address",
            "internalType" : "address",
            "type" : "address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_stake",
            "indexed" : false
         }
      ],
      "anonymous" : false
   },
   {
      "name" : "AddValidator",
      "type" : "event",
      "inputs" : [
         {
            "name" : "_address",
            "indexed" : false,
            "internalType" : "address",
            "type" : "address"
         },
         {
            "name" : "_stake",
            "indexed" : false,
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "anonymous" : false
   },
   {
      "inputs" : [
         {
            "name" : "_address",
            "indexed" : false,
            "type" : "address",
            "internalType" : "address"
         },
         {
            "indexed" : false,
            "name" : "_amount",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "type" : "event",
      "anonymous" : false,
      "name" : "BlockReward"
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
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType",
            "indexed" : false,
            "name" : "_oldType"
         },
         {
            "indexed" : false,
            "name" : "_newType",
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8"
         }
      ],
      "anonymous" : false,
      "name" : "ChangeUserType"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "name" : "_address",
            "indexed" : false,
            "type" : "address",
            "internalType" : "address"
         },
         {
            "name" : "_amount",
            "indexed" : false,
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "type" : "event",
      "name" : "MintStake"
   },
   {
      "name" : "RedeemStake",
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
      ],
      "type" : "event",
      "anonymous" : false
   },
   {
      "anonymous" : false,
      "type" : "event",
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address",
            "name" : "_address",
            "indexed" : false
         },
         {
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType",
            "indexed" : false,
            "name" : "_type"
         }
      ],
      "name" : "RemoveUser"
   },
   {
      "name" : "SetCommissionRate",
      "anonymous" : false,
      "type" : "event",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "indexed" : false,
            "name" : "_address"
         },
         {
            "indexed" : false,
            "name" : "_value",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ]
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "indexed" : false,
            "name" : "_gasPrice"
         }
      ],
      "type" : "event",
      "name" : "SetMinimumGasPrice"
   },
   {
      "inputs" : [
         {
            "name" : "from",
            "indexed" : true,
            "internalType" : "address",
            "type" : "address"
         },
         {
            "type" : "address",
            "internalType" : "address",
            "name" : "to",
            "indexed" : true
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "value",
            "indexed" : false
         }
      ],
      "type" : "event",
      "anonymous" : false,
      "name" : "Transfer"
   },
   {
      "inputs" : [
         {
            "type" : "string",
            "internalType" : "string",
            "indexed" : false,
            "name" : "version"
         }
      ],
      "type" : "event",
      "anonymous" : false,
      "name" : "Version"
   },
   {
      "stateMutability" : "payable",
      "type" : "fallback"
   },
   {
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
      "type" : "function",
      "stateMutability" : "nonpayable",
      "name" : "addParticipant",
      "outputs" : []
   },
   {
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address payable",
            "name" : "_address"
         },
         {
            "name" : "_enode",
            "internalType" : "string",
            "type" : "string"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_stake"
         }
      ],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "name" : "addStakeholder",
      "outputs" : []
   },
   {
      "type" : "function",
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address payable",
            "name" : "_address"
         },
         {
            "name" : "_stake",
            "internalType" : "uint256",
            "type" : "uint256"
         },
         {
            "name" : "_enode",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "stateMutability" : "nonpayable",
      "name" : "addValidator",
      "outputs" : []
   },
   {
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "name" : "bondingPeriod",
      "stateMutability" : "view",
      "inputs" : [],
      "type" : "function"
   },
   {
      "outputs" : [],
      "name" : "changeUserType",
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_address"
         },
         {
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType",
            "name" : "newUserType"
         }
      ],
      "type" : "function"
   },
   {
      "stateMutability" : "view",
      "inputs" : [
         {
            "name" : "_account",
            "internalType" : "address",
            "type" : "address"
         }
      ],
      "type" : "function",
      "outputs" : [
         {
            "internalType" : "bool",
            "type" : "bool",
            "name" : ""
         }
      ],
      "name" : "checkMember"
   },
   {
      "name" : "committeeSize",
      "outputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
         }
      ],
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "name" : "computeCommittee",
      "outputs" : [],
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "nonpayable"
   },
   {
      "stateMutability" : "view",
      "inputs" : [],
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
      "type" : "function",
      "inputs" : [],
      "stateMutability" : "view",
      "name" : "deployer",
      "outputs" : [
         {
            "name" : "",
            "type" : "address",
            "internalType" : "address"
         }
      ]
   },
   {
      "type" : "function",
      "inputs" : [],
      "stateMutability" : "view",
      "name" : "dumpEconomicsMetricData",
      "outputs" : [
         {
            "type" : "tuple",
            "internalType" : "struct Autonity.EconomicsMetricData",
            "name" : "economics",
            "components" : [
               {
                  "type" : "address[]",
                  "internalType" : "address[]",
                  "name" : "accounts"
               },
               {
                  "internalType" : "enum Autonity.UserType[]",
                  "type" : "uint8[]",
                  "name" : "usertypes"
               },
               {
                  "internalType" : "uint256[]",
                  "type" : "uint256[]",
                  "name" : "stakes"
               },
               {
                  "internalType" : "uint256[]",
                  "type" : "uint256[]",
                  "name" : "commissionrates"
               },
               {
                  "internalType" : "uint256",
                  "type" : "uint256",
                  "name" : "mingasprice"
               },
               {
                  "name" : "stakesupply",
                  "internalType" : "uint256",
                  "type" : "uint256"
               }
            ]
         }
      ]
   },
   {
      "outputs" : [
         {
            "name" : "",
            "internalType" : "string",
            "type" : "string"
         }
      ],
      "name" : "enodesWhitelist",
      "stateMutability" : "view",
      "type" : "function",
      "inputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ]
   },
   {
      "outputs" : [
         {
            "internalType" : "bool",
            "type" : "bool",
            "name" : ""
         },
         {
            "type" : "tuple[]",
            "internalType" : "struct Autonity.CommitteeMember[]",
            "components" : [
               {
                  "name" : "addr",
                  "type" : "address",
                  "internalType" : "address payable"
               },
               {
                  "type" : "uint256",
                  "internalType" : "uint256",
                  "name" : "votingPower"
               }
            ],
            "name" : ""
         }
      ],
      "name" : "finalize",
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "name" : "_amount",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ]
   },
   {
      "name" : "getAccountStake",
      "outputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
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
      "stateMutability" : "view"
   },
   {
      "name" : "getCommittee",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "address payable",
                  "type" : "address",
                  "name" : "addr"
               },
               {
                  "name" : "votingPower",
                  "type" : "uint256",
                  "internalType" : "uint256"
               }
            ],
            "name" : "",
            "internalType" : "struct Autonity.CommitteeMember[]",
            "type" : "tuple[]"
         }
      ],
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ],
      "name" : "getCurrentCommiteeSize",
      "stateMutability" : "view",
      "type" : "function",
      "inputs" : []
   },
   {
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ],
      "name" : "getMaxCommitteeSize",
      "stateMutability" : "view",
      "inputs" : [],
      "type" : "function"
   },
   {
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ],
      "name" : "getMinimumGasPrice",
      "stateMutability" : "view",
      "inputs" : [],
      "type" : "function"
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
      "stateMutability" : "view",
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address",
            "name" : "_account"
         }
      ],
      "type" : "function"
   },
   {
      "type" : "function",
      "inputs" : [],
      "stateMutability" : "view",
      "name" : "getStake",
      "outputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
         }
      ]
   },
   {
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view",
      "name" : "getStakeholders",
      "outputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : ""
         }
      ]
   },
   {
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view",
      "name" : "getValidators",
      "outputs" : [
         {
            "internalType" : "address[]",
            "type" : "address[]",
            "name" : ""
         }
      ]
   },
   {
      "stateMutability" : "view",
      "type" : "function",
      "inputs" : [],
      "outputs" : [
         {
            "type" : "string",
            "internalType" : "string",
            "name" : ""
         }
      ],
      "name" : "getVersion"
   },
   {
      "stateMutability" : "view",
      "type" : "function",
      "inputs" : [],
      "outputs" : [
         {
            "type" : "string[]",
            "internalType" : "string[]",
            "name" : ""
         }
      ],
      "name" : "getWhitelist"
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "name" : "_account",
            "internalType" : "address",
            "type" : "address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_amount"
         }
      ],
      "outputs" : [],
      "name" : "mintStake"
   },
   {
      "stateMutability" : "view",
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8"
         }
      ],
      "name" : "myUserType"
   },
   {
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view",
      "name" : "operatorAccount",
      "outputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : ""
         }
      ]
   },
   {
      "type" : "function",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         },
         {
            "name" : "_amount",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "stateMutability" : "nonpayable",
      "name" : "redeemStake",
      "outputs" : []
   },
   {
      "outputs" : [],
      "name" : "removeUser",
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "name" : "_address",
            "internalType" : "address",
            "type" : "address"
         }
      ]
   },
   {
      "name" : "retrieveContract",
      "outputs" : [
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : ""
         }
      ],
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view"
   },
   {
      "name" : "retrieveState",
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
            "name" : "",
            "type" : "uint256[]",
            "internalType" : "uint256[]"
         },
         {
            "type" : "uint256[]",
            "internalType" : "uint256[]",
            "name" : ""
         },
         {
            "internalType" : "uint256[]",
            "type" : "uint256[]",
            "name" : ""
         },
         {
            "name" : "",
            "internalType" : "address",
            "type" : "address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : ""
         }
      ],
      "type" : "function",
      "inputs" : [],
      "stateMutability" : "view"
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_recipient"
         },
         {
            "name" : "_amount",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "outputs" : [
         {
            "type" : "bool",
            "internalType" : "bool",
            "name" : ""
         }
      ],
      "name" : "send"
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_rate"
         }
      ],
      "outputs" : [
         {
            "name" : "",
            "type" : "bool",
            "internalType" : "bool"
         }
      ],
      "name" : "setCommissionRate"
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_size"
         }
      ],
      "outputs" : [],
      "name" : "setCommitteeSize"
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "name" : "_value",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "outputs" : [],
      "name" : "setMinimumGasPrice"
   },
   {
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ],
      "name" : "totalSupply",
      "stateMutability" : "view",
      "type" : "function",
      "inputs" : []
   },
   {
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "name" : "_bytecode",
            "internalType" : "string",
            "type" : "string"
         },
         {
            "name" : "_abi",
            "internalType" : "string",
            "type" : "string"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_version"
         }
      ],
      "outputs" : [
         {
            "type" : "bool",
            "internalType" : "bool",
            "name" : ""
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
