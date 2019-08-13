package rest_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lunjon/httpreq/internal/rest"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
	server *httptest.Server
	url    string
}

func (suite *ServerTestSuite) SetupSuite() {
	server := httptest.NewServer(&handler{})
	suite.server = server
	suite.url = server.URL
}

func (suite *ServerTestSuite) SetupTest() {
}

func (suite *ServerTestSuite) TearDownSuite() {
	suite.server.Close()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) TestGet() {
	req, err := rest.BuildRequest(http.MethodGet, suite.url, nil, nil)
	suite.NoError(err)
	suite.NotNil(req)

	res := rest.SendRequest(req)
	suite.NotNil(res)
	suite.True(res.Successful())
	suite.False(res.HasError())
	suite.Nil(res.Error())
}

func (suite *ServerTestSuite) TestSignedEnv() {
	req, err := rest.BuildRequest(http.MethodGet, suite.url, nil, nil)
	suite.NoError(err)
	suite.NotNil(req)

	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAIX2HKAKAEGOTTLOL")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "/HAHAOKAYdB914ByyI+F9/LOLvafanKETCHUP")
	err = rest.SignRequest(req, nil, "eu-west-1", "")
	suite.NoError(err)

	res := rest.SendRequest(req)
	suite.True(res.Successful())
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
