// Modifications Copyright 2020 The klaytn Authors
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

package nodecmd

import (
	"encoding/json"

	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/pkg/errors"

	"gopkg.in/urfave/cli.v1"
)

var (
	dbFlags = []cli.Flag{
		// src DB
		utils.DbTypeFlag,
		utils.SingleDBFlag,
		utils.NumStateTrieShardsFlag,
		utils.DynamoDBTableNameFlag,
		utils.DynamoDBRegionFlag,
		utils.DynamoDBIsProvisionedFlag,
		utils.DynamoDBReadCapacityFlag,
		utils.DynamoDBWriteCapacityFlag,
		utils.LevelDBCompressionTypeFlag,
		utils.DataDirFlag,
	}
	dbMigrationFlags = append(dbFlags, DBMigrationFlags...)

	MigrationCommand = cli.Command{
		Name:     "db-migration",
		Usage:    "db migration",
		Flags:    []cli.Flag{},
		Category: "DB MIGRATION COMMANDS",
		Description: `
The migration command migrates a DB to another DB.
The type of DBs can be different.
(e.g. LevelDB -> LevelDB, LevelDB -> BadgerDB, LevelDB -> DynamoDB)
Note: This feature is only provided when srcDB is single LevelDB.
Note: Do not use db migration while a node is executing.
`,
		Subcommands: []cli.Command{
			{
				Name:   "start",
				Usage:  "Start db migration",
				Flags:  dbMigrationFlags,
				Action: utils.MigrateFlags(startMigration),
				Description: `
This command starts DB migration.

Even if db dir names are changed in srcDB, the original db dir names are used in dstDB.
(e.g. use 'statetrie' instead of 'statetrie_migrated_xxxxx')
If dst db is singleDB, you should set dst.datadir or db.dst.dynamo.tablename
to the original db dir name.
(e.g. Data dir : 'chaindata/klay/statetrie', Dynamo table name : 'klaytn-statetrie')

Note: This feature is only provided when srcDB is single LevelDB.`,
			},
		},
	}
)

func startMigration(ctx *cli.Context) error {
	srcDBManager, dstDBManager, err := createDBManagerForMigration(ctx)
	if err != nil {
		return err
	}
	defer srcDBManager.Close()
	defer dstDBManager.Close()

	return srcDBManager.StartDBMigration(dstDBManager)
}

func createDBManagerForMigration(ctx *cli.Context) (database.DBManager, database.DBManager, error) {
	// create db config from ctx
	srcDBConfig, dstDBConfig, dbManagerCreationErr := createDBConfigForMigration(ctx)
	if dbManagerCreationErr != nil {
		return nil, nil, dbManagerCreationErr
	}

	// log
	s, _ := json.Marshal(srcDBConfig)
	d, _ := json.Marshal(dstDBConfig)
	logger.Info("dbManager created", "\nsrcDB", string(s), "\ndstDB", string(d))

	// create DBManager
	srcDBManager := database.NewDBManager(srcDBConfig)
	dstDBManager := database.NewDBManager(dstDBConfig)

	return srcDBManager, dstDBManager, nil
}

func createDBConfigForMigration(ctx *cli.Context) (*database.DBConfig, *database.DBConfig, error) {
	// srcDB
	srcDBC := &database.DBConfig{
		Dir:                ctx.GlobalString(utils.DataDirFlag.Name),
		DBType:             database.DBType(ctx.GlobalString(utils.DbTypeFlag.Name)).ToValid(),
		SingleDB:           ctx.GlobalBool(utils.SingleDBFlag.Name),
		NumStateTrieShards: ctx.GlobalUint(utils.NumStateTrieShardsFlag.Name),
		OpenFilesLimit:     database.GetOpenFilesLimit(),

		LevelDBCacheSize:    ctx.GlobalInt(utils.LevelDBCacheSizeFlag.Name),
		LevelDBCompression:  database.LevelDBCompressionType(ctx.GlobalInt(utils.LevelDBCompressionTypeFlag.Name)),
		EnableDBPerfMetrics: !ctx.IsSet(utils.DBNoPerformanceMetricsFlag.Name),

		DynamoDBConfig: &database.DynamoDBConfig{
			TableName:          ctx.GlobalString(utils.DynamoDBTableNameFlag.Name),
			Region:             ctx.GlobalString(utils.DynamoDBRegionFlag.Name),
			IsProvisioned:      ctx.GlobalBool(utils.DynamoDBIsProvisionedFlag.Name),
			ReadCapacityUnits:  ctx.GlobalInt64(utils.DynamoDBReadCapacityFlag.Name),
			WriteCapacityUnits: ctx.GlobalInt64(utils.DynamoDBWriteCapacityFlag.Name),
			PerfCheck:          !ctx.IsSet(utils.DBNoPerformanceMetricsFlag.Name),
		},
	}
	if len(srcDBC.DBType) == 0 { // changed to invalid type
		return nil, nil, errors.New("srcDB is not specified or invalid : " + ctx.GlobalString(utils.DbTypeFlag.Name))
	}

	// dstDB
	dstDBC := &database.DBConfig{
		Dir:                ctx.GlobalString(utils.DstDataDirFlag.Name),
		DBType:             database.DBType(ctx.GlobalString(utils.DstDbTypeFlag.Name)).ToValid(),
		SingleDB:           ctx.GlobalBool(utils.DstSingleDBFlag.Name),
		NumStateTrieShards: ctx.GlobalUint(utils.DstNumStateTrieShardsFlag.Name),
		OpenFilesLimit:     database.GetOpenFilesLimit(),

		LevelDBCacheSize:    ctx.GlobalInt(utils.DstLevelDBCacheSizeFlag.Name),
		LevelDBCompression:  database.LevelDBCompressionType(ctx.GlobalInt(utils.DstLevelDBCompressionTypeFlag.Name)),
		EnableDBPerfMetrics: !ctx.IsSet(utils.DBNoPerformanceMetricsFlag.Name),

		DynamoDBConfig: &database.DynamoDBConfig{
			TableName:          ctx.GlobalString(utils.DstDynamoDBTableNameFlag.Name),
			Region:             ctx.GlobalString(utils.DstDynamoDBRegionFlag.Name),
			IsProvisioned:      ctx.GlobalBool(utils.DstDynamoDBIsProvisionedFlag.Name),
			ReadCapacityUnits:  ctx.GlobalInt64(utils.DstDynamoDBReadCapacityFlag.Name),
			WriteCapacityUnits: ctx.GlobalInt64(utils.DstDynamoDBWriteCapacityFlag.Name),
			PerfCheck:          !ctx.IsSet(utils.DBNoPerformanceMetricsFlag.Name),
		},
	}
	if len(dstDBC.DBType) == 0 { // changed to invalid type
		return nil, nil, errors.New("dstDB is not specified or invalid : " + ctx.GlobalString(utils.DstDbTypeFlag.Name))
	}

	return srcDBC, dstDBC, nil
}

// TODO When it is stopped, store previous db migration info.
//      Continue migration on next call with the same setting.
