package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"math/big"
)

func (c *EthJsonRpcClient) GetBlockByNumber(
	ctx context.Context,
	height big.Int,
	shouldIncludeTxs bool,
) (*types.Block, error) {
	req := NewJsonRpcRequest(MethodEthGetBlockByNumber, []interface{}{
		"0x" + height.Text(16),
		shouldIncludeTxs,
	})
	res, err := c.call(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.Error != nil {
		if res.Error.Message == "Resource not found." {
			return nil, nil
		}

		return nil, fmt.Errorf("JSON RPC server returned an error, code=%d, message=%s", res.Error.Code, res.Error.Message)
	}

	if string(res.Result) == "null" {
		return nil, nil
	}

	block := &types.Block{}
	if err := json.Unmarshal(res.Result, block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json, %s: %w", string(res.Result), err)
	}

	return block, nil
}
