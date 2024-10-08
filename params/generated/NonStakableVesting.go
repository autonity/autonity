package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var NonStakableVestingBytecode = common.Hex2Bytes("60806040523480156200001157600080fd5b5060405162001bf538038062001bf583398101604081905262000034916200007f565b600280546001600160a01b039384166001600160a01b03199182161790915560038054929093169116179055620000be565b6001600160a01b03811681146200007c57600080fd5b50565b600080604083850312156200009357600080fd5b8251620000a08162000066565b6020840151909250620000b38162000066565b809150509250929050565b611b2780620000ce6000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80637b8d474411610097578063b0c9300911610066578063b0c930091461024f578063c5ca93a714610258578063d934047d146102cc578063f968f493146102d557600080fd5b80637b8d4744146101f6578063995e21a414610209578063a9f45b621461021c578063aad557261461022f57600080fd5b806321ec4487116100d357806321ec4487146101695780632afbbacb1461018a5780634e974657146101ad578063635bf933146101c057600080fd5b806307ae499f146101055780630a30959d1461011a5780630b5d0e421461012d578063213fe2b714610140575b600080fd5b61011861011336600461177e565b6102f2565b005b6101186101283660046117a0565b61036d565b61011861013b3660046117b9565b610392565b61015361014e366004611814565b61068a565b604051610160919061182f565b60405180910390f35b61017c6101773660046118b4565b61080b565b604051908152602001610160565b61019d6101983660046118b4565b610828565b6040519015158152602001610160565b6101186101bb3660046117a0565b610864565b61017c6101ce366004611814565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b6101186102043660046118de565b6108d0565b6101186102173660046117a0565b610b2e565b61011861022a366004611911565b610b9a565b61024261023d3660046118b4565b610c15565b604051610160919061194d565b61017c60045481565b61026b6102663660046117a0565b610cce565b6040516101609190600061010082019050825182526020830151602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015260c083015160c083015260e083015160e083015292915050565b61017c60055481565b6102dd610deb565b60408051928352602083019190915201610160565b60006102fe3384610fb2565b905061030981611070565b82111561035d5760405162461bcd60e51b815260206004820152601960248201527f6e6f7420656e6f75676820756e6c6f636b65642066756e64730000000000000060448201526064015b60405180910390fd5b61036781836110ba565b50505050565b60006103793383610fb2565b905061038d8161038883611070565b6110ba565b505050565b60035473ffffffffffffffffffffffffffffffffffffffff1633146103f95760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610354565b8360045410156104715760405162461bcd60e51b815260206004820152602960248201527f6e6f7420656e6f7567682066756e647320746f206372656174652061206e657760448201527f207363686564756c6500000000000000000000000000000000000000000000006064820152608401610354565b8060055410156104e95760405162461bcd60e51b815260206004820152603460248201527f7363686564756c6520746f74616c206475726174696f6e20657863656564732060448201527f6d617820616c6c6f776564206475726174696f6e0000000000000000000000006064820152608401610354565b6040805161010081018252848152602081018481529181018381526060820187815260808301888152600060a0850181815260c0860182815260e087018381526006805460018101825590855297517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f60089099029889015597517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4088015594517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4187015592517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4286015590517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4385015590517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4484015590517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d4583015591517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d46909101556004805486929061067f9084906119c1565b909155505050505050565b73ffffffffffffffffffffffffffffffffffffffff8116600090815260208190526040812080546060929067ffffffffffffffff8111156106cd576106cd6119d4565b60405190808252806020026020018201604052801561073957816020015b6107266040518060c0016040528060008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b8152602001906001900390816106eb5790505b50905060005b815181101561080357600183828154811061075c5761075c611a03565b90600052602060002001548154811061077757610777611a03565b60009182526020918290206040805160c081018252600690930290910180548352600181015493830193909352600283015490820152600382015460608201526004820154608082015260059091015460ff16151560a082015282518390839081106107e5576107e5611a03565b602002602001018190525080806107fb90611a32565b91505061073f565b509392505050565b600061081f61081a8484610fb2565b611070565b90505b92915050565b600060016108368484610fb2565b8154811061084657610846611a03565b600091825260209091206006909102016005015460ff169392505050565b60035473ffffffffffffffffffffffffffffffffffffffff1633146108cb5760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610354565b600555565b60035473ffffffffffffffffffffffffffffffffffffffff1633146109375760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610354565b60065481106109885760405162461bcd60e51b815260206004820152601360248201527f696e76616c6964207363686564756c65204944000000000000000000000000006044820152606401610354565b60006006828154811061099d5761099d611a03565b906000526020600020906008020190508281600401541015610a275760405162461bcd60e51b815260206004820152603860248201527f6e6f7420656e6f7567682066756e647320746f206372656174652061206e657760448201527f20636f6e747261637420756e646572207363686564756c6500000000000000006064820152608401610354565b6000610a4485858460000154856001015486600201546000611126565b600081815260076020526040902084905560018301548354919250610a6891611a6a565b826007015410610b0e57600060018281548110610a8757610a87611a03565b906000526020600020906006020190506000836004015484600601548360000154610ab29190611a7d565b610abc9190611a94565b905080846006016000828254610ad291906119c1565b9091555050815481908390600090610aeb9084906119c1565b9250508190555080826001016000828254610b069190611a6a565b909155505050505b83826004016000828254610b2291906119c1565b90915550505050505050565b60035473ffffffffffffffffffffffffffffffffffffffff163314610b955760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610354565b600455565b60035473ffffffffffffffffffffffffffffffffffffffff163314610c015760405162461bcd60e51b815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f720000000000006044820152606401610354565b61038d610c0e8484610fb2565b8483611306565b610c506040518060c0016040528060008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b6001610c5c8484610fb2565b81548110610c6c57610c6c611a03565b60009182526020918290206040805160c081018252600690930290910180548352600181015493830193909352600283015490820152600382015460608201526004820154608082015260059091015460ff16151560a0820152905092915050565b610d1660405180610100016040528060008152602001600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b6006548210610d675760405162461bcd60e51b815260206004820152601760248201527f7363686564756c6520646f6573206e6f742065786973740000000000000000006044820152606401610354565b60068281548110610d7a57610d7a611a03565b906000526020600020906008020160405180610100016040529081600082015481526020016001820154815260200160028201548152602001600382015481526020016004820154815260200160058201548152602001600682015481526020016007820154815250509050919050565b600254600090819073ffffffffffffffffffffffffffffffffffffffff163314610e7d5760405162461bcd60e51b815260206004820152602860248201527f66756e6374696f6e207265737472696374656420746f204175746f6e6974792060448201527f636f6e74726163740000000000000000000000000000000000000000000000006064820152608401610354565b426000805b600654811015610f9f57600060068281548110610ea157610ea1611a03565b906000526020600020906008020190508381600001548260010154610ec69190611a6a565b1180610ed9575080600501548160030154145b15610ee45750610f8d565b8381600701819055506000610f078260000154836002015487856003015461147d565b90508160050154811015610f1c575060058101545b6005820154610f2b90826119c1565b610f359085611a6a565b60058301829055825460028401546004850154929650610f5692889061147d565b90508160060154811015610f6b575060068101545b6006820154610f7a90826119c1565b610f849087611a6a565b60069092015593505b80610f9781611a32565b915050610e82565b50610faa83826119c1565b935050509091565b73ffffffffffffffffffffffffffffffffffffffff821660009081526020819052604081205482106110265760405162461bcd60e51b815260206004820152601360248201527f696e76616c696420636f6e7472616374206964000000000000000000000000006044820152606401610354565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260208190526040902080548390811061105d5761105d611a03565b9060005260206000200154905092915050565b60006108228261107f846114c0565b6000858152600760205260409020546006805490919081106110a3576110a3611a03565b906000526020600020906008020160070154611501565b600080600184815481106110d0576110d0611a03565b90600052602060002090600602019050806000015483111561110e5780546110f890846119c1565b9150611109843383600001546115d2565b61111f565b821561111f5761111f8433856115d2565b5092915050565b60008383116111775760405162461bcd60e51b815260206004820152601e60248201527f656e64206d7573742062652067726561746572207468616e20636c69666600006044820152606401610354565b50600180546040805160c08101825297885260006020808a018281528a8401998a5260608b0198895260808b0197885295151560a08b01908152848601865585835299517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6600686029081019190915595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf787015597517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf886015595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf985015593517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfa84015595517fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cfb90920180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169215159290921790915573ffffffffffffffffffffffffffffffffffffffff95909516825281835281208054948501815581522090910181905590565b73ffffffffffffffffffffffffffffffffffffffff82166000908152602081905260408120805490919061133c906001906119c1565b67ffffffffffffffff811115611354576113546119d4565b60405190808252806020026020018201604052801561137d578160200160208202803683370190505b5090506000805b835481101561140b57868482815481106113a0576113a0611a03565b906000526020600020015403156113f9578381815481106113c3576113c3611a03565b90600052602060002001548383806113da90611a32565b9450815181106113ec576113ec611a03565b6020026020010181815250505b8061140381611a32565b915050611384565b5073ffffffffffffffffffffffffffffffffffffffff851660009081526020818152604090912083516114409285019061171e565b5050505073ffffffffffffffffffffffffffffffffffffffff16600090815260208181526040822080546001810182559083529120019190915550565b60006114898585611a6a565b83106114965750806114b8565b836114a186856119c1565b6114ab9084611a7d565b6114b59190611a94565b90505b949350505050565b600080600183815481106114d6576114d6611a03565b90600052602060002090600602019050806001015481600001546114fa9190611a6a565b9392505050565b6000806001858154811061151757611517611a03565b906000526020600020906006020190508060030154816002015461153b9190611a6a565b83101561158a5760405162461bcd60e51b815260206004820152601c60248201527f636c69666620706572696f64206e6f74207265616368656420796574000000006044820152606401610354565b60006115a082600201548360040154868861147d565b905081600101548111156115c65760018201546115bd90826119c1565b925050506114fa565b50600095945050505050565b6000600184815481106115e7576115e7611a03565b906000526020600020906006020190508181600001600082825461160b91906119c1565b92505081905550818160010160008282546116269190611a6a565b90915550610367905083836002546040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff848116600483015260248201849052600092169063a9059cbb906044016020604051808303816000875af11580156116ab573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906116cf9190611acf565b90508061038d5760405162461bcd60e51b815260206004820152601360248201527f4e544e206e6f74207472616e73666572726564000000000000000000000000006044820152606401610354565b828054828255906000526020600020908101928215611759579160200282015b8281111561175957825182559160200191906001019061173e565b50611765929150611769565b5090565b5b80821115611765576000815560010161176a565b6000806040838503121561179157600080fd5b50508035926020909101359150565b6000602082840312156117b257600080fd5b5035919050565b600080600080608085870312156117cf57600080fd5b5050823594602084013594506040840135936060013592509050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461180f57600080fd5b919050565b60006020828403121561182657600080fd5b61081f826117eb565b6020808252825182820181905260009190848201906040850190845b818110156118a857611895838551805182526020810151602083015260408101516040830152606081015160608301526080810151608083015260a0810151151560a08301525050565b9284019260c0929092019160010161184b565b50909695505050505050565b600080604083850312156118c757600080fd5b6118d0836117eb565b946020939093013593505050565b6000806000606084860312156118f357600080fd5b6118fc846117eb565b95602085013595506040909401359392505050565b60008060006060848603121561192657600080fd5b61192f846117eb565b925060208401359150611944604085016117eb565b90509250925092565b60c081016108228284805182526020810151602083015260408101516040830152606081015160608301526080810151608083015260a0810151151560a08301525050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561082257610822611992565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611a6357611a63611992565b5060010190565b8082018082111561082257610822611992565b808202811582820484141761082257610822611992565b600082611aca577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600060208284031215611ae157600080fd5b815180151581146114fa57600080fdfea264697066735822122007d77e443f978150a808e49dfa2558a03367444453afe47c4fd949510189d04564736f6c63430008150033")

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
      "inputs" : [
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
                  "internalType" : "uint256",
                  "name" : "amount",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "unsubscribedAmount",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "totalUnlocked",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "totalUnlockedUnsubscribed",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "lastUnlockTime",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct NonStakableVesting.Schedule",
            "name" : "",
            "type" : "tuple"
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
