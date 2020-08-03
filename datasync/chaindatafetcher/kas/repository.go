package kas

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/klaytn/klaytn/log"
)

const (
	maxTransactionCount         = int64(1000000)
	maxTransactionLogCount      = int64(100000)
	maxInternalTransactionCount = int64(10000)

	maxPlaceholders = 65535

	placeholdersPerTxItem = 13
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)

type repository struct {
	db *gorm.DB
}

func getEndpoint(user, password, host, port, name string) string {
	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + name + "?parseTime=True&charset=utf8mb4"
}

func NewRepository(user, password, host, port, name string) *repository {
	endpoint := getEndpoint(user, password, host, port, name)
	db, err := gorm.Open("mysql", endpoint)
	if err != nil {
		logger.Crit("Connecting to DB is failed", "endpoint", endpoint, "err", err)
	}

	// TODO-ChainDataFetcher insert other options such as maxOpen, maxIdle, maxLifetime, etc.
	//db.DB().SetMaxOpenConns(maxOpen)
	//db.DB().SetMaxIdleConns(maxIdle)
	//db.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	return &repository{db: db}
}
