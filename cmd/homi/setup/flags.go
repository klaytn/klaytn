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
	"github.com/klaytn/klaytn/blockchain/system"
	"github.com/klaytn/klaytn/params"
	"github.com/urfave/cli/v2"
)

var (
	dockerImageId string
	outputPath    string
)

var (
	homiYamlFlag = &cli.StringFlag{
		Name:    "homi-yaml",
		Usage:   "Import homi.yaml to generate the config files to run the nodes",
		Aliases: []string{"yaml"},
	}
	genTypeFlag = &cli.StringFlag{
		Name:    "gen-type",
		Usage:   "Generate environment files according to the type (docker, local, remote, deploy)",
		Value:   "docker",
		Aliases: []string{},
	}

	cypressTestFlag = &cli.BoolFlag{
		Name:    "cypress-test",
		Usage:   "Generate genesis.json similar to the one used for Cypress with shorter intervals for testing",
		Aliases: []string{},
	}
	cypressFlag = &cli.BoolFlag{
		Name:    "cypress",
		Usage:   "Generate genesis.json similar to the one used for Cypress",
		Aliases: []string{},
	}
	baobabTestFlag = &cli.BoolFlag{
		Name:    "baobab-test",
		Usage:   "Generate genesis.json similar to the one used for Baobab with shorter intervals for testing",
		Aliases: []string{},
	}
	baobabFlag = &cli.BoolFlag{
		Name:    "baobab",
		Usage:   "Generate genesis.json similar to the one used for Baobab",
		Aliases: []string{},
	}
	serviceChainFlag = &cli.BoolFlag{
		Name:    "servicechain",
		Usage:   "Generate genesis.json similar to the one used for Serivce Chain",
		Aliases: []string{},
	}
	serviceChainTestFlag = &cli.BoolFlag{
		Name:    "servicechain-test",
		Usage:   "Generate genesis.json similar to the one used for Serivce Chain with shorter intervals for testing",
		Aliases: []string{},
	}
	cliqueFlag = &cli.BoolFlag{
		Name:    "clique",
		Usage:   "Use Clique consensus",
		Aliases: []string{"clique.enable"},
	}

	numOfCNsFlag = &cli.IntFlag{
		Name:    "cn-num",
		Usage:   "Number of consensus nodes",
		Value:   0,
		Aliases: []string{"topology.cn-num"},
	}
	numOfValidatorsFlag = &cli.IntFlag{
		Name:    "validators-num",
		Usage:   "Number of validators. If not set, it will be set the number of option cn-num",
		Value:   0,
		Aliases: []string{"topology.validators-num"},
	}
	numOfPNsFlag = &cli.IntFlag{
		Name:    "pn-num",
		Usage:   "Number of proxy node",
		Value:   0,
		Aliases: []string{"topology.pn-num"},
	}
	numOfENsFlag = &cli.IntFlag{
		Name:    "en-num",
		Usage:   "Number of end-point node",
		Value:   0,
		Aliases: []string{"topology.en-num"},
	}
	numOfSCNsFlag = &cli.IntFlag{
		Name:    "scn-num",
		Usage:   "Number of service chain nodes",
		Value:   0,
		Aliases: []string{"topology.scn-num"},
	}
	numOfSPNsFlag = &cli.IntFlag{
		Name:    "spn-num",
		Usage:   "Number of service chain proxy nodes",
		Value:   0,
		Aliases: []string{"topology.spn-num"},
	}
	numOfSENsFlag = &cli.IntFlag{
		Name:    "sen-num",
		Usage:   "Number of service chain end-point nodes",
		Value:   0,
		Aliases: []string{"topology.sen-num"},
	}
	numOfTestKeyFlag = &cli.IntFlag{
		Name:    "test-num",
		Usage:   "Number of test key",
		Value:   0,
		Aliases: []string{"topology.test-num"},
	}

	genesisTypeFlag = &cli.StringFlag{
		Name:    "genesis-type",
		Usage:   "Set the type of genesis.json to generate (cypress-test, cypress, baobab-test, baobab, clique, servicechain, servicechain-test, istanbul)",
		Aliases: []string{"genesis.type"},
	}
	mnemonic = &cli.StringFlag{
		Name:  "mnemonic",
		Usage: "Use given mnemonic to derive node keys",
		Value: "",
	}

	mnemonicPath = &cli.StringFlag{
		Name:  "mnemonic-path",
		Usage: "Use given path/coin to derive node keys (format: m/44'/60'/0'/0/). Effective only if --mnemonic is given",
		Value: "eth",
	}

	cnNodeKeyDirFlag = &cli.StringFlag{
		Name:  "cn-nodekey-dir",
		Usage: "CN nodekey dir containing nodekey*",
		Value: "",
	}

	pnNodeKeyDirFlag = &cli.StringFlag{
		Name:  "pn-nodekey-dir",
		Usage: "PN nodekey dir containing nodekey*",
		Value: "",
	}

	enNodeKeyDirFlag = &cli.StringFlag{
		Name:  "en-nodekey-dir",
		Usage: "EN nodekey dir containing nodekey*",
		Value: "",
	}

	chainIDFlag = &cli.Uint64Flag{
		Name:    "chainID",
		Usage:   "ChainID",
		Value:   1000,
		Aliases: []string{"genesis.chain-id"},
	}

	serviceChainIDFlag = &cli.Uint64Flag{
		Name:    "serviceChainID",
		Usage:   "Service Chain ID",
		Value:   1001,
		Aliases: []string{"genesis.service-chain-id"},
	}

	unitPriceFlag = &cli.Uint64Flag{
		Name:    "unitPrice",
		Usage:   "Price of unit",
		Value:   0,
		Aliases: []string{"genesis.unit-price"},
	}

	deriveShaImplFlag = &cli.IntFlag{
		Name:    "deriveShaImpl",
		Usage:   "Implementation of DeriveSha [0:Original, 1:Simple, 2:Concat]",
		Value:   0,
		Aliases: []string{"genesis.derive-sha-impl"},
	}

	outputPathFlag = &cli.StringFlag{
		Name:        "output",
		Usage:       "homi's result saved at this output folder",
		Aliases:     []string{"o"},
		Value:       "homi-output",
		Destination: &outputPath,
	}

	fundingAddrFlag = &cli.StringFlag{
		Name:    "funding-addr",
		Value:   "",
		Usage:   "Give initial funds to the given comma-separated list of addresses",
		Aliases: []string{"genesis.funding-addr"},
	}

	patchAddressBookFlag = &cli.BoolFlag{
		Name:    "patch-address-book",
		Usage:   "Patch genesis AddressBook's constructContract function",
		Aliases: []string{"genesis.patched-address-book.enable"},
	}

	patchAddressBookAddrFlag = &cli.StringFlag{
		Name:    "patch-address-book-addr",
		Usage:   "The address to inject in AddressBook's constructContract function [default: first CN's address]",
		Value:   "",
		Aliases: []string{"genesis.patched-address-book.addr"},
	}
	// Deprecated: it doesn't used anymore
	subGroupSizeFlag = &cli.IntFlag{
		Name:    "subgroup-size",
		Usage:   "CN's Subgroup size",
		Value:   21,
		Aliases: []string{},
	}

	addressBookMockFlag = &cli.BoolFlag{
		Name:  "address-book-mock",
		Usage: "Allocate an AddressBookMock at the genesis block",
	}

	dockerImageIdFlag = &cli.StringFlag{
		Name:        "docker-image-id",
		Value:       "klaytn/klaytn:latest", // https://hub.docker.com/r/klaytn/klaytn
		Usage:       "Base docker image ID (Image[:tag]), e.g., klaytn/klaytn:v1.5.3",
		Aliases:     []string{"deploy.docker.image-id"},
		Destination: &dockerImageId,
	}
	fasthttpFlag = &cli.BoolFlag{
		Name:    "fasthttp",
		Usage:   "(docker only) Use High performance http module",
		Aliases: []string{"deploy.docker.fasthttp"},
	}
	networkIdFlag = &cli.IntFlag{
		Name:    "network-id",
		Usage:   "(docker only) network identifier (default : 2018)",
		Value:   2018,
		Aliases: []string{"deploy.docker.network-id"},
	}
	nografanaFlag = &cli.BoolFlag{
		Name:    "no-grafana",
		Usage:   "(docker only) Do not make grafana container",
		Aliases: []string{"deploy.docker.no-grafana"},
	}
	useTxGenFlag = &cli.BoolFlag{
		Name:    "txgen",
		Usage:   "(docker only) Add txgen container",
		Aliases: []string{"deploy.docker.tx-gen.enable"},
	}
	txGenRateFlag = &cli.IntFlag{
		Name:    "txgen-rate",
		Usage:   "(docker only) txgen's rate option [default : 2000]",
		Value:   2000,
		Aliases: []string{"deploy.docker.tx-gen.rate"},
	}
	txGenConnFlag = &cli.IntFlag{
		Name:    "txgen-conn",
		Usage:   "(docker only) txgen's connection size option [default : 100]",
		Value:   100,
		Aliases: []string{"deploy.docker.tx-gen.connetions"},
	}
	txGenDurFlag = &cli.StringFlag{
		Name:    "txgen-dur",
		Usage:   "(docker only) txgen's duration option [default : 1m]",
		Value:   "1m",
		Aliases: []string{"deploy.docker.tx-gen.duration"},
	}
	txGenThFlag = &cli.IntFlag{
		Name:    "txgen-th",
		Usage:   "(docker-only) txgen's thread size option [default : 2]",
		Value:   2,
		Aliases: []string{"deploy.docker.tx-gen.thread"},
	}

	rpcPortFlag = &cli.IntFlag{
		Name:    "rpc-port",
		Usage:   "klay.conf - Klaytn node's rpc port [default: 8551] ",
		Value:   8551,
		Aliases: []string{"deploy.rpc-port"},
	}

	wsPortFlag = &cli.IntFlag{
		Name:    "ws-port",
		Usage:   "klay.conf - Klaytn node's ws port [default: 8552]",
		Value:   8552,
		Aliases: []string{"deploy.ws-port"},
	}

	p2pPortFlag = &cli.IntFlag{
		Name:    "p2p-port",
		Usage:   "klay.conf - Klaytn node's p2p port [default: 32323]",
		Value:   32323,
		Aliases: []string{"deploy.p2p-port"},
	}

	dataDirFlag = &cli.StringFlag{
		Name:    "data-dir",
		Usage:   "klay.conf - Klaytn node's data directory path [default : /var/klay/data]",
		Value:   "/var/klay/data",
		Aliases: []string{"deploy.data-dir"},
	}

	logDirFlag = &cli.StringFlag{
		Name:    "log-dir",
		Usage:   "klay.conf - Klaytn node's log directory path [default : /var/klay/log]",
		Value:   "/var/klay/log",
		Aliases: []string{"deploy.log-dir"},
	}

	// Governance flags
	governanceFlag = &cli.BoolFlag{
		Name:    "governance",
		Usage:   "governance field is added in the genesis file if this flag is set",
		Aliases: []string{"genesis.governance"},
	}

	govModeFlag = &cli.StringFlag{
		Name:    "gov-mode",
		Usage:   "governance mode (none, single, ballot) [default: none]",
		Value:   params.DefaultGovernanceMode,
		Aliases: []string{"genesis.governance-mode"},
	}

	governingNodeFlag = &cli.StringFlag{
		Name:    "governing-node",
		Usage:   "the governing node [default: 0x0000000000000000000000000000000000000000]",
		Value:   params.DefaultGoverningNode,
		Aliases: []string{"genesis.governing-node"},
	}

	govParamContractFlag = &cli.StringFlag{
		Name:  "gov-param-contract",
		Usage: "the GovParam contract address [default: 0x0000000000000000000000000000000000000000]",
		Value: params.DefaultGovParamContract,
	}

	rewardMintAmountFlag = &cli.StringFlag{
		Name:    "reward-mint-amount",
		Usage:   "governance minting amount",
		Value:   "9600000000000000000",
		Aliases: []string{"genesis.reward.mint-amount"},
	}
	rewardRatioFlag = &cli.StringFlag{
		Name:    "reward-ratio",
		Usage:   "governance ratio [default: 100/0/0]",
		Value:   params.DefaultRatio,
		Aliases: []string{"genesis.reward.ratio"},
	}

	rewardKip82RatioFlag = &cli.StringFlag{
		Name:  "reward-kip82-ratio",
		Usage: "kip82 ratio [default: 20/80]",
		Value: params.DefaultKip82Ratio,
	}

	rewardGiniCoeffFlag = &cli.BoolFlag{
		Name:    "reward-gini-coeff",
		Usage:   "governance gini-coefficient",
		Aliases: []string{"genesis.reward.gini-coefficient"},
	}
	rewardDeferredTxFeeFlag = &cli.BoolFlag{
		Name:    "reward-deferred-tx",
		Usage:   "governance deferred transaction",
		Aliases: []string{"genesis.reward.deferred-tx"},
	}
	rewardStakingFlag = &cli.Uint64Flag{
		Name:    "reward-staking-interval",
		Usage:   "reward staking update interval flag",
		Value:   86400,
		Aliases: []string{"genesis.reward.staking-interval"},
	}
	rewardProposerFlag = &cli.Uint64Flag{
		Name:    "reward-proposer-interval",
		Usage:   "reward proposer update interval flag",
		Value:   3600,
		Aliases: []string{"genesis.reward.proposer-inteval"},
	}
	rewardMinimumStakeFlag = &cli.StringFlag{
		Name:    "reward-minimum-stake",
		Usage:   "reward minimum stake flag",
		Value:   "2000000",
		Aliases: []string{"genesis.reward.minimum-stake"},
	}

	magmaLowerBoundBaseFeeFlag = &cli.Uint64Flag{
		Name:    "lower-bound-base-fee",
		Usage:   "lower bound base fee flag",
		Value:   params.DefaultLowerBoundBaseFee,
		Aliases: []string{"genesis.kip71.lower-bound-base-fee"},
	}

	magmaUpperBoundBaseFeeFlag = &cli.Uint64Flag{
		Name:    "upper-bound-base-fee",
		Usage:   "upper bound base fee flag",
		Value:   params.DefaultUpperBoundBaseFee,
		Aliases: []string{"genesis.kip71.upper-bound-base-fee"},
	}

	magmaGasTarget = &cli.Uint64Flag{
		Name:    "gas-target",
		Usage:   "gas target flag",
		Value:   params.DefaultGasTarget,
		Aliases: []string{"genesis.kip71.gas-target"},
	}

	magmaMaxBlockGasUsedForBaseFee = &cli.Uint64Flag{
		Name:    "block-gas-limit",
		Usage:   "block gas limit flag",
		Value:   params.DefaultMaxBlockGasUsedForBaseFee,
		Aliases: []string{"genesis.kip71.block-gas-limit"},
	}

	magmaBaseFeeDenominator = &cli.Uint64Flag{
		Name:    "base-fee-denominator",
		Usage:   "base fee denominator flag",
		Value:   params.DefaultBaseFeeDenominator,
		Aliases: []string{"genesis.kip71.base-fee-denominator"},
	}

	istEpochFlag = &cli.Uint64Flag{
		Name:    "ist-epoch",
		Usage:   "governance epoch [default: 604800]",
		Value:   params.DefaultEpoch,
		Aliases: []string{"genesis.consensus.istanbul.epoch"},
	}

	istProposerPolicyFlag = &cli.Uint64Flag{
		Name:    "ist-proposer-policy",
		Usage:   "governance proposer policy (0: RoundRobin, 1: Sticky, 2: WeightedRandom) [default: 0]",
		Value:   params.DefaultProposerPolicy,
		Aliases: []string{"genesis.consensus.istanbul.policy"},
	}

	istSubGroupFlag = &cli.Uint64Flag{
		Name:    "ist-subgroup",
		Usage:   "governance subgroup size [default: 21]",
		Value:   params.DefaultSubGroupSize,
		Aliases: []string{"genesis.consensus.istanbul.subgroup"},
	}

	cliqueEpochFlag = &cli.Uint64Flag{
		Name:    "clique-epoch",
		Usage:   "clique epoch",
		Value:   params.DefaultEpoch,
		Aliases: []string{"genesis.consensus.clique.epoch"},
	}

	cliquePeriodFlag = &cli.Uint64Flag{
		Name:    "clique-period",
		Usage:   "clique period",
		Value:   params.DefaultPeriod,
		Aliases: []string{"genesis.consensus.clique.period"},
	}

	istanbulCompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "istanbul-compatible-blocknumber",
		Usage:   "istanbulCompatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.istanbul-compatible-blocknumber"},
	}

	londonCompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "london-compatible-blocknumber",
		Usage:   "londonCompatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.london-compatible-blocknumber"},
	}

	ethTxTypeCompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "eth-tx-type-compatible-blocknumber",
		Usage:   "ethTxTypeCompatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.eth-tx-type-compatible-blocknumber"},
	}

	magmaCompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "magma-compatible-blocknumber",
		Usage:   "magmaCompatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.magma-compatible-blocknumber"},
	}

	koreCompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "kore-compatible-blocknumber",
		Usage:   "koreCompatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.kore-compatible-blocknumber"},
	}

	shanghaiCompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "shanghai-compatible-blocknumber",
		Usage:   "shanghaiCompatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.shanghai-compatible-blocknumber"},
	}

	cancunCompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "cancun-compatible-blocknumber",
		Usage:   "cancunCompatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.cancun-compatible-blocknumber"},
	}

	// KIP103 hardfork is optional
	kip103CompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "kip103-compatible-blocknumber",
		Usage:   "kip103Compatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.kip103-compatible-blocknumber"},
	}

	kip103ContractAddressFlag = &cli.StringFlag{
		Name:    "kip103-contract-address",
		Usage:   "kip103 contract address",
		Aliases: []string{"genesis.hardfork.kip103-contract-address"},
	}

	randaoCompatibleBlockNumberFlag = &cli.Int64Flag{
		Name:    "randao-compatible-blocknumber",
		Usage:   "randaoCompatible blockNumber",
		Value:   0,
		Aliases: []string{"genesis.hardfork.randao-compatible-blocknumber"},
	}

	kip113ProxyAddressFlag = &cli.StringFlag{
		Name:    "kip113-proxy-contract-address",
		Usage:   "kip113 proxy contract address",
		Value:   system.Kip113ProxyAddrMock.String(),
		Aliases: []string{"genesis.hardfork.kip113-proxy-contract-address"},
	}

	kip113LogicAddressFlag = &cli.StringFlag{
		Name:    "kip113-logic-contract-address",
		Usage:   "kip113 logic contract address",
		Value:   system.Kip113LogicAddrMock.String(),
		Aliases: []string{"genesis.hardfork.kip113-logic-contract-address"},
	}

	kip113MockFlag = &cli.BoolFlag{
		Name:  "kip113-mock",
		Usage: "Allocate an Kip113Mock at the genesis block",
	}

	registryMockFlag = &cli.BoolFlag{
		Name:  "registry-mock",
		Usage: "Allocate an RegistryMock at the genesis block",
	}
)
