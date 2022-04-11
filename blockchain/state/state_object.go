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
// This file is derived from core/state/state_object.go (2018/06/04).
// Modified and improved for the klaytn development.

package state

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/rlp"
)

var emptyCodeHash = crypto.Keccak256(nil)

var (
	errAccountDoesNotExist = errors.New("account does not exist")
)

type Code []byte

func (self Code) String() string {
	return string(self) //strings.Join(Disassemble(self), " ")
}

type Storage map[common.Hash]common.Hash

func (self Storage) String() (str string) {
	for key, value := range self {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

func (self Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range self {
		cpy[key] = value
	}

	return cpy
}

// stateObject represents a Klaytn account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
// Finally, call CommitStorageTrie to write the modified storage trie into a database.
type stateObject struct {
	address  common.Address
	addrHash common.Hash // hash of the address of the account
	account  account.Account
	db       *StateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Write caches.
	storageTrie Trie // storage trie, which becomes non-nil on first access
	code        Code // contract bytecode, which gets set when code is loaded

	originStorage Storage // Storage cache of original entries to dedup rewrites
	dirtyStorage  Storage // Storage entries that need to be flushed to disk
	fakeStorage   Storage // Fake storage which constructed by caller for debugging purpose.

	// Cache flags.
	// When an object is marked suicided it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool

	encoded atomic.Value // RLP-encoded data
}

type encodedData struct {
	err         error  // RLP-encoding error from stateObjectEncoder
	data        []byte // RLP-encoded stateObject
	trieHashKey []byte // hash of key used to update trie
	trieHexKey  []byte // hex string of tireHashKey
}

// empty returns whether the account is considered empty.
func (s *stateObject) empty() bool {
	return s.account.Empty()
}

// newObject creates a state object.
func newObject(db *StateDB, address common.Address, data account.Account) *stateObject {
	return &stateObject{
		db:            db,
		address:       address,
		addrHash:      crypto.Keccak256Hash(address[:]),
		account:       data,
		originStorage: make(Storage),
		dirtyStorage:  make(Storage),
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *stateObject) EncodeRLP(w io.Writer) error {
	serializer := account.NewAccountSerializerWithAccount(c.account)
	return rlp.Encode(w, serializer)
}

// setError remembers the first non-nil error it is called with.
func (self *stateObject) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *stateObject) markSuicided() {
	self.suicided = true
}

func (c *stateObject) touch() {
	c.db.journal.append(touchChange{
		account: &c.address,
	})
	if c.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		c.db.journal.dirty(c.address)
	}
}

func (c *stateObject) openStorageTrie(hash common.Hash, db Database) (Trie, error) {
	if c.db.prefetching {
		return db.OpenStorageTrieForPrefetching(hash)
	}
	return db.OpenStorageTrie(hash)
}

func (c *stateObject) getStorageTrie(db Database) Trie {
	if c.storageTrie == nil {
		if acc := account.GetProgramAccount(c.account); acc != nil {
			var err error
			c.storageTrie, err = c.openStorageTrie(acc.GetStorageRoot(), db)
			if err != nil {
				c.storageTrie, _ = c.openStorageTrie(common.Hash{}, db)
				c.setError(fmt.Errorf("can't create storage trie: %v", err))
			}
		} else {
			// not a contract account, just returns the empty trie.
			c.storageTrie, _ = c.openStorageTrie(common.Hash{}, db)
		}
	}
	return c.storageTrie
}

// GetState retrieves a value from the account storage trie.
func (self *stateObject) GetState(db Database, key common.Hash) common.Hash {
	// If we have a dirty value for this state entry, return it
	value, dirty := self.dirtyStorage[key]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return self.GetCommittedState(db, key)
}

// GetCommittedState retrieves a value from the committed account storage trie.
func (self *stateObject) GetCommittedState(db Database, key common.Hash) common.Hash {
	// If we have the original value cached, return that
	value, cached := self.originStorage[key]
	if cached {
		return value
	}
	// Track the amount of time wasted on reading the storage trie
	var (
		enc   []byte
		err   error
		meter *time.Duration
	)
	readStart := time.Now()
	if EnabledExpensive {
		// If the snap is 'under construction', the first lookup may fail. If that
		// happens, we don't want to double-count the time elapsed. Thus this
		// dance with the metering.
		defer func(start time.Time) {
			if meter != nil {
				*meter += time.Since(readStart)
			}
		}(time.Now())
	}
	if self.db.snap != nil {
		if EnabledExpensive {
			meter = &self.db.SnapshotStorageReads
		}
		// If the object was destructed in *this* block (and potentially resurrected),
		// the storage has been cleared out, and we should *not* consult the previous
		// snapshot about any storage values. The only possible alternatives are:
		//   1) resurrect happened, and new slot values were set -- those should
		//      have been handles via pendingStorage above.
		//   2) we don't have new values, and can deliver empty response back
		if _, destructed := self.db.snapDestructs[self.addrHash]; destructed {
			return common.Hash{}
		}
		enc, err = self.db.snap.Storage(self.addrHash, crypto.Keccak256Hash(key.Bytes()))
	}
	// If the snapshot is unavailable or reading from it fails, load from the database.
	if self.db.snap == nil || err != nil {
		if EnabledExpensive {
			if meter != nil {
				// If we already spent time checking the snapshot, account for it
				// and reset the readStart
				*meter += time.Since(readStart)
				readStart = time.Now()
			}
			meter = &self.db.StorageReads
		}
		// Load from DB in case it is missing.
		if enc, err = self.getStorageTrie(db).TryGet(key[:]); err != nil {
			self.setError(err)
			return common.Hash{}
		}
	}
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			self.setError(err)
		}
		value.SetBytes(content)
	}
	self.originStorage[key] = value
	return value
}

// SetState updates a value in account trie.
func (self *stateObject) SetState(db Database, key, value common.Hash) {
	// If the new value is the same as old, don't set
	prev := self.GetState(db, key)
	if prev == value {
		return
	}
	// New value is different, update and journal the change
	self.db.journal.append(storageChange{
		account:  &self.address,
		key:      key,
		prevalue: prev,
	})
	self.setState(key, value)
}

// SetStorage replaces the entire state storage with the given one.
//
// After this function is called, all original state will be ignored and state
// lookup only happens in the fake state storage.
//
// Note this function should only be used for debugging purpose.
func (self *stateObject) SetStorage(storage map[common.Hash]common.Hash) {
	// Allocate fake storage if it's nil.
	if self.fakeStorage == nil {
		self.fakeStorage = make(Storage)
	}
	for key, value := range storage {
		self.fakeStorage[key] = value
	}
	// Don't bother journal since this function should only be used for
	// debugging and the `fake` storage won't be committed to database.
}

// IsContractAccount returns true is the account has a non-empty codeHash.
func (self *stateObject) IsContractAccount() bool {
	acc := account.GetProgramAccount(self.account)
	if acc != nil && !bytes.Equal(acc.GetCodeHash(), emptyCodeHash) {
		return true
	}
	return false
}

// IsContractAvailable returns true if the account has a smart contract code hash and didn't self-destruct
func (self *stateObject) IsContractAvailable() bool {
	acc := account.GetProgramAccount(self.account)
	if acc != nil && !bytes.Equal(acc.GetCodeHash(), emptyCodeHash) && self.suicided == false {
		return true
	}
	return false
}

// IsProgramAccount returns true if the account implements ProgramAccount.
func (self *stateObject) IsProgramAccount() bool {
	return account.GetProgramAccount(self.account) != nil
}

func (self *stateObject) GetKey() accountkey.AccountKey {
	if ak := account.GetAccountWithKey(self.account); ak != nil {
		return ak.GetKey()
	}
	return accountkey.NewAccountKeyLegacy()
}

func (self *stateObject) setState(key, value common.Hash) {
	self.dirtyStorage[key] = value
}

func (self *stateObject) UpdateKey(newKey accountkey.AccountKey, currentBlockNumber uint64) error {
	return self.account.UpdateKey(newKey, currentBlockNumber)
}

// updateStorageTrie writes cached storage modifications into the object's storage trie.
func (self *stateObject) updateStorageTrie(db Database) Trie {
	// Track the amount of time wasted on updating the storage trie
	if EnabledExpensive {
		defer func(start time.Time) { self.db.StorageUpdates += time.Since(start) }(time.Now())
	}
	// The snapshot storage map for the object
	var storage map[common.Hash][]byte
	// Insert all the pending updates into the trie
	tr := self.getStorageTrie(db)
	for key, value := range self.dirtyStorage {
		delete(self.dirtyStorage, key)

		// Skip noop changes, persist actual changes
		if value == self.originStorage[key] {
			continue
		}
		self.originStorage[key] = value

		var v []byte
		if (value == common.Hash{}) {
			self.setError(tr.TryDelete(key[:]))
		} else {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ = rlp.EncodeToBytes(bytes.TrimLeft(value[:], "\x00"))
			self.setError(tr.TryUpdate(key[:], v))
		}
		// If state snapshotting is active, cache the data til commit
		if self.db.snap != nil {
			if storage == nil {
				// Retrieve the old storage map, if available, create a new one otherwise
				if storage = self.db.snapStorage[self.addrHash]; storage == nil {
					storage = make(map[common.Hash][]byte)
					self.db.snapStorage[self.addrHash] = storage
				}
			}
			storage[crypto.Keccak256Hash(key[:])] = v // v will be nil if it's deleted
		}
	}
	return tr
}

// updateStorageRoot sets the storage trie root to the newly updated one.
func (self *stateObject) updateStorageRoot(db Database) {
	if acc := account.GetProgramAccount(self.account); acc != nil {
		// Track the amount of time wasted on hashing the storage trie
		if EnabledExpensive {
			defer func(start time.Time) { self.db.StorageHashes += time.Since(start) }(time.Now())
		}
		acc.SetStorageRoot(self.storageTrie.Hash())
	}
}

// setStorageRoot calls SetStorageRoot if updateStorageRoot flag is given true.
// Otherwise, it just marks the object and update their root hash later.
func (self *stateObject) setStorageRoot(updateStorageRoot bool, objectsToUpdate map[common.Address]struct{}) {
	if acc := account.GetProgramAccount(self.account); acc != nil {
		if updateStorageRoot {
			// Track the amount of time wasted on hashing the storage trie
			if EnabledExpensive {
				defer func(start time.Time) { self.db.StorageHashes += time.Since(start) }(time.Now())
			}
			acc.SetStorageRoot(self.storageTrie.Hash())
			return
		}
		// If updateStorageRoot == false, it just marks the object and updates its storage root later.
		objectsToUpdate[self.Address()] = struct{}{}
	}
}

// CommitStorageTrie writes the storage trie of the object to db.
// This updates the storage trie root.
func (self *stateObject) CommitStorageTrie(db Database) error {
	self.updateStorageTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}
	// Track the amount of time wasted on committing the storage trie
	if EnabledExpensive {
		defer func(start time.Time) { self.db.StorageCommits += time.Since(start) }(time.Now())
	}
	if acc := account.GetProgramAccount(self.account); acc != nil {
		root, err := self.storageTrie.Commit(nil)
		if err != nil {
			return err
		}
		acc.SetStorageRoot(root)
	}
	return nil
}

// AddBalance removes amount from c's balance.
// It is used to add funds to the destination account of a transfer.
func (c *stateObject) AddBalance(amount *big.Int) {
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if c.empty() {
			c.touch()
		}

		return
	}
	c.SetBalance(new(big.Int).Add(c.Balance(), amount))
}

// SubBalance removes amount from c's balance.
// It is used to remove funds from the origin account of a transfer.
func (c *stateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	c.SetBalance(new(big.Int).Sub(c.Balance(), amount))
}

func (self *stateObject) SetBalance(amount *big.Int) {
	self.db.journal.append(balanceChange{
		account: &self.address,
		prev:    new(big.Int).Set(self.account.GetBalance()),
	})
	self.setBalance(amount)
}

func (self *stateObject) setBalance(amount *big.Int) {
	self.account.SetBalance(amount)
}

// Return the gas back to the origin. Used by the Virtual machine or Closures
func (c *stateObject) ReturnGas(gas *big.Int) {}

func (self *stateObject) deepCopy(db *StateDB) *stateObject {
	stateObject := newObject(db, self.address, self.account.DeepCopy())
	if self.storageTrie != nil {
		stateObject.storageTrie = db.db.CopyTrie(self.storageTrie)
	}
	stateObject.code = self.code
	stateObject.dirtyStorage = self.dirtyStorage.Copy()
	stateObject.originStorage = self.originStorage.Copy()
	stateObject.suicided = self.suicided
	stateObject.dirtyCode = self.dirtyCode
	stateObject.deleted = self.deleted
	return stateObject
}

//
// Attribute accessors
//

// Returns the address of the contract/account
func (c *stateObject) Address() common.Address {
	return c.address
}

// Code returns the contract code associated with this object, if any.
func (self *stateObject) Code(db Database) []byte {
	if self.code != nil {
		return self.code
	}
	if bytes.Equal(self.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(common.BytesToHash(self.CodeHash()))
	if err != nil {
		self.setError(fmt.Errorf("can't load code hash %x: %v", self.CodeHash(), err))
	}
	self.code = code
	return code
}

// CodeSize returns the size of the contract code associated with this object,
// or zero if none. This method is an almost mirror of Code, but uses a cache
// inside the database to avoid loading codes seen recently.
func (self *stateObject) CodeSize(db Database) int {
	if self.code != nil {
		return len(self.code)
	}
	if bytes.Equal(self.CodeHash(), emptyCodeHash) {
		return 0
	}
	size, err := db.ContractCodeSize(common.BytesToHash(self.CodeHash()))
	if err != nil {
		self.setError(fmt.Errorf("can't load code size %x: %v", self.CodeHash(), err))
	}
	return size
}

func (self *stateObject) SetCode(codeHash common.Hash, code []byte) error {
	prevcode := self.Code(self.db.db)
	self.db.journal.append(codeChange{
		account:  &self.address,
		prevhash: self.CodeHash(),
		prevcode: prevcode,
	})
	return self.setCode(codeHash, code)
}

func (self *stateObject) setCode(codeHash common.Hash, code []byte) error {
	acc := account.GetProgramAccount(self.account)
	if acc == nil {
		logger.Error("setCode() should be called only to a ProgramAccount!", "account address", self.address)
		return kerrors.ErrNotProgramAccount
	}
	self.code = code
	acc.SetCodeHash(codeHash[:])
	self.dirtyCode = true
	return nil
}

// IncNonce increases the nonce of the account by one with making a journal of the previous nonce.
func (self *stateObject) IncNonce() {
	nonce := self.account.GetNonce()
	self.db.journal.append(nonceChange{
		account: &self.address,
		prev:    nonce,
	})
	self.setNonce(nonce + 1)
}

func (self *stateObject) SetNonce(nonce uint64) {
	self.db.journal.append(nonceChange{
		account: &self.address,
		prev:    self.account.GetNonce(),
	})
	self.setNonce(nonce)
}

func (self *stateObject) setNonce(nonce uint64) {
	self.account.SetNonce(nonce)
}

func (self *stateObject) CodeHash() []byte {
	if acc := account.GetProgramAccount(self.account); acc != nil {
		return acc.GetCodeHash()
	}
	return emptyCodeHash
}

func (self *stateObject) Balance() *big.Int {
	return self.account.GetBalance()
}

//func (self *stateObject) HumanReadable() bool {
//	return self.account.GetHumanReadable()
//}

func (self *stateObject) Nonce() uint64 {
	return self.account.GetNonce()
}

// Never called, but must be present to allow stateObject to be used
// as a vm.Account interface that also satisfies the vm.ContractRef
// interface. Interfaces are awesome.
func (self *stateObject) Value() *big.Int {
	panic("Value on stateObject should never be called")
}
