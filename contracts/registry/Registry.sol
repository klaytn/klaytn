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
 * 1. Register: Registers a new system contract.
 *    - Only governance can register a new system contract.
 *    - Can't register a contract with the name of registered or active state contract.
 *
 * 2. Activate: Sets activation block number of the contract.
 *    - Only can be activated by systemTx.
 *    - Only can activate a contract that is in registered state.
 *    - Can't activate a contract from past.
 *
 * 3. Deprecate: Sets deprecate block number of the contract.
 *    - Only can be deprecated by systemTx.
 *    - Only can deprecate a contract that is in active state.
 *    - Can't deprecate a contract from past.
 *
 * 4. Replace: Replaces the contract to new one.
 *    - Only can be replaced by systemTx.
 *    - Deprecate the old contract and activate the new contract at the same block.
 *
 * Code organization
 *    - States (Types)
 *    - Modifiers
 *    - Mutators
 *    - Getters
 */
contract Registry is IRegistry {
    mapping(string => SystemContract) public registry;
    string[] public contractNames;
    address public governance;

    /* ========== MODIFIERS ========== */
    /**
     * @dev Throws if called by any account other than the governance.
     */
    modifier onlyGovernance() {
        require(msg.sender == governance, "Not a governance");
        _;
    }

    /**
     * @dev Throws if not called by systemTx.
     */
    modifier onlySystemTx() {
        // Details will be implemented after finalizing spec of systemTx.
        _;
    }

    /**
     * @dev Throws if the given contract name has already been registered.
     */
    modifier notDuplicate(string memory name) {
        require(registry[name].addr == address(0), "Already registered");
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

    /**
     * @dev Throws if the given address is zero address.
     */
    modifier notZeroAddress(address addr) {
        require(addr != address(0), "Zero address");
        _;
    }

    /**
     * @dev Throws if the given contract is not at given state.
     */
    modifier onlyGivenState(string memory name, State state) {
        require(stateAt(name, block.number) == state, "Not at given state");
        _;
    }

    /* ========== CONSTRUCTOR ========== */
    /**
     * @dev Initializes the contract setting the deployer as the initial governance.
     * @param _governance The address of the governance.
     * TODO: Add secretary authentication logic as below.
     * require(msg.sender == 0x235...., "...");
     */
    function constructContract(address _governance) external notZeroAddress(_governance) {
        governance = _governance;

        emit ConstructContract(_governance);
    }

    /* ========== MUTATORS ========== */
    /**
     * @dev Update governance contract address.
     * @param _newGovernance New governance contract address.
     */
    function updateGovernance(address _newGovernance) external override onlyGovernance notZeroAddress(_newGovernance) {
        governance = _newGovernance;

        emit UpdateGovernance(_newGovernance);
    }

    /**
     * @dev Registers a new system contract to the registry.
     * @param name The name of the contract to register.
     * @param addr The address of the contract to register.
     * @param activation If true, the contract will be activated right after registration.
     */
    function register(
        string memory name,
        address addr,
        bool activation
    ) external override onlyGovernance notDuplicate(name) notEmptyString(name) notZeroAddress(addr) {
        SystemContract storage s = registry[name];

        s.addr = addr;

        // Use next block number to avoid state conflict in the same block.
        if (activation) {
            s.activationBlockNumber = block.number + 1;
        }

        contractNames.push(name);

        emit Registered(name, addr);
    }

    /**
     * @dev Sets activation block number of the contract.
     * @param name The name of the contract to activate.
     * @param activationBlockNumber The activation block number of the contract.
     */
    function activate(
        string memory name,
        uint256 activationBlockNumber
    ) public override onlySystemTx onlyGivenState(name, State.Registered) {
        require(activationBlockNumber >= block.number, "Can't activate contract from past");

        SystemContract storage s = registry[name];

        s.activationBlockNumber = activationBlockNumber;

        emit Activated(name, activationBlockNumber);
    }

    /**
     * @dev Sets deprecate block number of the contract.
     * @param name The name of the contract to deprecate.
     * @param deprecateBlockNumber The deprecate block number of the contract.
     */
    function deprecate(
        string memory name,
        uint256 deprecateBlockNumber
    ) public override onlySystemTx onlyGivenState(name, State.Active) {
        require(deprecateBlockNumber >= block.number, "Can't deprecate contract from past");

        SystemContract storage s = registry[name];

        s.deprecateBlockNumber = deprecateBlockNumber;

        emit Deprecated(name, deprecateBlockNumber);
    }

    /**
     * @dev Replaces the contract with the new contract.
     * @param prevName The name of the previous contract to deprecate.
     * @param newName The name of the new contract to activate.
     * @param replaceBlockNumber The replacement block number of the contract.
     */
    function replace(
        string memory prevName,
        string memory newName,
        uint256 replaceBlockNumber
    ) external override onlySystemTx {
        // Deprecate the previous contract.
        deprecate(prevName, replaceBlockNumber);
        // Activate the new contract.
        activate(newName, replaceBlockNumber);

        emit Replaced(prevName, newName, replaceBlockNumber);
    }

    /* ========== GETTERS ========== */
    /**
     * @dev Returns the state of the contract at given block number.
     * @param name The name of the contract to check.
     * @param blockNumber The block number to check.
     */
    function stateAt(
        string memory name,
        uint256 blockNumber
    ) public view override notEmptyString(name) notZeroAddress(registry[name].addr) returns (State) {
        uint256 activationBlockNumber = registry[name].activationBlockNumber;
        uint256 deprecateBlockNumber = registry[name].deprecateBlockNumber;

        if (activationBlockNumber == 0) {
            return State.Registered;
        } else {
            if (deprecateBlockNumber == 0) {
                if (blockNumber < activationBlockNumber) {
                    return State.Registered;
                } else {
                    return State.Active;
                }
            } else {
                if (blockNumber < activationBlockNumber) {
                    return State.Registered;
                } else if (blockNumber >= activationBlockNumber && blockNumber < deprecateBlockNumber) {
                    return State.Active;
                } else {
                    return State.Depreacted;
                }
            }
        }
    }

    /**
     * @dev Returns the address of contract if active.
     * @param name The name of the contract to check.
     * NOTE: Doesn't need to have modifier since it calls stateAt internally.
     */
    function getContractIfActive(string memory name) external view override returns (address) {
        require(stateAt(name, block.number) == State.Active, "Not active contract");

        return registry[name].addr;
    }

    /**
     * @dev Returns the names and addresses of all contracts at given state.
     * @param blockNumber The block number to check.
     * @param state The state to check.
     */
    function readAllContractsAtGivenState(
        uint256 blockNumber,
        State state
    ) external view override returns (string[] memory, address[] memory) {
        uint256 count = 0;
        for (uint256 i = 0; i < contractNames.length; i++) {
            if (stateAt(contractNames[i], blockNumber) == state) {
                count++;
            }
        }

        string[] memory names = new string[](count);
        address[] memory addrs = new address[](count);

        uint256 index = 0;
        for (uint256 i = 0; i < contractNames.length; i++) {
            if (stateAt(contractNames[i], blockNumber) == state) {
                names[index] = contractNames[i];
                addrs[index] = registry[contractNames[i]].addr;
                index++;
            }
        }

        return (names, addrs);
    }

    /**
     * @dev Returns the all system contract names. (include deprecated contracts)
     */
    function getAllContractNames() external view override returns (string[] memory) {
        return contractNames;
    }
}
