// Modifications Copyright 2020 The klaytn Authors
// Copyright 2018 The go-ethereum Authors
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
// This file is derived from eth/tracers/internal/tracers/trigram_tracer.js (2020/11/18).
// Modified and improved for the klaytn development.
//
// trigram returns trigram(a group of three consecutive written units such as letters, syllables, or words.) counting number
// Example:
//  > debug.traceTransaction( "0x214e597e35da083692f5386141e69f47e973b2c56e7a8073b1ea08fd7571e9de", {tracer: "trigramTracer"})
//  {
//     --PUSH1: 1,
//     -PUSH1-MSTORE: 1,
//     ADD-SWAP1-DUP1: 1,
//     ADD-SWAP1-SWAP3: 1,
//     CALLDATALOAD-PUSH1-SHR: 1,
//     CALLDATALOAD-SWAP1-PUSH1: 1,
//  }
{
    // hist is the map of trigram counters
    hist: {},
    // lastOp is last operation
    lastOps: ['',''],
        lastDepth: 0,
    // step is invoked for every opcode that the VM executes.
    step: function(log, db) {
    var depth = log.getDepth();
    if (depth != this.lastDepth){
        this.lastOps = ['',''];
        this.lastDepth = depth;
        return;
    }
    var op = log.op.toString();
    var key = this.lastOps[0]+'-'+this.lastOps[1]+'-'+op;
    if (this.hist[key]){
        this.hist[key]++;
    }
    else {
        this.hist[key] = 1;
    }
    this.lastOps[0] = this.lastOps[1];
    this.lastOps[1] = op;
},
    // fault is invoked when the actual execution of an opcode fails.
    fault: function(log, db) {},
    // result is invoked when all the opcodes have been iterated over and returns
    // the final result of the tracing.
    result: function(ctx) {
        return this.hist;
    },
}
