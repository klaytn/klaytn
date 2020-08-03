package kas

const TxTableName = "klay_transfers"

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
