package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/ethdb"
	lru "github.com/hashicorp/golang-lru"
)

type HeaderGetter interface {
	GetHeader(hash common.Hash, number uint64) *types.Header
}

func NewHeaderGetter(chainDb ethdb.Database) (*headerGetter, error) {
	headerCache, err := lru.New(headerCacheLimit)
	if err != nil {
		return nil, err
	}
	numberCache, err := lru.New(numberCacheLimit)
	if err != nil {
		return nil, err
	}
	return &headerGetter{
		headerCache: headerCache,
		numberCache: numberCache,
		chainDb:     chainDb,
	}, nil
}

type headerGetter struct {
	headerCache *lru.Cache // Cache for the most recent block headers
	numberCache *lru.Cache // Cache for the most recent block numbers
	chainDb     ethdb.Database
}

// GetHeader retrieves a block header from the database by hash and number,
// caching it if found.
func (hg *headerGetter) GetHeader(hash common.Hash, number uint64) *types.Header {
	// Short circuit if the header's already in the cache, retrieve otherwise
	if header, ok := hg.headerCache.Get(hash); ok {
		return header.(*types.Header)
	}
	header := rawdb.ReadHeader(hg.chainDb, hash, number)
	if header == nil {
		return nil
	}
	// Cache the found header for next time and return
	hg.headerCache.Add(hash, header)
	return header
}

// GetHeaderByHash retrieves a block header from the database by hash, caching it if
// found.
func (hg *headerGetter) GetHeaderByHash(hash common.Hash) *types.Header {
	number := hg.GetBlockNumber(hash)
	if number == nil {
		return nil
	}
	return hg.GetHeader(hash, *number)
}

// HasHeader checks if a block header is present in the database or not.
// In theory, if header is present in the database, all relative components
// like td and hash->number should be present too.
func (hg *headerGetter) HasHeader(hash common.Hash, number uint64) bool {
	if hg.numberCache.Contains(hash) || hg.headerCache.Contains(hash) {
		return true
	}
	return rawdb.HasHeader(hg.chainDb, hash, number)
}

// GetHeaderByNumber retrieves a block header from the database by number,
// caching it (associated with its hash) if found.
func (hg *headerGetter) GetHeaderByNumber(number uint64) *types.Header {
	hash := rawdb.ReadCanonicalHash(hg.chainDb, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return hg.GetHeader(hash, number)
}

// GetBlockNumber retrieves the block number belonging to the given hash
// from the cache or database
func (hg *headerGetter) GetBlockNumber(hash common.Hash) *uint64 {
	if cached, ok := hg.numberCache.Get(hash); ok {
		number := cached.(uint64)
		return &number
	}
	number := rawdb.ReadHeaderNumber(hg.chainDb, hash)
	if number != nil {
		hg.numberCache.Add(hash, *number)
	}
	return number
}

func (hg *headerGetter) AddToCache(hash common.Hash, number uint64, header *types.Header) {
	hg.headerCache.Add(hash, header)
	hg.numberCache.Add(hash, number)
}
