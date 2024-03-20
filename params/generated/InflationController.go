package generated

import (
	"strings"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
)

var InflationControllerBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b506040516113b93803806113b983398101604081905261002f9161005c565b805160005560208101516001556040810151600255606081015160035560800151600455426005556100d9565b600060a0828403121561006e57600080fd5b60405160a081016001600160401b038111828210171561009e57634e487b7160e01b600052604160045260246000fd5b806040525082518152602083015160208201526040830151604082015260608301516060820152608083015160808201528091505092915050565b6112d1806100e86000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806392eff3cd1461003b578063cff0ab9614610061575b600080fd5b61004e610049366004611107565b6100a3565b6040519081526020015b60405180910390f35b60005460015460025460035460045461007b949392919085565b604080519586526020860194909452928401919091526060830152608082015260a001610058565b6000806100bc600554856100b79190611168565b610150565b905060006100d1600554856100b79190611168565b60035490915081136100f1576100e8878383610241565b92505050610148565b60035482121561013857600061010d8884600060030154610241565b905060006101218860006003015485610341565b905061012d818361117b565b945050505050610148565b610143868383610341565b925050505b949350505050565b6000610184670de0b6b3a76400007f80000000000000000000000000000000000000000000000000000000000000006111bd565b8212156101c5576040517f99474eeb000000000000000000000000000000000000000000000000000000008152600481018390526024015b60405180910390fd5b6101f7670de0b6b3a76400007f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff6111bd565b821315610233576040517f9d581091000000000000000000000000000000000000000000000000000000008152600481018390526024016101bc565b50670de0b6b3a76400000290565b60008061025b6000600201546102576000610150565b1490565b1561029757600054600154610290919061028b906102839061027d9084610365565b88610374565b6003546104e0565b610610565b9050610309565b60006102ab61028360006002015487610374565b905060006102e36102cd6102be84610625565b6102c86001610150565b610365565b6102de6102be600060020154610625565b6104e0565b6000546001549192506103049161028b906102fe9083610365565b84610374565b925050505b600061032a6103208361031b89610150565b610374565b61031b8688610365565b905061033581610693565b925050505b9392505050565b600061014861035d61035561032087610150565b600454610374565b610693565b90565b600061033a610362838561124c565b600082827f80000000000000000000000000000000000000000000000000000000000000008214806103c557507f800000000000000000000000000000000000000000000000000000000000000081145b156103fc576040517fa6070c2500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000841261040d5783610412565b836000035b9150600083126104225782610427565b826000035b9050600061043583836106a7565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81111561049b576040517f120b5b4300000000000000000000000000000000000000000000000000000000815260048101899052602481018890526044016101bc565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff858518136104d3816104cf578260000390565b8290565b9998505050505050505050565b600082827f800000000000000000000000000000000000000000000000000000000000000082148061053157507f800000000000000000000000000000000000000000000000000000000000000081145b15610568576040517f9fe2b45000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060008412610579578361057e565b836000035b91506000831261058e5782610593565b826000035b905060006105aa83670de0b6b3a7640000846107ae565b90507f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81111561049b576040517fd49c26b300000000000000000000000000000000000000000000000000000000815260048101899052602481018890526044016101bc565b600061033a6103628385611273565b92915050565b600081680736ea4425c11ac63081131561066e576040517fca7ec0c5000000000000000000000000000000000000000000000000000000008152600481018490526024016101bc565b6714057b7ef767814f810261014861068e670de0b6b3a7640000835b0590565b6108b9565b600061061f670de0b6b3a7640000836111bd565b600080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff848609848602925082811083820303915050806000036106f95750670de0b6b3a76400009004905061061f565b670de0b6b3a76400008110610744576040517f5173648d00000000000000000000000000000000000000000000000000000000815260048101869052602481018590526044016101bc565b6000670de0b6b3a7640000858709620400008185030493109091037d40000000000000000000000000000000000000000000000000000000000002919091177faccb18165bd6fe31ae1cf318dc5b51eee0e1ba569b88cd74c1773b91fac106690291505092915050565b600080807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff85870985870292508281108382030391505080600003610806578382816107fc576107fc61118e565b049250505061033a565b838110610850576040517f63a057780000000000000000000000000000000000000000000000000000000081526004810187905260248101869052604481018590526064016101bc565b60008486880960026001871981018816978890046003810283188082028403028082028403028082028403028082028403028082028403029081029092039091026000889003889004909101858311909403939093029303949094049190911702949350505050565b6000818181121561092c577ffffffffffffffffffffffffffffffffffffffffffffffffcc22e87f6eb468eeb8112156108f55750600092915050565b61092561090861036261068e8460000390565b6ec097ce7bc90715b34b9f10000000008161068a5761068a61118e565b915061098d565b680a688906bd8affffff811315610972576040517f0360d028000000000000000000000000000000000000000000000000000000008152600481018490526024016101bc565b670de0b6b3a7640000604082901b0561014861036282610993565b50919050565b7780000000000000000000000000000000000000000000000067ff00000000000000821615610ab4576780000000000000008216156109db5768016a09e667f3bcc9090260401c5b6740000000000000008216156109fa576801306fe0a31b7152df0260401c5b672000000000000000821615610a19576801172b83c7d517adce0260401c5b671000000000000000821615610a385768010b5586cf9890f62a0260401c5b670800000000000000821615610a57576801059b0d31585743ae0260401c5b670400000000000000821615610a7657680102c9a3e778060ee70260401c5b670200000000000000821615610a955768010163da9fb33356d80260401c5b670100000000000000821615610ab457680100b1afa5abcbed610260401c5b66ff000000000000821615610bb3576680000000000000821615610ae15768010058c86da1c09ea20260401c5b6640000000000000821615610aff576801002c605e2e8cec500260401c5b6620000000000000821615610b1d57680100162f3904051fa10260401c5b6610000000000000821615610b3b576801000b175effdc76ba0260401c5b6608000000000000821615610b5957680100058ba01fb9f96d0260401c5b6604000000000000821615610b775768010002c5cc37da94920260401c5b6602000000000000821615610b95576801000162e525ee05470260401c5b6601000000000000821615610bb35768010000b17255775c040260401c5b65ff0000000000821615610ca95765800000000000821615610bde576801000058b91b5bc9ae0260401c5b65400000000000821615610bfb57680100002c5c89d5ec6d0260401c5b65200000000000821615610c185768010000162e43f4f8310260401c5b65100000000000821615610c3557680100000b1721bcfc9a0260401c5b65080000000000821615610c525768010000058b90cf1e6e0260401c5b65040000000000821615610c6f576801000002c5c863b73f0260401c5b65020000000000821615610c8c57680100000162e430e5a20260401c5b65010000000000821615610ca9576801000000b1721835510260401c5b64ff00000000821615610d9657648000000000821615610cd257680100000058b90c0b490260401c5b644000000000821615610cee5768010000002c5c8601cc0260401c5b642000000000821615610d0a576801000000162e42fff00260401c5b641000000000821615610d265768010000000b17217fbb0260401c5b640800000000821615610d42576801000000058b90bfce0260401c5b640400000000821615610d5e57680100000002c5c85fe30260401c5b640200000000821615610d7a5768010000000162e42ff10260401c5b640100000000821615610d9657680100000000b17217f80260401c5b63ff000000821615610e7a576380000000821615610dbd5768010000000058b90bfc0260401c5b6340000000821615610dd8576801000000002c5c85fe0260401c5b6320000000821615610df357680100000000162e42ff0260401c5b6310000000821615610e0e576801000000000b17217f0260401c5b6308000000821615610e2957680100000000058b90c00260401c5b6304000000821615610e445768010000000002c5c8600260401c5b6302000000821615610e5f576801000000000162e4300260401c5b6301000000821615610e7a5768010000000000b172180260401c5b62ff0000821615610f555762800000821615610e9f576801000000000058b90c0260401c5b62400000821615610eb957680100000000002c5c860260401c5b62200000821615610ed35768010000000000162e430260401c5b62100000821615610eed57680100000000000b17210260401c5b62080000821615610f075768010000000000058b910260401c5b62040000821615610f21576801000000000002c5c80260401c5b62020000821615610f3b57680100000000000162e40260401c5b62010000821615610f55576801000000000000b1720260401c5b61ff0082161561102757618000821615610f7857680100000000000058b90260401c5b614000821615610f915768010000000000002c5d0260401c5b612000821615610faa576801000000000000162e0260401c5b611000821615610fc35768010000000000000b170260401c5b610800821615610fdc576801000000000000058c0260401c5b610400821615610ff557680100000000000002c60260401c5b61020082161561100e57680100000000000001630260401c5b61010082161561102757680100000000000000b10260401c5b60ff8216156110f057608082161561104857680100000000000000590260401c5b6040821615611060576801000000000000002c0260401c5b602082161561107857680100000000000000160260401c5b6010821615611090576801000000000000000b0260401c5b60088216156110a857680100000000000000060260401c5b60048216156110c057680100000000000000030260401c5b60028216156110d857680100000000000000010260401c5b60018216156110f057680100000000000000010260401c5b670de0b6b3a76400000260409190911c60bf031c90565b6000806000806080858703121561111d57600080fd5b5050823594602084013594506040840135936060013592509050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561061f5761061f611139565b8082018082111561061f5761061f611139565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826111f3577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff83147f80000000000000000000000000000000000000000000000000000000000000008314161561124757611247611139565b500590565b818103600083128015838313168383128216171561126c5761126c611139565b5092915050565b808201828112600083128015821682158216171561129357611293611139565b50509291505056fea26469706673582212206d0cd18427e74f6a70c81d82ceedb94150fd0b5c5693c6d1505d32a1fce5038c64736f6c63430008150033")

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