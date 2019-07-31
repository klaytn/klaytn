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
// This file is derived from common/mclock/mclock.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package mclock is a wrapper for a monotonic clock source.

mclock package provides a Now() function which returns the current time in nanoseconds from a monotonic clock in Duration type.
The returned time is based on some arbitrary platform-specific point in the
past.  The returned time is guaranteed to increase monotonically at a
constant rate, unlike time.Now() from the Go standard library, which may
slow down, speed up, jump forward or backward, due to NTP activity or leap
seconds.
*/
package mclock
