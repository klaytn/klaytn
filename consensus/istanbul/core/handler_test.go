package core

import (
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	mock_istanbul "github.com/klaytn/klaytn/consensus/istanbul/mocks"
	"github.com/klaytn/klaytn/consensus/istanbul/validator"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

// newMockBackend create a mock-backend initialized with default values
func newMockBackend(t *testing.T, validatorAddrs []common.Address) (*mock_istanbul.MockBackend, *gomock.Controller) {
	committeeSize := uint64(len(validatorAddrs) / 3)

	istExtra := &types.IstanbulExtra{
		Validators:    validatorAddrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}
	extra, err := rlp.EncodeToBytes(istExtra)
	if err != nil {
		t.Fatal(err)
	}

	initBlock := types.NewBlockWithHeader(&types.Header{
		ParentHash: common.Hash{},
		Number:     common.Big0,
		GasUsed:    0,
		Extra:      append(make([]byte, types.IstanbulExtraVanity), extra...),
		Time:       new(big.Int).SetUint64(1234),
		BlockScore: common.Big0,
	})

	eventMux := new(event.TypeMux)
	validatorSet := validator.NewWeightedCouncil(validatorAddrs, validatorAddrs, nil, nil,
		istanbul.WeightedRandom, committeeSize, 0, 0, &blockchain.BlockChain{})

	mockCtrl := gomock.NewController(t)
	mockBackend := mock_istanbul.NewMockBackend(mockCtrl)

	// Consider the last proposal is "initBlock" and the owner of mockBackend is validatorAddrs[0]
	mockBackend.EXPECT().Address().Return(validatorAddrs[0]).AnyTimes()
	mockBackend.EXPECT().LastProposal().Return(initBlock, validatorAddrs[0]).AnyTimes()
	mockBackend.EXPECT().Validators(initBlock).Return(validatorSet).AnyTimes()
	mockBackend.EXPECT().NodeType().Return(common.CONSENSUSNODE).AnyTimes()

	// Set an eventMux in which istanbul core will subscribe istanbul events
	mockBackend.EXPECT().EventMux().Return(eventMux).AnyTimes()

	// Just for bypassing an unused function
	mockBackend.EXPECT().SetCurrentView(gomock.Any()).Return().AnyTimes()

	// Always return nil for broadcasting related functions
	mockBackend.EXPECT().Sign(gomock.Any()).Return(nil, nil).AnyTimes()
	mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockBackend.EXPECT().GossipSubPeer(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	// Verify checks whether the proposal of the preprepare message is a valid block. Consider it valid.
	mockBackend.EXPECT().Verify(gomock.Any()).Return(time.Duration(0), nil).AnyTimes()

	return mockBackend, mockCtrl
}

// genValidators returns a set of addresses and corresponding keys used for generating a validator set
func genValidators(n int) ([]common.Address, map[common.Address]*ecdsa.PrivateKey) {
	addrs := make([]common.Address, n)
	keyMap := make(map[common.Address]*ecdsa.PrivateKey, n)

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(key.PublicKey)
		keyMap[addrs[i]] = key
	}
	return addrs, keyMap
}

// getRandomValidator selects a validator in the given validator set.
// `isCommittee` determines whether it returns a committee or a non-committee.
func getRandomValidator(isCommittee bool, valSet istanbul.ValidatorSet, prevHash common.Hash, view *istanbul.View) istanbul.Validator {
	committee := valSet.SubList(prevHash, view)

	if isCommittee {
		return committee[rand.Int()%(len(committee)-1)]
	}

	for _, val := range valSet.List() {
		for _, com := range committee {
			if val.Address() == com.Address() {
				isCommittee = true
			}
		}
		if !isCommittee {
			return val
		}
		isCommittee = false
	}

	// it should not be happened
	return nil
}

// signBlock signs the given block with the given private key
func signBlock(block *types.Block, privateKey *ecdsa.PrivateKey) (*types.Block, error) {
	var hash common.Hash
	header := block.Header()
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal
	rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])

	seal, err := crypto.Sign(crypto.Keccak256([]byte(hash.Bytes())), privateKey)
	if err != nil {
		return nil, err
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
	}
	istanbulExtra.Seal = seal

	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return nil, err
	}

	header.Extra = append(header.Extra[:types.IstanbulExtraVanity], payload...)
	return block.WithSeal(header), nil
}

// genBlock generates a signed block indicating prevBlock with ParentHash
func genBlock(prevBlock *types.Block, signerKey *ecdsa.PrivateKey) (*types.Block, error) {
	block := types.NewBlockWithHeader(&types.Header{
		ParentHash: prevBlock.Hash(),
		Number:     new(big.Int).Add(prevBlock.Number(), common.Big1),
		GasUsed:    0,
		Extra:      prevBlock.Extra(),
		Time:       new(big.Int).Add(prevBlock.Time(), common.Big1),
		BlockScore: new(big.Int).Add(prevBlock.BlockScore(), common.Big1),
	})
	return signBlock(block, signerKey)
}

// genIstanbulMsg generates an istanbul message with given values
func genIstanbulMsg(msgType uint64, prevHash common.Hash, proposal *types.Block, signerAddr common.Address, signerKey *ecdsa.PrivateKey) (istanbul.MessageEvent, error) {
	var subject interface{}

	if msgType == msgPreprepare {
		subject = &istanbul.Preprepare{
			View: &istanbul.View{
				Round:    big.NewInt(0),
				Sequence: proposal.Number(),
			},
			Proposal: proposal,
		}
	} else {
		subject = &istanbul.Subject{
			View: &istanbul.View{
				Round:    big.NewInt(0),
				Sequence: proposal.Number(),
			},
			Digest:   proposal.Hash(),
			PrevHash: prevHash,
		}
	}

	encodedSubject, err := Encode(subject)
	if err != nil {
		return istanbul.MessageEvent{}, err
	}

	msg := &message{
		Hash:    prevHash,
		Code:    msgType,
		Msg:     encodedSubject,
		Address: signerAddr,
	}

	data, err := msg.PayloadNoSig()
	if err != nil {
		return istanbul.MessageEvent{}, err
	}

	msg.Signature, err = crypto.Sign(crypto.Keccak256([]byte(data)), signerKey)
	if err != nil {
		return istanbul.MessageEvent{}, err
	}

	encodedPayload, err := msg.Payload()
	if err != nil {
		return istanbul.MessageEvent{}, err
	}

	istMsg := istanbul.MessageEvent{
		Hash:    msg.Hash,
		Payload: encodedPayload,
	}

	return istMsg, nil
}

// TestCore_handleEvents_scenario_invalidSender tests `handleEvents` function of `istanbul.core` with a scenario.
// It posts an invalid message and a valid message of each istanbul message type.
func TestCore_handleEvents_scenario_invalidSender(t *testing.T) {
	validatorAddrs, validatorKeyMap := genValidators(30)
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	// When the istanbul core started, a message handling loop in `handleEvents()` waits istanbul messages
	istCore := New(mockBackend, istConfig).(*core)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}

	// Get variables initialized on `newMockBackend()`
	eventMux := mockBackend.EventMux()
	lastProposal, _ := mockBackend.LastProposal()
	lastBlock := lastProposal.(*types.Block)
	validators := mockBackend.Validators(lastBlock)

	// Preprepare message originated from invalid sender
	{
		msgSender := getRandomValidator(false, validators, lastBlock.Hash(), istCore.currentView())
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		newProposal, err := genBlock(lastBlock, msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		istanbulMsg, err := genIstanbulMsg(msgPreprepare, lastProposal.Hash(), newProposal, msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		if err := eventMux.Post(istanbulMsg); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		assert.Nil(t, istCore.current.Preprepare)
	}

	// Preprepare message originated from valid sender and set a new proposal in the istanbul core
	{
		msgSender := validators.GetProposer()
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		newProposal, err := genBlock(lastBlock, msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		istanbulMsg, err := genIstanbulMsg(msgPreprepare, lastBlock.Hash(), newProposal, msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		if err := eventMux.Post(istanbulMsg); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		assert.Equal(t, istCore.current.Preprepare.Proposal.Header().String(), newProposal.Header().String())
	}

	// Prepare message originated from invalid sender
	{
		msgSender := getRandomValidator(false, validators, lastBlock.Hash(), istCore.currentView())
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		istanbulMsg, err := genIstanbulMsg(msgPrepare, lastBlock.Hash(), istCore.current.Preprepare.Proposal.(*types.Block), msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		if err := eventMux.Post(istanbulMsg); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		assert.Equal(t, 0, len(istCore.current.Prepares.messages))
	}

	// Prepare message originated from valid sender
	{
		msgSender := getRandomValidator(true, validators, lastBlock.Hash(), istCore.currentView())
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		istanbulMsg, err := genIstanbulMsg(msgPrepare, lastBlock.Hash(), istCore.current.Preprepare.Proposal.(*types.Block), msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		if err := eventMux.Post(istanbulMsg); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		assert.Equal(t, 1, len(istCore.current.Prepares.messages))
	}

	// Commit message originated from invalid sender
	{
		msgSender := getRandomValidator(false, validators, lastBlock.Hash(), istCore.currentView())
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		istanbulMsg, err := genIstanbulMsg(msgCommit, lastBlock.Hash(), istCore.current.Preprepare.Proposal.(*types.Block), msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		if err := eventMux.Post(istanbulMsg); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		assert.Equal(t, 0, len(istCore.current.Commits.messages))
	}

	// Commit message originated from valid sender
	{
		msgSender := getRandomValidator(true, validators, lastBlock.Hash(), istCore.currentView())
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		istanbulMsg, err := genIstanbulMsg(msgCommit, lastBlock.Hash(), istCore.current.Preprepare.Proposal.(*types.Block), msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		if err := eventMux.Post(istanbulMsg); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		assert.Equal(t, 1, len(istCore.current.Commits.messages))
	}

	//// RoundChange message originated from invalid sender
	//{
	//	msgSender := getRandomValidator(false, validators, lastBlock.Hash(), istCore.currentView())
	//	msgSenderKey := validatorKeyMap[msgSender.Address()]
	//
	//	istanbulMsg, err := genIstanbulMsg(msgRoundChange, lastBlock.Hash(), istCore.current.Preprepare.Proposal.(*types.Block), msgSender.Address(), msgSenderKey)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	if err := eventMux.Post(istanbulMsg); err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	time.Sleep(time.Second)
	//	assert.Nil(t, istCore.roundChangeSet.roundChanges[0]) // round is set to 0 in this test
	//}

	// RoundChange message originated from valid sender
	{
		msgSender := getRandomValidator(true, validators, lastBlock.Hash(), istCore.currentView())
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		istanbulMsg, err := genIstanbulMsg(msgRoundChange, lastBlock.Hash(), istCore.current.Preprepare.Proposal.(*types.Block), msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		if err := eventMux.Post(istanbulMsg); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		assert.Equal(t, 1, len(istCore.roundChangeSet.roundChanges[0].messages)) // round is set to 0 in this test
	}
}

func TestCore_handlerMsg(t *testing.T) {
	validatorAddrs, validatorKeyMap := genValidators(10)
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}

	lastProposal, _ := mockBackend.LastProposal()
	lastBlock := lastProposal.(*types.Block)
	validators := mockBackend.Validators(lastBlock)

	// invalid format
	{
		invalidMsg := []byte{0x1, 0x2, 0x3, 0x4}
		err := istCore.handleMsg(invalidMsg)
		assert.NotNil(t, err)
	}

	// invali sender (non-validator)
	{
		newAddr, keyMap := genValidators(1)
		nonValidatorAddr := newAddr[0]
		nonValidatorKey := keyMap[nonValidatorAddr]

		newProposal, err := genBlock(lastBlock, nonValidatorKey)
		if err != nil {
			t.Fatal(err)
		}

		istanbulMsg, err := genIstanbulMsg(msgPreprepare, lastBlock.Hash(), newProposal, nonValidatorAddr, nonValidatorKey)
		if err != nil {
			t.Fatal(err)
		}

		err = istCore.handleMsg(istanbulMsg.Payload)
		assert.NotNil(t, err)
	}

	// valid message
	{
		msgSender := validators.GetProposer()
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		newProposal, err := genBlock(lastBlock, msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		istanbulMsg, err := genIstanbulMsg(msgPreprepare, lastBlock.Hash(), newProposal, msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		err = istCore.handleMsg(istanbulMsg.Payload)
		assert.Nil(t, err)
	}
}
