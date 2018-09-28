package soma

import (
	"encoding/hex"
	"testing"

	"gitlab.clearmatics.net/oss/autonity/common"
	"gitlab.clearmatics.net/oss/autonity/core/vm/runtime"
)

func TestEVMRuntimeCall(t *testing.T) {
	/*
		pragma solidity ^0.4.25;

		contract Test {
			function test() public pure returns(string) {
				return "Hello Test!!!";
			}
		}
	*/
	contractBytecode := "608060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063f8a8fd6d14610046575b600080fd5b34801561005257600080fd5b5061005b6100d6565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561009b578082015181840152602081019050610080565b50505050905090810190601f1680156100c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040805190810160405280600d81526020017f48656c6c6f2054657374212121000000000000000000000000000000000000008152509050905600a165627a7a723058207d86d1462ac765f7f77965f34f8ad38a8fa270361ddfe7def03b516d6d6e4d120029"
	// (new Buffer(utils.sha3('test()'), 16)).toString().slice(0,8+2)
	input, err := hex.DecodeString("f8a8fd6d")
	if err != nil {
		t.Log(err)
	}

	ret, _, err := runtime.Execute(common.Hex2Bytes(contractBytecode), input, nil)
	if err != nil {
		t.Log(err)
	}
	// firstPart := ret[:32]
	// secondPart := ret[32:(32*2)] // size of the string (which is 13)
	retStr := ret[(32 * 2) : (32*2)+13] // third part the data itself
	if "Hello Test!!!" != string(retStr) {
		t.Error("Call() result different from expected: ", ret)
	}
}
