// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/klaytn/klaytn/api (interfaces: Backend)

// Package mock_api is a generated GoMock package.
package mock_api

import (
	context "context"
	big "math/big"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	klaytn "github.com/klaytn/klaytn"
	accounts "github.com/klaytn/klaytn/accounts"
	blockchain "github.com/klaytn/klaytn/blockchain"
	state "github.com/klaytn/klaytn/blockchain/state"
	types "github.com/klaytn/klaytn/blockchain/types"
	vm "github.com/klaytn/klaytn/blockchain/vm"
	common "github.com/klaytn/klaytn/common"
	consensus "github.com/klaytn/klaytn/consensus"
	event "github.com/klaytn/klaytn/event"
	rpc "github.com/klaytn/klaytn/networks/rpc"
	params "github.com/klaytn/klaytn/params"
	database "github.com/klaytn/klaytn/storage/database"
)

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// AccountManager mocks base method.
func (m *MockBackend) AccountManager() accounts.AccountManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccountManager")
	ret0, _ := ret[0].(accounts.AccountManager)
	return ret0
}

// AccountManager indicates an expected call of AccountManager.
func (mr *MockBackendMockRecorder) AccountManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccountManager", reflect.TypeOf((*MockBackend)(nil).AccountManager))
}

// BlockByHash mocks base method.
func (m *MockBackend) BlockByHash(arg0 context.Context, arg1 common.Hash) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockByHash", arg0, arg1)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockByHash indicates an expected call of BlockByHash.
func (mr *MockBackendMockRecorder) BlockByHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockByHash", reflect.TypeOf((*MockBackend)(nil).BlockByHash), arg0, arg1)
}

// BlockByNumber mocks base method.
func (m *MockBackend) BlockByNumber(arg0 context.Context, arg1 rpc.BlockNumber) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockByNumber", arg0, arg1)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockByNumber indicates an expected call of BlockByNumber.
func (mr *MockBackendMockRecorder) BlockByNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockByNumber", reflect.TypeOf((*MockBackend)(nil).BlockByNumber), arg0, arg1)
}

// BlockByNumberOrHash mocks base method.
func (m *MockBackend) BlockByNumberOrHash(arg0 context.Context, arg1 rpc.BlockNumberOrHash) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockByNumberOrHash", arg0, arg1)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockByNumberOrHash indicates an expected call of BlockByNumberOrHash.
func (mr *MockBackendMockRecorder) BlockByNumberOrHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockByNumberOrHash", reflect.TypeOf((*MockBackend)(nil).BlockByNumberOrHash), arg0, arg1)
}

// ChainConfig mocks base method.
func (m *MockBackend) ChainConfig() *params.ChainConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChainConfig")
	ret0, _ := ret[0].(*params.ChainConfig)
	return ret0
}

// ChainConfig indicates an expected call of ChainConfig.
func (mr *MockBackendMockRecorder) ChainConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChainConfig", reflect.TypeOf((*MockBackend)(nil).ChainConfig))
}

// ChainDB mocks base method.
func (m *MockBackend) ChainDB() database.DBManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChainDB")
	ret0, _ := ret[0].(database.DBManager)
	return ret0
}

// ChainDB indicates an expected call of ChainDB.
func (mr *MockBackendMockRecorder) ChainDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChainDB", reflect.TypeOf((*MockBackend)(nil).ChainDB))
}

// CurrentBlock mocks base method.
func (m *MockBackend) CurrentBlock() *types.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentBlock")
	ret0, _ := ret[0].(*types.Block)
	return ret0
}

// CurrentBlock indicates an expected call of CurrentBlock.
func (mr *MockBackendMockRecorder) CurrentBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentBlock", reflect.TypeOf((*MockBackend)(nil).CurrentBlock))
}

func (m *MockBackend) Engine() consensus.Engine {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Engine")
	ret0, _ := ret[0].(consensus.Engine)
	return ret0
}

// Engine indicates an expected call of Engine.
func (mr *MockBackendMockRecorder) Engine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Engine", reflect.TypeOf((*MockBackend)(nil).Engine))
}

// EventMux mocks base method.
func (m *MockBackend) EventMux() *event.TypeMux {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventMux")
	ret0, _ := ret[0].(*event.TypeMux)
	return ret0
}

// EventMux indicates an expected call of EventMux.
func (mr *MockBackendMockRecorder) EventMux() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventMux", reflect.TypeOf((*MockBackend)(nil).EventMux))
}

// FeeHistory mocks base method.
func (m *MockBackend) FeeHistory(arg0 context.Context, arg1 int, arg2 rpc.BlockNumber, arg3 []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FeeHistory", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].([][]*big.Int)
	ret2, _ := ret[2].([]*big.Int)
	ret3, _ := ret[3].([]float64)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// FeeHistory indicates an expected call of FeeHistory.
func (mr *MockBackendMockRecorder) FeeHistory(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FeeHistory", reflect.TypeOf((*MockBackend)(nil).FeeHistory), arg0, arg1, arg2, arg3)
}

// GetBlockReceipts mocks base method.
func (m *MockBackend) GetBlockReceipts(arg0 context.Context, arg1 common.Hash) types.Receipts {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockReceipts", arg0, arg1)
	ret0, _ := ret[0].(types.Receipts)
	return ret0
}

// GetBlockReceipts indicates an expected call of GetBlockReceipts.
func (mr *MockBackendMockRecorder) GetBlockReceipts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockReceipts", reflect.TypeOf((*MockBackend)(nil).GetBlockReceipts), arg0, arg1)
}

// GetBlockReceiptsInCache mocks base method.
func (m *MockBackend) GetBlockReceiptsInCache(arg0 common.Hash) types.Receipts {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockReceiptsInCache", arg0)
	ret0, _ := ret[0].(types.Receipts)
	return ret0
}

// GetBlockReceiptsInCache indicates an expected call of GetBlockReceiptsInCache.
func (mr *MockBackendMockRecorder) GetBlockReceiptsInCache(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockReceiptsInCache", reflect.TypeOf((*MockBackend)(nil).GetBlockReceiptsInCache), arg0)
}

// GetEVM mocks base method.
func (m *MockBackend) GetEVM(arg0 context.Context, arg1 blockchain.Message, arg2 *state.StateDB, arg3 *types.Header, arg4 vm.Config) (*vm.EVM, func() error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEVM", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*vm.EVM)
	ret1, _ := ret[1].(func() error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetEVM indicates an expected call of GetEVM.
func (mr *MockBackendMockRecorder) GetEVM(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEVM", reflect.TypeOf((*MockBackend)(nil).GetEVM), arg0, arg1, arg2, arg3, arg4)
}

// GetPoolNonce mocks base method.
func (m *MockBackend) GetPoolNonce(arg0 context.Context, arg1 common.Address) uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPoolNonce", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetPoolNonce indicates an expected call of GetPoolNonce.
func (mr *MockBackendMockRecorder) GetPoolNonce(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPoolNonce", reflect.TypeOf((*MockBackend)(nil).GetPoolNonce), arg0, arg1)
}

// GetPoolTransaction mocks base method.
func (m *MockBackend) GetPoolTransaction(arg0 common.Hash) *types.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPoolTransaction", arg0)
	ret0, _ := ret[0].(*types.Transaction)
	return ret0
}

// GetPoolTransaction indicates an expected call of GetPoolTransaction.
func (mr *MockBackendMockRecorder) GetPoolTransaction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPoolTransaction", reflect.TypeOf((*MockBackend)(nil).GetPoolTransaction), arg0)
}

// GetPoolTransactions mocks base method.
func (m *MockBackend) GetPoolTransactions() (types.Transactions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPoolTransactions")
	ret0, _ := ret[0].(types.Transactions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPoolTransactions indicates an expected call of GetPoolTransactions.
func (mr *MockBackendMockRecorder) GetPoolTransactions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPoolTransactions", reflect.TypeOf((*MockBackend)(nil).GetPoolTransactions))
}

// GetTd mocks base method.
func (m *MockBackend) GetTd(arg0 common.Hash) *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTd", arg0)
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// GetTd indicates an expected call of GetTd.
func (mr *MockBackendMockRecorder) GetTd(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTd", reflect.TypeOf((*MockBackend)(nil).GetTd), arg0)
}

// GetTxAndLookupInfo mocks base method.
func (m *MockBackend) GetTxAndLookupInfo(arg0 common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxAndLookupInfo", arg0)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(common.Hash)
	ret2, _ := ret[2].(uint64)
	ret3, _ := ret[3].(uint64)
	return ret0, ret1, ret2, ret3
}

// GetTxAndLookupInfo indicates an expected call of GetTxAndLookupInfo.
func (mr *MockBackendMockRecorder) GetTxAndLookupInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxAndLookupInfo", reflect.TypeOf((*MockBackend)(nil).GetTxAndLookupInfo), arg0)
}

// GetTxAndLookupInfoInCache mocks base method.
func (m *MockBackend) GetTxAndLookupInfoInCache(arg0 common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxAndLookupInfoInCache", arg0)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(common.Hash)
	ret2, _ := ret[2].(uint64)
	ret3, _ := ret[3].(uint64)
	return ret0, ret1, ret2, ret3
}

// GetTxAndLookupInfoInCache indicates an expected call of GetTxAndLookupInfoInCache.
func (mr *MockBackendMockRecorder) GetTxAndLookupInfoInCache(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxAndLookupInfoInCache", reflect.TypeOf((*MockBackend)(nil).GetTxAndLookupInfoInCache), arg0)
}

// GetTxLookupInfoAndReceipt mocks base method.
func (m *MockBackend) GetTxLookupInfoAndReceipt(arg0 context.Context, arg1 common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxLookupInfoAndReceipt", arg0, arg1)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(common.Hash)
	ret2, _ := ret[2].(uint64)
	ret3, _ := ret[3].(uint64)
	ret4, _ := ret[4].(*types.Receipt)
	return ret0, ret1, ret2, ret3, ret4
}

// GetTxLookupInfoAndReceipt indicates an expected call of GetTxLookupInfoAndReceipt.
func (mr *MockBackendMockRecorder) GetTxLookupInfoAndReceipt(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxLookupInfoAndReceipt", reflect.TypeOf((*MockBackend)(nil).GetTxLookupInfoAndReceipt), arg0, arg1)
}

// GetTxLookupInfoAndReceiptInCache mocks base method.
func (m *MockBackend) GetTxLookupInfoAndReceiptInCache(arg0 common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxLookupInfoAndReceiptInCache", arg0)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(common.Hash)
	ret2, _ := ret[2].(uint64)
	ret3, _ := ret[3].(uint64)
	ret4, _ := ret[4].(*types.Receipt)
	return ret0, ret1, ret2, ret3, ret4
}

// GetTxLookupInfoAndReceiptInCache indicates an expected call of GetTxLookupInfoAndReceiptInCache.
func (mr *MockBackendMockRecorder) GetTxLookupInfoAndReceiptInCache(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxLookupInfoAndReceiptInCache", reflect.TypeOf((*MockBackend)(nil).GetTxLookupInfoAndReceiptInCache), arg0)
}

// HeaderByHash mocks base method.
func (m *MockBackend) HeaderByHash(arg0 context.Context, arg1 common.Hash) (*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeaderByHash", arg0, arg1)
	ret0, _ := ret[0].(*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeaderByHash indicates an expected call of HeaderByHash.
func (mr *MockBackendMockRecorder) HeaderByHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeaderByHash", reflect.TypeOf((*MockBackend)(nil).HeaderByHash), arg0, arg1)
}

// HeaderByNumber mocks base method.
func (m *MockBackend) HeaderByNumber(arg0 context.Context, arg1 rpc.BlockNumber) (*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeaderByNumber", arg0, arg1)
	ret0, _ := ret[0].(*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeaderByNumber indicates an expected call of HeaderByNumber.
func (mr *MockBackendMockRecorder) HeaderByNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeaderByNumber", reflect.TypeOf((*MockBackend)(nil).HeaderByNumber), arg0, arg1)
}

// HeaderByNumberOrHash mocks base method.
func (m *MockBackend) HeaderByNumberOrHash(arg0 context.Context, arg1 rpc.BlockNumberOrHash) (*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeaderByNumberOrHash", arg0, arg1)
	ret0, _ := ret[0].(*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeaderByNumberOrHash indicates an expected call of HeaderByNumberOrHash.
func (mr *MockBackendMockRecorder) HeaderByNumberOrHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeaderByNumberOrHash", reflect.TypeOf((*MockBackend)(nil).HeaderByNumberOrHash), arg0, arg1)
}

// IsParallelDBWrite mocks base method.
func (m *MockBackend) IsParallelDBWrite() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsParallelDBWrite")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsParallelDBWrite indicates an expected call of IsParallelDBWrite.
func (mr *MockBackendMockRecorder) IsParallelDBWrite() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsParallelDBWrite", reflect.TypeOf((*MockBackend)(nil).IsParallelDBWrite))
}

// IsSenderTxHashIndexingEnabled mocks base method.
func (m *MockBackend) IsSenderTxHashIndexingEnabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSenderTxHashIndexingEnabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSenderTxHashIndexingEnabled indicates an expected call of IsSenderTxHashIndexingEnabled.
func (mr *MockBackendMockRecorder) IsSenderTxHashIndexingEnabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSenderTxHashIndexingEnabled", reflect.TypeOf((*MockBackend)(nil).IsSenderTxHashIndexingEnabled))
}

// Progress mocks base method.
func (m *MockBackend) Progress() klaytn.SyncProgress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Progress")
	ret0, _ := ret[0].(klaytn.SyncProgress)
	return ret0
}

// Progress indicates an expected call of Progress.
func (mr *MockBackendMockRecorder) Progress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Progress", reflect.TypeOf((*MockBackend)(nil).Progress))
}

// ProtocolVersion mocks base method.
func (m *MockBackend) ProtocolVersion() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProtocolVersion")
	ret0, _ := ret[0].(int)
	return ret0
}

// ProtocolVersion indicates an expected call of ProtocolVersion.
func (mr *MockBackendMockRecorder) ProtocolVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProtocolVersion", reflect.TypeOf((*MockBackend)(nil).ProtocolVersion))
}

// RPCGasCap mocks base method.
func (m *MockBackend) RPCGasCap() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RPCGasCap")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// RPCGasCap indicates an expected call of RPCGasCap.
func (mr *MockBackendMockRecorder) RPCGasCap() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RPCGasCap", reflect.TypeOf((*MockBackend)(nil).RPCGasCap))
}

// RPCTxFeeCap mocks base method.
func (m *MockBackend) RPCTxFeeCap() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RPCTxFeeCap")
	ret0, _ := ret[0].(float64)
	return ret0
}

// RPCTxFeeCap indicates an expected call of RPCTxFeeCap.
func (mr *MockBackendMockRecorder) RPCTxFeeCap() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RPCTxFeeCap", reflect.TypeOf((*MockBackend)(nil).RPCTxFeeCap))
}

// SendTx mocks base method.
func (m *MockBackend) SendTx(arg0 context.Context, arg1 *types.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTx", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendTx indicates an expected call of SendTx.
func (mr *MockBackendMockRecorder) SendTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTx", reflect.TypeOf((*MockBackend)(nil).SendTx), arg0, arg1)
}

// SetHead mocks base method.
func (m *MockBackend) SetHead(arg0 uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetHead", arg0)
}

// SetHead indicates an expected call of SetHead.
func (mr *MockBackendMockRecorder) SetHead(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHead", reflect.TypeOf((*MockBackend)(nil).SetHead), arg0)
}

// StateAndHeaderByNumber mocks base method.
func (m *MockBackend) StateAndHeaderByNumber(arg0 context.Context, arg1 rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateAndHeaderByNumber", arg0, arg1)
	ret0, _ := ret[0].(*state.StateDB)
	ret1, _ := ret[1].(*types.Header)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// StateAndHeaderByNumber indicates an expected call of StateAndHeaderByNumber.
func (mr *MockBackendMockRecorder) StateAndHeaderByNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateAndHeaderByNumber", reflect.TypeOf((*MockBackend)(nil).StateAndHeaderByNumber), arg0, arg1)
}

// StateAndHeaderByNumberOrHash mocks base method.
func (m *MockBackend) StateAndHeaderByNumberOrHash(arg0 context.Context, arg1 rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateAndHeaderByNumberOrHash", arg0, arg1)
	ret0, _ := ret[0].(*state.StateDB)
	ret1, _ := ret[1].(*types.Header)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// StateAndHeaderByNumberOrHash indicates an expected call of StateAndHeaderByNumberOrHash.
func (mr *MockBackendMockRecorder) StateAndHeaderByNumberOrHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateAndHeaderByNumberOrHash", reflect.TypeOf((*MockBackend)(nil).StateAndHeaderByNumberOrHash), arg0, arg1)
}

// Stats mocks base method.
func (m *MockBackend) Stats() (int, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// Stats indicates an expected call of Stats.
func (mr *MockBackendMockRecorder) Stats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockBackend)(nil).Stats))
}

// SubscribeChainEvent mocks base method.
func (m *MockBackend) SubscribeChainEvent(arg0 chan<- blockchain.ChainEvent) event.Subscription {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeChainEvent", arg0)
	ret0, _ := ret[0].(event.Subscription)
	return ret0
}

// SubscribeChainEvent indicates an expected call of SubscribeChainEvent.
func (mr *MockBackendMockRecorder) SubscribeChainEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeChainEvent", reflect.TypeOf((*MockBackend)(nil).SubscribeChainEvent), arg0)
}

// SubscribeChainHeadEvent mocks base method.
func (m *MockBackend) SubscribeChainHeadEvent(arg0 chan<- blockchain.ChainHeadEvent) event.Subscription {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeChainHeadEvent", arg0)
	ret0, _ := ret[0].(event.Subscription)
	return ret0
}

// SubscribeChainHeadEvent indicates an expected call of SubscribeChainHeadEvent.
func (mr *MockBackendMockRecorder) SubscribeChainHeadEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeChainHeadEvent", reflect.TypeOf((*MockBackend)(nil).SubscribeChainHeadEvent), arg0)
}

// SubscribeChainSideEvent mocks base method.
func (m *MockBackend) SubscribeChainSideEvent(arg0 chan<- blockchain.ChainSideEvent) event.Subscription {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeChainSideEvent", arg0)
	ret0, _ := ret[0].(event.Subscription)
	return ret0
}

// SubscribeChainSideEvent indicates an expected call of SubscribeChainSideEvent.
func (mr *MockBackendMockRecorder) SubscribeChainSideEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeChainSideEvent", reflect.TypeOf((*MockBackend)(nil).SubscribeChainSideEvent), arg0)
}

// SubscribeNewTxsEvent mocks base method.
func (m *MockBackend) SubscribeNewTxsEvent(arg0 chan<- blockchain.NewTxsEvent) event.Subscription {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeNewTxsEvent", arg0)
	ret0, _ := ret[0].(event.Subscription)
	return ret0
}

// SubscribeNewTxsEvent indicates an expected call of SubscribeNewTxsEvent.
func (mr *MockBackendMockRecorder) SubscribeNewTxsEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeNewTxsEvent", reflect.TypeOf((*MockBackend)(nil).SubscribeNewTxsEvent), arg0)
}

// SuggestPrice mocks base method.
func (m *MockBackend) SuggestPrice(arg0 context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuggestPrice", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestPrice indicates an expected call of SuggestPrice.
func (mr *MockBackendMockRecorder) SuggestPrice(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestPrice", reflect.TypeOf((*MockBackend)(nil).SuggestPrice), arg0)
}

// TxPoolContent mocks base method.
func (m *MockBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxPoolContent")
	ret0, _ := ret[0].(map[common.Address]types.Transactions)
	ret1, _ := ret[1].(map[common.Address]types.Transactions)
	return ret0, ret1
}

// TxPoolContent indicates an expected call of TxPoolContent.
func (mr *MockBackendMockRecorder) TxPoolContent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxPoolContent", reflect.TypeOf((*MockBackend)(nil).TxPoolContent))
}
