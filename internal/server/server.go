package server

import "Kourin1996/simple-go-eth-block-aggregator/internal/types"

type EthTransactionsServer struct {}

type TransactionStorage interface {
	GetTransactions(types.Address) ([]types.Transaction, error)
}

func NewServer() *EthTransactionsServer {
	return &EthTransactionsServer{}
}

func (s *EthTransactionsServer) Start() error {
	return nil
}