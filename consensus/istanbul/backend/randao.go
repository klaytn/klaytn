package backend

import (
	"bytes"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/bls"
)

// For testing without KIP-113 contract setup
type BlsPubkeyProvider interface {
	GetBlsPubkey(chain consensus.ChainReader, proposer common.Address) (bls.PublicKey, error)
	ResetBlsCache()
}

// The default implementation for BlsPubkeyFunc.
// Queries KIP-113 contract and verifies the PoP.
func (sb *backend) GetBlsPubkey(chain consensus.ChainReader, proposer common.Address) (bls.PublicKey, error) {
	logger.Crit("not implemented")
	return nil, errNoBlsPub
}

func (sb *backend) ResetBlsCache() {
	logger.Crit("not implemented")
}

// Calculate KIP-114 Randao header fields
// https://github.com/klaytn/kips/blob/kip114/KIPs/kip-114.md
func (sb *backend) CalcRandao(number *big.Int, prevMixHash []byte) ([]byte, []byte, error) {
	if sb.blsSecretKey == nil {
		return nil, nil, errNoBlsKey
	}
	if len(prevMixHash) != 32 {
		logger.Error("invalid prevMixHash", "number", number.Uint64(), "prevMixHash", hexutil.Encode(prevMixHash))
		return nil, nil, errInvalidRandaoFields
	}

	// block_num_to_bytes() = num.to_bytes(32, byteorder="big")
	msg := calcRandaoMsg(number)

	// calc_random_reveal() = sign(privateKey, headerNumber)
	randomReveal := bls.Sign(sb.blsSecretKey, msg[:]).Marshal()

	// calc_mix_hash() = xor(prevMixHash, keccak256(randomReveal))
	mixHash := calcMixHash(randomReveal, prevMixHash)

	return randomReveal, mixHash, nil
}

func (sb *backend) VerifyRandao(chain consensus.ChainReader, header *types.Header, prevMixHash []byte) error {
	proposer, err := sb.Author(header)
	if err != nil {
		return err
	}

	// [proposerPubkey, proposerPop] = get_proposer_pubkey_pop()
	// if not pop_verify(proposerPubkey, proposerPop): return False
	proposerPub, err := sb.blsPubkeyProvider.GetBlsPubkey(chain, proposer)
	if err != nil {
		return err
	}

	// if not verify(proposerPubkey, newHeader.number, newHeader.randomReveal): return False
	sig := header.RandomReveal
	msg := calcRandaoMsg(header.Number)
	ok, err := bls.VerifySignature(sig, msg, proposerPub)
	if err != nil {
		return err
	} else if !ok {
		return errInvalidRandaoFields
	}

	// if not newHeader.mixHash == calc_mix_hash(prevMixHash, newHeader.randomReveal): return False
	mixHash := calcMixHash(header.RandomReveal, prevMixHash)
	if !bytes.Equal(header.MixHash, mixHash) {
		return errInvalidRandaoFields
	}

	return nil
}

// block_num_to_bytes() = num.to_bytes(32, byteorder="big")
func calcRandaoMsg(number *big.Int) common.Hash {
	return common.BytesToHash(number.Bytes())
}

// calc_mix_hash() = xor(prevMixHash, keccak256(randomReveal))
func calcMixHash(randomReveal, prevMixHash []byte) []byte {
	mixHash := make([]byte, 32)
	revealHash := crypto.Keccak256(randomReveal)
	for i := 0; i < 32; i++ {
		mixHash[i] = prevMixHash[i] ^ revealHash[i]
	}
	return mixHash
}
