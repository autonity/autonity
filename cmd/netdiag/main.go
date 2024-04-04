package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"os/user"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
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
	gcpUsername = cli.StringFlag{
		Name:  "gcp-username",
		Value: "root",
		Usage: "Username to access gcp instances",
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
		Usage:     "Setup a new netdiag runner cluster",
		ArgsUsage: "",
		Flags: []cli.Flag{
			peersFlag,
			configFlag,
			gcpProjectIDFlag,
			gcpInstanceTemplateFlag,
			gcpUsername,
		},
		Description: `
The setup command deploys a new network of nodes.`,
	}
	// update command
	updateCommand = cli.Command{
		Action:    update,
		Name:      "update",
		Usage:     "Update a deployed cluster",
		ArgsUsage: "",
		Flags: []cli.Flag{
			configFlag,
		},
		Description: `
Update a cluster with current binary.`,
	}
	// control command
	controlCommand = cli.Command{
		Action:    control,
		Name:      "control",
		Usage:     "Control a network via rpc.",
		ArgsUsage: "",
		Flags: []cli.Flag{
			configFlag,
		},
		Description: `
The control command starts the netdiag command center.`,
	}
	// cleanup command

	// run command is to start a runner
	runCommand = cli.Command{
		Action:    run,
		Name:      "run",
		Usage:     "Start a runner instance.",
		ArgsUsage: "",
		Flags: []cli.Flag{
			configFlag,
			idFlag,
		},
		Description: `
The run command start a local runner`,
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
		runCommand,
		controlCommand,
		updateCommand,
		// run flag
	}
}

type nodeConfig struct {
	Enode        string
	Ip           string
	Key          string
	Zone         string
	InstanceName string
}
type config struct {
	Nodes        []nodeConfig
	GcpProjectId string
	GcpUsername  string
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

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	log.Info("==========================================")
	log.Info("=== Autonity Network Diagnosis Utility ===")
	log.Info("==========================================")
	if err := app.Run(os.Args); err != nil {
		log.Error("critical failure", "err", err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	user, _ := user.Current()
	fmt.Printf("Username: %s\n", user.Username)
	localId := c.Int(idFlag.Name)
	log.Info("Runner started", "cmd", strings.Join(os.Args, " "), "id", localId, "user", user.Username, "uid", user.Uid)
	// Listen for Ctrl-C.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	cfg := readConfigFile(c.String(configFlag.Name))
	key := cfg.Nodes[localId].Key
	skey, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Crit("can't load key", "key", key)
	}
	engine := newEngine(cfg, skey)
	if err := rpc.Register(&P2POp{engine}); err != nil {
		log.Error("can't register RPC", "err", err)
		os.Exit(1)
	}
	ln, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Error("listen error:", "err", err)
	}
	log.Info("listening rpc on port 1337")
	go rpc.Accept(ln)
	if err := engine.start(); err != nil {
		log.Error("engine start error", "err", err)
	}
	log.Info("P2P server started")
	// Block and wait for interrupt signal.
	<-sigCh
	log.Info("Shutdown signal received, exiting...")
	return nil
}

func control(c *cli.Context) error {
	// This is very ugly and need to be refactored :(
	targetPeer := 0
	cfg := readConfigFile(c.String(configFlag.Name))
	client, err := rpc.Dial("tcp", cfg.Nodes[targetPeer].Ip+":1337")
	if err != nil {
		log.Error("Dialing error", "err", err)
		return err
	}
	log.Info("Connected!", "ip", cfg.Nodes[targetPeer].Ip)
	reader := bufio.NewReader(os.Stdin)
	p := &P2POp{}
	typeName := reflect.TypeOf(p).Elem().Name()
	methods := reflect.TypeOf(p)

	peerRegexCmd := regexp.MustCompile(`^p(\d+)$`)

	for {
		// List available methods
		fmt.Println("Available commands:")
		for i := 0; i < methods.NumMethod(); i++ {
			method := methods.Method(i)
			fmt.Printf("[%d] %s\n", i+1, method.Name)
		}

		// User selects a method
		fmt.Printf("\n%s(%d)>> ", cfg.Nodes[targetPeer].Ip, targetPeer)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		matches := peerRegexCmd.FindStringSubmatch(input)
		if len(matches) == 2 {
			var err error
			if targetPeer, err = strconv.Atoi(matches[1]); err != nil {
				fmt.Printf("Invalid peer")
				return err
			}
			if client, err = rpc.Dial("tcp", cfg.Nodes[targetPeer].Ip+":1337"); err != nil {
				log.Error("Dialing error", "err", err)
				return err
			}
			fmt.Printf("Connected to peer %d\n", targetPeer)
			continue
		}

		methodIndex, err := strconv.Atoi(input)
		if err != nil || methodIndex < 1 || methodIndex > methods.NumMethod() {
			fmt.Printf("Invalid method selection.")
			return errInvalidInput
		}
		method := methods.Method(methodIndex - 1)
		argType := method.Func.Type().In(1) // Assuming first is receiver, second is context (if present)
		args := reflect.New(argType.Elem()).Interface()
		if userArg, ok := args.(Argument); ok {
			if err := userArg.AskUserInput(); err != nil {
				return err
			}
		}

		replyType := method.Func.Type().In(2)
		var reply reflect.Value
		if replyType.Kind() == reflect.Ptr {
			reply = reflect.New(replyType.Elem())
		} else {
			reply = reflect.New(replyType)
		}

		err = client.Call(typeName+"."+method.Name, args, reply.Interface())
		if err != nil {
			fmt.Printf("RPC call failed: %s\n", err)
			return err
		}

		fmt.Printf("%v", reply.Interface())
		fmt.Printf("----------------------------------------\n")
	}
	return nil
}

func setup(c *cli.Context) error {
	log.Info("New network setup")
	configFileName := c.String(configFlag.Name)
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
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		log.Error("NewInstancesRESTClient: %v", err)
		return err
	}
	vms := make([]*vm, n)
	var wg sync.WaitGroup
	// create VM instances on GCP
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			var err error
			name := "netdiag-runner-" + uuid.New().String()
			vms[id], err = deployVM(ctx, instancesClient, id, projectId, name, zones[id%len(zones)], instanceTemplate, c.String(gcpUsername.Name))
			wg.Done()
			if err != nil {
				log.Crit("error deploying VM", "id", id, "err", err)
			}
		}(i)
	}
	wg.Wait()
	log.Info("Instances deployment completed")
	// generate keys and enodes
	cfg := config{
		Nodes:        make([]nodeConfig, n),
		GcpProjectId: projectId,
		GcpUsername:  c.String(gcpUsername.Name),
	}
	for i := 0; i < n; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			return err
		}
		cfg.Nodes[i] = nodeConfig{
			Ip:           vms[i].ip,
			Key:          hex.EncodeToString(crypto.FromECDSA(key)),
			Enode:        fmt.Sprintf("enode://%x@%s:31337", crypto.FromECDSAPub(&key.PublicKey)[1:], vms[i].ip),
			Zone:         vms[i].zone,
			InstanceName: vms[i].instanceName,
		}
	}

	configFile, err := os.Create(configFileName)
	if err != nil {
		wd, _ := os.Getwd()
		log.Error("can't create config file", "err", err, "file", configFileName, "wd", wd)
		return err
	}

	defer configFile.Close() // Ensure the file is closed when the function exits
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(cfg); err != nil {
		log.Error("can't encode config to file", "err", err)
		return err
	}
	log.Info("Config generated")
	log.Info("Waiting 1 min")
	time.Sleep(60 * time.Second)
	// deploy runners
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			if err := vms[id].deployRunner(configFileName, false, false); err != nil {
				log.Crit("error deploying runner", "id", id, "err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// start runners
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			if err := vms[id].startRunner(configFileName); err != nil {
				log.Crit("error starting runner", "id", id, "err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("Finished!")
	for i := range vms {
		log.Info("Netdiag runner deployed", "id", i, "ip", vms[i].ip)
	}
	return nil
}

func update(c *cli.Context) error {
	log.Info("Updating cluster")
	cfg := readConfigFile(c.String(configFlag.Name))
	vms := make([]*vm, len(cfg.Nodes))
	var wg sync.WaitGroup
	now := time.Now()
	for i, n := range cfg.Nodes {
		wg.Add(1)
		vms[i] = newVM(i, n.Ip, n.InstanceName, n.Zone, cfg.GcpUsername)
		go func(id int) {
			log.Info("Killing runner", "id", id)
			if err := vms[id].killRunner(c.String(configFlag.Name)); err != nil {
				log.Crit("error starting runner", "id", id, "err", err)
			}
			if err := vms[id].deployRunner(c.String(configFlag.Name), true, true); err != nil {
				log.Crit("error deploying runner", "id", id, "err", err)
			}
			log.Info("Runner binary deployed", "id", id)
			if err := vms[id].startRunner(c.String(configFlag.Name)); err != nil {
				log.Crit("error starting runner", "id", id, "err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("Cluster deployed and running!", "duration(s)", time.Now().Sub(now).Seconds())
	return nil
}
