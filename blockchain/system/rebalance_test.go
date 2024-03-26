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

func TestRebalanceTreasuryKIP103(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	rebalanceAddress := common.HexToAddress("0x1030")
	senderKey, _ := crypto.GenerateKey()
	sender := bind.NewKeyedTransactor(senderKey)
	rebalanceTreasury(t,
		sender,
		&params.ChainConfig{
			Kip103CompatibleBlock: big.NewInt(1),
			Kip103ContractAddress: rebalanceAddress,
		},
		rebalanceAddress,
		Kip103MockCode,
	)
}

func TestRebalanceTreasuryKIP160(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	rebalanceAddress := common.HexToAddress("0x1030")
	senderKey, _ := crypto.GenerateKey()
	sender := bind.NewKeyedTransactor(senderKey)
	config := &params.ChainConfig{
		UnitPrice:                25000000000,
		ChainID:                  big.NewInt(1),
		IstanbulCompatibleBlock:  big.NewInt(0),
		LondonCompatibleBlock:    big.NewInt(0),
		EthTxTypeCompatibleBlock: big.NewInt(0),
		MagmaCompatibleBlock:     big.NewInt(0),
		KoreCompatibleBlock:      big.NewInt(0),
		ShanghaiCompatibleBlock:  big.NewInt(0),
		CancunCompatibleBlock:    big.NewInt(0),
		DragonCompatibleBlock:    big.NewInt(1),
		RandaoRegistry: &params.RegistryConfig{
			Records: map[string]common.Address{
				"KIP160": rebalanceAddress,
			},
			Owner: sender.From,
		},
	}
	config.SetDefaults()
	rebalanceTreasury(t, sender, config, rebalanceAddress, Kip160MockCode)
}

func rebalanceTreasury(t *testing.T, sender *bind.TransactOpts, config *params.ChainConfig, rebalanceAddress common.Address, rebalanceCode []byte) {
	var (
		senderAddr = sender.From

		zeroeds = []struct {
			addr    common.Address
			balance *big.Int
		}{
			{common.HexToAddress("0xaa00"), big.NewInt(4_000_000)},
			{common.HexToAddress("0xaa11"), big.NewInt(2_000_000)},
			{common.HexToAddress("0xaa22"), big.NewInt(1_000_000)},
		}

		allocateds = []struct {
			addr    common.Address
			balance *big.Int
		}{
			{common.HexToAddress("0xbb00"), big.NewInt(12_345)},
			{common.HexToAddress("0xbb11"), big.NewInt(0)},
		}

		alloc = blockchain.GenesisAlloc{
			senderAddr:         {Balance: big.NewInt(params.KLAY)},
			rebalanceAddress:   {Code: rebalanceCode, Balance: common.Big0},
			zeroeds[0].addr:    {Balance: zeroeds[0].balance},
			zeroeds[1].addr:    {Balance: zeroeds[1].balance},
			zeroeds[2].addr:    {Balance: zeroeds[2].balance},
			allocateds[0].addr: {Balance: allocateds[0].balance},
			allocateds[1].addr: {Balance: allocateds[1].balance},
		}
	)

	// TODO-system: add a case when zeroeds < allocateds
	testCases := []struct {
		rebalanceBlockNumber *big.Int
		status               uint8
		allocatedAmounts     []*big.Int

		expectedErr             error
		expectZeroedsAmounts    []*big.Int
		expectAllocatedsAmounts []*big.Int
		expectNonce             uint64
		expectBurnt             *big.Int
		expectSuccess           bool
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
			[]*big.Int{zeroeds[0].balance, zeroeds[1].balance, zeroeds[2].balance},
			[]*big.Int{allocateds[0].balance, allocateds[1].balance},
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
			[]*big.Int{zeroeds[0].balance, zeroeds[1].balance, zeroeds[2].balance},
			[]*big.Int{allocateds[0].balance, allocateds[1].balance},
			9,
			big.NewInt(0),
			false,
		},
		{
			// failed case - rebalancing aborted doe to bigger allocation than zeroeds
			big.NewInt(1),
			EnumRebalanceStatus_Registered,
			[]*big.Int{big.NewInt(2_000_000), big.NewInt(5_000_000 + 1)},

			ErrRebalanceBadStatus,
			[]*big.Int{zeroeds[0].balance, zeroeds[1].balance, zeroeds[2].balance},
			[]*big.Int{allocateds[0].balance, allocateds[1].balance},
			9,
			big.NewInt(0),
			false,
		},
	}

	for _, tc := range testCases {
		var (
			db                          = database.NewMemoryDBManager()
			backend                     = backends.NewSimulatedBackendWithDatabase(db, alloc, config)
			chain                       = backend.BlockChain()
			zeroedAddrs, allocatedAddrs []common.Address
		)

		// Deploy TreasuryRebalanceMock contract at block 0 and transit to block 1
		for _, entry := range zeroeds {
			zeroedAddrs = append(zeroedAddrs, entry.addr)
		}
		for _, entry := range allocateds {
			allocatedAddrs = append(allocatedAddrs, entry.addr)
		}

		if chain.Config().DragonCompatibleBlock != nil {
			contract, _ := system_contracts.NewTreasuryRebalanceMockV2Transactor(rebalanceAddress, backend)
			_, err := contract.TestSetAll(sender, zeroedAddrs, allocatedAddrs, tc.allocatedAmounts, tc.rebalanceBlockNumber, tc.status)
			assert.Nil(t, err)
		} else {
			contract, _ := system_contracts.NewTreasuryRebalanceMockTransactor(rebalanceAddress, backend)
			_, err := contract.TestSetAll(sender, zeroedAddrs, allocatedAddrs, tc.allocatedAmounts, tc.rebalanceBlockNumber, tc.status)
			assert.Nil(t, err)
		}

		backend.Commit()
		assert.Equal(t, uint64(1), backend.BlockChain().CurrentBlock().NumberU64())

		// Get state and check post-transition state
		state, err := chain.State()
		assert.Nil(t, err)

		res, err := RebalanceTreasury(state, chain, chain.CurrentHeader())
		assert.Equal(t, tc.expectedErr, err)

		for i, addr := range zeroedAddrs {
			assert.Equal(t, tc.expectZeroedsAmounts[i], state.GetBalance(addr))
		}
		for i, addr := range allocatedAddrs {
			assert.Equal(t, tc.expectAllocatedsAmounts[i], state.GetBalance(addr))
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
