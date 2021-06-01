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

package blockchain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
)

func TestGetVMerrFromReceiptStatus(t *testing.T) {
	err := GetVMerrFromReceiptStatus(types.ReceiptStatusFailed)
	expectedErr := ErrInvalidReceiptStatus
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Invalid err, want %s, got %s", expectedErr, err)
	}

	// Invalid ReceiptStatus
	err = GetVMerrFromReceiptStatus(types.ReceiptStatusLast)
	expectedErr = ErrInvalidReceiptStatus
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Invalid err, want %s, got %s", expectedErr, err)
	}

	err = GetVMerrFromReceiptStatus(types.ReceiptStatusSuccessful)
	if err != nil {
		t.Fatalf("Invalid err, want nil, got %s", err)
	}

	err = GetVMerrFromReceiptStatus(types.ReceiptStatusErrDefault)
	expectedErr = ErrVMDefault
	if err.Error() != expectedErr.Error() {
		t.Fatalf("Invalid err, want %s, got %s", expectedErr, err)
	}
}

func TestGetReceiptStatusFromVMerr(t *testing.T) {
	status := getReceiptStatusFromErrTxFailed(nil)
	expectedStatus := types.ReceiptStatusSuccessful
	if status != expectedStatus {
		t.Fatalf("Invalid receipt status, want %d, got %d", expectedStatus, status)
	}

	status = getReceiptStatusFromErrTxFailed(vm.ErrMaxCodeSizeExceeded)
	expectedStatus = types.ReceiptStatuserrMaxCodeSizeExceed
	if status != expectedStatus {
		t.Fatalf("Invalid receipt status, want %d, got %d", expectedStatus, status)
	}

	// Unknown VM error
	status = getReceiptStatusFromErrTxFailed(errors.New("Unknown VM error"))
	expectedStatus = types.ReceiptStatusErrDefault
	if status != expectedStatus {
		t.Fatalf("Invalid receipt status, want %d, got %d", expectedStatus, status)
	}
}

// TestPrintErrorCodeTable prints the error code table in a format of a markdown table.
func TestPrintErrorCodeTable(t *testing.T) {
	if testing.Verbose() {
		fmt.Println("| ErrorCode | Description |")
		fmt.Println("|---|---|")
		for i := uint(types.ReceiptStatusErrDefault); i < types.ReceiptStatusLast; i++ {
			fmt.Printf("|0x%02x|%s|\n", i, receiptstatus2errTxFailed[i])
		}
	}
}
