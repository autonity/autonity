package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"

	"github.com/google/uuid"
	"gopkg.in/urfave/cli.v1"

	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
)

var (
	app = cli.NewApp()

	configFlag = cli.StringFlag{
		Name:  "config",
		Value: "./config.json",
		Usage: "Configuration for netdiag runner",
	}
	gcpProjectIDFlag = cli.StringFlag{
		Name:  "gcp-project-id",
		Value: "",
		Usage: "GCP project id",
	}
	gcpInstanceTemplateFlag = cli.StringFlag{
		Name:  "gcp-template",
		Value: "",
		Usage: "GCP VM instance template",
	}
	peersFlag = cli.IntFlag{
		Name:  "peers",
		Value: 7,
		Usage: "Number of runner instances to deploy",
	}
	idFlag = cli.IntFlag{
		Name:  "id",
		Value: 0,
		Usage: "Index of the local runner in the configuration",
	}

	setupCommand = cli.Command{
		Action:    setup,
		Name:      "setup",
		Usage:     "Setup a new autonity diagnosis deployment",
		ArgsUsage: "",
		Flags: []cli.Flag{
			peersFlag,
			configFlag,
			gcpProjectIDFlag,
			gcpInstanceTemplateFlag,
		},
		Description: `
The setup command deploys a new network of nodes.`,
	}

	// control command

	// cleanup command

	// run command is to start a runner
	runCommand = cli.Command{
		Action:    setup,
		Name:      "init",
		Usage:     "Start a runner instance.",
		ArgsUsage: "",
		Flags: []cli.Flag{
			configFlag,
			idFlag,
		},
		Description: `
The setup command deploys a new network of nodes.`,
	}
)

var (
	errInvalidInput = errors.New("invalid input provided")
)

func init() {
	app.Name = "NetDiag"
	app.Usage = "Autonity Network Diagnosis Utility"
	app.Flags = []cli.Flag{}

	app.Action = run
	app.Commands = []cli.Command{
		setupCommand,
		// run flag
	}
}

type nodeConfig struct {
	enode string
	key   string
}
type config struct {
	nodes []nodeConfig
}

func readConfigFile(file string) config {
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Error("can't open config file", "err", err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	// Read the file into a byte slice
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Crit("error reading json", "err", err)
	}
	// Unmarshal the byte slice into a Person struct
	conf := &config{}
	if err := json.Unmarshal(byteValue, &conf); err != nil {
		log.Crit("error unmarshalling config", "err", err)
	}
	return *conf
}

func run(c *cli.Context) error {
	log.Info("Runner started")
	// Listen for Ctrl-C.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		for _ = range sigCh {
			return
		}
	}()
	cfg := readConfigFile(c.String(configFlag.Name))
	newEngine(cfg).start()
	// listen here for RPC commands
	return nil
}

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	log.Info("Autonity Network Diagnosis Utility")
	if err := app.Run(os.Args); err != nil {
		log.Error("critical failure", "err", err)
		os.Exit(1)
	}
}

func setup(c *cli.Context) error {
	n := c.Int(peersFlag.Name)
	if n <= 0 {
		fmt.Printf("--%s flag not provided or invalid.\n", peersFlag.Name)
		fmt.Print("How many runners to deploy? ")
		if _, err := fmt.Scan(&n); err != nil {
			return err
		}
		if n <= 0 {
			return errInvalidInput
		}
	}
	instanceTemplate := c.String(gcpInstanceTemplateFlag.Name)
	if instanceTemplate == "" {
		fmt.Printf("--%s flag not provided or invalid.\n", gcpInstanceTemplateFlag.Name)
		fmt.Print("Insert the GCP instance template for VMs: ")
		if _, err := fmt.Scan(&instanceTemplate); err != nil {
			return err
		}
	}
	projectId := c.String(gcpProjectIDFlag.Name)
	if projectId == "" {
		fmt.Printf("--%s flag not provided or invalid.\n", gcpProjectIDFlag.Name)
		fmt.Print("Insert the GCP project id: ")
		if _, err := fmt.Scan(&projectId); err != nil {
			return err
		}
	}
	zones, err := listZones(projectId)
	if err != nil {
		log.Crit("can't retrieve available zone list", "err", err)
	}
	log.Info("Deploying new runner network", "count", n, "template", instanceTemplate, "project-id", projectId)
	vms := make([]*vm, n)
	var wg sync.WaitGroup
	// create VM instances on GCP
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			var err error
			name := "netdiag-runner-" + uuid.New().String()
			vms[id], err = deployVM(projectId, name, zones[id%len(zones)], instanceTemplate)
			wg.Done()
			if err != nil {
				log.Crit("error deploying VM", "id", id, "err", err)
			}
		}(i)
	}
	wg.Wait()
	// generate keys and enodes
	cfg := config{nodes: make([]nodeConfig, n)}
	for i := 0; i < n; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			return err
		}
		cfg.nodes[i] = nodeConfig{
			key:   hex.EncodeToString(crypto.FromECDSA(key)),
			enode: fmt.Sprintf("enode://%x:", crypto.FromECDSAPub(&key.PublicKey)[1:]),
		}
	}

	configFile, err := os.Create(c.String(configFlag.Name))
	if err != nil {
		wd, _ := os.Getwd()
		log.Error("can't create config file", "err", err, "file", c.String(configFlag.Name), "wd", wd)
		return err
	}

	defer configFile.Close() // Ensure the file is closed when the function exits

	if err := json.NewEncoder(configFile).Encode(cfg); err != nil {
		log.Error("can't encode config to file", "err", err)
		return err
	}

	// deploy runners
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			if err := vms[id].deployRunner(c.String(configFlag.Name)); err != nil {
				log.Crit("error deploying runner", "id", id, "err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// start runners
	for i := 0; i < n; i++ {
		go func(id int) {
			wg.Add(1)
			if err := vms[id].startRunner(); err != nil {
				log.Crit("error deploying runner", "id", id, "err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return nil
}
