package main

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"ground-x/go-gxplatform/cmd/utils"
	gxplatform "ground-x/go-gxplatform/gxp"
	"ground-x/go-gxplatform/node"
	"ground-x/go-gxplatform/params"
	"os"
	"reflect"
	"unicode"

	"github.com/naoina/toml"
)

var (
	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

type gxpConfig struct {
	Gxp  gxplatform.Config
	Node node.Config
}

func loadConfig(file string, cfg *gxpConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit)
	cfg.HTTPModules = append(cfg.HTTPModules, "gxp", "shh")
	cfg.WSModules = append(cfg.WSModules, "gxp", "shh")
	cfg.IPCPath = "gxp.ipc"
	return cfg
}

func makeConfigNode(ctx *cli.Context) (*node.Node, gxpConfig) {
	// Load defaults.
	cfg := gxpConfig{
		Gxp:  gxplatform.DefaultConfig,
		Node: defaultNodeConfig(),
	}

	// Load config file.
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)
	stack, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}

	return stack, cfg
}

func makeFullNode(ctx *cli.Context) *node.Node {
	stack, cfg := makeConfigNode(ctx)

	utils.RegisterGxpService(stack, &cfg.Gxp)

	return stack
}
