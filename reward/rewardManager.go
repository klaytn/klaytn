package reward

import "github.com/klaytn/klaytn/log"

var logger = log.NewModuleLogger(log.Reward)

type governanceHelper interface {
	Epoch() uint64
	GetItemAtNumberByIntKey(num uint64, key int) (interface{}, error)
	DeferredTxFee() bool
}
