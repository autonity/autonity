// Copyright 2014 The go-ethereum Authors
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

// geth is the official command-line client for Ethereum.
package main

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/autonity/autonity/accounts"
	"github.com/autonity/autonity/cmd/utils"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/console/prompt"
	"github.com/autonity/autonity/eth/downloader"
	"github.com/autonity/autonity/ethclient"
	"github.com/autonity/autonity/internal/debug"
	"github.com/autonity/autonity/internal/flags"
	"github.com/autonity/autonity/internal/version"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/node"

	// Force-load the tracer engines to trigger registration
	_ "github.com/autonity/autonity/eth/tracers/js"
	_ "github.com/autonity/autonity/eth/tracers/native"
)

const (
	clientIdentifier = "autonity" // Client identifier to advertise over the network
	syncModeLight    = "light"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	gitDate   = ""
	// The app that holds all commands and flags.
	app = flags.NewApp("the autonity command line interface")
	// flags that configure the node
	nodeFlags = slices.Concat([]cli.Flag{
		utils.IdentityFlag,
		utils.PasswordFileFlag,
		utils.BootnodesFlag,
		utils.InitGenesisFlag,
		utils.MinFreeDiskSpaceFlag,
		utils.KeyStoreDirFlag,
		utils.ExternalSignerFlag,
		utils.USBFlag,
		utils.SmartCardDaemonPathFlag,
		utils.TxPoolLocalsFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceLimitFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolLifetimeFlag,
		utils.SyncModeFlag,
		utils.ExitWhenSyncedFlag,
		utils.GCModeFlag,
		utils.SnapshotFlag,
		utils.EthRequiredBlocksFlag,
		utils.BloomFilterSizeFlag,
		utils.CacheFlag,
		utils.CacheDatabaseFlag,
		utils.CacheTrieFlag,
		utils.CacheGCFlag,
		utils.CacheSnapshotFlag,
		utils.CacheNoPrefetchFlag,
		utils.CachePreimagesFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.MiningEnabledFlag,
		utils.MinerNotifyFlag,
		utils.LightKDFFlag,
		utils.MinerGasPriceFlag,
		utils.MinerExtraDataFlag,
		utils.MinerRecommitIntervalFlag,
		utils.MinerNoVerifyFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.AutonityKeysFileFlag,
		utils.AutonityKeysHexFlag,
		utils.OracleKeyFileFlag,
		utils.OracleKeyHexFlag,
		utils.WriteAddrFlag,
		utils.DNSDiscoveryFlag,
		utils.DeveloperFlag,
		utils.DeveloperGasLimitFlag,
		utils.DeveloperEtherbaseFlag,
		utils.VMEnableDebugFlag,
		utils.NetworkIdFlag,
		utils.EthStatsURLFlag,
		utils.NoCompactionFlag,
		utils.GpoBlocksFlag,
		utils.GpoPercentileFlag,
		utils.GpoMaxGasPriceFlag,
		utils.GpoIgnoreGasPriceFlag,
		utils.MinerNotifyFullFlag,
		utils.ConsensusListenPortFlag,
		utils.ConsensusNATFlag,
		configFileFlag,
	}, utils.NetworkFlags, utils.DatabaseFlags)

	rpcFlags = []cli.Flag{
		utils.HTTPEnabledFlag,
		utils.HTTPListenAddrFlag,
		utils.HTTPPortFlag,
		utils.HTTPCORSDomainFlag,
		utils.HTTPVirtualHostsFlag,
		utils.GraphQLEnabledFlag,
		utils.GraphQLCORSDomainFlag,
		utils.GraphQLVirtualHostsFlag,
		utils.HTTPApiFlag,
		utils.HTTPPathPrefixFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.WSPathPrefixFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
		utils.RPCGlobalGasCapFlag,
		utils.RPCGlobalEVMTimeoutFlag,
		utils.RPCGlobalTxFeeCapFlag,
		utils.AllowUnprotectedTxs,
	}

	metricsFlags = []cli.Flag{
		utils.MetricsEnabledFlag,
		utils.MetricsEnabledExpensiveFlag,
		utils.MetricsHTTPFlag,
		utils.MetricsPortFlag,
		utils.MetricsEnableInfluxDBFlag,
		utils.MetricsInfluxDBEndpointFlag,
		utils.MetricsInfluxDBDatabaseFlag,
		utils.MetricsInfluxDBUsernameFlag,
		utils.MetricsInfluxDBPasswordFlag,
		utils.MetricsInfluxDBTagsFlag,
		utils.MetricsEnableInfluxDBV2Flag,
		utils.MetricsInfluxDBTokenFlag,
		utils.MetricsInfluxDBBucketFlag,
		utils.MetricsInfluxDBOrganizationFlag,
	}
)

func init() {
	// Initialize the CLI app and start Autonity
	app.Action = autonity
	app.HideVersion = true // we have a command to print the version
	app.Copyright = ""
	app.Commands = []*cli.Command{
		// See chaincmd.go:
		importCommand,
		exportCommand,
		importPreimagesCommand,
		removedbCommand,
		dumpCommand,
		// See accountcmd.go:
		accountCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		// See misccmd.go:
		versionCommand,
		licenseCommand,
		ownershipProofCommand,
		genAutonityKeysCommand,
		// See config.go
		dumpConfigCommand,
		// see dbcmd.go
		dbCommand,
		// See snapshot.go
		snapshotCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = slices.Concat(
		nodeFlags,
		rpcFlags,
		consoleFlags,
		debug.Flags,
		metricsFlags,
	)
	flags.AutoEnvVars(app.Flags, "AUT")

	app.Before = func(ctx *cli.Context) error {
		maxprocs.Set() // Automatically set GOMAXPROCS to match Linux container CPU quota.
		flags.MigrateGlobalFlags(ctx)
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		flags.CheckEnvVars(ctx, app.Flags, "AUT")
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		prompt.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// prepare manipulates memory cache allowance and setups metric system.
// This function should be called before launching devp2p stack.
func prepare(ctx *cli.Context) {
	// If we're running a known preset, log it for convenience.
	switch {
	case ctx.IsSet(utils.PiccadillyFlag.Name):
		log.Info(`Starting Autonity on Piccadilly Testnet

ααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααα		
ααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααα

                        888                     d8b 888             
                        888                     Y8P 888             
                        888                         888             
       8888b.  888  888 888888 .d88b.  88888b.  888 888888 888  888 
          "88b 888  888 888   d88""88b 888 "88b 888 888    888  888 
      .d888888 888  888 888   888  888 888  888 888 888    888  888 
      888  888 Y88b 888 Y88b. Y88..88P 888  888 888 Y88b.  Y88b 888 
      "Y888888  "Y88888  "Y888 "Y88P"  888  888 888  "Y888  "Y88888 
                                                                888 
                         autonity.org                      Y8b d88P 
                                                            "Y88P"  

ααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααα
ααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααα

Take part in the Piccadilly Tiber Challenges! More infos @ https://autonity.org
Please remain tuned to our social channels for announcements
Discord:  https://discord.com/invite/autonity 
Telegram: https://t.me/autonity
X:        https://twitter.com/autonity_


`)
	case ctx.IsSet(utils.BakerlooFlag.Name):
		log.Info("Starting Autonity on Bakerloo testnet")
	case ctx.IsSet(utils.DeveloperFlag.Name):
		log.Info("Starting Autonity in ephemeral dev mode")
		log.Warn(`You are running autonity in --dev mode. Please note the following:

  1. This mode is only intended for fast, iterative development without assumptions on
     security or persistence.
  2. The database is created in memory. Therefore, shutting down your computer or losing 
	 power will wipe your entire block data and chain state for your dev environment.
  3. A random(unless specified explicitly with --dev.etherbase flag), pre-allocated developer 
     account will be available and unlocked as eth.coinbase, which can be used for testing. 
     The random dev account is temporary, stored on a ramdisk, and will be lost if your machine 
     is restarted.
  4. Networking is disabled; there is no listen-address, the maximum number of peers is set
     to 0, and discovery is disabled.
  5. Optionally you can use an external account, which would be pre-allocated and unlocked as
     eth.coinbase, set following flags for external account:
	--dev.etherbase <public address of account in hex format>
	--password <password file path>
	--keystore <account's keystore directory path>
`)
	case !ctx.IsSet(utils.NetworkIdFlag.Name):
		version, _ := version.Info()
		log.Info(`Starting Autonity on Mainnet`, "version", version)
		log.Info(banner)
	default:
		log.Info("Starting the Autonity node client", "version", version.Semantic, "networkid", ctx.Int(utils.NetworkIdFlag.Name))
	}
	// If we're a full node on mainnet without --cache specified, bump default cache allowance
	if ctx.String(utils.SyncModeFlag.Name) != syncModeLight && !ctx.IsSet(utils.CacheFlag.Name) &&
		!ctx.IsSet(utils.NetworkIdFlag.Name) && !ctx.IsSet(utils.DeveloperFlag.Name) {
		// Make sure we're not on any supported preconfigured testnet either
		ctx.Set(utils.CacheFlag.Name, strconv.Itoa(4096))
	}

	// Start system runtime metrics collection
	go metrics.CollectProcessMetrics(3 * time.Second)
}

// autonity is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func autonity(ctx *cli.Context) error {
	if args := ctx.Args().Slice(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}

	prepare(ctx)

	stack := makeFullNode(ctx)
	defer stack.Close()

	startNode(ctx, stack, false)
	stack.Wait()

	return nil
}

// startNode boots up the system node and all registered protocols, after which
// it starts the RPC/IPC interfaces and the miner.
func startNode(ctx *cli.Context, stack *node.Node, isConsole bool) {
	// Start up the node itself
	utils.StartNode(ctx, stack, isConsole)

	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	// Create a client to interact with local geth node.
	rpcClient := stack.Attach()
	ethClient := ethclient.NewClient(rpcClient)

	go func() {
		// Open any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				var derivationPaths []accounts.DerivationPath
				if event.Wallet.URL().Scheme == "ledger" {
					derivationPaths = append(derivationPaths, accounts.LegacyLedgerBaseDerivationPath)
				}
				derivationPaths = append(derivationPaths, accounts.DefaultBaseDerivationPath)

				event.Wallet.SelfDerive(derivationPaths, ethClient)

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

	// Spawn a standalone goroutine for status synchronization monitoring,
	// close the node when synchronization is complete if user required.
	if ctx.Bool(utils.ExitWhenSyncedFlag.Name) {
		go func() {
			sub := stack.EventMux().Subscribe(downloader.DoneEvent{})
			defer sub.Unsubscribe()
			for {
				event := <-sub.Chan()
				if event == nil {
					continue
				}
				done, ok := event.Data.(downloader.DoneEvent)
				if !ok {
					continue
				}
				if timestamp := time.Unix(int64(done.Latest.Time), 0); time.Since(timestamp) < 10*time.Minute {
					log.Info("Synchronisation completed", "latestnum", done.Latest.Number, "latesthash", done.Latest.Hash(),
						"age", common.PrettyAge(timestamp))
					stack.Close()
				}
			}
		}()
	}
}

const banner = `

ααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααα		
ααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααα

                        888                     d8b 888             
                        888                     Y8P 888             
                        888                         888             
       8888b.  888  888 888888 .d88b.  88888b.  888 888888 888  888 
          "88b 888  888 888   d88""88b 888 "88b 888 888    888  888 
      .d888888 888  888 888   888  888 888  888 888 888    888  888 
      888  888 Y88b 888 Y88b. Y88..88P 888  888 888 Y88b.  Y88b 888 
      "Y888888  "Y88888  "Y888 "Y88P"  888  888 888  "Y888  "Y88888 
                                                                888 
                         autonity.org                      Y8b d88P 
                                                            "Y88P"  

ααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααα
ααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααααα

Discord:  https://discord.com/invite/autonity 
Telegram: https://t.me/autonity
X:        https://twitter.com/autonity_

`
