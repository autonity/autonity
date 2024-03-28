package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OracleBytecode = common.Hex2Bytes("6080604052600160ff1b600755600160ff1b6008553480156200002157600080fd5b50604051620034c6380380620034c6833981016040819052620000449162000639565b600280546001600160a01b038087166001600160a01b03199283161790925560038054928616929091169190911790558151620000899060009060208501906200035f565b5081516200009f9060019060208501906200035f565b5080600981905550620000c485600060018851620000be91906200074e565b62000181565b8451620000d9906004906020880190620003bc565b508451620000ef906005906020880190620003bc565b5060016006819055600d8054909101815560009081525b855181101562000175576001600b60008884815181106200012b576200012b6200076a565b6020908102919091018101516001600160a01b03168252810191909152604001600020600201805460ff1916911515919091179055806200016c8162000780565b91505062000106565b505050505050620009c3565b8082126200018e57505050565b81816000856002620001a185856200079c565b620001ad9190620007c6565b620001b9908762000806565b81518110620001cc57620001cc6200076a565b602002602001015190505b8183136200032b575b806001600160a01b0316868481518110620001ff57620001ff6200076a565b60200260200101516001600160a01b031610156200022c5782620002238162000831565b935050620001e0565b806001600160a01b03168683815181106200024b576200024b6200076a565b60200260200101516001600160a01b031611156200027857816200026f816200084c565b9250506200022c565b81831362000325578582815181106200029557620002956200076a565b6020026020010151868481518110620002b257620002b26200076a565b6020026020010151878581518110620002cf57620002cf6200076a565b60200260200101888581518110620002eb57620002eb6200076a565b6001600160a01b0393841660209182029290920101529116905282620003118162000831565b935050818062000321906200084c565b9250505b620001d7565b8185121562000341576200034186868462000181565b8383121562000357576200035786848662000181565b505050505050565b828054828255906000526020600020908101928215620003aa579160200282015b82811115620003aa5782518290620003999082620008f7565b509160200191906001019062000380565b50620003b892915062000422565b5090565b82805482825590600052602060002090810192821562000414579160200282015b828111156200041457825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620003dd565b50620003b892915062000443565b80821115620003b85760006200043982826200045a565b5060010162000422565b5b80821115620003b8576000815560010162000444565b50805462000468906200086c565b6000825580601f1062000479575050565b601f01602090049060005260206000209081019062000499919062000443565b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620004dd57620004dd6200049c565b604052919050565b60006001600160401b038211156200050157620005016200049c565b5060051b60200190565b80516001600160a01b03811681146200052357600080fd5b919050565b6000601f83818401126200053b57600080fd5b82516020620005546200054e83620004e5565b620004b2565b82815260059290921b850181019181810190878411156200057457600080fd5b8287015b848110156200062d5780516001600160401b03808211156200059a5760008081fd5b818a0191508a603f830112620005b05760008081fd5b8582015181811115620005c757620005c76200049c565b620005da818a01601f19168801620004b2565b915080825260408c81838601011115620005f45760008081fd5b60005b8281101562000614578481018201518482018a01528801620005f7565b5050600090820187015284525091830191830162000578565b50979650505050505050565b600080600080600060a086880312156200065257600080fd5b85516001600160401b03808211156200066a57600080fd5b818801915088601f8301126200067f57600080fd5b81516020620006926200054e83620004e5565b82815260059290921b8401810191818101908c841115620006b257600080fd5b948201945b83861015620006db57620006cb866200050b565b82529482019490820190620006b7565b9950620006ec90508a82016200050b565b97505050620006fe604089016200050b565b945060608801519150808211156200071557600080fd5b50620007248882890162000528565b925050608086015190509295509295909350565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000764576200076462000738565b92915050565b634e487b7160e01b600052603260045260246000fd5b60006001820162000795576200079562000738565b5060010190565b8181036000831280158383131683831282161715620007bf57620007bf62000738565b5092915050565b600082620007e457634e487b7160e01b600052601260045260246000fd5b600160ff1b82146000198414161562000801576200080162000738565b500590565b808201828112600083128015821682158216171562000829576200082962000738565b505092915050565b60006001600160ff1b01820162000795576200079562000738565b6000600160ff1b820162000864576200086462000738565b506000190190565b600181811c908216806200088157607f821691505b602082108103620008a257634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620008f257600081815260208120601f850160051c81016020861015620008d15750805b601f850160051c820191505b818110156200035757828155600101620008dd565b505050565b81516001600160401b038111156200091357620009136200049c565b6200092b816200092484546200086c565b84620008a8565b602080601f8311600181146200096357600084156200094a5750858301515b600019600386901b1c1916600185901b17855562000357565b600085815260208120601f198616915b82811015620009945788860151825594840194600190910190840162000973565b5085821015620009b35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b612af380620009d36000396000f3fe6080604052600436106101565760003560e01c80638d4f75d2116100bf578063b3ab15fb11610079578063cdd7225311610056578063cdd7225314610413578063df7f710e14610435578063e6a02a281461045757005b8063b3ab15fb146103be578063b78dec52146103de578063ccce413b146103f357005b80639f8743f7116100a75780639f8743f71461037d578063a781358714610392578063aa2f89b5146103a857005b80638d4f75d2146103475780639670c0bc1461036757005b80634bb278f3116101105780635281b5c6116100f85780635281b5c61461029e5780635412b3ae146102cb578063845023f21461032757005b80634bb278f3146102315780634c56ea561461025657005b8063307de9b61161013e578063307de9b61461019e57806333f98c77146101be5780633c8510fd1461021157005b806308f21ff51461015f578063146ca5311461018857005b3661015d57005b005b34801561016b57600080fd5b5061017560085481565b6040519081526020015b60405180910390f35b34801561019457600080fd5b5061017560065481565b3480156101aa57600080fd5b5061015d6101b9366004611f82565b61046d565b3480156101ca57600080fd5b506101de6101d9366004612113565b610711565b60405161017f91908151815260208083015190820152604080830151908201526060918201519181019190915260800190565b34801561021d57600080fd5b506101de61022c366004612148565b6108c8565b34801561023d57600080fd5b50610246610a66565b604051901515815260200161017f565b34801561026257600080fd5b506101756102713660046121b8565b8151602081840181018051600c825292820194820194909420919093529091526000908152604090205481565b3480156102aa57600080fd5b506102be6102b9366004612206565b610cb1565b60405161017f919061228d565b3480156102d757600080fd5b5061030a6102e63660046122a7565b600b6020526000908152604090208054600182015460029092015490919060ff1683565b60408051938452602084019290925215159082015260600161017f565b34801561033357600080fd5b5061015d6103423660046122e6565b610d5d565b34801561035357600080fd5b5061015d610362366004612383565b610ea5565b34801561037357600080fd5b5062989680610175565b34801561038957600080fd5b50600654610175565b34801561039e57600080fd5b5061017560095481565b3480156103b457600080fd5b5061017560075481565b3480156103ca57600080fd5b5061015d6103d93660046122a7565b61107e565b3480156103ea57600080fd5b50600954610175565b3480156103ff57600080fd5b506102be61040e366004612206565b61116c565b34801561041f57600080fd5b5061042861117c565b60405161017f9190612434565b34801561044157600080fd5b5061044a6111eb565b60405161017f91906124e3565b34801561046357600080fd5b50610175600a5481565b336000908152600b602052604090206002015460ff166104ee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f746572730000000000000060448201526064015b60405180910390fd5b600654336000908152600b602052604090205403610568576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f616c726561647920766f7465640000000000000000000000000000000000000060448201526064016104e5565b336000908152600b602052604081206001810180549087905581546006549092559181900361059857505061070b565b60005484146105a857505061070b565b60016006546105b79190612525565b811415806105f45750848484336040516020016105d79493929190612538565b6040516020818303038152906040528051906020012060001c8214155b156106895760005b600054811015610681577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600c6000838154811061063c5761063c612596565b906000526020600020016040516106539190612618565b9081526040805160209281900383019020336000908152925290205580610679816126ac565b9150506105fc565b50505061070b565b60005b84811015610707578585828181106106a6576106a6612596565b90506020020135600c600083815481106106c2576106c2612596565b906000526020600020016040516106d99190612618565b90815260408051602092819003830190203360009081529252902055806106ff816126ac565b91505061068c565b5050505b50505050565b61073c6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000600d600160065461074f9190612525565b8154811061075f5761075f612596565b906000526020600020018360405161077791906126c6565b908152602001604051809103902060405180606001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff1660018111156107c8576107c86126e2565b60018111156107d9576107d96126e2565b8152505090508060200151600003610873576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f64617461206f66207468652061736b65642073796d626f6c206973206e6f742060448201527f617661696c61626c65000000000000000000000000000000000000000000000060648201526084016104e5565b60006040518060800160405280600160065461088f9190612525565b81526020018360000151815260200183602001518152602001836040015160018111156108be576108be6126e2565b9052949350505050565b6108f36040518060800160405280600081526020016000815260200160008152602001600081525090565b6000600d848154811061090857610908612596565b906000526020600020018360405161092091906126c6565b908152602001604051809103902060405180606001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff166001811115610971576109716126e2565b6001811115610982576109826126e2565b8152505090508060200151600003610a1c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602960248201527f64617461206f66207468652061736b65642073796d626f6c206973206e6f742060448201527f617661696c61626c65000000000000000000000000000000000000000000000060648201526084016104e5565b60006040518060800160405280868152602001836000015181526020018360200151815260200183604001516001811115610a5957610a596126e2565b9052925050505b92915050565b60025460009073ffffffffffffffffffffffffffffffffffffffff163314610b10576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e5565b600954600a54610b209190612711565b4310610cab5760005b600054811015610b4e57610b3c816113a9565b610b47600182612711565b9050610b29565b5060065460075403610bf65760005b600554811015610bf4576001600b600060058481548110610b8057610b80612596565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905580610bec816126ac565b915050610b5d565b505b600654600754610c07906001612724565b03610c1457610c1461171e565b43600a81905550600160066000828254610c2e9190612711565b9091555050600854610c41906002612724565b60065403610c5b5760018054610c5991600091611daf565b505b60065460095460408051928352436020840152429083015260608201527fb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e59060800160405180910390a150600190565b50600090565b60018181548110610cc157600080fd5b906000526020600020016000915090508054610cdc906125c5565b80601f0160208091040260200160405190810160405280929190818152602001828054610d08906125c5565b8015610d555780601f10610d2a57610100808354040283529160200191610d55565b820191906000526020600020905b815481529060010190602001808311610d3857829003601f168201915b505050505081565b60025473ffffffffffffffffffffffffffffffffffffffff163314610e04576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e5565b8051600003610e6f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f566f746572732063616e277420626520656d707479000000000000000000000060448201526064016104e5565b610e8881600060018451610e839190612525565b611985565b8051610e9b906005906020840190611e07565b5050600654600755565b60035473ffffffffffffffffffffffffffffffffffffffff163314610f26576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064016104e5565b8051600003610f91576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f73796d626f6c732063616e277420626520656d7074790000000000000000000060448201526064016104e5565b600654600854610fa2906001612724565b14158015610fb4575060065460085414155b61101a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e64000060448201526064016104e5565b805161102d906001906020840190611e8d565b5060065460088190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d908290611065906001612711565b60405161107392919061274c565b60405180910390a150565b60025473ffffffffffffffffffffffffffffffffffffffff163314611125576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e5565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60008181548110610cc157600080fd5b606060058054806020026020016040519081016040528092919081815260200182805480156111e157602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff1681526001909101906020018083116111b6575b5050505050905090565b606060065460085460016111ff9190612724565b036112db576001805480602002602001604051908101604052809291908181526020016000905b828210156112d2578382906000526020600020018054611245906125c5565b80601f0160208091040260200160405190810160405280929190818152602001828054611271906125c5565b80156112be5780601f10611293576101008083540402835291602001916112be565b820191906000526020600020905b8154815290600101906020018083116112a157829003601f168201915b505050505081526020019060010190611226565b50505050905090565b6000805480602002602001604051908101604052809291908181526020016000905b828210156112d257838290600052602060002001805461131c906125c5565b80601f0160208091040260200160405190810160405280929190818152602001828054611348906125c5565b80156113955780601f1061136a57610100808354040283529160200191611395565b820191906000526020600020905b81548152906001019060200180831161137857829003601f168201915b5050505050815260200190600101906112fd565b60008082815481106113bd576113bd612596565b9060005260206000200180546113d2906125c5565b80601f01602080910402602001604051908101604052809291908181526020018280546113fe906125c5565b801561144b5780601f106114205761010080835404028352916020019161144b565b820191906000526020600020905b81548152906001019060200180831161142e57829003601f168201915b50505050509050600060048054905067ffffffffffffffff81111561147257611472612007565b60405190808252806020026020018201604052801561149b578160200160208202803683370190505b5090506000805b6004548110156115f2576000600482815481106114c1576114c1612596565b600091825260208083209091015460065473ffffffffffffffffffffffffffffffffffffffff909116808452600b90925260409092205490925014158061156857507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600c8660405161153491906126c6565b908152604080516020928190038301902073ffffffffffffffffffffffffffffffffffffffff851660009081529252902054145b1561157357506115e0565b600c8560405161158391906126c6565b908152604080516020928190038301902073ffffffffffffffffffffffffffffffffffffffff84166000908152925290205484846115c0816126ac565b9550815181106115d2576115d2612596565b602002602001018181525050505b806115ea816126ac565b9150506114a2565b506000600d60016006546116069190612525565b8154811061161657611616612596565b906000526020600020018460405161162e91906126c6565b90815260405190819003602001902054905060018215611659576116528484611b73565b9150600090505b600d8054600190810182556000919091526040805160608101825284815242602082015291908201908390811115611693576116936126e2565b815250600d600654815481106116ab576116ab612596565b90600052602060002001866040516116c391906126c6565b9081526020016040518091039020600082015181600001556020820151816001015560408201518160020160006101000a81548160ff0219169083600181111561170f5761170f6126e2565b02179055505050505050505050565b6000805b60045482108015611734575060055481105b156118d8576005818154811061174c5761174c612596565b6000918252602090912001546004805473ffffffffffffffffffffffffffffffffffffffff909216918490811061178557611785612596565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16036117cc57816117b6816126ac565b92505080806117c4906126ac565b915050611722565b600581815481106117df576117df612596565b6000918252602090912001546004805473ffffffffffffffffffffffffffffffffffffffff909216918490811061181857611818612596565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1610156118ce57600b60006004848154811061185757611857612596565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff1683528201929092526040018120818155600181019190915560020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055816118c6816126ac565b925050611722565b806117c4816126ac565b60045482101561197057600b6000600484815481106118f9576118f9612596565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff1683528201929092526040018120818155600181019190915560020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905581611968816126ac565b9250506118d8565b6005805461198091600491611ed3565b505050565b80821261199157505050565b818160008560026119a2858561276e565b6119ac91906127c4565b6119b69087612724565b815181106119c6576119c6612596565b602002602001015190505b818313611b45575b8073ffffffffffffffffffffffffffffffffffffffff16868481518110611a0257611a02612596565b602002602001015173ffffffffffffffffffffffffffffffffffffffff161015611a385782611a308161280e565b9350506119d9565b8073ffffffffffffffffffffffffffffffffffffffff16868381518110611a6157611a61612596565b602002602001015173ffffffffffffffffffffffffffffffffffffffff161115611a975781611a8f8161283f565b925050611a38565b818313611b4057858281518110611ab057611ab0612596565b6020026020010151868481518110611aca57611aca612596565b6020026020010151878581518110611ae457611ae4612596565b60200260200101888581518110611afd57611afd612596565b73ffffffffffffffffffffffffffffffffffffffff93841660209182029290920101529116905282611b2e8161280e565b9350508180611b3c9061283f565b9250505b6119d1565b81851215611b5857611b58868684611985565b83831215611b6b57611b6b868486611985565b505050505050565b600081600003611b8557506000610a60565b611b9b836000611b96600186612525565b611c37565b6000611ba8600284612878565b9050611bb560028461288c565b15611bd957838181518110611bcc57611bcc612596565b6020026020010151611c2f565b6002848281518110611bed57611bed612596565b602002602001015185600184611c039190612525565b81518110611c1357611c13612596565b6020026020010151611c259190612724565b611c2f91906127c4565b949350505050565b8181808203611c47575050505050565b6000856002611c56878761276e565b611c6091906127c4565b611c6a9087612724565b81518110611c7a57611c7a612596565b602002602001015190505b818313611d89575b80868481518110611ca057611ca0612596565b60200260200101511215611cc05782611cb88161280e565b935050611c8d565b858281518110611cd257611cd2612596565b6020026020010151811215611cf35781611ceb8161283f565b925050611cc0565b818313611d8457858281518110611d0c57611d0c612596565b6020026020010151868481518110611d2657611d26612596565b6020026020010151878581518110611d4057611d40612596565b60200260200101888581518110611d5957611d59612596565b60209081029190910101919091525281611d728161283f565b9250508280611d809061280e565b9350505b611c85565b81851215611d9c57611d9c868684611c37565b83831215611b6b57611b6b868486611c37565b828054828255906000526020600020908101928215611df75760005260206000209182015b82811115611df75781611de784826128e6565b5091600101919060010190611dd4565b50611e03929150611f13565b5090565b828054828255906000526020600020908101928215611e81579160200282015b82811115611e8157825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190611e27565b50611e03929150611f30565b828054828255906000526020600020908101928215611df7579160200282015b82811115611df75782518290611ec390826129e5565b5091602001919060010190611ead565b828054828255906000526020600020908101928215611e815760005260206000209182015b82811115611e81578254825591600101919060010190611ef8565b80821115611e03576000611f278282611f45565b50600101611f13565b5b80821115611e035760008155600101611f31565b508054611f51906125c5565b6000825580601f10611f61575050565b601f016020900490600052602060002090810190611f7f9190611f30565b50565b60008060008060608587031215611f9857600080fd5b84359350602085013567ffffffffffffffff80821115611fb757600080fd5b818701915087601f830112611fcb57600080fd5b813581811115611fda57600080fd5b8860208260051b8501011115611fef57600080fd5b95986020929092019750949560400135945092505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561207d5761207d612007565b604052919050565b600082601f83011261209657600080fd5b813567ffffffffffffffff8111156120b0576120b0612007565b6120e160207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601612036565b8181528460208386010111156120f657600080fd5b816020850160208301376000918101602001919091529392505050565b60006020828403121561212557600080fd5b813567ffffffffffffffff81111561213c57600080fd5b611c2f84828501612085565b6000806040838503121561215b57600080fd5b82359150602083013567ffffffffffffffff81111561217957600080fd5b61218585828601612085565b9150509250929050565b803573ffffffffffffffffffffffffffffffffffffffff811681146121b357600080fd5b919050565b600080604083850312156121cb57600080fd5b823567ffffffffffffffff8111156121e257600080fd5b6121ee85828601612085565b9250506121fd6020840161218f565b90509250929050565b60006020828403121561221857600080fd5b5035919050565b60005b8381101561223a578181015183820152602001612222565b50506000910152565b6000815180845261225b81602086016020860161221f565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006122a06020830184612243565b9392505050565b6000602082840312156122b957600080fd5b6122a08261218f565b600067ffffffffffffffff8211156122dc576122dc612007565b5060051b60200190565b600060208083850312156122f957600080fd5b823567ffffffffffffffff81111561231057600080fd5b8301601f8101851361232157600080fd5b803561233461232f826122c2565b612036565b81815260059190911b8201830190838101908783111561235357600080fd5b928401925b82841015612378576123698461218f565b82529284019290840190612358565b979650505050505050565b6000602080838503121561239657600080fd5b823567ffffffffffffffff808211156123ae57600080fd5b818501915085601f8301126123c257600080fd5b81356123d061232f826122c2565b81815260059190911b830184019084810190888311156123ef57600080fd5b8585015b838110156124275780358581111561240b5760008081fd5b6124198b89838a0101612085565b8452509186019186016123f3565b5098975050505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561248257835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101612450565b50909695505050505050565b600081518084526020808501808196508360051b8101915082860160005b858110156124d65782840389526124c4848351612243565b988501989350908401906001016124ac565b5091979650505050505050565b6020815260006122a0602083018461248e565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b81810381811115610a6057610a606124f6565b60008186825b8781101561255c57813583526020928301929091019060010161253e565b5050938452505060601b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016602082015260340192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c908216806125d957607f821691505b602082108103612612577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b6000808354612626816125c5565b6001828116801561263e5760018114612671576126a0565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841687528215158302870194506126a0565b8760005260208060002060005b858110156126975781548a82015290840190820161267e565b50505082870194505b50929695505050505050565b600060001982036126bf576126bf6124f6565b5060010190565b600082516126d881846020870161221f565b9190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b80820180821115610a6057610a606124f6565b8082018281126000831280158216821582161715612744576127446124f6565b505092915050565b60408152600061275f604083018561248e565b90508260208301529392505050565b818103600083128015838313168383128216171561278e5761278e6124f6565b5092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826127d3576127d3612795565b60001983147f800000000000000000000000000000000000000000000000000000000000000083141615612809576128096124f6565b500590565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036126bf576126bf6124f6565b60007f80000000000000000000000000000000000000000000000000000000000000008203612870576128706124f6565b506000190190565b60008261288757612887612795565b500490565b60008261289b5761289b612795565b500690565b601f82111561198057600081815260208120601f850160051c810160208610156128c75750805b601f850160051c820191505b81811015611b6b578281556001016128d3565b8181036128f1575050565b6128fb82546125c5565b67ffffffffffffffff81111561291357612913612007565b6129278161292184546125c5565b846128a0565b6000601f82116001811461295b57600083156129435750848201545b600019600385901b1c1916600184901b1784556129de565b6000858152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0841690600086815260209020845b838110156129b35782860154825560019586019590910190602001612993565b50858310156129d15781850154600019600388901b60f8161c191681555b50505060018360011b0184555b5050505050565b815167ffffffffffffffff8111156129ff576129ff612007565b612a0d8161292184546125c5565b602080601f831160018114612a425760008415612a2a5750858301515b600019600386901b1c1916600185901b178555611b6b565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b82811015612a8f57888601518255948401946001909101908401612a70565b5085821015612aad5787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea26469706673582212202525f1a3007308bc292376bd51bda13f4b596caab156565ca3196faaa8f3f62464736f6c63430008150033")

var OracleAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address[]",
            "name" : "_voters",
            "type" : "address[]"
         },
         {
            "internalType" : "address",
            "name" : "_autonity",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_operator",
            "type" : "address"
         },
         {
            "internalType" : "string[]",
            "name" : "_symbols",
            "type" : "string[]"
         },
         {
            "internalType" : "uint256",
            "name" : "_votePeriod",
            "type" : "uint256"
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
            "internalType" : "uint256",
            "name" : "_round",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_height",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_timestamp",
            "type" : "uint256"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_votePeriod",
            "type" : "uint256"
         }
      ],
      "name" : "NewRound",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : false,
            "internalType" : "string[]",
            "name" : "_symbols",
            "type" : "string[]"
         },
         {
            "indexed" : false,
            "internalType" : "uint256",
            "name" : "_round",
            "type" : "uint256"
         }
      ],
      "name" : "NewSymbols",
      "type" : "event"
   },
   {
      "anonymous" : false,
      "inputs" : [
         {
            "indexed" : true,
            "internalType" : "address",
            "name" : "_voter",
            "type" : "address"
         },
         {
            "indexed" : false,
            "internalType" : "int256[]",
            "name" : "_votes",
            "type" : "int256[]"
         }
      ],
      "name" : "Voted",
      "type" : "event"
   },
   {
      "stateMutability" : "payable",
      "type" : "fallback"
   },
   {
      "inputs" : [],
      "name" : "finalize",
      "outputs" : [
         {
            "internalType" : "bool",
            "name" : "",
            "type" : "bool"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getPrecision",
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
      "inputs" : [],
      "name" : "getRound",
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
            "name" : "_round",
            "type" : "uint256"
         },
         {
            "internalType" : "string",
            "name" : "_symbol",
            "type" : "string"
         }
      ],
      "name" : "getRoundData",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "uint256",
                  "name" : "round",
                  "type" : "uint256"
               },
               {
                  "internalType" : "int256",
                  "name" : "price",
                  "type" : "int256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "timestamp",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "status",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct IOracle.RoundData",
            "name" : "data",
            "type" : "tuple"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getSymbols",
      "outputs" : [
         {
            "internalType" : "string[]",
            "name" : "",
            "type" : "string[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "getVotePeriod",
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
      "name" : "getVoters",
      "outputs" : [
         {
            "internalType" : "address[]",
            "name" : "",
            "type" : "address[]"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "lastRoundBlock",
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
      "name" : "lastVoterUpdateRound",
      "outputs" : [
         {
            "internalType" : "int256",
            "name" : "",
            "type" : "int256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "string",
            "name" : "_symbol",
            "type" : "string"
         }
      ],
      "name" : "latestRoundData",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "uint256",
                  "name" : "round",
                  "type" : "uint256"
               },
               {
                  "internalType" : "int256",
                  "name" : "price",
                  "type" : "int256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "timestamp",
                  "type" : "uint256"
               },
               {
                  "internalType" : "uint256",
                  "name" : "status",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct IOracle.RoundData",
            "name" : "data",
            "type" : "tuple"
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
      "name" : "newSymbols",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         },
         {
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ],
      "name" : "reports",
      "outputs" : [
         {
            "internalType" : "int256",
            "name" : "",
            "type" : "int256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "round",
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
            "internalType" : "string[]",
            "name" : "_symbols",
            "type" : "string[]"
         }
      ],
      "name" : "setSymbols",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "address[]",
            "name" : "_newVoters",
            "type" : "address[]"
         }
      ],
      "name" : "setVoters",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "symbolUpdatedRound",
      "outputs" : [
         {
            "internalType" : "int256",
            "name" : "",
            "type" : "int256"
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
      "name" : "symbols",
      "outputs" : [
         {
            "internalType" : "string",
            "name" : "",
            "type" : "string"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_commit",
            "type" : "uint256"
         },
         {
            "internalType" : "int256[]",
            "name" : "_reports",
            "type" : "int256[]"
         },
         {
            "internalType" : "uint256",
            "name" : "_salt",
            "type" : "uint256"
         }
      ],
      "name" : "vote",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "votePeriod",
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
      "name" : "votingInfo",
      "outputs" : [
         {
            "internalType" : "uint256",
            "name" : "round",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "commit",
            "type" : "uint256"
         },
         {
            "internalType" : "bool",
            "name" : "isVoter",
            "type" : "bool"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "stateMutability" : "payable",
      "type" : "receive"
   }
]
`))
