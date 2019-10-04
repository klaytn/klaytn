pragma solidity 0.4.24;

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

    constructor() public {
        owner = msg.sender;
    }

    function insertIdentity(string _serialNumber, string _publicKey, string _hash)
        public 
    {
        require(bytes(_serialNumber).length > 0);
        require(bytes(_publicKey).length > 0);
        require(bytes(_hash).length > 0);

        identites[_serialNumber] = Identity(_publicKey, _hash);
    }

    function getIdentity(string _serialNumber)
        public
        view
        returns (string, string) 
    {
        require(bytes(_serialNumber).length > 0);

        Identity memory identity = identites[_serialNumber];
        return (identity.publicKey, identity.hash);
    }

    function deleteIdentity(string _serialNumber)
        public
        onlyOwner 
    {
        require(bytes(_serialNumber).length > 0);

        delete identites[_serialNumber];
    }

    function insertCaCertificate(string _caKey, string _caCert)
        public
        onlyOwner 
    {
        require(bytes(_caKey).length > 0);
        require(bytes(_caCert).length > 0);

        caCertificates[_caKey] = _caCert;
    }

    function getCaCertificate(string _caKey)
        public
        view
        returns (string) 
    {
        require(bytes(_caKey).length > 0);

        return caCertificates[_caKey];
    }

    function deleteCaCertificate(string _caKey)
        public
        onlyOwner 
    {
        require(bytes(_caKey).length > 0);

        delete caCertificates[_caKey];
    }
}
