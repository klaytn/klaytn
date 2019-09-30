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

/*
Package log provides an opinionated, simple toolkit for best-practice logging.

Overview of log package

Logger is the interface for various implementation of logging.
There are log15 and zap logger and log15 is the default logger.

Source Files

  - format.go         : contains formatting functions for log15 logger
  - handler.go        : provides various types of log handling methods for log15 logger
  - handler_glog.go   : contains implementation of GlogHandler
  - handler_go13.go   : contains implementation of swapHandler in go1.3
  - handler_go14.go   : contains implementation of swapHandler in go1.4
  - handler_syslog.go : contains functions to use syslog package, but currently not used
  - interface.go      : defines Logger interface to support various implementations of loggers
  - log15_logger.go   : contains implementation of log15 logger, modified by go-ethereum
  - log_modules.go    : defines log modules, used to categorize logs
  - zap_logger.go     : contains functions and variables to use zap logger with Logger interface
*/
package log
