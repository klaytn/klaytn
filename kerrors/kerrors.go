// Copyright 2019 The klaytn Authors
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

package kerrors

import "errors"

// TODO-Klaytn: Use integer for error codes.
// TODO-Klaytn: Integrate all universally accessible errors into kerrors package.
var (
	ErrNotHumanReadableAddress    = errors.New("Human-readable address is not supported now")
	ErrHumanReadableNotSupported  = errors.New("Human-readable address is not supported now")
	ErrInvalidContractAddress     = errors.New("contract deploy transaction can't have a recipient address")
	ErrOutOfGas                   = errors.New("out of gas")
	ErrMaxKeysExceed              = errors.New("the number of keys exceeds the limit")
	ErrMaxKeysExceedInValidation  = errors.New("the number of keys exceeds the limit in the validation check")
	ErrMaxFeeRatioExceeded        = errors.New("fee ratio exceeded the maximum")
	ErrEmptySlice                 = errors.New("slice is empty")
	ErrNotForProgramAccount       = errors.New("this type transaction cannot be sent to contract addresses")
	ErrNotProgramAccount          = errors.New("not a program account (e.g., an account having code and storage)")
	ErrPrecompiledContractAddress = errors.New("the address is reserved for pre-compiled contracts")
	ErrInvalidCodeFormat          = errors.New("smart contract code format is invalid")

	// Error codes related to account keys.
	ErrAccountAlreadyExists                 = errors.New("account already exists")
	ErrFeeRatioOutOfRange                   = errors.New("fee ratio is out of range [1, 99]")
	ErrAccountKeyFailNotUpdatable           = errors.New("AccountKeyFail is not updatable")
	ErrDifferentAccountKeyType              = errors.New("different account key type")
	ErrAccountKeyNilUninitializable         = errors.New("AccountKeyNil cannot be initialized to an account")
	ErrNotOnCurve                           = errors.New("public key is not on curve")
	ErrZeroKeyWeight                        = errors.New("key weight is zero")
	ErrUnserializableKey                    = errors.New("key is not serializable")
	ErrDuplicatedKey                        = errors.New("duplicated key")
	ErrWeightedSumOverflow                  = errors.New("weighted sum overflow")
	ErrUnsatisfiableThreshold               = errors.New("unsatisfiable threshold. Weighted sum of keys is less than the threshold.")
	ErrZeroLength                           = errors.New("length is zero")
	ErrLengthTooLong                        = errors.New("length too long")
	ErrNestedCompositeType                  = errors.New("nested composite type")
	ErrLegacyTransactionMustBeWithLegacyKey = errors.New("a legacy transaction must be with a legacy account key")

	ErrDeprecated   = errors.New("deprecated feature")
	ErrNotSupported = errors.New("not supported")
)
