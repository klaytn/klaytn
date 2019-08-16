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

pragma solidity ^0.4.24;

contract ValidateSenderContract {

    function ValidateSender(address sender, bytes32 msgHash, bytes sigs) public returns (bool) {
        require(sigs.length % 65 == 0);
        bytes memory data = new bytes(20+32+sigs.length);
        uint idx = 0;
        uint i;
        for( i = 0; i < 20; i++) {
            data[idx++] = (bytes20)(sender)[i];
        }
        for( i = 0; i < 32; i++ ) {
            data[idx++] = msgHash[i];
        }
        for( i = 0; i < sigs.length; i++) {
            data[idx++] = sigs[i];
        }
        assembly {
            // skip length header.
            let ptr := add(data, 0x20)
            if iszero(call(gas, 0x0b, 0, ptr, idx, 31, 1)) {
              invalid()
            }
            return(0, 32)
        }

    }
}
