package generated

import (
	"strings"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
)

var UpgradeManagerBytecode = common.Hex2Bytes("608060405234801561000f575f80fd5b5060405161052f38038061052f83398101604081905261002e91610079565b5f80546001600160a01b039384166001600160a01b031991821617909155600180549290931691161790556100aa565b80516001600160a01b0381168114610074575f80fd5b919050565b5f806040838503121561008a575f80fd5b6100938361005e565b91506100a16020840161005e565b90509250929050565b610478806100b75f395ff3fe608060405234801561000f575f80fd5b506004361061004a575f3560e01c806355463ceb1461004e578063570ca735146100965780636e3d9ff0146100b6578063b3ab15fb146100cb575b5f80fd5b5f5461006d9073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b60015461006d9073ffffffffffffffffffffffffffffffffffffffff1681565b6100c96100c43660046102ed565b6100de565b005b6100c96100d93660046103c7565b6101ab565b60015473ffffffffffffffffffffffffffffffffffffffff163314610164576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f63616c6c6572206973206e6f7420746865206f70657261746f7200000000000060448201526064015b60405180910390fd5b60405160f9905f9061017c90859085906020016103e7565b60405160208183030381529060405290505f80825160208401855af43d5f803e8080156101a7573d5ff35b3d5ffd5b5f5473ffffffffffffffffffffffffffffffffffffffff163314610251576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602360248201527f63616c6c6572206973206e6f7420746865204175746f6e69747920636f6e747260448201527f6163740000000000000000000000000000000000000000000000000000000000606482015260840161015b565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b803573ffffffffffffffffffffffffffffffffffffffff811681146102bb575f80fd5b919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f80604083850312156102fe575f80fd5b61030783610298565b9150602083013567ffffffffffffffff80821115610323575f80fd5b818501915085601f830112610336575f80fd5b813581811115610348576103486102c0565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561038e5761038e6102c0565b816040528281528860208487010111156103a6575f80fd5b826020860160208301375f6020848301015280955050505050509250929050565b5f602082840312156103d7575f80fd5b6103e082610298565b9392505050565b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008360601b1681525f82515f5b818110156104315760208186018101516014868401015201610414565b505f9201601401918252509291505056fea2646970667358221220c7f4c45807952dc1031e4526c44d3cd18993e8e9e2b0d57102011f249160735f64736f6c63430008150033")

var UpgradeManagerAbi, _ = abi.JSON(strings.NewReader(`[
   {
      "inputs" : [
         {
            "internalType" : "address",
            "name" : "_autonity",
            "type" : "address"
         },
         {
            "internalType" : "address",
            "name" : "_operator",
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
            "internalType" : "address",
            "name" : "",
            "type" : "address"
         }
      ],
      "stateMutability" : "view",
      "type" : "function"
   },
   {
      "inputs" : [],
      "name" : "operator",
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
            "internalType" : "address",
            "name" : "_account",
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
            "internalType" : "address",
            "name" : "_target",
            "type" : "address"
         },
         {
            "internalType" : "string",
            "name" : "_data",
            "type" : "string"
         }
      ],
      "name" : "upgrade",
      "outputs" : [],
      "stateMutability" : "nonpayable",
      "type" : "function"
   }
]
`))
