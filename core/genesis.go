// Copyright 2014 The go-ethereum Authors
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

package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/autonity/autonity/accounts/keystore"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/ethdb"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"
	"github.com/autonity/autonity/trie"
)

//go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config     *params.ChainConfig `json:"config"`
	Nonce      uint64              `json:"nonce"`
	Timestamp  uint64              `json:"timestamp"`
	ExtraData  []byte              `json:"extraData"`
	GasLimit   uint64              `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int            `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash         `json:"mixHash"`
	Coinbase   common.Address      `json:"coinbase"`
	Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`

	BaseFee *big.Int `json:"baseFeePerGas"`
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount

func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

func (ga *GenesisAlloc) ToGenesisBonds() autonity.GenesisBonds {
	ret := make([]autonity.GenesisBond, 0, len(*ga))
	for addr, alloc := range *ga {
		delegations := make([]autonity.Delegation, 0)
		for validator, amount := range alloc.Bonds {
			delegations = append(delegations, autonity.Delegation{Validator: validator, Amount: amount})
		}
		slices.SortFunc(delegations, func(a, b autonity.Delegation) int {
			return strings.Compare(a.Validator.String(), b.Validator.String())
		})
		ret = append(ret, autonity.GenesisBond{
			Staker:        addr,
			NewtonBalance: alloc.NewtonBalance,
			Bonds:         delegations,
		})
	}
	slices.SortFunc(ret, func(a, b autonity.GenesisBond) int {
		return strings.Compare(a.Staker.String(), b.Staker.String())
	})
	return ret
}

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code          []byte                      `json:"code,omitempty"`
	Storage       map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance       *big.Int                    `json:"balance" gencodec:"required"`
	NewtonBalance *big.Int                    `json:"newtonBalance"`
	// validator address to amount bond to this validator
	Bonds      map[common.Address]*big.Int `json:"bonds"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Nonce      math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasLimit   math.HexOrDecimal64
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Difficulty *math.HexOrDecimal256
	BaseFee    *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code          hexutil.Bytes
	Balance       *math.HexOrDecimal256
	NewtonBalance *math.HexOrDecimal256
	Bonds         map[common.Address]*math.HexOrDecimal256
	Nonce         math.HexOrDecimal64
	Storage       map[storageJSON]storageJSON
	PrivateKey    hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database contains incompatible genesis (have %x, new %x)", e.Stored, e.New)
}

// SetupGenesisBlock writes or updates the genesis block in db.
// The behavior of it is:
//
//	                       genesis == nil         genesis != nil
//	                  +------------------------------------------
//	db has no genesis |  Return An Error      |  apply genesis to db
//	db has genesis    |  Use genesis from DB  |  apply genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(db ethdb.Database, genesis *Genesis) (*params.ChainConfig, common.Hash, error) {
	return SetupGenesisBlockWithOverride(db, genesis, nil, nil)
}

func SetupGenesisBlockWithOverride(db ethdb.Database, genesis *Genesis, overrideArrowGlacier, overrideTerminalTotalDifficulty *big.Int) (*params.ChainConfig, common.Hash, error) {
	if genesis != nil && genesis.Config == nil {
		return params.AllEthashProtocolChanges, common.Hash{}, errGenesisNoConfig
	}
	// Just commit the new block if there is no stored genesis block.
	stored := rawdb.ReadCanonicalHash(db, 0)
	if (stored == common.Hash{}) {
		if genesis == nil {
			log.Info("Writing default main-net genesis block")
			genesis = DefaultMainnetGenesisBlock()
		} else {
			log.Info("Writing custom genesis block")
		}
		block, err := genesis.Commit(db)
		if err != nil {
			return genesis.Config, common.Hash{}, err
		}
		return genesis.Config, block.Hash(), nil
	}
	// We have the genesis block in database(perhaps in ancient database)
	// but the corresponding state is missing.
	header := rawdb.ReadHeader(db, stored, 0)
	if _, err := state.New(header.Root, state.NewDatabaseWithConfig(db, nil), nil); err != nil {
		if genesis == nil {
			genesis = DefaultGenesisBlock()
		}
		// Ensure the stored genesis matches with the given one.
		b, err := genesis.ToBlock(nil)
		if err != nil {
			return nil, common.Hash{}, err
		}
		hash := b.Hash()
		if hash != stored {
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
		block, err := genesis.Commit(db)
		if err != nil {
			return genesis.Config, hash, err
		}
		return genesis.Config, block.Hash(), nil
	}
	// Check whether the genesis block is already written.
	if genesis != nil {
		b, err := genesis.ToBlock(nil)
		if err != nil {
			return nil, common.Hash{}, err
		}
		hash := b.Hash()
		if hash != stored {
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
	}
	// Get the existing chain configuration.
	newcfg := genesis.configOrDefault(stored)
	if overrideArrowGlacier != nil {
		newcfg.ArrowGlacierBlock = overrideArrowGlacier
	}
	if overrideTerminalTotalDifficulty != nil {
		newcfg.TerminalTotalDifficulty = overrideTerminalTotalDifficulty
	}
	if err := newcfg.CheckConfigForkOrder(); err != nil {
		return newcfg, common.Hash{}, err
	}
	storedcfg := rawdb.ReadChainConfig(db, stored)
	if storedcfg == nil {
		log.Warn("Found genesis block without chain config")
		rawdb.WriteChainConfig(db, stored, newcfg)
		return newcfg, stored, nil
	}
	// Special case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	if genesis == nil && stored != params.MainnetGenesisHash {
		return storedcfg, stored, nil
	}
	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := rawdb.ReadHeaderNumber(db, rawdb.ReadHeadHeaderHash(db))
	if height == nil {
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}
	compatErr := storedcfg.CheckCompatible(newcfg, *height)
	if compatErr != nil && *height != 0 && compatErr.RewindTo != 0 {
		return newcfg, stored, compatErr
	}
	rawdb.WriteChainConfig(db, stored, newcfg)
	return newcfg, stored, nil
}

func (g *Genesis) configOrDefault(ghash common.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	case ghash == params.MainnetGenesisHash:
		return params.MainnetChainConfig
	case ghash == params.RopstenGenesisHash:
		return params.RopstenChainConfig
	case ghash == params.SepoliaGenesisHash:
		return params.SepoliaChainConfig
	case ghash == params.RinkebyGenesisHash:
		return params.RinkebyChainConfig
	case ghash == params.GoerliGenesisHash:
		return params.GoerliChainConfig
	default:
		return params.AllEthashProtocolChanges
	}
}

// ToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func (g *Genesis) ToBlock(db ethdb.Database) (*types.Block, error) {
	g.setDefaultHardforks()
	g.Config.SetDefaults()
	if err := g.Config.Prepare(); err != nil {
		return nil, err
	}

	if g.Difficulty == nil {
		g.Difficulty = params.GenesisDifficulty
	}
	if g.Difficulty.Cmp(big.NewInt(0)) != 0 {
		return nil, fmt.Errorf("autonity requires genesis to have a difficulty of 0, instead got %v", g.Difficulty)
	}
	if db == nil {
		db = rawdb.NewMemoryDatabase()
	}
	statedb, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
	if err != nil {
		panic(err)
	}
	for addr, account := range g.Alloc {
		if account.Balance != nil {
			statedb.AddBalance(addr, account.Balance)
		}
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}

	genesisBonds := g.Alloc.ToGenesisBonds()
	evm := genesisEVM(g, statedb)
	if err := autonity.ExecuteGenesisSequence(g.Config, genesisBonds, evm); err != nil {
		return nil, fmt.Errorf("cannot execute genesis sequence: %w", err)
	}

	committee, err := autonity.CallGetCommittee(evm)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve genesis committee: %w", err)
	}

	root := statedb.IntermediateRoot(false)
	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       g.Timestamp,
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   g.GasLimit,
		GasUsed:    g.GasUsed,
		BaseFee:    g.BaseFee,
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		Root:       root,
		Round:      0,
	}

	epoch := &types.Epoch{
		Committee:          committee,
		PreviousEpochBlock: common.Big0,
		NextEpochBlock:     new(big.Int).SetUint64(g.Config.AutonityContractConfig.EpochPeriod),
		OmissionDelta:      new(big.Int).SetUint64(g.Config.OmissionAccountabilityConfig.Delta),
		Eip1559: &types.Eip1559Params{
			MinBaseFee:               new(big.Int).SetUint64(g.Config.AutonityContractConfig.MinBaseFee),
			BaseFeeChangeDenominator: new(big.Int).SetUint64(g.Config.AutonityContractConfig.BaseFeeChangeDenominator),
			ElasticityMultiplier:     new(big.Int).SetUint64(g.Config.AutonityContractConfig.ElasticityMultiplier),
			GasLimitBoundDivisor:     new(big.Int).SetUint64(g.Config.AutonityContractConfig.GasLimitBoundDivisor),
		},
	}
	head.Epoch = epoch

	if g.GasLimit == 0 {
		head.GasLimit = params.GenesisGasLimit
	}
	if g.Difficulty == nil && g.Mixhash == (common.Hash{}) {
		head.Difficulty = params.GenesisDifficulty
	}
	if g.Config != nil && g.Config.IsLondon(common.Big0) {
		if g.BaseFee != nil {
			head.BaseFee = g.BaseFee
		} else {
			head.BaseFee = new(big.Int).SetUint64(params.InitialBaseFee)
		}
	}
	statedb.Commit(false)
	statedb.Database().TrieDB().Commit(root, true, nil)

	return types.NewBlock(head, nil, nil, nil, trie.NewStackTrie(nil)), nil
}

func genesisEVM(genesis *Genesis, statedb vm.StateDB) *vm.EVM {
	evmContext := vm.BlockContext{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     func(n uint64) common.Hash { return common.Hash{} },
		Coinbase:    genesis.Coinbase,
		BlockNumber: big.NewInt(0),
		Time:        new(big.Int).SetUint64(genesis.Timestamp),
		GasLimit:    genesis.GasLimit,
		Difficulty:  genesis.Difficulty,

		ActivityProof:      nil,
		ActivityProofRound: 0,
	}
	txContext := vm.TxContext{
		Origin:   params.DeployerAddress,
		GasPrice: new(big.Int).SetUint64(0x0),
	}
	return vm.NewEVM(evmContext, txContext, statedb, genesis.Config, vm.Config{})
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(db ethdb.Database) (*types.Block, error) {
	if g.Config == nil {
		g.Config = params.TestChainConfig
	}

	if err := g.Config.CheckConfigForkOrder(); err != nil {
		return nil, err
	}

	block, err := g.ToBlock(db)
	if err != nil {
		return nil, err
	}

	if block.Number().Sign() != 0 {
		return nil, errors.New("can't commit genesis block with number > 0")
	}

	rawdb.WriteTd(db, block.Hash(), block.NumberU64(), g.Difficulty)
	rawdb.WriteBlock(db, block)
	rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), nil)
	rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
	rawdb.WriteHeadBlockHash(db, block.Hash())
	rawdb.WriteHeadFastBlockHash(db, block.Hash())
	rawdb.WriteHeadHeaderHash(db, block.Hash())
	rawdb.WriteEpochHeaderHash(db, block.Hash())
	rawdb.WriteChainConfig(db, block.Hash(), g.Config)
	rawdb.WriteContractsConfig(db, block.NumberU64(), &types.ContractsConfig{
		EpochPeriod:         new(big.Int).SetUint64(g.Config.AutonityContractConfig.EpochPeriod),
		BlockPeriod:         new(big.Int).SetUint64(g.Config.AutonityContractConfig.BlockPeriod),
		GasLimit:            new(big.Int).SetUint64(g.Config.AutonityContractConfig.GasLimit),
		ClusteringThreshold: new(big.Int).SetUint64(g.Config.AutonityContractConfig.ClusteringThreshold),
		Accountability: types.AccountabilityParams{
			Range:       new(big.Int).SetUint64(g.Config.AccountabilityConfig.Range),
			Delta:       new(big.Int).SetUint64(g.Config.AccountabilityConfig.Delta),
			GracePeriod: new(big.Int),
		},
		Eip1559: types.Eip1559Params{
			MinBaseFee:               new(big.Int).SetUint64(g.Config.AutonityContractConfig.MinBaseFee),
			BaseFeeChangeDenominator: new(big.Int).SetUint64(g.Config.AutonityContractConfig.BaseFeeChangeDenominator),
			ElasticityMultiplier:     new(big.Int).SetUint64(g.Config.AutonityContractConfig.ElasticityMultiplier),
			GasLimitBoundDivisor:     new(big.Int).SetUint64(g.Config.AutonityContractConfig.GasLimitBoundDivisor),
		},
	})
	return block, nil
}

func (g *Genesis) setDefaultHardforks() {
	if g.Config.ByzantiumBlock == nil {
		g.Config.ByzantiumBlock = new(big.Int)
	}
	if g.Config.HomesteadBlock == nil {
		g.Config.HomesteadBlock = new(big.Int)
	}
	if g.Config.ConstantinopleBlock == nil {
		g.Config.ConstantinopleBlock = new(big.Int)
	}
	if g.Config.PetersburgBlock == nil {
		g.Config.PetersburgBlock = new(big.Int)
	}
	if g.Config.IstanbulBlock == nil {
		g.Config.IstanbulBlock = new(big.Int)
	}
	if g.Config.MuirGlacierBlock == nil {
		g.Config.MuirGlacierBlock = new(big.Int)
	}
	if g.Config.BerlinBlock == nil {
		g.Config.BerlinBlock = new(big.Int)
	}
	if g.Config.LondonBlock == nil {
		g.Config.LondonBlock = new(big.Int)
	}
	if g.Config.ArrowGlacierBlock == nil {
		g.Config.ArrowGlacierBlock = new(big.Int)
	}
	if g.Config.EIP158Block == nil {
		g.Config.EIP158Block = new(big.Int)
	}
	if g.Config.EIP150Block == nil {
		g.Config.EIP150Block = new(big.Int)
	}
	if g.Config.EIP155Block == nil {
		g.Config.EIP155Block = new(big.Int)
	}
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db ethdb.Database) *types.Block {
	block, err := g.Commit(db)
	if err != nil {
		panic(err)
	}
	return block
}

// GenesisBlockForTesting creates and writes a block in which addr has the given wei balance.
func GenesisBlockForTesting(db ethdb.Database, addr common.Address, balance *big.Int) *types.Block {
	g := Genesis{
		Alloc:   GenesisAlloc{addr: {Balance: balance}},
		Config:  params.TestChainConfig,
		BaseFee: big.NewInt(params.InitialBaseFee),
		Mixhash: types.BFTDigest,
	}
	return g.MustCommit(db)
}

// DefaultGenesisBlock returns a default genesis block for testing purposes.
func DefaultGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.TestChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa"),
		GasLimit:   5000,
		Difficulty: big.NewInt(0),
		BaseFee:    big.NewInt(params.InitialBaseFee),
		Alloc:      decodePrealloc(mainnetAllocData),
	}
}

// DefaultPiccadillyGenesisBlock returns the Piccadilly network genesis block.
func DefaultPiccadillyGenesisBlock() *Genesis {
	var (
		sdpAccount = common.HexToAddress("0x8b914020A7099E4723f45561E897fa2740885A55")
	)
	g := &Genesis{
		Config:     params.PiccadillyChainConfig,
		Timestamp:  uint64(params.PiccadillyGenesisUnixTimestamp),
		Nonce:      0,
		GasLimit:   20_000_000,
		Difficulty: big.NewInt(0),
		Mixhash:    types.BFTDigest,
		Alloc: map[common.Address]GenesisAccount{
			sdpAccount: { // SDP Simulator Account
				Bonds: make(map[common.Address]*big.Int),
			},
		},
	}

	for _, alloc := range params.PiccadillyATNallocs {
		if prev, ok := g.Alloc[alloc.Address]; ok {
			prev.Balance = alloc.Value
			g.Alloc[alloc.Address] = prev
		} else {
			g.Alloc[alloc.Address] = GenesisAccount{
				Balance: alloc.Value,
			}
		}
	}

	for _, alloc := range params.PiccadillyNTNallocs {
		if prev, ok := g.Alloc[alloc.Address]; ok {
			prev.NewtonBalance = alloc.Value
			g.Alloc[alloc.Address] = prev
		} else {
			g.Alloc[alloc.Address] = GenesisAccount{
				NewtonBalance: alloc.Value,
			}
		}
	}

	for _, v := range params.PiccadillySDPDelegations {
		g.Alloc[sdpAccount].Bonds[v.Address] = v.Value
	}
	return g
}

func DefaultBakerlooGenesisBlock() *Genesis {
	var (
		sdpAccount     = common.HexToAddress("0xC1122BEb440E68c4c85415DBEF80ef97247aac85")
		genesisTime, _ = time.Parse(time.RFC3339, "2025-08-06T13:00:00Z")
	)
	params.BakerlooChainConfig.AutonityContractConfig.Schedules = params.BakerlooSchedules()
	g := &Genesis{
		Config:     params.BakerlooChainConfig,
		Timestamp:  uint64(genesisTime.Unix()), //nolint
		Nonce:      0,
		GasLimit:   30_000_000,
		BaseFee:    big.NewInt(10_000_000_000),
		Difficulty: big.NewInt(0),
		Mixhash:    types.BFTDigest,
		Alloc: map[common.Address]GenesisAccount{
			// EXT01-BAK
			common.HexToAddress("D1d46d82a92AfA9EC1AE3d83Ae23250D8b4D6317"): {
				NewtonBalance: new(big.Int).Mul(big.NewInt(25000), params.NtnPrecision),
				Balance:       new(big.Int).Mul(big.NewInt(2), params.ATNPrecision),
			},
			// CTL-VAULT-BAK1
			common.HexToAddress("0x001f0994E8Fc36D123A461B8995948E9F01054C8"): {
				NewtonBalance: new(big.Int).Mul(big.NewInt(1000000), params.NtnPrecision),
				Balance:       new(big.Int).Mul(big.NewInt(10), params.ATNPrecision),
			},
			// CTL-VAULT-BAK2
			common.HexToAddress("0xd0C72d80B0D24f329a5405A3aEF2E1B93551fc6d"): {
				NewtonBalance: new(big.Int).Mul(big.NewInt(1000000), params.NtnPrecision),
			},
			// LC-OPS-BAK
			common.HexToAddress("0x879417AFBB552D64000E0F1Ae449615983b94633"): {
				NewtonBalance: new(big.Int).Mul(big.NewInt(2000000), params.NtnPrecision),
			},
			// BOOTSTRAP-BAK
			common.HexToAddress("0x298550A7e719dca9C2F8E01a70c1Fe31345A9da8"): {
				NewtonBalance: new(big.Int).Mul(big.NewInt(208000), params.NtnPrecision),
				Balance:       new(big.Int).Mul(big.NewInt(100), params.ATNPrecision),
			},
			// SDP-OPS-BAK
			sdpAccount: {
				NewtonBalance: new(big.Int).Mul(big.NewInt(162000), params.NtnPrecision),
				Balance:       new(big.Int).Mul(big.NewInt(100), params.ATNPrecision),
				Bonds:         make(map[common.Address]*big.Int),
			},
			// AT-VAULT1-BAK
			common.HexToAddress("0x0EB1a6E69ce1b9Cb8540863EC32e483d819704D2"): {
				NewtonBalance: new(big.Int).Mul(big.NewInt(1257955633333332), big.NewInt(10000000000)),
				Balance:       new(big.Int).Mul(big.NewInt(100), params.ATNPrecision),
			},
			// AT-VAULT2-BAK
			common.HexToAddress("0xb8056F4fFBffe5c42637208286D7e7067Ed879C4"): {
				NewtonBalance: new(big.Int).Mul(big.NewInt(10000000), params.NtnPrecision),
			},
			// PO-BAK
			common.HexToAddress("0x83e5e0eab996Bb894814fa8F0AC96a0D314f06F3"): {
				Balance: new(big.Int).Mul(big.NewInt(10), params.ATNPrecision),
			},
			// Safe Singleton Factory
			common.HexToAddress("914d7Fec6aaC8cd542e72Bca78B30650d45643d7"): {
				Code: common.Hex2Bytes("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf3"),
			},
		},
	}

	// GVs initial ATN Allocs
	for _, v := range params.BakerlooValidators {
		g.Alloc[*v.NodeAddress] = GenesisAccount{Balance: big.NewInt(params.Ether)}
		g.Alloc[v.OracleAddress] = GenesisAccount{Balance: big.NewInt(params.Ether)}
		g.Alloc[v.Treasury] = GenesisAccount{Balance: big.NewInt(params.Ether)}
	}
	// SDP allocations
	for _, v := range params.BakerlooValidators {
		g.Alloc[sdpAccount].Bonds[*v.NodeAddress] = new(big.Int).Mul(big.NewInt(60_000), params.NtnPrecision)
	}
	return g
}

func DefaultMainnetGenesisBlock() *Genesis {
	var (
		sdpAccount     = common.HexToAddress("0x3AFcA309da96E9C55da6A9B34f172C76f5CB857c")
		genesisTime, _ = time.Parse(time.RFC3339, "2025-08-12T13:00:00Z")
	)
	params.AutMainnetChainConfig.AutonityContractConfig.Schedules = params.MainnetSchedules()
	g := &Genesis{
		Config:     params.AutMainnetChainConfig,
		Timestamp:  uint64(genesisTime.Unix()), //nolint
		Nonce:      0,
		GasLimit:   30_000_000,
		BaseFee:    big.NewInt(10_000_000_000),
		Difficulty: big.NewInt(0),
		Mixhash:    types.BFTDigest,
		Alloc: map[common.Address]GenesisAccount{
			// 82
			common.HexToAddress("0x3CE366baf7167e2f5A228b3f9ae27c6dD42F5DF6"): {
				NewtonBalance: tokenToBig(25000, 0),
				Balance:       tokenToBig(2, 0),
			},
			// 83
			common.HexToAddress("0xa06596342d532601ac2A56D9A222fddcA0030Db5"): {
				NewtonBalance: tokenToBig(100000, 0),
				Balance:       tokenToBig(2, 0),
			},
			// 84
			common.HexToAddress("0xB4ddd87536ac80d8Ed6c8295eD9475aa5516E184"): {
				NewtonBalance: tokenToBig(100000, 0),
				Balance:       tokenToBig(2, 0),
			},
			// 85
			common.HexToAddress("0x4ad5E2D166590bC59065d2848a5b20cE073DD702"): {
				NewtonBalance: tokenToBig(62018553506, 8),
				Balance:       tokenToBig(2, 0),
			},
			// 86
			common.HexToAddress("0xE3793Ce105661f13d637C544ef1236EE53BCf920"): {
				NewtonBalance: tokenToBig(25000, 0),
				Balance:       tokenToBig(2, 0),
			},
			// 87
			common.HexToAddress("0xcbFd1903A6C73190b459758E942D319043a42141"): {
				NewtonBalance: tokenToBig(50000, 0),
				Balance:       tokenToBig(2, 0),
			},
			// 88
			common.HexToAddress("0xa644A95eEDa1CfCE135dE1e81781254eABE89c1F"): {
				NewtonBalance: tokenToBig(100000, 0),
				Balance:       tokenToBig(2, 0),
			},
			// 89
			common.HexToAddress("0xD16689903C878dCd81F03bDe25E2841d0A399792"): {
				NewtonBalance: tokenToBig(1985007023835, 8),
				Balance:       tokenToBig(2, 0),
			},
			// 90
			common.HexToAddress("0xB8FaB011db3d31ba8437EFE783b087743B627B0A"): {
				NewtonBalance: tokenToBig(25000, 0),
				Balance:       tokenToBig(2, 0),
			},
			// 91
			common.HexToAddress("0x2e2dFb3fc658AD6E64b08370dBB8FF27Fb3758bf"): {
				NewtonBalance: tokenToBig(404806029861, 8),
				Balance:       tokenToBig(2, 0)},
			// 92
			common.HexToAddress("0x2f3133Fbfa33D9DdB5C5C4BdbDB9Bcb6393D5507"): {
				NewtonBalance: tokenToBig(2397197789349, 8),
				Balance:       tokenToBig(2, 0),
			},
			// 93
			common.HexToAddress("0xA39DbbDbC5225aB9A4b9A08Efdf1a5fdd3D97452"): {
				NewtonBalance: tokenToBig(1865472199569, 8),
				Balance:       tokenToBig(2, 0),
			},
			// 94
			common.HexToAddress("0x490E99D342482dC1e12843787904DbCC8bA02052"): {
				NewtonBalance: tokenToBig(50000, 0),
				Balance:       tokenToBig(2, 0),
			},
			// 95
			common.HexToAddress("0x08042E939a9dD38A71D4b5c548EC2a0cEb6eFeCb"): {
				NewtonBalance: tokenToBig(854224389593, 8),
				Balance:       tokenToBig(2, 0),
			},
			// 96
			common.HexToAddress("0x3abCDd9932E47Bb43552a404a7958C26AC065bb0"): {
				NewtonBalance: tokenToBig(1389607347619, 8),
				Balance:       tokenToBig(2, 0),
			},
			// 97
			common.HexToAddress("0xA079FD9bFdcb13402b4d805FF7810C0eb8458fbd"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 98
			common.HexToAddress("0x21dA5579c46C6A20324f20246c5F709D5a964975"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 99
			common.HexToAddress("0x2f8F16D2CbA5F02A3cFF144C26cc1252610f2EdF"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 100
			common.HexToAddress("0x1af6100F63b9e534D306FF895C319Fa91EDfD4e4"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 101
			common.HexToAddress("0xE631fAf0A8a738bc73D2Df4088909A3ebda92A35"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 102
			common.HexToAddress("0x189b26e54C747e2998D3D5e0d8A7EaA98945e0b7"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 103
			common.HexToAddress("0x4B2B38BdB50d7AAE589bC980c57C7EDEF550b1ad"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 104
			common.HexToAddress("0x24ef0F90B69c5E5710cCB30e1558e92B1cd86871"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 105
			common.HexToAddress("0x5A2ABf6eDcF95515ae0d07DBd286a6c907462276"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 106
			common.HexToAddress("0x1246C78d3CdDBD7BD391ad2fA35e450DB5A302A4"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 107
			common.HexToAddress("0xF639Dc6580732fBcCE32438a089Ab8c3A27Fd648"): {
				NewtonBalance: tokenToBig(1000000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 108
			common.HexToAddress("0xE710d8c5d85B5c318ff54d8158Cf092906D306B5"): {
				NewtonBalance: tokenToBig(100000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 109
			common.HexToAddress("0xA58dE4900B2284503fb61DC5CcdC0aBa605d0fF4"): {
				NewtonBalance: tokenToBig(30000, 0),
				Balance:       tokenToBig(10, 0),
			},
			// 111
			common.HexToAddress("0x946b4f0e4593F630Dac3F6100A154B7D417709cC"): {
				NewtonBalance: tokenToBig(4000000, 0),
			},
			// 112
			common.HexToAddress("0x930e9Ec3A2C3E496973cA547fc2F10245F44727d"): {
				NewtonBalance: tokenToBig(379000, 0),
				Balance:       tokenToBig(100, 0),
			},
			// 113
			sdpAccount: {
				NewtonBalance: tokenToBig(162000, 0),
				Balance:       tokenToBig(100, 0),
				Bonds:         make(map[common.Address]*big.Int),
			},
			// 134
			common.HexToAddress("0x051B6eD709C72961e2b061b4d07d4a1D407CA39d"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 135
			common.HexToAddress("0x8387C4cE6Ce1E56E151F241a50A805a9402b857A"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 136
			common.HexToAddress("0xd1564bF2A7748da79d79EAC1fE050Aed96271a12"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 137
			common.HexToAddress("0x408B1994Ebf30b8F2aA150B73Ae03D4B2E1A83Be"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 138
			common.HexToAddress("0x2BFeda4efC5a379E5C9CEbE169352a7527222f7F"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 139
			common.HexToAddress("0xc6d8c3fc26C251db394E745e2405c20C1d33Cf72"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 140
			common.HexToAddress("0x4A3d8e78B8DbA425d97b5b3eBb3b6dde33F4De39"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 141
			common.HexToAddress("0x7C6A432125E52ca197EAEfA6fE4C444B97B11770"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 142
			common.HexToAddress("0x5124EA7CaA3dFc1e3bfCd0839Cb060Ae57174255"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 143
			common.HexToAddress("0xF7e2c4A224429f8E7BEB042e239Ca743F5E36Dcf"): {
				NewtonBalance: tokenToBig(1000000, 0),
			},
			// 144
			common.HexToAddress("0x43fFDB2DF482b12D2dbbCeaA2A9B2F730cC98b65"): {
				NewtonBalance: tokenToBig(500000, 0),
			},
			// 145
			common.HexToAddress("0x4F4CA7Bf349Cd3BC312ce4aAfA6D42eA25a8AA6b"): {
				NewtonBalance: tokenToBig(238973, 0),
				Balance:       tokenToBig(100, 0),
			},
			// 146
			common.HexToAddress("0x5fB82096CdFc95755b7b766E94F80265C9Fc4bFC"): {
				Balance: tokenToBig(10, 0),
			},
			// Safe Singleton Factory
			common.HexToAddress("914d7Fec6aaC8cd542e72Bca78B30650d45643d7"): {
				Code: common.Hex2Bytes("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf3"),
			},
		},
	}

	// GVs initial ATN Allocs
	for _, v := range params.MainnetValidators {
		g.Alloc[*v.NodeAddress] = GenesisAccount{Balance: big.NewInt(params.Ether)}
		g.Alloc[v.OracleAddress] = GenesisAccount{Balance: big.NewInt(params.Ether)}
		g.Alloc[v.Treasury] = GenesisAccount{Balance: big.NewInt(params.Ether)}
	}
	// SDP allocations
	for _, v := range params.MainnetValidators {
		g.Alloc[sdpAccount].Bonds[*v.NodeAddress] = new(big.Int).Mul(big.NewInt(60_000), params.NtnPrecision)
	}
	return g
}

// DefaultRopstenGenesisBlock returns the Ropsten network genesis block.
func DefaultRopstenGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.RopstenChainConfig,
		Nonce:      66,
		ExtraData:  hexutil.MustDecode("0x3535353535353535353535353535353535353535353535353535353535353535"),
		GasLimit:   16777216,
		Difficulty: big.NewInt(1048576),
		Alloc:      decodePrealloc(ropstenAllocData),
	}
}

// DefaultRinkebyGenesisBlock returns the Rinkeby network genesis block.
func DefaultRinkebyGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.RinkebyChainConfig,
		Timestamp:  1492009146,
		ExtraData:  hexutil.MustDecode("0x52657370656374206d7920617574686f7269746168207e452e436172746d616e42eb768f2244c8811c63729a21a3569731535f067ffc57839b00206d1ad20c69a1981b489f772031b279182d99e65703f0076e4812653aab85fca0f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   4700000,
		Difficulty: big.NewInt(1),
		Alloc:      decodePrealloc(rinkebyAllocData),
	}
}

// DefaultGoerliGenesisBlock returns the Görli network genesis block.
func DefaultGoerliGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.GoerliChainConfig,
		Timestamp:  1548854791,
		ExtraData:  hexutil.MustDecode("0x22466c6578692069732061207468696e6722202d204166726900000000000000e0a2bd4258d2768837baa26a28fe71dc079f84c70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   10485760,
		Difficulty: big.NewInt(1),
		Alloc:      decodePrealloc(goerliAllocData),
	}
}

// DeveloperGenesisBlock returns the 'autonity --dev' genesis block.
func DeveloperGenesisBlock(gasLimit uint64, faucet *keystore.Key) *Genesis {
	validatorEnode := enode.NewV4(&faucet.PrivateKey.PublicKey, net.ParseIP("0.0.0.0"), 0, 0)
	testAutonityContractConfig := params.AutonityContractGenesis{
		MaxCommitteeSize:         1,
		BlockPeriod:              1,
		UnbondingPeriod:          120,
		EpochPeriod:              60, //seconds
		GasLimit:                 params.DefaultGenesisGasLimit,
		GasLimitBoundDivisor:     params.DefaultGasLimitBoundDivisor,
		BaseFeeChangeDenominator: params.DefaultBaseFeeChangeDenominator,
		ElasticityMultiplier:     params.DefaultElasticityMultiplier,
		ClusteringThreshold:      params.DefaultClusteringThreshold,
		DelegationRate:           1200,             // 12%
		WithholdingThreshold:     0,                // 0%, no tolerance
		ProposerRewardRate:       1000,             // 10%
		OracleRewardRate:         1000,             // 10%
		TreasuryFee:              1500000000000000, // 0.15%,
		MinBaseFee:               10000000000,
		Operator:                 faucet.Address,
		Treasury:                 faucet.Address,
		WithheldRewardsPool:      faucet.Address,
		InitialInflationReserve:  params.TestAutonityContractConfig.InitialInflationReserve,
		SkipGenesisVerification:  true,
		Validators: []*params.Validator{
			{
				Treasury:      faucet.Address,
				OracleAddress: faucet.Address,
				Enode:         validatorEnode.String(),
				BondedStake:   new(big.Int).SetUint64(1000),
				ConsensusKey:  params.TestValidatorConsensusKey.PublicKey().Marshal(),
			},
		},
	}
	testChainConfig := &params.ChainConfig{
		ChainID:                      big.NewInt(65111111),
		HomesteadBlock:               big.NewInt(0),
		EIP150Block:                  big.NewInt(0),
		EIP155Block:                  big.NewInt(0),
		EIP158Block:                  big.NewInt(0),
		ByzantiumBlock:               big.NewInt(0),
		ConstantinopleBlock:          big.NewInt(0),
		PetersburgBlock:              big.NewInt(0),
		IstanbulBlock:                big.NewInt(0),
		MuirGlacierBlock:             big.NewInt(0),
		BerlinBlock:                  big.NewInt(0),
		LondonBlock:                  big.NewInt(0),
		ArrowGlacierBlock:            big.NewInt(0),
		AutonityContractConfig:       &testAutonityContractConfig,
		AccountabilityConfig:         params.DefaultAccountabilityConfig,
		OracleContractConfig:         params.DefaultGenesisOracleConfig,
		OmissionAccountabilityConfig: params.DefaultOmissionAccountabilityConfig,
	}
	if err := testChainConfig.Prepare(); err != nil {
		log.Error("Error preparing chain configuration, err:", err)
	}

	return &Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		Mixhash:    types.BFTDigest,
		ExtraData:  []byte{},
		GasLimit:   gasLimit,
		BaseFee:    big.NewInt(15000000000),
		Difficulty: big.NewInt(0),
		Alloc: map[common.Address]GenesisAccount{
			faucet.Address: {Balance: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(9))},
		},
		Config: testChainConfig,
	}
}

func decodePrealloc(data string) GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
	}
	return ga
}

func tokenToBig(num int64, scale int64) *big.Int {
	s := new(big.Int).Exp(big.NewInt(10), big.NewInt(18-scale), nil)
	return new(big.Int).Mul(big.NewInt(num), s)
}
