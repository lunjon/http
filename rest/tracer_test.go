package rest

import (
	"crypto/tls"
	"net/http/httptrace"
	"testing"

	"github.com/lunjon/http/logging"
	"github.com/stretchr/testify/require"
)

func TestTracer(t *testing.T) {
	logger := logging.NewLogger()
	client := NewClient(server.Client(), logger, logger)

	url, _ := ParseURL(server.URL, nil)
	req, err := client.BuildRequest("GET", url, nil, nil)
	require.NoError(t, err)

	res := client.SendRequest(req)
	require.Nil(t, res.Error())
}

func TestTracerDNS(t *testing.T) {
	logger := logging.NewLogger()
	tracer := newTracer(logger)

	tracer.DNSStart(httptrace.DNSStartInfo{})
	tracer.DNSDone(httptrace.DNSDoneInfo{})

	require.NotZero(t, tracer.dnsStart)
	require.NotZero(t, tracer.dnsDuration)
}

func TestTracerTLS(t *testing.T) {
	logger := logging.NewLogger()
	tracer := newTracer(logger)

	tracer.TLSHandshakeStart()
	require.NotZero(t, tracer.tlsStart)

	tracer.TLSHandshakeDone(tls.ConnectionState{}, nil)
}
