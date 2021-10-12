package MintBurnWithPermission

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/common/math"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultGasLimit = 5000000
	AccountNum      = 5
)

// CheckReceipt can check if the tx receipt has expected status.
func CheckReceipt(b bind.DeployBackend, tx *types.Transaction, duration time.Duration, expectedStatus uint, t *testing.T) {
	timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), duration)
	defer cancelTimeout()

	receipt, err := bind.WaitMined(timeoutContext, b, tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedStatus, receipt.Status)
}

type mintBurnTestENV struct {
	backend          *backends.SimulatedBackend
	minter           *bind.TransactOpts
	burner           *bind.TransactOpts
	burnee           *bind.TransactOpts
	accounts         []*bind.TransactOpts
	mintBurnContract *MintBurnWithPermission
	mintAmount       *big.Int
	burnAmount       *big.Int
}

func generateMintBurnTestEnv(t *testing.T, accountNum uint) *mintBurnTestENV {
	// generate keys
	//adminKey, _ := crypto.GenerateKey()
	//admin := bind.NewKeyedTransactor(adminKey)
	//admin.GasLimit = DefaultGasLimit

	keys := make([]*ecdsa.PrivateKey, accountNum+4)
	accounts := make([]*bind.TransactOpts, accountNum+4)
	for i := uint(0); i < accountNum+3; i++ {
		keys[i], _ = crypto.GenerateKey()
		accounts[i] = bind.NewKeyedTransactor(keys[i])
		accounts[i].GasLimit = DefaultGasLimit
	}

	// generate backend with deployed
	alloc := blockchain.GenesisAlloc{
		vm.PermissionedCallerAddress: {
			Code:    common.FromHex(MintBurnWithPermissionBinRuntime),
			Balance: big.NewInt(0),
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"): common.BytesToHash(accounts[0].From.Bytes()),
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"): common.BytesToHash(accounts[1].From.Bytes()),
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"): common.BytesToHash(accounts[2].From.Bytes()),
			},
		},
	}
	backend := backends.NewSimulatedBackend(alloc)

	// Deploy MintBurn contract
	mintBurnContract, err := NewMintBurnWithPermission(vm.PermissionedCallerAddress, backend)
	assert.NoError(t, err)

	return &mintBurnTestENV{
		backend:          backend,
		minter:           accounts[0],
		burner:           accounts[1],
		burnee:           accounts[2],
		accounts:         accounts[3:],
		mintBurnContract: mintBurnContract,
		mintAmount:       big.NewInt(params.KLAY),
		burnAmount:       big.NewInt(params.Ston),
	}
}

// minter로 등록된 address만 출 가능함
func TestMintBurn_Mint_onlyMinter(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	minter := env.minter
	newMinter := env.accounts[0]
	mintee := env.accounts[1]
	mintAmount := env.mintAmount
	contract := env.mintBurnContract

	// firstminter는 mint 호출 가능함
	{
		tx, err := contract.Mint(minter, mintee.From, mintAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// minter로 등록된 address는 mint 호출 가능함
	{
		tx, err := contract.SetMinter(minter, newMinter.From)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		tx, err = contract.Mint(newMinter, mintee.From, mintAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// minter로 등록되지 않은 address는 mint 호출 불가능
	{
		// mintee address로 mint 호출 불가
		tx, err := contract.Mint(mintee, mintee.From, mintAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

		// 임의의 address로 mint 출호출 불가
		for i := 0; i < 10; i++ {
			key, _ := crypto.GenerateKey()
			account := bind.NewKeyedTransactor(key)
			account.GasLimit = DefaultGasLimit

			tx, err := contract.Mint(account, mintee.From, mintAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
		}
	}
}

// PermissionedCallerAddress로 등록된 contract만 mint 할 수 있는지 확인
func TestMintBurn_Mint_PermissionedCallerAddress(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	minter := env.minter
	mintee := env.accounts[0]
	newContractMinter := env.accounts[1]
	mintAmount := env.mintAmount
	contract := env.mintBurnContract

	// PermissionedCallerAddress로 등록된 contract로 mint 가능
	{
		tx, err := contract.Mint(minter, mintee.From, mintAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// PermissionedCallerAddress로 등록안된 contract로 mint 불가능
	{
		// Deploy
		_, tx, contract, err := DeployMintBurnTest(newContractMinter, backend)
		assert.NoError(t, err)
		backend.Commit()
		assert.Nil(t, bind.CheckWaitMined(backend, tx))

		tx, err = contract.Mint(newContractMinter, mintee.From, mintAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}
}

// mint의 amount가 uint256 범위 내에 있는지 확인함
func TestMintBurn_Mint_AmountRange(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	minter := env.minter
	mintee := env.accounts[0]
	contract := env.mintBurnContract

	// amount가 0일 때 mint가 실패함
	{
		tx, err := contract.Mint(minter, mintee.From, big.NewInt(0))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// amount가 양의 값일 때 mint가 성공함
	{
		mintAmounts := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(100),
			big.NewInt(math.MaxUint32),
			big.NewInt(params.Ston),
			big.NewInt(params.UKLAY),
			big.NewInt(params.MiliKLAY),
			big.NewInt(params.KLAY),
			math.MaxBig256,
		}

		for _, mintAmount := range mintAmounts {
			// mintAmount의 값은 uint256보다 작아야 하기 때문에 mint 전 reset 시켜줌
			tx, err := contract.ResetMintAmount(minter)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			tx, err = contract.Mint(minter, mintee.From, mintAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
		}
	}
}

// mint 시 입력한 amount만큼 balance에 추가 되는지 확인함
func TestMintBurn_Mint_AddBalance(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	minter := env.minter
	mintee := env.accounts[0]
	mintAmount := env.mintAmount
	contract := env.mintBurnContract

	// mint 시 입력한 amount만큼 balance에 추가 되는지 확인함
	{
		previousBalance, err := backend.BalanceAt(context.Background(), mintee.From, nil)
		assert.NoError(t, err)

		tx, err := contract.Mint(minter, mintee.From, mintAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		afterBalance, err := backend.BalanceAt(context.Background(), mintee.From, nil)
		assert.NoError(t, err)
		assert.Equal(t, mintAmount, afterBalance.Sub(afterBalance, previousBalance), "balance change does not match")
	}
}

// ResetMintAmount 실행 후 uint256_max 값을 burnee에게 mint
func mintMax(t *testing.T, env *mintBurnTestENV) {
	tx, err := env.mintBurnContract.ResetMintAmount(env.minter)
	assert.NoError(t, err)
	env.backend.Commit()
	CheckReceipt(env.backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	tx, err = env.mintBurnContract.Mint(env.minter, env.burnee.From, math.MaxBig256)
	assert.NoError(t, err)
	env.backend.Commit()
	CheckReceipt(env.backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
}

// burner로 등록된 address만 출 가능함
func TestMintBurn_Burn_onlyBurner(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	burner := env.burner
	newBurner := env.accounts[0]
	burnee := env.burnee
	burnAmount := env.burnAmount
	contract := env.mintBurnContract

	// Burn하려는 값보다 많은 값을 Mint (Mint한 값이 Burn한 값보다 많아야 함)
	mintMax(t, env)

	// firstburner는 burn 호출 가능함
	{
		tx, err := contract.Burn(burner, burnAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// burner로 등록된 address는 burn 호출 가능함
	{
		tx, err := contract.SetBurner(burner, newBurner.From)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		tx, err = contract.Burn(newBurner, burnAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// minter로 등록되지 않은 address는 mint 호출 불가능
	{
		// burnee address로 mint 호출 불가
		tx, err := contract.Burn(burnee, burnAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

		// 임의의 address로 mint 출호출 불가
		for i := 0; i < 10; i++ {
			key, _ := crypto.GenerateKey()
			account := bind.NewKeyedTransactor(key)
			account.GasLimit = DefaultGasLimit

			tx, err := contract.Burn(account, burnAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
		}
	}
}

// PermissionedCallerAddress로 등록된 contract만 burn 할 수 있는지 확인
func TestMintBurn_Burn_PermissionedCallerAddress(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	burner := env.burner
	burnee := env.burnee
	burnAmount := env.burnAmount
	contract := env.mintBurnContract

	// Burn하려는 값보다 많은 값을 Mint (Mint한 값이 Burn한 값보다 많아야 함)
	mintMax(t, env)

	// PermissionedCallerAddress로 등록된 contract로 burn 가능
	{
		tx, err := contract.Burn(burner, burnAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// PermissionedCallerAddress로 등록안된 contract로 burn 불가능
	{
		// Deploy
		_, tx, contract, err := DeployMintBurnTest(burner, backend)
		assert.NoError(t, err)
		backend.Commit()
		assert.Nil(t, bind.CheckWaitMined(backend, tx))

		tx, err = contract.Burn(burner, burnee.From, burnAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}
}

// burn의 amount가 uint256 범위 내에 있는지 확인함
func TestMintBurn_Burn_AmountRange(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	burner := env.burner
	contract := env.mintBurnContract

	// amount가 0일 때 burn가 실패함
	{
		tx, err := contract.Burn(burner, big.NewInt(0))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// amount가 양의 값일 때 burn가 성공함
	{
		burnAmounts := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(100),
			big.NewInt(math.MaxUint32),
			big.NewInt(params.Ston),
			big.NewInt(params.UKLAY),
			big.NewInt(params.MiliKLAY),
			big.NewInt(params.KLAY),
			math.MaxBig256,
		}

		for _, burnAmount := range burnAmounts {
			// Burn하려는 값보다 많은 값을 Mint (Mint한 값이 Burn한 값보다 많아야 함)
			mintMax(t, env)

			// burnAmount의 값은 uint256보다 작아야 하기 때문에 burn 전 reset 시켜줌
			tx, err := contract.ResetBurnAmount(burner)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			// Burn
			tx, err = contract.Burn(burner, burnAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
		}
	}
}

// burn 시 입력한 amount만큼 balance에 감소 되는지 확인함
func TestMintBurn_Burn_SubBalance(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	burner := env.burner
	burnee := env.burnee
	burnAmount := env.burnAmount
	contract := env.mintBurnContract

	// Burn하려는 값보다 많은 값을 Mint (Mint한 값이 Burn한 값보다 많아야 함)
	mintMax(t, env)

	// burn 시 입력한 amount만큼 balance에 감소 되는지 확인함
	{
		previousBalance, err := backend.BalanceAt(context.Background(), burnee.From, nil)
		assert.NoError(t, err)

		tx, err := contract.Burn(burner, burnAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		afterBalance, err := backend.BalanceAt(context.Background(), burnee.From, nil)
		assert.NoError(t, err)
		assert.Equal(t, burnAmount, previousBalance.Sub(previousBalance, afterBalance), "balance change does not match")
	}
}

// burn 시 보유량 보다 큰 값을 입력하면 실패하는지 확인함
func TestMintBurn_Burn_AmountSmallerThanBalance(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	burner := env.burner
	burnee := env.burnee
	mintAmount := env.mintAmount
	burnAmount := env.burnAmount
	contract := env.mintBurnContract

	// Burn하려는 값보다 많은 값을 Mint (Mint한 값이 Burn한 값보다 많아야 함)
	tx, err := env.mintBurnContract.Mint(env.minter, env.burnee.From, mintAmount)
	assert.NoError(t, err)
	env.backend.Commit()
	CheckReceipt(env.backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	// amount가 balance보다 크면 burn 실패함
	{
		previousBalance, err := backend.BalanceAt(context.Background(), burnee.From, nil)
		assert.NoError(t, err)

		tx, err := contract.Burn(burner, previousBalance.Add(previousBalance, big.NewInt(1)))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// amount가 balance보다 작으 burn 성공함
	{
		tx, err := contract.Burn(burner, burnAmount)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}
	// amount가 balance 같으면 burn 성공함
	{
		previousBalance, err := backend.BalanceAt(context.Background(), burnee.From, nil)
		assert.NoError(t, err)

		tx, err = contract.Burn(burner, previousBalance)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}
}

// minter가 정상적으로 설정되는지 확인함
func TestMintBurn_Minter(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	minter := env.minter
	newMinter := env.accounts[0]
	contract := env.mintBurnContract

	// genesis.json으로 설정한 minter 주소가 minter로 등록되었음을 확인함
	getMinter, err := contract.GetMinter(nil)
	assert.NoError(t, err)
	assert.Equal(t, minter.From, getMinter)

	// setMinter를 호출하여 다른 minter로 주소를 바뀌었는지 확인함
	tx, err := contract.SetMinter(minter, newMinter.From)
	assert.NoError(t, err)
	env.backend.Commit()
	CheckReceipt(env.backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	getMinter, err = contract.GetMinter(nil)
	assert.NoError(t, err)
	assert.Equal(t, newMinter.From, getMinter)
}

// burner가 정상적으로 설정되는지 확인함
func TestMintBurn_Burner(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	burner := env.burner
	newBurner := env.accounts[0]
	contract := env.mintBurnContract

	// genesis.json으로 설정한 burner 주소가 burner로 등록되었음을 확인함
	getBurner, err := contract.GetBurner(nil)
	assert.NoError(t, err)
	assert.Equal(t, burner.From, getBurner)

	// setMinter를 호출하여 다른 burner로 주소를 바뀌었는지 확인함
	tx, err := contract.SetBurner(burner, newBurner.From)
	assert.NoError(t, err)
	env.backend.Commit()
	CheckReceipt(env.backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	getBurner, err = contract.GetBurner(nil)
	assert.NoError(t, err)
	assert.Equal(t, newBurner.From, getBurner)
}

// burnee가 정상적으로 설정되는지 확인함
func TestMintBurn_Burnee(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	burner := env.burner
	burnee := env.burnee
	newBurnee := env.accounts[0]
	contract := env.mintBurnContract

	// genesis.json으로 설정한 burnee 주소가 burnee로 등록되었음을 확인함
	getBurnee, err := contract.GetBurnee(nil)
	assert.NoError(t, err)
	assert.Equal(t, burnee.From, getBurnee)

	// setMinter를 호출하여 다른 burnee로 주소를 바뀌었는지 확인함
	tx, err := contract.SetBurnee(burner, newBurnee.From)
	assert.NoError(t, err)
	env.backend.Commit()
	CheckReceipt(env.backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

	getBurnee, err = contract.GetBurnee(nil)
	assert.NoError(t, err)
	assert.Equal(t, newBurnee.From, getBurnee)
}

// resetMintAmount가 minted 값을 초기화 하는지 확인함
func TestMintBurn_resetMintAmount(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	minter := env.minter
	mintee := env.accounts[0]
	contract := env.mintBurnContract

	// resetMintAmount 호출 후에 minted 값이 0이 되는 것을 확인함
	{
		mintAmounts := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(100),
			big.NewInt(math.MaxUint32),
			big.NewInt(params.Ston),
			big.NewInt(params.UKLAY),
			big.NewInt(params.MiliKLAY),
			big.NewInt(params.KLAY),
			math.MaxBig256,
		}

		for _, mintAmount := range mintAmounts {
			// mint
			tx, err := contract.Mint(minter, mintee.From, mintAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			// reset mint amount
			tx, err = contract.ResetMintAmount(minter)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			// check if minted is equal to 0
			minted, err := contract.Minted(nil)
			assert.NoError(t, err)
			assert.True(t, big.NewInt(0).Cmp(minted) == 0, "mintAmount should be reset to 0")
		}
	}
}

// resetMintAmount가 burnt 값을 초기화 하는지 확인함
func TestMintBurn_resetBurnAmount(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	burner := env.burner
	contract := env.mintBurnContract

	// resetBurnAmount 호출 후에 burnt 값이 0이 되는 것을 확인함
	{
		burnAmounts := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(100),
			big.NewInt(math.MaxUint32),
			big.NewInt(params.Ston),
			big.NewInt(params.UKLAY),
			big.NewInt(params.MiliKLAY),
			big.NewInt(params.KLAY),
			math.MaxBig256,
		}

		for _, burnAmount := range burnAmounts {
			// mintMax value so that there will be no problem
			mintMax(t, env)

			// mint
			tx, err := contract.Burn(burner, burnAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			// reset mint amount
			tx, err = contract.ResetBurnAmount(burner)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			// check if minted is equal to 0
			burnt, err := contract.Burnt(nil)
			assert.NoError(t, err)
			assert.True(t, big.NewInt(0).Cmp(burnt) == 0, "burnAmount should be reset to 0")
		}
	}
}

// mintAmount가 mint한 값의 합인지 확인
func TestMintBurn_minted(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	minter := env.minter
	mintee := env.burnee
	contract := env.mintBurnContract

	// mint한 amount의 합인지 확인
	{
		mintAmounts := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(100),
			big.NewInt(math.MaxUint32),
			big.NewInt(params.Peb),
			big.NewInt(params.Ston),
			big.NewInt(params.KLAY),
		}

		totalAmount := big.NewInt(0)
		for _, mintAmount := range mintAmounts {
			tx, err := contract.Mint(minter, mintee.From, mintAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			totalAmount = totalAmount.Add(totalAmount, mintAmount)
		}

		minted, err := contract.Minted(nil)
		assert.NoError(t, err)
		assert.Equal(t, totalAmount, minted, "balance change does not match")
	}
}

// num이 uint256 범위 내에 있는지 확인
func checkUint256Range(t *testing.T, num *big.Int) {
	assert.True(t, big.NewInt(0).Cmp(num) <= 0, "balance change does not match")
	assert.True(t, num.Cmp(math.MaxBig256) <= 0, "balance change does not match")
}

// mint 수행 중 overflow가 발생하지 않으며, minted의 값이 uint256 범위 내에 있는 것을 확인
func TestMintBurn_minted_overflow(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	minter := env.minter
	mintee := env.burnee
	contract := env.mintBurnContract

	// mint (1) uint256_max and (2) uint256_max
	{
		// reset mint amount
		tx, err := contract.ResetMintAmount(minter)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// mint uint256 max
		tx, err = contract.Mint(minter, mintee.From, math.MaxBig256)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err := contract.Minted(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)

		// mint uint256 max
		tx, err = contract.Mint(minter, mintee.From, math.MaxBig256)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

		minted, err = contract.Minted(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)
	}

	// mint (1) uint256_max and (2) 1
	{
		// reset mint amount
		tx, err := contract.ResetMintAmount(minter)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// mint uint256 max
		tx, err = contract.Mint(minter, mintee.From, math.MaxBig256)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err := contract.Minted(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)

		// mint uint256 max
		tx, err = contract.Mint(minter, mintee.From, big.NewInt(1))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

		minted, err = contract.Minted(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)
	}

	// mint (1) 1 and (2) uint256_max
	{
		// reset mint amount
		tx, err := contract.ResetMintAmount(minter)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// mint uint256 max
		tx, err = contract.Mint(minter, mintee.From, big.NewInt(1))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err := contract.Minted(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)

		// mint uint256 max
		tx, err = contract.Mint(minter, mintee.From, math.MaxBig256)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

		minted, err = contract.Minted(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)
	}

	// mint (1) uint256_max / 2 and (2) uint256_max / 2
	{
		// reset mint amount
		tx, err := contract.ResetMintAmount(minter)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// mint uint256 max
		tx, err = contract.Mint(minter, mintee.From, big.NewInt(0).Div(math.MaxBig256, big.NewInt(2)))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err := contract.Minted(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)

		// mint uint256 max
		tx, err = contract.Mint(minter, mintee.From, big.NewInt(0).Div(math.MaxBig256, big.NewInt(2)))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err = contract.Minted(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)
	}
}

// burnAmount가 burn한 값의 합인지 확인
func TestMintBurn_burnt(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	burner := env.burner
	contract := env.mintBurnContract

	mintMax(t, env)

	// burn한 amount의 합인지 확인
	{
		burnAmounts := []*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(100),
			big.NewInt(params.Peb),
			big.NewInt(params.Ston),
		}

		totalAmount := big.NewInt(0)
		for _, burnAmount := range burnAmounts {
			tx, err := contract.Burn(burner, burnAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			totalAmount = totalAmount.Add(totalAmount, burnAmount)
		}

		burned, err := contract.Burnt(nil)
		assert.NoError(t, err)
		assert.Equal(t, totalAmount, burned, "balance change does not match")
	}
}

// burn 수행 중 overflow가 발생하지 않으며, burnt의 값이 uint256 범위 내에 있는 것을 확인
func TestMintBurn_burnt_overflow(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	burner := env.burner
	contract := env.mintBurnContract

	// burn (1) uint256_max and (2) uint256_max
	{
		// mint max amount
		mintMax(t, env)

		// reset burn amount
		tx, err := contract.ResetBurnAmount(burner)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// mint uint256 max
		tx, err = contract.Burn(burner, math.MaxBig256)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err := contract.Burnt(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)

		// mint uint256 max
		tx, err = contract.Burn(burner, math.MaxBig256)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

		minted, err = contract.Burnt(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)
	}

	// burn (1) uint256_max and (2) 1
	{
		// mint max amount
		mintMax(t, env)

		// reset burn amount
		tx, err := contract.ResetBurnAmount(burner)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// mint uint256 max
		tx, err = contract.Burn(burner, math.MaxBig256)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err := contract.Burnt(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)

		// mint uint256 max
		tx, err = contract.Burn(burner, big.NewInt(1))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

		minted, err = contract.Burnt(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)
	}

	// burn (1) 1 and (2) uint256_max
	{
		// mint max amount
		mintMax(t, env)

		// reset burn amount
		tx, err := contract.ResetBurnAmount(burner)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// mint uint256 max
		tx, err = contract.Burn(burner, big.NewInt(1))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err := contract.Burnt(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)

		// mint uint256 max
		tx, err = contract.Burn(burner, math.MaxBig256)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)

		minted, err = contract.Burnt(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)
	}

	// burn (1) uint256_max / 2 and (2) uint256_max / 2
	{
		// mint max amount
		mintMax(t, env)

		// reset burn amount
		tx, err := contract.ResetBurnAmount(burner)
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		// mint uint256 max
		tx, err = contract.Burn(burner, big.NewInt(0).Div(math.MaxBig256, big.NewInt(2)))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err := contract.Burnt(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)

		// mint uint256 max
		tx, err = contract.Burn(burner, big.NewInt(0).Div(math.MaxBig256, big.NewInt(2)))
		assert.NoError(t, err)
		backend.Commit()
		CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		minted, err = contract.Burnt(nil)
		assert.NoError(t, err)
		checkUint256Range(t, minted)
	}
}

// totalSupply의 값이 Minted - Burnt와 같은지 확인
// totalSupply의 값이 (Mint의 합) - (Burnt의 합)과 같은지 확인
func TestMintBurn_totalSupply(t *testing.T) {
	env := generateMintBurnTestEnv(t, AccountNum)
	defer env.backend.Close()

	backend := env.backend
	minter := env.minter
	burner := env.burner
	mintee := env.burnee
	contract := env.mintBurnContract

	mintAmounts := []*big.Int{
		big.NewInt(1000 * params.Ston),
		big.NewInt(123 * params.Ston),
		big.NewInt(456 * params.Peb),
		big.NewInt(978 * params.Peb),
	}

	burnAmounts := []*big.Int{
		big.NewInt(147 * params.Ston),
		big.NewInt(258 * params.Peb),
		big.NewInt(369 * params.Peb),
	}

	// Mint and Burn
	totalMint := big.NewInt(0)
	totalBurn := big.NewInt(0)
	previousBalance, err := backend.BalanceAt(context.Background(), mintee.From, nil)
	assert.NoError(t, err)
	{

		for _, mintAmount := range mintAmounts {
			tx, err := contract.Mint(minter, mintee.From, mintAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			totalMint = totalMint.Add(totalMint, mintAmount)
		}

		for _, burnAmount := range burnAmounts {
			tx, err := contract.Burn(burner, burnAmount)
			assert.NoError(t, err)
			backend.Commit()
			CheckReceipt(backend, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			totalBurn = totalBurn.Add(totalBurn, burnAmount)
		}
	}
	afterBalance, err := backend.BalanceAt(context.Background(), mintee.From, nil)
	assert.NoError(t, err)

	// 호출 결과가 (mint한 amount) - (burn한 amount)인 지 확인
	assert.Equal(t, big.NewInt(0).Sub(totalMint, totalBurn), big.NewInt(0).Sub(afterBalance, previousBalance), "balance change does not match")

	// 호출 결과가 minted - burnt인 지 확인
	minted, err := contract.Minted(nil)
	assert.NoError(t, err)
	burned, err := contract.Burnt(nil)
	assert.NoError(t, err)
	assert.Equal(t, totalMint.Sub(totalMint, totalBurn), minted.Sub(minted, burned), "balance change does not match")
}
