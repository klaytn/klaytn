// Modifications Copyright 2021 The klaytn Authors
// Copyright 2019 The go-ethereum Authors
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
// This file is derived from core/state/snapshot/wipe.go (2021/10/21).
// Modified and improved for the klaytn development.

package snapshot

import (
	"github.com/klaytn/klaytn/storage/database"
	"github.com/rcrowley/go-metrics"
)

// wipeKeyRange deletes a range of keys from the database starting with prefix
// and having a specific total key length. The start and limit is optional for
// specifying a particular key range for deletion.
//
// Origin is included for wiping and limit is excluded if they are specified.
func wipeKeyRange(db database.DBManager, kind string, prefix []byte, origin []byte, limit []byte, keylen int, meter metrics.Meter, report bool) error {
	// TODO-Klaytn-Snapshot port wipeKeyRange
	return nil
}
