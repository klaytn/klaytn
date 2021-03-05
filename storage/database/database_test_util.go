// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package database

import (
	"fmt"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

// NewLevelDBManagerForTest returns a DBManager, consisted of only LevelDB.
// It also accepts LevelDB option, opt.Options.
func NewLevelDBManagerForTest(dbc *DBConfig, levelDBOption *opt.Options) (DBManager, error) {
	dbm := newDatabaseManager(dbc)

	checkDBEntryConfigRatio()

	var ldb *levelDB
	var err error
	for i := 0; i < int(databaseEntryTypeSize); i++ {
		if dbm.config.SingleDB {
			if i == 0 {
				ldb, err = NewLevelDBWithOption(dbc.Dir, levelDBOption)
			}
		} else {
			dir := filepath.Join(dbc.Dir, dbBaseDirs[i])
			LDBOption := getLevelDBOptionByDBType(levelDBOption, DBEntryType(i))
			LDBOption.Compression = getCompressionType(dbc.LevelDBCompression, DBEntryType(i))

			ldb, err = NewLevelDBWithOption(dir, LDBOption)

			fmt.Printf("Database: %-15s Compression: %-15s DisableBufferPool: %v \n", dbBaseDirs[i], LDBOption.Compression, levelDBOption.DisableBufferPool)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create new LevelDB with options. err: %v", err)
		}

		dbm.dbs[i] = ldb
	}

	return dbm, nil
}

// getLevelDBOptionByDBType returns scaled LevelDB option from the given LevelDB option.
// Some settings are not changed since they are not globally shared resources.
// e.g., NoSync or CompactionTableSizeMultiplier
func getLevelDBOptionByDBType(levelDBOption *opt.Options, i DBEntryType) *opt.Options {
	copiedLevelDBOption := *levelDBOption
	ratio := dbConfigRatio[i]
	copiedLevelDBOption.WriteBuffer = levelDBOption.WriteBuffer * ratio / 100

	return &copiedLevelDBOption
}
