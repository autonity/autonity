package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OracleBytecode = common.Hex2Bytes("6080604052600160ff1b600755600160ff1b6008553480156200002157600080fd5b506040516200339e3803806200339e833981016040819052620000449162000639565b600280546001600160a01b038087166001600160a01b03199283161790925560038054928616929091169190911790558151620000899060009060208501906200035f565b5081516200009f9060019060208501906200035f565b5080600981905550620000c485600060018851620000be91906200074e565b62000181565b8451620000d9906004906020880190620003bc565b508451620000ef906005906020880190620003bc565b5060016006819055600d8054909101815560009081525b855181101562000175576001600b60008884815181106200012b576200012b6200076a565b6020908102919091018101516001600160a01b03168252810191909152604001600020600201805460ff1916911515919091179055806200016c8162000780565b91505062000106565b505050505050620009c3565b8082126200018e57505050565b81816000856002620001a185856200079c565b620001ad9190620007c6565b620001b9908762000806565b81518110620001cc57620001cc6200076a565b602002602001015190505b8183136200032b575b806001600160a01b0316868481518110620001ff57620001ff6200076a565b60200260200101516001600160a01b031610156200022c5782620002238162000831565b935050620001e0565b806001600160a01b03168683815181106200024b576200024b6200076a565b60200260200101516001600160a01b031611156200027857816200026f816200084c565b9250506200022c565b81831362000325578582815181106200029557620002956200076a565b6020026020010151868481518110620002b257620002b26200076a565b6020026020010151878581518110620002cf57620002cf6200076a565b60200260200101888581518110620002eb57620002eb6200076a565b6001600160a01b0393841660209182029290920101529116905282620003118162000831565b935050818062000321906200084c565b9250505b620001d7565b8185121562000341576200034186868462000181565b8383121562000357576200035786848662000181565b505050505050565b828054828255906000526020600020908101928215620003aa579160200282015b82811115620003aa5782518290620003999082620008f7565b509160200191906001019062000380565b50620003b892915062000422565b5090565b82805482825590600052602060002090810192821562000414579160200282015b828111156200041457825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620003dd565b50620003b892915062000443565b80821115620003b85760006200043982826200045a565b5060010162000422565b5b80821115620003b8576000815560010162000444565b50805462000468906200086c565b6000825580601f1062000479575050565b601f01602090049060005260206000209081019062000499919062000443565b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620004dd57620004dd6200049c565b604052919050565b60006001600160401b038211156200050157620005016200049c565b5060051b60200190565b80516001600160a01b03811681146200052357600080fd5b919050565b6000601f83818401126200053b57600080fd5b82516020620005546200054e83620004e5565b620004b2565b82815260059290921b850181019181810190878411156200057457600080fd5b8287015b848110156200062d5780516001600160401b03808211156200059a5760008081fd5b818a0191508a603f830112620005b05760008081fd5b8582015181811115620005c757620005c76200049c565b620005da818a01601f19168801620004b2565b915080825260408c81838601011115620005f45760008081fd5b60005b8281101562000614578481018201518482018a01528801620005f7565b5050600090820187015284525091830191830162000578565b50979650505050505050565b600080600080600060a086880312156200065257600080fd5b85516001600160401b03808211156200066a57600080fd5b818801915088601f8301126200067f57600080fd5b81516020620006926200054e83620004e5565b82815260059290921b8401810191818101908c841115620006b257600080fd5b948201945b83861015620006db57620006cb866200050b565b82529482019490820190620006b7565b9950620006ec90508a82016200050b565b97505050620006fe604089016200050b565b945060608801519150808211156200071557600080fd5b50620007248882890162000528565b925050608086015190509295509295909350565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000764576200076462000738565b92915050565b634e487b7160e01b600052603260045260246000fd5b60006001820162000795576200079562000738565b5060010190565b8181036000831280158383131683831282161715620007bf57620007bf62000738565b5092915050565b600082620007e457634e487b7160e01b600052601260045260246000fd5b600160ff1b82146000198414161562000801576200080162000738565b500590565b808201828112600083128015821682158216171562000829576200082962000738565b505092915050565b60006001600160ff1b01820162000795576200079562000738565b6000600160ff1b820162000864576200086462000738565b506000190190565b600181811c908216806200088157607f821691505b602082108103620008a257634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620008f257600081815260208120601f850160051c81016020861015620008d15750805b601f850160051c820191505b818110156200035757828155600101620008dd565b505050565b81516001600160401b038111156200091357620009136200049c565b6200092b816200092484546200086c565b84620008a8565b602080601f8311600181146200096357600084156200094a5750858301515b600019600386901b1c1916600185901b17855562000357565b600085815260208120601f198616915b82811015620009945788860151825594840194600190910190840162000973565b5085821015620009b35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6129cb80620009d36000396000f3fe6080604052600436106101565760003560e01c80638d4f75d2116100bf578063b3ab15fb11610079578063cdd7225311610056578063cdd7225314610413578063df7f710e14610435578063e6a02a281461045757005b8063b3ab15fb146103be578063b78dec52146103de578063ccce413b146103f357005b80639f8743f7116100a75780639f8743f71461037d578063a781358714610392578063aa2f89b5146103a857005b80638d4f75d2146103475780639670c0bc1461036757005b80634bb278f3116101105780635281b5c6116100f85780635281b5c61461029e5780635412b3ae146102cb578063845023f21461032757005b80634bb278f3146102315780634c56ea561461025657005b8063307de9b61161013e578063307de9b61461019e57806333f98c77146101be5780633c8510fd1461021157005b806308f21ff51461015f578063146ca5311461018857005b3661015d57005b005b34801561016b57600080fd5b5061017560085481565b6040519081526020015b60405180910390f35b34801561019457600080fd5b5061017560065481565b3480156101aa57600080fd5b5061015d6101b9366004611e5a565b61046d565b3480156101ca57600080fd5b506101de6101d9366004611feb565b610711565b60405161017f91908151815260208083015190820152604080830151908201526060918201519181019190915260800190565b34801561021d57600080fd5b506101de61022c366004612020565b610834565b34801561023d57600080fd5b5061024661093e565b604051901515815260200161017f565b34801561026257600080fd5b50610175610271366004612090565b8151602081840181018051600c825292820194820194909420919093529091526000908152604090205481565b3480156102aa57600080fd5b506102be6102b93660046120de565b610b89565b60405161017f9190612165565b3480156102d757600080fd5b5061030a6102e636600461217f565b600b6020526000908152604090208054600182015460029092015490919060ff1683565b60408051938452602084019290925215159082015260600161017f565b34801561033357600080fd5b5061015d6103423660046121be565b610c35565b34801561035357600080fd5b5061015d61036236600461225b565b610d7d565b34801561037357600080fd5b5062989680610175565b34801561038957600080fd5b50600654610175565b34801561039e57600080fd5b5061017560095481565b3480156103b457600080fd5b5061017560075481565b3480156103ca57600080fd5b5061015d6103d936600461217f565b610f56565b3480156103ea57600080fd5b50600954610175565b3480156103ff57600080fd5b506102be61040e3660046120de565b611044565b34801561041f57600080fd5b50610428611054565b60405161017f919061230c565b34801561044157600080fd5b5061044a6110c3565b60405161017f91906123bb565b34801561046357600080fd5b50610175600a5481565b336000908152600b602052604090206002015460ff166104ee576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f746572730000000000000060448201526064015b60405180910390fd5b600654336000908152600b602052604090205403610568576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f616c726561647920766f7465640000000000000000000000000000000000000060448201526064016104e5565b336000908152600b602052604081206001810180549087905581546006549092559181900361059857505061070b565b60005484146105a857505061070b565b60016006546105b791906123fd565b811415806105f45750848484336040516020016105d79493929190612410565b6040516020818303038152906040528051906020012060001c8214155b156106895760005b600054811015610681577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600c6000838154811061063c5761063c61246e565b9060005260206000200160405161065391906124f0565b908152604080516020928190038301902033600090815292529020558061067981612584565b9150506105fc565b50505061070b565b60005b84811015610707578585828181106106a6576106a661246e565b90506020020135600c600083815481106106c2576106c261246e565b906000526020600020016040516106d991906124f0565b90815260408051602092819003830190203360009081529252902055806106ff81612584565b91505061068c565b5050505b50505050565b61073c6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000600d600160065461074f91906123fd565b8154811061075f5761075f61246e565b9060005260206000200183604051610777919061259e565b908152602001604051809103902060405180606001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff1660018111156107c8576107c86125ba565b60018111156107d9576107d96125ba565b8152505090506000604051806080016040528060016006546107fb91906123fd565b815260200183600001518152602001836020015181526020018360400151600181111561082a5761082a6125ba565b9052949350505050565b61085f6040518060800160405280600081526020016000815260200160008152602001600081525090565b6000600d84815481106108745761087461246e565b906000526020600020018360405161088c919061259e565b908152602001604051809103902060405180606001604052908160008201548152602001600182015481526020016002820160009054906101000a900460ff1660018111156108dd576108dd6125ba565b60018111156108ee576108ee6125ba565b81525050905060006040518060800160405280868152602001836000015181526020018360200151815260200183604001516001811115610931576109316125ba565b9052925050505b92915050565b60025460009073ffffffffffffffffffffffffffffffffffffffff1633146109e8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e5565b600954600a546109f891906125e9565b4310610b835760005b600054811015610a2657610a1481611281565b610a1f6001826125e9565b9050610a01565b5060065460075403610ace5760005b600554811015610acc576001600b600060058481548110610a5857610a5861246e565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905580610ac481612584565b915050610a35565b505b600654600754610adf9060016125fc565b03610aec57610aec6115f6565b43600a81905550600160066000828254610b0691906125e9565b9091555050600854610b199060026125fc565b60065403610b335760018054610b3191600091611c87565b505b60065460095460408051928352436020840152429083015260608201527fb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e59060800160405180910390a150600190565b50600090565b60018181548110610b9957600080fd5b906000526020600020016000915090508054610bb49061249d565b80601f0160208091040260200160405190810160405280929190818152602001828054610be09061249d565b8015610c2d5780601f10610c0257610100808354040283529160200191610c2d565b820191906000526020600020905b815481529060010190602001808311610c1057829003601f168201915b505050505081565b60025473ffffffffffffffffffffffffffffffffffffffff163314610cdc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e5565b8051600003610d47576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f566f746572732063616e277420626520656d707479000000000000000000000060448201526064016104e5565b610d6081600060018451610d5b91906123fd565b61185d565b8051610d73906005906020840190611cdf565b5050600654600755565b60035473ffffffffffffffffffffffffffffffffffffffff163314610dfe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064016104e5565b8051600003610e69576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f73796d626f6c732063616e277420626520656d7074790000000000000000000060448201526064016104e5565b600654600854610e7a9060016125fc565b14158015610e8c575060065460085414155b610ef2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e64000060448201526064016104e5565b8051610f05906001906020840190611d65565b5060065460088190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d908290610f3d9060016125e9565b604051610f4b929190612624565b60405180910390a150565b60025473ffffffffffffffffffffffffffffffffffffffff163314610ffd576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e5565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60008181548110610b9957600080fd5b606060058054806020026020016040519081016040528092919081815260200182805480156110b957602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161108e575b5050505050905090565b606060065460085460016110d791906125fc565b036111b3576001805480602002602001604051908101604052809291908181526020016000905b828210156111aa57838290600052602060002001805461111d9061249d565b80601f01602080910402602001604051908101604052809291908181526020018280546111499061249d565b80156111965780601f1061116b57610100808354040283529160200191611196565b820191906000526020600020905b81548152906001019060200180831161117957829003601f168201915b5050505050815260200190600101906110fe565b50505050905090565b6000805480602002602001604051908101604052809291908181526020016000905b828210156111aa5783829060005260206000200180546111f49061249d565b80601f01602080910402602001604051908101604052809291908181526020018280546112209061249d565b801561126d5780601f106112425761010080835404028352916020019161126d565b820191906000526020600020905b81548152906001019060200180831161125057829003601f168201915b5050505050815260200190600101906111d5565b60008082815481106112955761129561246e565b9060005260206000200180546112aa9061249d565b80601f01602080910402602001604051908101604052809291908181526020018280546112d69061249d565b80156113235780601f106112f857610100808354040283529160200191611323565b820191906000526020600020905b81548152906001019060200180831161130657829003601f168201915b50505050509050600060048054905067ffffffffffffffff81111561134a5761134a611edf565b604051908082528060200260200182016040528015611373578160200160208202803683370190505b5090506000805b6004548110156114ca576000600482815481106113995761139961246e565b600091825260208083209091015460065473ffffffffffffffffffffffffffffffffffffffff909116808452600b90925260409092205490925014158061144057507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600c8660405161140c919061259e565b908152604080516020928190038301902073ffffffffffffffffffffffffffffffffffffffff851660009081529252902054145b1561144b57506114b8565b600c8560405161145b919061259e565b908152604080516020928190038301902073ffffffffffffffffffffffffffffffffffffffff841660009081529252902054848461149881612584565b9550815181106114aa576114aa61246e565b602002602001018181525050505b806114c281612584565b91505061137a565b506000600d60016006546114de91906123fd565b815481106114ee576114ee61246e565b9060005260206000200184604051611506919061259e565b908152604051908190036020019020549050600182156115315761152a8484611a4b565b9150600090505b600d805460019081018255600091909152604080516060810182528481524260208201529190820190839081111561156b5761156b6125ba565b815250600d600654815481106115835761158361246e565b906000526020600020018660405161159b919061259e565b9081526020016040518091039020600082015181600001556020820151816001015560408201518160020160006101000a81548160ff021916908360018111156115e7576115e76125ba565b02179055505050505050505050565b6000805b6004548210801561160c575060055481105b156117b057600581815481106116245761162461246e565b6000918252602090912001546004805473ffffffffffffffffffffffffffffffffffffffff909216918490811061165d5761165d61246e565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16036116a4578161168e81612584565b925050808061169c90612584565b9150506115fa565b600581815481106116b7576116b761246e565b6000918252602090912001546004805473ffffffffffffffffffffffffffffffffffffffff90921691849081106116f0576116f061246e565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff1610156117a657600b60006004848154811061172f5761172f61246e565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff1683528201929092526040018120818155600181019190915560020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558161179e81612584565b9250506115fa565b8061169c81612584565b60045482101561184857600b6000600484815481106117d1576117d161246e565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff1683528201929092526040018120818155600181019190915560020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558161184081612584565b9250506117b0565b6005805461185891600491611dab565b505050565b80821261186957505050565b8181600085600261187a8585612646565b611884919061269c565b61188e90876125fc565b8151811061189e5761189e61246e565b602002602001015190505b818313611a1d575b8073ffffffffffffffffffffffffffffffffffffffff168684815181106118da576118da61246e565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1610156119105782611908816126e6565b9350506118b1565b8073ffffffffffffffffffffffffffffffffffffffff168683815181106119395761193961246e565b602002602001015173ffffffffffffffffffffffffffffffffffffffff16111561196f578161196781612717565b925050611910565b818313611a18578582815181106119885761198861246e565b60200260200101518684815181106119a2576119a261246e565b60200260200101518785815181106119bc576119bc61246e565b602002602001018885815181106119d5576119d561246e565b73ffffffffffffffffffffffffffffffffffffffff93841660209182029290920101529116905282611a06816126e6565b9350508180611a1490612717565b9250505b6118a9565b81851215611a3057611a3086868461185d565b83831215611a4357611a4386848661185d565b505050505050565b600081600003611a5d57506000610938565b611a73836000611a6e6001866123fd565b611b0f565b6000611a80600284612750565b9050611a8d600284612764565b15611ab157838181518110611aa457611aa461246e565b6020026020010151611b07565b6002848281518110611ac557611ac561246e565b602002602001015185600184611adb91906123fd565b81518110611aeb57611aeb61246e565b6020026020010151611afd91906125fc565b611b07919061269c565b949350505050565b8181808203611b1f575050505050565b6000856002611b2e8787612646565b611b38919061269c565b611b4290876125fc565b81518110611b5257611b5261246e565b602002602001015190505b818313611c61575b80868481518110611b7857611b7861246e565b60200260200101511215611b985782611b90816126e6565b935050611b65565b858281518110611baa57611baa61246e565b6020026020010151811215611bcb5781611bc381612717565b925050611b98565b818313611c5c57858281518110611be457611be461246e565b6020026020010151868481518110611bfe57611bfe61246e565b6020026020010151878581518110611c1857611c1861246e565b60200260200101888581518110611c3157611c3161246e565b60209081029190910101919091525281611c4a81612717565b9250508280611c58906126e6565b9350505b611b5d565b81851215611c7457611c74868684611b0f565b83831215611a4357611a43868486611b0f565b828054828255906000526020600020908101928215611ccf5760005260206000209182015b82811115611ccf5781611cbf84826127be565b5091600101919060010190611cac565b50611cdb929150611deb565b5090565b828054828255906000526020600020908101928215611d59579160200282015b82811115611d5957825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190611cff565b50611cdb929150611e08565b828054828255906000526020600020908101928215611ccf579160200282015b82811115611ccf5782518290611d9b90826128bd565b5091602001919060010190611d85565b828054828255906000526020600020908101928215611d595760005260206000209182015b82811115611d59578254825591600101919060010190611dd0565b80821115611cdb576000611dff8282611e1d565b50600101611deb565b5b80821115611cdb5760008155600101611e09565b508054611e299061249d565b6000825580601f10611e39575050565b601f016020900490600052602060002090810190611e579190611e08565b50565b60008060008060608587031215611e7057600080fd5b84359350602085013567ffffffffffffffff80821115611e8f57600080fd5b818701915087601f830112611ea357600080fd5b813581811115611eb257600080fd5b8860208260051b8501011115611ec757600080fd5b95986020929092019750949560400135945092505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611f5557611f55611edf565b604052919050565b600082601f830112611f6e57600080fd5b813567ffffffffffffffff811115611f8857611f88611edf565b611fb960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611f0e565b818152846020838601011115611fce57600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215611ffd57600080fd5b813567ffffffffffffffff81111561201457600080fd5b611b0784828501611f5d565b6000806040838503121561203357600080fd5b82359150602083013567ffffffffffffffff81111561205157600080fd5b61205d85828601611f5d565b9150509250929050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461208b57600080fd5b919050565b600080604083850312156120a357600080fd5b823567ffffffffffffffff8111156120ba57600080fd5b6120c685828601611f5d565b9250506120d560208401612067565b90509250929050565b6000602082840312156120f057600080fd5b5035919050565b60005b838110156121125781810151838201526020016120fa565b50506000910152565b600081518084526121338160208601602086016120f7565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000612178602083018461211b565b9392505050565b60006020828403121561219157600080fd5b61217882612067565b600067ffffffffffffffff8211156121b4576121b4611edf565b5060051b60200190565b600060208083850312156121d157600080fd5b823567ffffffffffffffff8111156121e857600080fd5b8301601f810185136121f957600080fd5b803561220c6122078261219a565b611f0e565b81815260059190911b8201830190838101908783111561222b57600080fd5b928401925b828410156122505761224184612067565b82529284019290840190612230565b979650505050505050565b6000602080838503121561226e57600080fd5b823567ffffffffffffffff8082111561228657600080fd5b818501915085601f83011261229a57600080fd5b81356122a86122078261219a565b81815260059190911b830184019084810190888311156122c757600080fd5b8585015b838110156122ff578035858111156122e35760008081fd5b6122f18b89838a0101611f5d565b8452509186019186016122cb565b5098975050505050505050565b6020808252825182820181905260009190848201906040850190845b8181101561235a57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101612328565b50909695505050505050565b600081518084526020808501808196508360051b8101915082860160005b858110156123ae57828403895261239c84835161211b565b98850198935090840190600101612384565b5091979650505050505050565b6020815260006121786020830184612366565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b81810381811115610938576109386123ce565b60008186825b87811015612434578135835260209283019290910190600101612416565b5050938452505060601b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016602082015260340192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c908216806124b157607f821691505b6020821081036124ea577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60008083546124fe8161249d565b60018281168015612516576001811461254957612578565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0084168752821515830287019450612578565b8760005260208060002060005b8581101561256f5781548a820152908401908201612556565b50505082870194505b50929695505050505050565b60006000198203612597576125976123ce565b5060010190565b600082516125b08184602087016120f7565b9190910192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b80820180821115610938576109386123ce565b808201828112600083128015821682158216171561261c5761261c6123ce565b505092915050565b6040815260006126376040830185612366565b90508260208301529392505050565b8181036000831280158383131683831282161715612666576126666123ce565b5092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826126ab576126ab61266d565b60001983147f8000000000000000000000000000000000000000000000000000000000000000831416156126e1576126e16123ce565b500590565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203612597576125976123ce565b60007f80000000000000000000000000000000000000000000000000000000000000008203612748576127486123ce565b506000190190565b60008261275f5761275f61266d565b500490565b6000826127735761277361266d565b500690565b601f82111561185857600081815260208120601f850160051c8101602086101561279f5750805b601f850160051c820191505b81811015611a43578281556001016127ab565b8181036127c9575050565b6127d3825461249d565b67ffffffffffffffff8111156127eb576127eb611edf565b6127ff816127f9845461249d565b84612778565b6000601f821160018114612833576000831561281b5750848201545b600019600385901b1c1916600184901b1784556128b6565b6000858152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0841690600086815260209020845b8381101561288b578286015482556001958601959091019060200161286b565b50858310156128a95781850154600019600388901b60f8161c191681555b50505060018360011b0184555b5050505050565b815167ffffffffffffffff8111156128d7576128d7611edf565b6128e5816127f9845461249d565b602080601f83116001811461291a57600084156129025750858301515b600019600386901b1c1916600185901b178555611a43565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561296757888601518255948401946001909101908401612948565b50858210156129855787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea264697066735822122090ddb13e62244a38249b92947fe5aa99d9f56b1b77ba366ef314b2d8835d210364736f6c63430008150033")

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
