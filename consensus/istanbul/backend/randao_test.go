package backend

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/stretchr/testify/assert"
)

// Test low-level computation components
func TestCalcRandao(t *testing.T) {
	var (
		skhex = hexutil.MustDecode("0x6c605527c8e4f31c959478801d51384d690a22dfc6438604646f7709032c893a")
		sk, _ = bls.SecretKeyFromBytes(skhex)
		pk    = sk.PublicKey()

		// block_num_to_bytes() = num.to_bytes(32, byteorder="big")
		num = big.NewInt(31337)
		msg = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000007a69")

		// mix2 = xor(mix1, keccak256(sig))
		sig  = hexutil.MustDecode("0xadfe25ced45819332cbf088f01cdd2807686dd6309b11d7440237dd623624f401d4753747f5fb92374235e997edcd18318bae2806a1675b1e685e792abd1fbdf5c50ec1e148cc7fe861984d8bc3204c1b2136725b176902bc52eeb595919df3b")
		mix1 = hexutil.MustDecode("0x8019df1a2a9f833dc7f400a15b33e54a5c80295165c5953dc23891aab9203810")
		mix2 = hexutil.MustDecode("0x8772d58248bdf34e81ecbf36f28299cfa758b61ccf3f64e1dc0646687a55892f")
	)

	// Calculate RandomReveal and MixHash
	assert.Equal(t, msg, calcRandaoMsg(num))
	assert.Equal(t, sig, bls.Sign(sk, msg[:]).Marshal())
	assert.Equal(t, mix2, calcMixHash(sig, mix1))

	// Verify signature
	ok, err := bls.VerifySignature(sig, msg, pk)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestRandao_Prepare(t *testing.T) {
	config := testRandaoConfig.Copy()

	ctx := newTestContext(1, config, nil)
	chain, engine := ctx.chain, ctx.engine
	defer ctx.Cleanup()

	header := ctx.MakeHeader(chain.Genesis())
	assert.Nil(t, engine.Prepare(chain, header))

	assert.Len(t, header.RandomReveal, 96)
	assert.Len(t, header.MixHash, 32)
	assert.NotEqual(t, header.RandomReveal, make([]byte, 96))
	assert.NotEqual(t, header.MixHash, make([]byte, 32))
}

func TestRandao_Verify(t *testing.T) {
	config := testRandaoConfig.Copy()

	ctx := newTestContext(1, config, nil)
	chain, engine := ctx.chain, ctx.engine
	defer ctx.Cleanup()

	block := ctx.MakeBlockWithCommittedSeals(chain.Genesis())
	header := block.Header()
	assert.Nil(t, engine.VerifyHeader(chain, header, false))
}
