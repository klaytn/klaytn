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

import "./IKIP113.sol";
import "./IAddressBook.sol";
import "../lib/Initializable.sol";
import "../lib/UUPSUpgradeable.sol";
import "../lib/OwnableUpgradeable.sol";

contract KIP113 is Initializable, UUPSUpgradeable, OwnableUpgradeable, IKIP113 {
    IAddressBook public constant abook = IAddressBook(0x0000000000000000000000000000000000000400);
    bytes32 public constant ZERO48HASH = 0xc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd293; // keccak256(hex"00"*48)
    bytes32 public constant ZERO96HASH = 0x46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21; // keccak256(hex"00"*96)

    address[] public allNodeIds;
    mapping(address => BlsPublicKeyInfo) public record; // cnNodeId => BlsPublicKeyInfo

    event Registered(address cnNodeId, bytes publicKey, bytes pop);
    event Unregistered(address cnNodeId, bytes publicKey, bytes pop);

    modifier onlyValidPublicKey(bytes calldata publicKey) {
        require(publicKey.length == 48, "Public key must be 48 bytes");
        require(keccak256(publicKey) != ZERO48HASH, "Public key cannot be zero");
        // TODO: verify more with EIP-2537
        _;
    }

    modifier onlyValidPop(bytes calldata pop) {
        require(pop.length == 96, "Pop must be 96 bytes");
        require(keccak256(pop) != ZERO96HASH, "Pop cannot be zero");
        // TODO: verify more with EIP-2537
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize() public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();
    }

    function register(
        address cnNodeId,
        bytes calldata publicKey,
        bytes calldata pop
    ) external virtual onlyOwner onlyValidPublicKey(publicKey) onlyValidPop(pop) {
        require(isCN(cnNodeId), "cnNodeId is not in AddressBook");
        if (record[cnNodeId].publicKey.length == 0) {
            allNodeIds.push(cnNodeId);
        }

        record[cnNodeId] = BlsPublicKeyInfo(publicKey, pop);
        emit Registered(cnNodeId, publicKey, pop);
    }

    function unregister(address cnNodeId) external virtual onlyOwner {
        require(!isCN(cnNodeId), "CN is still in AddressBook");
        require(record[cnNodeId].publicKey.length != 0, "CN is not registered");

        _removeCnNodeId(cnNodeId);
        emit Unregistered(cnNodeId, record[cnNodeId].publicKey, record[cnNodeId].pop);
        delete record[cnNodeId];
    }

    function _removeCnNodeId(address cnNodeId) private {
        for (uint256 i = 0; i < allNodeIds.length; i++) {
            if (allNodeIds[i] == cnNodeId) {
                allNodeIds[i] = allNodeIds[allNodeIds.length - 1];
                allNodeIds.pop();
                return;
            }
        }
    }

    function getAllBlsInfo()
        external
        view
        virtual
        returns (address[] memory nodeIdList, BlsPublicKeyInfo[] memory pubkeyList)
    {
        nodeIdList = new address[](allNodeIds.length);
        pubkeyList = new BlsPublicKeyInfo[](allNodeIds.length);

        for (uint256 i = 0; i < nodeIdList.length; i++) {
            nodeIdList[i] = allNodeIds[i];
            pubkeyList[i] = record[allNodeIds[i]];
        }
    }

    function isCN(address target) internal view virtual returns (bool) {
        try abook.getCnInfo(target) {
            return true;
        } catch {
            return false;
        }
    }

    function _authorizeUpgrade(address newImplementation) internal virtual override onlyOwner {}
}
