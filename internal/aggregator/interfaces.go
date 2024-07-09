package aggregator

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"math/big"
)

type EthClient interface {
	GetLatestHeight() (*big.Int, error)
	GetBlocks(from big.Int, to *big.Int) ([]types.Block, error)
}

type EthTransactionStorage interface {
	InsertTransactions([]types.Transaction) error
}
