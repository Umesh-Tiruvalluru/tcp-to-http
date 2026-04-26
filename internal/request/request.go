package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type parserState string


type RequestLine struct {
	Method        string
	RequestTarget string // path
	HTTPVersion   string
}


type Request struct {
	RequestLine RequestLine
	state       parserState
	// Headers map[string]string
	// Body []byte
}

var (
	ErrorMalformedRequestLine = fmt.Errorf("malformed request line")
	ErrorUnsupportedHTTPVersion = fmt.Errorf("unsupported http version")
)

const (
	StateInit parserState = "Init"
	StateDone parserState = "Done"
)

var CRLF = []byte("\r\n")

// ParseRequestLine takes a string representing the request line and parses it into a RequestLine struct.
func ParseRequestLine(data []byte) (RequestLine, int, error) {
	idx := bytes.Index(data, CRLF)
	if idx == -1 {
		return  RequestLine{}, 0, nil
	}

	requestLine := string(data[:idx])
	consumedBytes := idx+len(CRLF)

	parts := strings.Split(requestLine, " ")
	
	if len(parts) != 3 {
		return  RequestLine{}, consumedBytes, ErrorMalformedRequestLine
	}

	httpParts := strings.Split(parts[2], "/")
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return RequestLine{}, consumedBytes, ErrorUnsupportedHTTPVersion
	}

	return RequestLine{
		Method: parts[0],
		RequestTarget: parts[1],
		HTTPVersion: httpParts[1],
	}, consumedBytes, nil
}


func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case StateInit:
		requestLine, consumedBytes, err := ParseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if requestLine.Method == "" {
			return consumedBytes, nil
		}
		r.RequestLine = requestLine
		r.state = StateDone
		return consumedBytes, nil
	default:
		panic("invalid state")
	}
}

// RequestFromReader reads from the provided reader and constructs a Request struct.
func RequestFromReader(reader io.Reader) (*Request, error) {

	request := &Request{
		state: StateInit,
	}


	buf := make([]byte, 1024) // buffer to hold incoming data
	bufLen := 0 // number of bytes currently in the buffer

	for request.state != StateDone{
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("unexpected EOF")
			}
			return nil, err
		}

		if n > 0 {
			bufLen += n
			_, err := request.parse(buf[:bufLen])
			if err != nil {
				return nil, err
			}
		}

	}

	return request, nil
}



/* 
- Read from reader until we have a full request line (ending with CRLF)
	- function ReadFromReader (reader buf.Reader) (*Request, error) 
	- function ParseRequestLine (line string) (RequestLine, error)
	- function (r *Request) parse() string
*/
