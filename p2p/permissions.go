package p2p

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/params"
)

const (
	NODE_NAME_LENGTH    = 32
	PERMISSIONED_CONFIG = "../permissioned-nodes.json"
)

//
func queryDb(datadir string) (*state.StateDB, *types.Header) {
	log.Info("QueryDb")
	db, err := ethdb.NewLDBDatabase(datadir+"/geth/chaindata", 768, 512)
	if err != nil {
		log.Info("err", "err", err)
	}

	hash := rawdb.ReadHeadBlockHash(db)
	// golog.Printf("hash: \t%x\n", hash)

	number := rawdb.ReadHeaderNumber(db, hash)
	// golog.Printf("Thing: \t%v\n", number)

	header := rawdb.ReadHeader(db, hash, *number)
	// golog.Printf("Thing: \t%v\n", header)

	sdb := state.NewDatabase(db)
	statedb, err := state.New(header.Root, sdb)
	if err != nil {
		log.Info("err", "err", err)
	}
	return statedb, header
}

// check if a given node is permissioned to connect to the change
func isNodePermissioned(nodename string, currentNode string, datadir string, direction string) bool {

	var permissionedList []string
	nodes := parsePermissionedNodes(datadir)
	for _, v := range nodes {
		permissionedList = append(permissionedList, v.ID.String())
	}

	log.Debug("isNodePermissioned", "permissionedList", permissionedList)
	for _, v := range permissionedList {
		if v == nodename {
			log.Debug("isNodePermissioned", "connection", direction, "nodename", nodename[:NODE_NAME_LENGTH], "ALLOWED-BY", currentNode[:NODE_NAME_LENGTH])
			return true
		}
		log.Debug("isNodePermissioned", "connection", direction, "nodename", nodename[:NODE_NAME_LENGTH], "DENIED-BY", currentNode[:NODE_NAME_LENGTH])
	}
	log.Debug("isNodePermissioned", "connection", direction, "nodename", nodename[:NODE_NAME_LENGTH], "DENIED-BY", currentNode[:NODE_NAME_LENGTH])
	return false
}

// contractCallAddress queries the Soma contract, for any functions that take and address as argument.
// Returns true/false if the the address is an active validator and false if not.
func callGlienicke(functionSig string, userAddr common.Address, contractAddress common.Address, header *types.Header, statedb *state.StateDB) bool {
	sender := vm.AccountRef(userAddr)
	gas := uint64(0xFFFFFFFF)
	evm := getEVM(header, userAddr, userAddr, statedb)

	// Pad address for ABI encoding
	// encodedAddress := [32]byte{}
	// copy(encodedAddress[12:], userAddr[:])
	input := crypto.Keccak256Hash([]byte(functionSig)).Bytes()
	// inputData := append(input[:], encodedAddress[:]...)

	// Call ActiveValidators()
	ret, gas, vmerr := evm.StaticCall(sender, contractAddress, input, gas)
	if len(ret) == 0 {
		log.Info("contractCallAddress(): No return value", "Block", header.Number.Int64())
		return false
	}

	if vmerr != nil {
		return false
	}

	const def = `[{ "name" : "method", "outputs": [{ "type": "bool" }] }]`
	funcAbi, err := abi.JSON(strings.NewReader(def))
	if err != nil {
		return false
	}

	var output bool
	err = funcAbi.Unpack(&output, "method", ret)
	if err != nil {
		return false
	}

	return output
}

// getEVM Instantiates a new EVM object which is required when creating or calling a deployed contract
func getEVM(header *types.Header, coinbase, origin common.Address, statedb *state.StateDB) *vm.EVM {
	gasPrice := new(big.Int).SetUint64(0x0)
	evmContext := vm.Context{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     func(n uint64) common.Hash { return header.Root }, // since this a one time thing, no point in adding complex functions to get the hash
		Origin:      origin,
		Coinbase:    coinbase,
		BlockNumber: header.Number,
		Time:        header.Time,
		GasLimit:    header.GasLimit,
		Difficulty:  header.Difficulty,
		GasPrice:    gasPrice,
	}

	chainConfig := params.AllCliqueProtocolChanges
	vmconfig := vm.Config{}
	evm := vm.NewEVM(evmContext, statedb, chainConfig, vmconfig)
	return evm
}

// func isNodePermissioned(nodename string, currentNode string, datadir string, direction string) bool {

// 	var permissionedList []string
// 	nodes := parsePermissionedNodes(datadir)
// 	for _, v := range nodes {
// 		permissionedList = append(permissionedList, v.ID.String())
// 	}

// 	log.Debug("isNodePermissioned", "permissionedList", permissionedList)
// 	for _, v := range permissionedList {
// 		if v == nodename {
// 			log.Debug("isNodePermissioned", "connection", direction, "nodename", nodename[:NODE_NAME_LENGTH], "ALLOWED-BY", currentNode[:NODE_NAME_LENGTH])
// 			return true
// 		}
// 		log.Debug("isNodePermissioned", "connection", direction, "nodename", nodename[:NODE_NAME_LENGTH], "DENIED-BY", currentNode[:NODE_NAME_LENGTH])
// 	}
// 	log.Debug("isNodePermissioned", "connection", direction, "nodename", nodename[:NODE_NAME_LENGTH], "DENIED-BY", currentNode[:NODE_NAME_LENGTH])
// 	return false
// }

//this is a shameless copy from the config.go. It is a duplication of the code
//for the timebeing to allow reload of the permissioned nodes while the server is running

func parsePermissionedNodes(DataDir string) []*discover.Node {

	log.Debug("parsePermissionedNodes", "DataDir", DataDir, "file", PERMISSIONED_CONFIG)

	path := filepath.Join(DataDir, PERMISSIONED_CONFIG)
	if _, err := os.Stat(path); err != nil {
		log.Error("Read Error for permissioned-nodes.json file. This is because 'permissioned' flag is specified but no permissioned-nodes.json file is present.", "err", err)
		return nil
	}
	// Load the nodes from the config file
	blob, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("parsePermissionedNodes: Failed to access nodes", "err", err)
		return nil
	}

	nodelist := []string{}
	if err := json.Unmarshal(blob, &nodelist); err != nil {
		log.Error("parsePermissionedNodes: Failed to load nodes", "err", err)
		return nil
	}
	// Interpret the list as a discovery node array
	var nodes []*discover.Node
	for _, url := range nodelist {
		if url == "" {
			log.Error("parsePermissionedNodes: Node URL blank")
			continue
		}
		node, err := discover.ParseNode(url)
		if err != nil {
			log.Error("parsePermissionedNodes: Node URL", "url", url, "err", err)
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}
