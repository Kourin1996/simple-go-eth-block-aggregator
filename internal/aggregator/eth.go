package aggregator

import (
	"context"
	"math/big"
)

type EthTransactionAggregator struct {
	ethClient EthClient
	storage   EthTransactionStorage
}

func New(
	ethClient EthClient,
	storage EthTransactionStorage,
) *EthTransactionAggregator {
	return &EthTransactionAggregator{
		ethClient: ethClient,
		storage:   storage,
	}
}

func (a *EthTransactionAggregator) Start(
	beginningHeight *big.Int,
) error {
	panic("not implemented!")
	return nil
}

func (a *EthTransactionAggregator) Stop(ctx context.Context) error {
	panic("not implemented!")
	return nil
}
