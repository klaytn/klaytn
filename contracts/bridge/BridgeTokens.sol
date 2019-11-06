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

pragma solidity 0.5.6;

import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract BridgeTokens is Ownable {
    mapping(address => address) public registeredTokens; // <token, counterpart token>
    mapping(address => uint) public indexOfTokens; // <token, index>
    address[] public registeredTokenList;
    mapping(address => bool) public lockedTokens;

    event TokenRegistered(address indexed token);
    event TokenDeregistered(address indexed token);
    event TokenLocked(address indexed token);
    event TokenUnlocked(address indexed token);

    modifier onlyRegisteredToken(address _token) {
        require(registeredTokens[_token] != address(0), "not allowed token");
        _;
    }

    modifier onlyNotRegisteredToken(address _token) {
        require(registeredTokens[_token] == address(0), "allowed token");
        _;
    }

    modifier onlyLockedToken(address _token) {
        require(lockedTokens[_token] == true, "unlocked token");
        _;
    }

    modifier onlyUnlockedToken(address _token) {
        require(lockedTokens[_token] == false, "locked token");
        _;
    }

    function getRegisteredTokenList() external view returns(address[] memory) {
        return registeredTokenList;
    }

    // registerToken can update the allowed token with the counterpart token.
    function registerToken(address _token, address _cToken)
        external
        onlyOwner
        onlyNotRegisteredToken(_token)
    {
        registeredTokens[_token] = _cToken;
        indexOfTokens[_token] = registeredTokenList.length;
        registeredTokenList.push(_token);

        emit TokenRegistered(_token);
    }

    // deregisterToken can remove the token in registeredToken list.
    function deregisterToken(address _token)
        external
        onlyOwner
        onlyRegisteredToken(_token)
    {
        delete registeredTokens[_token];
        delete lockedTokens[_token];

        uint idx = indexOfTokens[_token];
        delete indexOfTokens[_token];

        if (idx < registeredTokenList.length-1) {
            registeredTokenList[idx] = registeredTokenList[registeredTokenList.length-1];
            indexOfTokens[registeredTokenList[idx]] = idx;
        }
        registeredTokenList.length--;

        emit TokenDeregistered(_token);
    }

    // lockToken can lock the token to prevent request token transferring.
    function lockToken(address _token)
        external
        onlyOwner
        onlyRegisteredToken(_token)
        onlyUnlockedToken(_token)
    {
        lockedTokens[_token] = true;

        emit TokenLocked(_token);
    }

    // unlockToken can unlock the token to  request token transferring.
    function unlockToken(address _token)
        external
        onlyOwner
        onlyRegisteredToken(_token)
        onlyLockedToken(_token)
    {
        delete lockedTokens[_token];

        emit TokenUnlocked(_token);
    }
}
