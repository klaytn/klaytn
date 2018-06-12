package setup

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/p2p/discover"
	istcommon "ground-x/go-gxplatform/cmd/istanbul/common"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"path"
	"strconv"
	"ground-x/go-gxplatform/cmd/istanbul/genesis"
)

type validatorInfo struct {
	Address  common.Address
	Nodekey  string
	NodeInfo string
}

var (
	SetupCommand = cli.Command{
		Name:  "setup",
		Usage: "Setup your Istanbul network in seconds",
		Description: `This tool helps generate:
		* Genesis block
		* Static nodes for all validators
		* Validator details
	    for Istanbul consensus.
`,
		Action: gen,
		Flags: []cli.Flag{
			numOfValidatorsFlag,
			staticNodesFlag,
			verboseFlag,
			bftFlag,
			saveFlag,
		},
	}
)

func gen(ctx *cli.Context) error {
	num := ctx.Int(numOfValidatorsFlag.Name)

	keys, nodekeys, addrs := istcommon.GenerateKeys(num)
	var nodes []string

	if ctx.Bool(verboseFlag.Name) {
		fmt.Println("validators")
	}

	for i := 0; i < num; i++ {
		v := &validatorInfo{
			Address: addrs[i],
			Nodekey: nodekeys[i],
			NodeInfo: discover.NewNode(
				discover.PubkeyID(&keys[i].PublicKey),
				net.ParseIP("0.0.0.0"),
				0,
				uint16(30303) + uint16(i) * uint16(10101)).String(),
		}

		nodes = append(nodes, string(v.NodeInfo))

		if ctx.Bool(verboseFlag.Name) {
			str, _ := json.MarshalIndent(v, "", "\t")
			fmt.Println(string(str))

			if ctx.Bool(saveFlag.Name) {
				folderName := strconv.Itoa(i)
				os.MkdirAll(folderName, os.ModePerm)
				ioutil.WriteFile(path.Join(folderName, "nodekey"), []byte(nodekeys[i]), os.ModePerm)
			}
		}
	}

	if ctx.Bool(verboseFlag.Name) {
		fmt.Print("\n\n\n")
	}

	staticNodes, _ := json.MarshalIndent(nodes, "", "\t")
	if ctx.Bool(staticNodesFlag.Name) {
		name := "static-nodes.json"
		fmt.Println(name)
		fmt.Println(string(staticNodes))
		fmt.Print("\n\n\n")

		if ctx.Bool(saveFlag.Name) {
			ioutil.WriteFile(name, staticNodes, os.ModePerm)
		}
	}

	var jsonBytes []byte
	isBFT := ctx.Bool(bftFlag.Name)
	g := genesis.New(
		genesis.Validators(addrs...),
		genesis.Alloc(addrs, new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil)),
	)

	if isBFT {
		jsonBytes, _ = json.MarshalIndent(genesis.ToBFT(g, true), "", "    ")
	} else {
		jsonBytes, _ = json.MarshalIndent(g, "", "    ")
	}
	fmt.Println("genesis.json")
	fmt.Println(string(jsonBytes))

	if ctx.Bool(saveFlag.Name) {
		ioutil.WriteFile("genesis.json", jsonBytes, os.ModePerm)
	}

	return nil
}
