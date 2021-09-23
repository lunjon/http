package command

import (
	"io"
	"net/http"
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
