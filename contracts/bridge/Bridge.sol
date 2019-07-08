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

import "../servicechain_nft/INFTReceiver.sol";
import "../servicechain_token/ITokenReceiver.sol";

contract Bridge is ITokenReceiver, INFTReceiver, Ownable {
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

    constructor (bool _modeMintBurn) public payable {
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
        string uri);

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
        string tokenURI
    )
        external
        onlyOwner
    {
        require(handleNonce == _requestNonce, "mismatched handle / request nonce");

        emit HandleValueTransfer(_to, TokenKind.ERC721, _contractAddress, _uid, handleNonce);
        lastHandledRequestBlockNumber = _requestBlockNumber;
        handleNonce++;

        if (modeMintBurn) {
            ERC721MetadataMintable(_contractAddress).mintWithTokenURI(_to, _uid, tokenURI);
        } else {
            IERC721(_contractAddress).safeTransferFrom(address(this), _to, _uid);
        }
    }

    bytes4 constant TOKEN_RECEIVED = 0xbc04f0af;

    function onTokenReceived(
        address _from,
        uint256 _amount,
        address _to
    )
        public
        returns (bytes4)
    {
        require(isRunning, "stopped bridge");
        require(allowedTokens[msg.sender] != address(0), "Not a valid token");
        require(_amount > 0, "zero amount");

        emit RequestValueTransfer(TokenKind.ERC20, _from, _amount, msg.sender, _to, requestNonce,"");
        requestNonce++;
        return TOKEN_RECEIVED;
    }

    // Receiver function of NFT for 1-step deposits to the Bridge
    bytes4 private constant ERC721_RECEIVED = 0x150b7a02;

    function onNFTReceived(
        address from,
        uint256 tokenId,
        address to
    )
        public
        returns(bytes4)
    {
        require(isRunning, "stopped bridge");
        require(allowedTokens[msg.sender] != address(0), "Not a valid token");

        emit RequestValueTransfer(TokenKind.ERC721, from, tokenId, msg.sender, to, requestNonce,"");
        requestNonce++;
        return ERC721_RECEIVED;
    }

    // () requests transfer KLAY to msg.sender address on relative chain.
    function () external payable {
        require(isRunning, "stopped bridge");
        require(msg.value > 0, "zero msg.value");

        emit RequestValueTransfer(
            TokenKind.KLAY,
            msg.sender,
            msg.value,
            address(0),
            msg.sender,
            requestNonce,
            ""
        );
        requestNonce++;
    }

    // requestKLAYTransfer requests transfer KLAY to _to on relative chain.
    function requestKLAYTransfer(address _to) external payable {
        require(isRunning, "stopped bridge");
        require(msg.value > 0, "zero msg.value");

        emit RequestValueTransfer(
            TokenKind.KLAY,
            msg.sender,
            msg.value,
            address(0),
            _to,
            requestNonce,
            ""
        );
        requestNonce++;
    }

    // requestERC20Transfer requests transfer ERC20 to _to on relative chain.
    function requestERC20Transfer(address _contractAddress, address _to, uint256 _amount) external {
        require(isRunning, "stopped bridge");
        require(msg.value > 0, "zero msg.value");

        IERC20(_contractAddress).transferFrom(msg.sender, address(this), _amount);

        if (modeMintBurn) {
            ERC20Burnable(_contractAddress).burn(_amount);
        }

        emit RequestValueTransfer(
            TokenKind.KLAY,
            msg.sender,
            msg.value,
            address(0),
            _to,
            requestNonce,
            ""
        );
        requestNonce++;
    }

    // requestERC721Transfer requests transfer ERC721 to _to on relative chain.
    function requestERC721Transfer(address _contractAddress, address _to, uint256 _uid) external {
        require(isRunning, "stopped bridge");
        require(msg.value > 0, "zero msg.value");

        IERC721(_contractAddress).transferFrom(msg.sender, address(this), _uid);

        string memory uri = ERC721Metadata(_contractAddress).tokenURI(_uid);

        if (modeMintBurn) {
            ERC721Burnable(_contractAddress).burn(_uid);
        }

        emit RequestValueTransfer(
            TokenKind.KLAY,
            msg.sender,
            msg.value,
            address(0),
            _to,
            requestNonce,
            uri
        );
        requestNonce++;
    }

    // chargeWithoutEvent sends KLAY to this contract without event for increasing
    // the withdrawal limit.
    function chargeWithoutEvent() external payable {}
}
