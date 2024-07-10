# simple-go-eth-block-aggregator

Ethereum blockchain parser that will allow to query transactions for subscribed
addresses.  
This code also has API server in order to communicate with parser easily.  

Currently, this project doesn't depend on any external packages except for standard packages

## Version

- go 1.22.0

## How to start

Make sure you have the following environment variables set before starting server

```bash
export JSON_RPC_URL=<Ethereum JSON RPC Server URL>
export API_PORT=<Port for API to serve>
export BEGINNING_HEIGHT=<Starting block to fetch (decimal or hex)>
```

```
$ make run
```

## Project Structures

```
.
├── cmd/
│   └── main.go     # Entrypoint
├── internal/
│   ├── jsonrpc     # Ethereum JSON-RPC client
│   ├── server      # API for communicating with parser
│   ├── txstorage   # Transaction storage (supports only in-memory storage for now)
│   └── types       # Common types
├── pkg/
│   └── parser      # Ethereum block & transactions collector
├── go.mod
├── go.sum
└── README.md
```

## API

This project has REST API to check how parser works

### GET /current

Returns the height of the block which Parser processed in the last.

response:
```json
14392947
```

### POST /subscribe

Subscribes to the given address. If the address is not registered yet, it returns true.

request:
```json
{
    "address": "0x65d4Ec89Ce26763B4BEa27692E5981D8CD3A58C7"
}
```

response:
```json
{
    "ok": true
}
```

### POST /transactions

Returns transactions associated with given address

request:
```json
{
    "address": "0x65d4Ec89Ce26763B4BEa27692E5981D8CD3A58C7"
}
```

response:
```json
{
    "transactions": [
        {
            "blockHash": "0x0ded9e2ee8421d6d66a5ebd90c8de2804197d259f23d879c4a9e687570106d48",
            "blockNumber": "0xdba0c5",
            "from": "0x65d4ec89ce26763b4bea27692e5981d8cd3a58c7",
            "gas": "0x5208",
            "gasPrice": "0x65d5970",
            "maxFeePerGas": "0x65fd887",
            "maxPriorityFeePerGas": "0x65c1ce0",
            "hash": "0x753c0bf41226ac0f71fbc940d1eef23a70fafb2c148bdf4d3773d897c05c141b",
            "input": "0x",
            "nonce": "0x3e",
            "to": "0x6f0609f6a920101faf5a64f6f69bdcf5d4470ec6",
            "transactionIndex": "0x5",
            "value": "0x38d7ea4c68000",
            "type": "0x2",
            "chainId": "0xaa37dc",
            "v": "0x0",
            "r": "0x9579a8c9e0fa5613775aad8cbb0bedd768e42c593ad32810229c8c9a29e96427",
            "s": "0x6b6cf5da4a54ae8aa1ede819e86b8176442c56cf2e3e1125684f8ba46812d356",
            "yParity": "0x0"
        },
        {
            "blockHash": "0xcfbd71892b65dcf0572d5b94e84de9130a7c578ff11626a739f8b44095912477",
            "blockNumber": "0xdba0df",
            "from": "0x65d4ec89ce26763b4bea27692e5981d8cd3a58c7",
            "gas": "0x5208",
            "gasPrice": "0x65d529d",
            "maxFeePerGas": "0x65fd023",
            "maxPriorityFeePerGas": "0x65c1ce0",
            "hash": "0xaa0a22fac8ee1e8b15b8f36ff87f7810f76642880663ad44fb8f5ac464a1f6b0",
            "input": "0x",
            "nonce": "0x3f",
            "to": "0x6f0609f6a920101faf5a64f6f69bdcf5d4470ec6",
            "transactionIndex": "0x4",
            "value": "0x38d7ea4c68000",
            "type": "0x2",
            "chainId": "0xaa37dc",
            "v": "0x1",
            "r": "0xdcd90094948c604a12a258ddd4bcbed04c9a4881271e1e455843e958e113c1a0",
            "s": "0x123a6500f1d44b66ff3ce0c7b4598277ffd95ee498f8f44544d82911031c25ad",
            "yParity": "0x1"
        },
        {
            "blockHash": "0x884b9ed4185fca7f3a595a28296c4a0a4281dfe4982bf1f4bd67e6120eac25b7",
            "blockNumber": "0xdba0f5",
            "from": "0x6f0609f6a920101faf5a64f6f69bdcf5d4470ec6",
            "gas": "0x5208",
            "gasPrice": "0x65d51af",
            "maxFeePerGas": "0x65fba90",
            "maxPriorityFeePerGas": "0x65c1ce0",
            "hash": "0x91cd65a05b74483b0be43b863874d41e14434e7365cf3100bf84973c31bcf71c",
            "input": "0x",
            "nonce": "0x0",
            "to": "0x65d4ec89ce26763b4bea27692e5981d8cd3a58c7",
            "transactionIndex": "0x5",
            "value": "0x38d7ea4c68000",
            "type": "0x2",
            "chainId": "0xaa37dc",
            "v": "0x0",
            "r": "0x8657d9c6996352984d91b97702b99d8b6c7ad422186901e13293a8664140ac54",
            "s": "0x2e4075edde05ace7833a29a174e9be2d503051d0bba60e3fd1eebc5aa162159c",
            "yParity": "0x0"
        }
    ]
}
```