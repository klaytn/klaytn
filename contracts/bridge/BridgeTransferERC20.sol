pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";

import "../sc_erc20/IERC20BridgeReceiver.sol";
import "./BridgeTransferCommon.sol";


contract BridgeTransferERC20 is IERC20BridgeReceiver, BridgeTransfer {
    // handleERC20Transfer sends the token by the request.
    function handleERC20Transfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _value,
        uint64 _requestedNonce,
        uint64 _requestedBlockNumber,
        uint256[] _extraData
    )
        public
        onlyOperators
    {
        if (!voteValueTransfer(_requestedNonce)) {
            return;
        }

        emit HandleValueTransfer(
            _requestTxHash,
            TokenType.ERC20,
            _from,
            _to,
            _tokenAddress,
            _value,
            _requestedNonce,
            _extraData
        );
        _setHandledRequestTxHash(_requestTxHash);
        lastHandledRequestBlockNumber = _requestedBlockNumber;

        updateHandleNonce(_requestedNonce);

        if (modeMintBurn) {
            ERC20Mintable(_tokenAddress).mint(_to, _value);
        } else {
            IERC20(_tokenAddress).transfer(_to, _value);
        }
    }

    // _requestERC20Transfer requests transfer ERC20 to _to on relative chain.
    function _requestERC20Transfer(
        address _tokenAddress,
        address _from,
        address _to,
        uint256 _value,
        uint256 _feeLimit,
        uint256[] _extraData
    )
        internal
    {
        require(isRunning, "stopped bridge");
        require(_value > 0, "zero msg.value");
        require(allowedTokens[_tokenAddress] != address(0), "invalid token");

        uint256 fee = _payERC20FeeAndRefundChange(_from, _tokenAddress, _feeLimit);

        if (modeMintBurn) {
            ERC20Burnable(_tokenAddress).burn(_value);
        }

        emit RequestValueTransfer(
            TokenType.ERC20,
            _from,
            _to,
            _tokenAddress,
            _value,
            requestNonce,
            "",
            fee,
            _extraData
        );
        requestNonce++;
    }

    // onERC20Received function of ERC20 token for 1-step deposits to the Bridge.
    function onERC20Received(
        address _from,
        address _to,
        uint256 _value,
        uint256 _feeLimit,
        uint256[] _extraData
    )
        public
    {
        _requestERC20Transfer(msg.sender, _from, _to, _value, _feeLimit, _extraData);
    }

    // requestERC20Transfer requests transfer ERC20 to _to on relative chain.
    function requestERC20Transfer(
        address _tokenAddress,
        address _to,
        uint256 _value,
        uint256 _feeLimit,
        uint256[] _extraData
    )
        external
    {
        IERC20(_tokenAddress).transferFrom(msg.sender, address(this), _value.add(_feeLimit));
        _requestERC20Transfer(_tokenAddress, msg.sender, _to, _value, _feeLimit, _extraData);
    }


    // setERC20Fee sets the fee of the token transfer.
    function setERC20Fee(address _token, uint256 _fee, uint64 _requestNonce)
        external
        onlyOperators
    {
        if (!voteConfiguration(_requestNonce)) {
            return;
        }
        _setERC20Fee(_token, _fee);
    }
}
