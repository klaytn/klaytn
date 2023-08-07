package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/contracts/kip103"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

type mockKip103ContractCaller struct {
	abi        abi.ABI
	funcSigMap map[string]string
	retMap     map[string][]interface{}
}

var _ (bind.PendingContractCaller) = (*mockKip103ContractCaller)(nil)

func (caller *mockKip103ContractCaller) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return []byte(kip103.TreasuryRebalanceBinRuntime), nil
}

func (caller *mockKip103ContractCaller) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	return caller.CodeAt(ctx, contract, nil)
}

func (caller *mockKip103ContractCaller) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	funcSig := hexutil.Encode(call.Data[:4])[2:]
	funcName := strings.Split(caller.funcSigMap[funcSig], "(")[0]
	mockKey := funcName

	if funcName == "retirees" || funcName == "newbies" {
		lastByte := call.Data[len(call.Data)-1]
		mockKey = fmt.Sprintf("%s%v", funcName, lastByte)
	}

	mockRet := caller.retMap[mockKey]
	return caller.abi.Methods[funcName].Outputs.Pack(mockRet...)
}

func (caller *mockKip103ContractCaller) PendingCallContract(ctx context.Context, call klaytn.CallMsg) ([]byte, error) {
	return caller.CallContract(ctx, call, nil)
}

func TestRebalanceTreasury(t *testing.T) {
	bc, istBackend := newBlockChain(1)
	defer func() {
		istBackend.Stop()
		bc.Stop()
	}()

	parsed, err := abi.JSON(strings.NewReader(kip103.TreasuryRebalanceABI))
	if err != nil {
		t.Fatal(err)
	}

	// time to generate blocks
	time.Sleep(2 * time.Second)

	block := bc.CurrentBlock()

	bc.Config().Kip103CompatibleBlock = block.Number()
	bc.Config().Kip103ContractAddress = common.Address{}

	retireds := []struct {
		addr    common.Address
		balance *big.Int
	}{
		{
			addr:    common.HexToAddress("0x9712f943b296758aaae79944ec975884188d3a96"),
			balance: new(big.Int).Mul(big.NewInt(1000000), big.NewInt(params.KLAY)),
		},
		{
			addr:    common.HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed"),
			balance: new(big.Int).Mul(big.NewInt(700000), big.NewInt(params.KLAY)),
		},
		{
			addr:    common.HexToAddress("0xfb6916095ca1df60bb79ce92ce3ea74c37c5d359"),
			balance: new(big.Int).Mul(big.NewInt(12345), big.NewInt(params.KLAY)),
		},
	}

	totalRetiredBalance := big.NewInt(0)
	for _, retired := range retireds {
		totalRetiredBalance.Add(totalRetiredBalance, retired.balance)
	}

	defaultReturnMap := make(map[string][]interface{})
	defaultReturnMap["getRetiredCount"] = []interface{}{big.NewInt(3)}
	defaultReturnMap["retirees0"] = []interface{}{retireds[0].addr}
	defaultReturnMap["retirees1"] = []interface{}{retireds[1].addr}
	defaultReturnMap["retirees2"] = []interface{}{retireds[2].addr}

	defaultReturnMap["getNewbieCount"] = []interface{}{big.NewInt(3)}
	defaultReturnMap["newbies0"] = []interface{}{common.HexToAddress("0x819104a190255e0cedbdd9d5f59a557633d79db1"), new(big.Int).Mul(big.NewInt(1000000), big.NewInt(params.KLAY))}
	defaultReturnMap["newbies1"] = []interface{}{common.HexToAddress("0x75c3098be5e4b63fbac05838daaee378dd48098d"), new(big.Int).Mul(big.NewInt(500000), big.NewInt(params.KLAY))}
	defaultReturnMap["newbies2"] = []interface{}{common.HexToAddress("0xceB7ADDFBa9665d8767173D47dE4453D7b7B900D"), new(big.Int).Mul(big.NewInt(123412), big.NewInt(params.KLAY))}

	defaultReturnMap["rebalanceBlockNumber"] = []interface{}{block.Number()}
	defaultReturnMap["status"] = []interface{}{uint8(2)}

	testCases := []struct {
		modifier func(retMap map[string][]interface{})
		// TODO-aidn: add result checker also
		expectedErr error
	}{
		{
			func(retMap map[string][]interface{}) {
				// do nothing
			},
			nil,
		},
		{
			func(retMap map[string][]interface{}) {
				retMap["status"] = []interface{}{uint8(1)}
			},
			errNotProperStatus,
		},
		{
			func(retMap map[string][]interface{}) {
				retMap["newbies0"][1] = totalRetiredBalance
				retMap["newbies1"][1] = big.NewInt(0)
				retMap["newbies2"][1] = big.NewInt(0)
			},
			nil,
		},
		{
			func(retMap map[string][]interface{}) {
				retMap["newbies0"][1] = totalRetiredBalance
				retMap["newbies1"][1] = big.NewInt(0)
				retMap["newbies2"][1] = big.NewInt(1)
			},
			errNotEnoughRetiredBal,
		},
	}

	for _, tc := range testCases {
		// reset state
		state, err := bc.StateAt(block.Root())
		if err != nil {
			t.Fatal(err)
		}

		// initializing retireds' assets
		for i := range retireds {
			state.SetBalance(retireds[i].addr, retireds[i].balance)
		}

		// modification for each test case
		mockRetMap := make(map[string][]interface{})
		for k, v := range defaultReturnMap {
			mockRetMap[k] = v
		}

		tc.modifier(mockRetMap)

		c := &mockKip103ContractCaller{abi: parsed, funcSigMap: kip103.TreasuryRebalanceFuncSigs, retMap: mockRetMap}
		ret, err := RebalanceTreasury(state, bc, block.Header(), c)
		assert.Equal(t, tc.expectedErr, err)

		// balance check
		if ret.Success {
			totalRetired := big.NewInt(0)
			for addr, amount := range ret.Retired {
				totalRetired.Add(totalRetired, amount)
				assert.Equal(t, big.NewInt(0), state.GetBalance(addr))
			}

			totalNewbie := big.NewInt(0)
			for _, amount := range ret.Newbie {
				// TODO: compare mockRetMap balance with state.GetBalance(addr)
				totalNewbie.Add(totalNewbie, amount)
			}
			assert.Equal(t, totalRetired, new(big.Int).Add(totalNewbie, ret.Burnt))
		}

		memo, _ := json.Marshal(ret)
		t.Log(string(memo))
	}
}
