package generated

const UpgradeTestAbi = `[
   {
      "inputs" : [],
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
            "indexed" : false,
            "internalType" : "address",
            "name" : "addr",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         }
      ],
      "name" : "BurnedStake",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "validator",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "rate",
            "type" : "uint256"
         }
      ],
      "name" : "CommissionRateChange",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "gasPrice",
            "type" : "uint256"
         }
      ],
      "name" : "MinimumBaseFeeUpdated",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "addr",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         }
      ],
      "name" : "MintedStake",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "treasury",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "addr",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "effectiveBlock",
            "type" : "uint256"
         }
      ],
      "name" : "PausedValidator",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "treasury",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "addr",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "string",
            "name" : "enode",
            "type" : "string"
         },
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "liquidContract",
            "type" : "address"
         }
      ],
      "name" : "RegisteredValidator",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "addr",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         }
      ],
      "name" : "Rewarded",
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
      "stateMutability" : "payable",
      "type" : "fallback"
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
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_address",
            "type" : "address"
         }
      ],
      "name" : "activateValidator",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "owner",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "spender",
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
            "name" : "spender",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "amount",
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
            "name" : "_addr",
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
            "name" : "_validator",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "bond",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_addr",
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
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_validator",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_rate",
            "type" : "uint256"
         }
      ],
      "name" : "changeCommissionRate",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "completeContractUpgrade",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "computeCommittee",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "config",
      "outputs" : [
         {
            "internalType" : "address",
            "name" : "operatorAccount",
            "type" : "address"
         },
         {
            "internalType" : "address payable",
            "name" : "treasuryAccount",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "treasuryFee",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "minBaseFee",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "delegationRate",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "epochPeriod",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "unbondingPeriod",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "committeeSize",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "contractVersion",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "blockPeriod",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "deployer",
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
      "name" : "epochID",
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
      "name" : "epochReward",
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
      "name" : "epochTotalBondedStake",
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
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         }
      ],
      "name" : "finalize",
      "outputs" : [
         {
            "internalType" : "bool",
            "name" : "",
            "type" : "bool"
         },
         {
            "components" : [
               {
                  "internalType" : "address",
                  "name" : "addr",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "name" : "votingPower",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct Autonity.CommitteeMember[]",
            "name" : "",
            "type" : "tuple[]"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "startId",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "lastId",
            "type" : "uint256"
         }
      ],
      "name" : "getBondingReq",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "address payable",
                  "name" : "delegator",
                  "type" : "address"
               },
               {
                  "internalType" : "address",
                  "name" : "delegatee",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "name" : "amount",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "startBlock",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct Autonity.Staking[]",
            "name" : "",
            "type" : "tuple[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getCommittee",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "address",
                  "name" : "addr",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "name" : "votingPower",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct Autonity.CommitteeMember[]",
            "name" : "",
            "type" : "tuple[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getCommitteeEnodes",
      "outputs" : [
         {
            "internalType" : "string[]",
            "name" : "",
            "type" : "string[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getLastEpochBlock",
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
      "name" : "getMaxCommitteeSize",
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
      "name" : "getMinimumBaseFee",
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
      "name" : "getNewContract",
      "outputs" : [
         {
            "internalType" : "bytes",
            "name" : "",
            "type" : "bytes"
         },
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getOperator",
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
            "internalType" : "uint256",
            "name" : "height",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "round",
            "type" : "uint256"
         }
      ],
      "name" : "getProposer",
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
            "internalType" : "uint256",
            "name" : "startId",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "lastId",
            "type" : "uint256"
         }
      ],
      "name" : "getUnbondingReq",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "address payable",
                  "name" : "delegator",
                  "type" : "address"
               },
               {
                  "internalType" : "address",
                  "name" : "delegatee",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "name" : "amount",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "startBlock",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct Autonity.Staking[]",
            "name" : "",
            "type" : "tuple[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_addr",
            "type" : "address"
         }
      ],
      "name" : "getValidator",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "address payable",
                  "name" : "treasury",
                  "type" : "address"
               },
               {
                  "internalType" : "address",
                  "name" : "addr",
                  "type" : "address"
               },
               {
                  "internalType" : "string",
                  "name" : "enode",
                  "type" : "string"
               },
               {
                  "internalType" : "uint256",
                  "name" : "commissionRate",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "bondedStake",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "totalSlashed",
                  "type" : "uint256"
               },
               {
                  "internalType" : "contract Liquid",
                  "name" : "liquidContract",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "name" : "liquidSupply",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "registrationBlock",
                  "type" : "uint256"
               },
               {
                  "internalType" : "enum Autonity.ValidatorState",
                  "name" : "state",
                  "type" : "uint8"
               }
            ],
            "internalType" : "struct Autonity.Validator",
            "name" : "",
            "type" : "tuple"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getValidators",
      "outputs" : [
         {
            "internalType" : "address[]",
            "name" : "",
            "type" : "address[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getVersion",
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
      "name" : "headBondingID",
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
      "name" : "headUnbondingID",
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
      "name" : "lastEpochBlock",
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
            "name" : "_addr",
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
      "name" : "name",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "stateMutability" : "pure",
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
      "name" : "pauseValidator",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "string",
            "name" : "_enode",
            "type" : "string"
         },
         {
            "internalType" : "bytes",
            "name" : "_proof",
            "type" : "bytes"
         }
      ],
      "name" : "registerValidator",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "resetContractUpgrade",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
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
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_period",
            "type" : "uint256"
         }
      ],
      "name" : "setEpochPeriod",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_price",
            "type" : "uint256"
         }
      ],
      "name" : "setMinimumBaseFee",
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
         }
      ],
      "name" : "setOperatorAccount",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_account",
            "type" : "address"
         }
      ],
      "name" : "setTreasuryAccount",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_treasuryFee",
            "type" : "uint256"
         }
      ],
      "name" : "setTreasuryFee",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_period",
            "type" : "uint256"
         }
      ],
      "name" : "setUnbondingPeriod",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "symbol",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "stateMutability" : "pure",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "tailBondingID",
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
      "name" : "tailUnbondingID",
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
      "name" : "totalRedistributed",
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
            "name" : "_recipient",
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
            "name" : "sender",
            "type" : "address"
         },
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
      "name" : "transferFrom",
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
            "name" : "_validator",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "unbond",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "bytes",
            "name" : "_bytecode",
            "type" : "bytes"
         },
         {
            "internalType" : "string",
            "name" : "_abi",
            "type" : "string"
         }
      ],
      "name" : "upgradeContract",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "stateMutability" : "payable",
      "type" : "receive"
   }
]
`
