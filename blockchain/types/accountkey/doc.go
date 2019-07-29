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
  - account_key_fail.go               : An AccountKey for fail key is defined. If an account has the fail key, the account's transaction validation process always fails.
  - account_key_legacy.go             : An AccountKey for legacy key is defined.
  - account_key_nil.go                : An AccountKey for nil key is defined.
  - account_key_public.go             : An AccountKey for public key is defined.
  - account_key_role_based.go         : An AccountKey for role based key is defined.
  - account_key_serializer.go         : AccountKeySerializer is defined for serialization of AccountKey.
  - account_key_weighted_multi_sig.go : An AccountKey for weighted multi sig key is defined.
  - public_key.go                     : PublicKeySerializable is defined for serialization of public key.


For more information on AccountKey, please see the document below.
https://docs.klaytn.com/klaytn/design/account#account-key
*/
package accountkey
