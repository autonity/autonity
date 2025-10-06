package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/cmd/utils"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/internal/flags"
	"gopkg.in/urfave/cli.v1"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	gitDate   = ""

	solcBin   = "solc_static_linux_v0.8.30"
	solcFlags = []string{"--overwrite", "--optimize", "--optimize-runs", "10000", "--evm-version", "london", "--abi", "--bin", "--userdoc", "--devdoc", "-o"}

	app *cli.App
)

func init() {
	app = flags.NewApp(gitCommit, gitDate, "contract upgrade verification tool")
	app.Flags = []cli.Flag{}
	app.Action = utils.MigrateFlags(upcheck)
	cli.CommandHelpTemplate = flags.OriginCommandHelpTemplate
}

func upcheck(c *cli.Context) error {
	if len(c.Args()) != 2 {
		utils.Fatalf("invalid arguments. usage ./upcheck [CONTRACT] [ADDRESS]")
	}
	fmt.Fprintf(os.Stderr, "** Upgrade Verification Tool **\n")
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	path := c.Args().First()
	fmt.Fprintf(os.Stderr, "Contract: %s \n", c.Args()[0])

	// 1st step: Retrieve contrat's deployment bytecode
	bin := runSolc(exPath, path)

	// 2nd step: Construct ABI-encoded deployment calldata
	// for now, only simple constructors are supported with this tool

	abi, err := abi.JSON(strings.NewReader(`[{ "type" : "function", "name" : ""}]`))
	if err != nil {
		panic(err)
	}
	packedArgs, err := abi.Pack("")
	if err != nil {
		panic(err)
	}
	packed := append(bin, packedArgs...)
	// fmt.Println("Packed Deployment Bytecode:", common.Bytes2Hex(packed))
	// 3rd step: Construct ABI-encoded upgrade transaction
	upgraderAbi, err := bindings.UpgradeManagerMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	finalPacked, err := upgraderAbi.Pack("upgrade", common.HexToAddress(c.Args()[1]), string(packed))
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stderr, "Upgrade TX Calldata:")
	fmt.Println(common.Bytes2Hex(finalPacked))
	return nil
}

func runSolc(self, path string) []byte {
	solcPath := strings.Join([]string{self, solcBin}, "/")
	outputDir, err := os.MkdirTemp("", "upcheck")
	if err != nil {
		panic(err)
	}
	cmd := &exec.Cmd{
		Path:   solcPath,
		Args:   append(solcFlags, outputDir, path),
		Stdout: os.Stderr,
		Stderr: os.Stderr,
	}
	fmt.Fprintf(os.Stderr, "Running solc (%s) ...\n ....", solcPath)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	//Retrieving correct binary
	solidityFile := filepath.Base(path)
	solidityContract := strings.TrimSuffix(solidityFile, filepath.Ext(path))

	//
	bytecode, err := os.ReadFile(fmt.Sprintf("%s/%s.bin", outputDir, solidityContract))
	if err != nil {
		panic(err)
	}
	return common.Hex2Bytes(string(bytecode))
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
