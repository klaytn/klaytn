package core

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	mock_istanbul "github.com/klaytn/klaytn/consensus/istanbul/mocks"
	"github.com/klaytn/klaytn/fork"
	"github.com/klaytn/klaytn/params"
)

func TestCore_sendCommit(t *testing.T) {
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
		istCore.sendCommit()

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
		mockBackend.EXPECT().Sign(gomock.Any()).Return(nil, nil).Times(2)
		mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		istCore.backend = mockBackend
		istCore.sendCommit()

		// methods of mockBackend should be executed given times
		mockCtrl.Finish()
	}
}
