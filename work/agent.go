// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from miner/agent.go (2018/06/04).
// Modified and improved for the klaytn development.

package work

import (
	"sync"
	"sync/atomic"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
)

type CpuAgent struct {
	mu sync.Mutex

	workCh        chan *Task
	stop          chan struct{}
	quitCurrentOp chan struct{}
	returnCh      chan<- *Result

	chain  consensus.ChainReader
	engine consensus.Engine

	isMining int32 // isMining indicates whether the agent is currently mining

	nodetype common.ConnType
}

func NewCpuAgent(chain consensus.ChainReader, engine consensus.Engine, nodetype common.ConnType) *CpuAgent {
	miner := &CpuAgent{
		chain:    chain,
		engine:   engine,
		stop:     make(chan struct{}, 1),
		workCh:   make(chan *Task, 1),
		nodetype: nodetype,
	}
	return miner
}

func (self *CpuAgent) Work() chan<- *Task            { return self.workCh }
func (self *CpuAgent) SetReturnCh(ch chan<- *Result) { self.returnCh = ch }

func (self *CpuAgent) Stop() {
	if !atomic.CompareAndSwapInt32(&self.isMining, 1, 0) {
		return // agent already stopped
	}
	self.stop <- struct{}{}
done:
	// Empty work channel
	for {
		select {
		case <-self.workCh:
		default:
			break done
		}
	}
}

func (self *CpuAgent) Start() {
	if !atomic.CompareAndSwapInt32(&self.isMining, 0, 1) {
		return // agent already started
	}
	go self.update()
}

func (self *CpuAgent) update() {
out:
	for {
		select {
		case work := <-self.workCh:
			self.mu.Lock()
			if self.quitCurrentOp != nil {
				close(self.quitCurrentOp)
			}
			self.quitCurrentOp = make(chan struct{})
			go self.mine(work, self.quitCurrentOp)
			self.mu.Unlock()
		case <-self.stop:
			self.mu.Lock()
			if self.quitCurrentOp != nil {
				close(self.quitCurrentOp)
				self.quitCurrentOp = nil
			}
			self.mu.Unlock()
			break out
		}
	}
}

func (self *CpuAgent) mine(work *Task, stop <-chan struct{}) {
	// TODO-Klaytn drop or missing tx and remove mining on PN and EN
	if self.nodetype != common.CONSENSUSNODE {
		ResultChGauge.Update(ResultChGauge.Value() + 1)
		self.returnCh <- &Result{work, nil}
		return
	}

	if result, err := self.engine.Seal(self.chain, work.Block, stop); result != nil {
		logger.Info("Successfully sealed new block", "number", result.Number(), "hash", result.Hash())

		ResultChGauge.Update(ResultChGauge.Value() + 1)
		self.returnCh <- &Result{work, result}
	} else {
		if err != nil {
			logger.Warn("Block sealing failed", "err", err)
		}
		ResultChGauge.Update(ResultChGauge.Value() + 1)
		self.returnCh <- nil
	}
}

func (self *CpuAgent) GetHashRate() int64 {
	if pow, ok := self.engine.(consensus.PoW); ok {
		return int64(pow.Hashrate())
	}
	return 0
}
