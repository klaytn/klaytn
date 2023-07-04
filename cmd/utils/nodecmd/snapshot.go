// Modifications Copyright 2022 The klaytn Authors
// Copyright 2020 The go-ethereum Authors
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
// This file is derived from cmd/utils/nodecmd/snapshot.go (2022/07/08).
// Modified and improved for the klaytn development.

package nodecmd

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/cmd/utils"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/snapshot"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"gopkg.in/urfave/cli.v1"
)

var SnapshotCommand = cli.Command{
	Name:        "snapshot",
	Usage:       "A set of commands based on the snapshot",
	Description: "",
	Subcommands: []cli.Command{
		{
			Name:      "verify-state",
			Usage:     "Recalculate state hash based on the snapshot for verification",
			ArgsUsage: "<root>",
			Action:    utils.MigrateFlags(verifyState),
			Flags: []cli.Flag{
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
			},
			Description: `
klay snapshot verify-state <state-root>
will traverse the whole accounts and storages set based on the specified
snapshot and recalculate the root hash of state for verification.
In other words, this command does the snapshot to trie conversion.
`,
		},
		{
			Name:      "trace-trie",
			Usage:     "trace all trie nodes for verification",
			ArgsUsage: "<root>",
			Action:    utils.MigrateFlags(traceTrie),
			Flags: []cli.Flag{
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
			},
			Description: `
klaytn statedb trace-trie <state-root>
trace all account and storage nodes to find missing data
during the migration process.
Start tracing from the state root of the last block,
reading all nodes and logging the missing nodes.
`,
		},
		{
			Name:      "iterate-triedb",
			Usage:     "Iterate StateTrie DB for node count",
			ArgsUsage: "<root>",
			Action:    utils.MigrateFlags(iterateTrie),
			Flags: []cli.Flag{
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
			},
			Description: `
klaytn statedb iterate-triedb
Count the number of nodes in the state-trie db.
`,
		},
	},
}

var (
	midAccountCnt  = uint64(0)
	midStorageCnt  = uint64(0)
	codeCnt        = uint64(0)
	leafAccountCnt = uint64(0)
	leafStorageCnt = uint64(0)
	unknownCnt     = uint64(0)
	mutex          = &sync.Mutex{}
)

// getConfig returns a database config with the given context.
func getConfig(ctx *cli.Context) *database.DBConfig {
	return &database.DBConfig{
		Dir:                "chaindata",
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
}

// parseRoot parse the given hex string to hash.
func parseRoot(input string) (common.Hash, error) {
	var h common.Hash
	if err := h.UnmarshalText([]byte(input)); err != nil {
		return h, err
	}
	return h, nil
}

// verifyState verifies if the stored snapshot data is correct or not.
// if a root hash isn't given, the root hash of current block is investigated.
func verifyState(ctx *cli.Context) error {
	stack := MakeFullNode(ctx)
	db := stack.OpenDatabase(getConfig(ctx))
	head := db.ReadHeadBlockHash()
	if head == (common.Hash{}) {
		// Corrupt or empty database, init from scratch
		return errors.New("empty database")
	}
	// Make sure the entire head block is available
	headBlock := db.ReadBlockByHash(head)
	if headBlock == nil {
		return fmt.Errorf("head block missing: %v", head.String())
	}

	snaptree, err := snapshot.New(db, statedb.NewDatabase(db), 256, headBlock.Root(), false, false, false)
	if err != nil {
		logger.Error("Failed to open snapshot tree", "err", err)
		return err
	}
	if ctx.NArg() > 1 {
		logger.Error("Too many arguments given")
		return errors.New("too many arguments")
	}
	root := headBlock.Root()
	if ctx.NArg() == 1 {
		root, err = parseRoot(ctx.Args().First())
		if err != nil {
			logger.Error("Failed to resolve state root", "err", err)
			return err
		}
	}
	if err := snaptree.Verify(root); err != nil {
		logger.Error("Failed to verify state", "root", root, "err", err)
		return err
	}
	logger.Info("Verified the state", "root", root)
	return nil
}

func traceTrie(ctx *cli.Context) error {
	var childWait, logWait sync.WaitGroup

	stack := MakeFullNode(ctx)
	dbm := stack.OpenDatabase(getConfig(ctx))
	head := dbm.ReadHeadBlockHash()
	if head == (common.Hash{}) {
		// Corrupt or empty database, init from scratch
		return errors.New("empty database")
	}
	// Make sure the entire head block is available
	tmpHeadBlock := dbm.ReadBlockByHash(head)
	if tmpHeadBlock == nil {
		return fmt.Errorf("tmp head block missing: %v", head.String())
	}

	blockNumber := (tmpHeadBlock.NumberU64() / 128) * 128
	headBlock := dbm.ReadBlockByNumber(blockNumber)
	if headBlock == nil {
		return fmt.Errorf("head block missing: %v", head.String())
	}

	root := headBlock.Root()
	if root == (common.Hash{}) {
		// Corrupt or empty database, init from scratch
		return errors.New("empty root")
	}

	logger.Info("Trace Start", "BlockNum", blockNumber)

	sdb, err := state.New(root, state.NewDatabase(dbm), nil, nil)
	if err != nil {
		return fmt.Errorf("Failed to open newDB trie : %v", err)
	}
	trieDB := sdb.Database().TrieDB()

	// Get root-node childrens to create goroutine by number of childrens
	children, err := trieDB.NodeChildren(root)
	if err != nil {
		return fmt.Errorf("Fail get childrens of root : %v", err)
	}

	midAccountCnt, midStorageCnt, codeCnt, leafAccountCnt, leafStorageCnt, unknownCnt = 0, 0, 0, 0, 0, 0
	endFlag := false

	childWait.Add(len(children))
	logWait.Add(1)
	// create logging goroutine
	go func() {
		defer logWait.Done()
		for !endFlag {
			time.Sleep(time.Second * 5)
			logger.Info("Trie Tracer", "AccNode", midAccountCnt, "AccLeaf", leafAccountCnt, "StrgNode", midStorageCnt, "StrgLeaf", leafStorageCnt, "Unknown", unknownCnt, "CodeAcc", codeCnt)
		}
		logger.Info("Trie Tracer Finished", "AccNode", midAccountCnt, "AccLeaf", leafAccountCnt, "StrgNode", midStorageCnt, "StrgLeaf", leafStorageCnt, "Unknown", unknownCnt, "CodeAcc", codeCnt)
	}()

	// Create goroutine by number of childrens
	for _, child := range children {
		go func(child common.Hash) {
			defer childWait.Done()
			doTraceTrie(sdb.Database(), child)
		}(child)
	}

	childWait.Wait()
	endFlag = true
	logWait.Wait()
	return nil
}

func doTraceTrie(db state.Database, root common.Hash) (resultErr error) {
	logger.Info("Trie Tracer Start", "Hash Root", root)
	// Create and iterate a state trie rooted in a sub-node
	oldState, err := state.New(root, db, nil, nil)
	if err != nil {
		logger.Error("can not open trie DB", err.Error())
		panic(err)
	}

	oldIt := state.NewNodeIterator(oldState)

	for oldIt.Next() {
		mutex.Lock()
		switch oldIt.Type {
		case "state":
			midAccountCnt++
		case "storage":
			midStorageCnt++
		case "code":
			codeCnt++
		case "state_leaf":
			leafAccountCnt++
		case "storage_leaf":
			leafStorageCnt++
		default:
			unknownCnt++
		}
		mutex.Unlock()
	}
	if oldIt.Error != nil {
		logger.Error("Error Finished", "Root Hash", root, "Message", oldIt.Error)
	}
	logger.Info("Trie Tracer Finished", "Root Hash", root, "AccNode", midAccountCnt, "AccLeaf", leafAccountCnt, "StrgNode", midStorageCnt, "StrgLeaf", leafStorageCnt, "Unknown", unknownCnt, "CodeAcc", codeCnt)
	return nil
}

func iterateTrie(ctx *cli.Context) error {
	stack := MakeFullNode(ctx)
	dbm := stack.OpenDatabase(getConfig(ctx))
	sdb, err := state.New(common.Hash{}, state.NewDatabase(dbm), nil, nil)
	if err != nil {
		return fmt.Errorf("Failed to open newDB trie : %v", err)
	}

	logger.Info("TrieDB Iterator Start", "node count : all node count, nil node count : key or value is nil node count")
	cnt, nilCnt := uint64(0), uint64(0)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			logger.Info("TrieDB Iterator", "node count", cnt, "nil node count", nilCnt)
		}
	}()

	it := sdb.Database().TrieDB().DiskDB().GetStateTrieDB().NewIterator(nil, nil)
	defer it.Release()
	for it.Next() {
		cnt++
		if it.Key() == nil || it.Value() == nil || bytes.Equal(it.Key(), []byte("")) || bytes.Equal(it.Value(), []byte("")) {
			nilCnt++
		}
	}
	logger.Info("TrieDB Iterator finished", "total node count", cnt, "nil node count", nilCnt)
	return nil
}
