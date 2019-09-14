package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		url       string
		exptected *URL
		str       string
	}{
		{"/api/path", &URL{
			Scheme: HTTP,
			Port:   80,
			Host:   localhost,
			Path:   "/api/path",
		}, "http://localhost/api/path"},

		{":1234/api", &URL{
			Scheme: HTTP,
			Port:   1234,
			Host:   localhost,
			Path:   "/api",
		}, "http://localhost:1234/api"},

		{"localhost/path", &URL{
			Scheme: HTTP,
			Port:   80,
			Host:   localhost,
			Path:   "/path",
		}, "http://localhost/path"},

		{"localhost:1234/api/path", &URL{
			Scheme: HTTP,
			Port:   1234,
			Host:   localhost,
			Path:   "/api/path",
		}, "http://localhost:1234/api/path"},

		{"http://host.com/path/id", &URL{
			Scheme: HTTP,
			Port:   80,
			Host:   "host.com",
			Path:   "/path/id",
		}, "http://host.com/path/id"},

		{"host.com:1234/path/id", &URL{
			Scheme: HTTPS,
			Port:   1234,
			Host:   "host.com",
			Path:   "/path/id",
		}, "https://host.com:1234/path/id"},

		{"host.com/api/path", &URL{
			Scheme: HTTPS,
			Port:   443,
			Host:   "host.com",
			Path:   "/api/path",
		}, "https://host.com/api/path"},

		{"http://127.0.0.1:50126/path", &URL{
			Scheme: HTTP,
			Port:   50126,
			Host:   "127.0.0.1",
			Path:   "/path",
		}, "http://127.0.0.1:50126/path"},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			url, err := ParseURL(tt.url)
			assert.NoError(t, err)
			assert.Equal(t, tt.exptected.Scheme, url.Scheme)
			assert.Equal(t, tt.exptected.Port, url.Port)
			assert.Equal(t, tt.exptected.Host, url.Host)
			assert.Equal(t, tt.exptected.Path, url.Path)
			assert.Equal(t, tt.str, url.String())

		})
	}
}
