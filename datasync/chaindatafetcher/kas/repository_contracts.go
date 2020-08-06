package kas

import (
	"fmt"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"strings"
	"time"
)

// filterKIPContracts filters the deployed contracts to KIP7, KIP17 and others.
func filterKIPContracts(api *api.PublicBlockChainAPI, event blockchain.ChainEvent) ([]*FT, []*NFT, []*Contract, error) {
	var (
		kip7s  []*FT
		kip17s []*NFT
		others []*Contract
	)
	caller := newContractCaller(api)
	for _, receipt := range event.Receipts {
		if receipt.ContractAddress == (common.Address{}) {
			continue
		}
		blockNumber := event.Block.Number()
		contract := receipt.ContractAddress
		isKIP13, err := caller.isKIP13(contract, blockNumber)
		if err != nil {
			return nil, nil, nil, err
		} else if !isKIP13 {
			others = append(others, &Contract{Address: contract.Bytes()})
			continue
		}

		if isKIP7, err := caller.isKIP7(contract, blockNumber); err != nil {
			return nil, nil, nil, err
		} else if isKIP7 {
			kip7s = append(kip7s, &FT{Address: contract.Bytes()})
			continue
		}

		if isKIP17, err := caller.isKIP17(contract, blockNumber); err != nil {
			return nil, nil, nil, err
		} else if isKIP17 {
			kip17s = append(kip17s, &NFT{Address: contract.Bytes()})
			continue
		}
		others = append(others, &Contract{Address: contract.Bytes()})
	}
	return kip7s, kip17s, others, nil
}

// InsertContracts inserts deployed contracts in the given chain event into KAS database.
func (r *repository) InsertContracts(event blockchain.ChainEvent) error {
	kip7s, kip17s, others, err := filterKIPContracts(r.blockchainApi, event)
	if err != nil {
		logger.Error("Failed to filter KIP contracts", "err", err, "blockNumber", event.Block.NumberU64())
		return err
	}

	if err := r.insertFTs(kip7s); err != nil {
		logger.Error("Failed to insert KIP7 contracts", "err", err, "blockNumber", event.Block.NumberU64(), "numKIP7s", len(kip7s))
		return err
	}

	if err := r.insertNFTs(kip17s); err != nil {
		logger.Error("Failed to insert KIP17 contracts", "err", err, "blockNumber", event.Block.NumberU64(), "numKIP17s", len(kip17s))
		return err
	}

	if err := r.insertContracts(others); err != nil {
		logger.Error("Failed to insert other contracts", "err", err, "blockNumber", event.Block.NumberU64(), "numContracts", len(others))
		return err
	}

	return nil
}

// insertContracts inserts the contracts which are divided into chunkUnit because of max number of placeholders.
func (r *repository) insertContracts(contracts []*Contract) error {
	chunkUnit := maxPlaceholders / placeholdersPerContractItem
	var chunks []*Contract

	for contracts != nil {
		if placeholdersPerContractItem*len(contracts) > maxPlaceholders {
			chunks = contracts[:chunkUnit]
			contracts = contracts[chunkUnit:]
		} else {
			chunks = contracts
			contracts = nil
		}

		if err := r.bulkInsertContracts(chunks); err != nil {
			return err
		}
	}

	return nil
}

// bulkInsertTransactions inserts the given contracts in multiple rows at once.
func (r *repository) bulkInsertContracts(contracts []*Contract) error {
	if len(contracts) == 0 {
		return nil
	}
	var valueStrings []string
	var valueArgs []interface{}

	for _, ft := range contracts {
		valueStrings = append(valueStrings, "(?)")
		valueArgs = append(valueArgs, ft.Address)
	}

	rawQuery := `
			INSERT INTO contract(address)
			VALUES %s
			ON DUPLICATE KEY
			UPDATE address=address`
	query := fmt.Sprintf(rawQuery, strings.Join(valueStrings, ","))

	if _, err := r.db.DB().Exec(query, valueArgs...); err != nil {
		return err
	}
	return nil
}

// insertFTs inserts the FT contracts which are divided into chunkUnit because of max number of placeholders.
func (r *repository) insertFTs(fts []*FT) error {
	chunkUnit := maxPlaceholders / placeholdersPerFTItem
	var chunks []*FT

	for fts != nil {
		if placeholdersPerFTItem*len(fts) > maxPlaceholders {
			chunks = fts[:chunkUnit]
			fts = fts[chunkUnit:]
		} else {
			chunks = fts
			fts = nil
		}

		if err := r.bulkInsertFTs(chunks); err != nil {
			return err
		}
	}

	return nil
}

// bulkInsertFTs inserts the given FT contracts in multiple rows at once.
func (r *repository) bulkInsertFTs(fts []*FT) error {
	if len(fts) == 0 {
		return nil
	}
	var valueStrings []string
	var valueArgs []interface{}

	now := time.Now()
	for _, ft := range fts {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, ft.Address)
		valueArgs = append(valueArgs, &now)
		valueArgs = append(valueArgs, &now)
	}

	rawQuery := `
			INSERT INTO kct_ft_metadata(address, createdAt, updatedAt)
			VALUES %s
			ON DUPLICATE KEY
			UPDATE address=address`
	query := fmt.Sprintf(rawQuery, strings.Join(valueStrings, ","))

	if _, err := r.db.DB().Exec(query, valueArgs...); err != nil {
		return err
	}
	return nil
}

// insertNFTs inserts the NFT contracts which are divided into chunkUnit because of max number of placeholders.
func (r *repository) insertNFTs(nfts []*NFT) error {
	chunkUnit := maxPlaceholders / placeholdersPerNFTItem
	var chunks []*NFT

	for nfts != nil {
		if placeholdersPerNFTItem*len(nfts) > maxPlaceholders {
			chunks = nfts[:chunkUnit]
			nfts = nfts[chunkUnit:]
		} else {
			chunks = nfts
			nfts = nil
		}

		if err := r.bulkInsertNFTs(chunks); err != nil {
			return err
		}
	}

	return nil
}

// bulkInsertNFTs inserts the given NFT contracts in multiple rows at once.
func (r *repository) bulkInsertNFTs(nfts []*NFT) error {
	if len(nfts) == 0 {
		return nil
	}
	var valueStrings []string
	var valueArgs []interface{}

	now := time.Now()
	for _, nft := range nfts {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, nft.Address)
		valueArgs = append(valueArgs, &now)
		valueArgs = append(valueArgs, &now)
	}

	rawQuery := `
			INSERT INTO kct_nft_metadata(address, createdAt, updatedAt)
			VALUES %s
			ON DUPLICATE KEY
			UPDATE address=address`
	query := fmt.Sprintf(rawQuery, strings.Join(valueStrings, ","))

	if _, err := r.db.DB().Exec(query, valueArgs...); err != nil {
		return err
	}
	return nil
}
