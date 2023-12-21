// Copyright 2023 The klaytn Authors
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

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

import "./IRegistry.sol";

contract RegistryMock is IRegistry {
    function transferOwnership(address newOwner) external override {
        _owner = newOwner;
    }

    function register(string memory name, address addr, uint256 activation) public override {
        if (records[name].length == 0) {
            names.push(name);
        }
        records[name].push(Record(addr, activation));
    }

    function getActiveAddr(string memory name) external view override returns (address) {
        uint256 len = records[name].length;
        if (len == 0) {
            return address(0);
        } else {
            return records[name][len - 1].addr;
        }
    }

    function getAllRecords(string memory name) external view override returns (Record[] memory) {
        return records[name];
    }

    function getAllNames() external view override returns (string[] memory) {
        return names;
    }

    function owner() external view override returns (address) {
        return _owner;
    }
}
