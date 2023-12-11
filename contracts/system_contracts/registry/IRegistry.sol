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

abstract contract IRegistry {
    /* ========== VARIABLES ========== */
    /// The following variables are baked here because their storage layouts matter in protocol consensus
    /// when inject initial states (pre-deployed system contracts, owner) of the Registry.
    /// @dev Mapping of system contracts
    mapping(string => Record[]) public records;

    /// @dev Array of system contract names
    string[] public names;

    /// @dev Owner of contract
    address internal _owner;

    /* ========== TYPES ========== */
    /// @dev Struct of system contracts
    struct Record {
        address addr;
        uint256 activation;
    }

    /* ========== EVENTS ========== */
    /// @dev Emitted when the contract owner is updated by `transferOwnership`.
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    /// @dev Emitted when a new system contract is registered.
    event Registered(string name, address indexed addr, uint256 indexed activation);

    /* ========== MUTATORS ========== */
    /// @dev Registers a new system contract.
    function register(string memory name, address addr, uint256 activation) external virtual;

    /// @dev Transfers ownership to newOwner.
    function transferOwnership(address newOwner) external virtual;

    /* ========== GETTERS ========== */
    /// @dev Returns an address for active system contracts registered as name if exists.
    /// It returns a zero address if there's no active system contract with name.
    function getActiveAddr(string memory name) external virtual returns (address);

    /// @dev Returns all system contracts registered as name.
    function getAllRecords(string memory name) external view virtual returns (Record[] memory);

    /// @dev Returns all names of registered system contracts.
    function getAllNames() external view virtual returns (string[] memory);

    /// @dev Returns owner of contract.
    function owner() external view virtual returns (address);
}
