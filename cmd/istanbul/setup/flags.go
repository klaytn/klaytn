// Copyright 2017 AMIS Technologies
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

package setup

import "github.com/urfave/cli"

var (
	numOfValidatorsFlag = cli.IntFlag{
		Name:  "num",
		Usage: "Number of validators",
	}

	verboseFlag = cli.BoolFlag{
		Name:  "verbose",
		Usage: "Print validator details",
	}

	staticNodesFlag = cli.BoolFlag{
		Name:  "nodes",
		Usage: "Print static nodes template",
	}

	bftFlag = cli.BoolFlag{
		Name:  "bft",
		Usage: "Use BFT",
	}

	saveFlag = cli.BoolFlag{
		Name:  "save",
		Usage: "Save to files",
	}
)
