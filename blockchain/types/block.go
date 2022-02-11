// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/types/block.go (2018/06/04).
// Modified and improved for the klaytn development.

package types

import (
	"fmt"
	"io"
	"math/big"
	"sort"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto/sha3"

	"github.com/klaytn/klaytn/rlp"
)

const (
	Engine_IBFT int = iota
	Engine_Clique
	Engine_Gxhash
)

var (
	// EmptyRootHash is transaction/receipt root hash when there is no transaction.
	// This value is initialized in InitDeriveSha.
	EmptyRootHash = common.Hash{}
	EngineType    = Engine_IBFT
)

//go:generate gencodec -type Header -field-override headerMarshaling -out gen_header_json.go

// Header represents a block header in the Klaytn blockchain.
type Header struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	Rewardbase  common.Address `json:"reward"           gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	BlockScore  *big.Int       `json:"blockScore"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        *big.Int       `json:"timestamp"        gencodec:"required"`
	// TimeFoS represents a fraction of a second since `Time`.
	TimeFoS    uint8  `json:"timestampFoS"     gencodec:"required"`
	Extra      []byte `json:"extraData"        gencodec:"required"`
	Governance []byte `json:"governanceData"        gencodec:"required"`
	Vote       []byte `json:"voteData,omitempty"`
}

// field type overrides for gencodec
type headerMarshaling struct {
	BlockScore *hexutil.Big
	Number     *hexutil.Big
	GasUsed    hexutil.Uint64
	Time       *hexutil.Big
	TimeFoS    hexutil.Uint
	Extra      hexutil.Bytes
	Hash       common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
	Governance hexutil.Bytes
	Vote       hexutil.Bytes
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	// If the mix digest is equivalent to the predefined Istanbul digest, use Istanbul
	// specific hash calculation.
	if EngineType == Engine_IBFT {
		// Seal is reserved in extra-data. To prove block is signed by the proposer.
		if istanbulHeader := IstanbulFilteredHeader(h, true); istanbulHeader != nil {
			return rlpHash(istanbulHeader)
		}
	}
	return rlpHash(h)
}

// HashNoNonce returns the hash which is used as input for the proof-of-work search.
func (h *Header) HashNoNonce() common.Hash {
	return rlpHash([]interface{}{
		h.ParentHash,
		h.Rewardbase,
		h.Root,
		h.TxHash,
		h.ReceiptHash,
		h.Bloom,
		h.BlockScore,
		h.Number,
		h.GasUsed,
		h.Time,
		h.TimeFoS,
		h.Extra,
	})
}

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header) Size() common.StorageSize {
	return common.StorageSize(unsafe.Sizeof(*h)) + common.StorageSize(len(h.Extra)+(h.BlockScore.BitLen()+h.Number.BitLen()+h.Time.BitLen())/8)
}

func (h *Header) Round() byte {
	return byte(h.Extra[IstanbulExtraVanity-1])
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// prefixedRlpHash writes the prefix into the hasher before rlp-encoding the
// given interface. It's used for ethereum typed transactions.
func prefixedRlpHash(prefix byte, x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	hw.Reset()
	hw.Write([]byte{prefix})
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// EmptyBody returns true if there is no additional 'body' to complete the header
// that is: no transactions.
func (h *Header) EmptyBody() bool {
	return h.TxHash == EmptyRootHash
}

// EmptyReceipts returns true if there are no receipts for this header/block.
func (h *Header) EmptyReceipts() bool {
	return h.ReceiptHash == EmptyRootHash
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions) together.
type Body struct {
	Transactions []*Transaction
}

// Block represents an entire block in the Klaytn blockchain.
type Block struct {
	header       *Header
	transactions Transactions

	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package blockchain to store the total blockscore
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used to track inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

type Result struct {
	Block *Block
	Round int64
}

// extblock represents external block encoding used for Klaytn protocol, etc.
type extblock struct {
	Header *Header
	Txs    []*Transaction
}

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs and receipts.
func NewBlock(header *Header, txs []*Transaction, receipts []*Receipt) *Block {
	b := &Block{header: CopyHeader(header), td: new(big.Int)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.header.TxHash = EmptyRootHash
	} else {
		b.header.TxHash = DeriveSha(Transactions(txs))
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptHash = EmptyRootHash
	} else {
		b.header.ReceiptHash = DeriveSha(Receipts(receipts))
		b.header.Bloom = CreateBloom(receipts)
	}

	return b
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.Time = new(big.Int); h.Time != nil {
		cpy.Time.Set(h.Time)
	}
	if cpy.BlockScore = new(big.Int); h.BlockScore != nil {
		cpy.BlockScore.Set(h.BlockScore)
	}
	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	if len(h.Governance) > 0 {
		cpy.Governance = make([]byte, len(h.Governance))
		copy(cpy.Governance, h.Governance)
	}
	if len(h.Vote) > 0 {
		cpy.Vote = make([]byte, len(h.Vote))
		copy(cpy.Vote, h.Vote)
	}
	return &cpy
}

// DecodeRLP decodes the Klaytn
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.transactions = eb.Header, eb.Txs
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes a block into the Klaytn RLP block format.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock{
		Header: b.header,
		Txs:    b.transactions,
	})
}

func (b *Block) Transactions() Transactions { return b.transactions }

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}

func (b *Block) Number() *big.Int     { return new(big.Int).Set(b.header.Number) }
func (b *Block) GasUsed() uint64      { return b.header.GasUsed }
func (b *Block) BlockScore() *big.Int { return new(big.Int).Set(b.header.BlockScore) }
func (b *Block) Time() *big.Int       { return new(big.Int).Set(b.header.Time) }
func (b *Block) TimeFoS() uint8       { return b.header.TimeFoS }

func (b *Block) NumberU64() uint64          { return b.header.Number.Uint64() }
func (b *Block) Bloom() Bloom               { return b.header.Bloom }
func (b *Block) Rewardbase() common.Address { return b.header.Rewardbase }
func (b *Block) Root() common.Hash          { return b.header.Root }
func (b *Block) ParentHash() common.Hash    { return b.header.ParentHash }
func (b *Block) TxHash() common.Hash        { return b.header.TxHash }
func (b *Block) ReceiptHash() common.Hash   { return b.header.ReceiptHash }
func (b *Block) Extra() []byte              { return common.CopyBytes(b.header.Extra) }

func (b *Block) Header() *Header { return CopyHeader(b.header) }

// Body returns the non-header content of the block.
func (b *Block) Body() *Body { return &Body{b.transactions} }

func (b *Block) HashNoNonce() common.Hash {
	return b.header.HashNoNonce()
}

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	cpy := *header

	return &Block{
		header:       &cpy,
		transactions: b.transactions,
	}
}

// WithBody returns a new block with the given transactions.
func (b *Block) WithBody(transactions []*Transaction) *Block {
	block := &Block{
		header:       CopyHeader(b.header),
		transactions: make([]*Transaction, len(transactions)),
	}
	copy(block.transactions, transactions)
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

func (b *Block) String() string {
	str := fmt.Sprintf(`Block(#%v): Size: %v {
MinerHash: %x
%v
Transactions:
%v
}
`, b.Number(), b.Size(), b.header.HashNoNonce(), b.header, b.transactions)
	return str
}

func (h *Header) String() string {
	return fmt.Sprintf(`Header(%x):
[
	ParentHash:       %x
	Rewardbase:       %x
	Root:             %x
	TxSha:            %x
	ReceiptSha:       %x
	Bloom:            %x
	BlockScore:       %v
	Number:           %v
	GasUsed:          %v
	Time:             %v
	TimeFoS:          %v
	Extra:            %s
	Governance:       %x
	Vote:             %x
]`, h.Hash(), h.ParentHash, h.Rewardbase, h.Root, h.TxHash, h.ReceiptHash, h.Bloom, h.BlockScore, h.Number, h.GasUsed, h.Time, h.TimeFoS, h.Extra, h.Governance, h.Vote)
}

type Blocks []*Block

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks Blocks) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self blockSorter) Len() int { return len(self.blocks) }
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

func Number(b1, b2 *Block) bool { return b1.header.Number.Cmp(b2.header.Number) < 0 }
