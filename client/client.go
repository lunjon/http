package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"net/http/httptrace"
	"strings"
)

var (
	supportedMethods = map[string]bool{
		http.MethodGet:    true,
		http.MethodHead:   true,
		http.MethodPost:   true,
		http.MethodPatch:  true,
		http.MethodPut:    true,
		http.MethodDelete: true,
	}
)

type Client struct {
	httpClient   *http.Client
	tracer       *Tracer
	clientTrace  *httptrace.ClientTrace
	clientLogger *log.Logger
	traceLogger  *log.Logger
}

func NewClient(
	httpClient *http.Client,
	clientLogger *log.Logger,
	traceLogger *log.Logger,
) *Client {
	t := newTracer(traceLogger)
	trace := &httptrace.ClientTrace{
		TLSHandshakeStart: t.TLSHandshakeStart,
		TLSHandshakeDone:  t.TLSHandshakeDone,
		ConnectStart:      t.ConnectStart,
		ConnectDone:       t.ConnectDone,
		DNSStart:          t.DNSStart,
		DNSDone:           t.DNSDone,
	}

	return &Client{
		httpClient:   httpClient,
		tracer:       t,
		clientLogger: clientLogger,
		clientTrace:  trace,
	}
}

func (client *Client) BuildRequest(method string, u *url.URL, body []byte, header http.Header) (*http.Request, error) {
	method = strings.ToUpper(strings.TrimSpace(method))
	supported, found := supportedMethods[method]
	if !(supported && found) {
		return nil, fmt.Errorf("invalid or unsupported method: %s", method)
	}

	client.clientLogger.Printf("Building request: %s %s", method, u.String())

	var b io.Reader
	if body != nil {
		client.clientLogger.Printf("Using request body: %s", string(body))
		b = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, u.String(), b)
	if err != nil {
		client.clientLogger.Printf("Failed to build request: %v", err)
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	return req, nil
}

func (client *Client) Send(req *http.Request) (*http.Response, error) {
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), client.clientTrace))

	client.clientLogger.Printf("Sending request: %s %s", req.Method, req.URL.String())
	if len(req.Header) > 0 {
		b := strings.Builder{}
		fmt.Fprintln(&b, "Request headers:")
		for name, value := range req.Header {
			fmt.Fprintf(&b, "  %s: %s\n", name, value)
		}
		client.clientLogger.Print(b.String())
	}

	start := time.Now()
	res, err := client.httpClient.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		client.clientLogger.Printf("Request failed: %v", err)
		return nil, err
	}

	client.clientLogger.Printf("Response status: %s", res.Status)
	client.tracer.Report(elapsed)

	if err == nil && res != nil {
		b := strings.Builder{}
		fmt.Fprintln(&b, "Response headers:")
		for name, value := range res.Header {
			fmt.Fprintf(&b, "  %s:\t%s\n", name, value)
		}
		client.clientLogger.Print(b.String())
	}

	return res, err
}
