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
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zlManager = zapLoggerManager{"stderr", // Use stderr to outputPath instead of stdout to be aligned with log15.
	"json", zapcore.InfoLevel,
	sync.Mutex{}, make(map[ModuleID][]*zapLogger),
}

type zapLoggerManager struct {
	outputPath   string
	encodingType string
	logLevel     zapcore.Level
	mutex        sync.Mutex
	loggersMap   map[ModuleID][]*zapLogger
}

type zapLogger struct {
	mi  ModuleID
	cfg *zap.Config
	sl  *zap.SugaredLogger
}

// A zapLogger generated from NewWith inherits InitialFields and ModuleID from its parent.
func (zl *zapLogger) NewWith(keysAndValues ...interface{}) Logger {
	zlManager.mutex.Lock()
	defer zlManager.mutex.Unlock()

	newCfg := genDefaultConfig()
	for k, v := range zl.cfg.InitialFields {
		newCfg.InitialFields[k] = v
	}
	newCfg.Level.SetLevel(zl.cfg.Level.Level())
	return genLoggerZap(zl.mi, newCfg)
}

func (zl *zapLogger) newModuleLogger(mi ModuleID) Logger {
	zlManager.mutex.Lock()
	defer zlManager.mutex.Unlock()

	zapCfg := genDefaultConfig()
	zapCfg.InitialFields[module] = mi
	return genLoggerZap(mi, zapCfg)
}

func (zl *zapLogger) Trace(msg string, keysAndValues ...interface{}) {
	zl.sl.Debugw(msg, keysAndValues...)
}

func (zl *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	zl.sl.Debugw(msg, keysAndValues...)
}

func (zl *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	zl.sl.Infow(msg, keysAndValues...)
}

func (zl *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	zl.sl.Warnw(msg, keysAndValues...)
}

func (zl *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	zl.sl.Errorw(msg, keysAndValues...)
}

func (zl *zapLogger) ErrorWithStack(msg string, keysAndValues ...interface{}) {
	zl.sl.Errorw(msg, keysAndValues...)
}

func (zl *zapLogger) Crit(msg string, keysAndValues ...interface{}) {
	zl.sl.Fatalw(msg, keysAndValues...)
}

func (zl *zapLogger) CritWithStack(msg string, keysAndValues ...interface{}) {
	zl.sl.Fatalw(msg, keysAndValues...)
}

// GetHandler and SetHandler do nothing but exist to make consistency in Logger interface.
func (zl *zapLogger) GetHandler() Handler {
	return nil
}

func (zl *zapLogger) SetHandler(h Handler) {
}

func (zl *zapLogger) setLevel(lvl Lvl) {
	zl.cfg.Level.SetLevel(lvlToZapLevel(lvl))
}

// register registers the receiver to zapLoggerManager.
func (zl *zapLogger) register() {
	zlManager.loggersMap[zl.mi] = append(zlManager.loggersMap[zl.mi], zl)
}

func genBaseLoggerZap() Logger {
	return genLoggerZap(BaseLogger, genDefaultConfig())
}

// genLoggerZap creates a zapLogger with given ModuleID and Config.
func genLoggerZap(mi ModuleID, cfg *zap.Config) Logger {
	logger, err := cfg.Build()
	if err != nil {
		Fatalf("Error while building zapLogger from the config. ModuleID: %v, err: %v", mi, err)
	}
	newLogger := &zapLogger{mi, cfg, logger.Sugar()}
	newLogger.register()
	return newLogger
}

// ChangeLogLevelWithName changes the log level of loggers with given ModuleName.
func ChangeLogLevelWithName(moduleName string, lvl Lvl) error {
	if err := levelCheck(lvl); err != nil {
		return err
	}
	mi := GetModuleID(moduleName)
	if mi == ModuleNameLen {
		return errors.New("entered module name does not match with any existing log module")
	}
	return ChangeLogLevelWithID(mi, lvl)
}

// ChangeLogLevelWithName changes the log level of loggers with given ModuleID.
func ChangeLogLevelWithID(mi ModuleID, lvl Lvl) error {
	if err := levelCheck(lvl); err != nil {
		return err
	}
	if err := idCheck(mi); err != nil {
		return err
	}
	loggers := zlManager.loggersMap[mi]
	for _, logger := range loggers {
		logger.setLevel(lvl)
	}
	return nil
}

func ChangeGlobalLogLevel(glogger *GlogHandler, lvl Lvl) error {
	if err := levelCheck(lvl); err != nil {
		return err
	}
	for _, loggers := range zlManager.loggersMap {
		for _, logger := range loggers {
			logger.setLevel(lvl)
		}
	}

	if glogger != nil {
		glogger.Verbosity(lvl)
	}
	return nil
}

func levelCheck(lvl Lvl) error {
	if lvl >= LvlEnd {
		return errors.New(fmt.Sprintf("insert log level less than %d", LvlEnd))
	}
	if lvl < LvlCrit {
		return errors.New(fmt.Sprintf("insert log level greater than or equal to %d", LvlCrit))
	}
	return nil
}

func idCheck(mi ModuleID) error {
	if mi >= ModuleNameLen {
		return errors.New(fmt.Sprintf("insert log level less than %d", ModuleNameLen))
	}
	if mi <= BaseLogger {
		return errors.New(fmt.Sprintf("insert log level greater than %d", BaseLogger))
	}
	return nil
}

func lvlToZapLevel(lvl Lvl) zapcore.Level {
	switch lvl {
	case LvlCrit:
		return zapcore.FatalLevel
	case LvlError:
		return zapcore.ErrorLevel
	case LvlWarn:
		return zapcore.WarnLevel
	case LvlInfo:
		return zapcore.InfoLevel
	case LvlDebug:
		return zapcore.DebugLevel
	case LvlTrace:
		return zapcore.DebugLevel
	default:
		baseLogger.Error("Unexpected log level entered. Use InfoLevel instead.", "entered level", lvl)
		return zapcore.InfoLevel
	}
}

func genDefaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:  "ts",
		LevelKey: "level",
		NameKey:  "logger",
		//CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		//EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func genDefaultConfig() *zap.Config {
	encoderConfig := genDefaultEncoderConfig()
	return &zap.Config{
		Encoding:      zlManager.encodingType,
		Level:         zap.NewAtomicLevelAt(zlManager.logLevel),
		OutputPaths:   []string{zlManager.outputPath},
		Development:   false,
		EncoderConfig: encoderConfig,
		InitialFields: make(map[string]interface{}),
	}
}
