package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// call sends JSON-RPC request to server and returns response
func (c *EthJsonRpcClient) call(ctx context.Context, request *JsonRpcRequest) (*JsonRpcResponse, error) {
	// serialize request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize JSON RPC request: %w", err)
	}

	// build request
	req, err := http.NewRequest("POST", c.jsonRpcUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new JSON RPC request: %w", err)
	}

	// send request
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call JSON RPC: %w", err)
	}

	defer resp.Body.Close()

	// server should return Ok
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rpc server returns not 200 status: %d", resp.StatusCode)
	}

	// parse response json
	result := &JsonRpcResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON RPC response: %w", err)
	}

	return result, nil
}
