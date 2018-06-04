package istanbul

import (
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/event"
)

type Backend interface {

	Address() common.Address

	Validator(proposal Proposal) ValidatorSet

	EventMux() *event.TypeMux

	Broadcast(valSet ValidatorSet, payload []byte) error
}