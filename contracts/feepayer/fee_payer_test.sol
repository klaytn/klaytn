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

pragma solidity >0.8.0;

contract FeePayer {
    function GetFeePayerDirect() public returns (address) {
        assembly {
            if iszero(call(gas(), 0x0a, 0, 0, 0, 12, 20)) {
              invalid()
            }
            return(0, 32)
        }
    }

    function feePayer() internal returns (address addr) {
        assembly {
            let freemem := mload(0x40)
            let start_addr := add(freemem, 12)
            if iszero(call(gas(), 0x0a, 0, 0, 0, start_addr, 20)) {
              invalid()
            }
            addr := mload(freemem)
        }
    }

    function GetFeePayer() public returns (address) {
        return feePayer();
    }
}
