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

import "../klaytn-contracts/contracts/token/ERC20/ERC20.sol";
import "../klaytn-contracts/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "../klaytn-contracts/contracts/access/AccessControl.sol";
import "../klaytn-contracts/contracts/utils/Address.sol";

import "./ERC20ServiceChain.sol";

contract ServiceChainToken is ERC20, ERC20Burnable, AccessControl, IERC20Mint, ERC20ServiceChain {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    string public constant NAME = "ServiceChainToken";
    string public constant SYMBOL = "SCT";
    uint8 public constant DECIMALS = 18;

    // one billion in initial supply
    uint256 public constant INITIAL_SUPPLY = 1000000000 * (10 ** uint256(DECIMALS));

    constructor(address _bridge) ERC20(NAME, SYMBOL) ERC20ServiceChain(_bridge) {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
        _mint(msg.sender, INITIAL_SUPPLY);
    }

    function mint(address to, uint256 amount) public onlyRole(MINTER_ROLE) {
        _mint(to, amount);
    }
}
