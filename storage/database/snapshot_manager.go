package database

import (
	"encoding/binary"

	"github.com/klaytn/klaytn/common"
)

type SnapshotManager interface {
	ReadSnapshotRoot() common.Hash
	WriteSnapshotRoot(root common.Hash)
	DeleteSnapshotRoot()

	ReadAccountSnapshot(hash common.Hash) []byte
	WriteAccountSnapshot(hash common.Hash, entry []byte)
	DeleteAccountSnapshot(hash common.Hash)

	WriteStorageSnapshot(accountHash, storageHash common.Hash, entry []byte)
	ReadStorageSnapshot(accountHash, storageHash common.Hash) []byte
	DeleteStorageSnapshot(accountHash, storageHash common.Hash)

	ReadSnapshotJournal() []byte
	WriteSnapshotJournal(journal []byte)
	DeleteSnapshotJournal()

	ReadSnapshotGenerator() []byte
	WriteSnapshotGenerator(generator []byte)
	DeleteSnapshotGenerator()

	ReadSnapshotRecoveryNumber() *uint64
	WriteSnapshotRecoveryNumber(number uint64)
	DeleteSnapshotRecoveryNumber()

	IterateStorageSnapshots(accountHash common.Hash) Iterator

	NewSnapshotDBBatch() SnapshotDBBatch
}

func (dbm *databaseManager) NewSnapshotDBBatch() SnapshotDBBatch {
	return &snapshotDBBatch{dbm.NewBatch(SnapshotDB)}
}

type SnapshotDBBatch interface {
	Batch
	WriteSnapshotRoot(root common.Hash)
	DeleteSnapshotRoot()

	WriteAccountSnapshot(hash common.Hash, entry []byte)
	DeleteAccountSnapshot(hash common.Hash)

	WriteStorageSnapshot(accountHash, storageHash common.Hash, entry []byte)
	DeleteStorageSnapshot(accountHash, storageHash common.Hash)
}

type snapshotDBBatch struct {
	Batch
}

func (batch *snapshotDBBatch) WriteSnapshotRoot(root common.Hash) {
	writeSnapshotRoot(batch, root)
}

func (batch *snapshotDBBatch) DeleteSnapshotRoot() {
	deleteSnapshotRoot(batch)
}

func (batch *snapshotDBBatch) WriteAccountSnapshot(hash common.Hash, entry []byte) {
	writeAccountSnapshot(batch, hash, entry)
}
func (batch *snapshotDBBatch) DeleteAccountSnapshot(hash common.Hash) {
	deleteAccountSnapshot(batch, hash)
}

func (batch *snapshotDBBatch) WriteStorageSnapshot(accountHash, storageHash common.Hash, entry []byte) {
	writeStorageSnapshot(batch, accountHash, storageHash, entry)
}

func (batch *snapshotDBBatch) DeleteStorageSnapshot(accountHash, storageHash common.Hash) {
	deleteStorageSnapshot(batch, accountHash, storageHash)
}

func writeSnapshotRoot(db KeyValueWriter, root common.Hash) {
	if err := db.Put(snapshotRootKey, root[:]); err != nil {
		logger.Crit("Failed to store snapshot root", "err", err)
	}
}

func deleteSnapshotRoot(db KeyValueWriter) {
	if err := db.Delete(snapshotRootKey); err != nil {
		logger.Crit("Failed to remove snapshot root", "err", err)
	}
}

func writeAccountSnapshot(db KeyValueWriter, hash common.Hash, entry []byte) {
	if err := db.Put(accountSnapshotKey(hash), entry); err != nil {
		logger.Crit("Failed to store account snapshot", "err", err)
	}
}

func deleteAccountSnapshot(db KeyValueWriter, hash common.Hash) {
	if err := db.Delete(accountSnapshotKey(hash)); err != nil {
		logger.Crit("Failed to delete account snapshot", "err", err)
	}
}

func writeStorageSnapshot(db KeyValueWriter, accountHash, storageHash common.Hash, entry []byte) {
	if err := db.Put(storageSnapshotKey(accountHash, storageHash), entry); err != nil {
		logger.Crit("Failed to store storage snapshot", "err", err)
	}
}

func deleteStorageSnapshot(db KeyValueWriter, accountHash, storageHash common.Hash) {
	if err := db.Delete(storageSnapshotKey(accountHash, storageHash)); err != nil {
		logger.Crit("Failed to delete storage snapshot", "err", err)
	}
}

// ReadSnapshotRoot retrieves the root of the block whose state is contained in
// the persisted snapshot.
func (dbm *databaseManager) ReadSnapshotRoot() common.Hash {
	db := dbm.GetDatabase(SnapshotDB)
	data, _ := db.Get(snapshotRootKey)
	if len(data) != common.HashLength {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteSnapshotRoot stores the root of the block whose state is contained in
// the persisted snapshot.
func (dbm *databaseManager) WriteSnapshotRoot(root common.Hash) {
	db := dbm.GetDatabase(SnapshotDB)
	writeSnapshotRoot(db, root)
}

// DeleteSnapshotRoot deletes the hash of the block whose state is contained in
// the persisted snapshot. Since snapshots are not immutable, this  method can
// be used during updates, so a crash or failure will mark the entire snapshot
// invalid.
func (dbm *databaseManager) DeleteSnapshotRoot() {
	db := dbm.GetDatabase(SnapshotDB)
	deleteSnapshotRoot(db)
}

// ReadAccountSnapshot retrieves the snapshot entry of an account trie leaf.
func (dbm *databaseManager) ReadAccountSnapshot(hash common.Hash) []byte {
	db := dbm.GetDatabase(SnapshotDB)
	data, _ := db.Get(accountSnapshotKey(hash))
	return data
}

// WriteAccountSnapshot stores the snapshot entry of an account trie leaf.
func (dbm *databaseManager) WriteAccountSnapshot(hash common.Hash, entry []byte) {
	db := dbm.GetDatabase(SnapshotDB)
	writeAccountSnapshot(db, hash, entry)
}

// DeleteAccountSnapshot removes the snapshot entry of an account trie leaf.
func (dbm *databaseManager) DeleteAccountSnapshot(hash common.Hash) {
	db := dbm.GetDatabase(SnapshotDB)
	deleteAccountSnapshot(db, hash)
}

// ReadStorageSnapshot retrieves the snapshot entry of an storage trie leaf.
func (dbm *databaseManager) ReadStorageSnapshot(accountHash, storageHash common.Hash) []byte {
	db := dbm.GetDatabase(SnapshotDB)
	data, _ := db.Get(storageSnapshotKey(accountHash, storageHash))
	return data
}

// WriteStorageSnapshot stores the snapshot entry of an storage trie leaf.
func (dbm *databaseManager) WriteStorageSnapshot(accountHash, storageHash common.Hash, entry []byte) {
	db := dbm.GetDatabase(SnapshotDB)
	writeStorageSnapshot(db, accountHash, storageHash, entry)

}

// DeleteStorageSnapshot removes the snapshot entry of an storage trie leaf.
func (dbm *databaseManager) DeleteStorageSnapshot(accountHash, storageHash common.Hash) {
	db := dbm.GetDatabase(SnapshotDB)
	deleteStorageSnapshot(db, accountHash, storageHash)
}

// IterateStorageSnapshots returns an iterator for walking the entire storage
// space of a specific account.
func (dbm *databaseManager) IterateStorageSnapshots(accountHash common.Hash) Iterator {
	db := dbm.GetDatabase(SnapshotDB)
	return db.NewIterator(storageSnapshotsKey(accountHash), nil)
}

// ReadSnapshotJournal retrieves the serialized in-memory diff layers saved at
// the last shutdown. The blob is expected to be max a few 10s of megabytes.
func (dbm *databaseManager) ReadSnapshotJournal() []byte {
	db := dbm.GetDatabase(SnapshotDB)
	data, _ := db.Get(snapshotJournalKey)
	return data
}

// WriteSnapshotJournal stores the serialized in-memory diff layers to save at
// shutdown. The blob is expected to be max a few 10s of megabytes.
func (dbm *databaseManager) WriteSnapshotJournal(journal []byte) {
	db := dbm.GetDatabase(SnapshotDB)
	if err := db.Put(snapshotJournalKey, journal); err != nil {
		logger.Crit("Failed to store snapshot journal", "err", err)
	}
}

// DeleteSnapshotJournal deletes the serialized in-memory diff layers saved at
// the last shutdown.
func (dbm *databaseManager) DeleteSnapshotJournal() {
	db := dbm.GetDatabase(SnapshotDB)
	if err := db.Delete(snapshotJournalKey); err != nil {
		logger.Crit("Failed to remove snapshot journal", "err", err)
	}
}

// ReadSnapshotGenerator retrieves the serialized snapshot generator saved at
// the last shutdown.
func (dbm *databaseManager) ReadSnapshotGenerator() []byte {
	db := dbm.GetDatabase(SnapshotDB)
	data, _ := db.Get(SnapshotGeneratorKey)
	return data
}

// WriteSnapshotGenerator stores the serialized snapshot generator to save at
// shutdown.
func (dbm *databaseManager) WriteSnapshotGenerator(generator []byte) {
	db := dbm.GetDatabase(SnapshotDB)
	if err := db.Put(SnapshotGeneratorKey, generator); err != nil {
		logger.Crit("Failed to store snapshot generator", "err", err)
	}
}

// DeleteSnapshotGenerator deletes the serialized snapshot generator saved at
// the last shutdown
func (dbm *databaseManager) DeleteSnapshotGenerator() {
	db := dbm.GetDatabase(SnapshotDB)
	if err := db.Delete(SnapshotGeneratorKey); err != nil {
		logger.Crit("Failed to remove snapshot generator", "err", err)
	}
}

// ReadSnapshotRecoveryNumber retrieves the block number of the last persisted
// snapshot layer.
func (dbm *databaseManager) ReadSnapshotRecoveryNumber() *uint64 {
	db := dbm.GetDatabase(SnapshotDB)
	data, _ := db.Get(snapshotRecoveryKey)
	if len(data) == 0 {
		return nil
	}
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

// WriteSnapshotRecoveryNumber stores the block number of the last persisted
// snapshot layer.
func (dbm *databaseManager) WriteSnapshotRecoveryNumber(number uint64) {
	db := dbm.GetDatabase(SnapshotDB)
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], number)
	if err := db.Put(snapshotRecoveryKey, buf[:]); err != nil {
		logger.Crit("Failed to store snapshot recovery number", "err", err)
	}
}

// DeleteSnapshotRecoveryNumber deletes the block number of the last persisted
// snapshot layer.
func (dbm *databaseManager) DeleteSnapshotRecoveryNumber() {
	db := dbm.GetDatabase(SnapshotDB)
	if err := db.Delete(snapshotRecoveryKey); err != nil {
		logger.Crit("Failed to remove snapshot recovery number", "err", err)
	}
}

// ReadSnapshotSyncStatus retrieves the serialized sync status saved at shutdown.
func (dbm *databaseManager) ReadSnapshotSyncStatus() []byte {
	db := dbm.GetDatabase(SnapshotDB)
	data, _ := db.Get(snapshotSyncStatusKey)
	return data
}

// WriteSnapshotSyncStatus stores the serialized sync status to save at shutdown.
func (dbm *databaseManager) WriteSnapshotSyncStatus(status []byte) {
	db := dbm.GetDatabase(SnapshotDB)
	if err := db.Put(snapshotSyncStatusKey, status); err != nil {
		logger.Crit("Failed to store snapshot sync status", "err", err)
	}
}

// DeleteSnapshotSyncStatus deletes the serialized sync status saved at the last
// shutdown
func (dbm *databaseManager) DeleteSnapshotSyncStatus() {
	db := dbm.GetDatabase(SnapshotDB)
	if err := db.Delete(snapshotSyncStatusKey); err != nil {
		logger.Crit("Failed to remove snapshot sync status", "err", err)
	}
}
