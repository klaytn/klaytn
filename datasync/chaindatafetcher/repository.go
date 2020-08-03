package chaindatafetcher

import "github.com/klaytn/klaytn/blockchain"

type Repository interface {
	InsertTransactions(event blockchain.ChainEvent) error
}
