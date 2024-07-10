package types

type Address = string // 20 bytes hex (beginning with 0x)

type Block struct {
	Height string
	Txs    []*Transaction
}

type Transaction struct {
	Hash string
	From string
	To   string
}
