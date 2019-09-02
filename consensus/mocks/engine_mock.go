// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/klaytn/klaytn/consensus (interfaces: Engine)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	state "github.com/klaytn/klaytn/blockchain/state"
	types "github.com/klaytn/klaytn/blockchain/types"
	common "github.com/klaytn/klaytn/common"
	consensus "github.com/klaytn/klaytn/consensus"
	rpc "github.com/klaytn/klaytn/networks/rpc"
	big "math/big"
	reflect "reflect"
)

// MockEngine is a mock of Engine interface
type MockEngine struct {
	ctrl     *gomock.Controller
	recorder *MockEngineMockRecorder
}

// MockEngineMockRecorder is the mock recorder for MockEngine
type MockEngineMockRecorder struct {
	mock *MockEngine
}

// NewMockEngine creates a new mock instance
func NewMockEngine(ctrl *gomock.Controller) *MockEngine {
	mock := &MockEngine{ctrl: ctrl}
	mock.recorder = &MockEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEngine) EXPECT() *MockEngineMockRecorder {
	return m.recorder
}

// APIs mocks base method
func (m *MockEngine) APIs(arg0 consensus.ChainReader) []rpc.API {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APIs", arg0)
	ret0, _ := ret[0].([]rpc.API)
	return ret0
}

// APIs indicates an expected call of APIs
func (mr *MockEngineMockRecorder) APIs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APIs", reflect.TypeOf((*MockEngine)(nil).APIs), arg0)
}

// Author mocks base method
func (m *MockEngine) Author(arg0 *types.Header) (common.Address, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Author", arg0)
	ret0, _ := ret[0].(common.Address)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Author indicates an expected call of Author
func (mr *MockEngineMockRecorder) Author(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Author", reflect.TypeOf((*MockEngine)(nil).Author), arg0)
}

// CalcBlockScore mocks base method
func (m *MockEngine) CalcBlockScore(arg0 consensus.ChainReader, arg1 uint64, arg2 *types.Header) *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalcBlockScore", arg0, arg1, arg2)
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// CalcBlockScore indicates an expected call of CalcBlockScore
func (mr *MockEngineMockRecorder) CalcBlockScore(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalcBlockScore", reflect.TypeOf((*MockEngine)(nil).CalcBlockScore), arg0, arg1, arg2)
}

// Finalize mocks base method
func (m *MockEngine) Finalize(arg0 consensus.ChainReader, arg1 *types.Header, arg2 *state.StateDB, arg3 []*types.Transaction, arg4 []*types.Receipt) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Finalize", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Finalize indicates an expected call of Finalize
func (mr *MockEngineMockRecorder) Finalize(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Finalize", reflect.TypeOf((*MockEngine)(nil).Finalize), arg0, arg1, arg2, arg3, arg4)
}

// Prepare mocks base method
func (m *MockEngine) Prepare(arg0 consensus.ChainReader, arg1 *types.Header) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Prepare indicates an expected call of Prepare
func (mr *MockEngineMockRecorder) Prepare(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockEngine)(nil).Prepare), arg0, arg1)
}

// Protocol mocks base method
func (m *MockEngine) Protocol() consensus.Protocol {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Protocol")
	ret0, _ := ret[0].(consensus.Protocol)
	return ret0
}

// Protocol indicates an expected call of Protocol
func (mr *MockEngineMockRecorder) Protocol() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Protocol", reflect.TypeOf((*MockEngine)(nil).Protocol))
}

// Seal mocks base method
func (m *MockEngine) Seal(arg0 consensus.ChainReader, arg1 *types.Block, arg2 <-chan struct{}) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seal", arg0, arg1, arg2)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seal indicates an expected call of Seal
func (mr *MockEngineMockRecorder) Seal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seal", reflect.TypeOf((*MockEngine)(nil).Seal), arg0, arg1, arg2)
}

// VerifyHeader mocks base method
func (m *MockEngine) VerifyHeader(arg0 consensus.ChainReader, arg1 *types.Header, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyHeader", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyHeader indicates an expected call of VerifyHeader
func (mr *MockEngineMockRecorder) VerifyHeader(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyHeader", reflect.TypeOf((*MockEngine)(nil).VerifyHeader), arg0, arg1, arg2)
}

// VerifyHeaders mocks base method
func (m *MockEngine) VerifyHeaders(arg0 consensus.ChainReader, arg1 []*types.Header, arg2 []bool) (chan<- struct{}, <-chan error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyHeaders", arg0, arg1, arg2)
	ret0, _ := ret[0].(chan<- struct{})
	ret1, _ := ret[1].(<-chan error)
	return ret0, ret1
}

// VerifyHeaders indicates an expected call of VerifyHeaders
func (mr *MockEngineMockRecorder) VerifyHeaders(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyHeaders", reflect.TypeOf((*MockEngine)(nil).VerifyHeaders), arg0, arg1, arg2)
}

// VerifySeal mocks base method
func (m *MockEngine) VerifySeal(arg0 consensus.ChainReader, arg1 *types.Header) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySeal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifySeal indicates an expected call of VerifySeal
func (mr *MockEngineMockRecorder) VerifySeal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySeal", reflect.TypeOf((*MockEngine)(nil).VerifySeal), arg0, arg1)
}
