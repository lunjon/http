package runner_test

import (
	"os"
	"testing"

	"github.com/lunjon/httpreq/internal/runner"
	"github.com/stretchr/testify/assert"
)

func TestRequest_SetBaseURL(t *testing.T) {
	req := &runner.RequestTarget{URL: "https://api.example.com/path"}
	req.SetBaseURL("http://localhost:1234")
	assert.Equal(t, "http://localhost:1234/path", req.URL)
}

func TestSpec_Validate(t *testing.T) {
	os.Setenv("SECRET_TOKEN", "token")

	spec, err := runner.LoadSpec("testdata/spec_test.json")
	assert.NoError(t, err)

	err = spec.Validate()
	assert.NoError(t, err)

	assert.Equal(t, "token", spec.Headers["x-token"])

	req := spec.Requests[0]
	assert.Equal(t, "api.example.com/path", req.URL)
}
