// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/klaytn/klaytn/node/sc (interfaces: BridgePeer)

// Package sc is a generated GoMock package.
package sc

import (
	gomock "github.com/golang/mock/gomock"
	types "github.com/klaytn/klaytn/blockchain/types"
	common "github.com/klaytn/klaytn/common"
	p2p "github.com/klaytn/klaytn/networks/p2p"
	discover "github.com/klaytn/klaytn/networks/p2p/discover"
	big "math/big"
	reflect "reflect"
)

// MockBridgePeer is a mock of BridgePeer interface
type MockBridgePeer struct {
	ctrl     *gomock.Controller
	recorder *MockBridgePeerMockRecorder
}

// MockBridgePeerMockRecorder is the mock recorder for MockBridgePeer
type MockBridgePeerMockRecorder struct {
	mock *MockBridgePeer
}

// NewMockBridgePeer creates a new mock instance
func NewMockBridgePeer(ctrl *gomock.Controller) *MockBridgePeer {
	mock := &MockBridgePeer{ctrl: ctrl}
	mock.recorder = &MockBridgePeerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBridgePeer) EXPECT() *MockBridgePeerMockRecorder {
	return m.recorder
}

// AddToKnownTxs mocks base method
func (m *MockBridgePeer) AddToKnownTxs(arg0 common.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddToKnownTxs", arg0)
}

// AddToKnownTxs indicates an expected call of AddToKnownTxs
func (mr *MockBridgePeerMockRecorder) AddToKnownTxs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToKnownTxs", reflect.TypeOf((*MockBridgePeer)(nil).AddToKnownTxs), arg0)
}

// Close mocks base method
func (m *MockBridgePeer) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockBridgePeerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockBridgePeer)(nil).Close))
}

// ConnType mocks base method
func (m *MockBridgePeer) ConnType() common.ConnType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnType")
	ret0, _ := ret[0].(common.ConnType)
	return ret0
}

// ConnType indicates an expected call of ConnType
func (mr *MockBridgePeerMockRecorder) ConnType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnType", reflect.TypeOf((*MockBridgePeer)(nil).ConnType))
}

// GetAddr mocks base method
func (m *MockBridgePeer) GetAddr() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddr")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetAddr indicates an expected call of GetAddr
func (mr *MockBridgePeerMockRecorder) GetAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddr", reflect.TypeOf((*MockBridgePeer)(nil).GetAddr))
}

// GetChainID mocks base method
func (m *MockBridgePeer) GetChainID() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainID")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// GetChainID indicates an expected call of GetChainID
func (mr *MockBridgePeerMockRecorder) GetChainID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainID", reflect.TypeOf((*MockBridgePeer)(nil).GetChainID))
}

// GetID mocks base method
func (m *MockBridgePeer) GetID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetID indicates an expected call of GetID
func (mr *MockBridgePeerMockRecorder) GetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetID", reflect.TypeOf((*MockBridgePeer)(nil).GetID))
}

// GetP2PPeer mocks base method
func (m *MockBridgePeer) GetP2PPeer() *p2p.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetP2PPeer")
	ret0, _ := ret[0].(*p2p.Peer)
	return ret0
}

// GetP2PPeer indicates an expected call of GetP2PPeer
func (mr *MockBridgePeerMockRecorder) GetP2PPeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetP2PPeer", reflect.TypeOf((*MockBridgePeer)(nil).GetP2PPeer))
}

// GetP2PPeerID mocks base method
func (m *MockBridgePeer) GetP2PPeerID() discover.NodeID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetP2PPeerID")
	ret0, _ := ret[0].(discover.NodeID)
	return ret0
}

// GetP2PPeerID indicates an expected call of GetP2PPeerID
func (mr *MockBridgePeerMockRecorder) GetP2PPeerID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetP2PPeerID", reflect.TypeOf((*MockBridgePeer)(nil).GetP2PPeerID))
}

// GetRW mocks base method
func (m *MockBridgePeer) GetRW() p2p.MsgReadWriter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRW")
	ret0, _ := ret[0].(p2p.MsgReadWriter)
	return ret0
}

// GetRW indicates an expected call of GetRW
func (mr *MockBridgePeerMockRecorder) GetRW() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRW", reflect.TypeOf((*MockBridgePeer)(nil).GetRW))
}

// GetVersion mocks base method
func (m *MockBridgePeer) GetVersion() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetVersion indicates an expected call of GetVersion
func (mr *MockBridgePeerMockRecorder) GetVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockBridgePeer)(nil).GetVersion))
}

// Handle mocks base method
func (m *MockBridgePeer) Handle(arg0 *MainBridge) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle
func (mr *MockBridgePeerMockRecorder) Handle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockBridgePeer)(nil).Handle), arg0)
}

// Handshake mocks base method
func (m *MockBridgePeer) Handshake(arg0 uint64, arg1, arg2 *big.Int, arg3 common.Hash) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handshake", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handshake indicates an expected call of Handshake
func (mr *MockBridgePeerMockRecorder) Handshake(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handshake", reflect.TypeOf((*MockBridgePeer)(nil).Handshake), arg0, arg1, arg2, arg3)
}

// Head mocks base method
func (m *MockBridgePeer) Head() (common.Hash, *big.Int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Head")
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(*big.Int)
	return ret0, ret1
}

// Head indicates an expected call of Head
func (mr *MockBridgePeerMockRecorder) Head() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Head", reflect.TypeOf((*MockBridgePeer)(nil).Head))
}

// Info mocks base method
func (m *MockBridgePeer) Info() *BridgePeerInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info")
	ret0, _ := ret[0].(*BridgePeerInfo)
	return ret0
}

// Info indicates an expected call of Info
func (mr *MockBridgePeerMockRecorder) Info() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockBridgePeer)(nil).Info))
}

// KnowsTx mocks base method
func (m *MockBridgePeer) KnowsTx(arg0 common.Hash) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KnowsTx", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// KnowsTx indicates an expected call of KnowsTx
func (mr *MockBridgePeerMockRecorder) KnowsTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KnowsTx", reflect.TypeOf((*MockBridgePeer)(nil).KnowsTx), arg0)
}

// Send mocks base method
func (m *MockBridgePeer) Send(arg0 uint64, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockBridgePeerMockRecorder) Send(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockBridgePeer)(nil).Send), arg0, arg1)
}

// SendRequestRPC mocks base method
func (m *MockBridgePeer) SendRequestRPC(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRequestRPC", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendRequestRPC indicates an expected call of SendRequestRPC
func (mr *MockBridgePeerMockRecorder) SendRequestRPC(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRequestRPC", reflect.TypeOf((*MockBridgePeer)(nil).SendRequestRPC), arg0)
}

// SendResponseRPC mocks base method
func (m *MockBridgePeer) SendResponseRPC(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendResponseRPC", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendResponseRPC indicates an expected call of SendResponseRPC
func (mr *MockBridgePeerMockRecorder) SendResponseRPC(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendResponseRPC", reflect.TypeOf((*MockBridgePeer)(nil).SendResponseRPC), arg0)
}

// SendServiceChainInfoRequest mocks base method
func (m *MockBridgePeer) SendServiceChainInfoRequest(arg0 *common.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendServiceChainInfoRequest", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendServiceChainInfoRequest indicates an expected call of SendServiceChainInfoRequest
func (mr *MockBridgePeerMockRecorder) SendServiceChainInfoRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendServiceChainInfoRequest", reflect.TypeOf((*MockBridgePeer)(nil).SendServiceChainInfoRequest), arg0)
}

// SendServiceChainInfoResponse mocks base method
func (m *MockBridgePeer) SendServiceChainInfoResponse(arg0 *parentChainInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendServiceChainInfoResponse", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendServiceChainInfoResponse indicates an expected call of SendServiceChainInfoResponse
func (mr *MockBridgePeerMockRecorder) SendServiceChainInfoResponse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendServiceChainInfoResponse", reflect.TypeOf((*MockBridgePeer)(nil).SendServiceChainInfoResponse), arg0)
}

// SendServiceChainReceiptRequest mocks base method
func (m *MockBridgePeer) SendServiceChainReceiptRequest(arg0 []common.Hash) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendServiceChainReceiptRequest", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendServiceChainReceiptRequest indicates an expected call of SendServiceChainReceiptRequest
func (mr *MockBridgePeerMockRecorder) SendServiceChainReceiptRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendServiceChainReceiptRequest", reflect.TypeOf((*MockBridgePeer)(nil).SendServiceChainReceiptRequest), arg0)
}

// SendServiceChainReceiptResponse mocks base method
func (m *MockBridgePeer) SendServiceChainReceiptResponse(arg0 []*types.ReceiptForStorage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendServiceChainReceiptResponse", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendServiceChainReceiptResponse indicates an expected call of SendServiceChainReceiptResponse
func (mr *MockBridgePeerMockRecorder) SendServiceChainReceiptResponse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendServiceChainReceiptResponse", reflect.TypeOf((*MockBridgePeer)(nil).SendServiceChainReceiptResponse), arg0)
}

// SendServiceChainTxs mocks base method
func (m *MockBridgePeer) SendServiceChainTxs(arg0 types.Transactions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendServiceChainTxs", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendServiceChainTxs indicates an expected call of SendServiceChainTxs
func (mr *MockBridgePeerMockRecorder) SendServiceChainTxs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendServiceChainTxs", reflect.TypeOf((*MockBridgePeer)(nil).SendServiceChainTxs), arg0)
}

// SetAddr mocks base method
func (m *MockBridgePeer) SetAddr(arg0 common.Address) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetAddr", arg0)
}

// SetAddr indicates an expected call of SetAddr
func (mr *MockBridgePeerMockRecorder) SetAddr(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAddr", reflect.TypeOf((*MockBridgePeer)(nil).SetAddr), arg0)
}
