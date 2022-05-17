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
	"errors"
	"os"
	"strings"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"gopkg.in/urfave/cli.v1"
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
			utils.SingleDBFlag,
			utils.NumStateTrieShardsFlag,
			utils.DynamoDBTableNameFlag,
			utils.DynamoDBRegionFlag,
			utils.DynamoDBIsProvisionedFlag,
			utils.DynamoDBReadCapacityFlag,
			utils.DynamoDBWriteCapacityFlag,
			utils.DynamoDBReadOnlyFlag,
			utils.LevelDBCompressionTypeFlag,
			utils.DataDirFlag,
			utils.OverwriteGenesisFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The init command initializes a new genesis block and definition for the network.
This is a destructive action and changes the network in which you will be
participating.

It expects the genesis file as argument.`,
	}

	DumpGenesisCommand = cli.Command{
		Action:    utils.MigrateFlags(dumpGenesis),
		Name:      "dumpgenesis",
		Usage:     "Dumps genesis block JSON configuration to stdout",
		ArgsUsage: "",
		Flags: []cli.Flag{
			utils.CypressFlag,
			utils.BaobabFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The dumpgenesis command dumps the genesis block configuration in JSON format to stdout.`,
	}
)

// initGenesis will initialise the given JSON format genesis file and writes it as
// the zero'd block (i.e. genesis) or will fail hard if it can't succeed.
func initGenesis(ctx *cli.Context) error {
	// Make sure we have a valid genesis JSON
	genesisPath := ctx.Args().First()
	if len(genesisPath) == 0 {
		logger.Crit("Must supply path to genesis JSON file")
	}
	file, err := os.Open(genesisPath)
	if err != nil {
		logger.Crit("Failed to read genesis file", "err", err)
	}
	defer file.Close()

	genesis := new(blockchain.Genesis)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		logger.Crit("Invalid genesis file", "err", err)
		return err
	}
	if genesis.Config == nil {
		logger.Crit("Genesis config is not set")
	}

	// Update undefined config with default values
	genesis.Config.SetDefaults()

	// Validate config values
	if err := ValidateGenesisConfig(genesis); err != nil {
		logger.Crit("Invalid genesis", "err", err)
	}

	// Set genesis.Governance and reward intervals
	govSet := governance.GetGovernanceItemsFromChainConfig(genesis.Config)
	govItemBytes, err := json.Marshal(govSet.Items())
	if err != nil {
		logger.Crit("Failed to json marshaling governance data", "err", err)
	}
	if genesis.Governance, err = rlp.EncodeToBytes(govItemBytes); err != nil {
		logger.Crit("Failed to encode initial settings. Check your genesis.json", "err", err)
	}
	params.SetStakingUpdateInterval(genesis.Config.Governance.Reward.StakingUpdateInterval)
	params.SetProposerUpdateInterval(genesis.Config.Governance.Reward.ProposerUpdateInterval)

	// Open an initialise both full and light databases
	stack := MakeFullNode(ctx)
	parallelDBWrite := !ctx.GlobalIsSet(utils.NoParallelDBWriteFlag.Name)
	singleDB := ctx.GlobalIsSet(utils.SingleDBFlag.Name)
	numStateTrieShards := ctx.GlobalUint(utils.NumStateTrieShardsFlag.Name)
	overwriteGenesis := ctx.GlobalBool(utils.OverwriteGenesisFlag.Name)

	dbtype := database.DBType(ctx.GlobalString(utils.DbTypeFlag.Name)).ToValid()
	if len(dbtype) == 0 {
		logger.Crit("invalid dbtype", "dbtype", ctx.GlobalString(utils.DbTypeFlag.Name))
	}

	var dynamoDBConfig *database.DynamoDBConfig
	if dbtype == database.DynamoDB {
		dynamoDBConfig = &database.DynamoDBConfig{
			TableName:          ctx.GlobalString(utils.DynamoDBTableNameFlag.Name),
			Region:             ctx.GlobalString(utils.DynamoDBRegionFlag.Name),
			IsProvisioned:      ctx.GlobalBool(utils.DynamoDBIsProvisionedFlag.Name),
			ReadCapacityUnits:  ctx.GlobalInt64(utils.DynamoDBReadCapacityFlag.Name),
			WriteCapacityUnits: ctx.GlobalInt64(utils.DynamoDBWriteCapacityFlag.Name),
			ReadOnly:           ctx.GlobalBool(utils.DynamoDBReadOnlyFlag.Name),
		}
	}

	for _, name := range []string{"chaindata"} { // Removed "lightchaindata" since Klaytn doesn't use it
		dbc := &database.DBConfig{Dir: name, DBType: dbtype, ParallelDBWrite: parallelDBWrite,
			SingleDB: singleDB, NumStateTrieShards: numStateTrieShards,
			LevelDBCacheSize: 0, OpenFilesLimit: 0, DynamoDBConfig: dynamoDBConfig}
		chainDB := stack.OpenDatabase(dbc)

		// Initialize DeriveSha implementation
		blockchain.InitDeriveSha(genesis.Config.DeriveShaImpl)

		_, hash, err := blockchain.SetupGenesisBlock(chainDB, genesis, params.UnusedNetworkId, false, overwriteGenesis)
		if err != nil {
			logger.Crit("Failed to write genesis block", "err", err)
		}

		// Write governance items to database
		// If governance data already exist, it'll be skipped with an error log and will not return an error
		gov := governance.NewMixedEngineNoInit(genesis.Config, chainDB)
		if err := gov.WriteGovernance(0, govSet, governance.NewGovernanceSet()); err != nil {
			logger.Crit("Failed to write governance items", "err", err)
		}

		logger.Info("Successfully wrote genesis state", "database", name, "hash", hash.String())
		chainDB.Close()
	}
	return nil
}

func dumpGenesis(ctx *cli.Context) error {
	genesis := utils.MakeGenesis(ctx)
	if genesis == nil {
		genesis = blockchain.DefaultGenesisBlock()
	}
	if err := json.NewEncoder(os.Stdout).Encode(genesis); err != nil {
		logger.Crit("could not encode genesis")
	}
	return nil
}

func ValidateGenesisConfig(g *blockchain.Genesis) error {
	if g.Config.ChainID == nil {
		return errors.New("chainID is not specified")
	}

	if g.Config.Clique == nil && g.Config.Istanbul == nil {
		return errors.New("consensus engine should be configured")
	}

	if g.Config.Clique != nil && g.Config.Istanbul != nil {
		return errors.New("only one consensus engine can be configured")
	}

	if g.Config.Governance == nil || g.Config.Governance.Reward == nil {
		return errors.New("governance and reward policies should be configured")
	}

	if g.Config.Governance.Reward.ProposerUpdateInterval == 0 || g.Config.Governance.Reward.
		StakingUpdateInterval == 0 {
		return errors.New("proposerUpdateInterval and stakingUpdateInterval cannot be zero")
	}

	if g.Config.GetConsensusEngine() == params.UseIstanbul {
		if err := governance.CheckGenesisValues(g.Config); err != nil {
			return err
		}

		// TODO-Klaytn: Add validation logic for other GovernanceModes
		// Check if governingNode is properly set
		if strings.ToLower(g.Config.Governance.GovernanceMode) == "single" {
			var found bool

			istanbulExtra, err := types.ExtractIstanbulExtra(&types.Header{Extra: g.ExtraData})
			if err != nil {
				return err
			}

			for _, v := range istanbulExtra.Validators {
				if v == g.Config.Governance.GoverningNode {
					found = true
					break
				}
			}
			if !found {
				return errors.New("governingNode is not in the validator list")
			}
		}
	}
	return nil
}
