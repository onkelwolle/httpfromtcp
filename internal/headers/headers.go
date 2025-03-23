// Package headers provides functionality for parsing and managing HTTP headers.
// It implements RFC 7230 compliant header parsing with validation for header names
// and values.
package headers

import (
	"bytes"
	"fmt"
	"strings"
)

// Common errors that can occur during header parsing
var (
	ErrMissingColon     = fmt.Errorf("invalid header format, missing ':'")
	ErrEmptyKey         = fmt.Errorf("invalid header format, empty key")
	ErrSpaceBeforeColon = fmt.Errorf("invalid header format, no spaces allowed before colon")
)

// Headers represents a collection of HTTP headers where keys are case-insensitive
// and stored in lowercase form as per RFC 7230.
type Headers map[string]string

const (
	// crlf represents the HTTP line ending sequence
	crlf = "\r\n"
)

// Parse attempts to parse a single header line from the given byte slice.
// It returns:
// - n: number of bytes consumed
// - done: true if the header section is complete (empty line found)
// - err: any error that occurred during parsing
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Check for end of headers section (empty line)
	if index := bytes.Index(data, []byte(crlf)); index == -1 {
		return 0, false, nil // Not enough data
	} else if index == 0 {
		return 2, true, nil // Found end of headers
	} else {
		return h.parseHeaderLine(data, index)
	}
}

// parseHeaderLine parses a single header line and adds it to the Headers map.
// The line must end with CRLF and contain a valid key-value pair separated by ':'.
func (h Headers) parseHeaderLine(data []byte, endIndex int) (n int, done bool, err error) {
	line := strings.TrimSpace(string(data[:endIndex]))

	key, value, err := h.extractKeyValue(line)
	if err != nil {
		return 0, false, err
	}

	if err := h.validateHeaderName(key); err != nil {
		return 0, false, err
	}

	h.Set(key, value)
	return endIndex + 2, false, nil
}

// extractKeyValue splits a header line into its key and value components.
func (h Headers) extractKeyValue(line string) (key, value string, err error) {
	colonIndex := strings.IndexByte(line, ':')
	if colonIndex == -1 {
		return "", "", ErrMissingColon
	}

	key = line[:colonIndex]
	if key == "" {
		return "", "", ErrEmptyKey
	}
	if strings.HasSuffix(key, " ") {
		return "", "", ErrSpaceBeforeColon
	}

	value = strings.TrimSpace(line[colonIndex+1:])
	return key, value, nil
}

// validateHeaderName checks if the header name contains only valid token characters
// as defined in RFC 7230.
func (h Headers) validateHeaderName(name string) error {
	for _, c := range name {
		if !isValidTokenChar(c) {
			return fmt.Errorf("invalid character in header name: %c", c)
		}
	}
	return nil
}

// Set adds or updates a header value. The key is automatically converted to lowercase
// to ensure case-insensitive matching as per RFC 7230.
func (h Headers) Set(key, value string) {
	lowercaseKey := strings.ToLower(key)
	h[lowercaseKey] = value
}

// NewHeaders creates a new empty Headers map.
func NewHeaders() Headers {
	return make(map[string]string)
}

// isValidTokenChar checks if a character is valid in an HTTP header token
// as defined in RFC 7230 section 3.2.6.
func isValidTokenChar(c rune) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c >= '0' && c <= '9':
		return true
	}

	// Special characters allowed in tokens
	switch c {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}
	return false
}
