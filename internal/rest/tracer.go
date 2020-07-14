package rest

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"net/http/httptrace"
)

type Tracer struct {
	currentRequest *http.Request
	logger         *log.Logger

	dnsStart        time.Time
	dnsDuration     time.Duration
	tlsStart        time.Time
	tlsDuration     time.Duration
	connectStart    time.Time
	connectDuration time.Duration
}

func newTracer(logger *log.Logger) *Tracer {
	return &Tracer{
		logger: logger,
	}
}

func (t *Tracer) Report(total time.Duration) {
	t.logger.Print("Request duration:")
	t.logger.Printf("\tTotal: %v", total)
	t.logger.Printf("\tDNS lookup: %v", t.dnsDuration)
	t.logger.Printf("\tTLS handshake: %v", t.tlsDuration)
	t.logger.Printf("\tConnection: %v", t.connectDuration)
}

func (t *Tracer) RoundTrip(req *http.Request) (*http.Response, error) {
	t.currentRequest = req
	return http.DefaultTransport.RoundTrip(req)
}

func (t *Tracer) TLSHandshakeStart() {
	t.tlsStart = time.Now()
	t.logger.Print("Starting TLS handshake")
}

func (t *Tracer) TLSHandshakeDone(state tls.ConnectionState, err error) {
	t.tlsDuration += time.Since(t.tlsStart)
	msg := fmt.Sprintf("TLS handshake done after %v", t.tlsDuration)
	if err != nil {
		msg += fmt.Sprintf("with error: %v", err)
	}
	t.logger.Print(msg)
}

func (t *Tracer) DNSStart(info httptrace.DNSStartInfo) {
	t.dnsStart = time.Now()
	t.logger.Printf("Resolving DNS for host %s", info.Host)
}

func (t *Tracer) DNSDone(info httptrace.DNSDoneInfo) {
	t.dnsDuration += time.Since(t.dnsStart)
	if info.Err != nil {
		t.logger.Printf("Failed to during DNS lookup: %v", info.Err)
	} else {
		t.logger.Printf(
			"DNS lookup done after %v: %s (coalesced = %v)",
			t.dnsDuration,
			info.Addrs,
			info.Coalesced)
	}
}

func (t *Tracer) ConnectStart(network, addr string) {
	t.connectStart = time.Now()
	t.logger.Printf("Attempting to connect on %s to %s", network, addr)
}

func (t *Tracer) ConnectDone(network, addr string, err error) {
	t.connectDuration += time.Since(t.connectStart)
	if err != nil {
		t.logger.Printf("Failed to connect on %s to %s: %v", network, addr, err)
	} else {
		t.logger.Printf(
			"%s connection to %s established successfully after %v",
			network,
			addr,
			t.connectDuration)
	}
}
