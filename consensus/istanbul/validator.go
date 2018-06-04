package istanbul

import "ground-x/go-gxplatform/common"

type Validator interface {

	Address() common.Address

	String() string
}

type ValidatorSet interface {

}
