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

// revertTracer returns the string of REVERT.
// If not reverted, returns an empty string "".
{
    revertString: "",

    toAscii: function(hex) {
        var str = "";
        var i = 0, l = hex.length;
        if (hex.substring(0, 2) === "0x") {
            i = 2;
        }
        for (; i < l; i+=2) {
            var code = parseInt(hex.substr(i, 2), 16);
            str += String.fromCharCode(code);
        }

        return str;
    },

    // step is invoked for every opcode that the VM executes.
    step: function(log, db) { },

    // fault is invoked when the actual execution of an opcode fails.
    fault: function(log, db) { },

    // result is invoked when all the opcodes have been iterated over and returns
    // the final result of the tracing.
    result: function(ctx, db) {
        // Revert output example:
        // 0x08c379a0
        //   0000000000000000000000000000000000000000000000000000000000000020 stringOffset
        //   0000000000000000000000000000000000000000000000000000000000000008 stringLength
        //   4141414141414141
        if (ctx.error == "evm: execution reverted") {
            outputHex = toHex(ctx.output);
            if (outputHex.slice(2,10) == "08c379a0") {
                defaultOffset = 10;
                stringOffsetHex = "0x"+outputHex.slice(defaultOffset, defaultOffset + 32*2);
                stringLengthHex = "0x"+outputHex.slice(defaultOffset + 32*2, defaultOffset + 32*2 + 32*2);
                try {
                    stringOffset = parseInt(bigInt(stringOffsetHex).toString());
                    stringLength = parseInt(bigInt(stringLengthHex).toString());
                    start = defaultOffset + 32*2 + stringOffset*2;
                    end = start + stringLength*2;
                    this.revertString = this.toAscii(outputHex.slice(start,end));
                } catch (e) {
                    this.revertString = "";
                }
            }
        }
        return this.revertString;
    }
}
