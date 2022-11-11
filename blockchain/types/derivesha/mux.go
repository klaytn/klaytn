package derivesha

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
)

// TODO-Klaytn: Make DeriveShaMux state-less
// As-is: InitDS(type) + DS(list) + ERH
// To-be: InitDS() + DS(list, type) + ERH(type)

type IDeriveSha interface {
	DeriveSha(list types.DerivableList) common.Hash
}

var deriveShaObj IDeriveSha = nil
var logger = log.NewModuleLogger(log.Blockchain)

func InitDeriveSha(implType int) {
	switch implType {
	case types.ImplDeriveShaOriginal:
		deriveShaObj = DeriveShaOrig{}
	case types.ImplDeriveShaSimple:
		deriveShaObj = DeriveShaSimple{}
	case types.ImplDeriveShaConcat:
		deriveShaObj = DeriveShaConcat{}
	default:
		logger.Error("Unrecognized deriveShaImpl, falling back to Orig", "impl", implType)
		deriveShaObj = DeriveShaOrig{}
	}

	types.DeriveSha = DeriveShaMux
	types.EmptyRootHash = DeriveShaMux(types.Transactions{})
}

func DeriveShaMux(list types.DerivableList) common.Hash {
	return deriveShaObj.DeriveSha(list)
}
