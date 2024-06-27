package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OmissionAccountabilityBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b5060405161041138038061041183398101604081905261002f91610089565b600780546001600160a01b0319166001600160a01b0393909316929092179091558051600055602081015160015560408101516002556060810151600355608081015160045560a081015160055560c0015160065561014d565b60008082840361010081121561009e57600080fd5b83516001600160a01b03811681146100b557600080fd5b925060e0601f19820112156100c957600080fd5b5060405160e081016001600160401b03811182821017156100fa57634e487b7160e01b600052604160045260246000fd5b80604052506020840151815260408401516020820152606084015160408201526080840151606082015260a0840151608082015260c084015160a082015260e084015160c0820152809150509250929050565b6102b58061015c6000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80637516f5801461003b57806379502c5514610050575b600080fd5b61004e610049366004610189565b6100ab565b005b6000546001546002546003546004546005546006546100729695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e00160405180910390f35b60075473ffffffffffffffffffffffffffffffffffffffff163314610156576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e74726163740000000000000000000000000000000000000000606482015260840160405180910390fd5b5050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000806040838503121561019c57600080fd5b823580151581146101ac57600080fd5b915060208381013567ffffffffffffffff808211156101ca57600080fd5b818601915086601f8301126101de57600080fd5b8135818111156101f0576101f061015a565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f830116810181811085821117156102335761023361015a565b60405291825284820192508381018501918983111561025157600080fd5b938501935b8285101561026f57843584529385019392850192610256565b809650505050505050925092905056fea264697066735822122059af9c14b7753d747d7bb0c06f140844d9ba15a08288c18b50d2978cc23f38a764736f6c63430008150033")

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
   }
]
`))
