// Copyright 2023 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @title KIP-113 BLS public key registry
/// @dev See https://github.com/klaytn/kips/issues/113
abstract contract IKIP113 {

    struct BlsPublicKeyInfo {
        /// @dev compressed BLS12-381 public key (48 bytes)
        bytes publicKey;
        /// @dev proof-of-possession (96 bytes)
        ///  must be a result of PopProve algorithm as per
        ///  draft-irtf-cfrg-bls-signature-05 section 3.3.3.
        ///  with ciphersuite "BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_"
        bytes pop;
    }

    mapping(address => BlsPublicKeyInfo) infos;
    address[] addrs;

    event PubkeyRegistered(address addr, bytes publicKey, bytes pop);
    event PubkeyUnregistered(address addr, bytes publicKey, bytes pop);

    /// @dev Registers the given public key associated with the `msg.sender` address.
    ///  The function validates the following requirements:
    ///  - The function MUST revert if `publicKey.length != 48` or `pop.length != 96`.
    ///  - The function MUST revert if `publicKey` or `pop` is equal to zero.
    ///  - The function SHOULD authenticate and authorize `msg.sender`.
    ///  The function emits a `PubkeyRegistered` event.
    ///  _Note_ The function is not able to verify the validity of the public key and the proof-of-possession due to the lack of [EIP-2537](https://eips.ethereum.org/EIPS/eip-2537). See [validation](https://kips.klaytn.foundation/KIPs/kip-113#validation) for off-chain validation.
    function registerPublicKey(bytes calldata publicKey, bytes calldata pop) external virtual;

    /// @dev Unregisters the public key associated with the given `addr` address.
    ///  The function validates the following requirements:
    ///  - The function MUST revert if `addr` has not been registered.
    ///  - The function SHOULD authenticate and authorize `msg.sender`.
    ///  The function emits a `PubkeyUnregistered` event.
    function unregisterPublicKey(address addr) external virtual;

    /// @dev Returns the public key and proof-of-possession given a `addr`.
    ///  The function MUST revert if given `addr` is not registered.
    ///  _Note_ The function is not able to verify the validity of the public key and the proof-of-possession due to the lack of [EIP-2537](https://eips.ethereum.org/EIPS/eip-2537). See [validation](https://kips.klaytn.foundation/KIPs/kip-113#validation) for off-chain validation.
    function getInfo(address addr) external virtual view returns (BlsPublicKeyInfo memory pubkey);

    /// @dev Returns all the stored addresses, public keys, and proof-of-possessions at once.
    ///  _Note_ The function is not able to verify the validity of the public key and the proof-of-possession due to the lack of [EIP-2537](https://eips.ethereum.org/EIPS/eip-2537). See [validation](https://kips.klaytn.foundation/KIPs/kip-113#validation) for off-chain validation.
    function getAllInfo() external virtual view returns (address[] memory addrList, BlsPublicKeyInfo[] memory pubkeyList);
}
