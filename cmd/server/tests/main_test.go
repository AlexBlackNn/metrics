package tests

import (
	"github.com/AlexBlackNn/metrics/app/server"
	"github.com/AlexBlackNn/metrics/cmd/server/router"
	"github.com/AlexBlackNn/metrics/internal/config/configserver"
	"github.com/AlexBlackNn/metrics/internal/logger"
	"github.com/AlexBlackNn/metrics/pkg/storage/memstorage"
	"github.com/AlexBlackNn/metrics/pkg/storage/postgres"
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
	cfg := &configserver.Config{
		Env:                   "local",
		ServerAddr:            ":8080",
		ServerReadTimeout:     10,
		ServerWriteTimeout:    10,
		ServerIdleTimeout:     10,
		ServerStoreInterval:   2,
		ServerFileStoragePath: "/tmp/metrics-db.json",
		ServerRestore:         true,
		ServerRateLimit:       10000,
		ServerDataBaseDSN:     "postgresql://postgres:postgres@127.0.0.1:5432/postgres",
	}
	log := logger.New(cfg.Env)

	memStorage, _ := memstorage.New(cfg, log)
	postgresStorage, err := postgres.New(cfg, log)
	ms.Suite.NoError(err)

	ms.application, err = server.NewAppInitStorage(memStorage, postgresStorage, cfg, log)
	ms.Suite.NoError(err)
	ms.client = http.Client{Timeout: 3 * time.Second}
}

func (ms *MetricsSuite) BeforeTest(suiteName, testName string) {
	// Starts server with first random port.
	ms.srv = httptest.NewServer(router.NewChiRouter(
		ms.application.Cfg,
		ms.application.Log,
		ms.application.HandlersV1,
		ms.application.HandlersV2,
		ms.application.HealthHandlersV2,
		ms.application.HandlersV3,
	))
}

func (ms *MetricsSuite) AfterTest(suiteName, testName string) {
	ms.srv = nil
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(MetricsSuite))
}
