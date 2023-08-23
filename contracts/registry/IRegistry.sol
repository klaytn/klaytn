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

interface IRegistry {
  /* ========== TYPES ========== */
  // State of contract
  enum State {
    Registered,
    Active,
    Depreacted
  }

  struct SystemContract {
    address addr;
    uint256 activationBlockNumber;
    uint256 deprecateBlockNumber;
  }

  /* ========== EVENTS ========== */
  event ConstructContract(address indexed governance);
  event UpdateGovernance(address indexed newGovernance);
  event Registered(string name, address indexed addr);
  event Activated(string name, uint256 indexed activationBlockNumber);
  event Deprecated(string name, uint256 indexed deprecateBlockNumber);
  event Replaced(string prevName, string newName, uint256 indexed replaceBlockNumber);

  /* ========== MUTATORS ========== */
  // Update governance contract address
  function updateGovernance(address newGovernance) external;

  // Register a new contract by governance
  function register(string memory name, address addr, bool activation) external;

  // Activate/Deprecate/Replace contracts by systemTx
  function activate(string memory name, uint256 activationBlockNumber) external;

  function deprecate(string memory name, uint256 deprecateBlockNumber) external;

  function replace(string memory prevName, string memory newName, uint256 replaceBlockNumber) external;

  /* ========== GETTERS ========== */
  function registry(string memory name) external returns (address, uint256, uint256);

  function contractNames(uint256 index) external returns (string memory);

  function getContractIfActive(string memory name) external returns (address);

  function stateAt(string memory name, uint256 blockNumber) external view returns (State);

  function readAllContractsAtGivenState(
    uint256 blockNumber,
    State state
  ) external view returns (string[] memory, address[] memory);

  function getAllContractNames() external view returns (string[] memory);
}
