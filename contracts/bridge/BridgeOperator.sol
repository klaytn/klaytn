pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract BridgeOperator is Ownable {
    mapping(address => bool) public operators;
    mapping(bytes32 => mapping(address => bool)) public signedTxs; // <sha3(type, args, nonce), <singer, vote>>
    mapping(bytes32 => uint64) public signedTxsCounts; // <sha3(type, args, nonce)>
    mapping(uint64 => bool) public valueTransferTxs; // nonce
    mapping(uint8 => uint64) public operatorThresholds; // <tx type>
    uint64 public configurationNonce;

    enum VoteType {
        ValueTransfer,
        Configuration,
        Max
    }

    constructor() internal {
        for (uint64 i = 0; i < uint64(VoteType.Max); i++) {
            operatorThresholds[uint8(i)] = 1;
        }
    }

    modifier onlyOperators()
    {
        require(operators[msg.sender]);
        _;
    }

    // voteValueTransfer votes value transfer transaction with the operator.
    function voteValueTransfer(uint64 _requestNonce, bytes32 _voteKey) internal returns(bool) {
        if (valueTransferTxs[_requestNonce] || signedTxs[_voteKey][msg.sender]) {
            return false;
        }

        signedTxs[_voteKey][msg.sender] = true;
        signedTxsCounts[_voteKey]++;

        if (signedTxsCounts[_voteKey] == operatorThresholds[uint8(VoteType.ValueTransfer)]) {
            valueTransferTxs[_requestNonce] = true;
            return true;
        }

        return false;
    }

    // voteConfiguration votes contract configuration transaction with the operator.
    function voteConfiguration(bytes32 _voteKey, uint64 _requestNonce)
        internal
        returns(bool)
    {
        require(configurationNonce == _requestNonce, "nonce mismatch");

        if (signedTxs[_voteKey][msg.sender]) {
            return false;
        }

        signedTxs[_voteKey][msg.sender] = true;
        signedTxsCounts[_voteKey]++;

        if (signedTxsCounts[_voteKey] == operatorThresholds[uint8(VoteType.Configuration)]) {
            configurationNonce++;
            return true;
        }
        return false;
    }
}