package command

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
	c := client.NewClient(server.Client(), logger, logger)

	state := &state{}
	failFunc := func() {
		state.failCalled = true
	}

	fm := &formatterMock{}
	sm := &signerMock{}

	infos := &strings.Builder{}
	errors := &strings.Builder{}

	handler := NewHandler(
		c,
		fm,
		sm,
		logger,
		infos,
		errors,
		testAliasFilepath,
		false,
		failFunc,
		1,
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

	err := fixture.handler.handleRequest("get", server.URL, "")
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestGetErrorWithFail(t *testing.T) {
	fixture := setupRequestTest(t)
	fixture.handler.fail = true

	err := fixture.handler.handleRequest("get", server.URL+"/error", "")
	require.NoError(t, err)
	require.True(t, fixture.state.failCalled)
	require.Empty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestWithHeaders(t *testing.T) {
	fixture := setupRequestTest(t)
	os.Setenv("DEFAULT_HEADERS", "x-custom: value | authorization: bearer token")

	err := fixture.handler.handleRequest("get", server.URL, "")
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
		err := fixture.handler.handleRequest(http.MethodPost, server.URL, bodyflag)
		require.NoError(t, err)
		require.NotEmpty(t, fixture.infos.String())
		require.Empty(t, fixture.errors.String())
	}
}

func TestGetDefaultHeaders(t *testing.T) {
	fixture := setupRequestTest(t)
	os.Setenv("DEFAULT_HEADERS", "x-custom: value | authorization: bearer token")
	header, err := fixture.handler.getHeaders()
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(header), 2)
}
