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
	"io"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fjl/memsize/memsizeui"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/log/term"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Memsize memsizeui.Handler

var (
	verbosityFlag = cli.IntFlag{
		Name:    "verbosity",
		Usage:   "Logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail",
		Value:   3,
		Aliases: []string{"debug-profile.verbosity"},
		EnvVars: []string{"KLAYTN_VERBOSITY"},
	}
	vmoduleFlag = cli.StringFlag{
		Name:    "vmodule",
		Usage:   "Per-module verbosity: comma-separated list of <pattern>=<level> (e.g. klay/*=5,p2p=4)",
		Value:   "",
		Aliases: []string{"debug-profile.vmodule"},
		EnvVars: []string{"KLAYTN_VMODULE"},
	}
	backtraceAtFlag = cli.StringFlag{
		Name:    "backtrace",
		Usage:   "Request a stack trace at a specific logging statement (e.g. \"block.go:271\")",
		Value:   "",
		Aliases: []string{"debug-profile.backtrace"},
		EnvVars: []string{"KLAYTN_BACKTRACE"},
	}
	debugFlag = cli.BoolFlag{
		Name:    "debug",
		Usage:   "Prepends log messages with call-site location (file and line number)",
		Aliases: []string{"debug-profile.print-site"},
		EnvVars: []string{"KLAYTN_DEBUG"},
	}
	logFormatFlag = cli.StringFlag{
		Name:   "log.format",
		Usage:  "Log format to use (json|logfmt|terminal)",
		Value:  "terminal",
		EnvVar: "KLAYTN_LOGFORMAT",
	}
	logFileFlag = cli.StringFlag{
		Name:   "log.file",
		Usage:  "Write logs to a file",
		EnvVar: "KLAYTN_LOGFILE",
	}
	logRotateFlag = cli.BoolFlag{
		Name:   "log.rotate",
		Usage:  "Enables log file rotation",
		EnvVar: "KLAYTN_LOGROTATE",
	}
	logMaxSizeMBsFlag = cli.IntFlag{
		Name:   "log.maxsize",
		Usage:  "Maximum size in MBs of a single log file (use with --log.rotate flag)",
		Value:  100,
		EnvVar: "KLAYTN_LOGMAXSIZE",
	}
	logMaxBackupsFlag = cli.IntFlag{
		Name:   "log.maxbackups",
		Usage:  "Maximum number of log files to retain (use with --log.rotate flag)",
		Value:  10,
		EnvVar: "KLAYTN_LOGMAXBACKUPS",
	}
	logMaxAgeFlag = cli.IntFlag{
		Name:   "log.maxage",
		Usage:  "Maximum number of days to retain a log file (use with --log.rotate flag)",
		Value:  30,
		EnvVar: "KLAYTN_LOGMAXAGE",
	}
	logCompressFlag = cli.BoolFlag{
		Name:   "log.compress",
		Usage:  "Compress the log files (use with --log.rotate flag)",
		EnvVar: "KLAYTN_LOGCOMPRESS",
	}
	pprofFlag = cli.BoolFlag{
		Name:    "pprof",
		Usage:   "Enable the pprof HTTP server",
		Aliases: []string{"debug-profile.pprof.enable"},
		EnvVars: []string{"KLAYTN_PPROF"},
	}
	pprofPortFlag = cli.IntFlag{
		Name:    "pprofport",
		Usage:   "pprof HTTP server listening port",
		Value:   6060,
		Aliases: []string{"debug-profile.pprof.port"},
		EnvVars: []string{"KLAYTN_PPROFPORT"},
	}
	pprofAddrFlag = cli.StringFlag{
		Name:    "pprofaddr",
		Usage:   "pprof HTTP server listening interface",
		Value:   "127.0.0.1",
		Aliases: []string{"debug-profile.pprof.addr"},
		EnvVars: []string{"KLAYTN_PPROFADDR"},
	}
	memprofileFlag = cli.StringFlag{
		Name:    "memprofile",
		Usage:   "Write memory profile to the given file",
		Aliases: []string{"debug-profile.mem-profile.file-name"},
		EnvVars: []string{"KLAYTN_MEMPROFILE"},
	}
	memprofilerateFlag = cli.IntFlag{
		Name:    "memprofilerate",
		Usage:   "Turn on memory profiling with the given rate",
		Value:   runtime.MemProfileRate,
		Aliases: []string{"debug-profile.mem-profile.rate"},
		EnvVars: []string{"KLAYTN_MEMPROFILERATE"},
	}
	blockprofilerateFlag = cli.IntFlag{
		Name:    "blockprofilerate",
		Usage:   "Turn on block profiling with the given rate",
		Aliases: []string{"debug-profile.block-profile.rate"},
		EnvVars: []string{"KLAYTN_BLOCKPROFILERATE"},
	}
	cpuprofileFlag = cli.StringFlag{
		Name:    "cpuprofile",
		Usage:   "Write CPU profile to the given file",
		Aliases: []string{"debug-profile.cpu-profile.file-name"},
		EnvVars: []string{"KLAYTN_CPUPROFILE"},
	}
	traceFlag = cli.StringFlag{
		Name:    "trace",
		Usage:   "Write execution trace to the given file",
		Aliases: []string{"debug-profile.trace.file-name"},
		EnvVars: []string{"KLAYTN_TRACE"},
	}
)

// Flags holds all command-line flags required for debugging.
var Flags = []cli.Flag{
	altsrc.NewIntFlag(&verbosityFlag),
	altsrc.NewStringFlag(&vmoduleFlag),
	altsrc.NewStringFlag(&backtraceAtFlag),
	altsrc.NewBoolFlag(&debugFlag),
	altsrc.NewStringFlag(&logFormatFlag),
	altsrc.NewStringFlag(&logFileFlag),
	altsrc.NewBoolFlag(&logRotateFlag),
	altsrc.NewIntFlag(&logMaxSizeMBsFlag),
	altsrc.NewIntFlag(&logMaxBackupsFlag),
	altsrc.NewIntFlag(&logMaxAgeFlag),
	altsrc.NewBoolFlag(&logCompressFlag),
	altsrc.NewBoolFlag(&pprofFlag),
	altsrc.NewStringFlag(&pprofAddrFlag),
	altsrc.NewIntFlag(&pprofPortFlag),
	altsrc.NewStringFlag(&memprofileFlag),
	altsrc.NewIntFlag(&memprofilerateFlag),
	altsrc.NewIntFlag(&blockprofilerateFlag),
	altsrc.NewStringFlag(&cpuprofileFlag),
	altsrc.NewStringFlag(&traceFlag),
}

var glogger *log.GlogHandler

func init() {
	glogger = log.NewGlogHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(false)))
	glogger.Verbosity(log.LvlInfo)
	log.Root().SetHandler(glogger)
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
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		logger.Warn("Failed to create a directory", "logDir", logDir, "err", err)
	}
}

// Setup initializes profiling and logging based on the CLI flags.
// It should be called as early as possible in the program.
func Setup(ctx *cli.Context) error {
	var (
		output     = io.Writer(os.Stderr)
		logFmtFlag = ctx.String(logFormatFlag.Name)
		logfmt     log.Format
	)

	switch logFmtFlag {
	case "json":
		logfmt = log.JsonFormat()
	case "logfmt":
		logfmt = log.LogfmtFormat()
	case "terminal":
		useColor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
		if useColor {
			output = colorable.NewColorableStdout()
		}
		logfmt = log.TerminalFormat(useColor)
	default:
		// Unknown log format specified
		return fmt.Errorf("unknown log format: %v", ctx.String(logFormatFlag.Name))
	}

	var (
		ostream  = log.StreamHandler(output, logfmt)
		logFile  = ctx.String(logFileFlag.Name)
		rotation = ctx.Bool(logRotateFlag.Name)
		context  = []interface{}{"rotate", rotation}
	)

	// if logFile is not set when rotation, set default log name
	if rotation && len(logFile) == 0 {
		logFile = filepath.Join(os.TempDir(), "klaytn-lumberjack.log")
	}
	if len(logFile) > 0 {
		context = append(context, "format", logFmtFlag, "location", logFile)
		if err := validateLogLocation(logFile); err != nil {
			return fmt.Errorf("tried to create a temporary file to verify that the log path is writable, but it failed: %v", err)
		}
		if rotation {
			ostream = log.StreamHandler(&lumberjack.Logger{
				Filename:   logFile,
				MaxSize:    ctx.Int(logMaxSizeMBsFlag.Name),
				MaxBackups: ctx.Int(logMaxBackupsFlag.Name),
				MaxAge:     ctx.Int(logMaxAgeFlag.Name),
				Compress:   ctx.Bool(logCompressFlag.Name),
			}, logfmt)
		} else {
			logOutputStream, err := log.FileHandler(logFile, logfmt)
			if err != nil {
				return err
			}
			ostream = logOutputStream
		}
	}
	glogger = log.NewGlogHandler(ostream)

	// logging
	log.PrintOrigins(ctx.Bool(debugFlag.Name))
	if err := log.ChangeGlobalLogLevel(glogger, log.Lvl(ctx.Int(verbosityFlag.Name))); err != nil {
		return err
	}
	if err := glogger.Vmodule(ctx.String(vmoduleFlag.Name)); err != nil {
		return err
	}
	if len(ctx.String(backtraceAtFlag.Name)) != 0 {
		if err := glogger.BacktraceAt(ctx.String(backtraceAtFlag.Name)); err != nil {
			return err
		}
	}
	log.Root().SetHandler(glogger)

	// profiling, tracing
	runtime.MemProfileRate = ctx.Int(memprofilerateFlag.Name)
	Handler.SetBlockProfileRate(ctx.Int(blockprofilerateFlag.Name))
	if traceFile := ctx.String(traceFlag.Name); traceFile != "" {
		if err := Handler.StartGoTrace(traceFile); err != nil {
			return err
		}
	}
	if cpuFile := ctx.String(cpuprofileFlag.Name); cpuFile != "" {
		if err := Handler.StartCPUProfile(cpuFile); err != nil {
			return err
		}
	}
	Handler.memFile = ctx.String(memprofileFlag.Name)

	// pprof server
	if ctx.Bool(pprofFlag.Name) {
		addr := ctx.String(pprofAddrFlag.Name)
		port := ctx.Int(pprofPortFlag.Name)
		Handler.StartPProf(&addr, &port)
	}
	if len(logFile) > 0 {
		logger.Info("Logging configured", context...)
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

func validateLogLocation(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("error creating the directory: %w", err)
	}
	// Check if the path is writable by trying to create a temporary file
	tmp := path + ".temp"
	if f, err := os.Create(tmp); err != nil {
		return err
	} else {
		f.Close()
	}
	return os.Remove(tmp)
}
