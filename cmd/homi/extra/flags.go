// Copyright 2018 The klaytn Authors
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

package extra

import (
	"strings"

	"gopkg.in/urfave/cli.v1"
)

var (
	configFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}

	extraDataFlag = cli.StringFlag{
		Name:  "extradata",
		Usage: "Hex string for RLP encoded Istanbul extraData",
	}

	validatorsFlag = cli.StringFlag{
		Name:  "validators",
		Usage: "Validators for RLP encoded Istanbul extraData",
	}

	vanityFlag = cli.StringFlag{
		Name:  "vanity",
		Usage: "Vanity for RLP encoded Istanbul extraData",
		Value: "0x00",
	}
)

func splitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
}
