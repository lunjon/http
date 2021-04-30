package rest

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"net/http/httptrace"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
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
	httpClient  *http.Client
	tracer      *Tracer
	clientTrace *httptrace.ClientTrace
	logger      *log.Logger
}

func NewClient(httpClient *http.Client, logger *log.Logger) *Client {
	t := newTracer(logger)
	trace := &httptrace.ClientTrace{
		TLSHandshakeStart: t.TLSHandshakeStart,
		TLSHandshakeDone:  t.TLSHandshakeDone,
		ConnectStart:      t.ConnectStart,
		ConnectDone:       t.ConnectDone,
		DNSStart:          t.DNSStart,
		DNSDone:           t.DNSDone,
	}

	return &Client{
		httpClient:  httpClient,
		tracer:      t,
		clientTrace: trace,
		logger:      logger,
	}
}

func (handler *Client) Timeout(timeout time.Duration) {
	handler.httpClient.Timeout = timeout
}

func (client *Client) BuildRequest(method string, url *URL, json []byte, header http.Header) (*http.Request, error) {
	method = strings.ToUpper(strings.TrimSpace(method))
	supported, found := supportedMethods[method]
	if !(supported && found) {
		return nil, fmt.Errorf("invalid or unsupported method: %s", method)
	}

	client.logger.Printf("Building request: %s %s", method, url)

	client.logger.Printf("Parsed URL: %v", url.String())

	var body io.Reader
	if json != nil {
		client.logger.Printf("Using request body: %s", string(json))
		body = bytes.NewReader(json)
	}

	req, err := http.NewRequest(method, url.String(), body)
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
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), client.clientTrace))

	client.logger.Printf("Sending request: %s %s", req.Method, req.URL.String())
	if len(req.Header) > 0 {
		b := strings.Builder{}
		fmt.Fprintln(&b, "Request headers:")
		for name, value := range req.Header {
			fmt.Fprintf(&b, "  %s: %s\n", name, value)
		}
		client.logger.Print(b.String())
	}

	start := time.Now()
	res, err := client.httpClient.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		client.logger.Printf("Request failed: %v", err)
		return &Result{elapsed: elapsed, err: err}
	}

	client.logger.Printf("Response status: %s", res.Status)
	client.tracer.Report(elapsed)

	if err == nil && res != nil {
		b := strings.Builder{}
		fmt.Fprintln(&b, "Response headers:")
		for name, value := range res.Header {
			fmt.Fprintf(&b, "  %s:\t%s\n", name, value)
		}
		client.logger.Print(b.String())
	}

	return &Result{
		response: res,
		elapsed:  elapsed,
		err:      err,
	}
}
