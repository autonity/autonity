// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/internal/ethapi"
	"github.com/autonity/autonity/rlp"
	"github.com/autonity/autonity/rpc"
	"github.com/autonity/autonity/trie"
)

// PublicEthereumAPI provides an API to access Ethereum full node-related
// information.
type PublicEthereumAPI struct {
	e *Ethereum
}

// NewPublicEthereumAPI creates a new Ethereum protocol API for full nodes.
func NewPublicEthereumAPI(e *Ethereum) *PublicEthereumAPI {
	return &PublicEthereumAPI{e}
}

// Etherbase is the address that mining rewards will be send to
func (api *PublicEthereumAPI) Etherbase() (common.Address, error) {
	return api.e.Etherbase()
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (api *PublicEthereumAPI) Coinbase() (common.Address, error) {
	return api.Etherbase()
}

// Hashrate returns the POW hashrate
func (api *PublicEthereumAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(api.e.Miner().Hashrate())
}

// PublicMinerAPI provides an API to control the miner.
// It offers only methods that operate on data that pose no security risk when it is publicly accessible.
type PublicMinerAPI struct {
	e *Ethereum
}

// NewPublicMinerAPI create a new PublicMinerAPI instance.
func NewPublicMinerAPI(e *Ethereum) *PublicMinerAPI {
	return &PublicMinerAPI{e}
}

// Mining returns an indication if this node is currently mining.
func (api *PublicMinerAPI) Mining() bool {
	return api.e.IsMining()
}

// PrivateMinerAPI provides private RPC methods to control the miner.
// These methods can be abused by external users and must be considered insecure for use by untrusted users.
type PrivateMinerAPI struct {
	e *Ethereum
}

// NewPrivateMinerAPI create a new RPC service which controls the miner of this node.
func NewPrivateMinerAPI(e *Ethereum) *PrivateMinerAPI {
	return &PrivateMinerAPI{e: e}
}

// Start starts the miner with the given number of threads. If threads is nil,
// the number of workers started is equal to the number of logical CPUs that are
// usable by this process. If mining is already running, this method adjust the
// number of threads allowed to use and updates the minimum price required by the
// transaction pool.
// NOTE: will start mining even if the node is out of sync with the chain head
func (api *PrivateMinerAPI) Start(threads *int) error {
	if threads == nil {
		return api.e.StartMining(runtime.NumCPU())
	}
	return api.e.StartMining(*threads)
}

// Stop terminates the miner, both at the consensus engine level as well as at
// the block creation level.
func (api *PrivateMinerAPI) Stop() {
	api.e.StopMining()
}

// SetExtra sets the extra data string that is included when this miner mines a block.
func (api *PrivateMinerAPI) SetExtra(extra string) (bool, error) {
	if err := api.e.Miner().SetExtra([]byte(extra)); err != nil {
		return false, err
	}
	return true, nil
}

// SetGasPrice sets the minimum accepted gas price for the miner.
func (api *PrivateMinerAPI) SetGasPrice(gasPrice hexutil.Big) bool {
	api.e.lock.Lock()
	api.e.gasPrice = (*big.Int)(&gasPrice)
	api.e.lock.Unlock()

	api.e.txPool.SetGasPrice((*big.Int)(&gasPrice))
	return true
}

// SetGasLimit sets the gaslimit to target towards during mining.
func (api *PrivateMinerAPI) SetGasLimit(gasLimit hexutil.Uint64) bool {
	api.e.Miner().SetGasCeil(uint64(gasLimit))
	return true
}

// SetRecommitInterval updates the interval for miner sealing work recommitting.
func (api *PrivateMinerAPI) SetRecommitInterval(interval int) {
	api.e.Miner().SetRecommitInterval(time.Duration(interval) * time.Millisecond)
}

// PrivateAdminAPI is the collection of Ethereum full node-related APIs
// exposed over the private admin endpoint.
type PrivateAdminAPI struct {
	eth *Ethereum
}

// NewPrivateAdminAPI creates a new API definition for the full node private
// admin methods of the Ethereum service.
func NewPrivateAdminAPI(eth *Ethereum) *PrivateAdminAPI {
	return &PrivateAdminAPI{eth: eth}
}

// ExportChain exports the current blockchain into a local file,
// or a range of blocks if first and last are non-nil
func (api *PrivateAdminAPI) ExportChain(file string, first *uint64, last *uint64) (bool, error) {
	if first == nil && last != nil {
		return false, errors.New("last cannot be specified without first")
	}
	if first != nil && last == nil {
		head := api.eth.BlockChain().CurrentHeader().Number.Uint64()
		last = &head
	}
	if _, err := os.Stat(file); err == nil {
		// File already exists. Allowing overwrite could be a DoS vector,
		// since the 'file' may point to arbitrary paths on the drive
		return false, errors.New("location would overwrite an existing file")
	}
	// Make sure we can create the file to export into
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer out.Close()

	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	// Export the blockchain
	if first != nil {
		if err := api.eth.BlockChain().ExportN(writer, *first, *last); err != nil {
			return false, err
		}
	} else if err := api.eth.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}

func hasAllBlocks(chain *core.BlockChain, bs []*types.Block) bool {
	for _, b := range bs {
		if !chain.HasBlock(b.Hash(), b.NumberU64()) {
			return false
		}
	}

	return true
}

// ImportChain imports a blockchain from a local file.
func (api *PrivateAdminAPI) ImportChain(file string) (bool, error) {
	// Make sure the can access the file to import
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()

	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}

	// Run actual the import in pre-configured batches
	stream := rlp.NewStream(reader, 0)

	blocks, index := make([]*types.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		// Load a batch of blocks from the input file
		for len(blocks) < cap(blocks) {
			block := new(types.Block)
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}

		if hasAllBlocks(api.eth.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		// Import the batch and reset the buffer
		if _, err := api.eth.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}

// PublicDebugAPI is the collection of Ethereum full node APIs exposed
// over the public debugging endpoint.
type PublicDebugAPI struct {
	eth *Ethereum
}

// NewPublicDebugAPI creates a new API definition for the full node-
// related public debug methods of the Ethereum service.
func NewPublicDebugAPI(eth *Ethereum) *PublicDebugAPI {
	return &PublicDebugAPI{eth: eth}
}

// DumpBlock retrieves the entire state of the database at a given block.
func (api *PublicDebugAPI) DumpBlock(blockNr rpc.BlockNumber) (state.Dump, error) {
	opts := &state.DumpConfig{
		OnlyWithAddresses: true,
		Max:               AccountRangeMaxResults, // Sanity limit over RPC
	}
	if blockNr == rpc.PendingBlockNumber {
		// If we're dumping the pending state, we need to request
		// both the pending block as well as the pending state from
		// the miner and operate on those
		_, stateDb := api.eth.miner.Pending()
		return stateDb.RawDump(opts), nil
	}
	var block *types.Block
	if blockNr == rpc.LatestBlockNumber {
		block = api.eth.blockchain.CurrentBlock()
	} else {
		block = api.eth.blockchain.GetBlockByNumber(uint64(blockNr))
	}
	if block == nil {
		return state.Dump{}, fmt.Errorf("block #%d not found", blockNr)
	}
	stateDb, err := api.eth.BlockChain().StateAt(block.Root())
	if err != nil {
		return state.Dump{}, err
	}
	return stateDb.RawDump(opts), nil
}

// PrivateDebugAPI is the collection of Ethereum full node APIs exposed over
// the private debugging endpoint.
type PrivateDebugAPI struct {
	eth *Ethereum
}

// NewPrivateDebugAPI creates a new API definition for the full node-related
// private debug methods of the Ethereum service.
func NewPrivateDebugAPI(eth *Ethereum) *PrivateDebugAPI {
	return &PrivateDebugAPI{eth: eth}
}

// Preimage is a debug API function that returns the preimage for a sha3 hash, if known.
func (api *PrivateDebugAPI) Preimage(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	if preimage := rawdb.ReadPreimage(api.eth.ChainDb(), hash); preimage != nil {
		return preimage, nil
	}
	return nil, errors.New("unknown preimage")
}

// BadBlockArgs represents the entries in the list returned when bad blocks are queried.
type BadBlockArgs struct {
	Hash  common.Hash            `json:"hash"`
	Block map[string]interface{} `json:"block"`
	RLP   string                 `json:"rlp"`
}

// GetBadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
// and returns them as a JSON list of block-hashes
func (api *PrivateDebugAPI) GetBadBlocks(ctx context.Context) ([]*BadBlockArgs, error) {
	var (
		err     error
		blocks  = rawdb.ReadAllBadBlocks(api.eth.chainDb)
		results = make([]*BadBlockArgs, 0, len(blocks))
	)
	for _, block := range blocks {
		var (
			blockRlp  string
			blockJSON map[string]interface{}
		)
		if rlpBytes, err := rlp.EncodeToBytes(block); err != nil {
			blockRlp = err.Error() // Hacky, but hey, it works
		} else {
			blockRlp = fmt.Sprintf("0x%x", rlpBytes)
		}
		if blockJSON, err = ethapi.RPCMarshalBlock(block, true, true, api.eth.APIBackend.ChainConfig()); err != nil {
			blockJSON = map[string]interface{}{"error": err.Error()}
		}
		results = append(results, &BadBlockArgs{
			Hash:  block.Hash(),
			RLP:   blockRlp,
			Block: blockJSON,
		})
	}
	return results, nil
}

// AccountRangeMaxResults is the maximum number of results to be returned per call
const AccountRangeMaxResults = 256

// AccountRange enumerates all accounts in the given block and start point in paging request
func (api *PublicDebugAPI) AccountRange(blockNrOrHash rpc.BlockNumberOrHash, start []byte, maxResults int, nocode, nostorage, incompletes bool) (state.IteratorDump, error) {
	var stateDb *state.StateDB
	var err error

	if number, ok := blockNrOrHash.Number(); ok {
		if number == rpc.PendingBlockNumber {
			// If we're dumping the pending state, we need to request
			// both the pending block as well as the pending state from
			// the miner and operate on those
			_, stateDb = api.eth.miner.Pending()
		} else {
			var block *types.Block
			if number == rpc.LatestBlockNumber {
				block = api.eth.blockchain.CurrentBlock()
			} else {
				block = api.eth.blockchain.GetBlockByNumber(uint64(number))
			}
			if block == nil {
				return state.IteratorDump{}, fmt.Errorf("block #%d not found", number)
			}
			stateDb, err = api.eth.BlockChain().StateAt(block.Root())
			if err != nil {
				return state.IteratorDump{}, err
			}
		}
	} else if hash, ok := blockNrOrHash.Hash(); ok {
		block := api.eth.blockchain.GetBlockByHash(hash)
		if block == nil {
			return state.IteratorDump{}, fmt.Errorf("block %s not found", hash.Hex())
		}
		stateDb, err = api.eth.BlockChain().StateAt(block.Root())
		if err != nil {
			return state.IteratorDump{}, err
		}
	} else {
		return state.IteratorDump{}, errors.New("either block number or block hash must be specified")
	}

	opts := &state.DumpConfig{
		SkipCode:          nocode,
		SkipStorage:       nostorage,
		OnlyWithAddresses: !incompletes,
		Start:             start,
		Max:               uint64(maxResults),
	}
	if maxResults > AccountRangeMaxResults || maxResults <= 0 {
		opts.Max = AccountRangeMaxResults
	}
	return stateDb.IteratorDump(opts), nil
}

// StorageRangeResult is the result of a debug_storageRangeAt API call.
type StorageRangeResult struct {
	Storage storageMap   `json:"storage"`
	NextKey *common.Hash `json:"nextKey"` // nil if Storage includes the last key in the trie.
}

type storageMap map[common.Hash]storageEntry

type storageEntry struct {
	Key   *common.Hash `json:"key"`
	Value common.Hash  `json:"value"`
}

// StorageRangeAt returns the storage at the given block height and transaction index.
func (api *PrivateDebugAPI) StorageRangeAt(blockHash common.Hash, txIndex int, contractAddress common.Address, keyStart hexutil.Bytes, maxResult int) (StorageRangeResult, error) {
	// Retrieve the block
	block := api.eth.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return StorageRangeResult{}, fmt.Errorf("block %#x not found", blockHash)
	}
	_, _, statedb, err := api.eth.stateAtTransaction(block, txIndex, 0)
	if err != nil {
		return StorageRangeResult{}, err
	}
	st := statedb.StorageTrie(contractAddress)
	if st == nil {
		return StorageRangeResult{}, fmt.Errorf("account %x doesn't exist", contractAddress)
	}
	return storageRangeAt(st, keyStart, maxResult)
}

func storageRangeAt(st state.Trie, start []byte, maxResult int) (StorageRangeResult, error) {
	it := trie.NewIterator(st.NodeIterator(start))
	result := StorageRangeResult{Storage: storageMap{}}
	for i := 0; i < maxResult && it.Next(); i++ {
		_, content, _, err := rlp.Split(it.Value)
		if err != nil {
			return StorageRangeResult{}, err
		}
		e := storageEntry{Value: common.BytesToHash(content)}
		if preimage := st.GetKey(it.Key); preimage != nil {
			preimage := common.BytesToHash(preimage)
			e.Key = &preimage
		}
		result.Storage[common.BytesToHash(it.Key)] = e
	}
	// Add the 'next key' so clients can continue downloading.
	if it.Next() {
		next := common.BytesToHash(it.Key)
		result.NextKey = &next
	}
	return result, nil
}

// GetModifiedAccountsByNumber returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateDebugAPI) GetModifiedAccountsByNumber(startNum uint64, endNum *uint64) ([]common.Address, error) {
	var startBlock, endBlock *types.Block

	startBlock = api.eth.blockchain.GetBlockByNumber(startNum)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startNum)
	}

	if endNum == nil {
		endBlock = startBlock
		startBlock = api.eth.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.eth.blockchain.GetBlockByNumber(*endNum)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %d not found", *endNum)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

// GetModifiedAccountsByHash returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateDebugAPI) GetModifiedAccountsByHash(startHash common.Hash, endHash *common.Hash) ([]common.Address, error) {
	var startBlock, endBlock *types.Block
	startBlock = api.eth.blockchain.GetBlockByHash(startHash)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startHash)
	}

	if endHash == nil {
		endBlock = startBlock
		startBlock = api.eth.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.eth.blockchain.GetBlockByHash(*endHash)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %x not found", *endHash)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

func (api *PrivateDebugAPI) getModifiedAccounts(startBlock, endBlock *types.Block) ([]common.Address, error) {
	if startBlock.Number().Uint64() >= endBlock.Number().Uint64() {
		return nil, fmt.Errorf("start block height (%d) must be less than end block height (%d)", startBlock.Number().Uint64(), endBlock.Number().Uint64())
	}
	triedb := api.eth.BlockChain().StateCache().TrieDB()

	oldTrie, err := trie.NewSecure(startBlock.Root(), triedb)
	if err != nil {
		return nil, err
	}
	newTrie, err := trie.NewSecure(endBlock.Root(), triedb)
	if err != nil {
		return nil, err
	}
	diff, _ := trie.NewDifferenceIterator(oldTrie.NodeIterator([]byte{}), newTrie.NodeIterator([]byte{}))
	iter := trie.NewIterator(diff)

	var dirty []common.Address
	for iter.Next() {
		key := newTrie.GetKey(iter.Key)
		if key == nil {
			return nil, fmt.Errorf("no preimage found for hash %x", iter.Key)
		}
		dirty = append(dirty, common.BytesToAddress(key))
	}
	return dirty, nil
}

// GetAccessibleState returns the first number where the node has accessible
// state on disk. Note this being the post-state of that block and the pre-state
// of the next block.
// The (from, to) parameters are the sequence of blocks to search, which can go
// either forwards or backwards
func (api *PrivateDebugAPI) GetAccessibleState(from, to rpc.BlockNumber) (uint64, error) {
	db := api.eth.ChainDb()
	var pivot uint64
	if p := rawdb.ReadLastPivotNumber(db); p != nil {
		pivot = *p
		log.Info("Found fast-sync pivot marker", "number", pivot)
	}
	var resolveNum = func(num rpc.BlockNumber) (uint64, error) {
		// We don't have state for pending (-2), so treat it as latest
		if num.Int64() < 0 {
			block := api.eth.blockchain.CurrentBlock()
			if block == nil {
				return 0, fmt.Errorf("current block missing")
			}
			return block.NumberU64(), nil
		}
		return uint64(num.Int64()), nil
	}
	var (
		start   uint64
		end     uint64
		delta   = int64(1)
		lastLog time.Time
		err     error
	)
	if start, err = resolveNum(from); err != nil {
		return 0, err
	}
	if end, err = resolveNum(to); err != nil {
		return 0, err
	}
	if start == end {
		return 0, fmt.Errorf("from and to needs to be different")
	}
	if start > end {
		delta = -1
	}
	for i := int64(start); i != int64(end); i += delta {
		if time.Since(lastLog) > 8*time.Second {
			log.Info("Finding roots", "from", start, "to", end, "at", i)
			lastLog = time.Now()
		}
		if i < int64(pivot) {
			continue
		}
		h := api.eth.BlockChain().GetHeaderByNumber(uint64(i))
		if h == nil {
			return 0, fmt.Errorf("missing header %d", i)
		}
		if ok, _ := api.eth.ChainDb().Has(h.Root[:]); ok {
			return uint64(i), nil
		}
	}
	return 0, fmt.Errorf("No state found")
}

// AutonityContractAPI implements rpc.Methods to expose view functions of the
// autonity contract through the rpc api. Note, although it looks like this
// struct would be better defined in the rpc package or in the autonity
// package, circular dependencies make it infeasible.
type AutonityContractAPI struct {
	calls map[string]reflect.Value
}

// NewAutonityContractAPI builds a map of function name to method representing
// the view functions of the autonity contract, the map is then used as the
// return value for AllMethods. The methods are dynamically generated from the
// autonity contract ABI and assume that they are called with an instance of
// AutonityContractAPI as their method receiver, even though the functions
// themselves make no use of the method receiver. This design is required to be
// able to fit into the current approach taken for registering rpc services.
// See rpc.Server.RegisterName().
func NewAutonityContractAPI(bc *core.BlockChain, ac *autonity.ProtocolContracts) *AutonityContractAPI {
	var contractABI = ac.ABI()
	var contractViewMethods = make(map[string]reflect.Value)

	for n, m := range contractABI.Methods {
		functionName := n
		// Only expose read-only functions.
		if m.IsConstant() {
			// The RPC service expect the first argument of an API method to be the receiver object.
			inArgs := []reflect.Type{reflect.TypeOf(&AutonityContractAPI{})}
			inArgs = append(inArgs, m.Inputs.Types()...)
			sig := reflect.FuncOf(inArgs, []reflect.Type{
				reflect.TypeOf((*interface{})(nil)).Elem(),
				reflect.TypeOf((*error)(nil)).Elem(),
			}, false)

			contractViewMethods[functionName] = reflect.MakeFunc(sig,
				func(args []reflect.Value) []reflect.Value {
					// makereturn converts the return types to reflect.Value
					makereturn := func(res interface{}, err error) []reflect.Value {
						return []reflect.Value{reflect.ValueOf(&res).Elem(), reflect.ValueOf(&err).Elem()}
					}
					stateDB, err := bc.State()
					if err != nil {
						return makereturn(nil, err)
					}
					var iargs []interface{}
					// args[0] is the reflect.Value of *AutonityContractAPI.
					for i, arg := range args[1:] {
						// If the argument is a pointer it is then an optional parameter for the RPC handler. The
						// json unmarshalling function set it to nil if the argument isn't set in the RPC call.
						// There are no optional parameters for the Autonity contract methods. Solidity doesn't
						// even support them and the packing function will crash if nil is passed.
						if arg.Kind() == reflect.Ptr && arg.IsNil() {
							return makereturn(nil, fmt.Errorf("missing value for required argument %d", i))
						}
						iargs = append(iargs, arg.Interface())
					}

					// Pack the arguments call the function and then unpack the result and return it.
					packedArgs, err := contractABI.Pack(functionName, iargs...)
					if err != nil {
						return makereturn(nil, err)
					}
					packedResult, _, err := ac.CallContractFunc(stateDB, bc.CurrentHeader(), params.AutonityContractAddress, packedArgs)
					if err != nil {
						return makereturn(nil, err)
					}
					result, err := contractABI.Unpack(functionName, packedResult)

					// If the result slice contains only one element then just return the element.
					if len(result) == 1 {
						return makereturn(result[0], err)
					}
					return makereturn(result, err)
				})
		}
	}
	return &AutonityContractAPI{calls: contractViewMethods}
}

func (a *AutonityContractAPI) AllMethods() map[string]reflect.Value {
	return a.calls
}
