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

const (
	TxTableName          = "klay_transfers"
	KctTransferTableName = "kct_transfers"
	RevertedTxTableName  = "reverted_transactions"
)

type Tx struct {
	TransactionId   int64  `gorm:"column:transactionId;type:BIGINT;INDEX:idIdx;NOT NULL;PRIMARY_KEY"`
	FromAddr        []byte `gorm:"column:fromAddr;type:VARBINARY(20);INDEX:txFromAddrIdx"`
	ToAddr          []byte `gorm:"column:toAddr;type:VARBINARY(20);INDEX:txToAddrIdx"`
	Value           string `gorm:"column:value;type:VARCHAR(80)"`
	TransactionHash []byte `gorm:"column:transactionHash;type:VARBINARY(32);INDEX:tHashIdx;NOT NULL"`
	Status          int    `gorm:"column:status;type:SMALLINT"`
	Timestamp       int64  `gorm:"column:timestamp;type:INT(11)"`
	TypeInt         int    `gorm:"column:typeInt;INDEX:tTypeIdx;NOT NULL"`
	GasPrice        uint64 `gorm:"column:gasPrice;type:BIGINT"`
	GasUsed         uint64 `gorm:"column:gasUsed;type:BIGINT"`
	FeePayer        []byte `gorm:"column:feePayer;type:VARBINARY(20)"`
	FeeRatio        uint   `gorm:"column:feeRatio;type:INT"`
	Internal        bool   `gorm:"column:internal;type:TINYINT(1);DEFAULT:0"`
}

func (Tx) TableName() string {
	return TxTableName
}

type KCTTransfer struct {
	ContractAddress  []byte `gorm:"column:contractAddress;type:VARBINARY(20);INDEX:ttFromCompIdx,ttToCompIdx;NOT NULL"`
	From             []byte `gorm:"column:fromAddr;type:VARBINARY(20);INDEX:ttFromCompIdx,ttFromIdx"`
	To               []byte `gorm:"column:toAddr;type:VARBINARY(20);INDEX:ttToCompIdx,ttToIdx"`
	TransactionLogId int64  `gorm:"column:transactionLogId;type:BIGINT;PRIMARY_KEY;INDEX:ttFromCompIdx,ttToCompIdx"`
	Value            string `gorm:"column:value;type:VARCHAR(80)"`
	TransactionHash  []byte `gorm:"column:transactionHash;type:VARBINARY(32);INDEX:ttHashIdx;NOT NULL"`
	Timestamp        int64  `gorm:"column:timestamp;type:INT(11)"`
}

func (KCTTransfer) TableName() string {
	return KctTransferTableName
}

type RevertedTx struct {
	TransactionHash []byte `gorm:"column:transactionHash;type:VARBINARY(32);NOT NULL;PRIMARY_KEY"`
	BlockNumber     int64  `gorm:"column:blockNumber;type:BIGINT"`
	RevertMessage   string `gorm:"column:revertMessage;type:VARCHAR(1024)"`
	ContractAddress []byte `gorm:"column:contractAddress;type:VARBINARY(20);NOT NULL"`
	Timestamp       int64  `gorm:"column:timestamp;type:INT(11)"`
}

func (RevertedTx) TableName() string {
	return RevertedTxTableName
}
