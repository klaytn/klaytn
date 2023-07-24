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
	"github.com/urfave/cli/v2"
)

var (
	dockerImageId string
	outputPath    string
)

var (
	homiYamlFlag = cli.StringFlag{
		Name:    "homi-yaml",
		Usage:   "Import homi.yaml to generate the config files to run the nodes",
		Aliases: []string{"yaml"},
	}
	genTypeFlag = cli.StringFlag{
		Name:    "gen-type",
		Usage:   "Generate environment files according to the type (docker, local, remote, deploy)",
		Aliases: []string{},
		Value:   "docker",
	}

	cypressTestFlag = cli.BoolFlag{
		Name:    "cypress-test",
		Usage:   "Generate genesis.json similar to the one used for Cypress with shorter intervals for testing",
		Aliases: []string{},
	}
	cypressFlag = cli.BoolFlag{
		Name:    "cypress",
		Usage:   "Generate genesis.json similar to the one used for Cypress",
		Aliases: []string{},
	}
	baobabTestFlag = cli.BoolFlag{
		Name:    "baobab-test",
		Usage:   "Generate genesis.json similar to the one used for Baobab with shorter intervals for testing",
		Aliases: []string{},
	}
	baobabFlag = cli.BoolFlag{
		Name:    "baobab",
		Usage:   "Generate genesis.json similar to the one used for Baobab",
		Aliases: []string{},
	}
	serviceChainFlag = cli.BoolFlag{
		Name:    "servicechain",
		Usage:   "Generate genesis.json similar to the one used for Serivce Chain",
		Aliases: []string{},
	}
	serviceChainTestFlag = cli.BoolFlag{
		Name:    "servicechain-test",
		Usage:   "Generate genesis.json similar to the one used for Serivce Chain with shorter intervals for testing",
		Aliases: []string{},
	}
	cliqueFlag = cli.BoolFlag{
		Name:    "clique",
		Usage:   "Use Clique consensus",
		Aliases: []string{},
	}

	numOfCNsFlag = cli.IntFlag{
		Name:    "cn-num",
		Usage:   "Number of consensus nodes",
		Aliases: []string{"topology.cn-num"},
		Value:   0,
	}
	numOfValidatorsFlag = cli.IntFlag{
		Name:    "validators-num",
		Usage:   "Number of validators. If not set, it will be set the number of option cn-num",
		Aliases: []string{"topology.validators-num"},
		Value:   0,
	}
	numOfPNsFlag = cli.IntFlag{
		Name:    "pn-num",
		Usage:   "Number of proxy node",
		Aliases: []string{"topology.pn-num"},
		Value:   0,
	}
	numOfENsFlag = cli.IntFlag{
		Name:    "en-num",
		Usage:   "Number of end-point node",
		Aliases: []string{"topology.en-num"},
		Value:   0,
	}
	numOfSCNsFlag = cli.IntFlag{
		Name:    "scn-num",
		Usage:   "Number of service chain nodes",
		Aliases: []string{"topology.scn-num"},
		Value:   0,
	}
	numOfSPNsFlag = cli.IntFlag{
		Name:    "spn-num",
		Usage:   "Number of service chain proxy nodes",
		Aliases: []string{"topology.spn-num"},
		Value:   0,
	}
	numOfSENsFlag = cli.IntFlag{
		Name:    "sen-num",
		Usage:   "Number of service chain end-point nodes",
		Aliases: []string{"topology.sen-num"},
		Value:   0,
	}
	numOfTestKeyFlag = cli.IntFlag{
		Name:    "test-num",
		Usage:   "Number of test key",
		Aliases: []string{"topology.test-num"},
		Value:   0,
	}

	genesisTypeFlag = cli.StringFlag{
		Name:    "genesis-type",
		Usage:   "Set the type of genesis.json to generate (cypress-test, cypress, baobab-test, baobab, clique, servicechain, servicechain-test, istanbul)",
		Aliases: []string{"genesis.type"},
	}
	chainIDFlag = cli.Uint64Flag{
		Name:    "chainID",
		Usage:   "ChainID",
		Aliases: []string{"genesis.chain-id"},
		Value:   1000,
	}

	serviceChainIDFlag = cli.Uint64Flag{
		Name:    "serviceChainID",
		Usage:   "Service Chain ID",
		Aliases: []string{"genesis.service-chain-id"},
		Value:   1001,
	}

	unitPriceFlag = cli.Uint64Flag{
		Name:    "unitPrice",
		Usage:   "Price of unit",
		Aliases: []string{"genesis.unit-price"},
		Value:   0,
	}

	deriveShaImplFlag = cli.IntFlag{
		Name:    "deriveShaImpl",
		Usage:   "Implementation of DeriveSha [0:Original, 1:Simple, 2:Concat]",
		Aliases: []string{"genesis.derive-sha-impl"},
		Value:   0,
	}

	outputPathFlag = cli.StringFlag{
		Name:        "output",
		Usage:       "homi's result saved at this output folder",
		Aliases:     []string{"o"},
		Value:       "homi-output",
		Destination: &outputPath,
	}

	fundingAddrFlag = cli.StringFlag{
		Name:    "funding-addr",
		Value:   "",
		Usage:   "Give initial fund to the given addr",
		Aliases: []string{"genesis.funding-addr"},
	}

	patchAddressBookFlag = cli.BoolFlag{
		Name:    "patch-address-book",
		Usage:   "Patch genesis AddressBook's constructContract function",
		Aliases: []string{"genesis.patched-address-book.enable"},
	}

	patchAddressBookAddrFlag = cli.StringFlag{
		Name:    "patch-address-book-addr",
		Usage:   "The address to inject in AddressBook's constructContract function [default: first CN's address]",
		Aliases: []string{"genesis.patched-address-book.addr"},
		Value:   "",
	}
	// Deprecated: it doesn't used anymore
	subGroupSizeFlag = cli.IntFlag{
		Name:    "subgroup-size",
		Usage:   "CN's Subgroup size",
		Aliases: []string{},
		Value:   21,
	}

	dockerImageIdFlag = cli.StringFlag{
		Name:        "docker-image-id",
		Value:       "klaytn/klaytn:latest", // https://hub.docker.com/r/klaytn/klaytn
		Usage:       "Base docker image ID (Image[:tag]), e.g., klaytn/klaytn:v1.5.3",
		Aliases:     []string{"deploy.docker.image-id"},
		Destination: &dockerImageId,
	}
	fasthttpFlag = cli.BoolFlag{
		Name:    "fasthttp",
		Usage:   "(docker only) Use High performance http module",
		Aliases: []string{"deploy.docker.fasthttp"},
	}
	networkIdFlag = cli.IntFlag{
		Name:    "network-id",
		Usage:   "(docker only) network identifier (default : 2018)",
		Aliases: []string{"deploy.docker.network-id"},
		Value:   2018,
	}
	nografanaFlag = cli.BoolFlag{
		Name:    "no-grafana",
		Usage:   "(docker only) Do not make grafana container",
		Aliases: []string{"deploy.docker.no-grafana"},
	}
	useTxGenFlag = cli.BoolFlag{
		Name:    "txgen",
		Usage:   "(docker only) Add txgen container",
		Aliases: []string{"deploy.docker.tx-gen.enable"},
	}
	txGenRateFlag = cli.IntFlag{
		Name:    "txgen-rate",
		Usage:   "(docker only) txgen's rate option [default : 2000]",
		Aliases: []string{"deploy.docker.tx-gen.rate"},
		Value:   2000,
	}
	txGenConnFlag = cli.IntFlag{
		Name:    "txgen-conn",
		Usage:   "(docker only) txgen's connection size option [default : 100]",
		Aliases: []string{"deploy.docker.tx-gen.connetions"},
		Value:   100,
	}
	txGenDurFlag = cli.StringFlag{
		Name:    "txgen-dur",
		Usage:   "(docker only) txgen's duration option [default : 1m]",
		Aliases: []string{"deploy.docker.tx-gen.duration"},
		Value:   "1m",
	}
	txGenThFlag = cli.IntFlag{
		Name:    "txgen-th",
		Usage:   "(docker-only) txgen's thread size option [default : 2]",
		Aliases: []string{"deploy.docker.tx-gen.thread"},
		Value:   2,
	}

	rpcPortFlag = cli.IntFlag{
		Name:    "rpc-port",
		Usage:   "klay.conf - Klaytn node's rpc port [default: 8551] ",
		Aliases: []string{"deploy.rpc-port"},
		Value:   8551,
	}

	wsPortFlag = cli.IntFlag{
		Name:    "ws-port",
		Usage:   "klay.conf - Klaytn node's ws port [default: 8552]",
		Aliases: []string{"deploy.ws-port"},
		Value:   8552,
	}

	p2pPortFlag = cli.IntFlag{
		Name:    "p2p-port",
		Usage:   "klay.conf - Klaytn node's p2p port [default: 32323]",
		Aliases: []string{"deploy.p2p-port"},
		Value:   32323,
	}

	dataDirFlag = cli.StringFlag{
		Name:    "data-dir",
		Usage:   "klay.conf - Klaytn node's data directory path [default : /var/klay/data]",
		Aliases: []string{"deploy.data-dir"},
		Value:   "/var/klay/data",
	}

	logDirFlag = cli.StringFlag{
		Name:    "log-dir",
		Usage:   "klay.conf - Klaytn node's log directory path [default : /var/klay/log]",
		Aliases: []string{"deploy.log-dir"},
		Value:   "/var/klay/log",
	}

	// Governance flags
	governanceFlag = cli.BoolFlag{
		Name:    "governance",
		Usage:   "governance field is added in the genesis file if this flag is set",
		Aliases: []string{"genesis.governance"},
	}

	govModeFlag = cli.StringFlag{
		Name:    "gov-mode",
		Usage:   "governance mode (none, single, ballot) [default: none]",
		Aliases: []string{"genesis.governance-mode"},
		Value:   params.DefaultGovernanceMode,
	}

	governingNodeFlag = cli.StringFlag{
		Name:    "governing-node",
		Usage:   "the governing node [default: 0x0000000000000000000000000000000000000000]",
		Aliases: []string{"genesis.governing-node"},
		Value:   params.DefaultGoverningNode,
	}

	rewardMintAmountFlag = cli.StringFlag{
		Name:    "reward-mint-amount",
		Usage:   "governance minting amount",
		Aliases: []string{"genesis.reward.mint-amount"},
		Value:   "9600000000000000000",
	}
	rewardRatioFlag = cli.StringFlag{
		Name:    "reward-ratio",
		Usage:   "governance ratio [default: 100/0/0]",
		Aliases: []string{"genesis.reward.ratio"},
		Value:   params.DefaultRatio,
	}
	rewardGiniCoeffFlag = cli.BoolFlag{
		Name:    "reward-gini-coeff",
		Usage:   "governance gini-coefficient",
		Aliases: []string{"genesis.reward.gini-coefficient"},
	}
	rewardDeferredTxFeeFlag = cli.BoolFlag{
		Name:    "reward-deferred-tx",
		Usage:   "governance deferred transaction",
		Aliases: []string{"genesis.reward.deferred-tx"},
	}
	rewardStakingFlag = cli.Uint64Flag{
		Name:    "reward-staking-interval",
		Usage:   "reward staking update interval flag",
		Aliases: []string{"genesis.reward.staking-interval"},
		Value:   86400,
	}
	rewardProposerFlag = cli.Uint64Flag{
		Name:    "reward-proposer-interval",
		Usage:   "reward proposer update interval flag",
		Aliases: []string{"genesis.reward.proposer-inteval"},
		Value:   3600,
	}
	rewardMinimumStakeFlag = cli.StringFlag{
		Name:    "reward-minimum-stake",
		Usage:   "reward minimum stake flag",
		Aliases: []string{"genesis.reward.minimum-stake"},
		Value:   "2000000",
	}

	magmaLowerBoundBaseFeeFlag = cli.Uint64Flag{
		Name:    "lower-bound-base-fee",
		Usage:   "lower bound base fee flag",
		Aliases: []string{"genesis.kip71.lower-bound-base-fee"},
		Value:   params.DefaultLowerBoundBaseFee,
	}

	magmaUpperBoundBaseFeeFlag = cli.Uint64Flag{
		Name:    "upper-bound-base-fee",
		Usage:   "upper bound base fee flag",
		Aliases: []string{"genesis.kip71.upper-bound-base-fee"},
		Value:   params.DefaultUpperBoundBaseFee,
	}

	magmaGasTarget = cli.Uint64Flag{
		Name:    "gas-target",
		Usage:   "gas target flag",
		Aliases: []string{"genesis.kip71.gas-target"},
		Value:   params.DefaultGasTarget,
	}

	magmaMaxBlockGasUsedForBaseFee = cli.Uint64Flag{
		Name:    "block-gas-limit",
		Usage:   "block gas limit flag",
		Aliases: []string{"genesis.kip71.block-gas-limit"},
		Value:   params.DefaultMaxBlockGasUsedForBaseFee,
	}

	magmaBaseFeeDenominator = cli.Uint64Flag{
		Name:    "base-fee-denominator",
		Usage:   "base fee denominator flag",
		Aliases: []string{"genesis.kip71.base-fee-denominator"},
		Value:   params.DefaultBaseFeeDenominator,
	}

	istEpochFlag = cli.Uint64Flag{
		Name:    "ist-epoch",
		Usage:   "governance epoch [default: 604800]",
		Aliases: []string{"genesis.consensus.istanbul.epoch"},
		Value:   params.DefaultEpoch,
	}

	istProposerPolicyFlag = cli.Uint64Flag{
		Name:    "ist-proposer-policy",
		Usage:   "governance proposer policy (0: RoundRobin, 1: Sticky, 2: WeightedRandom) [default: 0]",
		Aliases: []string{"genesis.consensus.istanbul.policy"},
		Value:   params.DefaultProposerPolicy,
	}

	istSubGroupFlag = cli.Uint64Flag{
		Name:    "ist-subgroup",
		Usage:   "governance subgroup size [default: 21]",
		Aliases: []string{"genesis.consensus.istanbul.subgroup"},
		Value:   params.DefaultSubGroupSize,
	}

	cliqueEpochFlag = cli.Uint64Flag{
		Name:    "clique-epoch",
		Usage:   "clique epoch",
		Aliases: []string{"genesis.consensus.clique.epoch"},
		Value:   params.DefaultEpoch,
	}

	cliquePeriodFlag = cli.Uint64Flag{
		Name:    "clique-period",
		Usage:   "clique period",
		Aliases: []string{"genesis.consensus.clique.period"},
		Value:   params.DefaultPeriod,
	}

	istanbulCompatibleBlockNumberFlag = cli.Int64Flag{
		Name:    "istanbul-compatible-blocknumber",
		Usage:   "istanbulCompatible blockNumber",
		Aliases: []string{"genesis.hardfork.istanbul-compatible-blocknumber"},
		Value:   0,
	}

	londonCompatibleBlockNumberFlag = cli.Int64Flag{
		Name:    "london-compatible-blocknumber",
		Usage:   "londonCompatible blockNumber",
		Aliases: []string{"genesis.hardfork.london-compatible-blocknumber"},
		Value:   0,
	}

	ethTxTypeCompatibleBlockNumberFlag = cli.Int64Flag{
		Name:    "eth-tx-type-compatible-blocknumber",
		Usage:   "ethTxTypeCompatible blockNumber",
		Aliases: []string{"genesis.hardfork.eth-tx-type-compatible-blocknumber"},
		Value:   0,
	}

	magmaCompatibleBlockNumberFlag = cli.Int64Flag{
		Name:    "magma-compatible-blocknumber",
		Usage:   "magmaCompatible blockNumber",
		Aliases: []string{"genesis.hardfork.magma-compatible-blocknumber"},
		Value:   0,
	}

	koreCompatibleBlockNumberFlag = cli.Int64Flag{
		Name:    "kore-compatible-blocknumber",
		Usage:   "koreCompatible blockNumber",
		Aliases: []string{"genesis.hardfork.kore-compatible-blocknumber"},
		Value:   0,
	}
)
