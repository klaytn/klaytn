package core

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	mock_istanbul "github.com/klaytn/klaytn/consensus/istanbul/mocks"
	"github.com/klaytn/klaytn/fork"
	"github.com/klaytn/klaytn/params"
	"gotest.tools/assert/cmp"
)

func TestCore_sendPrepare(t *testing.T) {
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	validatorAddrs, validatorKeyMap := genValidators(6)
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}
	defer istCore.Stop()

	lastProposal, lastProposer := mockBackend.LastProposal()
	proposal, err := genBlock(lastProposal.(*types.Block), validatorKeyMap[validatorAddrs[0]])
	if err != nil {
		t.Fatal(err)
	}

	istCore.current.Preprepare = &istanbul.Preprepare{
		View:     istCore.currentView(),
		Proposal: proposal,
	}

	mockCtrl.Finish()

	// invalid case - not committee
	{
		// Increase round number until the owner of istanbul.core is not a member of the committee
		for istCore.valSet.CheckInSubList(lastProposal.Hash(), istCore.currentView(), istCore.Address()) {
			istCore.current.round.Add(istCore.current.round, common.Big1)
			istCore.valSet.CalcProposer(lastProposer, istCore.current.round.Uint64())
		}

		mockCtrl := gomock.NewController(t)
		mockBackend := mock_istanbul.NewMockBackend(mockCtrl)
		mockBackend.EXPECT().Sign(gomock.Any()).Return(nil, nil).Times(0)
		mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(0)

		istCore.backend = mockBackend
		istCore.sendPrepare()

		// methods of mockBackend should be executed given times
		mockCtrl.Finish()
	}

	// valid case
	{
		// Increase round number until the owner of istanbul.core become a member of the committee
		for !istCore.valSet.CheckInSubList(lastProposal.Hash(), istCore.currentView(), istCore.Address()) {
			istCore.current.round.Add(istCore.current.round, common.Big1)
			istCore.valSet.CalcProposer(lastProposer, istCore.current.round.Uint64())
		}

		mockCtrl := gomock.NewController(t)
		mockBackend := mock_istanbul.NewMockBackend(mockCtrl)
		mockBackend.EXPECT().Sign(gomock.Any()).Return(nil, nil).Times(1)
		mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		istCore.backend = mockBackend
		istCore.sendPrepare()

		// methods of mockBackend should be executed given times
		mockCtrl.Finish()
	}
}

func BenchmarkMsgCmp(b *testing.B) {
	getEmptySubject := func() istanbul.Subject {
		return istanbul.Subject{
			View: &istanbul.View{
				Round:    big.NewInt(0),
				Sequence: big.NewInt(0),
			},
			Digest:   common.HexToHash("1"),
			PrevHash: common.HexToHash("2"),
		}
	}
	s1, s2 := getEmptySubject(), getEmptySubject()

	// Worst
	b.Run("reflect.DeepEqual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflect.DeepEqual(s1, s2)
		}
	})

	// Better
	b.Run("cmp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cmp.Equal(&s1, &s2)
		}
	})

	// Best
	b.Run("own", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s1.Equal(&s2)
		}
	})
}
