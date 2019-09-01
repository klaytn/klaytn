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
Package core implements the core functionality of Istanbul consensus engine.
The package defines `core` struct which stores all current consensus status such as validator set and round status.
It also has accompanying methods to handle and broadcast messages between nodes.

In Istanbul consensus, there are 3 phases including Preprepare, Prepare and Commit.
Each phase has its own message and by handling these messages the consensus can be made.

When a consensus is made for a given block, core communicates with the Istanbul backend to proceed.

Source Files

Core package is composed of following files
 - `backlog.go`: Implements core methods handling future messages. The future message is a message which has future timestamp or in a different phase
 - `commit.go`: Implements core methods which send, receive, handle, verify and accept commit messages
 - `core.go`: Defines core struct and its methods related to timer setup, start new round and round state update
 - `errors.go`: Defines consensus message related errors
 - `events.go`: Defines backlog event and timeout event
 - `final_committed.go`: Start a new round when a final committed proposal is stored
 - `handler.go`: Implements core.Engine.Start and Stop. Provides event and message hendlers
 - `message_set.go`: Defines messageSet struct which has a validator set and messages from other nodes
 - `prepare.go`: Implements core methods which send, receive, handle, verify and accept prepare phase messages
 - `preprepare.go`: Implements core methods which send, handle and accept preprepare messages
 - `request.go`: Implements core methods which handle, check, store and process preprepare messages
 - `roundchange.go`: Implement core methods receiving and handling roundchange messages
 - `roundstate.go`: Defines roundState struct which has messages of each phase for a round
 - `types.go`: Defines Engine interface and message, State type
*/
package core
