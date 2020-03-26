// Copyright 2018 The klaytn Authors
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

/*
Package accountkey implements the AccountKey used in Klaytn.
Inside the package, types, functions and interfaces associated with the AccountKey are defined.

Type of AccountKey

The AccountKey types used in Klaytn are as follows:
  - AccountKeyTypeNil
  - AccountKeyTypeLegacy
  - AccountKeyTypePublic
  - AccountKeyTypeFail
  - AccountKeyTypeWeightedMultiSig
  - AccountKeyTypeRoleBased

Each AccountKey type implements the AccountKey interface.

Source Files

AccountKey related functions and variables are defined in the files listed below.
  - account_key.go                    : Defines the AccountKey types, the AccountKey interface and the functions related to AccountKey.
  - account_key_fail.go               : An AccountKey for AccountKeyFail type is defined. If an account has the fail key, the account's transaction validation process always fails.
  - account_key_legacy.go             : An AccountKey for AccountKeyLegacy type is defined. If an account has the legacy key, the account's key pair should be coupled with its address.
  - account_key_nil.go                : An AccountKey for AccountKeyNil type is defined. The nil key is used only for TxTypeAccountUpdate transactions representing an empty key.
  - account_key_public.go             : An AccountKey for AccountKeyPublic type is defined. If an account contains a public key as an account key, the public key will be used in the account's transaction validation process.
  - account_key_role_based.go         : An AccountKey for AccountKeyRoleBased type is defined. AccountKeyRoleBased contains keys that have three roles: RoleTransaction, RoleAccountUpdate, and RoleFeePayer. If an account has a role-based key that consists of more than one key, the account's transaction validation process will use one key in the role-based key depends on the transaction type.
  - account_key_serializer.go         : AccountKeySerializer is defined for serialization of AccountKey.
  - account_key_weighted_multi_sig.go : An AccountKey for AccountKeyWeightedMultiSig type is defined. AccountKeyWeightedMultiSig contains Threshold and WeightedPublicKeys.
  - public_key.go                     : PublicKeySerializable is defined for serialization of public key.

For more information on AccountKey, please see the document below.
https://docs.klaytn.com/klaytn/design/accounts#account-key
*/
package accountkey
