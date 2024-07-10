package types

type Address = string // 20 bytes hex (beginning with 0x)

type Block struct {
	Txs []*Transaction
}

type Transaction struct {
	From string
	To   string
}
