package reward

import (
	"github.com/klaytn/klaytn/blockchain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestBlockChain() *blockchain.BlockChain {
	return &blockchain.BlockChain{}
}

func TestAddressBookManager_makeMsgToAddressBook(t *testing.T) {
	targetAddress := "0x0000000000000000000000000000000000000400" // address of addressBook which the message has to be sent to
	addressBookManager := newAddressBookManager(newTestBlockChain(), nil)
	msg, error := addressBookManager.makeMsgToAddressBook()
	if error != nil {
		t.Errorf("error has occurred. error : %v", error)
	} else {
		assert.Equal(t, msg.To().String(), targetAddress)
	}
}
