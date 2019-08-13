package rest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

func BuildRequest(method, url string, json []byte, header http.Header) (*http.Request, error) {
	url, err := parseURL(url)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if json != nil {
		body = bytes.NewReader(json)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	return req, nil
}

func SignRequest(req *http.Request, body []byte, region, profile string) error {
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

type Result struct {
	response *http.Response
	elapsed  time.Duration
	err      error
	body     []byte
}

func (res *Result) Successful() bool {
	return res.response.StatusCode < 400
}

func (res *Result) HasError() bool {
	return res.err != nil
}

func (res *Result) Error() error {
	return res.err
}

func (res *Result) ElapsedMilliseconds() float64 {
	return res.elapsed.Seconds() * 1000
}

func (res *Result) Body() ([]byte, error) {
	if res.body != nil {
		return res.body, nil
	}

	b, err := ioutil.ReadAll(res.response.Body)
	defer res.response.Body.Close()
	if err != nil {
		return nil, err
	}

	res.body = b
	return b, nil
}

func (res *Result) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintln(res.response.Request.Method, "\t", res.response.Request.URL.String()))
	builder.WriteString(fmt.Sprintln("Status", "\t", res.response.Status))
	builder.WriteString(fmt.Sprintf("Elapsed  %.02f ms", res.ElapsedMilliseconds()))
	return builder.String()
}

func SendRequest(req *http.Request) *Result {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	start := time.Now()
	res, err := client.Do(req)
	elapsed := time.Since(start)

	return &Result{
		response: res,
		elapsed:  elapsed,
		err:      err,
	}
}
