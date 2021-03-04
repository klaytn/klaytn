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
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/klaytn/klaytn/common"
	"github.com/naoina/toml"
)

var (
	ExtraCommand = cli.Command{
		Name:  "extra",
		Usage: "Istanbul extraData manipulation",
		Subcommands: []cli.Command{
			{
				Action:    decode,
				Name:      "decode",
				Usage:     "To decode an Istanbul extraData",
				ArgsUsage: "<extra data>",
				Flags: []cli.Flag{
					extraDataFlag,
				},
				Description: `
		This command decodes extraData to vanity and validators.
		`,
			},
			{
				Action:    encode,
				Name:      "encode",
				Usage:     "To encode an Istanbul extraData",
				ArgsUsage: "<config file> or \"0xValidator1,0xValidator2...\"",
				Flags: []cli.Flag{
					configFlag,
					validatorsFlag,
					vanityFlag,
				},
				Description: `
		This command encodes vanity and validators to extraData. Please refer to example/config.toml.
		`,
			},
		},
	}
)

func encode(ctx *cli.Context) error {
	path := ctx.String(configFlag.Name)
	validators := ctx.String(validatorsFlag.Name)
	if len(path) == 0 && len(validators) == 0 {
		return cli.NewExitError("Must supply config file or enter validators", 0)
	}

	if len(path) != 0 {
		extraData, err := fromConfig(path)
		if err != nil {
			return cli.NewExitError("Failed to encode from config data", 0)
		}
		fmt.Println("Encoded Istanbul extra-data:", extraData)
	}

	if len(validators) != 0 {
		extraData, err := fromRawData(ctx.String(vanityFlag.Name), validators)
		if err != nil {
			return cli.NewExitError("Failed to encode from flags", 0)
		}
		fmt.Println("Encoded Istanbul extra-data:", extraData)
	}
	return nil
}

func fromRawData(vanity string, validators string) (string, error) {
	vs := splitAndTrim(validators)

	addrs := make([]common.Address, len(vs))
	for i, v := range vs {
		addrs[i] = common.HexToAddress(v)
	}
	return Encode(vanity, addrs)
}

func fromConfig(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", cli.NewExitError(fmt.Sprintf("Failed to read config file: %v", err), 1)
	}
	defer file.Close()

	var config struct {
		Vanity     string
		Validators []common.Address
	}

	if err := toml.NewDecoder(file).Decode(&config); err != nil {
		return "", cli.NewExitError(fmt.Sprintf("Failed to parse config file: %v", err), 2)
	}

	return Encode(config.Vanity, config.Validators)
}

func decode(ctx *cli.Context) error {
	if !ctx.IsSet(extraDataFlag.Name) {
		return cli.NewExitError("Must supply extra data", 10)
	}

	extraString := ctx.String(extraDataFlag.Name)
	vanity, istanbulExtra, err := Decode(extraString)
	if err != nil {
		return err
	}

	fmt.Println("vanity: ", "0x"+common.Bytes2Hex(vanity))

	for _, v := range istanbulExtra.Validators {
		fmt.Println("validator: ", v.Hex())
	}

	if len(istanbulExtra.Seal) != 0 {
		fmt.Println("seal:", "0x"+common.Bytes2Hex(istanbulExtra.Seal))
	}

	for _, seal := range istanbulExtra.CommittedSeal {
		fmt.Println("committed seal: ", "0x"+common.Bytes2Hex(seal))
	}

	return nil
}
