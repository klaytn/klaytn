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

package setup

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/blockchain"
	istcommon "github.com/klaytn/klaytn/cmd/homi/common"
	"github.com/klaytn/klaytn/cmd/homi/docker/compose"
	"github.com/klaytn/klaytn/cmd/homi/docker/service"
	"github.com/klaytn/klaytn/cmd/homi/genesis"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/params"
	"gopkg.in/urfave/cli.v1"
)

type ValidatorInfo struct {
	Address  common.Address
	Nodekey  string
	NodeInfo string
}

type GrafanaFile struct {
	url  string
	name string
}

var (
	SetupCommand = cli.Command{
		Name:  "setup",
		Usage: "Generate klaytn CN's init files",
		Description: `This tool helps generate:
		* Genesis Block (genesis.json)
		* Static nodes for all CNs(Consensus Node)
		* CN details
		* Docker-compose

		for Klaytn Consensus Node.

Args :
		type : [local | remote | deploy | docker (default)]
`,
		Action: gen,
		Flags: []cli.Flag{
			cypressTestFlag,
			cypressFlag,
			baobabTestFlag,
			baobabFlag,
			serviceChainFlag,
			serviceChainTestFlag,
			cliqueFlag,
			numOfCNsFlag,
			numOfValidatorsFlag,
			numOfPNsFlag,
			numOfENsFlag,
			numOfSCNsFlag,
			numOfSPNsFlag,
			numOfSENsFlag,
			numOfTestKeyFlag,
			chainIDFlag,
			serviceChainIDFlag,
			unitPriceFlag,
			deriveShaImplFlag,
			fundingAddrFlag,
			outputPathFlag,
			dockerImageIdFlag,
			fasthttpFlag,
			networkIdFlag,
			nografanaFlag,
			useTxGenFlag,
			txGenRateFlag,
			txGenThFlag,
			txGenConnFlag,
			txGenDurFlag,
			rpcPortFlag,
			wsPortFlag,
			p2pPortFlag,
			dataDirFlag,
			logDirFlag,
			governanceFlag,
			govModeFlag,
			governingNodeFlag,
			rewardMintAmountFlag,
			rewardRatioFlag,
			rewardGiniCoeffFlag,
			rewardStakingFlag,
			rewardProposerFlag,
			rewardMinimumStakeFlag,
			rewardDeferredTxFeeFlag,
			istEpochFlag,
			istProposerPolicyFlag,
			istSubGroupFlag,
			cliqueEpochFlag,
			cliquePeriodFlag,
			istanbulCompatibleBlockNumberFlag,
			londonCompatibleBlockNumberFlag,
			ethTxTypeCompatibleBlockNumberFlag,
		},
		ArgsUsage: "type",
	}
)

const (
	baobabOperatorAddress = "0x79deccfacd0599d3166eb76972be7bb20f51b46f"
	baobabOperatorKey     = "199fd187fdb2ce5f577797e1abaf4dd50e62275949c021f0112be40c9721e1a2"
)
const (
	DefaultTcpPort uint16 = 32323
	TypeNotDefined        = -1
	TypeDocker            = 0
	TypeLocal             = 1
	TypeRemote            = 2
	TypeDeploy            = 3
	DirScript             = "scripts"
	DirKeys               = "keys"
	DirPnScript           = "scripts_pn"
	DirPnKeys             = "keys_pn"
	DirTestKeys           = "keys_test"
	CNIpNetwork           = "10.11.2"
	PNIpNetwork1          = "10.11.10"
	PNIpNetwork2          = "10.11.11"
)

var Types = [4]string{"docker", "local", "remote", "deploy"}

var GrafanaFiles = [...]GrafanaFile{
	{
		url:  "https://s3.ap-northeast-2.amazonaws.com/klaytn-tools/Klaytn.json",
		name: "Klaytn.json",
	},
	{
		url:  "https://s3.ap-northeast-2.amazonaws.com/klaytn-tools/klaytn_txpool.json",
		name: "Klaytn_txpool.json",
	},
}

var lastIssuedPortNum = DefaultTcpPort

func genRewardConfig(ctx *cli.Context) *params.RewardConfig {
	mintingAmount := new(big.Int)
	mintingAmountString := ctx.String(rewardMintAmountFlag.Name)
	if _, ok := mintingAmount.SetString(mintingAmountString, 10); !ok {
		log.Fatalf("Minting amount must be a number", "value", mintingAmountString)
	}
	ratio := ctx.String(rewardRatioFlag.Name)
	giniCoeff := ctx.Bool(rewardGiniCoeffFlag.Name)
	deferredTxFee := ctx.Bool(rewardDeferredTxFeeFlag.Name)
	stakingInterval := ctx.Uint64(rewardStakingFlag.Name)
	proposalInterval := ctx.Uint64(rewardProposerFlag.Name)
	minimumStake := new(big.Int)
	minimumStakeString := ctx.String(rewardMinimumStakeFlag.Name)
	if _, ok := minimumStake.SetString(minimumStakeString, 10); !ok {
		log.Fatalf("Minimum stake must be a number", "value", minimumStakeString)
	}

	return &params.RewardConfig{
		MintingAmount:          mintingAmount,
		Ratio:                  ratio,
		UseGiniCoeff:           giniCoeff,
		DeferredTxFee:          deferredTxFee,
		StakingUpdateInterval:  stakingInterval,
		ProposerUpdateInterval: proposalInterval,
		MinimumStake:           minimumStake,
	}
}

func genIstanbulConfig(ctx *cli.Context) *params.IstanbulConfig {
	epoch := ctx.Uint64(istEpochFlag.Name)
	policy := ctx.Uint64(istProposerPolicyFlag.Name)
	subGroup := ctx.Uint64(istSubGroupFlag.Name)

	return &params.IstanbulConfig{
		Epoch:          epoch,
		ProposerPolicy: policy,
		SubGroupSize:   subGroup,
	}
}

func genGovernanceConfig(ctx *cli.Context) *params.GovernanceConfig {
	govMode := ctx.String(govModeFlag.Name)
	governingNode := ctx.String(governingNodeFlag.Name)
	if !common.IsHexAddress(governingNode) {
		log.Fatalf("Governing Node is invalid hex address", "value", governingNode)
	}
	return &params.GovernanceConfig{
		GoverningNode:  common.HexToAddress(governingNode),
		GovernanceMode: govMode,
		Reward:         genRewardConfig(ctx),
	}
}

func genCliqueConfig(ctx *cli.Context) *params.CliqueConfig {
	epoch := ctx.Uint64(cliqueEpochFlag.Name)
	period := ctx.Uint64(cliquePeriodFlag.Name)

	return &params.CliqueConfig{
		Period: period,
		Epoch:  epoch,
	}
}

func genIstanbulGenesis(ctx *cli.Context, nodeAddrs, testAddrs []common.Address, chainId uint64) *blockchain.Genesis {
	unitPrice := ctx.Uint64(unitPriceFlag.Name)
	chainID := new(big.Int).SetUint64(chainId)
	deriveShaImpl := ctx.Int(deriveShaImplFlag.Name)

	config := genGovernanceConfig(ctx)
	if len(nodeAddrs) > 0 && config.GoverningNode.String() == params.DefaultGoverningNode {
		config.GoverningNode = nodeAddrs[0]
	}

	options := []genesis.Option{
		genesis.Validators(nodeAddrs...),
		genesis.Alloc(append(nodeAddrs, testAddrs...), new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil)),
		genesis.DeriveShaImpl(deriveShaImpl),
		genesis.UnitPrice(unitPrice),
		genesis.ChainID(chainID),
	}

	if ok := ctx.Bool(governanceFlag.Name); ok {
		options = append(options, genesis.Governance(config))
	}
	options = append(options, genesis.Istanbul(genIstanbulConfig(ctx)))

	return genesis.New(options...)
}

func genCliqueGenesis(ctx *cli.Context, nodeAddrs, testAddrs []common.Address, chainId uint64) *blockchain.Genesis {
	config := genCliqueConfig(ctx)
	unitPrice := ctx.Uint64(unitPriceFlag.Name)
	chainID := new(big.Int).SetUint64(chainId)

	if ok := ctx.Bool(governanceFlag.Name); ok {
		log.Fatalf("Currently, governance is not supported for clique consensus", "--governance", ok)
	}

	genesisJson := genesis.NewClique(
		genesis.ValidatorsOfClique(nodeAddrs...),
		genesis.Alloc(append(nodeAddrs, testAddrs...), new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil)),
		genesis.UnitPrice(unitPrice),
		genesis.ChainID(chainID),
		genesis.Clique(config),
	)
	return genesisJson
}

func genValidatorKeystore(privKeys []*ecdsa.PrivateKey) {
	path := path.Join(outputPath, DirKeys)
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)

	for i, pk := range privKeys {
		pwdStr := RandStringRunes(params.PasswordLength)
		ks.ImportECDSA(pk, pwdStr)
		WriteFile([]byte(pwdStr), DirKeys, "passwd"+strconv.Itoa(i+1))
	}
}

func genCypressCommonGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	mintingAmount, _ := new(big.Int).SetString("9600000000000000000", 10)
	genesisJson := &blockchain.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		BlockScore: big.NewInt(genesis.InitBlockScore),
		Alloc:      make(blockchain.GenesisAlloc),
		Config: &params.ChainConfig{
			ChainID:       big.NewInt(10000),
			DeriveShaImpl: 2,
			Governance: &params.GovernanceConfig{
				GoverningNode:  nodeAddrs[0],
				GovernanceMode: "single",
				Reward: &params.RewardConfig{
					MintingAmount: mintingAmount,
					Ratio:         "34/54/12",
					UseGiniCoeff:  true,
					DeferredTxFee: true,
				},
			},
			Istanbul: &params.IstanbulConfig{
				ProposerPolicy: 2,
				SubGroupSize:   22,
			},
			UnitPrice: 25000000000,
		},
	}
	assignExtraData := genesis.Validators(nodeAddrs...)
	assignExtraData(genesisJson)

	return genesisJson
}

func genCypressGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	genesisJson := genCypressCommonGenesis(nodeAddrs, testAddrs)
	genesisJson.Config.Istanbul.Epoch = 604800
	genesisJson.Config.Governance.Reward.StakingUpdateInterval = 86400
	genesisJson.Config.Governance.Reward.ProposerUpdateInterval = 3600
	genesisJson.Config.Governance.Reward.MinimumStake = new(big.Int).SetUint64(5000000)
	allocationFunction := genesis.AllocWithCypressContract(append(nodeAddrs, testAddrs...), new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil))
	allocationFunction(genesisJson)
	return genesisJson
}

func genServiceChainCommonGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	genesisJson := &blockchain.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		BlockScore: big.NewInt(genesis.InitBlockScore),
		Alloc:      make(blockchain.GenesisAlloc),
		Config: &params.ChainConfig{
			ChainID:       big.NewInt(1000),
			DeriveShaImpl: 2,
			Istanbul: &params.IstanbulConfig{
				ProposerPolicy: 0,
				SubGroupSize:   22,
			},
			UnitPrice: 0,
		},
	}
	assignExtraData := genesis.Validators(nodeAddrs...)
	assignExtraData(genesisJson)

	return genesisJson
}

func genServiceChainGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	genesisJson := genServiceChainCommonGenesis(nodeAddrs, testAddrs)
	genesisJson.Config.Istanbul.Epoch = 3600
	allocationFunction := genesis.Alloc(append(nodeAddrs, testAddrs...), new(big.Int).Exp(big.NewInt(10), big.NewInt(10), nil))
	allocationFunction(genesisJson)
	return genesisJson
}

func genServiceChainTestGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	genesisJson := genServiceChainCommonGenesis(nodeAddrs, testAddrs)
	genesisJson.Config.Istanbul.Epoch = 30
	allocationFunction := genesis.Alloc(append(nodeAddrs, testAddrs...), new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil))
	allocationFunction(genesisJson)
	return genesisJson
}

func genCypressTestGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	testGenesis := genCypressCommonGenesis(nodeAddrs, testAddrs)
	testGenesis.Config.Istanbul.Epoch = 30
	testGenesis.Config.Governance.Reward.StakingUpdateInterval = 60
	testGenesis.Config.Governance.Reward.ProposerUpdateInterval = 30
	testGenesis.Config.Governance.Reward.MinimumStake = new(big.Int).SetUint64(5000000)
	allocationFunction := genesis.AllocWithPrecypressContract(append(nodeAddrs, testAddrs...), new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil))
	allocationFunction(testGenesis)
	return testGenesis
}

func genBaobabCommonGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	mintingAmount, _ := new(big.Int).SetString("9600000000000000000", 10)
	genesisJson := &blockchain.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		BlockScore: big.NewInt(genesis.InitBlockScore),
		Alloc:      make(blockchain.GenesisAlloc),
		Config: &params.ChainConfig{
			ChainID:       big.NewInt(2019),
			DeriveShaImpl: 2,
			Governance: &params.GovernanceConfig{
				GoverningNode:  nodeAddrs[0],
				GovernanceMode: "single",
				Reward: &params.RewardConfig{
					MintingAmount: mintingAmount,
					Ratio:         "34/54/12",
					UseGiniCoeff:  false,
					DeferredTxFee: true,
				},
			},
			Istanbul: &params.IstanbulConfig{
				ProposerPolicy: 2,
				SubGroupSize:   13,
			},
			UnitPrice: 25000000000,
		},
	}
	assignExtraData := genesis.Validators(nodeAddrs...)
	assignExtraData(genesisJson)

	return genesisJson
}

func genBaobabGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	genesisJson := genBaobabCommonGenesis(nodeAddrs, testAddrs)
	genesisJson.Config.Istanbul.Epoch = 604800
	genesisJson.Config.Governance.Reward.StakingUpdateInterval = 86400
	genesisJson.Config.Governance.Reward.ProposerUpdateInterval = 3600
	genesisJson.Config.Governance.Reward.MinimumStake = new(big.Int).SetUint64(5000000)
	allocationFunction := genesis.AllocWithBaobabContract(append(nodeAddrs, testAddrs...), new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil))
	allocationFunction(genesisJson)
	return genesisJson
}

func genBaobabTestGenesis(nodeAddrs, testAddrs []common.Address) *blockchain.Genesis {
	testGenesis := genBaobabCommonGenesis(nodeAddrs, testAddrs)
	testGenesis.Config.Istanbul.Epoch = 30
	testGenesis.Config.Governance.Reward.StakingUpdateInterval = 60
	testGenesis.Config.Governance.Reward.ProposerUpdateInterval = 30
	testGenesis.Config.Governance.Reward.MinimumStake = new(big.Int).SetUint64(5000000)
	allocationFunction := genesis.AllocWithPrebaobabContract(append(nodeAddrs, testAddrs...), new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil))
	allocationFunction(testGenesis)
	WriteFile([]byte(baobabOperatorAddress), "baobab_operator", "address")
	WriteFile([]byte(baobabOperatorKey), "baobab_operator", "private")
	return testGenesis
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+{}|[]")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func gen(ctx *cli.Context) error {
	genType := findGenType(ctx)

	cliqueFlag := ctx.Bool(cliqueFlag.Name)
	cnNum := ctx.Int(numOfCNsFlag.Name)
	numValidators := ctx.Int(numOfValidatorsFlag.Name)
	pnNum := ctx.Int(numOfPNsFlag.Name)
	enNum := ctx.Int(numOfENsFlag.Name)
	scnNum := ctx.Int(numOfSCNsFlag.Name)
	spnNum := ctx.Int(numOfSPNsFlag.Name)
	senNum := ctx.Int(numOfSENsFlag.Name)
	numTestAccs := ctx.Int(numOfTestKeyFlag.Name)
	baobab := ctx.Bool(baobabFlag.Name)
	baobabTest := ctx.Bool(baobabTestFlag.Name)
	cypress := ctx.Bool(cypressFlag.Name)
	cypressTest := ctx.Bool(cypressTestFlag.Name)
	serviceChain := ctx.Bool(serviceChainFlag.Name)
	serviceChainTest := ctx.Bool(serviceChainTestFlag.Name)
	chainid := ctx.Uint64(chainIDFlag.Name)
	serviceChainId := ctx.Uint64(serviceChainIDFlag.Name)

	if cnNum == 0 && scnNum == 0 {
		return fmt.Errorf("needed at least one consensus node (--cn-num 1) or one service chain consensus node (--scn-num 1) ")
	}

	if numValidators == 0 {
		numValidators = cnNum
	}
	if numValidators > cnNum {
		return fmt.Errorf("validators-num(%d) cannot be greater than num(%d)", numValidators, cnNum)
	}

	privKeys, nodeKeys, nodeAddrs := istcommon.GenerateKeys(cnNum)
	testPrivKeys, testKeys, testAddrs := istcommon.GenerateKeys(numTestAccs)

	var (
		genesisJson      *blockchain.Genesis
		genesisJsonBytes []byte
	)

	validatorNodeAddrs := make([]common.Address, numValidators)
	copy(validatorNodeAddrs, nodeAddrs[:numValidators])

	if cypressTest {
		genesisJson = genCypressTestGenesis(validatorNodeAddrs, testAddrs)
	} else if cypress {
		genesisJson = genCypressGenesis(validatorNodeAddrs, testAddrs)
	} else if baobabTest {
		genesisJson = genBaobabTestGenesis(validatorNodeAddrs, testAddrs)
	} else if baobab {
		genesisJson = genBaobabGenesis(validatorNodeAddrs, testAddrs)
	} else if cliqueFlag {
		genesisJson = genCliqueGenesis(ctx, validatorNodeAddrs, testAddrs, chainid)
	} else if serviceChain {
		genesisJson = genServiceChainGenesis(validatorNodeAddrs, testAddrs)
	} else if serviceChainTest {
		genesisJson = genServiceChainTestGenesis(validatorNodeAddrs, testAddrs)
	} else {
		genesisJson = genIstanbulGenesis(ctx, validatorNodeAddrs, testAddrs, chainid)
	}

	genesisJson.Config.IstanbulCompatibleBlock = big.NewInt(ctx.Int64(istanbulCompatibleBlockNumberFlag.Name))
	genesisJson.Config.LondonCompatibleBlock = big.NewInt(ctx.Int64(londonCompatibleBlockNumberFlag.Name))
	genesisJson.Config.EthTxTypeCompatibleBlock = big.NewInt(ctx.Int64(ethTxTypeCompatibleBlockNumberFlag.Name))

	genesisJsonBytes, _ = json.MarshalIndent(genesisJson, "", "    ")
	genValidatorKeystore(privKeys)
	lastIssuedPortNum = uint16(ctx.Int(p2pPortFlag.Name))

	switch genType {
	case TypeDocker:
		validators := makeValidators(cnNum, false, nodeAddrs, nodeKeys, privKeys)
		pnValidators, proxyNodeKeys := makeProxys(pnNum, false)
		nodeInfos := filterNodeInfo(validators)
		staticNodesJsonBytes, _ := json.MarshalIndent(nodeInfos, "", "\t")
		address := filterAddressesString(validators)
		pnInfos := filterNodeInfo(pnValidators)
		enValidators, enKeys := makeEndpoints(enNum, false)
		enInfos := filterNodeInfo(enValidators)

		scnValidators, scnKeys := makeSCNs(scnNum, false)
		scnInfos := filterNodeInfo(scnValidators)
		scnAddress := filterAddresses(scnValidators)

		spnValidators, spnKeys := makeSPNs(spnNum, false)
		spnInfos := filterNodeInfo(spnValidators)

		senValidators, senKeys := makeSENs(senNum, false)
		senInfos := filterNodeInfo(senValidators)

		staticPNJsonBytes, _ := json.MarshalIndent(pnInfos, "", "\t")
		staticENJsonBytes, _ := json.MarshalIndent(enInfos, "", "\t")
		staticSCNJsonBytes, _ := json.MarshalIndent(scnInfos, "", "\t")
		staticSPNJsonBytes, _ := json.MarshalIndent(spnInfos, "", "\t")
		staticSENJsonBytes, _ := json.MarshalIndent(senInfos, "", "\t")
		var bridgeNodesJsonBytes []byte
		if len(enInfos) != 0 {
			bridgeNodesJsonBytes, _ = json.MarshalIndent(enInfos[:1], "", "\t")
		}
		scnGenesisJsonBytes, _ := json.MarshalIndent(genIstanbulGenesis(ctx, scnAddress, nil, serviceChainId), "", "\t")

		dockerImageId := ctx.String(dockerImageIdFlag.Name)

		compose := compose.New(
			"172.16.239",
			cnNum,
			"bb98a0b6442386d0cdf8a31b267892c1",
			address,
			nodeKeys,
			removeSpacesAndLines(genesisJsonBytes),
			removeSpacesAndLines(scnGenesisJsonBytes),
			removeSpacesAndLines(staticNodesJsonBytes),
			removeSpacesAndLines(staticPNJsonBytes),
			removeSpacesAndLines(staticENJsonBytes),
			removeSpacesAndLines(staticSCNJsonBytes),
			removeSpacesAndLines(staticSPNJsonBytes),
			removeSpacesAndLines(staticSENJsonBytes),
			removeSpacesAndLines(bridgeNodesJsonBytes),
			dockerImageId,
			ctx.Bool(fasthttpFlag.Name),
			ctx.Int(networkIdFlag.Name),
			int(chainid),
			!ctx.BoolT(nografanaFlag.Name),
			proxyNodeKeys,
			enKeys,
			scnKeys,
			spnKeys,
			senKeys,
			ctx.Bool(useTxGenFlag.Name),
			service.TxGenOption{
				TxGenRate:       ctx.Int(txGenRateFlag.Name),
				TxGenThreadSize: ctx.Int(txGenThFlag.Name),
				TxGenConnSize:   ctx.Int(txGenConnFlag.Name),
				TxGenDuration:   ctx.String(txGenDurFlag.Name),
			})
		os.MkdirAll(outputPath, os.ModePerm)
		ioutil.WriteFile(path.Join(outputPath, "docker-compose.yml"), []byte(compose.String()), os.ModePerm)
		fmt.Println("Created : ", path.Join(outputPath, "docker-compose.yml"))
		ioutil.WriteFile(path.Join(outputPath, "prometheus.yml"), []byte(compose.PrometheusService.Config.String()), os.ModePerm)
		fmt.Println("Created : ", path.Join(outputPath, "prometheus.yml"))
		downLoadGrafanaJson()
	case TypeLocal:
		writeNodeFiles(true, cnNum, pnNum, nodeAddrs, nodeKeys, privKeys, genesisJsonBytes)
		writeTestKeys(DirTestKeys, testPrivKeys, testKeys)
		downLoadGrafanaJson()
	case TypeRemote:
		writeNodeFiles(false, cnNum, pnNum, nodeAddrs, nodeKeys, privKeys, genesisJsonBytes)
		writeTestKeys(DirTestKeys, testPrivKeys, testKeys)
		downLoadGrafanaJson()
	case TypeDeploy:
		writeCNInfoKey(cnNum, nodeAddrs, nodeKeys, privKeys, genesisJsonBytes)
		writeKlayConfig(ctx.Int(networkIdFlag.Name), ctx.Int(rpcPortFlag.Name), ctx.Int(wsPortFlag.Name), ctx.Int(p2pPortFlag.Name),
			ctx.String(dataDirFlag.Name), ctx.String(logDirFlag.Name), "CN")
		writeKlayConfig(ctx.Int(networkIdFlag.Name), ctx.Int(rpcPortFlag.Name), ctx.Int(wsPortFlag.Name), ctx.Int(p2pPortFlag.Name),
			ctx.String(dataDirFlag.Name), ctx.String(logDirFlag.Name), "PN")
		writePNInfoKey(ctx.Int(numOfPNsFlag.Name))
		writePrometheusConfig(cnNum, ctx.Int(numOfPNsFlag.Name))
	}

	return nil
}

func downLoadGrafanaJson() {
	for _, file := range GrafanaFiles {
		resp, err := http.Get(file.url)
		if err != nil {
			fmt.Printf("Failed to download the imgs dashboard file(%s) - %v\n", file.url, err)
		} else if resp.StatusCode != 200 {
			fmt.Printf("Failed to download the imgs dashboard file(%s) [%s] - %v\n", file.url, resp.Status, err)
		} else {
			bytes, e := ioutil.ReadAll(resp.Body)
			if e != nil {
				fmt.Println("Failed to read http response", e)
			} else {
				fileName := file.name
				ioutil.WriteFile(path.Join(outputPath, fileName), bytes, os.ModePerm)
				fmt.Println("Created : ", path.Join(outputPath, fileName))
			}
			resp.Body.Close()
		}
	}
}

func writeCNInfoKey(num int, nodeAddrs []common.Address, nodeKeys []string, privKeys []*ecdsa.PrivateKey,
	genesisJsonBytes []byte) {
	const DirCommon = "common"
	WriteFile(genesisJsonBytes, DirCommon, "genesis.json")

	validators := makeValidatorsWithIp(num, false, nodeAddrs, nodeKeys, privKeys, []string{CNIpNetwork})
	staticNodesJsonBytes, _ := json.MarshalIndent(filterNodeInfo(validators), "", "\t")
	WriteFile(staticNodesJsonBytes, DirCommon, "static-nodes.json")

	for i, v := range validators {
		parentDir := fmt.Sprintf("cn%02d", i+1)
		WriteFile([]byte(nodeKeys[i]), parentDir, "nodekey")
		str, _ := json.MarshalIndent(v, "", "\t")
		WriteFile([]byte(str), parentDir, "validator")
	}
}

func writePNInfoKey(num int) {
	privKeys, nodeKeys, nodeAddrs := istcommon.GenerateKeys(num)
	validators := makeValidatorsWithIp(num, false, nodeAddrs, nodeKeys, privKeys, []string{PNIpNetwork1, PNIpNetwork2})
	for i, v := range validators {
		parentDir := fmt.Sprintf("pn%02d", i+1)
		WriteFile([]byte(nodeKeys[i]), parentDir, "nodekey")
		str, _ := json.MarshalIndent(v, "", "\t")
		WriteFile([]byte(str), parentDir, "validator")
	}
}

func writeKlayConfig(networkId int, rpcPort int, wsPort int, p2pPort int, dataDir string, logDir string, nodeType string) {
	kConfig := KlaytnConfig{
		networkId,
		rpcPort,
		wsPort,
		p2pPort,
		dataDir,
		logDir,
		"/var/run/klay",
		nodeType,
	}
	WriteFile([]byte(kConfig.String()), strings.ToLower(nodeType), "klay.conf")
}

func writePrometheusConfig(cnNum int, pnNum int) {
	pConf := NewPrometheusConfig(cnNum, CNIpNetwork, pnNum, PNIpNetwork1, PNIpNetwork2)
	WriteFile([]byte(pConf.String()), "monitoring", "prometheus.yml")
}

func writeNodeFiles(isWorkOnSingleHost bool, num int, pnum int, nodeAddrs []common.Address, nodeKeys []string,
	privKeys []*ecdsa.PrivateKey, genesisJsonBytes []byte) {
	WriteFile(genesisJsonBytes, DirScript, "genesis.json")

	validators := makeValidators(num, isWorkOnSingleHost, nodeAddrs, nodeKeys, privKeys)
	nodeInfos := filterNodeInfo(validators)
	staticNodesJsonBytes, _ := json.MarshalIndent(nodeInfos, "", "\t")
	writeValidatorsAndNodesToFile(validators, DirKeys, nodeKeys)
	WriteFile(staticNodesJsonBytes, DirScript, "static-nodes.json")

	if pnum > 0 {
		proxys, proxyNodeKeys := makeProxys(pnum, isWorkOnSingleHost)
		pNodeInfos := filterNodeInfo(proxys)
		staticPNodesJsonBytes, _ := json.MarshalIndent(pNodeInfos, "", "\t")
		writeValidatorsAndNodesToFile(proxys, DirPnKeys, proxyNodeKeys)
		WriteFile(staticPNodesJsonBytes, DirPnScript, "static-nodes.json")
	}
}

func filterAddresses(validatorInfos []*ValidatorInfo) []common.Address {
	var addresses []common.Address
	for _, v := range validatorInfos {
		addresses = append(addresses, v.Address)
	}
	return addresses
}

func filterAddressesString(validatorInfos []*ValidatorInfo) []string {
	var address []string
	for _, v := range validatorInfos {
		address = append(address, v.Address.String())
	}
	return address
}

func filterNodeInfo(validatorInfos []*ValidatorInfo) []string {
	var nodes []string
	for _, v := range validatorInfos {
		nodes = append(nodes, string(v.NodeInfo))
	}
	return nodes
}

func makeValidators(num int, isWorkOnSingleHost bool, nodeAddrs []common.Address, nodeKeys []string,
	keys []*ecdsa.PrivateKey) []*ValidatorInfo {
	var validatorPort uint16
	var validators []*ValidatorInfo
	for i := 0; i < num; i++ {
		if isWorkOnSingleHost {
			validatorPort = lastIssuedPortNum
			lastIssuedPortNum++
		} else {
			validatorPort = DefaultTcpPort
		}

		v := &ValidatorInfo{
			Address: nodeAddrs[i],
			Nodekey: nodeKeys[i],
			NodeInfo: discover.NewNode(
				discover.PubkeyID(&keys[i].PublicKey),
				net.ParseIP("0.0.0.0"),
				0,
				validatorPort,
				nil,
				discover.NodeTypeCN).String(),
		}
		validators = append(validators, v)
	}
	return validators
}

func makeValidatorsWithIp(num int, isWorkOnSingleHost bool, nodeAddrs []common.Address, nodeKeys []string,
	keys []*ecdsa.PrivateKey, networkIds []string) []*ValidatorInfo {
	var validatorPort uint16
	var validators []*ValidatorInfo
	for i := 0; i < num; i++ {
		if isWorkOnSingleHost {
			validatorPort = lastIssuedPortNum
			lastIssuedPortNum++
		} else {
			validatorPort = DefaultTcpPort
		}

		nn := len(networkIds)
		idx := (i + 1) % nn
		if nn > 1 {
			if idx == 0 {
				idx = nn - 1
			} else { // idx > 0
				idx = idx - 1
			}
		}
		v := &ValidatorInfo{
			Address: nodeAddrs[i],
			Nodekey: nodeKeys[i],
			NodeInfo: discover.NewNode(
				discover.PubkeyID(&keys[i].PublicKey),
				net.ParseIP(fmt.Sprintf("%s.%d", networkIds[idx], 100+(i/nn)+1)),
				0,
				validatorPort,
				nil,
				discover.NodeTypeCN).String(),
		}
		validators = append(validators, v)
	}
	return validators
}

func makeProxys(num int, isWorkOnSingleHost bool) ([]*ValidatorInfo, []string) {
	privKeys, nodeKeys, nodeAddrs := istcommon.GenerateKeys(num)

	var p2pPort uint16
	var proxies []*ValidatorInfo
	var proxyNodeKeys []string
	for i := 0; i < num; i++ {
		if isWorkOnSingleHost {
			p2pPort = lastIssuedPortNum
			lastIssuedPortNum++
		} else {
			p2pPort = DefaultTcpPort
		}

		v := &ValidatorInfo{
			Address: nodeAddrs[i],
			Nodekey: nodeKeys[i],
			NodeInfo: discover.NewNode(
				discover.PubkeyID(&privKeys[i].PublicKey),
				net.ParseIP("0.0.0.0"),
				0,
				p2pPort,
				nil,
				discover.NodeTypePN).String(),
		}
		proxies = append(proxies, v)
		proxyNodeKeys = append(proxyNodeKeys, v.Nodekey)
	}
	return proxies, proxyNodeKeys
}

func makeEndpoints(num int, isWorkOnSingleHost bool) ([]*ValidatorInfo, []string) {
	privKeys, nodeKeys, nodeAddrs := istcommon.GenerateKeys(num)

	var p2pPort uint16
	var endpoints []*ValidatorInfo
	var endpointsNodeKeys []string
	for i := 0; i < num; i++ {
		if isWorkOnSingleHost {
			p2pPort = lastIssuedPortNum
			lastIssuedPortNum++
		} else {
			p2pPort = DefaultTcpPort
		}

		v := &ValidatorInfo{
			Address: nodeAddrs[i],
			Nodekey: nodeKeys[i],
			NodeInfo: discover.NewNode(
				discover.PubkeyID(&privKeys[i].PublicKey),
				net.ParseIP("0.0.0.0"),
				0,
				p2pPort,
				nil,
				discover.NodeTypeEN).String(),
		}
		endpoints = append(endpoints, v)
		endpointsNodeKeys = append(endpointsNodeKeys, v.Nodekey)
	}
	return endpoints, endpointsNodeKeys
}

func makeSCNs(num int, isWorkOnSingleHost bool) ([]*ValidatorInfo, []string) {
	privKeys, nodeKeys, nodeAddrs := istcommon.GenerateKeys(num)

	var p2pPort uint16
	var scn []*ValidatorInfo
	var scnKeys []string
	for i := 0; i < num; i++ {
		if isWorkOnSingleHost {
			p2pPort = lastIssuedPortNum
			lastIssuedPortNum++
		} else {
			p2pPort = DefaultTcpPort
		}

		v := &ValidatorInfo{
			Address: nodeAddrs[i],
			Nodekey: nodeKeys[i],
			NodeInfo: discover.NewNode(
				discover.PubkeyID(&privKeys[i].PublicKey),
				net.ParseIP("0.0.0.0"),
				0,
				p2pPort,
				nil,
				discover.NodeTypeUnknown).String(),
		}
		scn = append(scn, v)
		scnKeys = append(scnKeys, v.Nodekey)
	}
	return scn, scnKeys
}

func makeSPNs(num int, isWorkOnSingleHost bool) ([]*ValidatorInfo, []string) {
	privKeys, nodeKeys, nodeAddrs := istcommon.GenerateKeys(num)

	var p2pPort uint16
	var proxies []*ValidatorInfo
	var proxyNodeKeys []string
	for i := 0; i < num; i++ {
		if isWorkOnSingleHost {
			p2pPort = lastIssuedPortNum
			lastIssuedPortNum++
		} else {
			p2pPort = DefaultTcpPort
		}

		v := &ValidatorInfo{
			Address: nodeAddrs[i],
			Nodekey: nodeKeys[i],
			NodeInfo: discover.NewNode(
				discover.PubkeyID(&privKeys[i].PublicKey),
				net.ParseIP("0.0.0.0"),
				0,
				p2pPort,
				nil,
				discover.NodeTypeUnknown).String(),
		}
		proxies = append(proxies, v)
		proxyNodeKeys = append(proxyNodeKeys, v.Nodekey)
	}
	return proxies, proxyNodeKeys
}

func makeSENs(num int, isWorkOnSingleHost bool) ([]*ValidatorInfo, []string) {
	privKeys, nodeKeys, nodeAddrs := istcommon.GenerateKeys(num)

	var p2pPort uint16
	var endpoints []*ValidatorInfo
	var endpointsNodeKeys []string
	for i := 0; i < num; i++ {
		if isWorkOnSingleHost {
			p2pPort = lastIssuedPortNum
			lastIssuedPortNum++
		} else {
			p2pPort = DefaultTcpPort
		}

		v := &ValidatorInfo{
			Address: nodeAddrs[i],
			Nodekey: nodeKeys[i],
			NodeInfo: discover.NewNode(
				discover.PubkeyID(&privKeys[i].PublicKey),
				net.ParseIP("0.0.0.0"),
				0,
				p2pPort,
				nil,
				discover.NodeTypeUnknown).String(),
		}
		endpoints = append(endpoints, v)
		endpointsNodeKeys = append(endpointsNodeKeys, v.Nodekey)
	}
	return endpoints, endpointsNodeKeys
}

func writeValidatorsAndNodesToFile(validators []*ValidatorInfo, parentDir string, nodekeys []string) {
	parentPath := path.Join(outputPath, parentDir)
	os.MkdirAll(parentPath, os.ModePerm)

	for i, v := range validators {
		nodeKeyFilePath := path.Join(parentPath, "nodekey"+strconv.Itoa(i+1))
		ioutil.WriteFile(nodeKeyFilePath, []byte(nodekeys[i]), os.ModePerm)
		fmt.Println("Created : ", nodeKeyFilePath)

		str, _ := json.MarshalIndent(v, "", "\t")
		validatorInfoFilePath := path.Join(parentPath, "validator"+strconv.Itoa(i+1))
		ioutil.WriteFile(validatorInfoFilePath, []byte(str), os.ModePerm)
		fmt.Println("Created : ", validatorInfoFilePath)
	}
}

func writeTestKeys(parentDir string, privKeys []*ecdsa.PrivateKey, keys []string) {
	parentPath := path.Join(outputPath, parentDir)
	os.MkdirAll(parentPath, os.ModePerm)

	for i, key := range keys {
		testKeyFilePath := path.Join(parentPath, "testkey"+strconv.Itoa(i+1))
		ioutil.WriteFile(testKeyFilePath, []byte(key), os.ModePerm)
		fmt.Println("Created : ", testKeyFilePath)

		pk := privKeys[i]
		ksPath := path.Join(parentPath, "keystore"+strconv.Itoa(i+1))
		ks := keystore.NewKeyStore(ksPath, keystore.StandardScryptN, keystore.StandardScryptP)
		pwdStr := RandStringRunes(params.PasswordLength)
		ks.ImportECDSA(pk, pwdStr)
		WriteFile([]byte(pwdStr), path.Join(parentDir, "keystore"+strconv.Itoa(i+1)), crypto.PubkeyToAddress(pk.PublicKey).String())
	}
}

func WriteFile(content []byte, parentFolder string, fileName string) {
	filePath := path.Join(outputPath, parentFolder, fileName)
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	ioutil.WriteFile(filePath, content, os.ModePerm)
	fmt.Println("Created : ", filePath)
}

func findGenType(ctx *cli.Context) int {
	var genType = TypeNotDefined
	if len(ctx.Args()) >= 1 {
		for i, t := range Types {
			if t == ctx.Args()[0] {
				genType = i
				break
			}
		}
		if genType == TypeNotDefined {
			fmt.Printf("Wrong Type : %s\nSupported Types : [docker, local, remote, deploy]\n\n", ctx.Args()[0])
			cli.ShowSubcommandHelp(ctx)
			os.Exit(1)
		}
	} else {
		genType = TypeDocker
	}
	return genType
}

func removeSpacesAndLines(b []byte) string {
	out := string(b)
	out = strings.Replace(out, " ", "", -1)
	out = strings.Replace(out, "\t", "", -1)
	out = strings.Replace(out, "\n", "", -1)
	return out
}
