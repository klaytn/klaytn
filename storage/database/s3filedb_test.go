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
//
// For local test, please run the below.
//    $ docker run -d -p 4566:4566 localstack/localstack:0.11.5

package database

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/storage"
	"github.com/stretchr/testify/suite"
)

type SuiteS3FileDB struct {
	suite.Suite
	s3DB           *s3FileDB
	testBucketName *string
}

func (s *SuiteS3FileDB) SetupSuite() {
	// use local test configs
	region := "us-east-1"
	endpoint := "http://localhost:4566"
	testBucketName := aws.String("test-bucket")

	s3DB, err := newS3FileDB(region, endpoint, *testBucketName)
	if err != nil {
		s.Fail("failed to create s3Database", "err", err)
	}

	_, err = s3DB.s3.DeleteBucketPolicy(&s3.DeleteBucketPolicyInput{Bucket: testBucketName})
	if err != nil {
		s.Fail("failed to delete the bucket policy for the test", "err", err, "bucketName", *testBucketName)
	}

	if err != nil {
		s.Fail("failed to create the test bucket", "err", err)
	}

	s.s3DB = s3DB
	s.testBucketName = testBucketName
}

func (s *SuiteS3FileDB) TearDownSuite() {
	if _, err := s.s3DB.s3.DeleteBucket(&s3.DeleteBucketInput{Bucket: s.testBucketName}); err != nil {
		s.Fail("failed to delete the test bucket", "err", err, "bucketName", *s.testBucketName)
	}
}

func TestSuiteS3FileDB(t *testing.T) {
	storage.SkipLocalTest(t)

	suite.Run(t, new(SuiteS3FileDB))
}

func (s *SuiteS3FileDB) TestS3FileDB() {
	testKey := common.MakeRandomBytes(32)
	testVal := common.MakeRandomBytes(1024 * 1024)

	_, err := s.s3DB.read(testKey)
	if err == nil || !strings.Contains(err.Error(), s3.ErrCodeNoSuchKey) {
		s.Fail("test key already exist", "bucketName", *s.testBucketName)
	}

	uris, err := s.s3DB.write(item{key: testKey, val: testVal})
	if err != nil {
		s.Fail("Failed to write", "err", err, "bucketName", *s.testBucketName)
	}
	defer s.s3DB.delete(testKey)

	if uris == "" {
		s.Fail("Unexpected amount of uris are returned", "len(uris)", len(uris))
	}

	val, err := s.s3DB.read(testKey)
	if err != nil {
		s.Fail("Failed to read", "err", err, "bucketName", *s.testBucketName)
	}
	s.True(bytes.Equal(testVal, val))

	if s.s3DB.delete(testKey) != nil {
		s.Fail("Failed to delete", "err", err, "bucketName", *s.testBucketName)
	}
}

func (s *SuiteS3FileDB) TestS3FileDB_Overwrite() {
	testKey := common.MakeRandomBytes(32)
	var testVals [][]byte
	for i := 0; i < 10; i++ {
		testVals = append(testVals, common.MakeRandomBytes(1024*1024))
	}

	_, err := s.s3DB.read(testKey)
	if err == nil || !strings.Contains(err.Error(), s3.ErrCodeNoSuchKey) {
		s.Fail("test key already exist", "bucketName", *s.testBucketName)
	}

	var uris []string
	for _, testVal := range testVals {
		uri, err := s.s3DB.write(item{key: testKey, val: testVal})
		if err != nil {
			s.Fail("failed to write the data to s3FileDB", "err", err)
		}
		uris = append(uris, uri)
	}
	defer s.s3DB.delete(testKey)

	returnedVal, err := s.s3DB.read(testKey)
	if err != nil {
		s.Fail("failed to read the data from s3FileDB", "err", err)
	}

	s.Equal(testVals[len(testVals)-1], returnedVal)
	s.Equal(len(testVals), len(uris))
}

func (s *SuiteS3FileDB) TestS3FileDB_EmptyDelete() {
	testKey := common.MakeRandomBytes(256)
	s.NoError(s.s3DB.delete(testKey))
	s.NoError(s.s3DB.delete(testKey))
}
