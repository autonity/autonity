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
	"encoding/hex"
	"fmt"
	"github.com/autonity/autonity/cmd/utils"
	"github.com/autonity/autonity/crypto"
	"gopkg.in/urfave/cli.v1"
)

type outputAutInspect struct {
	NodeAddress         string
	NodePublicKey       string
	NodePrivateKey      string
	ConsensusPublicKey  string
	ConsensusPrivateKey string
}

const HexPrefix = "0x"

var commandAutInspect = cli.Command{
	Name:      "autinspect",
	Usage:     "inspect an AutonityKeys file",
	ArgsUsage: "<autonityKeys file>",
	Description: `
Print various information about the AutonityKeys file.

Private keys information can be printed by using the --private flag;
make sure to use this feature with great caution!`,
	Flags: []cli.Flag{
		passphraseFlag,
		jsonFlag,
		cli.BoolFlag{
			Name:  "private",
			Usage: "include the private keys in the output",
		},
	},
	Action: func(ctx *cli.Context) error {
		keyfilepath := ctx.Args().First()

		nodeKey, consensusKey, err := crypto.LoadAutonityKeys(keyfilepath)
		if err != nil {
			utils.Fatalf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
		}

		// Output all relevant information we can retrieve.
		showPrivate := ctx.Bool("private")
		out := outputAutInspect{
			NodeAddress:        crypto.PubkeyToAddress(nodeKey.PublicKey).Hex(),
			NodePublicKey:      HexPrefix + hex.EncodeToString(crypto.FromECDSAPub(&nodeKey.PublicKey)),
			ConsensusPublicKey: consensusKey.PublicKey().Hex(),
		}
		if showPrivate {
			out.NodePrivateKey = HexPrefix + hex.EncodeToString(crypto.FromECDSA(nodeKey))
			out.ConsensusPrivateKey = consensusKey.Hex()
		}

		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			fmt.Println("Node address: ", out.NodeAddress)
			fmt.Println("Node public key: ", out.NodePublicKey)
			fmt.Println("Consensus public key: ", out.ConsensusPublicKey)

			if showPrivate {
				fmt.Println("Node private key: ", out.NodePrivateKey)
				fmt.Println("Consensus private key: ", out.ConsensusPrivateKey)
			}
		}
		return nil
	},
}
