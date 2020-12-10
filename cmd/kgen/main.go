// Copyright 2019 The klaytn Authors
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

package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"

	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/cmd/utils/nodecmd"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"gopkg.in/urfave/cli.v1"
)

type validatorInfo struct {
	Address  common.Address
	Nodekey  string
	NodeInfo string
}

const (
	dirKeys = "keys" // directory name where the created files are stored.
)

var (
	logger   = log.NewModuleLogger(log.CMDKGEN)
	fileFlag = cli.BoolFlag{
		Name:  "file",
		Usage: `Generate a nodekey and a Klaytn node information as files`,
	}
	portFlag = cli.IntFlag{
		Name:  "port",
		Usage: `Specify a tcp port number`,
		Value: 32323,
	}
	ipFlag = cli.StringFlag{
		Name:  "ip",
		Usage: `Specify an ip address`,
		Value: "0.0.0.0",
	}
)

func init() {
	cli.AppHelpTemplate = utils.KgenHelpTemplate
	cli.HelpPrinter = utils.NewHelpPrinter(nil)
}

func main() {
	app := cli.NewApp()
	app.Name = "kgen"
	app.Usage = "The command line interface to generate nodekey information for Klaytn"
	app.Copyright = "Copyright 2018-2019 The klaytn Authors"
	app.Action = genNodeKey
	app.Flags = []cli.Flag{
		fileFlag,
		ipFlag,
		portFlag,
	}
	app.Commands = []cli.Command{
		nodecmd.VersionCommand,
	}
	app.HideVersion = true
	//app.CustomAppHelpTemplate = kgenHelper
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// writeNodeKeyInfoToFile writes `nodekey` and `node_info.json` as files under the `parentDir` folder.
// The validator is a json format file containing address, nodekey and nodeinfo.
func writeNodeKeyInfoToFile(validator *validatorInfo, parentDir string, nodekey string) error {
	parentPath := path.Join("", parentDir)
	err := os.MkdirAll(parentPath, os.ModePerm)
	if err != nil {
		return err
	}

	nodeKeyFilePath := path.Join(parentPath, "nodekey")
	if err = ioutil.WriteFile(nodeKeyFilePath, []byte(nodekey), os.ModePerm); err != nil {
		return err
	}
	fmt.Println("Created : ", nodeKeyFilePath)

	str, err := json.MarshalIndent(validator, "", "\t")
	if err != nil {
		return err
	}
	validatorInfoFilePath := path.Join(parentPath, "node_info.json")
	if err = ioutil.WriteFile(validatorInfoFilePath, []byte(str), os.ModePerm); err != nil {
		return err
	}

	fmt.Println("Created : ", validatorInfoFilePath)
	return nil
}

// makeNodeInfo creates a validator with the given parameters.
func makeNodeInfo(nodeAddr common.Address, nodeKey string, privKey *ecdsa.PrivateKey, ip string, port uint16) *validatorInfo {
	return &validatorInfo{
		Address: nodeAddr,
		Nodekey: nodeKey,
		NodeInfo: discover.NewNode(
			discover.PubkeyID(&privKey.PublicKey),
			net.ParseIP(ip),
			0,
			port,
			nil,
			discover.NodeTypeUnknown).String(),
	}
}

// genNodeKey creates a validator which is printed as json format or is stored into files(nodekey, validator).
func genNodeKey(ctx *cli.Context) error {
	pk, nk, addr, err := generateNodeInfoContents()
	if err != nil {
		return err
	}
	ip := ctx.String(ipFlag.Name)
	if net.ParseIP(ip).To4() == nil {
		return fmt.Errorf("IP address is not valid")
	}
	port := ctx.Uint(portFlag.Name)
	if port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	nodeinfo := makeNodeInfo(addr, nk, pk, ip, uint16(port))
	if ctx.Bool(fileFlag.Name) {
		if err := writeNodeKeyInfoToFile(nodeinfo, dirKeys, nk); err != nil {
			return err
		}
	} else {
		str, err := json.MarshalIndent(nodeinfo, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(str))
	}
	return nil
}

// randomBytes creates random bytes as long as `len`.
func randomBytes(len int) ([]byte, error) {
	b := make([]byte, len)
	_, _ = rand.Read(b)

	return b, nil
}

// randomHex creates a random 32-bytes hexadecimal string.
func randomHex() string {
	b, _ := randomBytes(32)
	return common.BytesToHash(b).Hex()
}

// generateNodeInfoContents generates contents of a validator.
func generateNodeInfoContents() (*ecdsa.PrivateKey, string, common.Address, error) {
	nodekey := randomHex()[2:] // truncate `0x` prefix

	key, err := crypto.HexToECDSA(nodekey)
	if err != nil {
		logger.Error("Failed to generate key", "err", err)
		return nil, "", common.Address{}, err
	}

	addr := crypto.PubkeyToAddress(key.PublicKey)

	return key, nodekey, addr, nil
}
