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

	"github.com/lunjon/http/internal/types"
)

func init() {
	supportedMethods = make(map[string]bool)
	for _, m := range SupportedMethods {
		supportedMethods[m] = true
	}
}

var (
	supportedMethods = map[string]bool{}
	SupportedMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
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
		taber := types.NewTaber("  ")
		taber.Writef("Request headers:\n")
		for name, value := range req.Header {
			line := []string{name + ":"}
			line = append(line, value...)
			taber.WriteLine(line...)
		}
		client.clientLogger.Print(taber.String())
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

	if len(res.Header) > 0 {
		taber := types.NewTaber("  ")
		taber.Writef("Response headers:\n")
		for name, value := range res.Header {
			line := []string{name + ":"}
			line = append(line, value...)
			taber.WriteLine(line...)
		}
		client.clientLogger.Print(taber.String())
	}

	return res, err
}
