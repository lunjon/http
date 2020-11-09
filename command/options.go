package command

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type HeaderOption struct {
	values http.Header
}

func NewHeaderOption() *HeaderOption {
	return &HeaderOption{
		values: make(http.Header),
	}
}

// Append adds the provided value as a header if it is valid
func (h *HeaderOption) Set(s string) error {
	key, value, err := parseHeader(s)
	if err != nil {
		return err
	}
	h.values.Add(key, value)
	return nil
}

func (h *HeaderOption) Type() string {
	return "Header"
}

func (h *HeaderOption) String() string {
	return ""
}

// Parse string s into a header name and value.
func parseHeader(h string) (string, string, error) {
	h = strings.TrimSpace(h)
	if len(h) == 0 {
		return "", "", fmt.Errorf("empty")
	}

	re := regexp.MustCompile(`^([a-zA-Z0-9\-_]+)\s*[:=]\s*(\S[\s\S]*)+$`)

	match := re.FindAllStringSubmatch(h, -1)
	if match == nil {
		return "", "", fmt.Errorf("invalid header format: %s", h)
	}

	key := strings.TrimSpace(match[0][1])
	value := strings.TrimSpace(match[0][2])
	return key, value, nil
}
