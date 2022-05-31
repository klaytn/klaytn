package core

import (
	"crypto/ecdsa"
	"fmt"
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
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/log/term"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// genBlockParams generates a signed block indicating prevBlock with ParentHash with additional parameters.
func genBlockParams(prevBlock *types.Block, signerKey *ecdsa.PrivateKey, gasUsed uint64, time int64, blockScore int64) (*types.Block, error) {
	block := types.NewBlockWithHeader(&types.Header{
		ParentHash: prevBlock.Hash(),
		Number:     new(big.Int).Add(prevBlock.Number(), common.Big1),
		GasUsed:    gasUsed,
		Extra:      prevBlock.Extra(),
		Time:       new(big.Int).Add(prevBlock.Time(), big.NewInt(time)),
		BlockScore: new(big.Int).Add(prevBlock.BlockScore(), big.NewInt(blockScore)),
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

// splitSubList splits a committee into two groups w/o proposer
// one for n nodes, the other for len(committee) - n - 1 nodes
func splitSubList(committee []istanbul.Validator, n int, proposerAddr common.Address) ([]istanbul.Validator, []istanbul.Validator) {
	var subCN, remainingCN []istanbul.Validator

	for _, val := range committee {
		if val.Address() == proposerAddr {
			// proposer is not included in any group
			continue
		}
		if len(subCN) < n {
			subCN = append(subCN, val)
		} else {
			remainingCN = append(remainingCN, val)
		}
	}
	return subCN, remainingCN
}

// Simulate a proposer that receives messages from disagreeing groups of CNs.
func simulateMaliciousCN(t *testing.T, numValidators int, numMalicious int) State {
	if testing.Verbose() {
		enableLog()
	}

	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	// Note that genValidators(n) will generate n/3 validators.
	// We want n validators, thus calling genValidators(3n).
	validatorAddrs, validatorKeyMap := genValidators(numValidators * 3)

	// Add more EXPECT()s to remove unexpected call error
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	mockBackend.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockBackend.EXPECT().HasBadProposal(gomock.Any()).Return(true).AnyTimes()
	defer mockCtrl.Finish()

	var (
		// it creates two pre-defined blocks: one for benign CNs, the other for the malicious
		// newProposal is a block which the proposer has created
		// malProposal is an incorrect block that malicious CNs use to try stop consensus
		lastProposal, _ = mockBackend.LastProposal()
		lastBlock       = lastProposal.(*types.Block)
		validators      = mockBackend.Validators(lastBlock)
		proposer        = validators.GetProposer()
		proposerKey     = validatorKeyMap[proposer.Address()]
		// the proposer generates a block as newProposal
		// malicious CNs does not accept the proposer's block and use malProposal's hash value for consensus
		newProposal, _ = genBlockParams(lastBlock, proposerKey, 0, 1, 1)
		malProposal, _ = genBlockParams(lastBlock, proposerKey, 0, 0, 0)
	)

	// Start istanbul core
	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom
	istCore := New(mockBackend, istConfig).(*core)
	err := istCore.Start()
	require.Nil(t, err)
	defer istCore.Stop()

	// Step 1 - Pre-prepare with correct message

	// Create pre-prepare message
	istanbulMsg, err := genIstanbulMsg(msgPreprepare, lastBlock.Hash(), newProposal, proposer.Address(), proposerKey)
	require.Nil(t, err)

	// Handle pre-prepare message
	err = istCore.handleMsg(istanbulMsg.Payload)
	require.Nil(t, err)

	// splitSubList split current committee into benign CNs and malicious CNs
	subList := validators.SubList(lastBlock.Hash(), istCore.currentView())
	maliciousCNs, benignCNs := splitSubList(subList, numMalicious, proposer.Address())
	benignCNs = append(benignCNs, proposer)

	// Shortcut for sending consensus message to everyone in `CNList`
	sendMessages := func(state uint64, proposal *types.Block, CNList []istanbul.Validator) {
		for _, val := range CNList {
			valKey := validatorKeyMap[val.Address()]
			istanbulMsg, err := genIstanbulMsg(state, lastBlock.Hash(), proposal, val.Address(), valKey)
			assert.Nil(t, err)
			err = istCore.handleMsg(istanbulMsg.Payload)
			// assert.Nil(t, err)
		}
	}

	// Step 2 - Receive disagreeing prepare messages

	sendMessages(msgPrepare, newProposal, benignCNs)
	sendMessages(msgPrepare, malProposal, maliciousCNs)

	if istCore.state.Cmp(StatePreprepared) == 0 {
		t.Logf("State stuck at preprepared")
		return istCore.state
	}

	// Step 3 - Receive disagreeing commit messages

	sendMessages(msgCommit, newProposal, benignCNs)
	sendMessages(msgCommit, malProposal, maliciousCNs)
	return istCore.state
}

// TestCore_MalCN tests whether the proposer can commit when malicious CNs exist.
func TestCore_malCN(t *testing.T) {
	// If there are less than 'f' malicious CNs, proposer can commit.
	state := simulateMaliciousCN(t, 4, 1)
	assert.Equal(t, StateCommitted, state)

	// If there are more than 'f' malicious CNs, the proposer cannot commit, stuck at preprepared state.
	state = simulateMaliciousCN(t, 4, 3)
	assert.Equal(t, StatePreprepared, state)
}

// Simulate chain split depending on the number of numValidators
func simulateChainSplit(t *testing.T, numValidators int) (State, State) {
	if testing.Verbose() {
		enableLog()
	}

	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	// Note that genValidators(n) will generate n/3 validators.
	// We want n validators, thus calling genValidators(3n).
	validatorAddrs, validatorKeyMap := genValidators(numValidators * 3)

	// Add more EXPECT()s to remove unexpected call error
	mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
	mockBackend.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockBackend.EXPECT().HasBadProposal(gomock.Any()).Return(true).AnyTimes()
	defer mockCtrl.Finish()

	var (
		lastProposal, _ = mockBackend.LastProposal()
		lastBlock       = lastProposal.(*types.Block)
		validators      = mockBackend.Validators(lastBlock)
		proposer        = validators.GetProposer()
		proposerKey     = validatorKeyMap[proposer.Address()]
	)

	// Start istanbul core
	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom
	coreProposer := New(mockBackend, istConfig).(*core)
	coreA := New(mockBackend, istConfig).(*core)
	coreB := New(mockBackend, istConfig).(*core)
	require.Nil(t,
		coreProposer.Start(),
		coreA.Start(),
		coreB.Start())
	defer coreProposer.Stop()
	defer coreA.Stop()
	defer coreB.Stop()

	// make two groups
	// the number of group size is (numValidators-1/2) + 1
	// groupA consists of proposer, coreA, unnamed node(s)
	// groupB consists of proposer, coreB, unnamed node(s)
	subList := validators.SubList(lastBlock.Hash(), coreProposer.currentView())
	groupA, groupB := splitSubList(subList, (numValidators-1)/2, proposer.Address())
	groupA = append(groupA, proposer)
	groupB = append(groupB, proposer)

	// Step 1 - the malicious proposer generates two blocks
	proposalA, err := genBlockParams(lastBlock, proposerKey, 0, 0, 1)
	assert.Nil(t, err)

	proposalB, err := genBlockParams(lastBlock, proposerKey, 1000, 10, 1)
	assert.Nil(t, err)

	// Shortcut for sending message `proposal` to core `c`
	sendMessages := func(state uint64, proposal *types.Block, CNList []istanbul.Validator, c *core) {
		for _, val := range CNList {
			valKey := validatorKeyMap[val.Address()]
			if state == msgPreprepare {
				istanbulMsg, _ := genIstanbulMsg(state, lastBlock.Hash(), proposal, proposer.Address(), valKey)
				err = c.handleMsg(istanbulMsg.Payload)
			} else {
				istanbulMsg, _ := genIstanbulMsg(state, lastBlock.Hash(), proposal, val.Address(), valKey)
				err = c.handleMsg(istanbulMsg.Payload)
			}
			if err != nil {
				t.Logf("handleMsg error: %s", err)
			}
		}
	}
	// Step 2 - exchange consensus messages inside each group

	// the proposer sends two different blocks to each group
	// each group receives a block and handles the message
	// when chain split occurs, their states become StateCommitted
	// otherwise, their states stay StatePreprepared
	sendMessages(msgPreprepare, proposalA, groupA, coreA)
	sendMessages(msgPrepare, proposalA, groupA, coreA)
	if coreA.state.Cmp(StatePrepared) == 0 {
		sendMessages(msgCommit, proposalA, groupA, coreA)
	}

	sendMessages(msgPreprepare, proposalB, groupB, coreB)
	sendMessages(msgPrepare, proposalB, groupB, coreB)
	if coreB.state.Cmp(StatePrepared) == 0 {
		sendMessages(msgCommit, proposalB, groupB, coreB)
	}

	return coreA.state, coreB.state
}

// TestCore_chainSplit tests whether a chain split occurs in a certain conditions:
// 1) the number of validators does not consist of 3f+1;
//     e.g. if the number of validator is 5, it consists of 3f+2 (f=1)
// 2) the proposer is malicious; it sends two different blocks to each group
func TestCore_chainSplit(t *testing.T) {
	// If the number of validators is not 3f+1, the chain can be split.
	stateA, stateB := simulateChainSplit(t, 5)
	assert.Equal(t, StateCommitted, stateA)
	assert.Equal(t, StateCommitted, stateB)

	// If the number of validators is 3f+1, the chain cannot be split.
	stateA, stateB = simulateChainSplit(t, 7)
	fmt.Println(stateA, stateB)
	assert.Equal(t, StatePreprepared, stateA)
	assert.Equal(t, StatePreprepared, stateB)
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
