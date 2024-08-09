package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OmissionAccountabilityBytecode = common.Hex2Bytes("60806040523480156200001157600080fd5b5060405162002b1a38038062002b1a83398101604081905262000034916200025e565b601480546001600160a01b038088166001600160a01b03199283161790925560048054918216919092161790558051600d55602080820151600e556040820151600f556060820151601055608082015160115560a082015160125560c08201516013558351620000ab9160009190860190620000d4565b508151620000c1906001906020850190620000d4565b5050600e54600355506200036092505050565b8280548282559060005260206000209081019282156200012c579160200282015b828111156200012c57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620000f5565b506200013a9291506200013e565b5090565b5b808211156200013a57600081556001016200013f565b6001600160a01b03811681146200016b57600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60405160e081016001600160401b0381118282101715620001a957620001a96200016e565b60405290565b600082601f830112620001c157600080fd5b815160206001600160401b0380831115620001e057620001e06200016e565b8260051b604051601f19603f830116810181811084821117156200020857620002086200016e565b6040529384528581018301938381019250878511156200022757600080fd5b83870191505b8482101562000253578151620002438162000155565b835291830191908301906200022d565b979650505050505050565b60008060008060008587036101608112156200027957600080fd5b8651620002868162000155565b6020880151909650620002998162000155565b60408801519095506001600160401b0380821115620002b757600080fd5b620002c58a838b01620001af565b95506060890151915080821115620002dc57600080fd5b50620002eb89828a01620001af565b93505060e0607f19820112156200030157600080fd5b506200030c62000184565b6080870151815260a0870151602082015260c0870151604082015260e08701516060820152610100870151608082015261012087015160a082015261014087015160c0820152809150509295509295909350565b6127aa80620003706000396000f3fe60806040526004361061018b5760003560e01c8063b3ab15fb116100d6578063d7eaef491161007f578063eeb9223311610059578063eeb92233146104f6578063f85cffe214610509578063f95bbd7f1461052957600080fd5b8063d7eaef4914610489578063e1e8cac6146104a9578063eb231a1a146104c957600080fd5b8063ce4b5bbe116100b0578063ce4b5bbe14610426578063d2aaca571461043c578063d5baf9081461046957600080fd5b8063b3ab15fb1461039b578063b8d5712a146103bb578063c1a482451461040657600080fd5b806370432e8b116101385780637f5e2f11116101125780637f5e2f11146103305780638bbde7e5146103455780639a11e0e61461036557600080fd5b806370432e8b1461027f57806379502c55146102ac5780637e7168231461031057600080fd5b806348fa71271161016957806348fa7127146102085780635426b5ea146102285780635ca1809c1461025557600080fd5b80631ede5a1a14610190578063278112dc146101b9578063482893c7146101e6575b600080fd5b34801561019c57600080fd5b506101a660085481565b6040519081526020015b60405180910390f35b3480156101c557600080fd5b506101a66101d4366004611ed7565b600a6020526000908152604090205481565b3480156101f257600080fd5b50610206610201366004611efb565b610559565b005b34801561021457600080fd5b50610206610223366004611efb565b6105bd565b34801561023457600080fd5b506101a6610243366004611ed7565b60076020526000908152604090205481565b34801561026157600080fd5b50600e546003546040805183815291909214156020820152016101b0565b34801561028b57600080fd5b506101a661029a366004611ed7565b600b6020526000908152604090205481565b3480156102b857600080fd5b50600d54600e54600f546010546011546012546013546102db9695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e0016101b0565b34801561031c57600080fd5b5061020661032b366004611efb565b61066e565b34801561033c57600080fd5b506127106101a6565b34801561035157600080fd5b50610206610360366004611efb565b6106cd565b34801561037157600080fd5b506101a6610380366004611ed7565b6001600160a01b03166000908152600a602052604090205490565b3480156103a757600080fd5b506102066103b6366004611ed7565b610841565b3480156103c757600080fd5b506103f66103d6366004611f14565b600660209081526000928352604080842090915290825290205460ff1681565b60405190151581526020016101b0565b34801561041257600080fd5b50610206610421366004612083565b6108fb565b34801561043257600080fd5b506101a661271081565b34801561044857600080fd5b506101a6610457366004611ed7565b60096020526000908152604090205481565b34801561047557600080fd5b50610206610484366004612104565b610afe565b34801561049557600080fd5b506102066104a4366004611efb565b610baa565b3480156104b557600080fd5b506102066104c4366004611efb565b610c5b565b3480156104d557600080fd5b506101a66104e4366004611ed7565b600c6020526000908152604090205481565b610206610504366004611efb565b610ce0565b34801561051557600080fd5b50610206610524366004611efb565b610fed565b34801561053557600080fd5b506103f6610544366004611efb565b60056020526000908152604090205460ff1681565b6004546001600160a01b031633146105b85760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064015b60405180910390fd5b601055565b6004546001600160a01b031633146106175760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064016105af565b6127108111156106695760405162461bcd60e51b815260206004820152601a60248201527f63616e6e6f7420657863656564207363616c6520666163746f7200000000000060448201526064016105af565b601255565b6004546001600160a01b031633146106c85760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064016105af565b601155565b6004546001600160a01b031633146107275760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064016105af565b601454604080517fdfb1a4d200000000000000000000000000000000000000000000000000000000815290516000926001600160a01b03169163dfb1a4d29160048083019260209291908290030181865afa15801561078a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107ae9190612168565b905060016107bd83600a6121b0565b6107c791906121c9565b811161083b5760405162461bcd60e51b815260206004820152603c60248201527f65706f636820706572696f64206e6565647320746f206265206772656174657260448201527f207468616e2044454c54412b6c6f6f6b6261636b57696e646f772d310000000060648201526084016105af565b50600355565b6014546001600160a01b031633146108c15760405162461bcd60e51b815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e7472616374000000000000000000000000000000000000000060648201526084016105af565b600480547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b6014546001600160a01b0316331461097b5760405162461bcd60e51b815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e7472616374000000000000000000000000000000000000000060648201526084016105af565b6000610988600a436121c9565b905082156109f357600081815260056020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556001600160a01b0388168352600790915281208054916109e9836121dc565b9190505550610a76565b600081815260056020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556001600160a01b0388168352600990915281208054869290610a4d9084906121b0565b925050819055508360086000828254610a6691906121b0565b90915550610a7690508682611118565b8115610af6576000610a8661131e565b9050610a9181611523565b60005b600054811015610aed57600060076000808481548110610ab657610ab6612214565b60009182526020808320909101546001600160a01b0316835282019290925260400190205580610ae5816121dc565b915050610a94565b5050600354600e555b505050505050565b6014546001600160a01b03163314610b7e5760405162461bcd60e51b815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e7472616374000000000000000000000000000000000000000060648201526084016105af565b8151610b91906000906020850190611e2d565b508051610ba5906001906020840190611e2d565b505050565b6004546001600160a01b03163314610c045760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064016105af565b612710811115610c565760405162461bcd60e51b815260206004820152601a60248201527f63616e6e6f7420657863656564207363616c6520666163746f7200000000000060448201526064016105af565b600d55565b6014546001600160a01b03163314610cdb5760405162461bcd60e51b815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e7472616374000000000000000000000000000000000000000060648201526084016105af565b600255565b6014546001600160a01b03163314610d605760405162461bcd60e51b815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e7472616374000000000000000000000000000000000000000060648201526084016105af565b3460005b600054811015610fe357600060096000808481548110610d8657610d86612214565b60009182526020808320909101546001600160a01b031683528201929092526040019020541115610fd15760006008548360096000808681548110610dcd57610dcd612214565b60009182526020808320909101546001600160a01b03168352820192909252604001902054610dfc9190612243565b610e06919061225a565b905060006008548560096000808781548110610e2457610e24612214565b60009182526020808320909101546001600160a01b03168352820192909252604001902054610e539190612243565b610e5d919061225a565b905060018381548110610e7257610e72612214565b60009182526020822001546040516001600160a01b03909116916108fc918591818181858888f193505050503d8060008114610eca576040519150601f19603f3d011682016040523d82523d6000602084013e610ecf565b606091505b5050601454600180546001600160a01b03909216925063a9059cbb9186908110610efb57610efb612214565b60009182526020909120015460405160e083901b7fffffffff000000000000000000000000000000000000000000000000000000001681526001600160a01b039091166004820152602481018490526044016020604051808303816000875af1158015610f6c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610f909190612295565b50600060096000808681548110610fa957610fa9612214565b60009182526020808320909101546001600160a01b0316835282019290925260400190205550505b80610fdb816121dc565b915050610d64565b5050600060085550565b6004546001600160a01b031633146110475760405162461bcd60e51b815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064016105af565b6127108111156110995760405162461bcd60e51b815260206004820152601a60248201527f63616e6e6f7420657863656564207363616c6520666163746f7200000000000060448201526064016105af565b600d54811115611113576040805162461bcd60e51b81526020600482015260248101919091527f70617374506572666f726d616e63655765696768742063616e6e6f742062652060448201527f67726561746572207468616e20696e61637469766974795468726573686f6c6460648201526084016105af565b600f55565b60005b82518110156111ac57600082815260066020526040812084516001929086908590811061114a5761114a612214565b6020908102919091018101516001600160a01b0316825281019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055806111a4816121dc565b91505061111b565b50600e546002546111bd91906121b0565b8110156111c8575050565b60005b8251811015610ba557600e5460019060006111e683866121c9565b90505b6111f382866121c9565b8111156112ad5760008181526005602052604090205460ff161561124057816002548661122091906121c9565b1161122e57600092506112ad565b81611238816121dc565b92505061129b565b60066000828152602001908152602001600020600087868151811061126757611267612214565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff1661129b57600092506112ad565b806112a5816122b2565b9150506111e9565b50811561130957600760008685815181106112ca576112ca612214565b60200260200101516001600160a01b03166001600160a01b031681526020019081526020016000206000815480929190611303906121dc565b91905055505b50508080611316906121dc565b9150506111cb565b600080601460009054906101000a90046001600160a01b03166001600160a01b031663dfb1a4d26040518163ffffffff1660e01b8152600401602060405180830381865afa158015611374573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113989190612168565b90506000805b60005481101561151c57600e54600090600a906113bb90866121c9565b6113c69060016121b0565b6113d091906121c9565b612710600760008086815481106113e9576113e9612214565b60009182526020808320909101546001600160a01b031683528201929092526040019020546114189190612243565b611422919061225a565b905061271081111561143357506127105b6000612710600d60020154600a600080878154811061145457611454612214565b60009182526020808320909101546001600160a01b031683528201929092526040019020546114839190612243565b600f54611492906127106121c9565b61149c9085612243565b6114a691906121b0565b6114b0919061225a565b600d549091508111156114cb57836114c7816121dc565b9450505b80600a60008086815481106114e2576114e2612214565b60009182526020808320909101546001600160a01b0316835282019290925260400190205550819050611514816121dc565b91505061139e565b5092915050565b60005b600054811015611a12576014546000805490916001600160a01b031690631904bb2e9083908590811061155b5761155b612214565b60009182526020909120015460405160e083901b7fffffffff000000000000000000000000000000000000000000000000000000001681526001600160a01b039091166004820152602401600060405180830381865afa1580156115c3573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261160991908101906123b4565b90506002816102600151600381111561162457611624612532565b148061164657506003816102600151600381111561164457611644612532565b145b156116515750611a00565b600d60000154600a600080858154811061166d5761166d612214565b60009182526020808320909101546001600160a01b03168352820192909252604001902054116117a9576000600b60008085815481106116af576116af612214565b60009182526020808320909101546001600160a01b0316835282019290925260400190205411156117a457600b60008084815481106116f0576116f0612214565b60009182526020808320909101546001600160a01b031683528201929092526040018120805491611720836122b2565b9190505550600b600080848154811061173b5761173b612214565b60009182526020808320909101546001600160a01b0316835282019290925260400181205490036117a4576000600c600080858154811061177e5761177e612214565b60009182526020808320909101546001600160a01b031683528201929092526040019020555b6119fe565b600c60008084815481106117bf576117bf612214565b60009182526020808320909101546001600160a01b0316835282019290925260400181208054916117ef836121dc565b91905055506000600c600080858154811061180c5761180c612214565b60009182526020808320909101546001600160a01b0316835282019290925260400181205481549091600c9181908790811061184a5761184a612214565b60009182526020808320909101546001600160a01b031683528201929092526040019020546118799190612243565b9050600081600d6003015461188e9190612243565b9050600082600d600401546118a39190612243565b90506118af82436121b0565b6102008501526002610260850181815250506000600b60008088815481106118d9576118d9612214565b60009182526020808320909101546001600160a01b03168352820192909252604001902054111561192e57611929848785600d6005015461191a9190612243565b6119249190612243565b611a16565b6119aa565b6014546040517f35be16e00000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906335be16e0906119779087906004016125e6565b600060405180830381600087803b15801561199157600080fd5b505af11580156119a5573d6000803e3d6000fd5b505050505b80600b60008088815481106119c1576119c1612214565b60009182526020808320909101546001600160a01b03168352820192909252604001812080549091906119f59084906121b0565b90915550505050505b505b80611a0a816121dc565b915050611526565b5050565b601354811115611a2557506013545b60008261012001518360c001518460a00151611a4191906121b0565b611a4b91906121b0565b601354909150600090611a5e8385612243565b611a68919061225a565b9050600081118015611a7957508181145b15611b9b57600060a085018190526101008501819052610120850181905260c08501526101e084018051829190611ab19083906121b0565b905250600361026085015260006102008501526014546040517f35be16e00000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906335be16e090611b0d9087906004016125e6565b600060405180830381600087803b158015611b2757600080fd5b505af1158015611b3b573d6000803e3d6000fd5b50505050602084810151604080516001600160a01b0390921682529181018390526000818301526001606082015290517f3cac37f432247a020a7112d5052bc279f35e1e3b80b0aab0eca49d1773ed3e3f9181900360800190a150505050565b61012084015181908111611bc857808561012001818151611bbc91906121c9565b90525060009050611be3565b610120850151611bd890826121c9565b600061012087015290505b8015611c60578085610100015110611c2b57808561010001818151611c0891906121c9565b90525060a085018051829190611c1f9083906121c9565b90525060009050611c60565b610100850151611c3b90826121c9565b90508461010001518560a001818151611c5491906121c9565b90525060006101008601525b600081118015611c83575060008560a001518660c00151611c8191906121b0565b115b15611d2f5760008560a001518660c00151611c9e91906121b0565b60c0870151611cad9084612243565b611cb7919061225a565b905060008660a001518760c00151611ccf91906121b0565b60a0880151611cde9085612243565b611ce8919061225a565b9050818760c001818151611cfc91906121c9565b90525060a087018051829190611d139083906121c9565b905250611d2081836121b0565b611d2a90846121c9565b925050505b611d3981836121c9565b915081856101e001818151611d4e91906121b0565b9052506014546040517f35be16e00000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906335be16e090611d9a9088906004016125e6565b600060405180830381600087803b158015611db457600080fd5b505af1158015611dc8573d6000803e3d6000fd5b50505050602085810151610200870151604080516001600160a01b039093168352928201859052818301526000606082015290517f3cac37f432247a020a7112d5052bc279f35e1e3b80b0aab0eca49d1773ed3e3f9181900360800190a15050505050565b828054828255906000526020600020908101928215611e9a579160200282015b82811115611e9a57825182547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03909116178255602090920191600190910190611e4d565b50611ea6929150611eaa565b5090565b5b80821115611ea65760008155600101611eab565b6001600160a01b0381168114611ed457600080fd5b50565b600060208284031215611ee957600080fd5b8135611ef481611ebf565b9392505050565b600060208284031215611f0d57600080fd5b5035919050565b60008060408385031215611f2757600080fd5b823591506020830135611f3981611ebf565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610280810167ffffffffffffffff81118282101715611f9757611f97611f44565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611fe457611fe4611f44565b604052919050565b600082601f830112611ffd57600080fd5b8135602067ffffffffffffffff82111561201957612019611f44565b8160051b612028828201611f9d565b928352848101820192828101908785111561204257600080fd5b83870192505b8483101561206a57823561205b81611ebf565b82529183019190830190612048565b979650505050505050565b8015158114611ed457600080fd5b600080600080600060a0868803121561209b57600080fd5b853567ffffffffffffffff8111156120b257600080fd5b6120be88828901611fec565b95505060208601356120cf81611ebf565b93506040860135925060608601356120e681612075565b915060808601356120f681612075565b809150509295509295909350565b6000806040838503121561211757600080fd5b823567ffffffffffffffff8082111561212f57600080fd5b61213b86838701611fec565b9350602085013591508082111561215157600080fd5b5061215e85828601611fec565b9150509250929050565b60006020828403121561217a57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156121c3576121c3612181565b92915050565b818103818111156121c3576121c3612181565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361220d5761220d612181565b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80820281158282048414176121c3576121c3612181565b600082612290577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b6000602082840312156122a757600080fd5b8151611ef481612075565b6000816122c1576122c1612181565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b80516122f281611ebf565b919050565b60005b838110156123125781810151838201526020016122fa565b50506000910152565b600082601f83011261232c57600080fd5b815167ffffffffffffffff81111561234657612346611f44565b61237760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611f9d565b81815284602083860101111561238c57600080fd5b61239d8260208301602087016122f7565b949350505050565b8051600481106122f257600080fd5b6000602082840312156123c657600080fd5b815167ffffffffffffffff808211156123de57600080fd5b9083019061028082860312156123f357600080fd5b6123fb611f73565b612404836122e7565b8152612412602084016122e7565b6020820152612423604084016122e7565b604082015260608301518281111561243a57600080fd5b6124468782860161231b565b6060830152506080830151608082015260a083015160a082015260c083015160c082015260e083015160e08201526101008084015181830152506101208084015181830152506101408084015181830152506101608084015181830152506101806124b28185016122e7565b908201526101a083810151908201526101c080840151908201526101e0808401519082015261020080840151908201526102208084015190820152610240808401518381111561250157600080fd5b61250d8882870161231b565b82840152505061026091506125238284016123a5565b91810191909152949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600081518084526125798160208601602086016122f7565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600481106125e2577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b602081526126006020820183516001600160a01b03169052565b6000602083015161261c60408401826001600160a01b03169052565b5060408301516001600160a01b038116606084015250606083015161028080608085015261264e6102a0850183612561565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a06126d5818701836001600160a01b03169052565b8601516101c0868101919091528601516101e0808701919091528601516102008087019190915286015161022080870191909152860151610240808701919091528601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001610260808801919091529091506127558483612561565b93508087015191505061276a828601826125ab565b509094935050505056fea26469706673582212207aa65b1adc63bf0e3451a3ec09730ca01744ca7b7b1efc7062a3485899e36c3764736f6c63430008150033")

var OmissionAccountabilityAbi, _ = abi.JSON(strings.NewReader(`[
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
         },
         {
            "internalType" : "address[]",
            "name" : "_nodeAddresses",
            "type" : "address[]"
         },
         {
            "internalType" : "address[]",
            "name" : "_treasuries",
            "type" : "address[]"
         },
         {
            "components" : [
               {
                  "internalType" : "uint256",
                  "name" : "inactivityThreshold",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "lookbackWindow",
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
               },
               {
                  "internalType" : "uint256",
                  "name" : "slashingRatePrecision",
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
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "address",
            "name" : "validator",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "amount",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "releaseBlock",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "bool",
            "name" : "isJailbound",
            "type" : "bool"
         }
      ],
      "name" : "InactivitySlashingEvent",
      "type" : "event"
   },
   {
      "inputs" : [],
      "name" : "SCALE_FACTOR",
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
            "name" : "inactivityThreshold",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "lookbackWindow",
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
         },
         {
            "internalType" : "uint256",
            "name" : "slashingRatePrecision",
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
            "name" : "_ntnReward",
            "type" : "uint256"
         }
      ],
      "name" : "distributeProposerRewards",
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
            "name" : "_absentees",
            "type" : "address[]"
         },
         {
            "internalType" : "address",
            "name" : "_proposer",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_proposerEffort",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "_isProposerOmissionFaulty",
            "type" : "bool"
         },
         {
            "internalType" : "bool",
            "name" : "_epochEnded",
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
            "internalType" : "address",
            "name" : "_validator",
            "type" : "address"
         }
      ],
      "name" : "getInactivityScore",
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
      "name" : "getLookbackWindow",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         },
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
      "name" : "getScaleFactor",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "",
            "type" : "uint256"
         }
      ],
      "stateMutability" : "pure",
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
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ],
      "name" : "inactiveValidators",
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
      "name" : "inactivityScores",
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
      "name" : "probationPeriods",
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
      "name" : "proposerEffort",
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
            "name" : "_nodeAddresses",
            "type" : "address[]"
         },
         {
            "internalType" : "address[]",
            "name" : "_treasuries",
            "type" : "address[]"
         }
      ],
      "name" : "setCommittee",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_inactivityThreshold",
            "type" : "uint256"
         }
      ],
      "name" : "setInactivityThreshold",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_initialJailingPeriod",
            "type" : "uint256"
         }
      ],
      "name" : "setInitialJailingPeriod",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_initialProbationPeriod",
            "type" : "uint256"
         }
      ],
      "name" : "setInitialProbationPeriod",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_initialSlashingRate",
            "type" : "uint256"
         }
      ],
      "name" : "setInitialSlashingRate",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_lastEpochBlock",
            "type" : "uint256"
         }
      ],
      "name" : "setLastEpochBlock",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_lookbackWindow",
            "type" : "uint256"
         }
      ],
      "name" : "setLookbackWindow",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_operator",
            "type" : "address"
         }
      ],
      "name" : "setOperator",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_pastPerformanceWeight",
            "type" : "uint256"
         }
      ],
      "name" : "setPastPerformanceWeight",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "totalEffort",
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
