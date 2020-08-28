package nodecmd

import (
	"github.com/klaytn/klaytn/cmd/utils"
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
	return errors.New("Not implemented")
}
