package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OracleBytecode = common.Hex2Bytes("6080604052600160ff1b600755600160ff1b6008553480156200002157600080fd5b50604051620032df380380620032df833981016040819052620000449162000639565b600280546001600160a01b038087166001600160a01b03199283161790925560038054928616929091169190911790558151620000899060009060208501906200035f565b5081516200009f9060019060208501906200035f565b5080600981905550620000c485600060018851620000be91906200074e565b62000181565b8451620000d9906004906020880190620003bc565b508451620000ef906005906020880190620003bc565b5060016006819055600d8054909101815560009081525b855181101562000175576001600b60008884815181106200012b576200012b6200076a565b6020908102919091018101516001600160a01b03168252810191909152604001600020600201805460ff1916911515919091179055806200016c8162000780565b91505062000106565b505050505050620009c3565b8082126200018e57505050565b81816000856002620001a185856200079c565b620001ad9190620007c6565b620001b9908762000806565b81518110620001cc57620001cc6200076a565b602002602001015190505b8183136200032b575b806001600160a01b0316868481518110620001ff57620001ff6200076a565b60200260200101516001600160a01b031610156200022c5782620002238162000831565b935050620001e0565b806001600160a01b03168683815181106200024b576200024b6200076a565b60200260200101516001600160a01b031611156200027857816200026f816200084c565b9250506200022c565b81831362000325578582815181106200029557620002956200076a565b6020026020010151868481518110620002b257620002b26200076a565b6020026020010151878581518110620002cf57620002cf6200076a565b60200260200101888581518110620002eb57620002eb6200076a565b6001600160a01b0393841660209182029290920101529116905282620003118162000831565b935050818062000321906200084c565b9250505b620001d7565b8185121562000341576200034186868462000181565b8383121562000357576200035786848662000181565b505050505050565b828054828255906000526020600020908101928215620003aa579160200282015b82811115620003aa5782518290620003999082620008f7565b509160200191906001019062000380565b50620003b892915062000422565b5090565b82805482825590600052602060002090810192821562000414579160200282015b828111156200041457825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620003dd565b50620003b892915062000443565b80821115620003b85760006200043982826200045a565b5060010162000422565b5b80821115620003b8576000815560010162000444565b50805462000468906200086c565b6000825580601f1062000479575050565b601f01602090049060005260206000209081019062000499919062000443565b50565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715620004dd57620004dd6200049c565b604052919050565b60006001600160401b038211156200050157620005016200049c565b5060051b60200190565b80516001600160a01b03811681146200052357600080fd5b919050565b6000601f83818401126200053b57600080fd5b82516020620005546200054e83620004e5565b620004b2565b82815260059290921b850181019181810190878411156200057457600080fd5b8287015b848110156200062d5780516001600160401b03808211156200059a5760008081fd5b818a0191508a603f830112620005b05760008081fd5b8582015181811115620005c757620005c76200049c565b620005da818a01601f19168801620004b2565b915080825260408c81838601011115620005f45760008081fd5b60005b8281101562000614578481018201518482018a01528801620005f7565b5050600090820187015284525091830191830162000578565b50979650505050505050565b600080600080600060a086880312156200065257600080fd5b85516001600160401b03808211156200066a57600080fd5b818801915088601f8301126200067f57600080fd5b81516020620006926200054e83620004e5565b82815260059290921b8401810191818101908c841115620006b257600080fd5b948201945b83861015620006db57620006cb866200050b565b82529482019490820190620006b7565b9950620006ec90508a82016200050b565b97505050620006fe604089016200050b565b945060608801519150808211156200071557600080fd5b50620007248882890162000528565b925050608086015190509295509295909350565b634e487b7160e01b600052601160045260246000fd5b8181038181111562000764576200076462000738565b92915050565b634e487b7160e01b600052603260045260246000fd5b60006001820162000795576200079562000738565b5060010190565b8181036000831280158383131683831282161715620007bf57620007bf62000738565b5092915050565b600082620007e457634e487b7160e01b600052601260045260246000fd5b600160ff1b82146000198414161562000801576200080162000738565b500590565b808201828112600083128015821682158216171562000829576200082962000738565b505092915050565b60006001600160ff1b01820162000795576200079562000738565b6000600160ff1b820162000864576200086462000738565b506000190190565b600181811c908216806200088157607f821691505b602082108103620008a257634e487b7160e01b600052602260045260246000fd5b50919050565b601f821115620008f257600081815260208120601f850160051c81016020861015620008d15750805b601f850160051c820191505b818110156200035757828155600101620008dd565b505050565b81516001600160401b038111156200091357620009136200049c565b6200092b816200092484546200086c565b84620008a8565b602080601f8311600181146200096357600084156200094a5750858301515b600019600386901b1c1916600185901b17855562000357565b600085815260208120601f198616915b82811015620009945788860151825594840194600190910190840162000973565b5085821015620009b35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b61290c80620009d36000396000f3fe6080604052600436106101565760003560e01c80638d4f75d2116100bf578063b3ab15fb11610079578063cdd7225311610056578063cdd7225314610415578063df7f710e14610437578063e6a02a281461045957005b8063b3ab15fb146103c0578063b78dec52146103e0578063ccce413b146103f557005b80639f8743f7116100a75780639f8743f71461037f578063a781358714610394578063aa2f89b5146103aa57005b80638d4f75d2146103495780639670c0bc1461036957005b80634bb278f3116101105780635281b5c6116100f85780635281b5c6146102a05780635412b3ae146102cd578063845023f21461032957005b80634bb278f3146102335780634c56ea561461025857005b8063307de9b61161013e578063307de9b61461019e57806333f98c77146101be5780633c8510fd1461021357005b806308f21ff51461015f578063146ca5311461018857005b3661015d57005b005b34801561016b57600080fd5b5061017560085481565b6040519081526020015b60405180910390f35b34801561019457600080fd5b5061017560065481565b3480156101aa57600080fd5b5061015d6101b9366004611dca565b61046f565b3480156101ca57600080fd5b506101de6101d9366004611f5b565b610713565b60405161017f919081518152602080830151908201526040808301519082015260609182015115159181019190915260800190565b34801561021f57600080fd5b506101de61022e366004611f90565b6107fd565b34801561023f57600080fd5b506102486108be565b604051901515815260200161017f565b34801561026457600080fd5b50610175610273366004612000565b8151602081840181018051600c825292820194820194909420919093529091526000908152604090205481565b3480156102ac57600080fd5b506102c06102bb36600461204e565b610b09565b60405161017f91906120d5565b3480156102d957600080fd5b5061030c6102e83660046120ef565b600b6020526000908152604090208054600182015460029092015490919060ff1683565b60408051938452602084019290925215159082015260600161017f565b34801561033557600080fd5b5061015d61034436600461212e565b610bb5565b34801561035557600080fd5b5061015d6103643660046121cb565b610cfd565b34801561037557600080fd5b5062989680610175565b34801561038b57600080fd5b50600654610175565b3480156103a057600080fd5b5061017560095481565b3480156103b657600080fd5b5061017560075481565b3480156103cc57600080fd5b5061015d6103db3660046120ef565b610ed6565b3480156103ec57600080fd5b50600954610175565b34801561040157600080fd5b506102c061041036600461204e565b610fc4565b34801561042157600080fd5b5061042a610fd4565b60405161017f919061227c565b34801561044357600080fd5b5061044c611043565b60405161017f919061232b565b34801561046557600080fd5b50610175600a5481565b336000908152600b602052604090206002015460ff166104f0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f7265737472696374656420746f206f6e6c7920766f746572730000000000000060448201526064015b60405180910390fd5b600654336000908152600b60205260409020540361056a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f616c726561647920766f7465640000000000000000000000000000000000000060448201526064016104e7565b336000908152600b602052604081206001810180549087905581546006549092559181900361059a57505061070d565b60005484146105aa57505061070d565b60016006546105b9919061236d565b811415806105f65750848484336040516020016105d99493929190612380565b6040516020818303038152906040528051906020012060001c8214155b1561068b5760005b600054811015610683577f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600c6000838154811061063e5761063e6123de565b906000526020600020016040516106559190612460565b908152604080516020928190038301902033600090815292529020558061067b816124f4565b9150506105fe565b50505061070d565b60005b84811015610709578585828181106106a8576106a86123de565b90506020020135600c600083815481106106c4576106c46123de565b906000526020600020016040516106db9190612460565b9081526040805160209281900383019020336000908152925290205580610701816124f4565b91505061068e565b5050505b50505050565b61074060405180608001604052806000815260200160008152602001600081526020016000151581525090565b6000600d6001600654610753919061236d565b81548110610763576107636123de565b906000526020600020018360405161077b919061250e565b90815260408051918290036020908101832060608401835280548452600180820154928501929092526002015460ff1615158383015281516080810190925260065492935060009282916107ce9161236d565b815260200183600001518152602001836020015181526020018360400151151581525090508092505050919050565b61082a60405180608001604052806000815260200160008152602001600081526020016000151581525090565b6000600d848154811061083f5761083f6123de565b9060005260206000200183604051610857919061250e565b9081526040805191829003602090810183206060808501845281548552600182015485840190815260029092015460ff16151585850190815284516080810186528a8152955193860193909352905192840192909252511515908201529150505b92915050565b60025460009073ffffffffffffffffffffffffffffffffffffffff163314610968576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e7565b600954600a54610978919061252a565b4310610b035760005b6000548110156109a65761099481611201565b61099f60018261252a565b9050610981565b5060065460075403610a4e5760005b600554811015610a4c576001600b6000600584815481106109d8576109d86123de565b60009182526020808320919091015473ffffffffffffffffffffffffffffffffffffffff168352820192909252604001902060020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905580610a44816124f4565b9150506109b5565b505b600654600754610a5f90600161253d565b03610a6c57610a6c611566565b43600a81905550600160066000828254610a86919061252a565b9091555050600854610a9990600261253d565b60065403610ab35760018054610ab191600091611bf7565b505b60065460095460408051928352436020840152429083015260608201527fb5d8636ab45e6cac7a4a61cb7c77f77f61a454d73aa2e6139ff8dcaf463537e59060800160405180910390a150600190565b50600090565b60018181548110610b1957600080fd5b906000526020600020016000915090508054610b349061240d565b80601f0160208091040260200160405190810160405280929190818152602001828054610b609061240d565b8015610bad5780601f10610b8257610100808354040283529160200191610bad565b820191906000526020600020905b815481529060010190602001808311610b9057829003601f168201915b505050505081565b60025473ffffffffffffffffffffffffffffffffffffffff163314610c5c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e7565b8051600003610cc7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f566f746572732063616e277420626520656d707479000000000000000000000060448201526064016104e7565b610ce081600060018451610cdb919061236d565b6117cd565b8051610cf3906005906020840190611c4f565b5050600654600755565b60035473ffffffffffffffffffffffffffffffffffffffff163314610d7e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f7265737472696374656420746f206f70657261746f720000000000000000000060448201526064016104e7565b8051600003610de9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f73796d626f6c732063616e277420626520656d7074790000000000000000000060448201526064016104e7565b600654600854610dfa90600161253d565b14158015610e0c575060065460085414155b610e72576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f63616e2774206265207570646174656420696e207468697320726f756e64000060448201526064016104e7565b8051610e85906001906020840190611cd5565b5060065460088190557faa278e424da680ce5dad66510415760e78e0bd87d45c786c6e88bdde82f9342d908290610ebd90600161252a565b604051610ecb929190612565565b60405180910390a150565b60025473ffffffffffffffffffffffffffffffffffffffff163314610f7d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f7265737472696374656420746f20746865206175746f6e69747920636f6e747260448201527f616374000000000000000000000000000000000000000000000000000000000060648201526084016104e7565b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60008181548110610b1957600080fd5b6060600580548060200260200160405190810160405280929190818152602001828054801561103957602002820191906000526020600020905b815473ffffffffffffffffffffffffffffffffffffffff16815260019091019060200180831161100e575b5050505050905090565b60606006546008546001611057919061253d565b03611133576001805480602002602001604051908101604052809291908181526020016000905b8282101561112a57838290600052602060002001805461109d9061240d565b80601f01602080910402602001604051908101604052809291908181526020018280546110c99061240d565b80156111165780601f106110eb57610100808354040283529160200191611116565b820191906000526020600020905b8154815290600101906020018083116110f957829003601f168201915b50505050508152602001906001019061107e565b50505050905090565b6000805480602002602001604051908101604052809291908181526020016000905b8282101561112a5783829060005260206000200180546111749061240d565b80601f01602080910402602001604051908101604052809291908181526020018280546111a09061240d565b80156111ed5780601f106111c2576101008083540402835291602001916111ed565b820191906000526020600020905b8154815290600101906020018083116111d057829003601f168201915b505050505081526020019060010190611155565b6000808281548110611215576112156123de565b90600052602060002001805461122a9061240d565b80601f01602080910402602001604051908101604052809291908181526020018280546112569061240d565b80156112a35780601f10611278576101008083540402835291602001916112a3565b820191906000526020600020905b81548152906001019060200180831161128657829003601f168201915b50505050509050600060048054905067ffffffffffffffff8111156112ca576112ca611e4f565b6040519080825280602002602001820160405280156112f3578160200160208202803683370190505b5090506000805b60045481101561144a57600060048281548110611319576113196123de565b600091825260208083209091015460065473ffffffffffffffffffffffffffffffffffffffff909116808452600b9092526040909220549092501415806113c057507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600c8660405161138c919061250e565b908152604080516020928190038301902073ffffffffffffffffffffffffffffffffffffffff851660009081529252902054145b156113cb5750611438565b600c856040516113db919061250e565b908152604080516020928190038301902073ffffffffffffffffffffffffffffffffffffffff8416600090815292529020548484611418816124f4565b95508151811061142a5761142a6123de565b602002602001018181525050505b80611442816124f4565b9150506112fa565b506000600d600160065461145e919061236d565b8154811061146e5761146e6123de565b9060005260206000200184604051611486919061250e565b908152604051908190036020019020549050600082156114b1576114aa84846119bb565b9150600190505b600d80546001018082556000829052604080516060810182528581524260208201528415159181019190915260065490929181106114f1576114f16123de565b9060005260206000200186604051611509919061250e565b9081526040805160209281900383019020835181559183015160018301559190910151600290910180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016911515919091179055505050505050565b6000805b6004548210801561157c575060055481105b156117205760058181548110611594576115946123de565b6000918252602090912001546004805473ffffffffffffffffffffffffffffffffffffffff90921691849081106115cd576115cd6123de565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff160361161457816115fe816124f4565b925050808061160c906124f4565b91505061156a565b60058181548110611627576116276123de565b6000918252602090912001546004805473ffffffffffffffffffffffffffffffffffffffff9092169184908110611660576116606123de565b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16101561171657600b60006004848154811061169f5761169f6123de565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff1683528201929092526040018120818155600181019190915560020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690558161170e816124f4565b92505061156a565b8061160c816124f4565b6004548210156117b857600b600060048481548110611741576117416123de565b600091825260208083209091015473ffffffffffffffffffffffffffffffffffffffff1683528201929092526040018120818155600181019190915560020180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055816117b0816124f4565b925050611720565b600580546117c891600491611d1b565b505050565b8082126117d957505050565b818160008560026117ea8585612587565b6117f491906125dd565b6117fe908761253d565b8151811061180e5761180e6123de565b602002602001015190505b81831361198d575b8073ffffffffffffffffffffffffffffffffffffffff1686848151811061184a5761184a6123de565b602002602001015173ffffffffffffffffffffffffffffffffffffffff161015611880578261187881612627565b935050611821565b8073ffffffffffffffffffffffffffffffffffffffff168683815181106118a9576118a96123de565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1611156118df57816118d781612658565b925050611880565b818313611988578582815181106118f8576118f86123de565b6020026020010151868481518110611912576119126123de565b602002602001015187858151811061192c5761192c6123de565b60200260200101888581518110611945576119456123de565b73ffffffffffffffffffffffffffffffffffffffff9384166020918202929092010152911690528261197681612627565b935050818061198490612658565b9250505b611819565b818512156119a0576119a08686846117cd565b838312156119b3576119b38684866117cd565b505050505050565b6000816000036119cd575060006108b8565b6119e38360006119de60018661236d565b611a7f565b60006119f0600284612691565b90506119fd6002846126a5565b15611a2157838181518110611a1457611a146123de565b6020026020010151611a77565b6002848281518110611a3557611a356123de565b602002602001015185600184611a4b919061236d565b81518110611a5b57611a5b6123de565b6020026020010151611a6d919061253d565b611a7791906125dd565b949350505050565b8181808203611a8f575050505050565b6000856002611a9e8787612587565b611aa891906125dd565b611ab2908761253d565b81518110611ac257611ac26123de565b602002602001015190505b818313611bd1575b80868481518110611ae857611ae86123de565b60200260200101511215611b085782611b0081612627565b935050611ad5565b858281518110611b1a57611b1a6123de565b6020026020010151811215611b3b5781611b3381612658565b925050611b08565b818313611bcc57858281518110611b5457611b546123de565b6020026020010151868481518110611b6e57611b6e6123de565b6020026020010151878581518110611b8857611b886123de565b60200260200101888581518110611ba157611ba16123de565b60209081029190910101919091525281611bba81612658565b9250508280611bc890612627565b9350505b611acd565b81851215611be457611be4868684611a7f565b838312156119b3576119b3868486611a7f565b828054828255906000526020600020908101928215611c3f5760005260206000209182015b82811115611c3f5781611c2f84826126ff565b5091600101919060010190611c1c565b50611c4b929150611d5b565b5090565b828054828255906000526020600020908101928215611cc9579160200282015b82811115611cc957825182547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff909116178255602090920191600190910190611c6f565b50611c4b929150611d78565b828054828255906000526020600020908101928215611c3f579160200282015b82811115611c3f5782518290611d0b90826127fe565b5091602001919060010190611cf5565b828054828255906000526020600020908101928215611cc95760005260206000209182015b82811115611cc9578254825591600101919060010190611d40565b80821115611c4b576000611d6f8282611d8d565b50600101611d5b565b5b80821115611c4b5760008155600101611d79565b508054611d999061240d565b6000825580601f10611da9575050565b601f016020900490600052602060002090810190611dc79190611d78565b50565b60008060008060608587031215611de057600080fd5b84359350602085013567ffffffffffffffff80821115611dff57600080fd5b818701915087601f830112611e1357600080fd5b813581811115611e2257600080fd5b8860208260051b8501011115611e3757600080fd5b95986020929092019750949560400135945092505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611ec557611ec5611e4f565b604052919050565b600082601f830112611ede57600080fd5b813567ffffffffffffffff811115611ef857611ef8611e4f565b611f2960207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601611e7e565b818152846020838601011115611f3e57600080fd5b816020850160208301376000918101602001919091529392505050565b600060208284031215611f6d57600080fd5b813567ffffffffffffffff811115611f8457600080fd5b611a7784828501611ecd565b60008060408385031215611fa357600080fd5b82359150602083013567ffffffffffffffff811115611fc157600080fd5b611fcd85828601611ecd565b9150509250929050565b803573ffffffffffffffffffffffffffffffffffffffff81168114611ffb57600080fd5b919050565b6000806040838503121561201357600080fd5b823567ffffffffffffffff81111561202a57600080fd5b61203685828601611ecd565b92505061204560208401611fd7565b90509250929050565b60006020828403121561206057600080fd5b5035919050565b60005b8381101561208257818101518382015260200161206a565b50506000910152565b600081518084526120a3816020860160208601612067565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b6020815260006120e8602083018461208b565b9392505050565b60006020828403121561210157600080fd5b6120e882611fd7565b600067ffffffffffffffff82111561212457612124611e4f565b5060051b60200190565b6000602080838503121561214157600080fd5b823567ffffffffffffffff81111561215857600080fd5b8301601f8101851361216957600080fd5b803561217c6121778261210a565b611e7e565b81815260059190911b8201830190838101908783111561219b57600080fd5b928401925b828410156121c0576121b184611fd7565b825292840192908401906121a0565b979650505050505050565b600060208083850312156121de57600080fd5b823567ffffffffffffffff808211156121f657600080fd5b818501915085601f83011261220a57600080fd5b81356122186121778261210a565b81815260059190911b8301840190848101908883111561223757600080fd5b8585015b8381101561226f578035858111156122535760008081fd5b6122618b89838a0101611ecd565b84525091860191860161223b565b5098975050505050505050565b6020808252825182820181905260009190848201906040850190845b818110156122ca57835173ffffffffffffffffffffffffffffffffffffffff1683529284019291840191600101612298565b50909695505050505050565b600081518084526020808501808196508360051b8101915082860160005b8581101561231e57828403895261230c84835161208b565b988501989350908401906001016122f4565b5091979650505050505050565b6020815260006120e860208301846122d6565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156108b8576108b861233e565b60008186825b878110156123a4578135835260209283019290910190600101612386565b5050938452505060601b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016602082015260340192915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600181811c9082168061242157607f821691505b60208210810361245a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b600080835461246e8161240d565b6001828116801561248657600181146124b9576124e8565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00841687528215158302870194506124e8565b8760005260208060002060005b858110156124df5781548a8201529084019082016124c6565b50505082870194505b50929695505050505050565b600060001982036125075761250761233e565b5060010190565b60008251612520818460208701612067565b9190910192915050565b808201808211156108b8576108b861233e565b808201828112600083128015821682158216171561255d5761255d61233e565b505092915050565b60408152600061257860408301856122d6565b90508260208301529392505050565b81810360008312801583831316838312821617156125a7576125a761233e565b5092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826125ec576125ec6125ae565b60001983147f8000000000000000000000000000000000000000000000000000000000000000831416156126225761262261233e565b500590565b60007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036125075761250761233e565b60007f800000000000000000000000000000000000000000000000000000000000000082036126895761268961233e565b506000190190565b6000826126a0576126a06125ae565b500490565b6000826126b4576126b46125ae565b500690565b601f8211156117c857600081815260208120601f850160051c810160208610156126e05750805b601f850160051c820191505b818110156119b3578281556001016126ec565b81810361270a575050565b612714825461240d565b67ffffffffffffffff81111561272c5761272c611e4f565b6127408161273a845461240d565b846126b9565b6000601f821160018114612774576000831561275c5750848201545b600019600385901b1c1916600184901b1784556127f7565b6000858152602090207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0841690600086815260209020845b838110156127cc57828601548255600195860195909101906020016127ac565b50858310156127ea5781850154600019600388901b60f8161c191681555b50505060018360011b0184555b5050505050565b815167ffffffffffffffff81111561281857612818611e4f565b6128268161273a845461240d565b602080601f83116001811461285b57600084156128435750858301515b600019600386901b1c1916600185901b1785556119b3565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156128a857888601518255948401946001909101908401612889565b50858210156128c65787850151600019600388901b60f8161c191681555b5050505050600190811b0190555056fea26469706673582212206892b7bfcaf5b2d94cf490a1b5d3f4c525d4e551760ac643fe2d78526f3c710b64736f6c63430008150033")

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
                  "internalType" : "bool",
                  "name" : "isValid",
                  "type" : "bool"
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
                  "internalType" : "bool",
                  "name" : "isValid",
                  "type" : "bool"
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
