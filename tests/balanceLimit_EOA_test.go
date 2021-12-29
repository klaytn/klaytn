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
	"github.com/klaytn/klaytn/params"
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

func BalanceLimitUpdate(account *TestAccountType, balanceLimit *big.Int, signer types.EIP155Signer, keys []*ecdsa.PrivateKey, t *testing.T) *types.Transaction {
	valueMapForCreation, _ := genMapForTxTypes(account, nil, types.TxTypeBalanceLimitUpdate)
	valueMapForCreation[types.TxValueKeyGasLimit] = uint64(0) // BalanceLimit does not require gas
	valueMapForCreation[types.TxValueKeyGasPrice] = big.NewInt(0)
	valueMapForCreation[types.TxValueKeyBalanceLimit] = balanceLimit

	tx, err := types.NewTransactionWithMap(types.TxTypeBalanceLimitUpdate, valueMapForCreation)
	assert.Equal(t, nil, err)

	if len(keys) == 0 {
		keys = account.Keys
	}

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
		tx := BalanceLimitUpdate(receiver, newBalanceLimit, signer, nil, t)

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
		tx := BalanceLimitUpdate(sender, newBalanceLimit, signer, nil, t)

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

// 한 블록에 setBalanceLimit을 여러번 호출했을 때 마지막 값만 반영되는 것을 확인
func TestBalanceLimit_pending_multi_SetBalanceLimit(t *testing.T) {
	env := generateBalanceLimitEOATestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	receiver := env.receivers[0]

	// 한 account에 대하여 BalanceLimitUpdate tx를 100개 생성
	tryNum := int64(100)

	var txs types.Transactions
	var balanceLimit int64

	for i := int64(1); i <= tryNum; i++ {
		min := 0
		max := math.MaxInt32
		transferAmount := int64(rand.Intn(max-min) + min)

		// balance limit update
		{
			balanceLimit += transferAmount
			tx := BalanceLimitUpdate(receiver, big.NewInt(balanceLimit), signer, nil, t)

			txs = append(txs, tx)

			receiver.AddNonce()
		}
	}

	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}

	// klay_getAccount 실행하여 설정한 balanceLimit 조회
	{
		statedb, err := bcdata.bc.State()
		assert.NoError(t, err)
		acc := statedb.GetAccount(receiver.Addr)
		eoa, ok := acc.(account.EOA)
		assert.True(t, ok)
		assert.Equal(t, balanceLimit, eoa.GetBalanceLimit().Int64())
		assert.Equal(t, tryNum, int64(eoa.GetNonce()))
	}
}

// 한 블록에 setBalanceLimit과 Account Update를 반복했을 때, roleAccountUpdate가 setBalanceLimit할 수 있는 것을 확인
func TestBalanceLimit_pending_setBalanceLimit_accountUpdate(t *testing.T) {
	env := generateBalanceLimitEOATestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	keyRoleUpdate := sender.GetUpdateKeys()

	gasPrice := new(big.Int).SetUint64(25 * params.Ston)
	tryNum := int64(100)

	var txs types.Transactions

	// 한 account에 대하여 다음 tx를 100개 생성
	for i := int64(0); i < tryNum; i++ {
		// balance limit update
		{
			min := 0
			max := math.MaxInt32
			balanceLimit := int64(rand.Intn(max-min) + min)
			tx := BalanceLimitUpdate(sender, big.NewInt(balanceLimit), signer, keyRoleUpdate, t)

			txs = append(txs, tx)

			sender.AddNonce()
		}

		// account update
		{
			// generate a role-based key
			prvKeys := genTestKeys(3)
			roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
				accountkey.NewAccountKeyPublicWithValue(&prvKeys[0].PublicKey),
				accountkey.NewAccountKeyPublicWithValue(&keyRoleUpdate[0].PublicKey),
				accountkey.NewAccountKeyPublicWithValue(&prvKeys[2].PublicKey),
			})
			// create tx
			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:      sender.Nonce,
				types.TxValueKeyFrom:       sender.Addr,
				types.TxValueKeyGasLimit:   gasLimit,
				types.TxValueKeyGasPrice:   gasPrice,
				types.TxValueKeyAccountKey: roleKey,
			}

			tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
			assert.Equal(t, nil, err)

			// sign with the previous update key
			err = tx.SignWithKeys(signer, keyRoleUpdate)
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			sender.AddNonce()
		}
	}

	// 블록 생성 (모든 tx 하나의 블록에)
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}
}

// 한 블록에 Account Update, setBalanceLimit, setAccountStatus를 반복했을 때, roleAccountUpdate가 setBalanceLimit, setAccountStatus 할 수 있는 것을 확인
func TestBalanceLimit_pending_setBalanceLimit_accountUpdate_setAccountStatus(t *testing.T) {
	env := generateBalanceLimitEOATestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	keyRoleUpdate := sender.GetUpdateKeys()

	gasPrice := new(big.Int).SetUint64(25 * params.Ston)
	var txs types.Transactions

	// 한 account에 대하여 다음 tx를 100개 생성
	tryNum := int64(100)
	for i := int64(0); i < tryNum; i++ {
		// account update
		{
			// generate a role-based key
			prvKeys := genTestKeys(3)
			roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
				accountkey.NewAccountKeyPublicWithValue(&prvKeys[0].PublicKey),
				accountkey.NewAccountKeyPublicWithValue(&keyRoleUpdate[0].PublicKey),
				accountkey.NewAccountKeyPublicWithValue(&prvKeys[2].PublicKey),
			})
			// create tx
			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:      sender.Nonce,
				types.TxValueKeyFrom:       sender.Addr,
				types.TxValueKeyGasLimit:   gasLimit,
				types.TxValueKeyGasPrice:   gasPrice,
				types.TxValueKeyAccountKey: roleKey,
			}

			tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
			assert.Equal(t, nil, err)

			// sign with the previous update key
			err = tx.SignWithKeys(signer, keyRoleUpdate)
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			sender.AddNonce()
		}

		// balance limit update
		{
			min := 0
			max := math.MaxInt32
			balanceLimit := int64(rand.Intn(max-min) + min)
			tx := BalanceLimitUpdate(sender, big.NewInt(balanceLimit), signer, keyRoleUpdate, t)

			txs = append(txs, tx)

			sender.AddNonce()
		}

		// account status update (stop)
		{
			tx := AccountStatusUpdate(sender, account.AccountStatusStop, signer, keyRoleUpdate, t)

			txs = append(txs, tx)

			sender.AddNonce()
		}

		// account status update (active)
		{
			tx := AccountStatusUpdate(sender, account.AccountStatusActive, signer, keyRoleUpdate, t)

			txs = append(txs, tx)

			sender.AddNonce()
		}
	}

	// 블록 생성 (모든 tx 하나의 블록에)
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}
	txs = types.Transactions{}
}
