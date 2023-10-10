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

import "./IKIP113.sol";

contract KIP113Mock is IKIP113 {
    function registerPublicKey(bytes calldata publicKey, bytes calldata pop) external virtual override {
        address addr = msg.sender;

        if (infos[addr].publicKey.length == 0) {
            addrs.push(addr);
        }

        infos[addr] = BlsPublicKeyInfo(publicKey, pop);
    }

    function unregisterPublicKey(address addr) external virtual override {
        delete infos[addr];
    }

    function getInfo(address addr) external virtual override view returns (BlsPublicKeyInfo memory pubkey) {
        return infos[addr];
    }

    function getAllInfo() external virtual override view returns (address[] memory addrList, BlsPublicKeyInfo[] memory pubkeyList) {
        uint count = addrs.length;

        addrList = new address[](count);
        pubkeyList = new BlsPublicKeyInfo[](count);

        for (uint i = 0; i < count; i++) {
            addrList[i] = addrs[i];
            pubkeyList[i] = infos[addrs[i]];
        }
        return (addrList, pubkeyList);
    }
}

