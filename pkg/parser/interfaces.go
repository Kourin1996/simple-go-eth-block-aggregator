package parser

import (
	"context"
	"math/big"

	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"
)

type EthClient interface {
	GetBlockNumber(ctx context.Context) (*big.Int, error)
	GetBlockByNumber(context.Context, big.Int, bool) (*types.Block, error)
}

type EthTransactionStorage interface {
	InsertTransactions([]*types.Transaction) error
	GetTransactionsByAddress(string) []types.Transaction
}
