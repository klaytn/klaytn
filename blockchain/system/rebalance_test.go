package system

import (
	"encoding/json"
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

func TestRebalanceTreasury(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	var (
		rebalanceAddress = common.HexToAddress("0x1030")
		senderKey, _     = crypto.GenerateKey()
		sender           = bind.NewKeyedTransactor(senderKey)
		senderAddr       = sender.From

		retirees = []struct {
			addr    common.Address
			balance *big.Int
		}{
			{common.HexToAddress("0xaa00"), big.NewInt(4_000_000)},
			{common.HexToAddress("0xaa11"), big.NewInt(2_000_000)},
			{common.HexToAddress("0xaa22"), big.NewInt(1_000_000)},
		}

		newbies = []struct {
			addr    common.Address
			balance *big.Int
		}{
			{common.HexToAddress("0xbb00"), big.NewInt(12_345)},
			{common.HexToAddress("0xbb11"), big.NewInt(0)},
		}

		config = &params.ChainConfig{
			Kip103CompatibleBlock: big.NewInt(1),
			Kip103ContractAddress: rebalanceAddress,
		}

		alloc = blockchain.GenesisAlloc{
			senderAddr:       {Balance: big.NewInt(params.KLAY)},
			rebalanceAddress: {Code: Kip103MockCode, Balance: common.Big0},
			retirees[0].addr: {Balance: retirees[0].balance},
			retirees[1].addr: {Balance: retirees[1].balance},
			retirees[2].addr: {Balance: retirees[2].balance},
			newbies[0].addr:  {Balance: newbies[0].balance},
			newbies[1].addr:  {Balance: newbies[1].balance},
		}
	)

	// TODO-system: add a case when retirees < newbies
	testCases := []struct {
		rebalanceBlockNumber *big.Int
		status               uint8
		newbieAmounts        []*big.Int

		expectedErr           error
		expectRetireesAmounts []*big.Int
		expectNewbiesAmounts  []*big.Int
		expectNonce           uint64
		expectBurnt           *big.Int
		expectSuccess         bool
	}{
		{
			// normal case
			big.NewInt(1),
			EnumRebalanceStatus_Approved,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000)},

			nil,
			[]*big.Int{big.NewInt(0), big.NewInt(0), big.NewInt(0)},
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000)},
			10,
			big.NewInt(12345),
			true,
		},
		{
			// failed case - rebalancing aborted due to wrong rebalanceBlockNumber
			big.NewInt(2),
			EnumRebalanceStatus_Approved,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000)},

			ErrRebalanceIncorrectBlock,
			[]*big.Int{retirees[0].balance, retirees[1].balance, retirees[2].balance},
			[]*big.Int{newbies[0].balance, newbies[1].balance},
			8,
			big.NewInt(0),
			false,
		},
		{
			// failed case - rebalancing aborted doe to wrong state
			big.NewInt(1),
			EnumRebalanceStatus_Registered,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000)},

			ErrRebalanceBadStatus,
			[]*big.Int{retirees[0].balance, retirees[1].balance, retirees[2].balance},
			[]*big.Int{newbies[0].balance, newbies[1].balance},
			9,
			big.NewInt(0),
			false,
		},
		{
			// failed case - rebalancing aborted doe to bigger allocation than retirees
			big.NewInt(1),
			EnumRebalanceStatus_Registered,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000 + 1)},

			ErrRebalanceBadStatus,
			[]*big.Int{retirees[0].balance, retirees[1].balance, retirees[2].balance},
			[]*big.Int{newbies[0].balance, newbies[1].balance},
			9,
			big.NewInt(0),
			false,
		},
	}

	for _, tc := range testCases {
		var (
			db                        = database.NewMemoryDBManager()
			backend                   = backends.NewSimulatedBackendWithDatabase(db, alloc, config)
			chain                     = backend.BlockChain()
			contract, _               = system_contracts.NewTreasuryRebalanceMockTransactor(rebalanceAddress, backend)
			retireeAddrs, newbieAddrs []common.Address
		)

		// Deploy TreasuryRebalanceMock contract at block 0 and transit to block 1
		for _, entry := range retirees {
			retireeAddrs = append(retireeAddrs, entry.addr)
		}
		for _, entry := range newbies {
			newbieAddrs = append(newbieAddrs, entry.addr)
		}

		_, err := contract.TestSetAll(sender, retireeAddrs, newbieAddrs, tc.newbieAmounts, tc.rebalanceBlockNumber, tc.status)
		assert.Nil(t, err)

		backend.Commit()
		assert.Equal(t, uint64(1), backend.BlockChain().CurrentBlock().NumberU64())

		// Get state and check post-transition state
		state, err := chain.State()
		assert.Nil(t, err)

		header := chain.CurrentHeader()
		c := &Kip103ContractCaller{state, chain, header}
		res, err := RebalanceTreasury(state, chain, header, c)
		assert.Equal(t, tc.expectedErr, err)

		for i, addr := range retireeAddrs {
			assert.Equal(t, tc.expectRetireesAmounts[i], state.GetBalance(addr))
		}
		for i, addr := range newbieAddrs {
			assert.Equal(t, tc.expectNewbiesAmounts[i], state.GetBalance(addr))
		}

		if chain.Config().DragonCompatibleBlock != nil {
			assert.Equal(t, uint64(0x0), state.GetNonce(common.HexToAddress("0x0")))
		} else {
			assert.Equal(t, tc.expectNonce, state.GetNonce(common.HexToAddress("0x0")))
		}
		assert.Equal(t, tc.expectBurnt, res.Burnt)
		assert.Equal(t, tc.expectSuccess, res.Success)

		memo, _ := json.Marshal(res)
		t.Log(string(memo))
	}
}
