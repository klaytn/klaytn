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
// This file is derived from internal/debug/api.go (2018/06/04).
// Modified and improved for the klaytn development.

package debug

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/metrics/exp"
	"github.com/klaytn/klaytn/params"
	"github.com/rcrowley/go-metrics"
)

// Handler is the global debugging handler.
var Handler = new(HandlerT)
var logger = log.NewModuleLogger(log.APIDebug)

// HandlerT implements the debugging API.
// Do not create values of this type, use the one
// in the Handler variable instead.
type HandlerT struct {
	mu        sync.Mutex
	cpuW      io.WriteCloser
	cpuFile   string
	memFile   string
	traceW    io.WriteCloser
	traceFile string

	// For the pprof http server
	handlerInited bool
	pprofServer   *http.Server

	logDir    string   // log directory path
	vmLogFile *os.File // a file descriptor of the vmlog output file
}

// Verbosity sets the log verbosity ceiling. The verbosity of individual packages
// and source files can be raised using Vmodule.
func (*HandlerT) Verbosity(level int) error {
	return log.ChangeGlobalLogLevel(glogger, log.Lvl(level))
}

// VerbosityByName sets the verbosity of log module with given name.
// Please note that VerbosityByName only works with zapLogger.
func (*HandlerT) VerbosityByName(mn string, level int) error {
	return log.ChangeLogLevelWithName(mn, log.Lvl(level))
}

// VerbosityByID sets the verbosity of log module with given ModuleID.
// Please note that VerbosityByID only works with zapLogger.
func (*HandlerT) VerbosityByID(mi int, level int) error {
	return log.ChangeLogLevelWithID(log.ModuleID(mi), log.Lvl(level))
}

// Vmodule sets the log verbosity pattern. See package log for details on the
// pattern syntax.
func (*HandlerT) Vmodule(pattern string) error {
	return glogger.Vmodule(pattern)
}

// BacktraceAt sets the log backtrace location. See package log for details on
// the pattern syntax.
func (*HandlerT) BacktraceAt(location string) error {
	return glogger.BacktraceAt(location)
}

// MemStats returns detailed runtime memory statistics.
func (*HandlerT) MemStats() *runtime.MemStats {
	s := new(runtime.MemStats)
	runtime.ReadMemStats(s)
	return s
}

// GcStats returns GC statistics.
func (*HandlerT) GcStats() *debug.GCStats {
	s := new(debug.GCStats)
	debug.ReadGCStats(s)
	return s
}

// StartPProf starts the pprof server.
func (h *HandlerT) StartPProf(ptrAddr *string, ptrPort *int) error {
	// Set the default server address and port if they are not set
	var (
		address string
		port    int
	)
	if ptrAddr == nil || *ptrAddr == "" {
		address = pprofAddrFlag.Value
	} else {
		address = *ptrAddr
	}

	if ptrPort == nil || *ptrPort == 0 {
		port = pprofPortFlag.Value
	} else {
		port = *ptrPort
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.pprofServer != nil {
		return errors.New("pprof server is already running")
	}

	serverAddr := fmt.Sprintf("%s:%d", address, port)
	httpServer := &http.Server{Addr: serverAddr}

	if !h.handlerInited {
		// Hook go-metrics into expvar on any /debug/metrics request, load all vars
		// from the registry into expvar, and execute regular expvar handler.
		exp.Exp(metrics.DefaultRegistry)
		http.Handle("/memsize/", http.StripPrefix("/memsize", &Memsize))
		h.handlerInited = true
	}

	logger.Info("Starting pprof server", "addr", fmt.Sprintf("http://%s/debug/pprof", serverAddr))
	go func(handle *HandlerT) {
		if err := httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logger.Info("pprof server is closed")
			} else {
				logger.Error("Failure in running pprof server", "err", err)
			}
		}
		h.mu.Lock()
		h.pprofServer = nil
		h.mu.Unlock()
	}(h)

	h.pprofServer = httpServer

	return nil
}

// StopPProf stops the pprof server.
func (h *HandlerT) StopPProf() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.pprofServer == nil {
		return errors.New("pprof server is not running")
	}

	logger.Info("Shutting down pprof server")
	h.pprofServer.Close()

	return nil
}

// IsPProfRunning returns true if the pprof HTTP server is running and false otherwise.
func (h *HandlerT) IsPProfRunning() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.pprofServer != nil
}

// CpuProfile turns on CPU profiling for nsec seconds and writes
// profile data to file.
func (h *HandlerT) CpuProfile(file string, nsec uint) error {
	if err := h.StartCPUProfile(file); err != nil {
		return err
	}
	time.Sleep(time.Duration(nsec) * time.Second)
	h.StopCPUProfile()
	return nil
}

// StartCPUProfile turns on CPU profiling, writing to the given file.
func (h *HandlerT) StartCPUProfile(file string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.cpuW != nil {
		return errors.New("CPU profiling already in progress")
	}
	f, err := os.Create(expandHome(file))
	if err != nil {
		return err
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return err
	}
	h.cpuW = f
	h.cpuFile = file
	logger.Info("CPU profiling started", "dump", h.cpuFile)
	return nil
}

// StopCPUProfile stops an ongoing CPU profile.
func (h *HandlerT) StopCPUProfile() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	pprof.StopCPUProfile()
	if h.cpuW == nil {
		return errors.New("CPU profiling not in progress")
	}
	logger.Info("Done writing CPU profile", "dump", h.cpuFile)
	h.cpuW.Close()
	h.cpuW = nil
	h.cpuFile = ""
	return nil
}

// GoTrace turns on tracing for nsec seconds and writes
// trace data to file.
func (h *HandlerT) GoTrace(file string, nsec uint) error {
	if err := h.StartGoTrace(file); err != nil {
		return err
	}
	time.Sleep(time.Duration(nsec) * time.Second)
	h.StopGoTrace()
	return nil
}

// BlockProfile turns on goroutine profiling for nsec seconds and writes profile data to
// file. It uses a profile rate of 1 for most accurate information. If a different rate is
// desired, set the rate and write the profile manually.
func (*HandlerT) BlockProfile(file string, nsec uint) error {
	runtime.SetBlockProfileRate(1)
	time.Sleep(time.Duration(nsec) * time.Second)
	defer runtime.SetBlockProfileRate(0)
	return writeProfile("block", file)
}

// SetBlockProfileRate sets the rate of goroutine block profile data collection.
// rate 0 disables block profiling.
func (*HandlerT) SetBlockProfileRate(rate int) {
	runtime.SetBlockProfileRate(rate)
}

// WriteBlockProfile writes a goroutine blocking profile to the given file.
func (*HandlerT) WriteBlockProfile(file string) error {
	return writeProfile("block", file)
}

// MutexProfile turns on mutex profiling for nsec seconds and writes profile data to file.
// It uses a profile rate of 1 for most accurate information. If a different rate is
// desired, set the rate and write the profile manually.
func (*HandlerT) MutexProfile(file string, nsec uint) error {
	runtime.SetMutexProfileFraction(1)
	time.Sleep(time.Duration(nsec) * time.Second)
	defer runtime.SetMutexProfileFraction(0)
	return writeProfile("mutex", file)
}

// SetMutexProfileFraction sets the rate of mutex profiling.
func (*HandlerT) SetMutexProfileFraction(rate int) {
	runtime.SetMutexProfileFraction(rate)
}

// WriteMutexProfile writes a goroutine blocking profile to the given file.
func (*HandlerT) WriteMutexProfile(file string) error {
	return writeProfile("mutex", file)
}

// WriteMemProfile writes an allocation profile to the given file.
// Note that the profiling rate cannot be set through the API,
// it must be set on the command line.
func (*HandlerT) WriteMemProfile(file string) error {
	return writeProfile("heap", file)
}

// Stacks returns a printed representation of the stacks of all goroutines.
func (*HandlerT) Stacks() string {
	buf := make([]byte, 1024*1024)
	buf = buf[:runtime.Stack(buf, true)]
	return string(buf)
}

// FreeOSMemory returns unused memory to the OS.
func (*HandlerT) FreeOSMemory() {
	debug.FreeOSMemory()
}

// SetGCPercent sets the garbage collection target percentage. It returns the previous
// setting. A negative value disables GC.
func (*HandlerT) SetGCPercent(v int) int {
	return debug.SetGCPercent(v)
}

func writeProfile(name, file string) error {
	p := pprof.Lookup(name)
	logger.Info("Writing profile records", "count", p.Count(), "type", name, "dump", file)
	f, err := os.Create(expandHome(file))
	if err != nil {
		return err
	}
	defer f.Close()
	return p.WriteTo(f, 0)
}

// expands home directory in file paths.
// ~someuser/tmp will not be expanded.
func expandHome(p string) string {
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		home := os.Getenv("HOME")
		if home == "" {
			if usr, err := user.Current(); err == nil {
				home = usr.HomeDir
			}
		}
		if home != "" {
			p = home + p[1:]
		}
	}
	return filepath.Clean(p)
}

// WriteVMLog writes msg to a vmlog output file.
func (h *HandlerT) WriteVMLog(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.vmLogFile != nil {
		if _, err := h.vmLogFile.WriteString(msg + "\n"); err != nil {
			// Since vmlog is a debugging feature, write failure can be treated as a warning.
			logger.Warn("Failed to write to a vmlog file", "msg", msg, "err", err)
		}
	}
}

// openVMLogFile opens a file for vmlog output as the append mode.
func (h *HandlerT) openVMLogFile() {
	var err error
	filename := filepath.Join(h.logDir, "vm.log")
	Handler.vmLogFile, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		logger.Warn("Failed to open a file", "filename", filename, "err", err)
	}
}

func vmLogTargetToString(target int) string {
	switch target {
	case 0:
		return "no output"
	case params.VMLogToFile:
		return "file"
	case params.VMLogToStdout:
		return "stdout"
	case params.VMLogToAll:
		return "both file and stdout"
	default:
		return ""
	}
}

// SetVMLogTarget sets the output target of vmlog.
func (h *HandlerT) SetVMLogTarget(target int) (string, error) {
	if target < 0 || target > params.VMLogToAll {
		return vmLogTargetToString(params.VMLogTarget), fmt.Errorf("target should be between 0 and %d", params.VMLogToAll)
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if (target & params.VMLogToFile) != 0 {
		if h.vmLogFile == nil {
			h.openVMLogFile()
		}
	} else {
		if h.vmLogFile != nil {
			if err := Handler.vmLogFile.Close(); err != nil {
				logger.Warn("Failed to close the vmlog file", "err", err)
			}
			Handler.vmLogFile = nil
		}
	}

	params.VMLogTarget = target
	return vmLogTargetToString(target), nil
}
