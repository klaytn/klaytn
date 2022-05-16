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
	"github.com/klaytn/klaytn/params"
	"gopkg.in/urfave/cli.v1"
)

var fundingAddr string
var dockerImageId string
var outputPath string

var (
	cypressTestFlag = cli.BoolFlag{
		Name:  "cypress-test",
		Usage: "Generate genesis.json similar to the one used for Cypress with shorter intervals for testing",
	}

	cypressFlag = cli.BoolFlag{
		Name:  "cypress",
		Usage: "Generate genesis.json similar to the one used for Cypress",
	}

	baobabTestFlag = cli.BoolFlag{
		Name:  "baobab-test",
		Usage: "Generate genesis.json similar to the one used for Baobab with shorter intervals for testing",
	}

	baobabFlag = cli.BoolFlag{
		Name:  "baobab",
		Usage: "Generate genesis.json similar to the one used for Baobab",
	}

	serviceChainFlag = cli.BoolFlag{
		Name:  "servicechain",
		Usage: "Generate genesis.json similar to the one used for Serivce Chain",
	}

	serviceChainTestFlag = cli.BoolFlag{
		Name:  "servicechain-test",
		Usage: "Generate genesis.json similar to the one used for Serivce Chain with shorter intervals for testing",
	}

	cliqueFlag = cli.BoolFlag{
		Name:  "clique",
		Usage: "Use Clique consensus",
	}

	numOfCNsFlag = cli.IntFlag{
		Name:  "cn-num",
		Usage: "Number of consensus nodes",
		Value: 0,
	}

	numOfValidatorsFlag = cli.IntFlag{
		Name:  "validators-num",
		Usage: "Number of validators. If not set, it will be set the number of option cn-num",
		Value: 0,
	}

	numOfPNsFlag = cli.IntFlag{
		Name:  "pn-num",
		Usage: "Number of proxy node",
		Value: 0,
	}

	numOfENsFlag = cli.IntFlag{
		Name:  "en-num",
		Usage: "Number of end-point node",
		Value: 0,
	}

	numOfSCNsFlag = cli.IntFlag{
		Name:  "scn-num",
		Usage: "Number of service chain nodes",
		Value: 0,
	}

	numOfSPNsFlag = cli.IntFlag{
		Name:  "spn-num",
		Usage: "Number of service chain proxy nodes",
		Value: 0,
	}

	numOfSENsFlag = cli.IntFlag{
		Name:  "sen-num",
		Usage: "Number of service chain end-point nodes",
		Value: 0,
	}

	numOfTestKeyFlag = cli.IntFlag{
		Name:  "test-num",
		Usage: "Number of test key",
		Value: 0,
	}

	chainIDFlag = cli.Uint64Flag{
		Name:  "chainID",
		Usage: "ChainID",
		Value: 1000,
	}

	serviceChainIDFlag = cli.Uint64Flag{
		Name:  "serviceChainID",
		Usage: "Service Chain ID",
		Value: 1001,
	}

	unitPriceFlag = cli.Uint64Flag{
		Name:  "unitPrice",
		Usage: "Price of unit",
		Value: 0,
	}

	deriveShaImplFlag = cli.IntFlag{
		Name:  "deriveShaImpl",
		Usage: "Implementation of DeriveSha [0:Original, 1:Simple, 2:Concat]",
		Value: 0,
	}

	outputPathFlag = cli.StringFlag{
		Name:        "output, o",
		Usage:       "homi's result saved at this output folder",
		Value:       "homi-output",
		Destination: &outputPath,
	}

	fundingAddrFlag = cli.StringFlag{
		Name:        "fundingAddr",
		Value:       "75a59b94889a05c03c66c3c84e9d2f8308ca4abd",
		Usage:       "Give initial fund to the given addr",
		Destination: &fundingAddr,
	}

	dockerImageIdFlag = cli.StringFlag{
		Name:        "docker-image-id",
		Value:       "klaytn/klaytn:latest", // https://hub.docker.com/r/klaytn/klaytn
		Usage:       "Base docker image ID (Image[:tag]), e.g., klaytn/klaytn:v1.5.3",
		Destination: &dockerImageId,
	}

	subGroupSizeFlag = cli.IntFlag{
		Name:  "subgroup-size",
		Usage: "CN's Subgroup size",
		Value: 21,
	}

	fasthttpFlag = cli.BoolFlag{
		Name:  "fasthttp",
		Usage: "(docker only) Use High performance http module",
	}

	networkIdFlag = cli.IntFlag{
		Name:  "network-id",
		Usage: "(docker only) network identifier (default : 2018)",
		Value: 2018,
	}

	nografanaFlag = cli.BoolFlag{
		Name:  "no-grafana",
		Usage: "(docker only) Do not make grafana container",
	}

	useTxGenFlag = cli.BoolFlag{
		Name:  "txgen",
		Usage: "(docker only) Add txgen container",
	}

	txGenRateFlag = cli.IntFlag{
		Name:  "txgen-rate",
		Usage: "(docker only) txgen's rate option [default : 2000]",
		Value: 2000,
	}

	txGenConnFlag = cli.IntFlag{
		Name:  "txgen-conn",
		Usage: "(docker only) txgen's connection size option [default : 100]",
		Value: 100,
	}

	txGenDurFlag = cli.StringFlag{
		Name:  "txgen-dur",
		Usage: "(docker only) txgen's duration option [default : 1m]",
		Value: "1m",
	}

	txGenThFlag = cli.IntFlag{
		Name:  "txgen-th",
		Usage: "(docker-only) txgen's thread size option [default : 2]",
		Value: 2,
	}

	rpcPortFlag = cli.IntFlag{
		Name:  "rpc-port",
		Usage: "klay.conf - Klaytn node's rpc port [default: 8551] ",
		Value: 8551,
	}

	wsPortFlag = cli.IntFlag{
		Name:  "ws-port",
		Usage: "klay.conf - Klaytn node's ws port [default: 8552]",
		Value: 8552,
	}

	p2pPortFlag = cli.IntFlag{
		Name:  "p2p-port",
		Usage: "klay.conf - Klaytn node's p2p port [default: 32323]",
		Value: 32323,
	}

	dataDirFlag = cli.StringFlag{
		Name:  "data-dir",
		Usage: "klay.conf - Klaytn node's data directory path [default : /var/klay/data]",
		Value: "/var/klay/data",
	}

	logDirFlag = cli.StringFlag{
		Name:  "log-dir",
		Usage: "klay.conf - Klaytn node's log directory path [default : /var/klay/log]",
		Value: "/var/klay/log",
	}

	// Governance flags
	governanceFlag = cli.BoolFlag{
		Name:  "governance",
		Usage: "governance field is added in the genesis file if this flag is set",
	}

	govModeFlag = cli.StringFlag{
		Name:  "gov-mode",
		Usage: "governance mode (none, single, ballot) [default: none]",
		Value: params.DefaultGovernanceMode,
	}

	governingNodeFlag = cli.StringFlag{
		Name:  "governing-node",
		Usage: "the governing node [default: 0x0000000000000000000000000000000000000000]",
		Value: params.DefaultGoverningNode,
	}

	rewardMintAmountFlag = cli.StringFlag{
		Name:  "reward-mint-amount",
		Usage: "governance minting amount",
		Value: "9600000000000000000",
	}

	rewardRatioFlag = cli.StringFlag{
		Name:  "reward-ratio",
		Usage: "governance ratio [default: 100/0/0]",
		Value: params.DefaultRatio,
	}

	rewardGiniCoeffFlag = cli.BoolFlag{
		Name:  "reward-gini-coeff",
		Usage: "governance gini-coefficient",
	}

	rewardDeferredTxFeeFlag = cli.BoolFlag{
		Name:  "reward-deferred-tx",
		Usage: "governance deferred transaction",
	}

	rewardStakingFlag = cli.Uint64Flag{
		Name:  "reward-staking-interval",
		Usage: "reward staking update interval flag",
		Value: 86400,
	}

	rewardProposerFlag = cli.Uint64Flag{
		Name:  "reward-proposer-interval",
		Usage: "reward proposer update interval flag",
		Value: 3600,
	}

	rewardMinimumStakeFlag = cli.StringFlag{
		Name:  "reward-minimum-stake",
		Usage: "reward minimum stake flag",
		Value: "2000000",
	}

	istEpochFlag = cli.Uint64Flag{
		Name:  "ist-epoch",
		Usage: "governance epoch [default: 604800]",
		Value: params.DefaultEpoch,
	}

	istProposerPolicyFlag = cli.Uint64Flag{
		Name:  "ist-proposer-policy",
		Usage: "governance proposer policy (0: RoundRobin, 1: Sticky, 2: WeightedRandom) [default: 0]",
		Value: params.DefaultProposerPolicy,
	}

	istSubGroupFlag = cli.Uint64Flag{
		Name:  "ist-subgroup",
		Usage: "governance subgroup size [default: 21]",
		Value: params.DefaultSubGroupSize,
	}

	cliqueEpochFlag = cli.Uint64Flag{
		Name:  "clique-epoch",
		Usage: "clique epoch",
		Value: params.DefaultEpoch,
	}

	cliquePeriodFlag = cli.Uint64Flag{
		Name:  "clique-period",
		Usage: "clique period",
		Value: params.DefaultPeriod,
	}

	istanbulCompatibleBlockNumberFlag = cli.Int64Flag{
		Name:  "istanbul-compatible-blocknumber",
		Usage: "istanbulCompatible blockNumber",
		Value: 0,
	}

	londonCompatibleBlockNumberFlag = cli.Int64Flag{
		Name:  "london-compatible-blocknumber",
		Usage: "londonCompatible blockNumber",
		Value: 0,
	}

	ethTxTypeCompatibleBlockNumberFlag = cli.Int64Flag{
		Name:  "eth-tx-type-compatible-blocknumber",
		Usage: "ethTxTypeCompatible blockNumber",
		Value: 0,
	}
)
