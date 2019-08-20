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

pragma solidity ^0.4.24;
import "./SafeMath.sol";

/**
 * @title AddressBook
 */

contract AddressBook {
    using SafeMath for uint256;
    /*
     *  Events
     */
    event DeployContract(string contractType, address[] adminList, uint256 requirement);
    event AddAdmin (address indexed admin);
    event DeleteAdmin(address indexed admin);
    event UpdateRequirement(uint256 requirement);
    event ClearRequest();
    event SubmitRequest(bytes32 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers);
    event ExpiredRequest(bytes32 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers);
    event RevokeRequest(bytes32 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers);
    event CancelRequest(bytes32 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg);
    event ExecuteRequestSuccess(bytes32 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg);
    event ExecuteRequestFailure(bytes32 indexed id, address indexed from, Functions functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg);

    event ActivateAddressBook();
    event UpdatePocContract(address prevPocContractAddress, uint256 prevVersion, address curPocContractAddress, uint256 curVersion);
    event UpdateKirContract(address prevKirContractAddress, uint256 prevVersion, address curKirContractAddress, uint256 curVersion);
    event UpdateSpareContract(address spareContractAddress);
    event RegisterCnStakingContract(address cnNodeId, address cnStakingContractAddress, address cnRewardAddress);
    event UnregisterCnStakingContract(address cnNodeId);
    event ReviseRewardAddress(address cnNodeId, address prevRewardAddress, address curRewardAddress);

    /*
     *  Constants
     */
    uint256 constant public MAX_ADMIN = 50;
    uint256 constant public MAX_PENDING_REQUEST = 100;
    string constant public CONTRACT_TYPE = "AddressBook";
    uint8 constant public CN_NODE_ID_TYPE = 0; 
    uint8 constant public CN_STAKING_ADDRESS_TYPE = 1; 
    uint8 constant public CN_REWARD_ADDRESS_TYPE = 2; 
    uint8 constant public POC_CONTRACT_TYPE = 3; 
    uint8 constant public KIR_CONTRACT_TYPE = 4;
    uint256 constant public ONE_WEEK = 1 weeks;
    uint256 constant public TWO_WEEKS = 2 weeks;
    uint256 constant public VERSION = 1;

    enum RequestState {Unknown, NotConfirmed, Executed, ExecutionFailed, Expired}
    enum Functions {Unknown, AddAdmin, DeleteAdmin, UpdateRequirement, ClearRequest, ActivateAddressBook, UpdatePocContract, UpdateKirContract, RegisterCnStakingContract, UnregisterCnStakingContract, UpdateSpareContract}

    struct Request {
        Functions functionId;
        bytes32 firstArg;
        bytes32 secondArg;
        bytes32 thirdArg;
        address[] confirmers; 
        uint256 initialProposedTime;
        RequestState state;
    }

    /*
     *  Storage
     */
    address[] private adminList;
    uint256 public requirement;
    mapping (address => bool) private isAdmin;
    mapping(bytes32 => Request) private requestMap;
    bytes32 [] private pendingRequestList;

    address public pocContractAddress;
    address public kirContractAddress;
    address public spareContractAddress;

    mapping(address => uint256) private cnIndexMap;
    address[] private cnNodeIdList;
    address[] private cnStakingContractList;
    address[] private cnRewardAddressList;

    bool public isActivated;
    bool public isConstructed;

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
        require(!isAdmin[_admin], "Admin already exits.");
        _;
    }

    modifier notNull(address _address) {
        require(_address != 0, "Address is null.");
        _;
    }

    modifier validRequirement(uint256 _adminCount, uint256 _requirement) {
        require(_adminCount <= MAX_ADMIN
            && _requirement <= _adminCount
            && _requirement != 0
            && _adminCount != 0, "Invalid requirement.");
        _;
    }


    /*
     *  Constructor
     */
    function constructContract(address[] _adminList, uint256 _requirement) external 
    validRequirement(_adminList.length, _requirement) {
        require(msg.sender == 0x88bb3838aa0a140aCb73EEb3d4B25a8D3aFD58D4, "Invalid sender.");
        require(isConstructed == false, "Already constructed.");
        uint256 adminListCnt = _adminList.length;

        isActivated = false;
        for (uint256 i = 0; i < adminListCnt; i++) {
            require(!isAdmin[_adminList[i]] && _adminList[i] != 0, "Address is null or not unique.");
            isAdmin[_adminList[i]] = true;
        }
        adminList = _adminList;
        requirement = _requirement;
        isConstructed = true;
        emit DeployContract(CONTRACT_TYPE, adminList, requirement);
    }

    /*
     *  Private functions
     */
    function deleteFromPendingRequestList(bytes32 id) private {
        uint256 pendingRequestListCnt = pendingRequestList.length;
        for(uint256 i = 0; i < pendingRequestListCnt; i++){
            if(id == pendingRequestList[i]){
                if(i != pendingRequestListCnt - 1) {
                    pendingRequestList[i] = pendingRequestList[pendingRequestListCnt - 1];
                }
                delete pendingRequestList[pendingRequestListCnt - 1];
                pendingRequestList.length = pendingRequestList.length.sub(1);
                break;
            }
        }
    }

    function getId(Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) private pure returns(bytes32) {
        return keccak256(abi.encodePacked(_functionId,_firstArg,_secondArg,_thirdArg));
    }

    function submitRequest(Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) private {
        bytes32 id = getId(_functionId,_firstArg,_secondArg,_thirdArg);

        if(requestMap[id].initialProposedTime != 0) {
            if (requestMap[id].initialProposedTime + TWO_WEEKS < now) {
                deleteFromPendingRequestList(id);
                delete requestMap[id];
            }
            else if (requestMap[id].initialProposedTime + ONE_WEEK < now) {
                if (requestMap[id].state != RequestState.Expired) {
                    requestMap[id].state = RequestState.Expired;
                }
                emit ExpiredRequest(id, msg.sender, _functionId, _firstArg, _secondArg, _thirdArg, requestMap[id].confirmers);
            }
            // Confirm
            else if (requestMap[id].initialProposedTime <= now){
                uint256 confirmersCnt = requestMap[id].confirmers.length;
                for(uint256 i = 0; i < confirmersCnt; i++){
                    require(msg.sender != requestMap[id].confirmers[i], "Msg.sender already requested.");
                }
                requestMap[id].confirmers.push(msg.sender);
                emit SubmitRequest(id, msg.sender, _functionId, _firstArg, _secondArg, _thirdArg, requestMap[id].confirmers);
            }
        }

        if (requestMap[id].initialProposedTime == 0) {
            if (pendingRequestList.length >= MAX_PENDING_REQUEST) {
                require(_functionId == Functions.ClearRequest, "Request list is full.");
            }
            requestMap[id] = Request({
                functionId : _functionId,
                firstArg : _firstArg,
                secondArg : _secondArg,
                thirdArg : _thirdArg,
                initialProposedTime : now,
                confirmers : new address[](0),
                state : RequestState.NotConfirmed
            });
            requestMap[id].confirmers.push(msg.sender);
            pendingRequestList.push(id);
            emit SubmitRequest(id, msg.sender, _functionId, _firstArg, _secondArg, _thirdArg, requestMap[id].confirmers);
        }
    }
    
    function executeRequest(bytes32 _id) private {
        bool executed = false;
        Request memory _executeRequest = requestMap[_id];

        if (_executeRequest.functionId == Functions.AddAdmin) {
            //bytes4(keccak256("addAdmin(address)")) => 0x70480275
            executed = address(this).call(0x70480275, address(_executeRequest.firstArg));
        }
        else if (_executeRequest.functionId == Functions.DeleteAdmin) {
            //bytes4(keccak256("deleteAdmin(address)")) => 0x27e1f7df
            executed = address(this).call(0x27e1f7df, address(_executeRequest.firstArg));
        }
        else if (_executeRequest.functionId == Functions.UpdateRequirement) {
            //bytes4(keccak256("updateRequirement(uint256)")) => 0xc47afb3a
            executed = address(this).call(0xc47afb3a, uint256(_executeRequest.firstArg));
        }
        else if (_executeRequest.functionId == Functions.ClearRequest) {
            //bytes4(keccak256("clearRequest()")) => 0x4f97638f
            executed = address(this).call(0x4f97638f);
        }
        else if (_executeRequest.functionId == Functions.ActivateAddressBook) {
            //bytes4(keccak256("activateAddressBook()")) => 0xcec92466
            executed = address(this).call(0xcec92466);
        }
        else if (_executeRequest.functionId == Functions.UpdatePocContract) {
            //bytes4(keccak256("updatePocContract(address,uint256)")) => 0xc7e9de75
            executed = address(this).call(0xc7e9de75, address(_executeRequest.firstArg), uint256(_executeRequest.secondArg));
        }
        else if (_executeRequest.functionId == Functions.UpdateKirContract) {
            //bytes4(keccak256("updateKirContract(address,uint256)")) => 0x4c5d435c
            executed = address(this).call(0x4c5d435c, address(_executeRequest.firstArg), uint256(_executeRequest.secondArg));
        }
        else if (_executeRequest.functionId == Functions.RegisterCnStakingContract) {
            //bytes4(keccak256("registerCnStakingContract(address,address,address)")) => 0x298b3c61
            executed = address(this).call(0x298b3c61, address(_executeRequest.firstArg), address(_executeRequest.secondArg), address(_executeRequest.thirdArg));
        }
        else if (_executeRequest.functionId == Functions.UnregisterCnStakingContract) {
            //bytes4(keccak256("unregisterCnStakingContract(address)")) => 0x579740db
            executed = address(this).call(0x579740db, address(_executeRequest.firstArg));
        }
        else if (_executeRequest.functionId == Functions.UpdateSpareContract) {
            //bytes4(keccak256("updateSpareContract(address)")) => 0xafaaf330
            executed = address(this).call(0xafaaf330, address(_executeRequest.firstArg));
        }

        deleteFromPendingRequestList(_id);
        if(executed) {
            if(requestMap[_id].initialProposedTime != 0) {
                requestMap[_id].state = RequestState.Executed;
            }
            emit ExecuteRequestSuccess(_id, msg.sender, _executeRequest.functionId, _executeRequest.firstArg, _executeRequest.secondArg, _executeRequest.thirdArg);
        } else {
            if(requestMap[_id].initialProposedTime != 0) {
                requestMap[_id].state = RequestState.ExecutionFailed;
            }
            emit ExecuteRequestFailure(_id, msg.sender, _executeRequest.functionId, _executeRequest.firstArg, _executeRequest.secondArg, _executeRequest.thirdArg);
        }
    }

    function checkQuorum(bytes32 _id) private view returns(bool) {
        return (requestMap[_id].confirmers.length >= requirement); 
    }

    /*
     *  external functions
     */
    function revokeRequest(Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) external 
    onlyAdmin(msg.sender) {
        bytes32 id = getId(_functionId,_firstArg,_secondArg,_thirdArg);

        require(requestMap[id].initialProposedTime != 0, "Invalid request.");
        require(requestMap[id].state == RequestState.NotConfirmed, "Must be at not-confirmed state.");
        bool foundIt = false;
        uint256 confirmerCnt = requestMap[id].confirmers.length;

        for(uint256 i = 0; i < confirmerCnt; i++){
            if(msg.sender == requestMap[id].confirmers[i]){
                foundIt = true;

                if (requestMap[id].initialProposedTime + ONE_WEEK < now) {
                    if (requestMap[id].initialProposedTime + TWO_WEEKS < now) {
                        deleteFromPendingRequestList(id);
                        delete requestMap[id];
                    }
                    else {
                        requestMap[id].state = RequestState.Expired;
                    }

                    emit ExpiredRequest(id, msg.sender, _functionId, _firstArg, _secondArg, _thirdArg, requestMap[id].confirmers);
                }
                else {
                    if(i != confirmerCnt - 1) {
                        requestMap[id].confirmers[i] = requestMap[id].confirmers[confirmerCnt - 1];
                    }
                    delete requestMap[id].confirmers[confirmerCnt - 1];
                    requestMap[id].confirmers.length = requestMap[id].confirmers.length.sub(1);

                    emit RevokeRequest(id, msg.sender, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg, requestMap[id].confirmers);

                    if(requestMap[id].confirmers.length == 0) {
                        deleteFromPendingRequestList(id);
                        delete requestMap[id];
                        emit CancelRequest(id, msg.sender, requestMap[id].functionId, requestMap[id].firstArg, requestMap[id].secondArg, requestMap[id].thirdArg);
                    }
                }
                break;
            }
        }
        require(foundIt, "Msg.sender has not requested.");
    }

    /*
     *  submit request functions
     */
    function submitAddAdmin(address _admin) external 
    onlyAdmin(msg.sender)
    adminDoesNotExist(_admin)
    notNull(_admin)
    validRequirement(adminList.length.add(1), requirement) {
        bytes32 id = getId(Functions.AddAdmin,bytes32(_admin),0,0);

        submitRequest(Functions.AddAdmin,bytes32(_admin),0,0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitDeleteAdmin(address _admin) external 
    onlyAdmin(_admin)
    onlyAdmin(msg.sender)
    notNull(_admin)
    validRequirement(adminList.length.sub(1), requirement) {
        bytes32 id = getId(Functions.DeleteAdmin,bytes32(_admin),0,0);

        submitRequest(Functions.DeleteAdmin,bytes32(_admin),0,0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUpdateRequirement(uint256 _requirement) external 
    onlyAdmin(msg.sender)
    validRequirement(adminList.length, _requirement) {
        require(requirement != _requirement, "Same requirement.");
        bytes32 id = getId(Functions.UpdateRequirement,bytes32(_requirement),0,0);

        submitRequest(Functions.UpdateRequirement,bytes32(_requirement),0,0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitClearRequest() external 
    onlyAdmin(msg.sender) {
        bytes32 id = getId(Functions.ClearRequest,0,0,0);

        submitRequest(Functions.ClearRequest,0,0,0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitActivateAddressBook() external
    onlyAdmin(msg.sender) {
        require(isActivated == false, "Already activated.");
        require(adminList.length != 0, "No admin is listed.");
        require(pocContractAddress != 0, "PoC contract is not registered.");
        require(kirContractAddress != 0, "KIR contract is not registered.");
        require(cnNodeIdList.length != 0, "No node ID is listed.");
        require(cnNodeIdList.length == cnStakingContractList.length, "Invalid length between node IDs and staking contracts.");
        require(cnStakingContractList.length == cnRewardAddressList.length, "Invalid length between staking contracts and reward addresses.");

        bytes32 id = getId(Functions.ActivateAddressBook,0,0,0);

        submitRequest(Functions.ActivateAddressBook,0,0,0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUpdatePocContract(address _pocContractAddress, uint256 _version) external 
    notNull(_pocContractAddress)
    onlyAdmin(msg.sender) {
        require(PocContractInterface(_pocContractAddress).getPocVersion() == _version, "Invalid PoC version.");
        
        bytes32 id = getId(Functions.UpdatePocContract,bytes32(_pocContractAddress),bytes32(_version),0);

        submitRequest(Functions.UpdatePocContract,bytes32(_pocContractAddress),bytes32(_version),0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUpdateKirContract(address _kirContractAddress, uint256 _version) external 
    notNull(_kirContractAddress)
    onlyAdmin(msg.sender) {
        require(KirContractInterface(_kirContractAddress).getKirVersion() == _version, "Invalid KIR version.");
        
        bytes32 id = getId(Functions.UpdateKirContract,bytes32(_kirContractAddress),bytes32(_version),0);

        submitRequest(Functions.UpdateKirContract,bytes32(_kirContractAddress),bytes32(_version),0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUpdateSpareContract(address _spareContractAddress) external 
    onlyAdmin(msg.sender) {        
        bytes32 id = getId(Functions.UpdateSpareContract,bytes32(_spareContractAddress),0,0);

        submitRequest(Functions.UpdateSpareContract,bytes32(_spareContractAddress),0,0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitRegisterCnStakingContract(address _cnNodeId, address _cnStakingContractAddress, address _cnRewardAddress) external
    notNull(_cnNodeId)
    notNull(_cnStakingContractAddress)
    notNull(_cnRewardAddress)  
    onlyAdmin(msg.sender) {
        if (cnNodeIdList.length > 0) {
            require(cnNodeIdList[cnIndexMap[_cnNodeId]] != _cnNodeId, "CN node ID already exist.");
        }
        require(CnStakingContractInterface(_cnStakingContractAddress).nodeId() == _cnNodeId, "Invalid CN node ID.");
        require(CnStakingContractInterface(_cnStakingContractAddress).rewardAddress() == _cnRewardAddress, "Invalid CN reward address.");
        require(CnStakingContractInterface(_cnStakingContractAddress).isInitialized() == true, "CN contract is not initialized.");

        bytes32 id = getId(Functions.RegisterCnStakingContract,bytes32(_cnNodeId),bytes32(_cnStakingContractAddress),bytes32(_cnRewardAddress));

        submitRequest(Functions.RegisterCnStakingContract,bytes32(_cnNodeId),bytes32(_cnStakingContractAddress),bytes32(_cnRewardAddress));
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUnregisterCnStakingContract(address _cnNodeId) external 
    notNull(_cnNodeId) 
    onlyAdmin(msg.sender) {
        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        require(cnNodeIdList.length > 1, "CN should be more than one.");

        bytes32 id = getId(Functions.UnregisterCnStakingContract,bytes32(_cnNodeId),0,0);

        submitRequest(Functions.UnregisterCnStakingContract,bytes32(_cnNodeId),0,0);
        if(checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
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
        uint256 adminCnt = adminList.length;
        isAdmin[_admin] = false;
        
        for (uint256 i=0; i < adminCnt - 1; i++) {
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
        require(requirement != _requirement, "Same requirement.");
        requirement = _requirement;
        clearRequest();
        emit UpdateRequirement(_requirement);
    }

    function clearRequest() public 
    onlyMultisigTx() {
        uint256 pendingRequestCnt = pendingRequestList.length;

        for (uint256 i = 0; i < pendingRequestCnt; i++){
            delete requestMap[pendingRequestList[i]];
        }
        delete pendingRequestList;
        emit ClearRequest();
    }

    function activateAddressBook() external 
    onlyMultisigTx() {
        require(isActivated == false, "Already activated.");
        require(adminList.length != 0, "No admin is listed.");
        require(pocContractAddress != 0, "PoC contract is not registered.");
        require(kirContractAddress != 0, "KIR contract is not registered.");
        require(cnNodeIdList.length != 0, "No node ID is listed.");
        require(cnNodeIdList.length == cnStakingContractList.length, "Invalid length between node IDs and staking contracts.");
        require(cnStakingContractList.length == cnRewardAddressList.length, "Invalid length between staking contracts and reward addresses.");
        isActivated = true;
        
        emit ActivateAddressBook();
    }

    function updatePocContract(address _pocContractAddress, uint256 _version) external 
    onlyMultisigTx() {
        require(PocContractInterface(_pocContractAddress).getPocVersion() == _version, "Invalid PoC version.");
        
        address prevPocContractAddress = pocContractAddress;
        pocContractAddress = _pocContractAddress;
        uint256 prevVersion = 0;

        if(prevPocContractAddress != 0) {
            prevVersion = PocContractInterface(prevPocContractAddress).getPocVersion();
        }
        emit UpdatePocContract(prevPocContractAddress, prevVersion, _pocContractAddress, _version);
    }

    function updateKirContract(address _kirContractAddress, uint256 _version) external 
    onlyMultisigTx() {
        require(KirContractInterface(_kirContractAddress).getKirVersion() == _version, "Invalid KIR version.");
        
        address prevKirContractAddress = kirContractAddress;
        kirContractAddress = _kirContractAddress;
        uint256 prevVersion = 0;

        if(prevKirContractAddress != 0) {
            prevVersion = KirContractInterface(prevKirContractAddress).getKirVersion();
        }
        emit UpdateKirContract(prevKirContractAddress, prevVersion, _kirContractAddress, _version);
    }

    function updateSpareContract(address _spareContractAddress) external 
    onlyMultisigTx() {
        spareContractAddress = _spareContractAddress;
        emit UpdateSpareContract(spareContractAddress);
    }

    function registerCnStakingContract(address _cnNodeId, address _cnStakingContractAddress, address _cnRewardAddress) external 
    onlyMultisigTx() {
        if (cnNodeIdList.length > 0) {
            require(cnNodeIdList[cnIndexMap[_cnNodeId]] != _cnNodeId, "CN node ID already exist.");
        }
        require(CnStakingContractInterface(_cnStakingContractAddress).nodeId() == _cnNodeId, "Invalid CN node ID.");
        require(CnStakingContractInterface(_cnStakingContractAddress).rewardAddress() == _cnRewardAddress, "Invalid CN reward address.");
        require(CnStakingContractInterface(_cnStakingContractAddress).isInitialized() == true, "CN contract is not initialized.");

        uint256 index = cnNodeIdList.length;
        cnIndexMap[_cnNodeId] = index;
        cnNodeIdList.push(_cnNodeId);
        cnStakingContractList.push(_cnStakingContractAddress);
        cnRewardAddressList.push(_cnRewardAddress);

        emit RegisterCnStakingContract(_cnNodeId, _cnStakingContractAddress, _cnRewardAddress);
    }

    function unregisterCnStakingContract(address _cnNodeId) external 
    onlyMultisigTx() {
        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        require(cnNodeIdList.length > 1, "CN should be more than one.");

        if (index < cnNodeIdList.length - 1) {
            cnNodeIdList[index] = cnNodeIdList[cnNodeIdList.length-1];
            cnStakingContractList[index] = cnStakingContractList[cnNodeIdList.length-1];
            cnRewardAddressList[index] = cnRewardAddressList[cnNodeIdList.length-1];

            cnIndexMap[cnNodeIdList[cnNodeIdList.length-1]] = index;
        }

        delete cnIndexMap[_cnNodeId];
        delete cnNodeIdList[cnNodeIdList.length-1];
        cnNodeIdList.length = cnNodeIdList.length.sub(1);
        delete cnStakingContractList[cnStakingContractList.length-1];
        cnStakingContractList.length = cnStakingContractList.length.sub(1);
        delete cnRewardAddressList[cnRewardAddressList.length-1];
        cnRewardAddressList.length = cnRewardAddressList.length.sub(1);

        emit UnregisterCnStakingContract(_cnNodeId);
    }

    /*
     * External function
     */
    function reviseRewardAddress(address _rewardAddress) external 
    notNull(_rewardAddress) {
        bool foundIt = false;
        uint256 index = 0;
        uint256 cnStakingContractListCnt = cnStakingContractList.length;
        for(uint256 i = 0; i < cnStakingContractListCnt; i++){
            if (cnStakingContractList[i] == msg.sender) {
                foundIt = true;
                index = i;
                break;
            }
        }
        require(foundIt, "Msg.sender is not CN contract.");
        address prevAddress = cnRewardAddressList[index];
        cnRewardAddressList[index] = _rewardAddress;

        emit ReviseRewardAddress(cnNodeIdList[index], prevAddress, cnRewardAddressList[index]);
    }

    /*
     * Getter functions
     */
    function getState() external view returns(address[], uint256) {
        return (adminList, requirement);
    }

    function getPendingRequestList() external view returns(bytes32[]) {
        return pendingRequestList;
    }

    function getRequestInfo(bytes32 _id) external view returns(Functions,bytes32,bytes32,bytes32,address[],uint256,RequestState) {
        return(
            requestMap[_id].functionId,
            requestMap[_id].firstArg,
            requestMap[_id].secondArg,
            requestMap[_id].thirdArg,
            requestMap[_id].confirmers,
            requestMap[_id].initialProposedTime,
            requestMap[_id].state
        );
    }

    function getRequestInfoByArgs(Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) external view returns(bytes32,address[],uint256,RequestState) {
        bytes32 _id = getId(_functionId,_firstArg,_secondArg,_thirdArg);
        return(
            _id,
            requestMap[_id].confirmers,
            requestMap[_id].initialProposedTime,
            requestMap[_id].state
        );
    }

    function getAllAddress() external view returns(uint8[], address[]) {
        uint8[] memory typeList;
        address[] memory addressList;
        if(isActivated == false) {
            typeList = new uint8[](0);
            addressList = new address[](0);
        } else {
            typeList = new uint8[](cnNodeIdList.length * 3 + 2);
            addressList = new address[](cnNodeIdList.length * 3 + 2);
            uint256 cnNodeCnt = cnNodeIdList.length;
            for (uint256 i = 0; i < cnNodeCnt; i ++){
                //add node id and its type number to array
                typeList[i * 3] = uint8(CN_NODE_ID_TYPE);
                addressList[i * 3] = address(cnNodeIdList[i]);
                //add staking address and its type number to array
                typeList[i * 3 + 1] = uint8(CN_STAKING_ADDRESS_TYPE);
                addressList[i * 3 + 1] = address(cnStakingContractList[i]);
                //add reward address and its type number to array
                typeList[i * 3 + 2] = uint8(CN_REWARD_ADDRESS_TYPE);
                addressList[i * 3 + 2] = address(cnRewardAddressList[i]);
            }
            typeList[cnNodeCnt *3] = uint8(POC_CONTRACT_TYPE);
            addressList[cnNodeCnt * 3] = address(pocContractAddress);
            typeList[cnNodeCnt * 3 + 1] = uint8(KIR_CONTRACT_TYPE);
            addressList[cnNodeCnt * 3 + 1] = address(kirContractAddress);
        }
        return (typeList, addressList);
    }

    function getAllAddressInfo() external view returns(address[], address[], address[], address, address) {
        return (cnNodeIdList, cnStakingContractList, cnRewardAddressList, pocContractAddress, kirContractAddress); 
    }

    function getCnInfo(address _cnNodeId) external notNull(_cnNodeId) view returns(address, address, address) {
        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        return(cnNodeIdList[index], cnStakingContractList[index], cnRewardAddressList[index]);
    }
}

interface CnStakingContractInterface {
    function nodeId() external view returns(address);
    function rewardAddress() external view returns(address);
    function isInitialized() external view returns(bool);
}
interface PocContractInterface {
    function getPocVersion() external pure returns(uint256);
}
interface KirContractInterface {
    function getKirVersion() external pure returns(uint256);
}
