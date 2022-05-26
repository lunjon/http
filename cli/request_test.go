package cli

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/logging"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type state struct {
	failCalled bool
}

type fixture struct {
	handler *RequestHandler
	root    *cobra.Command
	logs    *strings.Builder
	infos   *strings.Builder
	errors  *strings.Builder
	state   *state
}

func setupRequestTest(t *testing.T, cfgs ...config.Config) *fixture {
	logs := &strings.Builder{}
	infos := &strings.Builder{}
	errors := &strings.Builder{}

	logger := logging.New(io.Discard)
	logger.SetOutput(logs)

	c := client.NewClient(testServer.Client(), logger, logger)

	state := &state{}

	failFunc := func(int) {
		state.failCalled = true
	}

	formatter := &formatterMock{}
	signer := &signerMock{}

	cfg := config.New()
	if len(cfgs) == 1 {
		cfg = cfgs[0]
	}

	handler := newRequestHandler(
		c,
		formatter,
		signer,
		logger,
		cfg,
		http.Header{},
		infos,
		"",
		failFunc,
	)

	return &fixture{
		handler: handler,
		infos:   infos,
		errors:  errors,
		state:   state,
	}
}

func TestGet(t *testing.T) {
	fixture := setupRequestTest(t)

	err := fixture.handler.handleRequest("get", testServer.URL, "")
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestGetErrorWithFail(t *testing.T) {
	fixture := setupRequestTest(t, config.New().UseFail(true))

	err := fixture.handler.handleRequest("get", testServer.URL+"/error", "")
	require.NoError(t, err)
	require.True(t, fixture.state.failCalled)
	require.Empty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestPost(t *testing.T) {
	bodies := []string{
		"",                  // empty
		`{"body":"string"}`, // as string
		"command.go",        // filepath
	}
	fixture := setupRequestTest(t)

	for _, bodyflag := range bodies {
		err := fixture.handler.handleRequest(http.MethodPost, testServer.URL, bodyflag)
		require.NoError(t, err)
		require.NotEmpty(t, fixture.infos.String())
		require.Empty(t, fixture.errors.String())
	}
}
