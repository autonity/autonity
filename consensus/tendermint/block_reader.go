package tendermint

import (
	"fmt"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/state"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/ethdb"
)

type BlockReader struct {
	db      ethdb.Database
	statedb state.Database
}

// NewBlockReader constructs a new BlockReader with the db and state. The state
// instance must be the same state instance that is used to construct the
// core.BlockChain as the second parameter to NewBlockChainWithState. This
func NewBlockReader(db ethdb.Database, state state.Database) *BlockReader {
	return &BlockReader{
		db:      db,
		statedb: state,
	}
}

func (l *BlockReader) LatestBlock() (*types.Block, error) {
	hash := rawdb.ReadHeadBlockHash(l.db)
	if hash == (common.Hash{}) {
		return nil, fmt.Errorf("empty database")
	}

	number := rawdb.ReadHeaderNumber(l.db, hash)
	if number == nil {
		return nil, fmt.Errorf("failed to find number for block Hash %s", hash.String())
	}

	block := rawdb.ReadBlock(l.db, hash, *number)
	if block == nil {
		return nil, fmt.Errorf("failed to read block content for block number %d with Hash %s", *number, hash.String())
	}

	_, err := l.statedb.OpenTrie(block.Root())
	if err != nil {
		return nil, fmt.Errorf("missing state for block number %d with Hash %s err: %v", *number, hash.String(), err)
	}
	return block, nil
}

func (l *BlockReader) BlockByHeight(height uint64) (*types.Block, error) {
	hash := rawdb.ReadCanonicalHash(l.db, height)
	if hash == (common.Hash{}) {
		return nil, fmt.Errorf("no block for number %d", height)
	}
	block := rawdb.ReadBlock(l.db, hash, height)
	if block == nil {
		panic(fmt.Errorf("missing block in db at number %d but found Hash %v", height, hash.String()))
	}

	_, err := l.statedb.OpenTrie(block.Root())
	if err != nil {
		return nil, fmt.Errorf("missing state for block number %d with Hash %s err: %v", height, hash.String(), err)
	}
	return block, nil
}

func (l *BlockReader) BlockState(stateRoot common.Hash) (*state.StateDB, error) {
	// state.New has started takig a snapshot.Tree but it seems to be only for
	// performance, see - https://github.com/ethereum/go-ethereum/pull/20152
	return state.New(stateRoot, l.statedb, nil)
}
