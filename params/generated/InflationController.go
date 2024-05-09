package generated

import (
	"strings"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
)

var InflationControllerBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b506040516114f03803806114f083398101604081905261002f91610058565b8051600055602081015160015560408101516002556060810151600355608001516004556100d5565b600060a0828403121561006a57600080fd5b60405160a081016001600160401b038111828210171561009a57634e487b7160e01b600052604160045260246000fd5b806040525082518152602083015160208201526040830151604082015260608301516060820152608083015160808201528091505092915050565b61140c806100e46000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635ed8f8db1461004657806392eff3cd1461006c578063cff0ab961461007f575b600080fd5b610059610054366004611229565b6100c1565b6040519081526020015b60405180910390f35b61005961007a366004611255565b61020b565b600054600154600254600354600454610099949392919085565b604080519586526020860194909452928401919091526060830152608082015260a001610063565b6000806100cd846102b0565b905060006100da846102b0565b905060006100f86100f0600060020154856103a1565b600354610510565b9050600061010e6100f0600060020154856103a1565b905060006101d161017c61012f600080015461012a888a610640565b6103a1565b6101776101526101486000800154600060010154610640565b61012a8a8c610640565b61017261016360006002015461064f565b61016d60016102b0565b610640565b610510565b6106bd565b6101776101b56101a06101986000600101546000800154610640565b6003546103a1565b61012a6101ac8861064f565b61016d8a61064f565b6101726101c961016360006002015461064f565b6002546103a1565b90506101fc6101f76101ee6101e58461064f565b61012a8d6102b0565b61016d8c6102b0565b6106d2565b955050505050505b9392505050565b600080610217846102b0565b90506000610224846102b0565b60035490915081136102445761023b8783836106e6565b925050506102a8565b60035482128015610256575060035481135b1561029857600061026d88846000600301546106e6565b9050600061028188600060030154856107ad565b905061028d81836112b6565b9450505050506102a8565b6102a38683836107ad565b925050505b949350505050565b60006102e4670de0b6b3a76400007f80000000000000000000000000000000000000000000000000000000000000006112f8565b821215610325576040517f99474eeb000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b610357670de0b6b3a76400007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6112f8565b821315610393576040517f9d5810910000000000000000000000000000000000000000000000000000000081526004810183905260240161031c565b50670de0b6b3a76400000290565b600082827f80000000000000000000000000000000000000000000000000000000000000008214806103f257507f800000000000000000000000000000000000000000000000000000000000000081145b15610429576040517fa6070c2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000841261043a578361043f565b836000035b91506000831261044f5782610454565b826000035b9050600061046283836107c9565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8111156104c8576040517f120b5b43000000000000000000000000000000000000000000000000000000008152600481018990526024810188905260440161031c565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff85851813610503816104fe5782600003610500565b825b90565b9998505050505050505050565b600082827f800000000000000000000000000000000000000000000000000000000000000082148061056157507f800000000000000000000000000000000000000000000000000000000000000081145b15610598576040517f9fe2b45000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080600084126105a957836105ae565b836000035b9150600083126105be57826105c3565b826000035b905060006105da83670de0b6b3a7640000846108d0565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8111156104c8576040517fd49c26b3000000000000000000000000000000000000000000000000000000008152600481018990526024810188905260440161031c565b60006102046105008385611387565b600081680736ea4425c11ac630811315610698576040517fca7ec0c50000000000000000000000000000000000000000000000000000000081526004810184905260240161031c565b6714057b7ef767814f81026102a86106b8670de0b6b3a7640000835b0590565b6109db565b600061020461050083856113ae565b92915050565b60006106cc670de0b6b3a7640000836112f8565b6000806107006000600201546106fc60006102b0565b1490565b1561072f576000546001546107289190610177906100f0906107229084610640565b886103a1565b905061077c565b60006107436100f0600060020154876103a1565b905060006107566101526101638461064f565b60005460015491925061077791610177906107719083610640565b846103a1565b925050505b600061079861078e8361012a896102b0565b61012a8688610640565b90506107a3816106d2565b9695505050505050565b60006102a86101f76107c161078e876102b0565b6004546103a1565b600080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8486098486029250828110838203039150508060000361081b5750670de0b6b3a7640000900490506106cc565b670de0b6b3a76400008110610866576040517f5173648d000000000000000000000000000000000000000000000000000000008152600481018690526024810185905260440161031c565b6000670de0b6b3a7640000858709620400008185030493109091037d40000000000000000000000000000000000000000000000000000000000002919091177faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106690291505092915050565b600080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff858709858702925082811083820303915050806000036109285783828161091e5761091e6112c9565b0492505050610204565b838110610972576040517f63a0577800000000000000000000000000000000000000000000000000000000815260048101879052602481018690526044810185905260640161031c565b60008486880960026001871981018816978890046003810283188082028403028082028403028082028403028082028403028082028403029081029092039091026000889003889004909101858311909403939093029303949094049190911702949350505050565b60008181811215610a4e577ffffffffffffffffffffffffffffffffffffffffffffffffcc22e87f6eb468eeb811215610a175750600092915050565b610a47610a2a6105006106b88460000390565b6ec097ce7bc90715b34b9f1000000000816106b4576106b46112c9565b9150610aaf565b680a688906bd8affffff811315610a94576040517f0360d0280000000000000000000000000000000000000000000000000000000081526004810184905260240161031c565b670de0b6b3a7640000604082901b056102a861050082610ab5565b50919050565b7780000000000000000000000000000000000000000000000067ff00000000000000821615610bd657678000000000000000821615610afd5768016a09e667f3bcc9090260401c5b674000000000000000821615610b1c576801306fe0a31b7152df0260401c5b672000000000000000821615610b3b576801172b83c7d517adce0260401c5b671000000000000000821615610b5a5768010b5586cf9890f62a0260401c5b670800000000000000821615610b79576801059b0d31585743ae0260401c5b670400000000000000821615610b9857680102c9a3e778060ee70260401c5b670200000000000000821615610bb75768010163da9fb33356d80260401c5b670100000000000000821615610bd657680100b1afa5abcbed610260401c5b66ff000000000000821615610cd5576680000000000000821615610c035768010058c86da1c09ea20260401c5b6640000000000000821615610c21576801002c605e2e8cec500260401c5b6620000000000000821615610c3f57680100162f3904051fa10260401c5b6610000000000000821615610c5d576801000b175effdc76ba0260401c5b6608000000000000821615610c7b57680100058ba01fb9f96d0260401c5b6604000000000000821615610c995768010002c5cc37da94920260401c5b6602000000000000821615610cb7576801000162e525ee05470260401c5b6601000000000000821615610cd55768010000b17255775c040260401c5b65ff0000000000821615610dcb5765800000000000821615610d00576801000058b91b5bc9ae0260401c5b65400000000000821615610d1d57680100002c5c89d5ec6d0260401c5b65200000000000821615610d3a5768010000162e43f4f8310260401c5b65100000000000821615610d5757680100000b1721bcfc9a0260401c5b65080000000000821615610d745768010000058b90cf1e6e0260401c5b65040000000000821615610d91576801000002c5c863b73f0260401c5b65020000000000821615610dae57680100000162e430e5a20260401c5b65010000000000821615610dcb576801000000b1721835510260401c5b64ff00000000821615610eb857648000000000821615610df457680100000058b90c0b490260401c5b644000000000821615610e105768010000002c5c8601cc0260401c5b642000000000821615610e2c576801000000162e42fff00260401c5b641000000000821615610e485768010000000b17217fbb0260401c5b640800000000821615610e64576801000000058b90bfce0260401c5b640400000000821615610e8057680100000002c5c85fe30260401c5b640200000000821615610e9c5768010000000162e42ff10260401c5b640100000000821615610eb857680100000000b17217f80260401c5b63ff000000821615610f9c576380000000821615610edf5768010000000058b90bfc0260401c5b6340000000821615610efa576801000000002c5c85fe0260401c5b6320000000821615610f1557680100000000162e42ff0260401c5b6310000000821615610f30576801000000000b17217f0260401c5b6308000000821615610f4b57680100000000058b90c00260401c5b6304000000821615610f665768010000000002c5c8600260401c5b6302000000821615610f81576801000000000162e4300260401c5b6301000000821615610f9c5768010000000000b172180260401c5b62ff00008216156110775762800000821615610fc1576801000000000058b90c0260401c5b62400000821615610fdb57680100000000002c5c860260401c5b62200000821615610ff55768010000000000162e430260401c5b6210000082161561100f57680100000000000b17210260401c5b620800008216156110295768010000000000058b910260401c5b62040000821615611043576801000000000002c5c80260401c5b6202000082161561105d57680100000000000162e40260401c5b62010000821615611077576801000000000000b1720260401c5b61ff008216156111495761800082161561109a57680100000000000058b90260401c5b6140008216156110b35768010000000000002c5d0260401c5b6120008216156110cc576801000000000000162e0260401c5b6110008216156110e55768010000000000000b170260401c5b6108008216156110fe576801000000000000058c0260401c5b61040082161561111757680100000000000002c60260401c5b61020082161561113057680100000000000001630260401c5b61010082161561114957680100000000000000b10260401c5b60ff82161561121257608082161561116a57680100000000000000590260401c5b6040821615611182576801000000000000002c0260401c5b602082161561119a57680100000000000000160260401c5b60108216156111b2576801000000000000000b0260401c5b60088216156111ca57680100000000000000060260401c5b60048216156111e257680100000000000000030260401c5b60028216156111fa57680100000000000000010260401c5b600182161561121257680100000000000000010260401c5b670de0b6b3a76400000260409190911c60bf031c90565b60008060006060848603121561123e57600080fd5b505081359360208301359350604090920135919050565b6000806000806080858703121561126b57600080fd5b5050823594602084013594506040840135936060013592509050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b808201808211156106cc576106cc611287565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261132e577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561138257611382611287565b500590565b81810360008312801583831316838312821617156113a7576113a7611287565b5092915050565b80820182811260008312801582168215821617156113ce576113ce611287565b50509291505056fea2646970667358221220ada4b51361750c969a760fa646df4b03ee1a5db6f30c102c2f7e6d4a9e5145fd64736f6c63430008150033")

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
