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
// Database implementation of AWS DynamoDB READ ONLY.
// Calling put, delete, batch put and batch write does nothing and returns no error.
// Other functions such as get and has will call functions in dynamoDB.
//
// [WARN] Using this DB may cause pricing in your AWS account.
// [WARN] DynamoDB creates both Dynamo DB table and S3 bucket.
//
// You need to set AWS credentials to access to dynamoDB.
//    $ export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
//    $ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET

package database

// dynamoDBReadOnly uses dynamoDB.
// Calling put, delete, batch put and batch write does nothing and returns no error.
// Other functions such as get and has will call functions in dynamoDB.
type dynamoDBReadOnly struct {
	dynamoDB
}

// newDynamoDBReadOnly creates dynamoDBReadOnly. This uses dynamoDB to create dynamoDBReadOnly.
func newDynamoDBReadOnly(config *DynamoDBConfig) (*dynamoDBReadOnly, error) {
	config.ReadOnly = true
	dynamo, err := newDynamoDB(config)
	if err != nil {
		return nil, err
	}
	dynamo.logger.Warn("Read only flag is set for dynamoDB. Check for billing method and capacity.")
	return &dynamoDBReadOnly{*dynamo}, nil
}

func (dynamo *dynamoDBReadOnly) Put(key []byte, val []byte) error {
	return nil
}

func (dynamo *dynamoDBReadOnly) Delete(key []byte) error {
	return nil
}

func (dynamo *dynamoDBReadOnly) Close() {
}

func (dynamo *dynamoDBReadOnly) NewBatch() Batch {
	return &emptyBatch{}
}

type emptyBatch struct{}

func (batch *emptyBatch) Put(key, val []byte) error {
	return nil
}

func (batch *emptyBatch) Delete(key []byte) error {
	return nil
}

func (batch *emptyBatch) Write() error {
	return nil
}

func (batch *emptyBatch) ValueSize() int {
	return 0
}

func (batch *emptyBatch) Reset() {
}

func (batch *emptyBatch) Replay(w KeyValueWriter) error {
	return nil
}
