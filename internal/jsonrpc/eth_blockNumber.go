package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
)

// GetBlockNumber queries eth_blockNumber request to JSON-RPC server
func (c *EthJsonRpcClient) GetBlockNumber(ctx context.Context) (*big.Int, error) {
	req := NewJsonRpcRequest(MethodEthBlockNumber, nil)
	res, err := c.call(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, fmt.Errorf("JSON RPC server returned an error, code=%d, message=%s", res.Error.Code, res.Error.Message)
	}

	var hexHeight string
	if err := json.Unmarshal(res.Result, &hexHeight); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json, %s: %w", string(res.Result), err)
	}

	height, ok := (&big.Int{}).SetString(hexHeight, 0)
	if !ok {
		return nil, fmt.Errorf("failed to parse block height in hex, %s", hexHeight)
	}

	return height, nil
}
