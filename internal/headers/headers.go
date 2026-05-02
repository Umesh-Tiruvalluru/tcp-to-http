package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var crlf = []byte("\r\n")

type Headers struct {
	Headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		Headers: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.Headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
    key := strings.ToLower(name)
    if existing, ok := h.Headers[key]; ok {
        h.Headers[key] = existing + ", " + value
    } else {
        h.Headers[key] = value
    }
}


func isValid (name string) bool {
	for _, char := range name {
		if !(char >= 'A' && char <= 'Z') &&
			!(char >= 'a' && char <= 'z') &&
			!(char >= '0' && char <= '9') &&
			!(strings.ContainsRune("!#$%&'*+-.^_`|~", char)) {
			return false
		}
	}

	return true
}


func (h *Headers) ParseHeader(data []byte) (int, bool, error) {
	idx := bytes.Index(data, crlf)
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return len(crlf), true, nil
	}

	numByteRead := idx + len(crlf)
	fieldLine := string(data[:idx])

	name, value, ok := strings.Cut(fieldLine, ":")
	if !ok {
		return  numByteRead, false, fmt.Errorf("malformed field line")
	}

	value = strings.TrimSpace(value)

	if strings.ContainsAny(name, " \t") {
		return numByteRead, false, fmt.Errorf("Malformed field line, key must not contain any white space")
	}

	if !isValid(name) {
		return numByteRead, false, fmt.Errorf("Malformed field line, key contains invalid characters")
	}

	h.Set(name, value)
	return numByteRead, false, nil
}