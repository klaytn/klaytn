package core

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/log/term"
	"github.com/mattn/go-colorable"
	"io"
	"math/big"
	"math/rand"
	"os"
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
	"github.com/klaytn/klaytn/fork"
	"github.com/klaytn/klaytn/params"
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
	validatorSet := validator.NewWeightedCouncil(validatorAddrs, nil, validatorAddrs, nil, nil,
		istanbul.WeightedRandom, committeeSize, 0, 0, &blockchain.BlockChain{})

	mockCtrl := gomock.NewController(t)
	mockBackend := mock_istanbul.NewMockBackend(mockCtrl)

	// Consider the last proposal is "initBlock" and the owner of mockBackend is validatorAddrs[0]
	mockBackend.EXPECT().Address().Return(validatorAddrs[0]).AnyTimes()
	mockBackend.EXPECT().LastProposal().Return(initBlock, validatorAddrs[0]).AnyTimes()
	mockBackend.EXPECT().Validators(initBlock).Return(validatorSet).AnyTimes()
	mockBackend.EXPECT().NodeType().Return(common.CONSENSUSNODE).AnyTimes()
	mockBackend.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

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

	mockBackend.EXPECT().HasBadProposal(gomock.Any()).Return(true).AnyTimes()

	// Commit is added to remove unexpected call error

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

// genMaliciousBlock generates a modified block indicating prevBlock with ParentHash
func genMaliciousBlock(prevBlock *types.Block, signerKey *ecdsa.PrivateKey) (*types.Block, error) {
	block := types.NewBlockWithHeader(&types.Header{
		ParentHash: prevBlock.Hash(),
		//ParentHash: common.HexToHash("2"),
		Number:     common.Big0,
		GasUsed:    0,
		Extra:      prevBlock.Extra(),
		Time:       common.Big0,
		BlockScore: common.Big0,
	})
	return signBlock(block, signerKey)
}

// genBlock generates a signed block indicating prevBlock with ParentHash
func genBlockA(prevBlock *types.Block, signerKey *ecdsa.PrivateKey) (*types.Block, error) {
	block := types.NewBlockWithHeader(&types.Header{
		ParentHash: prevBlock.Hash(),
		Number:     new(big.Int).Add(prevBlock.Number(), common.Big1),
		GasUsed:    0,
		Extra:      prevBlock.Extra(),
		//Extra:      []byte{'a', 'b', 'c'},
		Time:       new(big.Int).Add(prevBlock.Time(), common.Big1),
		BlockScore: new(big.Int).Add(prevBlock.BlockScore(), common.Big1),
	})
	return signBlock(block, signerKey)
}

// genBlock generates a signed block indicating prevBlock with ParentHash
func genBlockB(prevBlock *types.Block, signerKey *ecdsa.PrivateKey) (*types.Block, error) {
	block := types.NewBlockWithHeader(&types.Header{
		ParentHash: prevBlock.Hash(),
		Number:     new(big.Int).Add(prevBlock.Number(), common.Big1),
		GasUsed:    1000,
		Extra:      prevBlock.Extra(),
		//Extra: []byte{'x', 'y', 'z'},

		Time:       new(big.Int).Add(prevBlock.Time(), big.NewInt(10)),
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
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

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
	defer istCore.Stop()

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
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	validatorAddrs, validatorKeyMap := genValidators(10)
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}
	defer istCore.Stop()

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

// TODO-Klaytn: To enable logging in the test code, we can use the following function.
// This function will be moved to somewhere utility functions are located.
func enableLog() {
	usecolor := term.IsTty(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	glogger := log.NewGlogHandler(log.StreamHandler(output, log.TerminalFormat(usecolor)))
	log.PrintOrigins(true)
	log.ChangeGlobalLogLevel(glogger, log.Lvl(3))
	glogger.Vmodule("")
	glogger.BacktraceAt("")
	log.Root().SetHandler(glogger)
}

func TestCore_MaliciousCN(t *testing.T) {

	enableLog()
	var maxAllowedRound = uint64(5)

	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	validatorAddrs, validatorKeyMap := genValidators(12)
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}
	defer istCore.Stop()

	lastProposal, _ := mockBackend.LastProposal()
	lastBlock := lastProposal.(*types.Block)
	validators := mockBackend.Validators(lastBlock)
	//fmt.Println(validators.SubList(lastBlock.Hash(), istCore.currentView()))

	// valid message
	{
		msgSender := validators.GetProposer()
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		newProposal, err := genBlock(lastBlock, msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}
		malProposal, err := genMaliciousBlock(lastBlock, msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		istanbulMsg, _ := genIstanbulMsg(msgPreprepare, lastBlock.Hash(), newProposal, msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}
		//fmt.Println("msgsender = ", msgSender.Address().Hex())
		//fmt.Println("validators = ", validators.SubList(lastBlock.Hash(), istCore.currentView()))
		istCore.handleMsg(istanbulMsg.Payload)
		//fmt.Println("current state = ", istCore.state)
		//validators := sb.getValidators(sb.chain.CurrentHeader().Number.Uint64(), sb.chain.CurrentHeader().Hash())
		//for _, val := range validators.List() {
		//	if addr == val.Address() {
		//		return nil
		//	}
		sendMessages := func(state uint64, malicious int) {
			cnt := 0
			//for k, v := range validatorKeyMap {
			for _, k := range validators.SubList(lastBlock.Hash(), istCore.currentView()) {
				v := validatorKeyMap[k.Address()]

				var msg *types.Block
				// the proposer does not send prepare message
				//if msgSender.Address().Hex() == k.Address().Hex() && state == msgPrepare {
				//	continue
				//}
				if cnt < malicious {
					msg = malProposal
					cnt++
				} else {
					msg = newProposal

				}
				istanbulMsg, _ := genIstanbulMsg(state, lastBlock.Hash(), msg, k.Address(), v)
				err = istCore.handleMsg(istanbulMsg.Payload)
				if err != nil {
					fmt.Println(err)
				}

			}
		}

		sendMessages(msgPrepare, 1)
		preparesSize := istCore.current.Prepares.Size()
		commitsSize := istCore.current.Commits.Size()

		fmt.Println("prepare size = ", preparesSize, " commit size: ", commitsSize)

		if istCore.state.Cmp(StatePrepared) < 0 {
			for {
				if istCore.currentView().Round.Uint64() > maxAllowedRound {
					t.Fatal("not prepared due to malicious cns")

				}
				//istCore.sendRoundChange(istCore.currentView().Round)
				istCore.sendEvent(timeoutEvent{&istanbul.View{
					Sequence: istCore.current.sequence,
					Round:    new(big.Int).Add(istCore.current.round, common.Big1),
				}})
				//istCore.sendRoundChange(new(big.Int).Add(istCore.currentView().Round, common.Big1))
				//fmt.Println("current round: ", istCore.currentView().Round)

				time.Sleep(1000 * time.Millisecond)
			}

		}

		//fmt.Println("prepared. state = ", istCore.state)

		sendMessages(msgCommit, 2)
		commitsSize = istCore.current.Commits.Size()

		fmt.Println("prepare size = ", preparesSize, " commit size: ", commitsSize)

		if istCore.state.Cmp(StateCommitted) < 0 {
			for {
				if istCore.currentView().Round.Uint64() > maxAllowedRound {
					t.Fatal("not committed due to malicious cns")

				}
				//istCore.sendRoundChange(istCore.currentView().Round)
				istCore.sendEvent(timeoutEvent{&istanbul.View{
					Sequence: istCore.current.sequence,
					Round:    new(big.Int).Add(istCore.current.round, common.Big1),
				}})
				//istCore.sendRoundChange(new(big.Int).Add(istCore.currentView().Round, common.Big1))
				//fmt.Println("current round: ", istCore.currentView().Round)

				time.Sleep(1000 * time.Millisecond)
			}

		}

		assert.Nil(t, err)

	}
}

func TestCore_MaliciousCN_5nodes(t *testing.T) {

	enableLog()

	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	validatorAddrs, validatorKeyMap := genValidators(15)
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	istCoreA := New(mockBackend, istConfig).(*core)
	istCoreB := New(mockBackend, istConfig).(*core)

	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}
	if err := istCoreA.Start(); err != nil {
		t.Fatal(err)
	}
	if err := istCoreB.Start(); err != nil {
		t.Fatal(err)
	}
	defer istCore.Stop()
	defer istCoreA.Stop()
	defer istCoreB.Stop()

	lastProposal, _ := mockBackend.LastProposal()
	lastBlock := lastProposal.(*types.Block)
	validators := mockBackend.Validators(lastBlock)
	//fmt.Println("validator addrs = ", validators.SubList(lastBlock.Hash(), istCore.currentView()))

	// validatorA, validatorB
	// msgSender sends a block of sequence #0 to validatorA
	// msgSender sends a block of sequence #1 to validatorB
	// node 5 f = 1 2f+1 = 3

	// valid message
	{
		msgSender := validators.GetProposer()
		msgSenderKey := validatorKeyMap[msgSender.Address()]

		newProposalA, err := genBlockA(lastBlock, msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}
		newProposalB, err := genBlockB(lastBlock, msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		istanbulMsgA, _ := genIstanbulMsg(msgPreprepare, lastBlock.Hash(), newProposalA, msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}
		istanbulMsgB, _ := genIstanbulMsg(msgPreprepare, lastBlock.Hash(), newProposalB, msgSender.Address(), msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		istCore.handleMsg(istanbulMsgA.Payload)
		istCoreA.handleMsg(istanbulMsgA.Payload)
		istCoreA.handleMsg(istanbulMsgA.Payload)
		istCore.handleMsg(istanbulMsgB.Payload)
		istCoreB.handleMsg(istanbulMsgB.Payload)
		istCoreB.handleMsg(istanbulMsgB.Payload)

		//{0,1,2,3,4} <- proposer이 항상 0 이라는 보장이 없어서 프로포저가 0이면 {0,1,2}{0,3,4} 프로포저가 1이면 {0,1,2}{1,3,4}

		committeeSize := len(validators.SubList(lastBlock.Hash(), istCore.currentView()))
		tmpList := validators.SubList(lastBlock.Hash(), istCore.currentView())
		//fmt.Println(tmpList)

		for i := range validators.SubList(lastBlock.Hash(), istCore.currentView()) {
			// both groupA and groupB includes proposer as validator
			if tmpList[i].Address() == msgSender.Address() {
				tmpList = append(tmpList[i+1:], tmpList[:i]...)
				break
			}
		}

		listA := make([]istanbul.Validator, (committeeSize-1)/2)
		listB := make([]istanbul.Validator, (committeeSize-1)/2)
		copy(listA, tmpList[:(committeeSize-1)/2])
		listA = append(listA, msgSender)

		copy(listB, tmpList[(committeeSize-1)/2:])
		listB = append(listB, msgSender)

		//fmt.Println(listA)
		//fmt.Println(listB)

		for _, k := range listA {
			v := validatorKeyMap[k.Address()]

			istanbulMsg, _ := genIstanbulMsg(msgPrepare, lastBlock.Hash(), newProposalA, k.Address(), v)
			err = istCoreA.handleMsg(istanbulMsg.Payload)
			if err != nil {
				fmt.Println(err)
			}
		}

		//sendMessages := func(state uint64, malicious int) {
		//	//for k, v := range validatorKeyMap {
		//	for _, k := range listB {
		//		v := validatorKeyMap[k.Address()]
		//
		//		var msg *types.Block
		//		// the proposer does not send prepare message
		//		//if msgSender.Address().Hex() == k.Address().Hex() && state == msgPrepare {
		//		//	continue
		//		//}
		//
		//		istanbulMsg, _ := genIstanbulMsg(state, lastBlock.Hash(), msg, k.Address(), v)
		//		err = istCoreA.handleMsg(istanbulMsg.Payload)
		//		if err != nil {
		//			fmt.Println(err)
		//		}
		//
		//	}
		//}

		//sendMessages(msgPrepare, 1)
		preparesSizeA := istCoreA.current.Prepares.Size()

		for _, k := range listA {
			v := validatorKeyMap[k.Address()]

			istanbulMsg, _ := genIstanbulMsg(msgCommit, lastBlock.Hash(), newProposalA, k.Address(), v)
			err = istCoreA.handleMsg(istanbulMsg.Payload)
			if err != nil {
				fmt.Println(err)
			}
		}
		commitsSizeA := istCoreA.current.Commits.Size()

		fmt.Println("group 1 prepare size = ", preparesSizeA, "group 1 commit size = ", commitsSizeA)

		for _, k := range listB {
			v := validatorKeyMap[k.Address()]

			istanbulMsg, _ := genIstanbulMsg(msgPrepare, lastBlock.Hash(), newProposalB, k.Address(), v)
			err = istCoreB.handleMsg(istanbulMsg.Payload)
			if err != nil {
				fmt.Println(err)
			}
		}

		for _, k := range listB {
			v := validatorKeyMap[k.Address()]

			istanbulMsg, _ := genIstanbulMsg(msgCommit, lastBlock.Hash(), newProposalB, k.Address(), v)
			err = istCoreB.handleMsg(istanbulMsg.Payload)
			if err != nil {
				fmt.Println(err)
			}
		}
		preparesSizeB := istCoreB.current.Prepares.Size()
		commitsSizeB := istCoreB.current.Commits.Size()
		fmt.Println("group 2 prepare size = ", preparesSizeB, "group 2 commit size = ", commitsSizeB)

		assert.Nil(t, err)

	}
}

// TestCore_handleTimeoutMsg_race tests a race condition between round change triggers.
// There should be no race condition when round change message and timeout event are handled simultaneously.

func TestCore_handleTimeoutMsg_race(t *testing.T) {
	// important variables to construct test cases
	const sleepTime = 200 * time.Millisecond
	const processingTime = 400 * time.Millisecond

	type testCase struct {
		name          string
		timeoutTime   time.Duration
		messageRound  int64
		expectedRound int64
	}
	testCases := []testCase{
		{
			// if timeoutTime < sleepTime,
			// timeout event will be posted and then round change message will be processed
			name:          "timeout before processing the (2f+1)th round change message",
			timeoutTime:   50 * time.Millisecond,
			messageRound:  10,
			expectedRound: 10,
		},
		{
			// if timeoutTime > sleepTime && timeoutTime < (processingTime + sleepTime),
			// timeout event will be posted during the processing of (2f+1)th round change message
			name:          "timeout during processing the (2f+1)th round change message",
			timeoutTime:   300 * time.Millisecond,
			messageRound:  20,
			expectedRound: 20,
		},
	}

	validatorAddrs, _ := genValidators(10)
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}
	defer istCore.Stop()

	eventMux := mockBackend.EventMux()
	lastProposal, _ := mockBackend.LastProposal()
	sequence := istCore.current.sequence.Int64()

	for _, tc := range testCases {
		handler := func(t *testing.T) {
			roundChangeTimer := istCore.roundChangeTimer.Load().(*time.Timer)

			// reset timeout timer of this round and wait some time
			roundChangeTimer.Reset(tc.timeoutTime)
			time.Sleep(sleepTime)

			// `istCore.validateFn` will be executed on processing a istanbul message
			istCore.validateFn = func(arg1 []byte, arg2 []byte) (common.Address, error) {
				// postpones the processing of a istanbul message
				time.Sleep(processingTime)
				return common.Address{}, nil
			}

			// prepare a round change message payload
			payload := makeRCMsgPayload(tc.messageRound, sequence, lastProposal.Hash(), validatorAddrs[0])
			if payload == nil {
				t.Fatal("failed to make a round change message payload")
			}

			// one round change message changes the round because the committee size of mockBackend is 3
			err := eventMux.Post(istanbul.MessageEvent{
				Hash:    lastProposal.Hash(),
				Payload: payload,
			})
			if err != nil {
				t.Fatal(err)
			}

			// wait until the istanbul message have processed
			time.Sleep(processingTime + sleepTime)
			roundChangeTimer.Stop()

			// check the result
			assert.Equal(t, tc.expectedRound, istCore.current.round.Int64())
		}
		t.Run(tc.name, handler)
	}
}

// makeRCMsgPayload makes a payload of round change message.
func makeRCMsgPayload(round int64, sequence int64, prevHash common.Hash, senderAddr common.Address) []byte {
	subject, err := Encode(&istanbul.Subject{
		View: &istanbul.View{
			Round:    big.NewInt(round),
			Sequence: big.NewInt(sequence),
		},
		Digest:   common.Hash{},
		PrevHash: prevHash,
	})
	if err != nil {
		return nil
	}

	msg := &message{
		Hash:    prevHash,
		Code:    msgRoundChange,
		Msg:     subject,
		Address: senderAddr,
	}

	payload, err := msg.Payload()
	if err != nil {
		return nil
	}

	return payload
}
