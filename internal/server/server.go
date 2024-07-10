package server

import (
	"context"
)

type EthTransactionsServer struct {
	storage TransactionStorage
}

func New(storage TransactionStorage) *EthTransactionsServer {
	return &EthTransactionsServer{
		storage: storage,
	}
}

func (s *EthTransactionsServer) Start(port int) error {
	panic("not implemented")
	return nil
}

func (s *EthTransactionsServer) Stop(ctx context.Context) error {
	panic("not implemented")
	return nil
}
