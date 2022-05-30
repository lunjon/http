package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type urlTest struct {
	url  string
	port string
	str  string
}

func TestParseURL_Valid(t *testing.T) {
	tests := []urlTest{
		{":9999/path", "9999", "http://localhost:9999/path"},
		{"localhost/path", "", "http://localhost/path"},
		{"127.0.0.1/path", "", "http://127.0.0.1/path"},
		{"https://127.0.0.1/path?query=value", "", "https://127.0.0.1/path?query=value"},
		{"http://localhost", "", "http://localhost"},
		{"http://localhost/path", "", "http://localhost/path"},
		{"https://localhost/path", "", "https://localhost/path"},
		{"http://127.0.0.1:50126/path", "50126", "http://127.0.0.1:50126/path"},
		{"http://127.0.0.1:50126/path?query=value", "50126", "http://127.0.0.1:50126/path?query=value"},
		{"http://api.host:5000?query=value", "5000", "http://api.host:5000?query=value"},
		{"api.host:5000?query=value", "5000", "https://api.host:5000?query=value"},
		{"https://api.com:5000/external/route", "5000", "https://api.com:5000/external/route"},
	}

	aliases := make(map[string]string)
	for i, tt := range tests {
		name := fmt.Sprintf("%d) ParseURL(%s)", i, tt.url)
		t.Run(name, func(t *testing.T) {
			url, err := ParseURL(tt.url, aliases)
			require.NoError(t, err)
			require.NotNil(t, url)
			require.Equal(t, tt.port, url.Port())
		})
	}
}

func TestParseURL_Invalid(t *testing.T) {
	tests := []string{
		"",
		"\n",
		"http",
		"https",
		"http:",
		"https:",
		"http://",
		"https://",
		"/path",
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d ParseURL(%s)", i, tt)
		t.Run(name, func(t *testing.T) {
			url, err := ParseURL(tt, nil)
			assert.Error(t, err)
			assert.Nil(t, url)
		})
	}
}

func TestParseURL_Alias(t *testing.T) {
	tests := []struct {
		url      string
		expected string
		aliases  map[string]string
	}{
		{"{test}/api", "http://localhost/api", map[string]string{"test": "http://localhost"}},
		{"https://{a}/api", "https://localhost/api", map[string]string{"a": "localhost"}},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("%d) ParseURL(%s)", i, tt)
		t.Run(name, func(t *testing.T) {
			url, err := ParseURL(tt.url, tt.aliases)
			assert.NoError(t, err)
			assert.Equal(t, url.String(), tt.expected)
		})
	}
}
