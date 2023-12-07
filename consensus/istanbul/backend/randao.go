package backend

import (
	"bytes"
	"errors"
	"math/big"

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain/system"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/klaytn/klaytn/params"
)

// For testing without KIP-113 contract setup
type BlsPubkeyProvider interface {
	// num should be the header number of the block to be verified.
	// Thus, since the state of num does not exist, the state of num-1 must be used.
	GetBlsPubkey(chain consensus.ChainReader, proposer common.Address, num *big.Int) (bls.PublicKey, error)
	ResetBlsCache()
}

type ChainBlsPubkeyProvider struct {
	cache *lru.ARCCache // Cached BlsPublicKeyInfos
}

func newChainBlsPubkeyProvider() *ChainBlsPubkeyProvider {
	cache, _ := lru.NewARC(128)
	return &ChainBlsPubkeyProvider{
		cache: cache,
	}
}

// GetBlsPubkey is the default implementation for BlsPubkeyFunc.
// Queries KIP-113 contract and verifies the PoP.
func (p *ChainBlsPubkeyProvider) GetBlsPubkey(chain consensus.ChainReader, proposer common.Address, num *big.Int) (bls.PublicKey, error) {
	infos, err := p.getBlsInfos(chain, num)
	if err != nil {
		return nil, err
	}

	info, ok := infos[proposer]
	if !ok {
		return nil, errNoBlsPub
	}
	if info.VerifyErr != nil {
		return nil, info.VerifyErr
	}
	return bls.PublicKeyFromBytes(info.PublicKey)
}

// getBlsInfosv returns all registered BLS info at the given block number.
// It retrieves cache first, and then retrieves the storage of KIP113 contract.
func (p *ChainBlsPubkeyProvider) getBlsInfos(chain consensus.ChainReader, num *big.Int) (system.BlsPublicKeyInfos, error) {
	if item, ok := p.cache.Get(num.Uint64()); ok {
		logger.Trace("BlsPublicKeyInfos cache hit", "number", num.Uint64())
		return item.(system.BlsPublicKeyInfos), nil
	}

	backend := backends.NewBlockchainContractBackend(chain, nil, nil)
	if common.Big0.Cmp(num) == 0 {
		return nil, errors.New("num cannot be zero")
	}
	parentNum := new(big.Int).Sub(num, common.Big1)

	var kip113Addr common.Address
	// Because the system contract Registry is installed at Finalize() of RandaoForkBlock,
	// it is not possible to read KIP113 address from the Registry at RandaoForkBlock.
	// Hence the ChainConfig fallback.
	if chain.Config().IsRandaoForkBlock(num) {
		var ok bool
		kip113Addr, ok = chain.Config().RandaoRegistry.Records[system.Kip113Name]
		if !ok {
			return nil, errors.New("KIP113 address not set in ChainConfig")
		}
	} else if chain.Config().IsRandaoForkEnabled(num) {
		var err error
		kip113Addr, err = system.ReadRegistryActiveAddr(backend, system.Kip113Name, parentNum)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Cannot read KIP113 address from registry before Randao fork")
	}

	infos, err := system.ReadKip113All(backend, kip113Addr, parentNum)
	if err != nil {
		return nil, err
	}
	logger.Trace("BlsPublicKeyInfos cache miss", "number", num.Uint64())
	p.cache.Add(num.Uint64(), infos)

	return infos, nil
}

func (p *ChainBlsPubkeyProvider) ResetBlsCache() {
	p.cache.Purge()
}

// CalcRandao calculates Randao-related header values specified in KIP-114.
// https://github.com/klaytn/kips/blob/kip114/KIPs/kip-114.md
func (sb *backend) CalcRandao(number *big.Int, prevMixHash []byte) ([]byte, []byte, error) {
	if sb.blsSecretKey == nil {
		return nil, nil, errNoBlsKey
	}
	if prevMixHash == nil {
		prevMixHash = params.ZeroMixHash
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

// VerifyRandao verifies whether header.RandomReveal is the same as expected.
func (sb *backend) VerifyRandao(chain consensus.ChainReader, header *types.Header, prevMixHash []byte) error {
	if header.Number.Sign() == 0 {
		return nil // Do not verify genesis block
	}
	if header.RandomReveal == nil || header.MixHash == nil {
		return errInvalidRandaoFields
	}
	// The following condition is only true when header's block number is the Randao hardfork block number.
	// Because of the above condition, prevMixHash cannot be nil after Randao hardfork block.
	if prevMixHash == nil {
		prevMixHash = params.ZeroMixHash
	}

	proposer, err := sb.Author(header)
	if err != nil {
		return err
	}

	// [proposerPubkey, proposerPop] = get_proposer_pubkey_pop()
	// if not pop_verify(proposerPubkey, proposerPop): return False
	proposerPub, err := sb.blsPubkeyProvider.GetBlsPubkey(chain, proposer, header.Number)
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
