// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package sc

import (
	"context"
	"errors"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	bridgecontract "github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/node/sc/bridgepool"
	"github.com/klaytn/klaytn/ser/rlp"
	"io"
	"math/big"
	"path"
	"sync"
	"time"
)

const (
	TokenEventChanSize  = 10000
	BridgeAddrJournal   = "bridge_addrs.rlp"
	maxPendingNonceDiff = 20000 // TODO-Klaytn-ServiceChain: update this limitation. Currently, 2 * 10000 TPS.
)

const (
	KLAY uint8 = iota
	TOKEN
	NFT
)

var (
	errNoBridgeInfo         = errors.New("bridge information does not exist")
	errDuplicatedBridgeInfo = errors.New("bridge information is duplicated")
	errNoRecovery           = errors.New("recovery does not exist")
	errAlreadySubscribed    = errors.New("already subscribed")
)

// RequestValueTransferEvent from Bridge contract
type RequestValueTransferEvent struct {
	TokenType    uint8
	ContractAddr common.Address
	TokenAddr    common.Address
	From         common.Address
	To           common.Address
	Amount       *big.Int // Amount is UID in NFT
	RequestNonce uint64
	BlockNumber  uint64
	txHash       common.Hash
}

func (rEv *RequestValueTransferEvent) Nonce() uint64 {
	return rEv.RequestNonce
}

// HandleValueTransferEvent from Bridge contract
type HandleValueTransferEvent struct {
	TokenType    uint8
	ContractAddr common.Address
	TokenAddr    common.Address
	Owner        common.Address
	Amount       *big.Int // Amount is UID in NFT
	HandleNonce  uint64
}

type BridgeJournal struct {
	LocalAddress  common.Address `json:"localAddress"`
	RemoteAddress common.Address `json:"remoteAddress"`
	Subscribed    bool           `json:"subscribed"`
}

type BridgeInfo struct {
	// TODO-Klaytn need to remove and replace subBridge in BridgeInfo
	// subBridge is used for only AddressManager().GetCounterPartToken.
	// Token pair information will be included by BridgeInfo instead of subBridge
	subBridge *SubBridge

	address            common.Address
	counterpartAddress common.Address // TODO-Klaytn need to set counterpart
	account            *accountInfo
	bridge             *bridgecontract.Bridge
	counterpartBridge  *bridgecontract.Bridge
	onServiceChain     bool
	subscribed         bool

	pendingRequestEvent *bridgepool.EventSortedMap // TODO-Klaytn Need to consider the nonce overflow(priority queue?) and the size overflow.
	nextHandleNonce     uint64                     // This nonce will be used for getting pending request value transfer events.

	isRunning                   bool
	handleNonce                 uint64 // the nonce from the handle value transfer event from the bridge.
	requestNonceFromCounterPart uint64 // the nonce from the request value transfer event from the counter part bridge.
	requestNonce                uint64 // the nonce from the request value transfer event from the counter part bridge.

	newEvent chan struct{}
	closed   chan struct{}
}

func NewBridgeInfo(subBridge *SubBridge, addr common.Address, bridge *bridgecontract.Bridge, cpAddr common.Address, cpBridge *bridgecontract.Bridge, account *accountInfo, local, subscribed bool) *BridgeInfo {
	bi := &BridgeInfo{
		subBridge,
		addr,
		cpAddr,
		account,
		bridge,
		cpBridge,
		local,
		subscribed,
		bridgepool.NewEventSortedMap(),
		0,
		true,
		0,
		0,
		0,
		make(chan struct{}),
		make(chan struct{}),
	}

	if err := bi.UpdateInfo(); err != nil {
		logger.Error("NewBridgeInfo can't be updated.", "err", err)
	}

	bi.nextHandleNonce = bi.handleNonce
	go bi.loop()

	return bi
}

func (bi *BridgeInfo) loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	logger.Info("start bridge loop", "addr", bi.address.String(), "onServiceChain", bi.onServiceChain)

	for {
		select {
		case <-bi.newEvent:
			bi.processingPendingRequestEvents()

		case <-ticker.C:
			bi.processingPendingRequestEvents()

		case <-bi.closed:
			logger.Info("stop bridge loop", "addr", bi.address.String(), "onServiceChain", bi.onServiceChain)
			return
		}
	}
}

func (bi *BridgeInfo) GetPendingRequestEvents(start uint64) []*RequestValueTransferEvent {
	ready := bi.pendingRequestEvent.Ready(start)
	var readyEvent []*RequestValueTransferEvent
	for _, item := range ready {
		readyEvent = append(readyEvent, item.(*RequestValueTransferEvent))
	}

	return readyEvent
}

// processingPendingRequestEvents handles pending request value transfer events of the bridge.
func (bi *BridgeInfo) processingPendingRequestEvents() error {
	ReadyEvent := bi.GetReadyRequestValueTransferEvents()
	if ReadyEvent == nil {
		return nil
	}

	logger.Debug("Get Pending request value transfer event", "len(readyEvent)", len(ReadyEvent), "nextHandleNonce", bi.nextHandleNonce, "len(pendingEvent)", bi.pendingRequestEvent.Len())

	diff := bi.requestNonceFromCounterPart - bi.handleNonce
	if diff > errorDiffRequestHandleNonce {
		logger.Error("Value transfer requested/handled nonce gap is too much.", "toSC", bi.onServiceChain, "diff", diff, "requestedNonce", bi.requestNonceFromCounterPart, "handledNonce", bi.handleNonce)
		// TODO-Klaytn need to consider starting value transfer recovery.
	}

	for idx, ev := range ReadyEvent {
		if ev.RequestNonce < bi.handleNonce {
			logger.Trace("past requests are ignored", "RequestNonce", ev.RequestNonce, "handleNonce", bi.handleNonce)
			continue
		}

		if ev.RequestNonce > bi.handleNonce+maxPendingNonceDiff {
			logger.Trace("nonce diff is too large", "limitation", maxPendingNonceDiff)
			return errors.New("nonce diff is too large")
		}

		if err := bi.handleRequestValueTransferEvent(ev); err != nil {
			bi.AddRequestValueTransferEvents(ReadyEvent[idx:])
			logger.Error("Failed handle request value transfer event", "err", err, "len(RePutEvent)", len(ReadyEvent[idx:]))
			return err
		}
		if bi.nextHandleNonce <= ev.RequestNonce {
			bi.nextHandleNonce = ev.RequestNonce + 1
		}
	}

	return nil
}

func (bi *BridgeInfo) UpdateInfo() error {
	if bi == nil {
		return errNoBridgeInfo
	}

	rn, err := bi.bridge.RequestNonce(nil)
	if err != nil {
		return err
	}
	bi.UpdateRequestNonce(rn)

	hn, err := bi.bridge.HandleNonce(nil)
	if err != nil {
		return err
	}
	bi.UpdateHandledNonce(hn)
	bi.UpdateRequestNonceFromCounterpart(hn)

	isRunning, err := bi.bridge.IsRunning(nil)
	if err != nil {
		return err
	}
	bi.isRunning = isRunning

	counterPart, err := bi.bridge.CounterpartBridge(nil)
	if err != nil {
		return err
	}
	bi.counterpartAddress = counterPart

	return nil
}

// handleRequestValueTransferEvent handles the given request value transfer event.
func (bi *BridgeInfo) handleRequestValueTransferEvent(ev *RequestValueTransferEvent) error {
	tokenType := ev.TokenType
	tokenAddr := bi.subBridge.AddressManager().GetCounterPartToken(ev.TokenAddr)
	if tokenType != KLAY && tokenAddr == (common.Address{}) {
		logger.Warn("Unregistered counter part token address.", "addr", tokenAddr.Hex())
		ctTokenAddr, err := bi.counterpartBridge.AllowedTokens(nil, ev.TokenAddr)
		if err != nil {
			return err
		}
		if ctTokenAddr == (common.Address{}) {
			return errors.New("can't get counterpart token from bridge")
		}

		if err := bi.subBridge.AddressManager().AddToken(ev.TokenAddr, ctTokenAddr); err != nil {
			return err
		}
		tokenAddr = ctTokenAddr
		logger.Info("Register counter part token address.", "addr", tokenAddr.Hex(), "cpAddr", ctTokenAddr.Hex())
	}

	to := ev.To

	bridgeAcc := bi.account

	bridgeAcc.Lock()
	defer bridgeAcc.UnLock()

	auth := bridgeAcc.GetTransactOpts()

	var handleTx *types.Transaction
	var err error

	switch tokenType {
	case KLAY:
		handleTx, err = bi.bridge.HandleKLAYTransfer(auth, ev.Amount, to, ev.RequestNonce, ev.BlockNumber)
		if err != nil {
			return err
		}
		logger.Trace("Bridge succeeded to HandleKLAYTransfer", "nonce", ev.RequestNonce, "tx", handleTx.Hash().String())

	case TOKEN:
		handleTx, err = bi.bridge.HandleTokenTransfer(auth, ev.Amount, to, tokenAddr, ev.RequestNonce, ev.BlockNumber)
		if err != nil {
			return err
		}
		logger.Trace("Bridge succeeded to HandleTokenTransfer", "nonce", ev.RequestNonce, "tx", handleTx.Hash().String())
	case NFT:
		handleTx, err = bi.bridge.HandleNFTTransfer(auth, ev.Amount, to, tokenAddr, ev.RequestNonce, ev.BlockNumber)
		if err != nil {
			return err
		}
		logger.Trace("Bridge succeeded to HandleNFTTransfer", "nonce", ev.RequestNonce, "tx", handleTx.Hash().String())
	default:
		logger.Warn("Got Unknown Token Type ReceivedEvent", "bridge", ev.ContractAddr, "nonce", ev.RequestNonce, "from", ev.From)
		return nil
	}

	bridgeAcc.IncNonce()

	bi.subBridge.ChainDB().WriteHandleTxHashFromRequestTxHash(ev.txHash, handleTx.Hash())
	return nil
}

// UpdateRequestNonceFromCounterpart updates the request nonce from counterpart bridge.
func (bi *BridgeInfo) UpdateRequestNonceFromCounterpart(nonce uint64) {
	if bi.requestNonceFromCounterPart < nonce {
		vtRequestNonceCount.Inc(int64(nonce - bi.requestNonceFromCounterPart))
		bi.requestNonceFromCounterPart = nonce
	}
}

// UpdateRequestNonce updates the request nonce of the bridge.
func (bi *BridgeInfo) UpdateRequestNonce(nonce uint64) {
	if bi.requestNonce < nonce {
		bi.requestNonce = nonce
	}
}

// UpdateHandledNonce updates the handled nonce with new nonce.
func (bi *BridgeInfo) UpdateHandledNonce(nonce uint64) {
	if bi.handleNonce < nonce {
		vtHandleNonceCount.Inc(int64(nonce - bi.handleNonce))
		bi.handleNonce = nonce
	}
}

// AddRequestValueTransferEvents adds events into the pendingRequestEvent.
func (bi *BridgeInfo) AddRequestValueTransferEvents(evs []*RequestValueTransferEvent) {
	// TODO-Klaytn Need to consider the nonce overflow(priority queue?) and the size overflow.
	// - If the size is full and received event has the omitted nonce, it can be allowed.
	for _, ev := range evs {
		bi.UpdateRequestNonceFromCounterpart(ev.RequestNonce + 1)
		bi.pendingRequestEvent.Put(ev)
	}

	vtPendingRequestEventMeter.Inc(int64(len(evs)))

	select {
	case bi.newEvent <- struct{}{}:
	default:
	}
}

// GetReadyRequestValueTransferEvents returns the processable events with the increasing nonce.
func (bi *BridgeInfo) GetReadyRequestValueTransferEvents() []*RequestValueTransferEvent {
	return bi.GetPendingRequestEvents(bi.nextHandleNonce)
}

// DecodeRLP decodes the Klaytn
func (b *BridgeJournal) DecodeRLP(s *rlp.Stream) error {
	var elem struct {
		LocalAddress  common.Address
		RemoteAddress common.Address
		Paired        bool
	}
	if err := s.Decode(&elem); err != nil {
		return err
	}
	b.LocalAddress, b.RemoteAddress, b.Subscribed = elem.LocalAddress, elem.RemoteAddress, elem.Paired
	return nil
}

// EncodeRLP serializes a BridgeJournal into the Klaytn RLP BridgeJournal format.
func (b *BridgeJournal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		b.LocalAddress,
		b.RemoteAddress,
		b.Subscribed,
	})
}

// BridgeManager manages Bridge SmartContracts
// for value transfer between service chain and parent chain
type BridgeManager struct {
	subBridge *SubBridge

	receivedEvents map[common.Address]event.Subscription
	withdrawEvents map[common.Address]event.Subscription
	bridges        map[common.Address]*BridgeInfo
	mu             sync.RWMutex

	tokenReceived event.Feed
	tokenWithdraw event.Feed

	scope event.SubscriptionScope

	journal    *bridgeAddrJournal
	recoveries map[common.Address]*valueTransferRecovery
	auth       *bind.TransactOpts
}

func NewBridgeManager(main *SubBridge) (*BridgeManager, error) {
	bridgeAddrJournal := newBridgeAddrJournal(path.Join(main.config.DataDir, BridgeAddrJournal))

	bridgeManager := &BridgeManager{
		subBridge:      main,
		receivedEvents: make(map[common.Address]event.Subscription),
		withdrawEvents: make(map[common.Address]event.Subscription),
		bridges:        make(map[common.Address]*BridgeInfo),
		journal:        bridgeAddrJournal,
		recoveries:     make(map[common.Address]*valueTransferRecovery),
	}

	logger.Info("Load Bridge Address from JournalFiles ", "path", bridgeManager.journal.path)
	bridgeManager.journal.cache = make(map[common.Address]*BridgeJournal)

	if err := bridgeManager.journal.load(func(gwjournal BridgeJournal) error {
		logger.Info("Load Bridge Address from JournalFiles ",
			"local address", gwjournal.LocalAddress.Hex(), "remote address", gwjournal.RemoteAddress.Hex())
		bridgeManager.journal.cache[gwjournal.LocalAddress] = &gwjournal
		return nil
	}); err != nil {
		logger.Error("fail to load bridge address", "err", err)
	}

	if err := bridgeManager.journal.rotate(bridgeManager.GetAllBridge()); err != nil {
		logger.Error("fail to rotate bridge journal", "err", err)
	}

	return bridgeManager, nil
}

// LogBridgeStatus logs the bridge contract requested/handled nonce status as an information.
func (bm *BridgeManager) LogBridgeStatus() {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if len(bm.bridges) == 0 {
		return
	}

	m2sTotalRequestNonce, m2sTotalHandleNonce := uint64(0), uint64(0)
	s2mTotalRequestNonce, s2mTotalHandleNonce := uint64(0), uint64(0)

	for bAddr, b := range bm.bridges {
		diffNonce := b.requestNonceFromCounterPart - b.handleNonce

		if b.subscribed {
			var headStr string
			if b.onServiceChain {
				headStr = "Bridge(Main -> Service Chain)"
				m2sTotalRequestNonce += b.requestNonceFromCounterPart
				m2sTotalHandleNonce += b.handleNonce
			} else {
				headStr = "Bridge(Service -> Main Chain)"
				s2mTotalRequestNonce += b.requestNonceFromCounterPart
				s2mTotalHandleNonce += b.handleNonce
			}
			logger.Debug(headStr, "bridge", bAddr.String(), "requestNonce", b.requestNonceFromCounterPart, "handleNonce", b.handleNonce, "pending", diffNonce)
		}
	}

	logger.Info("VT : Main -> Service Chain", "request", m2sTotalRequestNonce, "handle", m2sTotalHandleNonce, "pending", m2sTotalRequestNonce-m2sTotalHandleNonce)
	logger.Info("VT : Service -> Main Chain", "request", s2mTotalRequestNonce, "handle", s2mTotalHandleNonce, "pending", s2mTotalRequestNonce-s2mTotalHandleNonce)
}

// SubscribeTokenReceived registers a subscription of TokenReceivedEvent.
func (bm *BridgeManager) SubscribeTokenReceived(ch chan<- RequestValueTransferEvent) event.Subscription {
	return bm.scope.Track(bm.tokenReceived.Subscribe(ch))
}

// SubscribeTokenWithDraw registers a subscription of TokenTransferEvent.
func (bm *BridgeManager) SubscribeTokenWithDraw(ch chan<- HandleValueTransferEvent) event.Subscription {
	return bm.scope.Track(bm.tokenWithdraw.Subscribe(ch))
}

// GetAllBridge returns a slice of journal cache.
func (bm *BridgeManager) GetAllBridge() []*BridgeJournal {
	var gwjs []*BridgeJournal

	for _, journal := range bm.journal.cache {
		gwjs = append(gwjs, journal)
	}
	return gwjs
}

// GetBridge returns bridge contract of the specified address.
func (bm *BridgeManager) GetBridgeInfo(addr common.Address) (*BridgeInfo, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	bridge, ok := bm.bridges[addr]
	return bridge, ok
}

// DeleteBridgeInfo deletes the bridge info of the specified address.
func (bm *BridgeManager) DeleteBridgeInfo(addr common.Address) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	bi := bm.bridges[addr]
	if bi == nil {
		return errNoBridgeInfo
	}

	close(bi.closed)

	delete(bm.bridges, addr)
	return nil
}

// SetBridgeInfo stores the address and bridge pair with local/remote and subscription status.
func (bm *BridgeManager) SetBridgeInfo(addr common.Address, bridge *bridgecontract.Bridge, cpAddr common.Address, cpBridge *bridgecontract.Bridge, account *accountInfo, local bool, subscribed bool) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if bm.bridges[addr] != nil {
		return errDuplicatedBridgeInfo
	}
	bm.bridges[addr] = NewBridgeInfo(bm.subBridge, addr, bridge, cpAddr, cpBridge, account, local, subscribed)
	return nil
}

// RestoreBridges setups bridge subscription by using the journal cache.
func (bm *BridgeManager) RestoreBridges() error {
	bm.stopAllRecoveries()

	for _, journal := range bm.journal.cache {
		cBridgeAddr := journal.LocalAddress
		pBridgeAddr := journal.RemoteAddress
		cBridge, err := bridgecontract.NewBridge(cBridgeAddr, bm.subBridge.localBackend)
		if err != nil {
			logger.Error("local bridge creation is failed", err)
			continue
		}
		pBridge, err := bridgecontract.NewBridge(pBridgeAddr, bm.subBridge.remoteBackend)
		if err != nil {
			logger.Error("remote bridge creation is failed", err)
			continue
		}
		bam := bm.subBridge.bridgeAccountManager
		err = bm.SetBridgeInfo(cBridgeAddr, cBridge, pBridgeAddr, pBridge, bam.scAccount, true, false)
		if err != nil {
			logger.Error("setting local bridge info is failed", err)
			continue
		}
		err = bm.SetBridgeInfo(pBridgeAddr, pBridge, cBridgeAddr, cBridge, bam.mcAccount, false, false)
		if err != nil {
			logger.Error("setting remote bridge info is failed", err)
			continue
		}
		bm.subBridge.AddressManager().AddBridge(cBridgeAddr, pBridgeAddr)

		if journal.Subscribed {
			logger.Info("automatic bridge subscription", "local", cBridgeAddr, "remote", pBridgeAddr)
			if err := bm.subscribeEvent(cBridgeAddr, cBridge); err != nil {
				logger.Error("local bridge subscription is failed", err)
				continue
			}
			if err := bm.subscribeEvent(pBridgeAddr, pBridge); err != nil {
				// TODO-Klaytn need to consider how to retry.
				bm.subBridge.AddressManager().DeleteBridge(cBridgeAddr)
				bm.UnsubscribeEvent(cBridgeAddr)
				logger.Error("local bridge subscription is failed", err)
				continue
			}
			bm.AddRecovery(cBridgeAddr, pBridgeAddr)
		}
	}
	return nil
}

// SetJournal inserts or updates journal for a given addresses pair.
func (bm *BridgeManager) SetJournal(localAddress, remoteAddress common.Address) error {
	err := bm.journal.insert(localAddress, remoteAddress)
	return err
}

// AddRecovery starts value transfer recovery for a given addresses pair.
func (bm *BridgeManager) AddRecovery(localAddress, remoteAddress common.Address) error {
	if !bm.subBridge.config.VTRecovery {
		logger.Info("value transfer recovery is disabled")
		return nil
	}

	// Check if bridge information is exist.
	localBridgeInfo, ok := bm.GetBridgeInfo(localAddress)
	if !ok {
		return errNoBridgeInfo
	}
	remoteBridgeInfo, ok := bm.GetBridgeInfo(remoteAddress)
	if !ok {
		return errNoBridgeInfo
	}

	// Create and start value transfer recovery.
	recovery := NewValueTransferRecovery(bm.subBridge.config, localBridgeInfo, remoteBridgeInfo)
	recovery.Start()
	bm.recoveries[localAddress] = recovery // suppose local/remote is always a pair.
	return nil
}

// DeleteRecovery deletes the journal and stop the value transfer recovery for a given address pair.
func (bm *BridgeManager) DeleteRecovery(localAddress, remoteAddress common.Address) error {
	// Stop the recovery.
	recovery, ok := bm.recoveries[localAddress]
	if !ok {
		return errNoRecovery
	}
	recovery.Stop()
	delete(bm.recoveries, localAddress)

	return nil
}

// stopAllRecoveries stops the internal value transfer recoveries.
func (bm *BridgeManager) stopAllRecoveries() {
	for _, recovery := range bm.recoveries {
		recovery.Stop()
	}
	bm.recoveries = make(map[common.Address]*valueTransferRecovery)
}

// Deploy Bridge SmartContract on same node or remote node
func (bm *BridgeManager) DeployBridge(backend bind.ContractBackend, local bool) (*bridgecontract.Bridge, common.Address, error) {
	var acc *accountInfo
	if local {
		acc = bm.subBridge.bridgeAccountManager.scAccount
	} else {
		acc = bm.subBridge.bridgeAccountManager.mcAccount
	}

	addr, bridge, err := bm.deployBridge(acc, backend)
	if err != nil {
		return nil, common.Address{}, err
	}

	return bridge, addr, err
}

// DeployBridge handles actual smart contract deployment.
// To create contract, the chain ID, nonce, account key, private key, contract binding and gas price are used.
// The deployed contract address, transaction are returned. An error is also returned if any.
func (bm *BridgeManager) deployBridge(acc *accountInfo, backend bind.ContractBackend) (common.Address, *bridgecontract.Bridge, error) {
	acc.Lock()
	auth := acc.GetTransactOpts()
	addr, tx, contract, err := bridgecontract.DeployBridge(auth, backend)
	if err != nil {
		logger.Error("Failed to deploy contract.", "err", err)
		acc.UnLock()
		return common.Address{}, nil, err
	}
	acc.IncNonce()
	acc.UnLock()

	logger.Info("Bridge is deploying...", "addr", addr, "txHash", tx.Hash().String())

	back, ok := backend.(bind.DeployBackend)
	if !ok {
		logger.Warn("DeployBacked type assertion is failed. Skip WaitDeployed.")
		return addr, contract, nil
	}

	timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelTimeout()

	addr, err = bind.WaitDeployed(timeoutContext, back, tx)
	if err != nil {
		logger.Error("Failed to WaitDeployed.", "err", err, "txHash", tx.Hash().String())
		return common.Address{}, nil, err
	}
	logger.Info("Bridge is deployed.", "addr", addr, "txHash", tx.Hash().String())
	return addr, contract, nil
}

// SubscribeEvent registers a subscription of BridgeERC20Received and BridgeTokenWithdrawn
func (bm *BridgeManager) SubscribeEvent(addr common.Address) error {
	bridgeInfo, ok := bm.GetBridgeInfo(addr)
	if !ok {
		return errNoBridgeInfo
	}
	if bridgeInfo.subscribed {
		return errAlreadySubscribed
	}
	err := bm.subscribeEvent(addr, bridgeInfo.bridge)
	if err != nil {
		return err
	}

	return nil
}

// resetAllSubscribedEvents resets watch logs and recreates a goroutine loop to handle event messages.
func (bm *BridgeManager) ResetAllSubscribedEvents() error {
	logger.Info("ResetAllSubscribedEvents is called.")

	for _, journal := range bm.journal.cache {
		if journal.Subscribed {
			bm.UnsubscribeEvent(journal.LocalAddress)
			bm.UnsubscribeEvent(journal.RemoteAddress)

			bridgeInfo, ok := bm.GetBridgeInfo(journal.LocalAddress)
			if !ok {
				logger.Error("ResetAllSubscribedEvents failed to GetBridgeInfo", "localBridge", journal.LocalAddress.String())
				continue
			}
			bm.subscribeEvent(journal.LocalAddress, bridgeInfo.bridge)
			bridgeInfo, ok = bm.GetBridgeInfo(journal.RemoteAddress)
			if !ok {
				logger.Error("ResetAllSubscribedEvents failed to GetBridgeInfo", "remoteBridge", journal.RemoteAddress.String())
				bm.UnsubscribeEvent(journal.LocalAddress)
				continue
			}
			bm.subscribeEvent(journal.RemoteAddress, bridgeInfo.bridge)
		}
	}
	return nil
}

// SubscribeEvent sets watch logs and creates a goroutine loop to handle event messages.
func (bm *BridgeManager) subscribeEvent(addr common.Address, bridge *bridgecontract.Bridge) error {
	tokenReceivedCh := make(chan *bridgecontract.BridgeRequestValueTransfer, TokenEventChanSize)
	tokenWithdrawCh := make(chan *bridgecontract.BridgeHandleValueTransfer, TokenEventChanSize)

	receivedSub, err := bridge.WatchRequestValueTransfer(nil, tokenReceivedCh)
	if err != nil {
		logger.Error("Failed to pBridge.WatchERC20Received", "err", err)
		return err
	}
	bm.receivedEvents[addr] = receivedSub
	withdrawnSub, err := bridge.WatchHandleValueTransfer(nil, tokenWithdrawCh)
	if err != nil {
		logger.Error("Failed to pBridge.WatchTokenWithdrawn", "err", err)
		receivedSub.Unsubscribe()
		delete(bm.receivedEvents, addr)
		return err
	}
	bm.withdrawEvents[addr] = withdrawnSub
	bridgeInfo, ok := bm.GetBridgeInfo(addr)
	if !ok {
		receivedSub.Unsubscribe()
		withdrawnSub.Unsubscribe()
		delete(bm.receivedEvents, addr)
		delete(bm.withdrawEvents, addr)
		return errNoBridgeInfo
	}
	bridgeInfo.subscribed = true

	go bm.loop(addr, tokenReceivedCh, tokenWithdrawCh, bm.scope.Track(receivedSub), bm.scope.Track(withdrawnSub))

	return nil
}

// UnsubscribeEvent cancels the contract's watch logs and initializes the status.
func (bm *BridgeManager) UnsubscribeEvent(addr common.Address) {
	receivedSub := bm.receivedEvents[addr]
	if receivedSub != nil {
		receivedSub.Unsubscribe()
		delete(bm.receivedEvents, addr)
	}

	withdrawSub := bm.withdrawEvents[addr]
	if withdrawSub != nil {
		withdrawSub.Unsubscribe()
		delete(bm.withdrawEvents, addr)
	}

	bridgeInfo, ok := bm.GetBridgeInfo(addr)
	if ok {
		bridgeInfo.subscribed = false
	}
}

// Loop handles subscribed event messages.
func (bm *BridgeManager) loop(
	addr common.Address,
	receivedCh <-chan *bridgecontract.BridgeRequestValueTransfer,
	withdrawCh <-chan *bridgecontract.BridgeHandleValueTransfer,
	receivedSub event.Subscription,
	withdrawSub event.Subscription) {

	defer receivedSub.Unsubscribe()
	defer withdrawSub.Unsubscribe()

	// TODO-Klaytn change goroutine logic for performance
	for {
		select {
		case ev := <-receivedCh:
			receiveEvent := RequestValueTransferEvent{
				TokenType:    ev.Kind,
				ContractAddr: addr,
				TokenAddr:    ev.ContractAddress,
				From:         ev.From,
				To:           ev.To,
				Amount:       ev.Amount,
				RequestNonce: ev.RequestNonce,
				BlockNumber:  ev.Raw.BlockNumber,
				txHash:       ev.Raw.TxHash,
			}
			bm.tokenReceived.Send(receiveEvent)
		case ev := <-withdrawCh:
			withdrawEvent := HandleValueTransferEvent{
				TokenType:    ev.Kind,
				ContractAddr: addr,
				TokenAddr:    ev.ContractAddress,
				Owner:        ev.Owner,
				Amount:       ev.Value,
				HandleNonce:  ev.HandleNonce,
			}
			bm.tokenWithdraw.Send(withdrawEvent)
		case err := <-receivedSub.Err():
			logger.Info("Contract Event Loop Running Stop by receivedSub.Err()", "err", err)
			return
		case err := <-withdrawSub.Err():
			logger.Info("Contract Event Loop Running Stop by withdrawSub.Err()", "err", err)
			return
		}
	}
}

// Stop closes a subscribed event scope of the bridge manager.
func (bm *BridgeManager) Stop() {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	for _, bi := range bm.bridges {
		close(bi.closed)
	}

	bm.scope.Close()
}
