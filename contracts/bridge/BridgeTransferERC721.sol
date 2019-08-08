pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC721/IERC721.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721MetadataMintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Burnable.sol";

import "../sc_erc721/IERC721BridgeReceiver.sol";
import "./BridgeTransferCommon.sol";


contract BridgeTransferERC721 is IERC721BridgeReceiver, BridgeTransfer {
    // handleERC721Transfer sends the ERC721 by the request.
    function handleERC721Transfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _tokenId,
        uint64 _requestedNonce,
        uint64 _requestedBlockNumber,
        string _tokenURI,
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
            TokenType.ERC721,
            _from,
            _to,
            _tokenAddress,
            _tokenId,
            _requestedNonce,
            _extraData
        );
        _setHandledRequestTxHash(_requestTxHash);
        lastHandledRequestBlockNumber = _requestedBlockNumber;

        updateHandleNonce(_requestedNonce);

        if (modeMintBurn) {
            ERC721MetadataMintable(_tokenAddress).mintWithTokenURI(_to, _tokenId, _tokenURI);
        } else {
            IERC721(_tokenAddress).safeTransferFrom(address(this), _to, _tokenId);
        }
    }

    // _requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function _requestERC721Transfer(
        address _tokenAddress,
        address _from,
        address _to,
        uint256 _tokenId,
        uint256[] _extraData
    )
        internal
    {
        require(isRunning, "stopped bridge");
        require(allowedTokens[_tokenAddress] != address(0), "invalid token");

        string memory uri = ERC721Metadata(_tokenAddress).tokenURI(_tokenId);

        if (modeMintBurn) {
            ERC721Burnable(_tokenAddress).burn(_tokenId);
        }

        emit RequestValueTransfer(
            TokenType.ERC721,
            _from,
            _to,
            _tokenAddress,
            _tokenId,
            requestNonce,
            uri,
            0,
            _extraData
        );
        requestNonce++;
    }

    // onERC721Received function of ERC721 token for 1-step deposits to the Bridge
    function onERC721Received(
        address _from,
        uint256 _tokenId,
        address _to,
        uint256[] _extraData
    )
        public
    {
        _requestERC721Transfer(msg.sender, _from, _to, _tokenId, _extraData);
    }

    // requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function requestERC721Transfer(
        address _tokenAddress,
        address _to,
        uint256 _tokenId,
        uint256[] _extraData
    )
        external
    {
        IERC721(_tokenAddress).transferFrom(msg.sender, address(this), _tokenId);
        _requestERC721Transfer(_tokenAddress, msg.sender, _to, _tokenId, _extraData);
    }
}
