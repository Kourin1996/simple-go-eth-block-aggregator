package main

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/aggregator"
	"Kourin1996/simple-go-eth-block-aggregator/internal/jsonrpc"
	"Kourin1996/simple-go-eth-block-aggregator/internal/server"
	"Kourin1996/simple-go-eth-block-aggregator/internal/txstorage"
	"context"
	"errors"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	EnvKeyApiPort          = "API_PORT"
	EnvKeyBeginningHeights = "BEGINNING_HEIGHTS"
	EnvKeyJsonRpcUrl       = "JSON_RPC_URL"

	DefaultApiPort = 8000
)

func main() {
	envs, err := readEnvs()
	if err != nil {
		log.Fatalf("failed to read some envs: %+v", err)
	}

	client := &http.Client{}
	ethClient := jsonrpc.New(client, envs.JsonRpcUrl)
	store := txstorage.New()

	agg := aggregator.New(ethClient, store)
	srv := server.New(store)

	go agg.Start(envs.BeginningHeight)
	go srv.Start(envs.ApiPort)

	waitForTerminationSignal()

	if err := terminateServices([]Stoppable{agg, srv}); err != nil {
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
	panic("not implemented!")

	return &Env{}, nil
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
