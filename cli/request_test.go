package cli

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/lunjon/http/client"
	"github.com/lunjon/http/logging"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type state struct {
	failCalled bool
}

type fixture struct {
	handler *RequestHandler
	root    *cobra.Command
	infos   *strings.Builder
	errors  *strings.Builder
	state   *state
}

func setupRequestTest(t *testing.T) *fixture {
	logger := logging.NewLogger()
	logger.SetOutput(io.Discard)
	c := client.NewClient(testServer.Client(), logger, logger)

	state := &state{}
	failFunc := func(int) {
		state.failCalled = true
	}

	fm := &formatterMock{}
	sm := &signerMock{}

	infos := &strings.Builder{}
	errors := &strings.Builder{}

	cfg := config{
		fail:           false,
		repeat:         1,
		defaultHeaders: "x-custom: value | authorization: bearer token",
		headerOpt:      newHeaderOption(),
		aliasFilepath:  testAliasFilepath,
		logs:           io.Discard,
		infos:          infos,
		errs:           errors,
	}

	handler := newHandler(
		c,
		newAliasManagerMock(),
		fm,
		sm,
		logger,
		failFunc,
		&cfg,
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
	fixture := setupRequestTest(t)
	fixture.handler.fail = true

	err := fixture.handler.handleRequest("get", testServer.URL+"/error", "")
	require.NoError(t, err)
	require.True(t, fixture.state.failCalled)
	require.Empty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestWithHeaders(t *testing.T) {
	fixture := setupRequestTest(t)
	os.Setenv("DEFAULT_HEADERS", "x-custom: value | authorization: bearer token")

	err := fixture.handler.handleRequest("get", testServer.URL, "")
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
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

func TestGetDefaultHeaders(t *testing.T) {
	fixture := setupRequestTest(t)
	header, err := fixture.handler.getHeaders()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(header), 3)
	require.Contains(t, header, "X-Custom")
	require.Contains(t, header, userAgentHeader)
}
