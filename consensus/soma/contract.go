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
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

func deployContract(bytecodeStr string, statedb *state.StateDB) (common.Hash, error) {
	contractBytecode := common.Hex2Bytes(bytecodeStr[2:]) // [2:] removes 0x

	stateRoot := common.Hash{}
	log.Printf("State root: 0x%x\n", stateRoot)
	/*
		type Account struct {
			Nonce    uint64
			Balance  *big.Int
			Root     common.Hash // merkle root of the storage trie
			CodeHash []byte
		}
	*/

	userAddr := common.Address{}
	/*
		statedb, err := state.New(common.Hash{}, sdb)
		if err != nil {
			log.Printf("ERROR starting statedb! <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<,")
		}
	*/

	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		//GetHash:     core.GetHashFn(header,chainContext),
		GetHash:     func(n uint64) common.Hash { return stateRoot },
		Origin:      userAddr,
		Coinbase:    userAddr,
		BlockNumber: new(big.Int).SetUint64(0x00),
		Time:        new(big.Int).SetUint64(0x01),
		GasLimit:    uint64(0x0f4240),
		Difficulty:  new(big.Int).SetUint64(0x0100),
		GasPrice:    new(big.Int).SetUint64(0x3b9aca00),
	}
	chainConfig := params.AllSomaProtocolChanges
	vmconfig := vm.Config{}
	/*
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

	// commit makes current state saved into DB
	//root, err := statedb.Commit(true)
	//log.Printf("Trie root: 0x%x\n", root)

	//sdb := statedb.Database() //state.NewDatabase(db)
	//printDB(sdb)

	return common.Hash{}, nil
	//return root, err
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
