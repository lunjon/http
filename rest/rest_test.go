package rest

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	server *httptest.Server
)

type TestServer struct{}

func (ts *TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"result":{"data": "string"}}`))
}

func TestMain(m *testing.M) {
	router := &TestServer{}
	server = httptest.NewServer(router)

	status := m.Run()
	server.Close()

	os.Exit(status)
}
