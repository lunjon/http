package client

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Result struct {
	response *http.Response
	elapsed  time.Duration
	err      error
	body     []byte
}

func (res *Result) Status() string {
	return res.response.Status
}

func (res *Result) Request() *http.Request {
	return res.response.Request
}

func (res *Result) Successful() bool {
	return res.response.StatusCode < 400
}

func (res *Result) HasError() bool {
	return res.response != nil && res.err != nil
}

func (res *Result) Error() error {
	return res.err
}

func (res *Result) ElapsedMilliseconds() float64 {
	return res.elapsed.Seconds() * 1000
}

func (res *Result) Body() ([]byte, error) {
	if res.response == nil {
		return nil, res.err
	}

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
