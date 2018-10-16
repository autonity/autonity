package soma

import (
	"encoding/hex"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

func deployContract(bytecodeStr string, userAddr common.Address, header *types.Header, statedb *state.StateDB) (common.Address, *types.Transaction) {
	contractBytecode := common.Hex2Bytes(bytecodeStr[2:]) // [2:] removes 0x

	gasPrice := new(big.Int).SetUint64(0x0)
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		//GetHash:     core.GetHashFn(header,chainContext),
		GetHash:     func(n uint64) common.Hash { return header.Root }, // since this a one time thing, no point in adding complex functions to get the hash
		Origin:      userAddr,
		Coinbase:    userAddr,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    gasPrice,
	}
	chainConfig := params.AllSomaProtocolChanges
	vmconfig := vm.Config{}
	/*
		// vm.Config
			type Config struct {
				// Debug enabled debugging Interpreter options
				Debug bool
				// Tracer is the op code logger
				Tracer Tracer
				// NoRecursion disabled Interpreter call, callcode,
				// delegate call and create.
				NoRecursion bool
				// Enable recording of SHA3/keccak preimages
				EnablePreimageRecording bool
				// JumpTable contains the EVM instruction table. This
				// may be left uninitialised and will be set to the default
				// table.
				JumpTable [256]operation
			}
	*/
	evm := vm.NewEVM(evmContext, statedb, chainConfig, vmconfig)

	sender := vm.AccountRef(userAddr)
	data := contractBytecode
	gas := uint64(1000000)
	value := new(big.Int).SetUint64(0x00)
	ret, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)
	log.Println("====== CREATE =======")
	log.Printf("Contract:\n%s\n", hex.Dump(ret))
	log.Println("Address: ", contractAddress.String())
	log.Println("Gas: ", gas)
	log.Println("Error: ", vmerr)

	// create transaction
	contractTx := types.NewContractCreation(statedb.GetNonce(contractAddress), value, header.GasLimit, gasPrice, data)

	// CALL
	functionSig := "ActiveValidator(address)"
	log.Println("====== CALL =======", functionSig)
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()[:4]
	inputData := append(input[:], userAddr[:]...)
	ret, gas, vmerr = evm.Call(sender, contractAddress, inputData, gas, value)
	log.Printf("Result:\n%s\n", hex.Dump(ret))
	log.Println("User Address: ", userAddr)
	log.Println("Gas: ", gas)
	log.Println("Error: ", vmerr)

	//statedb.Commit(false)
	//printDB(statedb.Database())

	return contractAddress, contractTx
}

func callContract(contractAddress common.Address, userAddr common.Address, header *types.Header, statedb *state.StateDB) {
	gasPrice := new(big.Int).SetUint64(0x0)
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		//GetHash:     core.GetHashFn(header,chainContext),
		GetHash:     func(n uint64) common.Hash { return header.Root }, // since this a one time thing, no point in adding complex functions to get the hash
		Origin:      userAddr,
		Coinbase:    userAddr,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    gasPrice,
	}
	chainConfig := params.AllSomaProtocolChanges
	vmconfig := vm.Config{}

	evm := vm.NewEVM(evmContext, statedb, chainConfig, vmconfig)

	sender := vm.AccountRef(userAddr)
	gas := uint64(1000000)
	value := new(big.Int).SetUint64(0x00)
	// CALL
	functionSig := "ActiveValidator(address)"
	log.Println("====== CALL =======", functionSig)
	encodedAddress := [32]byte{}
	copy(encodedAddress[12:], userAddr[:])
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()[:4]
	inputData := append(input[:], encodedAddress[:]...)
	ret, gas, vmerr := evm.Call(sender, contractAddress, inputData, gas, value)
	log.Printf("Result:\n%s\n", hex.Dump(ret))
	log.Println("User Address: ", userAddr)
	log.Println("Gas: ", gas)
	log.Println("Error: ", vmerr)

}

func callActiveValidators(userAddr common.Address, contractAddress common.Address, header *types.Header, db ethdb.Database) (bool, error) {
	// Byte encoding of booleans
	trueResult := "0000000000000000000000000000000000000000000000000000000000000001"

	// Signature of function being called defined by Soma interface
	functionSig := "ActiveValidator(address)"

	// Instantiate new state database
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(header.Root, sdb)

	log.Println("====== QUERY THE SMART CONTRACT =======", functionSig)

	gasPrice := new(big.Int).SetUint64(0x0)
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		//GetHash:     core.GetHashFn(header,chainContext),
		GetHash:     func(n uint64) common.Hash { return header.Root }, // since this a one time thing, no point in adding complex functions to get the hash
		Origin:      userAddr,
		Coinbase:    userAddr,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    gasPrice,
	}
	chainConfig := params.AllSomaProtocolChanges
	vmconfig := vm.Config{}

	evm := vm.NewEVM(evmContext, statedb, chainConfig, vmconfig)

	sender := vm.AccountRef(userAddr)

	log.Println("Contract exists: ", statedb.Exist(contractAddress))
	gas := uint64(1000000)
	value := new(big.Int).SetUint64(0x00)

	// CALL
	log.Println("====== CALL =======", functionSig)
	encodedAddress := [32]byte{}
	copy(encodedAddress[12:], userAddr[:])
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()[:4]
	inputData := append(input[:], encodedAddress[:]...)
	ret, gas, vmerr := evm.Call(sender, contractAddress, inputData, gas, value)
	if vmerr != nil {
		return false, vmerr
	}
	// Parse result
	// result := strings.Compare(hex.EncodeToString(ret), trueResult)
	log.Printf("Comparison:\n%s\n%s ", hex.EncodeToString(ret), trueResult)
	if hex.EncodeToString(ret) == trueResult {
		return true, nil

	} else {
		return false, nil
	}

}

func printDebug(funcName string, chain consensus.ChainReader, header *types.Header) {
	printHeader := func(h *types.Header) {
		log.Printf("Header argument:\n\tnumber: %#v\n\tHash: 0x%x\n\tState Root: 0x%x\n\tParentHash: 0x%x\n", h.Number, h.Hash().Bytes(), h.Root.Bytes(), h.ParentHash.Bytes())
	}
	log.Printf("%s: =========================================================\n", funcName)
	// golog.Printf("%#v\n", chain)
	printHeader(header)
	printHeader(chain.CurrentHeader())
	//golog.Printf("%#v", chain.GetBlock(chain.GetHeaderByNumber(0).Hash(), 0))
	log.Printf("=========================================================\n")
}

func printDB(sdb state.Database) {
	log.Print("\n\n\n>>>>>>>>>>>>>>>>>>>>>>>>>> [START] printDB()")
	log.Print("Trie Nodes")
	for idx, node := range sdb.TrieDB().Nodes() {
		log.Print("\n\t====================================================================\n\t===================================================================\n")
		val, err := sdb.TrieDB().Node(node)
		if err != nil {
			log.Print("ERROR:", err)
		}
		var decodedValue [][]byte
		err = rlp.DecodeBytes(val, &decodedValue)
		if err != nil {
			log.Print("ERROR:", err)
		}
		log.Printf("node[%d]:\n", idx)
		log.Printf("\tkey:\t0x%x\n", node)
		log.Printf("\tvalue bytes:\t0x%x\n", val)
		for _, decodedProp := range decodedValue {
			log.Printf("\t\tdecoded prop:\t0x%x\n", decodedProp)
		}

		if len(decodedValue) != 0 {
			h := common.BytesToHash(decodedValue[0])
			log.Printf("\n\thash of address used as key in trie:\t0x%x\n", h)

			var acc state.Account
			err = rlp.DecodeBytes(decodedValue[1], &acc)
			if err != nil {
				log.Print("ERROR:", err)
			}
			log.Printf("\n\taccount of user form trie:\t%#v\n", acc)
			//log.Printf("node[%d]:\t%x\t%#v\n\t%#v\n%#v\n", idx, node, val, a, b)
		} else {
			log.Printf("\n\tunknown decoded value:\t%#v", decodedValue)
		}
	}
	log.Print("<<<<<<<<<<<<<<<<<<<<<<<<<< [END] printDB()\n\n\n")
}
