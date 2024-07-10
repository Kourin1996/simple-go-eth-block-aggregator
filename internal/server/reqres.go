package server

import (
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"
)

// PostSubscribeRequest is a request body for POST /subscribe API
type PostSubscribeRequest struct {
	Address string `json:"address"`
}

// PostSubscribeResponse is a response body for POST /subscribe API
type PostSubscribeResponse struct {
	Ok bool `json:"ok"`
}

// PostGetTransactionsRequest is a request body for POST /transactions API
type PostGetTransactionsRequest struct {
	Address string `json:"address"`
}

// PostGetTransactionsResponse is a response body for POST /transactions API
type PostGetTransactionsResponse struct {
	Transactions []types.Transaction `json:"transactions"`
}
