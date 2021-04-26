// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from node/service.go (2018/06/04).
// Modified and improved for the klaytn development.

package node

import (
	"crypto/ecdsa"
	"reflect"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/networks/p2p"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/storage/database"
)

type ServiceContext struct {
	config         *Config
	services       map[reflect.Type]Service
	EventMux       *event.TypeMux
	AccountManager *accounts.Manager
}

func NewServiceContext(conf *Config, srv map[reflect.Type]Service, mux *event.TypeMux, am *accounts.Manager) *ServiceContext {
	return &ServiceContext{conf, srv, mux, am}
}

// OpenDatabase opens an existing database with the given name (or creates one
// if no previous can be found) from within the node's data directory. If the
// node is an ephemeral one, a memory database is returned.
func (ctx *ServiceContext) OpenDatabase(dbc *database.DBConfig) database.DBManager {
	if ctx.config.DataDir == "" {
		return database.NewMemoryDBManager()
	}
	dbc.Dir = ctx.config.ResolvePath(dbc.Dir)
	return database.NewDBManager(dbc)
}

// ResolvePath resolves a user path into the data directory if that was relative
// and if the user actually uses persistent storage. It will return an empty string
// for emphemeral storage and the user's own input for absolute paths.
func (ctx *ServiceContext) ResolvePath(path string) string {
	return ctx.config.ResolvePath(path)
}

// Service retrieves a currently running service registered of a specific type.
func (ctx *ServiceContext) Service(service interface{}) error {
	element := reflect.ValueOf(service).Elem()
	if running, ok := ctx.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return ErrServiceUnknown
}

// NodeKey returns node key from config
func (ctx *ServiceContext) NodeKey() *ecdsa.PrivateKey {
	return ctx.config.NodeKey()
}

func (ctx *ServiceContext) NodeType() common.ConnType {
	return ctx.config.P2P.ConnectionType
}

// ServiceConstructor is the function signature of the constructors needed to be
// registered for service instantiation.
type ServiceConstructor func(ctx *ServiceContext) (Service, error)

// Service is an individual protocol that can be registered into a node.
//
// Notes:
//
// • Service life-cycle management is delegated to the node. The service is allowed to
// initialize itself upon creation, but no goroutines should be spun up outside of the
// Start method.
//
// • Restart logic is not required as the node will create a fresh instance
// every time a service is started.
type Service interface {
	// Protocols retrieves the P2P protocols the service wishes to start.
	Protocols() []p2p.Protocol

	// APIs retrieves the list of RPC descriptors the service provides
	APIs() []rpc.API

	// Start is called after all services have been constructed and the networking
	// layer was also initialized to spawn any goroutines required by the service.
	Start(server p2p.Server) error

	// Stop terminates all goroutines belonging to the service, blocking until they
	// are all terminated.
	Stop() error

	// retrieve components (blockchain, txpool, ..) from core service
	Components() []interface{}

	// set components (blockchain, txpool, ..) in core service
	SetComponents(components []interface{})
}
