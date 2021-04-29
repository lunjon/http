package rest_test

import (
	"fmt"
	"testing"

	"github.com/lunjon/http/rest"
	"github.com/stretchr/testify/assert"
)

type urlTest struct {
	url       string
	exptected *rest.URL
	str       string
}

func TestParseURL_Valid(t *testing.T) {
	tests := []urlTest{
		{":9999/path", &rest.URL{
			Scheme: rest.HTTP,
			Port:   9999,
			Host:   "localhost",
			Path:   "/path",
		}, "http://localhost:9999/path"},
		{"localhost/path", &rest.URL{
			Scheme: rest.HTTP,
			Port:   80,
			Host:   "localhost",
			Path:   "/path",
		}, "http://localhost/path"},
		{"127.0.0.1/path", &rest.URL{
			Scheme: rest.HTTP,
			Port:   80,
			Host:   "127.0.0.1",
			Path:   "/path",
			Query:  "",
		}, "http://127.0.0.1/path"},
		{"https://127.0.0.1/path?query=value", &rest.URL{
			Scheme: rest.HTTPS,
			Port:   443,
			Host:   "127.0.0.1",
			Path:   "/path",
			Query:  "query=value",
		}, "https://127.0.0.1/path?query=value"},
		{"http://localhost", &rest.URL{
			Scheme: rest.HTTP,
			Port:   80,
			Host:   "localhost",
			Path:   "",
		}, "http://localhost"},
		{"http://localhost/path", &rest.URL{
			Scheme: rest.HTTP,
			Port:   80,
			Host:   "localhost",
			Path:   "/path",
		}, "http://localhost/path"},
		{"https://localhost/path", &rest.URL{
			Scheme: rest.HTTPS,
			Port:   443,
			Host:   "localhost",
			Path:   "/path",
		}, "https://localhost/path"},
		{"http://127.0.0.1:50126/path", &rest.URL{
			Scheme: rest.HTTP,
			Port:   50126,
			Host:   "127.0.0.1",
			Path:   "/path",
		}, "http://127.0.0.1:50126/path"},
		{"http://127.0.0.1:50126/path?query=value", &rest.URL{
			Scheme: rest.HTTP,
			Port:   50126,
			Host:   "127.0.0.1",
			Path:   "/path",
			Query:  "query=value",
		}, "http://127.0.0.1:50126/path?query=value"},
		{"http://api.host:5000?query=value", &rest.URL{
			Scheme: rest.HTTP,
			Port:   5000,
			Host:   "api.host",
			Path:   "",
			Query:  "query=value",
		}, "http://api.host:5000?query=value"},
		{"api.host:5000?query=value", &rest.URL{
			Scheme: rest.HTTPS,
			Port:   5000,
			Host:   "api.host",
			Path:   "",
			Query:  "query=value",
		}, "https://api.host:5000?query=value"},
		{"https://api.com:5000/external/route", &rest.URL{
			Scheme: rest.HTTPS,
			Port:   5000,
			Host:   "api.com",
			Path:   "/external/route",
			Query:  "",
		}, "https://api.com:5000/external/route"},
	}

	aliases := make(map[string]string)
	for i, tt := range tests {
		name := fmt.Sprintf("%d) ParseURL(%s)", i, tt.url)
		t.Run(name, func(t *testing.T) {
			url, err := rest.ParseURL(tt.url, aliases)
			assert.NoError(t, err)
			assert.Equal(t, tt.exptected.Scheme, url.Scheme, "invalid scheme")
			assert.Equal(t, tt.exptected.Port, url.Port, "invalid port")
			assert.Equal(t, tt.exptected.Host, url.Host, "invalid host")
			assert.Equal(t, tt.exptected.Path, url.Path, "invalid path")
			assert.Equal(t, tt.exptected.Query, url.Query, "invalid query")
			assert.Equal(t, tt.str, url.String(), "invalid string representation")
		})
	}
}

func TestParseURL_Invalid(t *testing.T) {
	tests := []string{
		"",
		"\n",
		"http://",
		"https://",
		"/path",
	}

	for i, tt := range tests {
		name := fmt.Sprintf("%d) ParseURL(%s)", i, tt)
		t.Run(name, func(t *testing.T) {
			url, err := rest.ParseURL(tt, nil)
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
			url, err := rest.ParseURL(tt.url, tt.aliases)
			assert.NoError(t, err)
			assert.Equal(t, url.String(), tt.expected)
		})
	}
}
