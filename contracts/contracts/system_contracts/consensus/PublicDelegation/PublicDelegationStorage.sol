// Copyright 2024 The klaytn Authors
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
pragma solidity 0.8.25;

import "./IPublicDelegation.sol";

abstract contract PublicDelegationStorage is IPublicDelegation {
    /* ========== CONSTANTS ========== */

    string public constant CONTRACT_TYPE = "PublicDelegation";

    uint256 public constant VERSION = 1;

    uint256 public constant MAX_COMMISSION_RATE = 1e4; // 100%

    uint256 public constant COMMISSION_DENOMINATOR = 1e4;

    address internal constant _REGISTRY = 0x0000000000000000000000000000000000000401;

    /* ========== IMMUTABLES ========== */

    // Base CnStakingV3 of the public delegation
    ICnStakingV3 public immutable baseCnStakingV3;

    /* ========== STATE VARIABLES ========== */

    // Commission
    address public commissionTo;
    uint256 public commissionRate; // 1e4 = 100%

    // Withdrawal request
    mapping(address => uint256[]) public userRequestIds;
    mapping(uint256 => address) public requestIdToOwner;
}
