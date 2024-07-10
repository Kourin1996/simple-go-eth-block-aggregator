package main

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/jsonrpc"
	"Kourin1996/simple-go-eth-block-aggregator/internal/server"
	"Kourin1996/simple-go-eth-block-aggregator/internal/txstorage"
	"Kourin1996/simple-go-eth-block-aggregator/pkg/parser"
	"context"
	"errors"
	"fmt"
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

	DefaultApiPort = 8000
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
	srv := server.New(store)

	// start services
	if err := prs.Start(envs.BeginningHeight); err != nil {
		log.Fatalf("failed to start parser: %v", err)
	}
	if err := srv.Start(envs.ApiPort); err != nil {
		log.Fatalf("failed to start web server: %v", err)
	}

	// sleep until Ctrl + C is pressed
	waitForTerminationSignal()

	// terminate services
	if err := terminateServices([]Stoppable{prs, srv}); err != nil {
		log.Fatalf("some services failed to stop by timeout, err=%+v", err)
	}

	log.Printf("server stopped successfully, bye")
}

type Env struct {
	ApiPort         int
	BeginningHeight *big.Int
	JsonRpcUrl      string
}

func readEnvs() (*Env, error) {
	var (
		port            = DefaultApiPort
		beginningHeight *big.Int
		jsonRpcUrl      = ""
	)

	rawPort := os.Getenv(EnvKeyApiPort)
	if rawPort != "" {
		parsed, err := strconv.Atoi(rawPort)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", EnvKeyApiPort, err)
		}

		port = parsed
	}

	rawBeginningHeight := os.Getenv(EnvKeyBeginningHeight)
	if rawBeginningHeight != "" {
		height, ok := (&big.Int{}).SetString(rawBeginningHeight, 0)
		if !ok {
			return nil, fmt.Errorf("failed to parse %s", EnvKeyBeginningHeight)
		}

		beginningHeight = height
	}

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

func waitForTerminationSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}

type Stoppable interface {
	Stop(context.Context) error
}

func terminateServices(services []Stoppable) error {
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

	errs := make([]error, 0, num)
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
