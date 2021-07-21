// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/vm/contracts.go (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/big"
	"strconv"

	"github.com/klaytn/klaytn/api/debug"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/blake2b"
	"github.com/klaytn/klaytn/crypto/bn256"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"golang.org/x/crypto/ripemd160"
)

var logger = log.NewModuleLogger(log.VM)

var (
	errInputTooShort        = errors.New("input length is too short")
	errWrongSignatureLength = errors.New("wrong signature length")
)

// PrecompiledContract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
// If you want more information about Klaytn's precompiled contracts,
// please refer https://docs.klaytn.com/smart-contract/precompiled-contracts
type PrecompiledContract interface {
	// GetRequiredGasAndComputationCost returns the gas and computation cost
	// required to execute the precompiled contract.
	GetRequiredGasAndComputationCost(input []byte) (uint64, uint64)

	// Run runs the precompiled contract
	// contract, evm is only exists in klaytn, those are not used in go-ethereum
	Run(input []byte, contract *Contract, evm *EVM) ([]byte, error)
}

// PrecompiledContractsConstantinople contains the default set of pre-compiled Klaytn
// contracts based on Ethereum Constantinople.
var PrecompiledContractsConstantinople = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}):  &ecrecover{},
	common.BytesToAddress([]byte{2}):  &sha256hash{},
	common.BytesToAddress([]byte{3}):  &ripemd160hash{},
	common.BytesToAddress([]byte{4}):  &dataCopy{},
	common.BytesToAddress([]byte{5}):  &bigModExp{},
	common.BytesToAddress([]byte{6}):  &bn256AddConstantinople{},
	common.BytesToAddress([]byte{7}):  &bn256ScalarMulConstantinople{},
	common.BytesToAddress([]byte{8}):  &bn256PairingConstantinople{},
	common.BytesToAddress([]byte{9}):  &vmLog{},
	common.BytesToAddress([]byte{10}): &feePayer{},
	common.BytesToAddress([]byte{11}): &validateSender{},
}

// DO NOT USE 0x3FD, 0x3FE, 0x3FF ADDRESSES BEFORE ISTANBUL CHANGE ACTIVATED.
// PrecompiledContractsIstanbul contains the default set of pre-compiled Klaytn
// contracts based on Ethereum Istanbul.
var PrecompiledContractsIstanbul = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}):      &ecrecover{},
	common.BytesToAddress([]byte{2}):      &sha256hash{},
	common.BytesToAddress([]byte{3}):      &ripemd160hash{},
	common.BytesToAddress([]byte{4}):      &dataCopy{},
	common.BytesToAddress([]byte{5}):      &bigModExp{},
	common.BytesToAddress([]byte{6}):      &bn256AddIstanbul{},
	common.BytesToAddress([]byte{7}):      &bn256ScalarMulIstanbul{},
	common.BytesToAddress([]byte{8}):      &bn256PairingIstanbul{},
	common.BytesToAddress([]byte{9}):      &blake2F{},
	common.BytesToAddress([]byte{3, 253}): &vmLog{},
	common.BytesToAddress([]byte{3, 254}): &feePayer{},
	common.BytesToAddress([]byte{3, 255}): &validateSender{},
}

// RunPrecompiledContract runs and evaluates the output of a precompiled contract.
func RunPrecompiledContract(p PrecompiledContract, input []byte, contract *Contract, evm *EVM) (ret []byte, computationCost uint64, err error) {
	gas, computationCost := p.GetRequiredGasAndComputationCost(input)
	if contract.UseGas(gas) {
		ret, err = p.Run(input, contract, evm)
		return ret, computationCost, err
	}
	return nil, computationCost, kerrors.ErrOutOfGas
}

// ECRECOVER implemented as a native contract.
type ecrecover struct{}

func (c *ecrecover) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	return params.EcrecoverGas, params.EcrecoverComputationCost
}

func (c *ecrecover) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	const ecRecoverInputLength = 128

	input = common.RightPadBytes(input, ecRecoverInputLength)
	// "input" is (hash, v, r, s), each 32 bytes
	// but for ecrecover we want (r, s, v)

	r := new(big.Int).SetBytes(input[64:96])
	s := new(big.Int).SetBytes(input[96:128])
	v := input[63] - 27

	// tighter sig s values input homestead only apply to tx sigs
	if !allZero(input[32:63]) || !crypto.ValidateSignatureValues(v, r, s, false) {
		return nil, nil
	}
	// We must make sure not to modify the 'input', so placing the 'v' along with
	// the signature needs to be done on a new allocation
	sig := make([]byte, 65)
	copy(sig, input[64:128])
	sig[64] = v
	// v needs to be at the end for libsecp256k1
	pubKey, err := crypto.Ecrecover(input[:32], sig)
	// make sure the public key is a valid one
	if err != nil {
		return nil, nil
	}

	// the first byte of pubkey is bitcoin heritage
	return common.LeftPadBytes(crypto.Keccak256(pubKey[1:])[12:], 32), nil
}

// SHA256 implemented as a native contract.
type sha256hash struct{}

// GetRequiredGasAndComputationCost returns the gas required to execute the pre-compiled contract
// and the computation cost of the precompiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *sha256hash) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	n32Bytes := uint64(len(input)+31) / 32

	return n32Bytes*params.Sha256PerWordGas + params.Sha256BaseGas,
		n32Bytes*params.Sha256PerWordComputationCost + params.Sha256BaseComputationCost
}
func (c *sha256hash) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	h := sha256.Sum256(input)
	return h[:], nil
}

// RIPEMD160 implemented as a native contract.
type ripemd160hash struct{}

// GetRequiredGasAndComputationCost returns the gas required to execute the pre-compiled contract
// and the computation cost of the precompiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *ripemd160hash) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	n32Bytes := uint64(len(input)+31) / 32

	return n32Bytes*params.Ripemd160PerWordGas + params.Ripemd160BaseGas,
		n32Bytes*params.Ripemd160PerWordComputationCost + params.Ripemd160BaseComputationCost
}
func (c *ripemd160hash) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	ripemd := ripemd160.New()
	ripemd.Write(input)
	return common.LeftPadBytes(ripemd.Sum(nil), 32), nil
}

// data copy implemented as a native contract.
type dataCopy struct{}

// GetRequiredGasAndComputationCost returns the gas required to execute the pre-compiled contract
// and the computation cost of the precompiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *dataCopy) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	n32Bytes := uint64(len(input)+31) / 32
	return n32Bytes*params.IdentityPerWordGas + params.IdentityBaseGas,
		n32Bytes*params.IdentityPerWordComputationCost + params.IdentityBaseComputationCost
}
func (c *dataCopy) Run(in []byte, contract *Contract, evm *EVM) ([]byte, error) {
	return in, nil
}

// bigModExp implements a native big integer exponential modular operation.
type bigModExp struct{}

var (
	big1      = big.NewInt(1)
	big4      = big.NewInt(4)
	big8      = big.NewInt(8)
	big16     = big.NewInt(16)
	big32     = big.NewInt(32)
	big64     = big.NewInt(64)
	big96     = big.NewInt(96)
	big480    = big.NewInt(480)
	big1024   = big.NewInt(1024)
	big3072   = big.NewInt(3072)
	big199680 = big.NewInt(199680)
)

// GetRequiredGasAndComputationCost returns the gas required to execute the pre-compiled contract
// and the computation cost of the precompiled contract.
func (c *bigModExp) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	var (
		baseLen = new(big.Int).SetBytes(getData(input, 0, 32))
		expLen  = new(big.Int).SetBytes(getData(input, 32, 32))
		modLen  = new(big.Int).SetBytes(getData(input, 64, 32))
	)
	if len(input) > 96 {
		input = input[96:]
	} else {
		input = input[:0]
	}
	// Retrieve the head 32 bytes of exp for the adjusted exponent length
	var expHead *big.Int
	if big.NewInt(int64(len(input))).Cmp(baseLen) <= 0 {
		expHead = new(big.Int)
	} else {
		if expLen.Cmp(big32) > 0 {
			expHead = new(big.Int).SetBytes(getData(input, baseLen.Uint64(), 32))
		} else {
			expHead = new(big.Int).SetBytes(getData(input, baseLen.Uint64(), expLen.Uint64()))
		}
	}
	// Calculate the adjusted exponent length
	var msb int
	if bitlen := expHead.BitLen(); bitlen > 0 {
		msb = bitlen - 1
	}
	adjExpLen := new(big.Int)
	if expLen.Cmp(big32) > 0 {
		adjExpLen.Sub(expLen, big32)
		adjExpLen.Mul(big8, adjExpLen)
	}
	adjExpLen.Add(adjExpLen, big.NewInt(int64(msb)))

	// Calculate the gas cost of the operation
	gas := new(big.Int).Set(math.BigMax(modLen, baseLen))
	switch {
	case gas.Cmp(big64) <= 0:
		gas.Mul(gas, gas)
	case gas.Cmp(big1024) <= 0:
		gas = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(gas, gas), big4),
			new(big.Int).Sub(new(big.Int).Mul(big96, gas), big3072),
		)
	default:
		gas = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(gas, gas), big16),
			new(big.Int).Sub(new(big.Int).Mul(big480, gas), big199680),
		)
	}
	gas.Mul(gas, math.BigMax(adjExpLen, big1))
	gas.Div(gas, new(big.Int).SetUint64(params.ModExpQuadCoeffDiv))

	if gas.BitLen() > 64 {
		return math.MaxUint64, math.MaxUint64
	}
	return gas.Uint64(), (gas.Uint64() / 100) + params.BigModExpBaseComputationCost
}

func (c *bigModExp) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	var (
		baseLen = new(big.Int).SetBytes(getData(input, 0, 32)).Uint64()
		expLen  = new(big.Int).SetBytes(getData(input, 32, 32)).Uint64()
		modLen  = new(big.Int).SetBytes(getData(input, 64, 32)).Uint64()
	)
	if len(input) > 96 {
		input = input[96:]
	} else {
		input = input[:0]
	}
	// Handle a special case when both the base and mod length is zero
	if baseLen == 0 && modLen == 0 {
		return []byte{}, nil
	}
	// Retrieve the operands and execute the exponentiation
	var (
		base = new(big.Int).SetBytes(getData(input, 0, baseLen))
		exp  = new(big.Int).SetBytes(getData(input, baseLen, expLen))
		mod  = new(big.Int).SetBytes(getData(input, baseLen+expLen, modLen))
	)
	if mod.BitLen() == 0 {
		// Modulo 0 is undefined, return zero
		return common.LeftPadBytes([]byte{}, int(modLen)), nil
	}
	return common.LeftPadBytes(base.Exp(base, exp, mod).Bytes(), int(modLen)), nil
}

// newCurvePoint unmarshals a binary blob into a bn256 elliptic curve point,
// returning it, or an error if the point is invalid.
func newCurvePoint(blob []byte) (*bn256.G1, error) {
	p := new(bn256.G1)
	if _, err := p.Unmarshal(blob); err != nil {
		return nil, err
	}
	return p, nil
}

// newTwistPoint unmarshals a binary blob into a bn256 elliptic curve point,
// returning it, or an error if the point is invalid.
func newTwistPoint(blob []byte) (*bn256.G2, error) {
	p := new(bn256.G2)
	if _, err := p.Unmarshal(blob); err != nil {
		return nil, err
	}
	return p, nil
}

// runBn256Add implements the Bn256Add precompile, referenced by both
// Constantinople and Istanbul operations.
func runBn256Add(input []byte) ([]byte, error) {
	x, err := newCurvePoint(getData(input, 0, 64))
	if err != nil {
		return nil, err
	}
	y, err := newCurvePoint(getData(input, 64, 64))
	if err != nil {
		return nil, err
	}
	res := new(bn256.G1)
	res.Add(x, y)
	return res.Marshal(), nil
}

// bn256Add implements a native elliptic curve point addition conforming to
// Istanbul consensus rules.
type bn256AddIstanbul struct{}

func (c *bn256AddIstanbul) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	return params.Bn256AddGasIstanbul, params.Bn256AddComputationCost
}

func (c *bn256AddIstanbul) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	return runBn256Add(input)
}

// bn256AddByzantium implements a native elliptic curve point addition
// conforming to Byzantium consensus rules.
type bn256AddConstantinople struct{}

func (c *bn256AddConstantinople) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	return params.Bn256AddGasConstantinople, params.Bn256AddComputationCost
}

func (c *bn256AddConstantinople) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	return runBn256Add(input)
}

// runBn256ScalarMul implements the Bn256ScalarMul precompile, referenced by
// both Constantionple and Istanbul operations.
func runBn256ScalarMul(input []byte) ([]byte, error) {
	p, err := newCurvePoint(getData(input, 0, 64))
	if err != nil {
		return nil, err
	}
	res := new(bn256.G1)
	res.ScalarMult(p, new(big.Int).SetBytes(getData(input, 64, 32)))
	return res.Marshal(), nil
}

// bn256ScalarMulIstanbul implements a native elliptic curve scalar
// multiplication conforming to Istanbul consensus rules.
type bn256ScalarMulIstanbul struct{}

func (c *bn256ScalarMulIstanbul) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	return params.Bn256ScalarMulGasIstanbul, params.Bn256ScalarMulComputationCost
}

func (c *bn256ScalarMulIstanbul) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	return runBn256ScalarMul(input)
}

// bn256ScalarMulByzantium implements a native elliptic curve scalar
// multiplication conforming to Byzantium consensus rules.
type bn256ScalarMulConstantinople struct{}

func (c *bn256ScalarMulConstantinople) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	return params.Bn256ScalarMulGasConstantinople, params.Bn256ScalarMulComputationCost
}

func (c *bn256ScalarMulConstantinople) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	return runBn256ScalarMul(input)
}

var (
	// true32Byte is returned if the bn256 pairing check succeeds.
	true32Byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// false32Byte is returned if the bn256 pairing check fails.
	false32Byte = make([]byte, 32)

	// errBadPairingInput is returned if the bn256 pairing input is invalid.
	errBadPairingInput = errors.New("bad elliptic curve pairing size")
)

// runBn256Pairing implements the Bn256Pairing precompile, referenced by both
// Byzantium and Istanbul operations.
func runBn256Pairing(input []byte) ([]byte, error) {
	// Handle some corner cases cheaply
	if len(input)%192 > 0 {
		return nil, errBadPairingInput
	}
	// Convert the input into a set of coordinates
	var (
		cs []*bn256.G1
		ts []*bn256.G2
	)
	for i := 0; i < len(input); i += 192 {
		c, err := newCurvePoint(input[i : i+64])
		if err != nil {
			return nil, err
		}
		t, err := newTwistPoint(input[i+64 : i+192])
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
		ts = append(ts, t)
	}
	// Execute the pairing checks and return the results
	if bn256.PairingCheck(cs, ts) {
		return true32Byte, nil
	}
	return false32Byte, nil
}

// bn256PairingIstanbul implements a pairing pre-compile for the bn256 curve
// conforming to Istanbul consensus rules.
type bn256PairingIstanbul struct{}

func (c *bn256PairingIstanbul) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	numParings := uint64(len(input) / 192)
	return params.Bn256PairingBaseGasIstanbul + numParings*params.Bn256PairingPerPointGasIstanbul,
		params.Bn256ParingBaseComputationCost + numParings*params.Bn256ParingPerPointComputationCost
}

func (c *bn256PairingIstanbul) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	return runBn256Pairing(input)
}

// bn256PairingConstantinople implements a pairing pre-compile for the bn256 curve
// conforming to Constantinople consensus rules.
type bn256PairingConstantinople struct{}

func (c *bn256PairingConstantinople) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	numParings := uint64(len(input) / 192)
	return params.Bn256PairingBaseGasConstantinople + numParings*params.Bn256PairingPerPointGasConstantinople,
		params.Bn256ParingBaseComputationCost + numParings*params.Bn256ParingPerPointComputationCost
}

func (c *bn256PairingConstantinople) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	return runBn256Pairing(input)
}

type blake2F struct{}

const (
	blake2FInputLength        = 213
	blake2FFinalBlockBytes    = byte(1)
	blake2FNonFinalBlockBytes = byte(0)
)

var (
	errBlake2FInvalidInputLength = errors.New("invalid input length")
	errBlake2FInvalidFinalFlag   = errors.New("invalid final flag")
)

func (c *blake2F) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	// If the input is malformed, we can't calculate the gas, return 0 and let the
	// actual call choke and fault.
	if len(input) != blake2FInputLength {
		return 0, 0
	}
	gas := uint64(binary.BigEndian.Uint32(input[0:4]))
	return gas, params.Blake2bBaseComputationCost + params.Blake2bScaleComputationCost*gas
}

func (c *blake2F) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	// Make sure the input is valid (correct length and final flag)
	if len(input) != blake2FInputLength {
		return nil, errBlake2FInvalidInputLength
	}
	if input[212] != blake2FNonFinalBlockBytes && input[212] != blake2FFinalBlockBytes {
		return nil, errBlake2FInvalidFinalFlag
	}
	// Parse the input into the Blake2b call parameters
	var (
		rounds = binary.BigEndian.Uint32(input[0:4])
		final  = (input[212] == blake2FFinalBlockBytes)

		h [8]uint64
		m [16]uint64
		t [2]uint64
	)
	for i := 0; i < 8; i++ {
		offset := 4 + i*8
		h[i] = binary.LittleEndian.Uint64(input[offset : offset+8])
	}
	for i := 0; i < 16; i++ {
		offset := 68 + i*8
		m[i] = binary.LittleEndian.Uint64(input[offset : offset+8])
	}
	t[0] = binary.LittleEndian.Uint64(input[196:204])
	t[1] = binary.LittleEndian.Uint64(input[204:212])

	// Execute the compression function, extract and return the result
	blake2b.F(&h, m, t, final, rounds)

	output := make([]byte, 64)
	for i := 0; i < 8; i++ {
		offset := i * 8
		binary.LittleEndian.PutUint64(output[offset:offset+8], h[i])
	}
	return output, nil
}

// vmLog implemented as a native contract.
type vmLog struct{}

// GetRequiredGasAndComputationCost returns the gas required to execute the pre-compiled contract
// and the computation cost of the precompiled contract.
func (c *vmLog) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	l := uint64(len(input))
	return l*params.VMLogPerByteGas + params.VMLogBaseGas,
		l*params.VMLogPerByteComputationCost + params.VMLogBaseComputationCost
}

// Runs the vmLog contract.
func (c *vmLog) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	if (params.VMLogTarget & params.VMLogToFile) != 0 {
		prefix := "tx=" + evm.StateDB.GetTxHash().String() + " caller=" + contract.CallerAddress.String() + " msg="
		debug.Handler.WriteVMLog(prefix + string(input))
	}
	if (params.VMLogTarget & params.VMLogToStdout) != 0 {
		logger.Debug("vmlog", "tx", evm.StateDB.GetTxHash().String(),
			"caller", contract.CallerAddress.String(), "msg", strconv.QuoteToASCII(string(input)))
	}
	return nil, nil
}

type feePayer struct{}

func (c *feePayer) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	return params.FeePayerGas, params.FeePayerComputationCost
}

func (c *feePayer) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	return contract.FeePayerAddress.Bytes(), nil
}

type validateSender struct{}

func (c *validateSender) GetRequiredGasAndComputationCost(input []byte) (uint64, uint64) {
	numSigs := uint64(len(input) / common.SignatureLength)
	return numSigs * params.ValidateSenderGas,
		numSigs*params.ValidateSenderPerSigComputationCost + params.ValidateSenderBaseComputationCost
}

func (c *validateSender) Run(input []byte, contract *Contract, evm *EVM) ([]byte, error) {
	if err := c.validateSender(input, evm.StateDB, evm.BlockNumber.Uint64()); err != nil {
		// If return error makes contract execution failed, do not return the error.
		// Instead, print log.
		logger.Trace("validateSender failed", "err", err)
		return []byte{0}, nil
	}
	return []byte{1}, nil
}

func (c *validateSender) validateSender(input []byte, picker types.AccountKeyPicker, currentBlockNumber uint64) error {
	ptr := input

	// Parse the first 20 bytes. They represent an address to be verified.
	if len(ptr) < common.AddressLength {
		return errInputTooShort
	}
	from := common.BytesToAddress(input[0:common.AddressLength])
	ptr = ptr[common.AddressLength:]

	// Parse the next 32 bytes. They represent a message which was used to generate signatures.
	if len(ptr) < common.HashLength {
		return errInputTooShort
	}
	msg := ptr[0:common.HashLength]
	ptr = ptr[common.HashLength:]

	// Parse remaining bytes. The length should be divided by common.SignatureLength.
	if len(ptr)%common.SignatureLength != 0 {
		return errWrongSignatureLength
	}

	numSigs := len(ptr) / common.SignatureLength
	pubs := make([]*ecdsa.PublicKey, numSigs)
	for i := 0; i < numSigs; i++ {
		p, err := crypto.Ecrecover(msg, ptr[0:common.SignatureLength])
		if err != nil {
			return err
		}
		pubs[i], err = crypto.UnmarshalPubkey(p)
		if err != nil {
			return err
		}
		ptr = ptr[common.SignatureLength:]
	}

	k := picker.GetKey(from)
	if err := accountkey.ValidateAccountKey(currentBlockNumber, from, k, pubs, accountkey.RoleTransaction); err != nil {
		return err
	}

	return nil
}
