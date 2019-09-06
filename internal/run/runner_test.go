package run_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lunjon/httpreq/internal/rest"
	"github.com/lunjon/httpreq/internal/run"
	"github.com/stretchr/testify/suite"
)

type RunnerTestSuite struct {
	suite.Suite
	server *httptest.Server
	runner *run.Runner
	client *rest.Client
}

func (suite *RunnerTestSuite) SetupSuite() {
	h := HTTPTestHandler{}
	server := httptest.NewServer(h)
	client := rest.NewClient(server.Client())
	suite.server = server
	suite.client = client
}

func (suite *RunnerTestSuite) SetupTest() {
	runner, err := run.Load("testdata/runner_test.yaml")
	suite.NoError(err)

	// Change each request URL to the test server
	for _, req := range runner.Spec.Requests {
		req.URL = suite.server.URL
	}

	suite.runner = runner
}

func (suite *RunnerTestSuite) TearDownSuite() {
	suite.server.Close()
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(RunnerTestSuite))
}

func (suite *RunnerTestSuite) TestRunner_Run_all() {
	results, err := suite.runner.Run(suite.client)
	suite.NoError(err)
	suite.Len(results, 2)
}

func (suite *RunnerTestSuite) TestRunner_Run_Target1() {
	results, err := suite.runner.Run(suite.client, "target1")
	suite.NoError(err)
	suite.Len(results, 1)
}

func (suite *RunnerTestSuite) TestRunner_Run_Target2() {
	results, err := suite.runner.Run(suite.client, "target2")
	suite.NoError(err)
	suite.Len(results, 1)
}

// An HTTP test handler that always return status 200.
type HTTPTestHandler struct{}

func (h HTTPTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
