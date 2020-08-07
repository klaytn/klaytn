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
