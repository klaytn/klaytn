// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from cmd/bootnode/main.go (2018/06/04).
// Modified and improved for the klaytn development.

package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/cmd/utils/nodecmd"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/p2p/nat"
	"github.com/klaytn/klaytn/networks/rpc"
	"gopkg.in/urfave/cli.v1"
)

var (
	logger = log.NewModuleLogger(log.CMDKBN)
)

const (
	generateNodeKeySpecified = iota
	noPrivateKeyPathSpecified
	nodeKeyDuplicated
	writeOutAddress
	goodToGo
)

func bootnode(ctx *cli.Context) error {
	var (
		// Local variables
		err  error
		bcfg = bootnodeConfig{
			// Config variables
			networkID:    ctx.GlobalUint64(utils.NetworkIdFlag.Name),
			addr:         ctx.GlobalString(utils.BNAddrFlag.Name),
			genKeyPath:   ctx.GlobalString(utils.GenKeyFlag.Name),
			nodeKeyFile:  ctx.GlobalString(utils.NodeKeyFileFlag.Name),
			nodeKeyHex:   ctx.GlobalString(utils.NodeKeyHexFlag.Name),
			natFlag:      ctx.GlobalString(utils.NATFlag.Name),
			netrestrict:  ctx.GlobalString(utils.NetrestrictFlag.Name),
			writeAddress: ctx.GlobalBool(utils.WriteAddressFlag.Name),

			IPCPath:          "klay.ipc",
			DataDir:          ctx.GlobalString(utils.DataDirFlag.Name),
			HTTPPort:         DefaultHTTPPort,
			HTTPModules:      []string{"net"},
			HTTPVirtualHosts: []string{"localhost"},
			HTTPTimeouts:     rpc.DefaultHTTPTimeouts,
			WSPort:           DefaultWSPort,
			WSModules:        []string{"net"},
			GRPCPort:         DefaultGRPCPort,

			Logger: log.NewModuleLogger(log.CMDKBN),
		}
	)

	if err = nodecmd.CheckCommands(ctx); err != nil {
		return err
	}

	setIPC(ctx, &bcfg)
	// httptype is http or fasthttp
	if ctx.GlobalIsSet(utils.SrvTypeFlag.Name) {
		bcfg.HTTPServerType = ctx.GlobalString(utils.SrvTypeFlag.Name)
	}
	setHTTP(ctx, &bcfg)
	setWS(ctx, &bcfg)
	setgRPC(ctx, &bcfg)
	setAuthorizedNodes(ctx, &bcfg)

	// Check exit condition
	switch bcfg.checkCMDState() {
	case generateNodeKeySpecified:
		bcfg.generateNodeKey()
	case noPrivateKeyPathSpecified:
		return errors.New("Use --nodekey or --nodekeyhex to specify a private key")
	case nodeKeyDuplicated:
		return errors.New("Options --nodekey and --nodekeyhex are mutually exclusive")
	case writeOutAddress:
		bcfg.doWriteOutAddress()
	default:
		err = bcfg.readNodeKey()
		if err != nil {
			return err
		}
	}

	err = bcfg.validateNetworkParameter()
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", bcfg.listenAddr)
	if err != nil {
		log.Fatalf("Failed to ResolveUDPAddr: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to ListenUDP: %v", err)
	}

	realaddr := conn.LocalAddr().(*net.UDPAddr)
	if bcfg.natm != nil {
		if !realaddr.IP.IsLoopback() {
			go nat.Map(bcfg.natm, nil, "udp", realaddr.Port, realaddr.Port, "Klaytn node discovery")
		}
		// TODO: react to external IP changes over time.
		if ext, err := bcfg.natm.ExternalIP(); err == nil {
			realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
		}
	}

	cfg := discover.Config{
		NetworkID:       bcfg.networkID,
		PrivateKey:      bcfg.nodeKey,
		AnnounceAddr:    realaddr,
		NetRestrict:     bcfg.restrictList,
		Conn:            conn,
		Addr:            realaddr,
		Id:              discover.PubkeyID(&bcfg.nodeKey.PublicKey),
		NodeType:        p2p.ConvertNodeType(common.BOOTNODE),
		AuthorizedNodes: bcfg.AuthorizedNodes,
	}

	tab, err := discover.ListenUDP(&cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	node, err := New(&bcfg)
	if err != nil {
		return err
	}
	node.appendAPIs(NewBN(tab).APIs())
	if err := startNode(node); err != nil {
		return err
	}
	node.Wait()
	return nil
}

func startNode(node *Node) error {
	if err := node.Start(); err != nil {
		return err
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		logger.Info("Got interrupt, shutting down...")
		go node.Stop()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				logger.Info("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
	}()
	return nil
}

func main() {
	var (
		cliFlags = []cli.Flag{
			utils.SrvTypeFlag,
			utils.DataDirFlag,
			utils.GenKeyFlag,
			utils.NodeKeyFileFlag,
			utils.NodeKeyHexFlag,
			utils.WriteAddressFlag,
			utils.BNAddrFlag,
			utils.NATFlag,
			utils.NetrestrictFlag,
			utils.MetricsEnabledFlag,
			utils.PrometheusExporterFlag,
			utils.PrometheusExporterPortFlag,
			utils.AuthorizedNodesFlag,
			utils.NetworkIdFlag,
		}
	)
	// TODO-Klaytn: remove `help` command
	app := utils.NewApp("", "the Klaytn's bootnode command line interface")
	app.Name = "kbn"
	app.Copyright = "Copyright 2018 The klaytn Authors"
	app.UsageText = app.Name + " [global options] [commands]"
	app.Flags = append(app.Flags, cliFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	app.Flags = append(app.Flags, nodecmd.CommonRPCFlags...)
	app.Commands = []cli.Command{
		nodecmd.VersionCommand,
		nodecmd.AttachCommand,
	}

	app.Action = bootnode

	app.CommandNotFound = nodecmd.CommandNotExist
	app.OnUsageError = nodecmd.OnUsageError

	app.Before = nodecmd.BeforeRunBootnode

	app.After = func(c *cli.Context) error {
		debug.Exit()
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
