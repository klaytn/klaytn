// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from cmd/geth/chaincmd.go (2018/06/04).
// Modified and improved for the klaytn development.

package nodecmd

import (
	"encoding/json"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/urfave/cli"
	"os"
	"strings"
)

var logger = log.NewModuleLogger(log.CMDUtilsNodeCMD)

var (
	InitCommand = cli.Command{
		Action:    utils.MigrateFlags(initGenesis),
		Name:      "init",
		Usage:     "Bootstrap and initialize a new genesis block",
		ArgsUsage: "<genesisPath>",
		Flags: []cli.Flag{
			utils.DbTypeFlag,
			utils.NoPartitionedDBFlag,
			utils.NumStateTriePartitionsFlag,
			utils.LevelDBCompressionTypeFlag,
			utils.DataDirFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The init command initializes a new genesis block and definition for the network.
This is a destructive action and changes the network in which you will be
participating.

It expects the genesis file as argument.`,
	}
)

// initGenesis will initialise the given JSON format genesis file and writes it as
// the zero'd block (i.e. genesis) or will fail hard if it can't succeed.
func initGenesis(ctx *cli.Context) error {
	// Make sure we have a valid genesis JSON
	genesisPath := ctx.Args().First()
	if len(genesisPath) == 0 {
		log.Fatalf("Must supply path to genesis JSON file")
	}
	file, err := os.Open(genesisPath)
	if err != nil {
		log.Fatalf("Failed to read genesis file: %v", err)
	}
	defer file.Close()

	genesis := new(blockchain.Genesis)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		log.Fatalf("invalid genesis file: %v", err)
		return err
	}

	genesis = checkGenesisAndFillDefaultIfNeeded(genesis)

	if genesis.Config.Istanbul != nil {
		if err := governance.CheckGenesisValues(genesis.Config); err != nil {
			logger.Crit("Error in genesis json values", "err", err)
		}

		// Check if governingNode is properly set
		if strings.ToLower(genesis.Config.Governance.GovernanceMode) == "single" {
			if istanbulExtra, err := decodeExtra(genesis.ExtraData); err != nil {
				logger.Crit("Extra data couldn't be decoded. Please check if your extra data is correct", "err", err)
			} else {
				var found bool
				for _, v := range istanbulExtra.Validators {
					if v == genesis.Config.Governance.GoverningNode {
						found = true
						break
					}
				}
				if !found {
					logger.Crit("GoverningNode is not in the validator list. Please check if your governingNode address is correct")
				}
			}
		}
	}

	if genesis.Config.Governance.Reward.StakingUpdateInterval != 0 {
		params.SetStakingUpdateInterval(genesis.Config.Governance.Reward.StakingUpdateInterval)
	} else {
		genesis.Config.Governance.Reward.StakingUpdateInterval = params.StakingUpdateInterval()
	}

	if genesis.Config.Governance.Reward.ProposerUpdateInterval != 0 {
		params.SetProposerUpdateInterval(genesis.Config.Governance.Reward.ProposerUpdateInterval)
	} else {
		genesis.Config.Governance.Reward.ProposerUpdateInterval = params.ProposerUpdateInterval()
	}

	data := getGovernanceItemsFromGenesis(genesis)
	gbytes, err := json.Marshal(data.Items())
	if err != nil {
		logger.Crit("Failed to json marshaling governance data", "err", err)
	}
	if genesis.Governance, err = rlp.EncodeToBytes(gbytes); err != nil {
		logger.Crit("Failed to encode initial settings. Check your genesis.json", "err", err)
	}
	// Open an initialise both full and light databases
	stack := MakeFullNode(ctx)

	parallelDBWrite := !ctx.GlobalIsSet(utils.NoParallelDBWriteFlag.Name)
	partitioned := !ctx.GlobalIsSet(utils.NoPartitionedDBFlag.Name)
	numStateTriePartitions := ctx.GlobalUint(utils.NumStateTriePartitionsFlag.Name)
	for _, name := range []string{"chaindata", "lightchaindata"} {
		dbc := &database.DBConfig{Dir: name, DBType: database.LevelDB, ParallelDBWrite: parallelDBWrite,
			Partitioned: partitioned, NumStateTriePartitions: numStateTriePartitions,
			LevelDBCacheSize: 0, OpenFilesLimit: 0}
		chaindb := stack.OpenDatabase(dbc)
		// Initialize DeriveSha implementation
		blockchain.InitDeriveSha(genesis.Config.DeriveShaImpl)

		_, hash, err := blockchain.SetupGenesisBlock(chaindb, genesis, params.UnusedNetworkId, false)
		if err != nil {
			log.Fatalf("Failed to write genesis block: %v", err)
		}
		logger.Info("Successfully wrote genesis state", "database", name, "hash", hash.String())

		chaindb.Close()
	}
	return nil
}

func getGovernanceItemsFromGenesis(genesis *blockchain.Genesis) governance.GovernanceSet {
	g := governance.NewGovernanceSet()

	if genesis.Config.Governance != nil {
		governance := genesis.Config.Governance
		if err := g.SetValue(params.GovernanceMode, governance.GovernanceMode); err != nil {
			writeFailLog(params.GovernanceMode, err)
		}
		if err := g.SetValue(params.GoverningNode, governance.GoverningNode); err != nil {
			writeFailLog(params.GoverningNode, err)
		}
		if err := g.SetValue(params.UnitPrice, genesis.Config.UnitPrice); err != nil {
			writeFailLog(params.UnitPrice, err)
		}
		if err := g.SetValue(params.MintingAmount, governance.Reward.MintingAmount.String()); err != nil {
			writeFailLog(params.MintingAmount, err)
		}
		if err := g.SetValue(params.Ratio, governance.Reward.Ratio); err != nil {
			writeFailLog(params.Ratio, err)
		}
		if err := g.SetValue(params.UseGiniCoeff, governance.Reward.UseGiniCoeff); err != nil {
			writeFailLog(params.UseGiniCoeff, err)
		}
		if err := g.SetValue(params.DeferredTxFee, governance.Reward.DeferredTxFee); err != nil {
			writeFailLog(params.DeferredTxFee, err)
		}
		if err := g.SetValue(params.MinimumStake, governance.Reward.MinimumStake.String()); err != nil {
			writeFailLog(params.MinimumStake, err)
		}
		if err := g.SetValue(params.StakeUpdateInterval, governance.Reward.StakingUpdateInterval); err != nil {
			writeFailLog(params.StakeUpdateInterval, err)
		}
		if err := g.SetValue(params.ProposerRefreshInterval, governance.Reward.ProposerUpdateInterval); err != nil {
			writeFailLog(params.ProposerRefreshInterval, err)
		}
	}

	if genesis.Config.Istanbul != nil {
		istanbul := genesis.Config.Istanbul
		if err := g.SetValue(params.Epoch, istanbul.Epoch); err != nil {
			writeFailLog(params.Epoch, err)
		}
		if err := g.SetValue(params.Policy, istanbul.ProposerPolicy); err != nil {
			writeFailLog(params.Policy, err)
		}
		if err := g.SetValue(params.CommitteeSize, istanbul.SubGroupSize); err != nil {
			writeFailLog(params.CommitteeSize, err)
		}
	}
	return g
}

func writeFailLog(key int, err error) {
	msg := "Failed to set " + governance.GovernanceKeyMapReverse[key]
	logger.Crit(msg, "err", err)
}

func checkGenesisAndFillDefaultIfNeeded(genesis *blockchain.Genesis) *blockchain.Genesis {
	engine := params.UseIstanbul
	valueChanged := false
	if genesis.Config == nil {
		genesis.Config = new(params.ChainConfig)
	}

	// using Clique as a consensus engine
	if genesis.Config.Istanbul == nil && genesis.Config.Clique != nil {
		engine = params.UseClique
		if genesis.Config.Governance == nil {
			genesis.Config.Governance = governance.GetDefaultGovernanceConfig(engine)
		}
		valueChanged = true
	} else if genesis.Config.Istanbul == nil && genesis.Config.Clique == nil {
		engine = params.UseIstanbul
		genesis.Config.Istanbul = governance.GetDefaultIstanbulConfig()
		valueChanged = true
	} else if genesis.Config.Istanbul != nil && genesis.Config.Clique != nil {
		// Error case. Both istanbul and Clique exists
		logger.Crit("Both clique and istanbul configuration exists. Only one configuration can be applied. Exiting..")
	}

	// If we don't have governance config
	if genesis.Config.Governance == nil {
		genesis.Config.Governance = governance.GetDefaultGovernanceConfig(engine)
		valueChanged = true
	}

	if valueChanged {
		logger.Warn("Some input value of genesis.json have been set to default or changed")
	}
	return genesis
}

func decodeExtra(extraData []byte) (*types.IstanbulExtra, error) {
	istanbulExtra, err := types.ExtractIstanbulExtra(&types.Header{Extra: extraData})
	if err != nil {
		return nil, err
	}
	return istanbulExtra, nil
}
