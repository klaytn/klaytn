package common

import (
	"os"
	"github.com/inconshreveable/log15"
)

var log = New()
var defaultLogger = log15.New()

func init() {
	defaultLogger.SetHandler(
		log15.MultiHandler(
			log15.CallerFileHandler(log15.StreamHandler(
				os.Stdout, log15.TerminalFormat(),
			)),
		),
	)
}

func New(ctx ...interface{}) log15.Logger {
	return defaultLogger.New(ctx...)
}