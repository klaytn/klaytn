// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
)

var _ = (*headerMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (h Header) MarshalJSON() ([]byte, error) {
	type Header struct {
		ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
		Rewardbase  common.Address `json:"reward"           gencodec:"required"`
		Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
		TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
		Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
		BlockScore  *hexutil.Big   `json:"blockScore"       gencodec:"required"`
		Number      *hexutil.Big   `json:"number"           gencodec:"required"`
		GasUsed     hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
		Time        *hexutil.Big   `json:"timestamp"        gencodec:"required"`
		TimeFoS     hexutil.Uint   `json:"timestampFoS"              gencodec:"required"`
		Extra       hexutil.Bytes  `json:"extraData"                 gencodec:"required"`
		Governance  hexutil.Bytes  `json:"governanceData"            gencodec:"required"`
		Vote        hexutil.Bytes  `json:"voteData,omitempty"`
		BaseFee     *hexutil.Big   `json:"baseFeePerGas,omitempty"    rlp:"optional"`
		Hash        common.Hash    `json:"hash"`
	}
	var enc Header
	enc.ParentHash = h.ParentHash
	enc.Rewardbase = h.Rewardbase
	enc.Root = h.Root
	enc.TxHash = h.TxHash
	enc.ReceiptHash = h.ReceiptHash
	enc.Bloom = h.Bloom
	enc.BlockScore = (*hexutil.Big)(h.BlockScore)
	enc.Number = (*hexutil.Big)(h.Number)
	enc.GasUsed = hexutil.Uint64(h.GasUsed)
	enc.Time = (*hexutil.Big)(h.Time)
	enc.TimeFoS = hexutil.Uint(h.TimeFoS)
	enc.Extra = h.Extra
	enc.Governance = h.Governance
	enc.Vote = h.Vote
	enc.BaseFee = (*hexutil.Big)(h.BaseFee)
	enc.Hash = h.Hash()
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (h *Header) UnmarshalJSON(input []byte) error {
	type Header struct {
		ParentHash  *common.Hash    `json:"parentHash"       gencodec:"required"`
		Rewardbase  *common.Address `json:"reward"           gencodec:"required"`
		Root        *common.Hash    `json:"stateRoot"        gencodec:"required"`
		TxHash      *common.Hash    `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash *common.Hash    `json:"receiptsRoot"     gencodec:"required"`
		Bloom       *Bloom          `json:"logsBloom"        gencodec:"required"`
		BlockScore  *hexutil.Big    `json:"blockScore"       gencodec:"required"`
		Number      *hexutil.Big    `json:"number"           gencodec:"required"`
		GasUsed     *hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
		Time        *hexutil.Big    `json:"timestamp"        gencodec:"required"`
		TimeFoS     *hexutil.Uint   `json:"timestampFoS"              gencodec:"required"`
		Extra       *hexutil.Bytes  `json:"extraData"                 gencodec:"required"`
		Governance  *hexutil.Bytes  `json:"governanceData"            gencodec:"required"`
		Vote        *hexutil.Bytes  `json:"voteData,omitempty"`
		BaseFee     *hexutil.Big    `json:"baseFeePerGas,omitempty"    rlp:"optional"`
	}
	var dec Header
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for Header")
	}
	h.ParentHash = *dec.ParentHash
	if dec.Rewardbase == nil {
		return errors.New("missing required field 'reward' for Header")
	}
	h.Rewardbase = *dec.Rewardbase
	if dec.Root == nil {
		return errors.New("missing required field 'stateRoot' for Header")
	}
	h.Root = *dec.Root
	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionsRoot' for Header")
	}
	h.TxHash = *dec.TxHash
	if dec.ReceiptHash == nil {
		return errors.New("missing required field 'receiptsRoot' for Header")
	}
	h.ReceiptHash = *dec.ReceiptHash
	if dec.Bloom == nil {
		return errors.New("missing required field 'logsBloom' for Header")
	}
	h.Bloom = *dec.Bloom
	if dec.BlockScore == nil {
		return errors.New("missing required field 'blockScore' for Header")
	}
	h.BlockScore = (*big.Int)(dec.BlockScore)
	if dec.Number == nil {
		return errors.New("missing required field 'number' for Header")
	}
	h.Number = (*big.Int)(dec.Number)
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for Header")
	}
	h.GasUsed = uint64(*dec.GasUsed)
	if dec.Time == nil {
		return errors.New("missing required field 'timestamp' for Header")
	}
	h.Time = (*big.Int)(dec.Time)
	if dec.TimeFoS == nil {
		return errors.New("missing required field 'timestampFoS' for Header")
	}
	h.TimeFoS = uint8(*dec.TimeFoS)
	if dec.Extra == nil {
		return errors.New("missing required field 'extraData' for Header")
	}
	h.Extra = *dec.Extra
	if dec.Governance == nil {
		return errors.New("missing required field 'governanceData' for Header")
	}
	h.Governance = *dec.Governance
	if dec.Vote != nil {
		h.Vote = *dec.Vote
	}
	if dec.BaseFee != nil {
		h.BaseFee = (*big.Int)(dec.BaseFee)
	}
	return nil
}
