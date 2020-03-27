package database

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

var testRand *rand.Rand

type SuiteS3FileDB struct {
	suite.Suite
	s3DB           *s3FileDB
	testBucketName *string
}

func (s *SuiteS3FileDB) SetupSuite() {
	region := "ap-northeast-2"
	endpoint := "http://localhost:4572"
	testBucketName := aws.String("kas-test-bucket-" + strconv.Itoa(time.Now().Nanosecond()))

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
	testRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (s *SuiteS3FileDB) TearDownSuite() {
	if _, err := s.s3DB.s3.DeleteBucket(&s3.DeleteBucketInput{Bucket: s.testBucketName}); err != nil {
		s.Fail("failed to delete the test bucket", "err", err, "bucketName", *s.testBucketName)
	}
}

func TestSuiteS3FileDB(t *testing.T) {
	suite.Run(t, new(SuiteS3FileDB))
}

func (s *SuiteS3FileDB) TestS3FileDB() {
	testKey := randStrBytes(100)
	testVal := randStrBytes(1024 * 1024)

	_, err := s.s3DB.read(testKey)

	if err == nil || !strings.Contains(err.Error(), s3.ErrCodeNoSuchKey) {
		s.Fail("Failed to read", "err", err, "bucketName", *s.testBucketName)
	}

	var uris []uri
	uris, err = s.s3DB.write([]item{{key: testKey, val: testVal}})
	if err != nil {
		s.Fail("Failed to write", "err", err, "bucketName", *s.testBucketName)
	}

	if len(uris) != 1 {
		s.Fail("Unexpected amount of uris are returned", "len(uris)", len(uris))
	}

	var val []byte
	val, err = s.s3DB.read(testKey)
	s.True(bytes.Equal(testVal, val))

	err = s.s3DB.delete(testKey)
	if err != nil {
		s.Fail("Failed to delete", "err", err, "bucketName", *s.testBucketName)
	}
}

func (s *SuiteS3FileDB) TestS3FileDB_Overwrite() {
	testKey := randStrBytes(testRand.Intn(1000))
	var testVals [][]byte
	for i := 0; i < 10; i++ {
		testVals = append(testVals, randStrBytes(1024*1024))
	}

	_, err := s.s3DB.read(testKey)
	if err == nil || !strings.Contains(err.Error(), s3.ErrCodeNoSuchKey) {
		s.Fail("Failed to read", "err", err, "bucketName", *s.testBucketName)
	}

	// This is to ensure deleting the bucket after the tests
	defer s.s3DB.delete(testKey)

	var urisList [][]uri
	for _, testVal := range testVals {
		uris, err := s.s3DB.write([]item{{key: testKey, val: testVal}})
		if err != nil {
			s.Fail("failed to write the data to s3FileDB", "err", err)
		}
		urisList = append(urisList, uris)
	}

	returnedVal, err := s.s3DB.read(testKey)
	if err != nil {
		s.Fail("failed to read the data from s3FileDB", "err", err)
	}

	s.Equal(testVals[len(testVals)-1], returnedVal)
	s.Equal(len(testVals), len(urisList))
}

func (s *SuiteS3FileDB) TestS3FileDB_EmptyDelete() {
	testKey := randStrBytes(testRand.Intn(1000))
	s.NoError(s.s3DB.delete(testKey))
	s.NoError(s.s3DB.delete(testKey))
	s.NoError(s.s3DB.delete(testKey))
}
