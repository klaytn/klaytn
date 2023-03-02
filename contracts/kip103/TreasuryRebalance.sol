// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

import "./Ownable.sol";

/**
 * @title Smart contract to record the rebalance of treasury funds.
 * This contract is to mainly record the addresses which holds the treasury funds
 * before and after rebalancing. It facilates approval and redistributing to new addresses.
 * Core will execute the re-distribution reading from this contract.
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
     * Sender struct to store the details of sender and approver addresses
     */
    struct Sender {
        address sender;
        address[] approvers;
    }

    /**
     * Receiver struct to store reciever and amount
     */
    struct Receiver {
        address receiver;
        uint256 amount;
    }

    /**
     * Storage
     */
    Sender[] public senders; //array of sender structs
    Receiver[] public receivers; //array of receiver structs
    Status public status; //current status of the contract
    uint256 public rebalanceBlockNumber; //Block number of the execution of rebalancing
    string public memo; //result of the treasury fund rebalance

    /**
     * Events logs
     */
    event DeployContract(
        Status status,
        uint256 rebalanceBlockNumber,
        uint256 deployedBlockNumber
    );
    event RegisterSender(address sender, address[] approvers);
    event RemoveSender(address sender, uint256 senderCount);
    event RegisterReceiver(address receiver, uint256 fundAllocation);
    event RemoveReceiver(address receiver, uint256 receiverCount);
    event GetState(bool success, bytes result);
    event Approve(address sender, address approver, uint256 approversCount);
    event SetStatus(Status status);
    event Finalized(string memo, Status status);

    /**
     * Modifiers
     */
    modifier atStatus(Status _status) {
        require(status == _status, "function not allowed at this stage");
        _;
    }

    /**
     *  Constructor
     * @param _rebalanceBlockNumber is the target block number to execute the redistribution in Core
     */
    constructor(uint256 _rebalanceBlockNumber) {
        rebalanceBlockNumber = _rebalanceBlockNumber;
        status = Status.Initialized;
        emit DeployContract(status, _rebalanceBlockNumber, block.timestamp);
    }

    //State changing Functions
    /**
     * @dev registers sender details
     * @param _senderAddress is the address of the sender
     */
    function registerSender(address _senderAddress)
        public
        onlyOwner
        atStatus(Status.Initialized)
    {
        require(!senderExists(_senderAddress), "Sender is already registered");
        Sender storage sender = senders.push();
        sender.sender = _senderAddress;
        emit RegisterSender(sender.sender, sender.approvers);
    }

    /**
     * @dev remove the sender details from the array
     * @param _senderAddress is the address of the sender
     */
    function removeSender(address _senderAddress)
        public
        onlyOwner
        atStatus(Status.Initialized)
    {
        uint256 senderIndex = getSenderIndex(_senderAddress);
        senders[senderIndex] = senders[senders.length - 1];
        senders.pop();

        emit RemoveSender(_senderAddress, senders.length);
    }

    /**
     * @dev registers receiver address and its fund distribution
     * @param _receiverAddress is the address of the receiver
     * @param _amount is the fund to be allocated to the receiver
     */
    function registerReceiver(address _receiverAddress, uint256 _amount)
        public
        onlyOwner
        atStatus(Status.Initialized)
    {
        require(
            !receiverExists(_receiverAddress),
            "Receiver is already registered"
        );
        require(_amount != 0, "Amount cannot be set to 0");

        Receiver memory receiver = Receiver(_receiverAddress, _amount);
        receivers.push(receiver);

        emit RegisterReceiver(_receiverAddress, _amount);
    }

    /**
     * @dev remove the receiver details from the array
     * @param _receiverAddress is the address of the receiver
     */
    function removeReceiver(address _receiverAddress)
        public
        onlyOwner
        atStatus(Status.Initialized)
    {
        uint256 receiverIndex = getReceiverIndex(_receiverAddress);
        receivers[receiverIndex] = receivers[receivers.length - 1];
        receivers.pop();

        emit RemoveReceiver(_receiverAddress, receivers.length);
    }

    /**
     * @dev senderAddress can be a EOA or a contract address. To approve:
     *      If the senderAddress is a EOA, the msg.sender should be the EOA address
     *      If the senderAddress is a Contract, the msg.sender should be one of the contract `admin`.
     *      It uses the getState() function in the senderAddress contract to get the admin details.
     * @param _senderAddress is the address of the sender
     */
    function approve(address _senderAddress)
        public
        atStatus(Status.Registered)
    {
        require(
            senderExists(_senderAddress),
            "sender needs to be registered before approval"
        );

        //Check whether the sender address is EOA or contract address
        bool isContract = isContractAddr(_senderAddress);
        if (!isContract) {
            //check whether the msg.sender is the sender if its a EOA
            require(
                msg.sender == _senderAddress,
                "senderAddress is not the msg.sender"
            );
            _updateApprover(_senderAddress, msg.sender);
        } else {
            //check if the msg.sender is one of the admin of the senderAddress contract
            require(
                _validateAdmin(_senderAddress, msg.sender),
                "msg.sender is not the admin"
            );
            _updateApprover(_senderAddress, msg.sender);
        }
    }

    /**
     * @dev validate if the msg.sender is admin if the senderAddress is a contract
     * @param _senderAddress is the address of the contract
     * @param _approver is the msg.sender
     * @return isAdmin is true if the msg.sender is one of the admin
     */
    function _validateAdmin(address _senderAddress, address _approver)
        private
        returns (bool isAdmin)
    {
        (address[] memory adminList, ) = _getState(_senderAddress);
        require(adminList.length != 0, "admin list cannot be empty");
        for (uint8 i = 0; i < adminList.length; i++) {
            if (_approver == adminList[i]) {
                isAdmin = true;
            }
        }
    }

    /**
     * @dev gets the adminList and quorom by calling `getState()` method in senderAddress contract
     * @param _senderAddress is the address of the contract
     * @return adminList list of the senderAddress contract admins
     * @return req min required number of approvals
     */
    function _getState(address _senderAddress)
        private
        returns (address[] memory adminList, uint256 req)
    {
        //call getState() function in senderAddress contract to get the adminList
        bytes memory payload = abi.encodeWithSignature("getState()");
        (bool success, bytes memory result) = _senderAddress.staticcall(
            payload
        );
        emit GetState(success, result);
        require(success, "call failed");

        (adminList, req) = abi.decode(result, (address[], uint256));
    }

    /**
     * @dev Internal function to update the approver details of a sender
     * _senderAddress is the address of the sender
     * _approver is the admin of the senderAddress
     */
    function _updateApprover(address _senderAddress, address _approver)
        private
    {
        uint256 index = getSenderIndex(_senderAddress);
        address[] memory approvers = senders[index].approvers;
        for (uint256 i = 0; i < approvers.length; i++) {
            require(
                approvers[i] != _approver,
                "Duplicate approvers cannot be allowed"
            );
        }
        senders[index].approvers.push(_approver);
        emit Approve(
            _senderAddress,
            _approver,
            senders[index].approvers.length
        );
    }

    /**
     * @dev finalizeRegistration sets the status to Registered,
     *      After this stage, registrations will be restricted.
     */
    function finalizeRegistration()
        public
        onlyOwner
        atStatus(Status.Initialized)
    {
        status = Status.Registered;
        emit SetStatus(status);
    }

    /**
     * @dev finalizeApproval sets the status to Approved,
     *      After this stage, approvals will be restricted.
     */
    function finalizeApproval() public onlyOwner atStatus(Status.Registered) {
        require(
            getTreasuryAmount() < sumOfSenderBalance(),
            "treasury amount should be less than the sum of all sender address balances"
        );
        _isSendersApproved();
        status = Status.Approved;
        emit SetStatus(status);
    }

    /**
     * @dev verify if quorom reached for the sender approvals
     */
    function _isSendersApproved() private {
        for (uint256 i = 0; i < senders.length; i++) {
            Sender memory sender = senders[i];
            (, uint256 req) = _getState(sender.sender);
            if (sender.approvers.length >= req) {
                //if min quorom reached, make sure all approvers are still valid
                address[] memory approvers = sender.approvers;
                uint256 minApprovals = 0;
                for (uint256 j = 0; j < approvers.length; j++) {
                    _validateAdmin(senders[i].sender, approvers[j]);
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
        atStatus(Status.Approved)
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
        delete senders;
        delete receivers;
        delete memo;
        status = Status.Initialized;
    }

    //Getters
    /**
     * @dev to get sender details by senderAddress
     * @param _senderAddress is the address of the sender
     */
    function getSender(address _senderAddress)
        public
        view
        returns (address, address[] memory)
    {
        require(senderExists(_senderAddress), "Sender does not exist");
        uint256 index = getSenderIndex(_senderAddress);
        Sender memory sender = senders[index];
        return (sender.sender, sender.approvers);
    }

    /**
     * @dev check whether senderAddress is registered
     * @param _senderAddress is the address of the sender
     */
    function senderExists(address _senderAddress) public view returns (bool) {
        require(_senderAddress != address(0), "Invalid address");
        for (uint8 i = 0; i < senders.length; i++) {
            if (senders[i].sender == _senderAddress) {
                return true;
            }
        }
    }

    /**
     * @dev get index of the sender in the senders array
     * @param _senderAddress is the address of the sender
     */
    function getSenderIndex(address _senderAddress)
        public
        view
        returns (uint256)
    {
        for (uint256 i = 0; i < senders.length; i++) {
            if (senders[i].sender == _senderAddress) {
                return i;
            }
        }
        revert("Sender does not exist");
    }

    /**
     * @dev to calculate the sum of senders balances
     * @return sendersBalance the sum of balances of senders
     */
    function sumOfSenderBalance() public view returns (uint256 sendersBalance) {
        for (uint8 i = 0; i < senders.length; i++) {
            address senderAddress = senders[i].sender;
            sendersBalance += senderAddress.balance;
        }
        return sendersBalance;
    }

    /**
     * @dev to get receiver details by receiverAddress
     * @param _receiverAddress is the address of the receiver
     * @return receiver is the address of the receiver
     * @return amount is the fund allocated to the receiver
     
     */
    function getReceiver(address _receiverAddress)
        public
        view
        returns (address, uint256)
    {
        require(receiverExists(_receiverAddress), "Receiver does not exist");
        uint256 index = getReceiverIndex(_receiverAddress);
        Receiver memory receiver = receivers[index];
        return (receiver.receiver, receiver.amount);
    }

    /**
     * @dev check whether _receiverAddress is registered
     * @param _receiverAddress is the address of the receiver
     */
    function receiverExists(address _receiverAddress)
        public
        view
        returns (bool)
    {
        require(_receiverAddress != address(0), "Invalid address");
        for (uint8 i = 0; i < receivers.length; i++) {
            if (receivers[i].receiver == _receiverAddress) {
                return true;
            }
        }
    }

    /**
     * @dev get index of the receiver in the receivers array
     * @param _receiverAddress is the address of the receiver
     */
    function getReceiverIndex(address _receiverAddress)
        public
        view
        returns (uint256)
    {
        for (uint256 i = 0; i < receivers.length; i++) {
            if (receivers[i].receiver == _receiverAddress) {
                return i;
            }
        }
        revert("Receiver does not exist");
    }

    /**
     * @dev to calculate the sum of receiver funds
     * @return treasuryAmount the sum of funds allocated to receivers
     */
    function getTreasuryAmount() public view returns (uint256 treasuryAmount) {
        for (uint8 i = 0; i < receivers.length; i++) {
            treasuryAmount += receivers[i].amount;
        }
        return treasuryAmount;
    }

    /**
     * @dev gets the length of senders list
     */
    function getSenderCount() public view returns (uint256) {
        return senders.length;
    }

    /**
     * @dev gets the length of receivers list
     */
    function getReceiverCount() public view returns (uint256) {
        return receivers.length;
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
