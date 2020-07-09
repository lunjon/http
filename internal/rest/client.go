package rest

import (
    "bytes"
    "fmt"
    "log"
    "io"
    "net/http"
    "time"

    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/signer/v4"
    "github.com/lunjon/httpreq/pkg/parse"
    "net/http/httptrace"
	"strings"
)

// transport is an http.RoundTripper that tracks in-flight events.
type transport struct {
    current *http.Request
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
    log.Print("Setting new current request")
    t.current = req
    return http.DefaultTransport.RoundTrip(req)
}

func (t *transport) GotConn(info httptrace.GotConnInfo) {
    log.Printf("Connection reused for %v: %v", t.current.URL, info.Reused)
}

type Client struct {
    httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
    return &Client{
        httpClient: httpClient,
    }
}

func (client *Client) BuildRequest(method, url string, json []byte, header http.Header) (*http.Request, error) {
    log.Printf("Building request: %s %s", method, url)

    u, err := parse.ParseURL(url)
    if err != nil {
        log.Printf("Failed to parse url: %v", err)
        return nil, err
    }

    log.Printf("Parsed URL: %v", u.String())

    var body io.Reader
    if json != nil {
		log.Printf("Using request body: %s", string(json))
        body = bytes.NewReader(json)
    }

    req, err := http.NewRequest(method, u.String(), body)
    if err != nil {
		log.Printf("Failed to build request: %v", err)
        return nil, err
    }

    if header != nil {
        log.Printf("Using HTTP header: %+v", header)
        req.Header = header
    }

    return req, nil
}

func (client *Client) SignRequest(req *http.Request, body []byte, region, profile string) error {
    if region == "" {
        return fmt.Errorf("must specify an AWS region")
    }

    log.Print("Signing request using Sig V4")

    var credProvider credentials.Provider
    if profile != "" {
        log.Print("No AWS profile specified, trying default")
        credProvider = &credentials.SharedCredentialsProvider{
            Filename: "", // Use default, i.e. the configuration in use home directory
            Profile:  profile,
        }
    } else {
        log.Print("Using AWS credentials from environment")
        credProvider = &credentials.EnvProvider{}
    }

    creds := credentials.NewCredentials(credProvider)
    signer := v4.NewSigner(creds)
    _, err := signer.Sign(req, bytes.NewReader(body), "execute-api", region, time.Now())
    return err
}

func (client *Client) SendRequest(req *http.Request) *Result {
    t := &transport{}
    trace := &httptrace.ClientTrace{
        GotConn: t.GotConn,
    }

    req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
    client.httpClient.Transport = t

    log.Printf("Sending request: %s %s", req.Method, req.URL.String())
    var b strings.Builder
    fmt.Fprintln(&b, "Request headers:")
    for name, value := range req.Header {
        fmt.Fprintf(&b, "\t%s: %s\n", name, value)
    }
    log.Print(b.String())

    start := time.Now()
    res, err := client.httpClient.Do(req)
    elapsed := time.Since(start)

    log.Printf("Request duration: %v", elapsed)

    if err == nil && res != nil {
        var b strings.Builder
        fmt.Fprintln(&b, "Response headers:")
        for name, value := range res.Header {
            fmt.Fprintf(&b, "\t%s: %s\n", name, value)
        }
        log.Print(b.String())
    }

    return &Result{
        response: res,
        elapsed:  elapsed,
        err:      err,
    }
}
