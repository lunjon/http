package util

import (
	"fmt"
	"regexp"
	"strings"
)

// Map slice a using f.
func Map[T any](a []T, f func(T) T) []T {
	if len(a) == 0 {
		return []T{}
	}

	s := make([]T, len(a))
	for i, v := range a {
		s[i] = f(v)
	}
	return s
}

// Filter returns all values of a for which f returns true.
func Filter(a []string, f func(string) bool) []string {
	if len(a) == 0 {
		return []string{}
	}

	s := make([]string, 0)
	for _, v := range a {
		if f(v) {
			s = append(s, v)
		}
	}
	return s
}

func Contains(a []string, b string) bool {
	for _, v := range a {
		if v == b {
			return true
		}
	}
	return false
}

var (
	basePattern = `^([a-zA-Z0-9\-_]+)\s*%s\s*(\S[\s\S]*)+$`
)

type Splitter struct {
	re *regexp.Regexp
}

func NewSplitter(token string) *Splitter {
	re := regexp.MustCompile(fmt.Sprintf(basePattern, token))
	return &Splitter{re}
}

func (h *Splitter) ParseMany(values []string) (map[string]string, error) {
	m := map[string]string{}
	for _, v := range values {
		key, value, err := h.Parse(v)
		if err != nil {
			return nil, err
		}
		m[key] = value
	}
	return m, nil
}

func (h *Splitter) Parse(s string) (string, string, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return "", "", fmt.Errorf("empty string")
	}

	match := h.re.FindAllStringSubmatch(s, -1)
	if match == nil {
		return "", "", fmt.Errorf("invalid key-value format: %s", s)
	}

	key := strings.TrimSpace(match[0][1])
	value := strings.TrimSpace(match[0][2])
	return key, value, nil
}
