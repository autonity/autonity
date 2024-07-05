package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OmissionAccountabilityBytecode = common.Hex2Bytes("60806040523480156200001157600080fd5b5060405162000eee38038062000eee833981016040819052620000349162000221565b601480546001600160a01b0319166001600160a01b0385161790558051600c55602080820151600d556040820151600e556060820151600f55608082015160105560a082015160115560c082015160125560e08201516013558251620000a19160009190850190620000ab565b505050506200030c565b82805482825590600052602060002090810192821562000103579160200282015b828111156200010357825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620000cc565b506200011192915062000115565b5090565b5b8082111562000111576000815560010162000116565b6001600160a01b03811681146200014257600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b038111828210171562000186576200018662000145565b604052919050565b6000610100808385031215620001a357600080fd5b604051908101906001600160401b0382118183101715620001c857620001c862000145565b81604052809250835181526020840151602082015260408401516040820152606084015160608201526080840151608082015260a084015160a082015260c084015160c082015260e084015160e0820152505092915050565b600080600061014084860312156200023857600080fd5b835162000245816200012c565b602085810151919450906001600160401b03808211156200026557600080fd5b818701915087601f8301126200027a57600080fd5b8151818111156200028f576200028f62000145565b8060051b9150620002a28483016200015b565b818152918301840191848101908a841115620002bd57600080fd5b938501935b83851015620002eb5784519250620002da836200012c565b8282529385019390850190620002c2565b8097505050505050506200030385604086016200018e565b90509250925092565b610bd2806200031c6000396000f3fe608060405234801561001057600080fd5b50600436106100e95760003560e01c80638d0c2a311161008c578063e08b14ed11610066578063e08b14ed14610278578063eb231a1a1461028b578063f95bbd7f146102ab578063fd806677146102de57600080fd5b80638d0c2a3114610200578063904383d114610220578063c8425ac71461024057600080fd5b806354a2f945116100c857806354a2f945146101615780636a7ffba214610176578063754b1fd81461017f57806379502c551461019f57600080fd5b8062d049f3146100ee578063048a620c146101215780635426b5ea14610141575b600080fd5b61010e6100fc366004610880565b600a6020526000908152604090205481565b6040519081526020015b60405180910390f35b61010e61012f366004610880565b60066020526000908152604090205481565b61010e61014f366004610880565b60036020526000908152604090205481565b61017461016f3660046109a1565b6102fe565b005b61010e60045481565b61010e61018d366004610880565b60056020526000908152604090205481565b600c54600d54600e54600f546010546011546012546013546101c5979695949392919088565b604080519889526020890197909752958701949094526060860192909252608085015260a084015260c083015260e082015261010001610118565b61010e61020e366004610880565b60086020526000908152604090205481565b61010e61022e366004610880565b600b6020526000908152604090205481565b61025361024e366004610a2d565b6106b5565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610118565b610174610286366004610a4f565b6106fa565b61010e610299366004610880565b60096020526000908152604090205481565b6102ce6102b9366004610a8c565b60016020526000908152604090205460ff1681565b6040519015158152602001610118565b61010e6102ec366004610880565b60076020526000908152604090205481565b60145473ffffffffffffffffffffffffffffffffffffffff1633146103aa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e7472616374000000000000000000000000000000000000000060648201526084015b60405180910390fd5b60006103b68343610ad4565b9050841561042e57600081815260016020818152604080842080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690931790925573ffffffffffffffffffffffffffffffffffffffff8a16835260039052812080549161042483610aed565b91905055506106ab565b600081815260016020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905560028252909120895161047a928b01906107b8565b50600d546104889085610b25565b8111156106575760005b885181101561065557600d5460019060006104ad8386610ad4565b90505b6104ba8286610ad4565b81106105ca5760008181526001602052604090205460ff16156104e957816104e181610aed565b9250506105b8565b8781036104f957600092506105ca565b6000805b6000838152600260205260409020548110156105a657600083815260026020526040902080548290811061053357610533610b38565b6000918252602090912001548e5173ffffffffffffffffffffffffffffffffffffffff909116908f908890811061056c5761056c610b38565b602002602001015173ffffffffffffffffffffffffffffffffffffffff160361059457600191505b8061059e81610aed565b9150506104fd565b50806105b65760009350506105ca565b505b806105c281610b67565b9150506104b0565b50811561064057600360008c85815181106105e7576105e7610b38565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600081548092919061063a90610aed565b91905055505b5050808061064d90610aed565b915050610492565b505b73ffffffffffffffffffffffffffffffffffffffff87166000908152600560205260408120805488929061068c908490610b25565b9250508190555085600460008282546106a59190610b25565b90915550505b5050505050505050565b600260205281600052604060002081815481106106d157600080fd5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff169150829050565b60145473ffffffffffffffffffffffffffffffffffffffff1633146107a1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e7472616374000000000000000000000000000000000000000060648201526084016103a1565b80516107b49060009060208401906107b8565b5050565b828054828255906000526020600020908101928215610832579160200282015b8281111561083257825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff9091161782556020909201916001909101906107d8565b5061083e929150610842565b5090565b5b8082111561083e5760008155600101610843565b803573ffffffffffffffffffffffffffffffffffffffff8116811461087b57600080fd5b919050565b60006020828403121561089257600080fd5b61089b82610857565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f8301126108e257600080fd5b8135602067ffffffffffffffff808311156108ff576108ff6108a2565b8260051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f83011681018181108482111715610942576109426108a2565b60405293845285810183019383810192508785111561096057600080fd5b83870191505b848210156109865761097782610857565b83529183019190830190610966565b979650505050505050565b8035801515811461087b57600080fd5b600080600080600080600060e0888a0312156109bc57600080fd5b873567ffffffffffffffff8111156109d357600080fd5b6109df8a828b016108d1565b9750506109ee60208901610857565b955060408801359450610a0360608901610991565b93506080880135925060a08801359150610a1f60c08901610991565b905092959891949750929550565b60008060408385031215610a4057600080fd5b50508035926020909101359150565b600060208284031215610a6157600080fd5b813567ffffffffffffffff811115610a7857600080fd5b610a84848285016108d1565b949350505050565b600060208284031215610a9e57600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b81810381811115610ae757610ae7610aa5565b92915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610b1e57610b1e610aa5565b5060010190565b80820180821115610ae757610ae7610aa5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081610b7657610b76610aa5565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff019056fea26469706673582212207eb57c646af80e22f782668df2c1cc406610c08355afbc30e1815b94761e162e64736f6c63430008150033")

var OmissionAccountabilityAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_autonity",
            "type" : "address"
         },
         {
            "internalType" : "address[]",
            "name" : "_committee",
            "type" : "address[]"
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
                  "name" : "omissionLookBackWindow",
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
            "name" : "omissionLookBackWindow",
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
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "name" : "faultyProposers",
      "outputs" : [
         {
            "internalType" : "bool",
            "name" : "",
            "type" : "bool"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address[]",
            "name" : "absentees",
            "type" : "address[]"
         },
         {
            "internalType" : "address",
            "name" : "proposer",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "proposerEffort",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "isProposerOmissionFaulty",
            "type" : "bool"
         },
         {
            "internalType" : "uint256",
            "name" : "lastEpochBlock",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "delta",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "epochEnded",
            "type" : "bool"
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
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "name" : "inactiveValidators",
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
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ],
      "name" : "inactivityCounter",
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
      "name" : "proverEfforts",
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
   },
   {
      "inputs" : [
         {
            "internalType" : "address[]",
            "name" : "_committee",
            "type" : "address[]"
         }
      ],
      "name" : "setCommittee",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "totalAccumulatedEffort",
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
