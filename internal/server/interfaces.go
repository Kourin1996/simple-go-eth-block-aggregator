package server

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/types"
)

type TransactionStorage interface {
	GetTransactionsByAddress(address types.Address) []types.Transaction
}
