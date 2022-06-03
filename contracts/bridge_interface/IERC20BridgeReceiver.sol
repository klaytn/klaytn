// SPDX-License-Identifier: MIT

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

pragma solidity ^0.8.0;

// Implemented by bridge contract
interface IERC20BridgeReceiver {
    function onERC20Received(address _from, address _to, uint256 _amount, uint256 _feeLimit, bytes memory _extraData) external;
}

// Implemented by token contract
interface IERC20Mint {
    function mint(address to, uint256 amount) external;
}
