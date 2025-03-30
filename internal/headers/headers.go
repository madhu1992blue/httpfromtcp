package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const headerSpecialChars = "!#$%&'*+-.^_`|~"

func validateHeaderKey(key string) error {
	for _, c := range key {
		if !(c >= 'a' && c <= 'z') &&
			!(c >= 'A' && c <= 'Z') &&
			!(c >= '0' && c <= '9') &&
			!strings.Contains(headerSpecialChars, string(c)) {

			return fmt.Errorf("invalid character in header key: %c", c)
		}
	}
	return nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	n = 0
	crlfIndex := bytes.Index(data, []byte("\r\n"))
	if crlfIndex == -1 {
		// Not enough data to parse the headers
		return 0, false, nil
	}
	if crlfIndex == 0 {
		// Empty headers, done parsing
		return 2, true, nil
	}
	headerLine := string(data[:crlfIndex])
	parts := bytes.SplitN(data[:crlfIndex], []byte(":"), 2)
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("invalid header line: %s", headerLine)
	}
	if bytes.HasSuffix(parts[0], []byte(" ")) {
		return 0, false, fmt.Errorf("invalid spacing in key: %s", headerLine)
	}
	key := strings.ToLower(string(bytes.TrimSpace(parts[0])))
	value := string(bytes.TrimSpace(parts[1]))
	if validateHeaderKey(key) != nil {
		return 0, false, fmt.Errorf("invalid header key: %s", key)
	}
	prevHeaderValue, exists := h[key]
	if exists {
		value = fmt.Sprintf("%s, %s", prevHeaderValue, value)
	}
	h[key] = value
	n = crlfIndex + 2
	return n, false, nil
}
