package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"gopkg.in/urfave/cli.v1"

	"github.com/autonity/autonity/cmd/netdiag/api"
	"github.com/autonity/autonity/cmd/netdiag/core"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/log"
)

var (
	app = cli.NewApp()

	configFlag = cli.StringFlag{
		Name:  "Config",
		Value: "./Config.json",
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
	gcpUsernameFlag = cli.StringFlag{
		Name:  "gcp-username",
		Value: "root",
		Usage: "Username to access gcp instances",
	}
	gcpVmName = cli.StringFlag{
		Name:  "vm-name",
		Value: "",
		Usage: "vm instance name",
	}
	networkModeFlag = cli.StringFlag{
		Name:  "network",
		Value: "tcp",
		Usage: "Network type tcp/udp",
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
	restartOnlyFlag = cli.BoolFlag{
		Name: "restart-only",
	}
	pprofFlag = cli.BoolFlag{
		Name: "pprof",
	}

	logDownloadFlag = cli.BoolFlag{
		Name: "log-download",
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
			gcpUsernameFlag,
			networkModeFlag,
			gcpVmName,
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
			networkModeFlag,
			restartOnlyFlag,
			pprofFlag,
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
			idFlag,
		},
		Description: `
The control command starts the netdiag command center.`,
	}
	// cleanup command
	// run command is to start a runner
	cleanupCommand = cli.Command{
		Action:    cleanup,
		Name:      "cleanup",
		Usage:     "clean up runners",
		ArgsUsage: "",
		Flags: []cli.Flag{
			configFlag,
			logDownloadFlag,
		},
		Description: `
The cleanup command deletes all runners`,
	}

	// run command is to start a runner
	runCommand = cli.Command{
		Action:    run,
		Name:      "run",
		Usage:     "Start a runner instance.",
		ArgsUsage: "",
		Flags: []cli.Flag{
			configFlag,
			idFlag,
			networkModeFlag,
			pprofFlag,
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
		cleanupCommand,
		// run flag
	}
}

func readConfigFile(file string) core.Config {
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Error("can't open Config file", "err", err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	// Read the file into a byte slice
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Crit("error reading json", "err", err)
	}
	// Unmarshal the byte slice into a Person struct
	conf := &core.Config{}
	if err := json.Unmarshal(byteValue, &conf); err != nil {
		log.Crit("error unmarshalling Config", "err", err)
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
	if c.Bool(pprofFlag.Name) {
		f, err := os.Create("./profile.pprof")
		if err != nil {
			log.Crit("error pprof", err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	user, _ := user.Current()
	fmt.Printf("Username: %s\n", user.Username)
	localId := c.Int(idFlag.Name)
	log.Info("Runner started", "cmd", strings.Join(os.Args, " "), "id", localId, "user", user.Username, "uid", user.Uid, "network", c.String(networkModeFlag.Name))
	// Listen for Ctrl-C.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)

	cfg := readConfigFile(c.String(configFlag.Name))
	key := cfg.Nodes[localId].Key
	skey, err := crypto.HexToECDSA(key)
	if err != nil {
		log.Crit("can't load key", "key", key)
	}
	engine := core.NewEngine(cfg, localId, skey, c.String(networkModeFlag.Name))
	if err := rpc.Register(&api.P2POp{Engine: engine}); err != nil {
		log.Error("can't register RPC", "err", err)
		os.Exit(1)
	}
	ln, err := net.Listen("tcp", ":1337")
	if err != nil {
		log.Error("listen error:", "err", err)
	}
	log.Info("listening rpc on port 1337")
	go rpc.Accept(ln)
	if err := engine.Start(); err != nil {
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
	targetPeer := c.Int(idFlag.Name)
	cfg := readConfigFile(c.String(configFlag.Name))
	client, err := rpc.Dial("tcp", cfg.Nodes[targetPeer].Ip+":1337")
	if err != nil {
		log.Error("Dialing error", "err", err)
		return err
	}
	log.Info("Connected!", "ip", cfg.Nodes[targetPeer].Ip)
	reader := bufio.NewReader(os.Stdin)
	p := &api.P2POp{}
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
		fmt.Printf("\n%s|%s(%d)>> ", cfg.Nodes[targetPeer].Ip, cfg.Nodes[targetPeer].Zone, targetPeer)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Error("Error reading input", "err", err)
			return err
		}
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
		if userArg, ok := args.(api.Argument); ok {
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
	networkMode := c.String(networkModeFlag.Name)

	configFileName := c.String(configFlag.Name)
	if _, err := os.Stat(configFileName); err == nil {
		return fmt.Errorf("Config file:%s exists, cleanup and retry setup", configFileName)
	}
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
	temp, err := getInstanceTemplate(projectId, path.Base(instanceTemplate))
	if err != nil {
		log.Crit("can't retrieve templates", "err", err)
	}
	zones = filterZones(zones, temp.GetProperties().GetMachineType())
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
			name := "netdiag-runner-" + c.String(gcpVmName.Name) + "-" + uuid.New().String()
			attempts := 0
			for {
				zone := zones[(6*id+attempts)%len(zones)]
				vms[id], err = deployVM(ctx, instancesClient, id, projectId, name, zone.GetName(), instanceTemplate, c.String(gcpUsernameFlag.Name))
				if err == nil || attempts == 10 {
					break
				}
				attempts++
			}
			if err != nil {
				log.Error("error deploying VM after 10 attempts", "id", id, "err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("Instances deployment completed")
	// generate keys and enodes
	cfg := core.Config{
		Nodes:        make([]core.NodeConfig, n),
		GcpProjectId: projectId,
		GcpUsername:  c.String(gcpUsernameFlag.Name),
	}
	for i := 0; i < n; i++ {
		if vms[i] == nil {
			continue
		}
		key, err := crypto.GenerateKey()
		if err != nil {
			return err
		}
		cfg.Nodes[i] = core.NodeConfig{
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
		log.Error("can't create Config file", "err", err, "file", configFileName, "wd", wd)
		return err
	}

	defer configFile.Close() // Ensure the file is closed when the function exits
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(cfg); err != nil {
		log.Error("can't encode Config to file", "err", err)
		return err
	}
	log.Info("Config generated, wait 30 seconds")

	time.Sleep(30 * time.Second)
	// deploy runners
	for i := 0; i < n; i++ {
		if vms[i] == nil {
			continue
		}
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
		if vms[i] == nil {
			continue
		}
		wg.Add(1)
		go func(id int) {
			if err := vms[id].startRunner(configFileName, networkMode, ""); err != nil {
				log.Crit("error starting runner", "id", id, "err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("Finished!")
	for i := range vms {
		if vms[i] == nil {
			continue
		}
		log.Info("Netdiag runner deployed", "id", i, "ip", vms[i].ip)
	}
	return nil
}

func cleanup(c *cli.Context) error {
	log.Info("cleaning up cluster")

	cfg := readConfigFile(c.String(configFlag.Name))
	vms := make([]*vm, len(cfg.Nodes))
	var wg sync.WaitGroup
	now := time.Now()
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		log.Error("NewInstancesRESTClient: %v", err)
		return err
	}
	for i, n := range cfg.Nodes {
		n := n
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vms[id] = newVM(id, n.Ip, n.InstanceName, n.Zone, cfg.GcpUsername)
			if c.Bool(logDownloadFlag.Name) {
				vms[id].downloadLogs()
			}
			vms[id].deleteRunner(ctx, instancesClient, cfg.GcpProjectId)
		}(i)
	}

	wg.Wait()
	err = os.Remove(c.String(configFlag.Name))
	if err != nil {
		log.Error("Config removal failed", "error", err)
	}
	log.Info("removed Config", "file", c.String(configFlag.Name))
	log.Info("Cluster cleaned up", "duration(s)", time.Now().Sub(now).Seconds())
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
				log.Crit("error Killing runner", "id", id, "err", err)
			}
			time.Sleep(time.Second)
			if !c.Bool(restartOnlyFlag.Name) {
				if err := vms[id].deployRunner(c.String(configFlag.Name), false, true); err != nil {
					log.Crit("error deploying runner", "id", id, "err", err)
				}
				log.Info("Runner binary deployed", "id", id)
			}
			optFlag := ""
			if c.Bool(pprofFlag.Name) {
				optFlag += "--pprof"
			}
			if err := vms[id].startRunner(c.String(configFlag.Name), c.String(networkModeFlag.Name), optFlag); err != nil {
				log.Crit("error starting runner", "id", id, "err", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Info("Cluster deployed and running!", "duration(s)", time.Now().Sub(now).Seconds())
	return nil
}
func waitSSH(instanceName string, maxRetries int, retryInterval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		err := checkSSH(instanceName)
		if err == nil {
			return nil
		}
		log.Info("SSH connection failed", "Retrying", i+1, "error", err)
		time.Sleep(retryInterval)
	}
	return fmt.Errorf("failed to connect to SSH server after %d attempts", maxRetries)
}

func addKeyToAgent() error {
	// Build the command
	cmd := exec.Command("ssh-add", "~/.ssh/id_rsa-corp")

	// Capture output and errors
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error adding key to agent: %w\n%s", err, errBuf.String())
	}

	fmt.Println(out.String()) // Optional: Print output from ssh-add
	return nil
}

func checkSSH(instanceName string) error {
	// Set a short timeout for connection attempt
	timeout := time.Second * 5
	key, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa-corp"))
	if err != nil {
		log.Error("Failed to read ssh file, err:", err)
		return err
	}
	log.Info("key", "k", key)
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Error("Failed to parse private key", "err:", err)
		return err
	}

	userName := os.Getenv("USER")
	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	// Connect to SSH server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", instanceName), config)
	if err != nil {
		return err
	}

	defer client.Close()

	return nil
}
