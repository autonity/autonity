package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var OmissionAccountabilityBytecode = common.Hex2Bytes("60806040523480156200001157600080fd5b506040516200254a3803806200254a833981016040819052620000349162000244565b601280546001600160a01b0319166001600160a01b0386161790558051600b55602080820151600c556040820151600d556060820151600e556080820151600f5560a082015160105560c08201516011558351620000999160009190860190620000ba565b508151620000af906001906020850190620000ba565b50505050506200032f565b82805482825590600052602060002090810192821562000112579160200282015b828111156200011257825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190620000db565b506200012092915062000124565b5090565b5b8082111562000120576000815560010162000125565b6001600160a01b03811681146200015157600080fd5b50565b634e487b7160e01b600052604160045260246000fd5b60405160e081016001600160401b03811182821017156200018f576200018f62000154565b60405290565b600082601f830112620001a757600080fd5b815160206001600160401b0380831115620001c657620001c662000154565b8260051b604051601f19603f83011681018181108482111715620001ee57620001ee62000154565b6040529384528581018301938381019250878511156200020d57600080fd5b83870191505b848210156200023957815162000229816200013b565b8352918301919083019062000213565b979650505050505050565b6000806000808486036101408112156200025d57600080fd5b85516200026a816200013b565b60208701519095506001600160401b03808211156200028857600080fd5b6200029689838a0162000195565b95506040880151915080821115620002ad57600080fd5b50620002bc8882890162000195565b93505060e0605f1982011215620002d257600080fd5b50620002dd6200016a565b606086015181526080860151602082015260a0860151604082015260c0860151606082015260e0860151608082015261010086015160a082015261012086015160c08201528091505092959194509250565b61220b806200033f6000396000f3fe60806040526004361061010e5760003560e01c8063b8d5712a116100a5578063d5baf90811610074578063eb231a1a11610059578063eb231a1a14610382578063eeb92233146103af578063f95bbd7f146103c257600080fd5b8063d5baf90814610342578063e1e8cac61461036257600080fd5b8063b8d5712a146102a2578063c1a48245146102dd578063ce4b5bbe146102ff578063d2aaca571461031557600080fd5b806370432e8b116100e157806370432e8b146101c657806379502c55146101f35780637f5e2f11146102575780639a11e0e61461026c57600080fd5b80630c820006146101135780631ede5a1a14610148578063278112dc1461016c5780635426b5ea14610199575b600080fd5b34801561011f57600080fd5b5061013361012e366004611912565b6103f2565b60405190151581526020015b60405180910390f35b34801561015457600080fd5b5061015e60065481565b60405190815260200161013f565b34801561017857600080fd5b5061015e61018736600461193e565b60086020526000908152604090205481565b3480156101a557600080fd5b5061015e6101b436600461193e565b60056020526000908152604090205481565b3480156101d257600080fd5b5061015e6101e136600461193e565b60096020526000908152604090205481565b3480156101ff57600080fd5b50600b54600c54600d54600e54600f546010546011546102229695949392919087565b604080519788526020880196909652948601939093526060850191909152608084015260a083015260c082015260e00161013f565b34801561026357600080fd5b5061271061015e565b34801561027857600080fd5b5061015e61028736600461193e565b6001600160a01b031660009081526008602052604090205490565b3480156102ae57600080fd5b506101336102bd366004611962565b600460209081526000928352604080842090915290825290205460ff1681565b3480156102e957600080fd5b506102fd6102f8366004611ad1565b6104bf565b005b34801561030b57600080fd5b5061015e61271081565b34801561032157600080fd5b5061015e61033036600461193e565b60076020526000908152604090205481565b34801561034e57600080fd5b506102fd61035d366004611b52565b6106d6565b34801561036e57600080fd5b506102fd61037d366004611bb6565b61079c565b34801561038e57600080fd5b5061015e61039d36600461193e565b600a6020526000908152604090205481565b6102fd6103bd366004611bb6565b61083b565b3480156103ce57600080fd5b506101336103dd366004611bb6565b60036020526000908152604090205460ff1681565b60006103ff600a43611bfe565b8210610491576040517f08c379a0000000000000000000000000000000000000000000000000000000008152602060048201526024808201527f746f6f20667574757265206865696768742066726f6d2063757272656e74207360448201527f746174650000000000000000000000000000000000000000000000000000000060648201526084015b60405180910390fd5b5060008181526004602090815260408083206001600160a01b038616845290915290205460ff165b92915050565b6012546001600160a01b03163314610559576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e747261637400000000000000000000000000000000000000006064820152608401610488565b6000610566600a43611bfe565b905082156105d157600081815260036020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556001600160a01b0388168352600590915281208054916105c783611c11565b9190505550610654565b600081815260036020908152604080832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556001600160a01b038816835260079091528120805486929061062b908490611c49565b9250508190555083600660008282546106449190611c49565b9091555061065490508682610b62565b81156106ce576000610664610d68565b905061066f81610f5e565b60005b6000548110156106cb5760006005600080848154811061069457610694611c5c565b60009182526020808320909101546001600160a01b03168352820192909252604001902055806106c381611c11565b915050610672565b50505b505050505050565b6012546001600160a01b03163314610770576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e747261637400000000000000000000000000000000000000006064820152608401610488565b8151610783906000906020850190611868565b508051610797906001906020840190611868565b505050565b6012546001600160a01b03163314610836576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e747261637400000000000000000000000000000000000000006064820152608401610488565b600255565b6012546001600160a01b031633146108d5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602c60248201527f66756e6374696f6e207265737472696374656420746f20746865204175746f6e60448201527f69747920436f6e747261637400000000000000000000000000000000000000006064820152608401610488565b3460005b600054811015610b58576000600760008084815481106108fb576108fb611c5c565b60009182526020808320909101546001600160a01b031683528201929092526040019020541115610b46576000600654836007600080868154811061094257610942611c5c565b60009182526020808320909101546001600160a01b031683528201929092526040019020546109719190611c8b565b61097b9190611ca2565b90506000600654856007600080878154811061099957610999611c5c565b60009182526020808320909101546001600160a01b031683528201929092526040019020546109c89190611c8b565b6109d29190611ca2565b9050600183815481106109e7576109e7611c5c565b60009182526020822001546040516001600160a01b03909116916108fc918591818181858888f193505050503d8060008114610a3f576040519150601f19603f3d011682016040523d82523d6000602084013e610a44565b606091505b5050601254600180546001600160a01b03909216925063a9059cbb9186908110610a7057610a70611c5c565b60009182526020909120015460405160e083901b7fffffffff000000000000000000000000000000000000000000000000000000001681526001600160a01b039091166004820152602481018490526044016020604051808303816000875af1158015610ae1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b059190611cdd565b50600060076000808681548110610b1e57610b1e611c5c565b60009182526020808320909101546001600160a01b0316835282019290925260400190205550505b80610b5081611c11565b9150506108d9565b5050600060065550565b60005b8251811015610bf6576000828152600460205260408120845160019290869085908110610b9457610b94611c5c565b6020908102919091018101516001600160a01b0316825281019190915260400160002080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001691151591909117905580610bee81611c11565b915050610b65565b50600c54600254610c079190611c49565b811015610c12575050565b60005b825181101561079757600c546001906000610c308386611bfe565b90505b610c3d8286611bfe565b811115610cf75760008181526003602052604090205460ff1615610c8a578160025486610c6a9190611bfe565b11610c785760009250610cf7565b81610c8281611c11565b925050610ce5565b600460008281526020019081526020016000206000878681518110610cb157610cb1611c5c565b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff16610ce55760009250610cf7565b80610cef81611cfa565b915050610c33565b508115610d535760056000868581518110610d1457610d14611c5c565b60200260200101516001600160a01b03166001600160a01b031681526020019081526020016000206000815480929190610d4d90611c11565b91905055505b50508080610d6090611c11565b915050610c15565b600080601260009054906101000a90046001600160a01b03166001600160a01b031663dfb1a4d26040518163ffffffff1660e01b8152600401602060405180830381865afa158015610dbe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610de29190611d2f565b90506000805b600054811015610f5757600c54600090600a90610e059086611bfe565b610e10906001611c49565b610e1a9190611bfe565b61271060056000808681548110610e3357610e33611c5c565b60009182526020808320909101546001600160a01b03168352820192909252604001902054610e629190611c8b565b610e6c9190611ca2565b90506000612710600b6002015460086000808781548110610e8f57610e8f611c5c565b60009182526020808320909101546001600160a01b03168352820192909252604001902054610ebe9190611c8b565b600d54610ecd90612710611bfe565b610ed79085611c8b565b610ee19190611c49565b610eeb9190611ca2565b600b54909150811115610f065783610f0281611c11565b9450505b8060086000808681548110610f1d57610f1d611c5c565b60009182526020808320909101546001600160a01b0316835282019290925260400190205550819050610f4f81611c11565b915050610de8565b5092915050565b60005b60005481101561144d576012546000805490916001600160a01b031690631904bb2e90839085908110610f9657610f96611c5c565b60009182526020909120015460405160e083901b7fffffffff000000000000000000000000000000000000000000000000000000001681526001600160a01b039091166004820152602401600060405180830381865afa158015610ffe573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01682016040526110449190810190611e15565b90506002816102600151600381111561105f5761105f611f93565b148061108157506003816102600151600381111561107f5761107f611f93565b145b1561108c575061143b565b600b60000154600860008085815481106110a8576110a8611c5c565b60009182526020808320909101546001600160a01b03168352820192909252604001902054116111e4576000600960008085815481106110ea576110ea611c5c565b60009182526020808320909101546001600160a01b0316835282019290925260400190205411156111df576009600080848154811061112b5761112b611c5c565b60009182526020808320909101546001600160a01b03168352820192909252604001812080549161115b83611cfa565b91905055506009600080848154811061117657611176611c5c565b60009182526020808320909101546001600160a01b0316835282019290925260400181205490036111df576000600a60008085815481106111b9576111b9611c5c565b60009182526020808320909101546001600160a01b031683528201929092526040019020555b611439565b600a60008084815481106111fa576111fa611c5c565b60009182526020808320909101546001600160a01b03168352820192909252604001812080549161122a83611c11565b91905055506000600a600080858154811061124757611247611c5c565b60009182526020808320909101546001600160a01b0316835282019290925260400181205481549091600a9181908790811061128557611285611c5c565b60009182526020808320909101546001600160a01b031683528201929092526040019020546112b49190611c8b565b9050600081600b600301546112c99190611c8b565b9050600082600b600401546112de9190611c8b565b90506112ea8243611c49565b61020085015260026102608501818152505060006009600080888154811061131457611314611c5c565b60009182526020808320909101546001600160a01b03168352820192909252604001902054111561136957611364848785600b600501546113559190611c8b565b61135f9190611c8b565b611451565b6113e5565b6012546040517f35be16e00000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906335be16e0906113b2908790600401612047565b600060405180830381600087803b1580156113cc57600080fd5b505af11580156113e0573d6000803e3d6000fd5b505050505b80600960008088815481106113fc576113fc611c5c565b60009182526020808320909101546001600160a01b0316835282019290925260400181208054909190611430908490611c49565b90915550505050505b505b8061144581611c11565b915050610f61565b5050565b60115481111561146057506011545b60008261012001518360c001518460a0015161147c9190611c49565b6114869190611c49565b6011549091506000906114998385611c8b565b6114a39190611ca2565b90506000811180156114b457508181145b156115d657600060a085018190526101008501819052610120850181905260c08501526101e0840180518291906114ec908390611c49565b905250600361026085015260006102008501526012546040517f35be16e00000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906335be16e090611548908790600401612047565b600060405180830381600087803b15801561156257600080fd5b505af1158015611576573d6000803e3d6000fd5b50505050602084810151604080516001600160a01b0390921682529181018390526000818301526001606082015290517f3cac37f432247a020a7112d5052bc279f35e1e3b80b0aab0eca49d1773ed3e3f9181900360800190a150505050565b61012084015181908111611603578085610120018181516115f79190611bfe565b9052506000905061161e565b6101208501516116139082611bfe565b600061012087015290505b801561169b578085610100015110611666578085610100018181516116439190611bfe565b90525060a08501805182919061165a908390611bfe565b9052506000905061169b565b6101008501516116769082611bfe565b90508461010001518560a00181815161168f9190611bfe565b90525060006101008601525b6000811180156116be575060008560a001518660c001516116bc9190611c49565b115b1561176a5760008560a001518660c001516116d99190611c49565b60c08701516116e89084611c8b565b6116f29190611ca2565b905060008660a001518760c0015161170a9190611c49565b60a08801516117199085611c8b565b6117239190611ca2565b9050818760c0018181516117379190611bfe565b90525060a08701805182919061174e908390611bfe565b90525061175b8183611c49565b6117659084611bfe565b925050505b6117748183611bfe565b915081856101e0018181516117899190611c49565b9052506012546040517f35be16e00000000000000000000000000000000000000000000000000000000081526001600160a01b03909116906335be16e0906117d5908890600401612047565b600060405180830381600087803b1580156117ef57600080fd5b505af1158015611803573d6000803e3d6000fd5b50505050602085810151610200870151604080516001600160a01b039093168352928201859052818301526000606082015290517f3cac37f432247a020a7112d5052bc279f35e1e3b80b0aab0eca49d1773ed3e3f9181900360800190a15050505050565b8280548282559060005260206000209081019282156118d5579160200282015b828111156118d557825182547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b03909116178255602090920191600190910190611888565b506118e19291506118e5565b5090565b5b808211156118e157600081556001016118e6565b6001600160a01b038116811461190f57600080fd5b50565b6000806040838503121561192557600080fd5b8235611930816118fa565b946020939093013593505050565b60006020828403121561195057600080fd5b813561195b816118fa565b9392505050565b6000806040838503121561197557600080fd5b823591506020830135611987816118fa565b809150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051610280810167ffffffffffffffff811182821017156119e5576119e5611992565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715611a3257611a32611992565b604052919050565b600082601f830112611a4b57600080fd5b8135602067ffffffffffffffff821115611a6757611a67611992565b8160051b611a768282016119eb565b9283528481018201928281019087851115611a9057600080fd5b83870192505b84831015611ab8578235611aa9816118fa565b82529183019190830190611a96565b979650505050505050565b801515811461190f57600080fd5b600080600080600060a08688031215611ae957600080fd5b853567ffffffffffffffff811115611b0057600080fd5b611b0c88828901611a3a565b9550506020860135611b1d816118fa565b9350604086013592506060860135611b3481611ac3565b91506080860135611b4481611ac3565b809150509295509295909350565b60008060408385031215611b6557600080fd5b823567ffffffffffffffff80821115611b7d57600080fd5b611b8986838701611a3a565b93506020850135915080821115611b9f57600080fd5b50611bac85828601611a3a565b9150509250929050565b600060208284031215611bc857600080fd5b5035919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156104b9576104b9611bcf565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611c4257611c42611bcf565b5060010190565b808201808211156104b9576104b9611bcf565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b80820281158282048414176104b9576104b9611bcf565b600082611cd8577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600060208284031215611cef57600080fd5b815161195b81611ac3565b600081611d0957611d09611bcf565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b600060208284031215611d4157600080fd5b5051919050565b8051611d53816118fa565b919050565b60005b83811015611d73578181015183820152602001611d5b565b50506000910152565b600082601f830112611d8d57600080fd5b815167ffffffffffffffff811115611da757611da7611992565b611dd860207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116016119eb565b818152846020838601011115611ded57600080fd5b611dfe826020830160208701611d58565b949350505050565b805160048110611d5357600080fd5b600060208284031215611e2757600080fd5b815167ffffffffffffffff80821115611e3f57600080fd5b908301906102808286031215611e5457600080fd5b611e5c6119c1565b611e6583611d48565b8152611e7360208401611d48565b6020820152611e8460408401611d48565b6040820152606083015182811115611e9b57600080fd5b611ea787828601611d7c565b6060830152506080830151608082015260a083015160a082015260c083015160c082015260e083015160e0820152610100808401518183015250610120808401518183015250610140808401518183015250610160808401518183015250610180611f13818501611d48565b908201526101a083810151908201526101c080840151908201526101e08084015190820152610200808401519082015261022080840151908201526102408084015183811115611f6257600080fd5b611f6e88828701611d7c565b8284015250506102609150611f84828401611e06565b91810191909152949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b60008151808452611fda816020860160208601611d58565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b60048110612043577f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b9052565b602081526120616020820183516001600160a01b03169052565b6000602083015161207d60408401826001600160a01b03169052565b5060408301516001600160a01b03811660608401525060608301516102808060808501526120af6102a0850183611fc2565b9150608085015160a085015260a085015160c085015260c085015160e085015260e08501516101008181870152808701519150506101208181870152808701519150506101408181870152808701519150506101608181870152808701519150506101808181870152808701519150506101a0612136818701836001600160a01b03169052565b8601516101c0868101919091528601516101e0808701919091528601516102008087019190915286015161022080870191909152860151610240808701919091528601518584037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001610260808801919091529091506121b68483611fc2565b9350808701519150506121cb8286018261200c565b509094935050505056fea2646970667358221220ff6d7ffa806d10677d53eed166dbcead60c6820b6cdbfce44827e5c57a4b196364736f6c63430008150033")

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
            "name" : "_validator",
            "type" : "address"
         },
         {
            "internalType" : "uint256",
            "name" : "_height",
            "type" : "uint256"
         }
      ],
      "name" : "isValidatorFaultyOnHeight",
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
