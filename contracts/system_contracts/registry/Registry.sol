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

/**
 * @dev Registry is a contract that manages the addresses of system contracts.
 * Note: The pre-deployed system contracts will be directly injected into the registry in HF block.
 *
 * register: Registers a new system contract.
 *    - Only can be registered by governance.
 *    - If predecessor is not yet active, overwrite it.
 *
 * Code organization
 *    - Modifiers
 *    - Mutators
 *    - Getters
 */
contract Registry is IRegistry {
    /* ========== MODIFIERS ========== */
    // /**
    //  * @dev Throws if not called by systemTx.
    //  * TODO: Decide whether to use this modifier or not.
    //  */
    // modifier onlySystemTx() {
    //     _;
    // }

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(msg.sender == owner(), "Not owner");
        _;
    }

    /**
     * @dev Throws if the given string is empty.
     */
    modifier notEmptyString(string memory name) {
        bytes memory b = abi.encodePacked(name);
        require(b.length != 0, "Empty string");
        _;
    }

    /* ========== MUTATORS ========== */
    /**
     * @dev Registers a new system contract to the records.
     * @param name The name of the contract to register.
     * @param addr The address of the contract to register.
     * @param activation The activation block number of the contract.
     * NOTE: Register a zero address if you want to deprecate the contract without replacing it.
     */
    function register(
        string memory name,
        address addr,
        uint256 activation
    ) external override onlyOwner notEmptyString(name) {
        // Don't allow the current block since it affects to other txs in the same block.
        require(activation > block.number, "Can't register contract from past");

        uint256 length = records[name].length;

        if (length == 0) {
            names.push(name);
            records[name].push(Record(addr, activation));
        } else {
            Record storage last = records[name][length - 1];
            if (last.activation <= block.number) {
                // Last record is active. Append new record.
                records[name].push(Record(addr, activation));
            } else {
                // Last record is not yet active. Overwrite last record.
                last.addr = addr;
                last.activation = activation;
            }
        }

        emit Registered(name, addr, activation);
    }

    /**
     * @dev Transfers ownership of the contract to a newOwner.
     * @param newOwner The address to transfer ownership to.
     */
    function transferOwnership(address newOwner) external override onlyOwner {
        require(newOwner != address(0), "Zero address");
        _owner = newOwner;

        emit OwnershipTransferred(msg.sender, newOwner);
    }

    /* ========== GETTERS ========== */
    /**
     * @dev Returns the address of contract if active at current block.
     * @param name The name of the contract to check.
     * Note: If there is no active contract, it returns address(0).
     */
    function getActiveAddr(string memory name) public view virtual override returns (address) {
        uint256 length = records[name].length;

        // activation is always in ascending order.
        for (uint256 i = length; i > 0; i--) {
            if (records[name][i - 1].activation <= block.number) {
                return records[name][i - 1].addr;
            }
        }

        return address(0);
    }

    /**
     * @dev Returns all contract with same name.
     * @param name The name of the contract to check.
     */
    function getAllRecords(string memory name) public view override returns (Record[] memory) {
        return records[name];
    }

    /**
     * @dev Returns the all system contract names. (include deprecated contracts)
     */
    function getAllNames() public view override returns (string[] memory) {
        return names;
    }

    /**
     * @dev Returns the owner of the contract.
     */
    function owner() public view override returns (address) {
        return _owner;
    }
}
