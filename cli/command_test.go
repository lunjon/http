package cli

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/lunjon/http/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

var (
	testServer     *httptest.Server
	testdir        = "test-http"
	testConfigPath = path.Join(testdir, "config.toml")
)

type signerMock struct {
	called bool
}

func (f *signerMock) Sign(r *http.Request, body io.ReadSeeker) error {
	f.called = true
	return nil
}

type formatterMock struct {
	called bool
}

func (f *formatterMock) Format(r *http.Response) ([]byte, error) {
	f.called = true
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

type serverHandler struct{}

func (s *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/error":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"body": true}`))
	}
}

func TestMain(m *testing.M) {
	testServer = httptest.NewServer(&serverHandler{})
	if _, err := os.Stat(testdir); os.IsNotExist(err) {
		err := os.MkdirAll(testdir, 0700)
		checkErr(err, os.Stderr)
	}

	file, err := os.Create(testConfigPath)
	checkErr(err, os.Stderr)

	appConfig := config.New()
	err = appConfig.Write(file)
	checkErr(err, os.Stderr)

	status := m.Run()
	testServer.Close()
	os.RemoveAll(testdir)

	os.Exit(status)
}

type commandTestFixture struct {
	logs  *strings.Builder
	infos *strings.Builder
	errs  *strings.Builder
	cmd   *cobra.Command
}

func setupCommandTest(args ...string) *commandTestFixture {
	logs := &strings.Builder{}
	infos := &strings.Builder{}
	errs := &strings.Builder{}

	cliconf := cliConfig{
		configPath: testConfigPath,
		logs:       logs,
		infos:      infos,
		errors:     errs,
	}

	cmd := build("test", cliconf)
	cmd.SetArgs(args)

	return &commandTestFixture{
		logs:  logs,
		infos: infos,
		errs:  errs,
		cmd:   cmd,
	}
}

func TestDefaultBuild(t *testing.T) {
	cmd, err := Build("test")
	require.NoError(t, err)
	require.NotNil(t, cmd)
}

func TestRequestCommandGet(t *testing.T) {
	fixture := setupCommandTest("get", testServer.URL)

	err := fixture.cmd.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos)
}

func TestRequestGetSigned(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAKIAKAI")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "abcd//efgh/ijklmnopq//bca")
	os.Setenv("AWS_SESSION_TOKEN", "9bd58de0-20ab-4f29-bbd9-dedc700152e3")
	fixture := setupCommandTest("get", testServer.URL, "--aws-sigv4")

	err := fixture.cmd.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos)
	require.Empty(t, fixture.errs)
}
