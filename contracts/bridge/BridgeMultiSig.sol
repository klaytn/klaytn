pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract BridgeMultiSig is Ownable {
    mapping (address => bool) public signers;
    mapping (bytes32 => mapping (address => bool)) public signedTxs; // <sha3(type, args, nonce), <singer, vote>>
    mapping (bytes32 => uint64) public signedTxsCounts; // <sha3(type, args, nonce)>
    mapping (bytes32 => uint64) public committedTxs; // <sha3(type, nonce)>
    mapping (uint64 => uint64) public signerThresholds; // <tx type>
    mapping (uint64 => uint64) public configurationNonces; // <tx type, nonce>

    enum TransactionType {
        ValueTransfer,
        Configuration,
        ConfigurationRealtime,
        Max
    }

    constructor() internal {
        for (uint64 i = 0; i < uint64(TransactionType.Max); i++) {
            signerThresholds[uint64(i)] = 1;
        }
    }

    modifier onlySigners()
    {
        require(msg.sender == owner() || signers[msg.sender]);
        _;
    }

    // onlySequentialNonce checks sequential nonce increase.
    function onlySequentialNonce(TransactionType _txType, uint64 _requestNonce) internal view {
        require(configurationNonces[uint64(_txType)] == _requestNonce, "nonce mismatch");
    }

    // voteValueTransfer votes value transfer transaction with the signer.
    function voteValueTransfer(bytes32 _txKey, bytes32 _voteKey, address _signer) internal returns(bool) {
        if (committedTxs[_txKey] != 0 || signedTxs[_voteKey][_signer]) {
            return false;
        }

        signedTxs[_voteKey][_signer] = true;
        signedTxsCounts[_voteKey]++;

        if (signedTxsCounts[_voteKey] == signerThresholds[uint64(TransactionType.ValueTransfer)]) {
            committedTxs[_txKey] = 1;
            return true;
        }

        return false;
    }

    // voteGovernanceCommon votes contract governance transaction with the signer.
    // It does not need to check committedTxs since onlySequentialNonce checks it already with harder condition.
    function voteConfigurationCommon(bytes32 _voteKey, address _signer) internal returns(bool) {
        if (signedTxs[_voteKey][_signer]) {
            return false;
        }

        signedTxs[_voteKey][_signer] = true;
        signedTxsCounts[_voteKey]++;

        return true;
    }

    // voteGovernance votes contract governance transaction with the signer.
    function voteConfiguration(bytes32 _voteKey, address _signer)
    internal
    returns(bool)
    {
        if (!voteConfigurationCommon(_voteKey, _signer)) {
            return false;
        }
        if (signedTxsCounts[_voteKey] ==  signerThresholds[uint64(TransactionType.Configuration)]) {
            configurationNonces[uint64(TransactionType.Configuration)]++;
            return true;
        }
        return false;
    }

    // voteGovernanceRealtime votes frequent contract governance transaction with the signer.
    function voteConfigurationRealtime(bytes32 _voteKey, address _signer)
    internal
    returns(bool)
    {
        if (!voteConfigurationCommon(_voteKey, _signer)) {
            return false;
        }
        if (signedTxsCounts[_voteKey] ==  signerThresholds[uint64(TransactionType.ConfigurationRealtime)]) {
            configurationNonces[uint64(TransactionType.ConfigurationRealtime)]++;
            return true;
        }
        return false;
    }
}