// Copyright 2022 The klaytn Authors
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

pragma solidity 0.5.6;

contract BridgeVersionController {
    uint64 public VERSION = 0;

    constructor(uint64 version) public {
        VERSION = version;
    }


    // TODO: The currently above tentative implementation should be changed like below. 
    /*
    mapping(address => uint64) internal vMap;

    consturctor(addr address, utin64 version) public {
        require(addr.isContract(), "The given address is not a contract address");
        vMap[addr] = version;
    }
    // Version Getter
    function getVersion(addr address) public view returns(uint64) {
        return vMap[addr];
    }
    */
}
