// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

import "./Ownable.sol";
import "./ITreasuryRebalance.sol";

/**
 * @title Interface to get adminlist and quorom
 */
interface IRetiredContract {
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
contract TreasuryRebalance is Ownable, ITreasuryRebalance {
    /**
     * Storage
     */
    Retired[] public retirees; // array of the Retired struct
    Newbie[] public newbies; // array of Newbie struct
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
        rebalanceBlockNumber = _rebalanceBlockNumber;
        status = Status.Initialized;
        emit ContractDeployed(status, _rebalanceBlockNumber, block.timestamp);
    }

    //State changing Functions
    /**
     * @dev registers retired details
     * @param _retiredAddress is the address of the retired
     */
    function registerRetired(
        address _retiredAddress
    ) public onlyOwner onlyAtStatus(Status.Initialized) {
        require(
            !retiredExists(_retiredAddress),
            "Retired address is already registered"
        );
        Retired storage retired = retirees.push();
        retired.retired = _retiredAddress;
        emit RetiredRegistered(retired.retired);
    }

    /**
     * @dev remove the retired details from the array
     * @param _retiredAddress is the address of the retired
     */
    function removeRetired(
        address _retiredAddress
    ) public onlyOwner onlyAtStatus(Status.Initialized) {
        uint256 retiredIndex = getRetiredIndex(_retiredAddress);
        require(retiredIndex != type(uint256).max, "Retired not registered");
        retirees[retiredIndex] = retirees[retirees.length - 1];
        retirees.pop();

        emit RetiredRemoved(_retiredAddress);
    }

    /**
     * @dev registers newbie address and its fund distribution
     * @param _newbieAddress is the address of the newbie
     * @param _amount is the fund to be allocated to the newbie
     */
    function registerNewbie(
        address _newbieAddress,
        uint256 _amount
    ) public onlyOwner onlyAtStatus(Status.Initialized) {
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
    function removeNewbie(
        address _newbieAddress
    ) public onlyOwner onlyAtStatus(Status.Initialized) {
        uint256 newbieIndex = getNewbieIndex(_newbieAddress);
        require(newbieIndex != type(uint256).max, "Newbie not registered");
        newbies[newbieIndex] = newbies[newbies.length - 1];
        newbies.pop();

        emit NewbieRemoved(_newbieAddress);
    }

    /**
     * @dev retiredAddress can be a EOA or a contract address. To approve:
     *      If the retiredAddress is a EOA, the msg.sender should be the EOA address
     *      If the retiredAddress is a Contract, the msg.sender should be one of the contract `admin`.
     *      It uses the getState() function in the retiredAddress contract to get the admin details.
     * @param _retiredAddress is the address of the retired
     */
    function approve(
        address _retiredAddress
    ) public onlyAtStatus(Status.Registered) {
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
            (address[] memory adminList, ) = _getState(_retiredAddress);
            require(adminList.length != 0, "admin list cannot be empty");

            //check if the msg.sender is one of the admin of the retiredAddress contract
            require(
                _validateAdmin(msg.sender, adminList),
                "msg.sender is not the admin"
            );
            _updateApprover(_retiredAddress, msg.sender);
        }
    }

    /**
     * @dev validate if the msg.sender is admin if the retiredAddress is a contract
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
     * @dev gets the adminList and quorom by calling `getState()` method in retiredAddress contract
     * @param _retiredAddress is the address of the contract
     * @return adminList list of the retiredAddress contract admins
     * @return req min required number of approvals
     */
    function _getState(
        address _retiredAddress
    ) private view returns (address[] memory adminList, uint256 req) {
        IRetiredContract retiredContract = IRetiredContract(_retiredAddress);
        (adminList, req) = retiredContract.getState();
    }

    /**
     * @dev Internal function to update the approver details of a retired
     * _retiredAddress is the address of the retired
     * _approver is the admin of the retiredAddress
     */
    function _updateApprover(
        address _retiredAddress,
        address _approver
    ) private {
        uint256 index = getRetiredIndex(_retiredAddress);
        require(index != type(uint256).max, "Retired not registered");
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
            getTreasuryAmount() < sumOfRetiredBalance(),
            "treasury amount should be less than the sum of all retired address balances"
        );
        checkRetiredsApproved();
        status = Status.Approved;
        emit StatusChanged(status);
    }

    /**
     * @dev verify if quorom reached for the retired approvals
     */
    function checkRetiredsApproved() public view {
        for (uint256 i = 0; i < retirees.length; i++) {
            Retired memory retired = retirees[i];
            bool isContract = isContractAddr(retired.retired);
            if (isContract) {
                (address[] memory adminList, uint256 req) = _getState(
                    retired.retired
                );
                require(
                    retired.approvers.length >= req,
                    "min required admins should approve"
                );
                //if min quorom reached, make sure all approvers are still valid
                address[] memory approvers = retired.approvers;
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
                require(retired.approvers.length == 1, "EOA should approve");
            }
        }
    }

    /**
     * @dev sets the status of the contract to Finalize. Once finalized the storage data
     * of the contract cannot be modified
     * @param _memo is the result of the rebalance after executing successfully in the core.
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
    function getRetired(
        address _retiredAddress
    ) public view returns (address, address[] memory) {
        uint256 index = getRetiredIndex(_retiredAddress);
        require(index != type(uint256).max, "Retired not registered");
        Retired memory retired = retirees[index];
        return (retired.retired, retired.approvers);
    }

    /**
     * @dev check whether retiredAddress is registered
     * @param _retiredAddress is the address of the retired
     */
    function retiredExists(address _retiredAddress) public view returns (bool) {
        require(_retiredAddress != address(0), "Invalid address");
        for (uint256 i = 0; i < retirees.length; i++) {
            if (retirees[i].retired == _retiredAddress) {
                return true;
            }
        }
    }

    /**
     * @dev get index of the retired in the retirees array
     * @param _retiredAddress is the address of the retired
     */
    function getRetiredIndex(
        address _retiredAddress
    ) public view returns (uint256) {
        for (uint256 i = 0; i < retirees.length; i++) {
            if (retirees[i].retired == _retiredAddress) {
                return i;
            }
        }
        return type(uint256).max;
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
        for (uint256 i = 0; i < retirees.length; i++) {
            retireesBalance += retirees[i].retired.balance;
        }
        return retireesBalance;
    }

    /**
     * @dev to get newbie details by newbieAddress
     * @param _newbieAddress is the address of the newbie
     * @return newbie is the address of the newbie
     * @return amount is the fund allocated to the newbie
     */
    function getNewbie(
        address _newbieAddress
    ) public view returns (address, uint256) {
        uint256 index = getNewbieIndex(_newbieAddress);
        require(index != type(uint256).max, "Newbie not registered");
        Newbie memory newbie = newbies[index];
        return (newbie.newbie, newbie.amount);
    }

    /**
     * @dev check whether _newbieAddress is registered
     * @param _newbieAddress is the address of the newbie
     */
    function newbieExists(address _newbieAddress) public view returns (bool) {
        require(_newbieAddress != address(0), "Invalid address");
        for (uint256 i = 0; i < newbies.length; i++) {
            if (newbies[i].newbie == _newbieAddress) {
                return true;
            }
        }
    }

    /**
     * @dev get index of the newbie in the newbies array
     * @param _newbieAddress is the address of the newbie
     */
    function getNewbieIndex(
        address _newbieAddress
    ) public view returns (uint256) {
        for (uint256 i = 0; i < newbies.length; i++) {
            if (newbies[i].newbie == _newbieAddress) {
                return i;
            }
        }
        return type(uint256).max;
    }

    /**
     * @dev to calculate the sum of newbie funds
     * @return treasuryAmount the sum of funds allocated to newbies
     */
    function getTreasuryAmount() public view returns (uint256 treasuryAmount) {
        for (uint256 i = 0; i < newbies.length; i++) {
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
