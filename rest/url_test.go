package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL_Valid(t *testing.T) {
	tests := []struct {
		url       string
		exptected *URL
		str       string
	}{
		{"http://localhost/path", &URL{
			Scheme: HTTP,
			Port:   80,
			Host:   "localhost",
			Path:   "/path",
		}, "http://localhost/path"},
		{"http://127.0.0.1:50126/path", &URL{
			Scheme: HTTP,
			Port:   50126,
			Host:   "127.0.0.1",
			Path:   "/path",
		}, "http://127.0.0.1:50126/path"},
		{"http://127.0.0.1:50126/path?query=value", &URL{
			Scheme: HTTP,
			Port:   50126,
			Host:   "127.0.0.1",
			Path:   "/path",
			Query:  "query=value",
		}, "http://127.0.0.1:50126/path?query=value"},
		{"http://api.host:5000?query=value", &URL{
			Scheme: HTTP,
			Port:   5000,
			Host:   "api.host",
			Path:   "",
			Query:  "query=value",
		}, "http://api.host:5000?query=value"},
		{"https://api.com:5000/external/route", &URL{
			Scheme: HTTPS,
			Port:   5000,
			Host:   "api.com",
			Path:   "/external/route",
			Query:  "",
		}, "https://api.com:5000/external/route"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			url, err := ParseURL(tt.url)
			assert.NoError(t, err)
			assert.Equal(t, tt.exptected.Scheme, url.Scheme)
			assert.Equal(t, tt.exptected.Port, url.Port)
			assert.Equal(t, tt.exptected.Host, url.Host)
			assert.Equal(t, tt.exptected.Path, url.Path)
			assert.Equal(t, tt.exptected.Query, url.Query)
			assert.Equal(t, tt.str, url.String())

		})
	}
}

func TestParseURL_Invalid(t *testing.T) {
	tests := []string {
		"",
		"http://",
		"https://",
		"localhost",
		"/path",
		"api.com:8000/path",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			url, err := ParseURL(tt)
			assert.Error(t, err)
			assert.Nil(t, url)

		})
	}
}
