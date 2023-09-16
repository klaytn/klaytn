// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

import "./TreasuryRebalance.sol";

contract TreasuryRebalanceMock is TreasuryRebalance {
    constructor(uint256 _rebalanceBlockNumber) TreasuryRebalance(_rebalanceBlockNumber) {}

    function testSetAll(
        address[] calldata _retirees,
        address[] calldata _newbies,
        uint256[] calldata _amounts,
        uint256 _rebalanceBlockNumber,
        Status _status
    ) external {
        delete retirees;
        delete newbies;

        address[] memory empty = new address[](0);
        for (uint256 i = 0; i < _retirees.length; i++) {
            retirees.push(Retired(_retirees[i], empty));
        }
        for (uint256 i = 0; i < _newbies.length; i++) {
            newbies.push(Newbie(_newbies[i], _amounts[i]));
        }

        rebalanceBlockNumber = _rebalanceBlockNumber;
        status = _status;
    }
}
