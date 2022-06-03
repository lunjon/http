package client

import (
	"crypto/tls"
	"net/http/httptrace"
	"testing"

	"github.com/lunjon/http/internal/logging"
	"github.com/stretchr/testify/require"
)

func TestTracer(t *testing.T) {
	logger := logging.NewLogger()
	client, _ := NewClient(NewSettings(), logger, logger)

	url, _ := ParseURL(server.URL, nil)
	req, err := client.BuildRequest("GET", url, nil, nil)
	require.NoError(t, err)

	_, err = client.Send(req)
	require.NoError(t, err)
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
