package extra

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"ground-x/go-gxplatform/common"
	"github.com/naoina/toml"
)

var (
	ExtraCommand = cli.Command{
		Name:  "extra",
		Usage: "Istanbul extraData manipulation",
		Subcommands: []cli.Command{
			cli.Command{
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
			cli.Command{
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