package generated

const Abi = `[
   {
      "inputs" : [
         {
            "name" : "_participantAddress",
            "internalType" : "address[]",
            "type" : "address[]"
         },
         {
            "name" : "_participantEnode",
            "internalType" : "string[]",
            "type" : "string[]"
         },
         {
            "type" : "uint256[]",
            "internalType" : "uint256[]",
            "name" : "_participantType"
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
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_minGasPrice"
         },
         {
            "name" : "_bondingPeriod",
            "internalType" : "uint256",
            "type" : "uint256"
         },
         {
            "name" : "_committeeSize",
            "type" : "uint256",
            "internalType" : "uint256"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_contractVersion"
         }
      ],
      "type" : "constructor",
      "stateMutability" : "nonpayable"
   },
   {
      "name" : "AddParticipant",
      "anonymous" : false,
      "type" : "event",
      "inputs" : [
         {
            "indexed" : false,
            "type" : "address",
            "internalType" : "address",
            "name" : "_address"
         },
         {
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256",
            "name" : "_stake"
         }
      ]
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
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256",
            "name" : "_stake"
         }
      ],
      "type" : "event",
      "anonymous" : false,
      "name" : "AddStakeholder"
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
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256",
            "name" : "_stake"
         }
      ],
      "anonymous" : false,
      "name" : "AddValidator"
   },
   {
      "name" : "BlockReward",
      "anonymous" : false,
      "type" : "event",
      "inputs" : [
         {
            "internalType" : "address",
            "indexed" : false,
            "type" : "address",
            "name" : "_address"
         },
         {
            "name" : "_amount",
            "internalType" : "uint256",
            "type" : "uint256",
            "indexed" : false
         }
      ]
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
            "internalType" : "enum Autonity.UserType",
            "indexed" : false,
            "type" : "uint8",
            "name" : "_oldType"
         },
         {
            "name" : "_newType",
            "internalType" : "enum Autonity.UserType",
            "indexed" : false,
            "type" : "uint8"
         }
      ],
      "type" : "event",
      "name" : "ChangeUserType"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "internalType" : "address",
            "indexed" : false,
            "type" : "address",
            "name" : "_address"
         },
         {
            "name" : "_amount",
            "type" : "uint256",
            "indexed" : false,
            "internalType" : "uint256"
         }
      ],
      "type" : "event",
      "name" : "MintStake"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "indexed" : false,
            "internalType" : "address"
         },
         {
            "type" : "uint256",
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_amount"
         }
      ],
      "type" : "event",
      "name" : "RedeemStake"
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
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8",
            "indexed" : false,
            "name" : "_type"
         }
      ],
      "anonymous" : false,
      "name" : "RemoveUser"
   },
   {
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
            "name" : "_value",
            "indexed" : false,
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "name" : "SetCommissionRate"
   },
   {
      "name" : "SetMinimumGasPrice",
      "inputs" : [
         {
            "name" : "_gasPrice",
            "type" : "uint256",
            "indexed" : false,
            "internalType" : "uint256"
         }
      ],
      "type" : "event",
      "anonymous" : false
   },
   {
      "name" : "Transfer",
      "inputs" : [
         {
            "name" : "from",
            "type" : "address",
            "indexed" : true,
            "internalType" : "address"
         },
         {
            "name" : "to",
            "internalType" : "address",
            "indexed" : true,
            "type" : "address"
         },
         {
            "name" : "value",
            "internalType" : "uint256",
            "indexed" : false,
            "type" : "uint256"
         }
      ],
      "type" : "event",
      "anonymous" : false
   },
   {
      "type" : "event",
      "inputs" : [
         {
            "name" : "version",
            "internalType" : "string",
            "type" : "string",
            "indexed" : false
         }
      ],
      "anonymous" : false,
      "name" : "Version"
   },
   {
      "stateMutability" : "payable",
      "type" : "fallback"
   },
   {
      "name" : "addParticipant",
      "stateMutability" : "nonpayable",
      "outputs" : [],
      "inputs" : [
         {
            "name" : "_address",
            "internalType" : "address payable",
            "type" : "address"
         },
         {
            "name" : "_enode",
            "internalType" : "string",
            "type" : "string"
         }
      ],
      "type" : "function"
   },
   {
      "stateMutability" : "nonpayable",
      "outputs" : [],
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address payable",
            "name" : "_address"
         },
         {
            "internalType" : "string",
            "type" : "string",
            "name" : "_enode"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_stake"
         }
      ],
      "type" : "function",
      "name" : "addStakeholder"
   },
   {
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address payable",
            "name" : "_address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_stake"
         },
         {
            "name" : "_enode",
            "internalType" : "string",
            "type" : "string"
         }
      ],
      "type" : "function",
      "name" : "addValidator"
   },
   {
      "type" : "function",
      "inputs" : [],
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ],
      "name" : "bondingPeriod"
   },
   {
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "name" : "_address",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "internalType" : "enum Autonity.UserType",
            "type" : "uint8",
            "name" : "newUserType"
         }
      ],
      "type" : "function",
      "name" : "changeUserType"
   },
   {
      "name" : "checkMember",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
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
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ],
      "inputs" : [],
      "type" : "function",
      "name" : "committeeSize"
   },
   {
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "inputs" : [],
      "type" : "function",
      "name" : "computeCommittee"
   },
   {
      "name" : "contractVersion",
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "string",
            "internalType" : "string",
            "name" : ""
         }
      ]
   },
   {
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "address",
            "type" : "address"
         }
      ],
      "stateMutability" : "view",
      "name" : "deployer"
   },
   {
      "name" : "dumpEconomicsMetricData",
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "tuple",
            "internalType" : "struct Autonity.EconomicsMetricData",
            "name" : "economics",
            "components" : [
               {
                  "internalType" : "address[]",
                  "type" : "address[]",
                  "name" : "accounts"
               },
               {
                  "name" : "usertypes",
                  "internalType" : "enum Autonity.UserType[]",
                  "type" : "uint8[]"
               },
               {
                  "internalType" : "uint256[]",
                  "type" : "uint256[]",
                  "name" : "stakes"
               },
               {
                  "type" : "uint256[]",
                  "internalType" : "uint256[]",
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
      ],
      "inputs" : [],
      "type" : "function"
   },
   {
      "outputs" : [
         {
            "type" : "string",
            "internalType" : "string",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "type" : "function",
      "inputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "name" : "enodesWhitelist"
   },
   {
      "outputs" : [
         {
            "internalType" : "bool",
            "type" : "bool",
            "name" : ""
         },
         {
            "components" : [
               {
                  "type" : "address",
                  "internalType" : "address payable",
                  "name" : "addr"
               },
               {
                  "internalType" : "uint256",
                  "type" : "uint256",
                  "name" : "votingPower"
               }
            ],
            "name" : "",
            "internalType" : "struct Autonity.CommitteeMember[]",
            "type" : "tuple[]"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function",
      "inputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_amount"
         }
      ],
      "name" : "finalize"
   },
   {
      "outputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "inputs" : [
         {
            "internalType" : "address",
            "type" : "address",
            "name" : "_account"
         }
      ],
      "type" : "function",
      "name" : "getAccountStake"
   },
   {
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
                  "internalType" : "uint256",
                  "type" : "uint256"
               }
            ],
            "name" : "",
            "internalType" : "struct Autonity.CommitteeMember[]",
            "type" : "tuple[]"
         }
      ],
      "stateMutability" : "view",
      "inputs" : [],
      "type" : "function",
      "name" : "getCommittee"
   },
   {
      "inputs" : [],
      "type" : "function",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "name" : "getCurrentCommiteeSize"
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
      "type" : "function",
      "name" : "getMaxCommitteeSize"
   },
   {
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "inputs" : [],
      "type" : "function",
      "name" : "getMinimumGasPrice"
   },
   {
      "name" : "getRate",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "type" : "uint256",
            "internalType" : "uint256"
         }
      ],
      "inputs" : [
         {
            "name" : "_account",
            "type" : "address",
            "internalType" : "address"
         }
      ],
      "type" : "function"
   },
   {
      "name" : "getStake",
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "uint256",
            "type" : "uint256"
         }
      ],
      "type" : "function",
      "inputs" : []
   },
   {
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "name" : "getStakeholders"
   },
   {
      "name" : "getValidators",
      "outputs" : [
         {
            "internalType" : "address[]",
            "type" : "address[]",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "inputs" : [],
      "type" : "function"
   },
   {
      "name" : "getVersion",
      "outputs" : [
         {
            "internalType" : "string",
            "type" : "string",
            "name" : ""
         }
      ],
      "stateMutability" : "view",
      "type" : "function",
      "inputs" : []
   },
   {
      "name" : "getWhitelist",
      "type" : "function",
      "inputs" : [],
      "outputs" : [
         {
            "internalType" : "string[]",
            "type" : "string[]",
            "name" : ""
         }
      ],
      "stateMutability" : "view"
   },
   {
      "name" : "mintStake",
      "inputs" : [
         {
            "name" : "_account",
            "internalType" : "address",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_amount"
         }
      ],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "outputs" : []
   },
   {
      "name" : "myUserType",
      "type" : "function",
      "inputs" : [],
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "uint8",
            "internalType" : "enum Autonity.UserType",
            "name" : ""
         }
      ]
   },
   {
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "type" : "address",
            "internalType" : "address"
         }
      ],
      "type" : "function",
      "inputs" : [],
      "name" : "operatorAccount"
   },
   {
      "name" : "redeemStake",
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
      "stateMutability" : "nonpayable"
   },
   {
      "stateMutability" : "nonpayable",
      "outputs" : [],
      "inputs" : [
         {
            "type" : "address",
            "internalType" : "address",
            "name" : "_address"
         }
      ],
      "type" : "function",
      "name" : "removeUser"
   },
   {
      "stateMutability" : "view",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "string",
            "type" : "string"
         },
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "inputs" : [],
      "type" : "function",
      "name" : "retrieveContract"
   },
   {
      "name" : "retrieveState",
      "inputs" : [],
      "type" : "function",
      "outputs" : [
         {
            "type" : "address[]",
            "internalType" : "address[]",
            "name" : ""
         },
         {
            "type" : "string[]",
            "internalType" : "string[]",
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
            "internalType" : "uint256[]",
            "type" : "uint256[]",
            "name" : ""
         },
         {
            "name" : "",
            "type" : "address",
            "internalType" : "address"
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         },
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         },
         {
            "name" : "",
            "type" : "string",
            "internalType" : "string"
         }
      ],
      "stateMutability" : "view"
   },
   {
      "name" : "send",
      "type" : "function",
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
      ],
      "stateMutability" : "nonpayable",
      "outputs" : [
         {
            "name" : "",
            "internalType" : "bool",
            "type" : "bool"
         }
      ]
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_rate"
         }
      ],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "outputs" : [
         {
            "internalType" : "bool",
            "type" : "bool",
            "name" : ""
         }
      ],
      "name" : "setCommissionRate"
   },
   {
      "name" : "setCommitteeSize",
      "inputs" : [
         {
            "internalType" : "uint256",
            "type" : "uint256",
            "name" : "_size"
         }
      ],
      "type" : "function",
      "stateMutability" : "nonpayable",
      "outputs" : []
   },
   {
      "name" : "setMinimumGasPrice",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "inputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : "_value"
         }
      ],
      "type" : "function"
   },
   {
      "name" : "totalSupply",
      "stateMutability" : "view",
      "outputs" : [
         {
            "type" : "uint256",
            "internalType" : "uint256",
            "name" : ""
         }
      ],
      "type" : "function",
      "inputs" : []
   },
   {
      "name" : "upgradeContract",
      "inputs" : [
         {
            "name" : "_bytecode",
            "type" : "string",
            "internalType" : "string"
         },
         {
            "type" : "string",
            "internalType" : "string",
            "name" : "_abi"
         },
         {
            "name" : "_version",
            "type" : "string",
            "internalType" : "string"
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
      "stateMutability" : "nonpayable"
   },
   {
      "stateMutability" : "payable",
      "type" : "receive"
   }
]
`
