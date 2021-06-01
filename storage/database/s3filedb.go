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
// FileDB implementation of AWS S3.
//
// [WARN] Using this DB may cause pricing in your AWS account.
//
// You need to set AWS credentials to access to S3.
//    $ export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
//    $ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET

package database

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/klaytn/klaytn/common/hexutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/klaytn/klaytn/log"
)

// s3FileDB is an implementation of fileDB based on AWS S3.
// It stores the data to the designated AWS S3 bucket.
type s3FileDB struct {
	region   string
	endpoint string
	bucket   string
	s3       *s3.S3
	logger   log.Logger
}

// newS3FileDB returns a new s3FileDB with the given region, endpoint and bucketName.
// If the given bucket does not exist, it creates one.
func newS3FileDB(region, endpoint, bucketName string) (*s3FileDB, error) {
	localLogger := logger.NewWith("endpoint", endpoint, "bucketName", bucketName)
	sessionConf, err := session.NewSession(&aws.Config{
		Retryer: CustomRetryer{
			DefaultRetryer: client.DefaultRetryer{
				NumMaxRetries:    dynamoMaxRetry,
				MaxRetryDelay:    time.Second,
				MaxThrottleDelay: time.Second,
			},
		},
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		localLogger.Error("failed to create session", "region", region)
		return nil, err
	}

	s3DB := &s3FileDB{
		region:   region,
		endpoint: endpoint,
		bucket:   bucketName,
		s3:       s3.New(sessionConf),
		logger:   localLogger,
	}

	exist, err := s3DB.hasBucket(bucketName)
	if err != nil {
		localLogger.Error("failed to retrieve a bucket list", "err", err)
		return nil, err
	}

	if !exist {
		localLogger.Warn("creating a S3 bucket. You will be CHARGED until the bucket is deleted")
		_, err = s3DB.s3.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			localLogger.Error("failed to create a bucket", "err", err)
			return nil, err
		}
	}
	localLogger.Info("successfully created S3 session")
	return s3DB, nil
}

// hasBucket returns if the bucket exists in the endpoint of s3FileDB.
func (s3DB *s3FileDB) hasBucket(bucketName string) (bool, error) {
	output, err := s3DB.s3.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return false, err
	}

	bucketExist := false
	for _, bucket := range output.Buckets {
		if bucketName == *bucket.Name {
			bucketExist = true
			break
		}
	}
	return bucketExist, nil
}

// write puts list of items to its bucket and returns the list of URIs.
func (s3DB *s3FileDB) write(item item) (string, error) {
	o := &s3.PutObjectInput{
		Bucket:      aws.String(s3DB.bucket),
		Key:         aws.String(hexutil.Encode(item.key)),
		Body:        bytes.NewReader(item.val),
		ContentType: aws.String("application/octet-stream"),
	}

	if _, err := s3DB.s3.PutObject(o); err != nil {
		return "", fmt.Errorf("failed to write item to S3. key: %v, err: %w", string(item.key), err)
	}

	return hexutil.Encode(item.key), nil
}

// read gets the data from the bucket with the given key.
func (s3DB *s3FileDB) read(key []byte) ([]byte, error) {
	output, err := s3DB.s3.GetObject(&s3.GetObjectInput{
		Bucket:              aws.String(s3DB.bucket),
		Key:                 aws.String(hexutil.Encode(key)),
		ResponseContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		return nil, err
	}

	returnVal, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return returnVal, nil
}

// delete removes the data with the given key from the bucket.
// No error is returned if the data with the given key does not exist.
func (s3DB *s3FileDB) delete(key []byte) error {
	_, err := s3DB.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3DB.bucket),
		Key:    aws.String(hexutil.Encode(key)),
	})
	return err
}

// deleteBucket removes the bucket
func (s3DB *s3FileDB) deleteBucket() {
	if _, err := s3DB.s3.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(s3DB.bucket)}); err != nil {
		s3DB.logger.Error("failed to delete the test bucket", "err", err, "bucketName", s3DB.bucket)
	}
}
