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

package reward

import (
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func newTestBlockChain() *blockchain.BlockChain {
	return &blockchain.BlockChain{}
}

func TestAddressBookManager_makeMsgToAddressBook(t *testing.T) {
	targetAddress := "0x0000000000000000000000000000000000000400" // address of addressBook which the message has to be sent to
	ac := newAddressBookConnector(newTestBlockChain(), nil)
	msg, err := ac.makeMsgToAddressBook(params.Rules{IsIstanbul: true})
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	assert.Equal(t, targetAddress, msg.To().String())
}
