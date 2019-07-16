pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/math/SafeMath.sol";

import "../externals/openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";

import "../externals/openzeppelin-solidity/contracts/token/ERC721/IERC721.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721MetadataMintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Burnable.sol";

import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";

import "../sc_erc721/IERC721BridgeReceiver.sol";
import "../sc_erc20/IERC20BridgeReceiver.sol";
import "./BridgeFee.sol";

contract Bridge is IERC20BridgeReceiver, IERC721BridgeReceiver, Ownable, BridgeFee {
    uint64 public constant VERSION = 1;
    bool public modeMintBurn = false;
    address public counterpartBridge;
    bool public isRunning;

    mapping (address => address) public allowedTokens; // <token, counterpart token>

    using SafeMath for uint256;

    uint64 public requestNonce;
    uint64 public handleNonce;

    uint64 public lastHandledRequestBlockNumber;

    enum TokenKind {
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
     * Event to log the withdrawal of a token from the Bridge.
     * @param kind The type of token withdrawn (KLAY/TOKEN/NFT).
     * @param from is the requester of the request value transfer event.
     * @param contractAddress Address of token contract the token belong to.
     * @param amount is the amount for KLAY/TOKEN and the NFT ID for NFT.
     * @param requestNonce is the order number of the request value transfer.
     * @param uri is uri of ERC721 token.
     */
    event RequestValueTransfer(TokenKind kind,
        address from,
        uint256 amount,
        address contractAddress,
        address to,
        uint64 requestNonce,
        string uri,
        uint256 fee
    );

    /**
     * Event to log the withdrawal of a token from the Bridge.
     * @param owner Address of the entity that made the withdrawal.ga
     * @param kind The type of token withdrawn (KLAY/TOKEN/NFT).
     * @param contractAddress Address of token contract the token belong to.
     * @param value For KLAY/TOKEN this is the amount.
     * @param handleNonce is the order number of the handle value transfer.
     */
    event HandleValueTransfer(
        address owner,
        TokenKind kind,
        address contractAddress,
        uint256 value,
        uint64 handleNonce);

    // start allows the value transfer request.
    function start() external onlyOwner {
        isRunning = true;
    }

    // stop prevent the value transfer request.
    function stop() external onlyOwner {
        isRunning = false;
    }

    // stop prevent the value transfer request.
    function setCounterPartBridge(address _bridge) external onlyOwner {
        counterpartBridge = _bridge;
    }

    // registerToken can update the allowed token with the counterpart token.
    function registerToken(address _token, address _cToken) external onlyOwner {
        allowedTokens[_token] = _cToken;
    }

    // deregisterToken can remove the token in allowedToken list.
    function deregisterToken(address _token) external onlyOwner {
        delete allowedTokens[_token];
    }

    // handleERC20Transfer sends the token by the request.
    function handleERC20Transfer(
        uint256 _amount,
        address _to,
        address _contractAddress,
        uint64 _requestNonce,
        uint64 _requestBlockNumber
    )
        external
        onlyOwner
    {
        require(handleNonce == _requestNonce, "mismatched handle / request nonce");

        emit HandleValueTransfer(_to, TokenKind.ERC20, _contractAddress, _amount, handleNonce);
        lastHandledRequestBlockNumber = _requestBlockNumber;
        handleNonce++;

        if (modeMintBurn) {
            ERC20Mintable(_contractAddress).mint(_to, _amount);
        } else {
            IERC20(_contractAddress).transfer(_to, _amount);
        }
    }

    // handleKLAYTransfer sends the KLAY by the request.
    function handleKLAYTransfer(
        uint256 _amount,
        address _to,
        uint64 _requestNonce,
        uint64 _requestBlockNumber
    )
        external
        onlyOwner
    {
        require(handleNonce == _requestNonce, "mismatched handle / request nonce");

        emit HandleValueTransfer(_to, TokenKind.KLAY, address(0), _amount, handleNonce);
        lastHandledRequestBlockNumber = _requestBlockNumber;
        handleNonce++;

        _to.transfer(_amount);
    }

    // handleERC721Transfer sends the NFT by the request.
    function handleERC721Transfer(
        uint256 _uid,
        address _to,
        address _contractAddress,
        uint64 _requestNonce,
        uint64 _requestBlockNumber,
        string _tokenURI
    )
        external
        onlyOwner
    {
        require(handleNonce == _requestNonce, "mismatched handle / request nonce");

        emit HandleValueTransfer(_to, TokenKind.ERC721, _contractAddress, _uid, handleNonce);
        lastHandledRequestBlockNumber = _requestBlockNumber;
        handleNonce++;

        if (modeMintBurn) {
            ERC721MetadataMintable(_contractAddress).mintWithTokenURI(_to, _uid, _tokenURI);
        } else {
            IERC721(_contractAddress).safeTransferFrom(address(this), _to, _uid);
        }
    }

    // _requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function _requestKLAYTransfer(address _to, uint256 _fee) internal {
        require(isRunning, "stopped bridge");
        require(msg.value > 0, "zero msg.value");
        require(msg.value > _fee, "insufficient amount");
        require(feeOfKLAY == _fee, "invalid fee");

        _payKLAYFee(_fee);

        emit RequestValueTransfer(
            TokenKind.KLAY,
            msg.sender,
            msg.value.sub(_fee),
            address(0),
            _to,
            requestNonce,
            "",
            _fee
        );
        requestNonce++;
    }

    // () requests transfer KLAY to msg.sender address on relative chain.
    function () external payable {
        _requestKLAYTransfer(msg.sender, feeOfKLAY);
    }

    // requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function requestKLAYTransfer(address _to, uint256 _fee) external payable {
        _requestKLAYTransfer(_to, _fee);
    }

    // _requestERC20Transfer requests transfer ERC20 to _to on relative chain.
    function _requestERC20Transfer(address _contractAddress, address _from, address _to, uint256 _amount, uint256 _fee) internal {
        require(isRunning, "stopped bridge");
        require(_amount > 0, "zero msg.value");
        require(allowedTokens[_contractAddress] != address(0), "Not a valid token");
        require(_fee == feeOfERC20[_contractAddress], "invalid fee");

        _payERC20Fee(_contractAddress, _fee);

        if (modeMintBurn) {
            ERC20Burnable(_contractAddress).burn(_amount);
        }

        emit RequestValueTransfer(
            TokenKind.ERC20,
            _from,
            _amount,
            _contractAddress,
            _to,
            requestNonce,
            "",
            _fee
        );
        requestNonce++;
    }

    // Receiver function of ERC20 token for 1-step deposits to the Bridge
    function onERC20Received(
        address _from,
        uint256 _amount,
        address _to,
        uint256 _fee
    )
    public
    {
        _requestERC20Transfer(msg.sender, _from, _to, _amount, _fee);
    }

    // requestERC20Transfer requests transfer ERC20 to _to on relative chain.
    function requestERC20Transfer(address _contractAddress, address _to, uint256 _amount, uint256 _fee) external {
        IERC20(_contractAddress).transferFrom(msg.sender, address(this), _amount.add(_fee));
        _requestERC20Transfer(_contractAddress, msg.sender, _to, _amount, _fee);
    }

    // _requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function _requestERC721Transfer(address _contractAddress, address _from, address _to, uint256 _uid) internal {
        require(isRunning, "stopped bridge");
        require(allowedTokens[_contractAddress] != address(0), "Not a valid token");

        string memory uri = ERC721Metadata(_contractAddress).tokenURI(_uid);

        if (modeMintBurn) {
            ERC721Burnable(_contractAddress).burn(_uid);
        }

        emit RequestValueTransfer(
            TokenKind.ERC721,
            _from,
            _uid,
            _contractAddress,
            _to,
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
    function requestERC721Transfer(address _contractAddress, address _to, uint256 _uid) external {
        IERC721(_contractAddress).transferFrom(msg.sender, address(this), _uid);
        _requestERC721Transfer(_contractAddress, msg.sender, _to, _uid);
    }

    // chargeWithoutEvent sends KLAY to this contract without event for increasing
    // the withdrawal limit.
    function chargeWithoutEvent() external payable {}

    // setKLAYFee set the fee of KLAY tranfser
    function setKLAYFee(uint256 _fee) external onlyOwner {
        _setKLAYFee(_fee);
    }

    // setERC20Fee set the fee of the token transfer
    function setERC20Fee(address _token, uint256 _fee) external onlyOwner {
        _setERC20Fee(_token, _fee);
    }

    // setFeeReceiver set fee receiver.
    function setFeeReceiver(address _to) external onlyOwner {
        _setFeeReceiver(_to);
    }
}
