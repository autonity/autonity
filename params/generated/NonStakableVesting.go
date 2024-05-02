package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var NonStakableVestingBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b5060405162001786380380620017868339810160408190526100319161007a565b600280546001600160a01b039384166001600160a01b031991821617909155600380549290931691161790556100b4565b6001600160a01b038116811461007757600080fd5b50565b6000806040838503121561008d57600080fd5b825161009881610062565b60208401519092506100a981610062565b809150509250929050565b6116c280620000c46000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c8063472c513a1161008157806376461cea1161005b57806376461cea146101c8578063a17a5027146101db578063f968f493146101ee57600080fd5b8063472c513a146101525780636c04e040146101885780637264c4da146101a857600080fd5b806321ec4487116100b257806321ec4487146100f65780632afbbacb1461011c578063309681ef1461013f57600080fd5b806307ae499f146100ce5780630a30959d146100e3575b600080fd5b6100e16100dc36600461131f565b6101f6565b005b6100e16100f1366004611341565b61028b565b610109610104366004611383565b6102b0565b6040519081526020015b60405180910390f35b61012f61012a366004611383565b6102cd565b6040519015158152602001610113565b6100e161014d3660046113ad565b610309565b6101096101603660046113e0565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b61019b6101963660046113e0565b610469565b60405161011391906113fb565b6101bb6101b6366004611383565b6105ea565b6040516101139190611480565b6100e16101d63660046114c5565b6106a3565b6100e16101e93660046114f1565b61084d565b6101096108d9565b60006102023384610a39565b905061020d81610b11565b82111561027b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6e6f7420656e6f75676820756e6c6f636b65642066756e64730000000000000060448201526064015b60405180910390fd5b6102858183610b5b565b50505050565b60006102973383610a39565b90506102ab816102a683610b11565b610b5b565b505050565b60006102c46102bf8484610a39565b610b11565b90505b92915050565b600060016102db8484610a39565b815481106102eb576102eb61152d565b600091825260209091206006909102016005015460ff169392505050565b60035473ffffffffffffffffffffffffffffffffffffffff16331461038a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610272565b60045481106103f5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f696e76616c6964207363686564756c6520636c617373000000000000000000006044820152606401610272565b60006004828154811061040a5761040a61152d565b90600052602060002090600602019050600061043785858460000154856001015486600201546000610c32565b90508382600301600082825461044d919061158b565b9091555050600090815260056020526040902091909155505050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260208190526040812080546060929067ffffffffffffffff8111156104ac576104ac61159e565b60405190808252806020026020018201604052801561051857816020015b6105056040518060c0016040528060008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b8152602001906001900390816104ca5790505b50905060005b81518110156105e257600183828154811061053b5761053b61152d565b9060005260206000200154815481106105565761055661152d565b60009182526020918290206040805160c081018252600690930290910180548352600181015493830193909352600283015490820152600382015460608201526004820154608082015260059091015460ff16151560a082015282518390839081106105c4576105c461152d565b602002602001018190525080806105da906115cd565b91505061051e565b509392505050565b6106256040518060c0016040528060008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b60016106318484610a39565b815481106106415761064161152d565b60009182526020918290206040805160c081018252600690930290910180548352600181015493830193909352600283015490820152600382015460608201526004820154608082015260059091015460ff16151560a0820152905092915050565b60035473ffffffffffffffffffffffffffffffffffffffff163314610724576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610272565b6040805160c0810182529384526020840192835283019081526000606084018181526080850182815260a0860183815260048054600181018255945295517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b60069094029384015593517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19c83015591517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19d82015590517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19e82015590517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19f82015590517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd1a090910155565b60035473ffffffffffffffffffffffffffffffffffffffff1633146108ce576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610272565b6102ab838383610ebc565b60025460009073ffffffffffffffffffffffffffffffffffffffff163314610983576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602860248201527f66756e6374696f6e207265737472696374656420746f204175746f6e6974792060448201527f636f6e74726163740000000000000000000000000000000000000000000000006064820152608401610272565b4260005b600454811015610a34576000600482815481106109a6576109a661152d565b9060005260206000209060060201905082816001015411806109cf575080600401548160030154145b156109da5750610a22565b82816005018190555060006109fd82600001548360020154868560030154610ed5565b9050816004015481610a0f9190611605565b610a19908661158b565b60049092015592505b80610a2c816115cd565b915050610987565b505090565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602081905260408120548210610ac7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f696e76616c6964207363686564756c65206964000000000000000000000000006044820152606401610272565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602081905260409020805483908110610afe57610afe61152d565b9060005260206000200154905092915050565b60006102c782610b2084610f18565b600085815260056020526040902054600480549091908110610b4457610b4461152d565b906000526020600020906006020160050154610f59565b60008060018481548110610b7157610b7161152d565b906000526020600020906006020190504381600301541115610bef576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f636c69666620706572696f64206e6f74207265616368656420796574000000006044820152606401610272565b8054831115610c1a578054610c049084611605565b9150610c1584338360000154610fdd565b610c2b565b8215610c2b57610c2b843385610fdd565b5092915050565b600084841015610cc4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f636c696666206d7573742062652067726561746572207468616e206f7220657160448201527f75616c20746f20737461727400000000000000000000000000000000000000006064820152608401610272565b838311610d2d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f656e64206d7573742062652067726561746572207468616e20636c69666600006044820152606401610272565b50600180546040805160c08101825297885260006020808a018281528a8401998a5260608b0198895260808b0197885295151560a08b01908152848601865585835299517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6600686029081019190915595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf787015597517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf886015595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf985015593517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfa84015595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfb90920180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169215159290921790915573ffffffffffffffffffffffffffffffffffffffff95909516825281835281208054948501815581522090910181905590565b6000610ec88484610a39565b9050610285818584611041565b6000838310610ee5575080610f10565b610eef8585611605565b610ef98685611605565b610f039084611618565b610f0d919061162f565b90505b949350505050565b60008060018381548110610f2e57610f2e61152d565b9060005260206000209060060201905080600101548160000154610f52919061158b565b9392505050565b60008060018581548110610f6f57610f6f61152d565b906000526020600020906006020190508060030154831015610f95576000915050610f52565b6000610fab826002015483600401548688610ed5565b90508160010154811115610fd1576001820154610fc89082611605565b92505050610f52565b50600095945050505050565b600060018481548110610ff257610ff261152d565b90600052602060002090600602019050818160000160008282546110169190611605565b9250508190555081816001016000828254611031919061158b565b90915550610285905083836111b8565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602081905260408120805490919061107790600190611605565b67ffffffffffffffff81111561108f5761108f61159e565b6040519080825280602002602001820160405280156110b8578160200160208202803683370190505b5090506000805b835481101561114657868482815481106110db576110db61152d565b90600052602060002001540315611134578381815481106110fe576110fe61152d565b9060005260206000200154838380611115906115cd565b9450815181106111275761112761152d565b6020026020010181815250505b8061113e816115cd565b9150506110bf565b5073ffffffffffffffffffffffffffffffffffffffff8516600090815260208181526040909120835161117b928501906112bf565b5050505073ffffffffffffffffffffffffffffffffffffffff16600090815260208181526040822080546001810182559083529120019190915550565b6002546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015260248201849052600092169063a9059cbb906044016020604051808303816000875af1158015611232573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611256919061166a565b9050806102ab576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f4e544e206e6f74207472616e73666572726564000000000000000000000000006044820152606401610272565b8280548282559060005260206000209081019282156112fa579160200282015b828111156112fa5782518255916020019190600101906112df565b5061130692915061130a565b5090565b5b80821115611306576000815560010161130b565b6000806040838503121561133257600080fd5b50508035926020909101359150565b60006020828403121561135357600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461137e57600080fd5b919050565b6000806040838503121561139657600080fd5b61139f8361135a565b946020939093013593505050565b6000806000606084860312156113c257600080fd5b6113cb8461135a565b95602085013595506040909401359392505050565b6000602082840312156113f257600080fd5b6102c48261135a565b6020808252825182820181905260009190848201906040850190845b8181101561147457611461838551805182526020810151602083015260408101516040830152606081015160608301526080810151608083015260a0810151151560a08301525050565b9284019260c09290920191600101611417565b50909695505050505050565b60c081016102c78284805182526020810151602083015260408101516040830152606081015160608301526080810151608083015260a0810151151560a08301525050565b6000806000606084860312156114da57600080fd5b505081359360208301359350604090920135919050565b60008060006060848603121561150657600080fd5b61150f8461135a565b9250602084013591506115246040850161135a565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156102c7576102c761155c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036115fe576115fe61155c565b5060010190565b818103818111156102c7576102c761155c565b80820281158282048414176102c7576102c761155c565b600082611665577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60006020828403121561167c57600080fd5b81518015158114610f5257600080fdfea264697066735822122000afc6426089b162dd828e9563f31adcc7ea1fb4f62f5d3a68408f5e7a520b2464736f6c63430008150033")

var NonStakableVestingAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_autonity",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_operator",
            "type" : "address"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "constructor"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_beneficiary",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_id",
            "type" : "uint256"
         }
      ],
      "name" : "canStake",
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
            "internalType" : "address",
            "name" : "_beneficiary",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_id",
            "type" : "uint256"
         },
         {
            "internalType" : "address",
            "name" : "_recipient",
            "type" : "address"
         }
      ],
      "name" : "cancelSchedule",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_startTime",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_cliffTime",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_endTime",
            "type" : "uint256"
         }
      ],
      "name" : "createScheduleClass",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_beneficiary",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_id",
            "type" : "uint256"
         }
      ],
      "name" : "getSchedule",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "uint256",
                  "name" : "currentNTNAmount",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "withdrawnValue",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "start",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "cliff",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "end",
                  "type" : "uint256"
               },
               {
                  "internalType" : "bool",
                  "name" : "canStake",
                  "type" : "bool"
               }
            ],
            "internalType" : "struct ScheduleBase.Schedule",
            "name" : "",
            "type" : "tuple"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_beneficiary",
            "type" : "address"
         }
      ],
      "name" : "getSchedules",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "uint256",
                  "name" : "currentNTNAmount",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "withdrawnValue",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "start",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "cliff",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "end",
                  "type" : "uint256"
               },
               {
                  "internalType" : "bool",
                  "name" : "canStake",
                  "type" : "bool"
               }
            ],
            "internalType" : "struct ScheduleBase.Schedule[]",
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
            "name" : "_beneficiary",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_scheduleClass",
            "type" : "uint256"
         }
      ],
      "name" : "newSchedule",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_id",
            "type" : "uint256"
         }
      ],
      "name" : "releaseAllFunds",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_id",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "releaseFund",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_beneficiary",
            "type" : "address"
         }
      ],
      "name" : "totalSchedules",
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
      "name" : "unlockTokens",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "_totalNewUnlocked",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_beneficiary",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_id",
            "type" : "uint256"
         }
      ],
      "name" : "unlockedFunds",
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
