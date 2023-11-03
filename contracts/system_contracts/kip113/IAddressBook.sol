// Copyright 2022 The klaytn Authors
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

interface IAddressBook {
    enum RequestState {
        Unknown,
        NotConfirmed,
        Executed,
        ExecutionFailed,
        Expired
    }
    enum Functions {
        Unknown,
        AddAdmin,
        DeleteAdmin,
        UpdateRequirement,
        ClearRequest,
        ActivateAddressBook,
        UpdatePocContract,
        UpdateKirContract,
        RegisterCnStakingContract,
        UnregisterCnStakingContract,
        UpdateSpareContract
    }

    struct Request {
        Functions functionId;
        bytes32 firstArg;
        bytes32 secondArg;
        bytes32 thirdArg;
        address[] confirmers;
        uint256 initialProposedTime;
        RequestState state;
    }

    function pocContractAddress() external view returns (address);

    function kirContractAddress() external view returns (address);

    function spareContractAddress() external view returns (address);

    function isActivated() external view returns (bool);

    function isConstructed() external view returns (bool);

    function reviseRewardAddress(address _rewardAddress) external;

    function constructContract(address[] memory _adminList, uint256 _requirement) external;

    function submitAddAdmin(address _admin) external;

    function submitDeleteAdmin(address _admin) external;

    function submitUpdateRequirement(uint256 _requirement) external;

    function submitClearRequest() external;

    function submitActivateAddressBook() external;

    function submitUpdatePocContract(address _pocContractAddress, uint256 _version) external;

    function submitUpdateKirContract(address _kirContractAddress, uint256 _version) external;

    function submitRegisterCnStakingContract(
        address _cnNodeId,
        address _cnStakingContractAddress,
        address _cnRewardAddress
    ) external;

    function submitUnregisterCnStakingContract(address _cnNodeId) external;

    function submitUpdateSpareContract(address _spareContractAddress) external;

    function revokeRequest(Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) external;

    function getState() external view returns (address[] memory adminList, uint256 requirement);

    function getPendingRequestList() external view returns (bytes32[] memory pendingRequestList);

    function getRequestInfo(
        bytes32 _id
    )
        external
        view
        returns (
            Functions functionId,
            bytes32 firstArg,
            bytes32 secondArg,
            bytes32 thirdArg,
            address[] memory confirmers,
            uint256 initialProposedTime,
            RequestState state
        );

    function getRequestInfoByArgs(
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) external view returns (bytes32 id, address[] memory confirmers, uint256 initialProposedTime, RequestState state);

    function getAllAddress() external view returns (uint8[] memory typeList, address[] memory addressList);

    function getAllAddressInfo()
        external
        view
        returns (
            address[] memory cnNodeIdList,
            address[] memory cnStakingContractList,
            address[] memory cnRewardAddressList,
            address pocContractAddress,
            address kirContractAddress
        );

    function getCnInfo(
        address _cnNodeId
    ) external view returns (address cnNodeId, address cnStakingcontract, address cnRewardAddress);
}
