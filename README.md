# HTTP from TCP

A custom HTTP server implemented from scratch in Go using raw TCP sockets. This project demonstrates how HTTP works at the protocol level by implementing request parsing and response writing without using the standard library's `net/http` package.

## Overview

This project builds a fully functional HTTP/1.1 server by:
- Parsing HTTP request lines, headers, and bodies from raw TCP connections
- Writing HTTP responses including status lines, headers, and bodies
- Supporting chunked transfer encoding
- Handling multiple concurrent connections

## Project Structure

```
.
├── cmd/httpserver/           # Application entry point
│   └── main.go
├── internal/
│   ├── server/               # TCP server and request routing
│   │   └── server.go
│   ├── request/              # HTTP request parsing
│   │   └── request.go
│   ├── response/             # HTTP response writing
│   │   └── response.go
│   └── headers/              # HTTP header management
│       └── headers.go
├── go.mod
└── go.sum
```

## Features

- **Custom HTTP Parsing**: Implements request line, header, and body parsing from raw bytes
- **State Machine**: Uses state machines for both request parsing and response writing
- **HTTP/1.1 Support**: Supports HTTP/1.1 protocol version
- **Chunked Transfer Encoding**: Implements chunked responses for streaming data
- **Concurrent Connections**: Handles multiple clients simultaneously using goroutines

## Endpoints

| Path      | Description                     |
|-----------|---------------------------------|
| `/`       | Welcome message                 |
| `/about`  | About the server                |
| `/echo`   | Echoes the request method       |
| `/health` | Health check endpoint (returns OK) |

## Running the Server

```bash
go run cmd/httpserver/main.go
```

The server starts on port `42069` by default. Press `Ctrl+C` to stop gracefully.

## Testing

Run the included tests:

```bash
go test ./...
```

Or test specific packages:

```bash
go test ./internal/request/
go test ./internal/headers/
```

## Testing with curl

Once the server is running:

```bash
# Root endpoint
curl -v http://localhost:42069/

# About endpoint
curl -v http://localhost:42069/about

# Echo endpoint
curl -v http://localhost:42069/echo

# Health check
curl -v http://localhost:42069/health
```

## Technical Details

- **Request Parsing**: Uses a state machine (`StateParserRequestLine` → `StateParserHeader` → `StateParserBody` → `StateParserDone`) to incrementally parse incoming requests
- **Response Writing**: Follows a similar state machine pattern for writing responses in the correct order
- **Connection Handling**: Each connection is handled in a separate goroutine for concurrent request processing
- **Error Handling**: Returns appropriate HTTP status codes (400 for bad requests, 404 for not found)

## Dependencies

- Go 1.25+
- stretchr/testify (for testing)