package core

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	mock_istanbul "github.com/klaytn/klaytn/consensus/istanbul/mocks"
	"github.com/klaytn/klaytn/fork"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
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
	b.Run("EqualImpl", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s1.Equal(&s2)
		}
	})
}

func TestSubjectCmp(t *testing.T) {
	genRandomHash := func(n int) common.Hash {
		b := make([]byte, n)
		_, err := rand.Read(b)
		assert.Nil(t, err)
		return common.BytesToHash(b)
	}
	genRandomInt := func(min, max int) int64 {
		return int64(rand.Intn(max-min) + min)
	}
	genSubject := func(min, max int) *istanbul.Subject {
		round, seq := big.NewInt(genRandomInt(min, max)), big.NewInt(genRandomInt(min, max))
		digest, prevHash := genRandomHash(max), genRandomHash(max)
		return &istanbul.Subject{
			View: &istanbul.View{
				Round:    round,
				Sequence: seq,
			},
			Digest:   digest,
			PrevHash: prevHash,
		}
	}
	copySubject := func(s *istanbul.Subject) *istanbul.Subject {
		r := new(istanbul.Subject)
		v := new(istanbul.View)
		r.Digest = s.Digest
		r.PrevHash = s.PrevHash
		v.Round = new(big.Int).SetUint64(s.View.Round.Uint64())
		v.Sequence = new(big.Int).SetUint64(s.View.Sequence.Uint64())
		r.View = v
		return r
	}

	rand.Seed(time.Now().UnixNano())
	min, max, n := 1, 9999, 10000
	var identity bool
	var s1, s2 *istanbul.Subject
	for i := 0; i < n; i++ {
		s1 = genSubject(min, max)
		if rand.Intn(2) == 0 {
			identity = true
			s2 = copySubject(s1)
		} else {
			identity = false
			s2 = genSubject(max+1, max*2)
		}
		e := s1.Equal(s2)
		if identity {
			assert.Equal(t, e, true)
		} else {
			assert.Equal(t, e, false)
		}
		assert.Equal(t, e, reflect.DeepEqual(s1, s2))
	}
}

func TestNilSubjectCmp(t *testing.T) {
	sbj := istanbul.Subject{
		View: &istanbul.View{
			Round:    big.NewInt(0),
			Sequence: big.NewInt(0),
		},
		Digest:   common.HexToHash("1"),
		PrevHash: common.HexToHash("2"),
	}
	var nilSbj *istanbul.Subject = nil

	assert.Equal(t, sbj.Equal(nil), false)
	assert.Equal(t, sbj.Equal(nilSbj), false)
	assert.Equal(t, nilSbj.Equal(&sbj), false)
	assert.Equal(t, nilSbj.Equal(nilSbj), true)
	assert.Equal(t, nilSbj.Equal(nil), true)
}
