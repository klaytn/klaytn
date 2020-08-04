package kas

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/klaytn/klaytn/log"
	"time"
)

const (
	maxTxCountPerBlock      = int64(1000000)
	maxTxLogCountPerTx      = int64(100000)
	maxInternalTxCountPerTx = int64(10000)

	maxPlaceholders = 65535

	placeholdersPerTxItem          = 13
	placeholdersPerKCTTransferItem = 7

	maxDBRetryCount = 20
	DBRetryInterval = 1 * time.Second
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)

type repository struct {
	db *gorm.DB
}

func getEndpoint(user, password, host, port, name string) string {
	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + name + "?parseTime=True&charset=utf8mb4"
}

func NewRepository(user, password, host, port, name string) (*repository, error) {
	endpoint := getEndpoint(user, password, host, port, name)
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < maxDBRetryCount; i++ {
		db, err = gorm.Open("mysql", endpoint)
		if err != nil {
			logger.Warn("Retrying to connect DB", "endpoint", endpoint, "err", err)
			time.Sleep(DBRetryInterval)
		} else {
			// TODO-ChainDataFetcher insert other options such as maxOpen, maxIdle, maxLifetime, etc.
			//db.DB().SetMaxOpenConns(maxOpen)
			//db.DB().SetMaxIdleConns(maxIdle)
			//db.DB().SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

			return &repository{db: db}, nil
		}
	}
	logger.Error("Failed to connect to the database", "endpoint", endpoint, "err", err)
	return nil, err
}
