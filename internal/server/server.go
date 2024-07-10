package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type EthTransactionsServer struct {
	Parser  Parser
	Server  *http.Server
	ErrorCh chan error
}

func New(parser Parser, port uint) *EthTransactionsServer {
	handler := http.NewServeMux()

	srv := &EthTransactionsServer{
		Parser:  parser,
		Server:  &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler},
		ErrorCh: make(chan error),
	}

	handler.HandleFunc("/current", srv.handleGetCurrentBlock)
	handler.HandleFunc("/subscribe", srv.handlePostSubscribe)
	handler.HandleFunc("/transactions", srv.handlePostGetTransactions)

	return srv
}

func (s *EthTransactionsServer) Start() {
	go func() {
		log.Printf("starting web server, address=%s", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
			// maybe reach here in case that port is used by other application
			s.ErrorCh <- err
		}
	}()
}

// Stop closes server
func (s *EthTransactionsServer) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// ErrCh returns channel of error which is sent from
func (s *EthTransactionsServer) ErrCh() <-chan error {
	return s.ErrorCh
}

// handleGetCurrentBlock is a handler for GET /current
func (s *EthTransactionsServer) handleGetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	// validate request
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// get data
	height := s.Parser.GetCurrentBlock()

	log.Printf("/current is called, height=%d", height)

	// returns response
	s.writeResponse(w, &height)
}

// handlePostSubscribe is a handler for POST /subscribe
func (s *EthTransactionsServer) handlePostSubscribe(w http.ResponseWriter, r *http.Request) {
	// validate request
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse request body
	request := &PostSubscribeRequest{}
	if err := s.readRequestBody(r, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate request body
	if err := validateAddress(request.Address); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// register
	subscribed := s.Parser.Subscribe(request.Address)

	log.Printf("/subscribe is called, address=%s, subscribed=%t", request.Address, subscribed)

	// return response
	s.writeResponse(w, &PostSubscribeResponse{
		Ok: subscribed,
	})
}

// handlePostGetTransactions is a handler for POST /transactions
func (s *EthTransactionsServer) handlePostGetTransactions(w http.ResponseWriter, r *http.Request) {
	// validate request
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse request body
	request := &PostGetTransactionsRequest{}
	if err := s.readRequestBody(r, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate request body
	if err := validateAddress(request.Address); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// register
	transactions := s.Parser.GetTransactions(request.Address)

	log.Printf("/transactions is called, address=%s, num transactions=%d", request.Address, len(transactions))

	// return response
	s.writeResponse(w, &PostGetTransactionsResponse{
		Transactions: transactions,
	})
}

// readRequestBody is a helper function to read request body and map to given body object
func (s *EthTransactionsServer) readRequestBody(
	r *http.Request,
	body interface{},
) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		return err
	}

	return nil
}

// writeResponse is a helper function to write given data as a response body in json format
func (s *EthTransactionsServer) writeResponse(
	w http.ResponseWriter,
	response interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validateAddress checks that given address is correct format in Ethereum
func validateAddress(address string) error {
	if len(address) != 42 {
		return errors.New("given address is not 20 bytes hex")
	}
	if !isHex(address) {
		return errors.New("given address isn't represented in hex")
	}
	return nil
}

// isHex is a helper function to validate hex string
func isHex(s string) bool {
	hexPattern := `^0x[0-9a-fA-F]+$`
	matched, err := regexp.MatchString(hexPattern, s)
	if err != nil {
		return false
	}
	return matched
}
