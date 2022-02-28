package client

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestMIMETypeString(t *testing.T) {
	tests := []MIMEType{
		MIMETypeCSV,
		MIMETypeHTML,
		MIMETypeJSON,
		MIMETypeXML,
		MIMETypeUnknown,
	}
	for _, mime := range tests {
		t.Run(string(mime), func(t *testing.T) {
			require.NotEmpty(t, mime.String())
		})
	}
}
