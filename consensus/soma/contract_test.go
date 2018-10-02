package soma

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/core/vm/runtime"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
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
	// firstPart := ret[:32] // what is this?
	// secondPart := ret[32:(32*2)] // size of the string (which is 13)
	retStr := ret[(32 * 2) : (32*2)+13] // third part the data itself
	if "Hello Test!!!" != string(retStr) {
		t.Error("Call() result different from expected: ", ret)
	}
}

func MakePreState(db ethdb.Database, accounts core.GenesisAlloc) *state.StateDB {
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb)
	for addr, a := range accounts {
		statedb.SetCode(addr, a.Code)
		statedb.SetNonce(addr, a.Nonce)
		statedb.SetBalance(addr, a.Balance)
		for k, v := range a.Storage {
			statedb.SetState(addr, k, v)
		}
	}
	// Commit and re-open to start with a clean state.
	root, _ := statedb.Commit(false)
	statedb, _ = state.New(root, sdb)
	return statedb
}

func TestEVMContractDeployment(t *testing.T) {
	initialBalance := big.NewInt(1000000000)
	userKey, _ := crypto.GenerateKey()
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey)
	coinbaseKey, _ := crypto.GenerateKey()
	coinbaseAddr := crypto.PubkeyToAddress(coinbaseKey.PublicKey)
	originKey, _ := crypto.GenerateKey()
	originAddr := crypto.PubkeyToAddress(originKey.PublicKey)

	alloc := make(core.GenesisAlloc)
	alloc[userAddr] = core.GenesisAccount{
		Balance: initialBalance,
	}
	statedb := MakePreState(ethdb.NewMemDatabase(), alloc)
	initialBalanceUser := statedb.GetBalance(userAddr)
	t.Log("initialBalanceUser:\t", initialBalanceUser)
	initialBalanceCoinbase := statedb.GetBalance(coinbaseAddr)
	t.Log("initialBalanceCoinbase:\t", initialBalanceCoinbase)
	initialBalanceorigin := statedb.GetBalance(originAddr)
	t.Log("initialBalanceorigin:\t", initialBalanceorigin)

	vmTestBlockHash := func(n uint64) common.Hash {
		return common.BytesToHash(crypto.Keccak256([]byte(big.NewInt(int64(n)).String())))
	}
	context := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     vmTestBlockHash,
		Origin:      originAddr,
		Coinbase:    coinbaseAddr,
		BlockNumber: new(big.Int).SetUint64(0x00),
		Time:        new(big.Int).SetUint64(0x01),
		GasLimit:    uint64(0x0f4240),
		Difficulty:  new(big.Int).SetUint64(0x0100),
		GasPrice:    new(big.Int).SetUint64(0x3b9aca00),
	}
	vmconfig := vm.Config{}
	//vmconfig.NoRecursion = true
	evm := vm.NewEVM(context, statedb, params.MainnetChainConfig, vmconfig) //vmconfig)

	// CREATE
	sender := vm.AccountRef(userAddr)
	/*
		pragma solidity ^0.4.25;

		contract Test {
			function test() public pure returns(string) {
				return "Hello Test!!!";
			}

			int private count = 0;
			function incrementCounter() public {
				count += 1;
			}
			function decrementCounter() public {
				count -= 1;
			}
			function getCount() public view returns (int) {
				return count;
			}
		}
	*/
	contractBytecode := "60806040526000805534801561001457600080fd5b506101e6806100246000396000f300608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680635b34b96614610067578063a87d942c1461007e578063f5c5ad83146100a9578063f8a8fd6d146100c0575b600080fd5b34801561007357600080fd5b5061007c610150565b005b34801561008a57600080fd5b50610093610162565b6040518082815260200191505060405180910390f35b3480156100b557600080fd5b506100be61016b565b005b3480156100cc57600080fd5b506100d561017d565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101155780820151818401526020810190506100fa565b50505050905090810190601f1680156101425780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60016000808282540192505081905550565b60008054905090565b60016000808282540392505081905550565b60606040805190810160405280600d81526020017f48656c6c6f2054657374212121000000000000000000000000000000000000008152509050905600a165627a7a723058201b0858a814ecee293d6f73f3c8ed4b76a898989e7e0c3796fdb8db6a6c16884b0029"
	data := common.Hex2Bytes(contractBytecode)
	gas := uint64(1000000)
	value := new(big.Int).SetUint64(0x00)
	ret, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)
	t.Log("====== CREATE =======")
	t.Logf("Contract:\n%s\n", hex.Dump(ret))
	t.Log("Address: ", contractAddress.String())
	t.Log("Gas: ", gas)
	t.Log("Error: ", vmerr)
	// CALL
	functionSig := "test()"
	t.Log("====== CALL =======", functionSig)
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()
	ret, gas, vmerr = evm.Call(sender, contractAddress, input, gas, value)
	t.Logf("Result:\n%s\n", hex.Dump(ret))
	t.Log("Gas: ", gas)
	t.Log("Error: ", vmerr)
	// CALL
	functionSig = "getCount()"
	t.Log("====== CALL =======", functionSig)
	input = crypto.Keccak256Hash([]byte(functionSig)).Bytes()
	ret, gas, vmerr = evm.Call(sender, contractAddress, input, gas, value)
	t.Logf("Result:\n%s\n", hex.Dump(ret))
	t.Log("Gas: ", gas)
	t.Log("Error: ", vmerr)
	totalIncrements := new(big.Int).SetUint64(5)
	for i := uint64(0); i < totalIncrements.Uint64(); i++ {
		// CALL
		functionSig = "incrementCounter()"
		t.Log("====== CALL =======", functionSig)
		input = crypto.Keccak256Hash([]byte(functionSig)).Bytes()
		ret, gas, vmerr = evm.Call(sender, contractAddress, input, gas, value)
		t.Logf("Result:\n%s\n", hex.Dump(ret))
		t.Log("Gas: ", gas)
		t.Log("Error: ", vmerr)
	}
	// CALL
	functionSig = "getCount()"
	t.Log("====== CALL =======", functionSig)
	input = crypto.Keccak256Hash([]byte(functionSig)).Bytes()
	ret, gas, vmerr = evm.Call(sender, contractAddress, input, gas, value)
	t.Logf("Result:\n%s\n", hex.Dump(ret))
	t.Log("Gas: ", gas)
	t.Log("Error: ", vmerr)
	// CALL (TRANSFER)
	t.Log("====== CALL ======= TRANSFER")
	input = []byte{}
	value = new(big.Int).SetUint64(0x100)
	ret, gas, vmerr = evm.Call(sender, originAddr, nil, gas, value)
	t.Logf("Result:\n%s\n", hex.Dump(ret))
	t.Log("Gas: ", gas)
	t.Log("Error: ", vmerr)

	totalExpectedIncrements := new(big.Int).SetBytes(ret)
	if totalExpectedIncrements.Uint64() != totalIncrements.Uint64() {
		t.Error("Increments n smart contract and expected differ")
	}
	finalBalanceUser := statedb.GetBalance(userAddr)
	t.Log("finalBalanceUser:\t", finalBalanceUser)
	finalBalanceCoinbase := statedb.GetBalance(coinbaseAddr)
	t.Log("finalBalanceCoinbase:\t", finalBalanceCoinbase)
	finalBalanceorigin := statedb.GetBalance(originAddr)
	t.Log("finalBalanceorigin:\t", finalBalanceorigin)
}
