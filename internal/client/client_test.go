package client

import (
	"fmt"
	"testing"

	"github.com/lunjon/http/internal/logging"
	"github.com/stretchr/testify/require"
)

func setupClient(t *testing.T) *Client {
	logger := logging.NewLogger()
	client := NewClient(server.Client(), logger, logger)
	return client
}

func TestBuildRequest(t *testing.T) {
	client := setupClient(t)
	tests := []struct {
		method  string
		url     string
		body    string
		wantErr bool
	}{
		// Valid
		{"GET", "http://localhost", "", false},
		{"POST", "https://api.example.com:1234", "[]", false},
		{"post", "https://api.example.com:1234/path?query=something", `{"name": "lol"}`, false},
		{"DELETE", "https://api.example.com:1234/path?query=something", "", false},
		{"HEAD", "http://localhost/path", `{}`, false},
		{"Put", "http://localhost/path", `{"name": "lol"}`, false},
		{"Patch", "http://localhost/path", `{"name": "lol"}`, false},
		// Invalid
		{"", "", "", true},
		{"WHAT", "localhost/path", "", true},
	}

	var body []byte
	for _, test := range tests {
		if test.body != "" {
			body = []byte(test.body)
		}
		t.Run(test.method+" "+test.url, func(t *testing.T) {
			url, _ := ParseURL(test.url, nil)
			_, err := client.BuildRequest(test.method, url, body, nil)
			if (err != nil) != test.wantErr {
				t.Errorf("BuildRequest() error = %v, wantErr = %v", err, test.wantErr)
				return
			}
		})
	}
}

func TestClientGet(t *testing.T) {
	client := setupClient(t)
	u, err := parseURL(server.URL)
	require.NoError(t, err)

	req, err := client.BuildRequest("GET", u, nil, nil)
	require.NoError(t, err)

	res, err := client.Send(req)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestClientPost(t *testing.T) {
	client := setupClient(t)
	url, err := parseURL(server.URL)
	require.NoError(t, err)

	tests := []string{
		"{}",
		`{"name": "test"}`,
		`{"array": [1,2,3,4]}`,
		`{"array": [1,2,3,4], "bool": true}`,
	}
	for _, body := range tests {
		name := fmt.Sprintf("POST %s", url)

		t.Run(name, func(t *testing.T) {
			req, err := client.BuildRequest("POST", url, []byte(body), nil)
			if err != nil {
				t.Errorf("%s failed to build: %v", name, err)
				return
			}

			res, err := client.Send(req)
			require.NoError(t, err)
			require.NotNil(t, res)
		})
	}
}
