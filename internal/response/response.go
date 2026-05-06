package response

import (
	"errors"
	"fmt"
	"net"

	"github.com/Umesh-Tiruvalluru/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusNotFound            StatusCode = 404
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
	StatusMethodNotAllowed    StatusCode = 405
)

var statusText = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusNotFound:            "Not Found",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error",
	StatusMethodNotAllowed:    "Method Not Allowed",
}

const crlf = "\r\n"

type StateWriter int

const (
	StateWriterStatusLine StateWriter = iota
	StateWriterHeaders
	StateWriterBody
	StateWriterDone
	StateWriterChunked
	StateWriterChunkDone
)


type ResponseWriter struct {
	Conn  net.Conn
	State StateWriter
}

func (rw *ResponseWriter) WriteStatusLine(statusCode StatusCode) error {
    if rw.State != StateWriterStatusLine {
        return errors.New("WriteStatusLine must be called first")
    }

    text, ok := statusText[statusCode]
    if !ok {
        text = "Unknown"
    }

    line := fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode, text, crlf)
    _, err := rw.Conn.Write([]byte(line))
    if err != nil {
        return fmt.Errorf("writing status line: %w", err)
    }

    rw.State = StateWriterHeaders
    return nil
}

func (rw *ResponseWriter) WriteHeaders(h *headers.Headers) error {
    if rw.State != StateWriterHeaders {
        switch rw.State {
        case StateWriterStatusLine:
            return errors.New("WriteHeaders called before WriteStatusLine")
        case StateWriterBody, StateWriterDone:
            return errors.New("WriteHeaders called after body was written")
        }
    }

    buf := []byte{}
    for name, value := range h.Map() {
        buf = fmt.Appendf(buf, "%s: %s%s", name, value, crlf)
    }
    buf = append(buf, []byte(crlf)...) // blank line — end of headers

    _, err := rw.Conn.Write(buf)
    if err != nil {
        return fmt.Errorf("writing headers: %w", err)
    }

    rw.State = StateWriterBody
    return nil
}

func (rw *ResponseWriter) WriteBody(p []byte) (int, error) {
    if rw.State != StateWriterBody {
        switch rw.State {
        case StateWriterStatusLine:
            return 0, errors.New("WriteBody called before WriteStatusLine")
        case StateWriterHeaders:
            return 0, errors.New("WriteBody called before WriteHeaders")
        case StateWriterDone:
            return 0, errors.New("WriteBody called after response was completed")
        }
    }

    n, err := rw.Conn.Write(p)
    if err != nil {
        return n, fmt.Errorf("writing body: %w", err)
    }

    rw.State = StateWriterDone
    return n, nil
}

func (rw *ResponseWriter) StartChunked() error {
    if rw.State != StateWriterHeaders {
        return errors.New("StartChunked must be called after WriteHeaders")
    }

    h := headers.NewHeaders()
    h.Set("Transfer-Encoding", "chunked")
    h.Set("Content-Type", "text/plain")
    h.Set("Connection", "close")

    buf := []byte{}
    for name, value := range h.Map() {
        buf = fmt.Appendf(buf, "%s: %s%s", name, value, crlf)
    }
    buf = append(buf, []byte(crlf)...)

    _, err := rw.Conn.Write(buf)
    if err != nil {
        return fmt.Errorf("writing chunked headers: %w", err)
    }

    rw.State = StateWriterChunked
    return nil
}

func (rw *ResponseWriter) WriteChunk(p []byte) (int, error) {
    if rw.State != StateWriterChunked {
        return 0, errors.New("WriteChunk called before StartChunked")
    }

    chunkSize := fmt.Sprintf("%x\r\n", len(p))
    _, err := rw.Conn.Write([]byte(chunkSize))
    if err != nil {
        return 0, fmt.Errorf("writing chunk size: %w", err)
    }

    n, err := rw.Conn.Write(p)
    if err != nil {
        return n, fmt.Errorf("writing chunk data: %w", err)
    }

    _, err = rw.Conn.Write([]byte("\r\n"))
    if err != nil {
        return n, fmt.Errorf("writing chunk terminator: %w", err)
    }

    return n, nil
}

func (rw *ResponseWriter) EndChunked() error {
    if rw.State != StateWriterChunked {
        return errors.New("EndChunked called before StartChunked")
    }

    _, err := rw.Conn.Write([]byte("0\r\n\r\n"))
    if err != nil {
        return fmt.Errorf("writing chunked terminator: %w", err)
    }

    rw.State = StateWriterDone
    return nil
}

func DefaultHeaders(contentLen int) *headers.Headers {
    h := headers.NewHeaders()
    h.Set("Content-Length", fmt.Sprint(contentLen))
    h.Set("Content-Type", "text/plain")
    h.Set("Connection", "close")
    return h
}