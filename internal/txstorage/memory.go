package txstorage

import (
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"sync"
)

type InMemoryTransactionStorage struct {
	txMap             map[string]*types.Transaction
	txHashesByAddress map[string][]string // Address -> []BlockHash

	mutex sync.RWMutex
}

func New() *InMemoryTransactionStorage {
	return &InMemoryTransactionStorage{
		txMap:             make(map[string]*types.Transaction),
		txHashesByAddress: make(map[string][]string),
	}
}

func (s *InMemoryTransactionStorage) InsertTransactions(txs []*types.Transaction) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, tx := range txs {
		s.txMap[tx.Hash] = tx
		s.appendsTxHashForAddress(tx.From, tx.Hash)
		s.appendsTxHashForAddress(tx.To, tx.Hash)
	}

	return nil
}

// GetTransactionsByAddress returns list of transactions associated with given address
func (s *InMemoryTransactionStorage) GetTransactionsByAddress(target types.Address) []types.Transaction {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	blockHashes := s.txHashesByAddress[target]
	if len(blockHashes) == 0 {
		return nil
	}

	txs := make([]types.Transaction, len(blockHashes))
	for idx, hash := range blockHashes {
		txs[idx] = *s.txMap[hash]
	}

	return txs
}

func (s *InMemoryTransactionStorage) appendsTxHashForAddress(account string, txHash string) {
	_, ok := s.txHashesByAddress[account]
	if !ok {
		s.txHashesByAddress[account] = make([]string, 0)
	}

	s.txHashesByAddress[account] = append(s.txHashesByAddress[account], txHash)
}
