package rest

import (
	"bytes"
	"encoding/json"
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

func (res *Result) BodyFormatString() (string, error) {
	if res.body != nil {
		return "", nil
	}

	b, err := ioutil.ReadAll(res.response.Body)
	defer res.response.Body.Close()
	if err != nil {
		return "nil", err
	}

	dst := &bytes.Buffer{}
	err = json.Indent(dst, b, "", "  ")
	if err != nil {
		return "", err
	}

	return dst.String(), nil
}
