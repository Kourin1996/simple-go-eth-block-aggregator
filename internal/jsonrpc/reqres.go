package jsonrpc

import "encoding/json"

const (
	DefaultJsonRpcVersion = "2.0"
)

// JsonRpcRequest is a request body for JSON RPC call
type JsonRpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func NewJsonRpcRequest(
	method string,
	params []interface{},
) *JsonRpcRequest {
	if params == nil {
		params = make([]interface{}, 0)
	}

	return &JsonRpcRequest{
		Jsonrpc: DefaultJsonRpcVersion,
		Id:      1, // no expectations of using WebSocket right now
		Method:  method,
		Params:  params,
	}
}

// JsonRpcRequest is a response body for JSON RPC call
type JsonRpcResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *JsonRpcError   `json:"error,omitempty"`
}

// JsonRpcError is an error object in JSON RPC response
type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
