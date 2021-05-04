// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"unicode"

	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/log"
	"github.com/davecgh/go-spew/spew"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/clearmatics/autonity/cmd/utils"
	"github.com/clearmatics/autonity/eth"
	"github.com/clearmatics/autonity/internal/ethapi"
	"github.com/clearmatics/autonity/node"
	"github.com/clearmatics/autonity/params"
	"github.com/naoina/toml"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(nodeFlags, rpcFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

type autonityConfig struct {
	Eth      eth.Config
	Node     node.Config
	Ethstats ethstatsConfig
}

func loadConfig(file string, cfg *autonityConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "autonity.ipc"
	return cfg
}

func makeConfigNode(ctx *cli.Context) (*node.Node, autonityConfig) {
	// Load defaults.
	cfg := autonityConfig{
		Eth:  eth.DefaultConfig,
		Node: defaultNodeConfig(),
	}

	// Load config file.
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)
	stack, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}
	utils.SetEthConfig(ctx, stack, &cfg.Eth)
	if ctx.GlobalIsSet(utils.EthStatsURLFlag.Name) {
		cfg.Ethstats.URL = ctx.GlobalString(utils.EthStatsURLFlag.Name)
	}

	return stack, cfg
}

// loadGenesisFile will load and validate the given JSON format genesis file.
func loadGenesisFile(genesisPath string) (*core.Genesis, error) {
	file, err := os.Open(genesisPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	genesis := new(core.Genesis)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		return nil, err
	}
	// Make AutonityContract and Tendermint consensus mandatory for the time being.
	if genesis.Config == nil {
		return nil, fmt.Errorf("no Autonity Contract and Tendermint configs section in genesis")
	}
	if genesis.Config.AutonityContractConfig == nil {
		return nil, fmt.Errorf("no Autonity Contract config section in genesis")
	}

	if err := genesis.Config.AutonityContractConfig.Prepare(); err != nil {
		spew.Dump(genesis.Config.AutonityContractConfig)
		return nil, err
	}

	return genesis, nil
}

// applyGenesis attempts to apply the given genesis to the node database, the
// first time this is run for a node it initializes the genesis block in the
// database.
func applyGenesis(genesis *core.Genesis, node *node.Node) error {
	for _, name := range []string{"chaindata", "lightchaindata"} {
		chaindb, err := node.OpenDatabase(name, 0, 0, "")
		if err != nil {
			return fmt.Errorf("failed to open database: %v", err)
		}
		defer chaindb.Close()
		_, hash, err := core.SetupGenesisBlock(chaindb, genesis)
		if err != nil {
			return fmt.Errorf("failed to write genesis block: %v", err)
		}
		log.Info("Successfully wrote genesis state", "database", name, "hash", hash)
	}
	return nil
}

// If genesis flag is not set, node will load chain-data from data-dir. If the flat is set, node will load the
// genesis file, check if genesis file is match with genesis block, and check if chain configuration of the genesis
// file is compatible with current chain-data, apply new compatible chain configuration into chain db.
// Otherwise client will end up with a mis-match genesis error or an incompatible chain configuration error.
func initGenesisBlockOnStart(ctx *cli.Context, stack *node.Node) {
	genesisPath := ctx.GlobalString(utils.InitGenesisFlag.Name)
	if genesisPath != "" {
		log.Info("Trying to initialise genesis block with genesis file", "filepath", genesisPath)
		genesis, err := loadGenesisFile(genesisPath)
		if err != nil {
			utils.Fatalf("failed to validate genesis file: %v", err)
		}
		err = applyGenesis(genesis, stack)
		if err != nil {
			utils.Fatalf("failed to apply genesis file: %v", err)
		}
	}
}

func makeFullNode(ctx *cli.Context) (*node.Node, ethapi.Backend) {
	stack, cfg := makeConfigNode(ctx)

	// Combine init genesis on node start up.
	initGenesisBlockOnStart(ctx, stack)

	backend := utils.RegisterEthService(stack, &cfg.Eth)

	// Configure GraphQL if requested
	if ctx.GlobalIsSet(utils.GraphQLEnabledFlag.Name) {
		utils.RegisterGraphQLService(stack, backend, cfg.Node)
	}
	// Add the Ethereum Stats daemon if requested.
	if cfg.Ethstats.URL != "" {
		utils.RegisterEthStatsService(stack, backend, cfg.Ethstats.URL)
	}
	return stack, backend
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	_, cfg := makeConfigNode(ctx)
	comment := ""

	if cfg.Eth.Genesis != nil {
		cfg.Eth.Genesis = nil
		comment += "# Note: this config doesn't contain the genesis block.\n\n"
	}

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}

	dump := os.Stdout
	if ctx.NArg() > 0 {
		dump, err = os.OpenFile(ctx.Args().Get(0), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dump.Close()
	}
	dump.WriteString(comment)
	dump.Write(out)

	return nil
}
