# simple-go-eth-block-aggregator

An simple server for Ethereum transactions

## Version

- go 1.22.0

## How to run

```
$ go run ./cmd/*
```

## Structures

```
.
├── cmd/
│   └── main.go # Entrypoint
├── internal/
│   ├── aggregator  # Ethereum block & transactions collector
│   ├── jsonrpc     # Ethereum JSON-RPC client
│   ├── server      # API for retuning transactions
│   └── storage     # Ethereum transactions storage (in-memory for now)
├── go.mod
├── go.sum
└── README.md
```

