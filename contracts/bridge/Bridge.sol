pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";

import "../externals/openzeppelin-solidity/contracts/token/ERC721/IERC721.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721MetadataMintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Burnable.sol";

import "../sc_erc721/IERC721BridgeReceiver.sol";
import "../sc_erc20/IERC20BridgeReceiver.sol";
import "./BridgeFee.sol";
import "./BridgeOperator.sol";
import "./HandledRequests.sol";

contract Bridge is IERC20BridgeReceiver, IERC721BridgeReceiver, BridgeFee, BridgeOperator, HandledRequests {
    uint64 public constant VERSION = 1;
    bool public modeMintBurn = false;
    address public counterpartBridge;
    bool public isRunning;

    uint64 public requestNonce;
    uint64 public lastHandledRequestBlockNumber;
    uint64 public sequentialHandleNonce;
    uint64 public maxHandledRequestedNonce;
    mapping(uint64 => bool) public handledNonces;  // <handled nonce> history

    mapping(address => address) public allowedTokens; // <token, counterpart token>

    using SafeMath for uint256;

    enum TokenType {
        KLAY,
        ERC20,
        ERC721
    }

    // TODO-Klaytn-Service FeeReceiver should be passed by argument of constructor.
    constructor (bool _modeMintBurn) BridgeFee(address(0)) public payable {
        isRunning = true;
        modeMintBurn = _modeMintBurn;
    }

    /**
     * Event to log the request value transfer from the Bridge.
     * @param tokenType is the type of tokens (KLAY/ERC20/ERC721).
     * @param from is the requester of the request value transfer event.
     * @param to is the receiver of the value.
     * @param tokenAddress Address of token contract the token belong to.
     * @param valueOrTokenId is the value of KLAY/ERC20 or token ID of ERC721.
     * @param requestNonce is the order number of the request value transfer.
     * @param uri is uri of ERC721 token.
     */
    event RequestValueTransfer(
        TokenType tokenType,
        address from,
        address to,
        address tokenAddress,
        uint256 valueOrTokenId,
        uint64 requestNonce,
        string uri,
        uint256 fee
    );

    /**
     * Event to log the handle value transfer from the Bridge.
     * @param requestTxHash is a transaction hash of request value transfer.
     * @param tokenType is the type of tokens (KLAY/ERC20/ERC721).
     * @param from is an address of the account who requested the value transfer.
     * @param to is an address of the account who will received the value.
     * @param tokenAddress Address of token contract the token belong to.
     * @param valueOrTokenId is the value of KLAY/ERC20 or token ID of ERC721.
     * @param handleNonce is the order number of the handle value transfer.
     */
    event HandleValueTransfer(
        bytes32 requestTxHash,
        TokenType tokenType,
        address from,
        address to,
        address tokenAddress,
        uint256 valueOrTokenId,
        uint64 handleNonce
    );

    // start allows the value transfer request.
    function start(bool _status)
        external
        onlyOwner
    {
        isRunning = _status;
    }

    // stop prevent the value transfer request.
    function setCounterPartBridge(address _bridge)
        external
        onlyOwner
    {
        counterpartBridge = _bridge;
    }

    // registerToken can update the allowed token with the counterpart token.
    function registerToken(address _token, address _cToken)
        external
        onlyOwner
    {
        allowedTokens[_token] = _cToken;
    }

    // deregisterToken can remove the token in allowedToken list.
    function deregisterToken(address _token)
        external
        onlyOwner
    {
        delete allowedTokens[_token];
    }

    function updateHandleNonce(uint64 _requestedNonce) internal {
        uint64 i;
        handledNonces[_requestedNonce] = true;

        if (_requestedNonce > maxHandledRequestedNonce) {
            maxHandledRequestedNonce = _requestedNonce;
        }
        for (i = sequentialHandleNonce; i <= maxHandledRequestedNonce && handledNonces[i]; i++) { }
        sequentialHandleNonce = i;
    }

    // handleERC20Transfer sends the token by the request.
    function handleERC20Transfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _value,
        uint64 _requestedNonce,
        uint64 _requestedBlockNumber
    )
        public
        onlyOperators
    {
        bytes32 voteKey = keccak256(abi.encodePacked(VoteType.ValueTransfer, _from, _to, _tokenAddress, _value, _requestedNonce, _requestedBlockNumber));
        if (!voteValueTransfer(voteKey, _requestedNonce)) {
            return;
        }

        emit HandleValueTransfer(_requestTxHash, TokenType.ERC20, _from, _to, _tokenAddress, _value, _requestedNonce);
        _setHandledRequestTxHash(_requestTxHash);
        lastHandledRequestBlockNumber = _requestedBlockNumber;

        updateHandleNonce(_requestedNonce);

        if (modeMintBurn) {
            ERC20Mintable(_tokenAddress).mint(_to, _value);
        } else {
            IERC20(_tokenAddress).transfer(_to, _value);
        }
    }

    // handleKLAYTransfer sends the KLAY by the request.
    function handleKLAYTransfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        uint256 _value,
        uint64 _requestedNonce,
        uint64 _requestedBlockNumber
    )
    public
    onlyOperators
    {
        bytes32 voteKey = keccak256(abi.encodePacked(VoteType.ValueTransfer, _from, _to, _value, _requestedNonce, _requestedBlockNumber));
        if (!voteValueTransfer(voteKey, _requestedNonce)) {
            return;
        }

        emit HandleValueTransfer(_requestTxHash, TokenType.KLAY, _from, _to, address(0), _value, _requestedNonce);
        _setHandledRequestTxHash(_requestTxHash);
        lastHandledRequestBlockNumber = _requestedBlockNumber;

        updateHandleNonce(_requestedNonce);
        _to.transfer(_value);
    }

    // handleERC721Transfer sends the ERC721 by the request.
    function handleERC721Transfer(
        bytes32 _requestTxHash,
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _tokenId,
        uint64 _requestedNonce,
        uint64 _requestedBlockNumber,
        string _tokenURI
    )
        public
        onlyOperators
    {
        bytes32 voteKey = keccak256(abi.encodePacked(VoteType.ValueTransfer, _from, _to, _tokenAddress, _tokenId, _requestedNonce, _requestedBlockNumber, _tokenURI));
        if (!voteValueTransfer(voteKey, _requestedNonce)) {
            return;
        }

        emit HandleValueTransfer(_requestTxHash, TokenType.ERC721, _from, _to, _tokenAddress, _tokenId, _requestedNonce);
        _setHandledRequestTxHash(_requestTxHash);
        lastHandledRequestBlockNumber = _requestedBlockNumber;

        updateHandleNonce(_requestedNonce);

        if (modeMintBurn) {
            ERC721MetadataMintable(_tokenAddress).mintWithTokenURI(_to, _tokenId, _tokenURI);
        } else {
            IERC721(_tokenAddress).safeTransferFrom(address(this), _to, _tokenId);
        }
    }

    // _requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function _requestKLAYTransfer(address _to, uint256 _feeLimit) internal {
        require(isRunning, "stopped bridge");
        require(msg.value > _feeLimit, "insufficient amount");

        uint256 fee = _payKLAYFeeAndRefundChange(_feeLimit);

        emit RequestValueTransfer(
            TokenType.KLAY,
            msg.sender,
            _to,
            address(0),
            msg.value.sub(_feeLimit),
            requestNonce,
            "",
            fee
        );
        requestNonce++;
    }

    // () requests transfer KLAY to msg.sender address on relative chain.
    function () external payable {
        _requestKLAYTransfer(msg.sender, feeOfKLAY);
    }

    // requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function requestKLAYTransfer(address _to, uint256 _value) external payable {
        uint256 feeLimit = msg.value.sub(_value);
        _requestKLAYTransfer(_to, feeLimit);
    }

    // _requestERC20Transfer requests transfer ERC20 to _to on relative chain.
    function _requestERC20Transfer(address _tokenAddress, address _from, address _to, uint256 _value, uint256 _feeLimit) internal {
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
            fee
        );
        requestNonce++;
    }

    // Receiver function of ERC20 token for 1-step deposits to the Bridge
    function onERC20Received(
        address _from,
        uint256 _value,
        address _to,
        uint256 _feeLimit
    )
    public
    {
        _requestERC20Transfer(msg.sender, _from, _to, _value, _feeLimit);
    }

    // requestERC20Transfer requests transfer ERC20 to _to on relative chain.
    function requestERC20Transfer(address _tokenAddress, address _to, uint256 _value, uint256 _feeLimit) external {
        IERC20(_tokenAddress).transferFrom(msg.sender, address(this), _value.add(_feeLimit));
        _requestERC20Transfer(_tokenAddress, msg.sender, _to, _value, _feeLimit);
    }

    // _requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function _requestERC721Transfer(address _tokenAddress, address _from, address _to, uint256 _tokenId) internal {
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
            0
        );
        requestNonce++;
    }

    // Receiver function of ERC721 token for 1-step deposits to the Bridge
    function onERC721Received(
        address _from,
        uint256 _tokenId,
        address _to
    )
    public
    {
        _requestERC721Transfer(msg.sender, _from, _to, _tokenId);
    }

    // requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function requestERC721Transfer(address _tokenAddress, address _to, uint256 _tokenId) external {
        IERC721(_tokenAddress).transferFrom(msg.sender, address(this), _tokenId);
        _requestERC721Transfer(_tokenAddress, msg.sender, _to, _tokenId);
    }

    // chargeWithoutEvent sends KLAY to this contract without event for increasing
    // the withdrawal limit.
    function chargeWithoutEvent() external payable {}

    // setKLAYFee set the fee of KLAY tranfser
    function setKLAYFee(uint256 _fee, uint64 _requestNonce)
        external
        onlyOperators
    {
        bytes32 voteKey = keccak256(abi.encodePacked(this.setKLAYFee.selector, _fee, _requestNonce));
        if (!voteConfiguration(voteKey, _requestNonce)) {
            return;
        }
        _setKLAYFee(_fee);
    }

    // setERC20Fee set the fee of the token transfer
    function setERC20Fee(address _token, uint256 _fee, uint64 _requestNonce)
        external
        onlyOperators
    {
        bytes32 voteKey = keccak256(abi.encodePacked(this.setERC20Fee.selector, _token, _fee, _requestNonce));
        if (!voteConfiguration(voteKey, _requestNonce)) {
            return;
        }
        _setERC20Fee(_token, _fee);
    }

    // setFeeReceiver set fee receiver.
    function setFeeReceiver(address _feeReceiver)
        external
        onlyOwner
    {
        _setFeeReceiver(_feeReceiver);
    }
}
