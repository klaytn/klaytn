package chaindatafetcher

import "github.com/klaytn/klaytn/blockchain"

type repository interface {
	InsertTransactions(event blockchain.ChainEvent) error
	InsertTokenTransfers(event blockchain.ChainEvent) error
}
