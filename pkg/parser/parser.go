package parser

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"context"
	"log"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxRetry = 5
	Timeout  = 10 * time.Second
)

type Parser struct {
	ethClient EthClient
	storage   EthTransactionStorage

	addressMap         sync.Map
	currentBlockHeight *int32

	blockCh chan *types.Block

	// notification
	notifyCloseCh      chan struct{}
	notifyTerminatedCh chan struct{}
}

func New(
	ethClient EthClient,
	storage EthTransactionStorage,
) *Parser {
	return &Parser{
		ethClient: ethClient,
		storage:   storage,

		blockCh:            make(chan *types.Block, 1),
		notifyCloseCh:      make(chan struct{}),
		notifyTerminatedCh: make(chan struct{}),
	}
}

// last parsed block
func (p *Parser) GetCurrentBlock() int {
	h := atomic.LoadInt32(p.currentBlockHeight)

	return int(h)
}

// add address to observer
func (p *Parser) Subscribe(address string) bool {
	_, subscribed := p.addressMap.LoadOrStore(address, true)

	return !subscribed
}

// list of inbound or outbound transactions for an address
func (p *Parser) GetTransactions(address string) []types.Transaction {
	return p.storage.GetTransactionsByAddress(address)
}

func (p *Parser) Start(
	beginningHeight *big.Int,
) error {
	if beginningHeight == nil {
		height, err := p.ethClient.GetLatestHeight()
		if err != nil {
			return err
		}

		beginningHeight = height
	}

	go p.runScrapingProcess(*beginningHeight)
	go p.runStoringProcess()

	return nil
}

func (p *Parser) Stop(ctx context.Context) error {
	// emits close signal
	close(p.notifyCloseCh)

	// wait until background routine to be done or timeout comes
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.notifyTerminatedCh:
		return nil
	}
}

func (p *Parser) runScrapingProcess(
	beginningHeight big.Int,
) {
	current := beginningHeight

	// 2. loop until context.Done is called
	for {
		// fetch block
		block, err := p.fetchBlock(current)
		if err != nil {
			panic(err)
		}

		select {
		case <-p.notifyCloseCh:
			return
		default:
		}

		p.blockCh <- block
	}

	close(p.notifyTerminatedCh)
}

func (p *Parser) runStoringProcess() {
	for {
		// wait for new incoming block
		var block *types.Block
		select {
		case <-p.notifyCloseCh:
			return
		case block = <-p.blockCh:
		}

		// filter by address
		filtered := make([]*types.Transaction, 0, len(block.Txs))
		for _, tx := range block.Txs {
			if p.isSubscribingTo(tx.From) || p.isSubscribingTo(tx.To) {
				filtered = append(filtered, tx)
			}
		}

		// insert transactions to storage
		if err := p.storage.InsertTransactions(filtered); err != nil {
			log.Fatalf("failed to save transactions to storage: %v", err)
		}
	}
}

func (p *Parser) fetchBlock(
	height big.Int,
) (*types.Block, error) {
	ctx, cancel := context.WithCancel(context.Background())

	for {
		ctx, _ := context.WithTimeout(ctx, Timeout)

		p.ethClient.GetBlock(ctx, height)
	}

	panic("not implemented!")
}

func (p *Parser) isSubscribingTo(address string) bool {
	_, existing := p.addressMap.Load(address)

	return existing
}
