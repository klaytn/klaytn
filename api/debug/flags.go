// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from internal/debug/flags.go (2018/06/04).
// Modified and improved for the klaytn development.

package debug

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime"

	"io"

	"github.com/fjl/memsize/memsizeui"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/log/term"
	colorable "github.com/mattn/go-colorable"
	"gopkg.in/urfave/cli.v1"
)

var Memsize memsizeui.Handler

var (
	verbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail",
		Value: 3,
	}
	vmoduleFlag = cli.StringFlag{
		Name:  "vmodule",
		Usage: "Per-module verbosity: comma-separated list of <pattern>=<level> (e.g. klay/*=5,p2p=4)",
		Value: "",
	}
	backtraceAtFlag = cli.StringFlag{
		Name:  "backtrace",
		Usage: "Request a stack trace at a specific logging statement (e.g. \"block.go:271\")",
		Value: "",
	}
	debugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "Prepends log messages with call-site location (file and line number)",
	}
	pprofFlag = cli.BoolFlag{
		Name:  "pprof",
		Usage: "Enable the pprof HTTP server",
	}
	pprofPortFlag = cli.IntFlag{
		Name:  "pprofport",
		Usage: "pprof HTTP server listening port",
		Value: 6060,
	}
	pprofAddrFlag = cli.StringFlag{
		Name:  "pprofaddr",
		Usage: "pprof HTTP server listening interface",
		Value: "127.0.0.1",
	}
	memprofileFlag = cli.StringFlag{
		Name:  "memprofile",
		Usage: "Write memory profile to the given file",
	}
	memprofilerateFlag = cli.IntFlag{
		Name:  "memprofilerate",
		Usage: "Turn on memory profiling with the given rate",
		Value: runtime.MemProfileRate,
	}
	blockprofilerateFlag = cli.IntFlag{
		Name:  "blockprofilerate",
		Usage: "Turn on block profiling with the given rate",
	}
	cpuprofileFlag = cli.StringFlag{
		Name:  "cpuprofile",
		Usage: "Write CPU profile to the given file",
	}
	traceFlag = cli.StringFlag{
		Name:  "trace",
		Usage: "Write execution trace to the given file",
	}
)

// Flags holds all command-line flags required for debugging.
var Flags = []cli.Flag{
	verbosityFlag, vmoduleFlag, backtraceAtFlag, debugFlag,
	pprofFlag, pprofAddrFlag, pprofPortFlag,
	memprofileFlag, memprofilerateFlag,
	blockprofilerateFlag, cpuprofileFlag, traceFlag,
}

var glogger *log.GlogHandler

func init() {
	usecolor := term.IsTty(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	glogger = log.NewGlogHandler(log.StreamHandler(output, log.TerminalFormat(usecolor)))
}

func GetGlogger() (*log.GlogHandler, error) {
	if glogger != nil {
		return glogger, nil
	}
	return nil, fmt.Errorf("glogger is nil")
}

// CreateLogDir creates a directory whose path is logdir as well as empty log files.
func CreateLogDir(logDir string) {
	if logDir == "" {
		return
	}
	Handler.logDir = logDir

	// Currently failures on directory or file creation is treated as a warning.
	if err := os.MkdirAll(logDir, 0700); err != nil {
		logger.Warn("Failed to create a directory", "logDir", logDir, "err", err)
	}
}

// Setup initializes profiling and logging based on the CLI flags.
// It should be called as early as possible in the program.
func Setup(ctx *cli.Context) error {
	// logging
	log.PrintOrigins(ctx.GlobalBool(debugFlag.Name))
	if err := log.ChangeGlobalLogLevel(glogger, log.Lvl(ctx.GlobalInt(verbosityFlag.Name))); err != nil {
		return err
	}
	if err := glogger.Vmodule(ctx.GlobalString(vmoduleFlag.Name)); err != nil {
		return err
	}
	if len(ctx.GlobalString(backtraceAtFlag.Name)) != 0 {
		if err := glogger.BacktraceAt(ctx.GlobalString(backtraceAtFlag.Name)); err != nil {
			return err
		}
	}
	log.Root().SetHandler(glogger)

	// profiling, tracing
	runtime.MemProfileRate = ctx.GlobalInt(memprofilerateFlag.Name)
	Handler.SetBlockProfileRate(ctx.GlobalInt(blockprofilerateFlag.Name))
	if traceFile := ctx.GlobalString(traceFlag.Name); traceFile != "" {
		if err := Handler.StartGoTrace(traceFile); err != nil {
			return err
		}
	}
	if cpuFile := ctx.GlobalString(cpuprofileFlag.Name); cpuFile != "" {
		if err := Handler.StartCPUProfile(cpuFile); err != nil {
			return err
		}
	}
	Handler.memFile = ctx.GlobalString(memprofileFlag.Name)

	// pprof server
	if ctx.GlobalBool(pprofFlag.Name) {
		addr := ctx.GlobalString(pprofAddrFlag.Name)
		port := ctx.GlobalInt(pprofPortFlag.Name)
		Handler.StartPProf(&addr, &port)
	}
	return nil
}

// Exit stops all running profiles, flushing their output to the
// respective file.
func Exit() {
	if Handler.vmLogFile != nil {
		Handler.vmLogFile.Close()
		Handler.vmLogFile = nil
	}
	if Handler.memFile != "" {
		Handler.WriteMemProfile(Handler.memFile)
	}
	Handler.StopCPUProfile()
	Handler.StopGoTrace()
	Handler.StopPProf()
}
