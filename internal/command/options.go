package command

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type Header struct {
	values http.Header
}

func NewHeader() *Header {
	return &Header{
		values: make(http.Header),
	}
}

// Append adds the provided value as a header if it is valid
func (h *Header) Set(s string) error {
	key, value, err := parseHeader(s)
	if err != nil {
		return err
	}
	h.values.Add(key, value)
	return nil
}

func (h *Header) Type() string {
	return "Header"
}

func (h *Header) String() string {
	return "{}"
}

// Parse string s into a header name and value.
func parseHeader(h string) (string, string, error) {
	h = strings.TrimSpace(h)
	if len(h) == 0 {
		return "", "", fmt.Errorf("empty")
	}

	re := regexp.MustCompile(`^([a-zA-Z0-9\-_]+)\s*[:=]\s*(\S+)$`)

	match := re.FindAllStringSubmatch(h, -1)
	if match == nil {
		return "", "", fmt.Errorf("invalid header format: %s", h)
	}

	key := strings.TrimSpace(match[0][1])
	value := strings.TrimSpace(match[0][2])
	return key, value, nil
}
