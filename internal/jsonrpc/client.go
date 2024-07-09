package jsonrpc

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"math/big"
	"net/http"
)

type EthJsonRpcClient struct {
	client     *http.Client
	jsonRpcUrl string
}

func New(client *http.Client, jsonRpcUrl string) *EthJsonRpcClient {
	return &EthJsonRpcClient{
		client:     client,
		jsonRpcUrl: jsonRpcUrl,
	}
}

func (c *EthJsonRpcClient) GetLatestHeight() (*big.Int, error) {
	panic("not implemented")
	return nil, nil
}

func (c *EthJsonRpcClient) GetBlocks(from big.Int, to *big.Int) ([]types.Block, error) {
	panic("not implemented")
	return nil, nil
}
