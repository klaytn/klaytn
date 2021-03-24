// Copyright 2019 The klaytn Authors
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

package main

import (
	"github.com/klaytn/klaytn/networks/p2p/discover"
)

type PublicBootnodeAPI struct {
	bn *BN
}

func NewPublicBootnodeAPI(b *BN) *PublicBootnodeAPI {
	return &PublicBootnodeAPI{bn: b}
}

func (api *PublicBootnodeAPI) GetAuthorizedNodes() []*discover.Node {
	return api.bn.GetAuthorizedNodes()
}

type PrivateBootnodeAPI struct {
	bn *BN
}

func NewPrivateBootnodeAPI(b *BN) *PrivateBootnodeAPI {
	return &PrivateBootnodeAPI{bn: b}
}

func (api *PrivateBootnodeAPI) Name() string {
	return api.bn.Name()
}

func (api *PrivateBootnodeAPI) Resolve(target discover.NodeID, targetType discover.NodeType) *discover.Node {
	return api.bn.Resolve(target, targetType)
}

func (api *PrivateBootnodeAPI) Lookup(target discover.NodeID, targetType discover.NodeType) []*discover.Node {
	return api.bn.Lookup(target, targetType)
}

func (api *PrivateBootnodeAPI) ReadRandomNodes(nType discover.NodeType) []*discover.Node {
	var buf []*discover.Node
	api.bn.ReadRandomNodes(buf, nType)
	return buf
}

func (api *PrivateBootnodeAPI) CreateUpdateNodeOnDB(nodekni string) error {
	return api.bn.CreateUpdateNodeOnDB(nodekni)
}

func (api *PrivateBootnodeAPI) CreateUpdateNodeOnTable(nodekni string) error {
	return api.bn.CreateUpdateNodeOnTable(nodekni)
}

func (api *PrivateBootnodeAPI) GetNodeFromDB(id discover.NodeID) (*discover.Node, error) {
	return api.bn.GetNodeFromDB(id)
}

func (api *PrivateBootnodeAPI) GetTableEntries() []*discover.Node {
	return api.bn.GetTableEntries()
}

func (api *PrivateBootnodeAPI) GetTableReplacements() []*discover.Node {
	return api.bn.GetTableReplacements()
}

func (api *PrivateBootnodeAPI) DeleteNodeFromDB(nodekni string) error {
	return api.bn.DeleteNodeFromDB(nodekni)
}

func (api *PrivateBootnodeAPI) DeleteNodeFromTable(nodekni string) error {
	return api.bn.DeleteNodeFromTable(nodekni)
}

func (api *PrivateBootnodeAPI) PutAuthorizedNodes(rawurl string) error {
	return api.bn.PutAuthorizedNodes(rawurl)
}

func (api *PrivateBootnodeAPI) DeleteAuthorizedNodes(rawurl string) error {
	return api.bn.DeleteAuthorizedNodes(rawurl)
}
