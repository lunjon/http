package rest

import (
	"testing"
)

func TestParseRoute(t *testing.T) {
	tests := []struct {
		name    string
		route   string
		want    string
		wantErr bool
	}{
		// Valid
		{"only path", "/path", "http://localhost/path", false},
		{"only path, longer", "/api/path/", "http://localhost/api/path", false},
		{"only port", ":1234", "http://localhost:1234", false},
		{"only port with path", ":1234/path/", "http://localhost:1234/path", false},
		{"missing protocol, without port", "host.com/path/", "https://host.com/path", false},
		{"missing protocol, with port", "host.com:1234/path/", "https://host.com:1234/path", false},
		{"www", "www.google.com", "https://www.google.com", false},
		{"http, without port", "http://host.com/path/", "http://host.com/path", false},
		{"http, with port", "http://host.com:1234/path", "http://host.com:1234/path", false},
		{"https, without port", "https://host.com/path/", "https://host.com/path", false},
		{"https, with port", "https://host.com:1234/path", "https://host.com:1234/path", false},
		// localhost
		{"localhost, with port, without protocol", "localhost:1234/path", "http://localhost:1234/path", false},
		{"localhost, without port, without protocol", "localhost/path", "http://localhost/path", false},
		{"localhost ip, with port", "127.0.0.1:1234/path", "http://127.0.0.1:1234/path", false},
		{"localhost ip, without port", "127.0.0.1/path", "http://127.0.0.1/path", false},
		{"localhost ip, without port, with http", "http://127.0.0.1/path", "http://127.0.0.1/path", false},
		{"localhost ip, without port, with https", "https://127.0.0.1/path", "https://127.0.0.1/path", false},
		// Invalid
		{"empty", "", "", true},
		{"whitespace", "   ", "", true},
		{"only port, missing colon", "1234/path", "", true},
		{"invalid protocol", "lol://host.com:1234/path", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseURL(tt.route)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}
