package cli

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/lunjon/http/cli/options"
	"github.com/lunjon/http/internal/client"
	"github.com/lunjon/http/internal/config"
	"github.com/lunjon/http/internal/history"
	"github.com/lunjon/http/internal/logging"
	"github.com/stretchr/testify/require"
)

type testState struct {
	failCalled bool
}

type fixture struct {
	handler     *RequestHandler
	infos       *strings.Builder
	errors      *strings.Builder
	state       *testState
	historyMock history.Handler
}

func setupRequestTest(_ *testing.T, cfgs ...config.Config) *fixture {
	logs := &strings.Builder{}
	infos := &strings.Builder{}
	errors := &strings.Builder{}

	logger := logging.New(io.Discard)
	logger.SetOutput(logs)
	settings := client.NewSettings()

	c, _ := client.NewClient(settings, logger, logger)

	state := &testState{}
	failFunc := func(int) {
		state.failCalled = true
	}

	historyHandler := history.NewHandler(testHistoryPath)
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
		historyHandler,
		logger,
		cfg,
		http.Header{},
		infos,
		"",
		failFunc,
	)

	return &fixture{
		handler:     handler,
		infos:       infos,
		errors:      errors,
		state:       state,
		historyMock: historyHandler,
	}
}

func TestGet(t *testing.T) {
	fixture := setupRequestTest(t)

	err := fixture.handler.handleRequest("get", testServer.URL, options.DataOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())

	entry, err := fixture.historyMock.Latest()
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, entry.Method)
}

func TestGetErrorWithFail(t *testing.T) {
	fixture := setupRequestTest(t, config.New().UseFail(true))

	err := fixture.handler.handleRequest("get", testServer.URL+"/error", options.DataOptions{})
	require.NoError(t, err)
	require.True(t, fixture.state.failCalled)
	require.Empty(t, fixture.infos.String())
	require.Empty(t, fixture.errors.String())
}

func TestPost(t *testing.T) {
	datas := []options.DataOptions{
		options.NewDataOptions("", "", false, []string{}),
		options.NewDataOptions(`{"body":"string"}`, "", false, []string{}),
		options.NewDataOptions("", "command.go", false, []string{}),
	}
	fixture := setupRequestTest(t)

	for _, opts := range datas {
		err := fixture.handler.handleRequest(http.MethodPost, testServer.URL, opts)
		require.NoError(t, err)
		require.NotEmpty(t, fixture.infos.String())
		require.Empty(t, fixture.errors.String())
	}
}
