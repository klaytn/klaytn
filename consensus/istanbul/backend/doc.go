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

/*
Package backend defines backend struct which implements Backend interface of Istanbul consensus engine.
backend struct works as a backbone of the consensus engine having Istanbul core,
required objects to get information for making a consensus, recent messages and a broadcaster to send its message to peer nodes.

Source Files

Implementation of Backend interface and APIs are included in this package
 - `api.go`: Implements APIs which provide the states of Istanbul
 - `backend.go`: Defines backend struct which implements Backend interface working as a backbone of the consensus engine
 - `engine.go`: Implements various backend methods especially for verifying and building header information
 - `handler.go`: Implements backend methods for handling messages and broadcaster
 - `snapshot.go`: Defines snapshot struct which handles votes from nodes and makes governance changes

*/
package backend
