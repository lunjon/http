package rest

import (
	"bytes"
	"fmt"
	"log"
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

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
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
		body = bytes.NewReader(json)
	}

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
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
	log.Print("Sending request")
	start := time.Now()
	res, err := client.httpClient.Do(req)
	elapsed := time.Since(start)

	log.Printf("Request duration: %v", elapsed)

	return &Result{
		response: res,
		elapsed:  elapsed,
		err:      err,
	}
}
