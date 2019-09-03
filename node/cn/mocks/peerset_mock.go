// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/klaytn/klaytn/node/cn (interfaces: PeerSet)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	common "github.com/klaytn/klaytn/common"
	p2p "github.com/klaytn/klaytn/networks/p2p"
	cn "github.com/klaytn/klaytn/node/cn"
	reflect "reflect"
)

// MockPeerSet is a mock of PeerSet interface
type MockPeerSet struct {
	ctrl     *gomock.Controller
	recorder *MockPeerSetMockRecorder
}

// MockPeerSetMockRecorder is the mock recorder for MockPeerSet
type MockPeerSetMockRecorder struct {
	mock *MockPeerSet
}

// NewMockPeerSet creates a new mock instance
func NewMockPeerSet(ctrl *gomock.Controller) *MockPeerSet {
	mock := &MockPeerSet{ctrl: ctrl}
	mock.recorder = &MockPeerSetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPeerSet) EXPECT() *MockPeerSetMockRecorder {
	return m.recorder
}

// AnotherTypePeersWithTx mocks base method
func (m *MockPeerSet) AnotherTypePeersWithTx(arg0 common.Hash, arg1 p2p.ConnType) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnotherTypePeersWithTx", arg0, arg1)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// AnotherTypePeersWithTx indicates an expected call of AnotherTypePeersWithTx
func (mr *MockPeerSetMockRecorder) AnotherTypePeersWithTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AnotherTypePeersWithTx", reflect.TypeOf((*MockPeerSet)(nil).AnotherTypePeersWithTx), arg0, arg1)
}

// AnotherTypePeersWithoutTx mocks base method
func (m *MockPeerSet) AnotherTypePeersWithoutTx(arg0 common.Hash, arg1 p2p.ConnType) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnotherTypePeersWithoutTx", arg0, arg1)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// AnotherTypePeersWithoutTx indicates an expected call of AnotherTypePeersWithoutTx
func (mr *MockPeerSetMockRecorder) AnotherTypePeersWithoutTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AnotherTypePeersWithoutTx", reflect.TypeOf((*MockPeerSet)(nil).AnotherTypePeersWithoutTx), arg0, arg1)
}

// BestPeer mocks base method
func (m *MockPeerSet) BestPeer() cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BestPeer")
	ret0, _ := ret[0].(cn.Peer)
	return ret0
}

// BestPeer indicates an expected call of BestPeer
func (mr *MockPeerSetMockRecorder) BestPeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestPeer", reflect.TypeOf((*MockPeerSet)(nil).BestPeer))
}

// CNPeers mocks base method
func (m *MockPeerSet) CNPeers() map[common.Address]cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CNPeers")
	ret0, _ := ret[0].(map[common.Address]cn.Peer)
	return ret0
}

// CNPeers indicates an expected call of CNPeers
func (mr *MockPeerSetMockRecorder) CNPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CNPeers", reflect.TypeOf((*MockPeerSet)(nil).CNPeers))
}

// CNWithoutBlock mocks base method
func (m *MockPeerSet) CNWithoutBlock(arg0 common.Hash) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CNWithoutBlock", arg0)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// CNWithoutBlock indicates an expected call of CNWithoutBlock
func (mr *MockPeerSetMockRecorder) CNWithoutBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CNWithoutBlock", reflect.TypeOf((*MockPeerSet)(nil).CNWithoutBlock), arg0)
}

// CNWithoutTx mocks base method
func (m *MockPeerSet) CNWithoutTx(arg0 common.Hash) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CNWithoutTx", arg0)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// CNWithoutTx indicates an expected call of CNWithoutTx
func (mr *MockPeerSetMockRecorder) CNWithoutTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CNWithoutTx", reflect.TypeOf((*MockPeerSet)(nil).CNWithoutTx), arg0)
}

// Close mocks base method
func (m *MockPeerSet) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close
func (mr *MockPeerSetMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPeerSet)(nil).Close))
}

// ENPeers mocks base method
func (m *MockPeerSet) ENPeers() map[common.Address]cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ENPeers")
	ret0, _ := ret[0].(map[common.Address]cn.Peer)
	return ret0
}

// ENPeers indicates an expected call of ENPeers
func (mr *MockPeerSetMockRecorder) ENPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ENPeers", reflect.TypeOf((*MockPeerSet)(nil).ENPeers))
}

// ENWithoutBlock mocks base method
func (m *MockPeerSet) ENWithoutBlock(arg0 common.Hash) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ENWithoutBlock", arg0)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// ENWithoutBlock indicates an expected call of ENWithoutBlock
func (mr *MockPeerSetMockRecorder) ENWithoutBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ENWithoutBlock", reflect.TypeOf((*MockPeerSet)(nil).ENWithoutBlock), arg0)
}

// Len mocks base method
func (m *MockPeerSet) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len
func (mr *MockPeerSetMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockPeerSet)(nil).Len))
}

// PNPeers mocks base method
func (m *MockPeerSet) PNPeers() map[common.Address]cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PNPeers")
	ret0, _ := ret[0].(map[common.Address]cn.Peer)
	return ret0
}

// PNPeers indicates an expected call of PNPeers
func (mr *MockPeerSetMockRecorder) PNPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PNPeers", reflect.TypeOf((*MockPeerSet)(nil).PNPeers))
}

// PNWithoutBlock mocks base method
func (m *MockPeerSet) PNWithoutBlock(arg0 common.Hash) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PNWithoutBlock", arg0)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// PNWithoutBlock indicates an expected call of PNWithoutBlock
func (mr *MockPeerSetMockRecorder) PNWithoutBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PNWithoutBlock", reflect.TypeOf((*MockPeerSet)(nil).PNWithoutBlock), arg0)
}

// Peer mocks base method
func (m *MockPeerSet) Peer(arg0 string) cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peer", arg0)
	ret0, _ := ret[0].(cn.Peer)
	return ret0
}

// Peer indicates an expected call of Peer
func (mr *MockPeerSetMockRecorder) Peer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peer", reflect.TypeOf((*MockPeerSet)(nil).Peer), arg0)
}

// Peers mocks base method
func (m *MockPeerSet) Peers() map[string]cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peers")
	ret0, _ := ret[0].(map[string]cn.Peer)
	return ret0
}

// Peers indicates an expected call of Peers
func (mr *MockPeerSetMockRecorder) Peers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peers", reflect.TypeOf((*MockPeerSet)(nil).Peers))
}

// PeersWithoutBlock mocks base method
func (m *MockPeerSet) PeersWithoutBlock(arg0 common.Hash) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeersWithoutBlock", arg0)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// PeersWithoutBlock indicates an expected call of PeersWithoutBlock
func (mr *MockPeerSetMockRecorder) PeersWithoutBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeersWithoutBlock", reflect.TypeOf((*MockPeerSet)(nil).PeersWithoutBlock), arg0)
}

// PeersWithoutBlockExceptCN mocks base method
func (m *MockPeerSet) PeersWithoutBlockExceptCN(arg0 common.Hash) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeersWithoutBlockExceptCN", arg0)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// PeersWithoutBlockExceptCN indicates an expected call of PeersWithoutBlockExceptCN
func (mr *MockPeerSetMockRecorder) PeersWithoutBlockExceptCN(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeersWithoutBlockExceptCN", reflect.TypeOf((*MockPeerSet)(nil).PeersWithoutBlockExceptCN), arg0)
}

// PeersWithoutTx mocks base method
func (m *MockPeerSet) PeersWithoutTx(arg0 common.Hash) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeersWithoutTx", arg0)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// PeersWithoutTx indicates an expected call of PeersWithoutTx
func (mr *MockPeerSetMockRecorder) PeersWithoutTx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeersWithoutTx", reflect.TypeOf((*MockPeerSet)(nil).PeersWithoutTx), arg0)
}

// Register mocks base method
func (m *MockPeerSet) Register(arg0 cn.Peer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockPeerSetMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockPeerSet)(nil).Register), arg0)
}

// RegisterValidator mocks base method
func (m *MockPeerSet) RegisterValidator(arg0 p2p.ConnType, arg1 p2p.PeerTypeValidator) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterValidator", arg0, arg1)
}

// RegisterValidator indicates an expected call of RegisterValidator
func (mr *MockPeerSetMockRecorder) RegisterValidator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterValidator", reflect.TypeOf((*MockPeerSet)(nil).RegisterValidator), arg0, arg1)
}

// TypePeers mocks base method
func (m *MockPeerSet) TypePeers(arg0 p2p.ConnType) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TypePeers", arg0)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// TypePeers indicates an expected call of TypePeers
func (mr *MockPeerSetMockRecorder) TypePeers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TypePeers", reflect.TypeOf((*MockPeerSet)(nil).TypePeers), arg0)
}

// TypePeersWithTx mocks base method
func (m *MockPeerSet) TypePeersWithTx(arg0 common.Hash, arg1 p2p.ConnType) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TypePeersWithTx", arg0, arg1)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// TypePeersWithTx indicates an expected call of TypePeersWithTx
func (mr *MockPeerSetMockRecorder) TypePeersWithTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TypePeersWithTx", reflect.TypeOf((*MockPeerSet)(nil).TypePeersWithTx), arg0, arg1)
}

// TypePeersWithoutBlock mocks base method
func (m *MockPeerSet) TypePeersWithoutBlock(arg0 common.Hash, arg1 p2p.ConnType) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TypePeersWithoutBlock", arg0, arg1)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// TypePeersWithoutBlock indicates an expected call of TypePeersWithoutBlock
func (mr *MockPeerSetMockRecorder) TypePeersWithoutBlock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TypePeersWithoutBlock", reflect.TypeOf((*MockPeerSet)(nil).TypePeersWithoutBlock), arg0, arg1)
}

// TypePeersWithoutTx mocks base method
func (m *MockPeerSet) TypePeersWithoutTx(arg0 common.Hash, arg1 p2p.ConnType) []cn.Peer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TypePeersWithoutTx", arg0, arg1)
	ret0, _ := ret[0].([]cn.Peer)
	return ret0
}

// TypePeersWithoutTx indicates an expected call of TypePeersWithoutTx
func (mr *MockPeerSetMockRecorder) TypePeersWithoutTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TypePeersWithoutTx", reflect.TypeOf((*MockPeerSet)(nil).TypePeersWithoutTx), arg0, arg1)
}

// Unregister mocks base method
func (m *MockPeerSet) Unregister(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unregister", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unregister indicates an expected call of Unregister
func (mr *MockPeerSetMockRecorder) Unregister(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unregister", reflect.TypeOf((*MockPeerSet)(nil).Unregister), arg0)
}
