package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepare_Randao(t *testing.T) {
	config := testRandaoConfig.Copy()

	ctx := newTestContext(1, config, nil)
	chain, engine := ctx.chain, ctx.engine
	defer ctx.Cleanup()

	header := ctx.MakeHeader(chain.Genesis())
	assert.Nil(t, engine.Prepare(chain, header))

	assert.Len(t, header.RandomReveal, 96)
	assert.Len(t, header.MixHash, 32)
}

func TestVerify_Randao(t *testing.T) {
	config := testRandaoConfig.Copy()

	ctx := newTestContext(1, config, nil)
	chain, engine := ctx.chain, ctx.engine
	defer ctx.Cleanup()

	block := ctx.MakeBlockWithCommittedSeals(chain.Genesis())
	header := block.Header()
	assert.Nil(t, engine.VerifyHeader(chain, header, false))
}
