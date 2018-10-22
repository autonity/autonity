package soma

import (
	"encoding/hex"
	golog "log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

func deployContract(bytecodeStr string, userAddr common.Address, header *types.Header, statedb *state.StateDB) (common.Address, error) {
	contractBytecode := common.Hex2Bytes(bytecodeStr[2:]) // [2:] removes 0x

	// Initialise new Ethereum Virtual Machine
	gasPrice := new(big.Int).SetUint64(0x0)
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
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
	data := contractBytecode
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Soma validator governance contract
	_, contractAddress, gas, vmerr := evm.Create(sender, data, gas, value)

	if vmerr != nil {
		return contractAddress, vmerr
	}
	log.Info("Deployed Soma Governance Contract", "Address", contractAddress.String())

	return contractAddress, nil
}

func callActiveValidators(userAddr common.Address, contractAddress common.Address, header *types.Header, db ethdb.Database, chain consensus.ChainReader) (bool, error) {
	// Byte encoding of booleans
	expectedResult := "0000000000000000000000000000000000000000000000000000000000000001"

	// Signature of function being called defined by Soma interface
	functionSig := "ActiveValidator(address)"

	// Instantiate new state database
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(header.Root, sdb)

	// Initialise a new Ethereum Virtual Machine
	gasPrice := new(big.Int).SetUint64(0x0)

	// Implement my own stuff
	GetHashFn := func(ref *types.Header, chain consensus.ChainReader) func(n uint64) common.Hash {
		var cache map[uint64]common.Hash

		return func(n uint64) common.Hash {
			// If there's no hash cache yet, make one
			if cache == nil {
				cache = map[uint64]common.Hash{
					ref.Number.Uint64() - 1: ref.ParentHash,
				}
			}
			// Try to fulfill the request from the cache
			if hash, ok := cache[n]; ok {
				return hash
			}
			// Not cached, iterate the blocks and cache the hashes
			for header := chain.GetHeader(ref.ParentHash, ref.Number.Uint64()-1); header != nil; header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) {
				cache[header.Number.Uint64()-1] = header.ParentHash
				if n == header.Number.Uint64()-1 {
					return header.ParentHash
				}
			}
			return common.Hash{}
		}
	}

	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     GetHashFn(header, chain),
		// GetHash:     func(n uint64) common.Hash { return header.Root },
		Origin:      userAddr,
		Coinbase:    userAddr,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    gasPrice,
	}
	chainConfig := params.AllSomaProtocolChanges
	sender := vm.AccountRef(userAddr)
	gas := uint64(0xFFFFFFFF)
	// value := new(big.Int).SetUint64(0x00)
	vmconfig := vm.Config{}

	evm := vm.NewEVM(evmContext, statedb, chainConfig, vmconfig)

	// Pad address for ABI encoding
	encodedAddress := [32]byte{}
	copy(encodedAddress[12:], userAddr[:])
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()[:4]
	inputData := append(input[:], encodedAddress[:]...)

	// Call ActiveValidators()
	ret, gas, vmerr := evm.StaticCall(sender, contractAddress, inputData, gas)
	if vmerr != nil {
		return false, vmerr
	}

	// Check result
	if hex.EncodeToString(ret) == expectedResult {
		return true, nil

	} else {
		return false, nil
	}

}

func printDebug(funcName string, chain consensus.ChainReader, header *types.Header) {
	printHeader := func(h *types.Header) {
		golog.Printf("Header argument:\n\tnumber: %#v\n\tHash: 0x%x\n\tState Root: 0x%x\n\tParentHash: 0x%x\n", h.Number, h.Hash().Bytes(), h.Root.Bytes(), h.ParentHash.Bytes())
	}
	golog.Printf("%s: =========================================================\n", funcName)
	// gogolog.Printf("%#v\n", chain)
	printHeader(header)
	printHeader(chain.CurrentHeader())
	//gogolog.Printf("%#v", chain.GetBlock(chain.GetHeaderByNumber(0).Hash(), 0))
	golog.Printf("=========================================================\n")
}

func printDB(sdb state.Database) {
	golog.Print("\n\n\n>>>>>>>>>>>>>>>>>>>>>>>>>> [START] printDB()")
	golog.Print("Trie Nodes")
	for idx, node := range sdb.TrieDB().Nodes() {
		golog.Print("\n\t====================================================================\n\t===================================================================\n")
		val, err := sdb.TrieDB().Node(node)
		if err != nil {
			golog.Print("ERROR:", err)
		}
		var decodedValue [][]byte
		err = rlp.DecodeBytes(val, &decodedValue)
		if err != nil {
			golog.Print("ERROR:", err)
		}
		golog.Printf("node[%d]:\n", idx)
		golog.Printf("\tkey:\t0x%x\n", node)
		golog.Printf("\tvalue bytes:\t0x%x\n", val)
		for _, decodedProp := range decodedValue {
			golog.Printf("\t\tdecoded prop:\t0x%x\n", decodedProp)
		}

		if len(decodedValue) != 0 {
			h := common.BytesToHash(decodedValue[0])
			golog.Printf("\n\thash of address used as key in trie:\t0x%x\n", h)

			var acc state.Account
			err = rlp.DecodeBytes(decodedValue[1], &acc)
			if err != nil {
				golog.Print("ERROR:", err)
			}
			golog.Printf("\n\taccount of user form trie:\t%#v\n", acc)
			//golog.Printf("node[%d]:\t%x\t%#v\n\t%#v\n%#v\n", idx, node, val, a, b)
		} else {
			golog.Printf("\n\tunknown decoded value:\t%#v", decodedValue)
		}
	}
	golog.Print("<<<<<<<<<<<<<<<<<<<<<<<<<< [END] printDB()\n\n\n")
}
