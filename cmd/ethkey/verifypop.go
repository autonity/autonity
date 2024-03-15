package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"

	"github.com/autonity/autonity/accounts"
	"github.com/autonity/autonity/cmd/utils"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	enodes "github.com/autonity/autonity/p2p/enode"
)

type outputVerifyPOP struct {
	NodeKeyPOP      bool
	OracleKeyPOP    bool
	ConsensusKeyPOP bool
}

var commandVerifyPOP = cli.Command{
	Name:      "verifypop",
	Usage:     "verify proof of possessions",
	ArgsUsage: "<treasury> <enode> <oracle> <consensusKey> <proof>",
	Description: `
Print if the proof of possession for registering a validator is valid or not. Use --json for a json output.`,
	Flags: []cli.Flag{
		jsonFlag,
	},
	Action: func(ctx *cli.Context) error {
		// Checks arguments BEGIN
		if ctx.NArg() != 5 {
			fmt.Println("Incorrect number of arguments, 5 expected.")
			cli.ShowCommandHelpAndExit(ctx, "verifypop", 1)
		}
		treasuryArg := ctx.Args().Get(0)
		treasuryKey, err := hexutil.Decode(treasuryArg)
		if err != nil {
			utils.Fatalf("can't decode treasury key with error %v", err)
		}

		enodeArg := ctx.Args().Get(1)
		node, err := enodes.ParseV4NoResolve(enodeArg)
		if err != nil {
			utils.Fatalf("can't parse enode %v", err)
		}
		oracleArg := ctx.Args().Get(2)
		oracleKey, err := hexutil.Decode(oracleArg)
		if err != nil {
			utils.Fatalf("can't decode oracle key with error %v", err)
		}

		consensusKeyArg := ctx.Args().Get(3)
		consensusKey, err := hexutil.Decode(consensusKeyArg)
		if err != nil {
			utils.Fatalf("can't decode consensus key with error %v", err)
		}

		signaturesArg := ctx.Args().Get(4)
		signatures, err := hexutil.Decode(signaturesArg)
		if err != nil {
			utils.Fatalf("can't decode proof with error %v", err)
		}
		if len(signatures) != crypto.AutonityPOPLen {
			utils.Fatalf("invalid proof size")
		}
		// Checks arguments END
		out := outputVerifyPOP{
			NodeKeyPOP:      false,
			OracleKeyPOP:    false,
			ConsensusKeyPOP: false,
		}
		// compute the hash (the signed message)
		hash, _ := accounts.TextAndHash(treasuryKey)
		//
		// Step 1: node key pop
		//
		recoveredNodeKey, err := crypto.SigToPub(hash, signatures[:crypto.SignatureLength])
		if err != nil {
			utils.Fatalf("can't recover node key %v", err)
		}
		out.NodeKeyPOP = recoveredNodeKey.Equal(node.Pubkey())
		//
		// Step 2: oracle key pop
		//
		recoveredOracleAddress, err := crypto.SigToAddr(hash, signatures[crypto.SignatureLength:2*crypto.SignatureLength])
		if err != nil {
			utils.Fatalf("can't recover oracle key %v", err)
		}
		out.OracleKeyPOP = recoveredOracleAddress == common.BytesToAddress(oracleKey)

		//
		// Step 3: consensus\ key pop
		//
		key, err := blst.PublicKeyFromBytes(consensusKey)
		if err != nil {
			utils.Fatalf("invalid consensus key %v", err)
		}
		sig, err := blst.SignatureFromBytes(signatures[2*crypto.SignatureLength:])
		if err != nil {
			utils.Fatalf("can't recover signature key %v", err)
		}
		if sig.IsZero() {
			utils.Fatalf("can't recover signature key %v", err)
		}
		err = crypto.BLSPOPVerify(key, sig, treasuryKey)
		if err == nil {
			out.ConsensusKeyPOP = true
		}

		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			fmt.Println("===== POP Validation Results =====")
			fmt.Println("Node Key:      ", out.NodeKeyPOP)
			fmt.Println("Oracle Key:    ", out.OracleKeyPOP)
			fmt.Println("Consensus Key: ", out.ConsensusKeyPOP)
		}
		return nil
	},
}
