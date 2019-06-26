// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

const module = "module"
const (
	ZapLogger     = "zap"
	Log15Logger   = "log15"
	DefaultLogger = Log15Logger
)

var baseLogger Logger

type Logger interface {
	NewWith(keysAndValues ...interface{}) Logger
	newModuleLogger(mi ModuleID) Logger
	Trace(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	ErrorWithStack(msg string, keysAndValues ...interface{})
	Crit(msg string, keysAndValues ...interface{})
	CritWithStack(msg string, keysAndValues ...interface{})

	// GetHandler gets the handler associated with the logger.
	GetHandler() Handler
	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h Handler)
}

func init() {
	root.SetHandler(DiscardHandler())
	SetBaseLogger()
}

func SetBaseLogger() {
	switch DefaultLogger {
	case ZapLogger:
		baseLogger = genBaseLoggerZap()
	case Log15Logger:
		baseLogger = root
	default:
		baseLogger = genBaseLoggerZap()
	}
}

func NewModuleLogger(mi ModuleID) Logger {
	newLogger := baseLogger.newModuleLogger(mi)
	return newLogger
}

// Fatalf formats a message to standard error and exits the program.
// The message is also printed to standard output if standard error
// is redirected to a different file.
func Fatalf(format string, args ...interface{}) {
	w := io.MultiWriter(os.Stdout, os.Stderr)
	if runtime.GOOS == "windows" {
		// The SameFile check below doesn't work on Windows.
		// stdout is unlikely to get redirected though, so just print there.
		w = os.Stdout
	} else {
		outf, _ := os.Stdout.Stat()
		errf, _ := os.Stderr.Stat()
		if outf != nil && errf != nil && os.SameFile(outf, errf) {
			w = os.Stderr
		}
	}
	fmt.Fprintf(w, "Fatal: "+format+"\n", args...)
	os.Exit(1)
}
