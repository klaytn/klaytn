// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from quorum/consensus/istanbul/config.go (2018/06/04).
// Modified and improved for the klaytn development.

package istanbul

type ProposerPolicy uint64

const (
	RoundRobin ProposerPolicy = iota
	Sticky
	WeightedRandom
)

type Config struct {
	Timeout        uint64         `toml:",omitempty"` // The timeout for each Istanbul round in milliseconds.
	BlockPeriod    uint64         `toml:",omitempty"` // Default minimum difference between two consecutive block's timestamps in second
	ProposerPolicy ProposerPolicy `toml:",omitempty"` // The policy for proposer selection
	Epoch          uint64         `toml:",omitempty"` // The number of blocks after which to checkpoint and reset the pending votes
	SubGroupSize   uint64         `toml:",omitempty"`
}

// TODO-Klaytn-Istanbul: Do not use DefaultConfig except for assigning new config
var DefaultConfig = &Config{
	Timeout:        10000,
	BlockPeriod:    1,
	ProposerPolicy: RoundRobin,
	Epoch:          30000,
	SubGroupSize:   21,
}
