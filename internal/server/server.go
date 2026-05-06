package server

import (
	"fmt"
	"log"
	"net"

	"github.com/Umesh-Tiruvalluru/httpfromtcp/internal/request"
	"github.com/Umesh-Tiruvalluru/httpfromtcp/internal/response"
)


type HandleError struct {
	StatusCode response.StatusCode
	Message    string	 
}


type Server struct {
	closed   bool
	listener net.Listener
}


func handleConnections(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		rw := &response.ResponseWriter{Conn: conn, State: response.StateWriterStatusLine}
		rw.WriteStatusLine(response.StatusBadRequest)
		rw.WriteHeaders(response.DefaultHeaders(0))
		log.Printf("Failed to parse request: %v", err)
		return
	}

	rw := &response.ResponseWriter{Conn: conn, State: response.StateWriterStatusLine}

	handler, exists := routes[req.RequestLine.RequestTarget]
	if !exists {
		rw.WriteStatusLine(response.StatusNotFound)
		rw.WriteHeaders(response.DefaultHeaders(0))
		return
	}

	handler(rw, req)
}

type Handler func(*response.ResponseWriter, *request.Request)

var routes = map[string]Handler{
	"/":          handleRoot,
	"/about":     handleAbout,
	"/echo":      handleEcho,
	"/health":    handleHealth,
}

func handleRoot(rw *response.ResponseWriter, req *request.Request) {
	body := []byte("Welcome to my custom HTTP server!\n")
	rw.WriteStatusLine(response.StatusOK)
	rw.WriteHeaders(response.DefaultHeaders(len(body)))
	rw.WriteBody(body)
}

func handleAbout(rw *response.ResponseWriter, req *request.Request) {
	body := []byte("This is a Go HTTP server built from TCP sockets.\n")
	rw.WriteStatusLine(response.StatusOK)
	rw.WriteHeaders(response.DefaultHeaders(len(body)))
	rw.WriteBody(body)
}

func handleEcho(rw *response.ResponseWriter, req *request.Request) {
	body := []byte("Echo endpoint - method: " + req.RequestLine.Method)
	rw.WriteStatusLine(response.StatusOK)
	rw.WriteHeaders(response.DefaultHeaders(len(body)))
	rw.WriteBody(body)
}

func handleHealth(rw *response.ResponseWriter, req *request.Request) {
	body := []byte("OK")
	rw.WriteStatusLine(response.StatusOK)
	rw.WriteHeaders(response.DefaultHeaders(len(body)))
	rw.WriteBody(body)
}

func runServer(s *Server, listener net.Listener) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            if s.closed {
                return
            }
            log.Printf("accept error: %v", err)
            continue
        }
        go handleConnections(conn)
    }
}

func ServeHTTP (portNumber uint16) (*Server, error) {
	server := Server{}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		log.Fatal("error", err)
	}

	defer listener.Close()

	runServer(&server, listener)

	return &server, nil
} 

func (s *Server) Close () {
	s.closed = true
	s.listener.Close()
}


