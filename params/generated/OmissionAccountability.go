package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OmissionAccountabilityBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b5060405161021838038061021883398101604081905261002f91610089565b600780546001600160a01b0319166001600160a01b0393909316929092179091558051600055602081015160015560408101516002556060810151600355608081015160045560a081015160055560c0015160065561014d565b60008082840361010081121561009e57600080fd5b83516001600160a01b03811681146100b557600080fd5b925060e0601f19820112156100c957600080fd5b5060405160e081016001600160401b03811182821017156100fa57634e487b7160e01b600052604160045260246000fd5b80604052506020840151815260408401516020820152606084015160408201526080840151606082015260a0840151608082015260c084015160a082015260e084015160c0820152809150509250929050565b60bd8061015b6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806379502c5514602d575b600080fd5b600054600154600254600354600454600554600654604e9695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e00160405180910390f3fea2646970667358221220e4748475688f04a532bd2c78b0418ef8a6acd9293b545a682f3e7ad09edc3f5e64736f6c63430008150033")

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
   }
]
`))
