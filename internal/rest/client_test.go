package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lunjon/httpreq/internal/logging"
)

type TestServer struct{}

func (ts *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func setupClient(t *testing.T) (*Client, string) {
	logger := logging.NewLogger()
	router := &TestServer{}
	server := httptest.NewServer(router)
	client := NewClient(server.Client(), logger)

	t.Cleanup(func() {
		server.Close()
	})

	return client, server.URL
}

func TestBuildRequest(t *testing.T) {
	client, _ := setupClient(t)
	tests := []struct {
		method  string
		url     string
		body    string
		wantErr bool
	}{
		// Valid
		{"GET", "http://localhost", "", false},
		{"POST", "api.example.com:1234", "[]", false},
		{"post", "api.example.com:1234/path?query=something", `{"name": "lol"}`, false},
		{"DELETE", "https://api.example.com:1234/path?query=something", "", false},
		// Invalid
		{"", "", "", true},
		{"HEAD", "localhost/path", "", true},  // Unsupported
		{"Put", "localhost/path", "", true},   // Unsupported
		{"Patch", "localhost/path", "", true}, // Unsupported
		{"WHAT", "localhost/path", "", true},
	}

	for _, test := range tests {
		t.Run(test.method+" "+test.url, func(t *testing.T) {
			_, err := client.BuildRequest(test.method, test.url, nil, nil)
			if (err != nil) != test.wantErr {
				t.Errorf("BuildRequest() error = %v, wantErr = %v", err, test.wantErr)
				return
			}
		})
	}
}

func TestGet(t *testing.T) {
	client, url := setupClient(t)
	req, err := client.BuildRequest("GET", url, nil, nil)
	if err != nil {
		t.Errorf("failed to build: %v", err)
		return
	}

	res := client.SendRequest(req)
	if res.Error() != nil {
		t.Errorf("failed to send: %v", err)
		return
	}
	if !res.Successful() {
		t.Errorf("failed to send: %v", err)
		return
	}
}

func TestPost(t *testing.T) {
	client, url := setupClient(t)
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

			res := client.SendRequest(req)
			if res.Error() != nil {
				t.Errorf("%s failed to send: %v", name, err)
				return
			}
			if !res.Successful() {
				t.Errorf("%s failed to send: %v", name, err)
				return
			}
		})
	}
}
