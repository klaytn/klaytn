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

import "strings"

func (s *SuiteRepository) TestCheckpoint_Success() {
	// Clean up the table before test.
	s.NoError(s.repo.db.Delete(FetcherMetadata{}).Error)

	expected := int64(1912)
	err := s.repo.WriteCheckpoint(expected)
	s.NoError(err)

	actual, err := s.repo.ReadCheckpoint()
	s.NoError(err)
	s.Equal(expected, actual)
}

func (s *SuiteRepository) TestCheckpoint_Fail_RecordNotFound() {
	// Clean up the table before test.
	s.NoError(s.repo.db.Delete(FetcherMetadata{}).Error)

	// readCheckpoint returns an error if it is failed to read a checkpoint from database.
	_, err := s.repo.readCheckpoint()
	s.Error(err)
	s.True(strings.Contains(err.Error(), "record not found"))
}

func (s *SuiteRepository) TestCheckpoint_Success_RecordNotFound() {
	// Clean up the table before test.
	s.NoError(s.repo.db.Delete(FetcherMetadata{}).Error)

	// ReadCheckpoint filters "record not found" error and returns 0.
	checkpoint, err := s.repo.ReadCheckpoint()
	s.NoError(err)
	s.Equal(int64(0), checkpoint)
}
