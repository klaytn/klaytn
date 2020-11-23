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
// This file is derived from eth/tracers/internal/tracers/unigram_tracer.js (2020/11/18).
// Modified and improved for the klaytn development.
//
// unigramTracer returns opcode counting number
// Example:
//  > debug.traceTransaction( "0x214e597e35da083692f5386141e69f47e973b2c56e7a8073b1ea08fd7571e9de", {tracer: "unigramTracer"})
// {
//     ADD: 2,
//     CALLDATALOAD: 2,
//     CALLDATASIZE: 2,
//     CALLVALUE: 1,
//     DUP1: 8,
//     DUP2: 3,
// }
{
    // hist is the map of opcodes to counters
    hist: {},
    // nops counts number of ops
    nops: 0,
        // step is invoked for every opcode that the VM executes.
        step: function(log, db) {
    var op = log.op.toString();
    if (this.hist[op]){
        this.hist[op]++;
    }
    else {
        this.hist[op] = 1;
    }
    this.nops++;
},
    // fault is invoked when the actual execution of an opcode fails.
    fault: function(log, db) {},

    // result is invoked when all the opcodes have been iterated over and returns
    // the final result of the tracing.
    result: function(ctx) {
        if(this.nops > 0){
            return this.hist;
        }
    },
}
