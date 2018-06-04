package istanbul

import (
	"math/big"
	"ground-x/go-gxplatform/common"
	"io"
	"ground-x/go-gxplatform/rlp"
)

type Proposal interface {

	Number() *big.Int

	Hash() common.Hash

	EncodeRLP(w io.Writer) error

	DecodeRLP(s *rlp.Stream) error

	String() string
}

type Request struct {
	Proposal Proposal
}

type View struct {
	Round 	 *big.Int
	Sequence *big.Int
}

func (v *View) EncodeRLP(w io.Writer) error {

	return nil
}

func (v *View) DecodeRLP(s *rlp.Stream) error {

	return nil
}

func (v *View) String() string {
	return ""
}

// Cmp compares v and y and returns:
// -1 if v < y
//  0 if v == y
// +1 if v > y
func (v *View) Cmp(y *View) int {
	sdiff := v.Sequence.Cmp(y.Sequence)
	if sdiff != 0 {
		return sdiff
	}
	rdiff := v.Round.Cmp(y.Round)
	if rdiff != 0 {
		return rdiff
	}
	return 0
}

type Preprepare struct {
	View 	 *View
	Proposal Proposal
}

func (b *Preprepare) EncodeRLP(w io.Writer) error {
	return nil
}

func (b *Preprepare) DecodeRLP(s *rlp.Stream) error {
	return nil
}

func (b *Preprepare) String() string {
	return ""
}

type Subject struct {
	View *View
	Digest common.Hash
}

func (b *Subject) EncodeRLP(w io.Writer) error {
	return nil
}

func (b *Subject) DecodeRLP(s *rlp.Stream) error {
	return nil
}

func (b *Subject) String() string {
	return ""
}



