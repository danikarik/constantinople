package app

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AppSuite struct {
	suite.Suite
	app *App
}

func NewAppSuite() (*AppSuite, error) {
	viper.SetEnvPrefix("TEST")
	viper.AutomaticEnv()
	viper.SetDefault("ADDR", ":3000")
	viper.SetDefault("AUTH_SERV", "127.0.0.1:8000")
	viper.SetDefault("REDIS_HOST", "127.0.0.1:6379")
	viper.SetDefault("REDIS_PASS", "")
	viper.SetDefault("ORIGINS", []string{"*"})
	viper.SetDefault("USERNAME", "carbase")
	viper.SetDefault("PASSWORD", "FtgVCxZVZeRk5K9S")

	app, err := New(viper.GetString("ADDR"), Options{
		Origins:     viper.GetStringSlice("ORIGINS"),
		AuthService: viper.GetString("AUTH_SERV"),
		RedisHost:   viper.GetString("REDIS_HOST"),
		RedisPass:   viper.GetString("REDIS_PASS"),
		Username:    viper.GetString("USERNAME"),
		Password:    viper.GetString("PASSWORD"),
	})
	if err != nil {
		return nil, err
	}
	return &AppSuite{
		app: app,
	}, nil
}

func (s *AppSuite) TestPprof() {
	ts := httptest.NewServer(s.app.mux)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/debug", nil)
	if err != nil {
		s.Fail(err.Error())
	}
	req.SetBasicAuth(
		viper.GetString("USERNAME"),
		viper.GetString("PASSWORD"),
	)

	resp, _ := testRequest(req, s.T(), ts)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	req.SetBasicAuth(
		viper.GetString("USERNAME"),
		"debugfailed",
	)
	resp, _ = testRequest(req, s.T(), ts)
	assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (s *AppSuite) TestMetrics() {
	ts := httptest.NewServer(s.app.mux)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/metrics", nil)
	if err != nil {
		s.Fail(err.Error())
	}
	req.SetBasicAuth(
		viper.GetString("USERNAME"),
		viper.GetString("PASSWORD"),
	)

	resp, _ := testRequest(req, s.T(), ts)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	req.SetBasicAuth(
		viper.GetString("USERNAME"),
		"metricfailed",
	)
	resp, _ = testRequest(req, s.T(), ts)
	assert.Equal(s.T(), http.StatusUnauthorized, resp.StatusCode)
}

func TestAppSuite(t *testing.T) {
	as, err := NewAppSuite()
	if assert.NoError(t, err) {
		suite.Run(t, as)
	}
}

func testRequest(req *http.Request, t *testing.T, ts *httptest.Server) (*http.Response, string) {

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
