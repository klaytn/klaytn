// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
//
// This file is derived from core/bench_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"

	"crypto/ecdsa"
)

func BenchmarkInsertChain_empty_memDB(b *testing.B) {
	benchInsertChain(b, database.MemoryDB, nil)
}
func BenchmarkInsertChain_empty_levelDB(b *testing.B) {
	benchInsertChain(b, database.LevelDB, nil)
}
func BenchmarkInsertChain_empty_badgerDB(b *testing.B) {
	benchInsertChain(b, database.BadgerDB, nil)
}

func BenchmarkInsertChain_valueTx_memDB(b *testing.B) {
	benchInsertChain(b, database.MemoryDB, genValueTx(0))
}
func BenchmarkInsertChain_valueTx_levelDB(b *testing.B) {
	benchInsertChain(b, database.LevelDB, genValueTx(0))
}
func BenchmarkInsertChain_valueTx_badgerDB(b *testing.B) {
	benchInsertChain(b, database.BadgerDB, genValueTx(0))
}

func BenchmarkInsertChain_valueTx_10kB_memDB(b *testing.B) {
	benchInsertChain(b, database.MemoryDB, genValueTx(100*1024))
}
func BenchmarkInsertChain_valueTx_10kB_levelDB(b *testing.B) {
	benchInsertChain(b, database.LevelDB, genValueTx(100*1024))
}
func BenchmarkInsertChain_valueTx_10kB_badgerDB(b *testing.B) {
	benchInsertChain(b, database.BadgerDB, genValueTx(100*1024))
}

func BenchmarkInsertChain_ring200_memDB(b *testing.B) {
	benchInsertChain(b, database.MemoryDB, genTxRing(200))
}
func BenchmarkInsertChain_ring200_levelDB(b *testing.B) {
	benchInsertChain(b, database.LevelDB, genTxRing(200))
}
func BenchmarkInsertChain_ring200_badgerDB(b *testing.B) {
	benchInsertChain(b, database.BadgerDB, genTxRing(200))
}

func BenchmarkInsertChain_ring1000_memDB(b *testing.B) {
	benchInsertChain(b, database.MemoryDB, genTxRing(1000))
}
func BenchmarkInsertChain_ring1000_levelDB(b *testing.B) {
	benchInsertChain(b, database.LevelDB, genTxRing(1000))
}
func BenchmarkInsertChain_ring1000_badgerDB(b *testing.B) {
	benchInsertChain(b, database.BadgerDB, genTxRing(1000))
}

var (
	// This is the content of the genesis block used by the benchmarks.
	benchRootKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	benchRootAddr   = crypto.PubkeyToAddress(benchRootKey.PublicKey)
	benchRootFunds  = math.BigPow(2, 100)
)

// genValueTx returns a block generator that includes a single
// value-transfer transaction with n bytes of extra data in each block.
func genValueTx(nbytes int) func(int, *BlockGen) {
	return func(i int, gen *BlockGen) {
		toaddr := common.Address{}
		data := make([]byte, nbytes)
		gas, _ := types.IntrinsicGas(data, nil, false, params.TestChainConfig.Rules(big.NewInt(0)))
		signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
		tx, _ := types.SignTx(types.NewTransaction(gen.TxNonce(benchRootAddr), toaddr, big.NewInt(1), gas, nil, data), signer, benchRootKey)
		gen.AddTx(tx)
	}
}

var (
	ringKeys  = make([]*ecdsa.PrivateKey, 1000)
	ringAddrs = make([]common.Address, len(ringKeys))
)

func init() {
	ringKeys[0] = benchRootKey
	ringAddrs[0] = benchRootAddr
	for i := 1; i < len(ringKeys); i++ {
		ringKeys[i], _ = crypto.GenerateKey()
		ringAddrs[i] = crypto.PubkeyToAddress(ringKeys[i].PublicKey)
	}
}

// genTxRing returns a block generator that sends KLAY in a ring
// among n accounts. This is creates n entries in the state database
// and fills the blocks with many small transactions.
func genTxRing(naccounts int) func(int, *BlockGen) {
	from := 0
	signer := types.LatestSignerForChainID(params.TestChainConfig.ChainID)
	return func(i int, gen *BlockGen) {
		gas := uint64(1000000)
		for {
			gas -= params.TxGas
			if gas < params.TxGas {
				break
			}
			to := (from + 1) % naccounts
			tx := types.NewTransaction(
				gen.TxNonce(ringAddrs[from]),
				ringAddrs[to],
				benchRootFunds,
				params.TxGas,
				nil,
				nil,
			)
			tx, _ = types.SignTx(tx, signer, ringKeys[from])
			gen.AddTx(tx)
			from = to
		}
	}
}

func benchInsertChain(b *testing.B, dbType database.DBType, gen func(int, *BlockGen)) {
	// 1. Create the database
	dir := genTempDirForDB(b)
	defer os.RemoveAll(dir)

	db := genDBManagerForTest(dir, dbType)
	defer db.Close()

	// 2. Generate a chain of b.N blocks using the supplied block generator function.
	gspec := Genesis{
		Config: params.TestChainConfig,
		Alloc:  GenesisAlloc{benchRootAddr: {Balance: benchRootFunds}},
	}
	genesis := gspec.MustCommit(db)
	chain, _ := GenerateChain(gspec.Config, genesis, gxhash.NewFaker(), db, b.N, gen)

	// Time the insertion of the new chain.
	// State and blocks are stored in the same DB.
	chainman, _ := NewBlockChain(db, nil, gspec.Config, gxhash.NewFaker(), vm.Config{})
	defer chainman.Stop()
	b.ReportAllocs()
	b.ResetTimer()
	if i, err := chainman.InsertChain(chain); err != nil {
		b.Fatalf("insert error (block %d): %v\n", i, err)
	}
}

// BenchmarkChainRead Series
func BenchmarkChainRead_header_10k_levelDB(b *testing.B) {
	benchReadChain(b, false, database.LevelDB, 10000)
}
func BenchmarkChainRead_header_10k_badgerDB(b *testing.B) {
	benchReadChain(b, false, database.BadgerDB, 10000)
}

func BenchmarkChainRead_full_10k_levelDB(b *testing.B) {
	benchReadChain(b, true, database.LevelDB, 10000)
}
func BenchmarkChainRead_full_10k_badgerDB(b *testing.B) {
	benchReadChain(b, true, database.BadgerDB, 10000)
}

func BenchmarkChainRead_header_100k_levelDB(b *testing.B) {
	benchReadChain(b, false, database.LevelDB, 100000)
}
func BenchmarkChainRead_header_100k_badgerDB(b *testing.B) {
	benchReadChain(b, false, database.BadgerDB, 100000)
}

func BenchmarkChainRead_full_100k_levelDB(b *testing.B) {
	benchReadChain(b, true, database.LevelDB, 100000)
}
func BenchmarkChainRead_full_100k_badgerDB(b *testing.B) {
	benchReadChain(b, true, database.BadgerDB, 100000)
}

// Disabled because of too long test time
//func BenchmarkChainRead_header_500k_levelDB(b *testing.B) {
//	benchReadChain(b, false, database.LevelDB,500000)
//}
//func BenchmarkChainRead_header_500k_badgerDB(b *testing.B) {
//	benchReadChain(b, false, database.BadgerDB, 500000)
//}
//
//func BenchmarkChainRead_full_500k_levelDB(b *testing.B) {
//	benchReadChain(b, true, database.LevelDB,500000)
//}
//func BenchmarkChainRead_full_500k_badgerDB(b *testing.B) {
//	benchReadChain(b, true, database.BadgerDB,500000)
//}

// BenchmarkChainWrite Series
func BenchmarkChainWrite_header_10k_levelDB(b *testing.B) {
	benchWriteChain(b, false, database.LevelDB, 10000)
}
func BenchmarkChainWrite_header_10k_badgerDB(b *testing.B) {
	benchWriteChain(b, false, database.BadgerDB, 10000)
}

func BenchmarkChainWrite_full_10k_levelDB(b *testing.B) {
	benchWriteChain(b, true, database.LevelDB, 10000)
}
func BenchmarkChainWrite_full_10k_badgerDB(b *testing.B) {
	benchWriteChain(b, true, database.BadgerDB, 10000)
}

func BenchmarkChainWrite_header_100k_levelDB(b *testing.B) {
	benchWriteChain(b, false, database.LevelDB, 100000)
}
func BenchmarkChainWrite_header_100k_badgerDB(b *testing.B) {
	benchWriteChain(b, false, database.BadgerDB, 100000)
}

func BenchmarkChainWrite_full_100k_levelDB(b *testing.B) {
	benchWriteChain(b, true, database.LevelDB, 100000)
}
func BenchmarkChainWrite_full_100k_badgerDB(b *testing.B) {
	benchWriteChain(b, true, database.BadgerDB, 100000)
}

// Disabled because of too long test time
//func BenchmarkChainWrite_header_500k_levelDB(b *testing.B) {
//	benchWriteChain(b, false, database.LevelDB,500000)
//}
//func BenchmarkChainWrite_header_500k_badgerDB(b *testing.B) {
//	benchWriteChain(b, false, database.BadgerDB, 500000)
//}
//
//func BenchmarkChainWrite_full_500k_levelDB(b *testing.B) {
//	benchWriteChain(b, true, database.LevelDB,500000)
//}
//func BenchmarkChainWrite_full_500k_badgerDB(b *testing.B) {
//	benchWriteChain(b, true, database.BadgerDB,500000)
//}

// makeChainForBench writes a given number of headers or empty blocks/receipts
// into a database.
func makeChainForBench(db database.DBManager, full bool, count uint64) {
	var hash common.Hash
	for n := uint64(0); n < count; n++ {
		header := &types.Header{
			Number:      big.NewInt(int64(n)),
			ParentHash:  hash,
			BlockScore:  big.NewInt(1),
			TxHash:      types.EmptyRootHash,
			ReceiptHash: types.EmptyRootHash,
		}
		hash = header.Hash()

		db.WriteHeader(header)
		db.WriteCanonicalHash(hash, n)
		db.WriteTd(hash, n, big.NewInt(int64(n+1)))

		if full || n == 0 {
			db.WriteHeadBlockHash(hash)
			block := types.NewBlockWithHeader(header)
			db.WriteBody(hash, n, block.Body())
			db.WriteReceipts(hash, n, nil)
		}
	}
}

// write 'count' blocks to database 'b.N' times
func benchWriteChain(b *testing.B, full bool, databaseType database.DBType, count uint64) {
	for i := 0; i < b.N; i++ {
		dir := genTempDirForDB(b)

		db := genDBManagerForTest(dir, databaseType)
		makeChainForBench(db, full, count)

		db.Close()
		os.RemoveAll(dir)
	}
}

// write 'count' blocks to database and then read 'count' blocks
func benchReadChain(b *testing.B, full bool, databaseType database.DBType, count uint64) {
	dir := genTempDirForDB(b)
	defer os.RemoveAll(dir)

	db := genDBManagerForTest(dir, databaseType)
	makeChainForBench(db, full, count)
	db.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		db = genDBManagerForTest(dir, databaseType)

		chain, err := NewBlockChain(db, nil, params.TestChainConfig, gxhash.NewFaker(), vm.Config{})
		if err != nil {
			b.Fatalf("error creating chain: %v", err)
		}

		for n := uint64(0); n < count; n++ {
			header := chain.GetHeaderByNumber(n)
			if full {
				hash := header.Hash()
				db.ReadBody(hash, n)
				db.ReadReceipts(hash, n)
			}
		}
		chain.Stop()
		db.Close()
	}
}

// genTempDirForDB returns temp dir for database
func genTempDirForDB(b *testing.B) string {
	dir, err := ioutil.TempDir("", "klay-blockchain-bench")
	if err != nil {
		b.Fatalf("cannot create temporary directory: %v", err)
	}
	return dir
}

// genDBManagerForTest returns database.Database according to entered databaseType
func genDBManagerForTest(dir string, dbType database.DBType) database.DBManager {
	if dbType == database.MemoryDB {
		db := database.NewMemoryDBManager()
		return db
	} else {
		dbc := &database.DBConfig{Dir: dir, DBType: dbType, LevelDBCacheSize: 128, OpenFilesLimit: 128}
		return database.NewDBManager(dbc)
	}
}
