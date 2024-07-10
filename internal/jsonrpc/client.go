package jsonrpc

import (
	"context"
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"
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

func (c *EthJsonRpcClient) GetLatestHeight(ctx context.Context) (*big.Int, error) {
	panic("not implemented")
	return nil, nil
}

func (c *EthJsonRpcClient) GetBlock(ctx context.Context, height big.Int) (*types.Block, error) {
	panic("not implemented")
	return nil, nil
}
