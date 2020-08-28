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
The migration commands migrates a DB to a different kind of DB.
(e.g. LevelDB -> BadgerDB, LevelDB -> DynamoDB)
Note: This feature is only provided when srcDB is single LevelDB.
Note: Do not use db migration while a node is executing
`,
		Subcommands: []cli.Command{
			{
				Name:   "start",
				Usage:  "Start db migration",
				Flags:  dbMigrationFlags,
				Action: utils.MigrateFlags(startMigration),
				Description: `
This command starts migration.
Note: This feature is only provided when srcDB is single LevelDB..`,
			},
		},
	}
)

func startMigration(ctx *cli.Context) error {
	srcDBManager, dstDBManager, err := createDBConfig(ctx)
	if err != nil {
		return err
	}

	s, _ := json.Marshal(srcDBManager)
	d, _ := json.Marshal(dstDBManager)
	logger.Info("created DBConfig", "\nsrcDB", string(s), "\ndstDB", string(d))

	return errors.New("Not implemented")
}

func createDBConfig(ctx *cli.Context) (*database.DBConfig, *database.DBConfig, error) {
	// TODO enable for all dbs
	if !ctx.GlobalBool(utils.SingleDBFlag.Name) || !ctx.GlobalBool(utils.DstSingleDBFlag.Name) {
		return nil, nil, errors.New("this feature is provided for single db only")
	}

	// srcDB
	srcDBC := &database.DBConfig{
		Dir:                ctx.GlobalString(utils.DataDirFlag.Name),
		DBType:             database.DBType(ctx.GlobalString(utils.DbTypeFlag.Name)).ToValid(),
		SingleDB:           ctx.GlobalBool(utils.SingleDBFlag.Name),
		NumStateTrieShards: ctx.GlobalUint(utils.NumStateTrieShardsFlag.Name),
		OpenFilesLimit:     database.GetOpenFilesLimit(),

		LevelDBCacheSize:   ctx.GlobalInt(utils.LevelDBCacheSizeFlag.Name),
		LevelDBCompression: database.LevelDBCompressionType(ctx.GlobalInt(utils.LevelDBCompressionTypeFlag.Name)),

		DynamoDBConfig: &database.DynamoDBConfig{
			TableName:          ctx.GlobalString(utils.DynamoDBTableNameFlag.Name),
			Region:             ctx.GlobalString(utils.DynamoDBRegionFlag.Name),
			IsProvisioned:      ctx.GlobalBool(utils.DynamoDBIsProvisionedFlag.Name),
			ReadCapacityUnits:  ctx.GlobalInt64(utils.DynamoDBReadCapacityFlag.Name),
			WriteCapacityUnits: ctx.GlobalInt64(utils.DynamoDBWriteCapacityFlag.Name),
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

		LevelDBCacheSize:   ctx.GlobalInt(utils.DstLevelDBCacheSizeFlag.Name),
		LevelDBCompression: database.LevelDBCompressionType(ctx.GlobalInt(utils.DstLevelDBCompressionTypeFlag.Name)),

		DynamoDBConfig: &database.DynamoDBConfig{
			TableName:          ctx.GlobalString(utils.DstDynamoDBTableNameFlag.Name),
			Region:             ctx.GlobalString(utils.DstDynamoDBRegionFlag.Name),
			IsProvisioned:      ctx.GlobalBool(utils.DstDynamoDBIsProvisionedFlag.Name),
			ReadCapacityUnits:  ctx.GlobalInt64(utils.DstDynamoDBReadCapacityFlag.Name),
			WriteCapacityUnits: ctx.GlobalInt64(utils.DstDynamoDBWriteCapacityFlag.Name),
		},
	}
	if len(dstDBC.DBType) == 0 { // changed to invalid type
		return nil, nil, errors.New("dstDB is not specified or invalid : " + ctx.GlobalString(utils.DstDbTypeFlag.Name))
	}

	return srcDBC, dstDBC, nil
}
