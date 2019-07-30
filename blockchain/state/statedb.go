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
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/statedb"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
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
)

// TODO-Klaytn-StateDB Need to consider setting this value by command line option.
const maxCachedStateObjects = 40960

// StateDBs within the Klaytn protocol are used to cache stateObjects from Merkle Patricia Trie
// and mediates the operations to them.
type StateDB struct {
	db   Database
	trie Trie

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateObjects      map[common.Address]*stateObject
	stateObjectsDirty map[common.Address]struct{}

	// cachedStateObjects stores the most recent finalized stateObjects.
	cachedStateObjects common.Cache

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

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionId int

	lock sync.Mutex
}

// NewCachedStateObjects returns a new Common.Cache object for cachedStateObjects.
func NewCachedStateObjects() common.Cache {
	return common.NewCache(common.LRUConfig{CacheSize: maxCachedStateObjects})
}

// Create a new state from a given trie.
func New(root common.Hash, db Database) (*StateDB, error) {
	tr, err := db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	return &StateDB{
		db:                 db,
		trie:               tr,
		stateObjects:       make(map[common.Address]*stateObject),
		stateObjectsDirty:  make(map[common.Address]struct{}),
		cachedStateObjects: nil,
		logs:               make(map[common.Hash][]*types.Log),
		preimages:          make(map[common.Hash][]byte),
		journal:            newJournal(),
	}, nil
}

// NewWithCache creates a new state from a given trie with state object caching enabled.
func NewWithCache(root common.Hash, db Database, cachedStateObjects common.Cache) (*StateDB, error) {
	if stateDB, err := New(root, db); err != nil {
		return nil, err
	} else {
		stateDB.cachedStateObjects = cachedStateObjects
		return stateDB, nil
	}
}

// setError remembers the first non-nil error it is called with.
func (self *StateDB) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *StateDB) Error() error {
	return self.dbErr
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying state trie to avoid reloading data for the next operations.
func (self *StateDB) Reset(root common.Hash) error {
	tr, err := self.db.OpenTrie(root)
	if err != nil {
		return err
	}
	self.trie = tr
	self.stateObjects = make(map[common.Address]*stateObject)
	self.stateObjectsDirty = make(map[common.Address]struct{})
	self.cachedStateObjects = NewCachedStateObjects()
	self.thash = common.Hash{}
	self.bhash = common.Hash{}
	self.txIndex = 0
	self.logs = make(map[common.Hash][]*types.Log)
	self.logSize = 0
	self.preimages = make(map[common.Hash][]byte)
	self.clearJournalAndRefund()
	return nil
}

// ResetExceptCachedStateObjects clears out all ephemeral state objects from the state db,
// but keeps 1) underlying state trie and 2) cachedStateObjects to avoid reloading data.
func (self *StateDB) ResetExceptCachedStateObjects(root common.Hash) {
	cachedStateObjects := self.cachedStateObjects
	self.Reset(root)
	self.cachedStateObjects = cachedStateObjects
}

// UpdateTxPoolStateCache updates the caches used in TxPool.
func (self *StateDB) UpdateTxPoolStateCache(nonceCache common.Cache, balanceCache common.Cache) {
	if nonceCache == nil || balanceCache == nil {
		logger.ErrorWithStack("UpdateTxPoolStateCache should not be called! If nonceCache or balanceCache is nil")
		return
	}

	for addr, stateObj := range self.stateObjects {
		if stateObj.suicided || stateObj.deleted || stateObj.dbErr != nil {
			nonceCache.Add(addr, nil)
			balanceCache.Add(addr, nil)
		} else {
			nonceCache.Add(addr, stateObj.Nonce())
			balanceCache.Add(addr, stateObj.Balance())
		}
	}
}

// UpdateCachedStateObjects copies stateObjects to cachedStateObjects.
// And then, call ResetExceptCachedStateObjects() to clear out all ephemeral state objects,
// except for 1) underlying state trie and 2) cachedStateObjects.
func (self *StateDB) UpdateCachedStateObjects(root common.Hash) {
	if !self.UseCachedStateObjects() {
		logger.ErrorWithStack("UpdateCachedStateObjects should not be called! It is disabled!")
		return
	}
	for addr, stateObj := range self.stateObjects {
		if stateObj.suicided || stateObj.deleted || stateObj.dbErr != nil {
			self.cachedStateObjects.Add(addr, nil)
		} else {
			// TODO-Klaytn-StateDB cachedStorage also can be reused, but needs to be tested.
			newObj := newObject(self, addr, stateObj.account)
			newObj.storageTrie = stateObj.storageTrie
			newObj.code = stateObj.code
			self.cachedStateObjects.Add(addr, newObj)
		}
	}

	// Reset all member variables except for cachedStateObjects.
	self.ResetExceptCachedStateObjects(root)
}

// GetCachedStateObjects returns its cachedStateObjects.
func (self *StateDB) GetCachedStateObjects() common.Cache {
	return self.cachedStateObjects
}

// UseCachedStateObjects returns if it uses cachedStateObjects or not.
func (self *StateDB) UseCachedStateObjects() bool {
	return self.cachedStateObjects != nil
}

func (self *StateDB) AddLog(log *types.Log) {
	self.journal.append(addLogChange{txhash: self.thash})

	log.TxHash = self.thash
	log.BlockHash = self.bhash
	log.TxIndex = uint(self.txIndex)
	log.Index = self.logSize
	self.logs[self.thash] = append(self.logs[self.thash], log)
	self.logSize++
}

func (self *StateDB) GetLogs(hash common.Hash) []*types.Log {
	return self.logs[hash]
}

func (self *StateDB) Logs() []*types.Log {
	var logs []*types.Log
	for _, lgs := range self.logs {
		logs = append(logs, lgs...)
	}
	return logs
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (self *StateDB) AddPreimage(hash common.Hash, preimage []byte) {
	if _, ok := self.preimages[hash]; !ok {
		self.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		self.preimages[hash] = pi
	}
}

// Preimages returns a list of SHA3 preimages that have been submitted.
func (self *StateDB) Preimages() map[common.Hash][]byte {
	return self.preimages
}

func (self *StateDB) AddRefund(gas uint64) {
	self.journal.append(refundChange{prev: self.refund})
	self.refund += gas
}

// Exist reports whether the given account address exists in the state.
// Notably this also returns true for suicided accounts.
func (self *StateDB) Exist(addr common.Address) bool {
	return self.getStateObject(addr) != nil
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (self *StateDB) Empty(addr common.Address) bool {
	so := self.getStateObject(addr)
	return so == nil || so.empty()
}

// Retrieve the balance from the given address or 0 if object not found
func (self *StateDB) GetBalance(addr common.Address) *big.Int {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Balance()
	}
	return common.Big0
}

func (self *StateDB) GetNonce(addr common.Address) uint64 {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}
	return 0
}

func (self *StateDB) GetCode(addr common.Address) []byte {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Code(self.db)
	}
	return nil
}

func (self *StateDB) GetAccount(addr common.Address) account.Account {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.account
	}
	return nil
}

func (self *StateDB) IsContractAccount(addr common.Address) bool {
	stateObject := self.getStateObject(addr)
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

func (self *StateDB) GetCodeSize(addr common.Address) int {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return 0
	}
	if stateObject.code != nil {
		return len(stateObject.code)
	}
	size, err := self.db.ContractCodeSize(common.BytesToHash(stateObject.CodeHash()))
	if err != nil {
		self.setError(err)
	}
	return size
}

func (self *StateDB) GetCodeHash(addr common.Address) common.Hash {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return common.BytesToHash(emptyCodeHash)
	}
	return common.BytesToHash(stateObject.CodeHash())
}

func (self *StateDB) GetState(addr common.Address, bhash common.Hash) common.Hash {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetState(self.db, bhash)
	}
	return common.Hash{}
}

// IsContractAvailable returns true if the account corresponding to the given address implements ProgramAccount.
func (self *StateDB) IsContractAvailable(addr common.Address) bool {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.IsContractAvailable()
	}
	return false
}

// IsProgramAccount returns true if the account corresponding to the given address implements ProgramAccount.
func (self *StateDB) IsProgramAccount(addr common.Address) bool {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.IsProgramAccount()
	}
	return false
}

func (self *StateDB) IsValidCodeFormat(addr common.Address) bool {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		pa := account.GetProgramAccount(stateObject.account)
		if pa != nil {
			return pa.GetCodeFormat().Validate()
		}
		return false
	}
	return false
}

func (self *StateDB) GetKey(addr common.Address) accountkey.AccountKey {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetKey()
	}
	return accountkey.NewAccountKeyLegacy()
}

// Database retrieves the low level database supporting the lower level trie ops.
func (self *StateDB) Database() Database {
	return self.db
}

// StorageTrie returns the storage trie of an account.
// The return value is a copy and is nil for non-existent accounts.
func (self *StateDB) StorageTrie(addr common.Address) Trie {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return nil
	}
	cpy := stateObject.deepCopy(self)
	return cpy.updateStorageTrie(self.db)
}

func (self *StateDB) HasSuicided(addr common.Address) bool {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.suicided
	}
	return false
}

/*
 * SETTERS
 */

// AddBalance adds amount to the account associated with addr.
func (self *StateDB) AddBalance(addr common.Address, amount *big.Int) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (self *StateDB) SubBalance(addr common.Address, amount *big.Int) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)
	}
}

func (self *StateDB) SetBalance(addr common.Address, amount *big.Int) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetBalance(amount)
	}
}

// IncNonce increases the nonce of the account of the given address by one.
func (self *StateDB) IncNonce(addr common.Address) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.IncNonce()
	}
}

func (self *StateDB) SetNonce(addr common.Address, nonce uint64) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

func (self *StateDB) SetCode(addr common.Address, code []byte) error {
	stateObject := self.GetOrNewSmartContract(addr)
	if stateObject != nil {
		return stateObject.SetCode(crypto.Keccak256Hash(code), code)
	}

	return nil
}

func (self *StateDB) SetState(addr common.Address, key, value common.Hash) {
	stateObject := self.GetOrNewSmartContract(addr)
	if stateObject != nil {
		stateObject.SetState(self.db, key, value)
	}
}

// UpdateKey updates the account's key with the given key.
func (self *StateDB) UpdateKey(addr common.Address, newKey accountkey.AccountKey, currentBlockNumber uint64) error {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.UpdateKey(newKey, currentBlockNumber)
	}

	return errAccountDoesNotExist
}

// Suicide marks the given account as suicided.
// This clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (self *StateDB) Suicide(addr common.Address) bool {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return false
	}
	self.journal.append(suicideChange{
		account:     &addr,
		prev:        stateObject.suicided,
		prevbalance: new(big.Int).Set(stateObject.Balance()),
	})
	stateObject.markSuicided()
	stateObject.account.SetBalance(new(big.Int))

	return true
}

//
// Setting, updating & deleting state object methods.
//

// updateStateObject writes the given object to the statedb.
func (self *StateDB) updateStateObject(stateObject *stateObject) {
	addr := stateObject.Address()
	if data := stateObject.encoded.Load(); data != nil {
		encodedData := data.(*encodedData)
		if encodedData.err != nil {
			panic(fmt.Errorf("can't encode object at %x: %v", addr[:], encodedData.err))
		}
		self.setError(self.trie.TryUpdateWithKeys(addr[:],
			encodedData.trieHashKey, encodedData.trieHexKey, encodedData.data))
		stateObject.encoded = atomic.Value{}
	} else {
		data, err := rlp.EncodeToBytes(stateObject)
		if err != nil {
			panic(fmt.Errorf("can't encode object at %x: %v", addr[:], err))
		}
		self.setError(self.trie.TryUpdate(addr[:], data))
	}
}

// deleteStateObject removes the given object from the state trie.
func (self *StateDB) deleteStateObject(stateObject *stateObject) {
	stateObject.deleted = true
	addr := stateObject.Address()
	self.setError(self.trie.TryDelete(addr[:]))
}

// Retrieve a state object given by the address. Returns nil if not found.
func (self *StateDB) getStateObject(addr common.Address) *stateObject {
	// First, check stateObjects if there is "live" object.
	if obj := self.stateObjects[addr]; obj != nil {
		if obj.deleted {
			return nil
		}
		return obj
	}

	// Second, if not found in stateObjects, check cachedStateObjects.
	if self.UseCachedStateObjects() {
		if obj, exist := self.cachedStateObjects.Get(addr); exist && obj != nil {
			stateObj, ok := obj.(*stateObject)
			if !ok {
				logger.Error("cached state object is not *stateObject type!", "addr", addr.String())
				return nil
			}

			self.stateObjects[addr] = stateObj.deepCopy(self)
			return self.stateObjects[addr]
		}
	}

	// Third, the object for given address is not cached.
	// Load the object from the database.
	enc, err := self.trie.TryGet(addr[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	serializer := account.NewAccountSerializer()
	if err := rlp.DecodeBytes(enc, serializer); err != nil {
		logger.Error("Failed to decode state object", "addr", addr, "err", err)
		return nil
	}
	data := serializer.GetAccount()
	// Insert into the live set.
	obj := newObject(self, addr, data)
	self.setStateObject(obj)

	if self.UseCachedStateObjects() {
		self.cachedStateObjects.Add(addr, obj.deepCopy(self))
	}

	return obj
}

func (self *StateDB) setStateObject(object *stateObject) {
	self.stateObjects[object.Address()] = object
}

// Retrieve a state object or create a new state object if nil.
func (self *StateDB) GetOrNewStateObject(addr common.Address) *stateObject {
	stateObject := self.getStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		stateObject, _ = self.createObject(addr)
	}
	return stateObject
}

// Retrieve a state object or create a new state object if nil.
func (self *StateDB) GetOrNewSmartContract(addr common.Address) *stateObject {
	stateObject := self.getStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		stateObject, _ = self.createObjectWithMap(addr, account.SmartContractAccountType, map[account.AccountValueKeyType]interface{}{
			account.AccountValueKeyNonce:         uint64(1),
			account.AccountValueKeyHumanReadable: false,
			account.AccountValueKeyAccountKey:    accountkey.NewAccountKeyFail(),
		})
	}
	return stateObject
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (self *StateDB) createObject(addr common.Address) (newobj, prev *stateObject) {
	prev = self.getStateObject(addr)
	acc, err := account.NewAccountWithType(account.ExternallyOwnedAccountType)
	if err != nil {
		logger.Error("An error occurred on call NewAccountWithType", "err", err)
	}
	newobj = newObject(self, addr, acc)
	newobj.setNonce(0) // sets the object to dirty
	if prev == nil {
		self.journal.append(createObjectChange{account: &addr})
	} else {
		self.journal.append(resetObjectChange{prev: prev})
	}
	self.setStateObject(newobj)
	return newobj, prev
}

// createObjectWithMap creates a new state object with the given parameters (accountType and values).
// If there is an existing account with the given address, it is overwritten and
// returned as the second return value.
func (self *StateDB) createObjectWithMap(addr common.Address, accountType account.AccountType,
	values map[account.AccountValueKeyType]interface{}) (newobj, prev *stateObject) {
	prev = self.getStateObject(addr)
	acc, err := account.NewAccountWithMap(accountType, values)
	if err != nil {
		logger.Error("An error occurred on call NewAccountWithMap", "err", err)
	}
	newobj = newObject(self, addr, acc)
	newobj.setNonce(0) // sets the object to dirty
	if prev == nil {
		self.journal.append(createObjectChange{account: &addr})
	} else {
		self.journal.append(resetObjectChange{prev: prev})
	}
	self.setStateObject(newobj)
	return newobj, prev
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might arise that
// a contract does the following:
//
//   1. sends funds to sha(account ++ (nonce + 1))
//   2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (self *StateDB) CreateAccount(addr common.Address) {
	new, prev := self.createObject(addr)
	if prev != nil {
		new.setBalance(prev.account.GetBalance())
	}
}

func (self *StateDB) CreateEOA(addr common.Address, humanReadable bool, key accountkey.AccountKey) {
	values := map[account.AccountValueKeyType]interface{}{
		account.AccountValueKeyHumanReadable: humanReadable,
		account.AccountValueKeyAccountKey:    key,
	}
	new, prev := self.createObjectWithMap(addr, account.ExternallyOwnedAccountType, values)
	if prev != nil {
		new.setBalance(prev.account.GetBalance())
	}
}

func (self *StateDB) CreateSmartContractAccount(addr common.Address, format params.CodeFormat) {
	self.CreateSmartContractAccountWithKey(addr, false, accountkey.NewAccountKeyFail(), format)
}

func (self *StateDB) CreateSmartContractAccountWithKey(addr common.Address, humanReadable bool, key accountkey.AccountKey, format params.CodeFormat) {
	values := map[account.AccountValueKeyType]interface{}{
		account.AccountValueKeyNonce:         uint64(1),
		account.AccountValueKeyHumanReadable: humanReadable,
		account.AccountValueKeyAccountKey:    key,
		account.AccountValueKeyCodeFormat:    format,
	}
	new, prev := self.createObjectWithMap(addr, account.SmartContractAccountType, values)
	if prev != nil {
		new.setBalance(prev.account.GetBalance())
	}
}

func (db *StateDB) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) {
	so := db.getStateObject(addr)
	if so == nil {
		return
	}

	// When iterating over the storage check the cache first
	for h, value := range so.cachedStorage {
		cb(h, value)
	}

	it := statedb.NewIterator(so.getStorageTrie(db.db).NodeIterator(nil))
	for it.Next() {
		// ignore cached values
		key := common.BytesToHash(db.trie.GetKey(it.Key))
		if _, ok := so.cachedStorage[key]; !ok {
			cb(key, common.BytesToHash(it.Value))
		}
	}
}

// Copy creates a deep, independent copy of the state.
// Snapshots of the copied state cannot be applied to the copy.
func (self *StateDB) Copy() *StateDB {
	self.lock.Lock()
	defer self.lock.Unlock()

	// Copy all the basic fields, initialize the memory ones
	state := &StateDB{
		db:                 self.db,
		trie:               self.db.CopyTrie(self.trie),
		stateObjects:       make(map[common.Address]*stateObject, len(self.journal.dirties)),
		stateObjectsDirty:  make(map[common.Address]struct{}, len(self.journal.dirties)),
		cachedStateObjects: nil,
		refund:             self.refund,
		logs:               make(map[common.Hash][]*types.Log, len(self.logs)),
		logSize:            self.logSize,
		preimages:          make(map[common.Hash][]byte),
		journal:            newJournal(),
	}
	// Copy the dirty states, logs, and preimages
	for addr := range self.journal.dirties {
		// As documented [here](https://github.com/ethereum/go-ethereum/pull/16485#issuecomment-380438527),
		// and in the Finalise-method, there is a case where an object is in the journal but not
		// in the stateObjects: OOG after touch on ripeMD prior to Byzantium. Thus, we need to check for
		// nil
		if object, exist := self.stateObjects[addr]; exist {
			state.stateObjects[addr] = object.deepCopy(state)
			state.stateObjectsDirty[addr] = struct{}{}
		}
	}
	// Above, we don't copy the actual journal. This means that if the copy is copied, the
	// loop above will be a no-op, since the copy's journal is empty.
	// Thus, here we iterate over stateObjects, to enable copies of copies
	for addr := range self.stateObjectsDirty {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = self.stateObjects[addr].deepCopy(state)
			state.stateObjectsDirty[addr] = struct{}{}
		}
	}

	// NOTE-Klaytn-StateDB cachedStateObject is cache, so not copied to copied StateDB.

	deepCopyLogs(self, state)

	for hash, preimage := range self.preimages {
		state.preimages[hash] = preimage
	}

	// Use cachedStateObjects only if original StateDB uses.
	// However, cached values are not copied.
	if self.cachedStateObjects != nil {
		state.cachedStateObjects = NewCachedStateObjects()
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
func (self *StateDB) Snapshot() int {
	id := self.nextRevisionId
	self.nextRevisionId++
	self.validRevisions = append(self.validRevisions, revision{id, self.journal.length()})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (self *StateDB) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(self.validRevisions), func(i int) bool {
		return self.validRevisions[i].id >= revid
	})
	if idx == len(self.validRevisions) || self.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := self.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	self.journal.revert(self, snapshot)
	self.validRevisions = self.validRevisions[:idx]
}

// GetRefund returns the current value of the refund counter.
func (self *StateDB) GetRefund() uint64 {
	return self.refund
}

// Finalise finalises the state by removing the self destructed objects
// and clears the journal as well as the refunds.
func (s *StateDB) Finalise(deleteEmptyObjects bool) {
	for addr := range s.journal.dirties {
		stateObject, exist := s.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `s.journal.dirties` but not in `s.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}

		if stateObject.suicided || (deleteEmptyObjects && stateObject.empty()) {
			s.deleteStateObject(stateObject)
		} else {
			stateObject.updateStorageRoot(s.db)
			s.updateStateObject(stateObject)
		}
		s.stateObjectsDirty[addr] = struct{}{}
	}
	// Invalidate journal because reverting across transactions is not allowed.
	s.clearJournalAndRefund()
}

// IntermediateRoot computes the current root hash of the state statedb.
// It is called in between transactions to get the root hash that
// goes into transaction receipts.
func (s *StateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	s.Finalise(deleteEmptyObjects)
	return s.trie.Hash()
}

// Prepare sets the current transaction hash and index and block hash which is
// used when the EVM emits new state logs.
func (self *StateDB) Prepare(thash, bhash common.Hash, ti int) {
	self.thash = thash
	self.bhash = bhash
	self.txIndex = ti
}

func (s *StateDB) clearJournalAndRefund() {
	s.journal = newJournal()
	s.validRevisions = s.validRevisions[:0]
	s.refund = 0
}

// Commit writes the state to the underlying in-memory trie database.
func (s *StateDB) Commit(deleteEmptyObjects bool) (root common.Hash, err error) {
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
		case stateObject.suicided || (isDirty && deleteEmptyObjects && stateObject.empty()):
			// If the object has been removed, don't bother syncing it
			// and just mark it for deletion in the trie.
			s.deleteStateObject(stateObject)
		case isDirty:
			if stateObject.IsProgramAccount() {
				// Write any contract code associated with the state object.
				if stateObject.code != nil && stateObject.dirtyCode {
					s.db.TrieDB().InsertBlob(common.BytesToHash(stateObject.CodeHash()), stateObject.code)
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

	// Write trie changes.
	root, err = s.trie.Commit(func(leaf []byte, parent common.Hash) error {
		serializer := account.NewAccountSerializer()
		if err := rlp.DecodeBytes(leaf, serializer); err != nil {
			logger.Warn("RLP decode failed", "err", err, "leaf", string(leaf))
			return nil
		}
		acc := serializer.GetAccount()
		if pa := account.GetProgramAccount(acc); pa != nil {
			if pa.GetStorageRoot() != emptyState {
				s.db.TrieDB().Reference(pa.GetStorageRoot(), parent)
			}
			code := common.BytesToHash(pa.GetCodeHash())
			if code != emptyCode {
				s.db.TrieDB().Reference(code, parent)
			}
		}
		return nil
	})
	return root, err
}

// GetTxHash returns the hash of current running transaction.
func (s *StateDB) GetTxHash() common.Hash {
	return s.thash
}
