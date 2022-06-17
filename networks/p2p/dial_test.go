// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from p2p/dial_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package p2p

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/networks/p2p/discover"
	"github.com/klaytn/klaytn/networks/p2p/netutil"
)

func init() {
	spew.Config.Indent = "\t"
}

type dialtest struct {
	init   *dialstate // state before and after the test.
	rounds []round
}

type failedInfo struct {
	id        discover.NodeID
	failedTry int
}

type round struct {
	peers   []*Peer      // current peer set
	done    []task       // tasks that got done this round
	new     []task       // the result must match this one
	expired []failedInfo // task
}

func runDialTest(t *testing.T, test dialtest) {
	var (
		vTime   time.Time
		running int
	)
	pm := func(ps []*Peer) map[discover.NodeID]*Peer {
		m := make(map[discover.NodeID]*Peer)
		for _, p := range ps {
			m[p.ID()] = p
		}
		return m
	}
	for i, round := range test.rounds {
		for _, task := range round.done {
			running--
			if running < 0 {
				panic("running task counter underflow")
			}
			test.init.taskDone(task, vTime)
		}

		for _, info := range round.expired {
			if dt, ok := test.init.static[info.id]; ok {
				dt.failedTry = info.failedTry
			}
		}

		new := test.init.newTasks(running, pm(round.peers), vTime)
		if !sametasks(new, round.new) {
			t.Errorf("round %d: new tasks mismatch:\ngot %v\nwant %v\nstate: %v\nrunning: %v\n",
				i, spew.Sdump(new), spew.Sdump(round.new), spew.Sdump(test.init), spew.Sdump(running))
		}

		// Time advances by 16 seconds on every round.
		vTime = vTime.Add(16 * time.Second)
		running += len(new)
	}
}

type fakeTable []*discover.Node

func (t fakeTable) Name() string                                                    { return "fakeTable" }
func (t fakeTable) Self() *discover.Node                                            { return new(discover.Node) }
func (t fakeTable) Close()                                                          {}
func (t fakeTable) Lookup(discover.NodeID, discover.NodeType) []*discover.Node      { return nil }
func (t fakeTable) Resolve(discover.NodeID, discover.NodeType) *discover.Node       { return nil }
func (t fakeTable) GetNodes(targetType discover.NodeType, max int) []*discover.Node { return nil }
func (t fakeTable) ReadRandomNodes(buf []*discover.Node, nType discover.NodeType) int {
	return copy(buf, t)
}
func (t fakeTable) RetrieveNodes(target common.Hash, nType discover.NodeType, nresults int) []*discover.Node {
	return nil
}
func (t fakeTable) CreateUpdateNodeOnDB(n *discover.Node) error              { return nil }
func (t fakeTable) CreateUpdateNodeOnTable(n *discover.Node) error           { return nil }
func (t fakeTable) GetNodeFromDB(id discover.NodeID) (*discover.Node, error) { return nil, nil }
func (t fakeTable) DeleteNodeFromDB(n *discover.Node) error                  { return nil }
func (t fakeTable) DeleteNodeFromTable(n *discover.Node) error               { return nil }
func (t fakeTable) GetBucketEntries() []*discover.Node                       { return nil }
func (t fakeTable) GetReplacements() []*discover.Node                        { return nil }
func (t fakeTable) HasBond(id discover.NodeID) bool                          { return true }
func (t fakeTable) Bond(pinged bool, id discover.NodeID, addr *net.UDPAddr, tcpPort uint16, nType discover.NodeType) (*discover.Node, error) {
	return nil, nil
}
func (t fakeTable) IsAuthorized(id discover.NodeID, ntype discover.NodeType) bool {
	return true
}

func (t fakeTable) GetAuthorizedNodes() []*discover.Node         { return nil }
func (t fakeTable) PutAuthorizedNodes(nodes []*discover.Node)    {}
func (t fakeTable) DeleteAuthorizedNodes(nodes []*discover.Node) {}

// This test checks that dynamic dials are launched from discovery results.
func TestDialStateDynDial(t *testing.T) {
	runDialTest(t, dialtest{
		init: newDialState(nil, nil, fakeTable{}, 5, nil, nil, nil),
		rounds: []round{
			// A discovery query is launched.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(0)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
				},
				new: []task{&discoverTask{}},
			},
			// Dynamic dials are launched when it completes.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(0)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
				},
				done: []task{
					&discoverTask{results: []*discover.Node{
						{ID: uintID(2)}, // this one is already connected and not dialed.
						{ID: uintID(3)},
						{ID: uintID(4)},
						{ID: uintID(5)},
						{ID: uintID(6)}, // these are not tried because max dyn dials is 5
						{ID: uintID(7)}, // ...
					}},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(3)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(4)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(5)}},
				},
			},
			// Some of the dials complete but no new ones are launched yet because
			// the sum of active dial count and dynamic peer count is == maxDynDials.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(0)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(3)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(4)}}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(3)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(4)}},
				},
			},
			// No new dial tasks are launched in the this round because
			// maxDynDials has been reached.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(0)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(3)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(4)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(5)}}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(5)}},
				},
				new: []task{
					&waitExpireTask{Duration: 14 * time.Second},
				},
			},
			// In this round, the peer with id 2 drops off. The query
			// results from last discovery lookup are reused.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(0)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(3)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(4)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(5)}}},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(6)}},
				},
			},
			// More peers (3,4) drop off and dial for ID 6 completes.
			// The last query result from the discovery lookup is reused
			// and a new one is spawned because more candidates are needed.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(0)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(5)}}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(6)}},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(7)}},
					&discoverTask{},
				},
			},
			// Peer 7 is connected, but there still aren't enough dynamic peers
			// (4 out of 5). However, a discovery is already running, so ensure
			// no new is started.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(0)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(5)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(7)}}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(7)}},
				},
			},
			// Finish the running node discovery with an empty set. A new lookup
			// should be immediately requested.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(0)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(5)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(7)}}},
				},
				done: []task{
					&discoverTask{},
				},
				new: []task{
					&discoverTask{},
				},
			},
		},
	})
}

func TestDialStateDynDialFromTable(t *testing.T) {
	// This table always returns the same random nodes
	// in the order given below.
	table := fakeTable{
		{ID: uintID(1)},
		{ID: uintID(2)},
		{ID: uintID(3)},
		{ID: uintID(4)},
		{ID: uintID(5)},
		{ID: uintID(6)},
		{ID: uintID(7)},
		{ID: uintID(8)},
	}

	runDialTest(t, dialtest{
		init: newDialState(nil, nil, table, 10, nil, nil, nil),
		rounds: []round{
			// 5 out of 8 of the nodes returned by ReadRandomNodes are dialed.
			{
				new: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(1)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(2)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(3)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(4)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(5)}},
					&discoverTask{},
				},
			},
			// Dialing nodes 1,2 succeeds. Dials from the lookup are launched.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(1)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(2)}},
					&discoverTask{results: []*discover.Node{
						{ID: uintID(10)},
						{ID: uintID(11)},
						{ID: uintID(12)},
					}},
				},
				new: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(10)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(11)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(12)}},
					&discoverTask{},
				},
			},
			// Dialing nodes 3,4,5 fails. The dials from the lookup succeed.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(10)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(11)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(12)}}},
				},
				done: []task{
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(3)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(4)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(5)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(10)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(11)}},
					&dialTask{flags: dynDialedConn, dest: &discover.Node{ID: uintID(12)}},
				},
			},
			// Waiting for expiry. No waitExpireTask is launched because the
			// discovery query is still running.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(10)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(11)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(12)}}},
				},
			},
			// Nodes 3,4 are not tried again because only the first two
			// returned random nodes (nodes 1,2) are tried and they're
			// already connected.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(10)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(11)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(12)}}},
				},
			},
		},
	})
}

// This test checks that candidates that do not match the netrestrict list are not dialed.
func TestDialStateNetRestrict(t *testing.T) {
	// This table always returns the same random nodes
	// in the order given below.
	table := fakeTable{
		{ID: uintID(1), IP: net.ParseIP("127.0.0.1")},
		{ID: uintID(2), IP: net.ParseIP("127.0.0.2")},
		{ID: uintID(3), IP: net.ParseIP("127.0.0.3")},
		{ID: uintID(4), IP: net.ParseIP("127.0.0.4")},
		{ID: uintID(5), IP: net.ParseIP("127.0.2.5")},
		{ID: uintID(6), IP: net.ParseIP("127.0.2.6")},
		{ID: uintID(7), IP: net.ParseIP("127.0.2.7")},
		{ID: uintID(8), IP: net.ParseIP("127.0.2.8")},
	}
	restrict := new(netutil.Netlist)
	restrict.Add("127.0.2.0/24")

	runDialTest(t, dialtest{
		init: newDialState(nil, nil, table, 10, restrict, nil, nil),
		rounds: []round{
			{
				new: []task{
					&dialTask{flags: dynDialedConn, dest: table[4]},
					&discoverTask{},
				},
			},
		},
	})
}

// This test checks that static dials are launched.
func TestDialStateStaticDial(t *testing.T) {
	wantStatic := []*discover.Node{
		{ID: uintID(1)},
		{ID: uintID(2)},
		{ID: uintID(3)},
		{ID: uintID(4)},
		{ID: uintID(5)},
	}

	runDialTest(t, dialtest{
		init: newDialState(wantStatic, nil, fakeTable{}, 0, nil, nil, nil),
		rounds: []round{
			// Static dials are launched for the nodes that
			// aren't yet connected.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(3)}, dialType: DT_UNLIMITED},
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(4)}, dialType: DT_UNLIMITED},
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(5)}, dialType: DT_UNLIMITED},
				},
			},
			// No new tasks are launched in this round because all static
			// nodes are either connected or still being dialed.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(3)}}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(3)}, dialType: DT_UNLIMITED},
				},
			},
			// No new dial tasks are launched because all static
			// nodes are now connected.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(3)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(4)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(5)}}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(4)}, dialType: DT_UNLIMITED},
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(5)}, dialType: DT_UNLIMITED},
				},
				new: []task{
					&waitExpireTask{Duration: 14 * time.Second},
				},
			},
			// Wait a round for dial history to expire, no new tasks should spawn.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(3)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(4)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(5)}}},
				},
			},
			// If a static node is dropped, it should be immediately redialed,
			// irrespective whether it was originally static or dynamic.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(3)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(5)}}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(2)}, dialType: DT_UNLIMITED},
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(4)}, dialType: DT_UNLIMITED},
				},
			},
		},
	})
}

// This test check expired task removed from static
func TestDialStateTypeStaticDialExpired(t *testing.T) {
	tdt := dialType("test")

	tsMap := make(map[dialType]typedStatic)
	tsMap[tdt] = typedStatic{maxNodeCount: math.MaxInt64, maxTry: 3}

	ds := newDialState(nil, nil, fakeTable{}, 0, nil, nil, tsMap)
	dt := dialtest{
		init: ds,
		rounds: []round{
			{
				new: []task{
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(4)}, dialType: tdt, failedTry: 2},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(5)}, dialType: tdt, failedTry: 3},
					&discoverTypedStaticTask{name: tdt, max: math.MaxInt64 - 2},
				},
			},
		},
	}

	ds.addTypedStatic(&discover.Node{ID: uintID(4)}, tdt)
	ds.static[discover.Node{ID: uintID(4)}.ID].failedTry = 2
	ds.addTypedStatic(&discover.Node{ID: uintID(5)}, tdt)
	ds.static[discover.Node{ID: uintID(5)}.ID].failedTry = 3
	ds.addTypedStatic(&discover.Node{ID: uintID(6)}, tdt) // Expired DialTask, will removed from static
	ds.static[discover.Node{ID: uintID(6)}.ID].failedTry = 4

	runDialTest(t, dt)
}

func TestDialStateTypeStaticDialExpired2(t *testing.T) {
	tdt := dialType("test")

	tsMap := make(map[dialType]typedStatic)
	tsMap[tdt] = typedStatic{maxNodeCount: 6, maxTry: 3}

	ds := newDialState(nil, nil, fakeTable{}, 0, nil, nil, tsMap)
	ds.addTypedStatic(&discover.Node{ID: uintID(4)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(5)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(6)}, tdt)
	dt := dialtest{
		init: ds,
		rounds: []round{
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(1)}}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(4)}, dialType: tdt},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(5)}, dialType: tdt},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(6)}, dialType: tdt},
					&discoverTypedStaticTask{name: tdt, max: 6 - 3},
				},
			},
			{
				// all dialTask are failed to connect and 4, 6 dialTask is expired
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(4)}, dialType: tdt},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(5)}, dialType: tdt},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(6)}, dialType: tdt},
				},
				expired: []failedInfo{
					{discover.Node{ID: uintID(4)}.ID, 4},
					{discover.Node{ID: uintID(6)}.ID, 6},
				},
			},
			{
				// Wait expired task is returned to normal
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn | trustedConn, id: uintID(1)}}},
				},
			},
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn | trustedConn, id: uintID(1)}}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(5)}, dialType: tdt},
				},
			},
		},
	}
	runDialTest(t, dt)
}

func TestDialStateDiscoverTypedStatic(t *testing.T) {
	tdt := dialType("test")

	tsMap := make(map[dialType]typedStatic)
	tsMap[tdt] = typedStatic{maxNodeCount: 6, maxTry: 3}

	ds := newDialState(nil, nil, fakeTable{}, 0, nil, nil, tsMap)
	ds.addTypedStatic(&discover.Node{ID: uintID(4)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(5)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(6)}, tdt)

	dt := dialtest{
		init: ds,
		rounds: []round{
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(1)}}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(4)}, dialType: tdt},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(5)}, dialType: tdt},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(6)}, dialType: tdt},
					&discoverTypedStaticTask{name: tdt, max: 6 - 3},
				},
			},
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(1)}}},
				},
				done: []task{
					&discoverTypedStaticTask{name: tdt, max: 6 - 3},
				},
				new: []task{
					&discoverTypedStaticTask{name: tdt, max: 6 - 3},
				},
			},
		},
	}
	runDialTest(t, dt)
}

func TestDialStateTypedStaticMaxConn(t *testing.T) {
	tdt := dialType("test")

	tsMap := make(map[dialType]typedStatic)
	tsMap[tdt] = typedStatic{maxNodeCount: 3, maxTry: math.MaxInt64}

	ds := newDialState(nil, nil, fakeTable{}, 0, nil, nil, tsMap)
	ds.addTypedStatic(&discover.Node{ID: uintID(1)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(2)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(3)}, tdt) // will be delete cause exceed allowed typed conntection
	dt := dialtest{
		init: ds,
		rounds: []round{
			{
				new: []task{
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(1)}, dialType: tdt},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(2)}, dialType: tdt},
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(3)}, dialType: tdt},
				},
			},
		},
	}
	runDialTest(t, dt)
}

func TestDialStateTypedStaticMaxConn2(t *testing.T) {
	tdt := dialType("test")

	tsMap := make(map[dialType]typedStatic)
	tsMap[tdt] = typedStatic{maxNodeCount: 3, maxTry: math.MaxInt64}

	ds := newDialState(nil, nil, fakeTable{}, 0, nil, nil, tsMap)
	ds.addTypedStatic(&discover.Node{ID: uintID(1)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(2)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(3)}, tdt) // will be delete cause exceed allowed typed conntection
	dt := dialtest{
		init: ds,
		rounds: []round{
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn | trustedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: staticDialedConn | trustedConn, id: uintID(2)}}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn | trustedConn, dest: &discover.Node{ID: uintID(3)}, dialType: tdt},
				},
			},
		},
	}
	runDialTest(t, dt)
}

func TestDialStateTypedStaticMaxConn3(t *testing.T) {
	tdt := dialType("test")

	tsMap := make(map[dialType]typedStatic)
	tsMap[tdt] = typedStatic{maxNodeCount: 3, maxTry: math.MaxInt64}

	ds := newDialState(nil, nil, fakeTable{}, 0, nil, nil, tsMap)
	ds.addTypedStatic(&discover.Node{ID: uintID(1)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(2)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(3)}, tdt) // will be delete cause exceed allowed typed conntection
	ds.addTypedStatic(&discover.Node{ID: uintID(4)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(5)}, tdt)
	ds.addTypedStatic(&discover.Node{ID: uintID(6)}, tdt)
	dt := dialtest{
		init: ds,
		rounds: []round{
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(2)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(3)}}},
				},
				new: []task{},
			},
		},
	}
	runDialTest(t, dt)
}

// This test checks if static dials can ignore adding self ID to static node list.
func TestDialStateAddingSelfNode(t *testing.T) {
	privateKey := newkey()
	nodeID := discover.PubkeyID(&privateKey.PublicKey)

	fakePrivateKey := newkey()

	wantStatic := []*discover.Node{
		{ID: uintID(1)},
		{ID: uintID(2)},
		{ID: nodeID},
		{ID: uintID(3)},
		{ID: uintID(4)},
	}

	dialState_nil := newDialState(wantStatic, nil, fakeTable{}, 0, nil, nil, nil)
	dialState_func := newDialState(wantStatic, nil, fakeTable{}, 0, nil, privateKey, nil)
	dialState_normal := newDialState(wantStatic, nil, fakeTable{}, 0, nil, fakePrivateKey, nil)

	if len(dialState_nil.static) != 5 {
		t.Errorf("newDialState() can't process nil privateKey")
	}

	for _, n := range dialState_func.static {
		if n.dest.ID == nodeID {
			t.Errorf("newDialState() can't ignore adding self node to the static node list")
		}
	}

	if len(dialState_normal.static) != 5 {
		t.Errorf("newDialState() can't deal with normal case")
	}
}

// This test checks that static peers will be redialed immediately if they were re-added to a static list.
func TestDialStaticAfterReset(t *testing.T) {
	wantStatic := []*discover.Node{
		{ID: uintID(1)},
		{ID: uintID(2)},
	}

	rounds := []round{
		// Static dials are launched for the nodes that aren't yet connected.
		{
			peers: nil,
			new: []task{
				&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(1)}, dialType: DT_UNLIMITED},
				&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(2)}, dialType: DT_UNLIMITED},
			},
		},
		// No new dial tasks, all peers are connected.
		{
			peers: []*Peer{
				{rws: []*conn{{flags: staticDialedConn, id: uintID(1)}}},
				{rws: []*conn{{flags: staticDialedConn, id: uintID(2)}}},
			},
			done: []task{
				&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(1)}, dialType: DT_UNLIMITED},
				&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(2)}, dialType: DT_UNLIMITED},
			},
			new: []task{
				&waitExpireTask{Duration: 30 * time.Second},
			},
		},
	}
	dTest := dialtest{
		init:   newDialState(wantStatic, nil, fakeTable{}, 0, nil, nil, nil),
		rounds: rounds,
	}
	runDialTest(t, dTest)
	for _, n := range wantStatic {
		dTest.init.removeStatic(n)
		dTest.init.addStatic(n)
	}
	// without removing peers they will be considered recently dialed
	runDialTest(t, dTest)
}

// This test checks that past dials are not retried for some time.
func TestDialStateCache(t *testing.T) {
	wantStatic := []*discover.Node{
		{ID: uintID(1)},
		{ID: uintID(2)},
		{ID: uintID(3)},
	}

	runDialTest(t, dialtest{
		init: newDialState(wantStatic, nil, fakeTable{}, 0, nil, nil, nil),
		rounds: []round{
			// Static dials are launched for the nodes that
			// aren't yet connected.
			{
				peers: nil,
				new: []task{
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(1)}, dialType: DT_UNLIMITED},
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(2)}, dialType: DT_UNLIMITED},
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(3)}, dialType: DT_UNLIMITED},
				},
			},
			// No new tasks are launched in this round because all static
			// nodes are either connected or still being dialed.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: staticDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: staticDialedConn, id: uintID(2)}}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(1)}, dialType: DT_UNLIMITED},
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(2)}, dialType: DT_UNLIMITED},
				},
			},
			// A salvage task is launched to wait for node 3's history
			// entry to expire.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
				},
				done: []task{
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(3)}},
				},
				new: []task{
					&waitExpireTask{Duration: 14 * time.Second},
				},
			},
			// Still waiting for node 3's entry to expire in the cache.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
				},
			},
			// The cache entry for node 3 has expired and is retried.
			{
				peers: []*Peer{
					{rws: []*conn{{flags: dynDialedConn, id: uintID(1)}}},
					{rws: []*conn{{flags: dynDialedConn, id: uintID(2)}}},
				},
				new: []task{
					&dialTask{flags: staticDialedConn, dest: &discover.Node{ID: uintID(3)}, dialType: DT_UNLIMITED},
				},
			},
		},
	})
}

func TestDialResolve(t *testing.T) {
	resolved := discover.NewNode(uintID(1), net.IP{127, 0, 55, 234}, 3333, 4444, nil, discover.NodeTypeUnknown)
	table := &resolveMock{answer: resolved}
	state := newDialState(nil, nil, table, 0, nil, nil, nil)

	// Check that the task is generated with an incomplete ID.
	dest := discover.NewNode(uintID(1), nil, 0, 0, nil, discover.NodeTypeUnknown)
	state.addStatic(dest)
	tasks := state.newTasks(0, nil, time.Time{})
	if !reflect.DeepEqual(tasks, []task{&dialTask{flags: staticDialedConn, dest: dest, dialType: DT_UNLIMITED}}) {
		t.Fatalf("expected dial task, got %#v", tasks)
	}

	// Now run the task, it should resolve the ID once.
	config := Config{Dialer: TCPDialer{&net.Dialer{Deadline: time.Now().Add(-5 * time.Minute)}}}
	srv := &SingleChannelServer{&BaseServer{ntab: table, Config: config}}
	tasks[0].Do(srv)
	if !reflect.DeepEqual(table.resolveCalls, []discover.NodeID{dest.ID}) {
		t.Fatalf("wrong resolve calls, got %v", table.resolveCalls)
	}

	// Report it as done to the dialer, which should update the static node record.
	state.taskDone(tasks[0], time.Now())
	if state.static[uintID(1)].dest != resolved {
		t.Fatalf("state.dest not updated")
	}
}

// compares task lists but doesn't care about the order.
func sametasks(a, b []task) bool {
	if len(a) != len(b) {
		return false
	}
next:
	for _, ta := range a {
		for _, tb := range b {
			if reflect.DeepEqual(ta, tb) {
				continue next
			}
		}
		return false
	}
	return true
}

func uintID(i uint32) discover.NodeID {
	var id discover.NodeID
	binary.BigEndian.PutUint32(id[:], i)
	return id
}

// implements discoverTable for TestDialResolve
type resolveMock struct {
	resolveCalls []discover.NodeID
	answer       *discover.Node
}

func (t *resolveMock) GetNodes(targetType discover.NodeType, max int) []*discover.Node {
	panic("implement me")
}

func (t *resolveMock) LookupByType(target discover.NodeID, dt discover.DiscoveryType, nType discover.NodeType) []*discover.Node {
	panic("implement me")
}

func (t *resolveMock) Name() string {
	panic("implement me")
}

func (t *resolveMock) RetrieveNodes(target common.Hash, nType discover.NodeType, nresults int) []*discover.Node {
	panic("implement me")
}

func (t *resolveMock) HasBond(id discover.NodeID) bool {
	panic("implement me")
}

func (t *resolveMock) Bond(pinged bool, id discover.NodeID, addr *net.UDPAddr, tcpPort uint16, nType discover.NodeType) (*discover.Node, error) {
	panic("implement me")
}

func (t *resolveMock) CreateUpdateNodeOnDB(n *discover.Node) error {
	panic("implement me")
}

func (t *resolveMock) CreateUpdateNodeOnTable(n *discover.Node) error {
	panic("implement me")
}

func (t *resolveMock) GetNodeFromDB(id discover.NodeID) (*discover.Node, error) {
	panic("implement me")
}

func (t *resolveMock) DeleteNodeFromDB(n *discover.Node) error {
	panic("implement me")
}

func (t *resolveMock) DeleteNodeFromTable(n *discover.Node) error {
	panic("implement me")
}

func (t *resolveMock) GetBucketEntries() []*discover.Node {
	panic("implement me")
}

func (t *resolveMock) GetReplacements() []*discover.Node {
	panic("implement me")
}

func (t *resolveMock) Resolve(target discover.NodeID, nType discover.NodeType) *discover.Node {
	t.resolveCalls = append(t.resolveCalls, target)
	return t.answer
}

func (t *resolveMock) Self() *discover.Node { return new(discover.Node) }

func (t *resolveMock) Close()                     {}
func (t *resolveMock) Bootstrap([]*discover.Node) {}
func (t *resolveMock) Lookup(target discover.NodeID, nType discover.NodeType) []*discover.Node {
	return nil
}
func (t *resolveMock) ReadRandomNodes(buf []*discover.Node, nType discover.NodeType) int { return 0 }
func (t *resolveMock) IsAuthorized(id discover.NodeID, ntype discover.NodeType) bool {
	panic("implement me")
}

func (t *resolveMock) GetAuthorizedNodes() []*discover.Node {
	panic("implement me")
}

func (t *resolveMock) PutAuthorizedNodes(nodes []*discover.Node) {
	panic("implement me")
}

func (t *resolveMock) DeleteAuthorizedNodes(nodes []*discover.Node) {
	panic("implement me")
}
