package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OmissionAccountabilityBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b5060405161055038038061055083398101604081905261002f91610089565b601280546001600160a01b0319166001600160a01b0393909316929092179091558051600b556020810151600c556040810151600d556060810151600e556080810151600f5560a081015160105560c0015160115561014d565b60008082840361010081121561009e57600080fd5b83516001600160a01b03811681146100b557600080fd5b925060e0601f19820112156100c957600080fd5b5060405160e081016001600160401b03811182821017156100fa57634e487b7160e01b600052604160045260246000fd5b80604052506020840151815260408401516020820152606084015160408201526080840151606082015260a0840151608082015260c084015160a082015260e084015160c0820152809150509250929050565b6103f48061015c6000396000f3fe608060405234801561001057600080fd5b50600436106100a25760003560e01c80638c80e31d11610076578063904383d11161005b578063904383d114610163578063eb231a1a14610183578063f2d3f970146101a357600080fd5b80638c80e31d146101505780638d0f2bf71461015a57600080fd5b8062d049f3146100a757806337c3a9b1146100da5780637516f580146100e457806379502c55146100f9575b600080fd5b6100c76100b536600461025c565b60096020526000908152604090205481565b6040519081526020015b60405180910390f35b6000546100c79081565b6100f76100f23660046102c8565b6101ad565b005b600b54600c54600d54600e54600f5460105460115461011b9695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e0016100d1565b6002546100c79081565b6100c760045481565b6100c761017136600461025c565b600a6020526000908152604090205481565b6100c761019136600461025c565b60086020526000908152604090205481565b6005546100c79081565b60125473ffffffffffffffffffffffffffffffffffffffff163314610258576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e74726163740000000000000000000000000000000000000000606482015260840160405180910390fd5b5050565b60006020828403121561026e57600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461029257600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080604083850312156102db57600080fd5b823580151581146102eb57600080fd5b915060208381013567ffffffffffffffff8082111561030957600080fd5b818601915086601f83011261031d57600080fd5b81358181111561032f5761032f610299565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561037257610372610299565b60405291825284820192508381018501918983111561039057600080fd5b938501935b828510156103ae57843584529385019392850192610395565b809650505050505050925092905056fea2646970667358221220e45029f967aec240668cd5cb3ec28a6d32909811ff0bc885d422377b6c12965664736f6c63430008150033")

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
      "inputs" : [],
      "name" : "epochCollusionDegree",
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
      "name" : "epochInactivityScores",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "startEpochID",
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
