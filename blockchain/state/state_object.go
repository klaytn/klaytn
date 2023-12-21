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

var errAccountDoesNotExist = errors.New("account does not exist")

type Code []byte

func (c Code) String() string {
	return string(c) // strings.Join(Disassemble(self), " ")
}

type Storage map[common.Hash]common.Hash

func (s Storage) String() (str string) {
	for key, value := range s {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
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
	// When an object is marked self-destructed it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated

	// Flag whether the account was marked as self-destructed. The self-destructed account
	// is still accessible in the scope of same transaction.
	selfDestructed bool

	// Flag whether the account was marked as deleted. A self-destructed account
	// or an account that is considered as empty will be marked as deleted at
	// the end of transaction and no longer accessible anymore.
	deleted bool

	// Flag whether the object was created in the current transaction
	created bool

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
	// State objects are RLP encoded with ExtHash preserved.
	serializer := account.NewAccountSerializerExtWithAccount(c.account)
	return rlp.Encode(w, serializer)
}

// setError remembers the first non-nil error it is called with.
func (s *stateObject) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *stateObject) markSelfdestructed() {
	s.selfDestructed = true
}

func (s *stateObject) touch() {
	s.db.journal.append(touchChange{
		account: &s.address,
	})
	if s.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		s.db.journal.dirty(s.address)
	}
}

func (s *stateObject) openStorageTrie(hash common.ExtHash, db Database) (Trie, error) {
	return db.OpenStorageTrie(hash, s.db.trieOpts)
}

func (s *stateObject) getStorageTrie(db Database) Trie {
	if s.storageTrie == nil {
		if acc := account.GetProgramAccount(s.account); acc != nil {
			var err error
			s.storageTrie, err = s.openStorageTrie(acc.GetStorageRoot(), db)
			if err != nil {
				s.storageTrie, _ = s.openStorageTrie(common.ExtHash{}, db)
				s.setError(fmt.Errorf("can't create storage trie: %v", err))
			}
		} else {
			// not a contract account, just returns the empty trie.
			s.storageTrie, _ = s.openStorageTrie(common.ExtHash{}, db)
		}
	}
	return s.storageTrie
}

// GetState retrieves a value from the account storage trie.
func (s *stateObject) GetState(db Database, key common.Hash) common.Hash {
	// If we have a dirty value for this state entry, return it
	value, dirty := s.dirtyStorage[key]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return s.GetCommittedState(db, key)
}

// GetCommittedState retrieves a value from the committed account storage trie.
func (s *stateObject) GetCommittedState(db Database, key common.Hash) common.Hash {
	// If we have the original value cached, return that
	value, cached := s.originStorage[key]
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
	if s.db.snap != nil {
		if EnabledExpensive {
			meter = &s.db.SnapshotStorageReads
		}
		// If the object was destructed in *this* block (and potentially resurrected),
		// the storage has been cleared out, and we should *not* consult the previous
		// snapshot about any storage values. The only possible alternatives are:
		//   1) resurrect happened, and new slot values were set -- those should
		//      have been handles via pendingStorage above.
		//   2) we don't have new values, and can deliver empty response back
		if _, destructed := s.db.snapDestructs[s.addrHash]; destructed {
			return common.Hash{}
		}
		enc, err = s.db.snap.Storage(s.addrHash, crypto.Keccak256Hash(key.Bytes()))
	}
	// If the snapshot is unavailable or reading from it fails, load from the database.
	if s.db.snap == nil || err != nil {
		if EnabledExpensive {
			if meter != nil {
				// If we already spent time checking the snapshot, account for it
				// and reset the readStart
				*meter += time.Since(readStart)
				readStart = time.Now()
			}
			meter = &s.db.StorageReads
		}
		// Load from DB in case it is missing.
		if enc, err = s.getStorageTrie(db).TryGet(key[:]); err != nil {
			s.setError(err)
			return common.Hash{}
		}
	}
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			s.setError(err)
		}
		value.SetBytes(content)
	}
	s.originStorage[key] = value
	return value
}

// SetState updates a value in account trie.
func (s *stateObject) SetState(db Database, key, value common.Hash) {
	// If the new value is the same as old, don't set
	prev := s.GetState(db, key)
	if prev == value {
		return
	}
	// New value is different, update and journal the change
	s.db.journal.append(storageChange{
		account:  &s.address,
		key:      key,
		prevalue: prev,
	})
	s.setState(key, value)
}

// SetStorage replaces the entire state storage with the given one.
//
// After this function is called, all original state will be ignored and state
// lookup only happens in the fake state storage.
//
// Note this function should only be used for debugging purpose.
func (s *stateObject) SetStorage(storage map[common.Hash]common.Hash) {
	// Allocate fake storage if it's nil.
	if s.fakeStorage == nil {
		s.fakeStorage = make(Storage)
	}
	for key, value := range storage {
		s.fakeStorage[key] = value
	}
	// Don't bother journal since this function should only be used for
	// debugging and the `fake` storage won't be committed to database.
}

// IsContractAccount returns true is the account has a non-empty codeHash.
func (s *stateObject) IsContractAccount() bool {
	acc := account.GetProgramAccount(s.account)
	if acc != nil && !bytes.Equal(acc.GetCodeHash(), emptyCodeHash) {
		return true
	}
	return false
}

// IsContractAvailable returns true if the account has a smart contract code hash and didn't self-destruct
func (s *stateObject) IsContractAvailable() bool {
	acc := account.GetProgramAccount(s.account)
	if acc != nil && !bytes.Equal(acc.GetCodeHash(), emptyCodeHash) && !s.selfDestructed {
		return true
	}
	return false
}

// IsProgramAccount returns true if the account implements ProgramAccount.
func (s *stateObject) IsProgramAccount() bool {
	return account.GetProgramAccount(s.account) != nil
}

func (s *stateObject) GetKey() accountkey.AccountKey {
	if ak := account.GetAccountWithKey(s.account); ak != nil {
		return ak.GetKey()
	}
	return accountkey.NewAccountKeyLegacy()
}

func (s *stateObject) setState(key, value common.Hash) {
	s.dirtyStorage[key] = value
}

func (s *stateObject) UpdateKey(newKey accountkey.AccountKey, currentBlockNumber uint64) error {
	return s.account.UpdateKey(newKey, currentBlockNumber)
}

// updateStorageTrie writes cached storage modifications into the object's storage trie.
func (s *stateObject) updateStorageTrie(db Database) Trie {
	// Track the amount of time wasted on updating the storage trie
	if EnabledExpensive {
		defer func(start time.Time) { s.db.StorageUpdates += time.Since(start) }(time.Now())
	}
	// The snapshot storage map for the object
	var storage map[common.Hash][]byte
	// Insert all the pending updates into the trie
	tr := s.getStorageTrie(db)
	for key, value := range s.dirtyStorage {
		delete(s.dirtyStorage, key)

		// Skip noop changes, persist actual changes
		if value == s.originStorage[key] {
			continue
		}
		s.originStorage[key] = value

		var v []byte
		if (value == common.Hash{}) {
			s.setError(tr.TryDelete(key[:]))
		} else {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ = rlp.EncodeToBytes(bytes.TrimLeft(value[:], "\x00"))
			s.setError(tr.TryUpdate(key[:], v))
		}
		// If state snapshotting is active, cache the data til commit
		if s.db.snap != nil {
			if storage == nil {
				// Retrieve the old storage map, if available, create a new one otherwise
				if storage = s.db.snapStorage[s.addrHash]; storage == nil {
					storage = make(map[common.Hash][]byte)
					s.db.snapStorage[s.addrHash] = storage
				}
			}
			storage[crypto.Keccak256Hash(key[:])] = v // v will be nil if it's deleted
		}
	}
	return tr
}

// updateStorageRoot sets the storage trie root to the newly updated one.
func (s *stateObject) updateStorageRoot(db Database) {
	if acc := account.GetProgramAccount(s.account); acc != nil {
		// Track the amount of time wasted on hashing the storage trie
		if EnabledExpensive {
			defer func(start time.Time) { s.db.StorageHashes += time.Since(start) }(time.Now())
		}
		acc.SetStorageRoot(s.storageTrie.HashExt())
	}
}

// setStorageRoot calls SetStorageRoot if updateStorageRoot flag is given true.
// Otherwise, it just marks the object and update their root hash later.
func (s *stateObject) setStorageRoot(updateStorageRoot bool, objectsToUpdate map[common.Address]struct{}) {
	if acc := account.GetProgramAccount(s.account); acc != nil {
		if updateStorageRoot {
			// Track the amount of time wasted on hashing the storage trie
			if EnabledExpensive {
				defer func(start time.Time) { s.db.StorageHashes += time.Since(start) }(time.Now())
			}
			acc.SetStorageRoot(s.storageTrie.HashExt())
			return
		}
		// If updateStorageRoot == false, it just marks the object and updates its storage root later.
		objectsToUpdate[s.Address()] = struct{}{}
	}
}

// CommitStorageTrie writes the storage trie of the object to db.
// This updates the storage trie root.
func (s *stateObject) CommitStorageTrie(db Database) error {
	s.updateStorageTrie(db)
	if s.dbErr != nil {
		return s.dbErr
	}
	// Track the amount of time wasted on committing the storage trie
	if EnabledExpensive {
		defer func(start time.Time) { s.db.StorageCommits += time.Since(start) }(time.Now())
	}
	if acc := account.GetProgramAccount(s.account); acc != nil {
		root, err := s.storageTrie.CommitExt(nil)
		if err != nil {
			return err
		}
		acc.SetStorageRoot(root)
	}
	return nil
}

// AddBalance adds amount to c's balance.
// It is used to add funds to the destination account of a transfer.
func (s *stateObject) AddBalance(amount *big.Int) {
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if s.empty() {
			s.touch()
		}

		return
	}
	s.SetBalance(new(big.Int).Add(s.Balance(), amount))
}

// SubBalance removes amount from c's balance.
// It is used to remove funds from the origin account of a transfer.
func (s *stateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetBalance(new(big.Int).Sub(s.Balance(), amount))
}

func (s *stateObject) SetBalance(amount *big.Int) {
	s.db.journal.append(balanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.account.GetBalance()),
	})
	s.setBalance(amount)
}

func (s *stateObject) setBalance(amount *big.Int) {
	s.account.SetBalance(amount)
}

// Return the gas back to the origin. Used by the Virtual machine or Closures
func (s *stateObject) ReturnGas(gas *big.Int) {}

func (s *stateObject) deepCopy(db *StateDB) *stateObject {
	stateObject := newObject(db, s.address, s.account.DeepCopy())
	if s.storageTrie != nil {
		stateObject.storageTrie = db.db.CopyTrie(s.storageTrie)
	}
	stateObject.code = s.code
	stateObject.dirtyStorage = s.dirtyStorage.Copy()
	stateObject.originStorage = s.originStorage.Copy()
	stateObject.selfDestructed = s.selfDestructed
	stateObject.dirtyCode = s.dirtyCode
	stateObject.deleted = s.deleted
	return stateObject
}

//
// Attribute accessors
//

// Returns the address of the contract/account
func (s *stateObject) Address() common.Address {
	return s.address
}

// Code returns the contract code associated with this object, if any.
func (s *stateObject) Code(db Database) []byte {
	if s.code != nil {
		return s.code
	}
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.setError(fmt.Errorf("can't load code hash %x: %v", s.CodeHash(), err))
	}
	s.code = code
	return code
}

// CodeSize returns the size of the contract code associated with this object,
// or zero if none. This method is an almost mirror of Code, but uses a cache
// inside the database to avoid loading codes seen recently.
func (s *stateObject) CodeSize(db Database) int {
	if s.code != nil {
		return len(s.code)
	}
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return 0
	}
	size, err := db.ContractCodeSize(common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.setError(fmt.Errorf("can't load code size %x: %v", s.CodeHash(), err))
	}
	return size
}

func (s *stateObject) SetCode(codeHash common.Hash, code []byte) error {
	prevcode := s.Code(s.db.db)
	s.db.journal.append(codeChange{
		account:  &s.address,
		prevhash: s.CodeHash(),
		prevcode: prevcode,
	})
	return s.setCode(codeHash, code)
}

func (s *stateObject) setCode(codeHash common.Hash, code []byte) error {
	acc := account.GetProgramAccount(s.account)
	if acc == nil {
		logger.Error("setCode() should be called only to a ProgramAccount!", "account address", s.address)
		return kerrors.ErrNotProgramAccount
	}
	s.code = code
	acc.SetCodeHash(codeHash[:])
	s.dirtyCode = true
	return nil
}

// IncNonce increases the nonce of the account by one with making a journal of the previous nonce.
func (s *stateObject) IncNonce() {
	nonce := s.account.GetNonce()
	s.db.journal.append(nonceChange{
		account: &s.address,
		prev:    nonce,
	})
	s.setNonce(nonce + 1)
}

func (s *stateObject) SetNonce(nonce uint64) {
	s.db.journal.append(nonceChange{
		account: &s.address,
		prev:    s.account.GetNonce(),
	})
	s.setNonce(nonce)
}

func (s *stateObject) setNonce(nonce uint64) {
	s.account.SetNonce(nonce)
}

func (s *stateObject) CodeHash() []byte {
	if acc := account.GetProgramAccount(s.account); acc != nil {
		return acc.GetCodeHash()
	}
	return emptyCodeHash
}

func (s *stateObject) Balance() *big.Int {
	return s.account.GetBalance()
}

//func (self *stateObject) HumanReadable() bool {
//	return self.account.GetHumanReadable()
//}

func (s *stateObject) Nonce() uint64 {
	return s.account.GetNonce()
}

// Never called, but must be present to allow stateObject to be used
// as a vm.Account interface that also satisfies the vm.ContractRef
// interface. Interfaces are awesome.
func (s *stateObject) Value() *big.Int {
	panic("Value on stateObject should never be called")
}
