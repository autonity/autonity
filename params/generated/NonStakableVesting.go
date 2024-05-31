package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var NonStakableVestingBytecode = common.Hex2Bytes("60806040523480156200001157600080fd5b5060405162001a5938038062001a5983398101604081905262000034916200007f565b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055620000be565b6001600160a01b03811681146200007c57600080fd5b50565b600080604083850312156200009357600080fd5b8251620000a08162000066565b6020840151909250620000b38162000066565b809150509250929050565b61198b80620000ce6000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063635bf93311610097578063aad5572611610066578063aad5572614610224578063b0c9300914610244578063d934047d1461024d578063f968f4931461025657600080fd5b8063635bf933146101b55780637b8d4744146101eb578063995e21a4146101fe578063a9f45b621461021157600080fd5b8063213fe2b7116100d3578063213fe2b71461013557806321ec44871461015e5780632afbbacb1461017f5780634e974657146101a257600080fd5b806307ae499f146100fa5780630a30959d1461010f5780630b5d0e4214610122575b600080fd5b61010d6101083660046115e2565b610273565b005b61010d61011d366004611604565b6102ee565b61010d61013036600461161d565b610313565b610148610143366004611678565b61060b565b6040516101559190611693565b60405180910390f35b61017161016c366004611718565b61078c565b604051908152602001610155565b61019261018d366004611718565b6107a9565b6040519015158152602001610155565b61010d6101b0366004611604565b6107e5565b6101716101c3366004611678565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b61010d6101f9366004611742565b610851565b61010d61020c366004611604565b610aaf565b61010d61021f366004611775565b610b1b565b610237610232366004611718565b610b96565b60405161015591906117b1565b61017160045481565b61017160055481565b61025e610c4f565b60408051928352602083019190915201610155565b600061027f3384610e16565b905061028a81610ed4565b8211156102de5760405162461bcd60e51b815260206004820152601960248201527f6e6f7420656e6f75676820756e6c6f636b65642066756e64730000000000000060448201526064015b60405180910390fd5b6102e88183610f1e565b50505050565b60006102fa3383610e16565b905061030e8161030983610ed4565b610f1e565b505050565b60035473ffffffffffffffffffffffffffffffffffffffff16331461037a5760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064016102d5565b8360045410156103f25760405162461bcd60e51b815260206004820152602960248201527f6e6f7420656e6f7567682066756e647320746f206372656174652061206e657760448201527f207363686564756c65000000000000000000000000000000000000000000000060648201526084016102d5565b80600554101561046a5760405162461bcd60e51b815260206004820152603460248201527f7363686564756c6520746f74616c206475726174696f6e20657863656564732060448201527f6d617820616c6c6f776564206475726174696f6e00000000000000000000000060648201526084016102d5565b6040805161010081018252848152602081018481529181018381526060820187815260808301888152600060a0850181815260c0860182815260e087018381526006805460018101825590855297517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f60089099029889015597517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4088015594517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4187015592517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4286015590517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4385015590517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4484015590517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4583015591517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d469091015560048054869290610600908490611825565b909155505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260208190526040812080546060929067ffffffffffffffff81111561064e5761064e611838565b6040519080825280602002602001820160405280156106ba57816020015b6106a76040518060c0016040528060008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b81526020019060019003908161066c5790505b50905060005b81518110156107845760018382815481106106dd576106dd611867565b9060005260206000200154815481106106f8576106f8611867565b60009182526020918290206040805160c081018252600690930290910180548352600181015493830193909352600283015490820152600382015460608201526004820154608082015260059091015460ff16151560a0820152825183908390811061076657610766611867565b6020026020010181905250808061077c90611896565b9150506106c0565b509392505050565b60006107a061079b8484610e16565b610ed4565b90505b92915050565b600060016107b78484610e16565b815481106107c7576107c7611867565b600091825260209091206006909102016005015460ff169392505050565b60035473ffffffffffffffffffffffffffffffffffffffff16331461084c5760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064016102d5565b600555565b60035473ffffffffffffffffffffffffffffffffffffffff1633146108b85760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064016102d5565b60065481106109095760405162461bcd60e51b815260206004820152601360248201527f696e76616c6964207363686564756c652049440000000000000000000000000060448201526064016102d5565b60006006828154811061091e5761091e611867565b9060005260206000209060080201905082816004015410156109a85760405162461bcd60e51b815260206004820152603860248201527f6e6f7420656e6f7567682066756e647320746f206372656174652061206e657760448201527f20636f6e747261637420756e646572207363686564756c65000000000000000060648201526084016102d5565b60006109c585858460000154856001015486600201546000610f8a565b6000818152600760205260409020849055600183015483549192506109e9916118ce565b826007015410610a8f57600060018281548110610a0857610a08611867565b906000526020600020906006020190506000836004015484600601548360000154610a3391906118e1565b610a3d91906118f8565b905080846006016000828254610a539190611825565b9091555050815481908390600090610a6c908490611825565b9250508190555080826001016000828254610a8791906118ce565b909155505050505b83826004016000828254610aa39190611825565b90915550505050505050565b60035473ffffffffffffffffffffffffffffffffffffffff163314610b165760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064016102d5565b600455565b60035473ffffffffffffffffffffffffffffffffffffffff163314610b825760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064016102d5565b61030e610b8f8484610e16565b848361116a565b610bd16040518060c0016040528060008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b6001610bdd8484610e16565b81548110610bed57610bed611867565b60009182526020918290206040805160c081018252600690930290910180548352600181015493830193909352600283015490820152600382015460608201526004820154608082015260059091015460ff16151560a0820152905092915050565b600254600090819073ffffffffffffffffffffffffffffffffffffffff163314610ce15760405162461bcd60e51b815260206004820152602860248201527f66756e6374696f6e207265737472696374656420746f204175746f6e6974792060448201527f636f6e747261637400000000000000000000000000000000000000000000000060648201526084016102d5565b426000805b600654811015610e0357600060068281548110610d0557610d05611867565b906000526020600020906008020190508381600001548260010154610d2a91906118ce565b1180610d3d575080600501548160030154145b15610d485750610df1565b8381600701819055506000610d6b826000015483600201548785600301546112e1565b90508160050154811015610d80575060058101545b6005820154610d8f9082611825565b610d9990856118ce565b60058301829055825460028401546004850154929650610dba9288906112e1565b90508160060154811015610dcf575060068101545b6006820154610dde9082611825565b610de890876118ce565b60069092015593505b80610dfb81611896565b915050610ce6565b50610e0e8382611825565b935050509091565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602081905260408120548210610e8a5760405162461bcd60e51b815260206004820152601360248201527f696e76616c696420636f6e74726163742069640000000000000000000000000060448201526064016102d5565b73ffffffffffffffffffffffffffffffffffffffff83166000908152602081905260409020805483908110610ec157610ec1611867565b9060005260206000200154905092915050565b60006107a382610ee384611324565b600085815260076020526040902054600680549091908110610f0757610f07611867565b906000526020600020906008020160070154611365565b60008060018481548110610f3457610f34611867565b906000526020600020906006020190508060000154831115610f72578054610f5c9084611825565b9150610f6d84338360000154611436565b610f83565b8215610f8357610f83843385611436565b5092915050565b6000838311610fdb5760405162461bcd60e51b815260206004820152601e60248201527f656e64206d7573742062652067726561746572207468616e20636c696666000060448201526064016102d5565b50600180546040805160c08101825297885260006020808a018281528a8401998a5260608b0198895260808b0197885295151560a08b01908152848601865585835299517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6600686029081019190915595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf787015597517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf886015595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf985015593517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfa84015595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfb90920180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169215159290921790915573ffffffffffffffffffffffffffffffffffffffff95909516825281835281208054948501815581522090910181905590565b73ffffffffffffffffffffffffffffffffffffffff8216600090815260208190526040812080549091906111a090600190611825565b67ffffffffffffffff8111156111b8576111b8611838565b6040519080825280602002602001820160405280156111e1578160200160208202803683370190505b5090506000805b835481101561126f578684828154811061120457611204611867565b9060005260206000200154031561125d5783818154811061122757611227611867565b906000526020600020015483838061123e90611896565b94508151811061125057611250611867565b6020026020010181815250505b8061126781611896565b9150506111e8565b5073ffffffffffffffffffffffffffffffffffffffff851660009081526020818152604090912083516112a492850190611582565b5050505073ffffffffffffffffffffffffffffffffffffffff16600090815260208181526040822080546001810182559083529120019190915550565b60006112ed85856118ce565b83106112fa57508061131c565b836113058685611825565b61130f90846118e1565b61131991906118f8565b90505b949350505050565b6000806001838154811061133a5761133a611867565b906000526020600020906006020190508060010154816000015461135e91906118ce565b9392505050565b6000806001858154811061137b5761137b611867565b906000526020600020906006020190508060030154816002015461139f91906118ce565b8310156113ee5760405162461bcd60e51b815260206004820152601c60248201527f636c69666620706572696f64206e6f742072656163686564207965740000000060448201526064016102d5565b60006114048260020154836004015486886112e1565b9050816001015481111561142a5760018201546114219082611825565b9250505061135e565b50600095945050505050565b60006001848154811061144b5761144b611867565b906000526020600020906006020190508181600001600082825461146f9190611825565b925050819055508181600101600082825461148a91906118ce565b909155506102e8905083836002546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015260248201849052600092169063a9059cbb906044016020604051808303816000875af115801561150f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115339190611933565b90508061030e5760405162461bcd60e51b815260206004820152601360248201527f4e544e206e6f74207472616e736665727265640000000000000000000000000060448201526064016102d5565b8280548282559060005260206000209081019282156115bd579160200282015b828111156115bd5782518255916020019190600101906115a2565b506115c99291506115cd565b5090565b5b808211156115c957600081556001016115ce565b600080604083850312156115f557600080fd5b50508035926020909101359150565b60006020828403121561161657600080fd5b5035919050565b6000806000806080858703121561163357600080fd5b5050823594602084013594506040840135936060013592509050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461167357600080fd5b919050565b60006020828403121561168a57600080fd5b6107a08261164f565b6020808252825182820181905260009190848201906040850190845b8181101561170c576116f9838551805182526020810151602083015260408101516040830152606081015160608301526080810151608083015260a0810151151560a08301525050565b9284019260c092909201916001016116af565b50909695505050505050565b6000806040838503121561172b57600080fd5b6117348361164f565b946020939093013593505050565b60008060006060848603121561175757600080fd5b6117608461164f565b95602085013595506040909401359392505050565b60008060006060848603121561178a57600080fd5b6117938461164f565b9250602084013591506117a86040850161164f565b90509250925092565b60c081016107a38284805182526020810151602083015260408101516040830152606081015160608301526080810151608083015260a0810151151560a08301525050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156107a3576107a36117f6565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036118c7576118c76117f6565b5060010190565b808201808211156107a3576107a36117f6565b80820281158282048414176107a3576107a36117f6565b60008261192e577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b60006020828403121561194557600080fd5b8151801515811461135e57600080fdfea264697066735822122060a7c12eb89d4921ee5e980ca24db9d6aa9794c8976b3b21307174cf52e26c3664736f6c63430008150033")

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
      "name" : "changeContractBeneficiary",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_startTime",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_cliffDuration",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_totalDuration",
            "type" : "uint256"
         }
      ],
      "name" : "createSchedule",
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
      "name" : "getContract",
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
                  "name" : "cliffDuration",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "totalDuration",
                  "type" : "uint256"
               },
               {
                  "internalType" : "bool",
                  "name" : "canStake",
                  "type" : "bool"
               }
            ],
            "internalType" : "struct ContractBase.Contract",
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
      "name" : "getContracts",
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
                  "name" : "cliffDuration",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "totalDuration",
                  "type" : "uint256"
               },
               {
                  "internalType" : "bool",
                  "name" : "canStake",
                  "type" : "bool"
               }
            ],
            "internalType" : "struct ContractBase.Contract[]",
            "name" : "",
            "type" : "tuple[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "maxAllowedDuration",
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
            "name" : "_scheduleID",
            "type" : "uint256"
         }
      ],
      "name" : "newContract",
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
            "internalType" : "uint256",
            "name" : "_newMaxDuration",
            "type" : "uint256"
         }
      ],
      "name" : "setMaxAllowedDuration",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_totalNominal",
            "type" : "uint256"
         }
      ],
      "name" : "setTotalNominal",
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
      "name" : "totalContracts",
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
      "name" : "totalNominal",
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
            "name" : "_newUnlockedSubscribed",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_newUnlockedUnsubscribed",
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
