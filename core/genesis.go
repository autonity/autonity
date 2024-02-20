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
	ret := make(autonity.GenesisBonds, len(*ga))
	for addr, alloc := range *ga {
		ret[addr] = autonity.GenesisBond{
			NewtonBalance: alloc.NewtonBalance,
			Bonds:         alloc.Bonds,
		}
	}
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
			// No genesis from DB and configuration, don't start node with GETH main-net genesis block.
			return nil, common.Hash{}, fmt.Errorf("DB has no genesis block and there is no genesis file set by user")
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
	if g.Config.AutonityContractConfig == nil {
		return nil, fmt.Errorf("autonity config section missing in genesis")
	}
	g.setDefaultHardforks()
	if err := g.Config.AutonityContractConfig.Prepare(); err != nil {
		return nil, err
	}
	if g.Difficulty == nil {
		g.Difficulty = params.GenesisDifficulty
	}
	if g.Difficulty.Cmp(big.NewInt(0)) != 0 {
		return nil, fmt.Errorf("autonity requires genesis to have a difficulty of 0, instead got %v", g.Difficulty)
	}
	committee, err := extractCommittee(g.Config.AutonityContractConfig.Validators)
	if err != nil {
		return nil, err
	}
	if db == nil {
		db = rawdb.NewMemoryDatabase()
	}
	statedb, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
	if err != nil {
		panic(err)
	}
	for addr, account := range g.Alloc {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}

	genesisBonds := g.Alloc.ToGenesisBonds()

	evmProvider := func(statedb vm.StateDB) *vm.EVM {
		return genesisEVM(g, statedb)
	}

	evmContracts := autonity.NewGenesisEVMContract(evmProvider, statedb, db, g.Config)

	if err := autonity.DeployContracts(g.Config, genesisBonds, evmContracts); err != nil {
		return nil, fmt.Errorf("cannot deploy contracts: %w", err)
	}

	root := statedb.IntermediateRoot(false)
	head := &types.Header{
		Number:         new(big.Int).SetUint64(g.Number),
		Nonce:          types.EncodeNonce(g.Nonce),
		Time:           g.Timestamp,
		ParentHash:     g.ParentHash,
		Extra:          g.ExtraData,
		GasLimit:       g.GasLimit,
		GasUsed:        g.GasUsed,
		BaseFee:        g.BaseFee,
		Difficulty:     g.Difficulty,
		MixDigest:      g.Mixhash,
		Coinbase:       g.Coinbase,
		Root:           root,
		Round:          0,
		Committee:      committee,
		LastEpochBlock: common.Big0,
	}
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

	rawdb.WriteChainConfig(db, block.Hash(), g.Config)
	return block, nil
}

// extractCommittee takes a slice of autonity users and extracts the validators
// into a new type 'types.Committee' which is returned. It returns an error if
// the provided users contained no validators.
func extractCommittee(validators []*params.Validator) (*types.Committee, error) {
	committee := &types.Committee{Members: make([]*types.CommitteeMember, len(validators))}
	for i, v := range validators {
		committee.Members[i] = &types.CommitteeMember{
			Address:      *v.NodeAddress,
			VotingPower:  v.BondedStake,
			ConsensusKey: v.ConsensusKey,
		}
	}
	committee.Sort()

	if len(committee.Members) == 0 {
		return nil, fmt.Errorf("no validators specified in the initial autonity validators")
	}

	log.Info("Starting DPoS-BFT consensus protocol", "validators", committee)
	return committee, nil
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
	g := &Genesis{
		Config:     params.PiccaddillyChainConfig,
		Nonce:      0,
		GasLimit:   30_000_000,
		Difficulty: big.NewInt(0),
		Mixhash:    types.BFTDigest,
		Alloc: map[common.Address]GenesisAccount{
			params.PiccaddillyChainConfig.AutonityContractConfig.Operator: {Balance: new(big.Int).Mul(big.NewInt(3), big.NewInt(params.Ether))},
		},
	}
	for _, v := range g.Config.AutonityContractConfig.Validators {
		g.Alloc[*v.NodeAddress] = GenesisAccount{Balance: big.NewInt(params.Ether)}
		g.Alloc[v.OracleAddress] = GenesisAccount{Balance: big.NewInt(params.Ether)}
	}
	return g
}

func DefaultBakerlooGenesisBlock() *Genesis {
	g := &Genesis{
		Config:     params.BakerlooChainConfig,
		Nonce:      0,
		GasLimit:   30_000_000,
		Difficulty: big.NewInt(0),
		Mixhash:    types.BFTDigest,
		Alloc: map[common.Address]GenesisAccount{
			params.BakerlooChainConfig.AutonityContractConfig.Operator: {Balance: new(big.Int).Mul(big.NewInt(3), big.NewInt(params.Ether))},
		},
	}
	for _, v := range g.Config.AutonityContractConfig.Validators {
		g.Alloc[*v.NodeAddress] = GenesisAccount{Balance: big.NewInt(params.Ether)}
		g.Alloc[v.OracleAddress] = GenesisAccount{Balance: big.NewInt(params.Ether)}
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
		MaxCommitteeSize: 1,
		BlockPeriod:      1,
		UnbondingPeriod:  120,
		EpochPeriod:      30,               //seconds
		DelegationRate:   1200,             // 12%
		TreasuryFee:      1500000000000000, // 0.15%,
		MinBaseFee:       10000000000,
		Operator:         faucet.Address,
		Treasury:         faucet.Address,
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
	if err := testAutonityContractConfig.Prepare(); err != nil {
		log.Error("Error preparing contract, err:", err)
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
		Config: &params.ChainConfig{
			ChainID:                big.NewInt(65111111),
			HomesteadBlock:         big.NewInt(0),
			EIP150Block:            big.NewInt(0),
			EIP155Block:            big.NewInt(0),
			EIP158Block:            big.NewInt(0),
			ByzantiumBlock:         big.NewInt(0),
			ConstantinopleBlock:    big.NewInt(0),
			PetersburgBlock:        big.NewInt(0),
			IstanbulBlock:          big.NewInt(0),
			MuirGlacierBlock:       big.NewInt(0),
			BerlinBlock:            big.NewInt(0),
			LondonBlock:            big.NewInt(0),
			ArrowGlacierBlock:      big.NewInt(0),
			AutonityContractConfig: &testAutonityContractConfig,
			AccountabilityConfig:   params.DefaultAccountabilityConfig,
			OracleContractConfig:   params.DefaultGenesisOracleConfig,
		},
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
