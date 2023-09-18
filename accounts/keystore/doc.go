// Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from accounts/keystore/keystore.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package keystore implements encrypted storage of secp256k1 private keys.

Keys are stored as encrypted JSON files according to the Web3 Secret Storage specification.
See https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition for more information.

Source Files

Each file contains following contents
 - account_cache.go 	: Provides `accountCache` which contains a live index of all accounts in keystore folder
 - file_cache.go 	: Provides `fileCache` which contains information of all files in keystore folder
 - key.go 		: Defines `KeyV3` struct, `keyStore` interface and related functions
 - keyv4.go 		: Defines `KeyV4` struct.
 - keystore.go 		: Defines `KeyStore` which manages a key storage directory on disk and related functions
 - keystore_passphrase.go: Provides functions to encrypt and decrypt `Key` with a passphrase
 - keystore_plain.go 	: Deprecated
 - keystore_wallet.go 	: Defines `keystoreWallet` struct which implements accounts.Wallet interface. Wallet represents a software or hardware wallet that might contain one or more accounts
 - presale.go 		: Deprecated
 - watch.go 		: Provides a watcher which monitors any changes on the keystore folder
 - watch_fallback.go 	: Provides an empty watcher for unsupported platforms
*/
package keystore
