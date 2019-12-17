package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/lunjon/httpreq/pkg/parse"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client, timeout int) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}
	return &Client{httpClient: httpClient}
}

func (client *Client) BuildRequest(method, url string, json []byte, header http.Header) (*http.Request, error) {
	u, err := parse.ParseURL(url)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if json != nil {
		body = bytes.NewReader(json)
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
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

	var credProvider credentials.Provider
	if profile != "" {
		credProvider = &credentials.SharedCredentialsProvider{
			Filename: "", // Use default, i.e. the configuration in use home directory
			Profile:  profile,
		}
	} else {
		credProvider = &credentials.EnvProvider{}
	}

	creds := credentials.NewCredentials(credProvider)
	signer := v4.NewSigner(creds)
	_, err := signer.Sign(req, bytes.NewReader(body), "execute-api", region, time.Now())
	return err
}

func (client *Client) SendRequest(req *http.Request) *Result {
	start := time.Now()
	res, err := client.httpClient.Do(req)
	elapsed := time.Since(start)

	return &Result{
		response: res,
		elapsed:  elapsed,
		err:      err,
	}
}
