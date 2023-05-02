package generated

const Abi = `[
   {
      "inputs" : [
         {
            "internalType" : "address[]",
            "name" : "_voters",
            "type" : "address[]"
         },
         {
            "internalType" : "address",
            "name" : "_autonity",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_operator",
            "type" : "address"
         },
         {
            "internalType" : "string[]",
            "name" : "_symbols",
            "type" : "string[]"
         },
         {
            "internalType" : "uint256",
            "name" : "_votePeriod",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "constructor"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_round",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_height",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_timestamp",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_votePeriod",
            "type" : "uint256"
         }
      ],
      "name" : "NewRound",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "string[]",
            "name" : "_symbols",
            "type" : "string[]"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_round",
            "type" : "uint256"
         }
      ],
      "name" : "NewSymbols",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : true,
            "internalType" : "address",
            "name" : "_voter",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "int256[]",
            "name" : "_votes",
            "type" : "int256[]"
         }
      ],
      "name" : "Voted",
      "type" : "event"
   },
   {
      "stateMutability" : "payable",
      "type" : "fallback"
   },
   {
      "inputs" : [],
      "name" : "finalize",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getPrecision",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "pure",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getRound",
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
            "name" : "_round",
            "type" : "uint256"
         },
         {
            "internalType" : "string",
            "name" : "_symbol",
            "type" : "string"
         }
      ],
      "name" : "getRoundData",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "uint256",
                  "name" : "round",
                  "type" : "uint256"
               },
               {
                  "internalType" : "int256",
                  "name" : "price",
                  "type" : "int256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "timestamp",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "status",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct IOracle.RoundData",
            "name" : "data",
            "type" : "tuple"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getSymbols",
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
      "name" : "getVotePeriod",
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
      "name" : "getVoters",
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
      "name" : "lastRoundBlock",
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
      "name" : "lastVoterUpdateRound",
      "outputs" : [
         {
            "internalType" : "int256",
            "name" : "",
            "type" : "int256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "string",
            "name" : "_symbol",
            "type" : "string"
         }
      ],
      "name" : "latestRoundData",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "uint256",
                  "name" : "round",
                  "type" : "uint256"
               },
               {
                  "internalType" : "int256",
                  "name" : "price",
                  "type" : "int256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "timestamp",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "status",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct IOracle.RoundData",
            "name" : "data",
            "type" : "tuple"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "name" : "newSymbols",
      "outputs" : [
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
      "inputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         },
         {
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ],
      "name" : "reports",
      "outputs" : [
         {
            "internalType" : "int256",
            "name" : "",
            "type" : "int256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "round",
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
            "name" : "_operator",
            "type" : "address"
         }
      ],
      "name" : "setOperator",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "string[]",
            "name" : "_symbols",
            "type" : "string[]"
         }
      ],
      "name" : "setSymbols",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address[]",
            "name" : "_newVoters",
            "type" : "address[]"
         }
      ],
      "name" : "setVoters",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "symbolUpdatedRound",
      "outputs" : [
         {
            "internalType" : "int256",
            "name" : "",
            "type" : "int256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "name" : "symbols",
      "outputs" : [
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
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_commit",
            "type" : "uint256"
         },
         {
            "internalType" : "int256[]",
            "name" : "_reports",
            "type" : "int256[]"
         },
         {
            "internalType" : "uint256",
            "name" : "_salt",
            "type" : "uint256"
         }
      ],
      "name" : "vote",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "votePeriod",
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
            "name" : "",
            "type" : "address"
         }
      ],
      "name" : "votingInfo",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "round",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "commit",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "isVoter",
            "type" : "bool"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "stateMutability" : "payable",
      "type" : "receive"
   }
]
`
