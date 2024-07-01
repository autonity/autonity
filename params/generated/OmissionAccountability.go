package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OmissionAccountabilityBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b506040516105d63803806105d683398101604081905261002f91610091565b601280546001600160a01b0319166001600160a01b0393909316929092179091558051600a556020810151600b556040810151600c556060810151600d556080810151600e5560a0810151600f5560c081015160105560e00151601155610161565b6000808284036101208112156100a657600080fd5b83516001600160a01b03811681146100bd57600080fd5b9250610100601f1982018113156100d357600080fd5b60405191508082016001600160401b038111838210171561010457634e487b7160e01b600052604160045260246000fd5b80604052506020850151825260408501516020830152606085015160408301526080850151606083015260a0850151608083015260c085015160a083015260e085015160c08301528085015160e083015250809150509250929050565b610466806101706000396000f3fe608060405234801561001057600080fd5b50600436106100bd5760003560e01c80638c80e31d11610076578063904383d11161005b578063904383d1146101bf578063eb231a1a146101df578063fd806677146101ff57600080fd5b80638c80e31d146101955780638d0c2a311461019f57600080fd5b806337c3a9b1116100a757806337c3a9b1146101155780637516f5801461011f57806379502c551461013457600080fd5b8062d049f3146100c2578063048a620c146100f5575b600080fd5b6100e26100d03660046102ce565b60086020526000908152604090205481565b6040519081526020015b60405180910390f35b6100e26101033660046102ce565b60046020526000908152604090205481565b6000546100e29081565b61013261012d36600461033a565b61021f565b005b600a54600b54600c54600d54600e54600f5460105460115461015a979695949392919088565b604080519889526020890197909752958701949094526060860192909252608085015260a084015260c083015260e0820152610100016100ec565b6002546100e29081565b6100e26101ad3660046102ce565b60066020526000908152604090205481565b6100e26101cd3660046102ce565b60096020526000908152604090205481565b6100e26101ed3660046102ce565b60076020526000908152604090205481565b6100e261020d3660046102ce565b60056020526000908152604090205481565b60125473ffffffffffffffffffffffffffffffffffffffff1633146102ca576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e74726163740000000000000000000000000000000000000000606482015260840160405180910390fd5b5050565b6000602082840312156102e057600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461030457600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561034d57600080fd5b8235801515811461035d57600080fd5b915060208381013567ffffffffffffffff8082111561037b57600080fd5b818601915086601f83011261038f57600080fd5b8135818111156103a1576103a161030b565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156103e4576103e461030b565b60405291825284820192508381018501918983111561040257600080fd5b938501935b8285101561042057843584529385019392850192610407565b809650505050505050925092905056fea26469706673582212207bf66b9b7f6030c1c0bf24e9944c2dcaff60b82fcc25c1c2910006aabfdf9f2e64736f6c63430008150033")

var OmissionAccountabilityAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_autonity",
            "type" : "address"
         },
         {
            "components" : [
               {
                  "internalType" : "uint256",
                  "name" : "negligibleThreshold",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "omissionLoopBackWindow",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "activityProofRewardRate",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "maxCommitteeSize",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "pastPerformanceWeight",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "initialJailingPeriod",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "initialProbationPeriod",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "initialSlashingRate",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct OmissionAccountability.Config",
            "name" : "_config",
            "type" : "tuple"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "constructor"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ],
      "name" : "activityPercentage",
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
      "name" : "config",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "negligibleThreshold",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "omissionLoopBackWindow",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "activityProofRewardRate",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "maxCommitteeSize",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "pastPerformanceWeight",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "initialJailingPeriod",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "initialProbationPeriod",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "initialSlashingRate",
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
      "name" : "currentEpochInactivityScores",
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
      "name" : "epochInactiveBlocks",
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
      "name" : "epochProverEfforts",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "totalAccumulatedEfforts",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "bool",
            "name" : "isProposerOmissionFaulty",
            "type" : "bool"
         },
         {
            "internalType" : "uint256[]",
            "name" : "ids",
            "type" : "uint256[]"
         }
      ],
      "name" : "finalize",
      "outputs" : [],
      "stateMutability" : "nonpayable",
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
      "name" : "lastEpochInactivityScores",
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
      "name" : "lookBackWindow",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "start",
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
      "name" : "overallFaultyBlocks",
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
      "name" : "repeatedOffences",
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
