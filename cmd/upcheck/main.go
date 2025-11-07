package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/utils"
	"github.com/autonity/autonity/common"
	"github.com/kr/text"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v3"
)

var (
	showSourceFlag = &cli.BoolFlag{
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

// substitues big equalities diffs (> 6 lines) with:
// 3 lines of context
// [...]
// 3 lines of context
func dryUp(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	driedUpDiffs := make([]diffmatchpatch.Diff, 0, len(diffs))

	for _, diff := range diffs {

		// dry up only equalities
		if diff.Type != diffmatchpatch.DiffEqual {
			driedUpDiffs = append(driedUpDiffs, diff)
			continue
		}

		// count total number of lines.
		// kinda ugly but it works and no need to deal with
		// edge cases such as "does the last line end with \n or not"
		n := 0
		for range strings.Lines(diff.Text) {
			n++
		}

		// small equal diffs are left untouched
		if n <= 6 {
			driedUpDiffs = append(driedUpDiffs, diff)
			continue
		}

		// long equality diff here. dry it up
		var firsts string
		var lasts string
		i := 0
		for line := range strings.Lines(diff.Text) {
			if i < 3 {
				firsts += line
			}
			if i >= n-3 {
				lasts += line
			}
			i++
		}

		driedUpDiffs = append(driedUpDiffs, diffmatchpatch.Diff{
			Type: diff.Type, // == DiffEqual
			Text: firsts + "[...]\n" + lasts,
		})
	}
	return driedUpDiffs
}

// prints a line-mode diff
// see https://github.com/sergi/go-diff/issues/69#issuecomment-688602689
// and https://github.com/google/diff-match-patch/wiki/Line-or-Word-Diffs
func printCodeDiff(w io.Writer, base string, upgraded string) {
	dmp := diffmatchpatch.New()

	baseReduced, upgradedReduced, lines := dmp.DiffLinesToChars(base, upgraded)
	diffs := dmp.DiffMain(baseReduced, upgradedReduced, false)
	diffs = dmp.DiffCharsToLines(diffs, lines)
	diffs = dmp.DiffCleanupSemantic(diffs)

	diffs = dryUp(diffs)

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
		fmt.Fprintf(w, " \t%s\n", contract.Target.String())
		fmt.Fprintf(w, " \t \t%s\t%s\n", boldify("address"), contract.Target.Address().String())
		fmt.Fprintf(w, " \t \t%s\t%s\n", boldify("hash"), contract.Hash.String())
		fmt.Fprintf(w, " \t \t%s\t%s\n", boldify("version"), contract.VersionString)
		if showSourceFlag.IsSet() {
			fmt.Fprintf(w, "%s\n\n", boldify("source code diff:"))
			printCodeDiff(w, contract.BaseCode, contract.UpgradedCode)
		}
	}

	return nil
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

	calldata, err := autonity.Upgrades[upgradeNumber].Calldata()
	if err != nil {
		return fmt.Errorf("cannot assemble calldata for upgrade %d: %w", upgradeNumber, err)
	}

	fmt.Print(common.Bytes2Hex(calldata))

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
