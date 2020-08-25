package kas

import (
	"math/big"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

// TODO-ChainDataFetcher add more tx types test cases
func genNewValueTransfer(from, to common.Address, nonce, amount, gasLimit, gasPrice *big.Int) *types.Transaction {
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    nonce.Uint64(),
		types.TxValueKeyTo:       to,
		types.TxValueKeyAmount:   amount,
		types.TxValueKeyGasLimit: gasLimit.Uint64(),
		types.TxValueKeyGasPrice: gasPrice,
		types.TxValueKeyFrom:     from,
	})
	if err != nil {
		panic(err)
	}

	err = tx.Sign(signer, key)
	if err != nil {
		panic(err)
	}

	return tx
}

func checkValidChainEventsPosted(t *testing.T, expectedBlocks, expectedTxsPerBlock int, events []blockchain.ChainEvent) {
	assert.Equal(t, expectedBlocks, len(events))
	for _, ev := range events {
		assert.Equal(t, expectedTxsPerBlock, len(ev.Block.Transactions()))
	}
}

func TestRepository_TransformToTxs_Success(t *testing.T) {
	// define transaction contents
	from := address
	randKey, _ := crypto.GenerateKey()
	to := crypto.PubkeyToAddress(randKey.PublicKey)
	amount := new(big.Int).SetUint64(100)
	gasLimit := big.NewInt(500000)
	gasPrice := new(big.Int).SetUint64(1 * params.Peb)

	// make a chain event including a value transfer transaction
	numBlocks := 1
	numTxsPerBlock := 1
	events, err := makeChainEventsWithInternalTraces(numBlocks, func(i int, block *blockchain.BlockGen) {
		nonce := new(big.Int).SetUint64(block.TxNonce(from))
		tx := genNewValueTransfer(from, to, nonce, amount, gasLimit, gasPrice)
		block.AddTx(tx)
	})
	assert.NoError(t, err)

	checkValidChainEventsPosted(t, numBlocks, numTxsPerBlock, events)

	ev := events[0]
	txs, _ := transformToTxs(ev)

	for _, tx := range txs {
		assert.Equal(t, from.Bytes(), tx.FromAddr)
		assert.Equal(t, to.Bytes(), tx.ToAddr)
		assert.Equal(t, "0x"+amount.Text(16), tx.Value)
		assert.Equal(t, gasPrice.Uint64(), tx.GasPrice)
		t.Log(tx.Value)
	}
}

func (s *SuiteRepository) TestRepository_InsertTransactions_Panics_EmptyChainEvent() {
	s.Panics(func() { s.repo.InsertTransactions(blockchain.ChainEvent{}) })
}

func (s *SuiteRepository) TestRepository_bulkInsertTransaction_Success() {
	from := genRandomAddress()
	to := genRandomAddress()
	feePayer := genRandomAddress()

	expected := &Tx{
		TransactionId:   2020*(1000000*100000) + (0 * 10000) + 0,
		FromAddr:        from.Bytes(),
		ToAddr:          to.Bytes(),
		Value:           "0x1234",
		TransactionHash: common.HexToHash("0xd17153512821d290a29589e1decfbd26b6c2792faffe0d1e2aa664d3f4820fd1").Bytes(),
		Status:          1,
		Timestamp:       9999,
		TypeInt:         int(types.TxTypeValueTransfer),
		GasPrice:        8888,
		GasUsed:         7777,
		FeePayer:        feePayer.Bytes(),
		FeeRatio:        uint(30),
		Internal:        false,
	}

	err := s.repo.bulkInsertTransactions([]*Tx{expected})
	s.NoError(err)

	actual := &Tx{}
	s.NoError(s.repo.db.Where("fromAddr = ?", from.Bytes()).Take(actual).Error)
	s.Equal(expected, actual)
}

func (s *SuiteRepository) TestRepository_bulkInsertTransactions_Success_NoTransaction() {
	s.Nil(s.repo.bulkInsertTransactions(nil))
}

func (s *SuiteRepository) TestRepository_bulkInsertTransactions_Fail_NoTxhash() {
	tx := &Tx{
		FromAddr: genRandomAddress().Bytes(),
		ToAddr:   genRandomAddress().Bytes(),
		Value:    "0x1",
	}

	err := s.repo.bulkInsertTransactions([]*Tx{tx})
	s.Error(err)
	s.True(strings.Contains(err.Error(), "'transactionHash' cannot be null"))
}

func createChainEventsWithTooManyTxs() ([]blockchain.ChainEvent, error) {
	// define transaction contents
	from := address
	randKey, _ := crypto.GenerateKey()
	to := crypto.PubkeyToAddress(randKey.PublicKey)
	amount := new(big.Int).SetUint64(100)
	gasLimit := big.NewInt(500000)
	gasPrice := new(big.Int).SetUint64(1 * params.Peb)

	// make a chain event including a value transfer transaction
	numBlocks := 1
	return makeChainEventsWithInternalTraces(numBlocks, func(i int, block *blockchain.BlockGen) {
		for idx := 0; idx < maxPlaceholders/placeholdersPerTxItem+1; idx++ {
			nonce := new(big.Int).SetUint64(block.TxNonce(from))
			tx := genNewValueTransfer(from, to, nonce, amount, gasLimit, gasPrice)
			block.AddTx(tx)
		}
	})
}

func (s *SuiteRepository) TestRepository_bulkInsertTransactions_Fail_TooManyPlaceholders() {
	events, err := createChainEventsWithTooManyTxs()
	s.NoError(err)

	ev := events[0]
	txs, _ := transformToTxs(ev)
	err = s.repo.bulkInsertTransactions(txs)
	s.Error(err)
	s.True(strings.Contains(err.Error(), "Prepared statement contains too many placeholders"))
}

func (s *SuiteRepository) TestRepository_InsertTransactions_Success_TooManyPlaceholders() {
	events, err := createChainEventsWithTooManyTxs()
	s.NoError(err)

	ev := events[0]
	err = s.repo.InsertTransactions(ev)
	s.NoError(err)
}

func (s *SuiteRepository) TestRepository_InsertTransactions_Success_EmptyBlock() {
	numBlocks := 1
	numTxsPerBlock := 0
	events, err := makeChainEventsWithInternalTraces(numBlocks, nil)
	s.NoError(err)

	checkValidChainEventsPosted(s.T(), numBlocks, numTxsPerBlock, events)
	s.NoError(s.repo.InsertTransactions(events[0]))
}
