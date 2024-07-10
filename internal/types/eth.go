package types

type Block struct {
	BaseFeePerGas         string        `json:"baseFeePerGas"`
	BlobGasUsed           string        `json:"blobGasUsed"`
	Difficulty            string        `json:"difficulty"`
	ExcessBlobGas         string        `json:"excessBlobGas"`
	ExtraData             string        `json:"extraData"`
	GasLimit              string        `json:"gasLimit"`
	GasUsed               string        `json:"gasUsed"`
	Hash                  string        `json:"hash"`
	LogsBloom             string        `json:"logsBloom"`
	Miner                 string        `json:"miner"`
	MixHash               string        `json:"mixHash"`
	Nonce                 string        `json:"nonce"`
	Number                string        `json:"number"`
	ParentBeaconBlockRoot string        `json:"parentBeaconBlockRoot"`
	ParentHash            string        `json:"parentHash"`
	ReceiptsRoot          string        `json:"receiptsRoot"`
	Sha3Uncles            string        `json:"sha3Uncles"`
	Size                  string        `json:"size"`
	StateRoot             string        `json:"stateRoot"`
	Timestamp             string        `json:"timestamp"`
	TotalDifficulty       string        `json:"totalDifficulty"`
	Transactions          []Transaction `json:"transactions"`
	TransactionsRoot      string        `json:"transactionsRoot"`
	Uncles                []interface{} `json:"uncles"`
	Withdrawals           []Withdrawal  `json:"withdrawals"`
	WithdrawalsRoot       string        `json:"withdrawalsRoot"`
}

type Withdrawal struct {
	Index          string `json:"index"`
	ValidatorIndex string `json:"validatorIndex"`
	Address        string `json:"address"`
	Amount         string `json:"amount"`
}

type Transaction struct {
	BlockHash            string              `json:"blockHash"`
	BlockNumber          string              `json:"blockNumber"`
	From                 string              `json:"from"`
	Gas                  string              `json:"gas"`
	GasPrice             string              `json:"gasPrice"`
	MaxFeePerGas         string              `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string              `json:"maxPriorityFeePerGas,omitempty"`
	Hash                 string              `json:"hash"`
	Input                string              `json:"input"`
	Nonce                string              `json:"nonce"`
	To                   string              `json:"to"`
	TransactionIndex     string              `json:"transactionIndex"`
	Value                string              `json:"value"`
	Type                 string              `json:"type"`
	AccessList           []AccessListElement `json:"accessList,omitempty"`
	ChainId              string              `json:"chainId,omitempty"`
	V                    string              `json:"v"`
	R                    string              `json:"r"`
	S                    string              `json:"s"`
	YParity              string              `json:"yParity,omitempty"`
	MaxFeePerBlobGas     string              `json:"maxFeePerBlobGas,omitempty"`
	BlobVersionedHashes  []string            `json:"blobVersionedHashes,omitempty"`
}

type AccessListElement struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}
