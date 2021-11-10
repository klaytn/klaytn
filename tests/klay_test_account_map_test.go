// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
package tests

import (
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/reward/contract"
	"github.com/klaytn/klaytn/crypto"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// AddressBalanceMap
////////////////////////////////////////////////////////////////////////////////
type AccountInfo struct {
	balance       *big.Int
	nonce         uint64
	balanceLimit  *big.Int
	accountStatus account.AccountStatus
}

type AccountMap struct {
	m          map[common.Address]*AccountInfo
	rewardbase common.Address
}

func NewAccountMap() *AccountMap {
	return &AccountMap{
		m: make(map[common.Address]*AccountInfo),
	}
}

func (a *AccountMap) AddBalance(addr common.Address, v *big.Int) {
	if acc, ok := a.m[addr]; ok {
		acc.balance.Add(acc.balance, v)
	} else {
		// create an account
		a.Set(addr, v, 0)
	}
}

func (a *AccountMap) SubBalance(addr common.Address, v *big.Int) error {
	if acc, ok := a.m[addr]; ok {
		acc.balance.Sub(acc.balance, v)
	} else {
		return fmt.Errorf("trying to subtract balance from an uninitiailzed address (%s)", addr.Hex())
	}

	return nil
}

func (a *AccountMap) GetNonce(addr common.Address) uint64 {
	if acc, ok := a.m[addr]; ok {
		return acc.nonce
	}
	// 'StateDB.GetNonce' returns 0 when the address doesn't exist
	return 0
}

func (a *AccountMap) IncNonce(addr common.Address) {
	if acc, ok := a.m[addr]; ok {
		acc.nonce++
	}
}

func (a *AccountMap) Set(addr common.Address, v *big.Int, nonce uint64) {
	a.m[addr] = &AccountInfo{new(big.Int).Set(v), nonce, account.GetInitialBalanceLimit(), account.AccountStatusActive}
}

func (a *AccountMap) Initialize(bcdata *BCData) error {
	statedb, err := bcdata.bc.State()
	if err != nil {
		return err
	}

	// NOTE-Klaytn-Issue973 Developing Klaytn token economy
	// Add predefined accounts related to reward mechanism
	rewardContractAddr := common.HexToAddress(contract.RewardContractAddress)
	kirContractAddr := common.HexToAddress(contract.KIRContractAddress)
	pocContractAddr := common.HexToAddress(contract.PoCContractAddress)
	addrs := append(bcdata.addrs, &rewardContractAddr, &kirContractAddr, &pocContractAddr)

	for _, addr := range addrs {
		a.Set(*addr, statedb.GetBalance(*addr), statedb.GetNonce(*addr))
	}

	a.rewardbase = *bcdata.addrs[0]

	return nil
}

func (a *AccountMap) Update(txs types.Transactions, signer types.Signer, picker types.AccountKeyPicker, currentBlockNumber uint64) error {
	for _, tx := range txs {
		to := tx.To()
		v := tx.Value()

		gasFrom, err := tx.ValidateSender(signer, picker, currentBlockNumber)
		if err != nil {
			return err
		}
		from := tx.ValidatedSender()

		gasFeePayer := uint64(0)
		feePayer := from
		if tx.IsFeeDelegatedTransaction() {
			gasFeePayer, err = tx.ValidateFeePayer(signer, picker, currentBlockNumber)
			if err != nil {
				return err
			}
			feePayer = tx.ValidatedFeePayer()
		}

		// check if from/to account is in active state
		if err := a.validateAccountStatus(tx); err != nil {
			return err
		}

		// create a new address for deployed contract
		if to == nil && tx.Type() != types.TxTypeBalanceLimitUpdate {
			nonce := a.GetNonce(from)
			addr := crypto.CreateAddress(from, nonce)
			to = &addr
		}
		// Fail if we're trying to transfer more than the available balance
		if !a.canTransfer(from, v) {
			return vm.ErrInsufficientBalance // TODO-Klaytn-Issue615
		}
		// Fail if we're trying to transfer receiver's balance limit
		// This only checks balance limit of EOA. Smart cotracts doesn't have balance limit.
		if tx.Type() != types.TxTypeBalanceLimitUpdate {
			if !a.canBeTransferred(*to, v) {
				return vm.ErrExceedBalanceLimit
			}
			a.AddBalance(*to, v)
			a.SubBalance(from, v)
		}

		// TODO-Klaytn: This gas fee calculation is correct only if the transaction is a value transfer transaction.
		// Calculate the correct transaction fee by checking the corresponding receipt.
		intrinsicGas, err := tx.IntrinsicGas(currentBlockNumber)
		if err != nil {
			return err
		}

		intrinsicGas += gasFrom + gasFeePayer

		fee := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(intrinsicGas))

		feeRatio, ok := tx.FeeRatio()
		if ok {
			feeByFeePayer, feeBySender := types.CalcFeeWithRatio(feeRatio, fee)
			a.SubBalance(feePayer, feeByFeePayer)
			a.SubBalance(from, feeBySender)
		} else {
			a.SubBalance(feePayer, fee)
		}

		a.AddBalance(a.rewardbase, fee)

		a.IncNonce(from)

		if tx.IsLegacyTransaction() && tx.To() == nil {
			a.IncNonce(*to)
		}

		if tx.Type() == types.TxTypeSmartContractDeploy || tx.Type() == types.TxTypeFeeDelegatedSmartContractDeploy || tx.Type() == types.TxTypeFeeDelegatedSmartContractDeployWithRatio {
			a.IncNonce(*to)
		}
		if t, ok := tx.GetTxInternalData().(*types.TxInternalDataBalanceLimitUpdate); ok {
			a.m[from].balanceLimit = new(big.Int).Set(t.BalanceLimit)
		}
		if t, ok := tx.GetTxInternalData().(*types.TxInternalDataAccountStatusUpdate); ok {
			a.m[from].accountStatus = t.AccountStatus
		}
	}

	return nil
}

func (a *AccountMap) Verify(statedb *state.StateDB) error {
	for addr, acc := range a.m {
		if acc.nonce != statedb.GetNonce(addr) {
			return errors.New(fmt.Sprintf("[%s] nonce is different!! statedb(%d) != accountMap(%d).\n",
				addr.Hex(), statedb.GetNonce(addr), acc.nonce))
		}

		if acc.balance.Cmp(statedb.GetBalance(addr)) != 0 {
			return errors.New(fmt.Sprintf("[%s] balance is different!! statedb(%s) != accountMap(%s).\n",
				addr.Hex(), statedb.GetBalance(addr).String(), acc.balance.String()))
		}
	}

	return nil
}

func (a *AccountMap) validateAccountStatus(tx *types.Transaction) error {
	accountStatusGetter := func(addr common.Address) (account.AccountStatus, error) {
		if accountInfo, ok := a.m[addr]; ok {
			return accountInfo.accountStatus, nil
		}
		return account.AccountStatusUndefined, account.ErrNotEOA
	}

	// check if "from" is an active account
	if !blockchain.IsActiveAccount(accountStatusGetter, tx.ValidatedSender()) &&
		// RoleAccountUpdate can make transaction on any account status
		tx.GetRoleTypeForValidation() != accountkey.RoleAccountUpdate {
		return errors.Wrap(types.ErrAccountStatusStopSender, "stopped account is "+tx.ValidatedSender().String())
	}

	// check if "to" is an active account
	if to := tx.To(); to != nil && !blockchain.IsActiveAccount(accountStatusGetter, *to) {
		return errors.Wrap(types.ErrAccountStatusStopReceiver, "stopped account is "+to.String())
	}
	return nil
}

func (a *AccountMap) canTransfer(addr common.Address, amount *big.Int) bool {
	return a.m[addr].balance.Cmp(amount) >= 0
}

func (a *AccountMap) canBeTransferred(receiver common.Address, amount *big.Int) bool {
	balanceLimitGetter := func(addr common.Address) (*big.Int, error) {
		if accountInfo, ok := a.m[addr]; ok {
			return accountInfo.balanceLimit, nil
		}
		return nil, account.ErrNotEOA
	}
	balanceGetter := func(addr common.Address) *big.Int {
		if accountInfo, ok := a.m[addr]; ok {
			return accountInfo.balance
		}
		return new(big.Int)
	}
	return blockchain.CanBeTransferred(balanceLimitGetter, balanceGetter, receiver, amount)
}
