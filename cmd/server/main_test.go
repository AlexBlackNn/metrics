package main

import (
	"github.com/AlexBlackNn/metrics/app/server"
	"github.com/AlexBlackNn/metrics/cmd/server/router"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MetricsSuite struct {
	suite.Suite
	application *server.App
	client      http.Client
	srv         *httptest.Server
}

func (ms *MetricsSuite) SetupSuite() {
	var err error
	ms.application, err = server.New()
	if err != nil {
		ms.T().Fatal(err)
	}
	ms.client = http.Client{Timeout: 3 * time.Second}
}

func (ms *MetricsSuite) BeforeTest(suiteName, testName string) {
	// Starts server with first random port.
	ms.srv = httptest.NewServer(router.NewChiRouter(ms.application.Cfg, ms.application.Log, ms.application.HandlersV1, ms.application.HandlersV2))
}

func (ms *MetricsSuite) AfterTest(suiteName, testName string) {
	ms.srv = nil
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(MetricsSuite))
}
