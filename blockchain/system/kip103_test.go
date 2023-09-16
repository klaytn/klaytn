package system

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

func TestKip103(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

	toPeb := func(klay int64) *big.Int {
		return new(big.Int).Mul(big.NewInt(int64(klay)), big.NewInt(params.KLAY))
	}
	toPebArr := func(klay []int64) []*big.Int {
		result := make([]*big.Int, len(klay))
		for i, v := range klay {
			result[i] = toPeb(v)
		}
		return result
	}

	var ( // Genesis config
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)
		senderAddr   = sender.From
		kip103addr   = common.HexToAddress("0x1030")

		r1addr = common.HexToAddress("0xaa11")
		r2addr = common.HexToAddress("0xaa22")
		n1addr = common.HexToAddress("0xbb11")
		n2addr = common.HexToAddress("0xbb22")

		config = &params.ChainConfig{
			Kip103CompatibleBlock: big.NewInt(1),
			Kip103ContractAddress: kip103addr,
		}
		alloc = blockchain.GenesisAlloc{
			senderAddr: {Balance: toPeb(1)},
			kip103addr: {Code: Kip103MockCode, Balance: common.Big0},
			r1addr:     {Balance: toPeb(3_000_000)},
			r2addr:     {Balance: toPeb(4_000_000)},
			n1addr:     {Balance: toPeb(123)},
			n2addr:     {Balance: toPeb(0)},
		}
	)

	testcases := []struct {
		retirees     []common.Address
		newbies      []common.Address
		amounts      []int64 // in KLAY
		rebalanceNum int64
		status       uint8

		expectErr             error
		expectRetireesAmounts []int64
		expectNewbiesAmounts  []int64
		expectBurnt           int64
		// Because of the way Kip103ContractCaller is implemented, the nonce of
		// the null account (0x0) increses by the number of contract call.
		expectNonce uint64 // 5 + len(retirees) + len(newbies)
	}{
		{ // Normal case (3m + 4m -> 2m + 5m)
			[]common.Address{r1addr, r2addr},
			[]common.Address{n1addr, n2addr},
			[]int64{2_000_000, 5_000_000},
			1,
			EnumKip103Status_Approved,

			nil,
			[]int64{0, 0},
			[]int64{2_000_000, 5_000_000},
			123,
			9,
		},
		{ // Aborted due to block number
			[]common.Address{r1addr, r2addr},
			[]common.Address{n1addr, n2addr},
			[]int64{2_000_000, 5_000_000},
			2, // unequal to current block
			EnumKip103Status_Approved,

			ErrKip103IncorrectBlock,
			[]int64{3_000_000, 4_000_000}, // unchanged
			[]int64{123, 0},               // unchanged
			0,
			7,
		},
		{ // Aborted due to state
			[]common.Address{r1addr, r2addr},
			[]common.Address{n1addr, n2addr},
			[]int64{2_000_000, 5_000_000},
			1,
			EnumKip103Status_Registered, // bad state

			ErrKip103BadStatus,
			[]int64{3_000_000, 4_000_000}, // unchanged
			[]int64{123, 0},               // unchanged
			0,
			8,
		},
		{ // Not enough balance
			[]common.Address{r1addr, r2addr},
			[]common.Address{n1addr, n2addr},
			[]int64{2_000_000, 5_000_000 + 1},
			1,
			EnumKip103Status_Approved,

			ErrKip103NotEnoughBalance,
			[]int64{3_000_000, 4_000_000}, // unchanged
			[]int64{123, 0},               // unchanged
			0,
			9,
		},
	}

	for tcid, tc := range testcases {
		var (
			db          = database.NewMemoryDBManager()
			backend     = backends.NewSimulatedBackendWithDatabase(db, alloc, config)
			chain       = backend.BlockChain()
			contract, _ = system_contracts.NewTreasuryRebalanceMockTransactor(kip103addr, backend)
		)

		contract.TestSetAll(sender,
			tc.retirees,
			tc.newbies,
			toPebArr(tc.amounts),
			big.NewInt(tc.rebalanceNum),
			tc.status)
		backend.Commit()
		t.Log("blocknum", backend.BlockChain().CurrentBlock().NumberU64())

		state, _ := chain.State()
		header := chain.CurrentHeader()
		c := NewKip103ContractCaller(state, chain, chain.CurrentHeader())
		res, err := RebalanceTreasury(state, chain, header, c)

		// Check error
		if tc.expectErr != nil {
			assert.ErrorIs(t, err, tc.expectErr, tcid)
			assert.False(t, res.Success, tcid)
		} else {
			assert.Nil(t, err, tcid)
			assert.True(t, res.Success, tcid)
		}

		// Check post-transition state
		for i, addr := range tc.retirees {
			assert.Equal(t, toPeb(tc.expectRetireesAmounts[i]), state.GetBalance(addr), tcid)
		}
		for i, addr := range tc.newbies {
			assert.Equal(t, toPeb(tc.expectNewbiesAmounts[i]), state.GetBalance(addr), tcid)
		}
		assert.Equal(t, tc.expectNonce, state.GetNonce(common.HexToAddress("0x0")), tcid)

		// Check memo
		for _, addr := range tc.retirees { // Retired balances before transition
			assert.Equal(t, alloc[addr].Balance, res.Retired[addr], tcid)
		}
		for i, addr := range tc.newbies { // Allocated newbie balances
			assert.Equal(t, toPeb(tc.amounts[i]), res.Newbie[addr], tcid)
		}
		assert.Equal(t, toPeb(tc.expectBurnt), res.Burnt)
	}
}
