package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h *Headers) Parse(data []byte) (n int, done bool, err error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		// not enough data to start parsing
		return 0, false, nil
	}

	if index == 0 {
		// found end of the header
		return 2, true, nil
	}

	line := strings.TrimSpace(string(data[:index]))

	colon := strings.IndexByte(line, ':')
	if colon == -1 {
		return 0, false, fmt.Errorf("invalid header format, missing ':'")
	}

	key := line[:colon]
	if key == "" {
		return 0, false, fmt.Errorf("invalid header format, empty key")
	}
	if strings.HasSuffix(key, " ") {
		return 0, false, fmt.Errorf("invalid header format, no spaces allowed before colon")
	}

	value := strings.TrimSpace(line[colon+1:])

	if *h == nil {
		*h = NewHeaders()
	}
	(*h)[key] = value

	return index + 2, false, nil
}

func NewHeaders() Headers {
	return map[string]string{}
}
