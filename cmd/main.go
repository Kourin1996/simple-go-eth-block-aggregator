package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/jsonrpc"
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/server"
	"github.com/Kourin1996/simple-go-eth-block-aggregator/internal/txstorage"
	"github.com/Kourin1996/simple-go-eth-block-aggregator/pkg/parser"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	EnvKeyApiPort         = "API_PORT"
	EnvKeyBeginningHeight = "BEGINNING_HEIGHT"
	EnvKeyJsonRpcUrl      = "JSON_RPC_URL"

	DefaultApiPort uint = 8000
)

func main() {
	// read environment variables
	envs, err := readEnvs()
	if err != nil {
		log.Fatalf("failed to read some envs: %+v", err)
	}

	// create modules
	client := &http.Client{}
	ethClient := jsonrpc.New(client, envs.JsonRpcUrl)

	store := txstorage.New()
	prs := parser.New(ethClient, store)
	srv := server.New(prs, envs.ApiPort)

	// start services
	if err := prs.Start(envs.BeginningHeight); err != nil {
		log.Fatalf("failed to start parser: %v", err)
	}

	srv.Start()

	// wait until error occurs or terminate signal is sent
	waitForErrorOrTerminateSignal(prs, srv)

	// terminate services
	if err := terminateServices([]Stoppable{prs, srv}); err != nil {
		log.Fatalf("some services failed to stop by timeout, err=%+v", err)
	}

	log.Printf("all servicesc has stopped successfully, bye")
}

type Env struct {
	ApiPort         uint
	BeginningHeight *big.Int
	JsonRpcUrl      string
}

// readEnvs reads environment variables, parses, and returns Env
func readEnvs() (*Env, error) {
	var (
		port            = DefaultApiPort
		beginningHeight *big.Int
		jsonRpcUrl      string
	)

	// API port
	rawPort := os.Getenv(EnvKeyApiPort)
	if rawPort != "" {
		parsed, err := strconv.ParseUint(rawPort, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", EnvKeyApiPort, err)
		}

		port = uint(parsed)
	}

	// beginning height for block fetching
	rawBeginningHeight := os.Getenv(EnvKeyBeginningHeight)
	if rawBeginningHeight != "" {
		height, ok := (&big.Int{}).SetString(rawBeginningHeight, 0)
		if !ok {
			return nil, fmt.Errorf("failed to parse %s", EnvKeyBeginningHeight)
		}

		beginningHeight = height
	}

	// JSON RPC url
	jsonRpcUrl = os.Getenv(EnvKeyJsonRpcUrl)
	if jsonRpcUrl == "" {
		return nil, fmt.Errorf("%s is required", EnvKeyJsonRpcUrl)
	}

	return &Env{
		ApiPort:         port,
		BeginningHeight: beginningHeight,
		JsonRpcUrl:      jsonRpcUrl,
	}, nil
}

// waitForErrorOrTerminateSignal waits for SIGINT (Ctrl + c), or errors from services running as a background task
func waitForErrorOrTerminateSignal(
	p *parser.Parser,
	s *server.EthTransactionsServer,
) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("awaiting termination signals")

	select {
	case err := <-p.ErrCh():
		log.Printf("parser was terminated with error: %v", err)
	case err := <-s.ErrCh():
		log.Printf("server was terminated with error: %v", err)
	case <-signalCh:
		log.Printf("termination signal was sent")
	}
}

type Stoppable interface {
	Stop(context.Context) error
}

// terminateServices calls Stop method of each service
// and wait for them to shutdown gracefully
func terminateServices(services []Stoppable) error {
	log.Printf("terminating services...")

	num := len(services)

	var wg sync.WaitGroup
	wg.Add(num)

	errCh := make(chan error, num)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, srv := range services {
		srv := srv
		go func() {
			errCh <- srv.Stop(ctx)
			wg.Done()
		}()
	}

	wg.Wait()
	close(errCh)

	errs := make([]error, 0, num)
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
