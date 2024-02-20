// Copyright 2016 The go-ethereum Authors
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
	"crypto/ecdsa"
	"fmt"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	ethproto "github.com/autonity/autonity/eth/protocols/eth"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/autonity/autonity/cmd/utils"
	"github.com/autonity/autonity/consensus/ethash"
	"github.com/autonity/autonity/params"
	"gopkg.in/urfave/cli.v1"
)

var (
	versionCommand = cli.Command{
		Action:    utils.MigrateFlags(version),
		Name:      "version",
		Usage:     "Print version numbers",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is supposed to be machine-readable.
`,
	}

	licenseCommand = cli.Command{
		Action:    utils.MigrateFlags(license),
		Name:      "license",
		Usage:     "Display license information",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
	}

	ownershipProofCommand = cli.Command{
		Action: utils.MigrateFlags(genOwnershipProof),
		Name:   "genOwnershipProof",
		Usage:  "Generate enode proof",
		Flags: []cli.Flag{
			utils.AutonityKeysFileFlag,
			utils.AutonityKeysHexFlag,
			utils.OracleKeyFileFlag,
			utils.OracleKeyHexFlag,
		},
		Description: `
    	autonity genOwnershipProof
		Generates a ownership proof, given a node key file which contains the 
		node key and with or without the node consensus key appending at the file,
		oracle private key file and the treasury address. Proof is printed on 
		stdout in hex string format. This must be copied as it is and passed to 
		registerValidator. Note that, if the consensus key is missing from the node
		key file, a new consensus key will be generated and append to the input node
		key file.
		There are two ways to pass node private key:
			1. --autonitykeys <Autonity keys file name> 
			2. --autonitykeyshex <node keys in hex>
		Similarly there are two ways to pass oracle private key:
			1. --oraclekey <oracle key file name>
			2. --oraclekeyhex <oracle key in hex>`,
		ArgsUsage: "<treasury>",
		Category:  "MISCELLANEOUS COMMANDS",
	}

	genAutonityKeysCommand = cli.Command{
		Action: utils.MigrateFlags(genAutonityKeys),
		Name:   "genAutonityKeys",
		Usage:  "Generate autonity keys",
		Flags: []cli.Flag{
			utils.WriteAddrFlag,
		},
		Description: `
    	autonity genAutonityKeys <outkeyfile>
		Generate node key and its consensus key to the given file. Write out the
		node address, node public key of enode URL and the consensus key of registering
		a validator on stdout using	flag --writeaddress`,
		ArgsUsage: "<outkeyfile>",
		Category:  "MISCELLANEOUS COMMANDS",
	}
)

// makecache generates an ethash verification cache into the provided folder.
func makecache(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) != 2 {
		utils.Fatalf(`Usage: autonity makecache <block number> <outputdir>`)
	}
	block, err := strconv.ParseUint(args[0], 0, 64)
	if err != nil {
		utils.Fatalf("Invalid block number: %v", err)
	}
	ethash.MakeCache(block, args[1])

	return nil
}

// makedag gene
//
//	tes an ethash mining DAG into the provided folder.
func makedag(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) != 2 {
		utils.Fatalf(`Usage: autonity makedag <block number> <outputdir>`)
	}
	block, err := strconv.ParseUint(args[0], 0, 64)
	if err != nil {
		utils.Fatalf("Invalid block number: %v", err)
	}
	ethash.MakeDataset(block, args[1])

	return nil
}

func version(ctx *cli.Context) error {
	fmt.Println(strings.Title(clientIdentifier))
	fmt.Println("Version:", params.VersionWithMeta)
	if gitCommit != "" {
		fmt.Println("Git Commit:", gitCommit)
	}
	if gitDate != "" {
		fmt.Println("Git Commit Date:", gitDate)
	}
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Protocol Versions:", ethproto.ProtocolVersions)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}

func license(_ *cli.Context) error {
	fmt.Println(`Autonity is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Autonity is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with autonity. If not, see <http://www.gnu.org/licenses/>.`)
	return nil
}

// genOwnershipProof generates an ownership proof of the node key, oracle key and the consensus key for a node operator.
// If the input node key file is with a legacy format which missing a consensus key, the function will generate a random
// consensus secret key and append it in the legacy node key file.
func genOwnershipProof(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) != 1 {
		utils.Fatalf(`Usage: autonity genOwnershipProof [options] <treasuryAddress>`)
	}

	var nodePrivateKey, oraclePrivateKey *ecdsa.PrivateKey
	var consensusKey blst.SecretKey
	var err error
	// load node key and consensus key, if the consensus key is missing, it generates new one for legacy node key file.
	if nodeKeyFile := ctx.GlobalString(utils.AutonityKeysFileFlag.Name); nodeKeyFile != "" {
		s, err := os.Stat(nodeKeyFile)
		if err != nil {
			utils.Fatalf("Failed to load the node private key: %v", err)
		}

		// parse and load node key from legacy node key file, generate consensus key and append it.
		if s.Size() < crypto.AutonityKeysLenInChar {
			nodePrivateKey, err = crypto.LoadECDSA(nodeKeyFile)
			if err != nil {
				println("error: ", err.Error())
				utils.Fatalf("Failed to load the node private key: %v", err)
			}

			consensusKey, err = blst.RandKey()
			if err != nil {
				utils.Fatalf("Failed to generate node consensus key: %v", err)
			}

			if err = crypto.SaveAutonityKeys(nodeKeyFile, nodePrivateKey, consensusKey); err != nil {
				utils.Fatalf("Failed to generate node consensus key: %v", err)
			}
		}

		// parse and load node key and consensus key.
		if s.Size() >= crypto.AutonityKeysLenInChar {
			// load key from the node key file
			nodePrivateKey, consensusKey, err = crypto.LoadAutonityKeys(nodeKeyFile)
			if err != nil {
				utils.Fatalf("Failed to load the node private key: %v", err)
			}
		}
	} else if privateKeysHex := ctx.GlobalString(utils.AutonityKeysHexFlag.Name); privateKeysHex != "" {
		// if the consensus key is missing from the input hex string, terminate the execution.
		nodePrivateKey, consensusKey, err = crypto.HexToAutonityKeys(privateKeysHex)
		if err != nil {
			utils.Fatalf("Failed to parse the node private key: %v", err)
		}
	} else {
		utils.Fatalf(`Node key details are not provided`)
	}

	// load oracle node key from file or from input hex string.
	if oracleKeyFile := ctx.GlobalString(utils.OracleKeyFileFlag.Name); oracleKeyFile != "" {
		oraclePrivateKey, err = crypto.LoadECDSA(oracleKeyFile)
		if err != nil {
			utils.Fatalf("Failed to load the oracle private key: %v", err)
		}
	} else if oracleKeyHex := ctx.GlobalString(utils.OracleKeyHexFlag.Name); oracleKeyHex != "" {
		oraclePrivateKey, err = crypto.HexToECDSA(oracleKeyHex)
		if err != nil {
			utils.Fatalf("Failed to parse the oracle private key: %v", err)
		}
	} else {
		utils.Fatalf(`oracle key details are not provided`)
	}

	treasury := args[0]
	signatures, err := crypto.AutonityPOPProof(nodePrivateKey, oraclePrivateKey, treasury, consensusKey)
	if err != nil {
		if err == hexutil.ErrMissingPrefix {
			utils.Fatalf("Failed to decode: hex string without 0x prefix")
		}
		utils.Fatalf("Failed to generate Autonity POP: %v", err)
	}

	fmt.Println(hexutil.Encode(signatures))
	return nil
}

// genAutonityKeys generates a node key, and append its derived BLS private key (the validator key) in the key file.
func genAutonityKeys(ctx *cli.Context) error {
	outKeyFile := ctx.Args().First()
	if len(outKeyFile) == 0 {
		utils.Fatalf("Out key file must be provided!! Usage: autonity genAutonityKeys <outkeyfile> [options]")
	}

	nodeKey, consensusKey, err := crypto.GenAutonityKeys()
	if err != nil {
		utils.Fatalf("could not generate node key %v", err)
	}

	if err = crypto.SaveAutonityKeys(outKeyFile, nodeKey, consensusKey); err != nil {
		utils.Fatalf("could not save key %v", err)
	}

	writeAddr := ctx.GlobalBool(utils.WriteAddrFlag.Name)
	if writeAddr {
		fmt.Printf("Node address: %s\n", crypto.PubkeyToAddress(nodeKey.PublicKey).String())
		fmt.Printf("Node public key: 0x%x\n", crypto.FromECDSAPub(&nodeKey.PublicKey)[1:])
		fmt.Println("Consensus public key:", consensusKey.PublicKey().Hex())
	}
	fmt.Println("Node's validator key:", consensusKey.PublicKey().Hex())
	return nil
}
