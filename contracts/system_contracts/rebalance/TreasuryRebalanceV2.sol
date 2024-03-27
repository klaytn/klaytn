// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

import "../lib/Ownable_0.8.sol";

/**
 * @title Interface to get adminlist and quorom
 */
interface IZeroedContract {
    function getState()
        external
        view
        returns (address[] memory adminList, uint256 quorom);
}

/**
 * @title Smart contract to record the rebalance of treasury funds.
 * This contract is to mainly record the addresses which holds the treasury funds
 * before and after rebalancing. It facilates approval and redistributing to new addresses.
 * Core will execute the re-distribution by reading this contract.
 */
contract TreasuryRebalanceV2 is Ownable {
     // Status of the contract
    enum Status {
        Initialized,
        Registered,
        Approved,
        Finalized
    }

    /**
     * Zeroed struct to store zeroed address and their approver addresses
     */
    struct Zeroed {
        address addr;
        address[] approvers;
    }

    /**
     * Allocated struct to receiver address and their fund allocation
     */
    struct Allocated {
        address addr;
        uint256 amount;
    }

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
     * @dev Emitted when a Zeroed is registered
     */
    event ZeroedRegistered(address zeroed);

    /**
     * @dev Emitted when a Zeroed is removed
     */
    event ZeroedRemoved(address zeroed);

    /**
     * @dev Emitted when an Allocated is registered
     */
    event AllocatedRegistered(address allocated, uint256 fundAllocation);

    /**
     * @dev Emitted when an Allocated is removed
     */
    event AllocatedRemoved(address allocated);

    /**
     * @dev Emitted when an admin approves the zeroed address.
     */
    event Approved(address zeroed, address approver, uint256 approversCount);

    /**
     * @dev Emitted when the contract status changes
     */
    event StatusChanged(Status status);

    /**
     * @dev Emitted when the contract is finalized
     * memo - is the result of the treasury fund rebalancing
     */
    event Finalized(string memo, Status status);

    /**
     * Storage
     */
    Zeroed[] public zeroeds; // array of the Zeroed struct
    Allocated[] public allocateds; // array of Allocated struct
    Status public status; // current status of the contract
    uint256 public rebalanceBlockNumber; // the target block number of the execution of rebalancing.
    string public memo; // result of the treasury fund rebalance.

    /**
     * Modifiers
     */
    modifier onlyAtStatus(Status _status) {
        require(status == _status, "Not in the designated status");
        _;
    }

    /**
     *  Constructor
     * @param _rebalanceBlockNumber is the target block number of the execution the rebalance in Core
     */
    constructor(uint256 _rebalanceBlockNumber) {
        require(_rebalanceBlockNumber > block.number, "rebalance blockNumber should be greater than current block");
        rebalanceBlockNumber = _rebalanceBlockNumber;
        status = Status.Initialized;
        emit ContractDeployed(status, _rebalanceBlockNumber, block.timestamp);
    }

    //State changing Functions
    /**
     * @dev updates rebalance block number
     * @param _rebalanceBlockNumber is the updated target block number of the execution the rebalance in Core
     */
    function updateRebalanceBlocknumber(
        uint256 _rebalanceBlockNumber
    ) public onlyOwner {
        require(block.number < rebalanceBlockNumber, "current block shouldn't be past the currently set block number");
        require(block.number < _rebalanceBlockNumber, "rebalance blockNumber should be greater than current block");
        rebalanceBlockNumber = _rebalanceBlockNumber;
    }

    /**
     * @dev registers zeroed details
     * @param _zeroedAddress is the address of the zeroed
     */
    function registerZeroed(
        address _zeroedAddress
    ) public onlyOwner onlyAtStatus(Status.Initialized) {
        require(
            !zeroedExists(_zeroedAddress),
            "Zeroed address is already registered"
        );
        Zeroed storage zeroed = zeroeds.push();
        zeroed.addr = _zeroedAddress;
        emit ZeroedRegistered(zeroed.addr);
    }

    /**
     * @dev remove the zeroed details from the array
     * @param _zeroedAddress is the address of the zeroed
     */
    function removeZeroed(
        address _zeroedAddress
    ) public onlyOwner onlyAtStatus(Status.Initialized) {
        uint256 zeroedIndex = getZeroedIndex(_zeroedAddress);
        require(zeroedIndex != type(uint256).max, "Zeroed not registered");
        zeroeds[zeroedIndex] = zeroeds[zeroeds.length - 1];
        zeroeds.pop();
        emit ZeroedRemoved(_zeroedAddress);
    }

    /**
     * @dev registers allocated address and its fund distribution
     * @param _allocatedAddress is the address of the allocated
     * @param _amount is the fund to be allocated to the allocated
     */
    function registerAllocated(
        address _allocatedAddress,
        uint256 _amount
    ) public onlyOwner onlyAtStatus(Status.Initialized) {
        require(
            !allocatedExists(_allocatedAddress),
            "Allocated address is already registered"
        );
        require(_amount != 0, "Amount cannot be set to 0");
        Allocated memory allocated = Allocated(_allocatedAddress, _amount);
        allocateds.push(allocated);
        emit AllocatedRegistered(_allocatedAddress, _amount);
    }

    /**
     * @dev remove the allocated details from the array
     * @param _allocatedAddress is the address of the allocated
     */
    function removeAllocated(
        address _allocatedAddress
    ) public onlyOwner onlyAtStatus(Status.Initialized) {
        uint256 allocatedIndex = getAllocatedIndex(_allocatedAddress);
        require(allocatedIndex != type(uint256).max, "Allocated not registered");
        allocateds[allocatedIndex] = allocateds[allocateds.length - 1];
        allocateds.pop();
        emit AllocatedRemoved(_allocatedAddress);
    }

    /**
     * @dev zeroedAddress can be a EOA or a contract address. To approve:
     *      If the zeroedAddress is a EOA, the msg.sender should be the EOA address
     *      If the zeroedAddress is a Contract, the msg.sender should be one of the contract `admin`.
     *      It uses the getState() function in the zeroedAddress contract to get the admin details.
     * @param _zeroedAddress is the address of the zeroed
     */
    function approve(
        address _zeroedAddress
    ) public onlyAtStatus(Status.Registered) {
        require(
            zeroedExists(_zeroedAddress),
            "zeroed needs to be registered before approval"
        );
        //Check whether the zeroed address is EOA or contract address
        bool isContract = isContractAddr(_zeroedAddress);
        if (!isContract) {
            //check whether the msg.sender is the zeroed if its a EOA
            require(
                msg.sender == _zeroedAddress,
                "zeroedAddress is not the msg.sender"
            );
            _updateApprover(_zeroedAddress, msg.sender);
        } else {
            (address[] memory adminList, uint256 req) = _getState(_zeroedAddress);
            require(adminList.length != 0, "admin list cannot be empty");
            //check if the msg.sender is one of the admin of the zeroedAddress contract
            require(
                _validateAdmin(msg.sender, adminList),
                "msg.sender is not the admin"
            );
            _updateApprover(_zeroedAddress, msg.sender);
        }
    }

    /**
     * @dev validate if the msg.sender is admin if the zeroedAddress is a contract
     * @param _approver is the msg.sender
     * @return isAdmin is true if the msg.sender is one of the admin
     */
    function _validateAdmin(
        address _approver,
        address[] memory _adminList
    ) private pure returns (bool isAdmin) {
        for (uint256 i = 0; i < _adminList.length; i++) {
            if (_approver == _adminList[i]) {
                isAdmin = true;
            }
        }
    }

    /**
     * @dev gets the adminList and quorom by calling `getState()` method in zeroedAddress contract
     * @param _zeroedAddress is the address of the contract
     * @return adminList list of the zeroedAddress contract admins
     * @return req min required number of approvals
     */
    function _getState(
        address _zeroedAddress
    ) private view returns (address[] memory adminList, uint256 req) {
        IZeroedContract zeroedContract = IZeroedContract(_zeroedAddress);
        (adminList, req) = zeroedContract.getState();
    }

    /**
     * @dev Internal function to update the approver details of a zeroed
     * _zeroedAddress is the address of the zeroed
     * _approver is the admin of the zeroedAddress
     */
    function _updateApprover(
        address _zeroedAddress,
        address _approver
    ) private {
        uint256 index = getZeroedIndex(_zeroedAddress);
        require(index != type(uint256).max, "Zeroed not registered");
        address[] memory approvers = zeroeds[index].approvers;
        for (uint256 i = 0; i < approvers.length; i++) {
            require(approvers[i] != _approver, "Already approved");
        }
        zeroeds[index].approvers.push(_approver);
        emit Approved(
            _zeroedAddress,
            _approver,
            zeroeds[index].approvers.length
        );
    }

    /**
     * @dev finalizeRegistration sets the status to Registered,
     *      After this stage, registrations will be restricted.
     */
    function finalizeRegistration()
        public
        onlyOwner
        onlyAtStatus(Status.Initialized)
    {
        status = Status.Registered;
        emit StatusChanged(status);
    }

    /**
     * @dev finalizeApproval sets the status to Approved,
     *      After this stage, approvals will be restricted.
     */
    function finalizeApproval()
        public
        onlyOwner
        onlyAtStatus(Status.Registered)
    {
        require(
            getTreasuryAmount() < sumOfZeroedBalance(),
            "treasury amount should be less than the sum of all zeroed address balances"
        );
        checkZeroedsApproved();
        status = Status.Approved;
        emit StatusChanged(status);
    }

    /**
     * @dev verify if quorom reached for the zeroed approvals
     */
    function checkZeroedsApproved() public view {
        for (uint256 i = 0; i < zeroeds.length; i++) {
            Zeroed memory zeroed = zeroeds[i];
            bool isContract = isContractAddr(zeroed.addr);
            if (isContract) {
                (address[] memory adminList, uint256 req) = _getState(
                    zeroed.addr
                );
                require(
                    zeroed.approvers.length >= req,
                    "min required admins should approve"
                );
                //if min quorom reached, make sure all approvers are still valid
                address[] memory approvers = zeroed.approvers;
                uint256 validApprovals = 0;
                for (uint256 j = 0; j < approvers.length; j++) {
                    if (_validateAdmin(approvers[j], adminList)) {
                        validApprovals++;
                    }
                }
                require(
                    validApprovals >= req,
                    "min required admins should approve"
                );
            } else {
                require(zeroed.approvers.length == 1, "EOA should approve");
            }
        }
    }

    /**
     * @dev sets the status of the contract to Finalize. Once finalized the storage data
     * of the contract cannot be modified
     * @param _memo is the result of the rebalance after executing successfully in the core.
     * Can only be called by the current owner at Approved state after the execution of rebalance in the core
     *  - memo format: { "zeroeds": [ { "zeroed": "0xaddr", "balance": 0xamount },
     *                 { "zeroed": "0xaddr", "balance": 0xamount }, ... ],
     *                 "allocateds": [ { "allocated": "0xaddr", "fundAllocated": 0xamount },
     *                 { "allocated": "0xaddr", "fundAllocated": 0xamount }, ... ],
     *                 "burnt": 0xamount, "success": true/false }
     */
    function finalizeContract(
        string memory _memo
    ) public onlyOwner onlyAtStatus(Status.Approved) {
        memo = _memo;
        status = Status.Finalized;
        emit Finalized(memo, status);
        require(
            block.number > rebalanceBlockNumber,
            "Contract can only finalize after executing rebalancing"
        );
    }

    /**
     * @dev resets all storage values to empty objects except targetBlockNumber
     */
    function reset() public onlyOwner {
        //reset cannot be called at Finalized status or after target block.number
        require(status != Status.Finalized, "Contract is finalized, cannot reset values");
        require (block.number < rebalanceBlockNumber, "Rebalance blocknumber already passed");

        //`delete` keyword is used to set a storage variable or a dynamic array to its default value.
        delete zeroeds;
        delete allocateds;
        delete memo;
        status = Status.Initialized;
    }

    //Getters
    /**
     * @dev to get zeroed details by zeroedAddress
     * @param _zeroedAddress is the address of the zeroed
     */
    function getZeroed(
        address _zeroedAddress
    ) public view returns (address, address[] memory) {
        uint256 index = getZeroedIndex(_zeroedAddress);
        require(index != type(uint256).max, "Zeroed not registered");
        Zeroed memory zeroed = zeroeds[index];
        return (zeroed.addr, zeroed.approvers);
    }

    /**
     * @dev check whether zeroedAddress is registered
     * @param _zeroedAddress is the address of the zeroed
     */
    function zeroedExists(address _zeroedAddress) public view returns (bool) {
        require(_zeroedAddress != address(0), "Invalid address");
        for (uint256 i = 0; i < zeroeds.length; i++) {
            if (zeroeds[i].addr == _zeroedAddress) {
                return true;
            }
        }
    }

    /**
     * @dev get index of the zeroed in the zeroeds array
     * @param _zeroedAddress is the address of the zeroed
     */
    function getZeroedIndex(
        address _zeroedAddress
    ) public view returns (uint256) {
        for (uint256 i = 0; i < zeroeds.length; i++) {
            if (zeroeds[i].addr == _zeroedAddress) {
                return i;
            }
        }
        return type(uint256).max;
    }

    /**
     * @dev to calculate the sum of zeroeds balances
     * @return zeroedsBalance the sum of balances of zeroeds
     */
    function sumOfZeroedBalance()
        public
        view
        returns (uint256 zeroedsBalance)
    {
        for (uint256 i = 0; i < zeroeds.length; i++) {
            zeroedsBalance += zeroeds[i].addr.balance;
        }
        return zeroedsBalance;
    }

    /**
     * @dev to get allocated details by allocatedAddress
     * @param _allocatedAddress is the address of the allocated
     * @return allocated is the address of the allocated
     * @return amount is the fund allocated to the allocated
     */
    function getAllocated(
        address _allocatedAddress
    ) public view returns (address, uint256) {
        uint256 index = getAllocatedIndex(_allocatedAddress);
        require(index != type(uint256).max, "Allocated not registered");
        Allocated memory allocated = allocateds[index];
        return (allocated.addr, allocated.amount);
    }

    /**
     * @dev check whether _allocatedAddress is registered
     * @param _allocatedAddress is the address of the allocated
     */
    function allocatedExists(address _allocatedAddress) public view returns (bool) {
        require(_allocatedAddress != address(0), "Invalid address");
        for (uint256 i = 0; i < allocateds.length; i++) {
            if (allocateds[i].addr == _allocatedAddress) {
                return true;
            }
        }
    }

    /**
     * @dev get index of the allocated in the allocateds array
     * @param _allocatedAddress is the address of the allocated
     */
    function getAllocatedIndex(
        address _allocatedAddress
    ) public view returns (uint256) {
        for (uint256 i = 0; i < allocateds.length; i++) {
            if (allocateds[i].addr == _allocatedAddress) {
                return i;
            }
        }
        return type(uint256).max;
    }

    /**
     * @dev to calculate the sum of allocated funds
     * @return treasuryAmount the sum of funds allocated to allocateds
     */
    function getTreasuryAmount() public view returns (uint256 treasuryAmount) {
        for (uint256 i = 0; i < allocateds.length; i++) {
            treasuryAmount += allocateds[i].amount;
        }
        return treasuryAmount;
    }

    /**
     * @dev gets the length of zeroeds list
     */
    function getZeroedCount() public view returns (uint256) {
        return zeroeds.length;
    }

    /**
     * @dev gets the length of allocateds list
     */
    function getAllocatedCount() public view returns (uint256) {
        return allocateds.length;
    }

    /**
     * @dev fallback function to revert any payments
     */
    fallback() external payable {
        revert("This contract does not accept any payments");
    }

    /**
     * @dev Helper function to check the address is contract addr or EOA
     */
    function isContractAddr(address _addr) public view returns (bool) {
        uint256 size;
        assembly {
            size := extcodesize(_addr)
        }
        return size > 0;
    }
}
