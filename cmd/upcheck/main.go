package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/utils"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params/generated"
	"github.com/kr/text"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v3"
)

var (
	upgradeManagerABI = generated.UpgradeManager1Abi
	showSourceFlag    = &cli.BoolFlag{
		Name:    "show-source",
		Aliases: []string{"s"},
		Usage:   "show also source code diff",
	}
)

func boldify(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func listUpgrades(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 0 {
		return fmt.Errorf("list command takes no arguments")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 16, 0, '\t', 0)
	defer w.Flush()

	for i, upgrade := range autonity.Upgrades {
		fmt.Fprintf(w, "%s\t%s\n", boldify(fmt.Sprintf("upgrade %d", i)), upgrade.Name)
	}
	return nil
}

// TODO: this can be improved to show only the actual diffs and some context around it
func printCodeDiff(w io.Writer, base string, upgraded string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(base, upgraded, false)
	fmt.Fprintf(w, "%s", dmp.DiffPrettyText(diffs))
}

func detailUpgrade(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return fmt.Errorf("usage: upgrade <number>")
	}

	upgradeNumber, err := strconv.ParseUint(cmd.Args().First(), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse upgrade number: %w", err)
	}

	if upgradeNumber >= uint64(len(autonity.Upgrades)) {
		return fmt.Errorf("upgrade number out of bound, there are only %d upgrades available", len(autonity.Upgrades))
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 16, 0, '\t', 0)
	defer w.Flush()

	upgrade := autonity.Upgrades[upgradeNumber]
	fmt.Fprintf(w, "%s\t%s\n", boldify(fmt.Sprintf("upgrade %d", upgradeNumber)), upgrade.Name)
	fmt.Fprintf(w, "%s%s\n", boldify("description"), text.Indent(text.Wrap(upgrade.Description, 80), "\t"))
	fmt.Fprintf(w, "%s\t%v\n\n", boldify("excluded"), upgrade.ExclusionList)
	fmt.Fprintf(w, "%s\n", boldify(fmt.Sprintf("upgraded contracts (%d)", len(upgrade.Upgrades))))
	for _, contract := range upgrade.Upgrades {
		fmt.Fprintf(w, " \t%s (%s)\n", contract.Target.String(), contract.Target.Address().String())
		if showSourceFlag.IsSet() {
			fmt.Fprintf(w, "%s\n\n", boldify("source code diff:"))
			printCodeDiff(w, contract.BaseCode, contract.UpgradedCode)
		}
	}

	return nil
}

// assume indexes will not result in out of bound access
func upgradePayload(upgradeNumber uint64, contractNumber uint64) ([]byte, error) {
	contractUpgrade := autonity.Upgrades[upgradeNumber].Upgrades[contractNumber]

	// if no args are specified, `constructorArgs` will be == []
	constructorArgs, err := contractUpgrade.Abi.Pack("", contractUpgrade.Args...)
	if err != nil {
		return nil, fmt.Errorf("cannot pack upgrade args: %w", err)
	}
	return append(contractUpgrade.Bytecode, constructorArgs...), nil
}

func assemble(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return fmt.Errorf("usage: upgrade <number>")
	}

	upgradeNumber, err := strconv.ParseUint(cmd.Args().First(), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse upgrade number: %w", err)
	}

	if upgradeNumber >= uint64(len(autonity.Upgrades)) {
		return fmt.Errorf("upgrade number out of bound, there are only %d upgrades available", len(autonity.Upgrades))
	}

	upgrade := autonity.Upgrades[upgradeNumber]

	// check if 1 or more contracts need to be updated
	if len(upgrade.Upgrades) == 1 {
		payload, err := upgradePayload(upgradeNumber, 0)
		if err != nil {
			return fmt.Errorf("cannot build upgrade payload for upgrade %d - contract %d: %w", upgradeNumber, 0, err)
		}
		calldata, err := upgradeManagerABI.Pack("upgrade", upgrade.Upgrades[0].Target, string(payload))
		if err != nil {
			return fmt.Errorf("cannot build calldata for upgrade %d - contract %d: %w", upgradeNumber, 0, err)
		}
		fmt.Println(common.Bytes2Hex(calldata))
		return nil
	}

	// more than 1 contract needs to be upgraded
	// TODO: implement multiple contract upgrade
	//       - without version tag
	//       - with version tag

	return nil
}

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "list",
				Usage:   "list of available upgrades",
				Aliases: []string{"l"},
				Action:  listUpgrades,
			},
			{
				Name:      "detail",
				Usage:     "see the details of one upgrade",
				Aliases:   []string{"d"},
				Flags:     []cli.Flag{showSourceFlag},
				ArgsUsage: "<number>",
				Action:    detailUpgrade,
			},
			{
				Name:      "assemble",
				Aliases:   []string{"a"},
				Usage:     "assemble tx calldata for upgrades",
				ArgsUsage: "<number>",
				Action:    assemble,
			},
		},
		Name:  "upcheck",
		Usage: "contract upgrade verification tool",
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		utils.Fatalf("error %v", err)
	}
}
