package server

import (
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []types.Transaction
}
