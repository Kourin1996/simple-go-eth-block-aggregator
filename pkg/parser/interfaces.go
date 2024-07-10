package parser

import (
	"context"
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"math/big"
)

type EthClient interface {
	GetLatestHeight(ctx context.Context) (*big.Int, error)
	GetBlock(context.Context, big.Int) (*types.Block, error)
}

type EthTransactionStorage interface {
	InsertTransactions([]*types.Transaction) error
	GetTransactionsByAddress(string) []*types.Transaction
}
