// Copyright 2019 The klaytn Authors
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

pragma solidity 0.4.24;
import "./SafeMath.sol";

contract CnStakingContract {
    using SafeMath for uint256;
    /*
     *  Events
     */
    event DeployContract(string contractType, address contractValidator, address nodeId, address rewardAddress, address[] cnAdminList, uint256 requirement, uint256[] unlockTime, uint256[] unlockAmount);
    event ReviewInitialConditions(address indexed from);
    event CompleteReviewInitialConditions();
    event DepositLockupStakingAndInit(address from, uint256 value);

    event SubmitRequest(uint256 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg);
    event ConfirmRequest(uint256 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers);
    event RevokeConfirmation(uint256 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers);
    event CancelRequest(uint256 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg);
    event ExecuteRequestSuccess(uint256 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg);
    event ExecuteRequestFailure(uint256 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg);
    event ClearRequest();

    event AddAdmin (address indexed admin);
    event DeleteAdmin(address indexed admin);
    event UpdateRequirement(uint256 requirement);
    event WithdrawLockupStaking(address indexed to, uint256 value);
    event ApproveStakingWithdrawal(uint256 approvedWithdrawalId, address to, uint256 value, uint256 withdrawableFrom);
    event CancelApprovedStakingWithdrawal(uint256 approvedWithdrawalId, address to, uint256 value);
    event UpdateRewardAddress(address rewardAddress);

    event StakeKlay(address from, uint256 value);
    event WithdrawApprovedStaking(uint256 approvedWithdrawalId, address to, uint256 value);

    //AddressBook event
    event ReviseRewardAddress(address cnNodeId, address prevRewardAddress, address curRewardAddress);


    /*
     *  Constants
     */
    uint256 constant public MAX_ADMIN = 50;
    string constant public CONTRACT_TYPE = "CnStakingContract";
    uint256 constant public VERSION = 1;
    address constant public ADDRESS_BOOK_ADDRESS = 0x0000000000000000000000000000000000000400;
    uint256 constant public ONE_WEEK = 1 weeks;


    /*
     *  Enums
     */
    enum RequestState { Unknown, NotConfirmed, Executed, ExecutionFailed, Canceled }
    enum Functions { Unknown, AddAdmin, DeleteAdmin, UpdateRequirement, ClearRequest, WithdrawLockupStaking, ApproveStakingWithdrawal, CancelApprovedStakingWithdrawal, UpdateRewardAddress }
    enum WithdrawalStakingState { Unknown, Transferred, Canceled}

    /*
     *  Storage
     */
    address[] private adminList;
    uint256 public requirement;
    mapping (address => bool) private isAdmin;
    uint256 public lastClearedId;

    uint256 public requestCount;
    mapping(uint256 => Request) private requestMap;
    struct Request {
        Functions functionId;
        bytes32 firstArg;
        bytes32 secondArg;
        bytes32 thirdArg;
        address requestProposer;
        address[] confirmers;
        RequestState state;
    }

    address public contractValidator;
    bool public isInitialized;
    struct LockupConditions {
        uint256[] unlockTime;
        uint256[] unlockAmount;
        bool allReviewed;
        uint256 reviewedCount;
        mapping(address => bool) reviewedAdmin;
    }
    LockupConditions public lockupConditions;
    uint256 public initialLockupStaking;
    uint256 public remainingLockupStaking;
    address public nodeId;
    address public rewardAddress;


    uint256 public staking;
    uint256 public withdrawalRequestCount;
    mapping(uint256 => WithdrawalRequest) private withdrawalRequestMap;
    struct WithdrawalRequest {
        address to;
        uint256 value;
        uint256 withdrawableFrom;
        WithdrawalStakingState state;
    }


    /*
     *  Modifiers
     */
    modifier onlyMultisigTx() {
        require(msg.sender == address(this), "Not a multisig-transaction.");
        _;
    }

    modifier onlyAdmin(address _admin) {
        require(isAdmin[_admin], "Address is not admin.");
        _;
    }

    modifier adminDoesNotExist(address _admin) {
        require(!isAdmin[_admin], "Admin already exists.");
        _;
    }

    modifier notNull(address _address) {
        require(_address != 0, "Address is null");
        _;
    }

    modifier notConfirmedRequest(uint256 _id) {
        require(requestMap[_id].state == RequestState.NotConfirmed, "Must be at not-confirmed state.");
        _;
    }

    modifier validRequirement(uint256 _adminCount, uint256 _requirement) {
        require(_adminCount <= MAX_ADMIN
            && _requirement <= _adminCount
            && _requirement != 0
            && _adminCount != 0, "Invalid requirement.");
        _;
    }

    modifier beforeInit() {
        require(isInitialized == false, "Contract has been initialized.");
        _;
    }

    modifier afterInit() {
        require(isInitialized == true, "Contract is not initialized.");
        _;
    }


    /*
     *  Constructor
     */

    constructor(address _contractValidator, address _nodeId, address _rewardAddress, address[] _cnAdminlist, uint256 _requirement, uint256[] _unlockTime, uint256[] _unlockAmount) public
    validRequirement(_cnAdminlist.length, _requirement) 
    notNull(_nodeId) 
    notNull(_rewardAddress) {

        require(_contractValidator != 0, "Validator is null.");
        isAdmin[_contractValidator] = true;
        for(uint256 i = 0; i < _cnAdminlist.length; i++) {
            require(!isAdmin[_cnAdminlist[i]] && _cnAdminlist[i] != 0, "Address is null or not unique.");
            isAdmin[_cnAdminlist[i]] = true;
        }


        require(_unlockTime.length != 0 && _unlockAmount.length != 0 && _unlockTime.length == _unlockAmount.length, "Invalid unlock time and amount.");
        uint256 unlockTime = now;

        for (i = 0; i < _unlockAmount.length; i++) {
            require(unlockTime < _unlockTime[i], "Unlock time is not in ascending order.");
            require(_unlockAmount[i] > 0, "Amount is not positive number.");
            unlockTime = _unlockTime[i];
        }

        contractValidator = _contractValidator;
        nodeId = _nodeId;
        rewardAddress = _rewardAddress;
        adminList = _cnAdminlist;
        requirement = _requirement;
        lockupConditions.unlockTime = _unlockTime;
        lockupConditions.unlockAmount = _unlockAmount;

        isInitialized = false;
        emit DeployContract(CONTRACT_TYPE, _contractValidator, _nodeId, _rewardAddress, _cnAdminlist, _requirement, _unlockTime, _unlockAmount);
    }



    function reviewInitialConditions() external
    onlyAdmin(msg.sender)
    beforeInit() {
        require(lockupConditions.reviewedAdmin[msg.sender] == false, "Msg.sender already reviewed.");

        lockupConditions.reviewedAdmin[msg.sender] = true;
        lockupConditions.reviewedCount = lockupConditions.reviewedCount.add(1);
        
        emit ReviewInitialConditions(msg.sender);


        if(lockupConditions.reviewedCount == adminList.length + 1) {
            lockupConditions.allReviewed = true;
            emit CompleteReviewInitialConditions();
        }
    }


    function depositLockupStakingAndInit() external payable
    beforeInit() {
        uint256 requiredStakingAmount;
        uint256 cnt = lockupConditions.unlockAmount.length;
        for(uint256 i = 0; i < cnt; i++) {
            requiredStakingAmount = requiredStakingAmount.add(lockupConditions.unlockAmount[i]);
        }
        require(lockupConditions.allReviewed == true, "Reviewing is not finished.");
        require(msg.value == requiredStakingAmount, "Value does not match.");

        isAdmin[contractValidator] = false;
        delete contractValidator;

        initialLockupStaking = requiredStakingAmount;
        remainingLockupStaking = requiredStakingAmount;
        isInitialized = true;
        emit DepositLockupStakingAndInit(msg.sender, msg.value);
    }


    /*
     *  Submit multisig function request
     */

    /// @notice submit a request to add new admin at consensus node (multi-sig operation)
    /// @param _admin new admin address to be added
    function submitAddAdmin(address _admin) external
    afterInit()
    adminDoesNotExist(_admin)
    notNull(_admin)
    onlyAdmin(msg.sender)
    validRequirement(adminList.length.add(1), requirement) {
        uint256 id = requestCount;
        submitRequest(id, Functions.AddAdmin, bytes32(_admin), 0, 0);
        confirmRequest(id, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
    }

    /// @notice submit a request to delete an admin at consensus node (multi-sig operation)
    /// @param _admin address of the admin to be deleted
    function submitDeleteAdmin(address _admin) external
    afterInit()
    onlyAdmin(_admin)
    notNull(_admin)
    onlyAdmin(msg.sender)
    validRequirement(adminList.length.sub(1), requirement) {
        uint256 id = requestCount;
        submitRequest(id, Functions.DeleteAdmin, bytes32(_admin), 0, 0);
        confirmRequest(id, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
    }

    /// @notice submit a request to update the confirmation threshold (multi-sig operation)
    /// @param _requirement new confirmation threshold
    function submitUpdateRequirement(uint256 _requirement) external 
    afterInit()
    onlyAdmin(msg.sender)
    validRequirement(adminList.length, _requirement) {
        require(_requirement != requirement, "Invalid value");
        uint256 id = requestCount;
        submitRequest(id, Functions.UpdateRequirement, bytes32(_requirement), 0, 0);
        confirmRequest(id, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
    }

    /// @notice submit a request to clear all unfinalized request (multi-sig operation)
    function submitClearRequest() external 
    afterInit()
    onlyAdmin(msg.sender) {
        uint256 id = requestCount;
        submitRequest(id, Functions.ClearRequest, 0, 0, 0);
        confirmRequest(id, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
    }

    function submitWithdrawLockupStaking(address _to, uint256 _value) external
    afterInit()
    notNull(_to)
    onlyAdmin(msg.sender) {
        uint256 withdrawableStakingAmount;
        ( , , , ,withdrawableStakingAmount) = getLockupStakingInfo();
        require(_value > 0 && _value <= withdrawableStakingAmount, "Invalid value.");

        uint256 id = requestCount;
        submitRequest(id, Functions.WithdrawLockupStaking, bytes32(_to), bytes32(_value), 0);
        confirmRequest(id, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
    }

    /// @notice submit a request to withdraw staked KLAY (multi-sig operation).
    ///         7 days after "approveStakingWithdrawal" has been requested,
    ///         admins can call this function for actual withdrawal for another 7-day period.
    ///         If the admin doesn't withraw KLAY for that 7-day period, it expires.
    /// @param _to target address to receive KLAY
    /// @param _value withdrawl amount of KLAY
    function submitApproveStakingWithdrawal(address _to, uint256 _value) external
    afterInit()
    notNull(_to)
    onlyAdmin(msg.sender) {
        require(_value > 0 && _value <= staking, "Invalid value.");
        uint256 id = requestCount;
        submitRequest(id, Functions.ApproveStakingWithdrawal, bytes32(_to), bytes32(_value), 0);
        confirmRequest(id, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
    }

    /// @notice submit a request to cancel KLAY withdrawl request (multi-sig operation).
    /// @param _approvedWithdrawalId the withdrawal ID to cancel. The ID is acquired at the event log of ApproveStakingWithdrawal
    function submitCancelApprovedStakingWithdrawal(uint256 _approvedWithdrawalId) external
    afterInit()
    onlyAdmin(msg.sender) {
        require(withdrawalRequestMap[_approvedWithdrawalId].to != 0, "Withdrawal request does not exist.");
        require(withdrawalRequestMap[_approvedWithdrawalId].state == WithdrawalStakingState.Unknown, "Invalid state.");

        uint256 id = requestCount;
        submitRequest(id, Functions.CancelApprovedStakingWithdrawal, bytes32(_approvedWithdrawalId), 0, 0);
        confirmRequest(id, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
    }

    /// @notice submit a request to update the reward address of consensus node (multi-sig operation).
    /// @param _rewardAddress new reward address
    function submitUpdateRewardAddress(address _rewardAddress) external 
    afterInit()
    notNull(_rewardAddress) 
    onlyAdmin(msg.sender) {
        uint256 id = requestCount;
        submitRequest(id, Functions.UpdateRewardAddress, bytes32(_rewardAddress), 0, 0);
        confirmRequest(id, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
    }


    /*
     *  Confirm and revoke confirmation
     */
    /// @notice after requesting of proposer, other admin confirm the request by confirmRequest
    /// @param _id the request ID. It can be obtained by submitRequest event
    /// @param _functionId the function ID of the request. It can be obtained by submitRequest event or getRequestInfo getter function
    /// @param _firstArg the first argument of the request. It can be obtained by submitRequest event or getRequestInfo getter function
    /// @param _secondArg the second argument of the request. If there is no second argument, then it should be 0. It can be obtained by submitRequest event or getRequestInfo getter function
    /// @param _thirdArg the third argument of the request. If there is no second argument, then it should be 0. It can be obtained by submitRequest event or getRequestInfo getter function
    function confirmRequest(uint256 _id, Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) public
    notConfirmedRequest(_id) 
    onlyAdmin(msg.sender) {
        bool hasConfirmed = false;
        uint256 confirmersCnt = requestMap[_id].confirmers.length;
        for (uint256 i = 0; i < confirmersCnt; i++) {
            if(msg.sender == requestMap[_id].confirmers[i]) {
                hasConfirmed = true;
                break;
            }
        }
        require(!hasConfirmed, "Msg.sender already confirmed.");
        require(
            requestMap[_id].functionId == _functionId &&
            requestMap[_id].firstArg == _firstArg &&
            requestMap[_id].secondArg == _secondArg &&
            requestMap[_id].thirdArg == _thirdArg, "Function id and arguments do not match.");

        requestMap[_id].confirmers.push(msg.sender);

        address[] memory confirmers = requestMap[_id].confirmers;
        emit ConfirmRequest(_id, msg.sender, _functionId, _firstArg, _secondArg, _thirdArg, confirmers);


        if (checkQuorum(_id)) {
            executeRequest(_id);
        }
    }

    /// @notice revoking confirmation of each admin. If the admin is proposer, then the request will be cancled regardless of confirmation of other admins.
    /// @param _id the request ID. It can be obtained by submitRequest event
    /// @param _functionId the function ID of the request. It can be obtained by submitRequest event or getRequestInfo getter function
    /// @param _firstArg the first argument of the request. It can be obtained by submitRequest event or getRequestInfo getter function
    /// @param _secondArg the second argument of the request. If there is no second argument, then it should be 0. It can be obtained by submitRequest event or getRequestInfo getter function
    /// @param _thirdArg the third argument of the request. If there is no second argument, then it should be 0. It can be obtained by submitRequest event or getRequestInfo getter function
    function revokeConfirmation(uint256 _id, Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) external
    notConfirmedRequest(_id)
    onlyAdmin(msg.sender) {
        bool hasConfirmed = false;
        uint256 confirmersCnt = requestMap[_id].confirmers.length;
        for (uint256 i = 0; i < confirmersCnt; i++) {
            if(msg.sender == requestMap[_id].confirmers[i]) {
                hasConfirmed = true;
                break;
            }
        }
        require(hasConfirmed, "Msg.sender has not confirmed.");
        require(
            requestMap[_id].functionId == _functionId &&
            requestMap[_id].firstArg == _firstArg &&
            requestMap[_id].secondArg == _secondArg &&
            requestMap[_id].thirdArg == _thirdArg, "Function id and arguments do not match.");

        revokeHandler(_id);
    }

    // this function separated from revokeConfirmation function because of 'stack too deep' issue of solidity
    function revokeHandler(uint256 _id) private {

        if (requestMap[_id].requestProposer == msg.sender) {
            requestMap[_id].state = RequestState.Canceled;
            emit CancelRequest(_id, msg.sender, requestMap[_id].functionId, requestMap[_id].firstArg, requestMap[_id].secondArg, requestMap[_id].thirdArg);
        } else {
            deleteFromConfirmerList(_id, msg.sender);
            emit RevokeConfirmation(_id, msg.sender, requestMap[_id].functionId, requestMap[_id].firstArg, requestMap[_id].secondArg, requestMap[_id].thirdArg, requestMap[_id].confirmers);
        }
    }


    /*
     *  Multisig functions
     */
    function addAdmin(address _admin) external 
    onlyMultisigTx()
    adminDoesNotExist(_admin)
    validRequirement(adminList.length.add(1), requirement) {
        isAdmin[_admin] = true;
        adminList.push(_admin);
        clearRequest();
        emit AddAdmin(_admin);
    }

    function deleteAdmin(address _admin) external
    onlyMultisigTx()
    onlyAdmin(_admin)
    validRequirement(adminList.length.sub(1), requirement) {
        isAdmin[_admin] = false;
        
        uint256 adminCnt = adminList.length;
        for (uint256 i = 0; i < adminCnt - 1; i++) {
            if (adminList[i] == _admin) {
                adminList[i] = adminList[adminCnt - 1];
                break;
            }
        }
        delete adminList[adminCnt - 1];
        adminList.length = adminList.length.sub(1);
        clearRequest();
        emit DeleteAdmin(_admin);
    }

    function updateRequirement(uint256 _requirement) external 
    onlyMultisigTx()
    validRequirement(adminList.length, _requirement) {
        requirement = _requirement;
        clearRequest();
        emit UpdateRequirement(_requirement);
    }

    function clearRequest() public
    onlyMultisigTx() {
        for (uint256 i = lastClearedId; i < requestCount; i++){
            if (requestMap[i].state == RequestState.NotConfirmed) {
                requestMap[i].state = RequestState.Canceled;
            }
        }
        lastClearedId = requestCount;
        emit ClearRequest();
    }

    function withdrawLockupStaking(address _to, uint256 _value) external 
    onlyMultisigTx() {
        uint256 withdrawableStakingAmount;
        ( , , , ,withdrawableStakingAmount) = getLockupStakingInfo();
        require(withdrawableStakingAmount >= _value, "Value is not withdrawable.");
        remainingLockupStaking = remainingLockupStaking.sub(_value);

        _to.transfer(_value);
        emit WithdrawLockupStaking(_to, _value);
    }


    function approveStakingWithdrawal(address _to, uint256 _value) external 
    onlyMultisigTx() {
        require(_value <= staking, "Value is not withdrawable.");
        uint256 approvedWithdrawalId = withdrawalRequestCount;
        withdrawalRequestMap[approvedWithdrawalId] = WithdrawalRequest({
            to : _to,
            value : _value,
            withdrawableFrom : now + ONE_WEEK,
            state: WithdrawalStakingState.Unknown
        });
        withdrawalRequestCount = withdrawalRequestCount.add(1);
        emit ApproveStakingWithdrawal(approvedWithdrawalId, _to, _value, now + ONE_WEEK);
    }


    function cancelApprovedStakingWithdrawal(uint256 _approvedWithdrawalId) external 
    onlyMultisigTx() {
        require(withdrawalRequestMap[_approvedWithdrawalId].to != 0, "Withdrawal request does not exist.");
        require(withdrawalRequestMap[_approvedWithdrawalId].state == WithdrawalStakingState.Unknown, "Invalid state.");

        withdrawalRequestMap[_approvedWithdrawalId].state = WithdrawalStakingState.Canceled;
        emit CancelApprovedStakingWithdrawal(_approvedWithdrawalId, withdrawalRequestMap[_approvedWithdrawalId].to, withdrawalRequestMap[_approvedWithdrawalId].value);
    }


    function updateRewardAddress(address _rewardAddress) external 
    onlyMultisigTx() {
        rewardAddress = _rewardAddress;
        AddressBookInterface(ADDRESS_BOOK_ADDRESS).reviseRewardAddress(_rewardAddress);
        emit UpdateRewardAddress(rewardAddress);
    }


    /*
     * Public functions
     */
    /// @notice stake KLAY
    function stakeKlay() external payable 
    afterInit() {
        require(msg.value > 0, "Invalid amount.");
        staking = staking.add(msg.value);
        emit StakeKlay(msg.sender, msg.value);
    }

    /// @notice stake KLAY fallback function
    function () external payable 
    afterInit() {
        require(msg.value > 0, "Invalid amount.");
        staking = staking.add(msg.value);
        emit StakeKlay(msg.sender, msg.value);
    }

    /// @notice 7 days after "approveStakingWithdrawal" has been requested, admins can call this function for actual withdrawal. However, it's only available for 7 days
    /// @param _approvedWithdrawalId the withdrawal ID to excute. The ID is acquired at the event log of ApproveStakingWithdrawal
    function withdrawApprovedStaking(uint256 _approvedWithdrawalId) external 
    onlyAdmin(msg.sender) {
        require(withdrawalRequestMap[_approvedWithdrawalId].to != 0, "Withdrawal request does not exist.");
        require(withdrawalRequestMap[_approvedWithdrawalId].state == WithdrawalStakingState.Unknown, "Invalid state.");
        require(withdrawalRequestMap[_approvedWithdrawalId].value <= staking, "Value is not withdrawable.");
        require(now >= withdrawalRequestMap[_approvedWithdrawalId].withdrawableFrom, "Not withdrawable yet.");
        if (now >= withdrawalRequestMap[_approvedWithdrawalId].withdrawableFrom + ONE_WEEK) {

            withdrawalRequestMap[_approvedWithdrawalId].state = WithdrawalStakingState.Canceled;
            emit CancelApprovedStakingWithdrawal(_approvedWithdrawalId, withdrawalRequestMap[_approvedWithdrawalId].to, withdrawalRequestMap[_approvedWithdrawalId].value);
        } else {
            staking = staking.sub(withdrawalRequestMap[_approvedWithdrawalId].value);
            withdrawalRequestMap[_approvedWithdrawalId].state = WithdrawalStakingState.Transferred;
            withdrawalRequestMap[_approvedWithdrawalId].to.transfer(withdrawalRequestMap[_approvedWithdrawalId].value);
            emit WithdrawApprovedStaking(_approvedWithdrawalId, withdrawalRequestMap[_approvedWithdrawalId].to, withdrawalRequestMap[_approvedWithdrawalId].value);
        }
    }


    /*
     * Private functions
     */
    /// @dev 
    function deleteFromConfirmerList(uint256 _id, address _admin) private {
        uint256 confirmersCnt = requestMap[_id].confirmers.length;
        for(uint256 i = 0; i < confirmersCnt; i++){
            if(_admin == requestMap[_id].confirmers[i]){


                if(i != confirmersCnt - 1) {
                    requestMap[_id].confirmers[i] = requestMap[_id].confirmers[confirmersCnt - 1];
                }

                delete requestMap[_id].confirmers[confirmersCnt - 1];
                requestMap[_id].confirmers.length = confirmersCnt.sub(1);
                break;
            }
        }
    }

    function submitRequest(uint256 _id, Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) private {
        requestMap[_id] = Request({
            functionId : _functionId,
            firstArg : _firstArg,
            secondArg : _secondArg,
            thirdArg : _thirdArg,
            requestProposer : msg.sender,
            confirmers : new address[](0),
            state: RequestState.NotConfirmed
        });
        emit SubmitRequest(_id, msg.sender, _functionId, _firstArg, _secondArg, _thirdArg);

        requestCount = requestCount.add(1);
    }

    function executeRequest(uint256 _id) private {
        bool executed = false;
        if (requestMap[_id].functionId == Functions.AddAdmin) {
            //bytes4(keccak256("addAdmin(address)")) => 0x70480275
            executed = address(this).call(0x70480275, address(requestMap[_id].firstArg));
        }

        if (requestMap[_id].functionId == Functions.DeleteAdmin) {
            //bytes4(keccak256("deleteAdmin(address)")) => 0x27e1f7df
            executed = address(this).call(0x27e1f7df, address(requestMap[_id].firstArg));
        }

        if (requestMap[_id].functionId == Functions.UpdateRequirement) {
            //bytes4(keccak256("updateRequirement(uint256)")) => 0xc47afb3a
            executed = address(this).call(0xc47afb3a, uint256(requestMap[_id].firstArg));
        }

        if (requestMap[_id].functionId == Functions.ClearRequest) {
            //bytes4(keccak256("clearRequest()")) => 0x4f97638f
            executed = address(this).call(0x4f97638f);
        }

        if (requestMap[_id].functionId == Functions.WithdrawLockupStaking) {
            //bytes4(keccak256("withdrawLockupStaking(address,uint256)")) => 0x505ebed4
            executed = address(this).call(0x505ebed4, address(requestMap[_id].firstArg), uint256(requestMap[_id].secondArg)); 
        }

        if (requestMap[_id].functionId == Functions.ApproveStakingWithdrawal) {
            //bytes4(keccak256("approveStakingWithdrawal(address,uint256)")) => 0x5df8b09a
            executed = address(this).call(0x5df8b09a, address(requestMap[_id].firstArg), uint256(requestMap[_id].secondArg));
        }

        if (requestMap[_id].functionId == Functions.CancelApprovedStakingWithdrawal) {
            //bytes4(keccak256("cancelApprovedStakingWithdrawal(uint256)")) => 0xc804b115
            executed = address(this).call(0xc804b115, uint256(requestMap[_id].firstArg));
        }

        if (requestMap[_id].functionId == Functions.UpdateRewardAddress) {
            //bytes4(keccak256("updateRewardAddress(address)")) => 0x944dd5a2
            executed = address(this).call(0x944dd5a2, address(requestMap[_id].firstArg));
        }

        if(executed) {
            requestMap[_id].state = RequestState.Executed;
            emit ExecuteRequestSuccess(_id, msg.sender, requestMap[_id].functionId, requestMap[_id].firstArg, requestMap[_id].secondArg, requestMap[_id].thirdArg);
        } else {
            requestMap[_id].state = RequestState.ExecutionFailed;
            emit ExecuteRequestFailure(_id, msg.sender, requestMap[_id].functionId, requestMap[_id].firstArg, requestMap[_id].secondArg, requestMap[_id].thirdArg);
        }
    }

    function checkQuorum(uint256 _id) private view returns(bool) {

        return (requestMap[_id].confirmers.length >= requirement); 
    }


    /*
     * Getter functions
     */
    /// @dev 
    function getReviewers() external view 
    beforeInit() 
    returns(address[]) {
        address[] memory reviewers = new address[](lockupConditions.reviewedCount);
        uint256 id = 0;
        if(lockupConditions.reviewedAdmin[contractValidator] == true) {
            reviewers[id] = contractValidator;
            id ++;
        }
        for(uint256 i = 0; i < adminList.length; i ++) {
            if(lockupConditions.reviewedAdmin[adminList[i]] == true) {
                reviewers[id] = adminList[i];
                id ++;
            }
        }
        return reviewers;
    }

    /// @notice Queries for request id that matches entered state
    /// @param _from beginning index
    /// @param _to last index (if 0 or greater than total request count, it loops the whole list)
    /// @param _state request state
    /// @return uint256[] request IDs satisfying the conditions
    function getRequestIds(uint256 _from, uint256 _to, RequestState _state) external view returns(uint256[]) {
        uint256 lastIndex = _to;
        if (_to == 0 || _to >= requestCount) {
            lastIndex = requestCount;
        }
        require(lastIndex >= _from);

        uint256 cnt = 0;
        uint256 i;

        for (i = _from; i < lastIndex; i++) {
            if (requestMap[i].state == _state) {
                cnt += 1;
            }  
        }
        uint256[] memory requestIds = new uint256[](cnt);
        cnt = 0;
        for (i = _from; i < lastIndex; i++) {
            if (requestMap[i].state == _state) {
                requestIds[cnt] = i;
                cnt += 1;
            }
        }
        return requestIds;
    }

    /// @notice get details of a request
    /// @param _id request ID
    /// @return function ID, first argument, second argument, third argument, request proposer, confirmers of the request, state of the request
    function getRequestInfo(uint256 _id) external view returns(
        Functions,
        bytes32,
        bytes32,
        bytes32,
        address,
        address[],
        RequestState) {
        return(
            requestMap[_id].functionId,
            requestMap[_id].firstArg,
            requestMap[_id].secondArg,
            requestMap[_id].thirdArg,
            requestMap[_id].requestProposer,
            requestMap[_id].confirmers,
            requestMap[_id].state
        );
    }

    function getLockupStakingInfo() public view
    afterInit()
    returns(uint256[], uint256[], uint256, uint256, uint256) {
        uint256 currentTime = now;
        uint256 unlockedAmount = 0;

        uint256 cnt = lockupConditions.unlockTime.length;
        for (uint256 i = 0; i < cnt; i++){
            if(currentTime > lockupConditions.unlockTime[i]) {
                unlockedAmount = unlockedAmount.add(lockupConditions.unlockAmount[i]);
            }
        }
        uint256 amountWithdrawn = initialLockupStaking.sub(remainingLockupStaking);
        uint256 withdrawableLockupStaking = unlockedAmount.sub(amountWithdrawn);

        return (lockupConditions.unlockTime, lockupConditions.unlockAmount, initialLockupStaking, remainingLockupStaking, withdrawableLockupStaking);
    }

    /// @notice loops withdrawalRequestMap and returns aprroved withdrawal ids
    /// @param _from beginning index
    /// @param _to last index (if 0 or greater than total request count, it loops the whole list)
    /// @param _state withdrawal staking state
    /// @return withdrawal IDs satisfying the conditions
    function getApprovedStakingWithdrawalIds(uint256 _from, uint256 _to, WithdrawalStakingState _state) external view returns(uint256[]) {
        uint256 lastIndex = _to;
        if (_to == 0 || _to >= withdrawalRequestCount) {
            lastIndex = withdrawalRequestCount;
        }
        require(lastIndex >= _from, "Invalid index.");

        uint256 cnt = 0;
        uint256 i;

        for (i = _from; i < lastIndex; i++) {
            if (withdrawalRequestMap[i].state == _state) {
                cnt += 1;
            }  
        }
        uint256[] memory approvedWithdrawalIds = new uint256[](cnt);
        cnt = 0;
        for (i = _from; i < lastIndex; i++) {
            if (withdrawalRequestMap[i].state == _state) {
                approvedWithdrawalIds[cnt] = i;
                cnt += 1;
            }
        }
        return approvedWithdrawalIds;
    }

    /// @notice get details of approved staking withdrawal
    /// @param _index staking withdrawal ID
    /// @return withdrawal target address, wthdrawal KLAY value, time when it becomes available, withdrawal request state
    function getApprovedStakingWithdrawalInfo(uint256 _index) external view returns(address, uint256, uint256, WithdrawalStakingState) {
        return (
            withdrawalRequestMap[_index].to,
            withdrawalRequestMap[_index].value,
            withdrawalRequestMap[_index].withdrawableFrom,
            withdrawalRequestMap[_index].state
        );
    }

    function getState() external view 
    returns(address,
            address,
            address,
            address[],
            uint256,
            uint256[],
            uint256[],
            bool,
            bool) {
        return (
            contractValidator,
            nodeId,
            rewardAddress,
            adminList,
            requirement,
            lockupConditions.unlockTime,
            lockupConditions.unlockAmount,
            lockupConditions.allReviewed,
            isInitialized
        );
    }
}

interface AddressBookInterface {
    function reviseRewardAddress(address) external;
}
