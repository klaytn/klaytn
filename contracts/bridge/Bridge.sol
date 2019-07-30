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
import "./BridgeMultiSig.sol";

contract Bridge is IERC20BridgeReceiver, IERC721BridgeReceiver, BridgeFee, BridgeMultiSig {
    uint64 public constant VERSION = 1;
    bool public modeMintBurn = false;
    address public counterpartBridge;
    bool public isRunning;

    uint64 public requestNonce;
    uint64 public lastHandledRequestBlockNumber;
    uint64 public sequentialHandledNonce;
    uint64 public maxHandledNonce;
    mapping (uint64 => bool) public handledNonces;  // <handled nonce> history

    mapping (address => address) public allowedTokens; // <token, counterpart token>

    using SafeMath for uint256;

    enum TokenType {
        KLAY,
        ERC20,
        ERC721
    }

    // TODO-Klaytn-Service FeeReceiver should be passed by argument of constructor.
    constructor (bool _modeMintBurn) BridgeFee(address(0)) BridgeMultiSig() public payable {
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
     * @param tokenType is the type of tokens (KLAY/ERC20/ERC721).
     * @param from is an address of the account who requested the value transfer.
     * @param to is an address of the account who will received the value.
     * @param tokenAddress Address of token contract the token belong to.
     * @param valueOrTokenId is the value of KLAY/ERC20 or token ID of ERC721.
     * @param handleNonce is the order number of the handle value transfer.
     */
    event HandleValueTransfer(
        TokenType tokenType,
        address from,
        address to,
        address tokenAddress,
        uint256 valueOrTokenId,
        uint64 handleNonce
    );

    // start allows the value transfer request.
    function start(bool _status, uint64 _requestNonce)
        external
        onlySigners
    {
        onlySequentialNonce(TransactionType.Configuration, _requestNonce);
        bytes32 voteKey = keccak256(abi.encodePacked(this.start.selector, _status, _requestNonce));
        if (!voteConfiguration(voteKey, msg.sender)) {
            return;
        }
        isRunning = _status;
    }

    // stop prevent the value transfer request.
    function setCounterPartBridge(address _bridge, uint64 _requestNonce)
        external
        onlySigners
    {
        onlySequentialNonce(TransactionType.Configuration, _requestNonce);
        bytes32 voteKey = keccak256(abi.encodePacked(this.setCounterPartBridge.selector, _bridge, _requestNonce));
        if (!voteConfiguration(voteKey, msg.sender)) {
            return;
        }
        counterpartBridge = _bridge;
    }

    // registerToken can update the allowed token with the counterpart token.
    function registerToken(address _token, address _cToken, uint64 _requestNonce)
        external
        onlySigners
    {
        onlySequentialNonce(TransactionType.Configuration, _requestNonce);
        bytes32 voteKey = keccak256(abi.encodePacked(this.registerToken.selector, _token, _cToken, _requestNonce));
        if (!voteConfiguration(voteKey, msg.sender)) {
            return;
        }
        allowedTokens[_token] = _cToken;
    }

    // deregisterToken can remove the token in allowedToken list.
    function deregisterToken(address _token, uint64 _requestNonce)
        external
        onlySigners
    {
        onlySequentialNonce(TransactionType.Configuration, _requestNonce);
        bytes32 voteKey = keccak256(abi.encodePacked(this.deregisterToken.selector, _token, _requestNonce));
        if (!voteConfiguration(voteKey, msg.sender)) {
            return;
        }
        delete allowedTokens[_token];
    }

    // registerSigner registers new signer.
    function registerSigner(address _signer, uint64 _requestNonce)
    external
    onlySigners
    {
        onlySequentialNonce(TransactionType.Configuration, _requestNonce);
        require(_signer != address(0));
        bytes32 voteKey = keccak256(abi.encodePacked(this.registerSigner.selector, _signer, _requestNonce));
        if (!voteConfiguration(voteKey, msg.sender)) {
            return;
        }
        signers[_signer] = true;
    }

    // deregisterSigner deregisters the signer.
    function deregisterSigner(address _signer, uint64 _requestNonce)
    external
    onlySigners
    {
        onlySequentialNonce(TransactionType.Configuration, _requestNonce);
        require(_signer != address(0));
        bytes32 voteKey = keccak256(abi.encodePacked(this.registerSigner.selector, _signer, _requestNonce));
        if (!voteConfiguration(voteKey, msg.sender)) {
            return;
        }
        delete signers[_signer];
    }

    // setSignerThreshold sets signer threshold.
    function setSignerThreshold(TransactionType _txType, uint64 _threshold, uint64 _requestNonce)
    external
    onlySigners
    {
        onlySequentialNonce(TransactionType.Configuration, _requestNonce);
        require(_threshold > 0);
        bytes32 voteKey = keccak256(abi.encodePacked(this.setSignerThreshold.selector, _txType, _threshold, _requestNonce));
        if (!voteConfiguration(voteKey, msg.sender)) {
            return;
        }
        signerThresholds[uint64(_txType)] = _threshold;
    }

    function updateNonce(uint64 _requestNonce) internal {
        if (_requestNonce > maxHandledNonce) {
            maxHandledNonce = _requestNonce;
        }
        // TODO-Klaytn-ServiceChain: optimize this loop if possible.
        for (uint64 i = sequentialHandledNonce; i <= maxHandledNonce; i++) {
            if (!handledNonces[i]) {
                break;
            }
            sequentialHandledNonce = i+1;
        }
    }

    // handleKLAYTransfer sends the KLAY by the request.
    function handleKLAYTransfer(
        address _from,
        address _to,
        uint256 _value,
        uint64 _requestNonce,
        uint64 _requestBlockNumber
    )
        public
        onlySigners
    {
        bytes32 txKey = keccak256(abi.encodePacked(TransactionType.ValueTransfer, _requestNonce));
        bytes32 voteKey = keccak256(abi.encodePacked(TransactionType.ValueTransfer, _from, _to, _value, _requestNonce, _requestBlockNumber));
        if (!voteValueTransfer(txKey, voteKey, msg.sender)) {
            return;
        }

        emit HandleValueTransfer(TokenType.KLAY, _from, _to, address(0), _value, _requestNonce);
        lastHandledRequestBlockNumber = _requestBlockNumber;
        handledNonces[_requestNonce] = true;

        updateNonce(_requestNonce);
        _to.transfer(_value);
    }

    // handleERC20Transfer sends the token by the request.
    function handleERC20Transfer(
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _value,
        uint64 _requestNonce,
        uint64 _requestBlockNumber
    )
        public
        onlySigners
    {
        bytes32 txKey = keccak256(abi.encodePacked(TransactionType.ValueTransfer, _requestNonce));
        bytes32 voteKey = keccak256(abi.encodePacked(TransactionType.ValueTransfer, _from, _to, _tokenAddress, _value, _requestNonce, _requestBlockNumber));
        if (!voteValueTransfer(txKey, voteKey, msg.sender)) {
            return;
        }

        emit HandleValueTransfer(TokenType.ERC20, _from, _to, _tokenAddress, _value, _requestNonce);
        lastHandledRequestBlockNumber = _requestBlockNumber;
        handledNonces[_requestNonce] = true;

        updateNonce(_requestNonce);

        if (modeMintBurn) {
            ERC20Mintable(_tokenAddress).mint(_to, _value);
        } else {
            IERC20(_tokenAddress).transfer(_to, _value);
        }
    }

    // handleERC721Transfer sends the ERC721 by the request.
    function handleERC721Transfer(
        address _from,
        address _to,
        address _tokenAddress,
        uint256 _tokenId,
        uint64 _requestNonce,
        uint64 _requestBlockNumber,
        string _tokenURI
    )
        public
        onlySigners
    {
        bytes32 txKey = keccak256(abi.encodePacked(TransactionType.ValueTransfer, _requestNonce));
        bytes32 voteKey = keccak256(abi.encodePacked(TransactionType.ValueTransfer, _from, _to, _tokenAddress, _tokenId, _requestNonce, _requestBlockNumber, _tokenURI));
        if (!voteValueTransfer(txKey, voteKey, msg.sender)) {
            return;
        }

        emit HandleValueTransfer(TokenType.ERC721, _from, _to, _tokenAddress, _tokenId, _requestNonce);
        lastHandledRequestBlockNumber = _requestBlockNumber;
        handledNonces[_requestNonce] = true;

        updateNonce(_requestNonce);

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
        onlySigners {
        onlySequentialNonce(TransactionType.ConfigurationRealtime, _requestNonce);
        bytes32 voteKey = keccak256(abi.encodePacked(this.setKLAYFee.selector, _fee, _requestNonce));
        if (!voteConfigurationRealtime(voteKey, msg.sender)) {
            return;
        }
        _setKLAYFee(_fee);
    }

    // setERC20Fee set the fee of the token transfer
    function setERC20Fee(address _token, uint256 _fee, uint64 _requestNonce)
        external
        onlySigners
    {
        onlySequentialNonce(TransactionType.ConfigurationRealtime, _requestNonce);
        bytes32 voteKey = keccak256(abi.encodePacked(this.setERC20Fee.selector, _token, _fee, _requestNonce));
        if (!voteConfigurationRealtime(voteKey, msg.sender)) {
            return;
        }
        _setERC20Fee(_token, _fee);
    }

    // setFeeReceiver set fee receiver.
    function setFeeReceiver(address _feeReceiver, uint64 _requestNonce)
        external
        onlySigners
    {
        onlySequentialNonce(TransactionType.Configuration, _requestNonce);
        bytes32 voteKey = keccak256(abi.encodePacked(this.setFeeReceiver.selector, _feeReceiver, _requestNonce));
        if (!voteConfiguration(voteKey, msg.sender)) {
            return;
        }
        _setFeeReceiver(_feeReceiver);
    }
}
