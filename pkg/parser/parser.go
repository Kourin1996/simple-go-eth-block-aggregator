package parser

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/types"
	"log"
	"math"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxRetry                        = 10
	DefaultFetchTimeout             = 10 * time.Second
	DefaultBackoffTime              = 1 * time.Second
	DefaultNextBlockPollingInterval = 10 * time.Second
)

type Parser struct {
	ethClient EthClient
	storage   EthTransactionStorage

	addressMap         *sync.Map
	currentBlockHeight *atomic.Uint64

	blockCh chan *types.Block

	// notification
	notifyErrCh        chan error
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

		addressMap:         &sync.Map{},
		currentBlockHeight: &atomic.Uint64{},

		blockCh:            make(chan *types.Block, 1),
		notifyErrCh:        make(chan error, 1),
		notifyCloseCh:      make(chan struct{}),
		notifyTerminatedCh: make(chan struct{}),
	}
}

// ErrCh returns channel of error which is sent from goroutine
func (p *Parser) ErrCh() <-chan error {
	return p.notifyErrCh
}

// GetCurrentBlock returns last parsed block
func (p *Parser) GetCurrentBlock() int {
	h := p.currentBlockHeight.Load()

	return int(h)
}

// Subscribe adds address to observer
func (p *Parser) Subscribe(address string) bool {
	address = strings.ToLower(address)

	_, subscribed := p.addressMap.LoadOrStore(address, true)

	return !subscribed
}

// GetTransactions returns list of inbound or outbound transactions for an address
func (p *Parser) GetTransactions(address string) []types.Transaction {
	return p.storage.GetTransactionsByAddress(address)
}

// Start prepares required parameters and start background jobs
func (p *Parser) Start(beginningHeight *big.Int) error {
	if beginningHeight == nil {
		height, err := p.fetchLatestHeight()
		if err != nil {
			return err
		}

		beginningHeight = height
	}

	log.Printf("start fetching blocks from %d", beginningHeight.Uint64())

	go p.runScrapingProcess(*beginningHeight)
	go p.runStoringProcess()

	return nil
}

// Stop tries to terminate background job
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

// runScrapingProcess is a background job to fetch block in order and send it to channel
func (p *Parser) runScrapingProcess(beginningHeight big.Int) {
	current := &beginningHeight

	defer func() {
		log.Printf("scrapingProcess has been finished")
		close(p.notifyTerminatedCh)
	}()

	for {
		// fetch block
		block, err := p.fetchBlock(*current)
		if errors.Is(err, context.Canceled) {
			// Stop has been called, terminate process
			return
		} else if err != nil {
			// unrecoverable error occurred
			p.notifyErrCh <- err

			return
		}

		// next block is not created yet, wait certain time and retry
		if block == nil {
			log.Printf("next block is not created yet, retry in %d seconds...", uint(DefaultNextBlockPollingInterval.Seconds()))

			select {
			case <-time.After(DefaultNextBlockPollingInterval):
				continue
			case <-p.notifyCloseCh:
				// Stop has been called, terminate process
				return
			}
		}

		log.Printf("fetched new block, height=%d", current.Uint64())

		select {
		case <-p.notifyCloseCh:
			// Stop has been called, terminate process
			return
		case p.blockCh <- block:
			// Emit the fetched block
		}

		// increment next block height
		current = current.Add(current, big.NewInt(1))
	}
}

// runStoringProcess process fetched block and save transactions to storage
func (p *Parser) runStoringProcess() {
	for {
		// wait for new incoming block
		var block *types.Block
		select {
		case <-p.notifyCloseCh:
			// Stop has been called, terminate process
			return
		case block = <-p.blockCh:
		}

		// filter transactions by address
		filtered := make([]*types.Transaction, 0, len(block.Transactions))
		for _, tx := range block.Transactions {
			tx := tx

			if p.isSubscribingTo(tx.From) || p.isSubscribingTo(tx.To) {
				log.Printf("found a concerned transaction, hash=%s, from=%s, to=%s", tx.Hash, tx.From, tx.To)
				filtered = append(filtered, &tx)
			}
		}

		// insert transactions into storage
		if err := p.storage.InsertTransactions(filtered); err != nil {
			log.Printf("failed to save transactions to storage: %v", err)
			p.notifyErrCh <- err
		}

		// update current height
		if err := p.updateCurrentHeight(block.Number); err != nil {
			log.Printf("failed to store current block height: %v", err)
			p.notifyErrCh <- err
		}

		log.Printf("saved transactions of block, block height=%d", p.currentBlockHeight.Load())
	}
}

// fetchLatestHeight is a wrapper function to fetch latest block height by client
func (p *Parser) fetchLatestHeight() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultFetchTimeout)
	defer cancel()

	res, err := p.ethClient.GetBlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest block height: %w", err)
	}

	return res, nil
}

// fetchBlock fetches a block by given height with retry and backoff mechanisms
// It keeps attempting until it either succeeds, exceeds the maximum number of retries, or is cancelled
func (p *Parser) fetchBlock(height big.Int) (*types.Block, error) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-p.notifyCloseCh
		cancel()
	}()

	retryTime := 0 // number of attempt

	for {
		ctx, _ := context.WithTimeout(ctx, DefaultFetchTimeout) // nolint:govet

		block, err := p.ethClient.GetBlockByNumber(ctx, height, true)
		if err == nil {
			return block, nil
		}

		// error handling
		// cancelled by outside, exit function
		if errors.Is(err, context.Canceled) {
			return nil, err
		}

		// return error if retry times exceeds threshold, otherwise go to next loop for retry
		retryTime++
		if retryTime >= MaxRetry {
			return nil, fmt.Errorf("failed to acquire block after %d attempts: %w", MaxRetry, err)
		}

		// exponential backoff
		multiplier := math.Pow(2, float64(retryTime-1))
		delay := time.Duration(multiplier) * DefaultBackoffTime

		log.Printf("failed to fetch block, retry in %d seconds", uint(delay.Seconds()))

		select {
		case <-time.After(delay):
		case <-p.notifyCloseCh:
			return nil, nil
		}
	}
}

// isSubscribingTo is a helper function to read given address from map
func (p *Parser) isSubscribingTo(address string) bool {
	// just check the target address is stored or not
	_, existing := p.addressMap.Load(strings.ToLower(address))

	return existing
}

// updateCurrentHeight updates current maximum fetched block height
func (p *Parser) updateCurrentHeight(blockHeightHex string) error {
	height, ok := (&big.Int{}).SetString(blockHeightHex, 0)
	if !ok {
		return fmt.Errorf("failed to parse block height, %s", blockHeightHex)
	}

	p.currentBlockHeight.Store(height.Uint64())

	return nil
}
