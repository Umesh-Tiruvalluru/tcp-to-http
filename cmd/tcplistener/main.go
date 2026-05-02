package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Umesh-Tiruvalluru/httpfromtcp/internal/request"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("error reading request: %v", err)
		return
	}

	fmt.Println("---Request Line---")
	fmt.Printf("- Method: %s\n", req.RequestLine.Method)
	fmt.Printf("- Request Target: %s\n", req.RequestLine.RequestTarget)
	fmt.Printf("- HTTP Version: %s\n", req.RequestLine.HTTPVersion)
	fmt.Println("---Headers---")
	fmt.Println(req.Headers.Headers)
	fmt.Println("---Body---")
	if len(req.Body) == 0 {
		fmt.Println("(empty)")
	} else {
		fmt.Println(string(req.Body))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		fmt.Println("got connection from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}