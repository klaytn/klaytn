pragma solidity ^0.4.24;

import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721MetadataMintable.sol";
import "../externals/openzeppelin-solidity/contracts/token/ERC721/ERC721Burnable.sol";

import "../externals/openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./ERC721ServiceChain.sol";


contract ServiceChainNFT is ERC721Full("ServiceChainNFT", "SCN"), ERC721Burnable, ERC721MetadataMintable, ERC721ServiceChain {
    constructor (address _bridge) ERC721ServiceChain(_bridge) public {
    }

    // registerBulk registers (startID, endID-1) NFTs to the user once.
    // This is only for load test.
    function registerBulk(address _user, uint256 _startID, uint256 _endID) onlyOwner external {
        for (uint256 uid = _startID; uid < _endID; uid++) {
            mintWithTokenURI(_user, uid, "testURI");
        }
    }
}
