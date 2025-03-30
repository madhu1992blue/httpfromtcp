package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ParserState int

const (
	Initialized ParserState = iota
	Done
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	ParserState
}

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

func parseRequestLine(dataBytes []byte) (RequestLine, int, error) {
	crlfIndex := bytes.Index(dataBytes, []byte("\r\n"))
	if crlfIndex == -1 {
		// Not enough data to parse the request line
		return RequestLine{}, 0, nil
	}
	line := string(dataBytes[:crlfIndex])

	parts := strings.Fields(line)
	if len(parts) != 3 {
		return RequestLine{}, crlfIndex, fmt.Errorf("invalid request line: %s", string(line))
	}
	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return RequestLine{}, crlfIndex, fmt.Errorf("invalid method: %s", method)
		}
	}
	requestTarget := parts[1]
	httpVersion := parts[2]
	if !strings.HasPrefix(httpVersion, "HTTP/") {
		return RequestLine{}, crlfIndex, fmt.Errorf("invalid HTTP version: %s", httpVersion)
	}
	httpVersion = strings.TrimPrefix(httpVersion, "HTTP/")
	if httpVersion != "1.0" && httpVersion != "1.1" {
		return RequestLine{}, crlfIndex, fmt.Errorf("unsupported HTTP version: %s", httpVersion)
	}
	return RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}, crlfIndex, nil
}
func RequestFromReader(r io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := Request{
		ParserState: Initialized,
	}
	for req.ParserState != Done {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		index, err := r.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading from reader: %w", err)
		}
		readToIndex += index
		parsedSoFar, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing request: %w", err)
		}
		newBuf := make([]byte, len(buf)-parsedSoFar)
		copy(newBuf, buf[parsedSoFar:])
		buf = newBuf
		readToIndex -= parsedSoFar

	}
	return &req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState == Done {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}
	if r.ParserState != Initialized {
		return 0, fmt.Errorf("error: unknown parser state")
	}
	reqLine, offset, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if offset == 0 {
		// More data is needed to parse the request line
		return 0, nil
	}
	r.RequestLine = reqLine
	r.ParserState = Done
	return offset, nil
}
