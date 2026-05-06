package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Umesh-Tiruvalluru/httpfromtcp/internal/server"
)

var portNumber = uint16(42069)

func main () {
	_, err := server.ServeHTTP(portNumber)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	// defer server.Close()
	log.Println("Server started on port", portNumber)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}