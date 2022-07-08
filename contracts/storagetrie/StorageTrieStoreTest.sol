pragma solidity >0.8.0;

contract StorageTrieStoreTest {
    struct Identity {
        string publicKey;
        string hash;
    }

    address public owner; // ex : 0xe5e4a8f4ecc2b6298be33fed07a09599db4e46fa

    string public rootCaCertificate; // ROOT_CA
    mapping(string => string) caCertificates; // CA_C1

    mapping(string => Identity) identites;

    modifier onlyOwner {
        require(msg.sender == owner);
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    function insertIdentity(string calldata _serialNumber, string calldata _publicKey, string calldata _hash)
        public 
    {
        require(bytes(_serialNumber).length > 0);
        require(bytes(_publicKey).length > 0);
        require(bytes(_hash).length > 0);

        identites[_serialNumber] = Identity(_publicKey, _hash);
    }

    function getIdentity(string calldata _serialNumber)
        public
        view
        returns (string memory, string memory) 
    {
        require(bytes(_serialNumber).length > 0);

        Identity memory identity = identites[_serialNumber];
        return (identity.publicKey, identity.hash);
    }

    function deleteIdentity(string calldata _serialNumber)
        public
        onlyOwner 
    {
        require(bytes(_serialNumber).length > 0);

        delete identites[_serialNumber];
    }

    function insertCaCertificate(string calldata _caKey, string calldata _caCert)
        public
        onlyOwner 
    {
        require(bytes(_caKey).length > 0);
        require(bytes(_caCert).length > 0);

        caCertificates[_caKey] = _caCert;
    }

    function getCaCertificate(string calldata _caKey)
        public
        view
        returns (string memory) 
    {
        require(bytes(_caKey).length > 0);

        return caCertificates[_caKey];
    }

    function deleteCaCertificate(string calldata _caKey)
        public
        onlyOwner 
    {
        require(bytes(_caKey).length > 0);

        delete caCertificates[_caKey];
    }
}
