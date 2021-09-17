package rest

import (
	"net/http/httptest"
	"testing"

	"github.com/lunjon/http/logging"
	"github.com/stretchr/testify/require"
)

func TestTracer(t *testing.T) {
	logger := logging.NewLogger()
	router := &TestServer{}
	server := httptest.NewTLSServer(router)
	defer server.Close()

	client := NewClient(server.Client(), logger, logger)

	url, _ := ParseURL(server.URL, nil)
	req, err := client.BuildRequest("GET", url, nil, nil)
	require.NoError(t, err)

	res := client.SendRequest(req)
	require.Nil(t, res.Error())
}
