// Sources flattened with hardhat v2.11.1 https://hardhat.org

// File contracts/IGovParam.sol

pragma solidity ^0.8.0;

/**
 *
 * @dev Interface of the GovParam Contract
 *
 */
interface IGovParam {
    struct ParamState {
        bytes value;
        bool exists;
    }

    struct Param {
        uint64 activation;
        bool votable;
        ParamState prev; // before activation (exclusive)
        ParamState next; // after activation (inclusive)
    }

    event AddParam(string, bytes);
    event SetParam(string, bytes, uint64);
    event DeleteParam(string);
    event SetParamVotable(string, bool);

    function getAllStructParams() external view returns (string[] memory, Param[] memory);

    function getParam(string memory name) external view returns (bytes memory);

    function getAllParams() external view returns (string[] memory, bytes[] memory);

    function addParam(
        string calldata name,
        bytes calldata value
    ) external;

    function setParam(
        string calldata name,
        bytes calldata value,
        uint64 activation
    ) external;

    function deleteParam(
        string calldata name
    ) external;

    function setParamVotable(string memory name, bool votable) external;
}


// File @openzeppelin/contracts/utils/Context.sol@v4.7.3

// OpenZeppelin Contracts v4.4.1 (utils/Context.sol)

pragma solidity ^0.8.0;

/**
 * @dev Provides information about the current execution context, including the
 * sender of the transaction and its data. While these are generally available
 * via msg.sender and msg.data, they should not be accessed in such a direct
 * manner, since when dealing with meta-transactions the account sending and
 * paying for execution may not be the actual sender (as far as an application
 * is concerned).
 *
 * This contract is only required for intermediate, library-like contracts.
 */
abstract contract Context {
    function _msgSender() internal view virtual returns (address) {
        return msg.sender;
    }

    function _msgData() internal view virtual returns (bytes calldata) {
        return msg.data;
    }
}


// File @openzeppelin/contracts/access/Ownable.sol@v4.7.3

// OpenZeppelin Contracts (last updated v4.7.0) (access/Ownable.sol)

pragma solidity ^0.8.0;

/**
 * @dev Contract module which provides a basic access control mechanism, where
 * there is an account (an owner) that can be granted exclusive access to
 * specific functions.
 *
 * By default, the owner account will be the one that deploys the contract. This
 * can later be changed with {transferOwnership}.
 *
 * This module is used through inheritance. It will make available the modifier
 * `onlyOwner`, which can be applied to your functions to restrict their use to
 * the owner.
 */
abstract contract Ownable is Context {
    address private _owner;

    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    /**
     * @dev Initializes the contract setting the deployer as the initial owner.
     */
    constructor() {
        _transferOwnership(_msgSender());
    }

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        _checkOwner();
        _;
    }

    /**
     * @dev Returns the address of the current owner.
     */
    function owner() public view virtual returns (address) {
        return _owner;
    }

    /**
     * @dev Throws if the sender is not the owner.
     */
    function _checkOwner() internal view virtual {
        require(owner() == _msgSender(), "Ownable: caller is not the owner");
    }

    /**
     * @dev Leaves the contract without owner. It will not be possible to call
     * `onlyOwner` functions anymore. Can only be called by the current owner.
     *
     * NOTE: Renouncing ownership will leave the contract without an owner,
     * thereby removing any functionality that is only available to the owner.
     */
    function renounceOwnership() public virtual onlyOwner {
        _transferOwnership(address(0));
    }

    /**
     * @dev Transfers ownership of the contract to a new account (`newOwner`).
     * Can only be called by the current owner.
     */
    function transferOwnership(address newOwner) public virtual onlyOwner {
        require(newOwner != address(0), "Ownable: new owner is the zero address");
        _transferOwnership(newOwner);
    }

    /**
     * @dev Transfers ownership of the contract to a new account (`newOwner`).
     * Internal function without access restriction.
     */
    function _transferOwnership(address newOwner) internal virtual {
        address oldOwner = _owner;
        _owner = newOwner;
        emit OwnershipTransferred(oldOwner, newOwner);
    }
}


// File contracts/GovParam.sol

pragma solidity ^0.8.0;


/**
 *
 * @dev Contract to store and update governance parameters
 * This contract can be called by node to read the param values in the current block
 * Also, the governance contract can change the parameter values.
 *
 */
contract GovParam is Ownable, IGovParam
{
    string[] private _paramNames;
    mapping(string => Param) private _params;

    /**
     * @dev Returns all parameter names including non-existing ones
     */
    function paramNames()
        public
        view
        returns (string[] memory)
    {
        return _paramNames;
    }

    /**
     * @dev Returns all existing parameter names
     */
    function paramExistingNames()
        public
        view
        returns (string[] memory)
    {
        string[] memory src = paramNames();
        uint256 cnt;
        for (uint256 i = 0; i < src.length; i++) {
            if (_exists(src[i])) {
                cnt++;
            }
        }
        string[] memory dst = new string[](cnt);
        cnt = 0;
        for (uint256 i = 0; i < src.length; i++) {
            if (_exists(src[i])) {
                dst[cnt] = src[i];
                cnt++;
            }
        }
        return dst;
    }

    function params(string memory name)
        public
        view
        returns (Param memory)
    {
        return _params[name];
    }

    /**
     * @dev Returns all parameters as struct, including non-existing ones
     */
    function getAllStructParams()
        external
        view
        override
        returns (string[] memory, Param[] memory)
    {
        string[] memory names = paramNames();
        Param[] memory p = new Param[](names.length);
        for (uint256 i = 0; i < names.length; i++) {
            p[i] = params(names[i]);
        }
        return (names, p);
    }

    /**
     * @dev Returns parameter value viewed by the current block
     */
    function getParam(string memory name)
        public
        view
        override
        returns (bytes memory)
    {
        if (!_exists(name)) {
            return "";
        }
        return _state(name).value;
    }

    /**
     * @dev Returns all parameters value viewed by the current block
     */
    function getAllParams()
        external
        view
        override
        returns (string[] memory, bytes[] memory)
    {
        string[] memory names = paramExistingNames();
        bytes[] memory vals = new bytes[](names.length);
        for (uint256 i = 0; i < names.length; i++) {
            vals[i] = _state(names[i]).value;
        }
        return (names, vals);
    }

    /**
     * @dev Adds a new parameter
     *
     * @param name The name of the parameter (e.g., gasLimit)
     * @param value The value of the parameter (e.g., "0xff00")
     */
    function addParam(
        string calldata name,
        bytes calldata value
    ) external override onlyOwner {
        require(bytes(name).length > 0, "GovParam: name cannot be empty");
        require(!_exists(name), "GovParam: parameter already exists");

        _prepareOperation(name);

        ParamState storage next = _params[name].next;
        next.value  = value;
        next.exists = true;
        _params[name].activation = uint64(block.number + 1);
        _paramNames.push(name);

        emit AddParam(name, value);
    }

    /**
     * @dev Sets the value of a parameter and the value changes from the activation block
     *
     * @param name The name of the parameter (e.g., gasLimit)
     * @param value The value of the parameter (e.g., "0xff00")
     * @param activation The activation block number
     */
    function setParam(
        string calldata name,
        bytes calldata value,
        uint64 activation
    ) external override onlyOwner {
        require(bytes(name).length > 0, "GovParam: name cannot be empty");
        require(_exists(name), "GovParam: parameter does not exist");
        require(
            activation > block.number,
            "GovParam: activation must be in a future"
        );

        _prepareOperation(name);

        ParamState storage next = _params[name].next;
        next.value  = value;
        next.exists = true;
        _params[name].activation = activation;

        emit SetParam(name, value, activation);
    }

    /**
     * @dev Delete a parameter
     *
     * @param name The name of the parameter (e.g., gasLimit)
     */
    function deleteParam(
        string calldata name
    ) external override onlyOwner {
        require(bytes(name).length > 0, "GovParam: name cannot be empty");
        require(_exists(name), "GovParam: parameter does not exist");

        _prepareOperation(name);

        ParamState storage next = _params[name].next;
        next.value  = "";
        next.exists = false;
        _params[name].activation = uint64(block.number + 1);

        emit DeleteParam(name);
    }

    /**
     * @dev Make a parameter votable, so a parameter can be changed by governance contract
     *
     * @param name The name of the parameter (e.g., gasLimit)
     * - Requirements: By default, params(name).votable is false
     * @param votable True if votable, false otherwise
     */
    function setParamVotable(
        string memory name,
        bool votable
    ) external override onlyOwner {
        require(bytes(name).length > 0, "GovParam: name cannot be empty");
        require(_exists(name), "GovParam: parameter does not exist");

        _params[name].votable = votable;

        emit SetParamVotable(name, votable);
    }

    /**
     * @dev Migrate to initial state so that operations can be done
     */
    function _prepareOperation(string memory name) internal {
        if (block.number >= params(name).activation) {
            _params[name].prev = _params[name].next;
        }
    }

    /**
     * @dev Returns the current state of param
     */
    function _state(string memory name) internal view returns (ParamState memory) {
        if (block.number >= params(name).activation) {
            return params(name).next;
        }

        return params(name).prev;
    }

    function _exists(string memory name) internal view returns (bool) {
        return _state(name).exists;
    }
}
