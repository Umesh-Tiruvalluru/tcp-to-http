package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Umesh-Tiruvalluru/httpfromtcp/internal/headers"
)

type StateParser int

const (
	StateParserRequestLine StateParser = iota
	StateParserHeader
	StateParserBody
	StateParserDone
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HTTPVersion   string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	Body 		[]byte

	State       StateParser
}

var (
	ErrorMalformedRequestLine   = fmt.Errorf("malformed request line")
	ErrorUnsupportedHttpVersion = fmt.Errorf("unsupported http verison, only supports 1.1 version")
	ErrUnexpectedEOF            = fmt.Errorf("unexpected end of file")
)

var crlf = []byte("\r\n")

func parseRequestLine(data []byte) (RequestLine, int, error) {
	idx := bytes.Index(data, crlf)
	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	requestLine := string(data[:idx])
	numBytesRead := idx + len(crlf)

	requestLineParts := strings.Split(requestLine, " ")
	if len(requestLineParts) != 3 {
		return RequestLine{}, numBytesRead, ErrorMalformedRequestLine
	}

	httpVersionParts := strings.Split(requestLineParts[2], "/")
	if len(httpVersionParts) != 2 || httpVersionParts[0] != "HTTP" || httpVersionParts[1] != "1.1" {
		return RequestLine{}, numBytesRead, ErrorUnsupportedHttpVersion
	}

	return RequestLine{
		Method:        requestLineParts[0],
		RequestTarget: requestLineParts[1],
		HTTPVersion:   httpVersionParts[1],
	}, numBytesRead, nil
}

func parseBody(data []byte, contentLength int) ([]byte, int, error) {
    if len(data) < contentLength {
        return nil, 0, nil  
    }

    body := make([]byte, contentLength)
    copy(body, data[:contentLength])
    return body, contentLength, nil
}

func (r *Request) ParseRequest(buf []byte) (int, error) {
	switch r.State {
	case StateParserRequestLine:
		requestLine, numBytesRead, err := parseRequestLine(buf)
		if err != nil {
			return 0, err
		}

		if numBytesRead == 0 {
			return 0, nil
		}

		r.RequestLine = requestLine
		r.State = StateParserHeader

		return numBytesRead, nil

	case StateParserHeader:
		numBytesRead, done, err := r.Headers.ParseHeader(buf)
		if err != nil {
			return 0, err
		}

		if numBytesRead == 0 {
			return 0, nil
		}

		if done {
			r.State = StateParserBody
		}

		return numBytesRead, nil

case StateParserBody:
    contentLengthStr := r.Headers.Get("Content-Length")
    if contentLengthStr == "" {
        r.State = StateParserDone
        return 0, nil
    }

    contentLength, err := strconv.Atoi(strings.TrimSpace(contentLengthStr))
    if err != nil {
        return 0, fmt.Errorf("invalid Content-Length: %w", err)
    }

    body, numBytesRead, err := parseBody(buf, contentLength)
    if err != nil {
        return 0, err
    }

    if numBytesRead == 0 {
        return 0, nil  // not enough data yet
    }

    r.Body = body
    r.State = StateParserDone
    return numBytesRead, nil
	}

	return 0, nil
}

func RequestFromReader(conn io.Reader) (*Request, error) {
    request := Request{
        State:   StateParserRequestLine,
        Headers: headers.NewHeaders(),
    }

    buf := make([]byte, 4096)
    bufLen := 0

    for request.State != StateParserDone {
        n, readErr := conn.Read(buf[bufLen:])
        bufLen += n

        for bufLen > 0 {
            numBytesRead, err := request.ParseRequest(buf[:bufLen])
            if err != nil {
                return nil, err
            }
            if numBytesRead == 0 {
                break
            }
            copy(buf, buf[numBytesRead:bufLen])
            bufLen -= numBytesRead
        }

        if readErr != nil {
            if readErr == io.EOF {
                if request.State != StateParserDone {
                    return nil, ErrUnexpectedEOF
                }
                break
            }
            return nil, readErr
        }
    }

    return &request, nil
}