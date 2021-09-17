package command

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/lunjon/http/logging"
	"github.com/lunjon/http/rest"
	"github.com/spf13/cobra"
)

var (
	serverURL string
)

type serverHandler struct{}

func (s *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"body": true}`))
}

func TestMain(m *testing.M) {

	server := httptest.NewServer(&serverHandler{})
	serverURL = server.URL
	status := m.Run()
	server.Close()

	os.Exit(status)
}

type fixture struct {
	handler *Handler
	root    *cobra.Command
	infos   *strings.Builder
	errors  *strings.Builder
}

func setup(t *testing.T) *fixture {
	testdir := "tmp-http-test"
	if _, err := os.Stat(testdir); os.IsNotExist(err) {
		err := os.MkdirAll(testdir, 0700)
		checkErr(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(testdir)
	})

	logger := logging.NewLogger()
	logger.SetOutput(io.Discard)
	rc := rest.NewClient(&http.Client{}, logger, logger)

	infos := &strings.Builder{}
	errors := &strings.Builder{}
	handler := NewHandler(rc, logger, logger, infos, errors, testdir)
	root := build("test", handler)
	root.SetOutput(io.Discard)

	return &fixture{
		root:    root,
		handler: handler,
		infos:   infos,
		errors:  errors,
	}
}
