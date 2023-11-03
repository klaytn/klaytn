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
pragma solidity ^0.8.18;

import "./KIP113.sol";

contract KIP113Mock is KIP113 {
    function register(address addr, bytes calldata publicKey, bytes calldata pop) external override {
        if (record[addr].publicKey.length == 0) {
            allNodeIds.push(addr);
        }
        record[addr] = BlsPublicKeyInfo(publicKey, pop);
    }

    function getAllBlsInfo()
        external
        view
        virtual
        override
        returns (address[] memory nodeIdList, BlsPublicKeyInfo[] memory pubkeyList)
    {
        uint count = allNodeIds.length;

        nodeIdList = new address[](count);
        pubkeyList = new BlsPublicKeyInfo[](count);

        for (uint i = 0; i < count; i++) {
            nodeIdList[i] = allNodeIds[i];
            pubkeyList[i] = record[allNodeIds[i]];
        }
        return (nodeIdList, pubkeyList);
    }
}
