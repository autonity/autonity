package main

import (
	"encoding/hex"
	"fmt"

	"gopkg.in/urfave/cli.v1"

	"github.com/autonity/autonity/cmd/utils"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto"
)

type outputAutInspect struct {
	NodePublicKey       string
	NodeAddress         string
	NodePrivateKey      string `json:",omitempty"`
	ConsensusPublicKey  string
	ConsensusPrivateKey string `json:",omitempty"`
}

var commandAutInspect = cli.Command{
	Name:      "autinspect",
	Usage:     "inspect autonity keys file",
	ArgsUsage: "<keyfile>",
	Description: `
Print various information about the autonity keyfile.

Private key information can be printed by using the --private flag;
make sure to use this feature with great caution!`,
	Flags: []cli.Flag{
		jsonFlag,
		cli.BoolFlag{
			Name:  "private",
			Usage: "include the private keys in the output",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() != 1 {
			fmt.Println("Incorrect number of arguments, 1 expected.")
			cli.ShowCommandHelpAndExit(ctx, "autinspect", 1)
		}
		nodeKey, consensusKey, err := crypto.LoadAutonityKeys(ctx.Args().Get(0))
		if err != nil {
			utils.Fatalf("error opening key file %v", err)
		}
		out := outputAutInspect{
			NodePublicKey:      hexutil.Encode(crypto.FromECDSAPub(&nodeKey.PublicKey)[1:]),
			NodeAddress:        crypto.PubkeyToAddress(nodeKey.PublicKey).String(),
			ConsensusPublicKey: consensusKey.PublicKey().Hex(),
		}
		showPrivate := ctx.Bool("private")
		if showPrivate {
			out.ConsensusPrivateKey = hex.EncodeToString(consensusKey.Marshal())
			out.NodePrivateKey = hex.EncodeToString(crypto.FromECDSA(nodeKey))
		}
		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			fmt.Println("Node Address:          ", out.NodeAddress)
			fmt.Println("Node Public Key:       ", out.NodePublicKey)
			fmt.Println("Consensus Public Key:  ", out.ConsensusPublicKey)
			if showPrivate {
				fmt.Println("Consensus Private Key: ", out.ConsensusPrivateKey)
				fmt.Println("Node Private Key:      ", out.NodePrivateKey)
			}
		}
		return nil
	},
}
