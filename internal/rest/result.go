package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
		log.Printf("Failed to indent JSON: %v", err)
		return "", err
	}

	return dst.String(), nil
}

func (res *Result) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintln(res.response.Request.Method, "\t", res.response.Request.URL.String()))
	builder.WriteString(fmt.Sprintln("Status", "\t", res.response.Status))
	builder.WriteString(fmt.Sprintf("Elapsed  %.02f ms", res.ElapsedMilliseconds()))
	return builder.String()
}

func (res *Result) Info() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintln(res.response.Request.Method, "\t", res.response.Request.URL.String()))
	builder.WriteString(fmt.Sprintln("Status", "\t", res.response.Status))
	builder.WriteString(fmt.Sprintf("Elapsed  %.02f ms", res.ElapsedMilliseconds()))
	return builder.String()
}

