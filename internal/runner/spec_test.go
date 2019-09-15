package runner_test

import (
	"testing"

	"github.com/lunjon/httpreq/internal/runner"
	"github.com/stretchr/testify/assert"
)

func TestSetBaseURL(t *testing.T) {
	req := &runner.RequestTarget{URL: "https://api.example.com/path"}
	req.SetBaseURL("http://localhost:1234")
	assert.Equal(t, "http://localhost:1234/path", req.URL)
}
