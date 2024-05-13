package generated

import (
	"strings"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
)

var InflationControllerBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b5060405161151f38038061151f83398101604081905261002f9161005c565b805160005560208101516001556040810151600255606081015160035560800151600455426005556100d9565b600060a0828403121561006e57600080fd5b60405160a081016001600160401b038111828210171561009e57634e487b7160e01b600052604160045260246000fd5b806040525082518152602083015160208201526040830151604082015260608301516060820152608083015160808201528091505092915050565b611437806100e86000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635ed8f8db1461004657806392eff3cd1461006c578063cff0ab961461007f575b600080fd5b610059610054366004611241565b6100c1565b6040519081526020015b60405180910390f35b61005961007a36600461126d565b610220565b600054600154600254600354600454610099949392919085565b604080519586526020860194909452928401919091526060830152608082015260a001610063565b6000806100da600554856100d591906112ce565b6102c8565b905060006100ef600554856100d591906112ce565b9050600061010d610105600060020154856103b9565b600354610528565b90506000610123610105600060020154856103b9565b905060006101e6610191610144600080015461013f888a610658565b6103b9565b61018c61016761015d6000800154600060010154610658565b61013f8a8c610658565b610187610178600060020154610667565b61018260016102c8565b610658565b610528565b6106d5565b61018c6101ca6101b56101ad6000600101546000800154610658565b6003546103b9565b61013f6101c188610667565b6101828a610667565b6101876101de610178600060020154610667565b6002546103b9565b905061021161020c6102036101fa84610667565b61013f8d6102c8565b6101828c6102c8565b6106ea565b955050505050505b9392505050565b600080610234600554856100d591906112ce565b90506000610249600554856100d591906112ce565b6003549091508113610269576102608783836106fe565b925050506102c0565b6003548212156102b057600061028588846000600301546106fe565b9050600061029988600060030154856107c5565b90506102a581836112e1565b9450505050506102c0565b6102bb8683836107c5565b925050505b949350505050565b60006102fc670de0b6b3a76400007f8000000000000000000000000000000000000000000000000000000000000000611323565b82121561033d576040517f99474eeb000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b61036f670de0b6b3a76400007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff611323565b8213156103ab576040517f9d58109100000000000000000000000000000000000000000000000000000000815260048101839052602401610334565b50670de0b6b3a76400000290565b600082827f800000000000000000000000000000000000000000000000000000000000000082148061040a57507f800000000000000000000000000000000000000000000000000000000000000081145b15610441576040517fa6070c2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600084126104525783610457565b836000035b915060008312610467578261046c565b826000035b9050600061047a83836107e1565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8111156104e0576040517f120b5b430000000000000000000000000000000000000000000000000000000081526004810189905260248101889052604401610334565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8585181361051b816105165782600003610518565b825b90565b9998505050505050505050565b600082827f800000000000000000000000000000000000000000000000000000000000000082148061057957507f800000000000000000000000000000000000000000000000000000000000000081145b156105b0576040517f9fe2b45000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600084126105c157836105c6565b836000035b9150600083126105d657826105db565b826000035b905060006105f283670de0b6b3a7640000846108e8565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8111156104e0576040517fd49c26b30000000000000000000000000000000000000000000000000000000081526004810189905260248101889052604401610334565b600061021961051883856113b2565b600081680736ea4425c11ac6308113156106b0576040517fca7ec0c500000000000000000000000000000000000000000000000000000000815260048101849052602401610334565b6714057b7ef767814f81026102c06106d0670de0b6b3a7640000835b0590565b6109f3565b600061021961051883856113d9565b92915050565b60006106e4670de0b6b3a764000083611323565b60008061071860006002015461071460006102c8565b1490565b1561074757600054600154610740919061018c906101059061073a9084610658565b886103b9565b9050610794565b600061075b610105600060020154876103b9565b9050600061076e61016761017884610667565b60005460015491925061078f9161018c906107899083610658565b846103b9565b925050505b60006107b06107a68361013f896102c8565b61013f8688610658565b90506107bb816106ea565b9695505050505050565b60006102c061020c6107d96107a6876102c8565b6004546103b9565b600080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff848609848602925082811083820303915050806000036108335750670de0b6b3a7640000900490506106e4565b670de0b6b3a7640000811061087e576040517f5173648d0000000000000000000000000000000000000000000000000000000081526004810186905260248101859052604401610334565b6000670de0b6b3a7640000858709620400008185030493109091037d40000000000000000000000000000000000000000000000000000000000002919091177faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106690291505092915050565b600080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8587098587029250828110838203039150508060000361094057838281610936576109366112f4565b0492505050610219565b83811061098a576040517f63a05778000000000000000000000000000000000000000000000000000000008152600481018790526024810186905260448101859052606401610334565b60008486880960026001871981018816978890046003810283188082028403028082028403028082028403028082028403028082028403029081029092039091026000889003889004909101858311909403939093029303949094049190911702949350505050565b60008181811215610a66577ffffffffffffffffffffffffffffffffffffffffffffffffcc22e87f6eb468eeb811215610a2f5750600092915050565b610a5f610a426105186106d08460000390565b6ec097ce7bc90715b34b9f1000000000816106cc576106cc6112f4565b9150610ac7565b680a688906bd8affffff811315610aac576040517f0360d02800000000000000000000000000000000000000000000000000000000815260048101849052602401610334565b670de0b6b3a7640000604082901b056102c061051882610acd565b50919050565b7780000000000000000000000000000000000000000000000067ff00000000000000821615610bee57678000000000000000821615610b155768016a09e667f3bcc9090260401c5b674000000000000000821615610b34576801306fe0a31b7152df0260401c5b672000000000000000821615610b53576801172b83c7d517adce0260401c5b671000000000000000821615610b725768010b5586cf9890f62a0260401c5b670800000000000000821615610b91576801059b0d31585743ae0260401c5b670400000000000000821615610bb057680102c9a3e778060ee70260401c5b670200000000000000821615610bcf5768010163da9fb33356d80260401c5b670100000000000000821615610bee57680100b1afa5abcbed610260401c5b66ff000000000000821615610ced576680000000000000821615610c1b5768010058c86da1c09ea20260401c5b6640000000000000821615610c39576801002c605e2e8cec500260401c5b6620000000000000821615610c5757680100162f3904051fa10260401c5b6610000000000000821615610c75576801000b175effdc76ba0260401c5b6608000000000000821615610c9357680100058ba01fb9f96d0260401c5b6604000000000000821615610cb15768010002c5cc37da94920260401c5b6602000000000000821615610ccf576801000162e525ee05470260401c5b6601000000000000821615610ced5768010000b17255775c040260401c5b65ff0000000000821615610de35765800000000000821615610d18576801000058b91b5bc9ae0260401c5b65400000000000821615610d3557680100002c5c89d5ec6d0260401c5b65200000000000821615610d525768010000162e43f4f8310260401c5b65100000000000821615610d6f57680100000b1721bcfc9a0260401c5b65080000000000821615610d8c5768010000058b90cf1e6e0260401c5b65040000000000821615610da9576801000002c5c863b73f0260401c5b65020000000000821615610dc657680100000162e430e5a20260401c5b65010000000000821615610de3576801000000b1721835510260401c5b64ff00000000821615610ed057648000000000821615610e0c57680100000058b90c0b490260401c5b644000000000821615610e285768010000002c5c8601cc0260401c5b642000000000821615610e44576801000000162e42fff00260401c5b641000000000821615610e605768010000000b17217fbb0260401c5b640800000000821615610e7c576801000000058b90bfce0260401c5b640400000000821615610e9857680100000002c5c85fe30260401c5b640200000000821615610eb45768010000000162e42ff10260401c5b640100000000821615610ed057680100000000b17217f80260401c5b63ff000000821615610fb4576380000000821615610ef75768010000000058b90bfc0260401c5b6340000000821615610f12576801000000002c5c85fe0260401c5b6320000000821615610f2d57680100000000162e42ff0260401c5b6310000000821615610f48576801000000000b17217f0260401c5b6308000000821615610f6357680100000000058b90c00260401c5b6304000000821615610f7e5768010000000002c5c8600260401c5b6302000000821615610f99576801000000000162e4300260401c5b6301000000821615610fb45768010000000000b172180260401c5b62ff000082161561108f5762800000821615610fd9576801000000000058b90c0260401c5b62400000821615610ff357680100000000002c5c860260401c5b6220000082161561100d5768010000000000162e430260401c5b6210000082161561102757680100000000000b17210260401c5b620800008216156110415768010000000000058b910260401c5b6204000082161561105b576801000000000002c5c80260401c5b6202000082161561107557680100000000000162e40260401c5b6201000082161561108f576801000000000000b1720260401c5b61ff00821615611161576180008216156110b257680100000000000058b90260401c5b6140008216156110cb5768010000000000002c5d0260401c5b6120008216156110e4576801000000000000162e0260401c5b6110008216156110fd5768010000000000000b170260401c5b610800821615611116576801000000000000058c0260401c5b61040082161561112f57680100000000000002c60260401c5b61020082161561114857680100000000000001630260401c5b61010082161561116157680100000000000000b10260401c5b60ff82161561122a57608082161561118257680100000000000000590260401c5b604082161561119a576801000000000000002c0260401c5b60208216156111b257680100000000000000160260401c5b60108216156111ca576801000000000000000b0260401c5b60088216156111e257680100000000000000060260401c5b60048216156111fa57680100000000000000030260401c5b600282161561121257680100000000000000010260401c5b600182161561122a57680100000000000000010260401c5b670de0b6b3a76400000260409190911c60bf031c90565b60008060006060848603121561125657600080fd5b505081359360208301359350604090920135919050565b6000806000806080858703121561128357600080fd5b5050823594602084013594506040840135936060013592509050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b818103818111156106e4576106e461129f565b808201808211156106e4576106e461129f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b600082611359577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f8000000000000000000000000000000000000000000000000000000000000000831416156113ad576113ad61129f565b500590565b81810360008312801583831316838312821617156113d2576113d261129f565b5092915050565b80820182811260008312801582168215821617156113f9576113f961129f565b50509291505056fea2646970667358221220af4e5c7086dd4e119b99f2b28f3caa6c1a29137ee1d12d01a7a206fc96eb037464736f6c63430008150033")

var InflationControllerAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "components" : [
               {
                  "internalType" : "SD59x18",
                  "name" : "iInit",
                  "type" : "int256"
               },
               {
                  "internalType" : "SD59x18",
                  "name" : "iTrans",
                  "type" : "int256"
               },
               {
                  "internalType" : "SD59x18",
                  "name" : "aE",
                  "type" : "int256"
               },
               {
                  "internalType" : "SD59x18",
                  "name" : "T",
                  "type" : "int256"
               },
               {
                  "internalType" : "SD59x18",
                  "name" : "iPerm",
                  "type" : "int256"
               }
            ],
            "internalType" : "struct InflationController.Params",
            "name" : "_params",
            "type" : "tuple"
         }
      ],
      "stateMutability" : "nonpayable",
      "type" : "constructor"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "x",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "y",
            "type" : "uint256"
         }
      ],
      "name" : "PRBMath_MulDiv18_Overflow",
      "type" : "error"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "x",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "y",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "denominator",
            "type" : "uint256"
         }
      ],
      "name" : "PRBMath_MulDiv_Overflow",
      "type" : "error"
   },
   {
      "inputs" : [
         {
            "internalType" : "int256",
            "name" : "x",
            "type" : "int256"
         }
      ],
      "name" : "PRBMath_SD59x18_Convert_Overflow",
      "type" : "error"
   },
   {
      "inputs" : [
         {
            "internalType" : "int256",
            "name" : "x",
            "type" : "int256"
         }
      ],
      "name" : "PRBMath_SD59x18_Convert_Underflow",
      "type" : "error"
   },
   {
      "inputs" : [],
      "name" : "PRBMath_SD59x18_Div_InputTooSmall",
      "type" : "error"
   },
   {
      "inputs" : [
         {
            "internalType" : "SD59x18",
            "name" : "x",
            "type" : "int256"
         },
         {
            "internalType" : "SD59x18",
            "name" : "y",
            "type" : "int256"
         }
      ],
      "name" : "PRBMath_SD59x18_Div_Overflow",
      "type" : "error"
   },
   {
      "inputs" : [
         {
            "internalType" : "SD59x18",
            "name" : "x",
            "type" : "int256"
         }
      ],
      "name" : "PRBMath_SD59x18_Exp2_InputTooBig",
      "type" : "error"
   },
   {
      "inputs" : [
         {
            "internalType" : "SD59x18",
            "name" : "x",
            "type" : "int256"
         }
      ],
      "name" : "PRBMath_SD59x18_Exp_InputTooBig",
      "type" : "error"
   },
   {
      "inputs" : [],
      "name" : "PRBMath_SD59x18_Mul_InputTooSmall",
      "type" : "error"
   },
   {
      "inputs" : [
         {
            "internalType" : "SD59x18",
            "name" : "x",
            "type" : "int256"
         },
         {
            "internalType" : "SD59x18",
            "name" : "y",
            "type" : "int256"
         }
      ],
      "name" : "PRBMath_SD59x18_Mul_Overflow",
      "type" : "error"
   },
   {
      "inputs" : [
         {
            "internalType" : "uint256",
            "name" : "_currentSupply",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_inflationReserve",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_lastEpochTime",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_currentEpochTime",
            "type" : "uint256"
         }
      ],
      "name" : "calculateSupplyDelta",
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
            "name" : "_currentSupply",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_lastEpochTime",
            "type" : "uint256"
         },
         {
            "internalType" : "uint256",
            "name" : "_currentEpochTime",
            "type" : "uint256"
         }
      ],
      "name" : "calculateSupplyDeltaOLD",
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
      "name" : "params",
      "outputs" : [
         {
            "internalType" : "SD59x18",
            "name" : "iInit",
            "type" : "int256"
         },
         {
            "internalType" : "SD59x18",
            "name" : "iTrans",
            "type" : "int256"
         },
         {
            "internalType" : "SD59x18",
            "name" : "aE",
            "type" : "int256"
         },
         {
            "internalType" : "SD59x18",
            "name" : "T",
            "type" : "int256"
         },
         {
            "internalType" : "SD59x18",
            "name" : "iPerm",
            "type" : "int256"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   }
]
`))
