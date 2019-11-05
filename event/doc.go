// Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from event/event.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package event deals with subscriptions to real-time events.

Package event provides three different types of event dispatchers and different go-routines in a node can receive/send data by using it.

Source Files

Each file provides the following features

 - event.go: Provides `TypeMux` struct which dispatches events to registered receivers. Receivers can be registered to handle events of a certain type
 - feed.go: Provides `Feed` struct implements one-to-many subscriptions where the carrier of events is a channel. Values sent to a Feed are delivered to all subscribed channels simultaneously. Feeds can only be used with a single type
 - subscription.go: Provides `Subscription` interface which represents a stream of events. The carrier of the events is typically a channel but isn't part of the interface. Subscriptions can fail while established and failures are reported through an error channel

*/
package event
