pragma solidity ^0.4.24;


contract Callback {
    event RegisteredOffer(
        address owner,
        uint256 valueOrID,
        address tokenAddress,
        uint256 price
    );

    function registerOffer(
        address _owner,
        uint256 _valueOrID,
        address _tokenAddress,
        uint256 _price
    )
        public
    {
        emit RegisteredOffer(_owner, _valueOrID, _tokenAddress, _price);
    }
}
