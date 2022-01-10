package tests

import (
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type balanceLimitEOATestEnv struct {
	prof       *profile.Profiler
	bcdata     *BCData
	accountMap *AccountMap
	signer     types.EIP155Signer
	sender     *TestAccountType
	receivers  []*TestAccountType
}

func generateBalanceLimitEOATestEnv(t *testing.T) *balanceLimitEOATestEnv {
	if testing.Verbose() {
		enableLog()
	}

	prof := profile.NewProfiler()

	// Initialize blockchain
	initialBalance := new(big.Int).Mul(big.NewInt(100), account.GetInitialBalanceLimit())
	bcdata, err := NewBCData(6, 4, initialBalance)
	if err != nil {
		t.Fatal(err)
	}

	// Initialize address-balance map for verification
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}

	signer := types.NewEIP155Signer(bcdata.bc.Config().ChainID)

	// sender account
	sender := &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	// test account will be lack of KLAY
	receiverNum := 3
	receivers := make([]*TestAccountType, receiverNum)
	for i := 0; i < receiverNum; i++ {
		account, err := createDefaultAccount(accountkey.AccountKeyTypeLegacy)
		assert.Equal(t, nil, err)

		receivers[i] = account
		accountMap.Set(account.Addr, big.NewInt(0), account.Nonce)
	}

	return &balanceLimitEOATestEnv{
		prof:       prof,
		bcdata:     bcdata,
		accountMap: accountMap,
		signer:     signer,
		sender:     sender,
		receivers:  receivers,
	}
}

func ValueTransfer(sender *TestAccountType, receiver *TestAccountType, amount *big.Int, signer types.EIP155Signer, t *testing.T) *types.Transaction {
	valueMapForCreation, _ := genMapForTxTypes(sender, receiver, types.TxTypeValueTransfer)
	valueMapForCreation[types.TxValueKeyAmount] = amount

	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, valueMapForCreation)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, sender.Keys)
	assert.Equal(t, nil, err)

	return tx
}

func BalanceLimitUpdate(account *TestAccountType, balanceLimit *big.Int, signer types.EIP155Signer, t *testing.T) *types.Transaction {
	valueMapForCreation, _ := genMapForTxTypes(account, nil, types.TxTypeBalanceLimitUpdate)
	valueMapForCreation[types.TxValueKeyGasLimit] = uint64(0) // BalanceLimit does not require gas
	valueMapForCreation[types.TxValueKeyGasPrice] = big.NewInt(0)
	valueMapForCreation[types.TxValueKeyBalanceLimit] = balanceLimit

	tx, err := types.NewTransactionWithMap(types.TxTypeBalanceLimitUpdate, valueMapForCreation)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, account.Keys)
	assert.Equal(t, nil, err)

	return tx
}

// 존재하지 않는 account에 InitialBalanceLimit보다 큰 값을 전송할 수 없는 것을 확인
func TestBalanceLimit_EOA_InitialBalanceLimit_newAccount(t *testing.T) {
	env := generateBalanceLimitEOATestEnv(t)
	defer env.bcdata.Shutdown()

	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	receivers := env.receivers

	// InitialBalanceLimit보다 작은 값을 value transfer → 성공
	{
		min := 0
		max := math.MaxInt32
		transferAmount := big.NewInt(int64(rand.Intn(max-min) + min))
		tx := ValueTransfer(sender, receivers[0], transferAmount, signer, t)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, nil, err)
		assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	// InitialBalanceLimit보다 같은 값을 value transfer → 성공
	{
		transferAmount := account.GetInitialBalanceLimit()
		tx := ValueTransfer(sender, receivers[1], transferAmount, signer, t)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, nil, err)
		assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	// InitialBalanceLimit보다 큰 값을 value transfer → 실패
	{
		min := 1
		max := math.MaxInt32
		transferAmount := new(big.Int).Add(big.NewInt(int64(rand.Intn(max-min)+min)), account.GetInitialBalanceLimit())

		tx := ValueTransfer(sender, receivers[2], transferAmount, signer, t)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.NoError(t, err)
		assert.Equal(t, types.ReceiptStatusErrExceedBalanceLimit, receipt.Status)
	}
}

// 존재하는 account에 대해 setBalanceLimit이 보유한도를 변경하는 것을 확인
func TestBalanceLimit_EOA_setBalanceLimit(t *testing.T) {
	env := generateBalanceLimitEOATestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	receiver := env.receivers[0]

	initialBalanceLimt := account.GetInitialBalanceLimit()
	newBalanceLimit := new(big.Int).Mul(initialBalanceLimt, big.NewInt(2))

	min := 1
	max := math.MaxInt32
	transferAmount := new(big.Int).Add(big.NewInt(int64(rand.Intn(max-min)+min)), account.GetInitialBalanceLimit())

	// InitialBalanceLimit보다 많은 값을 보낼 수 없는 것을 확인
	{
		tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.NoError(t, err)
		assert.Equal(t, types.ReceiptStatusErrExceedBalanceLimit, receipt.Status)
	}

	// B에 setBalanceLimit을 이용하여 InitialBalanceLimit 더 높은 값(=NewBalanceLimit)으로 한도 변경
	{
		tx := BalanceLimitUpdate(receiver, newBalanceLimit, signer, t)

		txs := []*types.Transaction{tx}
		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		receiver.AddNonce()
	}

	// InitialBalanceLimit보다 큰 값을 보낼 수 있는 것을 확인
	{
		tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

		txs := []*types.Transaction{tx}
		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		sender.AddNonce()
	}

	// NewBalanceLimit보다 큰 값을 보낼 수 없는 것을 확인
	{
		transferAmount := new(big.Int).Add(newBalanceLimit, big.NewInt(1))
		tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.NoError(t, err)
		assert.Equal(t, types.ReceiptStatusErrExceedBalanceLimit, receipt.Status)
	}
}

// balanceLimit을 넘는 tx를 보내면 에러 발생하는 것을 확인 (txpool 에러 발생)
func TestBalanceLimit_EOA_ErrorMessage(t *testing.T) {
	env := generateBalanceLimitEOATestEnv(t)
	defer env.bcdata.Shutdown()

	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	receiver := env.receivers[0]
	initialBalanceLimt := account.GetInitialBalanceLimit()

	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	{
		transferAmount := new(big.Int).Add(initialBalanceLimt, big.NewInt(1))
		tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

		err := txpool.AddRemote(tx)
		assert.Error(t, err)
		assert.Equal(t, err, blockchain.ErrExceedBalanceLimit)
	}
}

// klay_getAccount에서 balanceLimit을 보여주는 것을 확인
func TestBalanceLimit_EOA_klay_account(t *testing.T) {
	env := generateBalanceLimitEOATestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender

	min := 0
	max := 100000000
	newBalanceLimit := big.NewInt(int64(rand.Intn(max-min) + min))

	// setBalanceLimit 호출 (셋팅하는 값은 rand로)
	{
		tx := BalanceLimitUpdate(sender, newBalanceLimit, signer, t)

		txs := []*types.Transaction{tx}
		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		sender.AddNonce()
	}

	// B에 대해 klay_getAccount 실행하여 설정한 balanceLimit 조회
	{
		statedb, err := bcdata.bc.State()
		assert.NoError(t, err)
		acc := statedb.GetAccount(sender.Addr)
		eoa, ok := acc.(account.EOA)
		assert.True(t, ok)
		assert.Equal(t, newBalanceLimit, eoa.GetBalanceLimit())
	}
}

// 여러 개의 트랜잭션을 한 블록에 커밋했을 때, 일부 트랜잭션만 실패하는지 확인하기 위함
func TestBalanceLimit_EOA_setBalanceLimit_pendingNonce(t *testing.T) {
	env := generateBalanceLimitEOATestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	receiver := env.receivers[0]

	min := 0
	max := math.MaxInt32

	transferAmounts := []*big.Int{
		big.NewInt(int64(rand.Intn(max-min) + min)),
		big.NewInt(int64(rand.Intn(max-min) + min)),
		big.NewInt(int64(rand.Intn(max-min) + min)),
		account.GetInitialBalanceLimit(), // 처리 못해야 함
		account.GetInitialBalanceLimit(), // 처리 못해야 함
	}

	var txs types.Transactions
	for _, transferAmount := range transferAmounts[:3] {
		tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

		txs = append(txs, tx)

		sender.AddNonce()
	}
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}

	sum := new(big.Int).Add(new(big.Int).Add(transferAmounts[0], transferAmounts[1]), transferAmounts[2])
	statedb, err := bcdata.bc.State()
	assert.NoError(t, err)
	assert.Equal(t, sum, statedb.GetBalance(receiver.Addr))
}
