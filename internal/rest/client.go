package rest

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/lunjon/httpreq/internal/parse"
	"net/http/httptrace"
	"strings"
)

type tracer struct {
	currentRequest *http.Request
	logger         *log.Logger
}

func newTracer(logger *log.Logger) *tracer {
	return &tracer{
		logger: logger,
	}
}

func (t *tracer) RoundTrip(req *http.Request) (*http.Response, error) {
	t.currentRequest = req
	return http.DefaultTransport.RoundTrip(req)
}

func (t *tracer) GotConn(info httptrace.GotConnInfo) {
	if info.Reused {
		t.logger.Printf("Connection reused for %s", t.currentRequest.URL.String())
	}
}

func (t *tracer) DNSStart(info httptrace.DNSStartInfo) {
	t.logger.Printf("Resolving DNS for host %s", info.Host)
}

func (t *tracer) DNSDone(info httptrace.DNSDoneInfo) {
	if info.Err != nil {
		t.logger.Printf("Failed to during DNS lookup: %v", info.Err)
	} else {
		t.logger.Printf("DNS lookup returned address: %s (coalesced = %v)", info.Addrs, info.Coalesced)
	}
}

func (t *tracer) ConnectStart(network, addr string) {
	t.logger.Printf("Attempting to connect on %s to %s", network, addr)
}

func (t *tracer) ConnectDone(network, addr string, err error) {
	if err != nil {
		t.logger.Printf("Failed to connect on %s to %s: %v", network, addr, err)
	} else {
		t.logger.Printf("Connected done on %s to %s", network, addr)
	}
}

type Client struct {
	httpClient *http.Client
	trace      *httptrace.ClientTrace
	logger     *log.Logger
}

func NewClient(httpClient *http.Client, logger *log.Logger) *Client {
	t := newTracer(logger)
	trace := &httptrace.ClientTrace{
		GotConn:      t.GotConn,
		ConnectStart: t.ConnectStart,
		ConnectDone:  t.ConnectDone,
		DNSStart:     t.DNSStart,
		DNSDone:      t.DNSDone,
	}

	return &Client{
		httpClient: httpClient,
		trace:      trace,
		logger:     logger,
	}
}

func (client *Client) BuildRequest(method, url string, json []byte, header http.Header) (*http.Request, error) {
	client.logger.Printf("Building request: %s %s", method, url)

	u, err := parse.ParseURL(url)
	if err != nil {
		client.logger.Printf("Failed to parse url: %v", err)
		return nil, err
	}

	client.logger.Printf("Parsed URL: %v", u.String())

	var body io.Reader
	if json != nil {
		client.logger.Printf("Using request body: %s", string(json))
		body = bytes.NewReader(json)
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		client.logger.Printf("Failed to build request: %v", err)
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	return req, nil
}

func (client *Client) SignRequest(req *http.Request, body []byte, region, profile string) error {
	if region == "" {
		return fmt.Errorf("must specify an AWS region")
	}

	client.logger.Print("Signing request using Sig V4")

	var credProvider credentials.Provider
	if profile != "" {
		client.logger.Print("No AWS profile specified, trying default")
		credProvider = &credentials.SharedCredentialsProvider{
			Filename: "", // Use default, i.e. the configuration in use home directory
			Profile:  profile,
		}
	} else {
		client.logger.Print("Using AWS credentials from environment")
		credProvider = &credentials.EnvProvider{}
	}

	creds := credentials.NewCredentials(credProvider)
	signer := v4.NewSigner(creds)
	_, err := signer.Sign(req, bytes.NewReader(body), "execute-api", region, time.Now())
	return err
}

func (client *Client) SendRequest(req *http.Request) *Result {

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), client.trace))

	client.logger.Printf("Sending request: %s %s", req.Method, req.URL.String())
	if len(req.Header) > 0 {
		var b strings.Builder
		fmt.Fprintln(&b, "Request headers:")
		for name, value := range req.Header {
			fmt.Fprintf(&b, "\t%s: %s\n", name, value)
		}
		client.logger.Print(b.String())
	}

	start := time.Now()
	res, err := client.httpClient.Do(req)
	elapsed := time.Since(start)

	client.logger.Printf("Request duration: %v", elapsed)

	if err == nil && res != nil {
		var b strings.Builder
		fmt.Fprintln(&b, "Response headers:")
		for name, value := range res.Header {
			fmt.Fprintf(&b, "\t%s: %s\n", name, value)
		}
		client.logger.Print(b.String())
	}

	return &Result{
		response: res,
		elapsed:  elapsed,
		err:      err,
	}
}
