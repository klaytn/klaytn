package kas

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/gxhash"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/suite"
)

var models = []interface{}{
	&Tx{},
	&FetcherMetadata{},
}

// configure and generate a test block chain
var (
	config      = params.TestChainConfig
	gendb       = database.NewMemoryDBManager()
	key, _      = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	address     = crypto.PubkeyToAddress(key.PublicKey)
	funds       = big.NewInt(100000000000000000)
	testGenesis = &blockchain.Genesis{
		Config: config,
		Alloc:  blockchain.GenesisAlloc{address: {Balance: funds}},
	}
	genesis = testGenesis.MustCommit(gendb)
	signer  = types.NewEIP155Signer(config.ChainID)
)

type SuiteRepository struct {
	suite.Suite
	repo *repository
}

func genRandomAddress() *common.Address {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return &addr
}

func genRandomHash() (h common.Hash) {
	hasher := sha3.NewKeccak256()

	r := rand.Uint64()
	rlp.Encode(hasher, r)
	hasher.Sum(h[:0])

	return h
}

func makeChainEventsWithInternalTraces(numBlocks int, genTxs func(i int, block *blockchain.BlockGen)) ([]blockchain.ChainEvent, error) {
	db := database.NewMemoryDBManager()
	testGenesis.MustCommit(db)

	// create new blockchain with enabled internal tx tracing option
	b, _ := blockchain.NewBlockChain(db, nil, testGenesis.Config, gxhash.NewFaker(), vm.Config{Debug: true, EnableInternalTxTracing: true})
	defer b.Stop()

	// subscribe a new chain event channel
	chainEventCh := make(chan blockchain.ChainEvent, numBlocks)
	subscription := b.SubscribeChainEvent(chainEventCh)
	defer subscription.Unsubscribe()

	// generate blocks
	blocks, _ := blockchain.GenerateChain(testGenesis.Config, genesis, gxhash.NewFaker(), gendb, numBlocks, genTxs)

	// insert the generated blocks into the test chain
	if _, err := b.InsertChain(blocks); err != nil {
		return nil, err
	}

	var events []blockchain.ChainEvent
	for i := 0; i < numBlocks; i++ {
		timer := time.NewTimer(1 * time.Second)
		select {
		case <-timer.C:
			return nil, fmt.Errorf("timeout. too late receive a chain event: %v block", i)
		case ev := <-chainEventCh:
			events = append(events, ev)
		}
		timer.Stop()
	}

	return events, nil
}

func setTestDatabase(t *testing.T, mysql *gorm.DB) {
	// Drop previous test database if possible.
	if err := mysql.Exec("DROP DATABASE IF EXISTS test").Error; err != nil {
		if !strings.Contains(err.Error(), "database doesn't exist") {
			t.Fatal("Unexpected error happened!", "err", err)
		}
	}
	// Create new test database.
	if err := mysql.Exec("CREATE DATABASE test DEFAULT CHARACTER SET UTF8").Error; err != nil {
		t.Fatal("Error while creating test database", "err", err)
	}
	// Use test database
	if err := mysql.Exec("USE test").Error; err != nil {
		t.Fatal("Error while setting test database", "err", err)
	}

	// Auto-migrate data model from model.DataModels
	if err := mysql.AutoMigrate(models...).Error; err != nil {
		t.Fatal("Error while auto migrating data models", "err", err)
	}
}

func (s *SuiteRepository) SetupSuite() {
	id := "root"

	mysql, err := gorm.Open("mysql", fmt.Sprintf("%s@/?parseTime=True", id))
	if err != nil {
		s.T().Log("Failed connecting to mysql", "id", id, "err", err)
		s.T().Skip()
	}

	setTestDatabase(s.T(), mysql)
	s.repo = &repository{db: mysql}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(SuiteRepository))
}
