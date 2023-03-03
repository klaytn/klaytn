// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

import "./Ownable.sol";

/**
 * @title Smart contract to record the rebalance of treasury funds.
 * This contract is to mainly record the addresses which holds the treasury funds
 * before and after rebalancing. It facilates approval and redistributing to new addresses.
 * Core will execute the re-distribution by reading this contract.
 */
contract TreasuryRebalance is Ownable {
    /**
     *  Enums to track the status of the contract
     */
    enum Status {
        Initialized,
        Registered,
        Approved,
        Finalized
    }

    /**
     * Retired struct to store retired addresses and their approver addresses
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

    /**
     * Storage
     */
    Retired[] public retirees; // array of the Retired struct
    Newbie[] public newbies; // array of Newbie struct
    Status public status; // current status of the contract
    uint256 public rebalanceBlockNumber; // the target block number of the execution of rebalancing.
    string public memo; // result of the treasury fund rebalance. eg : `Treasury Rebalanced Successfully`

    /**
     * Events logs
     */
    event ContractDeployed(
        Status status,
        uint256 rebalanceBlockNumber,
        uint256 deployedBlockNumber
    );
    event RetiredRegistered(address retired, address[] approvers);
    event RetiredRemoved(address retired, uint256 retiredCount);
    event NewbieRegistered(address newbie, uint256 fundAllocation);
    event NewbieRemoved(address newbie, uint256 newbieCount);
    event GetState(bool success, bytes result);
    event Approved(address retired, address approver, uint256 approversCount);
    event SetStatus(Status status);
    event Finalized(string memo, Status status);

    /**
     * Modifiers
     */
    modifier onlyAtStatus(Status _status) {
        require(status == _status, "function not allowed at this stage");
        _;
    }

    /**
     *  Constructor
     * @param _rebalanceBlockNumber is the target block number of the execution the rebalance in Core
     */
    constructor(uint256 _rebalanceBlockNumber) {
        rebalanceBlockNumber = _rebalanceBlockNumber;
        status = Status.Initialized;
        emit ContractDeployed(status, _rebalanceBlockNumber, block.timestamp);
    }

    //State changing Functions
    /**
     * @dev registers retired details
     * @param _retiredAddress is the address of the retired
     */
    function registerRetired(address _retiredAddress)
        public
        onlyOwner
        onlyAtStatus(Status.Initialized)
    {
        require(
            !retiredExists(_retiredAddress),
            "Retired address is already registered"
        );
        Retired storage retired = retirees.push();
        retired.retired = _retiredAddress;
        emit RetiredRegistered(retired.retired, retired.approvers);
    }

    /**
     * @dev remove the retired details from the array
     * @param _retiredAddress is the address of the retired
     */
    function removeRetired(address _retiredAddress)
        public
        onlyOwner
        onlyAtStatus(Status.Initialized)
    {
        uint256 retiredIndex = getRetiredIndex(_retiredAddress);
        retirees[retiredIndex] = retirees[retirees.length - 1];
        retirees.pop();

        emit RetiredRemoved(_retiredAddress, retirees.length);
    }

    /**
     * @dev registers newbie address and its fund distribution
     * @param _newbieAddress is the address of the newbie
     * @param _amount is the fund to be allocated to the newbie
     */
    function registerNewbie(address _newbieAddress, uint256 _amount)
        public
        onlyOwner
        onlyAtStatus(Status.Initialized)
    {
        require(
            !newbieExists(_newbieAddress),
            "Newbie address is already registered"
        );
        require(_amount != 0, "Amount cannot be set to 0");

        Newbie memory newbie = Newbie(_newbieAddress, _amount);
        newbies.push(newbie);

        emit NewbieRegistered(_newbieAddress, _amount);
    }

    /**
     * @dev remove the newbie details from the array
     * @param _newbieAddress is the address of the newbie
     */
    function removeNewbie(address _newbieAddress)
        public
        onlyOwner
        onlyAtStatus(Status.Initialized)
    {
        uint256 newbieIndex = getNewbieIndex(_newbieAddress);
        newbies[newbieIndex] = newbies[newbies.length - 1];
        newbies.pop();

        emit NewbieRemoved(_newbieAddress, newbies.length);
    }

    /**
     * @dev retiredAddress can be a EOA or a contract address. To approve:
     *      If the retiredAddress is a EOA, the msg.sender should be the EOA address
     *      If the retiredAddress is a Contract, the msg.sender should be one of the contract `admin`.
     *      It uses the getState() function in the retiredAddress contract to get the admin details.
     * @param _retiredAddress is the address of the retired
     */
    function approve(address _retiredAddress)
        public
        onlyAtStatus(Status.Registered)
    {
        require(
            retiredExists(_retiredAddress),
            "retired needs to be registered before approval"
        );

        //Check whether the retired address is EOA or contract address
        bool isContract = isContractAddr(_retiredAddress);
        if (!isContract) {
            //check whether the msg.sender is the retired if its a EOA
            require(
                msg.sender == _retiredAddress,
                "retiredAddress is not the msg.sender"
            );
            _updateApprover(_retiredAddress, msg.sender);
        } else {
            //check if the msg.sender is one of the admin of the retiredAddress contract
            require(
                _validateAdmin(_retiredAddress, msg.sender),
                "msg.sender is not the admin"
            );
            _updateApprover(_retiredAddress, msg.sender);
        }
    }

    /**
     * @dev validate if the msg.sender is admin if the retiredAddress is a contract
     * @param _retiredAddress is the address of the contract
     * @param _approver is the msg.sender
     * @return isAdmin is true if the msg.sender is one of the admin
     */
    function _validateAdmin(address _retiredAddress, address _approver)
        private
        returns (bool isAdmin)
    {
        (address[] memory adminList, ) = _getState(_retiredAddress);
        require(adminList.length != 0, "admin list cannot be empty");
        for (uint8 i = 0; i < adminList.length; i++) {
            if (_approver == adminList[i]) {
                isAdmin = true;
            }
        }
    }

    /**
     * @dev gets the adminList and quorom by calling `getState()` method in retiredAddress contract
     * @param _retiredAddress is the address of the contract
     * @return adminList list of the retiredAddress contract admins
     * @return req min required number of approvals
     */
    function _getState(address _retiredAddress)
        private
        returns (address[] memory adminList, uint256 req)
    {
        //call getState() function in retiredAddress contract to get the adminList
        bytes memory payload = abi.encodeWithSignature("getState()");
        (bool success, bytes memory result) = _retiredAddress.staticcall(
            payload
        );
        emit GetState(success, result);
        require(success, "call failed");

        (adminList, req) = abi.decode(result, (address[], uint256));
    }

    /**
     * @dev Internal function to update the approver details of a retired
     * _retiredAddress is the address of the retired
     * _approver is the admin of the retiredAddress
     */
    function _updateApprover(address _retiredAddress, address _approver)
        private
    {
        uint256 index = getRetiredIndex(_retiredAddress);
        address[] memory approvers = retirees[index].approvers;
        for (uint256 i = 0; i < approvers.length; i++) {
            require(approvers[i] != _approver, "Already approved");
        }
        retirees[index].approvers.push(_approver);
        emit Approved(
            _retiredAddress,
            _approver,
            retirees[index].approvers.length
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
        emit SetStatus(status);
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
            getTreasuryAmount() < sumOfRetiredBalance(),
            "treasury amount should be less than the sum of all retired address balances"
        );
        _isRetiredsApproved();
        status = Status.Approved;
        emit SetStatus(status);
    }

    /**
     * @dev verify if quorom reached for the retired approvals
     */
    function _isRetiredsApproved() private {
        for (uint256 i = 0; i < retirees.length; i++) {
            Retired memory retired = retirees[i];
            (, uint256 req) = _getState(retired.retired);
            if (retired.approvers.length >= req) {
                //if min quorom reached, make sure all approvers are still valid
                address[] memory approvers = retired.approvers;
                uint256 minApprovals = 0;
                for (uint256 j = 0; j < approvers.length; j++) {
                    _validateAdmin(retirees[i].retired, approvers[j]);
                    minApprovals++;
                }
                require(
                    minApprovals >= req,
                    "min required admins should approve"
                );
            } else {
                revert("min required admins should approve");
            }
        }
    }

    /**
     * @dev sets the status of the contract to Finalize. Once finalized the storage data
     * of the contract cannot be modified
     * @param _memo is the result of the rebalance after executing successfully in the core.
     */
    function finalizeContract(string memory _memo)
        public
        onlyOwner
        onlyAtStatus(Status.Approved)
    {
        memo = _memo;
        status = Status.Finalized;
        emit Finalized(memo, status);
    }

    /**
     * @dev resets all storage values to empty objects except targetBlockNumber
     */
    function reset() public onlyOwner {
        //reset cannot be called at Finalized status or after target block.number
        require(
            ((status != Status.Finalized) &&
                (block.number < rebalanceBlockNumber)),
            "Contract is finalized, cannot reset values"
        );

        //`delete` keyword is used to set a storage variable or a dynamic array to its default value.
        delete retirees;
        delete newbies;
        delete memo;
        status = Status.Initialized;
    }

    //Getters
    /**
     * @dev to get retired details by retiredAddress
     * @param _retiredAddress is the address of the retired
     */
    function getRetired(address _retiredAddress)
        public
        view
        returns (address, address[] memory)
    {
        require(retiredExists(_retiredAddress), "Retired does not exist");
        uint256 index = getRetiredIndex(_retiredAddress);
        Retired memory retired = retirees[index];
        return (retired.retired, retired.approvers);
    }

    /**
     * @dev check whether retiredAddress is registered
     * @param _retiredAddress is the address of the retired
     */
    function retiredExists(address _retiredAddress) public view returns (bool) {
        require(_retiredAddress != address(0), "Invalid address");
        for (uint8 i = 0; i < retirees.length; i++) {
            if (retirees[i].retired == _retiredAddress) {
                return true;
            }
        }
    }

    /**
     * @dev get index of the retired in the retirees array
     * @param _retiredAddress is the address of the retired
     */
    function getRetiredIndex(address _retiredAddress)
        public
        view
        returns (uint256)
    {
        for (uint256 i = 0; i < retirees.length; i++) {
            if (retirees[i].retired == _retiredAddress) {
                return i;
            }
        }
        revert("Retired does not exist");
    }

    /**
     * @dev to calculate the sum of retirees balances
     * @return retireesBalance the sum of balances of retireds
     */
    function sumOfRetiredBalance()
        public
        view
        returns (uint256 retireesBalance)
    {
        for (uint8 i = 0; i < retirees.length; i++) {
            address retiredAddress = retirees[i].retired;
            retireesBalance += retiredAddress.balance;
        }
        return retireesBalance;
    }

    /**
     * @dev to get newbie details by newbieAddress
     * @param _newbieAddress is the address of the newbie
     * @return newbie is the address of the newbie
     * @return amount is the fund allocated to the newbie
     
     */
    function getNewbie(address _newbieAddress)
        public
        view
        returns (address, uint256)
    {
        require(newbieExists(_newbieAddress), "Newbie does not exist");
        uint256 index = getNewbieIndex(_newbieAddress);
        Newbie memory newbie = newbies[index];
        return (newbie.newbie, newbie.amount);
    }

    /**
     * @dev check whether _newbieAddress is registered
     * @param _newbieAddress is the address of the newbie
     */
    function newbieExists(address _newbieAddress) public view returns (bool) {
        require(_newbieAddress != address(0), "Invalid address");
        for (uint8 i = 0; i < newbies.length; i++) {
            if (newbies[i].newbie == _newbieAddress) {
                return true;
            }
        }
    }

    /**
     * @dev get index of the newbie in the newbies array
     * @param _newbieAddress is the address of the newbie
     */
    function getNewbieIndex(address _newbieAddress)
        public
        view
        returns (uint256)
    {
        for (uint256 i = 0; i < newbies.length; i++) {
            if (newbies[i].newbie == _newbieAddress) {
                return i;
            }
        }
        revert("Newbie does not exist");
    }

    /**
     * @dev to calculate the sum of newbie funds
     * @return treasuryAmount the sum of funds allocated to newbies
     */
    function getTreasuryAmount() public view returns (uint256 treasuryAmount) {
        for (uint8 i = 0; i < newbies.length; i++) {
            treasuryAmount += newbies[i].amount;
        }
        return treasuryAmount;
    }

    /**
     * @dev gets the length of retirees list
     */
    function getRetiredCount() public view returns (uint256) {
        return retirees.length;
    }

    /**
     * @dev gets the length of newbies list
     */
    function getNewbieCount() public view returns (uint256) {
        return newbies.length;
    }

    /**
     * @dev allback function to revert any payments
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
