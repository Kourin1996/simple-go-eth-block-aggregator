package server

import "github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"

type PostSubscribeRequest struct {
	Address string `json:"address"`
}

type PostSubscribeResponse struct {
	Ok bool `json:"ok"`
}

type PostGetTransactionsRequest struct {
	Address string `json:"address"`
}

type PostGetTransactionsResponse struct {
	Transactions []types.Transaction `json:"transactions"`
}
