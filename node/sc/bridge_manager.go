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
	"io"
	"math/big"
	"path"
	"sync"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	bridgecontract "github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/node/sc/bridgepool"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
)

const (
	TokenEventChanSize  = 10000
	BridgeAddrJournal   = "bridge_addrs.rlp"
	maxPendingNonceDiff = 1000 // TODO-Klaytn-ServiceChain: update this limitation. Currently, 2 * 500 TPS.

	maxHandledEventSize = 10000000
)

const (
	KLAY uint8 = iota
	ERC20
	ERC721
)

const (
	voteTypeValueTransfer = 0
	voteTypeConfiguration = 1
)

var (
	ErrInvalidTokenPair     = errors.New("invalid token pair")
	ErrNoBridgeInfo         = errors.New("bridge information does not exist")
	ErrDuplicatedBridgeInfo = errors.New("bridge information is duplicated")
	ErrDuplicatedToken      = errors.New("token is duplicated")
	ErrNoRecovery           = errors.New("recovery does not exist")
	ErrAlreadySubscribed    = errors.New("already subscribed")
	ErrBridgeRestore        = errors.New("restoring bridges is failed")
)

// HandleValueTransferEvent from Bridge contract
type HandleValueTransferEvent struct {
	*bridgecontract.BridgeHandleValueTransfer
}

type BridgeJournal struct {
	ChildAddress  common.Address `json:"childAddress"`
	ParentAddress common.Address `json:"parentAddress"`
	Subscribed    bool           `json:"subscribed"`
}

type BridgeInfo struct {
	subBridge *SubBridge
	bridgeDB  database.DBManager

	counterpartBackend Backend
	address            common.Address
	counterpartAddress common.Address // TODO-Klaytn need to set counterpart
	account            *accountInfo
	bridge             *bridgecontract.Bridge
	counterpartBridge  *bridgecontract.Bridge
	onChildChain       bool
	subscribed         bool

	counterpartToken map[common.Address]common.Address

	pendingRequestEvent *bridgepool.ItemSortedMap

	isRunning                   bool
	handleNonce                 uint64 // the nonce from the handle value transfer event from the bridge.
	lowerHandleNonce            uint64 // the lower handle nonce from the bridge.
	requestNonceFromCounterPart uint64 // the nonce from the request value transfer event from the counter part bridge.
	requestNonce                uint64 // the nonce from the request value transfer event from the counter part bridge.

	newEvent chan struct{}
	closed   chan struct{}

	handledEvent *bridgepool.ItemSortedMap
}

type requestEvent struct {
	nonce uint64
}

func (ev requestEvent) Nonce() uint64 {
	return ev.nonce
}

func NewBridgeInfo(sb *SubBridge, addr common.Address, bridge *bridgecontract.Bridge, cpAddr common.Address, cpBridge *bridgecontract.Bridge, account *accountInfo, local, subscribed bool, cpBackend Backend) (*BridgeInfo, error) {
	bi := &BridgeInfo{
		sb,
		sb.chainDB,
		cpBackend,
		addr,
		cpAddr,
		account,
		bridge,
		cpBridge,
		local,
		subscribed,
		make(map[common.Address]common.Address),
		bridgepool.NewItemSortedMap(bridgepool.UnlimitedItemSortedMap),
		true,
		0,
		0,
		0,
		0,
		make(chan struct{}),
		make(chan struct{}),
		bridgepool.NewItemSortedMap(maxHandledEventSize),
	}

	if err := bi.UpdateInfo(); err != nil {
		return bi, err
	}

	go bi.loop()

	return bi, nil
}

func (bi *BridgeInfo) loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	logger.Info("start bridge loop", "addr", bi.address.String(), "onChildChain", bi.onChildChain)

	for {
		select {
		case <-bi.newEvent:
			bi.processingPendingRequestEvents()

		case <-ticker.C:
			bi.processingPendingRequestEvents()

		case <-bi.closed:
			logger.Info("stop bridge loop", "addr", bi.address.String(), "onChildChain", bi.onChildChain)
			return
		}
	}
}

func (bi *BridgeInfo) RegisterToken(token, counterpartToken common.Address) error {
	_, exist := bi.counterpartToken[token]
	if exist {
		return ErrDuplicatedToken
	}
	bi.counterpartToken[token] = counterpartToken
	return nil
}

func (bi *BridgeInfo) DeregisterToken(token, counterpartToken common.Address) error {
	_, exist := bi.counterpartToken[token]
	if !exist {
		return ErrInvalidTokenPair
	}
	delete(bi.counterpartToken, token)
	return nil
}

func (bi *BridgeInfo) GetCounterPartToken(token common.Address) common.Address {
	cpToken, exist := bi.counterpartToken[token]
	if !exist {
		return common.Address{}
	}
	return cpToken
}

func (bi *BridgeInfo) GetPendingRequestEvents() []IRequestValueTransferEvent {
	ready := bi.pendingRequestEvent.Pop(maxPendingNonceDiff / 2)
	readyEvent := make([]IRequestValueTransferEvent, len(ready))
	for i, item := range ready {
		readyEvent[i] = item.(IRequestValueTransferEvent)
	}
	vtPendingRequestEventCounter.Dec((int64)(len(ready)))
	return readyEvent
}

// processingPendingRequestEvents handles pending request value transfer events of the bridge.
func (bi *BridgeInfo) processingPendingRequestEvents() {
	ReadyEvent := bi.GetReadyRequestValueTransferEvents()
	if ReadyEvent == nil {
		return
	}

	logger.Trace("Get ready request value transfer event", "len(readyEvent)", len(ReadyEvent), "len(pendingEvent)", bi.pendingRequestEvent.Len())

	for idx, ev := range ReadyEvent {
		if ev.GetRequestNonce() < bi.lowerHandleNonce || bi.handledEvent.Exist(ev.GetRequestNonce()) {
			logger.Trace("handled requests can be ignored", "RequestNonce", ev.GetRequestNonce(), "lowerHandleNonce", bi.lowerHandleNonce)
			continue
		}

		if err := bi.handleRequestValueTransferEvent(ev); err != nil {
			bi.AddRequestValueTransferEvents(ReadyEvent[idx:])
			logger.Error("Failed handle request value transfer event", "err", err, "len(RePutEvent)", len(ReadyEvent[idx:]))
			return
		}
	}
}

func (bi *BridgeInfo) UpdateInfo() error {
	if bi == nil {
		return ErrNoBridgeInfo
	}

	rn, err := bi.bridge.RequestNonce(nil)
	if err != nil {
		return err
	}
	bi.SetRequestNonce(rn)

	hn, err := bi.bridge.LowerHandleNonce(nil)
	if err != nil {
		return err
	}

	bi.lowerHandleNonce = hn

	bi.SetHandleNonce(hn)
	bi.SetRequestNonceFromCounterpart(hn)

	isRunning, err := bi.bridge.IsRunning(nil)
	if err != nil {
		return err
	}
	bi.isRunning = isRunning

	return nil
}

// handleRequestValueTransferEvent handles the given request value transfer event.
func (bi *BridgeInfo) handleRequestValueTransferEvent(ev IRequestValueTransferEvent) error {
	var (
		tokenType                         = ev.GetTokenType()
		tokenAddr, from, to, contractAddr = ev.GetTokenAddress(), ev.GetFrom(), ev.GetTo(), ev.GetRaw().Address
		txHash                            = ev.GetRaw().TxHash
		valueOrTokenId                    = ev.GetValueOrTokenId()
		requestNonce, blkNumber           = ev.GetRequestNonce(), ev.GetRaw().BlockNumber
		extraData                         = ev.GetExtraData()
	)

	ctpartTokenAddr := bi.GetCounterPartToken(tokenAddr)
	// TODO-Klaytn-Servicechain Add counterpart token address in requestValueTransferEvent
	if tokenType != KLAY && ctpartTokenAddr == (common.Address{}) {
		logger.Warn("Unregistered counter part token address.", "addr", ctpartTokenAddr.Hex())
		ctTokenAddr, err := bi.counterpartBridge.RegisteredTokens(nil, tokenAddr)
		if err != nil {
			return err
		}
		if ctTokenAddr == (common.Address{}) {
			return errors.New("can't get counterpart token from bridge")
		}
		if err := bi.RegisterToken(tokenAddr, ctTokenAddr); err != nil {
			return err
		}
		ctpartTokenAddr = ctTokenAddr
		logger.Info("Register counter part token address.", "addr", ctpartTokenAddr.Hex(), "cpAddr", ctTokenAddr.Hex())
	}

	bridgeAcc := bi.account

	bridgeAcc.Lock()
	defer bridgeAcc.UnLock()

	auth := bridgeAcc.GenerateTransactOpts()

	var handleTx *types.Transaction
	var err error

	switch tokenType {
	case KLAY:
		handleTx, err = bi.bridge.HandleKLAYTransfer(auth, txHash, from, to, valueOrTokenId, requestNonce, blkNumber, extraData)
		if err != nil {
			return err
		}
		logger.Trace("Bridge succeeded to HandleKLAYTransfer", "nonce", requestNonce, "tx", handleTx.Hash().String())
	case ERC20:
		handleTx, err = bi.bridge.HandleERC20Transfer(auth, txHash, from, to, ctpartTokenAddr, valueOrTokenId, requestNonce, blkNumber, extraData)
		if err != nil {
			return err
		}
		logger.Trace("Bridge succeeded to HandleERC20Transfer", "nonce", requestNonce, "tx", handleTx.Hash().String())
	case ERC721:
		handleTx, err = bi.bridge.HandleERC721Transfer(auth, txHash, from, to, ctpartTokenAddr, valueOrTokenId, requestNonce, blkNumber, GetURI(ev), extraData)
		logger.Trace("Bridge succeeded to HandleERC721Transfer", "nonce", requestNonce, "tx", handleTx.Hash().String())
	default:
		logger.Error("Got Unknown Token Type ReceivedEvent", "bridge", contractAddr, "nonce", requestNonce, "from", from)
		return nil
	}

	bridgeAcc.IncNonce()

	bi.bridgeDB.WriteHandleTxHashFromRequestTxHash(txHash, handleTx.Hash())
	return nil
}

// SetRequestNonceFromCounterpart sets the request nonce from counterpart bridge.
func (bi *BridgeInfo) SetRequestNonceFromCounterpart(nonce uint64) {
	if bi.requestNonceFromCounterPart < nonce {
		vtRequestNonceCount.Inc(int64(nonce - bi.requestNonceFromCounterPart))
		bi.requestNonceFromCounterPart = nonce
	}
}

// SetRequestNonce sets the request nonce of the bridge.
func (bi *BridgeInfo) SetRequestNonce(nonce uint64) {
	if bi.requestNonce < nonce {
		bi.requestNonce = nonce
	}
}

// MarkHandledNonce marks the handled nonce and sets the handle nonce value.
func (bi *BridgeInfo) MarkHandledNonce(nonce uint64) {
	bi.SetHandleNonce(nonce + 1)
	bi.handledEvent.Put(requestEvent{nonce})
}

// SetHandleNonce sets the handled nonce with a new nonce.
func (bi *BridgeInfo) SetHandleNonce(nonce uint64) {
	if bi.handleNonce < nonce {
		vtHandleNonceCount.Inc(int64(nonce - bi.handleNonce))
		bi.handleNonce = nonce
	}
}

// UpdateLowerHandleNonce updates the lower handle nonce.
func (bi *BridgeInfo) UpdateLowerHandleNonce(nonce uint64) {
	if bi.lowerHandleNonce < nonce {
		vtLowerHandleNonceCount.Inc(int64(nonce - bi.lowerHandleNonce))
		bi.lowerHandleNonce = nonce

		bi.handledEvent.Forward(nonce)
	}
}

// AddRequestValueTransferEvents adds events into the pendingRequestEvent.
func (bi *BridgeInfo) AddRequestValueTransferEvents(evs []IRequestValueTransferEvent) {
	for _, ev := range evs {
		if bi.pendingRequestEvent.Len() > maxPendingNonceDiff {
			flatten := bi.pendingRequestEvent.Flatten()
			maxNonce := flatten[len(flatten)-1].Nonce()
			if ev.Nonce() >= maxNonce || bi.pendingRequestEvent.Exist(ev.Nonce()) {
				continue
			}
			bi.pendingRequestEvent.Remove(maxNonce)
			vtPendingRequestEventCounter.Dec(1)
			logger.Trace("List is full but add requestValueTransfer ", "newNonce", ev.Nonce(), "removedNonce", maxNonce)
		}

		bi.SetRequestNonceFromCounterpart(ev.GetRequestNonce() + 1)
		bi.pendingRequestEvent.Put(ev)
		vtPendingRequestEventCounter.Inc(1)
	}
	logger.Trace("added pending request events to the bridge info:", "bi.pendingRequestEvent", bi.pendingRequestEvent.Len())

	select {
	case bi.newEvent <- struct{}{}:
	default:
	}
}

// GetReadyRequestValueTransferEvents returns the processable events with the increasing nonce.
func (bi *BridgeInfo) GetReadyRequestValueTransferEvents() []IRequestValueTransferEvent {
	return bi.GetPendingRequestEvents()
}

// GetCurrentBlockNumber returns a current block number for each local and remote backend.
func (bi *BridgeInfo) GetCurrentBlockNumber() (uint64, error) {
	if bi.onChildChain {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return bi.subBridge.localBackend.CurrentBlockNumber(ctx)
	}
	return bi.subBridge.remoteBackend.CurrentBlockNumber(context.Background())
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
	b.ChildAddress, b.ParentAddress, b.Subscribed = elem.LocalAddress, elem.RemoteAddress, elem.Paired
	return nil
}

// EncodeRLP serializes a BridgeJournal into the Klaytn RLP BridgeJournal format.
func (b *BridgeJournal) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		b.ChildAddress,
		b.ParentAddress,
		b.Subscribed,
	})
}

// BridgeManager manages Bridge SmartContracts
// for value transfer between child chain and parent chain
type BridgeManager struct {
	subBridge *SubBridge

	receivedEvents map[common.Address][]event.Subscription
	withdrawEvents map[common.Address]event.Subscription
	bridges        map[common.Address]*BridgeInfo
	mu             sync.RWMutex

	reqVTevFeeder        event.Feed
	reqVTevEncodedFeeder event.Feed
	handleEventFeeder    event.Feed

	scope event.SubscriptionScope

	journal    *bridgeAddrJournal
	recoveries map[common.Address]*valueTransferRecovery
	auth       *bind.TransactOpts
}

func NewBridgeManager(main *SubBridge) (*BridgeManager, error) {
	bridgeAddrJournal := newBridgeAddrJournal(path.Join(main.config.DataDir, BridgeAddrJournal))

	bridgeManager := &BridgeManager{
		subBridge:      main,
		receivedEvents: make(map[common.Address][]event.Subscription),
		withdrawEvents: make(map[common.Address]event.Subscription),
		bridges:        make(map[common.Address]*BridgeInfo),
		journal:        bridgeAddrJournal,
		recoveries:     make(map[common.Address]*valueTransferRecovery),
	}

	logger.Info("Load Bridge Address from JournalFiles ", "path", bridgeManager.journal.path)
	bridgeManager.journal.cache = make(map[common.Address]*BridgeJournal)

	if err := bridgeManager.journal.load(func(gwjournal BridgeJournal) error {
		logger.Info("Load Bridge Address from JournalFiles ",
			"local address", gwjournal.ChildAddress.Hex(), "remote address", gwjournal.ParentAddress.Hex())
		bridgeManager.journal.cache[gwjournal.ChildAddress] = &gwjournal
		return nil
	}); err != nil {
		logger.Error("fail to load bridge address", "err", err)
	}

	if err := bridgeManager.journal.rotate(bridgeManager.GetAllBridge()); err != nil {
		logger.Error("fail to rotate bridge journal", "err", err)
	}

	return bridgeManager, nil
}

func (bm *BridgeManager) IsValidBridgePair(bridge1, bridge2 common.Address) bool {
	b1, ok1 := bm.GetBridgeInfo(bridge1)
	b2, ok2 := bm.GetBridgeInfo(bridge2)

	if ok1 && ok2 {
		if bridge1 == b2.counterpartAddress && bridge2 == b1.counterpartAddress {
			return true
		}
	}

	return false
}

func (bm *BridgeManager) GetCounterPartBridgeAddr(bridgeAddr common.Address) common.Address {
	bridge, ok := bm.GetBridgeInfo(bridgeAddr)
	if ok {
		return bridge.counterpartAddress
	}
	return common.Address{}
}

func (bm *BridgeManager) GetCounterPartBridge(bridgeAddr common.Address) *bridgecontract.Bridge {
	bridge, ok := bm.GetBridgeInfo(bridgeAddr)
	if ok {
		return bridge.counterpartBridge
	}
	return nil
}

// LogBridgeStatus logs the bridge contract requested/handled nonce status as an information.
func (bm *BridgeManager) LogBridgeStatus() {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if len(bm.bridges) == 0 {
		return
	}

	var p2cTotalRequestNonce, p2cTotalHandleNonce, p2cTotalLowerHandleNonce uint64
	var c2pTotalRequestNonce, c2pTotalHandleNonce, c2pTotalLowerHandleNonce uint64

	for bAddr, b := range bm.bridges {
		diffNonce := b.requestNonceFromCounterPart - b.handleNonce

		if b.subscribed {
			var headStr string
			if b.onChildChain {
				headStr = "Bridge(Parent -> Child Chain)"
				p2cTotalRequestNonce += b.requestNonceFromCounterPart
				p2cTotalHandleNonce += b.handleNonce
				p2cTotalLowerHandleNonce += b.lowerHandleNonce
			} else {
				headStr = "Bridge(Child -> Parent Chain)"
				c2pTotalRequestNonce += b.requestNonceFromCounterPart
				c2pTotalHandleNonce += b.handleNonce
				c2pTotalLowerHandleNonce += b.lowerHandleNonce
			}
			logger.Debug(headStr, "bridge", bAddr.String(), "requestNonce", b.requestNonceFromCounterPart, "lowerHandleNonce", b.lowerHandleNonce, "handleNonce", b.handleNonce, "pending", diffNonce)
		}
	}

	logger.Info("VT : Parent -> Child Chain", "request", p2cTotalRequestNonce, "handle", p2cTotalHandleNonce, "lowerHandle", p2cTotalLowerHandleNonce, "pending", p2cTotalRequestNonce-p2cTotalLowerHandleNonce)
	logger.Info("VT : Child -> Parent Chain", "request", c2pTotalRequestNonce, "handle", c2pTotalHandleNonce, "lowerHandle", c2pTotalLowerHandleNonce, "pending", c2pTotalRequestNonce-c2pTotalLowerHandleNonce)
}

// SubscribeReqVTev registers a subscription of RequestValueTransferEvent.
func (bm *BridgeManager) SubscribeReqVTev(ch chan<- RequestValueTransferEvent) event.Subscription {
	return bm.scope.Track(bm.reqVTevFeeder.Subscribe(ch))
}

// SubscribeReqVTencodedEv registers a subscription of RequestValueTransferEncoded.
func (bm *BridgeManager) SubscribeReqVTencodedEv(ch chan<- RequestValueTransferEncodedEvent) event.Subscription {
	return bm.scope.Track(bm.reqVTevEncodedFeeder.Subscribe(ch))
}

// SubscribeHandleVTev registers a subscription of HandleValueTransferEvent.
func (bm *BridgeManager) SubscribeHandleVTev(ch chan<- *HandleValueTransferEvent) event.Subscription {
	return bm.scope.Track(bm.handleEventFeeder.Subscribe(ch))
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
		return ErrNoBridgeInfo
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
		return ErrDuplicatedBridgeInfo
	}

	var counterpartBackend Backend
	if local {
		counterpartBackend = bm.subBridge.remoteBackend
	} else {
		counterpartBackend = bm.subBridge.localBackend
	}

	var err error
	bm.bridges[addr], err = NewBridgeInfo(bm.subBridge, addr, bridge, cpAddr, cpBridge, account, local, subscribed, counterpartBackend)
	return err
}

// RestoreBridges setups bridge subscription by using the journal cache.
func (bm *BridgeManager) RestoreBridges() error {
	if bm.subBridge.peers.Len() < 1 {
		logger.Debug("check peer connections to restore bridges")
		return ErrBridgeRestore
	}

	var counter = 0
	bm.stopAllRecoveries()

	for _, journal := range bm.journal.cache {
		cBridgeAddr := journal.ChildAddress
		pBridgeAddr := journal.ParentAddress
		bacc := bm.subBridge.bridgeAccounts

		// Set bridge info
		cBridgeInfo, cOk := bm.GetBridgeInfo(cBridgeAddr)
		pBridgeInfo, pOk := bm.GetBridgeInfo(pBridgeAddr)

		cBridge, err := bridgecontract.NewBridge(cBridgeAddr, bm.subBridge.localBackend)
		if err != nil {
			logger.Error("local bridge creation is failed", "err", err, "bridge", cBridge)
			break
		}

		pBridge, err := bridgecontract.NewBridge(pBridgeAddr, bm.subBridge.remoteBackend)
		if err != nil {
			logger.Error("remote bridge creation is failed", "err", err, "bridge", pBridge)
			break
		}

		if !cOk {
			err = bm.SetBridgeInfo(cBridgeAddr, cBridge, pBridgeAddr, pBridge, bacc.cAccount, true, false)
			if err != nil {
				logger.Error("setting local bridge info is failed", "err", err)
				bm.DeleteBridgeInfo(cBridgeAddr)
				break
			}
			cBridgeInfo, _ = bm.GetBridgeInfo(cBridgeAddr)
		}

		if !pOk {
			err = bm.SetBridgeInfo(pBridgeAddr, pBridge, cBridgeAddr, cBridgeInfo.bridge, bacc.pAccount, false, false)
			if err != nil {
				logger.Error("setting remote bridge info is failed", "err", err)
				bm.DeleteBridgeInfo(pBridgeAddr)
				break
			}
			pBridgeInfo, _ = bm.GetBridgeInfo(pBridgeAddr)
		}

		// Subscribe bridge events
		if journal.Subscribed {
			bm.UnsubscribeEvent(cBridgeAddr)
			bm.UnsubscribeEvent(pBridgeAddr)

			if !cBridgeInfo.subscribed {
				logger.Info("automatic local bridge subscription", "info", cBridgeInfo, "address", cBridgeInfo.address.String())
				if err := bm.subscribeEvent(cBridgeAddr, cBridgeInfo.bridge); err != nil {
					logger.Error("local bridge subscription is failed", "err", err)
					break
				}
			}
			if !pBridgeInfo.subscribed {
				logger.Info("automatic remote bridge subscription", "info", pBridgeInfo, "address", pBridgeInfo.address.String())
				if err := bm.subscribeEvent(pBridgeAddr, pBridgeInfo.bridge); err != nil {
					logger.Error("remote bridge subscription is failed", "err", err)
					bm.DeleteBridgeInfo(pBridgeAddr)
					break
				}
			}
			recovery := bm.recoveries[cBridgeAddr]
			if recovery == nil {
				bm.AddRecovery(cBridgeAddr, pBridgeAddr)
			}
		}

		counter++
	}

	if len(bm.journal.cache) == counter {
		logger.Info("succeeded to restore bridges", "pairs", counter)
		return nil
	}
	return ErrBridgeRestore
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
		return ErrNoBridgeInfo
	}
	remoteBridgeInfo, ok := bm.GetBridgeInfo(remoteAddress)
	if !ok {
		return ErrNoBridgeInfo
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
		return ErrNoRecovery
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

func (bm *BridgeManager) RegisterOperator(bridgeAddr, operatorAddr common.Address) (common.Hash, error) {
	bi, exist := bm.GetBridgeInfo(bridgeAddr)

	if !exist {
		return common.Hash{}, ErrNoBridgeInfo
	}

	bi.account.Lock()
	defer bi.account.UnLock()
	tx, err := bi.bridge.RegisterOperator(bi.account.GenerateTransactOpts(), operatorAddr)
	if err != nil {
		return common.Hash{}, err
	}
	bi.account.IncNonce()

	return tx.Hash(), nil
}

func (bm *BridgeManager) GetOperators(bridgeAddr common.Address) ([]common.Address, error) {
	bi, exist := bm.GetBridgeInfo(bridgeAddr)

	if !exist {
		return nil, ErrNoBridgeInfo
	}

	return bi.bridge.GetOperatorList(nil)
}

func (bm *BridgeManager) SetValueTransferOperatorThreshold(bridgeAddr common.Address, threshold uint8) (common.Hash, error) {
	bi, exist := bm.GetBridgeInfo(bridgeAddr)

	if !exist {
		return common.Hash{}, ErrNoBridgeInfo
	}

	bi.account.Lock()
	defer bi.account.UnLock()
	tx, err := bi.bridge.SetOperatorThreshold(bi.account.GenerateTransactOpts(), voteTypeValueTransfer, threshold)
	if err != nil {
		return common.Hash{}, err
	}
	bi.account.IncNonce()

	return tx.Hash(), nil
}

func (bm *BridgeManager) GetValueTransferOperatorThreshold(bridgeAddr common.Address) (uint8, error) {
	bi, exist := bm.GetBridgeInfo(bridgeAddr)

	if !exist {
		return 0, ErrNoBridgeInfo
	}

	threshold, err := bi.bridge.OperatorThresholds(nil, voteTypeValueTransfer)
	if err != nil {
		return 0, err
	}

	return threshold, nil
}

// Deploy Bridge SmartContract on same node or remote node
func (bm *BridgeManager) DeployBridge(auth *bind.TransactOpts, backend bind.ContractBackend, local bool) (*bridgecontract.Bridge, common.Address, error) {
	var acc *accountInfo
	var modeMintBurn bool

	if local {
		acc = bm.subBridge.bridgeAccounts.cAccount
		modeMintBurn = true
	} else {
		acc = bm.subBridge.bridgeAccounts.pAccount
		modeMintBurn = false
	}

	addr, bridge, err := bm.deployBridge(acc, auth, backend, modeMintBurn)
	if err != nil {
		return nil, common.Address{}, err
	}

	return bridge, addr, err
}

// DeployBridge handles actual smart contract deployment.
// To create contract, the chain ID, nonce, account key, private key, contract binding and gas price are used.
// The deployed contract address, transaction are returned. An error is also returned if any.
func (bm *BridgeManager) deployBridge(acc *accountInfo, auth *bind.TransactOpts, backend bind.ContractBackend, modeMintBurn bool) (common.Address, *bridgecontract.Bridge, error) {
	acc.Lock()
	addr, tx, contract, err := bridgecontract.DeployBridge(auth, backend, modeMintBurn)
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
		return ErrNoBridgeInfo
	}
	if bridgeInfo.subscribed {
		return ErrAlreadySubscribed
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
			bm.UnsubscribeEvent(journal.ChildAddress)
			bm.UnsubscribeEvent(journal.ParentAddress)

			bridgeInfo, ok := bm.GetBridgeInfo(journal.ChildAddress)
			if !ok {
				logger.Error("ResetAllSubscribedEvents failed to GetBridgeInfo", "localBridge", journal.ChildAddress.String())
				return ErrNoBridgeInfo
			}
			err := bm.subscribeEvent(journal.ChildAddress, bridgeInfo.bridge)
			if err != nil {
				return err
			}

			bridgeInfo, ok = bm.GetBridgeInfo(journal.ParentAddress)
			if !ok {
				logger.Error("ResetAllSubscribedEvents failed to GetBridgeInfo", "remoteBridge", journal.ParentAddress.String())
				bm.UnsubscribeEvent(journal.ChildAddress)
				return ErrNoBridgeInfo
			}
			err = bm.subscribeEvent(journal.ParentAddress, bridgeInfo.bridge)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SubscribeEvent sets watch logs and creates a goroutine loop to handle event messages.
func (bm *BridgeManager) subscribeEvent(addr common.Address, bridge *bridgecontract.Bridge) error {
	chanReqVT := make(chan *bridgecontract.BridgeRequestValueTransfer, TokenEventChanSize)
	chanReqVTencoded := make(chan *bridgecontract.BridgeRequestValueTransferEncoded, TokenEventChanSize)
	chanHandleVT := make(chan *bridgecontract.BridgeHandleValueTransfer, TokenEventChanSize)

	vtEv, err := bridge.WatchRequestValueTransfer(nil, chanReqVT, nil, nil, nil)
	if err != nil {
		logger.Error("Failed to watch RequestValueTransfer event", "err", err)
		return err
	}
	bm.receivedEvents[addr] = append(bm.receivedEvents[addr], vtEv)

	vtEncodedev, err := bridge.WatchRequestValueTransferEncoded(nil, chanReqVTencoded, nil, nil, nil)
	if err != nil {
		logger.Error("Failed to watch RequestValueTransferEncoded event", "err", err)
		return err
	}
	bm.receivedEvents[addr] = append(bm.receivedEvents[addr], vtEncodedev)

	withdrawnSub, err := bridge.WatchHandleValueTransfer(nil, chanHandleVT, nil, nil, nil)
	if err != nil {
		logger.Error("Failed to watch HandleValueTransfer event", "err", err)
		vtEv.Unsubscribe()
		vtEncodedev.Unsubscribe()
		delete(bm.receivedEvents, addr)
		return err
	}
	bm.withdrawEvents[addr] = withdrawnSub
	bridgeInfo, ok := bm.GetBridgeInfo(addr)
	if !ok {
		vtEv.Unsubscribe()
		vtEncodedev.Unsubscribe()
		withdrawnSub.Unsubscribe()
		delete(bm.receivedEvents, addr)
		delete(bm.withdrawEvents, addr)
		return ErrNoBridgeInfo
	}
	bridgeInfo.subscribed = true

	go bm.loop(addr, chanReqVT, chanReqVTencoded, chanHandleVT, vtEv, vtEncodedev, withdrawnSub)

	return nil
}

// UnsubscribeEvent cancels the contract's watch logs and initializes the status.
func (bm *BridgeManager) UnsubscribeEvent(addr common.Address) {
	receivedSub := bm.receivedEvents[addr]
	for _, sub := range receivedSub {
		sub.Unsubscribe()
	}
	delete(bm.receivedEvents, addr)

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
	chanReqVT <-chan *bridgecontract.BridgeRequestValueTransfer,
	chanReqVTencoded <-chan *bridgecontract.BridgeRequestValueTransferEncoded,
	chanHandleVT <-chan *bridgecontract.BridgeHandleValueTransfer,
	reqVTevSub, reqVTencodedEvSub event.Subscription,
	handleEventSub event.Subscription) {

	defer reqVTevSub.Unsubscribe()
	defer reqVTencodedEvSub.Unsubscribe()
	defer handleEventSub.Unsubscribe()

	bi, ok := bm.GetBridgeInfo(addr)
	if !ok {
		logger.Error("bridge information is missing")
		return
	}

	// TODO-Klaytn change goroutine logic for performance
	for {
		select {
		case <-bi.closed:
			return
		case ev := <-chanReqVT:
			bm.reqVTevFeeder.Send(RequestValueTransferEvent{ev})
		case ev := <-chanReqVTencoded:
			bm.reqVTevEncodedFeeder.Send(RequestValueTransferEncodedEvent{ev})
		case ev := <-chanHandleVT:
			bm.handleEventFeeder.Send(&HandleValueTransferEvent{ev})
		case err := <-reqVTevSub.Err():
			logger.Info("Contract Event Loop Running Stop by receivedSub.Err()", "err", err)
			return
		case err := <-reqVTencodedEvSub.Err():
			logger.Info("Contract Event Loop Running Stop by receivedSub.Err()", "err", err)
			return
		case err := <-handleEventSub.Err():
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

// SetERC20Fee set the ERC20 transfer fee on the bridge contract.
func (bm *BridgeManager) SetERC20Fee(bridgeAddr, tokenAddr common.Address, fee *big.Int) (common.Hash, error) {
	bi, ok := bm.GetBridgeInfo(bridgeAddr)
	if !ok {
		return common.Hash{}, ErrNoBridgeInfo
	}

	auth := bi.account
	auth.Lock()
	defer auth.UnLock()

	rn, err := bi.bridge.ConfigurationNonce(nil)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := bi.bridge.SetERC20Fee(auth.GenerateTransactOpts(), tokenAddr, fee, rn)
	if err != nil {
		return common.Hash{}, err
	}

	auth.IncNonce()

	return tx.Hash(), nil
}

// SetKLAYFee set the KLAY transfer fee on the bridge contract.
func (bm *BridgeManager) SetKLAYFee(bridgeAddr common.Address, fee *big.Int) (common.Hash, error) {
	bi, ok := bm.GetBridgeInfo(bridgeAddr)
	if !ok {
		return common.Hash{}, ErrNoBridgeInfo
	}

	auth := bi.account
	auth.Lock()
	defer auth.UnLock()

	rn, err := bi.bridge.ConfigurationNonce(nil)
	if err != nil {
		return common.Hash{}, err
	}

	tx, err := bi.bridge.SetKLAYFee(auth.GenerateTransactOpts(), fee, rn)
	if err != nil {
		return common.Hash{}, err
	}

	auth.IncNonce()

	return tx.Hash(), nil
}

// SetFeeReceiver set the fee receiver which can get the fee of value transfer request.
func (bm *BridgeManager) SetFeeReceiver(bridgeAddr, receiver common.Address) (common.Hash, error) {
	bi, ok := bm.GetBridgeInfo(bridgeAddr)
	if !ok {
		return common.Hash{}, ErrNoBridgeInfo
	}

	auth := bi.account
	auth.Lock()
	defer auth.UnLock()

	tx, err := bi.bridge.SetFeeReceiver(auth.GenerateTransactOpts(), receiver)
	if err != nil {
		return common.Hash{}, err
	}

	auth.IncNonce()

	return tx.Hash(), nil
}

// GetERC20Fee returns the ERC20 transfer fee on the bridge contract.
func (bm *BridgeManager) GetERC20Fee(bridgeAddr, tokenAddr common.Address) (*big.Int, error) {
	bi, ok := bm.GetBridgeInfo(bridgeAddr)
	if !ok {
		return nil, ErrNoBridgeInfo
	}

	return bi.bridge.FeeOfERC20(nil, tokenAddr)
}

// GetKLAYFee returns the KLAY transfer fee on the bridge contract.
func (bm *BridgeManager) GetKLAYFee(bridgeAddr common.Address) (*big.Int, error) {
	bi, ok := bm.GetBridgeInfo(bridgeAddr)
	if !ok {
		return nil, ErrNoBridgeInfo
	}

	return bi.bridge.FeeOfKLAY(nil)
}

// GetFeeReceiver returns the receiver which can get fee of value transfer request.
func (bm *BridgeManager) GetFeeReceiver(bridgeAddr common.Address) (common.Address, error) {
	bi, ok := bm.GetBridgeInfo(bridgeAddr)
	if !ok {
		return common.Address{}, ErrNoBridgeInfo
	}

	return bi.bridge.FeeReceiver(nil)
}

// IsInParentAddrs returns true if the bridgeAddr is in the list of parent bridge addresses and returns false if not.
func (bm *BridgeManager) IsInParentAddrs(bridgeAddr common.Address) bool {
	for _, journal := range bm.journal.cache {
		if journal.ParentAddress == bridgeAddr {
			return true
		}
	}
	return false
}

// IsInChildAddrs returns true if the bridgeAddr is in the list of child bridge addresses and returns false if not.
func (bm *BridgeManager) IsInChildAddrs(bridgeAddr common.Address) bool {
	for _, journal := range bm.journal.cache {
		if journal.ChildAddress == bridgeAddr {
			return true
		}
	}
	return false
}
