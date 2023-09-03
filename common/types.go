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
// This file is derived from common/types.go (2018/06/04).
// Modified and improved for the klaytn development.

package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto/sha3"
)

const (
	HashLength           = 32
	ExtHashCounterLength = 7
	ExtHashLength        = HashLength + ExtHashCounterLength
	AddressLength        = 20
	SignatureLength      = 65
)

var (
	hashT    = reflect.TypeOf(Hash{})
	extHashT = reflect.TypeOf(ExtHash{})
	addressT = reflect.TypeOf(Address{})
)

var (
	lastPrecompiledContractAddressHex = hexutil.MustDecode("0x00000000000000000000000000000000000003FF")

	// extHashLastCounter is the counter used to generate the counter for the ExtHash.
	// It starts off the most significant 7 bytes of the timestamp in nanoseconds at program startup.
	// It increments every time a new ExtHash counter is generated.
	//                      [b1 b2 b3 b4 b5 b6 b7 b8] = UnixNano()
	// extHashLastCounter = [00 b1 b2 b3 b4 b5 b6 b7]
	// nextCounter        =    [b1 b2 b3 b4 b5 b6 b7]
	extHashLastCounter = uint64(0)
	// extHashZeroCounter signifies the trie node referred by the ExtHash is actually
	// identified by the regular 32-byte Hash.
	extHashZeroCounter = ExtHashCounter{0, 0, 0, 0, 0, 0, 0}
)

func init() {
	extHashLastCounter = uint64(time.Now().UnixNano() >> 8)
	if extHashLastCounter == 0 {
		panic("Failed to retrieve current timestamp for ExtHashCounter")
	}
}

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x…%x", h[:3], h[29:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (h Hash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), h[:])
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

// getShardIndex returns the index of the shard.
// The address is arranged in the front or back of the array according to the initialization method.
// And the opposite is zero. In any case, to calculate the various shard index values,
// add both values and shift to calculate the shard index.
func (h Hash) getShardIndex(shardMask int) int {
	data1 := int(h[HashLength-1]) + int(h[0])
	data2 := int(h[HashLength-2]) + int(h[1])
	return ((data2 << 8) + data1) & shardMask
}

func EmptyHash(h Hash) bool {
	return h == Hash{}
}

// UnprefixedHash allows marshaling a Hash without 0x prefix.
type UnprefixedHash Hash

// UnmarshalText decodes the hash from hex. The 0x prefix is optional.
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// MarshalText encodes the hash as hex.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

/////////// ExtHash

type (
	// ExtHash is an extended hash composed of a 32 byte Hash and a 7 byte Counter.
	// ExtHash is used as the reference of Merkle Patricia Trie nodes to enable
	// the KIP-111 live state database pruning. The Hash component shall represent
	// the merkle hash of the node and the Counter component shall differentiate
	// nodes with the same merkle hash.
	ExtHash        [ExtHashLength]byte
	ExtHashCounter [ExtHashCounterLength]byte
)

// BytesToExtHash converts the byte array b to ExtHash.
// If len(b) is 0 or 32, then b is interpreted as a Hash and zero-extended.
// If len(b) is 39, then b is interpreted as an ExtHash.
// Otherwise, this function panics.
func BytesToExtHash(b []byte) (eh ExtHash) {
	if len(b) == 0 || len(b) == HashLength {
		return BytesToHash(b).ExtendZero()
	} else if len(b) == ExtHashLength {
		eh.SetBytes(b)
		return eh
	} else {
		logger.Crit("Invalid ExtHash bytes", "data", hexutil.Encode(b))
		return ExtHash{}
	}
}

func BytesToExtHashCounter(b []byte) (counter ExtHashCounter) {
	if len(b) == ExtHashCounterLength {
		copy(counter[:], b)
		return counter
	} else {
		logger.Crit("Invalid ExtHashCounter bytes", "data", hexutil.Encode(b))
		return ExtHashCounter{}
	}
}

func HexToExtHash(s string) ExtHash { return BytesToExtHash(FromHex(s)) }

func HexToExtHashCounter(s string) ExtHashCounter { return BytesToExtHashCounter(FromHex(s)) }

func (n ExtHashCounter) Bytes() []byte { return n[:] }

func (n ExtHashCounter) Hex() string { return hexutil.Encode(n[:]) }

func (eh ExtHash) Bytes() []byte { return eh[:] }

func (eh ExtHash) Hex() string { return hexutil.Encode(eh[:]) }

func (eh ExtHash) String() string { return eh.Hex() }

func (eh ExtHash) TerminalString() string {
	return fmt.Sprintf("%x…%x", eh[:3], eh[29:])
}

func (eh ExtHash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), eh[:])
}

func (eh *ExtHash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("ExtHash", input, eh[:])
}

func (eh *ExtHash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(extHashT, input, eh[:])
}

func (eh ExtHash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(eh[:]).MarshalText()
}

// SetBytes sets the ExtHash to the value of b.
// If b is larger than ExtHashLength, b will be cropped from the left.
// If b is smaller than ExtHashLength, b will be right aligned.
func (eh *ExtHash) SetBytes(b []byte) {
	if len(b) > ExtHashLength {
		b = b[len(b)-ExtHashLength:]
	}

	copy(eh[ExtHashLength-len(b):], b)
}

func (eh ExtHash) getShardIndex(shardMask int) int {
	return eh.Unextend().getShardIndex(shardMask)
}

func EmptyExtHash(eh ExtHash) bool {
	return EmptyHash(eh.Unextend())
}

// Unextend returns the 32 byte Hash component of an ExtHash
func (eh ExtHash) Unextend() (h Hash) {
	copy(h[:], eh[:HashLength])
	return h
}

// Counter returns the 7 byte counter component of an ExtHash
func (eh ExtHash) Counter() (counter ExtHashCounter) {
	copy(counter[:], eh[HashLength:])
	return counter
}

// IsZeroExtended returns true if the counter component of an ExtHash is zero.
// A zero counter signifies that the ExtHash is actually a Hash.
func (eh ExtHash) IsZeroExtended() bool {
	return bytes.Equal(eh.Counter().Bytes(), extHashZeroCounter[:])
}

// ResetExtHashCounterForTest sets the extHashCounter for deterministic testing
func ResetExtHashCounterForTest(counter uint64) {
	atomic.StoreUint64(&extHashLastCounter, counter)
}

func nextExtHashCounter() ExtHashCounter {
	num := atomic.AddUint64(&extHashLastCounter, 1)
	bin := make([]byte, 8)
	binary.BigEndian.PutUint64(bin, num)
	return BytesToExtHashCounter(bin[1:8])
}

// extend converts Hash to ExtHash by attaching a given counter
func (h Hash) extend(counter ExtHashCounter) (eh ExtHash) {
	copy(eh[:HashLength], h[:HashLength])
	copy(eh[HashLength:], counter[:])
	return eh
}

// Extend converts Hash to ExtHash by attaching an auto-generated counter
// Auto-generated counters must be different every time
func (h Hash) Extend() ExtHash {
	counter := nextExtHashCounter()
	eh := h.extend(counter)
	// logger.Trace("extend hash", "exthash", eh.Hex())
	return eh
}

// ExtendZero converts Hash to ExtHash by attaching the zero counter.
// A zero counter is attached to a 32-byte Hash of Trie nodes,
// later to be unextended back to a Hash.
func (h Hash) ExtendZero() ExtHash {
	return h.extend(extHashZeroCounter)
}

/////////// Address

// Address represents the 20 byte address of a Klaytn account.
type Address [AddressLength]byte

func EmptyAddress(a Address) bool {
	return a == Address{}
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}
func StringToAddress(s string) Address { return BytesToAddress([]byte(s)) }

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// IsPrecompiledContractAddress returns true if the input address is in the range of precompiled contract addresses.
func IsPrecompiledContractAddress(addr Address) bool {
	if bytes.Compare(addr.Bytes(), lastPrecompiledContractAddressHex) > 0 || addr == (Address{}) {
		return false
	}
	return true
}

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Klaytn address or not.
func IsHexAddress(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// SetBytesFromFront sets the address to the value of b.
// If len(b) is larger, take AddressLength bytes from front.
func (a *Address) SetBytesFromFront(b []byte) {
	if len(b) > AddressLength {
		b = b[:AddressLength]
	}
	copy(a[:], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

// getShardIndex returns the index of the shard.
// The address is arranged in the front or back of the array according to the initialization method.
// And the opposite is zero. In any case, to calculate the various shard index values,
// add both values and shift to calculate the shard index.
func (a Address) getShardIndex(shardMask int) int {
	data1 := int(a[AddressLength-1]) + int(a[0])
	data2 := int(a[AddressLength-2]) + int(a[1])
	return ((data2 << 8) + data1) & shardMask
}

// UnprefixedAddress allows marshaling an Address without 0x prefix.
type UnprefixedAddress Address

// UnmarshalText decodes the address from hex. The 0x prefix is optional.
func (a *UnprefixedAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAddress", input, a[:])
}

// MarshalText encodes the address as hex.
func (a UnprefixedAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

type ConnType int

const ConnTypeUndefined ConnType = -1

const (
	CONSENSUSNODE ConnType = iota
	ENDPOINTNODE
	PROXYNODE
	BOOTNODE
	UNKNOWNNODE // For error case
)

func (ct ConnType) Valid() bool {
	if int(ct) > 255 {
		return false
	}
	return true
}

func (ct ConnType) String() string {
	s := fmt.Sprintf("%d", int(ct))
	return s
}
