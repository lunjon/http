package command

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

var (
	server            *httptest.Server
	testdir           = "test-http"
	testAliasFilepath = path.Join(testdir, "aliases.json")
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
	server = httptest.NewServer(&serverHandler{})
	if _, err := os.Stat(testdir); os.IsNotExist(err) {
		err := os.MkdirAll(testdir, 0700)
		checkErr(err)
	}

	status := m.Run()
	server.Close()
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

	cfg := config{
		logs:          logs,
		infos:         infos,
		errs:          errs,
		aliasFilepath: testAliasFilepath,
		headerOpt:     newHeaderOption(),
	}

	cmd := build("test", &cfg)
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
	fixture := setupCommandTest("get", server.URL)

	err := fixture.cmd.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos)
}

func TestRequestCommandWithBrief(t *testing.T) {
	fixture := setupCommandTest("get", server.URL, "--brief")

	err := fixture.cmd.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos)
}

func TestRequestCommandWithSilent(t *testing.T) {
	fixture := setupCommandTest("get", server.URL, "--silent")

	err := fixture.cmd.Execute()
	require.NoError(t, err)
	require.Empty(t, fixture.infos)
}

func TestRequestGetSigned(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAKIAKAI")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "abcd//efgh/ijklmnopq//bca")
	os.Setenv("AWS_SESSION_TOKEN", "9bd58de0-20ab-4f29-bbd9-dedc700152e3")
	fixture := setupCommandTest("get", server.URL, "--aws-sigv4")

	err := fixture.cmd.Execute()
	require.NoError(t, err)
	require.NotEmpty(t, fixture.infos)
	require.Empty(t, fixture.errs)
}
