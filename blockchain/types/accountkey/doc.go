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
The package accountkey implements the AccountKey used in Klaytn.
Inside the package are the types, functions and interfaces associated with the AccountKey.

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
  - account_key.go                    : Defines the AccountKey type and the AccountKey interface, and defines the functions related to AccountKey.
  - account_key_fail.go               : AccountKeyFail is defined according to the AccountKey interface.
  - account_key_legacy.go             : AccountKeyLegacy is defined according to the AccountKey interface.
  - account_key_nil.go                : AccountKeyNil is defined according to the AccountKey interface.
  - account_key_public.go             : AccountKeyPublic is defined according to the AccountKey interface.
  - account_key_role_based.go         : AccountKeyRoleBased is defined according to the AccountKey interface.
  - account_key_serializer.go         : AccountKeySerializer is defined for serialization of AccountKey.
  - account_key_weighted_multi_sig.go : AccountKeyWeightedMultiSig is defined according to the AccountKey interface.
  - public_key.go                     : PublicKeySerializable is defined for serialization of PublicKey.


For more information on AccountKey, please see the document below.
https://docs.klaytn.com/klaytn/design/account#account-key
*/
package accountkey
