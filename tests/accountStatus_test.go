package tests

import (
	"crypto/ecdsa"
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common/profile"
	"github.com/klaytn/klaytn/params"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type accountStatusTestEnv struct {
	prof                *profile.Profiler
	bcdata              *BCData
	accountMap          *AccountMap
	signer              types.EIP155Signer
	sender              *TestAccountType
	receivers           []*TestAccountType
	initialBalanceLimit *big.Int
}

func generateAccountStatusTestEnv(t *testing.T) *accountStatusTestEnv {
	if testing.Verbose() {
		enableLog()
	}

	prof := profile.NewProfiler()

	// Initialize blockchain
	bcdata, err := NewBCData(6, 4, nil)
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

	return &accountStatusTestEnv{
		prof:       prof,
		bcdata:     bcdata,
		accountMap: accountMap,
		signer:     signer,
		sender:     sender,
		receivers:  receivers,
	}
}

func AccountStatusUpdate(account *TestAccountType, accountStatus account.AccountStatus, signer types.EIP155Signer, keys []*ecdsa.PrivateKey, t *testing.T) *types.Transaction {
	valueMapForCreation, _ := genMapForTxTypes(account, nil, types.TxTypeAccountStatusUpdate)
	valueMapForCreation[types.TxValueKeyGasLimit] = uint64(0) // AccountStatusUpdate does not require gas
	valueMapForCreation[types.TxValueKeyGasPrice] = big.NewInt(0)
	valueMapForCreation[types.TxValueKeyAccountStatus] = uint64(accountStatus)

	tx, err := types.NewTransactionWithMap(types.TxTypeAccountStatusUpdate, valueMapForCreation)
	assert.Equal(t, nil, err)

	if keys == nil {
		keys = account.Keys
	}
	err = tx.SignWithKeys(signer, keys)
	assert.Equal(t, nil, err)

	return tx
}

func GetAccountStatus(sender *TestAccountType, bcdata *BCData, t *testing.T) account.AccountStatus {
	statedb, err := bcdata.bc.State()
	assert.NoError(t, err)
	acc := statedb.GetAccount(sender.Addr)
	eoa, ok := acc.(account.EOA)
	assert.True(t, ok)

	return eoa.GetAccountStatus()
}

// EOA에 대하여 getAccountStatus을 호출할 때,  출력 ActiveAccount(0) 확인
func TestAccountStatus_getAccountStatus_EOA(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	bcdata := env.bcdata
	sender := env.sender

	// getAccountStatus 호출
	accountStatus := GetAccountStatus(sender, bcdata, t)
	assert.Equal(t, account.AccountStatusActive, accountStatus)
}

// EOA에 대하여 setAccountStatus을 호출할 때, AccountStatus 설정하는 것을 확인
func TestAccountStatus_setAccountStatus_EOA(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender

	for accountStatus := account.AccountStatusUndefined + 1; accountStatus < account.AccountStatusLast; accountStatus++ {
		// setAccountStatus 호출
		tx := AccountStatusUpdate(sender, accountStatus, signer, nil, t)

		txs := []*types.Transaction{tx}
		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		sender.AddNonce()

		// getAccountStatus 실행 시 설정한 값이 셋팅된 것을 확인
		getAccountStatus := GetAccountStatus(sender, bcdata, t)
		assert.Equal(t, accountStatus, getAccountStatus)
	}
}

// accountStatusLast 보다 작은 값만 설정할 수 있고, accountStatusLast와 같거나 큰 값은 설정할 수 없음을 확인
func TestAccountStatus_setAccountStatus_Range(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender

	// accountStatus가 account.AccountStatusUndefined 일 때 실패
	valueMapForCreation, _ := genMapForTxTypes(sender, nil, types.TxTypeAccountStatusUpdate)
	valueMapForCreation[types.TxValueKeyGasLimit] = uint64(0) // AccountStatusUpdate does not require gas
	valueMapForCreation[types.TxValueKeyGasPrice] = big.NewInt(0)
	valueMapForCreation[types.TxValueKeyAccountStatus] = uint64(account.AccountStatusUndefined)

	_, err := types.NewTransactionWithMap(types.TxTypeAccountStatusUpdate, valueMapForCreation)
	assert.Error(t, err)

	// accountStatusLast 보다 작은 값에 대해서 setAccountStatus 성공
	for accountStatus := account.AccountStatus(1); accountStatus < account.AccountStatusLast; accountStatus++ {
		tx := AccountStatusUpdate(sender, accountStatus, signer, nil, t)

		txs := []*types.Transaction{tx}
		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		sender.AddNonce()
	}

	// accountStatusLast와 같거나 큰 값에 대해서 setAccountStatus 실패
	for accountStatus := account.AccountStatusLast; accountStatus < account.AccountStatusLast+1; accountStatus++ {
		valueMapForCreation, _ := genMapForTxTypes(sender, nil, types.TxTypeAccountStatusUpdate)
		valueMapForCreation[types.TxValueKeyGasLimit] = uint64(0) // AccountStatusUpdate does not require gas
		valueMapForCreation[types.TxValueKeyGasPrice] = big.NewInt(0)
		valueMapForCreation[types.TxValueKeyAccountStatus] = uint64(accountStatus)

		// tx 생성 단계에서 실패
		_, err := types.NewTransactionWithMap(types.TxTypeAccountStatusUpdate, valueMapForCreation)
		assert.Error(t, err)
	}
}

// 존재하는 account에 대해 getAccountStatus을 호출할 때, 초기값 출력 ActiveAccount(0) 확인
func TestAccountStatus_InitialAccountStatus(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	receiver := env.receivers[0]

	// A → B로 value transfer (B를 EOA로 설정)
	min := 0
	max := math.MaxInt32
	transferAmount := big.NewInt(int64(rand.Intn(max-min) + min))
	tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

	txs := []*types.Transaction{tx}
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}
	sender.AddNonce()

	// B에 대해 getAccountStatus 실행하여 ActiveAccount(0) 조회
	getAccountStatus := GetAccountStatus(sender, bcdata, t)
	assert.Equal(t, account.AccountStatusActive, getAccountStatus)
}

// 송신인의 계정 상태가 stop일 때 tx를 보낼 수 없는 것을 확인
func TestAccountStatus_From(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	receiver := env.receivers[0]

	// tx 생성할 수 있음을 확인 함
	{
		// legacy transaction 생성
		{
			var txs types.Transactions

			valueMap, _ := genMapForTxTypes(sender, sender, types.TxTypeLegacyTransaction)
			tx, err := types.NewTransactionWithMap(types.TxTypeLegacyTransaction, valueMap)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, sender.Keys)
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
				t.Fatal(err)
			}
			sender.AddNonce()
		}

		// non-legacy transaction 생성
		{
			min := 0
			max := math.MaxInt32
			transferAmount := big.NewInt(int64(rand.Intn(max-min) + min))
			tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

			txs := []*types.Transaction{tx}
			if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
				t.Fatal(err)
			}
			sender.AddNonce()
		}
	}

	// AccountStatus를 Stop로 설정함
	{
		tx := AccountStatusUpdate(sender, account.AccountStatusStop, signer, nil, t)

		txs := []*types.Transaction{tx}
		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		sender.AddNonce()
	}

	// tx 생성할 수 없음을 확인 함
	{
		// legacy transaction에 대해 확인
		{
			valueMap, _ := genMapForTxTypes(sender, sender, types.TxTypeLegacyTransaction)
			tx, err := types.NewTransactionWithMap(types.TxTypeLegacyTransaction, valueMap)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, sender.Keys)
			assert.Equal(t, nil, err)

			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.NoError(t, err)
			assert.Equal(t, types.ReceiptStatusErrStoppedAccountFrom, receipt.Status)
		}

		// non-legacy transaction에 대해 확인
		{
			min := 0
			max := math.MaxInt32
			transferAmount := big.NewInt(int64(rand.Intn(max-min) + min))
			tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.NoError(t, err)
			assert.Equal(t, types.ReceiptStatusErrStoppedAccountFrom, receipt.Status)
		}
	}

	// RoleUpdate 실행가능한 것을 확인함
	{
		// AccountBalanceLimitUpdate 실행 가능 확인
		{
			valueMapForCreation, _ := genMapForTxTypes(sender, nil, types.TxTypeBalanceLimitUpdate)
			valueMapForCreation[types.TxValueKeyGasLimit] = uint64(0) // BalanceLimit does not require gas
			valueMapForCreation[types.TxValueKeyGasPrice] = big.NewInt(0)
			valueMapForCreation[types.TxValueKeyBalanceLimit] = new(big.Int).Mul(big.NewInt(params.KLAY), big.NewInt(params.KLAY))

			tx, err := types.NewTransactionWithMap(types.TxTypeBalanceLimitUpdate, valueMapForCreation)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, sender.Keys)
			assert.Equal(t, nil, err)

			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.NoError(t, err)
			assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
		}

		// AccountUpdate 실행인 가능 확인
		{
			newSender, err := createDefaultAccount(accountkey.AccountKeyTypePublic)
			assert.NoError(t, err)

			valueMap, _ := genMapForTxTypes(sender, sender, types.TxTypeAccountUpdate)
			valueMap[types.TxValueKeyAccountKey] = newSender.AccKey
			tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, valueMap)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, sender.Keys)
			assert.Equal(t, nil, err)

			receipt, _, err := applyTransaction(t, bcdata, tx)
			assert.NoError(t, err)
			assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
		}
	}

	// AccountStatus를 Active로 설정함
	{
		tx := AccountStatusUpdate(sender, account.AccountStatusActive, signer, nil, t)

		txs := []*types.Transaction{tx}
		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}

		sender.AddNonce()
	}

	// 다시 tx 생성할 수 있음을 확인 함
	{
		// legacy transaction 생성
		{
			var txs types.Transactions

			valueMap, _ := genMapForTxTypes(sender, sender, types.TxTypeLegacyTransaction)
			tx, err := types.NewTransactionWithMap(types.TxTypeLegacyTransaction, valueMap)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, sender.Keys)
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
				t.Fatal(err)
			}
			sender.AddNonce()
		}

		// non-legacy transaction 생성
		{
			min := 0
			max := math.MaxInt32
			transferAmount := big.NewInt(int64(rand.Intn(max-min) + min))
			tx := ValueTransfer(sender, receiver, transferAmount, signer, t)

			txs := []*types.Transaction{tx}
			if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
				t.Fatal(err)
			}
			sender.AddNonce()
		}
	}
}

// 수신인의 계정 상태가 stop일 때 tx를 보낼 수 없는 것을 확인
func TestAccountStatus_To(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	receiver := env.receivers[0]

	// 수신인 계정을 stop 상태로 setAccountStatus 호출
	tx := AccountStatusUpdate(receiver, account.AccountStatusStop, signer, nil, t)

	txs := []*types.Transaction{tx}
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}
	receiver.AddNonce()

	// 수신인에게 value transfer 실행 시 실패 확인
	min := 0
	max := math.MaxInt32
	transferAmount := big.NewInt(int64(rand.Intn(max-min) + min))
	tx = ValueTransfer(sender, receiver, transferAmount, signer, t)

	receipt, _, err := applyTransaction(t, bcdata, tx)
	assert.NoError(t, err)
	assert.Equal(t, types.ReceiptStatusErrStoppedAccountTo, receipt.Status)
}

// AccountStatus가 Stop일 때 tx를 보내면 에러 발생하는 것을 확인 (txpool 에러 발생)
func TestAccountStatus_ErrorMessage(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender

	// 계정을 stop 상태로 setAccountStatus 호출
	tx := AccountStatusUpdate(sender, account.AccountStatusStop, signer, nil, t)

	txs := []*types.Transaction{tx}
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}
	sender.AddNonce()

	// txpool에서 에러 발생하는 것을 확인
	txpool := blockchain.NewTxPool(blockchain.DefaultTxPoolConfig, bcdata.bc.Config(), bcdata.bc)

	min := 0
	max := math.MaxInt32
	transferAmount := big.NewInt(int64(rand.Intn(max-min) + min))
	tx = ValueTransfer(sender, sender, transferAmount, signer, t)

	err := txpool.AddRemote(tx)
	assert.True(t, errors.Is(err, types.ErrAccountStatusStopSender))
}

// klay_getAccount에서 AccountStatus을 보여주는 것을 확인
func TestAccountStatus_klay_account(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	newStatus := account.AccountStatusStop

	// 계정을 stop 상태로 setAccountStatus 호출
	tx := AccountStatusUpdate(sender, account.AccountStatusStop, signer, nil, t)

	txs := []*types.Transaction{tx}
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}
	sender.AddNonce()

	// klay_getAccount 실행하여 설정한 accountstatus 조회
	statedb, err := bcdata.bc.State()
	assert.NoError(t, err)
	acc := statedb.GetAccount(sender.Addr)
	eoa, ok := acc.(account.EOA)
	assert.True(t, ok)
	assert.Equal(t, newStatus, eoa.GetAccountStatus())
}

// role-based account에서 RoleUpdate를 가진 계정만 account status를 설정할 수 있는 것을 확인
func TestAccountStatus_RoleUpdate(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	sender := env.sender
	gasPrice := new(big.Int).SetUint64(25 * params.Ston)

	// generate a role-based key
	prvKeys := genTestKeys(3)
	roleKey := accountkey.NewAccountKeyRoleBasedWithValues(accountkey.AccountKeyRoleBased{
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[0].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[1].PublicKey),
		accountkey.NewAccountKeyPublicWithValue(&prvKeys[2].PublicKey),
	})
	keyRoleUpdate := []*ecdsa.PrivateKey{prvKeys[int(accountkey.RoleAccountUpdate)]}

	// RoleUpdate key로 SetBalanceAccount 동작 수행시 실패하는 것 확인
	{
		tx := AccountStatusUpdate(sender, account.AccountStatusStop, signer, keyRoleUpdate, t)

		_, _, err := applyTransaction(t, bcdata, tx)
		assert.Equal(t, types.ErrInvalidSigSender, err)
	}

	// role-based key를 이용하여 키 등록
	{
		var txs types.Transactions

		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      sender.Nonce,
			types.TxValueKeyFrom:       sender.Addr,
			types.TxValueKeyGasLimit:   gasLimit,
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: roleKey,
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		err = tx.SignWithKeys(signer, sender.Keys)
		assert.Equal(t, nil, err)

		txs = append(txs, tx)

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
		sender.Nonce += 1
	}

	// RoleUpdate key로 SetBalanceAccount 동작 수행시 성공하는 것 확인
	{
		tx := AccountStatusUpdate(sender, account.AccountStatusStop, signer, keyRoleUpdate, t)

		receipt, _, err := applyTransaction(t, bcdata, tx)
		assert.NoError(t, err)
		assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}
}

// 한 블록에 setAccountStatus을 여러번 호출했을 때 마지막 값만 반영되는 것을 확인
func TestAccountStatus_pending_multi_AccountStatus(t *testing.T) {
	env := generateAccountStatusTestEnv(t)
	defer env.bcdata.Shutdown()

	prof := env.prof
	accountMap := env.accountMap
	bcdata := env.bcdata
	signer := env.signer
	receiver := env.receivers[0]
	tryNum := int64(100)

	var txs types.Transactions

	for i := int64(1); i <= tryNum; i++ {
		// account status update (Active)
		{
			tx := AccountStatusUpdate(receiver, account.AccountStatusActive, signer, nil, t)

			txs = append(txs, tx)

			receiver.AddNonce()
		}

		// account status update (Stop)
		{
			tx := AccountStatusUpdate(receiver, account.AccountStatusStop, signer, nil, t)

			txs = append(txs, tx)

			receiver.AddNonce()
		}
	}

	// 블록 생성 (모든 tx 하나의 블록에)
	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}

	// 최종 AccountStatus 값과 Nonce 값 확인
	statedb, err := bcdata.bc.State()
	assert.NoError(t, err)
	acc := statedb.GetAccount(receiver.Addr)
	eoa, ok := acc.(account.EOA)
	assert.True(t, ok)
	assert.Equal(t, account.AccountStatusStop, eoa.GetAccountStatus())
	assert.Equal(t, tryNum*2, int64(eoa.GetNonce()))
}

// 한 블록에 setAccountStatus과 Account Update를 반복했을 때, roleAccountUpdate가 setAccountStatus할 수 있는 것을 확인
func TestAccountStatus_pending_setAccountStatus_accountUpdate(t *testing.T) {
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
