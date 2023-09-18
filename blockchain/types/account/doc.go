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
Package account implements Account used in Klaytn.

Inside the package, types, structures, functions, and interfaces associated with the Account are defined.

Type of Account

There are three types of Accounts used in Klaytn.

  - LegacyAccountType
  - ExternallyOwnedAccountType
  - SmartContractAccountType

AccountCommon implements the structure and functions common to Klaytn Account, which also implements the Account interface.

EOA (ExternallyOwnedAccount) and SCA (SmartContractAccount) are implemented in a structure that includes AccountCommon.

LegacyAccount is implemented separately according to the account interface.

Source Files

Account related functions and variables are defined in the files listed below.
  - account.go                  : Defines types, interfaces and functions associated with the Account.
  - account_common.go           : Data structures and functions that are common to EOA (ExternallyOwnedAccount) and SCA (SmartContractAccount) are defined as AccountCommon.
  - account_serializer.go       : AccountSerializer is defined for serializing Account.
  - externally_owned_account.go : ExternallyOwnedAccount containing an AccountCommon is defined.
  - legacy_account.go           : LegacyAccount that implements the Account interface is defined.
  - smart_contract_account.go   : SmartContractAccount containing an AccountCommon is defined.

For more information on Account, please refer to the document below.
https://docs.klaytn.com/klaytn/design/accounts#klaytn-accounts
*/
package account
