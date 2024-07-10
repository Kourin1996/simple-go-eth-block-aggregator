package parser

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"context"
	"math/big"
)

type EthClient interface {
	GetLatestHeight(ctx context.Context) (*big.Int, error)
	GetBlock(context.Context, big.Int) (*types.Block, error)
}

type EthTransactionStorage interface {
	InsertTransactions([]*types.Transaction) error
	GetTransactionsByAddress(string) []types.Transaction
}
