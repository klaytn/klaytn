// Modifications Copyright 2021 The klaytn Authors
// Copyright 2020 The go-ethereum Authors
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
// This file is derived from core/state/snapshot/generate_test.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/params"

	"github.com/klaytn/klaytn/log"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"golang.org/x/crypto/sha3"
)

func genExternallyOwnedAccount(nonce uint64, balance *big.Int) (account.Account, error) {
	return account.NewAccountWithMap(account.ExternallyOwnedAccountType, map[account.AccountValueKeyType]interface{}{
		account.AccountValueKeyNonce:         nonce,
		account.AccountValueKeyBalance:       balance,
		account.AccountValueKeyHumanReadable: false,
		account.AccountValueKeyAccountKey:    accountkey.NewAccountKeyLegacy(),
	})
}

func genSmartContractAccount(nonce uint64, balance *big.Int, storageRoot common.Hash, codeHash []byte) (account.Account, error) {
	return account.NewAccountWithMap(account.SmartContractAccountType, map[account.AccountValueKeyType]interface{}{
		account.AccountValueKeyNonce:         nonce,
		account.AccountValueKeyBalance:       balance,
		account.AccountValueKeyHumanReadable: false,
		account.AccountValueKeyAccountKey:    accountkey.NewAccountKeyLegacy(),
		account.AccountValueKeyStorageRoot:   storageRoot,
		account.AccountValueKeyCodeHash:      codeHash,
		account.AccountValueKeyCodeInfo:      params.CodeInfo(0),
	})
}

// Tests that snapshot generation from an empty database.
func TestGeneration(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	var (
		dbm    = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(dbm)
	)
	stTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	stTrie.Update([]byte("key-1"), []byte("val-1")) // 0x1314700b81afc49f94db3623ef1df38f3ed18b73a1b7ea2f6c095118cf6118a0
	stTrie.Update([]byte("key-2"), []byte("val-2")) // 0x18a0f4d79cff4459642dd7604f303886ad9d77c30cf3d7d7cedb3a693ab6d371
	stTrie.Update([]byte("key-3"), []byte("val-3")) // 0x51c71a47af0695957647fb68766d0becee77e953df17c29b3c2f25436f055c78
	stTrie.Commit(nil)                              // Root: 0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67

	accTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)

	acc1, _ := genSmartContractAccount(0, big.NewInt(1), stTrie.Hash(), emptyCode.Bytes())
	serializer := account.NewAccountSerializerWithAccount(acc1)
	val, _ := rlp.EncodeToBytes(serializer)
	accTrie.Update([]byte("acc-1"), val) // 0x30301e37c9af8ee5f609f1d60a3307d3e113bea03bef203e39aadc46bd5ad5ee

	acc2, _ := genExternallyOwnedAccount(0, big.NewInt(2))
	serializer2 := account.NewAccountSerializerWithAccount(acc2)
	val, _ = rlp.EncodeToBytes(serializer2)
	accTrie.Update([]byte("acc-2"), val) // 0x11944b79b322f047379973e18ba21657b8ad4b50ecd94217177bc56f89905228

	acc3, _ := genSmartContractAccount(0, big.NewInt(3), stTrie.Hash(), emptyCode.Bytes())
	serializer3 := account.NewAccountSerializerWithAccount(acc3)
	val, _ = rlp.EncodeToBytes(serializer3) // 0x8c2477df4801bbf88c6636445a2a9feff54c098cc218df403dc3f1007add780c
	accTrie.Update([]byte("acc-3"), val)

	root, _ := accTrie.Commit(nil) // Root: 0x4a651234bc4b8c7462b5ad4eb95bbb724eb636fed72bb5278d886f9ea4c345f8

	// TODO-Klaytn-Snapshot update proper block number
	triedb.Commit(root, false, 0)

	if have, want := root, common.HexToHash("0x4a651234bc4b8c7462b5ad4eb95bbb724eb636fed72bb5278d886f9ea4c345f8"); have != want {
		t.Fatalf("have %#x want %#x", have, want)
	}
	snap := generateSnapshot(dbm, triedb, 16, root)

	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}

	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

func hashData(input []byte) common.Hash {
	hasher := sha3.NewLegacyKeccak256()
	var hash common.Hash
	hasher.Reset()
	hasher.Write(input)
	hasher.Sum(hash[:0])
	return hash
}

// Tests that snapshot generation with existent flat state.
func TestGenerateExistentState(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	// We can't use statedb to make a test trie (circular dependency), so make
	// a fake one manually. We're going with a small account trie of 3 accounts,
	// two of which also has the same 3-slot storage trie attached.
	var (
		diskdb = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(diskdb)
	)
	stTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	stTrie.Update([]byte("key-1"), []byte("val-1")) // 0x1314700b81afc49f94db3623ef1df38f3ed18b73a1b7ea2f6c095118cf6118a0
	stTrie.Update([]byte("key-2"), []byte("val-2")) // 0x18a0f4d79cff4459642dd7604f303886ad9d77c30cf3d7d7cedb3a693ab6d371
	stTrie.Update([]byte("key-3"), []byte("val-3")) // 0x51c71a47af0695957647fb68766d0becee77e953df17c29b3c2f25436f055c78
	stTrie.Commit(nil)                              // Root: 0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67

	accTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	acc, _ := genSmartContractAccount(uint64(0), big.NewInt(1), stTrie.Hash(), emptyCode.Bytes())
	serializer := account.NewAccountSerializerWithAccount(acc)
	val, _ := rlp.EncodeToBytes(serializer)
	accTrie.Update([]byte("acc-1"), val) // 0x9250573b9c18c664139f3b6a7a8081b7d8f8916a8fcc5d94feec6c29f5fd4e9e
	diskdb.WriteAccountSnapshot(hashData([]byte("acc-1")), val)
	diskdb.WriteStorageSnapshot(hashData([]byte("acc-1")), hashData([]byte("key-1")), []byte("val-1"))
	diskdb.WriteStorageSnapshot(hashData([]byte("acc-1")), hashData([]byte("key-2")), []byte("val-2"))
	diskdb.WriteStorageSnapshot(hashData([]byte("acc-1")), hashData([]byte("key-3")), []byte("val-3"))

	acc, _ = genSmartContractAccount(uint64(0), big.NewInt(2), stTrie.Hash(), emptyCode.Bytes())
	serializer = account.NewAccountSerializerWithAccount(acc)
	val, _ = rlp.EncodeToBytes(serializer)
	accTrie.Update([]byte("acc-2"), val) // 0x65145f923027566669a1ae5ccac66f945b55ff6eaeb17d2ea8e048b7d381f2d7
	// diskdb.Put(hashData([]byte("acc-2")).Bytes(), val)
	diskdb.WriteAccountSnapshot(hashData([]byte("acc-2")), val)

	acc, _ = genSmartContractAccount(uint64(0), big.NewInt(3), stTrie.Hash(), emptyCode.Bytes())
	serializer = account.NewAccountSerializerWithAccount(acc)
	val, _ = rlp.EncodeToBytes(serializer)
	accTrie.Update([]byte("acc-3"), val) // 0x50815097425d000edfc8b3a4a13e175fc2bdcfee8bdfbf2d1ff61041d3c235b2
	diskdb.WriteAccountSnapshot(hashData([]byte("acc-3")), val)
	diskdb.WriteStorageSnapshot(hashData([]byte("acc-3")), hashData([]byte("key-1")), []byte("val-1"))
	diskdb.WriteStorageSnapshot(hashData([]byte("acc-3")), hashData([]byte("key-2")), []byte("val-2"))
	diskdb.WriteStorageSnapshot(hashData([]byte("acc-3")), hashData([]byte("key-3")), []byte("val-3"))

	root, _ := accTrie.Commit(nil) // Root: 0xe3712f1a226f3782caca78ca770ccc19ee000552813a9f59d479f8611db9b1fd
	// TODO-Klaytn-Snapshot put proper block number
	triedb.Commit(root, false, 0)

	snap := generateSnapshot(diskdb, triedb, 16, root)
	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

func checkSnapRoot(t *testing.T, snap *diskLayer, trieRoot common.Hash) {
	t.Helper()
	accIt := snap.AccountIterator(common.Hash{})
	defer accIt.Release()
	snapRoot, err := generateTrieRoot(accIt, common.Hash{}, trieGenerate,
		func(accountHash, codeHash common.Hash, stat *generateStats) (common.Hash, error) {
			storageIt, _ := snap.StorageIterator(accountHash, common.Hash{})
			defer storageIt.Release()

			hash, err := generateTrieRoot(storageIt, accountHash, trieGenerate, nil, stat, false)
			if err != nil {
				return common.Hash{}, err
			}
			return hash, nil
		}, newGenerateStats(), true)
	if err != nil {
		t.Fatal(err)
	}
	if snapRoot != trieRoot {
		t.Fatalf("snaproot: %#x != trieroot %#x", snapRoot, trieRoot)
	}
}

type testHelper struct {
	diskdb  database.DBManager
	triedb  *statedb.Database
	accTrie *statedb.SecureTrie
}

func newHelper() *testHelper {
	diskdb := database.NewMemoryDBManager()
	triedb := statedb.NewDatabase(diskdb)
	accTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	return &testHelper{
		diskdb:  diskdb,
		triedb:  triedb,
		accTrie: accTrie,
	}
}

func (t *testHelper) addTrieAccount(acckey string, acc account.Account) {
	val, _ := rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
	t.accTrie.Update([]byte(acckey), val)
}

func (t *testHelper) addSnapAccount(acckey string, acc account.Account) {
	val, _ := rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
	key := hashData([]byte(acckey))
	t.diskdb.WriteAccountSnapshot(key, val)
}

func (t *testHelper) addAccount(acckey string, acc account.Account) {
	t.addTrieAccount(acckey, acc)
	t.addSnapAccount(acckey, acc)
}

func (t *testHelper) addSnapStorage(accKey string, keys []string, vals []string) {
	accHash := hashData([]byte(accKey))
	for i, key := range keys {
		t.diskdb.WriteStorageSnapshot(accHash, hashData([]byte(key)), []byte(vals[i]))
	}
}

func (t *testHelper) makeStorageTrie(keys []string, vals []string) common.Hash {
	stTrie, _ := statedb.NewSecureTrie(common.Hash{}, t.triedb)
	for i, k := range keys {
		stTrie.Update([]byte(k), []byte(vals[i]))
	}
	root, _ := stTrie.Commit(nil)
	return root
}

func (t *testHelper) Generate() (common.Hash, *diskLayer) {
	root, _ := t.accTrie.Commit(nil)
	// TODO-Klaytn-Snapshot input proper block number
	t.triedb.Commit(root, false, 0)
	snap := generateSnapshot(t.diskdb, t.triedb, 16, root)
	return root, snap
}

// Tests that snapshot generation with existent flat state, where the flat state
// contains some errors:
// - the contract with empty storage root but has storage entries in the disk
// - the contract with non empty storage root but empty storage slots
// - the contract(non-empty storage) misses some storage slots
//   - miss in the beginning
//   - miss in the middle
//   - miss in the end
// - the contract(non-empty storage) has wrong storage slots
//   - wrong slots in the beginning
//   - wrong slots in the middle
//   - wrong slots in the end
// - the contract(non-empty storage) has extra storage slots
//   - extra slots in the beginning
//   - extra slots in the middle
//   - extra slots in the end
func TestGenerateExistentStateWithWrongStorage(t *testing.T) {
	helper := newHelper()
	stRoot := helper.makeStorageTrie([]string{"key-1", "key-2", "key-3"}, []string{"val-1", "val-2", "val-3"})

	// Account one, empty root but non-empty database
	acc1, _ := genExternallyOwnedAccount(0, big.NewInt(1))
	helper.addAccount("acc-1", acc1)
	helper.addSnapStorage("acc-1", []string{"key-1", "key-2", "key-3"}, []string{"val-1", "val-2", "val-3"})

	// Account two, non empty root but empty database
	acc2, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
	helper.addAccount("acc-2", acc2)

	// Miss slots
	{
		// Account three, non empty root but misses slots in the beginning
		acc3, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-3", acc3)
		helper.addSnapStorage("acc-3", []string{"key-2", "key-3"}, []string{"val-2", "val-3"})

		// Account four, non empty root but misses slots in the middle
		acc4, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-4", acc4)
		helper.addSnapStorage("acc-4", []string{"key-1", "key-3"}, []string{"val-1", "val-3"})

		// Account five, non empty root but misses slots in the end
		acc5, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-5", acc5)
		helper.addSnapStorage("acc-5", []string{"key-1", "key-2"}, []string{"val-1", "val-2"})
	}

	// Wrong storage slots
	{
		// Account six, non empty root but wrong slots in the beginning

		acc6, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-6", acc6)
		helper.addSnapStorage("acc-6", []string{"key-1", "key-2", "key-3"}, []string{"badval-1", "val-2", "val-3"})

		// Account seven, non empty root but wrong slots in the middle
		acc7, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-7", acc7)
		helper.addSnapStorage("acc-7", []string{"key-1", "key-2", "key-3"}, []string{"val-1", "badval-2", "val-3"})

		// Account eight, non empty root but wrong slots in the end
		acc8, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-8", acc8)
		helper.addSnapStorage("acc-8", []string{"key-1", "key-2", "key-3"}, []string{"val-1", "val-2", "badval-3"})

		// Account 9, non empty root but rotated slots
		acc9, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-9", acc9)
		helper.addSnapStorage("acc-9", []string{"key-1", "key-2", "key-3"}, []string{"val-1", "val-3", "val-2"})
	}

	// Extra storage slots
	{
		// Account 10, non empty root but extra slots in the beginning
		acc10, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-10", acc10)
		helper.addSnapStorage("acc-10", []string{"key-0", "key-1", "key-2", "key-3"}, []string{"val-0", "val-1", "val-2", "val-3"})

		// Account 11, non empty root but extra slots in the middle
		acc11, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-11", acc11)
		helper.addSnapStorage("acc-11", []string{"key-1", "key-2", "key-2-1", "key-3"}, []string{"val-1", "val-2", "val-2-1", "val-3"})

		// Account 12, non empty root but extra slots in the end
		acc12, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount("acc-12", acc12)
		helper.addSnapStorage("acc-12", []string{"key-1", "key-2", "key-3", "key-4"}, []string{"val-1", "val-2", "val-3", "val-4"})
	}

	root, snap := helper.Generate()
	t.Logf("Root: %#x\n", root) // Root = 0x8746cce9fd9c658b2cfd639878ed6584b7a2b3e73bb40f607fcfa156002429a0

	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

// Tests that snapshot generation with existent flat state, where the flat state
// contains some errors:
// - miss accounts
// - wrong accounts
// - extra accounts
func TestGenerateExistentStateWithWrongAccounts(t *testing.T) {
	helper := newHelper()
	stRoot := helper.makeStorageTrie([]string{"key-1", "key-2", "key-3"}, []string{"val-1", "val-2", "val-3"})

	// Trie accounts [acc-1, acc-2, acc-3, acc-4, acc-6]
	// Extra accounts [acc-0, acc-5, acc-7]

	// Missing accounts, only in the trie
	{

		acc1, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addTrieAccount("acc-1", acc1) // Beginning
		acc4, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addTrieAccount("acc-4", acc4) // Middle
		acc6, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addTrieAccount("acc-6", acc6) // End
	}

	// Wrong accounts
	{
		acc2, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addTrieAccount("acc-2", acc2)
		acc2, _ = genSmartContractAccount(0, big.NewInt(1), stRoot, common.Hex2Bytes("0x1234"))
		helper.addSnapAccount("acc-2", acc2)

		acc3, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addTrieAccount("acc-3", acc3)
		acc3, _ = genExternallyOwnedAccount(0, big.NewInt(1))
		helper.addSnapAccount("acc-3", acc3)
	}

	// Extra accounts, only in the snap
	{
		acc0, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyRoot.Bytes())
		helper.addSnapAccount("acc-0", acc0) // before the beginning
		acc5, _ := genSmartContractAccount(0, big.NewInt(1), emptyRoot, common.Hex2Bytes("0x1234"))
		helper.addSnapAccount("acc-5", acc5) // Middle
		acc7, _ := genSmartContractAccount(0, big.NewInt(1), emptyRoot, emptyRoot.Bytes())
		helper.addSnapAccount("acc-7", acc7) // after the end
	}

	root, snap := helper.Generate()
	t.Logf("Root: %#x\n", root) // Root = 0x825891472281463511e7ebcc7f109e4f9200c20fa384754e11fd605cd98464e8

	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)

	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

// Tests that snapshot generation errors out correctly in case of a missing trie
// node in the account trie.
func TestGenerateCorruptAccountTrie(t *testing.T) {
	// We can't use statedb to make a test trie (circular dependency), so make
	// a fake one manually. We're going with a small account trie of 3 accounts,
	// without any storage slots to keep the test smaller.
	var (
		diskdb = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(diskdb)
	)
	tr, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	acc1, _ := genExternallyOwnedAccount(0, big.NewInt(1))
	serializer1 := account.NewAccountSerializerWithAccount(acc1)
	val, _ := rlp.EncodeToBytes(serializer1)
	tr.Update([]byte("acc-1"), val) // 0xc7a30f39aff471c95d8a837497ad0e49b65be475cc0953540f80cfcdbdcd9074

	acc2, _ := genExternallyOwnedAccount(0, big.NewInt(2))
	serializer2 := account.NewAccountSerializerWithAccount(acc2)
	val, _ = rlp.EncodeToBytes(serializer2)
	tr.Update([]byte("acc-2"), val) // 0x65145f923027566669a1ae5ccac66f945b55ff6eaeb17d2ea8e048b7d381f2d7

	acc3, _ := genExternallyOwnedAccount(0, big.NewInt(3))
	serializer3 := account.NewAccountSerializerWithAccount(acc3)
	val, _ = rlp.EncodeToBytes(serializer3)
	tr.Update([]byte("acc-3"), val) // 0x19ead688e907b0fab07176120dceec244a72aff2f0aa51e8b827584e378772f4
	tr.Commit(nil)                  // Root: 0xa04693ea110a31037fb5ee814308a6f1d76bdab0b11676bdf4541d2de55ba978

	// Delete an account trie leaf and ensure the generator chokes
	// TODO-Klaytn-Snapshot put propoer block number
	triedb.Commit(common.HexToHash("0xa04693ea110a31037fb5ee814308a6f1d76bdab0b11676bdf4541d2de55ba978"), false, 0)
	diskdb.GetMemDB().Delete(common.HexToHash("0x65145f923027566669a1ae5ccac66f945b55ff6eaeb17d2ea8e048b7d381f2d7").Bytes())

	snap := generateSnapshot(diskdb, triedb, 16, common.HexToHash("0xa04693ea110a31037fb5ee814308a6f1d76bdab0b11676bdf4541d2de55ba978"))
	select {
	case <-snap.genPending:
		// Snapshot generation succeeded
		t.Errorf("Snapshot generated against corrupt account trie")

	case <-time.After(time.Second):
		// Not generated fast enough, hopefully blocked inside on missing trie node fail
	}
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

// Tests that snapshot generation errors out correctly in case of a missing root
// trie node for a storage trie. It's similar to internal corruption but it is
// handled differently inside the generator.
func TestGenerateMissingStorageTrie(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	// We can't use statedb to make a test trie (circular dependency), so make
	// a fake one manually. We're going with a small account trie of 3 accounts,
	// two of which also has the same 3-slot storage trie attached.
	var (
		diskdb = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(diskdb)
	)
	stTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	stTrie.Update([]byte("key-1"), []byte("val-1")) // 0x1314700b81afc49f94db3623ef1df38f3ed18b73a1b7ea2f6c095118cf6118a0
	stTrie.Update([]byte("key-2"), []byte("val-2")) // 0x18a0f4d79cff4459642dd7604f303886ad9d77c30cf3d7d7cedb3a693ab6d371
	stTrie.Update([]byte("key-3"), []byte("val-3")) // 0x51c71a47af0695957647fb68766d0becee77e953df17c29b3c2f25436f055c78
	stTrie.Commit(nil)                              // Root: 0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67

	accTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	acc, _ := genSmartContractAccount(0, big.NewInt(1), stTrie.Hash(), emptyCode.Bytes())
	val, _ := rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
	accTrie.Update([]byte("acc-1"), val) // 0x30301e37c9af8ee5f609f1d60a3307d3e113bea03bef203e39aadc46bd5ad5ee

	acc, _ = genSmartContractAccount(0, big.NewInt(2), stTrie.Hash(), emptyCode.Bytes())
	val, _ = rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
	accTrie.Update([]byte("acc-2"), val) // 0xff964e27aca3ef5766a7709de6beff44863d539c9af88ab1719865b80a55b6b2
	t.Log(accTrie.Hash().String())

	acc, _ = genSmartContractAccount(0, big.NewInt(3), stTrie.Hash(), emptyCode.Bytes())
	val, _ = rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
	accTrie.Update([]byte("acc-3"), val) // 0x8c2477df4801bbf88c6636445a2a9feff54c098cc218df403dc3f1007add780c
	accTrie.Commit(nil)                  // Root: 0xa2282b99de1fc11e32d26bee37707ef49a6978b2d375796a1b026a497193a2ef

	// We can only corrupt the disk database, so flush the tries out
	triedb.Reference(
		common.HexToHash("0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67"),
		common.HexToHash("0x30301e37c9af8ee5f609f1d60a3307d3e113bea03bef203e39aadc46bd5ad5ee"),
	)
	triedb.Reference(
		common.HexToHash("0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67"),
		common.HexToHash("0x8c2477df4801bbf88c6636445a2a9feff54c098cc218df403dc3f1007add780c"),
	)
	// TODO-Klaytn-Snapshot put proper block number
	triedb.Commit(common.HexToHash("0xa2282b99de1fc11e32d26bee37707ef49a6978b2d375796a1b026a497193a2ef"), false, 0)

	// Delete a storage trie root and ensure the generator chokes
	diskdb.GetMemDB().Delete(common.HexToHash("0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67").Bytes())

	snap := generateSnapshot(diskdb, triedb, 16, common.HexToHash("0xa2282b99de1fc11e32d26bee37707ef49a6978b2d375796a1b026a497193a2ef"))
	select {
	case <-snap.genPending:
		// Snapshot generation succeeded
		t.Errorf("Snapshot generated against corrupt storage trie")

	case <-time.After(time.Second):
		// Not generated fast enough, hopefully blocked inside on missing trie node fail
	}
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

// Tests that snapshot generation errors out correctly in case of a missing trie
// node in a storage trie.
func TestGenerateCorruptStorageTrie(t *testing.T) {
	// We can't use statedb to make a test trie (circular dependency), so make
	// a fake one manually. We're going with a small account trie of 3 accounts,
	// two of which also has the same 3-slot storage trie attached.
	var (
		diskdb = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(diskdb)
	)
	stTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	stTrie.Update([]byte("key-1"), []byte("val-1")) // 0x1314700b81afc49f94db3623ef1df38f3ed18b73a1b7ea2f6c095118cf6118a0
	stTrie.Update([]byte("key-2"), []byte("val-2")) // 0x18a0f4d79cff4459642dd7604f303886ad9d77c30cf3d7d7cedb3a693ab6d371
	stTrie.Update([]byte("key-3"), []byte("val-3")) // 0x51c71a47af0695957647fb68766d0becee77e953df17c29b3c2f25436f055c78
	stTrie.Commit(nil)                              // Root: 0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67

	accTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	acc, _ := genSmartContractAccount(0, big.NewInt(1), stTrie.Hash(), emptyCode.Bytes())
	serializer := account.NewAccountSerializerWithAccount(acc)
	val, _ := rlp.EncodeToBytes(serializer)
	accTrie.Update([]byte("acc-1"), val) // 0x30301e37c9af8ee5f609f1d60a3307d3e113bea03bef203e39aadc46bd5ad5ee

	acc, _ = genExternallyOwnedAccount(0, big.NewInt(2))
	serializer = account.NewAccountSerializerWithAccount(acc)
	val, _ = rlp.EncodeToBytes(serializer)
	accTrie.Update([]byte("acc-2"), val) // 0x11944b79b322f047379973e18ba21657b8ad4b50ecd94217177bc56f89905228

	acc, _ = genSmartContractAccount(0, big.NewInt(3), stTrie.Hash(), emptyCode.Bytes())
	serializer = account.NewAccountSerializerWithAccount(acc)
	val, _ = rlp.EncodeToBytes(serializer)
	accTrie.Update([]byte("acc-3"), val) // 0x8c2477df4801bbf88c6636445a2a9feff54c098cc218df403dc3f1007add780c
	accTrie.Commit(nil)                  // Root: 0x4a651234bc4b8c7462b5ad4eb95bbb724eb636fed72bb5278d886f9ea4c345f8

	// We can only corrupt the disk database, so flush the tries out
	triedb.Reference(
		common.HexToHash("0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67"),
		common.HexToHash("0x30301e37c9af8ee5f609f1d60a3307d3e113bea03bef203e39aadc46bd5ad5ee"),
	)
	triedb.Reference(
		common.HexToHash("0xddefcd9376dd029653ef384bd2f0a126bb755fe84fdcc9e7cf421ba454f2bc67"),
		common.HexToHash("0x8c2477df4801bbf88c6636445a2a9feff54c098cc218df403dc3f1007add780c"),
	)
	// TODO-Klaytn-Snapshot put proper block number
	triedb.Commit(common.HexToHash("0x4a651234bc4b8c7462b5ad4eb95bbb724eb636fed72bb5278d886f9ea4c345f8"), false, 0)

	// Delete a storage trie leaf and ensure the generator chokes
	diskdb.GetMemDB().Delete(common.HexToHash("0x18a0f4d79cff4459642dd7604f303886ad9d77c30cf3d7d7cedb3a693ab6d371").Bytes())

	snap := generateSnapshot(diskdb, triedb, 16, common.HexToHash("0x4a651234bc4b8c7462b5ad4eb95bbb724eb636fed72bb5278d886f9ea4c345f8"))
	select {
	case <-snap.genPending:
		// Snapshot generation succeeded
		t.Errorf("Snapshot generated against corrupt storage trie")

	case <-time.After(time.Second):
		// Not generated fast enough, hopefully blocked inside on missing trie node fail
	}
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

func getStorageTrie(n int, triedb *statedb.Database) *statedb.SecureTrie {
	stTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	for i := 0; i < n; i++ {
		k := fmt.Sprintf("key-%d", i)
		v := fmt.Sprintf("val-%d", i)
		stTrie.Update([]byte(k), []byte(v))
	}
	stTrie.Commit(nil)
	return stTrie
}

// Tests that snapshot generation when an extra account with storage exists in the snap state.
func TestGenerateWithExtraAccounts(t *testing.T) {
	var (
		diskdb = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(diskdb)
		stTrie = getStorageTrie(5, triedb)
	)
	accTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	{ // Account one in the trie
		acc, _ := genSmartContractAccount(0, big.NewInt(1), stTrie.Hash(), emptyCode.Bytes())
		val, _ := rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
		accTrie.Update([]byte("acc-1"), val) // 0xfa43eb0210d32b0013ae744d26ded52489ee3cab4a5bd9128a599135aba7c088
		// Identical in the snap
		key := hashData([]byte("acc-1"))
		diskdb.WriteAccountSnapshot(key, val)
		diskdb.WriteStorageSnapshot(key, hashData([]byte("key-1")), []byte("val-1"))
		diskdb.WriteStorageSnapshot(key, hashData([]byte("key-2")), []byte("val-2"))
		diskdb.WriteStorageSnapshot(key, hashData([]byte("key-3")), []byte("val-3"))
		diskdb.WriteStorageSnapshot(key, hashData([]byte("key-4")), []byte("val-4"))
		diskdb.WriteStorageSnapshot(key, hashData([]byte("key-5")), []byte("val-5"))
	}
	{ // Account two exists only in the snapshot
		acc, _ := genSmartContractAccount(0, big.NewInt(1), stTrie.Hash(), emptyCode.Bytes())
		val, _ := rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
		key := hashData([]byte("acc-2"))
		diskdb.WriteAccountSnapshot(key, val)
		diskdb.WriteStorageSnapshot(key, hashData([]byte("b-key-1")), []byte("b-val-1"))
		diskdb.WriteStorageSnapshot(key, hashData([]byte("b-key-2")), []byte("b-val-2"))
		diskdb.WriteStorageSnapshot(key, hashData([]byte("b-key-3")), []byte("b-val-3"))
	}
	root, _ := accTrie.Commit(nil)
	t.Logf("root: %x", root)
	// TODO-Klaytn-Snapshot put proper block number
	triedb.Commit(root, false, 0)
	// To verify the test: If we now inspect the snap db, there should exist extraneous storage items
	if data := diskdb.ReadStorageSnapshot(hashData([]byte("acc-2")), hashData([]byte("b-key-1"))); data == nil {
		t.Fatalf("expected snap storage to exist")
	}

	snap := generateSnapshot(diskdb, triedb, 16, root)
	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
	// If we now inspect the snap db, there should exist no extraneous storage items
	if data := diskdb.ReadStorageSnapshot(hashData([]byte("acc-2")), hashData([]byte("b-key-1"))); data != nil {
		t.Fatalf("expected slot to be removed, got %v", string(data))
	}
}

// Tests that snapshot generation when an extra account with storage exists in the snap state.
func TestGenerateWithManyExtraAccounts(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)

	var (
		diskdb = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(diskdb)
		stTrie = getStorageTrie(3, triedb)
	)
	accTrie, _ := statedb.NewSecureTrie(common.Hash{}, triedb)
	{ // Account one in the trie
		acc, _ := genSmartContractAccount(0, big.NewInt(1), stTrie.Hash(), emptyCode.Bytes())
		val, _ := rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
		accTrie.Update([]byte("acc-1"), val) // 0xe9ddf05eaf05dbd7c6fb1137e90b4f3b1ed43383aaaa88312de297e304599966
		// Identical in the snap
		key := hashData([]byte("acc-1"))
		diskdb.WriteAccountSnapshot(key, val)
		diskdb.WriteStorageSnapshot(key, hashData([]byte("key-1")), []byte("val-1"))
		diskdb.WriteStorageSnapshot(key, hashData([]byte("key-2")), []byte("val-2"))
		diskdb.WriteStorageSnapshot(key, hashData([]byte("key-3")), []byte("val-3"))
	}
	{ // 100 accounts exist only in snapshot
		for i := 0; i < 1000; i++ {
			acc, _ := genExternallyOwnedAccount(uint64(i), big.NewInt(int64(i)))
			val, _ := rlp.EncodeToBytes(acc)
			key := hashData([]byte(fmt.Sprintf("acc-%d", i)))
			diskdb.WriteAccountSnapshot(key, val)
		}
	}
	root, _ := accTrie.Commit(nil)
	t.Logf("root: %x", root)
	// TODO-Klaytn-Snapshot put proper block number
	triedb.Commit(root, false, 0)

	snap := generateSnapshot(diskdb, triedb, 16, root)
	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

// Tests this case
// maxAccountRange 3
// snapshot-accounts: 01, 02, 03, 04, 05, 06, 07
// trie-accounts:             03,             07
//
// We iterate three snapshot storage slots (max = 3) from the database. They are 0x01, 0x02, 0x03.
// The trie has a lot of deletions.
// So in trie, we iterate 2 entries 0x03, 0x07. We create the 0x07 in the database and abort the procedure, because the trie is exhausted.
// But in the database, we still have the stale storage slots 0x04, 0x05. They are not iterated yet, but the procedure is finished.
func TestGenerateWithExtraBeforeAndAfter(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	accountCheckRange = 3
	var (
		diskdb = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(diskdb)
	)
	accTrie, _ := statedb.NewTrie(common.Hash{}, triedb)
	{
		acc, _ := genExternallyOwnedAccount(0, big.NewInt(1))
		val, _ := rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
		accTrie.Update(common.HexToHash("0x03").Bytes(), val)
		accTrie.Update(common.HexToHash("0x07").Bytes(), val)

		diskdb.WriteAccountSnapshot(common.HexToHash("0x01"), val)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x02"), val)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x03"), val)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x04"), val)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x05"), val)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x06"), val)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x07"), val)
	}

	root, _ := accTrie.Commit(nil)
	t.Logf("root: %x", root)
	// TODO-Klaytn-Snapshot put proper block number
	triedb.Commit(root, false, 0)

	snap := generateSnapshot(diskdb, triedb, 16, root)
	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

// TestGenerateWithMalformedSnapdata tests what happes if we have some junk
// in the snapshot database, which cannot be parsed back to an account
func TestGenerateWithMalformedSnapdata(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	accountCheckRange = 3
	var (
		diskdb = database.NewMemoryDBManager()
		triedb = statedb.NewDatabase(diskdb)
	)
	accTrie, _ := statedb.NewTrie(common.Hash{}, triedb)
	{
		acc, _ := genExternallyOwnedAccount(0, big.NewInt(1))
		val, _ := rlp.EncodeToBytes(account.NewAccountSerializerWithAccount(acc))
		accTrie.Update(common.HexToHash("0x03").Bytes(), val)

		junk := make([]byte, 100)
		copy(junk, []byte{0xde, 0xad})
		diskdb.WriteAccountSnapshot(common.HexToHash("0x02"), junk)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x03"), junk)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x04"), junk)
		diskdb.WriteAccountSnapshot(common.HexToHash("0x05"), junk)
	}

	root, _ := accTrie.Commit(nil)
	t.Logf("root: %x", root)
	// TODO-Klaytn-Snapshot put proper block number
	triedb.Commit(root, false, 0)

	snap := generateSnapshot(diskdb, triedb, 16, root)
	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
	// If we now inspect the snap db, there should exist no extraneous storage items
	if data := diskdb.ReadStorageSnapshot(hashData([]byte("acc-2")), hashData([]byte("b-key-1"))); data != nil {
		t.Fatalf("expected slot to be removed, got %v", string(data))
	}
}

func TestGenerateFromEmptySnap(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	accountCheckRange = 10
	storageCheckRange = 20
	helper := newHelper()
	stRoot := helper.makeStorageTrie([]string{"key-1", "key-2", "key-3"}, []string{"val-1", "val-2", "val-3"})
	// Add 1K accounts to the trie
	for i := 0; i < 400; i++ {
		acc, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addTrieAccount(fmt.Sprintf("acc-%d", i), acc)
	}
	root, snap := helper.Generate()
	t.Logf("Root: %#x\n", root) // Root: 0x6f7af6d2e1a1bf2b84a3beb3f8b64388465fbc1e274ca5d5d3fc787ca78f59e4

	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}

// Tests that snapshot generation with existent flat state, where the flat state
// storage is correct, but incomplete.
// The incomplete part is on the second range
// snap: [ 0x01, 0x02, 0x03, 0x04] , [ 0x05, 0x06, 0x07, {missing}] (with storageCheck = 4)
// trie:  0x01, 0x02, 0x03, 0x04,  0x05, 0x06, 0x07, 0x08
// This hits a case where the snap verification passes, but there are more elements in the trie
// which we must also add.
func TestGenerateWithIncompleteStorage(t *testing.T) {
	storageCheckRange = 4
	helper := newHelper()
	stKeys := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	stVals := []string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8"}
	stRoot := helper.makeStorageTrie(stKeys, stVals)
	// We add 8 accounts, each one is missing exactly one of the storage slots. This means
	// we don't have to order the keys and figure out exactly which hash-key winds up
	// on the sensitive spots at the boundaries
	for i := 0; i < 8; i++ {
		accKey := fmt.Sprintf("acc-%d", i)
		acc, _ := genSmartContractAccount(0, big.NewInt(1), stRoot, emptyCode.Bytes())
		helper.addAccount(accKey, acc)
		var moddedKeys []string
		var moddedVals []string
		for ii := 0; ii < 8; ii++ {
			if ii != i {
				moddedKeys = append(moddedKeys, stKeys[ii])
				moddedVals = append(moddedVals, stVals[ii])
			}
		}
		helper.addSnapStorage(accKey, moddedKeys, moddedVals)
	}

	root, snap := helper.Generate()
	t.Logf("Root: %#x\n", root) // Root: 0xca73f6f05ba4ca3024ef340ef3dfca8fdabc1b677ff13f5a9571fd49c16e67ff

	select {
	case <-snap.genPending:
		// Snapshot generation succeeded

	case <-time.After(3 * time.Second):
		t.Errorf("Snapshot generation failed")
	}
	checkSnapRoot(t, snap, root)
	// Signal abortion to the generator and wait for it to tear down
	stop := make(chan *generatorStats)
	snap.genAbort <- stop
	<-stop
}
