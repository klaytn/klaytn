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
// This file is derived from core/state/statedb.go (2018/06/04).
// Modified and improved for the klaytn development.

package state

import (
	"fmt"
	"math/big"
	"sort"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/snapshot"
	"github.com/klaytn/klaytn/storage/statedb"
)

type revision struct {
	id           int
	journalIndex int
}

var (
	// emptyState is the known hash of an empty state trie entry.
	emptyState = crypto.Keccak256Hash(nil)

	// emptyCode is the known hash of the empty EVM bytecode.
	emptyCode = crypto.Keccak256Hash(nil)

	logger = log.NewModuleLogger(log.BlockchainState)

	// TODO-Klaytn EnabledExpensive and DBConfig.EnableDBPerfMetrics will be merged
	EnabledExpensive = false
)

// StateDBs within the Klaytn protocol are used to cache stateObjects from Merkle Patricia Trie
// and mediate the operations to them.
type StateDB struct {
	db       Database
	trie     Trie
	trieOpts *statedb.TrieOpts

	snaps         *snapshot.Tree
	snap          snapshot.Snapshot
	snapDestructs map[common.Hash]struct{}
	snapAccounts  map[common.Hash][]byte
	snapStorage   map[common.Hash]map[common.Hash][]byte

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateObjects             map[common.Address]*stateObject
	stateObjectsDirty        map[common.Address]struct{}
	stateObjectsDirtyStorage map[common.Address]struct{}

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash common.Hash
	txIndex      int
	logs         map[common.Hash][]*types.Log
	logSize      uint

	preimages map[common.Hash][]byte

	// Per-transaction access list
	accessList *accessList

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionId int

	prefetching bool

	// Measurements gathered during execution for debugging purposes
	AccountReads         time.Duration
	AccountHashes        time.Duration
	AccountUpdates       time.Duration
	AccountCommits       time.Duration
	StorageReads         time.Duration
	StorageHashes        time.Duration
	StorageUpdates       time.Duration
	StorageCommits       time.Duration
	SnapshotAccountReads time.Duration
	SnapshotStorageReads time.Duration
	SnapshotCommits      time.Duration
}

// Create a new state from a given trie.
func New(root common.Hash, db Database, snaps *snapshot.Tree, opts *statedb.TrieOpts) (*StateDB, error) {
	tr, err := db.OpenTrie(root, opts)
	if err != nil {
		return nil, err
	}
	sdb := &StateDB{
		db:                       db,
		trie:                     tr,
		trieOpts:                 opts,
		snaps:                    snaps,
		stateObjects:             make(map[common.Address]*stateObject),
		stateObjectsDirtyStorage: make(map[common.Address]struct{}),
		stateObjectsDirty:        make(map[common.Address]struct{}),
		logs:                     make(map[common.Hash][]*types.Log),
		preimages:                make(map[common.Hash][]byte),
		accessList:               newAccessList(),
		journal:                  newJournal(),
	}
	if sdb.snaps != nil {
		if sdb.snap = sdb.snaps.Snapshot(root); sdb.snap != nil {
			sdb.snapDestructs = make(map[common.Hash]struct{})
			sdb.snapAccounts = make(map[common.Hash][]byte)
			sdb.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
		}
	}
	if opts != nil && opts.Prefetching {
		sdb.prefetching = true
	}
	return sdb, nil
}

// RLockGCCachedNode locks the GC lock of CachedNode.
func (s *StateDB) LockGCCachedNode() {
	s.db.RLockGCCachedNode()
}

// RUnlockGCCachedNode unlocks the GC lock of CachedNode.
func (s *StateDB) UnlockGCCachedNode() {
	s.db.RUnlockGCCachedNode()
}

// setError remembers the first non-nil error it is called with.
func (s *StateDB) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *StateDB) Error() error {
	return s.dbErr
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying state trie to avoid reloading data for the next operations.
func (s *StateDB) Reset(root common.Hash) error {
	tr, err := s.db.OpenTrie(root, s.trieOpts)
	if err != nil {
		return err
	}
	s.trie = tr
	s.stateObjects = make(map[common.Address]*stateObject)
	s.stateObjectsDirty = make(map[common.Address]struct{})
	s.thash = common.Hash{}
	s.bhash = common.Hash{}
	s.txIndex = 0
	s.logs = make(map[common.Hash][]*types.Log)
	s.logSize = 0
	s.preimages = make(map[common.Hash][]byte)
	s.clearJournalAndRefund()
	s.accessList = newAccessList()
	return nil
}

func (s *StateDB) AddLog(log *types.Log) {
	s.journal.append(addLogChange{txhash: s.thash})

	log.TxHash = s.thash
	log.BlockHash = s.bhash
	log.TxIndex = uint(s.txIndex)
	log.Index = s.logSize
	s.logs[s.thash] = append(s.logs[s.thash], log)
	s.logSize++
}

func (s *StateDB) GetLogs(hash common.Hash) []*types.Log {
	return s.logs[hash]
}

func (s *StateDB) Logs() []*types.Log {
	var logs []*types.Log
	for _, lgs := range s.logs {
		logs = append(logs, lgs...)
	}
	return logs
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (s *StateDB) AddPreimage(hash common.Hash, preimage []byte) {
	if _, ok := s.preimages[hash]; !ok {
		s.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		s.preimages[hash] = pi
	}
}

// Preimages returns a list of SHA3 preimages that have been submitted.
func (s *StateDB) Preimages() map[common.Hash][]byte {
	return s.preimages
}

// AddRefund adds gas to the refund counter
func (s *StateDB) AddRefund(gas uint64) {
	s.journal.append(refundChange{prev: s.refund})
	s.refund += gas
}

// SubRefund removes gas from the refund counter.
// This method will panic if the refund counter goes below zero
func (s *StateDB) SubRefund(gas uint64) {
	s.journal.append(refundChange{prev: s.refund})
	if gas > s.refund {
		panic("Refund counter below zero")
	}
	s.refund -= gas
}

// Exist reports whether the given account address exists in the state.
// Notably this also returns true for self-destructed accounts.
func (s *StateDB) Exist(addr common.Address) bool {
	return s.getStateObject(addr) != nil
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (s *StateDB) Empty(addr common.Address) bool {
	so := s.getStateObject(addr)
	return so == nil || so.empty()
}

// Retrieve the balance from the given address or 0 if object not found
func (s *StateDB) GetBalance(addr common.Address) *big.Int {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Balance()
	}
	return common.Big0
}

func (s *StateDB) GetNonce(addr common.Address) uint64 {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}
	return 0
}

func (s *StateDB) GetCode(addr common.Address) []byte {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Code(s.db)
	}
	return nil
}

func (s *StateDB) GetAccount(addr common.Address) account.Account {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.account
	}
	return nil
}

func (s *StateDB) IsContractAccount(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.IsContractAccount()
	}
	return false
}

//func (self *StateDB) IsHumanReadable(addr common.Address) bool {
//	stateObject := self.getStateObject(addr)
//	if stateObject != nil {
//		return stateObject.HumanReadable()
//	}
//	return false
//}

func (s *StateDB) GetCodeSize(addr common.Address) int {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.CodeSize(s.db)
	}
	return 0
}

func (s *StateDB) GetCodeHash(addr common.Address) common.Hash {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return common.BytesToHash(emptyCodeHash)
	}
	return common.BytesToHash(stateObject.CodeHash())
}

// GetState retrieves a value from the given account's storage trie.
func (s *StateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetState(s.db, hash)
	}
	return common.Hash{}
}

// GetCommittedState retrieves a value from the given account's committed storage trie.
func (s *StateDB) GetCommittedState(addr common.Address, hash common.Hash) common.Hash {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetCommittedState(s.db, hash)
	}
	return common.Hash{}
}

// IsContractAvailable returns true if the account corresponding to the given address implements ProgramAccount.
func (s *StateDB) IsContractAvailable(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.IsContractAvailable()
	}
	return false
}

// IsProgramAccount returns true if the account corresponding to the given address implements ProgramAccount.
func (s *StateDB) IsProgramAccount(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.IsProgramAccount()
	}
	return false
}

func (s *StateDB) IsValidCodeFormat(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		pa := account.GetProgramAccount(stateObject.account)
		if pa != nil {
			return pa.GetCodeFormat().Validate()
		}
		return false
	}
	return false
}

// GetVmVersion return false when getStateObject(addr) or GetProgramAccount(stateObject.account) is failed.
func (s *StateDB) GetVmVersion(addr common.Address) (params.VmVersion, bool) {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		pa := account.GetProgramAccount(stateObject.account)
		if pa != nil {
			return pa.GetVmVersion(), true
		}
		return params.VmVersion0, false
	}
	return params.VmVersion0, false
}

func (s *StateDB) GetKey(addr common.Address) accountkey.AccountKey {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetKey()
	}
	return accountkey.NewAccountKeyLegacy()
}

// Database retrieves the low level database supporting the lower level trie ops.
func (s *StateDB) Database() Database {
	return s.db
}

// StorageTrie returns the storage trie of an account.
// The return value is a copy and is nil for non-existent accounts.
func (s *StateDB) StorageTrie(addr common.Address) Trie {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return nil
	}
	cpy := stateObject.deepCopy(s)
	return cpy.updateStorageTrie(s.db)
}

func (s *StateDB) HasSelfDestructed(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.selfDestructed
	}
	return false
}

/*
 * SETTERS
 */

// AddBalance adds amount to the account associated with addr.
func (s *StateDB) AddBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (s *StateDB) SubBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)
	}
}

func (s *StateDB) SetBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetBalance(amount)
	}
}

// IncNonce increases the nonce of the account of the given address by one.
func (s *StateDB) IncNonce(addr common.Address) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.IncNonce()
	}
}

func (s *StateDB) SetNonce(addr common.Address, nonce uint64) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

func (s *StateDB) SetCode(addr common.Address, code []byte) error {
	stateObject := s.GetOrNewSmartContract(addr)
	if stateObject != nil {
		return stateObject.SetCode(crypto.Keccak256Hash(code), code)
	}

	return nil
}

func (s *StateDB) SetState(addr common.Address, key, value common.Hash) {
	stateObject := s.GetOrNewSmartContract(addr)
	if stateObject != nil {
		stateObject.SetState(s.db, key, value)
	}
}

// SetStorage replaces the entire storage for the specified account with given
// storage. This function should only be used for debugging.
func (s *StateDB) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetStorage(storage)
	}
}

// UpdateKey updates the account's key with the given key.
func (s *StateDB) UpdateKey(addr common.Address, newKey accountkey.AccountKey, currentBlockNumber uint64) error {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.UpdateKey(newKey, currentBlockNumber)
	}

	return errAccountDoesNotExist
}

// SelfDestruct marks the given account as self-destructed.
// This clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after SelfDestruct.
func (s *StateDB) SelfDestruct(addr common.Address) {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return
	}
	s.journal.append(selfDestructChange{
		account:     &addr,
		prev:        stateObject.selfDestructed,
		prevbalance: new(big.Int).Set(stateObject.Balance()),
	})
	stateObject.markSelfdestructed()
	stateObject.account.SetBalance(new(big.Int))
}

func (s *StateDB) SelfDestruct6780(addr common.Address) {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return
	}

	if stateObject.created {
		s.SelfDestruct(addr)
	}
}

//
// Setting, updating & deleting state object methods.
//

// updateStateObject writes the given object to the statedb.
func (s *StateDB) updateStateObject(stateObject *stateObject) {
	// Track the amount of time wasted on updating the account from the trie
	if EnabledExpensive {
		defer func(start time.Time) { s.AccountUpdates += time.Since(start) }(time.Now())
	}
	addr := stateObject.Address()
	var snapshotData []byte
	if data := stateObject.encoded.Load(); data != nil {
		encodedData := data.(*encodedData)
		if encodedData.err != nil {
			panic(fmt.Errorf("can't encode object at %x: %v", addr[:], encodedData.err))
		}
		s.setError(s.trie.TryUpdateWithKeys(addr[:],
			encodedData.trieHashKey, encodedData.trieHexKey, encodedData.data))
		stateObject.encoded = atomic.Value{}
		snapshotData = encodedData.data
	} else {
		data, err := rlp.EncodeToBytes(stateObject)
		if err != nil {
			panic(fmt.Errorf("can't encode object at %x: %v", addr[:], err))
		}
		s.setError(s.trie.TryUpdate(addr[:], data))
		snapshotData = data
	}

	// If state snapshotting is active, cache the data til commit. Note, this
	// update mechanism is not symmetric to the deletion, because whereas it is
	// enough to track account updates at commit time, deletions need tracking
	// at transaction boundary level to ensure we capture state clearing.
	if s.snap != nil {
		s.snapAccounts[stateObject.addrHash] = snapshotData
	}
}

// deleteStateObject removes the given object from the state trie.
func (s *StateDB) deleteStateObject(stateObject *stateObject) {
	// Track the amount of time wasted on deleting the account from the trie
	if EnabledExpensive {
		defer func(start time.Time) { s.AccountUpdates += time.Since(start) }(time.Now())
	}
	stateObject.deleted = true
	addr := stateObject.Address()
	s.setError(s.trie.TryDelete(addr[:]))
}

// getStateObject retrieves a state object given by the address, returning nil if
// the object is not found or was deleted in this execution context. If you need
// to differentiate between non-existent/just-deleted, use getDeletedStateObject.
func (s *StateDB) getStateObject(addr common.Address) *stateObject {
	if obj := s.getDeletedStateObject(addr); obj != nil && !obj.deleted {
		return obj
	}
	return nil
}

// getDeletedStateObject is similar to getStateObject, but instead of returning
// nil for a deleted state object, it returns the actual object with the deleted
// flag set. This is needed by the state journal to revert to the correct s-
// destructed object instead of wiping all knowledge about the state object.
func (s *StateDB) getDeletedStateObject(addr common.Address) *stateObject {
	// First, check stateObjects if there is "live" object.
	if obj := s.stateObjects[addr]; obj != nil {
		return obj
	}
	// If no live objects are available, attempt to use snapshots
	var (
		acc account.Account
		err error
	)
	if s.snap != nil {
		if EnabledExpensive {
			defer func(start time.Time) { s.SnapshotAccountReads += time.Since(start) }(time.Now())
		}
		if acc, err = s.snap.Account(crypto.Keccak256Hash(addr.Bytes())); err == nil {
			if acc == nil {
				return nil
			}
		}
	}
	// If snapshot unavailable or reading from it failed, load from the database
	if s.snap == nil || err != nil {
		// Track the amount of time wasted on loading the object from the database
		if EnabledExpensive {
			defer func(start time.Time) { s.AccountReads += time.Since(start) }(time.Now())
		}
		// Second, the object for given address is not cached.
		// Load the object from the database.
		enc, err := s.trie.TryGet(addr[:])
		if len(enc) == 0 {
			s.setError(err)
			return nil
		}
		serializer := account.NewAccountSerializer()
		if err := rlp.DecodeBytes(enc, serializer); err != nil {
			logger.Error("Failed to decode state object", "addr", addr, "err", err)
			return nil
		}
		acc = serializer.GetAccount()
	}
	// Insert into the live set.
	obj := newObject(s, addr, acc)
	s.setStateObject(obj)

	return obj
}

func (s *StateDB) setStateObject(object *stateObject) {
	s.stateObjects[object.Address()] = object
}

// Retrieve a state object or create a new state object if nil.
func (s *StateDB) GetOrNewStateObject(addr common.Address) *stateObject {
	stateObject := s.getStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		stateObject, _ = s.createObject(addr)
	}
	return stateObject
}

// Retrieve a state object or create a new state object if nil.
func (s *StateDB) GetOrNewSmartContract(addr common.Address) *stateObject {
	stateObject := s.getStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		stateObject, _ = s.createObjectWithMap(addr, account.SmartContractAccountType, map[account.AccountValueKeyType]interface{}{
			account.AccountValueKeyNonce:         uint64(1),
			account.AccountValueKeyHumanReadable: false,
			account.AccountValueKeyAccountKey:    accountkey.NewAccountKeyFail(),
		})
	}
	return stateObject
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (s *StateDB) createObject(addr common.Address) (newobj, prev *stateObject) {
	prev = s.getDeletedStateObject(addr) // Note, prev might have been deleted, we need that!

	var prevdestruct bool
	if s.snap != nil && prev != nil {
		_, prevdestruct = s.snapDestructs[prev.addrHash]
		if !prevdestruct {
			s.snapDestructs[prev.addrHash] = struct{}{}
		}
	}
	acc, err := account.NewAccountWithType(account.ExternallyOwnedAccountType)
	if err != nil {
		logger.Error("An error occurred on call NewAccountWithType", "err", err)
	}
	newobj = newObject(s, addr, acc)
	newobj.setNonce(0) // sets the object to dirty
	if prev == nil {
		s.journal.append(createObjectChange{account: &addr})
	} else {
		s.journal.append(resetObjectChange{prev: prev, prevdestruct: prevdestruct})
	}

	newobj.created = true

	s.setStateObject(newobj)
	if prev != nil && !prev.deleted {
		return newobj, prev
	}
	return newobj, nil
}

// createObjectWithMap creates a new state object with the given parameters (accountType and values).
// If there is an existing account with the given address, it is overwritten and
// returned as the second return value.
func (s *StateDB) createObjectWithMap(addr common.Address, accountType account.AccountType,
	values map[account.AccountValueKeyType]interface{},
) (newobj, prev *stateObject) {
	prev = s.getDeletedStateObject(addr) // Note, prev might have been deleted, we need that!

	var prevdestruct bool
	if s.snap != nil && prev != nil {
		_, prevdestruct = s.snapDestructs[prev.addrHash]
		if !prevdestruct {
			s.snapDestructs[prev.addrHash] = struct{}{}
		}
	}
	acc, err := account.NewAccountWithMap(accountType, values)
	if err != nil {
		logger.Error("An error occurred on call NewAccountWithMap", "err", err)
	}
	newobj = newObject(s, addr, acc)
	newobj.setNonce(0) // sets the object to dirty
	if prev == nil {
		s.journal.append(createObjectChange{account: &addr})
	} else {
		s.journal.append(resetObjectChange{prev: prev, prevdestruct: prevdestruct})
	}
	s.setStateObject(newobj)
	if prev != nil && !prev.deleted {
		return newobj, prev
	}
	return newobj, nil
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists, the balance is carried over to the new account.
// Carrying over the balance ensures that Ether doesn't disappear.
//
// CreateAccount is currently used for test code only. Instead,
// use CreateEOA, CreateSmartContractAccount, or CreateSmartContractAccountWithKey to create a typed account.
func (s *StateDB) CreateAccount(addr common.Address) {
	new, prev := s.createObject(addr)
	if prev != nil {
		new.setBalance(prev.account.GetBalance())
	}
}

func (s *StateDB) CreateEOA(addr common.Address, humanReadable bool, key accountkey.AccountKey) {
	values := map[account.AccountValueKeyType]interface{}{
		account.AccountValueKeyHumanReadable: humanReadable,
		account.AccountValueKeyAccountKey:    key,
	}
	new, prev := s.createObjectWithMap(addr, account.ExternallyOwnedAccountType, values)
	if prev != nil {
		new.setBalance(prev.account.GetBalance())
	}
}

func (s *StateDB) CreateSmartContractAccount(addr common.Address, format params.CodeFormat, r params.Rules) {
	s.CreateSmartContractAccountWithKey(addr, false, accountkey.NewAccountKeyFail(), format, r)
}

func (s *StateDB) CreateSmartContractAccountWithKey(addr common.Address, humanReadable bool, key accountkey.AccountKey, format params.CodeFormat, r params.Rules) {
	values := map[account.AccountValueKeyType]interface{}{
		account.AccountValueKeyNonce:         uint64(1),
		account.AccountValueKeyHumanReadable: humanReadable,
		account.AccountValueKeyAccountKey:    key,
		account.AccountValueKeyCodeInfo:      params.NewCodeInfoWithRules(format, r),
	}
	new, prev := s.createObjectWithMap(addr, account.SmartContractAccountType, values)
	if prev != nil {
		new.setBalance(prev.account.GetBalance())
	}
}

func (s *StateDB) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) {
	so := s.getStateObject(addr)
	if so == nil {
		return
	}
	it := statedb.NewIterator(so.getStorageTrie(s.db).NodeIterator(nil))
	for it.Next() {
		key := common.BytesToHash(s.trie.GetKey(it.Key))
		if value, dirty := so.dirtyStorage[key]; dirty {
			cb(key, value)
			continue
		}
		enc := it.Value
		if len(enc) > 0 {
			_, content, _, _ := rlp.Split(enc)
			cb(key, common.BytesToHash(content))
		} else {
			cb(key, common.Hash{})
		}
	}
}

// Copy creates a deep, independent copy of the state.
// Snapshots of the copied state cannot be applied to the copy.
func (s *StateDB) Copy() *StateDB {
	// Copy all the basic fields, initialize the memory ones
	state := &StateDB{
		db:                s.db,
		trie:              s.db.CopyTrie(s.trie),
		stateObjects:      make(map[common.Address]*stateObject, len(s.journal.dirties)),
		stateObjectsDirty: make(map[common.Address]struct{}, len(s.journal.dirties)),
		refund:            s.refund,
		logs:              make(map[common.Hash][]*types.Log, len(s.logs)),
		logSize:           s.logSize,
		preimages:         make(map[common.Hash][]byte),
		journal:           newJournal(),
	}
	// Copy the dirty states, logs, and preimages
	for addr := range s.journal.dirties {
		// As documented [here](https://github.com/ethereum/go-ethereum/pull/16485#issuecomment-380438527),
		// and in the Finalise-method, there is a case where an object is in the journal but not
		// in the stateObjects: OOG after touch on ripeMD prior to Byzantium. Thus, we need to check for
		// nil
		if object, exist := s.stateObjects[addr]; exist {
			state.stateObjects[addr] = object.deepCopy(state)
			state.stateObjectsDirty[addr] = struct{}{}
		}
	}
	// Above, we don't copy the actual journal. This means that if the copy is copied, the
	// loop above will be a no-op, since the copy's journal is empty.
	// Thus, here we iterate over stateObjects, to enable copies of copies
	for addr := range s.stateObjectsDirty {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = s.stateObjects[addr].deepCopy(state)
			state.stateObjectsDirty[addr] = struct{}{}
		}
	}

	deepCopyLogs(s, state)

	for hash, preimage := range s.preimages {
		state.preimages[hash] = preimage
	}

	// Do we need to copy the access list? In practice: No. At the start of a
	// transaction, the access list is empty. In practice, we only ever copy state
	// _between_ transactions/blocks, never in the middle of a transaction.
	// However, it doesn't cost us much to copy an empty list, so we do it anyway
	// to not blow up if we ever decide copy it in the middle of a transaction
	state.accessList = s.accessList.Copy()

	if s.snaps != nil {
		// In order for the miner to be able to use and make additions
		// to the snapshot tree, we need to copy that aswell.
		// Otherwise, any block mined by ourselves will cause gaps in the tree,
		// and force the miner to operate trie-backed only
		state.snaps = s.snaps
		state.snap = s.snap
		// deep copy needed
		state.snapDestructs = make(map[common.Hash]struct{})
		for k, v := range s.snapDestructs {
			state.snapDestructs[k] = v
		}
		state.snapAccounts = make(map[common.Hash][]byte)
		for k, v := range s.snapAccounts {
			state.snapAccounts[k] = v
		}
		state.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
		for k, v := range s.snapStorage {
			temp := make(map[common.Hash][]byte)
			for kk, vv := range v {
				temp[kk] = vv
			}
			state.snapStorage[k] = temp
		}
	}
	return state
}

// deepCopyLogs deep-copies StateDB.logs from the left to the right.
func deepCopyLogs(from, to *StateDB) {
	for hash, logs := range from.logs {
		copied := make([]*types.Log, len(logs))
		for i, log := range logs {
			copied[i] = new(types.Log)
			*copied[i] = *log
		}
		to.logs[hash] = copied
	}
}

// Snapshot returns an identifier for the current revision of the state.
func (s *StateDB) Snapshot() int {
	id := s.nextRevisionId
	s.nextRevisionId++
	s.validRevisions = append(s.validRevisions, revision{id, s.journal.length()})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (s *StateDB) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(s.validRevisions), func(i int) bool {
		return s.validRevisions[i].id >= revid
	})
	if idx == len(s.validRevisions) || s.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := s.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	s.journal.revert(s, snapshot)
	s.validRevisions = s.validRevisions[:idx]
}

// GetRefund returns the current value of the refund counter.
func (s *StateDB) GetRefund() uint64 {
	return s.refund
}

// Finalise finalises the state by removing the self destructed objects
// and clears the journal as well as the refunds.
func (stateDB *StateDB) Finalise(deleteEmptyObjects bool, setStorageRoot bool) {
	for addr := range stateDB.journal.dirties {
		so, exist := stateDB.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `stateDB.journal.dirties` but not in `stateDB.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}

		if so.selfDestructed || (deleteEmptyObjects && so.empty()) {
			stateDB.deleteStateObject(so)

			// If state snapshotting is active, also mark the destruction there.
			// Note, we can't do this only at the end of a block because multiple
			// transactions within the same block might self destruct and then
			// ressurrect an account; but the snapshotter needs both events.
			if stateDB.snap != nil {
				stateDB.snapDestructs[so.addrHash] = struct{}{} // We need to maintain account deletions explicitly (will remain set indefinitely)
				delete(stateDB.snapAccounts, so.addrHash)       // Clear out any previously updated account data (may be recreated via a ressurrect)
				delete(stateDB.snapStorage, so.addrHash)        // Clear out any previously updated storage data (may be recreated via a ressurrect)
			}
		} else {
			so.updateStorageTrie(stateDB.db)
			so.setStorageRoot(setStorageRoot, stateDB.stateObjectsDirtyStorage)
			stateDB.updateStateObject(so)
		}
		so.created = false
		stateDB.stateObjectsDirty[addr] = struct{}{}
	}
	// Invalidate journal because reverting across transactions is not allowed.
	stateDB.clearJournalAndRefund()

	if setStorageRoot && len(stateDB.stateObjectsDirtyStorage) > 0 {
		for addr := range stateDB.stateObjectsDirtyStorage {
			so, exist := stateDB.stateObjects[addr]
			if exist {
				so.updateStorageRoot(stateDB.db)
				stateDB.updateStateObject(so)
			}
		}
		stateDB.stateObjectsDirtyStorage = make(map[common.Address]struct{})
	}
}

// IntermediateRoot computes the current root hash of the state statedb.
// It is called in between transactions to get the root hash that
// goes into transaction receipts.
func (s *StateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	s.Finalise(deleteEmptyObjects, true)
	// Track the amount of time wasted on hashing the account trie
	if EnabledExpensive {
		defer func(start time.Time) { s.AccountHashes += time.Since(start) }(time.Now())
	}
	return s.trie.Hash()
}

// Prepare sets the current transaction hash and index and block hash which is
// used when the EVM emits new state logs.
func (s *StateDB) Prepare(thash, bhash common.Hash, ti int) {
	s.thash = thash
	s.bhash = bhash
	s.txIndex = ti
}

func (s *StateDB) clearJournalAndRefund() {
	s.journal = newJournal()
	s.validRevisions = s.validRevisions[:0]
	s.refund = 0
}

// Commit writes the state to the underlying in-memory trie database.
func (s *StateDB) Commit(deleteEmptyObjects bool) (root common.Hash, err error) {
	if s.dbErr != nil {
		return common.Hash{}, fmt.Errorf("commit aborted due to earlier error: %v", s.dbErr)
	}

	defer s.clearJournalAndRefund()

	for addr := range s.journal.dirties {
		s.stateObjectsDirty[addr] = struct{}{}
	}

	objectEncoder := getStateObjectEncoder(len(s.stateObjects))
	var stateObjectsToUpdate []*stateObject
	// Commit objects to the trie.
	for addr, stateObject := range s.stateObjects {
		_, isDirty := s.stateObjectsDirty[addr]
		switch {
		case stateObject.selfDestructed || (isDirty && deleteEmptyObjects && stateObject.empty()):
			// If the object has been removed, don't bother syncing it
			// and just mark it for deletion in the trie.
			s.deleteStateObject(stateObject)
		case isDirty:
			if stateObject.IsProgramAccount() {
				// Write any contract code associated with the state object.
				if stateObject.code != nil && stateObject.dirtyCode {
					s.db.TrieDB().DiskDB().WriteCode(common.BytesToHash(stateObject.CodeHash()), stateObject.code)
					stateObject.dirtyCode = false
				}
				// Write any storage changes in the state object to its storage trie.
				if err := stateObject.CommitStorageTrie(s.db); err != nil {
					return common.Hash{}, err
				}
			}
			// Update the object in the main account trie.
			stateObjectsToUpdate = append(stateObjectsToUpdate, stateObject)
			objectEncoder.encode(stateObject)
		}
		delete(s.stateObjectsDirty, addr)
	}

	for _, so := range stateObjectsToUpdate {
		s.updateStateObject(so)
	}

	// Write the account trie changes, measuring the amount of wasted time
	if EnabledExpensive {
		defer func(start time.Time) { s.AccountCommits += time.Since(start) }(time.Now())
	}
	root, err = s.trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent common.ExtHash, parentDepth int) error {
		serializer := account.NewAccountSerializer()
		if err := rlp.DecodeBytes(leaf, serializer); err != nil {
			logger.Warn("RLP decode failed", "err", err, "leaf", string(leaf))
			return nil
		}
		acc := serializer.GetAccount()
		if pa := account.GetProgramAccount(acc); pa != nil {
			if pa.GetStorageRoot().Unextend() != emptyState {
				s.db.TrieDB().Reference(pa.GetStorageRoot(), parent)
			}
		}
		return nil
	})

	// If snapshotting is enabled, update the snapshot tree with this new version
	if s.snap != nil {
		if EnabledExpensive {
			defer func(start time.Time) { s.SnapshotCommits += time.Since(start) }(time.Now())
		}
		// Only update if there's a state transition (skip empty Clique blocks)
		if parent := s.snap.Root(); parent != root {
			if err := s.snaps.Update(root, parent, s.snapDestructs, s.snapAccounts, s.snapStorage); err != nil {
				logger.Warn("Failed to update snapshot tree", "from", parent, "to", root, "err", err)
			}
			// Keep 128 diff layers in the memory, persistent layer is 129th.
			// - head layer is paired with HEAD state
			// - head-1 layer is paired with HEAD-1 state
			// - head-127 layer(bottom-most diff layer) is paired with HEAD-127 state
			if err := s.snaps.Cap(root, 128); err != nil {
				logger.Warn("Failed to cap snapshot tree", "root", root, "layers", 128, "err", err)
			}
		}
		s.snap, s.snapDestructs, s.snapAccounts, s.snapStorage = nil, nil, nil, nil
	}

	return root, err
}

// GetTxHash returns the hash of current running transaction.
func (s *StateDB) GetTxHash() common.Hash {
	return s.thash
}

var (
	errNotExistingAddress = fmt.Errorf("there is no account corresponding to the given address")
	errNotContractAddress = fmt.Errorf("given address is not a contract address")
)

func (s *StateDB) GetContractStorageRoot(contractAddr common.Address) (common.ExtHash, error) {
	acc := s.GetAccount(contractAddr)
	if acc == nil {
		return common.ExtHash{}, errNotExistingAddress
	}
	if acc.Type() != account.SmartContractAccountType {
		return common.ExtHash{}, errNotContractAddress
	}
	contract, true := acc.(*account.SmartContractAccount)
	if !true {
		return common.ExtHash{}, errNotContractAddress
	}
	return contract.GetStorageRoot(), nil
}

// PrepareAccessList handles the preparatory steps for executing a state transition with
// regards to EIP-2929:
//
// - Add sender to access list (2929)
// - Add feepayer to access list (only for klaytn)
// - Add destination to access list (2929)
// - Add precompiles to access list (2929)
//
// regards to EIP-3651:
// - Add coinbase to access list (EIP-3651)
func (s *StateDB) PrepareAccessList(rules params.Rules, sender, feepayer, coinbase common.Address, dst *common.Address, precompiles []common.Address, list types.AccessList) {
	if rules.IsKore {
		// Clear out any leftover from previous executions
		s.accessList = newAccessList()

		s.AddAddressToAccessList(sender)
		if !common.EmptyAddress(feepayer) {
			s.AddAddressToAccessList(feepayer)
		}
		if dst != nil {
			s.AddAddressToAccessList(*dst)
			// If it's a create-tx, the destination will be added inside evm.create
		}
		for _, addr := range precompiles {
			s.AddAddressToAccessList(addr)
		}

		if rules.IsCancun {
			// Optional accessList is the accessList mentioned through tx args.
			for _, el := range list {
				s.AddAddressToAccessList(el.Address)
				for _, key := range el.StorageKeys {
					s.AddSlotToAccessList(el.Address, key)
				}
			}
		}

		if rules.IsShanghai {
			s.AddAddressToAccessList(coinbase)
		}
	}
}

// AddAddressToAccessList adds the given address to the access list
func (s *StateDB) AddAddressToAccessList(addr common.Address) {
	if s.accessList.AddAddress(addr) {
		s.journal.append(accessListAddAccountChange{&addr})
	}
}

// AddSlotToAccessList adds the given (address, slot)-tuple to the access list
func (s *StateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	addrMod, slotMod := s.accessList.AddSlot(addr, slot)
	if addrMod {
		// In practice, this should not happen, since there is no way to enter the
		// scope of 'address' without having the 'address' become already added
		// to the access list (via call-variant, create, etc).
		// Better safe than sorry, though
		s.journal.append(accessListAddAccountChange{&addr})
	}
	if slotMod {
		s.journal.append(accessListAddSlotChange{
			address: &addr,
			slot:    &slot,
		})
	}
}

// AddressInAccessList returns true if the given address is in the access list.
func (s *StateDB) AddressInAccessList(addr common.Address) bool {
	return s.accessList.ContainsAddress(addr)
}

// SlotInAccessList returns true if the given (address, slot)-tuple is in the access list.
func (s *StateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressPresent bool, slotPresent bool) {
	return s.accessList.Contains(addr, slot)
}
