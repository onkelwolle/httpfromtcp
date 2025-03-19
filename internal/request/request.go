package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	crlf = "\r\n"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("error on reading: %s\n", err)
		return nil, err
	}

	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return nil, fmt.Errorf("no CRLF in request-line")
	}

	requestLineString := string(data[:index])
	requestLine, err := requestLineFromString(requestLineString)
	if err != nil {
		return nil, err
	}
	return requestLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line (parts: %d): %s", len(parts), str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	requestTarget := parts[1]

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
