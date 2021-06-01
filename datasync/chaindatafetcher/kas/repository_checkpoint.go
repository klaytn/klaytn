// Copyright 2020 The klaytn Authors
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

package kas

import (
	"errors"

	"github.com/jinzhu/gorm"
)

const checkpointKey = "checkpoint"

func (r *repository) WriteCheckpoint(checkpoint int64) error {
	data := &FetcherMetadata{
		Key:   checkpointKey,
		Value: checkpoint,
	}

	return r.db.Save(data).Error
}

func (r *repository) ReadCheckpoint() (int64, error) {
	checkpoint, err := r.readCheckpoint()
	if r.isRecordNotFoundError(err) {
		return 0, nil
	}
	return checkpoint, err
}

func (r *repository) readCheckpoint() (int64, error) {
	data := &FetcherMetadata{}
	if err := r.db.Where("`key` = ?", checkpointKey).First(data).Error; err != nil {
		return 0, err
	}
	return data.Value, nil
}

func (r *repository) isRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || gorm.IsRecordNotFoundError(err)
}
