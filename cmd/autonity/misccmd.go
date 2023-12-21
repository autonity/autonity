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
			utils.NodeKeyFileFlag,
			utils.NodeKeyHexFlag,
			utils.OracleKeyFileFlag,
			utils.OracleKeyHexFlag,
		},
		Description: `
    	autonity genOwnershipProof
		Generates a proof, given a node private key, oracle private key and the 
		treasury address. Proof is printed on stdout in hex string format. This 
		must be copied as it is and passed to registerValidator.
		There are two ways to pass node private key:
			1. --nodekey <node key file name> 
			2. --nodekeyhex <node key in hex>
		Similarly there are two ways to pass oracle private key:
			1. --oraclekey <oracle key file name>
			2. --oraclekeyhex <oracle key in hex>`,
		ArgsUsage: "<treasury>",
		Category:  "MISCELLANEOUS COMMANDS",
	}

	genKeyCommand = cli.Command{
		Action: utils.MigrateFlags(genNodeKey),
		Name:   "genNodeKey",
		Usage:  "Generate node key",
		Flags: []cli.Flag{
			utils.WriteAddrFlag,
		},
		Description: `
    	autonity genNodeKey <outkeyfile>
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

// genOwnershipProof generates an ownership proof of the node, oracle account and the consensus key. It assumes that
// the input node key is generated by the CLI sub command: genNodeKey. Otherwise, it will fail due to the consensus
// secret key is missing from the input node key file or from the input hex string of node key.
// todo: (Jason) do we need a forward compatibility for legacy node key format which does not have a consensus secret
//
//	key appending? In this case, we can generate a random consensus secret key and append it in the legacy node key file.
func genOwnershipProof(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) != 1 {
		utils.Fatalf(`Usage: autonity genOwnershipProof [options] <treasuryAddress>`)
	}

	// Load private key
	var nodePrivateKey, oraclePrivateKey *ecdsa.PrivateKey
	var consensusKey blst.SecretKey
	var err error
	if nodeKeyFile := ctx.GlobalString(utils.NodeKeyFileFlag.Name); nodeKeyFile != "" {
		// load key from the node key file
		nodePrivateKey, consensusKey, err = crypto.LoadNodeKey(nodeKeyFile)
		if err != nil {
			utils.Fatalf("Failed to load the node private key: %v", err)
		}
	} else if privateKeysHex := ctx.GlobalString(utils.NodeKeyHexFlag.Name); privateKeysHex != "" {
		nodePrivateKey, consensusKey, err = crypto.HexToNodeKey(privateKeysHex)
		if err != nil {
			utils.Fatalf("Failed to parse the node private key: %v", err)
		}
	} else {
		utils.Fatalf(`Node key details are not provided`)
	}

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

	fmt.Println("Node consensus key:", consensusKey.PublicKey().Hex())
	fmt.Println("Signatures:", hexutil.Encode(signatures))
	return nil
}

// genNodeKey generates a node key, and append its derived BLS private key (the validator key) in the key file.
func genNodeKey(ctx *cli.Context) error {
	outKeyFile := ctx.Args().First()
	if len(outKeyFile) == 0 {
		utils.Fatalf("Out key file must be provided!! Usage: autonity genNodeKey <outkeyfile> [options]")
	}

	nodeKey, consensusKey, err := crypto.GenAutonityNodeKey()
	if err != nil {
		utils.Fatalf("could not generate node key %v", err)
	}

	if err = crypto.SaveNodeKey(outKeyFile, nodeKey, consensusKey); err != nil {
		utils.Fatalf("could not save key %v", err)
	}

	writeAddr := ctx.GlobalBool(utils.WriteAddrFlag.Name)
	if writeAddr {
		fmt.Printf("Node address: %s\n", crypto.PubkeyToAddress(nodeKey.PublicKey).String())
		fmt.Printf("Node public key: 0x%x\n", crypto.FromECDSAPub(&nodeKey.PublicKey)[1:])
		fmt.Println("Node consensus key:", consensusKey.PublicKey().Hex())
	}
	return nil
}
