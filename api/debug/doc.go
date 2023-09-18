// Copyright 2018 The klaytn Authors
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

/*
Package debug provides interfaces for Go runtime debugging facilities.

Overview of debug package

This package is mostly glue code making these facilities available through the CLI and RPC subsystem.
If you want to use them from Go code, use package runtime instead.

Source Files

  - api.go                : defines the global debugging handler implementing debugging APIs.
  - flags.go              : defines command-line flags enabling debugging APIs.
  - loudpanic.go          : (deprecated) panics in a way that gets all goroutine stacks printed on stderr.
  - loudpanic_fallback.go : (deprecated) implements fallback of LoudPanic.
  - trace.go              : implements start/stop functions of go trace.
  - trace_fallback.go     : implements fallback of StartGoTrace and StopGoTrace for Go < 1.5.
*/
package debug
