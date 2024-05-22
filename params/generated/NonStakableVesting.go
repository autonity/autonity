package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var NonStakableVestingBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b50604051620016cb380380620016cb8339810160408190526100319161007a565b600280546001600160a01b039384166001600160a01b031991821617909155600380549290931691161790556100b4565b6001600160a01b038116811461007757600080fd5b50565b6000806040838503121561008d57600080fd5b825161009881610062565b60208401519092506100a981610062565b809150509250929050565b61160780620000c46000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c8063472c513a1161008157806376461cea1161005b57806376461cea146101c8578063a17a5027146101db578063f968f493146101ee57600080fd5b8063472c513a146101525780636c04e040146101885780637264c4da146101a857600080fd5b806321ec4487116100b257806321ec4487146100f65780632afbbacb1461011c578063309681ef1461013f57600080fd5b806307ae499f146100ce5780630a30959d146100e3575b600080fd5b6100e16100dc366004611264565b6101f6565b005b6100e16100f1366004611286565b610271565b6101096101043660046112c8565b610296565b6040519081526020015b60405180910390f35b61012f61012a3660046112c8565b6102b3565b6040519015158152602001610113565b6100e161014d3660046112f2565b6102ef565b610109610160366004611325565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b61019b610196366004611325565b61041b565b6040516101139190611340565b6101bb6101b63660046112c8565b61059c565b60405161011391906113c5565b6100e16101d636600461140a565b610655565b6100e16101e9366004611436565b6107e5565b610109610857565b6000610202338461099d565b905061020d81610a5b565b8211156102615760405162461bcd60e51b815260206004820152601960248201527f6e6f7420656e6f75676820756e6c6f636b65642066756e64730000000000000060448201526064015b60405180910390fd5b61026b8183610aa5565b50505050565b600061027d338361099d565b90506102918161028c83610a5b565b610aa5565b505050565b60006102aa6102a5848461099d565b610a5b565b90505b92915050565b600060016102c1848461099d565b815481106102d1576102d1611472565b600091825260209091206006909102016005015460ff169392505050565b60035473ffffffffffffffffffffffffffffffffffffffff1633146103565760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610258565b60045481106103a75760405162461bcd60e51b815260206004820152601660248201527f696e76616c6964207363686564756c6520636c617373000000000000000000006044820152606401610258565b6000600482815481106103bc576103bc611472565b9060005260206000209060060201905060006103e985858460000154856001015486600201546000610b11565b9050838260030160008282546103ff91906114d0565b9091555050600090815260056020526040902091909155505050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260208190526040812080546060929067ffffffffffffffff81111561045e5761045e6114e3565b6040519080825280602002602001820160405280156104ca57816020015b6104b76040518060c0016040528060008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b81526020019060019003908161047c5790505b50905060005b81518110156105945760018382815481106104ed576104ed611472565b90600052602060002001548154811061050857610508611472565b60009182526020918290206040805160c081018252600690930290910180548352600181015493830193909352600283015490820152600382015460608201526004820154608082015260059091015460ff16151560a0820152825183908390811061057657610576611472565b6020026020010181905250808061058c90611512565b9150506104d0565b509392505050565b6105d76040518060c0016040528060008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b60016105e3848461099d565b815481106105f3576105f3611472565b60009182526020918290206040805160c081018252600690930290910180548352600181015493830193909352600283015490820152600382015460608201526004820154608082015260059091015460ff16151560a0820152905092915050565b60035473ffffffffffffffffffffffffffffffffffffffff1633146106bc5760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610258565b6040805160c0810182529384526020840192835283019081526000606084018181526080850182815260a0860183815260048054600181018255945295517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b60069094029384015593517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19c83015591517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19d82015590517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19e82015590517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19f82015590517f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd1a090910155565b60035473ffffffffffffffffffffffffffffffffffffffff16331461084c5760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610258565b610291838383610ddd565b60025460009073ffffffffffffffffffffffffffffffffffffffff1633146108e75760405162461bcd60e51b815260206004820152602860248201527f66756e6374696f6e207265737472696374656420746f204175746f6e6974792060448201527f636f6e74726163740000000000000000000000000000000000000000000000006064820152608401610258565b4260005b6004548110156109985760006004828154811061090a5761090a611472565b906000526020600020906006020190508281600101541180610933575080600401548160030154145b1561093e5750610986565b828160050181905550600061096182600001548360020154868560030154610df6565b9050816004015481610973919061154a565b61097d90866114d0565b60049092015592505b8061099081611512565b9150506108eb565b505090565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602081905260408120548210610a115760405162461bcd60e51b815260206004820152601360248201527f696e76616c6964207363686564756c65206964000000000000000000000000006044820152606401610258565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602081905260409020805483908110610a4857610a48611472565b9060005260206000200154905092915050565b60006102ad82610a6a84610e39565b600085815260056020526040902054600480549091908110610a8e57610a8e611472565b906000526020600020906006020160050154610e7a565b60008060018481548110610abb57610abb611472565b906000526020600020906006020190508060000154831115610af9578054610ae3908461154a565b9150610af484338360000154610f3c565b610b0a565b8215610b0a57610b0a843385610f3c565b5092915050565b600042851015610b895760405162461bcd60e51b815260206004820152602860248201527f7363686564756c652063616e6e6f74207374617274206265666f72652063726560448201527f6174696e672069740000000000000000000000000000000000000000000000006064820152608401610258565b84841015610bff5760405162461bcd60e51b815260206004820152602c60248201527f636c696666206d7573742062652067726561746572207468616e206f7220657160448201527f75616c20746f20737461727400000000000000000000000000000000000000006064820152608401610258565b838311610c4e5760405162461bcd60e51b815260206004820152601e60248201527f656e64206d7573742062652067726561746572207468616e20636c69666600006044820152606401610258565b50600180546040805160c08101825297885260006020808a018281528a8401998a5260608b0198895260808b0197885295151560a08b01908152848601865585835299517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6600686029081019190915595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf787015597517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf886015595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf985015593517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfa84015595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfb90920180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169215159290921790915573ffffffffffffffffffffffffffffffffffffffff95909516825281835281208054948501815581522090910181905590565b6000610de9848461099d565b905061026b818584610fa0565b6000838310610e06575080610e31565b610e10858561154a565b610e1a868561154a565b610e24908461155d565b610e2e9190611574565b90505b949350505050565b60008060018381548110610e4f57610e4f611472565b9060005260206000209060060201905080600101548160000154610e7391906114d0565b9392505050565b60008060018581548110610e9057610e90611472565b906000526020600020906006020190508060030154831015610ef45760405162461bcd60e51b815260206004820152601c60248201527f636c69666620706572696f64206e6f74207265616368656420796574000000006044820152606401610258565b6000610f0a826002015483600401548688610df6565b90508160010154811115610f30576001820154610f27908261154a565b92505050610e73565b50600095945050505050565b600060018481548110610f5157610f51611472565b9060005260206000209060060201905081816000016000828254610f75919061154a565b9250508190555081816001016000828254610f9091906114d0565b9091555061026b90508383611117565b73ffffffffffffffffffffffffffffffffffffffff821660009081526020819052604081208054909190610fd69060019061154a565b67ffffffffffffffff811115610fee57610fee6114e3565b604051908082528060200260200182016040528015611017578160200160208202803683370190505b5090506000805b83548110156110a5578684828154811061103a5761103a611472565b906000526020600020015403156110935783818154811061105d5761105d611472565b906000526020600020015483838061107490611512565b94508151811061108657611086611472565b6020026020010181815250505b8061109d81611512565b91505061101e565b5073ffffffffffffffffffffffffffffffffffffffff851660009081526020818152604090912083516110da92850190611204565b5050505073ffffffffffffffffffffffffffffffffffffffff16600090815260208181526040822080546001810182559083529120019190915550565b6002546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015260248201849052600092169063a9059cbb906044016020604051808303816000875af1158015611191573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111b591906115af565b9050806102915760405162461bcd60e51b815260206004820152601360248201527f4e544e206e6f74207472616e73666572726564000000000000000000000000006044820152606401610258565b82805482825590600052602060002090810192821561123f579160200282015b8281111561123f578251825591602001919060010190611224565b5061124b92915061124f565b5090565b5b8082111561124b5760008155600101611250565b6000806040838503121561127757600080fd5b50508035926020909101359150565b60006020828403121561129857600080fd5b5035919050565b803573ffffffffffffffffffffffffffffffffffffffff811681146112c357600080fd5b919050565b600080604083850312156112db57600080fd5b6112e48361129f565b946020939093013593505050565b60008060006060848603121561130757600080fd5b6113108461129f565b95602085013595506040909401359392505050565b60006020828403121561133757600080fd5b6102aa8261129f565b6020808252825182820181905260009190848201906040850190845b818110156113b9576113a6838551805182526020810151602083015260408101516040830152606081015160608301526080810151608083015260a0810151151560a08301525050565b9284019260c0929092019160010161135c565b50909695505050505050565b60c081016102ad8284805182526020810151602083015260408101516040830152606081015160608301526080810151608083015260a0810151151560a08301525050565b60008060006060848603121561141f57600080fd5b505081359360208301359350604090920135919050565b60008060006060848603121561144b57600080fd5b6114548461129f565b9250602084013591506114696040850161129f565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156102ad576102ad6114a1565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611543576115436114a1565b5060010190565b818103818111156102ad576102ad6114a1565b80820281158282048414176102ad576102ad6114a1565b6000826115aa577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000602082840312156115c157600080fd5b81518015158114610e7357600080fdfea264697066735822122022e98020d1c10fb9bd48c40f22858d814ef3d9787f59b7639985f73b1355495d64736f6c63430008150033")

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
