// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package adapters

import (
	"errors"
	"fmt"
	"math"
	"net"
	"sync"

	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/p2p/simulations/pipes"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/node"
)

// CnAdapter is a NodeAdapter which creates in-memory simulation nodes and
// connects them using net.Pipe
type CnAdapter struct {
	pipe     func() (net.Conn, net.Conn, error)
	mtx      sync.RWMutex
	nodes    map[discover.NodeID]*CnNode
	services map[string]ServiceFunc
}

// NewCnAdapter creates a CnAdapter which is capable of running in-memory
// simulation nodes running any of the given services (the services to run on a
// particular node are passed to the NewNode function in the NodeConfig)
// the adapter uses a net.Pipe for in-memory simulated network connections
func NewCnAdapter(services map[string]ServiceFunc) *CnAdapter {
	return &CnAdapter{
		pipe:     pipes.NetPipe,
		nodes:    make(map[discover.NodeID]*CnNode),
		services: services,
	}
}

// Name returns the name of the adapter for logging purposes
func (s *CnAdapter) Name() string {
	return "cnsim-adapter"
}

// NewNode returns a new CnNode using the given config
func (s *CnAdapter) NewNode(config *NodeConfig) (Node, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// check a node with the ID doesn't already exist
	id := config.ID
	if _, exists := s.nodes[id]; exists {
		return nil, fmt.Errorf("node already exists: %s", id)
	}

	// check the services are valid
	if len(config.Services) == 0 {
		return nil, errors.New("node must have at least one service")
	}
	for _, service := range config.Services {
		if _, exists := s.services[service]; !exists {
			return nil, fmt.Errorf("unknown node service %q", service)
		}
	}

	n, err := node.New(&node.Config{
		P2P: p2p.Config{
			PrivateKey:             config.PrivateKey, // from p2psim client
			MaxPhysicalConnections: math.MaxInt32,
			NoDiscovery:            true,
			ListenAddr:             fmt.Sprintf(":%d", config.Port),
			//Dialer:          s,
			EnableMsgEvents: config.EnableMsgEvents,
		},
		//Logger: log.New("node.id", id.String()),
		Logger: logger.NewWith("node.name", config.Name),

		//Logger: log.New(),
	})
	if err != nil {
		return nil, err
	}

	CnNode := &CnNode{
		ID:      id,
		config:  config,
		node:    n,
		adapter: s,
		running: make(map[string]node.Service),
	}
	s.nodes[id] = CnNode
	return CnNode, nil
}

// TODO : NOT USED
// Dial implements the p2p.NodeDialer interface by connecting to the node using
// an in-memory net.Pipe
func (s *CnAdapter) Dial(dest *discover.Node) (conn net.Conn, err error) {
	node, ok := s.GetNode(dest.ID)
	if !ok {
		return nil, fmt.Errorf("unknown node: %s", dest.ID)
	}
	srv := node.Server()
	if srv == nil {
		return nil, fmt.Errorf("node not running: %s", dest.ID)
	}
	// CnAdapter.pipe is net.Pipe (NewCnAdapter)
	pipe1, pipe2, err := s.pipe()
	if err != nil {
		return nil, err
	}
	// this is simulated 'listening'
	// asynchronously call the dialed destintion node's p2p server
	// to set up connection on the 'listening' side
	go srv.SetupConn(pipe1, 0, nil)
	return pipe2, nil
}

// DialRPC implements the RPCDialer interface by creating an in-memory RPC
// client of the given node
func (s *CnAdapter) DialRPC(id discover.NodeID) (*rpc.Client, error) {
	node, ok := s.GetNode(id)
	if !ok {
		return nil, fmt.Errorf("unknown node: %s", id)
	}
	handler, err := node.node.RPCHandler()
	if err != nil {
		return nil, err
	}
	return rpc.DialInProc(handler), nil
}

// GetNode returns the node with the given ID if it exists
func (s *CnAdapter) GetNode(id discover.NodeID) (*CnNode, bool) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	node, ok := s.nodes[id]
	return node, ok
}

// CnNode is an in-memory simulation node which connects to other nodes using
// net.Pipe (see CnAdapter.Dial), running devp2p protocols directly over that
// pipe
type CnNode struct {
	lock         sync.RWMutex
	ID           discover.NodeID
	config       *NodeConfig
	adapter      *CnAdapter
	node         *node.Node
	running      map[string]node.Service
	client       *rpc.Client
	registerOnce sync.Once
}

// Addr returns the node's discovery address
func (sn *CnNode) Addr() []byte {
	return []byte(sn.Node().String())
}

// Node returns a discover.Node representing the CnNode
func (sn *CnNode) Node() *discover.Node {
	//return discover.NewNode(sn.ID, net.IP{127, 0, 0, 1}, 30303, 30303)
	return discover.NewNode(sn.ID, net.IP{127, 0, 0, 1}, sn.config.Port, sn.config.Port, nil, discover.NodeTypeCN)
}

// Client returns an rpc.Client which can be used to communicate with the
// underlying services (it is set once the node has started)
func (sn *CnNode) Client() (*rpc.Client, error) {
	sn.lock.RLock()
	defer sn.lock.RUnlock()
	if sn.client == nil {
		return nil, errors.New("node not started")
	}
	return sn.client, nil
}

// ServeRPC serves RPC requests over the given connection by creating an
// in-memory client to the node's RPC server
func (sn *CnNode) ServeRPC(conn net.Conn) error {
	handler, err := sn.node.RPCHandler()
	if err != nil {
		return err
	}
	handler.ServeCodec(rpc.NewJSONCodec(conn), rpc.OptionMethodInvocation|rpc.OptionSubscriptions)
	return nil
}

// Snapshots creates snapshots of the services by calling the
// simulation_snapshot RPC method
func (sn *CnNode) Snapshots() (map[string][]byte, error) {
	sn.lock.RLock()
	services := make(map[string]node.Service, len(sn.running))
	for name, service := range sn.running {
		services[name] = service
	}
	sn.lock.RUnlock()
	if len(services) == 0 {
		return nil, errors.New("no running services")
	}
	snapshots := make(map[string][]byte)
	for name, service := range services {
		if s, ok := service.(interface {
			Snapshot() ([]byte, error)
		}); ok {
			snap, err := s.Snapshot()
			if err != nil {
				return nil, err
			}
			snapshots[name] = snap
		}
	}
	return snapshots, nil
}

// Start registers the services and starts the underlying devp2p node
func (sn *CnNode) Start(snapshots map[string][]byte) error {
	newService := func(name string) func(ctx *node.ServiceContext) (node.Service, error) {
		return func(nodeCtx *node.ServiceContext) (node.Service, error) {
			ctx := &ServiceContext{
				RPCDialer:   sn.adapter,
				NodeContext: nodeCtx,
				Config:      sn.config,
			}
			if snapshots != nil {
				ctx.Snapshot = snapshots[name]
			}
			serviceFunc := sn.adapter.services[name]
			service, err := serviceFunc(ctx)
			if err != nil {
				return nil, err
			}
			sn.running[name] = service
			return service, nil
		}
	}

	// ensure we only register the services once in the case of the node
	// being stopped and then started again
	var regErr error
	sn.registerOnce.Do(func() {
		for _, name := range sn.config.Services {
			if err := sn.node.Register(newService(name)); err != nil {
				regErr = err
				break
			}
		}
	})
	if regErr != nil {
		return regErr
	}

	if err := sn.node.Start(); err != nil {
		return err
	}

	// create an in-process RPC client
	handler, err := sn.node.RPCHandler()
	if err != nil {
		return err
	}

	sn.lock.Lock()
	sn.client = rpc.DialInProc(handler)
	sn.lock.Unlock()

	return nil
}

// Stop closes the RPC client and stops the underlying devp2p node
func (sn *CnNode) Stop() error {
	sn.lock.Lock()
	if sn.client != nil {
		sn.client.Close()
		sn.client = nil
	}
	sn.lock.Unlock()
	return sn.node.Stop()
}

// Services returns a copy of the underlying services
func (sn *CnNode) Services() []node.Service {
	sn.lock.RLock()
	defer sn.lock.RUnlock()
	services := make([]node.Service, 0, len(sn.running))
	for _, service := range sn.running {
		services = append(services, service)
	}
	return services
}

// Server returns the underlying p2p.Server
func (sn *CnNode) Server() p2p.Server {
	return sn.node.Server()
}

// SubscribeEvents subscribes the given channel to peer events from the
// underlying p2p.Server
func (sn *CnNode) SubscribeEvents(ch chan *p2p.PeerEvent) event.Subscription {
	srv := sn.Server()
	if srv == nil {
		panic("node not running")
	}
	return srv.SubscribeEvents(ch)
}

// NodeInfo returns information about the node
func (sn *CnNode) NodeInfo() *p2p.NodeInfo {
	server := sn.Server()
	if server == nil {
		return &p2p.NodeInfo{
			ID:    sn.ID.String(),
			Enode: sn.Node().String(),
		}
	}
	return server.NodeInfo()
}

func (sn *CnNode) PeersInfo() []*p2p.PeerInfo {
	server := sn.Server()
	if server == nil {
		return nil
	}
	return server.PeersInfo()
}

func (sn *CnNode) GetPeerCount() int {
	srv := sn.Server()
	if srv == nil {
		panic("node not running")
	}
	return srv.PeerCount()
}

func (sn *CnNode) DisconnectPeer(destID discover.NodeID) {
	srv := sn.Server()
	if srv == nil {
		panic("node not running")
	}
	srv.Disconnect(destID)
}

/*
func setSocketBuffer(conn net.Conn, socketReadBuffer int, socketWriteBuffer int) error {
	switch v := conn.(type) {
	case *net.UnixConn:
		err := v.SetReadBuffer(socketReadBuffer)
		if err != nil {
			return err
		}
		err = v.SetWriteBuffer(socketWriteBuffer)
		if err != nil {
			return err
		}
	}
	return nil
}*/
