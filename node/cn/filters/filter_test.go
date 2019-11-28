package filters

import (
	"github.com/golang/mock/gomock"
	cn "github.com/klaytn/klaytn/node/cn/filters/mock"
	"github.com/stretchr/testify/assert"

	"github.com/klaytn/klaytn/common"
	"testing"
)

func TestFilter_New(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockBackend := cn.NewMockBackend(mockCtrl)
	defer mockCtrl.Finish()

	mockBackend.EXPECT().BloomStatus().Return(uint64(123), uint64(321)).Times(1)

	addr1 := common.HexToAddress("111")
	addr2 := common.HexToAddress("222")
	addrs := []common.Address{addr1, addr2}

	topic1 := addr1.Hash()
	topic2 := addr2.Hash()
	topics := [][]common.Hash{{topic1}, {topic2}}

	begin := int64(12345)
	end := int64(54321)

	newFilter := New(mockBackend, begin, end, addrs, topics)
	assert.NotNil(t, newFilter)
	assert.Equal(t, begin, newFilter.begin)
	assert.Equal(t, end, newFilter.end)
	assert.Equal(t, topics, newFilter.topics)
	assert.Equal(t, addrs, newFilter.addresses)
}

func TestFilter_Logs(t *testing.T) {

}
