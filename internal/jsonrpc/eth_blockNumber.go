package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
)

func (c *EthJsonRpcClient) GetBlockNumber(ctx context.Context) (*big.Int, error) {
	req := NewJsonRpcRequest(MethodEthBlockNumber, nil)
	res, err := c.call(ctx, req)
	if err != nil {
		return nil, err
	}

	var hexHeight string
	if err := json.Unmarshal(res, &hexHeight); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json, %s: %w", string(res), err)
	}

	height, ok := (&big.Int{}).SetString(hexHeight, 0)
	if !ok {
		return nil, fmt.Errorf("failed to parse block height in hex, %s", hexHeight)
	}

	return height, nil
}
