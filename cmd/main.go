package main

import (
	"Kourin1996/simple-go-eth-block-aggregator/internal/server"
	"log"
)

func main() {
	s := server.NewServer()
	if err := s.Start(); err != nil {
		log.Fatalf("stopped server with error: %+v", err)
	}

	log.Printf("server stopped successfully, bye")
}