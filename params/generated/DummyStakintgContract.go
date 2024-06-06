package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var DummyStakintgContractBytecode = common.Hex2Bytes("60806040523480156200001157600080fd5b50604051620017c7380380620017c7833981016040819052620000349162000227565b600080546001600160a01b0319166001600160a01b03831690811790915560408051630cef984560e41b8152905163cef98450916004808201926020929091908290030181865afa1580156200008e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620000b4919062000259565b6001556000546040805163386a827b60e01b815290516001600160a01b039092169163386a827b916004808201926020929091908290030181865afa15801562000102573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000128919062000259565b600255600054604080516371d1bc5960e01b815290516001600160a01b03909216916371d1bc59916004808201926020929091908290030181865afa15801562000176573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200019c919062000259565b60035560005460408051632def6e8b60e11b815290516001600160a01b0390921691635bdedd16916004808201926020929091908290030181865afa158015620001ea573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000210919062000259565b600455506000805460ff60a81b1916905562000273565b6000602082840312156200023a57600080fd5b81516001600160a01b03811681146200025257600080fd5b9392505050565b6000602082840312156200026c57600080fd5b5051919050565b61154480620002836000396000f3fe6080604052600436106101805760003560e01c806397beb714116100d6578063cef984501161007f578063f92dc41d11610059578063f92dc41d1461058d578063fbb271ad146105ad578063ff160a31146105cd57600080fd5b8063cef98450146104cb578063d18736ab146104e1578063da3bec821461050157600080fd5b8063a5d059ca116100b0578063a5d059ca14610478578063a65258d91461048b578063a8920241146104ab57600080fd5b806397beb7141461040b5780639dfd1b8e14610445578063a515366a1461046557600080fd5b806355463ceb116101385780636661568511610112578063666156851461036757806370f33ab2146103b357806371d1bc59146103f557600080fd5b806355463ceb146102ae5780635bdedd1614610300578063604e7a8d1461031657600080fd5b8063386a827b11610169578063386a827b1461022a5780633c54c2901461024e57806347de4d811461026e57600080fd5b8063161605e31461018557806323394eec14610187575b600080fd5b005b34801561019357600080fd5b506101e66101a23660046111aa565b60076020526000908152604090208054600182015460029092015473ffffffffffffffffffffffffffffffffffffffff909116919060ff8082169161010090041684565b6040805173ffffffffffffffffffffffffffffffffffffffff9095168552602085019390935290151591830191909152151560608201526080015b60405180910390f35b34801561023657600080fd5b5061024060025481565b604051908152602001610221565b34801561025a57600080fd5b506101856102693660046111d8565b6105ed565b34801561027a57600080fd5b506000546102a1907501000000000000000000000000000000000000000000900460ff1681565b604051610221919061123c565b3480156102ba57600080fd5b506000546102db9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610221565b34801561030c57600080fd5b5061024060045481565b34801561032257600080fd5b50610185600080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000179055565b34801561037357600080fd5b5061039e6103823660046111aa565b6009602052600090815260409020805460019091015460ff1682565b60408051928352901515602083015201610221565b3480156103bf57600080fd5b506000546103e59074010000000000000000000000000000000000000000900460ff1681565b6040519015158152602001610221565b34801561040157600080fd5b5061024060035481565b34801561041757600080fd5b50610185600080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff169055565b34801561045157600080fd5b506101856104603660046112a1565b6107be565b6101856104733660046112f8565b610a15565b6101856104863660046112f8565b610aba565b34801561049757600080fd5b506102406104a63660046111aa565b610b1b565b3480156104b757600080fd5b506101856104c6366004611322565b610b3c565b3480156104d757600080fd5b5061024060015481565b3480156104ed57600080fd5b506101856104fc366004611384565b610daf565b34801561050d57600080fd5b5061056161051c3660046111aa565b60086020526000908152604090205473ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900460ff1682565b6040805173ffffffffffffffffffffffffffffffffffffffff9093168352901515602083015201610221565b34801561059957600080fd5b506102db6105a8366004611467565b61105c565b3480156105b957600080fd5b506101856105c8366004611489565b6110a1565b3480156105d957600080fd5b506102406105e83660046111aa565b6110fb565b60005473ffffffffffffffffffffffffffffffffffffffff16331461067f5760405162461bcd60e51b815260206004820152602860248201527f66756e6374696f6e207265737472696374656420746f204175746f6e6974792060448201527f636f6e747261637400000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b6003545a11156106f75760405162461bcd60e51b815260206004820152602260248201527f6d6f726520676173207265636569766564207468616e206d617820616c6c6f7760448201527f65640000000000000000000000000000000000000000000000000000000000006064820152608401610676565b60005474010000000000000000000000000000000000000000900460ff16156107625760405162461bcd60e51b815260206004820152601960248201527f756e626f6e64696e6752656c65617365642072657665727473000000000000006044820152606401610676565b60408051808201825292835290151560208084019182526000948552600990529220905181559051600190910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff16331461084b5760405162461bcd60e51b815260206004820152602860248201527f66756e6374696f6e207265737472696374656420746f204175746f6e6974792060448201527f636f6e74726163740000000000000000000000000000000000000000000000006064820152608401610676565b6001545a11156108c35760405162461bcd60e51b815260206004820152602260248201527f6d6f726520676173207265636569766564207468616e206d617820616c6c6f7760448201527f65640000000000000000000000000000000000000000000000000000000000006064820152608401610676565b60005474010000000000000000000000000000000000000000900460ff161561092e5760405162461bcd60e51b815260206004820152601660248201527f626f6e64696e674170706c6965642072657665727473000000000000000000006044820152606401610676565b6040805160808101825273ffffffffffffffffffffffffffffffffffffffff9586168152602080820195865293151581830190815292151560608201908152600097885260079094529520945185547fffffffffffffffffffffffff00000000000000000000000000000000000000001694169390931784559051600184015590516002909201805491517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009092169215157fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169290921761010091151591909102179055565b6000546040517fa515366a00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152602482018490529091169063a515366a9034906044015b60206040518083038185885af1158015610a90573d6000803e3d6000fd5b50505050506040513d601f19601f82011682018060405250810190610ab591906114b1565b505050565b6000546040517fa5d059ca00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8481166004830152602482018490529091169063a5d059ca903490604401610a72565b60058181548110610b2b57600080fd5b600091825260209091200154905081565b60005473ffffffffffffffffffffffffffffffffffffffff163314610bc95760405162461bcd60e51b815260206004820152602860248201527f66756e6374696f6e207265737472696374656420746f204175746f6e6974792060448201527f636f6e74726163740000000000000000000000000000000000000000000000006064820152608401610676565b6002545a1115610c415760405162461bcd60e51b815260206004820152602260248201527f6d6f726520676173207265636569766564207468616e206d617820616c6c6f7760448201527f65640000000000000000000000000000000000000000000000000000000000006064820152608401610676565b60005474010000000000000000000000000000000000000000900460ff1615610cac5760405162461bcd60e51b815260206004820152601860248201527f756e626f6e64696e674170706c696564207265766572747300000000000000006044820152606401610676565b60408051808201825273ffffffffffffffffffffffffffffffffffffffff808516825283151560208084019182526000888152600890915293909320915182549351151574010000000000000000000000000000000000000000027fffffffffffffffffffffff00000000000000000000000000000000000000000090941691161791909117905560026000547501000000000000000000000000000000000000000000900460ff166002811115610d6657610d6661120d565b03610ab557600080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff1674010000000000000000000000000000000000000000179055505050565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e3c5760405162461bcd60e51b815260206004820152602860248201527f66756e6374696f6e207265737472696374656420746f204175746f6e6974792060448201527f636f6e74726163740000000000000000000000000000000000000000000000006064820152608401610676565b8051600454610e4b91906114ca565b5a1115610ec05760405162461bcd60e51b815260206004820152602260248201527f6d6f726520676173207265636569766564207468616e206d617820616c6c6f7760448201527f65640000000000000000000000000000000000000000000000000000000000006064820152608401610676565b60005474010000000000000000000000000000000000000000900460ff1615610f2b5760405162461bcd60e51b815260206004820152601a60248201527f72657761726473446973747269627574656420726576657274730000000000006044820152606401610676565b80600a60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663c9d97af46040518163ffffffff1660e01b8152600401602060405180830381865afa158015610f9c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fc091906114b1565b81526020019081526020016000209080519060200190610fe192919061110b565b5060016000547501000000000000000000000000000000000000000000900460ff1660028111156110145761101461120d565b0361105957600080547fffffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffff16740100000000000000000000000000000000000000001790555b50565b600a602052816000526040600020818154811061107857600080fd5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff169150829050565b600080548291907fffffffffffffffffffff00ffffffffffffffffffffffffffffffffffffffffff1675010000000000000000000000000000000000000000008360028111156110f3576110f361120d565b021790555050565b60068181548110610b2b57600080fd5b828054828255906000526020600020908101928215611185579160200282015b8281111561118557825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90911617825560209092019160019091019061112b565b50611191929150611195565b5090565b5b808211156111915760008155600101611196565b6000602082840312156111bc57600080fd5b5035919050565b803580151581146111d357600080fd5b919050565b6000806000606084860312156111ed57600080fd5b8335925060208401359150611204604085016111c3565b90509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b6020810160038310611277577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b91905290565b803573ffffffffffffffffffffffffffffffffffffffff811681146111d357600080fd5b600080600080600060a086880312156112b957600080fd5b853594506112c96020870161127d565b9350604086013592506112de606087016111c3565b91506112ec608087016111c3565b90509295509295909350565b6000806040838503121561130b57600080fd5b6113148361127d565b946020939093013593505050565b60008060006060848603121561133757600080fd5b833592506113476020850161127d565b9150611204604085016111c3565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000602080838503121561139757600080fd5b823567ffffffffffffffff808211156113af57600080fd5b818501915085601f8301126113c357600080fd5b8135818111156113d5576113d5611355565b8060051b6040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0603f8301168101818110858211171561141857611418611355565b60405291825284820192508381018501918883111561143657600080fd5b938501935b8285101561145b5761144c8561127d565b8452938501939285019261143b565b98975050505050505050565b6000806040838503121561147a57600080fd5b50508035926020909101359150565b60006020828403121561149b57600080fd5b8135600381106114aa57600080fd5b9392505050565b6000602082840312156114c357600080fd5b5051919050565b8082028115828204841417611508577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b9291505056fea2646970667358221220ffb7a33c22a27fe37b18a609470a9c7ecd909b4d90671b64fd5ebf4b00f546fd64736f6c63430008150033")

var DummyStakintgContractAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address payable",
            "name" : "_autonity",
            "type" : "address"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "constructor"
   },
   {
      "inputs" : [],
      "name" : "autonity",
      "outputs" : [
         {
            "internalType" : "contract Autonity",
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
            "name" : "_validator",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "bond",
      "outputs" : [],
      "stateMutability" : "payable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_bondingID",
            "type" : "uint256"
         },
         {
            "internalType" : "address",
            "name" : "_validator",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_liquid",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "_selfDelegation",
            "type" : "bool"
         },
         {
            "internalType" : "bool",
            "name" : "_rejected",
            "type" : "bool"
         }
      ],
      "name" : "bondingApplied",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "maxBondAppliedGas",
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
      "name" : "maxRewardsDistributionGas",
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
      "name" : "maxUnbondAppliedGas",
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
      "name" : "maxUnbondReleasedGas",
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
      "name" : "notifiedBondings",
      "outputs" : [
         {
            "internalType" : "address",
            "name" : "validator",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "liquid",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "selfDelegation",
            "type" : "bool"
         },
         {
            "internalType" : "bool",
            "name" : "rejected",
            "type" : "bool"
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
      "name" : "notifiedRelease",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "rejeced",
            "type" : "bool"
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
         },
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "name" : "notifiedRewardsDistribution",
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
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "name" : "notifiedUnbonding",
      "outputs" : [
         {
            "internalType" : "address",
            "name" : "validator",
            "type" : "address"
         },
         {
            "internalType" : "bool",
            "name" : "rejected",
            "type" : "bool"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "processStakingOperations",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "receiveATN",
      "outputs" : [],
      "stateMutability" : "payable",
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
      "name" : "requestedBondings",
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
      "name" : "requestedUnbondings",
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
      "name" : "revertClock",
      "outputs" : [
         {
            "internalType" : "enum DummyStakintgContract.Fail",
            "name" : "",
            "type" : "uint8"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "revertStaking",
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
      "inputs" : [],
      "name" : "revertStakingOperations",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address[]",
            "name" : "_validators",
            "type" : "address[]"
         }
      ],
      "name" : "rewardsDistributed",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "enum DummyStakintgContract.Fail",
            "name" : "_clock",
            "type" : "uint8"
         }
      ],
      "name" : "setRevertClock",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_validator",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         }
      ],
      "name" : "unbond",
      "outputs" : [],
      "stateMutability" : "payable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_unbondingID",
            "type" : "uint256"
         },
         {
            "internalType" : "address",
            "name" : "_validator",
            "type" : "address"
         },
         {
            "internalType" : "bool",
            "name" : "_rejected",
            "type" : "bool"
         }
      ],
      "name" : "unbondingApplied",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_unbondingID",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_amount",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "_rejected",
            "type" : "bool"
         }
      ],
      "name" : "unbondingReleased",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   }
]
`))
