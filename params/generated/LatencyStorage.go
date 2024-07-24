package generated

import "strings"
import "github.com/autonity/autonity/accounts/abi"
import "github.com/autonity/autonity/common"

var LatencyStorageBytecode = common.Hex2Bytes("608060405234801561001057600080fd5b506108da806100206000396000f3fe608060405234801561001057600080fd5b50600436106100675760003560e01c806383a8123e1161005057806383a8123e146100a5578063b07ee6e3146100ba578063bb0f5556146100da57600080fd5b806341ed39111461006c57806346e1e76214610092575b600080fd5b61007f61007a36600461052b565b6100ed565b6040519081526020015b60405180910390f35b61007f6100a036600461052b565b610125565b6100b86100b3366004610629565b6101b2565b005b6100cd6100c83660046106eb565b610278565b6040516100899190610796565b6100cd6100e83660046106eb565b6103fa565b73ffffffffffffffffffffffffffffffffffffffff808316600090815260208181526040808320938516835292905220545b92915050565b73ffffffffffffffffffffffffffffffffffffffff8083166000818152602081815260408083209486168352938152838220548282528483209383529290529182205481158015906101775750600081115b1561018f57610186818361082a565b9250505061011f565b811561019f57610186828061082a565b6101a9818061082a565b95945050505050565b3360005b8251811015610273578281815181106101d1576101d161083d565b6020026020010151602001516000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600085848151811061022f5761022f61083d565b6020908102919091018101515173ffffffffffffffffffffffffffffffffffffffff168252810191909152604001600020558061026b8161086c565b9150506101b6565b505050565b60606000825167ffffffffffffffff8111156102965761029661055e565b6040519080825280602002602001820160405280156102db57816020015b60408051808201909152600080825260208201528152602001906001900390816102b45790505b50905060005b81518110156103f25760405180604001604052808583815181106103075761030761083d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1681526020016000808873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600087858151811061037c5761037c61083d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020548152508282815181106103d4576103d461083d565b602002602001018190525080806103ea9061086c565b9150506102e1565b509392505050565b60606000825167ffffffffffffffff8111156104185761041861055e565b60405190808252806020026020018201604052801561045d57816020015b60408051808201909152600080825260208201528152602001906001900390816104365790505b50905060005b81518110156103f25760405180604001604052808583815181106104895761048961083d565b602002602001015173ffffffffffffffffffffffffffffffffffffffff1681526020016104cf878785815181106104c2576104c261083d565b6020026020010151610125565b8152508282815181106104e4576104e461083d565b602002602001018190525080806104fa9061086c565b915050610463565b803573ffffffffffffffffffffffffffffffffffffffff8116811461052657600080fd5b919050565b6000806040838503121561053e57600080fd5b61054783610502565b915061055560208401610502565b90509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040805190810167ffffffffffffffff811182821017156105b0576105b061055e565b60405290565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff811182821017156105fd576105fd61055e565b604052919050565b600067ffffffffffffffff82111561061f5761061f61055e565b5060051b60200190565b6000602080838503121561063c57600080fd5b823567ffffffffffffffff81111561065357600080fd5b8301601f8101851361066457600080fd5b803561067761067282610605565b6105b6565b81815260069190911b8201830190838101908783111561069657600080fd5b928401925b828410156106e057604084890312156106b45760008081fd5b6106bc61058d565b6106c585610502565b8152848601358682015282526040909301929084019061069b565b979650505050505050565b600080604083850312156106fe57600080fd5b61070783610502565b915060208084013567ffffffffffffffff81111561072457600080fd5b8401601f8101861361073557600080fd5b803561074361067282610605565b81815260059190911b8201830190838101908883111561076257600080fd5b928401925b828410156107875761077884610502565b82529284019290840190610767565b80955050505050509250929050565b602080825282518282018190526000919060409081850190868401855b828110156107ee578151805173ffffffffffffffffffffffffffffffffffffffff1685528601518685015292840192908501906001016107b3565b5091979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082018082111561011f5761011f6107fb565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361089d5761089d6107fb565b506001019056fea26469706673582212203059a30572798468ef41ca041bcedc7ed9cae2cf0824ae3693e0917f72dda3cd64736f6c63430008150033")

var LatencyStorageAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "constructor"
   },
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_sender",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_receiver",
            "type" : "address"
         }
      ],
      "name" : "getLatency",
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
            "name" : "_sender",
            "type" : "address"
         },
         {
            "internalType" : "address[]",
            "name" : "_receiver",
            "type" : "address[]"
         }
      ],
      "name" : "getMultipleLatency",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "address",
                  "name" : "receiver",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "name" : "latency",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct LatencyStorage.Latency[]",
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
            "internalType" : "address",
            "name" : "_sender",
            "type" : "address"
         },
         {
            "internalType" : "address[]",
            "name" : "_receiver",
            "type" : "address[]"
         }
      ],
      "name" : "getMultipleTotalLatency",
      "outputs" : [
         {
            "components" : [
               {
                  "internalType" : "address",
                  "name" : "receiver",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "name" : "latency",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct LatencyStorage.Latency[]",
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
            "internalType" : "address",
            "name" : "_sender",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_receiver",
            "type" : "address"
         }
      ],
      "name" : "getTotalLatency",
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
            "components" : [
               {
                  "internalType" : "address",
                  "name" : "receiver",
                  "type" : "address"
               },
               {
                  "internalType" : "uint256",
                  "name" : "latency",
                  "type" : "uint256"
               }
            ],
            "internalType" : "struct LatencyStorage.Latency[]",
            "name" : "_latency",
            "type" : "tuple[]"
         }
      ],
      "name" : "updateLatency",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   }
]
`))
