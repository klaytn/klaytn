package validator

import (
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/consensus/istanbul"
)

func New(addr common.Address) istanbul.Validator {
	return &defaultValidator{
		address: addr,
	}
}

func NewSet(addrs []common.Address, policy istanbul.ProposerPolicy) istanbul.ValidatorSet {
	return newDefaultSet(addrs, policy)
}
