package parser

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"context"
	"math/big"
)

type Parser struct {
	ethClient EthClient
	storage   EthTransactionStorage
}

func New(
	ethClient EthClient,
	storage EthTransactionStorage,
) *Parser {
	return &Parser{
		ethClient: ethClient,
		storage:   storage,
	}
}

// last parsed block
func (p *Parser) GetCurrentBlock() int {
	panic("not implemented")
}

// add address to observer
func (p *Parser) Subscribe(address string) bool {
	panic("not implemented")
}

// list of inbound or outbound transactions for an address
func (p *Parser) GetTransactions(address string) []types.Transaction {
	panic("not implemented")
}

func (p *Parser) Start(
	beginningHeight *big.Int,
) error {
	panic("not implemented!")
	return nil
}

func (p *Parser) Stop(ctx context.Context) error {
	panic("not implemented!")
	return nil
}
