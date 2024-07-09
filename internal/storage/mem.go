package storage

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/types"
)

type InMemoryTransactionStorage struct {}

// Hash of transactions  hash -> transaction
// Incoming transactions address -> hash list
// Outgoing transactions address -> hash list

func (s *InMemoryTransactionStorage) InsertTransactions(
	txs []types.Transaction,
) error {
	// 1. sort by block height, index ???
	// 2. save

	return nil
}

func (s *InMemoryTransactionStorage) QueryOutboundTransactions(
	target types.Address,
) ([]types.Transaction, error) {
	panic("not implemented!")

	return nil, nil
}

func (s *InMemoryTransactionStorage) QueryInboundTransactions() ([]types.Transaction, error) {
	panic("not implemented!")

	return nil, nil
}