// Modifications Copyright 2019 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
//
// This file is derived from cmd/geth/main.go (2018/06/04).
// Modified and improved for the klaytn development.

package main

import (
	"fmt"
	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/cmd/utils/nodecmd"
	"github.com/klaytn/klaytn/console"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/metrics"
	"github.com/klaytn/klaytn/metrics/prometheus"
	"github.com/klaytn/klaytn/node"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/urfave/cli.v1"
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"
)

var (
	logger = log.NewModuleLogger(log.CMDKSPN)

	// The app that holds all commands and flags.
	app = utils.NewApp(nodecmd.GetGitCommit(), "The command line interface for Klaytn Proxy Node")

	// flags that configure the node
	nodeFlags = append(nodecmd.CommonNodeFlags, nodecmd.KSPNFlags...)

	rpcFlags = nodecmd.CommonRPCFlags
)

var spnHelpFlagGroups = []utils.FlagGroup{
	{
		Name: "KLAY",
		Flags: []cli.Flag{
			utils.DbTypeFlag,
			utils.DataDirFlag,
			utils.KeyStoreDirFlag,
			utils.IdentityFlag,
			utils.SyncModeFlag,
			utils.GCModeFlag,
			utils.LightKDFFlag,
			utils.SrvTypeFlag,
			utils.ExtraDataFlag,
			nodecmd.ConfigFileFlag,
		},
	},
	{
		Name: "SERVICECHAIN",
		Flags: []cli.Flag{
			utils.ChildChainIndexingFlag,
			utils.MainBridgeFlag,
			utils.MainBridgeListenPortFlag,
			utils.SubBridgeFlag,
			utils.SubBridgeListenPortFlag,
			utils.AnchoringPeriodFlag,
			utils.SentChainTxsLimit,
			utils.ParentChainIDFlag,
			utils.VTRecoveryFlag,
			utils.VTRecoveryIntervalFlag,
			utils.ServiceChainAnchoringFlag,
		},
	},
	{
		Name: "ACCOUNT",
		Flags: []cli.Flag{
			utils.UnlockedAccountFlag,
			utils.PasswordFileFlag,
		},
	},
	{
		Name: "TXPOOL",
		Flags: []cli.Flag{
			utils.TxPoolNoLocalsFlag,
			utils.TxPoolJournalFlag,
			utils.TxPoolJournalIntervalFlag,
			utils.TxPoolPriceLimitFlag,
			utils.TxPoolPriceBumpFlag,
			utils.TxPoolExecSlotsAccountFlag,
			utils.TxPoolExecSlotsAllFlag,
			utils.TxPoolNonExecSlotsAccountFlag,
			utils.TxPoolNonExecSlotsAllFlag,
			utils.TxPoolLifetimeFlag,
			utils.TxResendIntervalFlag,
			utils.TxResendCountFlag,
			utils.TxResendUseLegacyFlag,
		},
	},
	{
		Name: "DATABASE",
		Flags: []cli.Flag{
			utils.LevelDBCacheSizeFlag,
			utils.NoPartitionedDBFlag,
			utils.NumStateTriePartitionsFlag,
			utils.LevelDBCompressionTypeFlag,
			utils.LevelDBNoBufferPoolFlag,
			utils.NoParallelDBWriteFlag,
			utils.SenderTxHashIndexingFlag,
		},
	},
	{
		Name: "STATE",
		Flags: []cli.Flag{
			utils.StateDBCachingFlag,
			utils.TrieMemoryCacheSizeFlag,
			utils.TrieBlockIntervalFlag,
		},
	},
	{
		Name: "CACHE",
		Flags: []cli.Flag{
			utils.CacheTypeFlag,
			utils.CacheScaleFlag,
			utils.CacheUsageLevelFlag,
			utils.MemorySizeFlag,
			utils.CacheWriteThroughFlag,
			utils.TxPoolStateCacheFlag,
			utils.TrieCacheLimitFlag,
		},
	},
	{
		Name: "NETWORKING",
		Flags: []cli.Flag{
			utils.BootnodesFlag,
			utils.ListenPortFlag,
			utils.SubListenPortFlag,
			utils.MultiChannelUseFlag,
			utils.MaxConnectionsFlag,
			utils.MaxPendingPeersFlag,
			utils.TargetGasLimitFlag,
			utils.NATFlag,
			utils.NoDiscoverFlag,
			utils.RWTimerWaitTimeFlag,
			utils.RWTimerIntervalFlag,
			utils.NetrestrictFlag,
			utils.NodeKeyFileFlag,
			utils.NodeKeyHexFlag,
			utils.NetworkIdFlag,
		},
	},
	{
		Name: "METRICS",
		Flags: []cli.Flag{
			utils.MetricsEnabledFlag,
			utils.PrometheusExporterFlag,
			utils.PrometheusExporterPortFlag,
		},
	},
	{
		Name: "VIRTUAL MACHINE",
		Flags: []cli.Flag{
			utils.VMEnableDebugFlag,
			utils.VMLogTargetFlag,
		},
	},
	{
		Name: "API AND CONSOLE",
		Flags: []cli.Flag{
			utils.RPCEnabledFlag,
			utils.RPCListenAddrFlag,
			utils.RPCPortFlag,
			utils.RPCCORSDomainFlag,
			utils.RPCVirtualHostsFlag,
			utils.RPCApiFlag,
			utils.IPCDisabledFlag,
			utils.IPCPathFlag,
			utils.WSEnabledFlag,
			utils.WSListenAddrFlag,
			utils.WSPortFlag,
			utils.WSApiFlag,
			utils.WSAllowedOriginsFlag,
			utils.GRPCEnabledFlag,
			utils.GRPCListenAddrFlag,
			utils.GRPCPortFlag,
			utils.JSpathFlag,
			utils.ExecFlag,
			utils.PreloadJSFlag,
		},
	},
	{
		Name:  "LOGGING AND DEBUGGING",
		Flags: debug.Flags,
	},
	{
		Name: "MISC",
		Flags: []cli.Flag{
			utils.GenKeyFlag,
			utils.WriteAddressFlag,
		},
	},
}

func init() {
	utils.InitHelper()
	// Initialize the CLI app and start kspn
	app.Action = nodecmd.RunKlaytnNode
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018-2019 The klaytn Authors"
	app.Commands = []cli.Command{
		// See utils/nodecmd/chaincmd.go:
		nodecmd.InitCommand,

		// See utils/nodecmd/accountcmd.go
		nodecmd.AccountCommand,

		// See utils/nodecmd/consolecmd.go:
		nodecmd.GetConsoleCommand(nodeFlags, rpcFlags),
		nodecmd.AttachCommand,

		// See utils/nodecmd/versioncmd.go:
		nodecmd.VersionCommand,

		// See utils/nodecmd/dumpconfigcmd.go:
		nodecmd.GetDumpConfigCommand(nodeFlags, rpcFlags),
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, nodecmd.ConsoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)

	cli.AppHelpTemplate = utils.GlobalAppHelpTemplate
	cli.HelpPrinter = utils.NewHelpPrinter(spnHelpFlagGroups)

	app.CommandNotFound = nodecmd.CommandNotExist

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		logDir := (&node.Config{DataDir: utils.MakeDataDir(ctx)}).ResolvePath("logs")
		debug.CreateLogDir(logDir)
		if err := debug.Setup(ctx); err != nil {
			return err
		}

		// Start prometheus exporter
		if metrics.Enabled {
			logger.Info("Enabling metrics collection")
			if metrics.EnabledPrometheusExport {
				logger.Info("Enabling Prometheus Exporter")
				pClient := prometheusmetrics.NewPrometheusProvider(metrics.DefaultRegistry, "klaytn",
					"", prometheus.DefaultRegisterer, 3*time.Second)
				go pClient.UpdatePrometheusMetrics()
				http.Handle("/metrics", promhttp.Handler())
				port := ctx.GlobalInt(metrics.PrometheusExporterPortFlag)

				go func() {
					err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
					if err != nil {
						logger.Error("PrometheusExporter starting failed:", "port", port, "err", err)
					}
				}()
			}
		}

		// Start system runtime metrics collection
		go metrics.CollectProcessMetrics(3 * time.Second)

		utils.SetupNetwork(ctx)
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	// Set NodeTypeFlag to spn
	utils.NodeTypeFlag.Value = "spn"

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
