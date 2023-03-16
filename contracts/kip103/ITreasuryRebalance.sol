// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

/**
 * @dev External interface of TreasuryRebalance
 */
interface ITreasuryRebalance {
    /**
     * @dev Emitted when the contract is deployed
     * `rebalanceBlockNumber` is the target block number of the execution the rebalance in Core
     * `deployedBlockNumber` is the current block number when its deployed
     */
    event ContractDeployed(
        Status status,
        uint256 rebalanceBlockNumber,
        uint256 deployedBlockNumber
    );

    /**
     * @dev Emitted when a Retired is registered
     */
    event RetiredRegistered(address retired);

    /**
     * @dev Emitted when a Retired is removed
     */
    event RetiredRemoved(address retired);

    /**
     * @dev Emitted when a Newbie is registered
     */
    event NewbieRegistered(address newbie, uint256 fundAllocation);

    /**
     * @dev Emitted when a Newbie is removed
     */
    event NewbieRemoved(address newbie);

    /**
     * @dev Emitted when a admin approves the retired address.
     */
    event Approved(address retired, address approver, uint256 approversCount);

    /**
     * @dev Emitted when the contract status changes
     */
    event StatusChanged(Status status);

    /**
     * @dev Emitted when the contract is finalized
     * memo - is the result of the treasury fund rebalancing
     */
    event Finalized(string memo, Status status);

    // Status of the contract
    enum Status {
        Initialized,
        Registered,
        Approved,
        Finalized
    }

    /**
     * Retired struct to store retired address and their approver addresses
     */
    struct Retired {
        address retired;
        address[] approvers;
    }

    /**
     * Newbie struct to newbie receiver address and their fund allocation
     */
    struct Newbie {
        address newbie;
        uint256 amount;
    }

    // State variables
    function status() external view returns (Status); // current status of the contract

    function rebalanceBlockNumber() external view returns (uint256); // the target block number of the execution of rebalancing

    function memo() external view returns (string memory); // result of the treasury fund rebalance

    /**
     * @dev to get retired details by retiredAddress
     */
    function getRetired(
        address retiredAddress
    ) external view returns (address, address[] memory);

    /**
     * @dev to get newbie details by newbieAddress
     */
    function getNewbie(
        address newbieAddress
    ) external view returns (address, uint256);

    /**
     * @dev returns the sum of retirees balances
     */
    function sumOfRetiredBalance()
        external
        view
        returns (uint256 retireesBalance);

    /**
     * @dev returns the sum of newbie funds
     */
    function getTreasuryAmount() external view returns (uint256 treasuryAmount);

    /**
     * @dev returns the length of retirees list
     */
    function getRetiredCount() external view returns (uint256);

    /**
     * @dev returns the length of newbies list
     */
    function getNewbieCount() external view returns (uint256);

    /**
     * @dev verify all retirees are approved by admin
     */
    function checkRetiredsApproved() external view;

    /**
     * @dev get the current value of the state variables
     */
    function getSnapshot()
        external
        view
        returns (
            Retired[] memory retirees,
            Newbie[] memory newbies,
            uint256 totalRetireesBalance,
            uint256 totalNewbiesAllocation,
            Status status,
            string memory memo
    );

    // State changing functions
    /**
     * @dev registers retired details
     * Can only be called by the current owner at Initialized state
     */
    function registerRetired(address retiredAddress) external;

    /**
     * @dev remove the retired details from the array
     * Can only be called by the current owner at Initialized state
     */
    function removeRetired(address retiredAddress) external;

    /**
     * @dev registers newbie address and its fund distribution
     * Can only be called by the current owner at Initialized state
     */
    function registerNewbie(address newbieAddress, uint256 amount) external;

    /**
     * @dev remove the newbie details from the array
     * Can only be called by the current owner at Initialized state
     */
    function removeNewbie(address newbieAddress) external;

    /**
     * @dev approves a retiredAddress,the address can be a EOA or a contract address.
     *  - If the retiredAddress is a EOA, the caller should be the EOA address
     *  - If the retiredAddress is a Contract, the caller should be one of the contract `admin`
     */
    function approve(address retiredAddress) external;

    /**
     * @dev sets the status to Registered,
     *      After this stage, registrations will be restricted.
     * Can only be called by the current owner at Initialized state
     */
    function finalizeRegistration() external;

    /**
     * @dev sets the status to Approved,
     * Can only be called by the current owner at Registered state
     */
    function finalizeApproval() external;

    /**
     * @dev sets the status of the contract to Finalize. Once finalized the storage data
     * of the contract cannot be modified
     * Can only be called by the current owner at Approved state after the execution of rebalance in the core
     *  - memo format: {"success": true/false, "retried":[{"0xaddr": "0xamount"}, ...],
     *     "newbie": [{"0xaddr":"0xamount"}, ...], "bunt": "0xamount"}
     */
    function finalizeContract(string memory memo) external;

    /**
     * @dev resets all storage values to empty objects except targetBlockNumber
     */
    function reset() external;
}
