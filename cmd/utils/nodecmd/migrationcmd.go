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
	"github.com/urfave/cli/v2"
)

var (
	dbMigrationFlags = append(utils.DBMigrationSrcFlags, utils.DBMigrationDstFlags...)

	MigrationCommand = &cli.Command{
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
		Subcommands: []*cli.Command{
			{
				Name:   "start",
				Usage:  "Start db migration",
				Flags:  dbMigrationFlags,
				Action: startMigration,
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
		Dir:                ctx.String(utils.DataDirFlag.Name),
		DBType:             database.DBType(ctx.String(utils.DbTypeFlag.Name)).ToValid(),
		SingleDB:           ctx.Bool(utils.SingleDBFlag.Name),
		NumStateTrieShards: ctx.Uint(utils.NumStateTrieShardsFlag.Name),
		OpenFilesLimit:     database.GetOpenFilesLimit(),

		LevelDBCacheSize:    ctx.Int(utils.LevelDBCacheSizeFlag.Name),
		LevelDBCompression:  database.LevelDBCompressionType(ctx.Int(utils.LevelDBCompressionTypeFlag.Name)),
		EnableDBPerfMetrics: !ctx.Bool(utils.DBNoPerformanceMetricsFlag.Name),

		DynamoDBConfig: &database.DynamoDBConfig{
			TableName:          ctx.String(utils.DynamoDBTableNameFlag.Name),
			Region:             ctx.String(utils.DynamoDBRegionFlag.Name),
			IsProvisioned:      ctx.Bool(utils.DynamoDBIsProvisionedFlag.Name),
			ReadCapacityUnits:  ctx.Int64(utils.DynamoDBReadCapacityFlag.Name),
			WriteCapacityUnits: ctx.Int64(utils.DynamoDBWriteCapacityFlag.Name),
			PerfCheck:          !ctx.Bool(utils.DBNoPerformanceMetricsFlag.Name),
		},

		RocksDBConfig: &database.RocksDBConfig{
			CacheSize:                 ctx.Uint64(utils.RocksDBCacheSizeFlag.Name),
			DumpMallocStat:            ctx.Bool(utils.RocksDBDumpMallocStatFlag.Name),
			DisableMetrics:            ctx.Bool(utils.RocksDBDisableMetricsFlag.Name),
			Secondary:                 ctx.Bool(utils.RocksDBSecondaryFlag.Name),
			CompressionType:           ctx.String(utils.RocksDBCompressionTypeFlag.Name),
			BottommostCompressionType: ctx.String(utils.RocksDBBottommostCompressionTypeFlag.Name),
			FilterPolicy:              ctx.String(utils.RocksDBFilterPolicyFlag.Name),
			MaxOpenFiles:              ctx.Int(utils.RocksDBMaxOpenFilesFlag.Name),
			CacheIndexAndFilter:       ctx.Bool(utils.RocksDBCacheIndexAndFilterFlag.Name),
		},
	}
	if len(srcDBC.DBType) == 0 { // changed to invalid type
		return nil, nil, errors.New("srcDB is not specified or invalid : " + ctx.String(utils.DbTypeFlag.Name))
	}

	// dstDB
	dstDBC := &database.DBConfig{
		Dir:                ctx.String(utils.DstDataDirFlag.Name),
		DBType:             database.DBType(ctx.String(utils.DstDbTypeFlag.Name)).ToValid(),
		SingleDB:           ctx.Bool(utils.DstSingleDBFlag.Name),
		NumStateTrieShards: ctx.Uint(utils.DstNumStateTrieShardsFlag.Name),
		OpenFilesLimit:     database.GetOpenFilesLimit(),

		LevelDBCacheSize:    ctx.Int(utils.DstLevelDBCacheSizeFlag.Name),
		LevelDBCompression:  database.LevelDBCompressionType(ctx.Int(utils.DstLevelDBCompressionTypeFlag.Name)),
		EnableDBPerfMetrics: !ctx.Bool(utils.DBNoPerformanceMetricsFlag.Name),

		DynamoDBConfig: &database.DynamoDBConfig{
			TableName:          ctx.String(utils.DstDynamoDBTableNameFlag.Name),
			Region:             ctx.String(utils.DstDynamoDBRegionFlag.Name),
			IsProvisioned:      ctx.Bool(utils.DstDynamoDBIsProvisionedFlag.Name),
			ReadCapacityUnits:  ctx.Int64(utils.DstDynamoDBReadCapacityFlag.Name),
			WriteCapacityUnits: ctx.Int64(utils.DstDynamoDBWriteCapacityFlag.Name),
			PerfCheck:          !ctx.Bool(utils.DBNoPerformanceMetricsFlag.Name),
		},

		RocksDBConfig: &database.RocksDBConfig{
			CacheSize:                 ctx.Uint64(utils.DstRocksDBCacheSizeFlag.Name),
			DumpMallocStat:            ctx.Bool(utils.DstRocksDBDumpMallocStatFlag.Name),
			DisableMetrics:            ctx.Bool(utils.DstRocksDBDisableMetricsFlag.Name),
			Secondary:                 ctx.Bool(utils.DstRocksDBSecondaryFlag.Name),
			CompressionType:           ctx.String(utils.DstRocksDBCompressionTypeFlag.Name),
			BottommostCompressionType: ctx.String(utils.DstRocksDBBottommostCompressionTypeFlag.Name),
			FilterPolicy:              ctx.String(utils.DstRocksDBFilterPolicyFlag.Name),
			MaxOpenFiles:              ctx.Int(utils.DstRocksDBMaxOpenFilesFlag.Name),
			CacheIndexAndFilter:       ctx.Bool(utils.DstRocksDBCacheIndexAndFilterFlag.Name),
		},
	}
	if len(dstDBC.DBType) == 0 { // changed to invalid type
		return nil, nil, errors.New("dstDB is not specified or invalid : " + ctx.String(utils.DstDbTypeFlag.Name))
	}

	return srcDBC, dstDBC, nil
}

// TODO When it is stopped, store previous db migration info.
//      Continue migration on next call with the same setting.
